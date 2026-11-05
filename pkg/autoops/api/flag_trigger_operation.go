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

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	ftdomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	ftstorage "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

// executeAutoOpsRuleOperation applies the action to the feature inside the
// caller's transaction and returns the feature domain event describing the
// change. The event is nil when the feature was already in the requested
// state. The caller must publish the event after the transaction commits.
func executeAutoOpsRuleOperation(
	ctx context.Context,
	ftStorage ftstorage.FeatureStorage,
	environmentId string,
	actionType autoopsproto.ActionType,
	feature *ftdomain.Feature,
	editor *eventproto.Editor,
) (*eventproto.Event, error) {
	switch actionType {
	case autoopsproto.ActionType_ENABLE:
		return enableFeature(ctx, ftStorage, environmentId, feature, editor)
	case autoopsproto.ActionType_DISABLE:
		return disableFeature(ctx, ftStorage, environmentId, feature, editor)
	}
	return nil, statusUnknownOpsType.Err()
}

func enableFeature(
	ctx context.Context,
	ftStorage ftstorage.FeatureStorage,
	environmentId string,
	feature *ftdomain.Feature,
	editor *eventproto.Editor,
) (*eventproto.Event, error) {
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
		nil,                                // maintainer
		nil,                                // ruleOrder
		nil,                                // variationValueSchemaUpdate
	)
	if err != nil {
		return nil, err
	}

	// Only update storage if there were actual changes
	if updatedFeature.Version <= feature.Version {
		return nil, nil
	}
	if err := ftStorage.UpdateFeature(ctx, updatedFeature, environmentId); err != nil {
		return nil, err
	}
	return newFeatureEvent(
		editor,
		feature,
		environmentId,
		updatedFeature,
		eventproto.Event_FEATURE_ENABLED,
		&eventproto.FeatureEnabledEvent{Id: feature.Id},
		"Enabled by auto operation",
	)
}

func disableFeature(
	ctx context.Context,
	ftStorage ftstorage.FeatureStorage,
	environmentId string,
	feature *ftdomain.Feature,
	editor *eventproto.Editor,
) (*eventproto.Event, error) {
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
		nil,                                 // maintainer
		nil,                                 // ruleOrder
		nil,                                 // variationValueSchemaUpdate
	)
	if err != nil {
		return nil, err
	}

	// Only update storage if there were actual changes
	if updatedFeature.Version <= feature.Version {
		return nil, nil
	}
	if err := ftStorage.UpdateFeature(ctx, updatedFeature, environmentId); err != nil {
		return nil, err
	}
	return newFeatureEvent(
		editor,
		feature,
		environmentId,
		updatedFeature,
		eventproto.Event_FEATURE_DISABLED,
		&eventproto.FeatureDisabledEvent{Id: feature.Id},
		"Disabled by auto operation",
	)
}

func newFeatureEvent(
	editor *eventproto.Editor,
	feature *ftdomain.Feature,
	environmentId string,
	updatedFeature *ftdomain.Feature,
	eventType eventproto.Event_Type,
	eventData proto.Message,
	comment string,
) (*eventproto.Event, error) {
	return domainevent.NewEvent(
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
}
