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

package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/storage/v2"
	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	pgmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

var errInternalPR = errors.New("internal")

func TestNewProgressiveRolloutStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	s := NewProgressiveRolloutStorage(pgmock.NewMockQueryExecer(mockController))
	assert.IsType(t, &progressiveRolloutStorage{}, s)
}

func TestCreateProgressiveRolloutPostgres(t *testing.T) {
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
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, pgstorage.ErrDuplicateEntry)
			},
			expectedErr: v2as.ErrProgressiveRolloutAlreadyExists,
		},
		{
			desc: "error: internal",
			setup: func(s *progressiveRolloutStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternalPR)
			},
			expectedErr: errInternalPR,
		},
		{
			desc: "success",
			setup: func(s *progressiveRolloutStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newProgressiveRolloutStorageWithMock(t, mockController)
			p.setup(s)
			err := s.CreateProgressiveRollout(context.Background(),
				&domain.ProgressiveRollout{ProgressiveRollout: &autoopsproto.ProgressiveRollout{}}, "ns0")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateProgressiveRolloutPostgres(t *testing.T) {
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
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternalPR)
			},
			expectedErr: errInternalPR,
		},
		{
			desc: "error: unexpected affected rows",
			setup: func(s *progressiveRolloutStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: v2as.ErrProgressiveRolloutUnexpectedAffectedRows,
		},
		{
			desc: "success",
			setup: func(s *progressiveRolloutStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newProgressiveRolloutStorageWithMock(t, mockController)
			p.setup(s)
			err := s.UpdateProgressiveRollout(context.Background(),
				&domain.ProgressiveRollout{ProgressiveRollout: &autoopsproto.ProgressiveRollout{}}, "ns0")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetProgressiveRolloutPostgres(t *testing.T) {
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
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(pgstorage.ErrNoRows)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: v2as.ErrProgressiveRolloutNotFound,
		},
		{
			desc: "error: internal",
			setup: func(s *progressiveRolloutStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errInternalPR)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: errInternalPR,
		},
		{
			desc: "success",
			setup: func(s *progressiveRolloutStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newProgressiveRolloutStorageWithMock(t, mockController)
			p.setup(s)
			_, err := s.GetProgressiveRollout(context.Background(), "id", "ns0")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteProgressiveRolloutPostgres(t *testing.T) {
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
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternalPR)
			},
			expectedErr: errInternalPR,
		},
		{
			desc: "error: unexpected affected rows",
			setup: func(s *progressiveRolloutStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: v2as.ErrProgressiveRolloutUnexpectedAffectedRows,
		},
		{
			desc: "success",
			setup: func(s *progressiveRolloutStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newProgressiveRolloutStorageWithMock(t, mockController)
			p.setup(s)
			err := s.DeleteProgressiveRollout(context.Background(), "id", "ns0")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListProgressiveRolloutsPostgres(t *testing.T) {
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
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternalPR)
			},
			params:      v2as.ListProgressiveRolloutsParams{EnvironmentID: "ns0", Cursor: "0"},
			expectedErr: errInternalPR,
		},
		{
			desc: "error: count",
			setup: func(s *progressiveRolloutStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("count error"))
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params:      v2as.ListProgressiveRolloutsParams{EnvironmentID: "ns0", Cursor: "0"},
			expectedErr: errors.New("count error"),
		},
		{
			desc: "success",
			setup: func(s *progressiveRolloutStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
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
			s := newProgressiveRolloutStorageWithMock(t, mockController)
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

func newProgressiveRolloutStorageWithMock(t *testing.T, mockController *gomock.Controller) *progressiveRolloutStorage {
	t.Helper()
	return &progressiveRolloutStorage{qe: pgmock.NewMockQueryExecer(mockController)}
}
