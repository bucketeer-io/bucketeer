// Copyright 2023 The Bucketeer Authors.
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

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/account"
)

var (
	ErrAccountAlreadyExists          = errors.New("account: account already exists")
	ErrAccountNotFound               = errors.New("account: account not found")
	ErrAccountUnexpectedAffectedRows = errors.New("account: account unexpected affected rows")
)

type AccountStorage interface {
	CreateAccount(ctx context.Context, a *domain.Account, environmentNamespace string) error
	UpdateAccount(ctx context.Context, a *domain.Account, environmentNamespace string) error
	GetAccount(ctx context.Context, id, environmentNamespace string) (*domain.Account, error)
	ListAccounts(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Account, int, int64, error)
}

type accountStorage struct {
	qe mysql.QueryExecer
}

func NewAccountStorage(qe mysql.QueryExecer) AccountStorage {
	return &accountStorage{qe}
}

func (s *accountStorage) CreateAccount(ctx context.Context, a *domain.Account, environmentNamespace string) error {
	query := `
		INSERT INTO account (
			id,
			email,
			name,
			role,
			disabled,
			created_at,
			updated_at,
			deleted,
			environment_namespace
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
		a.Id,
		a.Email,
		a.Name,
		int32(a.Role),
		a.Disabled,
		a.CreatedAt,
		a.UpdatedAt,
		a.Deleted,
		environmentNamespace,
	)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrAccountAlreadyExists
		}
		return err
	}
	return nil
}

func (s *accountStorage) UpdateAccount(ctx context.Context, a *domain.Account, environmentNamespace string) error {
	query := `
		UPDATE 
			account
		SET
			email = ?,
			name = ?,
			role = ?,
			disabled = ?,
			created_at = ?,
			updated_at = ?,
			deleted = ?
		WHERE
			id = ? AND
			environment_namespace = ?
	`
	result, err := s.qe.ExecContext(
		ctx,
		query,
		a.Email,
		a.Name,
		int32(a.Role),
		a.Disabled,
		a.CreatedAt,
		a.UpdatedAt,
		a.Deleted,
		a.Id,
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
		return ErrAccountUnexpectedAffectedRows
	}
	return nil
}

func (s *accountStorage) GetAccount(ctx context.Context, id, environmentNamespace string) (*domain.Account, error) {
	account := proto.Account{}
	var role int32
	query := `
		SELECT
			id,
			email,
			name,
			role,
			disabled,
			created_at,
			updated_at,
			deleted
		FROM
			account
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
		&account.Id,
		&account.Email,
		&account.Name,
		&role,
		&account.Disabled,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.Deleted,
	)
	if err != nil {
		if err == mysql.ErrNoRows {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	account.Role = proto.Account_Role(role)
	return &domain.Account{Account: &account}, nil
}

func (s *accountStorage) ListAccounts(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.Account, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(`
		SELECT
			id,
			email,
			name,
			role,
			disabled,
			created_at,
			updated_at,
			deleted
		FROM
			account
		%s %s %s
		`, whereSQL, orderBySQL, limitOffsetSQL,
	)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	accounts := make([]*proto.Account, 0, limit)
	for rows.Next() {
		account := proto.Account{}
		var role int32
		err := rows.Scan(
			&account.Id,
			&account.Email,
			&account.Name,
			&role,
			&account.Disabled,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.Deleted,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		account.Role = proto.Account_Role(role)
		accounts = append(accounts, &account)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(accounts)
	var totalCount int64
	countQuery := fmt.Sprintf(`
		SELECT
			COUNT(1)
		FROM
			account
		%s %s
		`, whereSQL, orderBySQL,
	)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return accounts, nextOffset, totalCount, nil
}
