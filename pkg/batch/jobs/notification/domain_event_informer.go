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

package notification

import (
	"context"
	"errors"
	"time"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	gcodes "google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	domaineventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
	senderproto "github.com/bucketeer-io/bucketeer/proto/notification/sender"
)

var (
	ErrUnknownSourceType = errors.New("batch-server: domain-event-informer unknown source type")
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

type domainEventInformer struct {
	environmentClient            environmentclient.Client
	sender                       sender.Sender
	group                        errgroup.Group
	maxMPS                       int
	runningDurationPerBatch      time.Duration
	project                      string
	domainSubscription           string
	domainTopic                  string
	pullerNumGoroutines          int
	pullerMaxOutstandingMessages int
	pullerMaxOutstandingBytes    int
	opts                         *options
	logger                       *zap.Logger
}

func NewDomainEventInformer(
	environmentClient environmentclient.Client,
	sender sender.Sender,
	maxMPS int,
	runningDurationPerBatch time.Duration,
	project string,
	domainSubscription string,
	domainTopic string,
	pullerNumGoroutines int,
	pullerMaxOutstandingMessages int,
	pullerMaxOutstandingBytes int,
	opts ...Option) jobs.Job {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.metrics != nil {
		registerMetrics(options.metrics)
	}
	return &domainEventInformer{
		environmentClient:            environmentClient,
		sender:                       sender,
		maxMPS:                       maxMPS,
		runningDurationPerBatch:      runningDurationPerBatch,
		project:                      project,
		domainSubscription:           domainSubscription,
		domainTopic:                  domainTopic,
		pullerNumGoroutines:          pullerNumGoroutines,
		pullerMaxOutstandingMessages: pullerMaxOutstandingMessages,
		pullerMaxOutstandingBytes:    pullerMaxOutstandingBytes,
		opts:                         &options,
		logger:                       options.logger.Named("sender"),
	}
}

func (i *domainEventInformer) Run(ctx context.Context) error {
	i.logger.Info("DomainEventInformer running")
	rateLimitedPuller, closeClient, err := i.createPuller(ctx)
	if err != nil {
		return err
	}
	cctx, cancel := context.WithCancel(ctx)
	time.AfterFunc(i.runningDurationPerBatch, func() {
		i.logger.Info(
			"DomainEventInformer stopping",
			zap.Duration("runningDurationPerBatch", i.runningDurationPerBatch),
		)
		cancel()
	})
	i.group.Go(func() error {
		return rateLimitedPuller.Run(cctx)
	})
	i.group.Go(func() error {
		for {
			select {
			case msg, ok := <-rateLimitedPuller.MessageCh():
				if !ok {
					return nil
				}
				receivedCounter.WithLabelValues(typeDomainEvent).Inc()
				i.handleMessage(msg)
			case <-cctx.Done():
				return nil
			}
		}
	})
	err = i.group.Wait()
	closeClient()
	i.logger.Info("DomainEventInformer stopped")
	return err
}

func (i *domainEventInformer) createPuller(ctx context.Context) (puller.RateLimitedPuller, func(), error) {
	pubsubClient, err := pubsub.NewClient(
		ctx,
		i.project,
		pubsub.WithMetrics(i.opts.metrics),
		pubsub.WithLogger(i.logger),
	)
	if err != nil {
		i.logger.Error("Failed to create pubsub client", zap.Error(err))
		return nil, nil, err
	}
	pubsubPuller, err := pubsubClient.CreatePuller(i.domainSubscription, i.domainTopic,
		pubsub.WithNumGoroutines(i.pullerNumGoroutines),
		pubsub.WithMaxOutstandingMessages(i.pullerMaxOutstandingMessages),
		pubsub.WithMaxOutstandingBytes(i.pullerMaxOutstandingBytes),
	)
	if err != nil {
		i.logger.Error("Failed to create pubsub puller", zap.Error(err))
		return nil, nil, err
	}
	closePubsubClient := func() {
		i.logger.Debug("Closing pubsub client",
			zap.String("subscription", i.domainSubscription),
			zap.String("topic", i.domainTopic),
		)
		if err := pubsubClient.Close(); err != nil {
			i.logger.Error("Failed to close pubsub client", zap.Error(err))
		}
	}
	rateLimitedPuller := puller.NewRateLimitedPuller(pubsubPuller, i.maxMPS)
	return rateLimitedPuller, closePubsubClient, nil
}

func (i *domainEventInformer) handleMessage(msg *puller.Message) {
	if id := msg.Attributes["id"]; id == "" {
		msg.Ack()
		handledCounter.WithLabelValues(codes.MissingID.String(), codes.BadMessage.String()).Inc()
		return
	}
	domainEvent, err := i.unmarshalMessage(msg)
	if err != nil {
		handledCounter.WithLabelValues(typeDomainEvent, codes.BadMessage.String()).Inc()
		msg.Ack()
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	environmentID := ""
	if !domainEvent.IsAdminEvent {
		environment, err := i.getEnvironment(ctx, domainEvent.EnvironmentNamespace)
		if err != nil {
			if code := gstatus.Code(err); code == gcodes.NotFound {
				handledCounter.WithLabelValues(typeDomainEvent, codes.BadMessage.String()).Inc()
				msg.Ack()
				return
			}
			handledCounter.WithLabelValues(typeDomainEvent, codes.RepeatableError.String()).Inc()
			msg.Nack()
			return
		}
		environmentID = environment.Id
	}
	ne, err := i.createNotificationEvent(domainEvent, environmentID, domainEvent.IsAdminEvent)
	if err != nil {
		handledCounter.WithLabelValues(typeDomainEvent, codes.BadMessage.String()).Inc()
		msg.Ack()
		return
	}
	if err := i.sender.Send(ctx, ne); err != nil {
		handledCounter.WithLabelValues(typeDomainEvent, codes.NonRepeatableError.String()).Inc()
		msg.Ack()
		i.logger.Error("Failed to send notification event", zap.Error(err))
		return
	}
	handledCounter.WithLabelValues(typeDomainEvent, codes.OK.String()).Inc()
	msg.Ack()
}

func (i *domainEventInformer) createNotificationEvent(
	event *domaineventproto.Event,
	environmentID string,
	isAdminEvent bool,
) (*senderproto.NotificationEvent, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	st, err := i.convSourceType(event.EntityType)
	if err != nil {
		i.logger.Error("Failed to convert source type", zap.Error(err))
		return nil, err
	}
	ne := &senderproto.NotificationEvent{
		Id:                   id.String(),
		EnvironmentNamespace: event.EnvironmentNamespace,
		SourceType:           st,
		Notification: &senderproto.Notification{
			Type: senderproto.Notification_DomainEvent,
			DomainEventNotification: &senderproto.DomainEventNotification{
				EnvironmentId: environmentID,
				Editor:        event.Editor,
				EntityType:    event.EntityType,
				EntityId:      event.EntityId,
				Type:          event.Type,
			},
		},
		IsAdminEvent: isAdminEvent,
	}
	return ne, nil
}

func (i *domainEventInformer) getEnvironment(
	ctx context.Context,
	environmentId string,
) (*environmentproto.EnvironmentV2, error) {
	resp, err := i.environmentClient.GetEnvironmentV2(ctx, &environmentproto.GetEnvironmentV2Request{
		Id: environmentId,
	})
	if err != nil {
		i.logger.Error(
			"Failed to get environment",
			zap.Error(err),
			zap.String("environmentId", environmentId),
		)
		return nil, err
	}
	return resp.Environment, nil
}

func (i *domainEventInformer) unmarshalMessage(msg *puller.Message) (*domaineventproto.Event, error) {
	event := &domaineventproto.Event{}
	err := proto.Unmarshal(msg.Data, event)
	if err != nil {
		i.logger.Error("Failed to unmarshal message", zap.Error(err), zap.String("msgID", msg.ID))
		return nil, err
	}
	return event, nil
}

func (i *domainEventInformer) convSourceType(
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
