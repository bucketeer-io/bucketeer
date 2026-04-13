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

func TestFilterBindSQL(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		input        *Filter
		expectedSQL  string
		expectedArgs []interface{}
	}{
		{
			desc:         "Empty",
			input:        &Filter{},
			expectedSQL:  "",
			expectedArgs: nil,
		},
		{
			desc: "Success",
			input: &Filter{
				Column:   "name",
				Operator: OperatorEqual,
				Value:    "feature",
			},
			expectedSQL:  "name = $1",
			expectedArgs: []interface{}{"feature"},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sql, args, _ := p.input.BindSQL(1)
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
		})
	}
}

func TestInFilterBindSQL(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		input        *InFilter
		expectedSQL  string
		expectedArgs []interface{}
	}{
		{
			desc:         "Empty",
			input:        &InFilter{},
			expectedSQL:  "",
			expectedArgs: nil,
		},
		{
			desc: "Success",
			input: &InFilter{
				Column: "name",
				Values: []interface{}{"v1", "v2"},
			},
			expectedSQL:  " name IN ($1, $2)",
			expectedArgs: []interface{}{"v1", "v2"},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sql, args, _ := p.input.BindSQL(1)
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
		})
	}
}

func TestNullFilterBindSQL(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		input        *NullFilter
		expectedSQL  string
		expectedArgs []interface{}
	}{
		{
			desc:         "Empty",
			input:        &NullFilter{},
			expectedSQL:  "",
			expectedArgs: nil,
		},
		{
			desc: "Success: null",
			input: &NullFilter{
				Column: "name",
				IsNull: true,
			},
			expectedSQL:  " name IS NULL",
			expectedArgs: nil,
		},
		{
			desc: "Success: not null",
			input: &NullFilter{
				Column: "name",
				IsNull: false,
			},
			expectedSQL:  " name IS NOT NULL",
			expectedArgs: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sql, args, _ := p.input.BindSQL(1)
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
		})
	}
}

func TestJSONFilterBindSQL(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		input        *JSONFilter
		expectedSQL  string
		expectedArgs []interface{}
	}{
		{
			desc:         "Empty",
			input:        &JSONFilter{},
			expectedSQL:  "",
			expectedArgs: nil,
		},
		{
			desc: "Success: JSONContainsNumber",
			input: &JSONFilter{
				Column: "enums",
				Func:   JSONContainsNumber,
				Values: []interface{}{1, 3},
			},
			expectedSQL:  "(enums::jsonb @> $1::jsonb)",
			expectedArgs: []interface{}{"[1, 3]"},
		},
		{
			desc: "Success: JSONContainsJSON",
			input: &JSONFilter{
				Column: "enums",
				Func:   JSONContainsJSON,
				Values: []interface{}{"{\"key1\":\"val1\", \"key2\":\"val2\"}"},
			},
			expectedSQL:  "(enums::jsonb @> $1::jsonb)",
			expectedArgs: []interface{}{"[{\"key1\":\"val1\", \"key2\":\"val2\"}]"},
		},
		{
			desc: "Success: JSONContainsString",
			input: &JSONFilter{
				Column: "enums",
				Func:   JSONContainsString,
				Values: []interface{}{"abc", "xyz"},
			},
			expectedSQL:  "(enums::jsonb @> $1::jsonb)",
			expectedArgs: []interface{}{`["abc", "xyz"]`},
		},
		{
			desc: "Success: JSONLengthGreaterThan empty",
			input: &JSONFilter{
				Column: "enums",
				Func:   JSONLengthGreaterThan,
				Values: []interface{}{},
			},
			expectedSQL:  "",
			expectedArgs: nil,
		},
		{
			desc: "Success: JSONLengthGreaterThan",
			input: &JSONFilter{
				Column: "enums",
				Func:   JSONLengthGreaterThan,
				Values: []interface{}{"1"},
			},
			expectedSQL:  "jsonb_array_length(enums::jsonb) > 1",
			expectedArgs: nil,
		},
		{
			desc: "Success: JSONLengthSmallerThan empty",
			input: &JSONFilter{
				Column: "enums",
				Func:   JSONLengthSmallerThan,
				Values: []interface{}{},
			},
			expectedSQL:  "",
			expectedArgs: nil,
		},
		{
			desc: "Success: JSONLengthSmallerThan",
			input: &JSONFilter{
				Column: "enums",
				Func:   JSONLengthSmallerThan,
				Values: []interface{}{"1"},
			},
			expectedSQL:  "jsonb_array_length(enums::jsonb) < 1",
			expectedArgs: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sql, args, _ := p.input.BindSQL(1)
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
		})
	}
}

func TestSearchQueryBindSQL(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		input        *SearchQuery
		expectedSQL  string
		expectedArgs []interface{}
	}{
		{
			desc:         "Empty",
			input:        &SearchQuery{},
			expectedSQL:  "",
			expectedArgs: nil,
		},
		{
			desc: "Success",
			input: &SearchQuery{
				Columns: []string{"id", "name"},
				Keyword: "test",
			},
			expectedSQL:  " (id LIKE $1 OR name LIKE $2)",
			expectedArgs: []interface{}{"%test%", "%test%"},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sql, args, _ := p.input.BindSQL(1)
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
		})
	}
}

func TestOrFilterBindSQL(t *testing.T) {
	t.Parallel()
	sql, args, _ := (&OrFilter{
		Queries: []WherePart{
			&Filter{Column: "a", Operator: OperatorEqual, Value: "1"},
			&Filter{Column: "b", Operator: OperatorEqual, Value: "2"},
		},
	}).BindSQL(1)
	assert.Equal(t, "(a = $1 OR b = $2)", sql)
	assert.Equal(t, []interface{}{"1", "2"}, args)

	sql, args, _ = (&OrFilter{}).BindSQL(1)
	assert.Equal(t, "", sql)
	assert.Nil(t, args)
}

func TestConstructWhereSQLString(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		input        []WherePart
		expectedSQL  string
		expectedArgs []interface{}
	}{
		{
			desc:         "Empty",
			input:        nil,
			expectedSQL:  "",
			expectedArgs: nil,
		},
		{
			desc: "Success",
			input: []WherePart{
				&Filter{Column: "name", Operator: OperatorEqual, Value: "feature"},
				&JSONFilter{Column: "enums", Func: JSONContainsNumber, Values: []interface{}{1, 3}},
			},
			expectedSQL:  " WHERE name = $1 AND (enums::jsonb @> $2::jsonb) ",
			expectedArgs: []interface{}{"feature", "[1, 3]"},
		},
		{
			desc: "multiple Filter sequential placeholders",
			input: []WherePart{
				&Filter{
					Column:   "feature.deleted",
					Operator: OperatorEqual,
					Value:    false,
				},
				&Filter{
					Column:   "feature.environment_id",
					Operator: OperatorEqual,
					Value:    "env-123",
				},
			},
			expectedSQL:  " WHERE feature.deleted = $1 AND feature.environment_id = $2 ",
			expectedArgs: []interface{}{false, "env-123"},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sql, args := ConstructWhereSQLString(p.input)
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
		})
	}
}

func TestConstructOrderBySQLString(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		input       []*Order
		expectedSQL string
	}{
		{
			desc:        "Empty",
			input:       nil,
			expectedSQL: "",
		},
		{
			desc: "Success",
			input: []*Order{
				NewOrder("created_at", OrderDirectionDesc),
				NewOrder("id", OrderDirectionAsc),
			},
			expectedSQL: " ORDER BY created_at DESC, id ASC ",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sql := ConstructOrderBySQLString(p.input)
			assert.Equal(t, p.expectedSQL, sql)
		})
	}
}

func TestConstructLimitOffsetSQLString(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		limit       int
		offset      int
		expectedSQL string
	}{
		{
			desc:        "no limit & no offset",
			limit:       0,
			offset:      0,
			expectedSQL: "",
		},
		{
			desc:        "no limit & offset",
			limit:       0,
			offset:      5,
			expectedSQL: " OFFSET 5",
		},
		{
			desc:        "limit & no offset",
			limit:       10,
			offset:      0,
			expectedSQL: " LIMIT 10",
		},
		{
			desc:        "limit & offset",
			limit:       10,
			offset:      5,
			expectedSQL: " LIMIT 10 OFFSET 5",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sql := ConstructLimitOffsetSQLString(p.limit, p.offset)
			assert.Equal(t, p.expectedSQL, sql)
		})
	}
}

func TestConstructCountQuery(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		baseQuery    string
		options      *ListOptions
		expectedSQL  string
		expectedArgs []interface{}
	}{
		{
			desc:         "Empty options",
			baseQuery:    "SELECT COUNT(1) FROM feature",
			options:      nil,
			expectedSQL:  "SELECT COUNT(1) FROM feature",
			expectedArgs: []interface{}{},
		},
		{
			desc:      "With filter",
			baseQuery: "SELECT COUNT(1) FROM feature",
			options: &ListOptions{
				Filters: []*Filter{
					{
						Column:   "name",
						Operator: OperatorEqual,
						Value:    "feature-1",
					},
				},
			},
			expectedSQL:  "SELECT COUNT(1) FROM feature WHERE name = $1 ",
			expectedArgs: []interface{}{"feature-1"},
		},
		{
			desc:      "With multiple filters",
			baseQuery: "SELECT COUNT(1) FROM feature",
			options: &ListOptions{
				Filters: []*Filter{
					{
						Column:   "name",
						Operator: OperatorEqual,
						Value:    "feature-1",
					},
					{
						Column:   "environment_id",
						Operator: OperatorEqual,
						Value:    "env-1",
					},
				},
			},
			expectedSQL:  "SELECT COUNT(1) FROM feature WHERE name = $1 AND environment_id = $2 ",
			expectedArgs: []interface{}{"feature-1", "env-1"},
		},
		{
			desc:      "With search query",
			baseQuery: "SELECT COUNT(1) FROM feature",
			options: &ListOptions{
				SearchQuery: &SearchQuery{
					Columns: []string{"name", "description"},
					Keyword: "test",
				},
			},
			expectedSQL:  "SELECT COUNT(1) FROM feature WHERE  (name LIKE $1 OR description LIKE $2) ",
			expectedArgs: []interface{}{"%test%", "%test%"},
		},
		{
			desc:      "With complex filters",
			baseQuery: "SELECT COUNT(1) FROM feature",
			options: &ListOptions{
				Filters: []*Filter{
					{
						Column:   "name",
						Operator: OperatorEqual,
						Value:    "feature-1",
					},
				},
				InFilters: []*InFilter{
					{
						Column: "environment_id",
						Values: []interface{}{"env-1", "env-2"},
					},
				},
				NullFilters: []*NullFilter{
					{
						Column: "deleted_at",
						IsNull: true,
					},
				},
			},
			expectedSQL:  "SELECT COUNT(1) FROM feature WHERE name = $1 AND  environment_id IN ($2, $3) AND  deleted_at IS NULL ",
			expectedArgs: []interface{}{"feature-1", "env-1", "env-2"},
		},
		{
			desc:      "Limit and offset are omitted from the query",
			baseQuery: "SELECT COUNT(1) FROM feature",
			options: &ListOptions{
				Limit:  10,
				Offset: 5,
			},
			expectedSQL:  "SELECT COUNT(1) FROM feature",
			expectedArgs: []interface{}{},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sql, args := ConstructCountQuery(p.baseQuery, p.options)
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
		})
	}
}
