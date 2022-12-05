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

package bigtable

import (
	"fmt"

	"cloud.google.com/go/bigtable"
)

func NewKey(environmentNamespace, key string) string {
	if environmentNamespace == "" {
		return fmt.Sprintf("default#%s", key)
	}
	return fmt.Sprintf("%s#%s", environmentNamespace, key)
}

// All fields are required
type ReadRequest struct {
	TableName    string
	ColumnFamily string
	RowSet       RowSet
	RowFilters   []RowFilter // Optional
}

// All fields are required
type ReadRowRequest struct {
	TableName    string
	ColumnFamily string
	RowKey       string
}

// All fields are required
type WriteRequest struct {
	TableName    string
	ColumnFamily string
	ColumnName   string
	Items        []*WriteItem
}

type RowFilter interface {
	get() bigtable.Filter
}

// LatestNFilter returns a filter that matches the most recent N cells in each column.
type LatestNFilter int

func (r LatestNFilter) get() bigtable.Filter {
	return bigtable.LatestNFilter(int(r))
}

// ColumnFilter returns a filter that matches cells whose column name
type ColumnFilter string

func (r ColumnFilter) get() bigtable.Filter {
	return bigtable.ColumnFilter(string(r))
}

// RowSet is a set of rows to be read.
// It is satisfied by RowKey, RowList, RowPrefix, and RowPrefixRange.
type RowSet interface {
	get() bigtable.RowSet
}

// The row key.
type RowKey string

func (r RowKey) get() bigtable.RowSet {
	// SingleRow returns a RowSet for reading a single row.
	return bigtable.SingleRow(string(r))
}

// RowList is a sequence of row keys.
type RowList []string

func (r RowList) get() bigtable.RowSet {
	keys := make(bigtable.RowList, 0, len(r))
	return append(keys, r...)
}

// PrefixRange returns a RowRange consisting of all keys starting with the prefix.
type RowPrefix string

func (r RowPrefix) get() bigtable.RowSet {
	return bigtable.PrefixRange(string(r))
}

// A RowRange is a half-open interval (start, limit) encompassing
// all the rows with keys at least as large as Start, and less than Limit.
type RowPrefixRange struct {
	Start string
	Limit string
}

func (r RowPrefixRange) get() bigtable.RowSet {
	return bigtable.NewRange(r.Start, r.Limit)
}
