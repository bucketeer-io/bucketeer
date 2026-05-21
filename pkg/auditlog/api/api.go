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
	"errors"
	"sort"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	v2als "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

const (
	// Maximum page size for audit logs. Also used as default when page_size is not set or exceeds this value.
	maxAuditLogPageSize = 200
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
	accountStorage v2as.AccountStorage,
	auditLogStorage v2als.AuditLogStorage,
	adminAuditLogStorage v2als.AdminAuditLogStorage,
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
		accountStorage:       accountStorage,
		auditLogStorage:      auditLogStorage,
		adminAuditLogStorage: adminAuditLogStorage,
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
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		s.logger.Error("Missing audit log id",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, statusMissingID.Err()
	}
	auditlog, err := s.auditLogStorage.GetAuditLog(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		s.logger.Error("Failed to get audit log",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		if errors.Is(err, v2als.ErrAuditLogNotFound) {
			return nil, statusAuditLogNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	auditlog.LocalizedMessage = domainevent.LocalizedMessage(auditlog.Type, localizer)

	accounts, err := s.getAccountMapByEmails(ctx, []string{auditlog.Editor.Email}, req.EnvironmentId)
	if err != nil {
		s.logger.Error("Failed to get account map by emails",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	// Use maximum page size as default when not provided, is 0, or exceeds the maximum
	limit := int(req.PageSize)
	if limit <= 0 || limit > maxAuditLogPageSize {
		limit = maxAuditLogPageSize
	}

	cursor := req.Cursor
	if cursor == "" {
		cursor = "0"
	}
	// Validate cursor before passing to storage
	if _, err := strconv.Atoi(cursor); err != nil {
		return nil, statusInvalidCursor.Err()
	}

	var entityType *int32
	if req.EntityType != nil {
		v := req.EntityType.Value
		entityType = &v
	}

	params := v2als.ListAuditLogsParams{
		EnvironmentID:  req.EnvironmentId,
		EntityType:     entityType,
		From:           req.From,
		To:             req.To,
		SearchKeyword:  req.SearchKeyword,
		OrderBy:        req.OrderBy,
		OrderDirection: req.OrderDirection,
		PageSize:       limit,
		Cursor:         cursor,
	}
	auditlogs, nextCursor, totalCount, err := s.auditLogStorage.ListAuditLogs(ctx, params)
	if err != nil {
		s.logger.Error(
			"Failed to list auditlogs",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	editorEmails := make([]string, 0, len(auditlogs))
	for _, auditlog := range auditlogs {
		editorEmails = append(editorEmails, auditlog.Editor.Email)
	}
	editorEmails = deDuplicateStrings(editorEmails)
	accounts, err := s.getAccountMapByEmails(ctx, editorEmails, req.EnvironmentId)
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

func (s *auditlogService) ListAdminAuditLogs(
	ctx context.Context,
	req *proto.ListAdminAuditLogsRequest,
) (*proto.ListAdminAuditLogsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}

	// Use maximum page size as default when not provided, is 0, or exceeds the maximum
	limit := int(req.PageSize)
	if limit <= 0 || limit > maxAuditLogPageSize {
		limit = maxAuditLogPageSize
	}

	cursor := req.Cursor
	if cursor == "" {
		cursor = "0"
	}
	// Validate cursor before passing to storage
	if _, err := strconv.Atoi(cursor); err != nil {
		return nil, statusInvalidCursor.Err()
	}

	var entityType *int32
	if req.EntityType != nil {
		v := req.EntityType.Value
		entityType = &v
	}

	params := v2als.ListAdminAuditLogsParams{
		EntityType:     entityType,
		From:           req.From,
		To:             req.To,
		SearchKeyword:  req.SearchKeyword,
		OrderBy:        req.OrderBy,
		OrderDirection: req.OrderDirection,
		PageSize:       limit,
		Cursor:         cursor,
	}
	auditlogs, nextCursor, totalCount, err := s.adminAuditLogStorage.ListAdminAuditLogs(ctx, params)
	if err != nil {
		s.logger.Error(
			"Failed to list admin auditlogs",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, api.NewGRPCStatus(err).Err()
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

func (s *auditlogService) ListFeatureHistory(
	ctx context.Context,
	req *proto.ListFeatureHistoryRequest,
) (*proto.ListFeatureHistoryResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	// Use maximum page size as default when not provided, is 0, or exceeds the maximum
	limit := int(req.PageSize)
	if limit <= 0 || limit > maxAuditLogPageSize {
		limit = maxAuditLogPageSize
	}

	cursor := req.Cursor
	if cursor == "" {
		cursor = "0"
	}
	// Validate cursor before passing to storage
	if _, err := strconv.Atoi(cursor); err != nil {
		return nil, statusInvalidCursor.Err()
	}

	entityType := int32(eventproto.Event_FEATURE)
	params := v2als.ListAuditLogsParams{
		EnvironmentID:  req.EnvironmentId,
		EntityType:     &entityType,
		EntityID:       req.FeatureId,
		From:           req.From,
		To:             req.To,
		SearchKeyword:  req.SearchKeyword,
		OrderBy:        proto.ListAuditLogsRequest_OrderBy(req.OrderBy),
		OrderDirection: proto.ListAuditLogsRequest_OrderDirection(req.OrderDirection),
		PageSize:       limit,
		Cursor:         cursor,
	}
	auditlogs, nextCursor, totalCount, err := s.auditLogStorage.ListAuditLogs(ctx, params)
	if err != nil {
		s.logger.Error(
			"Failed to list feature history",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("featureId", req.FeatureId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
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
) (accountMap map[string]*accountproto.AccountV2, err error) {
	accountMap = make(map[string]*accountproto.AccountV2)
	if len(emails) == 0 {
		return accountMap, nil
	}
	emailsArg := make([]interface{}, len(emails))
	for i, email := range emails {
		emailsArg[i] = email
	}
	// TODO: Refactor account storage to use DB-agnostic params when account package is migrated.
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Strings("emails", emails),
				zap.String("environmentId", environmentID),
			)...,
		)
		return accountMap, api.NewGRPCStatus(err).Err()
	}
	for i := range accounts {
		accountMap[accounts[i].Email] = accounts[i]
	}
	return accountMap, nil
}

func (s *auditlogService) checkEnvironmentRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Environment,
	environmentId string,
) (*eventproto.Editor, error) {
	return role.CheckEnvironmentRoleWithLog(
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
		},
		s.logger,
		statusUnauthenticated.Err(),
		statusPermissionDenied.Err(),
		func(err error) error { return api.NewGRPCStatus(err).Err() },
	)
}

func (s *auditlogService) checkSystemAdminRole(
	ctx context.Context,
) (*eventproto.Editor, error) {
	editor, err := role.CheckSystemAdminRole(ctx)
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, statusUnauthenticated.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, statusPermissionDenied.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, api.NewGRPCStatus(err).Err()
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
