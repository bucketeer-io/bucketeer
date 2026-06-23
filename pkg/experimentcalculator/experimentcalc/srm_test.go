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
//

package experimentcalc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// newRolloutFeature builds a Feature with a ROLLOUT default strategy carrying
// the given (variation_id, weight) pairs, in the order given.
func newRolloutFeature(t *testing.T, pairs ...struct {
	id     string
	weight int32
}) *featureproto.Feature {
	t.Helper()
	vars := make([]*featureproto.RolloutStrategy_Variation, 0, len(pairs))
	for _, p := range pairs {
		vars = append(vars, &featureproto.RolloutStrategy_Variation{
			Variation: p.id,
			Weight:    p.weight,
		})
	}
	return &featureproto.Feature{
		DefaultStrategy: &featureproto.Strategy{
			Type: featureproto.Strategy_ROLLOUT,
			RolloutStrategy: &featureproto.RolloutStrategy{
				Variations: vars,
			},
		},
	}
}

func vr(id string, evalUsers int64) *eventcounter.VariationResult {
	return &eventcounter.VariationResult{
		VariationId:     id,
		EvaluationCount: &eventcounter.VariationCount{UserCount: evalUsers},
	}
}

func TestComputeSRM_BalancedSplit_ReportsOK(t *testing.T) {
	t.Parallel()
	feature := newRolloutFeature(t,
		struct {
			id     string
			weight int32
		}{"vid1", 50}, struct {
			id     string
			weight int32
		}{"vid2", 50},
	)
	// 5000 + 4950 ~ 50/50; chi-square should be near 0, p_value near 1.
	results := []*eventcounter.VariationResult{
		vr("vid1", 5000), vr("vid2", 4950),
	}
	got := computeSRM(results, feature, DefaultSRMThreshold)

	assert.Equal(t, eventcounter.SrmResult_OK, got.Status)
	assert.InDelta(t, 0.001, got.Threshold, 1e-9)
	assert.EqualValues(t, 1, got.DegreesOfFreedom)
	// (5000-4975)^2/4975 + (4950-4975)^2/4975 = 625/4975 + 625/4975 ≈ 0.2513
	assert.InDelta(t, 0.2513, got.ChiSquare, 1e-3)
	// chi-square(1) CDF at 0.2513 ≈ 0.384, so p = 1-CDF ≈ 0.616 — well above
	// the 0.001 threshold, so this is clearly OK.
	assert.InDelta(t, 0.616, got.PValue, 0.01)
	// Per-variation breakdown is populated and normalized.
	if assert.Len(t, got.Variations, 2) {
		// Order is deterministic (sorted by variation_id).
		assert.Equal(t, "vid1", got.Variations[0].VariationId)
		assert.EqualValues(t, 5000, got.Variations[0].ObservedUserCount)
		assert.InDelta(t, 0.5, got.Variations[0].ExpectedWeight, 1e-9)
		assert.InDelta(t, 4975.0, got.Variations[0].ExpectedUserCount, 1e-9)
	}
}

func TestComputeSRM_SkewedSplit_ReportsMISMATCH(t *testing.T) {
	t.Parallel()
	feature := newRolloutFeature(t,
		struct {
			id     string
			weight int32
		}{"vid1", 50}, struct {
			id     string
			weight int32
		}{"vid2", 50},
	)
	// 50/50 intended but observed 5300/4700 — chi-square ≈ 36, p ≪ 0.001.
	results := []*eventcounter.VariationResult{
		vr("vid1", 5300), vr("vid2", 4700),
	}
	got := computeSRM(results, feature, DefaultSRMThreshold)

	assert.Equal(t, eventcounter.SrmResult_MISMATCH, got.Status)
	assert.Less(t, got.PValue, DefaultSRMThreshold,
		"36-σ-equivalent chi-square should produce p ≪ threshold")
	assert.Greater(t, got.ChiSquare, 30.0)
}

func TestComputeSRM_NonUniformWeights_ConvergesToExpected(t *testing.T) {
	t.Parallel()
	feature := newRolloutFeature(t,
		struct {
			id     string
			weight int32
		}{"vid1", 70}, struct {
			id     string
			weight int32
		}{"vid2", 30},
	)
	// Observed exactly 7000/3000 — chi-square = 0, p = 1.
	results := []*eventcounter.VariationResult{
		vr("vid1", 7000), vr("vid2", 3000),
	}
	got := computeSRM(results, feature, DefaultSRMThreshold)

	assert.Equal(t, eventcounter.SrmResult_OK, got.Status)
	assert.InDelta(t, 0.0, got.ChiSquare, 1e-9)
	assert.InDelta(t, 1.0, got.PValue, 1e-9)
	if assert.Len(t, got.Variations, 2) {
		assert.InDelta(t, 0.7, got.Variations[0].ExpectedWeight, 1e-9)
		assert.InDelta(t, 7000.0, got.Variations[0].ExpectedUserCount, 1e-9)
		assert.InDelta(t, 0.3, got.Variations[1].ExpectedWeight, 1e-9)
		assert.InDelta(t, 3000.0, got.Variations[1].ExpectedUserCount, 1e-9)
	}
}

func TestComputeSRM_ThreeVariations(t *testing.T) {
	t.Parallel()
	feature := newRolloutFeature(t,
		struct {
			id     string
			weight int32
		}{"vid1", 50}, struct {
			id     string
			weight int32
		}{"vid2", 25}, struct {
			id     string
			weight int32
		}{"vid3", 25},
	)
	// Roughly on-split (4980/2510/2510): df = K-1 = 2.
	results := []*eventcounter.VariationResult{
		vr("vid1", 4980), vr("vid2", 2510), vr("vid3", 2510),
	}
	got := computeSRM(results, feature, DefaultSRMThreshold)

	assert.Equal(t, eventcounter.SrmResult_OK, got.Status)
	assert.EqualValues(t, 2, got.DegreesOfFreedom)
	assert.Greater(t, got.PValue, 0.5)
}

func TestComputeSRM_DefaultThresholdFallback(t *testing.T) {
	t.Parallel()
	feature := newRolloutFeature(t,
		struct {
			id     string
			weight int32
		}{"vid1", 50}, struct {
			id     string
			weight int32
		}{"vid2", 50},
	)
	results := []*eventcounter.VariationResult{vr("vid1", 5000), vr("vid2", 5000)}
	// threshold <= 0 should fall back to DefaultSRMThreshold.
	got := computeSRM(results, feature, 0)
	assert.InDelta(t, DefaultSRMThreshold, got.Threshold, 1e-9)
}

func TestComputeSRM_VariationOrderIsStable(t *testing.T) {
	t.Parallel()
	feature := newRolloutFeature(t,
		struct {
			id     string
			weight int32
		}{"vid-c", 33}, struct {
			id     string
			weight int32
		}{"vid-a", 33}, struct {
			id     string
			weight int32
		}{"vid-b", 34},
	)
	results := []*eventcounter.VariationResult{
		vr("vid-a", 3300), vr("vid-b", 3400), vr("vid-c", 3300),
	}
	got := computeSRM(results, feature, DefaultSRMThreshold)

	if assert.Len(t, got.Variations, 3) {
		assert.Equal(t, "vid-a", got.Variations[0].VariationId)
		assert.Equal(t, "vid-b", got.Variations[1].VariationId)
		assert.Equal(t, "vid-c", got.Variations[2].VariationId)
	}
}

func TestComputeSRM_Skipped(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		feature  *featureproto.Feature
		results  []*eventcounter.VariationResult
		wantHint string
	}{
		{
			name:     "nil feature",
			feature:  nil,
			results:  []*eventcounter.VariationResult{vr("vid1", 5000), vr("vid2", 5000)},
			wantHint: "feature definition not available",
		},
		{
			name: "no default strategy",
			feature: &featureproto.Feature{
				DefaultStrategy: nil,
			},
			results:  []*eventcounter.VariationResult{vr("vid1", 5000), vr("vid2", 5000)},
			wantHint: "feature has no default strategy",
		},
		{
			name: "fixed strategy is not a rollout",
			feature: &featureproto.Feature{
				DefaultStrategy: &featureproto.Strategy{
					Type:          featureproto.Strategy_FIXED,
					FixedStrategy: &featureproto.FixedStrategy{Variation: "vid1"},
				},
			},
			results:  []*eventcounter.VariationResult{vr("vid1", 5000), vr("vid2", 5000)},
			wantHint: "default strategy is not a rollout",
		},
		{
			name: "all rollout weights zero",
			feature: newRolloutFeature(t, struct {
				id     string
				weight int32
			}{"vid1", 0}, struct {
				id     string
				weight int32
			}{"vid2", 0}),
			results:  []*eventcounter.VariationResult{vr("vid1", 5000), vr("vid2", 5000)},
			wantHint: "weights are zero",
		},
		{
			name: "total observed below minimum",
			feature: newRolloutFeature(t, struct {
				id     string
				weight int32
			}{"vid1", 50}, struct {
				id     string
				weight int32
			}{"vid2", 50}),
			results:  []*eventcounter.VariationResult{vr("vid1", 40), vr("vid2", 40)},
			wantHint: "below the minimum required",
		},
		{
			name: "single variation gives df < 1",
			feature: newRolloutFeature(t, struct {
				id     string
				weight int32
			}{"vid1", 100}),
			results:  []*eventcounter.VariationResult{vr("vid1", 5000)},
			wantHint: "fewer than 2 variations",
		},
		{
			name: "rollout strategy with zero variations",
			feature: &featureproto.Feature{
				DefaultStrategy: &featureproto.Strategy{
					Type:            featureproto.Strategy_ROLLOUT,
					RolloutStrategy: &featureproto.RolloutStrategy{},
				},
			},
			results:  []*eventcounter.VariationResult{vr("vid1", 5000)},
			wantHint: "no variations",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := computeSRM(tt.results, tt.feature, DefaultSRMThreshold)
			assert.Equal(t, eventcounter.SrmResult_SKIPPED, got.Status, "status")
			assert.True(t,
				strings.Contains(got.SkipReason, tt.wantHint),
				"expected skip_reason to contain %q, got %q", tt.wantHint, got.SkipReason)
		})
	}
}

func TestComputeSRM_HandlesMissingObservedCount(t *testing.T) {
	t.Parallel()
	feature := newRolloutFeature(t,
		struct {
			id     string
			weight int32
		}{"vid1", 50}, struct {
			id     string
			weight int32
		}{"vid2", 50},
	)
	// vid2 has no observed users (variation absent from results entirely).
	// The check should still run and detect the mismatch.
	results := []*eventcounter.VariationResult{vr("vid1", 5000)}
	got := computeSRM(results, feature, DefaultSRMThreshold)

	assert.Equal(t, eventcounter.SrmResult_MISMATCH, got.Status,
		"a variation receiving 0 of 5000 expected users should trigger MISMATCH")
	if assert.Len(t, got.Variations, 2) {
		var observedSum int64
		for _, v := range got.Variations {
			observedSum += v.ObservedUserCount
		}
		assert.EqualValues(t, 5000, observedSum)
	}
}
