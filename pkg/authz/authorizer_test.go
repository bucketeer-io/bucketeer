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

	accdomain "github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
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

// fakeAccounts answers by email; every environment belongs to org-1.
type fakeAccounts struct {
	accounts     map[string]*accountproto.AccountV2
	systemAdmins map[string]bool
	err          error
	calls        int
}

func (f *fakeAccounts) Account(_ context.Context, email string, scope ScopeID) (*Account, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	organizationID := scope.ID
	if scope.Kind == EnvironmentScope {
		organizationID = "org-1"
	}
	account, ok := f.accounts[email]
	if !ok || account.OrganizationId != organizationID {
		return nil, status.Error(codes.NotFound, "account not found")
	}
	return &Account{AccountV2: account, SystemAdmin: f.systemAdmins[email]}, nil
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
			Name:             "Viewer",
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
			Name:             "Member",
			OrganizationId:   "org-1",
			OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		},
		"admin@example.com": {
			Email:            "admin@example.com",
			Name:             "Admin",
			OrganizationId:   "org-1",
			OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
		},
		"disabled@example.com": {
			Email:            "disabled@example.com",
			Name:             "Disabled",
			OrganizationId:   "org-1",
			OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
			Disabled:         true,
		},
		"root@example.com": {
			Email:            "root@example.com",
			Name:             "Root",
			OrganizationId:   "org-admin",
			OrganizationRole: accountproto.AccountV2_Role_Organization_OWNER,
		},
	}
}

func user(email string) Caller {
	return Caller{Email: email, OrganizationID: "org-1"}
}

func root() Caller {
	return Caller{Email: "root@example.com", OrganizationID: "org-admin"}
}

func service() Caller {
	return Caller{Email: "service@example.com", Service: true}
}

func relayed(maintainer string) Caller {
	return Caller{Email: maintainer, PublicAPI: publicAPIEditor(maintainer)}
}

func publicAPIEditor(maintainer string) *eventproto.Editor_PublicAPIEditor {
	return &eventproto.Editor_PublicAPIEditor{Token: obfuscatedAPIKey, Maintainer: maintainer, Name: "key"}
}

func rootEditor() *eventproto.Editor {
	return &eventproto.Editor{Email: "root@example.com", Name: "Root", IsAdmin: true}
}

func TestAuthorize(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc          string
		caller        Caller
		systemAdmin   bool
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
			caller:      user("editor@example.com"),
			method:      "/bucketeer.feature.FeatureService/Unknown",
			req:         &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expectedErr: ErrNoRule,
		},
		{
			desc:   "public without caller",
			method: ping,
			req:    &featureproto.UpdateFeatureRequest{},
		},
		{
			desc:        "no caller",
			method:      updateFeature,
			req:         &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expectedErr: ErrUnauthenticated,
		},
		{
			desc:          "authenticated only",
			caller:        user("viewer@example.com"),
			method:        listNotifications,
			req:           &featureproto.UpdateFeatureRequest{},
			expected:      &eventproto.Editor{Email: "viewer@example.com", Name: "Viewer"},
			expectedCalls: 1,
		},
		{
			desc:          "authenticated only without an account",
			caller:        user("nobody@example.com"),
			method:        listNotifications,
			req:           &featureproto.UpdateFeatureRequest{},
			expectedErr:   ErrUnauthenticated,
			expectedCalls: 1,
		},
		{
			desc:          "system admin only denies a user",
			caller:        user("admin@example.com"),
			method:        createOrganization,
			req:           &featureproto.UpdateFeatureRequest{},
			expectedErr:   ErrPermissionDenied,
			expectedCalls: 1,
		},
		{
			desc:          "system admin only allows a system admin",
			caller:        root(),
			systemAdmin:   true,
			method:        createOrganization,
			req:           &featureproto.UpdateFeatureRequest{},
			expected:      rootEditor(),
			expectedCalls: 1,
		},
		{
			desc:     "system admin only allows a service identity",
			caller:   service(),
			method:   createOrganization,
			req:      &featureproto.UpdateFeatureRequest{},
			expected: &eventproto.Editor{Email: "service@example.com", IsAdmin: true},
		},
		{
			desc:     "service identity has full access",
			caller:   service(),
			method:   updateFeature,
			req:      &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expected: &eventproto.Editor{Email: "service@example.com", IsAdmin: true},
		},
		{
			desc:          "missing scope",
			caller:        user("editor@example.com"),
			method:        updateFeature,
			req:           &featureproto.UpdateFeatureRequest{},
			expectedErr:   ErrMissingScope,
			expectedCalls: 1,
		},
		{
			desc:          "editor writes its environment",
			caller:        user("editor@example.com"),
			method:        updateFeature,
			req:           &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expected:      memberEditor("editor@example.com"),
			expectedCalls: 2,
		},
		{
			desc:          "viewer cannot write",
			caller:        user("viewer@example.com"),
			method:        updateFeature,
			req:           &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expectedErr:   ErrPermissionDenied,
			expectedCalls: 2,
		},
		{
			desc:          "viewer reads",
			caller:        user("viewer@example.com"),
			method:        getFeature,
			req:           &featureproto.GetFeatureRequest{EnvironmentId: "env-1"},
			expected:      memberEditor("viewer@example.com"),
			expectedCalls: 2,
		},
		{
			desc:          "non-member",
			caller:        user("nobody@example.com"),
			method:        getFeature,
			req:           &featureproto.GetFeatureRequest{EnvironmentId: "env-1"},
			expectedErr:   ErrUnauthenticated,
			expectedCalls: 1,
		},
		{
			desc:          "disabled account",
			caller:        user("disabled@example.com"),
			method:        getFeature,
			req:           &featureproto.GetFeatureRequest{EnvironmentId: "env-1"},
			expectedErr:   ErrUnauthenticated,
			expectedCalls: 1,
		},
		{
			desc:          "account lookup fails",
			caller:        user("editor@example.com"),
			method:        getFeature,
			req:           &featureproto.GetFeatureRequest{EnvironmentId: "env-1"},
			loaderErr:     errors.New("db down"),
			expectedErr:   ErrInternal,
			expectedCalls: 1,
		},
		{
			desc:          "organization admin writes without an environment role",
			caller:        user("admin@example.com"),
			method:        updateFeature,
			req:           &featureproto.UpdateFeatureRequest{EnvironmentId: "env-9"},
			expected:      memberEditor("admin@example.com"),
			expectedCalls: 2,
		},
		{
			desc:          "system admin reads without membership",
			caller:        root(),
			systemAdmin:   true,
			method:        getFeature,
			req:           &featureproto.GetFeatureRequest{EnvironmentId: "env-1"},
			expected:      rootEditor(),
			expectedCalls: 1,
		},
		{
			desc:          "system admin writes only as a member",
			caller:        root(),
			systemAdmin:   true,
			method:        updateFeature,
			req:           &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expectedErr:   ErrUnauthenticated,
			expectedCalls: 2,
		},
		{
			desc:        "system admin member writes",
			caller:      user("editor@example.com"),
			systemAdmin: true,
			method:      updateFeature,
			req:         &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expected: &eventproto.Editor{
				Email:            "editor@example.com",
				Name:             "Edi Tor",
				IsAdmin:          true,
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				EnvironmentRoles: testAccounts()["editor@example.com"].EnvironmentRoles,
			},
			expectedCalls: 2,
		},
		{
			desc:          "organization member without an environment role is denied there",
			caller:        user("member@example.com"),
			method:        getSubscription,
			req:           &subscriptionproto.GetSubscriptionRequest{EnvironmentId: "env-1"},
			expectedErr:   ErrPermissionDenied,
			expectedCalls: 2,
		},
		{
			desc:          "organization member reads accounts",
			caller:        user("viewer@example.com"),
			method:        listAccounts,
			req:           &accountproto.ListAccountsV2Request{OrganizationId: "org-1"},
			expected:      memberEditor("viewer@example.com"),
			expectedCalls: 1,
		},
		{
			desc:          "member cannot manage another account",
			caller:        user("viewer@example.com"),
			method:        updateAccount,
			req:           &accountproto.UpdateAccountV2Request{OrganizationId: "org-1", Email: "editor@example.com"},
			expectedErr:   ErrPermissionDenied,
			expectedCalls: 1,
		},
		{
			desc:          "member updates self through the alternative",
			caller:        user("viewer@example.com"),
			method:        updateAccount,
			req:           &accountproto.UpdateAccountV2Request{OrganizationId: "org-1", Email: "viewer@example.com"},
			expected:      memberEditor("viewer@example.com"),
			expectedCalls: 1,
		},
		{
			desc:          "admin manages accounts",
			caller:        user("admin@example.com"),
			method:        updateAccount,
			req:           &accountproto.UpdateAccountV2Request{OrganizationId: "org-1", Email: "editor@example.com"},
			expected:      memberEditor("admin@example.com"),
			expectedCalls: 1,
		},
		{
			desc:        "service identity runs the guards",
			caller:      service(),
			method:      updateAccount,
			req:         &accountproto.UpdateAccountV2Request{OrganizationId: "org-1", Email: "editor@example.com"},
			guardErr:    errGuard,
			expectedErr: errGuard,
		},
		{
			desc:     "service identity passes the guards without a permission check",
			caller:   service(),
			method:   updateAccount,
			req:      &accountproto.UpdateAccountV2Request{OrganizationId: "org-1", Email: "editor@example.com"},
			expected: &eventproto.Editor{Email: "service@example.com", IsAdmin: true},
		},
		{
			desc:          "guard rejects",
			caller:        user("admin@example.com"),
			method:        updateAccount,
			req:           &accountproto.UpdateAccountV2Request{OrganizationId: "org-1", Email: "editor@example.com"},
			guardErr:      errGuard,
			expectedErr:   errGuard,
			expectedCalls: 1,
		},
		{
			desc:          "audit logs of the organization need an admin",
			caller:        user("viewer@example.com"),
			method:        listAuditLogs,
			req:           &auditlogproto.ListAuditLogsRequest{OrganizationId: "org-1"},
			expectedErr:   ErrPermissionDenied,
			expectedCalls: 1,
		},
		{
			desc:          "audit logs of the organization by an admin",
			caller:        user("admin@example.com"),
			method:        listAuditLogs,
			req:           &auditlogproto.ListAuditLogsRequest{OrganizationId: "org-1"},
			expected:      memberEditor("admin@example.com"),
			expectedCalls: 1,
		},
		{
			desc:          "audit logs of the organization by a system admin without membership",
			caller:        root(),
			systemAdmin:   true,
			method:        listAuditLogs,
			req:           &auditlogproto.ListAuditLogsRequest{OrganizationId: "org-1"},
			expected:      rootEditor(),
			expectedCalls: 1,
		},
		{
			desc:          "audit logs of an environment by a viewer",
			caller:        user("viewer@example.com"),
			method:        listAuditLogs,
			req:           &auditlogproto.ListAuditLogsRequest{EnvironmentId: "env-1"},
			expected:      memberEditor("viewer@example.com"),
			expectedCalls: 2,
		},
		{
			desc:          "audit logs without any scope",
			caller:        user("viewer@example.com"),
			method:        listAuditLogs,
			req:           &auditlogproto.ListAuditLogsRequest{},
			expectedErr:   ErrMissingScope,
			expectedCalls: 1,
		},
		{
			desc:          "clone needs read on the source and write on the target",
			caller:        user("editor@example.com"),
			method:        cloneFeature,
			req:           &featureproto.CloneFeatureRequest{EnvironmentId: "env-2", TargetEnvironmentId: "env-1"},
			expected:      memberEditor("editor@example.com"),
			expectedCalls: 3,
		},
		{
			desc:          "clone into an environment without write",
			caller:        user("editor@example.com"),
			method:        cloneFeature,
			req:           &featureproto.CloneFeatureRequest{EnvironmentId: "env-1", TargetEnvironmentId: "env-2"},
			expectedErr:   ErrPermissionDenied,
			expectedCalls: 3,
		},
		{
			desc:          "every listed environment is checked",
			caller:        user("editor@example.com"),
			method:        getInsights,
			req:           &insightsproto.GetInsightsMonthlySummaryRequest{EnvironmentIds: []string{"env-1", "env-2"}},
			expected:      memberEditor("editor@example.com"),
			expectedCalls: 3,
		},
		{
			desc:          "one listed environment without a role denies",
			caller:        user("viewer@example.com"),
			method:        getInsights,
			req:           &insightsproto.GetInsightsMonthlySummaryRequest{EnvironmentIds: []string{"env-1", "env-2"}},
			expectedErr:   ErrPermissionDenied,
			expectedCalls: 3,
		},
		{
			desc:          "a repeated environment is loaded once",
			caller:        user("viewer@example.com"),
			method:        getInsights,
			req:           &insightsproto.GetInsightsMonthlySummaryRequest{EnvironmentIds: []string{"env-1", "env-1"}},
			expected:      memberEditor("viewer@example.com"),
			expectedCalls: 2,
		},
		{
			desc:   "public API request acts as the key maintainer",
			caller: relayed("editor@example.com"),
			method: updateFeature,
			req:    &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expected: &eventproto.Editor{
				Email:            "editor@example.com",
				Name:             "Edi Tor",
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				EnvironmentRoles: testAccounts()["editor@example.com"].EnvironmentRoles,
				PublicApiEditor:  publicAPIEditor("editor@example.com"),
			},
			expectedCalls: 1,
		},
		{
			desc:          "public API request is limited to the maintainer's role",
			caller:        relayed("viewer@example.com"),
			method:        updateFeature,
			req:           &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expectedErr:   ErrPermissionDenied,
			expectedCalls: 1,
		},
		{
			desc:          "public API request with an unknown maintainer",
			caller:        relayed("nobody@example.com"),
			method:        updateFeature,
			req:           &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expectedErr:   ErrUnauthenticated,
			expectedCalls: 1,
		},
		{
			desc:         "rule composes a special kind",
			caller:       user("editor@example.com"),
			method:       invalidRule,
			req:          &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expectedCode: codes.Internal,
		},
		{
			desc:          "rule names a field the request lacks",
			caller:        user("editor@example.com"),
			method:        getFeature,
			req:           &accountproto.GetMeRequest{},
			expectedCode:  codes.Internal,
			expectedCalls: 1,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			loader := &fakeAccounts{
				accounts:     testAccounts(),
				systemAdmins: map[string]bool{p.caller.Email: p.systemAdmin},
				err:          p.loaderErr,
			}
			guard := guardFunc(func(context.Context, *Access) error { return p.guardErr })
			authorizer := New(testRules(guard), loader)
			editor, err := authorizer.Authorize(context.Background(), p.caller, p.method, p.req)
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
		Name:             (&accdomain.AccountV2{AccountV2: account}).GetAccountFullName(),
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

func (s scopedAccounts) Account(_ context.Context, _ string, scope ScopeID) (*Account, error) {
	if account, ok := s[scope]; ok {
		return &Account{AccountV2: account}, nil
	}
	return nil, status.Error(codes.NotFound, "account not found")
}

func TestAuthorizeEditorOfTheTargetScope(t *testing.T) {
	t.Parallel()
	viewer := accountproto.AccountV2_Role_Environment_VIEWER
	editor := accountproto.AccountV2_Role_Environment_EDITOR
	source := &accountproto.AccountV2{
		Email:            "clone@example.com",
		Name:             "Clone",
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
		{Kind: OrganizationScope, ID: "org-1"}: source,
		{Kind: EnvironmentScope, ID: "env-1"}:  source,
		{Kind: EnvironmentScope, ID: "env-2"}:  target,
	}
	authorizer := New(testRules(guardFunc(func(context.Context, *Access) error { return nil })), loader)
	actual, err := authorizer.Authorize(
		context.Background(),
		user("clone@example.com"),
		cloneFeature,
		&featureproto.CloneFeatureRequest{EnvironmentId: "env-1", TargetEnvironmentId: "env-2"},
	)
	require.NoError(t, err)
	assert.Equal(t, &eventproto.Editor{
		Email:            "clone@example.com",
		Name:             "Clone",
		OrganizationRole: target.OrganizationRole,
		EnvironmentRoles: target.EnvironmentRoles,
	}, actual)
}

func TestAccessAccountIsCached(t *testing.T) {
	t.Parallel()
	loader := &fakeAccounts{accounts: testAccounts()}
	access := &Access{Caller: Caller{Email: "editor@example.com"}, accounts: &accountCache{loader: loader}}
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

func TestCallerRelayed(t *testing.T) {
	t.Parallel()
	md := metadata.Pairs(APIKeyTokenMDKey, rawAPIKey, APIKeyMaintainerMDKey, "editor@example.com", APIKeyNameMDKey, "key")
	withKey := metadata.NewIncomingContext(context.Background(), md)
	patterns := []struct {
		desc     string
		ctx      context.Context
		caller   Caller
		expected Caller
	}{
		{desc: "service with an API key", ctx: withKey, caller: service(), expected: relayed("editor@example.com")},
		{desc: "service without an API key", ctx: context.Background(), caller: service(), expected: service()},
		{
			desc: "service with empty metadata", caller: service(), expected: service(),
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs(APIKeyMaintainerMDKey, "x")),
		},
		{desc: "user metadata is ignored", ctx: withKey, caller: user("viewer@example.com"), expected: user("viewer@example.com")},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expected, p.caller.Relayed(p.ctx))
		})
	}
}

func TestCallerFullAccess(t *testing.T) {
	t.Parallel()
	assert.True(t, service().FullAccess())
	assert.False(t, relayed("editor@example.com").FullAccess())
	assert.False(t, user("editor@example.com").FullAccess())
}
