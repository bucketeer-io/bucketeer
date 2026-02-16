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
		pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "internal"))
	statusMissingID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing id", "ID"))
	statusMissingIDs = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing ids", "ID"))
	statusInvalidID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "invalid id", "ID"))
	statusMissingUser = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.FeaturePackageName, "missing user", "User"))
	statusMissingUserID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing user id", "UserId"))
	statusMissingFeatureIDs = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing feature ids", "FeatureFlagID"))
	statusMissingDefaultOnVariation = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.FeaturePackageName, "missing default on variation", "DefaultOnVariation"))
	statusMissingDefaultOffVariation = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.FeaturePackageName, "missing default off variation", "DefaultOffVariation"))
	statusInvalidDefaultOnVariation = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.FeaturePackageName,
			"invalid default on variation",
			"DefaultOnVariation",
		))
	statusInvalidDefaultOffVariation = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.FeaturePackageName,
			"invalid default off variation",
			"DefaultOffVariation",
		))
	statusInvalidCursor = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "cursor is invalid", "Cursor"))
	statusInvalidOrderBy = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "order_by is invalid", "OrderBy"))
	statusMissingName = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing name", "Name"))
	statusMissingFeatureVariations = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.FeaturePackageName,
			"feature must contain more than one variation",
			"Variation",
		))
	statusMissingFeatureTags = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "feature must contain one or more tags", "Tag"))
	statusCommentRequiredForUpdating = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "a comment is required for updating"))
	statusMissingSegmentID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing segment id", "Segment"))
	statusMissingSegmentUsersData = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing segment users data", "SegmentUser"))
	statusExceededMaxSegmentUsersDataSize = api.NewGRPCStatus(
		pkgErr.NewErrorExceededMax(
			pkgErr.FeaturePackageName,
			"max segment users data size exceeded",
			"SegmentUserData",
			maxSegmentUsersDataSize,
		))
	statusUnknownSegmentUserState = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgUnknown(pkgErr.FeaturePackageName, "unknown segment user state", "SegmentUserState"))
	statusIncorrectDestinationEnvironment = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.FeaturePackageName,
			"destination environment is the same as origin one",
			"DestinationEnvironment",
		))
	statusExceededMaxPageSizePerRequest = api.NewGRPCStatus(
		pkgErr.NewErrorExceededMax(
			pkgErr.FeaturePackageName,
			"max page size allowed is exceeded",
			"PageSize",
			maxPageSizePerRequest,
		))
	statusFeatureNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.FeaturePackageName, "not found", "FeatureFlag"))
	statusSegmentNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.FeaturePackageName, "segment not found", "Segment"))
	statusAlreadyExists = api.NewGRPCStatus(
		pkgErr.NewErrorAlreadyExists(pkgErr.FeaturePackageName, "already exists"))
	statusSegmentUsersAlreadyUploading = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "segment users already uploading"))
	statusSegmentStatusNotSuceeded = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "segment status is not suceeded"))
	statusSegmentInUse = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "segment is in use"))
	statusUnauthenticated = api.NewGRPCStatus(
		pkgErr.NewErrorUnauthenticated(pkgErr.FeaturePackageName, "unauthenticated"))
	statusPermissionDenied = api.NewGRPCStatus(
		pkgErr.NewErrorPermissionDenied(pkgErr.FeaturePackageName, "permission denied"))
	statusWaitingOrRunningExperimentExists = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "experiment in waiting or running status exists"))
	statusInvalidArchive = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(
			pkgErr.FeaturePackageName,
			"can't archive because this feature is used as a prerequsite"))
	statusVariationInUseByOtherFeatures = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(
			pkgErr.FeaturePackageName,
			"can't remove this variation because it is used as a prerequisite or rule in other features",
		))
	// flag trigger
	statusMissingTriggerFeatureID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing trigger feature id", "FeatureFlagID"))
	statusMissingTriggerType = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing trigger type", "TriggerType"))
	statusMissingTriggerAction = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing trigger action", "TriggerAction"))
	statusMissingTriggerID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing trigger id", "TriggerId"))
	statusSecretRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "trigger secret is required", "TriggerSecret"))
	statusTriggerAlreadyDisabled = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "trigger already disabled"))
	statusTriggerNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.FeaturePackageName, "trigger not found", "FlagTrigger"))
	statusTriggerDisableFailed = api.NewGRPCStatus(
		pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "trigger disable failed"))
	statusTriggerEnableFailed = api.NewGRPCStatus(
		pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "trigger enable failed"))
	statusTriggerActionInvalid = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgUnknown(pkgErr.FeaturePackageName, "trigger action is invalid", "TriggerAction"))
	statusTriggerUsageUpdateFailed = api.NewGRPCStatus(
		pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "trigger usage update failed"))
	// scheduled flag change
	statusMissingScheduledFlagChangeID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.FeaturePackageName,
			"missing scheduled flag change id",
			"ScheduledFlagChangeId",
		))
	statusMissingFeatureID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing feature id", "FeatureFlagID"))
	statusMissingScheduledAt = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing scheduled at", "ScheduledAt"))
	statusMissingPayload = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.FeaturePackageName, "missing payload", "Payload"))
	statusScheduledTimeTooSoon = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.FeaturePackageName,
			"scheduled time must be at least 5 minutes in the future",
			"ScheduledAt",
		))
	statusScheduledTimeTooFar = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.FeaturePackageName,
			"scheduled time must be at most 1 year in the future",
			"ScheduledAt",
		))
	statusExceededMaxSchedulesPerFlag = api.NewGRPCStatus(
		pkgErr.NewErrorExceededMax(
			pkgErr.FeaturePackageName,
			"exceeded maximum number of schedules per flag",
			"SchedulesPerFlag",
			maxSchedulesPerFlag,
		))
	statusExceededMaxChangesPerSchedule = api.NewGRPCStatus(
		pkgErr.NewErrorExceededMax(
			pkgErr.FeaturePackageName,
			"exceeded maximum number of changes per schedule",
			"ChangesPerSchedule",
			maxChangesPerSchedule,
		))
	statusScheduledFlagChangeNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.FeaturePackageName, "scheduled flag change not found", "ScheduledFlagChange"))
	statusScheduledFlagChangeNotPending = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "scheduled flag change is not pending"))
	statusEmptyPayload = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "payload must contain at least one change", "Payload"))
	statusInvalidVariationReference = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.FeaturePackageName,
			"invalid variation reference in payload",
			"VariationId",
		))
	statusInvalidRuleReference = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.FeaturePackageName,
			"invalid rule reference in payload",
			"RuleId",
		))
)
