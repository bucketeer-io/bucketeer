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

package mysqlserver

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/migration/mysql/api"
	"github.com/bucketeer-io/bucketeer/pkg/migration/mysql/migrate"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/token"
)

const command = "mysql-server"

type server struct {
	*kingpin.CmdClause
	port                      *int
	githubUser                *string
	githubAccessTokenPath     *string
	githubMigrationSourcePath *string
	mysqlUser                 *string
	mysqlPass                 *string
	mysqlHost                 *string
	mysqlPort                 *int
	mysqlDBName               *string
	certPath                  *string
	keyPath                   *string
	oauthKeyPath              *string
	oauthClientID             *string
	oauthIssuer               *string
}

func RegisterServerCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the gRPC server")
	server := &server{
		CmdClause:             cmd,
		port:                  cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		githubUser:            cmd.Flag("github-user", "GitHub user.").Required().String(),
		githubAccessTokenPath: cmd.Flag("github-access-token-path", "Path to GitHub access token.").Required().String(),
		githubMigrationSourcePath: cmd.Flag(
			"github-migration-source-path",
			"Path to migration file in GitHub. (e.g. owner/repo/path#ref)",
		).Required().String(),
		mysqlUser:   cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:   cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:   cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:   cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName: cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		certPath:    cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:     cmd.Flag("key", "Path to TLS key.").Required().String(),
		oauthKeyPath: cmd.Flag(
			"oauth-key",
			"Path to public key used to verify oauth token.",
		).Required().String(),
		oauthClientID: cmd.Flag("oauth-client-id", "The oauth clientID registered at dex.").Required().String(),
		oauthIssuer:   cmd.Flag("oauth-issuer", "The url of dex issuer.").Required().String(),
	}
	r.RegisterCommand(server)
	return server
}

func (s *server) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()

	migrateClientFactory, err := migrate.NewClientFactory(
		*s.githubUser, *s.githubAccessTokenPath, *s.githubMigrationSourcePath,
		*s.mysqlUser, *s.mysqlPass, *s.mysqlHost, *s.mysqlPort, *s.mysqlDBName,
	)
	if err != nil {
		return err
	}

	service := api.NewMySQLService(
		migrateClientFactory,
		api.WithLogger(logger),
	)

	verifier, err := token.NewVerifier(*s.oauthKeyPath, *s.oauthIssuer, *s.oauthClientID)
	if err != nil {
		return err
	}

	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
	)
	go healthChecker.Run(ctx)

	server := rpc.NewServer(service, *s.certPath, *s.keyPath,
		"migrate-mysql",
		rpc.WithPort(*s.port),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
		rpc.WithService(healthChecker),
		rpc.WithHandler("/health", healthChecker),
	)
	defer server.Stop(10 * time.Second)
	go server.Run()

	<-ctx.Done()
	return nil
}
