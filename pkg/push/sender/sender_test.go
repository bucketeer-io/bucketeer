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

package sender

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	domaineventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestIsFeaturesLatest(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc           string
		features       *featureproto.Features
		featureID      string
		featureVersion int32
		expected       bool
	}{
		{
			desc: "no feature",
			features: &featureproto.Features{
				Features: []*featureproto.Feature{{Id: "wrong", Version: int32(1)}},
			},
			featureID:      "fid",
			featureVersion: int32(1),
			expected:       false,
		},
		{
			desc: "not the latest version",
			features: &featureproto.Features{
				Features: []*featureproto.Feature{{Id: "fid", Version: int32(1)}},
			},
			featureID:      "fid",
			featureVersion: int32(2),
			expected:       false,
		},
		{
			desc: "the latest version",
			features: &featureproto.Features{
				Features: []*featureproto.Feature{{Id: "fid", Version: int32(2)}},
			},
			featureID:      "fid",
			featureVersion: int32(2),
			expected:       true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &sender{}
			actual := s.isFeaturesLatest(p.features, p.featureID, p.featureVersion)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestExtractFeatureID(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc             string
		input            *domaineventproto.Event
		expectedID       string
		expectedIsTarget bool
	}{
		{
			desc: "not feature entity",
			input: &domaineventproto.Event{
				EntityType: domaineventproto.Event_EXPERIMENT,
				EntityId:   "fid",
				Type:       domaineventproto.Event_FEATURE_VERSION_INCREMENTED,
			},
			expectedID:       "",
			expectedIsTarget: false,
		},
		{
			desc: "not version incremented",
			input: &domaineventproto.Event{
				EntityType: domaineventproto.Event_EXPERIMENT,
				EntityId:   "fid",
				Type:       domaineventproto.Event_FEATURE_DESCRIPTION_CHANGED,
			},
			expectedID:       "",
			expectedIsTarget: false,
		},
		{
			desc: "is target",
			input: &domaineventproto.Event{
				EntityType: domaineventproto.Event_FEATURE,
				EntityId:   "fid",
				Type:       domaineventproto.Event_FEATURE_VERSION_INCREMENTED,
			},
			expectedID:       "fid",
			expectedIsTarget: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &sender{}
			actualID, actualIsTarget := s.extractFeatureID(p.input)
			assert.Equal(t, p.expectedID, actualID)
			assert.Equal(t, p.expectedIsTarget, actualIsTarget)
		})
	}
}

func TestListFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	now := time.Now()
	twentyNineDaysAgo := now.Add(-29 * 24 * time.Hour)
	thirtyOneDaysAgo := now.Add(-31 * 24 * time.Hour)
	ctx := context.TODO()
	envNS := "ns0"

	patterns := []struct {
		desc                 string
		setup                func(*sender)
		environmentNamespace string
		expected             []*featureproto.Feature
		expectedErr          error
	}{
		{
			desc: "listFeatures fails",
			setup: func(s *sender) {
				s.featureClient.(*featureclientmock.MockClient).EXPECT().ListFeatures(
					ctx,
					&featureproto.ListFeaturesRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: envNS,
					},
				).Return(
					nil, errors.New("test"),
				)
			},
			environmentNamespace: envNS,
			expectedErr:          errors.New("test"),
			expected:             nil,
		},
		{
			desc: "success: including off-variation features",
			setup: func(s *sender) {
				s.featureClient.(*featureclientmock.MockClient).EXPECT().ListFeatures(
					ctx,
					&featureproto.ListFeaturesRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: envNS,
					},
				).Return(
					&featureproto.ListFeaturesResponse{
						Features: []*featureproto.Feature{
							{
								Id:      "id-0",
								Enabled: true,
							},
							{
								Id:           "id-1",
								Enabled:      true,
								OffVariation: "",
							},
							{
								Id:           "id-2",
								Enabled:      false,
								OffVariation: "var-2",
							},
							{
								Id:           "id-3",
								Enabled:      false,
								OffVariation: "",
							},
						},
					}, nil,
				)
			},
			environmentNamespace: envNS,
			expectedErr:          nil,
			expected: []*featureproto.Feature{
				{
					Id:      "id-0",
					Enabled: true,
				},
				{
					Id:           "id-1",
					Enabled:      true,
					OffVariation: "",
				},
				{
					Id:           "id-2",
					Enabled:      false,
					OffVariation: "var-2",
				},
			},
		},
		{
			desc: "success: including archived features",
			setup: func(s *sender) {
				s.featureClient.(*featureclientmock.MockClient).EXPECT().ListFeatures(
					ctx,
					&featureproto.ListFeaturesRequest{
						PageSize:             listRequestSize,
						Cursor:               "",
						EnvironmentNamespace: envNS,
					},
				).Return(
					&featureproto.ListFeaturesResponse{
						Features: []*featureproto.Feature{
							{
								Id:       "id-0",
								Enabled:  true,
								Archived: false,
							},
							{
								Id:        "id-1",
								Enabled:   true,
								Archived:  true,
								UpdatedAt: twentyNineDaysAgo.Unix(),
							},
							{
								Id:        "id-2",
								Enabled:   true,
								Archived:  true,
								UpdatedAt: thirtyOneDaysAgo.Unix(),
							},
						},
					}, nil,
				)
			},
			environmentNamespace: envNS,
			expectedErr:          nil,
			expected: []*featureproto.Feature{
				{
					Id:       "id-0",
					Enabled:  true,
					Archived: false,
				},
				{
					Id:        "id-1",
					Enabled:   true,
					Archived:  true,
					UpdatedAt: twentyNineDaysAgo.Unix(),
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &sender{
				featureClient: featureclientmock.NewMockClient(mockController),
			}
			p.setup(s)
			actual, err := s.listFeatures(ctx, p.environmentNamespace)
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}
