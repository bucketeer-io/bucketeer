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
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	envclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	expclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	expproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

const (
	day = 24 * time.Hour
)

type experimentCacher struct {
	environmentClient envclient.Client
	experimentClient  expclient.Client
	caches            []cachev3.ExperimentsCache
	opts              *jobs.Options
	logger            *zap.Logger
}

func NewExperimentCacher(
	environmentClient envclient.Client,
	experimentClient expclient.Client,
	multiCaches []cache.MultiGetCache,
	opts ...jobs.Option,
) jobs.Job {
	dopts := &jobs.Options{
		Logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	caches := make([]cachev3.ExperimentsCache, 0, len(multiCaches))
	for _, cache := range multiCaches {
		caches = append(caches, cachev3.NewExperimentsCache(cache))
	}
	return &experimentCacher{
		environmentClient: environmentClient,
		experimentClient:  experimentClient,
		caches:            caches,
		opts:              dopts,
		logger:            dopts.Logger.Named("experiment-cacher"),
	}
}

func (c *experimentCacher) Run(ctx context.Context) (lastErr error) {
	startTime := time.Now()
	defer func() {
		jobs.RecordJob(jobs.JobExperimentCacher, lastErr, time.Since(startTime))
	}()
	envs, err := c.listAllEnvironments(ctx)
	if err != nil {
		c.logger.Error("Failed to list all environments")
		return err
	}
	for _, env := range envs {
		experiments, err := c.listExperiments(ctx, env.Id)
		if err != nil {
			c.logger.Error("Failed to list experiments", zap.String("environmentId", env.Id))
			return err
		}
		c.putCache(&expproto.Experiments{Experiments: experiments}, env.Id)
	}
	return nil
}

func (c *experimentCacher) listAllEnvironments(
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

// List only running, stopped, and not archived experiments
func (c *experimentCacher) listExperiments(
	ctx context.Context,
	environmentID string,
) ([]*expproto.Experiment, error) {
	req := &expproto.ListExperimentsRequest{
		// Because the evaluation and goal events may be sent with a delay
		// for many reasons from the client side, we still calculate
		// the results for two days after it stopped.
		StopAt:        time.Now().Add(-2 * day).Unix(),
		PageSize:      0,
		EnvironmentId: environmentID,
		Statuses: []expproto.Experiment_Status{
			expproto.Experiment_RUNNING,
			expproto.Experiment_STOPPED,
		},
		Archived: &wrapperspb.BoolValue{Value: false},
	}
	resp, err := c.experimentClient.ListExperiments(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Experiments, nil
}

// Save the experiments by environment in all redis instances
// Since the batch runs every minute, we don't handle erros when putting the cache
func (c *experimentCacher) putCache(experiments *expproto.Experiments, environmentID string) int {
	var updatedInstances int
	var mu sync.Mutex     // Mutex to safely update `updatedInstances` across goroutines
	var wg sync.WaitGroup // Use a WaitGroup to wait for all goroutines to finish
	for _, cache := range c.caches {
		wg.Add(1) // Increment the WaitGroup counter
		go func(cache cachev3.ExperimentsCache) {
			defer wg.Done()
			if err := cache.Put(experiments, environmentID); err != nil {
				// Log the error, but do not stop the other goroutines
				c.logger.Error("Failed to cache experiments",
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
