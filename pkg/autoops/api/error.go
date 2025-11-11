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
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
)

var (
	statusInternal = api.NewGRPCStatus(
		pkgErr.NewErrorInternal(pkgErr.AutoopsPackageName, "internal"),
	)
	statusUnknownOpsType = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgUnknown(pkgErr.AutoopsPackageName, "unknown ops type", "AutoOperationType"),
	)
	statusInvalidCursor = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AutoopsPackageName, "cursor is invalid", "Cursor"),
	)
	statusAutoOpsRuleIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AutoopsPackageName, "id must be specified", "ID"),
	)
	statusFeatureIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AutoopsPackageName, "featureId must be specified", "FeatureFlagID"),
	)
	statusClauseRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AutoopsPackageName, "at least one clause must be specified", "Clause"),
	)
	statusClauseRequiredForDateTime = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName, "at least one date time clause must be specified", "Datetime"),
	)
	statusClauseRequiredForEventRate = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName, "at least one event rate clause must be specified", "EventRate"),
	)
	statusClauseIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AutoopsPackageName, "clause id must be specified", "ClauseId"),
	)
	statusClauseNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.AutoopsPackageName, "clause not found", "Clause"),
	)
	statusClauseAlreadyExecuted = api.NewGRPCStatus(
		pkgErr.NewErrorAlreadyExists(pkgErr.AutoopsPackageName, "clause is already executed"),
	)
	statusIncompatibleOpsType = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName,
			"ops type is incompatible with ops clause", "AutoOperationType",
		))
	statusAutoOpsRuleCompleted = api.NewGRPCStatus(
		pkgErr.NewErrorUnavailable(pkgErr.AutoopsPackageName, "auto ops rule is status of complete"),
	)
	statusAutoOpsRuleFinished = api.NewGRPCStatus(
		pkgErr.NewErrorUnavailable(pkgErr.AutoopsPackageName, "auto ops rule is status of finished"),
	)
	statusOpsEventRateClauseRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.AutoopsPackageName,
			"ops event rate clause must be specified", "EventRate",
		),
	)
	statusOpsEventRateClauseVariationIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"ops event rate clause variation id must be specified",
			"VariationId",
		),
	)
	statusOpsEventRateClauseGoalIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName, "ops event rate clause goal id is required", "Goal"))
	statusOpsEventRateClauseMinCountRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName, "ops event rate clause min count must be specified", "EventRate"))
	statusOpsEventRateClauseInvalidThredshold = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName,
			"ops event rate clause thredshold must be >0 and <=1",
			"Threshold",
		))
	statusDatetimeClauseRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.AutoopsPackageName, "datetime clause must be specified", "Datetime"))
	statusDatetimeClauseInvalidTime = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName,
			"datetime clause time must be after now timestamp",
			"Datetime",
		))
	statusDatetimeClauseDuplicateTime = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName,
			"datetime clause time must be unique",
			"Datetime",
		))
	statusAutoOpsRuleNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.AutoopsPackageName, "auto ops rule not found", "AutoOperation"))
	statusAutoOpsRuleAlreadyDeleted = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.AutoopsPackageName, "auto ops rule already deleted", "AutoOperation"))
	statusOpsEventRateClauseGoalNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(
			pkgErr.AutoopsPackageName,
			"ops event rate clause goal does not exist",
			"Goal",
		))
	statusAlreadyExists = api.NewGRPCStatus(
		pkgErr.NewErrorAlreadyExists(pkgErr.AutoopsPackageName, "already exists"))
	statusUnauthenticated = api.NewGRPCStatus(
		pkgErr.NewErrorUnauthenticated(pkgErr.AutoopsPackageName, "unauthenticated"))
	statusPermissionDenied = api.NewGRPCStatus(
		pkgErr.NewErrorPermissionDenied(pkgErr.AutoopsPackageName, "permission denied"))
	statusProgressiveRolloutFeatureIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"feature id must be specified for a progressive rollout",
			"FeatureFlagID",
		))
	statusProgressiveRolloutClauseRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(
			pkgErr.AutoopsPackageName,
			"at least one clause must be specified for a progressive rollout",
			"Clause",
		))
	statusIncorrectProgressiveRolloutClause = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName,
			"only one clause must be specified for a progressive rollout",
			"Clause",
		))
	statusProgressiveRolloutAlreadyStopped = api.NewGRPCStatus(
		pkgErr.NewErrorUnavailable(pkgErr.AutoopsPackageName, "progressive rollout is already stopped"))
	statusProgressiveRolloutClauseVariationIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"clause variation id must be specified for a progressive rollout",
			"VariationId",
		))
	statusProgressiveRolloutClauseInvalidVariationID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"the clause variation id set in the progressive rollout is invalid for a progressive rollout",
			"VariationId",
		))
	statusProgressiveRolloutClauseSchedulesRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"at least one clause schedule must be specified for a progressive rollout",
			"Datetime",
		))
	statusProgressiveRolloutClauseInvalidIncrements = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName,
			"increments is invalid for a progressive rollout",
			"Increments",
		))
	statusProgressiveRolloutClauseUnknownInterval = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgUnknown(
			pkgErr.AutoopsPackageName,
			"interval is unknown for a progressive rollout",
			"Interval",
		))
	statusProgressiveRolloutWaitingOrRunningExperimentExists = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(
			pkgErr.AutoopsPackageName,
			"cannot create a progressive rollout when there is a scheduled or running experiment",
		))
	statusProgressiveRolloutInsufficientVariations = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(
			pkgErr.AutoopsPackageName,
			"the feature must have at least 2 variations when creating a progressive rollout",
		))
	statusProgressiveRolloutControlVariationRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"control variation id must be specified for a progressive rollout",
			"control_variation_id",
		))
	statusProgressiveRolloutTargetVariationRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"target variation id must be specified for a progressive rollout",
			"target_variation_id",
		))
	statusProgressiveRolloutVariationsMustBeDifferent = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(
			pkgErr.AutoopsPackageName,
			"control and target variations must be different for a progressive rollout",
		))
	statusProgressiveRolloutControlVariationNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"the control variation id set in the progressive rollout does not exist in the feature",
			"control_variation_id",
		))
	statusProgressiveRolloutTargetVariationNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"the target variation id set in the progressive rollout does not exist in the feature",
			"target_variation_id",
		))
	statusProgressiveRolloutInvalidScheduleSpans = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(
			pkgErr.AutoopsPackageName,
			"the span of time for each scheduled time must be at least 5 minutes for a progressive rollout",
		))
	statusProgressiveRolloutScheduleExecutedAtRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"schedule executed_at must be specified for a progressive rollout",
			"Datetime",
		))
	statusProgressiveRolloutScheduleInvalidWeight = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"schedule weight is invalid for a progressive rollout",
			"ScheduleWeight",
		))
	statusProgressiveRolloutAlreadyExists = api.NewGRPCStatus(
		pkgErr.NewErrorAlreadyExists(pkgErr.AutoopsPackageName, "progressive rollout already exists"))
	statusProgressiveRolloutIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AutoopsPackageName, "id must be specified for a progressive rollout", "ID"))
	statusProgressiveRolloutNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(
			pkgErr.AutoopsPackageName, "progressive rollout not found", "ProgressiveRollout"))
	statusProgressiveRolloutInvalidCursor = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName, "cursor is invalid for a progressive rollout", "Cursor"))
	statusProgressiveRolloutInvalidOrderBy = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName, "order_by is invalid for a progressive rollout", "OrderBy"))
	statusProgressiveRolloutScheduleIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName, "schedule id must be specified for a progressive rollout", "Schedule"))
)
