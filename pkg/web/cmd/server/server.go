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
	bqquerier "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/querier"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	tagapi "github.com/bucketeer-io/bucketeer/pkg/tag/api"
	"github.com/bucketeer-io/bucketeer/pkg/token"
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
		go notificationServer.Stop(serverShutDownTimeout)
		go pushServer.Stop(serverShutDownTimeout)
		go tagServer.Stop(serverShutDownTimeout)
		go webConsoleServer.Stop(serverShutDownTimeout)
		go codeReferenceServer.Stop(serverShutDownTimeout)
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
