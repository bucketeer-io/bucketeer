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

package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cli"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
)

const (
	envNamespace   = "e2e"
	organizationID = "e2e"
	prefixTestName = "e2e-test"
)

var (
	targetEntities = []*postgresE2EInfo{
		{table: "subscription", targetField: "name"},
		{table: "experiment_result", targetField: ""},
		{table: "push", targetField: "name"},
		{table: "ops_count", targetField: ""},
		{table: "auto_ops_rule", targetField: "feature_id"},
		{table: "segment_user", targetField: "user_id"},
		{table: "segment", targetField: "name"},
		{table: "goal", targetField: "id"},
		{table: "experiment", targetField: "feature_id"},
		{table: "tag", targetField: "name"},
		{table: "ops_progressive_rollout", targetField: "feature_id"},
		{table: "flag_trigger", targetField: "description"},
		{table: "code_reference", targetField: "feature_id"},
		{table: "scheduled_feature_change", targetField: "feature_id"},
		{table: "feature", targetField: "id"},
		{table: "api_key", targetField: "name"},
	}
	targetEntitiesInOrganization = []*postgresE2EInfo{
		{table: "account_v2", targetField: "email"},
	}
)

type postgresE2EInfo struct {
	table       string
	targetField string
}

type command struct {
	*kingpin.CmdClause
	postgresUser   *string
	postgresPass   *string
	postgresHost   *string
	postgresPort   *int
	postgresDBName *string
	testID         *string
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("delete", "delete e2e data")
	command := &command{
		CmdClause:      cmd,
		postgresUser:   cmd.Flag("postgres-user", "PostgreSQL user.").Required().String(),
		postgresPass:   cmd.Flag("postgres-pass", "PostgreSQL password.").Required().String(),
		postgresHost:   cmd.Flag("postgres-host", "PostgreSQL host.").Required().String(),
		postgresPort:   cmd.Flag("postgres-port", "PostgreSQL port.").Required().Int(),
		postgresDBName: cmd.Flag("postgres-db-name", "PostgreSQL database name.").Required().String(),
		testID:         cmd.Flag("test-id", "Test ID.").String(),
	}
	r.RegisterCommand(command)
	return command
}

func (c *command) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	client, err := c.createPostgresClient(ctx, logger)
	if err != nil {
		logger.Error("Failed to create postgres client", zap.Error(err))
		return err
	}
	defer client.Close()
	for _, target := range targetEntities {
		query, args := c.constructDeleteQuery(target)
		_, err := client.ExecContext(
			ctx,
			query,
			args...,
		)
		if err != nil {
			logger.Error("Failed to delete data", zap.Error(err), zap.String("table", target.table))
			return err
		}
	}
	for _, target := range targetEntitiesInOrganization {
		query, args := c.constructDeleteQueryForOrganization(target)
		_, err := client.ExecContext(
			ctx,
			query,
			args...,
		)
		if err != nil {
			logger.Error("Failed to delete data in organization", zap.Error(err), zap.String("table", target.table))
			return err
		}
	}
	logger.Info("Done")
	return nil
}

func (c *command) createPostgresClient(
	ctx context.Context,
	logger *zap.Logger,
) (postgres.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return postgres.NewClient(
		ctx,
		*c.postgresUser, *c.postgresPass, *c.postgresHost,
		*c.postgresPort,
		*c.postgresDBName,
		postgres.WithLogger(logger),
	)
}

func (c *command) constructDeleteQuery(target *postgresE2EInfo) (query string, args []interface{}) {
	if target.targetField != "" && *c.testID != "" {
		query = fmt.Sprintf(`
			DELETE FROM
				%s
			WHERE
				environment_id = $1 AND
				%s LIKE $2
		`, target.table, target.targetField)
		args = []interface{}{
			envNamespace,
			prefixTestName + "-" + *c.testID + "%",
		}
		return
	}
	query = fmt.Sprintf(`
		DELETE FROM
			%s
		WHERE
			environment_id = $1
	`, target.table)
	args = []interface{}{
		envNamespace,
	}
	return
}

func (c *command) constructDeleteQueryForOrganization(target *postgresE2EInfo) (query string, args []interface{}) {
	if target.targetField != "" && *c.testID != "" {
		query = fmt.Sprintf(`
			DELETE FROM
				%s
			WHERE
			    name != 'localenv' AND -- to avoid deleting the service account used for e2e test
				organization_id = $1 AND
				%s LIKE $2
		`, target.table, target.targetField)
		args = []interface{}{
			organizationID,
			prefixTestName + "-" + *c.testID + "%",
		}
		return
	}
	query = fmt.Sprintf(`
		DELETE FROM
			%s
		WHERE
			organization_id = $1
		AND
		    name != 'localenv' -- to avoid deleting the service account used for e2e test
	`, target.table)
	args = []interface{}{
		organizationID,
	}
	return
}
