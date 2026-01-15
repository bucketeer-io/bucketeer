// Copyright 2026 The Bucketeer Authors.
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

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	ftdomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	ftstorage "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

func executeAutoOpsRuleOperation(
	ctx context.Context,
	ftStorage ftstorage.FeatureStorage,
	environmentId string,
	actionType autoopsproto.ActionType,
	feature *ftdomain.Feature,
	logger *zap.Logger,
	publisher publisher.Publisher,
	editor *eventproto.Editor,
) error {
	switch actionType {
	case autoopsproto.ActionType_ENABLE:
		return enableFeature(ctx, ftStorage, environmentId, feature, logger, publisher, editor)
	case autoopsproto.ActionType_DISABLE:
		return disableFeature(ctx, ftStorage, environmentId, feature, logger, publisher, editor)
	}
	return statusUnknownOpsType.Err()
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
	publisher publisher.Publisher,
	editor *eventproto.Editor,
) error {
	// Use domain layer Update method which handles validation and version incrementing
	updatedFeature, err := feature.Update(
		nil,                                // name
		nil,                                // description
		nil,                                // tags
		&wrapperspb.BoolValue{Value: true}, // enabled
		nil,                                // archived
		nil,                                // defaultStrategy
		nil,                                // offVariation
		false,                              // resetSamplingSeed
		nil,                                // prerequisiteChanges
		nil,                                // targetChanges
		nil,                                // ruleChanges
		nil,                                // variationChanges
		nil,                                // tagChanges
	)
	if err != nil {
		return err
	}

	// Only update storage if there were actual changes
	if updatedFeature.Version > feature.Version {
		// Update storage first
		if err := ftStorage.UpdateFeature(ctx, updatedFeature, environmentId); err != nil {
			return err
		}

		// Publish feature domain event
		publishFeatureEvent(
			ctx,
			publisher,
			editor,
			feature,
			environmentId,
			updatedFeature,
			eventproto.Event_FEATURE_ENABLED,
			&eventproto.FeatureEnabledEvent{Id: feature.Id},
			"Enabled by auto operation",
			logger,
		)
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
	publisher publisher.Publisher,
	editor *eventproto.Editor,
) error {
	// Use domain layer Update method which handles validation and version incrementing
	updatedFeature, err := feature.Update(
		nil,                                 // name
		nil,                                 // description
		nil,                                 // tags
		&wrapperspb.BoolValue{Value: false}, // enabled
		nil,                                 // archived
		nil,                                 // defaultStrategy
		nil,                                 // offVariation
		false,                               // resetSamplingSeed
		nil,                                 // prerequisiteChanges
		nil,                                 // targetChanges
		nil,                                 // ruleChanges
		nil,                                 // variationChanges
		nil,                                 // tagChanges
	)
	if err != nil {
		return err
	}

	// Only update storage if there were actual changes
	if updatedFeature.Version > feature.Version {
		// Update storage first
		if err := ftStorage.UpdateFeature(ctx, updatedFeature, environmentId); err != nil {
			return err
		}

		// Publish feature domain event
		publishFeatureEvent(
			ctx,
			publisher,
			editor,
			feature,
			environmentId,
			updatedFeature,
			eventproto.Event_FEATURE_DISABLED,
			&eventproto.FeatureDisabledEvent{Id: feature.Id},
			"Disabled by auto operation",
			logger,
		)
	}

	return nil
}

func publishFeatureEvent(
	ctx context.Context,
	publisher publisher.Publisher,
	editor *eventproto.Editor,
	feature *ftdomain.Feature,
	environmentId string,
	updatedFeature *ftdomain.Feature,
	eventType eventproto.Event_Type,
	eventData proto.Message,
	comment string,
	logger *zap.Logger,
) {
	featureEvent, err := domainevent.NewEvent(
		editor,
		eventproto.Event_FEATURE,
		feature.Id,
		eventType,
		eventData,
		environmentId,
		updatedFeature.Feature, // current state (after)
		feature.Feature,        // previous state (before)
		domainevent.WithComment(comment),
		domainevent.WithNewVersion(updatedFeature.Version),
	)
	if err != nil {
		logger.Error("Failed to create feature domain event", zap.Error(err))
		// Don't return error to avoid breaking auto ops execution
	} else {
		// Publish feature domain event
		if err := publisher.Publish(ctx, featureEvent); err != nil {
			logger.Error("Failed to publish feature domain event", zap.Error(err))
			// Don't return error to avoid breaking auto ops execution
		}
	}
}
