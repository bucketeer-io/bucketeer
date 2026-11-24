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

package notification

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	environmentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	sendermock "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/sender/mock"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	senderproto "github.com/bucketeer-io/bucketeer/v2/proto/subscription/sender"
)

func TestCreateExperimentRunningNotification(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*testing.T, *experimentRunningWatcher)
		expectedErr error
	}{
		{
			desc: "no experiment",
			setup: func(t *testing.T, w *experimentRunningWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{{Id: "ns0", Name: "ns0"}},
						Cursor:       "",
					}, nil)
				w.experimentClient.(*experimentclientmock.MockClient).EXPECT().ListExperiments(
					gomock.Any(), gomock.Any()).Return(
					&experimentproto.ListExperimentsResponse{
						Experiments: []*experimentproto.Experiment{},
					}, nil)
			},
		},
		{
			desc: "experiments exist",
			setup: func(t *testing.T, w *experimentRunningWatcher) {
				w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{{Id: "ns0", Name: "ns0"}},
						Cursor:       "",
					}, nil)
				w.experimentClient.(*experimentclientmock.MockClient).EXPECT().ListExperiments(
					gomock.Any(), gomock.Any()).Return(
					&experimentproto.ListExperimentsResponse{
						Experiments: []*experimentproto.Experiment{{
							Id:   "eid",
							Name: "ename",
						}, {
							Id:   "eid1",
							Name: "ename1",
						}},
					}, nil)
				w.sender.(*sendermock.MockSender).
					EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			w := newExperimentRunningWatcherWithMock(t, mockController)
			if p.setup != nil {
				p.setup(t, w)
			}
			err := w.Run(context.Background())
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

// The notification must not include an experiment whose kind is HOLDOUT:
// it is internal to a holdout group, and no user runs it.
func TestRunListsOnlyExperimentsOnFeatureFlags(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	w := newExperimentRunningWatcherWithMock(t, mockController)
	w.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironmentsV2(
		gomock.Any(), gomock.Any(),
	).Return(&environmentproto.ListEnvironmentsV2Response{
		Environments: []*environmentproto.EnvironmentV2{{Id: "ns0", Name: "ns0"}},
	}, nil)
	experiments := []*experimentproto.Experiment{{Id: "eid", Name: "ename"}}
	w.experimentClient.(*experimentclientmock.MockClient).EXPECT().ListExperiments(
		gomock.Any(),
		&experimentproto.ListExperimentsRequest{
			PageSize:      listRequestSize,
			EnvironmentId: "ns0",
			Statuses: []experimentproto.Experiment_Status{
				experimentproto.Experiment_RUNNING,
			},
			Kinds: []experimentproto.Experiment_Kind{
				experimentproto.Experiment_KIND_EXPERIMENT,
			},
			Archived: wrapperspb.Bool(false),
		},
	).Return(&experimentproto.ListExperimentsResponse{Experiments: experiments}, nil)
	w.sender.(*sendermock.MockSender).EXPECT().Send(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, ne *senderproto.NotificationEvent) error {
			assert.Equal(t, experiments, ne.Notification.ExperimentRunningNotification.Experiments)
			return nil
		},
	)

	assert.NoError(t, w.Run(context.Background()))
}

func newExperimentRunningWatcherWithMock(t *testing.T, c *gomock.Controller) *experimentRunningWatcher {
	t.Helper()
	return &experimentRunningWatcher{
		environmentClient: environmentclientmock.NewMockClient(c),
		experimentClient:  experimentclientmock.NewMockClient(c),
		sender:            sendermock.NewMockSender(c),
		logger:            zap.NewNop(),
		opts: &jobs.Options{
			Timeout: 5 * time.Minute,
		},
	}
}
