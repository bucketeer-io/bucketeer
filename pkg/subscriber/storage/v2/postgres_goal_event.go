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

package v2

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
)

type postgresGoalEventStorage struct {
	qe postgres.QueryExecer
}

func NewPostgresGoalEventStorage(qe postgres.QueryExecer) GoalEventStorageV2 {
	return &postgresGoalEventStorage{qe: qe}
}

func (s *postgresGoalEventStorage) CreateGoalEvents(
	ctx context.Context,
	events []GoalEventParams,
) error {
	// Process in batches to avoid exceeding max SQL parameter limits
	for i := 0; i < len(events); i += goalBatchSize {
		j := i + goalBatchSize
		if j > len(events) {
			j = len(events)
		}
		if err := s.createGoalEventsBatch(ctx, events[i:j]); err != nil {
			return err
		}
	}
	return nil
}

func (s *postgresGoalEventStorage) createGoalEventsBatch(
	ctx context.Context,
	events []GoalEventParams,
) error {
	var query strings.Builder
	query.WriteString(goalEventSql)

	fieldsPerEvent := 13
	args := make([]interface{}, 0, len(events)*fieldsPerEvent) // 13 fields per event

	argCounter := 0
	for i, event := range events {
		if i > 0 {
			query.WriteString(",")
		}
		query.WriteString(postgres.WritePlaceHolder(
			"($%d, $%d, TO_TIMESTAMP($%d), $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			argCounter+1,
			fieldsPerEvent,
		))
		argCounter += fieldsPerEvent

		// Validate required fields
		if event.ID == "" || event.EnvironmentID == "" || event.GoalID == "" || event.UserID == "" {
			return fmt.Errorf("missing required fields: id=%s, envId=%s, goalId=%s, userId=%s",
				event.ID, event.EnvironmentID, event.GoalID, event.UserID)
		}

		// Handle potentially null fields
		userData := sql.NullString{String: event.UserData, Valid: event.UserData != ""}
		tag := sql.NullString{String: event.Tag, Valid: event.Tag != ""}
		sourceID := sql.NullString{String: event.SourceID, Valid: event.SourceID != ""}
		featureID := sql.NullString{String: event.FeatureID, Valid: event.FeatureID != ""}
		variationID := sql.NullString{String: event.VariationID, Valid: event.VariationID != ""}
		reason := sql.NullString{String: event.Reason, Valid: event.Reason != ""}

		// Convert timestamp from microseconds to seconds
		timestampSeconds := event.Timestamp / 1000000

		args = append(
			args,
			event.ID,
			event.EnvironmentID,
			timestampSeconds,
			event.GoalID,
			event.Value,
			event.UserID,
			userData,
			tag,
			sourceID,
			featureID,
			event.FeatureVersion,
			variationID,
			reason,
		)
	}

	_, err := s.qe.ExecContext(ctx, query.String(), args...)
	if err != nil {
		return fmt.Errorf("failed to execute batch insert: %w", err)
	}
	return nil
}
