// Copyright 2022 The Bucketeer Authors.
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

package bigtable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKey(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc                                    string
		key, environmentNamespace, columnFamily string
		expected                                string
	}{
		{
			desc:                 "Valid without environmentNamespace",
			key:                  "user-id#tag",
			environmentNamespace: "",
			expected:             "default#user-id#tag",
		},
		{
			desc:                 "Valid with environmentNamespace",
			key:                  "user-id#tag",
			environmentNamespace: "environmentNamespace",
			expected:             "environmentNamespace#user-id#tag",
		},
	}
	for _, p := range patterns {
		assert.Equal(t, p.expected, NewKey(p.environmentNamespace, p.key))
	}
}
