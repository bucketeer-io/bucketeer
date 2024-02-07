// Copyright 2024 The Bucketeer Authors.
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

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

const (
	envNamespace   = "e2e"
	organizationID = "e2e"
	prefixTestName = "e2e-test"
)

var (
	targetEntities = []*mysqlE2EInfo{
		{table: "subscription", targetField: "name"},
		{table: "experiment_result", targetField: ""},
		{table: "push", targetField: "name"},
		{table: "ops_count", targetField: ""},
		{table: "auto_ops_rule", targetField: "feature_id"},
		{table: "segment_user", targetField: "user_id"},
		{table: "segment", targetField: "name"},
		{table: "goal", targetField: "id"},
		{table: "experiment", targetField: "feature_id"},
		{table: "tag", targetField: ""},
		{table: "ops_progressive_rollout", targetField: "feature_id"},
		{table: "flag_trigger", targetField: "description"},
		{table: "feature", targetField: "id"},
	}
	targetEntitiesInOrganization = []*mysqlE2EInfo{
		{table: "account_v2", targetField: "email"},
	}
)

type mysqlE2EInfo struct {
	table       string
	targetField string
}

type command struct {
	*kingpin.CmdClause
	mysqlUser   *string
	mysqlPass   *string
	mysqlHost   *string
	mysqlPort   *int
	mysqlDBName *string
	testID      *string
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("delete", "delete e2e data")
	command := &command{
		CmdClause:   cmd,
		mysqlUser:   cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:   cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:   cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:   cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName: cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		testID:      cmd.Flag("test-id", "Test ID.").String(),
	}
	r.RegisterCommand(command)
	return command
}

func (c *command) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	client, err := c.createMySQLClient(ctx, logger)
	if err != nil {
		logger.Error("Failed to create mysql client", zap.Error(err))
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

func (c *command) createMySQLClient(
	ctx context.Context,
	logger *zap.Logger,
) (mysql.Client, error) {
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

func (c *command) constructDeleteQuery(target *mysqlE2EInfo) (query string, args []interface{}) {
	if target.targetField != "" && *c.testID != "" {
		query = fmt.Sprintf(`
			DELETE FROM
				%s
			WHERE
				environment_namespace = ? AND
				%s LIKE ?
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
			environment_namespace = ?
	`, target.table)
	args = []interface{}{
		envNamespace,
	}
	return
}

func (c *command) constructDeleteQueryForOrganization(target *mysqlE2EInfo) (query string, args []interface{}) {
	if target.targetField != "" && *c.testID != "" {
		query = fmt.Sprintf(`
			DELETE FROM
				%s
			WHERE
				organization_id = ? AND
				%s LIKE ?
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
			organization_id = ?
	`, target.table)
	args = []interface{}{
		organizationID,
	}
	return
}
