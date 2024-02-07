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

package api

import (
	"context"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	v2als "github.com/bucketeer-io/bucketeer/pkg/auditlog/storage/v2"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/role"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	proto "github.com/bucketeer-io/bucketeer/proto/auditlog"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
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

type AuditlogService interface {
	Register(*grpc.Server)
	ListAuditLogs(context.Context, *proto.ListAuditLogsRequest) (*proto.ListAuditLogsResponse, error)
	ListAdminAuditLogs(
		ctx context.Context,
		req *proto.ListAdminAuditLogsRequest,
	) (*proto.ListAdminAuditLogsResponse, error)
	ListFeatureHistory(
		ctx context.Context,
		req *proto.ListFeatureHistoryRequest,
	) (*proto.ListFeatureHistoryResponse, error)
}

type auditlogService struct {
	accountClient     accountclient.Client
	mysqlStorage      v2als.AuditLogStorage
	mysqlAdminStorage v2als.AdminAuditLogStorage
	opts              *options
	logger            *zap.Logger
}

func NewAuditLogService(
	accountClient accountclient.Client,
	mysqlClient mysql.Client,
	opts ...Option,
) AuditlogService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &auditlogService{
		accountClient:     accountClient,
		mysqlStorage:      v2als.NewAuditLogStorage(mysqlClient),
		mysqlAdminStorage: v2als.NewAdminAuditLogStorage(mysqlClient),
		opts:              dopts,
		logger:            dopts.logger.Named("api"),
	}
}

func (s *auditlogService) Register(server *grpc.Server) {
	proto.RegisterAuditLogServiceServer(server, s)
}

func (s *auditlogService) ListAuditLogs(
	ctx context.Context,
	req *proto.ListAuditLogsRequest,
) (*proto.ListAuditLogsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentNamespace, localizer)
	if err != nil {
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
	whereParts := []mysql.WherePart{
		mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace),
	}
	if req.From != 0 {
		whereParts = append(whereParts, mysql.NewFilter("timestamp", ">=", req.From))
	}
	if req.To != 0 {
		whereParts = append(whereParts, mysql.NewFilter("timestamp", "<=", req.To))
	}
	if req.EntityType != nil {
		whereParts = append(whereParts, mysql.NewFilter("entity_type", "=", req.EntityType.Value))
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"editor"}, req.SearchKeyword))
	}
	orders, err := s.newAuditLogListOrders(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	auditlogs, nextCursor, totalCount, err := s.mysqlStorage.ListAuditLogs(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list auditlogs",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
	for _, auditlog := range auditlogs {
		auditlog.LocalizedMessage = domainevent.LocalizedMessage(auditlog.Type, localizer)
	}
	return &proto.ListAuditLogsResponse{
		AuditLogs:  auditlogs,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *auditlogService) newAuditLogListOrders(
	orderBy proto.ListAuditLogsRequest_OrderBy,
	orderDirection proto.ListAuditLogsRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case proto.ListAuditLogsRequest_DEFAULT,
		proto.ListAuditLogsRequest_TIMESTAMP:
		column = "timestamp"
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
	direction := mysql.OrderDirectionDesc
	if orderDirection == proto.ListAuditLogsRequest_ASC {
		direction = mysql.OrderDirectionAsc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *auditlogService) ListAdminAuditLogs(
	ctx context.Context,
	req *proto.ListAdminAuditLogsRequest,
) (*proto.ListAdminAuditLogsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{}
	if req.From != 0 {
		whereParts = append(whereParts, mysql.NewFilter("timestamp", ">=", req.From))
	}
	if req.To != 0 {
		whereParts = append(whereParts, mysql.NewFilter("timestamp", "<=", req.To))
	}
	if req.EntityType != nil {
		whereParts = append(whereParts, mysql.NewFilter("entity_type", "=", req.EntityType.Value))
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"editor"}, req.SearchKeyword))
	}
	orders, err := s.newAdminAuditLogListOrders(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
	auditlogs, nextCursor, totalCount, err := s.mysqlAdminStorage.ListAdminAuditLogs(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list admin auditlogs",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
	for _, auditlog := range auditlogs {
		auditlog.LocalizedMessage = domainevent.LocalizedMessage(auditlog.Type, localizer)
	}
	return &proto.ListAdminAuditLogsResponse{
		AuditLogs:  auditlogs,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *auditlogService) newAdminAuditLogListOrders(
	orderBy proto.ListAdminAuditLogsRequest_OrderBy,
	orderDirection proto.ListAdminAuditLogsRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case proto.ListAdminAuditLogsRequest_DEFAULT,
		proto.ListAdminAuditLogsRequest_TIMESTAMP:
		column = "timestamp"
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
	direction := mysql.OrderDirectionDesc
	if orderDirection == proto.ListAdminAuditLogsRequest_ASC {
		direction = mysql.OrderDirectionAsc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *auditlogService) ListFeatureHistory(
	ctx context.Context,
	req *proto.ListFeatureHistoryRequest,
) (*proto.ListFeatureHistoryResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentNamespace,
		localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace),
		mysql.NewFilter("entity_type", "=", int32(eventproto.Event_FEATURE)),
		mysql.NewFilter("entity_id", "=", req.FeatureId),
	}
	if req.From != 0 {
		whereParts = append(whereParts, mysql.NewFilter("timestamp", ">=", req.From))
	}
	if req.To != 0 {
		whereParts = append(whereParts, mysql.NewFilter("timestamp", "<=", req.To))
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"editor"}, req.SearchKeyword))
	}
	orders, err := s.newFeatureHistoryAuditLogListOrders(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
	auditlogs, nextCursor, totalCount, err := s.mysqlStorage.ListAuditLogs(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list feature history",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.String("featureId", req.FeatureId),
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
	for _, auditlog := range auditlogs {
		auditlog.LocalizedMessage = domainevent.LocalizedMessage(auditlog.Type, localizer)
	}
	return &proto.ListFeatureHistoryResponse{
		AuditLogs:  auditlogs,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *auditlogService) newFeatureHistoryAuditLogListOrders(
	orderBy proto.ListFeatureHistoryRequest_OrderBy,
	orderDirection proto.ListFeatureHistoryRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case proto.ListFeatureHistoryRequest_DEFAULT,
		proto.ListFeatureHistoryRequest_TIMESTAMP:
		column = "timestamp"
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
	direction := mysql.OrderDirectionDesc
	if orderDirection == proto.ListFeatureHistoryRequest_ASC {
		direction = mysql.OrderDirectionAsc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *auditlogService) checkEnvironmentRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Environment,
	environmentNamespace string,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckEnvironmentRole(
		ctx,
		requiredRole,
		environmentNamespace,
		func(email string) (*accountproto.AccountV2, error) {
			resp, err := s.accountClient.GetAccountV2ByEnvironmentID(ctx, &accountproto.GetAccountV2ByEnvironmentIDRequest{
				Email:         email,
				EnvironmentId: environmentNamespace,
			})
			if err != nil {
				return nil, err
			}
			return resp.Account, nil
		})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Info(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
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
			s.logger.Info(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
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
					zap.String("environmentNamespace", environmentNamespace),
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

func (s *auditlogService) checkSystemAdminRole(
	ctx context.Context,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckSystemAdminRole(ctx)
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Info(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
			s.logger.Info(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
