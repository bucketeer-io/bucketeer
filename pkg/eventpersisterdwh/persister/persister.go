// Copyright 2023 The Bucketeer Authors.
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

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	storage "github.com/bucketeer-io/bucketeer/pkg/eventpersisterdwh/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
)

const (
	day = 24 * 60 * 60
)

var (
	ErrUnexpectedMessageType                     = errors.New("eventpersister: unexpected message type")
	ErrAutoOpsRulesNotFound                      = errors.New("eventpersister: auto ops rules not found")
	ErrEvaluationsAreEmpty                       = errors.New("eventpersister: evaluations are empty")
	ErrEvaluationEventIssuedAfterExperimentEnded = errors.New("eventpersister: evaluation event issued after experiment ended") //nolint:lll
	ErrExperimentNotFound                        = errors.New("eventpersister: experiment not found")
	ErrFailedToEvaluateUser                      = errors.New("eventpersister: failed to evaluate user")
	ErrNoAutoOpsRules                            = errors.New("eventpersister: no auto ops rules")
	ErrNothingToLink                             = errors.New("eventpersister: nothing to link")
	ErrInvalidEventTimestamp                     = errors.New("eventpersister: invalid event timestamp")
)

type PersisterDWH struct {
	puller              puller.Puller
	logger              *zap.Logger
	ctx                 context.Context
	mysqlClient         mysql.Client
	runningPullerCtx    context.Context
	runningPullerCancel func()
	isRunning           bool
	rateLimitedPuller   puller.RateLimitedPuller
	cancel              func()
	group               errgroup.Group
	doneCh              chan struct{}
	writer              Writer
	opts                *options
}

type eventMap map[string]proto.Message
type environmentEventMap map[string]eventMap

type options struct {
	maxMPS        int
	numWorkers    int
	flushSize     int
	checkInterval time.Duration
	flushInterval time.Duration
	flushTimeout  time.Duration
	metrics       metrics.Registerer
	logger        *zap.Logger
	batchSize     int
}

type Option func(*options)

func WithCheckInterval(interval time.Duration) Option {
	return func(opts *options) {
		opts.checkInterval = interval
	}
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

func WithFlushTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.flushTimeout = timeout
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

func WithBatchSize(size int) Option {
	return func(opts *options) {
		opts.batchSize = size
	}
}

func NewPersisterDWH(
	p puller.Puller,
	r metrics.Registerer,
	writer Writer,
	mysqlClient mysql.Client,
	opts ...Option,
) *PersisterDWH {
	dopts := &options{
		maxMPS:        1000,
		numWorkers:    1,
		flushSize:     50,
		flushInterval: 5 * time.Second,
		flushTimeout:  20 * time.Second,
		logger:        zap.NewNop(),
		batchSize:     10,
	}
	for _, opt := range opts {
		opt(dopts)
	}
	if dopts.metrics != nil {
		registerMetrics(dopts.metrics)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &PersisterDWH{
		puller:      p,
		mysqlClient: mysqlClient,
		logger:      dopts.logger.Named("persister"),
		ctx:         ctx,
		cancel:      cancel,
		doneCh:      make(chan struct{}),
		writer:      writer,
		opts:        dopts,
	}
}

func (p *PersisterDWH) Run() error {
	defer close(p.doneCh)
	timer := time.NewTimer(p.opts.checkInterval)
	defer timer.Stop()
	subscription := make(chan struct{})
	go p.subscribe(subscription)
	for {
		select {
		case <-timer.C:
			// check if there are running experiment
			exist, err := p.checkRunningExperiments(p.ctx)
			if err != nil {
				p.logger.Error("Failed to check experiments existence", zap.Error(err))
				continue
			}
			if exist {
				if !p.IsRunning() {
					subscription <- struct{}{}
				}
			} else {
				if p.IsRunning() {
					p.unsubscribe()
					err := p.group.Wait()
					if err != nil {
						p.logger.Error("Waiting for puller to finish error", zap.Error(err))
					}
					p.rateLimitedPuller = puller.NewRateLimitedPuller(p.puller, p.opts.maxMPS)
					p.group = errgroup.Group{}
				}
			}
			timer.Reset(p.opts.checkInterval)
		case <-p.ctx.Done():
			if p.IsRunning() {
				p.unsubscribe()
				err := p.group.Wait()
				if err != nil {
					p.logger.Error("Waiting for puller to finish error", zap.Error(err))
					return err
				}
			}
			return nil
		}
	}
}

func (p *PersisterDWH) Stop() {
	p.cancel()
	<-p.doneCh
}

func (p *PersisterDWH) Check(ctx context.Context) health.Status {
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

func (p *PersisterDWH) subscribe(subscription chan struct{}) {
	for {
		select {
		case <-subscription:
			p.isRunning = true
			ctx, cancel := context.WithCancel(context.Background())
			p.runningPullerCtx = ctx
			p.runningPullerCancel = cancel
			p.group.Go(func() error {
				return p.rateLimitedPuller.Run(ctx)
			})
			for i := 0; i < p.opts.numWorkers; i++ {
				p.group.Go(p.batch)
			}
			err := p.group.Wait()
			if err != nil {
				p.logger.Error("Running puller error", zap.Error(err))
			}
			p.isRunning = false
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *PersisterDWH) unsubscribe() {
	p.runningPullerCancel()
}

func (p *PersisterDWH) IsRunning() bool {
	return p.isRunning
}

func (p *PersisterDWH) checkRunningExperiments(ctx context.Context) (bool, error) {
	experimentStorage := storage.NewExperimentStorage(p.mysqlClient)
	count, err := experimentStorage.CountRunningExperiments(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (p *PersisterDWH) batch() error {
	batch := make(map[string]*puller.Message)
	timer := time.NewTimer(p.opts.flushInterval)
	defer timer.Stop()
	for {
		select {
		case msg, ok := <-p.rateLimitedPuller.MessageCh():
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
		case <-p.runningPullerCtx.Done():
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

func (p *PersisterDWH) send(messages map[string]*puller.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), p.opts.flushTimeout)
	defer cancel()
	envEvents := p.extractEvents(messages)
	if len(envEvents) == 0 {
		p.logger.Error("all messages were bad")
		return
	}
	fails := p.writer.Write(ctx, envEvents)
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

func (p *PersisterDWH) extractEvents(messages map[string]*puller.Message) environmentEventMap {
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
