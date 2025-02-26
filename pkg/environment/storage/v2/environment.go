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

	"github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
)

var (
	ErrEnvironmentAlreadyExists          = errors.New("environment: already exists")
	ErrEnvironmentNotFound               = errors.New("environment: not found")
	ErrEnvironmentUnexpectedAffectedRows = errors.New("environment: unexpected affected rows")

	//go:embed sql/environment/insert_environment.sql
	insertEnvironmentSQL string
	//go:embed sql/environment/update_environment.sql
	updateEnvironmentSQL string
	//go:embed sql/environment/select_environment.sql
	selectEnvironmentSQL string
	//go:embed sql/environment/select_environments.sql
	selectEnvironmentsSQL string
	//go:embed sql/environment/count_environments.sql
	countEnvironmentsSQL string
)

type EnvironmentStorage interface {
	CreateEnvironmentV2(ctx context.Context, e *domain.EnvironmentV2) error
	UpdateEnvironmentV2(ctx context.Context, e *domain.EnvironmentV2) error
	GetEnvironmentV2(ctx context.Context, id string) (*domain.EnvironmentV2, error)
	ListEnvironmentsV2(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.EnvironmentV2, int, int64, error)
}

type environmentStorage struct {
	qe mysql.QueryExecer
}

func NewEnvironmentStorage(qe mysql.QueryExecer) EnvironmentStorage {
	return &environmentStorage{qe}
}

func (s *environmentStorage) CreateEnvironmentV2(ctx context.Context, e *domain.EnvironmentV2) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertEnvironmentSQL,
		e.Id,
		e.Name,
		e.UrlCode,
		e.Description,
		e.ProjectId,
		e.OrganizationId,
		e.Archived,
		e.RequireComment,
		e.CreatedAt,
		e.UpdatedAt,
	)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrEnvironmentAlreadyExists
		}
		return err
	}
	return nil
}

func (s *environmentStorage) UpdateEnvironmentV2(ctx context.Context, e *domain.EnvironmentV2) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateEnvironmentSQL,
		e.Name,
		e.Description,
		e.Archived,
		e.RequireComment,
		e.CreatedAt,
		e.UpdatedAt,
		e.Id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrEnvironmentUnexpectedAffectedRows
	}
	return nil
}

func (s *environmentStorage) GetEnvironmentV2(ctx context.Context, id string) (*domain.EnvironmentV2, error) {
	environment := proto.EnvironmentV2{}
	err := s.qe.QueryRowContext(
		ctx,
		selectEnvironmentSQL,
		id,
	).Scan(
		&environment.Id,
		&environment.Name,
		&environment.UrlCode,
		&environment.Description,
		&environment.ProjectId,
		&environment.OrganizationId,
		&environment.Archived,
		&environment.RequireComment,
		&environment.CreatedAt,
		&environment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrEnvironmentNotFound
		}
		return nil, err
	}
	return &domain.EnvironmentV2{EnvironmentV2: &environment}, nil
}

func (s *environmentStorage) ListEnvironmentsV2(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.EnvironmentV2, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(selectEnvironmentsSQL, whereSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	environments := make([]*proto.EnvironmentV2, 0, limit)
	for rows.Next() {
		environment := proto.EnvironmentV2{}
		err := rows.Scan(
			&environment.Id,
			&environment.Name,
			&environment.UrlCode,
			&environment.Description,
			&environment.ProjectId,
			&environment.OrganizationId,
			&environment.Archived,
			&environment.RequireComment,
			&environment.CreatedAt,
			&environment.UpdatedAt,
			&environment.FeatureFlagCount,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		environments = append(environments, &environment)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(environments)
	var totalCount int64
	countQuery := fmt.Sprintf(countEnvironmentsSQL, whereSQL)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return environments, nextOffset, totalCount, nil
}
