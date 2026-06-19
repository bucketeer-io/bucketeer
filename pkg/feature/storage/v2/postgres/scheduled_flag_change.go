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

package postgres

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
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

type scheduledFlagChangeStorage struct {
	qe pgstorage.QueryExecer
}

// NewScheduledFlagChangeStorage creates a new ScheduledFlagChangeStorage
func NewScheduledFlagChangeStorage(qe pgstorage.QueryExecer) v2fs.ScheduledFlagChangeStorage {
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
		pgstorage.JSONObject{Val: sfc.Payload},
		sfc.Comment,
		int32(sfc.Status),
		sfc.FailureReason,
		sfc.FlagVersionAtCreation,
		pgstorage.JSONObject{Val: sfc.Conflicts},
		sfc.CreatedBy,
		sfc.CreatedAt,
		sfc.UpdatedBy,
		sfc.UpdatedAt,
		sfc.ExecutedAt,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrDuplicateEntry) {
			return v2fs.ErrScheduledFlagChangeAlreadyExists
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
		pgstorage.JSONObject{Val: sfc.Payload},
		sfc.Comment,
		int32(sfc.Status),
		sfc.FailureReason,
		pgstorage.JSONObject{Val: sfc.Conflicts},
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
		return v2fs.ErrScheduledFlagChangeUnexpectedAffectedRows
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
		return v2fs.ErrScheduledFlagChangeUnexpectedAffectedRows
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
		&pgstorage.JSONObject{Val: &sfc.Payload},
		&sfc.Comment,
		&status,
		&sfc.FailureReason,
		&sfc.FlagVersionAtCreation,
		&pgstorage.JSONObject{Val: &sfc.Conflicts},
		&sfc.CreatedBy,
		&sfc.CreatedAt,
		&sfc.UpdatedBy,
		&sfc.UpdatedAt,
		&sfc.ExecutedAt,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, v2fs.ErrScheduledFlagChangeNotFound
		}
		return nil, err
	}
	sfc.Status = proto.ScheduledFlagChangeStatus(status)
	return &domain.ScheduledFlagChange{ScheduledFlagChange: &sfc}, nil
}

func scheduledFlagChangesListOptions(p v2fs.ListScheduledFlagChangesParams) *pgstorage.ListOptions {
	var filters []*pgstorage.Filter
	if p.EnvironmentID != "" {
		filters = append(filters, &pgstorage.Filter{
			Column:   "environment_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.EnvironmentID,
		})
	}
	if p.FeatureID != "" {
		filters = append(filters, &pgstorage.Filter{
			Column:   "feature_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.FeatureID,
		})
	}
	if p.ExcludeFeatureID != "" {
		filters = append(filters, &pgstorage.Filter{
			Column:   "feature_id",
			Operator: pgstorage.OperatorNotEqual,
			Value:    p.ExcludeFeatureID,
		})
	}
	if p.FromScheduledAt > 0 {
		filters = append(filters, &pgstorage.Filter{
			Column:   "scheduled_at",
			Operator: pgstorage.OperatorGreaterThanOrEqual,
			Value:    p.FromScheduledAt,
		})
	}
	if p.ToScheduledAt > 0 {
		filters = append(filters, &pgstorage.Filter{
			Column:   "scheduled_at",
			Operator: pgstorage.OperatorLessThanOrEqual,
			Value:    p.ToScheduledAt,
		})
	}
	var inFilters []*pgstorage.InFilter
	if len(p.Statuses) > 0 {
		statusValues := make([]interface{}, 0, len(p.Statuses))
		for _, status := range p.Statuses {
			statusValues = append(statusValues, int32(status))
		}
		inFilters = append(inFilters, &pgstorage.InFilter{
			Column: "status",
			Values: statusValues,
		})
	}
	return &pgstorage.ListOptions{
		Filters:   filters,
		InFilters: inFilters,
		Orders:    scheduledFlagChangesOrders(p.OrderBy, p.OrderDirection),
		Limit:     p.PageSize,
		Offset:    p.Offset,
	}
}

func scheduledFlagChangesOrders(
	orderBy proto.ListScheduledFlagChangesRequest_OrderBy,
	orderDirection proto.ListScheduledFlagChangesRequest_OrderDirection,
) []*pgstorage.Order {
	direction := pgstorage.OrderDirectionAsc
	if orderDirection == proto.ListScheduledFlagChangesRequest_DESC {
		direction = pgstorage.OrderDirectionDesc
	}
	switch orderBy {
	case proto.ListScheduledFlagChangesRequest_CREATED_AT:
		return []*pgstorage.Order{pgstorage.NewOrder("created_at", direction)}
	case proto.ListScheduledFlagChangesRequest_SCHEDULED_AT:
		return []*pgstorage.Order{pgstorage.NewOrder("scheduled_at", direction)}
	default:
		return []*pgstorage.Order{pgstorage.NewOrder("scheduled_at", pgstorage.OrderDirectionAsc)}
	}
}

func (s *scheduledFlagChangeStorage) ListScheduledFlagChanges(
	ctx context.Context,
	params v2fs.ListScheduledFlagChangesParams,
) ([]*proto.ScheduledFlagChange, int, int64, error) {
	options := scheduledFlagChangesListOptions(params)
	whereParts := options.CreateWhereParts()
	whereSQL, whereArgs := pgstorage.ConstructWhereSQLString(whereParts)
	orderBySQL := pgstorage.ConstructOrderBySQLString(options.Orders)
	limitOffsetSQL := pgstorage.ConstructLimitOffsetSQLString(options.Limit, options.Offset)
	query := fmt.Sprintf(listScheduledFlagChangesSQL, whereSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	scheduledFlagChanges := make([]*proto.ScheduledFlagChange, 0, options.Limit)
	for rows.Next() {
		sfc := proto.ScheduledFlagChange{}
		var status int32
		err := rows.Scan(
			&sfc.Id,
			&sfc.FeatureId,
			&sfc.EnvironmentId,
			&sfc.ScheduledAt,
			&sfc.Timezone,
			&pgstorage.JSONObject{Val: &sfc.Payload},
			&sfc.Comment,
			&status,
			&sfc.FailureReason,
			&sfc.FlagVersionAtCreation,
			&pgstorage.JSONObject{Val: &sfc.Conflicts},
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

	nextOffset := options.Offset + len(scheduledFlagChanges)
	var totalCount int64
	countQuery := fmt.Sprintf(countScheduledFlagChangesSQL, whereSQL)
	if err := s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount); err != nil {
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
	lockExpiredAt := now - int64(v2fs.LockExpiration.Seconds())

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
			&pgstorage.JSONObject{Val: &sfc.Payload},
			&sfc.Comment,
			&status,
			&sfc.FailureReason,
			&sfc.FlagVersionAtCreation,
			&pgstorage.JSONObject{Val: &sfc.Conflicts},
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
	lockExpiredAt := now - int64(v2fs.LockExpiration.Seconds())

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
	id, lockedBy string,
) error {
	_, err := s.qe.ExecContext(
		ctx,
		unlockScheduledFlagChangeSQL,
		id,
		lockedBy,
	)
	return err
}
