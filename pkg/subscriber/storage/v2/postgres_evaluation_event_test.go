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

func TestNewPostgresEvaluationEventStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewPostgresEvaluationEventStorage(mock.NewMockTransaction(mockController))
	assert.IsType(t, &postgresEvaluationEventStorage{}, storage)
}

func TestPostgresEvaluationEventStorage_CreateEvaluationEvents(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := context.Background()

	patterns := []struct {
		desc        string
		events      []EvaluationEventParams
		setup       func(s *postgresEvaluationEventStorage)
		expectedErr error
	}{
		{
			desc: "success: single event with all fields",
			events: []EvaluationEventParams{
				{
					ID:             "event-id-1",
					EnvironmentID:  "env-id-1",
					Timestamp:      1000000000, // microseconds
					FeatureID:      "feature-id-1",
					FeatureVersion: 1,
					UserID:         "user-id-1",
					UserData:       `{"key":"value"}`,
					VariationID:    "variation-id-1",
					Reason:         "TARGET",
					Tag:            "tag-1",
					SourceID:       "source-id-1",
				},
			},
			setup: func(s *postgresEvaluationEventStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockTransaction).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: single event with optional fields empty",
			events: []EvaluationEventParams{
				{
					ID:             "event-id-2",
					EnvironmentID:  "env-id-2",
					Timestamp:      2000000000,
					FeatureID:      "feature-id-2",
					FeatureVersion: 2,
					UserID:         "user-id-2",
					UserData:       "",
					VariationID:    "",
					Reason:         "",
					Tag:            "",
					SourceID:       "",
				},
			},
			setup: func(s *postgresEvaluationEventStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockTransaction).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: multiple events",
			events: []EvaluationEventParams{
				{
					ID:             "event-id-3",
					EnvironmentID:  "env-id-3",
					Timestamp:      3000000000,
					FeatureID:      "feature-id-3",
					FeatureVersion: 3,
					UserID:         "user-id-3",
					UserData:       `{"name":"test"}`,
					VariationID:    "variation-id-3",
					Reason:         "DEFAULT",
					Tag:            "tag-3",
					SourceID:       "source-id-3",
				},
				{
					ID:             "event-id-4",
					EnvironmentID:  "env-id-4",
					Timestamp:      4000000000,
					FeatureID:      "feature-id-4",
					FeatureVersion: 4,
					UserID:         "user-id-4",
					UserData:       `{"name":"test2"}`,
					VariationID:    "variation-id-4",
					Reason:         "RULE",
					Tag:            "tag-4",
					SourceID:       "source-id-4",
				},
			},
			setup: func(s *postgresEvaluationEventStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockTransaction).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
		{
			desc:   "success: empty events list",
			events: []EvaluationEventParams{},
			setup: func(s *postgresEvaluationEventStorage) {
			},
			expectedErr: nil,
		},
		{
			desc: "error: missing required field ID",
			events: []EvaluationEventParams{
				{
					ID:             "",
					EnvironmentID:  "env-id",
					Timestamp:      1000000000,
					FeatureID:      "feature-id",
					FeatureVersion: 1,
					UserID:         "user-id",
				},
			},
			setup: func(s *postgresEvaluationEventStorage) {
			},
			expectedErr: errors.New("missing required fields: id=, envId=env-id, featureId=feature-id, userId=user-id"),
		},
		{
			desc: "error: missing required field EnvironmentID",
			events: []EvaluationEventParams{
				{
					ID:             "event-id",
					EnvironmentID:  "",
					Timestamp:      1000000000,
					FeatureID:      "feature-id",
					FeatureVersion: 1,
					UserID:         "user-id",
				},
			},
			setup: func(s *postgresEvaluationEventStorage) {
			},
			expectedErr: errors.New("missing required fields: id=event-id, envId=, featureId=feature-id, userId=user-id"),
		},
		{
			desc: "error: missing required field FeatureID",
			events: []EvaluationEventParams{
				{
					ID:             "event-id",
					EnvironmentID:  "env-id",
					Timestamp:      1000000000,
					FeatureID:      "",
					FeatureVersion: 1,
					UserID:         "user-id",
				},
			},
			setup: func(s *postgresEvaluationEventStorage) {
			},
			expectedErr: errors.New("missing required fields: id=event-id, envId=env-id, featureId=, userId=user-id"),
		},
		{
			desc: "error: missing required field UserID",
			events: []EvaluationEventParams{
				{
					ID:             "event-id",
					EnvironmentID:  "env-id",
					Timestamp:      1000000000,
					FeatureID:      "feature-id",
					FeatureVersion: 1,
					UserID:         "",
				},
			},
			setup: func(s *postgresEvaluationEventStorage) {
			},
			expectedErr: errors.New("missing required fields: id=event-id, envId=env-id, featureId=feature-id, userId="),
		},
		{
			desc: "error: exec context fails",
			events: []EvaluationEventParams{
				{
					ID:             "event-id-5",
					EnvironmentID:  "env-id-5",
					Timestamp:      5000000000,
					FeatureID:      "feature-id-5",
					FeatureVersion: 5,
					UserID:         "user-id-5",
				},
			},
			setup: func(s *postgresEvaluationEventStorage) {
				s.qe.(*mock.MockTransaction).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("database error"))
			},
			expectedErr: errors.New("failed to execute batch insert: database error"),
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &postgresEvaluationEventStorage{
				qe: mock.NewMockTransaction(mockController),
			}
			if p.setup != nil {
				p.setup(s)
			}
			err := s.CreateEvaluationEvents(ctx, p.events)
			if p.expectedErr != nil {
				assert.EqualError(t, err, p.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPostgresEvaluationEventStorage_CreateEvaluationEvents_Batching(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := context.Background()

	events := make([]EvaluationEventParams, evaluationBatchSize+100)
	for i := 0; i < len(events); i++ {
		events[i] = EvaluationEventParams{
			ID:             "event-id",
			EnvironmentID:  "env-id",
			Timestamp:      int64(i * 1000000),
			FeatureID:      "feature-id",
			FeatureVersion: int32(i),
			UserID:         "user-id",
		}
	}

	s := &postgresEvaluationEventStorage{
		qe: mock.NewMockTransaction(mockController),
	}

	result := mock.NewMockResult(mockController)
	s.qe.(*mock.MockTransaction).EXPECT().ExecContext(
		gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(result, nil).Times(2)

	err := s.CreateEvaluationEvents(ctx, events)
	assert.NoError(t, err)
}

func TestPostgresEvaluationEventStorage_CreateEvaluationEvents_BatchingWithError(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := context.Background()

	events := make([]EvaluationEventParams, evaluationBatchSize+100)
	for i := 0; i < len(events); i++ {
		events[i] = EvaluationEventParams{
			ID:             "event-id",
			EnvironmentID:  "env-id",
			Timestamp:      int64(i * 1000000),
			FeatureID:      "feature-id",
			FeatureVersion: int32(i),
			UserID:         "user-id",
		}
	}

	s := &postgresEvaluationEventStorage{
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

	err := s.CreateEvaluationEvents(ctx, events)
	assert.EqualError(t, err, "failed to execute batch insert: database error on second batch")
}
