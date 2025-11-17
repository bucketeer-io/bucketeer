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
	"database/sql"
	_ "embed"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
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
	es.logger.Debug("🔍 MYSQL: Querying for user evaluation",
		zap.String("environmentId", environmentID),
		zap.String("userId", userID),
		zap.String("featureId", featureID),
		zap.Int32("featureVersion", featureVersion),
		zap.Time("experimentStartAt", experimentStartAt),
		zap.Time("experimentEndAt", experimentEndAt),
		zap.String("experimentStartAtISO", experimentStartAt.Format(time.RFC3339)),
		zap.String("experimentEndAtISO", experimentEndAt.Format(time.RFC3339)),
	)
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
		// Query to check what evaluation events exist for this user/feature (any version)
		checkRows, checkErr := es.qe.QueryContext(
			ctx,
			"SELECT feature_version, UNIX_TIMESTAMP(timestamp) as ts, variation_id, user_id FROM evaluation_event WHERE environment_id = ? AND user_id = ? AND feature_id = ? AND timestamp BETWEEN ? AND ? ORDER BY timestamp DESC LIMIT 10",
			environmentID,
			userID,
			featureID,
			experimentStartAt,
			experimentEndAt,
		)
		if checkErr == nil && checkRows != nil {
			var foundVersions []int32
			var foundTimestamps []int64
			var foundUserIDs []string
			for checkRows.Next() {
				var fv int32
				var ts int64
				var vid sql.NullString
				var uid string
				if err := checkRows.Scan(&fv, &ts, &vid, &uid); err == nil {
					foundVersions = append(foundVersions, fv)
					foundTimestamps = append(foundTimestamps, ts)
					foundUserIDs = append(foundUserIDs, uid)
				}
			}
			checkRows.Close()
			
			// Also check if there are ANY evaluation events for this feature (any user) in the time window
			anyRows, anyErr := es.qe.QueryContext(
				ctx,
				"SELECT user_id, feature_version, UNIX_TIMESTAMP(timestamp) as ts FROM evaluation_event WHERE environment_id = ? AND feature_id = ? AND timestamp BETWEEN ? AND ? ORDER BY timestamp DESC LIMIT 10",
				environmentID,
				featureID,
				experimentStartAt,
				experimentEndAt,
			)
			var anyUserIDs []string
			var anyVersions []int32
			var anyTimestamps []int64
			if anyErr == nil && anyRows != nil {
				for anyRows.Next() {
					var uid string
					var fv int32
					var ts int64
					if err := anyRows.Scan(&uid, &fv, &ts); err == nil {
						anyUserIDs = append(anyUserIDs, uid)
						anyVersions = append(anyVersions, fv)
						anyTimestamps = append(anyTimestamps, ts)
					}
				}
				anyRows.Close()
			}
			
			es.logger.Warn("MYSQL: No evaluation found with exact featureVersion",
				zap.String("environmentId", environmentID),
				zap.String("requestedUserId", userID),
				zap.String("featureId", featureID),
				zap.Int32("requestedFeatureVersion", featureVersion),
				zap.Int64s("foundTimestamps", foundTimestamps),
				zap.Int32s("foundFeatureVersions", foundVersions),
				zap.Strings("foundUserIDs", foundUserIDs),
				zap.Time("experimentStartAt", experimentStartAt),
				zap.Time("experimentEndAt", experimentEndAt),
				zap.String("experimentStartAtISO", experimentStartAt.Format(time.RFC3339)),
				zap.String("experimentEndAtISO", experimentEndAt.Format(time.RFC3339)),
				zap.Strings("anyUserIDsForFeature", anyUserIDs),
				zap.Int32s("anyVersionsForFeature", anyVersions),
				zap.Int64s("anyTimestampsForFeature", anyTimestamps),
			)
		} else {
			es.logger.Warn("MYSQL: No evaluation found in database - check query failed",
				zap.String("environmentId", environmentID),
				zap.String("userId", userID),
				zap.String("featureId", featureID),
				zap.Int32("featureVersion", featureVersion),
				zap.Error(checkErr),
				zap.Time("experimentStartAt", experimentStartAt),
				zap.Time("experimentEndAt", experimentEndAt),
			)
		}
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

	es.logger.Debug("MYSQL: Evaluation found in database",
		zap.String("userId", ue.UserID),
		zap.String("featureId", ue.FeatureID),
		zap.Int32("featureVersion", ue.FeatureVersion),
		zap.String("variationId", ue.VariationID),
		zap.Int64("timestamp", ue.Timestamp),
		zap.String("timestampISO", time.Unix(ue.Timestamp, 0).Format(time.RFC3339)),
	)

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
