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
	//go:embed sql/subscription/insert_subscription_v2.sql
	insertSubscriptionV2SQLQuery string
	//go:embed sql/subscription/update_subscription_v2.sql
	updateSubscriptionV2SQLQuery string
	//go:embed sql/subscription/delete_subscription_v2.sql
	deleteSubscriptionV2SQLQuery string
	//go:embed sql/subscription/select_subscription_v2.sql
	selectSubscriptionV2SQLQuery string
	//go:embed sql/subscription/select_subscription_v2_any.sql
	selectSubscriptionV2AnySQLQuery string
	//go:embed sql/subscription/select_subscription_v2_count.sql
	selectSubscriptionV2CountSQLQuery string
)

type subscriptionStorage struct {
	qe postgres.QueryExecer
}

func NewSubscriptionStorage(qe postgres.QueryExecer) v2ns.SubscriptionStorage {
	return &subscriptionStorage{qe}
}

func (s *subscriptionStorage) CreateSubscription(
	ctx context.Context,
	e *domain.Subscription,
	environmentId string,
) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertSubscriptionV2SQLQuery,
		e.Id,
		e.CreatedAt,
		e.UpdatedAt,
		e.Disabled,
		postgres.JSONObject{Val: e.SourceTypes},
		postgres.JSONObject{Val: e.Recipient},
		e.Name,
		postgres.JSONObject{Val: e.FeatureFlagTags},
		environmentId,
	)
	if err != nil {
		if errors.Is(err, postgres.ErrDuplicateEntry) {
			return v2ns.ErrSubscriptionAlreadyExists
		}
		return err
	}
	return nil
}

func (s *subscriptionStorage) UpdateSubscription(
	ctx context.Context,
	e *domain.Subscription,
	environmentId string,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateSubscriptionV2SQLQuery,
		e.UpdatedAt,
		e.Disabled,
		postgres.JSONObject{Val: e.SourceTypes},
		postgres.JSONObject{Val: e.Recipient},
		e.Name,
		postgres.JSONObject{Val: e.FeatureFlagTags},
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
		return v2ns.ErrSubscriptionUnexpectedAffectedRows
	}
	return nil
}

func (s *subscriptionStorage) DeleteSubscription(
	ctx context.Context,
	id, environmentId string,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		deleteSubscriptionV2SQLQuery,
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
		return v2ns.ErrSubscriptionUnexpectedAffectedRows
	}
	return nil
}

func (s *subscriptionStorage) GetSubscription(
	ctx context.Context,
	id, environmentId string,
) (*domain.Subscription, error) {
	subscription := proto.Subscription{}
	err := s.qe.QueryRowContext(
		ctx,
		selectSubscriptionV2SQLQuery,
		id,
		environmentId,
	).Scan(
		&subscription.Id,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
		&subscription.Disabled,
		&postgres.JSONObject{Val: &subscription.SourceTypes},
		&postgres.JSONObject{Val: &subscription.Recipient},
		&subscription.Name,
		&postgres.JSONObject{Val: &subscription.FeatureFlagTags},
		&subscription.EnvironmentId,
		&subscription.EnvironmentName,
	)
	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return nil, v2ns.ErrSubscriptionNotFound
		}
		return nil, err
	}
	return &domain.Subscription{Subscription: &subscription}, nil
}

func (s *subscriptionStorage) ListSubscriptions(
	ctx context.Context,
	params v2ns.ListSubscriptionsParams,
) ([]*proto.Subscription, int, int64, error) {
	options, err := listSubscriptionsOptions(params)
	if err != nil {
		return nil, 0, 0, err
	}
	query, whereArgs := postgres.ConstructQueryAndWhereArgs(selectSubscriptionV2AnySQLQuery, options)
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
			&postgres.JSONObject{Val: &subscription.FeatureFlagTags},
			&subscription.EnvironmentId,
			&subscription.EnvironmentName,
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
	countQuery, countWhereArgs := postgres.ConstructCountQuery(selectSubscriptionV2CountSQLQuery, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return subscriptions, nextOffset, totalCount, nil
}

func listSubscriptionsOptions(params v2ns.ListSubscriptionsParams) (*postgres.ListOptions, error) {
	var filters []*postgres.Filter
	var inFilters []*postgres.InFilter
	if params.OrganizationID != "" {
		// console v3
		filters = append(filters, &postgres.Filter{
			Column:   "env.organization_id",
			Operator: postgres.OperatorEqual,
			Value:    params.OrganizationID,
		})
		if len(params.EnvironmentIDs) > 0 {
			envIDs := make([]interface{}, 0, len(params.EnvironmentIDs))
			for _, id := range params.EnvironmentIDs {
				envIDs = append(envIDs, id)
			}
			inFilters = append(inFilters, &postgres.InFilter{
				Column: "sub.environment_id",
				Values: envIDs,
			})
		}
	} else if len(params.EnvironmentIDs) > 0 {
		// console v2
		filters = append(filters, &postgres.Filter{
			Column:   "sub.environment_id",
			Operator: postgres.OperatorEqual,
			Value:    params.EnvironmentIDs[0],
		})
	}
	if params.Disabled != nil {
		filters = append(filters, &postgres.Filter{
			Column:   "sub.disabled",
			Operator: postgres.OperatorEqual,
			Value:    params.Disabled,
		})
	}
	var searchQuery *postgres.SearchQuery
	if params.SearchKeyword != "" {
		searchQuery = &postgres.SearchQuery{
			Columns: []string{"sub.name"},
			Keyword: params.SearchKeyword,
		}
	}
	var jsonFilters []*postgres.JSONFilter
	if len(params.SourceTypes) > 0 {
		sourceTypesValues := make([]interface{}, len(params.SourceTypes))
		for i, st := range params.SourceTypes {
			sourceTypesValues[i] = int32(st)
		}
		jsonFilters = append(jsonFilters, &postgres.JSONFilter{
			Column: "sub.source_types",
			Func:   postgres.JSONContainsNumber,
			Values: sourceTypesValues,
		})
	}
	orders, err := listSubscriptionsOrders(params.OrderBy, params.OrderDirection)
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
		InFilters:   inFilters,
		JSONFilters: jsonFilters,
		SearchQuery: searchQuery,
		Orders:      orders,
	}, nil
}

func listSubscriptionsOrders(
	orderBy proto.ListSubscriptionsRequest_OrderBy,
	orderDirection proto.ListSubscriptionsRequest_OrderDirection,
) ([]*postgres.Order, error) {
	var column string
	switch orderBy {
	case proto.ListSubscriptionsRequest_DEFAULT,
		proto.ListSubscriptionsRequest_NAME:
		column = "sub.name"
	case proto.ListSubscriptionsRequest_CREATED_AT:
		column = "sub.created_at"
	case proto.ListSubscriptionsRequest_UPDATED_AT:
		column = "sub.updated_at"
	case proto.ListSubscriptionsRequest_ENVIRONMENT:
		column = "env.name"
	case proto.ListSubscriptionsRequest_STATE:
		column = "sub.disabled"
	default:
		return nil, v2ns.ErrInvalidOrderBy
	}
	direction := postgres.OrderDirectionAsc
	if orderDirection == proto.ListSubscriptionsRequest_DESC {
		direction = postgres.OrderDirectionDesc
	}
	return []*postgres.Order{postgres.NewOrder(column, direction)}, nil
}
