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

	"go.uber.org/zap"

	featuredomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	// maxCrossFlagSchedulesToCheck limits the number of cross-flag
	// schedules evaluated per UpdateFeature call to bound the DB
	// queries and processing time.
	maxCrossFlagSchedulesToCheck = 100
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
	storage        v2fs.ScheduledFlagChangeStorage
	featureStorage v2fs.FeatureStorage
	logger         *zap.Logger
}

// NewConflictDetector creates a new ConflictDetector.
func NewConflictDetector(
	storage v2fs.ScheduledFlagChangeStorage,
) *ConflictDetector {
	return &ConflictDetector{
		storage: storage,
		logger:  zap.NewNop(),
	}
}

// NewConflictDetectorWithFeatureStorage creates a ConflictDetector
// with feature storage for cross-flag conflict detection
// (e.g., prerequisite validation across flags).
func NewConflictDetectorWithFeatureStorage(
	storage v2fs.ScheduledFlagChangeStorage,
	featureStorage v2fs.FeatureStorage,
	logger *zap.Logger,
) *ConflictDetector {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &ConflictDetector{
		storage:        storage,
		featureStorage: featureStorage,
		logger:         logger,
	}
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

	pending, err := d.listPendingAndConflictSchedules(ctx, flag.Id, environmentID)
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
// It checks each pending/conflict schedule to see if the flag change caused any of its
// references to become stale (e.g., a variation it references was deleted).
// Only marks a schedule as CONFLICT if there are actual invalid references,
// not merely because the flag version changed.
//
// Auto-recovery: If a previously CONFLICT schedule now has all valid references
// (e.g., a deleted variation was re-added), it is restored to PENDING.
//
// Returns the number of schedules whose status changed (newly marked CONFLICT + recovered).
func (d *ConflictDetector) DetectConflictsOnFlagChange(
	ctx context.Context,
	flag *proto.Feature,
	environmentID string,
) (int, error) {
	pending, err := d.listPendingAndConflictSchedules(ctx, flag.Id, environmentID)
	if err != nil {
		return 0, err
	}

	changedCount := 0
	now := time.Now().Unix()
	for _, schedule := range pending {
		isPending := schedule.Status == proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING
		isConflict := schedule.Status == proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT

		// For PENDING schedules: only check if created before this flag version
		if isPending && schedule.FlagVersionAtCreation >= flag.Version {
			continue
		}

		// Check if any references in the schedule's payload are now invalid
		invalidRefs := validatePayloadReferences(flag, schedule.Payload, now)

		if len(invalidRefs) > 0 && isPending {
			// PENDING → CONFLICT: references are now invalid
			sfcDomain := &featuredomain.ScheduledFlagChange{ScheduledFlagChange: schedule}
			sfcDomain.MarkConflict(invalidRefs)
			if err := d.storage.UpdateScheduledFlagChange(ctx, sfcDomain); err != nil {
				return changedCount, err
			}
			changedCount++
		} else if len(invalidRefs) == 0 && isConflict {
			// Same-flag refs are valid, but the schedule may have been
			// marked CONFLICT due to a cross-flag prerequisite issue.
			// Re-validate prerequisites before restoring to PENDING.
			if hasPrerequisiteChanges(schedule.Payload) &&
				d.featureStorage != nil {
				prereqConflicts := d.validatePrerequisiteReferences(
					ctx, schedule.Payload, environmentID, now,
				)
				if len(prereqConflicts) > 0 {
					continue
				}
			}
			sfcDomain := &featuredomain.ScheduledFlagChange{ScheduledFlagChange: schedule}
			sfcDomain.RestoreToPending()
			if err := d.storage.UpdateScheduledFlagChange(ctx, sfcDomain); err != nil {
				return changedCount, err
			}
			changedCount++
		}
	}
	return changedCount, nil
}

// DetectCrossFlagConflicts is called after a flag is updated to check schedules for OTHER flags
// in the same environment. This catches cross-flag prerequisite conflicts:
// e.g., Flag A schedules a prerequisite referencing Flag B's variation, and that variation
// is deleted from Flag B.
//
// It also handles auto-recovery: if a previously CONFLICT schedule on another flag now has
// all valid references, it is restored to PENDING.
//
// This method requires featureStorage to be set (via NewConflictDetectorWithFeatureStorage).
// It is best-effort: errors are returned but the caller should not fail the UpdateFeature request.
func (d *ConflictDetector) DetectCrossFlagConflicts(
	ctx context.Context,
	updatedFlagID string,
	environmentID string,
) (int, error) {
	if d.featureStorage == nil {
		return 0, nil
	}

	// List all pending/conflict schedules in the environment for OTHER flags
	schedules, err := d.listCrossFlagSchedules(ctx, updatedFlagID, environmentID)
	if err != nil {
		return 0, fmt.Errorf("failed to list cross-flag schedules: %w", err)
	}

	if len(schedules) == 0 {
		return 0, nil
	}

	changedCount := 0
	now := time.Now().Unix()

	// Cache flags we've already fetched to avoid repeated DB lookups
	flagCache := make(map[string]*proto.Feature)

	for _, schedule := range schedules {
		// Only process schedules that reference the updated flag
		// (via prerequisites or FEATURE_FLAG clauses in rules)
		if !scheduleReferencesFlag(schedule.Payload, updatedFlagID) {
			continue
		}

		// Get the schedule's own flag
		scheduleFlagID := schedule.FeatureId
		flag, ok := flagCache[scheduleFlagID]
		if !ok {
			flagDomain, err := d.featureStorage.GetFeature(
				ctx, scheduleFlagID, environmentID,
			)
			if err != nil {
				d.logger.Warn(
					"Failed to get feature for cross-flag conflict check",
					zap.String("featureId", scheduleFlagID),
					zap.String("scheduleId", schedule.Id),
					zap.Error(err),
				)
				continue
			}
			flag = flagDomain.Feature
			flagCache[scheduleFlagID] = flag
		}

		// Validate same-flag references
		invalidRefs := validatePayloadReferences(
			flag, schedule.Payload, now,
		)

		// Validate cross-flag prerequisite references
		prereqConflicts := d.validatePrerequisiteReferences(
			ctx, schedule.Payload, environmentID, now,
		)
		invalidRefs = append(invalidRefs, prereqConflicts...)

		// Validate cross-flag FEATURE_FLAG references in rule clauses
		featureFlagConflicts := d.validateFeatureFlagReferences(
			ctx, schedule.Payload, environmentID, now,
		)
		invalidRefs = append(invalidRefs, featureFlagConflicts...)

		isPending := schedule.Status ==
			proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING
		isConflict := schedule.Status ==
			proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT

		if len(invalidRefs) > 0 && isPending {
			sfcDomain := &featuredomain.ScheduledFlagChange{
				ScheduledFlagChange: schedule,
			}
			sfcDomain.MarkConflict(invalidRefs)
			if err := d.storage.UpdateScheduledFlagChange(
				ctx, sfcDomain,
			); err != nil {
				d.logger.Warn(
					"Failed to mark schedule as CONFLICT",
					zap.String("scheduleId", schedule.Id),
					zap.Error(err),
				)
				continue
			}
			changedCount++
		} else if len(invalidRefs) == 0 && isConflict {
			sfcDomain := &featuredomain.ScheduledFlagChange{
				ScheduledFlagChange: schedule,
			}
			sfcDomain.RestoreToPending()
			if err := d.storage.UpdateScheduledFlagChange(
				ctx, sfcDomain,
			); err != nil {
				d.logger.Warn(
					"Failed to restore schedule to PENDING",
					zap.String("scheduleId", schedule.Id),
					zap.Error(err),
				)
				continue
			}
			changedCount++
		}
	}

	return changedCount, nil
}

// validatePrerequisiteReferences checks if prerequisite changes in the payload
// reference flags and variations that still exist.
func (d *ConflictDetector) validatePrerequisiteReferences(
	ctx context.Context,
	payload *proto.ScheduledChangePayload,
	environmentID string,
	now int64,
) []*proto.ScheduledChangeConflict {
	if payload == nil || d.featureStorage == nil {
		return nil
	}

	var conflicts []*proto.ScheduledChangeConflict
	for _, pc := range payload.PrerequisiteChanges {
		if pc == nil || pc.Prerequisite == nil {
			continue
		}
		if pc.ChangeType != proto.ChangeType_CREATE && pc.ChangeType != proto.ChangeType_UPDATE {
			continue
		}

		prereqFeature, err := d.featureStorage.GetFeature(
			ctx, pc.Prerequisite.FeatureId, environmentID,
		)
		if err != nil {
			d.logger.Warn(
				"Failed to get prerequisite feature for conflict check",
				zap.String("prerequisiteFeatureId", pc.Prerequisite.FeatureId),
				zap.Error(err),
			)
			conflicts = append(conflicts, &proto.ScheduledChangeConflict{
				Type: proto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE,
				Description: fmt.Sprintf(
					"Prerequisite flag %s not found",
					pc.Prerequisite.FeatureId,
				),
				ConflictingField: "prerequisites",
				DetectedAt:       now,
			})
			continue
		}

		if pc.Prerequisite.VariationId != "" &&
			!featuredomain.VariationExists(prereqFeature.Feature, pc.Prerequisite.VariationId) {
			conflicts = append(conflicts, &proto.ScheduledChangeConflict{
				Type: proto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE,
				Description: fmt.Sprintf(
					"Prerequisite flag %s no longer has variation %s",
					pc.Prerequisite.FeatureId, pc.Prerequisite.VariationId,
				),
				ConflictingField: "prerequisites",
				DetectedAt:       now,
			})
		}
	}
	return conflicts
}

// validateFeatureFlagReferences checks if FEATURE_FLAG clauses in rules
// reference flags and variations that still exist.
// The FEATURE_FLAG operator is used to create dependencies on other flags:
// - clause.Attribute contains the referenced feature ID
// - clause.Values contains the variation IDs of that flag
func (d *ConflictDetector) validateFeatureFlagReferences(
	ctx context.Context,
	payload *proto.ScheduledChangePayload,
	environmentID string,
	now int64,
) []*proto.ScheduledChangeConflict {
	if payload == nil || d.featureStorage == nil {
		return nil
	}

	var conflicts []*proto.ScheduledChangeConflict
	for _, rc := range payload.RuleChanges {
		if rc == nil || rc.Rule == nil {
			continue
		}
		if rc.ChangeType != proto.ChangeType_CREATE && rc.ChangeType != proto.ChangeType_UPDATE {
			continue
		}

		for _, clause := range rc.Rule.Clauses {
			if clause == nil || clause.Operator != proto.Clause_FEATURE_FLAG {
				continue
			}

			referencedFlagID := clause.Attribute
			if referencedFlagID == "" {
				continue
			}

			// Get the referenced feature
			referencedFeature, err := d.featureStorage.GetFeature(
				ctx, referencedFlagID, environmentID,
			)
			if err != nil {
				d.logger.Warn(
					"Failed to get referenced feature for FEATURE_FLAG clause",
					zap.String("referencedFeatureId", referencedFlagID),
					zap.String("ruleId", rc.Rule.Id),
					zap.Error(err),
				)
				conflicts = append(conflicts, &proto.ScheduledChangeConflict{
					Type: proto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE,
					Description: fmt.Sprintf(
						"Rule references flag %s which no longer exists",
						referencedFlagID,
					),
					ConflictingField: "rules",
					DetectedAt:       now,
				})
				continue
			}

			// Validate that all referenced variation IDs exist in the referenced flag
			for _, variationID := range clause.Values {
				if variationID != "" &&
					!featuredomain.VariationExists(referencedFeature.Feature, variationID) {
					conflicts = append(conflicts, &proto.ScheduledChangeConflict{
						Type: proto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE,
						Description: fmt.Sprintf(
							"Rule references flag %s variation %s which no longer exists",
							referencedFlagID, variationID,
						),
						ConflictingField: "rules",
						DetectedAt:       now,
					})
				}
			}
		}
	}
	return conflicts
}

// hasPrerequisiteChanges returns true if the payload has any prerequisite CREATE/UPDATE changes.
func hasPrerequisiteChanges(payload *proto.ScheduledChangePayload) bool {
	if payload == nil {
		return false
	}
	for _, pc := range payload.PrerequisiteChanges {
		if pc != nil && (pc.ChangeType == proto.ChangeType_CREATE || pc.ChangeType == proto.ChangeType_UPDATE) {
			return true
		}
	}
	return false
}

// scheduleReferencesFlag checks if a schedule's payload references a specific flag ID.
// This includes both prerequisites and flag-as-rule (FEATURE_FLAG operator in clauses).
// This is used to determine if a cross-flag conflict check should be performed.
func scheduleReferencesFlag(payload *proto.ScheduledChangePayload, targetFlagID string) bool {
	if payload == nil {
		return false
	}

	// Check if any prerequisite changes reference the target flag
	for _, pc := range payload.PrerequisiteChanges {
		if pc != nil && pc.Prerequisite != nil &&
			pc.Prerequisite.FeatureId == targetFlagID {
			return true
		}
	}

	// Check if any rule clauses reference the target flag via FEATURE_FLAG operator
	// In this case, the clause's attribute field contains the referenced feature ID
	for _, rc := range payload.RuleChanges {
		if rc != nil && rc.Rule != nil {
			for _, clause := range rc.Rule.Clauses {
				if clause != nil &&
					clause.Operator == proto.Clause_FEATURE_FLAG &&
					clause.Attribute == targetFlagID {
					return true
				}
			}
		}
	}

	return false
}

// listCrossFlagSchedules lists pending/conflict schedules in the environment
// for flags OTHER than the specified flag ID.
func (d *ConflictDetector) listCrossFlagSchedules(
	ctx context.Context,
	excludeFeatureID, environmentID string,
) ([]*proto.ScheduledFlagChange, error) {
	filters := []*mysql.FilterV2{
		{
			Column:   "environment_id",
			Operator: mysql.OperatorEqual,
			Value:    environmentID,
		},
		{
			Column:   "feature_id",
			Operator: mysql.OperatorNotEqual,
			Value:    excludeFeatureID,
		},
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
		Orders: []*mysql.Order{
			mysql.NewOrder("scheduled_at", mysql.OrderDirectionAsc),
		},
		Limit:  maxCrossFlagSchedulesToCheck,
		Offset: mysql.QueryNoOffset,
	}
	sfcs, _, _, err := d.storage.ListScheduledFlagChanges(
		ctx, options,
	)
	return sfcs, err
}

func (d *ConflictDetector) listPendingAndConflictSchedules(
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

// validatePayloadReferences validates that all references in the payload (variations, rules, targets, etc.)
// exist in the flag. Returns a list of conflicts if any references are invalid.
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
