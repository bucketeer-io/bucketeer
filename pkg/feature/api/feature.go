// Copyright 2023 The Bucketeer Authors.
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
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	autoopsdomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	experimentdomain "github.com/bucketeer-io/bucketeer/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/feature/command"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

const (
	getMultiChunkSize = 1000
	listRequestSize   = 500
)

var errEvaluationNotFound = status.Error(codes.NotFound, "feature: evaluation not found")

func (s *FeatureService) GetFeature(
	ctx context.Context,
	req *featureproto.GetFeatureRequest,
) (*featureproto.GetFeatureResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetFeatureRequest(req, localizer); err != nil {
		return nil, err
	}
	featureStorage := v2fs.NewFeatureStorage(s.mysqlClient)
	feature, err := featureStorage.GetFeature(ctx, req.Id, req.EnvironmentNamespace)
	if err != nil {
		if err == v2fs.ErrFeatureNotFound {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to get feature",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if err := s.setLastUsedInfosToFeatureByChunk(
		ctx,
		[]*featureproto.Feature{feature.Feature},
		req.EnvironmentNamespace,
		localizer,
	); err != nil {
		return nil, err
	}
	return &featureproto.GetFeatureResponse{Feature: feature.Feature}, nil
}

func (s *FeatureService) GetFeatures(
	ctx context.Context,
	req *featureproto.GetFeaturesRequest,
) (*featureproto.GetFeaturesResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetFeaturesRequest(req, localizer); err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace),
	}
	ids := make([]interface{}, 0, len(req.Ids))
	for _, id := range req.Ids {
		ids = append(ids, id)
	}
	if len(ids) > 0 {
		whereParts = append(whereParts, mysql.NewInFilter("id", ids))
	}
	featureStorage := v2fs.NewFeatureStorage(s.mysqlClient)
	features, _, _, err := featureStorage.ListFeatures(
		ctx,
		whereParts,
		nil,
		mysql.QueryNoLimit,
		mysql.QueryNoOffset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to get feature",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &featureproto.GetFeaturesResponse{Features: features}, nil
}

func (s *FeatureService) ListFeatures(
	ctx context.Context,
	req *featureproto.ListFeaturesRequest,
) (*featureproto.ListFeaturesResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
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
			req.SearchKeyword,
			req.OrderBy,
			req.OrderDirection,
			req.EnvironmentNamespace,
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
			req.SearchKeyword,
			req.OrderBy,
			req.OrderDirection,
			req.HasExperiment.Value,
			req.EnvironmentNamespace,
		)
	}
	if err != nil {
		return nil, err
	}
	return &featureproto.ListFeaturesResponse{
		Features:   features,
		Cursor:     cursor,
		TotalCount: totalCount,
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
	searchKeyword string,
	orderBy featureproto.ListFeaturesRequest_OrderBy,
	orderDirection featureproto.ListFeaturesRequest_OrderDirection,
	environmentNamespace string,
) ([]*featureproto.Feature, string, int64, error) {
	localizer := locale.NewLocalizer(ctx)
	whereParts := []mysql.WherePart{
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_namespace", "=", environmentNamespace),
	}
	tagValues := make([]interface{}, 0, len(tags))
	for _, tag := range tags {
		tagValues = append(tagValues, tag)
	}
	if len(tagValues) > 0 {
		whereParts = append(
			whereParts,
			mysql.NewJSONFilter("tags", mysql.JSONContainsString, tagValues),
		)
	}
	if maintainer != "" {
		whereParts = append(whereParts, mysql.NewFilter("maintainer", "=", maintainer))
	}
	if enabled != nil {
		whereParts = append(whereParts, mysql.NewFilter("enabled", "=", enabled.Value))
	}
	if archived != nil {
		whereParts = append(whereParts, mysql.NewFilter("archived", "=", archived.Value))
	}
	if hasPrerequisites != nil {
		if hasPrerequisites.Value {
			whereParts = append(
				whereParts,
				mysql.NewJSONFilter("prerequisites", mysql.JSONLengthGreaterThan, []interface{}{"0"}),
			)
		} else {
			whereParts = append(
				whereParts,
				mysql.NewJSONFilter("prerequisites", mysql.JSONLengthSmallerThan, []interface{}{"1"}),
			)
		}
	}
	if searchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"id", "name", "description"}, searchKeyword))
	}
	orders, err := s.newListFeaturesOrdersMySQL(orderBy, orderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
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
		dt, err := statusInvalidCursor.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor"),
		})
		if err != nil {
			return nil, "", 0, statusInternal.Err()
		}
		return nil, "", 0, dt.Err()
	}
	featureStorage := v2fs.NewFeatureStorage(s.mysqlClient)
	features, nextCursor, totalCount, err := featureStorage.ListFeatures(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list features",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
			)...,
		)
		return nil, "", 0, err
	}
	if err = s.setLastUsedInfosToFeatureByChunk(ctx, features, environmentNamespace, localizer); err != nil {
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
	searchKeyword string,
	orderBy featureproto.ListFeaturesRequest_OrderBy,
	orderDirection featureproto.ListFeaturesRequest_OrderDirection,
	hasExperiment bool,
	environmentNamespace string,
) ([]*featureproto.Feature, string, int64, error) {
	localizer := locale.NewLocalizer(ctx)
	whereParts := []mysql.WherePart{
		mysql.NewFilter("feature.deleted", "=", false),
		mysql.NewFilter("experiment.deleted", "=", false),
		mysql.NewFilter("feature.environment_namespace", "=", environmentNamespace),
		mysql.NewNullFilter("experiment.id", !hasExperiment),
	}
	tagValues := make([]interface{}, 0, len(tags))
	for _, tag := range tags {
		tagValues = append(tagValues, tag)
	}
	if len(tagValues) > 0 {
		whereParts = append(
			whereParts,
			mysql.NewJSONFilter("feature.tags", mysql.JSONContainsString, tagValues),
		)
	}
	if maintainer != "" {
		whereParts = append(whereParts, mysql.NewFilter("feature.maintainer", "=", maintainer))
	}
	if enabled != nil {
		whereParts = append(whereParts, mysql.NewFilter("feature.enabled", "=", enabled.Value))
	}
	if archived != nil {
		whereParts = append(whereParts, mysql.NewFilter("feature.archived", "=", archived.Value))
	}
	if hasPrerequisites != nil {
		if hasPrerequisites.Value {
			whereParts = append(
				whereParts,
				mysql.NewJSONFilter("prerequisites", mysql.JSONLengthGreaterThan, []interface{}{"0"}),
			)
		} else {
			whereParts = append(
				whereParts,
				mysql.NewJSONFilter("prerequisites", mysql.JSONLengthSmallerThan, []interface{}{"1"}),
			)
		}
	}
	if searchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"id", "name", "description"}, searchKeyword))
	}
	orders, err := s.newListFeaturesOrdersMySQL(orderBy, orderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
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
		dt, err := statusInvalidCursor.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor"),
		})
		if err != nil {
			return nil, "", 0, statusInternal.Err()
		}
		return nil, "", 0, dt.Err()
	}
	featureStorage := v2fs.NewFeatureStorage(s.mysqlClient)
	features, nextCursor, totalCount, err := featureStorage.ListFeaturesFilteredByExperiment(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list features filtered by experiment",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
			)...,
		)
		return nil, "", 0, err
	}
	if err = s.setLastUsedInfosToFeatureByChunk(ctx, features, environmentNamespace, localizer); err != nil {
		return nil, "", 0, err
	}
	return features, strconv.Itoa(nextCursor), totalCount, nil
}

func (s *FeatureService) newListFeaturesOrdersMySQL(
	orderBy featureproto.ListFeaturesRequest_OrderBy,
	orderDirection featureproto.ListFeaturesRequest_OrderDirection,
	localizer locale.Localizer,
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
	default:
		dt, err := statusInvalidOrderBy.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "order_by"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
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
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("archived", "=", false),
		mysql.NewFilter("enabled", "=", true),
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace),
	}
	tagValues := make([]interface{}, 0, len(req.Tags))
	for _, tag := range req.Tags {
		tagValues = append(tagValues, tag)
	}
	if len(tagValues) > 0 {
		whereParts = append(
			whereParts,
			mysql.NewJSONFilter("tags", mysql.JSONContainsString, tagValues),
		)
	}
	limit := int(req.PageSize)
	cursor := req.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		dt, err := statusInvalidCursor.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	featureStorage := v2fs.NewFeatureStorage(s.mysqlClient)
	features, nextCursor, _, err := featureStorage.ListFeatures(
		ctx,
		whereParts,
		nil,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list enabled features",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	if err = s.setLastUsedInfosToFeatureByChunk(ctx, features, req.EnvironmentNamespace, localizer); err != nil {
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
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err = validateCreateFeatureRequest(req.Command, localizer); err != nil {
		return nil, err
	}
	feature, err := domain.NewFeature(
		req.Command.Id,
		req.Command.Name,
		req.Command.Description,
		req.Command.VariationType,
		req.Command.Variations,
		req.Command.Tags,
		int(req.Command.DefaultOnVariationIndex.Value),
		int(req.Command.DefaultOffVariationIndex.Value),
		editor.Email,
	)
	if err != nil {
		return nil, err
	}
	var handler *command.FeatureCommandHandler = command.NewEmptyFeatureCommandHandler()
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		if err := s.upsertTags(ctx, tx, req.Command.Tags, req.EnvironmentNamespace); err != nil {
			return err
		}

		featureStorage := v2fs.NewFeatureStorage(tx)
		if err := featureStorage.CreateFeature(ctx, feature, req.EnvironmentNamespace); err != nil {
			s.logger.Error(
				"Failed to store feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return err
		}
		handler = command.NewFeatureCommandHandler(editor, feature, req.EnvironmentNamespace, "")
		if err := handler.Handle(ctx, req.Command); err != nil {
			s.logger.Error(
				"Failed to create feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		if err == v2fs.ErrFeatureAlreadyExists {
			dt, err := statusAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to create feature",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	err = s.refreshFeaturesCache(ctx, req.EnvironmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to refresh features cache",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if errs := s.publishDomainEvents(ctx, handler.Events); len(errs) > 0 {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Any("errors", errs),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &featureproto.CreateFeatureResponse{}, nil
}

func (s *FeatureService) UpdateFeatureDetails(
	ctx context.Context,
	req *featureproto.UpdateFeatureDetailsRequest,
) (*featureproto.UpdateFeatureDetailsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		dt, err := statusMissingID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if err := s.validateFeatureStatus(ctx, req.Id, req.EnvironmentNamespace, localizer); err != nil {
		return nil, err
	}
	var handler *command.FeatureCommandHandler = command.NewEmptyFeatureCommandHandler()
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		featureStorage := v2fs.NewFeatureStorage(tx)
		feature, err := featureStorage.GetFeature(ctx, req.Id, req.EnvironmentNamespace)
		if err != nil {
			s.logger.Error(
				"Failed to get feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return err
		}
		handler = command.NewFeatureCommandHandler(editor, feature, req.EnvironmentNamespace, req.Comment)
		err = handler.Handle(ctx, &featureproto.IncrementFeatureVersionCommand{})
		if err != nil {
			s.logger.Error(
				"Failed to increment feature version",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return err
		}
		if req.RenameFeatureCommand != nil {
			err = handler.Handle(ctx, req.RenameFeatureCommand)
			if err != nil {
				s.logger.Error(
					"Failed to rename feature",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentNamespace", req.EnvironmentNamespace),
					)...,
				)
				return err
			}
		}
		if req.ChangeDescriptionCommand != nil {
			err = handler.Handle(ctx, req.ChangeDescriptionCommand)
			if err != nil {
				s.logger.Error(
					"Failed to change feature description",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentNamespace", req.EnvironmentNamespace),
					)...,
				)
				return err
			}
		}
		if req.AddTagCommands != nil {
			for i := range req.AddTagCommands {
				err = handler.Handle(ctx, req.AddTagCommands[i])
				if err != nil {
					s.logger.Error(
						"Failed to add tag to feature",
						log.FieldsFromImcomingContext(ctx).AddFields(
							zap.Error(err),
							zap.String("environmentNamespace", req.EnvironmentNamespace),
						)...,
					)
					return err
				}
			}
			tags := []string{}
			for _, c := range req.AddTagCommands {
				tags = append(tags, c.Tag)
			}
			if err := s.upsertTags(ctx, tx, tags, req.EnvironmentNamespace); err != nil {
				return err
			}
		}
		if req.RemoveTagCommands != nil {
			for i := range req.RemoveTagCommands {
				err = handler.Handle(ctx, req.RemoveTagCommands[i])
				if err != nil {
					s.logger.Error(
						"Failed to remove tag from feature",
						log.FieldsFromImcomingContext(ctx).AddFields(
							zap.Error(err),
							zap.String("environmentNamespace", req.EnvironmentNamespace),
						)...,
					)
					return err
				}
			}
		}
		err = featureStorage.UpdateFeature(ctx, feature, req.EnvironmentNamespace)
		if err != nil {
			s.logger.Error(
				"Failed to update feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = s.refreshFeaturesCache(ctx, req.EnvironmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to refresh features cache",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if errs := s.publishDomainEvents(ctx, handler.Events); len(errs) > 0 {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Any("errors", errs),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &featureproto.UpdateFeatureDetailsResponse{}, nil
}

func (s *FeatureService) existsRunningExperiment(
	ctx context.Context,
	featureID, environmentNamespace string,
) (bool, error) {
	experiments, err := s.listExperiments(ctx, environmentNamespace, featureID)
	if err != nil {
		s.logger.Error(
			"Failed to list experiments",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
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

func (s *FeatureService) existsRunningProgressiveRollout(
	ctx context.Context,
	featureID, environmentNamespace string,
) (bool, error) {
	progressiveRollouts, err := s.listProgressiveRollouts(ctx, environmentNamespace, featureID)
	if err != nil {
		s.logger.Error(
			"Failed to list progressiveRollouts",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
			)...,
		)
		return false, err
	}
	return containsRunningProgressiveRollout(progressiveRollouts), nil
}

func containsRunningProgressiveRollout(progressiveRollouts []*autoopsproto.ProgressiveRollout) bool {
	for _, p := range progressiveRollouts {
		dp := &autoopsdomain.ProgressiveRollout{
			ProgressiveRollout: p,
		}
		if !dp.IsFinished() {
			return true
		}
	}
	return false
}

func (s *FeatureService) listProgressiveRollouts(
	ctx context.Context,
	featureID, environmentNamespace string,
) ([]*autoopsproto.ProgressiveRollout, error) {
	progressiveRollouts := make([]*autoopsproto.ProgressiveRollout, 0)
	cursor := ""
	for {
		resp, err := s.autoOpsClient.ListProgressiveRollouts(
			ctx,
			&autoopsproto.ListProgressiveRolloutsRequest{
				EnvironmentNamespace: environmentNamespace,
				PageSize:             listRequestSize,
				Cursor:               cursor,
				FeatureIds:           []string{featureID},
			},
		)
		if err != nil {
			return nil, err
		}
		progressiveRollouts = append(progressiveRollouts, resp.ProgressiveRollouts...)
		size := len(progressiveRollouts)
		if size == 0 || size < listRequestSize {
			return progressiveRollouts, nil
		}
		cursor = resp.Cursor
	}
}

// FIXME: remove this API after the new console is released
// Deprecated
func (s *FeatureService) EnableFeature(
	ctx context.Context,
	req *featureproto.EnableFeatureRequest,
) (*featureproto.EnableFeatureResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	if err := validateEnableFeatureRequest(req, localizer); err != nil {
		return nil, err
	}
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.updateFeature(
		ctx,
		req.Command,
		req.Id,
		req.EnvironmentNamespace,
		req.Comment,
		localizer,
		editor,
	); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to enable feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
		}
		return nil, err
	}
	return &featureproto.EnableFeatureResponse{}, nil
}

// FIXME: remove this API after the new console is released
// Deprecated
func (s *FeatureService) DisableFeature(
	ctx context.Context,
	req *featureproto.DisableFeatureRequest,
) (*featureproto.DisableFeatureResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	if err := validateDisableFeatureRequest(req, localizer); err != nil {
		return nil, err
	}
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.updateFeature(
		ctx,
		req.Command,
		req.Id,
		req.EnvironmentNamespace,
		req.Comment,
		localizer,
		editor,
	); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to disable feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
		}
		return nil, err
	}
	return &featureproto.DisableFeatureResponse{}, nil
}

func (s *FeatureService) ArchiveFeature(
	ctx context.Context,
	req *featureproto.ArchiveFeatureRequest,
) (*featureproto.ArchiveFeatureResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	whereParts := []mysql.WherePart{
		mysql.NewFilter("archived", "=", false),
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace),
	}
	featureStorage := v2fs.NewFeatureStorage(s.mysqlClient)
	features, _, _, err := featureStorage.ListFeatures(
		ctx,
		whereParts,
		nil,
		mysql.QueryNoLimit,
		mysql.QueryNoOffset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list feature",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	if err := validateArchiveFeatureRequest(req, features, localizer); err != nil {
		return nil, err
	}
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.updateFeature(
		ctx,
		req.Command,
		req.Id,
		req.EnvironmentNamespace,
		req.Comment,
		localizer,
		editor,
	); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to archive feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
		}
		return nil, err
	}
	return &featureproto.ArchiveFeatureResponse{}, nil
}

func (s *FeatureService) UnarchiveFeature(
	ctx context.Context,
	req *featureproto.UnarchiveFeatureRequest,
) (*featureproto.UnarchiveFeatureResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	if err := validateUnarchiveFeatureRequest(req, localizer); err != nil {
		return nil, err
	}
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.updateFeature(
		ctx,
		req.Command,
		req.Id,
		req.EnvironmentNamespace,
		req.Comment,
		localizer,
		editor,
	); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to unarchive feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
		}
		return nil, err
	}
	return &featureproto.UnarchiveFeatureResponse{}, nil
}

func (s *FeatureService) DeleteFeature(
	ctx context.Context,
	req *featureproto.DeleteFeatureRequest,
) (*featureproto.DeleteFeatureResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	if err := validateDeleteFeatureRequest(req, localizer); err != nil {
		return nil, err
	}
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.updateFeature(
		ctx,
		req.Command,
		req.Id,
		req.EnvironmentNamespace,
		req.Comment,
		localizer,
		editor,
	); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to delete feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
		}
		return nil, err
	}
	return &featureproto.DeleteFeatureResponse{}, nil
}

func (s *FeatureService) updateFeature(
	ctx context.Context,
	cmd command.Command,
	id, environmentNamespace, comment string,
	localizer locale.Localizer,
	editor *eventproto.Editor,
) error {
	if id == "" {
		dt, err := statusMissingID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if cmd == nil {
		dt, err := statusMissingCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if err := s.validateFeatureStatus(ctx, id, environmentNamespace, localizer); err != nil {
		return err
	}
	var handler *command.FeatureCommandHandler = command.NewEmptyFeatureCommandHandler()
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		featureStorage := v2fs.NewFeatureStorage(tx)
		feature, err := featureStorage.GetFeature(ctx, id, environmentNamespace)
		if err != nil {
			s.logger.Error(
				"Failed to get feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			return err
		}
		handler = command.NewFeatureCommandHandler(editor, feature, environmentNamespace, comment)
		err = handler.Handle(ctx, &featureproto.IncrementFeatureVersionCommand{})
		if err != nil {
			s.logger.Error(
				"Failed to increment feature version",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			return err
		}
		if err := handler.Handle(ctx, cmd); err != nil {
			s.logger.Error(
				"Failed to handle command",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			return err
		}
		if err := featureStorage.UpdateFeature(ctx, feature, environmentNamespace); err != nil {
			s.logger.Error(
				"Failed to update feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		return s.convUpdateFeatureError(err, localizer)
	}
	err = s.refreshFeaturesCache(ctx, environmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to refresh features cache",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if errs := s.publishDomainEvents(ctx, handler.Events); len(errs) > 0 {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Any("errors", errs),
				zap.String("environmentNamespace", environmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *FeatureService) convUpdateFeatureError(err error, localizer locale.Localizer) error {
	switch err {
	case v2fs.ErrFeatureNotFound,
		v2fs.ErrFeatureUnexpectedAffectedRows,
		storage.ErrKeyNotFound:
		dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.NotFoundError),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	case domain.ErrAlreadyDisabled:
		dt, err := statusNothingChange.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.NothingToChange),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	case domain.ErrAlreadyEnabled:
		dt, err := statusNothingChange.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.NothingToChange),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	default:
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
}

func (s *FeatureService) UpdateFeatureVariations(
	ctx context.Context,
	req *featureproto.UpdateFeatureVariationsRequest,
) (*featureproto.UpdateFeatureVariationsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		dt, err := statusMissingID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if err := s.validateFeatureStatus(ctx, req.Id, req.EnvironmentNamespace, localizer); err != nil {
		return nil, err
	}
	commands := make([]command.Command, 0, len(req.Commands))
	for _, c := range req.Commands {
		cmd, err := command.UnmarshalCommand(c)
		if err != nil {
			s.logger.Error(
				"Failed to unmarshal command",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return nil, err
		}
		commands = append(commands, cmd)
	}
	var handler *command.FeatureCommandHandler = command.NewEmptyFeatureCommandHandler()
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace),
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		featureStorage := v2fs.NewFeatureStorage(tx)
		features, _, _, err := featureStorage.ListFeatures(
			ctx,
			whereParts,
			nil,
			mysql.QueryNoLimit,
			mysql.QueryNoOffset,
		)
		if err != nil {
			s.logger.Error(
				"Failed to list feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return err
		}
		for _, cmd := range commands {
			if err := validateFeatureVariationsCommand(features, cmd, localizer); err != nil {
				s.logger.Info(
					"Invalid argument",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentNamespace", req.EnvironmentNamespace),
					)...,
				)
				return err
			}
		}
		f, err := findFeature(features, req.Id, localizer)
		if err != nil {
			s.logger.Error(
				"Failed to find feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return err
		}
		feature := &domain.Feature{Feature: f}
		handler = command.NewFeatureCommandHandler(editor, feature, req.EnvironmentNamespace, req.Comment)
		err = handler.Handle(ctx, &featureproto.IncrementFeatureVersionCommand{})
		if err != nil {
			s.logger.Error(
				"Failed to increment feature version",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return err
		}
		for _, cmd := range commands {
			err = handler.Handle(ctx, cmd)
			if err != nil {
				// TODO: make this error log more specific.
				s.logger.Error(
					"Failed to handle command",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentNamespace", req.EnvironmentNamespace),
					)...,
				)
				return err
			}
		}
		err = featureStorage.UpdateFeature(ctx, feature, req.EnvironmentNamespace)
		if err != nil {
			s.logger.Error(
				"Failed to update feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = s.refreshFeaturesCache(ctx, req.EnvironmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to refresh features cache",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if errs := s.publishDomainEvents(ctx, handler.Events); len(errs) > 0 {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Any("errors", errs),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &featureproto.UpdateFeatureVariationsResponse{}, nil
}

func (s *FeatureService) publishDomainEvents(ctx context.Context, events []*eventproto.Event) map[string]error {
	messages := make([]publisher.Message, 0, len(events))
	for _, event := range events {
		messages = append(messages, event)
	}
	return s.domainPublisher.PublishMulti(ctx, messages)
}

func (s *FeatureService) UpdateFeatureTargeting(
	ctx context.Context,
	req *featureproto.UpdateFeatureTargetingRequest,
) (*featureproto.UpdateFeatureTargetingResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		dt, err := statusMissingID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	commands := make([]command.Command, 0, len(req.Commands))
	for _, c := range req.Commands {
		cmd, err := command.UnmarshalCommand(c)
		if err != nil {
			s.logger.Error(
				"Failed to unmarshal command",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return nil, err
		}
		commands = append(commands, cmd)
	}
	if err := s.validateFeatureStatus(ctx, req.Id, req.EnvironmentNamespace, localizer); err != nil {
		return nil, err
	}
	// TODO: clean this up.
	// Problem: Changes in the UI should be atomic meaning either all or no changes will be made.
	// This means a transaction spanning all changes is needed.
	// Also:
	// Normally each command should be usable alone (load the feature from the repository change it and save it).
	// Also here because many commands are run sequentially they all expect the same version of the feature.
	var handler *command.FeatureCommandHandler = command.NewEmptyFeatureCommandHandler()
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		whereParts := []mysql.WherePart{
			mysql.NewFilter("deleted", "=", false),
			mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace),
		}
		featureStorage := v2fs.NewFeatureStorage(tx)
		features, _, _, err := featureStorage.ListFeatures(
			ctx,
			whereParts,
			nil,
			mysql.QueryNoLimit,
			mysql.QueryNoOffset,
		)
		if err != nil {
			s.logger.Error(
				"Failed to list feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return err
		}
		f, err := findFeature(features, req.Id, localizer)
		if err != nil {
			s.logger.Error(
				"Failed to find feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return err
		}
		for _, cmd := range commands {
			if err := validateFeatureTargetingCommand(features, f, cmd, localizer); err != nil {
				s.logger.Info(
					"Invalid argument",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentNamespace", req.EnvironmentNamespace),
					)...,
				)
				return err
			}
		}
		feature := &domain.Feature{Feature: f}
		handler = command.NewFeatureCommandHandler(
			editor,
			feature,
			req.EnvironmentNamespace,
			req.Comment,
		)
		err = handler.Handle(ctx, &featureproto.IncrementFeatureVersionCommand{})
		if err != nil {
			s.logger.Error(
				"Failed to increment feature version",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return err
		}
		for _, cmd := range commands {
			err = handler.Handle(ctx, cmd)
			if err != nil {
				// TODO: same as above. Make it more specific.
				s.logger.Error(
					"Failed to handle command",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentNamespace", req.EnvironmentNamespace),
					)...,
				)
				return err
			}
		}
		err = featureStorage.UpdateFeature(ctx, feature, req.EnvironmentNamespace)
		if err != nil {
			s.logger.Error(
				"Failed to update feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = s.refreshFeaturesCache(ctx, req.EnvironmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to refresh features cache",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if errs := s.publishDomainEvents(ctx, handler.Events); len(errs) > 0 {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Any("errors", errs),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &featureproto.UpdateFeatureTargetingResponse{}, nil
}

func findFeature(fs []*featureproto.Feature, id string, localizer locale.Localizer) (*featureproto.Feature, error) {
	for _, f := range fs {
		if f.Id == id {
			return f, nil
		}
	}
	dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
		Locale:  localizer.GetLocale(),
		Message: localizer.MustLocalize(locale.InternalServerError),
	})
	if err != nil {
		return nil, statusInternal.Err()
	}
	return nil, dt.Err()
}

func (s *FeatureService) evaluateFeatures(
	ctx context.Context,
	features []*featureproto.Feature,
	user *userproto.User,
	environmentNamespace string,
	tag string,
	localizer locale.Localizer,
) (*featureproto.UserEvaluations, error) {
	mapIDs := make(map[string]struct{})
	for _, f := range features {
		feature := &domain.Feature{Feature: f}
		for _, id := range feature.ListSegmentIDs() {
			mapIDs[id] = struct{}{}
		}
	}
	mapSegmentUsers, err := s.listSegmentUsers(ctx, mapIDs, environmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to list segments",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
				zap.String("userId", user.Id),
				zap.String("tag", tag),
			)...,
		)
		return nil, err
	}
	userEvaluations, err := domain.EvaluateFeatures(features, user, mapSegmentUsers, tag)
	if err != nil {
		s.logger.Error(
			"Failed to evaluate",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
				zap.String("userId", user.Id),
				zap.String("tag", tag),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return userEvaluations, nil
}

func (s *FeatureService) getFeatures(
	ctx context.Context,
	environmentNamespace string,
) ([]*featureproto.Feature, error) {
	features, err := s.featuresCache.Get(environmentNamespace)
	if err == nil {
		return features.Features, nil
	}
	s.logger.Info(
		"No cached data for Features",
		log.FieldsFromImcomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
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
		"",
		featureproto.ListFeaturesRequest_DEFAULT,
		featureproto.ListFeaturesRequest_ASC,
		environmentNamespace,
	)
	if err != nil {
		s.logger.Error(
			"Failed to retrive features from storage",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
			)...,
		)
		return nil, err
	}
	return fs, nil
}

func (s *FeatureService) listSegmentUsers(
	ctx context.Context,
	mapSegmentIDs map[string]struct{},
	environmentNamespace string,
) (map[string][]*featureproto.SegmentUser, error) {
	if len(mapSegmentIDs) == 0 {
		return nil, nil
	}
	users := make(map[string][]*featureproto.SegmentUser)
	for segmentID := range mapSegmentIDs {
		s, err, _ := s.flightgroup.Do(
			s.segmentFlightID(environmentNamespace, segmentID),
			func() (interface{}, error) {
				return s.getSegmentUsers(ctx, segmentID, environmentNamespace)
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

func (s *FeatureService) segmentFlightID(environmentNamespace, segmentID string) string {
	return fmt.Sprintf("%s:%s", environmentNamespace, segmentID)
}

func (s *FeatureService) getSegmentUsers(
	ctx context.Context,
	segmentID, environmentNamespace string,
) ([]*featureproto.SegmentUser, error) {
	segmentUsers, err := s.segmentUsersCache.Get(segmentID, environmentNamespace)
	if err == nil {
		return segmentUsers.Users, nil
	}
	s.logger.Info(
		"No cached data for SegmentUsers",
		log.FieldsFromImcomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("segmentId", segmentID),
		)...,
	)
	req := &featureproto.ListSegmentUsersRequest{
		SegmentId:            segmentID,
		EnvironmentNamespace: environmentNamespace,
	}
	res, err := s.ListSegmentUsers(ctx, req)
	if err != nil {
		s.logger.Error(
			"Failed to retrieve segment users from storage",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
				zap.String("segmentId", segmentID),
			)...,
		)
		return nil, err
	}
	su := &featureproto.SegmentUsers{
		SegmentId: segmentID,
		Users:     res.Users,
	}
	if err := s.segmentUsersCache.Put(su, environmentNamespace); err != nil {
		s.logger.Error(
			"Failed to cache segment users",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
				zap.String("segmentId", segmentID),
			)...,
		)
	}
	return res.Users, nil
}

func (s *FeatureService) setLastUsedInfosToFeatureByChunk(
	ctx context.Context,
	features []*featureproto.Feature,
	environmentNamespace string,
	localizer locale.Localizer,
) error {
	for i := 0; i < len(features); i += getMultiChunkSize {
		end := i + getMultiChunkSize
		if end > len(features) {
			end = len(features)
		}
		if err := s.setLastUsedInfosToFeature(ctx, features[i:end], environmentNamespace, localizer); err != nil {
			return err
		}
	}
	return nil
}

func (s *FeatureService) setLastUsedInfosToFeature(
	ctx context.Context,
	features []*featureproto.Feature,
	environmentNamespace string,
	localizer locale.Localizer,
) error {
	ids := make([]string, 0, len(features))
	for _, f := range features {
		ids = append(ids, domain.FeatureLastUsedInfoID(f.Id, f.Version))
	}
	storage := v2fs.NewFeatureLastUsedInfoStorage(s.mysqlClient)
	fluiList, err := storage.GetFeatureLastUsedInfos(ctx, ids, environmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to get feature last used infos",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	for _, f := range fluiList {
		for _, feature := range features {
			if feature.Id == f.FeatureLastUsedInfo.FeatureId {
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
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateEvaluateFeatures(req, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	fs, err, _ := s.flightgroup.Do(
		req.EnvironmentNamespace,
		func() (interface{}, error) {
			return s.getFeatures(ctx, req.EnvironmentNamespace)
		},
	)
	if err != nil {
		s.logger.Error(
			"Failed to list features",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	// If the feature ID is set in the request, it will evaluate a single feature.
	features, err := s.getTargetFeatures(fs.([]*featureproto.Feature), req.FeatureId, localizer)
	if err != nil {
		s.logger.Error(
			"Failed to get target features",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	userEvaluations, err := s.evaluateFeatures(ctx, features, req.User, req.EnvironmentNamespace, req.Tag, localizer)
	if err != nil {
		s.logger.Error(
			"Failed to evaluate features",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	// If the feature ID is set, it will return a single evaluation
	if req.FeatureId != "" {
		eval, err := s.findEvaluation(userEvaluations.Evaluations, req.FeatureId)
		if err != nil {
			s.logger.Error(
				"Failed to find evaluation",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
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

func (s *FeatureService) getTargetFeatures(
	fs []*featureproto.Feature,
	id string,
	localizer locale.Localizer,
) ([]*featureproto.Feature, error) {
	if id == "" {
		return fs, nil
	}
	feature, err := findFeature(fs, id, localizer)
	if err != nil {
		return nil, err
	}
	if len(feature.Prerequisites) > 0 {
		// If we select only the prerequisite feature flags, we have to get them recursively.
		// Thus, we evaluate all features here to avoid complex logic.
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
	environmentNamespace, featureID string,
) ([]*experimentproto.Experiment, error) {
	experiments := []*experimentproto.Experiment{}
	cursor := ""
	for {
		resp, err := s.experimentClient.ListExperiments(ctx, &experimentproto.ListExperimentsRequest{
			FeatureId:            featureID,
			PageSize:             listRequestSize,
			Cursor:               cursor,
			EnvironmentNamespace: environmentNamespace,
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
	localizer := locale.NewLocalizer(ctx)
	if err := validateCloneFeatureRequest(req, localizer); err != nil {
		return nil, err
	}
	editor, err := s.checkRole(
		ctx,
		accountproto.AccountV2_Role_Environment_EDITOR,
		req.Command.EnvironmentNamespace,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	featureStorage := v2fs.NewFeatureStorage(s.mysqlClient)
	f, err := featureStorage.GetFeature(ctx, req.Id, req.EnvironmentNamespace)
	if err != nil {
		if err == v2fs.ErrFeatureNotFound {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to get feature",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	domainFeature := &domain.Feature{
		Feature: f.Feature,
	}
	feature, err := domainFeature.Clone(editor.Email)
	if err != nil {
		return nil, err
	}
	var handler *command.FeatureCommandHandler = command.NewEmptyFeatureCommandHandler()
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		if err := featureStorage.CreateFeature(ctx, feature, req.Command.EnvironmentNamespace); err != nil {
			s.logger.Error(
				"Failed to store feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.Command.EnvironmentNamespace),
				)...,
			)
			return err
		}
		handler = command.NewFeatureCommandHandler(editor, feature, req.Command.EnvironmentNamespace, "")
		if err := handler.Handle(ctx, req.Command); err != nil {
			s.logger.Error(
				"Failed to clone feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNameSpace", req.Command.EnvironmentNamespace),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		if err == v2fs.ErrFeatureAlreadyExists {
			dt, err := statusAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to clone feature",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.Command.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	err = s.refreshFeaturesCache(ctx, req.EnvironmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to refresh features cache",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if errs := s.publishDomainEvents(ctx, handler.Events); len(errs) > 0 {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Any("errors", errs),
				zap.String("environmentNameSpace", req.Command.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &featureproto.CloneFeatureResponse{}, nil
}

func (s *FeatureService) refreshFeaturesCache(ctx context.Context, environmentNamespace string) error {
	fs, _, _, err := s.listFeatures(
		ctx,
		mysql.QueryNoLimit,
		"",
		nil,
		"",
		nil,
		nil,
		nil,
		"",
		featureproto.ListFeaturesRequest_DEFAULT,
		featureproto.ListFeaturesRequest_ASC,
		environmentNamespace,
	)
	if err != nil {
		return err
	}
	filtered := make([]*featureproto.Feature, 0)
	for _, f := range fs {
		ff := domain.Feature{Feature: f}
		if ff.IsDisabledAndOffVariationEmpty() {
			continue
		}
		// To keep the cache size small, we exclude feature flags archived more than thirty days ago.
		if ff.IsArchivedBeforeLastThirtyDays() {
			continue
		}
		filtered = append(filtered, f)
	}
	features := &featureproto.Features{
		Features: filtered,
	}
	if err := s.featuresCache.Put(features, environmentNamespace); err != nil {
		return err
	}
	s.logger.Info("Success to refresh features cache",
		zap.String("environmentNamespace", environmentNamespace),
		zap.Int("numberOfFeatures", len(fs)),
	)
	return nil
}
