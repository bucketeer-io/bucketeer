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

	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	ecstorage "github.com/bucketeer-io/bucketeer/pkg/eventcounter/storage/v2"
	ec "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	ft "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/writer"
	"github.com/bucketeer-io/bucketeer/pkg/subscriber/storage"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/proto/eventpersisterdwh"
	exproto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

const (
	goalEventTable = "goal_event"
	retryInterval  = 1 * time.Minute // interval to retry pending goals
	maxDelay       = 6 * time.Hour   // max delay before giving up
)

// pendingGoal holds a goal event waiting for its evaluation
// to arrive in BigQuery.
type pendingGoal struct {
	event       *eventproto.GoalEvent
	id          string
	environment string
	receivedAt  time.Time
	experiments []*exproto.Experiment
}

type goalEvtWriter struct {
	writer           storage.GoalEventWriter
	eventStorage     ecstorage.EventStorage
	experimentClient ec.Client
	featureClient    ft.Client
	cache            cachev3.ExperimentsCache
	flightgroup      singleflight.Group
	location         *time.Location
	logger           *zap.Logger

	pendingMu sync.Mutex
	pending   map[string]pendingGoal
}

// NewGoalEventWriter initializes the writer and starts the retry loop.
func NewGoalEventWriter(
	ctx context.Context,
	logger *zap.Logger,
	eventStorage ecstorage.EventStorage,
	exClient ec.Client,
	ftClient ft.Client,
	cache cachev3.ExperimentsCache,
	project, dataSet string,
	size int,
	location *time.Location,
) (Writer, error) {
	evt := epproto.GoalEvent{}
	goalWriter, err := writer.NewWriter(
		ctx,
		project,
		dataSet,
		goalEventTable,
		evt.ProtoReflect().Descriptor(),
		writer.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}
	w := &goalEvtWriter{
		writer:           storage.NewGoalEventWriter(goalWriter, size),
		eventStorage:     eventStorage,
		experimentClient: exClient,
		featureClient:    ftClient,
		cache:            cache,
		location:         location,
		logger:           logger,
		pending:          make(map[string]pendingGoal),
	}
	// Start background retry loop
	go w.retryPending(ctx)
	return w, nil
}

// retryPending periodically retries linking pending goals.
func (w *goalEvtWriter) retryPending(ctx context.Context) {
	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			w.processPending(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// processPending attempts to link all pending goals.
func (w *goalEvtWriter) processPending(ctx context.Context) {
	w.pendingMu.Lock()
	pending := w.pending
	w.pending = make(map[string]pendingGoal)
	w.pendingMu.Unlock()

	for _, pg := range pending {
		// Give up if too old
		if time.Since(pg.receivedAt) > maxDelay {
			w.logger.Warn("giving up on pending goal after max delay", zap.String("id", pg.id))
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeUserEvaluationNotFound).Inc()
			continue
		}
		// Retry linking
		evals, retriable, err := w.linkGoalEventByExperiment(ctx, pg.event, pg.environment, pg.experiments)
		if err != nil && retriable {
			w.enqueuePending(pg)
			continue
		}
		if err != nil {
			w.logger.Error("failed to link pending goal", zap.Error(err), zap.String("id", pg.id))
			subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeFailedToAppendGoalEvents).Inc()
			continue
		}
		// Convert and write
		events := make([]*epproto.GoalEvent, 0, len(evals))
		for _, eval := range evals {
			e, _, err := w.convToGoalEvent(pg.event, eval, pg.id, pg.event.Tag, pg.environment)
			if err != nil {
				w.logger.Error("failed to conv pending goal event", zap.Error(err), zap.String("id", pg.id))
				continue
			}
			events = append(events, e)
		}
		if _, err := w.writer.AppendRows(ctx, events); err != nil {
			w.logger.Error("failed to append pending goal events", zap.Error(err))
		}
		subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeLinked).Inc()
	}
}

// enqueuePending adds a goal event to the pending store.
func (w *goalEvtWriter) enqueuePending(pg pendingGoal) {
	w.pendingMu.Lock()
	defer w.pendingMu.Unlock()
	w.pending[pg.id] = pg
}

// Write processes incoming events and enqueues or writes them.
func (w *goalEvtWriter) Write(
	ctx context.Context,
	envEvents environmentEventDWHMap,
) map[string]bool {
	var goalEvents []*epproto.GoalEvent
	fails := make(map[string]bool, len(envEvents))
	for environmentId, events := range envEvents {
		experiments, err := w.listExperiments(ctx, environmentId)
		if err != nil {
			// Retry all events next time
			for id := range events {
				fails[id] = true
			}
			continue
		}
		for id, evt := range events {
			ge, ok := evt.(*eventproto.GoalEvent)
			if !ok {
				fails[id] = false
				continue
			}
			// Attempt immediate link
			evals, retriable, err := w.linkGoalEvent(ctx, ge, environmentId, experiments)
			if err != nil {
				if errors.Is(err, ecstorage.ErrNoResultsFound) {
					// Enqueue for retry
					w.enqueuePending(pendingGoal{
						event:       ge,
						id:          id,
						environment: environmentId,
						receivedAt:  time.Now(),
						experiments: experiments,
					})
					subscriberHandledCounter.WithLabelValues(subscriberGoalEventDWH, codeUserEvaluationNotFound).Inc()
					continue
				}
				fails[id] = retriable
				continue
			}
			// Convert and buffer for write
			for _, eval := range evals {
				e, retriable, err := w.convToGoalEvent(ge, eval, id, ge.Tag, environmentId)
				if err != nil {
					fails[id] = retriable
					break
				}
				goalEvents = append(goalEvents, e)
			}
		}
	}
	// Bulk write what we have
	if fs, err := w.writer.AppendRows(ctx, goalEvents); err != nil {
		w.logger.Error("Failed to append rows to goal_event table", zap.Error(err))
	} else {
		for id, f := range fs {
			fails[id] = f
		}
	}
	return fails
}

// Convert one or more goal events
// func (w *goalEvtWriter) convToGoalEvents(
// 	ctx context.Context,
// 	e *eventproto.GoalEvent,
// 	id, environmentID string,
// 	experiments []*exproto.Experiment,
// ) ([]*epproto.GoalEvent, bool, error) {
// 	evals, retriable, err := w.linkGoalEvent(ctx, e, environmentID, experiments)
// 	if err != nil {
// 		return nil, retriable, err
// 	}
// 	events := make([]*epproto.GoalEvent, 0, len(evals))
// 	for _, eval := range evals {
// 		event, retriable, err := w.convToGoalEvent(e, eval, id, e.Tag, environmentID)
// 		if err != nil {
// 			return nil, retriable, err
// 		}
// 		events = append(events, event)
// 	}
// 	return events, false, nil
// }

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
	environmentID string,
	experiments []*exproto.Experiment,
) ([]*ecstorage.UserEvaluation, bool, error) {
	evalExp, retriable, err := w.linkGoalEventByExperiment(ctx, event, environmentID, experiments)
	if err != nil {
		return nil, retriable, err
	}
	return evalExp, false, nil
}

// Link one or more experiments by goal ID
func (w *goalEvtWriter) linkGoalEventByExperiment(
	ctx context.Context,
	event *eventproto.GoalEvent,
	environmentID string,
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
		ev, err := w.getUserEvaluation(
			ctx,
			environmentID,
			getUserID(event.UserId, event.User),
			exp.FeatureId,
			exp.FeatureVersion,
			exp.StartAt,
			exp.StopAt,
			event.Timestamp,
		)
		if err != nil {
			w.logger.Error("Failed to get user evaluation",
				zap.Error(err),
				zap.String("environmentId", environmentID),
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

func (w *goalEvtWriter) getUserEvaluation(
	ctx context.Context,
	environmentID, userID, featureID string,
	featureVersion int32,
	experimentStartAt, experimentEndAt, goalTimestamp int64,
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
