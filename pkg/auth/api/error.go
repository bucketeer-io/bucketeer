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
)
