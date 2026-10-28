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

package holdout

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	environmentclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/holdout/domain"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	environmentListPageSize int64 = 500
	// ListHoldoutGroups rejects a page size above 100.
	holdoutListPageSize int64 = 100
)

type statusUpdater struct {
	environmentClient environmentclient.Client
	featureClient     featureclient.Client
	opts              *jobs.Options
	logger            *zap.Logger
}

func NewHoldoutGroupStatusUpdater(
	environmentClient environmentclient.Client,
	featureClient featureclient.Client,
	opts ...jobs.Option,
) jobs.Job {
	dopts := &jobs.Options{
		Timeout: 1 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &statusUpdater{
		environmentClient: environmentClient,
		featureClient:     featureClient,
		opts:              dopts,
		logger:            dopts.Logger.Named("holdout-status-updater"),
	}
}

func (u *statusUpdater) Run(ctx context.Context) (lastErr error) {
	startTime := time.Now()
	defer func() {
		jobs.RecordJob(jobs.JobHoldoutGroupStatusUpdater, lastErr, time.Since(startTime))
	}()

	ctx, cancel := context.WithTimeout(ctx, u.opts.Timeout)
	defer cancel()
	environments, err := u.listEnvironments(ctx)
	if err != nil {
		u.logger.Error("Failed to list environments", zap.Error(err))
		return err
	}
	for _, env := range environments {
		// The legacy environment has an empty ID. The holdout API requires an
		// environment ID on every request, so no group can exist there; listing
		// it would only fail the whole run with InvalidArgument.
		if env.Id == "" {
			continue
		}
		groups, err := u.listGroups(ctx, env.Id)
		if err != nil {
			u.logger.Error("Failed to list holdout groups",
				zap.Error(err),
				zap.String("environmentId", env.Id),
			)
			lastErr = err
			continue
		}
		// Stops run before starts. List order is by name, so a due waiting group
		// can appear ahead of the expired group that still occupies its period.
		// Stopping that predecessor in this run lets the successor start now.
		now := time.Now()
		var stops, rest []*featureproto.HoldoutGroup
		for _, group := range groups {
			if domain.DueForScheduledStop(group, now) {
				stops = append(stops, group)
				continue
			}
			rest = append(rest, group)
		}
		for _, batch := range [][]*featureproto.HoldoutGroup{stops, rest} {
			for _, group := range batch {
				if err = u.updateStatus(ctx, env.Id, group); err != nil {
					lastErr = err
				}
			}
		}
	}
	return lastErr
}

// updateStatus advances a group by at most one transition per run.
func (u *statusUpdater) updateStatus(
	ctx context.Context,
	environmentID string,
	group *featureproto.HoldoutGroup,
) error {
	now := time.Now()
	switch {
	case domain.DueForScheduledStart(group, now):
		return u.start(ctx, environmentID, group)
	case domain.DueForScheduledStop(group, now):
		return u.stop(ctx, environmentID, group)
	case domain.DueForScheduledAnalysis(group, now):
		return u.startAnalysis(ctx, environmentID, group)
	default:
		return nil
	}
}

func (u *statusUpdater) start(
	ctx context.Context,
	environmentID string,
	group *featureproto.HoldoutGroup,
) error {
	_, err := u.featureClient.StartHoldoutGroup(ctx, &featureproto.StartHoldoutGroupRequest{
		Id:            group.Id,
		EnvironmentId: environmentID,
		Mode:          featureproto.StartHoldoutGroupRequest_SCHEDULED,
	})
	if err != nil {
		u.logger.Error("Failed to start holdout group",
			zap.Error(err),
			zap.String("environmentId", environmentID),
			zap.String("id", group.Id),
		)
	}
	return err
}

func (u *statusUpdater) startAnalysis(
	ctx context.Context,
	environmentID string,
	group *featureproto.HoldoutGroup,
) error {
	_, err := u.featureClient.StartHoldoutGroupAnalysis(ctx, &featureproto.StartHoldoutGroupAnalysisRequest{
		Id:            group.Id,
		EnvironmentId: environmentID,
		Mode:          featureproto.StartHoldoutGroupAnalysisRequest_SCHEDULED,
	})
	if err != nil {
		u.logger.Error("Failed to start holdout group analysis",
			zap.Error(err),
			zap.String("environmentId", environmentID),
			zap.String("id", group.Id),
		)
	}
	return err
}

func (u *statusUpdater) stop(
	ctx context.Context,
	environmentID string,
	group *featureproto.HoldoutGroup,
) error {
	_, err := u.featureClient.StopHoldoutGroup(ctx, &featureproto.StopHoldoutGroupRequest{
		Id:            group.Id,
		EnvironmentId: environmentID,
		Mode:          featureproto.StopHoldoutGroupRequest_SCHEDULED,
	})
	if err != nil {
		u.logger.Error("Failed to stop holdout group",
			zap.Error(err),
			zap.String("environmentId", environmentID),
			zap.String("id", group.Id),
		)
	}
	return err
}

func (u *statusUpdater) listGroups(
	ctx context.Context,
	environmentID string,
) ([]*featureproto.HoldoutGroup, error) {
	var groups []*featureproto.HoldoutGroup
	cursor := ""
	for {
		resp, err := u.featureClient.ListHoldoutGroups(ctx, &featureproto.ListHoldoutGroupsRequest{
			PageSize:      holdoutListPageSize,
			Cursor:        cursor,
			EnvironmentId: environmentID,
			Statuses: []featureproto.HoldoutGroup_Status{
				featureproto.HoldoutGroup_WAITING,
				featureproto.HoldoutGroup_RUNNING,
				featureproto.HoldoutGroup_ANALYSIS,
			},
		})
		if err != nil {
			return nil, err
		}
		groups = append(groups, resp.HoldoutGroups...)
		if len(resp.HoldoutGroups) == 0 || int64(len(resp.HoldoutGroups)) < holdoutListPageSize {
			return groups, nil
		}
		cursor = resp.Cursor
	}
}

func (u *statusUpdater) listEnvironments(ctx context.Context) ([]*environmentproto.EnvironmentV2, error) {
	var environments []*environmentproto.EnvironmentV2
	cursor := ""
	for {
		resp, err := u.environmentClient.ListEnvironmentsV2(ctx, &environmentproto.ListEnvironmentsV2Request{
			PageSize: environmentListPageSize,
			Cursor:   cursor,
			Archived: &wrapperspb.BoolValue{Value: false},
		})
		if err != nil {
			return nil, err
		}
		environments = append(environments, resp.Environments...)
		if len(resp.Environments) == 0 || int64(len(resp.Environments)) < environmentListPageSize {
			return environments, nil
		}
		cursor = resp.Cursor
	}
}
