// Copyright 2025 The Bucketeer Authors.
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

package auth

import (
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

var (
	StatusInternal                = gstatus.New(codes.Internal, "auth: internal")
	StatusMissingCode             = gstatus.New(codes.InvalidArgument, "auth: code must not be empty")
	StatusMissingState            = gstatus.New(codes.InvalidArgument, "auth: state must not be empty")
	StatusMissingAuthType         = gstatus.New(codes.InvalidArgument, "auth: missing authType")
	StatusUnknownAuthType         = gstatus.New(codes.InvalidArgument, "auth: unknown authType")
	StatusMissingRedirectURL      = gstatus.New(codes.InvalidArgument, "auth: missing redirectURL")
	StatusUnregisteredRedirectURL = gstatus.New(codes.InvalidArgument, "auth: unregistered redirectURL")
	StatusMissingRefreshToken     = gstatus.New(codes.InvalidArgument, "auth: refreshToken must not be empty")
	StatusInvalidCode             = gstatus.New(codes.InvalidArgument, "auth: invalid code")
	StatusInvalidRefreshToken     = gstatus.New(codes.InvalidArgument, "auth: invalid refresh token")
	StatusUnapprovedAccount       = gstatus.New(codes.PermissionDenied, "auth: unapproved account")
	StatusAccessDeniedEmail       = gstatus.New(codes.PermissionDenied, "auth: access denied email")
	StatusUnauthenticated         = gstatus.New(codes.Unauthenticated, "auth: not authenticated")
	StateMissingUsername          = gstatus.New(codes.InvalidArgument, "auth: missing username")
	StateMissingPassword          = gstatus.New(codes.InvalidArgument, "auth: missing password")
	StatusAccessDenied            = gstatus.New(codes.PermissionDenied, "auth: access denied")
	StatusInvalidOrganization     = gstatus.New(codes.InvalidArgument, "auth: invalid organization")

	// Password-related errors
	StatusPasswordTooWeak        = gstatus.New(codes.InvalidArgument, "auth: password too weak")
	StatusPasswordMismatch       = gstatus.New(codes.InvalidArgument, "auth: password mismatch")
	StatusPasswordAlreadyExists  = gstatus.New(codes.AlreadyExists, "auth: password already exists")
	StatusPasswordNotFound       = gstatus.New(codes.NotFound, "auth: password not found")
	StatusMissingCurrentPassword = gstatus.New(codes.InvalidArgument, "auth: current password must not be empty")
	StatusMissingNewPassword     = gstatus.New(codes.InvalidArgument, "auth: new password must not be empty")

	// Password reset errors
	StatusInvalidResetToken  = gstatus.New(codes.InvalidArgument, "auth: invalid reset token")
	StatusExpiredResetToken  = gstatus.New(codes.InvalidArgument, "auth: reset token expired")
	StatusResetTokenNotFound = gstatus.New(codes.NotFound, "auth: reset token not found")
	StatusMissingResetToken  = gstatus.New(codes.InvalidArgument, "auth: reset token must not be empty")

	// Email service errors
	StatusEmailServiceUnavailable = gstatus.New(codes.Unavailable, "auth: email service unavailable")
	StatusTooManyEmailRequests    = gstatus.New(codes.ResourceExhausted, "auth: too many email requests")
	StatusInvalidEmailConfig      = gstatus.New(codes.InvalidArgument, "auth: invalid email configuration")
	StatusMissingEmail            = gstatus.New(codes.InvalidArgument, "auth: email must not be empty")
)
