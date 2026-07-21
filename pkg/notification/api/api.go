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
	"strings"

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

func (s *NotificationService) ListDraftNotifications(
	ctx context.Context,
	req *proto.ListDraftNotificationsRequest,
) (*proto.ListDraftNotificationsResponse, error) {
	return nil, statusNotImplemented
}

func (s *NotificationService) CreateNotification(
	ctx context.Context,
	req *proto.CreateNotificationRequest,
) (*proto.CreateNotificationResponse, error) {
	editor, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := validateCreateNotificationRequest(req); err != nil {
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
		return s.notificationStorage.CreateNotification(ctxWithTx, notification)
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
	return &proto.CreateNotificationResponse{
		Notification: notification.Notification,
	}, nil
}

func validateCreateNotificationRequest(req *proto.CreateNotificationRequest) error {
	if len(req.Localizations) == 0 {
		return statusLocalizationRequired.Err()
	}
	languages := make(map[string]struct{}, len(req.Localizations))
	for _, l := range req.Localizations {
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

func (s *NotificationService) UpdateNotification(
	ctx context.Context,
	req *proto.UpdateNotificationRequest,
) (*proto.UpdateNotificationResponse, error) {
	return nil, statusNotImplemented
}

func (s *NotificationService) PublishNotification(
	ctx context.Context,
	req *proto.PublishNotificationRequest,
) (*proto.PublishNotificationResponse, error) {
	return nil, statusNotImplemented
}

func (s *NotificationService) DeleteNotification(
	ctx context.Context,
	req *proto.DeleteNotificationRequest,
) (*proto.DeleteNotificationResponse, error) {
	return nil, statusNotImplemented
}
