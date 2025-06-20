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
package storage

import (
	"context"
	_ "embed"
	"errors"

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/team/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/team"
)

var (
	ErrTeamNotFound               = errors.New("team: not found")
	ErrTeamUnexpectedAffectedRows = errors.New("team: unexpected affected rows")

	//go:embed sql/insert_team.sql
	insertTeamSQL string
	//go:embed sql/select_team.sql
	selectTeamSQL string
	//go:embed sql/select_teams.sql
	selectTeamsSQL string
	//go:embed sql/count_teams.sql
	countTeamsSQL string
	//go:embed sql/delete_team.sql
	deleteTeamSQL string
)

type TeamStorage interface {
	UpsertTeam(ctx context.Context, team *domain.Team) error
	GetTeam(ctx context.Context, id, organizationID string) (*domain.Team, error)
	ListTeams(
		ctx context.Context,
		options *mysql.ListOptions,
	) ([]*proto.Team, int, int64, error)
	DeleteTeam(ctx context.Context, id string) error
}

type teamStorage struct {
	qe mysql.QueryExecer
}

func NewTeamStorage(qe mysql.QueryExecer) TeamStorage {
	return &teamStorage{
		qe: qe,
	}
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
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}
	return &domain.Team{
		Team: &team,
	}, nil
}

func (t *teamStorage) ListTeams(
	ctx context.Context,
	options *mysql.ListOptions,
) ([]*proto.Team, int, int64, error) {
	query, whereArgs := mysql.ConstructQueryAndWhereArgs(selectTeamsSQL, options)

	rows, err := t.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	var limit, offset int
	if options != nil {
		limit = options.Limit
		offset = options.Offset
	}
	teams := make([]*proto.Team, 0, limit)
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
		return nil, 0, 0, err
	}
	nextOffset := offset + len(teams)
	var totalCount int64
	countQuery, countWhereArgs := mysql.ConstructCountQuery(countTeamsSQL, options)
	err = t.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
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
		return ErrTeamUnexpectedAffectedRows
	}
	return nil
}
