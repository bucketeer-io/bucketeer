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

package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/storage/v2"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/storage/v2/mock"
	dbmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

func TestGetExperimentReturnsHoldoutKind(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	service := createExperimentService(mockController)
	service.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperiment(
		gomock.Any(), "id-1", "ns0",
	).Return(&domain.Experiment{
		Experiment: &experimentproto.Experiment{Id: "id-1", Kind: experimentproto.Experiment_HOLDOUT},
	}, nil)

	resp, err := service.GetExperiment(
		createContextWithEditor(),
		&experimentproto.GetExperimentRequest{Id: "id-1", EnvironmentId: "ns0"},
	)
	require.NoError(t, err)
	assert.Equal(t, experimentproto.Experiment_HOLDOUT, resp.Experiment.Kind)
}

func TestListExperimentsKinds(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc  string
		kinds []experimentproto.Experiment_Kind
	}{
		{
			desc:  "no kinds: no filter",
			kinds: nil,
		},
		{
			desc:  "experiments on feature flags only",
			kinds: []experimentproto.Experiment_Kind{experimentproto.Experiment_KIND_EXPERIMENT},
		},
		{
			desc:  "holdout experiments only",
			kinds: []experimentproto.Experiment_Kind{experimentproto.Experiment_HOLDOUT},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createExperimentService(mockController)
			// The list and the summary receive the same kinds.
			service.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().ListExperiments(
				gomock.Any(), gomock.Any(),
			).DoAndReturn(func(
				_ context.Context,
				params v2es.ListExperimentsParams,
			) ([]*experimentproto.Experiment, int, int64, error) {
				assert.Equal(t, "ns0", params.EnvironmentID)
				assert.Equal(t, p.kinds, params.Kinds)
				return []*experimentproto.Experiment{{Id: "id-1"}}, 1, int64(1), nil
			})
			service.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperimentSummary(
				gomock.Any(), "ns0", p.kinds,
			).Return(&v2es.ExperimentSummary{TotalRunningCount: 1}, nil)

			resp, err := service.ListExperiments(
				createContextWithEditor(),
				&experimentproto.ListExperimentsRequest{EnvironmentId: "ns0", Kinds: p.kinds},
			)
			require.NoError(t, err)
			assert.Len(t, resp.Experiments, 1)
			assert.Equal(t, int64(1), resp.Summary.TotalRunningCount)
		})
	}
}

func TestUpdateExperimentRejectsHoldoutKind(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	requests := []struct {
		desc string
		req  *experimentproto.UpdateExperimentRequest
	}{
		{
			desc: "rename",
			req: &experimentproto.UpdateExperimentRequest{
				Id:            "id-1",
				EnvironmentId: "ns0",
				Name:          wrapperspb.String("new-name"),
			},
		},
		{
			desc: "force stop",
			req: &experimentproto.UpdateExperimentRequest{
				Id:            "id-1",
				EnvironmentId: "ns0",
				Status: &experimentproto.UpdateExperimentRequest_UpdatedStatus{
					Status: experimentproto.Experiment_FORCE_STOPPED,
				},
			},
		},
		{
			desc: "archive",
			req: &experimentproto.UpdateExperimentRequest{
				Id:            "id-1",
				EnvironmentId: "ns0",
				Archived:      wrapperspb.Bool(true),
			},
		},
	}
	for _, p := range requests {
		t.Run(p.desc, func(t *testing.T) {
			service := createExperimentService(mockController)
			service.dbClient.(*dbmock.MockClient).EXPECT().RunInTransactionV2(
				gomock.Any(), gomock.Any(),
			).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			})
			service.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperiment(
				gomock.Any(), "id-1", "ns0",
			).Return(&domain.Experiment{
				Experiment: &experimentproto.Experiment{
					Id:     "id-1",
					Kind:   experimentproto.Experiment_HOLDOUT,
					Status: experimentproto.Experiment_WAITING,
				},
			}, nil)
			// Neither the storage write nor the domain event is expected:
			// the mocks fail the test if either is called.

			resp, err := service.UpdateExperiment(createContextWithEditor(), p.req)
			assert.Nil(t, resp)
			assert.Equal(t, statusExperimentOwnedByHoldoutGroup.Err(), err)
			assert.Equal(t, codes.FailedPrecondition, status.Code(err))
		})
	}
}

func TestDeleteExperimentRejectsHoldoutKind(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	service := createExperimentService(mockController)
	service.dbClient.(*dbmock.MockClient).EXPECT().RunInTransactionV2(
		gomock.Any(), gomock.Any(),
	).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	})
	service.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperiment(
		gomock.Any(), "id-1", "ns0",
	).Return(&domain.Experiment{
		Experiment: &experimentproto.Experiment{
			Id:   "id-1",
			Kind: experimentproto.Experiment_HOLDOUT,
		},
	}, nil)
	// Neither the storage write nor the domain event is expected:
	// the mocks fail the test if either is called.

	resp, err := service.DeleteExperiment(
		createContextWithEditor(),
		&experimentproto.DeleteExperimentRequest{Id: "id-1", EnvironmentId: "ns0"},
	)
	assert.Nil(t, resp)
	assert.Equal(t, statusExperimentOwnedByHoldoutGroup.Err(), err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}
