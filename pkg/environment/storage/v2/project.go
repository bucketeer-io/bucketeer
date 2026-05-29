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
	ErrProjectAlreadyExists = pkgErr.NewErrorAlreadyExists(
		pkgErr.EnvironmentPackageName,
		"project already exists")
	ErrProjectNotFound = pkgErr.NewErrorNotFound(
		pkgErr.EnvironmentPackageName,
		"project not found", "project")
	ErrProjectUnexpectedAffectedRows = pkgErr.NewErrorUnexpectedAffectedRows(
		pkgErr.EnvironmentPackageName,
		"project unexpected affected rows")
)

type ProjectStorage interface {
	CreateProject(ctx context.Context, p *domain.Project) error
	UpdateProject(ctx context.Context, p *domain.Project) error
	GetProject(ctx context.Context, id string) (*domain.Project, error)
	GetTrialProjectByEmail(
		ctx context.Context,
		email string,
		disabled, trial bool,
	) (*domain.Project, error)
	ListProjects(
		ctx context.Context,
		params ListProjectsParams,
	) ([]*proto.Project, int, int64, error)
}

type ListProjectsParams struct {
	OrganizationIDs []string
	OrganizationID  string
	Disabled        *bool
	SearchKeyword   string
	OrderBy         proto.ListProjectsV2Request_OrderBy
	OrderDirection  proto.ListProjectsV2Request_OrderDirection
	PageSize        int
	Cursor          string
}
