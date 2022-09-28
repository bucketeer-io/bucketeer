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

package transformer

import (
	"context"
	"errors"
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	userclient "github.com/bucketeer-io/bucketeer/pkg/user/client"
	userdomain "github.com/bucketeer-io/bucketeer/pkg/user/domain"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	clienteventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

const (
	transformTimeout    = 30 * time.Second
	publishMultiTimeout = 5 * time.Minute
)

var (
	errFailedToMarshal    = errors.New("goalBatch: failed to marshal event")
	errFailedToCreateUUID = errors.New("goalBatch: failed to create UUID")
	errFailedToPublish    = errors.New("goalBatch: failed to publish events")
)

type options struct {
	maxMPS     int
	numWorkers int
	metrics    metrics.Registerer
	logger     *zap.Logger
}

var defaultOptions = options{
	maxMPS:     1000,
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

type Transformer interface {
	Check(context.Context) health.Status
	Run() error
	Stop()
}

type transformer struct {
	userClient userclient.Client
	puller     puller.RateLimitedPuller
	publisher  publisher.Publisher
	errgroup   errgroup.Group
	opts       *options
	logger     *zap.Logger
	ctx        context.Context
	cancel     func()
	doneCh     chan struct{}
}

func NewTransformer(
	userClient userclient.Client,
	p puller.Puller,
	publisher publisher.Publisher,
	opts ...Option) Transformer {

	ctx, cancel := context.WithCancel(context.Background())
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.metrics != nil {
		registerMetrics(options.metrics)
	}
	return &transformer{
		userClient: userClient,
		puller:     puller.NewRateLimitedPuller(p, options.maxMPS),
		publisher:  publisher,
		opts:       &options,
		logger:     options.logger.Named("transformer"),
		ctx:        ctx,
		cancel:     cancel,
		doneCh:     make(chan struct{}),
	}
}

func (t *transformer) Run() error {
	defer close(t.doneCh)
	t.errgroup.Go(func() error {
		return t.puller.Run(t.ctx)
	})
	for i := 0; i < t.opts.numWorkers; i++ {
		t.errgroup.Go(t.runWorker)
	}
	return t.errgroup.Wait()
}

func (t *transformer) Stop() {
	t.logger.Info("Stop started")
	t.cancel()
	<-t.doneCh
	t.logger.Info("Stop finished")
}

func (t *transformer) Check(ctx context.Context) health.Status {
	select {
	case <-t.ctx.Done():
		t.logger.Error("Unhealthy due to context Done is closed", zap.Error(t.ctx.Err()))
		return health.Unhealthy
	default:
		if t.errgroup.FinishedCount() > 0 {
			t.logger.Error("Unhealthy", zap.Int32("FinishedCount", t.errgroup.FinishedCount()))
			return health.Unhealthy
		}
		return health.Healthy
	}
}

func (t *transformer) runWorker() error {
	record := func(code codes.Code, startTime time.Time) {
		handledCounter.WithLabelValues(code.String()).Inc()
		handledHistogram.WithLabelValues(code.String()).Observe(time.Since(startTime).Seconds())
	}
	for {
		select {
		case msg, ok := <-t.puller.MessageCh():
			if !ok {
				return nil
			}
			receivedCounter.Inc()
			startTime := time.Now()
			if id := msg.Attributes["id"]; id == "" {
				msg.Ack()
				record(codes.MissingID, startTime)
				continue
			}
			event, environmentNamespace, err := t.unmarshalMessage(msg)
			if err != nil {
				msg.Ack()
				record(codes.BadMessage, startTime)
				continue
			}
			err = t.handle(event, environmentNamespace)
			if err != nil {
				if err == errFailedToMarshal || err == errFailedToCreateUUID {
					record(codes.NonRepeatableError, startTime)
					msg.Ack()
					continue
				}
				record(codes.RepeatableError, startTime)
				msg.Nack()
				continue
			}
			msg.Ack()
			record(codes.OK, startTime)
		case <-t.ctx.Done():
			return nil
		}
	}
}

func (t *transformer) handle(event *clienteventproto.GoalBatchEvent, environmentNamespace string) error {
	events, err := t.transform(event, environmentNamespace)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}
	messages := make([]publisher.Message, 0, len(events))
	for _, event := range events {
		messages = append(messages, event)
	}
	ctx, cancel := context.WithTimeout(context.Background(), publishMultiTimeout)
	defer cancel()
	if errs := t.publisher.PublishMulti(ctx, messages); len(errs) > 0 {
		t.logger.Error("Failed to publish goal events", zap.Any("errors", errs),
			zap.String("environmentNamespace", environmentNamespace))
		eventCounter.WithLabelValues(typeGoal, codeFail).Inc()
		return errFailedToPublish
	}
	eventCounter.WithLabelValues(typeGoal, codeOK).Inc()
	return nil
}

// In case the target user is not found,
// it will return the events empty
func (t *transformer) transform(
	event *clienteventproto.GoalBatchEvent,
	environmentNamespace string,
) ([]*clienteventproto.Event, error) {
	events := make([]*clienteventproto.Event, 0)
	if len(event.UserGoalEventsOverTags) == 0 {
		return events, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), transformTimeout)
	defer cancel()
	user, err := t.getUser(ctx, environmentNamespace, event.UserId)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == grpccodes.NotFound {
			t.logger.Warn("User not found", zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
				zap.String("userId", event.UserId),
			)
			return events, nil
		}
		t.logger.Error("Failed to get user", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("userId", event.UserId),
		)
		return nil, err
	}
	for _, ugeot := range event.UserGoalEventsOverTags {
		tag := ugeot.Tag
		u := t.getUserDataByTag(user, tag)
		for _, uge := range ugeot.UserGoalEvents {
			e, err := t.marshalGoalEvent(environmentNamespace, tag, uge, u)
			if err != nil {
				return nil, err
			}
			events = append(events, e)
		}
	}
	return events, nil
}

func (t *transformer) marshalGoalEvent(
	environmentNamespace, tag string,
	uge *clienteventproto.UserGoalEvent,
	user *userproto.User,
) (*clienteventproto.Event, error) {
	ge := &clienteventproto.GoalEvent{
		SourceId:  clienteventproto.SourceId_GOAL_BATCH,
		Tag:       tag,
		Timestamp: uge.Timestamp,
		GoalId:    uge.GoalId,
		UserId:    user.Id,
		Value:     uge.Value,
		User:      user,
	}
	any, err := ptypes.MarshalAny(ge)
	if err != nil {
		t.logger.Error("Failed to marshal goal event", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("userId", user.Id),
			zap.String("tag", tag),
		)
		return nil, errFailedToMarshal
	}
	id, err := uuid.NewUUID()
	if err != nil {
		t.logger.Error("Failed to create UUID", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("userId", user.Id),
			zap.String("tag", tag),
		)
		return nil, errFailedToCreateUUID
	}
	return &clienteventproto.Event{
		Id:                   id.String(),
		Event:                any,
		EnvironmentNamespace: environmentNamespace,
	}, nil
}

func (t *transformer) unmarshalMessage(msg *puller.Message) (*clienteventproto.GoalBatchEvent, string, error) {
	event := &clienteventproto.Event{}
	if err := proto.Unmarshal(msg.Data, event); err != nil {
		t.logger.Error("Failed to unmarshal message", zap.Error(err), zap.String("msgId", msg.ID))
		return nil, "", err
	}
	goalBatchEvent := &clienteventproto.GoalBatchEvent{}
	if err := ptypes.UnmarshalAny(event.Event, goalBatchEvent); err != nil {
		t.logger.Error("Failed to unmarshal goal event", zap.Error(err), zap.String("msgId", msg.ID))
		return nil, "", err
	}
	return goalBatchEvent, event.EnvironmentNamespace, nil
}

func (t *transformer) getUser(
	ctx context.Context,
	environmentNamespace, userID string,
) (*userproto.User, error) {
	resp, err := t.userClient.GetUser(ctx, &userproto.GetUserRequest{
		UserId:               userID,
		EnvironmentNamespace: environmentNamespace,
	})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

func (t *transformer) getUserDataByTag(user *userproto.User, tag string) *userproto.User {
	u := &userdomain.User{User: user}
	return &userproto.User{
		Id:   user.Id,
		Data: u.Data(tag),
	}
}
