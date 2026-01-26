// Copyright 2026 The Bucketeer Authors.
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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWritePlaceHolder(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		template string
		start    int
		count    int
		expected string
	}{
		{
			desc:     "two placeholders start at 1",
			template: "($%d, TO_TIMESTAMP($%d))",
			start:    1,
			count:    2,
			expected: "($1, TO_TIMESTAMP($2))",
		},
		{
			desc:     "three placeholders start at 3",
			template: "($%d, $%d, $%d)",
			start:    3,
			count:    3,
			expected: "($3, $4, $5)",
		},
		{
			desc:     "zero placeholders",
			template: "()",
			start:    1,
			count:    0,
			expected: "()",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := WritePlaceHolder(p.template, p.start, p.count)
			assert.Equal(t, p.expected, actual)
		})
	}
}
