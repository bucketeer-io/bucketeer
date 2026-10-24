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
	experimentdomain "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	aoproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

const (
	linkedEventsMemoryCacheTTL = 1 * time.Minute
	// Bounds reloads when evictions keep arriving during a load
	linkedTargetsMaxLoadAttempts = 3
	// Detached from the request since Track publishes in the background
	linkedEventsTrackPublishTimeout = 10 * time.Second

	topicEvaluationDWH = "EvaluationDWH"
	topicEvaluationOPS = "EvaluationOPS"
	topicGoalDWH       = "GoalDWH"
	topicGoalOPS       = "GoalOPS"

	typeLinkedTargets = "LinkedTargets"
	typeExperiments   = "Experiments"
	typeAutoOpsRules  = "AutoOpsRules"

	codeLoadFailed = "LoadFailed"
)

var (
	errLinkedTargetsLoadFailed = errors.New("gateway: failed to load linked targets")
	errInvalidLinkedEvent      = errors.New("gateway: invalid linked event")
)

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
	// Bumped on eviction so an in-flight load does not store a stale snapshot
	generations map[string]uint64
	flightgroup singleflight.Group
	now         func() time.Time
	logger      *zap.Logger
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
		autoOpsRulesRedisCache: cachev3.NewAutoOpsRulesCache(redisCache, cachev3.NonPersistentAutoOpsRulesCacheTTL),
		experimentClient:       experimentClient,
		autoOpsClient:          autoOpsClient,
		memoryCacheTTL:         linkedEventsMemoryCacheTTL,
		targets:                make(map[string]*linkedTargets),
		generations:            make(map[string]uint64),
		now:                    time.Now,
		logger:                 logger.Named("linked-event-publisher"),
	}
}

// Publish publishes the events linked to experiments or auto ops rules to the dedicated topics
// and returns the errors by event ID.
func (p *LinkedEventPublisher) Publish(
	ctx context.Context,
	environmentID string,
	goals []publisher.Message,
	evaluations []publisher.Message,
) map[string]error {
	errs := make(map[string]error)
	if p.topics.GoalDWH == nil && p.topics.GoalOPS == nil {
		goals = nil
	}
	if p.topics.EvaluationDWH == nil && p.topics.EvaluationOPS == nil {
		evaluations = nil
	}
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
	evalDWH, evalOPS := p.linkEvaluations(targets, environmentID, evaluations)
	goalDWH, goalOPS := p.linkGoals(targets, environmentID, goals)
	p.publish(ctx, p.topics.EvaluationDWH, evalDWH, topicEvaluationDWH, environmentID, errs)
	p.publish(ctx, p.topics.EvaluationOPS, evalOPS, topicEvaluationOPS, environmentID, errs)
	p.publish(ctx, p.topics.GoalDWH, goalDWH, topicGoalDWH, environmentID, errs)
	p.publish(ctx, p.topics.GoalOPS, goalOPS, topicGoalOPS, environmentID, errs)
	return errs
}

// linkEvaluations returns the evaluations linked to experiments (DWH) and auto ops rules (OPS).
func (p *LinkedEventPublisher) linkEvaluations(
	targets *linkedTargets,
	environmentID string,
	evaluations []publisher.Message,
) (dwh, ops []publisher.Message) {
	for _, m := range evaluations {
		e := &eventproto.EvaluationEvent{}
		if !p.unmarshalEvent(m, e, typeEvaluation, environmentID) {
			continue
		}
		if _, ok := targets.experimentFeatures[fmt.Sprintf("%s:%d", e.FeatureId, e.FeatureVersion)]; ok {
			dwh = append(dwh, m)
		}
		if _, ok := targets.opsFeatureVariations[e.FeatureId+":"+e.VariationId]; ok {
			ops = append(ops, m)
		}
	}
	return dwh, ops
}

// linkGoals returns the goals linked to experiments (DWH) and auto ops rules (OPS).
func (p *LinkedEventPublisher) linkGoals(
	targets *linkedTargets,
	environmentID string,
	goals []publisher.Message,
) (dwh, ops []publisher.Message) {
	for _, m := range goals {
		e := &eventproto.GoalEvent{}
		if !p.unmarshalEvent(m, e, typeGoal, environmentID) {
			continue
		}
		if _, ok := targets.experimentGoals[e.GoalId]; ok {
			dwh = append(dwh, m)
		}
		if _, ok := targets.opsGoals[e.GoalId]; ok {
			ops = append(ops, m)
		}
	}
	return dwh, ops
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

func (p *LinkedEventPublisher) unmarshalEvent(
	m publisher.Message,
	dst proto.Message,
	eventType, environmentID string,
) bool {
	err := errInvalidLinkedEvent
	if ev, ok := m.(*eventproto.Event); ok && ev.Event != nil {
		err = ev.Event.UnmarshalTo(dst)
	}
	if err != nil {
		p.logger.Error("Failed to unmarshal linked event",
			zap.Error(err),
			zap.String("environmentId", environmentID),
			zap.String("type", eventType),
			zap.String("id", m.GetId()),
		)
		linkedEventCounter.WithLabelValues(eventType, codeUnmarshalFailed).Inc()
		return false
	}
	return true
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

// Evict drops the environment's targets when its experiments or auto ops rules change.
func (p *LinkedEventPublisher) Evict(environmentID string) {
	p.mu.Lock()
	delete(p.targets, environmentID)
	p.generations[environmentID]++
	p.mu.Unlock()
	p.flightgroup.Forget(environmentID)
}

// loadTargets builds the environment's linked targets and stores them in memory.
// When the environment is evicted while loading, it reloads so waiting callers do not get the stale targets.
func (p *LinkedEventPublisher) loadTargets(ctx context.Context, environmentID string) (*linkedTargets, error) {
	for attempt := 1; ; attempt++ {
		p.mu.RLock()
		generation := p.generations[environmentID]
		p.mu.RUnlock()
		t, err := p.buildTargets(ctx, environmentID)
		if err != nil {
			return nil, err
		}
		p.mu.Lock()
		evicted := p.generations[environmentID] != generation
		if evicted && attempt < linkedTargetsMaxLoadAttempts {
			p.mu.Unlock()
			continue
		}
		if !evicted {
			p.targets[environmentID] = t
		}
		// Drop expired entries of idle environments
		now := p.now()
		for id, entry := range p.targets {
			if !now.Before(entry.expiresAt) {
				delete(p.targets, id)
			}
		}
		p.mu.Unlock()
		return t, nil
	}
}

func (p *LinkedEventPublisher) buildTargets(ctx context.Context, environmentID string) (*linkedTargets, error) {
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
	// Captured before reading so a newer snapshot from the cache refresher wins
	generation := p.now().UnixNano()
	resp, err := p.experimentClient.ListExperiments(ctx, &experimentproto.ListExperimentsRequest{
		StopAt:        p.now().Add(-experimentdomain.StoppedExperimentEventWindow).Unix(),
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
	exps := &experimentproto.Experiments{Experiments: resp.Experiments}
	if _, err := p.experimentsRedisCache.PutIfNewer(exps, environmentID, generation); err != nil {
		p.logger.Error("Failed to cache experiments",
			zap.Error(err),
			zap.String("environmentId", environmentID),
		)
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
	generation := p.now().UnixNano()
	resp, err := p.autoOpsClient.ListAutoOpsRules(ctx, &aoproto.ListAutoOpsRulesRequest{
		PageSize:      0,
		EnvironmentId: environmentID,
	})
	if err != nil {
		return nil, err
	}
	rules := &aoproto.AutoOpsRules{AutoOpsRules: resp.AutoOpsRules}
	if _, err := p.autoOpsRulesRedisCache.PutIfNewer(rules, environmentID, generation); err != nil {
		p.logger.Error("Failed to cache auto ops rules",
			zap.Error(err),
			zap.String("environmentId", environmentID),
		)
	}
	return resp.AutoOpsRules, nil
}

// excludeFailed drops the events that failed on the dedicated topics, so a retry does not publish them twice
func excludeFailed[T any](messages []publisher.Message, failed map[string]T) []publisher.Message {
	result := make([]publisher.Message, 0, len(messages))
	for _, m := range messages {
		if _, ok := failed[m.GetId()]; !ok {
			result = append(result, m)
		}
	}
	return result
}
