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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package mysql

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

const placeHolder = "?"

// searchKeywordMaxLength is the maximum number of characters of a search
// keyword that take part in the LIKE pattern; longer keywords are truncated.
const searchKeywordMaxLength = 255

// likeEscaper escapes the LIKE wildcards and the escape character itself so a
// keyword is matched literally. The pattern must be used with ESCAPE '\\'.
var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

type Operator int

const (
	// Operation to find the field is equal to the specified value.
	OperatorEqual = iota + 1
	// Operation to find the field isn't equal to the specified value.
	OperatorNotEqual
	// Operation to find ones that contain any one of the multiple values.
	OperatorIn
	// Operation to find ones that do not contain any of the specified multiple values.
	OperatorNotIn
	// Operation to find ones the field is greater than the specified value.
	OperatorGreaterThan
	// Operation to find ones the field is greater or equal than the specified value.
	OperatorGreaterThanOrEqual
	// Operation to find ones the field is less than the specified value.
	OperatorLessThan
	// Operation to find ones the field is less or equal than the specified value.
	OperatorLessThanOrEqual
	// Operation to find ones that have a specified value in its array.
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
	OperatorContains:           "MEMBER OF",
}

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

type FilterV2 struct {
	Column   string
	Operator Operator
	Value    interface{}
}

func (f *FilterV2) SQLString() (sql string, args []interface{}) {
	if f.Column == "" || f.Operator < OperatorEqual || f.Operator > OperatorContains {
		return "", nil
	}
	sql = fmt.Sprintf("%s %s %s", f.Column, operatorMap[f.Operator], placeHolder)
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
	var sb strings.Builder
	if f.Column == "" || len(f.Values) == 0 {
		return "", nil
	} else {
		sb.WriteString(fmt.Sprintf(" %s IN (", f.Column))
	}
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
	var sb strings.Builder
	if f.Column == "" {
		return "", nil
	} else {
		sb.WriteString(" ")
	}
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
		strs := make([]string, 0, len(f.Values))
		for _, v := range f.Values {
			s, ok := v.(string)
			if !ok {
				return "", nil
			}
			strs = append(strs, s)
		}
		payload, err := json.Marshal(strs)
		if err != nil {
			return "", nil
		}
		sql = fmt.Sprintf("JSON_CONTAINS(%s, ?)", f.Column)
		args = append(args, string(payload))
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
	case JSONLengthGreaterThan, JSONLengthSmallerThan:
		if len(f.Values) == 0 {
			return "", nil
		}
		op := ">"
		if f.Func == JSONLengthSmallerThan {
			op = "<"
		}
		sql = fmt.Sprintf("JSON_LENGTH(%s) %s ?", f.Column, op)
		args = append(args, f.Values[0])
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
	var sb strings.Builder
	if len(q.Columns) == 0 {
		return "", nil
	} else {
		sb.WriteString(" (")
	}
	pattern := searchKeywordPattern(q.Keyword)
	for i, col := range q.Columns {
		if i != 0 {
			sb.WriteString(" OR ")
		}
		sb.WriteString(fmt.Sprintf(`%s LIKE ? ESCAPE '\\'`, col))
		args = append(args, pattern)
	}
	sb.WriteString(")")
	sql = sb.String()
	return
}

// searchKeywordPattern builds the LIKE pattern that matches the keyword as a
// literal substring: the keyword is truncated to searchKeywordMaxLength
// characters and its wildcards are escaped before it is wrapped in %.
func searchKeywordPattern(keyword string) string {
	if runes := []rune(keyword); len(runes) > searchKeywordMaxLength {
		keyword = string(runes[:searchKeywordMaxLength])
	}
	return "%" + likeEscaper.Replace(keyword) + "%"
}

type ExistsFilter struct {
	Subquery  string
	Values    []interface{}
	NotExists bool
}

func NewExistsFilter(subquery string, notExists bool, values ...interface{}) WherePart {
	return &ExistsFilter{
		Subquery:  subquery,
		Values:    values,
		NotExists: notExists,
	}
}

func (f *ExistsFilter) SQLString() (sql string, args []interface{}) {
	if f.Subquery == "" {
		return "", nil
	}
	if f.NotExists {
		sql = fmt.Sprintf("NOT EXISTS (%s)", f.Subquery)
	} else {
		sql = fmt.Sprintf("EXISTS (%s)", f.Subquery)
	}
	args = f.Values
	return
}

type OrFilter struct {
	Queries []WherePart
}

func NewOrFilter(queries []WherePart) WherePart {
	return &OrFilter{
		Queries: queries,
	}
}

func (f *OrFilter) SQLString() (sql string, args []interface{}) {
	return joinWhereParts(f.Queries, " OR ")
}

type AndFilter struct {
	Queries []WherePart
}

func NewAndFilter(queries []WherePart) WherePart {
	return &AndFilter{
		Queries: queries,
	}
}

func (f *AndFilter) SQLString() (sql string, args []interface{}) {
	return joinWhereParts(f.Queries, " AND ")
}

// joinWhereParts wraps the non-empty fragments of the parts in parentheses,
// joined with the separator. Parts that produce no SQL (for example a filter
// without a column or without values) are skipped, so they never leave a
// dangling operator behind. It returns no SQL when every part is empty.
func joinWhereParts(parts []WherePart, separator string) (sql string, args []interface{}) {
	var sb strings.Builder
	for _, part := range parts {
		partSQL, partArgs := part.SQLString()
		if partSQL == "" {
			continue
		}
		if sb.Len() == 0 {
			sb.WriteString("(")
		} else {
			sb.WriteString(separator)
		}
		sb.WriteString(partSQL)
		args = append(args, partArgs...)
	}
	if sb.Len() == 0 {
		return "", nil
	}
	sb.WriteString(")")
	return sb.String(), args
}

func ConstructWhereSQLString(wps []WherePart) (sql string, args []interface{}) {
	var sb strings.Builder
	for _, wp := range wps {
		wpSQL, wpArgs := wp.SQLString()
		if wpSQL == "" {
			// A part without SQL must not leave a dangling WHERE or AND behind.
			continue
		}
		if sb.Len() == 0 {
			sb.WriteString(" WHERE ")
		} else {
			sb.WriteString(" AND ")
		}
		sb.WriteString(wpSQL)
		args = append(args, wpArgs...)
	}
	if sb.Len() == 0 {
		return "", nil
	}
	sql = sb.String() + " "
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
	var sb strings.Builder
	if len(orders) == 0 {
		return ""
	} else {
		sb.WriteString(" ORDER BY ")
	}
	for i, o := range orders {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(o.Column)
		sb.WriteString(" ")
		sb.WriteString(o.Direction.String())
	}
	constructStr := sb.String() + " "
	return constructStr
}

func ConstructQueryAndWhereArgs(baseQuery string, options *ListOptions) (query string, whereArgs []interface{}) {
	if options != nil {
		var whereQuery string
		whereParts := options.CreateWhereParts()
		whereQuery, whereArgs = ConstructWhereSQLString(whereParts)
		orderByQuery := ConstructOrderBySQLString(options.Orders)
		limitOffsetQuery := ConstructLimitOffsetSQLString(options.Limit, options.Offset)
		query = baseQuery + whereQuery + orderByQuery + limitOffsetQuery
	} else {
		query = baseQuery
		whereArgs = []interface{}{}
	}
	return
}

// ConstructCountQuery builds a count query with optional filtering.
// It takes a base count SQL query and ListOptions, and returns the complete query and its arguments.
// The base SQL should be a valid COUNT query (e.g., "SELECT COUNT(1) FROM table").
// Additional WHERE clauses from ListOptions will be appended to the base query.
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

type Orders struct {
	Orders []*Order
}

func (o *Orders) SQLString() (sql string, args []interface{}) {
	var sb strings.Builder
	if len(o.Orders) == 0 {
		return "", nil
	} else {
		sb.WriteString("ORDER BY ")
	}
	for i, o := range o.Orders {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(o.Column)
		sb.WriteString(" ")
		sb.WriteString(o.Direction.String())
	}
	return sb.String(), nil
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
		return fmt.Sprintf(" LIMIT %d OFFSET %d", queryLimitAllRows, offset)
	}
	if limit != QueryNoLimit && offset == QueryNoOffset {
		return fmt.Sprintf(" LIMIT %d", limit)
	}
	return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
}

type ListOptions struct {
	Limit         int
	Filters       []*FilterV2
	InFilters     []*InFilter
	NullFilters   []*NullFilter
	JSONFilters   []*JSONFilter
	ExistsFilters []*ExistsFilter
	SearchQuery   *SearchQuery
	OrFilters     []*OrFilter
	Orders        []*Order
	Offset        int
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
	if lo.ExistsFilters != nil {
		for _, f := range lo.ExistsFilters {
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
