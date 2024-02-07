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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package sender

import (
	"context"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	notificationclient "github.com/bucketeer-io/bucketeer/pkg/notification/client"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender/notifier"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
	senderproto "github.com/bucketeer-io/bucketeer/proto/notification/sender"
)

type options struct {
	metrics metrics.Registerer
	logger  *zap.Logger
}

const (
	listRequestSize = 500
)

var (
	defaultOptions = options{
		logger: zap.NewNop(),
	}
)

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

type Sender interface {
	Send(context.Context, *senderproto.NotificationEvent) error
}

type sender struct {
	notificationClient notificationclient.Client
	notifiers          []notifier.Notifier
	opts               *options
	logger             *zap.Logger
}

func NewSender(
	notificationClient notificationclient.Client,
	notifiers []notifier.Notifier,
	opts ...Option) Sender {

	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.metrics != nil {
		registerMetrics(options.metrics)
	}
	return &sender{
		notificationClient: notificationClient,
		notifiers:          notifiers,
		opts:               &options,
		logger:             options.logger.Named("sender"),
	}
}

func (s *sender) Send(ctx context.Context, notificationEvent *senderproto.NotificationEvent) error {
	receivedCounter.Inc()
	subscriptions := []*notificationproto.Subscription{}
	if notificationEvent.IsAdminEvent {
		adminSubs, err := s.listEnabledAdminSubscriptions(ctx, notificationEvent.SourceType)
		if err != nil {
			handledCounter.WithLabelValues(codeFail).Inc()
			return err
		}
		subscriptions = append(subscriptions, adminSubs...)
	} else {
		subs, err := s.listEnabledSubscriptions(
			ctx,
			notificationEvent.EnvironmentNamespace,
			notificationEvent.SourceType,
		)
		if err != nil {
			handledCounter.WithLabelValues(codeFail).Inc()
			return err
		}
		subscriptions = append(subscriptions, subs...)
	}
	var lastErr error
	for _, subscription := range subscriptions {
		if err := s.send(
			ctx,
			notificationEvent.Notification,
			subscription.Recipient,
			subscription.Recipient.Language,
		); err != nil {
			s.logger.Error("Failed to send notification", zap.Error(err),
				zap.String("environmentNamespace", notificationEvent.EnvironmentNamespace),
			)
			lastErr = err
			continue
		}
		s.logger.Info("Succeeded to send notification",
			zap.String("environmentNamespace", notificationEvent.EnvironmentNamespace),
		)
	}
	if lastErr != nil {
		handledCounter.WithLabelValues(codeFail).Inc()
		return lastErr
	}
	handledCounter.WithLabelValues(codeSuccess).Inc()
	return nil
}

func (s *sender) send(
	ctx context.Context,
	notification *senderproto.Notification,
	recipient *notificationproto.Recipient,
	language notificationproto.Recipient_Language,
) error {
	for _, notifier := range s.notifiers {
		if err := notifier.Notify(ctx, notification, recipient, language); err != nil {
			return err
		}
	}
	return nil
}

func (s *sender) listEnabledSubscriptions(
	ctx context.Context,
	environmentNamespace string,
	sourceType notificationproto.Subscription_SourceType) ([]*notificationproto.Subscription, error) {

	subscriptions := []*notificationproto.Subscription{}
	cursor := ""
	for {
		resp, err := s.notificationClient.ListEnabledSubscriptions(ctx, &notificationproto.ListEnabledSubscriptionsRequest{
			EnvironmentNamespace: environmentNamespace,
			SourceTypes:          []notificationproto.Subscription_SourceType{sourceType},
			PageSize:             listRequestSize,
			Cursor:               cursor,
		})
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, resp.Subscriptions...)
		size := len(resp.Subscriptions)
		if size == 0 || size < listRequestSize {
			return subscriptions, nil
		}
		cursor = resp.Cursor
	}
}

func (s *sender) listEnabledAdminSubscriptions(
	ctx context.Context,
	sourceType notificationproto.Subscription_SourceType) ([]*notificationproto.Subscription, error) {

	subscriptions := []*notificationproto.Subscription{}
	cursor := ""
	for {
		resp, err := s.notificationClient.ListEnabledAdminSubscriptions(
			ctx,
			&notificationproto.ListEnabledAdminSubscriptionsRequest{
				SourceTypes: []notificationproto.Subscription_SourceType{sourceType},
				PageSize:    listRequestSize,
				Cursor:      cursor,
			},
		)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, resp.Subscriptions...)
		size := len(resp.Subscriptions)
		if size == 0 || size < listRequestSize {
			return subscriptions, nil
		}
		cursor = resp.Cursor
	}
}
