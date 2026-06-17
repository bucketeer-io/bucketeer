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
	"strconv"

	"github.com/bucketeer-io/bucketeer/v2/pkg/notification/domain"
	v2ns "github.com/bucketeer-io/bucketeer/v2/pkg/notification/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
)

var (
	//go:embed sql/admin_subscription/insert_admin_subscription_v2.sql
	insertAdminSubscriptionV2SQLQuery string
	//go:embed sql/admin_subscription/update_admin_subscription_v2.sql
	updateAdminSubscriptionV2SQLQuery string
	//go:embed sql/admin_subscription/delete_admin_subscription_v2.sql
	deleteAdminSubscriptionV2SQLQuery string
	//go:embed sql/admin_subscription/select_admin_subscription_v2_any.sql
	selectAdminSubscriptionV2AnySQLQuery string
	//go:embed sql/admin_subscription/select_admin_subscription_v2.sql
	selectAdminSubscriptionV2SQLQuery string
	//go:embed sql/admin_subscription/select_admin_subscription_v2_count.sql
	selectAdminSubscriptionV2CountSQLQuery string
)

type adminSubscriptionStorage struct {
	qe postgres.QueryExecer
}

func NewAdminSubscriptionStorage(qe postgres.QueryExecer) v2ns.AdminSubscriptionStorage {
	return &adminSubscriptionStorage{qe}
}

func (s *adminSubscriptionStorage) CreateAdminSubscription(ctx context.Context, e *domain.Subscription) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertAdminSubscriptionV2SQLQuery,
		e.Id,
		e.CreatedAt,
		e.UpdatedAt,
		e.Disabled,
		postgres.JSONObject{Val: e.SourceTypes},
		postgres.JSONObject{Val: e.Recipient},
		e.Name,
	)
	if err != nil {
		if errors.Is(err, postgres.ErrDuplicateEntry) {
			return v2ns.ErrAdminSubscriptionAlreadyExists
		}
		return err
	}
	return nil
}

func (s *adminSubscriptionStorage) UpdateAdminSubscription(ctx context.Context, e *domain.Subscription) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateAdminSubscriptionV2SQLQuery,
		e.UpdatedAt,
		e.Disabled,
		postgres.JSONObject{Val: e.SourceTypes},
		postgres.JSONObject{Val: e.Recipient},
		e.Name,
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
		return v2ns.ErrAdminSubscriptionUnexpectedAffectedRows
	}
	return nil
}

func (s *adminSubscriptionStorage) DeleteAdminSubscription(ctx context.Context, id string) error {
	result, err := s.qe.ExecContext(
		ctx,
		deleteAdminSubscriptionV2SQLQuery,
		id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return v2ns.ErrAdminSubscriptionUnexpectedAffectedRows
	}
	return nil
}

func (s *adminSubscriptionStorage) GetAdminSubscription(ctx context.Context, id string) (*domain.Subscription, error) {
	subscription := proto.Subscription{}
	err := s.qe.QueryRowContext(
		ctx,
		selectAdminSubscriptionV2SQLQuery,
		id,
	).Scan(
		&subscription.Id,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
		&subscription.Disabled,
		&postgres.JSONObject{Val: &subscription.SourceTypes},
		&postgres.JSONObject{Val: &subscription.Recipient},
		&subscription.Name,
	)
	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return nil, v2ns.ErrAdminSubscriptionNotFound
		}
		return nil, err
	}
	return &domain.Subscription{Subscription: &subscription}, nil
}

func (s *adminSubscriptionStorage) ListAdminSubscriptions(
	ctx context.Context,
	params v2ns.ListAdminSubscriptionsParams,
) ([]*proto.Subscription, int, int64, error) {
	options, err := listAdminSubscriptionsOptions(params)
	if err != nil {
		return nil, 0, 0, err
	}
	query, whereArgs := postgres.ConstructQueryAndWhereArgs(selectAdminSubscriptionV2AnySQLQuery, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	subscriptions := make([]*proto.Subscription, 0, options.Limit)
	for rows.Next() {
		subscription := proto.Subscription{}
		err := rows.Scan(
			&subscription.Id,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
			&subscription.Disabled,
			&postgres.JSONObject{Val: &subscription.SourceTypes},
			&postgres.JSONObject{Val: &subscription.Recipient},
			&subscription.Name,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		subscriptions = append(subscriptions, &subscription)
	}
	if rows.Err() != nil {
		return nil, 0, 0, rows.Err()
	}
	nextOffset := options.Offset + len(subscriptions)
	var totalCount int64
	countQuery, countWhereArgs := postgres.ConstructCountQuery(selectAdminSubscriptionV2CountSQLQuery, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return subscriptions, nextOffset, totalCount, nil
}

func listAdminSubscriptionsOptions(params v2ns.ListAdminSubscriptionsParams) (*postgres.ListOptions, error) {
	var filters []*postgres.Filter
	if params.Disabled != nil {
		filters = append(filters, &postgres.Filter{
			Column:   "disabled",
			Operator: postgres.OperatorEqual,
			Value:    *params.Disabled,
		})
	}
	var jsonFilters []*postgres.JSONFilter
	if len(params.SourceTypes) > 0 {
		sourceTypesValues := make([]interface{}, len(params.SourceTypes))
		for i, st := range params.SourceTypes {
			sourceTypesValues[i] = int32(st)
		}
		jsonFilters = append(jsonFilters, &postgres.JSONFilter{
			Column: "source_types",
			Func:   postgres.JSONContainsNumber,
			Values: sourceTypesValues,
		})
	}
	var searchQuery *postgres.SearchQuery
	if params.SearchKeyword != "" {
		searchQuery = &postgres.SearchQuery{
			Columns: []string{"name"},
			Keyword: params.SearchKeyword,
		}
	}
	orders, err := listAdminSubscriptionsOrders(params.OrderBy, params.OrderDirection)
	if err != nil {
		return nil, err
	}
	cursor := params.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, v2ns.ErrInvalidCursor
	}
	return &postgres.ListOptions{
		Limit:       int(params.PageSize),
		Offset:      offset,
		Filters:     filters,
		JSONFilters: jsonFilters,
		SearchQuery: searchQuery,
		Orders:      orders,
	}, nil
}

func listAdminSubscriptionsOrders(
	orderBy proto.ListAdminSubscriptionsRequest_OrderBy,
	orderDirection proto.ListAdminSubscriptionsRequest_OrderDirection,
) ([]*postgres.Order, error) {
	var column string
	switch orderBy {
	case proto.ListAdminSubscriptionsRequest_DEFAULT,
		proto.ListAdminSubscriptionsRequest_NAME:
		column = "name"
	case proto.ListAdminSubscriptionsRequest_CREATED_AT:
		column = "created_at"
	case proto.ListAdminSubscriptionsRequest_UPDATED_AT:
		column = "updated_at"
	default:
		return nil, v2ns.ErrInvalidOrderBy
	}
	direction := postgres.OrderDirectionAsc
	if orderDirection == proto.ListAdminSubscriptionsRequest_DESC {
		direction = postgres.OrderDirectionDesc
	}
	return []*postgres.Order{postgres.NewOrder(column, direction)}, nil
}
