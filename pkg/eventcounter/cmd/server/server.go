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

package server

import (
	"context"
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/eventcounter/api"
	"github.com/bucketeer-io/bucketeer/pkg/eventcounter/druid"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	storagedruid "github.com/bucketeer-io/bucketeer/pkg/storage/druid"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/token"
)

const command = "server"

type server struct {
	*kingpin.CmdClause
	port                  *int
	project               *string
	mysqlUser             *string
	mysqlPass             *string
	mysqlHost             *string
	mysqlPort             *int
	mysqlDBName           *string
	experimentService     *string
	featureService        *string
	accountService        *string
	certPath              *string
	keyPath               *string
	serviceTokenPath      *string
	oauthKeyPath          *string
	oauthClientID         *string
	oauthIssuer           *string
	druidURL              *string
	druidDatasourcePrefix *string
	druidUsername         *string
	druidPassword         *string
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the server")
	server := &server{
		CmdClause:   cmd,
		port:        cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:     cmd.Flag("project", "Google Cloud project name.").String(),
		mysqlUser:   cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:   cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:   cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:   cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName: cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		experimentService: cmd.Flag(
			"experiment-service",
			"bucketeer-experiment-service address.",
		).Default("experiment:9090").String(),
		featureService: cmd.Flag(
			"feature-service",
			"bucketeer-feature-service address.",
		).Default("feature:9090").String(),
		accountService: cmd.Flag(
			"account-service",
			"bucketeer-account-service address.",
		).Default("account:9090").String(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		oauthKeyPath: cmd.Flag(
			"oauth-key",
			"Path to public key used to verify oauth token.",
		).Required().String(),
		oauthClientID:         cmd.Flag("oauth-client-id", "The oauth clientID registered at dex.").Required().String(),
		oauthIssuer:           cmd.Flag("oauth-issuer", "The url of dex issuer.").Required().String(),
		druidURL:              cmd.Flag("druid-url", "Druid URL.").String(),
		druidDatasourcePrefix: cmd.Flag("druid-datasource-prefix", "Druid datasource prefix.").String(),
		druidUsername:         cmd.Flag("druid-username", "Druid username.").String(),
		druidPassword:         cmd.Flag("druid-password", "Druid password.").String(),
	}
	r.RegisterCommand(server)
	return server
}

func (s *server) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()

	mysqlClient, err := s.createMySQLClient(ctx, registerer, logger)
	if err != nil {
		return err
	}
	defer mysqlClient.Close()

	creds, err := rpcclient.NewPerRPCCredentials(*s.serviceTokenPath)
	if err != nil {
		return err
	}

	experimentClient, err := experimentclient.NewClient(*s.experimentService, *s.certPath,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
		rpcclient.WithMetrics(registerer),
		rpcclient.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer experimentClient.Close()

	featureClient, err := featureclient.NewClient(*s.featureService, *s.certPath,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
		rpcclient.WithMetrics(registerer),
		rpcclient.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer featureClient.Close()

	accountClient, err := accountclient.NewClient(*s.accountService, *s.certPath,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
		rpcclient.WithMetrics(registerer),
		rpcclient.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer accountClient.Close()

	druidQuerier, err := s.createDruidQuerier(ctx, logger)
	if err != nil {
		logger.Error("Failed to create druid querier", zap.Error(err))
		return err
	}

	service := api.NewEventCounterService(
		mysqlClient,
		experimentClient,
		featureClient,
		accountClient,
		druidQuerier,
		registerer,
		logger,
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

func (s *server) createMySQLClient(
	ctx context.Context,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (mysql.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return mysql.NewClient(
		ctx,
		*s.mysqlUser, *s.mysqlPass, *s.mysqlHost,
		*s.mysqlPort,
		*s.mysqlDBName,
		mysql.WithLogger(logger),
		mysql.WithMetrics(registerer),
	)
}

func (s *server) createDruidQuerier(ctx context.Context, logger *zap.Logger) (druid.Querier, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	brokerClient, err := storagedruid.NewBrokerClient(ctx, *s.druidURL, *s.druidUsername, *s.druidPassword)
	if err != nil {
		logger.Error("Failed to create druid broker client", zap.Error(err))
		return nil, err
	}
	return druid.NewDruidQuerier(brokerClient, *s.druidDatasourcePrefix, druid.WithLogger(logger)), nil
}
