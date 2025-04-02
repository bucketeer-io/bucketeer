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
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

var (
	statusInternal          = gstatus.New(codes.Internal, "autoops: internal")
	statusUnknownOpsType    = gstatus.New(codes.Internal, "autoops: unknown ops type")
	statusInvalidCursor     = gstatus.New(codes.InvalidArgument, "autoops: cursor is invalid")
	statusIDRequired        = gstatus.New(codes.InvalidArgument, "autoops: id must be specified")
	statusFeatureIDRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: featureId must be specified",
	)
	statusClauseRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: at least one clause must be specified",
	)
	statusClauseRequiredForDateTime = gstatus.New(
		codes.InvalidArgument,
		"autoops: at least one date time clause must be specified",
	)
	statusClauseRequiredForEventDate = gstatus.New(
		codes.InvalidArgument,
		"autoops: at least one event rate clause must be specified",
	)
	statusClauseIDRequired      = gstatus.New(codes.InvalidArgument, "autoops: clause id must be specified")
	statusClauseNotFound        = gstatus.New(codes.NotFound, "autoops: clause not found")
	statusClauseAlreadyExecuted = gstatus.New(codes.InvalidArgument, "autoops: clause is already executed")
	statusIncompatibleOpsType   = gstatus.New(
		codes.InvalidArgument,
		"autoops: ops type is incompatible with ops clause",
	)
	statusShouldAddMoreClauses = gstatus.New(
		codes.InvalidArgument,
		"autoops: if existing clauses are deleted all, should add one or more clauses.",
	)
	statusAutoOpsRuleCompleted       = gstatus.New(codes.InvalidArgument, "autoops: auto ops rule is status of complete")
	statusAutoOpsRuleFinished        = gstatus.New(codes.InvalidArgument, "autoops: auto ops rule is status of finished")
	statusOpsEventRateClauseRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: ops event rate clause must be specified",
	)
	statusOpsEventRateClauseVariationIDRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: ops event rate clause variation id must be specified",
	)
	statusOpsEventRateClauseGoalIDRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: ops event rate clause goal id is required",
	)
	statusOpsEventRateClauseMinCountRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: ops event rate clause min count must be specified",
	)
	statusOpsEventRateClauseInvalidThredshold = gstatus.New(
		codes.InvalidArgument,
		"autoops: ops event rate clause thredshold must be >0 and <=1",
	)
	statusDatetimeClauseRequired    = gstatus.New(codes.InvalidArgument, "autoops: datetime clause must be specified")
	statusDatetimeClauseInvalidTime = gstatus.New(
		codes.InvalidArgument,
		"autoops: datetime clause time must be after now timestamp",
	)
	statusDatetimeClauseDuplicateTime = gstatus.New(
		codes.InvalidArgument,
		"autoops: datetime clause time must be unique",
	)
	statusNotFound                       = gstatus.New(codes.NotFound, "autoops: not found")
	statusAlreadyDeleted                 = gstatus.New(codes.NotFound, "autoops: already deleted")
	statusOpsEventRateClauseGoalNotFound = gstatus.New(
		codes.NotFound,
		"autoops: ops event rate clause goal does not exist",
	)
	statusAlreadyExists                       = gstatus.New(codes.AlreadyExists, "autoops: already exists")
	statusUnauthenticated                     = gstatus.New(codes.Unauthenticated, "autoops: unauthenticated")
	statusPermissionDenied                    = gstatus.New(codes.PermissionDenied, "autoops: permission denied")
	statusProgressiveRolloutFeatureIDRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: feature id must be specified for a progressive rollout",
	)
	statusProgressiveRolloutClauseRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: at least one clause must be specified for a progressive rollout",
	)
	statusIncorrectProgressiveRolloutClause = gstatus.New(
		codes.InvalidArgument,
		"autoops: only one clause must be specified for a progressive rollout",
	)
	statusProgressiveRolloutInternal = gstatus.New(
		codes.Internal,
		"autoops: internal error occurs for a progressive rollout",
	)
	statusProgressiveRolloutAlreadyStopped = gstatus.New(
		codes.Internal,
		"autoops: progressive rollout is already stopped",
	)
	statusProgressiveRolloutClauseVariationIDRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: clause variation id must be specified for a progressive rollout",
	)
	statusProgressiveRolloutClauseInvalidVariationID = gstatus.New(
		codes.InvalidArgument,
		"autoops: the clause variation id set in the progressive rollout is invalid for a progressive rollout",
	)
	statusProgressiveRolloutClauseSchedulesRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: at least one clause schedule must be specified for a progressive rollout",
	)
	statusProgressiveRolloutClauseInvalidIncrements = gstatus.New(
		codes.InvalidArgument,
		"autoops: increments is invalid for a progressive rollout",
	)
	statusProgressiveRolloutClauseUnknownInterval = gstatus.New(
		codes.InvalidArgument,
		"autoops: interval is unknown for a progressive rollout",
	)
	statusProgressiveRolloutWaitingOrRunningExperimentExists = gstatus.New(
		codes.FailedPrecondition,
		"autoops: cannot create a progressive rollout when there is a scheduled or running experiment",
	)
	statusProgressiveRolloutInvalidVariationSize = gstatus.New(
		codes.FailedPrecondition,
		"autoops progressive rollout: the number of variations must be equal to 2 when creating a progressive rollout",
	)
	statusProgressiveRolloutInvalidScheduleSpans = gstatus.New(
		codes.FailedPrecondition,
		"autoops: the span of time for each scheduled time must be at least 5 minutes for a progressive rollout",
	)
	statusProgressiveRolloutScheduleExecutedAtRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: schedule executed_at must be specified for a progressive rollout",
	)
	statusProgressiveRolloutScheduleInvalidWeight = gstatus.New(
		codes.InvalidArgument,
		"autoops: schedule weight is invalid for a progressive rollout",
	)
	statusProgressiveRolloutAlreadyExists = gstatus.New(
		codes.AlreadyExists,
		"autoops: progressive rollout already exists",
	)
	statusProgressiveRolloutIDRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: id must be specified for a progressive rollout",
	)
	statusProgressiveRolloutNotFound = gstatus.New(
		codes.NotFound,
		"autoops: progressive rollout not found",
	)
	statusProgressiveRolloutInvalidCursor = gstatus.New(
		codes.InvalidArgument,
		"autoops: cursor is invalid for a progressive rollout",
	)
	statusProgressiveRolloutInvalidOrderBy = gstatus.New(
		codes.InvalidArgument,
		"autoops: order_by is invalid for a progressive rollout",
	)
	statusProgressiveRolloutScheduleIDRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: schedule id must be specified for a progressive rollout",
	)
)
