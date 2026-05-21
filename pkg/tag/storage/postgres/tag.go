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
	"fmt"
	"strconv"

	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
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
	qe pgstorage.QueryExecer
}

func NewTagStorage(qe pgstorage.QueryExecer) tagstorage.TagStorage {
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
		if errors.Is(err, pgstorage.ErrNoRows) {
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
		if errors.Is(err, pgstorage.ErrNoRows) {
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
) ([]*pgstorage.Order, error) {
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
	direction := pgstorage.OrderDirectionAsc
	if orderDirection == proto.ListTagsRequest_DESC {
		direction = pgstorage.OrderDirectionDesc
	}
	return []*pgstorage.Order{pgstorage.NewOrder(column, direction)}, nil
}

func (s *tagStorage) ListTags(
	ctx context.Context,
	p tagstorage.ListTagsParams,
) ([]*proto.Tag, int, int64, error) {
	orders, err := listTagsOrders(p.OrderBy, p.OrderDirection)
	if err != nil {
		return nil, 0, 0, err
	}
	filters := []*pgstorage.Filter{}
	if p.EnvironmentID != "" {
		filters = append(filters, &pgstorage.Filter{
			Column:   "tag.environment_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.EnvironmentID,
		})
	}
	if p.OrganizationID != "" {
		filters = append(filters, &pgstorage.Filter{
			Column:   "env.organization_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.OrganizationID,
		})
	}
	if p.EntityType != proto.Tag_UNSPECIFIED {
		filters = append(filters, &pgstorage.Filter{
			Column:   "tag.entity_type",
			Operator: pgstorage.OperatorEqual,
			Value:    int32(p.EntityType),
		})
	}
	var inFilters []*pgstorage.InFilter
	if len(p.EnvironmentIDs) > 0 {
		values := make([]interface{}, 0, len(p.EnvironmentIDs))
		for _, id := range p.EnvironmentIDs {
			values = append(values, id)
		}
		inFilters = append(inFilters, &pgstorage.InFilter{
			Column: "tag.environment_id",
			Values: values,
		})
	}
	var searchQuery *pgstorage.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &pgstorage.SearchQuery{
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
	options := &pgstorage.ListOptions{
		Limit:       p.PageSize,
		Offset:      offset,
		Filters:     filters,
		InFilters:   inFilters,
		SearchQuery: searchQuery,
		Orders:      orders,
	}
	whereParts := options.CreateWhereParts()
	whereSQL, whereArgs := pgstorage.ConstructWhereSQLString(whereParts)
	orderBySQL := pgstorage.ConstructOrderBySQLString(options.Orders)
	limitOffsetSQL := pgstorage.ConstructLimitOffsetSQLString(options.Limit, options.Offset)
	query := fmt.Sprintf(selectTagsSQL, whereSQL, orderBySQL, limitOffsetSQL)
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
	countQuery := fmt.Sprintf(countTagsSQL, whereSQL)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
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
