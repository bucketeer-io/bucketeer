// Copyright 2024 The Bucketeer Authors.
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

package v2

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/account"
)

func TestNewAccountStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewAccountStorage(mock.NewMockClient(mockController))
	assert.IsType(t, &accountStorage{}, storage)
}

func TestCreateAccountV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		input       *domain.AccountV2
		expectedErr error
	}{
		{
			desc: "ErrAccountAlreadyExists",
			setup: func(s *accountStorage) {
				s.client.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: &domain.AccountV2{
				AccountV2: &proto.AccountV2{Email: "test@example.com"},
			},
			expectedErr: ErrAccountAlreadyExists,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.client.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input: &domain.AccountV2{
				AccountV2: &proto.AccountV2{Email: "test@example.com"},
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				s.client.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: &domain.AccountV2{
				AccountV2: &proto.AccountV2{Email: "test@example.com"},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateAccountV2(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCreateAccountMockV2(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		setup       func(s *accountStorage, testMock sqlmock.Sqlmock)
		input       *domain.AccountV2
		expectedErr error
	}{
		{
			desc: "ErrAccountAlreadyExists",
			setup: func(s *accountStorage, testMock sqlmock.Sqlmock) {
				testMock.ExpectExec("").WillReturnError(mysql.ErrDuplicateEntry)
			},
			input: &domain.AccountV2{
				AccountV2: &proto.AccountV2{Email: "test@example.com"},
			},
			expectedErr: ErrAccountAlreadyExists,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage, testMock sqlmock.Sqlmock) {
				testMock.ExpectExec("").WillReturnError(errors.New("error"))
			},
			input: &domain.AccountV2{
				AccountV2: &proto.AccountV2{Email: "test@example.com"},
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *accountStorage, testMock sqlmock.Sqlmock) {
				query := `INSERT INTO account_v2 ( email, name, avatar_image_url, organization_id, organization_role, environment_roles, disabled, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
				testMock.
					ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(
						"Mike",
						"a@email.com",
						"url",
						"org-0",
						1,
						[]byte(`[{"environment_id":"env-0","role":1}]`),
						false,
						5,
						6,
					).
					WillReturnResult(sqlmock.NewResult(1, 9))
			},
			input: &domain.AccountV2{
				AccountV2: &proto.AccountV2{
					Email:            "a@email.com",
					Name:             "Mike",
					AvatarImageUrl:   "url",
					OrganizationId:   "org-0",
					OrganizationRole: 1,
					EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
						{
							EnvironmentId: "env-0",
							Role:          1,
						},
					},
					Disabled:      false,
					CreatedAt:     5,
					UpdatedAt:     6,
					SearchFilters: nil,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage, testMock := newAccountStorageWithMock2(t)
			defer storage.client.Close()

			if p.setup != nil {
				p.setup(storage, testMock)
			}
			err := storage.CreateAccountV2(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAccountV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		input       *domain.AccountV2
		expectedErr error
	}{
		{
			desc: "ErrAccountUnexpectedAffectedRows",
			setup: func(s *accountStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.client.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.AccountV2{
				AccountV2: &proto.AccountV2{Email: "test@example.com"},
			},
			expectedErr: ErrAccountUnexpectedAffectedRows,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.client.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input: &domain.AccountV2{
				AccountV2: &proto.AccountV2{Email: "test@example.com"},
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.client.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.AccountV2{
				AccountV2: &proto.AccountV2{Email: "test@example.com"},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpdateAccountV2(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteMockV2(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		setup       func(s *accountStorage, testMock sqlmock.Sqlmock)
		input       *domain.AccountV2
		expectedErr error
	}{
		{
			desc: "ErrAccountAlreadyExists",
			setup: func(s *accountStorage, testMock sqlmock.Sqlmock) {
				testMock.ExpectExec("").WillReturnError(mysql.ErrDuplicateEntry)
			},
			input: &domain.AccountV2{
				AccountV2: &proto.AccountV2{Email: "test@example.com"},
			},
			expectedErr: mysql.ErrDuplicateEntry,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage, testMock sqlmock.Sqlmock) {
				testMock.ExpectExec("").WillReturnError(errors.New("error"))
			},
			input: &domain.AccountV2{
				AccountV2: &proto.AccountV2{Email: "test@example.com"},
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "ErrAccountUnexpectedAffectedRows",
			setup: func(s *accountStorage, testMock sqlmock.Sqlmock) {
				query := `DELETE FROM account_v2 WHERE email = ? AND organization_id = ?`
				testMock.
					ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(
						"a@email.com",
						"org-0",
					).
					WillReturnResult(sqlmock.NewResult(1, -1))
			},
			input: &domain.AccountV2{
				AccountV2: &proto.AccountV2{
					Email:            "a@email.com",
					Name:             "Mike",
					AvatarImageUrl:   "url",
					OrganizationId:   "org-0",
					OrganizationRole: 1,
					EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
						{
							EnvironmentId: "env-0",
							Role:          1,
						},
					},
					Disabled:      false,
					CreatedAt:     5,
					UpdatedAt:     6,
					SearchFilters: nil,
				},
			},
			expectedErr: ErrAccountUnexpectedAffectedRows,
		},
		{
			desc: "Success",
			setup: func(s *accountStorage, testMock sqlmock.Sqlmock) {
				query := `DELETE FROM account_v2 WHERE email = ? AND organization_id = ?`
				testMock.
					ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(
						"a@email.com",
						"org-0",
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: &domain.AccountV2{
				AccountV2: &proto.AccountV2{
					Email:            "a@email.com",
					Name:             "Mike",
					AvatarImageUrl:   "url",
					OrganizationId:   "org-0",
					OrganizationRole: 1,
					EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
						{
							EnvironmentId: "env-0",
							Role:          1,
						},
					},
					Disabled:      false,
					CreatedAt:     5,
					UpdatedAt:     6,
					SearchFilters: nil,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage, testMock := newAccountStorageWithMock2(t)
			defer storage.client.Close()

			if p.setup != nil {
				p.setup(storage, testMock)
			}
			err := storage.DeleteAccountV2(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAccountV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc           string
		setup          func(*accountStorage)
		email          string
		organizationID string
		expectedErr    error
	}{
		{
			desc: "ErrAccountNotFound",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.client.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			email:          "test@example.com",
			organizationID: "org-0",
			expectedErr:    ErrAccountNotFound,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.client.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			email:          "test@example.com",
			organizationID: "org-0",
			expectedErr:    errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.client.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			email:          "test@example.com",
			organizationID: "org-0",
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetAccountV2(context.Background(), p.email, p.organizationID)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAccountV2ByEnvironmentID(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc          string
		setup         func(*accountStorage)
		email         string
		environmentID string
		expectedErr   error
	}{
		{
			desc: "ErrAccountNotFound",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.client.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			email:         "test@example.com",
			environmentID: "env-0",
			expectedErr:   ErrAccountNotFound,
		},
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.client.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			email:         "test@example.com",
			environmentID: "env-0",
			expectedErr:   errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.client.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			email:         "test@example.com",
			environmentID: "env-0",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetAccountV2ByEnvironmentID(context.Background(), p.email, p.environmentID)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAccountsWithOrganization(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*accountStorage)
		email       string
		expectedErr error
	}{
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.client.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			email:       "test@example.com",
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.client.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			email:       "test@example.com",
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetAccountsWithOrganization(context.Background(), p.email)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAccountsV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc           string
		setup          func(*accountStorage)
		whereParts     []mysql.WherePart
		orders         []*mysql.Order
		limit          int
		offset         int
		expected       []*proto.AccountV2
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "Error",
			setup: func(s *accountStorage) {
				s.client.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			whereParts:     nil,
			orders:         nil,
			limit:          0,
			offset:         0,
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *accountStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)

				var nextCount = 0
				rows.EXPECT().Next().DoAndReturn(
					func() bool {
						nextCount++
						return nextCount < 2
					}).Times(2)
				rows.EXPECT().Err().Return(nil)
				rows.EXPECT().Scan(gomock.Any()).Return(nil)
				s.client.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(),
					//					gomock.Regex(`SELECT email, name, avatar_image_url, organization_id, organization_role, environment_roles, disabled, created_at, updated_at, search_filters FROM account_v2 WHERE num >= ? ORDER BY id ASC LIMIT 10 OFFSET 5`),
					gomock.Any(),
					5,
				).Return(rows, nil)

				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.client.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			whereParts: []mysql.WherePart{
				mysql.NewFilter("num", ">=", 5),
			},
			orders: []*mysql.Order{
				mysql.NewOrder("id", mysql.OrderDirectionAsc),
			},
			limit:  10,
			offset: 5,
			expected: []*proto.AccountV2{
				{
					Email:            "",
					Name:             "",
					AvatarImageUrl:   "",
					OrganizationId:   "",
					OrganizationRole: 0,
					EnvironmentRoles: nil,
					Disabled:         false,
					CreatedAt:        0,
					UpdatedAt:        0,
					SearchFilters:    nil,
				},
			},
			//expected: []*proto.AccountV2{
			//	{
			//		Email:            "a@email.com",
			//		Name:             "Mike",
			//		AvatarImageUrl:   "url",
			//		OrganizationId:   "org-0",
			//		OrganizationRole: 1,
			//		EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
			//			{
			//				EnvironmentId: "env-0",
			//				Role:          1,
			//			},
			//		},
			//		Disabled:      false,
			//		CreatedAt:     5,
			//		UpdatedAt:     6,
			//		SearchFilters: nil,
			//	},
			//},
			expectedCursor: 6,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAccountStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			accounts, cursor, _, err := storage.ListAccountsV2(
				context.Background(),
				p.whereParts,
				p.orders,
				p.limit,
				p.offset,
			)
			assert.Equal(t, p.expected, accounts)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAccountsMockV2(t *testing.T) {
	t.Parallel()

	//tttt := []proto.AccountV2_EnvironmentRole{
	//	{
	//		EnvironmentId: "env-0",
	//		Role:          1,
	//	},
	//}
	//json, _ := json.Marshal(tttt)
	tests := []struct {
		title          string
		setup          func(s *accountStorage, testMock sqlmock.Sqlmock)
		whereParts     []mysql.WherePart
		orders         []*mysql.Order
		limit          int
		offset         int
		expectedCursor int
		expectedErr    error
		expected       []*proto.AccountV2
	}{
		{
			title: "QueryContext Error",
			setup: func(s *accountStorage, testMock sqlmock.Sqlmock) {
				testMock.ExpectQuery("").WillReturnError(errors.New("error"))
			},
			whereParts:     nil,
			orders:         nil,
			limit:          0,
			offset:         0,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
			expected:       nil,
		},
		{
			title: "Success",
			setup: func(s *accountStorage, testMock sqlmock.Sqlmock) {
				columns1 := []string{
					"email",
					"name",
					"avatar_image_url",
					"organization_id",
					"organization_role",
					"environment_roles",
					"disabled",
					"created_at",
					"updated_at",
					"search_filters",
				}

				query1 := `SELECT email, name, avatar_image_url, organization_id, organization_role, environment_roles, disabled, created_at, updated_at, search_filters FROM account_v2 WHERE num >= ? ORDER BY id ASC LIMIT 10 OFFSET 5`
				mock1 := []driver.Value{
					"a@email.com",
					"Mike",
					"url",
					"org-0",
					proto.AccountV2_Role_Organization_OWNER,
					[]byte(`[{"environment_id":"env-0","role":1}]`),
					false,
					5,
					6,
					nil,
				}
				rows := sqlmock.NewRows(columns1).AddRow(mock1...)
				testMock.ExpectQuery(regexp.QuoteMeta(query1)).
					WithArgs(5).
					RowsWillBeClosed().
					WillReturnRows(rows)

				query2 := `SELECT COUNT(1) FROM account_v2 WHERE num >= ? ORDER BY id ASC`
				mockRow2 := []driver.Value{len(mock1)}
				columns2 := []string{"totalCount"}
				rows2 := sqlmock.NewRows(columns2).AddRow(mockRow2...)
				testMock.ExpectQuery(regexp.QuoteMeta(query2)).
					WithArgs(5).
					WillReturnRows(rows2)
			},
			whereParts: []mysql.WherePart{
				mysql.NewFilter("num", ">=", 5),
			},
			orders: []*mysql.Order{
				mysql.NewOrder("id", mysql.OrderDirectionAsc),
			},
			limit:          10,
			offset:         5,
			expectedCursor: 6,
			expectedErr:    nil,
			expected: []*proto.AccountV2{
				{
					Email:            "a@email.com",
					Name:             "Mike",
					AvatarImageUrl:   "url",
					OrganizationId:   "org-0",
					OrganizationRole: proto.AccountV2_Role_Organization_OWNER,
					EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
						{
							EnvironmentId: "env-0",
							Role:          1,
						},
					},
					Disabled:      false,
					CreatedAt:     5,
					UpdatedAt:     6,
					SearchFilters: nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			// mock
			storage, testMock := newAccountStorageWithMock2(t)
			defer storage.client.Close()

			if tt.setup != nil {
				tt.setup(storage, testMock)
			}
			//if tt.expectedErr != nil {
			//	testMock.ExpectQuery(regexp.QuoteMeta(tt.query)).
			//		WillReturnError(errors.New("error"))
			//} else {
			//	rows := sqlmock.NewRows(columns).AddRow(tt.mockRow...)
			//	testMock.ExpectQuery(regexp.QuoteMeta(tt.query)).
			//		WillReturnRows(rows)
			//}

			accounts, cursor, _, err := storage.ListAccountsV2(
				context.Background(),
				tt.whereParts,
				tt.orders,
				tt.limit,
				tt.offset,
			)

			// assert
			assert.Equal(t, tt.expectedErr, err)
			if err := testMock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			if err == nil {
				assert.Equal(t, tt.expected, accounts)
				assert.Equal(t, tt.expectedCursor, cursor)
			}
		})
	}
}

func newAccountStorageWithMock(t *testing.T, mockController *gomock.Controller) *accountStorage {
	t.Helper()
	return &accountStorage{mock.NewMockClient(mockController)}
}

type testClient struct {
	db *sql.DB
}

func (t *testClient) QueryContext(ctx context.Context, query string, args ...interface{}) (mysql.Rows, error) {
	return t.db.QueryContext(ctx, query, args...)
}

func (t *testClient) QueryRowContext(ctx context.Context, query string, args ...interface{}) mysql.Row {
	return t.db.QueryRowContext(ctx, query, args...)
}

func (t *testClient) ExecContext(ctx context.Context, query string, args ...interface{}) (mysql.Result, error) {
	return t.db.ExecContext(ctx, query, args...)
}

func (t *testClient) Close() error {
	return t.db.Close()
}

func (t *testClient) BeginTx(ctx context.Context) (mysql.Transaction, error) {
	panic("implement me")
}

func (t *testClient) RunInTransaction(ctx context.Context, tx mysql.Transaction, f func() error) error {
	panic("implement me")
}

func (t *testClient) TearDown() error {
	return t.db.Close()
}

func (t *testClient) NewRows(cloumns []string) *sqlmock.Rows {
	return sqlmock.NewRows(cloumns)
}

func newAccountStorageWithMock2(t *testing.T) (*accountStorage, sqlmock.Sqlmock) {
	t.Helper()
	db, m, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	client := &testClient{
		db: db,
	}

	return &accountStorage{client}, m
}

//func newAccountStorageWithMock2(t *testing.T, mockController *gomock.Controller) {
//	t.Helper()
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	client := &mysql.client{ {
//		db,
//
//	}
//
//	teardown := func() {
//		db.Close()
//	}
//
//
//
//	client := mysql.NewClient(
//		t,
//		db,
//		"",
//		"",
//		"",
//		mysql.WithLogger(logger)
//		)
//	return &accountStorage{client}
//}
