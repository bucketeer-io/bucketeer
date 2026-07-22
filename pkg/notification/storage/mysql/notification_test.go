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

package mysql

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/notification/domain"
	notificationstorage "github.com/bucketeer-io/bucketeer/v2/pkg/notification/storage"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
)

func TestNewNotificationStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewNotificationStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &notificationStorage{}, storage)
}

func TestCreateNotification(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	notification := &domain.Notification{
		Notification: &proto.Notification{
			Id:           "notification-id-0",
			Status:       proto.Notification_DRAFT,
			CreatedBy:    "admin@example.com",
			LastEditedBy: "admin@example.com",
			CreatedAt:    1,
			UpdatedAt:    1,
			Localizations: []*proto.NotificationLocalization{
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
			},
		},
	}

	patterns := []struct {
		desc        string
		setup       func(*notificationStorage)
		input       *domain.Notification
		expectedErr error
	}{
		{
			desc: "ErrNotificationAlreadyExists",
			setup: func(s *notificationStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input:       notification,
			expectedErr: notificationstorage.ErrNotificationAlreadyExists,
		},
		{
			desc: "Error: insert notification",
			setup: func(s *notificationStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       notification,
			expectedErr: errors.New("error"),
		},
		{
			desc: "Error: insert localization",
			setup: func(s *notificationStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), insertNotificationSQL, gomock.Any(),
				).Return(nil, nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), insertNotificationLocalizationSQL, gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       notification,
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *notificationStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					insertNotificationSQL,
					"notification-id-0",
					int32(proto.Notification_DRAFT),
					"admin@example.com",
					"admin@example.com",
					int64(1),
					int64(1),
				).Return(nil, nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					insertNotificationLocalizationSQL,
					"notification-id-0",
					"en",
					mysql.JSONObject{Val: notification.Localizations[0].Tags},
					"New feature",
					"# New feature",
				).Return(nil, nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					insertNotificationLocalizationSQL,
					"notification-id-0",
					"ja",
					mysql.JSONObject{Val: notification.Localizations[1].Tags},
					"新機能",
					"# 新機能",
				).Return(nil, nil)
			},
			input:       notification,
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &notificationStorage{qe: mock.NewMockQueryExecer(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateNotification(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListDraftAdminNotifications(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(*notificationStorage)
		params         notificationstorage.ListDraftAdminNotificationsParams
		expected       []*proto.Notification
		expectedCursor int
		expectedCount  int64
		expectedErr    error
	}{
		{
			desc: "ErrInvalidListDraftAdminNotificationsOrderBy",
			params: notificationstorage.ListDraftAdminNotificationsParams{
				OrderBy: proto.ListDraftAdminNotificationsRequest_OrderBy(99),
			},
			expectedErr: notificationstorage.ErrInvalidListDraftAdminNotificationsOrderBy,
		},
		{
			desc: "ErrInvalidListDraftAdminNotificationsCursor",
			params: notificationstorage.ListDraftAdminNotificationsParams{
				Cursor: "invalid",
			},
			expectedErr: notificationstorage.ErrInvalidListDraftAdminNotificationsCursor,
		},
		{
			desc: "ErrInvalidListDraftAdminNotificationsCursor: negative",
			params: notificationstorage.ListDraftAdminNotificationsParams{
				Cursor: "-1",
			},
			expectedErr: notificationstorage.ErrInvalidListDraftAdminNotificationsCursor,
		},
		{
			desc: "Success: negative page size clamped",
			setup: func(s *notificationStorage) {
				listRows := mock.NewMockRows(mockController)
				listRows.EXPECT().Close().Return(nil)
				listRows.EXPECT().Next().Return(false)
				listRows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(listRows, nil)
				countRow := mock.NewMockRow(mockController)
				countRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
					*args[0].(*int64) = int64(0)
					return nil
				})
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(countRow)
			},
			params: notificationstorage.ListDraftAdminNotificationsParams{
				PageSize: -1,
			},
			expected:       []*proto.Notification{},
			expectedCursor: 0,
			expectedCount:  0,
			expectedErr:    nil,
		},
		{
			desc: "Error",
			setup: func(s *notificationStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			params:      notificationstorage.ListDraftAdminNotificationsParams{},
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *notificationStorage) {
				var listNextCallCount = 0
				listRows := mock.NewMockRows(mockController)
				listRows.EXPECT().Close().Return(nil)
				listRows.EXPECT().Next().DoAndReturn(func() bool {
					listNextCallCount++
					return listNextCallCount <= 1
				}).Times(2)
				listRows.EXPECT().Err().Return(nil)
				listRows.EXPECT().Scan(
					gomock.Any(), // id
					gomock.Any(), // status
					gomock.Any(), // created_by
					gomock.Any(), // last_edited_by
					gomock.Any(), // created_at
					gomock.Any(), // updated_at
				).Do(func(args ...interface{}) {
					*args[0].(*string) = "notification-id-0"
					*args[1].(*int32) = int32(proto.Notification_DRAFT)
					*args[2].(*string) = "admin@example.com"
					*args[3].(*string) = "admin@example.com"
					*args[4].(*int64) = int64(1)
					*args[5].(*int64) = int64(2)
				}).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(listRows, nil)
				var locNextCallCount = 0
				locRows := mock.NewMockRows(mockController)
				locRows.EXPECT().Close().Return(nil)
				locRows.EXPECT().Next().DoAndReturn(func() bool {
					locNextCallCount++
					return locNextCallCount <= 1
				}).Times(2)
				locRows.EXPECT().Err().Return(nil)
				locRows.EXPECT().Scan(
					gomock.Any(), // notification_id
					gomock.Any(), // language
					gomock.Any(), // tags
					gomock.Any(), // title
					gomock.Any(), // content
				).Do(func(args ...interface{}) {
					*args[0].(*string) = "notification-id-0"
					*args[1].(*string) = "en"
					*args[3].(*string) = "New feature"
					*args[4].(*string) = "# New feature"
				}).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(locRows, nil)
				countRow := mock.NewMockRow(mockController)
				countRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
					*args[0].(*int64) = int64(1)
					return nil
				})
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(countRow)
			},
			params: notificationstorage.ListDraftAdminNotificationsParams{
				SearchKeyword: "feature",
				PageSize:      10,
			},
			expected: []*proto.Notification{
				{
					Id:           "notification-id-0",
					Status:       proto.Notification_DRAFT,
					CreatedBy:    "admin@example.com",
					LastEditedBy: "admin@example.com",
					CreatedAt:    1,
					UpdatedAt:    2,
					Localizations: []*proto.NotificationLocalization{
						{
							Language: "en",
							Title:    "New feature",
							Content:  "# New feature",
						},
					},
				},
			},
			expectedCursor: 1,
			expectedCount:  1,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &notificationStorage{qe: mock.NewMockQueryExecer(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			notifications, cursor, count, err := storage.ListDraftAdminNotifications(context.Background(), p.params)
			assert.Equal(t, p.expectedErr, err)
			if p.expectedErr == nil {
				assert.Equal(t, p.expected, notifications)
				assert.Equal(t, p.expectedCursor, cursor)
				assert.Equal(t, p.expectedCount, count)
			}
		})
	}
}
