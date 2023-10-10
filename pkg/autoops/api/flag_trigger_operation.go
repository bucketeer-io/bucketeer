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

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type Command interface{}

func ExecuteAutoOpsRuleOperation(
	ctx context.Context,
	environmentNamespace string,
	autoOpsRule *domain.AutoOpsRule,
	featureClient featureclient.Client,
	logger *zap.Logger,
	localizer locale.Localizer,
) error {
	switch autoOpsRule.OpsType {
	case autoopsproto.OpsType_ENABLE_FEATURE:
		return enableFeature(ctx, environmentNamespace, autoOpsRule, featureClient, logger)
	case autoopsproto.OpsType_DISABLE_FEATURE:
		return disableFeature(ctx, environmentNamespace, autoOpsRule, featureClient, logger)
	}
	dt, err := statusUnknownOpsType.WithDetails(&errdetails.LocalizedMessage{
		Locale:  localizer.GetLocale(),
		Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "ops_type"),
	})
	if err != nil {
		return statusInternal.Err()
	}
	return dt.Err()
}

func enableFeature(
	ctx context.Context,
	environmentNamespace string,
	autoOpsRule *domain.AutoOpsRule,
	featureClient featureclient.Client,
	logger *zap.Logger,
) error {
	req := &featureproto.EnableFeatureRequest{
		Id:                   autoOpsRule.FeatureId,
		Command:              &featureproto.EnableFeatureCommand{},
		EnvironmentNamespace: environmentNamespace,
	}
	_, err := featureClient.EnableFeature(ctx, req)
	if err != nil {
		if code := status.Code(err); code == codes.FailedPrecondition {
			logger.Warn(
				"Feature flag is already enabled",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("featureId", req.Id),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return nil
		}
		logger.Error(
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

func disableFeature(
	ctx context.Context,
	environmentNamespace string,
	autoOpsRule *domain.AutoOpsRule,
	featureClient featureclient.Client,
	logger *zap.Logger,
) error {
	req := &featureproto.DisableFeatureRequest{
		Id:                   autoOpsRule.FeatureId,
		Command:              &featureproto.DisableFeatureCommand{},
		EnvironmentNamespace: environmentNamespace,
	}
	_, err := featureClient.DisableFeature(ctx, req)
	if err != nil {
		if code := status.Code(err); code == codes.FailedPrecondition {
			logger.Warn(
				"Feature flag is already disabled",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("featureId", req.Id),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return nil
		}
		logger.Error(
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
