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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	accstoragemock "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
)

func TestCreateAPIKeyMySQL(t *testing.T) {
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
		desc          string
		setup         func(*AccountService)
		isSystemAdmin bool
		req           *accountproto.CreateAPIKeyRequest
		expectedErr   error
	}{
		{
			desc:          "errMissingAPIKeyName",
			isSystemAdmin: true,
			req: &accountproto.CreateAPIKeyRequest{
				Command: &accountproto.CreateAPIKeyCommand{Name: ""},
			},
			expectedErr: createError(statusMissingAPIKeyName, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_name")),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			isSystemAdmin: true,
			req: &accountproto.CreateAPIKeyRequest{
				Command: &accountproto.CreateAPIKeyCommand{
					Name: "name",
					Role: accountproto.APIKey_SDK_CLIENT,
				},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().CreateAPIKey(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			isSystemAdmin: true,
			req: &accountproto.CreateAPIKeyRequest{
				Command: &accountproto.CreateAPIKeyCommand{
					Name: "name",
					Role: accountproto.APIKey_SDK_CLIENT,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, p.isSystemAdmin)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.CreateAPIKey(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestCreateAPIKeyMySQLNoCommand(t *testing.T) {
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
		desc          string
		setup         func(*AccountService)
		isSystemAdmin bool
		req           *accountproto.CreateAPIKeyRequest
		expectedErr   error
	}{
		{
			desc:          "errMissingAPIKeyName",
			isSystemAdmin: true,
			req: &accountproto.CreateAPIKeyRequest{
				Name: "",
			},
			expectedErr: createError(statusMissingAPIKeyName, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_name")),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			isSystemAdmin: true,
			req: &accountproto.CreateAPIKeyRequest{
				Name: "name",
				Role: accountproto.APIKey_SDK_CLIENT,
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().CreateAPIKey(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)

				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			isSystemAdmin: true,
			req: &accountproto.CreateAPIKeyRequest{
				Name:        "name",
				Maintainer:  "bucketeer@bucketeer.io",
				Role:        accountproto.APIKey_SDK_CLIENT,
				Description: "test key",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, p.isSystemAdmin)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.CreateAPIKey(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestChangeAPIKeyNameMySQL(t *testing.T) {
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
		desc          string
		setup         func(*AccountService)
		isSystemAdmin bool
		req           *accountproto.ChangeAPIKeyNameRequest
		expectedErr   error
	}{
		{
			desc:          "errMissingAPIKeyID",
			isSystemAdmin: true,
			req: &accountproto.ChangeAPIKeyNameRequest{
				Id: "",
			},
			expectedErr: createError(statusMissingAPIKeyID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_id")),
		},
		{
			desc:          "errNoCommand",
			isSystemAdmin: true,
			req: &accountproto.ChangeAPIKeyNameRequest{
				Id:      "id",
				Command: nil,
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "errNotFound",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAPIKeyNotFound)
			},
			isSystemAdmin: true,
			req: &accountproto.ChangeAPIKeyNameRequest{
				Id: "id",
				Command: &accountproto.ChangeAPIKeyNameCommand{
					Name: "",
				},
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			isSystemAdmin: true,
			req: &accountproto.ChangeAPIKeyNameRequest{
				Id: "id",
				Command: &accountproto.ChangeAPIKeyNameCommand{
					Name: "new name",
				},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAPIKey(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.APIKey{
					APIKey: &accountproto.APIKey{},
				}, nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAPIKey(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			isSystemAdmin: true,
			req: &accountproto.ChangeAPIKeyNameRequest{
				Id: "id",
				Command: &accountproto.ChangeAPIKeyNameCommand{
					Name: "new name",
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, p.isSystemAdmin)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.ChangeAPIKeyName(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestEnableAPIKeyMySQL(t *testing.T) {
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
		desc          string
		setup         func(*AccountService)
		isSystemAdmin bool
		req           *accountproto.EnableAPIKeyRequest
		expectedErr   error
	}{
		{
			desc:          "errMissingAPIKeyID",
			isSystemAdmin: true,
			req: &accountproto.EnableAPIKeyRequest{
				Id: "",
			},
			expectedErr: createError(statusMissingAPIKeyID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_id")),
		},
		{
			desc:          "errNoCommand",
			isSystemAdmin: true,
			req: &accountproto.EnableAPIKeyRequest{
				Id:      "id",
				Command: nil,
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "errNotFound",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAPIKeyNotFound)
			},
			isSystemAdmin: true,
			req: &accountproto.EnableAPIKeyRequest{
				Id:      "id",
				Command: &accountproto.EnableAPIKeyCommand{},
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			isSystemAdmin: true,
			req: &accountproto.EnableAPIKeyRequest{
				Id:      "id",
				Command: &accountproto.EnableAPIKeyCommand{},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAPIKey(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.APIKey{
					APIKey: &accountproto.APIKey{},
				}, nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAPIKey(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			isSystemAdmin: true,
			req: &accountproto.EnableAPIKeyRequest{
				Id:      "id",
				Command: &accountproto.EnableAPIKeyCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx := setToken(ctx, p.isSystemAdmin)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.EnableAPIKey(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestDisableAPIKeyMySQL(t *testing.T) {
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
		desc          string
		setup         func(*AccountService)
		isSystemAdmin bool
		req           *accountproto.DisableAPIKeyRequest
		expectedErr   error
	}{
		{
			desc:          "errMissingAPIKeyID",
			isSystemAdmin: true,
			req: &accountproto.DisableAPIKeyRequest{
				Id: "",
			},
			expectedErr: createError(statusMissingAPIKeyID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_id")),
		},
		{
			desc:          "errNoCommand",
			isSystemAdmin: true,
			req: &accountproto.DisableAPIKeyRequest{
				Id:      "id",
				Command: nil,
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "errNotFound",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAPIKeyNotFound)
			},
			isSystemAdmin: true,
			req: &accountproto.DisableAPIKeyRequest{
				Id:      "id",
				Command: &accountproto.DisableAPIKeyCommand{},
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "errInternal",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			isSystemAdmin: true,
			req: &accountproto.DisableAPIKeyRequest{
				Id:      "id",
				Command: &accountproto.DisableAPIKeyCommand{},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AccountService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAPIKey(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.APIKey{
					APIKey: &accountproto.APIKey{},
				}, nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().UpdateAPIKey(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			isSystemAdmin: true,
			req: &accountproto.DisableAPIKeyRequest{
				Id:      "id",
				Command: &accountproto.DisableAPIKeyCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(ctx, p.isSystemAdmin)
			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.DisableAPIKey(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestGetAPIKeyMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	createError := func(localizer locale.Localizer, status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc           string
		context        context.Context
		setup          func(*AccountService)
		req            *accountproto.GetAPIKeyRequest
		getExpectedErr func(localizer locale.Localizer) error
	}{
		{
			desc:    "errMissingAPIKeyID",
			context: createContextWithDefaultToken(t, true),
			req:     &accountproto.GetAPIKeyRequest{Id: ""},
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(localizer, statusMissingAPIKeyID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_id"))
			},
		},
		{
			desc:    "errNotFound",
			context: createContextWithDefaultToken(t, true),
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAPIKey(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, v2as.ErrAPIKeyNotFound)
			},
			req: &accountproto.GetAPIKeyRequest{Id: "id"},
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(localizer, statusNotFound, localizer.MustLocalize(locale.NotFoundError))
			},
		},
		{
			desc:    "success",
			context: createContextWithDefaultToken(t, true),
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAPIKey(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.APIKey{
					APIKey: &accountproto.APIKey{
						Id: "id",
					},
				}, nil)
			},
			req: &accountproto.GetAPIKeyRequest{Id: "id"},
			getExpectedErr: func(localizer locale.Localizer) error {
				return nil
			},
		},
		{
			desc:    "success with viewer account",
			context: createContextWithDefaultToken(t, false),
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAPIKey(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.APIKey{
					APIKey: &accountproto.APIKey{
						Id: "id",
					},
				}, nil)
			},
			req: &accountproto.GetAPIKeyRequest{Id: "id", EnvironmentId: "ns0"},
			getExpectedErr: func(localizer locale.Localizer) error {
				return nil
			},
		},
		{
			desc:    "errPermissionDenied",
			context: createContextWithDefaultToken(t, false),
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2ByEnvironmentID(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_UNASSIGNED,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_UNASSIGNED,
							},
						},
					},
				}, nil).AnyTimes()
			},
			req: &accountproto.GetAPIKeyRequest{Id: "id", EnvironmentId: "ns0"},
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(localizer, statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied))
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx := p.context
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				"accept-language": []string{"ja"},
			})

			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			localizer := locale.NewLocalizer(ctx)
			res, err := service.GetAPIKey(ctx, p.req)
			assert.Equal(t, p.getExpectedErr(localizer), err)
			if err == nil {
				assert.NotNil(t, res)
			}
		})
	}
}

func TestListAPIKeysMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	createError := func(localizer locale.Localizer, status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc           string
		context        context.Context
		setup          func(*AccountService)
		input          *accountproto.ListAPIKeysRequest
		expected       *accountproto.ListAPIKeysResponse
		getExpectedErr func(localizer locale.Localizer) error
	}{
		{
			desc:    "errInvalidCursor",
			context: createContextWithDefaultToken(t, true),
			input: &accountproto.ListAPIKeysRequest{
				EnvironmentIds: []string{"ns0"},
				OrganizationId: "org0",
				Cursor:         "XXX",
			},
			expected: nil,
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(localizer, statusInvalidCursor, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor"))
			},
		},
		{
			desc:    "errInternal",
			context: createContextWithDefaultToken(t, true),
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().ListAPIKeys(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), errors.New("error"))
			},
			input: &accountproto.ListAPIKeysRequest{
				EnvironmentIds: []string{"ns0"},
				OrganizationId: "org0",
			},
			expected: nil,
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(localizer, statusInternal, localizer.MustLocalize(locale.InternalServerError))
			},
		},
		{
			desc:    "errPermissionDenied",
			context: createContextWithDefaultToken(t, false),
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_UNASSIGNED,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_UNASSIGNED,
							},
						},
					},
				}, nil).AnyTimes()
			},
			input: &accountproto.ListAPIKeysRequest{
				EnvironmentIds: []string{"ns0"},
				OrganizationId: "org0",
			},
			expected: nil,
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(localizer, statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied))
			},
		},
		{
			desc:    "success",
			context: createContextWithDefaultToken(t, true),
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().ListAPIKeys(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*accountproto.APIKey{}, 0, int64(0), nil)
			},
			input: &accountproto.ListAPIKeysRequest{
				OrganizationId: "org0",
				PageSize:       2,
				Cursor:         "",
			},
			expected: &accountproto.ListAPIKeysResponse{ApiKeys: []*accountproto.APIKey{}, Cursor: "0"},
			getExpectedErr: func(localizer locale.Localizer) error {
				return nil
			},
		},
		{
			desc:    "success with admin account",
			context: createContextWithDefaultToken(t, false),
			setup: func(s *AccountService) {
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AccountV2{
					AccountV2: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().ListAPIKeys(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*accountproto.APIKey{}, 0, int64(0), nil)
			},
			input: &accountproto.ListAPIKeysRequest{
				EnvironmentIds: []string{"ns0"},
				OrganizationId: "org0",
			},
			expected: &accountproto.ListAPIKeysResponse{ApiKeys: []*accountproto.APIKey{}, Cursor: "0"},
			getExpectedErr: func(localizer locale.Localizer) error {
				return nil
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx, cancel := context.WithCancel(p.context)
			defer cancel()
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				"accept-language": []string{"ja"},
			})
			localizer := locale.NewLocalizer(ctx)

			service := createAccountService(t, mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			actual, err := service.ListAPIKeys(ctx, p.input)
			assert.Equal(t, p.getExpectedErr(localizer), err, p.desc)
			assert.Equal(t, p.expected, actual, p.desc)
		})
	}
}
