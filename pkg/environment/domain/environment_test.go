// Copyright 2023 The Bucketeer Authors.
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
	"go.uber.org/zap"
)

func TestNewEnvironmentV2(t *testing.T) {
	t.Parallel()
	env, err := NewEnvironmentV2("name", "code", "desc", "project-id", zap.NewNop())
	assert.NoError(t, err)
	assert.IsType(t, &EnvironmentV2{}, env)
	assert.Equal(t, "name", env.Name)
	assert.Equal(t, "code", env.UrlCode)
	assert.Equal(t, "desc", env.Description)
	assert.Equal(t, "project-id", env.ProjectId)
	assert.Equal(t, false, env.Archived)
}

func TestRenameEnvironmentV2(t *testing.T) {
	t.Parallel()
	env, err := NewEnvironmentV2("name", "code", "desc", "project-id", zap.NewNop())
	assert.NoError(t, err)
	newName := "new-name"
	env.Rename(newName)
	assert.Equal(t, newName, env.Name)
}

func TestChangeDescriptionEnvironmentV2(t *testing.T) {
	t.Parallel()
	env, err := NewEnvironmentV2("name", "code", "desc", "project-id", zap.NewNop())
	assert.NoError(t, err)
	newDesc := "new desc"
	env.ChangeDescription(newDesc)
	assert.Equal(t, newDesc, env.Description)
}

func TestSetArchivedEnvironmentV2(t *testing.T) {
	t.Parallel()
	env, err := NewEnvironmentV2("name", "code", "desc", "project-id", zap.NewNop())
	assert.NoError(t, err)
	env.SetArchived()
	assert.Equal(t, true, env.Archived)
}
