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

	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

var (
	ErrOrganizationAlreadyExists = pkgErr.NewErrorAlreadyExists(
		pkgErr.EnvironmentPackageName,
		"organization already exists")
	ErrOrganizationNotFound = pkgErr.NewErrorNotFound(
		pkgErr.EnvironmentPackageName,
		"organization not found",
		"organization")
	ErrOrganizationUnexpectedAffectedRows = pkgErr.NewErrorUnexpectedAffectedRows(
		pkgErr.EnvironmentPackageName,
		"organization unexpected affected rows")
)

type OrganizationStorage interface {
	CreateOrganization(ctx context.Context, p *domain.Organization) error
	UpdateOrganization(ctx context.Context, p *domain.Organization) error
	GetOrganization(ctx context.Context, id string) (*domain.Organization, error)
	GetSystemAdminOrganization(ctx context.Context) (*domain.Organization, error)
	ListOrganizations(
		ctx context.Context,
		params ListOrganizationsParams,
	) ([]*proto.Organization, int, int64, error)
}

type ListOrganizationsParams struct {
	Disabled       *bool
	Archived       *bool
	SearchKeyword  string
	OrderBy        proto.ListOrganizationsRequest_OrderBy
	OrderDirection proto.ListOrganizationsRequest_OrderDirection
	PageSize       int
	Cursor         string
}
