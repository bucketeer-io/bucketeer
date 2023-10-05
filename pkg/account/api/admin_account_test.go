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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	v2as "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	ecmock "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	storagemock "github.com/bucketeer-io/bucketeer/pkg/storage/mock"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestGetMeV2MySQL(t *testing.T) {
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

	patterns := []struct {
		desc            string
		ctx             context.Context
		setup           func(*AccountService)
		input           *accountproto.GetMeV2Request
		expected        string
		expectedIsAdmin bool
		expectedErr     error
	}{
		{
			desc:        "errUnauthenticated",
			ctx:         context.Background(),
			setup:       nil,
			input:       &accountproto.GetMeV2Request{},
			expected:    "",
			expectedErr: createError(statusUnauthenticated, localizer.MustLocalize(locale.UnauthenticatedError)),
		},
		{
			desc:        "errInvalidEmail",
			ctx:         createContextWithInvalidEmailToken(t, accountproto.Account_OWNER),
			setup:       nil,
			input:       &accountproto.GetMeV2Request{},
			expected:    "",
			expectedErr: createError(statusInvalidEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errInternal",
			ctx:  createContextWithDefaultToken(t, accountproto.Account_OWNER),
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListProjects(
					gomock.Any(),
					gomock.Any(),
				).Return(
					nil,
					createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
				)
			},
			input:       &accountproto.GetMeV2Request{},
			expected:    "",
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "errInternal_no_projects",
			ctx:  createContextWithDefaultToken(t, accountproto.Account_OWNER),
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListProjects(
					gomock.Any(),
					gomock.Any(),
				).Return(
					&environmentproto.ListProjectsResponse{},
					nil,
				)
			},
			input:       &accountproto.GetMeV2Request{},
			expected:    "",
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "errInternal_no_environments",
			ctx:  createContextWithDefaultToken(t, accountproto.Account_OWNER),
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
					&environmentproto.ListEnvironmentsV2Response{},
					nil,
				)
			},
			input:       &accountproto.GetMeV2Request{},
			expected:    "",
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "errNotFound",
			ctx:  createContextWithDefaultToken(t, accountproto.Account_EDITOR),
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
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows).Times(3)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row).Times(3)
			},
			input:       &accountproto.GetMeV2Request{},
			expected:    "",
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
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
			actual, err := service.GetMeV2(p.ctx, p.input)
			assert.Equal(t, p.expectedErr, err, p.desc)
			if actual != nil {
				assert.Equal(t, p.expected, actual.Email, p.desc)
				assert.Equal(t, p.expectedIsAdmin, actual.IsAdmin, p.desc)
			}
		})
	}
}

func TestGetMeByEmailV2MySQL(t *testing.T) {
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

	patterns := []struct {
		desc            string
		ctx             context.Context
		setup           func(*AccountService)
		input           *accountproto.GetMeByEmailV2Request
		expected        string
		expectedIsAdmin bool
		expectedErr     error
	}{
		{
			desc:  "errInvalidEmail",
			ctx:   createContextWithDefaultToken(t, accountproto.Account_OWNER),
			setup: nil,
			input: &accountproto.GetMeByEmailV2Request{
				Email: "bucketeer@",
			},
			expected: "",
			expectedErr: createError(
				statusInvalidEmail,
				localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
			),
		},
		{
			desc: "errNotFound",
			ctx:  createContextWithDefaultToken(t, accountproto.Account_EDITOR),
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
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows).Times(3)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row).Times(3)
			},
			input: &accountproto.GetMeByEmailV2Request{
				Email: "bucketeer@example.com",
			},
			expected:    "",
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
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
			actual, err := service.GetMeByEmailV2(p.ctx, p.input)
			assert.Equal(t, p.expectedErr, err, p.desc)
			if actual != nil {
				assert.Equal(t, p.expected, actual.Email, p.desc)
				assert.Equal(t, p.expectedIsAdmin, actual.IsAdmin, p.desc)
			}
		})
	}
}

func TestCreateAdminAccountMySQL(t *testing.T) {
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
		ctxRole     accountproto.Account_Role
		req         *accountproto.CreateAdminAccountRequest
		expectedErr error
	}{
		{
			desc:    "errNoCommand",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAdminAccountRequest{
				Command: nil,
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:    "errEmailIsEmpty",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAdminAccountRequest{
				Command: &accountproto.CreateAdminAccountCommand{Email: ""},
			},
			expectedErr: createError(statusEmailIsEmpty, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc:    "errInvalidEmail",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAdminAccountRequest{
				Command: &accountproto.CreateAdminAccountCommand{Email: "bucketeer@"},
			},
			expectedErr: createError(statusInvalidEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(nil, createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)))
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAdminAccountRequest{
				Command: &accountproto.CreateAdminAccountCommand{
					Email: "bucketeer@example.com",
				},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "errAlreadyExists_EnvironmentAccount",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: getEnvironments(t),
					Cursor:       "",
				}, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAdminAccountRequest{
				Command: &accountproto.CreateAdminAccountCommand{
					Email: "bucketeer_environment@example.com",
				},
			},
			expectedErr: createError(statusAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "errAlreadyExists_AdminAccount",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: getEnvironments(t),
					Cursor:       "",
				}, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows).Times(2)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row).Times(2)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAdminAccountAlreadyExists)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAdminAccountRequest{
				Command: &accountproto.CreateAdminAccountCommand{
					Email: "bucketeer_admin@example.com",
				},
			},
			expectedErr: createError(statusAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: getEnvironments(t),
					Cursor:       "",
				}, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows).Times(2)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row).Times(2)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.CreateAdminAccountRequest{
				Command: &accountproto.CreateAdminAccountCommand{
					Email: "bucketeer@example.com",
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, p.ctxRole)
			service := createAccountService(t, mockController, storagemock.NewMockClient(mockController))
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.CreateAdminAccount(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestEnableAdminAccountMySQL(t *testing.T) {
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
		ctxRole     accountproto.Account_Role
		req         *accountproto.EnableAdminAccountRequest
		expectedErr error
	}{
		{
			desc:    "errMissingAccountID",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.EnableAdminAccountRequest{
				Id: "",
			},
			expectedErr: createError(statusMissingAccountID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "account_id")),
		},
		{
			desc:    "errNoCommand",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.EnableAdminAccountRequest{
				Id:      "id",
				Command: nil,
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "errNotFound",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAdminAccountNotFound)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.EnableAdminAccountRequest{
				Id:      "id",
				Command: &accountproto.EnableAdminAccountCommand{},
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.EnableAdminAccountRequest{
				Id:      "bucketeer@example.com",
				Command: &accountproto.EnableAdminAccountCommand{},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.EnableAdminAccountRequest{
				Id:      "bucketeer@example.com",
				Command: &accountproto.EnableAdminAccountCommand{},
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
			_, err := service.EnableAdminAccount(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestDisableAdminAccountMySQL(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	t.Parallel()

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
		ctxRole     accountproto.Account_Role
		req         *accountproto.DisableAdminAccountRequest
		expectedErr error
	}{
		{
			desc:    "errMissingAccountID",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.DisableAdminAccountRequest{
				Id: "",
			},
			expectedErr: createError(statusMissingAccountID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "account_id")),
		},
		{
			desc:    "errNoCommand",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.DisableAdminAccountRequest{
				Id:      "id",
				Command: nil,
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "errNotFound",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAdminAccountNotFound)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.DisableAdminAccountRequest{
				Id:      "id",
				Command: &accountproto.DisableAdminAccountCommand{},
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.DisableAdminAccountRequest{
				Id:      "bucketeer@example.com",
				Command: &accountproto.DisableAdminAccountCommand{},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.DisableAdminAccountRequest{
				Id:      "bucketeer@example.com",
				Command: &accountproto.DisableAdminAccountCommand{},
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
			_, err := service.DisableAdminAccount(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestConvertAccountMySQL(t *testing.T) {
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
		ctxRole     accountproto.Account_Role
		req         *accountproto.ConvertAccountRequest
		expectedErr error
	}{
		{
			desc:    "errMissingAccountID",
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.ConvertAccountRequest{
				Id:      "",
				Command: &accountproto.ConvertAccountCommand{},
			},
			expectedErr: createError(statusMissingAccountID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "account_id")),
		},
		{
			desc: "errNotFound",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(),
					gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: getEnvironments(t),
					Cursor:       "",
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.ConvertAccountRequest{
				Id:      "b@aa.jp",
				Command: &accountproto.ConvertAccountCommand{},
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.environmentClient.(*ecmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(),
					gomock.Any(),
				).Return(&environmentproto.ListEnvironmentsV2Response{
					Environments: getEnvironments(t),
					Cursor:       "",
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			ctxRole: accountproto.Account_OWNER,
			req: &accountproto.ConvertAccountRequest{
				Id:      "bucketeer@example.com",
				Command: &accountproto.ConvertAccountCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx := setToken(ctx, p.ctxRole)
			service := createAccountService(t, mockController, storagemock.NewMockClient(mockController))
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.ConvertAccount(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestGetAdminAccountMySQL(t *testing.T) {
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
		req         *accountproto.GetAdminAccountRequest
		expectedErr error
	}{
		{
			desc: "errMissingAccountID",
			req: &accountproto.GetAdminAccountRequest{
				Email: "",
			},
			expectedErr: createError(statusEmailIsEmpty, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errInvalidEmail",
			req: &accountproto.GetAdminAccountRequest{
				Email: "bucketeer@",
			},
			expectedErr: createError(statusInvalidEmail, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email")),
		},
		{
			desc: "errNotFound",
			setup: func(s *AccountService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req: &accountproto.GetAdminAccountRequest{
				Email: "service@example.com",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req: &accountproto.GetAdminAccountRequest{
				Email: "bucketeer@example.com",
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
			res, err := service.GetAdminAccount(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
			if res != nil {
				assert.NotNil(t, res)
			}
		})
	}
}

func TestListAdminAccountsMySQL(t *testing.T) {
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
		input       *accountproto.ListAdminAccountsRequest
		expected    *accountproto.ListAdminAccountsResponse
		expectedErr error
	}{
		{
			desc:        "errInvalidCursor",
			setup:       nil,
			input:       &accountproto.ListAdminAccountsRequest{Cursor: "xxx"},
			expected:    nil,
			expectedErr: createError(statusInvalidCursor, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor")),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("test"))
			},
			input:       &accountproto.ListAdminAccountsRequest{},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
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
			input:       &accountproto.ListAdminAccountsRequest{PageSize: 2, Cursor: ""},
			expected:    &accountproto.ListAdminAccountsResponse{Accounts: []*accountproto.Account{}, Cursor: "0", TotalCount: 0},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			actual, err := service.ListAdminAccounts(ctx, p.input)
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
