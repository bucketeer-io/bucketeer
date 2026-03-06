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

func TestGetSuggestions_Unauthenticated(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	svc := createAIChatServiceForTest(
		t, mockCtrl,
		accountproto.AccountV2_Role_Organization_ADMIN,
		accountproto.AccountV2_Role_Environment_VIEWER,
	)

	// No token in context
	ctx := context.Background()
	_, err := svc.GetSuggestions(ctx, &aichatproto.GetSuggestionsRequest{
		EnvironmentId: "env-1",
	})

	assert.Error(t, err)
	st, _ := gstatus.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestGetSuggestions_Success(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	svc := createAIChatServiceForTest(
		t, mockCtrl,
		accountproto.AccountV2_Role_Organization_ADMIN,
		accountproto.AccountV2_Role_Environment_VIEWER,
	)

	ctx := createContextWithToken(t)
	resp, err := svc.GetSuggestions(ctx, &aichatproto.GetSuggestionsRequest{
		EnvironmentId: "env-1",
		PageContext: &aichatproto.PageContext{
			PageType: aichatproto.PageContext_PAGE_TYPE_FEATURE_FLAGS,
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Suggestions)
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

func TestChat_EmptyEnvironmentID(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	svc := createAIChatServiceForTest(
		t, mockCtrl,
		accountproto.AccountV2_Role_Organization_ADMIN,
		accountproto.AccountV2_Role_Environment_VIEWER,
	)

	stream := &mockChatStream{ctx: createContextWithToken(t)}
	err := svc.Chat(&aichatproto.ChatRequest{EnvironmentId: ""}, stream)

	assert.Error(t, err)
	st, _ := gstatus.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestChat_EmptyMessages(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	svc := createAIChatServiceForTest(
		t, mockCtrl,
		accountproto.AccountV2_Role_Organization_ADMIN,
		accountproto.AccountV2_Role_Environment_VIEWER,
	)

	stream := &mockChatStream{ctx: createContextWithToken(t)}
	err := svc.Chat(&aichatproto.ChatRequest{
		EnvironmentId: "env-1",
		Messages:      []*aichatproto.ChatMessage{},
	}, stream)

	assert.Error(t, err)
	st, _ := gstatus.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestChat_TooManyMessages(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	svc := createAIChatServiceForTest(
		t, mockCtrl,
		accountproto.AccountV2_Role_Organization_ADMIN,
		accountproto.AccountV2_Role_Environment_VIEWER,
	)

	msgs := make([]*aichatproto.ChatMessage, maxMessages+1)
	for i := range msgs {
		msgs[i] = &aichatproto.ChatMessage{
			Role:    aichatproto.ChatMessage_ROLE_USER,
			Content: "test",
		}
	}

	stream := &mockChatStream{ctx: createContextWithToken(t)}
	err := svc.Chat(&aichatproto.ChatRequest{
		EnvironmentId: "env-1",
		Messages:      msgs,
	}, stream)

	assert.Error(t, err)
}

func TestChat_Unauthenticated(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	svc := createAIChatServiceForTest(
		t, mockCtrl,
		accountproto.AccountV2_Role_Organization_ADMIN,
		accountproto.AccountV2_Role_Environment_VIEWER,
	)

	stream := &mockChatStream{ctx: context.Background()}
	err := svc.Chat(&aichatproto.ChatRequest{
		EnvironmentId: "env-1",
		Messages:      []*aichatproto.ChatMessage{{Role: aichatproto.ChatMessage_ROLE_USER, Content: "hi"}},
	}, stream)

	assert.Error(t, err)
	st, _ := gstatus.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestChat_Success(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	svc := createAIChatServiceForTest(
		t, mockCtrl,
		accountproto.AccountV2_Role_Organization_ADMIN,
		accountproto.AccountV2_Role_Environment_VIEWER,
	)

	stream := &mockChatStream{ctx: createContextWithToken(t)}
	err := svc.Chat(&aichatproto.ChatRequest{
		EnvironmentId: "env-1",
		Messages: []*aichatproto.ChatMessage{
			{Role: aichatproto.ChatMessage_ROLE_USER, Content: "hello"},
		},
	}, stream)

	assert.NoError(t, err)
	require.NotEmpty(t, stream.responses)
	// The mock LLM returns a single "test response" chunk with Done=true
	last := stream.responses[len(stream.responses)-1]
	assert.True(t, last.Done)
}

func TestGetSuggestions_EmptyEnvironmentID(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	svc := createAIChatServiceForTest(
		t, mockCtrl,
		accountproto.AccountV2_Role_Organization_ADMIN,
		accountproto.AccountV2_Role_Environment_VIEWER,
	)

	ctx := createContextWithToken(t)
	_, err := svc.GetSuggestions(ctx, &aichatproto.GetSuggestionsRequest{
		EnvironmentId: "",
	})

	assert.Error(t, err)
	st, _ := gstatus.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}
