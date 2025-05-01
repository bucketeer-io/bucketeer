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

package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/tag/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/tag"
)

func TestNewTagStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewTagStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &tagStorage{}, storage)
}

func TestUpsertTag(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*tagStorage)
		input       *domain.Tag
		expectedErr error
	}{
		{
			desc: "ErrTagAlreadyExists",
			setup: func(s *tagStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: &domain.Tag{
				Tag: &proto.Tag{Id: "tag-id-0"},
			},
			expectedErr: mysql.ErrDuplicateEntry,
		},
		{
			desc: "Error",
			setup: func(s *tagStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input: &domain.Tag{
				Tag: &proto.Tag{Id: "tag-id-0"},
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *tagStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					insertTagSQL,
					"tag-id-0",
					gomock.Any(),
					int64(1),
					int64(2),
					int32(proto.Tag_FEATURE_FLAG),
					"env-0",
				).Return(nil, nil)
			},
			input: &domain.Tag{
				Tag: &proto.Tag{
					Id:            "tag-id-0",
					Name:          "test-tag",
					CreatedAt:     1,
					UpdatedAt:     2,
					EntityType:    proto.Tag_FEATURE_FLAG,
					EnvironmentId: "env-0",
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newTagStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpsertTag(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetTag(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc          string
		setup         func(*tagStorage)
		id            string
		environmentId string
		expectedTag   *domain.Tag
		expectedErr   error
	}{
		{
			desc: "ErrTagNotFound",
			setup: func(s *tagStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:            "tag-id-0",
			environmentId: "env-0",
			expectedTag:   nil,
			expectedErr:   ErrTagNotFound,
		},
		{
			desc: "Error",
			setup: func(s *tagStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:            "tag-id-0",
			environmentId: "env-0",
			expectedTag:   nil,
			expectedErr:   errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *tagStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(
					gomock.Any(), // id
					gomock.Any(), // name
					gomock.Any(), // created_at
					gomock.Any(), // updated_at
					gomock.Any(), // entity_type
					gomock.Any(), // environment_id
					gomock.Any(), // environment_name
				).Do(func(args ...interface{}) {
					*args[0].(*string) = "tag-id-0"
					*args[1].(*string) = "test-tag"
					*args[2].(*int64) = int64(1)
					*args[3].(*int64) = int64(2)
					*args[4].(*int32) = int32(proto.Tag_FEATURE_FLAG)
					*args[5].(*string) = "env-0"
					*args[6].(*string) = "test-env"
				}).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					selectTagSQL,
					"tag-id-0",
					"env-0",
				).Return(row)
			},
			id:            "tag-id-0",
			environmentId: "env-0",
			expectedTag: &domain.Tag{
				Tag: &proto.Tag{
					Id:              "tag-id-0",
					Name:            "test-tag",
					CreatedAt:       1,
					UpdatedAt:       2,
					EntityType:      proto.Tag_FEATURE_FLAG,
					EnvironmentId:   "env-0",
					EnvironmentName: "test-env",
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newTagStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			tag, err := storage.GetTag(context.Background(), p.id, p.environmentId)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.Equal(t, p.expectedTag, tag)
			}
		})
	}
}

func TestListTags(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(*tagStorage)
		whereParts     []mysql.WherePart
		orders         []*mysql.Order
		limit          int
		offset         int
		expectedCount  int
		expectedCursor int
		expectedErr    error
		expectedTags   []*proto.Tag
	}{
		{
			desc: "Error",
			setup: func(s *tagStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			whereParts:     nil,
			orders:         nil,
			limit:          10,
			offset:         0,
			expectedCount:  0,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
			expectedTags:   nil,
		},
		{
			desc: "Success",
			setup: func(s *tagStorage) {
				var nextCallCount = 0
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().DoAndReturn(func() bool {
					nextCallCount++
					return nextCallCount <= 1
				}).Times(2)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(rows, nil)
				rows.EXPECT().Scan(
					gomock.Any(), // id
					gomock.Any(), // name
					gomock.Any(), // created_at
					gomock.Any(), // updated_at
					gomock.Any(), // entity_type
					gomock.Any(), // environment_id
					gomock.Any(), // environment_name
				).Do(func(args ...interface{}) {
					*args[0].(*string) = "tag-id-0"
					*args[1].(*string) = "test-tag"
					*args[2].(*int64) = int64(1)
					*args[3].(*int64) = int64(2)
					*args[4].(*int32) = int32(proto.Tag_FEATURE_FLAG)
					*args[5].(*string) = "env-0"
					*args[6].(*string) = "test-env"
				}).Return(nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(row)
			},
			whereParts: []mysql.WherePart{
				mysql.NewFilter("tag.environment_id", "=", "env-0"),
			},
			orders: []*mysql.Order{
				mysql.NewOrder("tag.name", mysql.OrderDirectionAsc),
			},
			limit:          10,
			offset:         0,
			expectedCount:  1,
			expectedCursor: 1,
			expectedErr:    nil,
			expectedTags: []*proto.Tag{
				{
					Id:              "tag-id-0",
					Name:            "test-tag",
					CreatedAt:       1,
					UpdatedAt:       2,
					EntityType:      proto.Tag_FEATURE_FLAG,
					EnvironmentId:   "env-0",
					EnvironmentName: "test-env",
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newTagStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			tags, cursor, _, err := storage.ListTags(
				context.Background(),
				p.whereParts,
				p.orders,
				p.limit,
				p.offset,
			)
			assert.Equal(t, p.expectedCount, len(tags))
			if tags != nil {
				assert.IsType(t, []*proto.Tag{}, tags)
				assert.Equal(t, p.expectedTags, tags)
			}
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAllEnvironmentTags(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc            string
		setup           func(*tagStorage)
		expectedErr     error
		expectedEnvTags []*proto.EnvironmentTag
	}{
		{
			desc: "Error",
			setup: func(s *tagStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr:     errors.New("error"),
			expectedEnvTags: nil,
		},
		{
			desc: "Success:Empty",
			setup: func(s *tagStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(),
					selectAllEnvironmentTagsSQL,
				).Return(rows, nil)
			},
			expectedErr:     nil,
			expectedEnvTags: []*proto.EnvironmentTag{},
		},
		{
			desc: "Success:WithData",
			setup: func(s *tagStorage) {
				var nextCallCount = 0
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().DoAndReturn(func() bool {
					nextCallCount++
					return nextCallCount <= 1
				}).Times(2)
				rows.EXPECT().Scan(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Do(func(args ...interface{}) {
					*args[0].(*string) = "env-0"
					*args[1].(*string) = "tag-id-0"
					*args[2].(*string) = "tag-name"
					*args[3].(*int64) = int64(1)
					*args[4].(*int64) = int64(2)
					*args[5].(*int32) = int32(proto.Tag_FEATURE_FLAG)
					*args[6].(*string) = "tag-env-id-0"
					*args[7].(*string) = "test-env-name"
				})
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(),
					selectAllEnvironmentTagsSQL,
				).Return(rows, nil)
			},
			expectedErr: nil,
			expectedEnvTags: []*proto.EnvironmentTag{
				{
					EnvironmentId: "env-0",
					Tags: []*proto.Tag{
						{
							Id:              "tag-id-0",
							Name:            "tag-name",
							CreatedAt:       1,
							UpdatedAt:       2,
							EntityType:      proto.Tag_FEATURE_FLAG,
							EnvironmentId:   "tag-env-id-0",
							EnvironmentName: "test-env-name",
						},
					},
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newTagStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			envTags, err := storage.ListAllEnvironmentTags(context.Background())
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, envTags)
				assert.IsType(t, []*proto.EnvironmentTag{}, envTags)
				assert.Equal(t, p.expectedEnvTags, envTags)
			}
		})
	}
}

func TestDeleteTag(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*tagStorage)
		id          string
		expectedErr error
	}{
		{
			desc: "ErrTagUnexpectedAffectedRows",
			setup: func(s *tagStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			id:          "tag-id-0",
			expectedErr: ErrTagUnexpectedAffectedRows,
		},
		{
			desc: "Error",
			setup: func(s *tagStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			id:          "tag-id-0",
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *tagStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					deleteTagSQL,
					"tag-id-0",
				).Return(result, nil)
			},
			id:          "tag-id-0",
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newTagStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.DeleteTag(context.Background(), p.id)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newTagStorageWithMock(t *testing.T, mockController *gomock.Controller) *tagStorage {
	t.Helper()
	return &tagStorage{mock.NewMockQueryExecer(mockController)}
}
