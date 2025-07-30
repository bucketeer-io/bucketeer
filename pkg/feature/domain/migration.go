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

package domain

import (
	"time"

	"github.com/bucketeer-io/bucketeer/proto/feature"
)

// VariationCleanupResult contains details about what was cleaned up
type VariationCleanupResult struct {
	Changed              bool
	OrphanedTargets      int
	OrphanedRules        int
	OrphanedDefault      int
	OrphanedOffVar       bool
	OrphanedVariationIDs []string
}

// CleanupOrphanedVariationReferences removes references to variations that no longer exist.
// This fixes data corruption caused by the incomplete variation deletion bug.
// TODO: Remove this function after 6 months (around July 2025) when all corrupted data is cleaned up
func (f *Feature) CleanupOrphanedVariationReferences() VariationCleanupResult {
	result := VariationCleanupResult{
		OrphanedVariationIDs: []string{},
	}

	if f == nil || f.Feature == nil {
		return result
	}

	validVariationIDs := make(map[string]bool)
	orphanedVariationIDs := make(map[string]bool)

	// Build map of valid variation IDs
	for _, v := range f.Variations {
		validVariationIDs[v.Id] = true
	}

	// 1. Clean up orphaned targets
	validTargets := make([]*feature.Target, 0, len(f.Targets))
	for _, target := range f.Targets {
		if validVariationIDs[target.Variation] {
			validTargets = append(validTargets, target)
		} else {
			result.Changed = true
			result.OrphanedTargets++
			orphanedVariationIDs[target.Variation] = true
		}
	}
	f.Targets = validTargets

	// 2. Clean up orphaned variations in rules
	for _, rule := range f.Rules {
		if rule.Strategy != nil && rule.Strategy.Type == feature.Strategy_ROLLOUT && rule.Strategy.RolloutStrategy != nil {
			validRolloutVariations := make(
				[]*feature.RolloutStrategy_Variation,
				0,
				len(rule.Strategy.RolloutStrategy.Variations),
			)
			for _, v := range rule.Strategy.RolloutStrategy.Variations {
				if validVariationIDs[v.Variation] {
					validRolloutVariations = append(validRolloutVariations, v)
				} else {
					result.Changed = true
					result.OrphanedRules++
					orphanedVariationIDs[v.Variation] = true
				}
			}
			rule.Strategy.RolloutStrategy.Variations = validRolloutVariations
		}
	}

	// 3. Clean up orphaned variations in default strategy
	if f.DefaultStrategy != nil &&
		f.DefaultStrategy.Type == feature.Strategy_ROLLOUT &&
		f.DefaultStrategy.RolloutStrategy != nil {
		validDefaultVariations := make(
			[]*feature.RolloutStrategy_Variation,
			0,
			len(f.DefaultStrategy.RolloutStrategy.Variations),
		)
		for _, v := range f.DefaultStrategy.RolloutStrategy.Variations {
			if validVariationIDs[v.Variation] {
				validDefaultVariations = append(validDefaultVariations, v)
			} else {
				result.Changed = true
				result.OrphanedDefault++
				orphanedVariationIDs[v.Variation] = true
			}
		}
		f.DefaultStrategy.RolloutStrategy.Variations = validDefaultVariations
	}

	// 4. Check if OffVariation still exists
	if f.OffVariation != "" && !validVariationIDs[f.OffVariation] {
		result.OrphanedOffVar = true
		orphanedVariationIDs[f.OffVariation] = true

		// Reset to second available variation, fallback to first if only one exists
		if len(f.Variations) > 1 {
			f.OffVariation = f.Variations[1].Id
			result.Changed = true
		} else if len(f.Variations) > 0 {
			f.OffVariation = f.Variations[0].Id
			result.Changed = true
		}
	}

	// Convert orphaned variation IDs to slice
	for variationID := range orphanedVariationIDs {
		result.OrphanedVariationIDs = append(result.OrphanedVariationIDs, variationID)
	}

	// Update timestamp if any changes were made
	if result.Changed {
		f.UpdatedAt = time.Now().Unix()
	}

	return result
}

// CleanupOrphanedVariationReferencesSimple provides backward compatibility
// TODO: Remove this after updating all call sites to use detailed version
func (f *Feature) CleanupOrphanedVariationReferencesSimple() bool {
	result := f.CleanupOrphanedVariationReferences()
	return result.Changed
}

// ValidateVariationReferences checks if a feature has orphaned variation references.
// Returns a list of orphaned variation IDs found.
func (f *Feature) ValidateVariationReferences() []string {
	if f == nil || f.Feature == nil {
		return nil
	}

	validVariationIDs := make(map[string]bool)
	orphanedVariations := make(map[string]bool)

	// Build map of valid variation IDs
	for _, v := range f.Variations {
		validVariationIDs[v.Id] = true
	}

	// Check targets for orphaned references
	for _, target := range f.Targets {
		if !validVariationIDs[target.Variation] {
			orphanedVariations[target.Variation] = true
		}
	}

	// Check rules for orphaned references
	for _, rule := range f.Rules {
		if rule.Strategy != nil && rule.Strategy.Type == feature.Strategy_ROLLOUT && rule.Strategy.RolloutStrategy != nil {
			for _, v := range rule.Strategy.RolloutStrategy.Variations {
				if !validVariationIDs[v.Variation] {
					orphanedVariations[v.Variation] = true
				}
			}
		}
	}

	// Check default strategy for orphaned references
	if f.DefaultStrategy != nil &&
		f.DefaultStrategy.Type == feature.Strategy_ROLLOUT &&
		f.DefaultStrategy.RolloutStrategy != nil {
		for _, v := range f.DefaultStrategy.RolloutStrategy.Variations {
			if !validVariationIDs[v.Variation] {
				orphanedVariations[v.Variation] = true
			}
		}
	}

	// Check OffVariation
	if f.OffVariation != "" && !validVariationIDs[f.OffVariation] {
		orphanedVariations[f.OffVariation] = true
	}

	// Convert map to slice
	result := make([]string, 0, len(orphanedVariations))
	for variationID := range orphanedVariations {
		result = append(result, variationID)
	}

	return result
}

type VariationMigrationResult struct {
	Changed           bool
	AddedToRules      int
	AddedToDefault    int
	AddedVariationIDs []string
}

// EnsureVariationsInStrategies adds missing variations to rules and default strategy
// This fixes data corruption from the historical AddVariation bug where variations
// were added to the variations list and targets, but not to rollout strategies
func (f *Feature) EnsureVariationsInStrategies() VariationMigrationResult {
	result := VariationMigrationResult{
		Changed:           false,
		AddedToRules:      0,
		AddedToDefault:    0,
		AddedVariationIDs: []string{},
	}

	if len(f.Variations) == 0 {
		return result
	}

	// Create a set of all valid variation IDs
	validVariationIDs := make(map[string]bool)
	for _, variation := range f.Variations {
		validVariationIDs[variation.Id] = true
	}

	// 1. Add missing variations to rules with rollout strategies
	for _, rule := range f.Rules {
		if rule.Strategy != nil &&
			rule.Strategy.Type == feature.Strategy_ROLLOUT &&
			rule.Strategy.RolloutStrategy != nil {
			// TODO: Remove this after updating all call sites to use detailed version
			added := f.ensureVariationsInRolloutStrategy(rule.Strategy.RolloutStrategy, validVariationIDs)
			result.AddedToRules += added
			if added > 0 {
				result.Changed = true
			}
		}
	}

	// 2. Add missing variations to default strategy
	if f.DefaultStrategy != nil &&
		f.DefaultStrategy.Type == feature.Strategy_ROLLOUT &&
		f.DefaultStrategy.RolloutStrategy != nil {
		added := f.ensureVariationsInRolloutStrategy(f.DefaultStrategy.RolloutStrategy, validVariationIDs)
		result.AddedToDefault += added
		if added > 0 {
			result.Changed = true
		}
	}

	// 3. Collect variation IDs that were processed (for logging)
	if result.Changed {
		for variationID := range validVariationIDs {
			result.AddedVariationIDs = append(result.AddedVariationIDs, variationID)
		}
	}

	return result
}

// ensureVariationsInRolloutStrategy adds missing variations to a rollout strategy
// Returns the number of variations added
func (f *Feature) ensureVariationsInRolloutStrategy(
	strategy *feature.RolloutStrategy,
	validVariationIDs map[string]bool,
) int {
	if strategy == nil {
		return 0
	}

	// Create a map of existing weights to preserve them
	existingWeights := make(map[string]int32)
	for _, strategyVar := range strategy.Variations {
		existingWeights[strategyVar.Variation] = strategyVar.Weight
	}

	// Check if we need to add any missing variations
	missingVariations := []string{}
	for _, variation := range f.Variations {
		variationID := variation.Id
		if validVariationIDs[variationID] {
			if _, exists := existingWeights[variationID]; !exists {
				missingVariations = append(missingVariations, variationID)
			}
		}
	}

	// If no missing variations, no changes needed
	if len(missingVariations) == 0 {
		return 0
	}

	// Reconstruct strategy.Variations to match f.Variations order exactly
	// This ensures UI consistency and predictable ordering
	originalCount := len(strategy.Variations)
	newStrategyVariations := []*feature.RolloutStrategy_Variation{}

	for _, variation := range f.Variations {
		variationID := variation.Id
		if validVariationIDs[variationID] {
			weight := existingWeights[variationID] // Will be 0 for missing variations
			newStrategyVariations = append(newStrategyVariations, &feature.RolloutStrategy_Variation{
				Variation: variationID,
				Weight:    weight,
			})
		}
	}

	strategy.Variations = newStrategyVariations
	return len(strategy.Variations) - originalCount
}

// EnsureVariationsInStrategiesSimple is a convenience wrapper that returns only a boolean
func (f *Feature) EnsureVariationsInStrategiesSimple() bool {
	result := f.EnsureVariationsInStrategies()
	return result.Changed
}
