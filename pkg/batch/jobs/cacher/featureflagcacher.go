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

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	envclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	ftclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	envproto "github.com/bucketeer-io/bucketeer/proto/environment"
	ftproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type featureFlagCacher struct {
	environmentClient envclient.Client
	featureClient     ftclient.Client
	cache             cachev3.FeaturesCache
	opts              *jobs.Options
	logger            *zap.Logger
}

func NewFeatureFlagCacher(
	environmentClient envclient.Client,
	featureClient ftclient.Client,
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
	return &featureFlagCacher{
		environmentClient: environmentClient,
		featureClient:     featureClient,
		cache:             cachev3.NewFeaturesCache(cache),
		opts:              dopts,
		logger:            dopts.Logger.Named("feature-flag-cacher"),
	}
}

func (c *featureFlagCacher) Run(ctx context.Context) error {
	envs, err := c.listAllEnvironments(ctx)
	if err != nil {
		c.logger.Error("Failed to list all environments")
		return err
	}
	for _, env := range envs {
		features, err := c.listFeatures(ctx, env.Id)
		if err != nil {
			c.logger.Error("Failed to list features", zap.String("environmentId", env.Id))
			return err
		}
		if err := c.cache.Put(&ftproto.Features{Features: features}, env.Id); err != nil {
			c.logger.Error("Failed to cache features", zap.String("environmentId", env.Id))
			continue
		}
	}
	return nil
}

func (c *featureFlagCacher) listAllEnvironments(
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

func (c *featureFlagCacher) listFeatures(
	ctx context.Context,
	environmentID string,
) ([]*ftproto.Feature, error) {
	req := &ftproto.ListFeaturesRequest{
		PageSize:             0,
		EnvironmentNamespace: environmentID,
	}
	resp, err := c.featureClient.ListFeatures(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Features, nil
}
