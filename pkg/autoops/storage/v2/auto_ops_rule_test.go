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

	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

func TestNewAutoOpsRuleStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	client := mock.NewMockClient(mockController)
	db := NewAutoOpsRuleStorage(client)
	assert.IsType(t, &autoOpsRuleStorage{}, db)
}

func TestCreateAutoOpsRule(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup         func(*autoOpsRuleStorage)
		input         *domain.AutoOpsRule
		environmentId string
		expectedErr   error
	}{
		{
			setup: func(s *autoOpsRuleStorage) {
				qe := mock.NewMockQueryExecer(mockController)
				s.client.(*mock.MockClient).EXPECT().Qe(
					gomock.Any(),
				).Return(qe)
				qe.EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: &domain.AutoOpsRule{
				AutoOpsRule: &proto.AutoOpsRule{Id: "id-0"},
			},
			environmentId: "ns0",
			expectedErr:   ErrAutoOpsRuleAlreadyExists,
		},
		{
			setup: func(s *autoOpsRuleStorage) {
				qe := mock.NewMockQueryExecer(mockController)
				s.client.(*mock.MockClient).EXPECT().Qe(
					gomock.Any(),
				).Return(qe)
				qe.EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: &domain.AutoOpsRule{
				AutoOpsRule: &proto.AutoOpsRule{Id: "id-1"},
			},
			environmentId: "ns0",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		storage := newAutoOpsRuleStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		err := storage.CreateAutoOpsRule(context.Background(), p.input, p.environmentId)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestUpdateAutoOpsRule(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup         func(*autoOpsRuleStorage)
		input         *domain.AutoOpsRule
		environmentId string
		expectedErr   error
	}{
		{
			setup: func(s *autoOpsRuleStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				qe := mock.NewMockQueryExecer(mockController)
				s.client.(*mock.MockClient).EXPECT().Qe(
					gomock.Any(),
				).Return(qe)
				qe.EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.AutoOpsRule{
				AutoOpsRule: &proto.AutoOpsRule{Id: "id-0"},
			},
			environmentId: "ns",
			expectedErr:   ErrAutoOpsRuleUnexpectedAffectedRows,
		},
		{
			setup: func(s *autoOpsRuleStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				qe := mock.NewMockQueryExecer(mockController)
				s.client.(*mock.MockClient).EXPECT().Qe(
					gomock.Any(),
				).Return(qe)
				qe.EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.AutoOpsRule{
				AutoOpsRule: &proto.AutoOpsRule{Id: "id-0"},
			},
			environmentId: "ns",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		storage := newAutoOpsRuleStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		err := storage.UpdateAutoOpsRule(context.Background(), p.input, p.environmentId)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestGetAutoOpsRule(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup         func(*autoOpsRuleStorage)
		input         string
		environmentId string
		expected      *domain.AutoOpsRule
		expectedErr   error
	}{
		{
			setup: func(s *autoOpsRuleStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				qe := mock.NewMockQueryExecer(mockController)
				s.client.(*mock.MockClient).EXPECT().Qe(
					gomock.Any(),
				).Return(qe)
				qe.EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input:         "",
			environmentId: "ns0",
			expected:      nil,
			expectedErr:   ErrAutoOpsRuleNotFound,
		},
		{
			setup: func(s *autoOpsRuleStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				qe := mock.NewMockQueryExecer(mockController)
				s.client.(*mock.MockClient).EXPECT().Qe(
					gomock.Any(),
				).Return(qe)
				qe.EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input:         "id-0",
			environmentId: "ns0",
			expected: &domain.AutoOpsRule{
				AutoOpsRule: &proto.AutoOpsRule{Id: "id-0"},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		storage := newAutoOpsRuleStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		_, err := storage.GetAutoOpsRule(context.Background(), p.input, p.environmentId)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestListAutoOpsRules(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		setup          func(*autoOpsRuleStorage)
		listOpts       *mysql.ListOptions
		expected       []*proto.AutoOpsRule
		expectedCursor int
		expectedErr    error
	}{
		{
			setup: func(s *autoOpsRuleStorage) {
				qe := mock.NewMockQueryExecer(mockController)
				s.client.(*mock.MockClient).EXPECT().Qe(
					gomock.Any(),
				).Return(qe)
				qe.EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			listOpts:       nil,
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			setup: func(s *autoOpsRuleStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				qe := mock.NewMockQueryExecer(mockController)
				s.client.(*mock.MockClient).EXPECT().Qe(
					gomock.Any(),
				).Return(qe)
				qe.EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			listOpts: &mysql.ListOptions{
				Limit:  10,
				Offset: 5,
				Filters: []*mysql.FilterV2{
					{
						Column:   "num",
						Operator: mysql.OperatorGreaterThanOrEqual,
						Value:    5,
					},
				},
				InFilter:    nil,
				NullFilters: nil,
				JSONFilters: nil,
				SearchQuery: nil,
				Orders: []*mysql.Order{
					{
						Column:    "id",
						Direction: mysql.OrderDirectionAsc,
					},
				},
			},
			expected:       []*proto.AutoOpsRule{},
			expectedCursor: 5,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		storage := newAutoOpsRuleStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		autoOpsRules, cursor, err := storage.ListAutoOpsRules(
			context.Background(),
			p.listOpts,
		)
		assert.Equal(t, p.expected, autoOpsRules)
		assert.Equal(t, p.expectedCursor, cursor)
		assert.Equal(t, p.expectedErr, err)
	}
}

func newAutoOpsRuleStorageWithMock(t *testing.T, mockController *gomock.Controller) *autoOpsRuleStorage {
	t.Helper()
	return &autoOpsRuleStorage{mock.NewMockClient(mockController)}
}
