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

package scheduledflagchange

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	featuredomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	// maxDueScheduleBatchSize is the maximum number of due schedules to process per run.
	maxDueScheduleBatchSize = 100

	// executorLockID is the identifier for this executor instance (used in TryLock).
	executorLockID = "batch-executor"
)

// scheduledFlagChangeExecutor is a batch job that executes due scheduled flag changes.
// It polls for schedules where scheduled_at <= now and status = PENDING,
// then executes them via the feature service API.
// Schedules in CONFLICT status that are past due are marked as FAILED.
type scheduledFlagChangeExecutor struct {
	sfcStorage    v2fs.ScheduledFlagChangeStorage
	featureClient featureclient.Client
	opts          *jobs.Options
	logger        *zap.Logger
}

// NewScheduledFlagChangeExecutor creates a new scheduled flag change executor batch job.
func NewScheduledFlagChangeExecutor(
	mysqlClient mysql.Client,
	featureClient featureclient.Client,
	opts ...jobs.Option,
) jobs.Job {
	dopts := &jobs.Options{
		Timeout: 50 * time.Second,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &scheduledFlagChangeExecutor{
		sfcStorage:    v2fs.NewScheduledFlagChangeStorage(mysqlClient),
		featureClient: featureClient,
		opts:          dopts,
		logger:        dopts.Logger.Named("scheduled-flag-change-executor"),
	}
}

// Run executes the scheduled flag change executor batch job.
// It processes due schedules in two phases:
//  1. Execute PENDING schedules via the feature service API
//  2. Mark past-due CONFLICT schedules as FAILED
func (e *scheduledFlagChangeExecutor) Run(ctx context.Context) (lastErr error) {
	startTime := time.Now()
	defer func() {
		jobs.RecordJob(jobs.JobScheduledFlagChangeExecutor, lastErr, time.Since(startTime))
	}()

	ctx, cancel := context.WithTimeout(ctx, e.opts.Timeout)
	defer cancel()

	e.logger.Info("ScheduledFlagChangeExecutor start running")

	now := time.Now().Unix()

	// Execute due PENDING schedules
	executed, failed, err := e.executeDueSchedules(ctx, now)
	if err != nil {
		e.logger.Error("Error executing due schedules", zap.Error(err))
		lastErr = err
	}

	// Mark past-due CONFLICT schedules as FAILED
	skipped, err := e.skipConflictSchedules(ctx, now)
	if err != nil {
		e.logger.Error("Error skipping conflict schedules", zap.Error(err))
		lastErr = err
	}

	e.logger.Info("ScheduledFlagChangeExecutor finished",
		zap.Int("executed", executed),
		zap.Int("failed", failed),
		zap.Int("skipped_conflict", skipped),
		zap.Duration("elapsedTime", time.Since(startTime)),
	)
	return lastErr
}

// executeDueSchedules lists and executes all due PENDING schedules.
// Returns counts of executed and failed schedules.
func (e *scheduledFlagChangeExecutor) executeDueSchedules(ctx context.Context, now int64) (int, int, error) {
	dueSchedules, err := e.sfcStorage.ListDueScheduledFlagChanges(ctx, now, maxDueScheduleBatchSize)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to list due schedules: %w", err)
	}

	if len(dueSchedules) == 0 {
		e.logger.Debug("No due schedules found")
		return 0, 0, nil
	}

	e.logger.Info("Found due schedules", zap.Int("count", len(dueSchedules)))

	executed := 0
	failed := 0

	for _, schedule := range dueSchedules {
		if err := ctx.Err(); err != nil {
			e.logger.Warn("Context cancelled, stopping execution", zap.Error(err))
			return executed, failed, err
		}

		if err := e.executeOne(ctx, schedule); err != nil {
			e.logger.Error("Failed to execute scheduled flag change",
				zap.String("id", schedule.Id),
				zap.String("featureId", schedule.FeatureId),
				zap.String("environmentId", schedule.EnvironmentId),
				zap.Error(err),
			)
			failed++
			continue
		}
		executed++
	}

	return executed, failed, nil
}

// executeOne executes a single scheduled flag change.
// It acquires a lock, then calls the feature service API to execute the schedule.
func (e *scheduledFlagChangeExecutor) executeOne(
	ctx context.Context,
	schedule *featureproto.ScheduledFlagChange,
) error {
	// Acquire lock to prevent concurrent execution
	locked, err := e.sfcStorage.TryLock(ctx, schedule.Id, executorLockID)
	if err != nil {
		return fmt.Errorf("failed to acquire lock for schedule %s: %w", schedule.Id, err)
	}
	if !locked {
		e.logger.Debug("Schedule already locked by another executor",
			zap.String("id", schedule.Id),
		)
		return nil
	}
	defer func() {
		if unlockErr := e.sfcStorage.Unlock(ctx, schedule.Id, executorLockID); unlockErr != nil {
			e.logger.Error("Failed to release lock",
				zap.String("id", schedule.Id),
				zap.Error(unlockErr),
			)
		}
	}()

	// Execute the schedule via the feature service API.
	// ExecuteScheduledFlagChange handles:
	// - Fetching the feature
	// - Validating references
	// - Applying changes within a transaction
	// - Marking the schedule as EXECUTED
	// - Publishing domain events (FEATURE_UPDATED)
	// - Updating the feature flag cache
	_, err = e.featureClient.ExecuteScheduledFlagChange(ctx, &featureproto.ExecuteScheduledFlagChangeRequest{
		EnvironmentId: schedule.EnvironmentId,
		Id:            schedule.Id,
	})
	if err != nil {
		// The API method already marks the schedule as FAILED internally when
		// validation fails or the feature is not found. So we just log here.
		e.logger.Error("ExecuteScheduledFlagChange API returned error",
			zap.String("id", schedule.Id),
			zap.String("featureId", schedule.FeatureId),
			zap.String("environmentId", schedule.EnvironmentId),
			zap.Error(err),
		)
		return err
	}

	e.logger.Info("Successfully executed scheduled flag change",
		zap.String("id", schedule.Id),
		zap.String("featureId", schedule.FeatureId),
		zap.String("environmentId", schedule.EnvironmentId),
	)
	return nil
}

// skipConflictSchedules finds past-due schedules in CONFLICT status and marks them as FAILED.
// These schedules had unresolved conflicts when their execution time arrived.
func (e *scheduledFlagChangeExecutor) skipConflictSchedules(ctx context.Context, now int64) (int, error) {
	conflictSchedules, err := e.listDueConflictSchedules(ctx, now)
	if err != nil {
		return 0, fmt.Errorf("failed to list due conflict schedules: %w", err)
	}

	if len(conflictSchedules) == 0 {
		return 0, nil
	}

	e.logger.Info("Found past-due conflict schedules to skip",
		zap.Int("count", len(conflictSchedules)),
	)

	skipped := 0
	for _, schedule := range conflictSchedules {
		if err := ctx.Err(); err != nil {
			return skipped, err
		}

		sfcDomain := &featuredomain.ScheduledFlagChange{ScheduledFlagChange: schedule}
		sfcDomain.MarkFailed("Skipped due to unresolved conflict at scheduled execution time")

		if err := e.sfcStorage.UpdateScheduledFlagChange(ctx, sfcDomain); err != nil {
			e.logger.Error("Failed to mark conflict schedule as failed",
				zap.String("id", schedule.Id),
				zap.String("featureId", schedule.FeatureId),
				zap.Error(err),
			)
			continue
		}

		e.logger.Info("Skipped conflict schedule",
			zap.String("id", schedule.Id),
			zap.String("featureId", schedule.FeatureId),
			zap.String("environmentId", schedule.EnvironmentId),
		)
		skipped++
	}

	return skipped, nil
}

// listDueConflictSchedules lists scheduled flag changes that are past due and in CONFLICT status.
func (e *scheduledFlagChangeExecutor) listDueConflictSchedules(
	ctx context.Context,
	now int64,
) ([]*featureproto.ScheduledFlagChange, error) {
	filters := []*mysql.FilterV2{
		{
			Column:   "status",
			Operator: mysql.OperatorEqual,
			Value:    int32(featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT),
		},
		{
			Column:   "scheduled_at",
			Operator: mysql.OperatorLessThanOrEqual,
			Value:    now,
		},
	}
	options := &mysql.ListOptions{
		Filters: filters,
		Orders:  []*mysql.Order{mysql.NewOrder("scheduled_at", mysql.OrderDirectionAsc)},
		Limit:   maxDueScheduleBatchSize,
		Offset:  mysql.QueryNoOffset,
	}
	sfcs, _, _, err := e.sfcStorage.ListScheduledFlagChanges(ctx, options)
	return sfcs, err
}
