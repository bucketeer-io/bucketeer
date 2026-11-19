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

	"github.com/blang/semver"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/protobuf/types/known/wrapperspb"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
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

func holdoutGroups(group *featureproto.HoldoutGroup) []*featureproto.HoldoutGroup {
	if group == nil {
		return nil
	}
	return []*featureproto.HoldoutGroup{group}
}

func countHoldoutEvaluation(group *featureproto.HoldoutGroup, sourceID string) {
	if group == nil || len(group.Features) == 0 {
		return
	}
	holdoutEvaluationsCounter.WithLabelValues(sourceID).Inc()
}

// sdkHoldoutGroups is the holdout_groups list of a GetFeatureFlags response:
// the serving group, or nil when nothing is serving. A clearance snapshot and
// an unavailable snapshot are both "nothing is serving" here. A server SDK
// evaluates with its own copy, so an unreadable snapshot only means the SDK
// evaluates without the holdout until the snapshot is back; the request is
// not failed for it.
//
// The group is copied down to what local evaluation reads: the assignment,
// the revision and the current targets. The name, the description, the goals
// and the editors' email addresses stay on the server. The cached group is
// shared between requests and is never modified.
func sdkHoldoutGroups(group *featureproto.HoldoutGroup) []*featureproto.HoldoutGroup {
	serving := evaluation.ServingHoldout(group)
	if serving == nil {
		return nil
	}
	var targets []*featureproto.HoldoutGroupFeature
	for _, target := range serving.Features {
		if target == nil || target.RemovedAt != 0 {
			continue
		}
		targets = append(targets, &featureproto.HoldoutGroupFeature{
			Id:                 target.Id,
			HoldoutGroupId:     target.HoldoutGroupId,
			FeatureId:          target.FeatureId,
			EnvironmentId:      target.EnvironmentId,
			HoldoutVariationId: target.HoldoutVariationId,
			AddedAt:            target.AddedAt,
		})
	}
	return []*featureproto.HoldoutGroup{{
		Id: serving.Id,
		// The environment ID is part of the HASH bucketing input.
		EnvironmentId:  serving.EnvironmentId,
		AssignmentType: serving.AssignmentType,
		SamplingSeed:   serving.SamplingSeed,
		BucketStart:    serving.BucketStart,
		Percentage:     serving.Percentage,
		SegmentId:      serving.SegmentId,
		Status:         serving.Status,
		UpdatedAt:      serving.UpdatedAt,
		Features:       targets,
	}}
}

// holdoutChangedSince reports whether the holdout definition has a revision
// at or after since: the serving group was started or edited, or serving
// stopped (the clearance snapshot keeps the stop time).
//
// GetFeatureFlags uses it to tell why feature_flags_id changed while no flag
// did. The comparison includes since itself, like the flag diff, because
// updated_at is whole seconds. The unavailable snapshot is not a revision:
// its updated_at is earlier than any cursor that reaches the diff, so it
// reports false and the caller falls back to a full resend.
func holdoutChangedSince(group *featureproto.HoldoutGroup, since int64) bool {
	if group == nil || group.UpdatedAt == 0 {
		return false
	}
	return group.UpdatedAt >= since
}

// holdoutMinSDKVersions lists, per source, the first SDK release that applies
// the holdout_groups of a GetFeatureFlags response.
//
// It is empty because no released SDK does yet. A source that is not listed
// is reported as incompatible, which is also right for a wrapper such as an
// OpenFeature provider: it reports its own version, so it needs its own
// entry. Add an entry when an SDK release ships holdout support.
var holdoutMinSDKVersions = map[eventproto.SourceId]semver.Version{}

// holdoutCapableSDK reports whether the SDK applies holdout groups when it
// evaluates locally. A version that cannot be parsed is not capable.
func holdoutCapableSDK(
	minVersions map[eventproto.SourceId]semver.Version,
	sourceID eventproto.SourceId,
	sdkVersion string,
) bool {
	minVersion, ok := minVersions[sourceID]
	if !ok {
		return false
	}
	version, err := semver.ParseTolerant(sdkVersion)
	if err != nil {
		return false
	}
	return version.GTE(minVersion)
}

// countHoldoutIncompatibleSDK counts a GetFeatureFlags request from an SDK
// that ignores holdout_groups while a serving group has targets. Such an SDK
// serves the normal variation to held-out users (RFC 0048 §4.3.2). The
// request is still answered in full.
func countHoldoutIncompatibleSDK(
	groups []*featureproto.HoldoutGroup,
	minVersions map[eventproto.SourceId]semver.Version,
	environmentID string,
	sourceID eventproto.SourceId,
	sdkVersion string,
) {
	targeted := false
	for _, group := range groups {
		if group != nil && len(group.Features) > 0 {
			targeted = true
			break
		}
	}
	if !targeted || holdoutCapableSDK(minVersions, sourceID, sdkVersion) {
		return
	}
	holdoutIncompatibleSDKCounter.WithLabelValues(environmentID, sourceID.String(), sdkVersion).Inc()
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
