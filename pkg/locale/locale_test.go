// Copyright 2024 The Bucketeer Authors.
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

package locale

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetLocation(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		input    string
		expected *time.Location
		invalid  bool
	}{
		{
			desc:     "tokyo",
			input:    "Asia/Tokyo",
			expected: time.FixedZone("Asia/Tokyo", 9*60*60),
			invalid:  false,
		},
		{
			desc:     "UTC",
			input:    "UTC",
			expected: time.FixedZone("UTC", 0),
			invalid:  false,
		},
		{
			desc:     "invalid",
			input:    "invalid",
			expected: nil,
			invalid:  true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := GetLocation(p.input)
			assert.Equal(t, actual, p.expected)
			if p.invalid {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
