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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	aomock "github.com/bucketeer-io/bucketeer/pkg/autoops/client/mock"
	ecmock "github.com/bucketeer-io/bucketeer/pkg/experiment/client/mock"
	fcmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	featuredomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	ftmock "github.com/bucketeer-io/bucketeer/pkg/feature/storage/mock"
	pullermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/mock"
	btstorage "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigtable"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	aoproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	esproto "github.com/bucketeer-io/bucketeer/proto/event/service"
	exproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

var defaultOptions = options{
	logger: zap.NewNop(),
}

func TestMarshalUserEvent(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	layout := "2006-01-02 15:04:05 -0700 MST"
	t1, err := time.Parse(layout, "2014-01-17 23:02:03 +0000 UTC")
	require.NoError(t, err)
	patterns := []struct {
		desc               string
		input              interface{}
		expected           string
		expectedErr        error
		expectedRepeatable bool
	}{
		{
			desc: "success: user event",
			input: &esproto.UserEvent{
				UserId:   "uid",
				SourceId: eventproto.SourceId_ANDROID,
				Tag:      "tag",
				LastSeen: t1.Unix(),
			},
			expected: `{
				"environmentNamespace": "ns",
				"sourceId": "ANDROID",
				"tag": "tag",
				"timestamp": "2014-01-17T23:02:03Z",
				"userId":"uid"
			}`,
			expectedErr:        nil,
			expectedRepeatable: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			persister := newPersister(mockController)
			actual, repeatable, err := persister.marshalEvent(persister.ctx, p.input, "ns")
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expectedRepeatable, repeatable)
			buf := new(bytes.Buffer)
			err = json.Compact(buf, []byte(p.expected))
			assert.Equal(t, buf.String(), actual)
		})
	}
}

func TestMarshalEvaluationEvent(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	layout := "2006-01-02 15:04:05 -0700 MST"
	t1, err := time.Parse(layout, "2014-01-17 23:02:03 +0000 UTC")
	require.NoError(t, err)
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
		Timestamp:      t1.Unix(),
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
	patterns := []struct {
		desc               string
		setup              func(context.Context, *Persister)
		input              interface{}
		expected           string
		expectedErr        error
		expectedRepeatable bool
	}{
		{
			desc: "error: failed to upsert evaluation event",
			setup: func(ctx context.Context, p *Persister) {
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().UpsertUserEvaluation(
					ctx,
					evaluation,
					environmentNamespace,
					"tag",
				).Return(btstorage.ErrInternal)
			},
			input:              evaluationEvent,
			expected:           "",
			expectedErr:        btstorage.ErrInternal,
			expectedRepeatable: true,
		},
		{
			desc: "success: evaluation event",
			setup: func(ctx context.Context, p *Persister) {
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().UpsertUserEvaluation(
					ctx,
					evaluation,
					environmentNamespace,
					"tag",
				).Return(nil)
			},
			input: evaluationEvent,
			expected: `{
				"environmentNamespace":"ns",
				"featureId": "fid",
				"featureVersion": "1",
				"metric.userId": "uid",
				"ns.user.data.atr":"av",
				"reason":"CLIENT",
				"sourceId":"UNKNOWN",
				"tag":"tag",
				"timestamp":"2014-01-17T23:02:03Z",
				"userId":"uid",
				"variationId":"vid"
			}`,
			expectedErr:        nil,
			expectedRepeatable: false,
		},
		{
			desc:               "err: ErrUnexpectedMessageType",
			input:              "",
			expected:           "",
			expectedErr:        ErrUnexpectedMessageType,
			expectedRepeatable: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			persister := newPersister(mockController)
			if p.setup != nil {
				p.setup(persister.ctx, persister)
			}
			actual, repeatable, err := persister.marshalEvent(persister.ctx, p.input, environmentNamespace)
			assert.Equal(t, p.expectedRepeatable, repeatable)
			if err != nil {
				assert.Equal(t, actual, "")
				assert.Equal(t, p.expectedErr, err)
			} else {
				assert.Equal(t, p.expectedErr, err)
				buf := new(bytes.Buffer)
				err = json.Compact(buf, []byte(p.expected))
				require.NoError(t, err)
				assert.Equal(t, buf.String(), actual)
			}
		})
	}
}

func TestMarshalGoalEventWithExperiments(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	timeNow := time.Now()
	timeFormated := time.Unix(timeNow.Unix(), 0).Format(time.RFC3339)
	environmentNamespace := "ns"
	patterns := []struct {
		desc               string
		setup              func(context.Context, *Persister)
		input              interface{}
		expected           string
		expectedErr        error
		expectedRepeatable bool
	}{
		{
			desc: "err: list experiment internal",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
					},
				).Return(nil, errors.New("internal"))
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: time.Now().Unix(),
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
			expected:           "",
			expectedErr:        errors.New("internal"),
			expectedRepeatable: true,
		},
		{
			desc: "err: list experiment empty",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
					},
				).Return(&exproto.ListExperimentsResponse{}, nil)
				p.autoOpsClient.(*aomock.MockClient).EXPECT().ListAutoOpsRules(
					ctx,
					&aoproto.ListAutoOpsRulesRequest{
						EnvironmentNamespace: environmentNamespace,
						PageSize:             listRequestSize,
						Cursor:               "",
					},
				).Return(&aoproto.ListAutoOpsRulesResponse{}, nil)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: time.Now().Unix(),
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
			expected:           "",
			expectedErr:        ErrNothingToLink,
			expectedRepeatable: false,
		},
		{
			desc: "err: experiment not found",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
					},
				).Return(&exproto.ListExperimentsResponse{
					Experiments: []*exproto.Experiment{
						{
							Id:      "experiment-id",
							GoalIds: []string{"goal-id"},
						},
					},
				}, nil)
				p.autoOpsClient.(*aomock.MockClient).EXPECT().ListAutoOpsRules(
					ctx,
					&aoproto.ListAutoOpsRulesRequest{
						EnvironmentNamespace: environmentNamespace,
						PageSize:             listRequestSize,
						Cursor:               "",
					},
				).Return(&aoproto.ListAutoOpsRulesResponse{}, nil)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: time.Now().Unix(),
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
			expected:           "",
			expectedErr:        ErrNothingToLink,
			expectedRepeatable: false,
		},
		{
			desc: "err: get evaluation not found",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
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
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: time.Now().Unix(),
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
			expected:           "",
			expectedErr:        btstorage.ErrKeyNotFound,
			expectedRepeatable: true,
		},
		{
			desc: "err: get evaluation internal",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
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
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: time.Now().Unix(),
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
			expected:           "",
			expectedErr:        errors.New("internal"),
			expectedRepeatable: true,
		},
		{
			desc: "err: get evaluation internal using empty tag",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
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
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: time.Now().Unix(),
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
			expected:           "",
			expectedErr:        errors.New("internal"),
			expectedRepeatable: true,
		},
		{
			desc: "success",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
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
				p.autoOpsClient.(*aomock.MockClient).EXPECT().ListAutoOpsRules(
					ctx,
					&aoproto.ListAutoOpsRulesRequest{
						EnvironmentNamespace: environmentNamespace,
						PageSize:             listRequestSize,
						Cursor:               "",
					},
				).Return(&aoproto.ListAutoOpsRulesResponse{}, nil)
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
				SourceId:  eventproto.SourceId_ANDROID,
				Timestamp: timeNow.Unix(),
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
			expected: fmt.Sprintf(`{
				"environmentNamespace": "ns",
				"evaluations": ["fid:1:vid:TARGET"],
				"goalId": "gid",
				"metric.userId": "uid",
				"ns.user.data.atr":"av",
				"sourceId":"ANDROID",
				"tag": "tag",
				"timestamp": "%s",
				"userId":"uid",
				"value": "1.2"
			}`, timeFormated),
			expectedErr:        nil,
			expectedRepeatable: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			persister := newPersister(mockController)
			if p.setup != nil {
				p.setup(persister.ctx, persister)
			}
			actual, repeatable, err := persister.marshalEvent(persister.ctx, p.input, environmentNamespace)
			assert.Equal(t, p.expectedRepeatable, repeatable)
			if err != nil {
				assert.Equal(t, actual, "")
				assert.Equal(t, p.expectedErr, err)
			} else {
				assert.Equal(t, p.expectedErr, err)
				buf := new(bytes.Buffer)
				err = json.Compact(buf, []byte(p.expected))
				require.NoError(t, err)
				assert.Equal(t, buf.String(), actual)
			}
		})
	}
}

func TestMarshalGoalEventWithAutoOpsRules(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	timeNow := time.Now()
	timeFormated := time.Unix(timeNow.Unix(), 0).Format(time.RFC3339)
	convert := func(oerc *aoproto.OpsEventRateClause) []*aoproto.Clause {
		var clauses []*aoproto.Clause
		c, err := ptypes.MarshalAny(oerc)
		require.NoError(t, err)
		clauses = append(clauses, &aoproto.Clause{Clause: c})
		return clauses
	}
	environmentNamespace := "ns"
	patterns := []struct {
		desc               string
		setup              func(context.Context, *Persister)
		input              interface{}
		expected           string
		expectedErr        error
		expectedRepeatable bool
	}{
		{
			desc: "err: list auto ops rules internal",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
					},
				).Return(&exproto.ListExperimentsResponse{}, nil)
				p.autoOpsClient.(*aomock.MockClient).EXPECT().ListAutoOpsRules(
					ctx,
					&aoproto.ListAutoOpsRulesRequest{
						EnvironmentNamespace: environmentNamespace,
						PageSize:             listRequestSize,
						Cursor:               "",
					},
				).Return(nil, errors.New("internal"))
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: time.Now().Unix(),
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
			expected:           "",
			expectedErr:        errors.New("internal"),
			expectedRepeatable: true,
		},
		{
			desc: "err: list auto ops rules empty",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
					},
				).Return(&exproto.ListExperimentsResponse{}, nil)
				p.autoOpsClient.(*aomock.MockClient).EXPECT().ListAutoOpsRules(
					ctx,
					&aoproto.ListAutoOpsRulesRequest{
						EnvironmentNamespace: environmentNamespace,
						PageSize:             listRequestSize,
						Cursor:               "",
					},
				).Return(&aoproto.ListAutoOpsRulesResponse{}, nil)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: time.Now().Unix(),
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
			expected:           "",
			expectedErr:        ErrNothingToLink,
			expectedRepeatable: false,
		},
		{
			desc: "err: auto ops rules not found",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
					},
				).Return(&exproto.ListExperimentsResponse{}, nil)
				p.autoOpsClient.(*aomock.MockClient).EXPECT().ListAutoOpsRules(
					ctx,
					&aoproto.ListAutoOpsRulesRequest{
						EnvironmentNamespace: environmentNamespace,
						PageSize:             listRequestSize,
						Cursor:               "",
					},
				).Return(&aoproto.ListAutoOpsRulesResponse{
					AutoOpsRules: []*aoproto.AutoOpsRule{
						{
							FeatureId: "fid-1",
							Clauses: convert(&aoproto.OpsEventRateClause{
								GoalId: "goal-id-1",
							}),
						},
						{
							FeatureId: "fid-2",
							Clauses: convert(&aoproto.OpsEventRateClause{
								GoalId: "goal-id-2",
							}),
						},
					},
				}, nil)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: time.Now().Unix(),
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
			expected:           "",
			expectedErr:        ErrNothingToLink,
			expectedRepeatable: false,
		},
		{
			desc: "err: get features internal error",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
					},
				).Return(&exproto.ListExperimentsResponse{}, nil)
				p.autoOpsClient.(*aomock.MockClient).EXPECT().ListAutoOpsRules(
					ctx,
					&aoproto.ListAutoOpsRulesRequest{
						EnvironmentNamespace: environmentNamespace,
						PageSize:             listRequestSize,
						Cursor:               "",
					},
				).Return(&aoproto.ListAutoOpsRulesResponse{
					AutoOpsRules: []*aoproto.AutoOpsRule{
						{
							FeatureId: "fid",
							Clauses: convert(&aoproto.OpsEventRateClause{
								GoalId: "gid",
							}),
						},
					},
				}, nil)
				p.featureClient.(*fcmock.MockClient).EXPECT().GetFeatures(
					ctx,
					&featureproto.GetFeaturesRequest{
						EnvironmentNamespace: environmentNamespace,
						Ids:                  []string{"fid"},
					},
				).Return(nil, errors.New("internal"))
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: time.Now().Unix(),
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
			expected:           "",
			expectedErr:        errors.New("internal"),
			expectedRepeatable: true,
		},
		{
			desc: "err: get evaluation not found",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
					},
				).Return(&exproto.ListExperimentsResponse{}, nil)
				p.autoOpsClient.(*aomock.MockClient).EXPECT().ListAutoOpsRules(
					ctx,
					&aoproto.ListAutoOpsRulesRequest{
						EnvironmentNamespace: environmentNamespace,
						PageSize:             listRequestSize,
						Cursor:               "",
					},
				).Return(&aoproto.ListAutoOpsRulesResponse{
					AutoOpsRules: []*aoproto.AutoOpsRule{
						{
							FeatureId: "fid",
							Clauses: convert(&aoproto.OpsEventRateClause{
								GoalId: "gid",
							}),
						},
					},
				}, nil)
				p.featureClient.(*fcmock.MockClient).EXPECT().GetFeatures(
					ctx,
					&featureproto.GetFeaturesRequest{
						EnvironmentNamespace: environmentNamespace,
						Ids:                  []string{"fid"},
					},
				).Return(&featureproto.GetFeaturesResponse{
					Features: []*featureproto.Feature{
						{
							Id:      "fid",
							Version: int32(1),
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
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: time.Now().Unix(),
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
			expected:           "",
			expectedErr:        btstorage.ErrKeyNotFound,
			expectedRepeatable: true,
		},
		{
			desc: "err: get evaluation internal",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
					},
				).Return(&exproto.ListExperimentsResponse{}, nil)
				p.autoOpsClient.(*aomock.MockClient).EXPECT().ListAutoOpsRules(
					ctx,
					&aoproto.ListAutoOpsRulesRequest{
						EnvironmentNamespace: environmentNamespace,
						PageSize:             listRequestSize,
						Cursor:               "",
					},
				).Return(&aoproto.ListAutoOpsRulesResponse{
					AutoOpsRules: []*aoproto.AutoOpsRule{
						{
							FeatureId: "fid",
							Clauses: convert(&aoproto.OpsEventRateClause{
								GoalId: "gid",
							}),
						},
					},
				}, nil)
				p.featureClient.(*fcmock.MockClient).EXPECT().GetFeatures(
					ctx,
					&featureproto.GetFeaturesRequest{
						EnvironmentNamespace: environmentNamespace,
						Ids:                  []string{"fid"},
					},
				).Return(&featureproto.GetFeaturesResponse{
					Features: []*featureproto.Feature{
						{
							Id:      "fid",
							Version: int32(1),
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
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: time.Now().Unix(),
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
			expected:           "",
			expectedErr:        errors.New("internal"),
			expectedRepeatable: true,
		},
		{
			desc: "err: get evaluation internal using empty tag",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
					},
				).Return(&exproto.ListExperimentsResponse{}, nil)
				p.autoOpsClient.(*aomock.MockClient).EXPECT().ListAutoOpsRules(
					ctx,
					&aoproto.ListAutoOpsRulesRequest{
						EnvironmentNamespace: environmentNamespace,
						PageSize:             listRequestSize,
						Cursor:               "",
					},
				).Return(&aoproto.ListAutoOpsRulesResponse{
					AutoOpsRules: []*aoproto.AutoOpsRule{
						{
							FeatureId: "fid",
							Clauses: convert(&aoproto.OpsEventRateClause{
								GoalId: "gid",
							}),
						},
					},
				}, nil)
				p.featureClient.(*fcmock.MockClient).EXPECT().GetFeatures(
					ctx,
					&featureproto.GetFeaturesRequest{
						EnvironmentNamespace: environmentNamespace,
						Ids:                  []string{"fid"},
					},
				).Return(&featureproto.GetFeaturesResponse{
					Features: []*featureproto.Feature{
						{
							Id:      "fid",
							Version: int32(1),
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
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: time.Now().Unix(),
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
			expected:           "",
			expectedErr:        errors.New("internal"),
			expectedRepeatable: true,
		},
		{
			desc: "err: auto ops is already triggered",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
					},
				).Return(&exproto.ListExperimentsResponse{}, nil)
				p.autoOpsClient.(*aomock.MockClient).EXPECT().ListAutoOpsRules(
					ctx,
					&aoproto.ListAutoOpsRulesRequest{
						EnvironmentNamespace: environmentNamespace,
						PageSize:             listRequestSize,
						Cursor:               "",
					},
				).Return(&aoproto.ListAutoOpsRulesResponse{
					AutoOpsRules: []*aoproto.AutoOpsRule{
						{
							FeatureId:   "fid",
							TriggeredAt: int64(1),
							Clauses: convert(&aoproto.OpsEventRateClause{
								GoalId: "gid",
							}),
						},
					},
				}, nil)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_ANDROID,
				Timestamp: timeNow.Unix(),
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
			expected:           "",
			expectedErr:        ErrNothingToLink,
			expectedRepeatable: false,
		},
		{
			desc: "success",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
					},
				).Return(&exproto.ListExperimentsResponse{}, nil)
				p.autoOpsClient.(*aomock.MockClient).EXPECT().ListAutoOpsRules(
					ctx,
					&aoproto.ListAutoOpsRulesRequest{
						EnvironmentNamespace: environmentNamespace,
						PageSize:             listRequestSize,
						Cursor:               "",
					},
				).Return(&aoproto.ListAutoOpsRulesResponse{
					AutoOpsRules: []*aoproto.AutoOpsRule{
						{
							FeatureId: "fid",
							Clauses: convert(&aoproto.OpsEventRateClause{
								GoalId: "gid",
							}),
						},
						{
							FeatureId:   "fid-2",
							TriggeredAt: int64(1),
							Clauses: convert(&aoproto.OpsEventRateClause{
								GoalId: "gid",
							}),
						},
						{
							FeatureId: "fid-3",
							Clauses: convert(&aoproto.OpsEventRateClause{
								GoalId: "gid",
							}),
						},
					},
				}, nil)
				p.featureClient.(*fcmock.MockClient).EXPECT().GetFeatures(
					ctx,
					&featureproto.GetFeaturesRequest{
						EnvironmentNamespace: environmentNamespace,
						Ids:                  []string{"fid", "fid-3"},
					},
				).Return(&featureproto.GetFeaturesResponse{
					Features: []*featureproto.Feature{
						{
							Id:      "fid",
							Version: int32(1),
						},
						{
							Id:      "fid-3",
							Version: int32(1),
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
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluation(
					ctx,
					"uid",
					environmentNamespace,
					"tag",
					"fid-3",
					int32(1),
				).Return(&featureproto.Evaluation{
					FeatureId:      "fid-3",
					FeatureVersion: int32(1),
					VariationId:    "vid-3",
					Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
				}, nil)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_ANDROID,
				Timestamp: timeNow.Unix(),
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
			expected: fmt.Sprintf(`{
				"environmentNamespace": "ns",
				"evaluations": ["fid:1:vid:TARGET","fid-3:1:vid-3:TARGET"],
				"goalId": "gid",
				"metric.userId": "uid",
				"ns.user.data.atr":"av",
				"sourceId":"ANDROID",
				"tag": "tag",
				"timestamp": "%s",
				"userId":"uid",
				"value": "1.2"
			}`, timeFormated),
			expectedErr:        nil,
			expectedRepeatable: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			persister := newPersister(mockController)
			if p.setup != nil {
				p.setup(persister.ctx, persister)
			}
			actual, repeatable, err := persister.marshalEvent(persister.ctx, p.input, environmentNamespace)
			assert.Equal(t, p.expectedRepeatable, repeatable)
			if err != nil {
				assert.Equal(t, actual, "")
				assert.Equal(t, p.expectedErr, err)
			} else {
				assert.Equal(t, p.expectedErr, err)
				buf := new(bytes.Buffer)
				err = json.Compact(buf, []byte(p.expected))
				require.NoError(t, err)
				assert.Equal(t, buf.String(), actual)
			}
		})
	}
}

func TestMarshalGoalEventWithExperimentsAndAutoOpsRules(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	timeNow := time.Now()
	timeFormated := time.Unix(timeNow.Unix(), 0).Format(time.RFC3339)
	convert := func(oerc *aoproto.OpsEventRateClause) []*aoproto.Clause {
		var clauses []*aoproto.Clause
		c, err := ptypes.MarshalAny(oerc)
		require.NoError(t, err)
		clauses = append(clauses, &aoproto.Clause{Clause: c})
		return clauses
	}
	timeMoreThan24Hours := timeNow.AddDate(0, 0, -2)
	environmentNamespace := "ns"
	patterns := []struct {
		desc               string
		setup              func(context.Context, *Persister)
		input              interface{}
		expected           string
		expectedErr        error
		expectedRepeatable bool
	}{
		{
			desc:               "err: ErrUnexpectedMessageType",
			input:              "",
			expected:           "",
			expectedErr:        ErrUnexpectedMessageType,
			expectedRepeatable: false,
		},
		{
			desc:  "err: invalid goal event timestamp",
			setup: nil,
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: timeMoreThan24Hours.Unix(),
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
			expected:           "",
			expectedErr:        ErrInvalidGoalEventTimestamp,
			expectedRepeatable: false,
		},
		{
			desc: "success: using same feature flag id, version and goal id in the experiments and in the auto ops",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
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
				p.autoOpsClient.(*aomock.MockClient).EXPECT().ListAutoOpsRules(
					ctx,
					&aoproto.ListAutoOpsRulesRequest{
						EnvironmentNamespace: environmentNamespace,
						PageSize:             listRequestSize,
						Cursor:               "",
					},
				).Return(&aoproto.ListAutoOpsRulesResponse{
					AutoOpsRules: []*aoproto.AutoOpsRule{
						{
							FeatureId: "fid",
							Clauses: convert(&aoproto.OpsEventRateClause{
								GoalId: "gid",
							}),
						},
					},
				}, nil)
				p.featureClient.(*fcmock.MockClient).EXPECT().GetFeatures(
					ctx,
					&featureproto.GetFeaturesRequest{
						EnvironmentNamespace: environmentNamespace,
						Ids:                  []string{"fid"},
					},
				).Return(&featureproto.GetFeaturesResponse{
					Features: []*featureproto.Feature{
						{
							Id:      "fid",
							Version: int32(1),
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
				}, nil).MaxTimes(2)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_ANDROID,
				Timestamp: timeNow.Unix(),
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
			expected: fmt.Sprintf(`{
				"environmentNamespace": "ns",
				"evaluations": ["fid:1:vid:TARGET"],
				"goalId": "gid",
				"metric.userId": "uid",
				"ns.user.data.atr":"av",
				"sourceId":"ANDROID",
				"tag": "tag",
				"timestamp": "%s",
				"userId":"uid",
				"value": "1.2"
			}`, timeFormated),
			expectedErr:        nil,
			expectedRepeatable: false,
		},
		{
			desc: "success: using same goal id but different feature flags",
			setup: func(ctx context.Context, p *Persister) {
				p.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
					ctx,
					&exproto.ListExperimentsRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: environmentNamespace,
						Statuses: []exproto.Experiment_Status{
							exproto.Experiment_RUNNING,
							exproto.Experiment_FORCE_STOPPED,
							exproto.Experiment_STOPPED,
						},
						Archived: &wrappers.BoolValue{Value: false},
					},
				).Return(&exproto.ListExperimentsResponse{
					Experiments: []*exproto.Experiment{
						{
							Id:             "experiment-id",
							GoalIds:        []string{"gid"},
							FeatureId:      "fid-1",
							FeatureVersion: int32(1),
						},
					},
				}, nil)
				p.autoOpsClient.(*aomock.MockClient).EXPECT().ListAutoOpsRules(
					ctx,
					&aoproto.ListAutoOpsRulesRequest{
						EnvironmentNamespace: environmentNamespace,
						PageSize:             listRequestSize,
						Cursor:               "",
					},
				).Return(&aoproto.ListAutoOpsRulesResponse{
					AutoOpsRules: []*aoproto.AutoOpsRule{
						{
							FeatureId: "fid-2",
							Clauses: convert(&aoproto.OpsEventRateClause{
								GoalId: "gid",
							}),
						},
					},
				}, nil)
				p.featureClient.(*fcmock.MockClient).EXPECT().GetFeatures(
					ctx,
					&featureproto.GetFeaturesRequest{
						EnvironmentNamespace: environmentNamespace,
						Ids:                  []string{"fid-2"},
					},
				).Return(&featureproto.GetFeaturesResponse{
					Features: []*featureproto.Feature{
						{
							Id:      "fid-2",
							Version: int32(1),
						},
					},
				}, nil)
				// Query for experiments
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluation(
					ctx,
					"uid",
					environmentNamespace,
					"tag",
					"fid-1",
					int32(1),
				).Return(&featureproto.Evaluation{
					FeatureId:      "fid-1",
					FeatureVersion: int32(1),
					VariationId:    "vid-1",
					Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
				}, nil).MaxTimes(1)
				// Query for auto ops rules
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluation(
					ctx,
					"uid",
					environmentNamespace,
					"tag",
					"fid-2",
					int32(1),
				).Return(&featureproto.Evaluation{
					FeatureId:      "fid-2",
					FeatureVersion: int32(1),
					VariationId:    "vid-2",
					Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
				}, nil).MaxTimes(1)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_ANDROID,
				Timestamp: timeNow.Unix(),
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
			expected: fmt.Sprintf(`{
				"environmentNamespace": "ns",
				"evaluations": ["fid-2:1:vid-2:TARGET", "fid-1:1:vid-1:TARGET"],
				"goalId": "gid",
				"metric.userId": "uid",
				"ns.user.data.atr":"av",
				"sourceId":"ANDROID",
				"tag": "tag",
				"timestamp": "%s",
				"userId":"uid",
				"value": "1.2"
			}`, timeFormated),
			expectedErr:        nil,
			expectedRepeatable: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			persister := newPersister(mockController)
			if p.setup != nil {
				p.setup(persister.ctx, persister)
			}
			actual, repeatable, err := persister.marshalEvent(persister.ctx, p.input, environmentNamespace)
			assert.Equal(t, p.expectedRepeatable, repeatable)
			if err != nil {
				assert.Equal(t, actual, "")
				assert.Equal(t, p.expectedErr, err)
			} else {
				assert.Equal(t, p.expectedErr, err)
				buf := new(bytes.Buffer)
				err = json.Compact(buf, []byte(p.expected))
				require.NoError(t, err)
				assert.Equal(t, buf.String(), actual)
			}
		})
	}
}

func TestUpsertMAU(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(context.Context, *gomock.Controller) *Persister
		input       proto.Message
		expectedErr error
	}{
		{
			desc: "not executed: mysqlClient is nil",
			setup: func(ctx context.Context, ctrl *gomock.Controller) *Persister {
				return newPersister(ctrl)
			},
			input:       &esproto.UserEvent{},
			expectedErr: nil,
		},
		{
			desc: "not executed: message is not UserEvent",
			setup: func(ctx context.Context, ctrl *gomock.Controller) *Persister {
				return newPersisterWithMysqlClient(ctrl)
			},
			input:       &eventproto.EvaluationEvent{},
			expectedErr: nil,
		},
		{
			desc: "success upsert UserEvent",
			setup: func(ctx context.Context, ctrl *gomock.Controller) *Persister {
				p := newPersisterWithMysqlClient(ctrl)
				p.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
				return p
			},
			input:       &esproto.UserEvent{},
			expectedErr: nil,
		},
		{
			desc: "error upsert UserEvent",
			setup: func(ctx context.Context, ctrl *gomock.Controller) *Persister {
				p := newPersisterWithMysqlClient(ctrl)
				p.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("internal"))
				return p
			},
			input:       &esproto.UserEvent{},
			expectedErr: errors.New("internal"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			persister := p.setup(context.Background(), mockController)
			actualErr := persister.upsertMAU(context.Background(), p.input, "ns")
			assert.Equal(t, p.expectedErr, actualErr)
		})
	}
}

func TestConvToEvaluation(t *testing.T) {
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
		persister := newPersister(mockController)
		ev, tag := persister.convToEvaluation(context.Background(), p.input)
		assert.True(t, proto.Equal(p.expected, ev), p.desc)
		assert.Equal(t, p.expectedTag, tag, p.desc)
	}
}

func TestEvaluationCountkey(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	featureID := "feature_id"
	variationID := "variation_id"
	unix := time.Now().Unix()
	environmentNamespace := "en-1"
	now := time.Unix(unix, 0)
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jpLocation)
	patterns := []struct {
		desc                 string
		kind                 string
		featureID            string
		variationID          string
		environmentNamespace string
		timestamp            int64
		expected             string
	}{
		{
			desc:                 "userCount",
			kind:                 userCountKey,
			featureID:            featureID,
			variationID:          variationID,
			environmentNamespace: environmentNamespace,
			timestamp:            unix,
			expected:             fmt.Sprintf("%s:%s:%d:%s:%s", environmentNamespace, userCountKey, date.Unix(), featureID, variationID),
		},
		{
			desc:                 "eventCount",
			kind:                 eventCountKey,
			featureID:            featureID,
			variationID:          variationID,
			environmentNamespace: environmentNamespace,
			timestamp:            unix,
			expected:             fmt.Sprintf("%s:%s:%d:%s:%s", environmentNamespace, eventCountKey, date.Unix(), featureID, variationID),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			persister := newPersister(mockController)
			actual := persister.newEvaluationCountkey(p.kind, p.featureID, p.variationID, p.environmentNamespace, p.timestamp)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestGetVariationID(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		variationID string
		reason      featureproto.Reason_Type
		expected    string
	}{
		{
			desc:        "get given variation id if off variation",
			variationID: "vID1",
			reason:      featureproto.Reason_OFF_VARIATION,
			expected:    "vID1",
		},
		{
			desc:        "get given variation id if target",
			variationID: "vID1",
			reason:      featureproto.Reason_TARGET,
			expected:    "vID1",
		},
		{
			desc:        "get given variation id if rule",
			variationID: "vID1",
			reason:      featureproto.Reason_RULE,
			expected:    "vID1",
		},
		{
			desc:        "get given variation id if prerequisite",
			variationID: "vID1",
			reason:      featureproto.Reason_PREREQUISITE,
			expected:    "vID1",
		},
		{
			desc:        "get default variation id if client",
			variationID: "vID1",
			reason:      featureproto.Reason_CLIENT,
			expected:    defaultVariationID,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := getVariationID(p.reason, p.variationID)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func newPersister(c *gomock.Controller) *Persister {
	ctx, cancel := context.WithCancel(context.Background())
	return &Persister{
		experimentClient:      ecmock.NewMockClient(c),
		featureClient:         fcmock.NewMockClient(c),
		autoOpsClient:         aomock.NewMockClient(c),
		puller:                pullermock.NewMockRateLimitedPuller(c),
		datastore:             nil,
		userEvaluationStorage: ftmock.NewMockUserEvaluationsStorage(c),
		opts:                  &defaultOptions,
		logger:                defaultOptions.logger,
		ctx:                   ctx,
		cancel:                cancel,
		doneCh:                make(chan struct{}),
	}
}

func newPersisterWithMysqlClient(c *gomock.Controller) *Persister {
	ctx, cancel := context.WithCancel(context.Background())
	return &Persister{
		experimentClient:      ecmock.NewMockClient(c),
		featureClient:         fcmock.NewMockClient(c),
		autoOpsClient:         aomock.NewMockClient(c),
		puller:                pullermock.NewMockRateLimitedPuller(c),
		datastore:             nil,
		userEvaluationStorage: ftmock.NewMockUserEvaluationsStorage(c),
		opts:                  &defaultOptions,
		logger:                defaultOptions.logger,
		ctx:                   ctx,
		cancel:                cancel,
		doneCh:                make(chan struct{}),
		mysqlClient:           mysqlmock.NewMockClient(c),
	}
}
