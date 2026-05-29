// Copyright 2026 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rag

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// newTestServer creates an httptest server that serves a Trees API response
// and raw file contents for the given docs map (path -> content).
func newTestServer(t *testing.T, docs map[string]string) *httptest.Server {
	t.Helper()

	treeEntries := make([]gitTreeEntry, 0, len(docs))
	for p := range docs {
		treeEntries = append(treeEntries, gitTreeEntry{Path: p, Type: "blob"})
	}
	// Add some non-doc entries to verify filtering
	treeEntries = append(treeEntries,
		gitTreeEntry{Path: "README.md", Type: "blob"},
		gitTreeEntry{Path: "src", Type: "tree"},
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/repos/bucketeer-io/bucketeer-docs/git/trees/main", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(gitTreeResponse{
			SHA:  "abc123",
			Tree: treeEntries,
		})
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Raw content paths: /bucketeer-io/bucketeer-docs/main/{path}
		for p, content := range docs {
			if r.URL.Path == "/bucketeer-io/bucketeer-docs/main/"+p {
				w.Write([]byte(content))
				return
			}
		}
		http.NotFound(w, r)
	})

	return httptest.NewServer(mux)
}

func TestGitHubSearcherSearch(t *testing.T) {
	t.Parallel()

	t.Run("returns relevant docs by keyword match", func(t *testing.T) {
		t.Parallel()

		server := newTestServer(t, map[string]string{
			"docs/feature-flags/segments.mdx":       "---\ntitle: Segments\n---\n\n# Segments\n\nSegments allow you to group users.",
			"docs/getting-started/quickstart.mdx":   "---\ntitle: Quickstart\n---\n\n# Quickstart\n\nGet started with Bucketeer.",
			"docs/sdk/server-side/go/index.md":      "---\ntitle: Go SDK\n---\n\n# Go SDK\n\nInstall the Go SDK.",
			"docs/sdk/client-side/android/index.md": "---\ntitle: Android SDK\n---\n\n# Android SDK\n\nInstall the Android SDK.",
		})
		defer server.Close()

		searcher := NewGitHubSearcher(zap.NewNop(), "")
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		docs, err := searcher.Search(t.Context(), "segments", 2)
		require.NoError(t, err)
		require.NotEmpty(t, docs)

		// "segments" should match the segments doc (path + title + content match)
		assert.Equal(t, "Segments", docs[0].Metadata.Title)
		assert.Contains(t, docs[0].Content, "Segments allow you to group users")
		assert.Equal(t, "feature-flags", docs[0].Metadata.Category)
	})

	t.Run("returns SDK docs when querying sdk", func(t *testing.T) {
		t.Parallel()

		server := newTestServer(t, map[string]string{
			"docs/feature-flags/segments.mdx":       "---\ntitle: Segments\n---\n\n# Segments\n\nSegments info.",
			"docs/sdk/server-side/go/index.md":      "---\ntitle: Go SDK\n---\n\n# Go SDK\n\nInstall the Go SDK for server-side evaluation.",
			"docs/sdk/client-side/android/index.md": "---\ntitle: Android SDK\n---\n\n# Android SDK\n\nInstall the Android SDK for mobile apps.",
			"docs/sdk/index.mdx":                    "---\ntitle: SDKs Overview\n---\n\n# SDKs\n\nBucketeer provides client-side and server-side SDKs.",
		})
		defer server.Close()

		searcher := NewGitHubSearcher(zap.NewNop(), "")
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		docs, err := searcher.Search(t.Context(), "SDK", 3)
		require.NoError(t, err)
		require.NotEmpty(t, docs)

		// All returned docs should be SDK-related
		for _, doc := range docs {
			assert.Contains(t, doc.Content+doc.Metadata.Title, "SDK",
				"expected SDK-related doc but got: %s", doc.Metadata.Title)
		}
	})

	t.Run("respects topK limit", func(t *testing.T) {
		t.Parallel()

		server := newTestServer(t, map[string]string{
			"docs/a/one.mdx":   "# One\nContent about testing one.",
			"docs/b/two.mdx":   "# Two\nContent about testing two.",
			"docs/c/three.mdx": "# Three\nContent about testing three.",
		})
		defer server.Close()

		searcher := NewGitHubSearcher(zap.NewNop(), "")
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		docs, err := searcher.Search(t.Context(), "testing", 2)
		require.NoError(t, err)
		assert.Len(t, docs, 2)
	})

	t.Run("returns nil for empty query", func(t *testing.T) {
		t.Parallel()
		searcher := NewGitHubSearcher(zap.NewNop(), "")
		docs, err := searcher.Search(t.Context(), "", 3)
		assert.NoError(t, err)
		assert.Nil(t, docs)
	})

	t.Run("returns nil for zero topK", func(t *testing.T) {
		t.Parallel()
		searcher := NewGitHubSearcher(zap.NewNop(), "")
		docs, err := searcher.Search(t.Context(), "test", 0)
		assert.NoError(t, err)
		assert.Nil(t, docs)
	})

	t.Run("graceful degradation on Trees API error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer server.Close()

		searcher := NewGitHubSearcher(zap.NewNop(), "")
		searcher.apiBaseURL = server.URL

		docs, err := searcher.Search(t.Context(), "test", 3)
		assert.Error(t, err)
		assert.Nil(t, docs)
	})

	t.Run("extracts category from path", func(t *testing.T) {
		t.Parallel()

		server := newTestServer(t, map[string]string{
			"docs/sdk/server-side/go/index.md": "---\ntitle: Go SDK\n---\n\n# Go SDK\n\nGo SDK content.",
		})
		defer server.Close()

		searcher := NewGitHubSearcher(zap.NewNop(), "")
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		docs, err := searcher.Search(t.Context(), "go sdk", 1)
		require.NoError(t, err)
		require.Len(t, docs, 1)
		assert.Equal(t, "sdk/server-side/go", docs[0].Metadata.Category)
	})

	t.Run("skips invalid paths in tree", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		mux.HandleFunc("/repos/bucketeer-io/bucketeer-docs/git/trees/main", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(gitTreeResponse{
				SHA: "abc",
				Tree: []gitTreeEntry{
					{Path: "docs/../etc/passwd", Type: "blob"},
					{Path: "docs/valid.mdx", Type: "blob"},
					{Path: "not-docs/file.mdx", Type: "blob"},
				},
			})
		})
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("# Valid\nSome valid content for testing."))
		})
		server := httptest.NewServer(mux)
		defer server.Close()

		searcher := NewGitHubSearcher(zap.NewNop(), "")
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		docs, err := searcher.Search(t.Context(), "valid", 10)
		require.NoError(t, err)
		// Only docs/valid.mdx should be indexed
		assert.Len(t, docs, 1)
	})

	t.Run("uses cached index on subsequent calls", func(t *testing.T) {
		t.Parallel()

		treeCallCount := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/repos/bucketeer-io/bucketeer-docs/git/trees/main", func(w http.ResponseWriter, r *http.Request) {
			treeCallCount++
			json.NewEncoder(w).Encode(gitTreeResponse{
				SHA: "abc",
				Tree: []gitTreeEntry{
					{Path: "docs/test.mdx", Type: "blob"},
				},
			})
		})
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("# Test\nTest content for caching."))
		})
		server := httptest.NewServer(mux)
		defer server.Close()

		searcher := NewGitHubSearcher(zap.NewNop(), "")
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		// First call should fetch the tree
		_, err := searcher.Search(t.Context(), "test", 1)
		require.NoError(t, err)
		assert.Equal(t, 1, treeCallCount)

		// Second call should use cache
		_, err = searcher.Search(t.Context(), "test", 1)
		require.NoError(t, err)
		assert.Equal(t, 1, treeCallCount)
	})

	t.Run("returns empty for no matching query", func(t *testing.T) {
		t.Parallel()

		server := newTestServer(t, map[string]string{
			"docs/feature-flags/segments.mdx": "---\ntitle: Segments\n---\n\n# Segments\n\nSegments info.",
		})
		defer server.Close()

		searcher := NewGitHubSearcher(zap.NewNop(), "")
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		docs, err := searcher.Search(t.Context(), "zzzznonexistent", 3)
		require.NoError(t, err)
		assert.Empty(t, docs)
	})

	t.Run("partial fetch failure still indexes successful docs", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		mux.HandleFunc("/repos/bucketeer-io/bucketeer-docs/git/trees/main", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(gitTreeResponse{
				SHA: "abc",
				Tree: []gitTreeEntry{
					{Path: "docs/ok.mdx", Type: "blob"},
					{Path: "docs/missing.mdx", Type: "blob"},
				},
			})
		})
		mux.HandleFunc("/bucketeer-io/bucketeer-docs/main/docs/ok.mdx", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("# OK\nThis doc loaded successfully."))
		})
		mux.HandleFunc("/bucketeer-io/bucketeer-docs/main/docs/missing.mdx", func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})
		server := httptest.NewServer(mux)
		defer server.Close()

		searcher := NewGitHubSearcher(zap.NewNop(), "")
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		docs, err := searcher.Search(t.Context(), "successfully", 3)
		require.NoError(t, err)
		assert.Len(t, docs, 1)
		assert.Equal(t, "OK", docs[0].Metadata.Title)
	})
}

func TestTokenizeQuery(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		query    string
		expected []string
	}{
		{
			desc:     "simple words",
			query:    "feature flags",
			expected: []string{"feature", "flags"},
		},
		{
			desc:     "lowercases",
			query:    "Go SDK Android",
			expected: []string{"go", "sdk", "android"},
		},
		{
			desc:     "handles empty",
			query:    "",
			expected: nil,
		},
		{
			desc:     "deduplicates tokens",
			query:    "sdk SDK",
			expected: []string{"sdk"},
		},
		{
			desc:     "whitespace only",
			query:    "   ",
			expected: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			result := tokenizeQuery(p.query)
			assert.Equal(t, p.expected, result)
		})
	}
}

func TestExpandQueryTokens(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		tokens   []string
		expected []string
	}{
		{
			desc:     "a/b expands to experiment terms",
			tokens:   []string{"a/b", "test"},
			expected: []string{"a/b", "test", "experiment", "experiments"},
		},
		{
			desc:     "ab expands to experiment terms",
			tokens:   []string{"ab"},
			expected: []string{"ab", "experiment", "experiments"},
		},
		{
			desc:     "experiment expands to a/b",
			tokens:   []string{"experiment"},
			expected: []string{"experiment", "a/b"},
		},
		{
			desc:     "toggle expands to flag",
			tokens:   []string{"toggle"},
			expected: []string{"toggle", "flag", "flags"},
		},
		{
			desc:     "progressive expands to auto-operation and rollout",
			tokens:   []string{"progressive"},
			expected: []string{"progressive", "auto-operation", "rollout"},
		},
		{
			desc:     "automated expands to auto-operation",
			tokens:   []string{"automated"},
			expected: []string{"automated", "auto-operation"},
		},
		{
			desc:     "variant expands to variations",
			tokens:   []string{"variant"},
			expected: []string{"variant", "variations"},
		},
		{
			desc:     "activity expands to audit-logs and history",
			tokens:   []string{"activity"},
			expected: []string{"activity", "audit-logs", "history"},
		},
		{
			desc:     "webhook expands to trigger",
			tokens:   []string{"webhook"},
			expected: []string{"webhook", "trigger"},
		},
		{
			desc:     "no synonyms leaves tokens unchanged",
			tokens:   []string{"sdk", "android"},
			expected: []string{"sdk", "android"},
		},
		{
			desc:     "does not duplicate existing tokens",
			tokens:   []string{"a/b", "experiment"},
			expected: []string{"a/b", "experiment", "experiments"},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			result := expandQueryTokens(p.tokens)
			assert.ElementsMatch(t, p.expected, result)
		})
	}
}

func TestGitHubSearcherSearchQuery(t *testing.T) {
	t.Parallel()

	// Fixture mirrors Bucketeer's actual docs paths. Each doc is filed under
	// canonical terminology that doesn't always match how users phrase questions.
	// The subtests verify expandQueryTokens bridges that gap end-to-end.
	server := newTestServer(t, map[string]string{
		"docs/experimentation/experiments.md":                                  "---\ntitle: Create an Experiment\n---\n\n# Create an Experiment\n\nDefine variations and goals.",
		"docs/feature-flags/segments.mdx":                                      "---\ntitle: Segments\n---\n\n# Segments\n\nSegments allow grouping users.",
		"docs/feature-flags/audit-logs.mdx":                                    "---\ntitle: Audit Logs\n---\n\n# Audit Logs\n\nTrack changes to your flags over time.",
		"docs/feature-flags/creating-feature-flags/auto-operation/rollout.mdx": "---\ntitle: Progressive Rollout\n---\n\n# Auto-operation Rollout\n\nGradually increase traffic to a variation.",
		"docs/feature-flags/creating-feature-flags/variations.mdx":             "---\ntitle: Variations\n---\n\n# Variations\n\nDefine the possible values returned by a flag.",
		"docs/feature-flags/creating-feature-flags/trigger.mdx":                "---\ntitle: Triggers\n---\n\n# Triggers\n\nFire flag operations via webhook URLs.",
	})
	// t.Cleanup (not defer) so the server stays alive for parallel subtests
	// that haven't started yet when this function returns.
	t.Cleanup(server.Close)

	searcher := NewGitHubSearcher(zap.NewNop(), "")
	searcher.apiBaseURL = server.URL
	searcher.rawBaseURL = server.URL

	patterns := []struct {
		desc         string
		query        string
		wantTopTitle string
	}{
		{
			desc:         "A/B test query maps to experiments doc",
			query:        "How do I set up an A/B test?",
			wantTopTitle: "Create an Experiment",
		},
		{
			desc:         "progressive rollout query maps to auto-operation doc",
			query:        "How do I do a progressive rollout?",
			wantTopTitle: "Progressive Rollout",
		},
		{
			desc:         "variant query maps to variations doc",
			query:        "How do I add a variant?",
			wantTopTitle: "Variations",
		},
		{
			desc:         "activity query maps to audit logs doc",
			query:        "Where can I see flag activity?",
			wantTopTitle: "Audit Logs",
		},
		{
			desc:         "webhook query maps to trigger doc",
			query:        "How do I set up a webhook?",
			wantTopTitle: "Triggers",
		},
		{
			desc:         "direct keyword query still works without expansion",
			query:        "segments",
			wantTopTitle: "Segments",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			docs, err := searcher.Search(t.Context(), p.query, 3)
			require.NoError(t, err)
			require.NotEmpty(t, docs, "expected results for query %q", p.query)
			assert.Equal(t, p.wantTopTitle, docs[0].Metadata.Title,
				"for query %q", p.query)
		})
	}
}

func TestScoreDoc(t *testing.T) {
	t.Parallel()

	content := "Segments allow you to group users based on attributes. Use segments for targeting."
	doc := indexedDoc{
		path:         "docs/feature-flags/segments.mdx",
		title:        "Segments",
		content:      content,
		lowerTitle:   strings.ToLower("Segments"),
		lowerContent: strings.ToLower(content),
		pathSegments: extractPathSegments("docs/feature-flags/segments.mdx"),
	}

	patterns := []struct {
		desc      string
		tokens    []string
		checkFunc func(t *testing.T, score float64)
	}{
		{
			desc:   "high score for path match",
			tokens: []string{"segments"},
			checkFunc: func(t *testing.T, score float64) {
				assert.Greater(t, score, 0.0)
			},
		},
		{
			desc:   "zero score for no match",
			tokens: []string{"zzzznonexistent"},
			checkFunc: func(t *testing.T, score float64) {
				assert.Equal(t, 0.0, score)
			},
		},
		{
			desc:   "higher score for multiple matches",
			tokens: []string{"segments", "targeting"},
			checkFunc: func(t *testing.T, score float64) {
				singleScore := scoreDoc(doc, []string{"segments"})
				assert.Greater(t, score, singleScore)
			},
		},
		{
			desc:   "title match scores higher than content-only match",
			tokens: []string{"segments"},
			checkFunc: func(t *testing.T, score float64) {
				contentScore := scoreDoc(doc, []string{"attributes"})
				assert.Greater(t, score, contentScore)
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			score := scoreDoc(doc, p.tokens)
			p.checkFunc(t, score)
		})
	}
}

func TestExtractCategory(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		path     string
		expected string
	}{
		{"feature-flags", "docs/feature-flags/segments.mdx", "feature-flags"},
		{"sdk/server-side/go", "docs/sdk/server-side/go/index.md", "sdk/server-side/go"},
		{"empty for docs root", "docs/index.mdx", ""},
		{"getting-started", "docs/getting-started/quickstart.mdx", "getting-started"},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, p.expected, extractCategory(p.path))
		})
	}
}

func TestDocsSiteURL(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		path     string
		expected string
	}{
		{"segments", "docs/feature-flags/segments.mdx", "https://docs.bucketeer.io/feature-flags/segments"},
		{"go sdk", "docs/sdk/server-side/go/index.md", "https://docs.bucketeer.io/sdk/server-side/go"},
		{"root index", "docs/index.mdx", "https://docs.bucketeer.io"},
		{"quickstart", "docs/getting-started/quickstart.mdx", "https://docs.bucketeer.io/getting-started/quickstart"},
		{"best practices", "docs/best-practices/optimize-bucketeer-with-tags.mdx", "https://docs.bucketeer.io/best-practices/optimize-bucketeer-with-tags"},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, p.expected, docsSiteURL(p.path))
		})
	}
}

func TestIsValidDocPath(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc  string
		path  string
		valid bool
	}{
		{"valid mdx", "docs/feature-flags/segments.mdx", true},
		{"valid md", "docs/sdk/go/index.md", true},
		{"path traversal", "docs/../etc/passwd", false},
		{"not docs dir", "not-docs/file.mdx", false},
		{"txt file", "docs/file.txt", false},
		{"json file", "docs/file.json", false},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, p.valid, isValidDocPath(p.path))
		})
	}
}
