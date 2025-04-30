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
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
)

const (
	pubsubErrNotFound = "NotFound"
)

type OnDemandConfiguration struct {
	Configuration
	CheckInterval int `json:"checkInterval"`
}

type onDemandSubscriber struct {
	name                string
	configuration       OnDemandConfiguration
	rateLimitedPuller   puller.RateLimitedPuller
	processor           OnDemandProcessor
	ctx                 context.Context
	cancel              context.CancelFunc
	runningPullerCtx    context.Context
	runningPullerCancel func()
	client              *pubsub.Client
	isRunning           bool
	group               errgroup.Group
	opts                options
	logger              *zap.Logger
}

func NewOnDemandSubscriber(
	name string,
	configuration OnDemandConfiguration,
	processor OnDemandProcessor,
	opts ...Option,
) Subscriber {
	options := defaultOptions
	for _, o := range opts {
		o(&options)
	}
	return &onDemandSubscriber{
		name:          name,
		configuration: configuration,
		processor:     processor,
		opts:          options,
		logger:        options.logger.Named(name),
	}
}

func (s *onDemandSubscriber) Run(ctx context.Context) {
	s.logger.Debug("onDemandSubscriber starting",
		zap.String("name", s.name),
		zap.String("project", s.configuration.Project),
		zap.String("subscription", s.configuration.Subscription),
		zap.String("topic", s.configuration.Topic),
	)
	s.ctx, s.cancel = context.WithCancel(ctx)
	err := s.createPubSubClient(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to create pubsub client",
			zap.Error(err),
			zap.String("name", s.name),
		)
		return
	}
	ticker := time.NewTicker(time.Duration(s.configuration.CheckInterval) * time.Second)
	defer ticker.Stop()
	subscription := make(chan struct{})
	go s.subscribe(subscription)
	for {
		select {
		case <-ticker.C:
			start, err := s.processor.Switch(ctx)
			if err != nil {
				s.logger.Error("Failed to check switch status", zap.Error(err))
				continue
			}
			if start {
				if !s.IsRunning() {
					err = s.createPuller()
					if err != nil {
						s.logger.Error("Failed to create new puller",
							zap.String("name", s.name),
							zap.Error(err),
						)
						continue
					}
					s.group = errgroup.Group{}
					subscription <- struct{}{}
				}
			} else {
				if s.IsRunning() {
					s.unsubscribe()
				}
				// delete subscription if it exists
				exists, err := s.client.SubscriptionExists(s.configuration.Subscription)
				if err != nil {
					s.logger.Error("Failed to check subscription existence",
						zap.String("name", s.name),
						zap.Error(err),
					)
					continue
				}
				if exists {
					err = s.client.DeleteSubscription(s.configuration.Subscription)
					if err != nil {
						s.logger.Error("Failed to delete subscription",
							zap.String("name", s.name),
							zap.Error(err),
						)
						continue
					}
				}
			}
		case <-ctx.Done():
			s.logger.Debug("Context is done")
			if s.IsRunning() {
				s.logger.Debug("Puller is running, stop pulling messages")
				s.unsubscribe()
			}
		}
	}
}

func (s *onDemandSubscriber) subscribe(subscription chan struct{}) {
	for {
		select {
		case <-subscription:
			s.isRunning = true
			ctx, cancel := context.WithCancel(context.Background())
			s.runningPullerCtx = ctx
			s.runningPullerCancel = cancel
			s.group.Go(func() error {
				err := s.rateLimitedPuller.Run(ctx)
				if err != nil {
					if strings.Contains(err.Error(), pubsubErrNotFound) {
						s.unsubscribe()
						return nil
					}
					s.logger.Error("Failed to pull messages", zap.Error(err))
					return err
				}
				return nil
			})
			for i := 0; i < s.configuration.WorkerNum; i++ {
				s.group.Go(s.batch)
			}
			err := s.group.Wait()
			if err != nil {
				s.logger.Error("Failed while running pull messages", zap.Error(err))
			}
			s.isRunning = false
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *onDemandSubscriber) batch() error {
	return s.processor.Process(s.runningPullerCtx, s.rateLimitedPuller.MessageCh())
}

func (s *onDemandSubscriber) unsubscribe() {
	s.runningPullerCancel()
}

func (s *onDemandSubscriber) IsRunning() bool {
	return s.isRunning
}

func (s *onDemandSubscriber) createPubSubClient(ctx context.Context) error {
	pubsubClient, err := pubsub.NewClient(
		ctx,
		s.configuration.Project,
		pubsub.WithMetrics(s.opts.metrics),
		pubsub.WithLogger(s.logger),
	)
	if err != nil {
		s.logger.Error("Failed to create pubsub client", zap.Error(err))
		return err
	}
	s.client = pubsubClient
	return nil
}

func (s *onDemandSubscriber) createPuller() error {
	pubsubPuller, err := s.client.CreatePuller(
		s.configuration.Subscription,
		s.configuration.Topic,
		pubsub.WithNumGoroutines(s.configuration.PullerNumGoroutines),
		pubsub.WithMaxOutstandingMessages(s.configuration.PullerMaxOutstandingMessages),
		pubsub.WithMaxOutstandingBytes(s.configuration.PullerMaxOutstandingBytes),
	)
	if err != nil {
		s.logger.Error("Failed to create pubsub puller", zap.Error(err))
		return err
	}
	s.rateLimitedPuller = puller.NewRateLimitedPuller(pubsubPuller, s.configuration.MaxMPS)
	return nil
}

func (s *onDemandSubscriber) Stop() {
	s.cancel()
}
