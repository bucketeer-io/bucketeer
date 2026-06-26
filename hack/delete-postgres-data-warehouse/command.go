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

var (
	// Tables to truncate in the data warehouse
	dataWarehouseTables = []string{
		"evaluation_event",
		"goal_event",
	}
)

type command struct {
	*kingpin.CmdClause
	postgresUser   *string
	postgresPass   *string
	postgresHost   *string
	postgresPort   *int
	postgresDBName *string
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("truncate", "Truncate data warehouse tables")
	command := &command{
		CmdClause:      cmd,
		postgresUser:   cmd.Flag("postgres-user", "PostgreSQL user.").Required().String(),
		postgresPass:   cmd.Flag("postgres-pass", "PostgreSQL password.").Required().String(),
		postgresHost:   cmd.Flag("postgres-host", "PostgreSQL host.").Required().String(),
		postgresPort:   cmd.Flag("postgres-port", "PostgreSQL port.").Required().Int(),
		postgresDBName: cmd.Flag("postgres-db-name", "PostgreSQL database name.").Required().String(),
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
