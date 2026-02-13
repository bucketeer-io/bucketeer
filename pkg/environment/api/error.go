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
	statusInternal = api.NewGRPCStatus(
		pkgErr.NewErrorInternal(pkgErr.EnvironmentPackageName, "internal error"))
	statusNoCommand = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.EnvironmentPackageName, "no command", "Command"))
	statusInvalidCursor = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.EnvironmentPackageName, "cursor is invalid", "Cursor"))
	// Essentially, the id field is required, but no validation is performed because some older services do not have ID.
	//statusEnvironmentIDRequired = gstatus.New(codes.InvalidArgument, "environment: environment id must be specified")
	statusEnvironmentNameRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.EnvironmentPackageName,
			"environment name must be specified",
			"EnvironmentName",
		))
	statusInvalidEnvironmentName = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"invalid environment name",
			"EnvironmentName",
		))
	statusInvalidEnvironmentUrlCode = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"invalid environment url code",
			"EnvironmentUrlCode",
		))
	statusEnvironmentIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EnvironmentPackageName, "environment id must be specified", "environment_id"))
	statusProjectIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EnvironmentPackageName, "project id must be specified", "ProjectId"))
	statusProjectNameRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EnvironmentPackageName, "project name must be specified", "ProjectName"))
	statusInvalidProjectName = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.EnvironmentPackageName, "invalid project name", "ProjectName"))
	statusInvalidProjectUrlCode = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"invalid project url code",
			"ProjectUrlCode",
		))
	statusInvalidProjectCreatorEmail = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"invalid project creator email",
			"Email",
		))
	statusInvalidOrganizationCreatorEmail = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"invalid organization creator email",
			"Email",
		))
	statusInvalidOrderBy = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.EnvironmentPackageName, "order_by is invalid", "OrderBy"))
	statusOrganizationIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EnvironmentPackageName, "organization id must be specified", "OrganizationId"))
	statusOrganizationNameRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.EnvironmentPackageName,
			"organization name must be specified",
			"OrganizationName",
		))
	statusInvalidOrganizationName = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"invalid organization name",
			"OrganizationName",
		))
	statusInvalidOrganizationUrlCode = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(
			pkgErr.EnvironmentPackageName,
			"invalid organization url code",
			"OrganizationUrlCode",
		))
	statusCannotUpdateSystemAdmin = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(
			pkgErr.EnvironmentPackageName,
			"cannot update system admin organization",
		))
	statusCannotDeleteOrganization = api.NewGRPCStatus(
		pkgErr.NewErrorFailedPrecondition(
			pkgErr.EnvironmentPackageName,
			"cannot delete organization",
		))
	statusEnvironmentNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.EnvironmentPackageName, "environment not found", "Environment"))
	statusProjectNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.EnvironmentPackageName, "project not found", "Project"))
	statusOrganizationNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.EnvironmentPackageName, "organization not found", "Organization"))
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
	statusAccountNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.EnvironmentPackageName, "account not found", "Account"))
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
