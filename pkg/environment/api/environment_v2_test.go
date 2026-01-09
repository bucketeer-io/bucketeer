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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"

	acmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2/mock"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	ftstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestGetEnvironmentV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		id          string
		expectedErr error
	}{
		{
			desc: "err: ErrEnvironmentNotFound",
			setup: func(s *EnvironmentService) {
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().GetEnvironmentV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2es.ErrEnvironmentNotFound)
			},
			id:          "id-0",
			expectedErr: statusEnvironmentNotFound.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().GetEnvironmentV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error"))
			},
			id:          "id-1",
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().GetEnvironmentV2(
					gomock.Any(), gomock.Any(),
				).Return(&domain.EnvironmentV2{}, nil)
			},
			id:          "id-3",
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			s := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			req := &proto.GetEnvironmentV2Request{Id: p.id}
			resp, err := s.GetEnvironmentV2(ctx, req)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestListEnvironmentsV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		input       *proto.ListEnvironmentsV2Request
		expected    *proto.ListEnvironmentsV2Response
		expectedErr error
	}{
		{
			desc:        "err: ErrInvalidCursor",
			setup:       nil,
			input:       &proto.ListEnvironmentsV2Request{Cursor: "XXX"},
			expected:    nil,
			expectedErr: statusInvalidCursor.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error"))
			},
			input:       &proto.ListEnvironmentsV2Request{},
			expected:    nil,
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.EnvironmentV2{}, 0, int64(0), nil)
			},
			input:       &proto.ListEnvironmentsV2Request{PageSize: 2, Cursor: ""},
			expected:    &proto.ListEnvironmentsV2Response{Environments: []*proto.EnvironmentV2{}, Cursor: "0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.ListEnvironmentsV2(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestCreateEnvironmentV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	envExpectedTrue, err := domain.NewEnvironmentV2(
		"Env Name-dev01",
		"url-code-01",
		"description",
		"project-id01",
		"organization-id01",
		true,
		nil,
	)

	envExpectedFalse, err := domain.NewEnvironmentV2(
		"Env Name-dev01",
		"url-code-01",
		"description",
		"project-id01",
		"organization-id01",
		false,
		nil,
	)

	envExpectedTrue.Archived = false
	require.NoError(t, err)

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.CreateEnvironmentV2Request
		expected    *proto.EnvironmentV2
		expectedErr error
	}{
		{
			desc:  "err: ErrInvalidEnvironmentName: empty name",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: ""},
			},
			expectedErr: statusEnvironmentNameRequired.Err(),
		},
		{
			desc:  "err: ErrInvalidEnvironmentName: only space",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "    "},
			},
			expectedErr: statusEnvironmentNameRequired.Err(),
		},
		{
			desc:  "err: ErrInvalidEnvironmentName: max name length exceeded",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: strings.Repeat("a", 51)},
			},
			expectedErr: statusInvalidEnvironmentName.Err(),
		},
		{
			desc:  "err: ErrInvalidEnvironmentUrlCode: empty url code",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: ""},
			},
			expectedErr: statusInvalidEnvironmentUrlCode.Err(),
		},
		{
			desc:  "err: ErrInvalidEnvironmentUrlCode: can't use uppercase",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: "URLCODE"},
			},
			expectedErr: statusInvalidEnvironmentUrlCode.Err(),
		},
		{
			desc:  "err: ErrInvalidEnvironmentUrlCode: can't use space",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: "url code"},
			},
			expectedErr: statusInvalidEnvironmentUrlCode.Err(),
		},
		{
			desc:  "err: ErrInvalidEnvironmentUrlCode: max url code length exceeded",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: strings.Repeat("a", 51)},
			},
			expectedErr: statusInvalidEnvironmentUrlCode.Err(),
		},
		{
			desc:  "err: ErrProjectIDRequired",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: "url-code", ProjectId: ""},
			},
			expectedErr: statusProjectIDRequired.Err(),
		},
		{
			desc: "err: ErrProjectNotFound",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2es.ErrProjectNotFound)
			},
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: "url-code", ProjectId: "project-id"},
			},
			expectedErr: statusProjectNotFound.Err(),
		},
		{
			desc: "err: ErrEnvironmentAlreadyExists",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Project{
					Project: &proto.Project{Id: "project-id", OrganizationId: "organization-id01"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrEnvironmentAlreadyExists)
			},
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: "url-code", ProjectId: "project-id"},
			},
			expectedErr: statusEnvironmentAlreadyExists.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Project{
					Project: &proto.Project{Id: "project-id", OrganizationId: "organization-id01"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error"))
			},
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{Name: "name", UrlCode: "url-code", ProjectId: "project-id"},
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error")).Err(),
		},
		{
			desc: "success: require comment is true",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Project{
					Project: &proto.Project{Id: "project-id", OrganizationId: "organization-id01"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().CreateEnvironmentV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{
					Name:           envExpectedTrue.Name,
					UrlCode:        envExpectedTrue.UrlCode,
					Description:    envExpectedTrue.Description,
					ProjectId:      envExpectedTrue.ProjectId,
					RequireComment: true,
				},
			},
			expected: envExpectedTrue.EnvironmentV2,
		},
		{
			desc: "success: require comment is false",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Project{
					Project: &proto.Project{Id: "project-id", OrganizationId: "organization-id01"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().CreateEnvironmentV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.CreateEnvironmentV2Request{
				Command: &proto.CreateEnvironmentV2Command{
					Name:           envExpectedFalse.Name,
					UrlCode:        envExpectedFalse.UrlCode,
					Description:    envExpectedFalse.Description,
					ProjectId:      envExpectedFalse.ProjectId,
					RequireComment: false,
				},
			},
			expected: envExpectedFalse.EnvironmentV2,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.CreateEnvironmentV2(ctx, p.req)
			if resp != nil {
				assert.True(t, len(resp.Environment.Id) > 0)
				assert.Equal(t, p.expected.Name, resp.Environment.Name)
				assert.Equal(t, p.expected.UrlCode, resp.Environment.UrlCode)
				assert.Equal(t, p.expected.Description, resp.Environment.Description)
				assert.Equal(t, p.expected.ProjectId, resp.Environment.ProjectId)
				assert.Equal(t, p.expected.Archived, resp.Environment.Archived)
				assert.Equal(t, p.expected.RequireComment, resp.Environment.RequireComment)
				assert.True(t, resp.Environment.CreatedAt > 0)
				assert.True(t, resp.Environment.UpdatedAt > 0)
			}
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCreateEnvironmentV2NoCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	envExpectedTrue, err := domain.NewEnvironmentV2(
		"Env Name-dev01",
		"url-code-01",
		"description",
		"project-id01",
		"organization-id01",
		true,
		nil,
	)

	envExpectedFalse, err := domain.NewEnvironmentV2(
		"Env Name-dev01",
		"url-code-01",
		"description",
		"project-id01",
		"organization-id01",
		false,
		nil,
	)

	envExpectedTrue.Archived = false
	require.NoError(t, err)

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.CreateEnvironmentV2Request
		expected    *proto.EnvironmentV2
		expectedErr error
	}{
		{
			desc:  "err: ErrInvalidEnvironmentName: empty name",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Name: "",
			},
			expectedErr: statusEnvironmentNameRequired.Err(),
		},
		{
			desc:  "err: ErrInvalidEnvironmentName: only space",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Name: "    ",
			},
			expectedErr: statusEnvironmentNameRequired.Err(),
		},
		{
			desc:  "err: ErrInvalidEnvironmentName: max name length exceeded",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Name: strings.Repeat("a", 51),
			},
			expectedErr: statusInvalidEnvironmentName.Err(),
		},
		{
			desc:  "err: ErrInvalidEnvironmentUrlCode: empty url code",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Name:    "name",
				UrlCode: "",
			},
			expectedErr: statusInvalidEnvironmentUrlCode.Err(),
		},
		{
			desc:  "err: ErrInvalidEnvironmentUrlCode: can't use uppercase",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Name:    "name",
				UrlCode: "URLCODE",
			},
			expectedErr: statusInvalidEnvironmentUrlCode.Err(),
		},
		{
			desc:  "err: ErrInvalidEnvironmentUrlCode: can't use space",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Name:    "name",
				UrlCode: "url code",
			},
			expectedErr: statusInvalidEnvironmentUrlCode.Err(),
		},
		{
			desc:  "err: ErrInvalidEnvironmentUrlCode: max url code length exceeded",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Name:    "name",
				UrlCode: strings.Repeat("a", 51),
			},
			expectedErr: statusInvalidEnvironmentUrlCode.Err(),
		},
		{
			desc:  "err: ErrProjectIDRequired",
			setup: nil,
			req: &proto.CreateEnvironmentV2Request{
				Name:      "name",
				UrlCode:   "url-code",
				ProjectId: "",
			},
			expectedErr: statusProjectIDRequired.Err(),
		},
		{
			desc: "err: ErrProjectNotFound",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2es.ErrProjectNotFound)
			},
			req: &proto.CreateEnvironmentV2Request{
				Name:      "name",
				UrlCode:   "url-code",
				ProjectId: "project-id",
			},
			expectedErr: statusProjectNotFound.Err(),
		},
		{
			desc: "err: ErrEnvironmentAlreadyExists",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Project{
					Project: &proto.Project{Id: "project-id", OrganizationId: "organization-id01"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrEnvironmentAlreadyExists)
			},
			req: &proto.CreateEnvironmentV2Request{
				Name:      "name",
				UrlCode:   "url-code",
				ProjectId: "project-id",
			},
			expectedErr: statusEnvironmentAlreadyExists.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), gomock.Any(),
				).Return(nil, pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error"))
			},
			req: &proto.CreateEnvironmentV2Request{
				Name:      "name",
				UrlCode:   "url-code",
				ProjectId: "project-id",
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error")).Err(),
		},
		{
			desc: "success: require comment is true",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Project{
					Project: &proto.Project{Id: "project-id", OrganizationId: "organization-id01"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().CreateEnvironmentV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.CreateEnvironmentV2Request{
				Name:           envExpectedTrue.Name,
				UrlCode:        envExpectedTrue.UrlCode,
				Description:    envExpectedTrue.Description,
				ProjectId:      envExpectedTrue.ProjectId,
				RequireComment: true,
			},
			expected: envExpectedTrue.EnvironmentV2,
		},
		{
			desc: "success: require comment is false",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Project{
					Project: &proto.Project{Id: "project-id", OrganizationId: "organization-id01"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().CreateEnvironmentV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.CreateEnvironmentV2Request{
				Name:           envExpectedFalse.Name,
				UrlCode:        envExpectedFalse.UrlCode,
				Description:    envExpectedFalse.Description,
				ProjectId:      envExpectedFalse.ProjectId,
				RequireComment: false,
			},
			expected: envExpectedFalse.EnvironmentV2,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.CreateEnvironmentV2(ctx, p.req)
			if resp != nil {
				assert.True(t, len(resp.Environment.Id) > 0)
				assert.Equal(t, p.expected.Name, resp.Environment.Name)
				assert.Equal(t, p.expected.UrlCode, resp.Environment.UrlCode)
				assert.Equal(t, p.expected.Description, resp.Environment.Description)
				assert.Equal(t, p.expected.ProjectId, resp.Environment.ProjectId)
				assert.Equal(t, p.expected.Archived, resp.Environment.Archived)
				assert.Equal(t, p.expected.RequireComment, resp.Environment.RequireComment)
				assert.True(t, resp.Environment.CreatedAt > 0)
				assert.True(t, resp.Environment.UpdatedAt > 0)
			}
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateEnvironmentV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.UpdateEnvironmentV2Request
		expectedErr error
	}{
		{
			desc:  "err: ErrInvalidEnvironmentName: only space",
			setup: nil,
			req: &proto.UpdateEnvironmentV2Request{
				Id:                       "id01",
				RenameCommand:            &proto.RenameEnvironmentV2Command{Name: "  "},
				ChangeDescriptionCommand: &proto.ChangeDescriptionEnvironmentV2Command{Description: "desc-1"},
			},
			expectedErr: statusEnvironmentNameRequired.Err(),
		},
		{
			desc:  "err: ErrInvalidEnvironmentName: max name length exceeded",
			setup: nil,
			req: &proto.UpdateEnvironmentV2Request{
				Id:                       "id01",
				RenameCommand:            &proto.RenameEnvironmentV2Command{Name: strings.Repeat("a", 51)},
				ChangeDescriptionCommand: &proto.ChangeDescriptionEnvironmentV2Command{Description: "desc-1"},
			},
			expectedErr: statusInvalidEnvironmentName.Err(),
		},
		{
			desc: "err: ErrEnvironmentNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrEnvironmentNotFound)
			},
			req: &proto.UpdateEnvironmentV2Request{
				Id:            "id01",
				RenameCommand: &proto.RenameEnvironmentV2Command{Name: "name-0"},
			},
			expectedErr: statusEnvironmentNotFound.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error"))
			},
			req: &proto.UpdateEnvironmentV2Request{
				Id:            "id02",
				RenameCommand: &proto.RenameEnvironmentV2Command{Name: "name-1"},
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().GetEnvironmentV2(
					gomock.Any(), gomock.Any(),
				).Return(&domain.EnvironmentV2{
					EnvironmentV2: &proto.EnvironmentV2{},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().UpdateEnvironmentV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.UpdateEnvironmentV2Request{
				Id:                          "id01",
				RenameCommand:               &proto.RenameEnvironmentV2Command{Name: "name-1"},
				ChangeDescriptionCommand:    &proto.ChangeDescriptionEnvironmentV2Command{Description: "desc-1"},
				ChangeRequireCommentCommand: &proto.ChangeRequireCommentCommand{RequireComment: true},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UpdateEnvironmentV2(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateEnvironmentV2NoCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.UpdateEnvironmentV2Request
		expectedErr error
	}{
		{
			desc:  "err: ErrInvalidEnvironmentName: only space",
			setup: nil,
			req: &proto.UpdateEnvironmentV2Request{
				Id:          "id01",
				Name:        wrapperspb.String("  "),
				Description: wrapperspb.String("desc-1"),
			},
			expectedErr: statusEnvironmentNameRequired.Err(),
		},
		{
			desc:  "err: ErrInvalidEnvironmentName: max name length exceeded",
			setup: nil,
			req: &proto.UpdateEnvironmentV2Request{
				Id:          "id01",
				Name:        wrapperspb.String(strings.Repeat("a", 51)),
				Description: wrapperspb.String("desc-1"),
			},
			expectedErr: statusInvalidEnvironmentName.Err(),
		},
		{
			desc: "err: ErrEnvironmentNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrEnvironmentNotFound)
			},
			req: &proto.UpdateEnvironmentV2Request{
				Id:   "id01",
				Name: wrapperspb.String("name-0"),
			},
			expectedErr: statusEnvironmentNotFound.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error"))
			},
			req: &proto.UpdateEnvironmentV2Request{
				Id:   "id02",
				Name: wrapperspb.String("name-1"),
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.UpdateEnvironmentV2Request{
				Id:             "id01",
				Name:           wrapperspb.String("name-1"),
				Description:    wrapperspb.String("desc-1"),
				RequireComment: wrapperspb.Bool(true),
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UpdateEnvironmentV2(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestArchiveEnvironmentV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.ArchiveEnvironmentV2Request
		expectedErr error
	}{
		{
			desc: "err: ErrEnvironmentNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrEnvironmentNotFound)
			},
			req:         &proto.ArchiveEnvironmentV2Request{Id: "id01", Command: &proto.ArchiveEnvironmentV2Command{}},
			expectedErr: statusEnvironmentNotFound.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error"))
			},
			req:         &proto.ArchiveEnvironmentV2Request{Id: "id02", Command: &proto.ArchiveEnvironmentV2Command{}},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().GetEnvironmentV2(
					gomock.Any(), gomock.Any(),
				).Return(&domain.EnvironmentV2{
					EnvironmentV2: &proto.EnvironmentV2{},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().UpdateEnvironmentV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req:         &proto.ArchiveEnvironmentV2Request{Id: "id01", Command: &proto.ArchiveEnvironmentV2Command{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.ArchiveEnvironmentV2(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUnarchiveEnvironmentV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.UnarchiveEnvironmentV2Request
		expectedErr error
	}{
		{
			desc: "err: ErrEnvironmentNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrEnvironmentNotFound)
			},
			req:         &proto.UnarchiveEnvironmentV2Request{Id: "id01", Command: &proto.UnarchiveEnvironmentV2Command{}},
			expectedErr: statusEnvironmentNotFound.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error"))
			},
			req:         &proto.UnarchiveEnvironmentV2Request{Id: "id02", Command: &proto.UnarchiveEnvironmentV2Command{}},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req:         &proto.UnarchiveEnvironmentV2Request{Id: "id01", Command: &proto.UnarchiveEnvironmentV2Command{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UnarchiveEnvironmentV2(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestEnvironmentV2APIs_Unauthenticated(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	unauthCtx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
		"accept-language": []string{"ja"},
	})

	expectedErr := statusUnauthenticated.Err()

	patterns := []struct {
		desc      string
		setupFunc func(*EnvironmentService)
		testFunc  func(*EnvironmentService) error
	}{
		{
			desc: "CreateEnvironmentV2 - unauthenticated",
			setupFunc: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Project{
					Project: &proto.Project{Id: "project-id", OrganizationId: "organization-id01"},
				}, nil)
			},
			testFunc: func(s *EnvironmentService) error {
				_, err := s.CreateEnvironmentV2(unauthCtx, &proto.CreateEnvironmentV2Request{
					Command: &proto.CreateEnvironmentV2Command{
						Name:      "name",
						UrlCode:   "url-code",
						ProjectId: "project-id",
					},
				})
				return err
			},
		},
		{
			desc: "CreateEnvironmentV2NoCommand - unauthenticated",
			setupFunc: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Project{
					Project: &proto.Project{Id: "project-id", OrganizationId: "organization-id01"},
				}, nil)
			},
			testFunc: func(s *EnvironmentService) error {
				_, err := s.CreateEnvironmentV2(unauthCtx, &proto.CreateEnvironmentV2Request{
					Name:      "name",
					UrlCode:   "url-code",
					ProjectId: "project-id",
				})
				return err
			},
		},
		{
			desc: "UpdateEnvironmentV2 - unauthenticated",
			setupFunc: func(s *EnvironmentService) {
				// For unauthenticated users, the error should occur before any storage calls
			},
			testFunc: func(s *EnvironmentService) error {
				_, err := s.UpdateEnvironmentV2(unauthCtx, &proto.UpdateEnvironmentV2Request{
					Id:            "env-id",
					RenameCommand: &proto.RenameEnvironmentV2Command{Name: "new-name"},
				})
				return err
			},
		},
		{
			desc: "ArchiveEnvironmentV2 - unauthenticated",
			setupFunc: func(s *EnvironmentService) {
				// For unauthenticated users, the error should occur before any storage calls
			},
			testFunc: func(s *EnvironmentService) error {
				_, err := s.ArchiveEnvironmentV2(unauthCtx, &proto.ArchiveEnvironmentV2Request{
					Id:      "env-id",
					Command: &proto.ArchiveEnvironmentV2Command{},
				})
				return err
			},
		},
		{
			desc: "UnarchiveEnvironmentV2 - unauthenticated",
			setupFunc: func(s *EnvironmentService) {
				// For unauthenticated users, the error should occur before any storage calls
			},
			testFunc: func(s *EnvironmentService) error {
				_, err := s.UnarchiveEnvironmentV2(unauthCtx, &proto.UnarchiveEnvironmentV2Request{
					Id:      "env-id",
					Command: &proto.UnarchiveEnvironmentV2Command{},
				})
				return err
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			service := newEnvironmentService(t, mockController, nil)
			if p.setupFunc != nil {
				p.setupFunc(service)
			}
			err := p.testFunc(service)
			assert.Equal(t, expectedErr, err)
		})
	}
}

func TestEnvironmentV2APIs_PermissionDenied(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	roleTestCtx := metadata.NewIncomingContext(createContextWithTokenRoleUnassigned(t), metadata.MD{
		"accept-language": []string{"ja"},
	})

	expectedErr := statusPermissionDenied.Err()

	rolePatterns := []struct {
		role      string
		roleValue accountproto.AccountV2_Role_Organization
	}{
		{"member", accountproto.AccountV2_Role_Organization_MEMBER},
		{"unassigned", accountproto.AccountV2_Role_Organization_UNASSIGNED},
	}

	for _, rolePattern := range rolePatterns {
		t.Run(rolePattern.role, func(t *testing.T) {
			patterns := []struct {
				desc      string
				setupFunc func(*EnvironmentService)
				testFunc  func(*EnvironmentService) error
			}{
				{
					desc: "CreateEnvironmentV2 - permission denied",
					setupFunc: func(s *EnvironmentService) {
						s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
							gomock.Any(), gomock.Any(),
						).Return(&domain.Project{
							Project: &proto.Project{Id: "project-id", OrganizationId: "organization-id01"},
						}, nil)
						s.accountClient.(*acmock.MockClient).EXPECT().GetAccountV2(
							gomock.Any(), &accountproto.GetAccountV2Request{
								Email:          "email",
								OrganizationId: "organization-id01",
							},
						).Return(&accountproto.GetAccountV2Response{
							Account: &accountproto.AccountV2{
								OrganizationRole: rolePattern.roleValue,
							},
						}, nil)
					},
					testFunc: func(s *EnvironmentService) error {
						_, err := s.CreateEnvironmentV2(roleTestCtx, &proto.CreateEnvironmentV2Request{
							Command: &proto.CreateEnvironmentV2Command{
								Name:      "name",
								UrlCode:   "url-code",
								ProjectId: "project-id",
							},
						})
						return err
					},
				},
				{
					desc: "CreateEnvironmentV2NoCommand - permission denied",
					setupFunc: func(s *EnvironmentService) {
						s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
							gomock.Any(), gomock.Any(),
						).Return(&domain.Project{
							Project: &proto.Project{Id: "project-id", OrganizationId: "organization-id01"},
						}, nil)
						s.accountClient.(*acmock.MockClient).EXPECT().GetAccountV2(
							gomock.Any(), &accountproto.GetAccountV2Request{
								Email:          "email",
								OrganizationId: "organization-id01",
							},
						).Return(&accountproto.GetAccountV2Response{
							Account: &accountproto.AccountV2{
								OrganizationRole: rolePattern.roleValue,
							},
						}, nil)
					},
					testFunc: func(s *EnvironmentService) error {
						_, err := s.CreateEnvironmentV2(roleTestCtx, &proto.CreateEnvironmentV2Request{
							Name:      "name",
							UrlCode:   "url-code",
							ProjectId: "project-id",
						})
						return err
					},
				},
				{
					desc: "UpdateEnvironmentV2 - permission denied",
					setupFunc: func(s *EnvironmentService) {
						// Mock the QueryRowContext call that GetAccountV2ByEnvironmentID will make
						row := mysqlmock.NewMockRow(mockController)
						row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
							// Populate the account struct with the test role
							if len(args) >= 12 {
								*(args[11].(*int32)) = int32(rolePattern.roleValue)
							}
							return nil
						})
						s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
							gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
						).Return(row)
					},
					testFunc: func(s *EnvironmentService) error {
						_, err := s.UpdateEnvironmentV2(roleTestCtx, &proto.UpdateEnvironmentV2Request{
							Id:            "env-id",
							RenameCommand: &proto.RenameEnvironmentV2Command{Name: "new-name"},
						})
						return err
					},
				},
				{
					desc: "ArchiveEnvironmentV2 - permission denied",
					setupFunc: func(s *EnvironmentService) {
						// Mock the QueryRowContext call that GetAccountV2ByEnvironmentID will make
						row := mysqlmock.NewMockRow(mockController)
						row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
							// Populate the account struct with the test role
							if len(args) >= 12 {
								*(args[11].(*int32)) = int32(rolePattern.roleValue)
							}
							return nil
						})
						s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
							gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
						).Return(row)
					},
					testFunc: func(s *EnvironmentService) error {
						_, err := s.ArchiveEnvironmentV2(roleTestCtx, &proto.ArchiveEnvironmentV2Request{
							Id:      "env-id",
							Command: &proto.ArchiveEnvironmentV2Command{},
						})
						return err
					},
				},
				{
					desc: "UnarchiveEnvironmentV2 - permission denied",
					setupFunc: func(s *EnvironmentService) {
						// Mock the QueryRowContext call that GetAccountV2ByEnvironmentID will make
						row := mysqlmock.NewMockRow(mockController)
						row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
							// Populate the account struct with the test role
							if len(args) >= 12 {
								*(args[11].(*int32)) = int32(rolePattern.roleValue)
							}
							return nil
						})
						s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
							gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
						).Return(row)
					},
					testFunc: func(s *EnvironmentService) error {
						_, err := s.UnarchiveEnvironmentV2(roleTestCtx, &proto.UnarchiveEnvironmentV2Request{
							Id:      "env-id",
							Command: &proto.UnarchiveEnvironmentV2Command{},
						})
						return err
					},
				},
			}

			for _, p := range patterns {
				t.Run(p.desc, func(t *testing.T) {
					t.Parallel()
					service := newEnvironmentService(t, mockController, nil)
					if p.setupFunc != nil {
						p.setupFunc(service)
					}
					err := p.testFunc(service)
					assert.Equal(t, expectedErr, err)
				})
			}
		})
	}
}

func TestDeleteEnvironmentV2Data(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.DeleteEnvironmentDataRequest
		expectedErr error
		expected    *proto.DeleteEnvironmentDataResponse
	}{
		{
			desc: "Environment not found",
			setup: func(s *EnvironmentService) {
				for range 11 {
					s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().
						CountTargetEntitiesInEnvironmentV2(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(int64(0), nil)
				}
				s.fluiStorage.(*ftstoragemock.MockFeatureLastUsedInfoStorage).EXPECT().
					SelectFeatureLastUsedInfos(gomock.Any(), gomock.Any()).
					Return([]*ftproto.FeatureLastUsedInfo{}, nil)
				for range targetEntitiesInEnvironment {
					s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().
						DeleteTargetFromEnvironmentV2(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil)
				}
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().
					DeleteEnvironmentV2(gomock.Any(), gomock.Any()).
					Return(statusEnvironmentNotFound.Err())
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(statusEnvironmentNotFound.Err())
			},
			req: &proto.DeleteEnvironmentDataRequest{
				EnvironmentIds: []string{"env-id"},
			},
			expectedErr: statusEnvironmentNotFound.Err(),
			expected:    nil,
		},
		{
			desc: "success dry run",
			setup: func(s *EnvironmentService) {
				for range 11 {
					s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().
						CountTargetEntitiesInEnvironmentV2(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(int64(1), nil)
				}
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().
					Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			req: &proto.DeleteEnvironmentDataRequest{
				EnvironmentIds: []string{"env-id"},
				DryRun:         true,
			},
			expectedErr: nil,
			expected: &proto.DeleteEnvironmentDataResponse{
				Summaries: []*proto.EnvironmentDeletionSummary{
					{
						EnvironmentId:               "env-id",
						FeaturesDeleted:             1,
						ExperimentsDeleted:          1,
						SubscriptionsDeleted:        1,
						PushesDeleted:               1,
						TagsDeleted:                 1,
						SegmentsDeleted:             1,
						FlagTriggersDeleted:         1,
						ApiKeysDeleted:              1,
						OperationsDeleted:           1,
						FeatureLastUsedInfosDeleted: 1,
						GoalsDeleted:                1,
					},
				},
			},
		},
		{
			desc: "success no force",
			setup: func(s *EnvironmentService) {
				for range 11 {
					s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().
						CountTargetEntitiesInEnvironmentV2(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(int64(1), nil)
				}
				s.fluiStorage.(*ftstoragemock.MockFeatureLastUsedInfoStorage).EXPECT().
					SelectFeatureLastUsedInfos(gomock.Any(), gomock.Any()).
					Return([]*ftproto.FeatureLastUsedInfo{}, nil)
				for range targetEntitiesInEnvironment {
					s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().
						DeleteTargetFromEnvironmentV2(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil)
				}
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().
					DeleteEnvironmentV2(gomock.Any(), gomock.Any()).
					Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().
					Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			req: &proto.DeleteEnvironmentDataRequest{
				EnvironmentIds: []string{"env-id"},
			},
			expectedErr: nil,
			expected: &proto.DeleteEnvironmentDataResponse{
				Summaries: []*proto.EnvironmentDeletionSummary{
					{
						EnvironmentId:               "env-id",
						FeaturesDeleted:             1,
						ExperimentsDeleted:          1,
						SubscriptionsDeleted:        1,
						PushesDeleted:               1,
						TagsDeleted:                 1,
						SegmentsDeleted:             1,
						FlagTriggersDeleted:         1,
						ApiKeysDeleted:              1,
						OperationsDeleted:           1,
						FeatureLastUsedInfosDeleted: 1,
						GoalsDeleted:                1,
					},
				},
			},
		},
		{
			desc: "success with force",
			setup: func(s *EnvironmentService) {
				for range 11 {
					s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().
						CountTargetEntitiesInEnvironmentV2(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(int64(1), nil)
				}
				for range targetEntitiesInEnvironment {
					s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().
						DeleteTargetFromEnvironmentV2(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil)
				}
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().
					DeleteEnvironmentV2(gomock.Any(), gomock.Any()).
					Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().
					Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			req: &proto.DeleteEnvironmentDataRequest{
				EnvironmentIds: []string{"env-id"},
				Force:          true,
			},
			expectedErr: nil,
			expected: &proto.DeleteEnvironmentDataResponse{
				Summaries: []*proto.EnvironmentDeletionSummary{
					{
						EnvironmentId:               "env-id",
						FeaturesDeleted:             1,
						ExperimentsDeleted:          1,
						SubscriptionsDeleted:        1,
						PushesDeleted:               1,
						TagsDeleted:                 1,
						SegmentsDeleted:             1,
						FlagTriggersDeleted:         1,
						ApiKeysDeleted:              1,
						OperationsDeleted:           1,
						FeatureLastUsedInfosDeleted: 1,
						GoalsDeleted:                1,
					},
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.DeleteEnvironmentData(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, resp)
		})
	}
}
