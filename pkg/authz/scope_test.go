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
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	insightsproto "github.com/bucketeer-io/bucketeer/v2/proto/insights"
	pushproto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

var errProjectNotFound = status.Error(codes.NotFound, "project not found")

// fakeLookups maps project IDs to organization IDs.
type fakeLookups map[string]string

func (f fakeLookups) OrganizationOfProject(_ context.Context, projectID string) (string, error) {
	org, ok := f[projectID]
	if !ok {
		return "", errProjectNotFound
	}
	return org, nil
}

func TestScopeResolve(t *testing.T) {
	t.Parallel()
	lookups := fakeLookups{"proj-1": "org-1"}
	patterns := []struct {
		desc        string
		scope       Scope
		req         proto.Message
		lookups     Lookups
		expected    []ScopeID
		expectedErr error
	}{
		{
			desc:     "environment field",
			scope:    InEnv("environment_id"),
			req:      &featureproto.UpdateFeatureRequest{EnvironmentId: "env-1"},
			expected: []ScopeID{{Kind: EnvironmentScope, ID: "env-1"}},
		},
		{
			desc:     "organization field",
			scope:    InOrg("organization_id"),
			req:      &accountproto.ListAccountsV2Request{OrganizationId: "org-1"},
			expected: []ScopeID{{Kind: OrganizationScope, ID: "org-1"}},
		},
		{
			desc:  "repeated field",
			scope: EachEnv("environment_ids"),
			req:   &insightsproto.GetInsightsMonthlySummaryRequest{EnvironmentIds: []string{"env-1", "env-2"}},
			expected: []ScopeID{
				{Kind: EnvironmentScope, ID: "env-1"},
				{Kind: EnvironmentScope, ID: "env-2"},
			},
		},
		{
			desc:        "empty element of a repeated field",
			scope:       EachEnv("environment_ids").IfSet(),
			req:         &insightsproto.GetInsightsMonthlySummaryRequest{EnvironmentIds: []string{"env-1", "", "env-2"}},
			expectedErr: ErrMissingScope,
		},
		{
			desc:        "empty field",
			scope:       InEnv("environment_id"),
			req:         &featureproto.UpdateFeatureRequest{},
			expectedErr: ErrMissingScope,
		},
		{
			desc:        "empty repeated field",
			scope:       EachEnv("environment_ids"),
			req:         &insightsproto.GetInsightsMonthlySummaryRequest{},
			expectedErr: ErrMissingScope,
		},
		{
			desc:     "optional empty field",
			scope:    InEnv("environment_id").IfSet(),
			req:      &featureproto.UpdateFeatureRequest{},
			expected: []ScopeID{},
		},
		{
			desc:        "unknown field",
			scope:       InEnv("namespace"),
			req:         &featureproto.UpdateFeatureRequest{},
			expectedErr: ErrRule,
		},
		{
			desc:        "non-string field",
			scope:       InEnv("page_size"),
			req:         &accountproto.ListAccountsV2Request{},
			expectedErr: ErrRule,
		},
		{
			desc:     "organization of environment",
			scope:    InOrgOfEnv("id"),
			req:      &environmentproto.GetEnvironmentV2Request{Id: "env-1"},
			expected: []ScopeID{{Kind: OrganizationOfEnvironmentScope, ID: "env-1"}},
		},
		{
			desc:     "organization of project",
			scope:    InOrgOfProject("id"),
			req:      &environmentproto.GetProjectRequest{Id: "proj-1"},
			lookups:  lookups,
			expected: []ScopeID{{Kind: OrganizationScope, ID: "org-1"}},
		},
		{
			desc:        "organization of an unknown project",
			scope:       InOrgOfProject("id"),
			req:         &environmentproto.GetProjectRequest{Id: "proj-2"},
			lookups:     lookups,
			expectedErr: errProjectNotFound,
		},
		{
			desc:        "organization of project without lookups",
			scope:       InOrgOfProject("id"),
			req:         &environmentproto.GetProjectRequest{Id: "proj-1"},
			expectedErr: ErrRule,
		},
		{
			desc:        "organization of project with an empty field",
			scope:       InOrgOfProject("id"),
			req:         &environmentproto.GetProjectRequest{},
			lookups:     lookups,
			expectedErr: ErrMissingScope,
		},
		{
			desc:     "first set: organization",
			scope:    FirstSet(InOrg("organization_id"), InEnv("environment_id")),
			req:      &pushproto.ListPushesRequest{OrganizationId: "org-1", EnvironmentId: "env-1"},
			expected: []ScopeID{{Kind: OrganizationScope, ID: "org-1"}},
		},
		{
			desc:     "first set: environment",
			scope:    FirstSet(InOrg("organization_id"), InEnv("environment_id")),
			req:      &pushproto.ListPushesRequest{EnvironmentId: "env-1"},
			expected: []ScopeID{{Kind: EnvironmentScope, ID: "env-1"}},
		},
		{
			desc:        "first set: neither",
			scope:       FirstSet(InOrg("organization_id"), InEnv("environment_id")),
			req:         &pushproto.ListPushesRequest{},
			expectedErr: ErrMissingScope,
		},
		{
			desc:        "first set: invalid rule",
			scope:       FirstSet(InOrg("namespace"), InEnv("environment_id")),
			req:         &pushproto.ListPushesRequest{EnvironmentId: "env-1"},
			expectedErr: ErrRule,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := p.scope.Resolve(context.Background(), p.req, p.lookups)
			assert.ErrorIs(t, err, p.expectedErr)
			assert.Equal(t, p.expected, actual)
		})
	}
}
