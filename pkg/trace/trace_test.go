// Copyright 2024 The Bucketeer Authors.
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

package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"
)

func TestSampler(t *testing.T) {
	t.Parallel()
	filteringSpanName := "span-name"
	testcases := []struct {
		desc     string
		sampler  trace.Sampler
		name     string
		expected bool
	}{
		{
			desc: "false: filteringSpanName NeverSample",
			sampler: NewSampler(
				WithDefaultProbability(1.0),
				WithFilteringSampler(filteringSpanName, trace.NeverSample()),
			),
			name:     filteringSpanName,
			expected: false,
		},
		{
			desc: "false: filteringSpanName Probability=0.0",
			sampler: NewSampler(
				WithDefaultProbability(1.0),
				WithFilteringSampler(filteringSpanName, trace.ProbabilitySampler(0.0)),
			),
			name:     filteringSpanName,
			expected: false,
		},
		{
			desc: "true: filteringSpanName Probability=1.0",
			sampler: NewSampler(
				WithDefaultProbability(0.0),
				WithFilteringSampler(filteringSpanName, trace.ProbabilitySampler(1.0)),
			),
			name:     filteringSpanName,
			expected: true,
		},
		{
			desc: "false: default Probability=0.0",
			sampler: NewSampler(
				WithDefaultProbability(0.0),
				WithFilteringSampler(filteringSpanName, trace.ProbabilitySampler(1.0)),
			),
			name:     "default",
			expected: false,
		},
		{
			desc: "true: default Probability=1.0",
			sampler: NewSampler(
				WithDefaultProbability(1.0),
				WithFilteringSampler(filteringSpanName, trace.ProbabilitySampler(0.0)),
			),
			name:     "default",
			expected: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			decision := tc.sampler(trace.SamplingParameters{
				Name: tc.name,
			})
			assert.Equal(t, tc.expected, decision.Sample)
		})
	}
}
