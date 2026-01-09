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

package v2

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/push/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

func TestNewPushStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewPushStorage(mock.NewMockClient(mockController))
	assert.IsType(t, &pushStorage{}, storage)
}

func TestCreatePush(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc          string
		setup         func(*pushStorage)
		input         *domain.Push
		environmentId string
		expectedErr   error
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
			environmentId: "ns",
			expectedErr:   ErrPushAlreadyExists,
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
			environmentId: "ns",
			expectedErr:   errors.New("error"),
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
			environmentId: "ns",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newPushStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreatePush(context.Background(), p.input, p.environmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdatePush(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc          string
		setup         func(*pushStorage)
		input         *domain.Push
		environmentId string
		expectedErr   error
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
			environmentId: "ns",
			expectedErr:   ErrPushUnexpectedAffectedRows,
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
			environmentId: "ns",
			expectedErr:   errors.New("error"),
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
			environmentId: "ns",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newPushStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpdatePush(context.Background(), p.input, p.environmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetPush(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc          string
		setup         func(*pushStorage)
		id            string
		environmentId string
		expectedErr   error
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
			id:            "id-0",
			environmentId: "ns",
			expectedErr:   ErrPushNotFound,
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
			id:            "id-0",
			environmentId: "ns",
			expectedErr:   errors.New("error"),
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
			id:            "id-0",
			environmentId: "ns",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newPushStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetPush(context.Background(), p.id, p.environmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListPushes(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc           string
		setup          func(*pushStorage)
		options        *mysql.ListOptions
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
			options:        nil,
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
			options: &mysql.ListOptions{
				Limit:  10,
				Offset: 5,
				Filters: []*mysql.FilterV2{
					{
						Column:   "num",
						Operator: mysql.OperatorGreaterThanOrEqual,
						Value:    5,
					},
				},
				InFilters:   nil,
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
			expected:       []*proto.Push{},
			expectedCursor: 5,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newPushStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			pushs, cursor, _, err := storage.ListPushes(
				context.Background(),
				p.options,
			)
			assert.Equal(t, p.expected, pushs)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeletePush(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc          string
		setup         func(*pushStorage)
		id            string
		environmentId string
		expectedErr   error
	}{
		{
			desc: "Err push unexpected affected rows",
			setup: func(s *pushStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			id:            "id-0",
			environmentId: "ns",
			expectedErr:   ErrPushUnexpectedAffectedRows,
		},
		{
			desc: "Error",
			setup: func(s *pushStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, errors.New("error"))
			},
			id:            "id-0",
			environmentId: "ns",
			expectedErr:   errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *pushStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			id:            "id-0",
			environmentId: "ns",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newPushStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.DeletePush(context.Background(), p.id, p.environmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newPushStorageWithMock(t *testing.T, mockController *gomock.Controller) *pushStorage {
	t.Helper()
	return &pushStorage{mock.NewMockQueryExecer(mockController)}
}
