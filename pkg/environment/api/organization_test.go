package api

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestGetOrganizationMySQL(t *testing.T) {
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
			desc:        "err: ErrOrganizationIDRequired",
			setup:       nil,
			id:          "",
			expectedErr: createError(statusOrganizationIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrOrganizationNotFound",
			setup: func(s *EnvironmentService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "err-id-0",
			expectedErr: createError(statusOrganizationNotFound, localizer.MustLocalize(locale.NotFoundError)),
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
			req := &proto.GetOrganizationRequest{Id: p.id}
			resp, err := s.GetOrganization(ctx, req)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestListOrganizationsMySQL(t *testing.T) {
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
		input       *proto.ListOrganizationsRequest
		expected    *proto.ListOrganizationsResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrInvalidCursor",
			setup:       nil,
			input:       &proto.ListOrganizationsRequest{Cursor: "XXX"},
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
			input:       &proto.ListOrganizationsRequest{},
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
			input:       &proto.ListOrganizationsRequest{PageSize: 2, Cursor: ""},
			expected:    &proto.ListOrganizationsResponse{Organizations: []*proto.Organization{}, Cursor: "0", TotalCount: 0},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.ListOrganizations(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestCreateOrganizationMySQL(t *testing.T) {
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

	orgExpected, err := domain.NewOrganization(
		"name",
		"url-code",
		"description",
		false,
		false,
	)
	require.NoError(t, err)
	trialOrgExpected, err := domain.NewOrganization(
		"name2",
		"url-code2",
		"description2",
		true,
		false,
	)
	require.NoError(t, err)

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.CreateOrganizationRequest
		expected    *proto.Organization
		expectedErr error
	}{
		{
			desc:        "err: ErrNoCommand",
			setup:       nil,
			req:         &proto.CreateOrganizationRequest{},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrInvalidOrganizationName: empty name",
			setup: nil,
			req: &proto.CreateOrganizationRequest{
				Command: &proto.CreateOrganizationCommand{Name: ""},
			},
			expectedErr: createError(statusOrganizationNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc:  "err: ErrInvalidOrganizationName: only space",
			setup: nil,
			req: &proto.CreateOrganizationRequest{
				Command: &proto.CreateOrganizationCommand{Name: "    "},
			},
			expectedErr: createError(statusOrganizationNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc:  "err: ErrInvalidOrganizationName: max name length exceeded",
			setup: nil,
			req: &proto.CreateOrganizationRequest{
				Command: &proto.CreateOrganizationCommand{Name: strings.Repeat("a", 51)},
			},
			expectedErr: createError(statusInvalidOrganizationName, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name")),
		},
		{
			desc:  "err: ErrInvalidOrganizationUrlCode: can't use uppercase",
			setup: nil,
			req: &proto.CreateOrganizationRequest{
				Command: &proto.CreateOrganizationCommand{Name: "id-1", UrlCode: "CODE"},
			},
			expectedErr: createError(statusInvalidOrganizationUrlCode, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code")),
		},
		{
			desc:  "err: ErrInvalidOrganizationUrlCode: max id length exceeded",
			setup: nil,
			req: &proto.CreateOrganizationRequest{
				Command: &proto.CreateOrganizationCommand{Name: "id-1", UrlCode: strings.Repeat("a", 51)},
			},
			expectedErr: createError(statusInvalidOrganizationUrlCode, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code")),
		},
		{
			desc: "err: ErrOrganizationAlreadyExists: duplicate id",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrOrganizationAlreadyExists)
			},
			req: &proto.CreateOrganizationRequest{
				Command: &proto.CreateOrganizationCommand{Name: "id-0", UrlCode: "id-0"},
			},
			expectedErr: createError(statusOrganizationAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.CreateOrganizationRequest{
				Command: &proto.CreateOrganizationCommand{Name: "id-1", UrlCode: "id-1"},
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
			req: &proto.CreateOrganizationRequest{
				Command: &proto.CreateOrganizationCommand{
					Name:        orgExpected.Name,
					UrlCode:     orgExpected.UrlCode,
					Description: orgExpected.Description,
					IsTrial:     false,
				},
			},
			expected:    orgExpected.Organization,
			expectedErr: nil,
		},
		{
			desc: "success trial",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.CreateOrganizationRequest{
				Command: &proto.CreateOrganizationCommand{
					Name:        trialOrgExpected.Name,
					UrlCode:     trialOrgExpected.UrlCode,
					Description: trialOrgExpected.Description,
					IsTrial:     trialOrgExpected.Trial,
				},
			},
			expected:    trialOrgExpected.Organization,
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.CreateOrganization(ctx, p.req)
			if resp != nil {
				assert.True(t, len(resp.Organization.Name) > 0)
				assert.Equal(t, p.expected.Name, resp.Organization.Name)
				assert.Equal(t, p.expected.UrlCode, resp.Organization.UrlCode)
				assert.Equal(t, p.expected.Description, resp.Organization.Description)
				assert.Equal(t, p.expected.Disabled, resp.Organization.Disabled)
				assert.Equal(t, p.expected.Archived, resp.Organization.Archived)
				assert.Equal(t, p.expected.Trial, resp.Organization.Trial)
				assert.True(t, resp.Organization.CreatedAt > 0)
				assert.True(t, resp.Organization.UpdatedAt > 0)
			}
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateOrganizationMySQL(t *testing.T) {
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
		req         *proto.UpdateOrganizationRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrNoCommand",
			setup:       nil,
			req:         &proto.UpdateOrganizationRequest{Id: "id-0"},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrInvalidOrganizationName: empty name",
			setup: nil,
			req: &proto.UpdateOrganizationRequest{
				Id:            "id-0",
				RenameCommand: &proto.ChangeNameOrganizationCommand{Name: ""},
			},
			expectedErr: createError(statusOrganizationNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc:  "err: ErrInvalidOrganizationName: only space",
			setup: nil,
			req: &proto.UpdateOrganizationRequest{
				Id:            "id-0",
				RenameCommand: &proto.ChangeNameOrganizationCommand{Name: "    "},
			},
			expectedErr: createError(statusOrganizationNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc:  "err: ErrInvalidOrganizationName: max name length exceeded",
			setup: nil,
			req: &proto.UpdateOrganizationRequest{
				Id:            "id-0",
				RenameCommand: &proto.ChangeNameOrganizationCommand{Name: strings.Repeat("a", 51)},
			},
			expectedErr: createError(statusInvalidOrganizationName, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name")),
		},
		{
			desc: "err: ErrOrganizationNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrOrganizationNotFound)
			},
			req: &proto.UpdateOrganizationRequest{
				Id:            "err-id-0",
				RenameCommand: &proto.ChangeNameOrganizationCommand{Name: "id-0"},
			},
			expectedErr: createError(statusOrganizationNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.UpdateOrganizationRequest{
				Id:            "err-id-1",
				RenameCommand: &proto.ChangeNameOrganizationCommand{Name: "id-1"},
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
			req: &proto.UpdateOrganizationRequest{
				Id:            "success-id-0",
				RenameCommand: &proto.ChangeNameOrganizationCommand{Name: "id-0"},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.UpdateOrganization(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestEnableOrganizationMySQL(t *testing.T) {
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
		req         *proto.EnableOrganizationRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			req: &proto.EnableOrganizationRequest{
				Id: "id-0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrOrganizationIDRequired",
			setup: nil,
			req: &proto.EnableOrganizationRequest{
				Command: &proto.EnableOrganizationCommand{},
			},
			expectedErr: createError(statusOrganizationIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrOrganizationNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrOrganizationNotFound)
			},
			req: &proto.EnableOrganizationRequest{
				Id:      "id-0",
				Command: &proto.EnableOrganizationCommand{},
			},
			expectedErr: createError(statusOrganizationNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.EnableOrganizationRequest{
				Id:      "id-1",
				Command: &proto.EnableOrganizationCommand{},
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
			req: &proto.EnableOrganizationRequest{
				Id:      "id-1",
				Command: &proto.EnableOrganizationCommand{},
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
			_, err := service.EnableOrganization(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDisableOrganizationMySQL(t *testing.T) {
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
		req         *proto.DisableOrganizationRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			req: &proto.DisableOrganizationRequest{
				Id: "id-0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrOrganizationIDRequired",
			setup: nil,
			req: &proto.DisableOrganizationRequest{
				Command: &proto.DisableOrganizationCommand{},
			},
			expectedErr: createError(statusOrganizationIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrOrganizationNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrOrganizationNotFound)
			},
			req: &proto.DisableOrganizationRequest{
				Id:      "id-0",
				Command: &proto.DisableOrganizationCommand{},
			},
			expectedErr: createError(statusOrganizationNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrCannotUpdateSystemAdmin",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(domain.ErrCannotDisableSystemAdmin)
			},
			req: &proto.DisableOrganizationRequest{
				Id:      "id-0",
				Command: &proto.DisableOrganizationCommand{},
			},
			expectedErr: createError(statusCannotUpdateSystemAdmin, localizer.MustLocalize(locale.InvalidArgumentError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.DisableOrganizationRequest{
				Id:      "id-1",
				Command: &proto.DisableOrganizationCommand{},
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
			req: &proto.DisableOrganizationRequest{
				Id:      "id-1",
				Command: &proto.DisableOrganizationCommand{},
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
			_, err := service.DisableOrganization(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestArchiveOrganizationMySQL(t *testing.T) {
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
		req         *proto.ArchiveOrganizationRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			req: &proto.ArchiveOrganizationRequest{
				Id: "id-0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrOrganizationIDRequired",
			setup: nil,
			req: &proto.ArchiveOrganizationRequest{
				Command: &proto.ArchiveOrganizationCommand{},
			},
			expectedErr: createError(statusOrganizationIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrOrganizationNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrOrganizationNotFound)
			},
			req: &proto.ArchiveOrganizationRequest{
				Id:      "id-0",
				Command: &proto.ArchiveOrganizationCommand{},
			},
			expectedErr: createError(statusOrganizationNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrCannotUpdateSystemAdmin",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(domain.ErrCannotArchiveSystemAdmin)
			},
			req: &proto.ArchiveOrganizationRequest{
				Id:      "id-0",
				Command: &proto.ArchiveOrganizationCommand{},
			},
			expectedErr: createError(statusCannotUpdateSystemAdmin, localizer.MustLocalize(locale.InvalidArgumentError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.ArchiveOrganizationRequest{
				Id:      "id-1",
				Command: &proto.ArchiveOrganizationCommand{},
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
			req: &proto.ArchiveOrganizationRequest{
				Id:      "id-1",
				Command: &proto.ArchiveOrganizationCommand{},
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
			_, err := service.ArchiveOrganization(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUnarchiveOrganizationMySQL(t *testing.T) {
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
		req         *proto.UnarchiveOrganizationRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			req: &proto.UnarchiveOrganizationRequest{
				Id: "id-0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrOrganizationIDRequired",
			setup: nil,
			req: &proto.UnarchiveOrganizationRequest{
				Command: &proto.UnarchiveOrganizationCommand{},
			},
			expectedErr: createError(statusOrganizationIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrOrganizationNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrOrganizationNotFound)
			},
			req: &proto.UnarchiveOrganizationRequest{
				Id:      "id-0",
				Command: &proto.UnarchiveOrganizationCommand{},
			},
			expectedErr: createError(statusOrganizationNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.UnarchiveOrganizationRequest{
				Id:      "id-1",
				Command: &proto.UnarchiveOrganizationCommand{},
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
			req: &proto.UnarchiveOrganizationRequest{
				Id:      "id-1",
				Command: &proto.UnarchiveOrganizationCommand{},
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
			_, err := service.UnarchiveOrganization(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestConvertTrialOrganizationMySQL(t *testing.T) {
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
		req         *proto.ConvertTrialOrganizationRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			req: &proto.ConvertTrialOrganizationRequest{
				Id: "id-0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrOrganizationIDRequired",
			setup: nil,
			req: &proto.ConvertTrialOrganizationRequest{
				Command: &proto.ConvertTrialOrganizationCommand{},
			},
			expectedErr: createError(statusOrganizationIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrOrganizationNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2es.ErrOrganizationNotFound)
			},
			req: &proto.ConvertTrialOrganizationRequest{
				Id:      "id-0",
				Command: &proto.ConvertTrialOrganizationCommand{},
			},
			expectedErr: createError(statusOrganizationNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.ConvertTrialOrganizationRequest{
				Id:      "id-1",
				Command: &proto.ConvertTrialOrganizationCommand{},
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
			req: &proto.ConvertTrialOrganizationRequest{
				Id:      "id-1",
				Command: &proto.ConvertTrialOrganizationCommand{},
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
			_, err := service.ConvertTrialOrganization(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
