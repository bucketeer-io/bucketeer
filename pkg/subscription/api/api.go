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

	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database"
	v2 "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/storage/v2"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	subscriptionproto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
)

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type SubscriptionService struct {
	dbClient                 database.Client
	adminSubscriptionStorage v2.AdminSubscriptionStorage
	subscriptionStorage      v2.SubscriptionStorage
	domainEventPublisher     publisher.Publisher
	opts                     *options
	logger                   *zap.Logger
}

func NewSubscriptionService(
	dbClient database.Client,
	adminSubscriptionStorage v2.AdminSubscriptionStorage,
	subscriptionStorage v2.SubscriptionStorage,
	domainEventPublisher publisher.Publisher,
	opts ...Option,
) *SubscriptionService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &SubscriptionService{
		dbClient:                 dbClient,
		adminSubscriptionStorage: adminSubscriptionStorage,
		subscriptionStorage:      subscriptionStorage,
		domainEventPublisher:     domainEventPublisher,
		opts:                     dopts,
		logger:                   dopts.logger.Named("api"),
	}
}

func (s *SubscriptionService) Register(server *grpc.Server) {
	subscriptionproto.RegisterSubscriptionServiceServer(server, s)
}

func (s *SubscriptionService) publishDomainEvents(ctx context.Context, events []*eventproto.Event) map[string]error {
	messages := make([]publisher.Message, 0, len(events))
	for _, event := range events {
		messages = append(messages, event)
	}
	return s.domainEventPublisher.PublishMulti(ctx, messages)
}
