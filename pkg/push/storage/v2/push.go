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
	"errors"

	err "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/push/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

var (
	ErrPushAlreadyExists          = err.NewErrorAlreadyExists(err.PushPackageName, "push already exists")
	ErrPushNotFound               = err.NewErrorNotFound(err.PushPackageName, "push not found", "push")
	ErrPushUnexpectedAffectedRows = err.NewErrorUnexpectedAffectedRows(
		err.PushPackageName,
		"push unexpected affected rows",
	)

	ErrInvalidListPushesCursor = errors.New("push storage: invalid list pushes cursor")
)

// ListPushesParams carries list intent for ListPushes without database-specific types.
type ListPushesParams struct {
	// PageSize is row limit; use database.QueryNoLimit for an uncapped list.
	PageSize       int64
	Cursor         string
	OrganizationID string
	EnvironmentIDs []string
	SearchKeyword  string
	Disabled       *bool
	// Deleted filters on push.deleted (e.g. false for active pushes only).
	Deleted        bool
	OrderBy        *proto.ListPushesRequest_OrderBy
	OrderDirection proto.ListPushesRequest_OrderDirection
}

type PushStorage interface {
	CreatePush(ctx context.Context, e *domain.Push, environmentId string) error
	UpdatePush(ctx context.Context, e *domain.Push, environmentId string) error
	GetPush(ctx context.Context, id, environmentId string) (*domain.Push, error)
	ListPushes(ctx context.Context, p ListPushesParams) ([]*proto.Push, int, int64, error)
	DeletePush(ctx context.Context, id, environmentId string) error
}
