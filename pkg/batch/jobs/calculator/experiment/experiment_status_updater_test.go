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

package experiment

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	ecmock "github.com/bucketeer-io/bucketeer/pkg/experiment/client/mock"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
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
				u.experimentClient.(*ecmock.MockClient).EXPECT().StartExperiment(gomock.Any(), gomock.Any()).Return(
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
				u.experimentClient.(*ecmock.MockClient).EXPECT().FinishExperiment(gomock.Any(), gomock.Any()).Return(
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
				u.experimentClient.(*ecmock.MockClient).EXPECT().StartExperiment(gomock.Any(), gomock.Any()).Return(
					&experimentproto.StartExperimentResponse{}, nil)
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
				u.experimentClient.(*ecmock.MockClient).EXPECT().FinishExperiment(gomock.Any(), gomock.Any()).Return(
					&experimentproto.FinishExperimentResponse{}, nil)
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

func newMockExperimentStatusUpdater(t *testing.T, c *gomock.Controller) *experimentStatusUpdater {
	return &experimentStatusUpdater{
		experimentClient: ecmock.NewMockClient(c),
		opts: &jobs.Options{
			Timeout: 5 * time.Second,
		},
		logger: zap.NewNop().Named("test-experiment-status-updater"),
	}
}
