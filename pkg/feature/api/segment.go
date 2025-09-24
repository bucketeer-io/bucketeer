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

package api

import (
	"context"
	"errors"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	pkgErr "github.com/bucketeer-io/bucketeer/pkg/error"
	"github.com/bucketeer-io/bucketeer/pkg/feature/command"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	errFeatureIDsNotFound = pkgErr.NewErrorNotFound(pkgErr.FeaturePackageName, "segment feature ids not found", "segment")
	errFeatureNotFound    = pkgErr.NewErrorNotFound(pkgErr.FeaturePackageName, "segment feature not found", "segment")
)

func (s *FeatureService) CreateSegment(
	ctx context.Context,
	req *featureproto.CreateSegmentRequest,
) (*featureproto.CreateSegmentResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if req.Command == nil {
		return s.createSegmentNoCommand(ctx, req, editor, localizer)
	}
	if err = validateCreateSegmentRequest(req.Command, localizer); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	segment, err := domain.NewSegment(req.Command.Name, req.Command.Description)
	if err != nil {
		s.logger.Error(
			"Failed to create segment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		if err := s.segmentStorage.CreateSegment(contextWithTx, segment, req.EnvironmentId); err != nil {
			s.logger.Error(
				"Failed to store segment",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return err
		}
		handler, err := command.NewSegmentCommandHandler(
			editor,
			segment,
			s.domainPublisher,
			req.EnvironmentId,
		)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			s.logger.Error(
				"Failed to handle command",
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
		if errors.Is(err, v2fs.ErrSegmentAlreadyExists) {
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
			"Failed to create segment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &featureproto.CreateSegmentResponse{
		Segment: segment.Segment,
	}, nil
}

func (s *FeatureService) createSegmentNoCommand(
	ctx context.Context,
	req *featureproto.CreateSegmentRequest,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) (*featureproto.CreateSegmentResponse, error) {
	if err := validateCreateSegmentNoCommandRequest(req, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	segment, err := domain.NewSegment(req.Name, req.Description)
	if err != nil {
		s.logger.Error(
			"Failed to create segment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		if err := s.segmentStorage.CreateSegment(contextWithTx, segment, req.EnvironmentId); err != nil {
			s.logger.Error(
				"Failed to store segment",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return err
		}
		e, err := domainevent.NewEvent(
			editor,
			eventproto.Event_SEGMENT,
			segment.Id,
			eventproto.Event_SEGMENT_CREATED,
			&eventproto.SegmentCreatedEvent{
				Id:          segment.Id,
				Name:        segment.Name,
				Description: segment.Description,
			},
			req.EnvironmentId,
			segment.Segment,
			nil,
		)
		if err != nil {
			return nil
		}
		return s.domainPublisher.Publish(ctx, e)
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrSegmentAlreadyExists) {
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
			"Failed to create segment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &featureproto.CreateSegmentResponse{
		Segment: segment.Segment,
	}, nil
}

func (s *FeatureService) DeleteSegment(
	ctx context.Context,
	req *featureproto.DeleteSegmentRequest,
) (*featureproto.DeleteSegmentResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if req.Command == nil {
		return s.deleteSegmentNoCommand(ctx, req, editor, localizer)
	}
	if err := validateDeleteSegmentRequest(req, localizer); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	if err := s.checkSegmentInUse(ctx, req.Id, req.EnvironmentId, localizer); err != nil {
		return nil, err
	}
	if _, err := s.updateSegment(
		ctx,
		editor,
		[]command.Command{req.Command},
		req.Id,
		req.EnvironmentId,
		localizer,
	); err != nil {
		return nil, err
	}
	return &featureproto.DeleteSegmentResponse{}, nil
}

func (s *FeatureService) deleteSegmentNoCommand(
	ctx context.Context,
	req *featureproto.DeleteSegmentRequest,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) (*featureproto.DeleteSegmentResponse, error) {
	if err := validateDeleteSegmentRequest(req, localizer); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	if err := s.checkSegmentInUse(ctx, req.Id, req.EnvironmentId, localizer); err != nil {
		return nil, err
	}
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		segment, _, err := s.segmentStorage.GetSegment(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			s.logger.Error(
				"Failed to get segment",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return err
		}
		event, err := domainevent.NewEvent(
			editor,
			eventproto.Event_SEGMENT,
			segment.Id,
			eventproto.Event_SEGMENT_DELETED,
			&eventproto.SegmentDeletedEvent{
				Id: segment.Id,
			},
			req.EnvironmentId,
			nil,             // Current state: entity no longer exists
			segment.Segment, // Previous state: what was deleted
		)
		if err != nil {
			return nil
		}
		if err := s.domainPublisher.Publish(ctx, event); err != nil {
			return err
		}
		return s.segmentStorage.DeleteSegment(contextWithTx, segment.Id)
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrSegmentNotFound) || errors.Is(err, v2fs.ErrSegmentUnexpectedAffectedRows) {
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
			"Failed to delete segment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &featureproto.DeleteSegmentResponse{}, nil
}

func (s *FeatureService) checkSegmentInUse(
	ctx context.Context,
	segmentID, environmentId string,
	localizer locale.Localizer,
) error {
	featureStorage := v2fs.NewFeatureStorage(s.mysqlClient)
	filters := []*mysql.FilterV2{
		{
			Column:   "deleted",
			Operator: mysql.OperatorEqual,
			Value:    false,
		},
		{
			Column:   "feature.environment_id",
			Operator: mysql.OperatorEqual,
			Value:    environmentId,
		},
	}
	options := &mysql.ListOptions{
		Filters:     filters,
		JSONFilters: nil,
		Orders:      nil,
		NullFilters: nil,
		InFilters:   nil,
		Limit:       mysql.QueryNoLimit,
		Offset:      mysql.QueryNoOffset,
	}
	features, _, _, err := featureStorage.ListFeatures(ctx, options)
	if err != nil {
		s.logger.Error(
			"Failed to list features",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		return api.NewGRPCStatus(err).Err()
	}
	if s.containsInRules(segmentID, features) {
		s.logger.Warn(
			"Segment User in use",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("segmentId", segmentID),
				zap.String("environmentId", environmentId),
			)...,
		)
		dt, err := statusSegmentInUse.WithDetails(&errdetails.LocalizedMessage{
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

func (s *FeatureService) containsInRules(segmentID string, features []*featureproto.Feature) bool {
	for _, f := range features {
		for _, r := range f.Rules {
			for _, c := range r.Clauses {
				if c.Operator == featureproto.Clause_SEGMENT {
					for _, id := range c.Values {
						if segmentID == id {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func (s *FeatureService) UpdateSegment(
	ctx context.Context,
	req *featureproto.UpdateSegmentRequest,
) (*featureproto.UpdateSegmentResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		s.logger.Error(
			"Permission denied",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	if req.Commands == nil {
		return s.updateSegmentNoCommand(ctx, req, editor, localizer)
	}
	commands := make([]command.Command, 0, len(req.Commands))
	for _, c := range req.Commands {
		cmd, err := command.UnmarshalCommand(c)
		if err != nil {
			s.logger.Error(
				"Failed to unmarshal command",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return nil, api.NewGRPCStatus(err).Err()
		}
		commands = append(commands, cmd)
	}
	if err := validateUpdateSegment(req.Id, commands, localizer); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	segment, err := s.updateSegment(
		ctx,
		editor,
		commands,
		req.Id,
		req.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	return &featureproto.UpdateSegmentResponse{
		Segment: segment,
	}, nil
}

func (s *FeatureService) updateSegmentNoCommand(
	ctx context.Context,
	req *featureproto.UpdateSegmentRequest,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) (*featureproto.UpdateSegmentResponse, error) {
	err := validateUpdateSegmentNoCommand(req, localizer)
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
	var updatedSegment *featureproto.Segment
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		segment, _, err := s.segmentStorage.GetSegment(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			s.logger.Error(
				"Failed to get segment",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return err
		}
		updated, err := segment.UpdateSegment(req.Name, req.Description)
		if err != nil {
			s.logger.Error(
				"Failed to update segment",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return err
		}
		updatedSegment = updated.Segment
		e, err := domainevent.NewEvent(
			editor,
			eventproto.Event_SEGMENT,
			req.Id,
			eventproto.Event_SEGMENT_NAME_CHANGED,
			&eventproto.SegmentUpdatedEvent{
				Id:          req.Id,
				Name:        req.Name,
				Description: req.Description,
			},
			req.EnvironmentId,
			updated.Segment,
			segment.Segment,
		)
		if err != nil {
			return err
		}
		if err := s.domainPublisher.Publish(ctx, e); err != nil {
			return err
		}
		return s.segmentStorage.UpdateSegment(contextWithTx, updated, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrSegmentNotFound) || errors.Is(err, v2fs.ErrSegmentUnexpectedAffectedRows) {
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
			"Failed to update segment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &featureproto.UpdateSegmentResponse{
		Segment: updatedSegment,
	}, nil
}

func (s *FeatureService) updateSegment(
	ctx context.Context,
	editor *eventproto.Editor,
	commands []command.Command,
	segmentID, environmentId string,
	localizer locale.Localizer,
) (*featureproto.Segment, error) {
	var updatedSegment *featureproto.Segment
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		segment, _, err := s.segmentStorage.GetSegment(contextWithTx, segmentID, environmentId)
		if err != nil {
			s.logger.Error(
				"Failed to get segment",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return err
		}
		handler, err := command.NewSegmentCommandHandler(
			editor,
			segment,
			s.domainPublisher,
			environmentId,
		)
		if err != nil {
			return err
		}
		for _, cmd := range commands {
			if err := handler.Handle(ctx, cmd); err != nil {
				s.logger.Error(
					"Failed to handle command",
					log.FieldsFromIncomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentId", environmentId),
					)...,
				)
				return err
			}
		}
		updatedSegment = segment.Segment
		return s.segmentStorage.UpdateSegment(contextWithTx, segment, environmentId)
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrSegmentNotFound) || errors.Is(err, v2fs.ErrSegmentUnexpectedAffectedRows) {
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
			"Failed to update segment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return updatedSegment, nil
}

func (s *FeatureService) GetSegment(
	ctx context.Context,
	req *featureproto.GetSegmentRequest,
) (*featureproto.GetSegmentResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetSegmentRequest(req, localizer); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	segment, featureIDs, err := s.segmentStorage.GetSegment(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		if err == v2fs.ErrSegmentNotFound {
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
			"Failed to get segment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	if err := s.injectFeaturesIntoSegments(
		ctx,
		[]*featureproto.Segment{
			segment.Segment,
		},
		map[string][]string{
			segment.Id: featureIDs,
		},
		req.EnvironmentId,
	); err != nil {
		s.logger.Error(
			"Failed to inject features into segments",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &featureproto.GetSegmentResponse{Segment: segment.Segment}, nil

}

func (s *FeatureService) ListSegments(
	ctx context.Context,
	req *featureproto.ListSegmentsRequest,
) (*featureproto.ListSegmentsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateListSegmentsRequest(req, localizer); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	filters := []*mysql.FilterV2{
		{
			Column:   "seg.deleted",
			Operator: mysql.OperatorEqual,
			Value:    false,
		},
		{
			Column:   "seg.environment_id",
			Operator: mysql.OperatorEqual,
			Value:    req.EnvironmentId,
		},
	}
	if req.Status != nil {
		filters = append(filters, &mysql.FilterV2{
			Column:   "seg.status",
			Operator: mysql.OperatorEqual,
			Value:    req.Status.Value,
		})
	}
	var searchQuery *mysql.SearchQuery
	if req.SearchKeyword != "" {
		searchQuery = &mysql.SearchQuery{
			Columns: []string{"seg.name", "seg.description"},
			Keyword: req.SearchKeyword,
		}
	}
	orders, err := s.newSegmentListOrders(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
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
	var isInUseStatus *bool
	if req.IsInUseStatus != nil {
		isInUseStatus = &req.IsInUseStatus.Value
	}
	options := &mysql.ListOptions{
		Limit:       limit,
		Offset:      offset,
		Filters:     filters,
		NullFilters: nil,
		JSONFilters: nil,
		InFilters:   nil,
		SearchQuery: searchQuery,
		Orders:      orders,
	}
	segments, nextCursor, totalCount, featureIDsMap, err := s.segmentStorage.ListSegments(
		ctx,
		options,
		isInUseStatus,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list segments",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	if err := s.injectFeaturesIntoSegments(
		ctx,
		segments,
		featureIDsMap,
		req.EnvironmentId,
	); err != nil {
		s.logger.Error(
			"Failed to inject features into segments",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &featureproto.ListSegmentsResponse{
		Segments:   segments,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *FeatureService) newSegmentListOrders(
	orderBy featureproto.ListSegmentsRequest_OrderBy,
	orderDirection featureproto.ListSegmentsRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case featureproto.ListSegmentsRequest_DEFAULT,
		featureproto.ListSegmentsRequest_NAME:
		column = "seg.name"
	case featureproto.ListSegmentsRequest_CREATED_AT:
		column = "seg.created_at"
	case featureproto.ListSegmentsRequest_UPDATED_AT:
		column = "seg.updated_at"
	case featureproto.ListSegmentsRequest_USERS:
		column = "seg.included_user_count"
	case featureproto.ListSegmentsRequest_CONNECTIONS:
		column = "feature_ids"
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
	if orderDirection == featureproto.ListSegmentsRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *FeatureService) injectFeaturesIntoSegments(
	ctx context.Context,
	segments []*featureproto.Segment,
	featureIDsMap map[string][]string,
	environmentId string,
) error {
	allFeatures, err := s.listAllFeatures(
		ctx,
		environmentId,
	)
	if err != nil {
		return err
	}
	for _, segment := range segments {
		featureIDs, ok := featureIDsMap[segment.Id]
		if !ok {
			return errFeatureIDsNotFound
		}
		features := make([]*featureproto.Feature, 0, len(featureIDs))
		for _, fID := range featureIDs {
			feature, ok := allFeatures[fID]
			if !ok {
				return errFeatureNotFound
			}
			features = append(features, feature)
		}
		segment.Features = features
	}
	return nil
}

func (s *FeatureService) listAllFeatures(
	ctx context.Context,
	environmentId string,
) (map[string]*featureproto.Feature, error) {
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
		environmentId,
	)
	if err != nil {
		return nil, err
	}
	featuresMap := make(map[string]*featureproto.Feature)
	for _, f := range fs {
		featuresMap[f.Id] = f
	}
	return featuresMap, nil
}
