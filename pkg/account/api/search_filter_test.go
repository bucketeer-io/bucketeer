// Copyright 2024 The Bucketeer Authors.
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
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
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
				}, nil).AnyTimes()
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
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
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
				}, nil).AnyTimes()
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
			expectedErr: createError(statusEmailIsEmpty, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email")),
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
				}, nil).AnyTimes()
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
			expectedErr: createError(statusMissingOrganizationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id")),
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
				).Return(&account, nil).AnyTimes()

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
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
			expectedErr: createError(statusInternal, localizer.MustLocalizeWithTemplate(locale.InternalServerError)),
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
				).Return(&account, nil).AnyTimes()

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
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
			expectedErr: createError(statusNotFound, localizer.MustLocalizeWithTemplate(locale.NotFoundError)),
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
				}, nil).AnyTimes()
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command:        nil,
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
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
				}, nil).AnyTimes()
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				EnvironmentId:  "envID0",
				Command: &accountproto.CreateSearchFilterCommand{
					Name: "",
				},
			},
			expectedErr: createError(statusSearchFilterNameIsEmpty, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
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
				}, nil).AnyTimes()
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
			expectedErr: createError(statusSearchFilterQueryIsEmpty, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "query")),
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
				}, nil).AnyTimes()
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
			expectedErr: createError(statusSearchFilterTargetTypeIsRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "filter_target_type")),
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

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(nil)
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

				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(nil)
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
