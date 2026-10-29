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
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	v2asmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	v2als "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2"
	v2alsmock "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2/mock"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

const (
	scopeTestOrganizationID    = "org-a"
	scopeTestEnvironmentID     = "env-a"
	scopeTestOtherEnvironment  = "env-b"
	scopeTestOtherOrganization = "org-b"
)

// combinedScopeCase describes one ListAuditLogs / GetAuditLog request that sets
// both organization_id and environment_id, and the account lookups it triggers.
type combinedScopeCase struct {
	name          string
	isSystemAdmin bool
	environmentID string
	// accountInOrganization is what GetAccountV2 returns for organization_id.
	accountInOrganization *accountproto.AccountV2
	// accountInEnvironment is what GetAccountV2ByEnvironmentID returns for
	// environment_id, i.e. the caller's account in the organization that owns
	// the environment. Nil means the caller has no account there.
	accountInEnvironment *accountproto.AccountV2
	wantErr              error
}

func combinedScopeCases() []combinedScopeCase {
	organizationAdmin := &accountproto.AccountV2{
		Email:            "test@example.com",
		OrganizationId:   scopeTestOrganizationID,
		OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
	}
	organizationMember := &accountproto.AccountV2{
		Email:            "test@example.com",
		OrganizationId:   scopeTestOrganizationID,
		OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
	}
	otherOrganizationMember := &accountproto.AccountV2{
		Email:            "test@example.com",
		OrganizationId:   scopeTestOtherOrganization,
		OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
	}
	otherOrganizationViewer := &accountproto.AccountV2{
		Email:            "test@example.com",
		OrganizationId:   scopeTestOtherOrganization,
		OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{EnvironmentId: scopeTestOtherEnvironment, Role: accountproto.AccountV2_Role_Environment_VIEWER},
		},
	}
	return []combinedScopeCase{
		{
			name:                  "rejects environment of an organization the caller has no account in",
			environmentID:         scopeTestOtherEnvironment,
			accountInOrganization: organizationAdmin,
			accountInEnvironment:  nil,
			wantErr:               statusUnauthenticated.Err(),
		},
		{
			name:                  "rejects environment the caller has no viewer role in",
			environmentID:         scopeTestOtherEnvironment,
			accountInOrganization: organizationAdmin,
			accountInEnvironment:  otherOrganizationMember,
			wantErr:               statusPermissionDenied.Err(),
		},
		{
			name:                  "rejects organization member even for an environment of the same organization",
			environmentID:         scopeTestEnvironmentID,
			accountInOrganization: organizationMember,
			wantErr:               statusPermissionDenied.Err(),
		},
		{
			name:                  "allows organization admin for an environment of the same organization",
			environmentID:         scopeTestEnvironmentID,
			accountInOrganization: organizationAdmin,
			accountInEnvironment:  organizationAdmin,
		},
		{
			name:                  "allows organization admin with viewer role in another organization's environment",
			environmentID:         scopeTestOtherEnvironment,
			accountInOrganization: organizationAdmin,
			accountInEnvironment:  otherOrganizationViewer,
		},
		{
			name:          "allows system admin without membership in either organization",
			isSystemAdmin: true,
			environmentID: scopeTestOtherEnvironment,
		},
	}
}

func TestListAuditLogsCombinedScope(t *testing.T) {
	t.Parallel()
	for _, tc := range combinedScopeCases() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			s, auditLogStorage, accountStorage := newCombinedScopeService(t, mockController, tc)
			if tc.wantErr == nil {
				auditLogStorage.EXPECT().ListAuditLogs(
					gomock.Any(),
					v2als.ListAuditLogsParams{
						EnvironmentID:  tc.environmentID,
						OrganizationID: scopeTestOrganizationID,
						OrderBy:        proto.ListAuditLogsRequest_DEFAULT,
						OrderDirection: proto.ListAuditLogsRequest_DESC,
						PageSize:       2,
						Cursor:         "0",
					},
				).Return(createAuditLogs(t), 2, int64(10), nil)
				accountStorage.EXPECT().GetAvatarAccountsV2(
					gomock.Any(), gomock.Any(),
				).Return([]*accountproto.AccountV2{}, nil)
			}
			// No storage expectation is registered for the rejected requests, so
			// any read of the audit log storage fails the test.

			actual, err := s.ListAuditLogs(combinedScopeContext(t, tc.isSystemAdmin), &proto.ListAuditLogsRequest{
				PageSize:       2,
				EnvironmentId:  tc.environmentID,
				OrganizationId: scopeTestOrganizationID,
			})
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr, err)
				assert.Nil(t, actual)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, &proto.ListAuditLogsResponse{
				AuditLogs:  createAuditLogs(t),
				Cursor:     "2",
				TotalCount: 10,
			}, actual)
		})
	}
}

func TestGetAuditLogCombinedScope(t *testing.T) {
	t.Parallel()
	for _, tc := range combinedScopeCases() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			s, auditLogStorage, accountStorage := newCombinedScopeService(t, mockController, tc)
			if tc.wantErr == nil {
				auditLogStorage.EXPECT().GetAuditLog(
					gomock.Any(), "id-1", tc.environmentID, scopeTestOrganizationID,
				).Return(&proto.AuditLog{
					Id:     "id-1",
					Editor: &domaineventproto.Editor{Email: "editor@example.com"},
				}, nil)
				accountStorage.EXPECT().GetAvatarAccountsV2(
					gomock.Any(), gomock.Any(),
				).Return([]*accountproto.AccountV2{}, nil)
			}
			// No storage expectation is registered for the rejected requests, so
			// any read of the audit log storage fails the test.

			ctx := combinedScopeContext(t, tc.isSystemAdmin)
			actual, err := s.GetAuditLog(ctx, &proto.GetAuditLogRequest{
				Id:             "id-1",
				EnvironmentId:  tc.environmentID,
				OrganizationId: scopeTestOrganizationID,
			})
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr, err)
				assert.Nil(t, actual)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, &proto.AuditLog{
				Id:               "id-1",
				Editor:           &domaineventproto.Editor{Email: "editor@example.com"},
				LocalizedMessage: domainevent.LocalizedMessage(domaineventproto.Event_UNKNOWN, locale.NewLocalizer(ctx)),
			}, actual.AuditLog)
		})
	}
}

// newCombinedScopeService builds the service with the account lookups the case
// describes: GetAccountV2 answers the organization role check and
// GetAccountV2ByEnvironmentID answers the environment role check. System admins
// pass both checks without any lookup.
func newCombinedScopeService(
	t *testing.T,
	mockController *gomock.Controller,
	tc combinedScopeCase,
) (*auditlogService, *v2alsmock.MockAuditLogStorage, *v2asmock.MockAccountStorage) {
	t.Helper()
	accountClient := accountclientmock.NewMockClient(mockController)
	accountStorage := v2asmock.NewMockAccountStorage(mockController)
	auditLogStorage := v2alsmock.NewMockAuditLogStorage(mockController)
	if tc.accountInOrganization != nil {
		accountClient.EXPECT().GetAccountV2(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, req *accountproto.GetAccountV2Request, _ ...grpc.CallOption) (*accountproto.GetAccountV2Response, error) {
				assert.Equal(t, scopeTestOrganizationID, req.OrganizationId)
				return &accountproto.GetAccountV2Response{Account: tc.accountInOrganization}, nil
			},
		)
	}
	if tc.accountInOrganization != nil &&
		tc.accountInOrganization.OrganizationRole >= accountproto.AccountV2_Role_Organization_ADMIN {
		accountClient.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, req *accountproto.GetAccountV2ByEnvironmentIDRequest, _ ...grpc.CallOption) (*accountproto.GetAccountV2ByEnvironmentIDResponse, error) {
				assert.Equal(t, tc.environmentID, req.EnvironmentId)
				if tc.accountInEnvironment == nil {
					return nil, status.Error(codes.NotFound, "account not found")
				}
				return &accountproto.GetAccountV2ByEnvironmentIDResponse{Account: tc.accountInEnvironment}, nil
			},
		)
	}
	s := NewAuditLogService(accountClient, accountStorage, auditLogStorage, WithLogger(zap.NewNop())).(*auditlogService)
	return s, auditLogStorage, accountStorage
}

func combinedScopeContext(t *testing.T, isSystemAdmin bool) context.Context {
	t.Helper()
	ctx := createContextWithToken(t, isSystemAdmin)
	return metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
}
