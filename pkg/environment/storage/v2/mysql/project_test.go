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

func TestNewProjectStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewProjectStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &projectStorage{}, storage)
}

func TestCreateProjectMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*projectStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *projectStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: duplicate entry",
			setup: func(s *projectStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			expectedErr: v2es.ErrProjectAlreadyExists,
		},
		{
			desc: "success",
			setup: func(s *projectStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &projectStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			err := storage.CreateProject(
				context.Background(),
				&domain.Project{Project: &proto.Project{}},
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateProjectMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*projectStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *projectStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: unexpected affected rows",
			setup: func(s *projectStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: v2es.ErrProjectUnexpectedAffectedRows,
		},
		{
			desc: "success",
			setup: func(s *projectStorage) {
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
			storage := &projectStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			err := storage.UpdateProject(
				context.Background(),
				&domain.Project{Project: &proto.Project{}},
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetProjectMySQLStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*projectStorage)
		expectedErr error
	}{
		{
			desc: "error: ErrNoRows",
			setup: func(s *projectStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: v2es.ErrProjectNotFound,
		},
		{
			desc: "error: internal",
			setup: func(s *projectStorage) {
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
			setup: func(s *projectStorage) {
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
			storage := &projectStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			_, err := storage.GetProject(context.Background(), "project-id")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetTrialProjectByEmailMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*projectStorage)
		expectedErr error
	}{
		{
			desc: "error: ErrNoRows",
			setup: func(s *projectStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: v2es.ErrProjectNotFound,
		},
		{
			desc: "error: internal",
			setup: func(s *projectStorage) {
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
			setup: func(s *projectStorage) {
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
			storage := &projectStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			_, err := storage.GetTrialProjectByEmail(
				context.Background(), "test@example.com", false, true,
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListProjectsMySQLStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	disabled := false
	patterns := []struct {
		desc           string
		setup          func(*projectStorage)
		params         v2es.ListProjectsParams
		expected       []*proto.Project
		expectedCursor int
		expectedErr    error
	}{
		{
			desc:  "error: invalid order by",
			setup: nil,
			params: v2es.ListProjectsParams{
				OrderBy: proto.ListProjectsV2Request_OrderBy(99),
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    v2es.ErrInvalidOrderBy,
		},
		{
			desc:  "error: invalid cursor",
			setup: nil,
			params: v2es.ListProjectsParams{
				PageSize: 10,
				Cursor:   "invalid",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    v2es.ErrInvalidCursor,
		},
		{
			desc: "error: query",
			setup: func(s *projectStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			params: v2es.ListProjectsParams{
				PageSize: 10,
				Cursor:   "0",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "error: count",
			setup: func(s *projectStorage) {
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
			params: v2es.ListProjectsParams{
				PageSize: 10,
				Cursor:   "0",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("count error"),
		},
		{
			desc: "success",
			setup: func(s *projectStorage) {
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
			params: v2es.ListProjectsParams{
				PageSize:        10,
				Cursor:          "0",
				OrganizationID:  "org-id",
				OrganizationIDs: []string{"org-1", "org-2"},
				Disabled:        &disabled,
				SearchKeyword:   "demo",
			},
			expected:       []*proto.Project{},
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &projectStorage{qe: mock.NewMockClient(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			projects, cursor, _, err := storage.ListProjects(
				context.Background(),
				p.params,
			)
			assert.Equal(t, p.expected, projects)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
