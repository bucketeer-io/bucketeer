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
	statusInternal        = api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.AuditlogPackageName, "internal"))
	statusUnauthenticated = api.NewGRPCStatus(
		pkgErr.NewErrorUnauthenticated(pkgErr.AuditlogPackageName, "unauthenticated"))
	statusPermissionDenied = api.NewGRPCStatus(
		pkgErr.NewErrorPermissionDenied(pkgErr.AuditlogPackageName, "permission denied"))
	statusMissingID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuditlogPackageName, "missing ID", "ID"))
	statusNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.AuditlogPackageName, "not found", "auditlog"))
	statusInvalidCursor = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AuditlogPackageName, "cursor is invalid", "cursor"))
	statusInvalidOrderBy = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.AuditlogPackageName, "order_by is invalid", "order_by"))
)
