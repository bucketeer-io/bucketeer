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
	"errors"
	"strconv"

	mysqlstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/tag/domain"
	tagstorage "github.com/bucketeer-io/bucketeer/v2/pkg/tag/storage"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/tag"
)

var (
	//go:embed sql/insert_tag.sql
	insertTagSQL string
	//go:embed sql/select_tag.sql
	selectTagSQL string
	//go:embed sql/select_tag_by_name.sql
	selectTagByNameSQL string
	//go:embed sql/select_tags.sql
	selectTagsSQL string
	//go:embed sql/select_all_environment_tags.sql
	selectAllEnvironmentTagsSQL string
	//go:embed sql/count_tags.sql
	countTagsSQL string
	//go:embed sql/delete_tag.sql
	deleteTagSQL string
)

type tagStorage struct {
	qe mysqlstorage.QueryExecer
}

func NewTagStorage(qe mysqlstorage.QueryExecer) tagstorage.TagStorage {
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
		if errors.Is(err, mysqlstorage.ErrNoRows) {
			return nil, tagstorage.ErrTagNotFound
		}
		return nil, err
	}
	tag.EntityType = proto.Tag_EntityType(entityType)
	return &domain.Tag{Tag: &tag}, nil
}

func (s *tagStorage) GetTagByName(
	ctx context.Context,
	name, environmentId string,
	entityType proto.Tag_EntityType,
) (*domain.Tag, error) {
	var entityTypeInt int32
	tag := proto.Tag{}
	err := s.qe.QueryRowContext(
		ctx,
		selectTagByNameSQL,
		name,
		environmentId,
		int32(entityType),
	).Scan(
		&tag.Id,
		&tag.Name,
		&tag.CreatedAt,
		&tag.UpdatedAt,
		&entityTypeInt,
		&tag.EnvironmentId,
		&tag.EnvironmentName,
	)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrNoRows) {
			return nil, tagstorage.ErrTagNotFound
		}
		return nil, err
	}
	tag.EntityType = proto.Tag_EntityType(entityTypeInt)
	return &domain.Tag{Tag: &tag}, nil
}

func listTagsOrders(
	orderBy proto.ListTagsRequest_OrderBy,
	orderDirection proto.ListTagsRequest_OrderDirection,
) ([]*mysqlstorage.Order, error) {
	var column string
	switch orderBy {
	case proto.ListTagsRequest_DEFAULT,
		proto.ListTagsRequest_NAME:
		column = "tag.name"
	case proto.ListTagsRequest_CREATED_AT:
		column = "tag.created_at"
	case proto.ListTagsRequest_UPDATED_AT:
		column = "tag.updated_at"
	case proto.ListTagsRequest_ENTITY_TYPE:
		column = "tag.entity_type"
	case proto.ListTagsRequest_ENVIRONMENT:
		column = "env.name"
	default:
		return nil, tagstorage.ErrInvalidListTagsOrderBy
	}
	direction := mysqlstorage.OrderDirectionAsc
	if orderDirection == proto.ListTagsRequest_DESC {
		direction = mysqlstorage.OrderDirectionDesc
	}
	return []*mysqlstorage.Order{mysqlstorage.NewOrder(column, direction)}, nil
}

func (s *tagStorage) ListTags(
	ctx context.Context,
	p tagstorage.ListTagsParams,
) ([]*proto.Tag, int, int64, error) {
	orders, err := listTagsOrders(p.OrderBy, p.OrderDirection)
	if err != nil {
		return nil, 0, 0, err
	}
	filters := []*mysqlstorage.FilterV2{}
	if p.EnvironmentID != "" {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "tag.environment_id",
			Operator: mysqlstorage.OperatorEqual,
			Value:    p.EnvironmentID,
		})
	}
	if p.OrganizationID != "" {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "env.organization_id",
			Operator: mysqlstorage.OperatorEqual,
			Value:    p.OrganizationID,
		})
	}
	if p.EntityType != proto.Tag_UNSPECIFIED {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "tag.entity_type",
			Operator: mysqlstorage.OperatorEqual,
			Value:    p.EntityType,
		})
	}
	var inFilters []*mysqlstorage.InFilter
	if len(p.EnvironmentIDs) > 0 {
		values := make([]interface{}, 0, len(p.EnvironmentIDs))
		for _, id := range p.EnvironmentIDs {
			values = append(values, id)
		}
		inFilters = append(inFilters, &mysqlstorage.InFilter{
			Column: "tag.environment_id",
			Values: values,
		})
	}
	var searchQuery *mysqlstorage.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &mysqlstorage.SearchQuery{
			Columns: []string{"tag.name"},
			Keyword: p.SearchKeyword,
		}
	}
	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, 0, 0, tagstorage.ErrInvalidListTagsCursor
	}
	options := &mysqlstorage.ListOptions{
		Filters:     filters,
		InFilters:   inFilters,
		SearchQuery: searchQuery,
		Orders:      orders,
		Limit:       p.PageSize,
		Offset:      offset,
	}
	query, whereArgs := mysqlstorage.ConstructQueryAndWhereArgs(selectTagsSQL, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	tags := make([]*proto.Tag, 0, p.PageSize)
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
		return nil, 0, 0, rows.Err()
	}
	nextOffset := offset + len(tags)
	var totalCount int64
	countQuery, countWhereArgs := mysqlstorage.ConstructCountQuery(countTagsSQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
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
		return nil, rows.Err()
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
		return tagstorage.ErrTagUnexpectedAffectedRows
	}
	return nil
}
