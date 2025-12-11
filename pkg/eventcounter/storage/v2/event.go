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
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"

	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	bqquerier "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/bigquery/querier"
)

const (
	DataTypeEvaluationEvent = "evaluation_event"
	DataTypeGoalEvent       = "goal_event"
)

var (
	//go:embed sql/select_user_evaluation.sql
	userEvaluationSQL string

	//go:embed sql/evaluation_count.sql
	evaluationCountSQL string

	//go:embed sql/goal_count.sql
	goalCountSQL string

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

type eventStorage struct {
	querier bqquerier.Client
	dataset string
	logger  *zap.Logger
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

func NewEventStorage(querier bqquerier.Client, dataset string, logger *zap.Logger) EventStorage {
	return &eventStorage{
		querier: querier,
		dataset: dataset,
		logger:  logger.Named("storage"),
	}
}

func (es *eventStorage) QueryEvaluationCount(
	ctx context.Context,
	environmentId string,
	startAt, endAt time.Time,
	featureID string,
	featureVersion int32,
) ([]*EvaluationEventCount, error) {
	datasource := fmt.Sprintf("%s.%s", es.dataset, DataTypeEvaluationEvent)
	query := fmt.Sprintf(evaluationCountSQL, datasource)
	params := []bigquery.QueryParameter{
		{
			Name:  "environmentId",
			Value: environmentId,
		},
		{
			Name:  "startAt",
			Value: startAt,
		},
		{
			Name:  "endAt",
			Value: endAt,
		},
		{
			Name:  "featureID",
			Value: featureID,
		},
		{
			Name:  "featureVersion",
			Value: featureVersion,
		},
	}
	es.logger.Debug("Query evaluation count",
		zap.String("query", query),
		zap.Any("params", params),
	)
	iter, err := es.querier.ExecQuery(ctx, query, params)
	if err != nil {
		es.logger.Error(
			"Failed to query evaluation count",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("query", query),
				zap.Any("params", params),
			)...,
		)
		return nil, err
	}
	rows := make([]*EvaluationEventCount, 0, iter.TotalRows)
	for {
		var row EvaluationEventCount
		err := iter.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			es.logger.Error(
				"Failed to convert evaluation event count from the query result",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("query", query),
					zap.Any("params", params),
				)...,
			)
			return nil, err
		}
		rows = append(rows, &row)
	}
	return rows, nil
}

func (es *eventStorage) QueryGoalCount(
	ctx context.Context,
	environmentId string,
	startAt, endAt time.Time,
	goalID, featureID string,
	featureVersion int32,
) ([]*GoalEventCount, error) {
	datasource := fmt.Sprintf("%s.%s", es.dataset, DataTypeGoalEvent)
	query := fmt.Sprintf(goalCountSQL, datasource)
	params := []bigquery.QueryParameter{
		{
			Name:  "environmentId",
			Value: environmentId,
		},
		{
			Name:  "startAt",
			Value: startAt,
		},
		{
			Name:  "endAt",
			Value: endAt,
		},
		{
			Name:  "goalID",
			Value: goalID,
		},
		{
			Name:  "featureID",
			Value: featureID,
		},
		{
			Name:  "featureVersion",
			Value: featureVersion,
		},
	}
	es.logger.Debug("query goal count",
		zap.String("query", query),
		zap.Any("params", params),
	)
	iter, err := es.querier.ExecQuery(ctx, query, params)
	if err != nil {
		es.logger.Error(
			"Failed to query goal count",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("query", query),
				zap.Any("params", params),
			)...,
		)
		return nil, err
	}
	rows := make([]*GoalEventCount, 0, iter.TotalRows)
	for {
		var row GoalEventCount
		err := iter.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			es.logger.Error(
				"Failed to convert goal event count from the query result",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("query", query),
					zap.Any("params", params),
				)...,
			)
			return nil, err
		}
		rows = append(rows, &row)
	}
	return rows, nil
}

func (es *eventStorage) QueryUserEvaluation(
	ctx context.Context,
	environmentID, userID, featureID string,
	featureVersion int32,
	experimentStartAt, experimentEndAt time.Time,
) (*UserEvaluation, error) {
	datasource := fmt.Sprintf("%s.%s", es.dataset, DataTypeEvaluationEvent)
	query := fmt.Sprintf(userEvaluationSQL, datasource)
	params := []bigquery.QueryParameter{
		{
			Name:  "environmentId",
			Value: environmentID,
		},
		{
			Name:  "userId",
			Value: userID,
		},
		{
			Name:  "featureId",
			Value: featureID,
		},
		{
			Name:  "featureVersion",
			Value: featureVersion,
		},
		{
			Name:  "experimentStartAt",
			Value: experimentStartAt,
		},
		{
			Name:  "experimentEndAt",
			Value: experimentEndAt,
		},
	}
	es.logger.Debug("Query user evaluation",
		zap.String("query", query),
		zap.Any("params", params),
	)
	iter, err := es.querier.ExecQuery(ctx, query, params)
	if err != nil {
		es.logger.Error(
			"Failed to query user evaluation",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("query", query),
				zap.Any("params", params),
			)...,
		)
		return nil, err
	}

	// Check if there are unexpected multiple rows
	if iter.TotalRows > 1 {
		return nil, ErrBQUnexpectedMultipleResults
	}

	// Retrieve the single expected row
	var row UserEvaluation
	err = iter.Next(&row)
	if err == iterator.Done {
		es.logger.Error(
			"User evaluation not found",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("query", query),
				zap.Any("params", params),
			)...,
		)
		return nil, ErrBQNoResultsFound
	}
	if err != nil {
		es.logger.Error(
			"Failed to convert user evaluation from the query result",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("query", query),
				zap.Any("params", params),
			)...,
		)
		return nil, err
	}
	return &row, nil
}
