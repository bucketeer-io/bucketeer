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
	"google.golang.org/grpc/metadata"

	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestValidateVariationDeletion(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	variationID1 := "variation-1"
	variationID2 := "variation-2"
	variationValue1 := "true"

	patterns := []struct {
		desc             string
		variationChanges []*featureproto.VariationChange
		features         []*featureproto.Feature
		targetFeatureID  string
		expected         error
	}{
		{
			desc:             "success: no variation changes",
			variationChanges: []*featureproto.VariationChange{},
			features:         []*featureproto.Feature{},
			targetFeatureID:  "feature-1",
			expected:         nil,
		},
		{
			desc: "success: no deletion changes",
			variationChanges: []*featureproto.VariationChange{
				{
					ChangeType: featureproto.ChangeType_UPDATE,
					Variation: &featureproto.Variation{
						Id:    variationID1,
						Value: variationValue1,
					},
				},
			},
			features:        []*featureproto.Feature{},
			targetFeatureID: "feature-1",
			expected:        nil,
		},
		{
			desc: "success: no other features using deleted variation",
			variationChanges: []*featureproto.VariationChange{
				{
					ChangeType: featureproto.ChangeType_DELETE,
					Variation: &featureproto.Variation{
						Id:    variationID1,
						Value: variationValue1,
					},
				},
			},
			features:        []*featureproto.Feature{},
			targetFeatureID: "feature-1",
			expected:        nil,
		},
		{
			desc: "error: other feature has prerequisite using deleted variation",
			variationChanges: []*featureproto.VariationChange{
				{
					ChangeType: featureproto.ChangeType_DELETE,
					Variation: &featureproto.Variation{
						Id:    variationID1,
						Value: variationValue1,
					},
				},
			},
			features: []*featureproto.Feature{
				{
					Id: "feature-1", // Target feature (must be included)
					Variations: []*featureproto.Variation{
						{
							Id:    variationID1,
							Value: variationValue1,
						},
					},
				},
				{
					Id: "feature-2", // Dependent feature
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId:   "feature-1",
							VariationId: variationID1,
						},
					},
				},
			},
			targetFeatureID: "feature-1",
			expected:        statusVariationInUseByOtherFeatures.Err(),
		},
		{
			desc: "error: other feature has FEATURE_FLAG rule using deleted variation ID",
			variationChanges: []*featureproto.VariationChange{
				{
					ChangeType: featureproto.ChangeType_DELETE,
					Variation: &featureproto.Variation{
						Id:    variationID1,
						Value: variationValue1,
					},
				},
			},
			features: []*featureproto.Feature{
				{
					Id: "feature-1", // Target feature (must be included)
					Variations: []*featureproto.Variation{
						{
							Id:    variationID1,
							Value: variationValue1,
						},
					},
				},
				{
					Id: "feature-2", // Dependent feature
					Rules: []*featureproto.Rule{
						{
							Clauses: []*featureproto.Clause{
								{
									Operator:  featureproto.Clause_FEATURE_FLAG,
									Attribute: "feature-1",
									Values:    []string{variationID1}, // Fixed: Use variation ID, not value
								},
							},
						},
					},
				},
			},
			targetFeatureID: "feature-1",
			expected:        statusVariationInUseByOtherFeatures.Err(),
		},
		{
			desc: "success: target feature uses deleted variation (should be excluded)",
			variationChanges: []*featureproto.VariationChange{
				{
					ChangeType: featureproto.ChangeType_DELETE,
					Variation: &featureproto.Variation{
						Id:    variationID1,
						Value: variationValue1,
					},
				},
			},
			features: []*featureproto.Feature{
				{
					Id: "feature-1", // Same as targetFeatureID
					Prerequisites: []*featureproto.Prerequisite{
						{
							FeatureId:   "feature-1",
							VariationId: variationID1,
						},
					},
				},
			},
			targetFeatureID: "feature-1",
			expected:        nil,
		},
		{
			desc: "success: different variation ID in FEATURE_FLAG rule",
			variationChanges: []*featureproto.VariationChange{
				{
					ChangeType: featureproto.ChangeType_DELETE,
					Variation: &featureproto.Variation{
						Id:    variationID1,
						Value: variationValue1,
					},
				},
			},
			features: []*featureproto.Feature{
				{
					Id: "feature-2",
					Rules: []*featureproto.Rule{
						{
							Clauses: []*featureproto.Clause{
								{
									Operator:  featureproto.Clause_FEATURE_FLAG,
									Attribute: "feature-1",
									Values:    []string{variationID2}, // Fixed: Different variation ID
								},
							},
						},
					},
				},
			},
			targetFeatureID: "feature-1",
			expected:        nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := validateVariationDeletion(p.variationChanges, p.features, p.targetFeatureID)
			assert.Equal(t, p.expected, err)
		})
	}
}
