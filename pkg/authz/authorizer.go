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
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	accdomain "github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

var (
	ErrUnauthenticated  = status.Error(codes.Unauthenticated, "authz: unauthenticated")
	ErrPermissionDenied = status.Error(codes.PermissionDenied, "authz: permission denied")
	ErrMissingScope     = status.Error(codes.InvalidArgument, "authz: request has no environment or organization")
	ErrInternal         = status.Error(codes.Internal, "authz: internal")
	// ErrNoRule is returned for a method Table 2 has no row for.
	ErrNoRule = status.Error(codes.PermissionDenied, "authz: no rule for method")
	// ErrRule marks a rule that cannot be applied to its request, such as a
	// scope field the request does not have.
	ErrRule = errors.New("authz: invalid rule")
)

// Rules is Table 2: the requirement of every RPC, keyed by full method name
// ("/bucketeer.feature.FeatureService/UpdateFeature").
type Rules map[string]Requirement

// Account is the caller's account in one organization, as the authorizer
// loaded it.
type Account struct {
	*accountproto.AccountV2
}

// AccountLoader fetches the caller's accounts. Both methods return a
// NotFound status when the caller has no such account.
type AccountLoader interface {
	// Account is the caller's account in the organization that owns scope.
	Account(ctx context.Context, email string, scope ScopeID) (*accountproto.AccountV2, error)
	// SystemAdmin is the caller's account in a system admin organization.
	SystemAdmin(ctx context.Context, email string) (*accountproto.AccountV2, error)
}

// Authorizer applies Table 1 and Table 2 to a request.
type Authorizer struct {
	rules    Rules
	accounts AccountLoader
}

func New(rules Rules, accounts AccountLoader) *Authorizer {
	return &Authorizer{rules: rules, accounts: accounts}
}

// Rule returns the requirement of method.
func (a *Authorizer) Rule(method string) (Requirement, bool) {
	r, ok := a.rules[method]
	return r, ok
}

// Authorize decides whether caller may invoke method with req and returns
// the editor the handler acts as. The editor is nil for a public method.
// The error is a gRPC status: Unauthenticated for an anonymous caller or a
// non-member, PermissionDenied for an insufficient role, or the error a
// guard returned.
func (a *Authorizer) Authorize(
	ctx context.Context,
	caller Caller,
	method string,
	req proto.Message,
) (*eventproto.Editor, error) {
	rule, ok := a.rules[method]
	if !ok {
		return nil, ErrNoRule
	}
	if rule.err != nil {
		return nil, statusOf(rule.err)
	}
	if rule.kind == requirePublic {
		return nil, nil
	}
	if caller.Email == "" {
		return nil, ErrUnauthenticated
	}
	if caller.PublicAPI != nil && rule.kind != requirePermission {
		return nil, ErrPermissionDenied
	}
	fullAccess := caller.FullAccess()
	access := &Access{
		Method:      method,
		Request:     req,
		Caller:      caller,
		SystemAdmin: fullAccess,
		accounts:    &accountCache{loader: a.accounts},
	}
	var admin *Account
	if !fullAccess && caller.PublicAPI == nil {
		var err error
		if admin, err = access.accounts.systemAdmin(ctx, caller.Email); err != nil {
			return nil, err
		}
		access.SystemAdmin = admin != nil
	}
	switch rule.kind {
	case requireAuthenticated:
		return editorOf(access, admin, nil), nil
	case requireSystemAdmin:
		if !access.SystemAdmin {
			return nil, ErrPermissionDenied
		}
		return editorOf(access, admin, nil), nil
	}
	// A service identity skips the permission checks, but not the guards:
	// owner protection and the like apply whoever calls.
	if fullAccess && len(rule.guards) == 0 {
		return editorOf(access, nil, nil), nil
	}
	// The editor carries the account of the last scope checked: the target
	// of a clone, as the handlers choose today.
	var target *Account
	for _, n := range rule.needs {
		scopes, err := n.scope.Resolve(req)
		if err != nil {
			return nil, statusOf(err)
		}
		access.Scopes = append(access.Scopes, scopes...)
		if fullAccess {
			continue
		}
		for _, scope := range scopes {
			acc, err := a.check(ctx, access, n, scope)
			if err != nil {
				return nil, err
			}
			if acc != nil {
				target = acc
			}
		}
	}
	if len(access.Scopes) == 0 {
		return nil, ErrMissingScope
	}
	for _, g := range rule.guards {
		if err := g.Check(ctx, access); err != nil {
			return nil, statusOf(err)
		}
	}
	return editorOf(access, admin, target), nil
}

// check decides one need in one scope. A member's account is loaded first
// so a disabled or unknown account is rejected whatever the rule says; a
// system admin's is loaded only when the permission needs membership.
func (a *Authorizer) check(ctx context.Context, access *Access, n need, scope ScopeID) (*Account, error) {
	var account *Account
	if !access.SystemAdmin {
		var err error
		if account, err = access.Account(ctx, scope); err != nil {
			return nil, err
		}
	}
	if access.SystemAdmin && readableBySystemAdmin(n.permission) {
		return account, nil
	}
	for _, alternative := range n.alternatives {
		ok, err := holds(ctx, access, alternative)
		if err != nil {
			return nil, statusOf(err)
		}
		if ok {
			return account, nil
		}
	}
	if account == nil {
		var err error
		if account, err = access.Account(ctx, scope); err != nil {
			return nil, err
		}
	}
	if !granted(account.AccountV2, scope, n.permission) {
		return nil, ErrPermissionDenied
	}
	return account, nil
}

func holds(ctx context.Context, access *Access, conditions []Condition) (bool, error) {
	for _, c := range conditions {
		ok, err := c.Holds(ctx, access)
		if err != nil || !ok {
			return false, err
		}
	}
	return true, nil
}

// statusOf keeps gRPC statuses as they are and hides anything else, such as
// an ErrRule, behind Internal.
func statusOf(err error) error {
	if _, ok := status.FromError(err); ok {
		return err
	}
	return status.Error(codes.Internal, fmt.Sprintf("authz: %v", err))
}

// editorOf names the editor from the account of the scope checked, or from
// the system admin account.
func editorOf(access *Access, admin, target *Account) *eventproto.Editor {
	editor := &eventproto.Editor{
		Email:           access.Caller.Email,
		IsAdmin:         access.SystemAdmin,
		PublicApiEditor: access.Caller.PublicAPI,
	}
	named := target
	if named == nil {
		named = admin
	}
	if named != nil {
		editor.Name = (&accdomain.AccountV2{AccountV2: named.AccountV2}).GetAccountFullName()
	}
	if target != nil {
		editor.EnvironmentRoles = target.EnvironmentRoles
		editor.OrganizationRole = target.OrganizationRole
	}
	return editor
}

// accountCache loads each scope's account once per request.
type accountCache struct {
	loader AccountLoader
	loaded map[ScopeID]*Account
}

func (c *accountCache) get(ctx context.Context, email string, scope ScopeID) (*Account, error) {
	if account, ok := c.loaded[scope]; ok {
		return account, nil
	}
	if c.loader == nil {
		return nil, ErrInternal
	}
	account, err := active(c.loader.Account(ctx, email, scope))
	if err != nil {
		return nil, err
	}
	if c.loaded == nil {
		c.loaded = make(map[ScopeID]*Account)
	}
	c.loaded[scope] = account
	return account, nil
}

// systemAdmin returns the caller's system admin account, or nil when the
// caller has none or it is disabled.
func (c *accountCache) systemAdmin(ctx context.Context, email string) (*Account, error) {
	if c.loader == nil {
		return nil, ErrInternal
	}
	account, err := active(c.loader.SystemAdmin(ctx, email))
	if errors.Is(err, ErrUnauthenticated) {
		return nil, nil
	}
	return account, err
}

// active maps a loader result to an active account, ErrUnauthenticated for
// a missing or disabled one, or ErrInternal.
func active(account *accountproto.AccountV2, err error) (*Account, error) {
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrUnauthenticated
		}
		return nil, ErrInternal
	}
	if account == nil || account.Disabled {
		return nil, ErrUnauthenticated
	}
	return &Account{AccountV2: account}, nil
}

type editorContextKey struct{}

// WithEditor stores the editor the authorizer produced for the handler.
func WithEditor(ctx context.Context, editor *eventproto.Editor) context.Context {
	return context.WithValue(ctx, editorContextKey{}, editor)
}

// EditorFrom returns the editor stored by WithEditor, or nil.
func EditorFrom(ctx context.Context) *eventproto.Editor {
	editor, _ := ctx.Value(editorContextKey{}).(*eventproto.Editor)
	return editor
}
