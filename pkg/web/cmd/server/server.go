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
	"gopkg.in/yaml.v2"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

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
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/factory"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	pushapi "github.com/bucketeer-io/bucketeer/pkg/push/api"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/pkg/rest"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/gateway"
	gatewayapi "github.com/bucketeer-io/bucketeer/pkg/rpc/gateway"
	bqquerier "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/querier"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	tagapi "github.com/bucketeer-io/bucketeer/pkg/tag/api"
	teamapi "github.com/bucketeer-io/bucketeer/pkg/team/api"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	auditlogproto "github.com/bucketeer-io/bucketeer/proto/auditlog"
	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	coderefproto "github.com/bucketeer-io/bucketeer/proto/coderef"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventcounterproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
	pushproto "github.com/bucketeer-io/bucketeer/proto/push"
	tagproto "github.com/bucketeer-io/bucketeer/proto/tag"
	teamproto "github.com/bucketeer-io/bucketeer/proto/team"
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
	port                            *int
	project                         *string
	isDemoSiteEnabled               *bool
	timezone                        *string
	certPath                        *string
	keyPath                         *string
	serviceTokenPath                *string
	mysqlUser                       *string
	mysqlPass                       *string
	mysqlHost                       *string
	mysqlPort                       *int
	mysqlDBName                     *string
	persistentRedisServerName       *string
	persistentRedisAddr             *string
	persistentRedisPoolMaxIdle      *int
	persistentRedisPoolMaxActive    *int
	nonPersistentRedisServerName    *string
	nonPersistentRedisAddr          *string
	nonPersistentRedisPoolMaxIdle   *int
	nonPersistentRedisPoolMaxActive *int
	bigQueryDataSet                 *string
	bigQueryDataLocation            *string
	domainTopic                     *string
	bulkSegmentUsersReceivedTopic   *string
	accountServicePort              *int
	authServicePort                 *int
	auditLogServicePort             *int
	autoOpsServicePort              *int
	environmentServicePort          *int
	eventCounterServicePort         *int
	experimentServicePort           *int
	featureServicePort              *int
	notificationServicePort         *int
	pushServicePort                 *int
	webConsoleServicePort           *int
	dashboardServicePort            *int
	tagServicePort                  *int
	codeReferenceServicePort        *int
	teamServicePort                 *int
	webGrpcGatewayPort              *int
	accountService                  *string
	authService                     *string
	batchService                    *string
	environmentService              *string
	experimentService               *string
	featureService                  *string
	autoOpsService                  *string
	codeReferenceService            *string
	refreshTokenTTL                 *time.Duration
	emailFilter                     *string
	oauthConfigPath                 *string
	oauthPublicKeyPath              *string
	oauthPrivateKeyPath             *string
	webhookBaseURL                  *string
	webhookKMSResourceName          *string
	cloudService                    *string
	webConsoleEnvJSPath             *string
	pubSubType                      *string
	pubSubRedisServerName           *string
	pubSubRedisAddr                 *string
	pubSubRedisPoolSize             *int
	pubSubRedisMinIdle              *int
	pubSubRedisPartitionCount       *int
	dataWarehouseType               *string
	dataWarehouseConfigPath         *string
}

type DataWarehouseConfig struct {
	Type      string                      `yaml:"type"`
	BatchSize int                         `yaml:"batchSize"`
	Timezone  string                      `yaml:"timezone"`
	BigQuery  DataWarehouseBigQueryConfig `yaml:"bigquery"`
	MySQL     DataWarehouseMySQLConfig    `yaml:"mysql"`
}

type DataWarehouseBigQueryConfig struct {
	Project  string `yaml:"project"`
	Dataset  string `yaml:"dataset"`
	Location string `yaml:"location"`
}

type DataWarehouseMySQLConfig struct {
	UseMainConnection bool   `yaml:"useMainConnection"`
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	User              string `yaml:"user"`
	Password          string `yaml:"password"`
	Database          string `yaml:"database"`
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the server")
	server := &server{
		CmdClause: cmd,
		port:      cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:   cmd.Flag("project", "Google Cloud project name.").Required().String(),
		isDemoSiteEnabled: cmd.Flag(
			"demo-site-enabled",
			"Is demo site enabled").Default("false").Bool(),
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
		teamServicePort: cmd.Flag(
			"team-service-port",
			"Port to bind to team service.",
		).Default("9107").Int(),
		webGrpcGatewayPort: cmd.Flag(
			"web-grpc-gateway-port",
			"Port to bind to web gRPC gateway.",
		).Default("9089").Int(),
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
		dataWarehouseType: cmd.Flag(
			"data-warehouse-type",
			"Data warehouse type (bigquery, mysql).",
		).Default("bigquery").String(),
		dataWarehouseConfigPath: cmd.Flag(
			"data-warehouse-config-path",
			"Path to data warehouse configuration file.",
		).String(),
		timezone:         cmd.Flag("timezone", "Time zone").Required().String(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		oauthPublicKeyPath: cmd.Flag(
			"oauth-public-key",
			"Path to public key used to verify oauth token.",
		).Required().String(),
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
		webhookBaseURL: cmd.Flag("webhook-base-url", "the base url for incoming webhooks.").Required().String(),
		webhookKMSResourceName: cmd.Flag(
			"webhook-kms-resource-name",
			"Cloud KMS resource name to encrypt and decrypt webhook credentials.",
		).Required().String(),
		cloudService:        cmd.Flag("cloud-service", "Cloud Service info").Default(gcp).String(),
		webConsoleEnvJSPath: cmd.Flag("web-console-env-js-path", "console env js path").Required().String(),
		// PubSub configuration
		pubSubType: cmd.Flag("pubsub-type",
			"Type of PubSub to use (google or redis-stream).",
		).Default("google").String(),
		pubSubRedisServerName: cmd.Flag("pubsub-redis-server-name",
			"Name of the Redis server for PubSub.",
		).Default("non-persistent-redis").String(),
		pubSubRedisAddr: cmd.Flag("pubsub-redis-addr",
			"Address of the Redis server for PubSub.",
		).Default("localhost:6379").String(),
		pubSubRedisPoolSize: cmd.Flag("pubsub-redis-pool-size",
			"Maximum number of connections for Redis PubSub.",
		).Default("10").Int(),
		pubSubRedisMinIdle: cmd.Flag("pubsub-redis-min-idle",
			"Minimum number of idle connections for Redis PubSub.",
		).Default("5").Int(),
		pubSubRedisPartitionCount: cmd.Flag("pubsub-redis-partition-count",
			"Number of partitions for Redis Streams PubSub.",
		).Default("16").Int(),
	}
	r.RegisterCommand(server)
	return server
}

func (s *server) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()

	// dataWarehouse config
	dataWarehouseConfig, err := s.readDataWarehouseConfig(ctx, logger)
	if err != nil {
		logger.Error("Failed to read dataWarehouse config", zap.Error(err))
		return err
	}

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
	var bigQueryQuerier bqquerier.Client
	if dataWarehouseConfig.Type == "bigquery" {
		bigQueryQuerier, err = s.createBigQueryQuerier(
			ctx, *s.project, dataWarehouseConfig.BigQuery.Location, registerer, logger,
		)
		if err != nil {
			logger.Error("Failed to create BigQuery client",
				zap.Error(err),
				zap.String("project", *s.project),
				zap.String("location", dataWarehouseConfig.BigQuery.Location),
				zap.String("data-set", dataWarehouseConfig.BigQuery.Dataset),
			)
			return err
		}
	}
	// bigQueryDataSet
	bigQueryDataSet := dataWarehouseConfig.BigQuery.Dataset
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
	autoOpsServer := rpc.NewServer(autoOpsService, *s.certPath, *s.keyPath,
		"auto-ops-server",
		rpc.WithPort(*s.autoOpsServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	go autoOpsServer.Run()

	// environmentService
	environmentService, err := s.createEnvironmentService(
		mysqlClient,
		accountClient,
		domainTopicPublisher,
		oAuthConfig,
		verifier,
		oAuthConfig,
		logger,
	)
	if err != nil {
		return err
	}
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
		eventcounterapi.WithDataWarehouseConfig(s.convertToAPIDataWarehouseConfig(dataWarehouseConfig)),
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

	// teamService
	teamService := teamapi.NewTeamService(
		mysqlClient,
		accountClient,
		domainTopicPublisher,
		teamapi.WithLogger(logger),
	)
	teamServer := rpc.NewServer(teamService, *s.certPath, *s.keyPath,
		"team-server",
		rpc.WithPort(*s.teamServicePort),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
	)
	go teamServer.Run()

	// Start the web console and dashboard servers
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
		rest.WithService(NewDashboardService(*s.webConsoleEnvJSPath)),
		rest.WithMetrics(registerer),
	)
	go dashboardServer.Run()

	// Set up REST gateway
	restAddr := fmt.Sprintf(":%d", *s.webGrpcGatewayPort)

	webGrpcGateway, err := gateway.NewGateway(
		restAddr,
		gateway.WithLogger(logger.Named("web-grpc-gateway")),
		gateway.WithMetrics(registerer),
		gateway.WithCertPath(*s.certPath),
		gateway.WithKeyPath(*s.keyPath),
	)
	if err != nil {
		return fmt.Errorf("failed to create web gRPC gateway: %v", err)
	}

	if err := webGrpcGateway.Start(ctx, s.createGatewayHandlers()...); err != nil {
		return fmt.Errorf("failed to start web gRPC gateway: %v", err)
	}

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
		go teamServer.Stop(serverShutDownTimeout)
		go webGrpcGateway.Stop(serverShutDownTimeout)
		// Close clients
		go mysqlClient.Close()
		go persistentRedisClient.Close()
		go nonPersistentRedisClient.Close()
		if dataWarehouseConfig.Type == "bigquery" {
			go bigQueryQuerier.Close()
		}
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

	// Create PubSub client using the factory
	pubSubType := factory.PubSubType(*s.pubSubType)
	factoryOpts := []factory.Option{
		factory.WithPubSubType(pubSubType),
		factory.WithMetrics(registerer),
		factory.WithLogger(logger),
	}

	// Add provider-specific options
	if pubSubType == factory.Google {
		factoryOpts = append(factoryOpts, factory.WithProjectID(*s.project))
	} else if pubSubType == factory.RedisStream {
		redisClient, err := redisv3.NewClient(
			*s.pubSubRedisAddr,
			redisv3.WithPoolSize(*s.pubSubRedisPoolSize),
			redisv3.WithMinIdleConns(*s.pubSubRedisMinIdle),
			redisv3.WithServerName(*s.pubSubRedisServerName),
			redisv3.WithMetrics(registerer),
			redisv3.WithLogger(logger),
		)
		if err != nil {
			return nil, err
		}
		factoryOpts = append(factoryOpts, factory.WithRedisClient(redisClient))
		factoryOpts = append(factoryOpts, factory.WithPartitionCount(*s.pubSubRedisPartitionCount))
	}

	// Create the PubSub client using the factory
	pubsubClient, err := factory.NewClient(ctx, factoryOpts...)
	if err != nil {
		return nil, err
	}

	// Create publisher for the topic
	return pubsubClient.CreatePublisher(topic)
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
		authapi.WithDemoSiteEnabled(*s.isDemoSiteEnabled),
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

func (s *server) createEnvironmentService(
	mysqlClient mysql.Client,
	accountClient accountclient.Client,
	domainTopicPublisher publisher.Publisher,
	oAuthConfig *auth.OAuthConfig,
	verifier token.Verifier,
	config *auth.OAuthConfig,
	logger *zap.Logger,
) (rpc.Service, error) {
	signer, err := token.NewSigner(*s.oauthPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	serviceOptions := []environmentapi.Option{
		environmentapi.WithLogger(logger),
		environmentapi.WithRefreshTokenTTL(*s.refreshTokenTTL),
		environmentapi.WithDemoSiteEnabled(*s.isDemoSiteEnabled),
	}
	if *s.emailFilter != "" {
		filter, err := regexp.Compile(*s.emailFilter)
		if err != nil {
			return nil, err
		}
		serviceOptions = append(serviceOptions, environmentapi.WithEmailFilter(filter))
	}

	return environmentapi.NewEnvironmentService(
		accountClient,
		mysqlClient,
		domainTopicPublisher,
		oAuthConfig,
		config.Issuer,
		config.Audience,
		signer,
		verifier,
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

func (s *server) createGatewayHandlers() []gatewayapi.HandlerRegistrar {
	return []gatewayapi.HandlerRegistrar{
		func(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
			accountGrpcAddr := fmt.Sprintf("localhost:%d", *s.accountServicePort)
			return accountproto.RegisterAccountServiceHandlerFromEndpoint(ctx, mux, accountGrpcAddr, opts)
		},
		func(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
			authGrpcAddr := fmt.Sprintf("localhost:%d", *s.authServicePort)
			return authproto.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, authGrpcAddr, opts)
		},
		func(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
			auditLogGrpcAddr := fmt.Sprintf("localhost:%d", *s.auditLogServicePort)
			return auditlogproto.RegisterAuditLogServiceHandlerFromEndpoint(ctx, mux, auditLogGrpcAddr, opts)
		},
		func(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
			autoOpsGrpcAddr := fmt.Sprintf("localhost:%d", *s.autoOpsServicePort)
			return autoopsproto.RegisterAutoOpsServiceHandlerFromEndpoint(ctx, mux, autoOpsGrpcAddr, opts)
		},
		func(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
			environmentGrpcAddr := fmt.Sprintf("localhost:%d", *s.environmentServicePort)
			return environmentproto.RegisterEnvironmentServiceHandlerFromEndpoint(ctx, mux, environmentGrpcAddr, opts)
		},
		func(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
			eventCounterGrpcAddr := fmt.Sprintf("localhost:%d", *s.eventCounterServicePort)
			return eventcounterproto.RegisterEventCounterServiceHandlerFromEndpoint(ctx, mux, eventCounterGrpcAddr, opts)
		},
		func(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
			experimentGrpcAddr := fmt.Sprintf("localhost:%d", *s.experimentServicePort)
			return experimentproto.RegisterExperimentServiceHandlerFromEndpoint(ctx, mux, experimentGrpcAddr, opts)
		},
		func(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
			featureGrpcAddr := fmt.Sprintf("localhost:%d", *s.featureServicePort)
			return featureproto.RegisterFeatureServiceHandlerFromEndpoint(ctx, mux, featureGrpcAddr, opts)
		},
		func(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
			notificationGrpcAddr := fmt.Sprintf("localhost:%d", *s.notificationServicePort)
			return notificationproto.RegisterNotificationServiceHandlerFromEndpoint(ctx, mux, notificationGrpcAddr, opts)
		},
		func(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
			pushGrpcAddr := fmt.Sprintf("localhost:%d", *s.pushServicePort)
			return pushproto.RegisterPushServiceHandlerFromEndpoint(ctx, mux, pushGrpcAddr, opts)
		},
		func(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
			tagGrpcAddr := fmt.Sprintf("localhost:%d", *s.tagServicePort)
			return tagproto.RegisterTagServiceHandlerFromEndpoint(ctx, mux, tagGrpcAddr, opts)
		},
		func(ctx context.Context, mux *runtime.ServeMux, options []grpc.DialOption) error {
			teamGrpcAddr := fmt.Sprintf("localhost:%d", *s.teamServicePort)
			return teamproto.RegisterTeamServiceHandlerFromEndpoint(ctx, mux, teamGrpcAddr, options)
		},
		func(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
			codeRefGrpcAddr := fmt.Sprintf("localhost:%d", *s.codeReferenceServicePort)
			return coderefproto.RegisterCodeReferenceServiceHandlerFromEndpoint(ctx, mux, codeRefGrpcAddr, opts)
		},
	}
}

func (s *server) readDataWarehouseConfig(
	ctx context.Context,
	logger *zap.Logger,
) (*DataWarehouseConfig, error) {
	// If config path is provided, read from file
	if *s.dataWarehouseConfigPath != "" {
		bytes, err := os.ReadFile(*s.dataWarehouseConfigPath)
		if err != nil {
			logger.Error("Failed to read dataWarehouse config file",
				zap.Error(err),
				zap.String("path", *s.dataWarehouseConfigPath),
			)
			return nil, err
		}
		config := DataWarehouseConfig{}
		if err = yaml.Unmarshal(bytes, &config); err != nil {
			logger.Error("Failed to unmarshal dataWarehouse config",
				zap.Error(err),
			)
			return nil, err
		}
		return &config, nil
	}

	// Fallback to environment variables / command line flags
	config := &DataWarehouseConfig{
		Type:      *s.dataWarehouseType,
		BatchSize: 1000,
		Timezone:  *s.timezone,
	}

	// Set default configurations based on type
	switch config.Type {
	case "mysql":
		config.MySQL = DataWarehouseMySQLConfig{
			UseMainConnection: true, // Default to using main connection
		}
	case "bigquery":
		config.BigQuery = DataWarehouseBigQueryConfig{
			Dataset:  *s.bigQueryDataSet,
			Location: *s.bigQueryDataLocation,
		}
	}

	return config, nil
}

func (s *server) convertToAPIDataWarehouseConfig(config *DataWarehouseConfig) *eventcounterapi.DataWarehouseConfig {
	return &eventcounterapi.DataWarehouseConfig{
		Type:      config.Type,
		BatchSize: config.BatchSize,
		Timezone:  config.Timezone,
		MySQL: eventcounterapi.DataWarehouseMySQLConfig{
			UseMainConnection: config.MySQL.UseMainConnection,
			Host:              config.MySQL.Host,
			Port:              config.MySQL.Port,
			User:              config.MySQL.User,
			Password:          config.MySQL.Password,
			Database:          config.MySQL.Database,
		},
		BigQuery: eventcounterapi.DataWarehouseBigQueryConfig{
			Project:  config.BigQuery.Project,
			Dataset:  config.BigQuery.Dataset,
			Location: config.BigQuery.Location,
		},
	}
}
