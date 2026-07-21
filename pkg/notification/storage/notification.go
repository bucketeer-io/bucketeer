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
	ErrInvalidListDraftNotificationsCursor = errors.New(
		"notification storage: invalid list draft notifications cursor")
	ErrInvalidListDraftNotificationsOrderBy = errors.New(
		"notification storage: invalid list draft notifications order by")
)

type NotificationStorage interface {
	CreateNotification(ctx context.Context, notification *domain.Notification) error
	ListDraftNotifications(
		ctx context.Context,
		params ListDraftNotificationsParams,
	) ([]*proto.Notification, int, int64, error)
}

type ListDraftNotificationsParams struct {
	SearchKeyword  string
	OrderBy        proto.ListDraftNotificationsRequest_OrderBy
	OrderDirection proto.ListDraftNotificationsRequest_OrderDirection
	PageSize       int
	Cursor         string
}
