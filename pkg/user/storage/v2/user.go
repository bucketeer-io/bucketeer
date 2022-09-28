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

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/user/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/user"
)

var (
	ErrUserNotFound = errors.New("user: not found")
)

type UserStorage interface {
	GetUser(ctx context.Context, id, environmentNamespace string) (*domain.User, error)
	GetUsers(ctx context.Context, ids []string, environmentNamespace string) ([]*domain.User, error)
	UpsertUser(ctx context.Context, user *domain.User, environmentNamespace string) error
	ListUsers(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.User, int, error)
}

type userStorage struct {
	qe mysql.QueryExecer
}

func NewUserStorage(qe mysql.QueryExecer) UserStorage {
	return &userStorage{qe: qe}
}

func (s *userStorage) GetUser(ctx context.Context, id, environmentNamespace string) (*domain.User, error) {
	user := &proto.User{}
	query := `
		SELECT
			id,
			tagged_data,
			last_seen,
			created_at
		FROM
			user
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
		&user.Id,
		&mysql.JSONObject{Val: &user.TaggedData},
		&user.LastSeen,
		&user.CreatedAt,
	)
	if err != nil {
		if err == mysql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &domain.User{User: user}, nil
}

func (s *userStorage) GetUsers(ctx context.Context, ids []string, environmentNamespace string) ([]*domain.User, error) {
	inFilterIDs := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		inFilterIDs = append(inFilterIDs, id)
	}
	whereParts := []mysql.WherePart{
		mysql.NewInFilter("id", inFilterIDs),
		mysql.NewFilter("environment_namespace", "=", environmentNamespace),
	}
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	query := fmt.Sprintf(`
		SELECT
			id,
			tagged_data,
			last_seen,
			created_at
		FROM
			user
		%s
	`, whereSQL,
	)
	rows, err := s.qe.QueryContext(
		ctx,
		query,
		whereArgs...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	entries := make([]*proto.User, 0, len(ids))
	for rows.Next() {
		user := proto.User{}
		err := rows.Scan(
			&user.Id,
			&mysql.JSONObject{Val: &user.TaggedData},
			&user.LastSeen,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &user)
	}
	if rows.Err() != nil {
		return nil, err
	}
	// NOTE: If the performance matters, remove the following loop and return protos.
	domainUsers := make([]*domain.User, 0, len(entries))
	for _, e := range entries {
		domainUsers = append(domainUsers, &domain.User{User: e})
	}
	return domainUsers, nil
}

func (s *userStorage) UpsertUser(ctx context.Context, user *domain.User, environmentNamespace string) error {
	query := `
	INSERT INTO user (
		id,
		tagged_data,
		last_seen,
		created_at,
		environment_namespace
	) VALUES (
		?, ?, ?, ?, ?
	) ON DUPLICATE KEY UPDATE
		tagged_data = VALUES(tagged_data),
		last_seen = VALUES(last_seen)
`
	_, err := s.qe.ExecContext(
		ctx,
		query,
		user.Id,
		mysql.JSONObject{Val: user.TaggedData},
		user.LastSeen,
		user.CreatedAt,
		environmentNamespace,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *userStorage) ListUsers(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.User, int, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(`
		SELECT
			id,
			tagged_data,
			last_seen,
			created_at
		FROM
			user
		%s %s %s
		`, whereSQL, orderBySQL, limitOffsetSQL,
	)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	users := make([]*proto.User, 0, limit)
	for rows.Next() {
		user := proto.User{}
		err := rows.Scan(
			&user.Id,
			&mysql.JSONObject{Val: &user.TaggedData},
			&user.LastSeen,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}
	if rows.Err() != nil {
		return nil, 0, err
	}
	nextOffset := offset + len(users)
	return users, nextOffset, nil
}
