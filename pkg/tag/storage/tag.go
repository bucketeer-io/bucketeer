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

	pkgerr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/tag/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/tag"
)

var (
	ErrTagNotFound               = pkgerr.NewErrorNotFound(pkgerr.TagPackageName, "tag not found", "tag")
	ErrTagUnexpectedAffectedRows = pkgerr.NewErrorUnexpectedAffectedRows(
		pkgerr.TagPackageName,
		"tag unexpected affected rows",
	)
	ErrInvalidListTagsCursor  = errors.New("tag storage: invalid list tags cursor")
	ErrInvalidListTagsOrderBy = errors.New("tag storage: invalid list tags order by")
)

type TagStorage interface {
	UpsertTag(ctx context.Context, tag *domain.Tag) error
	GetTag(ctx context.Context, id, environmentId string) (*domain.Tag, error)
	GetTagByName(
		ctx context.Context,
		name, environmentId string,
		entityType proto.Tag_EntityType,
	) (*domain.Tag, error)
	ListTags(
		ctx context.Context,
		params ListTagsParams,
	) ([]*proto.Tag, int, int64, error)
	ListAllEnvironmentTags(ctx context.Context) ([]*proto.EnvironmentTag, error)
	DeleteTag(ctx context.Context, id string) error
}

type ListTagsParams struct {
	EnvironmentID  string
	OrganizationID string
	EnvironmentIDs []string
	EntityType     proto.Tag_EntityType
	SearchKeyword  string
	OrderBy        proto.ListTagsRequest_OrderBy
	OrderDirection proto.ListTagsRequest_OrderDirection
	PageSize       int
	Cursor         string
}
