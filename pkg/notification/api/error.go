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
	statusInternal            = gstatus.New(codes.Internal, "notification: internal")
	statusIDRequired          = gstatus.New(codes.InvalidArgument, "notification: id must be specified")
	statusNameRequired        = gstatus.New(codes.InvalidArgument, "notification: name must be specified")
	statusSourceTypesRequired = gstatus.New(
		codes.InvalidArgument,
		"notification: notification types must be specified",
	)
	statusUnknownRecipient  = gstatus.New(codes.InvalidArgument, "notification: unknown recipient")
	statusRecipientRequired = gstatus.New(
		codes.InvalidArgument,
		"notification: recipient must be specified",
	)
	statusSlackRecipientRequired = gstatus.New(
		codes.InvalidArgument,
		"notification: slack recipient must be specified",
	)
	statusSlackRecipientWebhookURLRequired = gstatus.New(
		codes.InvalidArgument,
		"notification: webhook URL must be specified",
	)
	statusInvalidCursor    = gstatus.New(codes.InvalidArgument, "notification: cursor is invalid")
	statusNoCommand        = gstatus.New(codes.InvalidArgument, "notification: no command")
	statusInvalidOrderBy   = gstatus.New(codes.InvalidArgument, "environment: order_by is invalid")
	statusNotFound         = gstatus.New(codes.NotFound, "notification: not found")
	statusAlreadyExists    = gstatus.New(codes.AlreadyExists, "notification: already exists")
	statusUnauthenticated  = gstatus.New(codes.Unauthenticated, "notification: unauthenticated")
	statusPermissionDenied = gstatus.New(codes.PermissionDenied, "notification: permission denied")

	errInternalJaJP = status.MustWithDetails(
		statusInternal,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "????????????????????????????????????",
		},
	)
	errIDRequiredJaJP = status.MustWithDetails(
		statusIDRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "id???????????????",
		},
	)
	errNameRequiredJaJP = status.MustWithDetails(
		statusNameRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "name???????????????",
		},
	)
	errSourceTypesRequiredJaJP = status.MustWithDetails(
		statusSourceTypesRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "notification type???????????????????????????",
		},
	)
	errUnknownRecipientJaJP = status.MustWithDetails(
		statusUnknownRecipient,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "?????????recipient??????",
		},
	)
	errRecipientRequiredJaJP = status.MustWithDetails(
		statusRecipientRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "recipient???????????????",
		},
	)
	errSlackRecipientRequiredJaJP = status.MustWithDetails(
		statusSlackRecipientRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "slack recipient???????????????",
		},
	)
	errSlackRecipientWebhookURLRequiredJaJP = status.MustWithDetails(
		statusSlackRecipientWebhookURLRequired,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "slack recipient???webhook url???????????????",
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
	case statusIDRequired:
		return errIDRequiredJaJP
	case statusNameRequired:
		return errNameRequiredJaJP
	case statusSourceTypesRequired:
		return errSourceTypesRequiredJaJP
	case statusUnknownRecipient:
		return errUnknownRecipientJaJP
	case statusRecipientRequired:
		return errRecipientRequiredJaJP
	case statusSlackRecipientRequired:
		return errSlackRecipientRequiredJaJP
	case statusSlackRecipientWebhookURLRequired:
		return errSlackRecipientWebhookURLRequiredJaJP
	case statusInvalidCursor:
		return errInvalidCursorJaJP
	case statusNoCommand:
		return errNoCommandJaJP
	case statusInvalidOrderBy:
		return errInvalidOrderByJaJP
	case statusNotFound:
		return errNotFoundJaJP
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
