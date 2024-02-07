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

package uuid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateUUID(t *testing.T) {
	uuid, err := NewUUID()
	require.NoError(t, err)
	patterns := []*struct {
		id       string
		expected error
	}{
		{
			id:       "0efe416e 2fd2 4996 c5c3 194f05444f1f",
			expected: ErrIncorrectUUIDFormat,
		},
		{
			id:       "0efe416e2fd24996b5c3194f05444f1f",
			expected: ErrIncorrectUUIDFormat,
		},
		{
			id:       "0efe416e_2fd2_4996_b5c3_194f05444f1f",
			expected: ErrIncorrectUUIDFormat,
		},
		{
			id:       "0efe416e-2fd2-4996-b5c3-194f05444f1f",
			expected: nil,
		},
		{
			id:       uuid.String(),
			expected: nil,
		},
	}
	for _, p := range patterns {
		err := ValidateUUID(p.id)
		assert.Equal(t, p.expected, err)
	}
}
