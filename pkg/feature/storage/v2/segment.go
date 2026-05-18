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

	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var (
	ErrSegmentAlreadyExists          = pkgErr.NewErrorAlreadyExists(pkgErr.FeaturePackageName, "segment already exists")
	ErrSegmentNotFound               = pkgErr.NewErrorNotFound(pkgErr.FeaturePackageName, "segment not found", "segment")
	ErrSegmentUnexpectedAffectedRows = pkgErr.NewErrorUnexpectedAffectedRows(
		pkgErr.FeaturePackageName,
		"segment unexpected affected rows",
	)

	ErrInvalidListSegmentsCursor  = errors.New("segment storage: invalid list segments cursor")
	ErrInvalidListSegmentsOrderBy = errors.New("segment storage: invalid list segments order by")
)

type SegmentStorage interface {
	CreateSegment(ctx context.Context, segment *domain.Segment, environmentId string) error
	UpdateSegment(ctx context.Context, segment *domain.Segment, environmentId string) error
	GetSegment(ctx context.Context, id, environmentId string) (*domain.Segment, []string, error)
	ListSegments(
		ctx context.Context,
		params ListSegmentsParams,
	) ([]*proto.Segment, int, int64, map[string][]string, error)
	DeleteSegment(ctx context.Context, id string) error
	// ListAllInUseSegments lists all segments that are in use (referenced by feature flags).
	// Returns lightweight segment info (id, environment_id, updated_at).
	ListAllInUseSegments(ctx context.Context) ([]*InUseSegment, error)
	// ListSegmentUsersBySegment lists all users for a specific segment.
	// This is called per-segment to avoid loading all users in a single query.
	ListSegmentUsersBySegment(ctx context.Context, segmentID, environmentID string) ([]*proto.SegmentUser, error)
}

// InUseSegment represents a segment that is in use by feature flags.
type InUseSegment struct {
	SegmentID     string
	EnvironmentID string
	UpdatedAt     int64
}

// ListSegmentsParams carries list intent for ListSegments without database-specific types.
type ListSegmentsParams struct {
	PageSize       int64
	Cursor         string
	EnvironmentID  string
	Status         *int32
	SearchKeyword  string
	IsInUseStatus  *bool
	OrderBy        proto.ListSegmentsRequest_OrderBy
	OrderDirection proto.ListSegmentsRequest_OrderDirection
}
