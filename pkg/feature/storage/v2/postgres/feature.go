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

package postgres

import (
	"context"
	_ "embed"
	"errors"
	"strconv"
	"time"

	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const activeDays = 7 * 24 * time.Hour

var (
	//go:embed sql/feature/create_feature.sql
	createFeatureSQLQuery string
	//go:embed sql/feature/update_feature.sql
	updateFeatureSQLQuery string
	//go:embed sql/feature/select_feature.sql
	selectFeatureSQLQuery string
	//go:embed sql/feature/select_feature_by_version.sql
	selectFeatureByVersionSQLQuery string
	//go:embed sql/feature/select_features.sql
	selectFeaturesSQLQuery string
	//go:embed sql/feature/select_features_by_experiment.sql
	selectFeaturesByExperimentSQLQuery string
	//go:embed sql/feature/select_features_by_environment.sql
	selectFeaturesByEnvironmentSQLQuery string
	//go:embed sql/feature/select_all_environment_features.sql
	selectAllEnvironmentFeaturesSQLQuery string
	//go:embed sql/feature/select_feature_count_by_status.sql
	selectFeatureCountByStatusSQLQuery string
	//go:embed sql/feature/count_features.sql
	countFeatureSQLQuery string
	//go:embed sql/feature/count_features_by_experiment.sql
	countFeaturesByExperimentSQLQuery string
)

type featureStorage struct {
	qe pgstorage.QueryExecer
}

func NewFeatureStorage(qe pgstorage.QueryExecer) v2fs.FeatureStorage {
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
		pgstorage.JSONObject{Val: feature.Variations},
		pgstorage.JSONObject{Val: feature.Targets},
		pgstorage.JSONObject{Val: feature.Rules},
		pgstorage.JSONObject{Val: feature.DefaultStrategy},
		feature.OffVariation,
		pgstorage.JSONObject{Val: feature.Tags},
		feature.Maintainer,
		feature.SamplingSeed,
		pgstorage.JSONObject{Val: feature.Prerequisites},
		environmentID,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrDuplicateEntry) {
			return v2fs.ErrFeatureAlreadyExists
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
		pgstorage.JSONObject{Val: feature.Variations},
		pgstorage.JSONObject{Val: feature.Targets},
		pgstorage.JSONObject{Val: feature.Rules},
		pgstorage.JSONObject{Val: feature.DefaultStrategy},
		feature.OffVariation,
		pgstorage.JSONObject{Val: feature.Tags},
		feature.Maintainer,
		feature.SamplingSeed,
		pgstorage.JSONObject{Val: feature.Prerequisites},
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
		return v2fs.ErrFeatureUnexpectedAffectedRows
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
		&pgstorage.JSONObject{Val: &feature.Variations},
		&pgstorage.JSONObject{Val: &feature.Targets},
		&pgstorage.JSONObject{Val: &feature.Rules},
		&pgstorage.JSONObject{Val: &feature.DefaultStrategy},
		&feature.OffVariation,
		&pgstorage.JSONObject{Val: &feature.Tags},
		&feature.Maintainer,
		&feature.SamplingSeed,
		&pgstorage.JSONObject{Val: &feature.Prerequisites},
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, v2fs.ErrFeatureNotFound
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
		&pgstorage.JSONObject{Val: &feature},
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, v2fs.ErrFeatureNotFound
		}
		return nil, err
	}
	return &domain.Feature{Feature: &feature}, nil
}

func listFeaturesOptionsFromParams(p v2fs.ListFeaturesParams) (*pgstorage.ListOptions, error) {
	filters := []*pgstorage.Filter{
		{
			Column:   "feature.environment_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.EnvironmentID,
		},
	}
	if p.Deleted != nil {
		filters = append(filters, &pgstorage.Filter{
			Column:   "feature.deleted",
			Operator: pgstorage.OperatorEqual,
			Value:    *p.Deleted,
		})
	}
	if p.Maintainer != "" {
		filters = append(filters, &pgstorage.Filter{
			Column:   "feature.maintainer",
			Operator: pgstorage.OperatorEqual,
			Value:    p.Maintainer,
		})
	}
	if p.Enabled != nil {
		filters = append(filters, &pgstorage.Filter{
			Column:   "feature.enabled",
			Operator: pgstorage.OperatorEqual,
			Value:    *p.Enabled,
		})
	}
	if p.Archived != nil {
		filters = append(filters, &pgstorage.Filter{
			Column:   "feature.archived",
			Operator: pgstorage.OperatorEqual,
			Value:    *p.Archived,
		})
	}
	if p.HasFeatureFlagAsRule != nil {
		filters = append(filters, &pgstorage.Filter{
			Column:   "feature.rules::jsonb @> '[{\"clauses\":[{\"operator\":11}]}]'",
			Operator: pgstorage.OperatorEqual,
			Value:    *p.HasFeatureFlagAsRule,
		})
	}

	var jsonFilters []*pgstorage.JSONFilter
	tagValues := make([]interface{}, 0, len(p.Tags))
	for _, tag := range p.Tags {
		tagValues = append(tagValues, tag)
	}
	if len(tagValues) > 0 {
		jsonFilters = append(jsonFilters, &pgstorage.JSONFilter{
			Column: "feature.tags",
			Func:   pgstorage.JSONContainsString,
			Values: tagValues,
		})
	}
	if p.HasPrerequisites != nil {
		if *p.HasPrerequisites {
			jsonFilters = append(jsonFilters, &pgstorage.JSONFilter{
				Column: "feature.prerequisites",
				Func:   pgstorage.JSONLengthGreaterThan,
				Values: []interface{}{"0"},
			})
		} else {
			jsonFilters = append(jsonFilters, &pgstorage.JSONFilter{
				Column: "feature.prerequisites",
				Func:   pgstorage.JSONLengthSmallerThan,
				Values: []interface{}{"1"},
			})
		}
	}

	var nullFilters []*pgstorage.NullFilter
	switch p.Status {
	case proto.FeatureLastUsedInfo_UNKNOWN:
	case proto.FeatureLastUsedInfo_NEW:
		nullFilters = append(nullFilters, &pgstorage.NullFilter{
			Column: "feature_last_used_info.id",
			IsNull: true,
		})
	case proto.FeatureLastUsedInfo_ACTIVE:
		filters = append(filters, &pgstorage.Filter{
			Column:   "feature_last_used_info.last_used_at",
			Operator: pgstorage.OperatorGreaterThanOrEqual,
			Value:    time.Now().Add(-activeDays).Unix(),
		})
	case proto.FeatureLastUsedInfo_NO_ACTIVITY:
		filters = append(filters, &pgstorage.Filter{
			Column:   "feature_last_used_info.last_used_at",
			Operator: pgstorage.OperatorLessThan,
			Value:    time.Now().Add(-activeDays).Unix(),
		})
	}

	var searchQuery *pgstorage.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &pgstorage.SearchQuery{
			Columns: []string{"feature.id", "feature.name", "feature.description"},
			Keyword: p.SearchKeyword,
		}
	}

	var inFilters []*pgstorage.InFilter
	if len(p.IDs) > 0 {
		ids := make([]interface{}, 0, len(p.IDs))
		for _, id := range p.IDs {
			ids = append(ids, id)
		}
		inFilters = append(inFilters, &pgstorage.InFilter{
			Column: "feature.id",
			Values: ids,
		})
	}

	orders, err := listFeaturesOrders(p.OrderBy, p.OrderDirection)
	if err != nil {
		return nil, err
	}

	limit := int(p.PageSize)
	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, v2fs.ErrInvalidListFeaturesCursor
	}

	return &pgstorage.ListOptions{
		Limit:       limit,
		Offset:      offset,
		Filters:     filters,
		JSONFilters: jsonFilters,
		NullFilters: nullFilters,
		InFilters:   inFilters,
		SearchQuery: searchQuery,
		Orders:      orders,
	}, nil
}

func listFeaturesOrders(
	orderBy proto.ListFeaturesRequest_OrderBy,
	orderDirection proto.ListFeaturesRequest_OrderDirection,
) ([]*pgstorage.Order, error) {
	var column string
	switch orderBy {
	case proto.ListFeaturesRequest_DEFAULT,
		proto.ListFeaturesRequest_NAME:
		column = "feature.name"
	case proto.ListFeaturesRequest_CREATED_AT:
		column = "feature.created_at"
	case proto.ListFeaturesRequest_UPDATED_AT:
		column = "feature.updated_at"
	case proto.ListFeaturesRequest_TAGS:
		column = "feature.tags"
	case proto.ListFeaturesRequest_ENABLED:
		column = "feature.enabled"
	case proto.ListFeaturesRequest_AUTO_OPS:
		column = "(COALESCE(auto_ops_counts.progressive_rollout_count, 0) + " +
			"COALESCE(auto_ops_counts.schedule_count, 0) + " +
			"COALESCE(auto_ops_counts.kill_switch_count, 0))"
	default:
		return nil, v2fs.ErrInvalidListFeaturesOrderBy
	}
	direction := pgstorage.OrderDirectionAsc
	if orderDirection == proto.ListFeaturesRequest_DESC {
		direction = pgstorage.OrderDirectionDesc
	}
	return []*pgstorage.Order{pgstorage.NewOrder(column, direction)}, nil
}

func (s *featureStorage) ListFeatures(
	ctx context.Context,
	p v2fs.ListFeaturesParams,
) ([]*proto.Feature, int, int64, error) {
	options, err := listFeaturesOptionsFromParams(p)
	if err != nil {
		return nil, 0, 0, err
	}
	query, whereArgs := pgstorage.ConstructQueryAndWhereArgs(selectFeaturesSQLQuery, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	var limit, offset int
	if options != nil {
		limit = options.Limit
		offset = options.Offset
	}
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
			&pgstorage.JSONObject{Val: &feature.Variations},
			&pgstorage.JSONObject{Val: &feature.Targets},
			&pgstorage.JSONObject{Val: &feature.Rules},
			&pgstorage.JSONObject{Val: &feature.DefaultStrategy},
			&feature.OffVariation,
			&pgstorage.JSONObject{Val: &feature.Tags},
			&feature.Maintainer,
			&feature.SamplingSeed,
			&pgstorage.JSONObject{Val: &feature.Prerequisites},
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
		if lastUsedInfo.FeatureId != "" {
			feature.LastUsedInfo = &lastUsedInfo
		}
		features = append(features, &feature)
	}
	if rows.Err() != nil {
		return nil, 0, 0, rows.Err()
	}
	nextOffset := offset + len(features)
	var totalCount int64
	countQuery, countWhereArgs := pgstorage.ConstructCountQuery(countFeatureSQLQuery, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
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
	p v2fs.ListFeaturesFilteredByExperimentParams,
) ([]*proto.Feature, int, int64, error) {
	options, err := listFeaturesOptionsFromParams(p.ListFeaturesParams)
	if err != nil {
		return nil, 0, 0, err
	}
	if options.NullFilters == nil {
		options.NullFilters = make([]*pgstorage.NullFilter, 0, 1)
	}
	options.NullFilters = append(options.NullFilters, &pgstorage.NullFilter{
		Column: "experiment.id",
		IsNull: !p.HasExperiment,
	})

	query, whereArgs := pgstorage.ConstructQueryAndWhereArgs(selectFeaturesByExperimentSQLQuery, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	var limit, offset int
	if options != nil {
		limit = options.Limit
		offset = options.Offset
	}
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
			&pgstorage.JSONObject{Val: &feature.Variations},
			&pgstorage.JSONObject{Val: &feature.Targets},
			&pgstorage.JSONObject{Val: &feature.Rules},
			&pgstorage.JSONObject{Val: &feature.DefaultStrategy},
			&feature.OffVariation,
			&pgstorage.JSONObject{Val: &feature.Tags},
			&feature.Maintainer,
			&feature.SamplingSeed,
			&pgstorage.JSONObject{Val: &feature.Prerequisites},
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
		if lastUsedInfo.FeatureId != "" {
			feature.LastUsedInfo = &lastUsedInfo
		}
		features = append(features, &feature)
	}
	if rows.Err() != nil {
		return nil, 0, 0, rows.Err()
	}
	nextOffset := offset + len(features)
	var totalCount int64
	countQuery, countWhereArgs := pgstorage.ConstructCountQuery(countFeaturesByExperimentSQLQuery, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return features, nextOffset, totalCount, nil
}

func (s *featureStorage) ListFeaturesByEnvironment(
	ctx context.Context,
	environmentID string,
) ([]*proto.Feature, error) {
	rows, err := s.qe.QueryContext(ctx, selectFeaturesByEnvironmentSQLQuery, environmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	features := make([]*proto.Feature, 0)
	for rows.Next() {
		feature := proto.Feature{}
		err := rows.Scan(
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
			&pgstorage.JSONObject{Val: &feature.Variations},
			&pgstorage.JSONObject{Val: &feature.Targets},
			&pgstorage.JSONObject{Val: &feature.Rules},
			&pgstorage.JSONObject{Val: &feature.DefaultStrategy},
			&feature.OffVariation,
			&pgstorage.JSONObject{Val: &feature.Tags},
			&feature.Maintainer,
			&feature.SamplingSeed,
			&pgstorage.JSONObject{Val: &feature.Prerequisites},
		)
		if err != nil {
			return nil, err
		}
		features = append(features, &feature)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return features, nil
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
			&pgstorage.JSONObject{Val: &feature.Variations},
			&pgstorage.JSONObject{Val: &feature.Targets},
			&pgstorage.JSONObject{Val: &feature.Rules},
			&pgstorage.JSONObject{Val: &feature.DefaultStrategy},
			&feature.OffVariation,
			&pgstorage.JSONObject{Val: &feature.Tags},
			&feature.Maintainer,
			&feature.SamplingSeed,
			&pgstorage.JSONObject{Val: &feature.Prerequisites},
		)
		if err != nil {
			return nil, err
		}
		envFeatures[envID] = append(envFeatures[envID], &feature)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
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
