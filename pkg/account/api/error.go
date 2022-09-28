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
	statusInternal          = gstatus.New(codes.Internal, "account: internal")
	statusInvalidCursor     = gstatus.New(codes.InvalidArgument, "account: cursor is invalid")
	statusNoCommand         = gstatus.New(codes.InvalidArgument, "account: command must not be empty")
	statusMissingAccountID  = gstatus.New(codes.InvalidArgument, "account: account id must be specified")
	statusEmailIsEmpty      = gstatus.New(codes.InvalidArgument, "account: email is empty")
	statusInvalidEmail      = gstatus.New(codes.InvalidArgument, "account: invalid email format")
	statusMissingAPIKeyID   = gstatus.New(codes.InvalidArgument, "account: apikey id must be specified")
	statusMissingAPIKeyName = gstatus.New(codes.InvalidArgument, "account: apikey name must be not empty")
	statusInvalidOrderBy    = gstatus.New(codes.InvalidArgument, "account: order_by is invalid")
	statusNotFound          = gstatus.New(codes.NotFound, "account: not found")
	statusAlreadyExists     = gstatus.New(codes.AlreadyExists, "account: already exists")
	statusUnauthenticated   = gstatus.New(codes.Unauthenticated, "account: unauthenticated")
	statusPermissionDenied  = gstatus.New(codes.PermissionDenied, "account: permission denied")

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
	errMissingAccountIDJaJP = status.MustWithDetails(
		statusMissingAccountID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "account idは必須です",
		},
	)
	errEmailIsEmptyJaJP = status.MustWithDetails(
		statusEmailIsEmpty,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "emailは必須です",
		},
	)
	errInvalidEmailJaJP = status.MustWithDetails(
		statusInvalidEmail,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なemailです",
		},
	)
	errMissingAPIKeyIDJaJP = status.MustWithDetails(
		statusMissingAPIKeyID,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "api keyのidは必須です",
		},
	)
	errMissingAPIKeyNameJaJP = status.MustWithDetails(
		statusMissingAPIKeyName,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "api keyのnameは必須です",
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
	case statusMissingAccountID:
		return errMissingAccountIDJaJP
	case statusEmailIsEmpty:
		return errEmailIsEmptyJaJP
	case statusInvalidEmail:
		return errInvalidEmailJaJP
	case statusMissingAPIKeyID:
		return errMissingAPIKeyIDJaJP
	case statusMissingAPIKeyName:
		return errMissingAPIKeyNameJaJP
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
