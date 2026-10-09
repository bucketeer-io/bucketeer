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

package v3

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachemock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/mock"
	redis "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	tag           = "bucketeer-tag"
	environmentId = "bucketeer-environment"
)

func TestGetFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	features := createFeatures(t)
	dataFeatures := marshalMessage(t, features)
	key := fmt.Sprintf("%s:%s", environmentId, featuresKind)

	patterns := []struct {
		desc                string
		setup               func(*featuresCache)
		expectedErr         error
		expectedErrContains string
	}{
		{
			desc: "error_get_not_found",
			setup: func(tf *featuresCache) {
				tf.cache.(*cachemock.MockMultiGetCache).EXPECT().Get(key).Return(nil, cache.ErrNotFound)
			},
			expectedErr: cache.ErrNotFound,
		},
		{
			desc: "error_invalid_type",
			setup: func(tf *featuresCache) {
				tf.cache.(*cachemock.MockMultiGetCache).EXPECT().Get(key).Return("test", nil)
			},
			expectedErr: cache.ErrInvalidType,
		},
		{
			desc: "error_unmarshal",
			setup: func(tf *featuresCache) {
				// A varint field header followed by nothing is a truncated message.
				tf.cache.(*cachemock.MockMultiGetCache).EXPECT().Get(key).Return([]byte{0x0a}, nil)
			},
			expectedErrContains: "proto",
		},
		{
			desc: "success",
			setup: func(tf *featuresCache) {
				tf.cache.(*cachemock.MockMultiGetCache).EXPECT().Get(key).Return(dataFeatures, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			tf := newFeaturesCache(t, mockController)
			p.setup(tf)
			actual, err := tf.Get(environmentId)
			if p.expectedErrContains != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), p.expectedErrContains)
				assert.Nil(t, actual)
				return
			}
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				require.Len(t, actual.Features, len(features.Features))
				assert.Equal(t, features.Features[0].Id, actual.Features[0].Id)
				assert.Equal(t, features.Features[0].Name, actual.Features[0].Name)
			}
		})
	}
}

func TestPutFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	features := createFeatures(t)
	dataFeatures := marshalMessage(t, features)
	key := fmt.Sprintf("%s:%s", environmentId, featuresKind)

	patterns := []struct {
		desc        string
		setup       func(*featuresCache)
		input       *featureproto.Features
		expectedErr error
	}{
		{
			desc:        "error_proto_message_nil",
			setup:       nil,
			input:       nil,
			expectedErr: errors.New("features cannot be nil"),
		},
		{
			desc: "success",
			setup: func(tf *featuresCache) {
				tf.cache.(*cachemock.MockMultiGetCache).EXPECT().Put(key, dataFeatures, time.Duration(0)).Return(nil)
			},
			input:       features,
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			tf := newFeaturesCache(t, mockController)
			if p.setup != nil {
				p.setup(tf)
			}
			err := tf.Put(p.input, environmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestPutIfNewerFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	features := createFeatures(t)
	dataFeatures := marshalMessage(t, features)
	key := fmt.Sprintf("%s:%s", environmentId, featuresKind)
	genKey := fmt.Sprintf("{%s:%s}:%s", environmentId, featuresKind, featuresGenerationSuffix)
	const generation = int64(1_791_600_000_000_000_000)
	backendErr := errors.New("redis down")

	patterns := []struct {
		desc             string
		setup            func(*featuresCache)
		input            *featureproto.Features
		ttl              time.Duration
		expectedAccepted bool
		expectedErr      error
	}{
		{
			desc:        "error_proto_message_nil",
			input:       nil,
			expectedErr: errors.New("features cannot be nil"),
		},
		{
			desc: "accepted",
			setup: func(tf *featuresCache) {
				tf.cache.(*cachemock.MockMultiGetCache).EXPECT().
					PutIfNewer(key, genKey, dataFeatures, generation, time.Duration(0)).
					Return(true, nil)
			},
			input:            features,
			expectedAccepted: true,
		},
		{
			desc: "rejected as stale is not an error",
			setup: func(tf *featuresCache) {
				tf.cache.(*cachemock.MockMultiGetCache).EXPECT().
					PutIfNewer(key, genKey, dataFeatures, generation, time.Duration(0)).
					Return(false, nil)
			},
			input:            features,
			expectedAccepted: false,
		},
		{
			desc: "ttl is forwarded",
			setup: func(tf *featuresCache) {
				tf.cache.(*cachemock.MockMultiGetCache).EXPECT().
					PutIfNewer(key, genKey, dataFeatures, generation, time.Minute).
					Return(true, nil)
			},
			input:            features,
			ttl:              time.Minute,
			expectedAccepted: true,
		},
		{
			desc: "backend error is returned",
			setup: func(tf *featuresCache) {
				tf.cache.(*cachemock.MockMultiGetCache).EXPECT().
					PutIfNewer(key, genKey, dataFeatures, generation, time.Duration(0)).
					Return(false, backendErr)
			},
			input:       features,
			expectedErr: backendErr,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			tf := newFeaturesCache(t, mockController)
			tf.ttl = p.ttl
			if p.setup != nil {
				p.setup(tf)
			}
			accepted, err := tf.PutIfNewer(p.input, environmentId, generation)
			assert.Equal(t, p.expectedAccepted, accepted)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

// proto3 rejects invalid UTF-8 in string fields at marshal time; both write
// paths must surface that error without touching the backend.
func TestFeaturesMarshalError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	backing := cachemock.NewMockMultiGetCache(ctrl) // no expectations: backend must not be called
	fc := NewFeaturesCache(backing, 0)
	invalid := &featureproto.Features{Features: []*featureproto.Feature{{Id: "\xff\xfe"}}}

	err := fc.Put(invalid, environmentId)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UTF-8")

	accepted, err := fc.PutIfNewer(invalid, environmentId, 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UTF-8")
	assert.False(t, accepted)
}

func TestEvictFeatures(t *testing.T) {
	t.Parallel()
	key := fmt.Sprintf("%s:%s", environmentId, featuresKind)
	genKey := fmt.Sprintf("{%s:%s}:%s", environmentId, featuresKind, featuresGenerationSuffix)
	deleteErr := errors.New("delete failed")

	t.Run("in-memory backend removes both keys", func(t *testing.T) {
		t.Parallel()
		backing := NewInMemoryCache()
		defer backing.Destroy()
		fc := NewFeaturesCache(backing, 0)
		accepted, err := fc.PutIfNewer(createFeatures(t), environmentId, 1)
		require.NoError(t, err)
		require.True(t, accepted)
		_, err = backing.Get(key)
		require.NoError(t, err)
		_, err = backing.Get(genKey)
		require.NoError(t, err)

		require.NoError(t, fc.Evict(environmentId))
		_, err = backing.Get(key)
		assert.Equal(t, cache.ErrNotFound, err)
		_, err = backing.Get(genKey)
		assert.Equal(t, cache.ErrNotFound, err)
	})
	t.Run("deleter backend deletes both keys", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		backing := cachemock.NewMockMultiGetDeleteCountCache(ctrl)
		backing.EXPECT().Delete(key).Return(nil)
		backing.EXPECT().Delete(genKey).Return(nil)
		fc := NewFeaturesCache(backing, 0)
		assert.NoError(t, fc.Evict(environmentId))
	})
	t.Run("value key delete failure short-circuits", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		backing := cachemock.NewMockMultiGetDeleteCountCache(ctrl)
		backing.EXPECT().Delete(key).Return(deleteErr)
		fc := NewFeaturesCache(backing, 0)
		assert.Equal(t, deleteErr, fc.Evict(environmentId))
	})
	t.Run("generation key delete failure is returned", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		backing := cachemock.NewMockMultiGetDeleteCountCache(ctrl)
		backing.EXPECT().Delete(key).Return(nil)
		backing.EXPECT().Delete(genKey).Return(deleteErr)
		fc := NewFeaturesCache(backing, 0)
		assert.Equal(t, deleteErr, fc.Evict(environmentId))
	})
	t.Run("unsupported backend returns error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		fc := NewFeaturesCache(cachemock.NewMockMultiGetCache(ctrl), 0)
		err := fc.Evict(environmentId)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported cache backend")
	})
}

// The Lua script touches both keys in one call, which Redis Cluster only
// allows when they hash to the same slot. The value key format cannot change
// (existing readers), so the generation key must adopt its slot via a hash
// tag.
func TestFeaturesGenerationKeySharesHashSlot(t *testing.T) {
	t.Parallel()
	fc := &featuresCache{}
	for _, envID := range []string{"", "production", "env-1", "01HZX0000000000000000000", "a:b"} {
		t.Run(fmt.Sprintf("env=%q", envID), func(t *testing.T) {
			valueKey := fc.key(envID)
			genKey := fc.generationKey(envID)
			assert.Equal(t, fmt.Sprintf("{%s}:gen", valueKey), genKey)
			assert.Equal(t, redis.HashSlot(valueKey), redis.HashSlot(genKey))
		})
	}
}

// End-to-end through the in-memory backend: the snapshot returned by Get is
// always the one with the highest generation, in either write order.
func TestFeaturesCachePutIfNewerRoundTrip(t *testing.T) {
	t.Parallel()
	newer := &featureproto.Features{Features: []*featureproto.Feature{{Id: "newer"}}}
	older := &featureproto.Features{Features: []*featureproto.Feature{{Id: "older"}}}

	t.Run("newer first", func(t *testing.T) {
		t.Parallel()
		backing := NewInMemoryCache()
		defer backing.Destroy()
		fc := NewFeaturesCache(backing, 0)
		accepted, err := fc.PutIfNewer(newer, environmentId, 2)
		require.NoError(t, err)
		assert.True(t, accepted)
		accepted, err = fc.PutIfNewer(older, environmentId, 1)
		require.NoError(t, err)
		assert.False(t, accepted)
		got, err := fc.Get(environmentId)
		require.NoError(t, err)
		assert.Equal(t, "newer", got.Features[0].Id)
	})
	t.Run("older first", func(t *testing.T) {
		t.Parallel()
		backing := NewInMemoryCache()
		defer backing.Destroy()
		fc := NewFeaturesCache(backing, 0)
		accepted, err := fc.PutIfNewer(older, environmentId, 1)
		require.NoError(t, err)
		assert.True(t, accepted)
		accepted, err = fc.PutIfNewer(newer, environmentId, 2)
		require.NoError(t, err)
		assert.True(t, accepted)
		got, err := fc.Get(environmentId)
		require.NoError(t, err)
		assert.Equal(t, "newer", got.Features[0].Id)
	})
	t.Run("unconditional Put then older conditional put is accepted", func(t *testing.T) {
		t.Parallel()
		backing := NewInMemoryCache()
		defer backing.Destroy()
		fc := NewFeaturesCache(backing, 0)
		accepted, err := fc.PutIfNewer(newer, environmentId, 100)
		require.NoError(t, err)
		require.True(t, accepted)
		require.NoError(t, fc.Put(older, environmentId))
		accepted, err = fc.PutIfNewer(newer, environmentId, 1)
		require.NoError(t, err)
		assert.True(t, accepted)
	})
}

func createFeatures(t *testing.T) *featureproto.Features {
	t.Helper()
	f := []*featureproto.Feature{}
	for i := 0; i < 5; i++ {
		feature := &featureproto.Feature{
			Id:   fmt.Sprintf("feature-id-%d", i),
			Name: fmt.Sprintf("feature-name-%d", i),
		}
		f = append(f, feature)
	}
	return &featureproto.Features{
		Features: f,
	}
}

func marshalMessage(t *testing.T, pb proto.Message) interface{} {
	t.Helper()
	buffer, err := proto.Marshal(pb)
	require.NoError(t, err)
	return buffer
}

func newFeaturesCache(t *testing.T, mockController *gomock.Controller) *featuresCache {
	t.Helper()
	return &featuresCache{
		cache: cachemock.NewMockMultiGetCache(mockController),
	}
}
