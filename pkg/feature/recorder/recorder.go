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

package recorder

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	ftstorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"

	"github.com/bucketeer-io/bucketeer/proto/event/client"
)

const (
	userDataAppVersion = "app_version"
)

type options struct {
	maxMPS          int
	numWorkers      int
	flushInterval   time.Duration
	startupInterval time.Duration
	metrics         metrics.Registerer
	logger          *zap.Logger
}

type Option func(*options)

func WithMaxMPS(mps int) Option {
	return func(opts *options) {
		opts.maxMPS = mps
	}
}

func WithFlushInterval(interval time.Duration) Option {
	return func(opts *options) {
		opts.flushInterval = interval
	}
}

func WithStartupInterval(interval time.Duration) Option {
	return func(opts *options) {
		opts.startupInterval = interval
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

type Recorder interface {
	Check(context.Context) health.Status
	Run() error
	Stop()
}

type recorder struct {
	puller        puller.RateLimitedPuller
	storageClient mysql.Client
	opts          *options
	group         errgroup.Group
	logger        *zap.Logger
	ctx           context.Context
	cancel        func()
	doneCh        chan struct{}
}

type lastUsedInfoCache map[string]*domain.FeatureLastUsedInfo
type environmentLastUsedInfoCache map[string]lastUsedInfoCache

func NewRecorder(p puller.Puller, sc mysql.Client, opts ...Option) Recorder {
	ctx, cancel := context.WithCancel(context.Background())
	dopts := &options{
		maxMPS:          1000,
		numWorkers:      1,
		logger:          zap.NewNop(),
		flushInterval:   time.Minute,
		startupInterval: time.Second,
	}
	for _, opt := range opts {
		opt(dopts)
	}
	if dopts.metrics != nil {
		registerMetrics(dopts.metrics)
	}
	return &recorder{
		puller:        puller.NewRateLimitedPuller(p, dopts.maxMPS),
		storageClient: sc,
		opts:          dopts,
		logger:        dopts.logger.Named("recorder"),
		ctx:           ctx,
		cancel:        cancel,
		doneCh:        make(chan struct{}),
	}
}

// Run starts workers.
// To distribute requests to DB, sleep for a second when starting each worker.
func (r *recorder) Run() error {
	defer close(r.doneCh)
	r.group.Go(func() error {
		return r.puller.Run(r.ctx)
	})
	for i := 0; i < r.opts.numWorkers; i++ {
		r.group.Go(r.runWorker)
		time.Sleep(r.opts.startupInterval)
	}
	return r.group.Wait()
}

func (r *recorder) Stop() {
	r.cancel()
	<-r.doneCh
}

func (r *recorder) Check(ctx context.Context) health.Status {
	select {
	case <-r.ctx.Done():
		r.logger.Error("Unhealthy due to context Done is closed", zap.Error(r.ctx.Err()))
		return health.Unhealthy
	default:
		if r.group.FinishedCount() > 0 {
			r.logger.Error("Unhealthy", zap.Int32("FinishedCount", r.group.FinishedCount()))
			return health.Unhealthy
		}
		return health.Healthy
	}
}

func (r *recorder) runWorker() error {
	timer := time.NewTimer(r.opts.flushInterval)
	defer timer.Stop()
	envCache := environmentLastUsedInfoCache{}
	defer r.writeEnvLastUsedInfo(envCache)
	for {
		select {
		case <-r.ctx.Done():
			return nil
		case msg, ok := <-r.puller.MessageCh():
			if !ok {
				return nil
			}
			receivedCounter.Inc()
			event, err := r.unmarshalMessage(msg)
			if err != nil {
				msg.Ack()
				handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
				continue
			}
			evaluationEvent, err := r.unmarshalEvent(event.Event)
			if err != nil {
				msg.Ack()
				handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
				continue
			}
			r.cacheEnvLastUsedInfo(evaluationEvent, envCache, event.EnvironmentNamespace)
			msg.Ack()
			handledCounter.WithLabelValues(codes.OK.String()).Inc()
		case <-timer.C:
			r.writeEnvLastUsedInfo(envCache)
			envCache = make(environmentLastUsedInfoCache, len(envCache))
			timer.Reset(r.opts.flushInterval)
		}
	}
}

func (r *recorder) cacheEnvLastUsedInfo(
	e *client.EvaluationEvent,
	envCache environmentLastUsedInfoCache,
	environmentNamespace string,
) {
	// FIXME: Until the Web SDK is released including the fix below,
	// We need to ignore the error, otherwise the Feature Flag Status used won't be updated
	// https://github.com/bucketeer-io/bucketeer/issues/1145
	var clientVersion string
	if e.User == nil {
		r.logger.Warn("Failed to cache last used info. User is nil.",
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("featureId", e.FeatureId),
			zap.Int32("featureVersion", e.FeatureVersion))
	} else {
		clientVersion = e.User.Data[userDataAppVersion]
	}
	id := domain.FeatureLastUsedInfoID(e.FeatureId, e.FeatureVersion)
	if cache, ok := envCache[environmentNamespace]; ok {
		if info, ok := cache[id]; ok {
			info.UsedAt(e.Timestamp)
			if err := info.SetClientVersion(clientVersion); err != nil {
				r.logger.Error("Failed to set client version",
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
					zap.String("featureId", info.FeatureId),
					zap.Int32("featureVersion", info.Version),
					zap.String("clientVersion", clientVersion))
			}
			return
		}
		cache[id] = domain.NewFeatureLastUsedInfo(e.FeatureId, e.FeatureVersion, e.Timestamp, clientVersion)
		return
	}
	cache := lastUsedInfoCache{}
	cache[id] = domain.NewFeatureLastUsedInfo(e.FeatureId, e.FeatureVersion, e.Timestamp, clientVersion)
	envCache[environmentNamespace] = cache
}

func (r *recorder) writeEnvLastUsedInfo(envCache environmentLastUsedInfoCache) {
	for environmentNamespace, cache := range envCache {
		info := make([]*domain.FeatureLastUsedInfo, 0, len(cache))
		for _, v := range cache {
			info = append(info, v)
		}
		if err := r.upsertMultiFeatureLastUsedInfo(context.Background(), info, environmentNamespace); err != nil {
			r.logger.Error("failed to write featureLastUsedInfo", zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace))
			continue
		}
	}
}

func (r *recorder) unmarshalMessage(msg *puller.Message) (*client.Event, error) {
	var event client.Event
	if err := proto.Unmarshal(msg.Data, &event); err != nil {
		r.logger.Error("bad message", zap.Error(err), zap.Any("msg", msg))
		return nil, err
	}
	return &event, nil
}

func (r *recorder) unmarshalEvent(event *any.Any) (*client.EvaluationEvent, error) {
	var evaluationEvent client.EvaluationEvent
	if err := ptypes.UnmarshalAny(event, &evaluationEvent); err != nil {
		r.logger.Error("unexpected event", zap.Error(err), zap.Any("event", event))
		return nil, err
	}
	return &evaluationEvent, nil
}

func (r *recorder) upsertMultiFeatureLastUsedInfo(
	ctx context.Context,
	featureLastUsedInfos []*domain.FeatureLastUsedInfo,
	environmentNamespace string,
) error {
	ids := make([]string, 0, len(featureLastUsedInfos))
	for _, f := range featureLastUsedInfos {
		ids = append(ids, f.ID())
	}
	tx, err := r.storageClient.BeginTx(ctx)
	if err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return err
	}
	err = r.storageClient.RunInTransaction(ctx, tx, func() error {
		storage := ftstorage.NewFeatureLastUsedInfoStorage(r.storageClient)
		updatedInfo := make([]*domain.FeatureLastUsedInfo, 0, len(ids))
		currentInfo, err := storage.GetFeatureLastUsedInfos(ctx, ids, environmentNamespace)
		if err != nil {
			return err
		}
		currentInfoMap := make(map[string]*domain.FeatureLastUsedInfo, len(currentInfo))
		for _, c := range currentInfo {
			currentInfoMap[c.ID()] = c
		}
		for _, f := range featureLastUsedInfos {
			v, ok := currentInfoMap[f.ID()]
			if !ok {
				updatedInfo = append(updatedInfo, f)
				continue
			}
			var update bool
			if v.LastUsedAt < f.LastUsedAt {
				update = true
				v.LastUsedAt = f.LastUsedAt
			}
			if v.ClientOldestVersion != f.ClientOldestVersion {
				update = true
				v.ClientOldestVersion = f.ClientOldestVersion
			}
			if v.ClientLatestVersion != f.ClientLatestVersion {
				update = true
				v.ClientLatestVersion = f.ClientLatestVersion
			}
			if update {
				updatedInfo = append(updatedInfo, v)
			}
		}
		for _, info := range updatedInfo {
			if err := storage.UpsertFeatureLastUsedInfo(ctx, info, environmentNamespace); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
