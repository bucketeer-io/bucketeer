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

type segmentUserCacher struct {
	environmentClient envclient.Client
	featureClient     ftclient.Client
	cache             cachev3.SegmentUsersCache
	opts              *jobs.Options
	logger            *zap.Logger
}

func NewSegmentUserCacher(
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
	return &segmentUserCacher{
		environmentClient: environmentClient,
		featureClient:     featureClient,
		cache:             cachev3.NewSegmentUsersCache(cache),
		opts:              dopts,
		logger:            dopts.Logger.Named("segment-user-cacher"),
	}
}

func (c *segmentUserCacher) Run(ctx context.Context) error {
	envs, err := c.listAllEnvironments(ctx)
	if err != nil {
		c.logger.Error("Failed to list all environments")
		return err
	}
	for _, env := range envs {
		users, err := c.listSegmentUsers(ctx, env.Id)
		if err != nil {
			c.logger.Error("Failed to list segment users", zap.String("environmentId", env.Id))
			return err
		}
		if err := c.cache.Put(&ftproto.SegmentUsers{Users: users}, env.Id); err != nil {
			c.logger.Error("Failed to cache segment users", zap.String("environmentId", env.Id))
			continue
		}
	}
	return nil
}

func (c *segmentUserCacher) listAllEnvironments(
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

func (c *segmentUserCacher) listSegmentUsers(
	ctx context.Context,
	environmentID string,
) ([]*ftproto.SegmentUser, error) {
	req := &ftproto.ListSegmentUsersRequest{
		PageSize:             0,
		EnvironmentNamespace: environmentID,
	}
	resp, err := c.featureClient.ListSegmentUsers(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Users, nil
}
