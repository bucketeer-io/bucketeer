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
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"

	acclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	autoopsclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	"github.com/bucketeer-io/bucketeer/pkg/batch/api"
	btclient "github.com/bucketeer-io/bucketeer/pkg/batch/client"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	cacher "github.com/bucketeer-io/bucketeer/pkg/batch/jobs/cacher"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/calculator"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/deleter"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/experiment"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/mau"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/notification"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/opsevent"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/rediscounter"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	"github.com/bucketeer-io/bucketeer/pkg/experimentcalculator/stan"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	notificationclient "github.com/bucketeer-io/bucketeer/pkg/notification/client"
	notificationsender "github.com/bucketeer-io/bucketeer/pkg/notification/sender"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender/notifier"
	"github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor"
	opsexecutor "github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor"
	pushclient "github.com/bucketeer-io/bucketeer/pkg/push/client"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/gateway"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	batchproto "github.com/bucketeer-io/bucketeer/proto/batch"
)

const (
	command               = "server"
	clientDialTimeout     = 30 * time.Second
	serverShutDownTimeout = 10 * time.Second
)

type server struct {
	*kingpin.CmdClause
	// Common
	port               *int
	project            *string
	certPath           *string
	keyPath            *string
	serviceTokenPath   *string
	timezone           *string
	refreshInterval    *time.Duration
	experimentLockTTL  *time.Duration
	webURL             *string
	oauthPublicKeyPath *string
	oauthAudience      *string
	oauthIssuer        *string
	stanHost           *string
	stanPort           *string
	stanModelID        *string
	// MySQL
	mysqlUser        *string
	mysqlPass        *string
	mysqlHost        *string
	mysqlPort        *int
	mysqlDBName      *string
	mysqlDBOpenConns *int
	// gRPC service
	accountService              *string
	environmentService          *string
	experimentService           *string
	autoOpsService              *string
	eventCounterService         *string
	pushService                 *string
	featureService              *string
	notificationService         *string
	experimentCalculatorService *string
	batchService                *string
	// Persistent Redis
	persistentRedisServerName    *string
	persistentRedisAddr          *string
	persistentRedisPoolMaxIdle   *int
	persistentRedisPoolMaxActive *int
	// Non Persistent Redis
	nonPersistentRedisServerName     *string
	nonPersistentRedisAddr           *string
	nonPersistentChildRedisAddresses *[]string
	nonPersistentRedisPoolMaxIdle    *int
	nonPersistentRedisPoolMaxActive  *int
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start batch server")
	server := &server{
		CmdClause:        cmd,
		port:             cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:          cmd.Flag("project", "Google Cloud project name.").String(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		timezone:         cmd.Flag("timezone", "Time zone").Required().String(),
		refreshInterval: cmd.Flag(
			"refresh-interval",
			"Interval between refreshing target objects.",
		).Default("1m").Duration(),
		webURL: cmd.Flag("web-url", "Web console URL.").Required().String(),
		oauthPublicKeyPath: cmd.Flag(
			"oauth-public-key",
			"Path to public key used to verify oauth token.",
		).Required().String(),
		oauthAudience: cmd.Flag(
			"oauth-audience",
			"The oauth audience registered in the token",
		).Required().String(),
		oauthIssuer:      cmd.Flag("oauth-issuer", "The issuer url").Required().String(),
		stanHost:         cmd.Flag("stan-host", "httpstan host.").Default("localhost").String(),
		stanPort:         cmd.Flag("stan-port", "httpstan port.").Default("8080").String(),
		stanModelID:      cmd.Flag("stan-model-id", "httpstan modelId.").Required().String(),
		mysqlUser:        cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:        cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:        cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:        cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName:      cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		mysqlDBOpenConns: cmd.Flag("mysql-db-open-conns", "MySQL open connections.").Required().Int(),
		accountService: cmd.Flag(
			"account-service",
			"bucketeer-account-service address.",
		).Default("account:9090").String(),
		environmentService: cmd.Flag(
			"environment-service",
			"bucketeer-environment-service address.",
		).Default("environment:9090").String(),
		experimentService: cmd.Flag(
			"experiment-service",
			"bucketeer-experiment-service address.",
		).Default("experiment:9090").String(),
		autoOpsService: cmd.Flag(
			"auto-ops-service",
			"bucketeer-auto-ops-service address.",
		).Default("auto-ops:9090").String(),
		eventCounterService: cmd.Flag(
			"event-counter-service",
			"bucketeer-event-counter-service address.",
		).Default("event-counter-server:9090").String(),
		pushService: cmd.Flag(
			"push-service",
			"bucketeer-push-service address.",
		).Default("push:9090").String(),
		featureService: cmd.Flag(
			"feature-service",
			"bucketeer-feature-service address.",
		).Default("feature:9090").String(),
		notificationService: cmd.Flag(
			"notification-service",
			"bucketeer-notification-service address.",
		).Default("notification:9090").String(),
		experimentCalculatorService: cmd.Flag(
			"experiment-calculator-service",
			"bucketeer-experiment-calculator-service address.",
		).Default("experiment-calculator:9090").String(),
		batchService: cmd.Flag(
			"batch-service",
			"bucketeer-batch-service address.",
		).Default("localhost:9001").String(),
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
		nonPersistentChildRedisAddresses: cmd.Flag(
			"non-persistent-child-redis-addresses",
			"A list of non-persistent child Redis addresses.",
		).Strings(),
		experimentLockTTL: cmd.Flag("experiment-lock-ttl",
			"The ttl for experiment calculator lock").
			Default("10m").Duration(),
	}
	r.RegisterCommand(server)
	return server
}

func (s *server) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	*s.serviceTokenPath = s.insertTelepresenceMountRoot(*s.serviceTokenPath)
	*s.keyPath = s.insertTelepresenceMountRoot(*s.keyPath)
	*s.certPath = s.insertTelepresenceMountRoot(*s.certPath)

	registerer := metrics.DefaultRegisterer()

	verifier, err := token.NewVerifier(*s.oauthPublicKeyPath, *s.oauthIssuer, *s.oauthAudience)
	if err != nil {
		return err
	}

	mysqlClient, err := s.createMySQLClient(ctx, registerer, logger)
	if err != nil {
		return err
	}

	creds, err := client.NewPerRPCCredentials(*s.serviceTokenPath)
	if err != nil {
		return err
	}

	accountClient, err := acclient.NewClient(*s.accountService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	notificationClient, err := notificationclient.NewClient(*s.notificationService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}

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

	eventCounterClient, err := ecclient.NewClient(*s.eventCounterService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	pushClient, err := pushclient.NewClient(*s.pushService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer pushClient.Close()

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

	autoOpsExecutor := opsexecutor.NewAutoOpsExecutor(
		autoOpsClient,
		opsexecutor.WithLogger(logger),
	)

	progressiveRolloutExecutor := opsexecutor.NewProgressiveRolloutExecutor(
		autoOpsClient,
		executor.WithLogger(logger),
	)

	slackNotifier := notifier.NewSlackNotifier(*s.webURL)

	notificationSender := notificationsender.NewSender(
		notificationClient,
		[]notifier.Notifier{slackNotifier},
		notificationsender.WithMetrics(registerer),
		notificationsender.WithLogger(logger),
	)

	location, err := locale.GetLocation(*s.timezone)
	if err != nil {
		return err
	}

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
	defer persistentRedisClient.Close()

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
	defer nonPersistentRedisClient.Close()

	// This slice contains all Redis instance caches
	nonPersistentRedisCaches := make(
		[]cache.MultiGetCache,
		0,
		len(*s.nonPersistentChildRedisAddresses),
	)
	// Append the main Redis instance
	nonPersistentRedisCaches = append(
		nonPersistentRedisCaches,
		cachev3.NewRedisCache(nonPersistentRedisClient),
	)
	// Initialize all the child Redis clients
	childRedisClients := make([]redisv3.Client, 0, len(*s.nonPersistentChildRedisAddresses))
	for _, address := range *s.nonPersistentChildRedisAddresses {
		// We use the same options used for the main non-persistent Redis, besides the name
		client, err := redisv3.NewClient(
			address,
			redisv3.WithPoolSize(*s.nonPersistentRedisPoolMaxActive),
			redisv3.WithMinIdleConns(*s.nonPersistentRedisPoolMaxIdle),
			redisv3.WithServerName(s.getRedisHostname(address)),
			redisv3.WithMetrics(registerer),
			redisv3.WithLogger(logger),
		)
		if err != nil {
			return err
		}
		childRedisClients = append(childRedisClients, client)
		nonPersistentRedisCaches = append(nonPersistentRedisCaches, cachev3.NewRedisCache(client))
	}
	// TODO: To be removed after checked it works
	logger.Debug("Redis main address", zap.String("address", *s.nonPersistentRedisAddr))
	logger.Debug("Redis child addresses", zap.Strings("addresses", *s.nonPersistentChildRedisAddresses))
	logger.Debug("Redis non persistent cache size", zap.Int("size", len(nonPersistentRedisCaches)))

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
	defer batchClient.Close()

	service := api.NewBatchService(
		experiment.NewExperimentStatusUpdater(
			environmentClient,
			experimentClient,
			jobs.WithLogger(logger),
		),
		notification.NewExperimentRunningWatcher(
			environmentClient,
			experimentClient,
			notificationSender,
			jobs.WithTimeout(1*time.Minute),
			jobs.WithLogger(logger),
		),
		notification.NewFeatureStaleWatcher(
			environmentClient,
			featureClient,
			notificationSender,
			jobs.WithTimeout(1*time.Minute),
			jobs.WithLogger(logger),
		),
		notification.NewMAUCountWatcher(
			environmentClient,
			eventCounterClient,
			notificationSender,
			location,
			jobs.WithTimeout(60*time.Minute),
			jobs.WithLogger(logger),
		),
		opsevent.NewDatetimeWatcher(
			environmentClient,
			autoOpsClient,
			autoOpsExecutor,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		opsevent.NewEventCountWatcher(
			mysqlClient,
			environmentClient,
			autoOpsClient,
			eventCounterClient,
			featureClient,
			autoOpsExecutor,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		opsevent.NewProgressiveRolloutWacher(
			environmentClient,
			autoOpsClient,
			progressiveRolloutExecutor,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		rediscounter.NewRedisCounterDeleter(
			cachev3.NewRedisCache(persistentRedisClient),
			environmentClient,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		calculator.NewExperimentCalculate(
			stan.NewStan(*s.stanHost, *s.stanPort),
			*s.stanModelID,
			environmentClient,
			experimentClient,
			eventCounterClient,
			mysqlClient,
			calculator.NewExperimentLock(nonPersistentRedisClient, *s.experimentLockTTL),
			location,
			jobs.WithTimeout(30*time.Minute),
			jobs.WithLogger(logger),
			jobs.WithMetrics(registerer),
		),
		mau.NewMAUSummarizer(
			mysqlClient,
			eventCounterClient,
			location,
			jobs.WithTimeout(30*time.Minute),
			jobs.WithLogger(logger),
		),
		mau.NewMAUPartitionDeleter(
			mysqlClient,
			location,
			jobs.WithTimeout(60*time.Minute),
			jobs.WithLogger(logger),
		),
		mau.NewMAUPartitionCreator(
			mysqlClient,
			location,
			jobs.WithTimeout(60*time.Minute),
			jobs.WithLogger(logger),
		),
		cacher.NewFeatureFlagCacher(
			mysqlClient,
			nonPersistentRedisCaches,
			jobs.WithLogger(logger),
		),
		cacher.NewSegmentUserCacher(
			environmentClient,
			featureClient,
			nonPersistentRedisCaches,
			jobs.WithLogger(logger),
		),
		cacher.NewAPIKeyCacher(
			mysqlClient,
			nonPersistentRedisCaches,
			jobs.WithLogger(logger),
		),
		cacher.NewExperimentCacher(
			environmentClient,
			experimentClient,
			nonPersistentRedisCaches,
			jobs.WithLogger(logger),
		),
		cacher.NewAutoOpsRulesCacher(
			environmentClient,
			autoOpsClient,
			// Because the event-perister-ops uses persistent redis
			// We must use the same instance for caching.
			cachev3.NewRedisCache(persistentRedisClient),
			jobs.WithLogger(logger),
		),
		deleter.NewTagDeleter(
			mysqlClient,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		logger,
	)

	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
		health.WithCheck("persistent-redis", persistentRedisClient.Check),
		health.WithCheck("non-persistent-redis", nonPersistentRedisClient.Check),
	)
	go healthChecker.Run(ctx)

	server := rpc.NewServer(service, *s.certPath, *s.keyPath,
		"batch-server",
		rpc.WithVerifier(verifier),
		rpc.WithPort(*s.port),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
		rpc.WithService(healthChecker),
		rpc.WithHandler("/health", healthChecker),
	)
	go server.Run()

	// Setup REST gateway for batch service
	restPort := *s.port + 1000
	restAddr := fmt.Sprintf(":%d", restPort)
	grpcAddr := fmt.Sprintf("localhost:%d", *s.port)

	// Create a HandlerRegistrar adapter function that matches gateway.HandlerRegistrar signature
	batchHandler := func(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
		return batchproto.RegisterBatchServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	}

	batchGateway, err := gateway.NewGateway(
		restAddr,
		gateway.WithLogger(logger.Named("batch-gateway")),
		gateway.WithMetrics(registerer),
		gateway.WithCertPath(*s.certPath),
		gateway.WithKeyPath(*s.keyPath),
	)
	if err != nil {
		return fmt.Errorf("failed to create batch gateway: %v", err)
	}

	go func() {
		if err := batchGateway.Start(
			ctx,
			batchHandler,
		); err != nil {
			logger.Error("failed to start batch gateway", zap.Error(err))
		}
	}()

	defer func() {
		server.Stop(serverShutDownTimeout)
		batchGateway.Stop(context.Background())
		accountClient.Close()
		notificationClient.Close()
		experimentClient.Close()
		environmentClient.Close()
		eventCounterClient.Close()
		featureClient.Close()
		autoOpsClient.Close()
		mysqlClient.Close()
		for _, client := range childRedisClients {
			client.Close()
		}
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
		mysql.WithMaxOpenConns(*s.mysqlDBOpenConns),
	)
}

func (s *server) insertTelepresenceMountRoot(path string) string {
	volumeRoot := os.Getenv("TELEPRESENCE_ROOT")
	if volumeRoot == "" {
		return path
	}
	return volumeRoot + path
}

func (s *server) getRedisHostname(redisAddress string) string {
	address := strings.Split(redisAddress, ":")
	if len(address) == 0 {
		return redisAddress
	}
	return address[0]
}
