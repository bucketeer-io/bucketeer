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
	storage "github.com/bucketeer-io/bucketeer/pkg/eventpersisterops/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
)

var (
	ErrAutoOpsRuleNotFound    = errors.New("eventpersister: auto ops rule not found")
	ErrFeatureEmptyList       = errors.New("eventpersister: list feature returned empty")
	ErrFeatureVersionNotFound = errors.New("eventpersister: feature version not found")
	ErrNoExperiments          = errors.New("eventpersister: no experiments")
	ErrNothingToLink          = errors.New("eventpersister: nothing to link")
	ErrUnexpectedMessageType  = errors.New("eventpersister: unexpected message type")
)

type eventMap map[string]proto.Message
type environmentEventMap map[string]eventMap

type persister struct {
	client                       *pubsub.Client
	topic                        string
	subscription                 string
	pullerNumGoroutines          int
	pullerMaxOutstandingMessages int
	pullerMaxOutstandingBytes    int
	ctx                          context.Context
	cancel                       func()
	updater                      Updater
	mysqlClient                  mysql.Client
	runningPullerCtx             context.Context
	runningPullerCancel          func()
	isRunning                    bool
	rateLimitedPuller            puller.RateLimitedPuller
	group                        errgroup.Group
	doneCh                       chan struct{}
	logger                       *zap.Logger
	opts                         *options
}

type options struct {
	maxMPS        int
	numWorkers    int
	flushSize     int
	checkInterval time.Duration
	flushInterval time.Duration
	flushTimeout  time.Duration
	metrics       metrics.Registerer
	logger        *zap.Logger
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

func NewPersister(
	updater Updater,
	mysqlClient mysql.Client,
	client *pubsub.Client,
	subscription string,
	topic string,
	pullerNumGoroutines int,
	pullerMaxOutstandingMessages int,
	pullerMaxOutstandingBytes int,
	opts ...Option,
) *persister {
	dopts := &options{
		maxMPS:        1000,
		numWorkers:    1,
		flushSize:     100,
		checkInterval: 1 * time.Minute,
		flushInterval: 2 * time.Second,
		flushTimeout:  600 * time.Second,
		logger:        zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	if dopts.metrics != nil {
		registerMetrics(dopts.metrics)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &persister{
		client:                       client,
		topic:                        topic,
		subscription:                 subscription,
		pullerNumGoroutines:          pullerNumGoroutines,
		pullerMaxOutstandingMessages: pullerMaxOutstandingMessages,
		pullerMaxOutstandingBytes:    pullerMaxOutstandingBytes,
		ctx:                          ctx,
		cancel:                       cancel,
		updater:                      updater,
		mysqlClient:                  mysqlClient,
		doneCh:                       make(chan struct{}),
		logger:                       dopts.logger.Named("persister"),
		opts:                         dopts,
	}
}

func (p *persister) Run() error {
	defer close(p.doneCh)
	timer := time.NewTimer(p.opts.checkInterval)
	defer timer.Stop()
	subscription := make(chan struct{})
	go p.subscribe(subscription)
	for {
		select {
		case <-timer.C:
			// check if there are not been triggerd auto ops rules
			exist, err := p.checkAutoOpsRules(p.ctx)
			if err != nil {
				p.logger.Error("Failed to check auto ops rules existence", zap.Error(err))
				continue
			}
			if exist {
				p.logger.Debug("There are untriggered auto ops rules")
				if !p.IsRunning() {
					p.group = errgroup.Group{}
					err := p.createNewPuller()
					if err != nil {
						p.logger.Error("Failed to create new puller", zap.Error(err))
						return err
					}
					subscription <- struct{}{}
					p.logger.Debug("Puller is not running, start pulling messages")
				}
			} else {
				p.logger.Debug("There are no untriggered auto ops rules")
				if p.IsRunning() {
					p.logger.Debug("Puller is running, stop pulling messages")
					p.unsubscribe()
				}
			}
			timer.Reset(p.opts.checkInterval)
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
		if p.group.FailedCount() > 0 {
			p.logger.Error("Unhealthy", zap.Int32("FailedCount", p.group.FailedCount()))
			return health.Unhealthy
		}
		return health.Healthy
	}
}

func (p *persister) createNewPuller() error {
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

func (p *persister) subscribe(subscription chan struct{}) {
	for {
		select {
		case <-subscription:
			p.isRunning = true
			p.logger.Debug("Puller start subscribing")
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
			p.logger.Debug("Puller stopped subscribing")
			p.isRunning = false
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *persister) unsubscribe() {
	p.runningPullerCancel()
	err := p.client.DeleteSubscriptionIfExist(p.subscription)
	if err != nil {
		p.logger.Error("Failed to delete subscription", zap.Error(err))
	}
}

func (p *persister) IsRunning() bool {
	return p.isRunning
}

func (p *persister) checkAutoOpsRules(ctx context.Context) (bool, error) {
	autoOpsRuleStorage := storage.NewAutoOpsRuleStorage(p.mysqlClient)
	autoOpsRuleCount, err := autoOpsRuleStorage.CountNotTriggeredAutoOpsRules(ctx)
	if err != nil {
		return false, err
	}
	return autoOpsRuleCount > 0, nil
}

func (p *persister) batch() error {
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

func (p *persister) send(messages map[string]*puller.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), p.opts.flushTimeout)
	defer cancel()
	envEvents := p.extractEvents(messages)
	if len(envEvents) == 0 {
		p.logger.Error("All messages were bad")
		return
	}
	fails := p.updater.UpdateUserCounts(ctx, envEvents)
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

func (p *persister) extractEvents(messages map[string]*puller.Message) environmentEventMap {
	envEvents := environmentEventMap{}
	handleBadMessage := func(m *puller.Message, err error) {
		m.Ack()
		p.logger.Error("Bad message", zap.Error(err), zap.Any("msg", m))
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
