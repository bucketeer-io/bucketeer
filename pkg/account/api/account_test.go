// Copyright 2023 The Bucketeer Authors.
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

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	accstoragemock "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2/mock"
	ecmock "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestCreateAccountMySQL(t *testing.T) {
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
	envs := []*environmentproto.EnvironmentV2{
		{
			Id:             "ns0",
			Name:           "env0",
			OrganizationId: "org0",
		},
	}

	patterns := []struct {
		desc        string
		setup       func(*AccountService)
		ctxRole     accountproto.Account_Role
		req         *accountproto.CreateAccountRequest
		expectedErr error
	}{
		{
			desc:    "errNoCommand",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAccountRequest{
				Command:              nil,
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:    "errInvalidIsEmpty",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAccountRequest{
				Command:              &accountproto.CreateAccountCommand{Email: ""},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusEmailIsEmpty, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc:    "errInvalidEmail",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAccountRequest{
				Command:              &accountproto.CreateAccountCommand{Email: "bucketeer@"},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusInvalidEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errAlreadyExists_AdminAccount",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: envs,
					Cursor:       "1",
					TotalCount:   1,
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAdminAccount(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Account{
					Account: &accountproto.Account{
						Id:    "bucketeer@example.com",
						Email: "bucketeer@example.com",
						Name:  "test",
						Role:  accountproto.Account_OWNER,
					},
				}, nil)

			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAccountRequest{
				Command:              &accountproto.CreateAccountCommand{Email: "bucketeer_admin@example.com"},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "errAlreadyExists_EnvironmentAccount",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: envs,
					Cursor:       "1",
					TotalCount:   1,
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAdminAccount(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrAdminAccountNotFound)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountAlreadyExists)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAccountRequest{
				Command:              &accountproto.CreateAccountCommand{Email: "bucketeer_environment@example.com"},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: envs,
					Cursor:       "1",
					TotalCount:   1,
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAdminAccount(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrAdminAccountNotFound)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAccountRequest{
				Command: &accountproto.CreateAccountCommand{
					Email: "bucketeer@example.com",
					Role:  accountproto.Account_OWNER,
				},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: envs,
					Cursor:       "1",
					TotalCount:   1,
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAdminAccount(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrAdminAccountNotFound)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAccountRequest{
				Command: &accountproto.CreateAccountCommand{
					Email: "bucketeer@example.com",
					Role:  accountproto.Account_OWNER,
				},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, p.ctxRole)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.CreateAccount(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestChangeAccountRoleMySQL(t *testing.T) {
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
	envs := []*environmentproto.EnvironmentV2{
		{
			Id:             "ns0",
			Name:           "env0",
			OrganizationId: "org0",
		},
	}

	patterns := []struct {
		desc        string
		setup       func(*AccountService)
		ctxRole     accountproto.Account_Role
		req         *accountproto.ChangeAccountRoleRequest
		expectedErr error
	}{
		{
			desc:    "errMissingAccountID",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.ChangeAccountRoleRequest{
				Id:                   "",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusMissingAccountID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "account_id")),
		},
		{
			desc:    "errNoCommand",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.ChangeAccountRoleRequest{
				Id:                   "id",
				Command:              nil,
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "errNotFound",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: envs,
					Cursor:       "1",
					TotalCount:   1,
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.ChangeAccountRoleRequest{
				Id: "id",
				Command: &accountproto.ChangeAccountRoleCommand{
					Role: accountproto.Account_VIEWER,
				},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: envs,
					Cursor:       "1",
					TotalCount:   1,
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.ChangeAccountRoleRequest{
				Id: "bucketeer@example.com",
				Command: &accountproto.ChangeAccountRoleCommand{
					Role: accountproto.Account_VIEWER,
				},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: envs,
					Cursor:       "1",
					TotalCount:   1,
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.ChangeAccountRoleRequest{
				Id: "bucketeer@example.com",
				Command: &accountproto.ChangeAccountRoleCommand{
					Role: accountproto.Account_VIEWER,
				},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, p.ctxRole)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.ChangeAccountRole(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestEnableAccountMySQL(t *testing.T) {
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
	envs := []*environmentproto.EnvironmentV2{
		{
			Id:             "ns0",
			Name:           "env0",
			OrganizationId: "org0",
		},
	}

	patterns := []struct {
		desc        string
		setup       func(*AccountService)
		ctxRole     accountproto.Account_Role
		req         *accountproto.EnableAccountRequest
		expectedErr error
	}{
		{
			desc:    "errMissingAccountID",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.EnableAccountRequest{
				Id:                   "",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusMissingAccountID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "account_id")),
		},
		{
			desc:    "errNoCommand",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.EnableAccountRequest{
				Id:                   "id",
				Command:              nil,
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "errNotFound",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: envs,
					Cursor:       "1",
					TotalCount:   1,
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
				//s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccount(
				//	gomock.Any(), gomock.Any(), gomock.Any(),
				//).Return(nil, v2as.ErrAccountNotFound)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.EnableAccountRequest{
				Id:                   "id",
				Command:              &accountproto.EnableAccountCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: envs,
					Cursor:       "1",
					TotalCount:   1,
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.EnableAccountRequest{
				Id:                   "bucketeer@example.com",
				Command:              &accountproto.EnableAccountCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: envs,
					Cursor:       "1",
					TotalCount:   1,
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.EnableAccountRequest{
				Id:                   "bucketeer@example.com",
				Command:              &accountproto.EnableAccountCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx := setToken(ctx, p.ctxRole)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.EnableAccount(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestDisableAccountMySQL(t *testing.T) {
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
	envs := []*environmentproto.EnvironmentV2{
		{
			Id:             "ns0",
			Name:           "env0",
			OrganizationId: "org0",
		},
	}

	patterns := []struct {
		desc        string
		setup       func(*AccountService)
		ctxRole     accountproto.Account_Role
		req         *accountproto.DisableAccountRequest
		expectedErr error
	}{
		{
			desc:    "errMissingAccountID",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.DisableAccountRequest{
				Id:                   "",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusMissingAccountID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "account_id")),
		},
		{
			desc:    "errNoCommand",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.DisableAccountRequest{
				Id:                   "id",
				Command:              nil,
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "errNotFound",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: envs,
					Cursor:       "1",
					TotalCount:   1,
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.DisableAccountRequest{
				Id:                   "id",
				Command:              &accountproto.DisableAccountCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: envs,
					Cursor:       "1",
					TotalCount:   1,
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.DisableAccountRequest{
				Id:                   "bucketeer@example.com",
				Command:              &accountproto.DisableAccountCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: envs,
					Cursor:       "1",
					TotalCount:   1,
				}, nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.DisableAccountRequest{
				Id:                   "bucketeer@example.com",
				Command:              &accountproto.DisableAccountCommand{},
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, p.ctxRole)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.DisableAccount(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestGetAccountMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithDefaultToken(t, accountproto.Account_OWNER)
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
		req         *accountproto.GetAccountRequest
		expectedErr error
	}{
		{
			desc: "errMissingAccountID",
			req: &accountproto.GetAccountRequest{
				Email:                "",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusEmailIsEmpty, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errInvalidEmail",
			req: &accountproto.GetAccountRequest{
				Email:                "bucketeer@",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusInvalidEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errNotFound",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccount(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrAccountNotFound)
			},
			req: &accountproto.GetAccountRequest{
				Email:                "service@example.com",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccount(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Account{
					Account: &accountproto.Account{
						Id:    "bucketeer@example.com",
						Email: "bucketeer@example.com",
						Name:  "test",
						Role:  accountproto.Account_OWNER,
					},
				}, nil)
			},
			req: &accountproto.GetAccountRequest{
				Email:                "bucketeer@example.com",
				EnvironmentNamespace: "ns0",
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
			res, err := service.GetAccount(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, res)
			}
		})
	}
}

func TestListAccountsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithDefaultToken(t, accountproto.Account_OWNER)
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
		input       *accountproto.ListAccountsRequest
		expected    *accountproto.ListAccountsResponse
		expectedErr error
	}{
		{
			desc:        "errInvalidCursor",
			setup:       nil,
			input:       &accountproto.ListAccountsRequest{EnvironmentNamespace: "ns0", Cursor: "XXX"},
			expected:    nil,
			expectedErr: createError(statusInvalidCursor, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor")),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().ListAccounts(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), errors.New("test"))
			},
			input:       &accountproto.ListAccountsRequest{EnvironmentNamespace: "ns0"},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().ListAccounts(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*accountproto.Account{}, 0, int64(0), nil)
			},
			input:       &accountproto.ListAccountsRequest{PageSize: 2, Cursor: "", EnvironmentNamespace: "ns0"},
			expected:    &accountproto.ListAccountsResponse{Accounts: []*accountproto.Account{}, Cursor: "0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			actual, err := service.ListAccounts(ctx, p.input)
			assert.Equal(t, p.expectedErr, err, p.desc)
			assert.Equal(t, p.expected, actual, p.desc)
		})
	}
}

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
			desc: "errNoCommand",
			req: &accountproto.CreateAccountV2Request{
				Command:        nil,
				OrganizationId: "org0",
			},
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
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
						Name:             "test",
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
						Name:             "test",
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountAlreadyExists)
			},
			req: &accountproto.CreateAccountV2Request{
				Command: &accountproto.CreateAccountV2Command{
					Email:            "bucketeer_environment@example.com",
					Name:             "name",
					OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			req: &accountproto.CreateAccountV2Request{
				Command: &accountproto.CreateAccountV2Command{
					Email:            "bucketeer@example.com",
					Name:             "name",
					OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &accountproto.CreateAccountV2Request{
				Command: &accountproto.CreateAccountV2Command{
					Email:            "bucketeer@example.com",
					Name:             "name",
					OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
				},
				OrganizationId: "org0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, accountproto.Account_UNASSIGNED)
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
						Name:             "test",
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@",
				OrganizationId: "org0",
				ChangeNameCommand: &accountproto.ChangeAccountV2NameCommand{
					Name: "newName",
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email: "bucketeer@example.com",
				ChangeNameCommand: &accountproto.ChangeAccountV2NameCommand{
					Name: "newName",
				},
			},
			expectedErr: createError(statusMissingOrganizationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id")),
		},
		{
			desc: "errNoCommand",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "errInvalidNewName",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				ChangeNameCommand: &accountproto.ChangeAccountV2NameCommand{
					Name: strings.Repeat("a", 251),
				},
			},
			expectedErr: createError(statusInvalidName, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name")),
		},
		{
			desc: "errInvalidNewOrganizationRole",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						Name:             "test",
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				ChangeNameCommand: &accountproto.ChangeAccountV2NameCommand{
					Name: "newName",
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				ChangeNameCommand: &accountproto.ChangeAccountV2NameCommand{
					Name: "newName",
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &accountproto.UpdateAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				ChangeNameCommand: &accountproto.ChangeAccountV2NameCommand{
					Name: "newName",
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, accountproto.Account_UNASSIGNED)
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
						Name:             "test",
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
						Name:             "test",
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
						Name:             "test",
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
			desc: "errNoCommand",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.EnableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "errAccountNotFound",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			req: &accountproto.EnableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command:        &accountproto.EnableAccountV2Command{},
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			req: &accountproto.EnableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command:        &accountproto.EnableAccountV2Command{},
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &accountproto.EnableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command:        &accountproto.EnableAccountV2Command{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, accountproto.Account_UNASSIGNED)
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
						Name:             "test",
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
						Name:             "test",
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
						Name:             "test",
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
			desc: "errNoCommand",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.DisableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "errAccountNotFound",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			req: &accountproto.DisableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command:        &accountproto.DisableAccountV2Command{},
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			req: &accountproto.DisableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command:        &accountproto.DisableAccountV2Command{},
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &accountproto.DisableAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command:        &accountproto.DisableAccountV2Command{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, accountproto.Account_UNASSIGNED)
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
						Name:             "test",
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
						Name:             "test",
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
						Name:             "test",
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
			desc: "errNoCommand",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)
			},
			req: &accountproto.DeleteAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "errAccountNotFound",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			req: &accountproto.DeleteAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command:        &accountproto.DeleteAccountV2Command{},
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			req: &accountproto.DeleteAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command:        &accountproto.DeleteAccountV2Command{},
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &accountproto.DeleteAccountV2Request{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command:        &accountproto.DeleteAccountV2Command{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, accountproto.Account_UNASSIGNED)
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
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						Name:             "test",
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
						Name:             "test",
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
						Name:             "test",
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), "bucketeer@example.com", gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "bucketeer@example.com",
						Name:             "test",
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
			ctx = setToken(ctx, accountproto.Account_UNASSIGNED)
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

	ctx := createContextWithDefaultToken(t, accountproto.Account_UNASSIGNED)
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
						Name:             "test",
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().ListAccountsV2(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), errors.New("test"))
			},
			input:       &accountproto.ListAccountsV2Request{OrganizationId: "org0"},
			expected:    nil,
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
						Name:             "test",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
				}, nil)

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().ListAccountsV2(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
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

func setToken(ctx context.Context, role accountproto.Account_Role) context.Context {
	t := &token.IDToken{
		Issuer:    "issuer",
		Subject:   "sub",
		Audience:  "audience",
		Expiry:    time.Now().AddDate(100, 0, 0),
		IssuedAt:  time.Now(),
		Email:     "email",
		AdminRole: role,
	}
	return context.WithValue(ctx, rpc.Key, t)
}
