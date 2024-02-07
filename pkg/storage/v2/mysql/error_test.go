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
	"testing"

	"github.com/VividCortex/mysqlerr"
	libmysql "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestConvertMySQLError(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		input    error
		expected error
	}{
		{
			desc:     "nil",
			input:    nil,
			expected: nil,
		},
		{
			desc:     "mysql error: ErrDuplicateEntry",
			input:    &libmysql.MySQLError{Number: mysqlerr.ER_DUP_ENTRY},
			expected: ErrDuplicateEntry,
		},
		{
			desc:     "non mysql error",
			input:    errors.New("non mysql error"),
			expected: errors.New("non mysql error"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := convertMySQLError(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}
