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

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	ftdomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	ftstorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

func executeAutoOpsRuleOperation(
	ctx context.Context,
	ftStorage ftstorage.FeatureStorage,
	environmentId string,
	actionType autoopsproto.ActionType,
	feature *ftdomain.Feature,
	logger *zap.Logger,
	localizer locale.Localizer,
) error {
	switch actionType {
	case autoopsproto.ActionType_ENABLE:
		return enableFeature(ctx, ftStorage, environmentId, feature, logger)
	case autoopsproto.ActionType_DISABLE:
		return disableFeature(ctx, ftStorage, environmentId, feature, logger)
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
	environmentId string,
	feature *ftdomain.Feature,
	logger *zap.Logger,
) error {
	if err := feature.Enable(); err != nil {
		// If the flag is already disabled, we skip the updating
		return nil
	}
	if err := ftStorage.UpdateFeature(ctx, feature, environmentId); err != nil {
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
	environmentId string,
	feature *ftdomain.Feature,
	logger *zap.Logger,
) error {
	if err := feature.Disable(); err != nil {
		// If the flag is already disabled, we skip the updating
		return nil
	}
	if err := ftStorage.UpdateFeature(ctx, feature, environmentId); err != nil {
		return err
	}
	return nil
}
