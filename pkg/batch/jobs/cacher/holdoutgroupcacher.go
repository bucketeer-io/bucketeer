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
//

package cacher

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	envclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	v2hs "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

type holdoutGroupCacher struct {
	environmentClient    envclient.Client
	storage              v2hs.HoldoutStorage
	caches               []cachev3.HoldoutGroupCache
	measuredGroupsCaches []cachev3.HoldoutMeasuredGroupsCache
	opts                 *jobs.Options
	logger               *zap.Logger
}

func NewHoldoutGroupCacher(
	environmentClient envclient.Client,
	storage v2hs.HoldoutStorage,
	multiCaches []cache.MultiGetCache,
	opts ...jobs.Option,
) jobs.Job {
	dopts := &jobs.Options{
		Logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	caches := make([]cachev3.HoldoutGroupCache, 0, len(multiCaches))
	measuredGroupsCaches := make([]cachev3.HoldoutMeasuredGroupsCache, 0, len(multiCaches))
	for _, multiCache := range multiCaches {
		caches = append(caches, cachev3.NewHoldoutGroupCache(multiCache, 0))
		measuredGroupsCaches = append(measuredGroupsCaches, cachev3.NewHoldoutMeasuredGroupsCache(multiCache))
	}
	return &holdoutGroupCacher{
		environmentClient:    environmentClient,
		storage:              storage,
		caches:               caches,
		measuredGroupsCaches: measuredGroupsCaches,
		opts:                 dopts,
		logger:               dopts.Logger.Named("holdout-group-cacher"),
	}
}

func (c *holdoutGroupCacher) Run(ctx context.Context) (lastErr error) {
	startTime := time.Now()
	defer func() {
		jobs.RecordJob(jobs.JobHoldoutGroupCacher, lastErr, time.Since(startTime))
	}()
	envs, err := c.listAllEnvironments(ctx)
	if err != nil {
		c.logger.Error("Failed to list all environments", zap.Error(err))
		return err
	}
	// One environment's failure must not leave the remaining environments
	// with a stale snapshot until the next run. Refresh every environment
	// and report the last error afterwards.
	for _, env := range envs {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		// The two caches are refreshed independently: a failure of one must
		// not leave the other stale.
		if err := c.refreshServingSnapshot(ctx, env.Id); err != nil {
			lastErr = err
		}
		if err := c.refreshMeasuredGroups(ctx, env.Id); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func (c *holdoutGroupCacher) refreshServingSnapshot(ctx context.Context, environmentID string) error {
	// Captured before the read so a slower writer cannot overwrite a newer snapshot.
	generation := time.Now().UnixNano()
	snapshot, err := v2hs.ServingSnapshot(ctx, c.storage, environmentID)
	if err != nil {
		c.logger.Error("Failed to read the live holdout group",
			zap.Error(err),
			zap.String("environmentId", environmentID),
		)
		return err
	}
	c.putCache(snapshot, environmentID, generation)
	return nil
}

// refreshMeasuredGroups stores the groups the analytics pipeline measures.
// The list is empty for an environment without such a group, which is every
// environment whose groups have no goals.
func (c *holdoutGroupCacher) refreshMeasuredGroups(ctx context.Context, environmentID string) error {
	// Captured before the read so a slower writer cannot overwrite a newer list.
	now := time.Now()
	groups, err := v2hs.MeasuredGroups(ctx, c.storage, environmentID, now)
	if err != nil {
		c.logger.Error("Failed to read the measured holdout groups",
			zap.Error(err),
			zap.String("environmentId", environmentID),
		)
		return err
	}
	c.putMeasuredGroups(groups, environmentID, now.UnixNano())
	return nil
}

func (c *holdoutGroupCacher) listAllEnvironments(
	ctx context.Context,
) ([]*envproto.EnvironmentV2, error) {
	resp, err := c.environmentClient.ListEnvironmentsV2(ctx, &envproto.ListEnvironmentsV2Request{
		PageSize: 0,
	})
	if err != nil {
		return nil, err
	}
	return resp.Environments, nil
}

func (c *holdoutGroupCacher) putCache(
	group *ftproto.HoldoutGroup,
	environmentID string,
	generation int64,
) int {
	var updatedInstances int
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, cache := range c.caches {
		wg.Add(1)
		go func(cache cachev3.HoldoutGroupCache) {
			defer wg.Done()
			accepted, err := cache.PutIfNewer(group, environmentID, generation)
			if err != nil {
				c.logger.Error("Failed to cache holdout group",
					zap.Error(err),
					zap.String("environmentId", environmentID),
				)
				return
			}
			if !accepted {
				c.logger.Info("Skipped holdout cache write: a newer snapshot is already cached",
					zap.String("environmentId", environmentID),
					zap.Int64("generation", generation),
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

func (c *holdoutGroupCacher) putMeasuredGroups(
	groups *ftproto.HoldoutGroups,
	environmentID string,
	generation int64,
) {
	var wg sync.WaitGroup
	for _, cache := range c.measuredGroupsCaches {
		wg.Add(1)
		go func(cache cachev3.HoldoutMeasuredGroupsCache) {
			defer wg.Done()
			accepted, err := cache.PutIfNewer(groups, environmentID, generation)
			if err != nil {
				c.logger.Error("Failed to cache measured holdout groups",
					zap.Error(err),
					zap.String("environmentId", environmentID),
				)
				return
			}
			if !accepted {
				c.logger.Info("Skipped measured holdout groups cache write: a newer list is already cached",
					zap.String("environmentId", environmentID),
					zap.Int64("generation", generation),
				)
			}
		}(cache)
	}
	wg.Wait()
}
