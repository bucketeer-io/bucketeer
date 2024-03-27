// Copyright 2024 The Bucketeer Authors.
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

	"github.com/bucketeer-io/bucketeer/pkg/notification/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/notification"
)

var (
	ErrSubscriptionAlreadyExists          = errors.New("subscription: subscription already exists")
	ErrSubscriptionNotFound               = errors.New("subscription: subscription not found")
	ErrSubscriptionUnexpectedAffectedRows = errors.New("subscription: subscription unexpected affected rows")

	//go:embed sql/subscription/insert_subscription_v2.sql
	insertSubscriptionV2SQLQuery string
	//go:embed sql/subscription/update_subscription_v2.sql
	updateSubscriptionV2SQLQuery string
	//go:embed sql/subscription/delete_subscription_v2.sql
	deleteSubscriptionV2SQLQuery string
	//go:embed sql/subscription/select_subscription_v2.sql
	selectSubscriptionV2IDSQLQuery string
	//go:embed sql/subscription/select_subscription_v2_any.sql
	selectSubscriptionV2AnySQLQuery string
	//go:embed sql/subscription/select_subscription_v2_count.sql
	selectSubscriptionV2CountSQLQuery string
)

type SubscriptionStorage interface {
	CreateSubscription(ctx context.Context, e *domain.Subscription, environmentNamespace string) error
	UpdateSubscription(ctx context.Context, e *domain.Subscription, environmentNamespace string) error
	DeleteSubscription(ctx context.Context, id, environmentNamespace string) error
	GetSubscription(ctx context.Context, id, environmentNamespace string) (*domain.Subscription, error)
	ListSubscriptions(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Subscription, int, int64, error)
}

type subscriptionStorage struct {
	qe mysql.QueryExecer
}

func NewSubscriptionStorage(qe mysql.QueryExecer) SubscriptionStorage {
	return &subscriptionStorage{qe}
}

func (s *subscriptionStorage) CreateSubscription(
	ctx context.Context,
	e *domain.Subscription,
	environmentNamespace string,
) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertSubscriptionV2SQLQuery,
		e.Id,
		e.CreatedAt,
		e.UpdatedAt,
		e.Disabled,
		mysql.JSONObject{Val: e.SourceTypes},
		mysql.JSONObject{Val: e.Recipient},
		e.Name,
		environmentNamespace,
	)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrSubscriptionAlreadyExists
		}
		return err
	}
	return nil
}

func (s *subscriptionStorage) UpdateSubscription(
	ctx context.Context,
	e *domain.Subscription,
	environmentNamespace string,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateSubscriptionV2SQLQuery,
		e.UpdatedAt,
		e.Disabled,
		mysql.JSONObject{Val: e.SourceTypes},
		mysql.JSONObject{Val: e.Recipient},
		e.Name,
		e.Id,
		environmentNamespace,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrSubscriptionUnexpectedAffectedRows
	}
	return nil
}

func (s *subscriptionStorage) DeleteSubscription(
	ctx context.Context,
	id, environmentNamespace string,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		deleteSubscriptionV2SQLQuery,
		id,
		environmentNamespace,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrSubscriptionUnexpectedAffectedRows
	}
	return nil
}

func (s *subscriptionStorage) GetSubscription(
	ctx context.Context,
	id, environmentNamespace string,
) (*domain.Subscription, error) {
	subscription := proto.Subscription{}
	err := s.qe.QueryRowContext(
		ctx,
		selectSubscriptionV2IDSQLQuery,
		id,
		environmentNamespace,
	).Scan(
		&subscription.Id,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
		&subscription.Disabled,
		&mysql.JSONObject{Val: &subscription.SourceTypes},
		&mysql.JSONObject{Val: &subscription.Recipient},
		&subscription.Name,
	)
	if err != nil {
		if err == mysql.ErrNoRows {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}
	return &domain.Subscription{Subscription: &subscription}, nil
}

func (s *subscriptionStorage) ListSubscriptions(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.Subscription, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(selectSubscriptionV2AnySQLQuery, whereSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	subscriptions := make([]*proto.Subscription, 0, limit)
	for rows.Next() {
		subscription := proto.Subscription{}
		err := rows.Scan(
			&subscription.Id,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
			&subscription.Disabled,
			&mysql.JSONObject{Val: &subscription.SourceTypes},
			&mysql.JSONObject{Val: &subscription.Recipient},
			&subscription.Name,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		subscriptions = append(subscriptions, &subscription)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(subscriptions)
	var totalCount int64
	countQuery := fmt.Sprintf(selectSubscriptionV2CountSQLQuery, whereSQL, orderBySQL)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return subscriptions, nextOffset, totalCount, nil
}
