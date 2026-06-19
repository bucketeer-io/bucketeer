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

	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	"github.com/bucketeer-io/bucketeer/v2/pkg/team/domain"
	teamstorage "github.com/bucketeer-io/bucketeer/v2/pkg/team/storage"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/team"
)

var (
	//go:embed sql/insert_team.sql
	insertTeamSQL string
	//go:embed sql/select_team.sql
	selectTeamSQL string
	//go:embed sql/select_team_by_name.sql
	selectTeamByNameSQL string
	//go:embed sql/select_teams.sql
	selectTeamsSQL string
	//go:embed sql/count_teams.sql
	countTeamsSQL string
	//go:embed sql/delete_team.sql
	deleteTeamSQL string
)

type teamStorage struct {
	qe pgstorage.QueryExecer
}

func NewTeamStorage(qe pgstorage.QueryExecer) teamstorage.TeamStorage {
	return &teamStorage{qe: qe}
}

func (t *teamStorage) UpsertTeam(ctx context.Context, team *domain.Team) error {
	_, err := t.qe.ExecContext(
		ctx,
		insertTeamSQL,
		team.Id,
		team.Name,
		team.Description,
		team.OrganizationId,
		team.CreatedAt,
		team.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (t *teamStorage) GetTeam(ctx context.Context, id, organizationID string) (*domain.Team, error) {
	team := proto.Team{}
	err := t.qe.QueryRowContext(
		ctx,
		selectTeamSQL,
		id,
		organizationID,
	).Scan(
		&team.Id,
		&team.Name,
		&team.Description,
		&team.CreatedAt,
		&team.UpdatedAt,
		&team.OrganizationId,
		&team.OrganizationName,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, teamstorage.ErrTeamNotFound
		}
		return nil, err
	}
	return &domain.Team{
		Team: &team,
	}, nil
}

func (t *teamStorage) GetTeamByName(ctx context.Context, name, organizationID string) (*domain.Team, error) {
	team := proto.Team{}
	err := t.qe.QueryRowContext(
		ctx,
		selectTeamByNameSQL,
		name,
		organizationID,
	).Scan(
		&team.Id,
		&team.Name,
		&team.Description,
		&team.CreatedAt,
		&team.UpdatedAt,
		&team.OrganizationId,
		&team.OrganizationName,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, teamstorage.ErrTeamNotFound
		}
		return nil, err
	}
	return &domain.Team{
		Team: &team,
	}, nil
}

func listTeamsOrders(
	orderBy proto.ListTeamsRequest_OrderBy,
	orderDirection proto.ListTeamsRequest_OrderDirection,
) ([]*pgstorage.Order, error) {
	var column string
	switch orderBy {
	case proto.ListTeamsRequest_DEFAULT,
		proto.ListTeamsRequest_NAME:
		column = "team.name"
	case proto.ListTeamsRequest_CREATED_AT:
		column = "team.created_at"
	case proto.ListTeamsRequest_UPDATED_AT:
		column = "team.updated_at"
	case proto.ListTeamsRequest_ORGANIZATION:
		column = "team.organization_id"
	default:
		return nil, teamstorage.ErrInvalidListTeamsOrderBy
	}
	direction := pgstorage.OrderDirectionAsc
	if orderDirection == proto.ListTeamsRequest_DESC {
		direction = pgstorage.OrderDirectionDesc
	}
	return []*pgstorage.Order{pgstorage.NewOrder(column, direction)}, nil
}

func (t *teamStorage) ListTeams(
	ctx context.Context,
	p teamstorage.ListTeamsParams,
) ([]*proto.Team, int, int64, error) {
	orders, err := listTeamsOrders(p.OrderBy, p.OrderDirection)
	if err != nil {
		return nil, 0, 0, err
	}
	filters := []*pgstorage.Filter{
		{
			Column:   "team.organization_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.OrganizationID,
		},
	}
	var searchQuery *pgstorage.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &pgstorage.SearchQuery{
			Columns: []string{"team.name", "team.description"},
			Keyword: p.SearchKeyword,
		}
	}
	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, 0, 0, teamstorage.ErrInvalidListTeamsCursor
	}
	options := &pgstorage.ListOptions{
		Filters:     filters,
		SearchQuery: searchQuery,
		Orders:      orders,
		Limit:       p.PageSize,
		Offset:      offset,
	}
	whereParts := options.CreateWhereParts()
	whereSQL, whereArgs := pgstorage.ConstructWhereSQLString(whereParts)
	orderBySQL := pgstorage.ConstructOrderBySQLString(options.Orders)
	limitOffsetSQL := pgstorage.ConstructLimitOffsetSQLString(options.Limit, options.Offset)
	query := fmt.Sprintf(selectTeamsSQL, whereSQL, orderBySQL, limitOffsetSQL)
	rows, err := t.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	teams := make([]*proto.Team, 0, p.PageSize)
	for rows.Next() {
		team := proto.Team{}
		err := rows.Scan(
			&team.Id,
			&team.Name,
			&team.Description,
			&team.CreatedAt,
			&team.UpdatedAt,
			&team.OrganizationId,
			&team.OrganizationName,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		teams = append(teams, &team)
	}
	if rows.Err() != nil {
		return nil, 0, 0, rows.Err()
	}
	nextOffset := offset + len(teams)
	var totalCount int64
	countQuery := fmt.Sprintf(countTeamsSQL, whereSQL)
	err = t.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return teams, nextOffset, totalCount, nil
}

func (t *teamStorage) DeleteTeam(ctx context.Context, id string) error {
	result, err := t.qe.ExecContext(ctx, deleteTeamSQL, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return teamstorage.ErrTeamUnexpectedAffectedRows
	}
	return nil
}
