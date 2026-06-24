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
package dwhstorage

import (
	"context"

	epproto "github.com/bucketeer-io/bucketeer/v2/proto/eventpersisterdwh"
)

// EvalEventWriter writes evaluation events to the data warehouse.
// Implemented per data-warehouse backend (bigquery, mysql, postgres).
type EvalEventWriter interface {
	AppendRows(ctx context.Context, events []*epproto.EvaluationEvent) (map[string]bool, error)
}

// GoalEventWriter writes goal events to the data warehouse.
// Implemented per data-warehouse backend (bigquery, mysql, postgres).
type GoalEventWriter interface {
	AppendRows(ctx context.Context, events []*epproto.GoalEvent) (map[string]bool, error)
}

// EvaluationEventStorageV2 is the SQL (mysql/postgres) insert contract for evaluation events.
type EvaluationEventStorageV2 interface {
	CreateEvaluationEvents(ctx context.Context, events []EvaluationEventParams) error
}

// GoalEventStorageV2 is the SQL (mysql/postgres) insert contract for goal events.
type GoalEventStorageV2 interface {
	CreateGoalEvents(ctx context.Context, events []GoalEventParams) error
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
