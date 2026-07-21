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

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	proto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
)

func TestNewNotification(t *testing.T) {
	t.Parallel()
	localizations := []*proto.NotificationLocalization{
		{
			Language: "en",
			Tags:     []*proto.NotificationTag{{Name: "Announcement", Color: "#3B82F6"}},
			Title:    "New feature",
			Content:  "# New feature\nWe released a new feature.",
		},
		{
			Language: "ja",
			Title:    "新機能",
			Content:  "# 新機能\n新機能をリリースしました。",
		},
	}
	notification, err := NewNotification("admin@example.com", localizations)
	assert.Nil(t, err)
	assert.NotEmpty(t, notification.Id)
	assert.Equal(t, proto.Notification_DRAFT, notification.Status)
	assert.Equal(t, "admin@example.com", notification.CreatedBy)
	assert.Equal(t, "admin@example.com", notification.LastEditedBy)
	assert.Empty(t, notification.PublishedBy)
	assert.Zero(t, notification.PublishedAt)
	assert.True(t, notification.CreatedAt > 0)
	assert.Equal(t, notification.CreatedAt, notification.UpdatedAt)
	assert.Equal(t, localizations, notification.Localizations)
}
