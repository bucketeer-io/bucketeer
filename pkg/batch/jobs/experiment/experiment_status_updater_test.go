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

package experiment

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	envmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	ecmock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

func TestUpdateStatus(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc     string
		setup    func(t *testing.T, u *experimentStatusUpdater)
		input    *experimentproto.Experiment
		expected error
	}{
		{
			desc: "error: StartExperiment fails",
			setup: func(t *testing.T, u *experimentStatusUpdater) {
				u.experimentClient.(*ecmock.MockClient).EXPECT().UpdateExperiment(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("test"))
			},
			input: &experimentproto.Experiment{
				Id:      "eid",
				Status:  experimentproto.Experiment_WAITING,
				StartAt: time.Date(2019, 12, 25, 00, 00, 00, 0, time.UTC).Unix(),
			},
			expected: errors.New("test"),
		},
		{
			desc: "error: FinishExperiment fails",
			setup: func(t *testing.T, u *experimentStatusUpdater) {
				u.experimentClient.(*ecmock.MockClient).EXPECT().UpdateExperiment(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("test"))
			},
			input: &experimentproto.Experiment{
				Id:      "eid",
				Status:  experimentproto.Experiment_RUNNING,
				StartAt: time.Date(2019, 12, 25, 00, 00, 00, 0, time.UTC).Unix(),
			},
			expected: errors.New("test"),
		},
		{
			desc: "success: no update waiting",
			input: &experimentproto.Experiment{
				Id:      "eid",
				Status:  experimentproto.Experiment_WAITING,
				StartAt: time.Date(2100, 12, 25, 00, 00, 00, 0, time.UTC).Unix(),
			},
			expected: nil,
		},
		{
			desc: "success: update waiting to running",
			setup: func(t *testing.T, u *experimentStatusUpdater) {
				u.experimentClient.(*ecmock.MockClient).EXPECT().UpdateExperiment(gomock.Any(), gomock.Any()).Return(
					&experimentproto.UpdateExperimentResponse{}, nil)
			},
			input: &experimentproto.Experiment{
				Id:      "eid",
				Status:  experimentproto.Experiment_WAITING,
				StartAt: time.Date(2019, 12, 25, 00, 00, 00, 0, time.UTC).Unix(),
			},
			expected: nil,
		},
		{
			desc: "success: no update running",
			input: &experimentproto.Experiment{
				Id:     "eid",
				Status: experimentproto.Experiment_RUNNING,
				StopAt: time.Date(2100, 12, 25, 00, 00, 00, 0, time.UTC).Unix(),
			},
			expected: nil,
		},
		{
			desc: "success: update running to stopped",
			setup: func(t *testing.T, u *experimentStatusUpdater) {
				u.experimentClient.(*ecmock.MockClient).EXPECT().UpdateExperiment(gomock.Any(), gomock.Any()).Return(
					&experimentproto.UpdateExperimentResponse{}, nil)
			},
			input: &experimentproto.Experiment{
				Id:      "eid",
				Status:  experimentproto.Experiment_RUNNING,
				StartAt: time.Date(2019, 12, 25, 00, 00, 00, 0, time.UTC).Unix(),
				StopAt:  time.Date(2019, 12, 26, 00, 00, 00, 0, time.UTC).Unix(),
			},
			expected: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			updater := newMockExperimentStatusUpdater(t, mockController)
			if p.setup != nil {
				p.setup(t, updater)
			}
			err := updater.updateStatus(context.Background(), "ns", p.input)
			assert.Equal(t, p.expected, err)
		})
	}
}

// The job must not list an experiment whose kind is HOLDOUT. The holdout service
// changes the status of such an experiment, and UpdateExperiment rejects it.
func TestRunListsOnlyExperimentsOnFeatureFlags(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	updater := newMockExperimentStatusUpdater(t, mockController)
	updater.environmentClient.(*envmock.MockClient).EXPECT().ListEnvironmentsV2(
		gomock.Any(), gomock.Any(),
	).Return(&environmentproto.ListEnvironmentsV2Response{
		Environments: []*environmentproto.EnvironmentV2{{Id: "ns0"}},
	}, nil)
	updater.experimentClient.(*ecmock.MockClient).EXPECT().ListExperiments(
		gomock.Any(),
		&experimentproto.ListExperimentsRequest{
			PageSize:      listRequestSize,
			EnvironmentId: "ns0",
			Statuses: []experimentproto.Experiment_Status{
				experimentproto.Experiment_WAITING,
				experimentproto.Experiment_RUNNING,
			},
			Kinds: []experimentproto.Experiment_Kind{
				experimentproto.Experiment_KIND_EXPERIMENT,
			},
		},
	).Return(&experimentproto.ListExperimentsResponse{
		Experiments: []*experimentproto.Experiment{{
			Id:      "eid",
			Status:  experimentproto.Experiment_WAITING,
			StartAt: time.Date(2019, 12, 25, 00, 00, 00, 0, time.UTC).Unix(),
		}},
	}, nil)
	updater.experimentClient.(*ecmock.MockClient).EXPECT().UpdateExperiment(
		gomock.Any(),
		&experimentproto.UpdateExperimentRequest{
			EnvironmentId: "ns0",
			Id:            "eid",
			Status: &experimentproto.UpdateExperimentRequest_UpdatedStatus{
				Status: experimentproto.Experiment_RUNNING,
			},
		},
	).Return(&experimentproto.UpdateExperimentResponse{}, nil)

	assert.NoError(t, updater.Run(context.Background()))
}

func newMockExperimentStatusUpdater(t *testing.T, c *gomock.Controller) *experimentStatusUpdater {
	return &experimentStatusUpdater{
		environmentClient: envmock.NewMockClient(c),
		experimentClient:  ecmock.NewMockClient(c),
		opts: &jobs.Options{
			Timeout: 5 * time.Second,
		},
		logger: zap.NewNop().Named("test-experiment-status-updater"),
	}
}
