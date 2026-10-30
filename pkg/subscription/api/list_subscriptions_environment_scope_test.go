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

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	v2ss "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/storage/v2"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/storage/v2/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
)

// listScopeAccount returns the account of the caller for the environment
// scoping tests: a member of the organization with roles in env0 and env1 only.
func listScopeAccount(orgRole accountproto.AccountV2_Role_Organization) *accountproto.AccountV2 {
	return &accountproto.AccountV2{
		Email:            "email",
		OrganizationRole: orgRole,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{EnvironmentId: "env0", Role: accountproto.AccountV2_Role_Environment_VIEWER},
			{EnvironmentId: "env1", Role: accountproto.AccountV2_Role_Environment_EDITOR},
		},
	}
}

func TestListSubscriptionsMemberEnvironmentScope(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name            string
		account         *accountproto.AccountV2
		req             *proto.ListSubscriptionsRequest
		wantEnvFilter   []string
		wantStorageCall bool
		wantErr         error
	}{
		{
			name:    "member requesting only environments outside their roles gets an empty page",
			account: listScopeAccount(accountproto.AccountV2_Role_Organization_MEMBER),
			req: &proto.ListSubscriptionsRequest{
				OrganizationId: "org0",
				EnvironmentIds: []string{"env2", "env3"},
				PageSize:       10,
				Cursor:         "5",
			},
		},
		{
			name: "member without any environment role gets an empty page",
			account: &accountproto.AccountV2{
				Email:            "email",
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			},
			req: &proto.ListSubscriptionsRequest{OrganizationId: "org0", PageSize: 10},
		},
		{
			name:    "member with no accessible environment still gets an invalid cursor rejected",
			account: listScopeAccount(accountproto.AccountV2_Role_Organization_MEMBER),
			req: &proto.ListSubscriptionsRequest{
				OrganizationId: "org0",
				EnvironmentIds: []string{"env2"},
				Cursor:         "XXX",
			},
			wantErr: statusInvalidCursor.Err(),
		},
		{
			name:    "member with no accessible environment still gets an invalid order rejected",
			account: listScopeAccount(accountproto.AccountV2_Role_Organization_MEMBER),
			req: &proto.ListSubscriptionsRequest{
				OrganizationId: "org0",
				EnvironmentIds: []string{"env2"},
				OrderBy:        proto.ListSubscriptionsRequest_OrderBy(999),
			},
			wantErr: statusInvalidOrderBy.Err(),
		},
		{
			name:    "member requesting a mix keeps only the environments they have a role in",
			account: listScopeAccount(accountproto.AccountV2_Role_Organization_MEMBER),
			req: &proto.ListSubscriptionsRequest{
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
			req:             &proto.ListSubscriptionsRequest{OrganizationId: "org0", PageSize: 10},
			wantEnvFilter:   []string{"env0", "env1"},
			wantStorageCall: true,
		},
		{
			name:            "admin without an environment filter is not restricted",
			account:         listScopeAccount(accountproto.AccountV2_Role_Organization_ADMIN),
			req:             &proto.ListSubscriptionsRequest{OrganizationId: "org0", PageSize: 10},
			wantEnvFilter:   []string{},
			wantStorageCall: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := newSubscriptionService(mockController, nil, nil, nil)
			service.accountClient.(*accountclientmock.MockClient).EXPECT().
				GetAccountV2(gomock.Any(), gomock.Any()).
				Return(&accountproto.GetAccountV2Response{Account: tc.account}, nil).AnyTimes()

			var listed *v2ss.ListSubscriptionsParams
			if tc.wantStorageCall {
				service.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().
					ListSubscriptions(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2ss.ListSubscriptionsParams) ([]*proto.Subscription, int, int64, error) {
						listed = &p
						return []*proto.Subscription{storedSubscription().Subscription}, 1, 1, nil
					})
			}

			ctx := setToken(t, newIncomingContext(t), false)
			resp, err := service.ListSubscriptions(ctx, tc.req)
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
				assert.Empty(t, resp.Subscriptions)
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
			require.Len(t, resp.Subscriptions, 1)
			assert.Equal(t, "sub-0", resp.Subscriptions[0].Id)
			assert.Equal(t, maskedSlackWebhookURL, resp.Subscriptions[0].Recipient.SlackChannelRecipient.WebhookUrl)
			assert.Equal(t, int64(1), resp.TotalCount)
		})
	}
}
