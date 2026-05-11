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

package v2

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strconv"

	"github.com/bucketeer-io/bucketeer/v2/pkg/push/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

var (
	//go:embed sql/postgres/insert_push.sql
	insertPushPostgres string
	//go:embed sql/postgres/update_push.sql
	updatePushPostgres string
	//go:embed sql/postgres/select_push.sql
	selectPushPostgres string
	//go:embed sql/postgres/list_pushes.sql
	listPushesPostgres string
	//go:embed sql/postgres/count_pushes.sql
	countPushesPostgres string
	//go:embed sql/postgres/delete_push.sql
	deletePushPostgres string
)

type postgresPushStorage struct {
	qe postgres.QueryExecer
}

// NewPostgresPushStorage returns push persistence backed by PostgreSQL.
func NewPostgresPushStorage(qe postgres.QueryExecer) PushStorage {
	return &postgresPushStorage{qe: qe}
}

func (s *postgresPushStorage) CreatePush(ctx context.Context, e *domain.Push, environmentId string) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertPushPostgres,
		e.Id,
		e.FcmServiceAccount,
		postgres.JSONObject{Val: e.Tags},
		e.Deleted,
		e.Name,
		e.CreatedAt,
		e.UpdatedAt,
		environmentId,
		e.Disabled,
	)
	if err != nil {
		if errors.Is(err, postgres.ErrDuplicateEntry) {
			return ErrPushAlreadyExists
		}
		return err
	}
	return nil
}

func (s *postgresPushStorage) UpdatePush(ctx context.Context, e *domain.Push, environmentId string) error {
	result, err := s.qe.ExecContext(
		ctx,
		updatePushPostgres,
		e.FcmServiceAccount,
		postgres.JSONObject{Val: e.Tags},
		e.Deleted,
		e.Name,
		e.CreatedAt,
		e.UpdatedAt,
		e.Disabled,
		e.Id,
		environmentId,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrPushUnexpectedAffectedRows
	}
	return nil
}

func (s *postgresPushStorage) GetPush(ctx context.Context, id, environmentId string) (*domain.Push, error) {
	push := proto.Push{}
	err := s.qe.QueryRowContext(
		ctx,
		selectPushPostgres,
		id,
		environmentId,
	).Scan(
		&push.Id,
		&push.FcmServiceAccount,
		&postgres.JSONObject{Val: &push.Tags},
		&push.Deleted,
		&push.Name,
		&push.CreatedAt,
		&push.UpdatedAt,
		&push.Disabled,
		&push.EnvironmentId,
		&push.EnvironmentName,
	)
	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return nil, ErrPushNotFound
		}
		return nil, err
	}
	return &domain.Push{Push: &push}, nil
}

func postgresListOptionsFromParams(p ListPushesParams) (*postgres.ListOptions, error) {
	var filters []*postgres.Filter
	var inFilters []*postgres.InFilter
	if p.OrganizationID != "" {
		filters = append(filters, &postgres.Filter{
			Column:   "env.organization_id",
			Operator: postgres.OperatorEqual,
			Value:    p.OrganizationID,
		})
		if len(p.EnvironmentIDs) > 0 {
			envIDs := make([]interface{}, 0, len(p.EnvironmentIDs))
			for _, id := range p.EnvironmentIDs {
				envIDs = append(envIDs, id)
			}
			inFilters = append(inFilters, &postgres.InFilter{
				Column: "push.environment_id",
				Values: envIDs,
			})
		}
	} else {
		if len(p.EnvironmentIDs) > 0 {
			filters = append(filters, &postgres.Filter{
				Column:   "push.environment_id",
				Operator: postgres.OperatorEqual,
				Value:    p.EnvironmentIDs[0],
			})
		}
	}
	if p.Disabled != nil {
		filters = append(filters, &postgres.Filter{
			Column:   "push.disabled",
			Operator: postgres.OperatorEqual,
			Value:    *p.Disabled,
		})
	}
	filters = append(filters, &postgres.Filter{
		Column:   "push.deleted",
		Operator: postgres.OperatorEqual,
		Value:    p.Deleted,
	})
	var searchQuery *postgres.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &postgres.SearchQuery{
			Columns: []string{"push.name"},
			Keyword: p.SearchKeyword,
		}
	}
	limit := int(p.PageSize)
	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, ErrInvalidListPushesCursor
	}
	orders, err := postgresOrdersFromParams(p)
	if err != nil {
		return nil, err
	}
	return &postgres.ListOptions{
		Limit:       limit,
		Offset:      offset,
		Filters:     filters,
		SearchQuery: searchQuery,
		InFilters:   inFilters,
		Orders:      orders,
		JSONFilters: nil,
		NullFilters: nil,
	}, nil
}

func postgresOrdersFromParams(p ListPushesParams) ([]*postgres.Order, error) {
	if p.OrderBy == nil {
		return nil, nil
	}
	var column string
	switch *p.OrderBy {
	case proto.ListPushesRequest_DEFAULT,
		proto.ListPushesRequest_NAME:
		column = "push.name"
	case proto.ListPushesRequest_CREATED_AT:
		column = "push.created_at"
	case proto.ListPushesRequest_UPDATED_AT:
		column = "push.updated_at"
	case proto.ListPushesRequest_ENVIRONMENT:
		column = "env.name"
	case proto.ListPushesRequest_STATE:
		column = "push.disabled"
	default:
		return nil, fmt.Errorf("list pushes: invalid order_by %v", *p.OrderBy)
	}
	direction := postgres.OrderDirectionAsc
	if p.OrderDirection == proto.ListPushesRequest_DESC {
		direction = postgres.OrderDirectionDesc
	}
	return []*postgres.Order{postgres.NewOrder(column, direction)}, nil
}

func (s *postgresPushStorage) ListPushes(
	ctx context.Context,
	p ListPushesParams,
) ([]*proto.Push, int, int64, error) {
	options, err := postgresListOptionsFromParams(p)
	if err != nil {
		return nil, 0, 0, err
	}
	query, whereArgs := postgres.ConstructQueryAndWhereArgs(listPushesPostgres, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	var offset int
	var limit int
	if options != nil {
		offset = options.Offset
		limit = options.Limit
	}

	defer rows.Close()
	pushes := make([]*proto.Push, 0, limit)
	for rows.Next() {
		push := proto.Push{}
		err := rows.Scan(
			&push.Id,
			&push.FcmServiceAccount,
			&postgres.JSONObject{Val: &push.Tags},
			&push.Deleted,
			&push.Name,
			&push.CreatedAt,
			&push.UpdatedAt,
			&push.Disabled,
			&push.EnvironmentId,
			&push.EnvironmentName,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		pushes = append(pushes, &push)
	}
	if rows.Err() != nil {
		return nil, 0, 0, rows.Err()
	}
	nextOffset := offset + len(pushes)
	var totalCount int64
	countQuery, countWhereArgs := postgres.ConstructCountQuery(countPushesPostgres, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return pushes, nextOffset, totalCount, nil
}

func (s *postgresPushStorage) DeletePush(ctx context.Context, id, environmentId string) error {
	result, err := s.qe.ExecContext(
		ctx,
		deletePushPostgres,
		id,
		environmentId,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrPushUnexpectedAffectedRows
	}
	return nil
}
