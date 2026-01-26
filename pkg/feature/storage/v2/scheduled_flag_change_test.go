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

package v2

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestNewScheduledFlagChangeStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewScheduledFlagChangeStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &scheduledFlagChangeStorage{}, storage)
}

func TestScheduledFlagChangeStorageCreateScheduledFlagChange(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc                string
		setup               func(storage *scheduledFlagChangeStorage)
		scheduledFlagChange *domain.ScheduledFlagChange
		expectedErr         error
	}{
		{
			desc: "error: general error",
			setup: func(s *scheduledFlagChangeStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			scheduledFlagChange: &domain.ScheduledFlagChange{
				ScheduledFlagChange: &proto.ScheduledFlagChange{
					Id:            "sfc-1",
					FeatureId:     "feature-1",
					EnvironmentId: "env-1",
				},
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: duplicate entry",
			setup: func(s *scheduledFlagChangeStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			scheduledFlagChange: &domain.ScheduledFlagChange{
				ScheduledFlagChange: &proto.ScheduledFlagChange{
					Id:            "sfc-1",
					FeatureId:     "feature-1",
					EnvironmentId: "env-1",
				},
			},
			expectedErr: ErrScheduledFlagChangeAlreadyExists,
		},
		{
			desc: "success",
			setup: func(s *scheduledFlagChangeStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			scheduledFlagChange: &domain.ScheduledFlagChange{
				ScheduledFlagChange: &proto.ScheduledFlagChange{
					Id:            "sfc-1",
					FeatureId:     "feature-1",
					EnvironmentId: "env-1",
					ScheduledAt:   1700000000,
					Timezone:      "UTC",
					Status:        proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &scheduledFlagChangeStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			err := storage.CreateScheduledFlagChange(context.Background(), p.scheduledFlagChange)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestScheduledFlagChangeStorageUpdateScheduledFlagChange(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc                string
		setup               func(storage *scheduledFlagChangeStorage)
		scheduledFlagChange *domain.ScheduledFlagChange
		expectedErr         error
	}{
		{
			desc: "error: exec error",
			setup: func(s *scheduledFlagChangeStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			scheduledFlagChange: &domain.ScheduledFlagChange{
				ScheduledFlagChange: &proto.ScheduledFlagChange{
					Id:            "sfc-1",
					EnvironmentId: "env-1",
				},
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: no rows affected",
			setup: func(s *scheduledFlagChangeStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
			},
			scheduledFlagChange: &domain.ScheduledFlagChange{
				ScheduledFlagChange: &proto.ScheduledFlagChange{
					Id:            "sfc-1",
					EnvironmentId: "env-1",
				},
			},
			expectedErr: ErrScheduledFlagChangeUnexpectedAffectedRows,
		},
		{
			desc: "success",
			setup: func(s *scheduledFlagChangeStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
			},
			scheduledFlagChange: &domain.ScheduledFlagChange{
				ScheduledFlagChange: &proto.ScheduledFlagChange{
					Id:            "sfc-1",
					EnvironmentId: "env-1",
					ScheduledAt:   1700000000,
					Status:        proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &scheduledFlagChangeStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			err := storage.UpdateScheduledFlagChange(context.Background(), p.scheduledFlagChange)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestScheduledFlagChangeStorageDeleteScheduledFlagChange(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc          string
		setup         func(storage *scheduledFlagChangeStorage)
		id            string
		environmentID string
		expectedErr   error
	}{
		{
			desc: "error: exec error",
			setup: func(s *scheduledFlagChangeStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			id:            "sfc-1",
			environmentID: "env-1",
			expectedErr:   errors.New("error"),
		},
		{
			desc: "error: no rows affected",
			setup: func(s *scheduledFlagChangeStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
			},
			id:            "sfc-1",
			environmentID: "env-1",
			expectedErr:   ErrScheduledFlagChangeUnexpectedAffectedRows,
		},
		{
			desc: "success",
			setup: func(s *scheduledFlagChangeStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
			},
			id:            "sfc-1",
			environmentID: "env-1",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &scheduledFlagChangeStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			err := storage.DeleteScheduledFlagChange(context.Background(), p.id, p.environmentID)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestScheduledFlagChangeStorageGetScheduledFlagChange(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc          string
		setup         func(storage *scheduledFlagChangeStorage)
		id            string
		environmentID string
		expectedErr   error
	}{
		{
			desc: "error: general error",
			setup: func(s *scheduledFlagChangeStorage) {
				row := mock.NewMockRow(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
			},
			id:            "sfc-1",
			environmentID: "env-1",
			expectedErr:   errors.New("error"),
		},
		{
			desc: "error: not found",
			setup: func(s *scheduledFlagChangeStorage) {
				row := mock.NewMockRow(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
			},
			id:            "sfc-1",
			environmentID: "env-1",
			expectedErr:   ErrScheduledFlagChangeNotFound,
		},
		{
			desc: "success",
			setup: func(s *scheduledFlagChangeStorage) {
				row := mock.NewMockRow(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
			},
			id:            "sfc-1",
			environmentID: "env-1",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &scheduledFlagChangeStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			_, err := storage.GetScheduledFlagChange(context.Background(), p.id, p.environmentID)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestScheduledFlagChangeStorageListScheduledFlagChanges(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(storage *scheduledFlagChangeStorage)
		options        *mysql.ListOptions
		expected       []*proto.ScheduledFlagChange
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "error: query error",
			setup: func(s *scheduledFlagChangeStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			options:     nil,
			expected:    nil,
			expectedErr: errors.New("error"),
		},
		{
			desc: "success: empty result",
			setup: func(s *scheduledFlagChangeStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			options: &mysql.ListOptions{
				Limit:   0,
				Offset:  0,
				Filters: []*mysql.FilterV2{},
				Orders:  []*mysql.Order{},
			},
			expected:       []*proto.ScheduledFlagChange{},
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &scheduledFlagChangeStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			expected, nextOffset, _, err := storage.ListScheduledFlagChanges(context.Background(), p.options)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, expected)
			assert.Equal(t, p.expectedCursor, nextOffset)
		})
	}
}

func TestScheduledFlagChangeStorageListDueScheduledFlagChanges(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(storage *scheduledFlagChangeStorage)
		now         int64
		limit       int
		expected    []*proto.ScheduledFlagChange
		expectedErr error
	}{
		{
			desc: "error: query error",
			setup: func(s *scheduledFlagChangeStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			now:         1700000000,
			limit:       100,
			expected:    nil,
			expectedErr: errors.New("error"),
		},
		{
			desc: "success: empty result",
			setup: func(s *scheduledFlagChangeStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			now:         1700000000,
			limit:       100,
			expected:    []*proto.ScheduledFlagChange{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &scheduledFlagChangeStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			expected, err := storage.ListDueScheduledFlagChanges(context.Background(), p.now, p.limit)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, expected)
		})
	}
}

func TestScheduledFlagChangeStorageTryLock(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(storage *scheduledFlagChangeStorage)
		id             string
		lockedBy       string
		expectedLocked bool
		expectedErr    error
	}{
		{
			desc: "error: exec error",
			setup: func(s *scheduledFlagChangeStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			id:             "sfc-1",
			lockedBy:       "executor-1",
			expectedLocked: false,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "success: lock not acquired",
			setup: func(s *scheduledFlagChangeStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
			},
			id:             "sfc-1",
			lockedBy:       "executor-1",
			expectedLocked: false,
			expectedErr:    nil,
		},
		{
			desc: "success: lock acquired",
			setup: func(s *scheduledFlagChangeStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
			},
			id:             "sfc-1",
			lockedBy:       "executor-1",
			expectedLocked: true,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &scheduledFlagChangeStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			locked, err := storage.TryLock(context.Background(), p.id, p.lockedBy)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expectedLocked, locked)
		})
	}
}

func TestScheduledFlagChangeStorageUnlock(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(storage *scheduledFlagChangeStorage)
		id          string
		expectedErr error
	}{
		{
			desc: "error: exec error",
			setup: func(s *scheduledFlagChangeStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			id:          "sfc-1",
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *scheduledFlagChangeStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			id:          "sfc-1",
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &scheduledFlagChangeStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			err := storage.Unlock(context.Background(), p.id)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
