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
	"math"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
)

func Test_requestCountQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc     string
		duration time.Duration
		expected string
	}{
		{
			desc:     "1 hour",
			duration: time.Hour,
			expected: `sum by (environment_id,source_id) (sum_over_time(environment_id:source_id:method:bucketeer_gateway_api_request_total:rate5m[3300s:5m])) * 300`,
		},
		{
			desc:     "31 days (January)",
			duration: 31 * 24 * time.Hour,
			expected: `sum by (environment_id,source_id) (sum_over_time(environment_id:source_id:method:bucketeer_gateway_api_request_total:rate5m[2678100s:5m])) * 300`,
		},
		{
			desc:     "30 days (April)",
			duration: 30 * 24 * time.Hour,
			expected: `sum by (environment_id,source_id) (sum_over_time(environment_id:source_id:method:bucketeer_gateway_api_request_total:rate5m[2591700s:5m])) * 300`,
		},
		{
			desc:     "29 days (February, leap year)",
			duration: 29 * 24 * time.Hour,
			expected: `sum by (environment_id,source_id) (sum_over_time(environment_id:source_id:method:bucketeer_gateway_api_request_total:rate5m[2505300s:5m])) * 300`,
		},
		{
			desc:     "28 days (February)",
			duration: 28 * 24 * time.Hour,
			expected: `sum by (environment_id,source_id) (sum_over_time(environment_id:source_id:method:bucketeer_gateway_api_request_total:rate5m[2418900s:5m])) * 300`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			q := requestCountQuery(tt.duration)
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
			desc:               "full April (30 days)",
			targetDate:         time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
			expectedDuration:   30 * 24 * time.Hour,
			expectedEvaluation: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
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
			desc: "float value rounded to int64",
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
				"env1": {"ANDROID": 124},
			},
		},
		{
			desc: "skip NaN value",
			vector: model.Vector{
				&model.Sample{
					Metric: model.Metric{
						"environment_id": "env1",
						"source_id":      "ANDROID",
					},
					Value: model.SampleValue(math.NaN()),
				},
				&model.Sample{
					Metric: model.Metric{
						"environment_id": "env1",
						"source_id":      "IOS",
					},
					Value: 100,
				},
			},
			expected: map[string]map[string]int64{
				"env1": {"IOS": 100},
			},
		},
		{
			desc: "skip Inf value",
			vector: model.Vector{
				&model.Sample{
					Metric: model.Metric{
						"environment_id": "env1",
						"source_id":      "ANDROID",
					},
					Value: model.SampleValue(math.Inf(1)),
				},
			},
			expected: map[string]map[string]int64{},
		},
		{
			desc: "skip negative value",
			vector: model.Vector{
				&model.Sample{
					Metric: model.Metric{
						"environment_id": "env1",
						"source_id":      "ANDROID",
					},
					Value: -10,
				},
			},
			expected: map[string]map[string]int64{},
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
