// Copyright 2025 The Bucketeer Authors.
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
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewEnvironmentV2(t *testing.T) {
	t.Parallel()
	env, err := NewEnvironmentV2(
		"name",
		"code",
		"desc",
		"project-id",
		"organization-id",
		false,
		zap.NewNop(),
	)
	assert.NoError(t, err)
	assert.IsType(t, &EnvironmentV2{}, env)
	assert.Equal(t, "name", env.Name)
	assert.Equal(t, "code", env.UrlCode)
	assert.Equal(t, "desc", env.Description)
	assert.Equal(t, "project-id", env.ProjectId)
	assert.Equal(t, false, env.Archived)
	// Auto-archive default values
	assert.Equal(t, false, env.AutoArchiveEnabled)
	assert.Equal(t, int32(90), env.AutoArchiveUnusedDays)
	assert.Equal(t, true, env.AutoArchiveCheckCodeRefs)
}

func TestUpdateEnvironmentV2(t *testing.T) {
	t.Parallel()
	env, err := NewEnvironmentV2(
		"name",
		"code",
		"desc",
		"project-id",
		"organization-id",
		true,
		zap.NewNop(),
	)
	assert.NoError(t, err)

	updated, err := env.Update(
		wrapperspb.String("new-name"),
		wrapperspb.String("new-desc"),
		wrapperspb.Bool(false),
		wrapperspb.Bool(true),
		wrapperspb.Bool(true),
		wrapperspb.Int32(30),
		wrapperspb.Bool(false),
	)
	assert.NoError(t, err)
	assert.Equal(t, "new-name", updated.Name)
	assert.Equal(t, "new-desc", updated.Description)
	assert.Equal(t, false, updated.RequireComment)
	assert.Equal(t, true, updated.Archived)
	// Auto-archive settings
	assert.Equal(t, true, updated.AutoArchiveEnabled)
	assert.Equal(t, int32(30), updated.AutoArchiveUnusedDays)
	assert.Equal(t, false, updated.AutoArchiveCheckCodeRefs)
}

func TestRenameEnvironmentV2(t *testing.T) {
	t.Parallel()
	env, err := NewEnvironmentV2(
		"name",
		"code",
		"desc",
		"project-id",
		"organization-id",
		false,
		zap.NewNop(),
	)
	assert.NoError(t, err)
	newName := "new-name"
	env.Rename(newName)
	assert.Equal(t, newName, env.Name)
}

func TestChangeDescriptionEnvironmentV2(t *testing.T) {
	t.Parallel()
	env, err := NewEnvironmentV2(
		"name",
		"code",
		"desc",
		"project-id",
		"organization-id",
		false,
		zap.NewNop(),
	)
	assert.NoError(t, err)
	newDesc := "new desc"
	env.ChangeDescription(newDesc)
	assert.Equal(t, newDesc, env.Description)
}

func TestChangeRequireCommentEnvironmentV2(t *testing.T) {
	t.Parallel()
	env, err := NewEnvironmentV2(
		"name",
		"code",
		"desc",
		"project-id",
		"organization-id",
		false,
		zap.NewNop(),
	)
	assert.NoError(t, err)
	env.ChangeRequireComment(true)
	assert.Equal(t, true, env.RequireComment)
}

func TestSetArchivedEnvironmentV2(t *testing.T) {
	t.Parallel()
	env, err := NewEnvironmentV2(
		"name",
		"code",
		"desc",
		"project-id",
		"organization-id",
		false,
		zap.NewNop(),
	)
	assert.NoError(t, err)
	env.SetArchived()
	assert.Equal(t, true, env.Archived)
}
