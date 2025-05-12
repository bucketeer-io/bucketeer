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

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestNewSegmentStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewSegmentStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &segmentStorage{}, storage)
}

func TestCreateSegment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	tests := []struct {
		desc        string
		setup       func(*mock.MockQueryExecer)
		segment     *domain.Segment
		expectedErr error
	}{
		{
			desc: "err: segment already exists",
			setup: func(mockQueryExecer *mock.MockQueryExecer) {
				mockQueryExecer.EXPECT().ExecContext(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			segment: &domain.Segment{
				Segment: &proto.Segment{},
			},
			expectedErr: ErrSegmentAlreadyExists,
		},
		{
			desc: "success",
			setup: func(mockQueryExecer *mock.MockQueryExecer) {
				mockQueryExecer.EXPECT().ExecContext(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, nil)
			},
			segment: &domain.Segment{
				Segment: &proto.Segment{},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockQueryExecer := mock.NewMockQueryExecer(mockController)
			tt.setup(mockQueryExecer)
			storage := NewSegmentStorage(mockQueryExecer)
			err := storage.CreateSegment(context.Background(), tt.segment, "test-environment-id")
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestUpdateSegment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	tests := []struct {
		desc        string
		setup       func(*mock.MockQueryExecer)
		segment     *domain.Segment
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(mockQueryExecer *mock.MockQueryExecer) {
				mockQueryExecer.EXPECT().ExecContext(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, errors.New("test error"))
			},
			segment: &domain.Segment{
				Segment: &proto.Segment{},
			},
			expectedErr: errors.New("test error"),
		},
		{
			desc: "success",
			setup: func(mockQueryExecer *mock.MockQueryExecer) {
				result := mock.NewMockResult(mockController)
				mockQueryExecer.EXPECT().ExecContext(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(result, nil)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
			},
			segment: &domain.Segment{
				Segment: &proto.Segment{},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockQueryExecer := mock.NewMockQueryExecer(mockController)
			tt.setup(mockQueryExecer)
			storage := NewSegmentStorage(mockQueryExecer)
			err := storage.UpdateSegment(context.Background(), tt.segment, "test-environment-id")
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestGetSegment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	tests := []struct {
		desc        string
		setup       func(*mock.MockQueryExecer)
		expected    *domain.Segment
		expectedErr error
	}{
		{
			desc: "err: segment not found",
			setup: func(mockQueryExecer *mock.MockQueryExecer) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				mockQueryExecer.EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(row)
			},
			expected:    nil,
			expectedErr: ErrSegmentNotFound,
		},
		{
			desc: "success",
			setup: func(mockQueryExecer *mock.MockQueryExecer) {
				row := mock.NewMockRow(mockController)
				mockQueryExecer.EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(row)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
			},
			expected: &domain.Segment{
				Segment: &proto.Segment{},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockQueryExecer := mock.NewMockQueryExecer(mockController)
			tt.setup(mockQueryExecer)
			storage := NewSegmentStorage(mockQueryExecer)
			segment, _, err := storage.GetSegment(context.Background(), "test-segment-id", "test-environment-id")
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, segment)
		})
	}
}

func TestListSegments(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	tests := []struct {
		desc        string
		setup       func(*mock.MockQueryExecer)
		options     *mysql.ListOptions
		expectedErr error
	}{
		{
			desc: "err: query error",
			setup: func(mockQueryExecer *mock.MockQueryExecer) {
				mockQueryExecer.EXPECT().QueryContext(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, errors.New("test error"))
			},
			options: &mysql.ListOptions{
				Limit:  10,
				Offset: 0,
			},
			expectedErr: errors.New("test error"),
		},
		{
			desc: "success",
			setup: func(mockQueryExecer *mock.MockQueryExecer) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				mockQueryExecer.EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				mockQueryExecer.EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			options: &mysql.ListOptions{
				Limit:  10,
				Offset: 0,
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockQueryExecer := mock.NewMockQueryExecer(mockController)
			tt.setup(mockQueryExecer)
			storage := NewSegmentStorage(mockQueryExecer)
			isInUseStatus := false
			_, _, _, _, err := storage.ListSegments(context.Background(), tt.options, &isInUseStatus)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestDeleteSegment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	tests := []struct {
		desc        string
		setup       func(*mock.MockQueryExecer)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(mockQueryExecer *mock.MockQueryExecer) {
				mockQueryExecer.EXPECT().ExecContext(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, errors.New("test error"))
			},
			expectedErr: errors.New("test error"),
		},
		{
			desc: "success",
			setup: func(mockQueryExecer *mock.MockQueryExecer) {
				result := mock.NewMockResult(mockController)
				mockQueryExecer.EXPECT().ExecContext(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(result, nil)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockQueryExecer := mock.NewMockQueryExecer(mockController)
			tt.setup(mockQueryExecer)
			storage := NewSegmentStorage(mockQueryExecer)
			err := storage.DeleteSegment(context.Background(), "test-segment-id")
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
