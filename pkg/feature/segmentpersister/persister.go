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

package segmentpersister

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/feature/command"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	domainproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	serviceevent "github.com/bucketeer-io/bucketeer/proto/event/service"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	maxUserIDLength = 100
)

var (
	errSegmentInUse            = errors.New("segment: segment is in use")
	errExceededMaxUserIDLength = fmt.Errorf("segment: max user id length allowed is %d", maxUserIDLength)
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
	domainPublisher   publisher.Publisher
	mysqlClient       mysql.Client
	segmentUsersCache cachev3.SegmentUsersCache
	group             errgroup.Group
	opts              *options
	logger            *zap.Logger
	ctx               context.Context
	cancel            func()
	doneCh            chan struct{}
}

func NewPersister(
	p puller.Puller,
	domainPublisher publisher.Publisher,
	mysqlClient mysql.Client,
	v3Cache cache.MultiGetCache,
	opts ...Option,
) *Persister {
	dopts := &options{
		maxMPS:        100,
		numWorkers:    2,
		flushSize:     2,
		flushInterval: 10 * time.Second,
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
		domainPublisher:   domainPublisher,
		mysqlClient:       mysqlClient,
		segmentUsersCache: cachev3.NewSegmentUsersCache(v3Cache),
		opts:              dopts,
		logger:            dopts.logger.Named("segment-persister"),
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
		p.logger.Error("unhealthy due to context Done is closed", zap.Error(p.ctx.Err()))
		return health.Unhealthy
	default:
		if p.group.FinishedCount() > 0 {
			p.logger.Error("unhealthy", zap.Int32("finishedCount", p.group.FinishedCount()))
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
				p.logger.Warn("message with duplicate id", zap.String("id", id))
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
			return nil
		}
	}
}

func (p *Persister) handleChunk(chunk map[string]*puller.Message) {
	for _, msg := range chunk {
		p.logger.Debug("handling a message", zap.String("msgID", msg.ID))
		event, err := p.unmarshalMessage(msg)
		if err != nil {
			msg.Ack()
			p.logger.Error("failed to unmarshal message", zap.Error(err), zap.String("msgID", msg.ID))
			handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
			continue
		}
		if !validateSegmentUserState(event.State) {
			msg.Ack()
			p.logger.Error(
				"invalid state",
				zap.String("environmentNamespace", event.EnvironmentNamespace),
				zap.Int32("state", int32(event.State)),
			)
			handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
			if err := p.updateSegmentStatus(
				p.ctx,
				event.Editor,
				event.EnvironmentNamespace,
				event.SegmentId,
				0,
				event.State,
				featureproto.Segment_FAILED,
			); err != nil {
				p.logger.Error(
					"failed to update segment status",
					zap.Error(err),
					zap.String("environmentNamespace", event.EnvironmentNamespace),
				)
			}
			continue
		}
		if err := p.handleEvent(p.ctx, event); err != nil {
			switch err {
			case storage.ErrKeyNotFound, v2fs.ErrSegmentNotFound:
				msg.Ack()
				p.logger.Warn("segment not found", zap.Error(err), zap.String("environmentNamespace", event.EnvironmentNamespace))
				handledCounter.WithLabelValues(codes.NonRepeatableError.String()).Inc()
			case errSegmentInUse:
				msg.Ack()
				p.logger.Warn(
					"segment is in use",
					zap.Error(err),
					zap.String("environmentNamespace", event.EnvironmentNamespace),
				)
				handledCounter.WithLabelValues(codes.NonRepeatableError.String()).Inc()
			case errExceededMaxUserIDLength:
				msg.Ack()
				p.logger.Warn(
					"exceeded max user id length",
					zap.Error(err),
					zap.String("environmentNamespace", event.EnvironmentNamespace),
				)
				handledCounter.WithLabelValues(codes.NonRepeatableError.String()).Inc()
				if err := p.updateSegmentStatus(
					p.ctx,
					event.Editor,
					event.EnvironmentNamespace,
					event.SegmentId,
					0,
					event.State,
					featureproto.Segment_FAILED,
				); err != nil {
					p.logger.Error(
						"failed to update segment status",
						zap.Error(err),
						zap.String("environmentNamespace", event.EnvironmentNamespace),
					)
				}
			default:
				// retryable
				msg.Nack()
				p.logger.Error(
					"failed to handle event",
					zap.Error(err),
					zap.String("environmentNamespace", event.EnvironmentNamespace),
				)
				handledCounter.WithLabelValues(codes.RepeatableError.String()).Inc()
			}
			continue
		}
		msg.Ack()
		p.logger.Debug(
			"suceeded to persist segment users",
			zap.String("msgID", msg.ID),
			zap.String("environmentNamespace", event.EnvironmentNamespace),
			zap.String("segmentId", event.SegmentId),
		)
		handledCounter.WithLabelValues(codes.OK.String()).Inc()
	}
}

func (p *Persister) unmarshalMessage(msg *puller.Message) (*serviceevent.BulkSegmentUsersReceivedEvent, error) {
	event := &serviceevent.BulkSegmentUsersReceivedEvent{}
	err := proto.Unmarshal(msg.Data, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func validateSegmentUserState(state featureproto.SegmentUser_State) bool {
	switch state {
	case featureproto.SegmentUser_INCLUDED:
		return true
	default:
		return false
	}
}

func (p *Persister) handleEvent(ctx context.Context, event *serviceevent.BulkSegmentUsersReceivedEvent) error {
	segmentStorage := v2fs.NewSegmentStorage(p.mysqlClient)
	segment, _, err := segmentStorage.GetSegment(ctx, event.SegmentId, event.EnvironmentNamespace)
	if err != nil {
		return err
	}
	if segment.IsInUseStatus {
		return errSegmentInUse
	}
	cnt, err := p.persistSegmentUsers(ctx, event.EnvironmentNamespace, event.SegmentId, event.Data, event.State)
	if err != nil {
		return err
	}
	return p.updateSegmentStatus(
		ctx,
		event.Editor,
		event.EnvironmentNamespace,
		event.SegmentId,
		cnt,
		event.State,
		featureproto.Segment_SUCEEDED,
	)
}

func (p *Persister) persistSegmentUsers(
	ctx context.Context,
	environmentNamespace string,
	segmentID string,
	data []byte,
	state featureproto.SegmentUser_State,
) (int64, error) {
	segmentUserIDs := strings.Split(
		strings.NewReplacer(
			",", "\n",
			"\r\n", "\n",
		).Replace(string(data)),
		"\n",
	)
	uniqueSegmentUserIDs := make(map[string]struct{}, len(segmentUserIDs))
	for _, id := range segmentUserIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if len(id) > maxUserIDLength {
			return 0, errExceededMaxUserIDLength
		}
		uniqueSegmentUserIDs[id] = struct{}{}
	}
	allSegmentUsers := make([]*featureproto.SegmentUser, 0, len(uniqueSegmentUserIDs))
	var cnt int64
	for id := range uniqueSegmentUserIDs {
		cnt++
		user := domain.NewSegmentUser(segmentID, id, state, false)
		allSegmentUsers = append(allSegmentUsers, user.SegmentUser)
	}
	tx, err := p.mysqlClient.BeginTx(ctx)
	if err != nil {
		p.logger.Error("Failed to begin transaction", zap.Error(err))
		return 0, err
	}
	err = p.mysqlClient.RunInTransaction(ctx, tx, func() error {
		segmentUserStorage := v2fs.NewSegmentUserStorage(tx)
		if err := segmentUserStorage.UpsertSegmentUsers(ctx, allSegmentUsers, environmentNamespace); err != nil {
			return err
		}
		return p.updateCache(segmentID, environmentNamespace, allSegmentUsers)
	})
	if err != nil {
		return 0, nil
	}
	return cnt, nil
}

func (p *Persister) updateSegmentStatus(
	ctx context.Context,
	editor *domainproto.Editor,
	environmentNamespace string,
	segmentID string,
	cnt int64,
	state featureproto.SegmentUser_State,
	status featureproto.Segment_Status,
) error {
	tx, err := p.mysqlClient.BeginTx(ctx)
	if err != nil {
		p.logger.Error("Failed to begin transaction", zap.Error(err))
		return err
	}
	return p.mysqlClient.RunInTransaction(ctx, tx, func() error {
		segmentStorage := v2fs.NewSegmentStorage(tx)
		segment, _, err := segmentStorage.GetSegment(ctx, segmentID, environmentNamespace)
		if err != nil {
			return err
		}
		changeCmd := &featureproto.ChangeBulkUploadSegmentUsersStatusCommand{
			Status: status,
			State:  state,
			Count:  cnt,
		}
		handler := command.NewSegmentCommandHandler(editor, segment, p.domainPublisher, environmentNamespace)
		if err := handler.Handle(ctx, changeCmd); err != nil {
			return err
		}
		return segmentStorage.UpdateSegment(ctx, segment, environmentNamespace)
	})
}

func (p *Persister) updateCache(segmentID, environmentNamespace string, users []*featureproto.SegmentUser) error {
	segmentUsers := &featureproto.SegmentUsers{
		SegmentId: segmentID,
		Users:     users,
	}
	if err := p.segmentUsersCache.Put(segmentUsers, environmentNamespace); err != nil {
		p.logger.Error(
			"Failed to cache segment users",
			zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
		)
		return err
	}
	p.logger.Info("Segment users successfully cached",
		zap.String("environmentNamespace", environmentNamespace),
		zap.String("segmentId", segmentID),
		zap.Int("size", len(users)),
	)
	return nil
}
