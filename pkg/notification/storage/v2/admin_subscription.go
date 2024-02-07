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
	"errors"
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/notification/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/notification"
)

var (
	ErrAdminSubscriptionAlreadyExists          = errors.New("subscription: admin subscription already exists")
	ErrAdminSubscriptionNotFound               = errors.New("subscription: admin subscription not found")
	ErrAdminSubscriptionUnexpectedAffectedRows = errors.New("subscription: admin subscription unexpected affected rows")
)

type AdminSubscriptionStorage interface {
	CreateAdminSubscription(ctx context.Context, e *domain.Subscription) error
	UpdateAdminSubscription(ctx context.Context, e *domain.Subscription) error
	DeleteAdminSubscription(ctx context.Context, id string) error
	GetAdminSubscription(ctx context.Context, id string) (*domain.Subscription, error)
	ListAdminSubscriptions(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Subscription, int, int64, error)
}

type adminSubscriptionStorage struct {
	qe mysql.QueryExecer
}

func NewAdminSubscriptionStorage(qe mysql.QueryExecer) AdminSubscriptionStorage {
	return &adminSubscriptionStorage{qe}
}

func (s *adminSubscriptionStorage) CreateAdminSubscription(ctx context.Context, e *domain.Subscription) error {
	query := `
		INSERT INTO admin_subscription (
			id,
			created_at,
			updated_at,
			disabled,
			source_types,
			recipient,
			name
		) VALUES (
			?, ?, ?, ?, ?, ?, ?
		)
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
		e.Id,
		e.CreatedAt,
		e.UpdatedAt,
		e.Disabled,
		mysql.JSONObject{Val: e.SourceTypes},
		mysql.JSONObject{Val: e.Recipient},
		e.Name,
	)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrAdminSubscriptionAlreadyExists
		}
		return err
	}
	return nil
}

func (s *adminSubscriptionStorage) UpdateAdminSubscription(ctx context.Context, e *domain.Subscription) error {
	query := `
		UPDATE 
			admin_subscription
		SET
			created_at = ?,
			updated_at = ?,
			disabled = ?,
			source_types = ?,
			recipient = ?,
			name = ?
		WHERE
			id = ?
	`
	result, err := s.qe.ExecContext(
		ctx,
		query,
		e.CreatedAt,
		e.UpdatedAt,
		e.Disabled,
		mysql.JSONObject{Val: e.SourceTypes},
		mysql.JSONObject{Val: e.Recipient},
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
		return ErrAdminSubscriptionUnexpectedAffectedRows
	}
	return nil
}

func (s *adminSubscriptionStorage) DeleteAdminSubscription(ctx context.Context, id string) error {
	query := `
		DELETE FROM 
			admin_subscription
		WHERE
			id = ?
	`
	result, err := s.qe.ExecContext(
		ctx,
		query,
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
		return ErrAdminSubscriptionUnexpectedAffectedRows
	}
	return nil
}

func (s *adminSubscriptionStorage) GetAdminSubscription(ctx context.Context, id string) (*domain.Subscription, error) {
	subscription := proto.Subscription{}
	query := `
		SELECT
			id,
			created_at,
			updated_at,
			disabled,
			source_types,
			recipient,
			name
		FROM
			admin_subscription
		WHERE
			id = ?
	`
	err := s.qe.QueryRowContext(
		ctx,
		query,
		id,
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
			return nil, ErrAdminSubscriptionNotFound
		}
		return nil, err
	}
	return &domain.Subscription{Subscription: &subscription}, nil
}

func (s *adminSubscriptionStorage) ListAdminSubscriptions(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.Subscription, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(`
		SELECT
			id,
			created_at,
			updated_at,
			disabled,
			source_types,
			recipient,
			name
		FROM
			admin_subscription
		%s %s %s
		`, whereSQL, orderBySQL, limitOffsetSQL,
	)
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
	countQuery := fmt.Sprintf(`
		SELECT
			COUNT(1)
		FROM
		admin_subscription
		%s %s
		`, whereSQL, orderBySQL,
	)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return subscriptions, nextOffset, totalCount, nil
}
