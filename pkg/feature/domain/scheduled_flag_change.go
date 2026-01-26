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
	"time"

	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// ScheduledFlagChange is the domain model wrapper for scheduled flag changes
type ScheduledFlagChange struct {
	*proto.ScheduledFlagChange
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

// DetermineCategory computes the category based on the payload content
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

	hasSettings := s.Payload.Enabled != nil ||
		s.Payload.Name != nil ||
		s.Payload.Description != nil ||
		len(s.Payload.TagChanges) > 0 ||
		s.Payload.Archived != nil ||
		s.Payload.ResetSamplingSeed ||
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

	return proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_UNSPECIFIED
}
