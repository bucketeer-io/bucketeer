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

package subscriber

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// redisStreamMaxLen must distinguish "not set" (backend default) from an
// explicit 0 (trimming disabled), so it is decoded as a pointer.
func TestConfigurationRedisStreamMaxLenJSON(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		input    string
		expected *int64
	}{
		{desc: "absent keeps nil", input: `{"redisPartitionCount":16}`, expected: nil},
		{desc: "explicit zero", input: `{"redisStreamMaxLen":0}`, expected: ptr(0)},
		{desc: "explicit value", input: `{"redisStreamMaxLen":250000}`, expected: ptr(250000)},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			var conf Configuration
			require.NoError(t, json.Unmarshal([]byte(p.input), &conf))
			assert.Equal(t, p.expected, conf.RedisStreamMaxLen)
		})
	}
}

func ptr(v int64) *int64 { return &v }
