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
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

func TestNormalInverseGamma(t *testing.T) {
	v1Mean, v1Sd := 12.0, 10.0
	v2Mean, v2Sd := 15.0, 12.0
	v1Dist := distuv.Normal{Mu: v1Mean, Sigma: v1Sd}
	v2Dist := distuv.Normal{Mu: v2Mean, Sigma: v2Sd}
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

	assert.Greater(t, vrs["vid1"].GoalValueSumPerUserProb.Median, 12.0)
	assert.Less(t, vrs["vid1"].GoalValueSumPerUserProb.Median, 13.5)
	assert.Greater(t, vrs["vid1"].GoalValueSumPerUserProb.Percentile025, -2.0)
	assert.Less(t, vrs["vid1"].GoalValueSumPerUserProb.Percentile025, -1.0)
	assert.Greater(t, vrs["vid1"].GoalValueSumPerUserProb.Percentile975, 26.5)
	assert.Less(t, vrs["vid1"].GoalValueSumPerUserProb.Percentile975, 28.0)
	assert.Greater(t, vrs["vid1"].GoalValueSumPerUserProbBest.Mean, 0.4)
	assert.Less(t, vrs["vid1"].GoalValueSumPerUserProbBest.Mean, 0.5)
	assert.Equal(t, vrs["vid1"].GoalValueSumPerUserProbBeatBaseline.Mean, 0.0)

	assert.Greater(t, vrs["vid2"].GoalValueSumPerUserProb.Median, 15.0)
	assert.Less(t, vrs["vid2"].GoalValueSumPerUserProb.Median, 17.0)
	assert.Greater(t, vrs["vid2"].GoalValueSumPerUserProb.Percentile025, -6.0)
	assert.Less(t, vrs["vid2"].GoalValueSumPerUserProb.Percentile025, -4.0)
	assert.Greater(t, vrs["vid2"].GoalValueSumPerUserProb.Percentile975, 36.0)
	assert.Less(t, vrs["vid2"].GoalValueSumPerUserProb.Percentile975, 37.5)
	assert.Greater(t, vrs["vid2"].GoalValueSumPerUserProbBest.Mean, 0.4)
	assert.Less(t, vrs["vid2"].GoalValueSumPerUserProbBest.Mean, 0.6)
	assert.Greater(t, vrs["vid2"].GoalValueSumPerUserProbBeatBaseline.Mean, 0.4)
	assert.Less(t, vrs["vid2"].GoalValueSumPerUserProbBeatBaseline.Mean, 0.6)
}
