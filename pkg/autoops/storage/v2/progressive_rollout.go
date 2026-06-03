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
	"errors"

	"context"

	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	err "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

var (
	ErrProgressiveRolloutAlreadyExists = err.NewErrorAlreadyExists(err.AutoopsPackageName, "already exists")
	ErrProgressiveRolloutNotFound      = err.NewErrorNotFound(
		err.AutoopsPackageName,
		"not found",
		"progressive_rollout",
	)
	ErrProgressiveRolloutUnexpectedAffectedRows = err.NewErrorUnexpectedAffectedRows(
		err.AutoopsPackageName,
		"unexpected affected rows",
	)
)

// Shared list-query errors returned by ProgressiveRolloutStorage implementations.
var ErrInvalidOrderBy = errors.New("autoops/storage/v2: invalid order by")

type ProgressiveRolloutStorage interface {
	CreateProgressiveRollout(
		ctx context.Context,
		progressiveRollout *domain.ProgressiveRollout,
		environmentId string,
	) error
	GetProgressiveRollout(ctx context.Context, id, environmentId string) (*domain.ProgressiveRollout, error)
	DeleteProgressiveRollout(ctx context.Context, id, environmentId string) error
	ListProgressiveRollouts(
		ctx context.Context,
		params ListProgressiveRolloutsParams,
	) ([]*autoopsproto.ProgressiveRollout, int64, int, error)
	UpdateProgressiveRollout(ctx context.Context,
		progressiveRollout *domain.ProgressiveRollout,
		environmentId string,
	) error
}

type ListProgressiveRolloutsParams struct {
	EnvironmentID  string
	FeatureIDs     []string
	Type           *autoopsproto.ProgressiveRollout_Type
	Status         *autoopsproto.ProgressiveRollout_Status
	OrderBy        autoopsproto.ListProgressiveRolloutsRequest_OrderBy
	OrderDirection autoopsproto.ListProgressiveRolloutsRequest_OrderDirection
	PageSize       int
	Cursor         string
}
