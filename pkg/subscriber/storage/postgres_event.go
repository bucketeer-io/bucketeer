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

package storage

import (
	"context"
	"encoding/json"
	"fmt"

	storagev2 "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/v2"
	epproto "github.com/bucketeer-io/bucketeer/v2/proto/eventpersisterdwh"
)

type postgresEvalEventWriter struct {
	storage storagev2.EvaluationEventStorageV2
}

type postgresGoalEventWriter struct {
	storage storagev2.GoalEventStorageV2
}

func NewPostgresEvalEventWriter(storage storagev2.EvaluationEventStorageV2) EvalEventWriter {
	return &postgresEvalEventWriter{
		storage: storage,
	}
}

func NewPostgresGoalEventWriter(storage storagev2.GoalEventStorageV2) GoalEventWriter {
	return &postgresGoalEventWriter{
		storage: storage,
	}
}

func (w *postgresEvalEventWriter) AppendRows(
	ctx context.Context,
	events []*epproto.EvaluationEvent,
) (map[string]bool, error) {
	fails := make(map[string]bool, len(events))

	// Early return if no events
	if len(events) == 0 {
		return fails, nil
	}

	// Prepare batch parameters
	batchEvents := make([]storagev2.EvaluationEventParams, 0, len(events))
	idToIndex := make(map[string]int, len(events))

	for i, evt := range events {
		userData, err := json.Marshal(evt.UserData)
		if err != nil {
			fails[evt.Id] = true
			continue
		}

		batchEvents = append(batchEvents, storagev2.EvaluationEventParams{
			ID:             evt.Id,
			EnvironmentID:  evt.EnvironmentId,
			Timestamp:      evt.Timestamp,
			FeatureID:      evt.FeatureId,
			FeatureVersion: evt.FeatureVersion,
			UserID:         evt.UserId,
			UserData:       string(userData),
			VariationID:    evt.VariationId,
			Reason:         evt.Reason,
			Tag:            evt.Tag,
			SourceID:       evt.SourceId,
		})
		idToIndex[evt.Id] = i
	}

	// Execute batch insert
	if len(batchEvents) > 0 {
		err := w.storage.CreateEvaluationEvents(ctx, batchEvents)
		if err != nil {
			fmt.Printf("Error inserting evaluation events: %v\n", err)
			// If batch fails, mark all events as failed
			for _, evt := range events {
				fails[evt.Id] = true
			}
		}
	}

	return fails, nil
}

func (w *postgresGoalEventWriter) AppendRows(
	ctx context.Context,
	events []*epproto.GoalEvent,
) (map[string]bool, error) {
	fails := make(map[string]bool, len(events))

	// Early return if no events
	if len(events) == 0 {
		return fails, nil
	}

	// Prepare batch parameters
	batchEvents := make([]storagev2.GoalEventParams, 0, len(events))
	idToIndex := make(map[string]int, len(events))

	for i, evt := range events {
		userData, err := json.Marshal(evt.UserData)
		if err != nil {
			fails[evt.Id] = true
			continue
		}

		batchEvents = append(batchEvents, storagev2.GoalEventParams{
			ID:             evt.Id,
			EnvironmentID:  evt.EnvironmentId,
			Timestamp:      evt.Timestamp,
			GoalID:         evt.GoalId,
			Value:          evt.Value,
			UserID:         evt.UserId,
			UserData:       string(userData),
			Tag:            evt.Tag,
			SourceID:       evt.SourceId,
			FeatureID:      evt.FeatureId,
			FeatureVersion: evt.FeatureVersion,
			VariationID:    evt.VariationId,
			Reason:         evt.Reason,
		})
		idToIndex[evt.Id] = i
	}

	// Execute batch insert
	if len(batchEvents) > 0 {
		err := w.storage.CreateGoalEvents(ctx, batchEvents)
		if err != nil {
			// If batch fails, mark all events as failed
			for _, evt := range events {
				fails[evt.Id] = true
			}
		}
	}

	return fails, nil
}
