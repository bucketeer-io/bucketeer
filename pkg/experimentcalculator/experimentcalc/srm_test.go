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

// newRolloutFeatureWithAudience builds a Feature with a ROLLOUT default
// strategy and a populated Audience block (audience.percentage =
// audiencePercentage, audience.default_variation = defaultVariation). Pass
// audiencePercentage = 100 (or 0) and defaultVariation = "" to get the
// no-filtering case — strategy_evaluator.go only filters when 0 < pct < 100.
func newRolloutFeatureWithAudience(
	t *testing.T,
	audiencePercentage int32,
	defaultVariation string,
	pairs ...struct {
		id     string
		weight int32
	},
) *featureproto.Feature {
	t.Helper()
	f := newRolloutFeature(t, pairs...)
	f.DefaultStrategy.RolloutStrategy.Audience = &featureproto.Audience{
		Percentage:       audiencePercentage,
		DefaultVariation: defaultVariation,
	}
	return f
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
			// 99/1 split with totalObserved=150 easily clears the
			// minSRMSampleSize floor but produces an expected count of
			// just 1.5 in the small cell — far below the chi-square
			// "all expected >= 5" reliability floor. The check must
			// SKIP rather than emit a p-value the user shouldn't trust.
			name: "skewed split with small expected cell",
			feature: newRolloutFeature(t, struct {
				id     string
				weight int32
			}{"vid1", 99}, struct {
				id     string
				weight int32
			}{"vid2", 1}),
			results:  []*eventcounter.VariationResult{vr("vid1", 148), vr("vid2", 2)},
			wantHint: "smallest expected = 1.50",
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
		{
			// audience.default_variation is set to a variation that is
			// NOT in the rollout's variation list. UI validation should
			// prevent this, but if a misconfigured flag reaches the
			// calculator we refuse to compute SRM (we'd otherwise
			// silently mis-attribute the out-of-audience traffic to a
			// variation that doesn't exist).
			name: "audience default_variation not in rollout",
			feature: newRolloutFeatureWithAudience(t, 50, "vid-not-in-rollout",
				struct {
					id     string
					weight int32
				}{"vid1", 50}, struct {
					id     string
					weight int32
				}{"vid2", 50}),
			results:  []*eventcounter.VariationResult{vr("vid1", 5000), vr("vid2", 5000)},
			wantHint: "audience default_variation is not one of the rollout variations",
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

func TestComputeSRM_IncludesUnexpectedVariationFromObserved(t *testing.T) {
	t.Parallel()
	// Rollout intends a 50/50 split between vid1 and vid2, but observed
	// traffic also contains 2000 users on `vid-leak` — a variation that
	// isn't in the rollout (experiment schema drift, leaked traffic from a
	// stale bucketing decision, etc.). The leaked users must count toward
	// totalObserved so the expected counts for the known variations reflect
	// the real denominator (and the mismatch isn't silently hidden), and
	// the per-variation breakdown must surface the unknown variation with
	// expected_weight = 0 so the UI can show it.
	feature := newRolloutFeature(t,
		struct {
			id     string
			weight int32
		}{"vid1", 50}, struct {
			id     string
			weight int32
		}{"vid2", 50},
	)
	results := []*eventcounter.VariationResult{
		vr("vid1", 4000), vr("vid2", 4000), vr("vid-leak", 2000),
	}
	got := computeSRM(results, feature, DefaultSRMThreshold)

	// totalObserved = 10000 (includes the leaked 2000), so each known
	// variation's expected count is 5000, well above the 4000 observed →
	// large chi-square, MISMATCH.
	assert.Equal(t, eventcounter.SrmResult_MISMATCH, got.Status,
		"leaked traffic to an unknown variation must not be silently dropped from totalObserved")
	if assert.Len(t, got.Variations, 3) {
		// Deterministic order: sorted by variation_id → vid-leak, vid1, vid2.
		assert.Equal(t, "vid-leak", got.Variations[0].VariationId)
		assert.EqualValues(t, 2000, got.Variations[0].ObservedUserCount)
		assert.InDelta(t, 0.0, got.Variations[0].ExpectedWeight, 1e-9,
			"unknown variation must surface with expected_weight=0")
		assert.InDelta(t, 0.0, got.Variations[0].ExpectedUserCount, 1e-9)

		assert.Equal(t, "vid1", got.Variations[1].VariationId)
		assert.EqualValues(t, 4000, got.Variations[1].ObservedUserCount)
		assert.InDelta(t, 5000.0, got.Variations[1].ExpectedUserCount, 1e-9,
			"expected count for known variations must use totalObserved=10000")
	}
	// dof = K_pos - 1 = 2 - 1 = 1 (only the two cells with positive
	// expected counts contribute; the leaked variation is excluded from the
	// chi-square sum because dividing by expected=0 is undefined).
	assert.EqualValues(t, 1, got.DegreesOfFreedom)
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

// TestExtractExpectedFractions covers the pure math of audience-aware expected
// fractions. The function returns a normalized fraction in [0, 1] per
// variation; across the rollout's variations the fractions must sum to
// exactly 1.0 (modulo float rounding), and the per-variation values must
// match the analytic formula:
//
//	expected_fraction(V_i) = a · p_i              if V_i != D
//	expected_fraction(D)   = a · p_D + (1 - a)
//
// where p_i = w_i / Σw and a = audience.percentage / 100 (or 1.0 when the
// audience filter doesn't apply, per strategy_evaluator.go semantics).
func TestExtractExpectedFractions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		feature *featureproto.Feature
		want    map[string]float64
	}{
		{
			name: "no audience block → raw weights normalized",
			feature: newRolloutFeature(t,
				struct {
					id     string
					weight int32
				}{"A", 50}, struct {
					id     string
					weight int32
				}{"B", 50}),
			want: map[string]float64{"A": 0.5, "B": 0.5},
		},
		{
			name: "audience.percentage = 100 → audience filter disabled, raw weights",
			feature: newRolloutFeatureWithAudience(t, 100, "A",
				struct {
					id     string
					weight int32
				}{"A", 70}, struct {
					id     string
					weight int32
				}{"B", 30}),
			want: map[string]float64{"A": 0.7, "B": 0.3},
		},
		{
			name: "audience.percentage = 0 → audience filter disabled, raw weights",
			feature: newRolloutFeatureWithAudience(t, 0, "A",
				struct {
					id     string
					weight int32
				}{"A", 50}, struct {
					id     string
					weight int32
				}{"B", 50}),
			want: map[string]float64{"A": 0.5, "B": 0.5},
		},
		{
			name: "audience.default_variation = \"\" → excluded users get no event, raw weights",
			feature: newRolloutFeatureWithAudience(t, 50, "",
				struct {
					id     string
					weight int32
				}{"A", 50}, struct {
					id     string
					weight int32
				}{"B", 50}),
			want: map[string]float64{"A": 0.5, "B": 0.5},
		},
		{
			// THE REGRESSION CASE for the user's concrete example:
			// audience = 50%, default = Control, 50/50 A/Control.
			// Out-of-audience users (50% of total) all go to Control,
			// so Control's expected share is 0.5*0.5 + 0.5 = 0.75 and
			// A's is 0.5*0.5 = 0.25.
			name: "audience = 50%, default = Control: Control gets the excluded 50%",
			feature: newRolloutFeatureWithAudience(t, 50, "Control",
				struct {
					id     string
					weight int32
				}{"A", 50}, struct {
					id     string
					weight int32
				}{"Control", 50}),
			want: map[string]float64{"A": 0.25, "Control": 0.75},
		},
		{
			// Same shape but default = the treatment (B). Confirms the
			// excluded-credit goes to whichever variation is named the
			// default — not always Control.
			name: "audience = 50%, default = treatment: treatment gets the excluded 50%",
			feature: newRolloutFeatureWithAudience(t, 50, "B",
				struct {
					id     string
					weight int32
				}{"A", 50}, struct {
					id     string
					weight int32
				}{"B", 50}),
			want: map[string]float64{"A": 0.25, "B": 0.75},
		},
		{
			// Three variations with uneven weights; default is the
			// largest. Verifies the formula generalizes beyond 2-way
			// splits and that the excluded credit goes to the named
			// default regardless of its weight share.
			name: "audience = 20%, three variations, default = largest",
			feature: newRolloutFeatureWithAudience(t, 20, "A",
				struct {
					id     string
					weight int32
				}{"A", 60}, struct {
					id     string
					weight int32
				}{"B", 30}, struct {
					id     string
					weight int32
				}{"C", 10}),
			// a=0.2; p_A=0.6, p_B=0.3, p_C=0.1
			// f_A = 0.2*0.6 + 0.8 = 0.92
			// f_B = 0.2*0.3       = 0.06
			// f_C = 0.2*0.1       = 0.02
			want: map[string]float64{"A": 0.92, "B": 0.06, "C": 0.02},
		},
		{
			// Edge: audience = 1% (smallest non-zero filter), default = A.
			// Nearly all users are out-of-audience → default dominates.
			name: "audience = 1%, default = A: A's share approaches 1.0",
			feature: newRolloutFeatureWithAudience(t, 1, "A",
				struct {
					id     string
					weight int32
				}{"A", 50}, struct {
					id     string
					weight int32
				}{"B", 50}),
			// a=0.01; p_A=p_B=0.5
			// f_A = 0.01*0.5 + 0.99 = 0.995
			// f_B = 0.01*0.5        = 0.005
			want: map[string]float64{"A": 0.995, "B": 0.005},
		},
		{
			// Edge: audience = 99% (largest non-100% filter), default = A.
			// Almost no excluded credit; fractions ≈ raw.
			name: "audience = 99%, default = A: fractions stay near raw weights",
			feature: newRolloutFeatureWithAudience(t, 99, "A",
				struct {
					id     string
					weight int32
				}{"A", 50}, struct {
					id     string
					weight int32
				}{"B", 50}),
			// a=0.99; p_A=p_B=0.5
			// f_A = 0.99*0.5 + 0.01 = 0.505
			// f_B = 0.99*0.5        = 0.495
			want: map[string]float64{"A": 0.505, "B": 0.495},
		},
		{
			// Uneven in-audience split (70/30) AND audience filter (50%).
			// Verifies the two layers compose correctly without bias.
			name: "audience = 50%, weights 70/30, default = minority",
			feature: newRolloutFeatureWithAudience(t, 50, "B",
				struct {
					id     string
					weight int32
				}{"A", 70}, struct {
					id     string
					weight int32
				}{"B", 30}),
			// a=0.5; p_A=0.7, p_B=0.3
			// f_A = 0.5*0.7       = 0.35
			// f_B = 0.5*0.3 + 0.5 = 0.65
			want: map[string]float64{"A": 0.35, "B": 0.65},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := extractExpectedFractions(tt.feature)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, len(tt.want), len(got), "variation count")
			var sum float64
			for vid, wantF := range tt.want {
				gotF, ok := got[vid]
				if assert.True(t, ok, "variation %q missing from output", vid) {
					assert.InDelta(t, wantF, gotF, 1e-9,
						"expected_fraction for %q", vid)
				}
				sum += gotF
			}
			// Invariant: fractions over the rollout's variations sum to 1.
			assert.InDelta(t, 1.0, sum, 1e-9,
				"expected_fractions must sum to exactly 1 across rollout variations")
		})
	}
}

// TestComputeSRM_AudienceAware verifies the end-to-end SRM behavior with
// Audience Traffic Allocation < 100%. The crown jewel of this suite is the
// regression case from the user-reported bug: a perfectly-configured 50%
// audience with default=Control and a 50/50 variation split produces 25/75
// observed (because out-of-audience users all land on Control), and SRM
// must report OK, NOT MISMATCH.
func TestComputeSRM_AudienceAware(t *testing.T) {
	t.Parallel()

	t.Run("regression: 25/75 observed under 50% audience is OK, not MISMATCH", func(t *testing.T) {
		t.Parallel()
		// Exactly the scenario in the bug report:
		//   Audience = 50%, default = Control, weights A/Control = 50/50
		//   10,000 SDK calls
		//   → 5,000 in audience: 2,500 A + 2,500 Control
		//   → 5,000 out:           0 A + 5,000 Control
		//   Observed: A=2,500, Control=7,500 (exactly the audience-adjusted
		//   expected split of 0.25 / 0.75 against totalObserved=10,000).
		feature := newRolloutFeatureWithAudience(t, 50, "Control",
			struct {
				id     string
				weight int32
			}{"A", 50}, struct {
				id     string
				weight int32
			}{"Control", 50})
		results := []*eventcounter.VariationResult{
			vr("A", 2500), vr("Control", 7500),
		}
		got := computeSRM(results, feature, DefaultSRMThreshold)

		assert.Equal(t, eventcounter.SrmResult_OK, got.Status,
			"audience-aware SRM must NOT flag MISMATCH on perfectly-configured 50%% audience")
		assert.InDelta(t, 0.0, got.ChiSquare, 1e-6,
			"observed == expected → chi-square should be ~0")
		assert.Greater(t, got.PValue, 0.5,
			"p-value should be high when observed matches expected exactly")
		// Per-variation breakdown must surface the audience-adjusted shares.
		if assert.Len(t, got.Variations, 2) {
			assert.Equal(t, "A", got.Variations[0].VariationId)
			assert.InDelta(t, 0.25, got.Variations[0].ExpectedWeight, 1e-9)
			assert.InDelta(t, 2500.0, got.Variations[0].ExpectedUserCount, 1e-9)
			assert.Equal(t, "Control", got.Variations[1].VariationId)
			assert.InDelta(t, 0.75, got.Variations[1].ExpectedWeight, 1e-9)
			assert.InDelta(t, 7500.0, got.Variations[1].ExpectedUserCount, 1e-9)
		}
	})

	t.Run("genuine SRM beyond what audience accounts for is still MISMATCH", func(t *testing.T) {
		t.Parallel()
		// Same audience config as above, but observed is skewed in a way
		// the audience can't explain: A got 1,500 (less than the
		// 2,500 expected), Control got 8,500. This is a real bucketing bug
		// and must be flagged.
		feature := newRolloutFeatureWithAudience(t, 50, "Control",
			struct {
				id     string
				weight int32
			}{"A", 50}, struct {
				id     string
				weight int32
			}{"Control", 50})
		results := []*eventcounter.VariationResult{
			vr("A", 1500), vr("Control", 8500),
		}
		got := computeSRM(results, feature, DefaultSRMThreshold)
		assert.Equal(t, eventcounter.SrmResult_MISMATCH, got.Status,
			"a genuine skew beyond what audience explains must still be detected")
		assert.Less(t, got.PValue, DefaultSRMThreshold)
	})

	t.Run("audience = 100% with skewed observed: MISMATCH (no behavior change vs no-audience case)", func(t *testing.T) {
		t.Parallel()
		// Sanity: audience=100% must be a strict no-op vs the
		// no-audience-block case. A 50/50 rollout observed as 53/47 with
		// n=10k should still detect MISMATCH.
		feature := newRolloutFeatureWithAudience(t, 100, "Control",
			struct {
				id     string
				weight int32
			}{"A", 50}, struct {
				id     string
				weight int32
			}{"Control", 50})
		results := []*eventcounter.VariationResult{
			vr("A", 5300), vr("Control", 4700),
		}
		got := computeSRM(results, feature, DefaultSRMThreshold)
		assert.Equal(t, eventcounter.SrmResult_MISMATCH, got.Status)
		assert.Less(t, got.PValue, DefaultSRMThreshold)
	})

	t.Run("audience = 50%, default_variation empty: raw-weight semantics", func(t *testing.T) {
		t.Parallel()
		// When default_variation is empty, the SDK fires no event for
		// excluded users (per strategy_evaluator.go), so observed only
		// contains in-audience users and the raw 50/50 weights are the
		// correct expected split. n must clear the 100-user floor.
		feature := newRolloutFeatureWithAudience(t, 50, "",
			struct {
				id     string
				weight int32
			}{"A", 50}, struct {
				id     string
				weight int32
			}{"B", 50})
		results := []*eventcounter.VariationResult{
			vr("A", 2500), vr("B", 2500),
		}
		got := computeSRM(results, feature, DefaultSRMThreshold)
		assert.Equal(t, eventcounter.SrmResult_OK, got.Status)
		assert.InDelta(t, 0.5, got.Variations[0].ExpectedWeight, 1e-9)
		assert.InDelta(t, 0.5, got.Variations[1].ExpectedWeight, 1e-9)
	})

	t.Run("audience-aware combined with unknown-variation leak: still MISMATCH (leak detector)", func(t *testing.T) {
		t.Parallel()
		// Compose the previous round's "leaked variation" fix with the
		// audience-aware fix: audience=50% default=A on a 50/50 A/B
		// rollout, but observed also has 1,000 users on vid-leak. The
		// expected fractions are A=0.75, B=0.25 (audience-adjusted), so
		// against totalObserved=10,000 we'd expect 7,500 A / 2,500 B /
		// 0 vid-leak. The leak inflates the denominator beyond the
		// 9,000 we'd expect under just the audience adjustment → MISMATCH.
		feature := newRolloutFeatureWithAudience(t, 50, "A",
			struct {
				id     string
				weight int32
			}{"A", 50}, struct {
				id     string
				weight int32
			}{"B", 50})
		results := []*eventcounter.VariationResult{
			vr("A", 6750), vr("B", 2250), vr("vid-leak", 1000),
		}
		got := computeSRM(results, feature, DefaultSRMThreshold)
		assert.Equal(t, eventcounter.SrmResult_MISMATCH, got.Status,
			"leaked traffic on top of audience filtering must still be detected")
		// vid-leak must surface in the breakdown with expected_weight=0.
		if assert.Len(t, got.Variations, 3) {
			// Sorted: A, B, vid-leak.
			assert.Equal(t, "vid-leak", got.Variations[2].VariationId)
			assert.InDelta(t, 0.0, got.Variations[2].ExpectedWeight, 1e-9)
			assert.EqualValues(t, 1000, got.Variations[2].ObservedUserCount)
		}
	})
}

// TestComputeSRM_ZeroExpectedCellLeakDetection covers the dedicated leak
// detector that runs alongside chi-square. Any variation with
// expected_weight == 0 (either an explicit zero-weight in the rollout or a
// variation absent from the rollout entirely) receiving non-trivial observed
// traffic is a bucketing bug — the SDK routed users somewhere the rollout
// says should get 0%. Chi-square can't model expected=0 cells, so this
// detector exists to close that gap.
//
// Threshold semantics: MISMATCH fires when leaked observed users strictly
// exceed max(leakNoiseFloor=5, leakRateFloor=0.0001 * totalObserved).
func TestComputeSRM_ZeroExpectedCellLeakDetection(t *testing.T) {
	t.Parallel()

	t.Run("zero-weight rollout variation getting traffic: MISMATCH (formerly silently SKIPPED)", func(t *testing.T) {
		t.Parallel()
		// Rollout says [A:100, B:0, C:0] — only A should receive any
		// traffic at all. Before this fix, B/C receiving traffic
		// triggered df=0 → SKIPPED "too few cells", hiding the bug.
		// The leak detector now catches it.
		feature := newRolloutFeature(t,
			struct {
				id     string
				weight int32
			}{"A", 100}, struct {
				id     string
				weight int32
			}{"B", 0}, struct {
				id     string
				weight int32
			}{"C", 0})
		results := []*eventcounter.VariationResult{
			vr("A", 9000), vr("B", 50), vr("C", 50),
		}
		got := computeSRM(results, feature, DefaultSRMThreshold)

		assert.Equal(t, eventcounter.SrmResult_MISMATCH, got.Status,
			"100 users routed to variations explicitly weighted 0 must be flagged")
		// Per-variation breakdown still surfaces the offending cells.
		// Sort order: A, B, C.
		if assert.Len(t, got.Variations, 3) {
			assert.Equal(t, "A", got.Variations[0].VariationId)
			assert.InDelta(t, 1.0, got.Variations[0].ExpectedWeight, 1e-9)
			assert.Equal(t, "B", got.Variations[1].VariationId)
			assert.InDelta(t, 0.0, got.Variations[1].ExpectedWeight, 1e-9)
			assert.EqualValues(t, 50, got.Variations[1].ObservedUserCount)
			assert.Equal(t, "C", got.Variations[2].VariationId)
			assert.InDelta(t, 0.0, got.Variations[2].ExpectedWeight, 1e-9)
			assert.EqualValues(t, 50, got.Variations[2].ObservedUserCount)
		}
		// chi-square has only 1 positive-expected cell (A), so df=0 and
		// the chi-square fields stay at zero values — but the leak
		// detector still produced a MISMATCH.
		assert.EqualValues(t, 0, got.DegreesOfFreedom,
			"with only one positive-expected cell, chi-square has df=0 and is not reported")
		assert.Empty(t, got.SkipReason,
			"MISMATCH should not populate skip_reason")
	})

	t.Run("small leak on top of well-balanced rollout: MISMATCH (formerly silently OK)", func(t *testing.T) {
		t.Parallel()
		// Rollout [A:50, B:50]. Observed perfectly balances A and B
		// (4975/4975 → chi-square ≈ 0.25, p ≈ 0.62, would be OK), but
		// 50 users leaked to unconfigured variation D. Chi-square
		// passes; the leak detector catches it.
		feature := newRolloutFeature(t,
			struct {
				id     string
				weight int32
			}{"A", 50}, struct {
				id     string
				weight int32
			}{"B", 50})
		results := []*eventcounter.VariationResult{
			vr("A", 4975), vr("B", 4975), vr("D", 50),
		}
		got := computeSRM(results, feature, DefaultSRMThreshold)

		assert.Equal(t, eventcounter.SrmResult_MISMATCH, got.Status,
			"a small leak to an unconfigured variation must NOT be hidden by a balanced main split")
		// Chi-square is still reported alongside the leak signal. With
		// the small deviation (4975 vs expected 4987.5), p should be
		// comfortably above the 0.001 threshold — proving the MISMATCH
		// came from the leak detector, not chi-square.
		assert.Greater(t, got.PValue, 0.001,
			"chi-square should pass; the MISMATCH must come from the leak detector")
		assert.EqualValues(t, 1, got.DegreesOfFreedom,
			"chi-square stays visible (df=1) alongside the leak signal")
	})

	t.Run("tiny leak below noise floor: still OK (no false positive)", func(t *testing.T) {
		t.Parallel()
		// Same shape, but only 3 leaked users — below the absolute
		// noise floor of 5. Could be a transient race during config
		// rollout, stale SDK cache, etc. Don't false-positive on this.
		feature := newRolloutFeature(t,
			struct {
				id     string
				weight int32
			}{"A", 50}, struct {
				id     string
				weight int32
			}{"B", 50})
		results := []*eventcounter.VariationResult{
			vr("A", 4998), vr("B", 4999), vr("D", 3),
		}
		got := computeSRM(results, feature, DefaultSRMThreshold)
		assert.Equal(t, eventcounter.SrmResult_OK, got.Status,
			"3 leaked users below the 5-user noise floor must not trigger MISMATCH")
	})

	t.Run("leak exactly at threshold: still OK (strict-greater semantics)", func(t *testing.T) {
		t.Parallel()
		// Leaked = 5 (exactly the noise floor for small n). The
		// detector uses strict-greater (>) not >=, so the threshold
		// itself is treated as "still in noise territory".
		feature := newRolloutFeature(t,
			struct {
				id     string
				weight int32
			}{"A", 50}, struct {
				id     string
				weight int32
			}{"B", 50})
		results := []*eventcounter.VariationResult{
			vr("A", 4998), vr("B", 4997), vr("D", 5),
		}
		got := computeSRM(results, feature, DefaultSRMThreshold)
		assert.Equal(t, eventcounter.SrmResult_OK, got.Status,
			"leak == noise floor (5) must NOT trigger MISMATCH; threshold uses >, not >=")
	})

	t.Run("leak threshold scales with totalObserved: 50 leaked out of 1M is OK", func(t *testing.T) {
		t.Parallel()
		// n = 1,000,000 → rate floor (0.01%) yields 100, which exceeds
		// the noise floor of 5 → effective threshold = 100. A 50-user
		// leak is below the rate floor so must NOT MISMATCH.
		feature := newRolloutFeature(t,
			struct {
				id     string
				weight int32
			}{"A", 50}, struct {
				id     string
				weight int32
			}{"B", 50})
		results := []*eventcounter.VariationResult{
			vr("A", 499975), vr("B", 499975), vr("D", 50),
		}
		got := computeSRM(results, feature, DefaultSRMThreshold)
		assert.Equal(t, eventcounter.SrmResult_OK, got.Status,
			"50 leaked out of 1M (0.005%) is below the 0.01%% rate floor")
	})

	t.Run("leak threshold scales with totalObserved: 200 leaked out of 1M is MISMATCH", func(t *testing.T) {
		t.Parallel()
		// Same n=1M but 200 leaked users → 0.02%, above the rate floor.
		feature := newRolloutFeature(t,
			struct {
				id     string
				weight int32
			}{"A", 50}, struct {
				id     string
				weight int32
			}{"B", 50})
		results := []*eventcounter.VariationResult{
			vr("A", 499900), vr("B", 499900), vr("D", 200),
		}
		got := computeSRM(results, feature, DefaultSRMThreshold)
		assert.Equal(t, eventcounter.SrmResult_MISMATCH, got.Status,
			"200 leaked out of 1M (0.02%%) exceeds the 0.01%% rate floor")
	})

	t.Run("chi-square mismatch AND leak: MISMATCH with chi-square still reported", func(t *testing.T) {
		t.Parallel()
		// Both detectors fire. Verify chi-square fields stay populated
		// so the UI can show both signals to the user.
		feature := newRolloutFeature(t,
			struct {
				id     string
				weight int32
			}{"A", 50}, struct {
				id     string
				weight int32
			}{"B", 50})
		results := []*eventcounter.VariationResult{
			// Main split badly skewed (chi-square will MISMATCH) AND
			// 1000 users leaked to D.
			vr("A", 3000), vr("B", 6000), vr("D", 1000),
		}
		got := computeSRM(results, feature, DefaultSRMThreshold)

		assert.Equal(t, eventcounter.SrmResult_MISMATCH, got.Status)
		// Chi-square stayed visible.
		assert.Greater(t, got.ChiSquare, 0.0,
			"chi-square statistic must stay reported when both checks fire")
		assert.EqualValues(t, 1, got.DegreesOfFreedom)
		assert.Less(t, got.PValue, DefaultSRMThreshold,
			"chi-square's own p-value should be below threshold here too")
	})

	t.Run("no leak, balanced rollout: OK (sanity)", func(t *testing.T) {
		t.Parallel()
		feature := newRolloutFeature(t,
			struct {
				id     string
				weight int32
			}{"A", 50}, struct {
				id     string
				weight int32
			}{"B", 50})
		results := []*eventcounter.VariationResult{
			vr("A", 5000), vr("B", 5000),
		}
		got := computeSRM(results, feature, DefaultSRMThreshold)
		assert.Equal(t, eventcounter.SrmResult_OK, got.Status)
	})
}
