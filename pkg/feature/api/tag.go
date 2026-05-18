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

package api

import (
	"context"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/tag/domain"
	tagstorage "github.com/bucketeer-io/bucketeer/v2/pkg/tag/storage"
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
	orderBy, err := toTagOrderBy(req.OrderBy)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("orderBy", req.OrderBy.String()),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, statusInvalidOrderBy.Err()
	}
	orderDirection := tagproto.ListTagsRequest_ASC
	if req.OrderDirection == featureproto.ListTagsRequest_DESC {
		orderDirection = tagproto.ListTagsRequest_DESC
	}

	tags, nextCursor, totalCount, err := s.tagStorage.ListTags(ctx, tagstorage.ListTagsParams{
		EnvironmentID:  req.EnvironmentId,
		SearchKeyword:  req.SearchKeyword,
		OrderBy:        orderBy,
		OrderDirection: orderDirection,
		PageSize:       int(req.PageSize),
		Cursor:         req.Cursor,
	})
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

func toTagOrderBy(orderBy featureproto.ListTagsRequest_OrderBy) (tagproto.ListTagsRequest_OrderBy, error) {
	switch orderBy {
	case featureproto.ListTagsRequest_DEFAULT:
		return tagproto.ListTagsRequest_DEFAULT, nil
	case featureproto.ListTagsRequest_NAME:
		return tagproto.ListTagsRequest_NAME, nil
	case featureproto.ListTagsRequest_CREATED_AT:
		return tagproto.ListTagsRequest_CREATED_AT, nil
	case featureproto.ListTagsRequest_UPDATED_AT:
		return tagproto.ListTagsRequest_UPDATED_AT, nil
	default:
		return 0, statusInvalidOrderBy.Err()
	}
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
