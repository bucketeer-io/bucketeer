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

// nolint:lll
//
//go:generate mockgen -source=$GOFILE -aux_files=github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql=client.go -package=mock -destination=./mock/$GOFILE
package mysql

import (
	"context"
	"database/sql"
)

type Transaction interface {
	QueryExecer
	Commit() error
	Rollback() error
}

type transaction struct {
	stx *sql.Tx
}

func (tx *transaction) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	var err error
	defer record()(operationExec, &err)
	sret, err := tx.stx.ExecContext(ctx, query, args...)
	err = convertMySQLError(err)
	return &result{sret}, err
}

func (tx *transaction) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	var err error
	defer record()(operationQuery, &err)
	srows, err := tx.stx.QueryContext(ctx, query, args...)
	return &rows{srows}, err
}

func (tx *transaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	var err error
	defer record()(operationQueryRow, &err)
	r := &row{tx.stx.QueryRowContext(ctx, query, args...)}
	err = r.Err()
	return r
}

func (tx *transaction) Commit() error {
	var err error
	defer record()(operationCommit, &err)
	err = tx.stx.Commit()
	if err == sql.ErrTxDone {
		err = ErrTxDone
	}
	return err
}

func (tx *transaction) Rollback() error {
	var err error
	defer record()(operationRollback, &err)
	err = tx.stx.Rollback()
	if err == sql.ErrTxDone {
		err = ErrTxDone
	}
	return err
}
