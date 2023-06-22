// Copyright 2023 The Bucketeer Authors.
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
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"time"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	accountapi "github.com/bucketeer-io/bucketeer/pkg/account/api"
	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	auditlogapi "github.com/bucketeer-io/bucketeer/pkg/auditlog/api"
	authapi "github.com/bucketeer-io/bucketeer/pkg/auth/api"
	authclient "github.com/bucketeer-io/bucketeer/pkg/auth/client"
	"github.com/bucketeer-io/bucketeer/pkg/auth/oidc"
	autoopsapi "github.com/bucketeer-io/bucketeer/pkg/autoops/api"
	"github.com/bucketeer-io/bucketeer/pkg/autoops/webhookhandler"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/crypto"
	environmentapi "github.com/bucketeer-io/bucketeer/pkg/environment/api"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	eventcounterapi "github.com/bucketeer-io/bucketeer/pkg/eventcounter/api"
	experimentapi "github.com/bucketeer-io/bucketeer/pkg/experiment/api"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/pkg/rest"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	bqquerier "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/querier"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/token"
)

const (
	command            = "server"
	gcp                = "gcp"
	aws                = "aws"
	autoOpsWebhookPath = "hook"
)

type server struct {
	*kingpin.CmdClause
	project                 *string
	mysqlUser               *string
	mysqlPass               *string
	mysqlHost               *string
	mysqlPort               *int
	mysqlDBName             *string
	redisServerName         *string
	redisAddr               *string
	redisPoolMaxIdle        *int
	redisPoolMaxActive      *int
	bigQueryDataSet         *string
	bigQueryDataLocation    *string
	domainTopic             *string
	accountServicePort      *int
	authServicePort         *int
	auditLogServicePort     *int
	autoOpsServicePort      *int
	environmentServicePort  *int
	eventCounterServicePort *int
	experimentServicePort   *int
	accountService          *string
	authService             *string
	environmentService      *string
	experimentService       *string
	featureService          *string
	timezone                *string
	certPath                *string
	keyPath                 *string
	serviceTokenPath        *string
	oauthPublicKeyPath      *string
	oauthClientID           *string
	oauthIssuer             *string
	// auth
	oauthIssuerCertPath *string
	emailFilter         *string
	oauthRedirectURLs   *[]string
	oauthClientSecret   *string
	oauthPrivateKeyPath *string
	// autoOps
	webhookBaseURL         *string
	webhookKMSResourceName *string
	cloudService           *string
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the server")
	server := &server{
		CmdClause:       cmd,
		project:         cmd.Flag("project", "Google Cloud project name.").String(),
		mysqlUser:       cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:       cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:       cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:       cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName:     cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		redisServerName: cmd.Flag("redis-server-name", "Name of the redis.").Required().String(),
		redisAddr:       cmd.Flag("redis-addr", "Address of the redis.").Required().String(),
		redisPoolMaxIdle: cmd.Flag(
			"redis-pool-max-idle",
			"Maximum number of idle connections in the pool.",
		).Default("5").Int(),
		redisPoolMaxActive: cmd.Flag(
			"redis-pool-max-active",
			"Maximum number of connections allocated by the pool at a given time.",
		).Default("10").Int(),
		bigQueryDataSet:      cmd.Flag("bigquery-data-set", "BigQuery DataSet Name").String(),
		bigQueryDataLocation: cmd.Flag("bigquery-data-location", "BigQuery DataSet Location").String(),
		domainTopic:          cmd.Flag("domain-topic", "PubSub topic to publish domain events.").Required().String(),
		accountServicePort: cmd.Flag(
			"account-service-port",
			"Port to bind to account service.",
		).Default("9091").Int(),
		authServicePort: cmd.Flag(
			"auth-service-port",
			"Port to bind to auth service.",
		).Default("9092").Int(),
		auditLogServicePort: cmd.Flag(
			"audit-log-service-port",
			"Port to bind to audit log service.",
		).Default("9093").Int(),
		autoOpsServicePort: cmd.Flag(
			"auto-ops-service-port",
			"Port to bind to auto ops service.",
		).Default("9094").Int(),
		environmentServicePort: cmd.Flag(
			"environment-service-port",
			"Port to bind to environment service.",
		).Default("9095").Int(),
		eventCounterServicePort: cmd.Flag(
			"event-counter-service-port",
			"Port to bind to event counter service.",
		).Default("9096").Int(),
		experimentServicePort: cmd.Flag(
			"experiment-service-port",
			"Port to bind to experiment service.",
		).Default("9097").Int(),
		accountService: cmd.Flag(
			"account-service",
			"bucketeer-account-service address.",
		).Default("localhost:9001").String(),
		authService: cmd.Flag(
			"auth-service",
			"bucketeer-auth-service address.",
		).Default("localhost:9001").String(),
		environmentService: cmd.Flag(
			"environment-service",
			"bucketeer-environment-service address.",
		).Default("localhost:9001").String(),
		experimentService: cmd.Flag(
			"experiment-service",
			"bucketeer-experiment-service address.",
		).Default("localhost:9001").String(),
		featureService: cmd.Flag(
			"feature-service",
			"bucketeer-feature-service address.",
		).Default("localhost:9001").String(),
		timezone:         cmd.Flag("timezone", "Time zone").Required().String(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		oauthPublicKeyPath: cmd.Flag(
			"oauth-public-key",
			"Path to public key used to verify oauth token.",
		).Required().String(),
		oauthClientID: cmd.Flag(
			"oauth-client-id",
			"The oauth clientID registered at dex.",
		).Required().String(),
		oauthIssuer: cmd.Flag("oauth-issuer", "The url of dex issuer.").Required().String(),
		// auth
		oauthIssuerCertPath: cmd.Flag("oauth-issuer-cert", "Path to TLS certificate of issuer.").Required().String(),
		emailFilter:         cmd.Flag("email-filter", "Regexp pattern for filtering email.").String(),
		oauthRedirectURLs:   cmd.Flag("oauth-redirect-urls", "The redirect urls registered at Dex.").Required().Strings(),
		oauthClientSecret: cmd.Flag(
			"oauth-client-secret",
			"The oauth client secret registered at Dex.",
		).Required().String(),
		oauthPrivateKeyPath: cmd.Flag(
			"oauth-private-key",
			"Path to private key for signing oauth token.",
		).Required().String(),
		// autoOps
		webhookBaseURL: cmd.Flag("webhook-base-url", "the base url for incoming webhooks.").Required().String(),
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
	// redisV3Client
	redisV3Client, err := redisv3.NewClient(
		*s.redisAddr,
		redisv3.WithPoolSize(*s.redisPoolMaxActive),
		redisv3.WithMinIdleConns(*s.redisPoolMaxIdle),
		redisv3.WithServerName(*s.redisServerName),
		redisv3.WithMetrics(registerer),
		redisv3.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer redisV3Client.Close()
	redisV3Cache := cachev3.NewRedisCache(redisV3Client)
	// bigQueryQuerier
	bigQueryQuerier, err := s.createBigQueryQuerier(ctx, *s.project, *s.bigQueryDataLocation, registerer, logger)
	if err != nil {
		logger.Error("Failed to create BigQuery client",
			zap.Error(err),
			zap.String("project", *s.project),
			zap.String("location", *s.bigQueryDataLocation),
			zap.String("data-set", *s.bigQueryDataSet),
		)
		return err
	}
	defer bigQueryQuerier.Close()
	// bigQueryDataSet
	bigQueryDataSet := *s.bigQueryDataSet
	// location
	location, err := locale.GetLocation(*s.timezone)
	if err != nil {
		return err
	}
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
	// authClient
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
	// experimentClient
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
	// featureClient
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
	// auditLogService
	auditLogService := auditlogapi.NewAuditLogService(
		accountClient,
		mysqlClient,
		auditlogapi.WithLogger(logger),
	)
	auditLogServer := rpc.NewServer(auditLogService, *s.certPath, *s.keyPath,
		rpc.WithPort(*s.auditLogServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	defer auditLogServer.Stop(10 * time.Second)
	go auditLogServer.Run()
	// autoOpsService
	autoOpsService, autoOpsWebhookHandler, err := s.createAutoOpsService(
		ctx,
		accountClient,
		authClient,
		experimentClient,
		featureClient,
		mysqlClient,
		domainTopicPublisher,
		verifier,
		logger,
	)
	if err != nil {
		return err
	}
	autoOpsServer := rpc.NewServer(autoOpsService, *s.certPath, *s.keyPath,
		rpc.WithPort(*s.autoOpsServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
		rpc.WithHandler(fmt.Sprintf("/%s", autoOpsWebhookPath), autoOpsWebhookHandler),
	)
	defer autoOpsServer.Stop(10 * time.Second)
	go autoOpsServer.Run()
	// environmentService
	environmentService := environmentapi.NewEnvironmentService(
		accountClient,
		mysqlClient,
		domainTopicPublisher,
		environmentapi.WithLogger(logger),
	)
	environmentServer := rpc.NewServer(environmentService, *s.certPath, *s.keyPath,
		rpc.WithPort(*s.environmentServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	defer environmentServer.Stop(10 * time.Second)
	go environmentServer.Run()
	// eventCounterService
	eventCounterService := eventcounterapi.NewEventCounterService(
		mysqlClient,
		experimentClient,
		featureClient,
		accountClient,
		bigQueryQuerier,
		bigQueryDataSet,
		registerer,
		redisV3Cache,
		location,
		logger,
	)
	eventCounterServer := rpc.NewServer(eventCounterService, *s.certPath, *s.keyPath,
		rpc.WithPort(*s.eventCounterServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	defer eventCounterServer.Stop(10 * time.Second)
	go eventCounterServer.Run()
	// experimentService
	experimentService := experimentapi.NewExperimentService(
		featureClient,
		accountClient,
		mysqlClient,
		domainTopicPublisher,
		experimentapi.WithLogger(logger),
	)
	experimentServer := rpc.NewServer(experimentService, *s.certPath, *s.keyPath,
		rpc.WithPort(*s.experimentServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	defer experimentServer.Stop(10 * time.Second)
	go experimentServer.Run()
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

func (s *server) createBigQueryQuerier(
	ctx context.Context,
	project, location string,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (bqquerier.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return bqquerier.NewClient(
		ctx,
		project,
		location,
		bqquerier.WithMetrics(registerer),
		bqquerier.WithLogger(logger),
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

func (s *server) createAutoOpsService(
	ctx context.Context,
	accountClient accountclient.Client,
	authClient authclient.Client,
	experimentClient experimentclient.Client,
	featureClient featureclient.Client,
	mysqlClient mysql.Client,
	domainTopicPublisher publisher.Publisher,
	verifier token.Verifier,
	logger *zap.Logger,
) (rpc.Service, http.Handler, error) {
	u, err := url.Parse(*s.webhookBaseURL)
	if err != nil {
		return nil, nil, err
	}
	u.Path = path.Join(u.Path, autoOpsWebhookPath)

	var webhookCryptoUtil crypto.EncrypterDecrypter
	switch *s.cloudService {
	case gcp:
		webhookCryptoUtil, err = crypto.NewCloudKMSCrypto(ctx, *s.webhookKMSResourceName)
		if err != nil {
			return nil, nil, err
		}
	case aws:
		// TODO: Get region from command-line flags
		webhookCryptoUtil, err = crypto.NewAwsKMSCrypto(ctx, *s.webhookKMSResourceName, "ap-northeast-1")
		if err != nil {
			return nil, nil, err
		}
	}
	autoOpsService := autoopsapi.NewAutoOpsService(
		mysqlClient,
		featureClient,
		experimentClient,
		accountClient,
		authClient,
		domainTopicPublisher,
		u,
		webhookCryptoUtil,
		autoopsapi.WithLogger(logger),
	)
	autoOpsWebhookHandler, err := webhookhandler.NewHandler(
		mysqlClient,
		authClient,
		featureClient,
		domainTopicPublisher,
		verifier,
		*s.serviceTokenPath,
		webhookCryptoUtil,
		webhookhandler.WithLogger(logger),
	)
	if err != nil {
		return nil, nil, err
	}
	return autoOpsService, autoOpsWebhookHandler, nil
}
