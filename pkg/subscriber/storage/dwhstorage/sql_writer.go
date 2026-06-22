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

package dwhstorage

import (
	"context"
	"encoding/json"

	epproto "github.com/bucketeer-io/bucketeer/v2/proto/eventpersisterdwh"
)

// sqlEvalEventWriter adapts an SQL EvaluationEventStorageV2 (mysql/postgres) to the
// EvalEventWriter contract: it converts proto events to params and performs a batch insert.
type sqlEvalEventWriter struct {
	storage EvaluationEventStorageV2
}

// NewEvalEventWriter returns an EvalEventWriter backed by an SQL evaluation-event storage.
func NewEvalEventWriter(storage EvaluationEventStorageV2) EvalEventWriter {
	return &sqlEvalEventWriter{storage: storage}
}

func (w *sqlEvalEventWriter) AppendRows(
	ctx context.Context,
	events []*epproto.EvaluationEvent,
) (map[string]bool, error) {
	fails := make(map[string]bool, len(events))

	// Early return if no events
	if len(events) == 0 {
		return fails, nil
	}

	// Prepare batch parameters
	batchEvents := make([]EvaluationEventParams, 0, len(events))
	for _, evt := range events {
		userData, err := json.Marshal(evt.UserData)
		if err != nil {
			fails[evt.Id] = true
			continue
		}

		batchEvents = append(batchEvents, EvaluationEventParams{
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
	}

	// Execute batch insert
	if len(batchEvents) > 0 {
		if err := w.storage.CreateEvaluationEvents(ctx, batchEvents); err != nil {
			// If batch fails, mark all events as failed
			for _, evt := range events {
				fails[evt.Id] = true
			}
		}
	}

	return fails, nil
}

// sqlGoalEventWriter adapts an SQL GoalEventStorageV2 (mysql/postgres) to the GoalEventWriter contract.
type sqlGoalEventWriter struct {
	storage GoalEventStorageV2
}

// NewGoalEventWriter returns a GoalEventWriter backed by an SQL goal-event storage.
func NewGoalEventWriter(storage GoalEventStorageV2) GoalEventWriter {
	return &sqlGoalEventWriter{storage: storage}
}

func (w *sqlGoalEventWriter) AppendRows(
	ctx context.Context,
	events []*epproto.GoalEvent,
) (map[string]bool, error) {
	fails := make(map[string]bool, len(events))

	// Early return if no events
	if len(events) == 0 {
		return fails, nil
	}

	// Prepare batch parameters
	batchEvents := make([]GoalEventParams, 0, len(events))
	for _, evt := range events {
		userData, err := json.Marshal(evt.UserData)
		if err != nil {
			fails[evt.Id] = true
			continue
		}

		batchEvents = append(batchEvents, GoalEventParams{
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
	}

	// Execute batch insert
	if len(batchEvents) > 0 {
		if err := w.storage.CreateGoalEvents(ctx, batchEvents); err != nil {
			// If batch fails, mark all events as failed
			for _, evt := range events {
				fails[evt.Id] = true
			}
		}
	}

	return fails, nil
}
