// Copyright 2022 The Bucketeer Authors.
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

package domainevent

import (
	"context"
	"errors"
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"go.uber.org/zap"
	gcodes "google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"

	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender/informer"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	domaineventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
	senderproto "github.com/bucketeer-io/bucketeer/proto/notification/sender"
)

var (
	ErrUnknownSourceType = errors.New("domain-event-informer: unknown source type")
)

type options struct {
	maxMPS     int
	numWorkers int
	metrics    metrics.Registerer
	logger     *zap.Logger
}

var defaultOptions = options{
	maxMPS:     10,
	numWorkers: 1,
	logger:     zap.NewNop(),
}

type Option func(*options)

func WithMaxMPS(mps int) Option {
	return func(opts *options) {
		opts.maxMPS = mps
	}
}

func WithNumWorkers(n int) Option {
	return func(opts *options) {
		opts.numWorkers = n
	}
}

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
	environmentClient environmentclient.Client
	puller            puller.RateLimitedPuller
	sender            sender.Sender
	group             errgroup.Group
	opts              *options
	logger            *zap.Logger
	ctx               context.Context
	cancel            func()
	doneCh            chan struct{}
}

func NewDomainEventInformer(
	environmentClient environmentclient.Client,
	p puller.Puller,
	sender sender.Sender,
	opts ...Option) informer.Informer {

	ctx, cancel := context.WithCancel(context.Background())
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.metrics != nil {
		registerMetrics(options.metrics)
	}
	return &domainEventInformer{
		environmentClient: environmentClient,
		puller:            puller.NewRateLimitedPuller(p, options.maxMPS),
		sender:            sender,
		opts:              &options,
		logger:            options.logger.Named("sender"),
		ctx:               ctx,
		cancel:            cancel,
		doneCh:            make(chan struct{}),
	}
}

func (i *domainEventInformer) Run() error {
	defer close(i.doneCh)
	i.logger.Info("DomainEventInformer start running")
	i.group.Go(func() error {
		return i.puller.Run(i.ctx)
	})
	for idx := 0; idx < i.opts.numWorkers; idx++ {
		i.group.Go(i.runWorker)
	}
	err := i.group.Wait()
	i.logger.Info("DomainEventInformer start stopping")
	return err
}

func (i *domainEventInformer) Stop() {
	i.logger.Info("DomainEventInformer start stopping")
	i.cancel()
	<-i.doneCh
}

func (i *domainEventInformer) Check(ctx context.Context) health.Status {
	select {
	case <-i.ctx.Done():
		i.logger.Error("Unhealthy due to context Done is closed", zap.Error(i.ctx.Err()))
		return health.Unhealthy
	default:
		if i.group.FinishedCount() > 0 {
			i.logger.Error("Unhealthy", zap.Int32("FinishedCount", i.group.FinishedCount()))
			return health.Unhealthy
		}
		return health.Healthy
	}
}

func (i *domainEventInformer) runWorker() error {
	for {
		select {
		case msg, ok := <-i.puller.MessageCh():
			if !ok {
				return nil
			}
			receivedCounter.WithLabelValues(typeDomainEvent).Inc()
			i.handleMessage(msg)
		case <-i.ctx.Done():
			return nil
		}
	}
}

func (i *domainEventInformer) handleMessage(msg *puller.Message) {
	if id := msg.Attributes["id"]; id == "" {
		msg.Ack()
		handledCounter.WithLabelValues(codes.MissingID.String()).Inc()
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

func (i *domainEventInformer) unmarshalMessage(msg *puller.Message) (*domaineventproto.Event, error) {
	event := &domaineventproto.Event{}
	err := proto.Unmarshal(msg.Data, event)
	if err != nil {
		i.logger.Error("Failed to unmarshal message", zap.Error(err), zap.String("msgID", msg.ID))
		return nil, err
	}
	return event, nil
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
	case domaineventproto.Event_WEBHOOK:
		return notificationproto.Subscription_DOMAIN_EVENT_WEBHOOK, nil
	}
	return notificationproto.Subscription_SourceType(0), ErrUnknownSourceType
}

func (i *domainEventInformer) getEnvironment(
	ctx context.Context,
	environmentNamespace string,
) (*environmentproto.Environment, error) {
	resp, err := i.environmentClient.GetEnvironmentByNamespace(ctx, &environmentproto.GetEnvironmentByNamespaceRequest{
		Namespace: environmentNamespace,
	})
	if err != nil {
		i.logger.Error(
			"Failed to get environment",
			zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
		)
		return nil, err
	}
	return resp.Environment, nil
}
