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
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachemock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	experimentdomain "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	holdoutdomain "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/domain"
	holdoutmock "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2/mock"
	redismock "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3/mock"
	dwhstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/dwhstorage/mock"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/v2/proto/eventpersisterdwh"
	exproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

const (
	exposureEnvID   = "env-1"
	exposureGroupID = "hg-1"
	exposureFlagID  = "flag-1"
	// exposureDayStart is 00:00:00 UTC of a day.
	exposureDayStart = int64(1_700_006_400)
	exposureDay      = exposureDayStart / secondsPerDay
	// The measured group started the day before and stops two days later.
	// Its analysis period starts at noon of exposureDayStart's day.
	exposureGroupStartAt    = exposureDayStart - secondsPerDay
	exposureAnalysisStartAt = exposureDayStart + 12*60*60
	exposureGroupStopAt     = exposureDayStart + 2*secondsPerDay
	exposureFlagVersion     = int32(3)
)

// exposureUsers are user IDs picked by where they hash for
// exposureHashGroup: in the holdout arc, in the comparison range, in neither.
type exposureUserIDs struct {
	heldOut    string
	comparison string
	outside    string
}

func exposureHashGroup() *ftproto.HoldoutGroup {
	return &ftproto.HoldoutGroup{
		Id:                   exposureGroupID,
		EnvironmentId:        exposureEnvID,
		AssignmentType:       ftproto.HoldoutGroup_HASH,
		SamplingSeed:         "seed",
		BucketStart:          95,
		Percentage:           10,
		Status:               ftproto.HoldoutGroup_ANALYSIS,
		ExperimentId:         "exp-measurement",
		AnalysisExperimentId: "exp-analysis",
		StartAt:              exposureGroupStartAt,
		StopAt:               exposureGroupStopAt,
		Features: []*ftproto.HoldoutGroupFeature{
			{Id: "target-1", HoldoutGroupId: exposureGroupID, FeatureId: exposureFlagID, AddedAt: exposureGroupStartAt},
		},
	}
}

func exposureSegmentGroup() *ftproto.HoldoutGroup {
	group := exposureHashGroup()
	group.AssignmentType = ftproto.HoldoutGroup_SEGMENT
	group.SegmentId = "segment-1"
	group.SamplingSeed = ""
	group.BucketStart = 0
	group.Percentage = 0
	return group
}

func exposureUsers(t *testing.T) exposureUserIDs {
	t.Helper()
	group := exposureHashGroup()
	var users exposureUserIDs
	for i := 0; users.heldOut == "" || users.comparison == "" || users.outside == ""; i++ {
		require.Less(t, i, 10000)
		id := fmt.Sprintf("user-%d", i)
		heldOut := evaluation.IsHeldOut(group, &userproto.User{Id: id}, nil)
		comparison := evaluation.IsInHoldoutComparisonRange(group, id)
		switch {
		case heldOut && users.heldOut == "":
			users.heldOut = id
		case comparison && users.comparison == "":
			users.comparison = id
		case !heldOut && !comparison && users.outside == "":
			users.outside = id
		}
	}
	return users
}

// The two internal experiments of the measured group. They share the feature
// ID and version; only their windows differ.
func exposureMeasurementExperiment() *exproto.Experiment {
	return &exproto.Experiment{
		Id:             "exp-measurement",
		Kind:           exproto.Experiment_HOLDOUT,
		FeatureId:      exposureGroupID,
		FeatureVersion: holdoutdomain.InternalExperimentFeatureVersion,
		Status:         exproto.Experiment_RUNNING,
		StartAt:        exposureGroupStartAt,
		StopAt:         exposureGroupStopAt,
	}
}

func exposureAnalysisExperiment() *exproto.Experiment {
	e := exposureMeasurementExperiment()
	e.Id = "exp-analysis"
	e.StartAt = exposureAnalysisStartAt
	return e
}

func exposureFlagExperiment() *exproto.Experiment {
	return &exproto.Experiment{
		Id:             "exp-flag",
		FeatureId:      exposureFlagID,
		FeatureVersion: exposureFlagVersion,
		Status:         exproto.Experiment_RUNNING,
		StartAt:        exposureGroupStartAt,
		StopAt:         exposureGroupStopAt,
	}
}

func exposureEvent(userID string, timestamp int64, reason *ftproto.Reason) *eventproto.EvaluationEvent {
	return &eventproto.EvaluationEvent{
		Timestamp:      timestamp,
		FeatureId:      exposureFlagID,
		FeatureVersion: exposureFlagVersion,
		UserId:         userID,
		VariationId:    "variation-served",
		User:           &userproto.User{Id: userID, Data: map[string]string{"plan": "pro"}},
		Reason:         reason,
		Tag:            "server",
		SourceId:       eventproto.SourceId_GO_SERVER,
	}
}

func heldOutReason(groupID string) *ftproto.Reason {
	return &ftproto.Reason{Type: ftproto.Reason_HOLDOUT, HoldoutGroupId: groupID}
}

func servedNormallyReason() *ftproto.Reason {
	return &ftproto.Reason{Type: ftproto.Reason_DEFAULT}
}

// fakeExposureKeys is the Redis side of the dedupe: SET NX and the
// compare-and-delete script, on a map.
type fakeExposureKeys struct {
	mu       sync.Mutex
	values   map[string]string
	ttls     map[string]time.Duration
	asked    []string
	released []string
	setErr   error
	evalErr  error
}

func newFakeExposureKeys() *fakeExposureKeys {
	return &fakeExposureKeys{values: map[string]string{}, ttls: map[string]time.Duration{}}
}

func (f *fakeExposureKeys) attach(client *redismock.MockClient) {
	client.EXPECT().SetNX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
			f.mu.Lock()
			defer f.mu.Unlock()
			f.asked = append(f.asked, key)
			if f.setErr != nil {
				return false, f.setErr
			}
			if _, ok := f.values[key]; ok {
				return false, nil
			}
			f.values[key] = value.(string)
			f.ttls[key] = ttl
			return true, nil
		},
	).AnyTimes()
	client.EXPECT().Eval(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, _ string, keys []string, args ...interface{}) *goredis.Cmd {
			f.mu.Lock()
			defer f.mu.Unlock()
			if f.evalErr != nil {
				return goredis.NewCmdResult(nil, f.evalErr)
			}
			if f.values[keys[0]] != args[0].(string) {
				return goredis.NewCmdResult(int64(0), nil)
			}
			delete(f.values, keys[0])
			delete(f.ttls, keys[0])
			f.released = append(f.released, keys[0])
			return goredis.NewCmdResult(int64(1), nil)
		},
	).AnyTimes()
}

func (f *fakeExposureKeys) held() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	keys := make([]string, 0, len(f.values))
	for key := range f.values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func exposureKey(periodStartAt, day int64, userID string) string {
	return fmt.Sprintf("%s:holdout_exposure:%s:%d:%d:%s", exposureEnvID, exposureGroupID, periodStartAt, day, userID)
}

type evalWriterHarness struct {
	t                   *testing.T
	writer              *evalEvtWriter
	experimentsCache    *cachemock.MockExperimentsCache
	experimentClient    *experimentclientmock.MockClient
	holdoutStorage      *holdoutmock.MockHoldoutStorage
	measuredGroupsCache *cachemock.MockHoldoutMeasuredGroupsCache
	redis               *redismock.MockClient
	dwh                 *dwhstoragemock.MockEvalEventWriter
	// appended are the rows of every AppendRows call, in call order.
	appended [][]*epproto.EvaluationEvent
}

// newEvalWriterHarness builds the writer on mocks. Nothing is expected of
// the holdout storage, the measured groups cache or Redis: a test that does
// not set an expectation on one of them fails when the writer calls it.
func newEvalWriterHarness(t *testing.T) *evalWriterHarness {
	t.Helper()
	controller := gomock.NewController(t)
	h := &evalWriterHarness{
		t:                   t,
		experimentsCache:    cachemock.NewMockExperimentsCache(controller),
		experimentClient:    experimentclientmock.NewMockClient(controller),
		holdoutStorage:      holdoutmock.NewMockHoldoutStorage(controller),
		measuredGroupsCache: cachemock.NewMockHoldoutMeasuredGroupsCache(controller),
		redis:               redismock.NewMockClient(controller),
		dwh:                 dwhstoragemock.NewMockEvalEventWriter(controller),
	}
	writer := NewEvalEventWriter(
		zap.NewNop(),
		h.experimentClient,
		h.experimentsCache,
		h.holdoutStorage,
		h.measuredGroupsCache,
		h.redis,
		time.UTC,
		h.dwh,
		"test",
	)
	h.writer = writer.(*evalEvtWriter)
	return h
}

func (h *evalWriterHarness) withExperiments(experiments ...*exproto.Experiment) *evalWriterHarness {
	h.experimentsCache.EXPECT().Get(exposureEnvID).
		Return(&exproto.Experiments{Experiments: experiments}, nil).AnyTimes()
	return h
}

func (h *evalWriterHarness) withCachedGroups(groups ...*ftproto.HoldoutGroup) *evalWriterHarness {
	h.measuredGroupsCache.EXPECT().Get(exposureEnvID).
		Return(&ftproto.HoldoutGroups{HoldoutGroups: groups}, nil).AnyTimes()
	return h
}

func (h *evalWriterHarness) withKeys() *fakeExposureKeys {
	keys := newFakeExposureKeys()
	keys.attach(h.redis)
	return keys
}

// withAppend records the rows of every AppendRows call and answers with the
// given result.
func (h *evalWriterHarness) withAppend(fails map[string]bool, err error) *evalWriterHarness {
	h.dwh.EXPECT().AppendRows(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, rows []*epproto.EvaluationEvent) (map[string]bool, error) {
			h.appended = append(h.appended, rows)
			return fails, err
		},
	).AnyTimes()
	return h
}

func (h *evalWriterHarness) write(events eventDWHMap) map[string]bool {
	return h.writer.Write(context.Background(), environmentEventDWHMap{exposureEnvID: events})
}

// lastRows returns the rows of the last AppendRows call by row ID.
func (h *evalWriterHarness) lastRows() map[string]*epproto.EvaluationEvent {
	h.t.Helper()
	require.NotEmpty(h.t, h.appended)
	rows := make(map[string]*epproto.EvaluationEvent)
	for _, row := range h.appended[len(h.appended)-1] {
		_, duplicated := rows[row.Id]
		require.False(h.t, duplicated, "row ID %s appended twice", row.Id)
		rows[row.Id] = row
	}
	return rows
}

// A group without goals has no internal experiment, so the environment's
// experiments list has none of kind HOLDOUT. Its evaluation events then cost
// the pipeline nothing: the harness expects no read of the measured groups
// cache, no read of the holdout storage and no Redis call, and the only row
// written is the one of the flag experiment.
func TestEvalEvtWriterGroupWithoutGoalsAddsNothing(t *testing.T) {
	t.Parallel()
	users := exposureUsers(t)
	events := func() eventDWHMap {
		return eventDWHMap{
			"held-out":   exposureEvent(users.heldOut, exposureDayStart+10, heldOutReason(exposureGroupID)),
			"comparison": exposureEvent(users.comparison, exposureDayStart+20, servedNormallyReason()),
			"outside":    exposureEvent(users.outside, exposureDayStart+30, servedNormallyReason()),
		}
	}

	t.Run("the flag has an experiment", func(t *testing.T) {
		t.Parallel()
		h := newEvalWriterHarness(t).withExperiments(exposureFlagExperiment()).withAppend(nil, nil)
		fails := h.write(events())
		assert.Empty(t, fails)
		rows := h.lastRows()
		// The users served normally are written for the flag experiment.
		// The held-out user is not: the user did not get a variation of it.
		require.Len(t, rows, 2)
		for _, id := range []string{"comparison", "outside"} {
			require.Contains(t, rows, id)
			assert.Equal(t, exposureFlagID, rows[id].FeatureId)
			assert.Equal(t, exposureFlagVersion, rows[id].FeatureVersion)
		}
	})

	t.Run("the environment has no experiment", func(t *testing.T) {
		t.Parallel()
		h := newEvalWriterHarness(t).withExperiments().withAppend(nil, nil)
		fails := h.write(events())
		assert.Empty(t, fails)
		assert.Empty(t, h.lastRows())
	})

	t.Run("only another flag has an experiment", func(t *testing.T) {
		t.Parallel()
		other := exposureFlagExperiment()
		other.FeatureId = "flag-other"
		h := newEvalWriterHarness(t).withExperiments(other).withAppend(nil, nil)
		fails := h.write(events())
		assert.Empty(t, fails)
		assert.Empty(t, h.lastRows())
	})
}

func TestEvalEvtWriterHoldoutExposureVariations(t *testing.T) {
	t.Parallel()
	users := exposureUsers(t)
	timestamp := exposureDayStart + 10

	patterns := []struct {
		desc      string
		group     *ftproto.HoldoutGroup
		userID    string
		reason    *ftproto.Reason
		variation string
	}{
		{
			desc:      "hash: the HOLDOUT reason of the group is in_holdout",
			group:     exposureHashGroup(),
			userID:    users.heldOut,
			reason:    heldOutReason(exposureGroupID),
			variation: holdoutdomain.InHoldoutVariationID,
		},
		{
			// The reason decides, not the hash: this user is outside both
			// arcs of the group as the subscriber sees it.
			desc:      "hash: the HOLDOUT reason of the group is in_holdout wherever the user hashes",
			group:     exposureHashGroup(),
			userID:    users.outside,
			reason:    heldOutReason(exposureGroupID),
			variation: holdoutdomain.InHoldoutVariationID,
		},
		{
			desc:   "hash: the HOLDOUT reason of another group counts for neither",
			group:  exposureHashGroup(),
			userID: users.heldOut,
			reason: heldOutReason("hg-other"),
		},
		{
			desc:   "hash: the HOLDOUT reason without a group counts for neither",
			group:  exposureHashGroup(),
			userID: users.heldOut,
			reason: heldOutReason(""),
		},
		{
			desc:      "hash: served normally in the comparison range is not_in_holdout",
			group:     exposureHashGroup(),
			userID:    users.comparison,
			reason:    servedNormallyReason(),
			variation: holdoutdomain.NotInHoldoutVariationID,
		},
		{
			desc:   "hash: served normally outside the comparison range counts for neither",
			group:  exposureHashGroup(),
			userID: users.outside,
			reason: servedNormallyReason(),
		},
		{
			// An SDK that does not report the HOLDOUT reason: the user is
			// not counted as held out, and is not in the comparison range.
			desc:   "hash: served normally in the holdout arc counts for neither",
			group:  exposureHashGroup(),
			userID: users.heldOut,
			reason: servedNormallyReason(),
		},
		{
			desc:      "hash: an event without a reason in the comparison range is not_in_holdout",
			group:     exposureHashGroup(),
			userID:    users.comparison,
			reason:    nil,
			variation: holdoutdomain.NotInHoldoutVariationID,
		},
		{
			desc:      "segment: the HOLDOUT reason of the group is in_holdout",
			group:     exposureSegmentGroup(),
			userID:    users.outside,
			reason:    heldOutReason(exposureGroupID),
			variation: holdoutdomain.InHoldoutVariationID,
		},
		{
			desc:      "segment: every user served normally is not_in_holdout",
			group:     exposureSegmentGroup(),
			userID:    users.outside,
			reason:    servedNormallyReason(),
			variation: holdoutdomain.NotInHoldoutVariationID,
		},
		{
			desc:   "segment: the HOLDOUT reason of another group counts for neither",
			group:  exposureSegmentGroup(),
			userID: users.outside,
			reason: heldOutReason("hg-other"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			h := newEvalWriterHarness(t).
				withExperiments(exposureMeasurementExperiment()).
				withCachedGroups(p.group).
				withAppend(nil, nil)
			var keys *fakeExposureKeys
			if p.variation != "" {
				// Without a variation no Redis call is expected at all.
				keys = h.withKeys()
			}
			fails := h.write(eventDWHMap{"evt": exposureEvent(p.userID, timestamp, p.reason)})
			assert.Empty(t, fails)
			rows := h.lastRows()
			if p.variation == "" {
				assert.Empty(t, rows)
				return
			}
			require.Len(t, rows, 1)
			row := rows["evt:"+exposureGroupID]
			require.NotNil(t, row)
			assert.Equal(t, &epproto.EvaluationEvent{
				Id:             "evt:" + exposureGroupID,
				FeatureId:      exposureGroupID,
				FeatureVersion: holdoutdomain.InternalExperimentFeatureVersion,
				UserData:       `{"plan":"pro"}`,
				UserId:         p.userID,
				VariationId:    p.variation,
				Reason:         p.reason.GetType().String(),
				Tag:            "server",
				SourceId:       eventproto.SourceId_GO_SERVER.String(),
				EnvironmentId:  exposureEnvID,
				Timestamp:      time.Unix(timestamp, 0).UnixMicro(),
			}, row)
			key := exposureKey(exposureGroupStartAt, exposureDay, p.userID)
			assert.Equal(t, []string{key}, keys.held())
			assert.Equal(t, 25*time.Hour, keys.ttls[key])
		})
	}
}

func TestEvalEvtWriterHoldoutExposureNeedsUserID(t *testing.T) {
	t.Parallel()
	h := newEvalWriterHarness(t).
		withExperiments(exposureMeasurementExperiment()).
		withCachedGroups(exposureSegmentGroup()).
		withAppend(nil, nil)
	keys := h.withKeys()
	users := exposureUsers(t)

	// The user ID is taken from the user of the event when the event has no
	// user ID of its own.
	fromUser := exposureEvent(users.outside, exposureDayStart+10, servedNormallyReason())
	fromUser.UserId = ""
	// An event without any user ID counts for nothing.
	without := exposureEvent("", exposureDayStart+10, servedNormallyReason())
	without.User = nil

	fails := h.write(eventDWHMap{"from-user": fromUser, "without": without})
	assert.Empty(t, fails)
	rows := h.lastRows()
	require.Len(t, rows, 1)
	require.Contains(t, rows, "from-user:"+exposureGroupID)
	assert.Equal(t, users.outside, rows["from-user:"+exposureGroupID].UserId)
	assert.Equal(t, []string{exposureKey(exposureGroupStartAt, exposureDay, users.outside)}, keys.held())
}

func TestEvalEvtWriterHoldoutExposureTargetsAndWindow(t *testing.T) {
	t.Parallel()
	users := exposureUsers(t)
	removedAt := exposureDayStart + 100
	readdedAt := exposureDayStart + 200

	patterns := []struct {
		desc        string
		features    []*ftproto.HoldoutGroupFeature
		experiments []*exproto.Experiment
		featureID   string
		timestamp   int64
		written     bool
	}{
		{
			desc:      "a flag the group does not target",
			featureID: "flag-other",
			timestamp: exposureDayStart + 10,
		},
		{
			desc:      "at the time the flag was added",
			timestamp: exposureGroupStartAt,
			written:   true,
		},
		{
			desc: "before the flag was added",
			features: []*ftproto.HoldoutGroupFeature{
				{FeatureId: exposureFlagID, AddedAt: exposureDayStart + 50},
			},
			timestamp: exposureDayStart + 49,
		},
		{
			desc: "just before the flag was removed",
			features: []*ftproto.HoldoutGroupFeature{
				{FeatureId: exposureFlagID, AddedAt: exposureGroupStartAt, RemovedAt: removedAt},
			},
			timestamp: removedAt - 1,
			written:   true,
		},
		{
			desc: "at the time the flag was removed",
			features: []*ftproto.HoldoutGroupFeature{
				{FeatureId: exposureFlagID, AddedAt: exposureGroupStartAt, RemovedAt: removedAt},
			},
			timestamp: removedAt,
		},
		{
			desc: "between a removal and a new addition of the flag",
			features: []*ftproto.HoldoutGroupFeature{
				{FeatureId: exposureFlagID, AddedAt: exposureGroupStartAt, RemovedAt: removedAt},
				{FeatureId: exposureFlagID, AddedAt: readdedAt},
			},
			timestamp: readdedAt - 1,
		},
		{
			desc: "after the flag was added again",
			features: []*ftproto.HoldoutGroupFeature{
				{FeatureId: exposureFlagID, AddedAt: exposureGroupStartAt, RemovedAt: removedAt},
				{FeatureId: exposureFlagID, AddedAt: readdedAt},
			},
			timestamp: readdedAt,
			written:   true,
		},
		{
			desc:      "before the internal experiment started",
			timestamp: exposureGroupStartAt - 1,
			features: []*ftproto.HoldoutGroupFeature{
				{FeatureId: exposureFlagID, AddedAt: exposureGroupStartAt - 100},
			},
		},
		{
			desc:      "the last second of the internal experiment",
			timestamp: exposureGroupStopAt - 1,
			written:   true,
		},
		{
			desc:      "at the time the internal experiment stopped",
			timestamp: exposureGroupStopAt,
		},
		{
			// A group stopped by hand: its experiment stops at the time of
			// the stop, and the flag is no longer served from the holdout.
			desc: "after an internal experiment that was stopped early",
			experiments: []*exproto.Experiment{{
				Id:             "exp-measurement",
				Kind:           exproto.Experiment_HOLDOUT,
				FeatureId:      exposureGroupID,
				FeatureVersion: holdoutdomain.InternalExperimentFeatureVersion,
				Status:         exproto.Experiment_STOPPED,
				StartAt:        exposureGroupStartAt,
				StopAt:         exposureDayStart,
			}},
			timestamp: exposureDayStart + 10,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			group := exposureHashGroup()
			if p.features != nil {
				group.Features = p.features
			}
			experiments := p.experiments
			if experiments == nil {
				experiments = []*exproto.Experiment{exposureMeasurementExperiment()}
			}
			h := newEvalWriterHarness(t).
				withExperiments(experiments...).
				withCachedGroups(group).
				withAppend(nil, nil)
			if p.written {
				h.withKeys()
			}
			evt := exposureEvent(users.comparison, p.timestamp, servedNormallyReason())
			if p.featureID != "" {
				evt.FeatureId = p.featureID
			}
			fails := h.write(eventDWHMap{"evt": evt})
			assert.Empty(t, fails)
			if p.written {
				assert.Len(t, h.lastRows(), 1)
			} else {
				assert.Empty(t, h.lastRows())
			}
		})
	}
}

func TestHoldoutMeasurementPeriod(t *testing.T) {
	t.Parallel()
	both := []*exproto.Experiment{exposureMeasurementExperiment(), exposureAnalysisExperiment()}
	reversed := []*exproto.Experiment{exposureAnalysisExperiment(), exposureMeasurementExperiment()}
	patterns := []struct {
		desc        string
		experiments []*exproto.Experiment
		timestamp   int64
		startAt     int64
		found       bool
	}{
		{desc: "no experiment", timestamp: exposureDayStart},
		{desc: "before the first window", experiments: both, timestamp: exposureGroupStartAt - 1},
		{
			desc:        "start is inclusive",
			experiments: both,
			timestamp:   exposureGroupStartAt,
			startAt:     exposureGroupStartAt,
			found:       true,
		},
		{
			desc:        "before the analysis period only the first window contains it",
			experiments: both,
			timestamp:   exposureAnalysisStartAt - 1,
			startAt:     exposureGroupStartAt,
			found:       true,
		},
		{
			desc:        "in both windows the newest start wins",
			experiments: both,
			timestamp:   exposureAnalysisStartAt,
			startAt:     exposureAnalysisStartAt,
			found:       true,
		},
		{
			desc:        "the order of the experiments does not matter",
			experiments: reversed,
			timestamp:   exposureAnalysisStartAt,
			startAt:     exposureAnalysisStartAt,
			found:       true,
		},
		{desc: "stop is exclusive", experiments: both, timestamp: exposureGroupStopAt},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			startAt, found := holdoutMeasurementPeriod(p.experiments, p.timestamp)
			assert.Equal(t, p.found, found)
			assert.Equal(t, p.startAt, startAt)
		})
	}
}

func TestEvalEvtWriterHoldoutExposureOneRowPerUserPeriodAndDay(t *testing.T) {
	t.Parallel()
	users := exposureUsers(t)
	newHarness := func(t *testing.T) (*evalWriterHarness, *fakeExposureKeys) {
		h := newEvalWriterHarness(t).
			withExperiments(exposureMeasurementExperiment(), exposureAnalysisExperiment()).
			withCachedGroups(exposureHashGroup()).
			withAppend(nil, nil)
		return h, h.withKeys()
	}

	t.Run("the events of a user in one batch write one row and ask Redis once", func(t *testing.T) {
		t.Parallel()
		h, keys := newHarness(t)
		fails := h.write(eventDWHMap{
			"evt-1": exposureEvent(users.comparison, exposureDayStart+10, servedNormallyReason()),
			"evt-2": exposureEvent(users.comparison, exposureDayStart+20, servedNormallyReason()),
			"evt-3": exposureEvent(users.comparison, exposureDayStart+30, servedNormallyReason()),
		})
		assert.Empty(t, fails)
		assert.Len(t, h.lastRows(), 1)
		assert.Equal(t, []string{exposureKey(exposureGroupStartAt, exposureDay, users.comparison)}, keys.asked)
	})

	t.Run("a later batch of the same day writes no row and nacks nothing", func(t *testing.T) {
		t.Parallel()
		h, _ := newHarness(t)
		h.write(eventDWHMap{"evt-1": exposureEvent(users.comparison, exposureDayStart+10, servedNormallyReason())})
		require.Len(t, h.lastRows(), 1)
		fails := h.write(eventDWHMap{
			"evt-2": exposureEvent(users.comparison, exposureDayStart+20, servedNormallyReason()),
		})
		assert.Empty(t, fails)
		assert.Empty(t, h.lastRows())
	})

	t.Run("a second delivery of the event writes no row", func(t *testing.T) {
		t.Parallel()
		h, _ := newHarness(t)
		evt := exposureEvent(users.heldOut, exposureDayStart+10, heldOutReason(exposureGroupID))
		h.write(eventDWHMap{"evt": evt})
		require.Len(t, h.lastRows(), 1)
		fails := h.write(eventDWHMap{"evt": evt})
		assert.Empty(t, fails)
		assert.Empty(t, h.lastRows())
	})

	t.Run("each user has a row", func(t *testing.T) {
		t.Parallel()
		h, _ := newHarness(t)
		fails := h.write(eventDWHMap{
			"evt-1": exposureEvent(users.comparison, exposureDayStart+10, servedNormallyReason()),
			"evt-2": exposureEvent(users.heldOut, exposureDayStart+10, heldOutReason(exposureGroupID)),
		})
		assert.Empty(t, fails)
		rows := h.lastRows()
		require.Len(t, rows, 2)
		assert.Equal(t, holdoutdomain.NotInHoldoutVariationID, rows["evt-1:"+exposureGroupID].VariationId)
		assert.Equal(t, holdoutdomain.InHoldoutVariationID, rows["evt-2:"+exposureGroupID].VariationId)
	})

	t.Run("the next day writes a row", func(t *testing.T) {
		t.Parallel()
		h, keys := newHarness(t)
		fails := h.write(eventDWHMap{
			"evt-1": exposureEvent(users.comparison, exposureDayStart-1, servedNormallyReason()),
			"evt-2": exposureEvent(users.comparison, exposureDayStart, servedNormallyReason()),
		})
		assert.Empty(t, fails)
		assert.Len(t, h.lastRows(), 2)
		assert.Equal(t, []string{
			exposureKey(exposureGroupStartAt, exposureDay-1, users.comparison),
			exposureKey(exposureGroupStartAt, exposureDay, users.comparison),
		}, keys.held())
	})

	// The analysis-period experiment reads rows from its own start. A user
	// whose row of the day is from before that start gets a second row at
	// the first event after it, or that experiment would miss the user on
	// its first day.
	t.Run("the start of the analysis period writes a row on the same day", func(t *testing.T) {
		t.Parallel()
		h, keys := newHarness(t)
		fails := h.write(eventDWHMap{
			"before": exposureEvent(users.comparison, exposureAnalysisStartAt-1, servedNormallyReason()),
			"after":  exposureEvent(users.comparison, exposureAnalysisStartAt, servedNormallyReason()),
			"later":  exposureEvent(users.comparison, exposureAnalysisStartAt+60, servedNormallyReason()),
		})
		assert.Empty(t, fails)
		rows := h.lastRows()
		require.Len(t, rows, 2)
		require.Contains(t, rows, "before:"+exposureGroupID)
		assert.Equal(t, []string{
			exposureKey(exposureGroupStartAt, exposureDay, users.comparison),
			exposureKey(exposureAnalysisStartAt, exposureDay, users.comparison),
		}, keys.held())
		var afterStart int
		for _, row := range rows {
			if row.Timestamp >= time.Unix(exposureAnalysisStartAt, 0).UnixMicro() {
				afterStart++
			}
		}
		assert.Equal(t, 1, afterStart, "the analysis-period experiment sees the user once")
	})
}

func TestEvalEvtWriterHoldoutExposureWithFlagExperiment(t *testing.T) {
	t.Parallel()
	users := exposureUsers(t)
	newHarness := func(t *testing.T) *evalWriterHarness {
		h := newEvalWriterHarness(t).
			withExperiments(exposureFlagExperiment(), exposureMeasurementExperiment()).
			withCachedGroups(exposureHashGroup()).
			withAppend(nil, nil)
		h.withKeys()
		return h
	}

	t.Run("a user served normally gets the row of the flag experiment and the derived row", func(t *testing.T) {
		t.Parallel()
		h := newHarness(t)
		fails := h.write(eventDWHMap{
			"evt": exposureEvent(users.comparison, exposureDayStart+10, servedNormallyReason()),
		})
		assert.Empty(t, fails)
		rows := h.lastRows()
		require.Len(t, rows, 2)
		flagRow := rows["evt"]
		require.NotNil(t, flagRow)
		assert.Equal(t, exposureFlagID, flagRow.FeatureId)
		assert.Equal(t, exposureFlagVersion, flagRow.FeatureVersion)
		assert.Equal(t, "variation-served", flagRow.VariationId)
		derived := rows["evt:"+exposureGroupID]
		require.NotNil(t, derived)
		assert.Equal(t, exposureGroupID, derived.FeatureId)
		assert.Equal(t, holdoutdomain.NotInHoldoutVariationID, derived.VariationId)
	})

	t.Run("a held-out user gets only the derived row", func(t *testing.T) {
		t.Parallel()
		h := newHarness(t)
		fails := h.write(eventDWHMap{
			"evt": exposureEvent(users.heldOut, exposureDayStart+10, heldOutReason(exposureGroupID)),
		})
		assert.Empty(t, fails)
		rows := h.lastRows()
		require.Len(t, rows, 1)
		require.Contains(t, rows, "evt:"+exposureGroupID)
		assert.Equal(t, holdoutdomain.InHoldoutVariationID, rows["evt:"+exposureGroupID].VariationId)
	})

	t.Run("a user outside both arcs gets only the row of the flag experiment", func(t *testing.T) {
		t.Parallel()
		h := newHarness(t)
		fails := h.write(eventDWHMap{
			"evt": exposureEvent(users.outside, exposureDayStart+10, servedNormallyReason()),
		})
		assert.Empty(t, fails)
		rows := h.lastRows()
		require.Len(t, rows, 1)
		require.Contains(t, rows, "evt")
	})

	// The internal experiment has the ID of the group as its feature ID. An
	// evaluation event carrying that ID and version is not an event of the
	// group and writes nothing.
	t.Run("an event whose feature ID is the group ID is not written", func(t *testing.T) {
		t.Parallel()
		h := newHarness(t)
		evt := exposureEvent(users.comparison, exposureDayStart+10, servedNormallyReason())
		evt.FeatureId = exposureGroupID
		evt.FeatureVersion = holdoutdomain.InternalExperimentFeatureVersion
		fails := h.write(eventDWHMap{"evt": evt})
		assert.Empty(t, fails)
		assert.Empty(t, h.lastRows())
	})
}

func TestEvalEvtWriterHoldoutExposureFailures(t *testing.T) {
	t.Parallel()
	users := exposureUsers(t)
	key := exposureKey(exposureGroupStartAt, exposureDay, users.comparison)
	event := func() eventDWHMap {
		return eventDWHMap{
			"evt": exposureEvent(users.comparison, exposureDayStart+10, servedNormallyReason()),
		}
	}
	newHarness := func(t *testing.T) *evalWriterHarness {
		return newEvalWriterHarness(t).
			withExperiments(exposureFlagExperiment(), exposureMeasurementExperiment()).
			withCachedGroups(exposureHashGroup())
	}

	t.Run("a Redis error nacks the event and writes neither of its rows", func(t *testing.T) {
		t.Parallel()
		h := newHarness(t).withAppend(nil, nil)
		keys := h.withKeys()
		keys.setErr = errors.New("redis")
		events := event()
		events["other"] = exposureEvent(users.outside, exposureDayStart+10, servedNormallyReason())
		fails := h.write(events)
		assert.Equal(t, map[string]bool{"evt": true}, fails)
		// The event of the other user needs no key and is written.
		rows := h.lastRows()
		require.Len(t, rows, 1)
		require.Contains(t, rows, "other")

		// The redelivery finds no key and writes both rows.
		keys.setErr = nil
		fails = h.write(event())
		assert.Empty(t, fails)
		assert.Len(t, h.lastRows(), 2)
	})

	t.Run("a failed append releases the key and nacks the event once", func(t *testing.T) {
		t.Parallel()
		h := newHarness(t)
		keys := h.withKeys()
		gomock.InOrder(
			h.dwh.EXPECT().AppendRows(gomock.Any(), gomock.Any()).Return(nil, errors.New("dwh")),
			h.dwh.EXPECT().AppendRows(gomock.Any(), gomock.Any()).DoAndReturn(
				func(_ context.Context, rows []*epproto.EvaluationEvent) (map[string]bool, error) {
					h.appended = append(h.appended, rows)
					return nil, nil
				},
			),
		)
		fails := h.write(event())
		// Keyed by the ID of the event, not by the ID of the derived row.
		assert.Equal(t, map[string]bool{"evt": true}, fails)
		assert.Equal(t, []string{key}, keys.released)
		assert.Empty(t, keys.held())

		fails = h.write(event())
		assert.Empty(t, fails)
		rows := h.lastRows()
		require.Len(t, rows, 2)
		require.Contains(t, rows, "evt:"+exposureGroupID)
		assert.Equal(t, []string{key}, keys.held())
	})

	t.Run("a derived row that failed releases its key and nacks its event", func(t *testing.T) {
		t.Parallel()
		h := newHarness(t).withAppend(map[string]bool{"evt:" + exposureGroupID: true}, nil)
		keys := h.withKeys()
		fails := h.write(event())
		assert.Equal(t, map[string]bool{"evt": true}, fails)
		assert.Equal(t, []string{key}, keys.released)
	})

	t.Run("a derived row that cannot be retried releases its key and does not nack", func(t *testing.T) {
		t.Parallel()
		h := newHarness(t).withAppend(map[string]bool{"evt:" + exposureGroupID: false}, nil)
		keys := h.withKeys()
		fails := h.write(event())
		assert.Equal(t, map[string]bool{"evt": false}, fails)
		assert.Equal(t, []string{key}, keys.released)
	})

	t.Run("an event is retried when one of its two rows can be retried", func(t *testing.T) {
		t.Parallel()
		for _, appendFails := range []map[string]bool{
			{"evt": false, "evt:" + exposureGroupID: true},
			{"evt": true, "evt:" + exposureGroupID: false},
		} {
			h := newHarness(t).withAppend(appendFails, errors.New("dwh"))
			h.withKeys()
			fails := h.write(event())
			assert.Equal(t, map[string]bool{"evt": true}, fails)
		}
	})

	t.Run("a failed flag row keeps the key of the derived row that was written", func(t *testing.T) {
		t.Parallel()
		h := newHarness(t).withAppend(map[string]bool{"evt": true}, nil)
		keys := h.withKeys()
		fails := h.write(event())
		assert.Equal(t, map[string]bool{"evt": true}, fails)
		assert.Empty(t, keys.released)
		assert.Equal(t, []string{key}, keys.held())
	})

	t.Run("a key that cannot be released still nacks the event", func(t *testing.T) {
		t.Parallel()
		h := newHarness(t).withAppend(nil, errors.New("dwh"))
		keys := h.withKeys()
		keys.evalErr = errors.New("redis")
		fails := h.write(event())
		assert.Equal(t, map[string]bool{"evt": true}, fails)
	})
}

func TestEvalEvtWriterMeasuredGroupsSource(t *testing.T) {
	t.Parallel()
	users := exposureUsers(t)
	now := time.Unix(exposureDayStart+3600, 0)
	endedSince := now.Add(-experimentdomain.StoppedExperimentEventWindow).Unix()
	event := func() eventDWHMap {
		return eventDWHMap{
			"evt": exposureEvent(users.comparison, exposureDayStart+10, servedNormallyReason()),
		}
	}
	newHarness := func(t *testing.T) *evalWriterHarness {
		h := newEvalWriterHarness(t).withExperiments(exposureMeasurementExperiment()).withAppend(nil, nil)
		h.writer.holdoutExposures.now = func() time.Time { return now }
		return h
	}
	stored := func(groups ...*ftproto.HoldoutGroup) []*holdoutdomain.HoldoutGroup {
		list := make([]*holdoutdomain.HoldoutGroup, 0, len(groups))
		for _, group := range groups {
			list = append(list, &holdoutdomain.HoldoutGroup{HoldoutGroup: group})
		}
		return list
	}

	t.Run("a cached list with the group is used without a storage read", func(t *testing.T) {
		t.Parallel()
		h := newHarness(t).withCachedGroups(exposureHashGroup())
		h.withKeys()
		assert.Empty(t, h.write(event()))
		assert.Len(t, h.lastRows(), 1)
	})

	t.Run("a cache miss reads the storage and writes the cache", func(t *testing.T) {
		t.Parallel()
		h := newHarness(t)
		h.withKeys()
		group := exposureHashGroup()
		h.measuredGroupsCache.EXPECT().Get(exposureEnvID).Return(nil, cache.ErrNotFound)
		h.holdoutStorage.EXPECT().ListMeasuredGroups(gomock.Any(), exposureEnvID, endedSince).
			Return(stored(group), nil)
		h.measuredGroupsCache.EXPECT().PutIfNewer(
			&ftproto.HoldoutGroups{HoldoutGroups: []*ftproto.HoldoutGroup{group}},
			exposureEnvID,
			now.UnixNano(),
		).Return(true, nil)
		assert.Empty(t, h.write(event()))
		assert.Len(t, h.lastRows(), 1)
	})

	t.Run("a cached list without the group of a listed experiment reads the storage", func(t *testing.T) {
		t.Parallel()
		h := newHarness(t)
		h.withKeys()
		other := exposureHashGroup()
		other.Id = "hg-other"
		h.measuredGroupsCache.EXPECT().Get(exposureEnvID).
			Return(&ftproto.HoldoutGroups{HoldoutGroups: []*ftproto.HoldoutGroup{other}}, nil)
		h.holdoutStorage.EXPECT().ListMeasuredGroups(gomock.Any(), exposureEnvID, endedSince).
			Return(stored(other, exposureHashGroup()), nil)
		h.measuredGroupsCache.EXPECT().PutIfNewer(gomock.Any(), exposureEnvID, now.UnixNano()).Return(true, nil)
		assert.Empty(t, h.write(event()))
		assert.Len(t, h.lastRows(), 1)
	})

	t.Run("a refused or failed cache write does not stop the rows", func(t *testing.T) {
		t.Parallel()
		for _, result := range []struct {
			accepted bool
			err      error
		}{{false, nil}, {false, errors.New("redis")}} {
			h := newHarness(t)
			h.withKeys()
			h.measuredGroupsCache.EXPECT().Get(exposureEnvID).Return(nil, errors.New("redis"))
			h.holdoutStorage.EXPECT().ListMeasuredGroups(gomock.Any(), exposureEnvID, endedSince).
				Return(stored(exposureHashGroup()), nil)
			h.measuredGroupsCache.EXPECT().PutIfNewer(gomock.Any(), exposureEnvID, gomock.Any()).
				Return(result.accepted, result.err)
			assert.Empty(t, h.write(event()))
			assert.Len(t, h.lastRows(), 1)
		}
	})

	t.Run("a storage error nacks every event of the environment and writes no row", func(t *testing.T) {
		t.Parallel()
		// The flag has an experiment too: its row is not written either, so
		// the redelivery writes both rows of the event together.
		h := newEvalWriterHarness(t).
			withExperiments(exposureFlagExperiment(), exposureMeasurementExperiment()).
			withAppend(nil, nil)
		h.measuredGroupsCache.EXPECT().Get(exposureEnvID).Return(nil, cache.ErrNotFound)
		h.holdoutStorage.EXPECT().ListMeasuredGroups(gomock.Any(), exposureEnvID, gomock.Any()).
			Return(nil, errors.New("db"))
		events := event()
		events["other"] = exposureEvent(users.outside, exposureDayStart+10, servedNormallyReason())
		fails := h.write(events)
		assert.Equal(t, map[string]bool{"evt": true, "other": true}, fails)
		assert.Empty(t, h.lastRows())
	})

	t.Run("a group the storage does not list either derives nothing", func(t *testing.T) {
		t.Parallel()
		h := newHarness(t)
		h.measuredGroupsCache.EXPECT().Get(exposureEnvID).Return(nil, cache.ErrNotFound)
		h.holdoutStorage.EXPECT().ListMeasuredGroups(gomock.Any(), exposureEnvID, endedSince).Return(nil, nil)
		// The empty list is cached, so the next batch does not read the
		// storage for a list that has the same groups.
		h.measuredGroupsCache.EXPECT().PutIfNewer(
			gomock.Cond(func(groups *ftproto.HoldoutGroups) bool {
				return groups != nil && len(groups.HoldoutGroups) == 0
			}),
			exposureEnvID,
			now.UnixNano(),
		).Return(true, nil)
		assert.Empty(t, h.write(event()))
		assert.Empty(t, h.lastRows())
	})

	t.Run("a group in the list without a listed internal experiment derives nothing", func(t *testing.T) {
		t.Parallel()
		// The cached list holds a second group whose experiment is not in
		// the experiments list. It targets the flag, but it is not measured.
		other := exposureHashGroup()
		other.Id = "hg-other"
		measured := exposureHashGroup()
		measured.Features = nil
		h := newHarness(t).withCachedGroups(other, measured)
		assert.Empty(t, h.write(event()))
		assert.Empty(t, h.lastRows())
	})
}

func TestEvalEvtWriterCheckExperiment(t *testing.T) {
	t.Parallel()
	w := &evalEvtWriter{}
	experiments := []*exproto.Experiment{exposureMeasurementExperiment(), exposureFlagExperiment()}
	patterns := []struct {
		desc     string
		event    *eventproto.EvaluationEvent
		expected error
	}{
		{
			desc:  "linked",
			event: exposureEvent("user", exposureDayStart, servedNormallyReason()),
		},
		{
			desc:  "linked without a reason",
			event: exposureEvent("user", exposureDayStart, nil),
		},
		{
			desc: "another version of the flag",
			event: func() *eventproto.EvaluationEvent {
				e := exposureEvent("user", exposureDayStart, servedNormallyReason())
				e.FeatureVersion++
				return e
			}(),
			expected: ErrExperimentNotFound,
		},
		{
			desc:     "after the experiment ended",
			event:    exposureEvent("user", exposureGroupStopAt+1, servedNormallyReason()),
			expected: ErrEvaluationEventIssuedAfterExperimentEnded,
		},
		{
			desc:     "held out",
			event:    exposureEvent("user", exposureDayStart, heldOutReason(exposureGroupID)),
			expected: ErrEvaluationEventHeldOut,
		},
		{
			desc: "an internal experiment is not the experiment of a flag",
			event: func() *eventproto.EvaluationEvent {
				e := exposureEvent("user", exposureDayStart, servedNormallyReason())
				e.FeatureId = exposureGroupID
				e.FeatureVersion = holdoutdomain.InternalExperimentFeatureVersion
				return e
			}(),
			expected: ErrExperimentNotFound,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, p.expected, w.checkExperiment(p.event, experiments))
		})
	}
}
