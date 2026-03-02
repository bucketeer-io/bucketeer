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

package v2

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
)

func TestNewPostgresGoalEventStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewPostgresGoalEventStorage(mock.NewMockTransaction(mockController))
	assert.IsType(t, &postgresGoalEventStorage{}, storage)
}

func TestPostgresGoalEventStorage_CreateGoalEvents(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := context.Background()

	patterns := []struct {
		desc        string
		events      []GoalEventParams
		setup       func(s *postgresGoalEventStorage)
		expectedErr error
	}{
		{
			desc: "success: single event with all fields",
			events: []GoalEventParams{
				{
					ID:             "goal-event-id-1",
					EnvironmentID:  "env-id-1",
					Timestamp:      1000000000, // microseconds
					GoalID:         "goal-id-1",
					Value:          100.5,
					UserID:         "user-id-1",
					UserData:       `{"key":"value"}`,
					Tag:            "tag-1",
					SourceID:       "source-id-1",
					FeatureID:      "feature-id-1",
					FeatureVersion: 1,
					VariationID:    "variation-id-1",
					Reason:         "TARGET",
				},
			},
			setup: func(s *postgresGoalEventStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockTransaction).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: single event with optional fields empty",
			events: []GoalEventParams{
				{
					ID:             "goal-event-id-2",
					EnvironmentID:  "env-id-2",
					Timestamp:      2000000000,
					GoalID:         "goal-id-2",
					Value:          0,
					UserID:         "user-id-2",
					UserData:       "",
					Tag:            "",
					SourceID:       "",
					FeatureID:      "",
					FeatureVersion: 0,
					VariationID:    "",
					Reason:         "",
				},
			},
			setup: func(s *postgresGoalEventStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockTransaction).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: multiple events",
			events: []GoalEventParams{
				{
					ID:             "goal-event-id-3",
					EnvironmentID:  "env-id-3",
					Timestamp:      3000000000,
					GoalID:         "goal-id-3",
					Value:          50.25,
					UserID:         "user-id-3",
					UserData:       `{"name":"test"}`,
					Tag:            "tag-3",
					SourceID:       "source-id-3",
					FeatureID:      "feature-id-3",
					FeatureVersion: 3,
					VariationID:    "variation-id-3",
					Reason:         "DEFAULT",
				},
				{
					ID:             "goal-event-id-4",
					EnvironmentID:  "env-id-4",
					Timestamp:      4000000000,
					GoalID:         "goal-id-4",
					Value:          75.75,
					UserID:         "user-id-4",
					UserData:       `{"name":"test2"}`,
					Tag:            "tag-4",
					SourceID:       "source-id-4",
					FeatureID:      "feature-id-4",
					FeatureVersion: 4,
					VariationID:    "variation-id-4",
					Reason:         "RULE",
				},
			},
			setup: func(s *postgresGoalEventStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockTransaction).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
		{
			desc:   "success: empty events list",
			events: []GoalEventParams{},
			setup: func(s *postgresGoalEventStorage) {
			},
			expectedErr: nil,
		},
		{
			desc: "error: missing required field ID",
			events: []GoalEventParams{
				{
					ID:            "",
					EnvironmentID: "env-id",
					Timestamp:     1000000000,
					GoalID:        "goal-id",
					UserID:        "user-id",
				},
			},
			setup: func(s *postgresGoalEventStorage) {
			},
			expectedErr: errors.New("missing required fields: id=, envId=env-id, goalId=goal-id, userId=user-id"),
		},
		{
			desc: "error: missing required field EnvironmentID",
			events: []GoalEventParams{
				{
					ID:            "goal-event-id",
					EnvironmentID: "",
					Timestamp:     1000000000,
					GoalID:        "goal-id",
					UserID:        "user-id",
				},
			},
			setup: func(s *postgresGoalEventStorage) {
			},
			expectedErr: errors.New("missing required fields: id=goal-event-id, envId=, goalId=goal-id, userId=user-id"),
		},
		{
			desc: "error: missing required field GoalID",
			events: []GoalEventParams{
				{
					ID:            "goal-event-id",
					EnvironmentID: "env-id",
					Timestamp:     1000000000,
					GoalID:        "",
					UserID:        "user-id",
				},
			},
			setup: func(s *postgresGoalEventStorage) {
			},
			expectedErr: errors.New("missing required fields: id=goal-event-id, envId=env-id, goalId=, userId=user-id"),
		},
		{
			desc: "error: missing required field UserID",
			events: []GoalEventParams{
				{
					ID:            "goal-event-id",
					EnvironmentID: "env-id",
					Timestamp:     1000000000,
					GoalID:        "goal-id",
					UserID:        "",
				},
			},
			setup: func(s *postgresGoalEventStorage) {
			},
			expectedErr: errors.New("missing required fields: id=goal-event-id, envId=env-id, goalId=goal-id, userId="),
		},
		{
			desc: "error: exec context fails",
			events: []GoalEventParams{
				{
					ID:            "goal-event-id-5",
					EnvironmentID: "env-id-5",
					Timestamp:     5000000000,
					GoalID:        "goal-id-5",
					Value:         200.0,
					UserID:        "user-id-5",
				},
			},
			setup: func(s *postgresGoalEventStorage) {
				s.qe.(*mock.MockTransaction).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("database error"))
			},
			expectedErr: errors.New("failed to execute batch insert: database error"),
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &postgresGoalEventStorage{
				qe: mock.NewMockTransaction(mockController),
			}
			if p.setup != nil {
				p.setup(s)
			}
			err := s.CreateGoalEvents(ctx, p.events)
			if p.expectedErr != nil {
				assert.EqualError(t, err, p.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPostgresGoalEventStorage_CreateGoalEvents_Batching(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := context.Background()

	events := make([]GoalEventParams, goalBatchSize+100)
	for i := 0; i < len(events); i++ {
		events[i] = GoalEventParams{
			ID:            "goal-event-id",
			EnvironmentID: "env-id",
			Timestamp:     int64(i * 1000000),
			GoalID:        "goal-id",
			Value:         float32(i),
			UserID:        "user-id",
		}
	}

	s := &postgresGoalEventStorage{
		qe: mock.NewMockTransaction(mockController),
	}

	result := mock.NewMockResult(mockController)
	s.qe.(*mock.MockTransaction).EXPECT().ExecContext(
		gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(result, nil).Times(2)

	err := s.CreateGoalEvents(ctx, events)
	assert.NoError(t, err)
}

func TestPostgresGoalEventStorage_CreateGoalEvents_BatchingWithError(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := context.Background()

	events := make([]GoalEventParams, goalBatchSize+100)
	for i := 0; i < len(events); i++ {
		events[i] = GoalEventParams{
			ID:            "goal-event-id",
			EnvironmentID: "env-id",
			Timestamp:     int64(i * 1000000),
			GoalID:        "goal-id",
			Value:         float32(i),
			UserID:        "user-id",
		}
	}

	s := &postgresGoalEventStorage{
		qe: mock.NewMockTransaction(mockController),
	}

	result := mock.NewMockResult(mockController)
	gomock.InOrder(
		s.qe.(*mock.MockTransaction).EXPECT().ExecContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(result, nil),
		s.qe.(*mock.MockTransaction).EXPECT().ExecContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(nil, errors.New("database error on second batch")),
	)

	err := s.CreateGoalEvents(ctx, events)
	assert.EqualError(t, err, "failed to execute batch insert: database error on second batch")
}
