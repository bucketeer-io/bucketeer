// Copyright 2022 The Bucketeer Authors.
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

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	ecmock "github.com/bucketeer-io/bucketeer/pkg/experiment/client/mock"
	featuredomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	ftmock "github.com/bucketeer-io/bucketeer/pkg/feature/storage/mock"
	btstorage "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigtable"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/proto/eventpersisterdwh"
	exproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

var defaultOptions = options{
	logger: zap.NewNop(),
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
		Timestamp:      t1.UnixMicro(),
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
	userData, err := json.Marshal(evaluationEvent.User.Data)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventID := "event-id"
	patterns := []struct {
		desc               string
		setup              func(context.Context, *evalEvtWriter)
		input              *eventproto.EvaluationEvent
		expected           *epproto.EvaluationEvent
		expectedErr        error
		expectedRepeatable bool
	}{
		{
			desc: "error: failed to list experiments",
			setup: func(ctx context.Context, p *evalEvtWriter) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
						},
					},
				).Return(nil, errors.New("internal"))
			},
			input:              evaluationEvent,
			expected:           nil,
			expectedErr:        errors.New("internal"),
			expectedRepeatable: true,
		},
		{
			desc: "error: experiment does not exist",
			setup: func(ctx context.Context, p *evalEvtWriter) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
						},
					},
				).Return(nil, ErrNoExperiments)
			},
			input:              evaluationEvent,
			expected:           nil,
			expectedErr:        ErrNoExperiments,
			expectedRepeatable: true,
		},
		{
			desc: "error: experiment not found",
			setup: func(ctx context.Context, p *evalEvtWriter) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
						},
					},
				).Return(&exproto.ListExperimentsResponse{
					Experiments: []*exproto.Experiment{
						{
							Id:             "experiment-id",
							GoalIds:        []string{"goal-id"},
							FeatureId:      "invalid",
							FeatureVersion: -1,
						},
					},
				}, nil)
			},
			input:              evaluationEvent,
			expected:           nil,
			expectedErr:        ErrExperimentNotFound,
			expectedRepeatable: false,
		},
		{
			desc: "error: failed to upsert user evaluation",
			setup: func(ctx context.Context, p *evalEvtWriter) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
						},
					},
				).Return(&exproto.ListExperimentsResponse{
					Experiments: []*exproto.Experiment{
						{
							Id:             "experiment-id",
							GoalIds:        []string{"goal-id"},
							FeatureId:      evaluationEvent.FeatureId,
							FeatureVersion: evaluation.FeatureVersion,
						},
					},
				}, nil)
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().UpsertUserEvaluation(
					ctx,
					evaluation,
					environmentNamespace,
					"tag",
				).Return(btstorage.ErrInternal)
			},
			input:              evaluationEvent,
			expected:           nil,
			expectedErr:        btstorage.ErrInternal,
			expectedRepeatable: true,
		},
		{
			desc: "success: evaluation event",
			setup: func(ctx context.Context, p *evalEvtWriter) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
						},
					},
				).Return(&exproto.ListExperimentsResponse{
					Experiments: []*exproto.Experiment{
						{
							Id:             "experiment-id",
							GoalIds:        []string{"goal-id"},
							FeatureId:      evaluationEvent.FeatureId,
							FeatureVersion: evaluation.FeatureVersion,
						},
					},
				}, nil)
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().UpsertUserEvaluation(
					ctx,
					evaluation,
					environmentNamespace,
					"tag",
				).Return(nil)
			},
			input: evaluationEvent,
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
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			persister := newEvalEventWriter(mockController)
			if p.setup != nil {
				p.setup(ctx, persister)
			}
			actual, repeatable, err := persister.convToEvaluationEvent(ctx, p.input, eventID, environmentNamespace)
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
	ctx, cancel := context.WithCancel(context.Background())
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
		expected           *epproto.GoalEvent
		expectedErr        error
		expectedRepeatable bool
	}{
		{
			desc: "err: list experiment internal",
			setup: func(ctx context.Context, p *goalEvtWriter) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
						},
					},
				).Return(nil, errors.New("internal"))
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
			expected:           nil,
			expectedErr:        errors.New("internal"),
			expectedRepeatable: true,
		},
		{
			desc: "err: list experiment empty",
			setup: func(ctx context.Context, p *goalEvtWriter) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
						},
					},
				).Return(&exproto.ListExperimentsResponse{}, nil)
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
			expected:           nil,
			expectedErr:        ErrNoExperiments,
			expectedRepeatable: false,
		},
		{
			desc: "err: experiment not found",
			setup: func(ctx context.Context, p *goalEvtWriter) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
						},
					},
				).Return(&exproto.ListExperimentsResponse{
					Experiments: []*exproto.Experiment{
						{
							Id:      "experiment-id",
							GoalIds: []string{"goal-id"},
						},
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
					Data: map[string]string{"atr": "av"},
				},
				Value:       float64(1.2),
				Evaluations: nil,
				Tag:         "tag",
			},
			expected:           nil,
			expectedErr:        ErrExperimentNotFound,
			expectedRepeatable: false,
		},
		{
			desc: "err: get evaluation not found",
			setup: func(ctx context.Context, p *goalEvtWriter) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
						},
					},
				).Return(&exproto.ListExperimentsResponse{
					Experiments: []*exproto.Experiment{
						{
							Id:             "experiment-id",
							GoalIds:        []string{"gid"},
							FeatureId:      "fid",
							FeatureVersion: int32(1),
						},
					},
				}, nil)
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluation(
					ctx,
					"uid",
					"ns",
					"tag",
					"fid",
					int32(1),
				).Return(nil, btstorage.ErrKeyNotFound)
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
			expected:           nil,
			expectedErr:        btstorage.ErrKeyNotFound,
			expectedRepeatable: true,
		},
		{
			desc: "err: get evaluation internal",
			setup: func(ctx context.Context, p *goalEvtWriter) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
						},
					},
				).Return(&exproto.ListExperimentsResponse{
					Experiments: []*exproto.Experiment{
						{
							Id:             "experiment-id",
							GoalIds:        []string{"gid"},
							FeatureId:      "fid",
							FeatureVersion: int32(1),
						},
					},
				}, nil)
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluation(
					ctx,
					"uid",
					environmentNamespace,
					"tag",
					"fid",
					int32(1),
				).Return(nil, errors.New("internal"))
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
			expected:           nil,
			expectedErr:        errors.New("internal"),
			expectedRepeatable: true,
		},
		{
			desc: "err: get evaluation internal using empty tag",
			setup: func(ctx context.Context, p *goalEvtWriter) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
						},
					},
				).Return(&exproto.ListExperimentsResponse{
					Experiments: []*exproto.Experiment{
						{
							Id:             "experiment-id",
							GoalIds:        []string{"gid"},
							FeatureId:      "fid",
							FeatureVersion: int32(1),
						},
					},
				}, nil)
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluation(
					ctx,
					"uid",
					environmentNamespace,
					"none",
					"fid",
					int32(1),
				).Return(nil, errors.New("internal"))
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
				Tag:         "",
			},
			expected:           nil,
			expectedErr:        errors.New("internal"),
			expectedRepeatable: true,
		},
		{
			desc: "success",
			setup: func(ctx context.Context, p *goalEvtWriter) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
						},
					},
				).Return(&exproto.ListExperimentsResponse{
					Experiments: []*exproto.Experiment{
						{
							Id:             "experiment-id",
							GoalIds:        []string{"gid"},
							FeatureId:      "fid",
							FeatureVersion: int32(1),
						},
					},
				}, nil)
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluation(
					ctx,
					"uid",
					environmentNamespace,
					"tag",
					"fid",
					int32(1),
				).Return(&featureproto.Evaluation{
					FeatureId:      "fid",
					FeatureVersion: int32(1),
					VariationId:    "vid",
					Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
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
			expected: &epproto.GoalEvent{
				SourceId:             eventproto.SourceId_ANDROID.String(),
				Id:                   eventID,
				GoalId:               "gid",
				UserId:               "uid",
				Value:                1.2,
				Tag:                  "tag",
				FeatureId:            "fid",
				FeatureVersion:       int32(1),
				VariationId:          "vid",
				Reason:               featureproto.Reason_TARGET.String(),
				UserData:             string(userData),
				EnvironmentNamespace: environmentNamespace,
				Timestamp:            time.Unix(now.Unix(), 0).UnixMicro(),
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
			actual, repeatable, err := persister.convToGoalEvent(ctx, p.input, eventID, environmentNamespace)
			assert.Equal(t, p.expectedRepeatable, repeatable)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestConvToEvaluationDwh(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	tag := "tag"
	evaluationEventWithTag := &eventproto.EvaluationEvent{
		FeatureId:      "feature-id",
		FeatureVersion: 2,
		UserId:         "user-id",
		VariationId:    "variation-id",
		User:           &userproto.User{Id: "user-id"},
		Reason: &featureproto.Reason{
			Type: featureproto.Reason_DEFAULT,
		},
		Tag:       tag,
		Timestamp: time.Now().Unix(),
	}
	evaluationEventWithoutTag := &eventproto.EvaluationEvent{
		FeatureId:      "feature-id",
		FeatureVersion: 2,
		UserId:         "user-id",
		VariationId:    "variation-id",
		User:           &userproto.User{Id: "user-id"},
		Reason: &featureproto.Reason{
			Type: featureproto.Reason_DEFAULT,
		},
		Timestamp: time.Now().Unix(),
	}
	patterns := []struct {
		desc        string
		input       *eventproto.EvaluationEvent
		expected    *featureproto.Evaluation
		expectedTag string
	}{
		{
			desc:  "success without tag",
			input: evaluationEventWithoutTag,
			expected: &featureproto.Evaluation{
				Id: featuredomain.EvaluationID(
					evaluationEventWithoutTag.FeatureId,
					evaluationEventWithoutTag.FeatureVersion,
					evaluationEventWithoutTag.UserId,
				),
				FeatureId:      evaluationEventWithoutTag.FeatureId,
				FeatureVersion: evaluationEventWithoutTag.FeatureVersion,
				UserId:         evaluationEventWithoutTag.UserId,
				VariationId:    evaluationEventWithoutTag.VariationId,
				Reason:         evaluationEventWithoutTag.Reason,
			},
			expectedTag: "none",
		},
		{
			desc:  "success with tag",
			input: evaluationEventWithTag,
			expected: &featureproto.Evaluation{
				Id: featuredomain.EvaluationID(
					evaluationEventWithTag.FeatureId,
					evaluationEventWithTag.FeatureVersion,
					evaluationEventWithTag.UserId,
				),
				FeatureId:      evaluationEventWithTag.FeatureId,
				FeatureVersion: evaluationEventWithTag.FeatureVersion,
				UserId:         evaluationEventWithTag.UserId,
				VariationId:    evaluationEventWithTag.VariationId,
				Reason:         evaluationEventWithTag.Reason,
			},
			expectedTag: tag,
		},
	}
	for _, p := range patterns {
		persister := newEvalEventWriter(mockController)
		ev, tag := persister.convToEvaluation(context.Background(), p.input)
		assert.True(t, proto.Equal(p.expected, ev), p.desc)
		assert.Equal(t, p.expectedTag, tag, p.desc)
	}
}

func newEvalEventWriter(c *gomock.Controller) *evalEvtWriter {
	return &evalEvtWriter{
		experimentClient:      ecmock.NewMockClient(c),
		userEvaluationStorage: ftmock.NewMockUserEvaluationsStorage(c),
		logger:                defaultOptions.logger,
	}
}

func newGoalEventWriter(c *gomock.Controller) *goalEvtWriter {
	return &goalEvtWriter{
		experimentClient:      ecmock.NewMockClient(c),
		userEvaluationStorage: ftmock.NewMockUserEvaluationsStorage(c),
		logger:                defaultOptions.logger,
	}
}
