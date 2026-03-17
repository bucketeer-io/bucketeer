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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/llm"
	llmmock "github.com/bucketeer-io/bucketeer/v2/pkg/aichat/llm/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	aichatproto "github.com/bucketeer-io/bucketeer/v2/proto/aichat"
)

func createContextWithToken(t *testing.T) context.Context {
	t.Helper()
	tk := &token.AccessToken{
		Issuer:        "issuer",
		Audience:      "audience",
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         "test@example.com",
		IsSystemAdmin: true,
	}
	return context.WithValue(context.TODO(), rpc.AccessTokenKey, tk)
}

func createAIChatServiceForTest(
	t *testing.T,
	c *gomock.Controller,
	orgRole accountproto.AccountV2_Role_Organization,
	envRole accountproto.AccountV2_Role_Environment,
) *AIChatService {
	t.Helper()

	accountClientMock := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountV2ByEnvironmentIDResponse{
		Account: &accountproto.AccountV2{
			Email:            "test@example.com",
			OrganizationRole: orgRole,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "env-1",
					Role:          envRole,
				},
			},
		},
	}
	accountClientMock.EXPECT().
		GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).
		Return(ar, nil).
		AnyTimes()

	mockLLM := llmmock.NewMockClient(c)
	mockLLM.EXPECT().
		StreamChat(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(
			_ context.Context,
			_ []llm.Message,
			_ llm.StreamOptions,
		) (<-chan llm.Chunk, <-chan error) {
			chunkChan := make(chan llm.Chunk, 2)
			errChan := make(chan error, 1)
			go func() {
				defer close(chunkChan)
				defer close(errChan)
				chunkChan <- llm.Chunk{Content: "test response", Done: true, FinishReason: "stop"}
			}()
			return chunkChan, errChan
		}).
		AnyTimes()

	return NewAIChatService(
		mockLLM,
		nil,
		ChatConfig{Model: "test", MaxTokens: 100, Temperature: 0.5},
		accountClientMock,
		nil,
		zap.NewNop(),
	)
}

// mockChatStream implements aichatproto.AIChatService_ChatServer for testing.
type mockChatStream struct {
	ctx       context.Context
	responses []*aichatproto.ChatStreamResponse
	grpc.ServerStream
}

func (m *mockChatStream) Context() context.Context { return m.ctx }
func (m *mockChatStream) Send(resp *aichatproto.ChatStreamResponse) error {
	m.responses = append(m.responses, resp)
	return nil
}

func TestGetSuggestions(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		req         *aichatproto.GetSuggestionsRequest
		expectedErr codes.Code
	}{
		{
			desc: "error: unauthenticated",
			ctx:  context.Background(),
			req: &aichatproto.GetSuggestionsRequest{
				EnvironmentId: "env-1",
			},
			expectedErr: codes.Unauthenticated,
		},
		{
			desc: "error: empty environment id",
			ctx:  createContextWithToken(t),
			req: &aichatproto.GetSuggestionsRequest{
				EnvironmentId: "",
			},
			expectedErr: codes.InvalidArgument,
		},
		{
			desc: "success",
			ctx:  createContextWithToken(t),
			req: &aichatproto.GetSuggestionsRequest{
				EnvironmentId: "env-1",
				PageContext: &aichatproto.PageContext{
					PageType: aichatproto.PageContext_PAGE_TYPE_FEATURE_FLAGS,
				},
			},
			expectedErr: codes.OK,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			svc := createAIChatServiceForTest(
				t, mockController,
				accountproto.AccountV2_Role_Organization_ADMIN,
				accountproto.AccountV2_Role_Environment_VIEWER,
			)
			resp, err := svc.GetSuggestions(p.ctx, p.req)
			if p.expectedErr == codes.OK {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.Suggestions)
			} else {
				assert.Error(t, err)
				st, _ := gstatus.FromError(err)
				assert.Equal(t, p.expectedErr, st.Code())
			}
		})
	}
}

func TestChat(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	tooManyMessages := make([]*aichatproto.ChatMessage, maxMessages+1)
	for i := range tooManyMessages {
		tooManyMessages[i] = &aichatproto.ChatMessage{
			Role:    aichatproto.ChatMessage_ROLE_USER,
			Content: "test",
		}
	}

	patterns := []struct {
		desc        string
		ctx         context.Context
		req         *aichatproto.ChatRequest
		expectedErr codes.Code
		checkResp   bool
	}{
		{
			desc: "error: unauthenticated",
			ctx:  context.Background(),
			req: &aichatproto.ChatRequest{
				EnvironmentId: "env-1",
				Messages: []*aichatproto.ChatMessage{
					{Role: aichatproto.ChatMessage_ROLE_USER, Content: "hi"},
				},
			},
			expectedErr: codes.Unauthenticated,
		},
		{
			desc: "error: empty environment id",
			ctx:  createContextWithToken(t),
			req: &aichatproto.ChatRequest{
				EnvironmentId: "",
			},
			expectedErr: codes.InvalidArgument,
		},
		{
			desc: "error: empty messages",
			ctx:  createContextWithToken(t),
			req: &aichatproto.ChatRequest{
				EnvironmentId: "env-1",
				Messages:      []*aichatproto.ChatMessage{},
			},
			expectedErr: codes.InvalidArgument,
		},
		{
			desc: "error: too many messages",
			ctx:  createContextWithToken(t),
			req: &aichatproto.ChatRequest{
				EnvironmentId: "env-1",
				Messages:      tooManyMessages,
			},
			// ErrorTypeExceededMax maps to codes.Unknown via convertStatusCode
			expectedErr: codes.Unknown,
		},
		{
			desc: "success",
			ctx:  createContextWithToken(t),
			req: &aichatproto.ChatRequest{
				EnvironmentId: "env-1",
				Messages: []*aichatproto.ChatMessage{
					{Role: aichatproto.ChatMessage_ROLE_USER, Content: "hello"},
				},
			},
			expectedErr: codes.OK,
			checkResp:   true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			svc := createAIChatServiceForTest(
				t, mockController,
				accountproto.AccountV2_Role_Organization_ADMIN,
				accountproto.AccountV2_Role_Environment_VIEWER,
			)
			stream := &mockChatStream{ctx: p.ctx}
			err := svc.Chat(p.req, stream)
			if p.expectedErr == codes.OK {
				assert.NoError(t, err)
				if p.checkResp {
					require.NotEmpty(t, stream.responses)
					last := stream.responses[len(stream.responses)-1]
					assert.True(t, last.Done)
				}
			} else {
				assert.Error(t, err)
				st, _ := gstatus.FromError(err)
				assert.Equal(t, p.expectedErr, st.Code())
			}
		})
	}
}
