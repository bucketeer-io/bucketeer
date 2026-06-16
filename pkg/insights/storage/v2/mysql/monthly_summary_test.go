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

package mysql

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	v2is "github.com/bucketeer-io/bucketeer/v2/pkg/insights/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
)

func TestNewMonthlySummaryStorageMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	s := NewMonthlySummaryStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &monthlySummaryStorage{}, s)
}

func TestUpsertMonthlySummaryBatchMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	records := []v2is.MonthlySummaryRecord{
		{Yearmonth: "202601", EnvironmentID: "env1", SourceID: "ANDROID", MAU: 10, Requests: 100},
	}

	patterns := []struct {
		desc        string
		setup       func(*monthlySummaryStorage)
		input       []v2is.MonthlySummaryRecord
		expectedErr error
	}{
		{
			desc:        "no records: noop",
			setup:       nil,
			input:       nil,
			expectedErr: nil,
		},
		{
			desc: "error",
			setup: func(s *monthlySummaryStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       records,
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *monthlySummaryStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input:       records,
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &monthlySummaryStorage{qe: mock.NewMockQueryExecer(mockController)}
			if p.setup != nil {
				p.setup(s)
			}
			err := s.UpsertMonthlySummaryBatch(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListMonthlySummariesMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(*monthlySummaryStorage)
		environmentIDs []string
		sourceIDs      []string
		expected       []v2is.ListMonthlySummaryResult
		expectedErr    error
	}{
		{
			desc:           "empty inputs: noop",
			setup:          nil,
			environmentIDs: nil,
			sourceIDs:      []string{"ANDROID"},
			expected:       nil,
			expectedErr:    nil,
		},
		{
			desc: "error: query",
			setup: func(s *monthlySummaryStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			environmentIDs: []string{"env1"},
			sourceIDs:      []string{"ANDROID"},
			expected:       nil,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *monthlySummaryStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(true)
				rows.EXPECT().Scan(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).DoAndReturn(func(dest ...any) error {
					*(dest[0].(*string)) = "env1"
					*(dest[1].(*string)) = "env-name"
					*(dest[2].(*string)) = "project-name"
					*(dest[3].(*string)) = "ANDROID"
					*(dest[4].(*string)) = "202601"
					*(dest[5].(*int64)) = int64(10)
					*(dest[6].(*int64)) = int64(100)
					return nil
				})
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			environmentIDs: []string{"env1"},
			sourceIDs:      []string{"ANDROID"},
			expected: []v2is.ListMonthlySummaryResult{
				{
					EnvironmentID:   "env1",
					EnvironmentName: "env-name",
					ProjectName:     "project-name",
					SourceID:        "ANDROID",
					Yearmonth:       "202601",
					MAU:             10,
					Requests:        100,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &monthlySummaryStorage{qe: mock.NewMockQueryExecer(mockController)}
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.ListMonthlySummaries(context.Background(), p.environmentIDs, p.sourceIDs)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}
