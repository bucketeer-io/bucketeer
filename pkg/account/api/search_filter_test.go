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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	accstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	mysql "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

func TestCreateSearchFilter(t *testing.T) {
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
		req         *accountproto.CreateSearchFilterRequest
		expectedErr error
	}{
		{
			desc: "err: role is not allowed to create search filter",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_UNASSIGNED,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.CreateSearchFilterRequest{
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.CreateSearchFilterCommand{
					Name:             "filter",
					Query:            "query",
					FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "envID0",
					DefaultFilter:    false,
				},
			},
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc: "err: email is empty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.CreateSearchFilterRequest{
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.CreateSearchFilterCommand{
					Name:             "filter",
					Query:            "query",
					FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "envID0",
					DefaultFilter:    false,
				},
			},
			expectedErr: statusEmailIsEmpty.Err(),
		},
		{
			desc: "err: organization_id is empty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:         "bucketeer@example.com",
				EnvironmentId: "envID0",
				Command: &accountproto.CreateSearchFilterCommand{
					Name:             "filter",
					Query:            "query",
					FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "envID0",
					DefaultFilter:    false,
				},
			},
			expectedErr: statusMissingOrganizationID.Err(),
		},
		{
			desc: "err: internal error",
			setup: func(s *AccountService) {
				account := domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&account, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrSystemAdminAccountNotFound)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.CreateSearchFilterCommand{
					Name:             "filter",
					Query:            "query",
					FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "envID0",
					DefaultFilter:    false,
				},
			},
			expectedErr: statusInternal.Err(),
		},
		{
			desc: "err: account not found",
			setup: func(s *AccountService) {
				account := domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&account, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrSystemAdminAccountNotFound)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.CreateSearchFilterCommand{
					Name:             "filter",
					Query:            "query",
					FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "envID0",
					DefaultFilter:    false,
				},
			},
			expectedErr: statusAccountNotFound.Err(),
		},
		{
			desc: "err: command is nil",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command:        nil,
			},
			expectedErr: statusNoCommand.Err(),
		},
		{
			desc: "err: SearchFilter Name is empty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.CreateSearchFilterCommand{
					Name: "",
				},
			},
			expectedErr: statusSearchFilterNameIsEmpty.Err(),
		},
		{
			desc: "err: SearchFilter Query is empty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.CreateSearchFilterCommand{
					Name:  "name",
					Query: "",
				},
			},
			expectedErr: statusSearchFilterQueryIsEmpty.Err(),
		},
		{
			desc: "err: SearchFilter targetFilter is unknown",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.CreateSearchFilterCommand{
					Name:             "name",
					Query:            "query",
					FilterTargetType: accountproto.FilterTargetType_UNKNOWN,
				},
			},
			expectedErr: statusSearchFilterTargetTypeIsRequired.Err(),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrSystemAdminAccountNotFound)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{},
				}, nil)
				s.publisher.(*mock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.CreateSearchFilterCommand{
					Name:             "filter",
					Query:            "query",
					FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "envID0",
					DefaultFilter:    false,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: system admin account",
			setup: func(s *AccountService) {
				email := "bucketeer@example.com"
				orgID := "system_admin_org_id"
				acc := &domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            email,
						OrganizationId:   orgID,
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), "email", gomock.Any(),
				).Return(acc, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), email,
				).Return(acc, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), email, orgID,
				).Return(acc, nil)

				s.publisher.(*mock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.CreateSearchFilterCommand{
					Name:             "filter",
					Query:            "query",
					FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "envID0",
					DefaultFilter:    false,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: already has default filter",
			setup: func(s *AccountService) {
				account := domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
						SearchFilters: []*accountproto.SearchFilter{
							{
								Id:               "id",
								Name:             "filter",
								Query:            "query",
								FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
								EnvironmentId:    "envID0",
								DefaultFilter:    true,
							},
						},
					},
				}
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&account, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrSystemAdminAccountNotFound)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{},
				}, nil)
				s.publisher.(*mock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.CreateSearchFilterCommand{
					Name:             "filter",
					Query:            "query",
					FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    "envID0",
					DefaultFilter:    false,
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
			_, err := service.CreateSearchFilter(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestUpdateSearchFilter(t *testing.T) {
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
		req         *accountproto.UpdateSearchFilterRequest
		expectedErr error
	}{
		{
			desc: "err: role is not allowed to update search filter",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_UNASSIGNED,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.UpdateSearchFilterRequest{
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				ChangeNameCommand: &accountproto.ChangeSearchFilterNameCommand{
					Id:   "searchFilterID",
					Name: "filter",
				},
			},
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc: "err: email is empty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.UpdateSearchFilterRequest{
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				ChangeNameCommand: &accountproto.ChangeSearchFilterNameCommand{
					Id:   "searchFilterID",
					Name: "filter",
				},
			},
			expectedErr: statusEmailIsEmpty.Err(),
		},
		{
			desc: "err: organization_id is empty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.UpdateSearchFilterRequest{
				Email:         "bucketeer@example.com",
				EnvironmentId: "envID0",
				ChangeNameCommand: &accountproto.ChangeSearchFilterNameCommand{
					Id:   "searchFilterID",
					Name: "filter",
				},
			},
			expectedErr: statusMissingOrganizationID.Err(),
		},
		{
			desc: "err: internal error",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrSystemAdminAccountNotFound)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal("account", "test"))
			},
			req: &accountproto.UpdateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				ChangeNameCommand: &accountproto.ChangeSearchFilterNameCommand{
					Id:   "searchFilterID",
					Name: "filter",
				},
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal("account", "test")).Err(),
		},
		{
			desc: "err: account not found",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrSystemAdminAccountNotFound)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			req: &accountproto.UpdateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				ChangeNameCommand: &accountproto.ChangeSearchFilterNameCommand{
					Id:   "searchFilterID",
					Name: "filter",
				},
			},
			expectedErr: statusAccountNotFound.Err(),
		},
		{
			desc: "err: command is empty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.UpdateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
			},
			expectedErr: statusNoCommand.Err(),
		},
		{
			desc: "err: SearchFilter ID is empty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.UpdateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				ChangeNameCommand: &accountproto.ChangeSearchFilterNameCommand{
					Id:   "",
					Name: "filter",
				},
			},
			expectedErr: statusSearchFilterIDIsEmpty.Err(),
		},
		{
			desc: "err: SearchFilter Name is empty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.UpdateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				ChangeNameCommand: &accountproto.ChangeSearchFilterNameCommand{
					Id:   "tesID",
					Name: "",
				},
			},
			expectedErr: statusSearchFilterNameIsEmpty.Err(),
		},
		{
			desc: "err: SearchFilter ID is empty for ChangeNameCommand",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.UpdateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				ChangeNameCommand: &accountproto.ChangeSearchFilterNameCommand{
					Id:   "",
					Name: "update-name",
				},
			},
			expectedErr: statusSearchFilterIDIsEmpty.Err(),
		},
		{
			desc: "err: SearchFilter Query is empty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.UpdateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				ChangeQueryCommand: &accountproto.ChangeSearchFilterQueryCommand{
					Id:    "tesID",
					Query: "",
				},
			},
			expectedErr: statusSearchFilterQueryIsEmpty.Err(),
		},
		{
			desc: "err: SearchFilter ID is empty for ChangeQueryCommand",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.UpdateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				ChangeQueryCommand: &accountproto.ChangeSearchFilterQueryCommand{
					Query: "update-query",
				},
			},
			expectedErr: statusSearchFilterIDIsEmpty.Err(),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrSystemAdminAccountNotFound)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					println(err)
				}).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
						SearchFilters: []*accountproto.SearchFilter{
							{
								Id: "tesID",
							},
						},
					},
				}, nil)
				s.publisher.(*mock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &accountproto.UpdateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				ChangeNameCommand: &accountproto.ChangeSearchFilterNameCommand{
					Id:   "tesID",
					Name: "update-name",
				},
				ChangeQueryCommand: &accountproto.ChangeSearchFilterQueryCommand{
					Id:    "tesID",
					Query: "query",
				},
				ChangeDefaultFilterCommand: &accountproto.ChangeDefaultSearchFilterCommand{
					Id:            "tesID",
					DefaultFilter: true,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: system admin account",
			setup: func(s *AccountService) {
				email := "bucketeer@example.com"
				orgID := "system_admin_org_id"
				acc := &domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            email,
						OrganizationId:   orgID,
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
						SearchFilters: []*accountproto.SearchFilter{
							{
								Id: "tesID",
							},
						},
					},
				}

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), "email", gomock.Any(),
				).Return(acc, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), email,
				).Return(acc, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), email, orgID,
				).Return(acc, nil)

				s.publisher.(*mock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil).Times(3)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
			},
			req: &accountproto.UpdateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				ChangeNameCommand: &accountproto.ChangeSearchFilterNameCommand{
					Id:   "tesID",
					Name: "update-name",
				},
				ChangeQueryCommand: &accountproto.ChangeSearchFilterQueryCommand{
					Id:    "tesID",
					Query: "query",
				},
				ChangeDefaultFilterCommand: &accountproto.ChangeDefaultSearchFilterCommand{
					Id:            "tesID",
					DefaultFilter: true,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: already has default filter",
			setup: func(s *AccountService) {
				account := domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
						SearchFilters: []*accountproto.SearchFilter{
							{
								Id:               "id",
								Name:             "filter",
								Query:            "query",
								FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
								EnvironmentId:    "envID0",
								DefaultFilter:    true,
							},
						},
					},
				}
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&account, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrSystemAdminAccountNotFound)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{},
				}, nil)
			},
			req: &accountproto.UpdateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				ChangeDefaultFilterCommand: &accountproto.ChangeDefaultSearchFilterCommand{
					Id:            "id",
					DefaultFilter: false,
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
			_, err := service.UpdateSearchFilter(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestDeleteSearchFilter(t *testing.T) {
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
		req         *accountproto.DeleteSearchFilterRequest
		expectedErr error
	}{
		{
			desc: "err: role is not allowed",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_UNASSIGNED,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.DeleteSearchFilterRequest{
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.DeleteSearchFilterCommand{
					Id: "filterID",
				},
			},
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc: "err: email is empty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.DeleteSearchFilterRequest{
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.DeleteSearchFilterCommand{
					Id: "filterID",
				},
			},
			expectedErr: statusEmailIsEmpty.Err(),
		},
		{
			desc: "err: organization_id is empty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.DeleteSearchFilterRequest{
				Email:         "bucketeer@example.com",
				EnvironmentId: "envID0",
				Command: &accountproto.DeleteSearchFilterCommand{
					Id: "filterID",
				},
			},
			expectedErr: statusMissingOrganizationID.Err(),
		},
		{
			desc: "err: internal error",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrSystemAdminAccountNotFound)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal("account", "test"))
			},
			req: &accountproto.DeleteSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.DeleteSearchFilterCommand{
					Id: "filterID",
				},
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal("account", "test")).Err(),
		},
		{
			desc: "err: account not found",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrSystemAdminAccountNotFound)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			req: &accountproto.DeleteSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.DeleteSearchFilterCommand{
					Id: "filterID",
				},
			},
			expectedErr: statusAccountNotFound.Err(),
		},
		{
			desc: "err: command is nil",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.DeleteSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command:        nil,
			},
			expectedErr: statusNoCommand.Err(),
		},
		{
			desc: "err: SearchFilterID is empty",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
			},
			req: &accountproto.DeleteSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.DeleteSearchFilterCommand{
					Id: "",
				},
			},
			expectedErr: statusSearchFilterIDIsEmpty.Err(),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrSystemAdminAccountNotFound)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
						SearchFilters: []*accountproto.SearchFilter{
							{
								Id:               "filterID",
								Name:             "filter",
								Query:            "query",
								FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
								EnvironmentId:    "envID0",
								DefaultFilter:    true,
							},
						},
					},
				}, nil)
				s.publisher.(*mock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &accountproto.DeleteSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.DeleteSearchFilterCommand{
					Id: "filterID",
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: system admin account",
			setup: func(s *AccountService) {
				email := "bucketeer@example.com"
				orgID := "system_admin_org_id"
				acc := &domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationId:   orgID,
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "envID0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
						SearchFilters: []*accountproto.SearchFilter{
							{
								Id: "filterID",
							},
						},
					},
				}
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), "email", gomock.Any(),
				).Return(acc, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), email,
				).Return(acc, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), email, orgID,
				).Return(acc, nil)

				s.publisher.(*mock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
			},
			req: &accountproto.DeleteSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.DeleteSearchFilterCommand{
					Id: "filterID",
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
			_, err := service.DeleteSearchFilter(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}
