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
			Message: "内部エラーが発生しました",
		},
	)
	errMissingIDJaJP = status.MustWithDetails(
		statusMissingID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "idは必須です",
		},
	)
	errMissingIDsJaJP = status.MustWithDetails(
		statusMissingIDs,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "idsは必須です",
		},
	)
	errInvalidIDJaJP = status.MustWithDetails(
		statusInvalidID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なidです",
		},
	)
	errMissingKeywordJaJP = status.MustWithDetails(
		statusMissingKeyword,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "keywordは必須です",
		},
	)
	errMissingUserJaJP = status.MustWithDetails(
		statusMissingUser,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "userは必須です",
		},
	)
	errMissingUserIDJaJP = status.MustWithDetails(
		statusMissingUserID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "user idは必須です",
		},
	)
	errMissingUserIDsJaJP = status.MustWithDetails(
		statusMissingUserIDs,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "user idのリストは必須です",
		},
	)
	errMissingCommandJaJP = status.MustWithDetails(
		statusMissingCommand,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "commandは必須です",
		},
	)
	errMissingDefaultOnVariationJaJP = status.MustWithDetails(
		statusMissingDefaultOnVariation,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "default variationは必須です",
		},
	)
	errMissingDefaultOffVariationJaJP = status.MustWithDetails(
		statusMissingDefaultOffVariation,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "off variationは必須です",
		},
	)
	errInvalidDefaultOnVariationJaJP = status.MustWithDetails(
		statusInvalidDefaultOnVariation,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なdefault variationです",
		},
	)
	errInvalidDefaultOffVariationJaJP = status.MustWithDetails(
		statusInvalidDefaultOffVariation,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なoff variationです",
		},
	)
	errMissingVariationIDJaJP = status.MustWithDetails(
		statusMissingVariationID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "variation idは必須です",
		},
	)
	errInvalidVariationIDJaJP = status.MustWithDetails(
		statusInvalidVariationID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なvariation idです",
		},
	)
	errDifferentVariationsSizeJaJP = status.MustWithDetails(
		statusDifferentVariationsSize,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "featureのvariationsとrolloutのvariationsの数が異なります",
		},
	)
	errExceededMaxVariationWeightJaJP = status.MustWithDetails(
		statusExceededMaxVariationWeight,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: fmt.Sprintf("全てのweightの合計の最大サイズ (%d) を超えています", totalVariationWeight),
		},
	)
	errIncorrectVariationWeightJaJP = status.MustWithDetails(
		statusIncorrectVariationWeight,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: fmt.Sprintf("weightは0から%dの間である必要があります", totalVariationWeight),
		},
	)
	errInvalidCursorJaJP = status.MustWithDetails(
		statusInvalidCursor,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なcursorです",
		},
	)
	errInvalidOrderByJaJP = status.MustWithDetails(
		statusInvalidOrderBy,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なソート順の指定です",
		},
	)
	errMissingNameJaJP = status.MustWithDetails(
		statusMissingName,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "nameは必須です",
		},
	)
	errMissingFeatureVariationsJaJP = status.MustWithDetails(
		statusMissingFeatureVariations,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "featureのvariationsは必須です",
		},
	)
	errMissingFeatureTagsJaJP = status.MustWithDetails(
		statusMissingFeatureTags,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "featureのtagsは必須です",
		},
	)
	errMissingFeatureTagJaJP = status.MustWithDetails(
		statusMissingFeatureTag,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "feature tagは必須です",
		},
	)
	errMissingEvaluationJaJP = status.MustWithDetails(
		statusMissingEvaluation,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "evaluationは必須です",
		},
	)
	errUnknownCommandJaJP = status.MustWithDetails(
		statusUnknownCommand,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不明なcommandです",
		},
	)
	errMissingRuleJaJP = status.MustWithDetails(
		statusMissingRule,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "ruleは必須です",
		},
	)
	errMissingRuleIDJaJP = status.MustWithDetails(
		statusMissingRuleID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "rule idは必須です",
		},
	)
	errMissingRuleClauseJaJP = status.MustWithDetails(
		statusMissingRuleClause,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "ruleの条件は必須です",
		},
	)
	errMissingClauseIDJaJP = status.MustWithDetails(
		statusMissingClauseID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "条件のidは必須です",
		},
	)
	errMissingClauseAttributeJaJP = status.MustWithDetails(
		statusMissingClauseAttribute,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "条件のattributeは必須です",
		},
	)
	errMissingClauseValuesJaJP = status.MustWithDetails(
		statusMissingClauseValues,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "条件の値のリストは必須です",
		},
	)
	errMissingClauseValueJaJP = status.MustWithDetails(
		statusMissingClauseValue,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "条件の値は必須です",
		},
	)
	errMissingSegmentIDJaJP = status.MustWithDetails(
		statusMissingSegmentID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "segment idは必須です",
		},
	)
	errMissingSegmentUsersDataJaJP = status.MustWithDetails(
		statusMissingSegmentUsersData,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "segment userリストのデータは必須です",
		},
	)
	errMissingRuleStrategy = status.MustWithDetails(
		statusMissingRuleStrategy,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "rule strategyは必須です",
		},
	)
	errUnknownStrategy = status.MustWithDetails(
		statusUnknownStrategy,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不明なstrategyです",
		},
	)
	errMissingFixedStrategy = status.MustWithDetails(
		statusMissingFixedStrategy,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "fixed stategyは必須です",
		},
	)
	errMissingRolloutStrategy = status.MustWithDetails(
		statusMissingRolloutStrategy,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "rollout strategyは必須です",
		},
	)
	errExceededMaxSegmentUsersDataSizeJaJP = status.MustWithDetails(
		statusExceededMaxSegmentUsersDataSize,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: fmt.Sprintf("segment userリストの最大データサイズ (%d bytes) を超えています", maxSegmentUsersDataSize),
		},
	)
	errUnknownSegmentUserStateJaJP = status.MustWithDetails(
		statusUnknownSegmentUserState,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不明なsegment userのstateです",
		},
	)
	errIncorrectUUIDFormatJaJP = status.MustWithDetails(
		statusIncorrectUUIDFormat,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なUUIDのフォーマットです",
		},
	)
	errExceededMaxUserIDsLengthJaJP = status.MustWithDetails(
		statusExceededMaxUserIDsLength,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: fmt.Sprintf("user idリストの最大数 (%d) を超えています", maxUserIDsLength),
		},
	)
	errIncorrectDestinationEnvironmentJaJP = status.MustWithDetails(
		statusIncorrectDestinationEnvironment,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "クローン元とクローン先のenvironmentが同じです",
		},
	)
	errExceededMaxPageSizePerRequestJaJP = status.MustWithDetails(
		statusExceededMaxPageSizePerRequest,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: fmt.Sprintf("page sizeの最大値 (%d) を超えています", maxPageSizePerRequest),
		},
	)
	errNotFoundJaJP = status.MustWithDetails(
		statusNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "データが存在しません",
		},
	)
	errSegmentNotFoundJaJP = status.MustWithDetails(
		statusSegmentNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "segmentが存在しません",
		},
	)
	errAlreadyExistsJaJP = status.MustWithDetails(
		statusAlreadyExists,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "同じidのデータがすでに存在します",
		},
	)
	errNothingChangeJaJP = status.MustWithDetails(
		statusNothingChange,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "変更点がありません",
		},
	)
	errSegmentUsersAlreadyUploadingJaJP = status.MustWithDetails(
		statusSegmentUsersAlreadyUploading,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "segment userのリストはすでにアップロード中です",
		},
	)
	errSegmentStatusNotSuceededJaJP = status.MustWithDetails(
		statusSegmentStatusNotSuceeded,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "segmentのstatusがsuceededではありません",
		},
	)
	errSegmentInUseJaJP = status.MustWithDetails(
		statusSegmentInUse,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "segmentがfeature flagで使用されているため、削除できません",
		},
	)
	errUnauthenticatedJaJP = status.MustWithDetails(
		statusUnauthenticated,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "認証されていません",
		},
	)
	errPermissionDeniedJaJP = status.MustWithDetails(
		statusPermissionDenied,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "権限がありません",
		},
	)
	errWaitingOrRunningExperimentExistsJaJP = status.MustWithDetails(
		statusWaitingOrRunningExperimentExists,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "開始予定、もしくは実行中のExperimentが存在します。更新する場合はExperimentを停止してください。",
		},
	)
	errInvalidArchiveJaJP = status.MustWithDetails(
		statusInvalidArchive,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "前提条件のフラグとして登録されているフラグをアーカイブすることはできません",
		},
	)
	errInvalidChangingVariationJaJP = status.MustWithDetails(
		statusInvalidChangingVariation,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "前提条件のフラグとして登録されているフラグのバリエーションを変更または削除することはできません",
		},
	)
	errInvalidPrerequisiteJaJP = status.MustWithDetails(
		statusInvalidPrerequisite,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なprerequisiteです",
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
