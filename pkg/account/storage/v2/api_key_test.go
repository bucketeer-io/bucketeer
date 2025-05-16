// Copyright 2025 The Bucketeer Authors.
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
		desc          string
		setup         func(*accountStorage)
		input         *domain.APIKey
		environmentId string
		expectedErr   error
	}{
		{
			desc: "ErrAPIKeyAlreadyExists",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: &domain.APIKey{
				APIKey: &proto.APIKey{Id: "aid-0"},
			},
			environmentId: "ns0",
			expectedErr:   ErrAPIKeyAlreadyExists,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input: &domain.APIKey{
				APIKey: &proto.APIKey{Id: "aid-0"},
			},
			environmentId: "ns0",
			expectedErr:   errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					gomock.Any(),
					"aid-0",
					"name",
					int32(0),
					false,
					int64(2),
					int64(3),
					"ns0",
					"aid-0",
					"demo@bucketeer.io",
					"test",
				).Return(nil, nil)
			},
			input: &domain.APIKey{
				APIKey: &proto.APIKey{
					Id:          "aid-0",
					Name:        "name",
					Role:        0,
					Disabled:    false,
					Maintainer:  "demo@bucketeer.io",
					ApiKey:      "aid-0",
					Description: "test",
					CreatedAt:   2,
					UpdatedAt:   3,
				},
			},
			environmentId: "ns0",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateAPIKey(context.Background(), p.input, p.environmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAPIKey(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	id := "aid-0"
	environmentId := "ns0"
	name := "name"
	role := proto.APIKey_Role(0)
	disabled := false
	description := "test"
	createdAt := int64(2)
	updatedAt := int64(3)
	maintainer := "demo@bucketeer.io"

	patterns := []struct {
		desc          string
		setup         func(*accountStorage)
		input         *domain.APIKey
		environmentId string
		expectedErr   error
	}{
		{
			desc: "ErrAPIKeyUnexpectedAffectedRows",
			setup: func(s *accountStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.APIKey{
				APIKey: &proto.APIKey{Id: id},
			},
			environmentId: environmentId,
			expectedErr:   ErrAPIKeyUnexpectedAffectedRows,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input: &domain.APIKey{
				APIKey: &proto.APIKey{Id: id},
			},
			environmentId: environmentId,
			expectedErr:   errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					gomock.Any(),
					name, int32(role), disabled, maintainer, description, updatedAt, id, environmentId,
				).Return(result, nil)
			},
			input: &domain.APIKey{
				APIKey: &proto.APIKey{Id: id, Name: name, Role: role, Disabled: disabled, Maintainer: maintainer, Description: description, CreatedAt: createdAt, UpdatedAt: updatedAt},
			},
			environmentId: environmentId,
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpdateAPIKey(context.Background(), p.input, p.environmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAPIKey(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc          string
		setup         func(*accountStorage)
		id            string
		environmentId string
		expectedErr   error
	}{
		{
			desc: "ErrAPIKeyNotFound",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:            "id-0",
			environmentId: "ns0",
			expectedErr:   ErrAPIKeyNotFound,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)

			},
			id:            "id-0",
			environmentId: "ns0",
			expectedErr:   errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Any(),
					"id-0", "ns0",
				).Return(row)
			},
			id:            "id-0",
			environmentId: "ns0",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetAPIKey(context.Background(), p.id, p.environmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAPIKeyByAPIKey(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc          string
		setup         func(*accountStorage)
		apiKey        string
		environmentId string
		expectedErr   error
	}{
		{
			desc: "ErrAPIKeyNotFound",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			apiKey:        "id-0",
			environmentId: "ns0",
			expectedErr:   ErrAPIKeyNotFound,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)

			},
			apiKey:        "id-0",
			environmentId: "ns0",
			expectedErr:   errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Any(),
					"id-0", "ns0",
				).Return(row)
			},
			apiKey:        "id-0",
			environmentId: "ns0",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetAPIKeyByAPIKey(context.Background(), p.apiKey, p.environmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAPIKeys(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	getSize := 2
	offset := 5
	limit := 10
	createdAt := 50
	updatedAt := 77

	patterns := []struct {
		desc           string
		setup          func(*accountStorage)
		options        *mysql.ListOptions
		expectedCount  int
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			options:        nil,
			expectedCount:  0,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				var nextCallCount = 0
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().DoAndReturn(func() bool {
					nextCallCount++
					if nextCallCount <= getSize {
						return true
					} else {
						return false
					}
				}).Times(getSize + 1)
				rows.EXPECT().Scan(gomock.Any()).Return(nil).Times(getSize)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(),
					gomock.Any(),
					createdAt, updatedAt,
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Any(),
					createdAt, updatedAt,
				).Return(row)
			},
			options: &mysql.ListOptions{
				Filters: []*mysql.FilterV2{
					{
						Column:   "create_at",
						Operator: mysql.OperatorGreaterThanOrEqual,
						Value:    createdAt,
					},
					{
						Column:   "update_at",
						Operator: mysql.OperatorLessThan,
						Value:    updatedAt,
					},
				},
				Orders: []*mysql.Order{
					{
						Column:    "id",
						Direction: mysql.OrderDirectionAsc,
					},
					{
						Column:    "name",
						Direction: mysql.OrderDirectionDesc,
					},
				},
				Limit:       limit,
				Offset:      offset,
				JSONFilters: nil,
				InFilters:   nil,
				NullFilters: nil,
				SearchQuery: nil,
			},
			expectedCount:  getSize,
			expectedCursor: offset + getSize,
			expectedErr:    nil,
		},
		{
			desc: "Success:No wereParts and no orderParts and no limit and no offset",
			setup: func(s *accountStorage) {
				var nextCallCount = 0
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().DoAndReturn(func() bool {
					nextCallCount++
					if nextCallCount <= getSize {
						return true
					} else {
						return false
					}
				}).Times(getSize + 1)
				rows.EXPECT().Scan(gomock.Any()).Return(nil).Times(getSize)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(),
					gomock.Any(),
					[]interface{}{},
				).Return(rows, nil)

				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Any(),
					[]interface{}{},
				).Return(row)
			},
			options: &mysql.ListOptions{
				Filters:     []*mysql.FilterV2{},
				Orders:      []*mysql.Order{},
				Limit:       0,
				Offset:      0,
				JSONFilters: nil,
				InFilters:   nil,
				NullFilters: nil,
				SearchQuery: nil,
			},
			expectedCount:  getSize,
			expectedCursor: getSize,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			apiKeys, cursor, _, err := storage.ListAPIKeys(context.Background(), p.options)
			assert.Equal(t, p.expectedCount, len(apiKeys))
			if len(apiKeys) > 0 {
				assert.IsType(t, apiKeys, []*proto.APIKey{})
			}
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
