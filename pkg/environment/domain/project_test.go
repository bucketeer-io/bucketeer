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

func TestNewProject(t *testing.T) {
	t.Parallel()
	project := NewProject("project-id", "project desc", "test@example.com", false)
	assert.IsType(t, &Project{}, project)
}

func TestChangeDescriptionProject(t *testing.T) {
	t.Parallel()
	project := NewProject("project-id", "project desc", "test@example.com", false)
	newDesc := "new env desc"
	project.ChangeDescription(newDesc)
	assert.Equal(t, newDesc, project.Description)
}

func TestEnableProject(t *testing.T) {
	t.Parallel()
	project := NewProject("project-id", "project desc", "test@example.com", false)
	project.Disabled = true
	project.Enable()
	assert.False(t, project.Disabled)
}

func TestDisableProject(t *testing.T) {
	t.Parallel()
	project := NewProject("project-id", "project desc", "test@example.com", false)
	project.Disable()
	assert.True(t, project.Disabled)
}

func TestConvertTrialProject(t *testing.T) {
	t.Parallel()
	project := NewProject("project-id", "project desc", "test@example.com", true)
	project.ConvertTrial()
	assert.False(t, project.Trial)
}
