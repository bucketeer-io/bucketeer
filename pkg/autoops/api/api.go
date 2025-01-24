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
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	authclient "github.com/bucketeer-io/bucketeer/pkg/auth/client"
	"github.com/bucketeer-io/bucketeer/pkg/autoops/command"
	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	v2as "github.com/bucketeer-io/bucketeer/pkg/autoops/storage/v2"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	ftstorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	v2os "github.com/bucketeer-io/bucketeer/pkg/opsevent/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/role"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
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
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateCreateAutoOpsRuleRequest(ctx, req, localizer); err != nil {
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
	opsEventRateClauses, err := autoOpsRule.ExtractOpsEventRateClauses()
	if err != nil {
		s.logger.Error(
			"Failed to extract opsEventRateClauses",
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
	for _, c := range opsEventRateClauses {
		exist, err := s.existGoal(ctx, req.EnvironmentId, c.GoalId)
		if err != nil {
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		if !exist {
			s.logger.Error(
				"Goal does not exist",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.String("environmentId", req.EnvironmentId))...,
			)
			dt, err := statusOpsEventRateClauseGoalNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
	}

	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(s.mysqlClient)
		handler, err := command.NewAutoOpsCommandHandler(editor, autoOpsRule, s.publisher, req.EnvironmentId)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return autoOpsRuleStorage.CreateAutoOpsRule(contextWithTx, autoOpsRule, req.EnvironmentId)
	})
	if err != nil {
		if err == v2as.ErrAutoOpsRuleAlreadyExists {
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
			"Failed to create autoOps",
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
	return &autoopsproto.CreateAutoOpsRuleResponse{}, nil
}

func (s *AutoOpsService) validateCreateAutoOpsRuleRequest(
	ctx context.Context,
	req *autoopsproto.CreateAutoOpsRuleRequest,
	localizer locale.Localizer,
) error {
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
	if req.Command.FeatureId == "" {
		dt, err := statusFeatureIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if len(req.Command.OpsEventRateClauses) == 0 &&
		len(req.Command.DatetimeClauses) == 0 {
		dt, err := statusClauseRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.OpsType == autoopsproto.OpsType_TYPE_UNKNOWN {
		dt, err := statusIncompatibleOpsType.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "ops_type"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.OpsType == autoopsproto.OpsType_EVENT_RATE {
		if len(req.Command.OpsEventRateClauses) == 0 {
			dt, err := statusClauseRequiredForEventDate.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "clause"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		if len(req.Command.DatetimeClauses) > 0 {
			dt, err := statusIncompatibleOpsType.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "ops_type"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
	}
	if req.Command.OpsType == autoopsproto.OpsType_SCHEDULE {
		if len(req.Command.DatetimeClauses) == 0 {
			dt, err := statusClauseRequiredForDateTime.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "clause"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		if len(req.Command.OpsEventRateClauses) > 0 {
			dt, err := statusIncompatibleOpsType.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "ops_type"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
	}
	if err := s.validateOpsEventRateClauses(req.Command.OpsEventRateClauses, localizer); err != nil {
		return err
	}
	if err := s.validateDatetimeClauses(req.Command.DatetimeClauses, localizer); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) validateOpsEventRateClauses(
	clauses []*autoopsproto.OpsEventRateClause,
	localizer locale.Localizer,
) error {
	for _, c := range clauses {
		if err := s.validateOpsEventRateClause(c, localizer); err != nil {
			return err
		}
	}
	return nil
}

func (s *AutoOpsService) validateOpsEventRateClause(
	clause *autoopsproto.OpsEventRateClause,
	localizer locale.Localizer,
) error {
	if clause.VariationId == "" {
		dt, err := statusOpsEventRateClauseVariationIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if clause.GoalId == "" {
		dt, err := statusOpsEventRateClauseGoalIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if clause.MinCount <= 0 {
		dt, err := statusOpsEventRateClauseMinCountRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "min_count"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if clause.ThreadsholdRate > 1 || clause.ThreadsholdRate <= 0 {
		dt, err := statusOpsEventRateClauseInvalidThredshold.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "threshold"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if clause.ActionType == autoopsproto.ActionType_UNKNOWN || clause.ActionType == autoopsproto.ActionType_ENABLE {
		dt, err := statusIncompatibleOpsType.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "action_type"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *AutoOpsService) validateDatetimeClauses(
	clauses []*autoopsproto.DatetimeClause,
	localizer locale.Localizer,
) error {
	checkTimes := make(map[int64]bool)
	for _, c := range clauses {
		if checkTimes[c.Time] {
			dt, err := statusDatetimeClauseDuplicateTime.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		if err := s.validateDatetimeClause(c, localizer); err != nil {
			return err
		}
		checkTimes[c.Time] = true
	}
	return nil
}

func (s *AutoOpsService) validateDatetimeClause(
	clause *autoopsproto.DatetimeClause,
	localizer locale.Localizer) error {
	if clause.Time <= time.Now().Unix() {
		dt, err := statusDatetimeClauseInvalidTime.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if clause.ActionType == autoopsproto.ActionType_UNKNOWN {
		dt, err := statusIncompatibleOpsType.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "action_type"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *AutoOpsService) StopAutoOpsRule(
	ctx context.Context,
	req *autoopsproto.StopAutoOpsRuleRequest,
) (*autoopsproto.StopAutoOpsRuleResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}

	if err := validateStopAutoOpsRuleRequest(req, localizer); err != nil {
		return nil, err
	}

	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(s.mysqlClient)
		autoOpsRule, err := autoOpsRuleStorage.GetAutoOpsRule(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		if autoOpsRule.IsFinished() {
			dt, err := statusAutoOpsRuleFinished.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InvalidArgumentError),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		handler, err := command.NewAutoOpsCommandHandler(editor, autoOpsRule, s.publisher, req.EnvironmentId)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return autoOpsRuleStorage.UpdateAutoOpsRule(contextWithTx, autoOpsRule, req.EnvironmentId)
	})

	if err != nil {
		return nil, err
	}
	return &autoopsproto.StopAutoOpsRuleResponse{}, nil
}

func validateStopAutoOpsRuleRequest(req *autoopsproto.StopAutoOpsRuleRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
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
	return nil
}

func (s *AutoOpsService) DeleteAutoOpsRule(
	ctx context.Context,
	req *autoopsproto.DeleteAutoOpsRuleRequest,
) (*autoopsproto.DeleteAutoOpsRuleResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteAutoOpsRuleRequest(req, localizer); err != nil {
		return nil, err
	}

	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(s.mysqlClient)
		autoOpsRule, err := autoOpsRuleStorage.GetAutoOpsRule(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		handler, err := command.NewAutoOpsCommandHandler(editor, autoOpsRule, s.publisher, req.EnvironmentId)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return autoOpsRuleStorage.UpdateAutoOpsRule(contextWithTx, autoOpsRule, req.EnvironmentId)
	})
	if err != nil {
		if err == v2as.ErrAutoOpsRuleNotFound || err == v2as.ErrAutoOpsRuleUnexpectedAffectedRows {
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
			"Failed to delete autoOpsRule",
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
	return &autoopsproto.DeleteAutoOpsRuleResponse{}, nil
}

func validateDeleteAutoOpsRuleRequest(req *autoopsproto.DeleteAutoOpsRuleRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
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
	return nil
}

func (s *AutoOpsService) UpdateAutoOpsRule(
	ctx context.Context,
	req *autoopsproto.UpdateAutoOpsRuleRequest,
) (*autoopsproto.UpdateAutoOpsRuleResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateUpdateAutoOpsRuleRequest(req, localizer); err != nil {
		return nil, err
	}
	opsEventRateClauses := []*autoopsproto.OpsEventRateClause{}
	for _, c := range req.AddOpsEventRateClauseCommands {
		opsEventRateClauses = append(opsEventRateClauses, c.OpsEventRateClause)
	}
	for _, c := range req.ChangeOpsEventRateClauseCommands {
		opsEventRateClauses = append(opsEventRateClauses, c.OpsEventRateClause)
	}
	for _, c := range opsEventRateClauses {
		exist, err := s.existGoal(ctx, req.EnvironmentId, c.GoalId)
		if err != nil {
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		if !exist {
			s.logger.Error(
				"Goal does not exist",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.String("environmentId", req.EnvironmentId))...,
			)
			dt, err := statusOpsEventRateClauseGoalNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
	}
	commands := s.createUpdateAutoOpsRuleCommands(req)

	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(s.mysqlClient)
		autoOpsRule, err := autoOpsRuleStorage.GetAutoOpsRule(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}

		if autoOpsRule.IsFinished() || autoOpsRule.IsStopped() {
			dt, err := statusAutoOpsRuleCompleted.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InvalidArgumentError),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		if autoOpsRule.OpsType == autoopsproto.OpsType_SCHEDULE {
			if len(req.AddOpsEventRateClauseCommands) > 0 || len(req.ChangeOpsEventRateClauseCommands) > 0 {
				dt, err := statusIncompatibleOpsType.WithDetails(&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "ops_type"),
				})
				if err != nil {
					return statusInternal.Err()
				}
				return dt.Err()
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
					dt, err := statusDatetimeClauseDuplicateTime.WithDetails(&errdetails.LocalizedMessage{
						Locale:  localizer.GetLocale(),
						Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time"),
					})
					if err != nil {
						return statusInternal.Err()
					}
					return dt.Err()
				}
			}
			for _, c := range req.ChangeDatetimeClauseCommands {
				actionType, hasSameTime := checkTimes[c.DatetimeClause.Time]
				if hasSameTime && actionType == c.DatetimeClause.ActionType {
					dt, err := statusDatetimeClauseDuplicateTime.WithDetails(&errdetails.LocalizedMessage{
						Locale:  localizer.GetLocale(),
						Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time"),
					})
					if err != nil {
						return statusInternal.Err()
					}
					return dt.Err()
				}
			}
		}
		if autoOpsRule.OpsType == autoopsproto.OpsType_EVENT_RATE {
			if len(req.AddDatetimeClauseCommands) > 0 || len(req.ChangeDatetimeClauseCommands) > 0 {
				dt, err := statusIncompatibleOpsType.WithDetails(&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "ops_type"),
				})
				if err != nil {
					return statusInternal.Err()
				}
				return dt.Err()
			}
		}

		if req.DeleteClauseCommands != nil && len(autoOpsRule.Clauses) == len(req.DeleteClauseCommands) &&
			len(req.AddOpsEventRateClauseCommands) == 0 && len(req.AddDatetimeClauseCommands) == 0 {
			// When deleting, at least one Clause must exist.
			dt, err := statusShouldAddMoreClauses.WithDetails(&errdetails.LocalizedMessage{
				Locale: localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(
					locale.InvalidArgumentError,
					"add_event_rate_clause_commands",
					"add_datetime_clause_commands"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
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
		return autoOpsRuleStorage.UpdateAutoOpsRule(contextWithTx, autoOpsRule, req.EnvironmentId)
	})
	if err != nil {
		if err == v2as.ErrAutoOpsRuleNotFound || err == v2as.ErrAutoOpsRuleUnexpectedAffectedRows {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		if status.Code(err) == codes.InvalidArgument {
			return nil, err
		}
		s.logger.Error(
			"Failed to update autoOpsRule",
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
	return &autoopsproto.UpdateAutoOpsRuleResponse{}, nil
}

func (s *AutoOpsService) validateUpdateAutoOpsRuleRequest(
	req *autoopsproto.UpdateAutoOpsRuleRequest,
	localizer locale.Localizer,
) error {
	if req.Id == "" {
		dt, err := statusIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if s.isNoUpdateAutoOpsRuleCommand(req) {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	for _, c := range req.AddOpsEventRateClauseCommands {
		if c.OpsEventRateClause == nil {
			dt, err := statusOpsEventRateClauseRequired.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "ops_event_rate_clause"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		if err := s.validateOpsEventRateClause(c.OpsEventRateClause, localizer); err != nil {
			return err
		}
	}
	for _, c := range req.ChangeOpsEventRateClauseCommands {
		if c.Id == "" {
			dt, err := statusClauseIDRequired.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause_id"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		if c.OpsEventRateClause == nil {
			dt, err := statusOpsEventRateClauseRequired.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "ops_event_rate_clause"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		if err := s.validateOpsEventRateClause(c.OpsEventRateClause, localizer); err != nil {
			return err
		}
	}
	for _, c := range req.DeleteClauseCommands {
		if c.Id == "" {
			dt, err := statusClauseIDRequired.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause_id"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
	}

	var checkDatetimeClauses []*autoopsproto.DatetimeClause
	for _, c := range req.AddDatetimeClauseCommands {
		if c.DatetimeClause == nil {
			dt, err := statusDatetimeClauseRequired.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "datetime_clause"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		checkDatetimeClauses = append(checkDatetimeClauses, c.DatetimeClause)
	}
	if err := s.validateDatetimeClauses(checkDatetimeClauses, localizer); err != nil {
		return err
	}

	for _, c := range req.ChangeDatetimeClauseCommands {
		if c.Id == "" {
			dt, err := statusClauseIDRequired.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause_id"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		if c.DatetimeClause == nil {
			dt, err := statusDatetimeClauseRequired.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "datetime_clause"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		checkDatetimeClauses = append(checkDatetimeClauses, c.DatetimeClause)
	}
	if err := s.validateDatetimeClauses(checkDatetimeClauses, localizer); err != nil {
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
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateGetAutoOpsRuleRequest(req, localizer); err != nil {
		return nil, err
	}
	autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(s.mysqlClient)
	autoOpsRule, err := autoOpsRuleStorage.GetAutoOpsRule(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		if err == v2as.ErrAutoOpsRuleNotFound {
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
	if autoOpsRule.Deleted {
		dt, err := statusAlreadyDeleted.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.AlreadyDeletedError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &autoopsproto.GetAutoOpsRuleResponse{
		AutoOpsRule: autoOpsRule.AutoOpsRule,
	}, nil
}

func (s *AutoOpsService) validateGetAutoOpsRuleRequest(
	req *autoopsproto.GetAutoOpsRuleRequest,
	localizer locale.Localizer,
) error {
	if req.Id == "" {
		dt, err := statusIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *AutoOpsService) ListAutoOpsRules(
	ctx context.Context,
	req *autoopsproto.ListAutoOpsRulesRequest,
) (*autoopsproto.ListAutoOpsRulesResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(s.mysqlClient)
	autoOpsRules, cursor, err := s.listAutoOpsRules(
		ctx,
		req.PageSize,
		req.Cursor,
		req.FeatureIds,
		req.EnvironmentId,
		localizer,
		autoOpsRuleStorage,
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
	localizer locale.Localizer,
	storage v2as.AutoOpsRuleStorage,
) ([]*autoopsproto.AutoOpsRule, string, error) {
	whereParts := []mysql.WherePart{
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_id", "=", environmentId),
	}
	fIDs := make([]interface{}, 0, len(featureIds))
	for _, fID := range featureIds {
		fIDs = append(fIDs, fID)
	}
	if len(fIDs) > 0 {
		whereParts = append(whereParts, mysql.NewInFilter("feature_id", fIDs))
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
			return nil, "", statusInternal.Err()
		}
		return nil, "", dt.Err()

	}
	autoOpsRules, nextCursor, err := storage.ListAutoOpsRules(
		ctx,
		whereParts,
		nil,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list autoOpsRules",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, "", statusInternal.Err()
		}
		return nil, "", dt.Err()
	}
	return autoOpsRules, strconv.Itoa(nextCursor), nil
}

func (s *AutoOpsService) ExecuteAutoOps(
	ctx context.Context,
	req *autoopsproto.ExecuteAutoOpsRequest,
) (*autoopsproto.ExecuteAutoOpsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateExecuteAutoOpsRequest(req, localizer); err != nil {
		return nil, err
	}
	triggered, err := s.checkIfHasAlreadyTriggered(ctx, localizer, req.Id, req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if triggered {
		return &autoopsproto.ExecuteAutoOpsResponse{AlreadyTriggered: true}, nil
	}

	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, tx mysql.Transaction) error {
		autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(s.mysqlClient)
		autoOpsRule, err := autoOpsRuleStorage.GetAutoOpsRule(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}

		var executeClause *autoopsproto.Clause = nil
		for _, c := range autoOpsRule.Clauses {
			if c.Id == req.ExecuteAutoOpsRuleCommand.ClauseId {
				executeClause = c
				break
			}
		}

		if executeClause == nil {
			dt, err := statusClauseNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}

		ftStorage := ftstorage.NewFeatureStorage(tx)
		feature, err := ftStorage.GetFeature(contextWithTx, autoOpsRule.FeatureId, req.EnvironmentId)
		if err != nil {
			return err
		}
		prStorage := v2as.NewProgressiveRolloutStorage(s.mysqlClient)
		// Stop the running progressive rollout if the operation type is disable
		if executeClause.ActionType == autoopsproto.ActionType_DISABLE {
			if err := s.stopProgressiveRollout(
				contextWithTx,
				req.EnvironmentId,
				autoOpsRule,
				prStorage,
				localizer,
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
			localizer,
		); err != nil {
			s.logger.Error(
				"Failed to execute auto ops rule operation",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
					zap.String("autoOpsRuleId", autoOpsRule.Id),
					zap.String("featureId", autoOpsRule.FeatureId),
				)...,
			)
			return err
		}
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

		if err = autoOpsRuleStorage.UpdateAutoOpsRule(contextWithTx, autoOpsRule, req.EnvironmentId); err != nil {
			if err == v2as.ErrAutoOpsRuleUnexpectedAffectedRows {
				s.logger.Warn(
					"No rows were affected",
					log.FieldsFromImcomingContext(ctx).AddFields(
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
		if err == v2as.ErrAutoOpsRuleNotFound {
			s.logger.Warn(
				"Auto Ops Rule not found",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", req.Id),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
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
			"Failed to execute autoOpsRule",
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
	return &autoopsproto.ExecuteAutoOpsResponse{AlreadyTriggered: false}, nil
}

func (s *AutoOpsService) stopProgressiveRollout(
	ctx context.Context,
	environmentId string,
	autoOpsRule *domain.AutoOpsRule,
	storage v2as.ProgressiveRolloutStorage,
	localizer locale.Localizer,
) error {
	// Check what operation is being executed
	// and the set progressive rollout `stoppedBy`
	var stoppedBy autoopsproto.ProgressiveRollout_StoppedBy
	hasScheduleOps, err := autoOpsRule.HasScheduleOps()
	if err != nil {
		s.logger.Error(
			"Failed to check operation type",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
				zap.String("autoOpsRuleId", autoOpsRule.Id),
				zap.String("featureId", autoOpsRule.FeatureId),
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
	if hasScheduleOps {
		stoppedBy = autoopsproto.ProgressiveRollout_OPS_SCHEDULE
	} else {
		stoppedBy = autoopsproto.ProgressiveRollout_OPS_KILL_SWITCH
	}
	if err := executeStopProgressiveRolloutOperation(
		ctx,
		storage,
		s.convToInterfaceSlice([]string{autoOpsRule.FeatureId}),
		environmentId,
		stoppedBy,
	); err != nil {
		s.logger.Error(
			"Failed to stop progressive rollout",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
				zap.String("autoOpsRuleId", autoOpsRule.Id),
				zap.String("featureId", autoOpsRule.FeatureId),
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

func (s *AutoOpsService) validateExecuteAutoOpsRequest(
	req *autoopsproto.ExecuteAutoOpsRequest,
	localizer locale.Localizer,
) error {
	if req.Id == "" {
		dt, err := statusIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.ExecuteAutoOpsRuleCommand == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.ExecuteAutoOpsRuleCommand != nil && req.ExecuteAutoOpsRuleCommand.ClauseId == "" {
		dt, err := statusClauseRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *AutoOpsService) checkIfHasAlreadyTriggered(
	ctx context.Context,
	localizer locale.Localizer,
	ruleID,
	environmentId string,
) (bool, error) {
	storage := v2as.NewAutoOpsRuleStorage(s.mysqlClient)
	autoOpsRule, err := storage.GetAutoOpsRule(ctx, ruleID, environmentId)
	if err != nil {
		if err == v2as.ErrAutoOpsRuleNotFound {
			s.logger.Warn(
				"Auto Ops Rule not found",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("ruleID", ruleID),
					zap.String("environmentId", environmentId),
				)...,
			)
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return false, statusInternal.Err()
			}
			return false, dt.Err()
		}
		s.logger.Error(
			"Failed to get auto ops rule",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return false, statusInternal.Err()
		}
		return false, dt.Err()
	}
	if autoOpsRule.IsFinished() || autoOpsRule.IsStopped() || autoOpsRule.Deleted {
		s.logger.Warn(
			"Auto Ops Rule already triggered",
			log.FieldsFromImcomingContext(ctx).AddFields(
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
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
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
		localizer,
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
	localizer locale.Localizer,
) ([]*autoopsproto.OpsCount, string, error) {
	whereParts := []mysql.WherePart{
		mysql.NewFilter("environment_id", "=", environmentId),
	}
	fIDs := make([]interface{}, 0, len(featureIDs))
	for _, fID := range featureIDs {
		fIDs = append(fIDs, fID)
	}
	if len(fIDs) > 0 {
		whereParts = append(whereParts, mysql.NewInFilter("feature_id", fIDs))
	}
	aIDs := make([]interface{}, 0, len(autoOpsRuleIDs))
	for _, aID := range autoOpsRuleIDs {
		aIDs = append(aIDs, aID)
	}
	if len(aIDs) > 0 {
		whereParts = append(whereParts, mysql.NewInFilter("auto_ops_rule_id", aIDs))
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
			return nil, "", statusInternal.Err()
		}
		return nil, "", dt.Err()

	}
	opsCountStorage := v2os.NewOpsCountStorage(s.mysqlClient)
	opsCounts, nextCursor, err := opsCountStorage.ListOpsCounts(
		ctx,
		whereParts,
		nil,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list opsCounts",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, "", statusInternal.Err()
		}
		return nil, "", dt.Err()
	}
	return opsCounts, strconv.Itoa(nextCursor), nil
}

func (s *AutoOpsService) existGoal(ctx context.Context, environmentId string, goalID string) (bool, error) {
	_, err := s.getGoal(ctx, environmentId, goalID)
	if err != nil {
		if err == storage.ErrKeyNotFound {
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
			log.FieldsFromImcomingContext(ctx).AddFields(
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
	localizer locale.Localizer,
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
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.UnauthenticatedError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			dt, err := statusPermissionDenied.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.PermissionDenied),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
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
	}
	return editor, nil
}
