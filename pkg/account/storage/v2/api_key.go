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

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

var (
	ErrAPIKeyAlreadyExists          = pkgErr.NewErrorAlreadyExists(pkgErr.AccountPackageName, "api key already exists")
	ErrAPIKeyNotFound               = pkgErr.NewErrorNotFound(pkgErr.AccountPackageName, "api key not found", "api_key")
	ErrAPIKeyUnexpectedAffectedRows = pkgErr.NewErrorUnexpectedAffectedRows(
		pkgErr.AccountPackageName,
		"api key unexpected affected rows",
	)
)

var (
	//go:embed sql/api_key_v2/insert_api_key_v2.sql
	insertAPIKeyV2SQLQuery string
	//go:embed sql/api_key_v2/update_api_key_v2.sql
	updateAPIKeyV2SQLQuery string
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
		if errors.Is(err, mysql.ErrDuplicateEntry) {
			return ErrAPIKeyAlreadyExists
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
		return ErrAPIKeyUnexpectedAffectedRows
	}
	return nil
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
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrAPIKeyNotFound
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
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrAPIKeyNotFound
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
		return nil, err
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
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrAPIKeyNotFound
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
	options *mysql.ListOptions,
) ([]*proto.APIKey, int, int64, error) {
	query, whereArgs := mysql.ConstructQueryAndWhereArgs(selectAPIKeyV2SQLQuery, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	var limit, offset int
	if options.Limit != 0 {
		limit = options.Limit
		offset = options.Offset
	}
	apiKeys := make([]*proto.APIKey, 0, limit)
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
		)
		if err != nil {
			return nil, 0, 0, err
		}
		apiKey.Role = proto.APIKey_Role(role)
		apiKeys = append(apiKeys, &apiKey)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(apiKeys)
	var totalCount int64
	countQuery, countWhereArgs := mysql.ConstructCountQuery(selectAPIKeyV2CountSQLQuery, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return apiKeys, nextOffset, totalCount, nil
}
