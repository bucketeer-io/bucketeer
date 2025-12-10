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
