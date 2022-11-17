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
	"net/url"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	authclient "github.com/bucketeer-io/bucketeer/pkg/auth/client"
	"github.com/bucketeer-io/bucketeer/pkg/autoops/api"
	"github.com/bucketeer-io/bucketeer/pkg/autoops/webhookhandler"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/crypto"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/token"
)

const (
	command     = "server"
	webhookPath = "hook"
	gcp         = "gcp"
	aws         = "aws"
)

type server struct {
	*kingpin.CmdClause
	port              *int
	project           *string
	mysqlUser         *string
	mysqlPass         *string
	mysqlHost         *string
	mysqlPort         *int
	mysqlDBName       *string
	domainEventTopic  *string
	accountService    *string
	authService       *string
	featureService    *string
	experimentService *string
	certPath          *string
	keyPath           *string
	serviceTokenPath  *string

	oauthKeyPath  *string
	oauthClientID *string
	oauthIssuer   *string

	webhookBaseURL         *string
	webhookKMSResourceName *string

	cloudService *string
}

func RegisterServerCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the gRPC server")
	server := &server{
		CmdClause:        cmd,
		port:             cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:          cmd.Flag("project", "Google Cloud project name.").Required().String(),
		mysqlUser:        cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:        cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:        cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:        cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName:      cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		domainEventTopic: cmd.Flag("domain-event-topic", "PubSub topic to publish domain events.").Required().String(),
		accountService:   cmd.Flag("account-service", "bucketeer-account-service address.").Default("account:9090").String(),
		authService:      cmd.Flag("auth-service", "bucketeer-auth-service address.").Default("auth:9090").String(),
		featureService:   cmd.Flag("feature-service", "bucketeer-feature-service address.").Default("feature:9090").String(),
		experimentService: cmd.Flag(
			"experiment-service",
			"bucketeer-experiment-service address.",
		).Default("experiment:9090").String(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		oauthKeyPath:     cmd.Flag("oauth-key", "Path to public key used to verify oauth token.").Required().String(),
		oauthClientID:    cmd.Flag("oauth-client-id", "The oauth clientID registered at dex.").Required().String(),
		oauthIssuer:      cmd.Flag("oauth-issuer", "The url of dex issuer.").Required().String(),
		webhookBaseURL:   cmd.Flag("webhook-base-url", "the base url for incoming webhooks.").Required().String(),
		webhookKMSResourceName: cmd.Flag(
			"webhook-kms-resource-name",
			"Cloud KMS resource name to encrypt and decrypt webhook credentials.",
		).Required().String(),
		cloudService: cmd.Flag("cloud-service", "Cloud Service info").Default(gcp).String(),
	}
	r.RegisterCommand(server)
	return server
}

func (s *server) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()

	*s.serviceTokenPath = s.insertTelepresenceMoutRoot(*s.serviceTokenPath)
	*s.oauthKeyPath = s.insertTelepresenceMoutRoot(*s.oauthKeyPath)
	*s.keyPath = s.insertTelepresenceMoutRoot(*s.keyPath)
	*s.certPath = s.insertTelepresenceMoutRoot(*s.certPath)

	mysqlClient, err := s.createMySQLClient(ctx, registerer, logger)
	if err != nil {
		return err
	}
	defer mysqlClient.Close()

	publisher, err := s.createDomainEventPublisher(ctx, registerer, logger)
	if err != nil {
		return err
	}
	defer publisher.Stop()

	creds, err := client.NewPerRPCCredentials(*s.serviceTokenPath)
	if err != nil {
		return err
	}

	featureClient, err := featureclient.NewClient(*s.featureService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer featureClient.Close()

	experimentClient, err := experimentclient.NewClient(*s.experimentService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer experimentClient.Close()

	accountClient, err := accountclient.NewClient(*s.accountService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer accountClient.Close()

	authClient, err := authclient.NewClient(*s.authService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer authClient.Close()

	u, err := url.Parse(*s.webhookBaseURL)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, webhookPath)

	var webhookCryptoUtil crypto.EncrypterDecrypter
	switch *s.cloudService {
	case gcp:
		webhookCryptoUtil, err = crypto.NewCloudKMSCrypto(ctx, *s.webhookKMSResourceName)
		if err != nil {
			return err
		}
	case aws:
		// TODO: Get region from command-line flags
		webhookCryptoUtil, err = crypto.NewAwsKMSCrypto(ctx, *s.webhookKMSResourceName, "ap-northeast-1")
		if err != nil {
			return err
		}
	}

	service := api.NewAutoOpsService(
		mysqlClient,
		featureClient,
		experimentClient,
		accountClient,
		authClient,
		publisher,
		u,
		webhookCryptoUtil,
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

	webhookHandler, err := webhookhandler.NewHandler(
		mysqlClient,
		authClient,
		featureClient,
		publisher,
		verifier,
		*s.serviceTokenPath,
		webhookCryptoUtil,
		webhookhandler.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	server := rpc.NewServer(service, *s.certPath, *s.keyPath,
		rpc.WithPort(*s.port),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
		rpc.WithService(healthChecker),
		rpc.WithHandler("/health", healthChecker),
		rpc.WithHandler("/hook", webhookHandler),
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

func (s *server) createDomainEventPublisher(
	ctx context.Context,
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
	domainPublisher, err := client.CreatePublisher(*s.domainEventTopic)
	if err != nil {
		return nil, err
	}
	return domainPublisher, nil
}

func (s *server) insertTelepresenceMoutRoot(path string) string {
	volumeRoot := os.Getenv("TELEPRESENCE_ROOT")
	if volumeRoot == "" {
		return path
	}
	return volumeRoot + path
}
