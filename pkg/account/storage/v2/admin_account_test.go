// Copyright 2022 The Bucketeer Authors.
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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/account"
)

func TestNewAdminAccountStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewAdminAccountStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &adminAccountStorage{}, storage)
}

func TestCreateAdminAccount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := map[string]struct {
		setup       func(*adminAccountStorage)
		input       *domain.Account
		expectedErr error
	}{
		"ErrAdminAccountAlreadyExists": {
			setup: func(s *adminAccountStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: &domain.Account{
				Account: &proto.Account{Id: "aid-0"},
			},
			expectedErr: ErrAdminAccountAlreadyExists,
		},
		"Error": {
			setup: func(s *adminAccountStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input: &domain.Account{
				Account: &proto.Account{Id: "aid-0"},
			},
			expectedErr: errors.New("error"),
		},
		"Success": {
			setup: func(s *adminAccountStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: &domain.Account{
				Account: &proto.Account{Id: "aid-0"},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			storage := newAdminAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateAdminAccount(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAdminAccount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := map[string]struct {
		setup       func(*adminAccountStorage)
		input       *domain.Account
		expectedErr error
	}{
		"ErrAdminAccountUnexpectedAffectedRows": {
			setup: func(s *adminAccountStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.Account{
				Account: &proto.Account{Id: "aid-0"},
			},
			expectedErr: ErrAdminAccountUnexpectedAffectedRows,
		},
		"Error": {
			setup: func(s *adminAccountStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input: &domain.Account{
				Account: &proto.Account{Id: "aid-0"},
			},
			expectedErr: errors.New("error"),
		},
		"Success": {
			setup: func(s *adminAccountStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.Account{
				Account: &proto.Account{Id: "aid-0"},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			storage := newAdminAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpdateAdminAccount(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAdminAccount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := map[string]struct {
		setup       func(*adminAccountStorage)
		id          string
		expectedErr error
	}{
		"ErrAdminAccountNotFound": {
			setup: func(s *adminAccountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "id-0",
			expectedErr: ErrAdminAccountNotFound,
		},
		"Error": {
			setup: func(s *adminAccountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)

			},
			id:          "id-0",
			expectedErr: errors.New("error"),
		},
		"Success": {
			setup: func(s *adminAccountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "id-0",
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			storage := newAdminAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetAdminAccount(context.Background(), p.id)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAdminAccounts(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := map[string]struct {
		setup          func(*adminAccountStorage)
		whereParts     []mysql.WherePart
		orders         []*mysql.Order
		limit          int
		offset         int
		expected       []*proto.Account
		expectedCursor int
		expectedErr    error
	}{
		"Error": {
			setup: func(s *adminAccountStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			whereParts:     nil,
			orders:         nil,
			limit:          0,
			offset:         0,
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		"Success": {
			setup: func(s *adminAccountStorage) {
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
			whereParts: []mysql.WherePart{
				mysql.NewFilter("num", ">=", 5),
			},
			orders: []*mysql.Order{
				mysql.NewOrder("id", mysql.OrderDirectionAsc),
			},
			limit:          10,
			offset:         5,
			expected:       []*proto.Account{},
			expectedCursor: 5,
			expectedErr:    nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			storage := newAdminAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			accounts, cursor, _, err := storage.ListAdminAccounts(
				context.Background(),
				p.whereParts,
				p.orders,
				p.limit,
				p.offset,
			)
			assert.Equal(t, p.expected, accounts)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newAdminAccountStorageWithMock(t *testing.T, mockController *gomock.Controller) *adminAccountStorage {
	t.Helper()
	return &adminAccountStorage{mock.NewMockQueryExecer(mockController)}
}
