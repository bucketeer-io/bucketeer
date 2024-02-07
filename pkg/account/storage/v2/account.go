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
	_ "embed"
	"errors"
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

var (
	//go:embed sql/account_v2/insert_account_v2.sql
	insertAccountV2SQL string
	//go:embed sql/account_v2/update_account_v2.sql
	updateAccountV2SQL string
	//go:embed sql/account_v2/delete_account_v2.sql
	deleteAccountV2SQL string
	//go:embed sql/account_v2/select_account_v2.sql
	selectAccountV2SQL string
	//go:embed sql/account_v2/select_account_v2_by_environment_id.sql
	selectAccountV2ByEnvironmentIDSQL string
	//go:embed sql/account_v2/select_accounts_v2.sql
	selectAccountsV2SQL string
	//go:embed sql/account_v2/count_accounts_v2.sql
	countAccountsV2SQL string
	//go:embed sql/account_v2/select_accounts_with_organization.sql
	selectAccountsWithOrganizationSQL string
)

var (
	ErrAccountAlreadyExists          = errors.New("account: account already exists")
	ErrAccountNotFound               = errors.New("account: account not found")
	ErrAccountUnexpectedAffectedRows = errors.New("account: account unexpected affected rows")
)

func (s *accountStorage) CreateAccountV2(ctx context.Context, a *domain.AccountV2) error {
	_, err := s.qe(ctx).ExecContext(
		ctx,
		insertAccountV2SQL,
		a.Email,
		a.Name,
		a.AvatarImageUrl,
		a.OrganizationId,
		int32(a.OrganizationRole),
		mysql.JSONObject{Val: a.EnvironmentRoles},
		a.Disabled,
		a.CreatedAt,
		a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrDuplicateEntry) {
			return ErrAccountAlreadyExists
		}
		return err
	}
	return nil
}

func (s *accountStorage) UpdateAccountV2(ctx context.Context, a *domain.AccountV2) error {
	result, err := s.qe(ctx).ExecContext(
		ctx,
		updateAccountV2SQL,
		a.Name,
		a.AvatarImageUrl,
		int32(a.OrganizationRole),
		mysql.JSONObject{Val: a.EnvironmentRoles},
		a.Disabled,
		a.UpdatedAt,
		a.Email,
		a.OrganizationId,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrAccountUnexpectedAffectedRows
	}
	return nil
}

func (s *accountStorage) DeleteAccountV2(ctx context.Context, a *domain.AccountV2) error {
	result, err := s.qe(ctx).ExecContext(
		ctx,
		deleteAccountV2SQL,
		a.Email,
		a.OrganizationId,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrAccountUnexpectedAffectedRows
	}
	return nil
}

func (s *accountStorage) GetAccountV2(ctx context.Context, email, organizationID string) (*domain.AccountV2, error) {
	account := proto.AccountV2{}
	var organizationRole int32
	err := s.qe(ctx).QueryRowContext(
		ctx,
		selectAccountV2SQL,
		email,
		organizationID,
	).Scan(
		&account.Email,
		&account.Name,
		&account.AvatarImageUrl,
		&account.OrganizationId,
		&organizationRole,
		&mysql.JSONObject{Val: &account.EnvironmentRoles},
		&account.Disabled,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	account.OrganizationRole = proto.AccountV2_Role_Organization(organizationRole)
	return &domain.AccountV2{AccountV2: &account}, nil
}

func (s *accountStorage) GetAccountV2ByEnvironmentID(
	ctx context.Context,
	email, environmentID string,
) (*domain.AccountV2, error) {
	account := proto.AccountV2{}
	var organizationRole int32
	err := s.qe(ctx).QueryRowContext(
		ctx,
		selectAccountV2ByEnvironmentIDSQL,
		email,
		environmentID,
	).Scan(
		&account.Email,
		&account.Name,
		&account.AvatarImageUrl,
		&account.OrganizationId,
		&organizationRole,
		&mysql.JSONObject{Val: &account.EnvironmentRoles},
		&account.Disabled,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	account.OrganizationRole = proto.AccountV2_Role_Organization(organizationRole)
	return &domain.AccountV2{AccountV2: &account}, nil
}

func (s *accountStorage) GetAccountsWithOrganization(
	ctx context.Context,
	email string,
) ([]*domain.AccountWithOrganization, error) {
	rows, err := s.qe(ctx).QueryContext(ctx, selectAccountsWithOrganizationSQL, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	accountsWithOrg := make([]*domain.AccountWithOrganization, 0)
	for rows.Next() {
		account := proto.AccountV2{}
		organization := environmentproto.Organization{}
		var organizationRole int32
		err := rows.Scan(
			&account.Email,
			&account.Name,
			&account.AvatarImageUrl,
			&account.OrganizationId,
			&organizationRole,
			&mysql.JSONObject{Val: &account.EnvironmentRoles},
			&account.Disabled,
			&account.CreatedAt,
			&account.UpdatedAt,
			&organization.Id,
			&organization.Name,
			&organization.UrlCode,
			&organization.Description,
			&organization.Disabled,
			&organization.Archived,
			&organization.Trial,
			&organization.SystemAdmin,
			&organization.CreatedAt,
			&organization.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		account.OrganizationRole = proto.AccountV2_Role_Organization(organizationRole)
		accountsWithOrg = append(accountsWithOrg, &domain.AccountWithOrganization{
			AccountV2:    &account,
			Organization: &organization,
		})
	}
	if rows.Err() != nil {
		return nil, err
	}
	return accountsWithOrg, nil
}

func (s *accountStorage) ListAccountsV2(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.AccountV2, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(
		selectAccountsV2SQL,
		whereSQL,
		orderBySQL,
		limitOffsetSQL,
	)
	rows, err := s.qe(ctx).QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	accounts := make([]*proto.AccountV2, 0, limit)
	for rows.Next() {
		account := proto.AccountV2{}
		var organizationRole int32
		err := rows.Scan(
			&account.Email,
			&account.Name,
			&account.AvatarImageUrl,
			&account.OrganizationId,
			&organizationRole,
			&mysql.JSONObject{Val: &account.EnvironmentRoles},
			&account.Disabled,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		account.OrganizationRole = proto.AccountV2_Role_Organization(organizationRole)
		accounts = append(accounts, &account)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(accounts)
	var totalCount int64
	countQuery := fmt.Sprintf(countAccountsV2SQL, whereSQL, orderBySQL)
	err = s.qe(ctx).QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return accounts, nextOffset, totalCount, nil
}
