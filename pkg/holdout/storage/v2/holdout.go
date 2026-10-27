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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/holdout.go
package v2

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strconv"

	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/holdout/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var (
	ErrHoldoutGroupNotFound = pkgErr.NewErrorNotFound(
		pkgErr.HoldoutPackageName, "holdout group not found", "holdout_group",
	)
	ErrHoldoutGroupAlreadyExists = pkgErr.NewErrorAlreadyExists(
		pkgErr.HoldoutPackageName, "holdout group already exists",
	)
	ErrHoldoutGroupUnexpectedAffectedRows = pkgErr.NewErrorUnexpectedAffectedRows(
		pkgErr.HoldoutPackageName, "unexpected affected rows",
	)
	ErrHoldoutFeatureNotFound = pkgErr.NewErrorNotFound(
		pkgErr.HoldoutPackageName, "holdout group feature not found", "feature_id",
	)
	ErrEnvironmentNotFound = pkgErr.NewErrorNotFound(
		pkgErr.HoldoutPackageName, "environment not found", "environment",
	)
	ErrInvalidCursor = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.HoldoutPackageName, "cursor is invalid", "cursor",
	)
	ErrInvalidOrderBy = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.HoldoutPackageName, "order_by is invalid", "order_by",
	)
)

type HoldoutStorage interface {
	CreateHoldoutGroup(ctx context.Context, group *domain.HoldoutGroup) error
	UpdateHoldoutGroup(ctx context.Context, group *domain.HoldoutGroup) error
	GetHoldoutGroup(ctx context.Context, id, environmentID string) (*domain.HoldoutGroup, error)
	GetHoldoutGroupForUpdate(ctx context.Context, id, environmentID string) (*domain.HoldoutGroup, error)
	ListHoldoutGroups(ctx context.Context, params ListHoldoutGroupsParams) ([]*domain.HoldoutGroup, int, int64, error)
	CreateHoldoutGroupFeatures(ctx context.Context, features []*domain.HoldoutGroupFeature) error
	RemoveHoldoutGroupFeature(
		ctx context.Context,
		environmentID, groupID, featureID string,
		removedAt int64,
	) error
	ListCurrentFeatures(
		ctx context.Context,
		environmentID string,
		groupIDs []string,
	) ([]*domain.HoldoutGroupFeature, error)
	// LockEnvironment locks the environment row so concurrent holdout writes serialize.
	LockEnvironment(ctx context.Context, environmentID string) error
	// GetEnvironmentSamplingSeed returns the canonical HASH seed, or empty when the environment has none.
	GetEnvironmentSamplingSeed(ctx context.Context, environmentID string) (string, error)
	ListActiveGroups(ctx context.Context, environmentID string) ([]domain.GroupSnapshot, error)
	// GetLiveServingGroup returns the newest RUNNING or ANALYSIS group that is
	// not deleted or archived, with its current target features attached.
	// ErrHoldoutGroupNotFound when the environment has no serving group.
	GetLiveServingGroup(ctx context.Context, environmentID string) (*domain.HoldoutGroup, error)
	// ListLiveFeatureTargets returns current target rows of WAITING, RUNNING, or ANALYSIS groups.
	ListLiveFeatureTargets(ctx context.Context, environmentID, featureID string) ([]domain.FeatureTarget, error)
}

type ListHoldoutGroupsParams struct {
	EnvironmentID  string
	Archived       *bool
	Statuses       []proto.HoldoutGroup_Status
	SearchKeyword  string
	OrderBy        proto.ListHoldoutGroupsRequest_OrderBy
	OrderDirection proto.ListHoldoutGroupsRequest_OrderDirection
	PageSize       int
	Cursor         string
}

// Decimal scans DECIMAL columns into a float64 with three fractional digits on write.
type Decimal float64

func (d *Decimal) Scan(src interface{}) error {
	if src == nil {
		*d = 0
		return nil
	}
	switch v := src.(type) {
	case float64:
		*d = Decimal(v)
	case float32:
		*d = Decimal(v)
	case int64:
		*d = Decimal(v)
	case []byte:
		parsed, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return err
		}
		*d = Decimal(parsed)
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		*d = Decimal(parsed)
	default:
		return fmt.Errorf("holdout: cannot scan %T into decimal", src)
	}
	return nil
}

func (d Decimal) Value() (driver.Value, error) {
	return strconv.FormatFloat(float64(d), 'f', 3, 64), nil
}

func NormalizeGoalIDs(ids []string) []string {
	if ids == nil {
		return []string{}
	}
	return ids
}
