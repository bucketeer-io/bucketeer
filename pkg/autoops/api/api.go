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
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	authclient "github.com/bucketeer-io/bucketeer/v2/pkg/auth/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/command"
	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/storage/v2"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	v2os "github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type AutoOpsService struct {
	mysqlClient      mysql.Client
	opsCountStorage  v2os.OpsCountStorage
	autoOpsStorage   v2as.AutoOpsRuleStorage
	prStorage        v2as.ProgressiveRolloutStorage
	featureStorage   v2fs.FeatureStorage
	featureClient    featureclient.Client
	experimentClient experimentclient.Client
	accountClient    accountclient.Client
	authClient       authclient.Client
	publisher        publisher.Publisher
	opts             *options
	logger           *zap.Logger
}

func NewAutoOpsService(
	mysqlClient mysql.Client,
	featureClient featureclient.Client,
	experimentClient experimentclient.Client,
	accountClient accountclient.Client,
	authClient authclient.Client,
	publisher publisher.Publisher,
	opts ...Option,
) *AutoOpsService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &AutoOpsService{
		mysqlClient:      mysqlClient,
		opsCountStorage:  v2os.NewOpsCountStorage(mysqlClient),
		featureStorage:   v2fs.NewFeatureStorage(mysqlClient),
		autoOpsStorage:   v2as.NewAutoOpsRuleStorage(mysqlClient),
		prStorage:        v2as.NewProgressiveRolloutStorage(mysqlClient),
		featureClient:    featureClient,
		experimentClient: experimentClient,
		accountClient:    accountClient,
		authClient:       authClient,
		publisher:        publisher,
		opts:             dopts,
		logger:           dopts.logger.Named("api"),
	}
}

func (s *AutoOpsService) Register(server *grpc.Server) {
	autoopsproto.RegisterAutoOpsServiceServer(server, s)
}

func (s *AutoOpsService) CreateAutoOpsRule(
	ctx context.Context,
	req *autoopsproto.CreateAutoOpsRuleRequest,
) (*autoopsproto.CreateAutoOpsRuleResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if req.Command == nil {
		return s.createAutoOpsRuleNoCommand(ctx, req, editor)
	}

	if err := s.validateCreateAutoOpsRuleRequest(ctx, req); err != nil {
		return nil, err
	}
	autoOpsRule, err := domain.NewAutoOpsRule(
		req.Command.FeatureId,
		req.Command.OpsType,
		req.Command.OpsEventRateClauses,
		req.Command.DatetimeClauses,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create a new autoOpsRule",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	opsEventRateClauses, err := autoOpsRule.ExtractOpsEventRateClauses()
	if err != nil {
		s.logger.Error(
			"Failed to extract opsEventRateClauses",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	for _, c := range opsEventRateClauses {
		exist, err := s.existGoal(ctx, req.EnvironmentId, c.GoalId)
		if err != nil {
			return nil, api.NewGRPCStatus(err).Err()
		}
		if !exist {
			s.logger.Error(
				"Goal does not exist",
				log.FieldsFromIncomingContext(ctx).AddFields(zap.String("environmentId", req.EnvironmentId))...,
			)
			return nil, statusOpsEventRateClauseGoalNotFound.Err()
		}
	}

	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		handler, err := command.NewAutoOpsCommandHandler(editor, autoOpsRule, s.publisher, req.EnvironmentId)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return s.autoOpsStorage.CreateAutoOpsRule(contextWithTx, autoOpsRule, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAutoOpsRuleAlreadyExists) {
			return nil, statusAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create autoOps",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &autoopsproto.CreateAutoOpsRuleResponse{
		AutoOpsRule: autoOpsRule.AutoOpsRule,
	}, nil
}

func (s *AutoOpsService) createAutoOpsRuleNoCommand(
	ctx context.Context,
	req *autoopsproto.CreateAutoOpsRuleRequest,
	editor *eventproto.Editor,
) (*autoopsproto.CreateAutoOpsRuleResponse, error) {
	if err := s.validateCreateAutoOpsRuleRequestNoCommand(ctx, req); err != nil {
		return nil, err
	}

	autoOpsRule, err := domain.NewAutoOpsRule(
		req.FeatureId,
		req.OpsType,
		req.OpsEventRateClauses,
		req.DatetimeClauses,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create a new autoOpsRule",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	opsEventRateClauses, err := autoOpsRule.ExtractOpsEventRateClauses()
	if err != nil {
		s.logger.Error(
			"Failed to extract opsEventRateClauses",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	for _, c := range opsEventRateClauses {
		exist, err := s.existGoal(ctx, req.EnvironmentId, c.GoalId)
		if err != nil {
			return nil, api.NewGRPCStatus(err).Err()
		}
		if !exist {
			s.logger.Error(
				"Goal does not exist",
				log.FieldsFromIncomingContext(ctx).AddFields(zap.String("environmentId", req.EnvironmentId))...,
			)
			return nil, statusOpsEventRateClauseGoalNotFound.Err()
		}
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		e, err := domainevent.NewEvent(
			editor,
			eventproto.Event_AUTOOPS_RULE,
			autoOpsRule.Id,
			eventproto.Event_AUTOOPS_RULE_CREATED,
			&eventproto.AutoOpsRuleCreatedEvent{
				FeatureId: autoOpsRule.FeatureId,
				OpsType:   autoOpsRule.OpsType,
				Clauses:   autoOpsRule.Clauses,
				CreatedAt: autoOpsRule.CreatedAt,
				UpdatedAt: autoOpsRule.UpdatedAt,
				OpsStatus: autoOpsRule.AutoOpsStatus,
			},
			req.EnvironmentId,
			autoOpsRule.AutoOpsRule,
			nil,
		)
		if err != nil {
			return err
		}
		if err := s.publisher.Publish(ctx, e); err != nil {
			return err
		}
		return s.autoOpsStorage.CreateAutoOpsRule(contextWithTx, autoOpsRule, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAutoOpsRuleAlreadyExists) {
			return nil, statusAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create autoOps",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &autoopsproto.CreateAutoOpsRuleResponse{
		AutoOpsRule: autoOpsRule.AutoOpsRule,
	}, nil
}

func (s *AutoOpsService) validateCreateAutoOpsRuleRequest(
	ctx context.Context,
	req *autoopsproto.CreateAutoOpsRuleRequest,
) error {
	if req.Command.FeatureId == "" {
		return statusFeatureIDRequired.Err()
	}
	if len(req.Command.OpsEventRateClauses) == 0 &&
		len(req.Command.DatetimeClauses) == 0 {
		return statusClauseRequired.Err()
	}
	if req.Command.OpsType == autoopsproto.OpsType_TYPE_UNKNOWN {
		return statusIncompatibleOpsType.Err()
	}
	if req.Command.OpsType == autoopsproto.OpsType_EVENT_RATE {
		if len(req.Command.OpsEventRateClauses) == 0 {
			return statusClauseRequiredForEventRate.Err()
		}
		if len(req.Command.DatetimeClauses) > 0 {
			return statusIncompatibleOpsType.Err()
		}
	}
	if req.Command.OpsType == autoopsproto.OpsType_SCHEDULE {
		if len(req.Command.DatetimeClauses) == 0 {
			return statusClauseRequiredForDateTime.Err()
		}
		if len(req.Command.OpsEventRateClauses) > 0 {
			return statusIncompatibleOpsType.Err()
		}
	}
	if err := s.validateOpsEventRateClauses(req.Command.OpsEventRateClauses); err != nil {
		return err
	}
	if err := s.validateDatetimeClauses(req.Command.DatetimeClauses); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) validateCreateAutoOpsRuleRequestNoCommand(
	ctx context.Context,
	req *autoopsproto.CreateAutoOpsRuleRequest,
) error {
	if req.FeatureId == "" {
		return statusFeatureIDRequired.Err()
	}
	if len(req.OpsEventRateClauses) == 0 &&
		len(req.DatetimeClauses) == 0 {
		return statusClauseRequired.Err()
	}
	if req.OpsType == autoopsproto.OpsType_TYPE_UNKNOWN {
		return statusIncompatibleOpsType.Err()
	}
	if req.OpsType == autoopsproto.OpsType_EVENT_RATE {
		if len(req.OpsEventRateClauses) == 0 {
			return statusClauseRequiredForEventRate.Err()
		}
		if len(req.DatetimeClauses) > 0 {
			return statusIncompatibleOpsType.Err()
		}
	}
	if req.OpsType == autoopsproto.OpsType_SCHEDULE {
		if len(req.DatetimeClauses) == 0 {
			return statusClauseRequiredForDateTime.Err()
		}
		if len(req.OpsEventRateClauses) > 0 {
			return statusIncompatibleOpsType.Err()
		}
	}
	if err := s.validateOpsEventRateClauses(req.OpsEventRateClauses); err != nil {
		return err
	}
	if err := s.validateDatetimeClauses(req.DatetimeClauses); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) validateOpsEventRateClauses(
	clauses []*autoopsproto.OpsEventRateClause,
) error {
	for _, c := range clauses {
		if err := s.validateOpsEventRateClause(c); err != nil {
			return err
		}
	}
	return nil
}

func (s *AutoOpsService) validateOpsEventRateClause(
	clause *autoopsproto.OpsEventRateClause,
) error {
	if clause.VariationId == "" {
		return statusOpsEventRateClauseVariationIDRequired.Err()
	}
	if clause.GoalId == "" {
		return statusOpsEventRateClauseGoalIDRequired.Err()
	}
	if clause.MinCount <= 0 {
		return statusOpsEventRateClauseMinCountRequired.Err()
	}
	if clause.ThreadsholdRate > 1 || clause.ThreadsholdRate <= 0 {
		return statusOpsEventRateClauseInvalidThredshold.Err()
	}
	if clause.ActionType == autoopsproto.ActionType_UNKNOWN || clause.ActionType == autoopsproto.ActionType_ENABLE {
		return statusIncompatibleOpsType.Err()
	}
	return nil
}

func (s *AutoOpsService) validateDatetimeClauses(
	clauses []*autoopsproto.DatetimeClause,
) error {
	checkTimes := make(map[int64]bool)
	for _, c := range clauses {
		if checkTimes[c.Time] {
			return statusDatetimeClauseDuplicateTime.Err()
		}
		if err := s.validateDatetimeClause(c); err != nil {
			return err
		}
		checkTimes[c.Time] = true
	}
	return nil
}

func (s *AutoOpsService) validateDatetimeClause(
	clause *autoopsproto.DatetimeClause,
) error {
	if clause.Time <= time.Now().Unix() {
		return statusDatetimeClauseInvalidTime.Err()
	}
	if clause.ActionType == autoopsproto.ActionType_UNKNOWN {
		return statusIncompatibleOpsType.Err()
	}
	return nil
}

func (s *AutoOpsService) StopAutoOpsRule(
	ctx context.Context,
	req *autoopsproto.StopAutoOpsRuleRequest,
) (*autoopsproto.StopAutoOpsRuleResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if err := validateStopAutoOpsRuleRequest(req); err != nil {
		return nil, err
	}

	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		autoOpsRule, err := s.autoOpsStorage.GetAutoOpsRule(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		if autoOpsRule.IsFinished() {
			return statusAutoOpsRuleFinished.Err()
		}
		stopStatus := autoopsproto.AutoOpsStatus_STOPPED
		updated, err := autoOpsRule.Update(&stopStatus, nil, nil)
		if err != nil {
			return err
		}
		event, err := domainevent.NewEvent(
			editor,
			eventproto.Event_AUTOOPS_RULE,
			autoOpsRule.Id,
			eventproto.Event_AUTOOPS_RULE_OPS_STATUS_CHANGED,
			&eventproto.AutoOpsRuleOpsStatusChangedEvent{
				OpsStatus: stopStatus,
			},
			req.EnvironmentId,
			autoOpsRule,
			autoOpsRule,
		)
		if err != nil {
			return err
		}
		if err := s.publisher.Publish(ctx, event); err != nil {
			return err
		}
		return s.autoOpsStorage.UpdateAutoOpsRule(contextWithTx, updated, req.EnvironmentId)
	})

	if err != nil {
		return nil, err
	}
	return &autoopsproto.StopAutoOpsRuleResponse{}, nil
}

func validateStopAutoOpsRuleRequest(req *autoopsproto.StopAutoOpsRuleRequest) error {
	if req.Id == "" {
		return statusAutoOpsRuleIDRequired.Err()
	}
	return nil
}

func (s *AutoOpsService) DeleteAutoOpsRule(
	ctx context.Context,
	req *autoopsproto.DeleteAutoOpsRuleRequest,
) (*autoopsproto.DeleteAutoOpsRuleResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteAutoOpsRuleRequest(req); err != nil {
		return nil, err
	}

	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		autoOpsRule, err := s.autoOpsStorage.GetAutoOpsRule(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		autoOpsRule.SetDeleted()
		e, err := domainevent.NewEvent(
			editor,
			eventproto.Event_AUTOOPS_RULE,
			autoOpsRule.Id,
			eventproto.Event_AUTOOPS_RULE_DELETED,
			&eventproto.AutoOpsRuleDeletedEvent{},
			req.EnvironmentId,
			nil,                     // Current state: entity no longer exists
			autoOpsRule.AutoOpsRule, // Previous state: what was deleted
		)
		if err != nil {
			return err
		}
		if err := s.publisher.Publish(ctx, e); err != nil {
			return err
		}
		return s.autoOpsStorage.UpdateAutoOpsRule(contextWithTx, autoOpsRule, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAutoOpsRuleNotFound) || errors.Is(err, v2as.ErrAutoOpsRuleUnexpectedAffectedRows) {
			return nil, statusAutoOpsRuleNotFound.Err()
		}
		s.logger.Error(
			"Failed to delete autoOpsRule",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)

	}
	return &autoopsproto.DeleteAutoOpsRuleResponse{}, nil
}

func validateDeleteAutoOpsRuleRequest(req *autoopsproto.DeleteAutoOpsRuleRequest) error {
	if req.Id == "" {
		return statusAutoOpsRuleIDRequired.Err()
	}
	return nil
}

func (s *AutoOpsService) UpdateAutoOpsRule(
	ctx context.Context,
	req *autoopsproto.UpdateAutoOpsRuleRequest,
) (*autoopsproto.UpdateAutoOpsRuleResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if s.isNoUpdateAutoOpsRuleCommand(req) {
		return s.updateAutoOpsRuleNoCommand(ctx, req, editor)
	}

	if err := s.validateUpdateAutoOpsRuleRequest(req); err != nil {
		return nil, err
	}
	var opsEventRateClauses []*autoopsproto.OpsEventRateClause
	for _, c := range req.AddOpsEventRateClauseCommands {
		opsEventRateClauses = append(opsEventRateClauses, c.OpsEventRateClause)
	}
	for _, c := range req.ChangeOpsEventRateClauseCommands {
		opsEventRateClauses = append(opsEventRateClauses, c.OpsEventRateClause)
	}
	for _, c := range opsEventRateClauses {
		exist, err := s.existGoal(ctx, req.EnvironmentId, c.GoalId)
		if err != nil {
			return nil, api.NewGRPCStatus(err).Err()
		}
		if !exist {
			s.logger.Error(
				"Goal does not exist",
				log.FieldsFromIncomingContext(ctx).AddFields(zap.String("environmentId", req.EnvironmentId))...,
			)
			return nil, statusOpsEventRateClauseGoalNotFound.Err()
		}
	}
	commands := s.createUpdateAutoOpsRuleCommands(req)

	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		autoOpsRule, err := s.autoOpsStorage.GetAutoOpsRule(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}

		if autoOpsRule.IsFinished() || autoOpsRule.IsStopped() {
			return statusAutoOpsRuleCompleted.Err()
		}
		if autoOpsRule.OpsType == autoopsproto.OpsType_SCHEDULE {
			if len(req.AddOpsEventRateClauseCommands) > 0 || len(req.ChangeOpsEventRateClauseCommands) > 0 {
				return statusIncompatibleOpsType.Err()
			}

			// Delete a deletion schedule from the currently held schedules
			extractDateTimeClauses, _ := autoOpsRule.ExtractDatetimeClauses()
			for _, deleteClause := range req.DeleteClauseCommands {
				delete(extractDateTimeClauses, deleteClause.Id)
			}
			checkTimes := make(map[int64]autoopsproto.ActionType)
			for _, c := range extractDateTimeClauses {
				checkTimes[c.Time] = c.ActionType
			}

			// Check if there is a schedule with the same date and time.
			for _, c := range req.AddDatetimeClauseCommands {
				actionType, hasSameTime := checkTimes[c.DatetimeClause.Time]
				if hasSameTime && actionType == c.DatetimeClause.ActionType {
					return statusDatetimeClauseDuplicateTime.Err()
				}
			}
			for _, c := range req.ChangeDatetimeClauseCommands {
				actionType, hasSameTime := checkTimes[c.DatetimeClause.Time]
				if hasSameTime && actionType == c.DatetimeClause.ActionType {
					return statusDatetimeClauseDuplicateTime.Err()
				}
			}
		}
		if autoOpsRule.OpsType == autoopsproto.OpsType_EVENT_RATE {
			if len(req.AddDatetimeClauseCommands) > 0 || len(req.ChangeDatetimeClauseCommands) > 0 {
				return statusIncompatibleOpsType.Err()
			}
		}

		if req.DeleteClauseCommands != nil && len(autoOpsRule.Clauses) == len(req.DeleteClauseCommands) &&
			len(req.AddOpsEventRateClauseCommands) == 0 && len(req.AddDatetimeClauseCommands) == 0 {
			// When deleting, at least one Clause must exist.
			return statusShouldAddMoreClauses.Err()
		}
		handler, err := command.NewAutoOpsCommandHandler(editor, autoOpsRule, s.publisher, req.EnvironmentId)
		if err != nil {
			return err
		}
		for _, com := range commands {
			if err := handler.Handle(ctx, com); err != nil {
				return err
			}
		}
		return s.autoOpsStorage.UpdateAutoOpsRule(contextWithTx, autoOpsRule, req.EnvironmentId)
	})
	if err != nil {
		return nil, s.returnUpdateAutoOpsRuleError(ctx, req, err)
	}
	return &autoopsproto.UpdateAutoOpsRuleResponse{}, nil
}

func (s *AutoOpsService) updateAutoOpsRuleNoCommand(
	ctx context.Context,
	req *autoopsproto.UpdateAutoOpsRuleRequest,
	editor *eventproto.Editor,
) (*autoopsproto.UpdateAutoOpsRuleResponse, error) {
	err := s.validateUpdateAutoOpsRuleRequestNoCommand(req)
	if err != nil {
		return nil, err
	}
	for _, c := range req.OpsEventRateClauseChanges {
		if c.ChangeType == autoopsproto.ChangeType_DELETE {
			continue
		}
		goal, err := s.getGoal(ctx, req.EnvironmentId, c.Clause.GoalId)
		if err != nil {
			return nil, err
		}
		if goal == nil || goal.ConnectionType != experimentproto.Goal_OPERATION {
			s.logger.Error(
				"Goal does not exist",
				log.FieldsFromIncomingContext(ctx).AddFields(zap.String("environmentId", req.EnvironmentId))...,
			)
			return nil, statusOpsEventRateClauseGoalNotFound.Err()
		}
	}
	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		autoOpsRule, err := s.autoOpsStorage.GetAutoOpsRule(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}

		if autoOpsRule.IsFinished() || autoOpsRule.IsStopped() {
			return statusAutoOpsRuleCompleted.Err()
		}
		if autoOpsRule.OpsType == autoopsproto.OpsType_SCHEDULE {
			if len(req.OpsEventRateClauseChanges) > 0 {
				return statusIncompatibleOpsType.Err()
			}
			// Delete a deletion schedule from the currently held schedules
			extractDateTimeClauses, _ := autoOpsRule.ExtractDatetimeClauses()
			for _, deleteClause := range req.DatetimeClauseChanges {
				if deleteClause.ChangeType == autoopsproto.ChangeType_DELETE {
					delete(extractDateTimeClauses, deleteClause.Id)
				}
			}
			checkTimes := make(map[int64]autoopsproto.ActionType)
			for _, c := range extractDateTimeClauses {
				checkTimes[c.Time] = c.ActionType
			}

			// Check if there is a schedule with the same date and time.
			for _, c := range req.DatetimeClauseChanges {
				if c.Clause != nil && c.ChangeType != autoopsproto.ChangeType_DELETE {
					actionType, hasSameTime := checkTimes[c.Clause.Time]
					if hasSameTime && actionType == c.Clause.ActionType {
						return statusDatetimeClauseDuplicateTime.Err()
					}
				}
			}
		}

		if autoOpsRule.OpsType == autoopsproto.OpsType_EVENT_RATE {
			if len(req.DatetimeClauseChanges) > 0 {
				return statusIncompatibleOpsType.Err()
			}
		}

		updated, err := autoOpsRule.Update(nil, req.OpsEventRateClauseChanges, req.DatetimeClauseChanges)
		if err != nil {
			return err
		}

		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_AUTOOPS_RULE,
			updated.Id,
			eventproto.Event_AUTOOPS_RULE_CREATED,
			&eventproto.AutoOpsRuleUpdatedEvent{
				Id:                        req.Id,
				OpsEventRateClauseChanges: req.OpsEventRateClauseChanges,
				DatetimeClauseChanges:     req.DatetimeClauseChanges,
			},
			req.EnvironmentId,
			updated,
			autoOpsRule,
		)
		if err != nil {
			return err
		}
		if err := s.publisher.Publish(ctx, event); err != nil {
			return err
		}
		return s.autoOpsStorage.UpdateAutoOpsRule(contextWithTx, updated, req.EnvironmentId)
	})
	if err != nil {
		return nil, s.returnUpdateAutoOpsRuleError(ctx, req, err)
	}
	err = s.publisher.Publish(ctx, event)
	if err != nil {
		s.logger.Error(
			"Failed to publish event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &autoopsproto.UpdateAutoOpsRuleResponse{}, nil
}

func (s *AutoOpsService) returnUpdateAutoOpsRuleError(
	ctx context.Context,
	req *autoopsproto.UpdateAutoOpsRuleRequest,
	err error,
) error {
	if errors.Is(err, v2as.ErrAutoOpsRuleNotFound) || errors.Is(err, v2as.ErrAutoOpsRuleUnexpectedAffectedRows) {
		return statusAutoOpsRuleNotFound.Err()
	}
	if status.Code(err) == codes.InvalidArgument {
		return err
	}
	s.logger.Error(
		"Failed to update autoOpsRule",
		log.FieldsFromIncomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentId", req.EnvironmentId),
		)...,
	)
	return api.NewGRPCStatus(err).Err()
}

func (s *AutoOpsService) validateUpdateAutoOpsRuleRequest(
	req *autoopsproto.UpdateAutoOpsRuleRequest,
) error {
	if req.Id == "" {
		return statusAutoOpsRuleIDRequired.Err()
	}
	for _, c := range req.AddOpsEventRateClauseCommands {
		if c.OpsEventRateClause == nil {
			return statusOpsEventRateClauseRequired.Err()
		}
		if err := s.validateOpsEventRateClause(c.OpsEventRateClause); err != nil {
			return err
		}
	}
	for _, c := range req.ChangeOpsEventRateClauseCommands {
		if c.Id == "" {
			return statusClauseIDRequired.Err()
		}
		if c.OpsEventRateClause == nil {
			return statusOpsEventRateClauseRequired.Err()
		}
		if err := s.validateOpsEventRateClause(c.OpsEventRateClause); err != nil {
			return err
		}
	}
	for _, c := range req.DeleteClauseCommands {
		if c.Id == "" {
			return statusClauseIDRequired.Err()
		}
	}

	var checkDatetimeClauses []*autoopsproto.DatetimeClause
	for _, c := range req.AddDatetimeClauseCommands {
		if c.DatetimeClause == nil {
			return statusDatetimeClauseRequired.Err()
		}
		checkDatetimeClauses = append(checkDatetimeClauses, c.DatetimeClause)
	}
	if err := s.validateDatetimeClauses(checkDatetimeClauses); err != nil {
		return err
	}

	for _, c := range req.ChangeDatetimeClauseCommands {
		if c.Id == "" {
			return statusClauseIDRequired.Err()
		}
		if c.DatetimeClause == nil {
			return statusDatetimeClauseRequired.Err()
		}
		checkDatetimeClauses = append(checkDatetimeClauses, c.DatetimeClause)
	}
	if err := s.validateDatetimeClauses(checkDatetimeClauses); err != nil {
		return err
	}

	return nil
}

func (s *AutoOpsService) validateUpdateAutoOpsRuleRequestNoCommand(
	req *autoopsproto.UpdateAutoOpsRuleRequest,
) error {
	if req.Id == "" {
		return statusAutoOpsRuleIDRequired.Err()
	}
	for _, c := range req.OpsEventRateClauseChanges {
		if c.Id == "" && c.ChangeType == autoopsproto.ChangeType_DELETE {
			return statusClauseIDRequired.Err()
		}
		if c.ChangeType != autoopsproto.ChangeType_DELETE && c.Clause == nil {
			return statusOpsEventRateClauseRequired.Err()
		}
		if err := s.validateOpsEventRateClause(c.Clause); err != nil {
			return err
		}
	}

	var checkDatetimeClauses []*autoopsproto.DatetimeClause
	for _, c := range req.DatetimeClauseChanges {
		if c.Id == "" && c.ChangeType == autoopsproto.ChangeType_DELETE {
			return statusClauseIDRequired.Err()
		}
		if c.ChangeType != autoopsproto.ChangeType_DELETE && c.Clause == nil {
			return statusDatetimeClauseRequired.Err()
		}
		if c.ChangeType != autoopsproto.ChangeType_DELETE {
			checkDatetimeClauses = append(checkDatetimeClauses, c.Clause)
		}
	}
	if err := s.validateDatetimeClauses(checkDatetimeClauses); err != nil {
		return err
	}

	return nil
}

func (s *AutoOpsService) isNoUpdateAutoOpsRuleCommand(req *autoopsproto.UpdateAutoOpsRuleRequest) bool {
	return len(req.AddOpsEventRateClauseCommands) == 0 &&
		len(req.ChangeOpsEventRateClauseCommands) == 0 &&
		len(req.DeleteClauseCommands) == 0 &&
		len(req.AddDatetimeClauseCommands) == 0 &&
		len(req.ChangeDatetimeClauseCommands) == 0
}

func (s *AutoOpsService) createUpdateAutoOpsRuleCommands(req *autoopsproto.UpdateAutoOpsRuleRequest) []command.Command {
	commands := make([]command.Command, 0)
	for _, c := range req.AddOpsEventRateClauseCommands {
		commands = append(commands, c)
	}
	for _, c := range req.ChangeOpsEventRateClauseCommands {
		commands = append(commands, c)
	}
	for _, c := range req.AddDatetimeClauseCommands {
		commands = append(commands, c)
	}
	for _, c := range req.ChangeDatetimeClauseCommands {
		commands = append(commands, c)
	}
	for _, c := range req.DeleteClauseCommands {
		commands = append(commands, c)
	}
	return commands
}

func (s *AutoOpsService) GetAutoOpsRule(
	ctx context.Context,
	req *autoopsproto.GetAutoOpsRuleRequest,
) (*autoopsproto.GetAutoOpsRuleResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := s.validateGetAutoOpsRuleRequest(req); err != nil {
		return nil, err
	}
	autoOpsRule, err := s.autoOpsStorage.GetAutoOpsRule(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2as.ErrAutoOpsRuleNotFound) {
			return nil, statusAutoOpsRuleNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	if autoOpsRule.Deleted {
		return nil, statusAutoOpsRuleAlreadyDeleted.Err()
	}
	return &autoopsproto.GetAutoOpsRuleResponse{
		AutoOpsRule: autoOpsRule.AutoOpsRule,
	}, nil
}

func (s *AutoOpsService) validateGetAutoOpsRuleRequest(
	req *autoopsproto.GetAutoOpsRuleRequest,
) error {
	if req.Id == "" {
		return statusAutoOpsRuleIDRequired.Err()
	}
	return nil
}

func (s *AutoOpsService) ListAutoOpsRules(
	ctx context.Context,
	req *autoopsproto.ListAutoOpsRulesRequest,
) (*autoopsproto.ListAutoOpsRulesResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	autoOpsRules, cursor, err := s.listAutoOpsRules(
		ctx,
		req.PageSize,
		req.Cursor,
		req.FeatureIds,
		req.EnvironmentId,
	)
	if err != nil {
		return nil, err
	}
	return &autoopsproto.ListAutoOpsRulesResponse{
		AutoOpsRules: autoOpsRules,
		Cursor:       cursor,
	}, nil
}

func (s *AutoOpsService) listAutoOpsRules(
	ctx context.Context,
	pageSize int64,
	cursor string,
	featureIds []string,
	environmentId string,
) ([]*autoopsproto.AutoOpsRule, string, error) {
	filters := []*mysql.FilterV2{
		{
			Column:   "aor.deleted",
			Operator: mysql.OperatorEqual,
			Value:    false,
		},
		{
			Column:   "aor.environment_id",
			Operator: mysql.OperatorEqual,
			Value:    environmentId,
		},
	}
	fIDs := make([]interface{}, 0, len(featureIds))
	for _, fID := range featureIds {
		fIDs = append(fIDs, fID)
	}
	var inFilters []*mysql.InFilter
	if len(fIDs) > 0 {
		inFilters = append(inFilters, &mysql.InFilter{
			Column: "aor.feature_id",
			Values: fIDs,
		})
	}
	limit := int(pageSize)
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, "", statusInvalidCursor.Err()
	}
	options := &mysql.ListOptions{
		Limit:       limit,
		Offset:      offset,
		Filters:     filters,
		InFilters:   inFilters,
		NullFilters: nil,
		JSONFilters: nil,
		SearchQuery: nil,
		Orders:      nil,
	}
	autoOpsRules, nextCursor, err := s.autoOpsStorage.ListAutoOpsRules(ctx, options)
	if err != nil {
		s.logger.Error(
			"Failed to list autoOpsRules",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		return nil, "", api.NewGRPCStatus(err).Err()
	}
	return autoOpsRules, strconv.Itoa(nextCursor), nil
}

func (s *AutoOpsService) ExecuteAutoOps(
	ctx context.Context,
	req *autoopsproto.ExecuteAutoOpsRequest,
) (*autoopsproto.ExecuteAutoOpsResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if req.ExecuteAutoOpsRuleCommand == nil {
		return s.executeAutoOpsNoCommand(ctx, req, editor)
	}

	if err := s.validateExecuteAutoOpsRequest(req); err != nil {
		return nil, err
	}
	triggered, err := s.checkIfHasAlreadyTriggered(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if triggered {
		return &autoopsproto.ExecuteAutoOpsResponse{AlreadyTriggered: true}, nil
	}

	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, tx mysql.Transaction) error {
		autoOpsRule, err := s.autoOpsStorage.GetAutoOpsRule(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}

		var executeClause *autoopsproto.Clause
		for _, c := range autoOpsRule.Clauses {
			if c.Id == req.ExecuteAutoOpsRuleCommand.ClauseId {
				executeClause = c
				break
			}
		}
		// Check if the clause exists
		if executeClause == nil {
			return statusClauseNotFound.Err()
		}
		// Check if the clause is already executed
		if executeClause.ExecutedAt != 0 {
			return statusClauseAlreadyExecuted.Err()
		}
		ftStorage := v2fs.NewFeatureStorage(tx)
		feature, err := ftStorage.GetFeature(contextWithTx, autoOpsRule.FeatureId, req.EnvironmentId)
		if err != nil {
			return err
		}
		// Stop the running progressive rollout if the operation type is disable
		if executeClause.ActionType == autoopsproto.ActionType_DISABLE {
			if err := s.stopProgressiveRollout(
				contextWithTx,
				req.EnvironmentId,
				autoOpsRule,
			); err != nil {
				return err
			}
		}
		if err := executeAutoOpsRuleOperation(
			contextWithTx,
			ftStorage,
			req.EnvironmentId,
			executeClause.ActionType,
			feature,
			s.logger,
			s.publisher,
			editor,
		); err != nil {
			s.logger.Error(
				"Failed to execute auto ops rule operation",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
					zap.String("autoOpsRuleId", autoOpsRule.Id),
					zap.String("featureId", autoOpsRule.FeatureId),
				)...,
			)
			return err
		}
		// Set the `executed_at`, so it won't be executed twice
		executeClause.ExecutedAt = time.Now().Unix()
		// Update the status if needed.
		// When it executes the last clause, it will change to finished status.
		opsStatus := autoopsproto.AutoOpsStatus_RUNNING
		if autoOpsRule.Clauses[len(autoOpsRule.Clauses)-1].Id == req.ExecuteAutoOpsRuleCommand.ClauseId {
			opsStatus = autoopsproto.AutoOpsStatus_FINISHED
		}
		handler, err := command.NewAutoOpsCommandHandler(editor, autoOpsRule, s.publisher, req.EnvironmentId)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, &autoopsproto.ChangeAutoOpsStatusCommand{Status: opsStatus}); err != nil {
			return err
		}

		if err = s.autoOpsStorage.UpdateAutoOpsRule(contextWithTx, autoOpsRule, req.EnvironmentId); err != nil {
			if errors.Is(err, v2as.ErrAutoOpsRuleUnexpectedAffectedRows) {
				s.logger.Warn(
					"No rows were affected",
					log.FieldsFromIncomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("id", req.Id),
						zap.String("environmentId", req.EnvironmentId),
					)...,
				)
				return nil
			}
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAutoOpsRuleNotFound) {
			s.logger.Warn(
				"Auto Ops Rule not found",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", req.Id),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return nil, statusAutoOpsRuleNotFound.Err()
		}
		s.logger.Error(
			"Failed to execute autoOpsRule",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &autoopsproto.ExecuteAutoOpsResponse{AlreadyTriggered: false}, nil
}

func (s *AutoOpsService) executeAutoOpsNoCommand(
	ctx context.Context,
	req *autoopsproto.ExecuteAutoOpsRequest,
	editor *eventproto.Editor,
) (*autoopsproto.ExecuteAutoOpsResponse, error) {
	if err := s.validateExecuteAutoOpsRequestNoCommand(req); err != nil {
		return nil, err
	}
	triggered, err := s.checkIfHasAlreadyTriggered(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if triggered {
		return &autoopsproto.ExecuteAutoOpsResponse{AlreadyTriggered: true}, nil
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, tx mysql.Transaction) error {
		autoOpsRule, err := s.autoOpsStorage.GetAutoOpsRule(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}

		var executeClause *autoopsproto.Clause = nil
		for _, c := range autoOpsRule.Clauses {
			if c.Id == req.ClauseId {
				executeClause = c
				break
			}
		}
		// Check if the clause exists
		if executeClause == nil {
			return statusClauseNotFound.Err()
		}
		// Check if the clause is already executed
		if executeClause.ExecutedAt != 0 {
			return statusClauseAlreadyExecuted.Err()
		}

		ftStorage := v2fs.NewFeatureStorage(tx)
		feature, err := ftStorage.GetFeature(contextWithTx, autoOpsRule.FeatureId, req.EnvironmentId)
		if err != nil {
			return err
		}
		// Stop the running progressive rollout if the operation type is disable
		if executeClause.ActionType == autoopsproto.ActionType_DISABLE {
			if err := s.stopProgressiveRollout(
				contextWithTx,
				req.EnvironmentId,
				autoOpsRule,
			); err != nil {
				return err
			}
		}
		if err := executeAutoOpsRuleOperation(
			contextWithTx,
			ftStorage,
			req.EnvironmentId,
			executeClause.ActionType,
			feature,
			s.logger,
			s.publisher,
			editor,
		); err != nil {
			s.logger.Error(
				"Failed to execute auto ops rule operation",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
					zap.String("autoOpsRuleId", autoOpsRule.Id),
					zap.String("featureId", autoOpsRule.FeatureId),
				)...,
			)
			return err
		}
		// Set the `executed_at`, so it won't be executd twice
		executeClause.ExecutedAt = time.Now().Unix()
		// Update the status if needed.
		// When it executes the last clause, it will change to finished status.
		opsStatus := autoopsproto.AutoOpsStatus_RUNNING
		if autoOpsRule.Clauses[len(autoOpsRule.Clauses)-1].Id == req.ClauseId {
			opsStatus = autoopsproto.AutoOpsStatus_FINISHED
		}
		updated, err := autoOpsRule.Update(&opsStatus, nil, nil)
		if err != nil {
			return err
		}
		event, err := domainevent.NewEvent(
			editor,
			eventproto.Event_AUTOOPS_RULE,
			autoOpsRule.Id,
			eventproto.Event_AUTOOPS_RULE_OPS_STATUS_CHANGED,
			&eventproto.AutoOpsRuleOpsStatusChangedEvent{
				OpsStatus: opsStatus,
			},
			req.EnvironmentId,
			autoOpsRule,
			autoOpsRule,
		)
		if err != nil {
			return err
		}
		if err := s.publisher.Publish(ctx, event); err != nil {
			return err
		}
		if err = s.autoOpsStorage.UpdateAutoOpsRule(contextWithTx, updated, req.EnvironmentId); err != nil {
			if errors.Is(err, v2as.ErrAutoOpsRuleUnexpectedAffectedRows) {
				s.logger.Warn(
					"No rows were affected",
					log.FieldsFromIncomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("id", req.Id),
						zap.String("environmentId", req.EnvironmentId),
					)...,
				)
				return nil
			}
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAutoOpsRuleNotFound) {
			s.logger.Warn(
				"Auto Ops Rule not found",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", req.Id),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return nil, statusAutoOpsRuleNotFound.Err()
		}
		s.logger.Error(
			"Failed to execute autoOpsRule",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &autoopsproto.ExecuteAutoOpsResponse{AlreadyTriggered: false}, nil
}

func (s *AutoOpsService) stopProgressiveRollout(
	ctx context.Context,
	environmentId string,
	autoOpsRule *domain.AutoOpsRule,
) error {
	// Check what operation is being executed
	// and the set progressive rollout `stoppedBy`
	var stoppedBy autoopsproto.ProgressiveRollout_StoppedBy
	hasScheduleOps, err := autoOpsRule.HasScheduleOps()
	if err != nil {
		s.logger.Error(
			"Failed to check operation type",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
				zap.String("autoOpsRuleId", autoOpsRule.Id),
				zap.String("featureId", autoOpsRule.FeatureId),
			)...,
		)
		return api.NewGRPCStatus(err).Err()
	}
	if hasScheduleOps {
		stoppedBy = autoopsproto.ProgressiveRollout_OPS_SCHEDULE
	} else {
		stoppedBy = autoopsproto.ProgressiveRollout_OPS_KILL_SWITCH
	}
	if err := executeStopProgressiveRolloutOperation(
		ctx,
		s.prStorage,
		s.convToInterfaceSlice([]string{autoOpsRule.FeatureId}),
		environmentId,
		stoppedBy,
	); err != nil {
		s.logger.Error(
			"Failed to stop progressive rollout",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
				zap.String("autoOpsRuleId", autoOpsRule.Id),
				zap.String("featureId", autoOpsRule.FeatureId),
			)...,
		)
		return api.NewGRPCStatus(err).Err()
	}
	return nil
}

func (s *AutoOpsService) validateExecuteAutoOpsRequest(
	req *autoopsproto.ExecuteAutoOpsRequest,
) error {
	if req.Id == "" {
		return statusAutoOpsRuleIDRequired.Err()
	}
	if req.ExecuteAutoOpsRuleCommand != nil && req.ExecuteAutoOpsRuleCommand.ClauseId == "" {
		return statusClauseIDRequired.Err()
	}
	return nil
}

func (s *AutoOpsService) validateExecuteAutoOpsRequestNoCommand(
	req *autoopsproto.ExecuteAutoOpsRequest,
) error {
	if req.Id == "" {
		return statusAutoOpsRuleIDRequired.Err()
	}
	if req.ClauseId == "" {
		return statusClauseIDRequired.Err()
	}
	return nil
}

func (s *AutoOpsService) checkIfHasAlreadyTriggered(
	ctx context.Context,
	ruleID,
	environmentId string,
) (bool, error) {
	autoOpsRule, err := s.autoOpsStorage.GetAutoOpsRule(ctx, ruleID, environmentId)
	if err != nil {
		if errors.Is(err, v2as.ErrAutoOpsRuleNotFound) {
			s.logger.Warn(
				"Auto Ops Rule not found",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("ruleID", ruleID),
					zap.String("environmentId", environmentId),
				)...,
			)
			return false, statusAutoOpsRuleNotFound.Err()
		}
		s.logger.Error(
			"Failed to get auto ops rule",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		return false, api.NewGRPCStatus(err).Err()
	}
	if autoOpsRule.IsFinished() || autoOpsRule.IsStopped() || autoOpsRule.Deleted {
		s.logger.Warn(
			"Auto Ops Rule already triggered",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("ruleID", ruleID),
				zap.String("environmentId", environmentId),
			)...,
		)
		return true, nil
	}
	return false, nil
}

func (s *AutoOpsService) ListOpsCounts(
	ctx context.Context,
	req *autoopsproto.ListOpsCountsRequest,
) (*autoopsproto.ListOpsCountsResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	opsCounts, cursor, err := s.listOpsCounts(
		ctx,
		req.PageSize,
		req.Cursor,
		req.EnvironmentId,
		req.FeatureIds,
		req.AutoOpsRuleIds,
	)
	if err != nil {
		return nil, err
	}
	return &autoopsproto.ListOpsCountsResponse{
		Cursor:    cursor,
		OpsCounts: opsCounts,
	}, nil
}

func (s *AutoOpsService) listOpsCounts(
	ctx context.Context,
	pageSize int64,
	cursor string,
	environmentId string,
	featureIDs []string,
	autoOpsRuleIDs []string,
) ([]*autoopsproto.OpsCount, string, error) {
	var infilters []*mysql.InFilter
	fIDs := make([]interface{}, 0, len(featureIDs))
	for _, fID := range featureIDs {
		fIDs = append(fIDs, fID)
	}
	if len(fIDs) > 0 {
		infilters = append(infilters, &mysql.InFilter{
			Column: "feature_id",
			Values: fIDs,
		})
	}
	aIDs := make([]interface{}, 0, len(autoOpsRuleIDs))
	for _, aID := range autoOpsRuleIDs {
		aIDs = append(aIDs, aID)
	}
	if len(aIDs) > 0 {
		infilters = append(infilters, &mysql.InFilter{
			Column: "auto_ops_rule_id",
			Values: aIDs,
		})
	}
	limit := int(pageSize)
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, "", statusInvalidCursor.Err()
	}
	options := &mysql.ListOptions{
		Limit:  limit,
		Offset: offset,
		Filters: []*mysql.FilterV2{
			{
				Column:   "environment_id",
				Operator: mysql.OperatorEqual,
				Value:    environmentId,
			},
		},
		InFilters:   infilters,
		NullFilters: nil,
		JSONFilters: nil,
		SearchQuery: nil,
		Orders:      nil,
	}
	opsCounts, nextCursor, err := s.opsCountStorage.ListOpsCounts(
		ctx,
		options,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list opsCounts",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		return nil, "", api.NewGRPCStatus(err).Err()
	}
	return opsCounts, strconv.Itoa(nextCursor), nil
}

func (s *AutoOpsService) existGoal(ctx context.Context, environmentId string, goalID string) (bool, error) {
	_, err := s.getGoal(ctx, environmentId, goalID)
	if err != nil {
		if errors.Is(err, storage.ErrKeyNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *AutoOpsService) getGoal(
	ctx context.Context,
	environmentId, goalID string,
) (*experimentproto.Goal, error) {
	resp, err := s.experimentClient.GetGoal(ctx, &experimentproto.GetGoalRequest{
		Id:            goalID,
		EnvironmentId: environmentId,
	})
	if err != nil {
		s.logger.Error(
			"Failed to get goal",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
				zap.String("goalId", goalID),
			)...,
		)
		return nil, err
	}
	return resp.Goal, nil
}

func (s *AutoOpsService) checkEnvironmentRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Environment,
	environmentId string,
) (*eventproto.Editor, error) {
	editor, err := role.CheckEnvironmentRole(
		ctx,
		requiredRole,
		environmentId,
		func(email string) (*accountproto.AccountV2, error) {
			resp, err := s.accountClient.GetAccountV2ByEnvironmentID(ctx, &accountproto.GetAccountV2ByEnvironmentIDRequest{
				Email:         email,
				EnvironmentId: environmentId,
			})
			if err != nil {
				return nil, err
			}
			return resp.Account, nil
		})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return nil, statusUnauthenticated.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return nil, statusPermissionDenied.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return nil, api.NewGRPCStatus(err).Err()
		}
	}
	return editor, nil
}
