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
	ErrProjectAlreadyExists          = errors.New("project: already exists")
	ErrProjectNotFound               = errors.New("project: not found")
	ErrProjectUnexpectedAffectedRows = errors.New("project: unexpected affected rows")

	//go:embed sql/project/insert_project.sql
	insertProjectSQL string
	//go:embed sql/project/update_project.sql
	updateProjectSQL string
	//go:embed sql/project/select_project.sql
	selectProjectSQL string
	//go:embed sql/project/select_trial_project_by_email.sql
	selectTrialProjectByEmailSQL string
	//go:embed sql/project/select_projects.sql
	selectProjectsSQL string
	//go:embed sql/project/count_projects.sql
	countProjectsSQL string
	//go:embed sql/project/delete_projects.sql
	deleteProjectsSQL string
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
		options *mysql.ListOptions,
	) ([]*proto.Project, int, int64, error)
	DeleteProjects(ctx context.Context, whereParts []mysql.WherePart) error
}

type projectStorage struct {
	qe mysql.QueryExecer
}

func NewProjectStorage(qe mysql.QueryExecer) ProjectStorage {
	return &projectStorage{qe}
}

func (s *projectStorage) CreateProject(ctx context.Context, p *domain.Project) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertProjectSQL,
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
	result, err := s.qe.ExecContext(
		ctx,
		updateProjectSQL,
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
	err := s.qe.QueryRowContext(
		ctx,
		selectProjectSQL,
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
	err := s.qe.QueryRowContext(
		ctx,
		selectTrialProjectByEmailSQL,
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
	options *mysql.ListOptions,
) ([]*proto.Project, int, int64, error) {
	// We do not use ConstructQueryAndWhereArgs() here,
	// because select_projects.sql defines the variable strings in a complex constructed way.
	var query string
	var whereArgs []any
	if options != nil {
		var whereSQL string
		whereParts := options.CreateWhereParts()
		whereSQL, whereArgs = mysql.ConstructWhereSQLString(whereParts)
		orderBySQL := mysql.ConstructOrderBySQLString(options.Orders)
		limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(options.Limit, options.Offset)
		query = fmt.Sprintf(selectProjectsSQL, whereSQL, orderBySQL, limitOffsetSQL)
	} else {
		query = selectProjectsSQL
		whereArgs = []interface{}{}
	}
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	var limit, offset int
	if options != nil {
		limit = options.Limit
		offset = options.Offset
	}
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
			&project.EnvironmentCount,
			&project.FeatureFlagCount,
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
	countQuery, countWhereArgs := mysql.ConstructCountQuery(countProjectsSQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return projects, nextOffset, totalCount, nil
}

func (s *projectStorage) DeleteProjects(ctx context.Context, whereParts []mysql.WherePart) error {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	query := fmt.Sprintf(deleteProjectsSQL, whereSQL)
	_, err := s.qe.ExecContext(
		ctx,
		query,
		whereArgs...,
	)
	if err != nil {
		return err
	}
	return nil
}
