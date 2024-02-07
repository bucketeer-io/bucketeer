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
)

type GoalStorage interface {
	CreateGoal(ctx context.Context, g *domain.Goal, environmentNamespace string) error
	UpdateGoal(ctx context.Context, g *domain.Goal, environmentNamespace string) error
	GetGoal(ctx context.Context, id, environmentNamespace string) (*domain.Goal, error)
	ListGoals(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
		isInUseStatus *bool,
		environmentNamespace string,
	) ([]*proto.Goal, int, int64, error)
}

type goalStorage struct {
	qe mysql.QueryExecer
}

func NewGoalStorage(qe mysql.QueryExecer) GoalStorage {
	return &goalStorage{qe: qe}
}

func (s *goalStorage) CreateGoal(ctx context.Context, g *domain.Goal, environmentNamespace string) error {
	query := `
		INSERT INTO goal (
			id,
			name,
			description,
			archived,
			deleted,
			created_at,
			updated_at,
			environment_namespace
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?
		)
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
		g.Id,
		g.Name,
		g.Description,
		g.Archived,
		g.Deleted,
		g.CreatedAt,
		g.UpdatedAt,
		environmentNamespace,
	)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrGoalAlreadyExists
		}
		return err
	}
	return nil
}

func (s *goalStorage) UpdateGoal(ctx context.Context, g *domain.Goal, environmentNamespace string) error {
	query := `
		UPDATE 
			goal
		SET
			name = ?,
			description = ?,
			archived = ?,
			deleted = ?,
			created_at = ?,
			updated_at = ?
		WHERE
			id = ? AND
			environment_namespace = ?
	`
	result, err := s.qe.ExecContext(
		ctx,
		query,
		g.Name,
		g.Description,
		g.Archived,
		g.Deleted,
		g.CreatedAt,
		g.UpdatedAt,
		g.Id,
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
		return ErrGoalUnexpectedAffectedRows
	}
	return nil
}

func (s *goalStorage) GetGoal(ctx context.Context, id, environmentNamespace string) (*domain.Goal, error) {
	goal := proto.Goal{}
	query := `
		SELECT
			id,
			name,
			description,
			archived,
			deleted,
			created_at,
			updated_at,
			CASE 
				WHEN (
					SELECT 
						COUNT(1)
					FROM 
						experiment
					WHERE
						environment_namespace = ? AND
						goal_ids LIKE concat("%", goal.id, "%")
				) > 0 THEN TRUE 
				ELSE FALSE 
			END AS is_in_use_status
		FROM
			goal
		WHERE
			id = ? AND
			environment_namespace = ?
	`
	err := s.qe.QueryRowContext(
		ctx,
		query,
		environmentNamespace,
		id,
		environmentNamespace,
	).Scan(
		&goal.Id,
		&goal.Name,
		&goal.Description,
		&goal.Archived,
		&goal.Deleted,
		&goal.CreatedAt,
		&goal.UpdatedAt,
		&goal.IsInUseStatus,
	)
	if err != nil {
		if err == mysql.ErrNoRows {
			return nil, ErrGoalNotFound
		}
		return nil, err
	}
	return &domain.Goal{Goal: &goal}, nil
}

func (s *goalStorage) ListGoals(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
	isInUseStatus *bool,
	environmentNamespace string,
) ([]*proto.Goal, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	prepareArgs := make([]interface{}, 0, len(whereArgs)+1)
	prepareArgs = append(prepareArgs, environmentNamespace)
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
	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			description,
			archived,
			deleted,
			created_at,
			updated_at,
			CASE 
				WHEN (
					SELECT 
						COUNT(1)
					FROM 
						experiment
					WHERE
						environment_namespace = ? AND
						goal_ids LIKE concat("%%", goal.id, "%%")
				) > 0 THEN TRUE 
				ELSE FALSE 
			END AS is_in_use_status
		FROM
			goal
		%s %s %s %s
		`, whereSQL, isInUseStatusSQL, orderBySQL, limitOffsetSQL,
	)
	rows, err := s.qe.QueryContext(ctx, query, prepareArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	goals := make([]*proto.Goal, 0, limit)
	for rows.Next() {
		goal := proto.Goal{}
		err := rows.Scan(
			&goal.Id,
			&goal.Name,
			&goal.Description,
			&goal.Archived,
			&goal.Deleted,
			&goal.CreatedAt,
			&goal.UpdatedAt,
			&goal.IsInUseStatus,
		)
		if err != nil {
			return nil, 0, 0, err
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
	countQuery := fmt.Sprintf(`
		SELECT
			COUNT(	
				CASE 
					WHEN (
						SELECT 
							COUNT(1)
						FROM 
							experiment
						WHERE
							environment_namespace = ? AND
							goal_ids LIKE concat("%%", goal.id, "%%")
					) %s
				END
			)
		FROM
			goal
		%s %s
		`, countConditionSQL, whereSQL, orderBySQL,
	)
	err = s.qe.QueryRowContext(ctx, countQuery, prepareArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return goals, nextOffset, totalCount, nil
}
