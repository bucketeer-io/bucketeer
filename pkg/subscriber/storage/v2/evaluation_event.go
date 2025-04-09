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
	//go:embed sql/evaluation_event.sql
	evaluationEventSql string
)

const (
	evaluationBatchSize = 1000
)

type EvaluationEventStorageV2 interface {
	CreateEvaluationEvents(ctx context.Context, events []EvaluationEventParams) error
}

type EvaluationEventParams struct {
	ID             string
	EnvironmentID  string
	Timestamp      int64
	FeatureID      string
	FeatureVersion int32
	UserID         string
	UserData       string
	VariationID    string
	Reason         string
	Tag            string
	SourceID       string
}

type mysqlEvaluationEventStorage struct {
	qe mysql.QueryExecer
}

func NewMysqlEvaluationEventStorage(qe mysql.QueryExecer) EvaluationEventStorageV2 {
	return &mysqlEvaluationEventStorage{qe: qe}
}

func (s *mysqlEvaluationEventStorage) CreateEvaluationEvents(
	ctx context.Context,
	events []EvaluationEventParams,
) error {
	// Process in batches to avoid exceeding max SQL parameter limits
	for i := 0; i < len(events); i += evaluationBatchSize {
		j := i + evaluationBatchSize
		if j > len(events) {
			j = len(events)
		}
		if err := s.createEvaluationEventsBatch(ctx, events[i:j]); err != nil {
			return err
		}
	}
	return nil
}

func (s *mysqlEvaluationEventStorage) createEvaluationEventsBatch(
	ctx context.Context,
	events []EvaluationEventParams,
) error {
	var query strings.Builder
	query.WriteString(evaluationEventSql)

	args := make([]interface{}, 0, len(events)*11) // 11 fields per event

	for i, event := range events {
		if i > 0 {
			query.WriteString(",")
		}
		query.WriteString("(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(
			args,
			event.ID,
			event.EnvironmentID,
			event.Timestamp,
			event.FeatureID,
			event.FeatureVersion,
			event.UserID,
			event.UserData,
			event.VariationID,
			event.Reason,
			event.Tag,
			event.SourceID,
		)
	}

	_, err := s.qe.ExecContext(ctx, query.String(), args...)
	return err
}
