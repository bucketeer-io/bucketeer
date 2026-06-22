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

	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
)

func TestNewAutoOpsRuleStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewAutoOpsRuleStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &autoOpsRuleStorage{}, storage)
}

func TestAutoOpsRuleStorage_CountOpsEventRate(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc          string
		setup         func(s *autoOpsRuleStorage)
		expectedCount int
		expectedErr   error
	}{
		{
			desc: "error",
			setup: func(s *autoOpsRuleStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedCount: 0,
			expectedErr:   errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *autoOpsRuleStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
					*dest[0].(*int) = 2
					return nil
				})
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedCount: 2,
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &autoOpsRuleStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(s)
			count, err := s.CountOpsEventRate(context.Background())
			assert.Equal(t, p.expectedCount, count)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
