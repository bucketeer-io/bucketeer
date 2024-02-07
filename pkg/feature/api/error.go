// Copyright 2024 The Bucketeer Authors.
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
	"fmt"

	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

var (
	statusInternal                   = gstatus.New(codes.Internal, "feature: internal")
	statusMissingFrom                = gstatus.New(codes.InvalidArgument, "feature: missing from")
	statusMissingID                  = gstatus.New(codes.InvalidArgument, "feature: missing id")
	statusMissingIDs                 = gstatus.New(codes.InvalidArgument, "feature: missing ids")
	statusInvalidID                  = gstatus.New(codes.InvalidArgument, "feature: invalid id")
	statusMissingUser                = gstatus.New(codes.InvalidArgument, "feature: missing user")
	statusMissingUserID              = gstatus.New(codes.InvalidArgument, "feature: missing user id")
	statusMissingUserIDs             = gstatus.New(codes.InvalidArgument, "feature: missing user ids")
	statusMissingCommand             = gstatus.New(codes.InvalidArgument, "feature: missing command")
	statusMissingDefaultOnVariation  = gstatus.New(codes.InvalidArgument, "feature: missing default on variation")
	statusMissingDefaultOffVariation = gstatus.New(codes.InvalidArgument, "feature: missing default off variation")
	statusInvalidDefaultOnVariation  = gstatus.New(codes.InvalidArgument, "feature: invalid default on variation")
	statusInvalidDefaultOffVariation = gstatus.New(codes.InvalidArgument, "feature: invalid default off variation")
	statusMissingVariationID         = gstatus.New(codes.InvalidArgument, "feature: missing variation id")
	statusInvalidVariationID         = gstatus.New(codes.InvalidArgument, "feature: invalid variation id")
	statusDifferentVariationsSize    = gstatus.New(
		codes.InvalidArgument,
		"feature: feature variations and rollout variations must have the same size",
	)
	statusExceededMaxVariationWeight = gstatus.New(
		codes.InvalidArgument,
		fmt.Sprintf("feature: the sum of all weights value is %d", totalVariationWeight),
	)
	statusIncorrectVariationWeight = gstatus.New(
		codes.InvalidArgument,
		fmt.Sprintf("command: variation weight must be between 0 and %d", totalVariationWeight),
	)
	statusInvalidCursor            = gstatus.New(codes.InvalidArgument, "feature: cursor is invalid")
	statusInvalidOrderBy           = gstatus.New(codes.InvalidArgument, "feature: order_by is invalid")
	statusMissingName              = gstatus.New(codes.InvalidArgument, "feature: missing name")
	statusMissingFeatureVariations = gstatus.New(
		codes.InvalidArgument,
		"feature: feature must contain more than one variation",
	)
	statusMissingFeatureTags = gstatus.New(
		codes.InvalidArgument,
		"feature: feature must contain one or more tags",
	)
	statusUnknownCommand                  = gstatus.New(codes.InvalidArgument, "feature: unknown command")
	statusMissingRule                     = gstatus.New(codes.InvalidArgument, "feature: missing rule")
	statusMissingRuleID                   = gstatus.New(codes.InvalidArgument, "feature: missing rule id")
	statusMissingRuleClause               = gstatus.New(codes.InvalidArgument, "feature: missing rule clause")
	statusMissingClauseID                 = gstatus.New(codes.InvalidArgument, "feature: missing clause id")
	statusMissingClauseAttribute          = gstatus.New(codes.InvalidArgument, "feature: missing clause attribute")
	statusMissingClauseValues             = gstatus.New(codes.InvalidArgument, "feature: missing clause values")
	statusMissingClauseValue              = gstatus.New(codes.InvalidArgument, "feature: missing clause value")
	statusMissingSegmentID                = gstatus.New(codes.InvalidArgument, "feature: missing segment id")
	statusMissingSegmentUsersData         = gstatus.New(codes.InvalidArgument, "feature: missing segment users data")
	statusMissingRuleStrategy             = gstatus.New(codes.InvalidArgument, "feature: missing rule strategy")
	statusUnknownStrategy                 = gstatus.New(codes.InvalidArgument, "feature: unknown strategy")
	statusMissingFixedStrategy            = gstatus.New(codes.InvalidArgument, "feature: missing fixed strategy")
	statusMissingRolloutStrategy          = gstatus.New(codes.InvalidArgument, "feature: missing rollout strategy")
	statusExceededMaxSegmentUsersDataSize = gstatus.New(
		codes.InvalidArgument,
		fmt.Sprintf("feature: max segment users data size allowed is %d bytes", maxSegmentUsersDataSize),
	)
	statusUnknownSegmentUserState = gstatus.New(codes.InvalidArgument, "feature: unknown segment user state")
	statusIncorrectUUIDFormat     = gstatus.New(
		codes.InvalidArgument,
		"feature: uuid format must be an uuid version 4",
	)
	statusExceededMaxUserIDsLength = gstatus.New(
		codes.InvalidArgument,
		fmt.Sprintf("feature: max user ids length allowed is %d", maxUserIDsLength),
	)
	statusIncorrectDestinationEnvironment = gstatus.New(
		codes.InvalidArgument,
		"feature: destination environment is the same as origin one",
	)
	statusExceededMaxPageSizePerRequest = gstatus.New(
		codes.InvalidArgument,
		fmt.Sprintf("feature: max page size allowed is %d", maxPageSizePerRequest),
	)
	statusNotFound                     = gstatus.New(codes.NotFound, "feature: not found")
	statusSegmentNotFound              = gstatus.New(codes.NotFound, "feature: segment not found")
	statusAlreadyExists                = gstatus.New(codes.AlreadyExists, "feature: already exists")
	statusNothingChange                = gstatus.New(codes.FailedPrecondition, "feature: no change")
	statusSegmentUsersAlreadyUploading = gstatus.New(
		codes.FailedPrecondition,
		"feature: segment users already uploading",
	)
	statusSegmentStatusNotSuceeded = gstatus.New(
		codes.FailedPrecondition,
		"feature: segment status is not suceeded",
	)
	statusSegmentInUse                     = gstatus.New(codes.FailedPrecondition, "feature: segment is in use")
	statusUnauthenticated                  = gstatus.New(codes.Unauthenticated, "feature: unauthenticated")
	statusPermissionDenied                 = gstatus.New(codes.PermissionDenied, "feature: permission denied")
	statusWaitingOrRunningExperimentExists = gstatus.New(
		codes.FailedPrecondition,
		"feature: experiment in waiting or running status exists",
	)
	statusCycleExists    = gstatus.New(codes.FailedPrecondition, "feature: circular dependency detected")
	statusInvalidArchive = gstatus.New(
		codes.FailedPrecondition,
		"feature: cant't archive because this feature is used as a prerequsite",
	)
	statusInvalidChangingVariation = gstatus.New(
		codes.FailedPrecondition,
		"feature: can't change or remove this variation because it is used as a prerequsite",
	)
	statusInvalidPrerequisite                     = gstatus.New(codes.FailedPrecondition, "feature: invalid prerequisite")
	statusProgressiveRolloutWaitingOrRunningState = gstatus.New(
		codes.FailedPrecondition,
		"feature: there is a progressive rollout in the waiting or running state",
	)
	// flag trigger
	statusMissingTriggerFeatureID   = gstatus.New(codes.InvalidArgument, "feature: missing trigger feature id")
	statusMissingTriggerType        = gstatus.New(codes.InvalidArgument, "feature: missing trigger type")
	statusMissingTriggerAction      = gstatus.New(codes.InvalidArgument, "feature: missing trigger action")
	statusMissingTriggerDescription = gstatus.New(codes.InvalidArgument, "feature: missing trigger description")
	statusMissingTriggerID          = gstatus.New(codes.InvalidArgument, "feature: missing trigger id")
	statusSecretRequired            = gstatus.New(codes.InvalidArgument, "feature: trigger secret is required")
	statusTriggerAlreadyDisabled    = gstatus.New(codes.FailedPrecondition, "feature: trigger already disabled")
	statusTriggerNotFound           = gstatus.New(codes.NotFound, "feature: trigger not found")
	statusTriggerDisableFailed      = gstatus.New(codes.Internal, "feature: trigger disable failed")
	statusTriggerEnableFailed       = gstatus.New(codes.Internal, "feature: trigger enable failed")
	statusTriggerActionInvalid      = gstatus.New(codes.InvalidArgument, "feature: trigger action is invalid")
	statusTriggerUsageUpdateFailed  = gstatus.New(codes.Internal, "feature: trigger usage update failed")
)
