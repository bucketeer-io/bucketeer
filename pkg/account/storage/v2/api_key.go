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
)

var (
	ErrAPIKeyAlreadyExists          = errors.New("apiKey: api key already exists")
	ErrAPIKeyNotFound               = errors.New("apiKey: api key not found")
	ErrAPIKeyUnexpectedAffectedRows = errors.New("apiKey: api key unexpected affected rows")
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
	//go:embed sql/api_key_v2/select_api_key_v2_by_token.sql
	selectAPIKeyV2ByTokenSQLQuery string
)

func (s *accountStorage) CreateAPIKey(ctx context.Context, k *domain.APIKey, environmentID string) error {
	_, err := s.qe(ctx).ExecContext(
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
	result, err := s.qe(ctx).ExecContext(
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

func (s *accountStorage) GetAPIKey(ctx context.Context, token, environmentID string) (*domain.APIKey, error) {
	apiKey := proto.APIKey{}
	var role int32
	err := s.qe(ctx).QueryRowContext(
		ctx,
		selectAPIKeyV2ByTokenSQLQuery,
		token,
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

func (s *accountStorage) ListAPIKeys(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.APIKey, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(selectAPIKeyV2SQLQuery, whereSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.qe(ctx).QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
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
	countQuery := fmt.Sprintf(selectAPIKeyV2CountSQLQuery, whereSQL)
	err = s.qe(ctx).QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return apiKeys, nextOffset, totalCount, nil
}
