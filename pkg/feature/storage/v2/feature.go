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
	ErrFeatureAlreadyExists          = pkgErr.NewErrorAlreadyExists(pkgErr.FeaturePackageName, "feature already exists")
	ErrFeatureNotFound               = pkgErr.NewErrorNotFound(pkgErr.FeaturePackageName, "feature not found", "feature")
	ErrFeatureUnexpectedAffectedRows = pkgErr.NewErrorUnexpectedAffectedRows(
		pkgErr.FeaturePackageName,
		"feature unexpected affected rows",
	)

	ErrInvalidListFeaturesCursor  = errors.New("feature storage: invalid list features cursor")
	ErrInvalidListFeaturesOrderBy = errors.New("feature storage: invalid list features order by")
)

// ListFeaturesParams carries list intent for ListFeatures without database-specific types.
type ListFeaturesParams struct {
	// PageSize is row limit; use database.QueryNoLimit for an uncapped list.
	PageSize             int64
	Cursor               string
	EnvironmentID        string
	IDs                  []string
	Tags                 []string
	Maintainer           string
	Enabled              *bool
	Archived             *bool
	Deleted              *bool
	HasPrerequisites     *bool
	HasFeatureFlagAsRule *bool
	HasAutoOps           *bool
	SearchKeyword        string
	Status               proto.FeatureLastUsedInfo_Status
	OrderBy              proto.ListFeaturesRequest_OrderBy
	OrderDirection       proto.ListFeaturesRequest_OrderDirection
}

// ListFeaturesFilteredByExperimentParams extends ListFeaturesParams with experiment filtering.
type ListFeaturesFilteredByExperimentParams struct {
	ListFeaturesParams
	HasExperiment bool
}

type FeatureStorage interface {
	CreateFeature(ctx context.Context, feature *domain.Feature, environmentID string) error
	UpdateFeature(ctx context.Context, feature *domain.Feature, environmentID string) error
	GetFeature(ctx context.Context, id, environmentID string) (*domain.Feature, error)
	GetFeatureByVersion(ctx context.Context, id string, version int32, environmentID string) (*domain.Feature, error)
	ListFeatures(ctx context.Context, p ListFeaturesParams) ([]*proto.Feature, int, int64, error)
	GetFeatureSummary(
		ctx context.Context,
		environmentID string,
	) (*proto.FeatureSummary, error)
	ListFeaturesFilteredByExperiment(
		ctx context.Context,
		p ListFeaturesFilteredByExperimentParams,
	) ([]*proto.Feature, int, int64, error)
	// ListFeaturesByEnvironment lists all non-deleted features for a specific environment.
	// This is more efficient than ListAllEnvironmentFeatures when only one environment is needed.
	ListFeaturesByEnvironment(ctx context.Context, environmentID string) ([]*proto.Feature, error)
	ListAllEnvironmentFeatures(
		ctx context.Context,
	) ([]*proto.EnvironmentFeature, error)
}
