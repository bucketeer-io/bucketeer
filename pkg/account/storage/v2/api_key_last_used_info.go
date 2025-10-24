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

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

var (
	//go:embed sql/api_key_last_used_info/upsert_api_key_last_used_info.sql
	upsertAPIKeyLastUsedInfoSQLQuery string
	//go:embed sql/api_key_last_used_info/select_api_key_last_used_infos.sql
	selectAPIKeyLastUsedInfosSQLQuery string
)

func (s *accountStorage) UpsertAPIKeyLastUsedInfo(ctx context.Context, f *domain.APIKeyLastUsedInfo) error {
	_, err := s.qe.ExecContext(
		ctx,
		upsertAPIKeyLastUsedInfoSQLQuery,
		f.ApiKeyId,
		f.LastUsedAt,
		f.CreatedAt,
		f.EnvironmentId,
	)
	return err
}

func (s *accountStorage) GetAPIKeyLastUsedInfos(
	ctx context.Context,
	options *mysql.ListOptions,
) ([]*proto.APIKeyLastUsedInfo, error) {
	query, whereArgs := mysql.ConstructQueryAndWhereArgs(selectAPIKeyLastUsedInfosSQLQuery, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	infos := make([]*proto.APIKeyLastUsedInfo, 0)
	for rows.Next() {
		info := proto.APIKeyLastUsedInfo{}
		err := rows.Scan(
			&info.ApiKeyId,
			&info.LastUsedAt,
			&info.CreatedAt,
			&info.EnvironmentId,
		)
		if err != nil {
			return nil, err
		}
		infos = append(infos, &info)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return infos, nil
}
