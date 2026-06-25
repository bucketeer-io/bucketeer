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

package experimentcalc

import (
	"math"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/stat/distuv"

	"github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
)

// ---------------------------------------------------------------------------
// cvrBayesFactor unit tests
// ---------------------------------------------------------------------------

func TestCvrBayesFactor_NoData(t *testing.T) {
	t.Parallel()
	// No data in either arm → no evidence → BF = 1.
	assert.Equal(t, 1.0, cvrBayesFactor(0, 0, 0, 0))
	assert.Equal(t, 1.0, cvrBayesFactor(5, 0, 3, 10))
	assert.Equal(t, 1.0, cvrBayesFactor(5, 10, 3, 0))
}

func TestCvrBayesFactor_InvalidInputs(t *testing.T) {
	t.Parallel()
	// Successes > trials: invalid → BF = 1.
	assert.Equal(t, 1.0, cvrBayesFactor(15, 10, 5, 10))
	// Negative successes: invalid → BF = 1.
	assert.Equal(t, 1.0, cvrBayesFactor(-1, 10, 5, 10))
}

func TestCvrBayesFactor_EqualArms_NullFavoured(t *testing.T) {
	t.Parallel()
	// Both arms identical (50% CVR each) → data strongly favour H₀ → BF < 1.
	bf := cvrBayesFactor(50, 100, 50, 100)
	assert.True(t, bf < 1.0, "equal arms should give BF < 1 (H₀ favoured), got %f", bf)
	assert.True(t, bf > 0.0, "BF must be positive")
}

func TestCvrBayesFactor_FullConversionBothArms(t *testing.T) {
	t.Parallel()
	// Both arms at 100% conversion: logBF = log B(51,1) + log B(51,1) - log B(101,1)
	//   = -log(51) - log(51) + log(101) = log(101/51²) ≈ log(0.039)
	// So BF ≈ 0.039 — strong evidence for H₀.
	bf := cvrBayesFactor(50, 50, 50, 50)
	assert.InDelta(t, 0.039, bf, 0.002,
		"100%% CVR in both arms should give BF ≈ 0.039 (strong H₀), got %f", bf)
}

func TestCvrBayesFactor_LargeTreatmentEffect_BFAboveThreshold(t *testing.T) {
	t.Parallel()
	// 10% vs 30% CVR at n=500 each — strong H₁ signal.
	// BF should comfortably exceed the stopping threshold of 20.
	bf := cvrBayesFactor(50, 500, 150, 500)
	assert.GreaterOrEqual(t, bf, DefaultSequentialBFThreshold,
		"large CVR difference should give BF >= threshold, got %f", bf)
}

func TestCvrBayesFactor_ClearEffect_BFAboveThreshold(t *testing.T) {
	t.Parallel()
	// 10% vs 20% CVR at n=500 each — clear effect, BF must exceed threshold.
	bf := cvrBayesFactor(50, 500, 100, 500)
	assert.GreaterOrEqual(t, bf, DefaultSequentialBFThreshold,
		"clear CVR difference at n=500 should exceed threshold, got %f", bf)
}

func TestCvrBayesFactor_Symmetry(t *testing.T) {
	t.Parallel()
	// BF(A vs B) should equal BF(B vs A) — the BF is symmetric in H₀/H₁
	// framing because both arms appear symmetrically in the formula.
	bf1 := cvrBayesFactor(30, 100, 60, 100)
	bf2 := cvrBayesFactor(60, 100, 30, 100)
	assert.InDelta(t, bf1, bf2, 1e-9, "BF should be symmetric in arm order")
}

// ---------------------------------------------------------------------------
// valueBayesFactor unit tests
// ---------------------------------------------------------------------------

func TestValueBayesFactor_NoData(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 1.0, valueBayesFactor(0, 10.0, 1.0, 10, 10.0, 1.0))
	assert.Equal(t, 1.0, valueBayesFactor(10, 10.0, 1.0, 0, 10.0, 1.0))
}

func TestValueBayesFactor_EqualArms_NullFavoured(t *testing.T) {
	t.Parallel()
	// Both arms identical — data should favour H₀ → BF < 1.
	bf := valueBayesFactor(100, 10.0, 2.0, 100, 10.0, 2.0)
	assert.True(t, bf < 1.0, "equal arms should give BF < 1 (H₀ favoured), got %f", bf)
	assert.True(t, bf > 0.0, "BF must be positive")
}

func TestValueBayesFactor_LargeDifference_BFAboveThreshold(t *testing.T) {
	t.Parallel()
	// Large mean difference (10 vs 15) at n=50 each with small variance.
	// This mirrors the n=50 e2e fixture and should decisively exceed threshold.
	bf := valueBayesFactor(50, 10.05, 0.00255, 50, 15.05, 0.00255)
	assert.GreaterOrEqual(t, bf, DefaultSequentialBFThreshold,
		"large value difference at n=50 should exceed threshold, got %f", bf)
	// The BF should be extremely large (>>10^6) given the huge effect size.
	assert.GreaterOrEqual(t, math.Log(bf), 10.0,
		"log(BF) should be >> 10 for near-perfect separation, got log(BF)=%f", math.Log(bf))
}

func TestValueBayesFactor_ModerateEffect(t *testing.T) {
	t.Parallel()
	// Mean difference of 1 at var=4, n=200 each:
	// SE = sqrt(4/200 + 4/200) = sqrt(0.04) = 0.2, t = 1/0.2 = 5, n_eff = 100.
	// BF = sqrt(1/101) * exp(25*100/202) ≈ 0.0995 * exp(12.38) ≈ 23,500 >> 20.
	bf := valueBayesFactor(200, 10.0, 4.0, 200, 11.0, 4.0)
	assert.GreaterOrEqual(t, bf, DefaultSequentialBFThreshold,
		"moderate value difference at n=200 should exceed threshold, got %f", bf)
}

func TestValueBayesFactor_Symmetry(t *testing.T) {
	t.Parallel()
	bf1 := valueBayesFactor(100, 10.0, 2.0, 100, 12.0, 2.0)
	bf2 := valueBayesFactor(100, 12.0, 2.0, 100, 10.0, 2.0)
	// The pooled H₀ model is symmetric; small floating-point differences are OK.
	assert.InDelta(t, bf1, bf2, bf1*1e-9, "value BF should be symmetric in arm order")
}

func TestValueBayesFactor_ZeroVarianceInput_FallsBackToOne(t *testing.T) {
	t.Parallel()
	// Zero variance can occur when n=1 — must not panic; BF = 1 is acceptable.
	bf := valueBayesFactor(1, 10.0, 0.0, 1, 15.0, 0.0)
	assert.True(t, isFiniteFloat(bf) && bf > 0,
		"zero variance (n=1) must return a finite positive BF, got %f", bf)
}

// ---------------------------------------------------------------------------
// valueBayesFactor t-statistic formula verification
// ---------------------------------------------------------------------------

func TestValueBayesFactor_KnownTStatistic(t *testing.T) {
	t.Parallel()
	// Verify the closed-form: t=5, n_eff=100, r=1:
	//   BF = sqrt(1/101) * exp(25*100/202) = sqrt(0.0099) * exp(12.376)
	//   = 0.09950 * 236,895 ≈ 23,571
	bf := valueBayesFactor(200, 10.0, 4.0, 200, 11.0, 4.0)
	// SE = sqrt(4/200 + 4/200) = 0.2, t = 1.0/0.2 = 5, n_eff = 100
	expected := math.Sqrt(1.0/101.0) * math.Exp(25.0*100.0/202.0)
	assert.InDelta(t, expected, bf, expected*1e-9,
		"BF should match closed-form computation")
}

func TestValueBayesFactor_GrowsWithSampleSize(t *testing.T) {
	t.Parallel()
	// At a fixed mean difference, BF should increase monotonically with n.
	bf100 := valueBayesFactor(100, 10.0, 4.0, 100, 10.5, 4.0)
	bf500 := valueBayesFactor(500, 10.0, 4.0, 500, 10.5, 4.0)
	bf2000 := valueBayesFactor(2000, 10.0, 4.0, 2000, 10.5, 4.0)
	assert.Greater(t, bf500, bf100, "BF should grow with n (100→500)")
	assert.Greater(t, bf2000, bf500, "BF should grow with n (500→2000)")
}

func TestValueBayesFactor_UnderNullExpectedBFNearOneForLargeN(t *testing.T) {
	t.Parallel()
	// When both arms have the same mean, BF < 1 (H₀ favoured).
	// E[BF | H₀, t~N(0,1)] = 1 analytically; with t=0, BF = sqrt(1/(1+n_eff)) < 1.
	bf := valueBayesFactor(500, 10.0, 4.0, 500, 10.0, 4.0)
	assert.Less(t, bf, 1.0, "identical arms should give BF < 1")
	assert.Greater(t, bf, 0.0, "BF must be positive")
}

// ---------------------------------------------------------------------------
// fillSequentialBayesFactors unit tests
// ---------------------------------------------------------------------------

func TestFillSequentialBayesFactors_BaselineIsOne(t *testing.T) {
	t.Parallel()
	vrs := []*eventcounter.VariationResult{
		{
			VariationId:     "baseline",
			EvaluationCount: &eventcounter.VariationCount{UserCount: 100},
			ExperimentCount: &eventcounter.VariationCount{
				UserCount: 50, ValueSumPerUserMean: 10.0, ValueSumPerUserVariance: 2.0,
			},
		},
		{
			VariationId:     "treatment",
			EvaluationCount: &eventcounter.VariationCount{UserCount: 100},
			ExperimentCount: &eventcounter.VariationCount{
				UserCount: 60, ValueSumPerUserMean: 12.0, ValueSumPerUserVariance: 2.0,
			},
		},
	}
	fillSequentialBayesFactors(vrs, "baseline")

	assert.Equal(t, 1.0, vrs[0].CvrSequentialBayesFactor,
		"baseline CVR BF must be exactly 1")
	assert.Equal(t, 1.0, vrs[0].ValueSequentialBayesFactor,
		"baseline value BF must be exactly 1")
}

func TestFillSequentialBayesFactors_TreatmentBFPopulated(t *testing.T) {
	t.Parallel()
	vrs := []*eventcounter.VariationResult{
		{
			VariationId:     "baseline",
			EvaluationCount: &eventcounter.VariationCount{UserCount: 100},
			ExperimentCount: &eventcounter.VariationCount{
				UserCount: 50, ValueSumPerUserMean: 10.0, ValueSumPerUserVariance: 2.0,
			},
		},
		{
			VariationId:     "treatment",
			EvaluationCount: &eventcounter.VariationCount{UserCount: 100},
			ExperimentCount: &eventcounter.VariationCount{
				UserCount: 70, ValueSumPerUserMean: 14.0, ValueSumPerUserVariance: 2.0,
			},
		},
	}
	fillSequentialBayesFactors(vrs, "baseline")

	// Treatment CVR BF: 70/100 vs 50/100 — treatment has higher CVR.
	assert.Greater(t, vrs[1].CvrSequentialBayesFactor, 1.0,
		"treatment with higher CVR should have BF > 1")
	// Treatment value BF: mean 14 vs 10 with n=50/70 — clear separation.
	assert.Greater(t, vrs[1].ValueSequentialBayesFactor, 1.0,
		"treatment with higher mean should have value BF > 1")
	assert.True(t, isFiniteFloat(vrs[1].CvrSequentialBayesFactor))
	assert.True(t, isFiniteFloat(vrs[1].ValueSequentialBayesFactor))
}

func TestFillSequentialBayesFactors_NilVariationSkipped(t *testing.T) {
	t.Parallel()
	vrs := []*eventcounter.VariationResult{
		nil,
		{
			VariationId:     "baseline",
			EvaluationCount: &eventcounter.VariationCount{UserCount: 100},
			ExperimentCount: &eventcounter.VariationCount{UserCount: 50},
		},
	}
	// Must not panic.
	require.NotPanics(t, func() {
		fillSequentialBayesFactors(vrs, "baseline")
	})
}

// TestFillSequentialBayesFactors_MultiArmValueSkippedWhenOneArmHasZeroStats
// mirrors the calcGoalResult all-or-nothing guard: if any variation in the
// goal has zero value stats, the value posterior is skipped for all variations.
// The value BF must be 1.0 for all treatment arms in that case, even if the
// treatment arm being examined has non-zero stats.
func TestFillSequentialBayesFactors_MultiArmValueSkippedWhenOneArmHasZeroStats(t *testing.T) {
	t.Parallel()
	vrs := []*eventcounter.VariationResult{
		{
			VariationId:     "baseline",
			EvaluationCount: &eventcounter.VariationCount{UserCount: 100},
			ExperimentCount: &eventcounter.VariationCount{
				UserCount: 50, ValueSumPerUserMean: 10.0, ValueSumPerUserVariance: 2.0,
			},
		},
		{
			VariationId:     "treatment-A",
			EvaluationCount: &eventcounter.VariationCount{UserCount: 100},
			ExperimentCount: &eventcounter.VariationCount{
				// Good stats — value BF would be large if computed in isolation.
				UserCount: 70, ValueSumPerUserMean: 14.0, ValueSumPerUserVariance: 2.0,
			},
		},
		{
			VariationId:     "treatment-B",
			EvaluationCount: &eventcounter.VariationCount{UserCount: 100},
			ExperimentCount: &eventcounter.VariationCount{
				// Zero variance: calcGoalResult would skip value for the whole goal.
				UserCount: 60, ValueSumPerUserMean: 12.0, ValueSumPerUserVariance: 0.0,
			},
		},
	}
	fillSequentialBayesFactors(vrs, "baseline")

	vrA := vrs[1]
	vrB := vrs[2]

	// CVR BFs are unaffected by value stats.
	assert.Greater(t, vrA.CvrSequentialBayesFactor, 1.0,
		"treatment-A with higher CVR should still have CVR BF > 1")

	// Value BFs must both be 1.0 because treatment-B has zero variance —
	// the value posterior was skipped for the whole goal, so no value BF
	// should fire and set ValueSafeToStop=true.
	assert.Equal(t, 1.0, vrA.ValueSequentialBayesFactor,
		"treatment-A value BF must be 1.0 when any arm has zero variance (goal-wide skip)")
	assert.Equal(t, 1.0, vrB.ValueSequentialBayesFactor,
		"treatment-B value BF must be 1.0 (zero-variance arm)")
}

// ---------------------------------------------------------------------------
// computeSafeToStop unit tests
// ---------------------------------------------------------------------------

func TestComputeSafeToStop_ReturnsTrueWhenTreatmentBFAboveThreshold(t *testing.T) {
	t.Parallel()
	vrs := []*eventcounter.VariationResult{
		{VariationId: "baseline", CvrSequentialBayesFactor: 1.0},
		{VariationId: "treatment", CvrSequentialBayesFactor: 25.0},
	}
	result := computeSafeToStop(
		vrs,
		func(vr *eventcounter.VariationResult) float64 { return vr.CvrSequentialBayesFactor },
		"baseline",
		DefaultSequentialBFThreshold,
	)
	assert.True(t, result)
}

func TestComputeSafeToStop_ReturnsFalseWhenBelowThreshold(t *testing.T) {
	t.Parallel()
	vrs := []*eventcounter.VariationResult{
		{VariationId: "baseline", CvrSequentialBayesFactor: 1.0},
		{VariationId: "treatment", CvrSequentialBayesFactor: 5.0},
	}
	result := computeSafeToStop(
		vrs,
		func(vr *eventcounter.VariationResult) float64 { return vr.CvrSequentialBayesFactor },
		"baseline",
		DefaultSequentialBFThreshold,
	)
	assert.False(t, result)
}

func TestComputeSafeToStop_BaselineExcludedEvenIfBFHigh(t *testing.T) {
	t.Parallel()
	// Artificially set baseline's BF to a huge value — must still return false
	// because the baseline is excluded from the stopping signal.
	vrs := []*eventcounter.VariationResult{
		{VariationId: "baseline", CvrSequentialBayesFactor: 1000.0},
		{VariationId: "treatment", CvrSequentialBayesFactor: 3.0},
	}
	result := computeSafeToStop(
		vrs,
		func(vr *eventcounter.VariationResult) float64 { return vr.CvrSequentialBayesFactor },
		"baseline",
		DefaultSequentialBFThreshold,
	)
	assert.False(t, result, "baseline must not contribute to the stopping signal")
}

func TestComputeSafeToStop_ExactlyAtThreshold(t *testing.T) {
	t.Parallel()
	// BF == threshold should trigger safe-to-stop (≥, not >).
	vrs := []*eventcounter.VariationResult{
		{VariationId: "baseline", CvrSequentialBayesFactor: 1.0},
		{VariationId: "treatment", CvrSequentialBayesFactor: DefaultSequentialBFThreshold},
	}
	result := computeSafeToStop(
		vrs,
		func(vr *eventcounter.VariationResult) float64 { return vr.CvrSequentialBayesFactor },
		"baseline",
		DefaultSequentialBFThreshold,
	)
	assert.True(t, result, "BF exactly at threshold must trigger safe-to-stop")
}

// ---------------------------------------------------------------------------
// Multi-θ FPR simulation (acceptance criterion)
//
// This is the primary acceptance test for Follow-up K: peeking every day must
// NOT inflate the false-positive rate above a documented bound.
//
// Statistical guarantee recap:
//   The Beta-Binomial BF is a martingale under the Beta(1,1) prior-predictive
//   measure for H₀. Ville's inequality therefore gives:
//
//     P(max_t BF(t) ≥ K | H₀, θ integrated over Beta(1,1)) ≤ 1/K
//
//   This is a Bayesian-averaged bound, NOT a per-θ worst-case bound. At fixed
//   extreme θ values the per-θ error may deviate from 1/K. We test a grid of
//   θ values and FAIL the test if any per-θ FPR exceeds the relaxed tolerance,
//   so that we can detect when a GROW/e-variable construction is needed.
//
// Baseline comparison (demonstrates the problem we are solving):
//   The naive rule "stop when P(θ_B > θ_A | data) > 0.95" is simulated with
//   a Normal approximation to the Beta posteriors (accurate for n >> 1).
//   Its FPR should be substantially inflated relative to the BF rule,
//   confirming that the BF rule is strictly better at controlling FPR.
// ---------------------------------------------------------------------------

func TestSequentialBFControlsOptionalStoppingFPR(t *testing.T) {
	const (
		numExperiments = 10_000
		nPerDay        = 50 // per arm per day
		numDays        = 28
		threshold      = DefaultSequentialBFThreshold

		// Per-θ loose tolerance: the Bayesian-averaged guarantee holds on
		// average across the prior; per-θ deviations are expected but should
		// not be catastrophic. If any θ exceeds 12%, flag it as a warning that
		// a proper e-variable construction may be needed for that regime.
		perThetaWarnLevel = 0.12

		// Average FPR across tested θ values (unweighted) must stay at or
		// below this bound, corresponding to ≤1/K=5% plus simulation slack.
		avgFPRLimit = 0.065
	)

	// Null rates to probe. Cover low (near-0), moderate, and high (near-1)
	// regimes. Per point #1 in the approval note: test multiple θ values.
	nullRates := []float64{0.01, 0.05, 0.10, 0.30, 0.50}

	rng := rand.New(rand.NewPCG(0xDEAD_BEEF, 0xCAFE_BABE))

	type result struct {
		theta float64
		bfFPR float64 // false-positive rate under BF-based stopping
		nvFPR float64 // false-positive rate under naive 0.95-threshold stopping
	}
	results := make([]result, 0, len(nullRates))

	for _, theta := range nullRates {
		bfFalsePositives := 0
		naiveFalsePositives := 0

		for i := 0; i < numExperiments; i++ {
			var sA, nA, sB, nB int64
			bfTriggered := false
			naiveTriggered := false

			for day := 0; day < numDays; day++ {
				// Add one day's worth of Binomial(nPerDay, theta) increments.
				for j := 0; j < nPerDay; j++ {
					nA++
					if rng.Float64() < theta {
						sA++
					}
					nB++
					if rng.Float64() < theta {
						sB++
					}
				}

				if !bfTriggered {
					bf := cvrBayesFactor(sA, nA, sB, nB)
					if bf >= threshold {
						bfFalsePositives++
						bfTriggered = true
					}
				}

				if !naiveTriggered {
					prob := cvrProbBeatBaselineNormalApprox(sB, nB, sA, nA)
					if prob > 0.95 {
						naiveFalsePositives++
						naiveTriggered = true
					}
				}
			}
		}

		bfFPR := float64(bfFalsePositives) / float64(numExperiments)
		nvFPR := float64(naiveFalsePositives) / float64(numExperiments)
		results = append(results, result{theta, bfFPR, nvFPR})

		t.Logf("theta=%.2f: BF FPR=%.4f (%d/%d), naive FPR=%.4f (%d/%d)",
			theta, bfFPR, bfFalsePositives, numExperiments,
			nvFPR, naiveFalsePositives, numExperiments)

		// Per-θ FPR check. We use a relaxed tolerance (perThetaWarnLevel)
		// because the 1/K guarantee is Bayesian-averaged, not per-θ.
		if bfFPR > perThetaWarnLevel {
			t.Errorf(
				"theta=%.2f: per-θ BF FPR=%.4f exceeds %.0f%% warn level. "+
					"Consider a GROW / e-variable construction for extreme θ.",
				theta, bfFPR, perThetaWarnLevel*100)
		}

		// The BF rule must be strictly better than the naive rule at some point.
		// (Both could be ~0 at very low θ; only assert when naive FPR is notable.)
		if nvFPR > 0.05 {
			assert.Less(t, bfFPR, nvFPR,
				"theta=%.2f: BF FPR should be lower than naive 0.95-threshold FPR", theta)
		}
	}

	// Bayesian-averaged FPR: unweighted mean across tested null rates.
	var sumBF float64
	for _, r := range results {
		sumBF += r.bfFPR
	}
	avgBFFPR := sumBF / float64(len(results))
	t.Logf("Bayesian-averaged FPR (unweighted mean over tested θ): %.4f", avgBFFPR)
	assert.LessOrEqualf(t, avgBFFPR, avgFPRLimit,
		"Bayesian-averaged BF FPR=%.4f should be ≤ %.3f (1/K + simulation slack)",
		avgBFFPR, avgFPRLimit)
}

// cvrProbBeatBaselineNormalApprox estimates P(θ_B > θ_A | data) using a
// Normal approximation to each Beta posterior (accurate for n >> 30).
// Used only in the simulation comparison baseline above.
func cvrProbBeatBaselineNormalApprox(sB, nB, sA, nA int64) float64 {
	if nA <= 0 || nB <= 0 {
		return 0.5
	}
	pA := float64(sA) / float64(nA)
	pB := float64(sB) / float64(nB)
	varA := pA * (1 - pA) / float64(nA)
	varB := pB * (1 - pB) / float64(nB)
	se := math.Sqrt(varA + varB)
	if se == 0 {
		if pB > pA {
			return 1.0
		}
		return 0.0
	}
	z := (pB - pA) / se
	norm := distuv.Normal{Mu: 0, Sigma: 1}
	return norm.CDF(z)
}

// ---------------------------------------------------------------------------
// Value-metric multi-θ simulation
// ---------------------------------------------------------------------------

func TestValueBFControlsOptionalStoppingFPR(t *testing.T) {
	const (
		numExperiments    = 5_000
		nPerDay           = 30
		numDays           = 28
		threshold         = DefaultSequentialBFThreshold
		perThetaWarnLevel = 0.12
		avgFPRLimit       = 0.07
	)

	// Null scenarios: both arms drawn from Normal(mu, sigma^2).
	type nullScenario struct {
		mu    float64
		sigma float64
	}
	scenarios := []nullScenario{
		{mu: 1.0, sigma: 0.5},
		{mu: 10.0, sigma: 2.0},
		{mu: 100.0, sigma: 30.0},
		{mu: 0.5, sigma: 0.1},
	}

	rng := rand.New(rand.NewPCG(0xABCD_1234, 0xEF01_5678))
	var sumFPR float64

	for _, sc := range scenarios {
		falsePositives := 0
		for i := 0; i < numExperiments; i++ {
			var nA, nB int64
			var sumA, sumSqA, sumB, sumSqB float64
			triggered := false

			for day := 0; day < numDays; day++ {
				for j := 0; j < nPerDay; j++ {
					xa := sc.mu + rng.NormFloat64()*sc.sigma
					xb := sc.mu + rng.NormFloat64()*sc.sigma
					nA++
					sumA += xa
					sumSqA += xa * xa
					nB++
					sumB += xb
					sumSqB += xb * xb
				}

				if !triggered {
					meanA, varA := welfordMeanVar(nA, sumA, sumSqA)
					meanB, varB := welfordMeanVar(nB, sumB, sumSqB)
					bf := valueBayesFactor(nA, meanA, varA, nB, meanB, varB)
					if bf >= threshold {
						falsePositives++
						triggered = true
					}
				}
			}
		}
		fpr := float64(falsePositives) / float64(numExperiments)
		sumFPR += fpr
		t.Logf("value null (mu=%.1f, sigma=%.1f): FPR=%.4f (%d/%d)",
			sc.mu, sc.sigma, fpr, falsePositives, numExperiments)

		if fpr > perThetaWarnLevel {
			t.Errorf(
				"value null (mu=%.1f, sigma=%.1f): FPR=%.4f exceeds %.0f%% warn level",
				sc.mu, sc.sigma, fpr, perThetaWarnLevel*100)
		}
	}

	avgFPR := sumFPR / float64(len(scenarios))
	t.Logf("value Bayesian-averaged FPR: %.4f", avgFPR)
	assert.LessOrEqualf(t, avgFPR, avgFPRLimit,
		"value Bayesian-averaged FPR=%.4f should be ≤ %.3f", avgFPR, avgFPRLimit)
}

// welfordMeanVar computes sample mean and variance from running sums.
// Uses Welford-equivalent formula: var = (sumSq - n*mean^2) / (n-1).
func welfordMeanVar(n int64, sum, sumSq float64) (mean, variance float64) {
	if n <= 0 {
		return 0, 0
	}
	mean = sum / float64(n)
	if n <= 1 {
		return mean, 0
	}
	variance = (sumSq - float64(n)*mean*mean) / float64(n-1)
	if variance < 0 {
		variance = 0 // floating-point noise at near-zero variance
	}
	return mean, variance
}
