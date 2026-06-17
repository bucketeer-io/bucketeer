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
	"github.com/bucketeer-io/bucketeer/v2/pkg/team/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/team"
)

var (
	ErrTeamNotFound               = bkterr.NewErrorNotFound(bkterr.TeamPackageName, "not found", "team")
	ErrTeamUnexpectedAffectedRows = bkterr.NewErrorUnexpectedAffectedRows(
		bkterr.TeamPackageName,
		"unexpected affected rows",
	)
	ErrInvalidListTeamsCursor  = errors.New("team storage: invalid list teams cursor")
	ErrInvalidListTeamsOrderBy = errors.New("team storage: invalid list teams order by")
)

type TeamStorage interface {
	UpsertTeam(ctx context.Context, team *domain.Team) error
	GetTeam(ctx context.Context, id, organizationID string) (*domain.Team, error)
	GetTeamByName(ctx context.Context, name, organizationID string) (*domain.Team, error)
	ListTeams(
		ctx context.Context,
		params ListTeamsParams,
	) ([]*proto.Team, int, int64, error)
	DeleteTeam(ctx context.Context, id string) error
}

type ListTeamsParams struct {
	OrganizationID string
	SearchKeyword  string
	OrderBy        proto.ListTeamsRequest_OrderBy
	OrderDirection proto.ListTeamsRequest_OrderDirection
	PageSize       int
	Cursor         string
}
