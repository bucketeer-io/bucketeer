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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/llm"
	llmmock "github.com/bucketeer-io/bucketeer/v2/pkg/aichat/llm/mock"
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

		cfg := ChatConfig{Model: "test-model", MaxTokens: 100, Temperature: 0.5}

		req := &aichatproto.ChatRequest{
			Messages: []*aichatproto.ChatMessage{
				{Role: aichatproto.ChatMessage_ROLE_USER, Content: "Hello"},
			},
			EnvironmentId: "env-1",
		}

		respChan, chatErrChan := streamChat(t.Context(), mockClient, nil, nil, cfg, req, zap.NewNop())

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

		respChan, chatErrChan := streamChat(t.Context(), mockClient, nil, nil, cfg, req, zap.NewNop())

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

		respChan, chatErrChan := streamChat(t.Context(), mockClient, nil, mockFeatureClient, cfg, req, zap.NewNop())
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

		respChan, chatErrChan := streamChat(t.Context(), mockClient, nil, mockFeatureClient, cfg, req, zap.NewNop())

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

		respChan, chatErrChan := streamChat(t.Context(), mockClient, nil, mockFeatureClient, cfg, req, zap.NewNop())
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

		respChan, chatErrChan := streamChat(t.Context(), mockClient, nil, nil, cfg, req, zap.NewNop())
		for range respChan {
		}
		<-chatErrChan
	})
}

func TestSanitizeUserInput(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "normal input",
			input:    "How do I use feature flags?",
			expected: "How do I use feature flags?",
		},
		{
			desc:     "removes newlines",
			input:    "line1\nline2\rline3",
			expected: "line1 line2 line3",
		},
		{
			desc:     "escapes HTML",
			input:    "<script>alert('xss')</script>",
			expected: "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
		{
			desc:     "trims whitespace",
			input:    "  hello  ",
			expected: "hello",
		},
		{
			desc:     "truncates long input",
			input:    strings.Repeat("あ", maxInputLength+100),
			expected: html.EscapeString(strings.Repeat("あ", maxInputLength)),
		},
		{
			desc:     "empty input",
			input:    "",
			expected: "",
		},
		{
			desc:     "whitespace only",
			input:    "   \n\r  ",
			expected: "",
		},
		{
			desc:     "japanese text",
			input:    "フラグの使い方を教えてください",
			expected: "フラグの使い方を教えてください",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			result := sanitizeUserInput(p.input)
			assert.Equal(t, p.expected, result)
		})
	}
}

func TestNormalizeInput(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "short input unchanged",
			input:    "Hello world",
			expected: "Hello world",
		},
		{
			desc:     "replaces newlines",
			input:    "line1\nline2\rline3",
			expected: "line1 line2 line3",
		},
		{
			desc:     "truncates long input",
			input:    strings.Repeat("a", maxInputLength+50),
			expected: strings.Repeat("a", maxInputLength),
		},
		{
			desc:     "does not HTML escape",
			input:    "<b>bold</b>",
			expected: "<b>bold</b>",
		},
		{
			desc:     "trims whitespace",
			input:    "  hello  ",
			expected: "hello",
		},
		{
			desc:     "handles multibyte truncation",
			input:    strings.Repeat("あ", maxInputLength+10),
			expected: strings.Repeat("あ", maxInputLength),
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			result := normalizeInput(p.input)
			assert.Equal(t, p.expected, result)
		})
	}
}

func TestExtractSearchQuery(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc        string
		setupClient func(ctrl *gomock.Controller) *llmmock.MockClient
		query       string
		expected    string
		expectErr   bool
	}{
		{
			desc:        "nil client returns original query",
			setupClient: nil,
			query:       "some query",
			expected:    "some query",
		},
		{
			desc: "empty query returns empty query",
			setupClient: func(ctrl *gomock.Controller) *llmmock.MockClient {
				return llmmock.NewMockClient(ctrl)
			},
			query:    "",
			expected: "",
		},
		{
			desc: "successful extraction",
			setupClient: func(ctrl *gomock.Controller) *llmmock.MockClient {
				mc := llmmock.NewMockClient(ctrl)
				mc.EXPECT().
					Chat(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("feature flags sdk", nil)
				return mc
			},
			query:    "フィーチャーフラグのSDKについて教えて",
			expected: "feature flags sdk",
		},
		{
			desc: "LLM returns empty string falls back to original query",
			setupClient: func(ctrl *gomock.Controller) *llmmock.MockClient {
				mc := llmmock.NewMockClient(ctrl)
				mc.EXPECT().
					Chat(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", nil)
				return mc
			},
			query:    "original query",
			expected: "original query",
		},
		{
			desc: "LLM returns error",
			setupClient: func(ctrl *gomock.Controller) *llmmock.MockClient {
				mc := llmmock.NewMockClient(ctrl)
				mc.EXPECT().
					Chat(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", assert.AnError)
				return mc
			},
			query:     "some query",
			expected:  "",
			expectErr: true,
		},
		{
			desc: "context cancellation",
			setupClient: func(ctrl *gomock.Controller) *llmmock.MockClient {
				mc := llmmock.NewMockClient(ctrl)
				mc.EXPECT().
					Chat(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", context.Canceled)
				return mc
			},
			query:     "some query",
			expected:  "",
			expectErr: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			var client llm.Client
			if p.setupClient != nil {
				client = p.setupClient(ctrl)
			}
			result, err := extractSearchQuery(t.Context(), client, p.query, "test-model")
			if p.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, p.expected, result)
		})
	}
}

func TestDefaultChatConfig(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc        string
		input       ChatConfig
		checkModel  string
		checkTokens int
		checkTemp   *float64
	}{
		{
			desc:        "fills empty fields with defaults",
			input:       ChatConfig{},
			checkModel:  "",
			checkTokens: 1000,
		},
		{
			desc: "preserves explicit values",
			input: ChatConfig{
				Model:       "gpt-4",
				MaxTokens:   500,
				Temperature: 0.3,
			},
			checkModel:  "gpt-4",
			checkTokens: 500,
			checkTemp:   float64Ptr(0.3),
		},
		{
			desc: "preserves zero temperature",
			input: ChatConfig{
				Model:       "gpt-4",
				MaxTokens:   500,
				Temperature: 0.0,
			},
			checkTemp: float64Ptr(0.0),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			cfg := defaultChatConfig(p.input)
			if p.checkModel != "" {
				assert.Equal(t, p.checkModel, cfg.Model)
			}
			if p.checkTokens != 0 {
				assert.Equal(t, p.checkTokens, cfg.MaxTokens)
			}
			if p.checkTemp != nil {
				assert.InDelta(t, *p.checkTemp, cfg.Temperature, 0.001)
			}
		})
	}
}

func float64Ptr(v float64) *float64 {
	return &v
}
