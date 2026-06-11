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
	_ "embed"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
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

type accountStorage struct {
	qe pgstorage.QueryExecer
}

func NewAccountStorage(qe pgstorage.QueryExecer) v2as.AccountStorage {
	return &accountStorage{qe}
}

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
		&pgstorage.JSONObject{Val: a.Tags},
		&pgstorage.JSONObject{Val: a.Teams},
		a.OrganizationId,
		int32(a.OrganizationRole),
		pgstorage.JSONObject{Val: a.EnvironmentRoles},
		a.Disabled,
		a.CreatedAt,
		a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrDuplicateEntry) {
			return v2as.ErrAccountAlreadyExists
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
		&pgstorage.JSONObject{Val: a.Tags},
		&pgstorage.JSONObject{Val: a.Teams},
		int32(a.OrganizationRole),
		pgstorage.JSONObject{Val: a.EnvironmentRoles},
		a.Disabled,
		a.UpdatedAt,
		a.LastSeen,
		pgstorage.JSONObject{Val: a.SearchFilters},
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
		return v2as.ErrAccountUnexpectedAffectedRows
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
		return v2as.ErrAccountUnexpectedAffectedRows
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
		&pgstorage.JSONObject{Val: &account.Tags},
		&pgstorage.JSONObject{Val: &account.Teams},
		&account.OrganizationId,
		&organizationRole,
		&pgstorage.JSONObject{Val: &account.EnvironmentRoles},
		&account.Disabled,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.LastSeen,
		&pgstorage.JSONObject{Val: &account.SearchFilters},
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, v2as.ErrAccountNotFound
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
		&pgstorage.JSONObject{Val: &account.Tags},
		&pgstorage.JSONObject{Val: &account.Teams},
		&account.OrganizationId,
		&organizationRole,
		&pgstorage.JSONObject{Val: &account.EnvironmentRoles},
		&account.Disabled,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.LastSeen,
		&pgstorage.JSONObject{Val: &account.SearchFilters},
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, v2as.ErrAccountNotFound
		}
		return nil, err
	}
	account.OrganizationRole = proto.AccountV2_Role_Organization(organizationRole)
	return &domain.AccountV2{AccountV2: &account}, nil
}

func (s *accountStorage) GetAvatarAccountsV2(
	ctx context.Context,
	params v2as.GetAvatarAccountsV2Params,
) ([]*proto.AccountV2, error) {
	emailsArg := make([]interface{}, len(params.Emails))
	for i, email := range params.Emails {
		emailsArg[i] = email
	}
	options := &pgstorage.ListOptions{
		Limit:  0,
		Offset: 0,
		InFilters: []*pgstorage.InFilter{
			{Column: "a.email", Values: emailsArg},
		},
		Filters: []*pgstorage.Filter{
			{Column: "e.id", Operator: pgstorage.OperatorEqual, Value: params.EnvironmentID},
		},
	}
	query, whereArgs := pgstorage.ConstructQueryAndWhereArgs(selectAvatarAccountsV2SQL, options)
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
		err := rows.Scan(
			&account.Email,
			&account.AvatarFileType,
			&account.AvatarImage,
		)
		if err != nil {
			return nil, err
		}
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
			&pgstorage.JSONObject{Val: &account.Tags},
			&pgstorage.JSONObject{Val: &account.Teams},
			&account.OrganizationId,
			&organizationRole,
			&pgstorage.JSONObject{Val: &account.EnvironmentRoles},
			&account.Disabled,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.LastSeen,
			&pgstorage.JSONObject{Val: &account.SearchFilters},
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
	params v2as.ListAccountsV2Params,
) ([]*proto.AccountV2, int, int64, error) {
	options, err := listAccountsV2OptionsFromParams(params)
	if err != nil {
		return nil, 0, 0, err
	}
	query, whereArgs := pgstorage.ConstructQueryAndWhereArgs(selectAccountsV2SQL, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	accounts := make([]*proto.AccountV2, 0, options.Limit)
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
			&pgstorage.JSONObject{Val: &account.Tags},
			&pgstorage.JSONObject{Val: &account.Teams},
			&account.OrganizationId,
			&organizationRole,
			&pgstorage.JSONObject{Val: &account.EnvironmentRoles},
			&account.Disabled,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.LastSeen,
			&pgstorage.JSONObject{Val: &account.SearchFilters},
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
	nextOffset := options.Offset + len(accounts)
	var totalCount int64
	countQuery, countWhereArgs := pgstorage.ConstructCountQuery(countAccountsV2SQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return accounts, nextOffset, totalCount, nil
}

type environmentRole struct {
	EnvironmentID *string `json:"environment_id"`
	Role          *int32  `json:"role"`
}

func listAccountsV2OptionsFromParams(p v2as.ListAccountsV2Params) (*pgstorage.ListOptions, error) {
	var filters []*pgstorage.Filter
	if p.OrganizationID != "" {
		filters = append(filters, &pgstorage.Filter{
			Column:   "organization_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.OrganizationID,
		})
	}
	if p.Disabled != nil {
		filters = append(filters, &pgstorage.Filter{
			Column:   "disabled",
			Operator: pgstorage.OperatorEqual,
			Value:    *p.Disabled,
		})
	}
	if p.OrganizationRole != nil {
		filters = append(filters, &pgstorage.Filter{
			Column:   "organization_role",
			Operator: pgstorage.OperatorEqual,
			Value:    *p.OrganizationRole,
		})
	}

	var jsonFilters []*pgstorage.JSONFilter
	if len(p.Tags) > 0 {
		tagValues := make([]interface{}, 0, len(p.Tags))
		for _, tag := range p.Tags {
			tagValues = append(tagValues, tag)
		}
		jsonFilters = append(jsonFilters, &pgstorage.JSONFilter{
			Column: "tags",
			Func:   pgstorage.JSONContainsString,
			Values: tagValues,
		})
	}
	if len(p.Teams) > 0 {
		teamValues := make([]interface{}, 0, len(p.Teams))
		for _, team := range p.Teams {
			teamValues = append(teamValues, team)
		}
		jsonFilters = append(jsonFilters, &pgstorage.JSONFilter{
			Column: "teams",
			Func:   pgstorage.JSONContainsString,
			Values: teamValues,
		})
	}

	var orFilters []*pgstorage.OrFilter
	if len(p.EnvironmentRoles) == 0 {
		// Admin user filtering: use JSONContainsJSON for a single environment role filter
		envRole := &environmentRole{}
		if p.EnvironmentID != nil {
			envRole.EnvironmentID = p.EnvironmentID
		}
		if p.EnvironmentRole != nil {
			envRole.Role = p.EnvironmentRole
		}
		jsonValues, err := json.Marshal(envRole)
		if err != nil {
			return nil, err
		}
		values := []interface{}{string(jsonValues)}
		if envRole.EnvironmentID != nil || envRole.Role != nil {
			jsonFilters = append(jsonFilters, &pgstorage.JSONFilter{
				Column: "environment_roles",
				Func:   pgstorage.JSONContainsJSON,
				Values: values,
			})
		}
	} else {
		// Member user filtering: build OR filter for each environment role + org_role >= admin
		orWhereParts := make([]pgstorage.WherePart, 0)
		for _, r := range p.EnvironmentRoles {
			role := int32(r.Role)
			envRole := &environmentRole{
				EnvironmentID: &r.EnvironmentId,
				Role:          &role,
			}
			jsonValues, err := json.Marshal(envRole)
			if err != nil {
				return nil, err
			}
			orWhereParts = append(orWhereParts, &pgstorage.JSONFilter{
				Column: "environment_roles",
				Func:   pgstorage.JSONContainsJSON,
				Values: []interface{}{string(jsonValues)},
			})
		}
		orWhereParts = append(
			orWhereParts,
			&pgstorage.Filter{
				Column:   "organization_role",
				Operator: pgstorage.OperatorGreaterThanOrEqual,
				Value:    proto.AccountV2_Role_Organization_ADMIN,
			},
		)
		orFilters = append(orFilters, &pgstorage.OrFilter{
			Queries: orWhereParts,
		})
	}

	var searchQuery *pgstorage.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &pgstorage.SearchQuery{
			Columns: []string{"email", "first_name", "last_name"},
			Keyword: p.SearchKeyword,
		}
	}

	var column string
	switch p.OrderBy {
	case proto.ListAccountsV2Request_DEFAULT,
		proto.ListAccountsV2Request_EMAIL:
		column = "email"
	case proto.ListAccountsV2Request_CREATED_AT:
		column = "created_at"
	case proto.ListAccountsV2Request_UPDATED_AT:
		column = "updated_at"
	case proto.ListAccountsV2Request_ORGANIZATION_ROLE:
		column = "organization_role"
	case proto.ListAccountsV2Request_ENVIRONMENT_COUNT:
		column = "environment_count"
	case proto.ListAccountsV2Request_LAST_SEEN:
		column = "last_seen"
	case proto.ListAccountsV2Request_STATE:
		column = "disabled"
	case proto.ListAccountsV2Request_TAGS:
		column = "tags"
	case proto.ListAccountsV2Request_TEAMS:
		column = "teams"
	default:
		return nil, v2as.ErrInvalidOrderBy
	}
	direction := pgstorage.OrderDirectionAsc
	if p.OrderDirection == proto.ListAccountsV2Request_DESC {
		direction = pgstorage.OrderDirectionDesc
	}

	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, v2as.ErrInvalidCursor
	}

	return &pgstorage.ListOptions{
		Limit:       p.PageSize,
		Offset:      offset,
		Filters:     filters,
		JSONFilters: jsonFilters,
		SearchQuery: searchQuery,
		OrFilters:   orFilters,
		Orders:      []*pgstorage.Order{pgstorage.NewOrder(column, direction)},
	}, nil
}
