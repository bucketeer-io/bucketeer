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

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachetesting "github.com/bucketeer-io/bucketeer/v2/pkg/cache/testing"
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
			desc: "err: failed to list features for environment",
			setup: func(fc *featureFlagCacher) {
				fc.ftStorage.(*mockftstorage.MockFeatureStorage).EXPECT().
					ListFeaturesByEnvironment(gomock.Any(), envID).
					Return(nil, internalErr)
			},
			expectedErr: internalErr,
		},
		{
			desc: "success: empty features for environment",
			setup: func(fc *featureFlagCacher) {
				fc.ftStorage.(*mockftstorage.MockFeatureStorage).EXPECT().
					ListFeaturesByEnvironment(gomock.Any(), envID).
					Return([]*ftproto.Feature{}, nil)
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
					PutIfNewer(gomock.Any(), envID, fixedGeneration).
					Return(true, nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: refresh cache for specific environment",
			setup: func(fc *featureFlagCacher) {
				fc.ftStorage.(*mockftstorage.MockFeatureStorage).EXPECT().
					ListFeaturesByEnvironment(gomock.Any(), envID).
					Return([]*ftproto.Feature{
						{Id: "ft-id-1", OffVariation: "var-1"},
						{Id: "ft-id-2", OffVariation: "var-2"},
					}, nil)
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
					PutIfNewer(gomock.Any(), envID, fixedGeneration).
					Return(true, nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: filters out old archived features",
			setup: func(fc *featureFlagCacher) {
				fc.ftStorage.(*mockftstorage.MockFeatureStorage).EXPECT().
					ListFeaturesByEnvironment(gomock.Any(), envID).
					Return([]*ftproto.Feature{
						{Id: "ft-id-1", OffVariation: "var-1"}, // valid
						{
							Id:           "ft-id-2",
							Archived:     true,
							OffVariation: "var-2",
							UpdatedAt:    time.Now().AddDate(0, 0, -31).Unix(), // older than 30 days
						},
						{Id: "ft-id-3", OffVariation: "var-3"}, // valid
					}, nil)
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
					PutIfNewer(gomock.Any(), envID, fixedGeneration).
					DoAndReturn(func(features *ftproto.Features, envID string, generation int64) (bool, error) {
						// Should only have 2 features (ft-id-2 filtered out)
						assert.Len(t, features.Features, 2)
						assert.Equal(t, "ft-id-1", features.Features[0].Id)
						assert.Equal(t, "ft-id-3", features.Features[1].Id)
						return true, nil
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

func TestRefreshEnvironmentCacheWithEmptyEnvID(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	// Test with empty environment ID - this is a valid case
	// because there is an environment with empty ID
	emptyEnvID := ""

	cacher := newFeatureFlagCacherWithMock(t, controller, 1)
	cacher.ftStorage.(*mockftstorage.MockFeatureStorage).EXPECT().
		ListFeaturesByEnvironment(gomock.Any(), emptyEnvID).
		Return([]*ftproto.Feature{
			{Id: "ft-id-1", OffVariation: "var-1"},
		}, nil)
	cacher.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
		PutIfNewer(gomock.Any(), emptyEnvID, fixedGeneration).
		Return(true, nil)

	err := cacher.RefreshEnvironmentCache(context.Background(), emptyEnvID)
	assert.NoError(t, err)
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
					PutIfNewer(gomock.Any(), "env-id-1", fixedGeneration).
					Return(true, nil)
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
					PutIfNewer(gomock.Any(), "env-id-2", fixedGeneration).
					Return(true, nil)
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
					PutIfNewer(features, envID, fixedGeneration).
					Return(true, nil)
			},
		},
		{
			desc: "err: cache put fails (logged but not returned)",
			setup: func(fc *featureFlagCacher) {
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
					PutIfNewer(features, envID, fixedGeneration).
					Return(false, errors.New("cache error"))
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cacher := newFeatureFlagCacherWithMock(t, controller, 1)
			p.setup(cacher)
			// putCache doesn't return error, it just logs and records metrics
			cacher.putCache(features, envID, len(features.Features), fixedGeneration)
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
					PutIfNewer(features, envID, fixedGeneration).
					Return(true, nil)
				fc.caches[1].(*mockcachev3.MockFeaturesCache).EXPECT().
					PutIfNewer(features, envID, fixedGeneration).
					Return(true, nil)
			},
		},
		{
			desc: "partial failure: one cache fails",
			setup: func(fc *featureFlagCacher) {
				fc.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
					PutIfNewer(features, envID, fixedGeneration).
					Return(false, errors.New("cache error"))
				fc.caches[1].(*mockcachev3.MockFeaturesCache).EXPECT().
					PutIfNewer(features, envID, fixedGeneration).
					Return(true, nil)
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cacher := newFeatureFlagCacherWithMock(t, controller, 2)
			p.setup(cacher)
			cacher.putCache(features, envID, len(features.Features), fixedGeneration)
		})
	}
}

// fixedNow is the deterministic clock used by the test cacher; fixedGeneration
// is the generation every PutIfNewer call is expected to carry.
var (
	fixedNow        = time.Date(2026, 10, 9, 1, 2, 3, 456_789_000, time.UTC)
	fixedGeneration = fixedNow.UnixNano()
)

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
		now:       func() time.Time { return fixedNow },
	}
}

func TestNewFeatureFlagCacherUsesWallClock(t *testing.T) {
	t.Parallel()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	backing := []cache.MultiGetCache{cachetesting.NewInMemoryCache(), cachetesting.NewInMemoryCache()}
	c := NewFeatureFlagCacher(nil, backing, logger).(*featureFlagCacher)
	require.NotNil(t, c.now)
	before := time.Now()
	got := c.now()
	assert.False(t, got.Before(before))
	assert.Len(t, c.caches, 2)
}

// End-to-end with real (test) cache backends: an older snapshot from a slow
// batch run must not replace a newer one written by a faster refresh, in
// either arrival order, on every configured instance.
func TestFeatureFlagCacherEndToEndMonotonic(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()
	logger, err := log.NewLogger()
	require.NoError(t, err)

	envID := "env-e2e"
	backing := []cache.MultiGetCache{cachetesting.NewInMemoryCache(), cachetesting.NewInMemoryCache()}
	storage := mockftstorage.NewMockFeatureStorage(controller)
	c := NewFeatureFlagCacher(storage, backing, logger).(*featureFlagCacher)

	clock := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	tick := 0
	c.now = func() time.Time { tick++; return clock.Add(time.Duration(tick) * time.Second) }

	// Run 1 (generation t+1s) sees flag enabled=false.
	storage.EXPECT().ListFeaturesByEnvironment(gomock.Any(), envID).
		Return([]*ftproto.Feature{{Id: "ft", Enabled: false, OffVariation: "v"}}, nil)
	require.NoError(t, c.RefreshEnvironmentCache(context.Background(), envID))

	// Run 2 (generation t+2s) sees flag enabled=true.
	storage.EXPECT().ListFeaturesByEnvironment(gomock.Any(), envID).
		Return([]*ftproto.Feature{{Id: "ft", Enabled: true, OffVariation: "v"}}, nil)
	require.NoError(t, c.RefreshEnvironmentCache(context.Background(), envID))

	// A slow all-env run that started before run 2 (generation forced older)
	// delivers the stale enabled=false snapshot afterwards.
	c.now = func() time.Time { return clock.Add(1500 * time.Millisecond) }
	storage.EXPECT().ListAllEnvironmentFeatures(gomock.Any()).
		Return([]*ftproto.EnvironmentFeature{{
			EnvironmentId: envID,
			Features:      []*ftproto.Feature{{Id: "ft", Enabled: false, OffVariation: "v"}},
		}}, nil)
	require.NoError(t, c.RefreshAllEnvironmentCaches(context.Background()))

	for i, b := range backing {
		got, err := cachev3.NewFeaturesCache(b, 0).Get(envID)
		require.NoError(t, err, "instance %d", i)
		require.Len(t, got.Features, 1)
		assert.True(t, got.Features[0].Enabled, "instance %d must keep the newer snapshot", i)
	}
}

// The generation must be captured before the storage read starts, otherwise a
// read that began earlier than a concurrent writer's could be stamped later
// and win incorrectly.
func TestRefreshEnvironmentCacheGenerationCapturedBeforeRead(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	envID := "env-id-1"
	cacher := newFeatureFlagCacherWithMock(t, controller, 1)
	first := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	calls := 0
	cacher.now = func() time.Time {
		calls++
		return first.Add(time.Duration(calls) * time.Hour)
	}
	// now() is called once (for both the duration metric and the generation)
	// before ListFeaturesByEnvironment; the storage call advances a fake clock
	// to prove the generation was not taken afterwards.
	cacher.ftStorage.(*mockftstorage.MockFeatureStorage).EXPECT().
		ListFeaturesByEnvironment(gomock.Any(), envID).
		DoAndReturn(func(context.Context, string) ([]*ftproto.Feature, error) {
			assert.Equal(t, 1, calls, "generation must be captured before the read")
			return []*ftproto.Feature{{Id: "ft-id-1", OffVariation: "v"}}, nil
		})
	cacher.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
		PutIfNewer(gomock.Any(), envID, first.Add(time.Hour).UnixNano()).
		Return(true, nil)

	require.NoError(t, cacher.RefreshEnvironmentCache(context.Background(), envID))
}

// A single run of RefreshAllEnvironmentCaches stamps every environment with
// the same generation, taken before ListAllEnvironmentFeatures.
func TestRefreshAllEnvironmentCachesSharesOneGeneration(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	cacher := newFeatureFlagCacherWithMock(t, controller, 1)
	first := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	calls := 0
	cacher.now = func() time.Time {
		calls++
		return first.Add(time.Duration(calls) * time.Hour)
	}
	cacher.ftStorage.(*mockftstorage.MockFeatureStorage).EXPECT().
		ListAllEnvironmentFeatures(gomock.Any()).
		DoAndReturn(func(context.Context) ([]*ftproto.EnvironmentFeature, error) {
			assert.Equal(t, 1, calls)
			return []*ftproto.EnvironmentFeature{
				{EnvironmentId: "env-id-1", Features: []*ftproto.Feature{{Id: "ft-1", OffVariation: "v"}}},
				{EnvironmentId: "env-id-2", Features: []*ftproto.Feature{{Id: "ft-2", OffVariation: "v"}}},
				{EnvironmentId: "env-id-3", Features: []*ftproto.Feature{}},
			}, nil
		})
	expectedGen := first.Add(time.Hour).UnixNano()
	for _, env := range []string{"env-id-1", "env-id-2", "env-id-3"} {
		cacher.caches[0].(*mockcachev3.MockFeaturesCache).EXPECT().
			PutIfNewer(gomock.Any(), env, expectedGen).
			Return(true, nil)
	}
	require.NoError(t, cacher.RefreshAllEnvironmentCaches(context.Background()))
}

// Not parallel: the stale counter is process-global (labelled only by writer),
// and other tests in this package (e.g. the end-to-end monotonic test) also
// increment it. Running in the sequential phase keeps the deltas exact.
func TestPutCacheStaleRejection(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	features := &ftproto.Features{Features: []*ftproto.Feature{{Id: "ft-id-1"}}}

	patterns := []struct {
		desc                 string
		envID                string
		results              []putResult
		expectedStaleDelta   float64
		expectedSuccessDelta float64
		expectedFailDelta    float64
		expectGaugeUpdated   bool
	}{
		{
			desc:                 "rejected as stale: counted, not a failure, gauge untouched",
			envID:                "stale-env-1",
			results:              []putResult{{accepted: false}},
			expectedStaleDelta:   1,
			expectedSuccessDelta: 1,
			expectGaugeUpdated:   false,
		},
		{
			desc:                 "accepted: gauge updated",
			envID:                "stale-env-2",
			results:              []putResult{{accepted: true}},
			expectedSuccessDelta: 1,
			expectGaugeUpdated:   true,
		},
		{
			desc:                 "one accepted, one stale: success and gauge updated, one stale counted",
			envID:                "stale-env-3",
			results:              []putResult{{accepted: true}, {accepted: false}},
			expectedStaleDelta:   1,
			expectedSuccessDelta: 1,
			expectGaugeUpdated:   true,
		},
		{
			desc:               "one stale, one error: failure recorded, stale still counted, gauge untouched",
			envID:              "stale-env-4",
			results:            []putResult{{accepted: false}, {err: errors.New("cache error")}},
			expectedStaleDelta: 1,
			expectedFailDelta:  1,
			expectGaugeUpdated: false,
		},
		{
			desc:                 "all stale across instances: counted per instance",
			envID:                "stale-env-5",
			results:              []putResult{{accepted: false}, {accepted: false}},
			expectedStaleDelta:   2,
			expectedSuccessDelta: 1,
			expectGaugeUpdated:   false,
		},
		{
			desc:               "accepted and error: failure wins, gauge untouched",
			envID:              "stale-env-6",
			results:            []putResult{{accepted: true}, {err: errors.New("cache error")}},
			expectedFailDelta:  1,
			expectGaugeUpdated: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cacher := newFeatureFlagCacherWithMock(t, controller, len(p.results))
			for i, r := range p.results {
				cacher.caches[i].(*mockcachev3.MockFeaturesCache).EXPECT().
					PutIfNewer(features, p.envID, fixedGeneration).
					Return(r.accepted, r.err)
			}

			staleCounter := cachev3.FeaturesPutRejectedStaleCounter(cachev3.FeaturesWriterBatchCacher)
			successCounter := cachePutCounter.WithLabelValues(cacherTypeFeatureFlag, p.envID, codeSuccess)
			failCounter := cachePutCounter.WithLabelValues(cacherTypeFeatureFlag, p.envID, codeFail)
			staleBefore := testutil.ToFloat64(staleCounter)
			successBefore := testutil.ToFloat64(successCounter)
			failBefore := testutil.ToFloat64(failCounter)

			cacher.putCache(features, p.envID, 7, fixedGeneration)

			assert.Equal(t, p.expectedStaleDelta, testutil.ToFloat64(staleCounter)-staleBefore, "stale")
			assert.Equal(t, p.expectedSuccessDelta, testutil.ToFloat64(successCounter)-successBefore, "success")
			assert.Equal(t, p.expectedFailDelta, testutil.ToFloat64(failCounter)-failBefore, "fail")
			gauge := testutil.ToFloat64(featuresUpdatedGauge.WithLabelValues(cacherTypeFeatureFlag, p.envID))
			if p.expectGaugeUpdated {
				assert.Equal(t, float64(7), gauge)
			} else {
				assert.Equal(t, float64(0), gauge)
			}
		})
	}
}

type putResult struct {
	accepted bool
	err      error
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
