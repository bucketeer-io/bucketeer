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
	"errors"
	"slices"
	"strconv"
	"strings"

	pb "github.com/golang/protobuf/proto" // nolint:staticcheck
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	ftstorage "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/tag/domain"
	tagstorage "github.com/bucketeer-io/bucketeer/v2/pkg/tag/storage"
	accproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	"github.com/bucketeer-io/bucketeer/v2/proto/feature"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/tag"
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
	mysqlClient    mysql.Client
	tagStorage     tagstorage.TagStorage
	featureStorage ftstorage.FeatureStorage
	accountClient  accclient.Client
	publisher      publisher.Publisher
	opts           *options
	logger         *zap.Logger
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
		mysqlClient:    mysqlClient,
		tagStorage:     tagstorage.NewTagStorage(mysqlClient),
		featureStorage: ftstorage.NewFeatureStorage(mysqlClient),
		accountClient:  accountClient,
		publisher:      publisher,
		opts:           dopts,
		logger:         dopts.logger.Named("api"),
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("name", req.Name),
				zap.Int32("entity_type", int32(req.EntityType)),
			)...,
		)
		return nil, err
	}
	tag, err := domain.NewTag(strings.TrimSpace(req.Name), req.EnvironmentId, req.EntityType)
	if err != nil {
		s.logger.Error(
			"Failed to create domain tag",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.Any("tag", tag),
			)...,
		)
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentId, localizer)
	}

	var event *eventproto.Event
	var actualTag *domain.Tag
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		// Check if tag exists before upsert to determine if it's create or update
		existingTag, err := s.tagStorage.GetTagByName(ctxWithTx, tag.Name, tag.EnvironmentId, tag.EntityType)
		isCreate := err != nil && errors.Is(err, tagstorage.ErrTagNotFound)
		if err != nil && !isCreate {
			return err
		}

		if err := s.tagStorage.UpsertTag(ctxWithTx, tag); err != nil {
			return err
		}

		// Fetch the actual tag from DB to get correct final state after upsert
		actualTag, err = s.tagStorage.GetTagByName(ctxWithTx, tag.Name, tag.EnvironmentId, tag.EntityType)
		if err != nil {
			return err
		}

		// Use appropriate event type with appropriate previous entity data based on whether it's create or update
		var previousEntityData interface{}
		var eventType eventproto.Event_Type
		var eventData pb.Message
		if isCreate {
			previousEntityData = nil
			eventType = eventproto.Event_TAG_CREATED
			eventData = &eventproto.TagCreatedEvent{
				Id:            actualTag.Id,
				Name:          actualTag.Name,
				CreatedAt:     actualTag.CreatedAt,
				UpdatedAt:     actualTag.UpdatedAt,
				EntityType:    actualTag.EntityType,
				EnvironmentId: actualTag.EnvironmentId,
			}
		} else {
			previousEntityData = existingTag
			eventType = eventproto.Event_TAG_UPDATED
			eventData = &eventproto.TagUpdatedEvent{
				Id:            actualTag.Id,
				Name:          actualTag.Name,
				UpdatedAt:     actualTag.UpdatedAt,
				EntityType:    actualTag.EntityType,
				EnvironmentId: actualTag.EnvironmentId,
			}
		}
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_TAG,
			actualTag.Id,
			eventType,
			eventData,
			actualTag.EnvironmentId,
			actualTag,
			previousEntityData,
		)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.logger.Error(
			"Failed to store the tag",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.Any("tag", tag),
			)...,
		)
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentId, localizer)
	}
	if err := s.publisher.Publish(ctx, event); err != nil {
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentId, localizer)
	}
	return &proto.CreateTagResponse{Tag: actualTag.Tag}, nil
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
	inFilters := make([]*mysql.InFilter, 0)
	filters := []*mysql.FilterV2{}
	if req.OrganizationId != "" {
		// New console
		editor, err := s.checkOrganizationRole(
			ctx, accproto.AccountV2_Role_Organization_MEMBER,
			req.OrganizationId, localizer)
		if err != nil {
			return nil, err
		}
		filters = append(filters, &mysql.FilterV2{
			Column:   "env.organization_id",
			Operator: mysql.OperatorEqual,
			Value:    req.OrganizationId,
		})
		filterEnvironmentIDs := s.getAllowedEnvironments([]string{req.EnvironmentId}, editor)
		values := make([]interface{}, 0)
		for _, id := range filterEnvironmentIDs {
			values = append(values, id)
		}
		if len(filterEnvironmentIDs) > 0 {
			inFilters = append(inFilters, &mysql.InFilter{
				Column: "tag.environment_id",
				Values: values,
			})
		}
	} else {
		// Current console
		_, err := s.checkEnvironmentRole(
			ctx, accproto.AccountV2_Role_Environment_VIEWER,
			req.EnvironmentId, localizer)
		if err != nil {
			return nil, err
		}
		filters = append(filters, &mysql.FilterV2{
			Column:   "tag.environment_id",
			Operator: mysql.OperatorEqual,
			Value:    req.EnvironmentId,
		})
	}
	var searchQuery *mysql.SearchQuery
	if req.SearchKeyword != "" {
		searchQuery = &mysql.SearchQuery{
			Columns: []string{"tag.name"},
			Keyword: req.SearchKeyword,
		}
	}
	if req.EntityType != proto.Tag_UNSPECIFIED {
		filters = append(filters, &mysql.FilterV2{
			Column:   "tag.entity_type",
			Operator: mysql.OperatorEqual,
			Value:    req.EntityType,
		})
	}
	orders, err := s.newListTagsOrdersMySQL(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Failed to valid list tags API. Invalid argument.",
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
		dt, err := statusInvalidCursor.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	options := &mysql.ListOptions{
		Filters:     filters,
		SearchQuery: searchQuery,
		Orders:      orders,
		Limit:       limit,
		Offset:      offset,
		JSONFilters: nil,
		InFilters:   inFilters,
		NullFilters: nil,
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}
	var tag *domain.Tag
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		tagDB, err := s.tagStorage.GetTag(ctxWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			s.logger.Error(
				"Failed to get tag",
				log.FieldsFromIncomingContext(ctxWithTx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
					zap.Any("id", req.Id),
				)...,
			)
			return err
		}
		tag = tagDB

		// Check if the tag is in use by any feature
		features, err := s.listFeaturesFromEnvironment(ctxWithTx, req.EnvironmentId)
		if err != nil {
			return err
		}
		var inUsed = false
		for _, f := range features {
			if slices.Contains(f.Tags, tagDB.Name) {
				inUsed = true
				break
			}
		}
		if inUsed {
			s.logger.Error(
				"Failed to delete the tag because it is in use by a feature",
				log.FieldsFromIncomingContext(ctxWithTx).AddFields(
					zap.String("environmentId", req.EnvironmentId),
					zap.Any("tag", tagDB),
				)...,
			)
			return statusTagInUsed.Err()
		}

		// Delete it from DB
		return s.tagStorage.DeleteTag(ctxWithTx, req.Id)
	})
	if err != nil {
		if errors.Is(err, statusTagInUsed.Err()) {
			dt, err := statusTagInUsed.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.Tag),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
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
		nil, // Current state: entity no longer exists
		tag, // Previous state: what was deleted
	)
	if err != nil {
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentId, localizer)
	}
	if err = s.publisher.Publish(ctx, event); err != nil {
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentId, localizer)
	}
	return &proto.DeleteTagResponse{}, nil
}

func (s *TagService) listFeaturesFromEnvironment(
	ctx context.Context,
	environmentID string,
) ([]*feature.Feature, error) {
	features, _, _, err := s.featureStorage.ListFeatures(ctx, &mysql.ListOptions{
		Filters: []*mysql.FilterV2{
			{
				Column:   "feature.environment_id",
				Operator: mysql.OperatorEqual,
				Value:    environmentID,
			},
		},
		Orders:      nil,
		JSONFilters: nil,
		NullFilters: nil,
		InFilters:   nil,
		SearchQuery: nil,
		Limit:       mysql.QueryNoLimit,
		Offset:      mysql.QueryNoOffset,
	})
	if err != nil {
		s.logger.Error(
			"Failed to list features",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentID),
			)...,
		)
		return nil, err
	}
	return features, nil
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
				log.FieldsFromIncomingContext(ctx).AddFields(
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
				log.FieldsFromIncomingContext(ctx).AddFields(
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
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return nil, api.NewGRPCStatus(err).Err()
		}
	}
	return editor, nil
}

func (s *TagService) checkOrganizationRole(
	ctx context.Context,
	requiredRole accproto.AccountV2_Role_Organization,
	organizationID string,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckOrganizationRole(ctx, requiredRole, func(
		email string,
	) (*accproto.GetAccountV2Response, error) {
		resp, err := s.accountClient.GetAccountV2(ctx, &accproto.GetAccountV2Request{
			Email:          email,
			OrganizationId: organizationID,
		})
		if err != nil {
			return nil, err
		}
		return resp, nil
	})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationID", organizationID),
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
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationID", organizationID),
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
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationID", organizationID),
				)...,
			)
			return nil, api.NewGRPCStatus(err).Err()
		}
	}
	return editor, nil
}

func (s *TagService) getAllowedEnvironments(
	reqEnvironmentIDs []string,
	editor *eventproto.Editor,
) []string {
	filterEnvironmentIDs := make([]string, 0)
	if editor.OrganizationRole == accproto.AccountV2_Role_Organization_MEMBER {
		// only show API keys in allowed environments for member.
		if len(reqEnvironmentIDs) > 0 {
			for _, id := range reqEnvironmentIDs {
				for _, e := range editor.EnvironmentRoles {
					if e.EnvironmentId == id {
						filterEnvironmentIDs = append(filterEnvironmentIDs, id)
						break
					}
				}
			}
		} else {
			for _, e := range editor.EnvironmentRoles {
				filterEnvironmentIDs = append(filterEnvironmentIDs, e.EnvironmentId)
			}
		}
	} else {
		// if the user is an admin or owner, no need to filter environments.
		filterEnvironmentIDs = append(filterEnvironmentIDs, reqEnvironmentIDs...)
	}
	return filterEnvironmentIDs
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
		log.FieldsFromIncomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentId", environmentId),
		)...,
	)
	st := api.NewGRPCStatus(err)
	dt, err := st.WithDetails(&errdetails.LocalizedMessage{
		Locale:  localizer.GetLocale(),
		Message: localizer.MustLocalize(locale.InternalServerError),
	})
	if err != nil {
		return statusInternal.Err()
	}
	return dt.Err()
}
