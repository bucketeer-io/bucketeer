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

	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/storage/v2"
	mysqlstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

func TestNewAutoOpsRuleStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	s := NewAutoOpsRuleStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &autoOpsRuleStorage{}, s)
}

func TestCreateAutoOpsRuleMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*autoOpsRuleStorage)
		expectedErr error
	}{
		{
			desc: "error: duplicate entry",
			setup: func(s *autoOpsRuleStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysqlstorage.ErrDuplicateEntry)
			},
			expectedErr: v2as.ErrAutoOpsRuleAlreadyExists,
		},
		{
			desc: "error: internal",
			setup: func(s *autoOpsRuleStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *autoOpsRuleStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &autoOpsRuleStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(s)
			err := s.CreateAutoOpsRule(context.Background(),
				&domain.AutoOpsRule{AutoOpsRule: &proto.AutoOpsRule{}}, "ns0")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAutoOpsRuleMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*autoOpsRuleStorage)
		expectedErr error
	}{
		{
			desc: "error: internal",
			setup: func(s *autoOpsRuleStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: unexpected affected rows",
			setup: func(s *autoOpsRuleStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: v2as.ErrAutoOpsRuleUnexpectedAffectedRows,
		},
		{
			desc: "success",
			setup: func(s *autoOpsRuleStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &autoOpsRuleStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(s)
			err := s.UpdateAutoOpsRule(context.Background(),
				&domain.AutoOpsRule{AutoOpsRule: &proto.AutoOpsRule{}}, "ns0")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAutoOpsRuleMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*autoOpsRuleStorage)
		expectedErr error
	}{
		{
			desc: "error: not found",
			setup: func(s *autoOpsRuleStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysqlstorage.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: v2as.ErrAutoOpsRuleNotFound,
		},
		{
			desc: "error: internal",
			setup: func(s *autoOpsRuleStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("internal"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: errors.New("internal"),
		},
		{
			desc: "success",
			setup: func(s *autoOpsRuleStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &autoOpsRuleStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(s)
			_, err := s.GetAutoOpsRule(context.Background(), "id", "ns0")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAutoOpsRulesMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(*autoOpsRuleStorage)
		params         v2as.ListAutoOpsRulesParams
		expected       []*proto.AutoOpsRule
		expectedCursor int
		expectedErr    error
	}{
		{
			desc:           "error: invalid cursor",
			setup:          nil,
			params:         v2as.ListAutoOpsRulesParams{Cursor: "invalid"},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    v2as.ErrInvalidCursor,
		},
		{
			desc:           "error: negative cursor",
			setup:          nil,
			params:         v2as.ListAutoOpsRulesParams{Cursor: "-1"},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    v2as.ErrInvalidCursor,
		},
		{
			desc: "error: query",
			setup: func(s *autoOpsRuleStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			params:         v2as.ListAutoOpsRulesParams{EnvironmentID: "ns0", Cursor: "0"},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *autoOpsRuleStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			params: v2as.ListAutoOpsRulesParams{
				EnvironmentID: "ns0",
				FeatureIDs:    []string{"f1", "f2"},
				PageSize:      10,
				Cursor:        "0",
			},
			expected:       []*proto.AutoOpsRule{},
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &autoOpsRuleStorage{qe: mock.NewMockQueryExecer(mockController)}
			if p.setup != nil {
				p.setup(s)
			}
			rules, cursor, err := s.ListAutoOpsRules(context.Background(), p.params)
			assert.Equal(t, p.expected, rules)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
