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
	"encoding/json"
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
		{
			desc: "IN not supported on Filter",
			input: &Filter{
				Column:   "name",
				Operator: OperatorIn,
				Value:    "x",
			},
			expectedSQL:  "",
			expectedArgs: nil,
		},
		{
			desc: "NOT IN not supported on Filter",
			input: &Filter{
				Column:   "name",
				Operator: OperatorNotIn,
				Value:    "x",
			},
			expectedSQL:  "",
			expectedArgs: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sql, args, next := p.input.BindSQL(1)
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
			if p.desc == "IN not supported on Filter" || p.desc == "NOT IN not supported on Filter" {
				assert.Equal(t, 1, next, "placeholder index must not advance when fragment is empty")
			}
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
		{
			desc: "single value",
			input: &InFilter{
				Column: "status",
				Values: []interface{}{"active"},
			},
			expectedSQL:  " status IN ($1)",
			expectedArgs: []interface{}{"active"},
		},
		{
			desc: "three values",
			input: &InFilter{
				Column: "id",
				Values: []interface{}{10, 20, 30},
			},
			expectedSQL:  " id IN ($1, $2, $3)",
			expectedArgs: []interface{}{10, 20, 30},
		},
		{
			desc: "placeholders start at next index",
			input: &InFilter{
				Column: "environment_id",
				Values: []interface{}{"a", "b"},
			},
			expectedSQL:  " environment_id IN ($4, $5)",
			expectedArgs: []interface{}{"a", "b"},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			start := 1
			if p.desc == "placeholders start at next index" {
				start = 4
			}
			sql, args, next := p.input.BindSQL(start)
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
			assert.Equal(t, start+len(p.input.Values), next)
		})
	}
}

func TestNotInFilterBindSQL(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		input        *NotInFilter
		start        int
		expectedSQL  string
		expectedArgs []interface{}
	}{
		{
			desc:         "Empty",
			input:        &NotInFilter{},
			start:        1,
			expectedSQL:  "",
			expectedArgs: nil,
		},
		{
			desc: "two values",
			input: &NotInFilter{
				Column: "role",
				Values: []interface{}{"banned", "deleted"},
			},
			start:        1,
			expectedSQL:  " role NOT IN ($1, $2)",
			expectedArgs: []interface{}{"banned", "deleted"},
		},
		{
			desc: "single value",
			input: &NotInFilter{
				Column: "state",
				Values: []interface{}{99},
			},
			start:        2,
			expectedSQL:  " state NOT IN ($2)",
			expectedArgs: []interface{}{99},
		},
		{
			desc: "three values with offset start",
			input: &NotInFilter{
				Column: "id",
				Values: []interface{}{1, 2, 3},
			},
			start:        10,
			expectedSQL:  " id NOT IN ($10, $11, $12)",
			expectedArgs: []interface{}{1, 2, 3},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sql, args, next := p.input.BindSQL(p.start)
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
			if p.expectedSQL != "" {
				assert.Equal(t, p.start+len(p.input.Values), next)
			} else {
				assert.Equal(t, p.start, next)
			}
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
			expectedArgs: []interface{}{"[1,3]"},
		},
		{
			desc: "Success: JSONContainsJSON",
			input: &JSONFilter{
				Column: "enums",
				Func:   JSONContainsJSON,
				Values: []interface{}{`{"key1":"val1","key2":"val2"}`},
			},
			expectedSQL:  "(enums::jsonb @> $1::jsonb)",
			expectedArgs: []interface{}{`[{"key1":"val1","key2":"val2"}]`},
		},
		{
			desc: "Success: JSONContainsString",
			input: &JSONFilter{
				Column: "enums",
				Func:   JSONContainsString,
				Values: []interface{}{"abc", "xyz"},
			},
			expectedSQL:  "(enums::jsonb @> $1::jsonb)",
			expectedArgs: []interface{}{`["abc","xyz"]`},
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
			expectedSQL:  "jsonb_array_length(enums::jsonb) > $1",
			expectedArgs: []interface{}{"1"},
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
			expectedSQL:  "jsonb_array_length(enums::jsonb) < $1",
			expectedArgs: []interface{}{"1"},
		},
		{
			desc: "JSONContainsNumber: mixed numeric types",
			input: &JSONFilter{
				Column: "nums",
				Func:   JSONContainsNumber,
				Values: []interface{}{1, 2.5, -3},
			},
			expectedSQL:  "(nums::jsonb @> $1::jsonb)",
			expectedArgs: []interface{}{"[1,2.5,-3]"},
		},
		{
			desc: "JSONContainsJSON: multiple object literals",
			input: &JSONFilter{
				Column: "tags",
				Func:   JSONContainsJSON,
				Values: []interface{}{
					`{"a":1}`,
					`{"b":2}`,
				},
			},
			expectedSQL:  "(tags::jsonb @> $1::jsonb)",
			expectedArgs: []interface{}{`[{"a":1},{"b":2}]`},
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

func TestJSONFilterBindSQL_stringValuesWithJSONEscaping(t *testing.T) {
	t.Parallel()
	want, err := json.Marshal([]string{`he said "hi"`, "line1\nline2", `c:\tmp`})
	assert.NoError(t, err)
	f := &JSONFilter{
		Column: "labels",
		Func:   JSONContainsString,
		Values: []interface{}{`he said "hi"`, "line1\nline2", `c:\tmp`},
	}
	sql, args, next := f.BindSQL(3)
	assert.Equal(t, "(labels::jsonb @> $3::jsonb)", sql)
	assert.Equal(t, []interface{}{string(want)}, args)
	assert.Equal(t, 4, next)
}

func TestJSONFilterBindSQL_containsStringNonStringValue(t *testing.T) {
	t.Parallel()
	f := &JSONFilter{
		Column: "labels",
		Func:   JSONContainsString,
		Values: []interface{}{"ok", 42},
	}
	sql, args, next := f.BindSQL(1)
	assert.Empty(t, sql)
	assert.Nil(t, args)
	assert.Equal(t, 1, next)
}

func TestJSONFilterBindSQL_jsonContainsRawMessageAndMap(t *testing.T) {
	t.Parallel()
	raw := json.RawMessage(`{"x":true}`)
	inner, err := json.Marshal(map[string]int{"k": 1})
	assert.NoError(t, err)
	want, err := json.Marshal([]json.RawMessage{raw, json.RawMessage(inner)})
	assert.NoError(t, err)

	f := &JSONFilter{
		Column: "payload",
		Func:   JSONContainsJSON,
		Values: []interface{}{raw, map[string]int{"k": 1}},
	}
	sql, args, next := f.BindSQL(1)
	assert.Equal(t, "(payload::jsonb @> $1::jsonb)", sql)
	assert.Equal(t, []interface{}{string(want)}, args)
	assert.Equal(t, 2, next)
}

func TestConstructWhereSQLString_INAndNOTIN(t *testing.T) {
	t.Parallel()
	sql, args := ConstructWhereSQLString([]WherePart{
		&InFilter{Column: "environment_id", Values: []interface{}{"e1", "e2"}},
		&NotInFilter{Column: "status", Values: []interface{}{0, -1}},
	})
	assert.Equal(t, " WHERE  environment_id IN ($1, $2) AND  status NOT IN ($3, $4) ", sql)
	assert.Equal(t, []interface{}{"e1", "e2", 0, -1}, args)
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

func TestExistsFilterBindSQL(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		input        *ExistsFilter
		expectedSQL  string
		expectedArgs []interface{}
	}{
		{
			desc:         "Empty",
			input:        &ExistsFilter{},
			expectedSQL:  "",
			expectedArgs: nil,
		},
		{
			desc: "exists",
			input: &ExistsFilter{
				Subquery: "SELECT 1 FROM auto_ops_rule WHERE feature_id = feature.id",
			},
			expectedSQL:  "EXISTS (SELECT 1 FROM auto_ops_rule WHERE feature_id = feature.id)",
			expectedArgs: nil,
		},
		{
			desc: "not exists",
			input: &ExistsFilter{
				Subquery:  "SELECT 1 FROM auto_ops_rule WHERE feature_id = feature.id",
				NotExists: true,
			},
			expectedSQL:  "NOT EXISTS (SELECT 1 FROM auto_ops_rule WHERE feature_id = feature.id)",
			expectedArgs: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sql, args, next := p.input.BindSQL(1)
			assert.Equal(t, p.expectedSQL, sql)
			assert.Equal(t, p.expectedArgs, args)
			assert.Equal(t, 1, next)
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

	sql, args, _ = (&OrFilter{
		Queries: []WherePart{
			&Filter{Column: "a", Operator: OperatorEqual, Value: "1"},
			&Filter{}, // empty fragment skipped
			&Filter{Column: "b", Operator: OperatorEqual, Value: "2"},
		},
	}).BindSQL(1)
	assert.Equal(t, "(a = $1 OR b = $2)", sql)
	assert.Equal(t, []interface{}{"1", "2"}, args)

	sql, args, _ = (&OrFilter{
		Queries: []WherePart{
			&Filter{},
			&Filter{},
		},
	}).BindSQL(1)
	assert.Equal(t, "", sql)
	assert.Nil(t, args)

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
			expectedArgs: []interface{}{"feature", "[1,3]"},
		},
		{
			desc: "skips empty where parts",
			input: []WherePart{
				&Filter{Column: "name", Operator: OperatorEqual, Value: "feature"},
				&Filter{}, // empty
				&Filter{Column: "id", Operator: OperatorEqual, Value: 42},
			},
			expectedSQL:  " WHERE name = $1 AND id = $2 ",
			expectedArgs: []interface{}{"feature", 42},
		},
		{
			desc:         "all parts empty yields no WHERE",
			input:        []WherePart{&Filter{}, &NullFilter{}},
			expectedSQL:  "",
			expectedArgs: nil,
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
			desc:      "With IN and NOT IN",
			baseQuery: "SELECT COUNT(1) FROM feature",
			options: &ListOptions{
				InFilters: []*InFilter{
					{Column: "environment_id", Values: []interface{}{"a", "b"}},
				},
				NotInFilters: []*NotInFilter{
					{Column: "id", Values: []interface{}{0, 1, 2}},
				},
			},
			expectedSQL:  "SELECT COUNT(1) FROM feature WHERE  environment_id IN ($1, $2) AND  id NOT IN ($3, $4, $5) ",
			expectedArgs: []interface{}{"a", "b", 0, 1, 2},
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
