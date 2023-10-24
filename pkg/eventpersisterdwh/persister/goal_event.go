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
	ft "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/writer"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/proto/eventpersisterdwh"
	exproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

const goalEventTable = "goal_event"

type goalEvtWriter struct {
	writer           storage.GoalEventWriter
	experimentClient ec.Client
	featureClient    ft.Client
	flightgroup      singleflight.Group
	location         *time.Location
	logger           *zap.Logger
}

func NewGoalEventWriter(
	ctx context.Context,
	r metrics.Registerer,
	l *zap.Logger,
	exClient ec.Client,
	ftClient ft.Client,
	project, ds string,
	size int,
	location *time.Location,
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
		writer:           storage.NewGoalEventWriter(goalWriter, size),
		experimentClient: exClient,
		featureClient:    ftClient,
		location:         location,
		logger:           l,
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
				e, retriable, err := w.convToGoalEvents(ctx, evt, id, environmentNamespace)
				if err != nil {
					if err == ErrNoExperiments || err == ErrExperimentNotFound {
						// If there is nothing to link, we don't report it as an error
						handledCounter.WithLabelValues(codeNoLink).Inc()
						w.logger.Debug(
							"There is no experiment to link",
							zap.Error(err),
							zap.String("id", id),
							zap.String("environmentNamespace", environmentNamespace),
							zap.Any("goalEvent", evt),
						)
						continue
					}
					if !retriable {
						w.logger.Error(
							"Failed to convert to goal event",
							zap.Error(err),
							zap.String("id", id),
							zap.String("environmentNamespace", environmentNamespace),
							zap.Any("goalEvent", evt),
						)
					}
					fails[id] = retriable
					continue
				}
				goalEvents = append(goalEvents, e...)
				handledCounter.WithLabelValues(codeLinked).Inc()
			default:
				w.logger.Error(
					"The event is an unexpected message type",
					zap.String("id", id),
					zap.String("environmentNamespace", environmentNamespace),
					zap.Any("goalEvent", evt),
				)
				fails[id] = false
			}
		}
	}
	fs, err := w.writer.AppendRows(ctx, goalEvents)
	if err != nil {
		handledCounter.WithLabelValues(codeFailedToAppendGoalEvents).Inc()
		w.logger.Error(
			"failed to append rows to goal event",
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
		w.logger.Error(
			"failed to append rows in the bigquery",
			zap.Any("goalEvents", failedToAppendMap),
		)
	}
	return fails
}

// Convert one or more goal events
func (w *goalEvtWriter) convToGoalEvents(
	ctx context.Context,
	e *eventproto.GoalEvent,
	id, environmentNamespace string,
) ([]*epproto.GoalEvent, bool, error) {
	evals, retriable, err := w.linkGoalEvent(ctx, e, environmentNamespace, e.Tag)
	if err != nil {
		return nil, retriable, err
	}
	events := make([]*epproto.GoalEvent, 0, len(evals))
	for _, eval := range evals {
		event, retriable, err := w.convToGoalEvent(ctx, e, eval, id, e.Tag, environmentNamespace)
		if err != nil {
			return nil, retriable, err
		}
		events = append(events, event)
	}
	return events, false, nil
}

func (w *goalEvtWriter) convToGoalEvent(
	ctx context.Context,
	e *eventproto.GoalEvent,
	eval *featureproto.Evaluation,
	id, tag, environmentNamespace string,
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
	if tag == "" {
		// Tag is optional, so we insert none when is empty.
		tag = "none"
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
		Reason:               eval.Reason.Type.String(),
	}, false, nil
}

func (w *goalEvtWriter) linkGoalEvent(
	ctx context.Context,
	event *eventproto.GoalEvent,
	environmentNamespace, tag string,
) ([]*featureproto.Evaluation, bool, error) {
	evalExp, retriable, err := w.linkGoalEventByExperiment(ctx, event, environmentNamespace, tag)
	if err != nil {
		return nil, retriable, err
	}
	return evalExp, false, nil
}

// Link one or more experiments by goal ID
func (w *goalEvtWriter) linkGoalEventByExperiment(
	ctx context.Context,
	event *eventproto.GoalEvent,
	environmentNamespace, tag string,
) ([]*featureproto.Evaluation, bool, error) {
	// List experiments
	experiments, err := w.listExperiments(ctx, environmentNamespace)
	if err != nil {
		w.logger.Error("failed to list experiments",
			zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.Any("goalEvent", event),
		)
		return nil, true, err
	}
	if len(experiments) == 0 {
		return nil, false, ErrNoExperiments
	}
	// Find the experiment by goal ID
	// TODO: we must change the console UI not to allow creating
	// multiple experiments running at the same time,
	// using the same feature flag id and goal id
	exps := []*exproto.Experiment{}
	for _, exp := range experiments {
		if w.findGoalID(event.GoalId, exp.GoalIds) {
			// If the goal event was issued before the experiment started running,
			// we ignore those events to avoid issues in the conversion rate
			if exp.StartAt > event.Timestamp {
				handledCounter.WithLabelValues(codeEventOlderThanExperiment).Inc()
				continue
			}
			// If the goal event was issued after the experiment ended,
			// we ignore those events to avoid issues in the conversion rate
			if exp.StopAt < event.Timestamp {
				handledCounter.WithLabelValues(codeGoalEventIssuedAfterExperimentEnded).Inc()
				continue
			}
			exps = append(exps, exp)
		}
	}
	if len(exps) == 0 {
		return nil, false, ErrExperimentNotFound
	}
	evals := make([]*featureproto.Evaluation, 0, len(exps))
	for _, exp := range exps {
		// Get the user evaluation using the experiment info
		ev, err := w.getUserEvaluation(
			ctx,
			event.User,
			environmentNamespace,
			tag,
			exp.FeatureId,
			exp.FeatureVersion,
		)
		if err != nil {
			if err == ErrEvaluationsAreEmpty {
				w.logger.Error("evaluations are empty",
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
					zap.Any("goalEvent", event),
				)
				return nil, false, err
			}
			w.logger.Error("failed to get user evaluation",
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
				zap.Any("goalEvent", event),
			)
			return nil, true, err
		}
		evals = append(evals, ev)
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

// TODO: Evaluate the user based on Feature Flag ID and version.
// By evaluating the user using the latest feature version,
// it could affect the experiment conversion accuracy
func (w *goalEvtWriter) getUserEvaluation(
	ctx context.Context,
	user *userproto.User,
	environmentNamespace, tag, featureID string,
	featureVersion int32,
) (*featureproto.Evaluation, error) {
	resp, err := w.featureClient.EvaluateFeatures(ctx, &featureproto.EvaluateFeaturesRequest{
		EnvironmentNamespace: environmentNamespace,
		FeatureId:            featureID,
		Tag:                  tag,
		User:                 user,
	})
	if err != nil {
		w.logger.Error(
			"Failed to evaluate user",
			zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("userId", user.Id),
			zap.String("featureId", featureID),
			zap.Int32("featureVersion", featureVersion),
			zap.String("tag", tag),
		)
		handledCounter.WithLabelValues(codeFailedToEvaluateUser).Inc()
		return nil, ErrFailedToEvaluateUser
	}
	if len(resp.UserEvaluations.Evaluations) == 0 {
		handledCounter.WithLabelValues(codeEvaluationsAreEmpty).Inc()
		w.logger.Error(
			"Evaluations are empty",
			zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("userId", user.Id),
			zap.String("featureId", featureID),
			zap.Int32("featureVersion", featureVersion),
			zap.String("tag", tag),
		)
		return nil, ErrEvaluationsAreEmpty
	}
	return resp.UserEvaluations.Evaluations[0], nil
}
