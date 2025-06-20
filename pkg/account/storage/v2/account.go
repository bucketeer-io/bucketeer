// Copyright 2025 The Bucketeer Authors.
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
	//go:embed sql/account_v2/select_avatar_accounts_v2.sql
	selectAvatarAccountsV2SQL string
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
	_, err := s.qe.ExecContext(
		ctx,
		insertAccountV2SQL,
		a.Email,
		a.Name,
		a.FirstName,
		a.LastName,
		a.Language,
		a.AvatarImageUrl,
		a.AvatarFileType,
		a.AvatarImage,
		&mysql.JSONObject{Val: a.Tags},
		&mysql.JSONObject{Val: a.Teams},
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
	result, err := s.qe.ExecContext(
		ctx,
		updateAccountV2SQL,
		a.Name,
		a.FirstName,
		a.LastName,
		a.Language,
		a.AvatarImageUrl,
		a.AvatarFileType,
		a.AvatarImage,
		&mysql.JSONObject{Val: a.Tags},
		&mysql.JSONObject{Val: a.Teams},
		int32(a.OrganizationRole),
		mysql.JSONObject{Val: a.EnvironmentRoles},
		a.Disabled,
		a.UpdatedAt,
		a.LastSeen,
		mysql.JSONObject{Val: a.SearchFilters},
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
	result, err := s.qe.ExecContext(
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
	err := s.qe.QueryRowContext(
		ctx,
		selectAccountV2SQL,
		email,
		organizationID,
	).Scan(
		&account.Email,
		&account.Name,
		&account.FirstName,
		&account.LastName,
		&account.Language,
		&account.AvatarImageUrl,
		&account.AvatarFileType,
		&account.AvatarImage,
		&mysql.JSONObject{Val: &account.Tags},
		&mysql.JSONObject{Val: &account.Teams},
		&account.OrganizationId,
		&organizationRole,
		&mysql.JSONObject{Val: &account.EnvironmentRoles},
		&account.Disabled,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.LastSeen,
		&mysql.JSONObject{Val: &account.SearchFilters},
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
	err := s.qe.QueryRowContext(
		ctx,
		selectAccountV2ByEnvironmentIDSQL,
		email,
		environmentID,
	).Scan(
		&account.Email,
		&account.Name,
		&account.FirstName,
		&account.LastName,
		&account.Language,
		&account.AvatarImageUrl,
		&account.AvatarFileType,
		&account.AvatarImage,
		&mysql.JSONObject{Val: &account.Tags},
		&mysql.JSONObject{Val: &account.Teams},
		&account.OrganizationId,
		&organizationRole,
		&mysql.JSONObject{Val: &account.EnvironmentRoles},
		&account.Disabled,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.LastSeen,
		&mysql.JSONObject{Val: &account.SearchFilters},
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

func (s *accountStorage) GetAvatarAccountsV2(
	ctx context.Context,
	options *mysql.ListOptions,
) ([]*proto.AccountV2, error) {
	query, whereArgs := mysql.ConstructQueryAndWhereArgs(selectAvatarAccountsV2SQL, options)
	rows, err := s.qe.QueryContext(
		ctx,
		query,
		whereArgs...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	accounts := make([]*proto.AccountV2, 0)
	for rows.Next() {
		account := &proto.AccountV2{}
		var organizationRole int32
		err := rows.Scan(
			&account.Email,
			&account.AvatarFileType,
			&account.AvatarImage,
		)
		if err != nil {
			return nil, err
		}
		account.OrganizationRole = proto.AccountV2_Role_Organization(organizationRole)
		accounts = append(accounts, account)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return accounts, nil
}

func (s *accountStorage) GetAccountsWithOrganization(
	ctx context.Context,
	email string,
) ([]*domain.AccountWithOrganization, error) {
	rows, err := s.qe.QueryContext(ctx, selectAccountsWithOrganizationSQL, email)
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
			&account.FirstName,
			&account.LastName,
			&account.Language,
			&account.AvatarImageUrl,
			&account.AvatarFileType,
			&account.AvatarImage,
			&mysql.JSONObject{Val: &account.Tags},
			&mysql.JSONObject{Val: &account.Teams},
			&account.OrganizationId,
			&organizationRole,
			&mysql.JSONObject{Val: &account.EnvironmentRoles},
			&account.Disabled,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.LastSeen,
			&mysql.JSONObject{Val: &account.SearchFilters},
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
		return nil, rows.Err()
	}
	return accountsWithOrg, nil
}

func (s *accountStorage) ListAccountsV2(
	ctx context.Context,
	options *mysql.ListOptions,
) ([]*proto.AccountV2, int, int64, error) {
	query, whereArgs := mysql.ConstructQueryAndWhereArgs(selectAccountsV2SQL, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	var limit, offset int
	if options != nil {
		offset = options.Offset
		limit = options.Limit
	}
	accounts := make([]*proto.AccountV2, 0, limit)
	for rows.Next() {
		account := proto.AccountV2{}
		var organizationRole int32
		err := rows.Scan(
			&account.Email,
			&account.Name,
			&account.FirstName,
			&account.LastName,
			&account.Language,
			&account.AvatarImageUrl,
			&account.AvatarFileType,
			&account.AvatarImage,
			&mysql.JSONObject{Val: &account.Tags},
			&mysql.JSONObject{Val: &account.Teams},
			&account.OrganizationId,
			&organizationRole,
			&mysql.JSONObject{Val: &account.EnvironmentRoles},
			&account.Disabled,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.LastSeen,
			&mysql.JSONObject{Val: &account.SearchFilters},
			&account.EnvironmentCount,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		account.OrganizationRole = proto.AccountV2_Role_Organization(organizationRole)
		accounts = append(accounts, &account)
	}
	if rows.Err() != nil {
		return nil, 0, 0, rows.Err()
	}
	nextOffset := offset + len(accounts)
	var totalCount int64
	countQuery, countWhereArgs := mysql.ConstructCountQuery(countAccountsV2SQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return accounts, nextOffset, totalCount, nil
}
