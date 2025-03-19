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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package redis

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	v3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
)

// Puller RedisPuller implements the puller.Puller interface for Redis
type Puller struct {
	redisClient  v3.Client
	subscription string
	topic        string
	logger       *zap.Logger
	metrics      metrics.Registerer
	pubsub       *v3.PubSub
	closed       bool
	done         chan struct{}
}

type PullerOption func(*Puller)

// WithRedisPullerLogger sets the logger for the puller
func WithRedisPullerLogger(logger *zap.Logger) PullerOption {
	return func(p *Puller) {
		p.logger = logger
	}
}

// WithRedisPullerMetrics sets the metrics registerer for the puller
func WithRedisPullerMetrics(registerer metrics.Registerer) PullerOption {
	return func(p *Puller) {
		p.metrics = registerer
	}
}

// NewRedisPuller creates a new Redis puller
func NewRedisPuller(
	client v3.Client,
	subscription string,
	topic string,
	opts ...PullerOption,
) puller.Puller {
	p := &Puller{
		redisClient:  client,
		subscription: subscription,
		topic:        topic,
		logger:       zap.NewNop(),
		done:         make(chan struct{}),
	}

	for _, opt := range opts {
		opt(p)
	}

	p.logger = p.logger.Named("redis-puller")

	return p
}

// Pull subscribes to the topic and calls the handler function for each message.
// It blocks until context is canceled or an error occurs.
func (p *Puller) Pull(ctx context.Context, handler func(context.Context, *puller.Message)) error {
	if p.closed {
		return fmt.Errorf("redis puller is closed")
	}

	if p.pubsub != nil {
		return fmt.Errorf("redis puller is already pulling")
	}

	// Subscribe to the topic
	pubsub, err := p.redisClient.Subscribe(ctx, p.topic)
	if err != nil {
		p.logger.Error("Failed to subscribe to topic",
			zap.Error(err),
			zap.String("subscription", p.subscription),
			zap.String("topic", p.topic),
		)
		return err
	}
	p.pubsub = pubsub

	// Create a channel for messages
	messageCh := pubsub.Channel()

	// Setup mechanism to handle context cancellation
	ctxDone := ctx.Done()

	// Process messages in the current goroutine to block until context is canceled
	// This is consistent with Google PubSub behavior
	for {
		select {
		case <-ctxDone:
			p.logger.Debug("Context canceled, stopping pull",
				zap.String("subscription", p.subscription),
				zap.String("topic", p.topic),
			)
			p.Close()
			return ctx.Err()

		case <-p.done:
			p.logger.Debug("Puller closed, stopping pull",
				zap.String("subscription", p.subscription),
				zap.String("topic", p.topic),
			)
			return nil

		case msg, ok := <-messageCh:
			if !ok {
				// Channel closed
				p.logger.Debug("Message channel closed",
					zap.String("subscription", p.subscription),
					zap.String("topic", p.topic),
				)
				return nil
			}

			// Create a message with Ack/Nack functions that do nothing for Redis
			// since Redis doesn't have message acknowledgement
			message := &puller.Message{
				ID:   fmt.Sprintf("%s:%d", msg.Channel, time.Now().UnixNano()),
				Data: msg.Payload,
				Attributes: map[string]string{
					"channel": msg.Channel,
				},
				Ack:  func() {}, // Redis PubSub doesn't require acknowledgement
				Nack: func() {}, // Redis PubSub doesn't support negative acknowledgement
			}

			// Handle message in a separate goroutine to not block receive loop
			handler(ctx, message)
		}
	}
}

// Close closes the puller
func (p *Puller) Close() error {
	if p.closed {
		return nil
	}

	p.closed = true
	close(p.done)

	if p.pubsub != nil {
		err := p.pubsub.Close()
		if err != nil {
			p.logger.Error("Failed to close Redis pubsub", zap.Error(err))
			return err
		}
		p.pubsub = nil
	}

	return nil
}

// SubscriptionName returns the name of the subscription
func (p *Puller) SubscriptionName() string {
	return fmt.Sprintf("%s:%s", p.subscription, p.topic)
}
