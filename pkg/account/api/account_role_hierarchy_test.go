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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	accstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	alstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2/mock"
	ecmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	dbmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

const (
	hierarchyOrgID          = "org0"
	hierarchyOwnerEmail     = "owner@example.com"
	hierarchyCoOwnerEmail   = "co-owner@example.com"
	hierarchyAdminEmail     = "admin@example.com"
	hierarchyMemberEmail    = "member@example.com"
	hierarchyNewMemberEmail = "new-member@example.com"
)

// hierarchyAccounts returns the organization's accounts keyed by email. The
// organization's owner email points at hierarchyOwnerEmail; hierarchyCoOwnerEmail
// holds the OWNER role without being the organization owner.
func hierarchyAccounts() map[string]*domain.AccountV2 {
	newAccount := func(email string, role accountproto.AccountV2_Role_Organization) *domain.AccountV2 {
		return &domain.AccountV2{AccountV2: &accountproto.AccountV2{
			Email:            email,
			OrganizationId:   hierarchyOrgID,
			FirstName:        "Test",
			LastName:         "User",
			Language:         "en",
			OrganizationRole: role,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{EnvironmentId: "env0", Role: accountproto.AccountV2_Role_Environment_EDITOR},
			},
		}}
	}
	return map[string]*domain.AccountV2{
		hierarchyOwnerEmail:   newAccount(hierarchyOwnerEmail, accountproto.AccountV2_Role_Organization_OWNER),
		hierarchyCoOwnerEmail: newAccount(hierarchyCoOwnerEmail, accountproto.AccountV2_Role_Organization_OWNER),
		hierarchyAdminEmail:   newAccount(hierarchyAdminEmail, accountproto.AccountV2_Role_Organization_ADMIN),
		hierarchyMemberEmail:  newAccount(hierarchyMemberEmail, accountproto.AccountV2_Role_Organization_MEMBER),
	}
}

// setupHierarchyLookups serves account lookups from hierarchyAccounts and the
// organization lookup with hierarchyOwnerEmail as the organization owner.
// Account lookups match the email case-insensitively, like the database
// collation does, and return the account with its stored email.
func setupHierarchyLookups(service *AccountService) {
	accounts := hierarchyAccounts()
	service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
		GetAccountV2(gomock.Any(), gomock.Any(), hierarchyOrgID).
		DoAndReturn(func(_ context.Context, email, _ string) (*domain.AccountV2, error) {
			for storedEmail, account := range accounts {
				if strings.EqualFold(storedEmail, email) {
					return account, nil
				}
			}
			return nil, v2as.ErrAccountNotFound
		}).AnyTimes()
	service.environmentClient.(*ecmock.MockClient).EXPECT().
		GetOrganization(gomock.Any(), gomock.Any()).
		Return(&environmentproto.GetOrganizationResponse{
			Organization: &environmentproto.Organization{
				Id:         hierarchyOrgID,
				OwnerEmail: hierarchyOwnerEmail,
			},
		}, nil).AnyTimes()
}

func TestCreateAccountV2OrganizationRoleHierarchy(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		editorEmail   string
		isSystemAdmin bool
		role          accountproto.AccountV2_Role_Organization
		wantErr       error
	}{
		{
			name:        "admin cannot create an owner",
			editorEmail: hierarchyAdminEmail,
			role:        accountproto.AccountV2_Role_Organization_OWNER,
			wantErr:     statusPermissionDenied.Err(),
		},
		{
			name:        "admin can create an admin",
			editorEmail: hierarchyAdminEmail,
			role:        accountproto.AccountV2_Role_Organization_ADMIN,
		},
		{
			name:        "admin can create a member",
			editorEmail: hierarchyAdminEmail,
			role:        accountproto.AccountV2_Role_Organization_MEMBER,
		},
		{
			name:        "owner can create an owner",
			editorEmail: hierarchyOwnerEmail,
			role:        accountproto.AccountV2_Role_Organization_OWNER,
		},
		{
			name:          "system admin who is an organization admin can create an owner",
			editorEmail:   hierarchyAdminEmail,
			isSystemAdmin: true,
			role:          accountproto.AccountV2_Role_Organization_OWNER,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			ctx := createContextWithEmailToken(t, tc.editorEmail, tc.isSystemAdmin)
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{})

			service := createAccountService(t, mockController, nil)
			setupHierarchyLookups(service)

			var created *domain.AccountV2
			if tc.wantErr == nil {
				service.dbClient.(*dbmock.MockClient).EXPECT().
					RunInTransactionV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
					CreateAccountV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, a *domain.AccountV2) error {
						created = a
						return nil
					})
				service.auditLogStorage.(*alstoragemock.MockAuditLogStorage).EXPECT().
					CreateAuditLog(gomock.Any(), gomock.Any()).Return(nil)
				service.publisher.(*publishermock.MockPublisher).EXPECT().
					Publish(gomock.Any(), gomock.Any()).Return(nil)
			}

			resp, err := service.CreateAccountV2(ctx, &accountproto.CreateAccountV2Request{
				OrganizationId:   hierarchyOrgID,
				Email:            hierarchyNewMemberEmail,
				FirstName:        "New",
				LastName:         "Member",
				Language:         "en",
				OrganizationRole: tc.role,
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{EnvironmentId: "env0", Role: accountproto.AccountV2_Role_Environment_VIEWER},
				},
			})
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Nil(t, created)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, created)
			assert.Equal(t, tc.role, created.OrganizationRole)
			assert.Equal(t, tc.role, resp.Account.OrganizationRole)
		})
	}
}

func TestUpdateAccountV2OrganizationRoleHierarchy(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		editorEmail   string
		isSystemAdmin bool
		req           *accountproto.UpdateAccountV2Request
		wantErr       error
		wantRole      accountproto.AccountV2_Role_Organization
		wantDisabled  bool
	}{
		{
			name:        "admin cannot update an owner",
			editorEmail: hierarchyAdminEmail,
			req: &accountproto.UpdateAccountV2Request{
				OrganizationId: hierarchyOrgID,
				Email:          hierarchyOwnerEmail,
				FirstName:      wrapperspb.String("Updated"),
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name:        "admin cannot promote a member to owner",
			editorEmail: hierarchyAdminEmail,
			req: &accountproto.UpdateAccountV2Request{
				OrganizationId: hierarchyOrgID,
				Email:          hierarchyMemberEmail,
				OrganizationRole: &accountproto.UpdateAccountV2Request_OrganizationRoleValue{
					Role: accountproto.AccountV2_Role_Organization_OWNER,
				},
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name:        "admin cannot promote themselves to owner",
			editorEmail: hierarchyAdminEmail,
			req: &accountproto.UpdateAccountV2Request{
				OrganizationId: hierarchyOrgID,
				Email:          hierarchyAdminEmail,
				OrganizationRole: &accountproto.UpdateAccountV2Request_OrganizationRoleValue{
					Role: accountproto.AccountV2_Role_Organization_OWNER,
				},
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name:        "owner cannot disable the organization owner",
			editorEmail: hierarchyCoOwnerEmail,
			req: &accountproto.UpdateAccountV2Request{
				OrganizationId: hierarchyOrgID,
				Email:          hierarchyOwnerEmail,
				Disabled:       wrapperspb.Bool(true),
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name:          "system admin cannot disable the organization owner",
			editorEmail:   hierarchyAdminEmail,
			isSystemAdmin: true,
			req: &accountproto.UpdateAccountV2Request{
				OrganizationId: hierarchyOrgID,
				Email:          hierarchyOwnerEmail,
				Disabled:       wrapperspb.Bool(true),
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name:        "owner cannot disable the organization owner using a differently cased email",
			editorEmail: hierarchyCoOwnerEmail,
			req: &accountproto.UpdateAccountV2Request{
				OrganizationId: hierarchyOrgID,
				Email:          strings.ToUpper(hierarchyOwnerEmail),
				Disabled:       wrapperspb.Bool(true),
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name:        "admin can promote a member to admin",
			editorEmail: hierarchyAdminEmail,
			req: &accountproto.UpdateAccountV2Request{
				OrganizationId: hierarchyOrgID,
				Email:          hierarchyMemberEmail,
				OrganizationRole: &accountproto.UpdateAccountV2Request_OrganizationRoleValue{
					Role: accountproto.AccountV2_Role_Organization_ADMIN,
				},
			},
			wantRole: accountproto.AccountV2_Role_Organization_ADMIN,
		},
		{
			name:        "admin can disable a member",
			editorEmail: hierarchyAdminEmail,
			req: &accountproto.UpdateAccountV2Request{
				OrganizationId: hierarchyOrgID,
				Email:          hierarchyMemberEmail,
				Disabled:       wrapperspb.Bool(true),
			},
			wantRole:     accountproto.AccountV2_Role_Organization_MEMBER,
			wantDisabled: true,
		},
		{
			name:        "owner can promote a member to owner",
			editorEmail: hierarchyOwnerEmail,
			req: &accountproto.UpdateAccountV2Request{
				OrganizationId: hierarchyOrgID,
				Email:          hierarchyMemberEmail,
				OrganizationRole: &accountproto.UpdateAccountV2Request_OrganizationRoleValue{
					Role: accountproto.AccountV2_Role_Organization_OWNER,
				},
			},
			wantRole: accountproto.AccountV2_Role_Organization_OWNER,
		},
		{
			name:        "owner can disable another owner who is not the organization owner",
			editorEmail: hierarchyOwnerEmail,
			req: &accountproto.UpdateAccountV2Request{
				OrganizationId: hierarchyOrgID,
				Email:          hierarchyCoOwnerEmail,
				Disabled:       wrapperspb.Bool(true),
			},
			wantRole:     accountproto.AccountV2_Role_Organization_OWNER,
			wantDisabled: true,
		},
		{
			name:          "system admin who is an organization admin can update an owner",
			editorEmail:   hierarchyAdminEmail,
			isSystemAdmin: true,
			req: &accountproto.UpdateAccountV2Request{
				OrganizationId: hierarchyOrgID,
				Email:          hierarchyOwnerEmail,
				FirstName:      wrapperspb.String("Updated"),
			},
			wantRole: accountproto.AccountV2_Role_Organization_OWNER,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			ctx := createContextWithEmailToken(t, tc.editorEmail, tc.isSystemAdmin)
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{})

			service := createAccountService(t, mockController, nil)
			setupHierarchyLookups(service)

			var persisted *domain.AccountV2
			if tc.wantErr == nil {
				service.dbClient.(*dbmock.MockClient).EXPECT().
					RunInTransactionV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
					UpdateAccountV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, a *domain.AccountV2) error {
						persisted = a
						return nil
					})
				service.auditLogStorage.(*alstoragemock.MockAuditLogStorage).EXPECT().
					CreateAuditLog(gomock.Any(), gomock.Any()).Return(nil)
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
			assert.Equal(t, tc.req.Email, persisted.Email)
			assert.Equal(t, tc.wantRole, persisted.OrganizationRole)
			assert.Equal(t, tc.wantDisabled, persisted.Disabled)
			assert.Equal(t, tc.wantRole, resp.Account.OrganizationRole)
		})
	}
}

func TestDisableAccountV2OrganizationRoleHierarchy(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		editorEmail   string
		isSystemAdmin bool
		targetEmail   string
		wantErr       error
	}{
		{
			name:        "admin cannot disable the organization owner",
			editorEmail: hierarchyAdminEmail,
			targetEmail: hierarchyOwnerEmail,
			wantErr:     statusPermissionDenied.Err(),
		},
		{
			name:        "admin cannot disable an owner",
			editorEmail: hierarchyAdminEmail,
			targetEmail: hierarchyCoOwnerEmail,
			wantErr:     statusPermissionDenied.Err(),
		},
		{
			name:        "owner cannot disable the organization owner",
			editorEmail: hierarchyCoOwnerEmail,
			targetEmail: hierarchyOwnerEmail,
			wantErr:     statusPermissionDenied.Err(),
		},
		{
			name:          "system admin cannot disable the organization owner",
			editorEmail:   hierarchyAdminEmail,
			isSystemAdmin: true,
			targetEmail:   hierarchyOwnerEmail,
			wantErr:       statusPermissionDenied.Err(),
		},
		{
			name:        "owner cannot disable the organization owner using a differently cased email",
			editorEmail: hierarchyCoOwnerEmail,
			targetEmail: strings.ToUpper(hierarchyOwnerEmail),
			wantErr:     statusPermissionDenied.Err(),
		},
		{
			name:        "admin can disable a member",
			editorEmail: hierarchyAdminEmail,
			targetEmail: hierarchyMemberEmail,
		},
		{
			name:        "admin can disable a member using a differently cased email",
			editorEmail: hierarchyAdminEmail,
			targetEmail: strings.ToUpper(hierarchyMemberEmail),
		},
		{
			name:        "owner can disable another owner who is not the organization owner",
			editorEmail: hierarchyOwnerEmail,
			targetEmail: hierarchyCoOwnerEmail,
		},
		{
			name:          "system admin who is an organization admin can disable an owner",
			editorEmail:   hierarchyAdminEmail,
			isSystemAdmin: true,
			targetEmail:   hierarchyCoOwnerEmail,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			ctx := createContextWithEmailToken(t, tc.editorEmail, tc.isSystemAdmin)
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{})

			service := createAccountService(t, mockController, nil)
			setupHierarchyLookups(service)

			var persisted *domain.AccountV2
			if tc.wantErr == nil {
				service.dbClient.(*dbmock.MockClient).EXPECT().
					RunInTransactionV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
					UpdateAccountV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, a *domain.AccountV2) error {
						persisted = a
						return nil
					})
				service.auditLogStorage.(*alstoragemock.MockAuditLogStorage).EXPECT().
					CreateAuditLog(gomock.Any(), gomock.Any()).Return(nil)
				service.publisher.(*publishermock.MockPublisher).EXPECT().
					Publish(gomock.Any(), gomock.Any()).Return(nil)
			}

			resp, err := service.DisableAccountV2(ctx, &accountproto.DisableAccountV2Request{
				OrganizationId: hierarchyOrgID,
				Email:          tc.targetEmail,
			})
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Nil(t, persisted)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, persisted)
			assert.True(t, strings.EqualFold(tc.targetEmail, persisted.Email))
			assert.True(t, persisted.Disabled)
			assert.True(t, resp.Account.Disabled)
		})
	}
}

func TestDeleteAccountV2OrganizationRoleHierarchy(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		editorEmail   string
		isSystemAdmin bool
		targetEmail   string
		wantErr       error
	}{
		{
			name:        "admin cannot delete the organization owner",
			editorEmail: hierarchyAdminEmail,
			targetEmail: hierarchyOwnerEmail,
			wantErr:     statusPermissionDenied.Err(),
		},
		{
			name:        "admin cannot delete an owner",
			editorEmail: hierarchyAdminEmail,
			targetEmail: hierarchyCoOwnerEmail,
			wantErr:     statusPermissionDenied.Err(),
		},
		{
			name:        "owner cannot delete the organization owner",
			editorEmail: hierarchyCoOwnerEmail,
			targetEmail: hierarchyOwnerEmail,
			wantErr:     statusPermissionDenied.Err(),
		},
		{
			name:          "system admin cannot delete the organization owner",
			editorEmail:   hierarchyAdminEmail,
			isSystemAdmin: true,
			targetEmail:   hierarchyOwnerEmail,
			wantErr:       statusPermissionDenied.Err(),
		},
		{
			name:        "owner cannot delete the organization owner using a differently cased email",
			editorEmail: hierarchyCoOwnerEmail,
			targetEmail: strings.ToUpper(hierarchyOwnerEmail),
			wantErr:     statusPermissionDenied.Err(),
		},
		{
			name:        "admin can delete a member",
			editorEmail: hierarchyAdminEmail,
			targetEmail: hierarchyMemberEmail,
		},
		{
			name:        "admin can delete a member using a differently cased email",
			editorEmail: hierarchyAdminEmail,
			targetEmail: strings.ToUpper(hierarchyMemberEmail),
		},
		{
			name:        "admin can delete an admin",
			editorEmail: hierarchyAdminEmail,
			targetEmail: hierarchyAdminEmail,
		},
		{
			name:        "owner can delete another owner who is not the organization owner",
			editorEmail: hierarchyOwnerEmail,
			targetEmail: hierarchyCoOwnerEmail,
		},
		{
			name:          "system admin who is an organization admin can delete an owner",
			editorEmail:   hierarchyAdminEmail,
			isSystemAdmin: true,
			targetEmail:   hierarchyCoOwnerEmail,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			ctx := createContextWithEmailToken(t, tc.editorEmail, tc.isSystemAdmin)
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{})

			service := createAccountService(t, mockController, nil)
			setupHierarchyLookups(service)

			var deleted *domain.AccountV2
			if tc.wantErr == nil {
				service.dbClient.(*dbmock.MockClient).EXPECT().
					RunInTransactionV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				service.publisher.(*publishermock.MockPublisher).EXPECT().
					Publish(gomock.Any(), gomock.Any()).Return(nil)
				service.accountStorage.(*accstoragemock.MockAccountStorage).EXPECT().
					DeleteAccountV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, a *domain.AccountV2) error {
						deleted = a
						return nil
					})
			}

			_, err := service.DeleteAccountV2(ctx, &accountproto.DeleteAccountV2Request{
				OrganizationId: hierarchyOrgID,
				Email:          tc.targetEmail,
			})
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Nil(t, deleted)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, deleted)
			assert.True(t, strings.EqualFold(tc.targetEmail, deleted.Email))
		})
	}
}

func TestCanManageOrganizationRole(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		editor *eventproto.Editor
		role   accountproto.AccountV2_Role_Organization
		want   bool
	}{
		{
			name:   "admin manages member",
			editor: &eventproto.Editor{OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN},
			role:   accountproto.AccountV2_Role_Organization_MEMBER,
			want:   true,
		},
		{
			name:   "admin manages admin",
			editor: &eventproto.Editor{OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN},
			role:   accountproto.AccountV2_Role_Organization_ADMIN,
			want:   true,
		},
		{
			name:   "admin does not manage owner",
			editor: &eventproto.Editor{OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN},
			role:   accountproto.AccountV2_Role_Organization_OWNER,
			want:   false,
		},
		{
			name:   "owner manages owner",
			editor: &eventproto.Editor{OrganizationRole: accountproto.AccountV2_Role_Organization_OWNER},
			role:   accountproto.AccountV2_Role_Organization_OWNER,
			want:   true,
		},
		{
			name:   "system admin without an organization role manages owner",
			editor: &eventproto.Editor{IsAdmin: true},
			role:   accountproto.AccountV2_Role_Organization_OWNER,
			want:   true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, canManageOrganizationRole(tc.editor, tc.role))
		})
	}
}
