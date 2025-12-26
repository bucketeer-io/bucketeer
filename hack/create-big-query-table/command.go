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
	"strings"

	"cloud.google.com/go/bigquery"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cli"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
)

type command struct {
	*kingpin.CmdClause
	bigQueryEmulator *string
	project          *string
	dataset          *string
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("create", "Create table in Big Query Emulator")
	command := &command{
		CmdClause:        cmd,
		bigQueryEmulator: cmd.Flag("bigquery-emulator", "Big Query Emulator Host").Default("http://localhost:9050").String(),
		project:          cmd.Flag("project", "Project ID").Default("bucketeer-dev").String(),
		dataset:          cmd.Flag("dataset", "Dataset ID").Default("bucketeer").String(),
	}
	r.RegisterCommand(command)
	return command
}

func (c *command) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	client, err := bigquery.NewClient(ctx, *c.project,
		option.WithEndpoint(*c.bigQueryEmulator),
		option.WithoutAuthentication())
	if err != nil {
		logger.Error("failed to create BigQuery client",
			zap.Error(err))
	}
	defer client.Close()

	// Create dataset
	err = createDataset(ctx, client, logger, *c.dataset)
	if err != nil {
		if isAlreadyExistsError(err) {
			logger.Info("dataset already exists, skipping creation",
				zap.String("dataset", *c.dataset))
		} else {
			logger.Warn("failed to create dataset, continuing anyway",
				zap.Error(err),
			)
		}
	}

	// Create tables
	evaluationEventSchema := bigquery.Schema{
		{Name: "id", Type: bigquery.StringFieldType},
		{Name: "environment_id", Type: bigquery.StringFieldType},
		{Name: "timestamp", Type: bigquery.TimestampFieldType},
		{Name: "feature_id", Type: bigquery.StringFieldType},
		{Name: "feature_version", Type: bigquery.IntegerFieldType},
		{Name: "user_id", Type: bigquery.StringFieldType},
		{Name: "user_data", Type: bigquery.JSONFieldType},
		{Name: "variation_id", Type: bigquery.StringFieldType},
		{Name: "reason", Type: bigquery.StringFieldType},
		{Name: "tag", Type: bigquery.StringFieldType},
		{Name: "source_id", Type: bigquery.StringFieldType},
	}
	err = createTable(ctx, client, logger, *c.project, *c.dataset, "evaluation_event", evaluationEventSchema)
	if err != nil {
		if isAlreadyExistsError(err) {
			logger.Info("evaluation_event table already exists, skipping creation")
		} else {
			logger.Error("failed to create evaluation_event table",
				zap.Error(err),
			)
			return err
		}
	} else {
		logger.Info("successfully created evaluation_event table")
	}

	goalEventSchema := bigquery.Schema{
		{Name: "id", Type: bigquery.StringFieldType},
		{Name: "environment_id", Type: bigquery.StringFieldType},
		{Name: "timestamp", Type: bigquery.TimestampFieldType},
		{Name: "goal_id", Type: bigquery.StringFieldType},
		{Name: "value", Type: bigquery.FloatFieldType},
		{Name: "user_id", Type: bigquery.StringFieldType},
		{Name: "user_data", Type: bigquery.JSONFieldType},
		{Name: "tag", Type: bigquery.StringFieldType},
		{Name: "source_id", Type: bigquery.StringFieldType},
		{Name: "feature_id", Type: bigquery.StringFieldType},
		{Name: "feature_version", Type: bigquery.IntegerFieldType},
		{Name: "variation_id", Type: bigquery.StringFieldType},
		{Name: "reason", Type: bigquery.StringFieldType},
	}
	err = createTable(ctx, client, logger, *c.project, *c.dataset, "goal_event", goalEventSchema)
	if err != nil {
		if isAlreadyExistsError(err) {
			logger.Info("goal_event table already exists, skipping creation")
		} else {
			logger.Error("failed to create goal_event table",
				zap.Error(err),
			)
			return err
		}
	} else {
		logger.Info("successfully created goal_event table")
	}

	logger.Info("BigQuery table creation completed successfully")
	return nil
}

func createDataset(ctx context.Context,
	client *bigquery.Client, logger *zap.Logger, datasetID string,
) error {
	meta := &bigquery.DatasetMetadata{}
	if err := client.Dataset(datasetID).Create(ctx, meta); err != nil {
		if isAlreadyExistsError(err) {
			logger.Info("dataset already exists, skipping creation",
				zap.String("dataset", datasetID))
			return nil
		}
		return err
	}
	logger.Info("successfully created dataset",
		zap.String("dataset", datasetID))
	return nil
}

func createTable(ctx context.Context,
	client *bigquery.Client, logger *zap.Logger,
	projectID, datasetID, tableID string,
	schema bigquery.Schema,
) error {
	tableMetadata := &bigquery.TableMetadata{
		Schema: schema,
	}
	tableRef := client.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, tableMetadata); err != nil {
		return err
	}
	return nil
}

// isAlreadyExistsError checks if the error indicates the resource already exists.
// BigQuery emulator may return different error messages, so we check for common patterns.
func isAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	// Check for common "already exists" patterns
	return strings.Contains(errStr, "already exists") ||
		strings.Contains(errStr, "already created") ||
		strings.Contains(errStr, "duplicate") ||
		strings.Contains(errStr, "409") // HTTP 409 Conflict
}
