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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/publisher.go
package redis

import (
	"context"
	"errors"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	v3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
)

var (
	ErrStreamPublisherBadMessage = errors.New("redis stream publisher: bad message")
)

// StreamPublisher is a Redis Stream-based implementation of the Publisher interface
type StreamPublisher struct {
	redisClient v3.Client
	stream      string
	logger      *zap.Logger
}

type StreamPublisherOption func(*StreamPublisher)

// WithStreamPublisherMetrics sets the metrics registerer for the publisher
func WithStreamPublisherMetrics(registerer metrics.Registerer) StreamPublisherOption {
	return func(p *StreamPublisher) {
		// Redis client metrics are already registered by the Redis client
	}
}

// WithStreamPublisherLogger sets the logger for the publisher
func WithStreamPublisherLogger(logger *zap.Logger) StreamPublisherOption {
	return func(p *StreamPublisher) {
		p.logger = logger
	}
}

// NewStreamPublisher creates a new Redis Stream publisher
func NewStreamPublisher(client v3.Client, stream string, opts ...StreamPublisherOption) publisher.Publisher {
	p := &StreamPublisher{
		redisClient: client,
		stream:      stream,
		logger:      zap.NewNop(),
	}

	for _, opt := range opts {
		opt(p)
	}

	p.logger = p.logger.Named("redis-stream-publisher")

	return p
}

// Publish publishes a message to the stream
func (p *StreamPublisher) Publish(ctx context.Context, msg publisher.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		p.logger.Error("Failed to marshal message", zap.Error(err), zap.Any("message", msg))
		return ErrStreamPublisherBadMessage
	}

	// Publish the message to the stream
	// Use message ID as field name and serialized message as value
	// The * indicates Redis should auto-generate the ID
	messageID := msg.GetId()
	if messageID == "" {
		messageID = "message"
	}

	values := map[string]interface{}{
		messageID: data,
	}

	// Add the message to the stream
	_, err = p.redisClient.XAdd(ctx, p.stream, values)
	if err != nil {
		p.logger.Error("Failed to add message to stream",
			zap.Error(err),
			zap.String("stream", p.stream),
			zap.String("id", msg.GetId()),
		)
		return err
	}

	return nil
}

// PublishMulti publishes multiple messages
func (p *StreamPublisher) PublishMulti(ctx context.Context, messages []publisher.Message) map[string]error {
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
func (p *StreamPublisher) Stop() {
	// Redis stream publisher doesn't need cleanup
}
