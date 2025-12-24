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

var (
	statusInternal  = api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal"))
	statusNoCommand = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.EnvironmentPackageName, "no command", "command"))
	statusInvalidCursor = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.EnvironmentPackageName, "cursor is invalid", "cursor"))
	// Essentially, the id field is required, but no validation is performed because some older services do not have ID.
	//statusEnvironmentIDRequired = gstatus.New(codes.InvalidArgument, "environment: environment id must be specified")
	statusEnvironmentNameRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.EnvironmentPackageName,
			"environment name must be specified",
			"environment_name",
		))
	statusInvalidEnvironmentName = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"invalid environment name",
			"environment_name",
		))
	statusInvalidEnvironmentUrlCode = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.EnvironmentPackageName,
			"invalid environment url code",
			"environment_url_code",
		))
	statusProjectIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EnvironmentPackageName, "project id must be specified", "project_id"))
	statusProjectNameRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EnvironmentPackageName, "project name must be specified", "project_name"))
	statusInvalidProjectName = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.EnvironmentPackageName, "invalid project name", "project_name"))
	statusInvalidProjectUrlCode = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"invalid project url code",
			"project_url_code",
		))
	statusInvalidProjectCreatorEmail = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"invalid project creator email",
			"project_creator_email",
		))
	statusInvalidOrganizationCreatorEmail = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"invalid organization creator email",
			"organization_creator_email",
		))
	statusInvalidOrderBy = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.EnvironmentPackageName, "order_by is invalid", "order_by"))
	statusOrganizationIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EnvironmentPackageName, "organization id must be specified", "organization_id"))
	statusOrganizationNameRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.EnvironmentPackageName,
			"organization name must be specified",
			"organization_name",
		))
	statusInvalidOrganizationName = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"invalid organization name",
			"organization_name",
		))
	statusInvalidOrganizationUrlCode = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"invalid organization url code",
			"organization_url_code",
		))
	statusCannotUpdateSystemAdmin = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"cannot update system admin organization",
			"system_admin_organization",
		))
	statusEnvironmentNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.EnvironmentPackageName, "environment not found", "environment"))
	statusProjectNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.EnvironmentPackageName, "project not found", "project"))
	statusOrganizationNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.EnvironmentPackageName, "organization not found", "organization"))
	statusEnvironmentAlreadyExists = api.NewGRPCStatus(
		pkgErr.NewErrorAlreadyExists(pkgErr.EnvironmentPackageName, "environment already exists"))
	statusProjectAlreadyExists = api.NewGRPCStatus(
		pkgErr.NewErrorAlreadyExists(pkgErr.EnvironmentPackageName, "project already exists"))
	statusOrganizationAlreadyExists = api.NewGRPCStatus(
		pkgErr.NewErrorAlreadyExists(pkgErr.EnvironmentPackageName, "organization already exists"))
	statusProjectDisabled = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(pkgErr.EnvironmentPackageName, "project disabled"))
	statusUnauthenticated = api.NewGRPCStatus(
		pkgErr.NewErrorUnauthenticated(pkgErr.EnvironmentPackageName, "unauthenticated"))
	statusPermissionDenied = api.NewGRPCStatus(
		pkgErr.NewErrorPermissionDenied(pkgErr.EnvironmentPackageName, "permission denied"))
	statusNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.EnvironmentPackageName, "not found", "account"))
	statusDemoSiteDisabled = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(pkgErr.EnvironmentPackageName, "demo site is not enabled"))
	statusUserAlreadyInOrganization = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(pkgErr.EnvironmentPackageName, "user already in organization"))
	statusInvalidAutoArchiveUnusedDays = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"auto_archive_unused_days must be greater than 0 when auto_archive is enabled",
			"auto_archive_unused_days",
		))
	statusAutoArchiveNotEnabled = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"cannot update auto-archive settings when auto_archive_enabled is false",
			"auto_archive_settings",
		))
)
