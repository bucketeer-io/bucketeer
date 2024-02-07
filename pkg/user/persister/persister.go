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
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	ustorage "github.com/bucketeer-io/bucketeer/pkg/user/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	ecproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/service"
)

type options struct {
	maxMPS        int
	numWorkers    int
	flushSize     int
	flushInterval time.Duration
	pubsubTimeout time.Duration
	metrics       metrics.Registerer
	logger        *zap.Logger
}

type Option func(*options)

var defaultOptions = &options{
	maxMPS:        1000,
	numWorkers:    1,
	flushSize:     1000000,
	flushInterval: time.Second,
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
	mysqlClient mysql.Client
	timeNow     func() time.Time
	newUUID     func() (*uuid.UUID, error)
	puller      puller.RateLimitedPuller
	group       errgroup.Group
	opts        *options
	logger      *zap.Logger
	ctx         context.Context
	cancel      func()
	doneCh      chan struct{}
}

func NewPersister(
	mysqlClient mysql.Client,
	p puller.Puller,
	opts ...Option) Persister {

	dopts := defaultOptions
	for _, opt := range opts {
		opt(dopts)
	}
	ctx, cancel := context.WithCancel(context.Background())
	if dopts.metrics != nil {
		registerMetrics(dopts.metrics)
	}
	return &persister{
		mysqlClient: mysqlClient,
		timeNow:     time.Now,
		newUUID:     uuid.NewUUID,
		puller:      puller.NewRateLimitedPuller(p, dopts.maxMPS),
		opts:        dopts,
		logger:      dopts.logger.Named("persister"),
		ctx:         ctx,
		cancel:      cancel,
		doneCh:      make(chan struct{}),
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
	chunk := make(map[string]*puller.Message, p.opts.flushSize)
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
				handledCounter.WithLabelValues(codes.MissingID.String()).Inc()
				continue
			}
			if pre, ok := chunk[id]; ok {
				pre.Ack()
				p.logger.Warn("Message with duplicate id", zap.String("id", id))
				handledCounter.WithLabelValues(codes.DuplicateID.String()).Inc()
			}
			chunk[id] = msg
			if len(chunk) >= p.opts.flushSize {
				p.handleChunk(chunk)
				chunk = make(map[string]*puller.Message, p.opts.flushSize)
				timer.Reset(p.opts.flushInterval)
			}
		case <-timer.C:
			if len(chunk) > 0 {
				p.handleChunk(chunk)
				chunk = make(map[string]*puller.Message, p.opts.flushSize)
			}
			timer.Reset(p.opts.flushInterval)
		case <-p.ctx.Done():
			chunkSize := len(chunk)
			p.logger.Info("Context is done", zap.Int("chunkSize", chunkSize))
			if chunkSize > 0 {
				p.handleChunk(chunk)
				p.logger.Info("All the left messages are processed successfully", zap.Int("chunkSize", chunkSize))
			}
			return nil
		}
	}
}

func (p *persister) handleChunk(chunk map[string]*puller.Message) {
	for _, msg := range chunk {
		event, err := p.unmarshalMessage(msg)
		if err != nil {
			handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
			continue
		}
		if !p.validateEvent(event) {
			handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
			continue
		}
		ok, repeatable := p.upsert(event)
		if !ok {
			if repeatable {
				msg.Nack()
				handledCounter.WithLabelValues(codes.RepeatableError.String()).Inc()
			} else {
				msg.Ack()
				handledCounter.WithLabelValues(codes.NonRepeatableError.String()).Inc()
			}
			continue
		}
		msg.Ack()
		handledCounter.WithLabelValues(codes.OK.String()).Inc()
	}
}

func (p *persister) validateEvent(event *eventproto.UserEvent) bool {
	if event.UserId == "" {
		p.logger.Warn("Message contains an empty User Id", zap.Any("event", event))
		return false
	}
	if event.LastSeen == 0 {
		p.logger.Warn("Message's LastSeen is zero", zap.Any("event", event))
		return false
	}
	return true
}

func (p *persister) unmarshalMessage(msg *puller.Message) (*eventproto.UserEvent, error) {
	event := &ecproto.Event{}
	err := proto.Unmarshal(msg.Data, event)
	if err != nil {
		return nil, err
	}
	var userEvent eventproto.UserEvent
	if err := ptypes.UnmarshalAny(event.Event, &userEvent); err != nil {
		p.logger.Error("Failed to unmarshal Event -> UserEvent", zap.Error(err), zap.Any("msg", msg))
		return nil, err
	}
	return &userEvent, err
}

func (p *persister) upsert(event *eventproto.UserEvent) (ok, repeatable bool) {
	if err := p.upsertMAU(event); err != nil {
		p.logger.Error(
			"Failed to store the mau",
			zap.Error(err),
			zap.String("environmentNamespace", event.EnvironmentNamespace),
			zap.String("userId", event.UserId),
			zap.String("tag", event.Tag),
		)
		return false, true
	}
	return true, false
}

func (p *persister) upsertMAU(event *eventproto.UserEvent) error {
	s := ustorage.NewMysqlMAUStorage(p.mysqlClient)
	return s.UpsertMAU(p.ctx, event, event.EnvironmentNamespace)
}
