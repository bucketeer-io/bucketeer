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
	"errors"
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

var (
	ErrGoalAlreadyExists          = errors.New("goal: already exists")
	ErrGoalNotFound               = errors.New("goal: not found")
	ErrGoalUnexpectedAffectedRows = errors.New("goal: unexpected affected rows")

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

type GoalStorage interface {
	CreateGoal(ctx context.Context, g *domain.Goal, environmentId string) error
	UpdateGoal(ctx context.Context, g *domain.Goal, environmentId string) error
	GetGoal(ctx context.Context, id, environmentId string) (*domain.Goal, error)
	ListGoals(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
		isInUseStatus *bool,
		environmentId string,
	) ([]*proto.Goal, int, int64, error)
	DeleteGoal(ctx context.Context, id, environmentId string) error
}

type goalStorage struct {
	client mysql.Client
}

func NewGoalStorage(client mysql.Client) GoalStorage {
	return &goalStorage{client: client}
}

func (s *goalStorage) CreateGoal(ctx context.Context, g *domain.Goal, environmentId string) error {
	_, err := s.client.Qe(ctx).ExecContext(
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
		if errors.Is(err, mysql.ErrDuplicateEntry) {
			return ErrGoalAlreadyExists
		}
		return err
	}
	return nil
}

func (s *goalStorage) UpdateGoal(ctx context.Context, g *domain.Goal, environmentId string) error {
	result, err := s.client.Qe(ctx).ExecContext(
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
		return ErrGoalUnexpectedAffectedRows
	}
	return nil
}

func (s *goalStorage) GetGoal(ctx context.Context, id, environmentId string) (*domain.Goal, error) {
	goal := proto.Goal{}
	var connectionType int32
	var experiments []experimentRef
	err := s.client.Qe(ctx).QueryRowContext(
		ctx,
		selectGoalSQL,
		environmentId, // Case query
		environmentId, // Subquery
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
		&mysql.JSONObject{Val: &goal.Experiments},
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrGoalNotFound
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
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
	isInUseStatus *bool,
	environmentId string,
) ([]*proto.Goal, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	prepareArgs := make([]interface{}, 0, len(whereArgs)+2)
	prepareArgs = append(prepareArgs, environmentId) // Case query
	prepareArgs = append(prepareArgs, environmentId) // Subquery
	prepareArgs = append(prepareArgs, whereArgs...)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	var isInUseStatusSQL string
	if isInUseStatus != nil {
		if *isInUseStatus {
			isInUseStatusSQL = "HAVING is_in_use_status = TRUE"
		} else {
			isInUseStatusSQL = "HAVING is_in_use_status = FALSE"
		}
	}
	query := fmt.Sprintf(selectGoalsSQL, whereSQL, isInUseStatusSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.client.Qe(ctx).QueryContext(ctx, query, prepareArgs...)
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
			&mysql.JSONObject{Val: &experiments},
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
		return nil, 0, 0, err
	}
	nextOffset := offset + len(goals)
	var totalCount int64
	countConditionSQL := "> 0 THEN 1 ELSE 1"
	if isInUseStatus != nil {
		if *isInUseStatus {
			countConditionSQL = "> 0 THEN 1 ELSE NULL"
		} else {
			countConditionSQL = "> 0 THEN NULL ELSE 1"
		}
	}
	prepareCountArgs := make([]interface{}, 0, len(whereArgs)+1)
	prepareCountArgs = append(prepareCountArgs, environmentId)
	prepareCountArgs = append(prepareCountArgs, whereArgs...)
	countQuery := fmt.Sprintf(countGoalSQL, countConditionSQL, whereSQL)
	err = s.client.Qe(ctx).QueryRowContext(ctx, countQuery, prepareCountArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return goals, nextOffset, totalCount, nil
}

func (s *goalStorage) DeleteGoal(ctx context.Context, id, environmentId string) error {
	result, err := s.client.Qe(ctx).ExecContext(
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
		return ErrGoalUnexpectedAffectedRows
	}
	return nil
}
