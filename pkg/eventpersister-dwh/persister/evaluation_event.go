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

	"github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/eventpersister-dwh/datastore"
	ec "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featuredomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	featurestorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/writer"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/proto/eventpersister-dwh"
	exproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const evaluationEventTable = "evaluation_event"

type evalEvtWriter struct {
	writer                datastore.EvalEventWriter
	userEvaluationStorage featurestorage.UserEvaluationsStorage
	experimentClient      ec.Client
	logger                *zap.Logger
}

func NewEvalEventWriter(
	ctx context.Context,
	r metrics.Registerer,
	l *zap.Logger,
	project, ds string,
	size int,
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
		writer.WithBatchSize(size),
	)
	if err != nil {
		return nil, err
	}
	return &evalEvtWriter{
		writer: datastore.NewEvalEventWriter(evalQuery),
		logger: l,
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
					if err == ErrNoExperiments {
						w.logger.Warn(
							"There is no running experiments",
							zap.Error(err),
							zap.String("id", id),
							zap.String("environmentNamespace", environmentNamespace),
						)
						continue
					}
					if !retriable {
						w.logger.Error(
							"failed to convert to evaluation event",
							zap.Error(err),
							zap.String("id", id),
							zap.String("environmentNamespace", environmentNamespace),
						)
					}
					fails[id] = retriable
					continue
				}
				evalEvents = append(evalEvents, e)
			}
		}
	}
	if err := w.writer.AppendRows(ctx, evalEvents); err != nil {
		w.logger.Error(
			"failed to append rows to evaluation event",
			zap.Error(err),
		)
	}
	return fails
}

func (w *evalEvtWriter) convToEvaluationEvent(
	ctx context.Context,
	e *eventproto.EvaluationEvent,
	id, environmentNamespace string,
) (*epproto.EvaluationEvent, bool, error) {
	exist, err := w.existExperiment(ctx, e, environmentNamespace)
	if err != nil {
		return nil, true, err
	}
	if !exist {
		return nil, false, ErrNoExperiments
	}
	if err := w.upsertUserEvaluation(ctx, e, environmentNamespace); err != nil {
		return nil, true, err
	}
	var ud []byte
	if e.User != nil {
		var err error
		ud, err = json.Marshal(e.User.Data)
		if err != nil {
			return nil, false, err
		}
	}
	tag := e.Tag
	if tag == "" {
		// For requests with no tag, it will insert "none" instead, until all old SDK clients are updated
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
		Timestamp:            e.Timestamp,
	}, false, nil
}

func (w *evalEvtWriter) existExperiment(
	ctx context.Context,
	event *eventproto.EvaluationEvent,
	environmentNamespace string,
) (bool, error) {
	resp, err := w.experimentClient.ListExperiments(ctx, &exproto.ListExperimentsRequest{
		FeatureId:            event.FeatureId,
		FeatureVersion:       &wrappers.Int32Value{Value: event.FeatureVersion},
		PageSize:             1,
		EnvironmentNamespace: environmentNamespace,
		Statuses: []exproto.Experiment_Status{
			exproto.Experiment_RUNNING,
		},
	})
	if err != nil {
		return false, err
	}
	return len(resp.Experiments) == 1, nil
}

func (w *evalEvtWriter) upsertUserEvaluation(
	ctx context.Context,
	event *eventproto.EvaluationEvent,
	environmentNamespace string,
) error {
	evaluation, tag := w.convToEvaluation(ctx, event)
	if err := w.userEvaluationStorage.UpsertUserEvaluation(
		ctx,
		evaluation,
		environmentNamespace,
		tag,
	); err != nil {
		handledCounter.WithLabelValues(codeUpsertUserEvaluationFailed).Inc()
		return err
	}
	return nil
}

func (w *evalEvtWriter) convToEvaluation(
	ctx context.Context,
	event *eventproto.EvaluationEvent,
) (*featureproto.Evaluation, string) {
	evaluation := &featureproto.Evaluation{
		Id: featuredomain.EvaluationID(
			event.FeatureId,
			event.FeatureVersion,
			event.UserId,
		),
		FeatureId:      event.FeatureId,
		FeatureVersion: event.FeatureVersion,
		UserId:         event.UserId,
		VariationId:    event.VariationId,
		Reason:         event.Reason,
	}
	// For requests that doesn't have the tag info,
	// it will insert none instead, until all SDK clients are updated
	var tag string
	if event.Tag == "" {
		tag = "none"
	} else {
		tag = event.Tag
	}
	return evaluation, tag
}
