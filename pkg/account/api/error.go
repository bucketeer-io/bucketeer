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
	"github.com/bucketeer-io/bucketeer/pkg/account"
	"github.com/bucketeer-io/bucketeer/pkg/api/api"
	pkgErr "github.com/bucketeer-io/bucketeer/pkg/error"
)

var (
	statusInternal                               = api.NewGRPCStatus(pkgErr.NewErrorInternal(account.PackageName, "internal"))
	statusInvalidCursor                          = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "cursor is invalid", pkgErr.InvalidTypeNotMatchFormat, "cursor"))
	statusNoCommand                              = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "command must not be empty", pkgErr.InvalidTypeEmpty, "command"))
	statusMissingOrganizationID                  = api.NewGRPCStatus(account.ErrMissingOrganizationID)
	statusEmailIsEmpty                           = api.NewGRPCStatus(account.ErrEmailIsEmpty)
	statusInvalidEmail                           = api.NewGRPCStatus(account.ErrEmailInvalidFormat)
	statusFirstNameIsEmpty                       = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "first name is empty", pkgErr.InvalidTypeEmpty, "first name"))
	statusInvalidFirstName                       = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "invalid first name format", pkgErr.InvalidTypeNotMatchFormat, "first name"))
	statusLastNameIsEmpty                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "last name is empty", pkgErr.InvalidTypeEmpty, "last name"))
	statusInvalidLastName                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "invalid last name format", pkgErr.InvalidTypeNotMatchFormat, "last name"))
	statusLanguageIsEmpty                        = api.NewGRPCStatus(account.ErrLanguageIsEmpty)
	statusInvalidOrganizationRole                = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "invalid organization roles", pkgErr.InvalidTypeEmpty, "organization role"))
	statusInvalidEnvironmentRole                 = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "environment role must be specified", pkgErr.InvalidTypeEmpty, "environment role"))
	statusInvalidUpdateEnvironmentRolesWriteType = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "update environment roles write type must be specified", pkgErr.InvalidTypeEmpty, "update environment roles write type"))
	statusMissingAPIKeyID                        = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "apikey id must be specified", pkgErr.InvalidTypeEmpty, "apikey id"))
	statusMissingAPIKeyName                      = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "apikey name must be not empty", pkgErr.InvalidTypeEmpty, "apikey name"))
	statusInvalidOrderBy                         = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "order_by is invalid", pkgErr.InvalidTypeNotMatchFormat, "order_by"))
	statusNotFound                               = api.NewGRPCStatus(account.ErrAccountNotFound)
	statusAlreadyExists                          = api.NewGRPCStatus(account.ErrAccountAlreadyExists)
	statusUnauthenticated                        = api.NewGRPCStatus(pkgErr.NewErrorUnauthenticated(account.PackageName, "account unauthenticated"))
	statusPermissionDenied                       = api.NewGRPCStatus(pkgErr.NewErrorPermissionDenied(account.PackageName, "permission denied"))
	statusSearchFilterNameIsEmpty                = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "search filter name is empty", pkgErr.InvalidTypeEmpty, "search filter name"))
	statusSearchFilterQueryIsEmpty               = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "search filter query is empty", pkgErr.InvalidTypeEmpty, "search filter query"))
	statusSearchFilterTargetTypeIsRequired       = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "search filter target type is required", pkgErr.InvalidTypeEmpty, "search filter target type"))
	statusSearchFilterIDIsEmpty                  = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "search filter ID is empty", pkgErr.InvalidTypeEmpty, "search filter ID"))
	statusSearchFilterIDNotFound                 = api.NewGRPCStatus(account.ErrSearchFilterNotFound)
	statusInvalidListAPIKeyRequest               = api.NewGRPCStatus(pkgErr.NewErrorInvalidAugment(account.PackageName, "invalid list api key request", pkgErr.InvalidTypeEmpty, "list api key request"))
)
