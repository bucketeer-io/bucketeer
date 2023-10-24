// Copyright 2023 The Bucketeer Authors.
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
	flightgroup      singleflight.Group
	location         *time.Location
	logger           *zap.Logger
}

func NewEvalEventWriter(
	ctx context.Context,
	r metrics.Registerer,
	l *zap.Logger,
	exClient ec.Client,
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
		for id, event := range events {
			switch evt := event.(type) {
			case *eventproto.EvaluationEvent:
				e, retriable, err := w.convToEvaluationEvent(ctx, evt, id, environmentNamespace)
				if err != nil {
					// If there is nothing to link, we don't report it as an error
					handledCounter.WithLabelValues(codeNoLink).Inc()
					if err == ErrNoExperiments ||
						err == ErrExperimentNotFound ||
						err == ErrGoalEventIssuedAfterExperimentEnded {
						w.logger.Debug(
							"There is no experiment to link",
							zap.Error(err),
							zap.String("id", id),
							zap.String("environmentNamespace", environmentNamespace),
							zap.Any("evalEvent", evt),
						)
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
) (*epproto.EvaluationEvent, bool, error) {
	experiments, err := w.listExperiments(ctx, environmentNamespace)
	if err != nil {
		w.logger.Error("failed to list experiments",
			zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.Any("evalEvent", e),
		)
		return nil, true, err
	}
	if len(experiments) == 0 {
		return nil, false, ErrNoExperiments
	}
	exp := w.existExperiment(experiments, e.FeatureId, e.FeatureVersion)
	if exp == nil {
		return nil, false, ErrExperimentNotFound
	}
	if exp.StopAt < e.Timestamp {
		handledCounter.WithLabelValues(codeEventIssuedAfterExperimentEnded).Inc()
		return nil, false, ErrGoalEventIssuedAfterExperimentEnded
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
			experiments := []*exproto.Experiment{}
			cursor := ""
			for {
				resp, err := w.experimentClient.ListExperiments(ctx, &exproto.ListExperimentsRequest{
					PageSize:             listRequestSize,
					Cursor:               cursor,
					EnvironmentNamespace: environmentNamespace,
					Statuses: []exproto.Experiment_Status{
						exproto.Experiment_RUNNING,
						exproto.Experiment_STOPPED,
					},
				})
				if err != nil {
					return nil, err
				}
				experiments = append(experiments, resp.Experiments...)
				experimentSize := len(resp.Experiments)
				if experimentSize == 0 || experimentSize < listRequestSize {
					return experiments, nil
				}
				cursor = resp.Cursor
			}
		},
	)
	if err != nil {
		handledCounter.WithLabelValues(codeFailedToListExperiments).Inc()
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
