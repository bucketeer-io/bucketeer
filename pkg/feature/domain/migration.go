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

// CleanupOrphanedVariationReferences removes references to variations that no longer exist.
// This fixes data corruption caused by the incomplete variation deletion bug.
func (f *Feature) CleanupOrphanedVariationReferences() bool {
	if f == nil || f.Feature == nil {
		return false
	}

	changed := false
	validVariationIDs := make(map[string]bool)

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
			changed = true
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
					changed = true
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
				changed = true
			}
		}
		f.DefaultStrategy.RolloutStrategy.Variations = validDefaultVariations
	}

	// 4. Check if OffVariation still exists
	if f.OffVariation != "" && !validVariationIDs[f.OffVariation] {
		// Reset to second available variation, fallback to first if only one exists
		if len(f.Variations) > 1 {
			f.OffVariation = f.Variations[1].Id
			changed = true
		} else if len(f.Variations) > 0 {
			f.OffVariation = f.Variations[0].Id
			changed = true
		}
	}

	// Update timestamp if any changes were made
	if changed {
		f.UpdatedAt = time.Now().Unix()
	}

	return changed
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
