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
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	autoopsclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachemock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/mock"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	aoproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

const linkedTestEnvID = "env-1"

type linkedTestMocks struct {
	cache         *cachemock.MockMultiGetCache
	expClient     *experimentclientmock.MockClient
	autoOpsClient *autoopsclientmock.MockClient
	evalDWH       *publishermock.MockPublisher
	evalOPS       *publishermock.MockPublisher
	goalDWH       *publishermock.MockPublisher
	goalOPS       *publishermock.MockPublisher
	published     map[string][]string
}

func TestLinkedEventPublisherPublish(t *testing.T) {
	t.Parallel()
	experiments := &experimentproto.Experiments{Experiments: []*experimentproto.Experiment{
		{Id: "exp-1", FeatureId: "feature-1", FeatureVersion: 1, GoalIds: []string{"goal-1"}},
	}}
	rules := &aoproto.AutoOpsRules{AutoOpsRules: []*aoproto.AutoOpsRule{
		newLinkedTestRule(t, "feature-2", "variation-a", "goal-2", aoproto.AutoOpsStatus_RUNNING),
		newLinkedTestRule(t, "feature-3", "variation-a", "goal-3", aoproto.AutoOpsStatus_FINISHED),
	}}
	evaluations := []publisher.Message{
		newLinkedTestEvaluation(t, "e-exp", "feature-1", 1, "variation-a"),
		newLinkedTestEvaluation(t, "e-old-version", "feature-1", 2, "variation-a"),
		newLinkedTestEvaluation(t, "e-ops", "feature-2", 1, "variation-a"),
		newLinkedTestEvaluation(t, "e-other-variation", "feature-2", 1, "variation-b"),
		newLinkedTestEvaluation(t, "e-finished-rule", "feature-3", 1, "variation-a"),
	}
	goals := []publisher.Message{
		newLinkedTestGoal(t, "g-exp", "goal-1"),
		newLinkedTestGoal(t, "g-ops", "goal-2"),
		newLinkedTestGoal(t, "g-finished-rule", "goal-3"),
		newLinkedTestGoal(t, "g-none", "goal-4"),
	}
	linked := map[string][]string{
		topicEvaluationDWH: {"e-exp"},
		topicEvaluationOPS: {"e-ops"},
		topicGoalDWH:       {"g-exp"},
		topicGoalOPS:       {"g-ops"},
	}
	errDB := errors.New("db error")
	errPublish := errors.New("publish error")

	patterns := []struct {
		desc              string
		setup             func(*linkedTestMocks)
		disableOPS        bool
		expectedPublished map[string][]string
		expectedErrs      map[string]error
	}{
		{
			desc: "success: from redis",
			setup: func(m *linkedTestMocks) {
				m.expectRedis(t, experiments, rules)
				m.expectPublish(nil)
			},
			expectedPublished: linked,
			expectedErrs:      map[string]error{},
		},
		{
			desc: "success: database fallback on redis miss",
			setup: func(m *linkedTestMocks) {
				m.cache.EXPECT().Get(linkedTestEnvID+":experiment").Return(nil, cache.ErrNotFound)
				m.cache.EXPECT().Get(linkedTestEnvID+":autoOpsRule").Return(nil, cache.ErrNotFound)
				m.expClient.EXPECT().ListExperiments(gomock.Any(), gomock.Any()).Return(
					&experimentproto.ListExperimentsResponse{Experiments: experiments.Experiments}, nil)
				m.autoOpsClient.EXPECT().ListAutoOpsRules(gomock.Any(), gomock.Any()).Return(
					&aoproto.ListAutoOpsRulesResponse{AutoOpsRules: rules.AutoOpsRules}, nil)
				// The database result is written back to redis
				for key, ttl := range map[string]time.Duration{
					linkedTestEnvID + ":experiment":  0,
					linkedTestEnvID + ":autoOpsRule": cachev3.NonPersistentAutoOpsRulesCacheTTL,
				} {
					m.cache.EXPECT().PutIfNewer(key, "{"+key+"}:gen", gomock.Any(), gomock.Any(), ttl).Return(true, nil)
				}
				m.expectPublish(nil)
			},
			expectedPublished: linked,
			expectedErrs:      map[string]error{},
		},
		{
			desc: "success: disabled topics are skipped",
			setup: func(m *linkedTestMocks) {
				m.cache.EXPECT().Get(linkedTestEnvID+":experiment").Return(marshalLinkedTest(t, experiments), nil)
				m.expectPublish(nil)
			},
			disableOPS: true,
			expectedPublished: map[string][]string{
				topicEvaluationDWH: {"e-exp"},
				topicGoalDWH:       {"g-exp"},
			},
			expectedErrs: map[string]error{},
		},
		{
			desc: "err: all events fail when targets cannot be loaded",
			setup: func(m *linkedTestMocks) {
				m.cache.EXPECT().Get(linkedTestEnvID+":experiment").Return(nil, cache.ErrNotFound)
				m.expClient.EXPECT().ListExperiments(gomock.Any(), gomock.Any()).Return(nil, errDB)
			},
			expectedPublished: map[string][]string{},
			expectedErrs:      allLinkedTestErrs(append(evaluations, goals...), errLinkedTargetsLoadFailed),
		},
		{
			desc: "err: publish failure",
			setup: func(m *linkedTestMocks) {
				m.expectRedis(t, experiments, rules)
				m.expectPublish(map[string]error{"e-ops": errPublish})
			},
			expectedPublished: linked,
			expectedErrs:      map[string]error{"e-ops": errPublish},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			m := newLinkedTestMocks(ctrl)
			p.setup(m)
			topics := LinkedEventTopics{EvaluationDWH: m.evalDWH, EvaluationOPS: m.evalOPS, GoalDWH: m.goalDWH, GoalOPS: m.goalOPS}
			if p.disableOPS {
				topics.EvaluationOPS, topics.GoalOPS = nil, nil
			}
			lp := NewLinkedEventPublisher(topics, m.cache, m.expClient, m.autoOpsClient, zap.NewNop())
			errs := lp.Publish(context.Background(), linkedTestEnvID, goals, evaluations)
			assert.Equal(t, p.expectedErrs, errs)
			assert.Equal(t, p.expectedPublished, m.published)
		})
	}
}

func TestLinkedEventPublisherMemoryCache(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	m := newLinkedTestMocks(ctrl)
	// Redis is read only once while the in-memory entry is valid
	m.expectRedis(t, &experimentproto.Experiments{}, &aoproto.AutoOpsRules{})
	lp := NewLinkedEventPublisher(
		LinkedEventTopics{EvaluationDWH: m.evalDWH, EvaluationOPS: m.evalOPS, GoalDWH: m.goalDWH, GoalOPS: m.goalOPS},
		m.cache, m.expClient, m.autoOpsClient, zap.NewNop(),
	)
	goals := []publisher.Message{newLinkedTestGoal(t, "g-1", "goal-1")}
	for range 3 {
		assert.Empty(t, lp.Publish(context.Background(), linkedTestEnvID, goals, nil))
	}
}

func TestNewLinkedEventPublisherDisabled(t *testing.T) {
	t.Parallel()
	assert.Nil(t, NewLinkedEventPublisher(LinkedEventTopics{}, nil, nil, nil, zap.NewNop()))
}

func TestExcludeFailed(t *testing.T) {
	t.Parallel()
	messages := []publisher.Message{&eventproto.Event{Id: "a"}, &eventproto.Event{Id: "b"}}
	actual := excludeFailed(messages, map[string]error{"a": errors.New("error")})
	require.Len(t, actual, 1)
	assert.Equal(t, "b", actual[0].GetId())
}

func newLinkedTestMocks(ctrl *gomock.Controller) *linkedTestMocks {
	return &linkedTestMocks{
		cache:         cachemock.NewMockMultiGetCache(ctrl),
		expClient:     experimentclientmock.NewMockClient(ctrl),
		autoOpsClient: autoopsclientmock.NewMockClient(ctrl),
		evalDWH:       publishermock.NewMockPublisher(ctrl),
		evalOPS:       publishermock.NewMockPublisher(ctrl),
		goalDWH:       publishermock.NewMockPublisher(ctrl),
		goalOPS:       publishermock.NewMockPublisher(ctrl),
		published:     map[string][]string{},
	}
}

func (m *linkedTestMocks) expectRedis(t *testing.T, exps *experimentproto.Experiments, rules *aoproto.AutoOpsRules) {
	m.cache.EXPECT().Get(linkedTestEnvID+":experiment").Return(marshalLinkedTest(t, exps), nil)
	m.cache.EXPECT().Get(linkedTestEnvID+":autoOpsRule").Return(marshalLinkedTest(t, rules), nil)
}

func (m *linkedTestMocks) expectPublish(errs map[string]error) {
	for topic, pub := range map[string]*publishermock.MockPublisher{
		topicEvaluationDWH: m.evalDWH,
		topicEvaluationOPS: m.evalOPS,
		topicGoalDWH:       m.goalDWH,
		topicGoalOPS:       m.goalOPS,
	} {
		pub.EXPECT().PublishMulti(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, messages []publisher.Message) map[string]error {
				failed := map[string]error{}
				for _, msg := range messages {
					m.published[topic] = append(m.published[topic], msg.GetId())
					if err, ok := errs[msg.GetId()]; ok {
						failed[msg.GetId()] = err
					}
				}
				sort.Strings(m.published[topic])
				return failed
			},
		).AnyTimes()
	}
}

func allLinkedTestErrs(messages []publisher.Message, err error) map[string]error {
	errs := make(map[string]error, len(messages))
	for _, m := range messages {
		errs[m.GetId()] = err
	}
	return errs
}

func marshalLinkedTest(t *testing.T, m proto.Message) []byte {
	b, err := proto.Marshal(m)
	require.NoError(t, err)
	return b
}

func newLinkedTestRule(t *testing.T, featureID, variationID, goalID string, status aoproto.AutoOpsStatus) *aoproto.AutoOpsRule {
	c, err := anypb.New(&aoproto.OpsEventRateClause{VariationId: variationID, GoalId: goalID})
	require.NoError(t, err)
	return &aoproto.AutoOpsRule{
		Id:            featureID + "-rule",
		FeatureId:     featureID,
		AutoOpsStatus: status,
		Clauses:       []*aoproto.Clause{{Id: featureID + "-clause", Clause: c}},
	}
}

func newLinkedTestEvaluation(t *testing.T, id, featureID string, version int32, variationID string) *eventproto.Event {
	e, err := anypb.New(&eventproto.EvaluationEvent{FeatureId: featureID, FeatureVersion: version, VariationId: variationID})
	require.NoError(t, err)
	return &eventproto.Event{Id: id, Event: e, EnvironmentId: linkedTestEnvID}
}

func newLinkedTestGoal(t *testing.T, id, goalID string) *eventproto.Event {
	e, err := anypb.New(&eventproto.GoalEvent{GoalId: goalID})
	require.NoError(t, err)
	return &eventproto.Event{Id: id, Event: e, EnvironmentId: linkedTestEnvID}
}

func TestLinkedEventPublisherSkipsDisabledRoutes(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	m := newLinkedTestMocks(ctrl)
	// No cache or client calls are expected because goal routes are disabled
	lp := NewLinkedEventPublisher(
		LinkedEventTopics{EvaluationDWH: m.evalDWH},
		m.cache, m.expClient, m.autoOpsClient, zap.NewNop(),
	)
	goals := []publisher.Message{newLinkedTestGoal(t, "g-1", "goal-1")}
	assert.Empty(t, lp.Publish(context.Background(), linkedTestEnvID, goals, nil))
}

func TestLinkedEventPublisherEvict(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		setup        func(*linkedTestMocks, func() *LinkedEventPublisher)
		evictBetween bool
	}{
		{
			desc:         "reloads after eviction",
			evictBetween: true,
			setup: func(m *linkedTestMocks, _ func() *LinkedEventPublisher) {
				m.cache.EXPECT().Get(linkedTestEnvID+":experiment").Return(
					marshalLinkedTest(t, &experimentproto.Experiments{}), nil).Times(2)
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			m := newLinkedTestMocks(ctrl)
			var lp *LinkedEventPublisher
			p.setup(m, func() *LinkedEventPublisher { return lp })
			lp = NewLinkedEventPublisher(LinkedEventTopics{GoalDWH: m.goalDWH}, m.cache, m.expClient, m.autoOpsClient, zap.NewNop())
			goals := []publisher.Message{newLinkedTestGoal(t, "g-1", "goal-1")}
			assert.Empty(t, lp.Publish(context.Background(), linkedTestEnvID, goals, nil))
			if p.evictBetween {
				lp.Evict(linkedTestEnvID)
			}
			assert.Empty(t, lp.Publish(context.Background(), linkedTestEnvID, goals, nil))
		})
	}
}

func TestCacheInvalidatorEvictsLinkedTargets(t *testing.T) {
	t.Parallel()
	for _, entityType := range []domaineventproto.Event_EntityType{
		domaineventproto.Event_EXPERIMENT,
		domaineventproto.Event_AUTOOPS_RULE,
	} {
		t.Run(entityType.String(), func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			m := newLinkedTestMocks(ctrl)
			lp := NewLinkedEventPublisher(LinkedEventTopics{GoalDWH: m.goalDWH}, m.cache, m.expClient, m.autoOpsClient, zap.NewNop())
			for _, kind := range []string{linkedTargetsExperiments, linkedTargetsAutoOpsRules} {
				lp.targets[linkedTargetsKey(linkedTestEnvID, kind)] = &linkedTargets{}
			}
			ci := NewCacheInvalidator(nil, nil, nil, lp, nil, zap.NewNop())
			require.NoError(t, ci.evict(&domaineventproto.Event{EntityType: entityType, EnvironmentId: linkedTestEnvID}))
			assert.Empty(t, lp.targets)
		})
	}
}

func TestLinkedEventPublisherDropsExpiredTargets(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	m := newLinkedTestMocks(ctrl)
	m.cache.EXPECT().Get(linkedTestEnvID+":experiment").Return(marshalLinkedTest(t, &experimentproto.Experiments{}), nil)
	lp := NewLinkedEventPublisher(LinkedEventTopics{GoalDWH: m.goalDWH}, m.cache, m.expClient, m.autoOpsClient, zap.NewNop())
	now := time.Now()
	lp.now = func() time.Time { return now }
	idleKey := linkedTargetsKey("idle-env", linkedTargetsExperiments)
	activeKey := linkedTargetsKey("active-env", linkedTargetsExperiments)
	lp.targets[idleKey] = &linkedTargets{expiresAt: now.Add(-time.Second)}
	lp.targets[activeKey] = &linkedTargets{expiresAt: now.Add(time.Second)}
	goals := []publisher.Message{newLinkedTestGoal(t, "g-1", "goal-1")}
	assert.Empty(t, lp.Publish(context.Background(), linkedTestEnvID, goals, nil))
	assert.NotContains(t, lp.targets, idleKey)
	assert.Contains(t, lp.targets, activeKey)
	assert.Contains(t, lp.targets, linkedTargetsKey(linkedTestEnvID, linkedTargetsExperiments))
}

func TestGrpcTrackIgnoresLinkedEventFailure(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	m := newLinkedTestMocks(ctrl)
	gs := newGrpcGatewayServiceWithMock(t, ctrl)
	gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
		&accountproto.EnvironmentAPIKey{
			Environment: &environmentproto.EnvironmentV2{Id: linkedTestEnvID},
			ApiKey:      &accountproto.APIKey{Id: "id-0", Role: accountproto.APIKey_SDK_CLIENT},
		}, nil)
	gs.goalPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	m.cache.EXPECT().Get(linkedTestEnvID+":experiment").Return(marshalLinkedTest(t, &experimentproto.Experiments{
		Experiments: []*experimentproto.Experiment{{Id: "exp-1", GoalIds: []string{"goal-1"}}},
	}), nil)
	// A retry would get a new event ID and duplicate the goal on the main topic
	release, done := make(chan struct{}), make(chan struct{})
	m.goalDWH.EXPECT().PublishMulti(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, messages []publisher.Message) map[string]error {
			defer close(done)
			<-release
			return map[string]error{messages[0].GetId(): errors.New("publish error")}
		})
	opts := *gs.opts
	opts.linkedEventPublisher = NewLinkedEventPublisher(
		LinkedEventTopics{GoalDWH: m.goalDWH}, m.cache, m.expClient, m.autoOpsClient, zap.NewNop(),
	)
	gs.opts = &opts
	actual, err := gs.Track(context.Background(), &gwproto.TrackRequest{
		Apikey:    "api-key",
		Userid:    "user-id",
		Goalid:    "goal-1",
		Tag:       "tag",
		Timestamp: time.Now().Unix(),
	})
	// Returns while the linked publish is still blocked
	require.NoError(t, err)
	assert.Equal(t, &gwproto.TrackResponse{}, actual)
	close(release)
	<-done
}

func TestLinkedEventPublisherSkipsInvalidEvents(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	m := newLinkedTestMocks(ctrl)
	m.cache.EXPECT().Get(linkedTestEnvID+":experiment").Return(marshalLinkedTest(t, &experimentproto.Experiments{
		Experiments: []*experimentproto.Experiment{{Id: "exp-1", GoalIds: []string{"goal-1"}}},
	}), nil)
	m.expectPublish(nil)
	lp := NewLinkedEventPublisher(LinkedEventTopics{GoalDWH: m.goalDWH}, m.cache, m.expClient, m.autoOpsClient, zap.NewNop())
	// An evaluation event sent as a goal cannot be unmarshaled
	goals := []publisher.Message{
		newLinkedTestGoal(t, "g-1", "goal-1"),
		newLinkedTestEvaluation(t, "g-invalid", "feature-1", 1, "variation-a"),
	}
	assert.Empty(t, lp.Publish(context.Background(), linkedTestEnvID, goals, nil))
	assert.Equal(t, map[string][]string{topicGoalDWH: {"g-1"}}, m.published)
}

func TestGrpcRegisterEventsSkipsMainTopicsOnLinkedFailure(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	m := newLinkedTestMocks(ctrl)
	gs := newGrpcGatewayServiceWithMock(t, ctrl)
	linkedID, plainID := newUUID(t), newUUID(t)
	gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
		&accountproto.EnvironmentAPIKey{
			Environment: &environmentproto.EnvironmentV2{Id: linkedTestEnvID},
			ApiKey:      &accountproto.APIKey{Id: "id-0", Role: accountproto.APIKey_SDK_CLIENT},
		}, nil)
	m.cache.EXPECT().Get(linkedTestEnvID+":experiment").Return(marshalLinkedTest(t, &experimentproto.Experiments{
		Experiments: []*experimentproto.Experiment{{Id: "exp-1", GoalIds: []string{"goal-linked"}}},
	}), nil)
	m.expectPublish(map[string]error{linkedID: errors.New("publish error")})
	var mainPublished []string
	gs.goalPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, messages []publisher.Message) map[string]error {
			for _, msg := range messages {
				mainPublished = append(mainPublished, msg.GetId())
			}
			return nil
		})
	gs.evaluationPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(gomock.Any(), gomock.Any()).Return(nil)
	opts := *gs.opts
	opts.linkedEventPublisher = NewLinkedEventPublisher(
		LinkedEventTopics{GoalDWH: m.goalDWH}, m.cache, m.expClient, m.autoOpsClient, zap.NewNop(),
	)
	gs.opts = &opts
	newGoal := func(id, goalID string) *eventproto.Event {
		e, err := anypb.New(&eventproto.GoalEvent{
			Timestamp: time.Now().Unix(),
			GoalId:    goalID,
			UserId:    "user-id",
			User:      &userproto.User{Id: "user-id"},
		})
		require.NoError(t, err)
		return &eventproto.Event{Id: id, Event: e}
	}
	actual, err := gs.RegisterEvents(
		metadata.NewIncomingContext(context.Background(), metadata.MD{"authorization": []string{"api-key"}}),
		&gwproto.RegisterEventsRequest{Events: []*eventproto.Event{
			newGoal(linkedID, "goal-linked"),
			newGoal(plainID, "goal-plain"),
		}},
	)
	require.NoError(t, err)
	// The failed event is retried, so it must not be counted on the main topic yet
	assert.Equal(t, []string{plainID}, mainPublished)
	assert.Equal(t, map[string]*gwproto.RegisterEventsResponse_Error{
		linkedID: {Retriable: true, Message: "Failed to publish event"},
	}, actual.Errors)
}

func TestLinkedEventPublisherReloadsWhenEvictedDuringLoad(t *testing.T) {
	t.Parallel()
	stale := &experimentproto.Experiments{}
	fresh := &experimentproto.Experiments{
		Experiments: []*experimentproto.Experiment{{Id: "exp-1", GoalIds: []string{"goal-1"}}},
	}
	patterns := []struct {
		desc string
		// Number of loads that see an eviction while loading
		evictions         int
		expectedLoads     int
		expectedPublished map[string][]string
		expectedErrs      map[string]error
		expectedStored    bool
	}{
		{
			desc:              "reloads and returns the fresh targets",
			evictions:         1,
			expectedLoads:     2,
			expectedPublished: map[string][]string{topicGoalDWH: {"g-1"}},
			expectedErrs:      map[string]error{},
			expectedStored:    true,
		},
		{
			desc:              "fails as retriable after the max attempts without storing",
			evictions:         linkedTargetsMaxLoadAttempts,
			expectedLoads:     linkedTargetsMaxLoadAttempts,
			expectedPublished: map[string][]string{},
			expectedErrs:      map[string]error{"g-1": errLinkedTargetsLoadFailed},
			expectedStored:    false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			m := newLinkedTestMocks(ctrl)
			var lp *LinkedEventPublisher
			loads := 0
			m.cache.EXPECT().Get(linkedTestEnvID + ":experiment").DoAndReturn(func(string) (interface{}, error) {
				loads++
				if loads <= p.evictions {
					lp.Evict(linkedTestEnvID)
					// The last attempt still reads the fresh experiments
					if loads < linkedTargetsMaxLoadAttempts {
						return marshalLinkedTest(t, stale), nil
					}
				}
				return marshalLinkedTest(t, fresh), nil
			}).Times(p.expectedLoads)
			m.expectPublish(nil)
			lp = NewLinkedEventPublisher(LinkedEventTopics{GoalDWH: m.goalDWH}, m.cache, m.expClient, m.autoOpsClient, zap.NewNop())
			goals := []publisher.Message{newLinkedTestGoal(t, "g-1", "goal-1")}
			assert.Equal(t, p.expectedErrs, lp.Publish(context.Background(), linkedTestEnvID, goals, nil))
			assert.Equal(t, p.expectedPublished, m.published)
			_, stored := lp.targets[linkedTargetsKey(linkedTestEnvID, linkedTargetsExperiments)]
			assert.Equal(t, p.expectedStored, stored)
		})
	}
}

func TestLinkedEventPublisherLoadsOnlyNeededTargets(t *testing.T) {
	t.Parallel()
	experiments := &experimentproto.Experiments{Experiments: []*experimentproto.Experiment{
		{Id: "exp-1", FeatureId: "feature-1", FeatureVersion: 1},
	}}
	evaluations := []publisher.Message{newLinkedTestEvaluation(t, "e-1", "feature-1", 1, "variation-a")}
	goals := []publisher.Message{newLinkedTestGoal(t, "g-1", "goal-1")}
	patterns := []struct {
		desc              string
		setup             func(*linkedTestMocks)
		goals             []publisher.Message
		expectedPublished map[string][]string
		expectedErrs      map[string]error
	}{
		{
			desc: "success: an evaluation-only batch does not load auto ops rules",
			setup: func(m *linkedTestMocks) {
				m.cache.EXPECT().Get(linkedTestEnvID+":experiment").Return(marshalLinkedTest(t, experiments), nil)
			},
			expectedPublished: map[string][]string{topicEvaluationDWH: {"e-1"}},
			expectedErrs:      map[string]error{},
		},
		{
			desc: "err: an auto ops rules failure only fails the goals",
			setup: func(m *linkedTestMocks) {
				m.cache.EXPECT().Get(linkedTestEnvID+":experiment").Return(marshalLinkedTest(t, experiments), nil)
				m.cache.EXPECT().Get(linkedTestEnvID+":autoOpsRule").Return(nil, cache.ErrNotFound)
				m.autoOpsClient.EXPECT().ListAutoOpsRules(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))
			},
			goals:             goals,
			expectedPublished: map[string][]string{topicEvaluationDWH: {"e-1"}},
			expectedErrs:      map[string]error{"g-1": errLinkedTargetsLoadFailed},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			m := newLinkedTestMocks(ctrl)
			p.setup(m)
			m.expectPublish(nil)
			lp := NewLinkedEventPublisher(
				LinkedEventTopics{EvaluationDWH: m.evalDWH, GoalOPS: m.goalOPS},
				m.cache, m.expClient, m.autoOpsClient, zap.NewNop(),
			)
			assert.Equal(t, p.expectedErrs, lp.Publish(context.Background(), linkedTestEnvID, p.goals, evaluations))
			assert.Equal(t, p.expectedPublished, m.published)
		})
	}
}
