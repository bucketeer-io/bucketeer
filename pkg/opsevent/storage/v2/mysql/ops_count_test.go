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

	"github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/domain"
	v2os "github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

func TestNewOpsCountStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	s := NewOpsCountStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &opsCountStorage{}, s)
}

func TestUpsertOpsCountMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*opsCountStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *opsCountStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *opsCountStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &opsCountStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(s)
			err := s.UpsertOpsCount(context.Background(), "env-1", &domain.OpsCount{
				OpsCount: &proto.OpsCount{},
			})
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListOpsCountsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(*opsCountStorage)
		params         v2os.ListOpsCountsParams
		expected       []*proto.OpsCount
		expectedCursor int
		expectedErr    error
	}{
		{
			desc:  "error: invalid cursor",
			setup: nil,
			params: v2os.ListOpsCountsParams{
				PageSize: 10,
				Cursor:   "invalid",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    v2os.ErrInvalidCursor,
		},
		{
			desc: "error: query",
			setup: func(s *opsCountStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			params: v2os.ListOpsCountsParams{
				PageSize:      10,
				Cursor:        "0",
				EnvironmentID: "env-1",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *opsCountStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			params: v2os.ListOpsCountsParams{
				PageSize:       10,
				Cursor:         "5",
				EnvironmentID:  "env-1",
				FeatureIDs:     []string{"ftr-1", "ftr-2"},
				AutoOpsRuleIDs: []string{"rule-1"},
			},
			expected:       []*proto.OpsCount{},
			expectedCursor: 5,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &opsCountStorage{qe: mock.NewMockQueryExecer(mockController)}
			if p.setup != nil {
				p.setup(s)
			}
			opsCounts, cursor, err := s.ListOpsCounts(context.Background(), p.params)
			assert.Equal(t, p.expected, opsCounts)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
