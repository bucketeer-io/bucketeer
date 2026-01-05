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
)

func TestNewTeam(t *testing.T) {
	t.Parallel()
	team, err := NewTeam("Test Team", "This is a test team", "org-123")
	assert.Nil(t, err)
	assert.Equal(t, "Test Team", team.Name)
	assert.Equal(t, "This is a test team", team.Description)
	assert.Equal(t, "org-123", team.OrganizationId)
	assert.True(t, team.CreatedAt > 0)
	assert.True(t, team.UpdatedAt > 0)
	assert.NotEmpty(t, team.Id)
}
