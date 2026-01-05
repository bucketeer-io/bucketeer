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

// Package factory provides a factory for creating PubSub clients.
package factory

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/redis"
	v3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
)

// PubSubType represents the type of PubSub implementation.
type PubSubType string

const (
	// Google represents Google Cloud PubSub.
	Google PubSubType = "google"
	// RedisStream represents Redis Streams.
	RedisStream PubSubType = "redis-stream"
)

// ClientFactory represents a factory for creating PubSub clients.
type ClientFactory interface {
	// CreateClient creates a PubSub client.
	CreateClient(ctx context.Context, opts ...Option) (Client, error)
}

// Client is an interface for PubSub operations.
type Client interface {
	// CreatePublisher creates a publisher for the given topic.
	CreatePublisher(topic string) (publisher.Publisher, error)
	// CreatePublisherInProject creates a publisher for the given topic in the specified project.
	// For Redis, this behaves the same as CreatePublisher.
	CreatePublisherInProject(topic, project string) (publisher.Publisher, error)
	// CreatePuller creates a puller for the given subscription and topic.
	CreatePuller(subscription, topic string) (puller.Puller, error)
	// SubscriptionExists checks if a subscription exists.
	SubscriptionExists(subscription string) (bool, error)
	// DeleteSubscription deletes a subscription.
	DeleteSubscription(subscription string) error
	// Close closes the client.
	Close() error
}

// Options represents options for creating a PubSub client.
type options struct {
	pubSubType     PubSubType
	projectID      string
	redisClient    v3.Client
	metrics        metrics.Registerer
	logger         *zap.Logger
	partitionCount int
}

// Option is a function that configures options.
type Option func(*options)

// WithPubSubType sets the PubSub type.
func WithPubSubType(pubSubType PubSubType) Option {
	return func(opts *options) {
		opts.pubSubType = pubSubType
	}
}

// WithProjectID sets the Google Cloud project ID.
func WithProjectID(projectID string) Option {
	return func(opts *options) {
		opts.projectID = projectID
	}
}

// WithRedisClient sets the Redis client.
func WithRedisClient(redisClient v3.Client) Option {
	return func(opts *options) {
		opts.redisClient = redisClient
	}
}

// WithMetrics sets the metrics registerer.
func WithMetrics(metrics metrics.Registerer) Option {
	return func(opts *options) {
		opts.metrics = metrics
	}
}

// WithLogger sets the logger.
func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

// WithPartitionCount sets the number of partitions for Redis Streams
func WithPartitionCount(count int) Option {
	return func(opts *options) {
		opts.partitionCount = count
	}
}

// NewClient creates a new PubSub client based on the provided options.
func NewClient(ctx context.Context, opts ...Option) (Client, error) {
	options := &options{
		pubSubType: Google, // Default to Google PubSub
		logger:     zap.NewNop(),
	}

	for _, opt := range opts {
		opt(options)
	}

	switch options.pubSubType {
	case Google:
		if options.projectID == "" {
			return nil, fmt.Errorf("project ID is required for Google PubSub")
		}
		var googleOpts []pubsub.Option
		if options.metrics != nil {
			googleOpts = append(googleOpts, pubsub.WithMetrics(options.metrics))
		}
		if options.logger != nil {
			googleOpts = append(googleOpts, pubsub.WithLogger(options.logger))
		}

		// Create Google PubSub client
		client, err := pubsub.NewClient(ctx, options.projectID, googleOpts...)
		if err != nil {
			return nil, err
		}

		// Wrap in adapter
		return &GoogleClientAdapter{
			client: client,
			logger: options.logger.Named("google-pubsub-adapter"),
		}, nil

	case RedisStream:
		// Handle RedisStream type
		if options.redisClient == nil {
			return nil, fmt.Errorf("redis client is required for Redis Stream")
		}
		streamOpts := []redis.StreamOption{}
		if options.metrics != nil {
			streamOpts = append(streamOpts, redis.WithStreamMetrics(options.metrics))
		}
		if options.logger != nil {
			streamOpts = append(streamOpts, redis.WithStreamLogger(options.logger))
		}
		if options.partitionCount > 0 {
			streamOpts = append(streamOpts, redis.WithStreamPartitionCount(options.partitionCount))
		}

		// Create Redis Stream client
		client, err := redis.NewStreamClient(ctx, options.redisClient, streamOpts...)
		if err != nil {
			return nil, err
		}

		// Redis Stream client already implements our interface
		return client, nil

	default:
		return nil, fmt.Errorf("unsupported PubSub type: %s", options.pubSubType)
	}
}

// googleClientAdapter adapts the Google PubSub client to our interface.
type GoogleClientAdapter struct {
	client *pubsub.Client
	logger *zap.Logger
}

// CreatePublisher creates a publisher for the given topic.
func (a *GoogleClientAdapter) CreatePublisher(topic string) (publisher.Publisher, error) {
	// Google PubSub requires PublishOptions, but we don't expose them in our interface
	// Use default options
	return a.client.CreatePublisher(topic)
}

// CreatePublisherInProject creates a publisher for the given topic in the specified project.
func (a *GoogleClientAdapter) CreatePublisherInProject(topic, project string) (publisher.Publisher, error) {
	// Google PubSub requires PublishOptions, but we don't expose them in our interface
	// Use default options
	return a.client.CreatePublisherInProject(topic, project)
}

// CreatePuller creates a puller for the given subscription and topic.
func (a *GoogleClientAdapter) CreatePuller(subscription, topic string) (puller.Puller, error) {
	// Google PubSub requires ReceiveOptions, but we don't expose them in our interface
	// Use default options
	return a.client.CreatePuller(subscription, topic)
}

// Close closes the client.
func (a *GoogleClientAdapter) Close() error {
	return a.client.Close()
}

// SubscriptionExists checks if a subscription exists.
func (a *GoogleClientAdapter) SubscriptionExists(subscription string) (bool, error) {
	return a.client.SubscriptionExists(subscription)
}

// DeleteSubscription deletes a subscription.
func (a *GoogleClientAdapter) DeleteSubscription(subscription string) error {
	return a.client.DeleteSubscription(subscription)
}

// GoogleClient returns the underlying Google PubSub client.
func (a *GoogleClientAdapter) GoogleClient() *pubsub.Client {
	return a.client
}
