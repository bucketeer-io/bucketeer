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

package experimentcalc

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	envclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	ecclient "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/client/mock"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/experimentcalculator/stan"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/experimentcalculator/storage/v2/mock"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	holdoutdomain "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/domain"
	metricsmock "github.com/bucketeer-io/bucketeer/v2/pkg/metrics/mock"
	"github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func holdoutVariationResults(inHoldout, notInHoldout int64) []*eventcounter.VariationResult {
	return []*eventcounter.VariationResult{
		vr(holdoutdomain.InHoldoutVariationID, inHoldout),
		vr(holdoutdomain.NotInHoldoutVariationID, notInHoldout),
	}
}

func TestHoldoutExpectedFractions(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc      string
		group     *featureproto.HoldoutGroup
		expected  map[string]float64
		expectErr error
	}{
		{desc: "no group", group: nil, expectErr: errSRMHoldoutGroupMissing},
		{
			// The comparison range has the size of the holdout range, so the
			// split is even for every percentage.
			desc: "hash group with a small percentage",
			group: &featureproto.HoldoutGroup{
				AssignmentType: featureproto.HoldoutGroup_HASH,
				Percentage:     5,
			},
			expected: map[string]float64{
				holdoutdomain.InHoldoutVariationID:    0.5,
				holdoutdomain.NotInHoldoutVariationID: 0.5,
			},
		},
		{
			desc: "hash group with 50 percent",
			group: &featureproto.HoldoutGroup{
				AssignmentType: featureproto.HoldoutGroup_HASH,
				Percentage:     50,
			},
			expected: map[string]float64{
				holdoutdomain.InHoldoutVariationID:    0.5,
				holdoutdomain.NotInHoldoutVariationID: 0.5,
			},
		},
		{
			desc: "segment group",
			group: &featureproto.HoldoutGroup{
				AssignmentType: featureproto.HoldoutGroup_SEGMENT,
				SegmentId:      "segment-1",
			},
			expectErr: errSRMHoldoutNotHash,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			fractions, err := holdoutExpectedFractions(p.group)
			assert.Equal(t, p.expectErr, err)
			assert.Equal(t, p.expected, fractions)
		})
	}
}

func TestComputeHoldoutSRM(t *testing.T) {
	t.Parallel()
	hash := &featureproto.HoldoutGroup{AssignmentType: featureproto.HoldoutGroup_HASH, Percentage: 5}

	t.Run("an even split is OK", func(t *testing.T) {
		t.Parallel()
		got := computeHoldoutSRM(holdoutVariationResults(5000, 4950), hash, DefaultSRMThreshold)
		assert.Equal(t, eventcounter.SrmResult_OK, got.Status)
		assert.EqualValues(t, 1, got.DegreesOfFreedom)
		// The same counts as TestComputeSRM_BalancedSplit_ReportsOK.
		assert.InDelta(t, 0.2513, got.ChiSquare, 1e-3)
		require.Len(t, got.Variations, 2)
		assert.Equal(t, holdoutdomain.InHoldoutVariationID, got.Variations[0].VariationId)
		assert.EqualValues(t, 5000, got.Variations[0].ObservedUserCount)
		assert.InDelta(t, 0.5, got.Variations[0].ExpectedWeight, 1e-9)
		assert.InDelta(t, 4975.0, got.Variations[0].ExpectedUserCount, 1e-9)
		assert.Equal(t, holdoutdomain.NotInHoldoutVariationID, got.Variations[1].VariationId)
		assert.EqualValues(t, 4950, got.Variations[1].ObservedUserCount)
	})

	// The split an environment gets when part of its held-out users are not
	// reported with the HOLDOUT reason: in_holdout is short of users.
	t.Run("too few held-out users is a MISMATCH", func(t *testing.T) {
		t.Parallel()
		got := computeHoldoutSRM(holdoutVariationResults(4700, 5300), hash, DefaultSRMThreshold)
		assert.Equal(t, eventcounter.SrmResult_MISMATCH, got.Status)
		assert.Less(t, got.PValue, DefaultSRMThreshold)
	})

	t.Run("a segment group is SKIPPED with the reason", func(t *testing.T) {
		t.Parallel()
		segment := &featureproto.HoldoutGroup{AssignmentType: featureproto.HoldoutGroup_SEGMENT}
		// Counts that would be a mismatch against an even split.
		got := computeHoldoutSRM(holdoutVariationResults(100, 9900), segment, DefaultSRMThreshold)
		assert.Equal(t, eventcounter.SrmResult_SKIPPED, got.Status)
		assert.Equal(t,
			"holdout group assigns users by segment (no expected split to test against)",
			got.SkipReason,
		)
		assert.InDelta(t, DefaultSRMThreshold, got.Threshold, 1e-9)
	})

	t.Run("no group is SKIPPED with the reason", func(t *testing.T) {
		t.Parallel()
		got := computeHoldoutSRM(holdoutVariationResults(5000, 5000), nil, DefaultSRMThreshold)
		assert.Equal(t, eventcounter.SrmResult_SKIPPED, got.Status)
		assert.Equal(t, "holdout group not available", got.SkipReason)
	})

	t.Run("too few users is SKIPPED like for a flag", func(t *testing.T) {
		t.Parallel()
		got := computeHoldoutSRM(holdoutVariationResults(3, 2), hash, DefaultSRMThreshold)
		assert.Equal(t, eventcounter.SrmResult_SKIPPED, got.Status)
		assert.Contains(t, got.SkipReason, errSRMInsufficientSamples.Error())
	})
}

func newHoldoutSRMCalculator(t *testing.T) (*ExperimentCalculator, *featureclient.MockClient) {
	t.Helper()
	mockController := gomock.NewController(t)
	registerer := metricsmock.NewMockRegisterer(mockController)
	registerer.EXPECT().MustRegister(gomock.Any()).Return().AnyTimes()
	mockFeatureClient := featureclient.NewMockClient(mockController)
	return NewExperimentCalculator(
		stan.NewStan("localhost", "8080", registerer, zap.NewNop()),
		stanModelID,
		envclient.NewMockClient(mockController),
		ecclient.NewMockClient(mockController),
		experimentclient.NewMockClient(mockController),
		mockFeatureClient,
		storagemock.NewMockExperimentResultStorage(mockController),
		registerer,
		jpLocation,
		zap.NewNop(),
	), mockFeatureClient
}

// The feature ID of an internal experiment is the ID of its group. The SRM
// check reads the group with that ID and never asks for a flag with it: the
// mock has no expectation on GetFeature, so such a call fails the test.
func TestComputeExperimentSRM_HoldoutExperiment(t *testing.T) {
	t.Parallel()
	exp := &experimentproto.Experiment{
		Id:             "exp-internal",
		Kind:           experimentproto.Experiment_HOLDOUT,
		FeatureId:      "hg-1",
		FeatureVersion: holdoutdomain.InternalExperimentFeatureVersion,
	}
	request := &featureproto.GetHoldoutGroupRequest{Id: "hg-1", EnvironmentId: "env-1"}
	goalResults := func(inHoldout, notInHoldout int64) []*eventcounter.GoalResult {
		return []*eventcounter.GoalResult{{VariationResults: holdoutVariationResults(inHoldout, notInHoldout)}}
	}

	t.Run("hash group with an even split is OK", func(t *testing.T) {
		t.Parallel()
		calc, client := newHoldoutSRMCalculator(t)
		client.EXPECT().GetHoldoutGroup(gomock.Any(), request).Return(&featureproto.GetHoldoutGroupResponse{
			HoldoutGroup: &featureproto.HoldoutGroup{
				Id:             "hg-1",
				AssignmentType: featureproto.HoldoutGroup_HASH,
				Percentage:     5,
			},
		}, nil)
		res := calc.computeExperimentSRM(context.Background(), "env-1", exp, goalResults(5000, 4950))
		require.NotNil(t, res)
		assert.Equal(t, eventcounter.SrmResult_OK, res.Status)
		assert.Len(t, res.Variations, 2)
	})

	t.Run("hash group with an uneven split is a MISMATCH", func(t *testing.T) {
		t.Parallel()
		calc, client := newHoldoutSRMCalculator(t)
		client.EXPECT().GetHoldoutGroup(gomock.Any(), request).Return(&featureproto.GetHoldoutGroupResponse{
			HoldoutGroup: &featureproto.HoldoutGroup{
				Id:             "hg-1",
				AssignmentType: featureproto.HoldoutGroup_HASH,
				Percentage:     5,
			},
		}, nil)
		res := calc.computeExperimentSRM(context.Background(), "env-1", exp, goalResults(4700, 5300))
		require.NotNil(t, res)
		assert.Equal(t, eventcounter.SrmResult_MISMATCH, res.Status)
	})

	t.Run("segment group is SKIPPED", func(t *testing.T) {
		t.Parallel()
		calc, client := newHoldoutSRMCalculator(t)
		client.EXPECT().GetHoldoutGroup(gomock.Any(), request).Return(&featureproto.GetHoldoutGroupResponse{
			HoldoutGroup: &featureproto.HoldoutGroup{
				Id:             "hg-1",
				AssignmentType: featureproto.HoldoutGroup_SEGMENT,
				SegmentId:      "segment-1",
			},
		}, nil)
		res := calc.computeExperimentSRM(context.Background(), "env-1", exp, goalResults(100, 9900))
		require.NotNil(t, res)
		assert.Equal(t, eventcounter.SrmResult_SKIPPED, res.Status)
		assert.Equal(t,
			"holdout group assigns users by segment (no expected split to test against)",
			res.SkipReason,
		)
	})

	t.Run("a response without a group is SKIPPED", func(t *testing.T) {
		t.Parallel()
		calc, client := newHoldoutSRMCalculator(t)
		client.EXPECT().GetHoldoutGroup(gomock.Any(), request).Return(&featureproto.GetHoldoutGroupResponse{}, nil)
		res := calc.computeExperimentSRM(context.Background(), "env-1", exp, goalResults(5000, 5000))
		require.NotNil(t, res)
		assert.Equal(t, eventcounter.SrmResult_SKIPPED, res.Status)
		assert.Equal(t, "holdout group not available", res.SkipReason)
	})

	t.Run("a failed fetch is SKIPPED without the error text", func(t *testing.T) {
		t.Parallel()
		calc, client := newHoldoutSRMCalculator(t)
		client.EXPECT().GetHoldoutGroup(gomock.Any(), request).Return(nil, status.Error(codes.Internal,
			"internal: backend db connection refused (host=db-internal-01:3306, user=admin)"))
		res := calc.computeExperimentSRM(context.Background(), "env-1", exp, goalResults(5000, 5000))
		require.NotNil(t, res)
		assert.Equal(t, eventcounter.SrmResult_SKIPPED, res.Status)
		assert.Equal(t, "could not fetch holdout group", res.SkipReason)
		for _, leak := range []string{"db-internal-01", "admin", "connection refused", "Internal"} {
			assert.False(t, strings.Contains(res.SkipReason, leak))
		}
		assert.InDelta(t, DefaultSRMThreshold, res.Threshold, 1e-9)
	})

	t.Run("no goal results is SKIPPED before the group is fetched", func(t *testing.T) {
		t.Parallel()
		calc, _ := newHoldoutSRMCalculator(t)
		res := calc.computeExperimentSRM(context.Background(), "env-1", exp, nil)
		require.NotNil(t, res)
		assert.Equal(t, eventcounter.SrmResult_SKIPPED, res.Status)
		assert.Equal(t, "no goal results available to derive observed traffic split", res.SkipReason)
	})
}

// An experiment of a flag still reads the flag and never a holdout group.
func TestComputeExperimentSRM_FlagExperimentDoesNotFetchHoldoutGroup(t *testing.T) {
	t.Parallel()
	calc, client := newHoldoutSRMCalculator(t)
	client.EXPECT().GetFeature(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.NotFound, "feature"))
	exp := &experimentproto.Experiment{Id: "exp-1", FeatureId: "feature-1", FeatureVersion: 3}
	res := calc.computeExperimentSRM(context.Background(), "env-1", exp, []*eventcounter.GoalResult{{
		VariationResults: []*eventcounter.VariationResult{vr("vid1", 5000), vr("vid2", 5000)},
	}})
	require.NotNil(t, res)
	assert.Equal(t, "could not fetch feature definition", res.SkipReason)
}
