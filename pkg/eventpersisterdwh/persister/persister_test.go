// Copyright 2024 The Bucketeer Authors.
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

package persister

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	cachemock "github.com/bucketeer-io/bucketeer/pkg/cache/v3/mock"
	ecmock "github.com/bucketeer-io/bucketeer/pkg/experiment/client/mock"
	ftmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	featuredomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/proto/eventpersisterdwh"
	exproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

var (
	jpLocation = time.FixedZone("Asia/Tokyo", 9*60*60)

	feature = &featureproto.Feature{
		Id:      "fid",
		Version: int32(1),
		Tags:    []string{"tag"},
		Rules: []*featureproto.Rule{
			{
				Id: "ruleID",
				Clauses: []*featureproto.Clause{
					{
						Id:       "clauseID",
						Values:   []string{"segmentID"},
						Operator: featureproto.Clause_SEGMENT,
					},
				},
				Strategy: &featureproto.Strategy{
					FixedStrategy: &featureproto.FixedStrategy{
						Variation: "variationID_B",
					},
				},
			},
		},
		Variations: []*featureproto.Variation{
			{
				Id:    "variationID_A",
				Value: "true",
			},
			{
				Id:    "variationID_B",
				Value: "false",
			},
		},
		VariationType: featureproto.Feature_BOOLEAN,
		DefaultStrategy: &featureproto.Strategy{
			FixedStrategy: &featureproto.FixedStrategy{
				Variation: "variationID_A",
			},
		},
	}

	feature2 = &featureproto.Feature{
		Id:      "fid-2",
		Version: int32(1),
		Tags:    []string{"tag"},
		Rules: []*featureproto.Rule{
			{
				Id: "ruleID-2",
				Clauses: []*featureproto.Clause{
					{
						Id:       "clauseID-2",
						Values:   []string{"segmentID-2"},
						Operator: featureproto.Clause_SEGMENT,
					},
				},
				Strategy: &featureproto.Strategy{
					FixedStrategy: &featureproto.FixedStrategy{
						Variation: "variationID-2_A",
					},
				},
			},
		},
		Variations: []*featureproto.Variation{
			{
				Id:    "variationID-2_A",
				Value: "true",
			},
			{
				Id:    "variationID-2_B",
				Value: "false",
			},
		},
		VariationType: featureproto.Feature_BOOLEAN,
		DefaultStrategy: &featureproto.Strategy{
			FixedStrategy: &featureproto.FixedStrategy{
				Variation: "variationID-2_A",
			},
		},
	}

	featureForError = &featureproto.Feature{
		Id: "fid-2",
		Rules: []*featureproto.Rule{
			{
				Id: "ruleID",
				Clauses: []*featureproto.Clause{
					{
						Id:       "clauseID",
						Values:   []string{"segmentID"},
						Operator: featureproto.Clause_SEGMENT,
					},
				},
				Strategy: &featureproto.Strategy{
					FixedStrategy: &featureproto.FixedStrategy{
						Variation: "variationID_B",
					},
				},
			},
		},
	}

	defaultOptions = options{
		logger: zap.NewNop(),
	}
)

func TestListExperiments(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := time.Now()
	environmentNamespace := "ns"
	patterns := []struct {
		desc        string
		setup       func(context.Context, *evalEvtWriter)
		expected    []*exproto.Experiment
		expectedErr error
	}{
		{
			desc: "error: failed to list experiments",
			setup: func(ctx context.Context, p *evalEvtWriter) {
				p.cache.(*cachemock.MockExperimentsCache).EXPECT().Get(environmentNamespace).Return(
					nil,
					errors.New("internal"),
				)
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             0,
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_STOPPED,
						},
					},
				).Return(nil, errors.New("internal"))
			},
			expected:    nil,
			expectedErr: errors.New("internal"),
		},
		{
			desc: "success: stop_at is older than 3 days",
			setup: func(ctx context.Context, p *evalEvtWriter) {
				experiments := &exproto.Experiments{
					Experiments: []*exproto.Experiment{
						{
							Id:             "experiment-id",
							GoalIds:        []string{"goal-id"},
							FeatureId:      "feature-id",
							FeatureVersion: 1,
							Status:         exproto.Experiment_STOPPED,
							StopAt:         now.Unix() - 3*day,
						},
						{
							Id:             "experiment-id-2",
							GoalIds:        []string{"goal-id-2"},
							FeatureId:      "feature-id-2",
							FeatureVersion: 1,
							Status:         exproto.Experiment_STOPPED,
							StopAt:         now.Unix() - 1*day,
						},
					},
				}
				p.cache.(*cachemock.MockExperimentsCache).EXPECT().Get(environmentNamespace).Return(
					experiments,
					nil,
				)
			},
			expected: []*exproto.Experiment{
				{
					Id:             "experiment-id-2",
					GoalIds:        []string{"goal-id-2"},
					FeatureId:      "feature-id-2",
					FeatureVersion: 1,
					Status:         exproto.Experiment_STOPPED,
					StopAt:         now.Unix() - 1*day,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: using the cache",
			setup: func(ctx context.Context, p *evalEvtWriter) {
				experiments := &exproto.Experiments{
					Experiments: []*exproto.Experiment{
						{
							Id:             "experiment-id",
							GoalIds:        []string{"goal-id"},
							FeatureId:      "feature-id",
							FeatureVersion: 1,
						},
					},
				}
				p.cache.(*cachemock.MockExperimentsCache).EXPECT().Get(environmentNamespace).Return(
					experiments,
					nil,
				)
			},
			expected: []*exproto.Experiment{
				{
					Id:             "experiment-id",
					GoalIds:        []string{"goal-id"},
					FeatureId:      "feature-id",
					FeatureVersion: 1,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: using the API",
			setup: func(ctx context.Context, p *evalEvtWriter) {
				experiments := &exproto.Experiments{
					Experiments: []*exproto.Experiment{
						{
							Id:             "experiment-id",
							GoalIds:        []string{"goal-id"},
							FeatureId:      "feature-id",
							FeatureVersion: 1,
						},
					},
				}
				p.cache.(*cachemock.MockExperimentsCache).EXPECT().Get(environmentNamespace).Return(
					nil,
					errors.New("internal"),
				)
				p.cache.(*cachemock.MockExperimentsCache).EXPECT().Put(experiments, environmentNamespace).Return(
					nil,
				)
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             0,
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_STOPPED,
						},
					},
				).Return(&exproto.ListExperimentsResponse{
					Experiments: experiments.Experiments,
				}, nil)
			},
			expected: []*exproto.Experiment{
				{
					Id:             "experiment-id",
					GoalIds:        []string{"goal-id"},
					FeatureId:      "feature-id",
					FeatureVersion: 1,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			persister := newEvalEventWriter(mockController)
			p.setup(ctx, persister)
			actual, err := persister.listExperiments(ctx, environmentNamespace)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestConvToEvaluationEvent(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	t1 := time.Now()
	environmentNamespace := "ns"
	evaluation := &featureproto.Evaluation{
		Id: featuredomain.EvaluationID(
			"fid",
			1,
			"uid",
		),
		FeatureId:      "fid",
		FeatureVersion: 1,
		UserId:         "uid",
		VariationId:    "vid",
		Reason:         &featureproto.Reason{Type: featureproto.Reason_CLIENT},
	}
	evaluationEvent := &eventproto.EvaluationEvent{
		Tag:            "tag",
		Timestamp:      t1.Add(-time.Hour).Unix(),
		FeatureId:      "fid",
		FeatureVersion: int32(1),
		UserId:         "uid",
		VariationId:    "vid",
		Reason:         &featureproto.Reason{Type: featureproto.Reason_CLIENT},
		User: &userproto.User{
			Id:   "uid",
			Data: map[string]string{"atr": "av"},
		},
	}
	evaluationEventWithTagEmpty := &eventproto.EvaluationEvent{
		Tag:            "",
		Timestamp:      t1.Add(-time.Hour).Unix(),
		FeatureId:      "fid",
		FeatureVersion: int32(1),
		UserId:         "uid",
		VariationId:    "vid",
		Reason:         &featureproto.Reason{Type: featureproto.Reason_CLIENT},
		User: &userproto.User{
			Id:   "uid",
			Data: nil,
		},
	}
	userData, err := json.Marshal(evaluationEvent.User.Data)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventID := "event-id"
	patterns := []struct {
		desc               string
		input              *eventproto.EvaluationEvent
		inputExperiment    []*exproto.Experiment
		expected           *epproto.EvaluationEvent
		expectedErr        error
		expectedRepeatable bool
	}{
		{
			desc:  "error: experiment not found",
			input: evaluationEvent,
			inputExperiment: []*exproto.Experiment{
				{
					Id:             "experiment-id",
					GoalIds:        []string{"goal-id"},
					FeatureId:      "invalid",
					FeatureVersion: -1,
				},
			},
			expected:           nil,
			expectedErr:        ErrExperimentNotFound,
			expectedRepeatable: false,
		},
		{
			desc:  "error: evaluation event was issued after the experiment ended",
			input: evaluationEvent,
			inputExperiment: []*exproto.Experiment{
				{
					Id:             "experiment-id",
					GoalIds:        []string{"goal-id"},
					FeatureId:      evaluationEvent.FeatureId,
					FeatureVersion: evaluation.FeatureVersion,
					Status:         exproto.Experiment_STOPPED,
					StartAt:        t1.Add(-time.Hour * 3).Unix(),
					StopAt:         t1.Add(-time.Hour * 2).Unix(),
				},
			},
			expected:           nil,
			expectedErr:        ErrEvaluationEventIssuedAfterExperimentEnded,
			expectedRepeatable: false,
		},
		{
			desc:  "success: evaluation event with running status",
			input: evaluationEvent,
			inputExperiment: []*exproto.Experiment{
				{
					Id:             "experiment-id",
					GoalIds:        []string{"goal-id"},
					FeatureId:      evaluationEvent.FeatureId,
					FeatureVersion: evaluation.FeatureVersion,
					Status:         exproto.Experiment_RUNNING,
					StartAt:        t1.Add(-time.Hour).Unix(),
					StopAt:         t1.Add(time.Hour).Unix(),
				},
			},
			expected: &epproto.EvaluationEvent{
				Id:                   eventID,
				FeatureId:            evaluationEvent.FeatureId,
				FeatureVersion:       evaluationEvent.FeatureVersion,
				UserData:             string(userData),
				UserId:               evaluationEvent.UserId,
				VariationId:          evaluationEvent.VariationId,
				Reason:               evaluationEvent.Reason.Type.String(),
				Tag:                  evaluationEvent.Tag,
				SourceId:             evaluationEvent.SourceId.String(),
				EnvironmentNamespace: environmentNamespace,
				Timestamp:            time.Unix(evaluationEvent.Timestamp, 0).UnixMicro(),
			},
			expectedErr:        nil,
			expectedRepeatable: false,
		},
		{
			desc:  "success: evaluation event with stopped status",
			input: evaluationEvent,
			inputExperiment: []*exproto.Experiment{
				{
					Id:             "experiment-id",
					GoalIds:        []string{"goal-id"},
					FeatureId:      evaluationEvent.FeatureId,
					FeatureVersion: evaluation.FeatureVersion,
					Status:         exproto.Experiment_RUNNING,
					StartAt:        t1.Add(-time.Hour).Unix(),
					StopAt:         t1.Add(time.Hour).Unix(),
				},
			},
			expected: &epproto.EvaluationEvent{
				Id:                   eventID,
				FeatureId:            evaluationEvent.FeatureId,
				FeatureVersion:       evaluationEvent.FeatureVersion,
				UserData:             string(userData),
				UserId:               evaluationEvent.UserId,
				VariationId:          evaluationEvent.VariationId,
				Reason:               evaluationEvent.Reason.Type.String(),
				Tag:                  evaluationEvent.Tag,
				SourceId:             evaluationEvent.SourceId.String(),
				EnvironmentNamespace: environmentNamespace,
				Timestamp:            time.Unix(evaluationEvent.Timestamp, 0).UnixMicro(),
			},
			expectedErr:        nil,
			expectedRepeatable: false,
		},
		{
			desc:  "success: with empty tag and user data is nil",
			input: evaluationEventWithTagEmpty,
			inputExperiment: []*exproto.Experiment{
				{
					Id:             "experiment-id",
					GoalIds:        []string{"goal-id"},
					FeatureId:      evaluationEventWithTagEmpty.FeatureId,
					FeatureVersion: evaluation.FeatureVersion,
					Status:         exproto.Experiment_STOPPED,
					StartAt:        t1.Add(-time.Hour * 2).Unix(),
					StopAt:         t1.Add(-time.Hour).Unix(),
				},
			},
			expected: &epproto.EvaluationEvent{
				Id:                   eventID,
				FeatureId:            evaluationEventWithTagEmpty.FeatureId,
				FeatureVersion:       evaluationEventWithTagEmpty.FeatureVersion,
				UserData:             "{}",
				UserId:               evaluationEventWithTagEmpty.UserId,
				VariationId:          evaluationEventWithTagEmpty.VariationId,
				Reason:               evaluationEventWithTagEmpty.Reason.Type.String(),
				Tag:                  "none",
				SourceId:             evaluationEventWithTagEmpty.SourceId.String(),
				EnvironmentNamespace: environmentNamespace,
				Timestamp:            time.Unix(evaluationEventWithTagEmpty.Timestamp, 0).UnixMicro(),
			},
			expectedErr:        nil,
			expectedRepeatable: false,
		},
		{
			desc:  "success: using cache",
			input: evaluationEventWithTagEmpty,
			inputExperiment: []*exproto.Experiment{
				{
					Id:             "experiment-id",
					GoalIds:        []string{"goal-id"},
					FeatureId:      evaluationEventWithTagEmpty.FeatureId,
					FeatureVersion: evaluation.FeatureVersion,
					Status:         exproto.Experiment_STOPPED,
					StartAt:        t1.Add(-time.Hour * 2).Unix(),
					StopAt:         t1.Add(-time.Hour).Unix(),
				},
			},
			expected: &epproto.EvaluationEvent{
				Id:                   eventID,
				FeatureId:            evaluationEventWithTagEmpty.FeatureId,
				FeatureVersion:       evaluationEventWithTagEmpty.FeatureVersion,
				UserData:             "{}",
				UserId:               evaluationEventWithTagEmpty.UserId,
				VariationId:          evaluationEventWithTagEmpty.VariationId,
				Reason:               evaluationEventWithTagEmpty.Reason.Type.String(),
				Tag:                  "none",
				SourceId:             evaluationEventWithTagEmpty.SourceId.String(),
				EnvironmentNamespace: environmentNamespace,
				Timestamp:            time.Unix(evaluationEventWithTagEmpty.Timestamp, 0).UnixMicro(),
			},
			expectedErr:        nil,
			expectedRepeatable: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			persister := newEvalEventWriter(mockController)
			actual, repeatable, err := persister.convToEvaluationEvent(ctx, p.input, eventID, environmentNamespace, p.inputExperiment)
			assert.Equal(t, p.expectedRepeatable, repeatable)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestConvToGoalEventWithExperiments(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	now := time.Now()

	environmentNamespace := "ns"
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eventID := "event-id"
	user := &userproto.User{
		Id:   "uid",
		Data: map[string]string{"atr": "av"},
	}
	userData, err := json.Marshal(user.Data)
	require.NoError(t, err)
	patterns := []struct {
		desc               string
		setup              func(context.Context, *goalEvtWriter)
		input              *eventproto.GoalEvent
		inputExperiment    []*exproto.Experiment
		expected           []*epproto.GoalEvent
		expectedErr        error
		expectedRepeatable bool
	}{
		{
			desc: "err: experiment not found",
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_ANDROID,
				Timestamp: now.Unix(),
				GoalId:    "gid",
				UserId:    "uid",
				User: &userproto.User{
					Id:   "uid",
					Data: map[string]string{"atr": "av"},
				},
				Value:       float64(1.2),
				Evaluations: nil,
				Tag:         "tag",
			},
			inputExperiment: []*exproto.Experiment{
				{
					Id:      "experiment-id",
					GoalIds: []string{"goal-id"},
				},
			},
			expected:           nil,
			expectedErr:        ErrExperimentNotFound,
			expectedRepeatable: false,
		},
		{
			desc: "err: ErrFailedToEvaluateUser",
			setup: func(ctx context.Context, p *goalEvtWriter) {
				p.featureClient.(*ftmock.MockClient).EXPECT().EvaluateFeatures(
					ctx,
					gomock.Any(),
				).Return(nil, errors.New("internal error"))
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_ANDROID,
				Timestamp: now.Unix(),
				GoalId:    "gid",
				UserId:    "uid",
				User: &userproto.User{
					Id:   "uid",
					Data: map[string]string{"atr": "av"},
				},
				Value:       float64(1.2),
				Evaluations: nil,
				Tag:         "tag",
			},
			inputExperiment: []*exproto.Experiment{
				{
					Id:             "experiment-id",
					GoalIds:        []string{"gid"},
					FeatureId:      "fid",
					FeatureVersion: int32(1),
					StartAt:        time.Now().Add(-time.Hour).Unix(),
					StopAt:         time.Now().Add(time.Hour).Unix(),
				},
			},
			expected:           nil,
			expectedErr:        ErrFailedToEvaluateUser,
			expectedRepeatable: true,
		},
		{
			desc: "err: ErrEvaluationsAreEmpty",
			setup: func(ctx context.Context, p *goalEvtWriter) {
				p.featureClient.(*ftmock.MockClient).EXPECT().EvaluateFeatures(
					ctx,
					gomock.Any(),
				).Return(&featureproto.EvaluateFeaturesResponse{
					UserEvaluations: &featureproto.UserEvaluations{},
				}, nil)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_ANDROID,
				Timestamp: now.Unix(),
				GoalId:    "gid",
				UserId:    "uid",
				User: &userproto.User{
					Id:   "uid",
					Data: map[string]string{"atr": "av"},
				},
				Value:       float64(1.2),
				Evaluations: nil,
				Tag:         "tag",
			},
			inputExperiment: []*exproto.Experiment{
				{
					Id:             "experiment-id",
					GoalIds:        []string{"gid"},
					FeatureId:      "fid",
					FeatureVersion: int32(1),
					StartAt:        time.Now().Add(-time.Hour).Unix(),
					StopAt:         time.Now().Add(time.Hour).Unix(),
				},
			},
			expected:           nil,
			expectedErr:        ErrEvaluationsAreEmpty,
			expectedRepeatable: false,
		},
		{
			desc: "error: goal event was issued after the experiment ended",
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_ANDROID,
				Timestamp: now.Unix(),
				GoalId:    "gid",
				UserId:    "uid",
				User: &userproto.User{
					Id:   "uid",
					Data: nil,
				},
				Value:       float64(1.2),
				Evaluations: nil,
				Tag:         "",
			},
			inputExperiment: []*exproto.Experiment{
				{
					Id:             "experiment-id",
					GoalIds:        []string{"gid"},
					FeatureId:      "fid",
					FeatureVersion: int32(1),
					Status:         exproto.Experiment_STOPPED,
					StartAt:        time.Now().Add(-time.Hour * 2).Unix(),
					StopAt:         time.Now().Add(-time.Hour * 1).Unix(),
				},
			},
			expected:           nil,
			expectedErr:        ErrExperimentNotFound,
			expectedRepeatable: false,
		},
		{
			desc: "success",
			setup: func(ctx context.Context, p *goalEvtWriter) {
				p.featureClient.(*ftmock.MockClient).EXPECT().EvaluateFeatures(
					ctx,
					&featureproto.EvaluateFeaturesRequest{
						EnvironmentNamespace: environmentNamespace,
						FeatureId:            "fid",
						Tag:                  "tag",
						User:                 user,
					},
				).Return(&featureproto.EvaluateFeaturesResponse{
					UserEvaluations: &featureproto.UserEvaluations{
						Id: "",
						Evaluations: []*featureproto.Evaluation{
							{
								Id:             "eval-id",
								FeatureId:      "fid",
								FeatureVersion: int32(1),
								VariationId:    "variationID_B",
								Reason: &featureproto.Reason{
									Type: featureproto.Reason_RULE,
								},
								UserId: user.Id,
							},
						},
						CreatedAt: time.Now().Unix(),
					},
				}, nil)
				p.featureClient.(*ftmock.MockClient).EXPECT().EvaluateFeatures(
					ctx,
					&featureproto.EvaluateFeaturesRequest{
						EnvironmentNamespace: environmentNamespace,
						FeatureId:            "fid-2",
						Tag:                  "tag",
						User:                 user,
					},
				).Return(&featureproto.EvaluateFeaturesResponse{
					UserEvaluations: &featureproto.UserEvaluations{
						Id: "",
						Evaluations: []*featureproto.Evaluation{
							{
								Id:             "eval-id",
								FeatureId:      "fid-2",
								FeatureVersion: int32(1),
								VariationId:    "variationID-2_A",
								Reason: &featureproto.Reason{
									Type: featureproto.Reason_DEFAULT,
								},
								UserId: user.Id,
							},
						},
						CreatedAt: time.Now().Unix(),
					},
				}, nil)
			},
			input: &eventproto.GoalEvent{
				SourceId:    eventproto.SourceId_ANDROID,
				Timestamp:   now.Unix(),
				GoalId:      "gid",
				UserId:      "uid",
				User:        user,
				Value:       float64(1.2),
				Evaluations: nil,
				Tag:         "tag",
			},
			inputExperiment: []*exproto.Experiment{
				{
					Id:             "experiment-id",
					GoalIds:        []string{"gid"},
					FeatureId:      "fid",
					FeatureVersion: int32(1),
					StartAt:        time.Now().Add(-time.Hour).Unix(),
					StopAt:         time.Now().Add(time.Hour).Unix(),
				},
				{
					Id:             "experiment-id-2",
					GoalIds:        []string{"gid"},
					FeatureId:      "fid-2",
					FeatureVersion: int32(1),
					StartAt:        time.Now().Add(-time.Hour).Unix(),
					StopAt:         time.Now().Add(time.Hour).Unix(),
				},
				// This experiment won't be computed
				// because the startAt is higher than the goal event timestamp
				{
					Id:             "experiment-id-3",
					GoalIds:        []string{"gid"},
					FeatureId:      "fid-2",
					FeatureVersion: int32(1),
					StartAt:        time.Now().Add(time.Hour).Unix(),
					StopAt:         time.Now().Add(time.Hour).Unix(),
				},
			},
			expected: []*epproto.GoalEvent{
				{
					SourceId:             eventproto.SourceId_ANDROID.String(),
					Id:                   eventID,
					GoalId:               "gid",
					UserId:               "uid",
					Value:                1.2,
					Tag:                  "tag",
					FeatureId:            "fid",
					FeatureVersion:       int32(1),
					VariationId:          "variationID_B",
					Reason:               featureproto.Reason_RULE.String(),
					UserData:             string(userData),
					EnvironmentNamespace: environmentNamespace,
					Timestamp:            time.Unix(now.Unix(), 0).UnixMicro(),
				},
				{
					SourceId:             eventproto.SourceId_ANDROID.String(),
					Id:                   eventID,
					GoalId:               "gid",
					UserId:               "uid",
					Value:                1.2,
					Tag:                  "tag",
					FeatureId:            "fid-2",
					FeatureVersion:       int32(1),
					VariationId:          "variationID-2_A",
					Reason:               featureproto.Reason_DEFAULT.String(),
					UserData:             string(userData),
					EnvironmentNamespace: environmentNamespace,
					Timestamp:            time.Unix(now.Unix(), 0).UnixMicro(),
				},
			},
			expectedErr:        nil,
			expectedRepeatable: false,
		},
		{
			desc: "success: using cache",
			setup: func(ctx context.Context, p *goalEvtWriter) {
				p.featureClient.(*ftmock.MockClient).EXPECT().EvaluateFeatures(
					ctx,
					&featureproto.EvaluateFeaturesRequest{
						EnvironmentNamespace: environmentNamespace,
						FeatureId:            "fid",
						Tag:                  "",
						User: &userproto.User{
							Id:   "uid",
							Data: nil,
						},
					},
				).Return(&featureproto.EvaluateFeaturesResponse{
					UserEvaluations: &featureproto.UserEvaluations{
						Id: "",
						Evaluations: []*featureproto.Evaluation{
							{
								Id:             "eval-id",
								FeatureId:      "fid",
								FeatureVersion: int32(1),
								VariationId:    "variationID_B",
								Reason: &featureproto.Reason{
									Type: featureproto.Reason_RULE,
								},
								UserId: user.Id,
							},
						},
						CreatedAt: time.Now().Unix(),
					},
				}, nil)
				p.featureClient.(*ftmock.MockClient).EXPECT().EvaluateFeatures(
					ctx,
					&featureproto.EvaluateFeaturesRequest{
						EnvironmentNamespace: environmentNamespace,
						FeatureId:            "fid-2",
						Tag:                  "",
						User: &userproto.User{
							Id:   "uid",
							Data: nil,
						},
					},
				).Return(&featureproto.EvaluateFeaturesResponse{
					UserEvaluations: &featureproto.UserEvaluations{
						Id: "",
						Evaluations: []*featureproto.Evaluation{
							{
								Id:             "eval-id",
								FeatureId:      "fid-2",
								FeatureVersion: int32(1),
								VariationId:    "variationID-2_A",
								Reason: &featureproto.Reason{
									Type: featureproto.Reason_DEFAULT,
								},
								UserId: user.Id,
							},
						},
						CreatedAt: time.Now().Unix(),
					},
				}, nil)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_ANDROID,
				Timestamp: now.Unix(),
				GoalId:    "gid",
				UserId:    "uid",
				User: &userproto.User{
					Id:   "uid",
					Data: nil,
				},
				Value:       float64(1.2),
				Evaluations: nil,
				Tag:         "",
			},
			inputExperiment: []*exproto.Experiment{
				{
					Id:             "experiment-id",
					GoalIds:        []string{"gid"},
					FeatureId:      "fid",
					FeatureVersion: int32(1),
					StartAt:        time.Now().Add(-time.Hour).Unix(),
					StopAt:         time.Now().Add(time.Hour).Unix(),
				},
				{
					Id:             "experiment-id-2",
					GoalIds:        []string{"gid"},
					FeatureId:      "fid-2",
					FeatureVersion: int32(1),
					StartAt:        time.Now().Add(-time.Hour).Unix(),
					StopAt:         time.Now().Add(time.Hour).Unix(),
				},
				// This experiment won't be computed
				// because the startAt is higher than the goal event timestamp
				{
					Id:             "experiment-id-3",
					GoalIds:        []string{"gid"},
					FeatureId:      "fid-2",
					FeatureVersion: int32(1),
					StartAt:        time.Now().Add(time.Hour).Unix(),
					StopAt:         time.Now().Add(time.Hour).Unix(),
				},
				// This experiment won't be computed
				// because the goal event was issued after the experiment ended
				{
					Id:             "experiment-id-3",
					GoalIds:        []string{"gid"},
					FeatureId:      "fid-2",
					FeatureVersion: int32(1),
					StartAt:        time.Now().Add(-time.Hour * 3).Unix(),
					StopAt:         time.Now().Add(-time.Hour * 2).Unix(),
				},
			},
			expected: []*epproto.GoalEvent{
				{
					SourceId:             eventproto.SourceId_ANDROID.String(),
					Id:                   eventID,
					GoalId:               "gid",
					UserId:               "uid",
					Value:                1.2,
					Tag:                  "none",
					FeatureId:            "fid",
					FeatureVersion:       int32(1),
					VariationId:          "variationID_B",
					Reason:               featureproto.Reason_RULE.String(),
					UserData:             "{}",
					EnvironmentNamespace: environmentNamespace,
					Timestamp:            time.Unix(now.Unix(), 0).UnixMicro(),
				},
				{
					SourceId:             eventproto.SourceId_ANDROID.String(),
					Id:                   eventID,
					GoalId:               "gid",
					UserId:               "uid",
					Value:                1.2,
					Tag:                  "none",
					FeatureId:            "fid-2",
					FeatureVersion:       int32(1),
					VariationId:          "variationID-2_A",
					Reason:               featureproto.Reason_DEFAULT.String(),
					UserData:             "{}",
					EnvironmentNamespace: environmentNamespace,
					Timestamp:            time.Unix(now.Unix(), 0).UnixMicro(),
				},
			},
			expectedErr:        nil,
			expectedRepeatable: false,
		},
		{
			desc: "success: with empty tag and user data is nil",
			setup: func(ctx context.Context, p *goalEvtWriter) {
				p.featureClient.(*ftmock.MockClient).EXPECT().EvaluateFeatures(
					ctx,
					&featureproto.EvaluateFeaturesRequest{
						EnvironmentNamespace: environmentNamespace,
						FeatureId:            "fid",
						Tag:                  "",
						User: &userproto.User{
							Id:   "uid",
							Data: nil,
						},
					},
				).Return(&featureproto.EvaluateFeaturesResponse{
					UserEvaluations: &featureproto.UserEvaluations{
						Id: "",
						Evaluations: []*featureproto.Evaluation{
							{
								Id:             "eval-id",
								FeatureId:      "fid",
								FeatureVersion: int32(1),
								VariationId:    "variationID_B",
								Reason: &featureproto.Reason{
									Type: featureproto.Reason_RULE,
								},
								UserId: user.Id,
							},
						},
						CreatedAt: time.Now().Unix(),
					},
				}, nil)
				p.featureClient.(*ftmock.MockClient).EXPECT().EvaluateFeatures(
					ctx,
					&featureproto.EvaluateFeaturesRequest{
						EnvironmentNamespace: environmentNamespace,
						FeatureId:            "fid-2",
						Tag:                  "",
						User: &userproto.User{
							Id:   "uid",
							Data: nil,
						},
					},
				).Return(&featureproto.EvaluateFeaturesResponse{
					UserEvaluations: &featureproto.UserEvaluations{
						Id: "",
						Evaluations: []*featureproto.Evaluation{
							{
								Id:             "eval-id",
								FeatureId:      "fid-2",
								FeatureVersion: int32(1),
								VariationId:    "variationID-2_A",
								Reason: &featureproto.Reason{
									Type: featureproto.Reason_DEFAULT,
								},
								UserId: user.Id,
							},
						},
						CreatedAt: time.Now().Unix(),
					},
				}, nil)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_ANDROID,
				Timestamp: now.Unix(),
				GoalId:    "gid",
				UserId:    "uid",
				User: &userproto.User{
					Id:   "uid",
					Data: nil,
				},
				Value:       float64(1.2),
				Evaluations: nil,
				Tag:         "",
			},
			inputExperiment: []*exproto.Experiment{
				{
					Id:             "experiment-id",
					GoalIds:        []string{"gid"},
					FeatureId:      "fid",
					FeatureVersion: int32(1),
					StartAt:        time.Now().Add(-time.Hour).Unix(),
					StopAt:         time.Now().Add(time.Hour).Unix(),
				},
				{
					Id:             "experiment-id-2",
					GoalIds:        []string{"gid"},
					FeatureId:      "fid-2",
					FeatureVersion: int32(1),
					StartAt:        time.Now().Add(-time.Hour).Unix(),
					StopAt:         time.Now().Add(time.Hour).Unix(),
				},
				// This experiment won't be computed
				// because the startAt is higher than the goal event timestamp
				{
					Id:             "experiment-id-3",
					GoalIds:        []string{"gid"},
					FeatureId:      "fid-2",
					FeatureVersion: int32(1),
					StartAt:        time.Now().Add(time.Hour).Unix(),
					StopAt:         time.Now().Add(time.Hour).Unix(),
				},
				// This experiment won't be computed
				// because the goal event was issued after the experiment ended
				{
					Id:             "experiment-id-3",
					GoalIds:        []string{"gid"},
					FeatureId:      "fid-2",
					FeatureVersion: int32(1),
					StartAt:        time.Now().Add(-time.Hour * 3).Unix(),
					StopAt:         time.Now().Add(-time.Hour * 2).Unix(),
				},
			},
			expected: []*epproto.GoalEvent{
				{
					SourceId:             eventproto.SourceId_ANDROID.String(),
					Id:                   eventID,
					GoalId:               "gid",
					UserId:               "uid",
					Value:                1.2,
					Tag:                  "none",
					FeatureId:            "fid",
					FeatureVersion:       int32(1),
					VariationId:          "variationID_B",
					Reason:               featureproto.Reason_RULE.String(),
					UserData:             "{}",
					EnvironmentNamespace: environmentNamespace,
					Timestamp:            time.Unix(now.Unix(), 0).UnixMicro(),
				},
				{
					SourceId:             eventproto.SourceId_ANDROID.String(),
					Id:                   eventID,
					GoalId:               "gid",
					UserId:               "uid",
					Value:                1.2,
					Tag:                  "none",
					FeatureId:            "fid-2",
					FeatureVersion:       int32(1),
					VariationId:          "variationID-2_A",
					Reason:               featureproto.Reason_DEFAULT.String(),
					UserData:             "{}",
					EnvironmentNamespace: environmentNamespace,
					Timestamp:            time.Unix(now.Unix(), 0).UnixMicro(),
				},
			},
			expectedErr:        nil,
			expectedRepeatable: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			persister := newGoalEventWriter(mockController)
			if p.setup != nil {
				p.setup(ctx, persister)
			}
			actual, repeatable, err := persister.convToGoalEvents(
				ctx,
				p.input,
				eventID,
				environmentNamespace,
				p.inputExperiment,
			)
			assert.Equal(t, p.expectedRepeatable, repeatable)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newEvalEventWriter(c *gomock.Controller) *evalEvtWriter {
	return &evalEvtWriter{
		experimentClient: ecmock.NewMockClient(c),
		cache:            cachemock.NewMockExperimentsCache(c),
		location:         jpLocation,
		logger:           defaultOptions.logger,
	}
}

func newGoalEventWriter(c *gomock.Controller) *goalEvtWriter {
	return &goalEvtWriter{
		experimentClient: ecmock.NewMockClient(c),
		featureClient:    ftmock.NewMockClient(c),
		cache:            cachemock.NewMockExperimentsCache(c),
		location:         jpLocation,
		logger:           defaultOptions.logger,
	}
}
