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

package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/pkg/notification/sender/notifier"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/subscriber"
	domainevent "github.com/bucketeer-io/bucketeer/proto/event/domain"
	domaineventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
	senderproto "github.com/bucketeer-io/bucketeer/proto/notification/sender"
)

type DemoOrganizationCreationNotifierConfig struct {
	WebURL          string `json:"webURL"`
	SlackWebhookURL string `json:"slackWebhookURL"`
}

type demoOrganizationCreationNotifier struct {
	slackNotifier                          notifier.Notifier
	demoOrganizationCreationNotifierConfig DemoOrganizationCreationNotifierConfig
	logger                                 *zap.Logger
}

func NewDemoOrganizationCreationNotifier(
	config interface{},
	logger *zap.Logger,
) subscriber.PubSubProcessor {
	jsonConfigMap, ok := config.(map[string]interface{})
	if !ok {
		logger.Error("demoOrganizationCreationNotifier: invalid config type, expected map[string]interface{}")
		return nil
	}
	configBytes, err := json.Marshal(jsonConfigMap)
	if err != nil {
		logger.Error("demoOrganizationCreationNotifier: failed to marshal config", zap.Error(err))
		return nil
	}
	var notifierConfig DemoOrganizationCreationNotifierConfig
	if err := json.Unmarshal(configBytes, &notifierConfig); err != nil {
		logger.Error("demoOrganizationCreationNotifier: failed to unmarshal config", zap.Error(err))
		return nil
	}
	slackNotifier := notifier.NewSlackNotifier(notifierConfig.WebURL)

	return &demoOrganizationCreationNotifier{
		slackNotifier:                          slackNotifier,
		demoOrganizationCreationNotifierConfig: notifierConfig,
		logger:                                 logger,
	}
}

func (d demoOrganizationCreationNotifier) Process(ctx context.Context, msgChan <-chan *puller.Message) error {
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				d.logger.Error("demoOrganizationCreationNotifier: message channel closed")
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(subscriberDemoOrganizationEvent).Inc()
			d.handleMessage(msg)
		case <-ctx.Done():
			d.logger.Debug("subscriber context done, stopped processing messages")
			return nil
		}
	}
}

func (d demoOrganizationCreationNotifier) handleMessage(msg *puller.Message) {
	if id := msg.Attributes["id"]; id == "" {
		msg.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberDemoOrganizationEvent, codes.BadMessage.String()).Inc()
		return
	}
	domainEvent, err := d.unmarshalMessage(msg)
	if err != nil {
		d.logger.Error("Failed to unmarshal message",
			zap.Error(err),
			zap.String("msgID", msg.ID),
			zap.String("attributes", fmt.Sprintf("%+v", msg.Attributes)),
		)
		subscriberHandledCounter.WithLabelValues(subscriberDemoOrganizationEvent, codes.BadMessage.String()).Inc()
		msg.Ack()
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if domainEvent.Type != domainevent.Event_DEMO_ORGANIZATION_CREATED {
		msg.Ack()
		return
	}

	var organizationCreatedEvent domaineventproto.OrganizationCreatedEvent
	if err := domainEvent.Data.UnmarshalTo(&organizationCreatedEvent); err != nil {
		d.logger.Error("Failed to unmarshal OrganizationCreatedEvent",
			zap.String("event id", domainEvent.Id),
			zap.Error(err),
		)
		subscriberHandledCounter.WithLabelValues(
			subscriberDemoOrganizationEvent,
			codes.NonRepeatableError.String(),
		).Inc()
		msg.Ack()
		return
	}

	recipient := &notificationproto.Recipient{
		Type:     notificationproto.Recipient_SlackChannel,
		Language: notificationproto.Recipient_ENGLISH,
		SlackChannelRecipient: &notificationproto.SlackChannelRecipient{
			WebhookUrl: d.demoOrganizationCreationNotifierConfig.SlackWebhookURL,
		},
	}
	fmt.Printf("?%+v\n", recipient)
	err = d.slackNotifier.Notify(ctx, &senderproto.Notification{
		Type: senderproto.Notification_DemoOrganizationCreation,
		DemoOrganizationCreationNotification: &senderproto.DemoOrganizationCreationNotification{
			OwnerEmail:       organizationCreatedEvent.OwnerEmail,
			OrganizationId:   organizationCreatedEvent.Id,
			OrganizationName: organizationCreatedEvent.Name,
		},
	}, recipient, recipient.Language)
	if err != nil {
		d.logger.Error("Failed to send notification",
			zap.Error(err),
			zap.String("event id", domainEvent.Id),
			zap.String("webhookURL", d.demoOrganizationCreationNotifierConfig.SlackWebhookURL),
			zap.String("organizationId", organizationCreatedEvent.Id),
		)
		subscriberHandledCounter.WithLabelValues(
			subscriberDemoOrganizationEvent,
			codes.NonRepeatableError.String(),
		).Inc()
		msg.Ack()
		return
	}
	fmt.Printf("?2")
	subscriberHandledCounter.WithLabelValues(
		subscriberDemoOrganizationEvent,
		codes.OK.String(),
	).Inc()
	msg.Ack()
}

func (d demoOrganizationCreationNotifier) unmarshalMessage(msg *puller.Message) (*domainevent.Event, error) {
	event := &domaineventproto.Event{}
	err := proto.Unmarshal(msg.Data, event)
	if err != nil {
		d.logger.Error("Failed to unmarshal message", zap.Error(err), zap.String("msgID", msg.ID))
		return nil, err
	}
	return event, nil
}
