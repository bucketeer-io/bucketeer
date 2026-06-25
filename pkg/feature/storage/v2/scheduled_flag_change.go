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
	"time"

	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var (
	ErrScheduledFlagChangeAlreadyExists = pkgErr.NewErrorAlreadyExists(
		pkgErr.FeaturePackageName,
		"scheduled flag change already exists",
	)
	ErrScheduledFlagChangeNotFound = pkgErr.NewErrorNotFound(
		pkgErr.FeaturePackageName,
		"scheduled flag change not found",
		"scheduled_flag_change",
	)
	ErrScheduledFlagChangeUnexpectedAffectedRows = pkgErr.NewErrorUnexpectedAffectedRows(
		pkgErr.FeaturePackageName,
		"scheduled flag change unexpected affected rows",
	)
)

// LockExpiration is the time after which a lock is considered expired (5 minutes)
const LockExpiration = 5 * time.Minute

// ListScheduledFlagChangesParams carries list intent for ListScheduledFlagChanges
// without database-specific types.
type ListScheduledFlagChangesParams struct {
	EnvironmentID string
	FeatureID     string
	// ExcludeFeatureID, when set, matches schedules whose feature_id is NOT this value.
	ExcludeFeatureID string
	FromScheduledAt  int64
	ToScheduledAt    int64
	Statuses         []proto.ScheduledFlagChangeStatus
	OrderBy          proto.ListScheduledFlagChangesRequest_OrderBy
	OrderDirection   proto.ListScheduledFlagChangesRequest_OrderDirection
	// PageSize is the row limit; use database.QueryNoLimit for an uncapped list.
	PageSize int
	Offset   int
}

// ScheduledFlagChangeStorage defines the interface for scheduled flag change storage operations
type ScheduledFlagChangeStorage interface {
	// CreateScheduledFlagChange creates a new scheduled flag change
	CreateScheduledFlagChange(ctx context.Context, sfc *domain.ScheduledFlagChange) error
	// UpdateScheduledFlagChange updates an existing scheduled flag change
	UpdateScheduledFlagChange(ctx context.Context, sfc *domain.ScheduledFlagChange) error
	// DeleteScheduledFlagChange deletes a scheduled flag change by ID
	DeleteScheduledFlagChange(ctx context.Context, id, environmentID string) error
	// GetScheduledFlagChange retrieves a scheduled flag change by ID
	GetScheduledFlagChange(ctx context.Context, id, environmentID string) (*domain.ScheduledFlagChange, error)
	// ListScheduledFlagChanges lists scheduled flag changes with filtering and pagination
	ListScheduledFlagChanges(
		ctx context.Context,
		params ListScheduledFlagChangesParams,
	) ([]*proto.ScheduledFlagChange, int, int64, error)
	// ListDueScheduledFlagChanges lists scheduled flag changes that are due for execution
	ListDueScheduledFlagChanges(ctx context.Context, now int64, limit int) ([]*proto.ScheduledFlagChange, error)
	// TryLock attempts to acquire a lock on a scheduled flag change for execution
	TryLock(ctx context.Context, id, lockedBy string) (bool, error)
	// Unlock releases the lock on a scheduled flag change (only if locked by the same executor)
	Unlock(ctx context.Context, id, lockedBy string) error
}
