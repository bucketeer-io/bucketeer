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

package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStreamClientCreatePublisherMaxLen(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		opts     []StreamOption
		expected int64
	}{
		{desc: "default", opts: nil, expected: defaultStreamMaxLen},
		{desc: "custom", opts: []StreamOption{WithStreamMaxLen(7)}, expected: 7},
		{desc: "zero disables trimming", opts: []StreamOption{WithStreamMaxLen(0)}, expected: 0},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			c, err := NewStreamClient(context.Background(), nil, p.opts...)
			require.NoError(t, err)
			pub, err := c.CreatePublisher("topic")
			require.NoError(t, err)
			assert.Equal(t, p.expected, pub.(*StreamPublisher).maxLen)
		})
	}
}
