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

const packageName = "insights"

var (
	statusInternal = api.NewGRPCStatus(
		pkgErr.NewErrorInternal(packageName, "internal error"))
	statusUnauthenticated = api.NewGRPCStatus(
		pkgErr.NewErrorUnauthenticated(packageName, "unauthenticated"))
	statusPermissionDenied = api.NewGRPCStatus(
		pkgErr.NewErrorPermissionDenied(packageName, "permission denied"))
	statusEnvironmentIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(packageName, "environment id is required", "EnvironmentId"))
	statusStartAtRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(packageName, "startAt is required", "StartAt"))
	statusEndAtRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(packageName, "endAt is required", "EndAt"))
	statusStartAtIsAfterEndAt = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(packageName, "startAt is after endAt", "StartAt"))
	statusQueryRangeTooLarge = api.NewGRPCStatus(
		pkgErr.NewErrorExceededMax(packageName, "query range exceeds the max days", "StartAt", maxQueryRangeDays))
	statusDataSourceNotConfigured = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(packageName, "data source is not configured", "DataSource"))
)
