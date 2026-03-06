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
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/llm"
	llmmock "github.com/bucketeer-io/bucketeer/v2/pkg/aichat/llm/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/ratelimit"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	aichatproto "github.com/bucketeer-io/bucketeer/v2/proto/aichat"
)

// mockVerifier implements token.Verifier for testing.
type mockVerifier struct {
	token *token.AccessToken
	err   error
}

func (m *mockVerifier) VerifyAccessToken(string) (*token.AccessToken, error) {
	return m.token, m.err
}

func (m *mockVerifier) VerifyRefreshToken(string) (*token.RefreshToken, error) {
	return nil, nil
}

func (m *mockVerifier) VerifyDemoCreationToken(string) (*token.DemoCreationToken, error) {
	return nil, nil
}

func TestChatHTTPService_MethodNotAllowed(t *testing.T) {
	t.Parallel()
	handler := NewChatHTTPService(nil, nil, ChatConfig{}, nil, nil, nil, zap.NewNop())

	req := httptest.NewRequest(http.MethodGet, "/v1/aichat/chat", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestChatHTTPService_Unauthorized_NoHeader(t *testing.T) {
	t.Parallel()
	handler := NewChatHTTPService(nil, nil, ChatConfig{}, nil, nil, nil, zap.NewNop())

	req := httptest.NewRequest(http.MethodPost, "/v1/aichat/chat", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestChatHTTPService_Unauthorized_InvalidFormat(t *testing.T) {
	t.Parallel()
	handler := NewChatHTTPService(nil, nil, ChatConfig{}, nil, nil, nil, zap.NewNop())

	req := httptest.NewRequest(http.MethodPost, "/v1/aichat/chat", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestChatHTTPService_Unauthorized_InvalidToken(t *testing.T) {
	t.Parallel()
	v := &mockVerifier{
		err: fmt.Errorf("invalid token"),
	}
	handler := NewChatHTTPService(nil, nil, ChatConfig{}, v, nil, nil, zap.NewNop())

	req := httptest.NewRequest(http.MethodPost, "/v1/aichat/chat", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestChatHTTPService_BadRequest_InvalidJSON(t *testing.T) {
	t.Parallel()
	v := &mockVerifier{
		token: &token.AccessToken{Email: "test@example.com"},
	}
	handler := NewChatHTTPService(nil, nil, ChatConfig{}, v, nil, nil, zap.NewNop())

	req := httptest.NewRequest(http.MethodPost, "/v1/aichat/chat", strings.NewReader("{invalid"))
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestChatHTTPService_BadRequest_EmptyMessages(t *testing.T) {
	t.Parallel()
	v := &mockVerifier{
		token: &token.AccessToken{Email: "test@example.com"},
	}
	handler := NewChatHTTPService(nil, nil, ChatConfig{}, v, nil, nil, zap.NewNop())

	body := `{"messages":[],"environmentId":"env-1"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/aichat/chat", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestChatHTTPService_BadRequest_EmptyEnvironmentID(t *testing.T) {
	t.Parallel()
	v := &mockVerifier{
		token: &token.AccessToken{Email: "test@example.com"},
	}
	handler := NewChatHTTPService(nil, nil, ChatConfig{}, v, nil, nil, zap.NewNop())

	body := `{"messages":[{"role":"user","content":"hello"}],"environmentId":""}`
	req := httptest.NewRequest(http.MethodPost, "/v1/aichat/chat", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestChatHTTPService_StreamSuccess(t *testing.T) {
	t.Parallel()
	v := &mockVerifier{
		token: &token.AccessToken{Email: "test@example.com"},
	}

	// Create a mock LLM client that returns streamed responses
	mockLLM := createMockLLMClient(t, []*aichatproto.ChatStreamResponse{
		{Content: "Hello", Done: false},
		{Content: " world", Done: false},
		{Content: "", Done: true, FinishReason: "stop"},
	}, nil)

	handler := NewChatHTTPService(mockLLM, nil, ChatConfig{Model: "test-model", MaxTokens: 100, Temperature: 0.5}, v, newMockAccountClient(t), nil, zap.NewNop())

	body := `{"messages":[{"role":"user","content":"hi"}],"environmentId":"env-1","pageContext":{"pageType":"feature_flags"}}`
	req := httptest.NewRequest(http.MethodPost, "/v1/aichat/chat", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache", rec.Header().Get("Cache-Control"))

	// Parse SSE events
	events := parseSSEEvents(t, rec.Body.String())
	require.GreaterOrEqual(t, len(events), 2) // at least one data event + [DONE]

	// Verify last event is [DONE]
	assert.Equal(t, "[DONE]", events[len(events)-1])

	// Verify first data event has content
	var firstEvent map[string]interface{}
	err := json.Unmarshal([]byte(events[0]), &firstEvent)
	require.NoError(t, err)
	assert.NotEmpty(t, firstEvent["content"])
}

func TestChatHTTPService_StreamError(t *testing.T) {
	t.Parallel()
	v := &mockVerifier{
		token: &token.AccessToken{Email: "test@example.com"},
	}

	mockLLM := createMockLLMClient(t, nil, fmt.Errorf("LLM error"))

	handler := NewChatHTTPService(mockLLM, nil, ChatConfig{Model: "test-model"}, v, newMockAccountClient(t), nil, zap.NewNop())

	body := `{"messages":[{"role":"user","content":"hi"}],"environmentId":"env-1"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/aichat/chat", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code) // SSE always returns 200, errors are sent as events
	assert.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))

	events := parseSSEEvents(t, rec.Body.String())
	require.GreaterOrEqual(t, len(events), 1)

	// Should contain error event
	found := false
	for _, e := range events {
		if e == "[DONE]" {
			continue
		}
		var ev map[string]interface{}
		if err := json.Unmarshal([]byte(e), &ev); err == nil {
			if _, ok := ev["error"]; ok {
				found = true
			}
		}
	}
	assert.True(t, found, "expected error event in SSE stream")
}

func TestChatHTTPService_ToProtoRequest(t *testing.T) {
	t.Parallel()
	handler := NewChatHTTPService(nil, nil, ChatConfig{}, nil, nil, nil, zap.NewNop())

	req := &ChatHTTPRequest{
		Messages: []ChatMessageHTTP{
			{Role: "user", Content: "hello"},
			{Role: "assistant", Content: "hi there"},
			{Role: "user", Content: "how to use feature flags?"},
		},
		PageContext: &PageContextHTTP{
			PageType:  "feature_flags",
			FeatureID: "flag-1",
			Metadata:  map[string]string{"key": "value"},
		},
		EnvironmentID: "env-1",
	}

	proto := handler.toProtoRequest(req)

	assert.Len(t, proto.Messages, 3)
	assert.Equal(t, aichatproto.ChatMessage_ROLE_USER, proto.Messages[0].Role)
	assert.Equal(t, "hello", proto.Messages[0].Content)
	assert.Equal(t, aichatproto.ChatMessage_ROLE_ASSISTANT, proto.Messages[1].Role)
	assert.Equal(t, aichatproto.ChatMessage_ROLE_USER, proto.Messages[2].Role)

	assert.Equal(t, aichatproto.PageContext_PAGE_TYPE_FEATURE_FLAGS, proto.PageContext.PageType)
	assert.Equal(t, "flag-1", proto.PageContext.FeatureId)
	assert.Equal(t, "value", proto.PageContext.Metadata["key"])
	assert.Equal(t, "env-1", proto.EnvironmentId)
}

func TestChatHTTPService_ToProtoRequest_AllPageTypes(t *testing.T) {
	t.Parallel()
	handler := NewChatHTTPService(nil, nil, ChatConfig{}, nil, nil, nil, zap.NewNop())

	tests := []struct {
		pageType string
		expected aichatproto.PageContext_PageType
	}{
		{"feature_flags", aichatproto.PageContext_PAGE_TYPE_FEATURE_FLAGS},
		{"targeting", aichatproto.PageContext_PAGE_TYPE_TARGETING},
		{"experiments", aichatproto.PageContext_PAGE_TYPE_EXPERIMENTS},
		{"segments", aichatproto.PageContext_PAGE_TYPE_SEGMENTS},
		{"autoops", aichatproto.PageContext_PAGE_TYPE_AUTOOPS},
		{"unknown", aichatproto.PageContext_PAGE_TYPE_UNSPECIFIED},
		{"", aichatproto.PageContext_PAGE_TYPE_UNSPECIFIED},
	}

	for _, tt := range tests {
		t.Run(tt.pageType, func(t *testing.T) {
			req := &ChatHTTPRequest{
				Messages:      []ChatMessageHTTP{{Role: "user", Content: "test"}},
				PageContext:   &PageContextHTTP{PageType: tt.pageType},
				EnvironmentID: "env-1",
			}
			proto := handler.toProtoRequest(req)
			assert.Equal(t, tt.expected, proto.PageContext.PageType)
		})
	}
}

func TestChatHTTPService_Register(t *testing.T) {
	t.Parallel()
	handler := NewChatHTTPService(nil, nil, ChatConfig{}, nil, nil, nil, zap.NewNop())

	mux := http.NewServeMux()
	handler.Register(mux)

	// Verify the handler was registered by making a request
	req := httptest.NewRequest(http.MethodGet, "/v1/aichat/chat", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// Should get MethodNotAllowed (405), not NotFound (404)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestChatHTTPService_RateLimited(t *testing.T) {
	t.Parallel()
	v := &mockVerifier{
		token: &token.AccessToken{Email: "test@example.com"},
	}

	limiter := ratelimit.NewLimiter(ratelimit.Config{
		RequestsPerMinute: 60,
		BurstSize:         1, // Only allow 1 request
	})

	mockLLM := createMockLLMClient(t, []*aichatproto.ChatStreamResponse{
		{Content: "test", Done: true},
	}, nil)

	handler := NewChatHTTPService(mockLLM, nil, ChatConfig{Model: "test-model"}, v, newMockAccountClient(t), nil, zap.NewNop(), WithRateLimiter(limiter))

	body := `{"messages":[{"role":"user","content":"hi"}],"environmentId":"env-1"}`

	// First request should succeed
	req1 := httptest.NewRequest(http.MethodPost, "/v1/aichat/chat", strings.NewReader(body))
	req1.Header.Set("Authorization", "Bearer valid-token")
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)
	assert.Equal(t, http.StatusOK, rec1.Code)

	// Second request should be rate limited
	req2 := httptest.NewRequest(http.MethodPost, "/v1/aichat/chat", strings.NewReader(body))
	req2.Header.Set("Authorization", "Bearer valid-token")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	assert.Equal(t, http.StatusTooManyRequests, rec2.Code)
}

func TestGetEnvironmentRole(t *testing.T) {
	t.Parallel()

	roles := []*accountproto.AccountV2_EnvironmentRole{
		{EnvironmentId: "env-1", Role: accountproto.AccountV2_Role_Environment_VIEWER},
		{EnvironmentId: "env-2", Role: accountproto.AccountV2_Role_Environment_EDITOR},
	}

	t.Run("returns matching role", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, accountproto.AccountV2_Role_Environment_VIEWER, getEnvironmentRole(roles, "env-1"))
		assert.Equal(t, accountproto.AccountV2_Role_Environment_EDITOR, getEnvironmentRole(roles, "env-2"))
	})

	t.Run("returns UNASSIGNED for unknown env", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, accountproto.AccountV2_Role_Environment_UNASSIGNED, getEnvironmentRole(roles, "env-unknown"))
	})

	t.Run("returns UNASSIGNED for nil roles", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, accountproto.AccountV2_Role_Environment_UNASSIGNED, getEnvironmentRole(nil, "env-1"))
	})
}

func TestChatHTTPService_BadRequest_TooManyMessages(t *testing.T) {
	t.Parallel()
	v := &mockVerifier{
		token: &token.AccessToken{Email: "test@example.com"},
	}
	handler := NewChatHTTPService(nil, nil, ChatConfig{}, v, nil, nil, zap.NewNop())

	// Build JSON array with maxMessages+1 messages
	var msgs strings.Builder
	msgs.WriteString("[")
	for i := 0; i < maxMessages+1; i++ {
		if i > 0 {
			msgs.WriteString(",")
		}
		msgs.WriteString(`{"role":"user","content":"hi"}`)
	}
	msgs.WriteString("]")
	body := fmt.Sprintf(`{"messages":%s,"environmentId":"env-1"}`, msgs.String())
	req := httptest.NewRequest(http.MethodPost, "/v1/aichat/chat", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// createMockLLMClient creates a gomock LLM client that returns the given
// responses or error via StreamChat.
func createMockLLMClient(
	t *testing.T,
	responses []*aichatproto.ChatStreamResponse,
	streamErr error,
) *llmmock.MockClient {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockLLM := llmmock.NewMockClient(ctrl)
	mockLLM.EXPECT().
		StreamChat(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(
			_ context.Context,
			_ []llm.Message,
			_ llm.StreamOptions,
		) (<-chan llm.Chunk, <-chan error) {
			chunkChan := make(chan llm.Chunk, len(responses)+1)
			errChan := make(chan error, 1)
			go func() {
				defer close(chunkChan)
				defer close(errChan)
				if streamErr != nil {
					errChan <- streamErr
					return
				}
				for _, r := range responses {
					chunkChan <- llm.Chunk{
						Content:      r.Content,
						Done:         r.Done,
						FinishReason: r.FinishReason,
					}
					if r.Done {
						return
					}
				}
			}()
			return chunkChan, errChan
		}).
		AnyTimes()
	return mockLLM
}

// parseSSEEvents parses SSE data from the response body.
func parseSSEEvents(t *testing.T, body string) []string {
	t.Helper()
	var events []string
	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			events = append(events, strings.TrimPrefix(line, "data: "))
		}
	}
	return events
}

// newMockAccountClient returns a gomock account client that allows
// GetAccountV2ByEnvironmentID and returns an active account with VIEWER role.
func newMockAccountClient(t *testing.T) *accountclientmock.MockClient {
	t.Helper()
	ctrl := gomock.NewController(t)
	mc := accountclientmock.NewMockClient(ctrl)
	mc.EXPECT().
		GetAccountV2ByEnvironmentID(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(&accountproto.GetAccountV2ByEnvironmentIDResponse{
			Account: &accountproto.AccountV2{
				Email:    "test@example.com",
				Disabled: false,
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{
						EnvironmentId: "env-1",
						Role:          accountproto.AccountV2_Role_Environment_VIEWER,
					},
				},
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			},
		}, nil).
		AnyTimes()
	return mc
}
