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
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/llm"
	llmmock "github.com/bucketeer-io/bucketeer/v2/pkg/aichat/llm/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/rag"
	featureclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	aichatproto "github.com/bucketeer-io/bucketeer/v2/proto/aichat"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestStreamChat(t *testing.T) {
	t.Parallel()

	t.Run("streams response successfully", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockClient := llmmock.NewMockClient(ctrl)

		// Mock embedding for RAG
		mockClient.EXPECT().
			CreateEmbeddings(gomock.Any(), gomock.Any(), gomock.Any()).
			Return([][]float32{{0.1, 0.2}}, nil).
			AnyTimes()

		// Mock streaming chat
		mockClient.EXPECT().
			StreamChat(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, _ []llm.Message, _ llm.StreamOptions) (<-chan llm.Chunk, <-chan error) {
				ch := make(chan llm.Chunk, 3)
				ech := make(chan error, 1)
				go func() {
					ch <- llm.Chunk{Content: "Hello", Done: false}
					ch <- llm.Chunk{Content: " World", Done: false}
					ch <- llm.Chunk{Content: "", Done: true, FinishReason: "stop"}
					close(ch)
					close(ech)
				}()
				return ch, ech
			})

		// RAG service with minimal docs
		ragSvc := createTestRAGService(t, mockClient)

		cfg := ChatConfig{Model: "test-model", MaxTokens: 100, Temperature: 0.5}

		req := &aichatproto.ChatRequest{
			Messages: []*aichatproto.ChatMessage{
				{Role: aichatproto.ChatMessage_ROLE_USER, Content: "Hello"},
			},
			EnvironmentId: "env-1",
		}

		respChan, chatErrChan := streamChat(context.Background(), mockClient, ragSvc, nil, cfg, req, zap.NewNop())

		var chunks []*aichatproto.ChatStreamResponse
		for chunk := range respChan {
			chunks = append(chunks, chunk)
		}

		err := <-chatErrChan
		assert.NoError(t, err)
		require.Len(t, chunks, 3)
		assert.Equal(t, "Hello", chunks[0].Content)
		assert.Equal(t, " World", chunks[1].Content)
		assert.True(t, chunks[2].Done)
	})

	t.Run("handles LLM error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockClient := llmmock.NewMockClient(ctrl)

		mockClient.EXPECT().
			StreamChat(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, _ []llm.Message, _ llm.StreamOptions) (<-chan llm.Chunk, <-chan error) {
				ch := make(chan llm.Chunk)
				ech := make(chan error, 1)
				go func() {
					ech <- assert.AnError
					close(ch)
					close(ech)
				}()
				return ch, ech
			})

		cfg := ChatConfig{Model: "test-model"}

		req := &aichatproto.ChatRequest{
			Messages: []*aichatproto.ChatMessage{
				{Role: aichatproto.ChatMessage_ROLE_USER, Content: "test"},
			},
			EnvironmentId: "env-1",
		}

		respChan, chatErrChan := streamChat(context.Background(), mockClient, nil, nil, cfg, req, zap.NewNop())

		// Drain response channel
		for range respChan {
		}

		err := <-chatErrChan
		assert.Error(t, err)
	})

	t.Run("includes feature context in system prompt", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockClient := llmmock.NewMockClient(ctrl)

		// Capture system prompt to verify feature context is included
		var capturedMessages []llm.Message
		mockClient.EXPECT().
			StreamChat(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, msgs []llm.Message, _ llm.StreamOptions) (<-chan llm.Chunk, <-chan error) {
				capturedMessages = msgs
				ch := make(chan llm.Chunk, 1)
				ech := make(chan error, 1)
				go func() {
					ch <- llm.Chunk{Content: "ok", Done: true, FinishReason: "stop"}
					close(ch)
					close(ech)
				}()
				return ch, ech
			})

		mockFeatureClient := featureclientmock.NewMockClient(ctrl)
		mockFeatureClient.EXPECT().
			GetFeature(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&featureproto.GetFeatureResponse{
				Feature: &featureproto.Feature{
					Id:            "flag-1",
					Name:          "Dark Mode",
					Description:   "Toggle dark mode",
					Enabled:       true,
					VariationType: featureproto.Feature_BOOLEAN,
					Variations: []*featureproto.Variation{
						{Id: "v1", Value: "true", Name: "ON"},
						{Id: "v2", Value: "false", Name: "OFF"},
					},
					Tags: []string{"ui"},
				},
			}, nil)

		cfg := ChatConfig{Model: "test-model"}
		req := &aichatproto.ChatRequest{
			Messages: []*aichatproto.ChatMessage{
				{Role: aichatproto.ChatMessage_ROLE_USER, Content: "tell me about this flag"},
			},
			EnvironmentId: "env-1",
			PageContext: &aichatproto.PageContext{
				PageType:  aichatproto.PageContext_PAGE_TYPE_TARGETING,
				FeatureId: "flag-1",
			},
		}

		respChan, chatErrChan := streamChat(context.Background(), mockClient, nil, mockFeatureClient, cfg, req, zap.NewNop())
		for range respChan {
		}
		<-chatErrChan

		require.NotEmpty(t, capturedMessages)
		systemPrompt := capturedMessages[0].Content
		assert.Contains(t, systemPrompt, "Feature Flag Details")
		assert.Contains(t, systemPrompt, "Dark Mode")
		assert.Contains(t, systemPrompt, "Toggle dark mode")
		assert.Contains(t, systemPrompt, "Tags: ui")
		// Privacy: variation values must not leak
		assert.NotContains(t, systemPrompt, "\"true\"")
		assert.NotContains(t, systemPrompt, "\"false\"")
	})

	t.Run("gracefully degrades when feature fetch fails", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockClient := llmmock.NewMockClient(ctrl)

		mockClient.EXPECT().
			StreamChat(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, _ []llm.Message, _ llm.StreamOptions) (<-chan llm.Chunk, <-chan error) {
				ch := make(chan llm.Chunk, 1)
				ech := make(chan error, 1)
				go func() {
					ch <- llm.Chunk{Content: "ok", Done: true, FinishReason: "stop"}
					close(ch)
					close(ech)
				}()
				return ch, ech
			})

		mockFeatureClient := featureclientmock.NewMockClient(ctrl)
		mockFeatureClient.EXPECT().
			GetFeature(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, assert.AnError)

		cfg := ChatConfig{Model: "test-model"}
		req := &aichatproto.ChatRequest{
			Messages: []*aichatproto.ChatMessage{
				{Role: aichatproto.ChatMessage_ROLE_USER, Content: "test"},
			},
			EnvironmentId: "env-1",
			PageContext: &aichatproto.PageContext{
				PageType:  aichatproto.PageContext_PAGE_TYPE_TARGETING,
				FeatureId: "flag-1",
			},
		}

		respChan, chatErrChan := streamChat(context.Background(), mockClient, nil, mockFeatureClient, cfg, req, zap.NewNop())

		var chunks []*aichatproto.ChatStreamResponse
		for chunk := range respChan {
			chunks = append(chunks, chunk)
		}
		err := <-chatErrChan

		// Should succeed despite feature fetch failure
		assert.NoError(t, err)
		require.NotEmpty(t, chunks)
		assert.True(t, chunks[len(chunks)-1].Done)
	})

	t.Run("skips feature fetch when no featureId", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockClient := llmmock.NewMockClient(ctrl)

		mockClient.EXPECT().
			StreamChat(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, _ []llm.Message, _ llm.StreamOptions) (<-chan llm.Chunk, <-chan error) {
				ch := make(chan llm.Chunk, 1)
				ech := make(chan error, 1)
				go func() {
					ch <- llm.Chunk{Content: "ok", Done: true, FinishReason: "stop"}
					close(ch)
					close(ech)
				}()
				return ch, ech
			})

		// featureClient should NOT be called because FeatureId is empty
		mockFeatureClient := featureclientmock.NewMockClient(ctrl)
		// No EXPECT — any call to GetFeature will cause test failure

		cfg := ChatConfig{Model: "test-model"}
		req := &aichatproto.ChatRequest{
			Messages: []*aichatproto.ChatMessage{
				{Role: aichatproto.ChatMessage_ROLE_USER, Content: "test"},
			},
			EnvironmentId: "env-1",
			PageContext: &aichatproto.PageContext{
				PageType: aichatproto.PageContext_PAGE_TYPE_FEATURE_FLAGS,
				// FeatureId intentionally empty
			},
		}

		respChan, chatErrChan := streamChat(context.Background(), mockClient, nil, mockFeatureClient, cfg, req, zap.NewNop())
		for range respChan {
		}
		<-chatErrChan
	})

	t.Run("sanitizes user input", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockClient := llmmock.NewMockClient(ctrl)

		mockClient.EXPECT().
			StreamChat(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, msgs []llm.Message, _ llm.StreamOptions) (<-chan llm.Chunk, <-chan error) {
				// Verify user message was sanitized (HTML escaped)
				for _, m := range msgs {
					if m.Role == "user" {
						assert.NotContains(t, m.Content, "<script>")
						assert.Contains(t, m.Content, "&lt;script&gt;")
					}
				}
				ch := make(chan llm.Chunk, 1)
				ech := make(chan error, 1)
				go func() {
					ch <- llm.Chunk{Content: "ok", Done: true, FinishReason: "stop"}
					close(ch)
					close(ech)
				}()
				return ch, ech
			})

		cfg := ChatConfig{Model: "test-model"}

		req := &aichatproto.ChatRequest{
			Messages: []*aichatproto.ChatMessage{
				{Role: aichatproto.ChatMessage_ROLE_USER, Content: "<script>alert('xss')</script>"},
			},
			EnvironmentId: "env-1",
		}

		respChan, chatErrChan := streamChat(context.Background(), mockClient, nil, nil, cfg, req, zap.NewNop())
		for range respChan {
		}
		<-chatErrChan
	})
}

func createTestRAGService(t *testing.T, client llm.Client) *rag.Service {
	t.Helper()
	ragSvc, err := rag.NewService(client, "test-model", zap.NewNop())
	require.NoError(t, err)
	return ragSvc
}

func TestSanitizeUserInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal input",
			input:    "How do I use feature flags?",
			expected: "How do I use feature flags?",
		},
		{
			name:     "removes newlines",
			input:    "line1\nline2\rline3",
			expected: "line1 line2 line3",
		},
		{
			name:     "escapes HTML",
			input:    "<script>alert('xss')</script>",
			expected: "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
		{
			name:     "trims whitespace",
			input:    "  hello  ",
			expected: "hello",
		},
		{
			name:     "truncates long input",
			input:    strings.Repeat("あ", maxInputLength+100),
			expected: html.EscapeString(strings.Repeat("あ", maxInputLength)),
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace only",
			input:    "   \n\r  ",
			expected: "",
		},
		{
			name:     "japanese text",
			input:    "フラグの使い方を教えてください",
			expected: "フラグの使い方を教えてください",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := sanitizeUserInput(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeUserInput_MaxLength(t *testing.T) {
	t.Parallel()

	// Create input longer than maxInputLength
	longInput := strings.Repeat("a", maxInputLength+100)
	result := sanitizeUserInput(longInput)
	assert.Equal(t, maxInputLength, utf8.RuneCountInString(result))
}

func TestLimitInputLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "short input unchanged",
			input:    "Hello world",
			expected: "Hello world",
		},
		{
			name:     "replaces newlines",
			input:    "line1\nline2\rline3",
			expected: "line1 line2 line3",
		},
		{
			name:     "truncates long input",
			input:    strings.Repeat("a", maxInputLength+50),
			expected: strings.Repeat("a", maxInputLength),
		},
		{
			name:     "does not HTML escape",
			input:    "<b>bold</b>",
			expected: "<b>bold</b>",
		},
		{
			name:     "trims whitespace",
			input:    "  hello  ",
			expected: "hello",
		},
		{
			name:     "handles multibyte truncation",
			input:    strings.Repeat("あ", maxInputLength+10),
			expected: strings.Repeat("あ", maxInputLength),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := limitInputLength(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultChatConfig(t *testing.T) {
	t.Parallel()

	t.Run("fills empty fields with defaults", func(t *testing.T) {
		t.Parallel()
		cfg := defaultChatConfig(ChatConfig{})
		assert.Equal(t, "gpt-4o-mini", cfg.Model)
		assert.Equal(t, 1000, cfg.MaxTokens)
	})

	t.Run("preserves explicit values", func(t *testing.T) {
		t.Parallel()
		cfg := defaultChatConfig(ChatConfig{
			Model:       "gpt-4",
			MaxTokens:   500,
			Temperature: 0.3,
		})
		assert.Equal(t, "gpt-4", cfg.Model)
		assert.Equal(t, 500, cfg.MaxTokens)
		assert.InDelta(t, 0.3, cfg.Temperature, 0.001)
	})

	t.Run("preserves zero temperature", func(t *testing.T) {
		t.Parallel()
		cfg := defaultChatConfig(ChatConfig{
			Model:       "gpt-4",
			MaxTokens:   500,
			Temperature: 0.0,
		})
		assert.InDelta(t, 0.0, cfg.Temperature, 0.001)
	})
}
