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
//

package experimentcalc

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"

	"github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

const (
	priorMean  = 30
	priorVar   = 2
	priorSize  = 20
	priorAlpha = 10
	priorBeta  = 1000
)

type distr struct {
	mu    float64
	nu    float64
	alpha float64
	beta  float64
	n     int
}

func normalInverseGamma(
	ctx context.Context,
	vids []string,
	means, vars []float64,
	sizes []int64,
	baselineIdx, postGenNum int,
) map[string]*eventcounter.VariationResult {
	startTime := time.Now()
	variationNum := len(means)
	variationResults := make(map[string]*eventcounter.VariationResult, variationNum)
	sampleSeries := make([]series.Series, 0, variationNum)
	for i := 0; i < variationNum; i++ {
		post := calcPosterior(
			sizes[i],
			means[i],
			vars[i],
			priorSize,
			priorMean,
			priorVar,
			priorAlpha,
			priorBeta,
		)
		nums := generateNormalGamma(postGenNum, post.mu, post.nu, post.alpha, post.beta)
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

func calcPosterior(
	thisN int64,
	thisMu, thisSigma float64,
	priorN int64,
	priorMu, priorNu, priorAlpha, priorBeta float64) distr {
	retN := thisN + priorN
	n2 := math.Log(float64(thisN)) / math.Log(1.1)
	postMu := (priorNu*priorMu + n2*thisMu) / (priorNu + n2)
	postNu := priorNu + n2
	postAlpha := priorAlpha + (n2 / 2)
	postBeta := priorBeta +
		(0.5 * thisSigma * thisSigma * n2) +
		((n2 * priorNu / (priorNu * n2)) * ((thisMu - priorMu) * (thisMu - priorMu)) / 2)
	return distr{
		mu:    postMu,
		nu:    postNu,
		alpha: postAlpha,
		beta:  postBeta,
		n:     int(retN),
	}
}

func generateNormalGamma(n int, mu float64, lambda float64, alpha float64, beta float64) []float64 {
	tauDist := distuv.Gamma{Alpha: alpha, Beta: beta}

	tauSamples := make([]float64, n)
	for i := 0; i < n; i++ {
		tauSamples[i] = tauDist.Rand()
	}

	x := make([]float64, n)
	for i, tau := range tauSamples {
		x[i] = mu + math.Sqrt(1/(tau*lambda))*distuv.Normal{Mu: 0, Sigma: 1}.Rand()
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
