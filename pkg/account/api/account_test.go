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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	accstoragemock "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2/mock"
	alstoragemock "github.com/bucketeer-io/bucketeer/pkg/auditlog/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
)

func TestCreateAccountV2MySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
		setup       func(*AccountService)
		req         *accountproto.CreateAccountV2Request
		expectedErr error
	}{
		{
			desc: "errEmailIsEmpty",
			req: &accountproto.CreateAccountV2Request{
				Command:        &accountproto.CreateAccountV2Command{Email: ""},
				OrganizationId: "org0",
			},
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			expectedErr: createError(statusEmailIsEmpty, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email")),
		},
		{
			desc: "errInvalidEmail",
			req: &accountproto.CreateAccountV2Request{
				Command:        &accountproto.CreateAccountV2Command{Email: "bucketeer@"},
				OrganizationId: "org0",
			},
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			expectedErr: createError(statusInvalidEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errAccountAlreadyExists",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountAlreadyExists)
			},
			req: &accountproto.CreateAccountV2Request{
				Command: &accountproto.CreateAccountV2Command{
					Email:            "bucketeer_environment@example.com",
					OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
					EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
						{
							Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							EnvironmentId: "test",
						},
					},
				},
				OrganizationId: "org0",
			},
			expectedErr: createError(statusAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			req: &accountproto.CreateAccountV2Request{
				Command: &accountproto.CreateAccountV2Command{
					Email:            "bucketeer@example.com",
					OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
						{
							Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							EnvironmentId: "test",
						},
					},
				},
				OrganizationId: "org0",
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().CreateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &accountproto.CreateAccountV2Request{
				Command: &accountproto.CreateAccountV2Command{
					Email:            "bucketeer@example.com",
					OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
						{
							Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							EnvironmentId: "test",
						},
					},
				},
				OrganizationId: "org0",
			},
			expectedErr: nil,
		},
		{
			desc: "success: with admin role",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().CreateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &accountproto.CreateAccountV2Request{
				Command: &accountproto.CreateAccountV2Command{
					Email:            "bucketeer@example.com",
					OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					EnvironmentRoles: nil,
				},
				OrganizationId: "org0",
			},
			expectedErr: nil,
		},
		{
			desc: "success: create admin account with environment roles",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().CreateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &accountproto.CreateAccountV2Request{
				Command: &accountproto.CreateAccountV2Command{
					Email:            "bucketeer@example.com",
					OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
						{
							Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							EnvironmentId: "test",
						},
					},
				},
				OrganizationId: "org0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, false)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.CreateAccountV2(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestCreateAccountV2NoCommandMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
		setup       func(*AccountService)
		req         *accountproto.CreateAccountV2Request
		expectedErr error
	}{
		{
			desc: "errEmailIsEmpty",
			req: &accountproto.CreateAccountV2Request{
				Email:          "",
				OrganizationId: "org0",
			},
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			expectedErr: createError(statusEmailIsEmpty, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email")),
		},
		{
			desc: "errInvalidEmail",
			req: &accountproto.CreateAccountV2Request{
				Email:          "bucketeer@",
				OrganizationId: "org0",
			},
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			expectedErr: createError(statusInvalidEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errAccountAlreadyExists",
			req: &accountproto.CreateAccountV2Request{
				Email:            "bucketeer_environment@example.com",
				FirstName:        "Test",
				LastName:         "User",
				Language:         "en",
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				OrganizationId:   "org0",
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{
						Role:          accountproto.AccountV2_Role_Environment_VIEWER,
						EnvironmentId: "test",
					},
				},
			},
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountAlreadyExists)
			},
			expectedErr: createError(statusAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "errInternal",
			req: &accountproto.CreateAccountV2Request{
				Email:            "bucketeer@example.com",
				FirstName:        "Test",
				LastName:         "User",
				Language:         "en",
				OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
				OrganizationId:   "org0",
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{
						Role:          accountproto.AccountV2_Role_Environment_VIEWER,
						EnvironmentId: "test",
					},
				},
			},
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			req: &accountproto.CreateAccountV2Request{
				Email:            "bucketeer@example.com",
				FirstName:        "Test",
				LastName:         "User",
				Language:         "en",
				OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
				OrganizationId:   "org0",
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{
						Role:          accountproto.AccountV2_Role_Environment_VIEWER,
						EnvironmentId: "test",
					},
				},
			},
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.adminAuditLogStorage.(*alstoragemock.MockAdminAuditLogStorage).EXPECT().CreateAdminAuditLog(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, false)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.CreateAccountV2(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestUpdateAccountV2MySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
		setup       func(*AccountService)
		req         *accountproto.UpdateAccountV2Request
		expectedErr error
	}{
		{
			desc: "errEmailIsEmpty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				OrganizationId: "org0",
			},
			expectedErr: createError(statusEmailIsEmpty, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email")),
		},
		{
			desc: "errInvalidEmail",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@",
				OrganizationId: "org0",
				ChangeFirstNameCommand: &accountproto.ChangeAccountV2FirstNameCommand{
					FirstName: "newFirstName",
				},
				ChangeLastNameCommand: &accountproto.ChangeAccountV2LastNameCommand{
					LastName: "newLastName",
				},
			},
			expectedErr: createError(statusInvalidEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errOrganizationIDIsEmpty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email: "bucketeer@example.com",
				ChangeFirstNameCommand: &accountproto.ChangeAccountV2FirstNameCommand{
					FirstName: "newFirstName",
				},
				ChangeLastNameCommand: &accountproto.ChangeAccountV2LastNameCommand{
					LastName: "newLastName",
				},
			},
			expectedErr: createError(statusMissingOrganizationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id")),
		},
		{
			desc: "errInvalidNewName",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				ChangeFirstNameCommand: &accountproto.ChangeAccountV2FirstNameCommand{
					FirstName: strings.Repeat("a", 251),
				},
			},
			expectedErr: createError(statusInvalidFirstName, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "first_name")),
		},
		{
			desc: "errInvalidNewOrganizationRole",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				ChangeOrganizationRoleCommand: &accountproto.ChangeAccountV2OrganizationRoleCommand{
					Role: accountproto.AccountV2_Role_Organization_UNASSIGNED,
				},
			},
			expectedErr: createError(statusInvalidOrganizationRole, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "organization_role")),
		},
		{
			desc: "errAccountNotFound",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				ChangeFirstNameCommand: &accountproto.ChangeAccountV2FirstNameCommand{
					FirstName: "newFirstName",
				},
				ChangeLastNameCommand: &accountproto.ChangeAccountV2LastNameCommand{
					LastName: "newLastName",
				},
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{
						Role: accountproto.AccountV2_Role_Environment_EDITOR,
					},
				},
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				ChangeFirstNameCommand: &accountproto.ChangeAccountV2FirstNameCommand{
					FirstName: "newFirstName",
				},
				ChangeLastNameCommand: &accountproto.ChangeAccountV2LastNameCommand{
					LastName: "newLastName",
				},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				ChangeFirstNameCommand: &accountproto.ChangeAccountV2FirstNameCommand{
					FirstName: "newFirstName",
				},
				ChangeLastNameCommand: &accountproto.ChangeAccountV2LastNameCommand{
					LastName: "newLastName",
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: update admin account",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				ChangeFirstNameCommand: &accountproto.ChangeAccountV2FirstNameCommand{
					FirstName: "newFirstName",
				},
				ChangeLastNameCommand: &accountproto.ChangeAccountV2LastNameCommand{
					LastName: "newLastName",
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: update member to admin",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)

				// This is the GetAccountV2 call inside the transaction
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "env1",
								Role:          accountproto.AccountV2_Role_Environment_EDITOR,
							},
						},
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				ChangeOrganizationRoleCommand: &accountproto.ChangeAccountV2OrganizationRoleCommand{
					Role: accountproto.AccountV2_Role_Organization_ADMIN,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, false)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UpdateAccountV2(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestUpdateAccountV2NoCommandMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
		setup       func(*AccountService)
		req         *accountproto.UpdateAccountV2Request
		expectedErr error
	}{
		{
			desc: "errEmailIsEmpty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				OrganizationId: "org0",
			},
			expectedErr: createError(statusEmailIsEmpty, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email")),
		},
		{
			desc: "errInvalidEmail",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@",
				OrganizationId: "org0",
				FirstName:      wrapperspb.String("newFirstName"),
				LastName:       wrapperspb.String("newLastName"),
			},
			expectedErr: createError(statusInvalidEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errOrganizationIDIsEmpty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:     "bucketeer@example.com",
				FirstName: wrapperspb.String("newFirstName"),
				LastName:  wrapperspb.String("newLastName"),
			},
			expectedErr: createError(statusMissingOrganizationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id")),
		},
		{
			desc: "errInvalidNewName",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				FirstName:      wrapperspb.String(strings.Repeat("a", 251)),
			},
			expectedErr: createError(statusInvalidFirstName, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "first_name")),
		},
		{
			desc: "errInvalidNewOrganizationRole",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				OrganizationRole: &accountproto.UpdateAccountV2Request_OrganizationRoleValue{
					Role: accountproto.AccountV2_Role_Organization_UNASSIGNED,
				},
			},
			expectedErr: createError(statusInvalidOrganizationRole, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "organization_role")),
		},
		{
			desc: "errAccountNotFound",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				FirstName:      wrapperspb.String("newFirstName"),
				LastName:       wrapperspb.String("newLastName"),
				OrganizationRole: &accountproto.UpdateAccountV2Request_OrganizationRoleValue{
					Role: accountproto.AccountV2_Role_Organization_ADMIN,
				},
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{
						Role: accountproto.AccountV2_Role_Environment_EDITOR,
					},
				},
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				FirstName:      wrapperspb.String("newFirstName"),
				LastName:       wrapperspb.String("newLastName"),
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{
						Role: accountproto.AccountV2_Role_Environment_EDITOR,
					},
				},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationId:   "org0",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationId:   "org0",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)

				s.adminAuditLogStorage.(*alstoragemock.MockAdminAuditLogStorage).EXPECT().CreateAdminAuditLog(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				FirstName:      wrapperspb.String("newFirstName"),
				LastName:       wrapperspb.String("newLastName"),
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{
						Role: accountproto.AccountV2_Role_Environment_EDITOR,
					},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, false)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UpdateAccountV2(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestEnableAccountV2MySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
		setup       func(*AccountService)
		req         *accountproto.EnableAccountV2Request
		expectedErr error
	}{
		{
			desc: "errEmailIsEmpty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.EnableAccountV2Request{
				OrganizationId: "org0",
			},
			expectedErr: createError(statusEmailIsEmpty, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email")),
		},
		{
			desc: "errInvalidEmail",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.EnableAccountV2Request{
				Email:          "bucketeer@",
				OrganizationId: "org0",
			},
			expectedErr: createError(statusInvalidEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errOrganizationIDIsEmpty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.EnableAccountV2Request{
				Email: "bucketeer@example.com",
			},
			expectedErr: createError(statusMissingOrganizationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id")),
		},
		{
			desc: "errAccountNotFound",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			req: &accountproto.EnableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			req: &accountproto.EnableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationId:   "org0",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil).Times(2)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)

				s.adminAuditLogStorage.(*alstoragemock.MockAdminAuditLogStorage).EXPECT().CreateAdminAuditLog(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &accountproto.EnableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, false)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.EnableAccountV2(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestDisableAccountV2MySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
		setup       func(*AccountService)
		req         *accountproto.DisableAccountV2Request
		expectedErr error
	}{
		{
			desc: "errEmailIsEmpty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.DisableAccountV2Request{
				OrganizationId: "org0",
			},
			expectedErr: createError(statusEmailIsEmpty, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email")),
		},
		{
			desc: "errInvalidEmail",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.DisableAccountV2Request{
				Email:          "bucketeer@",
				OrganizationId: "org0",
			},
			expectedErr: createError(statusInvalidEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errOrganizationIDIsEmpty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.DisableAccountV2Request{
				Email: "bucketeer@example.com",
			},
			expectedErr: createError(statusMissingOrganizationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id")),
		},
		{
			desc: "errAccountNotFound",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			req: &accountproto.DisableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			req: &accountproto.DisableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						Language:         "en",
						OrganizationId:   "org0",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)

				s.adminAuditLogStorage.(*alstoragemock.MockAdminAuditLogStorage).EXPECT().CreateAdminAuditLog(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &accountproto.DisableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, false)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.DisableAccountV2(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestDeleteAccountV2MySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
		setup       func(*AccountService)
		req         *accountproto.DeleteAccountV2Request
		expectedErr error
	}{
		{
			desc: "errEmailIsEmpty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.DeleteAccountV2Request{
				OrganizationId: "org0",
			},
			expectedErr: createError(statusEmailIsEmpty, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email")),
		},
		{
			desc: "errInvalidEmail",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.DeleteAccountV2Request{
				Email:          "bucketeer@",
				OrganizationId: "org0",
			},
			expectedErr: createError(statusInvalidEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errOrganizationIDIsEmpty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.DeleteAccountV2Request{
				Email: "bucketeer@example.com",
			},
			expectedErr: createError(statusMissingOrganizationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id")),
		},
		{
			desc: "errAccountNotFound",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			req: &accountproto.DeleteAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			req: &accountproto.DeleteAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationId:   "org0",
						Language:         "en",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil).Times(2)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().DeleteAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &accountproto.DeleteAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, false)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.DeleteAccountV2(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestGetAccountV2ByEnvironmentIDMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
		setup       func(*AccountService)
		req         *accountproto.GetAccountV2ByEnvironmentIDRequest
		expectedErr error
	}{
		{
			desc: "errInvalidEmail",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), "email", gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.GetAccountV2ByEnvironmentIDRequest{
				Email:         "bucketeer@",
				EnvironmentId: "env0",
			},
			expectedErr: createError(statusInvalidEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errAccountNotFound",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), "email", gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), "bucketeer@example.com", gomock.Any(),
				).Return(nil, v2as.ErrAccountNotFound)
			},
			req: &accountproto.GetAccountV2ByEnvironmentIDRequest{
				Email:         "bucketeer@example.com",
				EnvironmentId: "env0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), "email", gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), "bucketeer@example.com", gomock.Any(),
				).Return(nil, errors.New("test"))
			},
			req: &accountproto.GetAccountV2ByEnvironmentIDRequest{
				Email:         "bucketeer@example.com",
				EnvironmentId: "env0",
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), "email", gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), "bucketeer@example.com", gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.GetAccountV2ByEnvironmentIDRequest{
				Email:         "bucketeer@example.com",
				EnvironmentId: "env0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, false)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.GetAccountV2ByEnvironmentID(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestListAccountsV2MySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithDefaultToken(t, false)
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
		setup       func(*AccountService)
		input       *accountproto.ListAccountsV2Request
		expected    *accountproto.ListAccountsV2Response
		expectedErr error
	}{
		{
			desc: "errInvalidCursor",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			input:       &accountproto.ListAccountsV2Request{OrganizationId: "org0", Cursor: "XXX"},
			expected:    nil,
			expectedErr: createError(statusInvalidCursor, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor")),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().ListAccountsV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), errors.New("test"))
			},
			input:       &accountproto.ListAccountsV2Request{OrganizationId: "org0"},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success with member role",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().ListAccountsV2(
					gomock.Any(), gomock.Any(),
				).Return([]*accountproto.AccountV2{}, 0, int64(0), nil)
			},
			input: &accountproto.ListAccountsV2Request{
				PageSize:       2,
				Cursor:         "",
				OrganizationId: "org0",
			},
			expected: &accountproto.ListAccountsV2Response{Accounts: []*accountproto.AccountV2{}, Cursor: "0"},
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "Test",
						LastName:         "User",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().ListAccountsV2(
					gomock.Any(), gomock.Any(),
				).Return([]*accountproto.AccountV2{}, 0, int64(0), nil)
			},
			input:       &accountproto.ListAccountsV2Request{PageSize: 2, Cursor: "", OrganizationId: "org0"},
			expected:    &accountproto.ListAccountsV2Response{Accounts: []*accountproto.AccountV2{}, Cursor: "0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			actual, err := service.ListAccountsV2(ctx, p.input)
			assert.Equal(t, p.expectedErr, err, p.desc)
			assert.Equal(t, p.expected, actual, p.desc)
		})
	}
}

func setToken(ctx context.Context, isSystemAdmin bool) context.Context {
	t := &token.AccessToken{
		Issuer:        "issuer",
		Audience:      "audience",
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         "email",
		IsSystemAdmin: isSystemAdmin,
	}
	return context.WithValue(ctx, rpc.Key, t)
}
