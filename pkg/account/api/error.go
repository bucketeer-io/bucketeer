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

const packageName = "account"

var (
	statusInternal                               = api.NewGRPCStatus(pkgErr.NewErrorInternal(packageName, "internal"))
	statusInvalidCursor                          = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "cursor is invalid", pkgErr.InvalidTypeNotMatchFormat, "cursor"))
	statusNoCommand                              = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "command must not be empty", pkgErr.InvalidTypeEmpty, "command"))
	statusMissingOrganizationID                  = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "organization id must be specified", pkgErr.InvalidTypeEmpty, "organization id"))
	statusEmailIsEmpty                           = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "email is empty", pkgErr.InvalidTypeEmpty, "email"))
	statusInvalidEmail                           = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "invalid email format", pkgErr.InvalidTypeNotMatchFormat, "email"))
	statusFirstNameIsEmpty                       = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "first name is empty", pkgErr.InvalidTypeEmpty, "first name"))
	statusInvalidFirstName                       = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "invalid first name format", pkgErr.InvalidTypeNotMatchFormat, "first name"))
	statusLastNameIsEmpty                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "last name is empty", pkgErr.InvalidTypeEmpty, "last name"))
	statusInvalidLastName                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "invalid last name format", pkgErr.InvalidTypeNotMatchFormat, "last name"))
	statusLanguageIsEmpty                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "language is empty", pkgErr.InvalidTypeEmpty, "language"))
	statusInvalidOrganizationRole                = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "invalid organization roles", pkgErr.InvalidTypeEmpty, "organization role"))
	statusInvalidEnvironmentRole                 = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "environment role must be specified", pkgErr.InvalidTypeEmpty, "environment role"))
	statusInvalidUpdateEnvironmentRolesWriteType = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "update environment roles write type must be specified", pkgErr.InvalidTypeEmpty, "update environment roles write type"))
	statusMissingAPIKeyID                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "apikey id must be specified", pkgErr.InvalidTypeEmpty, "apikey id"))
	statusMissingAPIKeyName                      = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "apikey name must be not empty", pkgErr.InvalidTypeEmpty, "apikey name"))
	statusInvalidOrderBy                         = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "order_by is invalid", pkgErr.InvalidTypeNotMatchFormat, "order_by"))
	statusNotFound                               = api.NewGRPCStatus(pkgErr.NewErrorNotFound(packageName, "account not found", "account"))
	statusAlreadyExists                          = api.NewGRPCStatus(pkgErr.NewErrorAlreadyExists(packageName, "account already exists", "account"))
	statusUnauthenticated                        = api.NewGRPCStatus(pkgErr.NewErrorUnauthenticated(packageName, "account unauthenticated"))
	statusPermissionDenied                       = api.NewGRPCStatus(pkgErr.NewErrorPermissionDenied(packageName, "permission denied"))
	statusSearchFilterNameIsEmpty                = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "search filter name is empty", pkgErr.InvalidTypeEmpty, "search filter name"))
	statusSearchFilterQueryIsEmpty               = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "search filter query is empty", pkgErr.InvalidTypeEmpty, "search filter query"))
	statusSearchFilterTargetTypeIsRequired       = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "search filter target type is required", pkgErr.InvalidTypeEmpty, "search filter target type"))
	statusSearchFilterIDIsEmpty                  = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "search filter ID is empty", pkgErr.InvalidTypeEmpty, "search filter ID"))
	statusSearchFilterIDNotFound                 = api.NewGRPCStatus(pkgErr.NewErrorNotFound(packageName, "search filter not found", "search filter"))
	statusInvalidListAPIKeyRequest               = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(packageName, "invalid list api key request", pkgErr.InvalidTypeEmpty, "list api key request"))
)
