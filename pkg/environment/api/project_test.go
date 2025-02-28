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

	acmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	pubmock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
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
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "err-id-0",
			expectedErr: createError(statusProjectNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "err-id-1",
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       &proto.ListProjectsRequest{},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.CreateProjectRequest{
				Command: &proto.CreateProjectCommand{OrganizationId: "organization-id", Name: "id-1", UrlCode: "id-1"},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
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
				).Return(errors.New("internal error"))
			},
			req: &proto.CreateProjectRequest{
				Name:           expected.Project.Name,
				UrlCode:        expected.Project.UrlCode,
				OrganizationId: expected.Project.OrganizationId,
				Description:    expected.Project.Description,
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
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
				).Return(errors.New("publish failed"))
			},
			req: &proto.CreateProjectRequest{
				Name:           expected.Project.Name,
				UrlCode:        expected.Project.UrlCode,
				OrganizationId: expected.Project.OrganizationId,
				Description:    expected.Project.Description,
			},
			expectedErr: createError(
				statusInternal,
				localizer.MustLocalize(locale.InternalServerError),
			),
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
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Name: "id-0", Email: "test@example.com"},
			},
			expectedErr: createError(statusProjectAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "err: ErrProjectAlreadyExists: duplicated id",
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrProjectAlreadyExists)
			},
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Name: "id-0", Email: "test@example.com"},
			},
			expectedErr: createError(statusProjectAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req: &proto.CreateTrialProjectRequest{
				Command: &proto.CreateTrialProjectCommand{Name: "id-1", Email: "test@example.com"},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil).Times(4)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil).Times(4)
				s.accountClient.(*acmock.MockClient).EXPECT().CreateAccountV2(gomock.Any(), gomock.Any()).Return(
					&accountproto.CreateAccountV2Response{}, nil)
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
			ctx:   ctx,
			desc:  "err: empty name",
			setup: nil,
			req: &environmentproto.UpdateProjectRequest{
				Id:          "project-id",
				Name:        &wrappers.StringValue{Value: "    "},
				Description: &wrappers.StringValue{Value: "updated description"},
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
			req: &environmentproto.UpdateProjectRequest{
				Id:          "project-id",
				Name:        &wrappers.StringValue{Value: strings.Repeat("a", 51)},
				Description: &wrappers.StringValue{Value: "updated description"},
			},
			expectedErr: createError(
				statusInvalidEnvironmentName,
				localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name"),
			),
		},
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
					return errors.New("update failed")
				})
			},
			req: &environmentproto.UpdateProjectRequest{
				Id:          "project-id",
				Name:        &wrappers.StringValue{Value: "ValidName"},
				Description: &wrappers.StringValue{Value: "updated description"},
			},
			expectedErr: createError(
				statusInternal,
				localizer.MustLocalize(locale.InternalServerError),
			),
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
				).Return(errors.New("publish failed"))
			},
			req: &environmentproto.UpdateProjectRequest{
				Id:          "project-id",
				Name:        &wrappers.StringValue{Value: "ValidName"},
				Description: &wrappers.StringValue{Value: "updated description"},
			},
			expectedErr: createError(
				statusInternal,
				localizer.MustLocalize(locale.InternalServerError),
			),
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
			resp, err := service.updateProjectNoCommand(p.ctx, p.req, localizer, editor)
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
			desc: "err: ErrProjectNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrProjectNotFound)
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.UpdateProjectRequest{
				Id:                       "id-1",
				ChangeDescriptionCommand: &proto.ChangeDescriptionProjectCommand{Description: "desc"},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrProjectNotFound)
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.EnableProjectRequest{
				Id:      "id-1",
				Command: &proto.EnableProjectCommand{},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrProjectNotFound)
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.DisableProjectRequest{
				Id:      "id-1",
				Command: &proto.DisableProjectCommand{},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrProjectNotFound)
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.ConvertTrialProjectRequest{
				Id:      "id-1",
				Command: &proto.ConvertTrialProjectCommand{},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
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
		action   func(context.Context, *EnvironmentService) error
		expected error
	}{
		{
			desc: "CreateTrialProject",
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.CreateTrialProject(ctx, &proto.CreateTrialProjectRequest{})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "EnableProject",
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.EnableProject(ctx, &proto.EnableProjectRequest{})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "DisableProject",
			action: func(ctx context.Context, es *EnvironmentService) error {
				_, err := es.DisableProject(ctx, &proto.DisableProjectRequest{})
				return err
			},
			expected: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newEnvironmentService(t, mockController, nil)
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
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(true)
				rows.EXPECT().Scan(gomock.Any()).Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
					// Mock the TotalCount
					*dest[0].(*int64) = 1
					return nil
				})
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("internal error"))
			},
			input: &environmentproto.ListProjectsV2Request{
				PageSize:       10,
				Cursor:         "",
				OrganizationId: "org-1",
			},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
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
