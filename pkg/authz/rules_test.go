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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	_ "github.com/bucketeer-io/bucketeer/v2/proto/account"
	_ "github.com/bucketeer-io/bucketeer/v2/proto/aichat"
	_ "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	_ "github.com/bucketeer-io/bucketeer/v2/proto/auth"
	_ "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	_ "github.com/bucketeer-io/bucketeer/v2/proto/coderef"
	_ "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	_ "github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
	_ "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	_ "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	_ "github.com/bucketeer-io/bucketeer/v2/proto/insights"
	_ "github.com/bucketeer-io/bucketeer/v2/proto/notification"
	_ "github.com/bucketeer-io/bucketeer/v2/proto/push"
	_ "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
	_ "github.com/bucketeer-io/bucketeer/v2/proto/tag"
	_ "github.com/bucketeer-io/bucketeer/v2/proto/team"
)

// TestWebRules fails when a method of a web service has no row in WebRules,
// when a row names a method or a request field that does not exist, or when
// a row is not a valid composition, so Table 2 cannot drift from the
// services it covers.
func TestWebRules(t *testing.T) {
	t.Parallel()
	methods := make(map[string]protoreflect.MethodDescriptor)
	for _, name := range WebServices {
		desc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(name))
		require.NoError(t, err, name)
		svc, ok := desc.(protoreflect.ServiceDescriptor)
		require.True(t, ok, name)
		for i := 0; i < svc.Methods().Len(); i++ {
			m := svc.Methods().Get(i)
			methods["/"+name+"/"+string(m.Name())] = m
		}
	}
	for method := range methods {
		assert.Contains(t, WebRules, method, "%s has no rule", method)
	}
	for method, rule := range WebRules {
		m, ok := methods[method]
		if !assert.True(t, ok, "%s is not a method of a web service", method) {
			continue
		}
		assert.NoError(t, rule.err, method)
		req, err := protoregistry.GlobalTypes.FindMessageByName(m.Input().FullName())
		require.NoError(t, err, method)
		for _, n := range rule.needs {
			_, err := n.scope.Resolve(context.Background(), req.New().Interface(), fakeLookups{})
			assert.False(t, errors.Is(err, ErrRule), "%s: %v", method, err)
		}
	}
}

func TestWebRulesSpotChecks(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		method  string
		isNeed  bool
		scopes  []ScopeKind
		kind    requirementKind
		guarded bool
	}{
		{method: "/bucketeer.feature.FeatureService/UpdateFeature", isNeed: true, scopes: []ScopeKind{EnvironmentScope}},
		{method: "/bucketeer.feature.FeatureService/BulkCloneFeature", isNeed: true,
			scopes: []ScopeKind{EnvironmentScope, EnvironmentScope}},
		{method: "/bucketeer.feature.FeatureService/CreateHoldoutGroup", isNeed: true,
			scopes: []ScopeKind{OrganizationOfEnvironmentScope}},
		{method: "/bucketeer.feature.FeatureService/FlagTriggerWebhook", kind: requirePublic},
		{method: "/bucketeer.account.AccountService/GetMe", isNeed: true, scopes: []ScopeKind{OrganizationScope}},
		{method: "/bucketeer.account.AccountService/GetMyOrganizations", kind: requireAuthenticated},
		{method: "/bucketeer.account.AccountService/GetEnvironmentAPIKey", kind: requireSystemAdmin},
		{method: "/bucketeer.environment.EnvironmentService/CreateOrganization", kind: requireSystemAdmin},
		{method: "/bucketeer.environment.EnvironmentService/ExchangeDemoToken", kind: requirePublic},
		{method: "/bucketeer.auth.AuthService/SignIn", kind: requirePublic},
	}
	for _, p := range patterns {
		t.Run(p.method[strings.LastIndex(p.method, "/")+1:], func(t *testing.T) {
			rule, ok := WebRules[p.method]
			require.True(t, ok)
			assert.Equal(t, p.kind, rule.kind)
			if !p.isNeed {
				assert.Empty(t, rule.needs)
				return
			}
			require.Len(t, rule.needs, len(p.scopes))
			for i, n := range rule.needs {
				assert.Equal(t, []ScopeKind{p.scopes[i]}, n.scope.Kinds())
			}
		})
	}
}
