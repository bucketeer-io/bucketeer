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

func TestUpdateNotification(t *testing.T) {
	t.Parallel()
	notification, err := NewNotification("admin@example.com", []*proto.NotificationLocalization{
		{Language: "en", Title: "Old title", Content: "old content"},
	})
	assert.Nil(t, err)
	createdAt := notification.CreatedAt
	newLocalizations := []*proto.NotificationLocalization{
		{Language: "en", Title: "New title", Content: "new content"},
		{Language: "ja", Title: "新タイトル", Content: "新しい内容"},
	}
	notification.Update("editor@example.com", newLocalizations)
	assert.Equal(t, "admin@example.com", notification.CreatedBy)
	assert.Equal(t, "editor@example.com", notification.LastEditedBy)
	assert.Equal(t, createdAt, notification.CreatedAt)
	assert.True(t, notification.UpdatedAt >= createdAt)
	assert.Equal(t, newLocalizations, notification.Localizations)
}

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

func TestPublishNotification(t *testing.T) {
	t.Parallel()
	notification, err := NewNotification("admin@example.com", []*proto.NotificationLocalization{
		{Language: "en", Title: "New feature", Content: "# New feature"},
	})
	assert.Nil(t, err)
	createdAt := notification.CreatedAt
	notification.Publish("publisher@example.com")
	assert.Equal(t, proto.Notification_PUBLISHED, notification.Status)
	assert.Equal(t, "publisher@example.com", notification.PublishedBy)
	assert.True(t, notification.PublishedAt >= createdAt)
	assert.Equal(t, notification.PublishedAt, notification.UpdatedAt)
	assert.Equal(t, "admin@example.com", notification.CreatedBy)
	assert.Equal(t, "admin@example.com", notification.LastEditedBy)
	assert.Equal(t, createdAt, notification.CreatedAt)
}

func TestResolveLocalization(t *testing.T) {
	t.Parallel()
	en := &proto.NotificationLocalization{Language: "en", Title: "English"}
	ja := &proto.NotificationLocalization{Language: "ja", Title: "日本語"}
	fr := &proto.NotificationLocalization{Language: "fr", Title: "Français"}

	patterns := []struct {
		desc          string
		localizations []*proto.NotificationLocalization
		language      string
		expected      *proto.NotificationLocalization
	}{
		{
			desc:          "nil when empty",
			localizations: nil,
			language:      "en",
			expected:      nil,
		},
		{
			desc:          "exact match",
			localizations: []*proto.NotificationLocalization{en, ja},
			language:      "ja",
			expected:      ja,
		},
		{
			desc:          "fallback to English",
			localizations: []*proto.NotificationLocalization{en, ja},
			language:      "fr",
			expected:      en,
		},
		{
			desc:          "fallback to first by language code",
			localizations: []*proto.NotificationLocalization{ja, fr},
			language:      "en",
			expected:      fr,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expected, ResolveLocalization(p.localizations, p.language))
		})
	}
}
