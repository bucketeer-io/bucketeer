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
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	storage "github.com/bucketeer-io/bucketeer/pkg/eventpersisterdwh/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
)

const (
	day = 24 * 60 * 60

	pubsubErrNotFound = "NotFound"
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
	client                       *pubsub.Client
	topic                        string
	subscription                 string
	pullerNumGoroutines          int
	pullerMaxOutstandingMessages int
	pullerMaxOutstandingBytes    int
	logger                       *zap.Logger
	ctx                          context.Context
	mysqlClient                  mysql.Client
	runningPullerCtx             context.Context
	runningPullerCancel          func()
	isRunning                    bool
	rateLimitedPuller            puller.RateLimitedPuller
	cancel                       func()
	group                        errgroup.Group
	doneCh                       chan struct{}
	writer                       Writer
	opts                         *options
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
	client *pubsub.Client,
	r metrics.Registerer,
	writer Writer,
	mysqlClient mysql.Client,
	subscription string,
	topic string,
	pullerNumGoroutines int,
	pullerMaxOutstandingMessages int,
	pullerMaxOutstandingBytes int,
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
		client:                       client,
		topic:                        topic,
		subscription:                 subscription,
		pullerNumGoroutines:          pullerNumGoroutines,
		pullerMaxOutstandingMessages: pullerMaxOutstandingMessages,
		pullerMaxOutstandingBytes:    pullerMaxOutstandingBytes,
		mysqlClient:                  mysqlClient,
		logger:                       dopts.logger.Named("persister"),
		ctx:                          ctx,
		cancel:                       cancel,
		doneCh:                       make(chan struct{}),
		writer:                       writer,
		opts:                         dopts,
	}
}

func (p *PersisterDWH) Run() error {
	defer close(p.doneCh)
	ticker := time.NewTicker(p.opts.checkInterval)
	defer ticker.Stop()
	subscription := make(chan struct{})
	go p.subscribe(subscription)
	for {
		select {
		case <-ticker.C:
			// check if there are running experiments
			exist, err := p.checkRunningExperiments(p.ctx)
			if err != nil {
				p.logger.Error("Failed to check experiments existence", zap.Error(err))
				continue
			}
			if exist {
				p.logger.Debug("There are running experiments")
				if !p.IsRunning() {
					err = p.createNewPuller()
					if err != nil {
						p.logger.Error("Failed to create new puller", zap.Error(err))
						return err
					}
					p.group = errgroup.Group{}
					subscription <- struct{}{}
					p.logger.Debug("Puller is not running, start pulling messages")
				}
			} else {
				p.logger.Debug("There are no running experiments")
				if p.IsRunning() {
					p.logger.Debug("Puller is running, stop pulling messages")
					p.unsubscribe()
				}
				// delete subscription if it exists
				exists, err := p.client.SubscriptionExists(p.subscription)
				if err != nil {
					p.logger.Error("Failed to check subscription existence", zap.Error(err))
					continue
				}
				if exists {
					p.logger.Debug("Subscription exists, delete it now",
						zap.String("subscription", p.subscription),
					)
					err = p.client.DeleteSubscription(p.subscription)
					if err != nil {
						p.logger.Error("Failed to delete subscription", zap.Error(err))
						continue
					} else {
						p.logger.Debug("Subscription deleted successfully",
							zap.String("subscription", p.subscription),
						)
					}
				} else {
					p.logger.Debug("Subscription does not exist", zap.String("subscription", p.subscription))
				}
			}
		case <-p.ctx.Done():
			p.logger.Debug("Context is done")
			if p.IsRunning() {
				p.logger.Debug("Puller is running, stop pulling messages")
				p.unsubscribe()
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
		if p.group.FailedCount() > 0 {
			p.logger.Error("Unhealthy", zap.Int32("FailedCount", p.group.FailedCount()))
			return health.Unhealthy
		}
		return health.Healthy
	}
}

func (p *PersisterDWH) createNewPuller() error {
	pubsubPuller, err := p.client.CreatePuller(p.subscription, p.topic,
		pubsub.WithNumGoroutines(p.pullerNumGoroutines),
		pubsub.WithMaxOutstandingMessages(p.pullerMaxOutstandingMessages),
		pubsub.WithMaxOutstandingBytes(p.pullerMaxOutstandingBytes),
	)
	if err != nil {
		return err
	}
	p.rateLimitedPuller = puller.NewRateLimitedPuller(pubsubPuller, p.opts.maxMPS)
	return nil
}

func (p *PersisterDWH) subscribe(subscription chan struct{}) {
	for {
		select {
		case <-subscription:
			p.isRunning = true
			p.logger.Debug("Puller started subscribing")
			ctx, cancel := context.WithCancel(context.Background())
			p.runningPullerCtx = ctx
			p.runningPullerCancel = cancel
			p.group.Go(func() error {
				err := p.rateLimitedPuller.Run(ctx)
				if err != nil {
					if strings.Contains(err.Error(), pubsubErrNotFound) {
						p.logger.Debug("Failed to pull messages. Subscription does not exist",
							zap.String("subscription", p.subscription))
						p.unsubscribe()
						return nil
					}
					p.logger.Error("Failed to pull messages", zap.Error(err))
					return err
				}
				return nil
			})
			for i := 0; i < p.opts.numWorkers; i++ {
				p.group.Go(p.batch)
			}
			err := p.group.Wait()
			if err != nil {
				p.logger.Error("Failed while running pull messages", zap.Error(err))
			}
			p.logger.Debug("Puller stopped subscribing")
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
	ticker := time.NewTicker(p.opts.flushInterval)
	defer ticker.Stop()
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
		case <-ticker.C:
			if len(batch) > 0 {
				p.send(batch)
				batch = make(map[string]*puller.Message)
			}
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
