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

import (
	"fmt"
	"math"
	"strings"
)

const placeHolder = "?"

type WherePart interface {
	SQLString() (sql string, args []interface{})
}

type Filter struct {
	Column   string
	Operator string
	Value    interface{}
}

func NewFilter(column, operator string, value interface{}) WherePart {
	return &Filter{
		Column:   column,
		Operator: operator,
		Value:    value,
	}
}

func (f *Filter) SQLString() (sql string, args []interface{}) {
	if f.Column == "" || f.Operator == "" {
		return "", nil
	}
	sql = fmt.Sprintf("%s %s %s", f.Column, f.Operator, placeHolder)
	args = append(args, f.Value)
	return
}

type InFilter struct {
	Column string
	Values []interface{}
}

func NewInFilter(column string, values []interface{}) WherePart {
	return &InFilter{
		Column: column,
		Values: values,
	}
}

func (f *InFilter) SQLString() (sql string, args []interface{}) {
	if f.Column == "" {
		return "", nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s IN (", f.Column))
	for i := range f.Values {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("?")
	}
	sb.WriteString(")")
	sql = sb.String()
	args = f.Values
	return
}

type NullFilter struct {
	Column string
	IsNull bool
}

func NewNullFilter(column string, isNull bool) WherePart {
	return &NullFilter{
		Column: column,
		IsNull: isNull,
	}
}

func (f *NullFilter) SQLString() (sql string, args []interface{}) {
	if f.Column == "" {
		return "", nil
	}
	var sb strings.Builder
	if f.IsNull {
		sb.WriteString(fmt.Sprintf("%s IS NULL", f.Column))
	} else {
		sb.WriteString(fmt.Sprintf("%s IS NOT NULL", f.Column))
	}
	sql = sb.String()
	return
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

func NewJSONFilter(column string, f JSONFilterFunc, values []interface{}) WherePart {
	return &JSONFilter{
		Column: column,
		Func:   f,
		Values: values,
	}
}

func (f *JSONFilter) SQLString() (sql string, args []interface{}) {
	if f.Column == "" {
		return "", nil
	}
	switch f.Func {
	case JSONContainsNumber:
		sql = fmt.Sprintf("JSON_CONTAINS(%s, ?)", f.Column)
		var sb strings.Builder
		sb.WriteString("[")
		for i, v := range f.Values {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteString("]")
		args = append(args, sb.String())
		return
	case JSONContainsString:
		sql = fmt.Sprintf("JSON_CONTAINS(%s, ?)", f.Column)
		var sb strings.Builder
		sb.WriteString("[")
		for i, v := range f.Values {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf(`"%s"`, v))
		}
		sb.WriteString("]")
		args = append(args, sb.String())
		return
	case JSONContainsJSON:
		sql = fmt.Sprintf("JSON_CONTAINS(%s, ?)", f.Column)
		var sb strings.Builder
		sb.WriteString("[")
		for i, v := range f.Values {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteString("]")
		args = append(args, sb.String())
		return
	case JSONLengthGreaterThan:
		if len(f.Values) == 0 {
			return "", nil
		}
		sql = fmt.Sprintf("JSON_LENGTH(%s) > %s", f.Column, f.Values[0])
		return
	case JSONLengthSmallerThan:
		if len(f.Values) == 0 {
			return "", nil
		}
		sql = fmt.Sprintf("JSON_LENGTH(%s) < %s", f.Column, f.Values[0])
		return
	default:
		return "", nil
	}
}

type SearchQuery struct {
	Columns []string
	Keyword string
}

func NewSearchQuery(columns []string, keyword string) WherePart {
	return &SearchQuery{
		Columns: columns,
		Keyword: keyword,
	}
}

func (q *SearchQuery) SQLString() (sql string, args []interface{}) {
	if len(q.Columns) == 0 {
		return "", nil
	}
	var sb strings.Builder
	sb.WriteString("(")
	for i, col := range q.Columns {
		if i != 0 {
			sb.WriteString(" OR ")
		}
		sb.WriteString(fmt.Sprintf("%s LIKE ?", col))
		args = append(args, "%"+q.Keyword+"%")
	}
	sb.WriteString(")")
	sql = sb.String()
	return
}

func ConstructWhereSQLString(wps []WherePart) (sql string, args []interface{}) {
	if len(wps) == 0 {
		return "", nil
	}
	var sb strings.Builder
	sb.WriteString("WHERE ")
	for i, wp := range wps {
		if i != 0 {
			sb.WriteString(" AND ")
		}
		wpSQL, wpArgs := wp.SQLString()
		sb.WriteString(wpSQL)
		args = append(args, wpArgs...)
	}
	sql = sb.String()
	return
}

type OrderDirection int

const (
	// default asc
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
	if len(orders) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("ORDER BY ")
	for i, o := range orders {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(o.Column)
		sb.WriteString(" ")
		sb.WriteString(o.Direction.String())
	}
	return sb.String()
}

const (
	QueryNoLimit  = 0
	QueryNoOffset = 0

	// Workaround for MySQL not support offset without limit
	// ref: https://dev.mysql.com/doc/refman/8.0/en/select.html
	queryLimitAllRows = math.MaxInt64
)

func ConstructLimitOffsetSQLString(limit, offset int) string {
	if limit == QueryNoLimit && offset == QueryNoOffset {
		return ""
	}
	if limit == QueryNoLimit && offset != QueryNoOffset {
		return fmt.Sprintf("LIMIT %d OFFSET %d", queryLimitAllRows, offset)
	}
	if limit != QueryNoLimit && offset == QueryNoOffset {
		return fmt.Sprintf("LIMIT %d", limit)
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}
