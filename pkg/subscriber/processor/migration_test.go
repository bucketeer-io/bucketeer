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
		expectedEnv *environmentIDMigration
	}{
		{
			desc: "migration disabled by default when no env vars set",
			setupFunc: func() {
				os.Unsetenv(envMigrationEnvironmentIDEnabled)
				os.Unsetenv(envMigrationEnvironmentIDFrom)
				os.Unsetenv(envMigrationEnvironmentIDTo)
			},
			expectedEnv: &environmentIDMigration{
				Enabled:           false,
				FromEnvironmentID: "",
				ToEnvironmentID:   "",
				InvalidConfig:     false,
				InvalidReason:     "",
			},
		},
		{
			desc: "migration disabled when enabled is false (no UUID validation)",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "false")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "not-validated-when-disabled")
			},
			expectedEnv: &environmentIDMigration{
				Enabled:           false,
				FromEnvironmentID: "",
				ToEnvironmentID:   "not-validated-when-disabled",
				InvalidConfig:     false,
				InvalidReason:     "",
			},
		},
		{
			desc: "migration disabled when enabled is not exactly true",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "TRUE") // uppercase
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "not-validated-when-disabled")
			},
			expectedEnv: &environmentIDMigration{
				Enabled:           false,
				FromEnvironmentID: "",
				ToEnvironmentID:   "not-validated-when-disabled",
				InvalidConfig:     false,
				InvalidReason:     "",
			},
		},
		{
			desc: "migration enabled with valid UUID v4",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "a1b2c3d4-e5f6-4890-abcd-ef1234567890")
			},
			expectedEnv: &environmentIDMigration{
				Enabled:           true,
				FromEnvironmentID: "",
				ToEnvironmentID:   "a1b2c3d4-e5f6-4890-abcd-ef1234567890",
				InvalidConfig:     false,
				InvalidReason:     "",
			},
		},
		{
			desc: "migration enabled with non-empty from and valid UUID v4",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "old-env-id")
				os.Setenv(envMigrationEnvironmentIDTo, "12345678-1234-4123-8123-123456789012")
			},
			expectedEnv: &environmentIDMigration{
				Enabled:           true,
				FromEnvironmentID: "old-env-id",
				ToEnvironmentID:   "12345678-1234-4123-8123-123456789012",
				InvalidConfig:     false,
				InvalidReason:     "",
			},
		},
		{
			desc: "migration invalid when toEnvironmentID is not a valid UUID",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "not-a-valid-uuid")
			},
			expectedEnv: &environmentIDMigration{
				Enabled:           true,
				FromEnvironmentID: "",
				ToEnvironmentID:   "not-a-valid-uuid",
				InvalidConfig:     true,
				InvalidReason:     "toEnvironmentID must be a valid UUID v4",
			},
		},
		{
			desc: "migration invalid when toEnvironmentID is UUID v1 format",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				// UUID v1 has version 1 at position 13 (should be 4 for v4)
				os.Setenv(envMigrationEnvironmentIDTo, "a1b2c3d4-e5f6-1890-abcd-ef1234567890")
			},
			expectedEnv: &environmentIDMigration{
				Enabled:           true,
				FromEnvironmentID: "",
				ToEnvironmentID:   "a1b2c3d4-e5f6-1890-abcd-ef1234567890",
				InvalidConfig:     true,
				InvalidReason:     "toEnvironmentID must be a valid UUID v4",
			},
		},
		{
			desc: "migration valid when toEnvironmentID is empty (no validation needed)",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "")
			},
			expectedEnv: &environmentIDMigration{
				Enabled:           true,
				FromEnvironmentID: "",
				ToEnvironmentID:   "",
				InvalidConfig:     false,
				InvalidReason:     "",
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			// Reset the singleton before each test
			resetMigrationConfig()

			// Setup environment variables
			p.setupFunc()

			// Get the migration config
			config := getEnvironmentIDMigration()

			// Assert
			assert.Equal(t, p.expectedEnv.Enabled, config.Enabled)
			assert.Equal(t, p.expectedEnv.FromEnvironmentID, config.FromEnvironmentID)
			assert.Equal(t, p.expectedEnv.ToEnvironmentID, config.ToEnvironmentID)
			assert.Equal(t, p.expectedEnv.InvalidConfig, config.InvalidConfig)
			assert.Equal(t, p.expectedEnv.InvalidReason, config.InvalidReason)
		})
	}

	// Cleanup
	os.Unsetenv(envMigrationEnvironmentIDEnabled)
	os.Unsetenv(envMigrationEnvironmentIDFrom)
	os.Unsetenv(envMigrationEnvironmentIDTo)
	resetMigrationConfig()
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
				os.Setenv(envMigrationEnvironmentIDTo, "a1b2c3d4-e5f6-4890-abcd-ef1234567890")
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
				os.Setenv(envMigrationEnvironmentIDTo, "a1b2c3d4-e5f6-4890-abcd-ef1234567890")
			},
			inputEnvID:     "",
			expectedTarget: "a1b2c3d4-e5f6-4890-abcd-ef1234567890",
		},
		{
			desc: "returns target when migration is enabled and input matches from (non-empty)",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "old-env-id")
				os.Setenv(envMigrationEnvironmentIDTo, "12345678-1234-4123-8123-123456789012")
			},
			inputEnvID:     "old-env-id",
			expectedTarget: "12345678-1234-4123-8123-123456789012",
		},
		{
			desc: "returns empty when input does not match from",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "a1b2c3d4-e5f6-4890-abcd-ef1234567890")
			},
			inputEnvID:     "different-env-id",
			expectedTarget: "",
		},
		{
			desc: "returns empty when input does not match from (non-empty from)",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "old-env-id")
				os.Setenv(envMigrationEnvironmentIDTo, "12345678-1234-4123-8123-123456789012")
			},
			inputEnvID:     "another-env-id",
			expectedTarget: "",
		},
		{
			desc: "returns empty when config is invalid (invalid UUID)",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "not-a-valid-uuid")
			},
			inputEnvID:     "",
			expectedTarget: "",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			// Reset the singleton before each test
			resetMigrationConfig()

			// Setup environment variables
			p.setupFunc()

			// Get the migration target
			target := getMigrationTargetEnvironmentID(p.inputEnvID)

			// Assert
			assert.Equal(t, p.expectedTarget, target)
		})
	}

	// Cleanup
	os.Unsetenv(envMigrationEnvironmentIDEnabled)
	os.Unsetenv(envMigrationEnvironmentIDFrom)
	os.Unsetenv(envMigrationEnvironmentIDTo)
	resetMigrationConfig()
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
				os.Setenv(envMigrationEnvironmentIDTo, "a1b2c3d4-e5f6-4890-abcd-ef1234567890")
			},
		},
		{
			desc: "logs error when migration config is invalid",
			setupFunc: func() {
				os.Setenv(envMigrationEnvironmentIDEnabled, "true")
				os.Setenv(envMigrationEnvironmentIDFrom, "")
				os.Setenv(envMigrationEnvironmentIDTo, "invalid-uuid")
			},
		},
	}

	logger, _ := zap.NewDevelopment()

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			// Reset the singleton before each test
			resetMigrationConfig()

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
	resetMigrationConfig()
}

func TestMigrationConfigSingleton(t *testing.T) {
	// Reset before test
	resetMigrationConfig()

	// Set initial values with valid UUID
	os.Setenv(envMigrationEnvironmentIDEnabled, "true")
	os.Setenv(envMigrationEnvironmentIDFrom, "initial-from")
	os.Setenv(envMigrationEnvironmentIDTo, "a1b2c3d4-e5f6-4890-abcd-ef1234567890")

	// Get config first time
	config1 := getEnvironmentIDMigration()
	assert.Equal(t, true, config1.Enabled)
	assert.Equal(t, "initial-from", config1.FromEnvironmentID)
	assert.Equal(t, "a1b2c3d4-e5f6-4890-abcd-ef1234567890", config1.ToEnvironmentID)
	assert.Equal(t, false, config1.InvalidConfig)

	// Change environment variables
	os.Setenv(envMigrationEnvironmentIDEnabled, "false")
	os.Setenv(envMigrationEnvironmentIDFrom, "changed-from")
	os.Setenv(envMigrationEnvironmentIDTo, "12345678-1234-4123-8123-123456789012")

	// Get config second time - should return cached values (singleton behavior)
	config2 := getEnvironmentIDMigration()
	assert.Equal(t, true, config2.Enabled, "Should return cached enabled value")
	assert.Equal(t, "initial-from", config2.FromEnvironmentID, "Should return cached from value")
	assert.Equal(t, "a1b2c3d4-e5f6-4890-abcd-ef1234567890", config2.ToEnvironmentID, "Should return cached to value")
	assert.Equal(t, false, config2.InvalidConfig, "Should return cached invalid config value")

	// Verify it's the same instance
	assert.Same(t, config1, config2, "Should return the same instance")

	// Cleanup
	os.Unsetenv(envMigrationEnvironmentIDEnabled)
	os.Unsetenv(envMigrationEnvironmentIDFrom)
	os.Unsetenv(envMigrationEnvironmentIDTo)
	resetMigrationConfig()
}
