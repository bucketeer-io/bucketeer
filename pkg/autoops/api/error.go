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
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/status"
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
	statusOpsEventRateClauseFeatureVersionRequired = gstatus.New(
		codes.InvalidArgument,
		"autoops: ops event rate clause feature version must be specified",
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
	statusAlreadyExists    = gstatus.New(codes.AlreadyExists, "autoops: already exists")
	statusUnauthenticated  = gstatus.New(codes.Unauthenticated, "autoops: unauthenticated")
	statusPermissionDenied = gstatus.New(codes.PermissionDenied, "autoops: permission denied")
	statusInvalidRequest   = gstatus.New(codes.InvalidArgument, "autoops: invalid request")

	errInternalJaJP = status.MustWithDetails(
		statusInternal,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "内部エラーが発生しました",
		},
	)
	errUnknownOpsTypeJaJP = status.MustWithDetails(
		statusUnknownOpsType,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不明なオペレーションタイプです",
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
	errNoCommandJaJP = status.MustWithDetails(
		statusNoCommand,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "commandは必須です",
		},
	)
	errIDRequiredJaJP = status.MustWithDetails(
		statusIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "idは必須です",
		},
	)
	errFeatureIDRequiredJaJP = status.MustWithDetails(
		statusFeatureIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "feature idは必須です",
		},
	)
	errClauseRequiredJaJP = status.MustWithDetails(
		statusClauseRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "自動オペレーションルールは必須です",
		},
	)
	errClauseIDRequiredJaJP = status.MustWithDetails(
		statusClauseIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "自動オペレーションルールのidは必須です",
		},
	)
	errIncompatibleOpsTypeJaJP = status.MustWithDetails(
		statusIncompatibleOpsType,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "対象のオペレーションタイプに対応していない自動オペレーションルールがあります",
		},
	)
	errOpsEventRateClauseRequiredJaJP = status.MustWithDetails(
		statusOpsEventRateClauseRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "イベントレートルールは必須です",
		},
	)
	errOpsEventRateClauseFeatureVersionRequiredJaJP = status.MustWithDetails(
		statusOpsEventRateClauseFeatureVersionRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "イベントレートルールのfeature versionは必須です",
		},
	)
	errOpsEventRateClauseVariationIDRequiredJaJP = status.MustWithDetails(
		statusOpsEventRateClauseVariationIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "イベントレートルールのvariation idは必須です",
		},
	)
	errOpsEventRateClauseGoalIDRequiredJaJP = status.MustWithDetails(
		statusOpsEventRateClauseGoalIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "イベントレートルールのgoal idは必須です",
		},
	)
	errOpsEventRateClauseMinCountRequiredJaJP = status.MustWithDetails(
		statusOpsEventRateClauseMinCountRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "イベントレートルールのminimum countは必須です",
		},
	)
	errOpsEventRateClauseInvalidThredsholdJaJP = status.MustWithDetails(
		statusOpsEventRateClauseMinCountRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "イベントレートルールのしきい値が不正です",
		},
	)
	errDatetimeClauseRequiredJaJP = status.MustWithDetails(
		statusDatetimeClauseRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "日時ルールは必須です",
		},
	)
	errDatetimeClauseInvalidTimeJaJP = status.MustWithDetails(
		statusDatetimeClauseInvalidTime,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "日時ルールの日時が不正です",
		},
	)
	errWebhookNotFoundJaJP = status.MustWithDetails(
		statusWebhookNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "ウェブフックが存在しません",
		},
	)
	errWebhookClauseRequiredJaJP = status.MustWithDetails(
		statusWebhookClauseRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "ウェブフックルールは必須です",
		},
	)
	errWebhookClauseWebhookIDRequiredJaJP = status.MustWithDetails(
		statusWebhookClauseWebhookIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "ウェブフックルールのwebhook idは必須です",
		},
	)
	errWebhookClauseConditionRequiredJaJP = status.MustWithDetails(
		statusWebhookClauseConditionRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "ウェブフックルールのconditionは必須です",
		},
	)
	errWebhookClauseConditionFilterRequiredJaJP = status.MustWithDetails(
		statusWebhookClauseConditionFilterRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "ウェブフックルールのconditionのfilterは必須です",
		},
	)
	errWebhookClauseConditionValueRequiredJaJP = status.MustWithDetails(
		statusWebhookClauseConditionValueRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "ウェブフックルールのconditionのvalueは必須です",
		},
	)
	errWebhookClauseConditionInvalidOperatorJaJP = status.MustWithDetails(
		statusWebhookClauseConditionInvalidOperator,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "ウェブフックルールのconditionのoperatorが不正です",
		},
	)
	errNotFoundJaJP = status.MustWithDetails(
		statusNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "データが存在しません",
		},
	)
	errAlreadyDeletedJaJP = status.MustWithDetails(
		statusAlreadyDeleted,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "データがすでに削除済みです",
		},
	)
	errOpsEventRateClauseGoalNotFoundJaJP = status.MustWithDetails(
		statusOpsEventRateClauseGoalNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "イベントレートルールのgoalが存在しません",
		},
	)
	errAlreadyExistsJaJP = status.MustWithDetails(
		statusAlreadyExists,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "同じidのデータがすでに存在します",
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
)

func localizedError(s *gstatus.Status, loc string) error {
	// handle loc if multi-lang is necessary
	switch s {
	case statusInternal:
		return errInternalJaJP
	case statusUnknownOpsType:
		return errUnknownOpsTypeJaJP
	case statusInvalidCursor:
		return errInvalidCursorJaJP
	case statusInvalidOrderBy:
		return errInvalidOrderByJaJP
	case statusNoCommand:
		return errNoCommandJaJP
	case statusIDRequired:
		return errIDRequiredJaJP
	case statusFeatureIDRequired:
		return errFeatureIDRequiredJaJP
	case statusClauseRequired:
		return errClauseRequiredJaJP
	case statusClauseIDRequired:
		return errClauseIDRequiredJaJP
	case statusIncompatibleOpsType:
		return errIncompatibleOpsTypeJaJP
	case statusOpsEventRateClauseRequired:
		return errOpsEventRateClauseRequiredJaJP
	case statusOpsEventRateClauseFeatureVersionRequired:
		return errOpsEventRateClauseFeatureVersionRequiredJaJP
	case statusOpsEventRateClauseVariationIDRequired:
		return errOpsEventRateClauseVariationIDRequiredJaJP
	case statusOpsEventRateClauseGoalIDRequired:
		return errOpsEventRateClauseGoalIDRequiredJaJP
	case statusOpsEventRateClauseMinCountRequired:
		return errOpsEventRateClauseMinCountRequiredJaJP
	case statusOpsEventRateClauseMinCountRequired:
		return errOpsEventRateClauseInvalidThredsholdJaJP
	case statusDatetimeClauseRequired:
		return errDatetimeClauseRequiredJaJP
	case statusDatetimeClauseInvalidTime:
		return errDatetimeClauseInvalidTimeJaJP
	case statusWebhookNotFound:
		return errWebhookNotFoundJaJP
	case statusWebhookClauseRequired:
		return errWebhookClauseRequiredJaJP
	case statusWebhookClauseWebhookIDRequired:
		return errWebhookClauseWebhookIDRequiredJaJP
	case statusWebhookClauseConditionRequired:
		return errWebhookClauseConditionRequiredJaJP
	case statusWebhookClauseConditionFilterRequired:
		return errWebhookClauseConditionFilterRequiredJaJP
	case statusWebhookClauseConditionInvalidOperator:
		return errWebhookClauseConditionInvalidOperatorJaJP
	case statusNotFound:
		return errNotFoundJaJP
	case statusAlreadyDeleted:
		return errAlreadyDeletedJaJP
	case statusOpsEventRateClauseGoalNotFound:
		return errOpsEventRateClauseGoalNotFoundJaJP
	case statusAlreadyExists:
		return errAlreadyExistsJaJP
	case statusUnauthenticated:
		return errUnauthenticatedJaJP
	case statusPermissionDenied:
		return errPermissionDeniedJaJP
	default:
		return errInternalJaJP
	}
}
