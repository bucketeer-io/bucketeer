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

package processor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/protobuf/types/known/wrapperspb"

	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	ecstorage "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/storage/v2"
	ec "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	ft "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	redisv3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	bqquerier "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/bigquery/querier"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/bigquery/writer"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage"
	storagev2 "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/v2"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/v2/proto/eventpersisterdwh"
	exproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

const (
	goalEventTable                = "goal_event"
	defaultRetryGoalEventInterval = 5 * time.Minute
	minLockTTL                    = 15 * time.Second
)

type goalEvtWriter struct {
	writer                  storage.GoalEventWriter
	eventStorage            ecstorage.EventStorage
	experimentClient        ec.Client
	featureClient           ft.Client
	redisClient             redisv3.Client
	locker                  *GoalEventLocker
	cache                   cachev3.ExperimentsCache
	flightgroup             singleflight.Group
	location                *time.Location
	logger                  *zap.Logger
	maxRetryGoalEventPeriod time.Duration
	retryGoalEventInterval  time.Duration
}

type GoalEventWriterOption struct {
	DataWarehouseType string
	MySQLClient       mysql.Client
	PostgresClient    postgres.Client
	BatchSize         int
}

func NewGoalEventWriter(
	ctx context.Context,
	logger *zap.Logger,
	exClient ec.Client,
	ftClient ft.Client,
	cache cachev3.ExperimentsCache,
	project, bigQueryDataSet, bigQueryDataLocation string,
	bigQueryBatchSize int,
	location *time.Location,
	redisClient redisv3.Client,
	maxRetryGoalEventPeriod time.Duration,
	retryGoalEventInterval time.Duration,
	registerer metrics.Registerer,
	options ...GoalEventWriterOption,
) (Writer, error) {
	var option GoalEventWriterOption
	if len(options) > 0 {
		option = options[0]
	}

	switch option.DataWarehouseType {
	case "mysql":
		if option.MySQLClient == nil {
			return nil, errors.New("mysql client is required when using MySQL storage")
		}

		goalStorage := storagev2.NewMysqlGoalEventStorage(option.MySQLClient)
		mysqlEventStorage := ecstorage.NewMySQLEventStorage(option.MySQLClient, logger)

		// Calculate lock TTL as 80% of retry interval
		if retryGoalEventInterval == 0 {
			retryGoalEventInterval = defaultRetryGoalEventInterval
		}
		lockTTL := time.Duration(float64(retryGoalEventInterval) * 0.8)

		// Ensure lock TTL is at least 15 seconds to allow time for:
		// - listExperiments() network call
		// - linkGoalEventByExperiment() MySQL query
		// - writer.AppendRows() MySQL write
		// Increased to 15s to handle server load and retry processing
		originalLockTTL := lockTTL
		if lockTTL < minLockTTL {
			lockTTL = minLockTTL
			logger.Info("Adjusted lock TTL to minimum",
				zap.Duration("original", originalLockTTL),
				zap.Duration("adjusted", lockTTL),
			)
		}

		w := &goalEvtWriter{
			writer:                  storage.NewMysqlGoalEventWriter(goalStorage),
			eventStorage:            mysqlEventStorage,
			experimentClient:        exClient,
			featureClient:           ftClient,
			redisClient:             redisClient,
			locker:                  NewGoalEventLocker(redisClient, lockTTL),
			cache:                   cache,
			location:                location,
			logger:                  logger,
			maxRetryGoalEventPeriod: maxRetryGoalEventPeriod,
			retryGoalEventInterval:  retryGoalEventInterval,
		}
		w.StartRetryProcessor(ctx)
		return w, nil
	case "postgres":
		if option.PostgresClient == nil {
			return nil, errors.New("postgres client is required when using Postgres storage")
		}

		goalStorage := storagev2.NewPostgresGoalEventStorage(option.PostgresClient)
		postgresEventStorage := ecstorage.NewPostgresEventStorage(option.PostgresClient, logger)

		// Calculate lock TTL as 80% of retry interval
		if retryGoalEventInterval == 0 {
			retryGoalEventInterval = defaultRetryGoalEventInterval
		}
		lockTTL := time.Duration(float64(retryGoalEventInterval) * 0.8)

		w := &goalEvtWriter{
			writer:                  storage.NewPostgresGoalEventWriter(goalStorage),
			eventStorage:            postgresEventStorage,
			experimentClient:        exClient,
			featureClient:           ftClient,
			redisClient:             redisClient,
			locker:                  NewGoalEventLocker(redisClient, lockTTL),
			cache:                   cache,
			location:                location,
			logger:                  logger,
			maxRetryGoalEventPeriod: maxRetryGoalEventPeriod,
			retryGoalEventInterval:  retryGoalEventInterval,
		}
		w.StartRetryProcessor(ctx)
		return w, nil
	case "bigquery":
		// Fall through to BigQuery implementation below
	default:
		// Default to BigQuery for backward compatibility
	}

	// BigQuery implementation
	evt := epproto.GoalEvent{}
	goalWriter, err := writer.NewWriter(
		ctx,
		project,
		bigQueryDataSet,
		goalEventTable,
		evt.ProtoReflect().Descriptor(),
		writer.WithLogger(logger),
		writer.WithMetrics(registerer),
	)
	if err != nil {
		return nil, err
	}
	eventQuerier, err := bqquerier.NewClient(
		ctx,
		project,
		bigQueryDataLocation,
		bqquerier.WithLogger(logger),
		bqquerier.WithMetrics(registerer),
	)
	if err != nil {
		return nil, err
	}
	if retryGoalEventInterval == 0 {
		retryGoalEventInterval = defaultRetryGoalEventInterval
	}
	// Calculate lock TTL as 80% of retry interval
	lockTTL := time.Duration(float64(retryGoalEventInterval) * 0.8)

	// Ensure lock TTL is at least 15 seconds to allow time for:
	// - listExperiments() network call
	// - linkGoalEventByExperiment() BigQuery query
	// - writer.AppendRows() BigQuery write
	// Increased to 15s to handle server load and retry processing
	originalLockTTL := lockTTL
	if lockTTL < minLockTTL {
		lockTTL = minLockTTL
		logger.Info("Adjusted lock TTL to minimum",
			zap.Duration("original", originalLockTTL),
			zap.Duration("adjusted", lockTTL),
		)
	}

	w := &goalEvtWriter{
		writer:                  storage.NewGoalEventWriter(goalWriter, bigQueryBatchSize),
		eventStorage:            ecstorage.NewBigQueryEventStorage(eventQuerier, bigQueryDataSet, logger),
		experimentClient:        exClient,
		featureClient:           ftClient,
		redisClient:             redisClient,
		locker:                  NewGoalEventLocker(redisClient, lockTTL),
		cache:                   cache,
		location:                location,
		logger:                  logger,
		maxRetryGoalEventPeriod: maxRetryGoalEventPeriod,
		retryGoalEventInterval:  retryGoalEventInterval,
	}
	w.StartRetryProcessor(ctx)
	return w, nil
}

func (w *goalEvtWriter) Write(
	ctx context.Context,
	envEvents environmentEventDWHMap,
) map[string]bool {
	var goalEvents []*epproto.GoalEvent
	fails := make(map[string]bool, len(envEvents))
	for environmentId, events := range envEvents {
		experiments, err := w.listExperiments(ctx, environmentId)
		if err != nil {
			w.logger.Error("failed to list experiments",
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToListExperiments).Inc()
			// Make sure to retry all the events in the next pulling
			for id := range events {
				fails[id] = true
			}
			continue
		}
		if len(experiments) == 0 {
			continue
		}
		for id, event := range events {
			switch evt := event.(type) {
			case *eventproto.GoalEvent:
				e, retriable, err := w.convToGoalEvents(ctx, evt, id, environmentId, experiments)
				if err != nil {
					if errors.Is(err, ErrExperimentNotFound) {
						// If there is nothing to link, we don't report it as an error
						subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeExperimentNotFound).Inc()
						continue
					}
					if !retriable {
						w.logger.Error(
							"Failed to convert to goal event",
							zap.Error(err),
							zap.String("id", id),
							zap.String("environmentId", environmentId),
							zap.Any("goalEvent", evt),
						)
					}
					fails[id] = retriable
					continue
				}
				goalEvents = append(goalEvents, e...)
				subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeLinked).Inc()
			default:
				w.logger.Error(
					"The event is an unexpected message type",
					zap.String("id", id),
					zap.String("environmentId", environmentId),
					zap.Any("goalEvent", evt),
				)
				fails[id] = false
			}
		}
	}
	fs, err := w.writer.AppendRows(ctx, goalEvents)
	if err != nil {
		subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToAppendGoalEvents).Inc()
		w.logger.Error("Failed to append rows to goal_event table",
			zap.Error(err),
		)
	}
	failedToAppendMap := make(map[string]*epproto.GoalEvent)
	for id, f := range fs {
		// To log which event has failed to append in the BigQuery, we need to find the event
		for _, ge := range goalEvents {
			if id == ge.Id {
				failedToAppendMap[id] = ge
			}
		}
		// Update the fails map
		fails[id] = f
	}
	if len(failedToAppendMap) > 0 {
		w.logger.Error("Failed to append goal events",
			zap.Any("goalEvents", failedToAppendMap),
		)
	}
	return fails
}

// Convert one or more goal events
func (w *goalEvtWriter) convToGoalEvents(
	ctx context.Context,
	e *eventproto.GoalEvent,
	id, environmentID string,
	experiments []*exproto.Experiment,
) ([]*epproto.GoalEvent, bool, error) {
	evals, retriable, err := w.linkGoalEvent(ctx, e, id, environmentID, e.Tag, experiments)
	if err != nil {
		return nil, retriable, err
	}
	events := make([]*epproto.GoalEvent, 0, len(evals))
	for _, eval := range evals {
		event, retriable, err := w.convToGoalEvent(e, eval, id, e.Tag, environmentID)
		if err != nil {
			return nil, retriable, err
		}
		events = append(events, event)
	}
	return events, false, nil
}

func (w *goalEvtWriter) convToGoalEvent(
	e *eventproto.GoalEvent,
	eval *ecstorage.UserEvaluation,
	id, tag, environmentID string,
) (*epproto.GoalEvent, bool, error) {
	var ud []byte
	if e.User != nil {
		var err error
		userData := make(map[string]string)
		if e.User.Data != nil {
			userData = e.User.Data
		}
		ud, err = json.Marshal(userData)
		if err != nil {
			return nil, false, err
		}
	}
	userID := getUserID(e.UserId, e.User)
	if tag == "" {
		// Tag is optional, so we insert none when is empty.
		tag = "none"
	}
	return &epproto.GoalEvent{
		Id:             id,
		GoalId:         e.GoalId,
		Value:          float32(e.Value),
		UserData:       string(ud),
		UserId:         userID,
		Tag:            tag,
		SourceId:       e.SourceId.String(),
		EnvironmentId:  environmentID,
		Timestamp:      time.Unix(e.Timestamp, 0).UnixMicro(),
		FeatureId:      eval.FeatureID,
		FeatureVersion: eval.FeatureVersion,
		VariationId:    eval.VariationID,
		Reason:         eval.Reason,
	}, false, nil
}

func (w *goalEvtWriter) linkGoalEvent(
	ctx context.Context,
	event *eventproto.GoalEvent,
	id, environmentID, tag string,
	experiments []*exproto.Experiment,
) ([]*ecstorage.UserEvaluation, bool, error) {
	evalExp, retriable, err := w.linkGoalEventByExperiment(ctx, event, id, environmentID, tag, experiments)
	if err != nil {
		return nil, retriable, err
	}
	return evalExp, false, nil
}

// Link one or more experiments by goal ID
func (w *goalEvtWriter) linkGoalEventByExperiment(
	ctx context.Context,
	event *eventproto.GoalEvent,
	id, environmentID, tag string,
	experiments []*exproto.Experiment,
) ([]*ecstorage.UserEvaluation, bool, error) {
	// Find the experiment by goal ID
	// TODO: we must change the console UI not to allow creating
	// multiple experiments running at the same time,
	// using the same feature flag id and goal id
	var exps []*exproto.Experiment
	for _, exp := range experiments {
		if w.findGoalID(event.GoalId, exp.GoalIds) {
			// If the goal event was issued before the experiment started running,
			// we ignore those events to avoid issues in the conversion rate
			if exp.StartAt > event.Timestamp {
				subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeEventOlderThanExperiment).Inc()
				continue
			}
			// If the goal event was issued after the experiment ended,
			// we ignore those events to avoid issues in the conversion rate
			if exp.StopAt < event.Timestamp {
				subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeGoalEventIssuedAfterExperimentEnded).Inc()
				continue
			}
			exps = append(exps, exp)
		}
	}
	if len(exps) == 0 {
		return nil, false, ErrExperimentNotFound
	}
	evals := make([]*ecstorage.UserEvaluation, 0, len(exps))
	for _, exp := range exps {
		// Get the user evaluation using the experiment info
		userID := getUserID(event.UserId, event.User)
		eval, err := w.getUserEvaluation(
			ctx,
			environmentID,
			userID,
			exp.FeatureId,
			exp.FeatureVersion,
			exp.StartAt,
			exp.StopAt,
		)
		if err != nil {
			if errors.Is(err, ecstorage.ErrBQNoResultsFound) {
				w.logger.Error("Evaluation not found",
					zap.Error(err),
					zap.String("environmentId", environmentID),
					zap.Any("goalEvent", event),
				)
				subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeUserEvaluationNotFound).Inc()
				if err := w.storeRetryMessage(&retryMessage{
					GoalEvent:     event,
					EnvironmentID: environmentID,
					RetryCount:    0,
					ID:            id,
				}); err != nil {
					subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToStoreRetryMessage).Inc()
					w.logger.Error("Failed to store retry message",
						zap.Error(err),
						zap.String("environmentId", environmentID),
						zap.Any("goalEvent", event),
					)
					return nil, true, err
				}
				return nil, false, err
			}
			w.logger.Error("failed to get user evaluation",
				zap.Error(err),
				zap.String("environmentId", environmentID),
				zap.Any("goalEvent", event),
			)
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToGetUserEvaluation).Inc()
			return nil, true, err
		}
		w.logger.Debug("Evaluation found",
			zap.String("environmentId", environmentID),
			zap.Any("goalEvent", event),
			zap.Any("evaluation", eval),
		)
		// Skip goal events that occurred before the evaluation timestamp.
		// This is intentional to trigger retry logic for proper event linking.
		if event.Timestamp < eval.Timestamp {
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeGoalEventIssuedBeforeEvaluation).Inc()
			w.logger.Error("Goal event issued before evaluation",
				zap.String("environmentId", environmentID),
				zap.Any("goalEvent", event),
				zap.Any("evaluation", eval),
			)
			continue
		}
		evals = append(evals, eval)
	}
	return evals, false, nil
}

func (w *goalEvtWriter) findGoalID(id string, goalIDs []string) bool {
	for _, goalID := range goalIDs {
		if id == goalID {
			return true
		}
	}
	return false
}

func (w *goalEvtWriter) listExperiments(
	ctx context.Context,
	environmentId string,
) ([]*exproto.Experiment, error) {
	experiments, err, _ := w.flightgroup.Do(
		fmt.Sprintf("%s:%s", environmentId, "listExperiments"),
		func() (interface{}, error) {
			// Get the experiment cache
			expList, err := w.cache.Get(environmentId)
			if err == nil {
				return expList.Experiments, nil
			}
			// Get the experiments from the DB
			resp, err := w.experimentClient.ListExperiments(ctx, &exproto.ListExperimentsRequest{
				// Because the evaluation and goal events may be sent with a delay
				// for many reasons from the client side, we still calculate
				// the results for two days after it stopped.
				StopAt:        time.Now().In(w.location).Add(-2 * day).Unix(),
				PageSize:      0,
				EnvironmentId: environmentId,
				Statuses: []exproto.Experiment_Status{
					exproto.Experiment_RUNNING,
					exproto.Experiment_STOPPED,
				},
				Archived: wrapperspb.Bool(false),
			})
			if err != nil {
				return nil, err
			}
			return resp.Experiments, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return experiments.([]*exproto.Experiment), nil
}

func (w *goalEvtWriter) getUserEvaluation(
	ctx context.Context,
	environmentID, userID, featureID string,
	featureVersion int32,
	experimentStartAt, experimentEndAt int64,
) (*ecstorage.UserEvaluation, error) {
	eval, err := w.eventStorage.QueryUserEvaluation(
		ctx,
		environmentID,
		userID,
		featureID,
		featureVersion,
		time.Unix(experimentStartAt, 0),
		time.Unix(experimentEndAt, 0),
	)
	if err != nil {
		return nil, err
	}
	return eval, nil
}
