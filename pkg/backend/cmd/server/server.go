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
	autoopsclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	"github.com/bucketeer-io/bucketeer/pkg/autoops/webhookhandler"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/crypto"
	environmentapi "github.com/bucketeer-io/bucketeer/pkg/environment/api"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	eventcounterapi "github.com/bucketeer-io/bucketeer/pkg/eventcounter/api"
	experimentapi "github.com/bucketeer-io/bucketeer/pkg/experiment/api"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureapi "github.com/bucketeer-io/bucketeer/pkg/feature/api"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	migratemysqlapi "github.com/bucketeer-io/bucketeer/pkg/migration/mysql/api"
	"github.com/bucketeer-io/bucketeer/pkg/migration/mysql/migrate"
	notificationapi "github.com/bucketeer-io/bucketeer/pkg/notification/api"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	pushapi "github.com/bucketeer-io/bucketeer/pkg/push/api"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/pkg/rest"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	bqquerier "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/querier"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/token"
)

const (
	command               = "server"
	gcp                   = "gcp"
	aws                   = "aws"
	autoOpsWebhookPath    = "hook"
	healthCheckTimeout    = 1 * time.Second
	clientDialTimeout     = 30 * time.Second
	serverShutDownTimeout = 10 * time.Second
)

type server struct {
	*kingpin.CmdClause
	// Common
	project            *string
	timezone           *string
	certPath           *string
	keyPath            *string
	serviceTokenPath   *string
	oauthPublicKeyPath *string
	oauthClientID      *string
	oauthIssuer        *string
	// MySQL
	mysqlUser   *string
	mysqlPass   *string
	mysqlHost   *string
	mysqlPort   *int
	mysqlDBName *string
	// MySQL for Migration
	mysqlMigrationUser *string
	mysqlMigrationPass *string
	// Persistent Redis
	persistentRedisServerName    *string
	persistentRedisAddr          *string
	persistentRedisPoolMaxIdle   *int
	persistentRedisPoolMaxActive *int
	// Non Persistent Redis
	nonPersistentRedisServerName    *string
	nonPersistentRedisAddr          *string
	nonPersistentRedisPoolMaxIdle   *int
	nonPersistentRedisPoolMaxActive *int
	// BigQuery
	bigQueryDataSet      *string
	bigQueryDataLocation *string
	// PubSub
	domainTopic                   *string
	bulkSegmentUsersReceivedTopic *string
	// Port
	accountServicePort      *int
	authServicePort         *int
	auditLogServicePort     *int
	autoOpsServicePort      *int
	environmentServicePort  *int
	eventCounterServicePort *int
	experimentServicePort   *int
	featureServicePort      *int
	migrateMySQLServicePort *int
	notificationServicePort *int
	pushServicePort         *int
	// Service
	accountService     *string
	authService        *string
	environmentService *string
	experimentService  *string
	featureService     *string
	autoOpsService     *string
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
	// migration-mysql
	githubUser                *string
	githubAccessTokenPath     *string
	githubMigrationSourcePath *string
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
		persistentRedisServerName: cmd.Flag(
			"persistent-redis-server-name",
			"Name of the persistent redis.",
		).Required().String(),
		persistentRedisAddr: cmd.Flag(
			"persistent-redis-addr",
			"Address of the persistent redis.",
		).Required().String(),
		persistentRedisPoolMaxIdle: cmd.Flag(
			"persistent-redis-pool-max-idle",
			"Maximum number of idle in the persistent redis connections pool.",
		).Default("5").Int(),
		persistentRedisPoolMaxActive: cmd.Flag(
			"persistent-redis-pool-max-active",
			"Maximum number of connections allocated by the persistent redis connections pool at a given time.",
		).Default("10").Int(),
		nonPersistentRedisServerName: cmd.Flag(
			"non-persistent-redis-server-name",
			"Name of the non-persistent redis.",
		).Required().String(),
		nonPersistentRedisAddr: cmd.Flag(
			"non-persistent-redis-addr",
			"Address of the non-persistent redis.",
		).Required().String(),
		nonPersistentRedisPoolMaxIdle: cmd.Flag(
			"non-persistent-redis-pool-max-idle",
			"Maximum number of idle in the non-persistent redis connections pool.",
		).Default("5").Int(),
		nonPersistentRedisPoolMaxActive: cmd.Flag(
			"non-persistent-redis-pool-max-active",
			"Maximum number of connections allocated by the non-persistent redis connections pool at a given time.",
		).Default("10").Int(),
		bigQueryDataSet:      cmd.Flag("bigquery-data-set", "BigQuery DataSet Name").String(),
		bigQueryDataLocation: cmd.Flag("bigquery-data-location", "BigQuery DataSet Location").String(),
		domainTopic: cmd.Flag(
			"domain-topic",
			"PubSub topic to publish domain events.",
		).Required().String(),
		bulkSegmentUsersReceivedTopic: cmd.Flag(
			"bulk-segment-users-received-topic",
			"PubSub topic to publish bulk segment users received events.",
		).Required().String(),
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
		featureServicePort: cmd.Flag(
			"feature-service-port",
			"Port to bind to feature service.",
		).Default("9098").Int(),
		migrateMySQLServicePort: cmd.Flag(
			"migrate-mysql-service-port",
			"Port to bind to migrate mysql service.",
		).Default("9099").Int(),
		notificationServicePort: cmd.Flag(
			"notification-service-port",
			"Port to bind to notification service.",
		).Default("9100").Int(),
		pushServicePort: cmd.Flag(
			"push-service-port",
			"Port to bind to push service.",
		).Default("9101").Int(),
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
		autoOpsService: cmd.Flag(
			"autoops-service",
			"bucketeer-autoops-service address.",
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
		// migration-mysql
		githubUser:            cmd.Flag("github-user", "GitHub user.").Required().String(),
		githubAccessTokenPath: cmd.Flag("github-access-token-path", "Path to GitHub access token.").Required().String(),
		githubMigrationSourcePath: cmd.Flag(
			"github-migration-source-path",
			"Path to migration file in GitHub. (e.g. owner/repo/path#ref)",
		).Required().String(),
		mysqlMigrationUser: cmd.Flag("mysql-migration-user", "MySQL user for migration.").Required().String(),
		mysqlMigrationPass: cmd.Flag("mysql-migration-pass", "MySQL password for migration.").Required().String(),
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
		health.WithTimeout(healthCheckTimeout),
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
	go healthcheckServer.Run()
	// mysqlClient
	mysqlClient, err := s.createMySQLClient(ctx, registerer, logger)
	if err != nil {
		return err
	}
	// persistentRedisClient
	persistentRedisClient, err := redisv3.NewClient(
		*s.persistentRedisAddr,
		redisv3.WithPoolSize(*s.persistentRedisPoolMaxActive),
		redisv3.WithMinIdleConns(*s.persistentRedisPoolMaxIdle),
		redisv3.WithServerName(*s.persistentRedisServerName),
		redisv3.WithMetrics(registerer),
		redisv3.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	persistentRedisV3Cache := cachev3.NewRedisCache(persistentRedisClient)
	// nonPersistentRedisClient
	nonPersistentRedisClient, err := redisv3.NewClient(
		*s.nonPersistentRedisAddr,
		redisv3.WithPoolSize(*s.nonPersistentRedisPoolMaxActive),
		redisv3.WithMinIdleConns(*s.nonPersistentRedisPoolMaxIdle),
		redisv3.WithServerName(*s.nonPersistentRedisServerName),
		redisv3.WithMetrics(registerer),
		redisv3.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	nonPersistentRedisV3Cache := cachev3.NewRedisCache(nonPersistentRedisClient)
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
	// segmentUsersPublisher
	segmentUsersPublisher, err := s.createPublisher(ctx, *s.bulkSegmentUsersReceivedTopic, registerer, logger)
	if err != nil {
		return err
	}
	// credential for grpc
	creds, err := client.NewPerRPCCredentials(*s.serviceTokenPath)
	if err != nil {
		return err
	}
	// accountClient
	accountClient, err := accountclient.NewClient(*s.accountService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(clientDialTimeout),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger))
	if err != nil {
		return err
	}
	// authClient
	authClient, err := authclient.NewClient(*s.authService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(clientDialTimeout),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	// environmentClient
	environmentClient, err := environmentclient.NewClient(*s.environmentService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(clientDialTimeout),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	// experimentClient
	experimentClient, err := experimentclient.NewClient(*s.experimentService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(clientDialTimeout),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	// featureClient
	featureClient, err := featureclient.NewClient(*s.featureService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(clientDialTimeout),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer featureClient.Close()
	// autoOpsClient
	autoOpsClient, err := autoopsclient.NewClient(*s.autoOpsService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer autoOpsClient.Close()
	// authService
	authService, err := s.createAuthService(ctx, accountClient, logger)
	if err != nil {
		return err
	}
	authServer := rpc.NewServer(authService, *s.certPath, *s.keyPath,
		"auth-server",
		rpc.WithPort(*s.authServicePort),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	go authServer.Run()
	// accountService
	accountService := accountapi.NewAccountService(
		environmentClient,
		mysqlClient,
		domainTopicPublisher,
		accountapi.WithLogger(logger),
	)
	accountServer := rpc.NewServer(accountService, *s.certPath, *s.keyPath,
		"account-server",
		rpc.WithPort(*s.accountServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	go accountServer.Run()
	// auditLogService
	auditLogService := auditlogapi.NewAuditLogService(
		accountClient,
		mysqlClient,
		auditlogapi.WithLogger(logger),
	)
	auditLogServer := rpc.NewServer(auditLogService, *s.certPath, *s.keyPath,
		"audit-log-server",
		rpc.WithPort(*s.auditLogServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
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
		"auto-ops-server",
		rpc.WithPort(*s.autoOpsServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
		rpc.WithHandler(fmt.Sprintf("/%s", autoOpsWebhookPath), autoOpsWebhookHandler),
	)
	go autoOpsServer.Run()
	// environmentService
	environmentService := environmentapi.NewEnvironmentService(
		accountClient,
		mysqlClient,
		domainTopicPublisher,
		environmentapi.WithLogger(logger),
	)
	environmentServer := rpc.NewServer(environmentService, *s.certPath, *s.keyPath,
		"environment-server",
		rpc.WithPort(*s.environmentServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
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
		persistentRedisV3Cache,
		location,
		logger,
	)
	eventCounterServer := rpc.NewServer(eventCounterService, *s.certPath, *s.keyPath,
		"event-counter-server",
		rpc.WithPort(*s.eventCounterServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
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
		"experiment-server",
		rpc.WithPort(*s.experimentServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	go experimentServer.Run()
	// featureService
	featureService := featureapi.NewFeatureService(
		mysqlClient,
		accountClient,
		experimentClient,
		autoOpsClient,
		nonPersistentRedisV3Cache,
		segmentUsersPublisher,
		domainTopicPublisher,
		featureapi.WithLogger(logger),
	)
	featureServer := rpc.NewServer(featureService, *s.certPath, *s.keyPath,
		"feature-server",
		rpc.WithPort(*s.featureServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	go featureServer.Run()
	// migrateMySQLService
	migrateClientFactory, err := migrate.NewClientFactory(
		*s.githubUser, *s.githubAccessTokenPath, *s.githubMigrationSourcePath,
		*s.mysqlMigrationUser, *s.mysqlMigrationPass, *s.mysqlHost, *s.mysqlPort, *s.mysqlDBName,
	)
	if err != nil {
		return err
	}
	migrateMySQLService := migratemysqlapi.NewMySQLService(
		migrateClientFactory,
		migratemysqlapi.WithLogger(logger),
	)
	migrateMySQLServer := rpc.NewServer(migrateMySQLService, *s.certPath, *s.keyPath,
		"migrate-mysql-server",
		rpc.WithPort(*s.migrateMySQLServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	go migrateMySQLServer.Run()
	// notificationService
	notificationService := notificationapi.NewNotificationService(
		mysqlClient,
		accountClient,
		domainTopicPublisher,
		notificationapi.WithLogger(logger),
	)
	notificationServer := rpc.NewServer(notificationService, *s.certPath, *s.keyPath,
		"notification-server",
		rpc.WithPort(*s.notificationServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	go notificationServer.Run()
	// pushService
	pushService := pushapi.NewPushService(
		mysqlClient,
		featureClient,
		experimentClient,
		accountClient,
		domainTopicPublisher,
		pushapi.WithLogger(logger),
	)
	pushServer := rpc.NewServer(pushService, *s.certPath, *s.keyPath,
		"push-server",
		rpc.WithPort(*s.pushServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	go pushServer.Run()
	// To detach this pod from Kubernetes Service before the app servers stop, we stop the health check service first.
	// Then, after 10 seconds of sleep, the app servers can be shut down, as no new requests are expected to be sent.
	// In this case, the Readiness prove must fail within 10 seconds and the pod must be detached.
	defer func() {
		go healthcheckServer.Stop(serverShutDownTimeout)
		time.Sleep(serverShutDownTimeout)
		go authServer.Stop(serverShutDownTimeout)
		go accountServer.Stop(serverShutDownTimeout)
		go auditLogServer.Stop(serverShutDownTimeout)
		go autoOpsServer.Stop(serverShutDownTimeout)
		go environmentServer.Stop(serverShutDownTimeout)
		go experimentServer.Stop(serverShutDownTimeout)
		go eventCounterServer.Stop(serverShutDownTimeout)
		go featureServer.Stop(serverShutDownTimeout)
		go migrateMySQLServer.Stop(serverShutDownTimeout)
		go notificationServer.Stop(serverShutDownTimeout)
		go pushServer.Stop(serverShutDownTimeout)
		go mysqlClient.Close()
		go persistentRedisClient.Close()
		go nonPersistentRedisClient.Close()
		go bigQueryQuerier.Close()
		go domainTopicPublisher.Stop()
		go segmentUsersPublisher.Stop()
		go accountClient.Close()
		go authClient.Close()
		go environmentClient.Close()
		go experimentClient.Close()
		go featureClient.Close()
	}()
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
