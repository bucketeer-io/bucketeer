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
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestNewAccountStorage(t *testing.T) {
	t.Parallel()
	client, _, _ := mysql.NewSqlMockClient()
	storage := NewAccountStorage(client)
	assert.IsType(t, &accountStorage{}, storage)
}

func TestCreateAccountV2(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		setup       func(sqlmock.Sqlmock)
		account     *domain.AccountV2
		expectedErr error
	}{
		{
			desc: "ErrAccountAlreadyExists",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectExec("").WillReturnError(mysql.ErrDuplicateEntry)
			},
			account: &domain.AccountV2{
				AccountV2: &proto.AccountV2{},
			},
			expectedErr: ErrAccountAlreadyExists,
		},
		{
			desc: "Err: Other error",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectExec("").WillReturnError(mysql.ErrTxDone)
			},
			account: &domain.AccountV2{
				AccountV2: &proto.AccountV2{},
			},
			expectedErr: mysql.ErrTxDone,
		},
		{
			desc: "Success",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.
					ExpectExec(
						regexp.QuoteMeta(
							`INSERT INTO account_v2 ( email, name, avatar_image_url, organization_id, organization_role, environment_roles, disabled, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
						),
					).
					WithArgs(
						"email",
						"name",
						"avatarImageUrl",
						"organizationId",
						3,
						[]byte(`[{"environment_id":"env-0","role":1}]`),
						false,
						1,
						2,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			account: &domain.AccountV2{
				AccountV2: &proto.AccountV2{
					Email:            "email",
					Name:             "name",
					AvatarImageUrl:   "avatarImageUrl",
					OrganizationId:   "organizationId",
					OrganizationRole: proto.AccountV2_Role_Organization_OWNER,
					EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
						{
							EnvironmentId: "env-0",
							Role:          proto.AccountV2_Role_Environment_VIEWER,
						},
					},
					CreatedAt:     1,
					UpdatedAt:     2,
					Disabled:      false,
					SearchFilters: nil,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage, sMock := newAccountStorageWithMockClient(t)
			defer storage.client.Close()

			if p.setup != nil {
				p.setup(sMock)
			}
			err := storage.CreateAccountV2(context.Background(), p.account)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAccountV2(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		setup       func(sqlmock.Sqlmock)
		account     *domain.AccountV2
		expectedErr error
	}{
		{
			desc: "ErrTxDone",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectExec("").WillReturnError(mysql.ErrTxDone)
			},
			account: &domain.AccountV2{
				AccountV2: &proto.AccountV2{},
			},
			expectedErr: mysql.ErrTxDone,
		},
		{
			desc: "ErrAccountUnexpectedAffectedRows",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
			},
			account: &domain.AccountV2{
				AccountV2: &proto.AccountV2{},
			},
			expectedErr: ErrAccountUnexpectedAffectedRows,
		},
		{
			desc: "Success",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectExec(
					regexp.QuoteMeta(`UPDATE account_v2 SET name = ?, avatar_image_url = ?, organization_role = ?, environment_roles = ?, disabled = ?, updated_at = ?, search_filters = ? WHERE email = ? AND organization_id = ?`),
				).WithArgs(
					"name",
					"avatarImageUrl",
					3,
					[]byte(`[{"environment_id":"env-0","role":1}]`),
					false,
					2,
					[]byte(`[{"id":"searchId","name":"searchName","query":"searchQuery","filter_target_type":1,"environment_id":"envId","default_filter":false}]`),
					"email",
					"organizationId",
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			account: &domain.AccountV2{
				AccountV2: &proto.AccountV2{
					Email:            "email",
					Name:             "name",
					AvatarImageUrl:   "avatarImageUrl",
					OrganizationId:   "organizationId",
					OrganizationRole: proto.AccountV2_Role_Organization_OWNER,
					EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
						{
							EnvironmentId: "env-0",
							Role:          proto.AccountV2_Role_Environment_VIEWER,
						},
					},
					CreatedAt: 1,
					UpdatedAt: 2,
					Disabled:  false,
					SearchFilters: []*proto.SearchFilter{
						{
							Id:               "searchId",
							Name:             "searchName",
							Query:            "searchQuery",
							FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
							EnvironmentId:    "envId",
							DefaultFilter:    false,
						},
					},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage, sMock := newAccountStorageWithMockClient(t)
			defer storage.client.Close()

			if p.setup != nil {
				p.setup(sMock)
			}
			err := storage.UpdateAccountV2(context.Background(), p.account)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteAccountV2(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		setup       func(sqlmock.Sqlmock)
		account     *domain.AccountV2
		expectedErr error
	}{
		{
			desc: "ErrTxDone",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectExec("").WillReturnError(mysql.ErrTxDone)
			},
			account: &domain.AccountV2{
				AccountV2: &proto.AccountV2{},
			},
			expectedErr: mysql.ErrTxDone,
		},
		{
			desc: "ErrAccountUnexpectedAffectedRows",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
			},
			account: &domain.AccountV2{
				AccountV2: &proto.AccountV2{},
			},
			expectedErr: ErrAccountUnexpectedAffectedRows,
		},
		{
			desc: "Success",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectExec(
					regexp.QuoteMeta(`DELETE FROM account_v2 WHERE email = ? AND organization_id = ?`),
				).WithArgs(
					"email",
					"organizationId",
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			account: &domain.AccountV2{
				AccountV2: &proto.AccountV2{
					Email:          "email",
					OrganizationId: "organizationId",
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage, sMock := newAccountStorageWithMockClient(t)
			defer storage.client.Close()

			if p.setup != nil {
				p.setup(sMock)
			}
			err := storage.DeleteAccountV2(context.Background(), p.account)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAccountV2(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc            string
		setup           func(sqlmock.Sqlmock)
		account         *domain.AccountV2
		expectedAccount *domain.AccountV2
		expectedErr     error
	}{
		{
			desc: "ErrAccountNotFound",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectQuery("").WillReturnError(mysql.ErrNoRows)
			},
			account: &domain.AccountV2{
				AccountV2: &proto.AccountV2{},
			},
			expectedAccount: nil,
			expectedErr:     ErrAccountNotFound,
		},
		{
			desc: "Error: Other error",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectQuery("").WillReturnError(mysql.ErrTxDone)
			},
			account: &domain.AccountV2{
				AccountV2: &proto.AccountV2{},
			},
			expectedAccount: nil,
			expectedErr:     mysql.ErrTxDone,
		},
		{
			desc: "Success",
			setup: func(sMock sqlmock.Sqlmock) {
				columns := []string{
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

				mockRow := sqlmock.
					NewRows(columns).
					AddRow(
						"email",
						"name",
						"avatarImageUrl",
						"organizationId",
						3,
						[]byte(`[{"environment_id":"env-0","role":2}]`),
						false,
						1,
						2,
						nil,
					)
				sMock.ExpectQuery(
					regexp.QuoteMeta(
						`SELECT email, name, avatar_image_url, organization_id, organization_role, environment_roles, disabled, created_at, updated_at, search_filters FROM account_v2 WHERE email = ? AND organization_id = ?`,
					),
				).WithArgs(
					"email",
					"organizationId",
				).WillReturnRows(mockRow)
			},
			account: &domain.AccountV2{
				AccountV2: &proto.AccountV2{
					Email:          "email",
					OrganizationId: "organizationId",
				},
			},
			expectedErr: nil,
			expectedAccount: &domain.AccountV2{
				AccountV2: &proto.AccountV2{
					Email:            "email",
					Name:             "name",
					OrganizationId:   "organizationId",
					OrganizationRole: proto.AccountV2_Role_Organization_OWNER,
					AvatarImageUrl:   "avatarImageUrl",
					EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
						{
							EnvironmentId: "env-0",
							Role:          proto.AccountV2_Role_Environment_EDITOR,
						},
					},
					Disabled:      false,
					CreatedAt:     1,
					UpdatedAt:     2,
					SearchFilters: nil,
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage, sMock := newAccountStorageWithMockClient(t)
			defer storage.client.Close()

			if p.setup != nil {
				p.setup(sMock)
			}
			account, err := storage.GetAccountV2(context.Background(), p.account.Email, p.account.OrganizationId)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expectedAccount, account)
		})
	}
}

func TestGetAccountV2ByEnvironmentID(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc            string
		setup           func(sqlmock.Sqlmock)
		account         *domain.AccountV2
		expectedAccount *domain.AccountV2
		expectedErr     error
	}{
		{
			desc: "ErrAccountNotFound",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectQuery("").WillReturnError(mysql.ErrNoRows)
			},
			account: &domain.AccountV2{
				AccountV2: &proto.AccountV2{},
			},
			expectedAccount: nil,
			expectedErr:     ErrAccountNotFound,
		},
		{
			desc: "Error: Other error",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectQuery("").WillReturnError(mysql.ErrTxDone)
			},
			account: &domain.AccountV2{
				AccountV2: &proto.AccountV2{},
			},
			expectedAccount: nil,
			expectedErr:     mysql.ErrTxDone,
		},
		{
			desc: "Success",
			setup: func(sMock sqlmock.Sqlmock) {
				columns := []string{
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

				mockRow := sqlmock.
					NewRows(columns).
					AddRow(
						"email",
						"name",
						"avatarImageUrl",
						"organizationId",
						3,
						[]byte(`[{"environment_id":"env-0","role":2}]`),
						false,
						1,
						2,
						nil,
					)
				sMock.ExpectQuery(
					regexp.QuoteMeta(
						`SELECT a.email, a.name, a.avatar_image_url, a.organization_id, a.organization_role, a.environment_roles, a.disabled, a.created_at, a.updated_at, a.search_filters FROM account_v2 AS a INNER JOIN environment_v2 AS e ON a.organization_id = e.organization_id WHERE a.email = ? AND e.id = ?`,
					),
				).WithArgs(
					"email",
					"organizationId",
				).WillReturnRows(mockRow)
			},
			account: &domain.AccountV2{
				AccountV2: &proto.AccountV2{
					Email:          "email",
					OrganizationId: "organizationId",
				},
			},
			expectedErr: nil,
			expectedAccount: &domain.AccountV2{
				AccountV2: &proto.AccountV2{
					Email:            "email",
					Name:             "name",
					OrganizationId:   "organizationId",
					OrganizationRole: proto.AccountV2_Role_Organization_OWNER,
					AvatarImageUrl:   "avatarImageUrl",
					EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
						{
							EnvironmentId: "env-0",
							Role:          proto.AccountV2_Role_Environment_EDITOR,
						},
					},
					Disabled:      false,
					CreatedAt:     1,
					UpdatedAt:     2,
					SearchFilters: nil,
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage, sMock := newAccountStorageWithMockClient(t)
			defer storage.client.Close()

			if p.setup != nil {
				p.setup(sMock)
			}
			account, err := storage.GetAccountV2ByEnvironmentID(context.Background(), p.account.Email, p.account.OrganizationId)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expectedAccount, account)
		})
	}
}

func TestGetAccountsWithOrganization(t *testing.T) {
	t.Parallel()
	query := `SELECT a.email, a.name, a.avatar_image_url, a.organization_id, a.organization_role, a.environment_roles, a.disabled, a.created_at, a.updated_at, a.search_filters, o.id, o.name, o.url_code, o.description, o.disabled, o.archived, o.trial, o.system_admin, o.created_at, o.updated_at FROM account_v2 AS a INNER JOIN organization AS o ON a.organization_id=o.id WHERE email=?`

	patterns := []struct {
		desc        string
		setup       func(sqlmock.Sqlmock)
		email       string
		expected    []*domain.AccountWithOrganization
		expectedErr error
	}{
		{
			desc: "ErrNoRows",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectQuery("").WithArgs("email").WillReturnError(mysql.ErrNoRows)
			},
			email:       "email",
			expected:    nil,
			expectedErr: mysql.ErrNoRows,
		},
		{
			desc: "Success: No Rows",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{}))
			},
			email:       "email",
			expected:    []*domain.AccountWithOrganization{},
			expectedErr: nil,
		},
		{
			desc: "Success",
			setup: func(sMock sqlmock.Sqlmock) {
				columns := []string{
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
					"organization_id",
					"organization_name",
					"url_code",
					"description",
					"organization_disabled",
					"archived",
					"trial",
					"system_admin",
					"organization_created_at",
					"organization_updated_at",
				}

				mockRow := sqlmock.
					NewRows(columns).
					AddRow(
						"email",
						"name",
						"avatarImageUrl",
						"organizationId",
						3,
						[]byte(`[{"environment_id":"env-0","role":2}]`),
						false,
						1,
						2,
						nil,
						"organizationId",
						"organizationName",
						"urlCode",
						"description",
						false,
						false,
						false,
						false,
						4,
						5,
					)
				sMock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("email").WillReturnRows(mockRow)
			},
			email:       "email",
			expectedErr: nil,
			expected: []*domain.AccountWithOrganization{
				{
					AccountV2: &proto.AccountV2{
						Email:            "email",
						Name:             "name",
						OrganizationId:   "organizationId",
						OrganizationRole: proto.AccountV2_Role_Organization_OWNER,
						AvatarImageUrl:   "avatarImageUrl",
						EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "env-0",
								Role:          proto.AccountV2_Role_Environment_EDITOR,
							},
						},
						Disabled:      false,
						CreatedAt:     1,
						UpdatedAt:     2,
						SearchFilters: nil,
					},
					Organization: &environmentproto.Organization{
						Id:          "organizationId",
						Name:        "organizationName",
						UrlCode:     "urlCode",
						Description: "description",
						Disabled:    false,
						Archived:    false,
						Trial:       false,
						SystemAdmin: false,
						CreatedAt:   4,
						UpdatedAt:   5,
					},
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage, sMock := newAccountStorageWithMockClient(t)
			defer storage.client.Close()

			if p.setup != nil {
				p.setup(sMock)
			}
			organizations, err := storage.GetAccountsWithOrganization(context.Background(), p.email)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, organizations)
		})
	}
}

func TestListAccountsV2(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc               string
		whereParts         []mysql.WherePart
		orders             []*mysql.Order
		limit              int
		offset             int
		setup              func(sqlmock.Sqlmock)
		expected           []*proto.AccountV2
		expectedCursor     int
		expectedErr        error
		expectedTotalCount int64
	}{
		{
			desc:       "Error: Select Accounts",
			whereParts: nil,
			orders:     nil,
			limit:      0,
			offset:     0,
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectQuery("").WillReturnError(errors.New("error"))
			},
			expected:           nil,
			expectedCursor:     0,
			expectedTotalCount: 0,
			expectedErr:        errors.New("error"),
		},
		{
			whereParts: nil,
			orders:     nil,
			limit:      0,
			offset:     0,
			desc:       "Error: TotalCount Select",
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{}))
				sMock.ExpectQuery("").WillReturnError(errors.New("error"))
			},
			expected:           nil,
			expectedCursor:     0,
			expectedTotalCount: 0,
			expectedErr:        errors.New("error"),
		},
		{
			desc: "Success: No Rows",
			whereParts: []mysql.WherePart{
				mysql.NewFilter("num", ">=", 5),
			},
			orders: []*mysql.Order{
				mysql.NewOrder("id", mysql.OrderDirectionAsc),
			},
			limit:  10,
			offset: 0,
			setup: func(sMock sqlmock.Sqlmock) {
				sMock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{}))
				sMock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"totalCount"}).AddRow(0))
			},
			expected:           []*proto.AccountV2{},
			expectedCursor:     0,
			expectedTotalCount: 0,
			expectedErr:        nil,
		},
		{
			desc: "Success",
			whereParts: []mysql.WherePart{
				mysql.NewFilter("num", ">=", 5),
			},
			orders: []*mysql.Order{
				mysql.NewOrder("id", mysql.OrderDirectionAsc),
			},
			limit:  10,
			offset: 5,
			setup: func(sMock sqlmock.Sqlmock) {
				selectColumns := []string{
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
				selectMock1 := []driver.Value{
					"email",
					"name",
					"avatarImageUrl",
					"organizationId",
					3,
					[]byte(`[{"environment_id":"env-0","role":2}]`),
					false,
					1,
					2,
					nil,
				}
				selectMock2 := []driver.Value{
					"email2",
					"name2",
					"avatarImageUrl2",
					"organizationId2",
					1,
					[]byte(`[{"environment_id":"env-2","role":1}]`),
					true,
					7,
					8,
					[]byte(`[{"id":"searchId", "name":"searchName", "query": "searchQuery", "filter_target_type": 1, "environment_id": "envId", "default_filter": false }]`),
				}

				selectRows := sqlmock.NewRows(selectColumns).AddRows(selectMock1, selectMock2)

				selectQuery := `SELECT email, name, avatar_image_url, organization_id, organization_role, environment_roles, disabled, created_at, updated_at, search_filters FROM account_v2 WHERE num >= ? ORDER BY id ASC LIMIT 10 OFFSET 5`
				sMock.ExpectQuery(regexp.QuoteMeta(selectQuery)).
					WithArgs(5).
					WillReturnRows(selectRows)

				totalCountQuery := `SELECT COUNT(1) FROM account_v2 WHERE num >= ? ORDER BY id ASC`
				totalCountRows := sqlmock.NewRows([]string{"totalCount"}).AddRows([]driver.Value{7})
				sMock.ExpectQuery(regexp.QuoteMeta(totalCountQuery)).
					WithArgs(5).
					WillReturnRows(totalCountRows)
			},
			expectedErr:        nil,
			expectedCursor:     7,
			expectedTotalCount: 7,
			expected: []*proto.AccountV2{
				{
					Email:            "email",
					Name:             "name",
					OrganizationId:   "organizationId",
					OrganizationRole: proto.AccountV2_Role_Organization_OWNER,
					AvatarImageUrl:   "avatarImageUrl",
					EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
						{
							EnvironmentId: "env-0",
							Role:          proto.AccountV2_Role_Environment_EDITOR,
						},
					},
					Disabled:      false,
					CreatedAt:     1,
					UpdatedAt:     2,
					SearchFilters: nil,
				},
				{
					Email:            "email2",
					Name:             "name2",
					OrganizationId:   "organizationId2",
					OrganizationRole: proto.AccountV2_Role_Organization_MEMBER,
					AvatarImageUrl:   "avatarImageUrl2",
					EnvironmentRoles: []*proto.AccountV2_EnvironmentRole{
						{
							EnvironmentId: "env-2",
							Role:          proto.AccountV2_Role_Environment_VIEWER,
						},
					},
					Disabled:  true,
					CreatedAt: 7,
					UpdatedAt: 8,
					SearchFilters: []*proto.SearchFilter{
						{
							Id:               "searchId",
							Name:             "searchName",
							Query:            "searchQuery",
							FilterTargetType: proto.FilterTargetType_FEATURE_FLAG,
							EnvironmentId:    "envId",
							DefaultFilter:    false,
						},
					},
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage, sMock := newAccountStorageWithMockClient(t)
			defer storage.client.Close()

			if p.setup != nil {
				p.setup(sMock)
			}
			accounts, cursor, totalCount, err := storage.ListAccountsV2(
				context.Background(),
				p.whereParts,
				p.orders,
				p.limit,
				p.offset,
			)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, accounts)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedTotalCount, totalCount)
		})
	}
}

func newAccountStorageWithMockClient(t *testing.T) (*accountStorage, sqlmock.Sqlmock) {
	t.Helper()
	client, sqlMock, err := mysql.NewSqlMockClient()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return &accountStorage{client}, sqlMock
}
