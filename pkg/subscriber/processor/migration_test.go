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

package processor

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetEnvironmentIDMigration(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		setupFunc   func()
		expectedEnv *EnvironmentIDMigration
	}{
		{
			desc: "migration disabled by default when no env vars set",
			setupFunc: func() {
				os.Unsetenv(envMigrationEnvironmentIDEnabled)
				os.Unsetenv(envMigrationEnvironmentIDFrom)
				os.Unsetenv(envMigrationEnvironmentIDTo)
			},
			expectedEnv: &EnvironmentIDMigration{
				Enabled:           false,
				FromEnvironmentID: "",
				ToEnvironmentID:   "",
			},
		},
		{
			desc: "migration disabled when enabled is false",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "false")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "new-uuid-123")
			},
			expectedEnv: &EnvironmentIDMigration{
				Enabled:           false,
				FromEnvironmentID: "",
				ToEnvironmentID:   "new-uuid-123",
			},
		},
		{
			desc: "migration disabled when enabled is not exactly true",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "TRUE") // uppercase
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "new-uuid-123")
			},
			expectedEnv: &EnvironmentIDMigration{
				Enabled:           false,
				FromEnvironmentID: "",
				ToEnvironmentID:   "new-uuid-123",
			},
		},
		{
			desc: "migration enabled when enabled is true",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "new-uuid-123")
			},
			expectedEnv: &EnvironmentIDMigration{
				Enabled:           true,
				FromEnvironmentID: "",
				ToEnvironmentID:   "new-uuid-123",
			},
		},
		{
			desc: "migration enabled with non-empty from environment ID",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "old-env-id")
				os.Setenv(envMigrationEnvironmentIDTo, "new-uuid-456")
			},
			expectedEnv: &EnvironmentIDMigration{
				Enabled:           true,
				FromEnvironmentID: "old-env-id",
				ToEnvironmentID:   "new-uuid-456",
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			// Reset the singleton before each test
			ResetMigrationConfig()

			// Setup environment variables
			p.setupFunc()

			// Get the migration config
			config := GetEnvironmentIDMigration()

			// Assert
			assert.Equal(t, p.expectedEnv.Enabled, config.Enabled)
			assert.Equal(t, p.expectedEnv.FromEnvironmentID, config.FromEnvironmentID)
			assert.Equal(t, p.expectedEnv.ToEnvironmentID, config.ToEnvironmentID)
		})
	}

	// Cleanup
	os.Unsetenv(envMigrationEnvironmentIDEnabled)
	os.Unsetenv(envMigrationEnvironmentIDFrom)
	os.Unsetenv(envMigrationEnvironmentIDTo)
	ResetMigrationConfig()
}

func TestGetMigrationTargetEnvironmentID(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc           string
		setupFunc      func()
		inputEnvID     string
		expectedTarget string
	}{
		{
			desc: "returns empty when migration is disabled",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "false")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "new-uuid-123")
			},
			inputEnvID:     "",
			expectedTarget: "",
		},
		{
			desc: "returns empty when migration is enabled but to is not set",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "")
			},
			inputEnvID:     "",
			expectedTarget: "",
		},
		{
			desc: "returns target when migration is enabled and input matches from (empty string)",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "new-uuid-123")
			},
			inputEnvID:     "",
			expectedTarget: "new-uuid-123",
		},
		{
			desc: "returns target when migration is enabled and input matches from (non-empty)",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "old-env-id")
				os.Setenv(envMigrationEnvironmentIDTo, "new-uuid-456")
			},
			inputEnvID:     "old-env-id",
			expectedTarget: "new-uuid-456",
		},
		{
			desc: "returns empty when input does not match from",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "new-uuid-123")
			},
			inputEnvID:     "different-env-id",
			expectedTarget: "",
		},
		{
			desc: "returns empty when input does not match from (non-empty from)",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "old-env-id")
				os.Setenv(envMigrationEnvironmentIDTo, "new-uuid-456")
			},
			inputEnvID:     "another-env-id",
			expectedTarget: "",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			// Reset the singleton before each test
			ResetMigrationConfig()

			// Setup environment variables
			p.setupFunc()

			// Get the migration target
			target := GetMigrationTargetEnvironmentID(p.inputEnvID)

			// Assert
			assert.Equal(t, p.expectedTarget, target)
		})
	}

	// Cleanup
	os.Unsetenv(envMigrationEnvironmentIDEnabled)
	os.Unsetenv(envMigrationEnvironmentIDFrom)
	os.Unsetenv(envMigrationEnvironmentIDTo)
	ResetMigrationConfig()
}

func TestLogMigrationConfig(t *testing.T) {
	patterns := []struct {
		desc      string
		setupFunc func()
	}{
		{
			desc: "logs disabled when migration is disabled",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "false")
			},
		},
		{
			desc: "logs warning when enabled but to is not set",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "")
			},
		},
		{
			desc: "logs enabled when migration is properly configured",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "new-uuid-123")
			},
		},
	}

	logger, _ := zap.NewDevelopment()

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			// Reset the singleton before each test
			ResetMigrationConfig()

			// Setup environment variables
			p.setupFunc()

			// This should not panic
			assert.NotPanics(t, func() {
				LogMigrationConfig(logger)
			})
		})
	}

	// Cleanup
	os.Unsetenv(envMigrationEnvironmentIDEnabled)
	os.Unsetenv(envMigrationEnvironmentIDFrom)
	os.Unsetenv(envMigrationEnvironmentIDTo)
	ResetMigrationConfig()
}

func TestMigrationConfigSingleton(t *testing.T) {
	// Reset before test
	ResetMigrationConfig()

	// Set initial values
	os.Setenv(envMigrationEnvironmentIDEnabled, "true")
	os.Setenv(envMigrationEnvironmentIDFrom, "initial-from")
	os.Setenv(envMigrationEnvironmentIDTo, "initial-to")

	// Get config first time
	config1 := GetEnvironmentIDMigration()
	assert.Equal(t, true, config1.Enabled)
	assert.Equal(t, "initial-from", config1.FromEnvironmentID)
	assert.Equal(t, "initial-to", config1.ToEnvironmentID)

	// Change environment variables
	os.Setenv(envMigrationEnvironmentIDEnabled, "false")
	os.Setenv(envMigrationEnvironmentIDFrom, "changed-from")
	os.Setenv(envMigrationEnvironmentIDTo, "changed-to")

	// Get config second time - should return cached values (singleton behavior)
	config2 := GetEnvironmentIDMigration()
	assert.Equal(t, true, config2.Enabled, "Should return cached enabled value")
	assert.Equal(t, "initial-from", config2.FromEnvironmentID, "Should return cached from value")
	assert.Equal(t, "initial-to", config2.ToEnvironmentID, "Should return cached to value")

	// Verify it's the same instance
	assert.Same(t, config1, config2, "Should return the same instance")

	// Cleanup
	os.Unsetenv(envMigrationEnvironmentIDEnabled)
	os.Unsetenv(envMigrationEnvironmentIDFrom)
	os.Unsetenv(envMigrationEnvironmentIDTo)
	ResetMigrationConfig()
}
