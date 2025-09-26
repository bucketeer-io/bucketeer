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

	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

func TestNewProgressiveRolloutStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	db := NewProgressiveRolloutStorage(mock.NewMockClient(mockController))
	assert.IsType(t, &progressiveRolloutStorage{}, db)
}

func TestCreateProgressiveRollout(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc          string
		setup         func(*progressiveRolloutStorage)
		input         *domain.ProgressiveRollout
		environmentId string
		expectedErr   error
	}{
		{
			desc: "error",
			setup: func(s *progressiveRolloutStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: &domain.ProgressiveRollout{
				ProgressiveRollout: &proto.ProgressiveRollout{Id: "id-1"},
			},
			environmentId: "ns0",
			expectedErr:   ErrProgressiveRolloutAlreadyExists,
		},
		{
			desc: "success",
			setup: func(s *progressiveRolloutStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: &domain.ProgressiveRollout{
				ProgressiveRollout: &proto.ProgressiveRollout{Id: "id-1"},
			},
			environmentId: "ns0",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		storage := newProgressiveRolloutStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		err := storage.CreateProgressiveRollout(context.Background(), p.input, p.environmentId)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestListProgressiveRollouts(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup              func(*progressiveRolloutStorage)
		listOpts           *mysql.ListOptions
		expected           []*proto.ProgressiveRollout
		expectedCursor     int
		expectedTotalCount int64
		expectedErr        error
	}{
		{
			setup: func(s *progressiveRolloutStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			listOpts:           nil,
			expected:           nil,
			expectedCursor:     0,
			expectedTotalCount: 0,
			expectedErr:        errors.New("error"),
		},
		{
			setup: func(s *progressiveRolloutStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
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
			expected:           []*proto.ProgressiveRollout{},
			expectedCursor:     5,
			expectedTotalCount: 0,
			expectedErr:        nil,
		},
	}
	for _, p := range patterns {
		storage := newProgressiveRolloutStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		pr, totalCount, cursor, err := storage.ListProgressiveRollouts(context.Background(), p.listOpts)
		assert.Equal(t, p.expected, pr)
		assert.Equal(t, p.expectedCursor, cursor)
		assert.Equal(t, p.expectedTotalCount, totalCount)
		assert.Equal(t, p.expectedErr, err)
	}
}

func newProgressiveRolloutStorageWithMock(t *testing.T, mockController *gomock.Controller) *progressiveRolloutStorage {
	t.Helper()
	return &progressiveRolloutStorage{mock.NewMockQueryExecer(mockController)}
}
