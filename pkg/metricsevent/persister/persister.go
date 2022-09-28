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

package persister

import (
	"context"
	"errors"
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/metricsevent/storage"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
)

var (
	ErrUnknownEvent    = errors.New("metricsevent persister: unknown metrics event")
	ErrInvalidDuration = errors.New("metricsevent persister: invalid duration")

	getEvaluationLatencyMetricsEvent = &eventproto.GetEvaluationLatencyMetricsEvent{}
	getEvaluationSizeMetricsEvent    = &eventproto.GetEvaluationSizeMetricsEvent{}
	timeoutErrorCountMetricsEvent    = &eventproto.TimeoutErrorCountMetricsEvent{}
	internalErrorCountMetricsEvent   = &eventproto.InternalErrorCountMetricsEvent{}
)

type options struct {
	maxMPS        int
	numWorkers    int
	pubsubTimeout time.Duration
	metrics       metrics.Registerer
	logger        *zap.Logger
}

type Option func(*options)

var defaultOptions = &options{
	maxMPS:        1000,
	numWorkers:    1,
	pubsubTimeout: 20 * time.Second,
	logger:        zap.NewNop(),
}

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

type Persister interface {
	Run() error
	Stop()
	Check(context.Context) health.Status
}

type persister struct {
	puller  puller.RateLimitedPuller
	storage storage.Storage
	group   errgroup.Group
	opts    *options
	logger  *zap.Logger
	ctx     context.Context
	cancel  func()
	doneCh  chan struct{}
}

func NewPersister(p puller.Puller, opts ...Option) Persister {
	dopts := defaultOptions
	for _, opt := range opts {
		opt(dopts)
	}
	ctx, cancel := context.WithCancel(context.Background())
	if dopts.metrics != nil {
		registerMetrics(dopts.metrics)
	}
	logger := dopts.logger.Named("persister")
	return &persister{
		puller:  puller.NewRateLimitedPuller(p, dopts.maxMPS),
		storage: storage.NewStorage(logger, dopts.metrics),
		opts:    dopts,
		logger:  logger,
		ctx:     ctx,
		cancel:  cancel,
		doneCh:  make(chan struct{}),
	}
}

func (p *persister) Run() error {
	defer close(p.doneCh)
	p.group.Go(func() error {
		return p.puller.Run(p.ctx)
	})
	for i := 0; i < p.opts.numWorkers; i++ {
		p.group.Go(p.runWorker)
	}
	return p.group.Wait()
}

func (p *persister) Stop() {
	p.cancel()
	<-p.doneCh
}

func (p *persister) Check(ctx context.Context) health.Status {
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

func (p *persister) runWorker() error {
	record := func(code codes.Code, startTime time.Time) {
		handledCounter.WithLabelValues(code.String()).Inc()
		handledHistogram.WithLabelValues(code.String()).Observe(time.Since(startTime).Seconds())
	}
	for {
		select {
		case msg, ok := <-p.puller.MessageCh():
			if !ok {
				return nil
			}
			receivedCounter.Inc()
			startTime := time.Now()
			if id := msg.Attributes["id"]; id == "" {
				p.logger.Error("message has no id")
				msg.Ack()
				record(codes.MissingID, startTime)
				continue
			}
			err := p.handle(msg)
			if err != nil {
				msg.Ack()
				record(codes.NonRepeatableError, startTime)
				continue
			}
			msg.Ack()
			record(codes.OK, startTime)
		case <-p.ctx.Done():
			return nil
		}
	}
}

func (p *persister) handle(message *puller.Message) error {
	metricsEvents, err := p.unmarshalMessage(message)
	if err != nil {
		p.logger.Error("message is bad")
		return err
	}
	err = p.saveMetrics(metricsEvents)
	if err != nil {
		p.logger.Error("could not store data to prometheus client", zap.Error(err))
		return err
	}
	return nil
}

func (p *persister) unmarshalMessage(message *puller.Message) (*eventproto.MetricsEvent, error) {
	event := &eventproto.Event{}
	if err := proto.Unmarshal(message.Data, event); err != nil {
		p.logger.Error("ummarshal event failed",
			zap.Error(err),
			zap.Any("msg", message),
		)
		handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
		return nil, err
	}
	me := &eventproto.MetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, me); err != nil {
		p.logger.Error("ummarshal metrics event failed",
			zap.Error(err),
			zap.Any("msg", message),
		)
		handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
		return nil, err
	}
	return me, nil
}

func (p *persister) saveMetrics(event *eventproto.MetricsEvent) error {
	if ptypes.Is(event.Event, getEvaluationLatencyMetricsEvent) {
		return p.saveGetEvaluationLatencyMetricsEvent(event)
	}
	if ptypes.Is(event.Event, getEvaluationSizeMetricsEvent) {
		return p.saveGetEvaluationSizeMetricsEvent(event)
	}
	if ptypes.Is(event.Event, timeoutErrorCountMetricsEvent) {
		return p.saveTimeoutErrorCountMetricsEvent(event)
	}
	if ptypes.Is(event.Event, internalErrorCountMetricsEvent) {
		return p.saveInternalErrorCountMetricsEvent(event)
	}
	return ErrUnknownEvent
}

func (p *persister) saveGetEvaluationLatencyMetricsEvent(event *eventproto.MetricsEvent) error {
	ev := &eventproto.GetEvaluationLatencyMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	if ev.Duration == nil {
		return ErrInvalidDuration
	}
	var tag, status string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
		status = ev.Labels["state"]
	}
	dur, err := ptypes.Duration(ev.Duration)
	if err != nil {
		return ErrInvalidDuration
	}
	p.storage.SaveGetEvaluationLatencyMetricsEvent(tag, status, dur)
	return nil
}

func (p *persister) saveGetEvaluationSizeMetricsEvent(event *eventproto.MetricsEvent) error {
	ev := &eventproto.GetEvaluationSizeMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	var tag, status string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
		status = ev.Labels["state"]
	}
	p.storage.SaveGetEvaluationSizeMetricsEvent(tag, status, ev.SizeByte)
	return nil
}

func (p *persister) saveTimeoutErrorCountMetricsEvent(event *eventproto.MetricsEvent) error {
	ev := &eventproto.TimeoutErrorCountMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	p.storage.SaveTimeoutErrorCountMetricsEvent(ev.Tag)
	return nil
}

func (p *persister) saveInternalErrorCountMetricsEvent(event *eventproto.MetricsEvent) error {
	ev := &eventproto.InternalErrorCountMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	p.storage.SaveInternalErrorCountMetricsEvent(ev.Tag)
	return nil
}
