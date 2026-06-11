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

	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

var (
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
	//go:embed sql/environment/list_auto_archive_enabled_environments.sql
	listAutoArchiveEnabledEnvironmentsSQL string
)

type environmentStorage struct {
	qe pgstorage.QueryExecer
}

func NewEnvironmentStorage(qe pgstorage.QueryExecer) v2es.EnvironmentStorage {
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
		e.AutoArchiveEnabled,
		e.AutoArchiveUnusedDays,
		e.AutoArchiveCheckCodeRefs,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrDuplicateEntry) {
			return v2es.ErrEnvironmentAlreadyExists
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
		e.AutoArchiveEnabled,
		e.AutoArchiveUnusedDays,
		e.AutoArchiveCheckCodeRefs,
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
		return v2es.ErrEnvironmentUnexpectedAffectedRows
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
		&environment.AutoArchiveEnabled,
		&environment.AutoArchiveUnusedDays,
		&environment.AutoArchiveCheckCodeRefs,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, v2es.ErrEnvironmentNotFound
		}
		return nil, err
	}
	return &domain.EnvironmentV2{EnvironmentV2: &environment}, nil
}

func (s *environmentStorage) ListEnvironmentsV2(
	ctx context.Context,
	params v2es.ListEnvironmentsV2Params,
) ([]*proto.EnvironmentV2, int, int64, error) {
	options, err := listEnvironmentsOptionsFromParams(params)
	if err != nil {
		return nil, 0, 0, err
	}
	whereParts := options.CreateWhereParts()
	whereSQL, whereArgs := pgstorage.ConstructWhereSQLString(whereParts)
	orderBySQL := pgstorage.ConstructOrderBySQLString(options.Orders)
	limitOffsetSQL := pgstorage.ConstructLimitOffsetSQLString(options.Limit, options.Offset)
	query := fmt.Sprintf(selectEnvironmentsSQL, whereSQL, orderBySQL, limitOffsetSQL)

	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	environments := make([]*proto.EnvironmentV2, 0, options.Limit)
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
			&environment.AutoArchiveEnabled,
			&environment.AutoArchiveUnusedDays,
			&environment.AutoArchiveCheckCodeRefs,
			&environment.FeatureFlagCount,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		environments = append(environments, &environment)
	}
	if rows.Err() != nil {
		return nil, 0, 0, rows.Err()
	}
	nextOffset := options.Offset + len(environments)
	var totalCount int64
	countQuery, countWhereArgs := pgstorage.ConstructCountQuery(countEnvironmentsSQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return environments, nextOffset, totalCount, nil
}

func (s *environmentStorage) ListAutoArchiveEnabledEnvironments(ctx context.Context) ([]*domain.EnvironmentV2, error) {
	rows, err := s.qe.QueryContext(ctx, listAutoArchiveEnabledEnvironmentsSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	environments := make([]*domain.EnvironmentV2, 0)
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
			&environment.AutoArchiveEnabled,
			&environment.AutoArchiveUnusedDays,
			&environment.AutoArchiveCheckCodeRefs,
		)
		if err != nil {
			return nil, err
		}
		environments = append(environments, &domain.EnvironmentV2{EnvironmentV2: &environment})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return environments, nil
}

func listEnvironmentsOptionsFromParams(p v2es.ListEnvironmentsV2Params) (*pgstorage.ListOptions, error) {
	var filters []*pgstorage.Filter
	if p.ProjectID != "" {
		filters = append(filters, &pgstorage.Filter{
			Column:   "environment_v2.project_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.ProjectID,
		})
	}
	if p.OrganizationID != "" {
		filters = append(filters, &pgstorage.Filter{
			Column:   "environment_v2.organization_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.OrganizationID,
		})
	}
	if p.Archived != nil {
		filters = append(filters, &pgstorage.Filter{
			Column:   "environment_v2.archived",
			Operator: pgstorage.OperatorEqual,
			Value:    *p.Archived,
		})
	}

	var searchQuery *pgstorage.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &pgstorage.SearchQuery{
			Columns: []string{
				"environment_v2.id",
				"environment_v2.name",
				"environment_v2.url_code",
				"environment_v2.description",
			},
			Keyword: p.SearchKeyword,
		}
	}

	var column string
	switch p.OrderBy {
	case proto.ListEnvironmentsV2Request_DEFAULT,
		proto.ListEnvironmentsV2Request_NAME:
		column = "environment_v2.name"
	case proto.ListEnvironmentsV2Request_ID:
		column = "environment_v2.id"
	case proto.ListEnvironmentsV2Request_URL_CODE:
		column = "environment_v2.url_code"
	case proto.ListEnvironmentsV2Request_CREATED_AT:
		column = "environment_v2.created_at"
	case proto.ListEnvironmentsV2Request_UPDATED_AT:
		column = "environment_v2.updated_at"
	case proto.ListEnvironmentsV2Request_FEATURE_COUNT:
		column = "feature_count"
	default:
		return nil, v2es.ErrInvalidOrderBy
	}
	direction := pgstorage.OrderDirectionAsc
	if p.OrderDirection == proto.ListEnvironmentsV2Request_DESC {
		direction = pgstorage.OrderDirectionDesc
	}

	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, v2es.ErrInvalidCursor
	}

	return &pgstorage.ListOptions{
		Limit:       p.PageSize,
		Offset:      offset,
		Filters:     filters,
		SearchQuery: searchQuery,
		Orders:      []*pgstorage.Order{pgstorage.NewOrder(column, direction)},
	}, nil
}
