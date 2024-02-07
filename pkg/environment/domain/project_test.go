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
)

func TestNewProject(t *testing.T) {
	t.Parallel()
	project, err := NewProject(
		"project-name",
		"project-code",
		"project desc",
		"test@example.com",
		"organization-id",
		false,
	)
	assert.NoError(t, err)
	assert.IsType(t, &Project{}, project)
	assert.NotEqual(t, "project-name", project.Id)
	assert.Equal(t, "project-name", project.Name)
	assert.Equal(t, "project-code", project.UrlCode)
}

func TestChangeDescriptionProject(t *testing.T) {
	t.Parallel()
	project, err := NewProject(
		"project-name",
		"project-code",
		"project desc",
		"test@example.com",
		"organization-id",
		false,
	)
	assert.NoError(t, err)
	newDesc := "new env desc"
	project.ChangeDescription(newDesc)
	assert.Equal(t, newDesc, project.Description)
}

func TestRenameProject(t *testing.T) {
	t.Parallel()
	project, err := NewProject(
		"project-name",
		"project-code",
		"project desc",
		"test@example.com",
		"organization-id",
		false,
	)
	assert.NoError(t, err)
	newName := "new-project-name"
	project.Rename(newName)
	assert.Equal(t, newName, project.Name)
}

func TestEnableProject(t *testing.T) {
	t.Parallel()
	project, err := NewProject(
		"project-name",
		"project-code",
		"project desc",
		"test@example.com",
		"organization-id",
		false,
	)
	assert.NoError(t, err)
	project.Disabled = true
	project.Enable()
	assert.False(t, project.Disabled)
}

func TestDisableProject(t *testing.T) {
	t.Parallel()
	project, err := NewProject(
		"project-name",
		"project-code",
		"project desc",
		"test@example.com",
		"organization-id",
		false,
	)
	assert.NoError(t, err)
	project.Disable()
	assert.True(t, project.Disabled)
}

func TestConvertTrialProject(t *testing.T) {
	t.Parallel()
	project, err := NewProject(
		"project-name",
		"project-code",
		"project desc",
		"test@example.com",
		"organization-id",
		true,
	)
	assert.NoError(t, err)
	project.ConvertTrial()
	assert.False(t, project.Trial)
}
