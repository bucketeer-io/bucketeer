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

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestRenameGoal(t *testing.T) {
	t.Parallel()
	g := newGoal(t)
	newName := "newGName"
	err := g.Rename(newName)
	assert.NoError(t, err)
	assert.Equal(t, newName, g.Name)
}

func TestChangeDescriptionGoal(t *testing.T) {
	t.Parallel()
	g := newGoal(t)
	newDesc := "newGDesc"
	err := g.ChangeDescription(newDesc)
	assert.NoError(t, err)
	assert.Equal(t, newDesc, g.Description)
}

func TestSetArchivedGoal(t *testing.T) {
	t.Parallel()
	g := newGoal(t)
	err := g.SetArchived()
	assert.NoError(t, err)
	assert.True(t, g.Archived)
}

func TestSetDeletedGoal(t *testing.T) {
	t.Parallel()
	g := newGoal(t)
	err := g.SetDeleted()
	assert.NoError(t, err)
	assert.True(t, g.Deleted)
}

func newGoal(t *testing.T) *Goal {
	t.Helper()
	g, err := NewGoal("gID", "gName", "gDesc")
	require.NoError(t, err)
	return g
}
