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

package api

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	autoopsclient "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client"
	aodomain "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	aoproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

const (
	linkedEventsMemoryCacheTTL = 1 * time.Minute
	// Same window as the experiment cacher, since events may arrive late
	linkedExperimentsStopAtWindow = 2 * 24 * time.Hour

	topicEvaluationDWH = "EvaluationDWH"
	topicEvaluationOPS = "EvaluationOPS"
	topicGoalDWH       = "GoalDWH"
	topicGoalOPS       = "GoalOPS"

	typeLinkedTargets = "LinkedTargets"
	typeExperiments   = "Experiments"
	typeAutoOpsRules  = "AutoOpsRules"

	codeLoadFailed = "LoadFailed"
)

var errLinkedTargetsLoadFailed = errors.New("gateway: failed to load linked targets")

// LinkedEventTopics holds the dedicated publishers. A nil publisher disables that topic.
type LinkedEventTopics struct {
	EvaluationDWH publisher.Publisher
	EvaluationOPS publisher.Publisher
	GoalDWH       publisher.Publisher
	GoalOPS       publisher.Publisher
}

func (t LinkedEventTopics) enabled() bool {
	return t.EvaluationDWH != nil || t.EvaluationOPS != nil || t.GoalDWH != nil || t.GoalOPS != nil
}

// LinkedEventPublisher publishes events linked to experiments or auto ops rules to dedicated topics.
type LinkedEventPublisher struct {
	topics                 LinkedEventTopics
	experimentsRedisCache  cachev3.ExperimentsCache
	autoOpsRulesRedisCache cachev3.AutoOpsRulesCache
	experimentClient       experimentclient.Client
	autoOpsClient          autoopsclient.Client
	memoryCacheTTL         time.Duration
	mu                     sync.RWMutex
	targets                map[string]*linkedTargets
	flightgroup            singleflight.Group
	now                    func() time.Time
	logger                 *zap.Logger
}

type linkedTargets struct {
	// featureID:featureVersion of experiments
	experimentFeatures map[string]struct{}
	experimentGoals    map[string]struct{}
	// featureID:variationID of active event rate clauses
	opsFeatureVariations map[string]struct{}
	opsGoals             map[string]struct{}
	expiresAt            time.Time
}

// NewLinkedEventPublisher returns nil when no dedicated topic is configured.
func NewLinkedEventPublisher(
	topics LinkedEventTopics,
	redisCache cache.MultiGetCache,
	experimentClient experimentclient.Client,
	autoOpsClient autoopsclient.Client,
	logger *zap.Logger,
) *LinkedEventPublisher {
	if !topics.enabled() {
		return nil
	}
	return &LinkedEventPublisher{
		topics:                 topics,
		experimentsRedisCache:  cachev3.NewExperimentsCache(redisCache),
		autoOpsRulesRedisCache: cachev3.NewAutoOpsRulesCache(redisCache),
		experimentClient:       experimentClient,
		autoOpsClient:          autoOpsClient,
		memoryCacheTTL:         linkedEventsMemoryCacheTTL,
		targets:                make(map[string]*linkedTargets),
		now:                    time.Now,
		logger:                 logger.Named("linked-event-publisher"),
	}
}

// Publish publishes the linked events and returns the errors by event ID.
func (p *LinkedEventPublisher) Publish(
	ctx context.Context,
	environmentID string,
	goals []publisher.Message,
	evaluations []publisher.Message,
) map[string]error {
	errs := make(map[string]error)
	if len(goals) == 0 && len(evaluations) == 0 {
		return errs
	}
	targets, err := p.getTargets(ctx, environmentID)
	if err != nil {
		p.logger.Error("Failed to load linked targets",
			zap.Error(err),
			zap.String("environmentId", environmentID),
		)
		for _, m := range goals {
			errs[m.GetId()] = errLinkedTargetsLoadFailed
		}
		for _, m := range evaluations {
			errs[m.GetId()] = errLinkedTargetsLoadFailed
		}
		linkedEventCounter.WithLabelValues(codeAll, codeLoadFailed).Add(float64(len(errs)))
		return errs
	}
	var evalDWH, evalOPS, goalDWH, goalOPS []publisher.Message
	for _, m := range evaluations {
		e := &eventproto.EvaluationEvent{}
		if !unmarshalEvent(m, e) {
			continue
		}
		if targets.linkEvaluationDWH(e) {
			evalDWH = append(evalDWH, m)
		}
		if targets.linkEvaluationOPS(e) {
			evalOPS = append(evalOPS, m)
		}
	}
	for _, m := range goals {
		e := &eventproto.GoalEvent{}
		if !unmarshalEvent(m, e) {
			continue
		}
		if _, ok := targets.experimentGoals[e.GoalId]; ok {
			goalDWH = append(goalDWH, m)
		}
		if _, ok := targets.opsGoals[e.GoalId]; ok {
			goalOPS = append(goalOPS, m)
		}
	}
	p.publish(ctx, p.topics.EvaluationDWH, evalDWH, topicEvaluationDWH, environmentID, errs)
	p.publish(ctx, p.topics.EvaluationOPS, evalOPS, topicEvaluationOPS, environmentID, errs)
	p.publish(ctx, p.topics.GoalDWH, goalDWH, topicGoalDWH, environmentID, errs)
	p.publish(ctx, p.topics.GoalOPS, goalOPS, topicGoalOPS, environmentID, errs)
	return errs
}

func (p *LinkedEventPublisher) publish(
	ctx context.Context,
	pub publisher.Publisher,
	messages []publisher.Message,
	topic, environmentID string,
	errs map[string]error,
) {
	if pub == nil || len(messages) == 0 {
		return
	}
	multiErrs := pub.PublishMulti(ctx, messages)
	var repeatableErrors, nonRepeatableErrors float64
	for id, err := range multiErrs {
		if errors.Is(err, publisher.ErrBadMessage) {
			nonRepeatableErrors++
		} else {
			repeatableErrors++
		}
		if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			p.logger.Error("Failed to publish linked event",
				zap.Error(err),
				zap.String("environmentId", environmentID),
				zap.String("topic", topic),
				zap.String("id", id),
			)
		}
		// Prefer a retriable error when the event failed on multiple topics
		if _, ok := errs[id]; !ok || !errors.Is(err, publisher.ErrBadMessage) {
			errs[id] = err
		}
	}
	linkedEventCounter.WithLabelValues(topic, codeNonRepeatableError).Add(nonRepeatableErrors)
	linkedEventCounter.WithLabelValues(topic, codeRepeatableError).Add(repeatableErrors)
	linkedEventCounter.WithLabelValues(topic, codeOK).Add(float64(len(messages) - len(multiErrs)))
}

func unmarshalEvent(m publisher.Message, dst proto.Message) bool {
	ev, ok := m.(*eventproto.Event)
	if !ok || ev.Event == nil {
		return false
	}
	return ev.Event.UnmarshalTo(dst) == nil
}

func (t *linkedTargets) linkEvaluationDWH(e *eventproto.EvaluationEvent) bool {
	_, ok := t.experimentFeatures[fmt.Sprintf("%s:%d", e.FeatureId, e.FeatureVersion)]
	return ok
}

func (t *linkedTargets) linkEvaluationOPS(e *eventproto.EvaluationEvent) bool {
	_, ok := t.opsFeatureVariations[e.FeatureId+":"+e.VariationId]
	return ok
}

func (p *LinkedEventPublisher) getTargets(ctx context.Context, environmentID string) (*linkedTargets, error) {
	p.mu.RLock()
	t, ok := p.targets[environmentID]
	p.mu.RUnlock()
	if ok && p.now().Before(t.expiresAt) {
		cacheCounter.WithLabelValues(callerGatewayService, typeLinkedTargets, cacheLayerInMemory, codeHit).Inc()
		return t, nil
	}
	cacheCounter.WithLabelValues(callerGatewayService, typeLinkedTargets, cacheLayerInMemory, codeMiss).Inc()
	ch := p.flightgroup.DoChan(environmentID, func() (interface{}, error) {
		innerCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), singleflightFetchTimeout)
		defer cancel()
		return p.loadTargets(innerCtx, environmentID)
	})
	select {
	case res := <-ch:
		if res.Err != nil {
			return nil, res.Err
		}
		return res.Val.(*linkedTargets), nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (p *LinkedEventPublisher) loadTargets(ctx context.Context, environmentID string) (*linkedTargets, error) {
	t := &linkedTargets{
		experimentFeatures:   make(map[string]struct{}),
		experimentGoals:      make(map[string]struct{}),
		opsFeatureVariations: make(map[string]struct{}),
		opsGoals:             make(map[string]struct{}),
	}
	if p.topics.EvaluationDWH != nil || p.topics.GoalDWH != nil {
		experiments, err := p.listExperiments(ctx, environmentID)
		if err != nil {
			return nil, err
		}
		for _, exp := range experiments {
			t.experimentFeatures[fmt.Sprintf("%s:%d", exp.FeatureId, exp.FeatureVersion)] = struct{}{}
			for _, goalID := range exp.GoalIds {
				t.experimentGoals[goalID] = struct{}{}
			}
		}
	}
	if p.topics.EvaluationOPS != nil || p.topics.GoalOPS != nil {
		rules, err := p.listAutoOpsRules(ctx, environmentID)
		if err != nil {
			return nil, err
		}
		for _, rule := range rules {
			r := &aodomain.AutoOpsRule{AutoOpsRule: rule}
			if r.IsFinished() || r.IsStopped() {
				continue
			}
			clauses, err := r.ExtractOpsEventRateClauses()
			if err != nil {
				continue
			}
			for _, clause := range clauses {
				t.opsFeatureVariations[rule.FeatureId+":"+clause.VariationId] = struct{}{}
				t.opsGoals[clause.GoalId] = struct{}{}
			}
		}
	}
	t.expiresAt = p.now().Add(p.memoryCacheTTL)
	p.mu.Lock()
	p.targets[environmentID] = t
	p.mu.Unlock()
	return t, nil
}

func (p *LinkedEventPublisher) listExperiments(
	ctx context.Context,
	environmentID string,
) ([]*experimentproto.Experiment, error) {
	if exps, err := p.experimentsRedisCache.Get(environmentID); err == nil {
		cacheCounter.WithLabelValues(callerGatewayService, typeExperiments, cacheLayerExternal, codeHit).Inc()
		return exps.Experiments, nil
	}
	cacheCounter.WithLabelValues(callerGatewayService, typeExperiments, cacheLayerExternal, codeMiss).Inc()
	resp, err := p.experimentClient.ListExperiments(ctx, &experimentproto.ListExperimentsRequest{
		StopAt:        p.now().Add(-linkedExperimentsStopAtWindow).Unix(),
		PageSize:      0,
		EnvironmentId: environmentID,
		Statuses: []experimentproto.Experiment_Status{
			experimentproto.Experiment_RUNNING,
			experimentproto.Experiment_STOPPED,
		},
		Archived: &wrapperspb.BoolValue{Value: false},
	})
	if err != nil {
		return nil, err
	}
	return resp.Experiments, nil
}

func (p *LinkedEventPublisher) listAutoOpsRules(
	ctx context.Context,
	environmentID string,
) ([]*aoproto.AutoOpsRule, error) {
	if rules, err := p.autoOpsRulesRedisCache.Get(environmentID); err == nil {
		cacheCounter.WithLabelValues(callerGatewayService, typeAutoOpsRules, cacheLayerExternal, codeHit).Inc()
		return rules.AutoOpsRules, nil
	}
	cacheCounter.WithLabelValues(callerGatewayService, typeAutoOpsRules, cacheLayerExternal, codeMiss).Inc()
	resp, err := p.autoOpsClient.ListAutoOpsRules(ctx, &aoproto.ListAutoOpsRulesRequest{
		PageSize:      0,
		EnvironmentId: environmentID,
	})
	if err != nil {
		return nil, err
	}
	return resp.AutoOpsRules, nil
}

// excludeFailed drops the events that failed on the main topics since the SDK retries them
func excludeFailed[T any](messages []publisher.Message, failed map[string]T) []publisher.Message {
	result := make([]publisher.Message, 0, len(messages))
	for _, m := range messages {
		if _, ok := failed[m.GetId()]; !ok {
			result = append(result, m)
		}
	}
	return result
}
