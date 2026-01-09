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

package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/team/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/team"
)

func TestNewTeamStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewTeamStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &teamStorage{}, storage)
}

func TestUpsertTeam(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*teamStorage)
		input       *domain.Team
		expectedErr error
	}{
		{
			desc: "ErrTeamAlreadyExists",
			setup: func(s *teamStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: &domain.Team{
				Team: &proto.Team{Id: "team-id-0"},
			},
			expectedErr: mysql.ErrDuplicateEntry,
		},
		{
			desc: "Error",
			setup: func(s *teamStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input: &domain.Team{
				Team: &proto.Team{Id: "team-id-0"},
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *teamStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					insertTeamSQL,
					"team-id-0",
					"team-name-0",
					"team-description-0",
					"org-0",
					int64(1),
					int64(2),
				).Return(nil, nil)
			},
			input: &domain.Team{
				Team: &proto.Team{
					Id:             "team-id-0",
					Name:           "team-name-0",
					Description:    "team-description-0",
					OrganizationId: "org-0",
					CreatedAt:      1,
					UpdatedAt:      2,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newTeamStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpsertTeam(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetTeam(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(*teamStorage)
		id             string
		organizationId string
		expectedTeam   *domain.Team
		expectedErr    error
	}{
		{
			desc: "ErrTeamNotFound",
			setup: func(s *teamStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:             "team-id-0",
			organizationId: "org-0",
			expectedTeam:   nil,
			expectedErr:    ErrTeamNotFound,
		},
		{
			desc: "Error",
			setup: func(s *teamStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:             "team-id-0",
			organizationId: "org-0",
			expectedTeam:   nil,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *teamStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(
					gomock.Any(), // id
					gomock.Any(), // name
					gomock.Any(), // description
					gomock.Any(), // created_at
					gomock.Any(), // updated_at
					gomock.Any(), // organization_id
					gomock.Any(), // organization_name
				).Do(func(args ...interface{}) {
					*args[0].(*string) = "team-id-0"
					*args[1].(*string) = "test-team"
					*args[2].(*string) = "test-description"
					*args[3].(*int64) = int64(1)
					*args[4].(*int64) = int64(2)
					*args[5].(*string) = "org-0"
					*args[6].(*string) = "test-org"
				}).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					selectTeamSQL,
					"team-id-0",
					"org-0",
				).Return(row)
			},
			id:             "team-id-0",
			organizationId: "org-0",
			expectedTeam: &domain.Team{
				Team: &proto.Team{
					Id:               "team-id-0",
					Name:             "test-team",
					Description:      "test-description",
					CreatedAt:        1,
					UpdatedAt:        2,
					OrganizationId:   "org-0",
					OrganizationName: "test-org",
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newTeamStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			team, err := storage.GetTeam(context.Background(), p.id, p.organizationId)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.Equal(t, p.expectedTeam, team)
			}
		})
	}
}

func TestListTeams(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(*teamStorage)
		options        *mysql.ListOptions
		expectedCount  int
		expectedCursor int
		expectedErr    error
		expectedTeams  []*proto.Team
	}{
		{
			desc: "Error",
			setup: func(s *teamStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			options:        nil,
			expectedCount:  0,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
			expectedTeams:  nil,
		},
		{
			desc: "Success",
			setup: func(s *teamStorage) {
				var nextCallCount = 0
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().DoAndReturn(func() bool {
					nextCallCount++
					return nextCallCount <= 1
				}).Times(2)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(rows, nil)
				rows.EXPECT().Scan(
					gomock.Any(), // id
					gomock.Any(), // name
					gomock.Any(), // description
					gomock.Any(), // created_at
					gomock.Any(), // updated_at
					gomock.Any(), // organization_id
					gomock.Any(), // organization_name
				).Do(func(args ...interface{}) {
					*args[0].(*string) = "team-id-0"
					*args[1].(*string) = "test-team"
					*args[2].(*string) = "test-description"
					*args[3].(*int64) = int64(1)
					*args[4].(*int64) = int64(2)
					*args[5].(*string) = "org-0"
					*args[6].(*string) = "test-org"
				}).Return(nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(row)
			},
			options: &mysql.ListOptions{
				Filters: []*mysql.FilterV2{
					{
						Column:   "team.organization_id",
						Operator: mysql.OperatorEqual,
						Value:    "org-0",
					},
				},
				Orders: []*mysql.Order{
					{
						Column:    "team.name",
						Direction: mysql.OrderDirectionAsc,
					},
				},
				Offset:      0,
				Limit:       10,
				JSONFilters: nil,
				InFilters:   nil,
				NullFilters: nil,
				SearchQuery: nil,
			},
			expectedCount:  1,
			expectedCursor: 1,
			expectedErr:    nil,
			expectedTeams: []*proto.Team{
				{
					Id:               "team-id-0",
					Name:             "test-team",
					Description:      "test-description",
					CreatedAt:        1,
					UpdatedAt:        2,
					OrganizationId:   "org-0",
					OrganizationName: "test-org",
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newTeamStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			teams, cursor, _, err := storage.ListTeams(context.Background(), p.options)
			assert.Equal(t, p.expectedCount, len(teams))
			if teams != nil {
				assert.IsType(t, []*proto.Team{}, teams)
				assert.Equal(t, p.expectedTeams, teams)
			}
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetTeamByName(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(*teamStorage)
		name           string
		organizationID string
		expectedTeam   *domain.Team
		expectedErr    error
	}{
		{
			desc: "ErrTeamNotFound",
			setup: func(s *teamStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			name:           "test-team",
			organizationID: "env-0",
			expectedTeam:   nil,
			expectedErr:    ErrTeamNotFound,
		},
		{
			desc: "Error",
			setup: func(s *teamStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			name:           "test-team",
			organizationID: "env-0",
			expectedTeam:   nil,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *teamStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(
					gomock.Any(), // id
					gomock.Any(), // name
					gomock.Any(), // description
					gomock.Any(), // created_at
					gomock.Any(), // updated_at
					gomock.Any(), // organization_id
					gomock.Any(), // organization_name
				).Do(func(args ...interface{}) {
					*args[0].(*string) = "team-id-0"
					*args[1].(*string) = "test-team"
					*args[2].(*string) = "test-description"
					*args[3].(*int64) = int64(2)
					*args[4].(*int64) = int64(3)
					*args[5].(*string) = "test-org"
					*args[6].(*string) = "test-org-name"
				}).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					selectTeamByNameSQL,
					"test-team",
					"test-org",
				).Return(row)
			},
			name:           "test-team",
			organizationID: "test-org",
			expectedTeam: &domain.Team{
				Team: &proto.Team{
					Id:               "team-id-0",
					Name:             "test-team",
					Description:      "test-description",
					CreatedAt:        2,
					UpdatedAt:        3,
					OrganizationId:   "test-org",
					OrganizationName: "test-org-name",
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			storage := newTeamStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			team, err := storage.GetTeamByName(context.Background(), p.name, p.organizationID)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.Equal(t, p.expectedTeam, team)
			}
		})
	}
}

func newTeamStorageWithMock(t *testing.T, mockController *gomock.Controller) *teamStorage {
	t.Helper()
	return &teamStorage{mock.NewMockQueryExecer(mockController)}
}
