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
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
)

var (
	statusInternal = api.NewGRPCStatus(
		pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "internal"))
	statusMissingFrom = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgUnknown(pkgErr.FeaturePackageName, "missing from", "From"))
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
	statusMissingUserIDs = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing user ids", "UserId"))
	statusMissingFeatureIDs = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing feature ids", "FeatureFlagID"))
	statusMissingCommand = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "missing command", "Command"))
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
	statusMissingVariationID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing variation id", "VariationId"))
	statusInvalidVariationID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgUnknown(pkgErr.FeaturePackageName, "invalid variation id", "VariationId"))
	statusDifferentVariationsSize = api.NewGRPCStatus(
		pkgErr.NewErrorDifferentVariationsSize(
			pkgErr.FeaturePackageName,
			"feature variations and rollout variations must have the same size",
		))
	statusExceededMaxVariationWeight = api.NewGRPCStatus(
		pkgErr.NewErrorExceededMax(
			pkgErr.FeaturePackageName,
			"the sum of all weights value exceeded",
			"SumOfWeights",
			int(totalVariationWeight),
		))
	statusIncorrectVariationWeight = api.NewGRPCStatus(
		pkgErr.NewErrorOutOfRange(
			pkgErr.FeaturePackageName,
			"variation weight must be between 0 and 100",
			"VariationWeight",
			0,
			100,
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
	statusUnknownCommand = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgUnknown(pkgErr.FeaturePackageName, "unknown command", "Command"))
	statusCommentRequiredForUpdating = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "a comment is required for updating"))
	statusMissingRule = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.FeaturePackageName, "missing rule", "Rule"))
	statusMissingRuleID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing rule id", "RuleId"))
	statusMissingRuleClause = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing rule clause", "RuleClause"))
	statusMissingClauseID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing clause id", "ClauseId"))
	statusMissingClauseAttribute = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing clause attribute", "RuleAttribute"))
	statusMissingClauseValues = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing clause values", "ClauseValue"))
	statusMissingClauseValue = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing clause value", "ClauseValue"))
	statusMissingSegmentID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing segment id", "Segment"))
	statusMissingSegmentUsersData = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.FeaturePackageName, "missing segment users data", "SegmentUser"))
	statusMissingRuleStrategy = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.FeaturePackageName, "missing rule strategy", "RuleStrategy"))
	statusUnknownStrategy = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgUnknown(pkgErr.FeaturePackageName, "unknown strategy", "Strategy"))
	statusMissingFixedStrategy = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.FeaturePackageName, "missing fixed strategy", "FixedStrategy"))
	statusMissingRolloutStrategy = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.FeaturePackageName, "missing rollout strategy", "RolloutStrategy"))
	statusExceededMaxSegmentUsersDataSize = api.NewGRPCStatus(
		pkgErr.NewErrorExceededMax(
			pkgErr.FeaturePackageName,
			"max segment users data size exceeded",
			"SegmentUserData",
			maxSegmentUsersDataSize,
		))
	statusUnknownSegmentUserState = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgUnknown(pkgErr.FeaturePackageName, "unknown segment user state", "SegmentUserState"))
	statusIncorrectUUIDFormat = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "uuid format must be an uuid version 4", "UUID"))
	statusExceededMaxUserIDsLength = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.FeaturePackageName, "max user ids length allowed is %d", "UserId"))
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
	statusCycleExists = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "circular dependency detected"))
	statusInvalidArchive = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(
			pkgErr.FeaturePackageName,
			"cant't archive because this feature is used as a prerequsite"))
	statusInvalidChangingVariation = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(
			pkgErr.FeaturePackageName,
			"can't change or remove this variation because it is used as a prerequsite"))
	statusVariationInUseByOtherFeatures = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(
			pkgErr.FeaturePackageName,
			"can't remove this variation because it is used as a prerequisite or rule in other features",
		))
	statusInvalidPrerequisite = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(pkgErr.FeaturePackageName, "invalid prerequisite"))
	statusProgressiveRolloutWaitingOrRunningState = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(
			pkgErr.FeaturePackageName,
			"there is a progressive rollout in the waiting or running state",
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
)
