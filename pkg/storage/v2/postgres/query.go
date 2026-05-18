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
	"fmt"
	"strings"
)

// WritePlaceHolder formats a placeholder template using sequential indices.
// Example: WritePlaceHolder("($%d, $%d)", 1, 2) -> "($1, $2)".
func WritePlaceHolder(template string, start, count int) string {
	args := make([]interface{}, 0, count)
	for i := 0; i < count; i++ {
		args = append(args, start+i)
	}
	return fmt.Sprintf(template, args...)
}

// WherePart builds PostgreSQL WHERE fragments using numbered placeholders ($1, $2, ...).
// ConstructWhereSQLString chains BindSQL with a running index so parameters stay unique.
type WherePart interface {
	BindSQL(next int) (sql string, args []interface{}, nextAfter int)
}

type Operator int

const (
	OperatorEqual = iota + 1
	OperatorNotEqual
	OperatorIn
	OperatorNotIn
	OperatorGreaterThan
	OperatorGreaterThanOrEqual
	OperatorLessThan
	OperatorLessThanOrEqual
	OperatorContains
)

var operatorMap = map[Operator]string{
	OperatorEqual:              "=",
	OperatorNotEqual:           "!=",
	OperatorIn:                 "IN",
	OperatorNotIn:              "NOT IN",
	OperatorGreaterThan:        ">",
	OperatorGreaterThanOrEqual: ">=",
	OperatorLessThan:           "<",
	OperatorLessThanOrEqual:    "<=",
}

type Filter struct {
	Column   string
	Operator Operator
	Value    interface{}
}

func (f *Filter) BindSQL(next int) (sql string, args []interface{}, nextAfter int) {
	if f.Column == "" || f.Operator < OperatorEqual || f.Operator > OperatorContains {
		return "", nil, next
	}
	// IN / NOT IN require multiple placeholders; use InFilter / NotInFilter instead.
	if f.Operator == OperatorIn || f.Operator == OperatorNotIn {
		return "", nil, next
	}
	if f.Operator == OperatorContains {
		b, err := json.Marshal(f.Value)
		if err != nil {
			return "", nil, next
		}
		sql = fmt.Sprintf("(%s::jsonb @> $%d::jsonb)", f.Column, next)
		args = append(args, string(b))
		return sql, args, next + 1
	}
	sql = fmt.Sprintf("%s %s $%d", f.Column, operatorMap[f.Operator], next)
	args = append(args, f.Value)
	return sql, args, next + 1
}

type InFilter struct {
	Column string
	Values []interface{}
}

func (f *InFilter) BindSQL(next int) (sql string, args []interface{}, nextAfter int) {
	if f.Column == "" || len(f.Values) == 0 {
		return "", nil, next
	}
	var sb strings.Builder
	sb.WriteString(" ")
	sb.WriteString(f.Column)
	sb.WriteString(" IN (")
	for i := range f.Values {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("$%d", next+i))
	}
	sb.WriteString(")")
	sql = sb.String()
	args = append(args, f.Values...)
	return sql, args, next + len(f.Values)
}

// NotInFilter builds "col NOT IN ($n,...)" with one placeholder per value.
type NotInFilter struct {
	Column string
	Values []interface{}
}

func (f *NotInFilter) BindSQL(next int) (sql string, args []interface{}, nextAfter int) {
	if f.Column == "" || len(f.Values) == 0 {
		return "", nil, next
	}
	var sb strings.Builder
	sb.WriteString(" ")
	sb.WriteString(f.Column)
	sb.WriteString(" NOT IN (")
	for i := range f.Values {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("$%d", next+i))
	}
	sb.WriteString(")")
	sql = sb.String()
	args = append(args, f.Values...)
	return sql, args, next + len(f.Values)
}

type NullFilter struct {
	Column string
	IsNull bool
}

func (f *NullFilter) BindSQL(next int) (sql string, args []interface{}, nextAfter int) {
	if f.Column == "" {
		return "", nil, next
	}
	var sb strings.Builder
	sb.WriteString(" ")
	if f.IsNull {
		sb.WriteString(fmt.Sprintf("%s IS NULL", f.Column))
	} else {
		sb.WriteString(fmt.Sprintf("%s IS NOT NULL", f.Column))
	}
	return sb.String(), nil, next
}

type JSONFilterFunc int

const (
	_ JSONFilterFunc = iota
	JSONContainsNumber
	JSONContainsString
	JSONLengthGreaterThan
	JSONLengthSmallerThan
	JSONContainsJSON
)

type JSONFilter struct {
	Column string
	Func   JSONFilterFunc
	Values []interface{}
}

func (f *JSONFilter) BindSQL(next int) (sql string, args []interface{}, nextAfter int) {
	if f.Column == "" {
		return "", nil, next
	}
	switch f.Func {
	case JSONContainsNumber:
		payload, err := json.Marshal(f.Values)
		if err != nil {
			return "", nil, next
		}
		sql = fmt.Sprintf("(%s::jsonb @> $%d::jsonb)", f.Column, next)
		args = append(args, string(payload))
		return sql, args, next + 1
	case JSONContainsString:
		strs := make([]string, 0, len(f.Values))
		for _, v := range f.Values {
			s, ok := v.(string)
			if !ok {
				return "", nil, next
			}
			strs = append(strs, s)
		}
		payload, err := json.Marshal(strs)
		if err != nil {
			return "", nil, next
		}
		sql = fmt.Sprintf("(%s::jsonb @> $%d::jsonb)", f.Column, next)
		args = append(args, string(payload))
		return sql, args, next + 1
	case JSONContainsJSON:
		elems := make([]json.RawMessage, 0, len(f.Values))
		for _, v := range f.Values {
			switch t := v.(type) {
			case string:
				elems = append(elems, json.RawMessage(t))
			case json.RawMessage:
				elems = append(elems, t)
			default:
				b, err := json.Marshal(t)
				if err != nil {
					return "", nil, next
				}
				elems = append(elems, json.RawMessage(b))
			}
		}
		payload, err := json.Marshal(elems)
		if err != nil {
			return "", nil, next
		}
		sql = fmt.Sprintf("(%s::jsonb @> $%d::jsonb)", f.Column, next)
		args = append(args, string(payload))
		return sql, args, next + 1
	case JSONLengthGreaterThan, JSONLengthSmallerThan:
		if len(f.Values) == 0 {
			return "", nil, next
		}
		op := ">"
		if f.Func == JSONLengthSmallerThan {
			op = "<"
		}
		sql = fmt.Sprintf("jsonb_array_length(%s::jsonb) %s $%d", f.Column, op, next)
		args = append(args, f.Values[0])
		return sql, args, next + 1
	default:
		return "", nil, next
	}
}

type SearchQuery struct {
	Columns []string
	Keyword string
}

func (q *SearchQuery) BindSQL(next int) (sql string, args []interface{}, nextAfter int) {
	if len(q.Columns) == 0 {
		return "", nil, next
	}
	var sb strings.Builder
	sb.WriteString(" (")
	for i, col := range q.Columns {
		if i != 0 {
			sb.WriteString(" OR ")
		}
		sb.WriteString(fmt.Sprintf("%s LIKE $%d", col, next+i))
		args = append(args, "%"+q.Keyword+"%")
	}
	sb.WriteString(")")
	return sb.String(), args, next + len(q.Columns)
}

type OrFilter struct {
	Queries []WherePart
}

func (f *OrFilter) BindSQL(next int) (sql string, args []interface{}, nextAfter int) {
	if len(f.Queries) == 0 {
		return "", nil, next
	}
	var sb strings.Builder
	cur := next
	wrote := false
	for _, q := range f.Queries {
		qs, qa, after := q.BindSQL(cur)
		if qs == "" {
			continue
		}
		if !wrote {
			sb.WriteString("(")
			wrote = true
		} else {
			sb.WriteString(" OR ")
		}
		sb.WriteString(qs)
		args = append(args, qa...)
		cur = after
	}
	if !wrote {
		return "", nil, next
	}
	sb.WriteString(")")
	return sb.String(), args, cur
}

func ConstructWhereSQLString(wps []WherePart) (sql string, args []interface{}) {
	if len(wps) == 0 {
		return "", nil
	}
	var sb strings.Builder
	next := 1
	first := true
	for _, wp := range wps {
		frag, fragArgs, after := wp.BindSQL(next)
		if frag == "" {
			continue
		}
		if first {
			sb.WriteString(" WHERE ")
			first = false
		} else {
			sb.WriteString(" AND ")
		}
		sb.WriteString(frag)
		args = append(args, fragArgs...)
		next = after
	}
	if first {
		return "", nil
	}
	sql = sb.String() + " "
	return sql, args
}

type OrderDirection int

const (
	OrderDirectionAsc OrderDirection = iota
	OrderDirectionDesc
)

func (o OrderDirection) String() string {
	switch o {
	case OrderDirectionAsc:
		return "ASC"
	case OrderDirectionDesc:
		return "DESC"
	default:
		return ""
	}
}

type Order struct {
	Column    string
	Direction OrderDirection
}

func NewOrder(column string, direction OrderDirection) *Order {
	return &Order{
		Column:    column,
		Direction: direction,
	}
}

func ConstructOrderBySQLString(orders []*Order) string {
	var sb strings.Builder
	if len(orders) == 0 {
		return ""
	}
	sb.WriteString(" ORDER BY ")
	for i, o := range orders {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(o.Column)
		sb.WriteString(" ")
		sb.WriteString(o.Direction.String())
	}
	return sb.String() + " "
}

func ConstructQueryAndWhereArgs(baseQuery string, options *ListOptions) (query string, whereArgs []interface{}) {
	if options != nil {
		whereParts := options.CreateWhereParts()
		whereQuery, whereArgs := ConstructWhereSQLString(whereParts)
		orderByQuery := ConstructOrderBySQLString(options.Orders)
		limitOffsetQuery := ConstructLimitOffsetSQLString(options.Limit, options.Offset)
		query = baseQuery + whereQuery + orderByQuery + limitOffsetQuery
		return query, whereArgs
	}
	return baseQuery, []interface{}{}
}

// ConstructCountQuery builds a count query with optional filtering.
func ConstructCountQuery(baseQuery string, options *ListOptions) (query string, whereArgs []interface{}) {
	if options != nil {
		whereQuery, whereArgs := ConstructWhereSQLString(options.CreateWhereParts())
		if whereArgs == nil {
			whereArgs = []interface{}{}
		}
		return baseQuery + whereQuery, whereArgs
	}
	return baseQuery, []interface{}{}
}

func ConstructLimitOffsetSQLString(limit, offset int) string {
	if limit == 0 && offset == 0 {
		return ""
	}
	if limit == 0 {
		return fmt.Sprintf(" OFFSET %d", offset)
	}
	if offset == 0 {
		return fmt.Sprintf(" LIMIT %d", limit)
	}
	return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
}

type ListOptions struct {
	Limit        int
	Filters      []*Filter
	InFilters    []*InFilter
	NotInFilters []*NotInFilter
	NullFilters  []*NullFilter
	JSONFilters  []*JSONFilter
	SearchQuery  *SearchQuery
	OrFilters    []*OrFilter
	Orders       []*Order
	Offset       int
}

func (lo *ListOptions) CreateWhereParts() []WherePart {
	var whereParts []WherePart
	if lo.Filters != nil {
		for _, f := range lo.Filters {
			whereParts = append(whereParts, f)
		}
	}
	if lo.InFilters != nil {
		for _, f := range lo.InFilters {
			whereParts = append(whereParts, f)
		}
	}
	if lo.NotInFilters != nil {
		for _, f := range lo.NotInFilters {
			whereParts = append(whereParts, f)
		}
	}
	if lo.NullFilters != nil {
		for _, f := range lo.NullFilters {
			whereParts = append(whereParts, f)
		}
	}
	if lo.JSONFilters != nil {
		for _, f := range lo.JSONFilters {
			whereParts = append(whereParts, f)
		}
	}
	if lo.SearchQuery != nil {
		whereParts = append(whereParts, lo.SearchQuery)
	}
	if lo.OrFilters != nil {
		for _, f := range lo.OrFilters {
			whereParts = append(whereParts, f)
		}
	}
	return whereParts
}
