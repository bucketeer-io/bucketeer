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
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cli"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
)

const (
	dialTimeout = 10 * time.Second

	// Numeric values match the proto enums in proto/account/account.proto:
	// Role_Organization: UNASSIGNED=0, MEMBER=1, ADMIN=2, OWNER=3
	// Role_Environment:  UNASSIGNED=0, VIEWER=1, EDITOR=2
	roleOrganizationADMIN = 2
	roleOrganizationOWNER = 3
	roleEnvironmentEDITOR = 2
)

// Upserts a row in account_v2; MySQL's composite PK is (email, organization_id),
// so re-running against a DB that already has the row is safe.
//
// `tags` is json NOT NULL with no DB default (added by migration
// 20250115040347_update_account_v2_table.sql), so it must be set explicitly.
// The other columns added by later migrations (first_name, last_name, language,
// last_seen, avatar_file_type, avatar_image, teams, search_filters) are either
// nullable or have DEFAULT values, so they don't need to be set here.
//
//go:embed sql/upsert_account_v2.sql
var upsertAccountV2SQL string

type command struct {
	*kingpin.CmdClause
	mysqlUser             *string
	mysqlPass             *string
	mysqlHost             *string
	mysqlPort             *int
	mysqlDBName           *string
	email                 *string
	defaultOrganizationID *string
	e2eOrganizationID     *string
	e2eEnvironmentID      *string
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command(
		"create",
		"Bootstrap an account with ADMIN role in a default organization and OWNER role in a system-admin organization.",
	)
	command := &command{
		CmdClause:   cmd,
		mysqlUser:   cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:   cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:   cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:   cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName: cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		email: cmd.Flag(
			"email",
			"Email of the account to create.",
		).Required().String(),
		defaultOrganizationID: cmd.Flag(
			"default-organization-id",
			"ID of the default organization where the account gets ADMIN role + EDITOR role on the e2e environment.",
		).Required().String(),
		e2eOrganizationID: cmd.Flag(
			"e2e-organization-id",
			"ID of the e2e (system-admin) organization where the account gets OWNER role.",
		).Required().String(),
		e2eEnvironmentID: cmd.Flag(
			"e2e-environment-id",
			"ID of the e2e environment that the default-org account gets EDITOR access to.",
		).Required().String(),
	}
	r.RegisterCommand(command)
	return command
}

func (c *command) Run(ctx context.Context, _ metrics.Metrics, logger *zap.Logger) error {
	client, err := c.createMySQLClient(ctx, logger)
	if err != nil {
		logger.Error("Failed to create mysql client", zap.Error(err))
		return err
	}
	defer client.Close()

	name := strings.Split(*c.email, "@")[0]
	now := time.Now().Unix()

	// Membership in the default org: ADMIN role + EDITOR role on the e2e environment.
	// environment_roles is JSON; the column shape matches what pkg/account/domain
	// reads ([]*AccountV2_EnvironmentRole serialized to JSON).
	envRoles, err := json.Marshal([]map[string]interface{}{
		{
			"environment_id": *c.e2eEnvironmentID,
			"role":           roleEnvironmentEDITOR,
		},
	})
	if err != nil {
		return fmt.Errorf("marshal environment_roles: %w", err)
	}
	if err := c.upsertAccount(
		ctx, client, logger,
		*c.email, name, *c.defaultOrganizationID,
		roleOrganizationADMIN, string(envRoles), now,
	); err != nil {
		return err
	}

	// Membership in the e2e (system-admin) org: OWNER role, no environment roles.
	// The e2e organization is expected to be flagged as a system-admin organization,
	// so any account in it is granted system admin privileges.
	if err := c.upsertAccount(
		ctx, client, logger,
		*c.email, name, *c.e2eOrganizationID,
		roleOrganizationOWNER, "[]", now,
	); err != nil {
		return err
	}

	logger.Info(
		"Account is ready",
		zap.String("email", *c.email),
		zap.String("defaultOrganizationId", *c.defaultOrganizationID),
		zap.String("e2eOrganizationId", *c.e2eOrganizationID),
		zap.String("e2eEnvironmentId", *c.e2eEnvironmentID),
	)
	return nil
}

func (c *command) upsertAccount(
	ctx context.Context,
	client mysql.Client,
	logger *zap.Logger,
	email, name, orgID string,
	orgRole int,
	envRolesJSON string,
	now int64,
) error {
	_, err := client.ExecContext(
		ctx,
		upsertAccountV2SQL,
		email,
		name,
		"",   // avatar_image_url
		"[]", // tags (json NOT NULL, no default)
		orgID,
		orgRole,
		envRolesJSON,
		0, // disabled
		now,
		now,
	)
	if err != nil {
		logger.Error(
			"Failed to upsert account",
			zap.Error(err),
			zap.String("email", email),
			zap.String("organizationId", orgID),
		)
		return err
	}
	logger.Info(
		"Account upserted",
		zap.String("email", email),
		zap.String("organizationId", orgID),
		zap.Int("organizationRole", orgRole),
	)
	return nil
}

func (c *command) createMySQLClient(
	ctx context.Context,
	logger *zap.Logger,
) (mysql.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, dialTimeout)
	defer cancel()
	return mysql.NewClient(
		ctx,
		*c.mysqlUser, *c.mysqlPass, *c.mysqlHost,
		*c.mysqlPort,
		*c.mysqlDBName,
		mysql.WithLogger(logger),
	)
}
