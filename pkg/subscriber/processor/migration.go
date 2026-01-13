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
)

const (
	// Environment variable names for migration configuration
	// These are set via Helm chart and read by the subscriber
	envMigrationEnvironmentIDEnabled = "BUCKETEER_SUBSCRIBER_MIGRATION_ENVIRONMENT_ID_ENABLED"
	envMigrationEnvironmentIDFrom    = "BUCKETEER_SUBSCRIBER_MIGRATION_ENVIRONMENT_ID_FROM"
	envMigrationEnvironmentIDTo      = "BUCKETEER_SUBSCRIBER_MIGRATION_ENVIRONMENT_ID_TO"
)

// EnvironmentIDMigration holds the configuration for migrating
// Redis keys from an old environment ID to a new one.
// This is used during the migration phase where we need to
// double-write to both old and new key formats.
type EnvironmentIDMigration struct {
	// Enabled indicates if migration is active (must be explicitly set to true)
	Enabled bool
	// FromEnvironmentID is the old environment ID (can be empty string "")
	FromEnvironmentID string
	// ToEnvironmentID is the new environment ID (UUID)
	ToEnvironmentID string
}

var (
	migrationConfig     *EnvironmentIDMigration
	migrationConfigOnce sync.Once
)

// GetEnvironmentIDMigration returns the migration configuration.
// It reads from environment variables on first call and caches the result.
func GetEnvironmentIDMigration() *EnvironmentIDMigration {
	migrationConfigOnce.Do(func() {
		// Migration must be explicitly enabled via the ENABLED flag
		enabled := os.Getenv(envMigrationEnvironmentIDEnabled) == "true"
		fromEnvID := os.Getenv(envMigrationEnvironmentIDFrom)
		toEnvID := os.Getenv(envMigrationEnvironmentIDTo)

		migrationConfig = &EnvironmentIDMigration{
			Enabled:           enabled,
			FromEnvironmentID: fromEnvID,
			ToEnvironmentID:   toEnvID,
		}
	})
	return migrationConfig
}

// GetMigrationTargetEnvironmentID returns the target environment ID
// if the given environmentID should be migrated.
// Returns empty string if no migration is needed.
func GetMigrationTargetEnvironmentID(environmentID string) string {
	config := GetEnvironmentIDMigration()
	if !config.Enabled {
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
	config := GetEnvironmentIDMigration()
	if config.Enabled {
		if config.ToEnvironmentID == "" {
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

// ResetMigrationConfig resets the migration configuration.
// This is only used for testing purposes.
func ResetMigrationConfig() {
	migrationConfigOnce = sync.Once{}
	migrationConfig = nil
}
