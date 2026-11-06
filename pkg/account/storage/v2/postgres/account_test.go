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

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	pgmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

var errInternal = errors.New("internal error")

func TestNewAccountStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewAccountStorage(pgmock.NewMockQueryExecer(mockController))
	assert.IsType(t, &accountStorage{}, storage)
}

func TestCreateAccountV2Postgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		account     *domain.AccountV2
		expectedErr error
	}{
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			account:     &domain.AccountV2{AccountV2: &proto.AccountV2{}},
			expectedErr: errInternal,
		},
		{
			desc: "ErrAccountAlreadyExists",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, postgres.ErrDuplicateEntry)
			},
			account:     &domain.AccountV2{AccountV2: &proto.AccountV2{}},
			expectedErr: v2as.ErrAccountAlreadyExists,
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			account:     &domain.AccountV2{AccountV2: &proto.AccountV2{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateAccountV2(context.Background(), p.account)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAccountV2Postgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		account     *domain.AccountV2
		expectedErr error
	}{
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			account:     &domain.AccountV2{AccountV2: &proto.AccountV2{}},
			expectedErr: errInternal,
		},
		{
			desc: "ErrUnexpectedAffectedRows",
			setup: func(s *accountStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			account:     &domain.AccountV2{AccountV2: &proto.AccountV2{}},
			expectedErr: v2as.ErrAccountUnexpectedAffectedRows,
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			account:     &domain.AccountV2{AccountV2: &proto.AccountV2{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpdateAccountV2(context.Background(), p.account)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteAccountV2Postgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		account     *domain.AccountV2
		expectedErr error
	}{
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			account:     &domain.AccountV2{AccountV2: &proto.AccountV2{}},
			expectedErr: errInternal,
		},
		{
			desc: "ErrUnexpectedAffectedRows",
			setup: func(s *accountStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			account:     &domain.AccountV2{AccountV2: &proto.AccountV2{}},
			expectedErr: v2as.ErrAccountUnexpectedAffectedRows,
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			account:     &domain.AccountV2{AccountV2: &proto.AccountV2{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.DeleteAccountV2(context.Background(), p.account)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAccountV2Postgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		email       string
		orgID       string
		expectedErr error
	}{
		{
			desc: "ErrAccountNotFound",
			setup: func(s *accountStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(postgres.ErrNoRows)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			email:       "test@example.com",
			orgID:       "org-1",
			expectedErr: v2as.ErrAccountNotFound,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errInternal)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			email:       "test@example.com",
			orgID:       "org-1",
			expectedErr: errInternal,
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			email:       "test@example.com",
			orgID:       "org-1",
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetAccountV2(context.Background(), p.email, p.orgID)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAccountV2ByEnvironmentIDPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc          string
		setup         func(*accountStorage)
		email         string
		environmentID string
		expectedErr   error
	}{
		{
			desc: "ErrAccountNotFound",
			setup: func(s *accountStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(postgres.ErrNoRows)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			email:         "test@example.com",
			environmentID: "env-1",
			expectedErr:   v2as.ErrAccountNotFound,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errInternal)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			email:         "test@example.com",
			environmentID: "env-1",
			expectedErr:   errInternal,
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			email:         "test@example.com",
			environmentID: "env-1",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetAccountV2ByEnvironmentID(context.Background(), p.email, p.environmentID)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAvatarAccountsV2Postgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		params      v2as.GetAvatarAccountsV2Params
		expected    []*proto.AccountV2
		expectedErr error
	}{
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			params: v2as.GetAvatarAccountsV2Params{
				Emails:        []string{"test@example.com"},
				EnvironmentID: "env-1",
			},
			expected:    nil,
			expectedErr: errInternal,
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			params: v2as.GetAvatarAccountsV2Params{
				Emails:        []string{"test@example.com"},
				EnvironmentID: "env-1",
			},
			expected:    []*proto.AccountV2{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			accounts, err := storage.GetAvatarAccountsV2(context.Background(), p.params)
			assert.Equal(t, p.expected, accounts)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAccountsV2Postgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc           string
		setup          func(*accountStorage)
		params         v2as.ListAccountsV2Params
		expected       []*proto.AccountV2
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "QueryError",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			params: v2as.ListAccountsV2Params{
				PageSize:       10,
				Cursor:         "0",
				OrganizationID: "org-1",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errInternal,
		},
		{
			desc: "CountError",
			setup: func(s *accountStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errInternal)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params: v2as.ListAccountsV2Params{
				PageSize:       10,
				Cursor:         "0",
				OrganizationID: "org-1",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errInternal,
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
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
			params: v2as.ListAccountsV2Params{
				PageSize:       10,
				Cursor:         "0",
				OrganizationID: "org-1",
			},
			expected:       []*proto.AccountV2{},
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			accounts, cursor, _, err := storage.ListAccountsV2(
				context.Background(),
				p.params,
			)
			assert.Equal(t, p.expected, accounts)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newAccountStorageWithMock(t *testing.T, mockController *gomock.Controller) *accountStorage {
	t.Helper()
	return &accountStorage{qe: pgmock.NewMockQueryExecer(mockController)}
}

// TestListAccountsV2OptionsEnvironmentRoleFilters pins the JSON containment
// candidates built for the environment_roles column: a member filter matches
// the environment only unless a role is requested, and a nil role is never
// serialized as null, which would match nothing.
func TestListAccountsV2OptionsEnvironmentRoleFilters(t *testing.T) {
	t.Parallel()

	jsonValues := func(parts []postgres.WherePart) []string {
		values := make([]string, 0, len(parts))
		for _, part := range parts {
			if f, ok := part.(*postgres.JSONFilter); ok {
				values = append(values, f.Values[0].(string))
			}
		}
		return values
	}
	cases := []struct {
		desc         string
		params       v2as.ListAccountsV2Params
		wantMember   []string // candidates OR-ed for a member filter
		wantAdminEnv string   // candidate of the admin environment filter
	}{
		{
			desc: "member filter matches the environment with any role",
			params: v2as.ListAccountsV2Params{
				OrganizationID: "org-1", PageSize: 10, Cursor: "0",
				EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
					{EnvironmentId: "env-a", Role: proto.AccountV2_Role_Environment_UNASSIGNED},
					{EnvironmentId: "env-b", Role: proto.AccountV2_Role_Environment_UNASSIGNED},
				},
			},
			wantMember: []string{`{"environment_id":"env-a"}`, `{"environment_id":"env-b"}`},
		},
		{
			desc: "member filter with a requested role matches that role",
			params: v2as.ListAccountsV2Params{
				OrganizationID: "org-1", PageSize: 10, Cursor: "0",
				EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
					{EnvironmentId: "env-a", Role: proto.AccountV2_Role_Environment_EDITOR},
				},
			},
			wantMember: []string{`{"environment_id":"env-a","role":2}`},
		},
		{
			desc: "admin environment filter without a role matches any role",
			params: v2as.ListAccountsV2Params{
				OrganizationID: "org-1", PageSize: 10, Cursor: "0",
				EnvironmentID: ptr("env-a"),
			},
			wantAdminEnv: `{"environment_id":"env-a"}`,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			options, err := listAccountsV2OptionsFromParams(tc.params)
			assert.NoError(t, err)
			if tc.wantMember != nil {
				if assert.Len(t, options.OrFilters, 1) {
					assert.Equal(t, tc.wantMember, jsonValues(options.OrFilters[0].Queries))
				}
				return
			}
			if assert.Len(t, options.JSONFilters, 1) {
				assert.Equal(t, tc.wantAdminEnv, options.JSONFilters[0].Values[0])
			}
		})
	}
}

func ptr[T any](v T) *T { return &v }
