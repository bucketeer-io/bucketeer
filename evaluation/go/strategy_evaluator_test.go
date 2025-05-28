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

package evaluation

import (
	"testing"

	"github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestStrategyEvaluator_Evaluate_Fixed(t *testing.T) {
	evaluator := &strategyEvaluator{}
	variations := []*feature.Variation{
		{Id: "variation-a", Value: "a"},
		{Id: "variation-b", Value: "b"},
	}
	strategy := &feature.Strategy{
		Type: feature.Strategy_FIXED,
		FixedStrategy: &feature.FixedStrategy{
			Variation: "variation-a",
		},
	}

	result, err := evaluator.Evaluate(strategy, "user-1", variations, "feature-1", "seed")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.Id != "variation-a" {
		t.Fatalf("Expected variation-a, got %s", result.Id)
	}
}

func TestStrategyEvaluator_Evaluate_Rollout_NoAudience(t *testing.T) {
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
			// No audience configuration (nil means no traffic control)
			Audience: nil,
		},
	}

	// Test with a user that would normally get variation-a
	result, err := evaluator.Evaluate(strategy, "user-1", variations, "feature-1", "seed")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// This should follow normal A/B split logic
	if result.Id != "variation-a" && result.Id != "variation-b" {
		t.Fatalf("Expected variation-a or variation-b, got %s", result.Id)
	}
}

func TestStrategyEvaluator_Evaluate_Rollout_WithAudience(t *testing.T) {
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

	for i := 0; i < totalUsers; i++ {
		userID := "user-" + string(rune(i))
		result, err := evaluator.Evaluate(strategy, userID, variations, "feature-1", "seed")
		if err != nil {
			t.Fatalf("Expected no error for user %s, got %v", userID, err)
		}

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

	if inExperimentCount < expectedInExperiment-tolerance || inExperimentCount > expectedInExperiment+tolerance {
		t.Fatalf("Expected approximately %d users in experiment, got %d (out of %d total)",
			expectedInExperiment, inExperimentCount, totalUsers)
	}

	if outOfExperimentCount < totalUsers-expectedInExperiment-tolerance || outOfExperimentCount > totalUsers-expectedInExperiment+tolerance {
		t.Fatalf("Expected approximately %d users out of experiment, got %d (out of %d total)",
			totalUsers-expectedInExperiment, outOfExperimentCount, totalUsers)
	}
}

func TestStrategyEvaluator_Evaluate_Rollout_Audience_NoDefaultVariation(t *testing.T) {
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
	for i := 0; i < 100; i++ {
		userID := "user-" + string(rune(i))
		_, err := evaluator.Evaluate(strategy, userID, variations, "feature-1", "seed")

		// We expect some users to get ErrVariationNotFound when they're outside traffic
		// and no default variation is specified
		if err == ErrVariationNotFound {
			// This is expected behavior
			return
		} else if err != nil {
			t.Fatalf("Unexpected error for user %s: %v", userID, err)
		}
	}
}

func TestStrategyEvaluator_Evaluate_Rollout_FullAudience(t *testing.T) {
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
			// 100% audience should behave like no audience control
			Audience: &feature.Audience{
				Percentage:       100,
				DefaultVariation: "variation-default",
			},
		},
	}

	// With 100% audience, all users should be in experiment
	for i := 0; i < 10; i++ {
		userID := "user-" + string(rune(i))
		result, err := evaluator.Evaluate(strategy, userID, variations, "feature-1", "seed")
		if err != nil {
			t.Fatalf("Expected no error for user %s, got %v", userID, err)
		}

		// Should never get default variation with 100% audience
		if result.Id == "variation-default" {
			t.Fatalf("Unexpected default variation for user %s with 100%% audience", userID)
		}
		if result.Id != "variation-a" && result.Id != "variation-b" {
			t.Fatalf("Expected variation-a or variation-b for user %s, got %s", userID, result.Id)
		}
	}
}
