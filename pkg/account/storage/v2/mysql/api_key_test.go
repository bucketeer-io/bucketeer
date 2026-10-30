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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

// testAPIKey has the length of a generated key (32 random bytes, hex encoded).
const testAPIKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func TestCreateAPIKeyMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: duplicate entry",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			expectedErr: v2as.ErrAPIKeyAlreadyExists,
		},
		{
			desc: "success",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: the digest and the prefix are stored next to the key",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), insertAPIKeyV2SQLQuery, gomock.Any(),
				).DoAndReturn(func(_ context.Context, _ string, args ...interface{}) (mysql.Result, error) {
					require.Len(t, args, 12)
					assert.Equal(t, testAPIKey, args[7], "api_key")
					assert.Equal(t, domain.HashAPIKey(testAPIKey), args[10], "api_key_hash")
					assert.NotEqual(t, testAPIKey, args[10], "the digest must differ from the key")
					assert.Equal(t, testAPIKey[:8], args[11], "api_key_prefix")
					return nil, nil
				})
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &accountStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			err := storage.CreateAPIKey(
				context.Background(),
				&domain.APIKey{APIKey: &proto.APIKey{ApiKey: testAPIKey}},
				"env-1",
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAPIKeyMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: unexpected affected rows",
			setup: func(s *accountStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: v2as.ErrAPIKeyUnexpectedAffectedRows,
		},
		{
			desc: "success",
			setup: func(s *accountStorage) {
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
			storage := &accountStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			err := storage.UpdateAPIKey(
				context.Background(),
				&domain.APIKey{APIKey: &proto.APIKey{}},
				"env-1",
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAPIKeyLastUsedAtMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc            string
		setup           func(*accountStorage)
		expectedUpdated bool
		expectedErr     error
	}{
		{
			desc: "error",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedUpdated: false,
			expectedErr:     errors.New("error"),
		},
		{
			desc: "success: updated",
			setup: func(s *accountStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedUpdated: true,
			expectedErr:     nil,
		},
		{
			desc: "success: not updated",
			setup: func(s *accountStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedUpdated: false,
			expectedErr:     nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &accountStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			updated, err := storage.UpdateAPIKeyLastUsedAt(
				context.Background(), "api-key-id", "env-1", 1234567890,
			)
			assert.Equal(t, p.expectedUpdated, updated)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAPIKeyMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		expectedErr error
	}{
		{
			desc: "error: ErrNoRows",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: v2as.ErrAPIKeyNotFound,
		},
		{
			desc: "error: internal",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("internal error"))
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			expectedErr: errors.New("internal error"),
		},
		{
			desc: "success",
			setup: func(s *accountStorage) {
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
			storage := &accountStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			_, err := storage.GetAPIKey(context.Background(), "api-key-id", "env-1")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

// scanResult returns a row whose Scan fails with err (nil for a found row).
func scanResult(mockController *gomock.Controller, err error) *mock.MockRow {
	row := mock.NewMockRow(mockController)
	row.EXPECT().Scan(gomock.Any()).Return(err)
	return row
}

// The fallback lookups must match on the key column alone. Requiring the digest
// to be missing would make a row that is backfilled between the digest lookup
// and the fallback match neither query.
func TestAPIKeyPlaintextFallbackQueriesMatchOnKeyOnlyMySQL(t *testing.T) {
	t.Parallel()
	for _, query := range []string{
		selectEnvironmentAPIKeyPlaintextSQLQuery,
		selectAPIKeyV2ByAPIKeyPlaintextSQLQuery,
	} {
		_, where, found := strings.Cut(query, "WHERE")
		require.True(t, found)
		assert.NotContains(t, where, "api_key_hash")
		assert.Contains(t, where, "api_key = ?")
	}
}

func TestGetAPIKeyByAPIKeyMySQL(t *testing.T) {
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
			desc: "success: matched on the digest of the key",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeySQLQuery, hash, "env-1",
				).Return(scanResult(mockController, nil))
			},
			expectedErr: nil,
		},
		{
			desc: "success: row without a digest yet is matched on the key itself",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeySQLQuery, hash, "env-1",
				).Return(scanResult(mockController, mysql.ErrNoRows))
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeyPlaintextSQLQuery, testAPIKey, "env-1",
				).Return(scanResult(mockController, nil))
			},
			expectedErr: nil,
		},
		{
			desc: "error: ErrNoRows from both lookups",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeySQLQuery, hash, "env-1",
				).Return(scanResult(mockController, mysql.ErrNoRows))
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeyPlaintextSQLQuery, testAPIKey, "env-1",
				).Return(scanResult(mockController, mysql.ErrNoRows))
			},
			expectedErr: v2as.ErrAPIKeyNotFound,
		},
		{
			desc: "error: internal error on the digest lookup is returned without falling back",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeySQLQuery, hash, "env-1",
				).Return(scanResult(mockController, errors.New("internal error")))
			},
			expectedErr: errors.New("internal error"),
		},
		{
			desc: "error: internal error on the fallback lookup",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeySQLQuery, hash, "env-1",
				).Return(scanResult(mockController, mysql.ErrNoRows))
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectAPIKeyV2ByAPIKeyPlaintextSQLQuery, testAPIKey, "env-1",
				).Return(scanResult(mockController, errors.New("internal error")))
			},
			expectedErr: errors.New("internal error"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &accountStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			_, err := storage.GetAPIKeyByAPIKey(
				context.Background(), testAPIKey, "env-1",
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetEnvironmentAPIKeyMySQL(t *testing.T) {
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
			desc: "success: matched on the digest of the key",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeySQLQuery, hash,
				).Return(scanResult(mockController, nil))
			},
			expectedErr: nil,
		},
		{
			desc: "success: row without a digest yet is matched on the key itself",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeySQLQuery, hash,
				).Return(scanResult(mockController, mysql.ErrNoRows))
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeyPlaintextSQLQuery, testAPIKey,
				).Return(scanResult(mockController, nil))
			},
			expectedErr: nil,
		},
		{
			desc: "error: ErrNoRows from both lookups",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeySQLQuery, hash,
				).Return(scanResult(mockController, mysql.ErrNoRows))
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeyPlaintextSQLQuery, testAPIKey,
				).Return(scanResult(mockController, mysql.ErrNoRows))
			},
			expectedErr: v2as.ErrAPIKeyNotFound,
		},
		{
			desc: "error: internal error on the digest lookup is returned without falling back",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeySQLQuery, hash,
				).Return(scanResult(mockController, errors.New("internal error")))
			},
			expectedErr: errors.New("internal error"),
		},
		{
			desc: "error: internal error on the fallback lookup",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeySQLQuery, hash,
				).Return(scanResult(mockController, mysql.ErrNoRows))
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), selectEnvironmentAPIKeyPlaintextSQLQuery, testAPIKey,
				).Return(scanResult(mockController, errors.New("internal error")))
			},
			expectedErr: errors.New("internal error"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &accountStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			_, err := storage.GetEnvironmentAPIKey(
				context.Background(), testAPIKey,
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAllEnvironmentAPIKeysMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		expectedErr error
	}{
		{
			desc: "error: query",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *accountStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &accountStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			keys, err := storage.ListAllEnvironmentAPIKeys(context.Background())
			if p.expectedErr != nil {
				assert.Nil(t, keys)
			} else {
				assert.NotNil(t, keys)
			}
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAPIKeysMySQL(t *testing.T) {
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
			desc: "error: query",
			setup: func(s *accountStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			params: v2as.ListAPIKeysParams{
				PageSize:       10,
				Cursor:         "0",
				OrganizationID: "org-1",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "error: count",
			setup: func(s *accountStorage) {
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
			params: v2as.ListAPIKeysParams{
				PageSize:       10,
				Cursor:         "0",
				OrganizationID: "org-1",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("count error"),
		},
		{
			desc: "success",
			setup: func(s *accountStorage) {
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
			storage := &accountStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
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
