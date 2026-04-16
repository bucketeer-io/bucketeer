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
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

//go:embed sql/mysql/insert_push.sql
var insertPushMySQL string

//go:embed sql/mysql/update_push.sql
var updatePushMySQL string

//go:embed sql/mysql/select_push.sql
var selectPushMySQL string

//go:embed sql/mysql/list_pushes.sql
var listPushesMySQL string

//go:embed sql/mysql/count_pushes.sql
var countPushesMySQL string

//go:embed sql/mysql/delete_push.sql
var deletePushMySQL string

type mysqlPushStorage struct {
	qe mysql.QueryExecer
}

func (s *mysqlPushStorage) CreatePush(ctx context.Context, e *domain.Push, environmentId string) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertPushMySQL,
		e.Id,
		e.FcmServiceAccount,
		mysql.JSONObject{Val: e.Tags},
		e.Deleted,
		e.Name,
		e.CreatedAt,
		e.UpdatedAt,
		environmentId,
		e.Disabled,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrDuplicateEntry) {
			return ErrPushAlreadyExists
		}
		return err
	}
	return nil
}

func (s *mysqlPushStorage) UpdatePush(ctx context.Context, e *domain.Push, environmentId string) error {
	result, err := s.qe.ExecContext(
		ctx,
		updatePushMySQL,
		e.FcmServiceAccount,
		mysql.JSONObject{Val: e.Tags},
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

func (s *mysqlPushStorage) GetPush(ctx context.Context, id, environmentId string) (*domain.Push, error) {
	push := proto.Push{}
	err := s.qe.QueryRowContext(
		ctx,
		selectPushMySQL,
		id,
		environmentId,
	).Scan(
		&push.Id,
		&push.FcmServiceAccount,
		&mysql.JSONObject{Val: &push.Tags},
		&push.Deleted,
		&push.Name,
		&push.CreatedAt,
		&push.UpdatedAt,
		&push.Disabled,
		&push.EnvironmentId,
		&push.EnvironmentName,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrPushNotFound
		}
		return nil, err
	}
	return &domain.Push{Push: &push}, nil
}

func mysqlListOptionsFromParams(p ListPushesParams) (*mysql.ListOptions, error) {
	var filters []*mysql.FilterV2
	var inFilters []*mysql.InFilter
	if p.OrganizationID != "" {
		// console v3
		filters = append(filters, &mysql.FilterV2{
			Column:   "env.organization_id",
			Operator: mysql.OperatorEqual,
			Value:    p.OrganizationID,
		})
		if len(p.EnvironmentIDs) > 0 {
			envIDs := make([]interface{}, 0, len(p.EnvironmentIDs))
			for _, id := range p.EnvironmentIDs {
				envIDs = append(envIDs, id)
			}
			inFilters = append(inFilters, &mysql.InFilter{
				Column: "push.environment_id",
				Values: envIDs,
			})
		}
	} else {
		// console v2
		if len(p.EnvironmentIDs) > 0 {
			filters = append(filters, &mysql.FilterV2{
				Column:   "push.environment_id",
				Operator: mysql.OperatorEqual,
				Value:    p.EnvironmentIDs[0],
			})
		}
	}
	if p.Disabled != nil {
		filters = append(filters, &mysql.FilterV2{
			Column:   "push.disabled",
			Operator: mysql.OperatorEqual,
			Value:    *p.Disabled,
		})
	}
	filters = append(filters, &mysql.FilterV2{
		Column:   "push.deleted",
		Operator: mysql.OperatorEqual,
		Value:    p.Deleted,
	})
	var searchQuery *mysql.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &mysql.SearchQuery{
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
	orders, err := mysqlOrdersFromParams(p)
	if err != nil {
		return nil, err
	}
	return &mysql.ListOptions{
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

func mysqlOrdersFromParams(p ListPushesParams) ([]*mysql.Order, error) {
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
	direction := mysql.OrderDirectionAsc
	if p.OrderDirection == proto.ListPushesRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *mysqlPushStorage) ListPushes(
	ctx context.Context,
	p ListPushesParams,
) ([]*proto.Push, int, int64, error) {
	options, err := mysqlListOptionsFromParams(p)
	if err != nil {
		return nil, 0, 0, err
	}
	query, whereArgs := mysql.ConstructQueryAndWhereArgs(listPushesMySQL, options)
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
			&mysql.JSONObject{Val: &push.Tags},
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
		return nil, 0, 0, err
	}
	nextOffset := offset + len(pushes)
	var totalCount int64
	countQuery, countWhereArgs := mysql.ConstructCountQuery(countPushesMySQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return pushes, nextOffset, totalCount, nil
}

func (s *mysqlPushStorage) DeletePush(ctx context.Context, id, environmentId string) error {
	result, err := s.qe.ExecContext(
		ctx,
		deletePushMySQL,
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
