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
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	aichatproto "github.com/bucketeer-io/bucketeer/v2/proto/aichat"
)

// Chat implements the streaming chat RPC.
func (s *AIChatService) Chat(
	req *aichatproto.ChatRequest,
	stream aichatproto.AIChatService_ChatServer,
) error {
	ctx := stream.Context()

	// Input validation (before authorization to avoid unnecessary RPC)
	if req.EnvironmentId == "" {
		return statusMissingEnvironmentID.Err()
	}
	if len(req.Messages) == 0 {
		return statusMissingMessages.Err()
	}
	if len(req.Messages) > maxMessages {
		return statusTooManyMessages.Err()
	}

	// Rate limit check (before auth to avoid unnecessary RPC on hot path)
	if s.rateLimiter != nil {
		token, ok := rpc.GetAccessToken(ctx)
		if ok && !s.rateLimiter.Allow(token.Email) {
			return statusRateLimitExceeded.Err()
		}
	}

	// Authorization check (Viewer role minimum)
	_, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId,
	)
	if err != nil {
		return err
	}

	// Stream chat
	responseChan, errChan := streamChat(ctx, s.llmClient, s.ragSearcher, s.featureClient, s.chatConfig, req, s.logger)

	for responseChan != nil || errChan != nil {
		select {
		case <-ctx.Done():
			return statusRequestCanceled.Err()
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
				continue
			}
			if err != nil {
				s.logger.Error("chat stream error", zap.Error(err))
				return statusChatFailed.Err()
			}
		case resp, ok := <-responseChan:
			if !ok {
				responseChan = nil
				continue
			}
			if err := stream.Send(resp); err != nil {
				return err
			}
			if resp.Done {
				return nil
			}
		}
	}
	return nil
}
