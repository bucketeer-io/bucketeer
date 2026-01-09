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

package evaluation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc            string
		targetFeatureID string
		variationIDs    []string
		flagVariations  map[string]string
		expected        bool
		expectedErr     error
	}{
		{
			desc:            "err: depended feature not found",
			targetFeatureID: "feature-1",
			variationIDs: []string{
				"variation-1",
				"variation-2",
			},
			flagVariations: map[string]string{},
			expected:       false,
			expectedErr:    ErrFeatureNotFound,
		},
		{
			desc:            "not matched",
			targetFeatureID: "feature-1",
			variationIDs: []string{
				"variation-1",
				"variation-2",
			},
			flagVariations: map[string]string{
				"feature-1": "variation-3",
			},
			expected:    false,
			expectedErr: nil,
		},
		{
			desc:            "success",
			targetFeatureID: "feature-1",
			variationIDs: []string{
				"variation-1",
				"variation-2",
			},
			flagVariations: map[string]string{
				"feature-1": "variation-2",
			},
			expected:    true,
			expectedErr: nil,
		},
	}
	eval := &dependencyEvaluator{}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := eval.Evaluate(p.targetFeatureID, p.variationIDs, p.flagVariations)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
