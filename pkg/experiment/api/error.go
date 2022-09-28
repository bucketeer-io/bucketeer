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
	statusInternal             = gstatus.New(codes.Internal, "experiment: internal")
	statusInvalidCursor        = gstatus.New(codes.InvalidArgument, "experiment: cursor is invalid")
	statusNoCommand            = gstatus.New(codes.InvalidArgument, "experiment: must contain at least one command")
	statusUnknownCommand       = gstatus.New(codes.InvalidArgument, "experiment: unknown command")
	statusFeatureIDRequired    = gstatus.New(codes.InvalidArgument, "experiment: feature id must be specified")
	statusExperimentIDRequired = gstatus.New(codes.InvalidArgument, "experiment: experiment id must be specified")
	statusGoalIDRequired       = gstatus.New(codes.InvalidArgument, "experiment: goal id must be specified")
	statusInvalidGoalID        = gstatus.New(codes.InvalidArgument, "experiment: invalid goal id")
	statusGoalNameRequired     = gstatus.New(codes.InvalidArgument, "experiment: goal name must be specified")
	statusPeriodTooLong        = gstatus.New(codes.InvalidArgument, "experiment: period too long")
	statusInvalidOrderBy       = gstatus.New(codes.InvalidArgument, "expriment: order_by is invalid")
	statusNotFound             = gstatus.New(codes.NotFound, "experiment: not found")
	statusGoalNotFound         = gstatus.New(codes.NotFound, "experiment: goal not found")
	statusFeatureNotFound      = gstatus.New(codes.NotFound, "experiment: feature not found")
	statusAlreadyExists        = gstatus.New(codes.AlreadyExists, "experiment: already exists")
	statusUnauthenticated      = gstatus.New(codes.Unauthenticated, "experiment: unauthenticated")
	statusPermissionDenied     = gstatus.New(codes.PermissionDenied, "experiment: permission denied")

	errInternalJaJP = status.MustWithDetails(
		statusInternal,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "内部エラーが発生しました",
		},
	)
	errInvalidCursorJaJP = status.MustWithDetails(
		statusInvalidCursor,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なcursorです",
		},
	)
	errNoCommandJaJP = status.MustWithDetails(
		statusNoCommand,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "commandは必須です",
		},
	)
	errUnknownCommandJaJP = status.MustWithDetails(
		statusUnknownCommand,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不明なcommandです",
		},
	)
	errFeatureIDRequiredJaJP = status.MustWithDetails(
		statusFeatureIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "feature idは必須です",
		},
	)
	errExperimentIDRequiredJaJP = status.MustWithDetails(
		statusExperimentIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "experiment idは必須です",
		},
	)
	errGoalIDRequiredJaJP = status.MustWithDetails(
		statusGoalIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "goal idは必須です",
		},
	)
	errInvalidGoalIDJaJP = status.MustWithDetails(
		statusInvalidGoalID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なgoal idです",
		},
	)
	errGoalNameRequiredJaJP = status.MustWithDetails(
		statusGoalNameRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "goal nameは必須です",
		},
	)
	errPeriodTooLongJaJP = status.MustWithDetails(
		statusPeriodTooLong,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: fmt.Sprintf("experiment期間は%d日以内で設定してください", maxExperimentPeriodDays),
		},
	)
	errInvalidOrderByJaJP = status.MustWithDetails(
		statusInvalidOrderBy,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なソート順の指定です",
		},
	)
	errNotFoundJaJP = status.MustWithDetails(
		statusNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "データが存在しません",
		},
	)
	errFeatureNotFoundJaJP = status.MustWithDetails(
		statusFeatureNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "featureが存在しません",
		},
	)
	errGoalNotFoundJaJP = status.MustWithDetails(
		statusGoalNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "goalが存在しません",
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
	case statusInvalidCursor:
		return errInvalidCursorJaJP
	case statusNoCommand:
		return errNoCommandJaJP
	case statusUnknownCommand:
		return errUnknownCommandJaJP
	case statusFeatureIDRequired:
		return errFeatureIDRequiredJaJP
	case statusExperimentIDRequired:
		return errExperimentIDRequiredJaJP
	case statusGoalIDRequired:
		return errGoalIDRequiredJaJP
	case statusInvalidGoalID:
		return errInvalidGoalIDJaJP
	case statusGoalNameRequired:
		return errGoalNameRequiredJaJP
	case statusPeriodTooLong:
		return errPeriodTooLongJaJP
	case statusInvalidOrderBy:
		return errInvalidOrderByJaJP
	case statusNotFound:
		return errNotFoundJaJP
	case statusFeatureNotFound:
		return errFeatureNotFoundJaJP
	case statusGoalNotFound:
		return errGoalNotFoundJaJP
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
