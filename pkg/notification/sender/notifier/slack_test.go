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

package notifier

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLastDays(t *testing.T) {
	patterns := []struct {
		desc     string
		inputNow time.Time
		expected int
	}{
		{
			desc:     "now is after stopAt",
			inputNow: time.Date(2019, 12, 26, 00, 00, 00, 0, time.UTC),
			expected: 0,
		},
		{
			desc:     "now equals to stopAt",
			inputNow: time.Date(2019, 12, 25, 23, 59, 59, 0, time.UTC),
			expected: 0,
		},
		{
			desc:     "0",
			inputNow: time.Date(2019, 12, 25, 23, 00, 00, 0, time.UTC),
			expected: 0,
		},
		{
			desc:     "1",
			inputNow: time.Date(2019, 12, 24, 00, 00, 00, 0, time.UTC),
			expected: 1,
		},
	}
	stopAt := time.Date(2019, 12, 25, 23, 59, 59, 0, time.UTC)
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := lastDays(p.inputNow, stopAt)
			assert.Equal(t, p.expected, actual)
		})
	}
}
