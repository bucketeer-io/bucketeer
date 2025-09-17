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

package api

import (
	"github.com/bucketeer-io/bucketeer/pkg/api/api"
	pkgErr "github.com/bucketeer-io/bucketeer/pkg/error"
)

var (
	statusInternal = api.NewGRPCStatus(
		pkgErr.NewErrorInternal(pkgErr.AutoopsPackageName, "internal"),
	)
	statusUnknownOpsType = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgUnknown(pkgErr.AutoopsPackageName, "unknown ops type", "ops_type"),
	)
	statusInvalidCursor = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AutoopsPackageName, "cursor is invalid", "cursor"),
	)
	statusIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AutoopsPackageName, "id must be specified", "id"),
	)
	statusFeatureIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AutoopsPackageName, "featureId must be specified", "featureId"),
	)
	statusClauseRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AutoopsPackageName, "at least one clause must be specified", "clause"),
	)
	statusClauseRequiredForDateTime = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName, "at least one date time clause must be specified", "date_time_clause"),
	)
	statusClauseRequiredForEventDate = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName, "at least one event rate clause must be specified", "event_rate_clause"),
	)
	statusClauseIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AutoopsPackageName, "clause id must be specified", "clause_id"),
	)
	statusClauseNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.AutoopsPackageName, "clause not found", "clause"),
	)
	statusClauseAlreadyExecuted = api.NewGRPCStatus(
		pkgErr.NewErrorAlreadyExists(pkgErr.AutoopsPackageName, "clause is already executed"),
	)
	statusIncompatibleOpsType = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName,
			"ops type is incompatible with ops clause", "ops_type",
		))
	statusShouldAddMoreClauses = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName, "if existing clauses are deleted all, should add one or more clauses.", "clause"),
	)
	statusAutoOpsRuleCompleted = api.NewGRPCStatus(
		pkgErr.NewErrorUnavailable(pkgErr.AutoopsPackageName, "auto ops rule is status of complete"),
	)
	statusAutoOpsRuleFinished = api.NewGRPCStatus(
		pkgErr.NewErrorUnavailable(pkgErr.AutoopsPackageName, "auto ops rule is status of finished"),
	)
	statusOpsEventRateClauseRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.AutoopsPackageName,
			"ops event rate clause must be specified", "ops_event_rate_clause",
		),
	)
	statusOpsEventRateClauseVariationIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"ops event rate clause variation id must be specified",
			"ops_event_rate_clause_variation_id",
		),
	)
	statusOpsEventRateClauseGoalIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName, "ops event rate clause goal id is required", "ops_event_rate_clause_goal_id"))
	statusOpsEventRateClauseMinCountRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName, "ops event rate clause min count must be specified", "ops_event_rate_clause_min_count"))
	statusOpsEventRateClauseInvalidThredshold = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName,
			"ops event rate clause thredshold must be >0 and <=1",
			"ops_event_rate_clause_thredshold",
		))
	statusDatetimeClauseRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.AutoopsPackageName, "datetime clause must be specified", "datetime_clause"))
	statusDatetimeClauseInvalidTime = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName,
			"datetime clause time must be after now timestamp",
			"datetime_clause_time",
		))
	statusDatetimeClauseDuplicateTime = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName,
			"datetime clause time must be unique",
			"datetime_clause_time",
		))
	statusNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.AutoopsPackageName, "not found", "auto_ops_rule"))
	statusAlreadyDeleted = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.AutoopsPackageName, "already deleted", "deleted_auto_ops_rule"))
	statusOpsEventRateClauseGoalNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(
			pkgErr.AutoopsPackageName,
			"ops event rate clause goal does not exist",
			"ops_event_rate_clause_goal",
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
			"feature_id",
		))
	statusProgressiveRolloutClauseRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(
			pkgErr.AutoopsPackageName,
			"at least one clause must be specified for a progressive rollout",
			"clause",
		))
	statusIncorrectProgressiveRolloutClause = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName,
			"only one clause must be specified for a progressive rollout",
			"clause",
		))
	statusProgressiveRolloutInternal = api.NewGRPCStatus(
		pkgErr.NewErrorInternal(pkgErr.AutoopsPackageName, "internal error occurs for a progressive rollout"))
	statusProgressiveRolloutAlreadyStopped = api.NewGRPCStatus(
		pkgErr.NewErrorUnavailable(pkgErr.AutoopsPackageName, "progressive rollout is already stopped"))
	statusProgressiveRolloutClauseVariationIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"clause variation id must be specified for a progressive rollout",
			"clause_variation_id",
		))
	statusProgressiveRolloutClauseInvalidVariationID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"the clause variation id set in the progressive rollout is invalid for a progressive rollout",
			"clause_variation_id",
		))
	statusProgressiveRolloutClauseSchedulesRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"at least one clause schedule must be specified for a progressive rollout",
			"clause_schedule",
		))
	statusProgressiveRolloutClauseInvalidIncrements = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName,
			"increments is invalid for a progressive rollout",
			"increments",
		))
	statusProgressiveRolloutClauseUnknownInterval = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgUnknown(
			pkgErr.AutoopsPackageName,
			"interval is unknown for a progressive rollout",
			"interval",
		))
	statusProgressiveRolloutWaitingOrRunningExperimentExists = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(
			pkgErr.AutoopsPackageName,
			"cannot create a progressive rollout when there is a scheduled or running experiment",
		))
	statusProgressiveRolloutInvalidVariationSize = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(
			pkgErr.AutoopsPackageName,
			"the number of variations must be equal to 2 when creating a progressive rollout",
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
			"schedule_executed_at",
		))
	statusProgressiveRolloutScheduleInvalidWeight = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName,
			"schedule weight is invalid for a progressive rollout",
			"schedule_weight",
		))
	statusProgressiveRolloutAlreadyExists = api.NewGRPCStatus(
		pkgErr.NewErrorAlreadyExists(pkgErr.AutoopsPackageName, "progressive rollout already exists"))
	statusProgressiveRolloutIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AutoopsPackageName, "id must be specified for a progressive rollout", "id"))
	statusProgressiveRolloutNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(
			pkgErr.AutoopsPackageName, "progressive rollout not found", "progressive_rollout"))
	statusProgressiveRolloutInvalidCursor = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName, "cursor is invalid for a progressive rollout", "cursor"))
	statusProgressiveRolloutInvalidOrderBy = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AutoopsPackageName, "order_by is invalid for a progressive rollout", "order_by"))
	statusProgressiveRolloutScheduleIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.AutoopsPackageName, "schedule id must be specified for a progressive rollout", "schedule_id"))
)
