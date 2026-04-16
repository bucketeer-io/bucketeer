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
	"fmt"
	"strings"
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

// CreatePuller creates a puller for the given subscription and topic.
// PullerOption is accepted for interface compatibility but ignored for Redis Streams.
func (c *StreamClient) CreatePuller(subscription, topic string, _ ...puller.PullerOption) (puller.Puller, error) {
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

// DeleteSubscription removes the consumer group named subscription from every
// partition stream for topic (same layout as StreamPuller). Missing streams or
// groups are ignored. Attempts all partitions and returns a joined error if any
// non-benign destroy fails (best-effort shutdown cleanup).
func (c *StreamClient) DeleteSubscription(subscription, topic string) error {
	if subscription == "" {
		return ErrInvalidStreamSubscription
	}
	if topic == "" {
		return ErrInvalidStreamTopic
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var errs []error
	for partition := 0; partition < c.partitionCount; partition++ {
		streamKey := fmt.Sprintf("%s-%d{stream}", topic, partition)
		if err := c.redisClient.XGroupDestroy(ctx, streamKey, subscription); err != nil {
			if isBenignXGroupDestroyErr(err) {
				continue
			}
			c.logger.Warn("Failed to destroy consumer group on stream partition",
				zap.String("stream", streamKey),
				zap.String("group", subscription),
				zap.Error(err),
			)
			errs = append(errs, fmt.Errorf("partition %s: %w", streamKey, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("redis stream delete subscription: %w", errors.Join(errs...))
	}
	return nil
}

func isBenignXGroupDestroyErr(err error) bool {
	if err == nil {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "nogroup") ||
		strings.Contains(msg, "no such key") ||
		strings.Contains(msg, "stream key not found") ||
		strings.Contains(msg, "the xgroup subcommand requires the key to exist") ||
		strings.Contains(msg, "requires the key to exist")
}
