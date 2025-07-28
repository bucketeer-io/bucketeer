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
	statusInternal                               = api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.AccountPackageName, "internal", nil))
	statusInvalidCursor                          = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "cursor is invalid", pkgErr.InvalidTypeNotMatchFormat, nil, "cursor"))
	statusNoCommand                              = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "command must not be empty", pkgErr.InvalidTypeEmpty, nil, "command"))
	statusMissingOrganizationID                  = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "organization id must be specified", pkgErr.InvalidTypeEmpty, nil, "organization_id"))
	statusEmailIsEmpty                           = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "email is empty", pkgErr.InvalidTypeEmpty, nil, "email"))
	statusInvalidEmail                           = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "invalid email format", pkgErr.InvalidTypeNotMatchFormat, nil, "email"))
	statusFirstNameIsEmpty                       = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "first name is empty", pkgErr.InvalidTypeEmpty, nil, "first name"))
	statusInvalidFirstName                       = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "invalid first name format", pkgErr.InvalidTypeNotMatchFormat, nil, "first name"))
	statusLastNameIsEmpty                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "last name is empty", pkgErr.InvalidTypeEmpty, nil, "last name"))
	statusInvalidLastName                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "invalid last name format", pkgErr.InvalidTypeNotMatchFormat, nil, "last name"))
	statusLanguageIsEmpty                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "language is empty", pkgErr.InvalidTypeEmpty, nil, "language"))
	statusInvalidOrganizationRole                = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "invalid organization roles", pkgErr.InvalidTypeEmpty, nil, "organization role"))
	statusInvalidEnvironmentRole                 = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "environment role must be specified", pkgErr.InvalidTypeEmpty, nil, "environment role"))
	statusInvalidUpdateEnvironmentRolesWriteType = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "update environment roles write type must be specified", pkgErr.InvalidTypeEmpty, nil, "update environment roles write type"))
	statusMissingAPIKeyID                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "apikey id must be specified", pkgErr.InvalidTypeEmpty, nil, "apikey id"))
	statusMissingAPIKeyName                      = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "apikey name must be not empty", pkgErr.InvalidTypeEmpty, nil, "apikey name"))
	statusInvalidOrderBy                         = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "order_by is invalid", pkgErr.InvalidTypeNotMatchFormat, nil, "order_by"))
	statusNotFound                               = api.NewGRPCStatus(pkgErr.NewErrorNotFound(pkgErr.AccountPackageName, "account not found", nil, "account"))
	statusAlreadyExists                          = api.NewGRPCStatus(pkgErr.NewErrorAlreadyExists(pkgErr.AccountPackageName, "account already exists", nil, "account"))
	statusUnauthenticated                        = api.NewGRPCStatus(pkgErr.NewErrorUnauthenticated(pkgErr.AccountPackageName, "account unauthenticated", nil))
	statusPermissionDenied                       = api.NewGRPCStatus(pkgErr.NewErrorPermissionDenied(pkgErr.AccountPackageName, "permission denied", nil))
	statusSearchFilterNameIsEmpty                = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "search filter name is empty", pkgErr.InvalidTypeEmpty, nil, "search filter name"))
	statusSearchFilterQueryIsEmpty               = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "search filter query is empty", pkgErr.InvalidTypeEmpty, nil, "search filter query"))
	statusSearchFilterTargetTypeIsRequired       = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "search filter target type is required", pkgErr.InvalidTypeEmpty, nil, "search filter target type"))
	statusSearchFilterIDIsEmpty                  = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "search filter ID is empty", pkgErr.InvalidTypeEmpty, nil, "search filter ID"))
	statusSearchFilterIDNotFound                 = api.NewGRPCStatus(pkgErr.NewErrorNotFound(pkgErr.AccountPackageName, "search filter not found", nil, "search_filter"))
	statusInvalidListAPIKeyRequest               = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(pkgErr.AccountPackageName, "invalid list api key request", pkgErr.InvalidTypeEmpty, nil, "list api key request"))
)
