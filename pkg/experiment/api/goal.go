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

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/command"
	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

var goalIDRegex = regexp.MustCompile("^[a-zA-Z0-9-]+$")

func (s *experimentService) GetGoal(ctx context.Context, req *proto.GetGoalRequest) (*proto.GetGoalResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		return nil, statusGoalIDRequired.Err()
	}
	goal, err := s.getGoalMySQL(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2es.ErrGoalNotFound) {
			return nil, statusGoalNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.mapConnectedOperations(ctx, []*proto.Goal{goal.Goal}, req.EnvironmentId)
	if err != nil {
		s.logger.Error("Failed to map connected operations", zap.Error(err))
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.GetGoalResponse{Goal: goal.Goal}, nil
}

func (s *experimentService) getGoalMySQL(
	ctx context.Context,
	goalID, environmentId string,
) (*domain.Goal, error) {
	goal, err := s.goalStorage.GetGoal(ctx, goalID, environmentId)
	if err != nil {
		s.logger.Error(
			"Failed to get goal",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
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
	if req.ConnectionType != proto.Goal_UNKNOWN {
		whereParts = append(whereParts, mysql.NewFilter("connection_type", "=", req.ConnectionType))
	}
	orders, err := s.newGoalListOrders(req.OrderBy, req.OrderDirection)
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
		return nil, statusInvalidCursor.Err()
	}
	var isInUseStatus *bool
	if req.IsInUseStatus != nil {
		isInUseStatus = &req.IsInUseStatus.Value
	}
	goals, nextCursor, totalCount, err := s.goalStorage.ListGoals(
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.mapConnectedOperations(ctx, goals, req.EnvironmentId)
	if err != nil {
		s.logger.Error("Failed to map connected operations", zap.Error(err))
		return nil, api.NewGRPCStatus(err).Err()
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
		return nil, statusInvalidOrderBy.Err()
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
	goalOpsMap := make(map[string][]*proto.Goal_AutoOpsRuleReference)
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
				goalOpsMap[c.GoalId] = append(goalOpsMap[c.GoalId], &proto.Goal_AutoOpsRuleReference{
					Id:            rule.Id,
					FeatureId:     rule.FeatureId,
					FeatureName:   rule.FeatureName,
					AutoOpsStatus: rule.AutoOpsStatus,
				})
			}
		}
	}
	for _, goal := range goals {
		if ops, ok := goalOpsMap[goal.Id]; ok {
			goal.AutoOpsRules = ops
			goal.IsInUseStatus = true
		}
	}
	return nil
}

func (s *experimentService) CreateGoal(
	ctx context.Context,
	req *proto.CreateGoalRequest,
) (*proto.CreateGoalResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if req.Command == nil {
		return s.createGoalNoCommand(ctx, req, editor)
	}
	if err := validateCreateGoalRequest(req); err != nil {
		return nil, err
	}
	goal, err := domain.NewGoal(req.Command.Id, req.Command.Name, req.Command.Description, req.Command.ConnectionType)
	if err != nil {
		s.logger.Error(
			"Failed to create a new goal",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		handler, err := command.NewGoalCommandHandler(editor, goal, s.publisher, req.EnvironmentId)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return s.goalStorage.CreateGoal(ctxWithTx, goal, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrGoalAlreadyExists) {
			return nil, statusAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create goal",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.CreateGoalResponse{
		Goal: goal.Goal,
	}, nil
}

func (s *experimentService) createGoalNoCommand(
	ctx context.Context,
	req *proto.CreateGoalRequest,
	editor *eventproto.Editor,
) (*proto.CreateGoalResponse, error) {
	if err := validateCreateGoalNoCommandRequest(req); err != nil {
		return nil, err
	}
	goal, err := domain.NewGoal(req.Id, req.Name, req.Description, req.ConnectionType)
	if err != nil {
		s.logger.Error(
			"Failed to create a new goal",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
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
			nil,
		)
		if err != nil {
			return err
		}
		if err := s.publisher.Publish(ctx, e); err != nil {
			return err
		}
		return s.goalStorage.CreateGoal(ctxWithTx, goal, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrGoalAlreadyExists) {
			return nil, statusAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create goal",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.CreateGoalResponse{
		Goal: goal.Goal,
	}, nil
}

func validateCreateGoalRequest(req *proto.CreateGoalRequest) error {
	if req.Command.Id == "" {
		return statusGoalIDRequired.Err()
	}
	if !goalIDRegex.MatchString(req.Command.Id) {
		return statusInvalidGoalID.Err()
	}
	if req.Command.Name == "" {
		return statusGoalNameRequired.Err()
	}
	return nil
}

func validateCreateGoalNoCommandRequest(req *proto.CreateGoalRequest) error {
	if req.Id == "" {
		return statusGoalIDRequired.Err()
	}
	if !goalIDRegex.MatchString(req.Id) {
		return statusInvalidGoalID.Err()
	}
	if req.Name == "" {
		return statusGoalNameRequired.Err()
	}
	return nil
}

func (s *experimentService) UpdateGoal(
	ctx context.Context,
	req *proto.UpdateGoalRequest,
) (*proto.UpdateGoalResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if req.ChangeDescriptionCommand == nil && req.RenameCommand == nil {
		return s.updateGoalNoCommand(ctx, req, editor)
	}
	if req.Id == "" {
		return nil, statusGoalIDRequired.Err()
	}
	commands := make([]command.Command, 0)
	if req.RenameCommand != nil {
		commands = append(commands, req.RenameCommand)
	}
	if req.ChangeDescriptionCommand != nil {
		commands = append(commands, req.ChangeDescriptionCommand)
	}
	if len(commands) == 0 {
		return nil, statusNoCommand.Err()
	}
	err = s.updateGoal(
		ctx,
		editor,
		req.EnvironmentId,
		req.Id,
		commands,
	)
	if err != nil {
		s.logger.Error(
			"Failed to update goal",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
) (*proto.UpdateGoalResponse, error) {
	err := s.validateUpdateGoalNoCommandRequest(req)
	if err != nil {
		return nil, err
	}
	var updatedGoal *proto.Goal
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		goal, err := s.goalStorage.GetGoal(ctxWithTx, req.Id, req.EnvironmentId)
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
		return s.goalStorage.UpdateGoal(ctxWithTx, updated, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrGoalNotFound) || errors.Is(err, v2es.ErrGoalUnexpectedAffectedRows) {
			return nil, statusGoalNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.UpdateGoalResponse{
		Goal: updatedGoal,
	}, nil
}

func (s *experimentService) validateUpdateGoalNoCommandRequest(
	req *proto.UpdateGoalRequest,
) error {
	if req.Id == "" {
		return statusGoalIDRequired.Err()
	}
	if req.Name != nil && req.Name.Value == "" {
		return statusGoalNameRequired.Err()
	}
	return nil
}

func (s *experimentService) ArchiveGoal(
	ctx context.Context,
	req *proto.ArchiveGoalRequest,
) (*proto.ArchiveGoalResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		return nil, statusGoalIDRequired.Err()
	}
	if req.Command == nil {
		return nil, statusNoCommand.Err()
	}
	err = s.updateGoal(
		ctx,
		editor,
		req.EnvironmentId,
		req.Id,
		[]command.Command{req.Command},
	)
	if err != nil {
		s.logger.Error(
			"Failed to archive goal",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		return nil, statusGoalIDRequired.Err()
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		goal, err := s.goalStorage.GetGoal(ctxWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		e, err := domainevent.NewEvent(
			editor,
			eventproto.Event_GOAL,
			goal.Id,
			eventproto.Event_GOAL_DELETED,
			&eventproto.GoalDeletedEvent{
				Id: goal.Id,
			},
			req.EnvironmentId,
			nil,       // Current state: entity no longer exists
			goal.Goal, // Previous state: what was deleted
		)
		if err != nil {
			return err
		}
		if err := s.publisher.Publish(ctxWithTx, e); err != nil {
			return err
		}
		return s.goalStorage.DeleteGoal(ctxWithTx, req.Id, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrGoalNotFound) || errors.Is(err, v2es.ErrGoalUnexpectedAffectedRows) {
			return nil, statusGoalNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.DeleteGoalResponse{}, nil
}

func (s *experimentService) updateGoal(
	ctx context.Context,
	editor *eventproto.Editor,
	environmentId, goalID string,
	commands []command.Command,
) error {
	err := s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		goal, err := s.goalStorage.GetGoal(ctxWithTx, goalID, environmentId)
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
		return s.goalStorage.UpdateGoal(ctxWithTx, goal, environmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrGoalNotFound) || errors.Is(err, v2es.ErrGoalUnexpectedAffectedRows) {
			return statusGoalNotFound.Err()
		}
		return api.NewGRPCStatus(err).Err()
	}
	return nil
}
