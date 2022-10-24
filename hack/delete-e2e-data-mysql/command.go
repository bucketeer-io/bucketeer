// Copyright 2022 The Bucketeer Authors.
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
	"strings"
	"time"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

const (
	envNamespace   = "e2e"
	prefixTestName = "e2e-test"
)

var (
	targetEntities = []*mysqlE2EInfo{
		{table: "subscription", targetField: "name", hasCreatedAt: false},
		{table: "experiment_result", targetField: "", hasCreatedAt: false},
		{table: "push", targetField: "name", hasCreatedAt: true},
		{table: "ops_count", targetField: "", hasCreatedAt: false},
		{table: "auto_ops_rule", targetField: "feature_id", hasCreatedAt: true},
		{table: "segment_user", targetField: "user_id", hasCreatedAt: false},
		{table: "segment", targetField: "name", hasCreatedAt: true},
		{table: "goal", targetField: "id", hasCreatedAt: true},
		{table: "experiment", targetField: "feature_id", hasCreatedAt: true},
		{table: "tag", targetField: "", hasCreatedAt: true},
		{table: "feature", targetField: "id", hasCreatedAt: true},
		{table: "webhook", targetField: "name", hasCreatedAt: true},
	}
)

type mysqlE2EInfo struct {
	table        string
	targetField  string
	hasCreatedAt bool
}

type command struct {
	*kingpin.CmdClause
	mysqlUser        *string
	mysqlPass        *string
	mysqlHost        *string
	mysqlPort        *int
	mysqlDBName      *string
	testID           *string
	retentionSeconds *int
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("delete", "delete e2e data")
	command := &command{
		CmdClause:        cmd,
		mysqlUser:        cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:        cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:        cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:        cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName:      cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		testID:           cmd.Flag("test-id", "Test ID.").String(),
		retentionSeconds: cmd.Flag("retention-seconds", "Test data retention period(seconds)").Int(),
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
		if err := c.deleteData(ctx, client, target); err != nil {
			logger.Error("Failed to delete data", zap.Error(err), zap.String("table", target.table))
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

func (c *command) deleteData(ctx context.Context, client mysql.Client, target *mysqlE2EInfo) error {
	query, args := c.constructDeleteQuery(target)
	_, err := client.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *command) constructDeleteQuery(target *mysqlE2EInfo) (query string, args []interface{}) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`
		DELETE FROM
			%s
		WHERE
			environment_namespace = ?
	`, target.table))
	args = []interface{}{
		envNamespace,
	}

	if target.targetField != "" && *c.testID != "" {
		sb.WriteString("AND " + target.targetField + " LIKE ?\n")
		targetName := prefixTestName + "-" + *c.testID + "%"
		args = append(args, targetName)
	}

	if target.hasCreatedAt && *c.retentionSeconds > 0 {
		sb.WriteString("AND created_at < ?\n")
		t := time.Now().Add(-1 * time.Duration(*c.retentionSeconds) * time.Second).Unix()
		args = append(args, t)
	}

	query = sb.String()
	return
}
