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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestHoldoutGroupCache(t *testing.T) {
	t.Parallel()
	mem := NewInMemoryCache()
	c := NewHoldoutGroupCache(mem, 0)

	_, err := c.Get("env")
	assert.ErrorIs(t, err, cache.ErrNotFound)

	group := &ftproto.HoldoutGroup{
		Id:            "hg",
		EnvironmentId: "env",
		Status:        ftproto.HoldoutGroup_RUNNING,
		Features: []*ftproto.HoldoutGroupFeature{{
			FeatureId:          "flag",
			HoldoutVariationId: "hold",
		}},
	}
	require.NoError(t, c.Put(group, "env"))
	got, err := c.Get("env")
	require.NoError(t, err)
	assert.Equal(t, "hg", got.Id)
	require.Len(t, got.Features, 1)
	assert.Equal(t, "hold", got.Features[0].HoldoutVariationId)

	require.NoError(t, c.Put(&ftproto.HoldoutGroup{}, "env"))
	empty, err := c.Get("env")
	require.NoError(t, err)
	assert.Empty(t, empty.Id)

	require.NoError(t, c.Evict("env"))
	_, err = c.Get("env")
	assert.ErrorIs(t, err, cache.ErrNotFound)
}

func TestHoldoutGroupCachePutIfNewer(t *testing.T) {
	t.Parallel()
	mem := NewInMemoryCache()
	c := NewHoldoutGroupCache(mem, 0)
	newer := &ftproto.HoldoutGroup{Id: "new", UpdatedAt: 2}
	older := &ftproto.HoldoutGroup{Id: "old", UpdatedAt: 1}

	accepted, err := c.PutIfNewer(newer, "env", 200)
	require.NoError(t, err)
	assert.True(t, accepted)
	accepted, err = c.PutIfNewer(older, "env", 100)
	require.NoError(t, err)
	assert.False(t, accepted)
	got, err := c.Get("env")
	require.NoError(t, err)
	assert.Equal(t, "new", got.Id)

	require.NoError(t, c.Evict("env"))
	_, err = c.Get("env")
	assert.ErrorIs(t, err, cache.ErrNotFound)
	// Evict removes the generation marker, so a later write is accepted.
	accepted, err = c.PutIfNewer(older, "env", 100)
	require.NoError(t, err)
	assert.True(t, accepted)

	accepted, err = c.PutIfNewer(nil, "env", 1)
	require.Error(t, err)
	assert.False(t, accepted)
}
