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
	"time"

	"go.uber.org/zap"

	accstorage "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	accproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

type apiKeyCacher struct {
	accStorage accstorage.AccountStorage
	caches     []cachev3.EnvironmentAPIKeyCache
	opts       *jobs.Options
	logger     *zap.Logger
}

func NewAPIKeyCacher(
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
	caches := make([]cachev3.EnvironmentAPIKeyCache, 0, len(multiCaches))
	for _, cache := range multiCaches {
		caches = append(caches, cachev3.NewEnvironmentAPIKeyCache(cache))
	}
	return &apiKeyCacher{
		accStorage: accstorage.NewAccountStorage(mysqlClient),
		caches:     caches,
		opts:       dopts,
		logger:     dopts.Logger.Named("api-key-cacher"),
	}
}

func (c *apiKeyCacher) Run(ctx context.Context) (lastErr error) {
	startTime := time.Now()
	defer func() {
		jobs.RecordJob(jobs.JobAPIKeyCacher, lastErr, time.Since(startTime))
	}()
	envAPIKeys, err := c.accStorage.ListAllEnvironmentAPIKeys(ctx)
	if err != nil {
		return err
	}
	for _, envAPIKey := range envAPIKeys {
		updatedInstances := c.putCache(envAPIKey.EnvironmentAPIKey)
		c.logger.Debug("Caching environment api key",
			zap.Any("environmentApiKey", envAPIKey.EnvironmentAPIKey),
			zap.Int("updatedInstances", updatedInstances),
		)
	}
	return nil
}

// Save the environment API key in all redis instances
// Since the batch runs every minute, we don't handle erros when putting the cache
func (c *apiKeyCacher) putCache(envAPIKey *accproto.EnvironmentAPIKey) int {
	var updatedInstances int
	var mu sync.Mutex     // Mutex to safely update `updatedInstances` across goroutines
	var wg sync.WaitGroup // Use a WaitGroup to wait for all goroutines to finish
	for _, cache := range c.caches {
		wg.Add(1) // Increment the WaitGroup counter
		go func(cache cachev3.EnvironmentAPIKeyCache) {
			defer wg.Done()
			if err := cache.Put(envAPIKey); err != nil {
				// Log the error, but do not stop the other goroutines
				c.logger.Error("Failed to cache environment api key",
					zap.Error(err),
					zap.Any("envAPIKey", envAPIKey),
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
