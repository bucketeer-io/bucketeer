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
//

package cacher

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"

	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	mockcachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var (
	ftsValid = []*ftproto.Feature{
		{
			Id:           "ft-id-1",
			Archived:     true,
			OffVariation: "variation-id-1",
			UpdatedAt:    time.Now().AddDate(0, 0, -20).Unix(),
		},
		{
			Id:           "ft-id-2",
			Archived:     true,
			OffVariation: "variation-id-2",
			UpdatedAt:    time.Now().AddDate(0, 0, -30).Unix(),
		},
		{
			Id:           "ft-id-3",
			Archived:     true,
			OffVariation: "variation-id-3",
			UpdatedAt:    time.Now().AddDate(0, 0, -10).Unix(),
		},
	}

	ftsOlderThan30Days = []*ftproto.Feature{
		{
			Id:           "ft-id-1",
			Archived:     true,
			OffVariation: "variation-id-1",
			UpdatedAt:    time.Now().AddDate(0, 0, -20).Unix(),
		},
		{
			Id:           "ft-id-2",
			Archived:     true,
			OffVariation: "variation-id-2",
			UpdatedAt:    time.Now().AddDate(0, 0, -31).Unix(),
		},
		{
			Id:           "ft-id-3",
			Archived:     true,
			OffVariation: "variation-id-3",
			UpdatedAt:    time.Now().AddDate(0, 0, -10).Unix(),
		},
	}

	ftsInvalid = []*ftproto.Feature{
		{
			Id:           "ft-id-1",
			Archived:     true,
			Enabled:      false,
			OffVariation: "",
			UpdatedAt:    time.Now().AddDate(0, 0, -20).Unix(),
		},
		{
			Id:           "ft-id-3",
			Archived:     true,
			Enabled:      false,
			OffVariation: "variation-id-3",
			UpdatedAt:    time.Now().AddDate(0, 0, -10).Unix(),
		},
	}
)

func TestFeatureFlagsPutCache(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	envID := "env-id"
	features := &ftproto.Features{
		Features: []*ftproto.Feature{
			{
				Id: "ft-id-1",
			},
			{
				Id: "ft-id-2",
			},
		},
	}

	patterns := []struct {
		desc     string
		setup    func(*featureFlagCacher)
		expected int
	}{
		{
			desc: "err: error at index 0",
			setup: func(fc *featureFlagCacher) {
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().Put(features, envID).
					Return(errors.New("internal error"))
				fc.caches[1].(*mockcachev3.MockFeaturesCache).EXPECT().Put(features, envID).
					Return(nil)
			},
			expected: 1,
		},
		{
			desc: "err: error at index 1",
			setup: func(fc *featureFlagCacher) {
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().Put(features, envID).
					Return(nil)
				fc.caches[1].(*mockcachev3.MockFeaturesCache).EXPECT().Put(features, envID).
					Return(errors.New("internal error"))
			},
			expected: 1,
		},
		{
			desc: "success",
			setup: func(fc *featureFlagCacher) {
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().Put(features, envID).
					Return(nil)
				fc.caches[1].(*mockcachev3.MockFeaturesCache).EXPECT().Put(features, envID).
					Return(nil)
			},
			expected: 2,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cacher := newFeatureFlagCacher(t, controller)
			p.setup(cacher)
			updatedInstances := cacher.putCache(features, envID)
			assert.Equal(t, p.expected, updatedInstances)
		})
	}
}

func TestRemoveOldFeatures(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	patterns := []struct {
		desc     string
		input    []*ftproto.Feature
		expected []*ftproto.Feature
	}{
		{
			desc:  "remove old feature",
			input: ftsOlderThan30Days,
			expected: []*ftproto.Feature{
				{
					Id:           "ft-id-1",
					Archived:     true,
					OffVariation: "variation-id-1",
					UpdatedAt:    time.Now().AddDate(0, 0, -20).Unix(),
				},
				{
					Id:           "ft-id-3",
					Archived:     true,
					OffVariation: "variation-id-3",
					UpdatedAt:    time.Now().AddDate(0, 0, -10).Unix(),
				},
			},
		},
		{
			desc:  "remove invalid feature",
			input: ftsInvalid,
			expected: []*ftproto.Feature{
				{
					Id:           "ft-id-3",
					Archived:     true,
					OffVariation: "variation-id-3",
					UpdatedAt:    time.Now().AddDate(0, 0, -10).Unix(),
				},
			},
		},
		{
			desc:     "remove nothing",
			input:    ftsValid,
			expected: ftsValid,
		},
		{
			desc: "remove all",
			input: []*ftproto.Feature{
				{
					Id:           "ft-id-1",
					Archived:     true,
					OffVariation: "variation-id-1",
					UpdatedAt:    time.Now().AddDate(0, 0, -31).Unix(),
				},
				{
					Id:           "ft-id-2",
					Archived:     true,
					OffVariation: "variation-id-2",
					UpdatedAt:    time.Now().AddDate(0, 0, -31).Unix(),
				},
				{
					Id:           "ft-id-3",
					Archived:     true,
					OffVariation: "variation-id-3",
					UpdatedAt:    time.Now().AddDate(0, 0, -31).Unix(),
				},
				{
					Id:           "ft-id-4",
					Archived:     true,
					OffVariation: "",
					UpdatedAt:    time.Now().AddDate(0, 0, -20).Unix(),
				},
			},
			expected: []*ftproto.Feature{},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cacher := &featureFlagCacher{}
			actual := cacher.removeOldFeatures(p.input)
			assert.True(t, compareFeatureSlicesSerialized(t, p.expected, actual))
		})
	}
}

func compareFeatureSlicesSerialized(t *testing.T, slice1, slice2 []*ftproto.Feature) bool {
	t.Helper()
	// Check if the slices have different lengths
	if len(slice1) != len(slice2) {
		t.Fatalf("Different slice size")
	}
	// Serialize and compare each message in the slices
	for i := 0; i < len(slice1); i++ {
		data1, err := proto.Marshal(slice1[i])
		if err != nil {
			t.Fatalf("Failed to serialize slice1[%d]: %v", i, err)
		}
		data2, err := proto.Marshal(slice2[i])
		if err != nil {
			t.Fatalf("Failed to serialize slice2[%d]: %v", i, err)
		}
		if !bytes.Equal(data1, data2) {
			return false
		}
	}
	// If no differences were found, the slices are equal
	return true
}

func newFeatureFlagCacher(t *testing.T, controller *gomock.Controller) *featureFlagCacher {
	t.Helper()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &featureFlagCacher{
		caches: []cachev3.FeaturesCache{
			mockcachev3.NewMockFeaturesCache(controller),
			mockcachev3.NewMockFeaturesCache(controller),
		},
		logger: logger,
	}
}
