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

package holdout

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	envmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	fcmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestUpdateStatus(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	past := time.Date(2019, 12, 25, 0, 0, 0, 0, time.UTC).Unix()
	future := time.Date(2100, 12, 25, 0, 0, 0, 0, time.UTC).Unix()

	patterns := []struct {
		desc     string
		setup    func(t *testing.T, u *statusUpdater)
		input    *featureproto.HoldoutGroup
		expected error
	}{
		{
			desc: "waiting group is not due",
			input: &featureproto.HoldoutGroup{
				Id:      "hg",
				Status:  featureproto.HoldoutGroup_WAITING,
				StartAt: future,
				StopAt:  future + 100,
			},
		},
		{
			desc: "waiting group starts",
			setup: func(t *testing.T, u *statusUpdater) {
				u.featureClient.(*fcmock.MockClient).EXPECT().
					StartHoldoutGroup(gomock.Any(), &featureproto.StartHoldoutGroupRequest{
						Id:            "hg",
						EnvironmentId: "env",
						Mode:          featureproto.StartHoldoutGroupRequest_SCHEDULED,
					}).
					Return(&featureproto.StartHoldoutGroupResponse{}, nil)
			},
			input: &featureproto.HoldoutGroup{
				Id:      "hg",
				Status:  featureproto.HoldoutGroup_WAITING,
				StartAt: past,
				StopAt:  future,
			},
		},
		{
			desc: "start failure is returned",
			setup: func(t *testing.T, u *statusUpdater) {
				u.featureClient.(*fcmock.MockClient).EXPECT().
					StartHoldoutGroup(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("start failed"))
			},
			input: &featureproto.HoldoutGroup{
				Id:      "hg",
				Status:  featureproto.HoldoutGroup_WAITING,
				StartAt: past,
				StopAt:  future,
			},
			expected: errors.New("start failed"),
		},
		{
			desc: "running group starts analysis",
			setup: func(t *testing.T, u *statusUpdater) {
				u.featureClient.(*fcmock.MockClient).EXPECT().
					StartHoldoutGroupAnalysis(gomock.Any(), &featureproto.StartHoldoutGroupAnalysisRequest{
						Id:            "hg",
						EnvironmentId: "env",
						Mode:          featureproto.StartHoldoutGroupAnalysisRequest_SCHEDULED,
					}).
					Return(&featureproto.StartHoldoutGroupAnalysisResponse{}, nil)
			},
			input: &featureproto.HoldoutGroup{
				Id:              "hg",
				Status:          featureproto.HoldoutGroup_RUNNING,
				StartAt:         past,
				StartAnalysisAt: past,
				StopAt:          future,
			},
		},
		{
			desc: "running group stops when stop_at has passed",
			setup: func(t *testing.T, u *statusUpdater) {
				u.featureClient.(*fcmock.MockClient).EXPECT().
					StopHoldoutGroup(gomock.Any(), &featureproto.StopHoldoutGroupRequest{
						Id:            "hg",
						EnvironmentId: "env",
						Mode:          featureproto.StopHoldoutGroupRequest_SCHEDULED,
					}).
					Return(&featureproto.StopHoldoutGroupResponse{}, nil)
			},
			input: &featureproto.HoldoutGroup{
				Id:              "hg",
				Status:          featureproto.HoldoutGroup_RUNNING,
				StartAt:         past,
				StartAnalysisAt: past,
				StopAt:          past,
			},
		},
		{
			desc: "analysis group stops when stop_at has passed",
			setup: func(t *testing.T, u *statusUpdater) {
				u.featureClient.(*fcmock.MockClient).EXPECT().
					StopHoldoutGroup(gomock.Any(), gomock.Any()).
					Return(&featureproto.StopHoldoutGroupResponse{}, nil)
			},
			input: &featureproto.HoldoutGroup{
				Id:      "hg",
				Status:  featureproto.HoldoutGroup_ANALYSIS,
				StartAt: past,
				StopAt:  past,
			},
		},
		{
			desc: "running group is not due",
			input: &featureproto.HoldoutGroup{
				Id:      "hg",
				Status:  featureproto.HoldoutGroup_RUNNING,
				StartAt: past,
				StopAt:  future,
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			updater := newMockStatusUpdater(t, mockController)
			if p.setup != nil {
				p.setup(t, updater)
			}
			err := updater.updateStatus(context.Background(), "env", p.input)
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestRun(t *testing.T) {
	t.Parallel()
	past := time.Date(2019, 12, 25, 0, 0, 0, 0, time.UTC).Unix()
	future := time.Date(2100, 12, 25, 0, 0, 0, 0, time.UTC).Unix()

	t.Run("environment list failure", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		updater := newStatusUpdater(t, controller)
		updater.environmentClient.(*envmock.MockClient).EXPECT().ListEnvironmentsV2(
			gomock.Any(), gomock.Any(),
		).Return(nil, errors.New("list failed"))
		assert.EqualError(t, updater.Run(context.Background()), "list failed")
	})

	t.Run("starts a due waiting group", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		updater := newStatusUpdater(t, controller)
		updater.environmentClient.(*envmock.MockClient).EXPECT().ListEnvironmentsV2(
			gomock.Any(), &environmentproto.ListEnvironmentsV2Request{
				PageSize: environmentListPageSize,
				Archived: &wrapperspb.BoolValue{Value: false},
			},
		).Return(&environmentproto.ListEnvironmentsV2Response{
			Environments: []*environmentproto.EnvironmentV2{{Id: "env"}},
		}, nil)
		updater.featureClient.(*fcmock.MockClient).EXPECT().ListHoldoutGroups(
			gomock.Any(), gomock.Any(),
		).Return(&featureproto.ListHoldoutGroupsResponse{
			HoldoutGroups: []*featureproto.HoldoutGroup{
				{Id: "later", Status: featureproto.HoldoutGroup_WAITING, StartAt: future, StopAt: future + 10},
				{Id: "due", Status: featureproto.HoldoutGroup_WAITING, StartAt: past, StopAt: future},
			},
		}, nil)
		updater.featureClient.(*fcmock.MockClient).EXPECT().StartHoldoutGroup(
			gomock.Any(), &featureproto.StartHoldoutGroupRequest{
				Id: "due", EnvironmentId: "env", Mode: featureproto.StartHoldoutGroupRequest_SCHEDULED,
			},
		).Return(&featureproto.StartHoldoutGroupResponse{}, nil)
		assert.NoError(t, updater.Run(context.Background()))
	})

	t.Run("skips the legacy environment with an empty ID", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		updater := newStatusUpdater(t, controller)
		updater.environmentClient.(*envmock.MockClient).EXPECT().ListEnvironmentsV2(
			gomock.Any(), gomock.Any(),
		).Return(&environmentproto.ListEnvironmentsV2Response{
			Environments: []*environmentproto.EnvironmentV2{{Id: ""}, {Id: "env"}},
		}, nil)
		// ListHoldoutGroups rejects an empty environment_id, so the legacy
		// environment must not be queried; the other environments still are.
		updater.featureClient.(*fcmock.MockClient).EXPECT().ListHoldoutGroups(
			gomock.Any(), gomock.Any(),
		).DoAndReturn(func(_ context.Context, req *featureproto.ListHoldoutGroupsRequest, _ ...grpc.CallOption) (*featureproto.ListHoldoutGroupsResponse, error) {
			assert.Equal(t, "env", req.EnvironmentId)
			return &featureproto.ListHoldoutGroupsResponse{}, nil
		}).Times(1)
		assert.NoError(t, updater.Run(context.Background()))
	})

	t.Run("group list failure is returned", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		updater := newStatusUpdater(t, controller)
		updater.environmentClient.(*envmock.MockClient).EXPECT().ListEnvironmentsV2(
			gomock.Any(), gomock.Any(),
		).Return(&environmentproto.ListEnvironmentsV2Response{
			Environments: []*environmentproto.EnvironmentV2{{Id: "env"}},
		}, nil)
		updater.featureClient.(*fcmock.MockClient).EXPECT().ListHoldoutGroups(
			gomock.Any(), gomock.Any(),
		).Return(nil, errors.New("groups failed"))
		assert.EqualError(t, updater.Run(context.Background()), "groups failed")
	})

	t.Run("stops an expired predecessor before starting the successor listed first", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		updater := newStatusUpdater(t, controller)
		updater.environmentClient.(*envmock.MockClient).EXPECT().ListEnvironmentsV2(
			gomock.Any(), gomock.Any(),
		).Return(&environmentproto.ListEnvironmentsV2Response{
			Environments: []*environmentproto.EnvironmentV2{{Id: "env"}},
		}, nil)
		updater.featureClient.(*fcmock.MockClient).EXPECT().ListHoldoutGroups(
			gomock.Any(), gomock.Any(),
		).Return(&featureproto.ListHoldoutGroupsResponse{
			HoldoutGroups: []*featureproto.HoldoutGroup{
				{Id: "successor", Status: featureproto.HoldoutGroup_WAITING, StartAt: past, StopAt: future},
				{Id: "predecessor", Status: featureproto.HoldoutGroup_RUNNING, StartAt: past - 100, StopAt: past},
				{Id: "missed", Status: featureproto.HoldoutGroup_WAITING, StartAt: past - 200, StopAt: past},
			},
		}, nil)
		stopPredecessor := updater.featureClient.(*fcmock.MockClient).EXPECT().StopHoldoutGroup(
			gomock.Any(), &featureproto.StopHoldoutGroupRequest{
				Id: "predecessor", EnvironmentId: "env", Mode: featureproto.StopHoldoutGroupRequest_SCHEDULED,
			},
		).Return(&featureproto.StopHoldoutGroupResponse{}, nil)
		stopMissed := updater.featureClient.(*fcmock.MockClient).EXPECT().StopHoldoutGroup(
			gomock.Any(), &featureproto.StopHoldoutGroupRequest{
				Id: "missed", EnvironmentId: "env", Mode: featureproto.StopHoldoutGroupRequest_SCHEDULED,
			},
		).Return(&featureproto.StopHoldoutGroupResponse{}, nil)
		startSuccessor := updater.featureClient.(*fcmock.MockClient).EXPECT().StartHoldoutGroup(
			gomock.Any(), &featureproto.StartHoldoutGroupRequest{
				Id: "successor", EnvironmentId: "env", Mode: featureproto.StartHoldoutGroupRequest_SCHEDULED,
			},
		).Return(&featureproto.StartHoldoutGroupResponse{}, nil)
		gomock.InOrder(stopPredecessor, stopMissed, startSuccessor)
		assert.NoError(t, updater.Run(context.Background()))
	})
}

func TestListPagination(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	updater := newStatusUpdater(t, controller)

	environments := make([]*environmentproto.EnvironmentV2, int(environmentListPageSize))
	for i := range environments {
		environments[i] = &environmentproto.EnvironmentV2{Id: "env"}
	}
	updater.environmentClient.(*envmock.MockClient).EXPECT().ListEnvironmentsV2(
		gomock.Any(), gomock.Any(),
	).Return(&environmentproto.ListEnvironmentsV2Response{Environments: environments, Cursor: "next"}, nil)
	updater.environmentClient.(*envmock.MockClient).EXPECT().ListEnvironmentsV2(
		gomock.Any(), gomock.Any(),
	).Return(&environmentproto.ListEnvironmentsV2Response{
		Environments: []*environmentproto.EnvironmentV2{{Id: "last"}},
	}, nil)
	listed, err := updater.listEnvironments(context.Background())
	require.NoError(t, err)
	assert.Len(t, listed, int(environmentListPageSize)+1)

	groups := make([]*featureproto.HoldoutGroup, int(holdoutListPageSize))
	for i := range groups {
		groups[i] = &featureproto.HoldoutGroup{Id: "g"}
	}
	updater.featureClient.(*fcmock.MockClient).EXPECT().ListHoldoutGroups(
		gomock.Any(), gomock.Any(),
	).Return(&featureproto.ListHoldoutGroupsResponse{HoldoutGroups: groups, Cursor: "next"}, nil)
	updater.featureClient.(*fcmock.MockClient).EXPECT().ListHoldoutGroups(
		gomock.Any(), gomock.Any(),
	).Return(&featureproto.ListHoldoutGroupsResponse{
		HoldoutGroups: []*featureproto.HoldoutGroup{{Id: "last"}},
	}, nil)
	listedGroups, err := updater.listGroups(context.Background(), "env")
	require.NoError(t, err)
	assert.Len(t, listedGroups, int(holdoutListPageSize)+1)
}

func newStatusUpdater(t *testing.T, controller *gomock.Controller) *statusUpdater {
	t.Helper()
	return &statusUpdater{
		environmentClient: envmock.NewMockClient(controller),
		featureClient:     fcmock.NewMockClient(controller),
		opts:              &jobs.Options{Timeout: 5 * time.Second},
		logger:            zap.NewNop(),
	}
}

func newMockStatusUpdater(t *testing.T, c *gomock.Controller) *statusUpdater {
	t.Helper()
	return &statusUpdater{
		featureClient: fcmock.NewMockClient(c),
		opts: &jobs.Options{
			Timeout: 5 * time.Second,
		},
		logger: zap.NewNop(),
	}
}
