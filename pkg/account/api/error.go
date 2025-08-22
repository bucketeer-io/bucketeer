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

package api

import (
	"github.com/bucketeer-io/bucketeer/pkg/api/api"
	pkgErr "github.com/bucketeer-io/bucketeer/pkg/error"
)

// nolint:lll
var (
	statusInternal                               = api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.AccountPackageName, "internal"))
	statusInvalidCursor                          = api.NewGRPCStatus(pkgErr.NewErrorInvalidNotMatchFormat(pkgErr.AccountPackageName, "cursor is invalid", "cursor"))
	statusNoCommand                              = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "command must not be empty", "command"))
	statusMissingOrganizationID                  = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "organization id must be specified", "organization_id"))
	statusEmailIsEmpty                           = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "email is empty", "email"))
	statusInvalidEmail                           = api.NewGRPCStatus(pkgErr.NewErrorInvalidNotMatchFormat(pkgErr.AccountPackageName, "invalid email format", "email"))
	statusFirstNameIsEmpty                       = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "first name is empty", "first name"))
	statusInvalidFirstName                       = api.NewGRPCStatus(pkgErr.NewErrorInvalidNotMatchFormat(pkgErr.AccountPackageName, "invalid first name format", "first name"))
	statusLastNameIsEmpty                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "last name is empty", "last name"))
	statusInvalidLastName                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidNotMatchFormat(pkgErr.AccountPackageName, "invalid last name format", "last name"))
	statusLanguageIsEmpty                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "language is empty", "language"))
	statusInvalidOrganizationRole                = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "invalid organization roles", "organization role"))
	statusInvalidEnvironmentRole                 = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "environment role must be specified", "environment role"))
	statusInvalidUpdateEnvironmentRolesWriteType = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "update environment roles write type must be specified", "update environment roles write type"))
	statusMissingAPIKeyID                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "apikey id must be specified", "apikey id"))
	statusMissingAPIKeyName                      = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "apikey name must be not empty", "apikey name"))
	statusInvalidOrderBy                         = api.NewGRPCStatus(pkgErr.NewErrorInvalidNotMatchFormat(pkgErr.AccountPackageName, "order_by is invalid", "order_by"))
	statusNotFound                               = api.NewGRPCStatus(pkgErr.NewErrorNotFound(pkgErr.AccountPackageName, "account not found", "account"))
	statusAlreadyExists                          = api.NewGRPCStatus(pkgErr.NewErrorAlreadyExists(pkgErr.AccountPackageName, "account already exists"))
	statusUnauthenticated                        = api.NewGRPCStatus(pkgErr.NewErrorUnauthenticated(pkgErr.AccountPackageName, "account unauthenticated"))
	statusPermissionDenied                       = api.NewGRPCStatus(pkgErr.NewErrorPermissionDenied(pkgErr.AccountPackageName, "permission denied"))
	statusSearchFilterNameIsEmpty                = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "search filter name is empty", "search filter name"))
	statusSearchFilterQueryIsEmpty               = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "search filter query is empty", "search filter query"))
	statusSearchFilterTargetTypeIsRequired       = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "search filter target type is required", "search filter target type"))
	statusSearchFilterIDIsEmpty                  = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "search filter ID is empty", "search filter ID"))
	statusSearchFilterIDNotFound                 = api.NewGRPCStatus(pkgErr.NewErrorNotFound(pkgErr.AccountPackageName, "search filter not found", "search_filter"))
	statusInvalidListAPIKeyRequest               = api.NewGRPCStatus(pkgErr.NewErrorInvalidEmpty(pkgErr.AccountPackageName, "invalid list api key request", "list api key request"))
)
