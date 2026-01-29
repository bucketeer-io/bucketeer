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

package opsevent

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	aoclient "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client"
	autoopsdomain "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	envclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	ftcacher "github.com/bucketeer-io/bucketeer/v2/pkg/feature/cacher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/batch/executor"
	aoproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

type progressiveRolloutWatcher struct {
	envClient                  envclient.Client
	aoClient                   aoclient.Client
	progressiveRolloutExecutor executor.ProgressiveRolloutExecutor
	ftCacher                   ftcacher.FeatureFlagCacher
	opts                       *jobs.Options
	logger                     *zap.Logger
}

func NewProgressiveRolloutWacher(
	envClient envclient.Client,
	aoClient aoclient.Client,
	progressiveRolloutExecutor executor.ProgressiveRolloutExecutor,
	ftCacher ftcacher.FeatureFlagCacher,
	opts ...jobs.Option,
) jobs.Job {
	dopts := &jobs.Options{
		Timeout: 5 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &progressiveRolloutWatcher{
		envClient:                  envClient,
		aoClient:                   aoClient,
		progressiveRolloutExecutor: progressiveRolloutExecutor,
		ftCacher:                   ftCacher,
		opts:                       dopts,
		logger:                     dopts.Logger.Named("progressive-rollout-watcher"),
	}
}

func (w *progressiveRolloutWatcher) Run(
	ctx context.Context,
) (lastErr error) {
	ctx, cancel := context.WithTimeout(ctx, w.opts.Timeout)
	defer cancel()
	environments, err := w.listEnvironments(ctx)
	if err != nil {
		lastErr = err
		return
	}
	for _, e := range environments {
		progressiveRollouts, err := w.listProgressiveRollouts(ctx, e.Id)
		if err != nil {
			lastErr = err
			return
		}
		var executed bool
		for _, p := range progressiveRollouts {
			progressiveRollout := &autoopsdomain.ProgressiveRollout{ProgressiveRollout: p}
			// Skip finished or stopped progressive rollouts to avoid unnecessary processing
			// This is consistent with datetime_watcher and event_count_watcher behavior
			if progressiveRollout.IsFinished() || progressiveRollout.IsStopped() {
				continue
			}
			wasExecuted, err := w.executeProgressiveRollout(ctx, p, e.Id)
			if err != nil {
				lastErr = err
			}
			if wasExecuted {
				executed = true
			}
		}
		// Update Redis cache immediately after successful progressive rollout execution
		// This ensures SDKs receive the updated flags without waiting for the periodic cache refresh
		if executed && w.ftCacher != nil {
			if err := w.ftCacher.RefreshEnvironmentCache(ctx, e.Id); err != nil {
				w.logger.Error("Failed to update feature flag cache after progressive rollout execution",
					zap.Error(err),
					zap.String("environmentId", e.Id),
				)
				// Don't set lastErr here - the progressive rollout execution succeeded,
				// cache update failure shouldn't be reported as main job failure
			}
		}
	}
	return lastErr
}

func (w *progressiveRolloutWatcher) listEnvironments(ctx context.Context) ([]*envproto.EnvironmentV2, error) {
	resp, err := w.envClient.ListEnvironmentsV2(ctx, &envproto.ListEnvironmentsV2Request{
		PageSize: 0,
		Archived: wrapperspb.Bool(false),
	})
	if err != nil {
		return nil, err
	}
	return resp.Environments, nil
}

func (w *progressiveRolloutWatcher) listProgressiveRollouts(
	ctx context.Context,
	environmentID string,
) ([]*aoproto.ProgressiveRollout, error) {
	resp, err := w.aoClient.ListProgressiveRollouts(
		ctx,
		&aoproto.ListProgressiveRolloutsRequest{
			EnvironmentId: environmentID,
			PageSize:      0,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp.ProgressiveRollouts, nil
}

// executeProgressiveRollout returns (executed bool, err error)
// // executed is true if the progressive rollout operation was actually executed
// (i.e. ExecuteProgressiveRollout was called and succeeded)
func (w *progressiveRolloutWatcher) executeProgressiveRollout(
	ctx context.Context,
	progressiveRollout *aoproto.ProgressiveRollout,
	environmentId string,
) (bool, error) {
	pr := &autoopsdomain.ProgressiveRollout{ProgressiveRollout: progressiveRollout}
	schedules, err := pr.ExtractSchedules()
	if err != nil {
		w.logger.Error("Failed to extract schedules", zap.Error(err),
			zap.String("environmentId", environmentId),
			zap.String("featureId", progressiveRollout.FeatureId),
			zap.String("progressiveRolloutId", progressiveRollout.Id),
		)
		return false, err
	}
	now := time.Now().Unix()
	for _, s := range schedules {
		if s.TriggeredAt == 0 && s.ExecuteAt <= now {
			w.logger.Info("Executing scheduled progressive rollout",
				zap.String("environmentId", environmentId),
				zap.String("featureId", progressiveRollout.FeatureId),
				zap.String("progressiveRolloutId", progressiveRollout.Id),
				zap.String("scheduleId", s.ScheduleId),
			)
			if err := w.progressiveRolloutExecutor.ExecuteProgressiveRollout(
				ctx,
				environmentId,
				progressiveRollout.Id,
				s.ScheduleId,
			); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}
