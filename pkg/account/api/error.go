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
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
)

// nolint:lll
var (
	statusInternal                               = api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.AccountPackageName, "internal"))
	statusInvalidCursor                          = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AccountPackageName, "cursor is invalid", "cursor"))
	statusNoCommand                              = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "command must not be empty", "command"))
	statusMissingOrganizationID                  = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "organization id must be specified", "organization_id"))
	statusEmailIsEmpty                           = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "email is empty", "email"))
	statusInvalidEmail                           = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AccountPackageName, "invalid email format", "email"))
	statusFirstNameIsEmpty                       = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "first name is empty", "first_name"))
	statusInvalidFirstName                       = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AccountPackageName, "invalid first name format", "first_name"))
	statusLastNameIsEmpty                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "last name is empty", "last_name"))
	statusInvalidLastName                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AccountPackageName, "invalid last name format", "last_name"))
	statusLanguageIsEmpty                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "language is empty", "language"))
	statusInvalidOrganizationRole                = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "invalid organization roles", "organization_role"))
	statusInvalidEnvironmentRole                 = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "environment role must be specified", "environment_role"))
	statusInvalidUpdateEnvironmentRolesWriteType = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "update environment roles write type must be specified", "update_environment_roles_write_type"))
	statusMissingAPIKeyID                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "apikey id must be specified", "apikey_id"))
	statusMissingAPIKeyName                      = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "apikey name must be not empty", "apikey_name"))
	statusInvalidOrderBy                         = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AccountPackageName, "order_by is invalid", "order_by"))
	statusNotFound                               = api.NewGRPCStatus(pkgErr.NewErrorNotFound(pkgErr.AccountPackageName, "account not found", "account"))
	statusAlreadyExists                          = api.NewGRPCStatus(pkgErr.NewErrorAlreadyExists(pkgErr.AccountPackageName, "account already exists"))
	statusUnauthenticated                        = api.NewGRPCStatus(pkgErr.NewErrorUnauthenticated(pkgErr.AccountPackageName, "account unauthenticated"))
	statusPermissionDenied                       = api.NewGRPCStatus(pkgErr.NewErrorPermissionDenied(pkgErr.AccountPackageName, "permission denied"))
	statusSearchFilterNameIsEmpty                = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "search filter name is empty", "search_filter_name"))
	statusSearchFilterQueryIsEmpty               = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "search filter query is empty", "search_filter_query"))
	statusSearchFilterTargetTypeIsRequired       = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "search filter target type is required", "search_filter_target_type"))
	statusSearchFilterIDIsEmpty                  = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "search filter ID is empty", "search_filter_ID"))
	statusSearchFilterIDNotFound                 = api.NewGRPCStatus(pkgErr.NewErrorNotFound(pkgErr.AccountPackageName, "search filter not found", "search_filter"))
	statusInvalidListAPIKeyRequest               = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "invalid list api key request", "list_api_key_request"))
)
