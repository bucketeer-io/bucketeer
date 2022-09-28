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
	statusInternal                = gstatus.New(codes.Internal, "auth: internal")
	statusMissingCode             = gstatus.New(codes.InvalidArgument, "auth: code must not be empty")
	statusMissingState            = gstatus.New(codes.InvalidArgument, "auth: state must not be empty")
	statusMissingRedirectURL      = gstatus.New(codes.InvalidArgument, "auth: missing redirectURL")
	statusUnregisteredRedirectURL = gstatus.New(codes.InvalidArgument, "auth: unregistered redirectURL")
	statusMissingRefreshToken     = gstatus.New(codes.InvalidArgument, "auth: refreshToken must not be empty")
	statusInvalidCode             = gstatus.New(codes.InvalidArgument, "auth: invalid code")
	statusInvalidRefreshToken     = gstatus.New(codes.InvalidArgument, "auth: invalid refresh token")
	statusUnapprovedAccount       = gstatus.New(codes.PermissionDenied, "auth: unapproved account")
	statusAccessDeniedEmail       = gstatus.New(codes.PermissionDenied, "auth: access denied email")
	statusMissingEncryptedSecret  = gstatus.New(codes.InvalidArgument, "auth: encrypted secret must not be empty")
	statusUnauthenticated         = gstatus.New(codes.Unauthenticated, "auth: unauthenticated")
	statusPermissionDenied        = gstatus.New(codes.PermissionDenied, "auth: permission denied")

	errInternalJaJP = status.MustWithDetails(
		statusInternal,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "内部エラーが発生しました",
		},
	)
	errMissingCodeJaJP = status.MustWithDetails(
		statusMissingCode,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "認証codeは必須です",
		},
	)
	errMissingStateJaJP = status.MustWithDetails(
		statusMissingState,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "stateは必須です",
		},
	)
	errMissingRedirectURLJaJP = status.MustWithDetails(
		statusMissingRedirectURL,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "redirect urlは必須です",
		},
	)
	errUnregisteredRedirectURLJaJP = status.MustWithDetails(
		statusUnregisteredRedirectURL,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "許可されていないredirect urlです",
		},
	)
	errMissingRefreshTokenJaJP = status.MustWithDetails(
		statusMissingRefreshToken,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "refresh tokenは必須です",
		},
	)
	errInvalidCodeJaJP = status.MustWithDetails(
		statusInvalidCode,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なcodeです",
		},
	)
	errInvalidRefreshTokenJaJP = status.MustWithDetails(
		statusInvalidRefreshToken,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "不正なrefresh tokenです",
		},
	)
	errUnapprovedAccountJaJP = status.MustWithDetails(
		statusUnapprovedAccount,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "許可されていないaccountです",
		},
	)
	errAccessDeniedEmailJaJP = status.MustWithDetails(
		statusAccessDeniedEmail,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "許可されていないemailです",
		},
	)
	errMissingEncryptedSecretJaJP = status.MustWithDetails(
		statusMissingEncryptedSecret,
		&errdetails.LocalizedMessage{
			Locale:  locale.JaJP,
			Message: "encrypted secretは必須です",
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
	case statusMissingCode:
		return errMissingCodeJaJP
	case statusMissingState:
		return errMissingStateJaJP
	case statusMissingRedirectURL:
		return errMissingRedirectURLJaJP
	case statusUnregisteredRedirectURL:
		return errUnregisteredRedirectURLJaJP
	case statusMissingRefreshToken:
		return errMissingRefreshTokenJaJP
	case statusInvalidCode:
		return errInvalidCodeJaJP
	case statusInvalidRefreshToken:
		return errInvalidRefreshTokenJaJP
	case statusUnapprovedAccount:
		return errUnapprovedAccountJaJP
	case statusAccessDeniedEmail:
		return errAccessDeniedEmailJaJP
	case statusMissingEncryptedSecret:
		return errMissingEncryptedSecretJaJP
	case statusUnauthenticated:
		return errUnauthenticatedJaJP
	case statusPermissionDenied:
		return errPermissionDeniedJaJP
	default:
		return errInternalJaJP
	}
}
