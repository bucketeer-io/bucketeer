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
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	insightsproto "github.com/bucketeer-io/bucketeer/v2/proto/insights"
)

func TestFieldScopeResolve(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		scope       Scope
		req         proto.Message
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
			req:   &insightsproto.GetInsightsMonthlySummaryRequest{EnvironmentIds: []string{"env-1", "", "env-2"}},
			expected: []ScopeID{
				{Kind: EnvironmentScope, ID: "env-1"},
				{Kind: EnvironmentScope, ID: "env-2"},
			},
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
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := p.scope.Resolve(p.req)
			assert.ErrorIs(t, err, p.expectedErr)
			assert.Equal(t, p.expected, actual)
		})
	}
}
