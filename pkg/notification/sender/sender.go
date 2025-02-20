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
package sender

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	notificationclient "github.com/bucketeer-io/bucketeer/pkg/notification/client"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender/notifier"
	ftproto "github.com/bucketeer-io/bucketeer/proto/feature"
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
	errFailedToUnmarshal      = errors.New("sender: failed to unmarshal")
	errFeatureFlagTagNotFound = errors.New("sender: feature flag tag not found")

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
			notificationEvent.EnvironmentId,
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
		// When a flag changes it must be checked before sending notifications
		send, err := s.checkForFeatureDomainEvent(
			subscription,
			notificationEvent.SourceType,
			notificationEvent.Notification.DomainEventNotification.EntityData,
		)
		if err != nil {
			return err
		}
		// Check if the subcription tag is configured in the feature flag
		// If not, we skip the notification
		if !send {
			continue
		}
		if err := s.send(
			ctx,
			notificationEvent.Notification,
			subscription.Recipient,
			subscription.Recipient.Language,
		); err != nil {
			s.logger.Error("Failed to send notification", zap.Error(err),
				zap.String("environmentId", notificationEvent.EnvironmentId),
			)
			lastErr = err
			continue
		}
		s.logger.Info("Succeeded to send notification",
			zap.String("environmentId", notificationEvent.EnvironmentId),
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
	environmentId string,
	sourceType notificationproto.Subscription_SourceType) ([]*notificationproto.Subscription, error) {

	subscriptions := []*notificationproto.Subscription{}
	cursor := ""
	for {
		resp, err := s.notificationClient.ListEnabledSubscriptions(ctx, &notificationproto.ListEnabledSubscriptionsRequest{
			EnvironmentId: environmentId,
			SourceTypes:   []notificationproto.Subscription_SourceType{sourceType},
			PageSize:      listRequestSize,
			Cursor:        cursor,
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

// When a flag changes it must be checked before sending notifications
func (s *sender) checkForFeatureDomainEvent(
	sub *notificationproto.Subscription,
	sourceType notificationproto.Subscription_SourceType,
	entityData string,
) (bool, error) {
	// Different domain event
	if sourceType != notificationproto.Subscription_DOMAIN_EVENT_FEATURE ||
		len(sub.FeatureFlagTags) == 0 {
		s.logger.Debug(
			"Sending notification. The source type is not a feature domain event or the subscription's tags are empty",
			zap.String("environmentId", sub.EnvironmentId),
			zap.String("subscriptionId", sub.Id),
			zap.String("subscriptionName", sub.Name),
			zap.Strings("subscriptionTags", sub.FeatureFlagTags),
			zap.String("entityData", entityData),
		)
		return true, nil
	}
	// Unmarshal the JSON string into the Feature message
	var feature ftproto.Feature
	if err := protojson.Unmarshal([]byte(entityData), &feature); err != nil {
		s.logger.Error("Failed to unmarshal feature message", zap.Error(err),
			zap.String("environmentId", sub.EnvironmentId),
			zap.String("subscriptionId", sub.Id),
			zap.String("subscriptionName", sub.Name),
			zap.String("entityData", entityData),
		)
		return false, errFailedToUnmarshal
	}
	// Check if the subcription tag is configured in the feature flag
	// If not, we skip the notification
	if containsTags(sub.FeatureFlagTags, feature.Tags) {
		s.logger.Debug(
			"Sending notification. Flag's tag matched with the tags configured in the subscription",
			zap.String("environmentId", sub.EnvironmentId),
			zap.String("subscriptionId", sub.Id),
			zap.String("subscriptionName", sub.Name),
			zap.Strings("subscriptionTags", sub.FeatureFlagTags),
			zap.String("entityData", entityData),
		)
		return true, nil
	}
	s.logger.Debug(
		"Skipping notification. Subscription's tags weren't found in the Flag's tags",
		zap.String("environmentId", sub.EnvironmentId),
		zap.String("subscriptionId", sub.Id),
		zap.String("subscriptionName", sub.Name),
		zap.Strings("subscriptionTags", sub.FeatureFlagTags),
		zap.String("entityData", entityData),
	)
	return false, errFeatureFlagTagNotFound
}

func containsTags(subTags []string, ftTags []string) bool {
	for _, subTag := range subTags {
		for _, ftTag := range ftTags {
			if subTag == ftTag {
				return true
			}
		}
	}
	return false
}
