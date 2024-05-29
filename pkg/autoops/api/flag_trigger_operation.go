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

	ftdomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	ftstorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

func executeAutoOpsRuleOperation(
	ctx context.Context,
	ftStorage ftstorage.FeatureStorage,
	environmentNamespace string,
	actionType autoopsproto.ActionType,
	feature *ftdomain.Feature,
	logger *zap.Logger,
	localizer locale.Localizer,
) error {
	logger.Debug("Start Execute auto ops rule operation actionType",
		zap.Int("actionType", int(actionType)))
	switch actionType {
	case autoopsproto.ActionType_ENABLE:
		return enableFeature(ctx, ftStorage, environmentNamespace, feature, logger)
	case autoopsproto.ActionType_DISABLE:
		return disableFeature(ctx, ftStorage, environmentNamespace, feature, logger)
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

// If the flag is already enabled, we only print an error log.
// Otherwise, the Operation state won't change to TRIGGERED, keeping in an infinite loop
// trying to enable the flag.
func enableFeature(
	ctx context.Context,
	ftStorage ftstorage.FeatureStorage,
	environmentNamespace string,
	feature *ftdomain.Feature,
	logger *zap.Logger,
) error {
	if err := feature.Enable(); err != nil {
		logger.Error(
			"Failed to enable feature flag",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", feature.Id),
				zap.String("environmentNamespace", environmentNamespace),
			)...,
		)
		return nil
	}
	if err := ftStorage.UpdateFeature(ctx, feature, environmentNamespace); err != nil {
		return err
	}
	return nil
}

// If the flag is already disabled, we only print an error log.
// Otherwise, the Operation state won't change to TRIGGERED, keeping in an infinite loop
// trying to disable the flag.
func disableFeature(
	ctx context.Context,
	ftStorage ftstorage.FeatureStorage,
	environmentNamespace string,
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
		return nil
	}
	if err := ftStorage.UpdateFeature(ctx, feature, environmentNamespace); err != nil {
		return err
	}
	return nil
}
