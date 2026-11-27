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
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestHoldoutMeasuredGroupsCache(t *testing.T) {
	t.Parallel()
	mem := NewInMemoryCache()
	c := NewHoldoutMeasuredGroupsCache(mem)

	_, err := c.Get("env")
	assert.ErrorIs(t, err, cache.ErrNotFound)

	groups := &ftproto.HoldoutGroups{HoldoutGroups: []*ftproto.HoldoutGroup{
		{
			Id:            "stopped",
			EnvironmentId: "env",
			Status:        ftproto.HoldoutGroup_FORCE_STOPPED,
			ExperimentId:  "e1",
			Features: []*ftproto.HoldoutGroupFeature{
				{FeatureId: "flag", AddedAt: 10, RemovedAt: 20},
				{FeatureId: "flag", AddedAt: 30},
			},
		},
		{Id: "running", EnvironmentId: "env", Status: ftproto.HoldoutGroup_RUNNING, ExperimentId: "e2"},
	}}
	accepted, err := c.PutIfNewer(groups, "env", 100)
	require.NoError(t, err)
	assert.True(t, accepted)
	got, err := c.Get("env")
	require.NoError(t, err)
	assert.True(t, proto.Equal(groups, got), "the list is read back as it was written, removed targets included")

	// The list of another environment is a separate entry.
	_, err = c.Get("env-2")
	assert.ErrorIs(t, err, cache.ErrNotFound)

	// It does not share a key with the serving snapshot of the environment.
	serving := NewHoldoutGroupCache(mem, 0)
	_, err = serving.Get("env")
	assert.ErrorIs(t, err, cache.ErrNotFound)
	require.NoError(t, serving.Put(&ftproto.HoldoutGroup{Id: "running"}, "env"))
	got, err = c.Get("env")
	require.NoError(t, err)
	assert.Len(t, got.HoldoutGroups, 2)
}

func TestHoldoutMeasuredGroupsCachePutIfNewer(t *testing.T) {
	t.Parallel()
	mem := NewInMemoryCache()
	c := NewHoldoutMeasuredGroupsCache(mem)
	newer := &ftproto.HoldoutGroups{HoldoutGroups: []*ftproto.HoldoutGroup{{Id: "new"}}}
	older := &ftproto.HoldoutGroups{HoldoutGroups: []*ftproto.HoldoutGroup{{Id: "old"}}}

	accepted, err := c.PutIfNewer(newer, "env", 200)
	require.NoError(t, err)
	assert.True(t, accepted)
	accepted, err = c.PutIfNewer(older, "env", 100)
	require.NoError(t, err)
	assert.False(t, accepted, "a list read earlier does not replace a list read later")
	got, err := c.Get("env")
	require.NoError(t, err)
	require.Len(t, got.HoldoutGroups, 1)
	assert.Equal(t, "new", got.HoldoutGroups[0].Id)

	// An empty list is a value: the environment has no measured group.
	accepted, err = c.PutIfNewer(&ftproto.HoldoutGroups{}, "env", 300)
	require.NoError(t, err)
	assert.True(t, accepted)
	got, err = c.Get("env")
	require.NoError(t, err)
	assert.Empty(t, got.HoldoutGroups)

	accepted, err = c.PutIfNewer(nil, "env", 400)
	require.Error(t, err)
	assert.False(t, accepted)
}
