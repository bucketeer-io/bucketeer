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
	"context"
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

		searcher := NewGitHubSearcher(zap.NewNop())
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		docs, err := searcher.Search(context.Background(), "segments", 2)
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

		searcher := NewGitHubSearcher(zap.NewNop())
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		docs, err := searcher.Search(context.Background(), "SDK", 3)
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

		searcher := NewGitHubSearcher(zap.NewNop())
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		docs, err := searcher.Search(context.Background(), "testing", 2)
		require.NoError(t, err)
		assert.Len(t, docs, 2)
	})

	t.Run("returns nil for empty query", func(t *testing.T) {
		t.Parallel()
		searcher := NewGitHubSearcher(zap.NewNop())
		docs, err := searcher.Search(context.Background(), "", 3)
		assert.NoError(t, err)
		assert.Nil(t, docs)
	})

	t.Run("returns nil for zero topK", func(t *testing.T) {
		t.Parallel()
		searcher := NewGitHubSearcher(zap.NewNop())
		docs, err := searcher.Search(context.Background(), "test", 0)
		assert.NoError(t, err)
		assert.Nil(t, docs)
	})

	t.Run("graceful degradation on Trees API error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer server.Close()

		searcher := NewGitHubSearcher(zap.NewNop())
		searcher.apiBaseURL = server.URL

		docs, err := searcher.Search(context.Background(), "test", 3)
		assert.Error(t, err)
		assert.Nil(t, docs)
	})

	t.Run("extracts category from path", func(t *testing.T) {
		t.Parallel()

		server := newTestServer(t, map[string]string{
			"docs/sdk/server-side/go/index.md": "---\ntitle: Go SDK\n---\n\n# Go SDK\n\nGo SDK content.",
		})
		defer server.Close()

		searcher := NewGitHubSearcher(zap.NewNop())
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		docs, err := searcher.Search(context.Background(), "go sdk", 1)
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

		searcher := NewGitHubSearcher(zap.NewNop())
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		docs, err := searcher.Search(context.Background(), "valid", 10)
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

		searcher := NewGitHubSearcher(zap.NewNop())
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		// First call should fetch the tree
		_, err := searcher.Search(context.Background(), "test", 1)
		require.NoError(t, err)
		assert.Equal(t, 1, treeCallCount)

		// Second call should use cache
		_, err = searcher.Search(context.Background(), "test", 1)
		require.NoError(t, err)
		assert.Equal(t, 1, treeCallCount)
	})

	t.Run("returns empty for no matching query", func(t *testing.T) {
		t.Parallel()

		server := newTestServer(t, map[string]string{
			"docs/feature-flags/segments.mdx": "---\ntitle: Segments\n---\n\n# Segments\n\nSegments info.",
		})
		defer server.Close()

		searcher := NewGitHubSearcher(zap.NewNop())
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		docs, err := searcher.Search(context.Background(), "zzzznonexistent", 3)
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

		searcher := NewGitHubSearcher(zap.NewNop())
		searcher.apiBaseURL = server.URL
		searcher.rawBaseURL = server.URL

		docs, err := searcher.Search(context.Background(), "successfully", 3)
		require.NoError(t, err)
		assert.Len(t, docs, 1)
		assert.Equal(t, "OK", docs[0].Metadata.Title)
	})
}

func TestTokenizeQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		query    string
		expected []string
	}{
		{
			name:     "simple words",
			query:    "feature flags",
			expected: []string{"feature", "flags"},
		},
		{
			name:     "removes stop words",
			query:    "how to use the sdk",
			expected: []string{"use", "sdk"},
		},
		{
			name:     "lowercases",
			query:    "Go SDK Android",
			expected: []string{"go", "sdk", "android"},
		},
		{
			name:     "handles empty",
			query:    "",
			expected: nil,
		},
		{
			name:     "all stop words returns empty",
			query:    "how is the",
			expected: nil,
		},
		{
			name:     "japanese katakana separated and translated",
			query:    "セグメント 使い方",
			expected: []string{"セグメント", "segment", "使い方"},
		},
		{
			name:     "extracts ASCII and splits katakana from Japanese",
			query:    "SDKについて調べて",
			expected: []string{"sdk", "について調べて"},
		},
		{
			name:     "mixed English and Japanese with spaces",
			query:    "Go SDKの使い方",
			expected: []string{"go", "sdk", "の使い方"},
		},
		{
			name:     "katakana loanwords translated to English",
			query:    "タグでフラグを整理する",
			expected: []string{"タグ", "tag", "フラグ", "flag", "を整理する"},
		},
		{
			name:     "katakana segment translated",
			query:    "セグメントの使い方",
			expected: []string{"セグメント", "segment", "の使い方"},
		},
		{
			name:     "deduplicates tokens",
			query:    "sdk SDK",
			expected: []string{"sdk"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tokenizeQuery(tt.query)
			assert.Equal(t, tt.expected, result)
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

	t.Run("high score for path match", func(t *testing.T) {
		t.Parallel()
		score := scoreDoc(doc, []string{"segments"})
		assert.Greater(t, score, 0.0)
	})

	t.Run("higher score for multiple matches", func(t *testing.T) {
		t.Parallel()
		singleScore := scoreDoc(doc, []string{"segments"})
		multiScore := scoreDoc(doc, []string{"segments", "targeting"})
		assert.Greater(t, multiScore, singleScore)
	})

	t.Run("zero score for no match", func(t *testing.T) {
		t.Parallel()
		score := scoreDoc(doc, []string{"zzzznonexistent"})
		assert.Equal(t, 0.0, score)
	})

	t.Run("title match scores higher than content-only match", func(t *testing.T) {
		t.Parallel()
		titleScore := scoreDoc(doc, []string{"segments"})     // matches path + title + content
		contentScore := scoreDoc(doc, []string{"attributes"}) // matches content only
		assert.Greater(t, titleScore, contentScore)
	})
}

func TestExtractCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path     string
		expected string
	}{
		{"docs/feature-flags/segments.mdx", "feature-flags"},
		{"docs/sdk/server-side/go/index.md", "sdk/server-side/go"},
		{"docs/index.mdx", ""},
		{"docs/getting-started/quickstart.mdx", "getting-started"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, extractCategory(tt.path))
		})
	}
}

func TestDocsSiteURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path     string
		expected string
	}{
		{"docs/feature-flags/segments.mdx", "https://docs.bucketeer.io/feature-flags/segments"},
		{"docs/sdk/server-side/go/index.md", "https://docs.bucketeer.io/sdk/server-side/go"},
		{"docs/index.mdx", "https://docs.bucketeer.io"},
		{"docs/getting-started/quickstart.mdx", "https://docs.bucketeer.io/getting-started/quickstart"},
		{"docs/best-practices/optimize-bucketeer-with-tags.mdx", "https://docs.bucketeer.io/best-practices/optimize-bucketeer-with-tags"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, docsSiteURL(tt.path))
		})
	}
}

func TestIsValidDocPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path  string
		valid bool
	}{
		{"docs/feature-flags/segments.mdx", true},
		{"docs/sdk/go/index.md", true},
		{"docs/../etc/passwd", false},
		{"not-docs/file.mdx", false},
		{"docs/file.txt", false},
		{"docs/file.json", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.valid, isValidDocPath(tt.path))
		})
	}
}
