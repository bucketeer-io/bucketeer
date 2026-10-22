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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	autoopsclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachemock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	aoproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
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
