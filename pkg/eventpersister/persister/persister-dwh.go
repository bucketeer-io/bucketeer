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
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/eventpersister/datastore"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
)

type PersisterDwh struct {
	puller          puller.RateLimitedPuller
	logger          *zap.Logger
	ctx             context.Context
	cancel          func()
	group           errgroup.Group
	doneCh          chan struct{}
	evalEventWriter datastore.EvalEventWriter
	goalEventWriter datastore.GoalEventWriter
	opts            *options
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
				e, err := p.convToGoalEvent(ctx, evt, id, environmentNamespace)
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
	ud, err := json.Marshal(e.User.Data)
	if err != nil {
		return nil, err
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

func (*PersisterDwh) convToGoalEvent(
	ctx context.Context,
	e *eventproto.GoalEvent,
	id, environmentNamespace string,
) (*ecproto.GoalEvent, error) {
	return &ecproto.GoalEvent{}, nil
}
