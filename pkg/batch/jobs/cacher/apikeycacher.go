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
//

package cacher

import (
	"context"

	"go.uber.org/zap"

	accstorage "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
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

func (c *apiKeyCacher) Run(ctx context.Context) error {
	envAPIKeys, err := c.accStorage.ListAllEnvironmentAPIKeys(ctx)
	if err != nil {
		return err
	}
	c.logger.Debug("Caching environment api keys",
		zap.Any("environmentApiKeys", envAPIKeys),
	)
	for _, envAPIKey := range envAPIKeys {
		for _, cache := range c.caches {
			if err := cache.Put(envAPIKey.EnvironmentAPIKey); err != nil {
				c.logger.Error("Failed to cache environment api key",
					zap.Error(err),
					zap.Any("envAPIKey", envAPIKey),
				)
			}
			continue
		}
	}
	return nil
}
