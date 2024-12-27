// Copyright 2025 The Bucketeer Authors.
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
//

package cacher

import (
	"context"
	"sync"

	"go.uber.org/zap"

	evaluation "github.com/bucketeer-io/bucketeer/evaluation/go"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	ftstorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	ftproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type featureFlagCacher struct {
	ftStorage ftstorage.FeatureStorage
	caches    []cachev3.FeaturesCache
	opts      *jobs.Options
	logger    *zap.Logger
}

func NewFeatureFlagCacher(
	mysqlClient mysql.Client,
	multiCaches []cache.MultiGetCache,
	opts ...jobs.Option,
) jobs.Job {
	dopts := &jobs.Options{
		Logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	caches := make([]cachev3.FeaturesCache, 0, len(multiCaches))
	for _, cache := range multiCaches {
		caches = append(caches, cachev3.NewFeaturesCache(cache))
	}
	return &featureFlagCacher{
		ftStorage: ftstorage.NewFeatureStorage(mysqlClient),
		caches:    caches,
		opts:      dopts,
		logger:    dopts.Logger.Named("feature-flag-cacher"),
	}
}

func (c *featureFlagCacher) Run(ctx context.Context) error {
	envFts, err := c.ftStorage.ListAllEnvironmentFeatures(ctx)
	if err != nil {
		c.logger.Error("Failed to all environment features")
		return err
	}
	for _, envFt := range envFts {
		fts := &ftproto.Features{
			Id:       evaluation.GenerateFeaturesID(envFt.Features),
			Features: envFt.Features,
		}
		updatedInstances := c.putCache(fts, envFt.EnvironmentId)
		c.logger.Debug("Updated Redis instances", zap.Int("size", updatedInstances))
		c.logger.Debug("Caching features",
			zap.String("environmentId", envFt.EnvironmentId),
			zap.Int("featureCount", len(envFt.Features)),
		)
	}
	return nil
}

// Save the flags by environment in all redis instances
// Since the batch runs every minute, we don't handle erros when putting the cache
func (c *featureFlagCacher) putCache(features *ftproto.Features, environmentID string) int {
	var updatedInstances int
	var mu sync.Mutex     // Mutex to safely update `updatedInstances` across goroutines
	var wg sync.WaitGroup // Use a WaitGroup to wait for all goroutines to finish
	for _, cache := range c.caches {
		wg.Add(1) // Increment the WaitGroup counter
		go func(cache cachev3.FeaturesCache) {
			defer wg.Done()
			if err := cache.Put(features, environmentID); err != nil {
				// Log the error, but do not stop the other goroutines
				c.logger.Error("Failed to cache features",
					zap.Error(err),
					zap.String("environmentId", environmentID),
				)
				return
			}
			mu.Lock()
			updatedInstances++
			mu.Unlock()
		}(cache)
	}
	wg.Wait()
	return updatedInstances
}
