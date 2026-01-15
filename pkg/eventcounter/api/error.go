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
		pkgErr.NewErrorInternal(pkgErr.EventCounterPackageName, "eventcounter: internal"))
	statusFeatureIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EventCounterPackageName, "feature id is required", "FeatureFlagID"))
	statusFeatureVersionRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(
			pkgErr.EventCounterPackageName, "feature flag version is required", "FeatureFlagVersion"))
	statusVariationIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EventCounterPackageName, "variation id is required", "VariationId"))
	statusExperimentIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EventCounterPackageName, "experiment id is required", "ExperimentId"))
	statusMAUYearMonthRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EventCounterPackageName, "mau year month is required", "MauYearMonth"))
	statusGoalIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EventCounterPackageName, "goal id is required", "Goal"))
	statusStartAtRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EventCounterPackageName, "start at is required", "StartAt"))
	statusEndAtRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EventCounterPackageName, "end at is required", "EndAt"))
	statusStartAtIsAfterEndAt = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.EventCounterPackageName, "start at is after end at", "StartAt"))
	statusAutoOpsRuleIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EventCounterPackageName, "auto ops rule id is required", "AutoOperationRuleId"))
	statusClauseIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.EventCounterPackageName, "clause id is required", "ClauseId"))
	statusExperimentResultNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.EventCounterPackageName, "experiment result not found", "ExperimentResult"))
	statusUnauthenticated = api.NewGRPCStatus(
		pkgErr.NewErrorUnauthenticated(pkgErr.EventCounterPackageName, "unauthenticated"))
	statusPermissionDenied = api.NewGRPCStatus(
		pkgErr.NewErrorPermissionDenied(pkgErr.EventCounterPackageName, "permission denied"))
	statusUnknownTimeRange = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.EventCounterPackageName, "unknown time range", "TimeRange"))
)
