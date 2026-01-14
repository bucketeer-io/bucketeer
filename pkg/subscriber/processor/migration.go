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
	"sync"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
)

const (
	// Environment variable names for migration configuration
	// These are set via Helm chart and read by the subscriber
	envMigrationEnvironmentIDEnabled = "BUCKETEER_SUBSCRIBER_MIGRATION_ENVIRONMENT_ID_ENABLED"
	envMigrationEnvironmentIDFrom    = "BUCKETEER_SUBSCRIBER_MIGRATION_ENVIRONMENT_ID_FROM"
	envMigrationEnvironmentIDTo      = "BUCKETEER_SUBSCRIBER_MIGRATION_ENVIRONMENT_ID_TO"
)

// environmentIDMigration holds the configuration for migrating
// Redis keys from an old environment ID to a new one.
// This is used during the migration phase where we need to
// double-write to both old and new key formats.
type environmentIDMigration struct {
	// Enabled indicates if migration is active (must be explicitly set to true)
	Enabled bool
	// FromEnvironmentID is the old environment ID (can be empty string "")
	FromEnvironmentID string
	// ToEnvironmentID is the new environment ID (must be a valid UUID v4)
	ToEnvironmentID string
	// InvalidConfig indicates if the configuration is invalid (e.g., invalid UUID)
	InvalidConfig bool
	// InvalidReason contains the reason why the config is invalid
	InvalidReason string
}

var (
	migrationConfig     *environmentIDMigration
	migrationConfigOnce sync.Once
)

// getEnvironmentIDMigration returns the migration configuration.
// It reads from environment variables on first call and caches the result.
// If enabled, it validates that ToEnvironmentID is a valid UUID v4.
func getEnvironmentIDMigration() *environmentIDMigration {
	migrationConfigOnce.Do(func() {
		// Migration must be explicitly enabled via the ENABLED flag
		enabled := os.Getenv(envMigrationEnvironmentIDEnabled) == "true"
		fromEnvID := os.Getenv(envMigrationEnvironmentIDFrom)
		toEnvID := os.Getenv(envMigrationEnvironmentIDTo)

		var invalidConfig bool
		var invalidReason string

		// Validate ToEnvironmentID is a valid UUID v4 when migration is enabled
		if enabled && toEnvID != "" {
			if err := uuid.ValidateUUID(toEnvID); err != nil {
				invalidConfig = true
				invalidReason = "toEnvironmentID must be a valid UUID v4"
			}
		}

		migrationConfig = &environmentIDMigration{
			Enabled:           enabled,
			FromEnvironmentID: fromEnvID,
			ToEnvironmentID:   toEnvID,
			InvalidConfig:     invalidConfig,
			InvalidReason:     invalidReason,
		}
	})
	return migrationConfig
}

// getMigrationTargetEnvironmentID returns the target environment ID
// if the given environmentID should be migrated.
// Returns empty string if no migration is needed or config is invalid.
func getMigrationTargetEnvironmentID(environmentID string) string {
	config := getEnvironmentIDMigration()
	if !config.Enabled {
		return ""
	}
	// Don't run migration if config is invalid
	if config.InvalidConfig {
		return ""
	}
	// Validate that ToEnvironmentID is set when migration is enabled
	if config.ToEnvironmentID == "" {
		return ""
	}
	// Check if this environment ID matches the one being migrated FROM
	if environmentID == config.FromEnvironmentID {
		return config.ToEnvironmentID
	}
	return ""
}

// LogMigrationConfig logs the current migration configuration
func LogMigrationConfig(logger *zap.Logger) {
	config := getEnvironmentIDMigration()
	if config.Enabled {
		if config.InvalidConfig {
			logger.Error("Environment ID migration configuration is invalid - migration will not run",
				zap.String("reason", config.InvalidReason),
				zap.String("fromEnvironmentID", config.FromEnvironmentID),
				zap.String("toEnvironmentID", config.ToEnvironmentID),
			)
		} else if config.ToEnvironmentID == "" {
			logger.Warn("Environment ID migration is enabled but toEnvironmentID is not set - migration will not run",
				zap.String("fromEnvironmentID", config.FromEnvironmentID),
			)
		} else {
			logger.Info("Environment ID migration is enabled - double-write is active",
				zap.String("fromEnvironmentID", config.FromEnvironmentID),
				zap.String("toEnvironmentID", config.ToEnvironmentID),
			)
		}
	} else {
		logger.Info("Environment ID migration is disabled")
	}
}

// resetMigrationConfig resets the migration configuration.
// This is only used for testing purposes.
func resetMigrationConfig() {
	migrationConfigOnce = sync.Once{}
	migrationConfig = nil
}
