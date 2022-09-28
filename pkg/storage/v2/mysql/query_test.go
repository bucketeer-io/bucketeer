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

package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterSQLString(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		input        *Filter
		expectedSQL  string
		expectedArgs []interface{}
	}{
		"Empty": {
			input:        &Filter{},
			expectedSQL:  "",
			expectedArgs: nil,
		},
		"Success": {
			input: &Filter{
				Column:   "name",
				Operator: "=",
				Value:    "feature",
			},
			expectedSQL:  "name = ?",
			expectedArgs: []interface{}{"feature"},
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			sql, args := p.input.SQLString()
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
		})
	}
}

func TestInFilterSQLString(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		input        *InFilter
		expectedSQL  string
		expectedArgs []interface{}
	}{
		"Empty": {
			input:        &InFilter{},
			expectedSQL:  "",
			expectedArgs: nil,
		},
		"Success": {
			input: &InFilter{
				Column: "name",
				Values: []interface{}{"v1", "v2"},
			},
			expectedSQL:  "name IN (?, ?)",
			expectedArgs: []interface{}{"v1", "v2"},
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			sql, args := p.input.SQLString()
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
		})
	}
}

func TestNullFilterSQLString(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		input        *NullFilter
		expectedSQL  string
		expectedArgs []interface{}
	}{
		"Empty": {
			input:        &NullFilter{},
			expectedSQL:  "",
			expectedArgs: nil,
		},
		"Success: null": {
			input: &NullFilter{
				Column: "name",
				IsNull: true,
			},
			expectedSQL:  "name IS NULL",
			expectedArgs: nil,
		},
		"Success: not null": {
			input: &NullFilter{
				Column: "name",
				IsNull: false,
			},
			expectedSQL:  "name IS NOT NULL",
			expectedArgs: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			sql, args := p.input.SQLString()
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
		})
	}
}

func TestJSONFilterSQLString(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		input        *JSONFilter
		expectedSQL  string
		expectedArgs []interface{}
	}{
		"Empty": {
			input:        &JSONFilter{},
			expectedSQL:  "",
			expectedArgs: nil,
		},
		"Success: JSONContainsNumber": {
			input: &JSONFilter{
				Column: "enums",
				Func:   JSONContainsNumber,
				Values: []interface{}{1, 3},
			},
			expectedSQL:  "JSON_CONTAINS(enums, ?)",
			expectedArgs: []interface{}{"[1, 3]"},
		},
		"Success: JSONContainsString": {
			input: &JSONFilter{
				Column: "enums",
				Func:   JSONContainsString,
				Values: []interface{}{"abc", "xyz"},
			},
			expectedSQL:  "JSON_CONTAINS(enums, ?)",
			expectedArgs: []interface{}{`["abc", "xyz"]`},
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			sql, args := p.input.SQLString()
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
		})
	}
}

func TestSearchQuerySQLString(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		input        *SearchQuery
		expectedSQL  string
		expectedArgs []interface{}
	}{
		"Empty": {
			input:        &SearchQuery{},
			expectedSQL:  "",
			expectedArgs: nil,
		},
		"Success": {
			input: &SearchQuery{
				Columns: []string{"id", "name"},
				Keyword: "test",
			},
			expectedSQL:  "(id LIKE ? OR name LIKE ?)",
			expectedArgs: []interface{}{"%test%", "%test%"},
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			sql, args := p.input.SQLString()
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
		})
	}
}

func TestConstructWhereSQLString(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		input        []WherePart
		expectedSQL  string
		expectedArgs []interface{}
	}{
		"Empty": {
			input:        nil,
			expectedSQL:  "",
			expectedArgs: nil,
		},
		"Success": {
			input: []WherePart{
				NewFilter("name", "=", "feature"),
				NewJSONFilter("enums", JSONContainsNumber, []interface{}{1, 3}),
			},
			expectedSQL:  "WHERE name = ? AND JSON_CONTAINS(enums, ?)",
			expectedArgs: []interface{}{"feature", "[1, 3]"},
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			sql, args := ConstructWhereSQLString(p.input)
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
		})
	}
}

func TestConstructOrderBySQLString(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		input       []*Order
		expectedSQL string
	}{
		"Empty": {
			input:       nil,
			expectedSQL: "",
		},
		"Success": {
			input: []*Order{
				NewOrder("created_at", OrderDirectionDesc),
				NewOrder("id", OrderDirectionAsc),
			},
			expectedSQL: "ORDER BY created_at DESC, id ASC",
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			sql := ConstructOrderBySQLString(p.input)
			assert.Equal(t, p.expectedSQL, sql)
		})
	}
}

func TestConstructLimitOffsetSQLString(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		limit       int
		offset      int
		expectedSQL string
	}{
		"no limit & no offset": {
			limit:       0,
			offset:      0,
			expectedSQL: "",
		},
		"no limit & offset": {
			limit:       0,
			offset:      5,
			expectedSQL: "LIMIT 9223372036854775807 OFFSET 5",
		},
		"limit & no offset": {
			limit:       10,
			offset:      0,
			expectedSQL: "LIMIT 10",
		},
		"limit & offset": {
			limit:       10,
			offset:      5,
			expectedSQL: "LIMIT 10 OFFSET 5",
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			sql := ConstructLimitOffsetSQLString(p.limit, p.offset)
			assert.Equal(t, p.expectedSQL, sql)
		})
	}
}
