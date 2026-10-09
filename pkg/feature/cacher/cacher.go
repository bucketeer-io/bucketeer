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

// Package cacher provides functionality to sync feature flags from MySQL to Redis cache.
//
//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package cacher

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	ftdomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	ftstorage "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// FeatureFlagCacher provides functionality to sync feature flags from MySQL to Redis.
// This is used by:
// - The batch job to periodically refresh the cache for all environments
// - Auto-ops watchers to immediately update the cache after flag changes
type FeatureFlagCacher interface {
	// RefreshEnvironmentCache updates the Redis cache for a specific environment.
	// This should be called after auto operations (Schedule, Kill Switch, Progressive Rollout)
	// to ensure SDKs receive the updated flags immediately.
	RefreshEnvironmentCache(ctx context.Context, environmentID string) error

	// RefreshAllEnvironmentCaches updates the Redis cache for all environments.
	// This is used by the periodic batch job.
	RefreshAllEnvironmentCaches(ctx context.Context) error
}

type featureFlagCacher struct {
	ftStorage ftstorage.FeatureStorage
	caches    []cachev3.FeaturesCache
	logger    *zap.Logger
	// now supplies the snapshot generation; overridable in tests.
	now func() time.Time
}

// NewFeatureFlagCacher creates a new FeatureFlagCacher.
func NewFeatureFlagCacher(
	ftStorage ftstorage.FeatureStorage,
	multiCaches []cache.MultiGetCache,
	logger *zap.Logger,
) FeatureFlagCacher {
	caches := make([]cachev3.FeaturesCache, 0, len(multiCaches))
	for _, c := range multiCaches {
		caches = append(caches, cachev3.NewFeaturesCache(c, 0))
	}
	return &featureFlagCacher{
		ftStorage: ftStorage,
		caches:    caches,
		logger:    logger.Named("feature-flag-cacher"),
		now:       time.Now,
	}
}

// RefreshEnvironmentCache updates the Redis cache for a specific environment.
func (c *featureFlagCacher) RefreshEnvironmentCache(ctx context.Context, environmentID string) error {
	startTime := c.now()
	// The generation is the time the read started: a snapshot read later
	// reflects at least everything an earlier read did, so "later read wins"
	// is the ordering FeaturesCache.PutIfNewer enforces. It must be captured
	// before the query is issued.
	generation := startTime.UnixNano()

	// Use targeted query for single environment instead of fetching all environments
	features, err := c.ftStorage.ListFeaturesByEnvironment(ctx, environmentID)
	if err != nil {
		c.logger.Error("Failed to list features for cache update",
			zap.Error(err),
			zap.String("environmentId", environmentID),
		)
		recordListFeatures(cacherTypeFeatureFlag, scopeSingle, environmentID, codeFail, time.Since(startTime).Seconds())
		return err
	}
	recordListFeatures(cacherTypeFeatureFlag, scopeSingle, environmentID, codeSuccess, time.Since(startTime).Seconds())

	filtered := c.removeOldFeatures(features)
	fts := &ftproto.Features{
		Id:       evaluation.GenerateFeaturesID(filtered),
		Features: filtered,
	}
	c.putCache(fts, environmentID, len(filtered), generation)

	return nil
}

// RefreshAllEnvironmentCaches updates the Redis cache for all environments.
func (c *featureFlagCacher) RefreshAllEnvironmentCaches(ctx context.Context) error {
	startTime := c.now()
	// One generation for the whole run: every environment's snapshot comes
	// from the same read, so they all share its start time.
	generation := startTime.UnixNano()

	envFts, err := c.ftStorage.ListAllEnvironmentFeatures(ctx)
	if err != nil {
		c.logger.Error("Failed to list all environment features")
		// Use scopeBatch with environmentIDAll for batch operations covering all environments
		recordListFeatures(cacherTypeFeatureFlag, scopeBatch, environmentIDAll, codeFail, time.Since(startTime).Seconds())
		return err
	}
	// Use scopeBatch with environmentIDAll for batch operations covering all environments
	recordListFeatures(cacherTypeFeatureFlag, scopeBatch, environmentIDAll, codeSuccess, time.Since(startTime).Seconds())

	for _, envFt := range envFts {
		filtered := c.removeOldFeatures(envFt.Features)
		fts := &ftproto.Features{
			Id:       evaluation.GenerateFeaturesID(filtered),
			Features: filtered,
		}
		c.putCache(fts, envFt.EnvironmentId, len(filtered), generation)
	}

	return nil
}

// removeOldFeatures filters out archived feature flags over thirty days ago.
func (c *featureFlagCacher) removeOldFeatures(features []*ftproto.Feature) []*ftproto.Feature {
	result := make([]*ftproto.Feature, 0, len(features))
	for _, f := range features {
		ft := ftdomain.Feature{Feature: f}
		if !ft.IsDisabledAndOffVariationEmpty() && !ft.IsArchivedBeforeLastThirtyDays() {
			result = append(result, f)
		}
	}
	return result
}

// putCache saves features to all Redis instances and records metrics.
//
// A write rejected as stale (a concurrent writer already cached a newer
// snapshot) is not a failure: the cache holds data at least as fresh as ours.
// It is counted in the stale-rejection metric and otherwise treated as
// success, except that the features-updated gauge is only set when at least
// one instance accepted our snapshot.
func (c *featureFlagCacher) putCache(
	features *ftproto.Features,
	environmentID string,
	featureCount int,
	generation int64,
) {
	var wg sync.WaitGroup
	var hasError bool
	var anyAccepted bool
	var mu sync.Mutex

	for _, cache := range c.caches {
		wg.Add(1)
		go func(cache cachev3.FeaturesCache) {
			defer wg.Done()
			accepted, err := cache.PutIfNewer(features, environmentID, generation)
			if err != nil {
				c.logger.Error("Failed to cache features",
					zap.Error(err),
					zap.String("environmentId", environmentID),
				)
				mu.Lock()
				hasError = true
				mu.Unlock()
				return
			}
			if !accepted {
				cachev3.RecordFeaturesPutRejectedStale(cachev3.FeaturesWriterBatchCacher)
				c.logger.Debug("Skipped caching features: a newer snapshot is already cached",
					zap.String("environmentId", environmentID),
					zap.Int64("generation", generation),
				)
				return
			}
			mu.Lock()
			anyAccepted = true
			mu.Unlock()
		}(cache)
	}
	wg.Wait()

	// Record metrics based on overall success/failure
	if hasError {
		recordCachePut(cacherTypeFeatureFlag, environmentID, codeFail)
		return
	}
	recordCachePut(cacherTypeFeatureFlag, environmentID, codeSuccess)
	if anyAccepted {
		recordFeaturesUpdated(cacherTypeFeatureFlag, environmentID, featureCount)
	}
}
