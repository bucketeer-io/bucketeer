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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	pgmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

// testAPIKey has the length of a generated key (32 random bytes, hex encoded).
const testAPIKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

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
		{
			desc: "Success: the digest and the prefix are stored next to the key",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), insertAPIKeyV2SQLQuery, gomock.Any(),
				).DoAndReturn(func(_ context.Context, _ string, args ...interface{}) (postgres.Result, error) {
					require.Len(t, args, 12)
					assert.Equal(t, testAPIKey, args[7], "api_key")
					assert.Equal(t, domain.HashAPIKey(testAPIKey), args[10], "api_key_hash")
					assert.NotEqual(t, testAPIKey, args[10], "the digest must differ from the key")
					assert.Equal(t, testAPIKey[:8], args[11], "api_key_prefix")
					return nil, nil
				})
			},
			apiKey:        &domain.APIKey{APIKey: &proto.APIKey{ApiKey: testAPIKey}},
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

// scanResult returns a row whose Scan fails with err (nil for a found row).
func scanResult(mockController *gomock.Controller, err error) *pgmock.MockRow {
	row := pgmock.NewMockRow(mockController)
	row.EXPECT().Scan(gomock.Any()).Return(err)
	return row
}

// The fallback lookups must match on the key column alone. Requiring the digest
// to be missing would make a row that is backfilled between the digest lookup
// and the fallback match neither query.
func TestAPIKeyPlaintextFallbackQueriesMatchOnKeyOnlyPostgres(t *testing.T) {
	t.Parallel()
	for _, query := range []string{
		selectEnvironmentAPIKeyPlaintextSQLQuery,
		selectAPIKeyV2ByAPIKeyPlaintextSQLQuery,
	} {
		_, where, found := strings.Cut(query, "WHERE")
		require.True(t, found)
		assert.NotContains(t, where, "api_key_hash")
		assert.Contains(t, where, "api_key = $1")
	}
}

func TestGetAPIKeyByAPIKeyPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	hash := domain.HashAPIKey(testAPIKey)
	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		expectedErr error
	}{
		{
			desc: "Success: matched on the digest of the key",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeySQLQuery, hash, "env-1",
				).Return(scanResult(mockController, nil))
			},
			expectedErr: nil,
		},
		{
			desc: "Success: row without a digest yet is matched on the key itself",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeySQLQuery, hash, "env-1",
				).Return(scanResult(mockController, postgres.ErrNoRows))
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeyPlaintextSQLQuery, testAPIKey, "env-1",
				).Return(scanResult(mockController, nil))
			},
			expectedErr: nil,
		},
		{
			desc: "ErrAPIKeyNotFound: no row from either lookup",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeySQLQuery, hash, "env-1",
				).Return(scanResult(mockController, postgres.ErrNoRows))
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeyPlaintextSQLQuery, testAPIKey, "env-1",
				).Return(scanResult(mockController, postgres.ErrNoRows))
			},
			expectedErr: v2as.ErrAPIKeyNotFound,
		},
		{
			desc: "Error: internal error on the digest lookup is returned without falling back",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeySQLQuery, hash, "env-1",
				).Return(scanResult(mockController, errInternal))
			},
			expectedErr: errInternal,
		},
		{
			desc: "Error: internal error on the fallback lookup",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeySQLQuery, hash, "env-1",
				).Return(scanResult(mockController, postgres.ErrNoRows))
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeyPlaintextSQLQuery, testAPIKey, "env-1",
				).Return(scanResult(mockController, errInternal))
			},
			expectedErr: errInternal,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetAPIKeyByAPIKey(context.Background(), testAPIKey, "env-1")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetEnvironmentAPIKeyPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	hash := domain.HashAPIKey(testAPIKey)
	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		expectedErr error
	}{
		{
			desc: "Success: matched on the digest of the key",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeySQLQuery, hash,
				).Return(scanResult(mockController, nil))
			},
			expectedErr: nil,
		},
		{
			desc: "Success: row without a digest yet is matched on the key itself",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeySQLQuery, hash,
				).Return(scanResult(mockController, postgres.ErrNoRows))
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeyPlaintextSQLQuery, testAPIKey,
				).Return(scanResult(mockController, nil))
			},
			expectedErr: nil,
		},
		{
			desc: "ErrAPIKeyNotFound: no row from either lookup",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeySQLQuery, hash,
				).Return(scanResult(mockController, postgres.ErrNoRows))
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeyPlaintextSQLQuery, testAPIKey,
				).Return(scanResult(mockController, postgres.ErrNoRows))
			},
			expectedErr: v2as.ErrAPIKeyNotFound,
		},
		{
			desc: "Error: internal error on the digest lookup is returned without falling back",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeySQLQuery, hash,
				).Return(scanResult(mockController, errInternal))
			},
			expectedErr: errInternal,
		},
		{
			desc: "Error: internal error on the fallback lookup",
			setup: func(s *accountStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeySQLQuery, hash,
				).Return(scanResult(mockController, postgres.ErrNoRows))
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeyPlaintextSQLQuery, testAPIKey,
				).Return(scanResult(mockController, errInternal))
			},
			expectedErr: errInternal,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetEnvironmentAPIKey(context.Background(), testAPIKey)
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
