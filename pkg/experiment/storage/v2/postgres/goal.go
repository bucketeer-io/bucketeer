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
	"strconv"

	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/storage/v2"
	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

var (
	//go:embed sql/goal/select_goals.sql
	selectGoalsSQL string
	//go:embed sql/goal/select_goal.sql
	selectGoalSQL string
	//go:embed sql/goal/count_goals.sql
	countGoalSQL string
	//go:embed sql/goal/insert_goal.sql
	insertGoalSQL string
	//go:embed sql/goal/update_goal.sql
	updateGoalSQL string
	//go:embed sql/goal/delete_goal.sql
	deleteGoalSQL string
)

type experimentRef struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	FeatureId   string `json:"feature_id"`
	FeatureName string `json:"feature_name"`
	Status      int32  `json:"status"`
}

type goalStorage struct {
	qe pgstorage.QueryExecer
}

func NewGoalStorage(qe pgstorage.QueryExecer) v2es.GoalStorage {
	return &goalStorage{qe: qe}
}

func (s *goalStorage) CreateGoal(ctx context.Context, g *domain.Goal, environmentId string) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertGoalSQL,
		g.Id,
		g.Name,
		g.Description,
		g.ConnectionType,
		g.Archived,
		g.Deleted,
		g.CreatedAt,
		g.UpdatedAt,
		environmentId,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrDuplicateEntry) {
			return v2es.ErrGoalAlreadyExists
		}
		return err
	}
	return nil
}

func (s *goalStorage) UpdateGoal(ctx context.Context, g *domain.Goal, environmentId string) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateGoalSQL,
		g.Name,
		g.Description,
		g.Archived,
		g.Deleted,
		g.CreatedAt,
		g.UpdatedAt,
		g.Id,
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
		return v2es.ErrGoalUnexpectedAffectedRows
	}
	return nil
}

func (s *goalStorage) GetGoal(ctx context.Context, id, environmentId string) (*domain.Goal, error) {
	goal := proto.Goal{}
	var connectionType int32
	var experiments []experimentRef
	err := s.qe.QueryRowContext(
		ctx,
		selectGoalSQL,
		id,
		environmentId,
	).Scan(
		&goal.Id,
		&goal.Name,
		&goal.Description,
		&connectionType,
		&goal.Archived,
		&goal.Deleted,
		&goal.CreatedAt,
		&goal.UpdatedAt,
		&goal.IsInUseStatus,
		&pgstorage.JSONObject{Val: &goal.Experiments},
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, v2es.ErrGoalNotFound
		}
		return nil, err
	}
	goal.ConnectionType = proto.Goal_ConnectionType(connectionType)
	for i := range experiments {
		goal.Experiments = append(goal.Experiments, &proto.Goal_ExperimentReference{
			Id:          experiments[i].Id,
			Name:        experiments[i].Name,
			FeatureId:   experiments[i].FeatureId,
			FeatureName: experiments[i].FeatureName,
			Status:      proto.Experiment_Status(experiments[i].Status),
		})
	}
	return &domain.Goal{Goal: &goal}, nil
}

func (s *goalStorage) ListGoals(
	ctx context.Context,
	params v2es.ListGoalsParams,
) ([]*proto.Goal, int, int64, error) {
	whereParts := []pgstorage.WherePart{
		&pgstorage.Filter{Column: "deleted", Operator: pgstorage.OperatorEqual, Value: false},
		&pgstorage.Filter{Column: "environment_id", Operator: pgstorage.OperatorEqual, Value: params.EnvironmentID},
	}
	if params.Archived != nil {
		whereParts = append(whereParts, &pgstorage.Filter{
			Column:   "archived",
			Operator: pgstorage.OperatorEqual,
			Value:    *params.Archived,
		})
	}
	if params.SearchKeyword != "" {
		whereParts = append(whereParts, &pgstorage.SearchQuery{
			Columns: []string{"id", "name", "description"},
			Keyword: params.SearchKeyword,
		})
	}
	if params.ConnectionType != proto.Goal_UNKNOWN {
		whereParts = append(whereParts, &pgstorage.Filter{
			Column:   "connection_type",
			Operator: pgstorage.OperatorEqual,
			Value:    params.ConnectionType,
		})
	}
	orders, err := goalListOrders(params.OrderBy, params.OrderDirection)
	if err != nil {
		return nil, 0, 0, err
	}
	cursor := params.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil || offset < 0 {
		return nil, 0, 0, v2es.ErrInvalidCursor
	}
	limit := params.PageSize
	if limit < 0 {
		limit = 0
	}
	whereSQL, whereArgs := pgstorage.ConstructWhereSQLString(whereParts)
	orderBySQL := pgstorage.ConstructOrderBySQLString(orders)
	limitOffsetSQL := pgstorage.ConstructLimitOffsetSQLString(limit, offset)
	var isInUseStatusSQL string
	if params.IsInUseStatus != nil {
		if *params.IsInUseStatus {
			isInUseStatusSQL = "WHERE is_in_use_status = TRUE"
		} else {
			isInUseStatusSQL = "WHERE is_in_use_status = FALSE"
		}
	}
	query := fmt.Sprintf(selectGoalsSQL, whereSQL, isInUseStatusSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	goals := make([]*proto.Goal, 0, limit)

	for rows.Next() {
		goal := proto.Goal{}
		var connectionType int32
		var experiments []experimentRef
		err := rows.Scan(
			&goal.Id,
			&goal.Name,
			&goal.Description,
			&connectionType,
			&goal.Archived,
			&goal.Deleted,
			&goal.CreatedAt,
			&goal.UpdatedAt,
			&goal.IsInUseStatus,
			&pgstorage.JSONObject{Val: &experiments},
		)
		if err != nil {
			return nil, 0, 0, err
		}
		goal.ConnectionType = proto.Goal_ConnectionType(connectionType)
		for i := range experiments {
			goal.Experiments = append(goal.Experiments, &proto.Goal_ExperimentReference{
				Id:          experiments[i].Id,
				Name:        experiments[i].Name,
				FeatureId:   experiments[i].FeatureId,
				FeatureName: experiments[i].FeatureName,
				Status:      proto.Experiment_Status(experiments[i].Status),
			})
		}
		goals = append(goals, &goal)
	}
	if rows.Err() != nil {
		return nil, 0, 0, rows.Err()
	}
	nextOffset := offset + len(goals)
	var totalCount int64
	countConditionSQL := "> 0 THEN 1 ELSE 1"
	if params.IsInUseStatus != nil {
		if *params.IsInUseStatus {
			countConditionSQL = "> 0 THEN 1 ELSE NULL"
		} else {
			countConditionSQL = "> 0 THEN NULL ELSE 1"
		}
	}
	countQuery := fmt.Sprintf(countGoalSQL, countConditionSQL, whereSQL)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return goals, nextOffset, totalCount, nil
}

func (s *goalStorage) DeleteGoal(ctx context.Context, id, environmentId string) error {
	result, err := s.qe.ExecContext(
		ctx,
		deleteGoalSQL,
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
		return v2es.ErrGoalUnexpectedAffectedRows
	}
	return nil
}

func goalListOrders(
	orderBy proto.ListGoalsRequest_OrderBy,
	orderDirection proto.ListGoalsRequest_OrderDirection,
) ([]*pgstorage.Order, error) {
	var column string
	switch orderBy {
	case proto.ListGoalsRequest_DEFAULT,
		proto.ListGoalsRequest_NAME:
		column = "name"
	case proto.ListGoalsRequest_CREATED_AT:
		column = "created_at"
	case proto.ListGoalsRequest_UPDATED_AT:
		column = "updated_at"
	case proto.ListGoalsRequest_CONNECTION_TYPE:
		column = "connection_type"
	default:
		return nil, v2es.ErrInvalidOrderBy
	}
	direction := pgstorage.OrderDirectionAsc
	if orderDirection == proto.ListGoalsRequest_DESC {
		direction = pgstorage.OrderDirectionDesc
	}
	return []*pgstorage.Order{pgstorage.NewOrder(column, direction)}, nil
}
