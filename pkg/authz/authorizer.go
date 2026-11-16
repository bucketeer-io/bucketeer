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

// AccountLoader fetches the caller's account in the organization that owns
// a scope. It returns a NotFound status when the caller is not a member.
type AccountLoader interface {
	Account(ctx context.Context, email string, scope ScopeID) (*accountproto.AccountV2, error)
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

// Authorize decides whether the caller in ctx may invoke method with req and
// returns the editor the handler acts as. The editor is nil for a public
// method. The error is a gRPC status: Unauthenticated for a missing token or
// a non-member, PermissionDenied for an insufficient role, or the error a
// guard returned.
func (a *Authorizer) Authorize(ctx context.Context, method string, req proto.Message) (*eventproto.Editor, error) {
	rule, ok := a.rules[method]
	if !ok {
		return nil, ErrNoRule
	}
	if rule.kind == requirePublic {
		return nil, nil
	}
	caller, ok := callerFromContext(ctx)
	if !ok {
		return nil, ErrUnauthenticated
	}
	access := &Access{
		Method:   method,
		Request:  req,
		Caller:   caller,
		accounts: &accountCache{loader: a.accounts},
	}
	switch rule.kind {
	case requireAuthenticated:
		return tokenEditor(caller), nil
	case requireSystemAdmin:
		if !caller.IsSystemAdmin {
			return nil, ErrPermissionDenied
		}
		return tokenEditor(caller), nil
	}
	if caller.FullAccess() {
		return tokenEditor(caller), nil
	}
	var account *Account
	for _, n := range rule.needs {
		scopes, err := n.scope.Resolve(req)
		if err != nil {
			return nil, statusOf(err)
		}
		access.Scopes = append(access.Scopes, scopes...)
		for _, scope := range scopes {
			acc, err := a.check(ctx, access, n, scope)
			if err != nil {
				return nil, err
			}
			if account == nil {
				account = acc
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
	if account == nil {
		return tokenEditor(caller), nil
	}
	return accountEditor(caller, account), nil
}

// check decides one need in one scope. A member's account is loaded first
// so a disabled or unknown account is rejected whatever the rule says; a
// system admin's is loaded only when the permission needs membership.
func (a *Authorizer) check(ctx context.Context, access *Access, n need, scope ScopeID) (*Account, error) {
	caller := access.Caller
	var account *Account
	if !caller.IsSystemAdmin {
		var err error
		if account, err = access.Account(ctx, scope); err != nil {
			return nil, err
		}
	}
	if caller.IsSystemAdmin && readableBySystemAdmin(n.permission) {
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

func tokenEditor(caller *Caller) *eventproto.Editor {
	return &eventproto.Editor{
		Email:           caller.Email,
		Name:            caller.Name,
		IsAdmin:         caller.IsSystemAdmin,
		PublicApiEditor: caller.PublicAPI,
	}
}

func accountEditor(caller *Caller, account *Account) *eventproto.Editor {
	name := caller.Name
	if caller.PublicAPI != nil {
		name = (&accdomain.AccountV2{AccountV2: account.AccountV2}).GetAccountFullName()
	}
	return &eventproto.Editor{
		Email:            caller.Email,
		Name:             name,
		IsAdmin:          caller.IsSystemAdmin,
		PublicApiEditor:  caller.PublicAPI,
		EnvironmentRoles: account.EnvironmentRoles,
		OrganizationRole: account.OrganizationRole,
	}
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
	acc, err := c.loader.Account(ctx, email, scope)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrUnauthenticated
		}
		return nil, ErrInternal
	}
	if acc == nil || acc.Disabled {
		return nil, ErrUnauthenticated
	}
	account := &Account{AccountV2: acc}
	if c.loaded == nil {
		c.loaded = make(map[ScopeID]*Account)
	}
	c.loaded[scope] = account
	return account, nil
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
