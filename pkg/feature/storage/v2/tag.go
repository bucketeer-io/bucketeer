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
	"errors"
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	ErrTagAlreadyExists = errors.New("tag: already exists")
	ErrTagNotFound      = errors.New("tag: not found")
)

type TagStorage interface {
	UpsertTag(ctx context.Context, tag *domain.Tag, environmentNamespace string) error
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

func (s *tagStorage) UpsertTag(
	ctx context.Context,
	tag *domain.Tag,
	environmentNamespace string,
) error {
	// To get last tags, update `updated_at`.
	query := `
		INSERT INTO tag (
			id,
			created_at,
			updated_at,
			environment_namespace
		) VALUES (
			?, ?, ?, ?
		) ON DUPLICATE KEY UPDATE
			updated_at = VALUES(updated_at)
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
		tag.Id,
		tag.CreatedAt,
		tag.UpdatedAt,
		environmentNamespace,
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
	query := fmt.Sprintf(`
		SELECT
			id,
			created_at,
			updated_at
		FROM
			tag
		%s %s %s
		`, whereSQL, orderBySQL, limitOffsetSQL,
	)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	tags := make([]*proto.Tag, 0, limit)
	for rows.Next() {
		tag := proto.Tag{}
		err := rows.Scan(
			&tag.Id,
			&tag.CreatedAt,
			&tag.UpdatedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		tags = append(tags, &tag)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(tags)
	countQuery := fmt.Sprintf(`
		SELECT
			COUNT(1)
		FROM
			tag
		%s %s
		`, whereSQL, orderBySQL,
	)
	var totalCount int64
	if err := s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount); err != nil {
		return nil, 0, 0, err
	}
	return tags, nextOffset, totalCount, nil
}
