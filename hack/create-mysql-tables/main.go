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
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	if err := run(ctx, logger); err != nil {
		logger.Error("Failed to set up MySQL tables", zap.Error(err))
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *zap.Logger) error {
	// Get configuration from environment variables with defaults
	host := getEnv("MYSQL_HOST", "localhost")
	port := getEnv("MYSQL_PORT", "3306")
	user := getEnv("MYSQL_USER", "root")
	password := getEnv("MYSQL_PASSWORD", "")
	schemaFile := getEnv("SCHEMA_FILE", "./create_tables.sql")

	// Create DSN for MySQL connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, password, host, port)

	logger.Info("Connecting to MySQL",
		zap.String("host", host),
		zap.String("port", port),
		zap.String("user", user),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Error("Failed to connect to MySQL", zap.Error(err))
		return err
	}
	defer db.Close()

	// Check connection
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping MySQL server", zap.Error(err))
		return err
	}
	logger.Info("Connected to MySQL server")

	// Read SQL schema file
	schemaPath, err := filepath.Abs(schemaFile)
	if err != nil {
		logger.Error("Failed to get absolute path to schema file",
			zap.Error(err),
			zap.String("path", schemaFile),
		)
		return err
	}

	logger.Info("Reading schema file", zap.String("path", schemaPath))
	schemaSQL, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		logger.Error("Failed to read schema file",
			zap.Error(err),
			zap.String("path", schemaPath),
		)
		return err
	}

	// Execute SQL schema
	logger.Info("Executing SQL schema")
	_, err = db.ExecContext(ctx, string(schemaSQL))
	if err != nil {
		logger.Error("Failed to execute SQL schema", zap.Error(err))
		return err
	}

	logger.Info("Successfully created MySQL tables")
	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
