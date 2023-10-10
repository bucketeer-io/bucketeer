// Copyright 2023 The Bucketeer Authors.
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
	"go.uber.org/mock/gomock"

	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	"github.com/bucketeer-io/bucketeer/proto/autoops"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestGetRolloutStrategyVariations(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc              string
		feature           *featureproto.Feature
		targetVariationID string
		targetWeight      int32
		expected          []*featureproto.RolloutStrategy_Variation
		expectedErr       error
	}{
		{
			desc: "success: weight is max",
			feature: &featureproto.Feature{
				Variations: []*featureproto.Variation{
					{
						Id: "vid-1",
					},
					{
						Id: "vid-2",
					},
				},
			},
			targetVariationID: "vid-1",
			targetWeight:      totalVariationWeight,
			expected: []*featureproto.RolloutStrategy_Variation{
				{
					Variation: "vid-1",
					Weight:    totalVariationWeight,
				},
				{
					Variation: "vid-2",
					Weight:    0,
				},
			},
		},
		{
			desc: "success: weight is not max",
			feature: &featureproto.Feature{
				Variations: []*featureproto.Variation{
					{
						Id: "vid-1",
					},
					{
						Id: "vid-2",
					},
				},
			},
			targetVariationID: "vid-2",
			targetWeight:      20,
			expected: []*featureproto.RolloutStrategy_Variation{
				{
					Variation: "vid-2",
					Weight:    20,
				},
				{
					Variation: "vid-1",
					Weight:    totalVariationWeight - 20,
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := getRolloutStrategyVariations(p.feature, &autoops.ProgressiveRolloutSchedule{
				Weight: p.targetWeight,
			}, p.targetVariationID)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestUpdateFeatureTargeting(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mockController := gomock.NewController(t)
	environmentNamespace := "en"
	fID := "fid-1"
	patterns := []struct {
		desc              string
		setup             func(*featureclientmock.MockClient)
		feature           *featureproto.Feature
		targetVariationID string
		targetWeight      int32
		expectedErr       error
	}{
		{
			desc: "success: weight is max",
			setup: func(mock *featureclientmock.MockClient) {
				mock.EXPECT().UpdateFeatureTargeting(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			feature: &featureproto.Feature{
				Id: fID,
				Variations: []*featureproto.Variation{
					{
						Id: "vid-1",
					},
					{
						Id: "vid-2",
					},
				},
			},
			targetVariationID: "vid-1",
			targetWeight:      totalVariationWeight,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			featureClientMock := featureclientmock.NewMockClient(mockController)
			if p.setup != nil {
				p.setup(featureClientMock)
			}
			err := updateFeatureTargeting(
				ctx,
				&autoops.ProgressiveRolloutSchedule{
					Weight: p.targetWeight,
				},
				featureClientMock,
				p.feature,
				p.targetVariationID,
				environmentNamespace,
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
