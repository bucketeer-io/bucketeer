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
	ErrInvalidListDraftAdminNotificationsCursor = errors.New(
		"notification storage: invalid list draft admin notifications cursor")
	ErrInvalidListDraftAdminNotificationsOrderBy = errors.New(
		"notification storage: invalid list draft admin notifications order by")
)

type NotificationStorage interface {
	CreateNotification(ctx context.Context, notification *domain.Notification) error
	ListDraftAdminNotifications(
		ctx context.Context,
		params ListDraftAdminNotificationsParams,
	) ([]*proto.Notification, int, int64, error)
}

type ListDraftAdminNotificationsParams struct {
	SearchKeyword  string
	OrderBy        proto.ListDraftAdminNotificationsRequest_OrderBy
	OrderDirection proto.ListDraftAdminNotificationsRequest_OrderDirection
	PageSize       int
	Cursor         string
}
