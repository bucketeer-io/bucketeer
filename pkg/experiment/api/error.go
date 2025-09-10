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

var (
	statusInternal = api.NewGRPCStatus(
		pkgErr.NewErrorInternal(pkgErr.ExperimentPackageName, "internal error"))
	statusInvalidCursor = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.ExperimentPackageName, "cursor is invalid", "cursor"))
	statusNoCommand = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNil(pkgErr.ExperimentPackageName, "must contain at least one command", "command"))
	statusFeatureIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.ExperimentPackageName, "feature id must be specified", "feature_id"))
	statusExperimentIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.ExperimentPackageName, "experiment id must be specified", "experiment_id"))
	statusExperimentNameRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.ExperimentPackageName, "experiment name must be specified", "experiment_name"))
	statusGoalIDRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.ExperimentPackageName, "goal id must be specified", "goal_id"))
	statusGoalTypeMismatch = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.ExperimentPackageName, "goal type mismatch", "goal_type"))
	statusInvalidGoalID = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.ExperimentPackageName, "invalid goal id", "goal_id"))
	statusGoalNameRequired = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgEmpty(pkgErr.ExperimentPackageName, "goal name must be specified", "goal_name"))
	statusPeriodTooLong = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.ExperimentPackageName, "period too long", "period"))
	statusPeriodInvalid = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.ExperimentPackageName, "period is invalid", "period"))
	statusInvalidOrderBy = api.NewGRPCStatus(
		pkgErr.NewErrorInvalidArgNotMatchFormat(pkgErr.ExperimentPackageName, "order_by is invalid", "order_by"))
	statusNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.ExperimentPackageName, "not found", "experiment"))
	statusGoalNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.ExperimentPackageName, "goal not found", "goal_id"))
	statusFeatureNotFound = api.NewGRPCStatus(
		pkgErr.NewErrorNotFound(pkgErr.ExperimentPackageName, "feature not found", "feature"))
	statusAlreadyExists = api.NewGRPCStatus(
		pkgErr.NewErrorAlreadyExists(pkgErr.ExperimentPackageName, "already exists"))
	statusUnauthenticated = api.NewGRPCStatus(
		pkgErr.NewErrorUnauthenticated(pkgErr.ExperimentPackageName, "unauthenticated"))
	statusPermissionDenied = api.NewGRPCStatus(
		pkgErr.NewErrorPermissionDenied(pkgErr.ExperimentPackageName, "permission denied"))
)
