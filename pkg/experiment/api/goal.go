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
	"regexp"
	"strconv"

	pb "github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/experiment/command"
	"github.com/bucketeer-io/bucketeer/pkg/experiment/domain"
	v2es "github.com/bucketeer-io/bucketeer/pkg/experiment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

var goalIDRegex = regexp.MustCompile("^[a-zA-Z0-9-]+$")

func (s *experimentService) GetGoal(ctx context.Context, req *proto.GetGoalRequest) (*proto.GetGoalResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		dt, err := statusGoalIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	goal, err := s.getGoalMySQL(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2es.ErrGoalNotFound) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
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
	err = s.mapConnectedOperations(ctx, []*proto.Goal{goal.Goal}, req.EnvironmentId)
	if err != nil {
		s.logger.Error("Failed to map connected operations", zap.Error(err))
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &proto.GetGoalResponse{Goal: goal.Goal}, nil
}

func (s *experimentService) getGoalMySQL(
	ctx context.Context,
	goalID, environmentId string,
) (*domain.Goal, error) {
	goalStorage := v2es.NewGoalStorage(s.mysqlClient)
	goal, err := goalStorage.GetGoal(ctx, goalID, environmentId)
	if err != nil {
		s.logger.Error(
			"Failed to get goal",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
				zap.String("goalId", goalID),
			)...,
		)
	}
	return goal, err
}

func (s *experimentService) ListGoals(
	ctx context.Context,
	req *proto.ListGoalsRequest,
) (*proto.ListGoalsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_id", "=", req.EnvironmentId),
	}
	if req.Archived != nil {
		whereParts = append(whereParts, mysql.NewFilter("archived", "=", req.Archived.Value))
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"id", "name", "description"}, req.SearchKeyword))
	}
	orders, err := s.newGoalListOrders(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
	goalStorage := v2es.NewGoalStorage(s.mysqlClient)
	goals, nextCursor, totalCount, err := goalStorage.ListGoals(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
		isInUseStatus,
		req.EnvironmentId,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list goals",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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
	err = s.mapConnectedOperations(ctx, goals, req.EnvironmentId)
	if err != nil {
		s.logger.Error("Failed to map connected operations", zap.Error(err))
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}

	return &proto.ListGoalsResponse{
		Goals:      goals,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *experimentService) newGoalListOrders(
	orderBy proto.ListGoalsRequest_OrderBy,
	orderDirection proto.ListGoalsRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case proto.ListGoalsRequest_DEFAULT,
		proto.ListGoalsRequest_NAME:
		column = "name"
	case proto.ListGoalsRequest_CREATED_AT:
		column = "created_at"
	case proto.ListGoalsRequest_UPDATED_AT:
		column = "updated_at"
	case proto.ListGoalsRequest_CONNECTION_TYPE:
		column = "connection_type"
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
	if orderDirection == proto.ListGoalsRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *experimentService) mapConnectedOperations(
	ctx context.Context,
	goals []*proto.Goal,
	environmentID string,
) error {
	listAutoOpRulesResp, err := s.autoOpsClient.ListAutoOpsRules(ctx, &autoopsproto.ListAutoOpsRulesRequest{
		EnvironmentId: environmentID,
	})
	if err != nil {
		return err
	}
	autoOpsRules := listAutoOpRulesResp.AutoOpsRules
	goalOpsMap := make(map[string][]*autoopsproto.AutoOpsRule)
	for _, rule := range autoOpsRules {
		for _, clause := range rule.Clauses {
			if clause.Clause.MessageIs(&autoopsproto.OpsEventRateClause{}) {
				c := &autoopsproto.OpsEventRateClause{}
				if err := clause.Clause.UnmarshalTo(c); err != nil {
					return err
				}
				if c.GoalId == "" {
					continue
				}
				goalOpsMap[c.GoalId] = append(goalOpsMap[c.GoalId], rule)
			}
		}
	}
	for _, goal := range goals {
		if ops, ok := goalOpsMap[goal.Id]; ok {
			goal.AutoOpsRules = ops
		}
	}
	return nil
}

func (s *experimentService) CreateGoal(
	ctx context.Context,
	req *proto.CreateGoalRequest,
) (*proto.CreateGoalResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if req.Command == nil {
		return s.createGoalNoCommand(ctx, req, editor, localizer)
	}
	if err := validateCreateGoalRequest(req, localizer); err != nil {
		return nil, err
	}
	goal, err := domain.NewGoal(req.Command.Id, req.Command.Name, req.Command.Description, req.Command.ConnectionType)
	if err != nil {
		s.logger.Error(
			"Failed to create a new goal",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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
		goalStorage := v2es.NewGoalStorage(tx)
		handler, err := command.NewGoalCommandHandler(editor, goal, s.publisher, req.EnvironmentId)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return goalStorage.CreateGoal(ctx, goal, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrGoalAlreadyExists) {
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
			"Failed to create goal",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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
	return &proto.CreateGoalResponse{
		Goal: goal.Goal,
	}, nil
}

func (s *experimentService) createGoalNoCommand(
	ctx context.Context,
	req *proto.CreateGoalRequest,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) (*proto.CreateGoalResponse, error) {
	if err := validateCreateGoalNoCommandRequest(req, localizer); err != nil {
		return nil, err
	}
	goal, err := domain.NewGoal(req.Id, req.Name, req.Description, req.ConnectionType)
	if err != nil {
		s.logger.Error(
			"Failed to create a new goal",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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
		goalStorage := v2es.NewGoalStorage(tx)
		prev := &domain.Goal{}
		e, err := domainevent.NewEvent(
			editor,
			eventproto.Event_GOAL,
			goal.Id,
			eventproto.Event_GOAL_CREATED,
			&eventproto.GoalCreatedEvent{
				Id:             goal.Id,
				Name:           goal.Name,
				Description:    goal.Description,
				ConnectionType: goal.ConnectionType,
				Deleted:        goal.Deleted,
				CreatedAt:      goal.CreatedAt,
				UpdatedAt:      goal.UpdatedAt,
			},
			req.EnvironmentId,
			goal.Goal,
			prev,
		)
		if err != nil {
			return err
		}
		if err := s.publisher.Publish(ctx, e); err != nil {
			return err
		}
		return goalStorage.CreateGoal(ctx, goal, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrGoalAlreadyExists) {
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
			"Failed to create goal",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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
	return &proto.CreateGoalResponse{
		Goal: goal.Goal,
	}, nil
}

func validateCreateGoalRequest(req *proto.CreateGoalRequest, localizer locale.Localizer) error {
	if req.Command.Id == "" {
		dt, err := statusGoalIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !goalIDRegex.MatchString(req.Command.Id) {
		dt, err := statusInvalidGoalID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "goal_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.Name == "" {
		dt, err := statusGoalNameRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateCreateGoalNoCommandRequest(req *proto.CreateGoalRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusGoalIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !goalIDRegex.MatchString(req.Id) {
		dt, err := statusInvalidGoalID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "goal_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Name == "" {
		dt, err := statusGoalNameRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *experimentService) UpdateGoal(
	ctx context.Context,
	req *proto.UpdateGoalRequest,
) (*proto.UpdateGoalResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if req.ChangeDescriptionCommand == nil && req.RenameCommand == nil {
		return s.updateGoalNoCommand(ctx, req, editor, localizer)
	}
	if req.Id == "" {
		dt, err := statusGoalIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	commands := make([]command.Command, 0)
	if req.RenameCommand != nil {
		commands = append(commands, req.RenameCommand)
	}
	if req.ChangeDescriptionCommand != nil {
		commands = append(commands, req.ChangeDescriptionCommand)
	}
	if len(commands) == 0 {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	err = s.updateGoal(
		ctx,
		editor,
		req.EnvironmentId,
		req.Id,
		commands,
		localizer,
	)
	if err != nil {
		s.logger.Error(
			"Failed to update goal",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	return &proto.UpdateGoalResponse{}, nil
}

func (s *experimentService) updateGoalNoCommand(
	ctx context.Context,
	req *proto.UpdateGoalRequest,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) (*proto.UpdateGoalResponse, error) {
	err := s.validateUpdateGoalNoCommandRequest(req, localizer)
	if err != nil {
		return nil, err
	}

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
	var updatedGoal *proto.Goal
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		goalStorage := v2es.NewGoalStorage(tx)
		goal, err := goalStorage.GetGoal(ctx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		updated, err := goal.Update(
			req.Name,
			req.Description,
			req.Archived,
		)
		if err != nil {
			return err
		}
		updatedGoal = updated.Goal

		var event pb.Message
		if req.Archived != nil && req.Archived.Value {
			event = &eventproto.GoalArchivedEvent{Id: goal.Id}
		} else {
			event = &eventproto.GoalUpdatedEvent{
				Id:          goal.Id,
				Name:        req.Name,
				Description: req.Description,
			}
		}
		e, err := domainevent.NewEvent(
			editor,
			eventproto.Event_GOAL,
			goal.Id,
			eventproto.Event_GOAL_UPDATED,
			event,
			req.EnvironmentId,
			updated.Goal,
			goal.Goal,
		)
		if err != nil {
			return err
		}
		if err = s.publisher.Publish(ctx, e); err != nil {
			return err
		}
		return goalStorage.UpdateGoal(ctx, updated, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrGoalNotFound) || errors.Is(err, v2es.ErrGoalUnexpectedAffectedRows) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
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
	return &proto.UpdateGoalResponse{
		Goal: updatedGoal,
	}, nil
}

func (s *experimentService) validateUpdateGoalNoCommandRequest(
	req *proto.UpdateGoalRequest,
	localizer locale.Localizer,
) error {
	if req.Id == "" {
		dt, err := statusGoalIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Name != nil && req.Name.Value == "" {
		dt, err := statusGoalNameRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *experimentService) ArchiveGoal(
	ctx context.Context,
	req *proto.ArchiveGoalRequest,
) (*proto.ArchiveGoalResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		dt, err := statusGoalIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	err = s.updateGoal(
		ctx,
		editor,
		req.EnvironmentId,
		req.Id,
		[]command.Command{req.Command},
		localizer,
	)
	if err != nil {
		s.logger.Error(
			"Failed to archive goal",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	return &proto.ArchiveGoalResponse{}, nil
}

func (s *experimentService) DeleteGoal(
	ctx context.Context,
	req *proto.DeleteGoalRequest,
) (*proto.DeleteGoalResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		dt, err := statusGoalIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	err = s.updateGoal(
		ctx,
		editor,
		req.EnvironmentId,
		req.Id,
		[]command.Command{req.Command},
		localizer,
	)
	if err != nil {
		s.logger.Error(
			"Failed to delete goal",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	return &proto.DeleteGoalResponse{}, nil
}

func (s *experimentService) updateGoal(
	ctx context.Context,
	editor *eventproto.Editor,
	environmentId, goalID string,
	commands []command.Command,
	localizer locale.Localizer,
) error {
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
		goalStorage := v2es.NewGoalStorage(tx)
		goal, err := goalStorage.GetGoal(ctx, goalID, environmentId)
		if err != nil {
			return err
		}
		handler, err := command.NewGoalCommandHandler(editor, goal, s.publisher, environmentId)
		if err != nil {
			return err
		}
		for _, command := range commands {
			if err := handler.Handle(ctx, command); err != nil {
				return err
			}
		}
		return goalStorage.UpdateGoal(ctx, goal, environmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrGoalNotFound) || errors.Is(err, v2es.ErrGoalUnexpectedAffectedRows) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
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
