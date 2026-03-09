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

func TestKeyHashSlot(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		hashesWhat string
	}{
		{
			name:       "no hash tag",
			key:        "foo",
			hashesWhat: "foo",
		},
		{
			name:       "simple hash tag",
			key:        "{foo}",
			hashesWhat: "foo",
		},
		{
			name:       "hash tag in middle",
			key:        "foo{bar}zap",
			hashesWhat: "bar",
		},
		{
			name:       "empty hash tag hashes whole key",
			key:        "foo{}bar",
			hashesWhat: "foo{}bar",
		},
		{
			name:       "missing closing brace hashes whole key",
			key:        "foo{bar",
			hashesWhat: "foo{bar",
		},
		{
			name:       "closing brace before opening brace hashes whole key",
			key:        "foo}bar{",
			hashesWhat: "foo}bar{",
		},
		{
			name:       "nested opening brace uses first valid pair",
			key:        "foo{{bar}}zap",
			hashesWhat: "{bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := int(crc16([]byte(tt.hashesWhat)) % RedisClusterSlots)
			got := KeyHashSlot(tt.key)

			assert.Equal(t, want, got, "key=%q", tt.key)
		})
	}
}
