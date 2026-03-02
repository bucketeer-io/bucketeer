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
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
)

func TestMySQLEventStorageQueryEvaluationCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := context.Background()
	startAt := time.Unix(100, 0)
	endAt := time.Unix(200, 0)

	patterns := []struct {
		desc        string
		setup       func(s *mysqlEventStorage)
		expected    []*EvaluationEventCount
		expectedErr error
	}{
		{
			desc: "success",
			setup: func(s *mysqlEventStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(true)
				rows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(dest ...any) error {
						*(dest[0].(*string)) = "vid1"
						*(dest[1].(*int64)) = int64(1)
						*(dest[2].(*int64)) = int64(2)
						return nil
					},
				)
				rows.EXPECT().Next().Return(true)
				rows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(dest ...any) error {
						*(dest[0].(*string)) = "vid2"
						*(dest[1].(*int64)) = int64(3)
						*(dest[2].(*int64)) = int64(4)
						return nil
					},
				)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)

			},
			expected: []*EvaluationEventCount{
				{VariationID: "vid1", EvaluationUser: 1, EvaluationTotal: 2},
				{VariationID: "vid2", EvaluationUser: 3, EvaluationTotal: 4},
			},
			expectedErr: nil,
		},
		{
			desc: "error: query",
			setup: func(s *mysqlEventStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("query error"))
			},
			expectedErr: errors.New("query error"),
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &mysqlEventStorage{
				qe:     mock.NewMockQueryExecer(mockController),
				logger: zap.NewNop(),
			}
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.QueryEvaluationCount(ctx, "env", startAt, endAt, "fid", 1)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestMySQLEventStorageQueryGoalCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := context.Background()
	startAt := time.Unix(100, 0)
	endAt := time.Unix(200, 0)

	patterns := []struct {
		desc        string
		setup       func(s *mysqlEventStorage)
		expected    []*GoalEventCount
		expectedErr error
	}{
		{
			desc: "success",
			setup: func(s *mysqlEventStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(true)
				rows.EXPECT().Scan(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).DoAndReturn(
					func(dest ...any) error {
						*(dest[0].(*string)) = "vid1"
						*(dest[1].(*int64)) = int64(1)
						*(dest[2].(*int64)) = int64(2)
						*(dest[3].(*float64)) = float64(3.5)
						*(dest[4].(*float64)) = float64(1.75)
						*(dest[5].(*float64)) = float64(0.25)
						return nil
					},
				)
				rows.EXPECT().Next().Return(true)
				rows.EXPECT().Scan(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).DoAndReturn(
					func(dest ...any) error {
						*(dest[0].(*string)) = "vid2"
						*(dest[1].(*int64)) = int64(3)
						*(dest[2].(*int64)) = int64(4)
						*(dest[3].(*float64)) = float64(5.5)
						*(dest[4].(*float64)) = float64(2.75)
						*(dest[5].(*float64)) = float64(0.5)
						return nil
					},
				)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)

			},
			expected: []*GoalEventCount{
				{
					VariationID:       "vid1",
					GoalUser:          1,
					GoalTotal:         2,
					GoalValueTotal:    3.5,
					GoalValueMean:     1.75,
					GoalValueVariance: 0.25,
				},
				{
					VariationID:       "vid2",
					GoalUser:          3,
					GoalTotal:         4,
					GoalValueTotal:    5.5,
					GoalValueMean:     2.75,
					GoalValueVariance: 0.5,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "error: query",
			setup: func(s *mysqlEventStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("query error"))
			},
			expectedErr: errors.New("query error"),
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &mysqlEventStorage{
				qe:     mock.NewMockQueryExecer(mockController),
				logger: zap.NewNop(),
			}
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.QueryGoalCount(ctx, "env", startAt, endAt, "gid", "fid", 1)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestMySQLEventStorageQueryUserEvaluation(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := context.Background()
	startAt := time.Unix(100, 0)
	endAt := time.Unix(200, 0)

	patterns := []struct {
		desc        string
		setup       func(s *mysqlEventStorage)
		expected    *UserEvaluation
		expectedErr error
	}{
		{
			desc: "success",
			setup: func(s *mysqlEventStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(true)
				rows.EXPECT().Scan(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).DoAndReturn(
					func(dest ...any) error {
						*(dest[0].(*string)) = "uid"
						*(dest[1].(*string)) = "fid"
						*(dest[2].(*int32)) = int32(1)
						*(dest[3].(*string)) = "vid"
						*(dest[4].(*string)) = "reason"
						*(dest[5].(*int64)) = int64(123)
						return nil
					},
				)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expected: &UserEvaluation{
				UserID:         "uid",
				FeatureID:      "fid",
				FeatureVersion: 1,
				VariationID:    "vid",
				Reason:         "reason",
				Timestamp:      123,
			},
			expectedErr: nil,
		},
		{
			desc: "error: query",
			setup: func(s *mysqlEventStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("query error"))
			},
			expectedErr: errors.New("query error"),
		},
		{
			desc: "error: no results",
			setup: func(s *mysqlEventStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expectedErr: ErrMySQLNoResultsFound,
		},
		{
			desc: "error: scan",
			setup: func(s *mysqlEventStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(true)
				rows.EXPECT().Scan(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("scan error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expectedErr: errors.New("scan error"),
		},
		{
			desc: "error: unexpected multiple results",
			setup: func(s *mysqlEventStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(true)
				rows.EXPECT().Scan(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				rows.EXPECT().Next().Return(true)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expectedErr: ErrMySQLUnexpectedMultipleResults,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &mysqlEventStorage{
				qe:     mock.NewMockQueryExecer(mockController),
				logger: zap.NewNop(),
			}
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.QueryUserEvaluation(ctx, "env", "uid", "fid", 1, startAt, endAt)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}
