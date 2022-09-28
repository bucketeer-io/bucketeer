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

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvironment(t *testing.T) {
	t.Parallel()
	env := NewEnvironment("env-id", "env desc", "project-id")
	expectedNamespace := "envid"
	assert.IsType(t, &Environment{}, env)
	assert.Equal(t, expectedNamespace, env.Namespace)
}

func TestRenameEnvironment(t *testing.T) {
	t.Parallel()
	env := NewEnvironment("env-id", "env desc", "project-id")
	newName := "new-env-name"
	env.Rename(newName)
	assert.Equal(t, newName, env.Name)
}

func TestChangeDescriptionEnvironment(t *testing.T) {
	t.Parallel()
	env := NewEnvironment("env-id", "env desc", "project-id")
	newDesc := "new env desc"
	env.ChangeDescription(newDesc)
	assert.Equal(t, newDesc, env.Description)
}

func TestSetDeletedEnvironment(t *testing.T) {
	t.Parallel()
	env := NewEnvironment("env-id", "env desc", "project-id")
	env.SetDeleted()
	assert.True(t, env.Deleted)
}
