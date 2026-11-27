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
	"google.golang.org/grpc"

	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/llm"
	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/rag"
	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/ratelimit"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	aichatproto "github.com/bucketeer-io/bucketeer/v2/proto/aichat"
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
	ragSearcher   rag.Searcher
	chatConfig    ChatConfig
	featureClient featureclient.Client
	rateLimiter   *ratelimit.Limiter
	logger        *zap.Logger
}

// NewAIChatService creates a new AIChatService.
func NewAIChatService(
	llmClient llm.Client,
	ragSearcher rag.Searcher,
	chatConfig ChatConfig,
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
		ragSearcher:   ragSearcher,
		chatConfig:    defaultChatConfig(chatConfig),
		featureClient: featureClient,
		rateLimiter:   dopts.rateLimiter,
		logger:        logger.Named("api"),
	}
}

// Register registers the service with the gRPC server.
func (s *AIChatService) Register(server *grpc.Server) {
	aichatproto.RegisterAIChatServiceServer(server, s)
}
