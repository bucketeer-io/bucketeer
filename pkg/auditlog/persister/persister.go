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

	pb "github.com/golang/protobuf/proto" // nolint:staticcheck
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/auditlog/domain"
	v2als "github.com/bucketeer-io/bucketeer/pkg/auditlog/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	domainevent "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

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

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type Persister struct {
	puller            puller.RateLimitedPuller
	mysqlAdminStorage v2als.AdminAuditLogStorage
	mysqlStorage      v2als.AuditLogStorage
	group             errgroup.Group
	opts              *options
	logger            *zap.Logger
	ctx               context.Context
	cancel            func()
	doneCh            chan struct{}
}

func NewPersister(
	p puller.Puller,
	mysqlClient mysql.Client,
	opts ...Option,
) *Persister {
	dopts := &options{
		maxMPS:        1000,
		numWorkers:    1,
		flushSize:     100,
		flushInterval: time.Second,
		logger:        zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	ctx, cancel := context.WithCancel(context.Background())
	if dopts.metrics != nil {
		registerMetrics(dopts.metrics)
	}
	return &Persister{
		puller:            puller.NewRateLimitedPuller(p, dopts.maxMPS),
		mysqlAdminStorage: v2als.NewAdminAuditLogStorage(mysqlClient),
		mysqlStorage:      v2als.NewAuditLogStorage(mysqlClient),
		opts:              dopts,
		logger:            dopts.logger.Named("persister"),
		ctx:               ctx,
		cancel:            cancel,
		doneCh:            make(chan struct{}),
	}
}

func (p *Persister) Run() error {
	defer close(p.doneCh)
	p.group.Go(func() error {
		return p.puller.Run(p.ctx)
	})
	for i := 0; i < p.opts.numWorkers; i++ {
		p.group.Go(p.runWorker)
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

func (p *Persister) runWorker() error {
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
			if _, ok := chunk[id]; ok {
				p.logger.Warn("Message with duplicate id", zap.String("id", id))
				handledCounter.WithLabelValues(codes.DuplicateID.String()).Inc()
			}
			chunk[id] = msg
			if len(chunk) >= p.opts.flushSize {
				p.flushChunk(chunk)
				chunk = make(map[string]*puller.Message, p.opts.flushSize)
				timer.Reset(p.opts.flushInterval)
			}
		case <-timer.C:
			if len(chunk) > 0 {
				p.flushChunk(chunk)
				chunk = make(map[string]*puller.Message, p.opts.flushSize)
			}
			timer.Reset(p.opts.flushInterval)
		case <-p.ctx.Done():
			return nil
		}
	}
}

func (p *Persister) flushChunk(chunk map[string]*puller.Message) {
	auditlogs, adminAuditLogs, messages, adminMessages := p.extractAuditLogs(chunk)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Environment audit logs
	p.createAuditLogsMySQL(ctx, auditlogs, messages, p.mysqlStorage.CreateAuditLogs)
	// Admin audit logs
	p.createAuditLogsMySQL(ctx, adminAuditLogs, adminMessages, p.mysqlAdminStorage.CreateAdminAuditLogs)
}

func (p *Persister) extractAuditLogs(
	chunk map[string]*puller.Message,
) (auditlogs, adminAuditLogs []*domain.AuditLog, messages, adminMessages []*puller.Message) {
	for _, msg := range chunk {
		event := &domainevent.Event{}
		if err := pb.Unmarshal(msg.Data, event); err != nil {
			p.logger.Error("Failed to unmarshal message", zap.Error(err))
			msg.Ack()
			continue
		}
		if event.IsAdminEvent {
			adminAuditLogs = append(adminAuditLogs, domain.NewAuditLog(event, storage.AdminEnvironmentNamespace))
			adminMessages = append(adminMessages, msg)
		} else {
			auditlogs = append(auditlogs, domain.NewAuditLog(event, event.EnvironmentNamespace))
			messages = append(messages, msg)
		}
	}
	return
}

func (p *Persister) createAuditLogsMySQL(
	ctx context.Context,
	auditlogs []*domain.AuditLog,
	messages []*puller.Message,
	createFunc func(ctx context.Context, auditLogs []*domain.AuditLog) error,
) {
	if len(auditlogs) == 0 {
		return
	}
	if err := createFunc(ctx, auditlogs); err != nil {
		p.logger.Error("Failed to put admin audit logs", zap.Error(err))
		for _, msg := range messages {
			handledCounter.WithLabelValues(codes.RepeatableError.String()).Inc()
			msg.Nack()
		}
		return
	}
	for _, msg := range messages {
		handledCounter.WithLabelValues(codes.OK.String()).Inc()
		msg.Ack()
	}
}
