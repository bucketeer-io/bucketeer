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

package job

import (
	"context"
	"time"

	"go.uber.org/zap"

	autoopsdomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/pkg/job"
	"github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor"
	"github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/targetstore"
)

type progressiveRolloutWatcher struct {
	environmentLister          targetstore.EnvironmentLister
	progressiveRolloutLister   targetstore.ProgressiveRolloutLister
	progressiveRolloutExecutor executor.ProgressiveRolloutExecutor
	opts                       *options
	logger                     *zap.Logger
}

func NewProgressiveRolloutWacher(
	targetStore targetstore.TargetStore,
	progressiveRolloutExecutor executor.ProgressiveRolloutExecutor,
	opts ...Option,
) job.Job {
	dopts := &options{
		timeout: 5 * time.Minute,
		logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &progressiveRolloutWatcher{
		environmentLister:          targetStore,
		progressiveRolloutLister:   targetStore,
		progressiveRolloutExecutor: progressiveRolloutExecutor,
		opts:                       dopts,
		logger:                     dopts.logger.Named("progressive-rollout-watcher"),
	}
}

func (w *progressiveRolloutWatcher) Run(
	ctx context.Context,
) error {
	ctx, cancel := context.WithTimeout(ctx, w.opts.timeout)
	defer cancel()
	environments := w.environmentLister.GetEnvironments(ctx)
	var lastErr error
	for _, e := range environments {
		progressiveRollouts := w.progressiveRolloutLister.GetProgressiveRollouts(
			ctx,
			e.Namespace,
		)
		for _, p := range progressiveRollouts {
			lastErr = w.executeProgressiveRollout(ctx, p, e.Namespace)
		}
	}
	return lastErr
}

func (w *progressiveRolloutWatcher) executeProgressiveRollout(
	ctx context.Context,
	progressiveRollout *autoopsdomain.ProgressiveRollout,
	environmentNamespace string,
) error {
	schedules, err := progressiveRollout.ExtractSchedules()
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
