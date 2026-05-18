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

package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	pgmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestNewSegmentUserStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewSegmentUserStorage(pgmock.NewMockQueryExecer(mockController))
	assert.IsType(t, &segmentUserStorage{}, storage)
}

func TestUpsertSegmentUsers(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*segmentUserStorage)
		users       []*proto.SegmentUser
		expectedErr error
	}{
		{
			desc: "Error",
			setup: func(s *segmentUserStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			users:       []*proto.SegmentUser{{Id: "id-1", SegmentId: "seg-1", UserId: "user-1"}},
			expectedErr: errInternal,
		},
		{
			desc: "Success",
			setup: func(s *segmentUserStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			users:       []*proto.SegmentUser{{Id: "id-1", SegmentId: "seg-1", UserId: "user-1"}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newSegmentUserStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpsertSegmentUsers(context.Background(), p.users, "env")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetSegmentUserPg(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*segmentUserStorage)
		expected    *domain.SegmentUser
		expectedErr error
	}{
		{
			desc: "ErrSegmentUserNotFound",
			setup: func(s *segmentUserStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(postgres.ErrNoRows)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expected:    nil,
			expectedErr: v2fs.ErrSegmentUserNotFound,
		},
		{
			desc: "Error",
			setup: func(s *segmentUserStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errInternal)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expected:    nil,
			expectedErr: errInternal,
		},
		{
			desc: "Success",
			setup: func(s *segmentUserStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expected:    &domain.SegmentUser{SegmentUser: &proto.SegmentUser{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newSegmentUserStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			user, err := storage.GetSegmentUser(context.Background(), "id", "env")
			assert.Equal(t, p.expected, user)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListSegmentUsersPg(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*segmentUserStorage)
		params      v2fs.ListSegmentUsersParams
		expectedErr error
	}{
		{
			desc: "Error",
			setup: func(s *segmentUserStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			params: v2fs.ListSegmentUsersParams{
				SegmentID:     "seg-1",
				EnvironmentID: "env",
				PageSize:      10,
			},
			expectedErr: errInternal,
		},
		{
			desc: "Success",
			setup: func(s *segmentUserStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			params: v2fs.ListSegmentUsersParams{
				SegmentID:     "seg-1",
				EnvironmentID: "env",
				PageSize:      10,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newSegmentUserStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, _, err := storage.ListSegmentUsers(context.Background(), p.params)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newSegmentUserStorageWithMock(t *testing.T, mockController *gomock.Controller) *segmentUserStorage {
	t.Helper()
	return &segmentUserStorage{qe: pgmock.NewMockQueryExecer(mockController)}
}
