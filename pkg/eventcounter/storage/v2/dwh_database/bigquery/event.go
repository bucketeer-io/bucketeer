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

package bigquery

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"

	dwhdatabase "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/storage/v2/dwh_database"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	bqquerier "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/bigquery/querier"
)

var (
	//go:embed sql/select_user_evaluation.sql
	userEvaluationSQL string

	//go:embed sql/evaluation_count.sql
	evaluationCountSQL string

	//go:embed sql/goal_count.sql
	goalCountSQL string
)

type eventStorage struct {
	querier bqquerier.Client
	dataset string
	logger  *zap.Logger
}

func NewEventStorage(querier bqquerier.Client, dataset string, logger *zap.Logger) dwhdatabase.EventStorage {
	return &eventStorage{
		querier: querier,
		dataset: dataset,
		logger:  logger.Named("bigquery-event-storage"),
	}
}

func (es *eventStorage) QueryEvaluationCount(
	ctx context.Context,
	environmentId string,
	startAt, endAt time.Time,
	featureID string,
	featureVersion int32,
) ([]*dwhdatabase.EvaluationEventCount, error) {
	datasource := fmt.Sprintf("%s.%s", es.dataset, dwhdatabase.DataTypeEvaluationEvent)
	query := fmt.Sprintf(evaluationCountSQL, datasource)
	params := []bigquery.QueryParameter{
		{Name: "environmentId", Value: environmentId},
		{Name: "startAt", Value: startAt},
		{Name: "endAt", Value: endAt},
		{Name: "featureID", Value: featureID},
		{Name: "featureVersion", Value: featureVersion},
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
	rows := make([]*dwhdatabase.EvaluationEventCount, 0, iter.TotalRows)
	for {
		var row dwhdatabase.EvaluationEventCount
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
) ([]*dwhdatabase.GoalEventCount, error) {
	datasource := fmt.Sprintf("%s.%s", es.dataset, dwhdatabase.DataTypeGoalEvent)
	query := fmt.Sprintf(goalCountSQL, datasource)
	params := []bigquery.QueryParameter{
		{Name: "environmentId", Value: environmentId},
		{Name: "startAt", Value: startAt},
		{Name: "endAt", Value: endAt},
		{Name: "goalID", Value: goalID},
		{Name: "featureID", Value: featureID},
		{Name: "featureVersion", Value: featureVersion},
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
	rows := make([]*dwhdatabase.GoalEventCount, 0, iter.TotalRows)
	for {
		var row dwhdatabase.GoalEventCount
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
) (*dwhdatabase.UserEvaluation, error) {
	datasource := fmt.Sprintf("%s.%s", es.dataset, dwhdatabase.DataTypeEvaluationEvent)
	query := fmt.Sprintf(userEvaluationSQL, datasource)
	params := []bigquery.QueryParameter{
		{Name: "environmentId", Value: environmentID},
		{Name: "userId", Value: userID},
		{Name: "featureId", Value: featureID},
		{Name: "featureVersion", Value: featureVersion},
		{Name: "experimentStartAt", Value: experimentStartAt},
		{Name: "experimentEndAt", Value: experimentEndAt},
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

	if iter.TotalRows > 1 {
		return nil, dwhdatabase.ErrBQUnexpectedMultipleResults
	}

	var row dwhdatabase.UserEvaluation
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
		return nil, dwhdatabase.ErrBQNoResultsFound
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
