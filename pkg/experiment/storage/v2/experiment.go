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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v2

import (
	"context"
	"errors"

	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

var (
	ErrExperimentAlreadyExists = pkgErr.NewErrorAlreadyExists(
		pkgErr.ExperimentPackageName,
		"already exists",
	)
	ErrExperimentNotFound = pkgErr.NewErrorNotFound(
		pkgErr.ExperimentPackageName,
		"not found",
		"experiment",
	)
	ErrExperimentUnexpectedAffectedRows = pkgErr.NewErrorUnexpectedAffectedRows(
		pkgErr.ExperimentPackageName,
		"unexpected affected rows",
	)
	ErrExperimentCannotBeArchived = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.ExperimentPackageName,
		"cannot be archived",
		"experiment_status",
	)
)

// Shared list-query errors returned by ExperimentStorage / GoalStorage implementations.
var (
	ErrInvalidCursor  = errors.New("experiment/storage/v2: invalid cursor")
	ErrInvalidOrderBy = errors.New("experiment/storage/v2: invalid order by")
)

type ExperimentStorage interface {
	CreateExperiment(ctx context.Context, e *domain.Experiment, environmentId string) error
	UpdateExperiment(ctx context.Context, e *domain.Experiment, environmentId string) error
	GetExperiment(ctx context.Context, id, environmentId string) (*domain.Experiment, error)
	ListExperiments(
		ctx context.Context,
		params ListExperimentsParams,
	) ([]*proto.Experiment, int, int64, error)
	// GetExperimentSummary returns the total count of experiments by status.
	GetExperimentSummary(ctx context.Context, environmentID string) (*ExperimentSummary, error)
}

type ListExperimentsParams struct {
	EnvironmentID  string
	Archived       *bool
	FeatureID      string
	FeatureVersion *int32
	StartAt        int64
	StopAt         int64
	Maintainer     string
	Statuses       []proto.Experiment_Status
	SearchKeyword  string
	OrderBy        proto.ListExperimentsRequest_OrderBy
	OrderDirection proto.ListExperimentsRequest_OrderDirection
	PageSize       int
	Cursor         string
}

type ExperimentSummary struct {
	TotalWaitingCount int64
	TotalRunningCount int64
	TotalStoppedCount int64
}
