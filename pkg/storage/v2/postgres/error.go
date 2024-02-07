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

package postgres

import (
	"errors"

	"github.com/lib/pq"
)

var (
	ErrDuplicateEntry = errors.New("postgres: duplicate entry")
)

const uniqueViolation pq.ErrorCode = "23505"

func convertPostgresError(err error) error {
	if err == nil {
		return nil
	}
	if postgresErr, ok := err.(*pq.Error); ok {
		switch postgresErr.Code {
		case uniqueViolation:
			return ErrDuplicateEntry
		}
	}
	return err
}
