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
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	ec "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/writer"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/subscriber/storage"
	storagev2 "github.com/bucketeer-io/bucketeer/pkg/subscriber/storage/v2"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/proto/eventpersisterdwh"
	exproto "github.com/bucketeer-io/bucketeer/proto/experiment"
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

	switch option.DataWarehouseType {
	case "mysql":
		if option.MySQLClient == nil {
			return nil, errors.New("mysql client is required when using MySQL storage")
		}

		evalStorage := storagev2.NewMysqlEvaluationEventStorage(option.MySQLClient)

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
			continue
		}
		for id, event := range events {
			switch evt := event.(type) {
			case *eventproto.EvaluationEvent:
				e, retriable, err := w.convToEvaluationEvent(ctx, evt, id, environmentId, experiments)
				if err != nil {
					// If there is nothing to link, we don't report it as an error
					if errors.Is(err, ErrExperimentNotFound) {
						subscriberHandledCounter.WithLabelValues(subscriberEvaluationEventDWH, codeExperimentNotFound).Inc()
						continue
					}
					if errors.Is(err, ErrEvaluationEventIssuedAfterExperimentEnded) {
						subscriberHandledCounter.WithLabelValues(subscriberEvaluationEventDWH, codeEventIssuedAfterExperimentEnded).Inc()
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
	fs, err := w.writer.AppendRows(ctx, evalEvents)
	if err != nil {
		subscriberHandledCounter.WithLabelValues(subscriberEvaluationEventDWH, codeFailedToAppendEvaluationEvents).Inc()
		w.logger.Error("Failed to append rows to evaluation_event table",
			zap.Error(err),
		)
	}
	failedToAppendMap := make(map[string]*epproto.EvaluationEvent)
	for id, f := range fs {
		// To log which event has failed to append in the BigQuery, we need to find the event
		for _, ee := range evalEvents {
			if id == ee.Id {
				failedToAppendMap[id] = ee
			}
		}
		// Update the fails map
		fails[id] = f
	}
	if len(failedToAppendMap) > 0 {
		w.logger.Error("Failed to append evaluation events",
			zap.Any("evaluationEvents", failedToAppendMap),
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
	exp := w.existExperiment(experiments, e.FeatureId, e.FeatureVersion)
	if exp == nil {
		return nil, false, ErrExperimentNotFound
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
