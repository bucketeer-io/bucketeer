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
	"errors"
	"strconv"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

var (
	//go:embed sql/api_key_v2/insert_api_key_v2.sql
	insertAPIKeyV2SQLQuery string
	//go:embed sql/api_key_v2/update_api_key_v2.sql
	updateAPIKeyV2SQLQuery string
	//go:embed sql/api_key_v2/update_api_key_v2_last_used_at.sql
	updateAPIKeyLastUsedAtV2SQLQuery string
	//go:embed sql/api_key_v2/select_api_key_v2.sql
	selectAPIKeyV2SQLQuery string
	//go:embed sql/api_key_v2/select_api_key_v2_count.sql
	selectAPIKeyV2CountSQLQuery string
	//go:embed sql/api_key_v2/select_api_key_v2_by_api_key.sql
	selectAPIKeyV2ByAPIKeySQLQuery string
	//go:embed sql/api_key_v2/select_environment_api_key_v2.sql
	selectEnvironmentAPIKeySQLQuery string
	//go:embed sql/api_key_v2/select_all_environment_api_keys_v2.sql
	selectAllEnvironmentAPIKeysSQLQuery string
	//go:embed sql/api_key_v2/select_api_key_v2_by_id.sql
	selectAPIKeyV2ByIDSQLQuery string
)

func (s *accountStorage) CreateAPIKey(ctx context.Context, k *domain.APIKey, environmentID string) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertAPIKeyV2SQLQuery,
		k.Id,
		k.Name,
		int32(k.Role),
		k.Disabled,
		k.CreatedAt,
		k.UpdatedAt,
		environmentID,
		k.ApiKey,
		k.Maintainer,
		k.Description,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrDuplicateEntry) {
			return v2as.ErrAPIKeyAlreadyExists
		}
		return err
	}
	return nil
}

func (s *accountStorage) UpdateAPIKey(ctx context.Context, k *domain.APIKey, environmentID string) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateAPIKeyV2SQLQuery,
		k.Name,
		int32(k.Role),
		k.Disabled,
		k.Maintainer,
		k.Description,
		k.UpdatedAt,
		k.Id,
		environmentID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return v2as.ErrAPIKeyUnexpectedAffectedRows
	}
	return nil
}

func (s *accountStorage) UpdateAPIKeyLastUsedAt(
	ctx context.Context,
	id, environmentID string,
	lastUsedAt int64,
) (bool, error) {
	result, err := s.qe.ExecContext(
		ctx,
		updateAPIKeyLastUsedAtV2SQLQuery,
		lastUsedAt,
		id,
		environmentID,
		lastUsedAt,
	)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}

func (s *accountStorage) GetAPIKey(ctx context.Context, id, environmentID string) (*domain.APIKey, error) {
	apiKey := proto.APIKey{}
	var role int32
	err := s.qe.QueryRowContext(
		ctx,
		selectAPIKeyV2ByIDSQLQuery,
		id,
		environmentID,
	).Scan(
		&apiKey.Id,
		&apiKey.Name,
		&role,
		&apiKey.Disabled,
		&apiKey.CreatedAt,
		&apiKey.UpdatedAt,
		&apiKey.Description,
		&apiKey.ApiKey,
		&apiKey.Maintainer,
		&apiKey.LastUsedAt,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, v2as.ErrAPIKeyNotFound
		}
		return nil, err
	}
	apiKey.Role = proto.APIKey_Role(role)
	return &domain.APIKey{APIKey: &apiKey}, nil
}

func (s *accountStorage) GetAPIKeyByAPIKey(
	ctx context.Context,
	apiKey string,
	environmentID string,
) (*domain.APIKey, error) {
	apiKeyDB := proto.APIKey{}
	var role int32
	err := s.qe.QueryRowContext(
		ctx,
		selectAPIKeyV2ByAPIKeySQLQuery,
		apiKey,
		environmentID,
	).Scan(
		&apiKeyDB.Id,
		&apiKeyDB.Name,
		&role,
		&apiKeyDB.Disabled,
		&apiKeyDB.CreatedAt,
		&apiKeyDB.UpdatedAt,
		&apiKeyDB.Description,
		&apiKeyDB.ApiKey,
		&apiKeyDB.Maintainer,
		&apiKeyDB.LastUsedAt,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, v2as.ErrAPIKeyNotFound
		}
		return nil, err
	}
	apiKeyDB.Role = proto.APIKey_Role(role)
	return &domain.APIKey{APIKey: &apiKeyDB}, nil
}

func (s *accountStorage) ListAllEnvironmentAPIKeys(
	ctx context.Context,
) ([]*domain.EnvironmentAPIKey, error) {
	rows, err := s.qe.QueryContext(ctx, selectAllEnvironmentAPIKeysSQLQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	envApiKeys := []*domain.EnvironmentAPIKey{}
	for rows.Next() {
		apiKeyDB := proto.APIKey{}
		envDB := envproto.EnvironmentV2{}
		envApiKeyDB := proto.EnvironmentAPIKey{}
		var role int32
		err := rows.Scan(
			// API Key columns
			&apiKeyDB.Id,
			&apiKeyDB.Name,
			&role,
			&apiKeyDB.Disabled,
			&apiKeyDB.CreatedAt,
			&apiKeyDB.UpdatedAt,
			&apiKeyDB.Description,
			&apiKeyDB.ApiKey,
			&apiKeyDB.Maintainer,
			&apiKeyDB.LastUsedAt,

			// Environment columns
			&envDB.Id,
			&envDB.Name,
			&envDB.UrlCode,
			&envDB.Description,
			&envDB.ProjectId,
			&envDB.OrganizationId,
			&envDB.Archived,
			&envDB.RequireComment,
			&envDB.CreatedAt,
			&envDB.UpdatedAt,

			// Project columns
			&envApiKeyDB.ProjectId,
			&envApiKeyDB.ProjectUrlCode,
			&envApiKeyDB.EnvironmentDisabled,
		)
		if err != nil {
			return nil, err
		}
		envApiKeyDB.ApiKey = &apiKeyDB
		envApiKeyDB.ApiKey.Role = proto.APIKey_Role(role)
		envApiKeyDB.Environment = &envDB
		envApiKeys = append(envApiKeys, &domain.EnvironmentAPIKey{EnvironmentAPIKey: &envApiKeyDB})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return envApiKeys, nil
}

func (s *accountStorage) GetEnvironmentAPIKey(
	ctx context.Context,
	apiKey string,
) (*domain.EnvironmentAPIKey, error) {
	apiKeyDB := proto.APIKey{}
	envDB := envproto.EnvironmentV2{}
	envApiKeyDB := proto.EnvironmentAPIKey{}
	var role int32
	err := s.qe.QueryRowContext(
		ctx,
		selectEnvironmentAPIKeySQLQuery,
		apiKey,
	).Scan(
		// API Key columns
		&apiKeyDB.Id,
		&apiKeyDB.Name,
		&role,
		&apiKeyDB.Disabled,
		&apiKeyDB.CreatedAt,
		&apiKeyDB.UpdatedAt,
		&apiKeyDB.Description,
		&apiKeyDB.ApiKey,
		&apiKeyDB.Maintainer,
		&apiKeyDB.LastUsedAt,

		// Environment columns
		&envDB.Id,
		&envDB.Name,
		&envDB.UrlCode,
		&envDB.Description,
		&envDB.ProjectId,
		&envDB.OrganizationId,
		&envDB.Archived,
		&envDB.RequireComment,
		&envDB.CreatedAt,
		&envDB.UpdatedAt,

		// Project columns
		&envApiKeyDB.ProjectId,
		&envApiKeyDB.ProjectUrlCode,
		&envApiKeyDB.EnvironmentDisabled,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, v2as.ErrAPIKeyNotFound
		}
		return nil, err
	}
	envApiKeyDB.ApiKey = &apiKeyDB
	envApiKeyDB.ApiKey.Role = proto.APIKey_Role(role)
	envApiKeyDB.Environment = &envDB
	return &domain.EnvironmentAPIKey{EnvironmentAPIKey: &envApiKeyDB}, nil
}

func (s *accountStorage) ListAPIKeys(
	ctx context.Context,
	params v2as.ListAPIKeysParams,
) ([]*proto.APIKey, int, int64, error) {
	options, err := listAPIKeysOptionsFromParams(params)
	if err != nil {
		return nil, 0, 0, err
	}
	query, whereArgs := pgstorage.ConstructQueryAndWhereArgs(selectAPIKeyV2SQLQuery, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	apiKeys := make([]*proto.APIKey, 0, options.Limit)
	for rows.Next() {
		apiKey := proto.APIKey{}
		var role int32
		err := rows.Scan(
			&apiKey.Id,
			&apiKey.Name,
			&role,
			&apiKey.Disabled,
			&apiKey.CreatedAt,
			&apiKey.UpdatedAt,
			&apiKey.Description,
			&apiKey.EnvironmentName,
			&apiKey.ApiKey,
			&apiKey.Maintainer,
			&apiKey.LastUsedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		apiKey.Role = proto.APIKey_Role(role)
		apiKeys = append(apiKeys, &apiKey)
	}
	if rows.Err() != nil {
		return nil, 0, 0, rows.Err()
	}
	nextOffset := options.Offset + len(apiKeys)
	var totalCount int64
	countQuery, countWhereArgs := pgstorage.ConstructCountQuery(selectAPIKeyV2CountSQLQuery, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return apiKeys, nextOffset, totalCount, nil
}

func listAPIKeysOptionsFromParams(p v2as.ListAPIKeysParams) (*pgstorage.ListOptions, error) {
	var filters []*pgstorage.Filter
	if p.OrganizationID != "" {
		filters = append(filters, &pgstorage.Filter{
			Column:   "environment_v2.organization_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.OrganizationID,
		})
	}
	if p.Disabled != nil {
		filters = append(filters, &pgstorage.Filter{
			Column:   "api_key.disabled",
			Operator: pgstorage.OperatorEqual,
			Value:    *p.Disabled,
		})
	}
	if p.MaintainerEmail != "" {
		filters = append(filters, &pgstorage.Filter{
			Column:   "api_key.maintainer",
			Operator: pgstorage.OperatorEqual,
			Value:    p.MaintainerEmail,
		})
	}

	var inFilters []*pgstorage.InFilter
	if len(p.EnvironmentIDs) > 0 {
		environmentIDs := make([]interface{}, 0, len(p.EnvironmentIDs))
		for _, id := range p.EnvironmentIDs {
			environmentIDs = append(environmentIDs, id)
		}
		inFilters = append(inFilters, &pgstorage.InFilter{
			Column: "api_key.environment_id",
			Values: environmentIDs,
		})
	}

	var searchQuery *pgstorage.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &pgstorage.SearchQuery{
			Columns: []string{"api_key.name"},
			Keyword: p.SearchKeyword,
		}
	}

	var column string
	switch p.OrderBy {
	case proto.ListAPIKeysRequest_DEFAULT,
		proto.ListAPIKeysRequest_NAME:
		column = "api_key.name"
	case proto.ListAPIKeysRequest_CREATED_AT:
		column = "api_key.created_at"
	case proto.ListAPIKeysRequest_UPDATED_AT:
		column = "api_key.updated_at"
	case proto.ListAPIKeysRequest_ROLE:
		column = "api_key.role"
	case proto.ListAPIKeysRequest_ENVIRONMENT:
		column = "environment_v2.name"
	case proto.ListAPIKeysRequest_STATE:
		column = "api_key.disabled"
	case proto.ListAPIKeysRequest_LAST_USED_AT:
		column = "api_key.last_used_at"
	default:
		return nil, v2as.ErrInvalidOrderBy
	}
	direction := pgstorage.OrderDirectionAsc
	if p.OrderDirection == proto.ListAPIKeysRequest_DESC {
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
		InFilters:   inFilters,
		SearchQuery: searchQuery,
		Orders:      []*pgstorage.Order{pgstorage.NewOrder(column, direction)},
	}, nil
}
