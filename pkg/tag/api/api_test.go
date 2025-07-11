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

package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	featurestoragemock "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/tag/domain"
	tagstoragemock "github.com/bucketeer-io/bucketeer/pkg/tag/storage/mock"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	"github.com/bucketeer-io/bucketeer/proto/feature"
	proto "github.com/bucketeer-io/bucketeer/proto/tag"
)

func TestNewTagService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	mysqlClientMock := mysqlmock.NewMockClient(mockController)
	accountClientMock := accountclientmock.NewMockClient(mockController)
	p := publishermock.NewMockPublisher(mockController)
	logger := zap.NewNop()
	s := NewTagService(
		mysqlClientMock,
		accountClientMock,
		p,
		WithLogger(logger),
	)
	assert.IsType(t, &TagService{}, s)
}

func TestCreateTagMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

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
			expectedErr: createError(statusNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
			expectedTag: nil,
		},
		{
			desc: "err: ErrEntityTypeRequired",
			req: &proto.CreateTagRequest{
				EnvironmentId: "ns0",
				Name:          "test-tag",
			},
			expectedErr: createError(statusEntityTypeRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "entity_type")),
			expectedTag: nil,
		},
		{
			desc: "err: UpsertTag fails",
			setup: func(s *TagService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(context.Context, mysql.Transaction) error) error {
					return fn(ctx, nil)
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
				).Return(errors.New("storage error"))
			},
			req: &proto.CreateTagRequest{
				EnvironmentId: "ns0",
				Name:          "test-tag",
				EntityType:    proto.Tag_FEATURE_FLAG,
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
			expectedTag: nil,
		},
		{
			desc: "err: GetTagByName fails",
			setup: func(s *TagService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(context.Context, mysql.Transaction) error) error {
					return fn(ctx, nil)
				})
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().UpsertTag(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().GetTagByName(
					gomock.Any(), "test-tag", "ns0", proto.Tag_FEATURE_FLAG,
				).Return(nil, errors.New("get error"))
			},
			req: &proto.CreateTagRequest{
				EnvironmentId: "ns0",
				Name:          "test-tag",
				EntityType:    proto.Tag_FEATURE_FLAG,
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
			expectedTag: nil,
		},
		{
			desc: "err: Publish fails",
			setup: func(s *TagService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(context.Context, mysql.Transaction) error) error {
					return fn(ctx, nil)
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
						UpdatedAt:     2000,
						EntityType:    proto.Tag_FEATURE_FLAG,
						EnvironmentId: "ns0",
					},
				}, nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("publish error"))
			},
			req: &proto.CreateTagRequest{
				EnvironmentId: "ns0",
				Name:          "test-tag",
				EntityType:    proto.Tag_FEATURE_FLAG,
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
			expectedTag: nil,
		},
		{
			desc: "success: new tag creation",
			setup: func(s *TagService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(context.Context, mysql.Transaction) error) error {
					return fn(ctx, nil)
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
				}, nil)
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(context.Context, mysql.Transaction) error) error {
					return fn(ctx, nil)
				})
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().UpsertTag(
					gomock.Any(), gomock.Any(),
				).Return(nil)
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
				}, nil)
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
			assert.Equal(t, p.expectedErr, err)
			if p.expectedErr == nil {
				require.NotNil(t, resp)
				assert.Equal(t, p.expectedTag, resp.Tag)
			}
		})
	}
}

func TestListTagsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleUnassigned(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		service     *TagService
		setup       func(*TagService)
		req         *proto.ListTagsRequest
		expectedErr error
		expected    *proto.ListTagsResponse
	}{
		{
			desc:        "errPermissionDenied",
			service:     createServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			setup:       func(s *TagService) {},
			req:         &proto.ListTagsRequest{EnvironmentId: "ns0", PageSize: 10, Cursor: "0"},
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
			expected:    nil,
		},
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
			assert.Equal(t, p.expectedErr, err)
			if p.expectedErr == nil {
				assert.Equal(t, p.expected, resp)
			}
		})
	}
}

func TestDeleteTagMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

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
			expectedErr: createError(statusNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: GetTag",
			setup: func(s *TagService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) error {
					return fn(ctx, nil)
				})
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().GetTag(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			req: &proto.DeleteTagRequest{
				Id:            "tag-0",
				EnvironmentId: "ns0",
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "err: in used",
			setup: func(s *TagService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) error {
					return fn(ctx, nil)
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
			expectedErr: createError(statusTagInUsed, localizer.MustLocalize(locale.Tag)),
		},
		{
			desc: "success",
			setup: func(s *TagService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
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
					gomock.Any(), gomock.Any(),
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
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func createTagService(c *gomock.Controller) *TagService {
	mysqlClientMock := mysqlmock.NewMockClient(c)
	accountClientMock := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountV2ByEnvironmentIDResponse{
		Account: &accountproto.AccountV2{
			Email:            "email",
			OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "ns0",
					Role:          accountproto.AccountV2_Role_Environment_EDITOR,
				},
				{
					EnvironmentId: "",
					Role:          accountproto.AccountV2_Role_Environment_EDITOR,
				},
			},
		},
	}
	accountClientMock.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	p := publishermock.NewMockPublisher(c)
	logger := zap.NewNop()
	return &TagService{
		mysqlClient:    mysqlClientMock,
		tagStorage:     tagstoragemock.NewMockTagStorage(c),
		featureStorage: featurestoragemock.NewMockFeatureStorage(c),
		accountClient:  accountClientMock,
		publisher:      p,
		opts: &options{
			logger: zap.NewNop(),
		},
		logger: logger,
	}
}

func createServiceWithGetAccountByEnvironmentMock(c *gomock.Controller, ro accountproto.AccountV2_Role_Organization, re accountproto.AccountV2_Role_Environment) *TagService {
	mysqlClientMock := mysqlmock.NewMockClient(c)
	accountClientMock := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountV2ByEnvironmentIDResponse{
		Account: &accountproto.AccountV2{
			Email:            "email",
			OrganizationRole: ro,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "ns0",
					Role:          re,
				},
			},
		},
	}
	accountClientMock.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	p := publishermock.NewMockPublisher(c)
	logger := zap.NewNop()
	return &TagService{
		mysqlClient:   mysqlClientMock,
		tagStorage:    tagstoragemock.NewMockTagStorage(c),
		accountClient: accountClientMock,
		publisher:     p,
		opts: &options{
			logger: zap.NewNop(),
		},
		logger: logger,
	}
}

func createContextWithTokenRoleUnassigned(t *testing.T) context.Context {
	t.Helper()
	token := &token.AccessToken{
		Issuer:   "issuer",
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "email",
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}

func createContextWithTokenRoleOwner(t *testing.T) context.Context {
	t.Helper()
	token := &token.AccessToken{
		Issuer:        "issuer",
		Audience:      "audience",
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         "email",
		IsSystemAdmin: true,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}
