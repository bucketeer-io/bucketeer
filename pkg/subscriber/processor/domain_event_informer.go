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
	"errors"
	"time"

	"go.uber.org/zap"
	gcodes "google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/subscriber"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	domaineventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
	senderproto "github.com/bucketeer-io/bucketeer/proto/notification/sender"
)

var (
	ErrUnknownSourceType = errors.New("batch-server: domain-event-informer unknown source type")
)

type domainEventInformer struct {
	environmentClient environmentclient.Client
	sender            sender.Sender
	logger            *zap.Logger
}

func NewDomainEventInformer(
	environmentClient environmentclient.Client,
	sender sender.Sender,
	logger *zap.Logger,
) subscriber.Processor {
	return &domainEventInformer{
		environmentClient: environmentClient,
		sender:            sender,
		logger:            logger,
	}
}

func (d domainEventInformer) Process(ctx context.Context, msgChan <-chan *puller.Message) error {
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				d.logger.Error("domainEventInformer: message channel closed")
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(subscriberDomainEvent).Inc()
			d.handleMessage(msg)
		case <-ctx.Done():
			d.logger.Info("subscriber context done, stopped processing messages")
			return nil
		}
	}
}

func (d domainEventInformer) handleMessage(msg *puller.Message) {
	if id := msg.Attributes["id"]; id == "" {
		msg.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberDomainEvent, codes.BadMessage.String()).Inc()
		return
	}
	domainEvent, err := d.unmarshalMessage(msg)
	if err != nil {
		subscriberHandledCounter.WithLabelValues(subscriberDomainEvent, codes.BadMessage.String()).Inc()
		msg.Ack()
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	environmentName := ""
	// TODO: The environmentURLCode will be dynamic when the console v3 is ready.
	// Currently, it inserts the url code `admin` in the domain event URL util
	// https://github.com/bucketeer-io/bucketeer/blob/main/pkg/domainevent/domain/url.go#L36-L40
	environmentURLCode := ""
	if !domainEvent.IsAdminEvent {
		environment, err := d.getEnvironment(ctx, domainEvent.EnvironmentId)
		if err != nil {
			if code := gstatus.Code(err); code == gcodes.NotFound {
				subscriberHandledCounter.WithLabelValues(subscriberDomainEvent, codes.BadMessage.String()).Inc()
				msg.Ack()
				return
			}
			subscriberHandledCounter.WithLabelValues(subscriberDomainEvent, codes.RepeatableError.String()).Inc()
			msg.Nack()
			return
		}
		environmentName = environment.Name
		environmentURLCode = environment.UrlCode
	}
	ne, err := d.createNotificationEvent(
		domainEvent,
		environmentName,
		environmentURLCode,
		domainEvent.IsAdminEvent,
	)
	if err != nil {
		subscriberHandledCounter.WithLabelValues(subscriberDomainEvent, codes.BadMessage.String()).Inc()
		msg.Ack()
		return
	}
	if err := d.sender.Send(ctx, ne); err != nil {
		subscriberHandledCounter.WithLabelValues(subscriberDomainEvent, codes.NonRepeatableError.String()).Inc()
		msg.Ack()
		d.logger.Error("Failed to send notification event", zap.Error(err))
		return
	}
	subscriberHandledCounter.WithLabelValues(subscriberDomainEvent, codes.OK.String()).Inc()
	msg.Ack()
}

func (d domainEventInformer) createNotificationEvent(
	event *domaineventproto.Event,
	environmentName, environmentURLCode string,
	isAdminEvent bool,
) (*senderproto.NotificationEvent, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	st, err := d.convSourceType(event.EntityType)
	if err != nil {
		d.logger.Error("Failed to convert source type", zap.Error(err))
		return nil, err
	}
	ne := &senderproto.NotificationEvent{
		Id:            id.String(),
		EnvironmentId: event.EnvironmentId,
		SourceType:    st,
		Notification: &senderproto.Notification{
			Type: senderproto.Notification_DomainEvent,
			DomainEventNotification: &senderproto.DomainEventNotification{
				EnvironmentName:    environmentName,
				EnvironmentUrlCode: environmentURLCode,
				Editor:             event.Editor,
				EntityType:         event.EntityType,
				EntityId:           event.EntityId,
				Type:               event.Type,
			},
		},
		IsAdminEvent: isAdminEvent,
	}
	return ne, nil
}

func (d domainEventInformer) getEnvironment(
	ctx context.Context,
	environmentId string,
) (*environmentproto.EnvironmentV2, error) {
	resp, err := d.environmentClient.GetEnvironmentV2(ctx, &environmentproto.GetEnvironmentV2Request{
		Id: environmentId,
	})
	if err != nil {
		d.logger.Error(
			"Failed to get environment",
			zap.Error(err),
			zap.String("environmentId", environmentId),
		)
		return nil, err
	}
	return resp.Environment, nil
}

func (d domainEventInformer) unmarshalMessage(msg *puller.Message) (*domaineventproto.Event, error) {
	event := &domaineventproto.Event{}
	err := proto.Unmarshal(msg.Data, event)
	if err != nil {
		d.logger.Error("Failed to unmarshal message", zap.Error(err), zap.String("msgID", msg.ID))
		return nil, err
	}
	return event, nil
}

func (d domainEventInformer) convSourceType(
	entityType domaineventproto.Event_EntityType,
) (notificationproto.Subscription_SourceType, error) {
	switch entityType {
	case domaineventproto.Event_FEATURE:
		return notificationproto.Subscription_DOMAIN_EVENT_FEATURE, nil
	case domaineventproto.Event_GOAL:
		return notificationproto.Subscription_DOMAIN_EVENT_GOAL, nil
	case domaineventproto.Event_EXPERIMENT:
		return notificationproto.Subscription_DOMAIN_EVENT_EXPERIMENT, nil
	case domaineventproto.Event_ACCOUNT:
		return notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT, nil
	case domaineventproto.Event_APIKEY:
		return notificationproto.Subscription_DOMAIN_EVENT_APIKEY, nil
	case domaineventproto.Event_SEGMENT:
		return notificationproto.Subscription_DOMAIN_EVENT_SEGMENT, nil
	case domaineventproto.Event_ENVIRONMENT:
		return notificationproto.Subscription_DOMAIN_EVENT_ENVIRONMENT, nil
	case domaineventproto.Event_ADMIN_ACCOUNT:
		return notificationproto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT, nil
	case domaineventproto.Event_AUTOOPS_RULE:
		return notificationproto.Subscription_DOMAIN_EVENT_AUTOOPS_RULE, nil
	case domaineventproto.Event_PUSH:
		return notificationproto.Subscription_DOMAIN_EVENT_PUSH, nil
	case domaineventproto.Event_SUBSCRIPTION:
		return notificationproto.Subscription_DOMAIN_EVENT_SUBSCRIPTION, nil
	case domaineventproto.Event_ADMIN_SUBSCRIPTION:
		return notificationproto.Subscription_DOMAIN_EVENT_ADMIN_SUBSCRIPTION, nil
	case domaineventproto.Event_PROJECT:
		return notificationproto.Subscription_DOMAIN_EVENT_PROJECT, nil
	case domaineventproto.Event_PROGRESSIVE_ROLLOUT:
		return notificationproto.Subscription_DOMAIN_EVENT_PROGRESSIVE_ROLLOUT, nil
	case domaineventproto.Event_ORGANIZATION:
		return notificationproto.Subscription_DOMAIN_EVENT_ORGANIZATION, nil
	case domaineventproto.Event_FLAG_TRIGGER:
		return notificationproto.Subscription_DOMAIN_EVENT_FLAG_TRIGGER, nil
	}
	return notificationproto.Subscription_SourceType(0), ErrUnknownSourceType
}
