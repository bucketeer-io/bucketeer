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
//

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v2

import (
	"context"
	"errors"

	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"

	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var (
	ErrFlagTriggerAlreadyExists = pkgErr.NewErrorAlreadyExists(
		pkgErr.FeaturePackageName,
		"flag trigger already exists",
	)
	ErrFlagTriggerNotFound = pkgErr.NewErrorNotFound(
		pkgErr.FeaturePackageName,
		"flag trigger not found",
		"flag_trigger",
	)
	ErrFlagTriggerUnexpectedAffectedRows = pkgErr.NewErrorUnexpectedAffectedRows(
		pkgErr.FeaturePackageName,
		"flag trigger unexpected affected rows",
	)

	ErrInvalidListFlagTriggersCursor  = errors.New("flag trigger storage: invalid list flag triggers cursor")
	ErrInvalidListFlagTriggersOrderBy = errors.New("flag trigger storage: invalid list flag triggers order by")
)

// ListFlagTriggersParams carries list intent for ListFlagTriggers without database-specific types.
type ListFlagTriggersParams struct {
	FeatureID      string
	EnvironmentID  string
	OrderBy        proto.ListFlagTriggersRequest_OrderBy
	OrderDirection proto.ListFlagTriggersRequest_OrderDirection
	PageSize       int
	Cursor         string
}

type FlagTriggerStorage interface {
	CreateFlagTrigger(ctx context.Context, flagTrigger *domain.FlagTrigger) error
	UpdateFlagTrigger(ctx context.Context, flagTrigger *domain.FlagTrigger) error
	DeleteFlagTrigger(ctx context.Context, id, environmentId string) error
	GetFlagTrigger(ctx context.Context, id, environmentId string) (*domain.FlagTrigger, error)
	GetFlagTriggerByToken(ctx context.Context, token string) (*domain.FlagTrigger, error)
	ListFlagTriggers(
		ctx context.Context,
		params ListFlagTriggersParams,
	) ([]*proto.FlagTrigger, int, int64, error)
}
