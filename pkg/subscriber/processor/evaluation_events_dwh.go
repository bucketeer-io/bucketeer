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
	ec "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	experimentdomain "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	v2hs "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2"
	redisv3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/dwhstorage"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/v2/proto/eventpersisterdwh"
	exproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

type evalEvtWriter struct {
	writer           dwhstorage.EvalEventWriter
	experimentClient ec.Client
	cache            cachev3.ExperimentsCache
	holdoutExposures *holdoutExposureDeriver
	flightgroup      singleflight.Group
	location         *time.Location
	subscriberType   string
	logger           *zap.Logger
}

// NewEvalEventWriter wraps an already-initialized data-warehouse evaluation-event writer
// (built in the server) with experiment-linking enrichment.
//
// holdoutStorage and measuredGroupsCache give the holdout groups whose
// internal experiments are in the experiments list, and holdoutLockClient is
// the Redis instance that keeps the rows derived for them unique.
func NewEvalEventWriter(
	l *zap.Logger,
	exClient ec.Client,
	cache cachev3.ExperimentsCache,
	holdoutStorage v2hs.HoldoutStorage,
	measuredGroupsCache cachev3.HoldoutMeasuredGroupsCache,
	holdoutLockClient redisv3.Client,
	location *time.Location,
	evalWriter dwhstorage.EvalEventWriter,
	subscriberType string,
) Writer {
	return &evalEvtWriter{
		writer:           evalWriter,
		experimentClient: exClient,
		cache:            cache,
		holdoutExposures: newHoldoutExposureDeriver(holdoutStorage, measuredGroupsCache, holdoutLockClient, l),
		location:         location,
		subscriberType:   subscriberType,
		logger:           l,
	}
}

func (w *evalEvtWriter) Write(
	ctx context.Context,
	envEvents environmentEventDWHMap,
) map[string]bool {
	var evalEvents []*epproto.EvaluationEvent
	// exposures are the rows in evalEvents that were derived for a holdout
	// group, by row ID.
	exposures := make(map[string]*holdoutExposure)
	exposureKeys := make(map[string]struct{})
	fails := make(map[string]bool, len(envEvents))
	for environmentId, events := range envEvents {
		experiments, err := w.listExperiments(ctx, environmentId)
		if err != nil {
			w.logger.Error("failed to list experiments",
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)
			subscriberHandledCounter.WithLabelValues(w.subscriberType, codeFailedToListExperiments).Inc()
			// Make sure to retry all the events in the next pulling
			for id := range events {
				fails[id] = true
			}
			continue
		}
		if len(experiments) == 0 {
			continue
		}
		experiments, measuredGroups, failureCode, err := w.listMeasuredHoldoutGroups(ctx, environmentId, experiments)
		if err != nil {
			w.logger.Error("failed to list measured holdout groups",
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)
			subscriberHandledCounter.WithLabelValues(w.subscriberType, failureCode).Inc()
			// Make sure to retry all the events in the next pulling
			for id := range events {
				fails[id] = true
			}
			continue
		}
		for id, event := range events {
			switch evt := event.(type) {
			case *eventproto.EvaluationEvent:
				// If there is nothing to link, we don't report it as an error
				linkErr := w.checkExperiment(evt, experiments)
				switch {
				case errors.Is(linkErr, ErrExperimentNotFound):
					subscriberHandledCounter.WithLabelValues(w.subscriberType, codeExperimentNotFound).Inc()
				case errors.Is(linkErr, ErrEvaluationEventIssuedAfterExperimentEnded):
					subscriberHandledCounter.WithLabelValues(w.subscriberType, codeEventIssuedAfterExperimentEnded).Inc()
				case errors.Is(linkErr, ErrEvaluationEventHeldOut):
					subscriberHandledCounter.WithLabelValues(w.subscriberType, codeEvaluationEventHeldOut).Inc()
				}
				match := w.holdoutExposures.match(evt, measuredGroups)
				if linkErr != nil && match == nil {
					continue
				}
				e, err := w.convToEvaluationEvent(evt, id, environmentId)
				if err != nil {
					w.logger.Error(
						"Failed to convert to evaluation event",
						zap.Error(err),
						zap.String("id", id),
						zap.String("environmentId", environmentId),
						zap.Any("evalEvent", evt),
					)
					fails[id] = false
					continue
				}
				var exposure *holdoutExposure
				if match != nil {
					exposure, err = w.holdoutExposures.acquire(ctx, match, id, e, exposureKeys)
					if err != nil {
						// Neither row of the event is written, so the next
						// delivery starts from the same state.
						w.logger.Error("Failed to deduplicate a holdout exposure",
							zap.Error(err),
							zap.String("id", id),
							zap.String("environmentId", environmentId),
							zap.String("holdoutGroupId", match.groupID),
						)
						subscriberHandledCounter.WithLabelValues(
							w.subscriberType, codeFailedToDeduplicateHoldoutExposure,
						).Inc()
						fails[id] = true
						continue
					}
				}
				if linkErr == nil {
					evalEvents = append(evalEvents, e)
					subscriberHandledCounter.WithLabelValues(w.subscriberType, codeLinked).Inc()
				}
				if exposure != nil {
					evalEvents = append(evalEvents, exposure.row)
					exposures[exposure.row.Id] = exposure
					subscriberHandledCounter.WithLabelValues(w.subscriberType, codeHoldoutExposureLinked).Inc()
				}
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
		subscriberHandledCounter.WithLabelValues(w.subscriberType, codeFailedToAppendEvaluationEvents).Inc()
		w.logger.Error("Failed to append rows to evaluation_event table",
			zap.Error(err),
		)
		if len(fs) == 0 {
			for _, ee := range evalEvents {
				fails[w.failedEventID(ctx, ee.Id, exposures)] = true
			}
			return fails
		}
	}
	failedToAppendMap := make(map[string]*epproto.EvaluationEvent)
	for id, f := range fs {
		for _, ee := range evalEvents {
			if id == ee.Id {
				failedToAppendMap[id] = ee
			}
		}
		// An event is retried when one of its rows can be retried.
		eventID := w.failedEventID(ctx, id, exposures)
		fails[eventID] = fails[eventID] || f
	}
	if len(failedToAppendMap) > 0 {
		w.logger.Error("Failed to append evaluation events",
			zap.Any("evaluationEvents", failedToAppendMap),
		)
	}
	return fails
}

// listMeasuredHoldoutGroups returns the measured holdout groups of the
// environment and the experiments list to use for the batch, which is
// experiments unless it is older than the groups.
//
// The two lists are cached apart. When a group names an internal experiment
// that experiments lacks, for example the experiment of an analysis period
// that just started, an event would be given the intervals of the listed
// experiments only, and could take the key of an earlier row of the user
// where the missing experiment needs a row of its own. The experiments are
// then read from the experiment service, which has what the group names. On
// an error, the returned code is the one to count it under.
func (w *evalEvtWriter) listMeasuredHoldoutGroups(
	ctx context.Context,
	environmentID string,
	experiments []*exproto.Experiment,
) ([]*exproto.Experiment, []*measuredHoldoutGroup, string, error) {
	measured, err := w.holdoutExposures.measuredGroups(ctx, environmentID, experiments)
	if err != nil {
		return nil, nil, codeFailedToListMeasuredHoldoutGroups, err
	}
	if !lacksInternalExperiment(measured) {
		return experiments, measured, "", nil
	}
	subscriberHandledCounter.WithLabelValues(w.subscriberType, codeExperimentsOlderThanHoldoutGroup).Inc()
	experiments, err = w.fetchExperiments(ctx, environmentID)
	if err != nil {
		return nil, nil, codeFailedToListExperiments, err
	}
	measured, err = w.holdoutExposures.measuredGroups(ctx, environmentID, experiments)
	if err != nil {
		return nil, nil, codeFailedToListMeasuredHoldoutGroups, err
	}
	if lacksInternalExperiment(measured) {
		// Nothing newer can be read. The batch goes on with what is listed.
		w.logger.Warn("A measured holdout group names an internal experiment that is not listed",
			zap.String("environmentId", environmentID),
		)
	}
	return experiments, measured, "", nil
}

// failedEventID returns the ID of the evaluation event that a row that was
// not written belongs to. The row written for a flag experiment has the ID
// of its event. A row derived for a holdout group has its own ID, and its
// key is released so that the row is written later.
func (w *evalEvtWriter) failedEventID(
	ctx context.Context,
	rowID string,
	exposures map[string]*holdoutExposure,
) string {
	exposure, ok := exposures[rowID]
	if !ok {
		return rowID
	}
	w.holdoutExposures.release(ctx, exposure)
	return exposure.eventID
}

// checkExperiment reports whether the event is written for the experiment of
// its flag. It returns nil when it is, and the reason when it is not.
func (w *evalEvtWriter) checkExperiment(
	e *eventproto.EvaluationEvent,
	experiments []*exproto.Experiment,
) error {
	exp := w.existExperiment(experiments, e.FeatureId, e.FeatureVersion)
	if exp == nil {
		return ErrExperimentNotFound
	}
	if exp.StopAt < e.Timestamp {
		return ErrEvaluationEventIssuedAfterExperimentEnded
	}
	// A held-out user was served the holdout variation whatever the flag
	// says, so the experiment of the flag does not count the user.
	if e.GetReason().GetType() == featureproto.Reason_HOLDOUT {
		return ErrEvaluationEventHeldOut
	}
	return nil
}

func (w *evalEvtWriter) convToEvaluationEvent(
	e *eventproto.EvaluationEvent,
	id, environmentId string,
) (*epproto.EvaluationEvent, error) {
	var ud []byte
	if e.User != nil {
		var err error
		userData := make(map[string]string)
		if e.User.Data != nil {
			userData = e.User.Data
		}
		ud, err = json.Marshal(userData)
		if err != nil {
			return nil, err
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
		Reason:         e.GetReason().GetType().String(),
		Tag:            tag,
		SourceId:       e.SourceId.String(),
		EnvironmentId:  environmentId,
		Timestamp:      time.Unix(e.Timestamp, 0).UnixMicro(),
	}, nil
}

// existExperiment returns the experiment of the flag version. The internal
// experiment of a holdout group is never one: its feature ID is the ID of
// the group, and its rows are derived, not written from an event of that ID.
func (w *evalEvtWriter) existExperiment(
	es []*exproto.Experiment,
	fID string,
	fVersion int32,
) *exproto.Experiment {
	for _, e := range es {
		if e.Kind == exproto.Experiment_HOLDOUT {
			continue
		}
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
			return w.listExperimentsFromService(ctx, environmentId)
		},
	)
	if err != nil {
		return nil, err
	}
	return experiments.([]*exproto.Experiment), nil
}

// fetchExperiments is listExperiments without the cache.
func (w *evalEvtWriter) fetchExperiments(
	ctx context.Context,
	environmentId string,
) ([]*exproto.Experiment, error) {
	experiments, err, _ := w.flightgroup.Do(
		fmt.Sprintf("%s:%s", environmentId, "fetchExperiments"),
		func() (interface{}, error) {
			return w.listExperimentsFromService(ctx, environmentId)
		},
	)
	if err != nil {
		return nil, err
	}
	return experiments.([]*exproto.Experiment), nil
}

func (w *evalEvtWriter) listExperimentsFromService(
	ctx context.Context,
	environmentId string,
) ([]*exproto.Experiment, error) {
	// Get the experiments from the DB
	resp, err := w.experimentClient.ListExperiments(ctx, &exproto.ListExperimentsRequest{
		// Because the evaluation and goal events may be sent with a delay
		// for many reasons from the client side, we still calculate
		// the results for two days after it stopped.
		StopAt:        time.Now().In(w.location).Add(-experimentdomain.StoppedExperimentEventWindow).Unix(),
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
