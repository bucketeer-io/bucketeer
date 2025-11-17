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

package processor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
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
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage"
	storagev2 "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/v2"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/v2/proto/eventpersisterdwh"
	exproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

const (
	goalEventTable                = "goal_event"
	defaultRetryGoalEventInterval = 5 * time.Minute
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

	// experimentNotFoundCache is a negative cache keyed by "environmentId:goalId".
	// Once we've confirmed (via DB) that there's no experiment for a given goal,
	// we remember it so we don't keep hitting the DB when clients keep sending
	// events for deleted/nonexistent experiments.
	experimentNotFoundCache sync.Map // key: "envId:goalId" -> struct{}
}

type GoalEventWriterOption struct {
	DataWarehouseType string
	MySQLClient       mysql.Client
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

	logger.Info("NewGoalEventWriter: data warehouse configuration",
		zap.String("dataWarehouseType", option.DataWarehouseType),
		zap.Bool("hasMySQLClient", option.MySQLClient != nil),
		zap.Int("optionsCount", len(options)),
	)

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

	w := &goalEvtWriter{
		writer:                  storage.NewGoalEventWriter(goalWriter, bigQueryBatchSize),
		eventStorage:            ecstorage.NewEventStorage(eventQuerier, bigQueryDataSet, logger),
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
				userID := getUserID(evt.UserId, evt.User)
				w.logger.Info("📥 GOAL: Received from Redis Stream",
					zap.String("eventId", id),
					zap.String("environmentId", environmentId),
					zap.String("goalId", evt.GoalId),
					zap.String("userId", userID),
					zap.String("goalEventUserId", evt.UserId),
					zap.Int64("timestamp", evt.Timestamp),
					zap.String("timestampISO", time.Unix(evt.Timestamp, 0).Format(time.RFC3339)),
					zap.Float64("value", evt.Value),
					zap.String("tag", evt.Tag),
				)
				e, retriable, err := w.convToGoalEvents(ctx, evt, id, environmentId, experiments)
				if err != nil {
					if errors.Is(err, ErrExperimentNotFound) {
						// If there is nothing to link, we don't report it as an error
						subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeExperimentNotFound).Inc()
						w.logger.Debug("Goal event skipped: experiment not found",
							zap.String("eventId", id),
							zap.String("goalId", evt.GoalId),
						)
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
				w.logger.Info("GOAL: Linked successfully, will write to data warehouse",
					zap.String("eventId", id),
					zap.String("goalId", evt.GoalId),
					zap.String("userId", getUserID(evt.UserId, evt.User)),
					zap.Int("linkedEventCount", len(e)),
					zap.Strings("linkedEventIds", func() []string {
						ids := make([]string, len(e))
						for i, ge := range e {
							ids[i] = ge.Id
						}
						return ids
					}()),
				)
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
	if len(goalEvents) > 0 {
		w.logger.Info("GOAL: Writing to data warehouse",
			zap.Int("eventCount", len(goalEvents)),
			zap.Strings("eventIds", func() []string {
				ids := make([]string, len(goalEvents))
				for i, ge := range goalEvents {
					ids[i] = ge.Id
				}
				return ids
			}()),
		)
	}
	fs, err := w.writer.AppendRows(ctx, goalEvents)
	if err != nil {
		subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToAppendGoalEvents).Inc()
		w.logger.Error("GOAL: Failed to append rows to goal_event table",
			zap.Error(err),
			zap.Int("eventCount", len(goalEvents)),
		)
	}
	failedToAppendMap := make(map[string]*epproto.GoalEvent)
	successfulEvents := make([]string, 0)
	for id, f := range fs {
		// To log which event has failed to append in the BigQuery, we need to find the event
		for _, ge := range goalEvents {
			if id == ge.Id {
				if f {
					failedToAppendMap[id] = ge
				} else {
					successfulEvents = append(successfulEvents, id)
				}
			}
		}
		// Update the fails map
		fails[id] = f
	}
	if len(successfulEvents) > 0 {
		envID := "unknown"
		if len(goalEvents) > 0 {
			envID = goalEvents[0].EnvironmentId
		}
		w.logger.Info("GOAL: Successfully wrote to data warehouse",
			zap.Int("count", len(successfulEvents)),
			zap.Strings("eventIds", successfulEvents),
			zap.String("environmentId", envID),
		)
	}
	if len(failedToAppendMap) > 0 {
		w.logger.Error("GOAL: Failed to append goal events",
			zap.Int("failedCount", len(failedToAppendMap)),
			zap.Strings("failedEventIds", func() []string {
				ids := make([]string, 0, len(failedToAppendMap))
				for id := range failedToAppendMap {
					ids = append(ids, id)
				}
				return ids
			}()),
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
	// First try using the current snapshot (usually from cache).
	evals, retriable, err := w.linkGoalEvent(ctx, e, id, environmentID, e.Tag, experiments)
	if err != nil {
		// If experiments were not found in the snapshot, we may be seeing a race
		// between cache refresh and new experiment creation (especially in e2e).
		// Do a one-time slow-path DB lookup per (environmentID, goalId) before
		// we conclude "experiment really doesn't exist".
		if errors.Is(err, ErrExperimentNotFound) {
			if w.isExperimentKnownMissing(environmentID, e.GoalId) {
				// We've already confirmed this goal has no experiment in DB.
				return nil, false, ErrExperimentNotFound
			}

			w.logger.Debug("Goal event linking fallback: experiment not found in cache snapshot, reloading from DB",
				zap.String("environmentId", environmentID),
				zap.String("goalId", e.GoalId),
			)

			dbExperiments, dbErr := w.listExperimentsFromDB(ctx, environmentID)
			if dbErr != nil {
				// Treat as infra error: caller may decide to retry.
				w.logger.Error("Failed to list experiments from DB in fallback path",
					zap.Error(dbErr),
					zap.String("environmentId", environmentID),
				)
				return nil, true, dbErr
			}

			// Log DB fallback results for debugging
			w.logger.Info("DB fallback: Listed experiments from DB",
				zap.String("environmentId", environmentID),
				zap.String("goalId", e.GoalId),
				zap.Int("experimentCount", len(dbExperiments)),
				zap.Int64("goalTimestamp", e.Timestamp),
				zap.String("goalTimestampISO", time.Unix(e.Timestamp, 0).Format(time.RFC3339)),
			)
			// Log goal IDs from experiments for debugging
			for _, exp := range dbExperiments {
				w.logger.Debug("DB fallback: Experiment found",
					zap.String("environmentId", environmentID),
					zap.String("goalId", e.GoalId),
					zap.String("experimentId", exp.Id),
					zap.Strings("experimentGoalIds", exp.GoalIds),
					zap.Int64("expStartAt", exp.StartAt),
					zap.Int64("expStopAt", exp.StopAt),
					zap.String("expStartAtISO", time.Unix(exp.StartAt, 0).Format(time.RFC3339)),
					zap.String("expStopAtISO", time.Unix(exp.StopAt, 0).Format(time.RFC3339)),
				)
			}

			evals, retriable, err = w.linkGoalEvent(ctx, e, id, environmentID, e.Tag, dbExperiments)
			if err != nil {
				if errors.Is(err, ErrExperimentNotFound) {
					// Check if any experiment has matching goal ID (even if timestamp is outside window)
					hasMatchingGoalID := false
					for _, exp := range dbExperiments {
						if w.findGoalID(e.GoalId, exp.GoalIds) {
							hasMatchingGoalID = true
							break
						}
					}
					if !hasMatchingGoalID {
						// No experiment exists with this goal ID - mark as missing
						w.markExperimentMissing(environmentID, e.GoalId)
						w.logger.Debug("Confirmed experiment not found after DB fallback; caching negative result",
							zap.String("environmentId", environmentID),
							zap.String("goalId", e.GoalId),
						)
					} else {
						// Experiment exists but timestamp is outside window - don't mark as missing
						w.logger.Debug("Experiment exists but goal event timestamp outside window; skipping event without marking as missing",
							zap.String("environmentId", environmentID),
							zap.String("goalId", e.GoalId),
							zap.Int64("goalTimestamp", e.Timestamp),
							zap.String("goalTimestampISO", time.Unix(e.Timestamp, 0).Format(time.RFC3339)),
						)
					}
					return nil, false, ErrExperimentNotFound
				}
				return nil, retriable, err
			}
			// Successfully found experiment via DB fallback - clear any negative cache entry
			w.clearExperimentMissing(environmentID, e.GoalId)
		} else {
			return nil, retriable, err
		}
	}
	// Successfully found experiment via cache - clear any negative cache entry
	w.clearExperimentMissing(environmentID, e.GoalId)
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
		// Log why no experiments matched for debugging
		hasMatchingGoalID := false
		for _, exp := range experiments {
			if w.findGoalID(event.GoalId, exp.GoalIds) {
				hasMatchingGoalID = true
				w.logger.Debug("Experiment has matching goal ID but timestamp outside window",
					zap.String("environmentId", environmentID),
					zap.String("goalId", event.GoalId),
					zap.String("experimentId", exp.Id),
					zap.Int64("goalTimestamp", event.Timestamp),
					zap.String("goalTimestampISO", time.Unix(event.Timestamp, 0).Format(time.RFC3339)),
					zap.Int64("expStartAt", exp.StartAt),
					zap.Int64("expStopAt", exp.StopAt),
					zap.String("expStartAtISO", time.Unix(exp.StartAt, 0).Format(time.RFC3339)),
					zap.String("expStopAtISO", time.Unix(exp.StopAt, 0).Format(time.RFC3339)),
					zap.Bool("beforeStart", exp.StartAt > event.Timestamp),
					zap.Bool("afterStop", exp.StopAt < event.Timestamp),
				)
			}
		}
		if !hasMatchingGoalID {
			w.logger.Debug("No experiment found with matching goal ID",
				zap.String("environmentId", environmentID),
				zap.String("goalId", event.GoalId),
				zap.Int("totalExperiments", len(experiments)),
			)
		}
		return nil, false, ErrExperimentNotFound
	}
	evals := make([]*ecstorage.UserEvaluation, 0, len(exps))
	for _, exp := range exps {
		// Get the user evaluation using the experiment info
		userID := getUserID(event.UserId, event.User)
		w.logger.Info("🔍 GOAL: Querying for evaluation",
			zap.String("environmentId", environmentID),
			zap.String("userId", userID),
			zap.String("goalEventUserId", event.UserId),
			zap.String("featureId", exp.FeatureId),
			zap.Int32("featureVersion", exp.FeatureVersion),
			zap.Int64("expStartAt", exp.StartAt),
			zap.String("expStartAtISO", time.Unix(exp.StartAt, 0).Format(time.RFC3339)),
			zap.Int64("expStopAt", exp.StopAt),
			zap.String("expStopAtISO", time.Unix(exp.StopAt, 0).Format(time.RFC3339)),
			zap.String("goalId", event.GoalId),
			zap.Int64("goalTimestamp", event.Timestamp),
			zap.String("goalTimestampISO", time.Unix(event.Timestamp, 0).Format(time.RFC3339)),
		)
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
			w.logger.Debug("GOAL: Evaluation query result",
				zap.String("environmentId", environmentID),
				zap.String("userId", userID),
				zap.String("featureId", exp.FeatureId),
				zap.Error(err),
			)
			if errors.Is(err, ecstorage.ErrNoResultsFound) {
				// Use the transport message ID to identify this specific goal event instance.
				// This is stable per message and ensures idempotency across retries without
				// collapsing distinct goal events that happen within the same second.
				retryID := id
				if retryID == "" {
					// Fallback to a deterministic compound ID (no randomness to preserve idempotency)
					retryID = fmt.Sprintf("%s-%s-%d", event.GoalId, userID, event.Timestamp)
				}
				retryKey := fmt.Sprintf("%s:%s:%s", environmentID, retryGoalEventKeyKind, retryID)

				// Check if retry key already exists to avoid overwriting existing retries
				existingData, getErr := w.redisClient.Get(retryKey)
				if getErr == nil && existingData != nil {
					// Retry key already exists, don't overwrite it
					w.logger.Debug("Retry key already exists, skipping initial retry storage",
						zap.String("environmentId", environmentID),
						zap.String("goalId", event.GoalId),
						zap.String("userId", userID),
						zap.Int64("timestamp", event.Timestamp),
						zap.String("retryId", retryID),
						zap.String("retryKey", retryKey),
					)
					subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeUserEvaluationNotFound).Inc()
					return nil, false, err
				}
				// If getErr != nil (including ErrNil when key doesn't exist), proceed to store new retry

				w.logger.Info("GOAL: Evaluation not found, storing retry message",
					zap.String("environmentId", environmentID),
					zap.String("goalId", event.GoalId),
					zap.String("userId", userID),
					zap.Int64("timestamp", event.Timestamp),
					zap.String("timestampISO", time.Unix(event.Timestamp, 0).Format(time.RFC3339)),
					zap.String("retryId", retryID),
					zap.String("retryKey", retryKey),
					zap.String("featureId", exp.FeatureId),
					zap.Int32("featureVersion", exp.FeatureVersion),
					zap.Int64("expStartAt", exp.StartAt),
					zap.Int64("expStopAt", exp.StopAt),
					zap.String("expStartAtISO", time.Unix(exp.StartAt, 0).Format(time.RFC3339)),
					zap.String("expStopAtISO", time.Unix(exp.StopAt, 0).Format(time.RFC3339)),
				)
				subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeUserEvaluationNotFound).Inc()
				if err := w.storeRetryMessage(&retryMessage{
					GoalEvent:     event,
					EnvironmentID: environmentID,
					RetryCount:    0,
					ID:            retryID,
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
		w.logger.Info("GOAL: Evaluation found for linking",
			zap.String("environmentId", environmentID),
			zap.String("goalId", event.GoalId),
			zap.String("userId", userID),
			zap.Int64("goalTimestamp", event.Timestamp),
			zap.String("goalTimestampISO", time.Unix(event.Timestamp, 0).Format(time.RFC3339)),
			zap.String("featureId", eval.FeatureID),
			zap.Int32("featureVersion", eval.FeatureVersion),
			zap.String("variationId", eval.VariationID),
			zap.Int64("evalTimestamp", eval.Timestamp),
			zap.String("evalTimestampISO", time.Unix(eval.Timestamp, 0).Format(time.RFC3339)),
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

// listExperiments returns the experiment snapshot, using the cache when possible
// and falling back to the experiment service / DB on cache miss.
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

// listExperimentsFromDB bypasses the cache and always queries the experiment
// service / DB. This is used as a slow path when we see a cache snapshot that
// doesn't contain any experiment for a given goalId, to avoid losing events
// due to cache staleness (e.g. e2e races).
func (w *goalEvtWriter) listExperimentsFromDB(
	ctx context.Context,
	environmentId string,
) ([]*exproto.Experiment, error) {
	resp, err := w.experimentClient.ListExperiments(ctx, &exproto.ListExperimentsRequest{
		// Same time window and filters as listExperiments.
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
		time.Unix(experimentStartAt, 0).UTC(),
		time.Unix(experimentEndAt, 0).UTC(),
	)
	if err != nil {
		return nil, err
	}
	return eval, nil
}

// experimentMissingCacheKey builds the key for the negative cache.
func (w *goalEvtWriter) experimentMissingCacheKey(environmentID, goalID string) string {
	return environmentID + ":" + goalID
}

// isExperimentKnownMissing returns true if we've already confirmed that there
// is no experiment for this (environment, goalId) pair via a DB check.
func (w *goalEvtWriter) isExperimentKnownMissing(environmentID, goalID string) bool {
	if goalID == "" {
		return false
	}
	key := w.experimentMissingCacheKey(environmentID, goalID)
	_, ok := w.experimentNotFoundCache.Load(key)
	return ok
}

// markExperimentMissing records that we've confirmed (via DB) there is no
// experiment for this (environment, goalId) pair, so we can skip future DB
// lookups when clients keep sending goal events for deleted/nonexistent
// experiments.
func (w *goalEvtWriter) markExperimentMissing(environmentID, goalID string) {
	if goalID == "" {
		return
	}
	key := w.experimentMissingCacheKey(environmentID, goalID)
	w.experimentNotFoundCache.Store(key, struct{}{})
}

// clearExperimentMissing removes the entry from the negative cache when an
// experiment is found. This is important when experiments are created after
// we've marked them as missing, or when cache is updated with new experiments.
func (w *goalEvtWriter) clearExperimentMissing(environmentID, goalID string) {
	if goalID == "" {
		return
	}
	key := w.experimentMissingCacheKey(environmentID, goalID)
	w.experimentNotFoundCache.Delete(key)
	w.logger.Debug("Cleared experiment missing cache entry",
		zap.String("environmentId", environmentID),
		zap.String("goalId", goalID),
	)
}
