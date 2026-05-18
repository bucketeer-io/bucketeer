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

func TestNewSegmentStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewSegmentStorage(pgmock.NewMockQueryExecer(mockController))
	assert.IsType(t, &segmentStorage{}, storage)
}

func TestCreateSegment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*segmentStorage)
		segment     *domain.Segment
		expectedErr error
	}{
		{
			desc: "ErrSegmentAlreadyExists",
			setup: func(s *segmentStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, postgres.ErrDuplicateEntry)
			},
			segment:     &domain.Segment{Segment: &proto.Segment{}},
			expectedErr: v2fs.ErrSegmentAlreadyExists,
		},
		{
			desc: "Error",
			setup: func(s *segmentStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			segment:     &domain.Segment{Segment: &proto.Segment{}},
			expectedErr: errInternal,
		},
		{
			desc: "Success",
			setup: func(s *segmentStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			segment:     &domain.Segment{Segment: &proto.Segment{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newSegmentStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateSegment(context.Background(), p.segment, "env")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateSegment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*segmentStorage)
		segment     *domain.Segment
		expectedErr error
	}{
		{
			desc: "ErrSegmentUnexpectedAffectedRows",
			setup: func(s *segmentStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			segment:     &domain.Segment{Segment: &proto.Segment{}},
			expectedErr: v2fs.ErrSegmentUnexpectedAffectedRows,
		},
		{
			desc: "Error",
			setup: func(s *segmentStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			segment:     &domain.Segment{Segment: &proto.Segment{}},
			expectedErr: errInternal,
		},
		{
			desc: "Success",
			setup: func(s *segmentStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			segment:     &domain.Segment{Segment: &proto.Segment{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newSegmentStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpdateSegment(context.Background(), p.segment, "env")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetSegment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*segmentStorage)
		expected    *domain.Segment
		expectedErr error
	}{
		{
			desc: "ErrSegmentNotFound",
			setup: func(s *segmentStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(postgres.ErrNoRows)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expected:    nil,
			expectedErr: v2fs.ErrSegmentNotFound,
		},
		{
			desc: "Error",
			setup: func(s *segmentStorage) {
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
			setup: func(s *segmentStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expected:    &domain.Segment{Segment: &proto.Segment{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newSegmentStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			segment, _, err := storage.GetSegment(context.Background(), "segment-id", "env")
			assert.Equal(t, p.expected, segment)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListSegments(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*segmentStorage)
		params      v2fs.ListSegmentsParams
		expectedErr error
	}{
		{
			desc: "Error",
			setup: func(s *segmentStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			params: v2fs.ListSegmentsParams{
				PageSize:      10,
				EnvironmentID: "env",
			},
			expectedErr: errInternal,
		},
		{
			desc: "Success",
			setup: func(s *segmentStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params: v2fs.ListSegmentsParams{
				PageSize:      10,
				EnvironmentID: "env",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newSegmentStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, _, _, _, err := storage.ListSegments(context.Background(), p.params)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteSegment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*segmentStorage)
		expectedErr error
	}{
		{
			desc: "ErrSegmentUnexpectedAffectedRows",
			setup: func(s *segmentStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: v2fs.ErrSegmentUnexpectedAffectedRows,
		},
		{
			desc: "Error",
			setup: func(s *segmentStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			expectedErr: errInternal,
		},
		{
			desc: "Success",
			setup: func(s *segmentStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newSegmentStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.DeleteSegment(context.Background(), "segment-id")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAllInUseSegments(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*segmentStorage)
		expected    []*v2fs.InUseSegment
		expectedErr error
	}{
		{
			desc: "Error: query fails",
			setup: func(s *segmentStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			expected:    nil,
			expectedErr: errInternal,
		},
		{
			desc: "Success: empty result",
			setup: func(s *segmentStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expected:    []*v2fs.InUseSegment{},
			expectedErr: nil,
		},
		{
			desc: "Error: scan fails",
			setup: func(s *segmentStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(true)
				rows.EXPECT().Scan(gomock.Any()).Return(errInternal)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expected:    nil,
			expectedErr: errInternal,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newSegmentStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			segments, err := storage.ListAllInUseSegments(context.Background())
			assert.Equal(t, p.expected, segments)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListSegmentUsersBySegment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*segmentStorage)
		expected    []*proto.SegmentUser
		expectedErr error
	}{
		{
			desc: "Error: query fails",
			setup: func(s *segmentStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			expected:    nil,
			expectedErr: errInternal,
		},
		{
			desc: "Success: empty result",
			setup: func(s *segmentStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expected:    []*proto.SegmentUser{},
			expectedErr: nil,
		},
		{
			desc: "Error: scan fails",
			setup: func(s *segmentStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(true)
				rows.EXPECT().Scan(gomock.Any()).Return(errInternal)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expected:    nil,
			expectedErr: errInternal,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newSegmentStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			users, err := storage.ListSegmentUsersBySegment(context.Background(), "segment-id", "env")
			assert.Equal(t, p.expected, users)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newSegmentStorageWithMock(t *testing.T, mockController *gomock.Controller) *segmentStorage {
	t.Helper()
	return &segmentStorage{qe: pgmock.NewMockQueryExecer(mockController)}
}
