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

package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/authz"
	err "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	featurestoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	databasemock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/tag/domain"
	tagstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/tag/storage/mock"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	"github.com/bucketeer-io/bucketeer/v2/proto/feature"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/tag"
)

func TestNewTagService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	dbClientMock := databasemock.NewMockClient(mockController)
	tagStorageMock := tagstoragemock.NewMockTagStorage(mockController)
	p := publishermock.NewMockPublisher(mockController)
	logger := zap.NewNop()
	featureStorageMock := featurestoragemock.NewMockFeatureStorage(mockController)
	s := NewTagService(
		dbClientMock,
		tagStorageMock,
		featureStorageMock,
		p,
		WithLogger(logger),
	)
	assert.IsType(t, &TagService{}, s)
}

func TestCreateTagMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithEditor(t, &eventproto.Editor{Email: "email"})

	patterns := []struct {
		desc        string
		setup       func(*TagService)
		req         *proto.CreateTagRequest
		expectedErr error
		expectedTag *proto.Tag
	}{
		{
			desc: "err: ErrNameRequired",
			req: &proto.CreateTagRequest{
				EnvironmentId: "ns0",
				EntityType:    proto.Tag_FEATURE_FLAG,
			},
			expectedErr: statusNameRequired.Err(),
			expectedTag: nil,
		},
		{
			desc: "err: ErrEntityTypeRequired",
			req: &proto.CreateTagRequest{
				EnvironmentId: "ns0",
				Name:          "test-tag",
			},
			expectedErr: statusEntityTypeRequired.Err(),
			expectedTag: nil,
		},
		{
			desc: "err: UpsertTag fails",
			setup: func(s *TagService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().GetTagByName(
					gomock.Any(), "test-tag", "ns0", proto.Tag_FEATURE_FLAG,
				).Return(&domain.Tag{
					Tag: &proto.Tag{
						Id:            "actual-tag-id",
						Name:          "test-tag",
						CreatedAt:     1000,
						UpdatedAt:     2000,
						EntityType:    proto.Tag_FEATURE_FLAG,
						EnvironmentId: "ns0",
					},
				}, nil)
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().UpsertTag(
					gomock.Any(), gomock.Any(),
				).Return(err.NewErrorInternal(err.TagPackageName, "storage error"))
			},
			req: &proto.CreateTagRequest{
				EnvironmentId: "ns0",
				Name:          "test-tag",
				EntityType:    proto.Tag_FEATURE_FLAG,
			},
			expectedErr: api.NewGRPCStatus(err.NewErrorInternal(err.TagPackageName, "storage error")).Err(),
			expectedTag: nil,
		},
		{
			desc: "err: GetTagByName fails",
			setup: func(s *TagService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().GetTagByName(
					gomock.Any(), "test-tag", "ns0", proto.Tag_FEATURE_FLAG,
				).Return(nil, err.NewErrorNotFound(err.TagPackageName, "tag name not found", "tag_name"))
			},
			req: &proto.CreateTagRequest{
				EnvironmentId: "ns0",
				Name:          "test-tag",
				EntityType:    proto.Tag_FEATURE_FLAG,
			},
			expectedErr: api.NewGRPCStatus(err.NewErrorNotFound(err.TagPackageName, "tag name not found", "tag_name")).Err(),
			expectedTag: nil,
		},
		{
			desc: "err: Publish fails",
			setup: func(s *TagService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().GetTagByName(
					gomock.Any(), "test-tag", "ns0", proto.Tag_FEATURE_FLAG,
				).Return(&domain.Tag{
					Tag: &proto.Tag{
						Id:            "actual-tag-id",
						Name:          "test-tag",
						CreatedAt:     1000,
						UpdatedAt:     2000,
						EntityType:    proto.Tag_FEATURE_FLAG,
						EnvironmentId: "ns0",
					},
				}, nil).Times(2)
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().UpsertTag(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(err.NewErrorInternal(err.TagPackageName, "publish error"))
			},
			req: &proto.CreateTagRequest{
				EnvironmentId: "ns0",
				Name:          "test-tag",
				EntityType:    proto.Tag_FEATURE_FLAG,
			},
			expectedErr: api.NewGRPCStatus(err.NewErrorInternal(err.TagPackageName, "publish error")).Err(),
			expectedTag: nil,
		},
		{
			desc: "success: new tag creation",
			setup: func(s *TagService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().UpsertTag(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().GetTagByName(
					gomock.Any(), "test-tag", "ns0", proto.Tag_FEATURE_FLAG,
				).Return(&domain.Tag{
					Tag: &proto.Tag{
						Id:            "actual-tag-id",
						Name:          "test-tag",
						CreatedAt:     1000,
						UpdatedAt:     1000,
						EntityType:    proto.Tag_FEATURE_FLAG,
						EnvironmentId: "ns0",
					},
				}, nil).Times(2)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.CreateTagRequest{
				EnvironmentId: "ns0",
				Name:          "test-tag",
				EntityType:    proto.Tag_FEATURE_FLAG,
			},
			expectedErr: nil,
			expectedTag: &proto.Tag{
				Id:            "actual-tag-id",
				Name:          "test-tag",
				CreatedAt:     1000,
				UpdatedAt:     1000,
				EntityType:    proto.Tag_FEATURE_FLAG,
				EnvironmentId: "ns0",
			},
		},
		{
			desc: "success: tag upsert (existing tag)",
			setup: func(s *TagService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})
				// Simulate upsert of existing tag - same ID and created_at, but updated updated_at
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().GetTagByName(
					gomock.Any(), "existing-tag", "ns0", proto.Tag_FEATURE_FLAG,
				).Return(&domain.Tag{
					Tag: &proto.Tag{
						Id:            "original-tag-id",
						Name:          "existing-tag",
						CreatedAt:     1000, // Original creation time
						UpdatedAt:     2000, // Updated time
						EntityType:    proto.Tag_FEATURE_FLAG,
						EnvironmentId: "ns0",
					},
				}, nil).Times(2)
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().UpsertTag(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.CreateTagRequest{
				EnvironmentId: "ns0",
				Name:          "existing-tag",
				EntityType:    proto.Tag_FEATURE_FLAG,
			},
			expectedErr: nil,
			expectedTag: &proto.Tag{
				Id:            "original-tag-id",
				Name:          "existing-tag",
				CreatedAt:     1000, // Original creation time
				UpdatedAt:     2000, // Updated time
				EntityType:    proto.Tag_FEATURE_FLAG,
				EnvironmentId: "ns0",
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createTagService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			resp, err := s.CreateTag(ctx, p.req)
			if p.expectedErr == nil {
				assert.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, p.expectedTag, resp.Tag)
			} else {
				assert.Error(t, err)
				expectedStatus, _ := gstatus.FromError(p.expectedErr)
				actualStatus, _ := gstatus.FromError(err)
				assert.Equal(t, expectedStatus.Code(), actualStatus.Code())
				assert.Equal(t, expectedStatus.Message(), actualStatus.Message())
			}
		})
	}
}

func TestListTagsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithEditor(t, &eventproto.Editor{Email: "email"})

	patterns := []struct {
		desc        string
		service     *TagService
		setup       func(*TagService)
		req         *proto.ListTagsRequest
		expectedErr error
		expected    *proto.ListTagsResponse
	}{
		{
			desc:    "success",
			service: createTagService(mockController),
			setup: func(s *TagService) {
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().ListTags(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Tag{
					{
						Id:            "tag-0",
						Name:          "test-tag",
						EnvironmentId: "ns0",
						EntityType:    proto.Tag_FEATURE_FLAG,
					},
				}, 0, int64(0), nil)
			},
			req: &proto.ListTagsRequest{
				EnvironmentId: "ns0",
				PageSize:      10,
				Cursor:        "0",
			},
			expectedErr: nil,
			expected: &proto.ListTagsResponse{
				Tags: []*proto.Tag{
					{
						Id:            "tag-0",
						Name:          "test-tag",
						EnvironmentId: "ns0",
						EntityType:    proto.Tag_FEATURE_FLAG,
					},
				},
				Cursor:     "0",
				TotalCount: 0,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := p.service
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.ListTags(ctx, p.req)
			if p.expectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, p.expected, resp)
			} else {
				assert.Error(t, err)
				expectedStatus, _ := gstatus.FromError(p.expectedErr)
				actualStatus, _ := gstatus.FromError(err)
				assert.Equal(t, expectedStatus.Code(), actualStatus.Code())
				assert.Equal(t, expectedStatus.Message(), actualStatus.Message())
			}
		})
	}
}

func TestDeleteTagMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithEditor(t, &eventproto.Editor{Email: "email"})

	patterns := []struct {
		desc        string
		setup       func(*TagService)
		req         *proto.DeleteTagRequest
		expectedErr error
	}{
		{
			desc: "err: ErrIDRequired",
			req: &proto.DeleteTagRequest{
				EnvironmentId: "ns0",
			},
			expectedErr: statusNameRequired.Err(),
		},
		{
			desc: "err: GetTag",
			setup: func(s *TagService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().GetTag(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, err.NewErrorNotFound(err.TagPackageName, "tag not found", "tag"))
			},
			req: &proto.DeleteTagRequest{
				Id:            "tag-0",
				EnvironmentId: "ns0",
			},
			expectedErr: api.NewGRPCStatus(err.NewErrorNotFound(err.TagPackageName, "tag not found", "tag")).Err(),
		},
		{
			desc: "err: in used",
			setup: func(s *TagService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().GetTag(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Tag{
					Tag: &proto.Tag{
						Id:            "tag-0",
						Name:          "test-tag",
						EnvironmentId: "ns0",
						EntityType:    proto.Tag_FEATURE_FLAG,
					},
				}, nil)
				s.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return([]*feature.Feature{
					{
						Id:   "feature-0",
						Tags: []string{"test-tag"},
					},
				}, 0, int64(0), nil)
			},
			req: &proto.DeleteTagRequest{
				Id:            "tag-0",
				EnvironmentId: "ns0",
			},
			expectedErr: statusTagInUsed.Err(),
		},
		{
			desc: "success",
			setup: func(s *TagService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					_ = fn(ctx)
				}).Return(nil)
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().GetTag(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Tag{
					Tag: &proto.Tag{
						Id:            "tag-0",
						Name:          "test-tag",
						EnvironmentId: "ns0",
						EntityType:    proto.Tag_FEATURE_FLAG,
					},
				}, nil)
				s.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), nil)
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().DeleteTag(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.DeleteTagRequest{
				Id:            "tag-0",
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createTagService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.DeleteTag(ctx, p.req)
			if p.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				expectedStatus, _ := gstatus.FromError(p.expectedErr)
				actualStatus, _ := gstatus.FromError(err)
				assert.Equal(t, expectedStatus.Code(), actualStatus.Code())
				assert.Equal(t, expectedStatus.Message(), actualStatus.Message())
			}
		})
	}
}

func createTagService(c *gomock.Controller) *TagService {
	dbClientMock := databasemock.NewMockClient(c)
	p := publishermock.NewMockPublisher(c)
	logger := zap.NewNop()
	return &TagService{
		dbClient:       dbClientMock,
		tagStorage:     tagstoragemock.NewMockTagStorage(c),
		featureStorage: featurestoragemock.NewMockFeatureStorage(c),
		publisher:      p,
		opts: &options{
			logger: zap.NewNop(),
		},
		logger: logger,
	}
}

// createContextWithEditor is a context after the authz interceptor allowed the request.
func createContextWithEditor(t *testing.T, editor *eventproto.Editor) context.Context {
	t.Helper()
	return authz.WithEditor(context.TODO(), editor)
}
