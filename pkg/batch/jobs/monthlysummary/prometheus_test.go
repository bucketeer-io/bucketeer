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

package monthlysummary

import (
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
)

func Test_requestCountIncreaseQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc     string
		duration time.Duration
		expected string
	}{
		{
			desc:     "1 hour",
			duration: time.Hour,
			expected: `sum by (environment_id,source_id) (increase(bucketeer_gateway_api_request_total[3600s]))`,
		},
		{
			desc:     "31 days (January)",
			duration: 31 * 24 * time.Hour,
			expected: `sum by (environment_id,source_id) (increase(bucketeer_gateway_api_request_total[2678400s]))`,
		},
		{
			desc:     "28 days (February)",
			duration: 28 * 24 * time.Hour,
			expected: `sum by (environment_id,source_id) (increase(bucketeer_gateway_api_request_total[2419200s]))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			q := requestCountIncreaseQuery(tt.duration)
			assert.Equal(t, tt.expected, q)
		})
	}
}

func Test_calculateTimeParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc               string
		targetDate         time.Time
		expectedDuration   time.Duration
		expectedEvaluation time.Time
	}{
		{
			desc:               "full January (31 days)",
			targetDate:         time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			expectedDuration:   31 * 24 * time.Hour,
			expectedEvaluation: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:               "full February leap year (29 days)",
			targetDate:         time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			expectedDuration:   29 * 24 * time.Hour,
			expectedEvaluation: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:               "mid-month (Jan 15)",
			targetDate:         time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			expectedDuration:   15 * 24 * time.Hour,
			expectedEvaluation: time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:               "first day of month",
			targetDate:         time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			expectedDuration:   24 * time.Hour,
			expectedEvaluation: time.Date(2024, 3, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:               "targetDate with non-zero time is normalized",
			targetDate:         time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC),
			expectedDuration:   15 * 24 * time.Hour,
			expectedEvaluation: time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			duration, evalTime := calculateTimeParams(tt.targetDate)
			assert.Equal(t, tt.expectedDuration, duration)
			assert.Equal(t, tt.expectedEvaluation, evalTime)
		})
	}
}

func Test_parseRequestCountVector(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc     string
		vector   model.Vector
		expected map[string]map[string]int64
	}{
		{
			desc:     "empty vector",
			vector:   model.Vector{},
			expected: map[string]map[string]int64{},
		},
		{
			desc: "single sample",
			vector: model.Vector{
				&model.Sample{
					Metric: model.Metric{
						"environment_id": "env1",
						"source_id":      "ANDROID",
					},
					Value: 100,
				},
			},
			expected: map[string]map[string]int64{
				"env1": {"ANDROID": 100},
			},
		},
		{
			desc: "multiple environments and sources",
			vector: model.Vector{
				&model.Sample{
					Metric: model.Metric{
						"environment_id": "env1",
						"source_id":      "ANDROID",
					},
					Value: 100,
				},
				&model.Sample{
					Metric: model.Metric{
						"environment_id": "env1",
						"source_id":      "IOS",
					},
					Value: 200,
				},
				&model.Sample{
					Metric: model.Metric{
						"environment_id": "env2",
						"source_id":      "ANDROID",
					},
					Value: 50,
				},
			},
			expected: map[string]map[string]int64{
				"env1": {"ANDROID": 100, "IOS": 200},
				"env2": {"ANDROID": 50},
			},
		},
		{
			desc: "skip samples with empty environment_id",
			vector: model.Vector{
				&model.Sample{
					Metric: model.Metric{
						"environment_id": "",
						"source_id":      "ANDROID",
					},
					Value: 100,
				},
				&model.Sample{
					Metric: model.Metric{
						"environment_id": "env1",
						"source_id":      "ANDROID",
					},
					Value: 200,
				},
			},
			expected: map[string]map[string]int64{
				"env1": {"ANDROID": 200},
			},
		},
		{
			desc: "skip samples with empty source_id",
			vector: model.Vector{
				&model.Sample{
					Metric: model.Metric{
						"environment_id": "env1",
						"source_id":      "",
					},
					Value: 100,
				},
				&model.Sample{
					Metric: model.Metric{
						"environment_id": "env1",
						"source_id":      "ANDROID",
					},
					Value: 200,
				},
			},
			expected: map[string]map[string]int64{
				"env1": {"ANDROID": 200},
			},
		},
		{
			desc: "skip samples with missing labels",
			vector: model.Vector{
				&model.Sample{
					Metric: model.Metric{
						"source_id": "ANDROID",
					},
					Value: 100,
				},
				&model.Sample{
					Metric: model.Metric{
						"environment_id": "env1",
					},
					Value: 150,
				},
				&model.Sample{
					Metric: model.Metric{
						"environment_id": "env1",
						"source_id":      "ANDROID",
					},
					Value: 200,
				},
			},
			expected: map[string]map[string]int64{
				"env1": {"ANDROID": 200},
			},
		},
		{
			desc: "float value truncated to int64",
			vector: model.Vector{
				&model.Sample{
					Metric: model.Metric{
						"environment_id": "env1",
						"source_id":      "ANDROID",
					},
					Value: 123.999,
				},
			},
			expected: map[string]map[string]int64{
				"env1": {"ANDROID": 123},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			result := parseRequestCountVector(tt.vector)
			assert.Equal(t, tt.expected, result)
		})
	}
}
