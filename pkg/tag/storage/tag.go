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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package storage

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/tag/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/tag"
)

var (
	//go:embed sql/insert_tag.sql
	insertTagSQL string
	//go:embed sql/select_tags.sql
	listTagsSQL string
)

type TagStorage interface {
	UpsertTag(ctx context.Context, tag *domain.Tag) error
	ListTags(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Tag, int, int64, error)
}

type tagStorage struct {
	qe mysql.QueryExecer
}

func NewTagStorage(qe mysql.QueryExecer) TagStorage {
	return &tagStorage{qe: qe}
}

func (s *tagStorage) UpsertTag(ctx context.Context, tag *domain.Tag) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertTagSQL,
		tag.Id,
		tag.CreatedAt,
		tag.UpdatedAt,
		int32(tag.EntityType),
		tag.EnvironmentId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *tagStorage) ListTags(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.Tag, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(listTagsSQL, whereSQL, orderBySQL, limitOffsetSQL)

	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	tags := make([]*proto.Tag, 0, limit)
	for rows.Next() {
		var entityType int32
		tag := proto.Tag{}
		err := rows.Scan(
			&tag.Id,
			&tag.CreatedAt,
			&tag.UpdatedAt,
			&entityType,
			&tag.EnvironmentId,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		tag.EntityType = proto.Tag_EntityType(entityType)
		tags = append(tags, &tag)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(tags)
	countQuery := fmt.Sprintf(listTagsSQL, whereSQL)
	var totalCount int64
	if err := s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount); err != nil {
		return nil, 0, 0, err
	}
	return tags, nextOffset, totalCount, nil
}
