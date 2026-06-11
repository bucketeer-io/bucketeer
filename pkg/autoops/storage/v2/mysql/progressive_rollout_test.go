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

package mysql

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/storage/v2"
	mysqlstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

func TestNewProgressiveRolloutStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	s := NewProgressiveRolloutStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &progressiveRolloutStorage{}, s)
}

func TestCreateProgressiveRolloutMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*progressiveRolloutStorage)
		expectedErr error
	}{
		{
			desc: "error: duplicate entry",
			setup: func(s *progressiveRolloutStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysqlstorage.ErrDuplicateEntry)
			},
			expectedErr: v2as.ErrProgressiveRolloutAlreadyExists,
		},
		{
			desc: "error: internal",
			setup: func(s *progressiveRolloutStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *progressiveRolloutStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &progressiveRolloutStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(s)
			err := s.CreateProgressiveRollout(context.Background(),
				&domain.ProgressiveRollout{ProgressiveRollout: &autoopsproto.ProgressiveRollout{}}, "ns0")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateProgressiveRolloutMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*progressiveRolloutStorage)
		expectedErr error
	}{
		{
			desc: "error: internal",
			setup: func(s *progressiveRolloutStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: unexpected affected rows",
			setup: func(s *progressiveRolloutStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: v2as.ErrProgressiveRolloutUnexpectedAffectedRows,
		},
		{
			desc: "success",
			setup: func(s *progressiveRolloutStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &progressiveRolloutStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(s)
			err := s.UpdateProgressiveRollout(context.Background(),
				&domain.ProgressiveRollout{ProgressiveRollout: &autoopsproto.ProgressiveRollout{}}, "ns0")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetProgressiveRolloutMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*progressiveRolloutStorage)
		expectedErr error
	}{
		{
			desc: "error: not found",
			setup: func(s *progressiveRolloutStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysqlstorage.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: v2as.ErrProgressiveRolloutNotFound,
		},
		{
			desc: "error: internal",
			setup: func(s *progressiveRolloutStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("internal"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: errors.New("internal"),
		},
		{
			desc: "success",
			setup: func(s *progressiveRolloutStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &progressiveRolloutStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(s)
			_, err := s.GetProgressiveRollout(context.Background(), "id", "ns0")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteProgressiveRolloutMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*progressiveRolloutStorage)
		expectedErr error
	}{
		{
			desc: "error: internal",
			setup: func(s *progressiveRolloutStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: unexpected affected rows",
			setup: func(s *progressiveRolloutStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: v2as.ErrProgressiveRolloutUnexpectedAffectedRows,
		},
		{
			desc: "success",
			setup: func(s *progressiveRolloutStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &progressiveRolloutStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(s)
			err := s.DeleteProgressiveRollout(context.Background(), "id", "ns0")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListProgressiveRolloutsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(*progressiveRolloutStorage)
		params         v2as.ListProgressiveRolloutsParams
		expected       []*autoopsproto.ProgressiveRollout
		expectedTotal  int64
		expectedCursor int
		expectedErr    error
	}{
		{
			desc:  "error: invalid order by",
			setup: nil,
			params: v2as.ListProgressiveRolloutsParams{
				OrderBy: autoopsproto.ListProgressiveRolloutsRequest_OrderBy(999),
			},
			expectedErr: v2as.ErrInvalidOrderBy,
		},
		{
			desc:        "error: invalid cursor",
			setup:       nil,
			params:      v2as.ListProgressiveRolloutsParams{Cursor: "invalid"},
			expectedErr: v2as.ErrInvalidCursor,
		},
		{
			desc:        "error: negative cursor",
			setup:       nil,
			params:      v2as.ListProgressiveRolloutsParams{Cursor: "-1"},
			expectedErr: v2as.ErrInvalidCursor,
		},
		{
			desc: "error: query",
			setup: func(s *progressiveRolloutStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			params:      v2as.ListProgressiveRolloutsParams{EnvironmentID: "ns0", Cursor: "0"},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: count",
			setup: func(s *progressiveRolloutStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("count error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params:      v2as.ListProgressiveRolloutsParams{EnvironmentID: "ns0", Cursor: "0"},
			expectedErr: errors.New("count error"),
		},
		{
			desc: "success",
			setup: func(s *progressiveRolloutStorage) {
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
			params: v2as.ListProgressiveRolloutsParams{
				EnvironmentID: "ns0",
				FeatureIDs:    []string{"f1"},
				Cursor:        "0",
				PageSize:      10,
			},
			expected:       []*autoopsproto.ProgressiveRollout{},
			expectedTotal:  0,
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &progressiveRolloutStorage{qe: mock.NewMockQueryExecer(mockController)}
			if p.setup != nil {
				p.setup(s)
			}
			rollouts, total, cursor, err := s.ListProgressiveRollouts(context.Background(), p.params)
			assert.Equal(t, p.expected, rollouts)
			assert.Equal(t, p.expectedTotal, total)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
