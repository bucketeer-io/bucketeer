// Copyright 2024 The Bucketeer Authors.
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
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

var (
	//go:embed sql/ops_progressive_rollout/insert_ops_progressive_rollout.sql
	insertOpsProgressiveRolloutSQL string
	//go:embed sql/ops_progressive_rollout/update_ops_progressive_rollout.sql
	updateOpsProgressiveRolloutSQL string
	//go:embed sql/ops_progressive_rollout/select_ops_progressive_rollout.sql
	selectOpsProgressiveRolloutSQL string
	//go:embed sql/ops_progressive_rollout/select_ops_progressive_rollouts.sql
	selectOpsProgressiveRolloutsSQL string
	//go:embed sql/ops_progressive_rollout/count_ops_progressive_rollouts.sql
	countOpsProgressiveRolloutsSQL string
	//go:embed sql/ops_progressive_rollout/delete_ops_progressive_rollout.sql
	deleteOpsProgressiveRolloutSQL string
)

var (
	ErrProgressiveRolloutAlreadyExists          = errors.New("progressiveRollout: already exists")
	ErrProgressiveRolloutNotFound               = errors.New("progressiveRollout: not found")
	ErrProgressiveRolloutUnexpectedAffectedRows = errors.New("progressiveRollout: unexpected affected rows")
)

type progressiveRolloutStorage struct {
	qe mysql.QueryExecer
}

type ProgressiveRolloutStorage interface {
	CreateProgressiveRollout(
		ctx context.Context,
		progressiveRollout *domain.ProgressiveRollout,
		environmentNamespace string,
	) error
	GetProgressiveRollout(ctx context.Context, id, environmentNamespace string) (*domain.ProgressiveRollout, error)
	DeleteProgressiveRollout(ctx context.Context, id, environmentNamespace string) error
	ListProgressiveRollouts(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*autoopsproto.ProgressiveRollout, int64, int, error)
	UpdateProgressiveRollout(ctx context.Context,
		progressiveRollout *domain.ProgressiveRollout,
		environmentNamespace string,
	) error
}

func NewProgressiveRolloutStorage(qe mysql.QueryExecer) ProgressiveRolloutStorage {
	return &progressiveRolloutStorage{qe: qe}
}

func (s *progressiveRolloutStorage) CreateProgressiveRollout(
	ctx context.Context,
	progressiveRollout *domain.ProgressiveRollout,
	environmentNamespace string,
) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertOpsProgressiveRolloutSQL,
		progressiveRollout.Id,
		progressiveRollout.FeatureId,
		mysql.JSONObject{Val: progressiveRollout.Clause},
		int32(progressiveRollout.Status),
		int32(progressiveRollout.StoppedBy),
		int32(progressiveRollout.Type),
		progressiveRollout.StoppedAt,
		progressiveRollout.CreatedAt,
		progressiveRollout.UpdatedAt,
		environmentNamespace,
	)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrProgressiveRolloutAlreadyExists
		}
		return err
	}
	return nil
}

func (s *progressiveRolloutStorage) GetProgressiveRollout(
	ctx context.Context,
	id, environmentNamespace string,
) (*domain.ProgressiveRollout, error) {
	progressiveRollout := autoopsproto.ProgressiveRollout{}
	err := s.qe.QueryRowContext(
		ctx,
		selectOpsProgressiveRolloutSQL,
		id,
		environmentNamespace,
	).Scan(
		&progressiveRollout.Id,
		&progressiveRollout.FeatureId,
		&mysql.JSONObject{Val: &progressiveRollout.Clause},
		&progressiveRollout.Status,
		&progressiveRollout.StoppedBy,
		&progressiveRollout.Type,
		&progressiveRollout.StoppedAt,
		&progressiveRollout.CreatedAt,
		&progressiveRollout.UpdatedAt,
	)
	if err != nil {
		if err == mysql.ErrNoRows {
			return nil, ErrProgressiveRolloutNotFound
		}
		return nil, err
	}
	return &domain.ProgressiveRollout{ProgressiveRollout: &progressiveRollout}, nil
}

func (s *progressiveRolloutStorage) DeleteProgressiveRollout(
	ctx context.Context,
	id, environmentNamespace string,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		deleteOpsProgressiveRolloutSQL,
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
		return ErrProgressiveRolloutUnexpectedAffectedRows
	}
	return nil
}

func (s *progressiveRolloutStorage) ListProgressiveRollouts(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*autoopsproto.ProgressiveRollout, int64, int, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(selectOpsProgressiveRolloutsSQL, whereSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	progressiveRollouts := make([]*autoopsproto.ProgressiveRollout, 0, limit)
	for rows.Next() {
		progressiveRollout := autoopsproto.ProgressiveRollout{}
		err := rows.Scan(
			&progressiveRollout.Id,
			&progressiveRollout.FeatureId,
			&mysql.JSONObject{Val: &progressiveRollout.Clause},
			&progressiveRollout.Status,
			&progressiveRollout.StoppedBy,
			&progressiveRollout.Type,
			&progressiveRollout.StoppedAt,
			&progressiveRollout.CreatedAt,
			&progressiveRollout.UpdatedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		progressiveRollouts = append(progressiveRollouts, &progressiveRollout)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(progressiveRollouts)
	var totalCount int64
	countQuery := fmt.Sprintf(countOpsProgressiveRolloutsSQL, whereSQL, orderBySQL)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return progressiveRollouts, totalCount, nextOffset, nil
}

func (s *progressiveRolloutStorage) UpdateProgressiveRollout(
	ctx context.Context,
	progressiveRollout *domain.ProgressiveRollout,
	environmentNamespace string,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateOpsProgressiveRolloutSQL,
		&progressiveRollout.FeatureId,
		&mysql.JSONObject{Val: &progressiveRollout.Clause},
		&progressiveRollout.Status,
		&progressiveRollout.StoppedBy,
		&progressiveRollout.Type,
		&progressiveRollout.StoppedAt,
		&progressiveRollout.CreatedAt,
		&progressiveRollout.UpdatedAt,
		&progressiveRollout.Id,
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
		return ErrProgressiveRolloutUnexpectedAffectedRows
	}
	return nil
}
