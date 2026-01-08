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

	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

func TestNewExperimentStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	db := NewExperimentStorage(mock.NewMockClient(mockController))
	assert.IsType(t, &experimentStorage{}, db)
}

func TestCreateExperiment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup         func(*experimentStorage)
		input         *domain.Experiment
		environmentId string
		expectedErr   error
	}{
		{
			setup: func(s *experimentStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: &domain.Experiment{
				Experiment: &proto.Experiment{Id: "id-0"},
			},
			environmentId: "ns0",
			expectedErr:   ErrExperimentAlreadyExists,
		},
		{
			setup: func(s *experimentStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: &domain.Experiment{
				Experiment: &proto.Experiment{Id: "id-1"},
			},
			environmentId: "ns0",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		db := newExperimentStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(db)
		}
		err := db.CreateExperiment(ctx, p.input, p.environmentId)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestUpdateExperiment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup         func(*experimentStorage)
		input         *domain.Experiment
		environmentId string
		expectedErr   error
	}{
		{
			setup: func(s *experimentStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.Experiment{
				Experiment: &proto.Experiment{Id: "id-0"},
			},
			environmentId: "ns",
			expectedErr:   ErrExperimentUnexpectedAffectedRows,
		},
		{
			setup: func(s *experimentStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.Experiment{
				Experiment: &proto.Experiment{Id: "id-0"},
			},
			environmentId: "ns",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		storage := newExperimentStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		err := storage.UpdateExperiment(context.Background(), p.input, p.environmentId)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestGetExperiment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup         func(*experimentStorage)
		input         string
		environmentId string
		expected      *domain.Experiment
		expectedErr   error
	}{
		{
			setup: func(s *experimentStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input:         "",
			environmentId: "ns0",
			expected:      nil,
			expectedErr:   ErrExperimentNotFound,
		},
		{
			setup: func(s *experimentStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input:         "id-0",
			environmentId: "ns0",
			expected: &domain.Experiment{
				Experiment: &proto.Experiment{Id: "id-0"},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		storage := newExperimentStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		_, err := storage.GetExperiment(context.Background(), p.input, p.environmentId)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestListExperiments(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		setup          func(*experimentStorage)
		options        *mysql.ListOptions
		expected       []*proto.Experiment
		expectedCursor int
		expectedErr    error
	}{
		{
			setup: func(s *experimentStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			options:        nil,
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			setup: func(s *experimentStorage) {
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
				Limit:  10,
				Offset: 5,
				Filters: []*mysql.FilterV2{
					{
						Column:   "num",
						Operator: mysql.OperatorGreaterThanOrEqual,
						Value:    5,
					},
				},
				Orders: []*mysql.Order{
					{
						Column:    "id",
						Direction: mysql.OrderDirectionAsc,
					},
				},
			},
			expected:       []*proto.Experiment{},
			expectedCursor: 5,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		storage := newExperimentStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		experiments, cursor, _, err := storage.ListExperiments(
			context.Background(),
			p.options,
		)
		assert.Equal(t, p.expected, experiments)
		assert.Equal(t, p.expectedCursor, cursor)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestCountExperimentByStatus(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*experimentStorage)
		expected    *ExperimentSummary
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *experimentStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expected:    nil,
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *experimentStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expected:    &ExperimentSummary{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newExperimentStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			summary, err := storage.GetExperimentSummary(context.Background(), "ns0")
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, summary)
		})
	}
}

func newExperimentStorageWithMock(t *testing.T, mockController *gomock.Controller) *experimentStorage {
	t.Helper()
	return &experimentStorage{mock.NewMockQueryExecer(mockController)}
}
