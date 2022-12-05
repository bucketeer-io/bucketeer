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
	"testing"

	"cloud.google.com/go/bigtable"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestItem(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc          string
		row           *row
		expected      *ReadItem
		expectedError error
	}{
		{
			desc: "ErrColumnFamilyNotFound",
			row: &row{
				row: map[string][]bigtable.ReadItem{"columnFamily-2": {
					{
						Row:       "Row-1",
						Column:    "columnFamily:Column",
						Timestamp: 0,
						Value:     []byte("Value-1"),
					},
				}},
				logger: zap.NewNop(),
			},
			expected:      nil,
			expectedError: ErrColumnFamilyNotFound,
		},
		{
			desc: "Valid",
			row: &row{
				row: map[string][]bigtable.ReadItem{"columnFamily": {
					{
						Row:       "Row-1",
						Column:    "columnFamily:Column",
						Timestamp: 0,
						Value:     []byte("Value-1"),
					},
				}},
			},
			expected: &ReadItem{
				RowKey:    "Row-1",
				Column:    "columnFamily:Column",
				Timestamp: 0,
				Value:     []byte("Value-1"),
			},
			expectedError: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			items, err := p.row.ReadItem("columnFamily", "Column")
			assert.Equal(t, p.expected, items)
			assert.Equal(t, p.expectedError, err)
		})
	}
}

func TestItems(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc          string
		rows          *rows
		expected      []*ReadItem
		expectedError error
	}{
		{
			desc: "ErrColumnFamilyNotFound",
			rows: &rows{
				rows: []bigtable.Row{
					map[string][]bigtable.ReadItem{"columnFamily-2": {
						{
							Row:       "Row-1",
							Column:    "columnFamily:Column",
							Timestamp: 0,
							Value:     []byte("Value-1"),
						},
					}},
				},
				logger: zap.NewNop(),
			},
			expected:      nil,
			expectedError: ErrColumnFamilyNotFound,
		},
		{
			desc: "Valid",
			rows: &rows{
				rows: []bigtable.Row{
					map[string][]bigtable.ReadItem{"columnFamily": {
						{
							Row:       "Row-1",
							Column:    "columnFamily:Column",
							Timestamp: 0,
							Value:     []byte("Value-1"),
						},
					}},
					map[string][]bigtable.ReadItem{"columnFamily": {
						{
							Row:       "Row-2",
							Column:    "columnFamily:Column-1",
							Timestamp: 0,
							Value:     []byte("Value-1"),
						},
						{
							Row:       "Row-2",
							Column:    "columnFamily:Column",
							Timestamp: 0,
							Value:     []byte("Value-2"),
						},
					}},
					map[string][]bigtable.ReadItem{"columnFamily": {
						{
							Row:       "Row-3",
							Column:    "columnFamily:Column-1",
							Timestamp: 0,
							Value:     []byte("Value-1"),
						},
						{
							Row:       "Row-3",
							Column:    "columnFamily:Column-2",
							Timestamp: 0,
							Value:     []byte("Value-2"),
						},
						{
							Row:       "Row-3",
							Column:    "columnFamily:Column",
							Timestamp: 0,
							Value:     []byte("Value-4"),
						},
						{
							Row:       "Row-3",
							Column:    "columnFamily:Column",
							Timestamp: 0,
							Value:     []byte("Value-3"),
						},
					}},
				},
				logger: zap.NewNop(),
			},
			expected: []*ReadItem{
				{
					RowKey:    "Row-1",
					Column:    "columnFamily:Column",
					Timestamp: 0,
					Value:     []byte("Value-1"),
				},
				{
					RowKey:    "Row-2",
					Column:    "columnFamily:Column",
					Timestamp: 0,
					Value:     []byte("Value-2"),
				},
				{
					RowKey:    "Row-3",
					Column:    "columnFamily:Column",
					Timestamp: 0,
					Value:     []byte("Value-4"),
				},
				{
					RowKey:    "Row-3",
					Column:    "columnFamily:Column",
					Timestamp: 0,
					Value:     []byte("Value-3"),
				},
			},
			expectedError: nil,
		},
	}
	for _, p := range patterns {
		items, err := p.rows.ReadItems("columnFamily", "Column")
		assert.Equal(t, p.expected, items)
		assert.Equal(t, p.expectedError, err)
	}
}
