// Copyright 2025 The Bucketeer Authors.
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
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	accountapi "github.com/bucketeer-io/bucketeer/pkg/account/api"
	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	auditlogapi "github.com/bucketeer-io/bucketeer/pkg/auditlog/api"
	"github.com/bucketeer-io/bucketeer/pkg/auth"
	authapi "github.com/bucketeer-io/bucketeer/pkg/auth/api"
	authclient "github.com/bucketeer-io/bucketeer/pkg/auth/client"
	autoopsapi "github.com/bucketeer-io/bucketeer/pkg/autoops/api"
	autoopsclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	btclient "github.com/bucketeer-io/bucketeer/pkg/batch/client"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	coderefapi "github.com/bucketeer-io/bucketeer/pkg/coderef/api"
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
	notificationapi "github.com/bucketeer-io/bucketeer/pkg/notification/api"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	pushapi "github.com/bucketeer-io/bucketeer/pkg/push/api"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/pkg/rest"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	gatewayapi "github.com/bucketeer-io/bucketeer/pkg/rpc/gateway"
	bqquerier "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/querier"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	tagapi "github.com/bucketeer-io/bucketeer/pkg/tag/api"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	auditlogproto "github.com/bucketeer-io/bucketeer/proto/auditlog"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	coderefproto "github.com/bucketeer-io/bucketeer/proto/coderef"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventcounterproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
	pushproto "github.com/bucketeer-io/bucketeer/proto/push"
	tagproto "github.com/bucketeer-io/bucketeer/proto/tag"
)

const (
	command                       = "server"
	gcp                           = "gcp"
	featureFlagTriggerWebhookPath = "webhook/triggers"
	healthCheckTimeout            = 1 * time.Second
	clientDialTimeout             = 30 * time.Second
	serverShutDownTimeout         = 10 * time.Second
)

type server struct {
	*kingpin.CmdClause
	// Common
	project          *string
	timezone         *string
	certPath         *string
	keyPath          *string
	serviceTokenPath *string
	// MySQL
	mysqlUser   *string
	mysqlPass   *string
	mysqlHost   *string
	mysqlPort   *int
	mysqlDBName *string
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
	accountServicePort       *int
	authServicePort          *int
	auditLogServicePort      *int
	autoOpsServicePort       *int
	environmentServicePort   *int
	eventCounterServicePort  *int
	experimentServicePort    *int
	featureServicePort       *int
	notificationServicePort  *int
	pushServicePort          *int
	webConsoleServicePort    *int
	dashboardServicePort     *int
	tagServicePort           *int
	codeReferenceServicePort *int
	// Service
	accountService       *string
	authService          *string
	batchService         *string
	environmentService   *string
	experimentService    *string
	featureService       *string
	autoOpsService       *string
	codeReferenceService *string
	// auth
	refreshTokenTTL     *time.Duration
	emailFilter         *string
	oauthConfigPath     *string
	oauthPublicKeyPath  *string
	oauthPrivateKeyPath *string
	// autoOps
	webhookBaseURL         *string
	webhookKMSResourceName *string
	cloudService           *string
	// web console
	webConsoleEnvJSPath *string
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
		notificationServicePort: cmd.Flag(
			"notification-service-port",
			"Port to bind to notification service.",
		).Default("9100").Int(),
		pushServicePort: cmd.Flag(
			"push-service-port",
			"Port to bind to push service.",
		).Default("9101").Int(),
		webConsoleServicePort: cmd.Flag(
			"web-console-service-port",
			"Port to bind to console service.",
		).Default("9102").Int(),
		dashboardServicePort: cmd.Flag(
			"dashboard-service-port",
			"Port to bind to dashboard service.",
		).Default("9103").Int(),
		tagServicePort: cmd.Flag(
			"tag-service-port",
			"Port to bind to tag service.",
		).Default("9104").Int(),
		codeReferenceServicePort: cmd.Flag(
			"code-reference-service-port",
			"Port to bind to code reference service.",
		).Default("9105").Int(),
		accountService: cmd.Flag(
			"account-service",
			"bucketeer-account-service address.",
		).Default("localhost:9001").String(),
		authService: cmd.Flag(
			"auth-service",
			"bucketeer-auth-service address.",
		).Default("localhost:9001").String(),
		batchService: cmd.Flag(
			"batch-service",
			"bucketeer-batch-service address.",
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
		codeReferenceService: cmd.Flag(
			"code-reference-service",
			"bucketeer-code-reference-service address.",
		).Default("localhost:9001").String(),
		timezone:         cmd.Flag("timezone", "Time zone").Required().String(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		oauthPublicKeyPath: cmd.Flag(
			"oauth-public-key",
			"Path to public key used to verify oauth token.",
		).Required().String(),
		// auth
		refreshTokenTTL: cmd.Flag(
			"refresh-token-ttl",
			"TTL for refresh token.",
		).Default("168h").Duration(),
		emailFilter:     cmd.Flag("email-filter", "Regexp pattern for filtering email.").String(),
		oauthConfigPath: cmd.Flag("oauth-config-path", "Path to oauth config.").Required().String(),
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
		cloudService:        cmd.Flag("cloud-service", "Cloud Service info").Default(gcp).String(),
		webConsoleEnvJSPath: cmd.Flag("web-console-env-js-path", "console env js path").Required().String(),
	}
	r.RegisterCommand(server)
	return server
}

func (s *server) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()

	// oauth config
	oAuthConfig, err := s.readOAuthConfig(ctx, logger)
	if err != nil {
		logger.Error("Failed to read OAuth config", zap.Error(err))
		return err
	}

	// verifier
	// TODO: refactor to support multiple issuers
	verifier, err := token.NewVerifier(
		*s.oauthPublicKeyPath,
		oAuthConfig.Issuer,
		oAuthConfig.Audience,
	)
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
	// batchClient
	batchClient, err := btclient.NewClient(*s.batchService, *s.certPath,
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
	authService, err := s.createAuthService(mysqlClient, accountClient, verifier, oAuthConfig, logger)
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

	// Account REST Gateway
	accountRestAddr := fmt.Sprintf(":%d", *s.accountServicePort+1000) // REST on port 10091
	accountGrpcAddr := fmt.Sprintf("localhost:%d", *s.accountServicePort)
	accountGateway, err := gatewayapi.NewGateway(
		accountGrpcAddr,
		accountRestAddr,
		gatewayapi.WithLogger(logger.Named("account-gateway")),
		gatewayapi.WithMetrics(registerer),
	)
	if err != nil {
		return fmt.Errorf("failed to create account gateway: %v", err)
	}

	go func() {
		if err := accountGateway.Start(
			ctx,
			accountproto.RegisterAccountServiceHandlerFromEndpoint,
		); err != nil {
			logger.Error("failed to start account gateway", zap.Error(err))
		}
	}()

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

	// Audit Log REST Gateway
	auditLogRestAddr := fmt.Sprintf(":%d", *s.auditLogServicePort+1000) // REST on port 10093
	auditLogGrpcAddr := fmt.Sprintf("localhost:%d", *s.auditLogServicePort)
	auditLogGateway, err := gatewayapi.NewGateway(
		auditLogGrpcAddr,
		auditLogRestAddr,
		gatewayapi.WithLogger(logger.Named("auditlog-gateway")),
		gatewayapi.WithMetrics(registerer),
	)
	if err != nil {
		return fmt.Errorf("failed to create audit log gateway: %v", err)
	}

	go func() {
		if err := auditLogGateway.Start(
			ctx,
			auditlogproto.RegisterAuditLogServiceHandlerFromEndpoint,
		); err != nil {
			logger.Error("failed to start audit log gateway", zap.Error(err))
		}
	}()

	// autoOpsService
	autoOpsService := autoopsapi.NewAutoOpsService(
		mysqlClient,
		featureClient,
		experimentClient,
		accountClient,
		authClient,
		domainTopicPublisher,
		autoopsapi.WithLogger(logger),
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
	)
	go autoOpsServer.Run()

	// Auto Ops REST Gateway
	autoOpsRestAddr := fmt.Sprintf(":%d", *s.autoOpsServicePort+1000) // REST on port 10094
	autoOpsGrpcAddr := fmt.Sprintf("localhost:%d", *s.autoOpsServicePort)
	autoOpsGateway, err := gatewayapi.NewGateway(
		autoOpsGrpcAddr,
		autoOpsRestAddr,
		gatewayapi.WithLogger(logger.Named("autoops-gateway")),
		gatewayapi.WithMetrics(registerer),
	)
	if err != nil {
		return fmt.Errorf("failed to create auto ops gateway: %v", err)
	}

	go func() {
		if err := autoOpsGateway.Start(
			ctx,
			autoopsproto.RegisterAutoOpsServiceHandlerFromEndpoint,
		); err != nil {
			logger.Error("failed to start auto ops gateway", zap.Error(err))
		}
	}()

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

	// Environment REST Gateway
	environmentRestAddr := fmt.Sprintf(":%d", *s.environmentServicePort+1000) // REST on port 10095
	environmentGrpcAddr := fmt.Sprintf("localhost:%d", *s.environmentServicePort)
	environmentGateway, err := gatewayapi.NewGateway(
		environmentGrpcAddr,
		environmentRestAddr,
		gatewayapi.WithLogger(logger.Named("environment-gateway")),
		gatewayapi.WithMetrics(registerer),
	)
	if err != nil {
		return fmt.Errorf("failed to create environment gateway: %v", err)
	}

	go func() {
		if err := environmentGateway.Start(
			ctx,
			environmentproto.RegisterEnvironmentServiceHandlerFromEndpoint,
		); err != nil {
			logger.Error("failed to start environment gateway", zap.Error(err))
		}
	}()

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

	// Event Counter REST Gateway
	eventCounterRestAddr := fmt.Sprintf(":%d", *s.eventCounterServicePort+1000) // REST on port 10096
	eventCounterGrpcAddr := fmt.Sprintf("localhost:%d", *s.eventCounterServicePort)
	eventCounterGateway, err := gatewayapi.NewGateway(
		eventCounterGrpcAddr,
		eventCounterRestAddr,
		gatewayapi.WithLogger(logger.Named("eventcounter-gateway")),
		gatewayapi.WithMetrics(registerer),
	)
	if err != nil {
		return fmt.Errorf("failed to create event counter gateway: %v", err)
	}

	go func() {
		if err := eventCounterGateway.Start(
			ctx,
			eventcounterproto.RegisterEventCounterServiceHandlerFromEndpoint,
		); err != nil {
			logger.Error("failed to start event counter gateway", zap.Error(err))
		}
	}()

	// experimentService
	experimentService := experimentapi.NewExperimentService(
		featureClient,
		accountClient,
		autoOpsClient,
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

	// Experiment REST Gateway
	experimentRestAddr := fmt.Sprintf(":%d", *s.experimentServicePort+1000) // REST on port 10097
	experimentGrpcAddr := fmt.Sprintf("localhost:%d", *s.experimentServicePort)
	experimentGateway, err := gatewayapi.NewGateway(
		experimentGrpcAddr,
		experimentRestAddr,
		gatewayapi.WithLogger(logger.Named("experiment-gateway")),
		gatewayapi.WithMetrics(registerer),
	)
	if err != nil {
		return fmt.Errorf("failed to create experiment gateway: %v", err)
	}

	go func() {
		if err := experimentGateway.Start(
			ctx,
			experimentproto.RegisterExperimentServiceHandlerFromEndpoint,
		); err != nil {
			logger.Error("failed to start experiment gateway", zap.Error(err))
		}
	}()

	// featureService
	featureService, err := s.createFeatureService(
		ctx,
		accountClient,
		experimentClient,
		autoOpsClient,
		batchClient,
		environmentClient,
		nonPersistentRedisV3Cache,
		segmentUsersPublisher,
		domainTopicPublisher,
		mysqlClient,
		logger,
	)
	if err != nil {
		return err
	}
	featureServer := rpc.NewServer(featureService, *s.certPath, *s.keyPath,
		"feature-server",
		rpc.WithPort(*s.featureServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	go featureServer.Run()

	// Feature REST Gateway
	featureRestAddr := fmt.Sprintf(":%d", *s.featureServicePort+1000) // REST on port 10098
	featureGrpcAddr := fmt.Sprintf("localhost:%d", *s.featureServicePort)
	featureGateway, err := gatewayapi.NewGateway(
		featureGrpcAddr,
		featureRestAddr,
		gatewayapi.WithLogger(logger.Named("feature-gateway")),
		gatewayapi.WithMetrics(registerer),
	)
	if err != nil {
		return fmt.Errorf("failed to create feature gateway: %v", err)
	}

	go func() {
		if err := featureGateway.Start(
			ctx,
			featureproto.RegisterFeatureServiceHandlerFromEndpoint,
		); err != nil {
			logger.Error("failed to start feature gateway", zap.Error(err))
		}
	}()

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

	// Notification REST Gateway
	notificationRestAddr := fmt.Sprintf(":%d", *s.notificationServicePort+1000) // REST on port 10100
	notificationGrpcAddr := fmt.Sprintf("localhost:%d", *s.notificationServicePort)
	notificationGateway, err := gatewayapi.NewGateway(
		notificationGrpcAddr,
		notificationRestAddr,
		gatewayapi.WithLogger(logger.Named("notification-gateway")),
		gatewayapi.WithMetrics(registerer),
	)
	if err != nil {
		return fmt.Errorf("failed to create notification gateway: %v", err)
	}

	go func() {
		if err := notificationGateway.Start(
			ctx,
			notificationproto.RegisterNotificationServiceHandlerFromEndpoint,
		); err != nil {
			logger.Error("failed to start notification gateway", zap.Error(err))
		}
	}()

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

	// Push REST Gateway
	pushRestAddr := fmt.Sprintf(":%d", *s.pushServicePort+1000) // REST on port 10101
	pushGrpcAddr := fmt.Sprintf("localhost:%d", *s.pushServicePort)
	pushGateway, err := gatewayapi.NewGateway(
		pushGrpcAddr,
		pushRestAddr,
		gatewayapi.WithLogger(logger.Named("push-gateway")),
		gatewayapi.WithMetrics(registerer),
	)
	if err != nil {
		return fmt.Errorf("failed to create push gateway: %v", err)
	}

	go func() {
		if err := pushGateway.Start(
			ctx,
			pushproto.RegisterPushServiceHandlerFromEndpoint,
		); err != nil {
			logger.Error("failed to start push gateway", zap.Error(err))
		}
	}()

	// tagService
	tagService := tagapi.NewTagService(
		mysqlClient,
		accountClient,
		domainTopicPublisher,
		tagapi.WithLogger(logger),
	)
	tagServer := rpc.NewServer(tagService, *s.certPath, *s.keyPath,
		"tag-server",
		rpc.WithPort(*s.tagServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	go tagServer.Run()

	// Tag REST Gateway
	tagRestAddr := fmt.Sprintf(":%d", *s.tagServicePort+1000) // REST on port 10104
	tagGrpcAddr := fmt.Sprintf("localhost:%d", *s.tagServicePort)
	tagGateway, err := gatewayapi.NewGateway(
		tagGrpcAddr,
		tagRestAddr,
		gatewayapi.WithLogger(logger.Named("tag-gateway")),
		gatewayapi.WithMetrics(registerer),
	)
	if err != nil {
		return fmt.Errorf("failed to create tag gateway: %v", err)
	}

	go func() {
		if err := tagGateway.Start(
			ctx,
			tagproto.RegisterTagServiceHandlerFromEndpoint,
		); err != nil {
			logger.Error("failed to start tag gateway", zap.Error(err))
		}
	}()

	// codeReferenceService
	codeReferenceService := coderefapi.NewCodeReferenceService(
		accountClient,
		mysqlClient,
		domainTopicPublisher,
		coderefapi.WithLogger(logger),
	)
	codeReferenceServer := rpc.NewServer(codeReferenceService, *s.certPath, *s.keyPath,
		"code-reference-server",
		rpc.WithPort(*s.codeReferenceServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	go codeReferenceServer.Run()

	// Code Reference REST Gateway
	codeRefRestAddr := fmt.Sprintf(":%d", *s.codeReferenceServicePort+1000) // REST on port 10105
	codeRefGrpcAddr := fmt.Sprintf("localhost:%d", *s.codeReferenceServicePort)
	codeRefGateway, err := gatewayapi.NewGateway(
		codeRefGrpcAddr,
		codeRefRestAddr,
		gatewayapi.WithLogger(logger.Named("coderef-gateway")),
		gatewayapi.WithMetrics(registerer),
	)
	if err != nil {
		return fmt.Errorf("failed to create code reference gateway: %v", err)
	}

	go func() {
		if err := codeRefGateway.Start(
			ctx,
			coderefproto.RegisterCodeReferenceServiceHandlerFromEndpoint,
		); err != nil {
			logger.Error("failed to start code reference gateway", zap.Error(err))
		}
	}()

	webConsoleServer := rest.NewServer(
		*s.certPath, *s.keyPath,
		rest.WithLogger(logger),
		rest.WithPort(*s.webConsoleServicePort),
		rest.WithService(NewWebConsoleService(*s.webConsoleEnvJSPath)),
		rest.WithMetrics(registerer),
	)
	go webConsoleServer.Run()
	dashboardServer := rest.NewServer(
		*s.certPath, *s.keyPath,
		rest.WithLogger(logger),
		rest.WithPort(*s.dashboardServicePort),
		rest.WithService(NewDashboardService()),
		rest.WithMetrics(registerer),
	)
	go dashboardServer.Run()
	// To detach this pod from Kubernetes Service before the app servers stop, we stop the health check service first.
	// Then, after 10 seconds of sleep, the app servers can be shut down, as no new requests are expected to be sent.
	// In this case, the Readiness prove must fail within 10 seconds and the pod must be detached.
	defer func() {
		go healthcheckServer.Stop(serverShutDownTimeout)
		time.Sleep(serverShutDownTimeout)
		// Stop gRPC servers
		go authServer.Stop(serverShutDownTimeout)
		go accountServer.Stop(serverShutDownTimeout)
		go auditLogServer.Stop(serverShutDownTimeout)
		go autoOpsServer.Stop(serverShutDownTimeout)
		go environmentServer.Stop(serverShutDownTimeout)
		go experimentServer.Stop(serverShutDownTimeout)
		go eventCounterServer.Stop(serverShutDownTimeout)
		go featureServer.Stop(serverShutDownTimeout)
		go notificationServer.Stop(serverShutDownTimeout)
		go pushServer.Stop(serverShutDownTimeout)
		go tagServer.Stop(serverShutDownTimeout)
		go webConsoleServer.Stop(serverShutDownTimeout)
		go codeReferenceServer.Stop(serverShutDownTimeout)
		// Stop REST gateways
		go accountGateway.Stop(context.Background())
		go auditLogGateway.Stop(context.Background())
		go autoOpsGateway.Stop(context.Background())
		go environmentGateway.Stop(context.Background())
		go eventCounterGateway.Stop(context.Background())
		go experimentGateway.Stop(context.Background())
		go featureGateway.Stop(context.Background())
		go notificationGateway.Stop(context.Background())
		go pushGateway.Stop(context.Background())
		go tagGateway.Stop(context.Background())
		go codeRefGateway.Stop(context.Background())
		// Close clients
		go mysqlClient.Close()
		go persistentRedisClient.Close()
		go nonPersistentRedisClient.Close()
		go bigQueryQuerier.Close()
		go domainTopicPublisher.Stop()
		go segmentUsersPublisher.Stop()
		go accountClient.Close()
		go authClient.Close()
		go batchClient.Close()
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

func (s *server) readOAuthConfig(
	ctx context.Context,
	logger *zap.Logger,
) (*auth.OAuthConfig, error) {
	bytes, err := os.ReadFile(*s.oauthConfigPath)
	if err != nil {
		logger.Error("auth: failed to read auth config file",
			zap.Error(err),
		)
		return nil, err
	}
	config := auth.OAuthConfig{}
	if err = json.Unmarshal(bytes, &config); err != nil {
		logger.Error("auth: failed to unmarshal auth config",
			zap.Error(err),
		)
		return nil, err
	}
	return &config, nil
}

func (s *server) createAuthService(
	mysqlClient mysql.Client,
	accountClient accountclient.Client,
	verifier token.Verifier,
	config *auth.OAuthConfig,
	logger *zap.Logger,
) (rpc.Service, error) {
	signer, err := token.NewSigner(*s.oauthPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	serviceOptions := []authapi.Option{
		authapi.WithLogger(logger),
		authapi.WithRefreshTokenTTL(*s.refreshTokenTTL),
	}
	if *s.emailFilter != "" {
		filter, err := regexp.Compile(*s.emailFilter)
		if err != nil {
			return nil, err
		}
		serviceOptions = append(serviceOptions, authapi.WithEmailFilter(filter))
	}
	return authapi.NewAuthService(
		config.Issuer,
		config.Audience,
		signer,
		verifier,
		mysqlClient,
		accountClient,
		config,
		serviceOptions...,
	), nil
}

func (s *server) createFeatureService(
	ctx context.Context,
	accountClient accountclient.Client,
	experimentClient experimentclient.Client,
	autoOpsClient autoopsclient.Client,
	batchClient btclient.Client,
	environmentClient environmentclient.Client,
	nonPersistentRedisV3Cache cache.MultiGetDeleteCountCache,
	segmentUsersPublisher publisher.Publisher,
	domainTopicPublisher publisher.Publisher,
	mysqlClient mysql.Client,
	logger *zap.Logger,
) (rpc.Service, error) {
	featureService := featureapi.NewFeatureService(
		mysqlClient,
		accountClient,
		experimentClient,
		autoOpsClient,
		batchClient,
		environmentClient,
		nonPersistentRedisV3Cache,
		segmentUsersPublisher,
		domainTopicPublisher,
		fmt.Sprintf("%s/%s", *s.webhookBaseURL, featureFlagTriggerWebhookPath),
		featureapi.WithLogger(logger),
	)

	return featureService, nil
}
