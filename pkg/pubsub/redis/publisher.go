// Copyright 2025 The Bucketeer Authors.
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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package redis

import (
	"context"
	"errors"

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	v3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
)

var (
	// ErrRedisPublisherBadMessage is returned when a message cannot be marshaled
	ErrRedisPublisherBadMessage = errors.New("redis.publisher: bad message")
)

// Publisher RedisPublisher implements the publisher.Publisher interface for Redis
type Publisher struct {
	redisClient v3.Client
	topic       string
	logger      *zap.Logger
	metrics     metrics.Registerer
}

type PublisherOption func(*Publisher)

// WithRedisPublisherLogger sets the logger for the publisher
func WithRedisPublisherLogger(logger *zap.Logger) PublisherOption {
	return func(p *Publisher) {
		p.logger = logger
	}
}

// WithRedisPublisherMetrics sets the metrics registerer for the publisher
func WithRedisPublisherMetrics(registerer metrics.Registerer) PublisherOption {
	return func(p *Publisher) {
		p.metrics = registerer
	}
}

// NewRedisPublisher creates a new Redis publisher
func NewRedisPublisher(client v3.Client, topic string, opts ...PublisherOption) publisher.Publisher {
	p := &Publisher{
		redisClient: client,
		topic:       topic,
		logger:      zap.NewNop(),
	}

	for _, opt := range opts {
		opt(p)
	}

	p.logger = p.logger.Named("redis-publisher")

	return p
}

// Publish publishes a message to the topic
func (p *Publisher) Publish(ctx context.Context, msg publisher.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		p.logger.Error("Failed to marshal message", zap.Error(err), zap.Any("message", msg))
		return ErrRedisPublisherBadMessage
	}

	// Publish the message
	_, err = p.redisClient.Publish(ctx, p.topic, data)
	if err != nil {
		p.logger.Error("Failed to publish message",
			zap.Error(err),
			zap.String("topic", p.topic),
			zap.String("id", msg.GetId()),
		)
		return err
	}

	return nil
}

// PublishMulti publishes multiple messages
func (p *Publisher) PublishMulti(ctx context.Context, messages []publisher.Message) map[string]error {
	errors := make(map[string]error)

	for _, msg := range messages {
		id := msg.GetId()
		if err := p.Publish(ctx, msg); err != nil {
			errors[id] = err
		}
	}

	return errors
}

// Stop stops the publisher
func (p *Publisher) Stop() {
	// Redis publisher doesn't need cleanup
}
