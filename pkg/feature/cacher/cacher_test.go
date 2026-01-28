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

package cacher

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"

	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	mockcachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	mockftstorage "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestRefreshEnvironmentCache(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	envID := "env-id-1"
	internalErr := errors.New("internal error")

	patterns := []struct {
		desc        string
		setup       func(*featureFlagCacher)
		expectedErr error
	}{
		{
			desc: "err: failed to list all environment features",
			setup: func(fc *featureFlagCacher) {
				fc.ftStorage.(*mockftstorage.MockFeatureStorage).EXPECT().
					ListAllEnvironmentFeatures(gomock.Any()).
					Return(nil, internalErr)
			},
			expectedErr: internalErr,
		},
		{
			desc: "success: environment not found (empty cache)",
			setup: func(fc *featureFlagCacher) {
				fc.ftStorage.(*mockftstorage.MockFeatureStorage).EXPECT().
					ListAllEnvironmentFeatures(gomock.Any()).
					Return([]*ftproto.EnvironmentFeature{
						{
							EnvironmentId: "other-env-id",
							Features:      []*ftproto.Feature{{Id: "ft-id-1"}},
						},
					}, nil)
				// Should still update cache with empty features
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
					Put(gomock.Any(), envID).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: refresh cache for specific environment",
			setup: func(fc *featureFlagCacher) {
				fc.ftStorage.(*mockftstorage.MockFeatureStorage).EXPECT().
					ListAllEnvironmentFeatures(gomock.Any()).
					Return([]*ftproto.EnvironmentFeature{
						{
							EnvironmentId: envID,
							Features: []*ftproto.Feature{
								{Id: "ft-id-1", OffVariation: "var-1"},
								{Id: "ft-id-2", OffVariation: "var-2"},
							},
						},
						{
							EnvironmentId: "other-env-id",
							Features:      []*ftproto.Feature{{Id: "ft-id-3"}},
						},
					}, nil)
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
					Put(gomock.Any(), envID).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: filters out old archived features",
			setup: func(fc *featureFlagCacher) {
				fc.ftStorage.(*mockftstorage.MockFeatureStorage).EXPECT().
					ListAllEnvironmentFeatures(gomock.Any()).
					Return([]*ftproto.EnvironmentFeature{
						{
							EnvironmentId: envID,
							Features: []*ftproto.Feature{
								{Id: "ft-id-1", OffVariation: "var-1"}, // valid
								{
									Id:           "ft-id-2",
									Archived:     true,
									OffVariation: "var-2",
									UpdatedAt:    time.Now().AddDate(0, 0, -31).Unix(), // older than 30 days
								},
								{Id: "ft-id-3", OffVariation: "var-3"}, // valid
							},
						},
					}, nil)
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
					Put(gomock.Any(), envID).
					DoAndReturn(func(features *ftproto.Features, envID string) error {
						// Should only have 2 features (ft-id-2 filtered out)
						assert.Len(t, features.Features, 2)
						assert.Equal(t, "ft-id-1", features.Features[0].Id)
						assert.Equal(t, "ft-id-3", features.Features[1].Id)
						return nil
					})
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cacher := newFeatureFlagCacherWithMock(t, controller, 1)
			p.setup(cacher)
			err := cacher.RefreshEnvironmentCache(context.Background(), envID)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestRefreshAllEnvironmentCaches(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	internalErr := errors.New("internal error")

	patterns := []struct {
		desc        string
		setup       func(*featureFlagCacher)
		expectedErr error
	}{
		{
			desc: "err: failed to list all environment features",
			setup: func(fc *featureFlagCacher) {
				fc.ftStorage.(*mockftstorage.MockFeatureStorage).EXPECT().
					ListAllEnvironmentFeatures(gomock.Any()).
					Return(nil, internalErr)
			},
			expectedErr: internalErr,
		},
		{
			desc: "success: refresh cache for all environments",
			setup: func(fc *featureFlagCacher) {
				fc.ftStorage.(*mockftstorage.MockFeatureStorage).EXPECT().
					ListAllEnvironmentFeatures(gomock.Any()).
					Return([]*ftproto.EnvironmentFeature{
						{
							EnvironmentId: "env-id-1",
							Features: []*ftproto.Feature{
								{Id: "ft-id-1", OffVariation: "var-1"},
							},
						},
						{
							EnvironmentId: "env-id-2",
							Features: []*ftproto.Feature{
								{Id: "ft-id-2", OffVariation: "var-2"},
							},
						},
					}, nil)
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
					Put(gomock.Any(), "env-id-1").
					Return(nil)
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
					Put(gomock.Any(), "env-id-2").
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: empty environments",
			setup: func(fc *featureFlagCacher) {
				fc.ftStorage.(*mockftstorage.MockFeatureStorage).EXPECT().
					ListAllEnvironmentFeatures(gomock.Any()).
					Return([]*ftproto.EnvironmentFeature{}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cacher := newFeatureFlagCacherWithMock(t, controller, 1)
			p.setup(cacher)
			err := cacher.RefreshAllEnvironmentCaches(context.Background())
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestRemoveOldFeatures(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		input    []*ftproto.Feature
		expected []*ftproto.Feature
	}{
		{
			desc: "remove archived feature older than 30 days",
			input: []*ftproto.Feature{
				{
					Id:           "ft-id-1",
					Archived:     true,
					OffVariation: "var-1",
					UpdatedAt:    time.Now().AddDate(0, 0, -20).Unix(),
				},
				{
					Id:           "ft-id-2",
					Archived:     true,
					OffVariation: "var-2",
					UpdatedAt:    time.Now().AddDate(0, 0, -31).Unix(), // older than 30 days
				},
				{
					Id:           "ft-id-3",
					Archived:     true,
					OffVariation: "var-3",
					UpdatedAt:    time.Now().AddDate(0, 0, -10).Unix(),
				},
			},
			expected: []*ftproto.Feature{
				{
					Id:           "ft-id-1",
					Archived:     true,
					OffVariation: "var-1",
					UpdatedAt:    time.Now().AddDate(0, 0, -20).Unix(),
				},
				{
					Id:           "ft-id-3",
					Archived:     true,
					OffVariation: "var-3",
					UpdatedAt:    time.Now().AddDate(0, 0, -10).Unix(),
				},
			},
		},
		{
			desc: "remove disabled feature with empty off variation",
			input: []*ftproto.Feature{
				{
					Id:           "ft-id-1",
					Archived:     true,
					Enabled:      false,
					OffVariation: "", // empty
					UpdatedAt:    time.Now().AddDate(0, 0, -20).Unix(),
				},
				{
					Id:           "ft-id-2",
					Archived:     true,
					Enabled:      false,
					OffVariation: "var-2",
					UpdatedAt:    time.Now().AddDate(0, 0, -10).Unix(),
				},
			},
			expected: []*ftproto.Feature{
				{
					Id:           "ft-id-2",
					Archived:     true,
					OffVariation: "var-2",
					UpdatedAt:    time.Now().AddDate(0, 0, -10).Unix(),
				},
			},
		},
		{
			desc: "keep all valid features",
			input: []*ftproto.Feature{
				{
					Id:           "ft-id-1",
					Archived:     false,
					OffVariation: "var-1",
				},
				{
					Id:           "ft-id-2",
					Archived:     true,
					OffVariation: "var-2",
					UpdatedAt:    time.Now().AddDate(0, 0, -20).Unix(), // within 30 days
				},
			},
			expected: []*ftproto.Feature{
				{
					Id:           "ft-id-1",
					Archived:     false,
					OffVariation: "var-1",
				},
				{
					Id:           "ft-id-2",
					Archived:     true,
					OffVariation: "var-2",
					UpdatedAt:    time.Now().AddDate(0, 0, -20).Unix(),
				},
			},
		},
		{
			desc:     "empty input",
			input:    []*ftproto.Feature{},
			expected: []*ftproto.Feature{},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cacher := &featureFlagCacher{}
			actual := cacher.removeOldFeatures(p.input)
			assert.True(t, compareFeatureSlices(t, p.expected, actual))
		})
	}
}

func TestPutCache(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	envID := "env-id"
	features := &ftproto.Features{
		Features: []*ftproto.Feature{
			{Id: "ft-id-1"},
			{Id: "ft-id-2"},
		},
	}

	patterns := []struct {
		desc  string
		setup func(*featureFlagCacher)
	}{
		{
			desc: "success: put to single cache",
			setup: func(fc *featureFlagCacher) {
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
					Put(features, envID).
					Return(nil)
			},
		},
		{
			desc: "err: cache put fails (logged but not returned)",
			setup: func(fc *featureFlagCacher) {
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
					Put(features, envID).
					Return(errors.New("cache error"))
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cacher := newFeatureFlagCacherWithMock(t, controller, 1)
			p.setup(cacher)
			// putCache doesn't return error, it just logs and records metrics
			cacher.putCache(features, envID, len(features.Features))
		})
	}
}

func TestPutCacheMultipleInstances(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	envID := "env-id"
	features := &ftproto.Features{
		Features: []*ftproto.Feature{{Id: "ft-id-1"}},
	}

	patterns := []struct {
		desc  string
		setup func(*featureFlagCacher)
	}{
		{
			desc: "success: put to multiple caches",
			setup: func(fc *featureFlagCacher) {
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
					Put(features, envID).
					Return(nil)
				fc.caches[1].(*mockcachev3.MockFeaturesCache).EXPECT().
					Put(features, envID).
					Return(nil)
			},
		},
		{
			desc: "partial failure: one cache fails",
			setup: func(fc *featureFlagCacher) {
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
					Put(features, envID).
					Return(errors.New("cache error"))
				fc.caches[1].(*mockcachev3.MockFeaturesCache).EXPECT().
					Put(features, envID).
					Return(nil)
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cacher := newFeatureFlagCacherWithMock(t, controller, 2)
			p.setup(cacher)
			cacher.putCache(features, envID, len(features.Features))
		})
	}
}

func newFeatureFlagCacherWithMock(t *testing.T, controller *gomock.Controller, numCaches int) *featureFlagCacher {
	t.Helper()
	logger, err := log.NewLogger()
	require.NoError(t, err)

	caches := make([]cachev3.FeaturesCache, numCaches)
	for i := 0; i < numCaches; i++ {
		caches[i] = mockcachev3.NewMockFeaturesCache(controller)
	}

	return &featureFlagCacher{
		ftStorage: mockftstorage.NewMockFeatureStorage(controller),
		caches:    caches,
		logger:    logger,
	}
}

func compareFeatureSlices(t *testing.T, slice1, slice2 []*ftproto.Feature) bool {
	t.Helper()
	if len(slice1) != len(slice2) {
		t.Logf("Different slice size: %d vs %d", len(slice1), len(slice2))
		return false
	}
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
	return true
}
