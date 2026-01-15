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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	accstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	alstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2/mock"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
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
			expectedErr: statusEmailIsEmpty.Err(),
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
			expectedErr: statusInvalidEmail.Err(),
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
			expectedErr: statusAccountAlreadyExists.Err(),
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
				).Return(pkgErr.NewErrorInternal("account", "test"))
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal("account", "test")).Err(),
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

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().CreateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.adminAuditLogStorage.(*alstoragemock.MockAdminAuditLogStorage).EXPECT().CreateAdminAuditLog(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
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
			expectedErr: statusEmailIsEmpty.Err(),
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
			expectedErr: statusInvalidEmail.Err(),
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
			expectedErr: statusMissingOrganizationID.Err(),
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
			expectedErr: statusInvalidFirstName.Err(),
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
			expectedErr: statusInvalidOrganizationRole.Err(),
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
			expectedErr: statusAccountNotFound.Err(),
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
				).Return(pkgErr.NewErrorInternal("account", "test"))
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
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal("account", "test")).Err(),
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
			expectedErr: statusEmailIsEmpty.Err(),
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
			expectedErr: statusInvalidEmail.Err(),
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
			expectedErr: statusMissingOrganizationID.Err(),
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
			expectedErr: statusAccountNotFound.Err(),
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
				).Return(pkgErr.NewErrorInternal("account", "test"))
			},
			req: &accountproto.EnableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal("account", "test")).Err(),
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
			expectedErr: statusEmailIsEmpty.Err(),
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
			expectedErr: statusInvalidEmail.Err(),
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
			expectedErr: statusMissingOrganizationID.Err(),
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
			expectedErr: statusAccountNotFound.Err(),
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
				).Return(pkgErr.NewErrorInternal("account", "test"))
			},
			req: &accountproto.DisableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal("account", "test")).Err(),
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
			expectedErr: statusEmailIsEmpty.Err(),
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
			expectedErr: statusInvalidEmail.Err(),
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
			expectedErr: statusMissingOrganizationID.Err(),
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
			expectedErr: statusAccountNotFound.Err(),
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
				).Return(pkgErr.NewErrorInternal("account", "test"))
			},
			req: &accountproto.DeleteAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal("account", "test")).Err(),
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
			expectedErr: statusInvalidEmail.Err(),
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
			expectedErr: statusAccountNotFound.Err(),
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
				).Return(nil, pkgErr.NewErrorInternal("account", "test"))
			},
			req: &accountproto.GetAccountV2ByEnvironmentIDRequest{
				Email:         "bucketeer@example.com",
				EnvironmentId: "env0",
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal("account", "test")).Err(),
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
			expectedErr: statusInvalidCursor.Err(),
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
				).Return(nil, 0, int64(0), pkgErr.NewErrorInternal("account", "test"))
			},
			input:       &accountproto.ListAccountsV2Request{OrganizationId: "org0"},
			expected:    nil,
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal("account", "test")).Err(),
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
	return context.WithValue(ctx, rpc.AccessTokenKey, t)
}
