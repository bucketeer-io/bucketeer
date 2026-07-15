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

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"

	proto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
)

var statusNotImplemented = gstatus.Error(codes.Unimplemented, "notification: not implemented")

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
	opts   *options
	logger *zap.Logger
}

func NewNotificationService(opts ...Option) *NotificationService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &NotificationService{
		opts:   dopts,
		logger: dopts.logger.Named("api"),
	}
}

func (s *NotificationService) Register(server *grpc.Server) {
	proto.RegisterNotificationServiceServer(server, s)
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
	return nil, statusNotImplemented
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
