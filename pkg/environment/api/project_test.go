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
	"strings"
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	gstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	acmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2/mock"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	pubmock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

func TestGetProjectMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
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
		setup       func(*EnvironmentService)
		id          string
		expectedErr error
	}{
		{
			desc:        "err: ErrProjectIDRequired",
			setup:       nil,
			id:          "",
			expectedErr: createError(statusProjectIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrProjectNotFound",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "err-id-0",
				).Return(nil, v2es.ErrProjectNotFound)
			},
			id:          "err-id-0",
			expectedErr: createError(statusProjectNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "err-id-1",
				).Return(nil, pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
			},
			id:          "err-id-1",
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "success-id-0",
				).Return(&domain.Project{
					Project: &proto.Project{Id: "success-id-0", OrganizationId: "org-1"},
				}, nil)
			},
			id:          "success-id-0",
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			req := &proto.GetProjectRequest{Id: p.id}
			resp, err := s.GetProject(ctx, req)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestListProjectsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithToken(t)
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
		setup       func(*EnvironmentService)
		input       *proto.ListProjectsRequest
		expected    *proto.ListProjectsResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrInvalidCursor",
			setup:       nil,
			input:       &proto.ListProjectsRequest{Cursor: "XXX"},
			expected:    nil,
			expectedErr: createError(statusInvalidCursor, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor")),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().ListProjects(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
			},
			input:       &proto.ListProjectsRequest{},
			expected:    nil,
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().ListProjects(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Project{}, 0, int64(0), nil)
			},
			input:       &proto.ListProjectsRequest{PageSize: 2, Cursor: "", OrganizationIds: []string{"org-1", "org-2"}},
			expected:    &proto.ListProjectsResponse{Projects: []*proto.Project{}, Cursor: "0", TotalCount: 0},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.ListProjects(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestCreateProjectMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
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

	projExpected, err := domain.NewProject(
		"name",
		"url-code",
		"description",
		"email",
		"organizationID",
		false,
	)
	require.NoError(t, err)

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.CreateProjectRequest
		expected    *proto.Project
		expectedErr error
	}{
		{
			desc:  "err: ErrInvalidProjectName: empty name",
			setup: nil,
			req: &proto.CreateProjectRequest{
				Command: &proto.CreateProjectCommand{Name: ""},
			},
			expectedErr: createError(statusProjectNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc:  "err: ErrInvalidProjectName: only space",
			setup: nil,
			req: &proto.CreateProjectRequest{
				Command: &proto.CreateProjectCommand{Name: "    "},
			},
			expectedErr: createError(statusProjectNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc:  "err: ErrInvalidProjectName: max name length exceeded",
			setup: nil,
			req: &proto.CreateProjectRequest{
				Command: &proto.CreateProjectCommand{Name: strings.Repeat("a", 51)},
			},
			expectedErr: createError(statusInvalidProjectName, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name")),
		},
		{
			desc:  "err: ErrInvalidProjectUrlCode: can't use uppercase",
			setup: nil,
			req: &proto.CreateProjectRequest{
				Command: &proto.CreateProjectCommand{Name: "id-1", UrlCode: "CODE"},
			},
			expectedErr: createError(statusInvalidProjectUrlCode, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code")),
		},
		{
			desc:  "err: ErrInvalidProjectUrlCode: max id length exceeded",
			setup: nil,
			req: &proto.CreateProjectRequest{
				Command: &proto.CreateProjectCommand{Name: "id-1", UrlCode: strings.Repeat("a", 51)},
			},
			expectedErr: createError(statusInvalidProjectUrlCode, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code")),
		},
		{
			desc: "err: ErrProjectAlreadyExists: duplicate id",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrProjectAlreadyExists)
			},
			req: &proto.CreateProjectRequest{
				Command: &proto.CreateProjectCommand{OrganizationId: "organization-id", Name: "id-0", UrlCode: "id-0"},
			},
			expectedErr: createError(statusProjectAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
			},
			req: &proto.CreateProjectRequest{
				Command: &proto.CreateProjectCommand{OrganizationId: "organization-id", Name: "id-1", UrlCode: "id-1"},
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().CreateProject(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.CreateProjectRequest{
				Command: &proto.CreateProjectCommand{
					OrganizationId: "organization-id",
					Name:           projExpected.Name,
					UrlCode:        projExpected.UrlCode,
					Description:    projExpected.Description,
				},
			},
			expected:    projExpected.Project,
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.CreateProject(ctx, p.req)
			if resp != nil {
				assert.True(t, len(resp.Project.Name) > 0)
				assert.Equal(t, p.expected.Name, resp.Project.Name)
				assert.Equal(t, p.expected.UrlCode, resp.Project.UrlCode)
				assert.Equal(t, p.expected.Description, resp.Project.Description)
				assert.Equal(t, p.expected.CreatorEmail, resp.Project.CreatorEmail)
				assert.True(t, resp.Project.CreatedAt > 0)
				assert.True(t, resp.Project.UpdatedAt > 0)
				assert.Equal(t, p.expected.Disabled, resp.Project.Disabled)
				assert.Equal(t, p.expected.Trial, resp.Project.Trial)
			}
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCreateProjectNoCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
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

	expected, err := domain.NewProject(
		"project name",
		"project-url-code",
		"Description",
		"test@bucketer.io",
		"organization-id",
		false,
	)
	require.NoError(t, err)

	patterns := []struct {
		ctx         context.Context
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.CreateProjectRequest
		expected    *proto.Project
		expectedErr error
	}{
		{
			ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"accept-language": []string{"ja"},
			}),
			desc:  "err: unauthenticated",
			setup: nil,
			req:   &proto.CreateProjectRequest{},
			expectedErr: createError(
				statusUnauthenticated,
				localizer.MustLocalize(locale.UnauthenticatedError),
			),
		},
		{
			ctx: metadata.NewIncomingContext(createContextWithTokenRoleUnassigned(t), metadata.MD{
				"accept-language": []string{"ja"},
			}),
			desc: "err: permission denied",
			setup: func(s *EnvironmentService) {
				s.accountClient.(*acmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(&accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
					},
				}, nil)
			},
			req: &proto.CreateProjectRequest{},
			expectedErr: createError(
				statusPermissionDenied,
				localizer.MustLocalize(locale.PermissionDenied),
			),
		},
		{
			ctx:   ctx,
			desc:  "err: empty name",
			setup: nil,
			req: &proto.CreateProjectRequest{
				Name: "",
			},
			expectedErr: createError(
				statusEnvironmentNameRequired,
				localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
			),
		},
		{
			ctx:   ctx,
			desc:  "err: only space",
			setup: nil,
			req: &proto.CreateProjectRequest{
				Name: "    ",
			},
			expectedErr: createError(
				statusEnvironmentNameRequired,
				localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
			),
		},
		{
			ctx:   ctx,
			desc:  "err: max name length exceeded",
			setup: nil,
			req: &proto.CreateProjectRequest{
				Name: strings.Repeat("a", 51),
			},
			expectedErr: createError(
				statusInvalidEnvironmentName,
				localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name"),
			),
		},
		{
			ctx:   ctx,
			desc:  "err: empty url code",
			setup: nil,
			req: &proto.CreateProjectRequest{
				Name:    "name",
				UrlCode: "",
			},
			expectedErr: createError(
				statusInvalidEnvironmentUrlCode,
				localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code"),
			),
		},
		{
			ctx:   ctx,
			desc:  "err: url code can't use uppercase",
			setup: nil,
			req: &proto.CreateProjectRequest{
				Name:    "name",
				UrlCode: "URLCODE",
			},
			expectedErr: createError(
				statusInvalidEnvironmentUrlCode,
				localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code"),
			),
		},
		{
			ctx:   ctx,
			desc:  "err: url code can't use space",
			setup: nil,
			req: &proto.CreateProjectRequest{
				Name:    "name",
				UrlCode: "url code",
			},
			expectedErr: createError(
				statusInvalidEnvironmentUrlCode,
				localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code"),
			),
		},
		{
			ctx:   ctx,
			desc:  "err: max url code length exceeded",
			setup: nil,
			req: &proto.CreateProjectRequest{
				Name:    "name",
				UrlCode: strings.Repeat("a", 51),
			},
			expectedErr: createError(
				statusInvalidEnvironmentUrlCode,
				localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code"),
			),
		},
		{
			ctx:   ctx,
			desc:  "err: organization ID required",
			setup: nil,
			req: &proto.CreateProjectRequest{
				Name:           "name",
				UrlCode:        "url-code",
				OrganizationId: "",
			},
			expectedErr: createError(
				statusProjectIDRequired,
				localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
			),
		},
		{
			ctx:  ctx,
			desc: "err: project already exists",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrEnvironmentAlreadyExists)
			},
			req: &proto.CreateProjectRequest{
				Name:           expected.Project.Name,
				UrlCode:        expected.Project.UrlCode,
				OrganizationId: expected.Project.OrganizationId,
				Description:    expected.Project.Description,
			},
			expectedErr: createError(statusEnvironmentAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			ctx:  ctx,
			desc: "err: internal error",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal error"))
			},
			req: &proto.CreateProjectRequest{
				Name:           expected.Project.Name,
				UrlCode:        expected.Project.UrlCode,
				OrganizationId: expected.Project.OrganizationId,
				Description:    expected.Project.Description,
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal error")).Err(),
		},
		{
			ctx:  ctx,
			desc: "err: publish domain event failed",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				// Simulate a failure when publishing the update event.
				s.publisher.(*pubmock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "publish failed"))
			},
			req: &proto.CreateProjectRequest{
				Name:           expected.Project.Name,
				UrlCode:        expected.Project.UrlCode,
				OrganizationId: expected.Project.OrganizationId,
				Description:    expected.Project.Description,
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "publish failed")).Err(),
		},
		{
			ctx:  ctx,
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*pubmock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.CreateProjectRequest{
				Name:           expected.Project.Name,
				UrlCode:        expected.Project.UrlCode,
				OrganizationId: expected.Project.OrganizationId,
				Description:    expected.Project.Description,
			},
			expected: expected.Project,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.CreateProject(p.ctx, p.req)
			if resp != nil {
				assert.Equal(t, p.expected.Name, resp.Project.Name)
				assert.Equal(t, p.expected.UrlCode, resp.Project.UrlCode)
				assert.Equal(t, p.expected.Description, resp.Project.Description)
				assert.Equal(t, p.expected.OrganizationId, resp.Project.OrganizationId)
				assert.Equal(t, p.expected.Trial, resp.Project.Trial)
				assert.True(t, resp.Project.CreatedAt > 0)
				assert.True(t, resp.Project.UpdatedAt > 0)
			} else {
				assert.Equal(t, p.expectedErr, err)
			}
		})
	}
}

func TestCreateTrialProjectMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
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
		setup       func(*EnvironmentService)
		req         *proto.CreateTrialProjectRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			req: &proto.CreateTrialProjectRequest{
				Command: nil,
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrInvalidProjectName: empty name",
			setup: nil,
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Name: ""},
			},
			expectedErr: createError(statusProjectNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc:  "err: ErrInvalidProjectName: only space",
			setup: nil,
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Name: "   "},
			},
			expectedErr: createError(statusProjectNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc:  "err: ErrInvalidProjectName: max id length exceeded",
			setup: nil,
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Name: strings.Repeat("a", 51)},
			},
			expectedErr: createError(statusInvalidProjectName, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name")),
		},
		{
			desc:  "err: ErrInvalidProjectUrlCode: can't use uppercase",
			setup: nil,
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Name: "id-1", UrlCode: "CODE"},
			},
			expectedErr: createError(statusInvalidProjectUrlCode, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code")),
		},
		{
			desc:  "err: ErrInvalidProjectUrlCode: max id length exceeded",
			setup: nil,
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Name: "id-1", UrlCode: strings.Repeat("a", 51)},
			},
			expectedErr: createError(statusInvalidProjectUrlCode, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code")),
		},
		{
			desc:  "err: ErrInvalidProjectCreatorEmail",
			setup: nil,
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Name: "id-0", Email: "email"},
			},
			expectedErr: createError(statusInvalidProjectCreatorEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "owner_email")),
		},
		{
			desc: "err: ErrProjectAlreadyExists: trial exists",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetTrialProjectByEmail(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Project{
					Project: &proto.Project{Id: "id-0"},
				}, nil)
			},
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Name: "id-0", Email: "test@example.com"},
			},
			expectedErr: createError(statusProjectAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "err: ErrProjectAlreadyExists: duplicated id",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetTrialProjectByEmail(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Project{
					Project: &proto.Project{Id: "id-0"},
				}, nil)
			},
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Name: "id-0", Email: "test@example.com"},
			},
			expectedErr: createError(statusProjectAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetTrialProjectByEmail(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
			},
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Name: "id-1", Email: "test@example.com"},
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetTrialProjectByEmail(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Project{
					Project: nil,
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil).AnyTimes()
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().CreateProject(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				// CreateEnvironmentV2 is called for two purposes: development environment and production environment.
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil).Times(2)
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().CreateEnvironmentV2(
					gomock.Any(), gomock.Any(),
				).Return(nil).Times(2)

				s.accountClient.(*acmock.MockClient).EXPECT().CreateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(&accountproto.CreateAccountV2Response{}, nil)
			},
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Name: "Project Name_001", Email: "test@example.com"},
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
			_, err := service.CreateTrialProject(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateProjectNoCommand(t *testing.T) {
	// This test focuses on internal logic (database operations, domain events)
	// and assumes inputs are already validated by the API layer (UpdateProject).
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	// Set up a context with token and language metadata.
	ctx := createContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)

	// Helper to create errors with localized messages
	createError := func(st *status.Status, msg string) error {
		stWithDetails, err := st.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return stWithDetails.Err()
	}

	editor := &eventproto.Editor{
		Email: "test@bucketer.io",
	}

	// Create a test project for the internal method tests
	testProject := &domain.Project{
		Project: &environmentproto.Project{
			Id:             "project-id",
			Name:           "Test Project",
			Description:    "Test Description",
			OrganizationId: "org-1",
		},
	}

	// Define test patterns.
	patterns := []struct {
		ctx         context.Context
		desc        string
		setup       func(*EnvironmentService)
		req         *environmentproto.UpdateProjectRequest
		expected    *environmentproto.UpdateProjectResponse
		expectedErr error
	}{
		{
			ctx:  ctx,
			desc: "err: project not found",
			setup: func(s *EnvironmentService) {
				// Simulate that the transaction returns an error from GetProject.
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrProjectNotFound)
			},
			req: &environmentproto.UpdateProjectRequest{
				Id:          "nonexistent",
				Name:        &wrappers.StringValue{Value: "ValidName"},
				Description: &wrappers.StringValue{Value: "updated description"},
			},
			expectedErr: createError(
				statusProjectNotFound,
				localizer.MustLocalize(locale.NotFoundError),
			),
		},
		{
			ctx:  ctx,
			desc: "err: update failed",
			setup: func(s *EnvironmentService) {
				// Simulate an error during the transaction (e.g. when calling proj.Update).
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(context.Context, mysql.Transaction) error) error {
					return pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "update failed")
				})
			},
			req: &environmentproto.UpdateProjectRequest{
				Id:          "project-id",
				Name:        &wrappers.StringValue{Value: "ValidName"},
				Description: &wrappers.StringValue{Value: "updated description"},
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "update failed")).Err(),
		},
		{
			ctx:  ctx,
			desc: "err: publish domain event failed",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				// Simulate a failure when publishing the update event.
				s.publisher.(*pubmock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "publish failed"))
			},
			req: &environmentproto.UpdateProjectRequest{
				Id:          "project-id",
				Name:        &wrappers.StringValue{Value: "ValidName"},
				Description: &wrappers.StringValue{Value: "updated description"},
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "publish failed")).Err(),
		},
		{
			ctx:  ctx,
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*pubmock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &environmentproto.UpdateProjectRequest{
				Id:          "project-id",
				Name:        &wrappers.StringValue{Value: "ValidName"},
				Description: &wrappers.StringValue{Value: "updated description"},
			},
			expected: &environmentproto.UpdateProjectResponse{},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.updateProjectNoCommand(p.ctx, p.req, testProject, localizer, editor)
			if resp != nil {
				// For a successful update, compare the expected response.
				assert.Equal(t, p.expected, resp)
			} else {
				assert.Equal(t, p.expectedErr, err)
			}
		})
	}
}

func TestUpdateProjectMySQL(t *testing.T) {
	// This test covers API-level validation and security checks for UpdateProject.
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
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
		setup       func(*EnvironmentService)
		req         *proto.UpdateProjectRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrProjectIDRequired",
			setup: nil,
			req: &proto.UpdateProjectRequest{
				ChangeDescriptionCommand: &proto.ChangeDescriptionProjectCommand{Description: "desc"},
			},
			expectedErr: createError(statusProjectIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "err: ErrInvalidProjectName",
			setup: nil,
			req: &proto.UpdateProjectRequest{
				Id:            "id-0",
				RenameCommand: &proto.RenameProjectCommand{Name: strings.Repeat("a", 51)},
			},
			expectedErr: createError(statusInvalidProjectName, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name")),
		},
		{
			desc:  "err: empty name (no-command path)",
			setup: nil,
			req: &proto.UpdateProjectRequest{
				Id:          "project-id",
				Name:        &wrappers.StringValue{Value: "    "},
				Description: &wrappers.StringValue{Value: "updated description"},
			},
			expectedErr: createError(statusEnvironmentNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc:  "err: max name length exceeded (no-command path)",
			setup: nil,
			req: &proto.UpdateProjectRequest{
				Id:          "project-id",
				Name:        &wrappers.StringValue{Value: strings.Repeat("a", 51)},
				Description: &wrappers.StringValue{Value: "updated description"},
			},
			expectedErr: createError(statusInvalidEnvironmentName, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name")),
		},
		{
			desc: "err: ErrProjectNotFound",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "id-0",
				).Return(nil, v2es.ErrProjectNotFound)
			},
			req: &proto.UpdateProjectRequest{
				Id:                       "id-0",
				ChangeDescriptionCommand: &proto.ChangeDescriptionProjectCommand{Description: "desc"},
			},
			expectedErr: createError(statusProjectNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "id-1",
				).Return(&domain.Project{
					Project: &proto.Project{Id: "id-1", OrganizationId: "org-1"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
			},
			req: &proto.UpdateProjectRequest{
				Id:                       "id-1",
				ChangeDescriptionCommand: &proto.ChangeDescriptionProjectCommand{Description: "desc"},
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "id-1",
				).Return(&domain.Project{
					Project: &proto.Project{Id: "id-1", OrganizationId: "org-1", Description: "old desc"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().UpdateProject(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.UpdateProjectRequest{
				Id:                       "id-1",
				ChangeDescriptionCommand: &proto.ChangeDescriptionProjectCommand{Description: "desc"},
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
			_, err := service.UpdateProject(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestEnableProjectMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
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
		setup       func(*EnvironmentService)
		req         *proto.EnableProjectRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			req: &proto.EnableProjectRequest{
				Id: "id-0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrProjectIDRequired",
			setup: nil,
			req: &proto.EnableProjectRequest{
				Command: &proto.EnableProjectCommand{},
			},
			expectedErr: createError(statusProjectIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrProjectNotFound",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "id-0",
				).Return(nil, v2es.ErrProjectNotFound)
			},
			req: &proto.EnableProjectRequest{
				Id:      "id-0",
				Command: &proto.EnableProjectCommand{},
			},
			expectedErr: createError(statusProjectNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "id-1",
				).Return(&domain.Project{
					Project: &proto.Project{Id: "id-1", OrganizationId: "org-1"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
			},
			req: &proto.EnableProjectRequest{
				Id:      "id-1",
				Command: &proto.EnableProjectCommand{},
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "id-1",
				).Return(&domain.Project{
					Project: &proto.Project{Id: "id-1", OrganizationId: "org-1"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().UpdateProject(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.EnableProjectRequest{
				Id:      "id-1",
				Command: &proto.EnableProjectCommand{},
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
			_, err := service.EnableProject(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDisableProjectMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
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
		setup       func(*EnvironmentService)
		req         *proto.DisableProjectRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			req: &proto.DisableProjectRequest{
				Id: "id-0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrProjectIDRequired",
			setup: nil,
			req: &proto.DisableProjectRequest{
				Command: &proto.DisableProjectCommand{},
			},
			expectedErr: createError(statusProjectIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrProjectNotFound",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "id-0",
				).Return(nil, v2es.ErrProjectNotFound)
			},
			req: &proto.DisableProjectRequest{
				Id:      "id-0",
				Command: &proto.DisableProjectCommand{},
			},
			expectedErr: createError(statusProjectNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "id-1",
				).Return(&domain.Project{
					Project: &proto.Project{Id: "id-1", OrganizationId: "org-1"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
			},
			req: &proto.DisableProjectRequest{
				Id:      "id-1",
				Command: &proto.DisableProjectCommand{},
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "id-1",
				).Return(&domain.Project{
					Project: &proto.Project{Id: "id-1", OrganizationId: "org-1"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().UpdateProject(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.DisableProjectRequest{
				Id:      "id-1",
				Command: &proto.DisableProjectCommand{},
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
			_, err := service.DisableProject(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestConvertTrialProjectMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
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
		setup       func(*EnvironmentService)
		req         *proto.ConvertTrialProjectRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			req: &proto.ConvertTrialProjectRequest{
				Id: "id-0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrProjectIDRequired",
			setup: nil,
			req: &proto.ConvertTrialProjectRequest{
				Command: &proto.ConvertTrialProjectCommand{},
			},
			expectedErr: createError(statusProjectIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrProjectNotFound",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "id-0",
				).Return(nil, v2es.ErrProjectNotFound)
			},
			req: &proto.ConvertTrialProjectRequest{
				Id:      "id-0",
				Command: &proto.ConvertTrialProjectCommand{},
			},
			expectedErr: createError(statusProjectNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "id-1",
				).Return(&domain.Project{
					Project: &proto.Project{Id: "id-1", OrganizationId: "org-1"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
			},
			req: &proto.ConvertTrialProjectRequest{
				Id:      "id-1",
				Command: &proto.ConvertTrialProjectCommand{},
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "id-1",
				).Return(&domain.Project{
					Project: &proto.Project{Id: "id-1", OrganizationId: "org-1"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().UpdateProject(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.ConvertTrialProjectRequest{
				Id:      "id-1",
				Command: &proto.ConvertTrialProjectCommand{},
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
			_, err := service.ConvertTrialProject(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestProjectPermissionDeniedMySQL(t *testing.T) {
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
		desc     string
		setup    func(*EnvironmentService)
		action   func(context.Context, *EnvironmentService) error
		expected error
	}{
		{
			desc: "CreateTrialProject",
			setup: func(s *EnvironmentService) {
				// No setup needed - this fails at system admin check
			},
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.CreateTrialProject(ctx, &proto.CreateTrialProjectRequest{})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "EnableProject",
			setup: func(s *EnvironmentService) {
				// Mock project fetch to succeed so we can reach the permission check
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "project-id",
				).Return(&domain.Project{
					Project: &proto.Project{Id: "project-id", OrganizationId: "org-1"},
				}, nil)
				// Mock account client call to return permission denied
				s.accountClient.(*acmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), &accountproto.GetAccountV2Request{
						Email:          "email",
						OrganizationId: "org-1",
					},
				).Return(&accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						OrganizationRole: accountproto.AccountV2_Role_Organization_UNASSIGNED,
					},
				}, nil)
			},
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.EnableProject(ctx, &proto.EnableProjectRequest{
					Id:      "project-id",
					Command: &proto.EnableProjectCommand{},
				})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "DisableProject",
			setup: func(s *EnvironmentService) {
				// Mock project fetch to succeed so we can reach the permission check
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "project-id",
				).Return(&domain.Project{
					Project: &proto.Project{Id: "project-id", OrganizationId: "org-1"},
				}, nil)
				// Mock account client call to return permission denied
				s.accountClient.(*acmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), &accountproto.GetAccountV2Request{
						Email:          "email",
						OrganizationId: "org-1",
					},
				).Return(&accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						OrganizationRole: accountproto.AccountV2_Role_Organization_UNASSIGNED,
					},
				}, nil)
			},
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.DisableProject(ctx, &proto.DisableProjectRequest{
					Id:      "project-id",
					Command: &proto.DisableProjectCommand{},
				})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "UpdateProject",
			setup: func(s *EnvironmentService) {
				// Mock project fetch to succeed so we can reach the permission check
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "project-id",
				).Return(&domain.Project{
					Project: &proto.Project{Id: "project-id", OrganizationId: "org-1"},
				}, nil)
				// Mock account client call to return permission denied
				s.accountClient.(*acmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), &accountproto.GetAccountV2Request{
						Email:          "email",
						OrganizationId: "org-1",
					},
				).Return(&accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						OrganizationRole: accountproto.AccountV2_Role_Organization_UNASSIGNED,
					},
				}, nil)
			},
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.UpdateProject(ctx, &proto.UpdateProjectRequest{
					Id:            "project-id",
					RenameCommand: &proto.RenameProjectCommand{Name: "New Name"},
				})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "UpdateProjectNoCommand",
			setup: func(s *EnvironmentService) {
				// Mock project fetch to succeed so we can reach the permission check
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "project-id",
				).Return(&domain.Project{
					Project: &proto.Project{Id: "project-id", OrganizationId: "org-1"},
				}, nil)
				// Mock account client call to return permission denied
				s.accountClient.(*acmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), &accountproto.GetAccountV2Request{
						Email:          "email",
						OrganizationId: "org-1",
					},
				).Return(&accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						OrganizationRole: accountproto.AccountV2_Role_Organization_UNASSIGNED,
					},
				}, nil)
			},
			action: func(ctx context.Context, es *EnvironmentService) error {
				name := "Updated Name"
				description := "Updated Description"
				_, err := es.UpdateProject(ctx, &proto.UpdateProjectRequest{
					Id:          "project-id",
					Name:        &wrapperspb.StringValue{Value: name},
					Description: &wrapperspb.StringValue{Value: description},
				})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "ConvertTrialProject",
			setup: func(s *EnvironmentService) {
				// Mock project fetch to succeed so we can reach the permission check
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().GetProject(
					gomock.Any(), "project-id",
				).Return(&domain.Project{
					Project: &proto.Project{Id: "project-id", OrganizationId: "org-1"},
				}, nil)
				// Mock account client call to return permission denied
				s.accountClient.(*acmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), &accountproto.GetAccountV2Request{
						Email:          "email",
						OrganizationId: "org-1",
					},
				).Return(&accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						OrganizationRole: accountproto.AccountV2_Role_Organization_UNASSIGNED,
					},
				}, nil)
			},
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.ConvertTrialProject(ctx, &proto.ConvertTrialProjectRequest{
					Id:      "project-id",
					Command: &proto.ConvertTrialProjectCommand{},
				})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			actual := p.action(ctx, service)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestListProjectsV2(t *testing.T) {
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
		setup       func(*EnvironmentService)
		input       *environmentproto.ListProjectsV2Request
		expected    *environmentproto.ListProjectsV2Response
		expectedErr error
	}{
		{
			desc: "success: list projects",
			setup: func(s *EnvironmentService) {
				s.accountClient.(*acmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(&accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
					},
				}, nil)
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().ListProjects(
					gomock.Any(), gomock.Any(),
				).Return([]*environmentproto.Project{{}}, 1, int64(1), nil)
			},
			input: &environmentproto.ListProjectsV2Request{
				PageSize:       10,
				Cursor:         "",
				OrganizationId: "org-1",
			},
			expected: &environmentproto.ListProjectsV2Response{
				Projects:   []*environmentproto.Project{{}}, // Expect one project
				Cursor:     "1",
				TotalCount: 1,
			},
			expectedErr: nil,
		},
		{
			desc: "failure: permission denied",
			setup: func(s *EnvironmentService) {
				s.accountClient.(*acmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), &accountproto.GetAccountV2Request{
						Email:          "email",
						OrganizationId: "org-1",
					},
				).Return(&accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						OrganizationRole: accountproto.AccountV2_Role_Organization_UNASSIGNED,
					},
				}, nil)
			},
			input: &environmentproto.ListProjectsV2Request{
				PageSize:       10,
				Cursor:         "",
				OrganizationId: "org-1",
			},
			expected:    nil,
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "failure: invalid cursor",
			setup: func(s *EnvironmentService) {
				s.accountClient.(*acmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(&accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
					},
				}, nil)
			},
			input: &environmentproto.ListProjectsV2Request{
				PageSize:       10,
				Cursor:         "invalid",
				OrganizationId: "org-1",
			},
			expected:    nil,
			expectedErr: createError(statusInvalidCursor, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor")),
		},
		{
			desc: "failure: internal error",
			setup: func(s *EnvironmentService) {
				s.accountClient.(*acmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(&accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
					},
				}, nil)
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().ListProjects(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
			},
			input: &environmentproto.ListProjectsV2Request{
				PageSize:       10,
				Cursor:         "",
				OrganizationId: "org-1",
			},
			expected:    nil,
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal")).Err(),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			actual, err := service.ListProjectsV2(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}
