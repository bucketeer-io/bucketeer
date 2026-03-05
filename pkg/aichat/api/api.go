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

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	pkgapi "github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/llm"
	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/rag"
	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/ratelimit"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	aichatproto "github.com/bucketeer-io/bucketeer/v2/proto/aichat"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

type options struct {
	rateLimiter *ratelimit.Limiter
}

// Option is a functional option for AIChatService.
type Option func(*options)

// WithGRPCRateLimiter sets the rate limiter for the gRPC service.
func WithGRPCRateLimiter(l *ratelimit.Limiter) Option {
	return func(o *options) {
		o.rateLimiter = l
	}
}

// AIChatService implements the gRPC AIChatService.
type AIChatService struct {
	aichatproto.UnimplementedAIChatServiceServer
	llmClient     llm.Client
	ragService    *rag.Service
	chatConfig    ChatConfig
	accountClient accountclient.Client
	featureClient featureclient.Client
	rateLimiter   *ratelimit.Limiter
	logger        *zap.Logger
}

// NewAIChatService creates a new AIChatService.
func NewAIChatService(
	llmClient llm.Client,
	ragService *rag.Service,
	chatConfig ChatConfig,
	accountClient accountclient.Client,
	featureClient featureclient.Client,
	logger *zap.Logger,
	opts ...Option,
) *AIChatService {
	dopts := &options{}
	for _, opt := range opts {
		opt(dopts)
	}
	return &AIChatService{
		llmClient:     llmClient,
		ragService:    ragService,
		chatConfig:    defaultChatConfig(chatConfig),
		accountClient: accountClient,
		featureClient: featureClient,
		rateLimiter:   dopts.rateLimiter,
		logger:        logger.Named("api"),
	}
}

// Register registers the service with the gRPC server.
func (s *AIChatService) Register(server *grpc.Server) {
	aichatproto.RegisterAIChatServiceServer(server, s)
}

func (s *AIChatService) checkEnvironmentRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Environment,
	environmentId string,
) (*eventproto.Editor, error) {
	editor, err := role.CheckEnvironmentRole(
		ctx,
		requiredRole,
		environmentId,
		func(email string) (*accountproto.AccountV2, error) {
			resp, err := s.accountClient.GetAccountV2ByEnvironmentID(
				ctx,
				&accountproto.GetAccountV2ByEnvironmentIDRequest{
					Email:         email,
					EnvironmentId: environmentId,
				},
			)
			if err != nil {
				return nil, err
			}
			return resp.Account, nil
		},
	)
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return nil, statusUnauthenticated.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return nil, statusPermissionDenied.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return nil, pkgapi.NewGRPCStatus(err).Err()
		}
	}
	return editor, nil
}
