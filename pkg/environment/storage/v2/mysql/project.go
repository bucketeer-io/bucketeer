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
	"fmt"
	"strconv"

	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
	mysqlstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

var (
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
)

type projectStorage struct {
	qe mysqlstorage.QueryExecer
}

func NewProjectStorage(qe mysqlstorage.QueryExecer) v2es.ProjectStorage {
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
		if err == mysqlstorage.ErrDuplicateEntry {
			return v2es.ErrProjectAlreadyExists
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
		return v2es.ErrProjectUnexpectedAffectedRows
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
		if err == mysqlstorage.ErrNoRows {
			return nil, v2es.ErrProjectNotFound
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
		if err == mysqlstorage.ErrNoRows {
			return nil, v2es.ErrProjectNotFound
		}
		return nil, err
	}
	return &domain.Project{Project: &project}, nil
}

func (s *projectStorage) ListProjects(
	ctx context.Context,
	params v2es.ListProjectsParams,
) ([]*proto.Project, int, int64, error) {
	options, err := listProjectsOptionsFromParams(params)
	if err != nil {
		return nil, 0, 0, err
	}
	var query string
	var whereArgs []any
	if options != nil {
		var whereSQL string
		whereParts := options.CreateWhereParts()
		whereSQL, whereArgs = mysqlstorage.ConstructWhereSQLString(whereParts)
		orderBySQL := mysqlstorage.ConstructOrderBySQLString(options.Orders)
		limitOffsetSQL := mysqlstorage.ConstructLimitOffsetSQLString(options.Limit, options.Offset)
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
	projects := make([]*proto.Project, 0, options.Limit)
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
		return nil, 0, 0, rows.Err()
	}
	nextOffset := options.Offset + len(projects)
	var totalCount int64
	countQuery, countWhereArgs := mysqlstorage.ConstructCountQuery(countProjectsSQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return projects, nextOffset, totalCount, nil
}

func listProjectsOptionsFromParams(p v2es.ListProjectsParams) (*mysqlstorage.ListOptions, error) {
	var filters []*mysqlstorage.FilterV2
	if p.OrganizationID != "" {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "project.organization_id",
			Operator: mysqlstorage.OperatorEqual,
			Value:    p.OrganizationID,
		})
	}
	if p.Disabled != nil {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "project.disabled",
			Operator: mysqlstorage.OperatorEqual,
			Value:    *p.Disabled,
		})
	}

	var inFilters []*mysqlstorage.InFilter
	if len(p.OrganizationIDs) > 0 {
		values := make([]interface{}, len(p.OrganizationIDs))
		for i, id := range p.OrganizationIDs {
			values[i] = id
		}
		inFilters = append(inFilters, &mysqlstorage.InFilter{
			Column: "project.organization_id",
			Values: values,
		})
	}

	var searchQuery *mysqlstorage.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &mysqlstorage.SearchQuery{
			Columns: []string{
				"project.id",
				"project.name",
				"project.url_code",
				"project.creator_email",
			},
			Keyword: p.SearchKeyword,
		}
	}

	var column string
	switch p.OrderBy {
	case proto.ListProjectsV2Request_DEFAULT,
		proto.ListProjectsV2Request_NAME:
		column = "project.name"
	case proto.ListProjectsV2Request_URL_CODE:
		column = "project.url_code"
	case proto.ListProjectsV2Request_ID:
		column = "project.id"
	case proto.ListProjectsV2Request_CREATED_AT:
		column = "project.created_at"
	case proto.ListProjectsV2Request_UPDATED_AT:
		column = "project.updated_at"
	case proto.ListProjectsV2Request_ENVIRONMENT_COUNT:
		column = "environment_count"
	case proto.ListProjectsV2Request_FEATURE_COUNT:
		column = "feature_count"
	case proto.ListProjectsV2Request_CREATOR_EMAIL:
		column = "project.creator_email"
	default:
		return nil, v2es.ErrInvalidOrderBy
	}
	direction := mysqlstorage.OrderDirectionAsc
	if p.OrderDirection == proto.ListProjectsV2Request_DESC {
		direction = mysqlstorage.OrderDirectionDesc
	}

	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, v2es.ErrInvalidCursor
	}

	return &mysqlstorage.ListOptions{
		Limit:       p.PageSize,
		Offset:      offset,
		Filters:     filters,
		InFilters:   inFilters,
		SearchQuery: searchQuery,
		Orders:      []*mysqlstorage.Order{mysqlstorage.NewOrder(column, direction)},
	}, nil
}
