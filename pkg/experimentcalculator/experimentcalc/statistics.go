// Copyright 2025 The Bucketeer Authors.
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

	"github.com/go-gota/gota/dataframe"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"

	"github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

func createCvrProb(
	df dataframe.DataFrame,
	samples []dataframe.DataFrame,
	index int,
) *eventcounter.DistributionSummary {
	col := fmt.Sprintf("p.%d", index)
	p := df.Col(col)
	// calculate histogram
	nBins := 100
	dividers := make([]float64, nBins+1)
	ordered := p.Subset(p.Order(false)).Float()
	min := floats.Min(ordered)
	max := floats.Max(ordered)
	max++
	floats.Span(dividers, min, max)
	histFloats := stat.Histogram(nil, dividers, ordered, nil)
	hist := make([]int64, len(histFloats))
	for i, v := range histFloats {
		hist[i] = int64(v)
	}
	histogram := &eventcounter.Histogram{
		Hist: hist,
		Bins: dividers,
	}
	paramSamples := extractParamSample(samples, col)
	return &eventcounter.DistributionSummary{
		Mean:          p.Mean(),
		Sd:            p.StdDev(),
		Rhat:          rHat(paramSamples),
		Histogram:     histogram,
		Median:        p.Median(),
		Percentile025: stat.Quantile(0.025, stat.LinInterp, ordered, nil),
		Percentile975: stat.Quantile(0.975, stat.LinInterp, ordered, nil),
	}
}

func createCvrProbBest(
	df dataframe.DataFrame,
	samples []dataframe.DataFrame,
	index int,
) *eventcounter.DistributionSummary {
	col := fmt.Sprintf("prob_best.%d", index)
	probBest := df.Col(col)
	paramSamples := extractParamSample(samples, col)
	return &eventcounter.DistributionSummary{
		Mean: probBest.Mean(),
		Sd:   probBest.StdDev(),
		Rhat: rHat(paramSamples),
	}
}

func createCvrProbBeatBaseline(
	df dataframe.DataFrame,
	samples []dataframe.DataFrame,
	baselineIdx, index int,
) *eventcounter.DistributionSummary {
	probBeatBaseline := &eventcounter.DistributionSummary{}
	if baselineIdx == index {
		return probBeatBaseline
	}
	col := fmt.Sprintf("prob_upper.%d.%d", index, baselineIdx)
	probUpper := df.Col(col)
	paramSamples := extractParamSample(samples, col)
	probBeatBaseline.Mean = probUpper.Mean()
	probBeatBaseline.Sd = probUpper.StdDev()
	probBeatBaseline.Rhat = rHat(paramSamples)
	return probBeatBaseline
}

// rHat returns the split potential scale reduction (split R hat)
// for the specified parameter samples.
func rHat(samples [][]float64) float64 {
	// Check for empty input
	if len(samples) == 0 || len(samples[0]) == 0 {
		return 1.0 // Safe default: no samples implies convergence.
	}

	chains := len(samples)
	nsamples := len(samples[0])
	// Use the smallest chain length among all chains.
	for i := 1; i < chains; i++ {
		nsamples = int(math.Min(float64(nsamples), float64(len(samples[i]))))
	}
	// Ensure we have at least two samples per chain.
	if nsamples < 2 {
		return 1.0
	}
	// Use an even number of samples.
	if nsamples%2 == 1 {
		nsamples--
	}
	n := nsamples / 2

	// Allocate arrays for the split-chain means and variances.
	splitChainMean := make([]float64, 2*chains)
	splitChainVar := make([]float64, 2*chains)

	// Split each chain into two halves.
	for chain := 0; chain < chains; chain++ {
		splitChainMean[2*chain] = stat.Mean(samples[chain][:n], nil)
		splitChainMean[2*chain+1] = stat.Mean(samples[chain][n:], nil)

		splitChainVar[2*chain] = stat.Variance(samples[chain][:n], nil)
		splitChainVar[2*chain+1] = stat.Variance(samples[chain][n:], nil)
	}

	// Calculate between-chain and within-chain variances.
	varBetween := float64(n) * stat.Variance(splitChainMean, nil)
	varWithin := stat.Mean(splitChainVar, nil)

	// Use an epsilon to avoid division by a near-zero within-chain variance.
	epsilon := 1e-10
	if math.Abs(varWithin) < epsilon {
		// If both between and within are near zero, the chains are identical.
		if math.Abs(varBetween) < epsilon {
			return 1.0
		}
		// Otherwise, force varWithin to epsilon to avoid division by zero.
		varWithin = epsilon
	}

	// Compute the split R hat statistic using the reformulated expression:
	// rHat = sqrt(( (B/W) + (n-1) ) / n)
	rhat := math.Sqrt((varBetween/varWithin + float64(n-1)) / float64(n))

	// As a final safeguard, return 1.0 if rhat is not a finite number.
	if math.IsNaN(rhat) || math.IsInf(rhat, 0) {
		return 1.0
	}

	return rhat
}

func extractParamSample(samples []dataframe.DataFrame, col string) [][]float64 {
	paramSeries := make([][]float64, 0, len(samples))
	for _, sample := range samples {
		paramSeries = append(paramSeries, sample.Col(col).Float())
	}
	return paramSeries
}
