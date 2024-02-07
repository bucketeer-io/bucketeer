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
package mysql

import "database/sql"

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type result struct {
	sql.Result
}

type Row interface {
	Err() error
	Scan(dest ...interface{}) error
}

type row struct {
	srow *sql.Row
}

func (r *row) Err() error {
	err := r.srow.Err()
	if err == sql.ErrNoRows {
		return ErrNoRows
	}
	return err
}

func (r *row) Scan(dest ...interface{}) error {
	err := r.srow.Scan(dest...)
	if err == sql.ErrNoRows {
		return ErrNoRows
	}
	return err
}

type Rows interface {
	Close() error
	Err() error
	Next() bool
	Scan(dest ...interface{}) error
}

type rows struct {
	*sql.Rows
}
