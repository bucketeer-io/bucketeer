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
package publisher

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	idAttribute = "id"
)

var (
	ErrBadMessage = errors.New("publisher: bad message")
)

type Message interface {
	GetId() string
	proto.Message
}

type Publisher interface {
	Publish(ctx context.Context, msg Message) error
	PublishMulti(ctx context.Context, messages []Message) map[string]error
	PublishWithOrdering(ctx context.Context, msg *OrderingMessage) error
	PublishMultiWithOrdering(ctx context.Context, messages []*OrderingMessage) map[string]error
	Stop()
}

type publisher struct {
	topic  *pubsub.Topic
	logger *zap.Logger
}

type options struct {
	metrics metrics.Registerer
	logger  *zap.Logger
}

type Option func(*options)

func WithMetrics(registerer metrics.Registerer) Option {
	return func(opts *options) {
		opts.metrics = registerer
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func NewPublisher(topic *pubsub.Topic, opts ...Option) Publisher {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	if dopts.metrics != nil {
		registerMetrics(dopts.metrics)
	}
	return &publisher{
		topic:  topic,
		logger: dopts.logger.Named("publisher"),
	}
}

func (p *publisher) publishMessage(ctx context.Context, msg Message, orderingKey string) (err error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		p.logger.Error("Failed to marshal message", zap.Error(err), zap.Any("message", msg))
		return ErrBadMessage
	}
	p.topic.EnableMessageOrdering = orderingKey != ""
	message := &pubsub.Message{
		Data:       data,
		Attributes: map[string]string{idAttribute: msg.GetId()},
	}
	if orderingKey != "" {
		message.OrderingKey = orderingKey
	}
	res := p.topic.Publish(ctx, message)
	_, err = res.Get(ctx)
	return
}

func (p *publisher) publishMultiMessages(
	ctx context.Context,
	messages []Message,
	orderingKeys map[string]string,
) (errors map[string]error) {
	errors = make(map[string]error)
	results := make(map[string]*pubsub.PublishResult, len(messages))
	p.topic.EnableMessageOrdering = len(orderingKeys) > 0

	for _, msg := range messages {
		id := msg.GetId()
		data, err := proto.Marshal(msg)
		if err != nil {
			p.logger.Error("Failed to marshal message", zap.Error(err), zap.Any("message", msg))
			errors[id] = ErrBadMessage
			continue
		}
		message := &pubsub.Message{
			Data:       data,
			Attributes: map[string]string{idAttribute: id},
		}
		if orderingKey, ok := orderingKeys[id]; ok {
			message.OrderingKey = orderingKey
		}
		results[id] = p.topic.Publish(ctx, message)
	}

	for id, result := range results {
		if _, err := result.Get(ctx); err != nil {
			errors[id] = err
		}
	}
	return
}

func (p *publisher) Publish(ctx context.Context, msg Message) (err error) {
	startTime := time.Now()
	defer func() {
		topicID := p.topic.ID()
		code := convertErrorToCode(err)
		handledCounter.WithLabelValues(topicID, methodPublish, code).Inc()
		handledHistogram.WithLabelValues(topicID, methodPublish, code).Observe(time.Since(startTime).Seconds())
	}()
	return p.publishMessage(ctx, msg, "")
}

func (p *publisher) PublishMulti(ctx context.Context, messages []Message) (errors map[string]error) {
	startTime := time.Now()
	defer func() {
		topicID := p.topic.ID()
		for _, err := range errors {
			code := convertErrorToCode(err)
			handledCounter.WithLabelValues(topicID, methodPublishMulti, code).Inc()
		}
		if successes := len(messages) - len(errors); successes > 0 {
			handledCounter.WithLabelValues(topicID, methodPublishMulti, codeOK).Add(float64(successes))
		}
		histogramCode := codeOK
		if len(errors) > 0 {
			histogramCode = codeUnknown
		}
		handledHistogram.WithLabelValues(topicID, methodPublishMulti, histogramCode).Observe(time.Since(startTime).Seconds())
	}()
	return p.publishMultiMessages(ctx, messages, nil)
}

func (p *publisher) PublishWithOrdering(ctx context.Context, msg *OrderingMessage) (err error) {
	startTime := time.Now()
	defer func() {
		topicID := p.topic.ID()
		code := convertErrorToCode(err)
		handledCounter.WithLabelValues(topicID, methodPublishWithOrderingKey, code).Inc()
		handledHistogram.WithLabelValues(topicID, methodPublishWithOrderingKey, code).Observe(time.Since(startTime).Seconds())
	}()
	return p.publishMessage(ctx, msg.Message, msg.OrderingKey)
}

func (p *publisher) PublishMultiWithOrdering(
	ctx context.Context,
	messages []*OrderingMessage,
) (errors map[string]error) {
	startTime := time.Now()
	defer func() {
		topicID := p.topic.ID()
		for _, err := range errors {
			code := convertErrorToCode(err)
			handledCounter.WithLabelValues(topicID, methodPublishMultiWithOrderingKey, code).Inc()
		}
		if successes := len(messages) - len(errors); successes > 0 {
			handledCounter.WithLabelValues(topicID, methodPublishMultiWithOrderingKey, codeOK).Add(float64(successes))
		}
		histogramCode := codeOK
		if len(errors) > 0 {
			histogramCode = codeUnknown
		}
		handledHistogram.WithLabelValues(
			topicID,
			methodPublishMultiWithOrderingKey,
			histogramCode,
		).Observe(time.Since(startTime).Seconds())
	}()

	// Convert OrderingMessages to Messages and create ordering key map
	msgs := make([]Message, len(messages))
	orderingKeys := make(map[string]string, len(messages))
	for i, msg := range messages {
		msgs[i] = msg.Message
		orderingKeys[msg.Message.GetId()] = msg.OrderingKey
	}
	return p.publishMultiMessages(ctx, msgs, orderingKeys)
}

func (p *publisher) Stop() {
	p.topic.Stop()
}

type OrderingMessage struct {
	Message     Message
	OrderingKey string
}

func NewOrderingMessage(msg Message, orderingKey string) *OrderingMessage {
	return &OrderingMessage{
		Message:     msg,
		OrderingKey: orderingKey,
	}
}
