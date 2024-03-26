// Copyright 2024 The Bucketeer Authors.
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
//

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

type Processor interface {
	Process(msg *puller.Message)
}

type Configuration struct {
	Project                      string `json:"project"`
	Subscription                 string `json:"subscription"`
	Topic                        string `json:"topic"`
	PullerNumGoroutines          int    `json:"pullerNumGoroutines"`
	PullerMaxOutstandingMessages int    `json:"pullerMaxOutstandingMessages"`
	PullerMaxOutstandingBytes    int    `json:"pullerMaxOutstandingBytes"`
	MaxMPS                       int    `json:"maxMPS"`
}

type Subscriber struct {
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
) *Subscriber {
	options := defaultOptions
	for _, o := range opts {
		o(&options)
	}
	return &Subscriber{
		name:          name,
		configuration: configuration,
		processor:     processor,
		opts:          options,
		logger:        options.logger.Named(name),
	}
}

func (s Subscriber) Run(ctx context.Context) {
	s.logger.Info("Subscriber starting",
		zap.String("name", s.name),
		zap.String("project", s.configuration.Project),
		zap.String("subscription", s.configuration.Subscription),
		zap.String("topic", s.configuration.Topic),
	)
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	rateLimiterPuller := s.createPuller(ctx, s.configuration)
	group := errgroup.Group{}
	group.Go(func() error {
		return rateLimiterPuller.Run(ctx)
	})
	group.Go(func() error {
		for {
			select {
			case msg, ok := <-rateLimiterPuller.MessageCh():
				if !ok {
					s.logger.Error("Subscriber message channel closed",
						zap.String("name", s.name))
					return nil
				}
				s.processor.Process(msg)
			case <-ctx.Done():
				s.logger.Info("Subscriber context done, stopped processing messages",
					zap.String("name", s.name))
				return nil
			}
		}
	})
	err := group.Wait()
	if err != nil {
		s.logger.Error("Subscriber stopped with error",
			zap.String("name", s.name),
			zap.Error(err))
	}
	s.logger.Info("Subscriber stopped",
		zap.String("name", s.name))
}

func (s Subscriber) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s Subscriber) createPuller(
	ctx context.Context,
	configuration Configuration,
) puller.RateLimitedPuller {
	pubsubClient, err := pubsub.NewClient(
		ctx,
		configuration.Project,
		pubsub.WithMetrics(s.opts.metrics),
		pubsub.WithLogger(s.logger),
	)
	if err != nil {
		s.logger.Error("Failed to create pubsub client", zap.Error(err))
		return nil
	}
	pubsubPuller, err := pubsubClient.CreatePuller(
		configuration.Subscription,
		configuration.Topic,
		pubsub.WithNumGoroutines(configuration.PullerNumGoroutines),
		pubsub.WithMaxOutstandingMessages(configuration.PullerMaxOutstandingMessages),
		pubsub.WithMaxOutstandingBytes(configuration.PullerMaxOutstandingBytes),
	)
	if err != nil {
		s.logger.Error("Failed to create pubsub puller", zap.Error(err))
		return nil
	}
	return puller.NewRateLimitedPuller(pubsubPuller, configuration.MaxMPS)
}
