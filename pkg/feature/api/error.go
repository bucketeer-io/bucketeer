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
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/status"
)

var (
	statusInternal                   = gstatus.New(codes.Internal, "feature: internal")
	statusMissingID                  = gstatus.New(codes.InvalidArgument, "feature: missing id")
	statusMissingIDs                 = gstatus.New(codes.InvalidArgument, "feature: missing ids")
	statusInvalidID                  = gstatus.New(codes.InvalidArgument, "feature: invalid id")
	statusMissingKeyword             = gstatus.New(codes.InvalidArgument, "feature: missing keyword")
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
	statusMissingFeatureTag               = gstatus.New(codes.InvalidArgument, "feature: missing feature tag")
	statusMissingEvaluation               = gstatus.New(codes.InvalidArgument, "feature: missing evaluation")
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
	statusInvalidPrerequisite = gstatus.New(codes.FailedPrecondition, "feature: invalid prerequisite")

	errInternalJaJP = status.MustWithDetails(
		statusInternal,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "????????????????????????????????????",
		},
	)
	errMissingIDJaJP = status.MustWithDetails(
		statusMissingID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "id???????????????",
		},
	)
	errMissingIDsJaJP = status.MustWithDetails(
		statusMissingIDs,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "ids???????????????",
		},
	)
	errInvalidIDJaJP = status.MustWithDetails(
		statusInvalidID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????id??????",
		},
	)
	errMissingKeywordJaJP = status.MustWithDetails(
		statusMissingKeyword,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "keyword???????????????",
		},
	)
	errMissingUserJaJP = status.MustWithDetails(
		statusMissingUser,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "user???????????????",
		},
	)
	errMissingUserIDJaJP = status.MustWithDetails(
		statusMissingUserID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "user id???????????????",
		},
	)
	errMissingUserIDsJaJP = status.MustWithDetails(
		statusMissingUserIDs,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "user id???????????????????????????",
		},
	)
	errMissingCommandJaJP = status.MustWithDetails(
		statusMissingCommand,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "command???????????????",
		},
	)
	errMissingDefaultOnVariationJaJP = status.MustWithDetails(
		statusMissingDefaultOnVariation,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "default variation???????????????",
		},
	)
	errMissingDefaultOffVariationJaJP = status.MustWithDetails(
		statusMissingDefaultOffVariation,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "off variation???????????????",
		},
	)
	errInvalidDefaultOnVariationJaJP = status.MustWithDetails(
		statusInvalidDefaultOnVariation,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????default variation??????",
		},
	)
	errInvalidDefaultOffVariationJaJP = status.MustWithDetails(
		statusInvalidDefaultOffVariation,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????off variation??????",
		},
	)
	errMissingVariationIDJaJP = status.MustWithDetails(
		statusMissingVariationID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "variation id???????????????",
		},
	)
	errInvalidVariationIDJaJP = status.MustWithDetails(
		statusInvalidVariationID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????variation id??????",
		},
	)
	errDifferentVariationsSizeJaJP = status.MustWithDetails(
		statusDifferentVariationsSize,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "feature???variations???rollout???variations????????????????????????",
		},
	)
	errExceededMaxVariationWeightJaJP = status.MustWithDetails(
		statusExceededMaxVariationWeight,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: fmt.Sprintf("?????????weight??????????????????????????? (%d) ?????????????????????", totalVariationWeight),
		},
	)
	errIncorrectVariationWeightJaJP = status.MustWithDetails(
		statusIncorrectVariationWeight,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: fmt.Sprintf("weight???0??????%d????????????????????????????????????", totalVariationWeight),
		},
	)
	errInvalidCursorJaJP = status.MustWithDetails(
		statusInvalidCursor,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????cursor??????",
		},
	)
	errInvalidOrderByJaJP = status.MustWithDetails(
		statusInvalidOrderBy,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "????????????????????????????????????",
		},
	)
	errMissingNameJaJP = status.MustWithDetails(
		statusMissingName,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "name???????????????",
		},
	)
	errMissingFeatureVariationsJaJP = status.MustWithDetails(
		statusMissingFeatureVariations,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "feature???variations???????????????",
		},
	)
	errMissingFeatureTagsJaJP = status.MustWithDetails(
		statusMissingFeatureTags,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "feature???tags???????????????",
		},
	)
	errMissingFeatureTagJaJP = status.MustWithDetails(
		statusMissingFeatureTag,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "feature tag???????????????",
		},
	)
	errMissingEvaluationJaJP = status.MustWithDetails(
		statusMissingEvaluation,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "evaluation???????????????",
		},
	)
	errUnknownCommandJaJP = status.MustWithDetails(
		statusUnknownCommand,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????command??????",
		},
	)
	errMissingRuleJaJP = status.MustWithDetails(
		statusMissingRule,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "rule???????????????",
		},
	)
	errMissingRuleIDJaJP = status.MustWithDetails(
		statusMissingRuleID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "rule id???????????????",
		},
	)
	errMissingRuleClauseJaJP = status.MustWithDetails(
		statusMissingRuleClause,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "rule????????????????????????",
		},
	)
	errMissingClauseIDJaJP = status.MustWithDetails(
		statusMissingClauseID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????id???????????????",
		},
	)
	errMissingClauseAttributeJaJP = status.MustWithDetails(
		statusMissingClauseAttribute,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????attribute???????????????",
		},
	)
	errMissingClauseValuesJaJP = status.MustWithDetails(
		statusMissingClauseValues,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "???????????????????????????????????????",
		},
	)
	errMissingClauseValueJaJP = status.MustWithDetails(
		statusMissingClauseValue,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "???????????????????????????",
		},
	)
	errMissingSegmentIDJaJP = status.MustWithDetails(
		statusMissingSegmentID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "segment id???????????????",
		},
	)
	errMissingSegmentUsersDataJaJP = status.MustWithDetails(
		statusMissingSegmentUsersData,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "segment user????????????????????????????????????",
		},
	)
	errMissingRuleStrategy = status.MustWithDetails(
		statusMissingRuleStrategy,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "rule strategy???????????????",
		},
	)
	errUnknownStrategy = status.MustWithDetails(
		statusUnknownStrategy,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????strategy??????",
		},
	)
	errMissingFixedStrategy = status.MustWithDetails(
		statusMissingFixedStrategy,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "fixed stategy???????????????",
		},
	)
	errMissingRolloutStrategy = status.MustWithDetails(
		statusMissingRolloutStrategy,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "rollout strategy???????????????",
		},
	)
	errExceededMaxSegmentUsersDataSizeJaJP = status.MustWithDetails(
		statusExceededMaxSegmentUsersDataSize,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: fmt.Sprintf("segment user???????????????????????????????????? (%d bytes) ?????????????????????", maxSegmentUsersDataSize),
		},
	)
	errUnknownSegmentUserStateJaJP = status.MustWithDetails(
		statusUnknownSegmentUserState,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????segment user???state??????",
		},
	)
	errIncorrectUUIDFormatJaJP = status.MustWithDetails(
		statusIncorrectUUIDFormat,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????UUID???????????????????????????",
		},
	)
	errExceededMaxUserIDsLengthJaJP = status.MustWithDetails(
		statusExceededMaxUserIDsLength,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: fmt.Sprintf("user id????????????????????? (%d) ?????????????????????", maxUserIDsLength),
		},
	)
	errIncorrectDestinationEnvironmentJaJP = status.MustWithDetails(
		statusIncorrectDestinationEnvironment,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "????????????????????????????????????environment???????????????",
		},
	)
	errExceededMaxPageSizePerRequestJaJP = status.MustWithDetails(
		statusExceededMaxPageSizePerRequest,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: fmt.Sprintf("page size???????????? (%d) ?????????????????????", maxPageSizePerRequest),
		},
	)
	errNotFoundJaJP = status.MustWithDetails(
		statusNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "??????????????????????????????",
		},
	)
	errSegmentNotFoundJaJP = status.MustWithDetails(
		statusSegmentNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "segment?????????????????????",
		},
	)
	errAlreadyExistsJaJP = status.MustWithDetails(
		statusAlreadyExists,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "??????id???????????????????????????????????????",
		},
	)
	errNothingChangeJaJP = status.MustWithDetails(
		statusNothingChange,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "???????????????????????????",
		},
	)
	errSegmentUsersAlreadyUploadingJaJP = status.MustWithDetails(
		statusSegmentUsersAlreadyUploading,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "segment user???????????????????????????????????????????????????",
		},
	)
	errSegmentStatusNotSuceededJaJP = status.MustWithDetails(
		statusSegmentStatusNotSuceeded,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "segment???status???suceeded?????????????????????",
		},
	)
	errSegmentInUseJaJP = status.MustWithDetails(
		statusSegmentInUse,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "segment???feature flag??????????????????????????????????????????????????????",
		},
	)
	errUnauthenticatedJaJP = status.MustWithDetails(
		statusUnauthenticated,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "???????????????????????????",
		},
	)
	errPermissionDeniedJaJP = status.MustWithDetails(
		statusPermissionDenied,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "????????????????????????",
		},
	)
	errWaitingOrRunningExperimentExistsJaJP = status.MustWithDetails(
		statusWaitingOrRunningExperimentExists,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "???????????????????????????????????????Experiment??????????????????????????????????????????Experiment??????????????????????????????",
		},
	)
	errInvalidArchiveJaJP = status.MustWithDetails(
		statusInvalidArchive,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "???????????????????????????????????????????????????????????????????????????????????????????????????????????????",
		},
	)
	errInvalidChangingVariationJaJP = status.MustWithDetails(
		statusInvalidChangingVariation,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????",
		},
	)
	errInvalidPrerequisiteJaJP = status.MustWithDetails(
		statusInvalidPrerequisite,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????prerequisite??????",
		},
	)
)

func localizedError(s *gstatus.Status, loc string) error {
	// handle loc if multi-lang is necessary
	switch s {
	case statusInternal:
		return errInternalJaJP
	case statusMissingID:
		return errMissingIDJaJP
	case statusMissingIDs:
		return errMissingIDsJaJP
	case statusInvalidID:
		return errInvalidIDJaJP
	case statusMissingKeyword:
		return errMissingKeywordJaJP
	case statusMissingUser:
		return errMissingUserJaJP
	case statusMissingUserID:
		return errMissingUserIDJaJP
	case statusMissingUserIDs:
		return errMissingUserIDsJaJP
	case statusMissingCommand:
		return errMissingCommandJaJP
	case statusMissingDefaultOnVariation:
		return errMissingDefaultOnVariationJaJP
	case statusMissingDefaultOffVariation:
		return errMissingDefaultOffVariationJaJP
	case statusInvalidDefaultOnVariation:
		return errInvalidDefaultOnVariationJaJP
	case statusInvalidDefaultOffVariation:
		return errInvalidDefaultOffVariationJaJP
	case statusMissingVariationID:
		return errMissingVariationIDJaJP
	case statusInvalidVariationID:
		return errInvalidVariationIDJaJP
	case statusDifferentVariationsSize:
		return errDifferentVariationsSizeJaJP
	case statusExceededMaxVariationWeight:
		return errExceededMaxVariationWeightJaJP
	case statusIncorrectVariationWeight:
		return errIncorrectVariationWeightJaJP
	case statusInvalidCursor:
		return errInvalidCursorJaJP
	case statusInvalidOrderBy:
		return errInvalidOrderByJaJP
	case statusMissingName:
		return errMissingNameJaJP
	case statusMissingFeatureVariations:
		return errMissingFeatureVariationsJaJP
	case statusMissingFeatureTags:
		return errMissingFeatureTagsJaJP
	case statusMissingFeatureTag:
		return errMissingFeatureTagJaJP
	case statusMissingEvaluation:
		return errMissingEvaluationJaJP
	case statusUnknownCommand:
		return errUnknownCommandJaJP
	case statusMissingRule:
		return errMissingRuleJaJP
	case statusMissingRuleID:
		return errMissingRuleIDJaJP
	case statusMissingRuleClause:
		return errMissingRuleClauseJaJP
	case statusMissingClauseID:
		return errMissingClauseIDJaJP
	case statusMissingClauseAttribute:
		return errMissingClauseAttributeJaJP
	case statusMissingClauseValues:
		return errMissingClauseValuesJaJP
	case statusMissingClauseValue:
		return errMissingClauseValueJaJP
	case statusMissingSegmentID:
		return errMissingSegmentIDJaJP
	case statusMissingSegmentUsersData:
		return errMissingSegmentUsersDataJaJP
	case statusMissingRuleStrategy:
		return errMissingRuleStrategy
	case statusUnknownStrategy:
		return errUnknownStrategy
	case statusMissingFixedStrategy:
		return errMissingFixedStrategy
	case statusMissingRolloutStrategy:
		return errMissingRolloutStrategy
	case statusExceededMaxSegmentUsersDataSize:
		return errExceededMaxSegmentUsersDataSizeJaJP
	case statusUnknownSegmentUserState:
		return errUnknownSegmentUserStateJaJP
	case statusIncorrectUUIDFormat:
		return errIncorrectUUIDFormatJaJP
	case statusExceededMaxUserIDsLength:
		return errExceededMaxUserIDsLengthJaJP
	case statusIncorrectDestinationEnvironment:
		return errIncorrectDestinationEnvironmentJaJP
	case statusExceededMaxPageSizePerRequest:
		return errExceededMaxPageSizePerRequestJaJP
	case statusNotFound:
		return errNotFoundJaJP
	case statusSegmentNotFound:
		return errSegmentNotFoundJaJP
	case statusAlreadyExists:
		return errAlreadyExistsJaJP
	case statusNothingChange:
		return errNothingChangeJaJP
	case statusSegmentUsersAlreadyUploading:
		return errSegmentUsersAlreadyUploadingJaJP
	case statusSegmentStatusNotSuceeded:
		return errSegmentStatusNotSuceededJaJP
	case statusSegmentInUse:
		return errSegmentInUseJaJP
	case statusUnauthenticated:
		return errUnauthenticatedJaJP
	case statusPermissionDenied:
		return errPermissionDeniedJaJP
	case statusWaitingOrRunningExperimentExists:
		return errWaitingOrRunningExperimentExistsJaJP
	case statusInvalidArchive:
		return errInvalidArchiveJaJP
	case statusInvalidChangingVariation:
		return errInvalidChangingVariationJaJP
	case statusInvalidPrerequisite:
		return errInvalidPrerequisiteJaJP
	default:
		return errInternalJaJP
	}
}
