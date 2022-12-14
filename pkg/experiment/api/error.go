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
			Message: "????????????????????????????????????",
		},
	)
	errInvalidCursorJaJP = status.MustWithDetails(
		statusInvalidCursor,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????cursor??????",
		},
	)
	errNoCommandJaJP = status.MustWithDetails(
		statusNoCommand,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "command???????????????",
		},
	)
	errUnknownCommandJaJP = status.MustWithDetails(
		statusUnknownCommand,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????command??????",
		},
	)
	errFeatureIDRequiredJaJP = status.MustWithDetails(
		statusFeatureIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "feature id???????????????",
		},
	)
	errExperimentIDRequiredJaJP = status.MustWithDetails(
		statusExperimentIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "experiment id???????????????",
		},
	)
	errGoalIDRequiredJaJP = status.MustWithDetails(
		statusGoalIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "goal id???????????????",
		},
	)
	errInvalidGoalIDJaJP = status.MustWithDetails(
		statusInvalidGoalID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????goal id??????",
		},
	)
	errGoalNameRequiredJaJP = status.MustWithDetails(
		statusGoalNameRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "goal name???????????????",
		},
	)
	errPeriodTooLongJaJP = status.MustWithDetails(
		statusPeriodTooLong,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: fmt.Sprintf("experiment?????????%d????????????????????????????????????", maxExperimentPeriodDays),
		},
	)
	errInvalidOrderByJaJP = status.MustWithDetails(
		statusInvalidOrderBy,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "????????????????????????????????????",
		},
	)
	errNotFoundJaJP = status.MustWithDetails(
		statusNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "??????????????????????????????",
		},
	)
	errFeatureNotFoundJaJP = status.MustWithDetails(
		statusFeatureNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "feature?????????????????????",
		},
	)
	errGoalNotFoundJaJP = status.MustWithDetails(
		statusGoalNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "goal?????????????????????",
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
