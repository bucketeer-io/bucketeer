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

	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestNewSegmentUserStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewSegmentUserStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &segmentUserStorage{}, storage)
}

func TestUpsertSegmentUsers(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	errInternal := errors.New("test error")
	tests := []struct {
		desc        string
		setup       func(*mock.MockQueryExecer)
		users       []*proto.SegmentUser
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(m *mock.MockQueryExecer) {
				m.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errInternal)
			},
			users: []*proto.SegmentUser{
				{Id: "id-1", SegmentId: "seg-1", UserId: "user-1"},
			},
			expectedErr: errInternal,
		},
		{
			desc: "success",
			setup: func(m *mock.MockQueryExecer) {
				m.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			users: []*proto.SegmentUser{
				{Id: "id-1", SegmentId: "seg-1", UserId: "user-1"},
			},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			m := mock.NewMockQueryExecer(mockController)
			tt.setup(m)
			storage := NewSegmentUserStorage(m)
			err := storage.UpsertSegmentUsers(context.Background(), tt.users, "env")
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestGetSegmentUser(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	tests := []struct {
		desc        string
		setup       func(*mock.MockQueryExecer)
		expected    *domain.SegmentUser
		expectedErr error
	}{
		{
			desc: "err: not found",
			setup: func(m *mock.MockQueryExecer) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				m.EXPECT().QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
			},
			expected:    nil,
			expectedErr: v2fs.ErrSegmentUserNotFound,
		},
		{
			desc: "success",
			setup: func(m *mock.MockQueryExecer) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				m.EXPECT().QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
			},
			expected:    &domain.SegmentUser{SegmentUser: &proto.SegmentUser{}},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			m := mock.NewMockQueryExecer(mockController)
			tt.setup(m)
			storage := NewSegmentUserStorage(m)
			user, err := storage.GetSegmentUser(context.Background(), "id", "env")
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, user)
		})
	}
}

func TestListSegmentUsersMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	errInternal := errors.New("test error")
	tests := []struct {
		desc        string
		setup       func(*mock.MockQueryExecer)
		params      v2fs.ListSegmentUsersParams
		expectedErr error
	}{
		{
			desc: "err: query error",
			setup: func(m *mock.MockQueryExecer) {
				m.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errInternal)
			},
			params: v2fs.ListSegmentUsersParams{
				SegmentID:     "seg-1",
				EnvironmentID: "env",
				PageSize:      10,
			},
			expectedErr: errInternal,
		},
		{
			desc: "success",
			setup: func(m *mock.MockQueryExecer) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				m.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(rows, nil)
			},
			params: v2fs.ListSegmentUsersParams{
				SegmentID:     "seg-1",
				EnvironmentID: "env",
				PageSize:      10,
			},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			m := mock.NewMockQueryExecer(mockController)
			tt.setup(m)
			storage := NewSegmentUserStorage(m)
			_, _, err := storage.ListSegmentUsers(context.Background(), tt.params)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
