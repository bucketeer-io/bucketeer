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
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/tag/domain"
	tagstoragemock "github.com/bucketeer-io/bucketeer/pkg/tag/storage/mock"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
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
	}{
		{
			desc: "err: ErrNameRequired",
			req: &proto.CreateTagRequest{
				EnvironmentId: "ns0",
				EntityType:    proto.Tag_FEATURE_FLAG,
			},
			expectedErr: createError(statusNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc: "err: ErrEntityTypeRequired",
			req: &proto.CreateTagRequest{
				EnvironmentId: "ns0",
				Name:          "test-tag",
			},
			expectedErr: createError(statusEntityTypeRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "entity_type")),
		},
		{
			desc: "success",
			setup: func(s *TagService) {
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().UpsertTag(
					gomock.Any(), gomock.Any(),
				).Return(nil)
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
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createTagService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.CreateTag(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
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
	}{
		{
			desc:        "errPermissionDenied",
			service:     createServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			setup:       func(s *TagService) {},
			req:         &proto.ListTagsRequest{EnvironmentId: "ns0", PageSize: 10, Cursor: "0"},
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc:    "success",
			service: createTagService(mockController),
			setup: func(s *TagService) {
				s.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().ListTags(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*proto.Tag{}, 0, int64(0), nil)
			},
			req: &proto.ListTagsRequest{
				EnvironmentId: "ns0",
				PageSize:      10,
				Cursor:        "0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := p.service
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.ListTags(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
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
			desc: "success",
			setup: func(s *TagService) {
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
