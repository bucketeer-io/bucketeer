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
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ScopeKind says what a scope ID identifies.
type ScopeKind int

const (
	EnvironmentScope ScopeKind = iota + 1
	OrganizationScope
)

func (k ScopeKind) String() string {
	switch k {
	case EnvironmentScope:
		return "environment"
	case OrganizationScope:
		return "organization"
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

// Scope resolves the scope IDs of a request. It returns ErrMissingScope when
// the request carries no ID, unless it is optional.
type Scope interface {
	Resolve(req proto.Message) ([]ScopeID, error)
}

// FieldScope reads scope IDs from one string field of the request. The field
// is named as in the .proto file, and may be repeated.
type FieldScope struct {
	Kind     ScopeKind
	Field    string
	optional bool
}

// InEnv scopes the check to the environment named by field.
func InEnv(field string) FieldScope {
	return FieldScope{Kind: EnvironmentScope, Field: field}
}

// InOrg scopes the check to the organization named by field.
func InOrg(field string) FieldScope {
	return FieldScope{Kind: OrganizationScope, Field: field}
}

// EachEnv checks every environment listed in the repeated field.
func EachEnv(field string) FieldScope {
	return InEnv(field)
}

// IfSet makes an empty field resolve to no scope instead of an error, for
// requests whose scope is one of several fields. The authorizer still
// rejects a request that resolves no scope at all.
func (s FieldScope) IfSet() FieldScope {
	s.optional = true
	return s
}

func (s FieldScope) Resolve(req proto.Message) ([]ScopeID, error) {
	values, err := stringField(req, s.Field)
	if err != nil {
		return nil, err
	}
	ids := make([]ScopeID, 0, len(values))
	for _, v := range values {
		if v == "" {
			continue
		}
		ids = append(ids, ScopeID{Kind: s.Kind, ID: v})
	}
	if len(ids) == 0 && !s.optional {
		return nil, ErrMissingScope
	}
	return ids, nil
}

// stringField returns the values of a string field, one for a singular
// field and all elements for a repeated one. A field that does not exist or
// is not a string is a misconfigured rule.
func stringField(req proto.Message, name string) ([]string, error) {
	m := req.ProtoReflect()
	fd := m.Descriptor().Fields().ByName(protoreflect.Name(name))
	if fd == nil || fd.Kind() != protoreflect.StringKind {
		return nil, fmt.Errorf("%w: %s has no string field %q", ErrRule, m.Descriptor().FullName(), name)
	}
	if !fd.IsList() {
		return []string{m.Get(fd).String()}, nil
	}
	list := m.Get(fd).List()
	values := make([]string, 0, list.Len())
	for i := 0; i < list.Len(); i++ {
		values = append(values, list.Get(i).String())
	}
	return values, nil
}
