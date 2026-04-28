// Copyright 2026 The Bucketeer Authors.
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

package processor

import (
	"context"
	"errors"
	"hash/fnv"
	"net"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/types/known/wrapperspb"

	accstorage "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	autoopsclient "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	domaineventdomain "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	featuredomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	// experimentRefreshLookback matches the batch experiment-cacher (see
	// pkg/batch/jobs/cacher/experimentcacher.go): we keep STOPPED
	// experiments in cache for two days after stop so late-arriving
	// evaluation/goal events still resolve.
	experimentRefreshLookback = 2 * 24 * time.Hour
)

const (
	// cacheRefresherListPageSize matches the page size used by the API
	// service's listFeatures path so the on-disk request shape is identical.
	cacheRefresherListPageSize = 1000

	// cacheRefresherFetchTimeout caps how long a single refresh attempt may
	// spend talking to upstream services (feature service / MySQL). It is
	// generous on purpose: a few large environments take several seconds to
	// page through.
	cacheRefresherFetchTimeout = 30 * time.Second

	// cacheRefresherWorkerCount is the number of per-environment serial
	// workers. Events for the same environmentId are routed (by hash) to the
	// same worker so writes for that environment never race; events for
	// different environments are processed in parallel across workers.
	cacheRefresherWorkerCount = 16

	// cacheRefresherWorkerQueueSize is the buffer per worker. If the queue
	// fills, the dispatcher blocks (back-pressuring the puller); we prefer
	// pull-side back-pressure over dropping events.
	cacheRefresherWorkerQueueSize = 256
)

var errCacheRefresherBadMessage = errors.New("cache refresher bad message")

// cacheRefresher subscribes to the domain-event topic and, for each
// flag/segment/api-key change, refreshes the L2 (Redis) cache from
// MySQL and then announces the change on the cache-invalidation topic so
// every API pod can drop its L1 (in-memory) entry.
//
// This replaces the previous evict-only behaviour. With refresh-on-event,
// the L2 cache is always populated for hot paths; API pods reload from a
// warm L2 instead of fanning out to MySQL through the singleflight safety
// net, eliminating the cache-miss thundering herd that follows every flag
// change.
type cacheRefresher struct {
	featureClient              featureclient.Client
	experimentClient           experimentclient.Client
	autoOpsClient              autoopsclient.Client
	accountStorage             accstorage.AccountStorage
	featuresCache              cachev3.FeaturesCache
	segmentUsersCache          cachev3.SegmentUsersCache
	environmentAPIKeyCache     cachev3.EnvironmentAPIKeyCache
	experimentsCache           cachev3.ExperimentsCache
	autoOpsRulesCache          cachev3.AutoOpsRulesCache
	cacheInvalidationPublisher publisher.Publisher
	logger                     *zap.Logger
}

// NewCacheRefresher returns a processor that refreshes L2 caches and announces
// invalidations on the cache-invalidation topic.
//
// cacheInvalidationPublisher may be nil for tests that exercise only the
// refresh path; in production wiring it must be supplied.
func NewCacheRefresher(
	featureClient featureclient.Client,
	experimentClient experimentclient.Client,
	autoOpsClient autoopsclient.Client,
	accountStorage accstorage.AccountStorage,
	featuresCache cachev3.FeaturesCache,
	segmentUsersCache cachev3.SegmentUsersCache,
	environmentAPIKeyCache cachev3.EnvironmentAPIKeyCache,
	experimentsCache cachev3.ExperimentsCache,
	autoOpsRulesCache cachev3.AutoOpsRulesCache,
	cacheInvalidationPublisher publisher.Publisher,
	logger *zap.Logger,
) subscriber.PubSubProcessor {
	return &cacheRefresher{
		featureClient:              featureClient,
		experimentClient:           experimentClient,
		autoOpsClient:              autoOpsClient,
		accountStorage:             accountStorage,
		featuresCache:              featuresCache,
		segmentUsersCache:          segmentUsersCache,
		environmentAPIKeyCache:     environmentAPIKeyCache,
		experimentsCache:           experimentsCache,
		autoOpsRulesCache:          autoOpsRulesCache,
		cacheInvalidationPublisher: cacheInvalidationPublisher,
		logger:                     logger.Named("cache-refresher"),
	}
}

// Process consumes domain events and dispatches them to per-environment
// serial workers. Events for the same environmentId are processed in order;
// events for different environments are processed concurrently across
// workers.
func (c *cacheRefresher) Process(ctx context.Context, msgChan <-chan *puller.Message) error {
	workers := make([]chan *puller.Message, cacheRefresherWorkerCount)
	doneCh := make(chan struct{}, cacheRefresherWorkerCount)
	for i := 0; i < cacheRefresherWorkerCount; i++ {
		workers[i] = make(chan *puller.Message, cacheRefresherWorkerQueueSize)
		go c.runWorker(ctx, workers[i], doneCh)
	}
	defer func() {
		for _, w := range workers {
			close(w)
		}
		for i := 0; i < cacheRefresherWorkerCount; i++ {
			<-doneCh
		}
	}()
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				c.logger.Error("cacheRefresher: message channel closed")
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(subscriberCacheRefresher).Inc()
			c.dispatch(ctx, msg, workers)
		case <-ctx.Done():
			c.logger.Debug("cacheRefresher: context done, stopped processing messages")
			return nil
		}
	}
}

// dispatch routes a message to the worker responsible for its environmentId.
// If the message cannot be parsed enough to extract an environment, the
// message is acked (bad message — retrying will not help).
func (c *cacheRefresher) dispatch(
	ctx context.Context,
	msg *puller.Message,
	workers []chan *puller.Message,
) {
	envID := extractEnvironmentID(msg.Data)
	idx := workerIndex(envID, len(workers))
	select {
	case workers[idx] <- msg:
	case <-ctx.Done():
		// Shutdown in-flight: nack so the message is redelivered by the
		// next pod / restart. Avoids losing events on graceful shutdown.
		msg.Nack()
	}
}

func (c *cacheRefresher) runWorker(
	ctx context.Context,
	queue <-chan *puller.Message,
	doneCh chan<- struct{},
) {
	defer func() { doneCh <- struct{}{} }()
	for msg := range queue {
		// Honour context cancellation between messages; in-flight handling
		// will still complete (so we don't leave Redis half-written).
		if ctx.Err() != nil {
			msg.Nack()
			continue
		}
		c.handleMessage(ctx, msg)
	}
}

// extractEnvironmentID best-effort decodes the environment id off the wire
// for routing. We tolerate unmarshal failures here; handleMessage will
// re-decode and surface the real error.
func extractEnvironmentID(data []byte) string {
	event := &domaineventproto.Event{}
	if err := proto.Unmarshal(data, event); err != nil {
		return ""
	}
	return event.EnvironmentId
}

// workerIndex hashes envID into [0, n). Empty envID maps to bucket 0; this
// is fine because the bad-message path will ack-and-drop the event.
func workerIndex(envID string, n int) int {
	if n <= 0 {
		return 0
	}
	if envID == "" {
		return 0
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(envID))
	return int(h.Sum32() % uint32(n))
}

func (c *cacheRefresher) handleMessage(ctx context.Context, msg *puller.Message) {
	event := &domaineventproto.Event{}
	if err := proto.Unmarshal(msg.Data, event); err != nil {
		c.logger.Error("Failed to unmarshal domain event",
			zap.Error(err),
			zap.String("msgID", msg.ID),
		)
		subscriberHandledCounter.WithLabelValues(subscriberCacheRefresher, codes.BadMessage.String()).Inc()
		msg.Ack()
		return
	}
	if err := c.refresh(ctx, event); err != nil {
		if errors.Is(err, errCacheRefresherBadMessage) {
			subscriberHandledCounter.WithLabelValues(subscriberCacheRefresher, codes.BadMessage.String()).Inc()
			msg.Ack()
			return
		}
		if isRepeatable(err) {
			c.logger.Warn("Failed to refresh cache with repeatable error",
				zap.Error(err),
				zap.String("environmentId", event.EnvironmentId),
				zap.String("entityId", event.EntityId),
				zap.String("entityType", event.EntityType.String()),
				zap.String("type", event.Type.String()),
			)
			subscriberHandledCounter.WithLabelValues(subscriberCacheRefresher, codes.RepeatableError.String()).Inc()
			msg.Nack()
			return
		}
		c.logger.Error("Failed to refresh cache with non-repeatable error",
			zap.Error(err),
			zap.String("environmentId", event.EnvironmentId),
			zap.String("entityId", event.EntityId),
			zap.String("entityType", event.EntityType.String()),
			zap.String("type", event.Type.String()),
		)
		subscriberHandledCounter.WithLabelValues(subscriberCacheRefresher, codes.NonRepeatableError.String()).Inc()
		msg.Ack()
		return
	}
	subscriberHandledCounter.WithLabelValues(subscriberCacheRefresher, codes.OK.String()).Inc()
	msg.Ack()
}

// refresh performs the fetch -> Put -> publish sequence for one event.
// Per the plan, EXPERIMENT and AUTOOPS_RULE remain plain evict (they aren't
// on the SDK hot path); FEATURE / SEGMENT / APIKEY are all rewrite-then-
// announce.
func (c *cacheRefresher) refresh(ctx context.Context, event *domaineventproto.Event) error {
	switch event.EntityType {
	case domaineventproto.Event_FEATURE:
		return c.refreshFeatures(ctx, event)
	case domaineventproto.Event_SEGMENT:
		return c.refreshSegmentUsers(ctx, event)
	case domaineventproto.Event_APIKEY:
		return c.refreshAPIKey(ctx, event)
	case domaineventproto.Event_EXPERIMENT:
		return c.refreshExperiments(ctx, event)
	case domaineventproto.Event_AUTOOPS_RULE:
		return c.refreshAutoOpsRules(ctx, event)
	default:
		return nil
	}
}

func (c *cacheRefresher) refreshFeatures(
	ctx context.Context,
	event *domaineventproto.Event,
) error {
	fetchCtx, cancel := context.WithTimeout(ctx, cacheRefresherFetchTimeout)
	defer cancel()
	features, err := c.fetchAllFeatures(fetchCtx, event.EnvironmentId)
	if err != nil {
		return err
	}
	if err := c.featuresCache.Put(features, event.EnvironmentId); err != nil {
		return err
	}
	c.logger.Debug("Refreshed features redis cache",
		zap.String("environmentId", event.EnvironmentId),
		zap.String("entityId", event.EntityId),
		zap.String("type", event.Type.String()),
		zap.Int("featuresCount", len(features.Features)),
	)
	return c.publishCacheInvalidation(ctx, event)
}

// fetchAllFeatures pages through ListFeatures the same way the API service
// does (see grpcGatewayService.listFeatures), so the cached blob shape is
// byte-equivalent to what the API would have written on a cold miss.
func (c *cacheRefresher) fetchAllFeatures(
	ctx context.Context,
	environmentID string,
) (*featureproto.Features, error) {
	features := []*featureproto.Feature{}
	cursor := ""
	for {
		resp, err := c.featureClient.ListFeatures(ctx, &featureproto.ListFeaturesRequest{
			PageSize:      cacheRefresherListPageSize,
			Cursor:        cursor,
			EnvironmentId: environmentID,
		})
		if err != nil {
			return nil, err
		}
		for _, f := range resp.Features {
			ff := featuredomain.Feature{Feature: f}
			if ff.IsDisabledAndOffVariationEmpty() {
				continue
			}
			if ff.IsArchivedBeforeLastThirtyDays() {
				continue
			}
			features = append(features, f)
		}
		featureSize := len(resp.Features)
		if featureSize == 0 || featureSize < cacheRefresherListPageSize {
			break
		}
		cursor = resp.Cursor
	}
	return &featureproto.Features{Features: features}, nil
}

func (c *cacheRefresher) refreshSegmentUsers(
	ctx context.Context,
	event *domaineventproto.Event,
) error {
	if event.EntityId == "" {
		return errCacheRefresherBadMessage
	}
	// SEGMENT_DELETED removes the segment row in MySQL, so a subsequent
	// GetSegment would return NotFound and the refresh would be classified
	// as a non-repeatable error and dropped — leaving the stale users blob
	// in L2 indefinitely. Handle the delete explicitly: evict L2 and
	// announce the invalidation so api pods drop their L1 entries.
	if event.Type == domaineventproto.Event_SEGMENT_DELETED {
		if err := c.segmentUsersCache.Evict(event.EntityId, event.EnvironmentId); err != nil {
			return err
		}
		c.logger.Debug("Evicted segment users redis cache (segment deleted)",
			zap.String("environmentId", event.EnvironmentId),
			zap.String("segmentId", event.EntityId),
			zap.String("type", event.Type.String()),
		)
		return c.publishCacheInvalidation(ctx, event)
	}
	fetchCtx, cancel := context.WithTimeout(ctx, cacheRefresherFetchTimeout)
	defer cancel()
	listResp, err := c.featureClient.ListSegmentUsers(fetchCtx, &featureproto.ListSegmentUsersRequest{
		SegmentId:     event.EntityId,
		EnvironmentId: event.EnvironmentId,
	})
	if err != nil {
		return err
	}
	getResp, err := c.featureClient.GetSegment(fetchCtx, &featureproto.GetSegmentRequest{
		Id:            event.EntityId,
		EnvironmentId: event.EnvironmentId,
	})
	if err != nil {
		return err
	}
	segmentUsers := &featureproto.SegmentUsers{
		SegmentId: event.EntityId,
		Users:     listResp.Users,
		UpdatedAt: getResp.Segment.UpdatedAt,
	}
	if err := c.segmentUsersCache.Put(segmentUsers, event.EnvironmentId); err != nil {
		return err
	}
	c.logger.Debug("Refreshed segment users redis cache",
		zap.String("environmentId", event.EnvironmentId),
		zap.String("segmentId", event.EntityId),
		zap.String("type", event.Type.String()),
		zap.Int("usersCount", len(segmentUsers.Users)),
	)
	return c.publishCacheInvalidation(ctx, event)
}

func (c *cacheRefresher) refreshAPIKey(
	ctx context.Context,
	event *domaineventproto.Event,
) error {
	secrets, err := domaineventdomain.ExtractAPIKeySecrets(event)
	if err != nil {
		if len(secrets) > 0 {
			c.logger.Warn(
				"Partially failed to extract api_key from entity data; refreshing available secrets",
				zap.Error(err),
				zap.String("environmentId", event.EnvironmentId),
				zap.String("entityId", event.EntityId),
				zap.String("type", event.Type.String()),
			)
		} else {
			c.logger.Error(
				"Failed to extract api_key from entity data",
				zap.Error(err),
				zap.String("environmentId", event.EnvironmentId),
				zap.String("entityId", event.EntityId),
				zap.String("type", event.Type.String()),
			)
			return errCacheRefresherBadMessage
		}
	}
	if len(secrets) == 0 {
		return nil
	}
	fetchCtx, cancel := context.WithTimeout(ctx, cacheRefresherFetchTimeout)
	defer cancel()
	for _, secret := range secrets {
		domainEnvAPIKey, err := c.accountStorage.GetEnvironmentAPIKey(fetchCtx, secret)
		if err != nil {
			// If the secret was rotated away and the row is gone, evict so
			// stale L2 entries don't authorise a removed key.
			if errors.Is(err, accstorage.ErrAPIKeyNotFound) {
				if evictErr := c.environmentAPIKeyCache.Evict(secret); evictErr != nil {
					return evictErr
				}
				c.logger.Debug("Evicted environment API key (no longer in DB)",
					zap.String("environmentId", event.EnvironmentId),
					zap.String("entityId", event.EntityId),
					zap.String("type", event.Type.String()),
				)
				continue
			}
			return err
		}
		if err := c.environmentAPIKeyCache.Put(domainEnvAPIKey.EnvironmentAPIKey); err != nil {
			return err
		}
		c.logger.Debug("Refreshed environment API key redis cache",
			zap.String("environmentId", event.EnvironmentId),
			zap.String("entityId", event.EntityId),
			zap.String("type", event.Type.String()),
		)
	}
	return c.publishCacheInvalidation(ctx, event)
}

// refreshExperiments rebuilds the experiments L2 cache for the affected
// environment. Mirrors the batch experiment-cacher: only RUNNING and
// STOPPED experiments are cached, and STOPPED ones drop out
// experimentRefreshLookback after their stop_at because late-arriving
// evaluation/goal events still need to resolve them.
func (c *cacheRefresher) refreshExperiments(
	ctx context.Context,
	event *domaineventproto.Event,
) error {
	fetchCtx, cancel := context.WithTimeout(ctx, cacheRefresherFetchTimeout)
	defer cancel()
	resp, err := c.experimentClient.ListExperiments(fetchCtx, &experimentproto.ListExperimentsRequest{
		PageSize:      0,
		EnvironmentId: event.EnvironmentId,
		StopAt:        time.Now().Add(-experimentRefreshLookback).Unix(),
		Statuses: []experimentproto.Experiment_Status{
			experimentproto.Experiment_RUNNING,
			experimentproto.Experiment_STOPPED,
		},
		Archived: &wrapperspb.BoolValue{Value: false},
	})
	if err != nil {
		return err
	}
	experiments := &experimentproto.Experiments{Experiments: resp.Experiments}
	if err := c.experimentsCache.Put(experiments, event.EnvironmentId); err != nil {
		return err
	}
	c.logger.Debug("Refreshed experiments redis cache",
		zap.String("environmentId", event.EnvironmentId),
		zap.String("entityId", event.EntityId),
		zap.String("type", event.Type.String()),
		zap.Int("experimentsCount", len(experiments.Experiments)),
	)
	return c.publishCacheInvalidation(ctx, event)
}

// refreshAutoOpsRules rebuilds the auto-ops-rules L2 cache for the
// affected environment. Mirrors the batch auto-ops-rules-cacher.
func (c *cacheRefresher) refreshAutoOpsRules(
	ctx context.Context,
	event *domaineventproto.Event,
) error {
	fetchCtx, cancel := context.WithTimeout(ctx, cacheRefresherFetchTimeout)
	defer cancel()
	resp, err := c.autoOpsClient.ListAutoOpsRules(fetchCtx, &autoopsproto.ListAutoOpsRulesRequest{
		PageSize:      0,
		EnvironmentId: event.EnvironmentId,
	})
	if err != nil {
		return err
	}
	autoOpsRules := &autoopsproto.AutoOpsRules{AutoOpsRules: resp.AutoOpsRules}
	if err := c.autoOpsRulesCache.Put(autoOpsRules, event.EnvironmentId); err != nil {
		return err
	}
	c.logger.Debug("Refreshed auto ops rules redis cache",
		zap.String("environmentId", event.EnvironmentId),
		zap.String("entityId", event.EntityId),
		zap.String("type", event.Type.String()),
		zap.Int("autoOpsRulesCount", len(autoOpsRules.AutoOpsRules)),
	)
	return c.publishCacheInvalidation(ctx, event)
}

// publishCacheInvalidation announces the change on the cache-invalidation
// topic so every consumer (today: api pods' L1 in-memory cache) can
// reconcile its own layer. We reuse domain.Event verbatim to avoid a new
// proto; cache_invalidator on the consumer side already understands this
// shape (see pkg/api/api/cache_invalidator.go).
//
// The publisher may be nil in tests; in that case publish is a no-op.
func (c *cacheRefresher) publishCacheInvalidation(
	ctx context.Context,
	event *domaineventproto.Event,
) error {
	if c.cacheInvalidationPublisher == nil {
		return nil
	}
	return c.cacheInvalidationPublisher.Publish(ctx, event)
}

// isRepeatable returns true for transient errors that should be retried
// (NACKed). DeadlineExceeded, context.Canceled, network timeouts, and known
// connection-level strings are repeatable; cache lookup misses are not.
//
// Why context.Canceled is repeatable: the only thing that cancels a worker
// context here is graceful shutdown (SIGTERM). If a refresh has already
// written to L2 and the publish step then fails because shutdown is in
// progress, ack-and-drop would leave api pods with stale L1 entries until
// TTL. Nacking instead lets the message be redelivered to a healthy pod,
// which idempotently re-runs the refresh and the publish.
func isRepeatable(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) {
		return true
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	if errors.Is(err, cache.ErrNotFound) || errors.Is(err, cache.ErrInvalidType) {
		return false
	}
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		return true
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "eof") {
		return true
	}
	return false
}
