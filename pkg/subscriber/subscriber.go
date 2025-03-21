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

package subscriber

import (
	"context"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/factory"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
)

type options struct {
	metrics metrics.Registerer
	logger  *zap.Logger
}

// Default values
const (
	// DefaultPubSubType is the default PubSub implementation
	DefaultPubSubType = factory.Google
)

type Option func(*options)

func WithMetrics(r metrics.Registerer) Option {
	return func(o *options) {
		o.metrics = r
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

type Subscriber interface {
	Run(ctx context.Context)
	Stop()
}

type PubSubProcessor interface {
	Process(ctx context.Context, msgChan <-chan *puller.Message) error
}

type OnDemandProcessor interface {
	PubSubProcessor
	Switch(ctx context.Context) (bool, error)
}

type Configuration struct {
	// PubSubType specifies which PubSub implementation to use (google or redis)
	PubSubType                   string `json:"pubSubType"`
	Project                      string `json:"project"`
	Subscription                 string `json:"subscription"`
	Topic                        string `json:"topic"`
	PullerNumGoroutines          int    `json:"pullerNumGoroutines"`
	PullerMaxOutstandingMessages int    `json:"pullerMaxOutstandingMessages"`
	PullerMaxOutstandingBytes    int    `json:"pullerMaxOutstandingBytes"`
	MaxMPS                       int    `json:"maxMPS"`
	WorkerNum                    int    `json:"workerNum"`
	// Redis configuration (used when PubSubType is "redis")
	RedisAddr     string `json:"redisAddr,omitempty"`
	RedisPoolSize int    `json:"redisPoolSize,omitempty"`
	RedisMinIdle  int    `json:"redisMinIdle,omitempty"`
	RedisDB       int    `json:"redisDB,omitempty"`
}

type pubSubSubscriber struct {
	name          string
	configuration Configuration
	processor     PubSubProcessor
	cancel        context.CancelFunc
	opts          options
	logger        *zap.Logger
}

func NewPubSubSubscriber(
	name string,
	configuration Configuration,
	processor PubSubProcessor,
	opts ...Option,
) Subscriber {
	dopts := options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(&dopts)
	}
	if configuration.PubSubType == "" {
		configuration.PubSubType = string(DefaultPubSubType)
	}
	logger := dopts.logger.Named("subscriber").With(
		zap.String("name", name),
	)
	return &pubSubSubscriber{
		name:          name,
		configuration: configuration,
		processor:     processor,
		opts:          dopts,
		logger:        logger,
	}
}

func (s pubSubSubscriber) Run(ctx context.Context) {
	s.logger.Debug("subscriber starting",
		zap.String("name", s.name),
		zap.String("pubSubType", s.configuration.PubSubType),
		zap.String("project", s.configuration.Project),
		zap.String("subscription", s.configuration.Subscription),
		zap.String("topic", s.configuration.Topic),
	)
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	rateLimiterPuller := s.createPuller(ctx)
	if rateLimiterPuller == nil {
		s.logger.Error("Failed to create puller, stopping subscriber", zap.String("name", s.name))
		return
	}
	group := errgroup.Group{}
	group.Go(func() error {
		return rateLimiterPuller.Run(ctx)
	})
	for i := 0; i < s.configuration.WorkerNum; i++ {
		group.Go(func() error {
			return s.processor.Process(ctx, rateLimiterPuller.MessageCh())
		})
	}
	err := group.Wait()
	if err != nil {
		s.logger.Error("subscriber stopped with error",
			zap.String("name", s.name),
			zap.Error(err))
	}
	s.logger.Debug("subscriber stopped",
		zap.String("name", s.name))
}

func (s pubSubSubscriber) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s pubSubSubscriber) createPuller(
	ctx context.Context,
) puller.RateLimitedPuller {
	var pubsubClient factory.Client
	var err error

	// Create client based on configured PubSubType
	pubSubType := factory.PubSubType(s.configuration.PubSubType)

	factoryOpts := []factory.Option{
		factory.WithPubSubType(pubSubType),
		factory.WithMetrics(s.opts.metrics),
		factory.WithLogger(s.logger),
	}

	// Add provider-specific options
	if pubSubType == factory.Google {
		factoryOpts = append(factoryOpts, factory.WithProjectID(s.configuration.Project))
	} else if pubSubType == factory.Redis {
		// Create Redis client
		redisClient, redisErr := createRedisClient(ctx, s.configuration, s.logger)
		if redisErr != nil {
			s.logger.Error("Failed to create Redis client", zap.Error(redisErr))
			return nil
		}
		factoryOpts = append(factoryOpts, factory.WithRedisClient(redisClient))
	}

	// Create the PubSub client using the factory
	pubsubClient, err = factory.NewClient(ctx, factoryOpts...)
	if err != nil {
		s.logger.Error("Failed to create pubsub client",
			zap.Error(err),
			zap.String("pubSubType", string(pubSubType)),
		)
		return nil
	}

	// Create the puller using the client
	pubsubPuller, err := pubsubClient.CreatePuller(
		s.configuration.Subscription,
		s.configuration.Topic,
	)
	if err != nil {
		s.logger.Error("Failed to create puller",
			zap.Error(err),
			zap.String("subscription", s.configuration.Subscription),
			zap.String("topic", s.configuration.Topic),
		)
		return nil
	}

	// Create rate-limited puller (only MaxMPS is supported in the actual implementation)
	rateLimitedPuller := puller.NewRateLimitedPuller(pubsubPuller, s.configuration.MaxMPS)
	return rateLimitedPuller
}

// createRedisClient creates a Redis client from the configuration
func createRedisClient(ctx context.Context, config Configuration, logger *zap.Logger) (redisv3.Client, error) {
	redisOpts := []redisv3.Option{}

	// Add Redis options if configured
	if config.RedisPoolSize > 0 {
		redisOpts = append(redisOpts, redisv3.WithPoolSize(config.RedisPoolSize))
	}
	if config.RedisMinIdle > 0 {
		redisOpts = append(redisOpts, redisv3.WithMinIdleConns(config.RedisMinIdle))
	}

	// Add metrics and logger
	redisOpts = append(redisOpts, redisv3.WithServerName("pubsub-redis"))
	redisOpts = append(redisOpts, redisv3.WithLogger(logger.Named("redis-client")))

	// Create Redis client
	return redisv3.NewClient(
		config.RedisAddr,
		redisOpts...,
	)
}
