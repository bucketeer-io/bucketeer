// Copyright 2023 The Bucketeer Authors.
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
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	//go:embed sql/flag_trigger/insert_flag_trigger.sql
	insertFlagTriggerSQL string
	//go:embed sql/flag_trigger/update_flag_trigger.sql
	updateFlagTriggerSQL string
	//go:embed sql/flag_trigger/update_flag_trigger_usage.sql
	updateFlagTriggerUsageSQL string
	//go:embed sql/flag_trigger/reset_flag_trigger.sql
	resetFlagTriggerSQL string
	//go:embed sql/flag_trigger/delete_flag_trigger.sql
	deleteFlagTriggerSQL string
	//go:embed sql/flag_trigger/get_flag_trigger.sql
	getFlagTriggerSQL string
	//go:embed sql/flag_trigger/list_flag_trigger.sql
	listFlagTriggersSQL string
	//go:embed sql/flag_trigger/count_flag_trigger.sql
	countFlagTriggersSQL string
	//go:embed sql/flag_trigger/enable_flag_trigger.sql
	enableFlagTriggerSQL string
	//go:embed sql/flag_trigger/disable_flag_trigger.sql
	disableFlagTriggerSQL string
)

var (
	ErrFlagTriggerAlreadyExists          = errors.New("flag trigger: already exists")
	ErrFlagTriggerNotFound               = errors.New("flag trigger: not found")
	ErrFlagTriggerUnexpectedAffectedRows = errors.New("flag trigger: unexpected affected rows")
)

type FlagTriggerStorage interface {
	CreateFlagTrigger(ctx context.Context, flagTrigger *domain.FlagTrigger) error
	UpdateFlagTrigger(
		ctx context.Context,
		id, environmentNamespace, description string,
	) error
	UpdateFlagTriggerUsage(ctx context.Context,
		id, environmentNamespace string,
		triggerTimes int64,
	) error
	EnableFlagTrigger(ctx context.Context, id, environmentNamespace string) error
	DisableFlagTrigger(ctx context.Context, id, environmentNamespace string) error
	ResetFlagTrigger(
		ctx context.Context,
		id, environmentNamespace, uuid string,
	) error
	DeleteFlagTrigger(ctx context.Context, id, environmentNamespace string) error
	GetFlagTrigger(ctx context.Context, id, environmentNamespace string) (*domain.FlagTrigger, error)
	ListFlagTriggers(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.FlagTrigger, int, int64, error)
}

type flagTriggerStorage struct {
	qe mysql.QueryExecer
}

func NewFlagTriggerStorage(
	qe mysql.QueryExecer,
) FlagTriggerStorage {
	return &flagTriggerStorage{qe: qe}
}

func (f flagTriggerStorage) CreateFlagTrigger(
	ctx context.Context,
	flagTrigger *domain.FlagTrigger,
) error {
	_, err := f.qe.ExecContext(ctx, insertFlagTriggerSQL,
		flagTrigger.Id,
		flagTrigger.FeatureId,
		flagTrigger.EnvironmentNamespace,
		flagTrigger.Type,
		flagTrigger.Action,
		flagTrigger.Description,
		flagTrigger.TriggerTimes,
		flagTrigger.LastTriggeredAt,
		flagTrigger.Uuid,
		flagTrigger.Disabled,
		flagTrigger.Deleted,
		flagTrigger.CreatedAt,
		flagTrigger.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrDuplicateEntry) {
			return ErrFlagTriggerAlreadyExists
		}
		return err
	}
	return nil
}

func (f flagTriggerStorage) UpdateFlagTrigger(
	ctx context.Context,
	id, environmentNamespace, description string,
) error {
	result, err := f.qe.ExecContext(ctx, updateFlagTriggerSQL,
		description,
		time.Now().Unix(),
		id,
		environmentNamespace,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrFlagTriggerUnexpectedAffectedRows
	}
	return nil
}

func (f flagTriggerStorage) UpdateFlagTriggerUsage(
	ctx context.Context,
	id, environmentNamespace string,
	triggerTimes int64,
) error {
	now := time.Now().Unix()
	result, err := f.qe.ExecContext(ctx, updateFlagTriggerUsageSQL,
		triggerTimes,
		now,
		now,
		id,
		environmentNamespace,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrFlagTriggerUnexpectedAffectedRows
	}
	return nil
}

func (f flagTriggerStorage) EnableFlagTrigger(
	ctx context.Context,
	id, environmentNamespace string,
) error {
	result, err := f.qe.ExecContext(ctx, enableFlagTriggerSQL,
		time.Now().Unix(),
		id,
		environmentNamespace,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrFlagTriggerUnexpectedAffectedRows
	}
	return nil
}

func (f flagTriggerStorage) DisableFlagTrigger(
	ctx context.Context,
	id, environmentNamespace string,
) error {
	result, err := f.qe.ExecContext(ctx, disableFlagTriggerSQL,
		time.Now().Unix(),
		id,
		environmentNamespace,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrFlagTriggerUnexpectedAffectedRows
	}
	return nil
}

func (f flagTriggerStorage) ResetFlagTrigger(
	ctx context.Context,
	id, environmentNamespace, uuid string,
) error {
	result, err := f.qe.ExecContext(ctx, resetFlagTriggerSQL,
		uuid,
		time.Now().Unix(),
		id,
		environmentNamespace,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrFlagTriggerUnexpectedAffectedRows
	}
	return nil
}

func (f flagTriggerStorage) DeleteFlagTrigger(
	ctx context.Context,
	id, environmentNamespace string,
) error {
	result, err := f.qe.ExecContext(ctx, deleteFlagTriggerSQL,
		time.Now().Unix(),
		id,
		environmentNamespace,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrFlagTriggerUnexpectedAffectedRows
	}
	return nil
}

func (f flagTriggerStorage) GetFlagTrigger(
	ctx context.Context,
	id, environmentNamespace string,
) (*domain.FlagTrigger, error) {
	trigger := proto.FlagTrigger{}
	err := f.qe.QueryRowContext(
		ctx,
		getFlagTriggerSQL,
		id,
		environmentNamespace,
	).Scan(
		&trigger.Id,
		&trigger.FeatureId,
		&trigger.EnvironmentNamespace,
		&trigger.Type,
		&trigger.Action,
		&trigger.Description,
		&trigger.TriggerTimes,
		&trigger.LastTriggeredAt,
		&trigger.Uuid,
		&trigger.Disabled,
		&trigger.Deleted,
		&trigger.CreatedAt,
		&trigger.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrFlagTriggerNotFound
		}
		return nil, err
	}
	return &domain.FlagTrigger{FlagTrigger: &trigger}, nil
}

func (f flagTriggerStorage) ListFlagTriggers(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.FlagTrigger, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(listFlagTriggersSQL, whereSQL, orderBySQL, limitOffsetSQL)
	rows, err := f.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	flagTriggers := make([]*proto.FlagTrigger, 0, limit)
	for rows.Next() {
		trigger := proto.FlagTrigger{}
		err := rows.Scan(
			&trigger.Id,
			&trigger.FeatureId,
			&trigger.EnvironmentNamespace,
			&trigger.Type,
			&trigger.Action,
			&trigger.Description,
			&trigger.TriggerTimes,
			&trigger.LastTriggeredAt,
			&trigger.Uuid,
			&trigger.Disabled,
			&trigger.Deleted,
			&trigger.CreatedAt,
			&trigger.UpdatedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		flagTriggers = append(flagTriggers, &trigger)
	}
	nextOffset := offset + len(flagTriggers)
	var totalCount int64
	countQuery := fmt.Sprintf(countFlagTriggersSQL, whereSQL, orderBySQL)
	if err := f.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount); err != nil {
		return nil, 0, 0, err
	}
	return flagTriggers, nextOffset, totalCount, nil
}
