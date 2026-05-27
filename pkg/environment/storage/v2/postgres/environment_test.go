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

	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	pgmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

var errInternalEnv = errors.New("internal error")

func TestNewEnvironmentStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewEnvironmentStorage(pgmock.NewMockQueryExecer(mockController))
	assert.IsType(t, &environmentStorage{}, storage)
}

func TestCreateEnvironmentV2Postgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*environmentStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *environmentStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternalEnv)
			},
			expectedErr: errInternalEnv,
		},
		{
			desc: "error: duplicate entry",
			setup: func(s *environmentStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, postgres.ErrDuplicateEntry)
			},
			expectedErr: v2es.ErrEnvironmentAlreadyExists,
		},
		{
			desc: "success",
			setup: func(s *environmentStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newEnvironmentStorageWithMock(t, mockController)
			p.setup(storage)
			err := storage.CreateEnvironmentV2(
				context.Background(),
				&domain.EnvironmentV2{EnvironmentV2: &proto.EnvironmentV2{}},
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateEnvironmentV2Postgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*environmentStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *environmentStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternalEnv)
			},
			expectedErr: errInternalEnv,
		},
		{
			desc: "error: unexpected affected rows",
			setup: func(s *environmentStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: v2es.ErrEnvironmentUnexpectedAffectedRows,
		},
		{
			desc: "success",
			setup: func(s *environmentStorage) {
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
			storage := newEnvironmentStorageWithMock(t, mockController)
			p.setup(storage)
			err := storage.UpdateEnvironmentV2(
				context.Background(),
				&domain.EnvironmentV2{EnvironmentV2: &proto.EnvironmentV2{}},
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetEnvironmentV2Postgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*environmentStorage)
		expectedErr error
	}{
		{
			desc: "error: ErrNoRows",
			setup: func(s *environmentStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(postgres.ErrNoRows)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: v2es.ErrEnvironmentNotFound,
		},
		{
			desc: "error: internal",
			setup: func(s *environmentStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errInternalEnv)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: errInternalEnv,
		},
		{
			desc: "success",
			setup: func(s *environmentStorage) {
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
			storage := newEnvironmentStorageWithMock(t, mockController)
			p.setup(storage)
			_, err := storage.GetEnvironmentV2(context.Background(), "env-id")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListEnvironmentsV2Postgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(*environmentStorage)
		params         v2es.ListEnvironmentsV2Params
		expected       []*proto.EnvironmentV2
		expectedCursor int
		expectedErr    error
	}{
		{
			desc:  "error: invalid order by",
			setup: nil,
			params: v2es.ListEnvironmentsV2Params{
				OrderBy: proto.ListEnvironmentsV2Request_OrderBy(99),
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    v2es.ErrInvalidOrderBy,
		},
		{
			desc:  "error: invalid cursor",
			setup: nil,
			params: v2es.ListEnvironmentsV2Params{
				PageSize: 10,
				Cursor:   "invalid",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    v2es.ErrInvalidCursor,
		},
		{
			desc: "error: query",
			setup: func(s *environmentStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternalEnv)
			},
			params: v2es.ListEnvironmentsV2Params{
				PageSize: 10,
				Cursor:   "0",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errInternalEnv,
		},
		{
			desc: "error: count",
			setup: func(s *environmentStorage) {
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
			params: v2es.ListEnvironmentsV2Params{
				PageSize: 10,
				Cursor:   "0",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("count error"),
		},
		{
			desc: "success",
			setup: func(s *environmentStorage) {
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
			params: v2es.ListEnvironmentsV2Params{
				PageSize:       10,
				Cursor:         "0",
				ProjectID:      "project-id",
				OrganizationID: "org-id",
				SearchKeyword:  "demo",
			},
			expected:       []*proto.EnvironmentV2{},
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newEnvironmentStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			environments, cursor, _, err := storage.ListEnvironmentsV2(
				context.Background(),
				p.params,
			)
			assert.Equal(t, p.expected, environments)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAutoArchiveEnabledEnvironmentsPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*environmentStorage)
		expected    []*domain.EnvironmentV2
		expectedErr error
	}{
		{
			desc: "error: query",
			setup: func(s *environmentStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(),
				).Return(nil, errInternalEnv)
			},
			expected:    nil,
			expectedErr: errInternalEnv,
		},
		{
			desc: "success",
			setup: func(s *environmentStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expected:    []*domain.EnvironmentV2{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newEnvironmentStorageWithMock(t, mockController)
			p.setup(storage)
			environments, err := storage.ListAutoArchiveEnabledEnvironments(context.Background())
			assert.Equal(t, p.expected, environments)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newEnvironmentStorageWithMock(t *testing.T, mockController *gomock.Controller) *environmentStorage {
	t.Helper()
	return &environmentStorage{qe: pgmock.NewMockQueryExecer(mockController)}
}
