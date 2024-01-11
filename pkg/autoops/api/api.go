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
	"net/url"
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
	"github.com/bucketeer-io/bucketeer/pkg/crypto"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
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
	mysqlClient       mysql.Client
	featureClient     featureclient.Client
	experimentClient  experimentclient.Client
	accountClient     accountclient.Client
	authClient        authclient.Client
	publisher         publisher.Publisher
	webhookBaseURL    *url.URL
	webhookCryptoUtil crypto.EncrypterDecrypter
	opts              *options
	logger            *zap.Logger
}

func NewAutoOpsService(
	mysqlClient mysql.Client,
	featureClient featureclient.Client,
	experimentClient experimentclient.Client,
	accountClient accountclient.Client,
	authClient authclient.Client,
	publisher publisher.Publisher,
	webhookBaseURL *url.URL,
	webhookCryptoUtil crypto.EncrypterDecrypter,
	opts ...Option,
) *AutoOpsService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &AutoOpsService{
		mysqlClient:       mysqlClient,
		featureClient:     featureClient,
		experimentClient:  experimentClient,
		accountClient:     accountClient,
		authClient:        authClient,
		publisher:         publisher,
		webhookBaseURL:    webhookBaseURL,
		opts:              dopts,
		webhookCryptoUtil: webhookCryptoUtil,
		logger:            dopts.logger.Named("api"),
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
	editor, err := s.checkRole(ctx, accountproto.Account_EDITOR, req.EnvironmentNamespace, localizer)
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
		req.Command.WebhookClauses,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create a new autoOpsRule",
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
	opsEventRateClauses, err := autoOpsRule.ExtractOpsEventRateClauses()
	if err != nil {
		s.logger.Error(
			"Failed to extract opsEventRateClauses",
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
	for _, c := range opsEventRateClauses {
		exist, err := s.existGoal(ctx, req.EnvironmentNamespace, c.GoalId)
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
				log.FieldsFromImcomingContext(ctx).AddFields(zap.String("environmentNamespace", req.EnvironmentNamespace))...,
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
		autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(tx)
		handler := command.NewAutoOpsCommandHandler(editor, autoOpsRule, s.publisher, req.EnvironmentNamespace)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return autoOpsRuleStorage.CreateAutoOpsRule(ctx, autoOpsRule, req.EnvironmentNamespace)
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
	return &autoopsproto.CreateAutoOpsRuleResponse{}, nil
}

func (s *AutoOpsService) existsRunningProgressiveRollout(
	ctx context.Context,
	featureID, environmentNamespace string,
	localizer locale.Localizer,
) (bool, error) {
	progressiveRollouts, err := s.listProgressiveRolloutsByFeatureID(
		ctx,
		environmentNamespace, featureID,
		localizer,
	)
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
		dp := &domain.ProgressiveRollout{
			ProgressiveRollout: p,
		}
		if !dp.IsFinished() {
			return true
		}
	}
	return false
}

func (s *AutoOpsService) listProgressiveRolloutsByFeatureID(
	ctx context.Context,
	featureID, environmentNamespace string,
	localizer locale.Localizer,
) ([]*autoopsproto.ProgressiveRollout, error) {
	progressiveRollouts := make([]*autoopsproto.ProgressiveRollout, 0)
	cursor := ""
	for {
		progressiveRollout, _, nextOffset, err := s.listProgressiveRollouts(
			ctx,
			&autoopsproto.ListProgressiveRolloutsRequest{
				EnvironmentNamespace: environmentNamespace,
				PageSize:             listRequestSize,
				Cursor:               cursor,
				FeatureIds:           []string{featureID},
			},
			localizer,
		)
		if err != nil {
			return nil, err
		}
		progressiveRollouts = append(progressiveRollouts, progressiveRollout...)
		size := len(progressiveRollouts)
		if size == 0 || size < listRequestSize {
			return progressiveRollouts, nil
		}
		cursor = strconv.Itoa(nextOffset)
	}
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
		len(req.Command.DatetimeClauses) == 0 &&
		len(req.Command.WebhookClauses) == 0 {
		dt, err := statusClauseRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.OpsType == autoopsproto.OpsType_ENABLE_FEATURE && len(req.Command.OpsEventRateClauses) > 0 {
		dt, err := statusIncompatibleOpsType.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "ops_type"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if err := s.validateOpsEventRateClauses(req.Command.OpsEventRateClauses, localizer); err != nil {
		return err
	}
	if err := s.validateDatetimeClauses(req.Command.DatetimeClauses, localizer); err != nil {
		return err
	}
	if err := s.validateWebhookClauses(req.Command.WebhookClauses, localizer); err != nil {
		return err
	}
	runningProgressiveRolloutExists, err := s.existsRunningProgressiveRollout(
		ctx,
		req.Command.FeatureId,
		req.EnvironmentNamespace,
		localizer,
	)
	if err != nil {
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if runningProgressiveRolloutExists && (len(req.Command.DatetimeClauses) == 0 || len(req.Command.WebhookClauses) == 0) {
		dt, err := statusWaitingOrRunningProgressiveRolloutExists.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.AutoOpsWaitingOrRunningExperimentExists),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
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
	return nil
}

func (s *AutoOpsService) validateDatetimeClauses(
	clauses []*autoopsproto.DatetimeClause,
	localizer locale.Localizer,
) error {
	for _, c := range clauses {
		if err := s.validateDatetimeClause(c, localizer); err != nil {
			return err
		}
	}
	return nil
}

func (s *AutoOpsService) validateDatetimeClause(clause *autoopsproto.DatetimeClause, localizer locale.Localizer) error {
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
	return nil
}

func (s *AutoOpsService) validateWebhookClauses(
	clauses []*autoopsproto.WebhookClause,
	localizer locale.Localizer,
) error {
	for _, c := range clauses {
		if err := s.validateWebhookClause(c, localizer); err != nil {
			return err
		}
	}
	return nil
}

func (s *AutoOpsService) validateWebhookClause(clause *autoopsproto.WebhookClause, localizer locale.Localizer) error {
	if clause.WebhookId == "" {
		dt, err := statusWebhookClauseWebhookIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "webhook_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if len(clause.Conditions) == 0 {
		dt, err := statusWebhookClauseConditionRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "condition"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	for _, c := range clause.Conditions {
		if c.Filter == "" {
			dt, err := statusWebhookClauseConditionFilterRequired.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "condition_filter"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		if c.Value == "" {
			dt, err := statusWebhookClauseConditionValueRequired.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "condition_value"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		_, ok := autoopsproto.WebhookClause_Condition_Operator_name[int32(c.Operator)]
		if !ok {
			dt, err := statusWebhookClauseConditionInvalidOperator.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "condition_operator"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
	}
	return nil
}

func (s *AutoOpsService) DeleteAutoOpsRule(
	ctx context.Context,
	req *autoopsproto.DeleteAutoOpsRuleRequest,
) (*autoopsproto.DeleteAutoOpsRuleResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountproto.Account_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteAutoOpsRuleRequest(req, localizer); err != nil {
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
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(tx)
		autoOpsRule, err := autoOpsRuleStorage.GetAutoOpsRule(ctx, req.Id, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		handler := command.NewAutoOpsCommandHandler(editor, autoOpsRule, s.publisher, req.EnvironmentNamespace)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return autoOpsRuleStorage.UpdateAutoOpsRule(ctx, autoOpsRule, req.EnvironmentNamespace)
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
	editor, err := s.checkRole(ctx, accountproto.Account_EDITOR, req.EnvironmentNamespace, localizer)
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
		exist, err := s.existGoal(ctx, req.EnvironmentNamespace, c.GoalId)
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
				log.FieldsFromImcomingContext(ctx).AddFields(zap.String("environmentNamespace", req.EnvironmentNamespace))...,
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
		autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(tx)
		autoOpsRule, err := autoOpsRuleStorage.GetAutoOpsRule(ctx, req.Id, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		if req.ChangeAutoOpsRuleOpsTypeCommand != nil {
			if req.ChangeAutoOpsRuleOpsTypeCommand.OpsType == autoopsproto.OpsType_ENABLE_FEATURE &&
				len(req.AddOpsEventRateClauseCommands) > 0 {
				dt, err := statusIncompatibleOpsType.WithDetails(&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "ops_type"),
				})
				if err != nil {
					return statusInternal.Err()
				}
				return dt.Err()
			}
		} else if autoOpsRule.OpsType == autoopsproto.OpsType_ENABLE_FEATURE && len(req.AddOpsEventRateClauseCommands) > 0 {
			dt, err := statusIncompatibleOpsType.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "ops_type"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		handler := command.NewAutoOpsCommandHandler(editor, autoOpsRule, s.publisher, req.EnvironmentNamespace)
		for _, command := range commands {
			if err := handler.Handle(ctx, command); err != nil {
				return err
			}
		}
		return autoOpsRuleStorage.UpdateAutoOpsRule(ctx, autoOpsRule, req.EnvironmentNamespace)
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
		if err := s.validateDatetimeClause(c.DatetimeClause, localizer); err != nil {
			return err
		}
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
		if err := s.validateDatetimeClause(c.DatetimeClause, localizer); err != nil {
			return err
		}
	}
	for _, c := range req.AddWebhookClauseCommands {
		if c.WebhookClause == nil {
			dt, err := statusWebhookClauseRequired.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "webhook_clause"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		if err := s.validateWebhookClause(c.WebhookClause, localizer); err != nil {
			return err
		}
	}
	for _, c := range req.ChangeWebhookClauseCommands {
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
		if c.WebhookClause == nil {
			dt, err := statusWebhookClauseRequired.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "webhook_clause"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		if err := s.validateWebhookClause(c.WebhookClause, localizer); err != nil {
			return err
		}
	}
	return nil
}

func (s *AutoOpsService) isNoUpdateAutoOpsRuleCommand(req *autoopsproto.UpdateAutoOpsRuleRequest) bool {
	return req.ChangeAutoOpsRuleOpsTypeCommand == nil &&
		len(req.AddOpsEventRateClauseCommands) == 0 &&
		len(req.ChangeOpsEventRateClauseCommands) == 0 &&
		len(req.DeleteClauseCommands) == 0 &&
		len(req.AddDatetimeClauseCommands) == 0 &&
		len(req.ChangeDatetimeClauseCommands) == 0 &&
		len(req.AddWebhookClauseCommands) == 0 &&
		len(req.ChangeWebhookClauseCommands) == 0
}

func (s *AutoOpsService) createUpdateAutoOpsRuleCommands(req *autoopsproto.UpdateAutoOpsRuleRequest) []command.Command {
	commands := make([]command.Command, 0)
	if req.ChangeAutoOpsRuleOpsTypeCommand != nil {
		commands = append(commands, req.ChangeAutoOpsRuleOpsTypeCommand)
	}
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
	for _, c := range req.AddWebhookClauseCommands {
		commands = append(commands, c)
	}
	for _, c := range req.ChangeWebhookClauseCommands {
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
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateGetAutoOpsRuleRequest(req, localizer); err != nil {
		return nil, err
	}
	autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(s.mysqlClient)
	autoOpsRule, err := autoOpsRuleStorage.GetAutoOpsRule(ctx, req.Id, req.EnvironmentNamespace)
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
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(s.mysqlClient)
	autoOpsRules, cursor, err := s.listAutoOpsRules(
		ctx,
		req.PageSize,
		req.Cursor,
		req.FeatureIds,
		req.EnvironmentNamespace,
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
	environmentNamespace string,
	localizer locale.Localizer,
	storage v2as.AutoOpsRuleStorage,
) ([]*autoopsproto.AutoOpsRule, string, error) {
	whereParts := []mysql.WherePart{
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_namespace", "=", environmentNamespace),
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
				zap.String("environmentNamespace", environmentNamespace),
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
	editor, err := s.checkRole(ctx, accountproto.Account_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateExecuteAutoOpsRequest(req, localizer); err != nil {
		return nil, err
	}
	triggered, err := s.checkIfHasAlreadyTriggered(ctx, localizer, req.Id, req.EnvironmentNamespace)
	if err != nil {
		return nil, err
	}
	if triggered {
		return &autoopsproto.ExecuteAutoOpsResponse{AlreadyTriggered: true}, nil
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
		autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(tx)
		autoOpsRule, err := autoOpsRuleStorage.GetAutoOpsRule(ctx, req.Id, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		handler := command.NewAutoOpsCommandHandler(editor, autoOpsRule, s.publisher, req.EnvironmentNamespace)
		if err := handler.Handle(ctx, req.ChangeAutoOpsRuleTriggeredAtCommand); err != nil {
			return err
		}
		if err = autoOpsRuleStorage.UpdateAutoOpsRule(ctx, autoOpsRule, req.EnvironmentNamespace); err != nil {
			if err == v2as.ErrAutoOpsRuleUnexpectedAffectedRows {
				s.logger.Warn(
					"No rows were affected",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("id", req.Id),
						zap.String("environmentNamespace", req.EnvironmentNamespace),
					)...,
				)
				return nil
			}
			return err
		}
		return ExecuteAutoOpsRuleOperation(ctx, req.EnvironmentNamespace, autoOpsRule, s.featureClient, s.logger, localizer)
	})
	if err != nil {
		if err == v2as.ErrAutoOpsRuleNotFound {
			s.logger.Warn(
				"Auto Ops Rule not found",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", req.Id),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	return &autoopsproto.ExecuteAutoOpsResponse{AlreadyTriggered: false}, nil
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
	if req.ChangeAutoOpsRuleTriggeredAtCommand == nil {
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

func (s *AutoOpsService) checkIfHasAlreadyTriggered(
	ctx context.Context,
	localizer locale.Localizer,
	ruleID,
	environmentNamespace string,
) (bool, error) {
	storage := v2as.NewAutoOpsRuleStorage(s.mysqlClient)
	autoOpsRule, err := storage.GetAutoOpsRule(ctx, ruleID, environmentNamespace)
	if err != nil {
		if err == v2as.ErrAutoOpsRuleNotFound {
			s.logger.Warn(
				"Auto Ops Rule not found",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("ruleID", ruleID),
					zap.String("environmentNamespace", environmentNamespace),
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
				zap.String("environmentNamespace", environmentNamespace),
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
	if autoOpsRule.AlreadyTriggered() {
		s.logger.Warn(
			"Auto Ops Rule already triggered",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("ruleID", ruleID),
				zap.String("environmentNamespace", environmentNamespace),
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
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	opsCounts, cursor, err := s.listOpsCounts(
		ctx,
		req.PageSize,
		req.Cursor,
		req.EnvironmentNamespace,
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
	environmentNamespace string,
	featureIDs []string,
	autoOpsRuleIDs []string,
	localizer locale.Localizer,
) ([]*autoopsproto.OpsCount, string, error) {
	whereParts := []mysql.WherePart{
		mysql.NewFilter("environment_namespace", "=", environmentNamespace),
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
				zap.String("environmentNamespace", environmentNamespace),
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

func (s *AutoOpsService) existGoal(ctx context.Context, environmentNamespace string, goalID string) (bool, error) {
	_, err := s.getGoal(ctx, environmentNamespace, goalID)
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
	environmentNamespace, goalID string,
) (*experimentproto.Goal, error) {
	resp, err := s.experimentClient.GetGoal(ctx, &experimentproto.GetGoalRequest{
		Id:                   goalID,
		EnvironmentNamespace: environmentNamespace,
	})
	if err != nil {
		s.logger.Error(
			"Failed to get goal",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
				zap.String("goalId", goalID),
			)...,
		)
		return nil, err
	}
	return resp.Goal, nil
}

func (s *AutoOpsService) checkRole(
	ctx context.Context,
	requiredRole accountproto.Account_Role,
	environmentNamespace string,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckRole(ctx, requiredRole, func(email string) (*accountproto.GetAccountResponse, error) {
		return s.accountClient.GetAccount(ctx, &accountproto.GetAccountRequest{
			Email:                email,
			EnvironmentNamespace: environmentNamespace,
		})
	})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Info(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
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
			s.logger.Info(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
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
					zap.String("environmentNamespace", environmentNamespace),
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

func (s *AutoOpsService) reportInternalServerError(
	ctx context.Context,
	err error,
	environmentNamespace string,
	localizer locale.Localizer,
) error {
	s.logger.Error(
		"Internal server error",
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
