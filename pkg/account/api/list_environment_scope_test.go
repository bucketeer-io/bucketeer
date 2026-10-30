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
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	accstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

// listScopeAccount returns the account of the caller for the environment
// scoping tests: a member of org0 with roles in env0 and env1 only.
func listScopeAccount(orgRole accountproto.AccountV2_Role_Organization) *domain.AccountV2 {
	return &domain.AccountV2{AccountV2: &accountproto.AccountV2{
		Email:            "bucketeer@example.com",
		OrganizationId:   "org0",
		OrganizationRole: orgRole,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{EnvironmentId: "env0", Role: accountproto.AccountV2_Role_Environment_VIEWER},
			{EnvironmentId: "env1", Role: accountproto.AccountV2_Role_Environment_EDITOR},
		},
	}}
}

func TestListAPIKeysMemberEnvironmentScope(t *testing.T) {
	t.Parallel()

	storedKeys := []*accountproto.APIKey{{Id: "key0", Name: "key0", ApiKey: "1234567890abcdef"}}

	cases := []struct {
		name            string
		account         *domain.AccountV2
		req             *accountproto.ListAPIKeysRequest
		wantEnvFilter   []string
		wantStorageCall bool
		wantErr         error
	}{
		{
			name:    "member requesting only environments outside their roles gets an empty page",
			account: listScopeAccount(accountproto.AccountV2_Role_Organization_MEMBER),
			req: &accountproto.ListAPIKeysRequest{
				OrganizationId: "org0",
				EnvironmentIds: []string{"env2", "env3"},
				PageSize:       10,
				Cursor:         "5",
			},
		},
		{
			name: "member without any environment role gets an empty page",
			account: &domain.AccountV2{AccountV2: &accountproto.AccountV2{
				Email:            "bucketeer@example.com",
				OrganizationId:   "org0",
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			}},
			req: &accountproto.ListAPIKeysRequest{OrganizationId: "org0", PageSize: 10},
		},
		{
			name:    "member with no accessible environment still gets an invalid cursor rejected",
			account: listScopeAccount(accountproto.AccountV2_Role_Organization_MEMBER),
			req: &accountproto.ListAPIKeysRequest{
				OrganizationId: "org0",
				EnvironmentIds: []string{"env2"},
				Cursor:         "XXX",
			},
			wantErr: statusInvalidCursor.Err(),
		},
		{
			name:    "member requesting a mix keeps only the environments they have a role in",
			account: listScopeAccount(accountproto.AccountV2_Role_Organization_MEMBER),
			req: &accountproto.ListAPIKeysRequest{
				OrganizationId: "org0",
				EnvironmentIds: []string{"env2", "env1"},
				PageSize:       10,
			},
			wantEnvFilter:   []string{"env1"},
			wantStorageCall: true,
		},
		{
			name:            "member without an environment filter lists all of their environments",
			account:         listScopeAccount(accountproto.AccountV2_Role_Organization_MEMBER),
			req:             &accountproto.ListAPIKeysRequest{OrganizationId: "org0", PageSize: 10},
			wantEnvFilter:   []string{"env0", "env1"},
			wantStorageCall: true,
		},
		{
			name:            "admin without an environment filter is not restricted",
			account:         listScopeAccount(accountproto.AccountV2_Role_Organization_ADMIN),
			req:             &accountproto.ListAPIKeysRequest{OrganizationId: "org0", PageSize: 10},
			wantEnvFilter:   []string{},
			wantStorageCall: true,
		},
		{
			name:    "admin requesting an environment outside their roles is not restricted",
			account: listScopeAccount(accountproto.AccountV2_Role_Organization_ADMIN),
			req: &accountproto.ListAPIKeysRequest{
				OrganizationId: "org0",
				EnvironmentIds: []string{"env2"},
				PageSize:       10,
			},
			wantEnvFilter:   []string{"env2"},
			wantStorageCall: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			ctx := createContextWithDefaultToken(t, false)
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{})

			service := createAccountService(t, mockController, nil)
			service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
				GetAccountV2(gomock.Any(), "bucketeer@example.com", "org0").
				Return(tc.account, nil).AnyTimes()

			var listed *v2as.ListAPIKeysParams
			if tc.wantStorageCall {
				service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
					ListAPIKeys(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2as.ListAPIKeysParams) ([]*accountproto.APIKey, int, int64, error) {
						listed = &p
						return storedKeys, 1, 1, nil
					})
			}

			resp, err := service.ListAPIKeys(ctx, tc.req)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Nil(t, resp)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			if !tc.wantStorageCall {
				// Nothing is read from the storage and the page is empty.
				assert.Empty(t, resp.ApiKeys)
				assert.Equal(t, int64(0), resp.TotalCount)
				expectedCursor := tc.req.Cursor
				if expectedCursor == "" {
					expectedCursor = "0"
				}
				assert.Equal(t, expectedCursor, resp.Cursor)
				return
			}
			require.NotNil(t, listed)
			assert.Equal(t, "org0", listed.OrganizationID)
			assert.Equal(t, tc.wantEnvFilter, listed.EnvironmentIDs)
			require.Len(t, resp.ApiKeys, 1)
			assert.Equal(t, "key0", resp.ApiKeys[0].Id)
			assert.Equal(t, int64(1), resp.TotalCount)
		})
	}
}

func TestListAccountsV2MemberEnvironmentScope(t *testing.T) {
	t.Parallel()

	storedAccounts := []*accountproto.AccountV2{{Email: "other@example.com", OrganizationId: "org0"}}

	cases := []struct {
		name            string
		account         *domain.AccountV2
		req             *accountproto.ListAccountsV2Request
		wantEnvRoles    []string
		wantStorageCall bool
		wantErr         error
	}{
		{
			name:    "member requesting an environment outside their roles gets an empty page",
			account: listScopeAccount(accountproto.AccountV2_Role_Organization_MEMBER),
			req: &accountproto.ListAccountsV2Request{
				OrganizationId: "org0",
				EnvironmentId:  wrapperspb.String("env2"),
				PageSize:       10,
				Cursor:         "5",
			},
		},
		{
			name:    "member requesting a role they do not hold in the environment gets an empty page",
			account: listScopeAccount(accountproto.AccountV2_Role_Organization_MEMBER),
			req: &accountproto.ListAccountsV2Request{
				OrganizationId:  "org0",
				EnvironmentId:   wrapperspb.String("env0"),
				EnvironmentRole: wrapperspb.Int32(int32(accountproto.AccountV2_Role_Environment_EDITOR)),
				PageSize:        10,
			},
		},
		{
			name: "member without any environment role gets an empty page",
			account: &domain.AccountV2{AccountV2: &accountproto.AccountV2{
				Email:            "bucketeer@example.com",
				OrganizationId:   "org0",
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			}},
			req: &accountproto.ListAccountsV2Request{OrganizationId: "org0", PageSize: 10},
		},
		{
			name:    "member with no accessible environment still gets an invalid cursor rejected",
			account: listScopeAccount(accountproto.AccountV2_Role_Organization_MEMBER),
			req: &accountproto.ListAccountsV2Request{
				OrganizationId: "org0",
				EnvironmentId:  wrapperspb.String("env2"),
				Cursor:         "XXX",
			},
			wantErr: statusInvalidCursor.Err(),
		},
		{
			name:    "member requesting an environment they have a role in is scoped to it",
			account: listScopeAccount(accountproto.AccountV2_Role_Organization_MEMBER),
			req: &accountproto.ListAccountsV2Request{
				OrganizationId: "org0",
				EnvironmentId:  wrapperspb.String("env1"),
				PageSize:       10,
			},
			wantEnvRoles:    []string{"env1"},
			wantStorageCall: true,
		},
		{
			name:            "member without an environment filter lists all of their environments",
			account:         listScopeAccount(accountproto.AccountV2_Role_Organization_MEMBER),
			req:             &accountproto.ListAccountsV2Request{OrganizationId: "org0", PageSize: 10},
			wantEnvRoles:    []string{"env0", "env1"},
			wantStorageCall: true,
		},
		{
			name:            "admin without an environment filter is not restricted",
			account:         listScopeAccount(accountproto.AccountV2_Role_Organization_ADMIN),
			req:             &accountproto.ListAccountsV2Request{OrganizationId: "org0", PageSize: 10},
			wantEnvRoles:    []string{},
			wantStorageCall: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			ctx := createContextWithDefaultToken(t, false)
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{})

			service := createAccountService(t, mockController, nil)
			service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
				GetAccountV2(gomock.Any(), "bucketeer@example.com", "org0").
				Return(tc.account, nil).AnyTimes()

			var listed *v2as.ListAccountsV2Params
			if tc.wantStorageCall {
				service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
					ListAccountsV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2as.ListAccountsV2Params) ([]*accountproto.AccountV2, int, int64, error) {
						listed = &p
						return storedAccounts, 1, 1, nil
					})
			}

			resp, err := service.ListAccountsV2(ctx, tc.req)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Nil(t, resp)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			if !tc.wantStorageCall {
				// Nothing is read from the storage and the page is empty.
				assert.Empty(t, resp.Accounts)
				assert.Equal(t, int64(0), resp.TotalCount)
				expectedCursor := tc.req.Cursor
				if expectedCursor == "" {
					expectedCursor = "0"
				}
				assert.Equal(t, expectedCursor, resp.Cursor)
				return
			}
			require.NotNil(t, listed)
			assert.Equal(t, "org0", listed.OrganizationID)
			envIDs := make([]string, 0, len(listed.EnvironmentRoles))
			for _, r := range listed.EnvironmentRoles {
				envIDs = append(envIDs, r.EnvironmentId)
			}
			assert.Equal(t, tc.wantEnvRoles, envIDs)
			require.Len(t, resp.Accounts, 1)
			assert.Equal(t, "other@example.com", resp.Accounts[0].Email)
			assert.Equal(t, int64(1), resp.TotalCount)
		})
	}
}
