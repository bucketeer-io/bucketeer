// Copyright 2022 The Bucketeer Authors.
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

var errAlreadyTriggered = errors.New("auto ops Rule has already triggered")

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
	editor, err := s.checkRole(ctx, accountproto.Account_EDITOR, req.EnvironmentNamespace)
	if err != nil {
		return nil, err
	}
	if err := s.validateCreateAutoOpsRuleRequest(req); err != nil {
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
		return nil, localizedError(statusInternal, locale.JaJP)
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
		return nil, localizedError(statusInternal, locale.JaJP)
	}
	for _, c := range opsEventRateClauses {
		exist, err := s.existGoal(ctx, req.EnvironmentNamespace, c.GoalId)
		if err != nil {
			return nil, localizedError(statusInternal, locale.JaJP)
		}
		if !exist {
			s.logger.Error(
				"Goal does not exist",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.String("environmentNamespace", req.EnvironmentNamespace))...,
			)
			return nil, localizedError(statusOpsEventRateClauseGoalNotFound, locale.JaJP)
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
		return nil, localizedError(statusInternal, locale.JaJP)
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
			return nil, localizedError(statusAlreadyExists, locale.JaJP)
		}
		s.logger.Error(
			"Failed to create autoOps",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, localizedError(statusInternal, locale.JaJP)
	}
	return &autoopsproto.CreateAutoOpsRuleResponse{}, nil
}

func (s *AutoOpsService) validateCreateAutoOpsRuleRequest(req *autoopsproto.CreateAutoOpsRuleRequest) error {
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	if req.Command.FeatureId == "" {
		return localizedError(statusFeatureIDRequired, locale.JaJP)
	}
	if len(req.Command.OpsEventRateClauses) == 0 && len(req.Command.DatetimeClauses) == 0 && len(req.Command.WebhookClauses) == 0 {
		return localizedError(statusClauseRequired, locale.JaJP)
	}
	if req.Command.OpsType == autoopsproto.OpsType_ENABLE_FEATURE && len(req.Command.OpsEventRateClauses) > 0 {
		return localizedError(statusIncompatibleOpsType, locale.JaJP)
	}
	if err := s.validateOpsEventRateClauses(req.Command.OpsEventRateClauses); err != nil {
		return err
	}
	if err := s.validateDatetimeClauses(req.Command.DatetimeClauses); err != nil {
		return err
	}
	if err := s.validateWebhookClauses(req.Command.WebhookClauses); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) validateOpsEventRateClauses(clauses []*autoopsproto.OpsEventRateClause) error {
	for _, c := range clauses {
		if err := s.validateOpsEventRateClause(c); err != nil {
			return err
		}
	}
	return nil
}

func (s *AutoOpsService) validateOpsEventRateClause(clause *autoopsproto.OpsEventRateClause) error {
	if clause.VariationId == "" {
		return localizedError(statusOpsEventRateClauseVariationIDRequired, locale.JaJP)
	}
	if clause.GoalId == "" {
		return localizedError(statusOpsEventRateClauseGoalIDRequired, locale.JaJP)
	}
	if clause.MinCount <= 0 {
		return localizedError(statusOpsEventRateClauseMinCountRequired, locale.JaJP)
	}
	if clause.ThreadsholdRate > 1 || clause.ThreadsholdRate <= 0 {
		return localizedError(statusOpsEventRateClauseInvalidThredshold, locale.JaJP)
	}
	return nil
}

func (s *AutoOpsService) validateDatetimeClauses(clauses []*autoopsproto.DatetimeClause) error {
	for _, c := range clauses {
		if err := s.validateDatetimeClause(c); err != nil {
			return err
		}
	}
	return nil
}

func (s *AutoOpsService) validateDatetimeClause(clause *autoopsproto.DatetimeClause) error {
	if clause.Time <= time.Now().Unix() {
		return localizedError(statusDatetimeClauseInvalidTime, locale.JaJP)
	}
	return nil
}

func (s *AutoOpsService) validateWebhookClauses(clauses []*autoopsproto.WebhookClause) error {
	for _, c := range clauses {
		if err := s.validateWebhookClause(c); err != nil {
			return err
		}
	}
	return nil
}

func (s *AutoOpsService) validateWebhookClause(clause *autoopsproto.WebhookClause) error {
	if clause.WebhookId == "" {
		return localizedError(statusWebhookClauseWebhookIDRequired, locale.JaJP)
	}
	if len(clause.Conditions) == 0 {
		return localizedError(statusWebhookClauseConditionRequired, locale.JaJP)
	}
	for _, c := range clause.Conditions {
		if c.Filter == "" {
			return localizedError(statusWebhookClauseConditionFilterRequired, locale.JaJP)
		}
		if c.Value == "" {
			return localizedError(statusWebhookClauseConditionValueRequired, locale.JaJP)
		}
		_, ok := autoopsproto.WebhookClause_Condition_Operator_name[int32(c.Operator)]
		if !ok {
			return localizedError(statusWebhookClauseConditionInvalidOperator, locale.JaJP)
		}
	}
	return nil
}

func (s *AutoOpsService) DeleteAutoOpsRule(
	ctx context.Context,
	req *autoopsproto.DeleteAutoOpsRuleRequest,
) (*autoopsproto.DeleteAutoOpsRuleResponse, error) {
	editor, err := s.checkRole(ctx, accountproto.Account_EDITOR, req.EnvironmentNamespace)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteAutoOpsRuleRequest(req); err != nil {
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
		return nil, localizedError(statusInternal, locale.JaJP)
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
			return nil, localizedError(statusNotFound, locale.JaJP)
		}
		s.logger.Error(
			"Failed to delete autoOpsRule",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, localizedError(statusInternal, locale.JaJP)
	}
	return &autoopsproto.DeleteAutoOpsRuleResponse{}, nil
}

func validateDeleteAutoOpsRuleRequest(req *autoopsproto.DeleteAutoOpsRuleRequest) error {
	if req.Id == "" {
		return localizedError(statusIDRequired, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	return nil
}

func (s *AutoOpsService) UpdateAutoOpsRule(
	ctx context.Context,
	req *autoopsproto.UpdateAutoOpsRuleRequest,
) (*autoopsproto.UpdateAutoOpsRuleResponse, error) {
	editor, err := s.checkRole(ctx, accountproto.Account_EDITOR, req.EnvironmentNamespace)
	if err != nil {
		return nil, err
	}
	if err := s.validateUpdateAutoOpsRuleRequest(req); err != nil {
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
			return nil, localizedError(statusInternal, locale.JaJP)
		}
		if !exist {
			s.logger.Error(
				"Goal does not exist",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.String("environmentNamespace", req.EnvironmentNamespace))...,
			)
			return nil, localizedError(statusOpsEventRateClauseGoalNotFound, locale.JaJP)
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
		return nil, localizedError(statusInternal, locale.JaJP)
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
				return localizedError(statusIncompatibleOpsType, locale.JaJP)
			}
		} else if autoOpsRule.OpsType == autoopsproto.OpsType_ENABLE_FEATURE && len(req.AddOpsEventRateClauseCommands) > 0 {
			return localizedError(statusIncompatibleOpsType, locale.JaJP)
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
			return nil, localizedError(statusNotFound, locale.JaJP)
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
		return nil, localizedError(statusInternal, locale.JaJP)
	}
	return &autoopsproto.UpdateAutoOpsRuleResponse{}, nil
}

func (s *AutoOpsService) validateUpdateAutoOpsRuleRequest(req *autoopsproto.UpdateAutoOpsRuleRequest) error {
	if req.Id == "" {
		return localizedError(statusIDRequired, locale.JaJP)
	}
	if s.isNoUpdateAutoOpsRuleCommand(req) {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	for _, c := range req.AddOpsEventRateClauseCommands {
		if c.OpsEventRateClause == nil {
			return localizedError(statusOpsEventRateClauseRequired, locale.JaJP)
		}
		if err := s.validateOpsEventRateClause(c.OpsEventRateClause); err != nil {
			return err
		}
	}
	for _, c := range req.ChangeOpsEventRateClauseCommands {
		if c.Id == "" {
			return localizedError(statusClauseIDRequired, locale.JaJP)
		}
		if c.OpsEventRateClause == nil {
			return localizedError(statusOpsEventRateClauseRequired, locale.JaJP)
		}
		if err := s.validateOpsEventRateClause(c.OpsEventRateClause); err != nil {
			return err
		}
	}
	for _, c := range req.DeleteClauseCommands {
		if c.Id == "" {
			return localizedError(statusClauseIDRequired, locale.JaJP)
		}
	}
	for _, c := range req.AddDatetimeClauseCommands {
		if c.DatetimeClause == nil {
			return localizedError(statusDatetimeClauseRequired, locale.JaJP)
		}
		if err := s.validateDatetimeClause(c.DatetimeClause); err != nil {
			return err
		}
	}
	for _, c := range req.ChangeDatetimeClauseCommands {
		if c.Id == "" {
			return localizedError(statusClauseIDRequired, locale.JaJP)
		}
		if c.DatetimeClause == nil {
			return localizedError(statusDatetimeClauseRequired, locale.JaJP)
		}
		if err := s.validateDatetimeClause(c.DatetimeClause); err != nil {
			return err
		}
	}
	for _, c := range req.AddWebhookClauseCommands {
		if c.WebhookClause == nil {
			return localizedError(statusWebhookClauseRequired, locale.JaJP)
		}
		if err := s.validateWebhookClause(c.WebhookClause); err != nil {
			return err
		}
	}
	for _, c := range req.ChangeWebhookClauseCommands {
		if c.Id == "" {
			return localizedError(statusClauseIDRequired, locale.JaJP)
		}
		if c.WebhookClause == nil {
			return localizedError(statusWebhookClauseRequired, locale.JaJP)
		}
		if err := s.validateWebhookClause(c.WebhookClause); err != nil {
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
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace)
	if err != nil {
		return nil, err
	}
	if err := s.validateGetAutoOpsRuleRequest(req); err != nil {
		return nil, err
	}
	autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(s.mysqlClient)
	autoOpsRule, err := autoOpsRuleStorage.GetAutoOpsRule(ctx, req.Id, req.EnvironmentNamespace)
	if err != nil {
		if err == v2as.ErrAutoOpsRuleNotFound {
			return nil, localizedError(statusNotFound, locale.JaJP)
		}
		return nil, localizedError(statusInternal, locale.JaJP)
	}
	if autoOpsRule.Deleted {
		return nil, localizedError(statusAlreadyDeleted, locale.JaJP)
	}
	return &autoopsproto.GetAutoOpsRuleResponse{
		AutoOpsRule: autoOpsRule.AutoOpsRule,
	}, nil
}

func (s *AutoOpsService) validateGetAutoOpsRuleRequest(req *autoopsproto.GetAutoOpsRuleRequest) error {
	if req.Id == "" {
		return localizedError(statusIDRequired, locale.JaJP)
	}
	return nil
}

func (s *AutoOpsService) ListAutoOpsRules(
	ctx context.Context,
	req *autoopsproto.ListAutoOpsRulesRequest,
) (*autoopsproto.ListAutoOpsRulesResponse, error) {
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace)
	if err != nil {
		return nil, err
	}
	autoOpsRules, cursor, err := s.listAutoOpsRules(
		ctx,
		req.PageSize,
		req.Cursor,
		req.FeatureIds,
		req.EnvironmentNamespace)
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
		return nil, "", localizedError(statusInvalidCursor, locale.JaJP)
	}
	autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(s.mysqlClient)
	autoOpsRules, nextCursor, err := autoOpsRuleStorage.ListAutoOpsRules(
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
		return nil, "", localizedError(statusInternal, locale.JaJP)
	}
	return autoOpsRules, strconv.Itoa(nextCursor), nil
}

func (s *AutoOpsService) ExecuteAutoOps(
	ctx context.Context,
	req *autoopsproto.ExecuteAutoOpsRequest,
) (*autoopsproto.ExecuteAutoOpsResponse, error) {
	editor, err := s.checkRole(ctx, accountproto.Account_EDITOR, req.EnvironmentNamespace)
	if err != nil {
		return nil, err
	}
	if err := s.validateExecuteAutoOpsRequest(req); err != nil {
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
		return nil, localizedError(statusInternal, locale.JaJP)
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(tx)
		autoOpsRule, err := autoOpsRuleStorage.GetAutoOpsRule(ctx, req.Id, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		if autoOpsRule.AlreadyTriggered() {
			return errAlreadyTriggered
		}
		handler := command.NewAutoOpsCommandHandler(editor, autoOpsRule, s.publisher, req.EnvironmentNamespace)
		if err := handler.Handle(ctx, req.ChangeAutoOpsRuleTriggeredAtCommand); err != nil {
			return err
		}
		if err = autoOpsRuleStorage.UpdateAutoOpsRule(ctx, autoOpsRule, req.EnvironmentNamespace); err != nil {
			return err
		}
		return ExecuteOperation(ctx, req.EnvironmentNamespace, autoOpsRule, s.featureClient, s.logger)
	})
	if err != nil {
		if err == errAlreadyTriggered {
			return &autoopsproto.ExecuteAutoOpsResponse{AlreadyTriggered: true}, nil
		}
		if err == v2as.ErrAutoOpsRuleNotFound || err == v2as.ErrAutoOpsRuleUnexpectedAffectedRows {
			return nil, localizedError(statusNotFound, locale.JaJP)
		}
		s.logger.Error(
			"Failed to execute autoOpsRule",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, localizedError(statusInternal, locale.JaJP)
	}
	return &autoopsproto.ExecuteAutoOpsResponse{AlreadyTriggered: false}, nil
}

func (s *AutoOpsService) validateExecuteAutoOpsRequest(req *autoopsproto.ExecuteAutoOpsRequest) error {
	if req.Id == "" {
		return localizedError(statusIDRequired, locale.JaJP)
	}
	if req.ChangeAutoOpsRuleTriggeredAtCommand == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	return nil
}

func (s *AutoOpsService) ListOpsCounts(
	ctx context.Context,
	req *autoopsproto.ListOpsCountsRequest,
) (*autoopsproto.ListOpsCountsResponse, error) {
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace)
	if err != nil {
		return nil, err
	}
	opsCounts, cursor, err := s.listOpsCounts(
		ctx,
		req.PageSize,
		req.Cursor,
		req.EnvironmentNamespace,
		req.FeatureIds,
		req.AutoOpsRuleIds)
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
		return nil, "", localizedError(statusInvalidCursor, locale.JaJP)
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
		return nil, "", localizedError(statusInternal, locale.JaJP)
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
			return nil, localizedError(statusUnauthenticated, locale.JaJP)
		case codes.PermissionDenied:
			s.logger.Info(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			return nil, localizedError(statusPermissionDenied, locale.JaJP)
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			return nil, localizedError(statusInternal, locale.JaJP)
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
