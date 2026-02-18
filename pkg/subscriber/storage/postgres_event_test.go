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

package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	storagev2 "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/v2"
	mockv2 "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/v2/mock"
	epproto "github.com/bucketeer-io/bucketeer/v2/proto/eventpersisterdwh"
)

func TestPostgresEvalEventWriterAppendRows(t *testing.T) {
	t.Parallel()
	errStorage := errors.New("storage error")
	patterns := []struct {
		desc          string
		events        []*epproto.EvaluationEvent
		setupMock     func(m *mockv2.MockEvaluationEventStorageV2)
		expectedFails map[string]bool
		expectedErr   error
	}{
		{
			desc:   "empty events: no storage call, empty fails",
			events: []*epproto.EvaluationEvent{},
			setupMock: func(m *mockv2.MockEvaluationEventStorageV2) {
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
			setupMock: func(m *mockv2.MockEvaluationEventStorageV2) {
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
			setupMock: func(m *mockv2.MockEvaluationEventStorageV2) {
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
			setupMock: func(m *mockv2.MockEvaluationEventStorageV2) {
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
			mockStorage := mockv2.NewMockEvaluationEventStorageV2(mockCtrl)
			p.setupMock(mockStorage)
			writer := NewPostgresEvalEventWriter(mockStorage)
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

	expectedParams := []storagev2.EvaluationEventParams{
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

	mockStorage := mockv2.NewMockEvaluationEventStorageV2(mockCtrl)
	mockStorage.EXPECT().
		CreateEvaluationEvents(gomock.Any(), expectedParams).
		Return(nil)

	writer := NewPostgresEvalEventWriter(mockStorage)
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
		setupMock     func(m *mockv2.MockGoalEventStorageV2)
		expectedFails map[string]bool
		expectedErr   error
	}{
		{
			desc:   "empty events: no storage call, empty fails",
			events: []*epproto.GoalEvent{},
			setupMock: func(m *mockv2.MockGoalEventStorageV2) {
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
			setupMock: func(m *mockv2.MockGoalEventStorageV2) {
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
			setupMock: func(m *mockv2.MockGoalEventStorageV2) {
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
			setupMock: func(m *mockv2.MockGoalEventStorageV2) {
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
			mockStorage := mockv2.NewMockGoalEventStorageV2(mockCtrl)
			p.setupMock(mockStorage)
			writer := NewPostgresGoalEventWriter(mockStorage)
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

	expectedParams := []storagev2.GoalEventParams{
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

	mockStorage := mockv2.NewMockGoalEventStorageV2(mockCtrl)
	mockStorage.EXPECT().
		CreateGoalEvents(gomock.Any(), expectedParams).
		Return(nil)

	writer := NewPostgresGoalEventWriter(mockStorage)
	fails, err := writer.AppendRows(context.Background(), []*epproto.GoalEvent{evt})
	require.NoError(t, err)
	assert.Empty(t, fails)
}
