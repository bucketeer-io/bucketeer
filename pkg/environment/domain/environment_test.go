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
	assert.Equal(t, defaultAutoArchiveUnusedDays, env.AutoArchiveUnusedDays)
	assert.Equal(t, defaultAutoArchiveCheckCodeRefs, env.AutoArchiveCheckCodeRefs)
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

func TestUpdateEnvironmentV2_AutoArchiveValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                       string
		existingAutoArchiveEnabled bool
		autoArchiveEnabled         *wrapperspb.BoolValue
		autoArchiveUnusedDays      *wrapperspb.Int32Value
		autoArchiveCheckCodeRefs   *wrapperspb.BoolValue
		expectedError              error
	}{
		{
			name:                       "err: enable auto-archive without unused_days",
			existingAutoArchiveEnabled: false,
			autoArchiveEnabled:         wrapperspb.Bool(true),
			autoArchiveUnusedDays:      nil,
			autoArchiveCheckCodeRefs:   nil,
			expectedError:              ErrAutoArchiveUnusedDaysRequired,
		},
		{
			name:                       "err: enable auto-archive with unused_days=0",
			existingAutoArchiveEnabled: false,
			autoArchiveEnabled:         wrapperspb.Bool(true),
			autoArchiveUnusedDays:      wrapperspb.Int32(0),
			autoArchiveCheckCodeRefs:   nil,
			expectedError:              ErrAutoArchiveUnusedDaysRequired,
		},
		{
			name:                       "err: enable auto-archive with unused_days negative",
			existingAutoArchiveEnabled: false,
			autoArchiveEnabled:         wrapperspb.Bool(true),
			autoArchiveUnusedDays:      wrapperspb.Int32(-1),
			autoArchiveCheckCodeRefs:   nil,
			expectedError:              ErrAutoArchiveUnusedDaysRequired,
		},
		{
			name:                       "err: update unused_days when auto-archive is disabled",
			existingAutoArchiveEnabled: false,
			autoArchiveEnabled:         nil,
			autoArchiveUnusedDays:      wrapperspb.Int32(30),
			autoArchiveCheckCodeRefs:   nil,
			expectedError:              ErrAutoArchiveNotEnabled,
		},
		{
			name:                       "err: update check_code_refs when auto-archive is disabled",
			existingAutoArchiveEnabled: false,
			autoArchiveEnabled:         nil,
			autoArchiveUnusedDays:      nil,
			autoArchiveCheckCodeRefs:   wrapperspb.Bool(false),
			expectedError:              ErrAutoArchiveNotEnabled,
		},
		{
			name:                       "err: disable auto-archive and update other fields simultaneously",
			existingAutoArchiveEnabled: true,
			autoArchiveEnabled:         wrapperspb.Bool(false),
			autoArchiveUnusedDays:      wrapperspb.Int32(60),
			autoArchiveCheckCodeRefs:   nil,
			expectedError:              ErrAutoArchiveNotEnabled,
		},
		{
			name:                       "success: enable auto-archive with valid unused_days",
			existingAutoArchiveEnabled: false,
			autoArchiveEnabled:         wrapperspb.Bool(true),
			autoArchiveUnusedDays:      wrapperspb.Int32(30),
			autoArchiveCheckCodeRefs:   nil,
			expectedError:              nil,
		},
		{
			name:                       "success: update unused_days when auto-archive is already enabled",
			existingAutoArchiveEnabled: true,
			autoArchiveEnabled:         nil,
			autoArchiveUnusedDays:      wrapperspb.Int32(60),
			autoArchiveCheckCodeRefs:   nil,
			expectedError:              nil,
		},
		{
			name:                       "success: update check_code_refs when auto-archive is already enabled",
			existingAutoArchiveEnabled: true,
			autoArchiveEnabled:         nil,
			autoArchiveUnusedDays:      nil,
			autoArchiveCheckCodeRefs:   wrapperspb.Bool(false),
			expectedError:              nil,
		},
		{
			name:                       "success: disable auto-archive only",
			existingAutoArchiveEnabled: true,
			autoArchiveEnabled:         wrapperspb.Bool(false),
			autoArchiveUnusedDays:      nil,
			autoArchiveCheckCodeRefs:   nil,
			expectedError:              nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			// Set existing auto-archive state
			env.AutoArchiveEnabled = tt.existingAutoArchiveEnabled
			if tt.existingAutoArchiveEnabled {
				env.AutoArchiveUnusedDays = 90
				env.AutoArchiveCheckCodeRefs = true
			}

			_, err = env.Update(
				nil, // name
				nil, // description
				nil, // requireComment
				nil, // archived
				tt.autoArchiveEnabled,
				tt.autoArchiveUnusedDays,
				tt.autoArchiveCheckCodeRefs,
			)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
