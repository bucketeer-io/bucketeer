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

package auth

import (
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
)

var (
	StatusInternal    = api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.AuthPackageName, "internal"))
	StatusMissingCode = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "code must not be empty", "code"))
	StatusMissingState = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "state must not be empty", "state"))
	StatusMissingAuthType = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AuthPackageName, "missing authType", "authType"))
	StatusUnknownAuthType = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgUnknown(pkgErr.AuthPackageName, "unknown authType", "authType"))
	StatusMissingRedirectURL = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "missing redirectURL", "redirectURL"))
	StatusMissingRefreshToken = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "refreshToken must not be empty", "refreshToken"))
	StatusUnapprovedAccount = api.NewGRPCStatus(
		pkgErr.NewErrorPermissionDenied(pkgErr.AuthPackageName, "unapproved account"))
	StatusAccessDeniedEmail = api.NewGRPCStatus(
		pkgErr.NewErrorPermissionDenied(pkgErr.AuthPackageName, "access denied email"))
	StatusUnauthenticated = api.NewGRPCStatus(
		pkgErr.NewErrorUnauthenticated(pkgErr.AuthPackageName, "not authenticated"))
	StateMissingUsername = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "missing username", "username"))
	StateMissingPassword = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "missing password", "password"))
	StatusAccessDenied = api.NewGRPCStatus(
		pkgErr.NewErrorPermissionDenied(pkgErr.AuthPackageName, "access denied"))
)
