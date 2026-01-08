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
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
)

var (
	// Tables to truncate in the data warehouse
	dataWarehouseTables = []string{
		"evaluation_event",
		"goal_event",
	}
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
	cmd := p.Command("truncate", "Truncate data warehouse tables")
	command := &command{
		CmdClause:   cmd,
		mysqlUser:   cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:   cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:   cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:   cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName: cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
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

	for _, table := range dataWarehouseTables {
		query := fmt.Sprintf("TRUNCATE TABLE %s", table)
		_, err := client.ExecContext(ctx, query)
		if err != nil {
			logger.Error("Failed to truncate table", zap.Error(err), zap.String("table", table))
			return err
		}
		logger.Info("Truncated table", zap.String("table", table))
	}

	logger.Info("All data warehouse tables truncated successfully")
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
