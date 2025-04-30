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
	"sort"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	v2as "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
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
	accountClient        accountclient.Client
	accountStorage       v2as.AccountStorage
	auditLogStorage      v2als.AuditLogStorage
	adminAuditLogStorage v2als.AdminAuditLogStorage
	opts                 *options
	logger               *zap.Logger
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
		accountClient:        accountClient,
		accountStorage:       v2as.NewAccountStorage(mysqlClient),
		auditLogStorage:      v2als.NewAuditLogStorage(mysqlClient),
		adminAuditLogStorage: v2als.NewAdminAuditLogStorage(mysqlClient),
		opts:                 dopts,
		logger:               dopts.logger.Named("api"),
	}
}

func (s *auditlogService) Register(server *grpc.Server) {
	proto.RegisterAuditLogServiceServer(server, s)
}

func (s *auditlogService) GetAuditLog(
	ctx context.Context,
	req *proto.GetAuditLogRequest,
) (*proto.GetAuditLogResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		s.logger.Error("Missing audit log id",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		dt, err := statusMissingID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "id"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	auditlog, err := s.auditLogStorage.GetAuditLog(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		var dt *status.Status
		s.logger.Error("Failed to get audit log",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		if errors.Is(err, v2als.ErrAuditLogNotFound) {
			dt, err = statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.NotFoundError, "auditlog"),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
		} else {
			dt, err = statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
		}
		return nil, dt.Err()
	}
	auditlog.LocalizedMessage = domainevent.LocalizedMessage(auditlog.Type, localizer)

	accounts, err := s.getAccountMapByEmails(ctx, []string{auditlog.Editor.Email}, req.EnvironmentId, localizer)
	if err != nil {
		s.logger.Error("Failed to get account map by emails",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
				zap.String("email", auditlog.Editor.Email),
			)...,
		)
		// return without avatar image
		return &proto.GetAuditLogResponse{
			AuditLog: auditlog,
		}, nil
	}
	if account, ok := accounts[auditlog.Editor.Email]; ok {
		auditlog.Editor.AvatarImage = account.AvatarImage
		auditlog.Editor.AvatarFileType = account.AvatarFileType
		if auditlog.Editor.PublicApiEditor != nil {
			auditlog.Editor.PublicApiEditor.AvatarImage = account.AvatarImage
			auditlog.Editor.PublicApiEditor.AvatarFileType = account.AvatarFileType
		}
	}
	return &proto.GetAuditLogResponse{
		AuditLog: auditlog,
	}, nil
}

func (s *auditlogService) ListAuditLogs(
	ctx context.Context,
	req *proto.ListAuditLogsRequest,
) (*proto.ListAuditLogsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
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
		mysql.NewFilter("environment_id", "=", req.EnvironmentId),
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
	auditlogs, nextCursor, totalCount, err := s.auditLogStorage.ListAuditLogs(
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
	editorEmails := make([]string, 0, len(auditlogs))
	for _, auditlog := range auditlogs {
		editorEmails = append(editorEmails, auditlog.Editor.Email)
	}
	editorEmails = deDuplicateStrings(editorEmails)
	accounts, err := s.getAccountMapByEmails(ctx, editorEmails, req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}

	for i := range auditlogs {
		if account, ok := accounts[auditlogs[i].Editor.Email]; ok {
			auditlogs[i].Editor.AvatarImage = account.AvatarImage
			auditlogs[i].Editor.AvatarFileType = account.AvatarFileType
			if auditlogs[i].Editor.PublicApiEditor != nil {
				auditlogs[i].Editor.PublicApiEditor.AvatarImage = account.AvatarImage
				auditlogs[i].Editor.PublicApiEditor.AvatarFileType = account.AvatarFileType
			}
		}
		auditlogs[i].LocalizedMessage = domainevent.LocalizedMessage(auditlogs[i].Type, localizer)
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
	auditlogs, nextCursor, totalCount, err := s.adminAuditLogStorage.ListAdminAuditLogs(
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
		req.EnvironmentId,
		localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("environment_id", "=", req.EnvironmentId),
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
	auditlogs, nextCursor, totalCount, err := s.auditLogStorage.ListAuditLogs(
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
				zap.String("environmentId", req.EnvironmentId),
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

func (s *auditlogService) getAccountMapByEmails(
	ctx context.Context,
	emails []string,
	environmentID string,
	localizer locale.Localizer,
) (accountMap map[string]*accountproto.AccountV2, err error) {
	accountMap = make(map[string]*accountproto.AccountV2)
	if len(emails) == 0 {
		return accountMap, nil
	}
	emailsArg := make([]interface{}, len(emails))
	for i, email := range emails {
		emailsArg[i] = email
	}
	options := &mysql.ListOptions{
		Limit:  0,
		Offset: 0,
		Orders: nil,
		InFilters: []*mysql.InFilter{
			{
				Column: "a.email",
				Values: emailsArg,
			},
		},
		Filters: []*mysql.FilterV2{
			{
				Column:   "e.id",
				Operator: mysql.OperatorEqual,
				Value:    environmentID,
			},
		},
	}
	accounts, err := s.accountStorage.GetAvatarAccountsV2(ctx, options)
	if err != nil {
		s.logger.Error(
			"Failed to list feature history",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Strings("emails", emails),
				zap.String("environmentId", environmentID),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return accountMap, statusInternal.Err()
		}
		return accountMap, dt.Err()
	}
	for i := range accounts {
		accountMap[accounts[i].Email] = accounts[i]
	}
	return accountMap, nil
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
	environmentId string,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckEnvironmentRole(
		ctx,
		requiredRole,
		environmentId,
		func(email string) (*accountproto.AccountV2, error) {
			resp, err := s.accountClient.GetAccountV2ByEnvironmentID(ctx, &accountproto.GetAccountV2ByEnvironmentIDRequest{
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

func (s *auditlogService) checkSystemAdminRole(
	ctx context.Context,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckSystemAdminRole(ctx)
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
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
			s.logger.Error(
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

func deDuplicateStrings(args []string) []string {
	sort.Strings(args)
	var result []string
	for i := 0; i < len(args); i++ {
		if i == 0 || args[i] != args[i-1] {
			result = append(result, args[i])
		}
	}
	return result
}
