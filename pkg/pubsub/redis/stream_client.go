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

package redis

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	v3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
)

const (
	// Default number of partitions for Redis Streams
	defaultStreamPartitionCount = 16
)

var (
	ErrInvalidStreamTopic        = errors.New("redis stream: invalid topic")
	ErrInvalidStreamSubscription = errors.New("redis stream: invalid subscription")
)

// StreamClient is a Redis Streams-based implementation that can create publishers and pullers
type StreamClient struct {
	redisClient    v3.Client
	opts           *streamOptions
	logger         *zap.Logger
	partitionCount int
}

type streamOptions struct {
	metrics        metrics.Registerer
	logger         *zap.Logger
	partitionCount int
	idleTime       int // Idle time in seconds for pending message reclaim
}

type StreamOption func(*streamOptions)

// WithStreamMetrics sets the metrics registerer for the client
func WithStreamMetrics(registerer metrics.Registerer) StreamOption {
	return func(opts *streamOptions) {
		opts.metrics = registerer
	}
}

// WithStreamLogger sets the logger for the client
func WithStreamLogger(logger *zap.Logger) StreamOption {
	return func(opts *streamOptions) {
		opts.logger = logger
	}
}

// WithStreamPartitionCount sets the number of partitions for the streams
func WithStreamPartitionCount(count int) StreamOption {
	return func(opts *streamOptions) {
		opts.partitionCount = count
	}
}

// WithStreamIdleTime sets the idle time in seconds for pending message reclaim
func WithStreamIdleTime(idleTimeSeconds int) StreamOption {
	return func(opts *streamOptions) {
		opts.idleTime = idleTimeSeconds
	}
}

// NewStreamClient creates a new Redis Streams client
func NewStreamClient(ctx context.Context, redisClient v3.Client, opts ...StreamOption) (*StreamClient, error) {
	options := &streamOptions{
		logger:         zap.NewNop(),
		partitionCount: defaultStreamPartitionCount,
	}
	for _, opt := range opts {
		opt(options)
	}

	return &StreamClient{
		redisClient:    redisClient,
		opts:           options,
		logger:         options.logger.Named("redis-stream-pubsub"),
		partitionCount: options.partitionCount,
	}, nil
}

// CreatePublisher creates a publisher for the given topic
func (c *StreamClient) CreatePublisher(topic string) (publisher.Publisher, error) {
	if topic == "" {
		return nil, ErrInvalidStreamTopic
	}

	options := []StreamPublisherOption{
		WithStreamPublisherLogger(c.logger),
		WithStreamPublisherPartitionCount(c.partitionCount),
	}
	if c.opts.metrics != nil {
		options = append(options, WithStreamPublisherMetrics(c.opts.metrics))
	}

	return NewStreamPublisher(c.redisClient, topic, options...), nil
}

// CreatePuller creates a puller for the given subscription and topic
func (c *StreamClient) CreatePuller(subscription, topic string) (puller.Puller, error) {
	if subscription == "" {
		return nil, ErrInvalidStreamSubscription
	}

	if topic == "" {
		return nil, ErrInvalidStreamTopic
	}

	options := []StreamPullerOption{
		WithStreamPullerLogger(c.logger),
		WithStreamPullerPartitionCount(c.partitionCount),
	}
	if c.opts.metrics != nil {
		options = append(options, WithStreamPullerMetrics(c.opts.metrics))
	}
	if c.opts.idleTime > 0 {
		options = append(options, WithStreamPullerIdleTime(time.Duration(c.opts.idleTime)*time.Second))
	}

	return NewStreamPuller(c.redisClient, subscription, topic, options...), nil
}

// CreatePublisherInProject creates a publisher for the given topic
// For Redis, this behaves the same as CreatePublisher since Redis doesn't have the concept of projects
func (c *StreamClient) CreatePublisherInProject(topic, project string) (publisher.Publisher, error) {
	// For Redis, we ignore the project parameter
	return c.CreatePublisher(topic)
}

// Close closes the client
func (c *StreamClient) Close() error {
	// The client doesn't own the Redis client, so we don't close it
	return nil
}

// SubscriptionExists checks if a subscription exists
func (c *StreamClient) SubscriptionExists(subscription string) (bool, error) {
	// For Redis Streams, we need to check if the consumer group exists
	if subscription == "" {
		return false, ErrInvalidStreamSubscription
	}

	// We can't check if a consumer group exists without knowing the stream name
	// Just return true since the group will be created on demand when pulling
	return true, nil
}

// DeleteSubscription deletes a subscription
func (c *StreamClient) DeleteSubscription(subscription string) error {
	// For Redis Streams, we would need to delete the consumer group
	// This would require knowing the stream name, which we don't have
	// Just return nil to indicate success since the operation is idempotent
	return nil
}
