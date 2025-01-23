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
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/role"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/tag/domain"
	tagstorage "github.com/bucketeer-io/bucketeer/pkg/tag/storage"
	accproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/tag"
)

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type TagService struct {
	mysqlClient   mysql.Client
	accountClient accclient.Client
	publisher     publisher.Publisher
	opts          *options
	logger        *zap.Logger
}

func NewTagService(
	mysqlClient mysql.Client,
	accountClient accclient.Client,
	publisher publisher.Publisher,
	opts ...Option,
) *TagService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &TagService{
		mysqlClient:   mysqlClient,
		accountClient: accountClient,
		publisher:     publisher,
		opts:          dopts,
		logger:        dopts.logger.Named("api"),
	}
}

func (s *TagService) Register(server *grpc.Server) {
	proto.RegisterTagServiceServer(server, s)
}

func (s *TagService) CreateTag(
	ctx context.Context,
	req *proto.CreateTagRequest,
) (*proto.CreateTagResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateCreateTagRquest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to create a tag",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("name", req.Name),
				zap.Int32("entity_type", int32(req.EntityType)),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	tagStorage := tagstorage.NewTagStorage(s.mysqlClient)
	tag, err := domain.NewTag(strings.TrimSpace(req.Name), req.EnvironmentId, req.EntityType)
	if err != nil {
		s.logger.Error(
			"Failed to create domain tag",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.Any("tag", tag),
			)...,
		)
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentId, localizer)
	}
	// Save in the DB
	if err := tagStorage.UpsertTag(ctx, tag); err != nil {
		s.logger.Error(
			"Failed to store the tag",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.Any("tag", tag),
			)...,
		)
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentId, localizer)
	}
	// Publish event
	event, err := domainevent.NewEvent(
		editor,
		eventproto.Event_TAG,
		tag.Id,
		eventproto.Event_TAG_CREATED,
		&eventproto.TagCreatedEvent{
			Id:            tag.Id,
			Name:          tag.Name,
			CreatedAt:     tag.CreatedAt,
			UpdatedAt:     tag.UpdatedAt,
			EntityType:    tag.EntityType,
			EnvironmentId: tag.EnvironmentId,
		},
		tag.EnvironmentId,
		tag,
		tag,
	)
	if err != nil {
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentId, localizer)
	}
	if err = s.publisher.Publish(ctx, event); err != nil {
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentId, localizer)
	}
	return &proto.CreateTagResponse{Tag: tag.Tag}, nil
}

func (s *TagService) validateCreateTagRquest(req *proto.CreateTagRequest, localizer locale.Localizer) error {
	if len(strings.TrimSpace(req.Name)) == 0 {
		dt, err := statusNameRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.EntityType == proto.Tag_UNSPECIFIED {
		dt, err := statusEntityTypeRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "entity_type"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *TagService) ListTags(
	ctx context.Context,
	req *proto.ListTagsRequest,
) (*proto.ListTagsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{}
	if req.OrganizationId != "" {
		// New console
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"tag.organization_id"}, req.OrganizationId))
	} else {
		// Current console
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"tag.environment_id"}, req.EnvironmentId))
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"tag.name"}, req.SearchKeyword))
	}
	if req.EntityType != proto.Tag_UNSPECIFIED {
		whereParts = append(whereParts, mysql.NewFilter("tag.entity_type", "=", req.EntityType))
	}
	orders, err := s.newListTagsOrdersMySQL(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Failed to valid list tags API. Invalid argument.",
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
	tagStorage := tagstorage.NewTagStorage(s.mysqlClient)
	tags, nextCursor, totalCount, err := tagStorage.ListTags(
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
	return &proto.ListTagsResponse{
		Tags:       tags,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *TagService) DeleteTag(
	ctx context.Context,
	req *proto.DeleteTagRequest,
) (*proto.DeleteTagResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateDeleteTagRquest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to delete a tag",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	tagStorage := tagstorage.NewTagStorage(s.mysqlClient)
	tag, err := tagStorage.GetTag(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentId, localizer)
	}
	// Delete it from DB
	if err := tagStorage.DeleteTag(ctx, req.Id, req.EnvironmentId); err != nil {
		s.logger.Error(
			"Failed to store the tag",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.Any("tag", tag),
			)...,
		)
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentId, localizer)
	}
	// Publish event
	event, err := domainevent.NewEvent(
		editor,
		eventproto.Event_TAG,
		tag.Id,
		eventproto.Event_TAG_DELETED,
		&eventproto.TagDeletedEvent{
			Id:            tag.Id,
			EnvironmentId: tag.EnvironmentId,
		},
		tag.EnvironmentId,
		tag,
		tag,
	)
	if err != nil {
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentId, localizer)
	}
	if err = s.publisher.Publish(ctx, event); err != nil {
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentId, localizer)
	}
	return &proto.DeleteTagResponse{}, nil
}

func (s *TagService) validateDeleteTagRquest(req *proto.DeleteTagRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusNameRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *TagService) checkEnvironmentRole(
	ctx context.Context,
	requiredRole accproto.AccountV2_Role_Environment,
	environmentId string,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckEnvironmentRole(
		ctx,
		requiredRole,
		environmentId,
		func(email string) (*accproto.AccountV2, error) {
			resp, err := s.accountClient.GetAccountV2ByEnvironmentID(ctx, &accproto.GetAccountV2ByEnvironmentIDRequest{
				Email:         email,
				EnvironmentId: environmentId,
			})
			if err != nil {
				return nil, err
			}
			return resp.Account, nil
		})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.UnauthenticatedError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			dt, err := statusPermissionDenied.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.PermissionDenied),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
	}
	return editor, nil
}

func (s *TagService) newListTagsOrdersMySQL(
	orderBy proto.ListTagsRequest_OrderBy,
	orderDirection proto.ListTagsRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
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
	if orderDirection == proto.ListTagsRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *TagService) reportInternalServerError(
	ctx context.Context,
	err error,
	environmentId string,
	localizer locale.Localizer,
) error {
	s.logger.Error(
		"Internal server error",
		log.FieldsFromImcomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentId", environmentId),
		)...,
	)
	dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
		Locale:  localizer.GetLocale(),
		Message: localizer.MustLocalize(locale.InternalServerError),
	})
	if err != nil {
		return statusInternal.Err()
	}
	return dt.Err()
}
