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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v2

import (
	"context"
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
)

type EnvironmentStorage interface {
	CreateEnvironment(ctx context.Context, e *domain.Environment) error
	UpdateEnvironment(ctx context.Context, e *domain.Environment) error
	GetEnvironment(ctx context.Context, id string) (*domain.Environment, error)
	GetEnvironmentByNamespace(
		ctx context.Context,
		namespace string,
		deleted bool,
	) (*domain.Environment, error)
	ListEnvironments(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Environment, int, int64, error)
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

func (s *environmentStorage) CreateEnvironment(ctx context.Context, e *domain.Environment) error {
	query := `
		INSERT INTO environment (
			id,
			namespace,
			name,
			description,
			deleted,
			created_at,
			updated_at,
			project_id
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?
		)
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
		e.Id,
		e.Namespace,
		e.Name,
		e.Description,
		e.Deleted,
		e.CreatedAt,
		e.UpdatedAt,
		e.ProjectId,
	)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrEnvironmentAlreadyExists
		}
		return err
	}
	return nil
}

func (s *environmentStorage) UpdateEnvironment(ctx context.Context, e *domain.Environment) error {
	query := `
		UPDATE 
			environment
		SET
			namespace = ?,
			name = ?,
			description = ?,
			deleted = ?,
			created_at = ?,
			updated_at = ?,
			project_id = ?
		WHERE
			id = ?
	`
	result, err := s.qe.ExecContext(
		ctx,
		query,
		e.Namespace,
		e.Name,
		e.Description,
		e.Deleted,
		e.CreatedAt,
		e.UpdatedAt,
		e.ProjectId,
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

func (s *environmentStorage) GetEnvironment(ctx context.Context, id string) (*domain.Environment, error) {
	e := proto.Environment{}
	query := `
		SELECT
			id,
			namespace,
			name,
			description,
			deleted,
			created_at,
			updated_at,
			project_id
		FROM
			environment
		WHERE
			id = ?
	`
	err := s.qe.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&e.Id,
		&e.Namespace,
		&e.Name,
		&e.Description,
		&e.Deleted,
		&e.CreatedAt,
		&e.UpdatedAt,
		&e.ProjectId,
	)
	if err != nil {
		if err == mysql.ErrNoRows {
			return nil, ErrEnvironmentNotFound
		}
		return nil, err
	}
	return &domain.Environment{Environment: &e}, nil
}

func (s *environmentStorage) GetEnvironmentByNamespace(
	ctx context.Context,
	namespace string,
	deleted bool,
) (*domain.Environment, error) {
	e := proto.Environment{}
	query := `
		SELECT
			id,
			namespace,
			name,
			description,
			deleted,
			created_at,
			updated_at,
			project_id
		FROM
			environment
		WHERE
			namespace = ? AND
			deleted = ?
	`
	err := s.qe.QueryRowContext(
		ctx,
		query,
		namespace,
		deleted,
	).Scan(
		&e.Id,
		&e.Namespace,
		&e.Name,
		&e.Description,
		&e.Deleted,
		&e.CreatedAt,
		&e.UpdatedAt,
		&e.ProjectId,
	)
	if err != nil {
		if err == mysql.ErrNoRows {
			return nil, ErrEnvironmentNotFound
		}
		return nil, err
	}
	return &domain.Environment{Environment: &e}, nil

}

func (s *environmentStorage) ListEnvironments(ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.Environment, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(`
		SELECT
			id,
			namespace,
			name,
			description,
			deleted,
			created_at,
			updated_at,
			project_id
		FROM
			environment
		%s %s %s
		`, whereSQL, orderBySQL, limitOffsetSQL,
	)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	environments := make([]*proto.Environment, 0, limit)
	for rows.Next() {
		e := proto.Environment{}
		err := rows.Scan(
			&e.Id,
			&e.Namespace,
			&e.Name,
			&e.Description,
			&e.Deleted,
			&e.CreatedAt,
			&e.UpdatedAt,
			&e.ProjectId,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		environments = append(environments, &e)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(environments)
	var totalCount int64
	countQuery := fmt.Sprintf(`
		SELECT
			COUNT(1)
		FROM
			environment
		%s %s
		`, whereSQL, orderBySQL,
	)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return environments, nextOffset, totalCount, nil
}

func (s *environmentStorage) CreateEnvironmentV2(ctx context.Context, e *domain.EnvironmentV2) error {
	query := `
		INSERT INTO environment_v2 (
			id,
			name,
			url_code,
			description,
			project_id,
			archived,
			created_at,
			updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?
		)
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
		e.Id,
		e.Name,
		e.UrlCode,
		e.Description,
		e.ProjectId,
		e.Archived,
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
	query := `
		UPDATE 
			environment_v2
		SET
			name = ?,
			description = ?,
			archived = ?,
			created_at = ?,
			updated_at = ?
		WHERE
			id = ?
	`
	result, err := s.qe.ExecContext(
		ctx,
		query,
		e.Name,
		e.Description,
		e.Archived,
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
	e := proto.EnvironmentV2{}
	query := `
		SELECT
			id,
			name,
			url_code,
			description,
			project_id,
			archived,
			created_at,
			updated_at
		FROM
			environment_v2
		WHERE
			id = ?
	`
	err := s.qe.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&e.Id,
		&e.Name,
		&e.UrlCode,
		&e.Description,
		&e.ProjectId,
		&e.Archived,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		if err == mysql.ErrNoRows {
			return nil, ErrEnvironmentNotFound
		}
		return nil, err
	}
	return &domain.EnvironmentV2{EnvironmentV2: &e}, nil
}

func (s *environmentStorage) ListEnvironmentsV2(ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.EnvironmentV2, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			url_code,
			description,
			project_id,
			archived,
			created_at,
			updated_at
		FROM
			environment_v2
		%s %s %s
		`, whereSQL, orderBySQL, limitOffsetSQL,
	)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	environments := make([]*proto.EnvironmentV2, 0, limit)
	for rows.Next() {
		e := proto.EnvironmentV2{}
		err := rows.Scan(
			&e.Id,
			&e.Name,
			&e.UrlCode,
			&e.Description,
			&e.ProjectId,
			&e.Archived,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		environments = append(environments, &e)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(environments)
	var totalCount int64
	countQuery := fmt.Sprintf(`
		SELECT
			COUNT(1)
		FROM
			environment_v2
		%s %s
		`, whereSQL, orderBySQL,
	)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return environments, nextOffset, totalCount, nil
}
