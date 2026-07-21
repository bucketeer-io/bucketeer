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

package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/notification/domain"
	notificationstorage "github.com/bucketeer-io/bucketeer/v2/pkg/notification/storage"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
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
				).Return(nil, postgres.ErrDuplicateEntry)
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
					postgres.JSONObject{Val: notification.Localizations[0].Tags},
					"New feature",
					"# New feature",
				).Return(nil, nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					insertNotificationLocalizationSQL,
					"notification-id-0",
					"ja",
					postgres.JSONObject{Val: notification.Localizations[1].Tags},
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
