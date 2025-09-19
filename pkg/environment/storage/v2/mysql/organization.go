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
	"fmt"
	"strconv"

	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
	mysqlstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

var (
	//go:embed sql/organization/insert_organization.sql
	insertOrganizationSQL string
	//go:embed sql/organization/update_organization.sql
	updateOrganizationSQL string
	//go:embed sql/organization/select_organization.sql
	selectOrganizationSQL string
	//go:embed sql/organization/select_system_admin_organization.sql
	selectSystemAdminOrganizationSQL string
	//go:embed sql/organization/select_organizations.sql
	selectOrganizationsSQL string
	//go:embed sql/organization/count_organizations.sql
	countOrganizationsSQL string
)

type organizationStorage struct {
	qe mysqlstorage.QueryExecer
}

func NewOrganizationStorage(qe mysqlstorage.QueryExecer) v2es.OrganizationStorage {
	return &organizationStorage{qe}
}

func (s *organizationStorage) CreateOrganization(ctx context.Context, o *domain.Organization) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertOrganizationSQL,
		o.Id,
		o.Name,
		o.OwnerEmail,
		o.UrlCode,
		o.Description,
		o.Disabled,
		o.Archived,
		o.Trial,
		o.SystemAdmin,
		o.PasswordAuthenticationEnabled,
		o.CreatedAt,
		o.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrDuplicateEntry) {
			return v2es.ErrOrganizationAlreadyExists
		}
		return err
	}
	return nil
}

func (s *organizationStorage) UpdateOrganization(ctx context.Context, o *domain.Organization) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateOrganizationSQL,
		o.Name,
		o.OwnerEmail,
		o.Description,
		o.Disabled,
		o.Archived,
		o.Trial,
		o.PasswordAuthenticationEnabled,
		o.CreatedAt,
		o.UpdatedAt,
		o.Id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return v2es.ErrOrganizationUnexpectedAffectedRows
	}
	return nil
}

func (s *organizationStorage) GetOrganization(ctx context.Context, id string) (*domain.Organization, error) {
	organization := proto.Organization{}
	err := s.qe.QueryRowContext(
		ctx,
		selectOrganizationSQL,
		id,
	).Scan(
		&organization.Id,
		&organization.Name,
		&organization.OwnerEmail,
		&organization.UrlCode,
		&organization.Description,
		&organization.Disabled,
		&organization.Archived,
		&organization.Trial,
		&organization.SystemAdmin,
		&organization.PasswordAuthenticationEnabled,
		&organization.CreatedAt,
		&organization.UpdatedAt,
		&organization.ProjectCount,
		&organization.EnvironmentCount,
		&organization.UserCount,
	)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrNoRows) {
			return nil, v2es.ErrOrganizationNotFound
		}
		return nil, err
	}
	return &domain.Organization{Organization: &organization}, nil
}

func (s *organizationStorage) GetSystemAdminOrganization(ctx context.Context) (*domain.Organization, error) {
	organization := proto.Organization{}
	err := s.qe.QueryRowContext(
		ctx,
		selectSystemAdminOrganizationSQL,
	).Scan(
		&organization.Id,
		&organization.Name,
		&organization.OwnerEmail,
		&organization.UrlCode,
		&organization.Description,
		&organization.Disabled,
		&organization.Archived,
		&organization.Trial,
		&organization.SystemAdmin,
		&organization.PasswordAuthenticationEnabled,
		&organization.CreatedAt,
		&organization.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrNoRows) {
			return nil, v2es.ErrOrganizationNotFound
		}
		return nil, err
	}
	return &domain.Organization{Organization: &organization}, nil
}

func (s *organizationStorage) ListOrganizations(
	ctx context.Context,
	params v2es.ListOrganizationsParams,
) ([]*proto.Organization, int, int64, error) {
	options, err := listOrganizationsOptionsFromParams(params)
	if err != nil {
		return nil, 0, 0, err
	}
	whereParts := options.CreateWhereParts()
	whereSQL, whereArgs := mysqlstorage.ConstructWhereSQLString(whereParts)
	orderBySQL := mysqlstorage.ConstructOrderBySQLString(options.Orders)
	limitOffsetSQL := mysqlstorage.ConstructLimitOffsetSQLString(options.Limit, options.Offset)
	query := fmt.Sprintf(selectOrganizationsSQL, whereSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	organizations := make([]*proto.Organization, 0, options.Limit)
	for rows.Next() {
		organization := proto.Organization{}
		err := rows.Scan(
			&organization.Id,
			&organization.Name,
			&organization.OwnerEmail,
			&organization.UrlCode,
			&organization.Description,
			&organization.Disabled,
			&organization.Archived,
			&organization.Trial,
			&organization.SystemAdmin,
			&organization.PasswordAuthenticationEnabled,
			&organization.CreatedAt,
			&organization.UpdatedAt,
			&organization.ProjectCount,
			&organization.EnvironmentCount,
			&organization.UserCount,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		organizations = append(organizations, &organization)
	}
	if rows.Err() != nil {
		return nil, 0, 0, rows.Err()
	}
	nextOffset := options.Offset + len(organizations)
	var totalCount int64
	countQuery, countWhereArgs := mysqlstorage.ConstructCountQuery(countOrganizationsSQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return organizations, nextOffset, totalCount, nil
}

func listOrganizationsOptionsFromParams(p v2es.ListOrganizationsParams) (*mysqlstorage.ListOptions, error) {
	var filters []*mysqlstorage.FilterV2
	if p.Disabled != nil {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "organization.disabled",
			Operator: mysqlstorage.OperatorEqual,
			Value:    *p.Disabled,
		})
	}
	if p.Archived != nil {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "organization.archived",
			Operator: mysqlstorage.OperatorEqual,
			Value:    *p.Archived,
		})
	}

	var searchQuery *mysqlstorage.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &mysqlstorage.SearchQuery{
			Columns: []string{
				"organization.id",
				"organization.name",
				"organization.url_code",
			},
			Keyword: p.SearchKeyword,
		}
	}

	var column string
	switch p.OrderBy {
	case proto.ListOrganizationsRequest_DEFAULT,
		proto.ListOrganizationsRequest_NAME:
		column = "organization.name"
	case proto.ListOrganizationsRequest_URL_CODE:
		column = "organization.url_code"
	case proto.ListOrganizationsRequest_ID:
		column = "organization.id"
	case proto.ListOrganizationsRequest_CREATED_AT:
		column = "organization.created_at"
	case proto.ListOrganizationsRequest_UPDATED_AT:
		column = "organization.updated_at"
	case proto.ListOrganizationsRequest_ENVIRONMENT_COUNT:
		column = "environments"
	case proto.ListOrganizationsRequest_PROJECT_COUNT:
		column = "projects"
	case proto.ListOrganizationsRequest_USER_COUNT:
		column = "users"
	default:
		return nil, v2es.ErrInvalidOrderBy
	}
	direction := mysqlstorage.OrderDirectionAsc
	if p.OrderDirection == proto.ListOrganizationsRequest_DESC {
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
		SearchQuery: searchQuery,
		Orders:      []*mysqlstorage.Order{mysqlstorage.NewOrder(column, direction)},
	}, nil
}
