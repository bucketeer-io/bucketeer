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
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	fcmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	ftmock "github.com/bucketeer-io/bucketeer/pkg/feature/storage/mock"
	pullermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/mock"
	btstorage "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigtable"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	esproto "github.com/bucketeer-io/bucketeer/proto/event/service"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

var defaultOptions = options{
	logger: zap.NewNop(),
}

func TestMarshaEvent(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	layout := "2006-01-02 15:04:05 -0700 MST"
	t1, err := time.Parse(layout, "2014-01-17 23:02:03 +0000 UTC")
	require.NoError(t, err)
	patterns := map[string]struct {
		setup              func(context.Context, *Persister)
		input              interface{}
		expected           string
		expectedErr        error
		expectedRepeatable bool
	}{
		"success: user event": {
			setup: nil,
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
		"success evaluation event": {
			setup: nil,
			input: &eventproto.EvaluationEvent{
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
			},
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
		"err goal batch event: internal error from bigtable": {
			setup: func(ctx context.Context, p *Persister) {
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluations(
					ctx,
					"uid",
					"ns",
					"tag",
				).Return(nil, btstorage.ErrInternal).Times(1)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: t1.Unix(),
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
			expectedErr:        btstorage.ErrInternal,
			expectedRepeatable: true,
		},
		"success goal batch event: getting evaluations from bigtable": {
			setup: func(ctx context.Context, p *Persister) {
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluations(
					ctx,
					"uid",
					"ns",
					"tag",
				).Return([]*featureproto.Evaluation{
					{
						FeatureId:      "fid-0",
						FeatureVersion: int32(0),
						VariationId:    "vid-0",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_CLIENT},
					},
					{
						FeatureId:      "fid-1",
						FeatureVersion: int32(1),
						VariationId:    "vid-1",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
					},
				}, nil).Times(1)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: t1.Unix(),
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
			expected: `{
				"environmentNamespace": "ns",
				"evaluations": ["fid-0:0:vid-0:CLIENT","fid-1:1:vid-1:TARGET"],
				"goalId": "gid",
				"metric.userId": "uid",
				"ns.user.data.atr":"av",
				"sourceId":"GOAL_BATCH",
				"tag": "tag",
				"timestamp": "2014-01-17T23:02:03Z",
				"userId":"uid",
				"value": "1.2"
			}`,
			expectedErr:        nil,
			expectedRepeatable: false,
		},
		"success goal batch event: getting evaluations from evaluate process with segment users": {
			setup: func(ctx context.Context, p *Persister) {
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluations(
					ctx,
					"uid",
					"ns",
					"tag",
				).Return(nil, btstorage.ErrKeyNotFound).Times(1)
				p.featureClient.(*fcmock.MockClient).EXPECT().EvaluateFeatures(
					ctx,
					&featureproto.EvaluateFeaturesRequest{
						User: &userproto.User{
							Id:   "uid",
							Data: map[string]string{"atr": "av"},
						},
						EnvironmentNamespace: "ns",
						Tag:                  "tag",
					},
				).Return(
					&featureproto.EvaluateFeaturesResponse{
						UserEvaluations: &featureproto.UserEvaluations{
							Id: "uid",
							Evaluations: []*featureproto.Evaluation{
								{
									FeatureId:      "fid",
									FeatureVersion: int32(1),
									VariationId:    "vid-1",
									Reason:         &featureproto.Reason{Type: featureproto.Reason_RULE},
								},
							},
						},
					}, nil,
				).Times(1)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: t1.Unix(),
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
			expected: `{
				"environmentNamespace": "ns",
				"evaluations": ["fid:1:vid-1:RULE"],
				"goalId": "gid",
				"metric.userId": "uid",
				"ns.user.data.atr":"av",
				"sourceId":"GOAL_BATCH",
				"tag": "tag",
				"timestamp": "2014-01-17T23:02:03Z",
				"userId":"uid",
				"value": "1.2"
			}`,
			expectedErr:        nil,
			expectedRepeatable: false,
		},
		"err goal batch event: internal error from feature api": {
			setup: func(ctx context.Context, p *Persister) {
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluations(
					ctx,
					"uid",
					"ns",
					"tag",
				).Return(nil, btstorage.ErrKeyNotFound).Times(1)
				p.featureClient.(*fcmock.MockClient).EXPECT().EvaluateFeatures(
					ctx,
					&featureproto.EvaluateFeaturesRequest{
						User: &userproto.User{
							Id:   "uid",
							Data: map[string]string{"atr": "av"},
						},
						EnvironmentNamespace: "ns",
						Tag:                  "tag",
					},
				).Return(
					nil, btstorage.ErrInternal,
				).Times(1)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_GOAL_BATCH,
				Timestamp: t1.Unix(),
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
			expectedErr:        btstorage.ErrInternal,
			expectedRepeatable: false,
		},
		"success goal event: no tag info": {
			setup: nil,
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_ANDROID,
				Timestamp: t1.Unix(),
				GoalId:    "gid",
				UserId:    "uid",
				User: &userproto.User{
					Id:   "uid",
					Data: map[string]string{"atr": "av"},
				},
				Value: float64(1.2),
				Evaluations: []*featureproto.Evaluation{
					{
						FeatureId:      "fid-0",
						FeatureVersion: int32(0),
						VariationId:    "vid-0",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_CLIENT},
					},
					{
						FeatureId:      "fid-1",
						FeatureVersion: int32(1),
						VariationId:    "vid-1",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
					},
				},
				Tag: "",
			},
			expected: `{
				"environmentNamespace": "ns",
				"evaluations": ["fid-0:0:vid-0:CLIENT","fid-1:1:vid-1:TARGET"],
				"goalId": "gid",
				"metric.userId": "uid",
				"ns.user.data.atr":"av",
				"sourceId":"ANDROID",
				"tag": "",
				"timestamp": "2014-01-17T23:02:03Z",
				"userId":"uid",
				"value": "1.2"
			}`,
			expectedErr:        nil,
			expectedRepeatable: false,
		},
		"err goal event: internal": {
			setup: func(ctx context.Context, p *Persister) {
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluations(
					ctx,
					"uid",
					"ns",
					"tag",
				).Return(nil, btstorage.ErrInternal).Times(1)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_ANDROID,
				Timestamp: t1.Unix(),
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
			expectedErr:        btstorage.ErrInternal,
			expectedRepeatable: true,
		},
		"success goal event: key not found not in bigtable": {
			setup: func(ctx context.Context, p *Persister) {
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluations(
					ctx,
					"uid",
					"ns",
					"tag",
				).Return(nil, btstorage.ErrKeyNotFound).Times(1)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_ANDROID,
				Timestamp: t1.Unix(),
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
			expected: `{
				"environmentNamespace": "ns",
				"evaluations": [],
				"goalId": "gid",
				"metric.userId": "uid",
				"ns.user.data.atr":"av",
				"sourceId":"ANDROID",
				"tag": "tag",
				"timestamp": "2014-01-17T23:02:03Z",
				"userId":"uid",
				"value": "1.2"
			}`,
			expectedErr:        nil,
			expectedRepeatable: false,
		},
		"success goal event: getting evaluations from bigtable": {
			setup: func(ctx context.Context, p *Persister) {
				p.userEvaluationStorage.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluations(
					ctx,
					"uid",
					"ns",
					"tag",
				).Return([]*featureproto.Evaluation{
					{
						FeatureId:      "fid-0",
						FeatureVersion: int32(0),
						VariationId:    "vid-0",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_CLIENT},
					},
					{
						FeatureId:      "fid-1",
						FeatureVersion: int32(1),
						VariationId:    "vid-1",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
					},
				}, nil).Times(1)
			},
			input: &eventproto.GoalEvent{
				SourceId:  eventproto.SourceId_ANDROID,
				Timestamp: t1.Unix(),
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
			expected: `{
				"environmentNamespace": "ns",
				"evaluations": ["fid-0:0:vid-0:CLIENT","fid-1:1:vid-1:TARGET"],
				"goalId": "gid",
				"metric.userId": "uid",
				"ns.user.data.atr":"av",
				"sourceId":"ANDROID",
				"tag": "tag",
				"timestamp": "2014-01-17T23:02:03Z",
				"userId":"uid",
				"value": "1.2"
			}`,
			expectedErr:        nil,
			expectedRepeatable: false,
		},
		"err: ErrUnexpectedMessageType": {
			input:              "",
			expected:           "",
			expectedErr:        ErrUnexpectedMessageType,
			expectedRepeatable: false,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			persister := newPersister(mockController)
			if p.setup != nil {
				p.setup(persister.ctx, persister)
			}
			actual, repeatable, err := persister.marshalEvent(p.input, "ns")
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

func newPersister(c *gomock.Controller) *Persister {
	ctx, cancel := context.WithCancel(context.Background())
	return &Persister{
		featureClient:         fcmock.NewMockClient(c),
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
