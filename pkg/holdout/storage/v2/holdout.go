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
	"errors"
	"fmt"
	"strconv"
	"time"

	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	experimentdomain "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
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
	// GetClearedHoldoutUpdatedAt is the updated_at of the newest STOPPED or
	// FORCE_STOPPED group, including archived ones. The empty cache snapshot
	// stores it so a client whose evaluatedAt is still behind the stop gets a
	// full re-evaluation. Returns 0 when the environment has never stopped a group.
	GetClearedHoldoutUpdatedAt(ctx context.Context, environmentID string) (int64, error)
	// CreateSegmentUserSnapshot copies the users the segment includes into
	// holdout_group_segment_user under the group, in one statement that runs
	// in the database, and returns how many users it copied. The users are
	// not loaded into the service.
	CreateSegmentUserSnapshot(ctx context.Context, groupID, segmentID, environmentID string) (int64, error)
	// CountSegmentUsers counts the users CreateSegmentUserSnapshot would copy
	// for the segment and stops when it has counted limit of them: a result
	// equal to limit means limit or more. Rows of the segment that are not
	// counted, excluded or deleted users, are still read on the way; no index
	// covers state and deleted.
	CountSegmentUsers(ctx context.Context, segmentID, environmentID string, limit int64) (int64, error)
	// ListLiveFeatureTargets returns current target rows of WAITING, RUNNING, or ANALYSIS groups.
	ListLiveFeatureTargets(ctx context.Context, environmentID, featureID string) ([]domain.FeatureTarget, error)
	// ListMeasuredGroups returns the groups that have an internal experiment
	// and are serving, or stopped at endedSince or later, archived ones
	// included. A scheduled stop is dated by stop_at and a manual stop by
	// stopped_at, which is the stop_at of the internal experiments of the
	// group in each case. Every target row of a group is attached, removed
	// ones included, so a reader can tell which flags the group targeted at
	// a past time. A group without an internal experiment is never returned.
	// The order is start_at, then id; target rows are in added_at, then id
	// order.
	ListMeasuredGroups(ctx context.Context, environmentID string, endedSince int64) ([]*domain.HoldoutGroup, error)
}

// MeasuredGroups is the proto stored in the measured groups cache: the
// groups of the environment whose internal experiments the analytics
// pipeline still lists at now, which is until StoppedExperimentEventWindow
// after they stopped. Every writer (batch cacher, subscriber) builds the
// list here so they publish the same one.
func MeasuredGroups(
	ctx context.Context,
	storage HoldoutStorage,
	environmentID string,
	now time.Time,
) (*proto.HoldoutGroups, error) {
	groups, err := storage.ListMeasuredGroups(
		ctx,
		environmentID,
		now.Add(-experimentdomain.StoppedExperimentEventWindow).Unix(),
	)
	if err != nil {
		return nil, err
	}
	list := &proto.HoldoutGroups{HoldoutGroups: make([]*proto.HoldoutGroup, 0, len(groups))}
	for _, group := range groups {
		list.HoldoutGroups = append(list.HoldoutGroups, group.HoldoutGroup)
	}
	return list, nil
}

// ServingSnapshot is the proto stored in the holdout cache.
// A missing group becomes an empty snapshot whose updated_at is the newest
// stop time, so every writer (batch cacher, feature service, subscriber)
// publishes the same revision.
func ServingSnapshot(
	ctx context.Context,
	storage HoldoutStorage,
	environmentID string,
) (*proto.HoldoutGroup, error) {
	group, err := storage.GetLiveServingGroup(ctx, environmentID)
	if err != nil {
		if errors.Is(err, ErrHoldoutGroupNotFound) {
			return clearedServingSnapshot(ctx, storage, environmentID)
		}
		return nil, err
	}
	if group == nil || group.HoldoutGroup == nil {
		return clearedServingSnapshot(ctx, storage, environmentID)
	}
	return group.HoldoutGroup, nil
}

func clearedServingSnapshot(
	ctx context.Context,
	storage HoldoutStorage,
	environmentID string,
) (*proto.HoldoutGroup, error) {
	updatedAt, err := storage.GetClearedHoldoutUpdatedAt(ctx, environmentID)
	if err != nil {
		return nil, err
	}
	return &proto.HoldoutGroup{UpdatedAt: updatedAt}, nil
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
