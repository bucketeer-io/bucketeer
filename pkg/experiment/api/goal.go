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
	"regexp"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/experiment/command"
	"github.com/bucketeer-io/bucketeer/pkg/experiment/domain"
	v2es "github.com/bucketeer-io/bucketeer/pkg/experiment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

var goalIDRegex = regexp.MustCompile("^[a-zA-Z0-9-]+$")

func (s *experimentService) GetGoal(ctx context.Context, req *proto.GetGoalRequest) (*proto.GetGoalResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
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
	goal, err := s.getGoalMySQL(ctx, req.Id, req.EnvironmentNamespace)
	if err != nil {
		if err == v2es.ErrGoalNotFound {
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
	return &proto.GetGoalResponse{Goal: goal.Goal}, nil
}

func (s *experimentService) getGoalMySQL(
	ctx context.Context,
	goalID, environmentNamespace string,
) (*domain.Goal, error) {
	goalStorage := v2es.NewGoalStorage(s.mysqlClient)
	goal, err := goalStorage.GetGoal(ctx, goalID, environmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to get goal",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
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
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace),
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
		req.EnvironmentNamespace,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list goals",
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

func (s *experimentService) CreateGoal(
	ctx context.Context,
	req *proto.CreateGoalRequest,
) (*proto.CreateGoalResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateCreateGoalRequest(req, localizer); err != nil {
		return nil, err
	}
	goal, err := domain.NewGoal(req.Command.Id, req.Command.Name, req.Command.Description)
	if err != nil {
		s.logger.Error(
			"Failed to create a new goal",
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
		handler := command.NewGoalCommandHandler(editor, goal, s.publisher, req.EnvironmentNamespace)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return goalStorage.CreateGoal(ctx, goal, req.EnvironmentNamespace)
	})
	if err != nil {
		if err == v2es.ErrGoalAlreadyExists {
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
	return &proto.CreateGoalResponse{}, nil
}

func validateCreateGoalRequest(req *proto.CreateGoalRequest, localizer locale.Localizer) error {
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
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

func (s *experimentService) UpdateGoal(
	ctx context.Context,
	req *proto.UpdateGoalRequest,
) (*proto.UpdateGoalResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
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
		req.EnvironmentNamespace,
		req.Id,
		commands,
		localizer,
	)
	if err != nil {
		s.logger.Error(
			"Failed to update goal",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	return &proto.UpdateGoalResponse{}, nil
}

func (s *experimentService) ArchiveGoal(
	ctx context.Context,
	req *proto.ArchiveGoalRequest,
) (*proto.ArchiveGoalResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
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
		req.EnvironmentNamespace,
		req.Id,
		[]command.Command{req.Command},
		localizer,
	)
	if err != nil {
		s.logger.Error(
			"Failed to archive goal",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
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
		req.EnvironmentNamespace,
		req.Id,
		[]command.Command{req.Command},
		localizer,
	)
	if err != nil {
		s.logger.Error(
			"Failed to delete goal",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	return &proto.DeleteGoalResponse{}, nil
}

func (s *experimentService) updateGoal(
	ctx context.Context,
	editor *eventproto.Editor,
	environmentNamespace, goalID string,
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
		goal, err := goalStorage.GetGoal(ctx, goalID, environmentNamespace)
		if err != nil {
			return err
		}
		handler := command.NewGoalCommandHandler(editor, goal, s.publisher, environmentNamespace)
		for _, command := range commands {
			if err := handler.Handle(ctx, command); err != nil {
				return err
			}
		}
		return goalStorage.UpdateGoal(ctx, goal, environmentNamespace)
	})
	if err != nil {
		if err == v2es.ErrGoalNotFound || err == v2es.ErrGoalUnexpectedAffectedRows {
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
