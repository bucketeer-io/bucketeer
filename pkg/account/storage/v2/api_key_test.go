// Copyright 2024 The Bucketeer Authors.
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

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/account"
)

func TestCreateAPIKey(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc                 string
		setup                func(*accountStorage)
		input                *domain.APIKey
		environmentNamespace string
		expectedErr          error
	}{
		{
			desc: "ErrAPIKeyAlreadyExists",
			setup: func(s *accountStorage) {
				s.client.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: &domain.APIKey{
				APIKey: &proto.APIKey{Id: "aid-0"},
			},
			environmentNamespace: "ns0",
			expectedErr:          ErrAPIKeyAlreadyExists,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.client.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input: &domain.APIKey{
				APIKey: &proto.APIKey{Id: "aid-0"},
			},
			environmentNamespace: "ns0",
			expectedErr:          errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				s.client.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(),
					gomock.Regex("^INSERT INTO api_key[\\s\\S]\\s*?id,\\s*?name,\\s*?role,\\s*?disabled,\\s*?created_at,\\s*?updated_at,\\s*?environment_namespace[\\s\\S]*?(\\?[\\s\\S]*?){7}"),
					"aid-0", "name", int32(0), false, int64(2), int64(3), "ns0",
				).Return(nil, nil).Times(1)
			},
			input: &domain.APIKey{
				APIKey: &proto.APIKey{Id: "aid-0", Name: "name", Role: 0, Disabled: false, CreatedAt: 2, UpdatedAt: 3},
			},
			environmentNamespace: "ns0",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateAPIKey(context.Background(), p.input, p.environmentNamespace)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAPIKey(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	id := "aid-0"
	environmentNamespace := "ns0"
	name := "name"
	role := proto.APIKey_Role(0)
	disabled := false
	createdAt := int64(2)
	updatedAt := int64(3)

	patterns := []struct {
		desc                 string
		setup                func(*accountStorage)
		whereParts           []mysql.WherePart
		input                *domain.APIKey
		environmentNamespace string
		expectedErr          error
	}{
		{
			desc: "ErrAPIKeyUnexpectedAffectedRows",
			setup: func(s *accountStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.client.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.APIKey{
				APIKey: &proto.APIKey{Id: id},
			},
			environmentNamespace: environmentNamespace,
			expectedErr:          ErrAPIKeyUnexpectedAffectedRows,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.client.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input: &domain.APIKey{
				APIKey: &proto.APIKey{Id: id},
			},
			environmentNamespace: environmentNamespace,
			expectedErr:          errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil).Times(1)
				s.client.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(),
					"UPDATE api_key SET name = ?, role = ?, disabled = ?, created_at = ?, updated_at = ? WHERE id = ? AND environment_namespace = ?",
					[]interface{}{name, int32(role), disabled, createdAt, updatedAt}, []interface{}{id, environmentNamespace},
				).Return(result, nil).Times(1)
			},
			input: &domain.APIKey{
				APIKey: &proto.APIKey{Id: id, Name: name, Role: role, Disabled: disabled, CreatedAt: createdAt, UpdatedAt: updatedAt},
			},
			whereParts: []mysql.WherePart{
				mysql.NewFilter("id", "=", id),
				mysql.NewFilter("environment_namespace", "=", environmentNamespace),
			},
			environmentNamespace: environmentNamespace,
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpdateAPIKey(context.Background(), p.input, p.whereParts)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAPIKey(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc                 string
		setup                func(*accountStorage)
		id                   string
		environmentNamespace string
		expectedErr          error
	}{
		{
			desc: "ErrAPIKeyNotFound",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.client.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:                   "id-0",
			environmentNamespace: "ns0",
			expectedErr:          ErrAPIKeyNotFound,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.client.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)

			},
			id:                   "id-0",
			environmentNamespace: "ns0",
			expectedErr:          errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil).Times(1)
				s.client.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Any(),
					"id-0", "ns0",
				).Return(row).Times(1)
			},
			id:                   "id-0",
			environmentNamespace: "ns0",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetAPIKey(context.Background(), p.id, p.environmentNamespace)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

var (
// nextCount = 0
)

func TestListAPIKeys(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc           string
		setup          func(*accountStorage)
		whereParts     []mysql.WherePart
		orders         []*mysql.Order
		limit          int
		offset         int
		expected       []*proto.APIKey
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.client.(*mock.MockClient).EXPECT().QueryContext(
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
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				var nextCallCount = 0
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil).Times(1)
				rows.EXPECT().Next().DoAndReturn(func() bool {
					if nextCallCount < 2 {
						nextCallCount++
						return true
					} else {
						return false
					}
				}).Times(3)
				rows.EXPECT().Scan(gomock.Any()).Return(nil).Times(2)
				rows.EXPECT().Err().Return(nil)
				s.client.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(),
					gomock.Regex("SELECT[\\s\\S]\\s*?id,\\s*?name,\\s*?role,\\s*?disabled,\\s*?created_at,\\s*?updated_at\\s*?FROM\\s*?api_key[\\s\\S]*?$"),
					5,
				).Return(rows, nil).Times(1)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil).Times(1)
				s.client.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row).Times(1)
			},
			whereParts: []mysql.WherePart{
				mysql.NewFilter("num", ">=", 5),
			},
			orders: []*mysql.Order{
				mysql.NewOrder("id", mysql.OrderDirectionAsc),
			},
			limit:  10,
			offset: 5,
			expected: []*proto.APIKey{
				{Id: "", Name: "", Role: 0, Disabled: false, CreatedAt: 0, UpdatedAt: 0},
				{Id: "", Name: "", Role: 0, Disabled: false, CreatedAt: 0, UpdatedAt: 0},
			},
			expectedCursor: 7,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			apiKeys, cursor, _, err := storage.ListAPIKeys(
				context.Background(),
				p.whereParts,
				p.orders,
				p.limit,
				p.offset,
			)
			assert.Equal(t, p.expected, apiKeys)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
