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
	statusInternal                   = gstatus.New(codes.Internal, "environment: internal")
	statusNoCommand                  = gstatus.New(codes.InvalidArgument, "environment: no command")
	statusInvalidCursor              = gstatus.New(codes.InvalidArgument, "environment: cursor is invalid")
	statusEnvironmentIDRequired      = gstatus.New(codes.InvalidArgument, "environment: environment id must be specified")
	statusInvalidEnvironmentID       = gstatus.New(codes.InvalidArgument, "environment: invalid environment id")
	statusProjectIDRequired          = gstatus.New(codes.InvalidArgument, "environment: project id must be specified")
	statusInvalidProjectID           = gstatus.New(codes.InvalidArgument, "environment: invalid project id")
	statusInvalidProjectCreatorEmail = gstatus.New(codes.InvalidArgument, "environment: invalid project creator email")
	statusInvalidOrderBy             = gstatus.New(codes.InvalidArgument, "environment: order_by is invalid")
	statusEnvironmentNotFound        = gstatus.New(codes.NotFound, "environment: environment not found")
	statusProjectNotFound            = gstatus.New(codes.NotFound, "environment: project not found")
	statusEnvironmentAlreadyDeleted  = gstatus.New(codes.NotFound, "environment: environment already deleted")
	statusEnvironmentAlreadyExists   = gstatus.New(codes.AlreadyExists, "environment: environment already exists")
	statusProjectAlreadyExists       = gstatus.New(codes.AlreadyExists, "environment: project already exists")
	statusProjectDisabled            = gstatus.New(codes.FailedPrecondition, "environment: project disabled")
	statusUnauthenticated            = gstatus.New(codes.Unauthenticated, "environment: unauthenticated")
	statusPermissionDenied           = gstatus.New(codes.PermissionDenied, "environment: permission denied")

	errInternalJaJP = status.MustWithDetails(
		statusInternal,
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
	errInvalidCursorJaJP = status.MustWithDetails(
		statusInvalidCursor,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????cursor??????",
		},
	)
	errEnvironmentIDRequiredJaJP = status.MustWithDetails(
		statusEnvironmentIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "environment id???????????????",
		},
	)
	errInvalidEnvironmentIDJaJP = status.MustWithDetails(
		statusInvalidEnvironmentID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????environment id??????",
		},
	)
	errProjectIDRequiredJaJP = status.MustWithDetails(
		statusProjectIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "project id???????????????",
		},
	)
	errInvalidProjectIDJaJP = status.MustWithDetails(
		statusInvalidProjectID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????project id??????",
		},
	)
	errInvalidProjectCreatorEmailJaJP = status.MustWithDetails(
		statusInvalidProjectCreatorEmail,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "Project????????????email???????????????",
		},
	)
	errInvalidOrderByJaJP = status.MustWithDetails(
		statusInvalidOrderBy,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "????????????????????????????????????",
		},
	)
	errEnvironmentNotFoundJaJP = status.MustWithDetails(
		statusEnvironmentNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "environment?????????????????????????????????",
		},
	)
	errProjectNotFoundJaJP = status.MustWithDetails(
		statusProjectNotFound,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "project?????????????????????????????????",
		},
	)
	errEnvironmentAlreadyDeletedJaJP = status.MustWithDetails(
		statusEnvironmentAlreadyDeleted,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "environment??????????????????????????????????????????",
		},
	)
	errEnvironmentAlreadyExistsJaJP = status.MustWithDetails(
		statusEnvironmentAlreadyExists,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "??????id?????????namespace???environment???????????????????????????????????????",
		},
	)
	errProjectAlreadyExistsJaJP = status.MustWithDetails(
		statusProjectAlreadyExists,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "??????id???project???????????????????????????????????????",
		},
	)
	errProjectDisabledJaJp = status.MustWithDetails(
		statusProjectDisabled,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "project??????????????????????????????????????????",
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
	case statusNoCommand:
		return errNoCommandJaJP
	case statusInvalidCursor:
		return errInvalidCursorJaJP
	case statusEnvironmentIDRequired:
		return errEnvironmentIDRequiredJaJP
	case statusInvalidEnvironmentID:
		return errInvalidEnvironmentIDJaJP
	case statusProjectIDRequired:
		return errProjectIDRequiredJaJP
	case statusInvalidProjectID:
		return errInvalidProjectIDJaJP
	case statusInvalidProjectCreatorEmail:
		return errInvalidProjectCreatorEmailJaJP
	case statusInvalidOrderBy:
		return errInvalidOrderByJaJP
	case statusEnvironmentNotFound:
		return errEnvironmentNotFoundJaJP
	case statusProjectNotFound:
		return errProjectNotFoundJaJP
	case statusEnvironmentAlreadyDeleted:
		return errEnvironmentAlreadyDeletedJaJP
	case statusEnvironmentAlreadyExists:
		return errEnvironmentAlreadyExistsJaJP
	case statusProjectAlreadyExists:
		return errProjectAlreadyExistsJaJP
	case statusProjectDisabled:
		return errProjectDisabledJaJp
	case statusUnauthenticated:
		return errUnauthenticatedJaJP
	case statusPermissionDenied:
		return errPermissionDeniedJaJP
	default:
		return errInternalJaJP
	}
}
