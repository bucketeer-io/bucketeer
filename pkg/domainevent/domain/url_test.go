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

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	proto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func TestURL(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc            string
		inputEntityType proto.Event_EntityType
		expected        string
	}{
		{
			desc:            "feature",
			inputEntityType: proto.Event_FEATURE,
			expected:        "url/env/features/id",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := URL(p.inputEntityType, "url", "env", "id")
			assert.NoError(t, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

// TestImplementedURL checks if every domain entity type has a url.
func TestImplementedURL(t *testing.T) {
	t.Parallel()
	for k, v := range proto.Event_EntityType_name {
		t.Run(v, func(t *testing.T) {
			_, err := URL(proto.Event_EntityType(k), "url", "env", "id")
			assert.NoError(t, err)
		})
	}
}
