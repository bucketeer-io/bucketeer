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
	"time"

	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
)

type Notification struct {
	*proto.Notification
}

// Update replaces the localizations of a draft and stamps the editor.
func (n *Notification) Update(
	lastEditedBy string,
	localizations []*proto.NotificationLocalization,
) {
	n.LastEditedBy = lastEditedBy
	n.UpdatedAt = time.Now().Unix()
	n.Localizations = localizations
}

// Publish marks a draft as published and stamps the publisher.
func (n *Notification) Publish(publishedBy string) {
	now := time.Now().Unix()
	n.Status = proto.Notification_PUBLISHED
	n.PublishedBy = publishedBy
	n.PublishedAt = now
	n.UpdatedAt = now
}

// ResolveLocalization returns the localization for the given language,
// falling back to English, then to the first localization by language code.
// It returns nil when there are no localizations.
func ResolveLocalization(
	localizations []*proto.NotificationLocalization,
	language string,
) *proto.NotificationLocalization {
	if len(localizations) == 0 {
		return nil
	}
	var english *proto.NotificationLocalization
	for _, l := range localizations {
		if l.Language == language {
			return l
		}
		if l.Language == "en" {
			english = l
		}
	}
	if english != nil {
		return english
	}
	first := localizations[0]
	for _, l := range localizations[1:] {
		if l.Language < first.Language {
			first = l
		}
	}
	return first
}

func NewNotification(
	createdBy string,
	localizations []*proto.NotificationLocalization,
) (*Notification, error) {
	now := time.Now().Unix()
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &Notification{
		Notification: &proto.Notification{
			Id:            id.String(),
			Status:        proto.Notification_DRAFT,
			CreatedBy:     createdBy,
			LastEditedBy:  createdBy,
			CreatedAt:     now,
			UpdatedAt:     now,
			Localizations: localizations,
		},
	}, nil
}
