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
	_ "embed"
	"fmt"

	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	mysqlstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var (
	//go:embed sql/feature_last_used_info/select_feature_last_used_infos.sql
	selectFeatureLastUsedInfosSQL string
	//go:embed sql/feature_last_used_info/upsert_feature_last_used_info.sql
	upsertFeatureLastUsedInfoSQL string
)

type featureLastUsedInfoStorage struct {
	qe mysqlstorage.QueryExecer
}

func NewFeatureLastUsedInfoStorage(qe mysqlstorage.QueryExecer) v2fs.FeatureLastUsedInfoStorage {
	return &featureLastUsedInfoStorage{qe: qe}
}

func (s *featureLastUsedInfoStorage) GetFeatureLastUsedInfos(
	ctx context.Context,
	ids []string,
	environmentId string,
) ([]*domain.FeatureLastUsedInfo, error) {
	inFilterIDs := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		inFilterIDs = append(inFilterIDs, id)
	}
	whereParts := []mysqlstorage.WherePart{
		mysqlstorage.NewInFilter("id", inFilterIDs),
		mysqlstorage.NewFilter("environment_id", "=", environmentId),
	}
	whereSQL, whereArgs := mysqlstorage.ConstructWhereSQLString(whereParts)
	query := fmt.Sprintf(selectFeatureLastUsedInfosSQL, whereSQL)
	rows, err := s.qe.QueryContext(
		ctx,
		query,
		whereArgs...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	entries := make([]*proto.FeatureLastUsedInfo, 0, len(ids))
	for rows.Next() {
		flui := proto.FeatureLastUsedInfo{}
		err := rows.Scan(
			&flui.FeatureId,
			&flui.Version,
			&flui.LastUsedAt,
			&flui.ClientOldestVersion,
			&flui.ClientLatestVersion,
			&flui.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &flui)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	domainFeatureLastUsedInfos := make([]*domain.FeatureLastUsedInfo, 0, len(entries))
	for _, e := range entries {
		domainFeatureLastUsedInfos = append(
			domainFeatureLastUsedInfos,
			&domain.FeatureLastUsedInfo{FeatureLastUsedInfo: e},
		)
	}
	return domainFeatureLastUsedInfos, nil
}

func (s *featureLastUsedInfoStorage) UpsertFeatureLastUsedInfo(
	ctx context.Context,
	flui *domain.FeatureLastUsedInfo,
	environmentId string,
) error {
	_, err := s.qe.ExecContext(
		ctx,
		upsertFeatureLastUsedInfoSQL,
		flui.ID(),
		flui.FeatureId,
		flui.Version,
		flui.LastUsedAt,
		flui.ClientOldestVersion,
		flui.ClientLatestVersion,
		flui.CreatedAt,
		environmentId,
	)
	if err != nil {
		return err
	}
	return nil
}
