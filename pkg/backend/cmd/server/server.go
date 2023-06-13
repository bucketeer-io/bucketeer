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
	"regexp"
	"time"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	accountapi "github.com/bucketeer-io/bucketeer/pkg/account/api"
	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	authapi "github.com/bucketeer-io/bucketeer/pkg/auth/api"
	"github.com/bucketeer-io/bucketeer/pkg/auth/oidc"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/rest"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/token"
)

const (
	command = "server"
)

type server struct {
	*kingpin.CmdClause
	project             *string
	mysqlUser           *string
	mysqlPass           *string
	mysqlHost           *string
	mysqlPort           *int
	mysqlDBName         *string
	domainTopic         *string
	accountServicePort  *int
	authServicePort     *int
	accountService      *string
	environmentService  *string
	certPath            *string
	keyPath             *string
	serviceTokenPath    *string
	oauthPrivateKeyPath *string
	oauthPublicKeyPath  *string
	oauthClientID       *string
	oauthClientSecret   *string
	oauthRedirectURLs   *[]string
	oauthIssuer         *string
	oauthIssuerCertPath *string
	emailFilter         *string
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the server")
	server := &server{
		CmdClause:   cmd,
		project:     cmd.Flag("project", "Google Cloud project name.").String(),
		mysqlUser:   cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:   cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:   cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:   cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName: cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		domainTopic: cmd.Flag("domain-topic", "PubSub topic to publish domain events.").Required().String(),
		accountServicePort: cmd.Flag(
			"account-service-port",
			"Port to bind to account service.",
		).Default("9091").Int(),
		authServicePort: cmd.Flag(
			"auth-service-port",
			"Port to bind to auth service.",
		).Default("9092").Int(),
		accountService: cmd.Flag(
			"account-service",
			"bucketeer-account-service address.",
		).Default("account:9090").String(),
		environmentService: cmd.Flag(
			"environment-service",
			"bucketeer-environment-service address.",
		).Default("environment:9090").String(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		oauthPrivateKeyPath: cmd.Flag(
			"oauth-private-key",
			"Path to private key for signing oauth token.",
		).Required().String(),
		oauthPublicKeyPath: cmd.Flag(
			"oauth-public-key",
			"Path to public key used to verify oauth token.",
		).Required().String(),
		oauthClientID: cmd.Flag(
			"oauth-client-id",
			"The oauth clientID registered at dex.",
		).Required().String(),
		oauthClientSecret: cmd.Flag(
			"oauth-client-secret",
			"The oauth client secret registered at Dex.",
		).Required().String(),
		oauthRedirectURLs:   cmd.Flag("oauth-redirect-urls", "The redirect urls registered at Dex.").Required().Strings(),
		oauthIssuer:         cmd.Flag("oauth-issuer", "The url of dex issuer.").Required().String(),
		oauthIssuerCertPath: cmd.Flag("oauth-issuer-cert", "Path to TLS certificate of issuer.").Required().String(),
		emailFilter:         cmd.Flag("email-filter", "Regexp pattern for filtering email.").String(),
	}
	r.RegisterCommand(server)
	return server
}

func (s *server) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()
	// verifier
	verifier, err := token.NewVerifier(*s.oauthPublicKeyPath, *s.oauthIssuer, *s.oauthClientID)
	if err != nil {
		return err
	}
	// healthCheckService
	restHealthChecker := health.NewRestChecker(
		"", "",
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
	)
	go restHealthChecker.Run(ctx)
	// healthcheckService
	healthcheckServer := rest.NewServer(
		*s.certPath, *s.keyPath,
		rest.WithLogger(logger),
		rest.WithService(restHealthChecker),
		rest.WithMetrics(registerer),
	)
	defer healthcheckServer.Stop(10 * time.Second)
	go healthcheckServer.Run()
	// mysqlClient
	mysqlClient, err := s.createMySQLClient(ctx, registerer, logger)
	if err != nil {
		return err
	}
	defer mysqlClient.Close()
	// domainTopicPublisher
	domainTopicPublisher, err := s.createPublisher(ctx, *s.domainTopic, registerer, logger)
	if err != nil {
		return err
	}
	defer domainTopicPublisher.Stop()
	// credential for grpc
	creds, err := client.NewPerRPCCredentials(*s.serviceTokenPath)
	if err != nil {
		return err
	}
	// accountClient
	accountClient, err := accountclient.NewClient(*s.accountService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger))
	if err != nil {
		return err
	}
	defer accountClient.Close()
	// environmentClient
	environmentClient, err := environmentclient.NewClient(*s.environmentService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer environmentClient.Close()
	// authService
	authService, err := s.createAuthService(ctx, accountClient, logger)
	if err != nil {
		return err
	}
	authServer := rpc.NewServer(authService, *s.certPath, *s.keyPath,
		rpc.WithPort(*s.authServicePort),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	defer authServer.Stop(10 * time.Second)
	go authServer.Run()
	// accountService
	accountService := accountapi.NewAccountService(
		environmentClient,
		mysqlClient,
		domainTopicPublisher,
		accountapi.WithLogger(logger),
	)
	accountServer := rpc.NewServer(accountService, *s.certPath, *s.keyPath,
		rpc.WithPort(*s.accountServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	defer accountServer.Stop(10 * time.Second)
	go accountServer.Run()
	// other services...
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

func (s *server) createPublisher(
	ctx context.Context,
	topic string,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (publisher.Publisher, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	client, err := pubsub.NewClient(
		ctx,
		*s.project,
		pubsub.WithMetrics(registerer),
		pubsub.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}
	return client.CreatePublisher(topic)
}

func (s *server) createAuthService(
	ctx context.Context,
	accountClient accountclient.Client,
	logger *zap.Logger,
) (rpc.Service, error) {
	o, err := oidc.NewOIDC(
		ctx,
		*s.oauthIssuer,
		*s.oauthIssuerCertPath,
		*s.oauthClientID,
		*s.oauthClientSecret,
		*s.oauthRedirectURLs,
		oidc.WithLogger(logger))
	if err != nil {
		return nil, err
	}
	signer, err := token.NewSigner(*s.oauthPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	serviceOptions := []authapi.Option{
		authapi.WithLogger(logger),
	}
	if *s.emailFilter != "" {
		filter, err := regexp.Compile(*s.emailFilter)
		if err != nil {
			return nil, err
		}
		serviceOptions = append(serviceOptions, authapi.WithEmailFilter(filter))
	}
	return authapi.NewAuthService(o, signer, accountClient, serviceOptions...), nil
}
