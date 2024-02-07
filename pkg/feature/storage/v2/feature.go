// Copyright 2024 The Bucketeer Authors.
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

type FeatureStorage interface {
	CreateFeature(ctx context.Context, feature *domain.Feature, environmentNamespace string) error
	UpdateFeature(ctx context.Context, feature *domain.Feature, environmentNamespace string) error
	GetFeature(ctx context.Context, key, environmentNamespace string) (*domain.Feature, error)
	ListFeatures(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Feature, int, int64, error)
	ListFeaturesFilteredByExperiment(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Feature, int, int64, error)
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
	environmentNamespace string,
) error {
	query := `
		INSERT INTO feature (
			id,
			name,
			description,
			enabled,
			archived,
			deleted,
			evaluation_undelayable,
			ttl,
			version,
			created_at,
			updated_at,
			variation_type,
			variations,
			targets,
			rules,
			default_strategy,
			off_variation,
			tags,
			maintainer,
			sampling_seed,
			prerequisites,
			environment_namespace
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?
		)
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
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
		environmentNamespace,
	)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrFeatureAlreadyExists
		}
		return err
	}
	return nil
}

func (s *featureStorage) UpdateFeature(
	ctx context.Context,
	feature *domain.Feature,
	environmentNamespace string,
) error {
	query := `
		UPDATE
			feature
		SET
			name = ?,
			description = ?,
			enabled = ?,
			archived = ?,
			deleted = ?,
			evaluation_undelayable = ?,
			ttl = ?,
			version = ?,
			created_at = ?,
			updated_at = ?,
			variation_type = ?,
			variations = ?,
			targets = ?,
			rules = ?,
			default_strategy = ?,
			off_variation = ?,
			tags = ?,
			maintainer = ?,
			sampling_seed = ?,
			prerequisites = ?
		WHERE
			id = ? AND
			environment_namespace = ?
	`
	result, err := s.qe.ExecContext(
		ctx,
		query,
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
		environmentNamespace,
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
	key, environmentNamespace string,
) (*domain.Feature, error) {
	feature := proto.Feature{}
	query := `
		SELECT
			id,
			name,
			description,
			enabled,
			archived,
			deleted,
			evaluation_undelayable,
			ttl,
			version,
			created_at,
			updated_at,
			variation_type,
			variations,
			targets,
			rules,
			default_strategy,
			off_variation,
			tags,
			maintainer,
			sampling_seed,
			prerequisites
		FROM
			feature
		WHERE
			id = ? AND
			environment_namespace = ?
	`
	err := s.qe.QueryRowContext(
		ctx,
		query,
		key,
		environmentNamespace,
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
		if err == mysql.ErrNoRows {
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
	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			description,
			enabled,
			archived,
			deleted,
			evaluation_undelayable,
			ttl,
			version,
			created_at,
			updated_at,
			variation_type,
			variations,
			targets,
			rules,
			default_strategy,
			off_variation,
			tags,
			maintainer,
			sampling_seed,
			prerequisites
		FROM
			feature
		%s %s %s
		`, whereSQL, orderBySQL, limitOffsetSQL,
	)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	features := make([]*proto.Feature, 0, limit)
	for rows.Next() {
		feature := proto.Feature{}
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
		)
		if err != nil {
			return nil, 0, 0, err
		}
		features = append(features, &feature)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(features)
	var totalCount int64
	countQuery := fmt.Sprintf(`
		SELECT
			COUNT(1)
		FROM
			feature
		%s %s
		`, whereSQL, orderBySQL,
	)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return features, nextOffset, totalCount, nil
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
	query := fmt.Sprintf(`
		SELECT DISTINCT
			feature.id,
			feature.name,
			feature.description,
			feature.enabled,
			feature.archived,
			feature.deleted,
			feature.evaluation_undelayable,
			feature.ttl,
			feature.version,
			feature.created_at,
			feature.updated_at,
			feature.variation_type,
			feature.variations,
			feature.targets,
			feature.rules,
			feature.default_strategy,
			feature.off_variation,
			feature.tags,
			feature.maintainer,
			feature.sampling_seed,
			feature.prerequisites
		FROM
			feature
		LEFT OUTER JOIN
			experiment
		ON
			feature.id = experiment.feature_id AND
			feature.environment_namespace = experiment.environment_namespace
		%s %s %s
		`, whereSQL, orderBySQL, limitOffsetSQL,
	)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	features := make([]*proto.Feature, 0, limit)
	for rows.Next() {
		feature := proto.Feature{}
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
		)
		if err != nil {
			return nil, 0, 0, err
		}
		features = append(features, &feature)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(features)
	var totalCount int64
	countQuery := fmt.Sprintf(`
		SELECT
			COUNT(DISTINCT feature.id)
		FROM
			feature
		LEFT OUTER JOIN
			experiment
		ON
			feature.id = experiment.feature_id AND
			feature.environment_namespace = experiment.environment_namespace
		%s %s
		`, whereSQL, orderBySQL,
	)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return features, nextOffset, totalCount, nil
}
