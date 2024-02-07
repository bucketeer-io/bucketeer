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
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestConvertPostgresError(t *testing.T) {
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
			input:    &pq.Error{Code: uniqueViolation},
			expected: ErrDuplicateEntry,
		},
		{
			desc:     "non mysql error",
			input:    errors.New("non postgres error"),
			expected: errors.New("non postgres error"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := convertPostgresError(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}
