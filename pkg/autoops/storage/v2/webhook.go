// Copyright 2022 The Bucketeer Authors.
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

	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

var (
	ErrWebhookAlreadyExists          = errors.New("webhook: already exists")
	ErrWebhookNotFound               = errors.New("webhook: not found")
	ErrWebhookUnexpectedAffectedRows = errors.New("webhook: unexpected affected rows")
)

type WebhookStorage interface {
	CreateWebhook(ctx context.Context, webhook *domain.Webhook, environmentNamespace string) error
	UpdateWebhook(ctx context.Context, webhook *domain.Webhook, environmentNamespace string) error
	DeleteWebhook(ctx context.Context, id, environmentNamespace string) error
	GetWebhook(ctx context.Context, id, environmentNamespace string) (*domain.Webhook, error)
	ListWebhooks(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Webhook, int, int64, error)
}

type webhookStorage struct {
	qe mysql.QueryExecer
}

func NewWebhookStorage(
	qe mysql.QueryExecer,
) WebhookStorage {
	return &webhookStorage{qe: qe}
}

func (s *webhookStorage) CreateWebhook(
	ctx context.Context,
	webhook *domain.Webhook,
	environmentNamespace string,
) error {
	query := `
		INSERT INTO webhook (
			id,
			name,
			description,
			environment_namespace,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
		webhook.Id,
		webhook.Name,
		webhook.Description,
		environmentNamespace,
		webhook.CreatedAt,
		webhook.UpdatedAt,
	)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrWebhookAlreadyExists
		}
		return err
	}
	return nil
}

func (s *webhookStorage) UpdateWebhook(
	ctx context.Context,
	webhook *domain.Webhook,
	environmentNamespace string,
) error {
	query := `
		UPDATE 
			webhook
		SET
			name = ?,
			description = ?,
			updated_at = ?
		WHERE
			id = ? AND
			environment_namespace = ?
	`
	result, err := s.qe.ExecContext(
		ctx,
		query,
		webhook.Name,
		webhook.Description,
		webhook.UpdatedAt,
		webhook.Id,
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
		return ErrWebhookUnexpectedAffectedRows
	}
	return nil
}

func (s *webhookStorage) DeleteWebhook(
	ctx context.Context,
	id, environmentNamespace string,
) error {
	query := `
		DELETE FROM
			webhook
		WHERE
			id = ? AND
			environment_namespace = ?
	`
	result, err := s.qe.ExecContext(
		ctx,
		query,
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
		return ErrWebhookUnexpectedAffectedRows
	}
	return nil
}

func (s *webhookStorage) GetWebhook(
	ctx context.Context,
	id, environmentNamespace string,
) (*domain.Webhook, error) {
	query := `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM
			webhook
		WHERE
			id = ? AND
			environment_namespace = ?
	`
	webhook := proto.Webhook{}
	err := s.qe.QueryRowContext(
		ctx,
		query,
		id,
		environmentNamespace,
	).Scan(
		&webhook.Id,
		&webhook.Name,
		&webhook.Description,
		&webhook.CreatedAt,
		&webhook.UpdatedAt,
	)
	if err != nil {
		if err == mysql.ErrNoRows {
			return nil, ErrWebhookNotFound
		}
		return nil, err
	}
	return &domain.Webhook{Webhook: &webhook}, nil
}

func (s *webhookStorage) ListWebhooks(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.Webhook, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM
			webhook
		%s %s %s
		`, whereSQL, orderBySQL, limitOffsetSQL,
	)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	webhooks := make([]*proto.Webhook, 0, limit)
	for rows.Next() {
		webhook := proto.Webhook{}
		err := rows.Scan(
			&webhook.Id,
			&webhook.Name,
			&webhook.Description,
			&webhook.CreatedAt,
			&webhook.UpdatedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		webhooks = append(webhooks, &webhook)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(webhooks)
	countQuery := fmt.Sprintf(`
		SELECT
			COUNT(1)
		FROM
			webhook
		%s %s
		`, whereSQL, orderBySQL,
	)
	var totalCount int64
	if err := s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount); err != nil {
		return nil, 0, 0, err
	}
	return webhooks, nextOffset, totalCount, nil
}
