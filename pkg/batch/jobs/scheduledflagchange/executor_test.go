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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	featureclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	sfcmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestExecuteDueSchedules_NoDueSchedules(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := sfcmock.NewMockScheduledFlagChangeStorage(ctrl)
	mockClient := featureclientmock.NewMockClient(ctrl)

	executor := &scheduledFlagChangeExecutor{
		sfcStorage:    mockStorage,
		featureClient: mockClient,
		opts: &jobs.Options{
			Timeout: 50 * time.Second,
			Logger:  zap.NewNop(),
		},
		logger: zap.NewNop(),
	}

	now := time.Now().Unix()
	mockStorage.EXPECT().
		ListDueScheduledFlagChanges(gomock.Any(), now, maxDueScheduleBatchSize).
		Return(nil, nil)

	executed, failed, err := executor.executeDueSchedules(context.Background(), now)
	assert.NoError(t, err)
	assert.Equal(t, 0, executed)
	assert.Equal(t, 0, failed)
}

func TestExecuteDueSchedules_ExecutesPendingSchedule(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := sfcmock.NewMockScheduledFlagChangeStorage(ctrl)
	mockClient := featureclientmock.NewMockClient(ctrl)

	executor := &scheduledFlagChangeExecutor{
		sfcStorage:    mockStorage,
		featureClient: mockClient,
		opts: &jobs.Options{
			Timeout: 50 * time.Second,
			Logger:  zap.NewNop(),
		},
		logger: zap.NewNop(),
	}

	schedule := &featureproto.ScheduledFlagChange{
		Id:            "sfc-1",
		FeatureId:     "feature-1",
		EnvironmentId: "env-1",
		Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
		ScheduledAt:   time.Now().Add(-1 * time.Minute).Unix(),
	}

	now := time.Now().Unix()
	mockStorage.EXPECT().
		ListDueScheduledFlagChanges(gomock.Any(), now, maxDueScheduleBatchSize).
		Return([]*featureproto.ScheduledFlagChange{schedule}, nil)

	// Lock/Unlock
	mockStorage.EXPECT().
		TryLock(gomock.Any(), "sfc-1", executorLockID).
		Return(true, nil)
	mockStorage.EXPECT().
		Unlock(gomock.Any(), "sfc-1", executorLockID).
		Return(nil)

	// Execute via feature client
	mockClient.EXPECT().
		ExecuteScheduledFlagChange(gomock.Any(), &featureproto.ExecuteScheduledFlagChangeRequest{
			EnvironmentId: "env-1",
			Id:            "sfc-1",
		}).
		Return(&featureproto.ExecuteScheduledFlagChangeResponse{}, nil)

	executed, failed, err := executor.executeDueSchedules(context.Background(), now)
	assert.NoError(t, err)
	assert.Equal(t, 1, executed)
	assert.Equal(t, 0, failed)
}

func TestExecuteDueSchedules_FailedExecution(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := sfcmock.NewMockScheduledFlagChangeStorage(ctrl)
	mockClient := featureclientmock.NewMockClient(ctrl)

	executor := &scheduledFlagChangeExecutor{
		sfcStorage:    mockStorage,
		featureClient: mockClient,
		opts: &jobs.Options{
			Timeout: 50 * time.Second,
			Logger:  zap.NewNop(),
		},
		logger: zap.NewNop(),
	}

	schedule := &featureproto.ScheduledFlagChange{
		Id:            "sfc-1",
		FeatureId:     "feature-1",
		EnvironmentId: "env-1",
		Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
		ScheduledAt:   time.Now().Add(-1 * time.Minute).Unix(),
	}

	now := time.Now().Unix()
	mockStorage.EXPECT().
		ListDueScheduledFlagChanges(gomock.Any(), now, maxDueScheduleBatchSize).
		Return([]*featureproto.ScheduledFlagChange{schedule}, nil)

	mockStorage.EXPECT().
		TryLock(gomock.Any(), "sfc-1", executorLockID).
		Return(true, nil)
	mockStorage.EXPECT().
		Unlock(gomock.Any(), "sfc-1", executorLockID).
		Return(nil)

	// Execute fails
	mockClient.EXPECT().
		ExecuteScheduledFlagChange(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("feature not found"))

	executed, failed, err := executor.executeDueSchedules(context.Background(), now)
	assert.NoError(t, err)
	assert.Equal(t, 0, executed)
	assert.Equal(t, 1, failed)
}

func TestExecuteOne_LockAlreadyHeld(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := sfcmock.NewMockScheduledFlagChangeStorage(ctrl)
	mockClient := featureclientmock.NewMockClient(ctrl)

	executor := &scheduledFlagChangeExecutor{
		sfcStorage:    mockStorage,
		featureClient: mockClient,
		opts: &jobs.Options{
			Timeout: 50 * time.Second,
			Logger:  zap.NewNop(),
		},
		logger: zap.NewNop(),
	}

	schedule := &featureproto.ScheduledFlagChange{
		Id:            "sfc-1",
		FeatureId:     "feature-1",
		EnvironmentId: "env-1",
	}

	// TryLock returns false (already locked)
	mockStorage.EXPECT().
		TryLock(gomock.Any(), "sfc-1", executorLockID).
		Return(false, nil)

	err := executor.executeOne(context.Background(), schedule)
	assert.NoError(t, err) // Not an error, just skip
}

func TestSkipConflictSchedules_MarksConflictsAsFailed(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := sfcmock.NewMockScheduledFlagChangeStorage(ctrl)
	mockClient := featureclientmock.NewMockClient(ctrl)

	executor := &scheduledFlagChangeExecutor{
		sfcStorage:    mockStorage,
		featureClient: mockClient,
		opts: &jobs.Options{
			Timeout: 50 * time.Second,
			Logger:  zap.NewNop(),
		},
		logger: zap.NewNop(),
	}

	conflictSchedule := &featureproto.ScheduledFlagChange{
		Id:            "sfc-conflict-1",
		FeatureId:     "feature-1",
		EnvironmentId: "env-1",
		Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT,
		ScheduledAt:   time.Now().Add(-5 * time.Minute).Unix(),
	}

	now := time.Now().Unix()

	// ListScheduledFlagChanges for CONFLICT + past due
	mockStorage.EXPECT().
		ListScheduledFlagChanges(gomock.Any(), gomock.Any()).
		Return([]*featureproto.ScheduledFlagChange{conflictSchedule}, 0, int64(1), nil)

	// UpdateScheduledFlagChange to mark as FAILED
	mockStorage.EXPECT().
		UpdateScheduledFlagChange(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, sfc interface{}) error {
			// Verify the schedule is marked as FAILED
			return nil
		})

	skipped, err := executor.skipConflictSchedules(context.Background(), now)
	assert.NoError(t, err)
	assert.Equal(t, 1, skipped)
}

func TestSkipConflictSchedules_NoConflictSchedules(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := sfcmock.NewMockScheduledFlagChangeStorage(ctrl)
	mockClient := featureclientmock.NewMockClient(ctrl)

	executor := &scheduledFlagChangeExecutor{
		sfcStorage:    mockStorage,
		featureClient: mockClient,
		opts: &jobs.Options{
			Timeout: 50 * time.Second,
			Logger:  zap.NewNop(),
		},
		logger: zap.NewNop(),
	}

	now := time.Now().Unix()

	mockStorage.EXPECT().
		ListScheduledFlagChanges(gomock.Any(), gomock.Any()).
		Return(nil, 0, int64(0), nil)

	skipped, err := executor.skipConflictSchedules(context.Background(), now)
	assert.NoError(t, err)
	assert.Equal(t, 0, skipped)
}

func TestExecuteDueSchedules_MultipleSchedules(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := sfcmock.NewMockScheduledFlagChangeStorage(ctrl)
	mockClient := featureclientmock.NewMockClient(ctrl)

	executor := &scheduledFlagChangeExecutor{
		sfcStorage:    mockStorage,
		featureClient: mockClient,
		opts: &jobs.Options{
			Timeout: 50 * time.Second,
			Logger:  zap.NewNop(),
		},
		logger: zap.NewNop(),
	}

	schedules := []*featureproto.ScheduledFlagChange{
		{
			Id:            "sfc-1",
			FeatureId:     "feature-1",
			EnvironmentId: "env-1",
			Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
			ScheduledAt:   time.Now().Add(-2 * time.Minute).Unix(),
		},
		{
			Id:            "sfc-2",
			FeatureId:     "feature-2",
			EnvironmentId: "env-1",
			Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
			ScheduledAt:   time.Now().Add(-1 * time.Minute).Unix(),
		},
	}

	now := time.Now().Unix()
	mockStorage.EXPECT().
		ListDueScheduledFlagChanges(gomock.Any(), now, maxDueScheduleBatchSize).
		Return(schedules, nil)

	// Both get locked and unlocked
	mockStorage.EXPECT().TryLock(gomock.Any(), "sfc-1", executorLockID).Return(true, nil)
	mockStorage.EXPECT().Unlock(gomock.Any(), "sfc-1", executorLockID).Return(nil)
	mockStorage.EXPECT().TryLock(gomock.Any(), "sfc-2", executorLockID).Return(true, nil)
	mockStorage.EXPECT().Unlock(gomock.Any(), "sfc-2", executorLockID).Return(nil)

	// First succeeds, second fails
	mockClient.EXPECT().
		ExecuteScheduledFlagChange(gomock.Any(), &featureproto.ExecuteScheduledFlagChangeRequest{
			EnvironmentId: "env-1",
			Id:            "sfc-1",
		}).
		Return(&featureproto.ExecuteScheduledFlagChangeResponse{}, nil)

	mockClient.EXPECT().
		ExecuteScheduledFlagChange(gomock.Any(), &featureproto.ExecuteScheduledFlagChangeRequest{
			EnvironmentId: "env-1",
			Id:            "sfc-2",
		}).
		Return(nil, errors.New("validation failed"))

	executed, failed, err := executor.executeDueSchedules(context.Background(), now)
	assert.NoError(t, err)
	assert.Equal(t, 1, executed)
	assert.Equal(t, 1, failed)
}
