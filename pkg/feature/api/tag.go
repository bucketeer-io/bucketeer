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

package api

import (
	"context"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/tag/domain"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	tagproto "github.com/bucketeer-io/bucketeer/v2/proto/tag"
)

func (s *FeatureService) ListTags(
	ctx context.Context,
	req *featureproto.ListTagsRequest,
) (*featureproto.ListTagsResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	filters := []*mysql.FilterV2{
		{
			Column:   "environment_id",
			Operator: mysql.OperatorEqual,
			Value:    req.EnvironmentId,
		},
	}
	var searchQuery *mysql.SearchQuery
	if req.SearchKeyword != "" {
		searchQuery = &mysql.SearchQuery{
			Columns: []string{"id"},
			Keyword: req.SearchKeyword,
		}
	}
	orders, err := s.newListTagsOrdersMySQL(req.OrderBy, req.OrderDirection)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	limit := int(req.PageSize)
	cursor := req.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, statusInvalidCursor.Err()
	}
	options := &mysql.ListOptions{
		Filters:     filters,
		SearchQuery: searchQuery,
		Orders:      orders,
		Limit:       limit,
		Offset:      offset,
	}
	tags, nextCursor, totalCount, err := s.tagStorage.ListTags(ctx, options)
	if err != nil {
		s.logger.Error(
			"Failed to list tags",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentId)
	}
	convertedTags := make([]*featureproto.Tag, len(tags))
	for i, tag := range tags {
		convertedTags[i] = &featureproto.Tag{
			Id:   tag.Id,
			Name: tag.Name,
		}
	}
	return &featureproto.ListTagsResponse{
		Tags:       convertedTags,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *FeatureService) newListTagsOrdersMySQL(
	orderBy featureproto.ListTagsRequest_OrderBy,
	orderDirection featureproto.ListTagsRequest_OrderDirection,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case featureproto.ListTagsRequest_DEFAULT,
		featureproto.ListTagsRequest_NAME:
		column = "tag.name"
	case featureproto.ListTagsRequest_CREATED_AT:
		column = "tag.created_at"
	case featureproto.ListTagsRequest_UPDATED_AT:
		column = "tag.updated_at"
	default:
		return nil, statusInvalidOrderBy.Err()
	}
	direction := mysql.OrderDirectionAsc
	if orderDirection == featureproto.ListTagsRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *FeatureService) upsertTags(
	ctx context.Context,
	tags []string,
	environmentId string,
) error {
	for _, tag := range tags {
		trimed := strings.TrimSpace(tag)
		if trimed == "" {
			continue
		}
		t, err := domain.NewTag(trimed, environmentId, tagproto.Tag_FEATURE_FLAG)
		if err != nil {
			s.logger.Error(
				"Failed to create domain tag",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environment_id", environmentId),
					zap.String("tagID", tag),
				)...,
			)
			return err
		}
		if err := s.tagStorage.UpsertTag(ctx, t); err != nil {
			s.logger.Error(
				"Failed to store tag",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environment_id", environmentId),
					zap.String("tagID", tag),
				)...,
			)
			return err
		}
	}
	return nil
}
