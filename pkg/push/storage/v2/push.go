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

	"github.com/bucketeer-io/bucketeer/pkg/push/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/push"
)

var (
	ErrPushAlreadyExists          = errors.New("push: push already exists")
	ErrPushNotFound               = errors.New("push: push not found")
	ErrPushUnexpectedAffectedRows = errors.New("push: push unexpected affected rows")
)

type PushStorage interface {
	CreatePush(ctx context.Context, e *domain.Push, environmentNamespace string) error
	UpdatePush(ctx context.Context, e *domain.Push, environmentNamespace string) error
	GetPush(ctx context.Context, id, environmentNamespace string) (*domain.Push, error)
	ListPushes(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Push, int, int64, error)
}

type pushStorage struct {
	qe mysql.QueryExecer
}

func NewPushStorage(qe mysql.QueryExecer) PushStorage {
	return &pushStorage{qe: qe}
}

func (s *pushStorage) CreatePush(ctx context.Context, e *domain.Push, environmentNamespace string) error {
	query := `
		INSERT INTO push (
			id,
			fcm_api_key,
			tags,
			deleted,
			name,
			created_at,
			updated_at,
			environment_namespace 
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?
		)
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
		e.Id,
		e.FcmApiKey,
		mysql.JSONObject{Val: e.Tags},
		e.Deleted,
		e.Name,
		e.CreatedAt,
		e.UpdatedAt,
		environmentNamespace,
	)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrPushAlreadyExists
		}
		return err
	}
	return nil
}

func (s *pushStorage) UpdatePush(ctx context.Context, e *domain.Push, environmentNamespace string) error {
	query := `
		UPDATE 
			push
		SET
			fcm_api_key = ?,
			tags = ?,
			deleted = ?,
			name = ?,
			created_at = ?,
			updated_at = ?
		WHERE
			id = ? AND
			environment_namespace = ?
	`
	result, err := s.qe.ExecContext(
		ctx,
		query,
		e.FcmApiKey,
		mysql.JSONObject{Val: e.Tags},
		e.Deleted,
		e.Name,
		e.CreatedAt,
		e.UpdatedAt,
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
		return ErrPushUnexpectedAffectedRows
	}
	return nil
}

func (s *pushStorage) GetPush(ctx context.Context, id, environmentNamespace string) (*domain.Push, error) {
	push := proto.Push{}
	query := `
		SELECT
			id,
			fcm_api_key,
			tags,
			deleted,
			name,
			created_at,
			updated_at
		FROM
			push
		WHERE
			id = ? AND
			environment_namespace = ?
	`
	err := s.qe.QueryRowContext(
		ctx,
		query,
		id,
		environmentNamespace,
	).Scan(
		&push.Id,
		&push.FcmApiKey,
		&mysql.JSONObject{Val: &push.Tags},
		&push.Deleted,
		&push.Name,
		&push.CreatedAt,
		&push.UpdatedAt,
	)
	if err != nil {
		if err == mysql.ErrNoRows {
			return nil, ErrPushNotFound
		}
		return nil, err
	}
	return &domain.Push{Push: &push}, nil
}

func (s *pushStorage) ListPushes(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.Push, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(`
		SELECT
			id,
			fcm_api_key,
			tags,
			deleted,
			name,
			created_at,
			updated_at
		FROM
			push
		%s %s %s
		`, whereSQL, orderBySQL, limitOffsetSQL,
	)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	pushes := make([]*proto.Push, 0, limit)
	for rows.Next() {
		push := proto.Push{}
		err := rows.Scan(
			&push.Id,
			&push.FcmApiKey,
			&mysql.JSONObject{Val: &push.Tags},
			&push.Deleted,
			&push.Name,
			&push.CreatedAt,
			&push.UpdatedAt,
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
	countQuery := fmt.Sprintf(`
		SELECT
			COUNT(1)
		FROM
			push
		%s %s
		`, whereSQL, orderBySQL,
	)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return pushes, nextOffset, totalCount, nil
}
