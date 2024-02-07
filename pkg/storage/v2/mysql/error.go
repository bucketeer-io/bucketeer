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
	"errors"

	"github.com/VividCortex/mysqlerr"
	libmysql "github.com/go-sql-driver/mysql"
)

var (
	ErrNoRows = errors.New("mysql: no rows")
	ErrTxDone = errors.New("mysql: tx done")

	// errors converted from MySQLError
	ErrDuplicateEntry = errors.New("mysql: duplicate entry")
)

func convertMySQLError(err error) error {
	if err == nil {
		return nil
	}
	if mysqlErr, ok := err.(*libmysql.MySQLError); ok {
		switch mysqlErr.Number {
		case mysqlerr.ER_DUP_ENTRY:
			return ErrDuplicateEntry
		}
	}
	return err
}
