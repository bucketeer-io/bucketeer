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
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"

	v2 "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	v2as "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	accstoragemock "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2/mock"
	ecmock "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestGetMeMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	lang := "ja"
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{lang},
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
	org := environmentproto.Organization{Id: "org0"}

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*AccountService)
		input       *accountproto.GetMeRequest
		expected    *accountproto.GetMeResponse
		expectedErr error
	}{
		{
			desc:        "errUnauthenticated",
			ctx:         context.Background(),
			setup:       nil,
			input:       &accountproto.GetMeRequest{},
			expected:    nil,
			expectedErr: createError(statusUnauthenticated, localizer.MustLocalize(locale.UnauthenticatedError)),
		},
		{
			desc:        "errInvalidEmail",
			ctx:         createContextWithInvalidEmailToken(t),
			setup:       nil,
			input:       &accountproto.GetMeRequest{},
			expected:    nil,
			expectedErr: createError(statusInvalidEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errInternal",
			ctx:  createContextWithDefaultToken(t, true),
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListProjects(
					gomock.Any(),
					gomock.Any(),
				).Return(
					nil,
					createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
				)
			},
			input:       &accountproto.GetMeRequest{},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "err: account is disabled",
			ctx:  createContextWithDefaultToken(t, true),
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListProjects(
					gomock.Any(),
					gomock.Any(),
				).Return(
					&environmentproto.ListProjectsResponse{
						Projects: getProjects(t),
						Cursor:   "",
					},
					nil,
				)
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(),
					gomock.Any(),
				).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: getEnvironments(t),
						Cursor:       "",
					},
					nil,
				)
				s.environmentClient.(*ecmock.MockClient).EXPECT().GetOrganization(
					gomock.Any(), gomock.Any(),
				).Return(
					&environmentproto.GetOrganizationResponse{
						Organization: &org,
					},
					nil,
				)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrSystemAdminAccountNotFound)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						FirstName:        "test",
						LastName:         "test",
						Language:         "en",
						AvatarImageUrl:   "",
						OrganizationId:   "org0",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_EDITOR,
							},
						},
						Disabled: true,
					},
				}, nil)
			},
			input: &accountproto.GetMeRequest{
				OrganizationId: "org0",
			},
			expected:    nil,
			expectedErr: createError(statusUnauthenticated, localizer.MustLocalize(locale.UnauthenticatedError)),
		},
		{
			desc: "err: account not found",
			ctx:  createContextWithDefaultToken(t, true),
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListProjects(
					gomock.Any(),
					gomock.Any(),
				).Return(
					&environmentproto.ListProjectsResponse{
						Projects: getProjects(t),
						Cursor:   "",
					},
					nil,
				)
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(),
					gomock.Any(),
				).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: getEnvironments(t),
						Cursor:       "",
					},
					nil,
				)
				s.environmentClient.(*ecmock.MockClient).EXPECT().GetOrganization(
					gomock.Any(), gomock.Any(),
				).Return(
					&environmentproto.GetOrganizationResponse{
						Organization: &org,
					},
					nil,
				)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrSystemAdminAccountNotFound)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, v2.ErrAccountNotFound)
			},
			input: &accountproto.GetMeRequest{
				OrganizationId: "org0",
			},
			expected:    nil,
			expectedErr: createError(statusUnauthenticated, localizer.MustLocalize(locale.UnauthenticatedError)),
		},
		{
			desc: "success",
			ctx:  createContextWithDefaultToken(t, true),
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListProjects(
					gomock.Any(),
					gomock.Any(),
				).Return(
					&environmentproto.ListProjectsResponse{
						Projects: getProjects(t),
						Cursor:   "",
					},
					nil,
				)
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(),
					gomock.Any(),
				).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: getEnvironments(t),
						Cursor:       "",
					},
					nil,
				)
				s.environmentClient.(*ecmock.MockClient).EXPECT().GetOrganization(
					gomock.Any(), gomock.Any(),
				).Return(
					&environmentproto.GetOrganizationResponse{
						Organization: &org,
					},
					nil,
				)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetSystemAdminAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrSystemAdminAccountNotFound)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						AvatarImageUrl:   "",
						OrganizationId:   "org0",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_EDITOR,
							},
						},
						Disabled:  false,
						CreatedAt: 0,
						UpdatedAt: 0,
						SearchFilters: []*accountproto.SearchFilter{
							{
								Id: "search-filter-id",
							},
						},
						FirstName: "test",
						LastName:  "test",
						Language:  "en",
					},
				}, nil).Times(2)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &accountproto.GetMeRequest{
				OrganizationId: "org0",
			},
			expected: &accountproto.GetMeResponse{
				Account: &accountproto.ConsoleAccount{
					Email:            "bucketeer@example.com",
					FirstName:        "test",
					LastName:         "test",
					Language:         "en",
					AvatarUrl:        "",
					Organization:     &org,
					OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					EnvironmentRoles: []*accountproto.ConsoleAccount_EnvironmentRole{
						{
							Environment: &environmentproto.EnvironmentV2{
								Id:        "ns0",
								Name:      "ns0",
								ProjectId: "pj0",
							},
							Project: &environmentproto.Project{
								Id: "pj0",
							},
							Role: accountproto.AccountV2_Role_Environment_EDITOR,
						},
						{
							Environment: &environmentproto.EnvironmentV2{
								Id:        "ns1",
								Name:      "ns1",
								ProjectId: "pj0",
							},
							Project: &environmentproto.Project{
								Id: "pj0",
							},
							Role: accountproto.AccountV2_Role_Environment_EDITOR,
						},
					},
					SearchFilters: []*accountproto.SearchFilter{
						{
							Id: "search-filter-id",
						},
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			p.ctx = metadata.NewIncomingContext(p.ctx, metadata.MD{
				"accept-language": []string{lang},
			})
			actual, err := service.GetMe(p.ctx, p.input)
			assert.Equal(t, p.expectedErr, err, p.desc)
			if actual != nil {
				assert.Greater(t, actual.Account.LastSeen, int64(0), "LastSeen should be set")
				assert.Less(t, actual.Account.LastSeen, time.Now().Unix()+1, "LastSeen should not be in the future")

				actualAccount := actual.Account
				actualAccount.LastSeen = 0
				expectedAccount := p.expected.Account
				expectedAccount.LastSeen = 0
				assert.Equal(t, expectedAccount, actualAccount, p.desc)
			}
		})
	}
}

func TestGetMyOrganizationsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithDefaultToken(t, true)
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
		input       *accountproto.GetMyOrganizationsRequest
		expected    *accountproto.GetMyOrganizationsResponse
		expectedErr error
	}{
		{
			desc: "errInternal: GetAccountsWithOrganization",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountsWithOrganization(
					gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("test"))
			},
			input:       &accountproto.GetMyOrganizationsRequest{},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "errInternal: GetOrganizations from environment service",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountsWithOrganization(
					gomock.Any(), gomock.Any(),
				).Return([]*domain.AccountWithOrganization{
					{
						Organization: &environmentproto.Organization{
							Id:          "org0",
							SystemAdmin: true,
						},
						AccountV2: &accountproto.AccountV2{
							Email: "bucketeer@example.com",
						},
					},
				}, nil)
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListOrganizations(
					gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("test"))
			},
			input:       &accountproto.GetMyOrganizationsRequest{},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountsWithOrganization(
					gomock.Any(), gomock.Any(),
				).Return([]*domain.AccountWithOrganization{}, nil)
			},
			input:       &accountproto.GetMyOrganizationsRequest{},
			expected:    &accountproto.GetMyOrganizationsResponse{Organizations: []*environmentproto.Organization{}},
			expectedErr: nil,
		},
		{
			desc: "success: including system admin organization",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountsWithOrganization(
					gomock.Any(), gomock.Any(),
				).Return([]*domain.AccountWithOrganization{
					{
						Organization: &environmentproto.Organization{
							Id:          "org0",
							SystemAdmin: true,
						},
						AccountV2: &accountproto.AccountV2{
							Email: "bucketeer@example.com",
						},
					},
				}, nil)
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListOrganizations(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListOrganizationsResponse{
					Organizations: []*environmentproto.Organization{
						{
							Id:          "org0",
							SystemAdmin: true,
						},
					},
				}, nil)
			},
			input: &accountproto.GetMyOrganizationsRequest{},
			expected: &accountproto.GetMyOrganizationsResponse{Organizations: []*environmentproto.Organization{
				{
					Id:          "org0",
					SystemAdmin: true,
				},
			}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			actual, err := service.GetMyOrganizations(ctx, p.input)
			assert.Equal(t, p.expectedErr, err, p.desc)
			assert.Equal(t, p.expected, actual, p.desc)
		})
	}
}

func TestGetMyOrganizationsByEmailMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithDefaultToken(t, true)
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
		input       *accountproto.GetMyOrganizationsByEmailRequest
		expected    *accountproto.GetMyOrganizationsResponse
		expectedErr error
	}{
		{
			desc:        "errBadRequest: Invalid email format",
			input:       &accountproto.GetMyOrganizationsByEmailRequest{Email: "bucketeer"},
			expected:    nil,
			expectedErr: createError(statusInvalidEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errInternal: GetAccountsWithOrganization",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountsWithOrganization(
					gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("test"))
			},
			input:       &accountproto.GetMyOrganizationsByEmailRequest{Email: "bucketeer@example.com"},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "errInternal: GetOrganizations from environment service",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountsWithOrganization(
					gomock.Any(), gomock.Any(),
				).Return([]*domain.AccountWithOrganization{
					{
						Organization: &environmentproto.Organization{
							Id:          "org0",
							SystemAdmin: true,
						},
						AccountV2: &accountproto.AccountV2{
							Email: "bucketeer@example.com",
						},
					},
				}, nil)
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListOrganizations(
					gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("test"))
			},
			input:       &accountproto.GetMyOrganizationsByEmailRequest{Email: "bucketeer@example.com"},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountsWithOrganization(
					gomock.Any(), gomock.Any(),
				).Return([]*domain.AccountWithOrganization{}, nil)
			},
			input:       &accountproto.GetMyOrganizationsByEmailRequest{Email: "bucketeer@example.com"},
			expected:    &accountproto.GetMyOrganizationsResponse{Organizations: []*environmentproto.Organization{}},
			expectedErr: nil,
		},
		{
			desc: "success: including system admin organization",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountsWithOrganization(
					gomock.Any(), gomock.Any(),
				).Return([]*domain.AccountWithOrganization{
					{
						Organization: &environmentproto.Organization{
							Id:          "org0",
							SystemAdmin: true,
						},
						AccountV2: &accountproto.AccountV2{
							Email: "bucketeer@example.com",
						},
					},
				}, nil)
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListOrganizations(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListOrganizationsResponse{
					Organizations: []*environmentproto.Organization{
						{
							Id:          "org0",
							SystemAdmin: true,
						},
					},
				}, nil)
			},
			input: &accountproto.GetMyOrganizationsByEmailRequest{Email: "bucketeer@example.com"},
			expected: &accountproto.GetMyOrganizationsResponse{Organizations: []*environmentproto.Organization{
				{
					Id:          "org0",
					SystemAdmin: true,
				},
			}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			actual, err := service.GetMyOrganizationsByEmail(ctx, p.input)
			assert.Equal(t, p.expectedErr, err, p.desc)
			assert.Equal(t, p.expected, actual, p.desc)
		})
	}
}

func getProjects(t *testing.T) []*environmentproto.Project {
	t.Helper()
	return []*environmentproto.Project{
		{Id: "pj0"},
	}
}

func getEnvironments(t *testing.T) []*environmentproto.EnvironmentV2 {
	t.Helper()
	return []*environmentproto.EnvironmentV2{
		{Id: "ns0", Name: "ns0", ProjectId: "pj0"},
		{Id: "ns1", Name: "ns1", ProjectId: "pj0"},
	}
}
