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

	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

func TestNewAutoOpsRuleStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	db := NewAutoOpsRuleStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &autoOpsRuleStorage{}, db)
}

func TestCreateAutoOpsRule(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup                func(*autoOpsRuleStorage)
		input                *domain.AutoOpsRule
		environmentNamespace string
		expectedErr          error
	}{
		{
			setup: func(s *autoOpsRuleStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: &domain.AutoOpsRule{
				AutoOpsRule: &proto.AutoOpsRule{Id: "id-0"},
			},
			environmentNamespace: "ns0",
			expectedErr:          ErrAutoOpsRuleAlreadyExists,
		},
		{
			setup: func(s *autoOpsRuleStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: &domain.AutoOpsRule{
				AutoOpsRule: &proto.AutoOpsRule{Id: "id-1"},
			},
			environmentNamespace: "ns0",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		storage := newAutoOpsRuleStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		err := storage.CreateAutoOpsRule(context.Background(), p.input, p.environmentNamespace)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestUpdateAutoOpsRule(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup                func(*autoOpsRuleStorage)
		input                *domain.AutoOpsRule
		environmentNamespace string
		expectedErr          error
	}{
		{
			setup: func(s *autoOpsRuleStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.AutoOpsRule{
				AutoOpsRule: &proto.AutoOpsRule{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          ErrAutoOpsRuleUnexpectedAffectedRows,
		},
		{
			setup: func(s *autoOpsRuleStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.AutoOpsRule{
				AutoOpsRule: &proto.AutoOpsRule{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		storage := newAutoOpsRuleStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		err := storage.UpdateAutoOpsRule(context.Background(), p.input, p.environmentNamespace)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestGetAutoOpsRule(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup                func(*autoOpsRuleStorage)
		input                string
		environmentNamespace string
		expected             *domain.AutoOpsRule
		expectedErr          error
	}{
		{
			setup: func(s *autoOpsRuleStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input:                "",
			environmentNamespace: "ns0",
			expected:             nil,
			expectedErr:          ErrAutoOpsRuleNotFound,
		},
		{
			setup: func(s *autoOpsRuleStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input:                "id-0",
			environmentNamespace: "ns0",
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
		_, err := storage.GetAutoOpsRule(context.Background(), p.input, p.environmentNamespace)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestListAutoOpsRules(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		setup          func(*autoOpsRuleStorage)
		whereParts     []mysql.WherePart
		orders         []*mysql.Order
		limit          int
		offset         int
		expected       []*proto.AutoOpsRule
		expectedCursor int
		expectedErr    error
	}{
		{
			setup: func(s *autoOpsRuleStorage) {
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
		{
			setup: func(s *autoOpsRuleStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			whereParts: []mysql.WherePart{
				mysql.NewFilter("num", ">=", 5),
			},
			orders: []*mysql.Order{
				mysql.NewOrder("id", mysql.OrderDirectionAsc),
			},
			limit:          10,
			offset:         5,
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
			p.whereParts,
			p.orders,
			p.limit,
			p.offset,
		)
		assert.Equal(t, p.expected, autoOpsRules)
		assert.Equal(t, p.expectedCursor, cursor)
		assert.Equal(t, p.expectedErr, err)
	}
}

func newAutoOpsRuleStorageWithMock(t *testing.T, mockController *gomock.Controller) *autoOpsRuleStorage {
	t.Helper()
	return &autoOpsRuleStorage{mock.NewMockQueryExecer(mockController)}
}
