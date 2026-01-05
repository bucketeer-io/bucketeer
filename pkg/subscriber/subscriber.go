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

package subscriber

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/factory"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	redisv3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
)

type options struct {
	metrics metrics.Registerer
	logger  *zap.Logger
}

// Configuration types
const (
	// PubSubTypeGoogle is the Google Cloud PubSub implementation
	PubSubTypeGoogle = "google"
	// PubSubTypeRedisStream is the Redis Stream implementation
	PubSubTypeRedisStream = "redis-stream"

	// DefaultPubSubType is the default PubSub implementation
	DefaultPubSubType = PubSubTypeRedisStream
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
	RedisServerName     string `json:"redisServerName,omitempty"`
	RedisAddr           string `json:"redisAddr,omitempty"`
	RedisPoolSize       int    `json:"redisPoolSize,omitempty"`
	RedisMinIdle        int    `json:"redisMinIdle,omitempty"`
	RedisDB             int    `json:"redisDB,omitempty"`
	RedisPartitionCount int    `json:"redisPartitionCount,omitempty"`
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
	switch pubSubType {
	case factory.Google:
		factoryOpts = append(factoryOpts, factory.WithProjectID(s.configuration.Project))
	case factory.RedisStream:
		// Create Redis client
		redisClient, redisErr := createRedisClient(ctx, s.configuration, s.logger, s.opts.metrics)
		if redisErr != nil {
			s.logger.Error("Failed to create Redis client", zap.Error(redisErr))
			return nil
		}
		factoryOpts = append(factoryOpts, factory.WithRedisClient(redisClient))

		// Add partition count if configured
		if s.configuration.RedisPartitionCount > 0 {
			factoryOpts = append(factoryOpts, factory.WithPartitionCount(s.configuration.RedisPartitionCount))
		}
	}

	// Create the PubSub client using the factory with context.Background()
	// to ensure connections remain healthy until explicitly stopped during graceful shutdown
	pubsubClient, err = factory.NewClient(context.Background(), factoryOpts...)
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
func createRedisClient(ctx context.Context,
	conf Configuration,
	logger *zap.Logger,
	metrics metrics.Registerer) (redisv3.Client, error) {
	redisAddr := conf.RedisAddr
	if redisAddr == "" {
		return nil, fmt.Errorf("redis address is required for Redis PubSub")
	}

	redisPoolSize := 10
	if conf.RedisPoolSize > 0 {
		redisPoolSize = conf.RedisPoolSize
	}

	redisMinIdle := 3
	if conf.RedisMinIdle > 0 {
		redisMinIdle = conf.RedisMinIdle
	}

	pubSubType := conf.PubSubType
	if pubSubType == "" {
		pubSubType = string(DefaultPubSubType)
	}

	logger.Debug("Creating Redis client",
		zap.String("address", redisAddr),
		zap.Int("poolSize", redisPoolSize),
		zap.Int("minIdle", redisMinIdle),
		zap.String("serverName", conf.RedisServerName),
		zap.String("pubSubType", pubSubType),
	)

	// Create Redis client
	return redisv3.NewClient(
		redisAddr,
		redisv3.WithPoolSize(redisPoolSize),
		redisv3.WithMinIdleConns(redisMinIdle),
		redisv3.WithServerName(conf.RedisServerName),
		redisv3.WithMetrics(metrics),
		redisv3.WithLogger(logger),
	)
}
