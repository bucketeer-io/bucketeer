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

package authz

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	auditlogproto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	insightsproto "github.com/bucketeer-io/bucketeer/v2/proto/insights"
	subscriptionproto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
)

const (
	updateFeature      = "/bucketeer.feature.FeatureService/UpdateFeature"
	getFeature         = "/bucketeer.feature.FeatureService/GetFeature"
	cloneFeature       = "/bucketeer.feature.FeatureService/CloneFeature"
	getInsights        = "/bucketeer.insights.InsightsService/GetInsightsMonthlySummary"
	listAuditLogs      = "/bucketeer.auditlog.AuditLogService/ListAuditLogs"
	updateAccount      = "/bucketeer.account.AccountService/UpdateAccountV2"
	listAccounts       = "/bucketeer.account.AccountService/ListAccountsV2"
	getSubscription    = "/bucketeer.subscription.SubscriptionService/GetSubscription"
	listNotifications  = "/bucketeer.notification.NotificationService/ListNotifications"
	invalidRule        = "/bucketeer.feature.FeatureService/InvalidRule"
	createOrganization = "/bucketeer.environment.EnvironmentService/CreateOrganization"
	ping               = "/bucketeer.gateway.Gateway/Ping"

	rawAPIKey        = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	obfuscatedAPIKey = "0123....cdef"
)

var errGuard = status.Error(codes.FailedPrecondition, "guard")

type guardFunc func(ctx context.Context, access *Access) error

func (f guardFunc) Check(ctx context.Context, access *Access) error { return f(ctx, access) }

func testRules(guard Guard) Rules {
	return Rules{
		updateFeature: Need(FeatureWrite, InEnv("environment_id")),
		getFeature:    Need(FeatureRead, InEnv("environment_id")),
		cloneFeature: Need(FeatureRead, InEnv("environment_id")).
			And(Need(FeatureWrite, InEnv("target_environment_id"))),
		getInsights: Need(InsightsRead, EachEnv("environment_ids")),
		listAuditLogs: Need(AuditLogRead, InOrg("organization_id").IfSet()).Or(SystemAdmin()).
			And(Need(AuditLogRead, InEnv("environment_id").IfSet())),
		updateAccount:      Need(AccountManage, InOrg("organization_id")).Or(Self("email")).Guard(guard),
		listAccounts:       Need(AccountRead, InOrg("organization_id")),
		getSubscription:    Need(SubscriptionRead, InEnv("environment_id")),
		listNotifications:  Authenticated(),
		invalidRule:        Need(FeatureWrite, InEnv("environment_id")).And(SystemAdminOnly()),
		createOrganization: SystemAdminOnly(),
		ping:               Public(),
	}
}

// fakeAccounts answers by email, whatever the scope, like one organization.
type fakeAccounts struct {
	accounts map[string]*accountproto.AccountV2
	err      error
	calls    int
}

func (f *fakeAccounts) Account(_ context.Context, email string, _ ScopeID) (*accountproto.AccountV2, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	account, ok := f.accounts[email]
	if !ok {
		return nil, status.Error(codes.NotFound, "account not found")
	}
	return account, nil
}

func envRole(env string, role accountproto.AccountV2_Role_Environment) *accountproto.AccountV2_EnvironmentRole {
	return &accountproto.AccountV2_EnvironmentRole{EnvironmentId: env, Role: role}
}

func testAccounts() map[string]*accountproto.AccountV2 {
	viewer := accountproto.AccountV2_Role_Environment_VIEWER
	editor := accountproto.AccountV2_Role_Environment_EDITOR
	return map[string]*accountproto.AccountV2{
		"viewer@example.com": {
			Email:            "viewer@example.com",
			OrganizationId:   "org-1",
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{envRole("env-1", viewer)},
		},
		"editor@example.com": {
			Email:            "editor@example.com",
			Name:             "Editor",
			FirstName:        "Edi",
			LastName:         "Tor",
			OrganizationId:   "org-1",
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{envRole("env-1", editor), envRole("env-2", viewer)},
		},
		"member@example.com": {
			Email:            "member@example.com",
			OrganizationId:   "org-1",
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		},
		"admin@example.com": {
			Email:            "admin@example.com",
			OrganizationId:   "org-1",
			OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
		},
		"disabled@example.com": {
			Email:            "disabled@example.com",
			OrganizationId:   "org-1",
			OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
			Disabled:         true,
		},
	}
}

func userCtx(email string) context.Context {
	return context.WithValue(context.Background(), rpc.AccessTokenKey, &token.AccessToken{Email: email, Name: "Name"})
}

func systemAdminCtx(email string) context.Context {
	return context.WithValue(context.Background(), rpc.AccessTokenKey,
		&token.AccessToken{Email: email, Name: "Admin", IsSystemAdmin: true})
}

func serviceTokenCtx() context.Context {
	return context.WithValue(context.Background(), rpc.AccessTokenKey,
		&token.AccessToken{Email: "service@example.com", Name: "Service", IsSystemAdmin: true, IsServiceToken: true})
}

func relayedCtx(maintainer string) context.Context {
	md := metadata.Pairs(APIKeyTokenMDKey, rawAPIKey, APIKeyMaintainerMDKey, maintainer, APIKeyNameMDKey, "key")
	return metadata.NewIncomingContext(serviceTokenCtx(), md)
}

func TestAuthorize(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc          string
		ctx           context.Context
		method        string
		req           proto.Message
		guardErr      error
		loaderErr     error
		expected      *eventproto.Editor
		expectedErr   error
		expectedCode  codes.Code
		expectedCalls int
	}{
		{
			desc:        "no rule",
			ctx:         userCtx("editor@example.com"),
			method:      "/bucketeer.feature.FeatureService/Unknown",
			req:         &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expectedErr: ErrNoRule,
		},
		{
			desc:   "public without token",
			ctx:    context.Background(),
			method: ping,
			req:    &featureproto.UpdateFeatureRequest{},
		},
		{
			desc:        "no token",
			ctx:         context.Background(),
			method:      updateFeature,
			req:         &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expectedErr: ErrUnauthenticated,
		},
		{
			desc:     "authenticated only",
			ctx:      userCtx("nobody@example.com"),
			method:   listNotifications,
			req:      &featureproto.UpdateFeatureRequest{},
			expected: &eventproto.Editor{Email: "nobody@example.com", Name: "Name"},
		},
		{
			desc:        "system admin only denies a user",
			ctx:         userCtx("admin@example.com"),
			method:      createOrganization,
			req:         &featureproto.UpdateFeatureRequest{},
			expectedErr: ErrPermissionDenied,
		},
		{
			desc:     "system admin only allows a system admin",
			ctx:      systemAdminCtx("root@example.com"),
			method:   createOrganization,
			req:      &featureproto.UpdateFeatureRequest{},
			expected: &eventproto.Editor{Email: "root@example.com", Name: "Admin", IsAdmin: true},
		},
		{
			desc:     "service token has full access",
			ctx:      serviceTokenCtx(),
			method:   updateFeature,
			req:      &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expected: &eventproto.Editor{Email: "service@example.com", Name: "Service", IsAdmin: true},
		},
		{
			desc:        "missing scope",
			ctx:         userCtx("editor@example.com"),
			method:      updateFeature,
			req:         &featureproto.UpdateFeatureRequest{},
			expectedErr: ErrMissingScope,
		},
		{
			desc:   "editor writes its environment",
			ctx:    userCtx("editor@example.com"),
			method: updateFeature,
			req:    &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expected: &eventproto.Editor{
				Email:            "editor@example.com",
				Name:             "Name",
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				EnvironmentRoles: testAccounts()["editor@example.com"].EnvironmentRoles,
			},
			expectedCalls: 1,
		},
		{
			desc:          "viewer cannot write",
			ctx:           userCtx("viewer@example.com"),
			method:        updateFeature,
			req:           &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expectedErr:   ErrPermissionDenied,
			expectedCalls: 1,
		},
		{
			desc:          "viewer reads",
			ctx:           userCtx("viewer@example.com"),
			method:        getFeature,
			req:           &featureproto.GetFeatureRequest{EnvironmentId: "env-1"},
			expected:      memberEditor("viewer@example.com"),
			expectedCalls: 1,
		},
		{
			desc:          "non-member",
			ctx:           userCtx("nobody@example.com"),
			method:        getFeature,
			req:           &featureproto.GetFeatureRequest{EnvironmentId: "env-1"},
			expectedErr:   ErrUnauthenticated,
			expectedCalls: 1,
		},
		{
			desc:          "disabled account",
			ctx:           userCtx("disabled@example.com"),
			method:        getFeature,
			req:           &featureproto.GetFeatureRequest{EnvironmentId: "env-1"},
			expectedErr:   ErrUnauthenticated,
			expectedCalls: 1,
		},
		{
			desc:          "account lookup fails",
			ctx:           userCtx("editor@example.com"),
			method:        getFeature,
			req:           &featureproto.GetFeatureRequest{EnvironmentId: "env-1"},
			loaderErr:     errors.New("db down"),
			expectedErr:   ErrInternal,
			expectedCalls: 1,
		},
		{
			desc:          "organization admin writes without an environment role",
			ctx:           userCtx("admin@example.com"),
			method:        updateFeature,
			req:           &featureproto.UpdateFeatureRequest{EnvironmentId: "env-9"},
			expected:      memberEditor("admin@example.com"),
			expectedCalls: 1,
		},
		{
			desc:          "system admin reads without membership",
			ctx:           systemAdminCtx("root@example.com"),
			method:        getFeature,
			req:           &featureproto.GetFeatureRequest{EnvironmentId: "env-1"},
			expected:      &eventproto.Editor{Email: "root@example.com", Name: "Admin", IsAdmin: true},
			expectedCalls: 0,
		},
		{
			desc:          "system admin writes only as a member",
			ctx:           systemAdminCtx("root@example.com"),
			method:        updateFeature,
			req:           &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expectedErr:   ErrUnauthenticated,
			expectedCalls: 1,
		},
		{
			desc:   "system admin member writes",
			ctx:    systemAdminCtx("editor@example.com"),
			method: updateFeature,
			req:    &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expected: &eventproto.Editor{
				Email:            "editor@example.com",
				Name:             "Admin",
				IsAdmin:          true,
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				EnvironmentRoles: testAccounts()["editor@example.com"].EnvironmentRoles,
			},
			expectedCalls: 1,
		},
		{
			desc:          "organization member without an environment role is denied there",
			ctx:           userCtx("member@example.com"),
			method:        getSubscription,
			req:           &subscriptionproto.GetSubscriptionRequest{EnvironmentId: "env-1"},
			expectedErr:   ErrPermissionDenied,
			expectedCalls: 1,
		},
		{
			desc:          "organization member reads accounts",
			ctx:           userCtx("viewer@example.com"),
			method:        listAccounts,
			req:           &accountproto.ListAccountsV2Request{OrganizationId: "org-1"},
			expected:      memberEditor("viewer@example.com"),
			expectedCalls: 1,
		},
		{
			desc:          "member cannot manage another account",
			ctx:           userCtx("viewer@example.com"),
			method:        updateAccount,
			req:           &accountproto.UpdateAccountV2Request{OrganizationId: "org-1", Email: "editor@example.com"},
			expectedErr:   ErrPermissionDenied,
			expectedCalls: 1,
		},
		{
			desc:          "member updates self through the alternative",
			ctx:           userCtx("viewer@example.com"),
			method:        updateAccount,
			req:           &accountproto.UpdateAccountV2Request{OrganizationId: "org-1", Email: "viewer@example.com"},
			expected:      memberEditor("viewer@example.com"),
			expectedCalls: 1,
		},
		{
			desc:          "admin manages accounts",
			ctx:           userCtx("admin@example.com"),
			method:        updateAccount,
			req:           &accountproto.UpdateAccountV2Request{OrganizationId: "org-1", Email: "editor@example.com"},
			expected:      memberEditor("admin@example.com"),
			expectedCalls: 1,
		},
		{
			desc:        "service token runs the guards",
			ctx:         serviceTokenCtx(),
			method:      updateAccount,
			req:         &accountproto.UpdateAccountV2Request{OrganizationId: "org-1", Email: "editor@example.com"},
			guardErr:    errGuard,
			expectedErr: errGuard,
		},
		{
			desc:     "service token passes the guards without a permission check",
			ctx:      serviceTokenCtx(),
			method:   updateAccount,
			req:      &accountproto.UpdateAccountV2Request{OrganizationId: "org-1", Email: "editor@example.com"},
			expected: &eventproto.Editor{Email: "service@example.com", Name: "Service", IsAdmin: true},
		},
		{
			desc:          "guard rejects",
			ctx:           userCtx("admin@example.com"),
			method:        updateAccount,
			req:           &accountproto.UpdateAccountV2Request{OrganizationId: "org-1", Email: "editor@example.com"},
			guardErr:      errGuard,
			expectedErr:   errGuard,
			expectedCalls: 1,
		},
		{
			desc:          "audit logs of the organization need an admin",
			ctx:           userCtx("viewer@example.com"),
			method:        listAuditLogs,
			req:           &auditlogproto.ListAuditLogsRequest{OrganizationId: "org-1"},
			expectedErr:   ErrPermissionDenied,
			expectedCalls: 1,
		},
		{
			desc:          "audit logs of the organization by an admin",
			ctx:           userCtx("admin@example.com"),
			method:        listAuditLogs,
			req:           &auditlogproto.ListAuditLogsRequest{OrganizationId: "org-1"},
			expected:      memberEditor("admin@example.com"),
			expectedCalls: 1,
		},
		{
			desc:          "audit logs of the organization by a system admin without membership",
			ctx:           systemAdminCtx("root@example.com"),
			method:        listAuditLogs,
			req:           &auditlogproto.ListAuditLogsRequest{OrganizationId: "org-1"},
			expected:      &eventproto.Editor{Email: "root@example.com", Name: "Admin", IsAdmin: true},
			expectedCalls: 0,
		},
		{
			desc:          "audit logs of an environment by a viewer",
			ctx:           userCtx("viewer@example.com"),
			method:        listAuditLogs,
			req:           &auditlogproto.ListAuditLogsRequest{EnvironmentId: "env-1"},
			expected:      memberEditor("viewer@example.com"),
			expectedCalls: 1,
		},
		{
			desc:        "audit logs without any scope",
			ctx:         userCtx("viewer@example.com"),
			method:      listAuditLogs,
			req:         &auditlogproto.ListAuditLogsRequest{},
			expectedErr: ErrMissingScope,
		},
		{
			desc:          "clone needs read on the source and write on the target",
			ctx:           userCtx("editor@example.com"),
			method:        cloneFeature,
			req:           &featureproto.CloneFeatureRequest{EnvironmentId: "env-2", TargetEnvironmentId: "env-1"},
			expected:      memberEditor("editor@example.com"),
			expectedCalls: 2,
		},
		{
			desc:          "clone into an environment without write",
			ctx:           userCtx("editor@example.com"),
			method:        cloneFeature,
			req:           &featureproto.CloneFeatureRequest{EnvironmentId: "env-1", TargetEnvironmentId: "env-2"},
			expectedErr:   ErrPermissionDenied,
			expectedCalls: 2,
		},
		{
			desc:          "every listed environment is checked",
			ctx:           userCtx("editor@example.com"),
			method:        getInsights,
			req:           &insightsproto.GetInsightsMonthlySummaryRequest{EnvironmentIds: []string{"env-1", "env-2"}},
			expected:      memberEditor("editor@example.com"),
			expectedCalls: 2,
		},
		{
			desc:          "one listed environment without a role denies",
			ctx:           userCtx("viewer@example.com"),
			method:        getInsights,
			req:           &insightsproto.GetInsightsMonthlySummaryRequest{EnvironmentIds: []string{"env-1", "env-2"}},
			expectedErr:   ErrPermissionDenied,
			expectedCalls: 2,
		},
		{
			desc:          "a repeated environment is loaded once",
			ctx:           userCtx("viewer@example.com"),
			method:        getInsights,
			req:           &insightsproto.GetInsightsMonthlySummaryRequest{EnvironmentIds: []string{"env-1", "env-1"}},
			expected:      memberEditor("viewer@example.com"),
			expectedCalls: 1,
		},
		{
			desc:   "public API request acts as the key maintainer",
			ctx:    relayedCtx("editor@example.com"),
			method: updateFeature,
			req:    &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expected: &eventproto.Editor{
				Email:            "editor@example.com",
				Name:             "Edi Tor",
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				EnvironmentRoles: testAccounts()["editor@example.com"].EnvironmentRoles,
				PublicApiEditor: &eventproto.Editor_PublicAPIEditor{
					Token: obfuscatedAPIKey, Maintainer: "editor@example.com", Name: "key",
				},
			},
			expectedCalls: 1,
		},
		{
			desc:          "public API request is limited to the maintainer's role",
			ctx:           relayedCtx("viewer@example.com"),
			method:        updateFeature,
			req:           &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expectedErr:   ErrPermissionDenied,
			expectedCalls: 1,
		},
		{
			desc:          "public API request with an unknown maintainer",
			ctx:           relayedCtx("nobody@example.com"),
			method:        updateFeature,
			req:           &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expectedErr:   ErrUnauthenticated,
			expectedCalls: 1,
		},
		{
			desc:         "rule composes a special kind",
			ctx:          systemAdminCtx("editor@example.com"),
			method:       invalidRule,
			req:          &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expectedCode: codes.Internal,
		},
		{
			desc:         "rule names a field the request lacks",
			ctx:          userCtx("editor@example.com"),
			method:       getFeature,
			req:          &accountproto.GetMeRequest{},
			expectedCode: codes.Internal,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			loader := &fakeAccounts{accounts: testAccounts(), err: p.loaderErr}
			guard := guardFunc(func(context.Context, *Access) error { return p.guardErr })
			authorizer := New(testRules(guard), loader)
			editor, err := authorizer.Authorize(p.ctx, p.method, p.req)
			if p.expectedCode != codes.OK {
				require.Error(t, err)
				assert.Equal(t, p.expectedCode, status.Code(err))
			} else {
				assert.Equal(t, p.expectedErr, err)
			}
			assert.Equal(t, p.expected, editor)
			assert.Equal(t, p.expectedCalls, loader.calls)
		})
	}
}

func memberEditor(email string) *eventproto.Editor {
	account := testAccounts()[email]
	return &eventproto.Editor{
		Email:            email,
		Name:             "Name",
		OrganizationRole: account.OrganizationRole,
		EnvironmentRoles: account.EnvironmentRoles,
	}
}

func TestInvalidCompositions(t *testing.T) {
	t.Parallel()
	need := Need(FeatureWrite, InEnv("environment_id"))
	guard := guardFunc(func(context.Context, *Access) error { return nil })
	patterns := []struct {
		desc string
		rule Requirement
	}{
		{desc: "And on Public", rule: Public().And(need)},
		{desc: "And with SystemAdminOnly", rule: need.And(SystemAdminOnly())},
		{desc: "Or on Authenticated", rule: Authenticated().Or(SystemAdmin())},
		{desc: "Guard on SystemAdminOnly", rule: SystemAdminOnly().Guard(guard)},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.ErrorIs(t, p.rule.err, ErrRule)
		})
	}
	assert.NoError(t, need.Or(Self("email")).And(need).Guard(guard).err)
}

// scopedAccounts answers per scope, like a caller who is a member of two
// organizations with different roles.
type scopedAccounts map[ScopeID]*accountproto.AccountV2

func (s scopedAccounts) Account(_ context.Context, _ string, scope ScopeID) (*accountproto.AccountV2, error) {
	if account, ok := s[scope]; ok {
		return account, nil
	}
	return nil, status.Error(codes.NotFound, "account not found")
}

func TestAuthorizeEditorOfTheTargetScope(t *testing.T) {
	t.Parallel()
	viewer := accountproto.AccountV2_Role_Environment_VIEWER
	editor := accountproto.AccountV2_Role_Environment_EDITOR
	source := &accountproto.AccountV2{
		Email:            "clone@example.com",
		OrganizationId:   "org-1",
		OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{envRole("env-1", viewer)},
	}
	target := &accountproto.AccountV2{
		Email:            "clone@example.com",
		OrganizationId:   "org-2",
		OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{envRole("env-2", editor)},
	}
	loader := scopedAccounts{
		{Kind: EnvironmentScope, ID: "env-1"}: source,
		{Kind: EnvironmentScope, ID: "env-2"}: target,
	}
	authorizer := New(testRules(guardFunc(func(context.Context, *Access) error { return nil })), loader)
	actual, err := authorizer.Authorize(
		userCtx("clone@example.com"),
		cloneFeature,
		&featureproto.CloneFeatureRequest{EnvironmentId: "env-1", TargetEnvironmentId: "env-2"},
	)
	require.NoError(t, err)
	assert.Equal(t, &eventproto.Editor{
		Email:            "clone@example.com",
		Name:             "Name",
		OrganizationRole: target.OrganizationRole,
		EnvironmentRoles: target.EnvironmentRoles,
	}, actual)
}

func TestAccessAccountIsCached(t *testing.T) {
	t.Parallel()
	loader := &fakeAccounts{accounts: testAccounts()}
	access := &Access{Caller: &Caller{Email: "editor@example.com"}, accounts: &accountCache{loader: loader}}
	scope := ScopeID{Kind: EnvironmentScope, ID: "env-1"}
	first, err := access.Account(context.Background(), scope)
	require.NoError(t, err)
	second, err := access.Account(context.Background(), scope)
	require.NoError(t, err)
	assert.Same(t, first, second)
	assert.Equal(t, 1, loader.calls)
}

func TestEditorContext(t *testing.T) {
	t.Parallel()
	assert.Nil(t, EditorFrom(context.Background()))
	editor := &eventproto.Editor{Email: "editor@example.com"}
	assert.Same(t, editor, EditorFrom(WithEditor(context.Background(), editor)))
}
