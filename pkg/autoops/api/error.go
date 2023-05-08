// Copyright 2022 The Bucketeer Authors.
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
	statusInvalidOrderBy    = gstatus.New(codes.InvalidArgument, "autoops: order_by is invalid")
	statusNoCommand         = gstatus.New(codes.InvalidArgument, "autoops: no command")
	statusIDRequired        = gstatus.New(codes.InvalidArgument, "autoops: id must be specified")
	statusFeatureIDRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: featureId must be specified",
	)
	statusClauseRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: at least one clause must be specified",
	)
	statusClauseIDRequired    = gstatus.New(codes.InvalidArgument, "autoops: clause id must be specified")
	statusIncompatibleOpsType = gstatus.New(
		codes.InvalidArgument,
		"autoops: ops type is incompatible with ops event rate clause",
	)
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
	statusNotFound                       = gstatus.New(codes.NotFound, "autoops: not found")
	statusAlreadyDeleted                 = gstatus.New(codes.NotFound, "autoops: already deleted")
	statusOpsEventRateClauseGoalNotFound = gstatus.New(
		codes.NotFound,
		"autoops: ops event rate clause goal does not exist",
	)
	statusWebhookNotFound       = gstatus.New(codes.NotFound, "autoops: webhook not found")
	statusWebhookClauseRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: webhook clause is required",
	)
	statusWebhookClauseWebhookIDRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: webhook clause wehook id is required",
	)
	statusWebhookClauseConditionRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: webhook clause condition is required",
	)
	statusWebhookClauseConditionFilterRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: webhook clause condition filter is required",
	)
	statusWebhookClauseConditionValueRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: webhook clause condition value is required",
	)
	statusWebhookClauseConditionInvalidOperator = gstatus.New(
		codes.InvalidArgument,
		"autoops: webhook clause condition oerator is invalid",
	)
	statusAlreadyExists               = gstatus.New(codes.AlreadyExists, "autoops: already exists")
	statusUnauthenticated             = gstatus.New(codes.Unauthenticated, "autoops: unauthenticated")
	statusPermissionDenied            = gstatus.New(codes.PermissionDenied, "autoops: permission denied")
	statusInvalidRequest              = gstatus.New(codes.InvalidArgument, "autoops: invalid request")
	statusProgressiveRolloutNoCommand = gstatus.New(
		codes.InvalidArgument,
		"autoops progressive rollout: no command",
	)
	statusProgressiveRolloutFeatureIDRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops progressive rollout: featureId must be specified",
	)
	statusProgressiveRolloutClauseRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops progressive rollout: at least one clause must be specified",
	)
	statusProgressiveRolloutInternal = gstatus.New(
		codes.Internal,
		"autoops progressive rollout: internal",
	)
	statusProgressiveRolloutClauseVariationIDRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops progressive rollout: clause variation id must be specified",
	)
	statusProgressiveRolloutClauseSchedulesRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops progressive rollout: at least one clause schedule must be specified",
	)
	statusProgressiveRolloutClauseInvalidIncrements = gstatus.New(
		codes.InvalidArgument,
		"autoops progressive rollout: clause increments is invalid",
	)
	statusProgressiveRolloutClauseUnknownInterval = gstatus.New(
		codes.InvalidArgument,
		"autoops progressive rollout: clause interval is unknown.",
	)
	statusProgressiveRolloutInvalidVariationSize = gstatus.New(
		codes.FailedPrecondition,
		"autoops progressive rollout: the number of variations must be equal to 2",
	)
	statusProgressiveRolloutInvalidScheduleSpans = gstatus.New(
		codes.FailedPrecondition,
		"autoops progressive rollout: The span of time for each scheduled time must be at least 5 minutes.",
	)
	statusProgressiveRolloutScheduleExecutedAtRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops progressive rollout: schedule executed_at must be specified",
	)
	statusProgressiveRolloutScheduleInvalidWeight = gstatus.New(
		codes.InvalidArgument,
		"autoops progressive rollout: schedule weight is invalid",
	)
	statusProgressiveRolloutAlreadyExists = gstatus.New(
		codes.AlreadyExists,
		"autoops progressive rollout: already exists",
	)
	statusProgressiveRolloutIDRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops progressive rollout: id must be specified",
	)
	statusProgressiveRolloutNotFound = gstatus.New(
		codes.NotFound,
		"autoops progressive rollout: not found",
	)
	statusProgressiveRolloutInvalidCursor = gstatus.New(
		codes.InvalidArgument,
		"autoops progressive rollout: cursor is invalid",
	)
	statusProgressiveRolloutInvalidOrderBy = gstatus.New(
		codes.InvalidArgument,
		"autoops progressive rollout: order_by is invalid",
	)
	statusProgressiveRolloutScheduleIDRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops progressive rollout: schedule id must be specified",
	)
)
