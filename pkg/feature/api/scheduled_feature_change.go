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

package api

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"go.uber.org/zap"

	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func (s *FeatureService) CreateScheduledFlagChange(
	ctx context.Context,
	req *ftproto.CreateScheduledFlagChangeRequest,
) (*ftproto.CreateScheduledFlagChangeResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	// Validate request (can be done outside transaction - no DB access)
	if err := s.validateCreateScheduledFlagChangeRequest(req); err != nil {
		return nil, err
	}

	var sfc *domain.ScheduledFlagChange
	var featureName string

	// Use transaction to ensure atomic check-and-create:
	// - Feature lookup, count check, and create must be atomic to prevent race conditions
	// - Without transaction, concurrent requests could exceed the max schedules limit
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		// Get the feature to validate it exists and get its version
		feature, err := s.featureStorage.GetFeature(ctxWithTx, req.FeatureId, req.EnvironmentId)
		if err != nil {
			if errors.Is(err, v2fs.ErrFeatureNotFound) {
				return statusFeatureNotFound.Err()
			}
			s.logger.Error(
				"Failed to get feature for scheduled change",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("featureId", req.FeatureId),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return statusInternal.Err()
		}
		featureName = feature.Name

		// Validate payload references
		if err := s.validateScheduledChangePayload(req.Payload, feature.Feature); err != nil {
			return err
		}

		// Check max schedules per flag limit
		pendingCount, err := s.countPendingSchedulesForFeature(
			ctxWithTx,
			req.FeatureId,
			req.EnvironmentId,
		)
		if err != nil {
			s.logger.Error(
				"Failed to count pending schedules",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("featureId", req.FeatureId),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return statusInternal.Err()
		}
		if pendingCount >= maxSchedulesPerFlag {
			return statusExceededMaxSchedulesPerFlag.Err()
		}

		// Create the scheduled flag change
		sfc, err = domain.NewScheduledFlagChange(
			req.FeatureId,
			req.EnvironmentId,
			req.ScheduledAt,
			req.Timezone,
			req.Payload,
			req.Comment,
			feature.Version,
			editor.Email,
		)
		if err != nil {
			s.logger.Error(
				"Failed to create scheduled flag change domain object",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("featureId", req.FeatureId),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return statusInternal.Err()
		}

		// Set derived fields
		sfc.Category = sfc.DetermineCategory()
		sfc.ChangeSummaries = sfc.GenerateChangeSummaries(feature.Feature)

		// Store the scheduled flag change
		if err := s.scheduledFlagChangeStorage.CreateScheduledFlagChange(ctxWithTx, sfc); err != nil {
			if errors.Is(err, v2fs.ErrScheduledFlagChangeAlreadyExists) {
				return statusAlreadyExists.Err()
			}
			s.logger.Error(
				"Failed to store scheduled flag change",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("featureId", req.FeatureId),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return statusInternal.Err()
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Publish domain event for audit log
	if err := s.publishScheduledFlagChangeCreatedEvent(
		ctx,
		editor,
		sfc,
		featureName,
		req.EnvironmentId,
	); err != nil {
		s.logger.Error(
			"Failed to publish scheduled flag change created event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", sfc.Id),
				zap.String("featureId", req.FeatureId),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		// Don't fail the request if event publishing fails
	}

	// TODO: Detect conflicts with existing schedules (Phase 4)
	var detectedConflicts []*ftproto.ScheduledChangeConflict

	return &ftproto.CreateScheduledFlagChangeResponse{
		ScheduledFlagChange: sfc.ScheduledFlagChange,
		DetectedConflicts:   detectedConflicts,
	}, nil
}

func (s *FeatureService) GetScheduledFlagChange(
	ctx context.Context,
	req *ftproto.GetScheduledFlagChangeRequest,
) (*ftproto.GetScheduledFlagChangeResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if req.Id == "" {
		return nil, statusMissingScheduledFlagChangeID.Err()
	}

	sfc, err := s.scheduledFlagChangeStorage.GetScheduledFlagChange(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2fs.ErrScheduledFlagChangeNotFound) {
			return nil, statusScheduledFlagChangeNotFound.Err()
		}
		s.logger.Error(
			"Failed to get scheduled flag change",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, statusInternal.Err()
	}

	// Enrich with derived fields if not already set
	if len(sfc.ChangeSummaries) == 0 {
		feature, _ := s.featureStorage.GetFeature(ctx, sfc.FeatureId, req.EnvironmentId)
		var featureProto *ftproto.Feature
		if feature != nil {
			featureProto = feature.Feature
		}
		sfc.Category = sfc.DetermineCategory()
		sfc.ChangeSummaries = sfc.GenerateChangeSummaries(featureProto)
	}

	return &ftproto.GetScheduledFlagChangeResponse{
		ScheduledFlagChange: sfc.ScheduledFlagChange,
	}, nil
}

func (s *FeatureService) UpdateScheduledFlagChange(
	ctx context.Context,
	req *ftproto.UpdateScheduledFlagChangeRequest,
) (*ftproto.UpdateScheduledFlagChangeResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if req.Id == "" {
		return nil, statusMissingScheduledFlagChangeID.Err()
	}

	// Validate scheduled_at if provided
	if req.ScheduledAt != nil {
		if err := validateScheduledTime(req.ScheduledAt.Value); err != nil {
			return nil, err
		}
	}

	var sfc *domain.ScheduledFlagChange
	var feature *domain.Feature
	var previousScheduledAt int64

	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		var err error
		sfc, err = s.scheduledFlagChangeStorage.GetScheduledFlagChange(ctxWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			if errors.Is(err, v2fs.ErrScheduledFlagChangeNotFound) {
				return statusScheduledFlagChangeNotFound.Err()
			}
			return err
		}

		// Capture previous scheduled time for audit event
		previousScheduledAt = sfc.ScheduledAt

		// Only allow updating pending schedules
		if !sfc.IsPending() && !sfc.IsConflict() {
			return statusScheduledFlagChangeNotPending.Err()
		}

		// Get feature for validation and summary generation
		feature, err = s.featureStorage.GetFeature(ctxWithTx, sfc.FeatureId, req.EnvironmentId)
		if err != nil {
			return err
		}

		// Update schedule time if provided
		if req.ScheduledAt != nil {
			timezone := ""
			if req.Timezone != nil {
				timezone = req.Timezone.Value
			}
			sfc.UpdateSchedule(req.ScheduledAt.Value, timezone, editor.Email)
		}

		// Update payload if provided
		if req.Payload != nil {
			// Validate payload
			if err := s.validateScheduledChangePayload(req.Payload, feature.Feature); err != nil {
				return err
			}
			comment := ""
			if req.Comment != nil {
				comment = req.Comment.Value
			}
			sfc.UpdatePayload(req.Payload, comment, editor.Email)
		} else if req.Comment != nil {
			// Only update comment
			sfc.Comment = req.Comment.Value
			sfc.UpdatedBy = editor.Email
			sfc.UpdatedAt = time.Now().Unix()
		}

		// Clear conflict status if schedule is being updated
		if sfc.Status == ftproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT {
			sfc.Status = ftproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING
			sfc.Conflicts = nil
		}

		// Update derived fields
		sfc.Category = sfc.DetermineCategory()
		sfc.ChangeSummaries = sfc.GenerateChangeSummaries(feature.Feature)

		return s.scheduledFlagChangeStorage.UpdateScheduledFlagChange(ctxWithTx, sfc)
	})

	if err != nil {
		s.logger.Error(
			"Failed to update scheduled flag change",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}

	// Publish domain event for audit log
	featureName := ""
	if feature != nil {
		featureName = feature.Name
	}
	if err := s.publishScheduledFlagChangeUpdatedEvent(
		ctx,
		editor,
		sfc,
		featureName,
		previousScheduledAt,
		req.EnvironmentId,
	); err != nil {
		s.logger.Error(
			"Failed to publish scheduled flag change updated event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", sfc.Id),
				zap.String("featureId", sfc.FeatureId),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		// Don't fail the request if event publishing fails
	}

	// TODO: Detect conflicts with existing schedules (Phase 4)
	var detectedConflicts []*ftproto.ScheduledChangeConflict

	return &ftproto.UpdateScheduledFlagChangeResponse{
		ScheduledFlagChange: sfc.ScheduledFlagChange,
		DetectedConflicts:   detectedConflicts,
	}, nil
}

func (s *FeatureService) DeleteScheduledFlagChange(
	ctx context.Context,
	req *ftproto.DeleteScheduledFlagChangeRequest,
) (*ftproto.DeleteScheduledFlagChangeResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if req.Id == "" {
		return nil, statusMissingScheduledFlagChangeID.Err()
	}

	var sfc *domain.ScheduledFlagChange
	var featureName string

	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		var err error
		sfc, err = s.scheduledFlagChangeStorage.GetScheduledFlagChange(ctxWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			if errors.Is(err, v2fs.ErrScheduledFlagChangeNotFound) {
				return statusScheduledFlagChangeNotFound.Err()
			}
			return err
		}

		// Only allow deleting pending or conflict schedules
		if !sfc.IsPending() && !sfc.IsConflict() {
			return statusScheduledFlagChangeNotPending.Err()
		}

		// Get feature name for audit event
		feature, err := s.featureStorage.GetFeature(ctxWithTx, sfc.FeatureId, req.EnvironmentId)
		if err == nil && feature != nil {
			featureName = feature.Name
		}

		// Cancel the schedule instead of hard delete (for audit trail)
		sfc.Cancel(editor.Email, "Cancelled by user")
		return s.scheduledFlagChangeStorage.UpdateScheduledFlagChange(ctxWithTx, sfc)
	})

	if err != nil {
		s.logger.Error(
			"Failed to delete scheduled flag change",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}

	// Publish domain event for audit log
	if err := s.publishScheduledFlagChangeCancelledEvent(
		ctx,
		editor,
		sfc,
		featureName,
		req.EnvironmentId,
	); err != nil {
		s.logger.Error(
			"Failed to publish scheduled flag change cancelled event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", sfc.Id),
				zap.String("featureId", sfc.FeatureId),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		// Don't fail the request if event publishing fails
	}

	return &ftproto.DeleteScheduledFlagChangeResponse{}, nil
}

func (s *FeatureService) ListScheduledFlagChanges(
	ctx context.Context,
	req *ftproto.ListScheduledFlagChangesRequest,
) (*ftproto.ListScheduledFlagChangesResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	// Build filters
	filters := []*mysql.FilterV2{
		{
			Column:   "environment_id",
			Operator: mysql.OperatorEqual,
			Value:    req.EnvironmentId,
		},
	}

	if req.FeatureId != "" {
		filters = append(filters, &mysql.FilterV2{
			Column:   "feature_id",
			Operator: mysql.OperatorEqual,
			Value:    req.FeatureId,
		})
	}

	if req.FromScheduledAt > 0 {
		filters = append(filters, &mysql.FilterV2{
			Column:   "scheduled_at",
			Operator: mysql.OperatorGreaterThanOrEqual,
			Value:    req.FromScheduledAt,
		})
	}

	if req.ToScheduledAt > 0 {
		filters = append(filters, &mysql.FilterV2{
			Column:   "scheduled_at",
			Operator: mysql.OperatorLessThanOrEqual,
			Value:    req.ToScheduledAt,
		})
	}

	// Status filter
	var inFilters []*mysql.InFilter
	if len(req.Statuses) > 0 {
		statusValues := make([]interface{}, 0, len(req.Statuses))
		for _, status := range req.Statuses {
			statusValues = append(statusValues, int32(status))
		}
		inFilters = append(inFilters, &mysql.InFilter{
			Column: "status",
			Values: statusValues,
		})
	}

	// Order
	var orders []*mysql.Order
	switch req.OrderBy {
	case ftproto.ListScheduledFlagChangesRequest_SCHEDULED_AT:
		direction := mysql.OrderDirectionAsc
		if req.OrderDirection == ftproto.ListScheduledFlagChangesRequest_DESC {
			direction = mysql.OrderDirectionDesc
		}
		orders = append(orders, mysql.NewOrder("scheduled_at", direction))
	case ftproto.ListScheduledFlagChangesRequest_CREATED_AT:
		direction := mysql.OrderDirectionAsc
		if req.OrderDirection == ftproto.ListScheduledFlagChangesRequest_DESC {
			direction = mysql.OrderDirectionDesc
		}
		orders = append(orders, mysql.NewOrder("created_at", direction))
	default:
		// Default: order by scheduled_at ASC
		orders = append(orders, mysql.NewOrder("scheduled_at", mysql.OrderDirectionAsc))
	}

	// Pagination
	limit := int(req.PageSize)
	if limit <= 0 || limit > maxPageSizePerRequest {
		limit = maxPageSizePerRequest
	}

	offset := 0
	if req.Cursor != "" {
		var err error
		offset, err = strconv.Atoi(req.Cursor)
		if err != nil {
			return nil, statusInvalidCursor.Err()
		}
	}

	options := &mysql.ListOptions{
		Filters:   filters,
		InFilters: inFilters,
		Orders:    orders,
		Limit:     limit,
		Offset:    offset,
	}

	sfcs, nextOffset, totalCount, err := s.scheduledFlagChangeStorage.ListScheduledFlagChanges(ctx, options)
	if err != nil {
		s.logger.Error(
			"Failed to list scheduled flag changes",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, statusInternal.Err()
	}

	// Enrich with derived fields
	for _, sfc := range sfcs {
		if len(sfc.ChangeSummaries) == 0 {
			sfcDomain := &domain.ScheduledFlagChange{ScheduledFlagChange: sfc}
			feature, _ := s.featureStorage.GetFeature(ctx, sfc.FeatureId, req.EnvironmentId)
			var featureProto *ftproto.Feature
			if feature != nil {
				featureProto = feature.Feature
			}
			sfc.Category = sfcDomain.DetermineCategory()
			sfc.ChangeSummaries = sfcDomain.GenerateChangeSummaries(featureProto)
		}
	}

	return &ftproto.ListScheduledFlagChangesResponse{
		ScheduledFlagChanges: sfcs,
		Cursor:               strconv.Itoa(nextOffset),
		TotalCount:           totalCount,
	}, nil
}

func (s *FeatureService) ExecuteScheduledFlagChange(
	ctx context.Context,
	req *ftproto.ExecuteScheduledFlagChangeRequest,
) (*ftproto.ExecuteScheduledFlagChangeResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if req.Id == "" {
		return nil, statusMissingScheduledFlagChangeID.Err()
	}

	var sfc *domain.ScheduledFlagChange

	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		var err error
		sfc, err = s.scheduledFlagChangeStorage.GetScheduledFlagChange(ctxWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			if errors.Is(err, v2fs.ErrScheduledFlagChangeNotFound) {
				return statusScheduledFlagChangeNotFound.Err()
			}
			return err
		}

		// Only allow executing pending schedules
		if !sfc.IsPending() {
			return statusScheduledFlagChangeNotPending.Err()
		}

		// Get the feature
		feature, err := s.featureStorage.GetFeature(ctxWithTx, sfc.FeatureId, req.EnvironmentId)
		if err != nil {
			if errors.Is(err, v2fs.ErrFeatureNotFound) {
				sfc.MarkFailed("Feature not found")
				_ = s.scheduledFlagChangeStorage.UpdateScheduledFlagChange(ctxWithTx, sfc)
				return statusFeatureNotFound.Err()
			}
			return err
		}

		// Validate references still exist
		if err := s.validateScheduledChangePayload(sfc.Payload, feature.Feature); err != nil {
			sfc.MarkFailed(err.Error())
			_ = s.scheduledFlagChangeStorage.UpdateScheduledFlagChange(ctxWithTx, sfc)
			return err
		}

		// Apply the changes using the existing UpdateFeature flow
		updateReq := convertPayloadToUpdateRequest(sfc.Payload, sfc.FeatureId, req.EnvironmentId)
		updateReq.Comment = "Applied from scheduled change: " + sfc.Comment

		_, err = s.UpdateFeature(ctxWithTx, updateReq)
		if err != nil {
			sfc.MarkFailed(err.Error())
			_ = s.scheduledFlagChangeStorage.UpdateScheduledFlagChange(ctxWithTx, sfc)
			return err
		}

		// Mark as executed
		sfc.MarkExecuted()
		sfc.UpdatedBy = editor.Email
		return s.scheduledFlagChangeStorage.UpdateScheduledFlagChange(ctxWithTx, sfc)
	})

	if err != nil {
		s.logger.Error(
			"Failed to execute scheduled flag change",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}

	return &ftproto.ExecuteScheduledFlagChangeResponse{
		ScheduledFlagChange: sfc.ScheduledFlagChange,
	}, nil
}

func (s *FeatureService) GetScheduledFlagChangeSummary(
	ctx context.Context,
	req *ftproto.GetScheduledFlagChangeSummaryRequest,
) (*ftproto.GetScheduledFlagChangeSummaryResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if req.FeatureId == "" {
		return nil, statusMissingFeatureID.Err()
	}

	// Get all pending and conflict schedules for this feature
	filters := []*mysql.FilterV2{
		{
			Column:   "environment_id",
			Operator: mysql.OperatorEqual,
			Value:    req.EnvironmentId,
		},
		{
			Column:   "feature_id",
			Operator: mysql.OperatorEqual,
			Value:    req.FeatureId,
		},
	}

	inFilters := []*mysql.InFilter{
		{
			Column: "status",
			Values: []interface{}{
				int32(ftproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING),
				int32(ftproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT),
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

	sfcs, _, _, err := s.scheduledFlagChangeStorage.ListScheduledFlagChanges(ctx, options)
	if err != nil {
		s.logger.Error(
			"Failed to get scheduled flag change summary",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", req.FeatureId),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, statusInternal.Err()
	}

	summary := &ftproto.ScheduledFlagChangeSummary{
		FeatureId: req.FeatureId,
	}

	for _, sfc := range sfcs {
		switch sfc.Status {
		case ftproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING:
			summary.PendingCount++
		case ftproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT:
			summary.ConflictCount++
		}
	}

	// Set next scheduled info from the first pending/conflict schedule
	if len(sfcs) > 0 {
		summary.NextScheduledAt = sfcs[0].ScheduledAt
		sfcDomain := &domain.ScheduledFlagChange{ScheduledFlagChange: sfcs[0]}
		summary.NextCategory = sfcDomain.DetermineCategory()
	}

	return &ftproto.GetScheduledFlagChangeSummaryResponse{
		Summary: summary,
	}, nil
}

// Helper functions

func (s *FeatureService) validateCreateScheduledFlagChangeRequest(req *ftproto.CreateScheduledFlagChangeRequest) error {
	if req.FeatureId == "" {
		return statusMissingFeatureID.Err()
	}
	if req.ScheduledAt == 0 {
		return statusMissingScheduledAt.Err()
	}
	if req.Payload == nil {
		return statusMissingPayload.Err()
	}

	// Validate scheduled time
	if err := validateScheduledTime(req.ScheduledAt); err != nil {
		return err
	}

	// Validate payload is not empty
	sfcDomain := &domain.ScheduledFlagChange{
		ScheduledFlagChange: &ftproto.ScheduledFlagChange{Payload: req.Payload},
	}
	if sfcDomain.IsEmptyPayload() {
		return statusEmptyPayload.Err()
	}

	// Validate max changes per schedule
	if sfcDomain.CountChanges() > maxChangesPerSchedule {
		return statusExceededMaxChangesPerSchedule.Err()
	}

	return nil
}

func validateScheduledTime(scheduledAt int64) error {
	now := time.Now().Unix()
	minTime := now + int64(minScheduleTimeMinutes*60)
	maxTime := now + int64(maxScheduleTimeDays*24*60*60)

	if scheduledAt < minTime {
		return statusScheduledTimeTooSoon.Err()
	}
	if scheduledAt > maxTime {
		return statusScheduledTimeTooFar.Err()
	}
	return nil
}

func (s *FeatureService) validateScheduledChangePayload(
	payload *ftproto.ScheduledChangePayload,
	feature *ftproto.Feature,
) error {
	if payload == nil {
		return nil
	}

	// Validate variation references
	for _, vc := range payload.VariationChanges {
		if vc.ChangeType == ftproto.ChangeType_UPDATE || vc.ChangeType == ftproto.ChangeType_DELETE {
			if !variationExists(feature, vc.Variation.Id) {
				return statusInvalidVariationReference.Err()
			}
		}
	}

	// Validate rule references
	for _, rc := range payload.RuleChanges {
		if rc.ChangeType == ftproto.ChangeType_UPDATE || rc.ChangeType == ftproto.ChangeType_DELETE {
			if !ruleExists(feature, rc.Rule.Id) {
				return statusInvalidRuleReference.Err()
			}
		}
	}

	// Validate off variation reference
	if payload.OffVariation != nil && payload.OffVariation.Value != "" {
		if !variationExists(feature, payload.OffVariation.Value) {
			return statusInvalidVariationReference.Err()
		}
	}

	// Validate default strategy variation references
	if payload.DefaultStrategy != nil {
		if payload.DefaultStrategy.FixedStrategy != nil {
			if !variationExists(feature, payload.DefaultStrategy.FixedStrategy.Variation) {
				return statusInvalidVariationReference.Err()
			}
		}
		if payload.DefaultStrategy.RolloutStrategy != nil {
			for _, rv := range payload.DefaultStrategy.RolloutStrategy.Variations {
				if !variationExists(feature, rv.Variation) {
					return statusInvalidVariationReference.Err()
				}
			}
		}
	}

	return nil
}

func variationExists(feature *ftproto.Feature, variationID string) bool {
	if feature == nil {
		return false
	}
	for _, v := range feature.Variations {
		if v.Id == variationID {
			return true
		}
	}
	return false
}

func ruleExists(feature *ftproto.Feature, ruleID string) bool {
	if feature == nil {
		return false
	}
	for _, r := range feature.Rules {
		if r.Id == ruleID {
			return true
		}
	}
	return false
}

func (s *FeatureService) countPendingSchedulesForFeature(
	ctx context.Context,
	featureID, environmentID string,
) (int, error) {
	filters := []*mysql.FilterV2{
		{
			Column:   "environment_id",
			Operator: mysql.OperatorEqual,
			Value:    environmentID,
		},
		{
			Column:   "feature_id",
			Operator: mysql.OperatorEqual,
			Value:    featureID,
		},
	}

	inFilters := []*mysql.InFilter{
		{
			Column: "status",
			Values: []interface{}{
				int32(ftproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING),
				int32(ftproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT),
			},
		},
	}

	options := &mysql.ListOptions{
		Filters:   filters,
		InFilters: inFilters,
		Limit:     mysql.QueryNoLimit,
		Offset:    mysql.QueryNoOffset,
	}

	_, _, totalCount, err := s.scheduledFlagChangeStorage.ListScheduledFlagChanges(ctx, options)
	if err != nil {
		return 0, err
	}
	return int(totalCount), nil
}

func convertPayloadToUpdateRequest(
	payload *ftproto.ScheduledChangePayload,
	featureID, environmentID string,
) *ftproto.UpdateFeatureRequest {
	req := &ftproto.UpdateFeatureRequest{
		EnvironmentId:       environmentID,
		Id:                  featureID,
		VariationChanges:    payload.VariationChanges,
		RuleChanges:         payload.RuleChanges,
		PrerequisiteChanges: payload.PrerequisiteChanges,
		TargetChanges:       payload.TargetChanges,
		TagChanges:          payload.TagChanges,
		DefaultStrategy:     payload.DefaultStrategy,
		OffVariation:        payload.OffVariation,
		Enabled:             payload.Enabled,
		Name:                payload.Name,
		Description:         payload.Description,
		Archived:            payload.Archived,
		ResetSamplingSeed:   payload.ResetSamplingSeed,
		Maintainer:          payload.Maintainer,
	}
	return req
}

// cancelPendingScheduledFlagChanges cancels all pending/conflict scheduled flag changes for a feature
func (s *FeatureService) cancelPendingScheduledFlagChanges(
	ctx context.Context,
	featureID, environmentID, cancelledBy, reason string,
) error {
	filters := []*mysql.FilterV2{
		{
			Column:   "environment_id",
			Operator: mysql.OperatorEqual,
			Value:    environmentID,
		},
		{
			Column:   "feature_id",
			Operator: mysql.OperatorEqual,
			Value:    featureID,
		},
	}

	inFilters := []*mysql.InFilter{
		{
			Column: "status",
			Values: []interface{}{
				int32(ftproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING),
				int32(ftproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT),
			},
		},
	}

	options := &mysql.ListOptions{
		Filters:   filters,
		InFilters: inFilters,
		Limit:     mysql.QueryNoLimit,
		Offset:    mysql.QueryNoOffset,
	}

	sfcs, _, _, err := s.scheduledFlagChangeStorage.ListScheduledFlagChanges(ctx, options)
	if err != nil {
		return err
	}

	for _, sfc := range sfcs {
		sfcDomain := &domain.ScheduledFlagChange{ScheduledFlagChange: sfc}
		sfcDomain.Cancel(cancelledBy, reason)
		if err := s.scheduledFlagChangeStorage.UpdateScheduledFlagChange(ctx, sfcDomain); err != nil {
			s.logger.Error(
				"Failed to cancel scheduled flag change",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("scheduledFlagChangeId", sfc.Id),
					zap.String("featureId", featureID),
					zap.String("environmentId", environmentID),
				)...,
			)
			// Continue cancelling other schedules
		}
	}

	return nil
}

// publishScheduledFlagChangeCreatedEvent publishes a domain event for audit logging
func (s *FeatureService) publishScheduledFlagChangeCreatedEvent(
	ctx context.Context,
	editor *eventproto.Editor,
	sfc *domain.ScheduledFlagChange,
	featureName string,
	environmentID string,
) error {
	payloadJSON, _ := json.Marshal(sfc.Payload)
	changeSummaries := changeSummariesToStrings(sfc.ChangeSummaries)

	event, err := domainevent.NewEvent(
		editor,
		eventproto.Event_SCHEDULED_FLAG_CHANGE,
		sfc.Id,
		eventproto.Event_SCHEDULED_FLAG_CHANGE_CREATED,
		&eventproto.ScheduledFlagChangeCreatedEvent{
			Id:              sfc.Id,
			FeatureId:       sfc.FeatureId,
			FeatureName:     featureName,
			ScheduledAt:     sfc.ScheduledAt,
			Timezone:        sfc.Timezone,
			Category:        sfc.Category.String(),
			ChangeSummaries: changeSummaries,
			PayloadJson:     string(payloadJSON),
			ScheduledBy:     sfc.CreatedBy,
		},
		environmentID,
		sfc.ScheduledFlagChange,
		nil,
	)
	if err != nil {
		return err
	}
	return s.domainPublisher.Publish(ctx, event)
}

// publishScheduledFlagChangeUpdatedEvent publishes a domain event for audit logging
func (s *FeatureService) publishScheduledFlagChangeUpdatedEvent(
	ctx context.Context,
	editor *eventproto.Editor,
	sfc *domain.ScheduledFlagChange,
	featureName string,
	previousScheduledAt int64,
	environmentID string,
) error {
	payloadJSON, _ := json.Marshal(sfc.Payload)
	changeSummaries := changeSummariesToStrings(sfc.ChangeSummaries)

	event, err := domainevent.NewEvent(
		editor,
		eventproto.Event_SCHEDULED_FLAG_CHANGE,
		sfc.Id,
		eventproto.Event_SCHEDULED_FLAG_CHANGE_UPDATED,
		&eventproto.ScheduledFlagChangeUpdatedEvent{
			Id:                  sfc.Id,
			FeatureId:           sfc.FeatureId,
			FeatureName:         featureName,
			PreviousScheduledAt: previousScheduledAt,
			NewScheduledAt:      sfc.ScheduledAt,
			Timezone:            sfc.Timezone,
			ChangeSummaries:     changeSummaries,
			PayloadJson:         string(payloadJSON),
		},
		environmentID,
		sfc.ScheduledFlagChange,
		nil,
	)
	if err != nil {
		return err
	}
	return s.domainPublisher.Publish(ctx, event)
}

// publishScheduledFlagChangeCancelledEvent publishes a domain event for audit logging
func (s *FeatureService) publishScheduledFlagChangeCancelledEvent(
	ctx context.Context,
	editor *eventproto.Editor,
	sfc *domain.ScheduledFlagChange,
	featureName string,
	environmentID string,
) error {
	changeSummaries := changeSummariesToStrings(sfc.ChangeSummaries)

	event, err := domainevent.NewEvent(
		editor,
		eventproto.Event_SCHEDULED_FLAG_CHANGE,
		sfc.Id,
		eventproto.Event_SCHEDULED_FLAG_CHANGE_CANCELLED,
		&eventproto.ScheduledFlagChangeCancelledEvent{
			Id:                    sfc.Id,
			FeatureId:             sfc.FeatureId,
			FeatureName:           featureName,
			ScheduledAt:           sfc.ScheduledAt,
			Timezone:              sfc.Timezone,
			ChangeSummaries:       changeSummaries,
			OriginallyScheduledBy: sfc.CreatedBy,
			OriginallyCreatedAt:   sfc.CreatedAt,
		},
		environmentID,
		sfc.ScheduledFlagChange,
		nil,
	)
	if err != nil {
		return err
	}
	return s.domainPublisher.Publish(ctx, event)
}

// changeSummariesToStrings converts structured ChangeSummary objects to simple strings for events
func changeSummariesToStrings(summaries []*ftproto.ChangeSummary) []string {
	result := make([]string, 0, len(summaries))
	for _, s := range summaries {
		// Use the message key as the string representation
		// The audit log UI can look up translations if needed
		result = append(result, s.MessageKey)
	}
	return result
}
