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
	"errors"

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/tag/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/tag"
)

var (
	ErrTagNotFound               = errors.New("tag: not found")
	ErrTagUnexpectedAffectedRows = errors.New("tag: unexpected affected rows")

	//go:embed sql/insert_tag.sql
	insertTagSQL string
	//go:embed sql/select_tag.sql
	selectTagSQL string
	//go:embed sql/select_tags.sql
	selectTagsSQL string
	//go:embed sql/select_all_environment_tags.sql
	selectAllEnvironmentTagsSQL string
	//go:embed sql/count_tags.sql
	countTagsSQL string
	//go:embed sql/delete_tag.sql
	deleteTagSQL string
)

type TagStorage interface {
	UpsertTag(ctx context.Context, tag *domain.Tag) error
	GetTag(ctx context.Context, id, environmentId string) (*domain.Tag, error)
	ListTags(
		ctx context.Context,
		options *mysql.ListOptions,
	) ([]*proto.Tag, int, int64, error)
	ListAllEnvironmentTags(ctx context.Context) ([]*proto.EnvironmentTag, error)
	DeleteTag(ctx context.Context, id string) error
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
		&tag.Name,
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

func (s *tagStorage) GetTag(ctx context.Context, id, environmentId string) (*domain.Tag, error) {
	var entityType int32
	tag := proto.Tag{}
	err := s.qe.QueryRowContext(
		ctx,
		selectTagSQL,
		id,
		environmentId,
	).Scan(
		&tag.Id,
		&tag.Name,
		&tag.CreatedAt,
		&tag.UpdatedAt,
		&entityType,
		&tag.EnvironmentId,
		&tag.EnvironmentName,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}
	tag.EntityType = proto.Tag_EntityType(entityType)
	return &domain.Tag{Tag: &tag}, nil
}

func (s *tagStorage) ListTags(
	ctx context.Context,
	options *mysql.ListOptions,
) ([]*proto.Tag, int, int64, error) {
	query, whereArgs := mysql.ConstructQueryAndWhereArgs(selectTagsSQL, options)

	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	var limit, offset int
	if options != nil {
		limit = options.Limit
		offset = options.Offset
	}
	tags := make([]*proto.Tag, 0, limit)
	for rows.Next() {
		var entityType int32
		tag := proto.Tag{}
		err := rows.Scan(
			&tag.Id,
			&tag.Name,
			&tag.CreatedAt,
			&tag.UpdatedAt,
			&entityType,
			&tag.EnvironmentId,
			&tag.EnvironmentName,
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
	var totalCount int64
	countQuery, countWhereArgs := mysql.ConstructCountQuery(countTagsSQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return tags, nextOffset, totalCount, nil
}

func (s *tagStorage) ListAllEnvironmentTags(ctx context.Context) ([]*proto.EnvironmentTag, error) {
	rows, err := s.qe.QueryContext(ctx, selectAllEnvironmentTagsSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	envTags := map[string][]*proto.Tag{}
	for rows.Next() {
		var entityType int32
		var envID string
		tag := proto.Tag{}
		err := rows.Scan(
			&envID,
			&tag.Id,
			&tag.Name,
			&tag.CreatedAt,
			&tag.UpdatedAt,
			&entityType,
			&tag.EnvironmentId,
			&tag.EnvironmentName,
		)
		if err != nil {
			return nil, err
		}
		tag.EntityType = proto.Tag_EntityType(entityType)
		envTags[envID] = append(envTags[envID], &tag)
	}
	if rows.Err() != nil {
		return nil, err
	}
	environmentTags := make([]*proto.EnvironmentTag, 0, len(envTags))
	for key, tags := range envTags {
		envTag := &proto.EnvironmentTag{
			EnvironmentId: key,
			Tags:          tags,
		}
		environmentTags = append(environmentTags, envTag)
	}
	return environmentTags, nil
}

func (s *tagStorage) DeleteTag(ctx context.Context, id string) error {
	result, err := s.qe.ExecContext(ctx, deleteTagSQL, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrTagUnexpectedAffectedRows
	}
	return nil
}
