// Copyright 2026 The Bucketeer Authors.
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
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
)

var (
	statusMissingCode = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "code must not be empty", "Code"))
	statusMissingState = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "state must not be empty", "State"))
	statusMissingAuthType = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "auth type must be specified", "AuthType"))
	statusUnknownAuthType = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgUnknown(pkgErr.AuthPackageName, "unknown auth type", "AuthType"))
	statusMissingRedirectURL = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "redirect url must not be empty", "RedirectUrl"))
	statusMissingRefreshToken = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "refresh token must not be empty", "RefreshToken"))
	statusMissingUsername = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "username must not be empty", "Username"))
	statusMissingPassword = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "password must not be empty", "Password"))
	statusUnapprovedAccount = api.NewGRPCStatus(
		pkgErr.NewErrorPermissionDenied(pkgErr.AuthPackageName, "unapproved account"))
	statusAccessDeniedEmail = api.NewGRPCStatus(
		pkgErr.NewErrorPermissionDenied(pkgErr.AuthPackageName, "access denied email"))
	statusAccessDenied = api.NewGRPCStatus(
		pkgErr.NewErrorPermissionDenied(pkgErr.AuthPackageName, "access denied"))
	statusUnauthenticated = api.NewGRPCStatus(
		pkgErr.NewErrorUnauthenticated(pkgErr.AuthPackageName, "not authenticated"))

	// Password-related errors
	statusMissingCurrentPassword = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "current password must not be empty", "CurrentPassword"))
	statusMissingNewPassword = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "new password must not be empty", "NewPassword"))
	statusPasswordsIdentical = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.AuthPackageName,
			"new password must be different from current password",
			"NewPassword",
		))
	statusMissingEmail = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "email must not be empty", "Email"))
	statusInvalidEmailConfig = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AuthPackageName, "invalid email configuration", "EmailConfig"))
	statusPasswordNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.AuthPackageName, "password not found", "Password"))
	statusPasswordMismatch = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AuthPackageName, "password mismatch", "Password"))
	statusPasswordTooWeak = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AuthPackageName, "password too weak", "Password"))
	statusPasswordAlreadyExists = api.NewGRPCStatus(
		pkgErr.NewErrorAlreadyExists(pkgErr.AuthPackageName, "password already exists"))
	statusInvalidResetToken = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AuthPackageName, "invalid reset token", "ResetToken"))
	statusExpiredResetToken = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(pkgErr.AuthPackageName, "reset token expired"))
	statusMissingResetToken = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "reset token must not be empty", "ResetToken"))
	statusEmailServiceUnavailable = api.NewGRPCStatus(
		pkgErr.NewErrorUnavailable(pkgErr.AuthPackageName, "email service unavailable"))
)
