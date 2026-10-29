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

package api

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_latencyQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc      string
		envIDs    []string
		sourceIDs []string
		apiIDs    []string
		expected  string
	}{
		{
			desc:     "no filters",
			expected: `environment_id:source_id:method:bucketeer_gateway_api_handling_seconds:avg:rate5m`,
		},
		{
			desc:     "with env filter",
			envIDs:   []string{"env1"},
			expected: `environment_id:source_id:method:bucketeer_gateway_api_handling_seconds:avg:rate5m{environment_id=~"^env1$"}`,
		},
		{
			desc:      "all filters",
			envIDs:    []string{"env1"},
			sourceIDs: []string{"GO_SERVER"},
			apiIDs:    []string{"GetEvaluations", "GetFeatureFlags"},
			expected:  `environment_id:source_id:method:bucketeer_gateway_api_handling_seconds:avg:rate5m{environment_id=~"^env1$",source_id=~"^GO_SERVER$",method=~"^(GetEvaluations|GetFeatureFlags)$"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			q := latencyQuery(tt.envIDs, tt.sourceIDs, tt.apiIDs)
			assert.Equal(t, tt.expected, q)
		})
	}
}

func Test_requestCountQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc      string
		envIDs    []string
		sourceIDs []string
		apiIDs    []string
		expected  string
	}{
		{
			desc:     "no filters",
			expected: `environment_id:source_id:method:bucketeer_gateway_api_request_total:rate5m`,
		},
		{
			desc:      "with filters",
			envIDs:    []string{"env1", "env2"},
			sourceIDs: []string{"GO_SERVER"},
			apiIDs:    []string{"GetEvaluations", "GetFeatureFlags"},
			expected:  `environment_id:source_id:method:bucketeer_gateway_api_request_total:rate5m{environment_id=~"^(env1|env2)$",source_id=~"^GO_SERVER$",method=~"^(GetEvaluations|GetFeatureFlags)$"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			q := requestCountQuery(tt.envIDs, tt.sourceIDs, tt.apiIDs)
			assert.Equal(t, tt.expected, q)
		})
	}
}

func Test_evaluationsQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc      string
		envIDs    []string
		sourceIDs []string
		expected  string
	}{
		{
			desc:     "no filters",
			expected: `sum by (environment_id, source_id, evaluation_type) (environment_id:pod:evaluation_type:source_id:bucketeer_api_gateway_evaluations_total:rate5m)`,
		},
		{
			desc:      "with filters",
			envIDs:    []string{"env1", "env2"},
			sourceIDs: []string{"ANDROID"},
			expected:  `sum by (environment_id, source_id, evaluation_type) (environment_id:pod:evaluation_type:source_id:bucketeer_api_gateway_evaluations_total:rate5m{environment_id=~"^(env1|env2)$",source_id=~"^ANDROID$"})`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			q := evaluationsQuery(tt.envIDs, tt.sourceIDs)
			assert.Equal(t, tt.expected, q)
		})
	}
}

func Test_errorRatesQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc      string
		envIDs    []string
		sourceIDs []string
		apiIDs    []string
		expected  string
	}{
		{
			desc:     "no filters",
			expected: `environment_id:source_id:method:bucketeer_gateway_api_error_rate:rate5m`,
		},
		{
			desc:      "with filters",
			envIDs:    []string{"env1"},
			sourceIDs: []string{"ANDROID"},
			apiIDs:    []string{"GetEvaluations"},
			expected:  `environment_id:source_id:method:bucketeer_gateway_api_error_rate:rate5m{environment_id=~"^env1$",source_id=~"^ANDROID$",method=~"^GetEvaluations$"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			q := errorRatesQuery(tt.envIDs, tt.sourceIDs, tt.apiIDs)
			assert.Equal(t, tt.expected, q)
		})
	}
}

func Test_buildSelector(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc      string
		envIDs    []string
		sourceIDs []string
		apiIDs    []string
		expected  string
	}{
		{
			desc:     "no filters",
			expected: "",
		},
		{
			desc:     "single environment id",
			envIDs:   []string{"env1"},
			expected: `{environment_id=~"^env1$"}`,
		},
		{
			desc:     "uuid environment id",
			envIDs:   []string{"0b7a7c6e-0d7f-4d3e-9c1a-2b3c4d5e6f70"},
			expected: `{environment_id=~"^0b7a7c6e-0d7f-4d3e-9c1a-2b3c4d5e6f70$"}`,
		},
		{
			desc:      "all filters",
			envIDs:    []string{"env1", "env2"},
			sourceIDs: []string{"GO_SERVER"},
			apiIDs:    []string{"GetEvaluations", "GetFeatureFlags"},
			expected:  `{environment_id=~"^(env1|env2)$",source_id=~"^GO_SERVER$",method=~"^(GetEvaluations|GetFeatureFlags)$"}`,
		},
		{
			desc:     "environment id containing a double quote and PromQL syntax stays inside the string literal",
			envIDs:   []string{`env1"} or up{job=~".+`},
			expected: `{environment_id=~"^env1\"\\} or up\\{job=~\"\\.\\+$"}`,
		},
		{
			desc:     "environment id containing a double quote within multiple values",
			envIDs:   []string{"env1", `env2"`},
			expected: `{environment_id=~"^(env1|env2\")$"}`,
		},
		{
			desc:     "environment id containing a backslash",
			envIDs:   []string{`env\1`},
			expected: `{environment_id=~"^env\\\\1$"}`,
		},
		{
			desc:     "environment id containing control characters",
			envIDs:   []string{"env\n1\t"},
			expected: `{environment_id=~"^env\n1\t$"}`,
		},
		{
			desc:     "environment id containing regex metacharacters",
			envIDs:   []string{"env.1+"},
			expected: `{environment_id=~"^env\\.1\\+$"}`,
		},
		{
			desc:      "source id and method containing double quotes",
			sourceIDs: []string{`GO"SERVER`},
			apiIDs:    []string{`Get"Evaluations`},
			expected:  `{source_id=~"^GO\"SERVER$",method=~"^Get\"Evaluations$"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			s := buildSelector(tt.envIDs, tt.sourceIDs, tt.apiIDs)
			assert.Equal(t, tt.expected, s)
		})
	}
}

// Test_buildSelectorValueRoundTrip checks that, whatever the environment ids
// contain, the selector consists of exactly one environment_id matcher whose
// value is a single well-formed double-quoted string literal, and that the
// regex inside it matches exactly the given ids once the literal is unquoted
// (which is what Prometheus does before compiling the matcher).
func Test_buildSelectorValueRoundTrip(t *testing.T) {
	t.Parallel()

	// A quoted string literal: an opening quote, any number of non-quote,
	// non-backslash characters or escape sequences, and a closing quote.
	selectorPattern := regexp.MustCompile(`^\{environment_id=~("(?:[^"\\]|\\.)*")\}$`)

	tests := []struct {
		desc   string
		envIDs []string
	}{
		{
			desc:   "plain ids",
			envIDs: []string{"env1", "0b7a7c6e-0d7f-4d3e-9c1a-2b3c4d5e6f70"},
		},
		{
			desc:   "double quote followed by PromQL syntax",
			envIDs: []string{`env1"} or up{job=~".+`},
		},
		{
			desc:   "double quote and backslash",
			envIDs: []string{`env1\"`, `"env2"`},
		},
		{
			desc:   "backslash before a regex metacharacter",
			envIDs: []string{`env\.1`},
		},
		{
			desc:   "control characters",
			envIDs: []string{"env\n1", "env\t2\r"},
		},
		{
			desc:   "non-ascii",
			envIDs: []string{"環境1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			s := buildSelector(tt.envIDs, nil, nil)

			match := selectorPattern.FindStringSubmatch(s)
			require.Len(t, match, 2, "selector is not a single quoted environment_id matcher: %s", s)

			unquoted, err := strconv.Unquote(match[1])
			require.NoError(t, err)
			re, err := regexp.Compile(unquoted)
			require.NoError(t, err)

			for _, id := range tt.envIDs {
				assert.True(t, re.MatchString(id), "regex %q must match %q", unquoted, id)
				assert.False(t, re.MatchString(id+"x"), "regex %q must not match %q", unquoted, id+"x")
				assert.False(t, re.MatchString("x"+id), "regex %q must not match %q", unquoted, "x"+id)
			}
			assert.False(t, re.MatchString("unrelated"), "regex %q must not match an unrelated id", unquoted)
		})
	}
}
