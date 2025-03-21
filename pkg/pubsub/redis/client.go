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

package redis

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	v3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
)

var (
	ErrInvalidTopic        = errors.New("redis: invalid topic")
	ErrInvalidSubscription = errors.New("redis: invalid subscription")
)

// Client is a Redis-based implementation that can create publishers and pullers
type Client struct {
	redisClient v3.Client
	opts        *options
	logger      *zap.Logger
}

type options struct {
	metrics metrics.Registerer
	logger  *zap.Logger
}

type Option func(*options)

// WithMetrics sets the metrics registerer for the client
func WithMetrics(registerer metrics.Registerer) Option {
	return func(opts *options) {
		opts.metrics = registerer
	}
}

// WithLogger sets the logger for the client
func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

// NewClient creates a new Redis client
func NewClient(ctx context.Context, redisClient v3.Client, opts ...Option) (*Client, error) {
	options := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(options)
	}

	return &Client{
		redisClient: redisClient,
		opts:        options,
		logger:      options.logger.Named("redis-pubsub"),
	}, nil
}

// CreatePublisher creates a publisher for the given topic
func (c *Client) CreatePublisher(topic string) (publisher.Publisher, error) {
	if topic == "" {
		return nil, ErrInvalidTopic
	}

	options := []PublisherOption{
		WithRedisPublisherLogger(c.logger),
	}
	if c.opts.metrics != nil {
		options = append(options, WithRedisPublisherMetrics(c.opts.metrics))
	}

	return NewRedisPublisher(c.redisClient, topic, options...), nil
}

// CreatePuller creates a puller for the given subscription and topic
func (c *Client) CreatePuller(subscription, topic string) (puller.Puller, error) {
	if subscription == "" {
		return nil, ErrInvalidSubscription
	}

	if topic == "" {
		return nil, ErrInvalidTopic
	}

	options := []PullerOption{
		WithRedisPullerLogger(c.logger),
	}
	if c.opts.metrics != nil {
		options = append(options, WithRedisPullerMetrics(c.opts.metrics))
	}

	return NewRedisPuller(c.redisClient, subscription, topic, options...), nil
}

// CreatePublisherInProject creates a publisher for the given topic
// For Redis, this behaves the same as CreatePublisher since Redis doesn't have the concept of projects
func (c *Client) CreatePublisherInProject(topic, project string) (publisher.Publisher, error) {
	// For Redis, we ignore the project parameter
	return c.CreatePublisher(topic)
}

// Close closes the client
func (c *Client) Close() error {
	// The client doesn't own the Redis client, so we don't close it
	return nil
}

// SubscriptionExists checks if a subscription exists
func (c *Client) SubscriptionExists(subscription string) (bool, error) {
	// For Redis, subscriptions don't really "exist" in the same way as Google PubSub
	// We'll return true if the subscription is valid
	if subscription == "" {
		return false, ErrInvalidSubscription
	}
	return true, nil
}

// DeleteSubscription deletes a subscription
func (c *Client) DeleteSubscription(subscription string) error {
	// In Redis, subscriptions are just runtime constructs, so there's nothing to delete
	// Return nil to indicate success
	return nil
}
