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

package persister

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/eventpersisterdwh/storage"
	ec "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/writer"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/proto/eventpersisterdwh"
	exproto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

const evaluationEventTable = "evaluation_event"

type evalEvtWriter struct {
	writer           storage.EvalEventWriter
	experimentClient ec.Client
	cache            cachev3.ExperimentsCache
	flightgroup      singleflight.Group
	location         *time.Location
	logger           *zap.Logger
}

func NewEvalEventWriter(
	ctx context.Context,
	r metrics.Registerer,
	l *zap.Logger,
	exClient ec.Client,
	cache cachev3.ExperimentsCache,
	project, ds string,
	size int,
	location *time.Location,
) (Writer, error) {
	evt := epproto.EvaluationEvent{}
	evalQuery, err := writer.NewWriter(
		ctx,
		project,
		ds,
		evaluationEventTable,
		evt.ProtoReflect().Descriptor(),
		writer.WithMetrics(r),
		writer.WithLogger(l),
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
	envEvents environmentEventMap,
) map[string]bool {
	evalEvents := []*epproto.EvaluationEvent{}
	fails := make(map[string]bool, len(envEvents))
	for environmentNamespace, events := range envEvents {
		experiments, err := w.listExperiments(ctx, environmentNamespace)
		if err != nil {
			w.logger.Error("failed to list experiments",
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
			)
			handledCounter.WithLabelValues(codeFailedToListExperiments).Inc()
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
				e, retriable, err := w.convToEvaluationEvent(ctx, evt, id, environmentNamespace, experiments)
				if err != nil {
					// If there is nothing to link, we don't report it as an error
					if err == ErrExperimentNotFound {
						handledCounter.WithLabelValues(codeExperimentNotFound).Inc()
						continue
					}
					if err == ErrEvaluationEventIssuedAfterExperimentEnded {
						handledCounter.WithLabelValues(codeEventIssuedAfterExperimentEnded).Inc()
						continue
					}
					if !retriable {
						w.logger.Error(
							"Failed to convert to evaluation event",
							zap.Error(err),
							zap.String("id", id),
							zap.String("environmentNamespace", environmentNamespace),
							zap.Any("evalEvent", evt),
						)
					}
					fails[id] = retriable
					continue
				}
				evalEvents = append(evalEvents, e)
				handledCounter.WithLabelValues(codeLinked).Inc()
			default:
				w.logger.Error(
					"The event is an unexpected message type",
					zap.String("id", id),
					zap.String("environmentNamespace", environmentNamespace),
					zap.Any("evalEvent", evt),
				)
				fails[id] = false
			}
		}
	}
	fs, err := w.writer.AppendRows(ctx, evalEvents)
	if err != nil {
		handledCounter.WithLabelValues(codeFailedToAppendEvaluationEvents).Inc()
		w.logger.Error(
			"failed to append rows to evaluation event",
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
		w.logger.Error(
			"failed to append rows in the bigquery",
			zap.Any("goalEvents", failedToAppendMap),
		)
	}
	return fails
}

func (w *evalEvtWriter) convToEvaluationEvent(
	ctx context.Context,
	e *eventproto.EvaluationEvent,
	id, environmentNamespace string,
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
	tag := e.Tag
	if tag == "" {
		// Tag is optional, so we insert none when is empty.
		tag = "none"
	}
	return &epproto.EvaluationEvent{
		Id:                   id,
		FeatureId:            e.FeatureId,
		FeatureVersion:       e.FeatureVersion,
		UserData:             string(ud),
		UserId:               e.UserId,
		VariationId:          e.VariationId,
		Reason:               e.Reason.Type.String(),
		Tag:                  tag,
		SourceId:             e.SourceId.String(),
		EnvironmentNamespace: environmentNamespace,
		Timestamp:            time.Unix(e.Timestamp, 0).UnixMicro(),
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
	environmentNamespace string,
) ([]*exproto.Experiment, error) {
	exp, err, _ := w.flightgroup.Do(
		fmt.Sprintf("%s:%s", environmentNamespace, "listExperiments"),
		func() (interface{}, error) {
			// Get the experiment cache
			expList, err := w.cache.Get(environmentNamespace)
			if err == nil {
				return expList.Experiments, nil
			}
			// Get the experiments from the DB
			resp, err := w.experimentClient.ListExperiments(ctx, &exproto.ListExperimentsRequest{
				PageSize:             0,
				EnvironmentNamespace: environmentNamespace,
				Statuses: []exproto.Experiment_Status{
					exproto.Experiment_RUNNING,
					exproto.Experiment_STOPPED,
				},
			})
			if err != nil {
				return nil, err
			}
			// Cache the experiment for the next request
			experiments := &exproto.Experiments{
				Experiments: resp.Experiments,
			}
			if err := w.cache.Put(experiments, environmentNamespace); err != nil {
				w.logger.Error("Failed to cache experiments",
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)
			}
			return resp.Experiments, nil
		},
	)
	if err != nil {
		return nil, err
	}
	// Filter the stopped experiments
	// Because the evaluation and goal events may be sent with a delay for many reasons from the client side,
	// we still calculate the results for two days after it stopped.
	now := time.Now().In(w.location)
	experiments := make([]*exproto.Experiment, 0, len(exp.([]*exproto.Experiment)))
	for _, e := range exp.([]*exproto.Experiment) {
		if e.Status == exproto.Experiment_STOPPED && now.Unix()-e.StopAt > 2*day {
			continue
		}
		experiments = append(experiments, e)
	}
	return experiments, nil
}
