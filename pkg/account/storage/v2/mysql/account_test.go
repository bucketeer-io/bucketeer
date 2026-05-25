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

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

func TestNewAccountStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewAccountStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &accountStorage{}, storage)
}

func TestCreateAccountV2MySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: duplicate entry",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			expectedErr: v2as.ErrAccountAlreadyExists,
		},
		{
			desc: "success",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &accountStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			err := storage.CreateAccountV2(
				context.Background(),
				&domain.AccountV2{AccountV2: &proto.AccountV2{}},
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAccountV2MySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: unexpected affected rows",
			setup: func(s *accountStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: v2as.ErrAccountUnexpectedAffectedRows,
		},
		{
			desc: "success",
			setup: func(s *accountStorage) {
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
			storage := &accountStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			err := storage.UpdateAccountV2(
				context.Background(),
				&domain.AccountV2{AccountV2: &proto.AccountV2{}},
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteAccountV2MySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: unexpected affected rows",
			setup: func(s *accountStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: v2as.ErrAccountUnexpectedAffectedRows,
		},
		{
			desc: "success",
			setup: func(s *accountStorage) {
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
			storage := &accountStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			err := storage.DeleteAccountV2(
				context.Background(),
				&domain.AccountV2{AccountV2: &proto.AccountV2{}},
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAccountV2MySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		expectedErr error
	}{
		{
			desc: "error: ErrNoRows",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: v2as.ErrAccountNotFound,
		},
		{
			desc: "error: internal",
			setup: func(s *accountStorage) {
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
			setup: func(s *accountStorage) {
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
			storage := &accountStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			_, err := storage.GetAccountV2(context.Background(), "test@example.com", "org-1")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAccountV2ByEnvironmentIDMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		expectedErr error
	}{
		{
			desc: "error: ErrNoRows",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: v2as.ErrAccountNotFound,
		},
		{
			desc: "error: internal",
			setup: func(s *accountStorage) {
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
			setup: func(s *accountStorage) {
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
			storage := &accountStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			_, err := storage.GetAccountV2ByEnvironmentID(
				context.Background(), "test@example.com", "env-1",
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAvatarAccountsV2MySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		expected    []*proto.AccountV2
		expectedErr error
	}{
		{
			desc: "error: query",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expected:    nil,
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *accountStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expected:    []*proto.AccountV2{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &accountStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			accounts, err := storage.GetAvatarAccountsV2(
				context.Background(),
				v2as.GetAvatarAccountsV2Params{
					Emails:        []string{"test@example.com"},
					EnvironmentID: "env-1",
				},
			)
			assert.Equal(t, p.expected, accounts)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAccountsV2MySQL(t *testing.T) {
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
			desc: "error: query",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			params: v2as.ListAccountsV2Params{
				PageSize:       10,
				Cursor:         "0",
				OrganizationID: "org-1",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "error: count",
			setup: func(s *accountStorage) {
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
			params: v2as.ListAccountsV2Params{
				PageSize:       10,
				Cursor:         "0",
				OrganizationID: "org-1",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("count error"),
		},
		{
			desc: "success",
			setup: func(s *accountStorage) {
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
			storage := &accountStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
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
