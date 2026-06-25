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

	"github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
)

const (
	// DefaultSequentialBFThreshold is the Bayes Factor threshold for the
	// always-valid stopping rule. The martingale property of the Bayes Factor
	// under H₀ provides the following guarantee (Ville's inequality):
	//
	//   E[BF(t) | H₀ prior-predictive] = 1 at every look t
	//   ⟹ P(max_t BF(t) ≥ K | H₀, Bayesian-averaged) ≤ 1/K
	//
	// Setting K = 20 yields a Bayesian-averaged Type-I error rate of at most
	// 5%. This means the false-positive rate is controlled on average over the
	// H₀ prior on θ, NOT as a per-θ worst-case bound. At extreme fixed θ
	// values (near 0 or 1) the per-θ error may deviate from 5%; the
	// multi-θ simulation in sequential_bayes_factor_test.go audits this.
	//
	// Multiple treatment arms: the BF controls error PER comparison (treatment
	// vs baseline). With k treatment arms the family-wise error inflates roughly
	// linearly; full multiple-comparisons handling belongs to a future phase.
	//
	// References:
	//   Grünwald, De Heide & Koolen, "Safe testing", JRSS-B 2024.
	//   Turner et al., "A tutorial on Bayesian sequential tests using Bayes
	//   factors", Psychonomic Bulletin & Review, 2022.
	//   De Heide & Grünwald, "Why optional stopping can be harmless for the
	//   Bayesian", 2021.
	DefaultSequentialBFThreshold = 20.0

	// valueBFEffectSizeScale is the unit for the Normal prior on the
	// standardized difference used by valueBayesFactor. A value of 1.0 places
	// the H₁ prior at N(0, 1) on the standardized effect size δ/SE, which is
	// equivalent to the "unit information prior" (Kass & Wasserman 1995).
	// MUST be a fixed constant — deriving it from experiment data would void
	// the always-valid guarantee.
	valueBFEffectSizeScale = 1.0

	// minValueBFSampleSize is the minimum per-arm sample size required before
	// valueBayesFactor returns a meaningful BF. Below this floor the Normal
	// approximation to the t-distribution is unreliable: E[BF | H₀] > 1 by
	// a non-trivial amount that can cause spurious ValueSafeToStop=true
	// decisions on noisy early data. The FPR simulations are validated at
	// n_per_arm >= 30, matching this floor.
	minValueBFSampleSize = int64(30)
)

// cvrBayesFactor computes the Beta-Binomial Bayes Factor for the comparison
//
//	H₁: θ_A ~ Beta(1,1), θ_B ~ Beta(1,1)  independently
//	H₀: θ_A = θ_B = θ,  θ ~ Beta(1,1)
//
// The result BF₁₀ ≥ 0. Values >> 1 favour H₁ (arms differ); values << 1
// favour H₀ (arms identical).
//
// Always-valid property (Bayesian-averaged): under H₀ the sequence
// {BF(t)} is a martingale with respect to the Beta(1,1) prior-predictive
// measure, so P(max_t BF(t) ≥ K | H₀, θ integrated) ≤ 1/K. This is NOT
// a per-fixed-θ worst-case bound.
//
// Returns 1.0 when either arm has no data, or when inputs are invalid
// (negative counts, successes exceeding trials).
func cvrBayesFactor(sA, nA, sB, nB int64) float64 {
	if nA <= 0 || nB <= 0 {
		return 1.0
	}
	if sA < 0 || sA > nA || sB < 0 || sB > nB {
		return 1.0
	}
	// logBF = log B(sA+1, nA-sA+1) + log B(sB+1, nB-sB+1) - log B(sA+sB+1, nA+nB-sA-sB+1)
	// Computed in log-space via lgamma to avoid overflow for large n.
	logBF := logBetaFn(float64(sA)+1, float64(nA-sA)+1) +
		logBetaFn(float64(sB)+1, float64(nB-sB)+1) -
		logBetaFn(float64(sA+sB)+1, float64(nA+nB-sA-sB)+1)
	if math.IsNaN(logBF) || math.IsInf(logBF, 0) {
		return 1.0
	}
	return expBF(logBF)
}

// logBetaFn returns log B(a, b) = logΓ(a) + logΓ(b) − logΓ(a+b).
func logBetaFn(a, b float64) float64 {
	la, _ := math.Lgamma(a)
	lb, _ := math.Lgamma(b)
	lab, _ := math.Lgamma(a + b)
	return la + lb - lab
}

// valueBayesFactor computes the Normal-Normal Bayes Factor for testing whether
// the per-user value means of two arms differ, using the t-statistic approach:
//
//	H₀: δ = 0              (treatment effect is zero)
//	H₁: δ ~ Normal(0, r²) (unit-information prior on standardized effect)
//
// where δ = (μ_B − μ_A) / SE is the standardized mean difference and
// r = valueBFEffectSizeScale (fixed constant, = 1.0).
//
// The Bayes Factor has the closed-form expression:
//
//	BF₁₀ = sqrt(1/(1 + n_eff·r²)) · exp(t² · n_eff·r² / (2·(1 + n_eff·r²)))
//
// where:
//   - t  = (mean_B − mean_A) / SE  is the Welch two-sample t-statistic
//   - SE = sqrt(var_A/n_A + var_B/n_B)  is the Welch standard error
//   - n_eff = n_A·n_B / (n_A + n_B)  is the effective sample size
//
// Always-valid property: for n → ∞, t is approximately N(0,1) under H₀, and
// one can verify analytically that E[BF(t) | H₀, t ~ N(0,1)] = 1, i.e.
// {BF(t)} is a martingale under H₀ in the large-n limit (Kass & Wasserman
// 1995, §5). For finite n, the t-distribution tails introduce a small
// positive bias in E[BF | H₀]; the simulation in sequential_bayes_factor_test.go
// confirms FPR stays within the documented bounds for n_per_arm ≥ 30.
//
// This approach avoids Bartlett's paradox: the BF tests the DIFFERENCE δ
// relative to a fixed effect-size prior rather than comparing two absolute NIG
// models whose priors are sensitive to the distance between μ₀ and the data.
// The always-valid FPR guarantee requires valueBFEffectSizeScale to be a
// fixed constant — never derive it from experiment data.
//
// Returns 1.0 when either arm has fewer than minValueBFSampleSize users (the
// t-approximation is unreliable below this floor), when the SE is zero (no
// within-arm variance), or when any computation is non-finite.
func valueBayesFactor(nA int64, meanA, varA float64, nB int64, meanB, varB float64) float64 {
	// The t-statistic BF is only approximately always-valid for sufficiently
	// large per-arm sample sizes (validated in tests at n_per_arm >= 30).
	// Below this floor the t-distribution tails are heavy enough that
	// E[BF | H₀] > 1 non-trivially, which can trigger ValueSafeToStop=true
	// on noisy early data. Return 1.0 (no evidence) until n is large enough.
	const minN = minValueBFSampleSize
	if nA < minN || nB < minN {
		return 1.0
	}

	// Welch standard error of the mean difference.
	seSq := varA/float64(nA) + varB/float64(nB)
	if seSq <= 0 || !isFiniteFloat(seSq) {
		return 1.0
	}

	t := (meanB - meanA) / math.Sqrt(seSq)
	if !isFiniteFloat(t) {
		return 1.0
	}

	nEff := float64(nA) * float64(nB) / float64(nA+nB)
	r := valueBFEffectSizeScale
	nEffR2 := nEff * r * r

	logBF := 0.5*math.Log(1/(1+nEffR2)) + t*t*nEffR2/(2*(1+nEffR2))
	if !isFiniteFloat(logBF) {
		return 1.0
	}
	return expBF(logBF)
}

// fillSequentialBayesFactors populates CvrSequentialBayesFactor and
// ValueSequentialBayesFactor on each VariationResult using the per-arm
// sufficient statistics already stored in ExperimentCount / EvaluationCount.
//
// The baseline variation (identified by baselineID) is set to 1.0 for both
// metrics: the BF of the baseline against itself is 1 by construction and
// must not be treated as evidence of a treatment effect.
//
// BFs are computed once from the final accumulated data (all timestamps up to
// the current batch run). This function is called after the per-timestamp loop
// so that the VariationResult structs hold the most-recent sufficient stats.
func fillSequentialBayesFactors(
	vrs []*eventcounter.VariationResult,
	baselineID string,
) {
	// Find the baseline's sufficient statistics.
	var baseGoalN int64
	var baseEvalN int64
	var baseMean, baseVar float64
	for _, vr := range vrs {
		if vr == nil {
			continue
		}
		if vr.VariationId == baselineID {
			if vr.EvaluationCount != nil {
				baseEvalN = vr.EvaluationCount.UserCount
			}
			if vr.ExperimentCount != nil {
				baseGoalN = vr.ExperimentCount.UserCount
				baseMean = vr.ExperimentCount.ValueSumPerUserMean
				baseVar = vr.ExperimentCount.ValueSumPerUserVariance
			}
			break
		}
	}

	for _, vr := range vrs {
		if vr == nil {
			continue
		}
		if vr.VariationId == baselineID {
			// Baseline vs itself: BF = 1 by construction.
			vr.CvrSequentialBayesFactor = 1.0
			vr.ValueSequentialBayesFactor = 1.0
			continue
		}

		var evalN, goalN int64
		var mean, variance float64
		if vr.EvaluationCount != nil {
			evalN = vr.EvaluationCount.UserCount
		}
		if vr.ExperimentCount != nil {
			goalN = vr.ExperimentCount.UserCount
			mean = vr.ExperimentCount.ValueSumPerUserMean
			variance = vr.ExperimentCount.ValueSumPerUserVariance
		}

		// CVR BF: treatment successes (goal users) vs baseline, over
		// their respective trial counts (eval users).
		vr.CvrSequentialBayesFactor = cvrBayesFactor(goalN, evalN, baseGoalN, baseEvalN)

		// Value BF: valid only when both arms have non-zero sufficient
		// stats (the calculator skips value analysis when these are zero,
		// matching the guard in calcGoalResult).
		if goalN > 0 && mean != 0 && variance != 0 &&
			baseGoalN > 0 && baseMean != 0 && baseVar != 0 {
			vr.ValueSequentialBayesFactor = valueBayesFactor(
				goalN, mean, variance,
				baseGoalN, baseMean, baseVar,
			)
		} else {
			vr.ValueSequentialBayesFactor = 1.0
		}
	}
}

// computeSafeToStop returns true when at least one non-baseline variation has
// a Bayes Factor (selected by bfFn) at or above threshold, indicating that
// the sequential evidence criterion has been met for the corresponding metric.
//
// The baseline variation is excluded: its BF is 1 by construction and must
// not contribute to the stopping signal.
func computeSafeToStop(
	vrs []*eventcounter.VariationResult,
	bfFn func(*eventcounter.VariationResult) float64,
	baselineID string,
	threshold float64,
) bool {
	for _, vr := range vrs {
		if vr == nil || vr.VariationId == baselineID {
			continue
		}
		if bfFn(vr) >= threshold {
			return true
		}
	}
	return false
}

// expBF exponentiates a log-Bayes-Factor, clamping the result to
// math.MaxFloat64 rather than returning +Inf. +Inf is not representable in
// JSON (proto float64 fields are JSON-encoded) and can break downstream
// consumers; clamping is safe because any finite BF >> 20 already satisfies
// the stopping criterion.
func expBF(logBF float64) float64 {
	// math.Log(math.MaxFloat64) ≈ 709.78; beyond this Exp overflows to +Inf.
	if logBF >= math.Log(math.MaxFloat64) {
		return math.MaxFloat64
	}
	return math.Exp(logBF)
}

// isFiniteFloat reports whether v is a real finite number (not NaN or ±Inf).
func isFiniteFloat(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}
