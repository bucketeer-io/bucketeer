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
	"sort"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	aoclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	aodomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/eventpersister/datastore"
	ec "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	aoproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
	exproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	twentyFourHours = 24 * time.Hour
)

type PersisterDwh struct {
	experimentClient ec.Client
	featureClient    featureclient.Client
	autoOpsClient    aoclient.Client
	puller           puller.RateLimitedPuller
	logger           *zap.Logger
	ctx              context.Context
	cancel           func()
	group            errgroup.Group
	doneCh           chan struct{}
	evalEventWriter  datastore.EvalEventWriter
	goalEventWriter  datastore.GoalEventWriter
	evalEventReader  datastore.EvalEventReader
	flightgroup      singleflight.Group
	opts             *options
}

func NewPersisterDwh(
	p puller.Puller,
	evalEventWriter datastore.EvalEventWriter,
	goalEventWriter datastore.GoalEventWriter,
	evalEventReader datastore.EvalEventReader,
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
		evalEventReader: evalEventReader,
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

	for environmentNamespace, events := range envEvents {
		for id, event := range events {
			if evt, ok := event.(*eventproto.EvaluationEvent); ok {
				e, err := p.convToEvaluationEvent(ctx, evt, id, environmentNamespace)
				if err != nil {
					p.logger.Error(
						"failed to convert to evaluation event",
						zap.Error(err),
						zap.String("id", id),
						zap.String("environmentNamespace", environmentNamespace),
					)
				}
				evalEvents = append(evalEvents, e)
			}
		}
	}
	if err := p.evalEventWriter.Write(ctx, evalEvents); err != nil {

	}
	for environmentNamespace, events := range envEvents {
		for id, event := range events {
			if evt, ok := event.(*eventproto.GoalEvent); ok {
				e, _, err := p.convToGoalEvent(ctx, evt, id, environmentNamespace)
				if err != nil {
					p.logger.Error(
						"failed to convert to evaluation event",
						zap.Error(err),
						zap.String("id", id),
						zap.String("environmentNamespace", environmentNamespace),
					)
				}
				goalEvents = append(goalEvents, e)
			}
		}
	}
	if err := p.goalEventWriter.Write(ctx, goalEvents); err != nil {

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

func (*PersisterDwh) convToEvaluationEvent(
	ctx context.Context,
	e *eventproto.EvaluationEvent,
	id, environmentNamespace string,
) (*ecproto.EvaluationEvent, error) {
	var ud []byte
	if e.User != nil {
		var err error
		ud, err = json.Marshal(e.User.Data)
		if err != nil {
			return nil, err
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
	}, nil
}

func (p *PersisterDwh) convToGoalEvent(
	ctx context.Context,
	e *eventproto.GoalEvent,
	id, environmentNamespace string,
) (*ecproto.GoalEvent, bool, error) {
	ev, retriable, err := p.linkGoalEvent(ctx, e, environmentNamespace)
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
	return &ecproto.GoalEvent{
		Id:                   id,
		GoalId:               e.GoalId,
		Value:                float32(e.Value),
		UserData:             string(ud),
		UserId:               e.UserId,
		Evaluation:           strings.Join(ev, ","),
		Tag:                  e.Tag,
		SourceId:             e.SourceId.String(),
		EnvironmentNamespace: environmentNamespace,
		Timestamp:            e.Timestamp,
	}, false, nil
}

func (p *PersisterDwh) linkGoalEvent(
	ctx context.Context,
	event *eventproto.GoalEvent,
	environmentNamespace string,
) ([]string, bool, error) {
	if err := p.validateTimestamp(event.Timestamp); err != nil {
		handledCounter.WithLabelValues(codeInvalidGoalEventTimestamp).Inc()
		return nil, false, err
	}
	evaluations := []*featureproto.Evaluation{}
	evalExp, retriable, err := p.linkGoalEventByExperiment(
		ctx,
		event,
		environmentNamespace,
	)
	if err != nil {
		return nil, retriable, err
	}
	if evalExp != nil {
		evaluations = append(evaluations, evalExp)
	}
	evalAuto, retriable, err := p.linkGoalEventByAutoOps(
		ctx,
		event,
		environmentNamespace,
	)
	if err != nil {
		return nil, retriable, err
	}
	if evalAuto != nil {
		evaluations = append(evaluations, evalAuto...)
	}
	if len(evaluations) == 0 {
		handledCounter.WithLabelValues(codeNothingToLink).Inc()
		return nil, false, ErrNothingToLink
	}
	evalSet := p.getEvalSet(evaluations)
	// Sort the slice to avoid errors in the unit tests
	sort.Sort(sort.Reverse(sort.StringSlice(evalSet)))
	return evalSet, false, nil
}

func (*PersisterDwh) validateTimestamp(
	timestamp int64,
) error {
	actual := time.Unix(timestamp, 0)
	now := time.Now()
	min := now.Add(-twentyFourHours)
	max := now.Add(twentyFourHours)
	if actual.Before(min) || actual.After(max) {
		return ErrInvalidGoalEventTimestamp
	}
	return nil
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
		return &featureproto.Evaluation{}, false, nil
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
		return &featureproto.Evaluation{}, false, nil
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

func (*PersisterDwh) getEvalSet(evals []*featureproto.Evaluation) []string {
	evalsMap := make(map[string]struct{})
	for _, ev := range evals {
		// Check the reason
		reason := ""
		if ev.Reason != nil {
			reason = ev.Reason.Type.String()
		}
		eval := fmt.Sprintf("%s:%d:%s:%s", ev.FeatureId, ev.FeatureVersion, ev.VariationId, reason)
		// Remove duplicates if needed
		evalsMap[eval] = struct{}{}
	}
	// Convert it to slice
	evalSet := []string{}
	for key := range evalsMap {
		evalSet = append(evalSet, key)
	}
	return evalSet
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

func (p *PersisterDwh) getUserEvaluation(
	environmentNamespace,
	userID,
	tag,
	featureID string,
	featureVersion int32,
) (*featureproto.Evaluation, error) {
	
}

// Because the same goal can be used on other feature flags
// we must link the goal event to all flags using the same goal ID
// If it fails even once to link the goal event, it will retry until all events are linked.
func (p *PersisterDwh) linkGoalEventByAutoOps(
	ctx context.Context,
	event *eventproto.GoalEvent,
	environmentNamespace string,
) ([]*featureproto.Evaluation, bool, error) {
	// List all auto ops rules
	list, err := p.listAutoOpsRules(ctx, environmentNamespace)
	if err != nil {
		return nil, true, err
	}
	if len(list) == 0 {
		return nil, false, ErrNoAutoOpsRules
	}
	// Find the feature flags by goal ID
	featureIDs := p.findFeatureIDs(event.GoalId, list)
	if len(featureIDs) == 0 {
		return nil, false, ErrAutoOpsRulesNotFound
	}
	// Get the lastest feature version
	resp, err := p.featureClient.GetFeatures(ctx, &featureproto.GetFeaturesRequest{
		EnvironmentNamespace: environmentNamespace,
		Ids:                  featureIDs,
	})
	if err != nil {
		handledCounter.WithLabelValues(codeFailedToGetFeatures).Inc()
		return nil, true, err
	}
	// Get all user evaluations using the feature flag info
	evaluations := make([]*featureproto.Evaluation, 0, len(resp.Features))
	for _, feature := range resp.Features {
		ev, err := p.getUserEvaluation(
			environmentNamespace,
			event.UserId,
			event.Tag,
			feature.Id,
			feature.Version,
		)
		if err != nil {
			return nil, true, err
		}
		evaluations = append(evaluations, ev)
	}
	return evaluations, false, nil
}

func (p *PersisterDwh) findFeatureIDs(goalID string, listAutoOpsRules []*aoproto.AutoOpsRule) []string {
	featureIDs := []string{}
	for _, aor := range listAutoOpsRules {
		autoOpsRule := &aodomain.AutoOpsRule{AutoOpsRule: aor}
		// We ignore the rules that are already triggered
		if autoOpsRule.AlreadyTriggered() {
			continue
		}
		clauses, err := autoOpsRule.ExtractOpsEventRateClauses()
		if err != nil {
			handledCounter.WithLabelValues(codeFailedToExtractOpsEventRateClauses).Inc()
			continue
		}
		for _, clause := range clauses {
			if clause.GoalId == goalID {
				featureIDs = append(featureIDs, autoOpsRule.FeatureId)
			}
		}
	}
	return featureIDs
}

func (p *PersisterDwh) listAutoOpsRules(
	ctx context.Context,
	environmentNamespace string,
) ([]*aoproto.AutoOpsRule, error) {
	exp, err, _ := p.flightgroup.Do(
		fmt.Sprintf("%s:%s", environmentNamespace, "listAutoOpsRules"),
		func() (interface{}, error) {
			aor := []*aoproto.AutoOpsRule{}
			cursor := ""
			for {
				resp, err := p.autoOpsClient.ListAutoOpsRules(ctx, &aoproto.ListAutoOpsRulesRequest{
					EnvironmentNamespace: environmentNamespace,
					PageSize:             listRequestSize,
					Cursor:               cursor,
				})
				if err != nil {
					return nil, err
				}
				aor = append(aor, resp.AutoOpsRules...)
				aorSize := len(resp.AutoOpsRules)
				if aorSize == 0 || aorSize < listRequestSize {
					return aor, nil
				}
				cursor = resp.Cursor
			}
		},
	)
	if err != nil {
		handledCounter.WithLabelValues(codeFailedToListAutoOpsRules).Inc()
		return nil, err
	}
	return exp.([]*aoproto.AutoOpsRule), nil
}
