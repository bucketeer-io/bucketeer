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

package mysql

import (
	"context"
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

type SqlMockClient struct {
	db *sql.DB
}

func (s *SqlMockClient) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *SqlMockClient) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (s *SqlMockClient) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	var err error
	defer record()(operationExec, &err)
	sret, err := s.db.ExecContext(ctx, query, args...)
	err = convertMySQLError(err)
	return &result{sret}, err
}

func (s *SqlMockClient) Close() error {
	return s.db.Close()
}

func (s *SqlMockClient) BeginTx(ctx context.Context) (Transaction, error) {
	var err error
	defer record()(operationBeginTx, &err)
	stx, err := s.db.BeginTx(ctx, nil)
	return &transaction{stx}, err
}

func (s *SqlMockClient) RunInTransaction(ctx context.Context, tx Transaction, f func() error) error {
	var err error
	defer record()(operationRunInTransaction, &err)
	defer func() {
		if err != nil {
			tx.Rollback() // nolint:errcheck
		}
	}()
	if err = f(); err == nil {
		err = tx.Commit()
	}
	return err
}

func (s *SqlMockClient) TearDown() error {
	return s.db.Close()
}

func (s *SqlMockClient) NewRows(columns []string) *sqlmock.Rows {
	return sqlmock.NewRows(columns)
}

func NewSqlMockClient() (*SqlMockClient, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	return &SqlMockClient{db: db}, mock, nil
}
