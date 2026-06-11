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

package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	pgmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

func TestCreateAPIKeyPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc          string
		setup         func(*accountStorage)
		apiKey        *domain.APIKey
		environmentID string
		expectedErr   error
	}{
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			apiKey:        &domain.APIKey{APIKey: &proto.APIKey{}},
			environmentID: "env-1",
			expectedErr:   errInternal,
		},
		{
			desc: "ErrAPIKeyAlreadyExists",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, postgres.ErrDuplicateEntry)
			},
			apiKey:        &domain.APIKey{APIKey: &proto.APIKey{}},
			environmentID: "env-1",
			expectedErr:   v2as.ErrAPIKeyAlreadyExists,
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			apiKey:        &domain.APIKey{APIKey: &proto.APIKey{}},
			environmentID: "env-1",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateAPIKey(context.Background(), p.apiKey, p.environmentID)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAPIKeyPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc          string
		setup         func(*accountStorage)
		apiKey        *domain.APIKey
		environmentID string
		expectedErr   error
	}{
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			apiKey:        &domain.APIKey{APIKey: &proto.APIKey{}},
			environmentID: "env-1",
			expectedErr:   errInternal,
		},
		{
			desc: "ErrUnexpectedAffectedRows",
			setup: func(s *accountStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			apiKey:        &domain.APIKey{APIKey: &proto.APIKey{}},
			environmentID: "env-1",
			expectedErr:   v2as.ErrAPIKeyUnexpectedAffectedRows,
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			apiKey:        &domain.APIKey{APIKey: &proto.APIKey{}},
			environmentID: "env-1",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpdateAPIKey(context.Background(), p.apiKey, p.environmentID)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAPIKeyLastUsedAtPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc            string
		setup           func(*accountStorage)
		id              string
		environmentID   string
		lastUsedAt      int64
		expectedUpdated bool
		expectedErr     error
	}{
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			id:              "api-key-1",
			environmentID:   "env-1",
			lastUsedAt:      1000,
			expectedUpdated: false,
			expectedErr:     errInternal,
		},
		{
			desc: "Success_Updated",
			setup: func(s *accountStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			id:              "api-key-1",
			environmentID:   "env-1",
			lastUsedAt:      1000,
			expectedUpdated: true,
			expectedErr:     nil,
		},
		{
			desc: "Success_NotUpdated",
			setup: func(s *accountStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			id:              "api-key-1",
			environmentID:   "env-1",
			lastUsedAt:      1000,
			expectedUpdated: false,
			expectedErr:     nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			updated, err := storage.UpdateAPIKeyLastUsedAt(
				context.Background(), p.id, p.environmentID, p.lastUsedAt,
			)
			assert.Equal(t, p.expectedUpdated, updated)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAPIKeyPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc          string
		setup         func(*accountStorage)
		id            string
		environmentID string
		expectedErr   error
	}{
		{
			desc: "ErrAPIKeyNotFound",
			setup: func(s *accountStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(postgres.ErrNoRows)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:            "api-key-1",
			environmentID: "env-1",
			expectedErr:   v2as.ErrAPIKeyNotFound,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errInternal)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:            "api-key-1",
			environmentID: "env-1",
			expectedErr:   errInternal,
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:            "api-key-1",
			environmentID: "env-1",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetAPIKey(context.Background(), p.id, p.environmentID)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAPIKeyByAPIKeyPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc          string
		setup         func(*accountStorage)
		apiKey        string
		environmentID string
		expectedErr   error
	}{
		{
			desc: "ErrAPIKeyNotFound",
			setup: func(s *accountStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(postgres.ErrNoRows)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			apiKey:        "some-api-key",
			environmentID: "env-1",
			expectedErr:   v2as.ErrAPIKeyNotFound,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errInternal)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			apiKey:        "some-api-key",
			environmentID: "env-1",
			expectedErr:   errInternal,
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			apiKey:        "some-api-key",
			environmentID: "env-1",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetAPIKeyByAPIKey(context.Background(), p.apiKey, p.environmentID)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetEnvironmentAPIKeyPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		apiKey      string
		expectedErr error
	}{
		{
			desc: "ErrAPIKeyNotFound",
			setup: func(s *accountStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(postgres.ErrNoRows)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			apiKey:      "some-api-key",
			expectedErr: v2as.ErrAPIKeyNotFound,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errInternal)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			apiKey:      "some-api-key",
			expectedErr: errInternal,
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			apiKey:      "some-api-key",
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetEnvironmentAPIKey(context.Background(), p.apiKey)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAllEnvironmentAPIKeysPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		expected    []*domain.EnvironmentAPIKey
		expectedErr error
	}{
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			expected:    nil,
			expectedErr: errInternal,
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expected:    []*domain.EnvironmentAPIKey{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			keys, err := storage.ListAllEnvironmentAPIKeys(context.Background())
			assert.Equal(t, p.expected, keys)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAPIKeysPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc           string
		setup          func(*accountStorage)
		params         v2as.ListAPIKeysParams
		expected       []*proto.APIKey
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "QueryError",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			params: v2as.ListAPIKeysParams{
				PageSize:       10,
				Cursor:         "0",
				OrganizationID: "org-1",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errInternal,
		},
		{
			desc: "CountError",
			setup: func(s *accountStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errInternal)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params: v2as.ListAPIKeysParams{
				PageSize:       10,
				Cursor:         "0",
				OrganizationID: "org-1",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errInternal,
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params: v2as.ListAPIKeysParams{
				PageSize:       10,
				Cursor:         "0",
				OrganizationID: "org-1",
			},
			expected:       []*proto.APIKey{},
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			apiKeys, cursor, _, err := storage.ListAPIKeys(
				context.Background(),
				p.params,
			)
			assert.Equal(t, p.expected, apiKeys)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
