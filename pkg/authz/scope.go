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

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ScopeKind says what a scope ID identifies and which roles count there.
type ScopeKind int

const (
	EnvironmentScope ScopeKind = iota + 1
	OrganizationScope
	// OrganizationOfEnvironmentScope is the organization that owns the environment named by the ID.
	OrganizationOfEnvironmentScope
)

func (k ScopeKind) String() string {
	switch k {
	case EnvironmentScope:
		return "environment"
	case OrganizationScope:
		return "organization"
	case OrganizationOfEnvironmentScope:
		return "organization of environment"
	default:
		return fmt.Sprintf("ScopeKind(%d)", int(k))
	}
}

// ScopeID is the environment or organization a permission is checked in.
type ScopeID struct {
	Kind ScopeKind
	ID   string
}

func (s ScopeID) String() string {
	return s.Kind.String() + ":" + s.ID
}

// Lookups finds the scope of requests that carry only a resource ID.
type Lookups interface {
	// OrganizationOfProject returns the owner organization of the project, or a NotFound status.
	OrganizationOfProject(ctx context.Context, projectID string) (string, error)
}

// Scope resolves the scope IDs of a request. It returns ErrMissingScope when
// the request carries no ID, unless it is optional.
type Scope interface {
	// Kinds lists the kinds of scope Resolve may return.
	Kinds() []ScopeKind
	Resolve(ctx context.Context, req proto.Message, lookups Lookups) ([]ScopeID, error)
}

// FieldScope reads scope IDs from one string field of the request. The field
// is named as in the .proto file, and may be repeated.
type FieldScope struct {
	kind     ScopeKind
	field    string
	optional bool
}

// InEnv scopes the check to the environment named by field.
func InEnv(field string) FieldScope {
	return FieldScope{kind: EnvironmentScope, field: field}
}

// InOrg scopes the check to the organization named by field.
func InOrg(field string) FieldScope {
	return FieldScope{kind: OrganizationScope, field: field}
}

// InOrgOfEnv scopes the check to the organization that owns the environment named by field.
func InOrgOfEnv(field string) FieldScope {
	return FieldScope{kind: OrganizationOfEnvironmentScope, field: field}
}

// EachEnv checks every environment listed in the repeated field.
func EachEnv(field string) FieldScope {
	return InEnv(field)
}

// IfSet makes an empty singular field resolve to no scope instead of an
// error, for requests whose scope is one of several fields. The authorizer
// still rejects a request that resolves no scope at all.
func (s FieldScope) IfSet() FieldScope {
	s.optional = true
	return s
}

func (s FieldScope) Kinds() []ScopeKind {
	return []ScopeKind{s.kind}
}

// Resolve returns one scope per value of the field. An empty element of a
// repeated field is an error, because the handler would still use it; an
// empty singular field is an error unless the scope is optional.
func (s FieldScope) Resolve(_ context.Context, req proto.Message, _ Lookups) ([]ScopeID, error) {
	values, repeated, err := stringField(req, s.field)
	if err != nil {
		return nil, err
	}
	ids := make([]ScopeID, 0, len(values))
	for _, v := range values {
		if v == "" {
			if repeated {
				return nil, ErrMissingScope
			}
			continue
		}
		ids = append(ids, ScopeID{Kind: s.kind, ID: v})
	}
	if len(ids) == 0 && !s.optional {
		return nil, ErrMissingScope
	}
	return ids, nil
}

// ProjectScope scopes the check to the organization that owns the project named by field.
type ProjectScope struct {
	field string
}

// InOrgOfProject scopes the check to the organization that owns the project named by field.
func InOrgOfProject(field string) ProjectScope {
	return ProjectScope{field: field}
}

func (s ProjectScope) Kinds() []ScopeKind {
	return []ScopeKind{OrganizationScope}
}

func (s ProjectScope) Resolve(ctx context.Context, req proto.Message, lookups Lookups) ([]ScopeID, error) {
	ids, err := InOrg(s.field).Resolve(ctx, req, lookups)
	if err != nil {
		return nil, err
	}
	if lookups == nil {
		return nil, fmt.Errorf("%w: %s needs Lookups", ErrRule, s.field)
	}
	for i, id := range ids {
		organizationID, err := lookups.OrganizationOfProject(ctx, id.ID)
		if err != nil {
			return nil, err
		}
		ids[i] = ScopeID{Kind: OrganizationScope, ID: organizationID}
	}
	return ids, nil
}

// FirstSetScope resolves the first of its scopes that the request names.
type FirstSetScope struct {
	scopes []Scope
}

// FirstSet checks the first of scopes the request names, and only that one.
func FirstSet(scopes ...Scope) FirstSetScope {
	return FirstSetScope{scopes: scopes}
}

func (s FirstSetScope) Kinds() []ScopeKind {
	var kinds []ScopeKind
	for _, scope := range s.scopes {
		kinds = append(kinds, scope.Kinds()...)
	}
	return kinds
}

func (s FirstSetScope) Resolve(ctx context.Context, req proto.Message, lookups Lookups) ([]ScopeID, error) {
	for _, scope := range s.scopes {
		ids, err := scope.Resolve(ctx, req, lookups)
		if errors.Is(err, ErrMissingScope) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if len(ids) > 0 {
			return ids, nil
		}
	}
	return nil, ErrMissingScope
}

// stringField returns the values of a string field, one for a singular
// field and all elements for a repeated one. A field that does not exist or
// is not a string is a misconfigured rule.
func stringField(req proto.Message, name string) (values []string, repeated bool, err error) {
	m := req.ProtoReflect()
	fd := m.Descriptor().Fields().ByName(protoreflect.Name(name))
	if fd == nil || fd.Kind() != protoreflect.StringKind {
		return nil, false, fmt.Errorf("%w: %s has no string field %q", ErrRule, m.Descriptor().FullName(), name)
	}
	if !fd.IsList() {
		return []string{m.Get(fd).String()}, false, nil
	}
	list := m.Get(fd).List()
	values = make([]string, 0, list.Len())
	for i := 0; i < list.Len(); i++ {
		values = append(values, list.Get(i).String())
	}
	return values, true, nil
}
