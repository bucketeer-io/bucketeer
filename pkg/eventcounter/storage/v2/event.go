// Copyright 2024 The Bucketeer Authors.
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
	"embed"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	bqquerier "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/querier"
)

const (
	DataTypeEvaluationEvent = "evaluation_event"
	DataTypeGoalEvent       = "goal_event"
	EvaluationCountSQLFile  = "sql/evaluation_count.sql"
	GoalCountSQLFile        = "sql/goal_count.sql"
)

var (
	//go:embed sql
	sql embed.FS
)

type EventStorage interface {
	QueryEvaluationCount(
		ctx context.Context,
		environmentNamespace string,
		startAt, endAt time.Time,
		featureID string,
		featureVersion int32,
	) ([]*EvaluationEventCount, error)
	QueryGoalCount(
		ctx context.Context,
		environmentNamespace string,
		startAt, endAt time.Time,
		goalID, featureID string,
		featureVersion int32,
	) ([]*GoalEventCount, error)
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

func NewEventStorage(querier bqquerier.Client, dataset string, logger *zap.Logger) EventStorage {
	return &eventStorage{
		querier: querier,
		dataset: dataset,
		logger:  logger.Named("storage"),
	}
}

func (es *eventStorage) QueryEvaluationCount(
	ctx context.Context,
	environmentNamespace string,
	startAt, endAt time.Time,
	featureID string,
	featureVersion int32,
) ([]*EvaluationEventCount, error) {
	fileName := EvaluationCountSQLFile
	q, err := sql.ReadFile(fileName)
	if err != nil {
		es.logger.Error(
			"Failed to read file",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("fileName", fileName),
			)...,
		)
		return nil, err
	}
	datasource := fmt.Sprintf("%s.%s", es.dataset, DataTypeEvaluationEvent)
	query := fmt.Sprintf(string(q), datasource)
	params := []bigquery.QueryParameter{
		{
			Name:  "environmentNamespace",
			Value: environmentNamespace,
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
			log.FieldsFromImcomingContext(ctx).AddFields(
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
				log.FieldsFromImcomingContext(ctx).AddFields(
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
	environmentNamespace string,
	startAt, endAt time.Time,
	goalID, featureID string,
	featureVersion int32,
) ([]*GoalEventCount, error) {
	fileName := GoalCountSQLFile
	q, err := sql.ReadFile(fileName)
	if err != nil {
		es.logger.Error(
			"Failed to read file",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("fileName", fileName),
			)...,
		)
		return nil, err
	}
	datasource := fmt.Sprintf("%s.%s", es.dataset, DataTypeGoalEvent)
	query := fmt.Sprintf(string(q), datasource)
	params := []bigquery.QueryParameter{
		{
			Name:  "environmentNamespace",
			Value: environmentNamespace,
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
			log.FieldsFromImcomingContext(ctx).AddFields(
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
				log.FieldsFromImcomingContext(ctx).AddFields(
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
