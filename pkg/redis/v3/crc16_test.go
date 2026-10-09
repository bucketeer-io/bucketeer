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
)

func TestHashSlot(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		key      string
		expected int
	}{
		// Reference values from the Redis cluster spec / redis-cli CLUSTER KEYSLOT.
		{desc: "empty key", key: "", expected: 0},
		{desc: "foo", key: "foo", expected: 12182},
		{desc: "bar", key: "bar", expected: 5061},
		{desc: "hash tag uses only tagged part", key: "{foo}bar", expected: 12182},
		{desc: "hash tag anywhere in key", key: "prefix:{foo}:suffix", expected: 12182},
		{desc: "first tag wins", key: "{foo}{bar}", expected: 12182},
		{desc: "empty tag hashes whole key", key: "{}foo", expected: HashSlot("{}foo")},
		{desc: "unclosed brace hashes whole key", key: "{foo", expected: HashSlot("{foo")},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expected, HashSlot(p.key))
			assert.Equal(t, keyHashSlot(p.key), HashSlot(p.key))
		})
	}
	assert.NotEqual(t, HashSlot("{}foo"), HashSlot("foo"), "empty tag must not be treated as tag")
	assert.NotEqual(t, HashSlot("{foo"), HashSlot("foo"), "unclosed brace must not be treated as tag")
	assert.Less(t, HashSlot("anything"), RedisClusterSlots)
}
