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

// Package authz decides whether a caller may invoke an RPC (RFC 0049).
// Table 1 (this file) says which permissions each role grants, and Table 2
// (Rules) says which permission each RPC needs and in which scope.
package authz

import (
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

// Permission is "<resource>:<verb>". The verbs mirror the roles the checks
// are migrated from: read (VIEWER / MEMBER), write (EDITOR) and manage (ADMIN).
type Permission string

const (
	// Environment-scoped resources.
	FeatureRead              Permission = "feature:read"
	FeatureWrite             Permission = "feature:write"
	SegmentRead              Permission = "segment:read"
	SegmentWrite             Permission = "segment:write"
	ScheduledFlagChangeRead  Permission = "scheduled_flag_change:read"
	ScheduledFlagChangeWrite Permission = "scheduled_flag_change:write"
	HoldoutGroupRead         Permission = "holdout_group:read"
	HoldoutGroupManage       Permission = "holdout_group:manage"
	AutoOpsRuleRead          Permission = "autoops_rule:read"
	AutoOpsRuleWrite         Permission = "autoops_rule:write"
	ProgressiveRolloutRead   Permission = "progressive_rollout:read"
	ProgressiveRolloutWrite  Permission = "progressive_rollout:write"
	ExperimentRead           Permission = "experiment:read"
	ExperimentWrite          Permission = "experiment:write"
	GoalRead                 Permission = "goal:read"
	GoalWrite                Permission = "goal:write"
	EventCounterRead         Permission = "event_counter:read"
	InsightsRead             Permission = "insights:read"
	CodeReferenceRead        Permission = "code_reference:read"
	CodeReferenceWrite       Permission = "code_reference:write"
	TagRead                  Permission = "tag:read"
	TagWrite                 Permission = "tag:write"
	PushRead                 Permission = "push:read"
	PushWrite                Permission = "push:write"
	SubscriptionRead         Permission = "subscription:read"
	SubscriptionWrite        Permission = "subscription:write"
	SearchFilterWrite        Permission = "search_filter:write"
	AIChatRead               Permission = "ai_chat:read"
	UserAttributeRead        Permission = "user_attribute:read"
	// Audit logs of an environment need an environment role; audit logs of a
	// whole organization need the organization admin role.
	AuditLogRead Permission = "audit_log:read"

	// Organization-scoped resources.
	AccountRead        Permission = "account:read"
	AccountManage      Permission = "account:manage"
	APIKeyRead         Permission = "api_key:read"
	APIKeyManage       Permission = "api_key:manage"
	EnvironmentRead    Permission = "environment:read"
	EnvironmentManage  Permission = "environment:manage"
	ProjectRead        Permission = "project:read"
	ProjectManage      Permission = "project:manage"
	OrganizationRead   Permission = "organization:read"
	OrganizationManage Permission = "organization:manage"
	TeamRead           Permission = "team:read"
	TeamManage         Permission = "team:manage"
)

// Role is one of the account roles. Environment roles apply to one
// environment; organization roles apply to every environment of the
// organization, as pkg/role behaves today.
type Role string

const (
	EnvironmentViewer  Role = "environment:VIEWER"
	EnvironmentEditor  Role = "environment:EDITOR"
	OrganizationMember Role = "organization:MEMBER"
	OrganizationAdmin  Role = "organization:ADMIN"
	OrganizationOwner  Role = "organization:OWNER"
)

// The scope each permission is checked in: a permission asked in the other
// scope is never granted. A resource that is listed by organization as well
// as by environment (tags, pushes, subscriptions, API keys, audit logs) is in
// both. TestPermissionsAreScopedAndGranted keeps these lists and Table 1
// complete when a permission is added or moved.
var (
	environmentPermissions = []Permission{
		FeatureRead, FeatureWrite, SegmentRead, SegmentWrite,
		ScheduledFlagChangeRead, ScheduledFlagChangeWrite,
		HoldoutGroupRead, HoldoutGroupManage, AutoOpsRuleRead, AutoOpsRuleWrite,
		ProgressiveRolloutRead, ProgressiveRolloutWrite, ExperimentRead, ExperimentWrite,
		GoalRead, GoalWrite, EventCounterRead, InsightsRead,
		CodeReferenceRead, CodeReferenceWrite, TagRead, TagWrite, PushRead, PushWrite,
		SubscriptionRead, SubscriptionWrite, SearchFilterWrite, AIChatRead,
		UserAttributeRead, AuditLogRead, APIKeyRead,
	}
	organizationPermissions = []Permission{
		AccountRead, AccountManage, APIKeyRead, APIKeyManage,
		EnvironmentRead, EnvironmentManage, ProjectRead, ProjectManage,
		OrganizationRead, OrganizationManage, TeamRead, TeamManage,
		TagRead, PushRead, SubscriptionRead, AuditLogRead,
	}
	scoped = map[ScopeKind]map[Permission]struct{}{
		EnvironmentScope:  toSet(environmentPermissions),
		OrganizationScope: toSet(organizationPermissions),
	}
)

func toSet(list []Permission) map[Permission]struct{} {
	set := make(map[Permission]struct{}, len(list))
	for _, p := range list {
		set[p] = struct{}{}
	}
	return set
}

func inScope(permission Permission, kind ScopeKind) bool {
	_, ok := scoped[kind][permission]
	return ok
}

// Table 1: what each role grants. A role inherits every permission of the
// roles listed in inherits, so the entries only name what the role adds.
var (
	grants = map[Role][]Permission{
		EnvironmentViewer: {
			FeatureRead, SegmentRead, ScheduledFlagChangeRead, HoldoutGroupRead,
			AutoOpsRuleRead, ProgressiveRolloutRead, ExperimentRead, GoalRead,
			EventCounterRead, InsightsRead, CodeReferenceRead, TagRead, PushRead,
			SubscriptionRead, AuditLogRead, AIChatRead, UserAttributeRead, APIKeyRead,
			SearchFilterWrite,
		},
		EnvironmentEditor: {
			FeatureWrite, SegmentWrite, ScheduledFlagChangeWrite, AutoOpsRuleWrite,
			ProgressiveRolloutWrite, ExperimentWrite, GoalWrite, CodeReferenceWrite,
			TagWrite, PushWrite, SubscriptionWrite,
		},
		OrganizationMember: {
			AccountRead, APIKeyRead, EnvironmentRead, ProjectRead, OrganizationRead,
			TeamRead, TagRead, PushRead, SubscriptionRead,
		},
		OrganizationAdmin: {
			AccountManage, APIKeyManage, EnvironmentManage, ProjectManage, TeamManage,
			HoldoutGroupManage, AuditLogRead,
		},
		OrganizationOwner: {
			OrganizationManage,
		},
	}
	inherits = map[Role][]Role{
		EnvironmentEditor:  {EnvironmentViewer},
		OrganizationAdmin:  {OrganizationMember, EnvironmentEditor},
		OrganizationOwner:  {OrganizationAdmin},
		OrganizationMember: nil,
		EnvironmentViewer:  nil,
	}
	permissions = expandGrants()
)

func expandGrants() map[Role]map[Permission]struct{} {
	expanded := make(map[Role]map[Permission]struct{}, len(grants))
	var expand func(Role) map[Permission]struct{}
	expand = func(role Role) map[Permission]struct{} {
		if set, ok := expanded[role]; ok {
			return set
		}
		set := make(map[Permission]struct{})
		for _, p := range grants[role] {
			set[p] = struct{}{}
		}
		for _, parent := range inherits[role] {
			for p := range expand(parent) {
				set[p] = struct{}{}
			}
		}
		expanded[role] = set
		return set
	}
	for role := range grants {
		expand(role)
	}
	return expanded
}

// Grants reports whether role holds permission, directly or by inheritance.
func Grants(role Role, permission Permission) bool {
	_, ok := permissions[role][permission]
	return ok
}

// Permissions returns every permission role holds, in no particular order.
func Permissions(role Role) []Permission {
	set := permissions[role]
	list := make([]Permission, 0, len(set))
	for p := range set {
		list = append(list, p)
	}
	return list
}

func organizationRole(role accountproto.AccountV2_Role_Organization) (Role, bool) {
	switch role {
	case accountproto.AccountV2_Role_Organization_MEMBER:
		return OrganizationMember, true
	case accountproto.AccountV2_Role_Organization_ADMIN:
		return OrganizationAdmin, true
	case accountproto.AccountV2_Role_Organization_OWNER:
		return OrganizationOwner, true
	default:
		return "", false
	}
}

func environmentRole(roles []*accountproto.AccountV2_EnvironmentRole, environmentID string) (Role, bool) {
	for _, r := range roles {
		if r.EnvironmentId != environmentID {
			continue
		}
		switch r.Role {
		case accountproto.AccountV2_Role_Environment_VIEWER:
			return EnvironmentViewer, true
		case accountproto.AccountV2_Role_Environment_EDITOR:
			return EnvironmentEditor, true
		}
	}
	return "", false
}

// inheritsRole reports whether role inherits ancestor, directly or not.
func inheritsRole(role, ancestor Role) bool {
	for _, parent := range inherits[role] {
		if parent == ancestor || inheritsRole(parent, ancestor) {
			return true
		}
	}
	return false
}

// granted reports whether account holds permission in scope. The
// organization role counts in an environment only when it inherits an
// environment role (ADMIN, OWNER); a MEMBER counts in the organization alone,
// and needs an environment role for the environment, as in pkg/role.
func granted(account *accountproto.AccountV2, scope ScopeID, permission Permission) bool {
	if !inScope(permission, scope.Kind) {
		return false
	}
	if role, ok := organizationRole(account.OrganizationRole); ok && Grants(role, permission) {
		if scope.Kind == OrganizationScope || inheritsRole(role, EnvironmentViewer) {
			return true
		}
	}
	if scope.Kind != EnvironmentScope {
		return false
	}
	role, ok := environmentRole(account.EnvironmentRoles, scope.ID)
	return ok && Grants(role, permission)
}

// readableBySystemAdmin reports whether a system admin holds permission
// without being a member: everything the lowest role of either scope grants.
func readableBySystemAdmin(permission Permission) bool {
	return Grants(EnvironmentViewer, permission) || Grants(OrganizationMember, permission)
}
