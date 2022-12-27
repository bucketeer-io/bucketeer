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
	statusInternal             = gstatus.New(codes.Internal, "eventcounter: internal")
	statusFeatureIDRequired    = gstatus.New(codes.InvalidArgument, "eventcounter: feature id is required")
	statusExperimentIDRequired = gstatus.New(codes.InvalidArgument, "eventcounter: experiment id is required")
	statusMAUYearMonthRequired = gstatus.New(codes.InvalidArgument, "eventcounter: mau year month is required")
	statusGoalIDRequired       = gstatus.New(codes.InvalidArgument, "eventcounter: goal id is required")
	statusStartAtRequired      = gstatus.New(codes.InvalidArgument, "eventcounter: start at is required")
	statusEndAtRequired        = gstatus.New(codes.InvalidArgument, "eventcounter: end at is required")
	statusPeriodOutOfRange     = gstatus.New(codes.InvalidArgument, "eventcounter: period out of range")
	statusStartAtIsAfterEndAt  = gstatus.New(codes.InvalidArgument, "eventcounter: start at is after end at")
	statusNotFound             = gstatus.New(codes.NotFound, "eventcounter: not found")
	statusUnauthenticated      = gstatus.New(codes.Unauthenticated, "feature: unauthenticated")
	statusPermissionDenied     = gstatus.New(codes.PermissionDenied, "feature: permission denied")

	errInternalJaJP = status.MustWithDetails(
		statusInternal,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "内部エラーが発生しました",
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
	errStartAtRequiredJaJP = status.MustWithDetails(
		statusStartAtRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "start atは必須です",
		},
	)
	errEndAtRequiredJaJP = status.MustWithDetails(
		statusEndAtRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "end atは必須です",
		},
	)
	errStartAtIsAfterEndAtJaJP = status.MustWithDetails(
		statusStartAtIsAfterEndAt,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "start at はend at以前を指定してください。",
		},
	)
	errPeroidOutOfRangeJaJP = status.MustWithDetails(
		statusPeriodOutOfRange,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "期間は過去30日以内を選択してください。",
		},
	)
	errNotFoundJaJP = status.MustWithDetails(
		statusNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "データが存在しません",
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
	case statusFeatureIDRequired:
		return errFeatureIDRequiredJaJP
	case statusExperimentIDRequired:
		return errExperimentIDRequiredJaJP
	case statusGoalIDRequired:
		return errGoalIDRequiredJaJP
	case statusStartAtRequired:
		return errStartAtRequiredJaJP
	case statusEndAtRequired:
		return errEndAtRequiredJaJP
	case statusPeriodOutOfRange:
		return errPeroidOutOfRangeJaJP
	case statusStartAtIsAfterEndAt:
		return errStartAtIsAfterEndAtJaJP
	case statusNotFound:
		return errNotFoundJaJP
	case statusUnauthenticated:
		return errUnauthenticatedJaJP
	case statusPermissionDenied:
		return errPermissionDeniedJaJP
	default:
		return errInternalJaJP
	}
}
