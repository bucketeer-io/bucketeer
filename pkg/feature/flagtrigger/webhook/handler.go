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
//

package webhook

import (
	"context"
	"encoding/base64"
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/crypto"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type handler struct {
	mysqlClient       mysql.Client
	featureClient     featureclient.Client
	triggerCryptoUtil crypto.EncrypterDecrypter
	logger            *zap.Logger
}

type Option func(*options)

type options struct {
	logger *zap.Logger
}

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

func NewHandler(
	mysqlClient mysql.Client,
	featureClient featureclient.Client,
	triggerCryptoUtil crypto.EncrypterDecrypter,
	opts ...Option,
) *handler {
	options := &options{
		logger: zap.NewNop(),
	}
	for _, o := range opts {
		o(options)
	}
	return &handler{
		mysqlClient:       mysqlClient,
		featureClient:     featureClient,
		triggerCryptoUtil: triggerCryptoUtil,
		logger:            options.logger,
	}
}

func (h handler) ServeHTTP(resp http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	secret := request.URL.Query().Get("secret")
	if secret == "" {
		h.logger.Error(
			"Failed to get secret from query",
			log.FieldsFromImcomingContext(ctx)...,
		)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	triggerSecret, err := h.authSecret(ctx, secret)
	if err != nil {
		h.logger.Error(
			"Failed to auth trigger secret",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	storage := v2fs.NewFlagTriggerStorage(h.mysqlClient)
	trigger, err := storage.GetFlagTrigger(ctx, triggerSecret.GetID(), triggerSecret.GetEnvironmentNamespace())
	if err != nil {
		h.logger.Error(
			"Failed to get flag trigger",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	if trigger.GetDisabled() || trigger.GetDeleted() {
		h.logger.Error(
			"Flag trigger is disabled or deleted",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	// check trigger secret
	if trigger.GetFeatureId() != triggerSecret.GetFeatureID() ||
		int(trigger.GetAction()) != triggerSecret.GetAction() ||
		trigger.GetUuid() != triggerSecret.GetUUID() {
		h.logger.Error(
			"Failed to auth trigger secret",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	if trigger.GetAction() == featureproto.FlagTrigger_Action_ON {
		err := h.enableFeature(ctx, trigger.GetFeatureId(), trigger.GetEnvironmentNamespace())
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if trigger.GetAction() == featureproto.FlagTrigger_Action_OFF {
		err := h.disableFeature(ctx, trigger.GetFeatureId(), trigger.GetEnvironmentNamespace())
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		h.logger.Error(
			"Invalid action",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.updateTriggerUsageInfo(ctx, trigger)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func (h handler) authSecret(
	ctx context.Context,
	secret string,
) (*domain.FlagTriggerSecret, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(secret)
	if err != nil {
		h.logger.Error(
			"Failed to decode encrypted secret",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	decrypted, err := h.triggerCryptoUtil.Decrypt(ctx, decoded)
	if err != nil {
		h.logger.Error(
			"Failed to decrypt encrypted secret",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	triggerSecret, err := domain.UnmarshalFlagTriggerSecret(decrypted)
	if err != nil {
		h.logger.Error(
			"Failed to unmarshal trigger url secret",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	return triggerSecret, nil
}

func (h handler) updateTriggerUsageInfo(ctx context.Context, flagTrigger *domain.FlagTrigger) error {
	tx, err := h.mysqlClient.BeginTx(ctx)
	if err != nil {
		h.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).
				AddFields(zap.Error(err))...,
		)
		return nil
	}
	err = h.mysqlClient.RunInTransaction(ctx, tx, func() error {
		storage := v2fs.NewFlagTriggerStorage(tx)
		err := storage.UpdateFlagTriggerUsage(
			ctx,
			flagTrigger.GetId(),
			flagTrigger.GetEnvironmentNamespace(),
			int64(flagTrigger.GetTriggerTimes()+1),
		)
		if err != nil {
			h.logger.Error(
				"Failed to update flag trigger usage",
				log.FieldsFromImcomingContext(ctx).
					AddFields(zap.Error(err))...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (h handler) enableFeature(
	ctx context.Context,
	featureId, environmentNamespace string,
) error {
	req := &featureproto.EnableFeatureRequest{
		Id:                   featureId,
		Command:              &featureproto.EnableFeatureCommand{},
		EnvironmentNamespace: environmentNamespace,
	}
	_, err := h.featureClient.EnableFeature(ctx, req)
	if err != nil {
		if code := status.Code(err); code == codes.FailedPrecondition {
			h.logger.Warn(
				"Feature flag is already enabled",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("featureId", req.Id),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return nil
		}
		h.logger.Error(
			"Failed to enable feature flag",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", req.Id),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return err
	}
	return nil
}

func (h handler) disableFeature(ctx context.Context,
	featureId, environmentNamespace string,
) error {
	req := &featureproto.DisableFeatureRequest{
		Id:                   featureId,
		Command:              &featureproto.DisableFeatureCommand{},
		EnvironmentNamespace: environmentNamespace,
	}
	_, err := h.featureClient.DisableFeature(ctx, req)
	if err != nil {
		if code := status.Code(err); code == codes.FailedPrecondition {
			h.logger.Warn(
				"Feature flag is already disabled",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("featureId", req.Id),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return nil
		}
		h.logger.Error(
			"Failed to disable feature flag",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", req.Id),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return err
	}
	return nil
}
