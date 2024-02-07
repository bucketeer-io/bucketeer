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

package v2

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

func TestNewGoalStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	db := NewGoalStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &goalStorage{}, db)
}

func TestCreateGoal(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup                func(*goalStorage)
		input                *domain.Goal
		environmentNamespace string
		expectedErr          error
	}{
		{
			setup: func(s *goalStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: &domain.Goal{
				Goal: &proto.Goal{Id: "id-0"},
			},
			environmentNamespace: "ns0",
			expectedErr:          ErrGoalAlreadyExists,
		},
		{
			setup: func(s *goalStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: &domain.Goal{
				Goal: &proto.Goal{Id: "id-1"},
			},
			environmentNamespace: "ns0",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		db := newGoalStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(db)
		}
		err := db.CreateGoal(ctx, p.input, p.environmentNamespace)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestUpdateGoal(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup                func(*goalStorage)
		input                *domain.Goal
		environmentNamespace string
		expectedErr          error
	}{
		{
			setup: func(s *goalStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.Goal{
				Goal: &proto.Goal{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          ErrGoalUnexpectedAffectedRows,
		},
		{
			setup: func(s *goalStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.Goal{
				Goal: &proto.Goal{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		storage := newGoalStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		err := storage.UpdateGoal(context.Background(), p.input, p.environmentNamespace)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestGetGoal(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup                func(*goalStorage)
		input                string
		environmentNamespace string
		expected             *domain.Goal
		expectedErr          error
	}{
		{
			setup: func(s *goalStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input:                "",
			environmentNamespace: "ns0",
			expected:             nil,
			expectedErr:          ErrGoalNotFound,
		},
		{
			setup: func(s *goalStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input:                "id-0",
			environmentNamespace: "ns0",
			expected: &domain.Goal{
				Goal: &proto.Goal{Id: "id-0"},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		storage := newGoalStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		_, err := storage.GetGoal(context.Background(), p.input, p.environmentNamespace)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestListGoals(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		setup                func(*goalStorage)
		whereParts           []mysql.WherePart
		orders               []*mysql.Order
		limit                int
		offset               int
		isInUseStatus        *bool
		environmentNamespace string
		expected             []*proto.Goal
		expectedCursor       int
		expectedErr          error
	}{
		{
			setup: func(s *goalStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			whereParts:           nil,
			orders:               nil,
			limit:                0,
			offset:               0,
			isInUseStatus:        nil,
			environmentNamespace: "",
			expected:             nil,
			expectedCursor:       0,
			expectedErr:          errors.New("error"),
		},
		{
			setup: func(s *goalStorage) {
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
			whereParts: []mysql.WherePart{
				mysql.NewFilter("num", ">=", 5),
			},
			orders: []*mysql.Order{
				mysql.NewOrder("id", mysql.OrderDirectionAsc),
			},
			limit:                10,
			offset:               5,
			isInUseStatus:        nil,
			environmentNamespace: "ns0",
			expected:             []*proto.Goal{},
			expectedCursor:       5,
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		storage := newGoalStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		goals, cursor, _, err := storage.ListGoals(
			context.Background(),
			p.whereParts,
			p.orders,
			p.limit,
			p.offset,
			p.isInUseStatus,
			p.environmentNamespace,
		)
		assert.Equal(t, p.expected, goals)
		assert.Equal(t, p.expectedCursor, cursor)
		assert.Equal(t, p.expectedErr, err)
	}
}

func newGoalStorageWithMock(t *testing.T, mockController *gomock.Controller) *goalStorage {
	t.Helper()
	return &goalStorage{mock.NewMockQueryExecer(mockController)}
}
