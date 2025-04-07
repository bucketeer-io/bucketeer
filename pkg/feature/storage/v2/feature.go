// Copyright 2025 The Bucketeer Authors.
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
	_ "embed"
	"errors"
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	ErrFeatureAlreadyExists          = errors.New("feature: already exists")
	ErrFeatureNotFound               = errors.New("feature: not found")
	ErrFeatureUnexpectedAffectedRows = errors.New("feature: unexpected affected rows")
)

var (
	//go:embed sql/feature/create_feature.sql
	createFeatureSQLQuery string
	//go:embed sql/feature/update_feature.sql
	updateFeatureSQLQuery string
	//go:embed sql/feature/select_all_environment_features.sql
	selectAllEnvironmentFeaturesSQLQuery string
	//go:embed sql/feature/select_features.sql
	selectFeaturesSQLQuery string
	//go:embed sql/feature/select_features_by_experiment.sql
	selectFeaturesByExperimentSQLQuery string
	//go:embed sql/feature/select_feature_count_by_status.sql
	selectFeatureCountByStatusSQLQuery string
	//go:embed sql/feature/count_features.sql
	countFeatureSQLQuery string
	//go:embed sql/feature/count_features_by_experiment.sql
	countFeaturesByExperimentSQLQuery string
	//go:embed sql/feature/select_feature.sql
	selectFeatureSQLQuery string
	//go:embed sql/feature/select_feature_by_version.sql
	selectFeatureByVersionSQLQuery string
)

type FeatureStorage interface {
	CreateFeature(ctx context.Context, feature *domain.Feature, environmentID string) error
	UpdateFeature(ctx context.Context, feature *domain.Feature, environmentID string) error
	GetFeature(ctx context.Context, id, environmentID string) (*domain.Feature, error)
	GetFeatureByVersion(ctx context.Context, id string, version int32, environmentID string) (*domain.Feature, error)
	ListFeatures(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Feature, int, int64, error)
	GetFeatureSummary(
		ctx context.Context,
		environmentID string,
	) (*proto.FeatureSummary, error)
	ListFeaturesFilteredByExperiment(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Feature, int, int64, error)
	ListAllEnvironmentFeatures(
		ctx context.Context,
	) ([]*proto.EnvironmentFeature, error)
}

type featureStorage struct {
	qe mysql.QueryExecer
}

func NewFeatureStorage(qe mysql.QueryExecer) FeatureStorage {
	return &featureStorage{qe: qe}
}

func (s *featureStorage) CreateFeature(
	ctx context.Context,
	feature *domain.Feature,
	environmentID string,
) error {
	_, err := s.qe.ExecContext(
		ctx,
		createFeatureSQLQuery,
		feature.Id,
		feature.Name,
		feature.Description,
		feature.Enabled,
		feature.Archived,
		feature.Deleted,
		feature.EvaluationUndelayable,
		feature.Ttl,
		feature.Version,
		feature.CreatedAt,
		feature.UpdatedAt,
		int32(feature.VariationType),
		mysql.JSONObject{Val: feature.Variations},
		mysql.JSONObject{Val: feature.Targets},
		mysql.JSONObject{Val: feature.Rules},
		mysql.JSONObject{Val: feature.DefaultStrategy},
		feature.OffVariation,
		mysql.JSONObject{Val: feature.Tags},
		feature.Maintainer,
		feature.SamplingSeed,
		mysql.JSONObject{Val: feature.Prerequisites},
		environmentID,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrDuplicateEntry) {
			return ErrFeatureAlreadyExists
		}
		return err
	}
	return nil
}

func (s *featureStorage) UpdateFeature(
	ctx context.Context,
	feature *domain.Feature,
	environmentID string,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateFeatureSQLQuery,
		feature.Name,
		feature.Description,
		feature.Enabled,
		feature.Archived,
		feature.Deleted,
		feature.EvaluationUndelayable,
		feature.Ttl,
		feature.Version,
		feature.CreatedAt,
		feature.UpdatedAt,
		int32(feature.VariationType),
		mysql.JSONObject{Val: feature.Variations},
		mysql.JSONObject{Val: feature.Targets},
		mysql.JSONObject{Val: feature.Rules},
		mysql.JSONObject{Val: feature.DefaultStrategy},
		feature.OffVariation,
		mysql.JSONObject{Val: feature.Tags},
		feature.Maintainer,
		feature.SamplingSeed,
		mysql.JSONObject{Val: feature.Prerequisites},
		feature.Id,
		environmentID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrFeatureUnexpectedAffectedRows
	}
	return nil
}

func (s *featureStorage) GetFeature(
	ctx context.Context,
	id, environmentID string,
) (*domain.Feature, error) {
	feature := proto.Feature{}
	feature.AutoOpsSummary = &proto.AutoOpsSummary{}
	err := s.qe.QueryRowContext(
		ctx,
		selectFeatureSQLQuery,
		id,
		environmentID,
	).Scan(
		&feature.Id,
		&feature.Name,
		&feature.Description,
		&feature.Enabled,
		&feature.Archived,
		&feature.Deleted,
		&feature.EvaluationUndelayable,
		&feature.Ttl,
		&feature.Version,
		&feature.CreatedAt,
		&feature.UpdatedAt,
		&feature.VariationType,
		&mysql.JSONObject{Val: &feature.Variations},
		&mysql.JSONObject{Val: &feature.Targets},
		&mysql.JSONObject{Val: &feature.Rules},
		&mysql.JSONObject{Val: &feature.DefaultStrategy},
		&feature.OffVariation,
		&mysql.JSONObject{Val: &feature.Tags},
		&feature.Maintainer,
		&feature.SamplingSeed,
		&mysql.JSONObject{Val: &feature.Prerequisites},
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrFeatureNotFound
		}
		return nil, err
	}
	return &domain.Feature{Feature: &feature}, nil
}

func (s *featureStorage) GetFeatureByVersion(
	ctx context.Context,
	id string, version int32, environmentID string,
) (*domain.Feature, error) {
	feature := proto.Feature{}
	feature.AutoOpsSummary = &proto.AutoOpsSummary{}
	err := s.qe.QueryRowContext(
		ctx,
		selectFeatureByVersionSQLQuery,
		environmentID,
		id,
		version,
	).Scan(
		&mysql.JSONObject{Val: &feature},
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrFeatureNotFound
		}
		return nil, err
	}
	return &domain.Feature{Feature: &feature}, nil
}

func (s *featureStorage) ListFeatures(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.Feature, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(selectFeaturesSQLQuery, whereSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	features := make([]*proto.Feature, 0, limit)
	for rows.Next() {
		feature := proto.Feature{}
		feature.AutoOpsSummary = &proto.AutoOpsSummary{}
		lastUsedInfo := proto.FeatureLastUsedInfo{}
		err := rows.Scan(
			&feature.Id,
			&feature.Name,
			&feature.Description,
			&feature.Enabled,
			&feature.Archived,
			&feature.Deleted,
			&feature.EvaluationUndelayable,
			&feature.Ttl,
			&feature.Version,
			&feature.CreatedAt,
			&feature.UpdatedAt,
			&feature.VariationType,
			&mysql.JSONObject{Val: &feature.Variations},
			&mysql.JSONObject{Val: &feature.Targets},
			&mysql.JSONObject{Val: &feature.Rules},
			&mysql.JSONObject{Val: &feature.DefaultStrategy},
			&feature.OffVariation,
			&mysql.JSONObject{Val: &feature.Tags},
			&feature.Maintainer,
			&feature.SamplingSeed,
			&mysql.JSONObject{Val: &feature.Prerequisites},
			&feature.AutoOpsSummary.ProgressiveRolloutCount,
			&feature.AutoOpsSummary.ScheduleCount,
			&feature.AutoOpsSummary.KillSwitchCount,
			&lastUsedInfo.FeatureId,
			&lastUsedInfo.Version,
			&lastUsedInfo.LastUsedAt,
			&lastUsedInfo.CreatedAt,
			&lastUsedInfo.ClientOldestVersion,
			&lastUsedInfo.ClientLatestVersion,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		// Flags that haven't been evaluated yet won't have the status info.
		if lastUsedInfo.FeatureId != "" {
			feature.LastUsedInfo = &lastUsedInfo
		}
		features = append(features, &feature)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(features)
	var totalCount int64
	countQuery := fmt.Sprintf(countFeatureSQLQuery, whereSQL)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return features, nextOffset, totalCount, nil
}

func (s *featureStorage) GetFeatureSummary(
	ctx context.Context,
	environmentID string,
) (*proto.FeatureSummary, error) {
	var countByStatus proto.FeatureSummary
	err := s.qe.QueryRowContext(ctx, selectFeatureCountByStatusSQLQuery, environmentID).Scan(
		&countByStatus.Total,
		&countByStatus.Active,
		&countByStatus.Inactive,
	)
	if err != nil {
		return nil, err
	}
	return &countByStatus, nil
}

func (s *featureStorage) ListFeaturesFilteredByExperiment(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.Feature, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(selectFeaturesByExperimentSQLQuery, whereSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	features := make([]*proto.Feature, 0, limit)
	for rows.Next() {
		feature := proto.Feature{}
		feature.AutoOpsSummary = &proto.AutoOpsSummary{}
		lastUsedInfo := proto.FeatureLastUsedInfo{}
		err := rows.Scan(
			&feature.Id,
			&feature.Name,
			&feature.Description,
			&feature.Enabled,
			&feature.Archived,
			&feature.Deleted,
			&feature.EvaluationUndelayable,
			&feature.Ttl,
			&feature.Version,
			&feature.CreatedAt,
			&feature.UpdatedAt,
			&feature.VariationType,
			&mysql.JSONObject{Val: &feature.Variations},
			&mysql.JSONObject{Val: &feature.Targets},
			&mysql.JSONObject{Val: &feature.Rules},
			&mysql.JSONObject{Val: &feature.DefaultStrategy},
			&feature.OffVariation,
			&mysql.JSONObject{Val: &feature.Tags},
			&feature.Maintainer,
			&feature.SamplingSeed,
			&mysql.JSONObject{Val: &feature.Prerequisites},
			&feature.AutoOpsSummary.ProgressiveRolloutCount,
			&feature.AutoOpsSummary.ScheduleCount,
			&feature.AutoOpsSummary.KillSwitchCount,
			&lastUsedInfo.FeatureId,
			&lastUsedInfo.Version,
			&lastUsedInfo.LastUsedAt,
			&lastUsedInfo.CreatedAt,
			&lastUsedInfo.ClientOldestVersion,
			&lastUsedInfo.ClientLatestVersion,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		// Flags that haven't been evaluated yet won't have the status info.
		if lastUsedInfo.FeatureId != "" {
			feature.LastUsedInfo = &lastUsedInfo
		}
		features = append(features, &feature)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(features)
	var totalCount int64
	countQuery := fmt.Sprintf(countFeaturesByExperimentSQLQuery, whereSQL)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return features, nextOffset, totalCount, nil
}

func (s *featureStorage) ListAllEnvironmentFeatures(
	ctx context.Context,
) ([]*proto.EnvironmentFeature, error) {
	rows, err := s.qe.QueryContext(ctx, selectAllEnvironmentFeaturesSQLQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	envFeatures := map[string][]*proto.Feature{}
	for rows.Next() {
		feature := proto.Feature{}
		var envID string
		err := rows.Scan(
			&envID,
			// Feature columns
			&feature.Id,
			&feature.Name,
			&feature.Description,
			&feature.Enabled,
			&feature.Archived,
			&feature.Deleted,
			&feature.Version,
			&feature.CreatedAt,
			&feature.UpdatedAt,
			&feature.VariationType,
			&mysql.JSONObject{Val: &feature.Variations},
			&mysql.JSONObject{Val: &feature.Targets},
			&mysql.JSONObject{Val: &feature.Rules},
			&mysql.JSONObject{Val: &feature.DefaultStrategy},
			&feature.OffVariation,
			&mysql.JSONObject{Val: &feature.Tags},
			&feature.Maintainer,
			&feature.SamplingSeed,
			&mysql.JSONObject{Val: &feature.Prerequisites},
		)
		if err != nil {
			return nil, err
		}
		envFeatures[envID] = append(envFeatures[envID], &feature)
	}
	if rows.Err() != nil {
		return nil, err
	}
	envFts := make([]*proto.EnvironmentFeature, 0, len(envFeatures))
	for key, fts := range envFeatures {
		envFeature := &proto.EnvironmentFeature{
			EnvironmentId: key,
			Features:      fts,
		}
		envFts = append(envFts, envFeature)
	}
	return envFts, nil
}
