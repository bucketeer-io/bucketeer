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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v2

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"

	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
)

var (
	//go:embed sql/goal_event.sql
	goalEventSql string
)

const (
	goalBatchSize = 1000
)

type GoalEventStorageV2 interface {
	CreateGoalEvents(ctx context.Context, events []GoalEventParams) error
}

type GoalEventParams struct {
	ID             string
	EnvironmentID  string
	Timestamp      int64
	GoalID         string
	Value          float32
	UserID         string
	UserData       string
	Tag            string
	SourceID       string
	FeatureID      string
	FeatureVersion int32
	VariationID    string
	Reason         string
}

type mysqlGoalEventStorage struct {
	qe mysql.QueryExecer
}

func NewMysqlGoalEventStorage(qe mysql.QueryExecer) GoalEventStorageV2 {
	return &mysqlGoalEventStorage{qe: qe}
}

func (s *mysqlGoalEventStorage) CreateGoalEvents(
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

func (s *mysqlGoalEventStorage) createGoalEventsBatch(
	ctx context.Context,
	events []GoalEventParams,
) error {
	var query strings.Builder
	query.WriteString(goalEventSql)

	args := make([]interface{}, 0, len(events)*13) // 13 fields per event

	for i, event := range events {
		if i > 0 {
			query.WriteString(",")
		}
		query.WriteString("(?, ?, FROM_UNIXTIME(?), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

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
	return err
}
