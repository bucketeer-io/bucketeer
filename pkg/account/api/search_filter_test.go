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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
			},
			req: &accountproto.CreateSearchFilterRequest{
				OrganizationId:       "org0",
				EnvironmentNamespace: "envName0",
				Command: &accountproto.CreateSearchFilterCommand{
					SearchFilter: &accountproto.SearchFilter{
						Name:             "filter",
						Query:            "query",
						FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
						EnvironmentId:    "envID0",
						DefaultFilter:    false,
					},
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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:                "bucketeer@example.com",
				EnvironmentNamespace: "envName0",
				Command: &accountproto.CreateSearchFilterCommand{
					SearchFilter: &accountproto.SearchFilter{
						Name:             "filter",
						Query:            "query",
						FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
						EnvironmentId:    "envID0",
						DefaultFilter:    false,
					},
				},
			},
			expectedErr: createError(statusMissingOrganizationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id")),
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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:                "bucketeer@example.com",
				OrganizationId:       "org0",
				EnvironmentNamespace: "envName0",
				Command: &accountproto.CreateSearchFilterCommand{
					SearchFilter: &accountproto.SearchFilter{
						Name:             "filter",
						Query:            "query",
						FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
						EnvironmentId:    "envID0",
						DefaultFilter:    false,
					},
				},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalizeWithTemplate(locale.InternalServerError)),
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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:                "bucketeer@example.com",
				OrganizationId:       "org0",
				EnvironmentNamespace: "envName0",
				Command: &accountproto.CreateSearchFilterCommand{
					SearchFilter: &accountproto.SearchFilter{
						Name:             "filter",
						Query:            "query",
						FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
						EnvironmentId:    "envID0",
						DefaultFilter:    false,
					},
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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command:        nil,
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "err: SearchFilter is nil",
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
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command: &accountproto.CreateSearchFilterCommand{
					SearchFilter: nil,
				},
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "search_filter")),
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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command: &accountproto.CreateSearchFilterCommand{
					SearchFilter: &accountproto.SearchFilter{
						Name: "",
					},
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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command: &accountproto.CreateSearchFilterCommand{
					SearchFilter: &accountproto.SearchFilter{
						Name:  "name",
						Query: "",
					},
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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command: &accountproto.CreateSearchFilterCommand{
					SearchFilter: &accountproto.SearchFilter{
						Name:             "name",
						Query:            "query",
						FilterTargetType: accountproto.FilterTargetType_UNKNOWN,
					},
				},
			},
			expectedErr: createError(statusSearchFilterTargetTypeIsUnknown, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "filter_target_type")),
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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &accountproto.CreateSearchFilterRequest{
				Email:                "bucketeer@example.com",
				OrganizationId:       "org0",
				EnvironmentNamespace: "envName0",
				Command: &accountproto.CreateSearchFilterCommand{
					SearchFilter: &accountproto.SearchFilter{
						Name:             "filter",
						Query:            "query",
						FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
						EnvironmentId:    "envID0",
						DefaultFilter:    false,
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
			_, err := service.CreateSearchFilterV2(ctx, p.req)
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
		req         *accountproto.DeleteSearchFilterRequest
		expectedErr error
	}{
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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
			},
			req: &accountproto.DeleteSearchFilterRequest{
				OrganizationId:       "org0",
				EnvironmentNamespace: "envName0",
				Command: &accountproto.DeleteSearchFilterCommand{
					SearchFilterId: "filterID",
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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
			},
			req: &accountproto.DeleteSearchFilterRequest{
				Email:                "bucketeer@example.com",
				EnvironmentNamespace: "envName0",
				Command: &accountproto.DeleteSearchFilterCommand{
					SearchFilterId: "filterID",
				},
			},
			expectedErr: createError(statusMissingOrganizationID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id")),
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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("test"))
			},
			req: &accountproto.DeleteSearchFilterRequest{
				Email:                "bucketeer@example.com",
				OrganizationId:       "org0",
				EnvironmentNamespace: "envName0",
				Command: &accountproto.DeleteSearchFilterCommand{
					SearchFilterId: "filterID",
				},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalizeWithTemplate(locale.InternalServerError)),
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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(v2as.ErrAccountNotFound)
			},
			req: &accountproto.DeleteSearchFilterRequest{
				Email:                "bucketeer@example.com",
				OrganizationId:       "org0",
				EnvironmentNamespace: "envName0",
				Command: &accountproto.DeleteSearchFilterCommand{
					SearchFilterId: "filterID",
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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
			},
			req: &accountproto.DeleteSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command:        nil,
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
			},
			req: &accountproto.DeleteSearchFilterRequest{
				Email:          "bucketeer@example.com",
				OrganizationId: "org0",
				Command: &accountproto.DeleteSearchFilterCommand{
					SearchFilterId: "",
				},
			},
			expectedErr: createError(statusSearchFilterIDIsEmpty, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "search_filter_id")),
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
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil).AnyTimes()
				s.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &accountproto.DeleteSearchFilterRequest{
				Email:                "bucketeer@example.com",
				OrganizationId:       "org0",
				EnvironmentNamespace: "envName0",
				Command: &accountproto.DeleteSearchFilterCommand{
					SearchFilterId: "filterID",
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
			_, err := service.DeleteSearchFilterV2(ctx, p.req)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}

}
