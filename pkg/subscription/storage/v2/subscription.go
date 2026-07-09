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
	ErrSubscriptionAlreadyExists = err.NewErrorAlreadyExists(
		err.SubscriptionPackageName,
		"subscription already exists",
	)
	ErrSubscriptionNotFound = err.NewErrorNotFound(
		err.SubscriptionPackageName,
		"subscription not found",
		"subscription",
	)
	ErrSubscriptionUnexpectedAffectedRows = err.NewErrorUnexpectedAffectedRows(
		err.SubscriptionPackageName,
		"subscription unexpected affected rows",
	)
)

// ListSubscriptionsParams carries list intent without database-specific types.
type ListSubscriptionsParams struct {
	OrganizationID string
	EnvironmentIDs []string
	SourceTypes    []proto.Subscription_SourceType
	Disabled       *bool
	SearchKeyword  string
	OrderBy        proto.ListSubscriptionsRequest_OrderBy
	OrderDirection proto.ListSubscriptionsRequest_OrderDirection
	PageSize       int64
	Cursor         string
}

type SubscriptionStorage interface {
	CreateSubscription(ctx context.Context, e *domain.Subscription, environmentId string) error
	UpdateSubscription(ctx context.Context, e *domain.Subscription, environmentId string) error
	DeleteSubscription(ctx context.Context, id, environmentId string) error
	GetSubscription(ctx context.Context, id, environmentId string) (*domain.Subscription, error)
	ListSubscriptions(
		ctx context.Context,
		params ListSubscriptionsParams,
	) ([]*proto.Subscription, int, int64, error)
}
