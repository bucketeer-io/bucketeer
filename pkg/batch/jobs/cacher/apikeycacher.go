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
	"google.golang.org/protobuf/types/known/wrapperspb"

	acclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	envclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	acproto "github.com/bucketeer-io/bucketeer/proto/account"
	envproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

type apiKeyCacher struct {
	environmentClient envclient.Client
	accountClient     acclient.Client
	cache             cachev3.EnvironmentAPIKeyCache
	opts              *jobs.Options
	logger            *zap.Logger
}

func NewAPIKeyCacher(
	environmentClient envclient.Client,
	accountClient acclient.Client,
	cache cache.MultiGetCache,
	opts ...jobs.Option,
) jobs.Job {
	dopts := &jobs.Options{
		Logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &apiKeyCacher{
		environmentClient: environmentClient,
		accountClient:     accountClient,
		cache:             cachev3.NewEnvironmentAPIKeyCache(cache),
		opts:              dopts,
		logger:            dopts.Logger.Named("api-key-cacher"),
	}
}

func (c *apiKeyCacher) Run(ctx context.Context) error {
	envs, err := c.listAllEnvironments(ctx)
	if err != nil {
		c.logger.Error("Failed to list all environments", zap.Error(err))
		return err
	}
	for _, env := range envs {
		envAPIKeys, err := c.listEnvAPIKeys(ctx, env)
		if err != nil {
			c.logger.Error("Failed to list environment api keys",
				zap.Error(err),
				zap.String("environmentId", env.Id),
			)
			return err
		}
		for _, envAPIKey := range envAPIKeys {
			if err := c.cache.Put(envAPIKey); err != nil {
				c.logger.Error("Failed to cache environment api key",
					zap.Error(err),
					zap.Any("envAPIKey", envAPIKey),
				)
				continue
			}
		}
	}
	return nil
}

func (c *apiKeyCacher) listAllEnvironments(
	ctx context.Context,
) ([]*envproto.EnvironmentV2, error) {
	req := &envproto.ListEnvironmentsV2Request{
		PageSize: 0,
		Archived: wrapperspb.Bool(false),
	}
	resp, err := c.environmentClient.ListEnvironmentsV2(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Environments, nil
}

func (c *apiKeyCacher) listEnvAPIKeys(
	ctx context.Context,
	environment *envproto.EnvironmentV2,
) ([]*acproto.EnvironmentAPIKey, error) {
	req := &acproto.ListAPIKeysRequest{
		PageSize:       0,
		EnvironmentIds: []string{environment.Id},
	}
	resp, err := c.accountClient.ListAPIKeys(ctx, req)
	if err != nil {
		return nil, err
	}
	proj, err := c.getProject(ctx, environment.ProjectId)
	if err != nil {
		c.logger.Error("Failed to get project",
			zap.Error(err),
			zap.String("organizationId", environment.OrganizationId),
			zap.String("projectId", environment.ProjectId),
		)
		return nil, err
	}
	envAPIKeys := make([]*acproto.EnvironmentAPIKey, 0, len(resp.ApiKeys))
	for _, key := range resp.ApiKeys {
		envAPIKey := &acproto.EnvironmentAPIKey{
			ApiKey:              key,
			EnvironmentDisabled: proj.Disabled,
			ProjectId:           environment.ProjectId,
			ProjectUrlCode:      proj.UrlCode,
			Environment:         environment,
		}
		envAPIKeys = append(envAPIKeys, envAPIKey)
	}
	return envAPIKeys, nil
}

func (c *apiKeyCacher) getProject(
	ctx context.Context,
	projectID string,
) (*envproto.Project, error) {
	req := &envproto.GetProjectRequest{
		Id: projectID,
	}
	resp, err := c.environmentClient.GetProject(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Project, nil
}
