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

package api

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	experimentdomain "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	btproto "github.com/bucketeer-io/bucketeer/v2/proto/batch"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

const (
	getMultiChunkSize = 1000
	listRequestSize   = 500
	// after 7 days without request, the feature is considered as no activity
	activeDays = 7 * 24 * time.Hour
)

var errEvaluationNotFound = status.Error(codes.NotFound, "feature: evaluation not found")

func (s *FeatureService) GetFeature(
	ctx context.Context,
	req *featureproto.GetFeatureRequest,
) (*featureproto.GetFeatureResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateGetFeatureRequest(req); err != nil {
		return nil, err
	}
	featureStorage := v2fs.NewFeatureStorage(s.mysqlClient)
	var feature *domain.Feature
	if req.FeatureVersion != nil {
		feature, err = featureStorage.GetFeatureByVersion(
			ctx,
			req.Id,
			req.FeatureVersion.Value,
			req.EnvironmentId,
		)
	} else {
		feature, err = featureStorage.GetFeature(ctx, req.Id, req.EnvironmentId)
	}
	if err != nil {
		if errors.Is(err, v2fs.ErrFeatureNotFound) {
			return nil, statusFeatureNotFound.Err()
		}
		s.logger.Error(
			"Failed to get feature",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	// TEMPORARY: Clean up any orphaned variation references before returning to UI
	// Clean up any orphaned variation references before returning to UI
	// This prevents the UI from seeing corrupted data and sending it back in update requests
	cleanupResult := feature.CleanupOrphanedVariationReferences()
	if cleanupResult.Changed {
		s.logger.Warn(
			"Cleaned up orphaned variation references in feature during get (temporary migration)",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("featureId", req.Id),
				zap.String("environmentId", req.EnvironmentId),
				zap.Int("orphanedTargets", cleanupResult.OrphanedTargets),
				zap.Int("orphanedRules", cleanupResult.OrphanedRules),
				zap.Int("orphanedDefault", cleanupResult.OrphanedDefault),
				zap.Bool("orphanedOffVar", cleanupResult.OrphanedOffVar),
				zap.Strings("orphanedVariationIDs", cleanupResult.OrphanedVariationIDs),
			)...,
		)
	}

	// TEMPORARY: Ensure all variations are present in rollout strategies during reads
	// This fixes existing data corruption from the historical AddVariation bug
	// TODO: Remove this after DB migration is complete
	migrationResult := feature.EnsureVariationsInStrategies()
	if migrationResult.Changed {
		s.logger.Warn(
			"Added missing variations to rollout strategies in feature during get (temporary migration)",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("featureId", req.Id),
				zap.String("environmentId", req.EnvironmentId),
				zap.Int("addedToRules", migrationResult.AddedToRules),
				zap.Int("addedToDefault", migrationResult.AddedToDefault),
				zap.Strings("processedVariationIDs", migrationResult.AddedVariationIDs),
			)...,
		)
	}

	if err := s.setLastUsedInfosToFeatureByChunk(
		ctx,
		[]*featureproto.Feature{feature.Feature},
		req.EnvironmentId,
	); err != nil {
		return nil, err
	}
	return &featureproto.GetFeatureResponse{Feature: feature.Feature}, nil
}

func (s *FeatureService) GetFeatures(
	ctx context.Context,
	req *featureproto.GetFeaturesRequest,
) (*featureproto.GetFeaturesResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateGetFeaturesRequest(req); err != nil {
		return nil, err
	}
	filters := []*mysql.FilterV2{
		{
			Column:   "feature.environment_id",
			Operator: mysql.OperatorEqual,
			Value:    req.EnvironmentId,
		},
	}
	ids := make([]interface{}, 0, len(req.Ids))
	for _, id := range req.Ids {
		ids = append(ids, id)
	}
	var inFilters []*mysql.InFilter
	if len(ids) > 0 {
		inFilters = append(inFilters, &mysql.InFilter{
			Column: "feature.id",
			Values: ids,
		})
	}
	featureStorage := v2fs.NewFeatureStorage(s.mysqlClient)
	options := &mysql.ListOptions{
		Filters:     filters,
		Orders:      nil,
		JSONFilters: nil,
		NullFilters: nil,
		InFilters:   inFilters,
		SearchQuery: nil,
		Limit:       mysql.QueryNoLimit,
		Offset:      mysql.QueryNoOffset,
	}
	features, _, _, err := featureStorage.ListFeatures(ctx, options)
	if err != nil {
		s.logger.Error(
			"Failed to get feature",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	// TEMPORARY: Clean up any orphaned variation references before returning to UI
	// Clean up any orphaned variation references before returning to UI
	// This prevents the UI from seeing corrupted data and sending it back in update requests
	for _, f := range features {
		domainFeature := &domain.Feature{Feature: f}
		cleanupResult := domainFeature.CleanupOrphanedVariationReferences()
		if cleanupResult.Changed {
			s.logger.Warn(
				"Cleaned up orphaned variation references in feature during get multiple (temporary migration)",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.String("featureId", f.Id),
					zap.String("environmentId", req.EnvironmentId),
					zap.Int("orphanedTargets", cleanupResult.OrphanedTargets),
					zap.Int("orphanedRules", cleanupResult.OrphanedRules),
					zap.Int("orphanedDefault", cleanupResult.OrphanedDefault),
					zap.Bool("orphanedOffVar", cleanupResult.OrphanedOffVar),
					zap.Strings("orphanedVariationIDs", cleanupResult.OrphanedVariationIDs),
				)...,
			)
		}

		// TEMPORARY: Ensure all variations are present in rollout strategies during reads
		// This fixes existing data corruption from the historical AddVariation bug
		// TODO: Remove this after DB migration is complete
		migrationResult := domainFeature.EnsureVariationsInStrategies()
		if migrationResult.Changed {
			s.logger.Warn(
				"Added missing variations to rollout strategies in feature during get multiple (temporary migration)",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.String("featureId", f.Id),
					zap.String("environmentId", req.EnvironmentId),
					zap.Int("addedToRules", migrationResult.AddedToRules),
					zap.Int("addedToDefault", migrationResult.AddedToDefault),
					zap.Strings("processedVariationIDs", migrationResult.AddedVariationIDs),
				)...,
			)
		}
	}

	return &featureproto.GetFeaturesResponse{Features: features}, nil
}

func (s *FeatureService) ListFeatures(
	ctx context.Context,
	req *featureproto.ListFeaturesRequest,
) (*featureproto.ListFeaturesResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	var features []*featureproto.Feature
	var cursor string
	var totalCount int64
	if req.HasExperiment == nil {
		features, cursor, totalCount, err = s.listFeatures(
			ctx,
			req.PageSize,
			req.Cursor,
			req.Tags,
			req.Maintainer,
			req.Enabled,
			req.Archived,
			req.HasPrerequisites,
			req.HasFeatureFlagAsRule,
			req.SearchKeyword,
			req.Status,
			req.OrderBy,
			req.OrderDirection,
			req.EnvironmentId,
		)
	} else {
		features, cursor, totalCount, err = s.listFeaturesFilteredByExperiment(
			ctx,
			req.PageSize,
			req.Cursor,
			req.Tags,
			req.Maintainer,
			req.Enabled,
			req.Archived,
			req.HasPrerequisites,
			req.HasFeatureFlagAsRule,
			req.SearchKeyword,
			req.Status,
			req.OrderBy,
			req.OrderDirection,
			req.HasExperiment.Value,
			req.EnvironmentId,
		)
	}
	if err != nil {
		return nil, err
	}
	featureCount, err := s.featureStorage.GetFeatureSummary(ctx, req.EnvironmentId)
	if err != nil {
		s.logger.Error(
			"Failed to count features by status",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, statusInternal.Err()
	}

	// TEMPORARY: Clean up any orphaned variation references before returning to UI
	// This prevents the UI from seeing corrupted data and sending it back in update requests
	for _, f := range features {
		domainFeature := &domain.Feature{Feature: f}
		cleanupResult := domainFeature.CleanupOrphanedVariationReferences()
		if cleanupResult.Changed {
			s.logger.Warn(
				"Cleaned up orphaned variation references in feature during list (temporary migration)",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.String("featureId", f.Id),
					zap.String("environmentId", req.EnvironmentId),
					zap.Int("orphanedTargets", cleanupResult.OrphanedTargets),
					zap.Int("orphanedRules", cleanupResult.OrphanedRules),
					zap.Int("orphanedDefault", cleanupResult.OrphanedDefault),
					zap.Bool("orphanedOffVar", cleanupResult.OrphanedOffVar),
					zap.Strings("orphanedVariationIDs", cleanupResult.OrphanedVariationIDs),
				)...,
			)
		}

		// TEMPORARY: Ensure all variations are present in rollout strategies during reads
		// This fixes existing data corruption from the historical AddVariation bug
		// TODO: Remove this after DB migration is complete
		migrationResult := domainFeature.EnsureVariationsInStrategies()
		if migrationResult.Changed {
			s.logger.Warn(
				"Added missing variations to rollout strategies in feature during list (temporary migration)",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.String("featureId", f.Id),
					zap.String("environmentId", req.EnvironmentId),
					zap.Int("addedToRules", migrationResult.AddedToRules),
					zap.Int("addedToDefault", migrationResult.AddedToDefault),
					zap.Strings("processedVariationIDs", migrationResult.AddedVariationIDs),
				)...,
			)
		}
	}

	return &featureproto.ListFeaturesResponse{
		Features:             features,
		Cursor:               cursor,
		TotalCount:           totalCount,
		FeatureCountByStatus: featureCount,
	}, nil
}

func (s *FeatureService) listFeatures(
	ctx context.Context,
	pageSize int64,
	cursor string,
	tags []string,
	maintainer string,
	enabled *wrappers.BoolValue,
	archived *wrappers.BoolValue,
	hasPrerequisites *wrappers.BoolValue,
	hasFeatureFlagAsRule *wrappers.BoolValue,
	searchKeyword string,
	status featureproto.FeatureLastUsedInfo_Status,
	orderBy featureproto.ListFeaturesRequest_OrderBy,
	orderDirection featureproto.ListFeaturesRequest_OrderDirection,
	environmentId string,
) ([]*featureproto.Feature, string, int64, error) {
	filters := []*mysql.FilterV2{
		{
			Column:   "feature.deleted",
			Operator: mysql.OperatorEqual,
			Value:    false,
		},
		{
			Column:   "feature.environment_id",
			Operator: mysql.OperatorEqual,
			Value:    environmentId,
		},
	}
	tagValues := make([]interface{}, 0, len(tags))
	for _, tag := range tags {
		tagValues = append(tagValues, tag)
	}
	var jsonFilters []*mysql.JSONFilter
	if len(tagValues) > 0 {
		jsonFilters = append(jsonFilters, &mysql.JSONFilter{
			Column: "feature.tags",
			Func:   mysql.JSONContainsString,
			Values: tagValues,
		})
	}
	if maintainer != "" {
		filters = append(filters, &mysql.FilterV2{
			Column:   "feature.maintainer",
			Operator: mysql.OperatorEqual,
			Value:    maintainer,
		})
	}
	if enabled != nil {
		filters = append(filters, &mysql.FilterV2{
			Column:   "feature.enabled",
			Operator: mysql.OperatorEqual,
			Value:    enabled.Value,
		})
	}
	if archived != nil {
		filters = append(filters, &mysql.FilterV2{
			Column:   "feature.archived",
			Operator: mysql.OperatorEqual,
			Value:    archived.Value,
		})
	}
	if hasPrerequisites != nil {
		if hasPrerequisites.Value {
			jsonFilters = append(jsonFilters, &mysql.JSONFilter{
				Column: "feature.prerequisites",
				Func:   mysql.JSONLengthGreaterThan,
				Values: []interface{}{"0"},
			})
		} else {
			jsonFilters = append(jsonFilters, &mysql.JSONFilter{
				Column: "feature.prerequisites",
				Func:   mysql.JSONLengthSmallerThan,
				Values: []interface{}{"1"},
			})
		}
	}
	if hasFeatureFlagAsRule != nil {
		// 11 is feature flag rule operator
		if hasFeatureFlagAsRule.Value {
			filters = append(filters, &mysql.FilterV2{
				Column:   "JSON_CONTAINS(JSON_EXTRACT(rules, '$[*].clauses[*].operator'), '11')",
				Operator: mysql.OperatorEqual,
				Value:    true,
			})
		} else {
			filters = append(filters, &mysql.FilterV2{
				Column:   "JSON_CONTAINS(JSON_EXTRACT(rules, '$[*].clauses[*].operator'), '11')",
				Operator: mysql.OperatorEqual,
				Value:    false,
			})
		}
	}
	var searchQuery *mysql.SearchQuery
	if searchKeyword != "" {
		searchQuery = &mysql.SearchQuery{
			Columns: []string{"feature.id", "feature.name", "feature.description"},
			Keyword: searchKeyword,
		}
	}
	var nullFilters []*mysql.NullFilter
	switch status {
	case featureproto.FeatureLastUsedInfo_UNKNOWN:
	case featureproto.FeatureLastUsedInfo_NEW:
		nullFilters = append(nullFilters, &mysql.NullFilter{
			Column: "feature_last_used_info.id",
			IsNull: true,
		})
	case featureproto.FeatureLastUsedInfo_ACTIVE:
		filters = append(filters, &mysql.FilterV2{
			Column:   "feature_last_used_info.last_used_at",
			Operator: mysql.OperatorGreaterThanOrEqual,
			Value:    time.Now().Add(-activeDays).Unix(),
		})
	case featureproto.FeatureLastUsedInfo_NO_ACTIVITY:
		filters = append(filters, &mysql.FilterV2{
			Column:   "feature_last_used_info.last_used_at",
			Operator: mysql.OperatorLessThan,
			Value:    time.Now().Add(-activeDays).Unix(),
		})
	}
	orders, err := s.newListFeaturesOrdersMySQL(orderBy, orderDirection)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		return nil, "", 0, err
	}
	limit := int(pageSize)
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, "", 0, statusInvalidCursor.Err()
	}
	options := &mysql.ListOptions{
		Filters:     filters,
		Orders:      orders,
		JSONFilters: jsonFilters,
		NullFilters: nullFilters,
		InFilters:   nil,
		SearchQuery: searchQuery,
		Limit:       limit,
		Offset:      offset,
	}
	features, nextCursor, totalCount, err := s.featureStorage.ListFeatures(ctx, options)
	if err != nil {
		s.logger.Error(
			"Failed to list features",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		return nil, "", 0, err
	}
	return features, strconv.Itoa(nextCursor), totalCount, nil
}

func (s *FeatureService) listFeaturesFilteredByExperiment(
	ctx context.Context,
	pageSize int64,
	cursor string,
	tags []string,
	maintainer string,
	enabled *wrappers.BoolValue,
	archived *wrappers.BoolValue,
	hasPrerequisites *wrappers.BoolValue,
	hasFeatureFlagAsRule *wrappers.BoolValue,
	searchKeyword string,
	status featureproto.FeatureLastUsedInfo_Status,
	orderBy featureproto.ListFeaturesRequest_OrderBy,
	orderDirection featureproto.ListFeaturesRequest_OrderDirection,
	hasExperiment bool,
	environmentId string,
) ([]*featureproto.Feature, string, int64, error) {
	filters := []*mysql.FilterV2{
		{
			Column:   "feature.deleted",
			Operator: mysql.OperatorEqual,
			Value:    false,
		},
		{
			Column:   "feature.environment_id",
			Operator: mysql.OperatorEqual,
			Value:    environmentId,
		},
	}
	nullFilters := []*mysql.NullFilter{
		{
			Column: "experiment.id",
			IsNull: !hasExperiment,
		},
	}
	tagValues := make([]interface{}, 0, len(tags))
	for _, tag := range tags {
		tagValues = append(tagValues, tag)
	}
	var jsonFilters []*mysql.JSONFilter
	if len(tagValues) > 0 {
		jsonFilters = append(jsonFilters, &mysql.JSONFilter{
			Column: "feature.tags",
			Func:   mysql.JSONContainsString,
			Values: tagValues,
		})
	}
	if maintainer != "" {
		filters = append(filters, &mysql.FilterV2{
			Column:   "feature.maintainer",
			Operator: mysql.OperatorEqual,
			Value:    maintainer,
		})
	}
	if enabled != nil {
		filters = append(filters, &mysql.FilterV2{
			Column:   "feature.enabled",
			Operator: mysql.OperatorEqual,
			Value:    enabled.Value,
		})
	}
	if archived != nil {
		filters = append(filters, &mysql.FilterV2{
			Column:   "feature.archived",
			Operator: mysql.OperatorEqual,
			Value:    archived.Value,
		})
	}
	if hasPrerequisites != nil {
		if hasPrerequisites.Value {
			jsonFilters = append(jsonFilters, &mysql.JSONFilter{
				Column: "feature.prerequisites",
				Func:   mysql.JSONLengthGreaterThan,
				Values: []interface{}{"0"},
			})
		} else {
			jsonFilters = append(jsonFilters, &mysql.JSONFilter{
				Column: "feature.prerequisites",
				Func:   mysql.JSONLengthSmallerThan,
				Values: []interface{}{"1"},
			})
		}
	}
	if hasFeatureFlagAsRule != nil {
		// 11 is feature flag rule operator
		if hasFeatureFlagAsRule.Value {
			filters = append(filters, &mysql.FilterV2{
				Column:   "JSON_CONTAINS(JSON_EXTRACT(rules, '$[*].clauses[*].operator'), '11')",
				Operator: mysql.OperatorEqual,
				Value:    true,
			})
		} else {
			filters = append(filters, &mysql.FilterV2{
				Column:   "JSON_CONTAINS(JSON_EXTRACT(rules, '$[*].clauses[*].operator'), '11')",
				Operator: mysql.OperatorEqual,
				Value:    false,
			})
		}
	}
	var searchQuery *mysql.SearchQuery
	if searchKeyword != "" {
		searchQuery = &mysql.SearchQuery{
			Columns: []string{"feature.id", "feature.name", "feature.description"},
			Keyword: searchKeyword,
		}
	}
	switch status {
	case featureproto.FeatureLastUsedInfo_UNKNOWN:
	case featureproto.FeatureLastUsedInfo_NEW:
		nullFilters = append(nullFilters, &mysql.NullFilter{
			Column: "feature_last_used_info.id",
			IsNull: true,
		})
	case featureproto.FeatureLastUsedInfo_ACTIVE:
		filters = append(filters, &mysql.FilterV2{
			Column:   "feature_last_used_info.last_used_at",
			Operator: mysql.OperatorGreaterThanOrEqual,
			Value:    time.Now().Add(-activeDays).Unix(),
		})
	case featureproto.FeatureLastUsedInfo_NO_ACTIVITY:
		filters = append(filters, &mysql.FilterV2{
			Column:   "feature_last_used_info.last_used_at",
			Operator: mysql.OperatorLessThan,
			Value:    time.Now().Add(-activeDays).Unix(),
		})
	}
	orders, err := s.newListFeaturesOrdersMySQL(orderBy, orderDirection)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		return nil, "", 0, err
	}
	limit := int(pageSize)
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, "", 0, statusInvalidCursor.Err()
	}
	featureStorage := v2fs.NewFeatureStorage(s.mysqlClient)
	options := &mysql.ListOptions{
		Filters:     filters,
		Orders:      orders,
		JSONFilters: jsonFilters,
		NullFilters: nullFilters,
		InFilters:   nil,
		SearchQuery: searchQuery,
		Limit:       limit,
		Offset:      offset,
	}
	features, nextCursor, totalCount, err := featureStorage.ListFeaturesFilteredByExperiment(ctx, options)
	if err != nil {
		s.logger.Error(
			"Failed to list features filtered by experiment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		return nil, "", 0, err
	}
	return features, strconv.Itoa(nextCursor), totalCount, nil
}

func (s *FeatureService) newListFeaturesOrdersMySQL(
	orderBy featureproto.ListFeaturesRequest_OrderBy,
	orderDirection featureproto.ListFeaturesRequest_OrderDirection,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case featureproto.ListFeaturesRequest_DEFAULT,
		featureproto.ListFeaturesRequest_NAME:
		column = "feature.name"
	case featureproto.ListFeaturesRequest_CREATED_AT:
		column = "feature.created_at"
	case featureproto.ListFeaturesRequest_UPDATED_AT:
		column = "feature.updated_at"
	case featureproto.ListFeaturesRequest_TAGS:
		column = "feature.tags"
	case featureproto.ListFeaturesRequest_ENABLED:
		column = "feature.enabled"
	case featureproto.ListFeaturesRequest_AUTO_OPS:
		column = "(progressive_rollout_count + schedule_count + kill_switch_count)"
	default:
		return nil, statusInvalidOrderBy.Err()
	}
	direction := mysql.OrderDirectionAsc
	if orderDirection == featureproto.ListFeaturesRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *FeatureService) ListEnabledFeatures(
	ctx context.Context,
	req *featureproto.ListEnabledFeaturesRequest,
) (*featureproto.ListEnabledFeaturesResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	filters := []*mysql.FilterV2{
		{
			Column:   "archived",
			Operator: mysql.OperatorEqual,
			Value:    false,
		},
		{
			Column:   "enabled",
			Operator: mysql.OperatorEqual,
			Value:    true,
		},
		{
			Column:   "deleted",
			Operator: mysql.OperatorEqual,
			Value:    false,
		},
		{
			Column:   "feature.environment_id",
			Operator: mysql.OperatorEqual,
			Value:    req.EnvironmentId,
		},
	}
	tagValues := make([]interface{}, 0, len(req.Tags))
	for _, tag := range req.Tags {
		tagValues = append(tagValues, tag)
	}
	var jsonFilters []*mysql.JSONFilter
	if len(tagValues) > 0 {
		jsonFilters = append(
			jsonFilters,
			&mysql.JSONFilter{
				Column: "tags",
				Func:   mysql.JSONContainsString,
				Values: tagValues,
			})
	}
	limit := int(req.PageSize)
	cursor := req.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, statusInvalidCursor.Err()
	}
	featureStorage := v2fs.NewFeatureStorage(s.mysqlClient)
	options := &mysql.ListOptions{
		Filters:     filters,
		JSONFilters: jsonFilters,
		Orders:      nil,
		NullFilters: nil,
		InFilters:   nil,
		SearchQuery: nil,
		Limit:       limit,
		Offset:      offset,
	}
	features, nextCursor, _, err := featureStorage.ListFeatures(ctx, options)
	if err != nil {
		s.logger.Error(
			"Failed to list enabled features",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	if err = s.setLastUsedInfosToFeatureByChunk(ctx, features, req.EnvironmentId); err != nil {
		return nil, err
	}
	return &featureproto.ListEnabledFeaturesResponse{
		Features: features,
		Cursor:   strconv.Itoa(nextCursor),
	}, nil
}

func (s *FeatureService) CreateFeature(
	ctx context.Context,
	req *featureproto.CreateFeatureRequest,
) (*featureproto.CreateFeatureResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	err = validateCreateFeatureRequest(req)
	if err != nil {
		return nil, err
	}
	feature, err := domain.NewFeature(
		req.Id,
		req.Name,
		req.Description,
		req.VariationType,
		req.Variations,
		req.Tags,
		int(req.DefaultOnVariationIndex.Value),
		int(req.DefaultOffVariationIndex.Value),
		editor.Email,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create feature",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_FEATURE,
			feature.Id,
			eventproto.Event_FEATURE_CREATED,
			&eventproto.FeatureCreatedEvent{
				Id:                       feature.Id,
				Name:                     feature.Name,
				Description:              feature.Description,
				User:                     "default",
				Variations:               feature.Variations,
				DefaultOnVariationIndex:  req.DefaultOnVariationIndex,
				DefaultOffVariationIndex: req.DefaultOffVariationIndex,
				VariationType:            req.VariationType,
				Tags:                     feature.Tags,
				Prerequisites:            feature.Prerequisites,
				Targets:                  feature.Targets,
				Rules:                    feature.Rules,
			},
			req.EnvironmentId,
			feature,
			nil,
		)
		if err != nil {
			return err
		}
		if err := s.upsertTags(ctxWithTx, req.Tags, req.EnvironmentId); err != nil {
			return err
		}
		return s.featureStorage.CreateFeature(ctxWithTx, feature, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrFeatureAlreadyExists) {
			return nil, statusAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create feature",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.domainPublisher.Publish(ctx, event)
	if err != nil {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Any("errors", err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	s.updateFeatureFlagCache(ctx)
	return &featureproto.CreateFeatureResponse{Feature: feature.Feature}, nil
}

func (s *FeatureService) UpdateFeature(
	ctx context.Context,
	req *featureproto.UpdateFeatureRequest,
) (*featureproto.UpdateFeatureResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		return nil, statusMissingID.Err()
	}
	if err := s.validateFeatureStatus(ctx, req.Id, req.EnvironmentId); err != nil {
		return nil, err
	}
	if err := s.validateEnvironmentSettings(ctx, req.EnvironmentId, req.Comment); err != nil {
		return nil, err
	}
	var event *eventproto.Event
	var updatedpb *featureproto.Feature
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		filters := []*mysql.FilterV2{
			{
				Column:   "feature.deleted",
				Operator: mysql.OperatorEqual,
				Value:    false,
			},
			{
				Column:   "feature.environment_id",
				Operator: mysql.OperatorEqual,
				Value:    req.EnvironmentId,
			},
		}
		options := &mysql.ListOptions{
			Filters:     filters,
			JSONFilters: nil,
			Orders:      nil,
			NullFilters: nil,
			InFilters:   nil,
			SearchQuery: nil,
			Limit:       mysql.QueryNoLimit,
			Offset:      mysql.QueryNoOffset,
		}
		features, _, _, err := s.featureStorage.ListFeatures(ctxWithTx, options)
		if err != nil {
			s.logger.Error(
				"Failed to list features",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return err
		}
		var feature *domain.Feature
		for _, f := range features {
			if f.Id == req.Id {
				feature = &domain.Feature{Feature: f}
				break
			}
		}
		if feature == nil {
			err := statusFeatureNotFound.Err()
			s.logger.Error(
				"Failed to find feature",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", req.Id),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return err
		}

		// Clean up any orphaned variation references BEFORE validation
		// This fixes data corruption from the historical variation deletion bug
		cleanupResult := feature.CleanupOrphanedVariationReferences()
		if cleanupResult.Changed {
			s.logger.Warn(
				"Cleaned up orphaned variation references in feature during update",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.String("featureId", req.Id),
					zap.String("environmentId", req.EnvironmentId),
					zap.Int("orphanedTargets", cleanupResult.OrphanedTargets),
					zap.Int("orphanedRules", cleanupResult.OrphanedRules),
					zap.Int("orphanedDefault", cleanupResult.OrphanedDefault),
					zap.Bool("orphanedOffVar", cleanupResult.OrphanedOffVar),
					zap.Strings("orphanedVariationIDs", cleanupResult.OrphanedVariationIDs),
				)...,
			)
		}

		// Ensure all variations are present in rollout strategies BEFORE validation
		// This fixes data corruption from the historical AddVariation bug
		migrationResult := feature.EnsureVariationsInStrategies()
		if migrationResult.Changed {
			s.logger.Warn(
				"Added missing variations to rollout strategies in feature during update",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.String("featureId", req.Id),
					zap.String("environmentId", req.EnvironmentId),
					zap.Int("addedToRules", migrationResult.AddedToRules),
					zap.Int("addedToDefault", migrationResult.AddedToDefault),
					zap.Strings("processedVariationIDs", migrationResult.AddedVariationIDs),
				)...,
			)
		}

		// Check if this is an archive request and if other features depend on this one.
		// This check mirrors the validation in ArchiveFeature to ensure consistent behavior.
		if req.Archived != nil && req.Archived.GetValue() && !feature.Archived {
			if domain.HasFeaturesDependsOnTargets([]*featureproto.Feature{feature.Feature}, features) {
				return statusInvalidArchive.Err()
			}
		}

		updated, err := feature.Update(
			req.Name,
			req.Description,
			req.Tags,
			req.Enabled,
			req.Archived,
			req.DefaultStrategy,
			req.OffVariation,
			req.ResetSamplingSeed,
			req.PrerequisiteChanges,
			req.TargetChanges,
			req.RuleChanges,
			req.VariationChanges,
			req.TagChanges,
			req.Maintainer,
		)
		if err != nil {
			return err
		}
		if err := s.upsertTags(ctxWithTx, updated.Tags, req.EnvironmentId); err != nil {
			return err
		}
		// To check if the flag to be updated is a dependency of other flags, we must validate it before updating.
		// Exclude all the archived and deleted flags from the list.
		tgts := []*featureproto.Feature{}
		for _, f := range features {
			if f.Id == updated.Id {
				f = updated.Feature
			}
			if f.Archived || f.Deleted {
				continue
			}
			tgts = append(tgts, f)
		}
		if err := domain.ValidateFeatureDependencies(tgts); err != nil {
			return err
		}
		// Validate that variations being deleted are not used in other features
		if err := validateVariationDeletion(req.VariationChanges, features, req.Id); err != nil {
			return err
		}
		updatedpb = updated.Feature
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_FEATURE,
			feature.Id,
			eventproto.Event_FEATURE_UPDATED,
			&eventproto.FeatureUpdatedEvent{
				Id: req.Id,
			},
			req.EnvironmentId,
			updated.Feature,
			feature.Feature,
			// check require comment.
			domainevent.WithComment(req.Comment),
			domainevent.WithNewVersion(updated.Version),
		)
		if err != nil {
			return err
		}
		err = s.featureStorage.UpdateFeature(ctxWithTx, updated, req.EnvironmentId)
		if err != nil {
			s.logger.Error(
				"Failed to update feature",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if errs := s.publishDomainEvents(ctx, []*eventproto.Event{event}); len(errs) > 0 {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Any("errors", errs),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, statusInternal.Err()
	}
	s.updateFeatureFlagCache(ctx)
	return &featureproto.UpdateFeatureResponse{
		Feature: updatedpb,
	}, nil
}

func (s *FeatureService) existsRunningExperiment(
	ctx context.Context,
	featureID, environmentId string,
) (bool, error) {
	experiments, err := s.listExperiments(ctx, environmentId, featureID)
	if err != nil {
		s.logger.Error(
			"Failed to list experiments",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		return false, err
	}
	return containsRunningExperiment(experiments), nil
}

func containsRunningExperiment(experiments []*experimentproto.Experiment) bool {
	now := time.Now()
	for _, e := range experiments {
		de := &experimentdomain.Experiment{Experiment: e}
		if de.IsNotFinished(now) {
			return true
		}
	}
	return false
}

func (s *FeatureService) DeleteFeature(
	ctx context.Context,
	req *featureproto.DeleteFeatureRequest,
) (*featureproto.DeleteFeatureResponse, error) {
	if err := validateDeleteFeatureRequest(req); err != nil {
		return nil, err
	}
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := s.validateEnvironmentSettings(ctx, req.EnvironmentId, req.Comment); err != nil {
		return nil, err
	}
	var eventPb *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		feature, err := s.featureStorage.GetFeature(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			s.logger.Error(
				"Failed to get feature",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", req.Id),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return err
		}

		// Clean up any orphaned variation references BEFORE validation
		// This fixes data corruption from the historical variation deletion bug
		cleanupResult := feature.CleanupOrphanedVariationReferences()
		if cleanupResult.Changed {
			s.logger.Warn(
				"Cleaned up orphaned variation references in feature during details update",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.String("featureId", req.Id),
					zap.String("environmentId", req.EnvironmentId),
					zap.Int("orphanedTargets", cleanupResult.OrphanedTargets),
					zap.Int("orphanedRules", cleanupResult.OrphanedRules),
					zap.Int("orphanedDefault", cleanupResult.OrphanedDefault),
					zap.Bool("orphanedOffVar", cleanupResult.OrphanedOffVar),
					zap.Strings("orphanedVariationIDs", cleanupResult.OrphanedVariationIDs),
				)...,
			)
		}

		// Ensure all variations are present in rollout strategies BEFORE validation
		// This fixes data corruption from the historical AddVariation bug
		migrationResult := feature.EnsureVariationsInStrategies()
		if migrationResult.Changed {
			s.logger.Warn(
				"Added missing variations to rollout strategies in feature during details update",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.String("featureId", req.Id),
					zap.String("environmentId", req.EnvironmentId),
					zap.Int("addedToRules", migrationResult.AddedToRules),
					zap.Int("addedToDefault", migrationResult.AddedToDefault),
					zap.Strings("processedVariationIDs", migrationResult.AddedVariationIDs),
				)...,
			)
		}

		err = feature.Delete()
		if err != nil {
			s.logger.Error(
				"Failed to delete feature",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return err
		}
		eventPb, err = domainevent.NewEvent(
			editor,
			eventproto.Event_FEATURE,
			feature.Id,
			eventproto.Event_FEATURE_DELETED,
			&eventproto.FeatureDeletedEvent{
				Id: req.Id,
			},
			req.EnvironmentId,
			nil,
			feature.Feature,
			// check require comment.
			domainevent.WithComment(req.Comment),
			domainevent.WithNewVersion(feature.Version),
		)
		if err != nil {
			return err
		}

		if err := s.featureStorage.UpdateFeature(contextWithTx, feature, req.EnvironmentId); err != nil {
			s.logger.Error(
				"Failed to update feature",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, s.convUpdateFeatureError(err)
	}
	err = s.domainPublisher.Publish(ctx, eventPb)
	if err != nil {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Any("errors", err),
				zap.String("environmentId", req.EnvironmentId),
				zap.Any("event", eventPb),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	s.updateFeatureFlagCache(ctx)
	return &featureproto.DeleteFeatureResponse{}, nil
}

func (s *FeatureService) convUpdateFeatureError(err error) error {
	switch err {
	case v2fs.ErrFeatureNotFound,
		v2fs.ErrFeatureUnexpectedAffectedRows,
		storage.ErrKeyNotFound:
		return statusFeatureNotFound.Err()
	default:
		return api.NewGRPCStatus(err).Err()
	}
}

func (s *FeatureService) publishDomainEvents(ctx context.Context, events []*eventproto.Event) map[string]error {
	messages := make([]publisher.Message, 0, len(events))
	for _, event := range events {
		messages = append(messages, event)
	}
	return s.domainPublisher.PublishMulti(ctx, messages)
}

func findFeature(fs []*featureproto.Feature, id string) (*featureproto.Feature, error) {
	for _, f := range fs {
		if f.Id == id {
			return f, nil
		}
	}
	return nil, statusInternal.Err()
}

func (s *FeatureService) evaluateFeatures(
	ctx context.Context,
	features []*featureproto.Feature,
	user *userproto.User,
	EnvironmentId string,
	tag string,
) (*featureproto.UserEvaluations, error) {
	evaluator := evaluation.NewEvaluator()
	mapIDs := make(map[string]struct{})
	for _, f := range features {
		for _, id := range evaluator.ListSegmentIDs(f) {
			mapIDs[id] = struct{}{}
		}
	}
	mapSegmentUsers, err := s.listSegmentUsers(ctx, mapIDs, EnvironmentId)
	if err != nil {
		s.logger.Error(
			"Failed to list segments",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", EnvironmentId),
				zap.String("userId", user.Id),
				zap.String("tag", tag),
			)...,
		)
		return nil, err
	}
	userEvaluations, err := evaluator.EvaluateFeatures(features, user, mapSegmentUsers, tag)
	if err != nil {
		s.logger.Error(
			"Failed to evaluate",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", EnvironmentId),
				zap.String("userId", user.Id),
				zap.String("tag", tag),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return userEvaluations, nil
}

func (s *FeatureService) getFeatures(
	ctx context.Context,
	EnvironmentId string,
) ([]*featureproto.Feature, error) {
	features, err := s.featuresCache.Get(EnvironmentId)
	if err == nil {
		return features.Features, nil
	}
	s.logger.Warn(
		"No cached data for Features",
		log.FieldsFromIncomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentId", EnvironmentId),
		)...,
	)
	fs, _, _, err := s.listFeatures(
		ctx,
		mysql.QueryNoLimit,
		"",
		nil,
		"",
		nil,
		nil,
		nil,
		nil,
		"",
		featureproto.FeatureLastUsedInfo_UNKNOWN,
		featureproto.ListFeaturesRequest_DEFAULT,
		featureproto.ListFeaturesRequest_ASC,
		EnvironmentId,
	)
	if err != nil {
		s.logger.Error(
			"Failed to retrieve features from storage",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", EnvironmentId),
			)...,
		)
		return nil, err
	}
	return fs, nil
}

func (s *FeatureService) listSegmentUsers(
	ctx context.Context,
	mapSegmentIDs map[string]struct{},
	EnvironmentId string,
) (map[string][]*featureproto.SegmentUser, error) {
	if len(mapSegmentIDs) == 0 {
		return nil, nil
	}
	users := make(map[string][]*featureproto.SegmentUser)
	for segmentID := range mapSegmentIDs {
		s, err, _ := s.flightgroup.Do(
			s.segmentFlightID(EnvironmentId, segmentID),
			func() (interface{}, error) {
				return s.getSegmentUsers(ctx, segmentID, EnvironmentId)
			},
		)
		if err != nil {
			return nil, err
		}
		listUsers := s.([]*featureproto.SegmentUser)
		users[segmentID] = listUsers
	}
	return users, nil
}

func (s *FeatureService) segmentFlightID(EnvironmentId, segmentID string) string {
	return fmt.Sprintf("%s:%s", EnvironmentId, segmentID)
}

func (s *FeatureService) getSegmentUsers(
	ctx context.Context,
	segmentID, EnvironmentId string,
) ([]*featureproto.SegmentUser, error) {
	segmentUsers, err := s.segmentUsersCache.Get(segmentID, EnvironmentId)
	if err == nil {
		return segmentUsers.Users, nil
	}
	s.logger.Warn(
		"No cached data for SegmentUsers",
		log.FieldsFromIncomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentId", EnvironmentId),
			zap.String("segmentId", segmentID),
		)...,
	)
	req := &featureproto.ListSegmentUsersRequest{
		SegmentId:     segmentID,
		EnvironmentId: EnvironmentId,
	}
	res, storageErr := s.ListSegmentUsers(ctx, req)
	if storageErr != nil {
		s.logger.Error(
			"Failed to retrieve segment users from storage",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(storageErr),
				zap.String("environmentId", EnvironmentId),
				zap.String("segmentId", segmentID),
			)...,
		)
		return nil, err
	}
	return res.Users, nil
}

func (s *FeatureService) setLastUsedInfosToFeatureByChunk(
	ctx context.Context,
	features []*featureproto.Feature,
	EnvironmentId string,
) error {
	for i := 0; i < len(features); i += getMultiChunkSize {
		end := i + getMultiChunkSize
		if end > len(features) {
			end = len(features)
		}
		if err := s.setLastUsedInfosToFeature(ctx, features[i:end], EnvironmentId); err != nil {
			return err
		}
	}
	return nil
}

func (s *FeatureService) setLastUsedInfosToFeature(
	ctx context.Context,
	features []*featureproto.Feature,
	EnvironmentId string,
) error {
	ids := make([]string, 0, len(features))
	for _, f := range features {
		ids = append(ids, domain.FeatureLastUsedInfoID(f.Id, f.Version))
	}
	fluiList, err := s.fluiStorage.GetFeatureLastUsedInfos(ctx, ids, EnvironmentId)
	if err != nil {
		s.logger.Error(
			"Failed to get feature last used infos",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", EnvironmentId),
			)...,
		)
		return api.NewGRPCStatus(err).Err()
	}
	for _, f := range fluiList {
		for _, feature := range features {
			if feature.Id == f.FeatureId {
				feature.LastUsedInfo = f.FeatureLastUsedInfo
				break
			}
		}
	}
	return nil
}

func (s *FeatureService) EvaluateFeatures(
	ctx context.Context,
	req *featureproto.EvaluateFeaturesRequest,
) (*featureproto.EvaluateFeaturesResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateEvaluateFeatures(req); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	fs, err, _ := s.flightgroup.Do(
		req.EnvironmentId,
		func() (interface{}, error) {
			return s.getFeatures(ctx, req.EnvironmentId)
		},
	)
	if err != nil {
		s.logger.Error(
			"Failed to list features",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	// If the feature ID is set in the request, it will evaluate a single feature.
	features, err := s.getTargetFeatures(fs.([]*featureproto.Feature), req.FeatureId)
	if err != nil {
		s.logger.Error(
			"Failed to get target features",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	userEvaluations, err := s.evaluateFeatures(ctx, features, req.User, req.EnvironmentId, req.Tag)
	if err != nil {
		s.logger.Error(
			"Failed to evaluate features",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	// If the feature ID is set, it will return a single evaluation
	if req.FeatureId != "" {
		eval, err := s.findEvaluation(userEvaluations.Evaluations, req.FeatureId)
		if err != nil {
			s.logger.Error(
				"Failed to find evaluation",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return nil, api.NewGRPCStatus(err).Err()
		}
		return &featureproto.EvaluateFeaturesResponse{
			UserEvaluations: &featureproto.UserEvaluations{
				Id:          userEvaluations.Id,
				Evaluations: []*featureproto.Evaluation{eval},
				CreatedAt:   userEvaluations.CreatedAt,
			}}, nil
	}
	return &featureproto.EvaluateFeaturesResponse{UserEvaluations: userEvaluations}, nil
}

func (s *FeatureService) DebugEvaluateFeatures(
	ctx context.Context,
	req *featureproto.DebugEvaluateFeaturesRequest,
) (*featureproto.DebugEvaluateFeaturesResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	err = validateDebugEvaluateFeatures(req)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	fs, err, _ := s.flightgroup.Do(
		req.EnvironmentId,
		func() (interface{}, error) {
			return s.getFeatures(ctx, req.EnvironmentId)
		},
	)
	if err != nil {
		s.logger.Error(
			"Failed to list features",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	features := fs.([]*featureproto.Feature)
	var evaluations = make([]*featureproto.Evaluation, 0)
	var archivedFS = make([]string, 0)
	// If the feature ID is set in the request, it will evaluate a single feature.
	if len(req.FeatureIds) == 1 {
		features, err = s.getTargetFeatures(features, req.FeatureIds[0])
		if err != nil {
			s.logger.Error(
				"Failed to get target features",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return nil, api.NewGRPCStatus(err).Err()
		}
	}

	for i := range req.Users {
		userEvaluations, err := s.evaluateFeatures(
			ctx, features, req.Users[i], req.EnvironmentId, "",
		)
		if err != nil {
			s.logger.Error(
				"Failed to evaluate features",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return nil, api.NewGRPCStatus(err).Err()
		}

		evaluations = append(evaluations, userEvaluations.Evaluations...)
		archivedFS = append(archivedFS, userEvaluations.ArchivedFeatureIds...)
	}
	evaluationResults := make([]*featureproto.Evaluation, 0)
	for _, eval := range evaluations {
		if slices.Contains(req.FeatureIds, eval.FeatureId) {
			evaluationResults = append(evaluationResults, eval)
		}
	}

	return &featureproto.DebugEvaluateFeaturesResponse{
		Evaluations:        evaluationResults,
		ArchivedFeatureIds: archivedFS,
	}, nil
}

func (s *FeatureService) getTargetFeatures(
	fs []*featureproto.Feature,
	id string,
) ([]*featureproto.Feature, error) {
	if id == "" {
		return fs, nil
	}
	feature, err := findFeature(fs, id)
	if err != nil {
		return nil, err
	}
	// Check if the flag depends on other flags.
	// Thus, we evaluate all features here to avoid complex logic.
	df := &domain.Feature{Feature: feature}
	if len(df.FeatureIDsDependsOn()) > 0 {
		return fs, nil
	}
	return []*featureproto.Feature{feature}, nil
}

func (*FeatureService) findEvaluation(
	evals []*featureproto.Evaluation,
	id string,
) (*featureproto.Evaluation, error) {
	for _, e := range evals {
		if e.FeatureId == id {
			return e, nil
		}
	}
	return nil, errEvaluationNotFound
}

func (s *FeatureService) listExperiments(
	ctx context.Context,
	EnvironmentId, featureID string,
) ([]*experimentproto.Experiment, error) {
	experiments := []*experimentproto.Experiment{}
	cursor := ""
	for {
		resp, err := s.experimentClient.ListExperiments(ctx, &experimentproto.ListExperimentsRequest{
			FeatureId:     featureID,
			PageSize:      listRequestSize,
			Cursor:        cursor,
			EnvironmentId: EnvironmentId,
		})
		if err != nil {
			return nil, err
		}
		experiments = append(experiments, resp.Experiments...)
		featureSize := len(resp.Experiments)
		if featureSize == 0 || featureSize < listRequestSize {
			return experiments, nil
		}
		cursor = resp.Cursor
	}
}

func (s *FeatureService) CloneFeature(
	ctx context.Context,
	req *featureproto.CloneFeatureRequest,
) (*featureproto.CloneFeatureResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.TargetEnvironmentId,
	)
	if err != nil {
		return nil, err
	}
	err = validateCloneFeatureRequest(req)
	if err != nil {
		return nil, err
	}
	featureStorage := v2fs.NewFeatureStorage(s.mysqlClient)
	f, err := featureStorage.GetFeature(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2fs.ErrFeatureNotFound) {
			return nil, statusFeatureNotFound.Err()
		}
		s.logger.Error(
			"Failed to get feature",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("targetEnvironmentId", req.TargetEnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	domainFeature := &domain.Feature{
		Feature: f.Feature,
	}
	feature, err := domainFeature.Clone(editor.Email)
	if err != nil {
		s.logger.Error(
			"Failed to clone domain feature",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("targetEnvironmentId", req.TargetEnvironmentId),
			)...,
		)
		return nil, err
	}
	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_FEATURE,
			feature.Id,
			eventproto.Event_FEATURE_CLONED,
			&eventproto.FeatureClonedEvent{
				Id:                feature.Id,
				Name:              feature.Name,
				Description:       feature.Description,
				Variations:        feature.Variations,
				Targets:           feature.Targets,
				Rules:             feature.Rules,
				DefaultStrategy:   feature.DefaultStrategy,
				OffVariation:      feature.OffVariation,
				Tags:              feature.Tags,
				Maintainer:        feature.Maintainer,
				VariationType:     feature.VariationType,
				Prerequisites:     feature.Prerequisites,
				SourceEnvironment: req.EnvironmentId,
				TargetEnvironment: req.TargetEnvironmentId,
			},
			req.TargetEnvironmentId,
			feature,
			feature,
		)

		if err := featureStorage.CreateFeature(ctxWithTx, feature, req.TargetEnvironmentId); err != nil {
			s.logger.Error(
				"Failed to store feature",
				log.FieldsFromIncomingContext(ctxWithTx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
					zap.String("targetEnvironmentId", req.TargetEnvironmentId),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrFeatureAlreadyExists) {
			return nil, statusAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to clone feature",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	if err = s.domainPublisher.Publish(ctx, event); err != nil {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Any("errors", err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("targetEnvironmentId", req.TargetEnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	s.updateFeatureFlagCache(ctx)
	return &featureproto.CloneFeatureResponse{}, nil
}

// Even if the update request fails, the cronjob will keep trying
// to update the cache every minute, so we don't need to retry.
func (s *FeatureService) updateFeatureFlagCache(ctx context.Context) {
	req := &btproto.BatchJobRequest{
		Job: btproto.BatchJob_FeatureFlagCacher,
	}
	_, err := s.batchClient.ExecuteBatchJob(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update feature flag cache", zap.Error(err))
	}
}
