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
	"fmt"

	"google.golang.org/protobuf/proto"
)

// Condition is one clause of an Or alternative: when every condition of an
// alternative holds, the caller passes without the permission.
type Condition interface {
	Holds(ctx context.Context, access *Access) (bool, error)
}

// Guard is a target-dependent rule checked after the permission, such as
// owner protection. It returns nil to allow.
type Guard interface {
	Check(ctx context.Context, access *Access) error
}

type requirementKind int

const (
	requirePermission requirementKind = iota
	requirePublic
	requireAuthenticated
	requireSystemAdmin
)

type need struct {
	permission   Permission
	scope        Scope
	alternatives [][]Condition
}

// Requirement is one row of Table 2: what a caller needs to invoke an RPC.
// Or, And and Guard apply to Need only; used on Public, Authenticated or
// SystemAdminOnly they make the requirement invalid, and the authorizer then
// rejects every call of the RPC instead of silently dropping a clause.
type Requirement struct {
	kind   requirementKind
	needs  []need
	guards []Guard
	err    error
}

func (r Requirement) composable(op string) Requirement {
	if r.err == nil && r.kind != requirePermission {
		r.err = fmt.Errorf("%w: %s is not allowed on %s", ErrRule, op, r.kind)
	}
	return r
}

func (k requirementKind) String() string {
	switch k {
	case requirePublic:
		return "Public"
	case requireAuthenticated:
		return "Authenticated"
	case requireSystemAdmin:
		return "SystemAdminOnly"
	default:
		return "Need"
	}
}

// Need requires permission in the scope resolved from the request.
func Need(permission Permission, scope Scope) Requirement {
	r := Requirement{needs: []need{{permission: permission, scope: scope}}}
	for _, kind := range scope.Kinds() {
		if !inScope(permission, kind) {
			r.err = fmt.Errorf("%w: %s is not %s-scoped", ErrRule, permission, kind)
			break
		}
	}
	return r
}

// Public needs nothing, not even a caller.
func Public() Requirement {
	return Requirement{kind: requirePublic}
}

// Authenticated needs a verified caller and nothing else, for RPCs that act
// on the caller's own data and check the account themselves.
func Authenticated() Requirement {
	return Requirement{kind: requireAuthenticated}
}

// SystemAdminOnly needs a system admin caller.
func SystemAdminOnly() Requirement {
	return Requirement{kind: requireSystemAdmin}
}

// Or adds an alternative to the last Need: the caller passes it when every
// condition holds, even without the permission.
func (r Requirement) Or(conditions ...Condition) Requirement {
	r = r.composable("Or")
	if r.err != nil || len(r.needs) == 0 || len(conditions) == 0 {
		return r
	}
	needs := append([]need(nil), r.needs...)
	last := &needs[len(needs)-1]
	last.alternatives = append(append([][]Condition(nil), last.alternatives...), conditions)
	r.needs = needs
	return r
}

// And requires other as well, for requests that carry several scopes with
// different permissions (a source and a target environment).
func (r Requirement) And(other Requirement) Requirement {
	r = r.composable("And")
	other = other.composable("And")
	if r.err == nil {
		r.err = other.err
	}
	r.needs = append(append([]need(nil), r.needs...), other.needs...)
	r.guards = append(append([]Guard(nil), r.guards...), other.guards...)
	return r
}

// Guard adds a guard checked after every Need passed.
func (r Requirement) Guard(guard Guard) Requirement {
	r = r.composable("Guard")
	r.guards = append(append([]Guard(nil), r.guards...), guard)
	return r
}

// Access is what conditions and guards see: the request, the caller, the
// scopes it resolved to and the caller's accounts, loaded on demand.
type Access struct {
	Method  string
	Request proto.Message
	Caller  Caller
	Scopes  []ScopeID
	// HasSystemAdminAccess holds for a system admin and for a full-access service.
	HasSystemAdminAccess bool

	accounts *accountCache
}

// Account returns the caller's account in the organization of scope, or
// ErrUnauthenticated when the caller is not an active member there.
func (a *Access) Account(ctx context.Context, scope ScopeID) (*Account, error) {
	return a.accounts.get(ctx, a.Caller.Email, scope)
}

type selfCondition struct {
	field string
}

// Self holds when the string field of the request is the caller's email.
func Self(field string) Condition {
	return selfCondition{field: field}
}

func (c selfCondition) Holds(_ context.Context, access *Access) (bool, error) {
	values, _, err := stringField(access.Request, c.field)
	if err != nil {
		return false, err
	}
	return len(values) == 1 && values[0] != "" && values[0] == access.Caller.Email, nil
}

type systemAdminCondition struct{}

// SystemAdmin holds for a system admin, for RPCs a system admin may call
// without being a member of the organization.
func SystemAdmin() Condition {
	return systemAdminCondition{}
}

func (systemAdminCondition) Holds(_ context.Context, access *Access) (bool, error) {
	return access.HasSystemAdminAccess, nil
}
