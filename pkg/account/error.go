// Copyright 2025 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package account

import (
	pkgErr "github.com/bucketeer-io/bucketeer/pkg/error"
)

const PackageName = "account"

var (
	ErrEmailIsEmpty            = pkgErr.NewErrorInvalidAugment(PackageName, "email is empty", pkgErr.InvalidTypeEmpty, "email")
	ErrMissingOrganizationID   = pkgErr.NewErrorInvalidAugment(PackageName, "organization id must be specified", pkgErr.InvalidTypeEmpty, "organization_id")
	ErrEmailInvalidFormat      = pkgErr.NewErrorInvalidAugment(PackageName, "invalid email format", pkgErr.InvalidTypeNotMatchFormat, "email")
	ErrFullNameIsEmpty         = pkgErr.NewErrorInvalidAugment(PackageName, "full name is empty", pkgErr.InvalidTypeEmpty, "full_name")
	ErrFullNameInvalidFormat   = pkgErr.NewErrorInvalidAugment(PackageName, "invalid full name format", pkgErr.InvalidTypeNotMatchFormat, "full_name")
	ErrLastNameInvalidFormat   = pkgErr.NewErrorInvalidAugment(PackageName, "invalid last name format", pkgErr.InvalidTypeNotMatchFormat, "last_name")
	ErrLanguageIsEmpty         = pkgErr.NewErrorInvalidAugment(PackageName, "language is empty", pkgErr.InvalidTypeEmpty, "language")
	ErrOrganizationRoleInvalid = pkgErr.NewErrorInvalidAugment(PackageName, "invalid organization role", pkgErr.InvalidTypeEmpty, "organization_role")
	ErrSearchFilterNotFound    = pkgErr.NewErrorNotFound(PackageName, "search filter not found", "search_filter")
	ErrTeamNotFound            = pkgErr.NewErrorNotFound(PackageName, "team not found", "team")
)
