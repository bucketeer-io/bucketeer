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
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
)

type options struct {
	metrics metrics.Registerer
	logger  *zap.Logger
}

var defaultOptions = options{
	logger: zap.NewNop(),
}

type Option func(*options)

func WithMetrics(r metrics.Registerer) Option {
	return func(opts *options) {
		opts.metrics = r
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type Subscriber interface {
	Run(ctx context.Context)
	Stop()
}

type Processor interface {
	Process(ctx context.Context, msgChan <-chan *puller.Message) error
}

type OnDemandProcessor interface {
	Processor
	Switch(ctx context.Context) (bool, error)
}

type Configuration struct {
	Project                      string `json:"project"`
	Subscription                 string `json:"subscription"`
	Topic                        string `json:"topic"`
	PullerNumGoroutines          int    `json:"pullerNumGoroutines"`
	PullerMaxOutstandingMessages int    `json:"pullerMaxOutstandingMessages"`
	PullerMaxOutstandingBytes    int    `json:"pullerMaxOutstandingBytes"`
	MaxMPS                       int    `json:"maxMPS"`
	WorkerNum                    int    `json:"workerNum"`
}

type subscriber struct {
	name          string
	configuration Configuration
	processor     Processor
	cancel        context.CancelFunc
	opts          options
	logger        *zap.Logger
}

func NewSubscriber(
	name string,
	configuration Configuration,
	processor Processor,
	opts ...Option,
) Subscriber {
	options := defaultOptions
	for _, o := range opts {
		o(&options)
	}
	return &subscriber{
		name:          name,
		configuration: configuration,
		processor:     processor,
		opts:          options,
		logger:        options.logger.Named(name),
	}
}

func (s subscriber) Run(ctx context.Context) {
	s.logger.Debug("subscriber starting",
		zap.String("name", s.name),
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

func (s subscriber) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s subscriber) createPuller(
	ctx context.Context,
) puller.RateLimitedPuller {
	pubsubClient, err := pubsub.NewClient(
		ctx,
		s.configuration.Project,
		pubsub.WithMetrics(s.opts.metrics),
		pubsub.WithLogger(s.logger),
	)
	if err != nil {
		s.logger.Error("Failed to create pubsub client", zap.Error(err))
		return nil
	}
	pubsubPuller, err := pubsubClient.CreatePuller(
		s.configuration.Subscription,
		s.configuration.Topic,
		pubsub.WithNumGoroutines(s.configuration.PullerNumGoroutines),
		pubsub.WithMaxOutstandingMessages(s.configuration.PullerMaxOutstandingMessages),
		pubsub.WithMaxOutstandingBytes(s.configuration.PullerMaxOutstandingBytes),
	)
	if err != nil {
		s.logger.Error("Failed to create pubsub puller", zap.Error(err))
		return nil
	}
	return puller.NewRateLimitedPuller(pubsubPuller, s.configuration.MaxMPS)
}
