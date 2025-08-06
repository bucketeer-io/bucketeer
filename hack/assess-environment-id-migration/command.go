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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	v3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

type command struct {
	*kingpin.CmdClause
	mysqlUser     *string
	mysqlPass     *string
	mysqlHost     *string
	mysqlPort     *int
	mysqlDBName   *string
	redisAddress  *string
	redisPassword *string
	outputFormat  *string
	generateUUIDs *bool
}

type TableAssessment struct {
	TableName         string   `json:"table_name"`
	EmptyEnvironments int64    `json:"empty_environments"`
	TotalRecords      int64    `json:"total_records"`
	SampleEmptyIDs    []string `json:"sample_empty_ids,omitempty"`
}

type RedisAssessment struct {
	AdminKeys         int64            `json:"admin_keys"`
	KeysByPattern     map[string]int64 `json:"keys_by_pattern"`
	SampleKeys        []string         `json:"sample_keys"`
	EstimatedDataSize string           `json:"estimated_data_size"`
}

type MigrationAssessment struct {
	Timestamp       time.Time         `json:"timestamp"`
	MySQLAssessment []TableAssessment `json:"mysql_assessment"`
	RedisAssessment RedisAssessment   `json:"redis_assessment"`
	GeneratedUUIDs  []string          `json:"generated_uuids,omitempty"`
	Summary         AssessmentSummary `json:"summary"`
	Recommendations []string          `json:"recommendations"`
}

type AssessmentSummary struct {
	TotalEmptyEnvironments       int64  `json:"total_empty_environments"`
	TotalAffectedMySQLTables     int    `json:"total_affected_mysql_tables"`
	TotalAffectedRedisKeys       int64  `json:"total_affected_redis_keys"`
	EstimatedMigrationTime       string `json:"estimated_migration_time"`
	RecommendedMaintenanceWindow string `json:"recommended_maintenance_window"`
}

// Tables to check for empty environment_id
var targetTables = []string{
	"account",
	"api_key",
	"audit_log",
	"auto_ops_rule",
	"experiment",
	"experiment_result",
	"feature",
	"feature_last_used_info",
	"flag_trigger",
	"goal",
	"mau",
	"ops_count",
	"ops_progressive_rollout",
	"push",
	"segment",
	"segment_user",
	"subscription",
	"tag",
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("assess", "Assess environment_id migration requirements")
	command := &command{
		CmdClause:     cmd,
		mysqlUser:     cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:     cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:     cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:     cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName:   cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		redisAddress:  cmd.Flag("redis-address", "Redis cluster address.").Required().String(),
		redisPassword: cmd.Flag("redis-password", "Redis password.").String(),
		outputFormat:  cmd.Flag("output-format", "Output format (json|text)").Default("text").String(),
		generateUUIDs: cmd.Flag("generate-uuids", "Generate UUIDs for migration.").Default("true").Bool(),
	}
	r.RegisterCommand(command)
	return command
}

func (c *command) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	logger.Info("Starting environment_id migration assessment...")

	assessment := MigrationAssessment{
		Timestamp: time.Now(),
	}

	// Step 1: MySQL Assessment
	logger.Info("Analyzing MySQL tables...")
	mysqlAssessment, err := c.assessMySQL(ctx, logger)
	if err != nil {
		logger.Error("Failed to assess MySQL", zap.Error(err))
		return err
	}
	assessment.MySQLAssessment = mysqlAssessment

	// Step 2: Redis Assessment
	logger.Info("Analyzing Redis keys...")
	redisAssessment, err := c.assessRedis(ctx, logger)
	if err != nil {
		logger.Error("Failed to assess Redis", zap.Error(err))
		return err
	}
	assessment.RedisAssessment = redisAssessment

	// Step 3: Generate UUIDs if requested
	if *c.generateUUIDs {
		logger.Info("Generating UUIDs...")
		uuids, err := c.generateAndValidateUUIDs(ctx, logger)
		if err != nil {
			logger.Error("Failed to generate UUIDs", zap.Error(err))
			return err
		}
		assessment.GeneratedUUIDs = uuids
	}

	// Step 4: Create summary and recommendations
	assessment.Summary = c.createSummary(assessment)
	assessment.Recommendations = c.createRecommendations(assessment)

	// Step 5: Output results
	return c.outputAssessment(assessment, logger)
}

func (c *command) assessMySQL(ctx context.Context, logger *zap.Logger) ([]TableAssessment, error) {
	client, err := c.createMySQLClient(ctx, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create MySQL client: %w", err)
	}
	defer client.Close()

	var assessments []TableAssessment

	for _, table := range targetTables {
		logger.Info("Assessing table", zap.String("table", table))

		assessment, err := c.assessTable(ctx, client, table, logger)
		if err != nil {
			logger.Warn("Failed to assess table", zap.String("table", table), zap.Error(err))
			continue
		}

		if assessment.EmptyEnvironments > 0 {
			assessments = append(assessments, assessment)
			logger.Info("Found empty environment_ids",
				zap.String("table", table),
				zap.Int64("count", assessment.EmptyEnvironments))
		}
	}

	return assessments, nil
}

func (c *command) assessTable(ctx context.Context, client mysql.Client, tableName string, logger *zap.Logger) (TableAssessment, error) {
	assessment := TableAssessment{
		TableName: tableName,
	}

	// Check if table has environment_id column
	columnCheckQuery := `
		SELECT COUNT(*) 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND COLUMN_NAME = 'environment_id'
	`

	var hasColumn int64
	err := client.QueryRowContext(ctx, columnCheckQuery, *c.mysqlDBName, tableName).Scan(&hasColumn)
	if err != nil {
		return assessment, fmt.Errorf("failed to check column existence: %w", err)
	}

	if hasColumn == 0 {
		logger.Debug("Table does not have environment_id column", zap.String("table", tableName))
		return assessment, nil
	}

	// Count total records
	totalQuery := fmt.Sprintf("SELECT COUNT(*) FROM `%s`", tableName)
	err = client.QueryRowContext(ctx, totalQuery).Scan(&assessment.TotalRecords)
	if err != nil {
		return assessment, fmt.Errorf("failed to count total records: %w", err)
	}

	// Count empty environment_id records
	emptyQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s 
		WHERE environment_id = ? OR environment_id IS NULL
	`, tableName)

	err = client.QueryRowContext(ctx, emptyQuery, storage.AdminEnvironmentID).Scan(&assessment.EmptyEnvironments)
	if err != nil {
		return assessment, fmt.Errorf("failed to count empty environment_ids: %w", err)
	}

	// Get sample IDs if there are empty environments (limit to 5 samples)
	if assessment.EmptyEnvironments > 0 {
		sampleQuery := fmt.Sprintf(`
			SELECT id 
			FROM %s 
			WHERE environment_id = ? OR environment_id IS NULL 
			LIMIT 5
		`, tableName)

		rows, err := client.QueryContext(ctx, sampleQuery, storage.AdminEnvironmentID)
		if err != nil {
			logger.Warn("Failed to get sample IDs", zap.String("table", tableName), zap.Error(err))
		} else {
			defer rows.Close()
			for rows.Next() {
				var id string
				if err := rows.Scan(&id); err == nil {
					assessment.SampleEmptyIDs = append(assessment.SampleEmptyIDs, id)
				}
			}
		}
	}

	return assessment, nil
}

func (c *command) assessRedis(ctx context.Context, logger *zap.Logger) (RedisAssessment, error) {
	client, err := c.createRedisClient(logger)
	if err != nil {
		return RedisAssessment{}, fmt.Errorf("failed to create Redis client: %w", err)
	}
	defer client.Close()

	assessment := RedisAssessment{
		KeysByPattern: make(map[string]int64),
		SampleKeys:    make([]string, 0),
	}

	// Scan for keys that match admin pattern (no environment prefix)
	// Admin keys follow pattern: kind:id:... instead of environment_id:kind:id:...
	var cursor uint64
	var totalAdminKeys int64
	batchSize := int64(1000)
	_ = regexp.MustCompile(`^[^:]+:[^:]+:`)

	logger.Info("Scanning Redis keys for admin patterns...")

	for {
		nextCursor, keys, err := client.Scan(cursor, "*", batchSize)
		if err != nil {
			return assessment, fmt.Errorf("failed to scan Redis keys: %w", err)
		}

		for _, key := range keys {
			// Check if this looks like an admin key (no environment_id prefix)
			// Admin keys: kind:timestamp:feature_id:variation_id
			// Regular keys: environment_id:kind:timestamp:feature_id:variation_id
			parts := strings.Split(key, ":")
			if len(parts) >= 3 {
				// Check if first part looks like a kind (not UUID/environment_id)
				if c.isAdminKey(key) {
					totalAdminKeys++
					kind := parts[0]
					assessment.KeysByPattern[kind]++

					// Store sample keys (limit to 10)
					if len(assessment.SampleKeys) < 10 {
						assessment.SampleKeys = append(assessment.SampleKeys, key)
					}
				}
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}

		// Log progress every 10000 keys
		if totalAdminKeys%10000 == 0 && totalAdminKeys > 0 {
			logger.Info("Scanning progress", zap.Int64("admin_keys_found", totalAdminKeys))
		}
	}

	assessment.AdminKeys = totalAdminKeys
	assessment.EstimatedDataSize = c.estimateRedisDataSize(totalAdminKeys)

	logger.Info("Redis assessment completed",
		zap.Int64("total_admin_keys", totalAdminKeys),
		zap.Int("pattern_count", len(assessment.KeysByPattern)))

	return assessment, nil
}

func (c *command) isAdminKey(key string) bool {
	parts := strings.Split(key, ":")
	if len(parts) < 2 {
		return false
	}

	// Known admin key patterns based on cache kinds
	adminKinds := []string{
		"uc",  // user count
		"ec",  // evaluation count
		"gc",  // goal count
		"seg", // segment
		"ftr", // feature
	}

	firstPart := parts[0]
	for _, kind := range adminKinds {
		if firstPart == kind {
			return true
		}
	}

	// Also check if first part is not a UUID (environment_id)
	if _, err := uuid.Parse(firstPart); err != nil {
		// Not a UUID, likely an admin key if it matches expected patterns
		return len(parts) >= 3 && c.looksLikeTimestamp(parts[1])
	}

	return false
}

func (c *command) looksLikeTimestamp(s string) bool {
	// Check if string looks like a Unix timestamp (10 digits)
	if len(s) == 10 {
		for _, char := range s {
			if char < '0' || char > '9' {
				return false
			}
		}
		return true
	}
	return false
}

func (c *command) estimateRedisDataSize(keyCount int64) string {
	// Rough estimate: average key size ~100 bytes, value size ~200 bytes
	estimatedBytes := keyCount * 300

	if estimatedBytes < 1024 {
		return fmt.Sprintf("%d bytes", estimatedBytes)
	} else if estimatedBytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(estimatedBytes)/1024)
	} else if estimatedBytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(estimatedBytes)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float64(estimatedBytes)/(1024*1024*1024))
	}
}

func (c *command) generateAndValidateUUIDs(ctx context.Context, logger *zap.Logger) ([]string, error) {
	// Count how many empty environments we found
	emptyEnvCount := 0
	// For now, generate 1 UUID assuming most cases have 1 admin environment
	// This can be adjusted based on actual findings
	emptyEnvCount = 1

	client, err := c.createMySQLClient(ctx, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create MySQL client: %w", err)
	}
	defer client.Close()

	var generatedUUIDs []string

	for i := 0; i < emptyEnvCount; i++ {
		var newUUID string
		var attempts int
		maxAttempts := 10

		for attempts < maxAttempts {
			newUUID = uuid.New().String()

			// Check if UUID already exists in environment_v2 table
			var exists int
			checkQuery := "SELECT COUNT(*) FROM environment_v2 WHERE id = ?"
			err := client.QueryRowContext(ctx, checkQuery, newUUID).Scan(&exists)
			if err != nil {
				logger.Warn("Failed to check UUID uniqueness", zap.Error(err))
			}

			if exists == 0 {
				break // UUID is unique
			}

			attempts++
			logger.Warn("UUID collision detected, generating new one",
				zap.String("uuid", newUUID),
				zap.Int("attempt", attempts))
		}

		if attempts >= maxAttempts {
			return nil, fmt.Errorf("failed to generate unique UUID after %d attempts", maxAttempts)
		}

		generatedUUIDs = append(generatedUUIDs, newUUID)
		logger.Info("Generated UUID for migration", zap.String("uuid", newUUID))
	}

	return generatedUUIDs, nil
}

func (c *command) createSummary(assessment MigrationAssessment) AssessmentSummary {
	var totalEmpty int64
	affectedTables := len(assessment.MySQLAssessment)

	for _, table := range assessment.MySQLAssessment {
		totalEmpty += table.EmptyEnvironments
	}

	// Estimate migration time based on data volume
	estimatedMinutes := 30 // Base time

	// Add time for MySQL migration (1 minute per 100k records)
	if totalEmpty > 0 {
		estimatedMinutes += int(totalEmpty / 100000)
	}

	// Add time for Redis migration (1 minute per 10k keys)
	if assessment.RedisAssessment.AdminKeys > 0 {
		estimatedMinutes += int(assessment.RedisAssessment.AdminKeys / 10000)
	}

	estimatedTime := fmt.Sprintf("%d-%.0f minutes", estimatedMinutes, float64(estimatedMinutes)*1.5)

	return AssessmentSummary{
		TotalEmptyEnvironments:       totalEmpty,
		TotalAffectedMySQLTables:     affectedTables,
		TotalAffectedRedisKeys:       assessment.RedisAssessment.AdminKeys,
		EstimatedMigrationTime:       estimatedTime,
		RecommendedMaintenanceWindow: "1-2 hours",
	}
}

func (c *command) createRecommendations(assessment MigrationAssessment) []string {
	var recommendations []string

	if assessment.Summary.TotalEmptyEnvironments == 0 && assessment.Summary.TotalAffectedRedisKeys == 0 {
		recommendations = append(recommendations, "✅ No migration needed - no empty environment_ids found")
		return recommendations
	}

	recommendations = append(recommendations, "🚀 Migration Strategy Recommendations:")

	if assessment.Summary.TotalEmptyEnvironments > 100000 {
		recommendations = append(recommendations, "⚠️  Large dataset detected - consider batch processing for MySQL updates")
	}

	if assessment.Summary.TotalAffectedRedisKeys > 50000 {
		recommendations = append(recommendations, "⚠️  Large Redis dataset - implement progress monitoring")
	}

	recommendations = append(recommendations, "📋 Pre-migration checklist:")
	recommendations = append(recommendations, "   - Stop admin console and subscriber services")
	recommendations = append(recommendations, "   - Create database backup")
	recommendations = append(recommendations, "   - Test migration script on staging environment")

	recommendations = append(recommendations, "🔧 Implementation order:")
	recommendations = append(recommendations, "   1. Generate and validate UUIDs")
	recommendations = append(recommendations, "   2. Update MySQL tables in transaction batches")
	recommendations = append(recommendations, "   3. Migrate Redis keys with progress tracking")
	recommendations = append(recommendations, "   4. Validate data integrity")
	recommendations = append(recommendations, "   5. Restart services with environment_id validation")

	if assessment.Summary.TotalAffectedRedisKeys > 0 {
		recommendations = append(recommendations, "🔄 Redis migration approach:")
		recommendations = append(recommendations, "   - Use SCAN instead of KEYS for cluster safety")
		recommendations = append(recommendations, "   - Process in batches of 1000 keys")
		recommendations = append(recommendations, "   - Keep original keys as backup during migration")
	}

	return recommendations
}

func (c *command) outputAssessment(assessment MigrationAssessment, logger *zap.Logger) error {
	switch *c.outputFormat {
	case "json":
		output, err := json.MarshalIndent(assessment, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(output))

	case "text":
		c.printTextReport(assessment)

	default:
		return fmt.Errorf("unsupported output format: %s", *c.outputFormat)
	}

	return nil
}

func (c *command) printTextReport(assessment MigrationAssessment) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("🔍 ENVIRONMENT_ID MIGRATION ASSESSMENT REPORT")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("📅 Generated: %s\n\n", assessment.Timestamp.Format("2006-01-02 15:04:05"))

	// Summary
	fmt.Println("📊 SUMMARY")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Printf("Empty Environments Found: %d\n", assessment.Summary.TotalEmptyEnvironments)
	fmt.Printf("Affected MySQL Tables: %d\n", assessment.Summary.TotalAffectedMySQLTables)
	fmt.Printf("Affected Redis Keys: %d\n", assessment.Summary.TotalAffectedRedisKeys)
	fmt.Printf("Estimated Migration Time: %s\n", assessment.Summary.EstimatedMigrationTime)
	fmt.Printf("Recommended Window: %s\n\n", assessment.Summary.RecommendedMaintenanceWindow)

	// MySQL Assessment
	if len(assessment.MySQLAssessment) > 0 {
		fmt.Println("🗄️  MYSQL ASSESSMENT")
		fmt.Println(strings.Repeat("-", 30))
		for _, table := range assessment.MySQLAssessment {
			fmt.Printf("Table: %s\n", table.TableName)
			fmt.Printf("  ├─ Empty environment_ids: %d\n", table.EmptyEnvironments)
			fmt.Printf("  ├─ Total records: %d\n", table.TotalRecords)
			if len(table.SampleEmptyIDs) > 0 {
				fmt.Printf("  └─ Sample IDs: %v\n", table.SampleEmptyIDs)
			}
			fmt.Println()
		}
	}

	// Redis Assessment
	if assessment.RedisAssessment.AdminKeys > 0 {
		fmt.Println("🔴 REDIS ASSESSMENT")
		fmt.Println(strings.Repeat("-", 30))
		fmt.Printf("Total Admin Keys: %d\n", assessment.RedisAssessment.AdminKeys)
		fmt.Printf("Estimated Data Size: %s\n", assessment.RedisAssessment.EstimatedDataSize)

		if len(assessment.RedisAssessment.KeysByPattern) > 0 {
			fmt.Println("Keys by Pattern:")
			for pattern, count := range assessment.RedisAssessment.KeysByPattern {
				fmt.Printf("  ├─ %s: %d keys\n", pattern, count)
			}
		}

		if len(assessment.RedisAssessment.SampleKeys) > 0 {
			fmt.Printf("Sample Keys: %v\n", assessment.RedisAssessment.SampleKeys)
		}
		fmt.Println()
	}

	// Generated UUIDs
	if len(assessment.GeneratedUUIDs) > 0 {
		fmt.Println("🆔 GENERATED UUIDS")
		fmt.Println(strings.Repeat("-", 30))
		for i, uuid := range assessment.GeneratedUUIDs {
			fmt.Printf("UUID %d: %s\n", i+1, uuid)
		}
		fmt.Println()
	}

	// Recommendations
	fmt.Println("💡 RECOMMENDATIONS")
	fmt.Println(strings.Repeat("-", 30))
	for _, rec := range assessment.Recommendations {
		fmt.Println(rec)
	}

	fmt.Println("\n" + strings.Repeat("=", 70))
}

func (c *command) createMySQLClient(ctx context.Context, logger *zap.Logger) (mysql.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return mysql.NewClient(
		ctx,
		*c.mysqlUser, *c.mysqlPass, *c.mysqlHost,
		*c.mysqlPort,
		*c.mysqlDBName,
		mysql.WithLogger(logger),
	)
}

func (c *command) createRedisClient(logger *zap.Logger) (v3.Client, error) {
	opts := []v3.Option{v3.WithLogger(logger)}
	if *c.redisPassword != "" {
		opts = append(opts, v3.WithPassword(*c.redisPassword))
	}
	return v3.NewClient(*c.redisAddress, opts...)
}
