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

	"github.com/bucketeer-io/bucketeer/pkg/push/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/push"
)

func TestNewPushStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewPushStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &pushStorage{}, storage)
}

func TestCreatePush(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc                 string
		setup                func(*pushStorage)
		input                *domain.Push
		environmentNamespace string
		expectedErr          error
	}{
		{
			desc: "ErrPushAlreadyExists",
			setup: func(s *pushStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: &domain.Push{
				Push: &proto.Push{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          ErrPushAlreadyExists,
		},
		{
			desc: "Error",
			setup: func(s *pushStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))

			},
			input: &domain.Push{
				Push: &proto.Push{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *pushStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: &domain.Push{
				Push: &proto.Push{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newpushStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreatePush(context.Background(), p.input, p.environmentNamespace)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdatePush(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc                 string
		setup                func(*pushStorage)
		input                *domain.Push
		environmentNamespace string
		expectedErr          error
	}{
		{
			desc: "ErrPushUnexpectedAffectedRows",
			setup: func(s *pushStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.Push{
				Push: &proto.Push{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          ErrPushUnexpectedAffectedRows,
		},
		{
			desc: "Error",
			setup: func(s *pushStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))

			},
			input: &domain.Push{
				Push: &proto.Push{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *pushStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.Push{
				Push: &proto.Push{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newpushStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpdatePush(context.Background(), p.input, p.environmentNamespace)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetPush(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc                 string
		setup                func(*pushStorage)
		id                   string
		environmentNamespace string
		expectedErr          error
	}{
		{
			desc: "ErrPushNotFound",
			setup: func(s *pushStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:                   "id-0",
			environmentNamespace: "ns",
			expectedErr:          ErrPushNotFound,
		},
		{
			desc: "Error",
			setup: func(s *pushStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)

			},
			id:                   "id-0",
			environmentNamespace: "ns",
			expectedErr:          errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *pushStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:                   "id-0",
			environmentNamespace: "ns",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newpushStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetPush(context.Background(), p.id, p.environmentNamespace)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListPushs(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc           string
		setup          func(*pushStorage)
		whereParts     []mysql.WherePart
		orders         []*mysql.Order
		limit          int
		offset         int
		expected       []*proto.Push
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "Error",
			setup: func(s *pushStorage) {
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
			desc: "Success",
			setup: func(s *pushStorage) {
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
			expected:       []*proto.Push{},
			expectedCursor: 5,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newpushStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			pushs, cursor, _, err := storage.ListPushes(
				context.Background(),
				p.whereParts,
				p.orders,
				p.limit,
				p.offset,
			)
			assert.Equal(t, p.expected, pushs)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newpushStorageWithMock(t *testing.T, mockController *gomock.Controller) *pushStorage {
	t.Helper()
	return &pushStorage{mock.NewMockQueryExecer(mockController)}
}
