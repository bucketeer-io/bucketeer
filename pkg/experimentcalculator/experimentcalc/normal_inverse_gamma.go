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
	"fmt"
	"math"
	"math/rand/v2"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"

	"github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
)

// Fallback prior hyper-parameters used when the observed inputs are too
// degenerate to derive a full empirical-Bayes prior:
//   - `fallbackPriorMean` (=0) is only used when no variation has any usable
//     sample (totalN == 0) or the pooled mean is non-finite. Whenever any
//     usable sample exists, `computeEmpiricalBayesPriors` keeps μ₀ at the
//     pooled observed mean.
//   - `fallbackPriorAlpha` / `fallbackPriorBeta` (=1, =1) are used when a
//     usable mean exists but the pooled within-variation variance is
//     undefined or non-positive (e.g. every variation has n ≤ 1, or every
//     within-variation sample variance is 0).
//
// All fallback values are deliberately weak (1 pseudo-observation for the
// mean; an Inv-Gamma(1, 1) variance prior — undefined mean, mode 0.5 — i.e. a
// scale-1, scale-only "weak/unit-scale" prior, not a fixed unit variance) so
// the prior dies off as 1/(1+n) for the mean and even faster for the variance.
const (
	fallbackPriorMean  = 0.0
	fallbackPriorKappa = 1.0
	fallbackPriorAlpha = 1.0
	fallbackPriorBeta  = 1.0
)

// Pseudo-count weights for the empirical-Bayes prior. We use one pseudo
// observation for the mean (kappa0) and a matched one-degree-of-freedom prior
// for the variance (alpha0=1, beta0=pooled_var). With these weights the prior
// influence is ~50% at n=1, ~17% at n=5, ~9% at n=10, and ~1% at n=99, so the
// prior anchors small-sample estimates at the data's own scale without
// distorting moderate or large samples — and recovers Fix #1's large-n
// concentration behaviour exactly.
const (
	empiricalPriorKappa = 1.0
	empiricalPriorAlpha = 1.0
)

type distr struct {
	mu    float64
	nu    float64
	alpha float64
	beta  float64
}

// normalInverseGamma computes the value-metric posterior summaries. src seeds
// the Monte Carlo sampling; pass nil to use the global RNG (production) or a
// seeded source for deterministic tests.
func normalInverseGamma(
	src rand.Source,
	vids []string,
	means, vars []float64,
	sizes []int64,
	baselineIdx, postGenNum int,
) map[string]*eventcounter.VariationResult {
	startTime := time.Now()
	variationNum := len(means)
	variationResults := make(map[string]*eventcounter.VariationResult, variationNum)
	sampleSeries := make([]series.Series, 0, variationNum)
	// Derive the prior from the observed data (empirical Bayes) so the prior
	// auto-adapts to whatever scale this metric lives on (sub-unit conversions,
	// dollars, yen, seconds, …) instead of being anchored to an arbitrary
	// hardcoded constant.
	priorMu, priorKappa, priorAlpha, priorBeta := computeEmpiricalBayesPriors(means, vars, sizes)
	for i := 0; i < variationNum; i++ {
		post := calcPosterior(
			sizes[i],
			means[i],
			vars[i],
			priorMu,
			priorKappa,
			priorAlpha,
			priorBeta,
		)
		nums := generateNormalGamma(src, postGenNum, post.mu, post.nu, post.alpha, post.beta)
		sampleSeries = append(sampleSeries, series.Floats(nums))
	}
	samples := dataframe.New(sampleSeries...)
	best := samples.Rapply(calcBest)
	beatBaseline := samples.Rapply(func(s series.Series) series.Series {
		return calcBeatBaseline(s, baselineIdx)
	})
	for i := 0; i < variationNum; i++ {
		col := fmt.Sprintf("X%d", i)
		vr := &eventcounter.VariationResult{
			GoalValueSumPerUserProb:             createValueSumProb(samples.Col(col)),
			GoalValueSumPerUserProbBest:         createValueSumProbBest(best.Col(col)),
			GoalValueSumPerUserProbBeatBaseline: createValueSumProbBeatBaseline(beatBaseline.Col(col)),
		}
		variationResults[vids[i]] = vr
	}
	calculationHistogram.WithLabelValues(normalInverseGammaMethod).Observe(time.Since(startTime).Seconds())
	return variationResults
}

// calcPosterior performs the conjugate Normal-Inverse-Gamma update.
// priorKappa is the prior pseudo-count for the mean (kappa_0); thisVar is the
// per-user sample variance (divided by n-1), so the sum of squared deviations
// is (n-1) * thisVar.
func calcPosterior(
	thisN int64,
	thisMu, thisVar float64,
	priorMu, priorKappa, priorAlpha, priorBeta float64) distr {
	n := float64(thisN)
	kappaN := priorKappa + n
	postMu := (priorKappa*priorMu + n*thisMu) / kappaN
	postAlpha := priorAlpha + (n / 2)
	sumSquaredDev := 0.0
	if thisN > 1 {
		sumSquaredDev = float64(thisN-1) * thisVar
	}
	postBeta := priorBeta +
		(0.5 * sumSquaredDev) +
		((priorKappa * n / kappaN) * (thisMu - priorMu) * (thisMu - priorMu) / 2)
	return distr{
		mu:    postMu,
		nu:    kappaN,
		alpha: postAlpha,
		beta:  postBeta,
	}
}

// computeEmpiricalBayesPriors derives weakly-informative Normal-Inverse-Gamma
// prior hyper-parameters from the observed per-variation summaries:
//
//	mu0    = pooled (sample-size-weighted) mean across variations
//	kappa0 = 1                              // 1 pseudo-observation for the mean
//	alpha0 = 1                              // 1 pseudo-dof for the variance
//	beta0  = pooled within-variation sample variance
//	         = Σ(n_i - 1) · s_i² / Σ(n_i - 1)
//
// Inputs `vars[i]` are already sample variances (divisor n_i - 1, as enforced
// by Fix #1's VAR_SAMP standardisation), so Σ(n_i - 1) · s_i² recovers the
// classical pooled sum of squared deviations.
//
// We pool across all variations rather than using the baseline so the prior is
// symmetric and does not silently pull treatment posteriors toward control.
//
// Fallback layers (the function never returns NaN, Inf, or a non-positive
// variance — downstream NIG sampling would blow up on any of those):
//
//   - per-variation: skip variations with n ≤ 0 or with non-finite mean
//     entirely; for variance pooling, additionally skip variations with
//     non-finite or negative variance (their mean still contributes), so a
//     single bad row cannot poison the pooled estimate;
//   - mean: if no variation contributes a usable sample (totalN == 0) or the
//     pooled mean is non-finite, fall back to the full generic prior;
//   - variance: if the pooled within-variation variance is undefined
//     (pooledDoF == 0) or non-positive / non-finite, keep μ₀ at the pooled
//     mean but fall back to fallbackPriorAlpha / fallbackPriorBeta for the
//     variance prior.
func computeEmpiricalBayesPriors(
	means, vars []float64,
	sizes []int64,
) (mu, kappa, alpha, beta float64) {
	var totalN int64
	var sumNX float64
	var pooledSS float64
	var pooledDoF int64
	for i, n := range sizes {
		if n <= 0 {
			continue
		}
		m := means[i]
		if math.IsNaN(m) || math.IsInf(m, 0) {
			continue
		}
		totalN += n
		sumNX += float64(n) * m
		if n > 1 {
			v := vars[i]
			if math.IsNaN(v) || math.IsInf(v, 0) || v < 0 {
				continue
			}
			pooledSS += float64(n-1) * v
			pooledDoF += n - 1
		}
	}
	if totalN == 0 {
		return fallbackPriorMean, fallbackPriorKappa, fallbackPriorAlpha, fallbackPriorBeta
	}
	pooledMean := sumNX / float64(totalN)
	if math.IsNaN(pooledMean) || math.IsInf(pooledMean, 0) {
		return fallbackPriorMean, fallbackPriorKappa, fallbackPriorAlpha, fallbackPriorBeta
	}
	if pooledDoF == 0 {
		return pooledMean, empiricalPriorKappa, fallbackPriorAlpha, fallbackPriorBeta
	}
	pooledVar := pooledSS / float64(pooledDoF)
	if pooledVar <= 0 || math.IsNaN(pooledVar) || math.IsInf(pooledVar, 0) {
		return pooledMean, empiricalPriorKappa, fallbackPriorAlpha, fallbackPriorBeta
	}
	return pooledMean, empiricalPriorKappa, empiricalPriorAlpha, pooledVar
}

func generateNormalGamma(src rand.Source, n int, mu float64, lambda float64, alpha float64, beta float64) []float64 {
	tauDist := distuv.Gamma{Alpha: alpha, Beta: beta, Src: src}

	tauSamples := make([]float64, n)
	for i := 0; i < n; i++ {
		tauSamples[i] = tauDist.Rand()
	}

	normDist := distuv.Normal{Mu: 0, Sigma: 1, Src: src}
	x := make([]float64, n)
	for i, tau := range tauSamples {
		x[i] = mu + math.Sqrt(1/(tau*lambda))*normDist.Rand()
	}

	return x
}

func calcBest(s series.Series) series.Series {
	max := s.Max()
	samples := s.Float()
	maxArray := make([]int, len(samples))
	for i := 0; i < len(samples); i++ {
		if samples[i] == max {
			maxArray[i] = 1
		} else {
			maxArray[i] = 0
		}
	}
	return series.Ints(maxArray)
}

func calcBeatBaseline(s series.Series, baselineIdx int) series.Series {
	baseline := s.Val(baselineIdx).(float64)
	samples := s.Float()
	beatArray := make([]int, len(samples))
	for i := 0; i < len(samples); i++ {
		if samples[i] > baseline {
			beatArray[i] = 1
		} else {
			beatArray[i] = 0
		}
	}
	return series.Ints(beatArray)
}

func createValueSumProb(samples series.Series) *eventcounter.DistributionSummary {
	ordered := samples.Subset(samples.Order(false)).Float()
	return &eventcounter.DistributionSummary{
		Median:        samples.Median(),
		Percentile025: stat.Quantile(0.025, stat.LinInterp, ordered, nil),
		Percentile975: stat.Quantile(0.975, stat.LinInterp, ordered, nil),
	}
}

func createValueSumProbBest(samples series.Series) *eventcounter.DistributionSummary {
	return &eventcounter.DistributionSummary{
		Mean: samples.Mean(),
	}
}

func createValueSumProbBeatBaseline(samples series.Series) *eventcounter.DistributionSummary {
	return &eventcounter.DistributionSummary{
		Mean: samples.Mean(),
	}
}
