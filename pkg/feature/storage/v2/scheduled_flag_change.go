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
	_ "embed"
	"errors"
	"time"

	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var (
	//go:embed sql/scheduled_feature_change/insert_scheduled_feature_change.sql
	insertScheduledFlagChangeSQL string
	//go:embed sql/scheduled_feature_change/update_scheduled_feature_change.sql
	updateScheduledFlagChangeSQL string
	//go:embed sql/scheduled_feature_change/delete_scheduled_feature_change.sql
	deleteScheduledFlagChangeSQL string
	//go:embed sql/scheduled_feature_change/get_scheduled_feature_change.sql
	getScheduledFlagChangeSQL string
	//go:embed sql/scheduled_feature_change/list_scheduled_feature_changes.sql
	listScheduledFlagChangesSQL string
	//go:embed sql/scheduled_feature_change/count_scheduled_feature_changes.sql
	countScheduledFlagChangesSQL string
	//go:embed sql/scheduled_feature_change/list_due_scheduled_feature_changes.sql
	listDueScheduledFlagChangesSQL string
	//go:embed sql/scheduled_feature_change/try_lock_scheduled_feature_change.sql
	tryLockScheduledFlagChangeSQL string
	//go:embed sql/scheduled_feature_change/unlock_scheduled_feature_change.sql
	unlockScheduledFlagChangeSQL string
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
	ListScheduledFlagChanges(ctx context.Context, options *mysql.ListOptions) ([]*proto.ScheduledFlagChange, int, int64, error)
	// ListDueScheduledFlagChanges lists scheduled flag changes that are due for execution
	ListDueScheduledFlagChanges(ctx context.Context, now int64, limit int) ([]*proto.ScheduledFlagChange, error)
	// TryLock attempts to acquire a lock on a scheduled flag change for execution
	TryLock(ctx context.Context, id, lockedBy string) (bool, error)
	// Unlock releases the lock on a scheduled flag change
	Unlock(ctx context.Context, id string) error
}

type scheduledFlagChangeStorage struct {
	qe mysql.QueryExecer
}

// NewScheduledFlagChangeStorage creates a new ScheduledFlagChangeStorage
func NewScheduledFlagChangeStorage(qe mysql.QueryExecer) ScheduledFlagChangeStorage {
	return &scheduledFlagChangeStorage{qe: qe}
}

func (s *scheduledFlagChangeStorage) CreateScheduledFlagChange(
	ctx context.Context,
	sfc *domain.ScheduledFlagChange,
) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertScheduledFlagChangeSQL,
		sfc.Id,
		sfc.FeatureId,
		sfc.EnvironmentId,
		sfc.ScheduledAt,
		sfc.Timezone,
		mysql.JSONObject{Val: sfc.Payload},
		sfc.Comment,
		int32(sfc.Status),
		sfc.FailureReason,
		sfc.FlagVersionAtCreation,
		mysql.JSONObject{Val: sfc.Conflicts},
		sfc.CreatedBy,
		sfc.CreatedAt,
		sfc.UpdatedBy,
		sfc.UpdatedAt,
		sfc.ExecutedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrDuplicateEntry) {
			return ErrScheduledFlagChangeAlreadyExists
		}
		return err
	}
	return nil
}

func (s *scheduledFlagChangeStorage) UpdateScheduledFlagChange(
	ctx context.Context,
	sfc *domain.ScheduledFlagChange,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateScheduledFlagChangeSQL,
		sfc.ScheduledAt,
		sfc.Timezone,
		mysql.JSONObject{Val: sfc.Payload},
		sfc.Comment,
		int32(sfc.Status),
		sfc.FailureReason,
		mysql.JSONObject{Val: sfc.Conflicts},
		sfc.UpdatedBy,
		sfc.UpdatedAt,
		sfc.ExecutedAt,
		sfc.Id,
		sfc.EnvironmentId,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrScheduledFlagChangeUnexpectedAffectedRows
	}
	return nil
}

func (s *scheduledFlagChangeStorage) DeleteScheduledFlagChange(
	ctx context.Context,
	id, environmentID string,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		deleteScheduledFlagChangeSQL,
		id,
		environmentID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrScheduledFlagChangeUnexpectedAffectedRows
	}
	return nil
}

func (s *scheduledFlagChangeStorage) GetScheduledFlagChange(
	ctx context.Context,
	id, environmentID string,
) (*domain.ScheduledFlagChange, error) {
	sfc := proto.ScheduledFlagChange{}
	var status int32
	err := s.qe.QueryRowContext(
		ctx,
		getScheduledFlagChangeSQL,
		id,
		environmentID,
	).Scan(
		&sfc.Id,
		&sfc.FeatureId,
		&sfc.EnvironmentId,
		&sfc.ScheduledAt,
		&sfc.Timezone,
		&mysql.JSONObject{Val: &sfc.Payload},
		&sfc.Comment,
		&status,
		&sfc.FailureReason,
		&sfc.FlagVersionAtCreation,
		&mysql.JSONObject{Val: &sfc.Conflicts},
		&sfc.CreatedBy,
		&sfc.CreatedAt,
		&sfc.UpdatedBy,
		&sfc.UpdatedAt,
		&sfc.ExecutedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrScheduledFlagChangeNotFound
		}
		return nil, err
	}
	sfc.Status = proto.ScheduledFlagChangeStatus(status)
	return &domain.ScheduledFlagChange{ScheduledFlagChange: &sfc}, nil
}

func (s *scheduledFlagChangeStorage) ListScheduledFlagChanges(
	ctx context.Context,
	options *mysql.ListOptions,
) ([]*proto.ScheduledFlagChange, int, int64, error) {
	query, whereArgs := mysql.ConstructQueryAndWhereArgs(listScheduledFlagChangesSQL, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var limit, offset int
	if options != nil {
		limit = options.Limit
		offset = options.Offset
	}
	scheduledFlagChanges := make([]*proto.ScheduledFlagChange, 0, limit)
	for rows.Next() {
		sfc := proto.ScheduledFlagChange{}
		var status int32
		err := rows.Scan(
			&sfc.Id,
			&sfc.FeatureId,
			&sfc.EnvironmentId,
			&sfc.ScheduledAt,
			&sfc.Timezone,
			&mysql.JSONObject{Val: &sfc.Payload},
			&sfc.Comment,
			&status,
			&sfc.FailureReason,
			&sfc.FlagVersionAtCreation,
			&mysql.JSONObject{Val: &sfc.Conflicts},
			&sfc.CreatedBy,
			&sfc.CreatedAt,
			&sfc.UpdatedBy,
			&sfc.UpdatedAt,
			&sfc.ExecutedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		sfc.Status = proto.ScheduledFlagChangeStatus(status)
		scheduledFlagChanges = append(scheduledFlagChanges, &sfc)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	nextOffset := offset + len(scheduledFlagChanges)
	var totalCount int64
	countQuery, countWhereArgs := mysql.ConstructCountQuery(countScheduledFlagChangesSQL, options)
	if err := s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount); err != nil {
		return nil, 0, 0, err
	}
	return scheduledFlagChanges, nextOffset, totalCount, nil
}

func (s *scheduledFlagChangeStorage) ListDueScheduledFlagChanges(
	ctx context.Context,
	now int64,
	limit int,
) ([]*proto.ScheduledFlagChange, error) {
	// Lock expiration time (5 minutes ago)
	lockExpiredAt := now - int64(LockExpiration.Seconds())

	rows, err := s.qe.QueryContext(
		ctx,
		listDueScheduledFlagChangesSQL,
		now,
		lockExpiredAt,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scheduledFlagChanges := make([]*proto.ScheduledFlagChange, 0, limit)
	for rows.Next() {
		sfc := proto.ScheduledFlagChange{}
		var status int32
		err := rows.Scan(
			&sfc.Id,
			&sfc.FeatureId,
			&sfc.EnvironmentId,
			&sfc.ScheduledAt,
			&sfc.Timezone,
			&mysql.JSONObject{Val: &sfc.Payload},
			&sfc.Comment,
			&status,
			&sfc.FailureReason,
			&sfc.FlagVersionAtCreation,
			&mysql.JSONObject{Val: &sfc.Conflicts},
			&sfc.CreatedBy,
			&sfc.CreatedAt,
			&sfc.UpdatedBy,
			&sfc.UpdatedAt,
			&sfc.ExecutedAt,
		)
		if err != nil {
			return nil, err
		}
		sfc.Status = proto.ScheduledFlagChangeStatus(status)
		scheduledFlagChanges = append(scheduledFlagChanges, &sfc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return scheduledFlagChanges, nil
}

func (s *scheduledFlagChangeStorage) TryLock(
	ctx context.Context,
	id, lockedBy string,
) (bool, error) {
	now := time.Now().Unix()
	lockExpiredAt := now - int64(LockExpiration.Seconds())

	result, err := s.qe.ExecContext(
		ctx,
		tryLockScheduledFlagChangeSQL,
		now,
		lockedBy,
		id,
		lockExpiredAt,
	)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	// If rowsAffected is 1, the lock was acquired successfully
	return rowsAffected == 1, nil
}

func (s *scheduledFlagChangeStorage) Unlock(
	ctx context.Context,
	id string,
) error {
	_, err := s.qe.ExecContext(
		ctx,
		unlockScheduledFlagChangeSQL,
		id,
	)
	return err
}
