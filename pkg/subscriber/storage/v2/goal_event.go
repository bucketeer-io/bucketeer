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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v2

import (
	"context"
	_ "embed"
	"strings"

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
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
		query.WriteString("(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(
			args,
			event.ID,
			event.EnvironmentID,
			event.Timestamp,
			event.GoalID,
			event.Value,
			event.UserID,
			event.UserData,
			event.Tag,
			event.SourceID,
			event.FeatureID,
			event.FeatureVersion,
			event.VariationID,
			event.Reason,
		)
	}

	_, err := s.qe.ExecContext(ctx, query.String(), args...)
	return err
}
