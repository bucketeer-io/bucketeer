// Copyright 2024 The Bucketeer Authors.
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

	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	ftdomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type Command interface{}

func executeAutoOpsRuleOperation(
	ctx context.Context,
	environmentNamespace string,
	autoOpsRule *domain.AutoOpsRule,
	feature *ftdomain.Feature,
	logger *zap.Logger,
	localizer locale.Localizer,
) error {
	switch autoOpsRule.OpsType {
	case autoopsproto.OpsType_ENABLE_FEATURE:
		return enableFeature(ctx, environmentNamespace, autoOpsRule, feature, logger)
	case autoopsproto.OpsType_DISABLE_FEATURE:
		return disableFeature(ctx, environmentNamespace, autoOpsRule, feature, logger)
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
	feature *ftdomain.Feature,
	logger *zap.Logger,
) error {
	req := &featureproto.EnableFeatureRequest{
		Id:                   autoOpsRule.FeatureId,
		Command:              &featureproto.EnableFeatureCommand{},
		EnvironmentNamespace: environmentNamespace,
	}
	if err := feature.Enable(); err != nil {
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
	feature *ftdomain.Feature,
	logger *zap.Logger,
) error {
	if err := feature.Disable(); err != nil {
		logger.Error(
			"Failed to disable feature flag",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", feature.Id),
				zap.String("environmentNamespace", environmentNamespace),
			)...,
		)
		return err
	}
	return nil
}
