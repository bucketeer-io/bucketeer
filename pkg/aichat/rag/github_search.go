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
	"math"
	"net/http"
	"path"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode"
	"unicode/utf8"

	"go.uber.org/zap"
)

const (
	defaultAPIBaseURL  = "https://api.github.com"
	defaultRawBaseURL  = "https://raw.githubusercontent.com"
	docsRepo           = "bucketeer-io/bucketeer-docs"
	maxContentLength   = 2000
	maxIndexContent    = 5000 // content stored in index for scoring (longer than output)
	maxRawBodyBytes    = 512 * 1024
	maxTopK            = 10
	maxConcurrentFetch = 5
	defaultCacheTTL    = 24 * time.Hour
)

// katakanaToEnglish maps common Katakana loanwords to their English equivalents
// for cross-language search matching against English documentation.
var katakanaToEnglish = map[string]string{
	"タグ":       "tag",
	"フラグ":      "flag",
	"セグメント":    "segment",
	"ターゲティング":  "targeting",
	"ターゲット":    "target",
	"エクスペリメント": "experiment",
	"テスト":      "test",
	"ロールアウト":   "rollout",
	"ユーザー":     "user",
	"バケット":     "bucket",
	"バリエーション":  "variation",
	"プッシュ":     "push",
	"イベント":     "event",
	"ゴール":      "goal",
	"ダッシュボード":  "dashboard",
	"オートオプス":   "autoops",
	"トリガー":     "trigger",
	"ウェブフック":   "webhook",
	"クイックスタート": "quickstart",
	"インストール":   "install",
	"コンソール":    "console",
	"プロジェクト":   "project",
	"エンバイロメント": "environment",
	"オペレーション":  "operation",
	"オペレーションズ": "operations",
	"パフォーマンス":  "performance",
	"アナリティクス":  "analytics",
	"メトリクス":    "metrics",
	"チェンジログ":   "changelog",
	"ドキュメント":   "documentation",
}

// stopWords are common words excluded from search scoring.
var stopWords = map[string]bool{
	"a": true, "an": true, "the": true, "is": true, "are": true,
	"was": true, "were": true, "be": true, "been": true, "being": true,
	"have": true, "has": true, "had": true, "do": true, "does": true,
	"did": true, "will": true, "would": true, "could": true, "should": true,
	"may": true, "might": true, "can": true, "shall": true,
	"to": true, "of": true, "in": true, "for": true, "on": true,
	"with": true, "at": true, "by": true, "from": true, "as": true,
	"into": true, "about": true, "between": true, "through": true,
	"and": true, "but": true, "or": true, "nor": true, "not": true,
	"it": true, "its": true, "this": true, "that": true, "these": true,
	"i": true, "me": true, "my": true, "we": true, "our": true, "you": true, "your": true,
	"how": true, "what": true, "which": true, "who": true, "where": true, "when": true, "why": true,
}

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
	httpClient *http.Client
	apiBaseURL string // GitHub API base URL (for Trees API)
	rawBaseURL string // raw.githubusercontent.com base URL
	logger     *zap.Logger
	cacheTTL   time.Duration

	mu          sync.RWMutex
	docIndex    []indexedDoc
	lastRefresh time.Time
	refreshing  int32 // atomic flag to prevent concurrent refreshes
}

// NewGitHubSearcher creates a new GitHubSearcher.
func NewGitHubSearcher(logger *zap.Logger) *GitHubSearcher {
	return &GitHubSearcher{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiBaseURL: defaultAPIBaseURL,
		rawBaseURL: defaultRawBaseURL,
		cacheTTL:   defaultCacheTTL,
		logger:     logger.Named("github-searcher"),
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
	u := fmt.Sprintf("%s/repos/%s/git/trees/main?recursive=1", g.apiBaseURL, docsRepo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github trees API returned status %d", resp.StatusCode)
	}

	var treeResp gitTreeResponse
	if err := json.NewDecoder(resp.Body).Decode(&treeResp); err != nil {
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

// fetchAllDocs fetches raw content for all paths concurrently with a semaphore.
func (g *GitHubSearcher) fetchAllDocs(ctx context.Context, paths []string) []indexedDoc {
	type result struct {
		index int
		doc   indexedDoc
		err   error
	}

	results := make(chan result, len(paths))
	sem := make(chan struct{}, maxConcurrentFetch)
	var wg sync.WaitGroup

	for i, p := range paths {
		wg.Add(1)
		go func(idx int, docPath string) {
			defer wg.Done()
			sem <- struct{}{}        // acquire
			defer func() { <-sem }() // release

			doc, err := g.fetchRawDoc(ctx, docPath)
			results <- result{index: idx, doc: doc, err: err}
		}(i, p)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	docs := make([]indexedDoc, 0, len(paths))
	for r := range results {
		if r.err != nil {
			g.logger.Warn("Failed to fetch doc", zap.Error(r.err))
			continue
		}
		docs = append(docs, r.doc)
	}
	return docs
}

// fetchRawDoc fetches and processes a single document from raw.githubusercontent.com.
func (g *GitHubSearcher) fetchRawDoc(ctx context.Context, docPath string) (indexedDoc, error) {
	rawURL := fmt.Sprintf("%s/%s/main/%s", g.rawBaseURL, docsRepo, docPath)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return indexedDoc{}, err
	}

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

// reASCIIWord matches sequences of ASCII alphanumeric characters (including hyphens).
var reASCIIWord = regexp.MustCompile(`[a-zA-Z0-9][\w-]*`)

// tokenizeQuery splits a query into lowercase tokens, removing stop words.
// It handles both space-delimited languages (English) and non-space-delimited
// languages (Japanese, Chinese) by extracting ASCII words from mixed text
// and splitting CJK characters into individual tokens.
func tokenizeQuery(query string) []string {
	lower := strings.ToLower(query)
	seen := make(map[string]bool)
	tokens := make([]string, 0)

	addToken := func(t string) {
		if t != "" && !stopWords[t] && !seen[t] {
			seen[t] = true
			tokens = append(tokens, t)
		}
	}

	// First pass: extract space-delimited words (handles English and mixed text)
	for _, w := range strings.Fields(lower) {
		// Extract ASCII words from mixed tokens like "sdkについて"
		asciiWords := reASCIIWord.FindAllString(w, -1)
		if len(asciiWords) > 0 {
			for _, aw := range asciiWords {
				addToken(aw)
			}
		}
		// Extract CJK character runs as individual tokens
		cjkRun := extractCJKRuns(w)
		for _, run := range cjkRun {
			addToken(run)
			// Translate known Katakana loanwords to English equivalents
			if eng, ok := katakanaToEnglish[run]; ok {
				addToken(eng)
			}
		}
	}

	if len(tokens) == 0 {
		return nil
	}
	return tokens
}

// extractCJKRuns extracts CJK tokens from a string by isolating Katakana runs
// (which typically represent loanwords like タグ, フラグ, セグメント) as separate
// tokens. Non-Katakana CJK sequences (Hiragana + Kanji) are kept together.
// This yields: "タグでフラグを整理する" → ["タグ", "フラグ", "整理する"].
func extractCJKRuns(s string) []string {
	var runs []string
	var current []rune
	var isKatakana bool // true if current run is katakana

	flush := func() {
		if len(current) > 0 {
			runs = append(runs, string(current))
			current = current[:0]
		}
	}

	for _, r := range s {
		if !isCJK(r) {
			flush()
			continue
		}
		kata := unicode.Is(unicode.Katakana, r)
		if len(current) > 0 && kata != isKatakana {
			flush()
		}
		isKatakana = kata
		current = append(current, r)
	}
	flush()

	// Filter out short hiragana-only tokens (likely particles like で, を, の)
	filtered := runs[:0]
	for _, run := range runs {
		if isHiraganaOnly(run) && len([]rune(run)) <= 2 {
			continue
		}
		filtered = append(filtered, run)
	}
	return filtered
}

// isHiraganaOnly returns true if every rune in the string is Hiragana.
func isHiraganaOnly(s string) bool {
	for _, r := range s {
		if !unicode.Is(unicode.Hiragana, r) {
			return false
		}
	}
	return true
}

// isCJK returns true if the rune is a CJK character (Chinese, Japanese, Korean)
// including Hiragana, Katakana, CJK Unified Ideographs, and Hangul.
func isCJK(r rune) bool {
	return unicode.Is(unicode.Han, r) ||
		unicode.Is(unicode.Hiragana, r) ||
		unicode.Is(unicode.Katakana, r) ||
		unicode.Is(unicode.Hangul, r)
}

// scoreDoc computes a relevance score for a document against query tokens.
// Uses pre-computed lowercase fields and path segments from indexedDoc.
func scoreDoc(doc indexedDoc, queryTokens []string) float64 {
	var score float64

	for _, token := range queryTokens {
		// Path segment match (highest weight)
		for _, seg := range doc.pathSegments {
			if seg == token {
				score += 3.0
			} else if strings.Contains(seg, token) {
				score += 2.0
			}
		}
		// Title match
		if strings.Contains(doc.lowerTitle, token) {
			score += 2.5
		}
		// Content match (normalized, capped)
		count := strings.Count(doc.lowerContent, token)
		if count > 0 {
			score += math.Min(float64(count), 5.0) / 5.0
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
