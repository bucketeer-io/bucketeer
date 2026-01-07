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

package evaluation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestStrategyEvaluator_Evaluate_Fixed(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		strategy    *feature.Strategy
		userID      string
		variations  []*feature.Variation
		featureID   string
		seed        string
		expected    string
		expectedErr error
	}{
		{
			desc: "success: fixed strategy",
			strategy: &feature.Strategy{
				Type: feature.Strategy_FIXED,
				FixedStrategy: &feature.FixedStrategy{
					Variation: "variation-a",
				},
			},
			userID: "user-1",
			variations: []*feature.Variation{
				{Id: "variation-a", Value: "a"},
				{Id: "variation-b", Value: "b"},
			},
			featureID:   "feature-1",
			seed:        "seed",
			expected:    "variation-a",
			expectedErr: nil,
		},
	}

	evaluator := &strategyEvaluator{}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			result, err := evaluator.Evaluate(p.strategy, p.userID, p.variations, p.featureID, p.seed)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.Equal(t, p.expected, result.Id)
			}
		})
	}
}

func TestStrategyEvaluator_Evaluate_Rollout_NoAudience(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		strategy    *feature.Strategy
		userID      string
		variations  []*feature.Variation
		featureID   string
		seed        string
		expectedIDs []string
		expectedErr error
	}{
		{
			desc: "success: rollout strategy without audience",
			strategy: &feature.Strategy{
				Type: feature.Strategy_ROLLOUT,
				RolloutStrategy: &feature.RolloutStrategy{
					Variations: []*feature.RolloutStrategy_Variation{
						{Variation: "variation-a", Weight: 50000},
						{Variation: "variation-b", Weight: 50000},
					},
					// No audience configuration (nil means no traffic control)
					Audience: nil,
				},
			},
			userID: "user-1",
			variations: []*feature.Variation{
				{Id: "variation-a", Value: "a"},
				{Id: "variation-b", Value: "b"},
			},
			featureID:   "feature-1",
			seed:        "seed",
			expectedIDs: []string{"variation-a", "variation-b"},
			expectedErr: nil,
		},
	}

	evaluator := &strategyEvaluator{}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			result, err := evaluator.Evaluate(p.strategy, p.userID, p.variations, p.featureID, p.seed)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.Contains(t, p.expectedIDs, result.Id)
			}
		})
	}
}

func TestStrategyEvaluator_Evaluate_Rollout_WithAudience(t *testing.T) {
	t.Parallel()
	evaluator := &strategyEvaluator{}
	variations := []*feature.Variation{
		{Id: "variation-a", Value: "a"},
		{Id: "variation-b", Value: "b"},
		{Id: "variation-default", Value: "default"},
	}
	strategy := &feature.Strategy{
		Type: feature.Strategy_ROLLOUT,
		RolloutStrategy: &feature.RolloutStrategy{
			Variations: []*feature.RolloutStrategy_Variation{
				{Variation: "variation-a", Weight: 50000},
				{Variation: "variation-b", Weight: 50000},
			},
			// 10% audience configuration
			Audience: &feature.Audience{
				Percentage:       10,
				DefaultVariation: "variation-default",
			},
		},
	}

	// Test multiple users to verify audience control
	inExperimentCount := 0
	outOfExperimentCount := 0
	totalUsers := 1000

	for i := range totalUsers {
		userID := fmt.Sprintf("user-%d", i)
		result, err := evaluator.Evaluate(strategy, userID, variations, "feature-1", "seed")
		assert.NoError(t, err)

		if result.Id == "variation-default" {
			outOfExperimentCount++
		} else if result.Id == "variation-a" || result.Id == "variation-b" {
			inExperimentCount++
		} else {
			t.Fatalf("Unexpected variation %s for user %s", result.Id, userID)
		}
	}

	// Verify approximately 10% are in experiment (allow some variance)
	expectedInExperiment := totalUsers * 10 / 100
	tolerance := totalUsers * 5 / 100 // 5% tolerance

	assert.True(t, inExperimentCount >= expectedInExperiment-tolerance && inExperimentCount <= expectedInExperiment+tolerance,
		"Expected approximately %d users in experiment, got %d (out of %d total)",
		expectedInExperiment, inExperimentCount, totalUsers)

	assert.True(t, outOfExperimentCount >= totalUsers-expectedInExperiment-tolerance && outOfExperimentCount <= totalUsers-expectedInExperiment+tolerance,
		"Expected approximately %d users out of experiment, got %d (out of %d total)",
		totalUsers-expectedInExperiment, outOfExperimentCount, totalUsers)
}

func TestStrategyEvaluator_Evaluate_Rollout_Audience_NoDefaultVariation(t *testing.T) {
	t.Parallel()
	evaluator := &strategyEvaluator{}
	variations := []*feature.Variation{
		{Id: "variation-a", Value: "a"},
		{Id: "variation-b", Value: "b"},
	}
	strategy := &feature.Strategy{
		Type: feature.Strategy_ROLLOUT,
		RolloutStrategy: &feature.RolloutStrategy{
			Variations: []*feature.RolloutStrategy_Variation{
				{Variation: "variation-a", Weight: 50000},
				{Variation: "variation-b", Weight: 50000},
			},
			Audience: &feature.Audience{
				Percentage:       10,
				DefaultVariation: "", // No default variation specified
			},
		},
	}

	// Find a user that would be outside the experiment traffic
	foundExpectedError := false
	for i := range 100 {
		userID := fmt.Sprintf("user-%d", i)
		_, err := evaluator.Evaluate(strategy, userID, variations, "feature-1", "seed")

		// We expect some users to get ErrVariationNotFound when they're outside traffic
		// and no default variation is specified
		if err == ErrVariationNotFound {
			// This is expected behavior
			foundExpectedError = true
			break
		} else if err != nil {
			t.Fatalf("Unexpected error for user %s: %v", userID, err)
		}
	}

	assert.True(t, foundExpectedError, "Expected at least one user to get ErrVariationNotFound, but none did")
}

func TestStrategyEvaluator_Evaluate_Rollout_FullAudience(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		strategy     *feature.Strategy
		variations   []*feature.Variation
		userID       string
		featureID    string
		seed         string
		expectedIDs  []string
		unexpectedID string
		expectedErr  error
	}{
		{
			desc: "success: 100% audience should behave like no audience control",
			strategy: &feature.Strategy{
				Type: feature.Strategy_ROLLOUT,
				RolloutStrategy: &feature.RolloutStrategy{
					Variations: []*feature.RolloutStrategy_Variation{
						{Variation: "variation-a", Weight: 50000},
						{Variation: "variation-b", Weight: 50000},
					},
					// 100% audience should behave like no audience control
					Audience: &feature.Audience{
						Percentage:       100,
						DefaultVariation: "variation-default",
					},
				},
			},
			variations: []*feature.Variation{
				{Id: "variation-a", Value: "a"},
				{Id: "variation-b", Value: "b"},
				{Id: "variation-default", Value: "default"},
			},
			userID:       "user-1",
			featureID:    "feature-1",
			seed:         "seed",
			expectedIDs:  []string{"variation-a", "variation-b"},
			unexpectedID: "variation-default",
			expectedErr:  nil,
		},
	}

	evaluator := &strategyEvaluator{}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			// With 100% audience, all users should be in experiment
			for i := range 10 {
				userID := fmt.Sprintf("user-%d", i)
				result, err := evaluator.Evaluate(p.strategy, userID, p.variations, p.featureID, p.seed)
				assert.Equal(t, p.expectedErr, err)

				if err == nil {
					// Should never get default variation with 100% audience
					assert.NotEqual(t, p.unexpectedID, result.Id,
						"Unexpected default variation for user %s with 100%% audience", userID)
					assert.Contains(t, p.expectedIDs, result.Id,
						"Expected variation-a or variation-b for user %s, got %s", userID, result.Id)
				}
			}
		})
	}
}
