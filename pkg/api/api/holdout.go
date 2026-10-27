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
	"errors"
	"sync"

	"go.uber.org/zap"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
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

// commitHoldoutFill stores the Redis snapshot when no invalidation landed
// since epoch was observed. It reports whether the snapshot was stored.
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

// loadLiveHoldoutGroup reads the serving snapshot from memory, then redis.
// A miss is "no live group". A read error is an unavailable snapshot, so a
// disabled flag with no off variation stays omitted while the cache is down.
// The gateway does not query MySQL.
func loadLiveHoldoutGroup(
	memory, redis cachev3.HoldoutGroupCache,
	environmentID string,
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
	// A miss and a read error are both a missing snapshot. Serve the stand-in
	// clearance so a disabled flag with no off variation stays omitted, and
	// do not cache it: the next request should see Redis once it recovers.
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
			return evaluation.UnavailableHoldoutSnapshot()
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
	return evaluation.HoldoutTargetFeatureIDs(loadLiveHoldoutGroup(memory, redis, environmentID, logger))
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

func (s *gatewayService) liveHoldoutGroup(environmentID string) *featureproto.HoldoutGroup {
	return loadLiveHoldoutGroup(s.holdoutGroupCache, s.holdoutGroupRedisCache, environmentID, s.logger)
}

func (s *grpcGatewayService) liveHoldoutGroup(environmentID string) *featureproto.HoldoutGroup {
	return loadLiveHoldoutGroup(s.holdoutGroupCache, s.holdoutGroupRedisCache, environmentID, s.logger)
}
