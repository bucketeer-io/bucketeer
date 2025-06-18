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
	_ "embed"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

var (
	//go:embed create_events_table.sql
	eventTablesMigrationSQL string
)

type command struct {
	*kingpin.CmdClause
	mysqlUser   *string
	mysqlPass   *string
	mysqlHost   *string
	mysqlPort   *int
	mysqlDBName *string
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("create", "Create MySQL event tables for data warehouse")
	command := &command{
		CmdClause:   cmd,
		mysqlUser:   cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:   cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:   cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:   cmd.Flag("mysql-port", "MySQL port.").Default("3306").Int(),
		mysqlDBName: cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
	}
	r.RegisterCommand(command)
	return command
}

func (c *command) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	logger.Info("Starting MySQL event tables creation",
		zap.String("host", *c.mysqlHost),
		zap.Int("port", *c.mysqlPort),
		zap.String("database", *c.mysqlDBName))

	// Create MySQL client
	client, err := c.createMySQLClient(ctx, logger)
	if err != nil {
		logger.Error("Failed to create MySQL client", zap.Error(err))
		return err
	}
	defer client.Close()

	// Check if tables already exist
	existingTables, err := c.checkTablesExist(ctx, client, logger)
	if err != nil {
		logger.Error("Failed to check table existence", zap.Error(err))
		return err
	}

	if len(existingTables) > 0 {
		logger.Warn("Some tables already exist, skipping creation for existing tables",
			zap.Strings("existing_tables", existingTables))
		for _, table := range existingTables {
			logger.Info("Table already exists", zap.String("table", table))
		}

		// If all tables exist, no need to proceed
		if len(existingTables) == 2 { // evaluation_event and goal_event
			logger.Info("All required tables already exist, nothing to create")
			return nil
		}
	}

	// Execute the SQL statements (will skip existing ones)
	err = c.executeSQLStatements(ctx, client, eventTablesMigrationSQL, existingTables, logger)
	if err != nil {
		logger.Error("Failed to execute SQL statements", zap.Error(err))
		return err
	}

	logger.Info("Successfully completed MySQL event tables setup")
	return nil
}

func (c *command) checkTablesExist(ctx context.Context, client mysql.Client, logger *zap.Logger) ([]string, error) {
	tables := []string{"evaluation_event", "goal_event"}
	var existingTables []string

	for _, table := range tables {
		query := `SELECT COUNT(*) FROM information_schema.tables 
				 WHERE table_schema = ? AND table_name = ?`

		var count int
		err := client.QueryRowContext(ctx, query, *c.mysqlDBName, table).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("failed to check table %s existence: %w", table, err)
		}

		if count > 0 {
			existingTables = append(existingTables, table)
		}
	}

	return existingTables, nil
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

func (c *command) executeSQLStatements(ctx context.Context, client mysql.Client, sqlContent string, existingTables []string, logger *zap.Logger) error {
	// Split the SQL content into individual statements
	statements := c.splitSQLStatements(sqlContent)

	for i, statement := range statements {
		if statement == "" {
			continue
		}

		// Skip statements for tables that already exist
		if c.shouldSkipStatement(statement, existingTables) {
			logger.Info("Skipping statement for existing table", zap.Int("statement_number", i+1))
			continue
		}

		logger.Info("Executing SQL statement", zap.Int("statement_number", i+1))
		logger.Debug("SQL statement", zap.String("sql", statement))

		_, err := client.ExecContext(ctx, statement)
		if err != nil {
			return fmt.Errorf("failed to execute statement %d: %w", i+1, err)
		}
	}

	return nil
}

func (c *command) shouldSkipStatement(statement string, existingTables []string) bool {
	statement = strings.ToLower(strings.TrimSpace(statement))

	// Check if this is a CREATE TABLE statement for an existing table
	if strings.HasPrefix(statement, "create table") {
		for _, table := range existingTables {
			if strings.Contains(statement, "`"+table+"`") || strings.Contains(statement, table) {
				return true
			}
		}
	}

	return false
}

func (c *command) splitSQLStatements(sqlContent string) []string {
	// Remove comments and split by semicolons
	lines := strings.Split(sqlContent, "\n")
	var cleanLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		cleanLines = append(cleanLines, line)
	}

	// Join lines and split by semicolons
	cleanSQL := strings.Join(cleanLines, " ")
	statements := strings.Split(cleanSQL, ";")

	var result []string
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, stmt)
		}
	}

	return result
}
