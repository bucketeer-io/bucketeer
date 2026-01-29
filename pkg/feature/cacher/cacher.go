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
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
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
}

// NewFeatureFlagCacher creates a new FeatureFlagCacher.
func NewFeatureFlagCacher(
	mysqlClient mysql.Client,
	multiCaches []cache.MultiGetCache,
	logger *zap.Logger,
) FeatureFlagCacher {
	caches := make([]cachev3.FeaturesCache, 0, len(multiCaches))
	for _, c := range multiCaches {
		caches = append(caches, cachev3.NewFeaturesCache(c))
	}
	return &featureFlagCacher{
		ftStorage: ftstorage.NewFeatureStorage(mysqlClient),
		caches:    caches,
		logger:    logger.Named("feature-flag-cacher"),
	}
}

// RefreshEnvironmentCache updates the Redis cache for a specific environment.
func (c *featureFlagCacher) RefreshEnvironmentCache(ctx context.Context, environmentID string) error {
	startTime := time.Now()

	// Use targeted query for single environment instead of fetching all environments
	features, err := c.ftStorage.ListFeaturesByEnvironment(ctx, environmentID)
	if err != nil {
		c.logger.Error("Failed to list features for cache update",
			zap.Error(err),
			zap.String("environmentId", environmentID),
		)
		recordListFeatures(scopeSingle, environmentID, codeFail, time.Since(startTime).Seconds())
		return err
	}
	recordListFeatures(scopeSingle, environmentID, codeSuccess, time.Since(startTime).Seconds())

	filtered := c.removeOldFeatures(features)
	fts := &ftproto.Features{
		Id:       evaluation.GenerateFeaturesID(filtered),
		Features: filtered,
	}
	c.putCache(fts, environmentID, len(filtered))

	c.logger.Info("Successfully updated feature flag cache",
		zap.String("environmentId", environmentID),
		zap.Int("featureCount", len(filtered)),
		zap.String("featuresId", fts.Id),
	)

	return nil
}

// RefreshAllEnvironmentCaches updates the Redis cache for all environments.
func (c *featureFlagCacher) RefreshAllEnvironmentCaches(ctx context.Context) error {
	startTime := time.Now()

	envFts, err := c.ftStorage.ListAllEnvironmentFeatures(ctx)
	if err != nil {
		c.logger.Error("Failed to list all environment features")
		// Use scopeBatch with environmentIDAll for batch operations covering all environments
		recordListFeatures(scopeBatch, environmentIDAll, codeFail, time.Since(startTime).Seconds())
		return err
	}
	// Use scopeBatch with environmentIDAll for batch operations covering all environments
	recordListFeatures(scopeBatch, environmentIDAll, codeSuccess, time.Since(startTime).Seconds())

	for _, envFt := range envFts {
		filtered := c.removeOldFeatures(envFt.Features)
		fts := &ftproto.Features{
			Id:       evaluation.GenerateFeaturesID(filtered),
			Features: filtered,
		}
		c.putCache(fts, envFt.EnvironmentId, len(filtered))
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
func (c *featureFlagCacher) putCache(features *ftproto.Features, environmentID string, featureCount int) {
	var wg sync.WaitGroup
	var hasError bool
	var mu sync.Mutex

	for _, cache := range c.caches {
		wg.Add(1)
		go func(cache cachev3.FeaturesCache) {
			defer wg.Done()
			if err := cache.Put(features, environmentID); err != nil {
				c.logger.Error("Failed to cache features",
					zap.Error(err),
					zap.String("environmentId", environmentID),
				)
				mu.Lock()
				hasError = true
				mu.Unlock()
			}
		}(cache)
	}
	wg.Wait()

	// Record metrics based on overall success/failure
	if hasError {
		recordCachePut(environmentID, codeFail)
	} else {
		recordCachePut(environmentID, codeSuccess)
		recordFeaturesUpdated(environmentID, featureCount)
	}
}
