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

func TestCreateAdminNotification(t *testing.T) {
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
			err := storage.CreateAdminNotification(context.Background(), p.input)
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

func TestGetAdminNotification(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*notificationStorage)
		id          string
		expected    *domain.Notification
		expectedErr error
	}{
		{
			desc: "ErrNotificationNotFound",
			setup: func(s *notificationStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "notification-id-0",
			expectedErr: notificationstorage.ErrNotificationNotFound,
		},
		{
			desc: "Error",
			setup: func(s *notificationStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "notification-id-0",
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *notificationStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
					*args[0].(*string) = "notification-id-0"
					*args[1].(*int32) = int32(proto.Notification_DRAFT)
					*args[2].(*string) = "admin@example.com"
					*args[3].(*string) = "admin@example.com"
					*args[4].(*string) = ""
					*args[5].(*int64) = int64(0)
					*args[6].(*int64) = int64(1)
					*args[7].(*int64) = int64(2)
					return nil
				})
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectNotificationSQL, "notification-id-0",
				).Return(row)
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
			},
			id: "notification-id-0",
			expected: &domain.Notification{
				Notification: &proto.Notification{
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
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &notificationStorage{qe: mock.NewMockQueryExecer(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			notification, err := storage.GetAdminNotification(context.Background(), p.id)
			assert.Equal(t, p.expectedErr, err)
			if p.expectedErr == nil {
				assert.Equal(t, p.expected, notification)
			}
		})
	}
}

func TestUpdateAdminNotification(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	notification := &domain.Notification{
		Notification: &proto.Notification{
			Id:           "notification-id-0",
			Status:       proto.Notification_DRAFT,
			CreatedBy:    "admin@example.com",
			LastEditedBy: "editor@example.com",
			CreatedAt:    1,
			UpdatedAt:    5,
			Localizations: []*proto.NotificationLocalization{
				{
					Language: "en",
					Tags:     []*proto.NotificationTag{{Name: "Announcement", Color: "#3B82F6"}},
					Title:    "Updated title",
					Content:  "# Updated content",
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
			desc: "Error: update notification",
			setup: func(s *notificationStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       notification,
			expectedErr: errors.New("error"),
		},
		{
			desc: "ErrNotificationNotFound: already deleted",
			setup: func(s *notificationStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input:       notification,
			expectedErr: notificationstorage.ErrNotificationNotFound,
		},
		{
			desc: "Error: delete localizations",
			setup: func(s *notificationStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), updateNotificationSQL, gomock.Any(),
				).Return(result, nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), deleteNotificationLocalizationsSQL, gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       notification,
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *notificationStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					updateNotificationSQL,
					"editor@example.com",
					int64(5),
					"notification-id-0",
				).Return(result, nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					deleteNotificationLocalizationsSQL,
					"notification-id-0",
				).Return(nil, nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					insertNotificationLocalizationSQL,
					"notification-id-0",
					"en",
					mysql.JSONObject{Val: notification.Localizations[0].Tags},
					"Updated title",
					"# Updated content",
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
			err := storage.UpdateAdminNotification(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteAdminNotification(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*notificationStorage)
		id          string
		expectedErr error
	}{
		{
			desc: "Error: delete notification",
			setup: func(s *notificationStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			id:          "notification-id-0",
			expectedErr: errors.New("error"),
		},
		{
			desc: "ErrNotificationNotFound",
			setup: func(s *notificationStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			id:          "notification-id-0",
			expectedErr: notificationstorage.ErrNotificationNotFound,
		},
		{
			desc: "Error: rows affected",
			setup: func(s *notificationStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), errors.New("error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			id:          "notification-id-0",
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *notificationStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					deleteNotificationSQL,
					"editor@example.com",
					int64(5),
					"notification-id-0",
				).Return(result, nil)
			},
			id:          "notification-id-0",
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &notificationStorage{qe: mock.NewMockQueryExecer(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.DeleteAdminNotification(context.Background(), p.id, "editor@example.com", 5)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestPublishAdminNotification(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	notification := &domain.Notification{
		Notification: &proto.Notification{
			Id:           "notification-id-0",
			Status:       proto.Notification_PUBLISHED,
			CreatedBy:    "admin@example.com",
			LastEditedBy: "admin@example.com",
			PublishedBy:  "publisher@example.com",
			PublishedAt:  5,
			CreatedAt:    1,
			UpdatedAt:    5,
		},
	}

	patterns := []struct {
		desc        string
		setup       func(*notificationStorage)
		input       *domain.Notification
		expectedErr error
	}{
		{
			desc: "Error: publish notification",
			setup: func(s *notificationStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       notification,
			expectedErr: errors.New("error"),
		},
		{
			desc: "ErrNotificationAlreadyPublished: no longer a live draft",
			setup: func(s *notificationStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input:       notification,
			expectedErr: notificationstorage.ErrNotificationAlreadyPublished,
		},
		{
			desc: "Error: rows affected",
			setup: func(s *notificationStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), errors.New("error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input:       notification,
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *notificationStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					publishNotificationSQL,
					int32(proto.Notification_PUBLISHED),
					"publisher@example.com",
					int64(5),
					int64(5),
					"notification-id-0",
					int32(proto.Notification_DRAFT),
				).Return(result, nil)
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
			err := storage.PublishAdminNotification(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListNotifications(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	newReadIDRows := func(ids ...string) *mock.MockRows {
		rows := mock.NewMockRows(mockController)
		rows.EXPECT().Close().Return(nil)
		i := 0
		rows.EXPECT().Next().DoAndReturn(func() bool {
			i++
			return i <= len(ids)
		}).Times(len(ids) + 1)
		if len(ids) > 0 {
			j := 0
			rows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
				*args[0].(*string) = ids[j]
				j++
				return nil
			}).Times(len(ids))
		}
		rows.EXPECT().Err().Return(nil)
		return rows
	}
	newListRows := func() *mock.MockRows {
		rows := mock.NewMockRows(mockController)
		rows.EXPECT().Close().Return(nil)
		i := 0
		rows.EXPECT().Next().DoAndReturn(func() bool {
			i++
			return i <= 1
		}).Times(2)
		rows.EXPECT().Err().Return(nil)
		rows.EXPECT().Scan(
			gomock.Any(), // id
			gomock.Any(), // status
			gomock.Any(), // created_by
			gomock.Any(), // last_edited_by
			gomock.Any(), // published_by
			gomock.Any(), // published_at
			gomock.Any(), // created_at
			gomock.Any(), // updated_at
		).Do(func(args ...interface{}) {
			*args[0].(*string) = "notification-id-0"
			*args[1].(*int32) = int32(proto.Notification_PUBLISHED)
			*args[2].(*string) = "admin@example.com"
			*args[3].(*string) = "admin@example.com"
			*args[4].(*string) = "publisher@example.com"
			*args[5].(*int64) = int64(10)
			*args[6].(*int64) = int64(1)
			*args[7].(*int64) = int64(10)
		}).Return(nil)
		return rows
	}
	newLocRows := func() *mock.MockRows {
		rows := mock.NewMockRows(mockController)
		rows.EXPECT().Close().Return(nil)
		i := 0
		rows.EXPECT().Next().DoAndReturn(func() bool {
			i++
			return i <= 2
		}).Times(3)
		j := 0
		rows.EXPECT().Scan(
			gomock.Any(), // notification_id
			gomock.Any(), // language
			gomock.Any(), // tags
			gomock.Any(), // title
			gomock.Any(), // content
		).DoAndReturn(func(args ...interface{}) error {
			j++
			*args[0].(*string) = "notification-id-0"
			if j == 1 {
				*args[1].(*string) = "en"
				*args[3].(*string) = "New feature"
				*args[4].(*string) = "# New feature"
			} else {
				*args[1].(*string) = "ja"
				*args[3].(*string) = "新機能"
				*args[4].(*string) = "# 新機能"
			}
			return nil
		}).Times(2)
		rows.EXPECT().Err().Return(nil)
		return rows
	}
	newCountRow := func() *mock.MockRow {
		row := mock.NewMockRow(mockController)
		row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
			*args[0].(*int64) = int64(1)
			return nil
		})
		return row
	}
	published := func(localization *proto.NotificationLocalization, read bool) []*proto.Notification {
		return []*proto.Notification{
			{
				Id:           "notification-id-0",
				Status:       proto.Notification_PUBLISHED,
				CreatedBy:    "admin@example.com",
				LastEditedBy: "admin@example.com",
				PublishedBy:  "publisher@example.com",
				PublishedAt:  10,
				CreatedAt:    1,
				UpdatedAt:    10,
				Localization: localization,
				Read:         read,
			},
		}
	}

	patterns := []struct {
		desc           string
		setup          func(*notificationStorage)
		params         notificationstorage.ListNotificationsParams
		expected       []*proto.Notification
		expectedCursor int
		expectedCount  int64
		expectedErr    error
	}{
		{
			desc: "ErrInvalidListNotificationsOrderBy",
			params: notificationstorage.ListNotificationsParams{
				OrderBy: proto.ListNotificationsRequest_OrderBy(99),
			},
			expectedErr: notificationstorage.ErrInvalidListNotificationsOrderBy,
		},
		{
			desc: "ErrInvalidListNotificationsCursor",
			params: notificationstorage.ListNotificationsParams{
				Cursor: "invalid",
			},
			expectedErr: notificationstorage.ErrInvalidListNotificationsCursor,
		},
		{
			desc: "Success: all, resolved to requested language, read flag set",
			setup: func(s *notificationStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(newListRows(), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(newLocRows(), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(newReadIDRows("notification-id-0"), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(newCountRow())
			},
			params: notificationstorage.ListNotificationsParams{
				Email:    "viewer@example.com",
				Language: "ja",
				PageSize: 10,
			},
			expected: published(&proto.NotificationLocalization{
				Language: "ja",
				Title:    "新機能",
				Content:  "# 新機能",
			}, true),
			expectedCursor: 1,
			expectedCount:  1,
			expectedErr:    nil,
		},
		{
			desc: "Success: unread, fallback to English, read flag unset",
			setup: func(s *notificationStorage) {
				boundRow := mock.NewMockRow(mockController)
				boundRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
					*args[0].(*int64) = int64(5)
					return nil
				})
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), selectEarliestAccountCreatedAtSQL, "viewer@example.com",
				).Return(boundRow)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(newListRows(), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(newLocRows(), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(newReadIDRows(), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(newCountRow())
			},
			params: notificationstorage.ListNotificationsParams{
				Email:      "viewer@example.com",
				Language:   "fr",
				ReadStatus: proto.ListNotificationsRequest_UNREAD,
				PageSize:   10,
			},
			expected: published(&proto.NotificationLocalization{
				Language: "en",
				Title:    "New feature",
				Content:  "# New feature",
			}, false),
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
			notifications, cursor, count, err := storage.ListNotifications(context.Background(), p.params)
			assert.Equal(t, p.expectedErr, err)
			if p.expectedErr == nil {
				assert.Equal(t, p.expected, notifications)
				assert.Equal(t, p.expectedCursor, cursor)
				assert.Equal(t, p.expectedCount, count)
			}
		})
	}
}

func TestIsNotificationRead(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*notificationStorage)
		expected    bool
		expectedErr error
	}{
		{
			desc: "Error",
			setup: func(s *notificationStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expected:    false,
			expectedErr: errors.New("error"),
		},
		{
			desc: "Read",
			setup: func(s *notificationStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				i := 0
				rows.EXPECT().Next().DoAndReturn(func() bool {
					i++
					return i <= 1
				}).Times(2)
				rows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
					*args[0].(*string) = "notification-id-0"
					return nil
				})
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expected:    true,
			expectedErr: nil,
		},
		{
			desc: "Unread",
			setup: func(s *notificationStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expected:    false,
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &notificationStorage{qe: mock.NewMockQueryExecer(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			read, err := storage.IsNotificationRead(context.Background(), "notification-id-0", "viewer@example.com")
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, read)
		})
	}
}

func TestMarkNotificationsAsRead(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*notificationStorage)
		ids         []string
		expectedErr error
	}{
		{
			desc: "Error",
			setup: func(s *notificationStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			ids:         []string{"notification-id-0"},
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success: one upsert per id",
			setup: func(s *notificationStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					insertNotificationReadSQL,
					"viewer@example.com",
					int64(5),
					"notification-id-0",
					int32(proto.Notification_PUBLISHED),
				).Return(nil, nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					insertNotificationReadSQL,
					"viewer@example.com",
					int64(5),
					"notification-id-1",
					int32(proto.Notification_PUBLISHED),
				).Return(nil, nil)
			},
			ids:         []string{"notification-id-0", "notification-id-1"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &notificationStorage{qe: mock.NewMockQueryExecer(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.MarkNotificationsAsRead(context.Background(), p.ids, "viewer@example.com", 5)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
