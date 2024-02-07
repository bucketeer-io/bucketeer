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

	"github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
)

var (
	ErrProjectAlreadyExists          = errors.New("project: already exists")
	ErrProjectNotFound               = errors.New("project: not found")
	ErrProjectUnexpectedAffectedRows = errors.New("project: unexpected affected rows")
)

type ProjectStorage interface {
	CreateProject(ctx context.Context, p *domain.Project) error
	UpdateProject(ctx context.Context, p *domain.Project) error
	GetProject(ctx context.Context, id string) (*domain.Project, error)
	GetTrialProjectByEmail(
		ctx context.Context,
		email string,
		disabled, trial bool,
	) (*domain.Project, error)
	ListProjects(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Project, int, int64, error)
}

type projectStorage struct {
	qe mysql.QueryExecer
}

func NewProjectStorage(qe mysql.QueryExecer) ProjectStorage {
	return &projectStorage{qe}
}

func (s *projectStorage) CreateProject(ctx context.Context, p *domain.Project) error {
	query := `
		INSERT INTO project (
			id,
			name,
			url_code,
			description,
			disabled,
			trial,
			creator_email,
			organization_id,
			created_at,
			updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
		p.Id,
		p.Name,
		p.UrlCode,
		p.Description,
		p.Disabled,
		p.Trial,
		p.CreatorEmail,
		p.OrganizationId,
		p.CreatedAt,
		p.UpdatedAt,
	)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrProjectAlreadyExists
		}
		return err
	}
	return nil
}

func (s *projectStorage) UpdateProject(ctx context.Context, p *domain.Project) error {
	query := `
		UPDATE 
			project
		SET
			name = ?,
			description = ?,
			disabled = ?,
			trial = ?,
			creator_email = ?,
			created_at = ?,
			updated_at = ?
		WHERE
			id = ?
	`
	result, err := s.qe.ExecContext(
		ctx,
		query,
		p.Name,
		p.Description,
		p.Disabled,
		p.Trial,
		p.CreatorEmail,
		p.CreatedAt,
		p.UpdatedAt,
		p.Id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrProjectUnexpectedAffectedRows
	}
	return nil
}

func (s *projectStorage) GetProject(ctx context.Context, id string) (*domain.Project, error) {
	project := proto.Project{}
	query := `
		SELECT
			id,
			name,
			url_code,
			description,
			disabled,
			trial,
			creator_email,
			organization_id,
			created_at,
			updated_at
		FROM
			project
		WHERE
			id = ?
	`
	err := s.qe.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&project.Id,
		&project.Name,
		&project.UrlCode,
		&project.Description,
		&project.Disabled,
		&project.Trial,
		&project.CreatorEmail,
		&project.OrganizationId,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		if err == mysql.ErrNoRows {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	return &domain.Project{Project: &project}, nil
}

func (s *projectStorage) GetTrialProjectByEmail(
	ctx context.Context,
	email string,
	disabled, trial bool,
) (*domain.Project, error) {
	project := proto.Project{}
	query := `
		SELECT
			id,
			name,
			url_code,
			description,
			disabled,
			trial,
			creator_email,
			organization_id,
			created_at,
			updated_at
		FROM
			project
		WHERE
			creator_email = ? AND
			disabled = ? AND
			trial = ?
	`
	err := s.qe.QueryRowContext(
		ctx,
		query,
		email,
		disabled,
		trial,
	).Scan(
		&project.Id,
		&project.Name,
		&project.UrlCode,
		&project.Description,
		&project.Disabled,
		&project.Trial,
		&project.CreatorEmail,
		&project.OrganizationId,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		if err == mysql.ErrNoRows {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	return &domain.Project{Project: &project}, nil

}

func (s *projectStorage) ListProjects(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.Project, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			url_code,
			description,
			disabled,
			trial,
			creator_email,
			organization_id,
			created_at,
			updated_at
		FROM
			project
		%s %s %s
		`, whereSQL, orderBySQL, limitOffsetSQL,
	)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	projects := make([]*proto.Project, 0, limit)
	for rows.Next() {
		project := proto.Project{}
		err := rows.Scan(
			&project.Id,
			&project.Name,
			&project.UrlCode,
			&project.Description,
			&project.Disabled,
			&project.Trial,
			&project.CreatorEmail,
			&project.OrganizationId,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		projects = append(projects, &project)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(projects)
	var totalCount int64
	countQuery := fmt.Sprintf(`
		SELECT
			COUNT(1)
		FROM
			project
		%s %s
		`, whereSQL, orderBySQL,
	)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return projects, nextOffset, totalCount, nil
}
