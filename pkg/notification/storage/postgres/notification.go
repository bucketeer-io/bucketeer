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
	_ "embed"
	"errors"

	"github.com/bucketeer-io/bucketeer/v2/pkg/notification/domain"
	notificationstorage "github.com/bucketeer-io/bucketeer/v2/pkg/notification/storage"
	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
)

var (
	//go:embed sql/insert_notification.sql
	insertNotificationSQL string
	//go:embed sql/insert_notification_localization.sql
	insertNotificationLocalizationSQL string
)

type notificationStorage struct {
	qe pgstorage.QueryExecer
}

func NewNotificationStorage(qe pgstorage.QueryExecer) notificationstorage.NotificationStorage {
	return &notificationStorage{qe: qe}
}

func (s *notificationStorage) CreateNotification(
	ctx context.Context,
	notification *domain.Notification,
) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertNotificationSQL,
		notification.Id,
		int32(notification.Status),
		notification.CreatedBy,
		notification.LastEditedBy,
		notification.CreatedAt,
		notification.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrDuplicateEntry) {
			return notificationstorage.ErrNotificationAlreadyExists
		}
		return err
	}
	for _, l := range notification.Localizations {
		_, err := s.qe.ExecContext(
			ctx,
			insertNotificationLocalizationSQL,
			notification.Id,
			l.Language,
			pgstorage.JSONObject{Val: l.Tags},
			l.Title,
			l.Content,
		)
		if err != nil {
			if errors.Is(err, pgstorage.ErrDuplicateEntry) {
				return notificationstorage.ErrNotificationAlreadyExists
			}
			return err
		}
	}
	return nil
}
