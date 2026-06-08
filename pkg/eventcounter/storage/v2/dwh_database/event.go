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
package dwhdatabase

import (
	"context"
	"time"

	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
)

const (
	DataTypeEvaluationEvent = "evaluation_event"
	DataTypeGoalEvent       = "goal_event"
)

var (
	ErrBQUnexpectedMultipleResults = pkgErr.NewErrorInternal(
		pkgErr.EventCounterPackageName,
		"bigquery: unexpected multiple results")
	ErrBQNoResultsFound = pkgErr.NewErrorInternal(
		pkgErr.EventCounterPackageName,
		"bigquery: no results found")

	ErrMySQLUnexpectedMultipleResults = pkgErr.NewErrorInternal(
		pkgErr.EventCounterPackageName,
		"MySQL: unexpected multiple results")
	ErrMySQLNoResultsFound = pkgErr.NewErrorInternal(
		pkgErr.EventCounterPackageName,
		"MySQL: no results found")

	ErrPostgresUnexpectedMultipleResults = pkgErr.NewErrorInternal(
		pkgErr.EventCounterPackageName,
		"Postgres: unexpected multiple results")
	ErrPostgresNoResultsFound = pkgErr.NewErrorInternal(
		pkgErr.EventCounterPackageName,
		"Postgres: no results found")
)

type EventStorage interface {
	QueryEvaluationCount(
		ctx context.Context,
		environmentId string,
		startAt, endAt time.Time,
		featureID string,
		featureVersion int32,
	) ([]*EvaluationEventCount, error)
	QueryGoalCount(
		ctx context.Context,
		environmentId string,
		startAt, endAt time.Time,
		goalID, featureID string,
		featureVersion int32,
	) ([]*GoalEventCount, error)
	QueryUserEvaluation(
		ctx context.Context,
		environmentID, userID, featureID string,
		featureVersion int32,
		experimentStartAt, experimentEndAt time.Time,
	) (*UserEvaluation, error)
}

type EvaluationEventCount struct {
	VariationID     string
	EvaluationUser  int64
	EvaluationTotal int64
}

type GoalEventCount struct {
	VariationID       string
	GoalUser          int64
	GoalTotal         int64
	GoalValueTotal    float64
	GoalValueMean     float64
	GoalValueVariance float64
}

type UserEvaluation struct {
	UserID         string
	FeatureID      string
	FeatureVersion int32
	VariationID    string
	Reason         string
	Timestamp      int64
}
