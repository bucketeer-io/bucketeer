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

	aoclient "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	envclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	aoproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

type autoOpsRulesCacher struct {
	environmentClient envclient.Client
	autoOpsClient     aoclient.Client
	caches            []cachev3.AutoOpsRulesCache
	opts              *jobs.Options
	logger            *zap.Logger
}

func NewAutoOpsRulesCacher(
	environmentClient envclient.Client,
	autoOpsClient aoclient.Client,
	multiCaches []cache.MultiGetCache,
	opts ...jobs.Option,
) jobs.Job {
	dopts := &jobs.Options{
		Logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	caches := make([]cachev3.AutoOpsRulesCache, 0, len(multiCaches))
	for _, cache := range multiCaches {
		caches = append(caches, cachev3.NewAutoOpsRulesCache(cache))
	}
	return &autoOpsRulesCacher{
		environmentClient: environmentClient,
		autoOpsClient:     autoOpsClient,
		caches:            caches,
		opts:              dopts,
		logger:            dopts.Logger.Named("auto-ops-rules-cacher"),
	}
}

func (c *autoOpsRulesCacher) Run(ctx context.Context) (lastErr error) {
	startTime := time.Now()
	defer func() {
		jobs.RecordJob(jobs.JobAutoOpsRulesCacher, lastErr, time.Since(startTime))
	}()
	envs, err := c.listAllEnvironments(ctx)
	if err != nil {
		c.logger.Error("Failed to list all environments")
		return err
	}
	for _, env := range envs {
		autoOpsRules, err := c.listAutoOpsRules(ctx, env.Id)
		if err != nil {
			c.logger.Error("Failed to list auto ops rules", zap.String("environmentId", env.Id))
			return err
		}
		c.putCache(&aoproto.AutoOpsRules{AutoOpsRules: autoOpsRules}, env.Id)
	}
	return nil
}

func (c *autoOpsRulesCacher) listAllEnvironments(
	ctx context.Context,
) ([]*envproto.EnvironmentV2, error) {
	req := &envproto.ListEnvironmentsV2Request{
		PageSize: 0,
	}
	resp, err := c.environmentClient.ListEnvironmentsV2(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Environments, nil
}

func (c *autoOpsRulesCacher) listAutoOpsRules(
	ctx context.Context,
	environmentID string,
) ([]*aoproto.AutoOpsRule, error) {
	req := &aoproto.ListAutoOpsRulesRequest{
		PageSize:      0,
		EnvironmentId: environmentID,
	}
	resp, err := c.autoOpsClient.ListAutoOpsRules(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.AutoOpsRules, nil
}

// Save the auto ops rules by environment in all redis instances
func (c *autoOpsRulesCacher) putCache(rules *aoproto.AutoOpsRules, environmentID string) {
	var wg sync.WaitGroup
	for _, cache := range c.caches {
		wg.Add(1)
		go func(cache cachev3.AutoOpsRulesCache) {
			defer wg.Done()
			if err := cache.Put(rules, environmentID); err != nil {
				c.logger.Error("Failed to cache auto ops rules",
					zap.Error(err),
					zap.String("environmentId", environmentID),
				)
			}
		}(cache)
	}
	wg.Wait()
}
