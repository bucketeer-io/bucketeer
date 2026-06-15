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

	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

var (
	ErrGoalAlreadyExists          = errors.New("goal: already exists")
	ErrGoalNotFound               = errors.New("goal: not found")
	ErrGoalUnexpectedAffectedRows = errors.New("goal: unexpected affected rows")
)

type GoalStorage interface {
	CreateGoal(ctx context.Context, g *domain.Goal, environmentId string) error
	UpdateGoal(ctx context.Context, g *domain.Goal, environmentId string) error
	GetGoal(ctx context.Context, id, environmentId string) (*domain.Goal, error)
	ListGoals(
		ctx context.Context,
		params ListGoalsParams,
	) ([]*proto.Goal, int, int64, error)
	DeleteGoal(ctx context.Context, id, environmentId string) error
}

type ListGoalsParams struct {
	EnvironmentID  string
	Archived       *bool
	SearchKeyword  string
	ConnectionType proto.Goal_ConnectionType
	IsInUseStatus  *bool
	OrderBy        proto.ListGoalsRequest_OrderBy
	OrderDirection proto.ListGoalsRequest_OrderDirection
	PageSize       int
	Cursor         string
}
