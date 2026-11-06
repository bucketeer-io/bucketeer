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
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

const (
	visibilityOrgID       = "org-visibility"
	visibilityMemberEmail = "member@example.com"
	visibilityAdminEmail  = "admin@example.com"
	visibilityPeerEmail   = "peer@example.com"
	visibilityOtherEmail  = "other@example.com"
	visibilityUnknown     = "unknown@example.com"
	// A distinct account whose email differs from the member's only by case,
	// as a case-sensitive store allows.
	visibilityCaseTwin = "MEMBER@example.com"
)

// visibilityAccounts returns an organization where the member only has a role
// in env-a, the peer has a different role in env-a, and the other account only
// has a role in env-b.
func visibilityAccounts() map[string]*domain.AccountV2 {
	return map[string]*domain.AccountV2{
		visibilityMemberEmail: {AccountV2: &accountproto.AccountV2{
			Email:            visibilityMemberEmail,
			OrganizationId:   visibilityOrgID,
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{EnvironmentId: "env-a", Role: accountproto.AccountV2_Role_Environment_VIEWER},
			},
		}},
		visibilityAdminEmail: {AccountV2: &accountproto.AccountV2{
			Email:            visibilityAdminEmail,
			OrganizationId:   visibilityOrgID,
			OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
		}},
		visibilityPeerEmail: {AccountV2: &accountproto.AccountV2{
			Email:            visibilityPeerEmail,
			OrganizationId:   visibilityOrgID,
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{EnvironmentId: "env-a", Role: accountproto.AccountV2_Role_Environment_EDITOR},
				{EnvironmentId: "env-b", Role: accountproto.AccountV2_Role_Environment_EDITOR},
			},
		}},
		visibilityOtherEmail: {AccountV2: &accountproto.AccountV2{
			Email:            visibilityOtherEmail,
			OrganizationId:   visibilityOrgID,
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{EnvironmentId: "env-b", Role: accountproto.AccountV2_Role_Environment_EDITOR},
			},
		}},
		visibilityCaseTwin: {AccountV2: &accountproto.AccountV2{
			Email:            visibilityCaseTwin,
			OrganizationId:   visibilityOrgID,
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{EnvironmentId: "env-b", Role: accountproto.AccountV2_Role_Environment_EDITOR},
			},
		}},
	}
}

// setupVisibilityLookups serves the editor and target lookups of both reads
// from visibilityAccounts, by exact email like a case-sensitive store, within
// the organization.
func setupVisibilityLookups(service *AccountService) {
	accounts := visibilityAccounts()
	lookup := func(email string) (*domain.AccountV2, error) {
		if account, ok := accounts[email]; ok {
			return account, nil
		}
		return nil, v2as.ErrAccountNotFound
	}
	storage := service.accountStorage.(*accstoragemock.MockAccountStorage)
	storage.EXPECT().GetAccountV2(gomock.Any(), gomock.Any(), visibilityOrgID).
		DoAndReturn(func(_ context.Context, email, _ string) (*domain.AccountV2, error) {
			return lookup(email)
		}).AnyTimes()
	storage.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, email, _ string) (*domain.AccountV2, error) {
			return lookup(email)
		}).AnyTimes()
}

func TestGetAccountV2MemberVisibility(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		editorEmail string
		target      string
		wantErr     error
	}{
		{name: "member reads own account", editorEmail: visibilityMemberEmail, target: visibilityMemberEmail},
		{name: "member reads an organization admin", editorEmail: visibilityMemberEmail, target: visibilityAdminEmail},
		{name: "member reads a peer of a shared environment regardless of role", editorEmail: visibilityMemberEmail, target: visibilityPeerEmail},
		{
			name: "member cannot read an account of another environment", editorEmail: visibilityMemberEmail,
			target: visibilityOtherEmail, wantErr: statusAccountNotFound.Err(),
		},
		{
			// The invisible account and an unknown one get the same answer.
			name: "member reads an unknown account", editorEmail: visibilityMemberEmail,
			target: visibilityUnknown, wantErr: statusAccountNotFound.Err(),
		},
		{
			// A distinct account differing only by email case is not the member's own.
			name: "member cannot read a case-distinct account of another environment", editorEmail: visibilityMemberEmail,
			target: visibilityCaseTwin, wantErr: statusAccountNotFound.Err(),
		},
		{name: "admin reads an account of any environment", editorEmail: visibilityAdminEmail, target: visibilityOtherEmail},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			service := createAccountService(t, mockController, nil)
			setupVisibilityLookups(service)
			ctx := metadata.NewIncomingContext(createContextWithEmailToken(t, tc.editorEmail, false), metadata.MD{})

			res, err := service.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
				Email: tc.target, OrganizationId: visibilityOrgID,
			})
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Nil(t, res)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.target, res.Account.Email)
		})
	}
}

func TestGetAccountV2ByEnvironmentIDMemberVisibility(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		editorEmail   string
		environmentID string
		target        string
		wantErr       error
	}{
		{name: "member reads a peer through a shared environment", editorEmail: visibilityMemberEmail, environmentID: "env-a", target: visibilityPeerEmail},
		{name: "member reads an organization admin through own environment", editorEmail: visibilityMemberEmail, environmentID: "env-a", target: visibilityAdminEmail},
		{
			name: "member cannot read an account of another environment", editorEmail: visibilityMemberEmail,
			environmentID: "env-a", target: visibilityOtherEmail, wantErr: statusAccountNotFound.Err(),
		},
		{
			// Even a visible peer is not readable through an environment the member has no role in.
			name: "member cannot look up through an environment without a role", editorEmail: visibilityMemberEmail,
			environmentID: "env-b", target: visibilityPeerEmail, wantErr: statusAccountNotFound.Err(),
		},
		{name: "admin reads through any environment", editorEmail: visibilityAdminEmail, environmentID: "env-b", target: visibilityOtherEmail},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			service := createAccountService(t, mockController, nil)
			setupVisibilityLookups(service)
			ctx := metadata.NewIncomingContext(createContextWithEmailToken(t, tc.editorEmail, false), metadata.MD{})

			res, err := service.GetAccountV2ByEnvironmentID(ctx, &accountproto.GetAccountV2ByEnvironmentIDRequest{
				Email: tc.target, EnvironmentId: tc.environmentID,
			})
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Nil(t, res)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.target, res.Account.Email)
		})
	}
}

func TestCanReadAccount(t *testing.T) {
	t.Parallel()

	member := &eventproto.Editor{
		Email:            visibilityMemberEmail,
		OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{EnvironmentId: "env-a", Role: accountproto.AccountV2_Role_Environment_VIEWER},
		},
	}
	accounts := visibilityAccounts()
	cases := []struct {
		name    string
		editor  *eventproto.Editor
		account *accountproto.AccountV2
		want    bool
	}{
		{name: "member sees self", editor: member, account: accounts[visibilityMemberEmail].AccountV2, want: true},
		{name: "member does not see a case-distinct account without shared access", editor: member, account: &accountproto.AccountV2{Email: visibilityCaseTwin}, want: false},
		{name: "member sees organization admin", editor: member, account: accounts[visibilityAdminEmail].AccountV2, want: true},
		{name: "member sees peer sharing an environment", editor: member, account: accounts[visibilityPeerEmail].AccountV2, want: true},
		{name: "member does not see account of another environment", editor: member, account: accounts[visibilityOtherEmail].AccountV2, want: false},
		{
			// An UNASSIGNED entry grants no access to its environment.
			name: "member's unassigned environment entry does not grant visibility",
			editor: &eventproto.Editor{
				Email:            visibilityMemberEmail,
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{EnvironmentId: "env-b", Role: accountproto.AccountV2_Role_Environment_UNASSIGNED},
				},
			},
			account: accounts[visibilityOtherEmail].AccountV2, want: false,
		},
		{
			// The target's UNASSIGNED entry in the member's environment is not shared access.
			name:   "member does not see an account whose only entry in the shared environment is unassigned",
			editor: member,
			account: &accountproto.AccountV2{
				Email:            visibilityOtherEmail,
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{EnvironmentId: "env-a", Role: accountproto.AccountV2_Role_Environment_UNASSIGNED},
					{EnvironmentId: "env-b", Role: accountproto.AccountV2_Role_Environment_EDITOR},
				},
			},
			want: false,
		},
		{name: "organization admin sees everyone", editor: &eventproto.Editor{OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN}, account: accounts[visibilityOtherEmail].AccountV2, want: true},
		{name: "system admin sees everyone", editor: &eventproto.Editor{IsAdmin: true}, account: accounts[visibilityOtherEmail].AccountV2, want: true},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, canReadAccount(tc.editor, tc.account))
		})
	}
}
