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
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

func TestGrants(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc       string
		role       Role
		permission Permission
		expected   bool
	}{
		{desc: "viewer reads", role: EnvironmentViewer, permission: FeatureRead, expected: true},
		{desc: "viewer does not write", role: EnvironmentViewer, permission: FeatureWrite, expected: false},
		{desc: "viewer manages own search filters", role: EnvironmentViewer, permission: SearchFilterWrite, expected: true},
		{desc: "editor inherits read", role: EnvironmentEditor, permission: FeatureRead, expected: true},
		{desc: "editor writes", role: EnvironmentEditor, permission: FeatureWrite, expected: true},
		{desc: "editor does not manage accounts", role: EnvironmentEditor, permission: AccountManage, expected: false},
		{desc: "member reads accounts", role: OrganizationMember, permission: AccountRead, expected: true},
		{desc: "member has no environment permission", role: OrganizationMember, permission: FeatureRead, expected: false},
		{desc: "admin inherits member", role: OrganizationAdmin, permission: AccountRead, expected: true},
		{desc: "admin inherits editor", role: OrganizationAdmin, permission: FeatureWrite, expected: true},
		{desc: "admin manages", role: OrganizationAdmin, permission: EnvironmentManage, expected: true},
		{desc: "admin does not manage organization", role: OrganizationAdmin, permission: OrganizationManage, expected: false},
		{desc: "owner inherits admin", role: OrganizationOwner, permission: AccountManage, expected: true},
		{desc: "owner manages organization", role: OrganizationOwner, permission: OrganizationManage, expected: true},
		{desc: "unknown role", role: Role("nobody"), permission: FeatureRead, expected: false},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expected, Grants(p.role, p.permission))
		})
	}
}

// TestPermissionsAreScopedAndGranted fails when a permission is added to
// permission.go without a scope or without a role in Table 1, or when it is
// granted only by roles that do not count in its scope, so a resource cannot
// be added or moved without updating both lists and Table 1.
func TestPermissionsAreScopedAndGranted(t *testing.T) {
	t.Parallel()
	declared := declaredPermissions(t)
	grantedBy := make(map[Permission][]Role)
	for role, list := range grants {
		for _, p := range list {
			grantedBy[p] = append(grantedBy[p], role)
		}
	}
	for _, p := range declared {
		assert.NotEmpty(t, grantedBy[p], "%s is declared but no role grants it", p)
		assert.True(t, inScope(p, EnvironmentScope) || inScope(p, OrganizationScope),
			"%s is declared but in neither scope list", p)
	}
	for p, roles := range grantedBy {
		assert.Contains(t, declared, p, "%s is granted but not declared", p)
		for _, r := range roles {
			if (r == EnvironmentViewer || r == EnvironmentEditor) && !inScope(p, EnvironmentScope) {
				assert.Fail(t, "scope", "%s grants %s, which is not environment-scoped", r, p)
			}
			if r == OrganizationMember && !inScope(p, OrganizationScope) {
				assert.Fail(t, "scope", "%s grants %s, which is not organization-scoped", r, p)
			}
		}
	}
	for _, p := range environmentPermissions {
		assert.Contains(t, declared, p, "%s is environment-scoped but not declared", p)
		assert.True(t, grantedInEnvironment(grantedBy[p]),
			"%s is environment-scoped but only %v grant it", p, grantedBy[p])
	}
	for _, p := range organizationPermissions {
		assert.Contains(t, declared, p, "%s is organization-scoped but not declared", p)
		assert.True(t, grantedInOrganization(grantedBy[p]),
			"%s is organization-scoped but only %v grant it", p, grantedBy[p])
	}
}

func grantedInEnvironment(roles []Role) bool {
	for _, r := range roles {
		if r == EnvironmentViewer || r == EnvironmentEditor || inheritsRole(r, EnvironmentViewer) {
			return true
		}
	}
	return false
}

func grantedInOrganization(roles []Role) bool {
	for _, r := range roles {
		if r == OrganizationMember || inheritsRole(r, OrganizationMember) {
			return true
		}
	}
	return false
}

// declaredPermissions reads every Permission constant of permission.go.
func declaredPermissions(t *testing.T) []Permission {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), "permission.go", nil, 0)
	require.NoError(t, err)
	var declared []Permission
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			vs := spec.(*ast.ValueSpec)
			if ident, ok := vs.Type.(*ast.Ident); !ok || ident.Name != "Permission" {
				continue
			}
			for _, v := range vs.Values {
				value, err := strconv.Unquote(v.(*ast.BasicLit).Value)
				require.NoError(t, err)
				declared = append(declared, Permission(value))
			}
		}
	}
	require.NotEmpty(t, declared)
	return declared
}

func TestOwnerHoldsEverything(t *testing.T) {
	t.Parallel()
	for role := range grants {
		for _, p := range Permissions(role) {
			assert.True(t, Grants(OrganizationOwner, p), "owner lacks %s of %s", p, role)
		}
	}
}

func TestGranted(t *testing.T) {
	t.Parallel()
	env := ScopeID{Kind: EnvironmentScope, ID: "env-1"}
	otherEnv := ScopeID{Kind: EnvironmentScope, ID: "env-2"}
	org := ScopeID{Kind: OrganizationScope, ID: "org-1"}
	editor := &accountproto.AccountV2{
		OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{EnvironmentId: "env-1", Role: accountproto.AccountV2_Role_Environment_EDITOR},
		},
	}
	member := &accountproto.AccountV2{OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER}
	admin := &accountproto.AccountV2{OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN}
	unassigned := &accountproto.AccountV2{
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{EnvironmentId: "env-1", Role: accountproto.AccountV2_Role_Environment_UNASSIGNED},
		},
	}
	patterns := []struct {
		desc       string
		account    *accountproto.AccountV2
		scope      ScopeID
		permission Permission
		expected   bool
	}{
		{desc: "editor writes its environment", account: editor, scope: env, permission: FeatureWrite, expected: true},
		{desc: "editor has no role elsewhere", account: editor, scope: otherEnv, permission: FeatureRead, expected: false},
		{desc: "member reads organization", account: editor, scope: org, permission: AccountRead, expected: true},
		{desc: "member does not manage organization", account: editor, scope: org, permission: AccountManage, expected: false},
		{desc: "environment role does not count in organization scope",
			account: editor, scope: org, permission: FeatureWrite, expected: false},
		{desc: "member reads tags of the organization", account: member, scope: org, permission: TagRead, expected: true},
		{desc: "member without an environment role reads no tags there",
			account: member, scope: env, permission: TagRead, expected: false},
		{desc: "admin writes every environment", account: admin, scope: otherEnv, permission: FeatureWrite, expected: true},
		{desc: "admin manages organization scope", account: admin, scope: org, permission: AccountManage, expected: true},
		{desc: "environment permission is never granted in organization scope",
			account: admin, scope: org, permission: FeatureWrite, expected: false},
		{desc: "organization permission is never granted in environment scope",
			account: admin, scope: env, permission: AccountManage, expected: false},
		{desc: "unassigned environment role grants nothing",
			account: unassigned, scope: env, permission: FeatureRead, expected: false},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expected, granted(p.account, p.scope, p.permission))
		})
	}
}

func TestInheritsRole(t *testing.T) {
	t.Parallel()
	assert.True(t, inheritsRole(OrganizationOwner, EnvironmentViewer))
	assert.True(t, inheritsRole(OrganizationAdmin, EnvironmentViewer))
	assert.False(t, inheritsRole(OrganizationMember, EnvironmentViewer))
	assert.False(t, inheritsRole(EnvironmentViewer, EnvironmentViewer))
}

func TestReadableBySystemAdmin(t *testing.T) {
	t.Parallel()
	assert.True(t, readableBySystemAdmin(FeatureRead))
	assert.True(t, readableBySystemAdmin(AccountRead))
	assert.False(t, readableBySystemAdmin(FeatureWrite))
	assert.False(t, readableBySystemAdmin(AccountManage))
}
