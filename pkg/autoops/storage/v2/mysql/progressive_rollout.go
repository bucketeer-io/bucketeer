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
	"strconv"

	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/storage/v2"
	mysqlstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
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

type progressiveRolloutStorage struct {
	qe mysqlstorage.QueryExecer
}

func NewProgressiveRolloutStorage(qe mysqlstorage.QueryExecer) v2as.ProgressiveRolloutStorage {
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
		mysqlstorage.JSONObject{Val: progressiveRollout.Clause},
		int32(progressiveRollout.Status),
		int32(progressiveRollout.StoppedBy),
		int32(progressiveRollout.Type),
		progressiveRollout.StoppedAt,
		progressiveRollout.CreatedAt,
		progressiveRollout.UpdatedAt,
		environmentId,
	)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrDuplicateEntry) {
			return v2as.ErrProgressiveRolloutAlreadyExists
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
		&mysqlstorage.JSONObject{Val: &progressiveRollout.Clause},
		&progressiveRollout.Status,
		&progressiveRollout.StoppedBy,
		&progressiveRollout.Type,
		&progressiveRollout.StoppedAt,
		&progressiveRollout.CreatedAt,
		&progressiveRollout.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrNoRows) {
			return nil, v2as.ErrProgressiveRolloutNotFound
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
		return v2as.ErrProgressiveRolloutUnexpectedAffectedRows
	}
	return nil
}

func (s *progressiveRolloutStorage) ListProgressiveRollouts(
	ctx context.Context,
	params v2as.ListProgressiveRolloutsParams,
) ([]*autoopsproto.ProgressiveRollout, int64, int, error) {
	options, err := listProgressiveRolloutsOptionsFromParams(params)
	if err != nil {
		return nil, 0, 0, err
	}
	query, whereArgs := mysqlstorage.ConstructQueryAndWhereArgs(selectOpsProgressiveRolloutsSQL, options)
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
			&mysqlstorage.JSONObject{Val: &progressiveRollout.Clause},
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
		return nil, 0, 0, rows.Err()
	}
	nextOffset := options.Offset + len(progressiveRollouts)
	var totalCount int64
	countQuery, countWhereArgs := mysqlstorage.ConstructCountQuery(countOpsProgressiveRolloutsSQL, options)
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
		&mysqlstorage.JSONObject{Val: &progressiveRollout.Clause},
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
		return v2as.ErrProgressiveRolloutUnexpectedAffectedRows
	}
	return nil
}

func listProgressiveRolloutsOptionsFromParams(
	p v2as.ListProgressiveRolloutsParams,
) (*mysqlstorage.ListOptions, error) {
	filters := []*mysqlstorage.FilterV2{
		{
			Column:   "environment_id",
			Operator: mysqlstorage.OperatorEqual,
			Value:    p.EnvironmentID,
		},
	}
	if p.Type != nil {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "type",
			Operator: mysqlstorage.OperatorEqual,
			Value:    *p.Type,
		})
	}
	if p.Status != nil {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "status",
			Operator: mysqlstorage.OperatorEqual,
			Value:    *p.Status,
		})
	}
	var inFilters []*mysqlstorage.InFilter
	if len(p.FeatureIDs) > 0 {
		values := make([]interface{}, len(p.FeatureIDs))
		for i, id := range p.FeatureIDs {
			values[i] = id
		}
		inFilters = append(inFilters, &mysqlstorage.InFilter{
			Column: "feature_id",
			Values: values,
		})
	}

	var column string
	switch p.OrderBy {
	case autoopsproto.ListProgressiveRolloutsRequest_DEFAULT:
		column = "id"
	case autoopsproto.ListProgressiveRolloutsRequest_CREATED_AT:
		column = "created_at"
	case autoopsproto.ListProgressiveRolloutsRequest_UPDATED_AT:
		column = "updated_at"
	default:
		return nil, v2as.ErrInvalidOrderBy
	}
	direction := mysqlstorage.OrderDirectionAsc
	if p.OrderDirection == autoopsproto.ListProgressiveRolloutsRequest_DESC {
		direction = mysqlstorage.OrderDirectionDesc
	}
	orders := []*mysqlstorage.Order{mysqlstorage.NewOrder(column, direction)}

	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil || offset < 0 {
		return nil, v2as.ErrInvalidCursor
	}
	limit := p.PageSize
	if limit < 0 {
		limit = 0
	}
	return &mysqlstorage.ListOptions{
		Limit:     limit,
		Offset:    offset,
		Filters:   filters,
		InFilters: inFilters,
		Orders:    orders,
	}, nil
}
