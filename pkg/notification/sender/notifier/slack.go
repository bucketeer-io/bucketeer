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

package notifier

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"google.golang.org/grpc/metadata"

	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	featuredomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	notificationdomain "github.com/bucketeer-io/bucketeer/pkg/notification/domain"
	domainproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
	"github.com/bucketeer-io/bucketeer/proto/notification/sender"
	senderproto "github.com/bucketeer-io/bucketeer/proto/notification/sender"
)

const (
	linkTemplate = "<%s|%s>"
)

var (
	ErrUnknownNotification = errors.New("slacknotifier: unknown notification")
	ErrInvalidLanguage     = errors.New("slacknotifier: invalid language")
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

type slackNotifier struct {
	webURL string
	logger *zap.Logger
	opts   *options
}

func NewSlackNotifier(webURL string, opts ...Option) Notifier {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.metrics != nil {
		registerMetrics(options.metrics)
	}
	return &slackNotifier{
		webURL: webURL,
		opts:   &options,
		logger: options.logger.Named("sender"),
	}
}

func (n *slackNotifier) Notify(
	ctx context.Context,
	notification *senderproto.Notification,
	recipient *notificationproto.Recipient,
	language notificationproto.Recipient_Language,
) error {
	if recipient.Type != notificationproto.Recipient_SlackChannel {
		return nil
	}
	receivedCounter.WithLabelValues(typeSlack).Inc()
	if err := n.notify(ctx, notification, recipient.SlackChannelRecipient, language); err != nil {
		n.logger.Error("Failed to notify",
			zap.Error(err),
		)
		handledCounter.WithLabelValues(typeSlack, codeFail).Inc()
		return err
	}
	handledCounter.WithLabelValues(typeSlack, codeSuccess).Inc()
	return nil
}

func (n *slackNotifier) notify(
	ctx context.Context,
	notification *sender.Notification,
	slackRecipient *notificationproto.SlackChannelRecipient,
	language notificationproto.Recipient_Language,
) error {
	localizer, err := n.newLocalizer(ctx, language)
	if err != nil {
		return err
	}
	msg, err := n.createMessage(notification, slackRecipient, localizer)
	if err != nil {
		return err
	}
	if err = n.postWebhook(ctx, msg, slackRecipient.WebhookUrl); err != nil {
		// FIXME: Retry?
		return err
	}
	return nil
}

func (n *slackNotifier) newLocalizer(
	ctx context.Context,
	language notificationproto.Recipient_Language,
) (locale.Localizer, error) {
	var l string
	switch language {
	case notificationproto.Recipient_JAPANESE:
		l = locale.Ja
	case notificationproto.Recipient_ENGLISH:
		l = locale.En
	default:
		return nil, ErrInvalidLanguage
	}
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{l},
	})
	return locale.NewLocalizer(ctx), nil
}

func (n *slackNotifier) createMessage(
	notification *sender.Notification,
	slackRecipient *notificationproto.SlackChannelRecipient,
	localizer locale.Localizer,
) (*slack.WebhookMessage, error) {
	attachment, err := n.createAttachment(notification, localizer)
	if err != nil {
		return nil, err
	}
	msg := &slack.WebhookMessage{
		Attachments: []slack.Attachment{*attachment},
	}
	return msg, nil
}

func (n *slackNotifier) createAttachment(
	notification *sender.Notification,
	localizer locale.Localizer,
) (*slack.Attachment, error) {
	switch notification.Type {
	case sender.Notification_DomainEvent:
		return n.createDomainEventAttachment(notification.DomainEventNotification, localizer)
	case sender.Notification_FeatureStale:
		return n.createFeatureStaleAttachment(notification.FeatureStaleNotification)
	case sender.Notification_ExperimentRunning:
		return n.createExperimentRunningAttachment(notification.ExperimentRunningNotification)
	case sender.Notification_MauCount:
		return n.createMAUCountAttachment(notification.MauCountNotification)
	}
	return nil, ErrUnknownNotification
}

func (n *slackNotifier) createDomainEventAttachment(
	notification *senderproto.DomainEventNotification,
	localizer locale.Localizer,
) (*slack.Attachment, error) {
	// handle loc if multi-lang is necessary
	localizedMessage := domainevent.LocalizedMessage(notification.Type, localizer)
	url, err := domainevent.URL(
		notification.EntityType,
		n.webURL,
		notification.EnvironmentId,
		notification.EntityId,
	)
	if err != nil {
		return nil, err
	}
	attachment := &slack.Attachment{
		Color:      "#36a64f",
		AuthorName: notification.Editor.Email,
		Text: localizedMessage.Message + "\n\n" +
			"Environment: " + notification.EnvironmentId + "\n" +
			"Entity ID: " + notification.EntityId + "\n" +
			"URL: " + url,
	}
	return attachment, nil
}

func (n *slackNotifier) createFeatureStaleAttachment(
	notification *senderproto.FeatureStaleNotification,
) (*slack.Attachment, error) {
	featureListMsg := ""
	for _, feature := range notification.Features {
		url, err := domainevent.URL(
			domainproto.Event_FEATURE,
			n.webURL,
			notification.EnvironmentId,
			feature.Id,
		)
		if err != nil {
			return nil, err
		}
		newLine := "- ID: `" + feature.Id + "`, Name: *" + fmt.Sprintf(linkTemplate, url, feature.Name) + "*\n"
		featureListMsg = featureListMsg + newLine
	}
	// handle loc if multi-lang is necessary
	msg, err := localizedMessage(msgTypeFeatureStale, locale.Ja)
	if err != nil {
		return nil, err
	}
	replacedMsg := fmt.Sprintf(msg.Message, featuredomain.SecondsToStale/24/60/60)
	attachment := &slack.Attachment{
		Color:      "#F4D03F",
		MarkdownIn: []string{"text"},
		Text: replacedMsg + "\n\n" +
			"Environment: " + notification.EnvironmentId + "\n\n" +
			"Feature flags: \n\n" +
			featureListMsg,
	}
	return attachment, nil
}

func (n *slackNotifier) createExperimentRunningAttachment(
	notification *senderproto.ExperimentRunningNotification,
) (*slack.Attachment, error) {
	listMsg := ""
	now := time.Now()
	for _, e := range notification.Experiments {
		url, err := domainevent.URL(
			domainproto.Event_EXPERIMENT,
			n.webURL,
			notification.EnvironmentId,
			e.Id,
		)
		if err != nil {
			return nil, err
		}
		nameLink := fmt.Sprintf(linkTemplate, url, e.Name)
		newLine := fmt.Sprintf("- 残り `%d` 日, Name: *%s*\n", lastDays(now, time.Unix(e.StopAt, 0)), nameLink)
		listMsg = listMsg + newLine
	}
	// handle loc if multi-lang is necessary
	msg, err := localizedMessage(msgTypeExperimentResult, locale.Ja)
	if err != nil {
		return nil, err
	}
	attachment := &slack.Attachment{
		Color:      "#3498DB",
		MarkdownIn: []string{"text"},
		Text: msg.Message + "\n\n" +
			"Environment: " + notification.EnvironmentId + "\n\n" +
			"Experiments: \n\n" +
			listMsg,
	}
	return attachment, nil
}

func (n *slackNotifier) createMAUCountAttachment(
	notification *senderproto.MauCountNotification,
) (*slack.Attachment, error) {
	msg, err := localizedMessage(msgTypeMAUCount, locale.Ja)
	if err != nil {
		return nil, err
	}
	replacedMsg := fmt.Sprintf(msg.Message, notification.Month)
	p := message.NewPrinter(language.English)
	attachment := &slack.Attachment{
		Color:      "#3498DB",
		MarkdownIn: []string{"text"},
		Text: replacedMsg + "\n\n" +
			"Environment: " + notification.EnvironmentId + "\n" +
			p.Sprintf("Event count: %d", notification.EventCount) + "\n" +
			p.Sprintf("User count: %d", notification.UserCount),
	}
	return attachment, nil
}

func lastDays(now, stopAt time.Time) int {
	return int(stopAt.Sub(now).Hours() / 24)
}

func (n *slackNotifier) postWebhook(ctx context.Context, msg *slack.WebhookMessage, webhookURL string) error {
	if err := slack.PostWebhook(webhookURL, msg); err != nil {
		n.logger.Error("Failed to post a message",
			zap.Error(err),
			// Avoid logging a webhook URL which contains secret.
			zap.String("slackRecipientId", notificationdomain.SlackChannelRecipientID(webhookURL)),
		)
		return err
	}
	return nil
}
