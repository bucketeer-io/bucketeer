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
			Message: "????????????????????????????????????",
		},
	)
	errUnknownOpsTypeJaJP = status.MustWithDetails(
		statusUnknownOpsType,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????????????????????????????????????????",
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
	errNoCommandJaJP = status.MustWithDetails(
		statusNoCommand,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "command???????????????",
		},
	)
	errIDRequiredJaJP = status.MustWithDetails(
		statusIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "id???????????????",
		},
	)
	errFeatureIDRequiredJaJP = status.MustWithDetails(
		statusFeatureIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "feature id???????????????",
		},
	)
	errClauseRequiredJaJP = status.MustWithDetails(
		statusClauseRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "???????????????????????????????????????????????????",
		},
	)
	errClauseIDRequiredJaJP = status.MustWithDetails(
		statusClauseIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "???????????????????????????????????????id???????????????",
		},
	)
	errIncompatibleOpsTypeJaJP = status.MustWithDetails(
		statusIncompatibleOpsType,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "??????????????????????????????????????????????????????????????????????????????????????????????????????????????????",
		},
	)
	errOpsEventRateClauseRequiredJaJP = status.MustWithDetails(
		statusOpsEventRateClauseRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????????????????????????????????????????",
		},
	)
	errOpsEventRateClauseFeatureVersionRequiredJaJP = status.MustWithDetails(
		statusOpsEventRateClauseFeatureVersionRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????????????????????????????feature version???????????????",
		},
	)
	errOpsEventRateClauseVariationIDRequiredJaJP = status.MustWithDetails(
		statusOpsEventRateClauseVariationIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????????????????????????????variation id???????????????",
		},
	)
	errOpsEventRateClauseGoalIDRequiredJaJP = status.MustWithDetails(
		statusOpsEventRateClauseGoalIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????????????????????????????goal id???????????????",
		},
	)
	errOpsEventRateClauseMinCountRequiredJaJP = status.MustWithDetails(
		statusOpsEventRateClauseMinCountRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????????????????????????????minimum count???????????????",
		},
	)
	errOpsEventRateClauseInvalidThredsholdJaJP = status.MustWithDetails(
		statusOpsEventRateClauseMinCountRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "????????????????????????????????????????????????????????????",
		},
	)
	errDatetimeClauseRequiredJaJP = status.MustWithDetails(
		statusDatetimeClauseRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "??????????????????????????????",
		},
	)
	errDatetimeClauseInvalidTimeJaJP = status.MustWithDetails(
		statusDatetimeClauseInvalidTime,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "???????????????????????????????????????",
		},
	)
	errWebhookNotFoundJaJP = status.MustWithDetails(
		statusWebhookNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "???????????????????????????????????????",
		},
	)
	errWebhookClauseRequiredJaJP = status.MustWithDetails(
		statusWebhookClauseRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "??????????????????????????????????????????",
		},
	)
	errWebhookClauseWebhookIDRequiredJaJP = status.MustWithDetails(
		statusWebhookClauseWebhookIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "??????????????????????????????webhook id???????????????",
		},
	)
	errWebhookClauseConditionRequiredJaJP = status.MustWithDetails(
		statusWebhookClauseConditionRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "??????????????????????????????condition???????????????",
		},
	)
	errWebhookClauseConditionFilterRequiredJaJP = status.MustWithDetails(
		statusWebhookClauseConditionFilterRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "??????????????????????????????condition???filter???????????????",
		},
	)
	// nolint:deadcode,unused,varcheck
	errWebhookClauseConditionValueRequiredJaJP = status.MustWithDetails(
		statusWebhookClauseConditionValueRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "??????????????????????????????condition???value???????????????",
		},
	)
	errWebhookClauseConditionInvalidOperatorJaJP = status.MustWithDetails(
		statusWebhookClauseConditionInvalidOperator,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "??????????????????????????????condition???operator???????????????",
		},
	)
	errNotFoundJaJP = status.MustWithDetails(
		statusNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "??????????????????????????????",
		},
	)
	errAlreadyDeletedJaJP = status.MustWithDetails(
		statusAlreadyDeleted,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "???????????????????????????????????????",
		},
	)
	errOpsEventRateClauseGoalNotFoundJaJP = status.MustWithDetails(
		statusOpsEventRateClauseGoalNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????????????????????????????goal?????????????????????",
		},
	)
	errAlreadyExistsJaJP = status.MustWithDetails(
		statusAlreadyExists,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "??????id???????????????????????????????????????",
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
