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
	accstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	alstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	dbmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

func TestUpdateAccountV2SelfServicePrivilegeFields(t *testing.T) {
	t.Parallel()

	memberAccount := func() *domain.AccountV2 {
		return &domain.AccountV2{AccountV2: &accountproto.AccountV2{
			Email:            "member@example.com",
			OrganizationId:   "org0",
			FirstName:        "Low",
			LastName:         "Priv",
			Language:         "en",
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{EnvironmentId: "env0", Role: accountproto.AccountV2_Role_Environment_VIEWER},
				{EnvironmentId: "env1", Role: accountproto.AccountV2_Role_Environment_EDITOR},
			},
		}}
	}

	cases := []struct {
		name        string
		req         *accountproto.UpdateAccountV2Request
		wantErr     error
		wantUpdated bool
	}{
		{
			name: "rejects organization role change",
			req: &accountproto.UpdateAccountV2Request{
				Email:          "member@example.com",
				OrganizationId: "org0",
				OrganizationRole: &accountproto.UpdateAccountV2Request_OrganizationRoleValue{
					Role: accountproto.AccountV2_Role_Organization_OWNER,
				},
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name: "rejects environment role change",
			req: &accountproto.UpdateAccountV2Request{
				Email:          "member@example.com",
				OrganizationId: "org0",
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{EnvironmentId: "env0", Role: accountproto.AccountV2_Role_Environment_EDITOR},
					{EnvironmentId: "env1", Role: accountproto.AccountV2_Role_Environment_EDITOR},
				},
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name: "rejects disabled change",
			req: &accountproto.UpdateAccountV2Request{
				Email:          "member@example.com",
				OrganizationId: "org0",
				Disabled:       wrapperspb.Bool(true),
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name: "rejects combined privilege escalation",
			req: &accountproto.UpdateAccountV2Request{
				Email:          "member@example.com",
				OrganizationId: "org0",
				FirstName:      wrapperspb.String("Attacker"),
				OrganizationRole: &accountproto.UpdateAccountV2Request_OrganizationRoleValue{
					Role: accountproto.AccountV2_Role_Organization_OWNER,
				},
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{EnvironmentId: "env0", Role: accountproto.AccountV2_Role_Environment_EDITOR},
				},
				Disabled: wrapperspb.Bool(true),
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name: "allows profile update without privilege fields",
			req: &accountproto.UpdateAccountV2Request{
				Email:          "member@example.com",
				OrganizationId: "org0",
				FirstName:      wrapperspb.String("Updated"),
			},
			wantUpdated: true,
		},
		{
			name: "allows profile update when environment roles match current",
			req: &accountproto.UpdateAccountV2Request{
				Email:          "member@example.com",
				OrganizationId: "org0",
				FirstName:      wrapperspb.String("Updated"),
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{EnvironmentId: "env1", Role: accountproto.AccountV2_Role_Environment_EDITOR},
					{EnvironmentId: "env0", Role: accountproto.AccountV2_Role_Environment_VIEWER},
				},
			},
			wantUpdated: true,
		},
		{
			name: "allows profile update when organization role matches current",
			req: &accountproto.UpdateAccountV2Request{
				Email:          "member@example.com",
				OrganizationId: "org0",
				FirstName:      wrapperspb.String("Updated"),
				OrganizationRole: &accountproto.UpdateAccountV2Request_OrganizationRoleValue{
					Role: accountproto.AccountV2_Role_Organization_MEMBER,
				},
			},
			wantUpdated: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			ctx := createContextWithEmailToken(t, "member@example.com", false)
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{})

			service := createAccountService(t, mockController)

			service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
				GetAccountV2(gomock.Any(), "member@example.com", "org0").
				DoAndReturn(func(context.Context, string, string) (*domain.AccountV2, error) {
					return memberAccount(), nil
				}).AnyTimes()

			var persisted *domain.AccountV2
			if tc.wantUpdated {
				service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
					UpdateAccountV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, a *domain.AccountV2) error {
						persisted = a
						return nil
					})
				service.auditLogStorage.(*alstoragemock.MockAuditLogStorage).EXPECT().
					CreateAuditLog(gomock.Any(), gomock.Any()).Return(nil)
				service.dbClient.(*dbmock.MockClient).EXPECT().
					RunInTransactionV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				service.publisher.(*publishermock.MockPublisher).EXPECT().
					Publish(gomock.Any(), gomock.Any()).Return(nil)
			}

			resp, err := service.UpdateAccountV2(ctx, tc.req)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Nil(t, persisted)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, persisted)
			assert.Equal(t, "Updated", persisted.FirstName)
			assert.Equal(t, accountproto.AccountV2_Role_Organization_MEMBER, persisted.OrganizationRole)
			require.Len(t, persisted.EnvironmentRoles, 2)
			assert.Equal(t, accountproto.AccountV2_Role_Environment_VIEWER, persisted.EnvironmentRoles[0].Role)
			assert.Equal(t, accountproto.AccountV2_Role_Organization_MEMBER, resp.Account.OrganizationRole)
		})
	}
}

func TestEnvironmentRolesEqual(t *testing.T) {
	t.Parallel()
	base := []*accountproto.AccountV2_EnvironmentRole{
		{EnvironmentId: "env0", Role: accountproto.AccountV2_Role_Environment_VIEWER},
		{EnvironmentId: "env1", Role: accountproto.AccountV2_Role_Environment_EDITOR},
	}
	reordered := []*accountproto.AccountV2_EnvironmentRole{
		{EnvironmentId: "env1", Role: accountproto.AccountV2_Role_Environment_EDITOR},
		{EnvironmentId: "env0", Role: accountproto.AccountV2_Role_Environment_VIEWER},
	}
	assert.True(t, environmentRolesEqual(reordered, base))
	assert.False(t, environmentRolesEqual(
		[]*accountproto.AccountV2_EnvironmentRole{
			{EnvironmentId: "env0", Role: accountproto.AccountV2_Role_Environment_EDITOR},
		},
		base,
	))
}
