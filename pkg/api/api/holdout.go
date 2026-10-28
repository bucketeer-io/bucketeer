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

package api

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/protobuf/types/known/wrapperspb"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// holdoutFillGate stops a gateway request that read Redis before a stop from
// writing that snapshot back into memory after the invalidator evicts it.
// The invalidator bumps the epoch and evicts under the lock. A fill records
// the epoch before the Redis read and stores the result only when the epoch
// is unchanged.
type holdoutFillGate struct {
	mu    sync.Mutex
	epoch uint64
}

var holdoutFillGates sync.Map // environmentID -> *holdoutFillGate

func holdoutGate(environmentID string) *holdoutFillGate {
	value, _ := holdoutFillGates.LoadOrStore(environmentID, &holdoutFillGate{})
	return value.(*holdoutFillGate)
}

func loadHoldoutFillEpoch(environmentID string) uint64 {
	gate := holdoutGate(environmentID)
	gate.mu.Lock()
	defer gate.mu.Unlock()
	return gate.epoch
}

// invalidateHoldoutFill bumps the epoch and evicts while holding it, so a
// fill that already read the previous snapshot cannot store it afterwards.
func invalidateHoldoutFill(environmentID string, evict func()) {
	gate := holdoutGate(environmentID)
	gate.mu.Lock()
	defer gate.mu.Unlock()
	gate.epoch++
	if evict != nil {
		evict()
	}
}

// commitHoldoutFill stores a Redis snapshot when no invalidation landed
// since epoch was observed. It reports whether the snapshot was stored.
// Feature fills and holdout fills share the epoch: a holdout start evicts
// both caches, and an in-flight feature fill must not put the pre-start
// feature list back afterwards.
func commitHoldoutFill(environmentID string, epoch uint64, put func()) bool {
	gate := holdoutGate(environmentID)
	gate.mu.Lock()
	defer gate.mu.Unlock()
	if gate.epoch != epoch {
		return false
	}
	put()
	return true
}

// holdoutSnapshotFallback reads the serving snapshot when both caches miss.
// It is the L3 of the holdout snapshot, like the feature service is for the
// features cache. nil disables it. epoch is the fill epoch the caller
// observed before its cache reads; a shared read must not outlive it.
type holdoutSnapshotFallback func(ctx context.Context, epoch uint64) (*featureproto.HoldoutGroup, error)

// loadLiveHoldoutGroup reads the serving snapshot from memory, then redis,
// then the fallback. Only when every layer fails is the snapshot unavailable,
// which marks the issued evaluation IDs so the client reconciles later. Redis
// is kept warm by the batch cacher; the fallback covers the window before its
// first run and a flushed cache, so a missing key never degrades evaluation.
func loadLiveHoldoutGroup(
	ctx context.Context,
	memory, redis cachev3.HoldoutGroupCache,
	environmentID string,
	fallback holdoutSnapshotFallback,
	logger *zap.Logger,
) *featureproto.HoldoutGroup {
	readFailed := false
	if memory != nil {
		group, err := memory.Get(environmentID)
		if err == nil {
			return evaluation.EvaluationHoldout(group)
		}
		if !errors.Is(err, cache.ErrNotFound) {
			readFailed = true
			if logger != nil {
				logger.Debug("Failed to get the holdout group from memory",
					zap.Error(err),
					zap.String("environmentId", environmentID),
				)
			}
		}
	}
	if redis == nil {
		if readFailed {
			return evaluation.UnavailableHoldoutSnapshot()
		}
		return nil
	}
	var latest *featureproto.HoldoutGroup
	for attempt := 0; attempt < 2; attempt++ {
		epoch := loadHoldoutFillEpoch(environmentID)
		group, err := redis.Get(environmentID)
		if err != nil {
			if logger != nil && !errors.Is(err, cache.ErrNotFound) {
				logger.Debug("Failed to get the holdout group from redis",
					zap.Error(err),
					zap.String("environmentId", environmentID),
				)
			}
			// A miss and a read error are both a missing snapshot. Read it
			// from the feature service instead. Only memory is filled: Redis
			// belongs to the cacher and the refresher, whose generation
			// guard would otherwise have to arbitrate with every gateway pod.
			if fallback == nil {
				return evaluation.UnavailableHoldoutSnapshot()
			}
			group, err = fallback(ctx, epoch)
			if err != nil {
				if logger != nil {
					logger.Error("Failed to read the holdout snapshot from the feature service",
						zap.Error(err),
						zap.String("environmentId", environmentID),
					)
				}
				return evaluation.UnavailableHoldoutSnapshot()
			}
		}
		latest = group
		if memory == nil {
			return evaluation.EvaluationHoldout(group)
		}
		stored := commitHoldoutFill(environmentID, epoch, func() {
			if putErr := memory.Put(group, environmentID); putErr != nil && logger != nil {
				logger.Debug("Failed to store the holdout group in memory",
					zap.Error(putErr),
					zap.String("environmentId", environmentID),
				)
			}
		})
		if stored {
			return evaluation.EvaluationHoldout(group)
		}
	}
	// The epoch moved on both attempts. Serve the latest read and leave
	// memory to the invalidator, or to a fill that observes a stable epoch.
	return evaluation.EvaluationHoldout(latest)
}

// CachedHoldoutTargetIDs reads the serving snapshot and returns its current
// target feature IDs. A miss or a clearance snapshot returns nil.
func CachedHoldoutTargetIDs(
	memory, redis cachev3.HoldoutGroupCache,
	environmentID string,
	logger *zap.Logger,
) map[string]struct{} {
	return evaluation.HoldoutTargetFeatureIDs(
		loadLiveHoldoutGroup(context.Background(), memory, redis, environmentID, nil, logger),
	)
}

// FillFeatureSnapshot reads a feature snapshot and stores it in memory when
// no invalidation has landed since the read started. REST, gRPC, and the
// stream dispatcher fill L1 through this, so a holdout start cannot be
// overwritten by an in-flight pre-start snapshot. The read is retried once.
// When both attempts race an invalidation, the latest snapshot is returned
// and left out of memory.
func FillFeatureSnapshot(
	environmentID string,
	read func() (*featureproto.Features, error),
	put func(*featureproto.Features),
) (*featureproto.Features, error) {
	var latest *featureproto.Features
	for attempt := 0; attempt < 2; attempt++ {
		epoch := loadHoldoutFillEpoch(environmentID)
		features, err := read()
		if err != nil {
			return nil, err
		}
		latest = features
		if put == nil || commitHoldoutFill(environmentID, epoch, func() { put(features) }) {
			return features, nil
		}
	}
	return latest, nil
}

func holdoutGroups(group *featureproto.HoldoutGroup) []*featureproto.HoldoutGroup {
	if group == nil {
		return nil
	}
	return []*featureproto.HoldoutGroup{group}
}

func appendRemovedFeatureIDs(evals *featureproto.UserEvaluations, removed []string) {
	if evals == nil || len(removed) == 0 {
		return
	}
	evals.ArchivedFeatureIds = append(evals.ArchivedFeatureIds, removed...)
}

func hiddenFeatureIDs(hidden []*featureproto.Feature) []string {
	if len(hidden) == 0 {
		return nil
	}
	ids := make([]string, 0, len(hidden))
	for _, feature := range hidden {
		if feature != nil && feature.Id != "" {
			ids = append(ids, feature.Id)
		}
	}
	return ids
}

func countHoldoutEvaluation(group *featureproto.HoldoutGroup, sourceID string) {
	if group == nil || len(group.Features) == 0 {
		return
	}
	holdoutEvaluationsCounter.WithLabelValues(sourceID).Inc()
}

func (s *gatewayService) liveHoldoutGroup(ctx context.Context, environmentID string) *featureproto.HoldoutGroup {
	return loadLiveHoldoutGroup(
		ctx, s.holdoutGroupCache, s.holdoutGroupRedisCache, environmentID,
		holdoutSnapshotFromService(&s.flightgroup, s.featureClient, environmentID),
		s.logger,
	)
}

func (s *grpcGatewayService) liveHoldoutGroup(ctx context.Context, environmentID string) *featureproto.HoldoutGroup {
	return loadLiveHoldoutGroup(
		ctx, s.holdoutGroupCache, s.holdoutGroupRedisCache, environmentID,
		holdoutSnapshotFromService(&s.flightgroup, s.featureClient, environmentID),
		s.logger,
	)
}

// holdoutSnapshotFromService shares one feature-service read per environment
// and fill epoch across the requests that miss both caches at the same time.
// The epoch is part of the key: a request that starts after an invalidation
// must not join a read begun before it, or commitHoldoutFill would accept
// the pre-invalidation snapshot under the new epoch and restore a stopped
// group to memory. The shared read is detached from the first caller's
// cancellation like singleflightFetch, while each caller still gives up on
// its own context so a closed stream does not wait out the fetch timeout.
func holdoutSnapshotFromService(
	flightgroup *singleflight.Group,
	client featureclient.Client,
	environmentID string,
) holdoutSnapshotFallback {
	if client == nil {
		return nil
	}
	return func(ctx context.Context, epoch uint64) (*featureproto.HoldoutGroup, error) {
		key := fmt.Sprintf("holdout:%s:%d", environmentID, epoch)
		ch := flightgroup.DoChan(key, func() (interface{}, error) {
			fetchCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), singleflightFetchTimeout)
			defer cancel()
			return HoldoutServingSnapshot(fetchCtx, client, environmentID)
		})
		select {
		case res := <-ch:
			if res.Err != nil {
				return nil, res.Err
			}
			return res.Val.(*featureproto.HoldoutGroup), nil
		case <-ctx.Done():
			return nil, fmt.Errorf("%w: %w", errCallerCanceled, ctx.Err())
		}
	}
}

// HoldoutServingSnapshot rebuilds the serving snapshot through the feature
// service, with the same shape as storage.ServingSnapshot: the newest
// serving group with its current targets, or an empty snapshot carrying the
// updated_at of the newest stopped group so a client behind the stop still
// gets a full evaluation.
func HoldoutServingSnapshot(
	ctx context.Context,
	client featureclient.Client,
	environmentID string,
) (*featureproto.HoldoutGroup, error) {
	live, err := client.ListHoldoutGroups(ctx, &featureproto.ListHoldoutGroupsRequest{
		EnvironmentId:  environmentID,
		PageSize:       1,
		OrderBy:        featureproto.ListHoldoutGroupsRequest_UPDATED_AT,
		OrderDirection: featureproto.ListHoldoutGroupsRequest_DESC,
		Statuses: []featureproto.HoldoutGroup_Status{
			featureproto.HoldoutGroup_RUNNING,
			featureproto.HoldoutGroup_ANALYSIS,
		},
		Archived: wrapperspb.Bool(false),
	})
	if err != nil {
		return nil, err
	}
	if len(live.HoldoutGroups) > 0 {
		return live.HoldoutGroups[0], nil
	}
	stopped, err := client.ListHoldoutGroups(ctx, &featureproto.ListHoldoutGroupsRequest{
		EnvironmentId:  environmentID,
		PageSize:       1,
		OrderBy:        featureproto.ListHoldoutGroupsRequest_UPDATED_AT,
		OrderDirection: featureproto.ListHoldoutGroupsRequest_DESC,
		Statuses: []featureproto.HoldoutGroup_Status{
			featureproto.HoldoutGroup_STOPPED,
			featureproto.HoldoutGroup_FORCE_STOPPED,
		},
	})
	if err != nil {
		return nil, err
	}
	var updatedAt int64
	if len(stopped.HoldoutGroups) > 0 {
		updatedAt = stopped.HoldoutGroups[0].UpdatedAt
	}
	return evaluation.ClearanceSnapshot(updatedAt), nil
}
