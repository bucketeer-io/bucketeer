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

package persister

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/cache"
	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	ErrUnexpectedMessageType = errors.New("eventpersister: unexpected message type")
	ErrAutoOpsRulesNotFound  = errors.New("eventpersister: auto ops rules not found")
	ErrExperimentNotFound    = errors.New("eventpersister: experiment not found")
	ErrNoAutoOpsRules        = errors.New("eventpersister: no auto ops rules")
	ErrNoExperiments         = errors.New("eventpersister: no experiments")
	ErrNothingToLink         = errors.New("eventpersister: nothing to link")
	ErrReasonNil             = errors.New("eventpersister: reason is nil")
)

const (
	eventCountKey      = "ec"
	userCountKey       = "uc"
	defaultVariationID = "default"
)

type eventMap map[string]proto.Message
type environmentEventMap map[string]eventMap

type options struct {
	maxMPS        int
	numWorkers    int
	flushSize     int
	flushInterval time.Duration
	metrics       metrics.Registerer
	logger        *zap.Logger
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

func WithFlushSize(s int) Option {
	return func(opts *options) {
		opts.flushSize = s
	}
}

func WithFlushInterval(i time.Duration) Option {
	return func(opts *options) {
		opts.flushInterval = i
	}
}

func WithMetrics(r metrics.Registerer) Option {
	return func(opts *options) {
		opts.metrics = r
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type Persister struct {
	puller                puller.RateLimitedPuller
	group                 errgroup.Group
	opts                  *options
	logger                *zap.Logger
	ctx                   context.Context
	cancel                func()
	doneCh                chan struct{}
	evaluationCountCacher cache.MultiGetDeleteCountCache
}

func NewPersister(
	p puller.Puller,
	v3Cache cache.MultiGetDeleteCountCache,
	opts ...Option,
) *Persister {
	dopts := &options{
		maxMPS:        1000,
		numWorkers:    1,
		flushSize:     50,
		flushInterval: 5 * time.Second,
		logger:        zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	if dopts.metrics != nil {
		registerMetrics(dopts.metrics)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Persister{
		puller:                puller.NewRateLimitedPuller(p, dopts.maxMPS),
		opts:                  dopts,
		logger:                dopts.logger.Named("persister"),
		ctx:                   ctx,
		cancel:                cancel,
		doneCh:                make(chan struct{}),
		evaluationCountCacher: v3Cache,
	}
}

func (p *Persister) Run() error {
	defer close(p.doneCh)
	p.group.Go(func() error {
		return p.puller.Run(p.ctx)
	})
	for i := 0; i < p.opts.numWorkers; i++ {
		p.group.Go(p.batch)
	}
	return p.group.Wait()
}

func (p *Persister) Stop() {
	p.cancel()
	<-p.doneCh
}

func (p *Persister) Check(ctx context.Context) health.Status {
	select {
	case <-p.ctx.Done():
		p.logger.Error("Unhealthy due to context Done is closed", zap.Error(p.ctx.Err()))
		return health.Unhealthy
	default:
		if p.group.FinishedCount() > 0 {
			p.logger.Error("Unhealthy", zap.Int32("FinishedCount", p.group.FinishedCount()))
			return health.Unhealthy
		}
		return health.Healthy
	}
}

func (p *Persister) batch() error {
	batch := make(map[string]*puller.Message)
	timer := time.NewTimer(p.opts.flushInterval)
	defer timer.Stop()
	for {
		select {
		case msg, ok := <-p.puller.MessageCh():
			if !ok {
				return nil
			}
			receivedCounter.Inc()
			id := msg.Attributes["id"]
			if id == "" {
				msg.Ack()
				// TODO: better log format for msg data
				handledCounter.WithLabelValues(codes.MissingID.String()).Inc()
				continue
			}
			if previous, ok := batch[id]; ok {
				previous.Ack()
				p.logger.Warn("Message with duplicate id", zap.String("id", id))
				handledCounter.WithLabelValues(codes.DuplicateID.String()).Inc()
			}
			batch[id] = msg
			if len(batch) < p.opts.flushSize {
				continue
			}
			p.send(batch)
			batch = make(map[string]*puller.Message)
			timer.Reset(p.opts.flushInterval)
		case <-timer.C:
			if len(batch) > 0 {
				p.send(batch)
				batch = make(map[string]*puller.Message)
			}
			timer.Reset(p.opts.flushInterval)
		case <-p.ctx.Done():
			batchSize := len(batch)
			p.logger.Info("Context is done", zap.Int("batchSize", batchSize))
			if len(batch) > 0 {
				p.send(batch)
				p.logger.Info("All the left messages are processed successfully", zap.Int("batchSize", batchSize))
			}
			return nil
		}
	}
}

func (p *Persister) send(messages map[string]*puller.Message) {
	envEvents := p.extractEvents(messages)
	if len(envEvents) == 0 {
		p.logger.Error("all messages were bad")
		return
	}
	fails := make(map[string]bool, len(messages))
	for environmentNamespace, events := range envEvents {
		for id, event := range events {
			if err := p.upsertEvaluationCount(event, environmentNamespace); err != nil {
				p.logger.Error(
					"Failed to upsert an evaluation event in redis",
					zap.Error(err),
					zap.String("id", id),
					zap.String("environmentNamespace", environmentNamespace),
				)
				fails[id] = true
			}
		}
	}
	for id, m := range messages {
		if repeatable, ok := fails[id]; ok {
			if repeatable {
				m.Nack()
				handledCounter.WithLabelValues(codes.RepeatableError.String()).Inc()
			} else {
				m.Ack()
				handledCounter.WithLabelValues(codes.NonRepeatableError.String()).Inc()
			}
			continue
		}
		m.Ack()
		handledCounter.WithLabelValues(codes.OK.String()).Inc()
	}
}

func (p *Persister) extractEvents(messages map[string]*puller.Message) environmentEventMap {
	envEvents := environmentEventMap{}
	handleBadMessage := func(m *puller.Message, err error) {
		m.Ack()
		p.logger.Error("bad message", zap.Error(err), zap.Any("msg", m))
		handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
	}
	for _, m := range messages {
		event := &eventproto.Event{}
		if err := proto.Unmarshal(m.Data, event); err != nil {
			handleBadMessage(m, err)
			continue
		}
		var innerEvent ptypes.DynamicAny
		if err := ptypes.UnmarshalAny(event.Event, &innerEvent); err != nil {
			handleBadMessage(m, err)
			continue
		}
		if innerEvents, ok := envEvents[event.EnvironmentNamespace]; ok {
			innerEvents[event.Id] = innerEvent.Message
			continue
		}
		envEvents[event.EnvironmentNamespace] = eventMap{event.Id: innerEvent.Message}
	}
	return envEvents
}

func getVariationID(reason *featureproto.Reason, vID string) (string, error) {
	if reason == nil {
		return "", ErrReasonNil
	}
	if reason.Type == featureproto.Reason_CLIENT {
		return defaultVariationID, nil
	}
	return vID, nil
}

func (p *Persister) upsertEvaluationCount(event proto.Message, environmentNamespace string) error {
	if e, ok := event.(*eventproto.EvaluationEvent); ok {
		vID, err := getVariationID(e.Reason, e.VariationId)
		if err != nil {
			return err
		}
		// To avoid duplication when the request fails, we increment the event count in the end
		// because the user count is an unique count, and there is no problem adding the same event more than once
		uckv2 := p.newEvaluationCountkeyV2(userCountKey, e.FeatureId, vID, environmentNamespace, e.Timestamp)
		if err := p.countUser(uckv2, e.UserId); err != nil {
			return err
		}
		eckv2 := p.newEvaluationCountkeyV2(eventCountKey, e.FeatureId, vID, environmentNamespace, e.Timestamp)
		if err := p.countEvent(eckv2); err != nil {
			return err
		}
	}
	return nil
}

func (p *Persister) newEvaluationCountkeyV2(
	kind, featureID, variationID, environmentNamespace string,
	timestamp int64,
) string {
	t := time.Unix(timestamp, 0)
	date := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.UTC)
	return cache.MakeKey(
		kind,
		fmt.Sprintf("%d:%s:%s", date.Unix(), featureID, variationID),
		environmentNamespace,
	)
}

func (p *Persister) countEvent(key string) error {
	_, err := p.evaluationCountCacher.Increment(key)
	if err != nil {
		return err
	}
	return nil
}

func (p *Persister) countUser(key, userID string) error {
	_, err := p.evaluationCountCacher.PFAdd(key, userID)
	if err != nil {
		return err
	}
	return nil
}
