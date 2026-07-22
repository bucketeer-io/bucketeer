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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/notification/storage"
	notificationstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/notification/storage/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	databasemock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
)

func TestNewNotificationService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	dbClientMock := databasemock.NewMockClient(mockController)
	notificationStorageMock := notificationstoragemock.NewMockNotificationStorage(mockController)
	s := NewNotificationService(
		dbClientMock,
		notificationStorageMock,
		WithLogger(zap.NewNop()),
	)
	assert.IsType(t, &NotificationService{}, s)
}

func TestNotificationService_CreateNotification(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	adminCtx := metadata.NewIncomingContext(
		createContextWithToken(t, true),
		metadata.MD{"accept-language": []string{"en"}},
	)
	memberCtx := metadata.NewIncomingContext(
		createContextWithToken(t, false),
		metadata.MD{"accept-language": []string{"en"}},
	)

	validLocalizations := []*proto.NotificationLocalization{
		{
			Language: "en",
			Tags:     []*proto.NotificationTag{{Name: "Announcement", Color: "#3B82F6"}},
			Title:    "New feature",
			Content:  "# New feature",
		},
		{
			Language: "ja",
			Title:    "新機能",
			Content:  "# 新機能",
		},
	}

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*NotificationService)
		req         *proto.CreateNotificationRequest
		expectedErr error
	}{
		{
			desc: "err: unauthenticated",
			ctx:  context.TODO(),
			req: &proto.CreateNotificationRequest{
				Localizations: validLocalizations,
			},
			expectedErr: statusUnauthenticated.Err(),
		},
		{
			desc: "err: permission denied",
			ctx:  memberCtx,
			req: &proto.CreateNotificationRequest{
				Localizations: validLocalizations,
			},
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc:        "err: localization required",
			ctx:         adminCtx,
			req:         &proto.CreateNotificationRequest{},
			expectedErr: statusLocalizationRequired.Err(),
		},
		{
			desc: "err: language required",
			ctx:  adminCtx,
			req: &proto.CreateNotificationRequest{
				Localizations: []*proto.NotificationLocalization{
					{Language: " ", Title: "New feature", Content: "# New feature"},
				},
			},
			expectedErr: statusLanguageRequired.Err(),
		},
		{
			desc: "err: duplicated language",
			ctx:  adminCtx,
			req: &proto.CreateNotificationRequest{
				Localizations: []*proto.NotificationLocalization{
					{Language: "en", Title: "New feature", Content: "# New feature"},
					{Language: "en", Title: "Another", Content: "# Another"},
				},
			},
			expectedErr: statusDuplicatedLanguage.Err(),
		},
		{
			desc: "err: title required",
			ctx:  adminCtx,
			req: &proto.CreateNotificationRequest{
				Localizations: []*proto.NotificationLocalization{
					{Language: "en", Title: " ", Content: "# New feature"},
				},
			},
			expectedErr: statusTitleRequired.Err(),
		},
		{
			desc: "err: content required",
			ctx:  adminCtx,
			req: &proto.CreateNotificationRequest{
				Localizations: []*proto.NotificationLocalization{
					{Language: "en", Title: "New feature", Content: " "},
				},
			},
			expectedErr: statusContentRequired.Err(),
		},
		{
			desc: "err: already exists",
			ctx:  adminCtx,
			setup: func(s *NotificationService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				s.notificationStorage.(*notificationstoragemock.MockNotificationStorage).EXPECT().CreateNotification(
					gomock.Any(), gomock.Any(),
				).Return(storage.ErrNotificationAlreadyExists)
			},
			req: &proto.CreateNotificationRequest{
				Localizations: validLocalizations,
			},
			expectedErr: statusNotificationAlreadyExists.Err(),
		},
		{
			desc: "err: internal",
			ctx:  adminCtx,
			setup: func(s *NotificationService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &proto.CreateNotificationRequest{
				Localizations: validLocalizations,
			},
			expectedErr: api.NewGRPCStatus(errors.New("error")).Err(),
		},
		{
			desc: "success",
			ctx:  adminCtx,
			setup: func(s *NotificationService) {
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				s.notificationStorage.(*notificationstoragemock.MockNotificationStorage).EXPECT().CreateNotification(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &proto.CreateNotificationRequest{
				Localizations: validLocalizations,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createNotificationService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			resp, err := s.CreateNotification(p.ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
			if p.expectedErr == nil {
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.Notification.Id)
				assert.Equal(t, proto.Notification_DRAFT, resp.Notification.Status)
				assert.Equal(t, "email", resp.Notification.CreatedBy)
				assert.Equal(t, "email", resp.Notification.LastEditedBy)
				assert.Equal(t, p.req.Localizations, resp.Notification.Localizations)
			}
		})
	}
}

func createNotificationService(c *gomock.Controller) *NotificationService {
	return &NotificationService{
		dbClient:            databasemock.NewMockClient(c),
		notificationStorage: notificationstoragemock.NewMockNotificationStorage(c),
		opts: &options{
			logger: zap.NewNop(),
		},
		logger: zap.NewNop(),
	}
}

func createContextWithToken(t *testing.T, isSystemAdmin bool) context.Context {
	t.Helper()
	accessToken := &token.AccessToken{
		Issuer:        "issuer",
		Audience:      "audience",
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         "email",
		IsSystemAdmin: isSystemAdmin,
	}
	return context.WithValue(context.TODO(), rpc.AccessTokenKey, accessToken)
}

func TestNotificationService_ListDraftAdminNotifications(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	adminCtx := metadata.NewIncomingContext(
		createContextWithToken(t, true),
		metadata.MD{"accept-language": []string{"en"}},
	)
	memberCtx := metadata.NewIncomingContext(
		createContextWithToken(t, false),
		metadata.MD{"accept-language": []string{"en"}},
	)

	drafts := []*proto.Notification{
		{
			Id:           "notification-id-0",
			Status:       proto.Notification_DRAFT,
			CreatedBy:    "admin@example.com",
			LastEditedBy: "admin@example.com",
			CreatedAt:    1,
			UpdatedAt:    2,
			Localizations: []*proto.NotificationLocalization{
				{Language: "en", Title: "New feature", Content: "# New feature"},
			},
		},
	}

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*NotificationService)
		req         *proto.ListDraftAdminNotificationsRequest
		expectedRes *proto.ListDraftAdminNotificationsResponse
		expectedErr error
	}{
		{
			desc:        "err: unauthenticated",
			ctx:         context.TODO(),
			req:         &proto.ListDraftAdminNotificationsRequest{},
			expectedRes: nil,
			expectedErr: statusUnauthenticated.Err(),
		},
		{
			desc:        "err: permission denied",
			ctx:         memberCtx,
			req:         &proto.ListDraftAdminNotificationsRequest{},
			expectedRes: nil,
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc: "err: invalid cursor",
			ctx:  adminCtx,
			setup: func(s *NotificationService) {
				s.notificationStorage.(*notificationstoragemock.MockNotificationStorage).EXPECT().ListDraftAdminNotifications(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), storage.ErrInvalidListDraftAdminNotificationsCursor)
			},
			req:         &proto.ListDraftAdminNotificationsRequest{Cursor: "invalid"},
			expectedRes: nil,
			expectedErr: statusInvalidCursor.Err(),
		},
		{
			desc: "err: invalid order by",
			ctx:  adminCtx,
			setup: func(s *NotificationService) {
				s.notificationStorage.(*notificationstoragemock.MockNotificationStorage).EXPECT().ListDraftAdminNotifications(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), storage.ErrInvalidListDraftAdminNotificationsOrderBy)
			},
			req:         &proto.ListDraftAdminNotificationsRequest{},
			expectedRes: nil,
			expectedErr: statusInvalidOrderBy.Err(),
		},
		{
			desc: "err: internal",
			ctx:  adminCtx,
			setup: func(s *NotificationService) {
				s.notificationStorage.(*notificationstoragemock.MockNotificationStorage).EXPECT().ListDraftAdminNotifications(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), errors.New("error"))
			},
			req:         &proto.ListDraftAdminNotificationsRequest{},
			expectedRes: nil,
			expectedErr: api.NewGRPCStatus(errors.New("error")).Err(),
		},
		{
			desc: "success",
			ctx:  adminCtx,
			setup: func(s *NotificationService) {
				s.notificationStorage.(*notificationstoragemock.MockNotificationStorage).EXPECT().ListDraftAdminNotifications(
					gomock.Any(),
					storage.ListDraftAdminNotificationsParams{
						SearchKeyword:  "feature",
						OrderBy:        proto.ListDraftAdminNotificationsRequest_UPDATED_AT,
						OrderDirection: proto.ListDraftAdminNotificationsRequest_DESC,
						PageSize:       10,
						Cursor:         "0",
					},
				).Return(drafts, 1, int64(1), nil)
			},
			req: &proto.ListDraftAdminNotificationsRequest{
				PageSize:       10,
				Cursor:         "0",
				OrderBy:        proto.ListDraftAdminNotificationsRequest_UPDATED_AT,
				OrderDirection: proto.ListDraftAdminNotificationsRequest_DESC,
				SearchKeyword:  "feature",
			},
			expectedRes: &proto.ListDraftAdminNotificationsResponse{
				Notifications: drafts,
				NextCursor:    "1",
				TotalCount:    1,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createNotificationService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			res, err := s.ListDraftAdminNotifications(p.ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expectedRes, res)
		})
	}
}
