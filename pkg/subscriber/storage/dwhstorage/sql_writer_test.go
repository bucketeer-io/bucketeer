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

package dwhstorage_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	dwhstorage "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/dwhstorage"
	mock "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/dwhstorage/mock"
	epproto "github.com/bucketeer-io/bucketeer/v2/proto/eventpersisterdwh"
)

func TestPostgresEvalEventWriterAppendRows(t *testing.T) {
	t.Parallel()
	errStorage := errors.New("storage error")
	patterns := []struct {
		desc          string
		events        []*epproto.EvaluationEvent
		setupMock     func(m *mock.MockEvaluationEventStorageV2)
		expectedFails map[string]bool
		expectedErr   error
	}{
		{
			desc:   "empty events: no storage call, empty fails",
			events: []*epproto.EvaluationEvent{},
			setupMock: func(m *mock.MockEvaluationEventStorageV2) {
				// no calls expected
			},
			expectedFails: map[string]bool{},
			expectedErr:   nil,
		},
		{
			desc: "single event success: storage called, no fails",
			events: []*epproto.EvaluationEvent{
				{
					Id:             "eval-1",
					EnvironmentId:  "env-1",
					Timestamp:      1000,
					FeatureId:      "feature-1",
					FeatureVersion: 2,
					UserId:         "user-1",
					UserData:       "",
					VariationId:    "var-1",
					Reason:         "rule",
					Tag:            "tag-1",
					SourceId:       "src-1",
				},
			},
			setupMock: func(m *mock.MockEvaluationEventStorageV2) {
				m.EXPECT().CreateEvaluationEvents(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedFails: map[string]bool{},
			expectedErr:   nil,
		},
		{
			desc: "multiple events success: storage called with all params, no fails",
			events: []*epproto.EvaluationEvent{
				{Id: "eval-1", EnvironmentId: "env-1", Timestamp: 100, FeatureId: "f-1", FeatureVersion: 1, UserId: "u-1", VariationId: "v-1", Reason: "default", Tag: "t-1", SourceId: "s-1"},
				{Id: "eval-2", EnvironmentId: "env-2", Timestamp: 200, FeatureId: "f-2", FeatureVersion: 2, UserId: "u-2", VariationId: "v-2", Reason: "rule", Tag: "t-2", SourceId: "s-2"},
			},
			setupMock: func(m *mock.MockEvaluationEventStorageV2) {
				m.EXPECT().CreateEvaluationEvents(gomock.Any(), gomock.Len(2)).Return(nil)
			},
			expectedFails: map[string]bool{},
			expectedErr:   nil,
		},
		{
			desc: "storage error: all events marked as failed",
			events: []*epproto.EvaluationEvent{
				{Id: "eval-1", EnvironmentId: "env-1"},
				{Id: "eval-2", EnvironmentId: "env-2"},
			},
			setupMock: func(m *mock.MockEvaluationEventStorageV2) {
				m.EXPECT().CreateEvaluationEvents(gomock.Any(), gomock.Any()).Return(errStorage)
			},
			expectedFails: map[string]bool{
				"eval-1": true,
				"eval-2": true,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockStorage := mock.NewMockEvaluationEventStorageV2(mockCtrl)
			p.setupMock(mockStorage)
			writer := dwhstorage.NewEvalEventWriter(mockStorage)
			fails, err := writer.AppendRows(context.Background(), p.events)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expectedFails, fails)
		})
	}
}

func TestPostgresEvalEventWriterAppendRows_VerifyParams(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evt := &epproto.EvaluationEvent{
		Id:             "eval-1",
		EnvironmentId:  "env-1",
		Timestamp:      1234567890,
		FeatureId:      "feature-a",
		FeatureVersion: 3,
		UserId:         "user-abc",
		UserData:       "",
		VariationId:    "var-x",
		Reason:         "default",
		Tag:            "ios",
		SourceId:       "sdk",
	}

	expectedParams := []dwhstorage.EvaluationEventParams{
		{
			ID:             "eval-1",
			EnvironmentID:  "env-1",
			Timestamp:      1234567890,
			FeatureID:      "feature-a",
			FeatureVersion: 3,
			UserID:         "user-abc",
			UserData:       `""`,
			VariationID:    "var-x",
			Reason:         "default",
			Tag:            "ios",
			SourceID:       "sdk",
		},
	}

	mockStorage := mock.NewMockEvaluationEventStorageV2(mockCtrl)
	mockStorage.EXPECT().
		CreateEvaluationEvents(gomock.Any(), expectedParams).
		Return(nil)

	writer := dwhstorage.NewEvalEventWriter(mockStorage)
	fails, err := writer.AppendRows(context.Background(), []*epproto.EvaluationEvent{evt})
	require.NoError(t, err)
	assert.Empty(t, fails)
}

func TestPostgresGoalEventWriterAppendRows(t *testing.T) {
	t.Parallel()
	errStorage := errors.New("storage error")
	patterns := []struct {
		desc          string
		events        []*epproto.GoalEvent
		setupMock     func(m *mock.MockGoalEventStorageV2)
		expectedFails map[string]bool
		expectedErr   error
	}{
		{
			desc:   "empty events: no storage call, empty fails",
			events: []*epproto.GoalEvent{},
			setupMock: func(m *mock.MockGoalEventStorageV2) {
				// no calls expected
			},
			expectedFails: map[string]bool{},
			expectedErr:   nil,
		},
		{
			desc: "single event success: storage called, no fails",
			events: []*epproto.GoalEvent{
				{
					Id:             "goal-1",
					EnvironmentId:  "env-1",
					Timestamp:      1000,
					GoalId:         "g-1",
					Value:          1.5,
					UserId:         "user-1",
					UserData:       "",
					Tag:            "tag-1",
					SourceId:       "src-1",
					FeatureId:      "feature-1",
					FeatureVersion: 2,
					VariationId:    "var-1",
					Reason:         "rule",
				},
			},
			setupMock: func(m *mock.MockGoalEventStorageV2) {
				m.EXPECT().CreateGoalEvents(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedFails: map[string]bool{},
			expectedErr:   nil,
		},
		{
			desc: "multiple events success: storage called with all params, no fails",
			events: []*epproto.GoalEvent{
				{Id: "goal-1", EnvironmentId: "env-1", GoalId: "g-1", Value: 1.0},
				{Id: "goal-2", EnvironmentId: "env-2", GoalId: "g-2", Value: 2.0},
				{Id: "goal-3", EnvironmentId: "env-3", GoalId: "g-3", Value: 3.0},
			},
			setupMock: func(m *mock.MockGoalEventStorageV2) {
				m.EXPECT().CreateGoalEvents(gomock.Any(), gomock.Len(3)).Return(nil)
			},
			expectedFails: map[string]bool{},
			expectedErr:   nil,
		},
		{
			desc: "storage error: all events marked as failed",
			events: []*epproto.GoalEvent{
				{Id: "goal-1", EnvironmentId: "env-1"},
				{Id: "goal-2", EnvironmentId: "env-2"},
				{Id: "goal-3", EnvironmentId: "env-3"},
			},
			setupMock: func(m *mock.MockGoalEventStorageV2) {
				m.EXPECT().CreateGoalEvents(gomock.Any(), gomock.Any()).Return(errStorage)
			},
			expectedFails: map[string]bool{
				"goal-1": true,
				"goal-2": true,
				"goal-3": true,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockStorage := mock.NewMockGoalEventStorageV2(mockCtrl)
			p.setupMock(mockStorage)
			writer := dwhstorage.NewGoalEventWriter(mockStorage)
			fails, err := writer.AppendRows(context.Background(), p.events)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expectedFails, fails)
		})
	}
}

func TestPostgresGoalEventWriterAppendRows_VerifyParams(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evt := &epproto.GoalEvent{
		Id:             "goal-1",
		EnvironmentId:  "env-1",
		Timestamp:      9876543210,
		GoalId:         "goal-abc",
		Value:          3.14,
		UserId:         "user-xyz",
		UserData:       "",
		Tag:            "android",
		SourceId:       "sdk",
		FeatureId:      "feature-b",
		FeatureVersion: 5,
		VariationId:    "var-y",
		Reason:         "rule",
	}

	expectedParams := []dwhstorage.GoalEventParams{
		{
			ID:             "goal-1",
			EnvironmentID:  "env-1",
			Timestamp:      9876543210,
			GoalID:         "goal-abc",
			Value:          3.14,
			UserID:         "user-xyz",
			UserData:       `""`,
			Tag:            "android",
			SourceID:       "sdk",
			FeatureID:      "feature-b",
			FeatureVersion: 5,
			VariationID:    "var-y",
			Reason:         "rule",
		},
	}

	mockStorage := mock.NewMockGoalEventStorageV2(mockCtrl)
	mockStorage.EXPECT().
		CreateGoalEvents(gomock.Any(), expectedParams).
		Return(nil)

	writer := dwhstorage.NewGoalEventWriter(mockStorage)
	fails, err := writer.AppendRows(context.Background(), []*epproto.GoalEvent{evt})
	require.NoError(t, err)
	assert.Empty(t, fails)
}

func TestGoalEventWriterAppendRowsRowIDs(t *testing.T) {
	t.Parallel()
	longFeatureID := strings.Repeat("f", 255)
	patterns := []struct {
		desc   string
		events []*epproto.GoalEvent
		check  func(t *testing.T, ids []string)
	}{
		{
			desc: "a single association keeps the event ID",
			events: []*epproto.GoalEvent{
				{Id: "goal-1", FeatureId: "feature-1", FeatureVersion: 1},
			},
			check: func(t *testing.T, ids []string) {
				assert.Equal(t, []string{"goal-1"}, ids)
			},
		},
		{
			desc: "several associations get distinct fixed-length IDs",
			events: []*epproto.GoalEvent{
				{Id: "goal-1", FeatureId: "feature-1", FeatureVersion: 1},
				{Id: "goal-1", FeatureId: longFeatureID, FeatureVersion: 3},
				{Id: "goal-2", FeatureId: "feature-1", FeatureVersion: 1},
			},
			check: func(t *testing.T, ids []string) {
				require.Len(t, ids, 3)
				assert.NotEqual(t, ids[0], ids[1])
				assert.Len(t, ids[0], 36)
				assert.Len(t, ids[1], 36)
				assert.Equal(t, "goal-2", ids[2])
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			var deliveries [][]string
			mockCtrl := gomock.NewController(t)
			mockStorage := mock.NewMockGoalEventStorageV2(mockCtrl)
			mockStorage.EXPECT().CreateGoalEvents(gomock.Any(), gomock.Any()).DoAndReturn(
				func(_ context.Context, params []dwhstorage.GoalEventParams) error {
					ids := make([]string, 0, len(params))
					for _, param := range params {
						ids = append(ids, param.ID)
					}
					deliveries = append(deliveries, ids)
					return nil
				}).Times(2)
			writer := dwhstorage.NewGoalEventWriter(mockStorage)
			// A redelivery produces the same IDs
			for range 2 {
				fails, err := writer.AppendRows(context.Background(), p.events)
				require.NoError(t, err)
				assert.Empty(t, fails)
			}
			p.check(t, deliveries[0])
			assert.Equal(t, deliveries[0], deliveries[1])
		})
	}
}
