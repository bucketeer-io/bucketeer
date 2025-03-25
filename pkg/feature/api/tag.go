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
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func (s *FeatureService) ListTags(
	ctx context.Context,
	req *featureproto.ListTagsRequest,
) (*featureproto.ListTagsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("environment_id", "=", req.EnvironmentId),
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"id"}, req.SearchKeyword))
	}
	orders, err := s.newListTagsOrdersMySQL(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
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
		dt, err := statusInvalidCursor.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	tags, nextCursor, totalCount, err := s.tagStorage.ListTags(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list tags",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentId, localizer)
	}
	return &featureproto.ListTagsResponse{
		Tags:       tags,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *FeatureService) newListTagsOrdersMySQL(
	orderBy featureproto.ListTagsRequest_OrderBy,
	orderDirection featureproto.ListTagsRequest_OrderDirection,
	localizer locale.Localizer,
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
		dt, err := statusInvalidOrderBy.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "order_by"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
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
		t, err := domain.NewTag(trimed)
		if err != nil {
			s.logger.Error(
				"Failed to create domain tag",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environment_id", environmentId),
					zap.String("tagID", tag),
				)...,
			)
			return err
		}
		if err := s.tagStorage.UpsertTag(ctx, t, environmentId); err != nil {
			s.logger.Error(
				"Failed to store tag",
				log.FieldsFromImcomingContext(ctx).AddFields(
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
