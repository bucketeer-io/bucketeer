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
	"go.uber.org/zap"
)

// A ReadItem is returned by ReadItems. A ReadItem contains data from a specific row and column.
type ReadItem struct {
	RowKey, Column string
	Timestamp      int64
	Value          []byte
}

// A WriteItem contains the key and value to write an item to a specific row
type WriteItem struct {
	Key   string
	Value []byte
}

type Row interface {
	ReadItem(columnFamily, column string) (*ReadItem, error)
}

type row struct {
	row    bigtable.Row
	logger *zap.Logger
}

type Rows interface {
	ReadItems(columnFamily, column string) ([]*ReadItem, error)
}

type rows struct {
	rows   []bigtable.Row
	logger *zap.Logger
}

func (r *row) ReadItem(columnFamily, column string) (readItem *ReadItem, err error) {
	defer record()(operationReadItem, &err)
	items, err := getColumnItems(r.row, columnFamily, column)
	if err != nil {
		r.logger.Error("Failed to read item by column",
			zap.Error(err),
			zap.String("columnFamily", columnFamily),
			zap.String("column", column),
		)
		return nil, err
	}
	return items[0], nil
}

func (r *rows) ReadItems(columnFamily, column string) (readItems []*ReadItem, err error) {
	defer record()(operationReadItems, &err)
	for _, row := range r.rows {
		var items []*ReadItem
		items, err = getColumnItems(row, columnFamily, column)
		if err != nil {
			r.logger.Error("Failed to read items by column",
				zap.Error(err),
				zap.String("columnFamily", columnFamily),
				zap.String("column", column),
			)
			return nil, err
		}
		readItems = append(readItems, items...)
	}
	return readItems, nil
}

func getColumnItems(row bigtable.Row, columnFamily, column string) ([]*ReadItem, error) {
	items, ok := row[columnFamily]
	if !ok {
		return nil, ErrColumnFamilyNotFound
	}
	var readItems []*ReadItem
	col := fmt.Sprintf("%s:%s", columnFamily, column)
	for _, item := range items {
		if item.Column == col {
			i := &ReadItem{
				RowKey:    item.Row,
				Column:    item.Column,
				Timestamp: item.Timestamp.Time().Unix(),
				Value:     item.Value,
			}
			readItems = append(readItems, i)
		}
	}
	if len(readItems) == 0 {
		return nil, ErrColumnNotFound
	}
	return readItems, nil
}
