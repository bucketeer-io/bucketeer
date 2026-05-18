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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/tag/domain"
	tagstorage "github.com/bucketeer-io/bucketeer/v2/pkg/tag/storage"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/tag"
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
			desc: "ErrDuplicateEntry",
			setup: func(s *tagStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, pgstorage.ErrDuplicateEntry)
			},
			input: &domain.Tag{
				Tag: &proto.Tag{Id: "tag-id-0"},
			},
			expectedErr: pgstorage.ErrDuplicateEntry,
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
				row.EXPECT().Scan(gomock.Any()).Return(pgstorage.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:            "tag-id-0",
			environmentId: "env-0",
			expectedTag:   nil,
			expectedErr:   tagstorage.ErrTagNotFound,
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

func TestGetTagByName(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc          string
		setup         func(*tagStorage)
		name          string
		environmentId string
		entityType    proto.Tag_EntityType
		expectedTag   *domain.Tag
		expectedErr   error
	}{
		{
			desc: "ErrTagNotFound",
			setup: func(s *tagStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(pgstorage.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			name:          "test-tag",
			environmentId: "env-0",
			entityType:    proto.Tag_FEATURE_FLAG,
			expectedTag:   nil,
			expectedErr:   tagstorage.ErrTagNotFound,
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
			name:          "test-tag",
			environmentId: "env-0",
			entityType:    proto.Tag_FEATURE_FLAG,
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
					selectTagByNameSQL,
					"test-tag",
					"env-0",
					int32(proto.Tag_FEATURE_FLAG),
				).Return(row)
			},
			name:          "test-tag",
			environmentId: "env-0",
			entityType:    proto.Tag_FEATURE_FLAG,
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
			tag, err := storage.GetTagByName(context.Background(), p.name, p.environmentId, p.entityType)
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
		params         tagstorage.ListTagsParams
		expectedCount  int
		expectedCursor int
		expectedErr    error
		expectedTags   []*proto.Tag
	}{
		{
			desc: "ErrInvalidCursor",
			params: tagstorage.ListTagsParams{
				EnvironmentID: "env-0",
				PageSize:      10,
				Cursor:        "invalid",
			},
			expectedCount:  0,
			expectedCursor: 0,
			expectedErr:    tagstorage.ErrInvalidListTagsCursor,
			expectedTags:   nil,
		},
		{
			desc: "ErrInvalidOrderBy",
			params: tagstorage.ListTagsParams{
				EnvironmentID: "env-0",
				PageSize:      10,
				Cursor:        "0",
				OrderBy:       proto.ListTagsRequest_OrderBy(99),
			},
			expectedCount:  0,
			expectedCursor: 0,
			expectedErr:    tagstorage.ErrInvalidListTagsOrderBy,
			expectedTags:   nil,
		},
		{
			desc: "Error",
			setup: func(s *tagStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			params: tagstorage.ListTagsParams{
				EnvironmentID: "env-0",
				PageSize:      10,
				Cursor:        "0",
			},
			expectedCount:  0,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
			expectedTags:   nil,
		},
		{
			desc: "Success:Empty",
			setup: func(s *tagStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Do(func(args ...interface{}) {
					*args[0].(*int64) = int64(0)
				}).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params: tagstorage.ListTagsParams{
				EnvironmentID: "env-0",
				PageSize:      10,
				Cursor:        "0",
			},
			expectedCount:  0,
			expectedCursor: 0,
			expectedErr:    nil,
			expectedTags:   []*proto.Tag{},
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
				row.EXPECT().Scan(gomock.Any()).Do(func(args ...interface{}) {
					*args[0].(*int64) = int64(1)
				}).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(row)
			},
			params: tagstorage.ListTagsParams{
				EnvironmentID: "env-0",
				PageSize:      10,
				Cursor:        "0",
			},
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
		{
			desc: "Success:WithOffset",
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
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				rows.EXPECT().Scan(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Do(func(args ...interface{}) {
					*args[0].(*string) = "tag-id-5"
					*args[1].(*string) = "tag-five"
					*args[2].(*int64) = int64(5)
					*args[3].(*int64) = int64(6)
					*args[4].(*int32) = int32(proto.Tag_FEATURE_FLAG)
					*args[5].(*string) = "env-0"
					*args[6].(*string) = "test-env"
				}).Return(nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Do(func(args ...interface{}) {
					*args[0].(*int64) = int64(10)
				}).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params: tagstorage.ListTagsParams{
				EnvironmentID: "env-0",
				PageSize:      5,
				Cursor:        "5",
			},
			expectedCount:  1,
			expectedCursor: 6,
			expectedErr:    nil,
			expectedTags: []*proto.Tag{
				{
					Id:              "tag-id-5",
					Name:            "tag-five",
					CreatedAt:       5,
					UpdatedAt:       6,
					EntityType:      proto.Tag_FEATURE_FLAG,
					EnvironmentId:   "env-0",
					EnvironmentName: "test-env",
				},
			},
		},
		{
			desc: "Success:WithFilters",
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
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				rows.EXPECT().Scan(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
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
				row.EXPECT().Scan(gomock.Any()).Do(func(args ...interface{}) {
					*args[0].(*int64) = int64(1)
				}).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params: tagstorage.ListTagsParams{
				EnvironmentID:  "env-0",
				OrganizationID: "org-0",
				EntityType:     proto.Tag_FEATURE_FLAG,
				SearchKeyword:  "test",
				OrderBy:        proto.ListTagsRequest_NAME,
				OrderDirection: proto.ListTagsRequest_DESC,
				PageSize:       10,
				Cursor:         "0",
			},
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
			tags, cursor, _, err := storage.ListTags(context.Background(), p.params)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.Equal(t, p.expectedCount, len(tags))
				assert.IsType(t, []*proto.Tag{}, tags)
				assert.Equal(t, p.expectedTags, tags)
			}
			assert.Equal(t, p.expectedCursor, cursor)
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
		{
			desc: "Success:MultipleEnvironments",
			setup: func(s *tagStorage) {
				var nextCallCount = 0
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().DoAndReturn(func() bool {
					nextCallCount++
					return nextCallCount <= 2
				}).Times(3)
				gomock.InOrder(
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
						*args[2].(*string) = "tag-name-0"
						*args[3].(*int64) = int64(1)
						*args[4].(*int64) = int64(2)
						*args[5].(*int32) = int32(proto.Tag_FEATURE_FLAG)
						*args[6].(*string) = "env-0"
						*args[7].(*string) = "env-name-0"
					}).Return(nil),
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
						*args[0].(*string) = "env-1"
						*args[1].(*string) = "tag-id-1"
						*args[2].(*string) = "tag-name-1"
						*args[3].(*int64) = int64(3)
						*args[4].(*int64) = int64(4)
						*args[5].(*int32) = int32(proto.Tag_FEATURE_FLAG)
						*args[6].(*string) = "env-1"
						*args[7].(*string) = "env-name-1"
					}).Return(nil),
				)
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
							Name:            "tag-name-0",
							CreatedAt:       1,
							UpdatedAt:       2,
							EntityType:      proto.Tag_FEATURE_FLAG,
							EnvironmentId:   "env-0",
							EnvironmentName: "env-name-0",
						},
					},
				},
				{
					EnvironmentId: "env-1",
					Tags: []*proto.Tag{
						{
							Id:              "tag-id-1",
							Name:            "tag-name-1",
							CreatedAt:       3,
							UpdatedAt:       4,
							EntityType:      proto.Tag_FEATURE_FLAG,
							EnvironmentId:   "env-1",
							EnvironmentName: "env-name-1",
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
				// Use ElementsMatch because map iteration order is non-deterministic
				assert.ElementsMatch(t, p.expectedEnvTags, envTags)
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
			expectedErr: tagstorage.ErrTagUnexpectedAffectedRows,
		},
		{
			desc: "ErrRowsAffected",
			setup: func(s *tagStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), errors.New("rows affected error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			id:          "tag-id-0",
			expectedErr: errors.New("rows affected error"),
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
