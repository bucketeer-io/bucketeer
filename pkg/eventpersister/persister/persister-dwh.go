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

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	aoclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/eventpersister/datastore"
	ec "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	featuredomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	featurestorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	bigtable "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigtable"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
	exproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	twentyFourHours = 24 * time.Hour
)

type PersisterDwh struct {
	experimentClient      ec.Client
	featureClient         featureclient.Client
	autoOpsClient         aoclient.Client
	puller                puller.RateLimitedPuller
	logger                *zap.Logger
	ctx                   context.Context
	cancel                func()
	group                 errgroup.Group
	doneCh                chan struct{}
	userEvaluationStorage featurestorage.UserEvaluationsStorage
	evalEventWriter       datastore.EvalEventWriter
	goalEventWriter       datastore.GoalEventWriter
	flightgroup           singleflight.Group
	opts                  *options
}

func NewPersisterDwh(
	p puller.Puller,
	evalEventWriter datastore.EvalEventWriter,
	goalEventWriter datastore.GoalEventWriter,
	opts ...Option,
) *PersisterDwh {
	dopts := &options{
		maxMPS:        1000,
		numWorkers:    1,
		flushSize:     50,
		flushInterval: 5 * time.Second,
		flushTimeout:  20 * time.Second,
		logger:        zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	if dopts.metrics != nil {
		dwhRegisterMetrics(dopts.metrics)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &PersisterDwh{
		puller:          puller.NewRateLimitedPuller(p, dopts.maxMPS),
		logger:          dopts.logger.Named("persister-dwh"),
		ctx:             ctx,
		cancel:          cancel,
		doneCh:          make(chan struct{}),
		evalEventWriter: evalEventWriter,
		goalEventWriter: goalEventWriter,
	}
}

func (p *PersisterDwh) Run() error {
	defer close(p.doneCh)
	p.group.Go(func() error {
		return p.puller.Run(p.ctx)
	})
	for i := 0; i < p.opts.numWorkers; i++ {
		p.group.Go(p.batch)
	}
	return p.group.Wait()
}

func (p *PersisterDwh) Stop() {
	p.cancel()
	<-p.doneCh
}

func (p *PersisterDwh) Check(ctx context.Context) health.Status {
	select {
	case <-p.ctx.Done():
		p.logger.Error("Unhealthy due to context Done is closed", zap.Error(p.ctx.Err()))
		return health.Unhealthy
	default:
		if p.group.FinishedCount() > 0 {
			p.logger.Error("Unhealthy", zap.Int32("FinishedCount", p.group.FinishedCount()))
			return health.Unhealthy
		}
		return health.Healthy
	}
}

func (p *PersisterDwh) batch() error {
	batch := make(map[string]*puller.Message)
	timer := time.NewTimer(p.opts.flushInterval)
	defer timer.Stop()
	for {
		select {
		case msg, ok := <-p.puller.MessageCh():
			if !ok {
				return nil
			}
			// dwhReceivedCounter.Inc()
			id := msg.Attributes["id"]
			if id == "" {
				msg.Ack()
				// TODO: better log format for msg data
				// dwhHandledCounter.WithLabelValues(codes.MissingID.String()).Inc()
				continue
			}
			if previous, ok := batch[id]; ok {
				previous.Ack()
				p.logger.Warn("Message with duplicate id", zap.String("id", id))
				// dwhHandledCounter.WithLabelValues(codes.DuplicateID.String()).Inc()
			}
			batch[id] = msg
			if len(batch) < p.opts.flushSize {
				continue
			}
			p.send(batch)
			batch = make(map[string]*puller.Message)
			timer.Reset(p.opts.flushInterval)
		case <-timer.C:
			if len(batch) > 0 {
				p.send(batch)
				batch = make(map[string]*puller.Message)
			}
			timer.Reset(p.opts.flushInterval)
		case <-p.ctx.Done():
			return nil
		}
	}
}

func (p *PersisterDwh) send(messages map[string]*puller.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), p.opts.flushTimeout)
	defer cancel()

	envEvents := p.extractEvents(messages)
	if len(envEvents) == 0 {
		p.logger.Error("all messages were bad")
		return
	}

	evalEvents := []*ecproto.EvaluationEvent{}
	goalEvents := []*ecproto.GoalEvent{}

	fails := make(map[string]bool, len(envEvents))
	for environmentNamespace, events := range envEvents {
		for id, event := range events {
			switch evt := event.(type) {
			case *eventproto.EvaluationEvent:
				e, retriable, err := p.convToEvaluationEvent(ctx, evt, id, environmentNamespace)
				if err != nil {
					if err == ErrNoExperiments {
						continue
					}
					if !retriable {
						p.logger.Error(
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
			case *eventproto.GoalEvent:
				e, retriable, err := p.convToGoalEvent(ctx, evt, id, environmentNamespace)
				if err != nil {
					if !retriable {
						p.logger.Error(
							"failed to convert to goal event",
							zap.Error(err),
							zap.String("id", id),
							zap.String("environmentNamespace", environmentNamespace),
						)
					}
					fails[id] = retriable
					continue
				}
				goalEvents = append(goalEvents, e)
			}
		}
	}
	if err := p.evalEventWriter.AppendRows(ctx, evalEvents); err != nil {
		p.logger.Error(
			"failed to append rows to evaluation event",
			zap.Error(err),
		)
	}
	if err := p.goalEventWriter.AppendRows(ctx, goalEvents); err != nil {
		p.logger.Error(
			"failed to append rows to goal event",
			zap.Error(err),
		)
	}
	for id, m := range messages {
		if repeatable, ok := fails[id]; ok {
			if repeatable {
				m.Nack()
				handledCounter.WithLabelValues(codes.RepeatableError.String()).Inc()
			} else {
				m.Ack()
				handledCounter.WithLabelValues(codes.NonRepeatableError.String()).Inc()
			}
			continue
		}
		m.Ack()
		handledCounter.WithLabelValues(codes.OK.String()).Inc()
	}
}

func (p *PersisterDwh) extractEvents(messages map[string]*puller.Message) environmentEventMap {
	envEvents := environmentEventMap{}
	handleBadMessage := func(m *puller.Message, err error) {
		m.Ack()
		p.logger.Error("bad message", zap.Error(err), zap.Any("msg", m))
		// handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
	}
	for _, m := range messages {
		event := &eventproto.Event{}
		if err := proto.Unmarshal(m.Data, event); err != nil {
			handleBadMessage(m, err)
			continue
		}
		var innerEvent ptypes.DynamicAny
		if err := ptypes.UnmarshalAny(event.Event, &innerEvent); err != nil {
			handleBadMessage(m, err)
			continue
		}
		if innerEvents, ok := envEvents[event.EnvironmentNamespace]; ok {
			innerEvents[event.Id] = innerEvent.Message
			continue
		}
		envEvents[event.EnvironmentNamespace] = eventMap{event.Id: innerEvent.Message}
	}
	return envEvents
}

func (p *PersisterDwh) convToEvaluationEvent(
	ctx context.Context,
	e *eventproto.EvaluationEvent,
	id, environmentNamespace string,
) (*ecproto.EvaluationEvent, bool, error) {
	if err := p.validateTimestamp(e.Timestamp); err != nil {
		handledCounter.WithLabelValues(codeInvalidGoalEventTimestamp).Inc()
		return nil, false, err
	}
	exist, err := p.existExperiment(ctx, e, environmentNamespace)
	if err != nil {
		return nil, true, err
	}
	if !exist {
		return nil, false, ErrNoExperiments
	}
	if err := p.upsertUserEvaluation(ctx, e, environmentNamespace); err != nil {
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
	return &ecproto.EvaluationEvent{
		Id:                   id,
		FeatureId:            e.FeatureId,
		FeatureVersion:       e.FeatureVersion,
		UserData:             string(ud),
		UserId:               e.UserId,
		VariationId:          e.VariationId,
		Reason:               e.Reason.Type.String(),
		Tag:                  e.Tag,
		SourceId:             e.SourceId.String(),
		EnvironmentNamespace: environmentNamespace,
		Timestamp:            e.Timestamp,
	}, false, nil
}

func (p *PersisterDwh) convToGoalEvent(
	ctx context.Context,
	e *eventproto.GoalEvent,
	id, environmentNamespace string,
) (*ecproto.GoalEvent, bool, error) {
	eval, retriable, err := p.linkGoalEvent(ctx, e, environmentNamespace)
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
	return &ecproto.GoalEvent{
		Id:                   id,
		GoalId:               e.GoalId,
		Value:                float32(e.Value),
		UserData:             string(ud),
		UserId:               e.UserId,
		Tag:                  e.Tag,
		SourceId:             e.SourceId.String(),
		EnvironmentNamespace: environmentNamespace,
		Timestamp:            e.Timestamp,
		FeatureId:            eval.FeatureId,
		FeatureVersion:       eval.FeatureVersion,
		VariationId:          eval.VariationId,
		Reason:               reason,
	}, false, nil
}

func (p *PersisterDwh) linkGoalEvent(
	ctx context.Context,
	event *eventproto.GoalEvent,
	environmentNamespace string,
) (*featureproto.Evaluation, bool, error) {
	if err := p.validateTimestamp(event.Timestamp); err != nil {
		handledCounter.WithLabelValues(codeInvalidGoalEventTimestamp).Inc()
		return nil, false, err
	}
	evalExp, retriable, err := p.linkGoalEventByExperiment(ctx, event, environmentNamespace)
	// If there are no experiments or the goal ID didn't match, it will ignore the error
	// so it can try to link the goal event to auto ops.
	if err != nil && err != ErrNoExperiments && err != ErrExperimentNotFound {
		return nil, retriable, err
	}
	if evalExp == nil {
		handledCounter.WithLabelValues(codeNothingToLink).Inc()
		return nil, false, ErrNothingToLink
	}
	return evalExp, false, nil
}

func (*PersisterDwh) validateTimestamp(
	timestamp int64,
) error {
	actual := time.Unix(timestamp, 0)
	now := time.Now()
	min := now.Add(-twentyFourHours)
	max := now.Add(twentyFourHours)
	if actual.Before(min) || actual.After(max) {
		return ErrInvalidEventTimestamp
	}
	return nil
}

func (p *PersisterDwh) getUserEvaluation(
	environmentNamespace,
	userID,
	tag,
	featureID string,
	featureVersion int32,
) (*featureproto.Evaluation, error) {
	// For requests with no tag, it will insert "none" instead, until all old SDK clients are updated
	if tag == "" {
		tag = "none"
	}
	evaluation, err := p.userEvaluationStorage.GetUserEvaluation(
		p.ctx,
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

func (p *PersisterDwh) linkGoalEventByExperiment(
	ctx context.Context,
	event *eventproto.GoalEvent,
	environmentNamespace string,
) (*featureproto.Evaluation, bool, error) {
	// List experiments with the following status RUNNING, FORCE_STOPPED, and STOPPED
	experiments, err := p.listExperiments(ctx, environmentNamespace)
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
		if p.findGoalID(event.GoalId, exp.GoalIds) {
			experiment = exp
			break
		}
	}
	if experiment == nil {
		return nil, false, ErrExperimentNotFound
	}
	// Get the user evaluation using the experiment info
	ev, err := p.getUserEvaluation(
		environmentNamespace,
		event.UserId,
		event.Tag,
		experiment.FeatureId,
		experiment.FeatureVersion,
	)
	if err != nil {
		return nil, true, err
	}
	return ev, false, nil
}

func (*PersisterDwh) findGoalID(id string, goalIDs []string) bool {
	for _, goalID := range goalIDs {
		if id == goalID {
			return true
		}
	}
	return false
}

func (p *PersisterDwh) listExperiments(
	ctx context.Context,
	environmentNamespace string,
) ([]*exproto.Experiment, error) {
	exp, err, _ := p.flightgroup.Do(
		fmt.Sprintf("%s:%s", environmentNamespace, "listExperiments"),
		func() (interface{}, error) {
			experiments := []*exproto.Experiment{}
			cursor := ""
			for {
				resp, err := p.experimentClient.ListExperiments(ctx, &exproto.ListExperimentsRequest{
					PageSize:             listRequestSize,
					Cursor:               cursor,
					EnvironmentNamespace: environmentNamespace,
					Statuses: []exproto.Experiment_Status{
						exproto.Experiment_RUNNING,
						exproto.Experiment_FORCE_STOPPED,
						exproto.Experiment_STOPPED,
					},
					Archived: &wrappers.BoolValue{Value: false},
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

func (p *PersisterDwh) convToEvaluation(
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

func (p *PersisterDwh) existExperiment(
	ctx context.Context,
	event *eventproto.EvaluationEvent,
	environmentNamespace string,
) (bool, error) {
	resp, err := p.experimentClient.ListExperiments(ctx, &exproto.ListExperimentsRequest{
		FeatureId:            event.FeatureId,
		FeatureVersion:       &wrappers.Int32Value{Value: event.FeatureVersion},
		PageSize:             1,
		EnvironmentNamespace: environmentNamespace,
		Statuses: []exproto.Experiment_Status{
			exproto.Experiment_RUNNING,
			exproto.Experiment_FORCE_STOPPED,
			exproto.Experiment_STOPPED,
		},
		Archived: &wrappers.BoolValue{Value: false},
	})
	if err != nil {
		return false, err
	}
	return len(resp.Experiments) == 1, nil
}

func (p *PersisterDwh) upsertUserEvaluation(
	ctx context.Context,
	event *eventproto.EvaluationEvent,
	environmentNamespace string,
) error {
	evaluation, tag := p.convToEvaluation(ctx, event)
	if err := p.userEvaluationStorage.UpsertUserEvaluation(
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
