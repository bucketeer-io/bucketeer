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
package v2

import (
	"context"

	err "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscription/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
)

var (
	ErrAdminSubscriptionAlreadyExists = err.NewErrorAlreadyExists(
		err.NotificationPackageName,
		"admin subscription already exists",
	)
	ErrAdminSubscriptionNotFound = err.NewErrorNotFound(
		err.NotificationPackageName,
		"admin subscription not found",
		"admin_subscription",
	)
	ErrAdminSubscriptionUnexpectedAffectedRows = err.NewErrorUnexpectedAffectedRows(
		err.NotificationPackageName,
		"admin subscription unexpected affected rows",
	)
)

// ListAdminSubscriptionsParams carries list intent without database-specific types.
type ListAdminSubscriptionsParams struct {
	SourceTypes    []proto.Subscription_SourceType
	Disabled       *bool
	SearchKeyword  string
	OrderBy        proto.ListAdminSubscriptionsRequest_OrderBy
	OrderDirection proto.ListAdminSubscriptionsRequest_OrderDirection
	PageSize       int64
	Cursor         string
}

type AdminSubscriptionStorage interface {
	CreateAdminSubscription(ctx context.Context, e *domain.Subscription) error
	UpdateAdminSubscription(ctx context.Context, e *domain.Subscription) error
	DeleteAdminSubscription(ctx context.Context, id string) error
	GetAdminSubscription(ctx context.Context, id string) (*domain.Subscription, error)
	ListAdminSubscriptions(
		ctx context.Context,
		params ListAdminSubscriptionsParams,
	) ([]*proto.Subscription, int, int64, error)
}
