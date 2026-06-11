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

	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

func TestNewOrganizationStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewOrganizationStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &organizationStorage{}, storage)
}

func TestCreateOrganizationMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*organizationStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *organizationStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: duplicate entry",
			setup: func(s *organizationStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			expectedErr: v2es.ErrOrganizationAlreadyExists,
		},
		{
			desc: "success",
			setup: func(s *organizationStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &organizationStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			err := storage.CreateOrganization(
				context.Background(),
				&domain.Organization{Organization: &proto.Organization{}},
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateOrganizationMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*organizationStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *organizationStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: unexpected affected rows",
			setup: func(s *organizationStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: v2es.ErrOrganizationUnexpectedAffectedRows,
		},
		{
			desc: "success",
			setup: func(s *organizationStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &organizationStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			err := storage.UpdateOrganization(
				context.Background(),
				&domain.Organization{Organization: &proto.Organization{}},
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetOrganizationMySQLStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*organizationStorage)
		expectedErr error
	}{
		{
			desc: "error: ErrNoRows",
			setup: func(s *organizationStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: v2es.ErrOrganizationNotFound,
		},
		{
			desc: "error: internal",
			setup: func(s *organizationStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("internal error"))
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: errors.New("internal error"),
		},
		{
			desc: "success",
			setup: func(s *organizationStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &organizationStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			_, err := storage.GetOrganization(context.Background(), "org-id")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetSystemAdminOrganizationMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*organizationStorage)
		expectedErr error
	}{
		{
			desc: "error: ErrNoRows",
			setup: func(s *organizationStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: v2es.ErrOrganizationNotFound,
		},
		{
			desc: "error: internal",
			setup: func(s *organizationStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("internal error"))
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: errors.New("internal error"),
		},
		{
			desc: "success",
			setup: func(s *organizationStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &organizationStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			_, err := storage.GetSystemAdminOrganization(context.Background())
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListOrganizationsMySQLStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	disabled := false
	archived := false
	patterns := []struct {
		desc           string
		setup          func(*organizationStorage)
		params         v2es.ListOrganizationsParams
		expected       []*proto.Organization
		expectedCursor int
		expectedErr    error
	}{
		{
			desc:  "error: invalid order by",
			setup: nil,
			params: v2es.ListOrganizationsParams{
				OrderBy: proto.ListOrganizationsRequest_OrderBy(99),
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    v2es.ErrInvalidOrderBy,
		},
		{
			desc:  "error: invalid cursor",
			setup: nil,
			params: v2es.ListOrganizationsParams{
				PageSize: 10,
				Cursor:   "invalid",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    v2es.ErrInvalidCursor,
		},
		{
			desc: "error: query",
			setup: func(s *organizationStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			params: v2es.ListOrganizationsParams{
				PageSize: 10,
				Cursor:   "0",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "error: count",
			setup: func(s *organizationStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("count error"))
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params: v2es.ListOrganizationsParams{
				PageSize: 10,
				Cursor:   "0",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("count error"),
		},
		{
			desc: "success",
			setup: func(s *organizationStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params: v2es.ListOrganizationsParams{
				PageSize:      10,
				Cursor:        "0",
				Disabled:      &disabled,
				Archived:      &archived,
				SearchKeyword: "demo",
			},
			expected:       []*proto.Organization{},
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &organizationStorage{qe: mock.NewMockClient(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			organizations, cursor, _, err := storage.ListOrganizations(
				context.Background(),
				p.params,
			)
			assert.Equal(t, p.expected, organizations)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
