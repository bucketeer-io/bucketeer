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
	_ "embed"
	"errors"

	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"

	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var (
	//go:embed sql/flag_trigger/insert_flag_trigger.sql
	insertFlagTriggerSQL string
	//go:embed sql/flag_trigger/update_flag_trigger.sql
	updateFlagTriggerSQL string
	//go:embed sql/flag_trigger/delete_flag_trigger.sql
	deleteFlagTriggerSQL string
	//go:embed sql/flag_trigger/get_flag_trigger.sql
	getFlagTriggerSQL string
	//go:embed sql/flag_trigger/get_flag_trigger_by_token.sql
	getFlagTriggerByTokenSQL string
	//go:embed sql/flag_trigger/list_flag_trigger.sql
	listFlagTriggersSQL string
	//go:embed sql/flag_trigger/count_flag_trigger.sql
	countFlagTriggersSQL string
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
)

type FlagTriggerStorage interface {
	CreateFlagTrigger(ctx context.Context, flagTrigger *domain.FlagTrigger) error
	UpdateFlagTrigger(ctx context.Context, flagTrigger *domain.FlagTrigger) error
	DeleteFlagTrigger(ctx context.Context, id, environmentId string) error
	GetFlagTrigger(ctx context.Context, id, environmentId string) (*domain.FlagTrigger, error)
	GetFlagTriggerByToken(ctx context.Context, token string) (*domain.FlagTrigger, error)
	ListFlagTriggers(
		ctx context.Context,
		options *mysql.ListOptions,
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
		flagTrigger.EnvironmentId,
		flagTrigger.Type,
		flagTrigger.Action,
		flagTrigger.Description,
		flagTrigger.TriggerCount,
		flagTrigger.LastTriggeredAt,
		flagTrigger.Token,
		flagTrigger.Disabled,
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

func (f flagTriggerStorage) UpdateFlagTrigger(ctx context.Context, flagTrigger *domain.FlagTrigger) error {
	result, err := f.qe.ExecContext(ctx, updateFlagTriggerSQL,
		flagTrigger.FeatureId,
		flagTrigger.Type,
		flagTrigger.Action,
		flagTrigger.Description,
		flagTrigger.TriggerCount,
		flagTrigger.LastTriggeredAt,
		flagTrigger.Token,
		flagTrigger.Disabled,
		flagTrigger.CreatedAt,
		flagTrigger.UpdatedAt,
		flagTrigger.Id,
		flagTrigger.EnvironmentId,
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
	id, environmentId string,
) error {
	result, err := f.qe.ExecContext(ctx, deleteFlagTriggerSQL,
		id,
		environmentId,
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
	id, environmentId string,
) (*domain.FlagTrigger, error) {
	trigger := proto.FlagTrigger{}
	err := f.qe.QueryRowContext(
		ctx,
		getFlagTriggerSQL,
		id,
		environmentId,
	).Scan(
		&trigger.Id,
		&trigger.FeatureId,
		&trigger.EnvironmentId,
		&trigger.Type,
		&trigger.Action,
		&trigger.Description,
		&trigger.TriggerCount,
		&trigger.LastTriggeredAt,
		&trigger.Token,
		&trigger.Disabled,
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

func (f flagTriggerStorage) GetFlagTriggerByToken(
	ctx context.Context,
	token string,
) (*domain.FlagTrigger, error) {
	trigger := proto.FlagTrigger{}
	err := f.qe.QueryRowContext(
		ctx,
		getFlagTriggerByTokenSQL,
		token,
	).Scan(
		&trigger.Id,
		&trigger.FeatureId,
		&trigger.EnvironmentId,
		&trigger.Type,
		&trigger.Action,
		&trigger.Description,
		&trigger.TriggerCount,
		&trigger.LastTriggeredAt,
		&trigger.Token,
		&trigger.Disabled,
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
	options *mysql.ListOptions,
) ([]*proto.FlagTrigger, int, int64, error) {
	query, whereArgs := mysql.ConstructQueryAndWhereArgs(listFlagTriggersSQL, options)
	rows, err := f.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	var limit, offset int
	if options != nil {
		limit = options.Limit
		offset = options.Offset
	}
	flagTriggers := make([]*proto.FlagTrigger, 0, limit)
	for rows.Next() {
		trigger := proto.FlagTrigger{}
		err := rows.Scan(
			&trigger.Id,
			&trigger.FeatureId,
			&trigger.Type,
			&trigger.Action,
			&trigger.Description,
			&trigger.TriggerCount,
			&trigger.LastTriggeredAt,
			&trigger.Token,
			&trigger.Disabled,
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
	countQuery, countWhereArgs := mysql.ConstructCountQuery(countFlagTriggersSQL, options)
	if err := f.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount); err != nil {
		return nil, 0, 0, err
	}
	return flagTriggers, nextOffset, totalCount, nil
}
