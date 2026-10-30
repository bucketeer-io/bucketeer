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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	accstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	dbmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

func TestSearchFilterRequestsTargetCallerAccount(t *testing.T) {
	t.Parallel()

	const (
		callerEmail    = "viewer@example.com"
		callerOrgID    = "org0"
		envID          = "env0"
		otherEmail     = "other@example.com"
		otherOrgID     = "org1"
		systemAdminOrg = "system-admin-org"
		filterID       = "filter0"
	)

	callerAccount := func() *domain.AccountV2 {
		return &domain.AccountV2{AccountV2: &accountproto.AccountV2{
			Email:            callerEmail,
			OrganizationId:   callerOrgID,
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{EnvironmentId: envID, Role: accountproto.AccountV2_Role_Environment_VIEWER},
			},
			SearchFilters: []*accountproto.SearchFilter{
				{
					Id:               filterID,
					Name:             "mine",
					Query:            "query",
					FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
					EnvironmentId:    envID,
				},
			},
		}}
	}
	systemAdminAccount := func() *domain.AccountV2 {
		account := callerAccount()
		account.OrganizationId = systemAdminOrg
		return account
	}
	createCommand := &accountproto.CreateSearchFilterCommand{
		Name:             "new",
		Query:            "new-query",
		FilterTargetType: accountproto.FilterTargetType_FEATURE_FLAG,
		EnvironmentId:    envID,
		DefaultFilter:    true,
	}

	cases := []struct {
		name        string
		systemAdmin bool
		call        func(context.Context, *AccountService) error
		wantErr     error
		wantOrgID   string
		check       func(*testing.T, *domain.AccountV2)
	}{
		{
			name: "create rejects another account in another organization",
			call: func(ctx context.Context, s *AccountService) error {
				_, err := s.CreateSearchFilter(ctx, &accountproto.CreateSearchFilterRequest{
					Email:          otherEmail,
					OrganizationId: otherOrgID,
					EnvironmentId:  envID,
					Command:        createCommand,
				})
				return err
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name: "create rejects another account in the same organization",
			call: func(ctx context.Context, s *AccountService) error {
				_, err := s.CreateSearchFilter(ctx, &accountproto.CreateSearchFilterRequest{
					Email:          otherEmail,
					OrganizationId: callerOrgID,
					EnvironmentId:  envID,
					Command:        createCommand,
				})
				return err
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name: "update rejects another account",
			call: func(ctx context.Context, s *AccountService) error {
				_, err := s.UpdateSearchFilter(ctx, &accountproto.UpdateSearchFilterRequest{
					Email:          otherEmail,
					OrganizationId: otherOrgID,
					EnvironmentId:  envID,
					ChangeNameCommand: &accountproto.ChangeSearchFilterNameCommand{
						Id:   filterID,
						Name: "renamed",
					},
				})
				return err
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name: "delete rejects another account",
			call: func(ctx context.Context, s *AccountService) error {
				_, err := s.DeleteSearchFilter(ctx, &accountproto.DeleteSearchFilterRequest{
					Email:          otherEmail,
					OrganizationId: otherOrgID,
					EnvironmentId:  envID,
					Command:        &accountproto.DeleteSearchFilterCommand{Id: filterID},
				})
				return err
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name: "create stores the filter on the caller's account in the environment's organization",
			call: func(ctx context.Context, s *AccountService) error {
				_, err := s.CreateSearchFilter(ctx, &accountproto.CreateSearchFilterRequest{
					Email:          callerEmail,
					OrganizationId: otherOrgID,
					EnvironmentId:  envID,
					Command:        createCommand,
				})
				return err
			},
			wantOrgID: callerOrgID,
			check: func(t *testing.T, persisted *domain.AccountV2) {
				require.Len(t, persisted.SearchFilters, 2)
				assert.Equal(t, "new", persisted.SearchFilters[1].Name)
				assert.True(t, persisted.SearchFilters[1].DefaultFilter)
			},
		},
		{
			name: "update changes the caller's own filter",
			call: func(ctx context.Context, s *AccountService) error {
				_, err := s.UpdateSearchFilter(ctx, &accountproto.UpdateSearchFilterRequest{
					Email:          callerEmail,
					OrganizationId: callerOrgID,
					EnvironmentId:  envID,
					ChangeNameCommand: &accountproto.ChangeSearchFilterNameCommand{
						Id:   filterID,
						Name: "renamed",
					},
				})
				return err
			},
			wantOrgID: callerOrgID,
			check: func(t *testing.T, persisted *domain.AccountV2) {
				require.Len(t, persisted.SearchFilters, 1)
				assert.Equal(t, "renamed", persisted.SearchFilters[0].Name)
			},
		},
		{
			name: "delete removes the caller's own filter",
			call: func(ctx context.Context, s *AccountService) error {
				_, err := s.DeleteSearchFilter(ctx, &accountproto.DeleteSearchFilterRequest{
					Email:          callerEmail,
					OrganizationId: callerOrgID,
					EnvironmentId:  envID,
					Command:        &accountproto.DeleteSearchFilterCommand{Id: filterID},
				})
				return err
			},
			wantOrgID: callerOrgID,
			check: func(t *testing.T, persisted *domain.AccountV2) {
				assert.Empty(t, persisted.SearchFilters)
			},
		},
		{
			name:        "system admin stores the filter on the system admin account",
			systemAdmin: true,
			call: func(ctx context.Context, s *AccountService) error {
				_, err := s.CreateSearchFilter(ctx, &accountproto.CreateSearchFilterRequest{
					Email:          callerEmail,
					OrganizationId: callerOrgID,
					EnvironmentId:  envID,
					Command:        createCommand,
				})
				return err
			},
			wantOrgID: systemAdminOrg,
			check: func(t *testing.T, persisted *domain.AccountV2) {
				require.Len(t, persisted.SearchFilters, 2)
				assert.Equal(t, "new", persisted.SearchFilters[1].Name)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			ctx := createContextWithEmailToken(t, callerEmail, tc.systemAdmin)
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{})

			service := createAccountService(t, mockController)

			// Only the caller's own account is ever loaded. Looking up any other
			// email or organization is an unexpected call and fails the test.
			service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
				GetAccountV2ByEnvironmentID(gomock.Any(), callerEmail, envID).
				DoAndReturn(func(context.Context, string, string) (*domain.AccountV2, error) {
					return callerAccount(), nil
				}).AnyTimes()

			var persisted *domain.AccountV2
			if tc.wantErr == nil {
				sysAdminLookup := service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
					GetSystemAdminAccountV2(gomock.Any(), callerEmail)
				if tc.systemAdmin {
					sysAdminLookup.Return(systemAdminAccount(), nil)
				} else {
					sysAdminLookup.Return(nil, v2as.ErrSystemAdminAccountNotFound)
				}
				service.dbClient.(*dbmock.MockClient).EXPECT().
					RunInTransactionV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				service.publisher.(*publishermock.MockPublisher).EXPECT().
					Publish(gomock.Any(), gomock.Any()).Return(nil)
				service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
					UpdateAccountV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, a *domain.AccountV2) error {
						persisted = a
						return nil
					})
			}

			err := tc.call(ctx, service)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Nil(t, persisted)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, persisted)
			assert.Equal(t, callerEmail, persisted.Email)
			assert.Equal(t, tc.wantOrgID, persisted.OrganizationId)
			tc.check(t, persisted)
		})
	}
}
