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
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cli"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
)

const (
	dialTimeout = 10 * time.Second

	// Numeric values match the proto enums in proto/account/account.proto:
	// Role_Organization: UNASSIGNED=0, MEMBER=1, ADMIN=2, OWNER=3
	roleOrganizationMEMBER = 1
	roleOrganizationADMIN  = 2
	roleOrganizationOWNER  = 3

	// Role_Environment: UNASSIGNED=0, VIEWER=1, EDITOR=2
	roleEnvironmentVIEWER = 1
	roleEnvironmentEDITOR = 2

	// Org roles >= ADMIN already grant access to every environment in the org
	// (see pkg/account/domain/account.go: ChangeOrganizationRole), so the
	// org-owner rows carry an empty environment_roles array.
	emptyEnvironmentRolesJSON = "[]"

	// Access tokens are minted with a far-future expiry so a single bootstrap
	// run keeps working for the lifetime of a local/e2e environment.
	tokenTTLYears = 100
)

// environmentRole mirrors the JSON shape of proto AccountV2_EnvironmentRole as
// it is stored in the account_v2.environment_roles column
// (e.g. [{"environment_id":"e2e","role":2}]).
type environmentRole struct {
	EnvironmentID string `json:"environment_id"`
	Role          int    `json:"role"`
}

// Upserts a row in account_v2; MySQL's composite PK is (email, organization_id),
// so re-running against a DB that already has the row is safe.
//
// `tags` is json NOT NULL with no DB default (added by migration
// 20250115040347_update_account_v2_table.sql), so it must be set explicitly.
// The other columns added by later migrations (first_name, last_name, language,
// last_seen, avatar_file_type, avatar_image, teams, search_filters) are either
// nullable or have DEFAULT values, so they don't need to be set here.
//
// The query uses MySQL 8.0's row-alias form (`VALUES (...) AS new`) instead of
// the deprecated `VALUES(col)` reference inside ON DUPLICATE KEY UPDATE.
//
//go:embed sql/upsert_account_v2.sql
var upsertAccountV2SQL string

type command struct {
	*kingpin.CmdClause
	mysqlUser                  *string
	mysqlPass                  *string
	mysqlHost                  *string
	mysqlPort                  *int
	mysqlDBName                *string
	sysAdminEmail              *string
	orgOwnerEmail              *string
	envWriteEmail              *string
	envReadEmail               *string
	defaultOrganizationID      *string
	e2eOrganizationID          *string
	e2eEnvironmentID           *string
	oauthKeyPath               *string
	issuer                     *string
	audience                   *string
	sysAdminTokenOutput        *string
	orgOwnerDefaultTokenOutput *string
	orgOwnerE2ETokenOutput     *string
	envWriteTokenOutput        *string
	envReadTokenOutput         *string
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command(
		"create",
		"Bootstrap the local/e2e accounts (system admin, org owner, environment "+
			"editor, environment viewer) and generate their access tokens",
	)
	command := &command{
		CmdClause:   cmd,
		mysqlUser:   cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:   cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:   cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:   cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName: cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		sysAdminEmail: cmd.Flag(
			"sys-admin-email",
			"Email of the system admin account (OWNER of the e2e/system-admin organization; its token is a system admin).",
		).Required().String(),
		orgOwnerEmail: cmd.Flag(
			"org-owner-email",
			"Email of the organization owner account (OWNER in both the default and e2e organizations).",
		).Required().String(),
		envWriteEmail: cmd.Flag(
			"env-write-email",
			"Email of the environment editor account (MEMBER of the default organization with EDITOR role on the e2e environment).",
		).Required().String(),
		envReadEmail: cmd.Flag(
			"env-read-email",
			"Email of the environment viewer account (MEMBER of the default organization with VIEWER role on the e2e environment).",
		).Required().String(),
		defaultOrganizationID: cmd.Flag(
			"default-organization-id",
			"ID of the default organization that owns the e2e environment.",
		).Required().String(),
		e2eOrganizationID: cmd.Flag(
			"e2e-organization-id",
			"ID of the e2e organization where the org owner account also gets OWNER role.",
		).Required().String(),
		e2eEnvironmentID: cmd.Flag(
			"e2e-environment-id",
			"ID of the e2e environment used for the editor/viewer environment roles.",
		).Required().String(),
		oauthKeyPath: cmd.Flag(
			"oauth-key",
			"Path to the OAuth RSA private key used to sign the access tokens.",
		).Required().String(),
		issuer: cmd.Flag(
			"issuer",
			"Issuer URL set in the generated access tokens.",
		).Required().String(),
		audience: cmd.Flag(
			"audience",
			"OAuth audience set in the generated access tokens.",
		).Default("bucketeer").String(),
		sysAdminTokenOutput: cmd.Flag(
			"sys-admin-token-output",
			"Path of the file to write the system admin access token.",
		).Required().String(),
		orgOwnerDefaultTokenOutput: cmd.Flag(
			"org-owner-default-token-output",
			"Path of the file to write the org owner access token (scoped to the default organization).",
		).Required().String(),
		orgOwnerE2ETokenOutput: cmd.Flag(
			"org-owner-e2e-token-output",
			"Path of the file to write the org owner access token scoped to the e2e organization.",
		).Required().String(),
		envWriteTokenOutput: cmd.Flag(
			"env-write-token-output",
			"Path of the file to write the environment editor access token.",
		).Required().String(),
		envReadTokenOutput: cmd.Flag(
			"env-read-token-output",
			"Path of the file to write the environment viewer access token.",
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

	now := time.Now().Unix()

	editorEnvRoles, err := environmentRolesJSON(*c.e2eEnvironmentID, roleEnvironmentEDITOR)
	if err != nil {
		return err
	}
	viewerEnvRoles, err := environmentRolesJSON(*c.e2eEnvironmentID, roleEnvironmentVIEWER)
	if err != nil {
		return err
	}

	// System admin: OWNER of the e2e (system-admin) organization. Its token is
	// minted with is_system_admin=true, so it can call system-admin-only APIs
	// (e.g. creating organizations/environments) and read across organizations.
	if err := c.upsertAccount(
		ctx, client, logger,
		*c.sysAdminEmail, *c.e2eOrganizationID,
		roleOrganizationOWNER, emptyEnvironmentRolesJSON, now,
	); err != nil {
		return err
	}

	// Org owner: OWNER of the default org (which owns the e2e environment) and
	// OWNER of the e2e org. Org role >= ADMIN implies access to every
	// environment in the org, so environment_roles stays empty.
	if err := c.upsertAccount(
		ctx, client, logger,
		*c.orgOwnerEmail, *c.defaultOrganizationID,
		roleOrganizationOWNER, emptyEnvironmentRolesJSON, now,
	); err != nil {
		return err
	}
	if err := c.upsertAccount(
		ctx, client, logger,
		*c.orgOwnerEmail, *c.e2eOrganizationID,
		roleOrganizationOWNER, emptyEnvironmentRolesJSON, now,
	); err != nil {
		return err
	}

	// Environment editor: MEMBER of the default org with EDITOR role on the e2e environment.
	if err := c.upsertAccount(
		ctx, client, logger,
		*c.envWriteEmail, *c.defaultOrganizationID,
		roleOrganizationMEMBER, editorEnvRoles, now,
	); err != nil {
		return err
	}

	// Environment viewer: MEMBER of the default org with VIEWER role on the e2e environment.
	if err := c.upsertAccount(
		ctx, client, logger,
		*c.envReadEmail, *c.defaultOrganizationID,
		roleOrganizationMEMBER, viewerEnvRoles, now,
	); err != nil {
		return err
	}

	// Mint an access token per account. The system admin token is scoped to the
	// e2e (system-admin) org and carries is_system_admin=true. The other three
	// are scoped to the default org and are NOT system admins — the org owner
	// relies on its org OWNER role and the editor/viewer rely on their
	// environment roles, exercising the real RBAC path.
	signer, err := token.NewSigner(*c.oauthKeyPath)
	if err != nil {
		logger.Error("Failed to create token signer", zap.Error(err))
		return err
	}
	for _, t := range []struct {
		email          string
		output         string
		organizationID string
		isSystemAdmin  bool
	}{
		{*c.sysAdminEmail, *c.sysAdminTokenOutput, *c.e2eOrganizationID, true},
		{*c.orgOwnerEmail, *c.orgOwnerDefaultTokenOutput, *c.defaultOrganizationID, false},
		{*c.orgOwnerEmail, *c.orgOwnerE2ETokenOutput, *c.e2eOrganizationID, false},
		{*c.envWriteEmail, *c.envWriteTokenOutput, *c.defaultOrganizationID, false},
		{*c.envReadEmail, *c.envReadTokenOutput, *c.defaultOrganizationID, false},
	} {
		if err := c.writeAccessToken(logger, signer, t.email, t.output, t.organizationID, t.isSystemAdmin); err != nil {
			return err
		}
	}

	logger.Info(
		"Accounts and access tokens are ready",
		zap.String("sysAdminEmail", *c.sysAdminEmail),
		zap.String("orgOwnerEmail", *c.orgOwnerEmail),
		zap.String("envWriteEmail", *c.envWriteEmail),
		zap.String("envReadEmail", *c.envReadEmail),
		zap.String("defaultOrganizationId", *c.defaultOrganizationID),
		zap.String("e2eOrganizationId", *c.e2eOrganizationID),
		zap.String("e2eEnvironmentId", *c.e2eEnvironmentID),
	)
	return nil
}

func environmentRolesJSON(environmentID string, role int) (string, error) {
	b, err := json.Marshal([]environmentRole{{EnvironmentID: environmentID, Role: role}})
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *command) writeAccessToken(
	logger *zap.Logger,
	signer token.Signer,
	email, output, organizationID string,
	isSystemAdmin bool,
) error {
	now := time.Now()
	accessToken := &token.AccessToken{
		Issuer:         *c.issuer,
		Audience:       *c.audience,
		Expiry:         now.AddDate(tokenTTLYears, 0, 0),
		IssuedAt:       now,
		Email:          email,
		Name:           accountName(email),
		OrganizationID: organizationID,
		IsSystemAdmin:  isSystemAdmin,
		IsServiceToken: false,
	}
	signed, err := signer.SignAccessToken(accessToken)
	if err != nil {
		logger.Error("Failed to sign access token", zap.Error(err), zap.String("email", email))
		return err
	}
	if err := os.WriteFile(output, []byte(signed), 0644); err != nil {
		logger.Error("Failed to write access token", zap.Error(err), zap.String("output", output))
		return err
	}
	logger.Info("Access token generated", zap.String("email", email), zap.String("output", output))
	return nil
}

func (c *command) upsertAccount(
	ctx context.Context,
	client mysql.Client,
	logger *zap.Logger,
	email, orgID string,
	orgRole int,
	environmentRolesJSON string,
	now int64,
) error {
	_, err := client.ExecContext(
		ctx,
		upsertAccountV2SQL,
		email,
		accountName(email),
		"",   // avatar_image_url
		"[]", // tags (json NOT NULL, no default)
		orgID,
		orgRole,
		environmentRolesJSON,
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
		zap.String("environmentRoles", environmentRolesJSON),
	)
	return nil
}

func accountName(email string) string {
	return strings.Split(email, "@")[0]
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
