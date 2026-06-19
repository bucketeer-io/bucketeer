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
	"context"
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
	vrs := normalInverseGamma(context.TODO(), vids, means, vars, sizes, baselineIdx, 25000)

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
	vrs := normalInverseGamma(context.TODO(), vids, means, vars, sizes, baselineIdx, 25000)

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
