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

	evaluation "github.com/bucketeer-io/bucketeer/evaluation/go"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	envclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	ftclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
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
		Logger: zap.NewNop(),
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
		fts := &ftproto.Features{
			Id:       evaluation.GenerateFeaturesID(features),
			Features: features,
		}
		fids := make([]string, 0, len(fts.Features))
		for _, f := range fts.Features {
			fids = append(fids, f.Id)
		}
		c.logger.Info("Caching features",
			zap.String("environmentId", env.Id),
			zap.Strings("featureIds", fids),
		)
		if err := c.cache.Put(fts, env.Id); err != nil {
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
		PageSize:      0,
		EnvironmentId: environmentID,
	}
	resp, err := c.featureClient.ListFeatures(ctx, req)
	if err != nil {
		return nil, err
	}
	filtered := make([]*ftproto.Feature, 0, len(resp.Features))
	for _, f := range resp.Features {
		ff := domain.Feature{Feature: f}
		if ff.IsDisabledAndOffVariationEmpty() {
			continue
		}
		// We exclude archived feature flags over thirty days ago to keep the cache size small.
		if ff.IsArchivedBeforeLastThirtyDays() {
			continue
		}
		filtered = append(filtered, f)
	}
	return filtered, nil
}
