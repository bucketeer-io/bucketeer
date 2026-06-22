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
	"math"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

func TestNormalInverseGamma(t *testing.T) {
	t.Parallel()
	v1Mean, v1Sd := 12.0, 10.0
	v2Mean, v2Sd := 15.0, 12.0
	// Seed local RNG sources so the generated data is deterministic and the test
	// stays stable under t.Parallel() regardless of other concurrent tests.
	v1Dist := distuv.Normal{Mu: v1Mean, Sigma: v1Sd, Src: rand.NewPCG(1, 2)}
	v2Dist := distuv.Normal{Mu: v2Mean, Sigma: v2Sd, Src: rand.NewPCG(3, 4)}
	sampleNum := 20000
	v1 := make([]float64, sampleNum)
	v2 := make([]float64, sampleNum)

	for i := 0; i < sampleNum; i++ {
		if rand1 := v1Dist.Rand(); rand1 < 0 {
			v1[i] = 0
		} else {
			v1[i] = rand1
		}
		if rand2 := v2Dist.Rand(); rand2 < 0 {
			v2[i] = 0
		} else {
			v2[i] = rand2
		}
	}

	vids := []string{"vid1", "vid2"}
	means := []float64{stat.Mean(v1, nil), stat.Mean(v2, nil)}
	vars := []float64{stat.Variance(v1, nil), stat.Variance(v2, nil)}
	sizes := []int64{int64(len(v1)), int64(len(v2))}
	baselineIdx := 0
	vrs := normalInverseGamma(rand.NewPCG(9, 10), vids, means, vars, sizes, baselineIdx, 25000)

	// With the conjugate update using the real sample size, the posterior for the
	// mean concentrates around the observed per-user mean and the 95% credible
	// interval is narrow (well under 1 unit wide for ~20k samples), instead of the
	// previous behaviour where the interval spanned roughly [-2, 28].
	vid1 := vrs["vid1"]
	assert.Greater(t, vid1.GoalValueSumPerUserProb.Median, 11.5)
	assert.Less(t, vid1.GoalValueSumPerUserProb.Median, 13.5)
	assert.Less(
		t,
		vid1.GoalValueSumPerUserProb.Percentile975-vid1.GoalValueSumPerUserProb.Percentile025,
		1.5,
	)
	// vid1 (baseline) has the lower mean, so it is almost never the best variation.
	assert.Less(t, vid1.GoalValueSumPerUserProbBest.Mean, 0.05)
	assert.Equal(t, vid1.GoalValueSumPerUserProbBeatBaseline.Mean, 0.0)

	vid2 := vrs["vid2"]
	assert.Greater(t, vid2.GoalValueSumPerUserProb.Median, 14.5)
	assert.Less(t, vid2.GoalValueSumPerUserProb.Median, 16.5)
	assert.Less(
		t,
		vid2.GoalValueSumPerUserProb.Percentile975-vid2.GoalValueSumPerUserProb.Percentile025,
		1.5,
	)
	// vid2 has the clearly higher mean, so it is essentially certain to be the
	// best and to beat the baseline.
	assert.Greater(t, vid2.GoalValueSumPerUserProbBest.Mean, 0.95)
	assert.Greater(t, vid2.GoalValueSumPerUserProbBeatBaseline.Mean, 0.95)
}

// TestNormalInverseGammaHeavyTailed verifies the posterior stays sensible on
// heavy-tailed, strictly-positive value data (log-normal revenue per user):
// the posterior mean concentrates near the observed per-user mean, the credible
// interval is narrow and stays positive (no nonsensical negative-revenue
// bounds), and the higher-revenue variation is identified decisively.
func TestNormalInverseGammaHeavyTailed(t *testing.T) {
	t.Parallel()
	// True means are exp(mu + sigma^2/2): ~12.18 for vid1 and ~16.44 for vid2.
	// Seed local RNG sources for deterministic, parallel-safe sampling.
	v1Dist := distuv.LogNormal{Mu: 2.0, Sigma: 1.0, Src: rand.NewPCG(5, 6)}
	v2Dist := distuv.LogNormal{Mu: 2.3, Sigma: 1.0, Src: rand.NewPCG(7, 8)}
	sampleNum := 20000
	v1 := make([]float64, sampleNum)
	v2 := make([]float64, sampleNum)
	for i := 0; i < sampleNum; i++ {
		v1[i] = v1Dist.Rand()
		v2[i] = v2Dist.Rand()
	}

	vids := []string{"vid1", "vid2"}
	means := []float64{stat.Mean(v1, nil), stat.Mean(v2, nil)}
	vars := []float64{stat.Variance(v1, nil), stat.Variance(v2, nil)}
	sizes := []int64{int64(len(v1)), int64(len(v2))}
	baselineIdx := 0
	vrs := normalInverseGamma(rand.NewPCG(11, 12), vids, means, vars, sizes, baselineIdx, 25000)

	vid1 := vrs["vid1"]
	// Posterior median tracks the observed per-user mean.
	assert.InDelta(t, vid1.GoalValueSumPerUserProb.Median, means[0], 1.0)
	// Strictly-positive credible bounds (heavy tails must not push them negative).
	assert.Greater(t, vid1.GoalValueSumPerUserProb.Percentile025, 0.0)
	assert.Less(
		t,
		vid1.GoalValueSumPerUserProb.Percentile975-vid1.GoalValueSumPerUserProb.Percentile025,
		2.0,
	)
	assert.Less(t, vid1.GoalValueSumPerUserProbBest.Mean, 0.05)
	assert.Equal(t, vid1.GoalValueSumPerUserProbBeatBaseline.Mean, 0.0)

	vid2 := vrs["vid2"]
	assert.InDelta(t, vid2.GoalValueSumPerUserProb.Median, means[1], 1.0)
	assert.Greater(t, vid2.GoalValueSumPerUserProb.Percentile025, 0.0)
	assert.Less(
		t,
		vid2.GoalValueSumPerUserProb.Percentile975-vid2.GoalValueSumPerUserProb.Percentile025,
		2.0,
	)
	assert.Greater(t, vid2.GoalValueSumPerUserProbBest.Mean, 0.95)
	assert.Greater(t, vid2.GoalValueSumPerUserProbBeatBaseline.Mean, 0.95)
}

// TestNormalInverseGammaSmallNOffScale exercises the empirical-Bayes prior on
// a small sample (n=10) of an off-scale metric — per-user revenue around 5,000
// yen — to make sure the posterior tracks the data's natural scale instead of
// being yanked toward the hardcoded constant the previous implementation used.
//
// The test asserts that:
//  1. The posterior median is within ~10% of the per-variation observed mean
//     (a hardcoded prior at 30 would have biased small-n results downward by
//     hundreds-to-thousands of units here).
//  2. The 95% credible interval is on the order of the per-user SD, not blown
//     up by an absolute β₀=1000 prior that would dominate at this scale.
//  3. The clearly higher-mean variation is identified as the likely best.
func TestNormalInverseGammaSmallNOffScale(t *testing.T) {
	t.Parallel()
	const sampleNum = 10
	v1Dist := distuv.Normal{Mu: 5000, Sigma: 800, Src: rand.NewPCG(13, 14)}
	v2Dist := distuv.Normal{Mu: 6000, Sigma: 800, Src: rand.NewPCG(15, 16)}
	v1 := make([]float64, sampleNum)
	v2 := make([]float64, sampleNum)
	for i := 0; i < sampleNum; i++ {
		v1[i] = v1Dist.Rand()
		v2[i] = v2Dist.Rand()
	}

	vids := []string{"vid1", "vid2"}
	means := []float64{stat.Mean(v1, nil), stat.Mean(v2, nil)}
	vars := []float64{stat.Variance(v1, nil), stat.Variance(v2, nil)}
	sizes := []int64{int64(len(v1)), int64(len(v2))}
	vrs := normalInverseGamma(rand.NewPCG(17, 18), vids, means, vars, sizes, 0, 25000)

	vid1 := vrs["vid1"]
	assert.InDelta(t, vid1.GoalValueSumPerUserProb.Median, means[0], 0.10*means[0])
	v1Sd := math.Sqrt(vars[0])
	assert.Less(
		t,
		vid1.GoalValueSumPerUserProb.Percentile975-vid1.GoalValueSumPerUserProb.Percentile025,
		3.0*v1Sd,
	)
	vid2 := vrs["vid2"]
	assert.InDelta(t, vid2.GoalValueSumPerUserProb.Median, means[1], 0.10*means[1])
	v2Sd := math.Sqrt(vars[1])
	assert.Less(
		t,
		vid2.GoalValueSumPerUserProb.Percentile975-vid2.GoalValueSumPerUserProb.Percentile025,
		3.0*v2Sd,
	)
	assert.Greater(t, vid2.GoalValueSumPerUserProbBest.Mean, vid1.GoalValueSumPerUserProbBest.Mean)
	assert.Greater(t, vid2.GoalValueSumPerUserProbBeatBaseline.Mean, 0.5)
}

// TestNormalInverseGammaSmallNSubUnit exercises the empirical-Bayes prior on a
// sub-unit metric (per-user values around 0.3, e.g. value-sum-per-user for a
// low-value goal) at small n. Previously the hardcoded priors (mean=30,
// beta=1000) would have pulled the posterior median to a 1–2-order-of-magnitude
// wrong location and produced an enormous credible interval relative to the
// data.
func TestNormalInverseGammaSmallNSubUnit(t *testing.T) {
	t.Parallel()
	const sampleNum = 10
	v1Dist := distuv.Normal{Mu: 0.30, Sigma: 0.10, Src: rand.NewPCG(19, 20)}
	v2Dist := distuv.Normal{Mu: 0.35, Sigma: 0.10, Src: rand.NewPCG(21, 22)}
	v1 := make([]float64, sampleNum)
	v2 := make([]float64, sampleNum)
	for i := 0; i < sampleNum; i++ {
		v1[i] = v1Dist.Rand()
		v2[i] = v2Dist.Rand()
	}

	vids := []string{"vid1", "vid2"}
	means := []float64{stat.Mean(v1, nil), stat.Mean(v2, nil)}
	vars := []float64{stat.Variance(v1, nil), stat.Variance(v2, nil)}
	sizes := []int64{int64(len(v1)), int64(len(v2))}
	vrs := normalInverseGamma(rand.NewPCG(23, 24), vids, means, vars, sizes, 0, 25000)

	for i, vid := range vids {
		vr := vrs[vid]
		// Not pulled toward 30: median stays in the same order of magnitude as
		// the data (well under 1, not somewhere near the old prior).
		assert.Less(t, vr.GoalValueSumPerUserProb.Median, 1.0,
			"variation %s median should track the sub-unit metric scale", vid)
		assert.Greater(t, vr.GoalValueSumPerUserProb.Median, 0.0,
			"variation %s median should be positive", vid)
		// CI tracks the data's natural variability, not an absolute prior.
		assert.Less(
			t,
			vr.GoalValueSumPerUserProb.Percentile975-vr.GoalValueSumPerUserProb.Percentile025,
			1.0,
			"variation %s CI should be on the order of the data, not the old absolute prior", vid,
		)
		sd := math.Sqrt(vars[i])
		assert.InDelta(t, vr.GoalValueSumPerUserProb.Median, means[i], 5*sd)
	}
}

func TestComputeEmpiricalBayesPriors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		means     []float64
		vars      []float64
		sizes     []int64
		wantMu    float64
		wantKappa float64
		wantAlpha float64
		wantBeta  float64
	}{
		{
			name:      "pooled mean and variance across two variations",
			means:     []float64{10.0, 20.0},
			vars:      []float64{4.0, 9.0},
			sizes:     []int64{50, 50},
			wantMu:    15.0, // (50*10 + 50*20)/100
			wantKappa: 1.0,
			wantAlpha: 1.0,
			wantBeta:  (49*4.0 + 49*9.0) / 98, // pooled within-variation var
		},
		{
			name:      "size-weighted pooled mean",
			means:     []float64{10.0, 20.0},
			vars:      []float64{4.0, 4.0},
			sizes:     []int64{10, 90},
			wantMu:    (10*10.0 + 90*20.0) / 100, // 19.0
			wantKappa: 1.0,
			wantAlpha: 1.0,
			wantBeta:  (9*4.0 + 89*4.0) / 98, // ~4.0
		},
		{
			name:      "single variation: pooled mean and var equal that variation's",
			means:     []float64{7.5},
			vars:      []float64{2.25},
			sizes:     []int64{100},
			wantMu:    7.5,
			wantKappa: 1.0,
			wantAlpha: 1.0,
			wantBeta:  2.25,
		},
		{
			name:      "all sizes <= 1 falls back to weak variance prior but keeps pooled mean",
			means:     []float64{4.0, 6.0},
			vars:      []float64{0.0, 0.0},
			sizes:     []int64{1, 1},
			wantMu:    5.0, // (1*4 + 1*6)/2
			wantKappa: 1.0,
			wantAlpha: fallbackPriorAlpha,
			wantBeta:  fallbackPriorBeta,
		},
		{
			name:      "zero pooled variance falls back to weak variance prior but keeps pooled mean",
			means:     []float64{2.0, 2.0},
			vars:      []float64{0.0, 0.0},
			sizes:     []int64{10, 10},
			wantMu:    2.0,
			wantKappa: 1.0,
			wantAlpha: fallbackPriorAlpha,
			wantBeta:  fallbackPriorBeta,
		},
		{
			name:      "no usable variations falls back entirely",
			means:     []float64{0, 0},
			vars:      []float64{0, 0},
			sizes:     []int64{0, 0},
			wantMu:    fallbackPriorMean,
			wantKappa: fallbackPriorKappa,
			wantAlpha: fallbackPriorAlpha,
			wantBeta:  fallbackPriorBeta,
		},
		{
			name:      "ignores variations with n<=0 in both mean and variance pooling",
			means:     []float64{99.0, 10.0, 20.0},
			vars:      []float64{1000.0, 4.0, 9.0},
			sizes:     []int64{0, 50, 50},
			wantMu:    15.0,
			wantKappa: 1.0,
			wantAlpha: 1.0,
			wantBeta:  (49*4.0 + 49*9.0) / 98,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mu, kappa, alpha, beta := computeEmpiricalBayesPriors(tt.means, tt.vars, tt.sizes)
			assert.InDelta(t, tt.wantMu, mu, 1e-9, "mu")
			assert.Equal(t, tt.wantKappa, kappa, "kappa")
			assert.Equal(t, tt.wantAlpha, alpha, "alpha")
			assert.InDelta(t, tt.wantBeta, beta, 1e-9, "beta")
		})
	}
}
