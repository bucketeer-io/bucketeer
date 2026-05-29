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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// DocChunk represents a document chunk returned by a search.
type DocChunk struct {
	ID       string  `json:"id"`
	Content  string  `json:"content"`
	Metadata DocMeta `json:"metadata"`
}

// DocMeta contains metadata about a document chunk.
type DocMeta struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Category string `json:"category"`
}

// Searcher defines the interface for document search.
type Searcher interface {
	Search(ctx context.Context, query string, topK int) ([]DocChunk, error)
}

const (
	defaultAPIBaseURL    = "https://api.github.com"
	defaultRawBaseURL    = "https://raw.githubusercontent.com"
	docsRepo             = "bucketeer-io/bucketeer-docs"
	maxContentLength     = 2000
	maxIndexContent      = 5000 // content stored in index for scoring (longer than output)
	maxRawBodyBytes      = 512 * 1024
	maxTreeResponseBytes = 10 * 1024 * 1024
	maxTopK              = 10
	maxConcurrentFetch   = 5
	defaultCacheTTL      = 24 * time.Hour
)

// gitTreeResponse is the GitHub Trees API response.
type gitTreeResponse struct {
	SHA  string         `json:"sha"`
	Tree []gitTreeEntry `json:"tree"`
}

// gitTreeEntry is a single entry in the tree.
type gitTreeEntry struct {
	Path string `json:"path"`
	Type string `json:"type"` // "blob" or "tree"
}

// indexedDoc is an in-memory cached document used for local scoring.
type indexedDoc struct {
	path     string
	title    string
	content  string // full content for scoring (up to maxIndexContent)
	htmlURL  string
	category string
	// Pre-computed fields for scoring (avoid repeated work in hot path)
	lowerTitle   string
	lowerContent string
	pathSegments []string
}

// GitHubSearcher searches Bucketeer documentation by fetching all docs
// from the GitHub Trees API and scoring them locally against the query.
// No authentication is required for public repositories.
type GitHubSearcher struct {
	httpClient  *http.Client
	apiBaseURL  string // GitHub API base URL (for Trees API)
	rawBaseURL  string // raw.githubusercontent.com base URL
	githubToken string // optional GitHub token to increase rate limits
	logger      *zap.Logger
	cacheTTL    time.Duration

	mu          sync.RWMutex
	docIndex    []indexedDoc
	lastRefresh time.Time
	refreshing  int32 // atomic flag to prevent concurrent refreshes
}

// NewGitHubSearcher creates a new GitHubSearcher.
func NewGitHubSearcher(logger *zap.Logger, githubToken string) *GitHubSearcher {
	return &GitHubSearcher{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiBaseURL:  defaultAPIBaseURL,
		rawBaseURL:  defaultRawBaseURL,
		githubToken: githubToken,
		cacheTTL:    defaultCacheTTL,
		logger:      logger.Named("github-searcher"),
	}
}

// setAuthHeader adds an Authorization header if a GitHub token is configured.
func (g *GitHubSearcher) setAuthHeader(req *http.Request) {
	if g.githubToken != "" {
		req.Header.Set("Authorization", "Bearer "+g.githubToken)
	}
}

// Search finds relevant documentation by scoring cached docs against the query.
func (g *GitHubSearcher) Search(ctx context.Context, query string, topK int) ([]DocChunk, error) {
	if query == "" || topK <= 0 {
		return nil, nil
	}
	if topK > maxTopK {
		topK = maxTopK
	}

	index, err := g.ensureIndex(ctx)
	if err != nil {
		return nil, err
	}

	tokens := tokenizeQuery(query)
	if len(tokens) == 0 {
		return nil, nil
	}
	tokens = expandQueryTokens(tokens)

	// Score all indexed docs
	type scored struct {
		doc   indexedDoc
		score float64
	}
	results := make([]scored, 0, len(index))
	for _, d := range index {
		s := scoreDoc(d, tokens)
		if s > 0 {
			results = append(results, scored{doc: d, score: s})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	if topK < len(results) {
		results = results[:topK]
	}

	chunks := make([]DocChunk, len(results))
	for i, r := range results {
		content := r.doc.content
		if utf8.RuneCountInString(content) > maxContentLength {
			content = string([]rune(content)[:maxContentLength])
		}
		chunks[i] = DocChunk{
			ID:      r.doc.path,
			Content: content,
			Metadata: DocMeta{
				Title:    r.doc.title,
				URL:      r.doc.htmlURL,
				Category: r.doc.category,
			},
		}
	}
	return chunks, nil
}

// ensureIndex returns the cached document index, refreshing if needed.
func (g *GitHubSearcher) ensureIndex(ctx context.Context) ([]indexedDoc, error) {
	g.mu.RLock()
	if len(g.docIndex) > 0 && time.Since(g.lastRefresh) < g.cacheTTL {
		defer g.mu.RUnlock()
		return g.docIndex, nil
	}
	hasCache := len(g.docIndex) > 0
	cached := g.docIndex
	g.mu.RUnlock()

	// First time: blocking refresh
	if !hasCache {
		return g.refreshIndex(ctx)
	}

	// Subsequent: background refresh, return stale cache
	if atomic.CompareAndSwapInt32(&g.refreshing, 0, 1) {
		go func() {
			defer atomic.StoreInt32(&g.refreshing, 0)
			if _, err := g.refreshIndex(context.Background()); err != nil {
				g.logger.Warn("Background index refresh failed", zap.Error(err))
			}
		}()
	}
	return cached, nil
}

// refreshIndex fetches the full document tree and builds the in-memory index.
func (g *GitHubSearcher) refreshIndex(ctx context.Context) ([]indexedDoc, error) {
	paths, err := g.fetchTree(ctx)
	if err != nil {
		return nil, err
	}

	index := g.fetchAllDocs(ctx, paths)

	g.mu.Lock()
	g.docIndex = index
	g.lastRefresh = time.Now()
	g.mu.Unlock()

	g.logger.Info("Document index refreshed", zap.Int("count", len(index)))
	return index, nil
}

// fetchTree retrieves all doc file paths from the GitHub Trees API.
func (g *GitHubSearcher) fetchTree(ctx context.Context) ([]string, error) {
	u, err := url.JoinPath(g.apiBaseURL, "repos", docsRepo, "git/trees/main")
	if err != nil {
		return nil, err
	}
	u += "?recursive=1"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	g.setAuthHeader(req)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github trees API returned status %d", resp.StatusCode)
	}

	var treeResp gitTreeResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxTreeResponseBytes)).Decode(&treeResp); err != nil {
		return nil, err
	}

	var paths []string
	for _, entry := range treeResp.Tree {
		if entry.Type == "blob" && isValidDocPath(entry.Path) {
			paths = append(paths, entry.Path)
		}
	}

	return paths, nil
}

// fetchAllDocs fetches raw content for all paths concurrently.
func (g *GitHubSearcher) fetchAllDocs(ctx context.Context, paths []string) []indexedDoc {
	var mu sync.Mutex
	docs := make([]indexedDoc, 0, len(paths))

	eg, gCtx := errgroup.WithContext(ctx)
	eg.SetLimit(maxConcurrentFetch)

	for _, p := range paths {
		eg.Go(func() error {
			doc, err := g.fetchRawDoc(gCtx, p)
			if err != nil {
				g.logger.Warn("Failed to fetch doc", zap.Error(err))
				return nil
			}
			mu.Lock()
			docs = append(docs, doc)
			mu.Unlock()
			return nil
		})
	}
	_ = eg.Wait()
	return docs
}

// fetchRawDoc fetches and processes a single document from raw.githubusercontent.com.
func (g *GitHubSearcher) fetchRawDoc(ctx context.Context, docPath string) (indexedDoc, error) {
	rawURL, err := url.JoinPath(g.rawBaseURL, docsRepo, "main", docPath)
	if err != nil {
		return indexedDoc{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return indexedDoc{}, err
	}
	g.setAuthHeader(req)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return indexedDoc{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return indexedDoc{}, fmt.Errorf("raw fetch returned status %d for %s", resp.StatusCode, docPath)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxRawBodyBytes))
	if err != nil {
		return indexedDoc{}, err
	}

	rawContent := string(body)
	title := ExtractTitle(rawContent)
	content := StripMDX(rawContent)

	// Keep more content in index for scoring, truncate only on output
	if utf8.RuneCountInString(content) > maxIndexContent {
		content = string([]rune(content)[:maxIndexContent])
	}

	htmlURL := docsSiteURL(docPath)

	return indexedDoc{
		path:         docPath,
		title:        title,
		content:      content,
		htmlURL:      htmlURL,
		category:     extractCategory(docPath),
		lowerTitle:   strings.ToLower(title),
		lowerContent: strings.ToLower(content),
		pathSegments: extractPathSegments(docPath),
	}, nil
}

// querySynonyms maps a token to additional tokens that should also be searched.
// Bucketeer's docs use canonical terminology (e.g. "experiment") that doesn't always
// match how users phrase questions (e.g. "A/B test"). Expansion bridges that gap
// without requiring a semantic search backend.
//
// Keep this list narrow: only terms that are genuinely interchangeable in the docs
// domain. Avoid generic words like "test" — they would boost too many docs.
var querySynonyms = map[string][]string{
	// A/B testing -> experiments
	"a/b":         {"experiment", "experiments"},
	"ab":          {"experiment", "experiments"},
	"split":       {"experiment", "experiments"},
	"experiment":  {"a/b"},
	"experiments": {"a/b"},
	// Progressive / automated rollout -> auto-operation
	"progressive": {"auto-operation", "rollout"},
	"automated":   {"auto-operation"},
	"autoops":     {"auto-operation"},
	"automation":  {"auto-operation"},
	// Variant terminology -> variations
	"variant":  {"variations"},
	"variants": {"variations"},
	// Activity -> audit logs / history
	"activity": {"audit-logs", "history"},
	// Webhooks -> trigger (Bucketeer's flag-trigger feature)
	"webhook":  {"trigger"},
	"webhooks": {"trigger"},
}

// expandQueryTokens appends synonym tokens for any token in the input that has
// entries in querySynonyms. Existing tokens are preserved; duplicates are skipped.
func expandQueryTokens(tokens []string) []string {
	seen := make(map[string]bool, len(tokens))
	for _, t := range tokens {
		seen[t] = true
	}
	expanded := tokens
	for _, t := range tokens {
		for _, syn := range querySynonyms[t] {
			if !seen[syn] {
				seen[syn] = true
				expanded = append(expanded, syn)
			}
		}
	}
	return expanded
}

// tokenizeQuery splits a query into unique lowercase tokens,
// stripping punctuation from each token.
func tokenizeQuery(query string) []string {
	seen := make(map[string]bool)
	var tokens []string
	for _, w := range strings.Fields(strings.ToLower(query)) {
		// Strip leading/trailing punctuation
		w = strings.Trim(w, ".,;:!?\"'()[]{}/*-+=#@&")
		if w != "" && !seen[w] {
			seen[w] = true
			tokens = append(tokens, w)
		}
	}
	if len(tokens) == 0 {
		return nil
	}
	return tokens
}

// scoreDoc computes a relevance score for a document against query tokens.
// Uses pre-computed lowercase fields and path segments from indexedDoc.
func scoreDoc(doc indexedDoc, queryTokens []string) float64 {
	var score float64

	for _, token := range queryTokens {
		// Skip single-character tokens to reduce noise
		if utf8.RuneCountInString(token) <= 1 {
			continue
		}
		// Path segment match (highest weight — structural relevance)
		// Use HasPrefix for reverse direction to handle plurals (e.g. "sdks" has prefix "sdk")
		for _, seg := range doc.pathSegments {
			if seg == token {
				score += 10.0
			} else if strings.Contains(seg, token) || strings.HasPrefix(token, seg) {
				score += 5.0
			}
		}
		// Title match
		if strings.Contains(doc.lowerTitle, token) {
			score += 3.0
		}
		// Content match (low weight — presence only, not frequency)
		if strings.Contains(doc.lowerContent, token) {
			score += 0.5
		}
	}
	return score
}

// extractPathSegments returns meaningful segments from a doc path for scoring.
// e.g. "docs/feature-flags/segments.mdx" -> ["feature-flags", "segments"]
func extractPathSegments(p string) []string {
	trimmed := strings.TrimPrefix(p, "docs/")
	parts := strings.Split(trimmed, "/")
	segments := make([]string, 0, len(parts))
	for _, part := range parts {
		// Remove file extension
		if idx := strings.LastIndex(part, "."); idx >= 0 {
			part = part[:idx]
		}
		// Skip "index" as it's not meaningful
		if part != "" && part != "index" {
			segments = append(segments, strings.ToLower(part))
		}
	}
	return segments
}

// isValidDocPath validates that the path is a safe documentation file path.
func isValidDocPath(p string) bool {
	if strings.Contains(p, "..") {
		return false
	}
	return strings.HasPrefix(p, "docs/") && (strings.HasSuffix(p, ".mdx") || strings.HasSuffix(p, ".md"))
}

// docsSiteURL converts a repository doc path to a published documentation URL.
// e.g. "docs/feature-flags/segments.mdx" -> "https://docs.bucketeer.io/feature-flags/segments"
// e.g. "docs/sdk/server-side/go/index.md" -> "https://docs.bucketeer.io/sdk/server-side/go"
func docsSiteURL(docPath string) string {
	const docsBaseURL = "https://docs.bucketeer.io"
	trimmed := strings.TrimPrefix(docPath, "docs/")
	// Remove file extension (.md / .mdx)
	if idx := strings.LastIndex(trimmed, "."); idx >= 0 {
		trimmed = trimmed[:idx]
	}
	// Remove trailing /index (directory index pages)
	trimmed = strings.TrimSuffix(trimmed, "/index")
	// Handle root index
	if trimmed == "index" {
		return docsBaseURL
	}
	return docsBaseURL + "/" + trimmed
}

// extractCategory extracts the category from a docs path like "docs/feature-flags/segments.mdx".
func extractCategory(p string) string {
	trimmed := strings.TrimPrefix(p, "docs/")
	dir := path.Dir(trimmed)
	if dir == "." {
		return ""
	}
	return dir
}
