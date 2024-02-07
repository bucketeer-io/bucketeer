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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v2

import (
	"context"
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type FeatureLastUsedInfoStorage interface {
	GetFeatureLastUsedInfos(
		ctx context.Context,
		ids []string,
		environmentNamespace string,
	) ([]*domain.FeatureLastUsedInfo, error)
	UpsertFeatureLastUsedInfo(
		ctx context.Context,
		featureLastUsedInfos *domain.FeatureLastUsedInfo,
		environmentNamespace string,
	) error
}

type featureLastUsedInfoStorage struct {
	qe mysql.QueryExecer
}

func NewFeatureLastUsedInfoStorage(qe mysql.QueryExecer) FeatureLastUsedInfoStorage {
	return &featureLastUsedInfoStorage{qe: qe}
}

func (s *featureLastUsedInfoStorage) GetFeatureLastUsedInfos(
	ctx context.Context,
	ids []string,
	environmentNamespace string,
) ([]*domain.FeatureLastUsedInfo, error) {
	inFilterIDs := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		inFilterIDs = append(inFilterIDs, id)
	}
	whereParts := []mysql.WherePart{
		mysql.NewInFilter("id", inFilterIDs),
		mysql.NewFilter("environment_namespace", "=", environmentNamespace),
	}
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	query := fmt.Sprintf(`
		SELECT
			feature_id,
			version,
			last_used_at,
			client_oldest_version,
			client_latest_version,
			created_at
		FROM
			feature_last_used_info
		%s
	`, whereSQL,
	)
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
		return nil, err
	}
	// NOTE: If the performance matters, remove the following loop and return protos.
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
	environmentNamespace string,
) error {
	query := `
		INSERT INTO feature_last_used_info (
			id,
			feature_id,
			version,
			last_used_at,
			client_oldest_version,
			client_latest_version,
			created_at,
			environment_namespace
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?
		) ON DUPLICATE KEY UPDATE
			feature_id = VALUES(feature_id),
			version = VALUES(version),
			last_used_at = VALUES(last_used_at),
			client_oldest_version = VALUES(client_oldest_version),
			client_latest_version = VALUES(client_latest_version)
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
		flui.ID(),
		flui.FeatureId,
		flui.Version,
		flui.LastUsedAt,
		flui.ClientOldestVersion,
		flui.ClientLatestVersion,
		flui.CreatedAt,
		environmentNamespace,
	)
	if err != nil {
		return err
	}
	return nil
}
