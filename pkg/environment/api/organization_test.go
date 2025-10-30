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
	accountdomain "github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	accstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2/mock"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

// Helper function to create default authentication settings for tests
func newTestAuthSettings() *proto.AuthenticationSettings {
	return &proto.AuthenticationSettings{
		EnabledTypes: []proto.AuthenticationType{
			proto.AuthenticationType_AUTHENTICATION_TYPE_GOOGLE,
		},
	}
}

func TestGetOrganizationMySQL(t *testing.T) {
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
			desc:        "err: ErrOrganizationIDRequired",
			setup:       nil,
			id:          "",
			expectedErr: statusOrganizationIDRequired.Err(),
		},
		{
			desc: "err: ErrOrganizationNotFound",
			setup: func(s *EnvironmentService) {
				s.orgStorage.(*storagemock.MockOrganizationStorage).EXPECT().GetOrganization(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2es.ErrOrganizationNotFound)
			},
			id:          "err-id-0",
			expectedErr: statusOrganizationNotFound.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.orgStorage.(*storagemock.MockOrganizationStorage).EXPECT().GetOrganization(
					gomock.Any(), gomock.Any(),
				).Return(nil, pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
			},
			id:          "err-id-1",
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.orgStorage.(*storagemock.MockOrganizationStorage).EXPECT().GetOrganization(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Organization{
					Organization: &proto.Organization{Id: "success-id-0"},
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
			expectedErr: statusInvalidCursor.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.orgStorage.(*storagemock.MockOrganizationStorage).EXPECT().ListOrganizations(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
			},
			input:       &proto.ListOrganizationsRequest{},
			expected:    nil,
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.orgStorage.(*storagemock.MockOrganizationStorage).EXPECT().ListOrganizations(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Organization{}, 0, int64(0), nil)
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

	orgExpected, err := domain.NewOrganization(
		"name",
		"url-code",
		"test@example.com",
		"description",
		false,
		false,
		newTestAuthSettings(),
	)
	require.NoError(t, err)
	trialOrgExpected, err := domain.NewOrganization(
		"name2",
		"url-code2",
		"test@test.org",
		"description2",
		true,
		false,
		newTestAuthSettings(),
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
			desc:  "err: ErrInvalidOrganizationName: empty name",
			setup: nil,
			req: &proto.CreateOrganizationRequest{
				Name: "",
			},
			expectedErr: statusOrganizationNameRequired.Err(),
		},
		{
			desc:  "err: ErrInvalidOrganizationName: only space",
			setup: nil,
			req: &proto.CreateOrganizationRequest{
				Name: "    ",
			},
			expectedErr: statusOrganizationNameRequired.Err(),
		},
		{
			desc:  "err: ErrInvalidOrganizationName: max name length exceeded",
			setup: nil,
			req: &proto.CreateOrganizationRequest{
				Name: strings.Repeat("a", 51),
			},
			expectedErr: statusInvalidOrganizationName.Err(),
		},
		{
			desc:  "err: ErrInvalidOrganizationUrlCode: can't use uppercase",
			setup: nil,
			req: &proto.CreateOrganizationRequest{
				Name:    "id-1",
				UrlCode: "CODE",
			},
			expectedErr: statusInvalidOrganizationUrlCode.Err(),
		},
		{
			desc:  "err: ErrInvalidOrganizationUrlCode: max id length exceeded",
			setup: nil,
			req: &proto.CreateOrganizationRequest{
				Name:    "id-1",
				UrlCode: strings.Repeat("a", 51),
			},
			expectedErr: statusInvalidOrganizationUrlCode.Err(),
		},
		{
			desc: "err: ErrOrganizationAlreadyExists: duplicate id",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrOrganizationAlreadyExists)
			},
			req: &proto.CreateOrganizationRequest{
				Name:       "id-0",
				UrlCode:    "id-0",
				OwnerEmail: "test@test.org",
			},
			expectedErr: statusOrganizationAlreadyExists.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error"))
			},
			req: &proto.CreateOrganizationRequest{
				Name:       "id-1",
				UrlCode:    "id-1",
				OwnerEmail: "test@test.org",
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.orgStorage.(*storagemock.MockOrganizationStorage).EXPECT().CreateOrganization(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().CreateProject(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().CreateEnvironmentV2(
					gomock.Any(), gomock.Any(),
				).Times(2).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().CreateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &proto.CreateOrganizationRequest{
				Name:        orgExpected.Name,
				UrlCode:     orgExpected.UrlCode,
				Description: orgExpected.Description,
				IsTrial:     false,
				OwnerEmail:  "test@test.org",
			},
			expected:    orgExpected.Organization,
			expectedErr: nil,
		},
		{
			desc: "success trial",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.orgStorage.(*storagemock.MockOrganizationStorage).EXPECT().CreateOrganization(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().CreateProject(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().CreateEnvironmentV2(
					gomock.Any(), gomock.Any(),
				).Times(2).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().CreateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &proto.CreateOrganizationRequest{
				Name:        trialOrgExpected.Name,
				UrlCode:     trialOrgExpected.UrlCode,
				Description: trialOrgExpected.Description,
				IsTrial:     trialOrgExpected.Trial,
				OwnerEmail:  "test@test.org",
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

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.UpdateOrganizationRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrInvalidOrganizationName: empty name",
			setup: nil,
			req: &proto.UpdateOrganizationRequest{
				Id:   "id-0",
				Name: wrapperspb.String(""),
			},
			expectedErr: statusOrganizationNameRequired.Err(),
		},
		{
			desc:  "err: ErrInvalidOrganizationName: only space",
			setup: nil,
			req: &proto.UpdateOrganizationRequest{
				Id:   "id-0",
				Name: wrapperspb.String("    "),
			},
			expectedErr: statusOrganizationNameRequired.Err(),
		},
		{
			desc:  "err: ErrInvalidOrganizationName: max name length exceeded",
			setup: nil,
			req: &proto.UpdateOrganizationRequest{
				Id:   "id-0",
				Name: wrapperspb.String(strings.Repeat("a", 51)),
			},
			expectedErr: statusInvalidOrganizationName.Err(),
		},
		{
			desc: "err: ErrOrganizationNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrOrganizationNotFound)
			},
			req: &proto.UpdateOrganizationRequest{
				Id:   "err-id-0",
				Name: wrapperspb.String("id-0"),
			},
			expectedErr: statusOrganizationNotFound.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
			},
			req: &proto.UpdateOrganizationRequest{
				Id:   "err-id-1",
				Name: wrapperspb.String("id-1"),
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.UpdateOrganizationRequest{
				Id:   "success-id-0",
				Name: wrapperspb.String("id-0"),
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

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.EnableOrganizationRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrOrganizationIDRequired",
			setup: nil,
			req: &proto.EnableOrganizationRequest{
				Id: "",
			},
			expectedErr: statusOrganizationIDRequired.Err(),
		},
		{
			desc: "err: ErrOrganizationNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrOrganizationNotFound)
			},
			req: &proto.EnableOrganizationRequest{
				Id: "id-0",
			},
			expectedErr: statusOrganizationNotFound.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
			},
			req: &proto.EnableOrganizationRequest{
				Id: "id-1",
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &proto.EnableOrganizationRequest{
				Id: "id-1",
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

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.DisableOrganizationRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrOrganizationIDRequired",
			setup: nil,
			req: &proto.DisableOrganizationRequest{
				Id: "",
			},
			expectedErr: statusOrganizationIDRequired.Err(),
		},
		{
			desc: "err: ErrOrganizationNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrOrganizationNotFound)
			},
			req: &proto.DisableOrganizationRequest{
				Id: "id-0",
			},
			expectedErr: statusOrganizationNotFound.Err(),
		},
		{
			desc: "err: ErrCannotUpdateSystemAdmin",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(domain.ErrCannotDisableSystemAdmin)
			},
			req: &proto.DisableOrganizationRequest{
				Id: "id-0",
			},
			expectedErr: statusCannotUpdateSystemAdmin.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
			},
			req: &proto.DisableOrganizationRequest{
				Id: "id-1",
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &proto.DisableOrganizationRequest{
				Id: "id-1",
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

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.ArchiveOrganizationRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrOrganizationIDRequired",
			setup: nil,
			req: &proto.ArchiveOrganizationRequest{
				Id: "",
			},
			expectedErr: statusOrganizationIDRequired.Err(),
		},
		{
			desc: "err: ErrOrganizationNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrOrganizationNotFound)
			},
			req: &proto.ArchiveOrganizationRequest{
				Id: "id-0",
			},
			expectedErr: statusOrganizationNotFound.Err(),
		},
		{
			desc: "err: ErrCannotUpdateSystemAdmin",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(domain.ErrCannotArchiveSystemAdmin)
			},
			req: &proto.ArchiveOrganizationRequest{
				Id: "id-0",
			},
			expectedErr: statusCannotUpdateSystemAdmin.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
			},
			req: &proto.ArchiveOrganizationRequest{
				Id: "id-1",
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &proto.ArchiveOrganizationRequest{
				Id: "id-1",
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

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.UnarchiveOrganizationRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrOrganizationIDRequired",
			setup: nil,
			req: &proto.UnarchiveOrganizationRequest{
				Id: "",
			},
			expectedErr: statusOrganizationIDRequired.Err(),
		},
		{
			desc: "err: ErrOrganizationNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrOrganizationNotFound)
			},
			req: &proto.UnarchiveOrganizationRequest{
				Id: "id-0",
			},
			expectedErr: statusOrganizationNotFound.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error"))
			},
			req: &proto.UnarchiveOrganizationRequest{
				Id: "id-1",
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &proto.UnarchiveOrganizationRequest{
				Id: "id-1",
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

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.ConvertTrialOrganizationRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrOrganizationIDRequired",
			setup: nil,
			req: &proto.ConvertTrialOrganizationRequest{
				Id: "",
			},
			expectedErr: statusOrganizationIDRequired.Err(),
		},
		{
			desc: "err: ErrOrganizationNotFound",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrOrganizationNotFound)
			},
			req: &proto.ConvertTrialOrganizationRequest{
				Id: "id-0",
			},
			expectedErr: statusOrganizationNotFound.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error"))
			},
			req: &proto.ConvertTrialOrganizationRequest{
				Id: "id-1",
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &proto.ConvertTrialOrganizationRequest{
				Id: "id-1",
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

func TestEnvironmentService_CreateDemoOrganization(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createDemoContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	orgExpected, err := domain.NewOrganization(
		"name",
		"url-code",
		"test@example.com",
		"description",
		false,
		false,
		newTestAuthSettings(),
	)
	require.NoError(t, err)

	patterns := []struct {
		desc        string
		setup       func(*EnvironmentService)
		req         *proto.CreateDemoOrganizationRequest
		expected    *proto.Organization
		expectedErr error
	}{
		{
			desc:        "err: ErrInvalidOrganizationName: empty name",
			setup:       nil,
			req:         &proto.CreateDemoOrganizationRequest{Name: ""},
			expectedErr: statusOrganizationNameRequired.Err(),
		},
		{
			desc:        "err: ErrInvalidOrganizationName: only space",
			setup:       nil,
			req:         &proto.CreateDemoOrganizationRequest{Name: "    "},
			expectedErr: statusOrganizationNameRequired.Err(),
		},
		{
			desc:        "err: ErrInvalidOrganizationName: max name length exceeded",
			setup:       nil,
			req:         &proto.CreateDemoOrganizationRequest{Name: strings.Repeat("a", 51)},
			expectedErr: statusInvalidOrganizationName.Err(),
		},
		{
			desc:        "err: ErrInvalidOrganizationUrlCode: can't use uppercase",
			setup:       nil,
			req:         &proto.CreateDemoOrganizationRequest{Name: "id-1", UrlCode: "CODE"},
			expectedErr: statusInvalidOrganizationUrlCode.Err(),
		},
		{
			desc:        "err: ErrInvalidOrganizationUrlCode: max id length exceeded",
			setup:       nil,
			req:         &proto.CreateDemoOrganizationRequest{Name: "id-1", UrlCode: strings.Repeat("a", 51)},
			expectedErr: statusInvalidOrganizationUrlCode.Err(),
		},
		{
			desc: "err: ErrOrganizationAlreadyExists: duplicate id",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrOrganizationAlreadyExists)
			},
			req:         &proto.CreateDemoOrganizationRequest{Name: "id-0", UrlCode: "id-0"},
			expectedErr: statusOrganizationAlreadyExists.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error"))
			},
			req:         &proto.CreateDemoOrganizationRequest{Name: "id-1", UrlCode: "id-1"},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "error")).Err(),
		},
		{
			desc: "success",
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.orgStorage.(*storagemock.MockOrganizationStorage).EXPECT().CreateOrganization(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.projectStorage.(*storagemock.MockProjectStorage).EXPECT().CreateProject(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.environmentStorage.(*storagemock.MockEnvironmentStorage).EXPECT().CreateEnvironmentV2(
					gomock.Any(), gomock.Any(),
				).Return(nil).Times(2)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().CreateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.CreateDemoOrganizationRequest{
				Name:        orgExpected.Name,
				UrlCode:     orgExpected.UrlCode,
				Description: orgExpected.Description,
			},
			expected:    orgExpected.Organization,
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newDemoEnvironmentService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.CreateDemoOrganization(ctx, p.req)
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

func TestValidateOwnershipTransfer(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	// Create context with system admin token for most tests
	ctxAdmin := createContextWithToken(t)
	ctxAdmin = metadata.NewIncomingContext(ctxAdmin, metadata.MD{
		"accept-language": []string{"ja"},
	})

	// Create context with non-admin token for ownership validation tests
	ctxOwner := createContextWithTokenRoleUnassigned(t)
	ctxOwner = metadata.NewIncomingContext(ctxOwner, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*EnvironmentService)
		req         *proto.UpdateOrganizationRequest
		expectedErr error
	}{
		{
			desc: "success: no ownership transfer (name change only)",
			ctx:  ctxAdmin,
			setup: func(s *EnvironmentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &proto.UpdateOrganizationRequest{
				Id:   "org-1",
				Name: wrapperspb.String("New Organization Name"),
			},
			expectedErr: nil,
		},
		{
			desc: "err: no-op ownership transfer (same owner email)",
			ctx:  ctxAdmin,
			setup: func(s *EnvironmentService) {
				s.orgStorage.(*storagemock.MockOrganizationStorage).EXPECT().GetOrganization(
					gomock.Any(), "org-1",
				).Return(&domain.Organization{
					Organization: &proto.Organization{
						Id:         "org-1",
						OwnerEmail: "current-owner@example.com",
					},
				}, nil)
			},
			req: &proto.UpdateOrganizationRequest{
				Id:         "org-1",
				OwnerEmail: wrapperspb.String("current-owner@example.com"),
			},
			expectedErr: statusNoCommand.Err(),
		},
		{
			desc: "err: non-owner trying to transfer ownership",
			ctx:  ctxOwner,
			setup: func(s *EnvironmentService) {
				s.accountClient.(*acmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(&accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_OWNER,
					},
				}, nil)
				s.orgStorage.(*storagemock.MockOrganizationStorage).EXPECT().GetOrganization(
					gomock.Any(), "org-1",
				).Return(&domain.Organization{
					Organization: &proto.Organization{
						Id:         "org-1",
						OwnerEmail: "current-owner@example.com", // Different from token email
					},
				}, nil)
			},
			req: &proto.UpdateOrganizationRequest{
				Id:         "org-1",
				OwnerEmail: wrapperspb.String("new-owner@example.com"),
			},
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc: "err: new owner account not found",
			ctx:  ctxAdmin,
			setup: func(s *EnvironmentService) {
				s.orgStorage.(*storagemock.MockOrganizationStorage).EXPECT().GetOrganization(
					gomock.Any(), "org-1",
				).Return(&domain.Organization{
					Organization: &proto.Organization{
						Id:         "org-1",
						OwnerEmail: "current-owner@example.com",
					},
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), "new-owner@example.com", "org-1",
				).Return(nil, v2as.ErrAccountNotFound)
			},
			req: &proto.UpdateOrganizationRequest{
				Id:         "org-1",
				OwnerEmail: wrapperspb.String("new-owner@example.com"),
			},
			expectedErr: statusAccountNotFound.Err(),
		},
		{
			desc: "err: new owner account is disabled",
			ctx:  ctxAdmin,
			setup: func(s *EnvironmentService) {
				s.orgStorage.(*storagemock.MockOrganizationStorage).EXPECT().GetOrganization(
					gomock.Any(), "org-1",
				).Return(&domain.Organization{
					Organization: &proto.Organization{
						Id:         "org-1",
						OwnerEmail: "current-owner@example.com",
					},
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), "new-owner@example.com", "org-1",
				).Return(&accountdomain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:    "new-owner@example.com",
						Disabled: true,
					},
				}, nil)
			},
			req: &proto.UpdateOrganizationRequest{
				Id:         "org-1",
				OwnerEmail: wrapperspb.String("new-owner@example.com"),
			},
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc: "success: valid ownership transfer passes validation",
			ctx:  ctxAdmin,
			setup: func(s *EnvironmentService) {
				// Mock validation phase - validateOwnershipTransfer
				s.orgStorage.(*storagemock.MockOrganizationStorage).EXPECT().GetOrganization(
					gomock.Any(), "org-1",
				).Return(&domain.Organization{
					Organization: &proto.Organization{
						Id:         "org-1",
						OwnerEmail: "current-owner@example.com",
					},
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), "new-owner@example.com", "org-1",
				).Return(&accountdomain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:    "new-owner@example.com",
						Disabled: false,
					},
				}, nil)

				// Mock transaction execution (simplified)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)

				// Mock updateOwnerRole calls (these happen after transaction)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), "current-owner@example.com", "org-1",
				).Return(&accountdomain.AccountV2{
					AccountV2: &accountproto.AccountV2{Email: "current-owner@example.com"},
				}, nil).AnyTimes()
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil).AnyTimes()
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), "new-owner@example.com", "org-1",
				).Return(&accountdomain.AccountV2{
					AccountV2: &accountproto.AccountV2{Email: "new-owner@example.com"},
				}, nil).AnyTimes()
			},
			req: &proto.UpdateOrganizationRequest{
				Id:         "org-1",
				OwnerEmail: wrapperspb.String("new-owner@example.com"),
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
			_, err := service.UpdateOrganization(p.ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
