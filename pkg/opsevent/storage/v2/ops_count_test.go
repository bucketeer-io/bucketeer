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

	"github.com/bucketeer-io/bucketeer/pkg/opsevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

func TestNewOpsEventStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	db := NewOpsCountStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &opsCountStorage{}, db)
}

func TestUpsertOpsCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup                func(*opsCountStorage)
		input                *domain.OpsCount
		environmentNamespace string
		expectedErr          error
	}{
		{
			setup: func(s *opsCountStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input:                &domain.OpsCount{OpsCount: &proto.OpsCount{}},
			environmentNamespace: "ns",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		storage := newOpsCountStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		err := storage.UpsertOpsCount(context.Background(), p.environmentNamespace, p.input)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestListOpsCounts(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		setup          func(*opsCountStorage)
		whereParts     []mysql.WherePart
		orders         []*mysql.Order
		limit          int
		offset         int
		expected       []*proto.OpsCount
		expectedCursor int
		expectedErr    error
	}{
		{
			setup: func(s *opsCountStorage) {
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
			setup: func(s *opsCountStorage) {
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
			expected:       []*proto.OpsCount{},
			expectedCursor: 5,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		storage := newOpsCountStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		opsCounts, cursor, err := storage.ListOpsCounts(
			context.Background(),
			p.whereParts,
			p.orders,
			p.limit,
			p.offset,
		)
		assert.Equal(t, p.expected, opsCounts)
		assert.Equal(t, p.expectedCursor, cursor)
		assert.Equal(t, p.expectedErr, err)
	}
}

func newOpsCountStorageWithMock(t *testing.T, mockController *gomock.Controller) *opsCountStorage {
	t.Helper()
	return &opsCountStorage{mock.NewMockQueryExecer(mockController)}
}
