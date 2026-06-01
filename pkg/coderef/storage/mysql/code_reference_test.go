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

	"github.com/bucketeer-io/bucketeer/v2/pkg/coderef/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/coderef/storage"
	mysqlstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	coderefproto "github.com/bucketeer-io/bucketeer/v2/proto/coderef"
)

func TestNewCodeReferenceStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	s := NewCodeReferenceStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &codeReferenceStorage{}, s)
}

func TestCreateCodeReferenceMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*codeReferenceStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *codeReferenceStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *codeReferenceStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &codeReferenceStorage{qe: mock.NewMockClient(mockController)}
			p.setup(s)
			err := s.CreateCodeReference(context.Background(), &domain.CodeReference{
				CodeReference: coderefproto.CodeReference{},
			})
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateCodeReferenceMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*codeReferenceStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *codeReferenceStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: unexpected affected rows",
			setup: func(s *codeReferenceStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: storage.ErrCodeReferenceUnexpectedAffectedRows,
		},
		{
			desc: "success",
			setup: func(s *codeReferenceStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &codeReferenceStorage{qe: mock.NewMockClient(mockController)}
			p.setup(s)
			err := s.UpdateCodeReference(context.Background(), &domain.CodeReference{
				CodeReference: coderefproto.CodeReference{},
			})
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetCodeReferenceMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*codeReferenceStorage)
		expectedErr error
	}{
		{
			desc: "error: ErrNoRows",
			setup: func(s *codeReferenceStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysqlstorage.ErrNoRows)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: storage.ErrCodeReferenceNotFound,
		},
		{
			desc: "error: internal",
			setup: func(s *codeReferenceStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("internal"))
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: errors.New("internal"),
		},
		{
			desc: "success",
			setup: func(s *codeReferenceStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &codeReferenceStorage{qe: mock.NewMockClient(mockController)}
			p.setup(s)
			_, err := s.GetCodeReference(context.Background(), "cr-id")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListCodeReferencesMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(*codeReferenceStorage)
		params         storage.ListCodeReferencesParams
		expected       []*domain.CodeReference
		expectedCursor int
		expectedErr    error
	}{
		{
			desc:  "error: invalid order by",
			setup: nil,
			params: storage.ListCodeReferencesParams{
				OrderBy: coderefproto.ListCodeReferencesRequest_OrderBy(99),
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    storage.ErrInvalidOrderBy,
		},
		{
			desc:  "error: invalid cursor",
			setup: nil,
			params: storage.ListCodeReferencesParams{
				PageSize: 10,
				Cursor:   "invalid",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    storage.ErrInvalidCursor,
		},
		{
			desc:  "error: negative cursor",
			setup: nil,
			params: storage.ListCodeReferencesParams{
				PageSize: 10,
				Cursor:   "-1",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    storage.ErrInvalidCursor,
		},
		{
			desc: "error: query",
			setup: func(s *codeReferenceStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			params: storage.ListCodeReferencesParams{
				PageSize:      10,
				Cursor:        "0",
				EnvironmentID: "env-1",
				FeatureID:     "ftr-1",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "error: count",
			setup: func(s *codeReferenceStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("count error"))
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params: storage.ListCodeReferencesParams{
				PageSize:      10,
				Cursor:        "0",
				EnvironmentID: "env-1",
				FeatureID:     "ftr-1",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("count error"),
		},
		{
			desc: "success",
			setup: func(s *codeReferenceStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params: storage.ListCodeReferencesParams{
				PageSize:         10,
				Cursor:           "0",
				EnvironmentID:    "env-1",
				FeatureID:        "ftr-1",
				RepositoryName:   "repo",
				RepositoryOwner:  "owner",
				RepositoryType:   coderefproto.CodeReference_GITHUB,
				RepositoryBranch: "main",
				FileExtension:    "go",
				OrderBy:          coderefproto.ListCodeReferencesRequest_CREATED_AT,
				OrderDirection:   coderefproto.ListCodeReferencesRequest_DESC,
			},
			expected:       []*domain.CodeReference{},
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &codeReferenceStorage{qe: mock.NewMockClient(mockController)}
			if p.setup != nil {
				p.setup(s)
			}
			codeRefs, cursor, _, err := s.ListCodeReferences(context.Background(), p.params)
			assert.Equal(t, p.expected, codeRefs)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteCodeReferenceMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*codeReferenceStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *codeReferenceStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: not found",
			setup: func(s *codeReferenceStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: storage.ErrCodeReferenceNotFound,
		},
		{
			desc: "success",
			setup: func(s *codeReferenceStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &codeReferenceStorage{qe: mock.NewMockClient(mockController)}
			p.setup(s)
			err := s.DeleteCodeReference(context.Background(), "cr-id")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetCodeReferenceCountsByFeatureIDsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*codeReferenceStorage)
		featureIDs  []string
		expected    map[string]int64
		expectedErr error
	}{
		{
			desc:        "empty feature ids returns empty map",
			setup:       nil,
			featureIDs:  []string{},
			expected:    map[string]int64{},
			expectedErr: nil,
		},
		{
			desc: "error: query",
			setup: func(s *codeReferenceStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			featureIDs:  []string{"ftr-1"},
			expected:    nil,
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *codeReferenceStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			featureIDs:  []string{"ftr-1", "ftr-2"},
			expected:    map[string]int64{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &codeReferenceStorage{qe: mock.NewMockClient(mockController)}
			if p.setup != nil {
				p.setup(s)
			}
			counts, err := s.GetCodeReferenceCountsByFeatureIDs(
				context.Background(), "env-1", p.featureIDs,
			)
			assert.Equal(t, p.expected, counts)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
