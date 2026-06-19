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

package mysql

import (
	"context"
	_ "embed"
	"errors"
	"time"

	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	mysqlstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
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
	qe mysqlstorage.QueryExecer
}

// NewScheduledFlagChangeStorage creates a new ScheduledFlagChangeStorage
func NewScheduledFlagChangeStorage(qe mysqlstorage.QueryExecer) v2fs.ScheduledFlagChangeStorage {
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
		mysqlstorage.JSONObject{Val: sfc.Payload},
		sfc.Comment,
		int32(sfc.Status),
		sfc.FailureReason,
		sfc.FlagVersionAtCreation,
		mysqlstorage.JSONObject{Val: sfc.Conflicts},
		sfc.CreatedBy,
		sfc.CreatedAt,
		sfc.UpdatedBy,
		sfc.UpdatedAt,
		sfc.ExecutedAt,
	)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrDuplicateEntry) {
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
		mysqlstorage.JSONObject{Val: sfc.Payload},
		sfc.Comment,
		int32(sfc.Status),
		sfc.FailureReason,
		mysqlstorage.JSONObject{Val: sfc.Conflicts},
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
		&mysqlstorage.JSONObject{Val: &sfc.Payload},
		&sfc.Comment,
		&status,
		&sfc.FailureReason,
		&sfc.FlagVersionAtCreation,
		&mysqlstorage.JSONObject{Val: &sfc.Conflicts},
		&sfc.CreatedBy,
		&sfc.CreatedAt,
		&sfc.UpdatedBy,
		&sfc.UpdatedAt,
		&sfc.ExecutedAt,
	)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrNoRows) {
			return nil, v2fs.ErrScheduledFlagChangeNotFound
		}
		return nil, err
	}
	sfc.Status = proto.ScheduledFlagChangeStatus(status)
	return &domain.ScheduledFlagChange{ScheduledFlagChange: &sfc}, nil
}

func scheduledFlagChangesListOptions(p v2fs.ListScheduledFlagChangesParams) *mysqlstorage.ListOptions {
	var filters []*mysqlstorage.FilterV2
	if p.EnvironmentID != "" {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "environment_id",
			Operator: mysqlstorage.OperatorEqual,
			Value:    p.EnvironmentID,
		})
	}
	if p.FeatureID != "" {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "feature_id",
			Operator: mysqlstorage.OperatorEqual,
			Value:    p.FeatureID,
		})
	}
	if p.ExcludeFeatureID != "" {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "feature_id",
			Operator: mysqlstorage.OperatorNotEqual,
			Value:    p.ExcludeFeatureID,
		})
	}
	if p.FromScheduledAt > 0 {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "scheduled_at",
			Operator: mysqlstorage.OperatorGreaterThanOrEqual,
			Value:    p.FromScheduledAt,
		})
	}
	if p.ToScheduledAt > 0 {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "scheduled_at",
			Operator: mysqlstorage.OperatorLessThanOrEqual,
			Value:    p.ToScheduledAt,
		})
	}
	var inFilters []*mysqlstorage.InFilter
	if len(p.Statuses) > 0 {
		statusValues := make([]interface{}, 0, len(p.Statuses))
		for _, status := range p.Statuses {
			statusValues = append(statusValues, int32(status))
		}
		inFilters = append(inFilters, &mysqlstorage.InFilter{
			Column: "status",
			Values: statusValues,
		})
	}
	return &mysqlstorage.ListOptions{
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
) []*mysqlstorage.Order {
	direction := mysqlstorage.OrderDirectionAsc
	if orderDirection == proto.ListScheduledFlagChangesRequest_DESC {
		direction = mysqlstorage.OrderDirectionDesc
	}
	switch orderBy {
	case proto.ListScheduledFlagChangesRequest_CREATED_AT:
		return []*mysqlstorage.Order{mysqlstorage.NewOrder("created_at", direction)}
	case proto.ListScheduledFlagChangesRequest_SCHEDULED_AT:
		return []*mysqlstorage.Order{mysqlstorage.NewOrder("scheduled_at", direction)}
	default:
		return []*mysqlstorage.Order{mysqlstorage.NewOrder("scheduled_at", mysqlstorage.OrderDirectionAsc)}
	}
}

func (s *scheduledFlagChangeStorage) ListScheduledFlagChanges(
	ctx context.Context,
	params v2fs.ListScheduledFlagChangesParams,
) ([]*proto.ScheduledFlagChange, int, int64, error) {
	options := scheduledFlagChangesListOptions(params)
	query, whereArgs := mysqlstorage.ConstructQueryAndWhereArgs(listScheduledFlagChangesSQL, options)
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
			&mysqlstorage.JSONObject{Val: &sfc.Payload},
			&sfc.Comment,
			&status,
			&sfc.FailureReason,
			&sfc.FlagVersionAtCreation,
			&mysqlstorage.JSONObject{Val: &sfc.Conflicts},
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
	countQuery, countWhereArgs := mysqlstorage.ConstructCountQuery(countScheduledFlagChangesSQL, options)
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
			&mysqlstorage.JSONObject{Val: &sfc.Payload},
			&sfc.Comment,
			&status,
			&sfc.FailureReason,
			&sfc.FlagVersionAtCreation,
			&mysqlstorage.JSONObject{Val: &sfc.Conflicts},
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
