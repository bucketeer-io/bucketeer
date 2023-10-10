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

package webhookhandler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/autoops/command"
	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/pkg/crypto"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"

	authclient "github.com/bucketeer-io/bucketeer/pkg/auth/client"
	autoopsapi "github.com/bucketeer-io/bucketeer/pkg/autoops/api"
	autoopsdomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	v2as "github.com/bucketeer-io/bucketeer/pkg/autoops/storage/v2"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	event "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

const (
	urlParamKeyAuth = "auth"
)

var (
	errAuthKeyEmpty     = errors.New("autoops: auth key is empty")
	errAlreadyTriggered = errors.New("autoops: rule has already triggered")
	errPermissionDenied = errors.New("autoops: permission denied")
)

func WithLogger(logger *zap.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

type handler struct {
	mysqlClient       mysql.Client
	authClient        authclient.Client
	featureClient     featureclient.Client
	publisher         publisher.Publisher
	editor            *event.Editor
	webhookCryptoUtil crypto.EncrypterDecrypter
	logger            *zap.Logger
}

type Option func(*options)

type options struct {
	logger *zap.Logger
}

func NewHandler(
	mysqlClient mysql.Client,
	authClient authclient.Client,
	featureClient featureclient.Client,
	publisher publisher.Publisher,
	verifier token.Verifier,
	tokenPath string,
	webhookCryptoUtil crypto.EncrypterDecrypter,
	opts ...Option,
) (*handler, error) {
	options := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(options)
	}
	data, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, err
	}
	token, err := verifier.Verify(strings.TrimSpace(string(data)))
	if err != nil {
		return nil, err
	}
	if !token.IsAdmin() {
		return nil, errPermissionDenied
	}
	editor := &event.Editor{
		Email:   token.Email,
		Role:    accountproto.Account_OWNER,
		IsAdmin: true,
	}
	return &handler{
		mysqlClient:       mysqlClient,
		authClient:        authClient,
		featureClient:     featureClient,
		publisher:         publisher,
		editor:            editor,
		webhookCryptoUtil: webhookCryptoUtil,
		logger:            options.logger.Named("webhookhandler"),
	}, nil
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	localizer := locale.NewLocalizer(ctx)
	if ctx.Err() == context.Canceled {
		h.logger.Warn(
			"Request was canceled",
			log.FieldsFromImcomingContext(ctx)...,
		)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	secret, err := validateParams(req)
	if err != nil {
		h.logger.Warn(
			"Invalid url parameters",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	ws, err := h.authWebhook(ctx, secret)
	if err != nil {
		h.logger.Error(
			"Failed to get webhook configuration",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	tx, err := h.mysqlClient.BeginTx(ctx)
	if err != nil {
		h.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = h.mysqlClient.RunInTransaction(ctx, tx, func() error {
		webhookStorage := v2as.NewWebhookStorage(tx)
		webhook, err := webhookStorage.GetWebhook(ctx, ws.GetWebhookID(), ws.GetEnvironmentNamespace())
		if err != nil {
			return err
		}
		autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(tx)
		whereParts := []mysql.WherePart{
			mysql.NewFilter("deleted", "=", false),
			mysql.NewFilter("environment_namespace", "=", ws.GetEnvironmentNamespace()),
		}
		autoOpsRules, _, err := autoOpsRuleStorage.ListAutoOpsRules(
			ctx,
			whereParts,
			nil,
			mysql.QueryNoLimit,
			mysql.QueryNoOffset,
		)
		if err != nil {
			return err
		}
		var payload interface{}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			return err
		}
		// Handle webhook and assesses all rules
		// and return the last occurred error.
		var lastErr error
		for _, r := range autoOpsRules {
			rule := &autoopsdomain.AutoOpsRule{AutoOpsRule: r}
			asmt, err := h.assessAutoOpsRule(ctx, rule, webhook.Id, payload)
			if err != nil {
				lastErr = err
			}
			if asmt {
				if err = h.executeAutoOps(ctx, rule, ws.GetEnvironmentNamespace(), autoOpsRuleStorage, localizer); err != nil {
					lastErr = err
				}
			}
		}
		return lastErr
	})
	if err != nil {
		h.logger.Error(
			"Failed to execute autoOpsRule",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", ws.GetEnvironmentNamespace()),
			)...,
		)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func validateParams(req *http.Request) (string, error) {
	secret := req.URL.Query().Get(urlParamKeyAuth)
	if secret == "" {
		return "", errAuthKeyEmpty
	}
	return secret, nil
}

func (h *handler) authWebhook(
	ctx context.Context,
	secret string,
) (domain.WebhookSecret, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(secret)
	if err != nil {
		h.logger.Error(
			"Failed to decode encrypted secret",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	decrypted, err := h.webhookCryptoUtil.Decrypt(ctx, decoded)
	if err != nil {
		h.logger.Error(
			"Failed to decrypt encrypted secret",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	ws, err := domain.UnmarshalWebhookSecret(decrypted)
	if err != nil {
		h.logger.Error(
			"Failed to unmarshal webhook secret",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	return ws, nil
}

func (h *handler) assessAutoOpsRule(
	ctx context.Context,
	a *autoopsdomain.AutoOpsRule,
	tarId string,
	payload interface{},
) (bool, error) {
	webhookClauses, err := a.ExtractWebhookClauses()
	if err != nil {
		h.logger.Error("Failed to extract webhook clauses",
			zap.Error(err),
			zap.String("featureId", a.FeatureId),
			zap.String("autoOpsRuleId", a.Id),
		)
		return false, err
	}
	var lastErr error
	// All clauses are combined with implicit OR
	for _, w := range webhookClauses {
		if w.WebhookId != tarId {
			continue
		}
		asmt, err := evaluateClause(ctx, w, payload)
		if err != nil {
			h.logger.Error("Skipping evaluation because an error has occurred",
				zap.Error(err),
				zap.String("featureId", a.FeatureId),
				zap.String("autoOpsRuleId", a.Id),
			)
			lastErr = err
			continue
		}
		if asmt {
			h.logger.Info("Clause satisfies condition",
				zap.String("featureId", a.FeatureId),
				zap.String("autoOpsRuleId", a.Id),
				zap.Any("webhookClause", w),
			)
			return true, lastErr
		}
	}
	return false, lastErr
}

func (h *handler) executeAutoOps(
	ctx context.Context,
	rule *autoopsdomain.AutoOpsRule,
	environmentNamespace string,
	storage v2as.AutoOpsRuleStorage,
	localizer locale.Localizer,
) error {
	if rule.AlreadyTriggered() {
		return errAlreadyTriggered
	}
	handler := command.NewAutoOpsCommandHandler(h.editor, rule, h.publisher, environmentNamespace)
	if err := handler.Handle(ctx, &autoopsproto.ChangeAutoOpsRuleTriggeredAtCommand{}); err != nil {
		return err
	}
	if err := storage.UpdateAutoOpsRule(ctx, rule, environmentNamespace); err != nil {
		return err
	}
	return autoopsapi.ExecuteAutoOpsRuleOperation(ctx, environmentNamespace, rule, h.featureClient, h.logger, localizer)
}
