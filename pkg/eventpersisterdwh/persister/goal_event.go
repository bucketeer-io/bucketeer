// Copyright 2022 The Bucketeer Authors.
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
	featurestorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/writer"
	bigtable "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigtable"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/proto/eventpersisterdwh"
	exproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const goalEventTable = "goal_event"

type goalEvtWriter struct {
	writer                storage.GoalEventWriter
	userEvaluationStorage featurestorage.UserEvaluationsStorage
	experimentClient      ec.Client
	flightgroup           singleflight.Group
	logger                *zap.Logger
}

func NewGoalEventWriter(
	ctx context.Context,
	r metrics.Registerer,
	l *zap.Logger,
	userEvaluationStorage featurestorage.UserEvaluationsStorage,
	exClient ec.Client,
	project, ds string,
	size int,
) (Writer, error) {
	evt := epproto.GoalEvent{}
	goalWriter, err := writer.NewWriter(
		ctx,
		project,
		ds,
		goalEventTable,
		evt.ProtoReflect().Descriptor(),
		writer.WithMetrics(r),
		writer.WithLogger(l),
	)
	if err != nil {
		return nil, err
	}
	return &goalEvtWriter{
		writer:                storage.NewGoalEventWriter(goalWriter, size),
		userEvaluationStorage: userEvaluationStorage,
		experimentClient:      exClient,
		logger:                l,
	}, nil
}

func (w *goalEvtWriter) Write(
	ctx context.Context,
	envEvents environmentEventMap,
) map[string]bool {
	goalEvents := []*epproto.GoalEvent{}
	fails := make(map[string]bool, len(envEvents))
	for environmentNamespace, events := range envEvents {
		for id, event := range events {
			switch evt := event.(type) {
			case *eventproto.GoalEvent:
				e, retriable, err := w.convToGoalEvent(ctx, evt, id, environmentNamespace)
				if err != nil {
					if err == ErrNoExperiments || err == ErrExperimentNotFound {
						w.logger.Debug(
							"There is no experiment to link",
							zap.Error(err),
							zap.String("id", id),
							zap.String("environmentNamespace", environmentNamespace),
						)
						continue
					}
					if !retriable {
						w.logger.Error(
							"Failed to convert to goal event",
							zap.Error(err),
							zap.String("id", id),
							zap.String("environmentNamespace", environmentNamespace),
						)
					}
					fails[id] = retriable
					continue
				}
				goalEvents = append(goalEvents, e)
			default:
				w.logger.Error(
					"The event is an unexpected message type",
					zap.String("id", id),
					zap.String("environmentNamespace", environmentNamespace),
				)
				fails[id] = false
			}
		}
	}
	fs, err := w.writer.AppendRows(ctx, goalEvents)
	if err != nil {
		w.logger.Error(
			"failed to append rows to goal event",
			zap.Error(err),
		)
	}
	for id, f := range fs {
		fails[id] = f
	}
	return fails
}

func (w *goalEvtWriter) convToGoalEvent(
	ctx context.Context,
	e *eventproto.GoalEvent,
	id, environmentNamespace string,
) (*epproto.GoalEvent, bool, error) {
	tag := e.Tag
	if tag == "" {
		// For requests with no tag, it will insert "none" instead, until all old SDK clients are updated
		tag = "none"
	}
	eval, retriable, err := w.linkGoalEvent(ctx, e, environmentNamespace, tag)
	if err != nil {
		return nil, retriable, err
	}
	var ud []byte
	if e.User != nil {
		var err error
		ud, err = json.Marshal(e.User.Data)
		if err != nil {
			return nil, false, err
		}
	}
	reason := ""
	if eval.Reason != nil {
		reason = eval.Reason.Type.String()
	}
	return &epproto.GoalEvent{
		Id:                   id,
		GoalId:               e.GoalId,
		Value:                float32(e.Value),
		UserData:             string(ud),
		UserId:               e.UserId,
		Tag:                  tag,
		SourceId:             e.SourceId.String(),
		EnvironmentNamespace: environmentNamespace,
		Timestamp:            time.Unix(e.Timestamp, 0).UnixMicro(),
		FeatureId:            eval.FeatureId,
		FeatureVersion:       eval.FeatureVersion,
		VariationId:          eval.VariationId,
		Reason:               reason,
	}, false, nil
}

func (w *goalEvtWriter) linkGoalEvent(
	ctx context.Context,
	event *eventproto.GoalEvent,
	environmentNamespace, tag string,
) (*featureproto.Evaluation, bool, error) {
	evalExp, retriable, err := w.linkGoalEventByExperiment(ctx, event, environmentNamespace, tag)
	if err != nil {
		return nil, retriable, err
	}
	return evalExp, false, nil
}

func (w *goalEvtWriter) linkGoalEventByExperiment(
	ctx context.Context,
	event *eventproto.GoalEvent,
	environmentNamespace, tag string,
) (*featureproto.Evaluation, bool, error) {
	// List experiments
	experiments, err := w.listExperiments(ctx, environmentNamespace)
	if err != nil {
		return nil, true, err
	}
	if len(experiments) == 0 {
		return nil, false, ErrNoExperiments
	}
	// Find the experiment by goal ID
	// TODO: we must change the console UI not to allow creating
	// multiple experiments running at the same time,
	// using the same feature flag id and goal id
	var experiment *exproto.Experiment
	for _, exp := range experiments {
		if w.findGoalID(event.GoalId, exp.GoalIds) {
			experiment = exp
			break
		}
	}
	if experiment == nil {
		return nil, false, ErrExperimentNotFound
	}
	// Get the user evaluation using the experiment info
	ev, err := w.getUserEvaluation(
		ctx,
		environmentNamespace,
		event.UserId,
		tag,
		experiment.FeatureId,
		experiment.FeatureVersion,
	)
	if err != nil {
		return nil, true, err
	}
	return ev, false, nil
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
	return exp.([]*exproto.Experiment), nil
}

func (w *goalEvtWriter) getUserEvaluation(
	ctx context.Context,
	environmentNamespace,
	userID,
	tag,
	featureID string,
	featureVersion int32,
) (*featureproto.Evaluation, error) {
	evaluation, err := w.userEvaluationStorage.GetUserEvaluation(
		ctx,
		userID,
		environmentNamespace,
		tag,
		featureID,
		featureVersion,
	)
	if err != nil {
		if err == bigtable.ErrKeyNotFound {
			handledCounter.WithLabelValues(codeUserEvaluationNotFound).Inc()
		} else {
			handledCounter.WithLabelValues(codeFailedToGetUserEvaluation).Inc()
		}
		return nil, err
	}
	return evaluation, nil
}
