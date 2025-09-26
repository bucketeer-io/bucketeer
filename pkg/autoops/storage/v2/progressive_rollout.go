// Copyright 2025 The Bucketeer Authors.
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

	err "github.com/bucketeer-io/bucketeer/v2/pkg/error"

	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
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
	ErrProgressiveRolloutAlreadyExists = err.NewErrorAlreadyExists(err.AutoopsPackageName, "already exists")
	ErrProgressiveRolloutNotFound      = err.NewErrorNotFound(
		err.AutoopsPackageName,
		"not found",
		"progressive_rollout",
	)
	ErrProgressiveRolloutUnexpectedAffectedRows = err.NewErrorUnexpectedAffectedRows(
		err.AutoopsPackageName,
		"unexpected affected rows",
	)
)

type progressiveRolloutStorage struct {
	qe mysql.QueryExecer
}

type ProgressiveRolloutStorage interface {
	CreateProgressiveRollout(
		ctx context.Context,
		progressiveRollout *domain.ProgressiveRollout,
		environmentId string,
	) error
	GetProgressiveRollout(ctx context.Context, id, environmentId string) (*domain.ProgressiveRollout, error)
	DeleteProgressiveRollout(ctx context.Context, id, environmentId string) error
	ListProgressiveRollouts(
		ctx context.Context,
		options *mysql.ListOptions,
	) ([]*autoopsproto.ProgressiveRollout, int64, int, error)
	UpdateProgressiveRollout(ctx context.Context,
		progressiveRollout *domain.ProgressiveRollout,
		environmentId string,
	) error
}

func NewProgressiveRolloutStorage(qe mysql.QueryExecer) ProgressiveRolloutStorage {
	return &progressiveRolloutStorage{qe: qe}
}

func (s *progressiveRolloutStorage) CreateProgressiveRollout(
	ctx context.Context,
	progressiveRollout *domain.ProgressiveRollout,
	environmentId string,
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
		environmentId,
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
	id, environmentId string,
) (*domain.ProgressiveRollout, error) {
	progressiveRollout := autoopsproto.ProgressiveRollout{}
	err := s.qe.QueryRowContext(
		ctx,
		selectOpsProgressiveRolloutSQL,
		id,
		environmentId,
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
	id, environmentId string,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		deleteOpsProgressiveRolloutSQL,
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
		return ErrProgressiveRolloutUnexpectedAffectedRows
	}
	return nil
}

func (s *progressiveRolloutStorage) ListProgressiveRollouts(
	ctx context.Context,
	options *mysql.ListOptions,
) ([]*autoopsproto.ProgressiveRollout, int64, int, error) {
	query, whereArgs := mysql.ConstructQueryAndWhereArgs(selectOpsProgressiveRolloutsSQL, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	progressiveRollouts := make([]*autoopsproto.ProgressiveRollout, 0)
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
	var offset int
	if options != nil {
		offset = options.Offset
	}
	nextOffset := offset + len(progressiveRollouts)
	var totalCount int64
	countQuery, countWhereArgs := mysql.ConstructCountQuery(countOpsProgressiveRolloutsSQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return progressiveRollouts, totalCount, nextOffset, nil
}

func (s *progressiveRolloutStorage) UpdateProgressiveRollout(
	ctx context.Context,
	progressiveRollout *domain.ProgressiveRollout,
	environmentId string,
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
		return ErrProgressiveRolloutUnexpectedAffectedRows
	}
	return nil
}
