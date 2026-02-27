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

package domain

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// ScheduledFlagChange is the domain model wrapper for scheduled flag changes
type ScheduledFlagChange struct {
	*proto.ScheduledFlagChange
}

// ChangeSummaryOptions provides optional lookup maps for richer summary text.
// When a lookup value is missing, summary generation falls back to IDs.
type ChangeSummaryOptions struct {
	SegmentNames            map[string]string
	CrossFlagVariationNames map[string]map[string]string
}

// NewScheduledFlagChange creates a new ScheduledFlagChange domain object
func NewScheduledFlagChange(
	featureID string,
	environmentID string,
	scheduledAt int64,
	timezone string,
	payload *proto.ScheduledChangePayload,
	comment string,
	flagVersionAtCreation int32,
	createdBy string,
) (*ScheduledFlagChange, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	if timezone == "" {
		timezone = "UTC"
	}
	return &ScheduledFlagChange{
		ScheduledFlagChange: &proto.ScheduledFlagChange{
			Id:                    id.String(),
			FeatureId:             featureID,
			EnvironmentId:         environmentID,
			ScheduledAt:           scheduledAt,
			Timezone:              timezone,
			Payload:               payload,
			Comment:               comment,
			Status:                proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
			FlagVersionAtCreation: flagVersionAtCreation,
			CreatedBy:             createdBy,
			CreatedAt:             now,
			UpdatedAt:             now,
		},
	}, nil
}

// UpdateSchedule updates the schedule time and timezone
func (s *ScheduledFlagChange) UpdateSchedule(scheduledAt int64, timezone string, updatedBy string) {
	s.ScheduledAt = scheduledAt
	if timezone != "" {
		s.Timezone = timezone
	}
	s.UpdatedBy = updatedBy
	s.UpdatedAt = time.Now().Unix()
}

// UpdatePayload updates the payload and comment
func (s *ScheduledFlagChange) UpdatePayload(payload *proto.ScheduledChangePayload, comment string, updatedBy string) {
	if payload != nil {
		s.Payload = payload
	}
	s.Comment = comment
	s.UpdatedBy = updatedBy
	s.UpdatedAt = time.Now().Unix()
}

// Cancel sets the status to CANCELLED
func (s *ScheduledFlagChange) Cancel(updatedBy string, reason string) {
	s.Status = proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CANCELLED
	s.FailureReason = reason
	s.UpdatedBy = updatedBy
	s.UpdatedAt = time.Now().Unix()
}

// MarkExecuted sets the status to EXECUTED
func (s *ScheduledFlagChange) MarkExecuted() {
	s.Status = proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_EXECUTED
	now := time.Now().Unix()
	s.ExecutedAt = now
	s.UpdatedAt = now
}

// MarkFailed sets the status to FAILED with a reason
func (s *ScheduledFlagChange) MarkFailed(reason string) {
	s.Status = proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_FAILED
	s.FailureReason = reason
	s.UpdatedAt = time.Now().Unix()
}

// MarkConflict sets the status to CONFLICT and records conflicts
func (s *ScheduledFlagChange) MarkConflict(conflicts []*proto.ScheduledChangeConflict) {
	s.Status = proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT
	s.Conflicts = conflicts
	s.UpdatedAt = time.Now().Unix()
}

// RestoreToPending restores a CONFLICT schedule back to PENDING status.
// This is used for auto-recovery when previously invalid references become valid again
// (e.g., a deleted variation is re-added to the flag).
func (s *ScheduledFlagChange) RestoreToPending() {
	s.Status = proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING
	s.Conflicts = nil
	s.UpdatedAt = time.Now().Unix()
}

// IsPending returns true if the status is PENDING
func (s *ScheduledFlagChange) IsPending() bool {
	return s.Status == proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING
}

// IsConflict returns true if the status is CONFLICT
func (s *ScheduledFlagChange) IsConflict() bool {
	return s.Status == proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT
}

// IsDue returns true if the scheduled time has passed
func (s *ScheduledFlagChange) IsDue() bool {
	return s.ScheduledAt <= time.Now().Unix()
}

// DetermineCategory computes the category based on the payload content.
// Note: ResetSamplingSeed is intentionally excluded from category calculation
// because it's a modifier that can be used from both Targeting and Variations tabs.
// This prevents confusing MIXED categorization when users make targeting/variation
// changes and also check "Reset Sampling Seed" on the same tab.
func (s *ScheduledFlagChange) DetermineCategory() proto.ScheduledChangeCategory {
	if s.Payload == nil {
		return proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_UNSPECIFIED
	}

	hasTargeting := len(s.Payload.RuleChanges) > 0 ||
		len(s.Payload.TargetChanges) > 0 ||
		len(s.Payload.PrerequisiteChanges) > 0 ||
		s.Payload.DefaultStrategy != nil

	hasVariations := len(s.Payload.VariationChanges) > 0 ||
		s.Payload.OffVariation != nil

	// ResetSamplingSeed is excluded from hasSettings because it's a modifier
	// that accompanies targeting/variation changes, not a standalone setting
	hasSettings := s.Payload.Enabled != nil ||
		s.Payload.Name != nil ||
		s.Payload.Description != nil ||
		len(s.Payload.TagChanges) > 0 ||
		s.Payload.Archived != nil ||
		s.Payload.Maintainer != nil

	categoryCount := 0
	if hasTargeting {
		categoryCount++
	}
	if hasVariations {
		categoryCount++
	}
	if hasSettings {
		categoryCount++
	}

	if categoryCount > 1 {
		return proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_MIXED
	}
	if hasTargeting {
		return proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_TARGETING
	}
	if hasVariations {
		return proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_VARIATIONS
	}
	if hasSettings {
		return proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_SETTINGS
	}

	// If only ResetSamplingSeed is set (no other changes), treat as SETTINGS
	if s.Payload.ResetSamplingSeed {
		return proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_SETTINGS
	}

	return proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_UNSPECIFIED
}

// i18n message keys for change summaries
// Frontend uses these keys to look up translations and interpolate values
const (
	MsgKeyEnableFlag              = "ScheduledChange.EnableFlag"
	MsgKeyDisableFlag             = "ScheduledChange.DisableFlag"
	MsgKeyRenameFlag              = "ScheduledChange.RenameFlag"
	MsgKeyUpdateDescription       = "ScheduledChange.UpdateDescription"
	MsgKeyChangeMaintainer        = "ScheduledChange.ChangeMaintainer"
	MsgKeyArchiveFlag             = "ScheduledChange.ArchiveFlag"
	MsgKeyUnarchiveFlag           = "ScheduledChange.UnarchiveFlag"
	MsgKeyResetSamplingSeed       = "ScheduledChange.ResetSamplingSeed"
	MsgKeyAddTag                  = "ScheduledChange.AddTag"
	MsgKeyRemoveTag               = "ScheduledChange.RemoveTag"
	MsgKeyAddVariation            = "ScheduledChange.AddVariation"
	MsgKeyUpdateVariation         = "ScheduledChange.UpdateVariation"
	MsgKeyChangeVariationValue    = "ScheduledChange.ChangeVariationValue"
	MsgKeyRenameVariation         = "ScheduledChange.RenameVariation"
	MsgKeyDeleteVariation         = "ScheduledChange.DeleteVariation"
	MsgKeyChangeOffVariation      = "ScheduledChange.ChangeOffVariation"
	MsgKeyAddRule                 = "ScheduledChange.AddRule"
	MsgKeyUpdateRule              = "ScheduledChange.UpdateRule"
	MsgKeyDeleteRule              = "ScheduledChange.DeleteRule"
	MsgKeyAddClauseToRule         = "ScheduledChange.AddClauseToRule"
	MsgKeyUpdateClauseInRule      = "ScheduledChange.UpdateClauseInRule"
	MsgKeyRemoveClauseFromRule    = "ScheduledChange.RemoveClauseFromRule"
	MsgKeyAddFeatureFlagClause    = "ScheduledChange.AddFeatureFlagClauseToRule"
	MsgKeyUpdateFeatureFlagClause = "ScheduledChange.UpdateFeatureFlagClauseInRule"
	MsgKeyRemoveFeatureFlagClause = "ScheduledChange.RemoveFeatureFlagClauseFromRule"
	MsgKeyTargetUsers             = "ScheduledChange.TargetUsers"
	MsgKeyRemoveTargeting         = "ScheduledChange.RemoveTargeting"
	MsgKeyAddPrerequisite         = "ScheduledChange.AddPrerequisite"
	MsgKeyUpdatePrerequisite      = "ScheduledChange.UpdatePrerequisite"
	MsgKeyRemovePrerequisite      = "ScheduledChange.RemovePrerequisite"
	MsgKeyChangeDefaultStrategy   = "ScheduledChange.ChangeDefaultStrategy"
)

// Cancellation reasons - currently plain text, will be i18n-ready in future
// TODO(i18n): Convert these to message keys when implementing i18n for failure reasons
const (
	CancelReasonUserCancelled = "Cancelled by user"
	CancelReasonFlagArchived  = "Flag was archived"
)

// newChangeSummary creates a new ChangeSummary with the given message key and values
func newChangeSummary(messageKey string, values map[string]string) *proto.ChangeSummary {
	return &proto.ChangeSummary{
		MessageKey: messageKey,
		Values:     values,
	}
}

// GenerateChangeSummaries creates i18n-ready summaries from the payload.
// Returns structured ChangeSummary objects with message keys and interpolation values.
// The frontend uses these to render localized messages.
// The flag parameter is optional and used to resolve variation/rule names.
func (s *ScheduledFlagChange) GenerateChangeSummaries(flag *proto.Feature) []*proto.ChangeSummary {
	return s.GenerateChangeSummariesWithOptions(flag, nil)
}

// GenerateChangeSummariesWithOptions creates i18n-ready summaries from the payload
// with optional lookups for segment names and cross-flag variation names.
func (s *ScheduledFlagChange) GenerateChangeSummariesWithOptions(
	flag *proto.Feature,
	options *ChangeSummaryOptions,
) []*proto.ChangeSummary {
	if s.Payload == nil {
		return nil
	}
	var summaries []*proto.ChangeSummary

	// Flag state changes (Settings)
	if s.Payload.Enabled != nil {
		if s.Payload.Enabled.Value {
			summaries = append(summaries, newChangeSummary(MsgKeyEnableFlag, nil))
		} else {
			summaries = append(summaries, newChangeSummary(MsgKeyDisableFlag, nil))
		}
	}

	// Name change (Settings)
	if s.Payload.Name != nil {
		summaries = append(summaries, newChangeSummary(MsgKeyRenameFlag, map[string]string{
			"name": s.Payload.Name.Value,
		}))
	}

	// Description change (Settings)
	if s.Payload.Description != nil {
		summaries = append(summaries, newChangeSummary(MsgKeyUpdateDescription, nil))
	}

	// Maintainer change (Settings)
	if s.Payload.Maintainer != nil {
		summaries = append(summaries, newChangeSummary(MsgKeyChangeMaintainer, map[string]string{
			"maintainer": s.Payload.Maintainer.Value,
		}))
	}

	// Archive change (Settings)
	if s.Payload.Archived != nil {
		if s.Payload.Archived.Value {
			summaries = append(summaries, newChangeSummary(MsgKeyArchiveFlag, nil))
		} else {
			summaries = append(summaries, newChangeSummary(MsgKeyUnarchiveFlag, nil))
		}
	}

	// Reset sampling seed (Settings)
	if s.Payload.ResetSamplingSeed {
		summaries = append(summaries, newChangeSummary(MsgKeyResetSamplingSeed, nil))
	}

	// Tag changes (Settings)
	for _, tc := range s.Payload.TagChanges {
		switch tc.ChangeType {
		case proto.ChangeType_CREATE:
			summaries = append(summaries, newChangeSummary(MsgKeyAddTag, map[string]string{
				"tag": tc.Tag,
			}))
		case proto.ChangeType_DELETE:
			summaries = append(summaries, newChangeSummary(MsgKeyRemoveTag, map[string]string{
				"tag": tc.Tag,
			}))
		}
	}

	// Variation changes (Variations)
	for _, vc := range s.Payload.VariationChanges {
		if vc.Variation == nil {
			continue // Skip invalid entries
		}
		switch vc.ChangeType {
		case proto.ChangeType_CREATE:
			summaries = append(summaries, newChangeSummary(MsgKeyAddVariation, map[string]string{
				"name":  vc.Variation.Name,
				"value": vc.Variation.Value,
			}))
		case proto.ChangeType_UPDATE:
			originalVar := sfcFindVariation(flag, vc.Variation.Id)
			if originalVar != nil {
				valueChanged := originalVar.Value != vc.Variation.Value
				nameChanged := originalVar.Name != vc.Variation.Name

				// Generate separate summaries for each type of change
				if valueChanged {
					summaries = append(summaries, newChangeSummary(MsgKeyChangeVariationValue, map[string]string{
						"name":     vc.Variation.Name,
						"oldValue": originalVar.Value,
						"newValue": vc.Variation.Value,
					}))
				}
				if nameChanged {
					summaries = append(summaries, newChangeSummary(MsgKeyRenameVariation, map[string]string{
						"oldName": originalVar.Name,
						"newName": vc.Variation.Name,
					}))
				}
				// If neither changed, it's some other update (e.g., description in future)
				if !valueChanged && !nameChanged {
					summaries = append(summaries, newChangeSummary(MsgKeyUpdateVariation, map[string]string{
						"name": vc.Variation.Name,
					}))
				}
			} else {
				// Original variation not found, show generic update
				summaries = append(summaries, newChangeSummary(MsgKeyUpdateVariation, map[string]string{
					"name": vc.Variation.Name,
				}))
			}
		case proto.ChangeType_DELETE:
			originalVar := sfcFindVariation(flag, vc.Variation.Id)
			name := vc.Variation.Id
			if originalVar != nil {
				name = originalVar.Name
			}
			summaries = append(summaries, newChangeSummary(MsgKeyDeleteVariation, map[string]string{
				"name": name,
			}))
		}
	}

	// Off variation change (Variations)
	if s.Payload.OffVariation != nil {
		variationName := sfcGetVariationName(flag, s.Payload.OffVariation.Value)
		summaries = append(summaries, newChangeSummary(MsgKeyChangeOffVariation, map[string]string{
			"name": variationName,
		}))
	}

	// Rule changes (Targeting)
	for _, rc := range s.Payload.RuleChanges {
		if rc.Rule == nil {
			continue // Skip invalid entries
		}
		switch rc.ChangeType {
		case proto.ChangeType_CREATE:
			summaries = append(summaries, newChangeSummary(MsgKeyAddRule, map[string]string{
				"description": sfcDescribeRule(rc.Rule, flag, options),
			}))
		case proto.ChangeType_UPDATE:
			originalRule := sfcFindRule(flag, rc.Rule.Id)
			if originalRule == nil {
				summaries = append(summaries, newChangeSummary(MsgKeyUpdateRule, map[string]string{
					"description": sfcDescribeRule(rc.Rule, flag, options),
				}))
				continue
			}

			clauseSummaries := sfcBuildRuleClauseSummaries(originalRule, rc.Rule, flag, options)
			if len(clauseSummaries) == 0 {
				summaries = append(summaries, newChangeSummary(MsgKeyUpdateRule, map[string]string{
					"description": sfcDescribeRule(rc.Rule, flag, options),
				}))
				continue
			}
			summaries = append(summaries, clauseSummaries...)
		case proto.ChangeType_DELETE:
			originalRule := sfcFindRule(flag, rc.Rule.Id)
			desc := rc.Rule.Id
			if originalRule != nil {
				desc = sfcDescribeRule(originalRule, flag, options)
			}
			summaries = append(summaries, newChangeSummary(MsgKeyDeleteRule, map[string]string{
				"description": desc,
			}))
		}
	}

	// Target changes (Targeting - individual user targeting)
	for _, tc := range s.Payload.TargetChanges {
		if tc.Target == nil {
			continue // Skip invalid entries
		}
		variationName := sfcGetVariationName(flag, tc.Target.Variation)
		switch tc.ChangeType {
		case proto.ChangeType_CREATE, proto.ChangeType_UPDATE:
			summaries = append(summaries, newChangeSummary(MsgKeyTargetUsers, map[string]string{
				"count":         fmt.Sprintf("%d", len(tc.Target.Users)),
				"variationName": variationName,
			}))
		case proto.ChangeType_DELETE:
			summaries = append(summaries, newChangeSummary(MsgKeyRemoveTargeting, map[string]string{
				"variationName": variationName,
			}))
		}
	}

	// Prerequisite changes (Targeting)
	for _, pc := range s.Payload.PrerequisiteChanges {
		if pc.Prerequisite == nil {
			continue // Skip invalid entries
		}
		switch pc.ChangeType {
		case proto.ChangeType_CREATE:
			summaries = append(summaries, newChangeSummary(MsgKeyAddPrerequisite, map[string]string{
				"featureId": pc.Prerequisite.FeatureId,
			}))
		case proto.ChangeType_UPDATE:
			summaries = append(summaries, newChangeSummary(MsgKeyUpdatePrerequisite, map[string]string{
				"featureId": pc.Prerequisite.FeatureId,
			}))
		case proto.ChangeType_DELETE:
			summaries = append(summaries, newChangeSummary(MsgKeyRemovePrerequisite, map[string]string{
				"featureId": pc.Prerequisite.FeatureId,
			}))
		}
	}

	// Default strategy change (Targeting)
	if s.Payload.DefaultStrategy != nil {
		summaries = append(summaries, newChangeSummary(MsgKeyChangeDefaultStrategy, map[string]string{
			"description": sfcDescribeStrategy(s.Payload.DefaultStrategy, flag),
		}))
	}

	return summaries
}

// CountChanges returns the total number of changes in the payload
func (s *ScheduledFlagChange) CountChanges() int {
	if s.Payload == nil {
		return 0
	}
	count := 0
	count += len(s.Payload.RuleChanges)
	count += len(s.Payload.TargetChanges)
	count += len(s.Payload.PrerequisiteChanges)
	count += len(s.Payload.VariationChanges)
	count += len(s.Payload.TagChanges)
	if s.Payload.DefaultStrategy != nil {
		count++
	}
	if s.Payload.OffVariation != nil {
		count++
	}
	if s.Payload.Enabled != nil {
		count++
	}
	if s.Payload.Name != nil {
		count++
	}
	if s.Payload.Description != nil {
		count++
	}
	if s.Payload.Archived != nil {
		count++
	}
	if s.Payload.ResetSamplingSeed {
		count++
	}
	if s.Payload.Maintainer != nil {
		count++
	}
	return count
}

// IsEmptyPayload returns true if the payload has no changes
func (s *ScheduledFlagChange) IsEmptyPayload() bool {
	return s.CountChanges() == 0
}

// Helper functions for generating summaries (prefixed with sfc_ to avoid conflicts)

func sfcFindVariation(flag *proto.Feature, variationID string) *proto.Variation {
	if flag == nil {
		return nil
	}
	for _, v := range flag.Variations {
		if v.Id == variationID {
			return v
		}
	}
	return nil
}

func sfcFindRule(flag *proto.Feature, ruleID string) *proto.Rule {
	if flag == nil {
		return nil
	}
	for _, r := range flag.Rules {
		if r.Id == ruleID {
			return r
		}
	}
	return nil
}

func sfcGetVariationName(flag *proto.Feature, variationID string) string {
	v := sfcFindVariation(flag, variationID)
	if v != nil {
		return v.Name
	}
	return variationID
}

func sfcDescribeRule(rule *proto.Rule, flag *proto.Feature, options *ChangeSummaryOptions) string {
	if rule == nil || len(rule.Clauses) == 0 {
		return "(no conditions)"
	}
	return sfcDescribeClause(rule.Clauses[0], flag, options)
}

func sfcBuildRuleClauseSummaries(
	originalRule *proto.Rule,
	newRule *proto.Rule,
	flag *proto.Feature,
	options *ChangeSummaryOptions,
) []*proto.ChangeSummary {
	ruleLabel := sfcRuleLabel(flag, newRule.Id)
	var summaries []*proto.ChangeSummary

	oldClauseByKey := make(map[string]*proto.Clause, len(originalRule.Clauses))
	for i, clause := range originalRule.Clauses {
		oldClauseByKey[sfcClauseKey(clause, i)] = clause
	}
	newClauseByKey := make(map[string]*proto.Clause, len(newRule.Clauses))
	for i, clause := range newRule.Clauses {
		newClauseByKey[sfcClauseKey(clause, i)] = clause
	}

	for i, newClause := range newRule.Clauses {
		key := sfcClauseKey(newClause, i)
		oldClause, exists := oldClauseByKey[key]
		if !exists {
			msgKey := MsgKeyAddClauseToRule
			if sfcIsFeatureFlagClause(newClause) {
				msgKey = MsgKeyAddFeatureFlagClause
			}
			summaries = append(summaries, newChangeSummary(msgKey, map[string]string{
				"rule":   ruleLabel,
				"clause": sfcDescribeClause(newClause, flag, options),
			}))
			continue
		}
		if sfcIsSameClause(oldClause, newClause) {
			continue
		}
		msgKey := MsgKeyUpdateClauseInRule
		if sfcIsFeatureFlagClause(newClause) || sfcIsFeatureFlagClause(oldClause) {
			msgKey = MsgKeyUpdateFeatureFlagClause
		}
		summaries = append(summaries, newChangeSummary(msgKey, map[string]string{
			"rule":      ruleLabel,
			"oldClause": sfcDescribeClause(oldClause, flag, options),
			"newClause": sfcDescribeClause(newClause, flag, options),
		}))
	}

	for i, oldClause := range originalRule.Clauses {
		key := sfcClauseKey(oldClause, i)
		if _, exists := newClauseByKey[key]; exists {
			continue
		}
		msgKey := MsgKeyRemoveClauseFromRule
		if sfcIsFeatureFlagClause(oldClause) {
			msgKey = MsgKeyRemoveFeatureFlagClause
		}
		summaries = append(summaries, newChangeSummary(msgKey, map[string]string{
			"rule":   ruleLabel,
			"clause": sfcDescribeClause(oldClause, flag, options),
		}))
	}
	return summaries
}

func sfcRuleLabel(flag *proto.Feature, ruleID string) string {
	if flag == nil {
		return "rule"
	}
	for i, rule := range flag.Rules {
		if rule.Id == ruleID {
			return fmt.Sprintf("rule #%d", i+1)
		}
	}
	return "rule"
}

func sfcClauseKey(clause *proto.Clause, index int) string {
	if clause == nil {
		return fmt.Sprintf("idx-%d", index)
	}
	if clause.Id != "" {
		return clause.Id
	}
	return fmt.Sprintf(
		"idx-%d:%s:%d",
		index,
		clause.Attribute,
		int(clause.Operator),
	)
}

func sfcIsSameClause(oldClause, newClause *proto.Clause) bool {
	if oldClause == nil || newClause == nil {
		return oldClause == nil && newClause == nil
	}
	if oldClause.Attribute != newClause.Attribute ||
		oldClause.Operator != newClause.Operator ||
		len(oldClause.Values) != len(newClause.Values) {
		return false
	}
	oldValues := append([]string(nil), oldClause.Values...)
	newValues := append([]string(nil), newClause.Values...)
	sort.Strings(oldValues)
	sort.Strings(newValues)
	for i := range oldValues {
		if oldValues[i] != newValues[i] {
			return false
		}
	}
	return true
}

func sfcIsFeatureFlagClause(clause *proto.Clause) bool {
	return clause != nil && clause.Operator == proto.Clause_FEATURE_FLAG
}

func sfcDescribeClause(clause *proto.Clause, flag *proto.Feature, options *ChangeSummaryOptions) string {
	if clause == nil {
		return "(empty clause)"
	}
	switch clause.Operator {
	case proto.Clause_SEGMENT:
		return fmt.Sprintf(
			"segment match: %s",
			strings.Join(sfcResolveSegmentDisplayValues(clause.Values, options), ", "),
		)
	case proto.Clause_FEATURE_FLAG:
		if clause.Attribute == "" {
			return "feature flag clause"
		}
		if len(clause.Values) == 0 {
			return fmt.Sprintf("flag %s condition", clause.Attribute)
		}
		resolvedValues := make([]string, 0, len(clause.Values))
		for _, variationID := range clause.Values {
			resolvedValues = append(
				resolvedValues,
				sfcResolveFeatureFlagVariationDisplayName(clause.Attribute, variationID, options),
			)
		}
		return fmt.Sprintf("flag %s serves \"%s\"", clause.Attribute, strings.Join(resolvedValues, ", "))
	default:
		if clause.Attribute == "segment" {
			return fmt.Sprintf(
				"segment match: %s",
				strings.Join(sfcResolveSegmentDisplayValues(clause.Values, options), ", "),
			)
		}
		return fmt.Sprintf(
			"%s %s %s",
			clause.Attribute,
			clause.Operator.String(),
			strings.Join(clause.Values, ", "),
		)
	}
}

func sfcResolveSegmentDisplayValues(values []string, options *ChangeSummaryOptions) []string {
	resolved := make([]string, 0, len(values))
	for _, segmentID := range values {
		if options != nil && options.SegmentNames != nil {
			if name, ok := options.SegmentNames[segmentID]; ok && name != "" {
				resolved = append(resolved, name)
				continue
			}
		}
		resolved = append(resolved, segmentID)
	}
	return resolved
}

func sfcResolveFeatureFlagVariationDisplayName(
	featureID string,
	variationID string,
	options *ChangeSummaryOptions,
) string {
	if options == nil || options.CrossFlagVariationNames == nil {
		return variationID
	}
	variationNamesByID, ok := options.CrossFlagVariationNames[featureID]
	if !ok {
		return variationID
	}
	variationName, ok := variationNamesByID[variationID]
	if !ok || variationName == "" || variationName == variationID {
		return variationID
	}
	return fmt.Sprintf("%s (%s)", variationName, variationID)
}

func sfcDescribeStrategy(strategy *proto.Strategy, flag *proto.Feature) string {
	if strategy == nil {
		return "Unknown strategy"
	}
	if strategy.Type == proto.Strategy_FIXED && strategy.FixedStrategy != nil {
		variationName := sfcGetVariationName(flag, strategy.FixedStrategy.Variation)
		return fmt.Sprintf("Serve \"%s\"", variationName)
	}
	if strategy.Type == proto.Strategy_ROLLOUT && strategy.RolloutStrategy != nil {
		var parts []string
		for _, v := range strategy.RolloutStrategy.Variations {
			variationName := sfcGetVariationName(flag, v.Variation)
			// Weight is in thousandths (100000 = 100%)
			percentage := float64(v.Weight) / 1000
			parts = append(parts, fmt.Sprintf("%s: %.1f%%", variationName, percentage))
		}
		return fmt.Sprintf("Rollout (%s)", strings.Join(parts, ", "))
	}
	return "Unknown strategy"
}
