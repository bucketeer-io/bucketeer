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
	statusInvalidCursor = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.ExperimentPackageName, "cursor is invalid", "Cursor"))
	statusFeatureIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.ExperimentPackageName, "feature id must be specified", "FeatureFlagID"))
	statusExperimentIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.ExperimentPackageName, "experiment id must be specified", "Experiment"))
	statusExperimentNameRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.ExperimentPackageName, "experiment name must be specified", "Experiment"))
	statusGoalIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.ExperimentPackageName, "goal id must be specified", "Goal"))
	statusGoalTypeMismatch = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.ExperimentPackageName, "goal type mismatch", "Goal"))
	statusInvalidGoalID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.ExperimentPackageName, "invalid goal id", "Goal"))
	statusGoalNameRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.ExperimentPackageName, "goal name must be specified", "Goal"))
	statusExperimentPeriodOutOfRange = api.NewGRPCStatus(
		pkgErr.NewErrorOutOfRange(
			pkgErr.ExperimentPackageName,
			"period too long",
			"ExperimentPeriod",
			0,
			maxExperimentPeriod,
		))
	statusExperimentPeriodInvalid = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.ExperimentPackageName, "period is invalid", "ExperimentPeriod"))
	statusInvalidOrderBy = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.ExperimentPackageName, "order_by is invalid", "OrderBy"))
	statusExperimentNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.ExperimentPackageName, "experiment not found", "Experiment"))
	statusGoalNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.ExperimentPackageName, "goal not found", "Goal"))
	statusFeatureNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.ExperimentPackageName, "feature not found", "FeatureFlag"))
	statusAlreadyExists = api.NewGRPCStatus(
		pkgErr.NewErrorAlreadyExists(pkgErr.ExperimentPackageName, "already exists"))
	statusUnauthenticated = api.NewGRPCStatus(
		pkgErr.NewErrorUnauthenticated(pkgErr.ExperimentPackageName, "unauthenticated"))
	statusPermissionDenied = api.NewGRPCStatus(
		pkgErr.NewErrorPermissionDenied(pkgErr.ExperimentPackageName, "permission denied"))
)
