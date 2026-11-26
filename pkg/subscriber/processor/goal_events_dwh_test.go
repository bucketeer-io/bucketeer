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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	cachemock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	ecdwh "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/storage/v2/dwh_database"
	ecdwhmock "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/storage/v2/dwh_database/mock"
	redismock "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3/mock"
	dwhstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/dwhstorage/mock"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/v2/proto/eventpersisterdwh"
	exproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

const (
	testEnvironmentID = "env-1"
	testGoalID        = "goal-1"
	testUserID        = "user-1"
	testFeatureID     = "feature-1"
	testEventID       = "event-1"
)

func newTestGoalEvtWriter(
	eventStorage *ecdwhmock.MockEventStorage,
	redisClient *redismock.MockClient,
	cache *cachemock.MockExperimentsCache,
) *goalEvtWriter {
	return &goalEvtWriter{
		eventStorage:            eventStorage,
		redisClient:             redisClient,
		cache:                   cache,
		location:                time.UTC,
		logger:                  zap.NewNop(),
		maxRetryGoalEventPeriod: time.Hour,
		retryGoalEventInterval:  time.Minute,
		subscriberType:          subscriberGoalEventDWH,
	}
}

func newTestExperiment(now int64) *exproto.Experiment {
	return &exproto.Experiment{
		Id:             "experiment-1",
		GoalIds:        []string{testGoalID},
		FeatureId:      testFeatureID,
		FeatureVersion: 1,
		StartAt:        now - 3600,
		StopAt:         now + 3600,
	}
}

func newTestGoalEvent(timestamp int64) *eventproto.GoalEvent {
	return &eventproto.GoalEvent{
		GoalId:    testGoalID,
		UserId:    testUserID,
		Timestamp: timestamp,
		Tag:       "tag",
		Value:     1.0,
	}
}

func TestLinkGoalEventByExperimentLinksWhenEvaluationIsOlder(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now().Unix()
	eventStorage := ecdwhmock.NewMockEventStorage(ctrl)
	redisClient := redismock.NewMockClient(ctrl)
	w := newTestGoalEvtWriter(eventStorage, redisClient, nil)

	// The user was evaluated before the goal event was created: the event must link.
	eventStorage.EXPECT().QueryUserEvaluation(
		gomock.Any(), testEnvironmentID, testUserID, testFeatureID, int32(1), gomock.Any(), gomock.Any(),
	).Return(&ecdwh.UserEvaluation{
		UserID:         testUserID,
		FeatureID:      testFeatureID,
		FeatureVersion: 1,
		VariationID:    "variation-a",
		Reason:         "TARGET",
		Timestamp:      now - 60,
	}, nil)

	evals, retriable, err := w.linkGoalEventByExperiment(
		context.Background(),
		newTestGoalEvent(now),
		testEventID,
		testEnvironmentID,
		"tag",
		[]*exproto.Experiment{newTestExperiment(now)},
		true,
	)
	assert.NoError(t, err)
	assert.False(t, retriable)
	assert.Len(t, evals, 1)
	assert.Equal(t, "variation-a", evals[0].VariationID)
}

func TestLinkGoalEventByExperimentDiscardsGoalEventOlderThanEvaluation(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now().Unix()
	eventStorage := ecdwhmock.NewMockEventStorage(ctrl)
	// No expectations are set on the Redis client: storing a retry message
	// for this case would fail the test.
	redisClient := redismock.NewMockClient(ctrl)
	w := newTestGoalEvtWriter(eventStorage, redisClient, nil)

	// The latest evaluation is newer than the goal event: the client SDK sent
	// the goal event before evaluating the user, so it must be discarded.
	eventStorage.EXPECT().QueryUserEvaluation(
		gomock.Any(), testEnvironmentID, testUserID, testFeatureID, int32(1), gomock.Any(), gomock.Any(),
	).Return(&ecdwh.UserEvaluation{
		UserID:         testUserID,
		FeatureID:      testFeatureID,
		FeatureVersion: 1,
		VariationID:    "variation-a",
		Reason:         "TARGET",
		Timestamp:      now + 60,
	}, nil)

	evals, retriable, err := w.linkGoalEventByExperiment(
		context.Background(),
		newTestGoalEvent(now),
		testEventID,
		testEnvironmentID,
		"tag",
		[]*exproto.Experiment{newTestExperiment(now)},
		true,
	)
	assert.ErrorIs(t, err, ErrGoalEventOlderThanEvaluation)
	assert.False(t, retriable)
	assert.Nil(t, evals)
}

func TestLinkGoalEventByExperimentStoresRetryMessageWhenEvaluationNotFound(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now().Unix()
	eventStorage := ecdwhmock.NewMockEventStorage(ctrl)
	redisClient := redismock.NewMockClient(ctrl)
	w := newTestGoalEvtWriter(eventStorage, redisClient, nil)

	// The evaluation hasn't landed in the DWH yet (e.g. pub/sub delivered the
	// goal event first): a retry message must be stored.
	eventStorage.EXPECT().QueryUserEvaluation(
		gomock.Any(), testEnvironmentID, testUserID, testFeatureID, int32(1), gomock.Any(), gomock.Any(),
	).Return(nil, ecdwh.ErrBQNoResultsFound)
	redisClient.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	evals, retriable, err := w.linkGoalEventByExperiment(
		context.Background(),
		newTestGoalEvent(now),
		testEventID,
		testEnvironmentID,
		"tag",
		[]*exproto.Experiment{newTestExperiment(now)},
		true,
	)
	assert.True(t, ecdwh.IsNoResultsFound(err))
	assert.False(t, retriable)
	assert.Nil(t, evals)
}

func TestLinkGoalEventByExperimentDoesNotStoreRetryMessageForRetryProcessor(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now().Unix()
	eventStorage := ecdwhmock.NewMockEventStorage(ctrl)
	// No expectations are set on the Redis client: the retry processor owns its
	// retry message, so the linker must not overwrite it (that would reset the
	// backoff counter).
	redisClient := redismock.NewMockClient(ctrl)
	w := newTestGoalEvtWriter(eventStorage, redisClient, nil)

	eventStorage.EXPECT().QueryUserEvaluation(
		gomock.Any(), testEnvironmentID, testUserID, testFeatureID, int32(1), gomock.Any(), gomock.Any(),
	).Return(nil, ecdwh.ErrBQNoResultsFound)

	evals, retriable, err := w.linkGoalEventByExperiment(
		context.Background(),
		newTestGoalEvent(now),
		testEventID,
		testEnvironmentID,
		"tag",
		[]*exproto.Experiment{newTestExperiment(now)},
		false,
	)
	assert.True(t, ecdwh.IsNoResultsFound(err))
	assert.False(t, retriable)
	assert.Nil(t, evals)
}

func TestHandleNewRetryDiscardsGoalEventOlderThanEvaluation(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now().Unix()
	eventStorage := ecdwhmock.NewMockEventStorage(ctrl)
	redisClient := redismock.NewMockClient(ctrl)
	cache := cachemock.NewMockExperimentsCache(ctrl)
	w := newTestGoalEvtWriter(eventStorage, redisClient, cache)

	cache.EXPECT().Get(testEnvironmentID).Return(&exproto.Experiments{
		Experiments: []*exproto.Experiment{newTestExperiment(now)},
	}, nil)
	eventStorage.EXPECT().QueryUserEvaluation(
		gomock.Any(), testEnvironmentID, testUserID, testFeatureID, int32(1), gomock.Any(), gomock.Any(),
	).Return(&ecdwh.UserEvaluation{
		UserID:         testUserID,
		FeatureID:      testFeatureID,
		FeatureVersion: 1,
		VariationID:    "variation-a",
		Reason:         "TARGET",
		Timestamp:      now + 60,
	}, nil)
	// The retry message must be deleted instead of re-stored:
	// retrying can never link a goal event older than the latest evaluation.
	key := testEnvironmentID + ":" + retryGoalEventKeyKind + ":" + testEventID
	redisClient.EXPECT().Del(key).Return(nil).Times(1)

	msg := &retryMessage{
		GoalEvent:     newTestGoalEvent(now),
		EnvironmentID: testEnvironmentID,
		RetryCount:    1,
		ID:            testEventID,
		FirstRetryAt:  now - 60,
		RetryAt:       now - 30,
	}
	w.handleNewRetry(context.Background(), msg, key, w.logger)
}

func TestHandleNewRetryRequeuesWhenEvaluationNotFound(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now().Unix()
	eventStorage := ecdwhmock.NewMockEventStorage(ctrl)
	redisClient := redismock.NewMockClient(ctrl)
	cache := cachemock.NewMockExperimentsCache(ctrl)
	w := newTestGoalEvtWriter(eventStorage, redisClient, cache)

	cache.EXPECT().Get(testEnvironmentID).Return(&exproto.Experiments{
		Experiments: []*exproto.Experiment{newTestExperiment(now)},
	}, nil)
	eventStorage.EXPECT().QueryUserEvaluation(
		gomock.Any(), testEnvironmentID, testUserID, testFeatureID, int32(1), gomock.Any(), gomock.Any(),
	).Return(nil, ecdwh.ErrBQNoResultsFound)
	// The evaluation may still land later, so the retry message must be
	// re-stored exactly once with the incremented retry count.
	key := testEnvironmentID + ":" + retryGoalEventKeyKind + ":" + testEventID
	redisClient.EXPECT().Set(key, gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ string, value interface{}, _ time.Duration) error {
			data, ok := value.([]byte)
			if !ok {
				return errors.New("unexpected value type")
			}
			assert.Contains(t, string(data), `"retryCount":2`)
			return nil
		},
	).Times(1)

	msg := &retryMessage{
		GoalEvent:     newTestGoalEvent(now),
		EnvironmentID: testEnvironmentID,
		RetryCount:    1,
		ID:            testEventID,
		FirstRetryAt:  now - 60,
		RetryAt:       now - 30,
	}
	w.handleNewRetry(context.Background(), msg, key, w.logger)
}

// newTestHoldoutExperiment is an internal experiment of a holdout group that
// has the same goal as newTestExperiment. Its feature ID is the group ID.
func newTestHoldoutExperiment(now int64) *exproto.Experiment {
	return &exproto.Experiment{
		Id:             "holdout-experiment-1",
		Kind:           exproto.Experiment_HOLDOUT,
		GoalIds:        []string{testGoalID},
		FeatureId:      "holdout-group-1",
		FeatureVersion: 1,
		StartAt:        now - 3600,
		StopAt:         now + 3600,
		Status:         exproto.Experiment_RUNNING,
	}
}

// An internal experiment of a holdout group is not linked. A flag experiment
// that has the same goal is linked as if the internal experiment were absent:
// the only evaluation that is looked up is the flag's. gomock fails the test
// on a lookup with the group ID as the feature ID.
func TestLinkGoalEventByExperimentSkipsHoldoutExperiment(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	flagEvaluation := &ecdwh.UserEvaluation{
		UserID:         testUserID,
		FeatureID:      testFeatureID,
		FeatureVersion: 1,
		VariationID:    "variation-a",
		Reason:         "TARGET",
		Timestamp:      now - 60,
	}
	orders := map[string][]*exproto.Experiment{
		"holdout experiment first": {newTestHoldoutExperiment(now), newTestExperiment(now)},
		"flag experiment first":    {newTestExperiment(now), newTestHoldoutExperiment(now)},
	}
	for name, experiments := range orders {
		t.Run("a flag experiment sharing the goal is linked: "+name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			eventStorage := ecdwhmock.NewMockEventStorage(ctrl)
			redisClient := redismock.NewMockClient(ctrl)
			w := newTestGoalEvtWriter(eventStorage, redisClient, nil)
			eventStorage.EXPECT().QueryUserEvaluation(
				gomock.Any(), testEnvironmentID, testUserID, testFeatureID, int32(1), gomock.Any(), gomock.Any(),
			).Return(flagEvaluation, nil).Times(1)

			evals, retriable, err := w.linkGoalEventByExperiment(
				context.Background(), newTestGoalEvent(now), testEventID, testEnvironmentID, "tag", experiments, true,
			)
			assert.NoError(t, err)
			assert.False(t, retriable)
			assert.Len(t, evals, 1)
			assert.Equal(t, testFeatureID, evals[0].FeatureID)
			assert.Equal(t, "variation-a", evals[0].VariationID)
		})
	}

	for _, storeRetryOnMiss := range []bool{true, false} {
		name := "retry path"
		if storeRetryOnMiss {
			name = "fresh path"
		}
		t.Run("only a holdout experiment has the goal: "+name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			// No evaluation is looked up and no retry message is stored.
			eventStorage := ecdwhmock.NewMockEventStorage(ctrl)
			redisClient := redismock.NewMockClient(ctrl)
			w := newTestGoalEvtWriter(eventStorage, redisClient, nil)

			evals, retriable, err := w.linkGoalEventByExperiment(
				context.Background(),
				newTestGoalEvent(now),
				testEventID,
				testEnvironmentID,
				"tag",
				[]*exproto.Experiment{newTestHoldoutExperiment(now)},
				storeRetryOnMiss,
			)
			assert.ErrorIs(t, err, ErrExperimentNotFound)
			assert.False(t, retriable)
			assert.Nil(t, evals)
		})
	}

	t.Run("a missing flag evaluation stores one retry and looks up no holdout evaluation", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		eventStorage := ecdwhmock.NewMockEventStorage(ctrl)
		redisClient := redismock.NewMockClient(ctrl)
		w := newTestGoalEvtWriter(eventStorage, redisClient, nil)
		eventStorage.EXPECT().QueryUserEvaluation(
			gomock.Any(), testEnvironmentID, testUserID, testFeatureID, int32(1), gomock.Any(), gomock.Any(),
		).Return(nil, ecdwh.ErrBQNoResultsFound).Times(1)
		redisClient.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		evals, _, err := w.linkGoalEventByExperiment(
			context.Background(),
			newTestGoalEvent(now),
			testEventID,
			testEnvironmentID,
			"tag",
			[]*exproto.Experiment{newTestHoldoutExperiment(now), newTestExperiment(now)},
			true,
		)
		assert.True(t, ecdwh.IsNoResultsFound(err))
		assert.Nil(t, evals)
	})
}

// The retry processor reads the experiments from the cache, which holds the
// internal experiments too, and links the retried goal event to the flag
// experiment only.
func TestHandleNewRetrySkipsHoldoutExperiment(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	key := testEnvironmentID + ":" + retryGoalEventKeyKind + ":" + testEventID
	newMessage := func() *retryMessage {
		return &retryMessage{
			GoalEvent:     newTestGoalEvent(now),
			EnvironmentID: testEnvironmentID,
			RetryCount:    1,
			ID:            testEventID,
			FirstRetryAt:  now - 60,
			RetryAt:       now - 30,
		}
	}

	t.Run("a flag experiment sharing the goal is linked and written", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		eventStorage := ecdwhmock.NewMockEventStorage(ctrl)
		redisClient := redismock.NewMockClient(ctrl)
		cache := cachemock.NewMockExperimentsCache(ctrl)
		goalWriter := dwhstoragemock.NewMockGoalEventWriter(ctrl)
		w := newTestGoalEvtWriter(eventStorage, redisClient, cache)
		w.writer = goalWriter

		cache.EXPECT().Get(testEnvironmentID).Return(&exproto.Experiments{
			Experiments: []*exproto.Experiment{newTestHoldoutExperiment(now), newTestExperiment(now)},
		}, nil)
		eventStorage.EXPECT().QueryUserEvaluation(
			gomock.Any(), testEnvironmentID, testUserID, testFeatureID, int32(1), gomock.Any(), gomock.Any(),
		).Return(&ecdwh.UserEvaluation{
			UserID:         testUserID,
			FeatureID:      testFeatureID,
			FeatureVersion: 1,
			VariationID:    "variation-a",
			Reason:         "TARGET",
			Timestamp:      now - 60,
		}, nil).Times(1)
		goalWriter.EXPECT().AppendRows(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, events []*epproto.GoalEvent) (map[string]bool, error) {
				if assert.Len(t, events, 1) {
					assert.Equal(t, testFeatureID, events[0].FeatureId)
					assert.Equal(t, "variation-a", events[0].VariationId)
				}
				return nil, nil
			},
		)
		redisClient.EXPECT().Del(key).Return(nil).Times(1)

		w.handleNewRetry(context.Background(), newMessage(), key, w.logger)
	})

	t.Run("only a holdout experiment has the goal: the retry message is deleted", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		// No evaluation is looked up and the message is not stored again.
		eventStorage := ecdwhmock.NewMockEventStorage(ctrl)
		redisClient := redismock.NewMockClient(ctrl)
		cache := cachemock.NewMockExperimentsCache(ctrl)
		w := newTestGoalEvtWriter(eventStorage, redisClient, cache)

		cache.EXPECT().Get(testEnvironmentID).Return(&exproto.Experiments{
			Experiments: []*exproto.Experiment{newTestHoldoutExperiment(now)},
		}, nil)
		redisClient.EXPECT().Del(key).Return(nil).Times(1)

		w.handleNewRetry(context.Background(), newMessage(), key, w.logger)
	})
}

func TestScanAndProcessVisitsKeysFromEveryScanPage(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	redisClient := redismock.NewMockClient(ctrl)
	w := newTestGoalEvtWriter(nil, redisClient, nil)
	w.locker = NewGoalEventLocker(redisClient, time.Minute)

	redisClient.EXPECT().ScanAll("*:goal_event_retry:*", int64(scanBatchSize), gomock.Any()).DoAndReturn(
		func(_ string, _ int64, fn func([]string) error) error {
			// One page per cluster master.
			if err := fn([]string{testEnvironmentID + ":goal_event_retry:event-1"}); err != nil {
				return err
			}
			return fn([]string{testEnvironmentID + ":goal_event_retry:event-2"})
		})
	// Each key reaches processing; the lock is not acquired, so it stops there.
	redisClient.EXPECT().SetNX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil).Times(2)

	w.scanAndProcess(context.Background())
}

func TestScanAndProcessScanError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	redisClient := redismock.NewMockClient(ctrl)
	w := newTestGoalEvtWriter(nil, redisClient, nil)

	redisClient.EXPECT().ScanAll(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("scan failed"))

	w.scanAndProcess(context.Background())
}
