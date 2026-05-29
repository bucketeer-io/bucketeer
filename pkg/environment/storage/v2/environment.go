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

	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

var (
	ErrEnvironmentAlreadyExists = pkgErr.NewErrorAlreadyExists(
		pkgErr.EnvironmentPackageName,
		"environment already exists")
	ErrEnvironmentNotFound = pkgErr.NewErrorNotFound(
		pkgErr.EnvironmentPackageName,
		"environment not found",
		"environment")
	ErrEnvironmentUnexpectedAffectedRows = pkgErr.NewErrorUnexpectedAffectedRows(
		pkgErr.EnvironmentPackageName,
		"environment unexpected affected rows")
)

// Shared list-query errors returned by EnvironmentStorage, OrganizationStorage,
// and ProjectStorage implementations.
var (
	ErrInvalidOrderBy = errors.New("environment/storage/v2: invalid order by")
	ErrInvalidCursor  = errors.New("environment/storage/v2: invalid cursor")
)

type EnvironmentStorage interface {
	CreateEnvironmentV2(ctx context.Context, e *domain.EnvironmentV2) error
	UpdateEnvironmentV2(ctx context.Context, e *domain.EnvironmentV2) error
	GetEnvironmentV2(ctx context.Context, id string) (*domain.EnvironmentV2, error)
	ListEnvironmentsV2(
		ctx context.Context,
		params ListEnvironmentsV2Params,
	) ([]*proto.EnvironmentV2, int, int64, error)
	ListAutoArchiveEnabledEnvironments(ctx context.Context) ([]*domain.EnvironmentV2, error)
}

type ListEnvironmentsV2Params struct {
	ProjectID      string
	OrganizationID string
	Archived       *bool
	SearchKeyword  string
	OrderBy        proto.ListEnvironmentsV2Request_OrderBy
	OrderDirection proto.ListEnvironmentsV2Request_OrderDirection
	PageSize       int
	Cursor         string
}
