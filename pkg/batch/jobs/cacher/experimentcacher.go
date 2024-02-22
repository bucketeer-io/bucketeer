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
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	envclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	expclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	envproto "github.com/bucketeer-io/bucketeer/proto/environment"
	expproto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

type experimentCacher struct {
	environmentClient envclient.Client
	experimentClient  expclient.Client
	cache             cachev3.ExperimentsCache
	opts              *jobs.Options
	logger            *zap.Logger
}

func NewExperimentCacher(
	environmentClient envclient.Client,
	experimentClient expclient.Client,
	cache cache.MultiGetCache,
	opts ...jobs.Option,
) jobs.Job {
	dopts := &jobs.Options{
		Timeout: 1 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &experimentCacher{
		environmentClient: environmentClient,
		experimentClient:  experimentClient,
		cache:             cachev3.NewExperimentsCache(cache),
		opts:              dopts,
		logger:            dopts.Logger.Named("experiment-cacher"),
	}
}

func (c *experimentCacher) Run(ctx context.Context) error {
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
		if err := c.cache.Put(&expproto.Experiments{Experiments: experiments}, env.Id); err != nil {
			c.logger.Error("Failed to cache experiments", zap.String("environmentId", env.Id))
			continue
		}
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
		PageSize:             0,
		EnvironmentNamespace: environmentID,
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
