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
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	envclient "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client/mock"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/experimentcalculator/stan"
	metricsmock "github.com/bucketeer-io/bucketeer/pkg/metrics/mock"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/proto/eventcounter"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	jpLocation = time.FixedZone("Asia/Tokyo", 9*60*60)
)

const (
	stanModelID = "y3qsnd7m"
)

func creatExperimentCalculator(mockController *gomock.Controller) *ExperimentCalculator {
	registerer := metricsmock.NewMockRegisterer(mockController)
	registerer.EXPECT().MustRegister(gomock.Any()).Return().Times(2) // One for stan metrics, one for calculator metrics
	return NewExperimentCalculator(
		stan.NewStan("localhost", "8080", registerer, zap.NewNop()),
		stanModelID,
		envclient.NewMockClient(mockController),
		ecclient.NewMockClient(mockController),
		experimentclient.NewMockClient(mockController),
		mysqlmock.NewMockClient(mockController),
		registerer,
		jpLocation,
		zap.NewNop(),
	)
}

func TestExperimentCalculatorBinomialModelSample(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	experimentCalculator := creatExperimentCalculator(mockController)
	experiment := &experimentproto.Experiment{
		BaseVariationId: "vid1",
		Variations: []*featureproto.Variation{
			{
				Id: "vid1",
			},
			{
				Id: "vid2",
			},
		},
	}
	assert.NotNil(t, experimentCalculator, "ExperimentCalculator should not be nil")
	ctx := context.TODO()
	vrs, err := experimentCalculator.binomialModelSample(
		ctx,
		[]string{"vid1", "vid2"},
		[]int64{38, 51},
		[]int64{101, 99},
		0,
		experiment,
	)

	assert.NoError(t, err, "BinomialModelSample should not be error")

	assert.GreaterOrEqual(t, vrs["vid1"].CvrProb.Mean, 0.37)
	assert.LessOrEqual(t, vrs["vid1"].CvrProb.Mean, 0.38)
	assert.GreaterOrEqual(t, vrs["vid1"].CvrProb.Sd, 0.045)
	assert.LessOrEqual(t, vrs["vid1"].CvrProb.Sd, 0.05)
	assert.GreaterOrEqual(t, vrs["vid1"].CvrProb.Rhat, 0.9)
	assert.LessOrEqual(t, vrs["vid1"].CvrProb.Rhat, 1.1)
	assert.Equal(t, len(vrs["vid1"].CvrProb.Histogram.Hist), 100)
	assert.Equal(t, len(vrs["vid1"].CvrProb.Histogram.Bins), 101)
	assert.GreaterOrEqual(t, vrs["vid1"].CvrProbBest.Mean, 0.023)
	assert.LessOrEqual(t, vrs["vid1"].CvrProbBest.Mean, 0.026)
	assert.GreaterOrEqual(t, vrs["vid1"].CvrProbBest.Sd, 0.15)
	assert.LessOrEqual(t, vrs["vid1"].CvrProbBest.Sd, 0.16)
	assert.GreaterOrEqual(t, vrs["vid1"].CvrProbBest.Rhat, 0.9)
	assert.LessOrEqual(t, vrs["vid1"].CvrProbBest.Rhat, 1.1)
	assert.Equal(t, vrs["vid1"].CvrProbBeatBaseline.Mean, 0.0)
	assert.Equal(t, vrs["vid1"].CvrProbBeatBaseline.Sd, 0.0)
	assert.Equal(t, vrs["vid1"].CvrProbBeatBaseline.Rhat, 0.0)

	assert.GreaterOrEqual(t, vrs["vid2"].CvrProb.Mean, 0.49)
	assert.LessOrEqual(t, vrs["vid2"].CvrProb.Mean, 0.52)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProb.Sd, 0.045)
	assert.LessOrEqual(t, vrs["vid2"].CvrProb.Sd, 0.05)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProb.Rhat, 0.9)
	assert.LessOrEqual(t, vrs["vid2"].CvrProb.Rhat, 1.1)
	assert.Equal(t, len(vrs["vid2"].CvrProb.Histogram.Hist), 100)
	assert.Equal(t, len(vrs["vid2"].CvrProb.Histogram.Bins), 101)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProbBest.Mean, 0.97)
	assert.LessOrEqual(t, vrs["vid2"].CvrProbBest.Mean, 0.98)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProbBest.Sd, 0.15)
	assert.LessOrEqual(t, vrs["vid2"].CvrProbBest.Sd, 0.16)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProbBest.Rhat, 0.9)
	assert.LessOrEqual(t, vrs["vid2"].CvrProbBest.Rhat, 1.1)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProbBeatBaseline.Mean, 0.97)
	assert.LessOrEqual(t, vrs["vid2"].CvrProbBeatBaseline.Mean, 0.98)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProbBeatBaseline.Sd, 0.15)
	assert.LessOrEqual(t, vrs["vid2"].CvrProbBeatBaseline.Sd, 0.16)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProbBeatBaseline.Rhat, 0.9)
	assert.LessOrEqual(t, vrs["vid2"].CvrProbBeatBaseline.Rhat, 1.1)

}

func TestListEndAt(t *testing.T) {
	type args struct {
		startAt int64
		endAt   int64
		now     int64
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{
			name: "3 days 18 hours",
			args: args{
				startAt: 1614848400, // 2021-03-04 09:00:00Z
				endAt:   1615086000, // 2021-03-07 03:00:00Z
				now:     32508810000,
			},
			want: []int64{1614934800, 1615021200, 1615086000},
		},
		{
			name: "now is earlier than end_at",
			args: args{
				startAt: 1614848400, // 2021-03-04 09:00:00Z
				endAt:   1615086000, // 2021-03-07 03:00:00Z
				now:     1614967200, // 2021-03-06 03:00:00Z
			},
			want: []int64{1614934800, 1614967200},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, listEndAt(tt.args.startAt, tt.args.endAt, tt.args.now),
				"listEndAt(%v, %v, %v)", tt.args.startAt, tt.args.endAt, tt.args.now)
		})
	}
}

// TestCalculateExpectedLoss tests the calculateExpectedLoss function with various scenarios
func TestCalculateExpectedLoss(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name           string
		variations     []*eventcounter.VariationResult
		expectModified bool               // whether we expect the variations to be modified
		expected       map[string]float64 // expected expectedLoss values by variation ID
	}{
		{
			name: "basic expected loss calculation",
			variations: []*eventcounter.VariationResult{
				{
					VariationId: "var1",
					CvrSamples:  []float64{0.1, 0.2, 0.3},
				},
				{
					VariationId: "var2",
					CvrSamples:  []float64{0.2, 0.1, 0.4},
				},
			},
			expectModified: true,
			expected: map[string]float64{
				"var1": ((0.2 - 0.1) + (0.2 - 0.2) + (0.4 - 0.3)) * 100 / 3, // (0.1 + 0 + 0.1) * 100 / 3 = 6.67
				"var2": ((0.2 - 0.2) + (0.2 - 0.1) + (0.4 - 0.4)) * 100 / 3, // (0 + 0.1 + 0) * 100 / 3 = 3.33
			},
		},
		{
			name: "three variations with ties",
			variations: []*eventcounter.VariationResult{
				{
					VariationId: "var1",
					CvrSamples:  []float64{0.5, 0.5, 0.5},
				},
				{
					VariationId: "var2",
					CvrSamples:  []float64{0.5, 0.6, 0.4},
				},
				{
					VariationId: "var3",
					CvrSamples:  []float64{0.4, 0.4, 0.6},
				},
			},
			expectModified: true,
			expected: map[string]float64{
				"var1": ((0.5 - 0.5) + (0.6 - 0.5) + (0.6 - 0.5)) * 100 / 3, // (0 + 0.1 + 0.1) * 100 / 3 = 6.67
				"var2": ((0.5 - 0.5) + (0.6 - 0.6) + (0.6 - 0.4)) * 100 / 3, // (0 + 0 + 0.2) * 100 / 3 = 6.67
				"var3": ((0.5 - 0.4) + (0.6 - 0.4) + (0.6 - 0.6)) * 100 / 3, // (0.1 + 0.2 + 0) * 100 / 3 = 10.0
			},
		},
		{
			name: "single variation",
			variations: []*eventcounter.VariationResult{
				{VariationId: "solo", CvrSamples: []float64{0.1, 0.4, 0.3}},
			},
			expectModified: true, // we do run through the code and set ExpectedLoss to 0
			expected: map[string]float64{
				"solo": 0.0,
			},
		},
		{
			name:           "empty variations array",
			variations:     []*eventcounter.VariationResult{},
			expectModified: false,
			expected:       map[string]float64{},
		},
		{
			name: "inconsistent sample lengths",
			variations: []*eventcounter.VariationResult{
				{
					VariationId: "var1",
					CvrSamples:  []float64{0.1, 0.2, 0.3},
				},
				{
					VariationId: "var2",
					CvrSamples:  []float64{0.2, 0.1}, // Missing a sample
				},
			},
			expectModified: false,
			expected: map[string]float64{
				"var1": 0.0,
				"var2": 0.0,
			},
		},
		{
			name: "variation with no samples",
			variations: []*eventcounter.VariationResult{
				{
					VariationId: "var1",
					CvrSamples:  []float64{0.1, 0.2, 0.3},
				},
				{
					VariationId: "var2",
					CvrSamples:  []float64{}, // Empty samples
				},
			},
			expectModified: false,
			expected: map[string]float64{
				"var1": 0.0,
				"var2": 0.0,
			},
		},
	}

	// Helper function to capture initial expectedLoss values before calculation
	saveInitialValues := func(variations []*eventcounter.VariationResult) map[string]float64 {
		result := make(map[string]float64)
		for _, v := range variations {
			result[v.VariationId] = v.ExpectedLoss
		}
		return result
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create calculator with test logger
			calculator := ExperimentCalculator{
				logger: logger,
			}

			// Save initial values
			initialValues := saveInitialValues(tt.variations)

			// Run the calculation
			calculator.calculateExpectedLoss(tt.variations)

			// Check expected modifications
			for _, variation := range tt.variations {
				if tt.expectModified {
					// For cases where we expect modification, check against calculated values
					expected := tt.expected[variation.VariationId]
					assert.InDelta(t, expected, variation.ExpectedLoss, 0.01,
						"Expected loss for %s should be %f, got %f",
						variation.VariationId, expected, variation.ExpectedLoss)
				} else {
					// For cases where we don't expect modification, values should be unchanged
					assert.Equal(t, initialValues[variation.VariationId], variation.ExpectedLoss,
						"Expected loss should not change for %s", variation.VariationId)
				}
			}
		})
	}
}

// TestExpectedLossWithBinomialModelSamples tests the expected loss calculation with real CVR samples from binomialModelSample
func TestExpectedLossWithBinomialModelSamples(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	experimentCalculator := creatExperimentCalculator(mockController)
	experiment := &experimentproto.Experiment{
		BaseVariationId: "vid1",
		Variations: []*featureproto.Variation{
			{
				Id: "vid1",
			},
			{
				Id: "vid2",
			},
		},
	}
	ctx := context.TODO()

	// Get variation results with CVR samples from binomialModelSample
	vrs, err := experimentCalculator.binomialModelSample(
		ctx,
		[]string{"vid1", "vid2"},
		[]int64{38, 51},
		[]int64{101, 99},
		0,
		experiment,
	)
	assert.NoError(t, err, "BinomialModelSample should not return an error")

	// Convert map to slice for calculateExpectedLoss
	variationResults := []*eventcounter.VariationResult{
		vrs["vid1"],
		vrs["vid2"],
	}

	// Ensure we have CVR samples before calculation
	assert.True(t, len(variationResults[0].CvrSamples) > 0, "vid1 should have CVR samples")
	assert.True(t, len(variationResults[1].CvrSamples) > 0, "vid2 should have CVR samples")

	// Calculate expected loss
	experimentCalculator.calculateExpectedLoss(variationResults)

	// Verify expected loss calculation
	// Since vid2 has higher conversion rate (around 0.5) than vid1 (around 0.37),
	// the expected loss for vid1 should be higher and vid2 should be lower
	assert.Greater(t, variationResults[0].ExpectedLoss, variationResults[1].ExpectedLoss,
		"Expected loss for vid1 should be greater than vid2")

	// Expected loss for vid1 should be positive (around 10-13%)
	assert.Greater(t, variationResults[0].ExpectedLoss, 0.0,
		"Expected loss for vid1 should be positive")

	// Expected loss for vid2 (best variation) should be close to 0
	assert.InDelta(t, 0.0, variationResults[1].ExpectedLoss, 1.0,
		"Expected loss for vid2 should be close to 0")

	// Log the actual values for reference
	t.Logf("vid1 expected loss: %.2f%%", variationResults[0].ExpectedLoss)
	t.Logf("vid2 expected loss: %.2f%%", variationResults[1].ExpectedLoss)
}
