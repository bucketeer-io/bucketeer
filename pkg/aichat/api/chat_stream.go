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

package api

import (
	"context"
	"html"
	"strings"
	"time"
	"unicode/utf8"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/llm"
	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/rag"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	aichatproto "github.com/bucketeer-io/bucketeer/v2/proto/aichat"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const maxInputLength = 2000
const maxMessages = 50

var newlineReplacer = strings.NewReplacer("\n", " ", "\r", " ")

// ChatConfig holds configuration for chat streaming.
// It is an alias for llm.StreamOptions to avoid redundant type definitions.
type ChatConfig = llm.StreamOptions

func defaultChatConfig(cfg ChatConfig) ChatConfig {
	if cfg.Model == "" {
		cfg.Model = "gpt-4o-mini"
	}
	if cfg.MaxTokens == 0 {
		cfg.MaxTokens = 1000
	}
	// Temperature is intentionally not defaulted here because 0.0 is a valid
	// value (deterministic output). Callers must set it explicitly.
	return cfg
}

// streamChat generates a streaming chat response using LLM and optional RAG.
func streamChat(
	ctx context.Context,
	llmClient llm.Client,
	ragSearcher rag.Searcher,
	featureClient featureclient.Client,
	cfg ChatConfig,
	req *aichatproto.ChatRequest,
	logger *zap.Logger,
) (<-chan *aichatproto.ChatStreamResponse, <-chan error) {
	responseChan := make(chan *aichatproto.ChatStreamResponse, 100)
	errChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errChan)

		// Convert to LLM messages and find last user message for RAG
		var lastUserMessage string
		llmMessages := make([]llm.Message, 0, len(req.Messages)+1)
		for _, m := range req.Messages {
			role := llm.RoleUser
			content := m.Content
			if m.Role == aichatproto.ChatMessage_ROLE_ASSISTANT {
				role = llm.RoleAssistant
				content = limitInputLength(content)
			} else {
				lastUserMessage = limitInputLength(content) // clean text for RAG (no HTML escape)
				content = sanitizeUserInput(content)        // sanitized for LLM
			}
			llmMessages = append(llmMessages, llm.Message{
				Role:    role,
				Content: content,
			})
		}

		// Fetch RAG docs and feature details concurrently (both are independent network calls)
		var relevantDocs []rag.DocChunk
		var featureContext string
		g, gCtx := errgroup.WithContext(ctx)
		if ragSearcher != nil && lastUserMessage != "" {
			g.Go(func() error {
				// Extract English search keywords via LLM, then search
				searchQuery := lastUserMessage
				if extracted, err := extractSearchQuery(
					gCtx, llmClient, lastUserMessage, cfg.Model,
				); err != nil {
					logger.Warn("Keyword extraction failed, using raw query",
						zap.Error(err))
				} else {
					searchQuery = extracted
				}
				docs, err := ragSearcher.Search(gCtx, searchQuery, 3)
				if err != nil {
					logger.Warn("RAG search failed, continuing without context",
						zap.Error(err))
				} else {
					relevantDocs = docs
				}
				return nil // graceful degradation
			})
		}
		if featureClient != nil && req.PageContext != nil && req.PageContext.FeatureId != "" {
			g.Go(func() error {
				resp, err := featureClient.GetFeature(gCtx, &featureproto.GetFeatureRequest{
					EnvironmentId: req.EnvironmentId,
					Id:            req.PageContext.FeatureId,
				})
				if err != nil {
					logger.Warn("Failed to get feature for chat context, continuing without",
						zap.Error(err),
						zap.String("featureId", req.PageContext.FeatureId),
					)
				} else if resp.Feature != nil {
					featureContext = buildFeatureContext(resp.Feature)
				}
				return nil // graceful degradation
			})
		}
		_ = g.Wait()

		// Build system prompt with context, RAG, and feature details
		systemPrompt := buildSystemPrompt(req.PageContext, relevantDocs, featureContext)
		messages := append([]llm.Message{{Role: llm.RoleSystem, Content: systemPrompt}}, llmMessages...)

		// Stream chat from LLM
		chunkChan, llmErrChan := llmClient.StreamChat(ctx, messages, cfg)

		for chunkChan != nil || llmErrChan != nil {
			select {
			case <-ctx.Done():
				return
			case err, ok := <-llmErrChan:
				if !ok {
					llmErrChan = nil
					continue
				}
				if err != nil {
					logger.Error("LLM stream error", zap.Error(err))
					errChan <- err
					return
				}
			case chunk, ok := <-chunkChan:
				if !ok {
					chunkChan = nil
					continue
				}
				select {
				case responseChan <- &aichatproto.ChatStreamResponse{
					Content:      chunk.Content,
					Done:         chunk.Done,
					FinishReason: chunk.FinishReason,
				}:
				case <-ctx.Done():
					return
				}
				if chunk.Done {
					return
				}
			}
		}
	}()

	return responseChan, errChan
}

// normalizeInput replaces newlines with spaces and truncates to maxInputLength runes.
func normalizeInput(input string) string {
	input = newlineReplacer.Replace(input)
	if utf8.RuneCountInString(input) > maxInputLength {
		runes := []rune(input)
		input = string(runes[:maxInputLength])
	}
	return strings.TrimSpace(input)
}

// sanitizeUserInput cleans user input for safety (normalize + HTML escape).
func sanitizeUserInput(input string) string {
	return html.EscapeString(normalizeInput(input))
}

// limitInputLength applies normalization without HTML escaping.
// Used for assistant messages and RAG queries.
func limitInputLength(input string) string {
	return normalizeInput(input)
}

const keywordExtractionPrompt = `Extract 3-5 English search keywords from the user query.
Return only lowercase keywords separated by spaces. No explanation.`

// extractSearchQuery uses the LLM to convert a user query (which may be in
// any language) into English search keywords for RAG document scoring.
func extractSearchQuery(
	ctx context.Context,
	client llm.Client,
	query string,
	model string,
) (string, error) {
	if client == nil || query == "" {
		return query, nil
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	messages := []llm.Message{
		{Role: llm.RoleSystem, Content: keywordExtractionPrompt},
		{Role: llm.RoleUser, Content: query},
	}
	opts := llm.StreamOptions{
		Model:       model,
		MaxTokens:   50,
		Temperature: 0,
	}
	result, err := client.Chat(ctx, messages, opts)
	if err != nil {
		return "", err
	}
	result = strings.TrimSpace(result)
	if result == "" {
		return query, nil
	}
	return result, nil
}
