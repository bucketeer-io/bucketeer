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

package scheduled

import (
	"context"
	"fmt"
	"time"

	featuredomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// ConflictDetector detects conflicts between scheduled flag changes.
//
// Conflict types detected:
//   - DEPENDENCY_MISSING: An earlier schedule deletes a variation/rule that a later schedule
//     depends on (would cause the later schedule to fail at execution time).
//   - INVALID_REFERENCE: The schedule references a variation/rule that doesn't exist in the
//     current flag state.
//
// Not treated as conflicts:
//   - Multiple schedules modifying the same field at different times. This is a valid use case
//     (e.g., enable at 10am, disable at 2pm; gradual rollout 25% -> 50% -> 100%).
//   - Flag version mismatch alone. A version bump (e.g., description change) doesn't
//     necessarily affect a schedule's validity. Only actual stale references matter.
type ConflictDetector struct {
	storage v2fs.ScheduledFlagChangeStorage
}

// NewConflictDetector creates a new ConflictDetector.
func NewConflictDetector(storage v2fs.ScheduledFlagChangeStorage) *ConflictDetector {
	return &ConflictDetector{storage: storage}
}

// DetectConflictsOnCreate checks for conflicts when creating or updating a schedule.
// It detects:
//   - DEPENDENCY_MISSING: an earlier schedule deletes a variation/rule this schedule references
//   - INVALID_REFERENCE: the schedule references variations/rules that don't exist in the flag
func (d *ConflictDetector) DetectConflictsOnCreate(
	ctx context.Context,
	flag *proto.Feature,
	newPayload *proto.ScheduledChangePayload,
	scheduledAt int64,
	environmentID string,
	excludeScheduleID string, // When updating, exclude the current schedule
) ([]*proto.ScheduledChangeConflict, error) {
	var conflicts []*proto.ScheduledChangeConflict

	pending, err := d.listPendingSchedules(ctx, flag.Id, environmentID)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	for _, existing := range pending {
		if existing.Id == excludeScheduleID {
			continue
		}

		// Check dependency: earlier schedule deletes variation/rule that new schedule references.
		// This would cause the new schedule to fail at execution time.
		if existing.ScheduledAt < scheduledAt {
			if refs := findDeletedReferencesNeededByPayload(existing.Payload, newPayload); len(refs) > 0 {
				for _, ref := range refs {
					conflicts = append(conflicts, &proto.ScheduledChangeConflict{
						Type: proto.ScheduledChangeConflict_CONFLICT_TYPE_DEPENDENCY_MISSING,
						Description: fmt.Sprintf(
							"Earlier schedule (ID: %s) deletes %s which this schedule references",
							existing.Id, ref,
						),
						ConflictingScheduleId: existing.Id,
						ConflictingField:      ref,
						DetectedAt:            now,
					})
				}
			}
		}
	}

	// Validate references exist in current flag state (INVALID_REFERENCE)
	conflicts = append(conflicts, validatePayloadReferences(flag, newPayload, now)...)

	return conflicts, nil
}

// DetectConflictsOnFlagChange is called when the flag is directly modified.
// It checks each pending schedule to see if the flag change caused any of its
// references to become stale (e.g., a variation it references was deleted).
// Only marks a schedule as CONFLICT if there are actual invalid references,
// not merely because the flag version changed.
// Returns the number of schedules marked as CONFLICT.
func (d *ConflictDetector) DetectConflictsOnFlagChange(
	ctx context.Context,
	flag *proto.Feature,
	environmentID string,
) (int, error) {
	pending, err := d.listPendingSchedules(ctx, flag.Id, environmentID)
	if err != nil {
		return 0, err
	}

	markedCount := 0
	now := time.Now().Unix()
	for _, schedule := range pending {
		// Only check schedules created before this flag version -- they may have stale refs
		if schedule.FlagVersionAtCreation >= flag.Version {
			continue
		}

		// Check if any references in the schedule's payload are now invalid
		invalidRefs := validatePayloadReferences(flag, schedule.Payload, now)
		if len(invalidRefs) == 0 {
			continue // No stale references, this schedule is still valid
		}

		sfcDomain := &featuredomain.ScheduledFlagChange{ScheduledFlagChange: schedule}
		sfcDomain.MarkConflict(invalidRefs)
		if err := d.storage.UpdateScheduledFlagChange(ctx, sfcDomain); err != nil {
			return markedCount, err
		}
		markedCount++
	}
	return markedCount, nil
}

func (d *ConflictDetector) listPendingSchedules(
	ctx context.Context,
	featureID, environmentID string,
) ([]*proto.ScheduledFlagChange, error) {
	filters := []*mysql.FilterV2{
		{Column: "environment_id", Operator: mysql.OperatorEqual, Value: environmentID},
		{Column: "feature_id", Operator: mysql.OperatorEqual, Value: featureID},
	}
	inFilters := []*mysql.InFilter{
		{
			Column: "status",
			Values: []interface{}{
				int32(proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING),
				int32(proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT),
			},
		},
	}
	options := &mysql.ListOptions{
		Filters:   filters,
		InFilters: inFilters,
		Orders:    []*mysql.Order{mysql.NewOrder("scheduled_at", mysql.OrderDirectionAsc)},
		Limit:     mysql.QueryNoLimit,
		Offset:    mysql.QueryNoOffset,
	}
	sfcs, _, _, err := d.storage.ListScheduledFlagChanges(ctx, options)
	return sfcs, err
}

// findDeletedReferencesNeededByPayload checks if an earlier schedule's payload
// deletes variations or rules that the new schedule's payload depends on.
// Returns a list of human-readable references that would be missing.
func findDeletedReferencesNeededByPayload(earlierPayload, newPayload *proto.ScheduledChangePayload) []string {
	if earlierPayload == nil || newPayload == nil {
		return nil
	}
	var refs []string

	// Variations deleted by earlier schedule
	deletedVars := make(map[string]bool)
	for _, vc := range earlierPayload.VariationChanges {
		if vc != nil && vc.ChangeType == proto.ChangeType_DELETE && vc.Variation != nil {
			deletedVars[vc.Variation.Id] = true
		}
	}

	// Check if new payload references any of those deleted variations
	for _, vc := range newPayload.VariationChanges {
		if vc == nil || vc.Variation == nil {
			continue
		}
		if vc.ChangeType == proto.ChangeType_UPDATE || vc.ChangeType == proto.ChangeType_DELETE {
			if deletedVars[vc.Variation.Id] {
				refs = append(refs, "variation "+vc.Variation.Id)
			}
		}
	}
	if newPayload.OffVariation != nil && newPayload.OffVariation.Value != "" {
		if deletedVars[newPayload.OffVariation.Value] {
			refs = append(refs, "off_variation "+newPayload.OffVariation.Value)
		}
	}
	if newPayload.DefaultStrategy != nil {
		if newPayload.DefaultStrategy.FixedStrategy != nil {
			if deletedVars[newPayload.DefaultStrategy.FixedStrategy.Variation] {
				refs = append(refs, "default_strategy.variation")
			}
		}
		if newPayload.DefaultStrategy.RolloutStrategy != nil {
			for _, rv := range newPayload.DefaultStrategy.RolloutStrategy.Variations {
				if deletedVars[rv.Variation] {
					refs = append(refs, "default_strategy.rollout.variation "+rv.Variation)
				}
			}
		}
	}
	for _, tc := range newPayload.TargetChanges {
		if tc != nil && tc.Target != nil && deletedVars[tc.Target.Variation] {
			refs = append(refs, "target.variation "+tc.Target.Variation)
		}
	}
	for _, rc := range newPayload.RuleChanges {
		if rc != nil && rc.Rule != nil && rc.Rule.Strategy != nil {
			if rc.Rule.Strategy.FixedStrategy != nil && deletedVars[rc.Rule.Strategy.FixedStrategy.Variation] {
				refs = append(refs, "rule.variation")
			}
			if rc.Rule.Strategy.RolloutStrategy != nil {
				for _, rv := range rc.Rule.Strategy.RolloutStrategy.Variations {
					if deletedVars[rv.Variation] {
						refs = append(refs, "rule.rollout.variation "+rv.Variation)
					}
				}
			}
		}
	}

	// Rules deleted by earlier schedule
	deletedRules := make(map[string]bool)
	for _, rc := range earlierPayload.RuleChanges {
		if rc != nil && rc.ChangeType == proto.ChangeType_DELETE && rc.Rule != nil {
			deletedRules[rc.Rule.Id] = true
		}
	}
	for _, rc := range newPayload.RuleChanges {
		if rc == nil || rc.Rule == nil {
			continue
		}
		if rc.ChangeType == proto.ChangeType_UPDATE || rc.ChangeType == proto.ChangeType_DELETE {
			if deletedRules[rc.Rule.Id] {
				refs = append(refs, "rule "+rc.Rule.Id)
			}
		}
	}

	return refs
}

// validatePayloadReferences checks if the schedule's payload references
// variations/rules that exist in the current flag state.
// Returns INVALID_REFERENCE conflicts for any stale references.
func validatePayloadReferences(
	flag *proto.Feature,
	payload *proto.ScheduledChangePayload,
	now int64,
) []*proto.ScheduledChangeConflict {
	if payload == nil || flag == nil {
		return nil
	}
	var conflicts []*proto.ScheduledChangeConflict

	// Variation references (update/delete require the variation to exist)
	for _, vc := range payload.VariationChanges {
		if vc == nil || vc.Variation == nil {
			continue
		}
		if vc.ChangeType == proto.ChangeType_UPDATE || vc.ChangeType == proto.ChangeType_DELETE {
			if !featuredomain.VariationExists(flag, vc.Variation.Id) {
				conflicts = append(conflicts, &proto.ScheduledChangeConflict{
					Type:             proto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE,
					Description:      fmt.Sprintf("Variation %s does not exist", vc.Variation.Id),
					ConflictingField: "variations",
					DetectedAt:       now,
				})
			}
		}
	}

	// Rule references (update/delete require the rule to exist)
	// Also check that any variation references inside the rule's strategy are valid
	for _, rc := range payload.RuleChanges {
		if rc == nil || rc.Rule == nil {
			continue
		}
		if rc.ChangeType == proto.ChangeType_UPDATE || rc.ChangeType == proto.ChangeType_DELETE {
			if !featuredomain.RuleExists(flag, rc.Rule.Id) {
				conflicts = append(conflicts, &proto.ScheduledChangeConflict{
					Type:             proto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE,
					Description:      fmt.Sprintf("Rule %s does not exist", rc.Rule.Id),
					ConflictingField: "rules",
					DetectedAt:       now,
				})
			}
		}
		if rc.Rule.Strategy != nil {
			if rc.Rule.Strategy.FixedStrategy != nil {
				vid := rc.Rule.Strategy.FixedStrategy.Variation
				if vid != "" && !featuredomain.VariationExists(flag, vid) {
					conflicts = append(conflicts, &proto.ScheduledChangeConflict{
						Type:             proto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE,
						Description:      fmt.Sprintf("Rule strategy references non-existent variation %s", vid),
						ConflictingField: "rules.strategy",
						DetectedAt:       now,
					})
				}
			}
			if rc.Rule.Strategy.RolloutStrategy != nil {
				for _, rv := range rc.Rule.Strategy.RolloutStrategy.Variations {
					if !featuredomain.VariationExists(flag, rv.Variation) {
						conflicts = append(conflicts, &proto.ScheduledChangeConflict{
							Type:             proto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE,
							Description:      fmt.Sprintf("Rule strategy references non-existent variation %s", rv.Variation),
							ConflictingField: "rules.strategy",
							DetectedAt:       now,
						})
						break
					}
				}
			}
		}
	}

	// Target variation references (the variation a target points to must exist)
	for _, tc := range payload.TargetChanges {
		if tc == nil || tc.Target == nil {
			continue
		}
		vid := tc.Target.Variation
		if vid != "" && !featuredomain.VariationExists(flag, vid) {
			conflicts = append(conflicts, &proto.ScheduledChangeConflict{
				Type:             proto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE,
				Description:      fmt.Sprintf("Target references non-existent variation %s", vid),
				ConflictingField: "targets",
				DetectedAt:       now,
			})
		}
	}

	// Off variation reference
	if payload.OffVariation != nil && payload.OffVariation.Value != "" {
		if !featuredomain.VariationExists(flag, payload.OffVariation.Value) {
			conflicts = append(conflicts, &proto.ScheduledChangeConflict{
				Type:             proto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE,
				Description:      fmt.Sprintf("OFF variation %s does not exist", payload.OffVariation.Value),
				ConflictingField: "off_variation",
				DetectedAt:       now,
			})
		}
	}

	// Default strategy variation references
	if payload.DefaultStrategy != nil {
		if payload.DefaultStrategy.FixedStrategy != nil {
			if !featuredomain.VariationExists(flag, payload.DefaultStrategy.FixedStrategy.Variation) {
				conflicts = append(conflicts, &proto.ScheduledChangeConflict{
					Type:             proto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE,
					Description:      "Default strategy references non-existent variation",
					ConflictingField: "default_strategy",
					DetectedAt:       now,
				})
			}
		}
		if payload.DefaultStrategy.RolloutStrategy != nil {
			for _, rv := range payload.DefaultStrategy.RolloutStrategy.Variations {
				if !featuredomain.VariationExists(flag, rv.Variation) {
					conflicts = append(conflicts, &proto.ScheduledChangeConflict{
						Type:             proto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE,
						Description:      "Default strategy references non-existent variation",
						ConflictingField: "default_strategy",
						DetectedAt:       now,
					})
					break
				}
			}
		}
	}

	return conflicts
}
