// Copyright 2023 The Bucketeer Authors.
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

	aoclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	autoopsdomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	envclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	"github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor"
	aoproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	envproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

type progressiveRolloutWatcher struct {
	envClient                  envclient.Client
	aoClient                   aoclient.Client
	progressiveRolloutExecutor executor.ProgressiveRolloutExecutor
	opts                       *jobs.Options
	logger                     *zap.Logger
}

func NewProgressiveRolloutWacher(
	envClient envclient.Client,
	aoClient aoclient.Client,
	progressiveRolloutExecutor executor.ProgressiveRolloutExecutor,
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
		for _, p := range progressiveRollouts {
			lastErr = w.executeProgressiveRollout(ctx, p, e.Id)
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

func (s *progressiveRolloutWatcher) listProgressiveRollouts(
	ctx context.Context,
	environmentID string,
) ([]*aoproto.ProgressiveRollout, error) {
	resp, err := s.aoClient.ListProgressiveRollouts(
		ctx,
		&aoproto.ListProgressiveRolloutsRequest{
			EnvironmentNamespace: environmentID,
			PageSize:             0,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp.ProgressiveRollouts, nil
}

func (w *progressiveRolloutWatcher) executeProgressiveRollout(
	ctx context.Context,
	progressiveRollout *aoproto.ProgressiveRollout,
	environmentNamespace string,
) error {
	pr := &autoopsdomain.ProgressiveRollout{ProgressiveRollout: progressiveRollout}
	schedules, err := pr.ExtractSchedules()
	if err != nil {
		w.logger.Error("Failed to extract schedules", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("featureId", progressiveRollout.FeatureId),
			zap.String("progressiveRolloutId", progressiveRollout.Id),
		)
		return err
	}
	now := time.Now().Unix()
	for _, s := range schedules {
		if s.TriggeredAt == 0 && s.ExecuteAt <= now {
			w.logger.Info("scheduled time is passed",
				zap.String("environmentNamespace", environmentNamespace),
				zap.String("featureId", progressiveRollout.FeatureId),
				zap.String("progressiveRolloutId", progressiveRollout.Id),
				zap.Any("datetimeClause", s),
			)
			if err := w.progressiveRolloutExecutor.ExecuteProgressiveRollout(
				ctx,
				environmentNamespace,
				progressiveRollout.Id,
				s.ScheduleId,
			); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}
