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
	ec "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/bigquery/writer"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage"
	storagev2 "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/v2"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/v2/proto/eventpersisterdwh"
	exproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

const (
	day                  = 24 * time.Hour
	evaluationEventTable = "evaluation_event"
)

type evalEvtWriter struct {
	writer           storage.EvalEventWriter
	experimentClient ec.Client
	cache            cachev3.ExperimentsCache
	flightgroup      singleflight.Group
	location         *time.Location
	logger           *zap.Logger

	// experimentNotFoundCache is a negative cache keyed by "environmentId:featureId:v<version>".
	// Once we've confirmed (via DB) that there's no experiment for a given
	// (featureId, featureVersion) in an environment, we remember it so we don't
	// keep hitting the DB when clients keep sending evaluation events that can't
	// be linked.
	experimentNotFoundCache sync.Map // key: "envId:featureId:v<version>" -> struct{}
}

type EvalEventWriterOption struct {
	DataWarehouseType string
	MySQLClient       mysql.Client
	BatchSize         int
}

func NewEvalEventWriter(
	ctx context.Context,
	l *zap.Logger,
	exClient ec.Client,
	cache cachev3.ExperimentsCache,
	project, ds string,
	size int,
	location *time.Location,
	registerer metrics.Registerer,
	options ...EvalEventWriterOption,
) (Writer, error) {
	var option EvalEventWriterOption
	if len(options) > 0 {
		option = options[0]
	}

	l.Info("NewEvalEventWriter: data warehouse configuration",
		zap.String("dataWarehouseType", option.DataWarehouseType),
		zap.Bool("hasMySQLClient", option.MySQLClient != nil),
		zap.Int("optionsCount", len(options)),
	)

	switch option.DataWarehouseType {
	case "mysql":
		if option.MySQLClient == nil {
			return nil, errors.New("mysql client is required when using MySQL storage")
		}

		evalStorage := storagev2.NewMysqlEvaluationEventStorageWithLogger(option.MySQLClient, l)

		return &evalEvtWriter{
			writer:           storage.NewMysqlEvalEventWriter(evalStorage),
			experimentClient: exClient,
			cache:            cache,
			location:         location,
			logger:           l,
		}, nil

	case "bigquery":
		// Fall through to BigQuery implementation below
	default:
		// Default to BigQuery for backward compatibility
	}

	// BigQuery implementation
	evt := epproto.EvaluationEvent{}
	evalQuery, err := writer.NewWriter(
		ctx,
		project,
		ds,
		evaluationEventTable,
		evt.ProtoReflect().Descriptor(),
		writer.WithLogger(l),
		writer.WithMetrics(registerer),
	)
	if err != nil {
		return nil, err
	}
	return &evalEvtWriter{
		writer:           storage.NewEvalEventWriter(evalQuery, size),
		experimentClient: exClient,
		cache:            cache,
		location:         location,
		logger:           l,
	}, nil
}

func (w *evalEvtWriter) Write(
	ctx context.Context,
	envEvents environmentEventDWHMap,
) map[string]bool {
	var evalEvents []*epproto.EvaluationEvent
	fails := make(map[string]bool, len(envEvents))
	for environmentId, events := range envEvents {
		experiments, err := w.listExperiments(ctx, environmentId)
		if err != nil {
			w.logger.Error("failed to list experiments",
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)
			subscriberHandledCounter.WithLabelValues(subscriberEvaluationEventDWH, codeFailedToListExperiments).Inc()
			// Make sure to retry all the events in the next pulling
			for id := range events {
				fails[id] = true
			}
			continue
		}
		if len(experiments) == 0 {
			w.logger.Info("⚠️ No experiments found for environment, skipping evaluation events",
				zap.String("environmentId", environmentId),
				zap.Int("eventCount", len(events)),
				zap.Strings("eventFeatureIds", func() []string {
					ids := make([]string, 0, len(events))
					for _, event := range events {
						if evt, ok := event.(*eventproto.EvaluationEvent); ok {
							ids = append(ids, fmt.Sprintf("%s:v%d", evt.FeatureId, evt.FeatureVersion))
						}
					}
					return ids
				}()),
			)
			// ACK these messages since there's nothing to link them to
			// They won't be added to fails, so they'll be ACKed in events_dwh_persister.go
			continue
		}
		for id, event := range events {
			switch evt := event.(type) {
			case *eventproto.EvaluationEvent:
				w.logger.Info("📥 EVAL: Received from Redis Stream",
					zap.String("eventId", id),
					zap.String("environmentId", environmentId),
					zap.String("featureId", evt.FeatureId),
					zap.Int32("featureVersion", evt.FeatureVersion),
					zap.String("userId", evt.UserId),
					zap.String("evalEventUserId", evt.UserId), // Duplicate for easier filtering
					zap.String("variationId", evt.VariationId),
					zap.Int64("timestamp", evt.Timestamp),
					zap.String("timestampISO", time.Unix(evt.Timestamp, 0).Format(time.RFC3339)),
				)
				e, retriable, err := w.convToEvaluationEvent(ctx, evt, id, environmentId, experiments)
				if err != nil {
					// If there is nothing to link, we don't report it as an error
					if errors.Is(err, ErrExperimentNotFound) {
						subscriberHandledCounter.WithLabelValues(subscriberEvaluationEventDWH, codeExperimentNotFound).Inc()
						w.logger.Info("⏭️ EVAL: Skipped - experiment not found",
							zap.String("eventId", id),
							zap.String("environmentId", environmentId),
							zap.String("featureId", evt.FeatureId),
							zap.Int32("featureVersion", evt.FeatureVersion),
							zap.String("userId", evt.UserId),
							zap.Int("totalExperiments", len(experiments)),
							zap.Strings("experimentFeatureIds", func() []string {
								ids := make([]string, len(experiments))
								for i, exp := range experiments {
									ids[i] = fmt.Sprintf("%s:v%d", exp.FeatureId, exp.FeatureVersion)
								}
								return ids
							}()),
						)
						// This event will be ACKed (not added to fails)
						continue
					}
					if errors.Is(err, ErrEvaluationEventIssuedAfterExperimentEnded) {
						subscriberHandledCounter.WithLabelValues(subscriberEvaluationEventDWH, codeEventIssuedAfterExperimentEnded).Inc()
						w.logger.Info("⏭️ EVAL: Skipped - issued after experiment ended",
							zap.String("eventId", id),
							zap.String("environmentId", environmentId),
							zap.String("featureId", evt.FeatureId),
							zap.Int32("featureVersion", evt.FeatureVersion),
							zap.String("userId", evt.UserId),
							zap.Int64("eventTimestamp", evt.Timestamp),
						)
						// This event will be ACKed (not added to fails)
						continue
					}
					if !retriable {
						w.logger.Error(
							"Failed to convert to evaluation event",
							zap.Error(err),
							zap.String("id", id),
							zap.String("environmentId", environmentId),
							zap.Any("evalEvent", evt),
						)
					}
					fails[id] = retriable
					continue
				}
				evalEvents = append(evalEvents, e)
				subscriberHandledCounter.WithLabelValues(subscriberEvaluationEventDWH, codeLinked).Inc()
				w.logger.Debug("Evaluation event linked successfully, will write to data warehouse",
					zap.String("eventId", id),
					zap.String("featureId", evt.FeatureId),
					zap.String("userId", evt.UserId),
				)
			default:
				w.logger.Error(
					"The event is an unexpected message type",
					zap.String("id", id),
					zap.String("environmentId", environmentId),
					zap.Any("evalEvent", evt),
				)
				fails[id] = false
			}
		}
	}
	if len(evalEvents) > 0 {
		w.logger.Info("💾 EVAL: Writing to data warehouse",
			zap.Int("eventCount", len(evalEvents)),
			zap.Strings("eventIds", func() []string {
				ids := make([]string, len(evalEvents))
				for i, e := range evalEvents {
					ids[i] = e.Id
				}
				return ids
			}()),
		)
	}
	fs, err := w.writer.AppendRows(ctx, evalEvents)
	if err != nil {
		subscriberHandledCounter.WithLabelValues(subscriberEvaluationEventDWH, codeFailedToAppendEvaluationEvents).Inc()
		w.logger.Error("❌ EVAL: Failed to append rows to evaluation_event table",
			zap.Error(err),
			zap.Int("eventCount", len(evalEvents)),
		)
	}
	failedToAppendMap := make(map[string]*epproto.EvaluationEvent)
	successfulEvents := make([]string, 0)
	for id, f := range fs {
		// To log which event has failed to append in the BigQuery, we need to find the event
		for _, ee := range evalEvents {
			if id == ee.Id {
				if f {
					failedToAppendMap[id] = ee
				} else {
					successfulEvents = append(successfulEvents, id)
				}
			}
		}
		// Update the fails map
		fails[id] = f
	}
	if len(successfulEvents) > 0 {
		w.logger.Info("✅ EVAL: Successfully wrote to data warehouse",
			zap.Int("count", len(successfulEvents)),
			zap.Strings("eventIds", successfulEvents),
		)
	}
	if len(failedToAppendMap) > 0 {
		w.logger.Error("❌ EVAL: Failed to append evaluation events",
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

func (w *evalEvtWriter) convToEvaluationEvent(
	ctx context.Context,
	e *eventproto.EvaluationEvent,
	id, environmentId string,
	experiments []*exproto.Experiment,
) (*epproto.EvaluationEvent, bool, error) {
	// Try to find the experiment in the current snapshot (usually cache-backed).
	exp, retriable, err := w.findExperimentForEvaluation(
		ctx,
		environmentId,
		e.FeatureId,
		e.FeatureVersion,
		experiments,
	)
	if err != nil {
		return nil, retriable, err
	}
	if exp.StopAt < e.Timestamp {
		return nil, false, ErrEvaluationEventIssuedAfterExperimentEnded
	}

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
	tag := e.Tag
	if tag == "" {
		// Tag is optional, so we insert none when is empty.
		tag = "none"
	}
	variationID := e.VariationId
	if variationID == "" {
		// When the default value is returned to the app in the SDK,
		// it creates a default evaluation event with an empty variation ID.
		// Because we don't have the variation ID we insert "default" as the variation_id
		variationID = "default"
	}
	return &epproto.EvaluationEvent{
		Id:             id,
		FeatureId:      e.FeatureId,
		FeatureVersion: e.FeatureVersion,
		UserData:       string(ud),
		UserId:         userID,
		VariationId:    variationID,
		Reason:         e.Reason.Type.String(),
		Tag:            tag,
		SourceId:       e.SourceId.String(),
		EnvironmentId:  environmentId,
		Timestamp:      time.Unix(e.Timestamp, 0).UnixMicro(),
	}, false, nil
}

// findExperimentForEvaluation finds the experiment for a given (featureId,
// featureVersion) using the provided snapshot first, and if not found,
// optionally falls back to a direct DB/ListExperiments call. It also uses a
// negative cache so that experiments confirmed missing don't trigger repeated
// DB lookups when clients keep sending events for them.
func (w *evalEvtWriter) findExperimentForEvaluation(
	ctx context.Context,
	environmentId, featureID string,
	featureVersion int32,
	experiments []*exproto.Experiment,
) (*exproto.Experiment, bool, error) {
	// 1. Fast path: from the current snapshot (cache-backed).
	exp := w.existExperiment(experiments, featureID, featureVersion)
	if exp != nil {
		// Successfully found experiment via cache - clear any negative cache entry
		w.clearExperimentMissing(environmentId, featureID, featureVersion)
		return exp, false, nil
	}

	// 2. If we've already confirmed there's no experiment for this
	// (environment, featureId, featureVersion), don't hit the DB again.
	if w.isExperimentKnownMissing(environmentId, featureID, featureVersion) {
		return nil, false, ErrExperimentNotFound
	}

	w.logger.Debug("EVAL: Experiment not found in cache snapshot, reloading from DB",
		zap.String("environmentId", environmentId),
		zap.String("featureId", featureID),
		zap.Int32("featureVersion", featureVersion),
	)

	// 3. Slow path: bypass cache and query DB/experiment service directly.
	dbExperiments, err := w.listExperimentsFromDB(ctx, environmentId)
	if err != nil {
		// Treat as infra error: caller may decide to retry.
		w.logger.Error("EVAL: Failed to list experiments from DB in fallback path",
			zap.Error(err),
			zap.String("environmentId", environmentId),
		)
		return nil, true, err
	}

	exp = w.existExperiment(dbExperiments, featureID, featureVersion)
	if exp == nil {
		// 4. Still not found after DB: mark as missing to avoid future DB hits.
		w.markExperimentMissing(environmentId, featureID, featureVersion)
		w.logger.Debug("EVAL: Confirmed experiment not found after DB fallback; caching negative result",
			zap.String("environmentId", environmentId),
			zap.String("featureId", featureID),
			zap.Int32("featureVersion", featureVersion),
		)
		return nil, false, ErrExperimentNotFound
	}

	// Successfully found experiment via DB fallback - clear any negative cache entry
	w.clearExperimentMissing(environmentId, featureID, featureVersion)
	return exp, false, nil
}

func (w *evalEvtWriter) existExperiment(
	es []*exproto.Experiment,
	fID string,
	fVersion int32,
) *exproto.Experiment {
	for _, e := range es {
		if e.FeatureId == fID && e.FeatureVersion == fVersion {
			return e
		}
	}
	return nil
}

func (w *evalEvtWriter) listExperiments(
	ctx context.Context,
	environmentId string,
) ([]*exproto.Experiment, error) {
	experiments, err, _ := w.flightgroup.Do(
		fmt.Sprintf("%s:%s", environmentId, "listExperiments"),
		func() (interface{}, error) {
			// Get the experiment cache
			expList, err := w.cache.Get(environmentId)
			if err == nil {
				w.logger.Debug("Using experiment cache",
					zap.String("environmentId", environmentId),
					zap.Int("experimentCount", len(expList.Experiments)),
				)
				return expList.Experiments, nil
			}
			w.logger.Debug("Cache miss, querying DB for experiments",
				zap.String("environmentId", environmentId),
				zap.Error(err),
			)
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
			w.logger.Debug("Queried DB for experiments",
				zap.String("environmentId", environmentId),
				zap.Int("experimentCount", len(resp.Experiments)),
			)
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
// doesn't contain an experiment for a given (featureId, featureVersion),
// to avoid losing events due to cache staleness (e.g. e2e races).
func (w *evalEvtWriter) listExperimentsFromDB(
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

// experimentMissingCacheKey builds the key for the negative cache.
func (w *evalEvtWriter) experimentMissingCacheKey(
	environmentID, featureID string,
	featureVersion int32,
) string {
	return fmt.Sprintf("%s:%s:v%d", environmentID, featureID, featureVersion)
}

// isExperimentKnownMissing returns true if we've already confirmed that there
// is no experiment for this (environment, featureId, featureVersion) pair via a
// DB check.
func (w *evalEvtWriter) isExperimentKnownMissing(
	environmentID, featureID string,
	featureVersion int32,
) bool {
	if featureID == "" {
		return false
	}
	key := w.experimentMissingCacheKey(environmentID, featureID, featureVersion)
	_, ok := w.experimentNotFoundCache.Load(key)
	return ok
}

// markExperimentMissing records that we've confirmed (via DB) there is no
// experiment for this (environment, featureId, featureVersion) pair, so we can
// skip future DB lookups when clients keep sending evaluation events for
// deleted/nonexistent experiments.
func (w *evalEvtWriter) markExperimentMissing(
	environmentID, featureID string,
	featureVersion int32,
) {
	if featureID == "" {
		return
	}
	key := w.experimentMissingCacheKey(environmentID, featureID, featureVersion)
	w.experimentNotFoundCache.Store(key, struct{}{})
}

// clearExperimentMissing removes the entry from the negative cache when an
// experiment is found. This is important when experiments are created after
// we've marked them as missing, or when cache is updated with new experiments.
func (w *evalEvtWriter) clearExperimentMissing(
	environmentID, featureID string,
	featureVersion int32,
) {
	if featureID == "" {
		return
	}
	key := w.experimentMissingCacheKey(environmentID, featureID, featureVersion)
	w.experimentNotFoundCache.Delete(key)
	w.logger.Debug("Cleared experiment missing cache entry",
		zap.String("environmentId", environmentID),
		zap.String("featureId", featureID),
		zap.Int32("featureVersion", featureVersion),
	)
}
