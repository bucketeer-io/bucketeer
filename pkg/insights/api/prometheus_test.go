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
	"testing"

	"github.com/stretchr/testify/assert"
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
