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
	statusInvalidCursor                          = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AccountPackageName, "cursor is invalid", "Cursor"))
	statusNoCommand                              = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "command must not be empty", "Command"))
	statusMissingOrganizationID                  = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "organization id must be specified", "Organization"))
	statusEmailIsEmpty                           = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "email is empty", "Email"))
	statusInvalidEmail                           = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AccountPackageName, "invalid email format", "Email"))
	statusFirstNameIsEmpty                       = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "first name is empty", "MemberFirstName"))
	statusInvalidFirstName                       = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AccountPackageName, "invalid first name format", "MemberFirstName"))
	statusLastNameIsEmpty                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "last name is empty", "MemberLastName"))
	statusInvalidLastName                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AccountPackageName, "invalid last name format", "MemberLastName"))
	statusLanguageIsEmpty                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "language is empty", "MemberLanguage"))
	statusInvalidOrganizationRole                = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "invalid organization roles", "MemberOrganizationRole"))
	statusInvalidEnvironmentRole                 = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "environment role must be specified", "MemberEnvironmentRoles"))
	statusInvalidUpdateEnvironmentRolesWriteType = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "update environment roles write type must be specified", "UpdateEnvironmentRolesWriteType"))
	statusMissingAPIKeyID                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "apikey id must be specified", "APIKey"))
	statusMissingAPIKeyName                      = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "apikey name must be not empty", "APIKey"))
	statusInvalidOrderBy                         = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AccountPackageName, "order_by is invalid", "OrderBy"))
	statusAccountNotFound                        = api.NewGRPCStatus(pkgErr.NewErrorNotFound(pkgErr.AccountPackageName, "member not found", "Member"))
	statusAPIKeyNotFound                         = api.NewGRPCStatus(pkgErr.NewErrorNotFound(pkgErr.AccountPackageName, "api key not found", "APIKey"))
	statusOrganizationNotFound                   = api.NewGRPCStatus(pkgErr.NewErrorNotFound(pkgErr.AccountPackageName, "organization not found", "Organization"))
	statusAccountAlreadyExists                   = api.NewGRPCStatus(pkgErr.NewErrorAlreadyExists(pkgErr.AccountPackageName, "member already exists"))
	statusAPIKeyAlreadyExists                    = api.NewGRPCStatus(pkgErr.NewErrorAlreadyExists(pkgErr.AccountPackageName, "api key already exists"))
	statusUnauthenticated                        = api.NewGRPCStatus(pkgErr.NewErrorUnauthenticated(pkgErr.AccountPackageName, "member unauthenticated"))
	statusPermissionDenied                       = api.NewGRPCStatus(pkgErr.NewErrorPermissionDenied(pkgErr.AccountPackageName, "permission denied"))
	statusSearchFilterNameIsEmpty                = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "search filter name is empty", "MemberSearchFilterName"))
	statusSearchFilterQueryIsEmpty               = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "search filter query is empty", "MemberSearchFilterQuery"))
	statusSearchFilterTargetTypeIsRequired       = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "search filter target type is required", "SearchFilterTargetType"))
	statusSearchFilterIDIsEmpty                  = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "search filter ID is empty", "SearchFilterId"))
	statusSearchFilterIDNotFound                 = api.NewGRPCStatus(pkgErr.NewErrorNotFound(pkgErr.AccountPackageName, "search filter not found", "MemberSearchFilter"))
	statusInvalidListAPIKeyRequest               = api.NewGRPCStatus(pkgErr.NewErrorInvalidArgEmpty(pkgErr.AccountPackageName, "invalid list api key request", "ListAPIKeyRequest"))
)
