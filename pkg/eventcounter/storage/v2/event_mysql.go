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
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

const (
	DataTypeEvaluationEventMySQL = "evaluation_event"
	DataTypeGoalEventMySQL       = "goal_event"
)

var (
	//go:embed sql/evaluation_event_mysql.sql
	evaluationEventMySQLQuery string

	//go:embed sql/goal_event_mysql.sql
	goalEventMySQLQuery string

	//go:embed sql/user_evaluation_mysql.sql
	userEvaluationMySQLQuery string
)

type mysqlEventStorage struct {
	qe     mysql.QueryExecer
	logger *zap.Logger
}

func NewMySQLEventStorage(qe mysql.QueryExecer, logger *zap.Logger) EventStorage {
	return &mysqlEventStorage{
		qe:     qe,
		logger: logger.Named("mysql-event-storage"),
	}
}

func (es *mysqlEventStorage) QueryEvaluationCount(
	ctx context.Context,
	environmentId string,
	startAt, endAt time.Time,
	featureID string,
	featureVersion int32,
) ([]*EvaluationEventCount, error) {
	rows, err := es.qe.QueryContext(
		ctx,
		evaluationEventMySQLQuery,
		startAt,
		endAt,
		environmentId,
		featureID,
		featureVersion,
	)
	if err != nil {
		es.logger.Error(
			"Failed to query evaluation count",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
				zap.Time("startAt", startAt),
				zap.Time("endAt", endAt),
				zap.String("featureId", featureID),
				zap.Int32("featureVersion", featureVersion),
			)...,
		)
		return nil, err
	}
	defer rows.Close()

	results := make([]*EvaluationEventCount, 0)
	for rows.Next() {
		var ec EvaluationEventCount
		if err := rows.Scan(&ec.VariationID, &ec.EvaluationUser, &ec.EvaluationTotal); err != nil {
			es.logger.Error(
				"Failed to scan evaluation event count",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
				)...,
			)
			return nil, err
		}
		results = append(results, &ec)
	}

	if err := rows.Err(); err != nil {
		es.logger.Error(
			"Error after scanning evaluation event counts",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	return results, nil
}

func (es *mysqlEventStorage) QueryGoalCount(
	ctx context.Context,
	environmentId string,
	startAt, endAt time.Time,
	goalID, featureID string,
	featureVersion int32,
) ([]*GoalEventCount, error) {
	rows, err := es.qe.QueryContext(
		ctx,
		goalEventMySQLQuery,
		startAt,
		endAt,
		environmentId,
		goalID,
		featureID,
		featureVersion,
	)
	if err != nil {
		es.logger.Error(
			"Failed to query goal count",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
				zap.Time("startAt", startAt),
				zap.Time("endAt", endAt),
				zap.String("goalId", goalID),
				zap.String("featureId", featureID),
				zap.Int32("featureVersion", featureVersion),
			)...,
		)
		return nil, err
	}
	defer rows.Close()

	results := make([]*GoalEventCount, 0)
	for rows.Next() {
		var gc GoalEventCount
		if err := rows.Scan(
			&gc.VariationID,
			&gc.GoalUser,
			&gc.GoalTotal,
			&gc.GoalValueTotal,
			&gc.GoalValueMean,
			&gc.GoalValueVariance,
		); err != nil {
			es.logger.Error(
				"Failed to scan goal event count",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
				)...,
			)
			return nil, err
		}
		results = append(results, &gc)
	}

	if err := rows.Err(); err != nil {
		es.logger.Error(
			"Error after scanning goal event counts",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	return results, nil
}

func (es *mysqlEventStorage) QueryUserEvaluation(
	ctx context.Context,
	environmentID, userID, featureID string,
	featureVersion int32,
	experimentStartAt, experimentEndAt time.Time,
) (*UserEvaluation, error) {
	rows, err := es.qe.QueryContext(
		ctx,
		userEvaluationMySQLQuery,
		environmentID,
		featureID,
		featureVersion,
		userID,
		experimentStartAt,
		experimentEndAt,
	)
	if err != nil {
		es.logger.Error(
			"Failed to query user evaluation",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentID),
				zap.String("userId", userID),
				zap.String("featureId", featureID),
				zap.Int32("featureVersion", featureVersion),
				zap.Time("experimentStartAt", experimentStartAt),
				zap.Time("experimentEndAt", experimentEndAt),
			)...,
		)
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, ErrNoResultsFound
	}

	var ue UserEvaluation
	if err := rows.Scan(
		&ue.UserID,
		&ue.FeatureID,
		&ue.FeatureVersion,
		&ue.VariationID,
		&ue.Reason,
		&ue.Timestamp,
	); err != nil {
		es.logger.Error(
			"Failed to scan user evaluation",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	if rows.Next() {
		return nil, ErrUnexpectedMultipleResults
	}

	if err := rows.Err(); err != nil {
		es.logger.Error(
			"Error after scanning user evaluation",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	return &ue, nil
}
