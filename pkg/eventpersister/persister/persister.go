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
	"sync"
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/cache"
	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	ftdomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	ftstorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type lastUsedInfoCache map[string]*ftdomain.FeatureLastUsedInfo
type environmentLastUsedInfoCache map[string]lastUsedInfoCache

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
	userDataAppVersion = "app_version"
)

type eventMap map[string]*eventproto.EvaluationEvent
type environmentEventMap map[string]eventMap

type options struct {
	maxMPS             int
	numWorkers         int
	flushSize          int
	flushInterval      time.Duration
	writeCacheInterval time.Duration
	metrics            metrics.Registerer
	logger             *zap.Logger
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

func WithWriteCacheInterval(i time.Duration) Option {
	return func(opts *options) {
		opts.writeCacheInterval = i
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
	mysqlClient           mysql.Client
	group                 errgroup.Group
	opts                  *options
	logger                *zap.Logger
	ctx                   context.Context
	cancel                func()
	doneCh                chan struct{}
	envLastUsedCache      environmentLastUsedInfoCache
	evaluationCountCacher cache.MultiGetDeleteCountCache
	mutex                 sync.Mutex
}

func NewPersister(
	p puller.Puller,
	mysqlClient mysql.Client,
	v3Cache cache.MultiGetDeleteCountCache,
	opts ...Option,
) *Persister {
	dopts := &options{
		maxMPS:             1000,
		numWorkers:         1,
		flushSize:          50,
		flushInterval:      5 * time.Second,
		writeCacheInterval: 1 * time.Minute,
		logger:             zap.NewNop(),
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
		envLastUsedCache:      environmentLastUsedInfoCache{},
		mysqlClient:           mysqlClient,
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
	p.group.Go(p.writeFlagLastUsedInfoCache)
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
	updateEvaluationCounter := func(envEvents environmentEventMap) {
		// Increment the evaluation event count in the Redis
		fails := p.incrementEnvEvents(envEvents)
		// Check to Ack or Nack the messages
		p.checkMessages(batch, fails)
		// Reset the maps and the timer
		batch = make(map[string]*puller.Message)
		timer.Reset(p.opts.flushInterval)
	}
	for {
		select {
		case msg, ok := <-p.puller.MessageCh():
			if !ok {
				p.logger.Error("Failed to pull message")
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
			envEvents := p.extractEvents(batch)
			// Update the feature flag last-used cache
			p.cacheLastUsedInfoPerEnv(envEvents)
			updateEvaluationCounter(envEvents)
		case <-timer.C:
			envEvents := p.extractEvents(batch)
			// Update the feature flag last-used cache
			p.cacheLastUsedInfoPerEnv(envEvents)
			updateEvaluationCounter(envEvents)
		case <-p.ctx.Done():
			// Nack the messages to be redelivered
			for _, msg := range batch {
				msg.Nack()
			}
			p.logger.Info("All the left messages were Nack successfully before shutting down",
				zap.Int("batchSize", len(batch)))
			return nil
		}
	}
}

func (p *Persister) incrementEnvEvents(envEvents environmentEventMap) map[string]bool {
	fails := make(map[string]bool, len(envEvents))
	for environmentNamespace, events := range envEvents {
		for id, event := range events {
			// Increment the evaluation event count in the Redis
			if err := p.incrementEvaluationCount(event, environmentNamespace); err != nil {
				p.logger.Error(
					"Failed to increment the evaluation event in the Redis",
					zap.Error(err),
					zap.String("id", id),
					zap.String("environmentNamespace", environmentNamespace),
				)
				fails[id] = true
			}
		}
	}
	return fails
}

func (p *Persister) checkMessages(messages map[string]*puller.Message, fails map[string]bool) {
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
		p.logger.Error("Bad proto message", zap.Error(err), zap.Any("msg", m))
		handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
	}
	for _, m := range messages {
		event := &eventproto.Event{}
		if err := proto.Unmarshal(m.Data, event); err != nil {
			handleBadMessage(m, err)
			continue
		}
		innerEvent := &eventproto.EvaluationEvent{}
		if err := ptypes.UnmarshalAny(event.Event, innerEvent); err != nil {
			handleBadMessage(m, err)
			continue
		}
		if innerEvents, ok := envEvents[event.EnvironmentNamespace]; ok {
			innerEvents[event.Id] = innerEvent
			continue
		}
		envEvents[event.EnvironmentNamespace] = eventMap{event.Id: innerEvent}
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

func (p *Persister) incrementEvaluationCount(event proto.Message, environmentNamespace string) error {
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

func (p *Persister) cacheLastUsedInfoPerEnv(envEvents environmentEventMap) {
	for environmentNamespace, events := range envEvents {
		for _, event := range events {
			p.cacheEnvLastUsedInfo(event, environmentNamespace)
		}
		p.logger.Debug("Cache has been updated",
			zap.String("environmentNamespace", environmentNamespace),
			zap.Int("cacheSize", len(p.envLastUsedCache[environmentNamespace])),
			zap.Int("eventSize", len(events)),
		)
	}
}

func (p *Persister) cacheEnvLastUsedInfo(
	event *eventproto.EvaluationEvent,
	environmentNamespace string,
) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var clientVersion string
	if event.User == nil {
		p.logger.Warn("Failed to cache last used info. User is nil.",
			zap.String("environmentNamespace", environmentNamespace))
	} else {
		clientVersion = event.User.Data[userDataAppVersion]
	}
	id := ftdomain.FeatureLastUsedInfoID(event.FeatureId, event.FeatureVersion)
	if cache, ok := p.envLastUsedCache[environmentNamespace]; ok {
		if info, ok := cache[id]; ok {
			info.UsedAt(event.Timestamp)
			if err := info.SetClientVersion(clientVersion); err != nil {
				p.logger.Error("Failed to set client version",
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
					zap.String("featureId", info.FeatureId),
					zap.Int32("featureVersion", info.Version),
					zap.String("clientVersion", clientVersion))
			}
			return
		}
		cache[id] = ftdomain.NewFeatureLastUsedInfo(
			event.FeatureId,
			event.FeatureVersion,
			event.Timestamp,
			clientVersion,
		)
		return
	}
	cache := lastUsedInfoCache{}
	cache[id] = ftdomain.NewFeatureLastUsedInfo(
		event.FeatureId,
		event.FeatureVersion,
		event.Timestamp,
		clientVersion,
	)
	p.envLastUsedCache[environmentNamespace] = cache
}

// Write the feature flag last-used cache in the MySQL and reset the cache
func (p *Persister) writeFlagLastUsedInfoCache() error {
	timer := time.NewTimer(p.opts.writeCacheInterval)
	for {
		select {
		case <-p.ctx.Done():
			return nil
		case <-timer.C:
			p.logger.Debug("Write cache timer triggered")
			p.writeEnvLastUsedInfo()
			timer.Reset(p.opts.writeCacheInterval)
		}
	}
}

func (p *Persister) writeEnvLastUsedInfo() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for environmentNamespace, cache := range p.envLastUsedCache {
		info := make([]*ftdomain.FeatureLastUsedInfo, 0, len(cache))
		for _, v := range cache {
			info = append(info, v)
		}
		if err := p.upsertMultiFeatureLastUsedInfo(context.Background(), info, environmentNamespace); err != nil {
			p.logger.Error("Failed to write feature last-used info", zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace))
			continue
		}
		p.logger.Debug("Cache has been written",
			zap.String("environmentNamespace", environmentNamespace),
			zap.Int("cacheSize", len(info)),
		)
	}
	// Reset the cache
	p.envLastUsedCache = make(environmentLastUsedInfoCache)
}

func (p *Persister) upsertMultiFeatureLastUsedInfo(
	ctx context.Context,
	featureLastUsedInfos []*ftdomain.FeatureLastUsedInfo,
	environmentNamespace string,
) error {
	ids := make([]string, 0, len(featureLastUsedInfos))
	for _, f := range featureLastUsedInfos {
		ids = append(ids, f.ID())
	}
	storage := ftstorage.NewFeatureLastUsedInfoStorage(p.mysqlClient)
	updatedInfo := make([]*ftdomain.FeatureLastUsedInfo, 0, len(ids))
	currentInfo, err := storage.GetFeatureLastUsedInfos(ctx, ids, environmentNamespace)
	if err != nil {
		return err
	}
	currentInfoMap := make(map[string]*ftdomain.FeatureLastUsedInfo, len(currentInfo))
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
		if err := p.upsertFeatureLastUsedInfo(ctx, info, environmentNamespace); err != nil {
			return err
		}
	}
	return nil
}

func (p *Persister) upsertFeatureLastUsedInfo(
	ctx context.Context,
	featureLastUsedInfo *ftdomain.FeatureLastUsedInfo,
	environmentNamespace string,
) error {
	storage := ftstorage.NewFeatureLastUsedInfoStorage(p.mysqlClient)
	if err := storage.UpsertFeatureLastUsedInfo(
		ctx,
		featureLastUsedInfo,
		environmentNamespace,
	); err != nil {
		return err
	}
	return nil
}
