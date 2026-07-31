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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package storage

import (
	"context"
	"errors"

	bkterr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/notification/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
)

var (
	ErrNotificationNotFound = bkterr.NewErrorNotFound(
		bkterr.NotificationPackageName,
		"not found",
		"notification",
	)
	ErrNotificationAlreadyExists = bkterr.NewErrorAlreadyExists(
		bkterr.NotificationPackageName,
		"already exists",
	)
	ErrNotificationAlreadyPublished = bkterr.NewErrorFailedPrecondition(
		bkterr.NotificationPackageName,
		"already published",
	)
	ErrInvalidListDraftAdminNotificationsCursor = errors.New(
		"notification storage: invalid list draft admin notifications cursor")
	ErrInvalidListDraftAdminNotificationsOrderBy = errors.New(
		"notification storage: invalid list draft admin notifications order by")
	ErrInvalidListNotificationsCursor = errors.New(
		"notification storage: invalid list notifications cursor")
	ErrInvalidListNotificationsOrderBy = errors.New(
		"notification storage: invalid list notifications order by")
)

type NotificationStorage interface {
	CreateAdminNotification(ctx context.Context, notification *domain.Notification) error
	GetAdminNotification(ctx context.Context, id string) (*domain.Notification, error)
	UpdateAdminNotification(ctx context.Context, notification *domain.Notification) error
	// PublishAdminNotification persists the publish state transition
	// (status, published_by, published_at, updated_at). The update only
	// applies to a live draft, so concurrent publishes cannot overwrite
	// each other; losing the race returns ErrNotificationAlreadyPublished.
	PublishAdminNotification(ctx context.Context, notification *domain.Notification) error
	// DeleteAdminNotification soft-deletes a notification, recording who
	// deleted it and when; localizations and read markers stay intact.
	DeleteAdminNotification(ctx context.Context, id, lastEditedBy string, updatedAt int64) error
	ListDraftAdminNotifications(
		ctx context.Context,
		params ListDraftAdminNotificationsParams,
	) ([]*proto.Notification, int, int64, error)
	// ListNotifications lists published notifications for a viewer: the
	// localization is resolved to the requested language (falling back to
	// English, then to whichever localization exists) and the viewer's read
	// flag is attached. UNREAD only considers notifications published after
	// the viewer's account was created.
	ListNotifications(
		ctx context.Context,
		params ListNotificationsParams,
	) ([]*proto.Notification, int, int64, error)
}

type ListDraftAdminNotificationsParams struct {
	SearchKeyword  string
	OrderBy        proto.ListDraftAdminNotificationsRequest_OrderBy
	OrderDirection proto.ListDraftAdminNotificationsRequest_OrderDirection
	PageSize       int
	Cursor         string
}

type ListNotificationsParams struct {
	// Viewer identity; read state is keyed by email.
	Email string
	// BCP 47 console language used to resolve the localization.
	Language        string
	SearchKeyword   string
	ReadStatus      proto.ListNotificationsRequest_ReadStatus
	PublishedAtFrom int64
	PublishedAtTo   int64
	OrderBy         proto.ListNotificationsRequest_OrderBy
	OrderDirection  proto.ListNotificationsRequest_OrderDirection
	PageSize        int
	Cursor          string
}
