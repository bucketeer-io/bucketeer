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

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	deproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

func TestExtractAPIKeySecrets(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc        string
		event       *deproto.Event
		expected    []string
		expectError bool
	}{
		{
			desc: "entity data only",
			event: &deproto.Event{
				EntityData: `{"api_key": "secret-1"}`,
			},
			expected: []string{"secret-1"},
		},
		{
			desc: "previous entity data only",
			event: &deproto.Event{
				PreviousEntityData: `{"api_key": "old-secret"}`,
			},
			expected: []string{"old-secret"},
		},
		{
			desc: "both with different secrets",
			event: &deproto.Event{
				EntityData:         `{"api_key": "new-secret"}`,
				PreviousEntityData: `{"api_key": "old-secret"}`,
			},
			expected: []string{"old-secret", "new-secret"},
		},
		{
			desc: "both with same secret deduplicates",
			event: &deproto.Event{
				EntityData:         `{"api_key": "same-secret"}`,
				PreviousEntityData: `{"api_key": "same-secret"}`,
			},
			expected: []string{"same-secret"},
		},
		{
			desc:     "empty entity data returns nil without error",
			event:    &deproto.Event{},
			expected: nil,
		},
		{
			desc: "invalid JSON returns error",
			event: &deproto.Event{
				EntityData: `not json`,
			},
			expectError: true,
		},
		{
			desc: "valid JSON but missing api_key field returns empty",
			event: &deproto.Event{
				EntityData: `{"name": "test"}`,
			},
			expected: nil,
		},
		{
			desc: "one valid and one invalid snapshot returns secret with error",
			event: &deproto.Event{
				EntityData:         `{"api_key": "good-secret"}`,
				PreviousEntityData: `not json`,
			},
			expected:    []string{"good-secret"},
			expectError: true,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			result, err := ExtractAPIKeySecrets(p.event)
			if p.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, p.expected, result)
		})
	}
}
