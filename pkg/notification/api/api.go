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
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/notification/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/notification/storage"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
)

const (
	// Maximum page size for notifications. Also used as default when page_size
	// is not set or exceeds this value.
	maxNotificationPageSize = 200
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

type NotificationService struct {
	dbClient            database.Client
	notificationStorage storage.NotificationStorage
	opts                *options
	logger              *zap.Logger
}

func NewNotificationService(
	dbClient database.Client,
	notificationStorage storage.NotificationStorage,
	opts ...Option,
) *NotificationService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &NotificationService{
		dbClient:            dbClient,
		notificationStorage: notificationStorage,
		opts:                dopts,
		logger:              dopts.logger.Named("api"),
	}
}

func (s *NotificationService) Register(server *grpc.Server) {
	proto.RegisterNotificationServiceServer(server, s)
}

func (s *NotificationService) checkSystemAdminRole(
	ctx context.Context,
) (*eventproto.Editor, error) {
	editor, err := role.CheckSystemAdminRole(ctx)
	if err != nil {
		switch gstatus.Code(err) {
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

func (s *NotificationService) ListNotifications(
	ctx context.Context,
	req *proto.ListNotificationsRequest,
) (*proto.ListNotificationsResponse, error) {
	return nil, statusNotImplemented
}

func (s *NotificationService) GetNotification(
	ctx context.Context,
	req *proto.GetNotificationRequest,
) (*proto.GetNotificationResponse, error) {
	return nil, statusNotImplemented
}

func (s *NotificationService) GetNotificationUnreadCount(
	ctx context.Context,
	req *proto.GetNotificationUnreadCountRequest,
) (*proto.GetNotificationUnreadCountResponse, error) {
	return nil, statusNotImplemented
}

func (s *NotificationService) MarkNotificationsAsRead(
	ctx context.Context,
	req *proto.MarkNotificationsAsReadRequest,
) (*proto.MarkNotificationsAsReadResponse, error) {
	return nil, statusNotImplemented
}

func (s *NotificationService) MarkAllNotificationsAsRead(
	ctx context.Context,
	req *proto.MarkAllNotificationsAsReadRequest,
) (*proto.MarkAllNotificationsAsReadResponse, error) {
	return nil, statusNotImplemented
}

func (s *NotificationService) ListDraftAdminNotifications(
	ctx context.Context,
	req *proto.ListDraftAdminNotificationsRequest,
) (*proto.ListDraftAdminNotificationsResponse, error) {
	_, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	limit := int(req.PageSize)
	if limit <= 0 || limit > maxNotificationPageSize {
		limit = maxNotificationPageSize
	}
	params := storage.ListDraftAdminNotificationsParams{
		SearchKeyword:  req.SearchKeyword,
		OrderBy:        req.OrderBy,
		OrderDirection: req.OrderDirection,
		PageSize:       limit,
		Cursor:         req.Cursor,
	}
	notifications, nextOffset, totalCount, err := s.notificationStorage.ListDraftAdminNotifications(ctx, params)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidListDraftAdminNotificationsCursor) {
			return nil, statusInvalidCursor.Err()
		}
		if errors.Is(err, storage.ErrInvalidListDraftAdminNotificationsOrderBy) {
			return nil, statusInvalidOrderBy.Err()
		}
		s.logger.Error(
			"Failed to list draft notifications",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.ListDraftAdminNotificationsResponse{
		Notifications: notifications,
		NextCursor:    strconv.Itoa(nextOffset),
		TotalCount:    totalCount,
	}, nil
}

func (s *NotificationService) CreateAdminNotification(
	ctx context.Context,
	req *proto.CreateAdminNotificationRequest,
) (*proto.CreateAdminNotificationResponse, error) {
	editor, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := validateCreateAdminNotificationRequest(req); err != nil {
		s.logger.Error(
			"Failed to validate create notification request",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	notification, err := domain.NewNotification(editor.Email, req.Localizations)
	if err != nil {
		s.logger.Error(
			"Failed to create new notification",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.dbClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context) error {
		return s.notificationStorage.CreateAdminNotification(ctxWithTx, notification)
	})
	if err != nil {
		if errors.Is(err, storage.ErrNotificationAlreadyExists) {
			return nil, statusNotificationAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create notification",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("notificationId", notification.Id),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.CreateAdminNotificationResponse{
		Notification: notification.Notification,
	}, nil
}

func validateCreateAdminNotificationRequest(req *proto.CreateAdminNotificationRequest) error {
	return validateLocalizations(req.Localizations)
}

func validateLocalizations(localizations []*proto.NotificationLocalization) error {
	if len(localizations) == 0 {
		return statusLocalizationRequired.Err()
	}
	languages := make(map[string]struct{}, len(localizations))
	for _, l := range localizations {
		l.Language = strings.TrimSpace(l.Language)
		l.Title = strings.TrimSpace(l.Title)
		if l.Language == "" {
			return statusLanguageRequired.Err()
		}
		if _, ok := languages[l.Language]; ok {
			return statusDuplicatedLanguage.Err()
		}
		languages[l.Language] = struct{}{}
		if l.Title == "" {
			return statusTitleRequired.Err()
		}
		if strings.TrimSpace(l.Content) == "" {
			return statusContentRequired.Err()
		}
	}
	return nil
}

func (s *NotificationService) UpdateAdminNotification(
	ctx context.Context,
	req *proto.UpdateAdminNotificationRequest,
) (*proto.UpdateAdminNotificationResponse, error) {
	editor, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := validateUpdateAdminNotificationRequest(req); err != nil {
		s.logger.Error(
			"Failed to validate update notification request",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	var notification *domain.Notification
	err = s.dbClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context) error {
		var err error
		notification, err = s.notificationStorage.GetAdminNotification(ctxWithTx, req.Id)
		if err != nil {
			return err
		}
		if notification.Status != proto.Notification_DRAFT {
			return statusNotificationAlreadyPublished.Err()
		}
		notification.Update(editor.Email, req.Localizations)
		return s.notificationStorage.UpdateAdminNotification(ctxWithTx, notification)
	})
	if err != nil {
		if errors.Is(err, storage.ErrNotificationNotFound) {
			return nil, statusNotificationNotFound.Err()
		}
		if errors.Is(err, statusNotificationAlreadyPublished.Err()) {
			return nil, statusNotificationAlreadyPublished.Err()
		}
		s.logger.Error(
			"Failed to update notification",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("notificationId", req.Id),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.UpdateAdminNotificationResponse{
		Notification: notification.Notification,
	}, nil
}

func validateUpdateAdminNotificationRequest(req *proto.UpdateAdminNotificationRequest) error {
	if len(strings.TrimSpace(req.Id)) == 0 {
		return statusNotificationIDRequired.Err()
	}
	return validateLocalizations(req.Localizations)
}

func (s *NotificationService) PublishAdminNotification(
	ctx context.Context,
	req *proto.PublishAdminNotificationRequest,
) (*proto.PublishAdminNotificationResponse, error) {
	return nil, statusNotImplemented
}

func (s *NotificationService) DeleteAdminNotification(
	ctx context.Context,
	req *proto.DeleteAdminNotificationRequest,
) (*proto.DeleteAdminNotificationResponse, error) {
	editor, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(req.Id)) == 0 {
		return nil, statusNotificationIDRequired.Err()
	}
	err = s.dbClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context) error {
		return s.notificationStorage.DeleteAdminNotification(ctxWithTx, req.Id, editor.Email, time.Now().Unix())
	})
	if err != nil {
		if errors.Is(err, storage.ErrNotificationNotFound) {
			return nil, statusNotificationNotFound.Err()
		}
		s.logger.Error(
			"Failed to delete notification",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("notificationId", req.Id),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.DeleteAdminNotificationResponse{}, nil
}
