// Copyright 2023 The Bucketeer Authors.
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

package pubsub

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/backoff"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
)

var (
	ErrInvalidTopic = errors.New("pubsub: invalid topic")
)

type Client struct {
	*pubsub.Client
	opts   *options
	logger *zap.Logger
}

type options struct {
	backoff backoff.Backoff
	retries int
	metrics metrics.Registerer
	logger  *zap.Logger
}

func defaultOptions() *options {
	return &options{
		backoff: backoff.NewExponential(time.Second, 20*time.Second),
		retries: 3,
		logger:  zap.NewNop(),
	}
}

type Option func(*options)

func WithBackoff(bf backoff.Backoff) Option {
	return func(opts *options) {
		opts.backoff = bf
	}
}

func WithRetries(retries int) Option {
	return func(opts *options) {
		opts.retries = retries
	}
}

func WithMetrics(registerer metrics.Registerer) Option {
	return func(opts *options) {
		opts.metrics = registerer
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type receiveOptions = pubsub.ReceiveSettings

type ReceiveOption func(*receiveOptions)

func WithMaxExtension(d time.Duration) ReceiveOption {
	return func(opts *receiveOptions) {
		opts.MaxExtension = d
	}
}

func WithMaxOutstandingMessages(n int) ReceiveOption {
	return func(opts *receiveOptions) {
		opts.MaxOutstandingMessages = n
	}
}

func WithMaxOutstandingBytes(b int) ReceiveOption {
	return func(opts *receiveOptions) {
		opts.MaxOutstandingBytes = b
	}
}

func WithNumGoroutines(n int) ReceiveOption {
	return func(opts *receiveOptions) {
		opts.NumGoroutines = n
	}
}

type publishOptions = pubsub.PublishSettings

type PublishOption func(*publishOptions)

func WithPublishNumGoroutines(n int) PublishOption {
	return func(opts *publishOptions) {
		opts.NumGoroutines = n
	}
}

func WithPublishTimeout(timeout time.Duration) PublishOption {
	return func(opts *publishOptions) {
		opts.Timeout = timeout
	}
}

func NewClient(ctx context.Context, project string, opts ...Option) (*Client, error) {
	c, err := pubsub.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &Client{
		Client: c,
		opts:   options,
		logger: options.logger.Named("pubsub"),
	}, nil
}

func (c *Client) CreatePublisher(topic string, opts ...PublishOption) (publisher.Publisher, error) {
	t, err := c.topic(topic)
	if err != nil {
		c.logger.Error("Failed to create topic",
			zap.String("topic", topic),
			zap.Error(err))
		return nil, err
	}
	return c.createPublisher(t, opts...)
}

func (c *Client) CreatePublisherInProject(topic, project string, opts ...PublishOption) (publisher.Publisher, error) {
	t, err := c.topicInProject(topic, project)
	if err != nil {
		c.logger.Error("Failed to create topic",
			zap.String("topic", topic),
			zap.String("project", project),
			zap.Error(err))
		return nil, err
	}
	return c.createPublisher(t, opts...)
}

func (c *Client) createPublisher(topic *pubsub.Topic, opts ...PublishOption) (publisher.Publisher, error) {
	settings := (publishOptions)(pubsub.DefaultPublishSettings)
	for _, opt := range opts {
		opt(&settings)
	}
	topic.PublishSettings = settings
	options := []publisher.Option{publisher.WithLogger(c.logger)}
	if c.opts.metrics != nil {
		options = append(options, publisher.WithMetrics(c.opts.metrics))
	}
	return publisher.NewPublisher(topic, options...), nil
}

func (c *Client) CreatePuller(subscription, topic string, opts ...ReceiveOption) (puller.Puller, error) {
	s, err := c.subscription(subscription, topic)
	if err != nil {
		c.logger.Error("Failed to create puller",
			zap.String("subscription", subscription),
			zap.String("topic", topic),
			zap.Error(err))
		return nil, err
	}
	options := (receiveOptions)(pubsub.DefaultReceiveSettings)
	for _, opt := range opts {
		opt(&options)
	}
	s.ReceiveSettings = options
	c.logger.Info("Create a new puller", zap.Any("receiveSettings", options))
	return puller.NewPuller(
		s,
		puller.WithLogger(c.logger),
	), nil
}

func (c *Client) topic(id string) (*pubsub.Topic, error) {
	topic := c.Client.Topic(id)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ok, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if ok {
		return topic, nil
	}
	return nil, ErrInvalidTopic
}

func (c *Client) topicInProject(topicID, projectID string) (*pubsub.Topic, error) {
	topic := c.Client.TopicInProject(topicID, projectID)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ok, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if ok {
		return topic, nil
	}
	return nil, ErrInvalidTopic
}

// TODO: add metrics
func (c *Client) subscription(id, topicID string) (*pubsub.Subscription, error) {
	sub := c.Client.Subscription(id)
	topic := c.Client.Topic(topicID)
	var lastErr error
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	retry := backoff.NewRetry(ctx, c.opts.retries, c.opts.backoff.Clone())
	for retry.WaitNext() {
		ok, err := sub.Exists(ctx)
		if err != nil {
			continue
		}
		if ok {
			return sub, nil
		}
		_, err = c.Client.CreateSubscription(ctx, id, pubsub.SubscriptionConfig{
			Topic: topic,
		})
		if err == nil {
			return sub, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func (c *Client) SubscriptionExists(id string) (bool, error) {
	sub := c.Client.Subscription(id)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	exists, err := sub.Exists(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (c *Client) SubscriptionDetached(id string) (bool, error) {
	sub := c.Client.Subscription(id)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conf, err := sub.Config(ctx)
	if err != nil {
		return false, err
	}
	return conf.Detached, nil
}

func (c *Client) DeleteSubscription(id string) error {
	sub := c.Client.Subscription(id)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := sub.Delete(ctx)
	return err
}

func (c *Client) DeleteSubscriptionIfExist(id string) error {
	exists, err := c.SubscriptionExists(id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("subscription %s does not exist", id)
	}
	return c.DeleteSubscription(id)
}

func (c *Client) DetachSubscription(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Detach the subscription
	_, err := c.Client.DetachSubscription(ctx, name)
	if err != nil {
		return err
	}
	return nil
}
