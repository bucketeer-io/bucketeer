// Copyright 2026 The Bucketeer Authors.
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
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	accstorage "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	accountmysql "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mysql"
	accountpostgres "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/postgres"
	v2als "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2"
	auditlogmysql "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2/mysql"
	auditlogpostgres "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2/postgres"
	autoopsclient "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client"
	btclient "github.com/bucketeer-io/bucketeer/v2/pkg/batch/client"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cli"
	"github.com/bucketeer-io/bucketeer/v2/pkg/email"
	environmentclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	ecdwh "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/storage/v2/dwh_database"
	ecbigquery "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/storage/v2/dwh_database/bigquery"
	ecmysql "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/storage/v2/dwh_database/mysql"
	ecpostgres "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/storage/v2/dwh_database/postgres"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	featuremysql "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mysql"
	featurepostgres "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/postgres"
	"github.com/bucketeer-io/bucketeer/v2/pkg/health"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/factory"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	pushstorage "github.com/bucketeer-io/bucketeer/v2/pkg/push/storage/v2"
	redisv3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rest"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	bqquerier "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/bigquery/querier"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/processor"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/dwhstorage"
	dwhbigquery "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/dwhstorage/bigquery"
	dwhmysql "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/dwhstorage/mysql"
	dwhpostgres "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/dwhstorage/postgres"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/operationalstorage"
	opmysql "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/operationalstorage/mysql"
	oppostgres "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/operationalstorage/postgres"
	notificationclient "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/client"
	notificationsender "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/sender"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscription/sender/notifier"
)

const (
	command            = "server"
	healthCheckTimeout = 1 * time.Second
	clientDialTimeout  = 30 * time.Second
)

type server struct {
	*kingpin.CmdClause
	// Common
	port             *int
	project          *string
	certPath         *string
	keyPath          *string
	serviceTokenPath *string
	webURL           *string
	emailConfigPath  *string
	demoSiteEnabled  *bool
	// Operational database
	operationalDatabaseType *string
	// MySQL
	mysqlUser        *string
	mysqlPass        *string
	mysqlHost        *string
	mysqlPort        *int
	mysqlDBName      *string
	mysqlDBOpenConns *int
	// PostgreSQL
	postgresUser   *string
	postgresPass   *string
	postgresHost   *string
	postgresPort   *int
	postgresDBName *string
	// gRPC service
	environmentService          *string
	experimentService           *string
	autoOpsService              *string
	eventCounterService         *string
	pushService                 *string
	featureService              *string
	notificationService         *string
	experimentCalculatorService *string
	batchService                *string
	// PubSub config
	subscriberConfig         *string
	onDemandSubscriberConfig *string
	processorsConfig         *string
	onDemandProcessorsConfig *string
	// Persistent Redis
	persistentRedisServerName    *string
	persistentRedisAddr          *string
	persistentRedisPoolMaxIdle   *int
	persistentRedisPoolMaxActive *int
	persistentRedisMode          *string
	// Non Persistent Redis
	nonPersistentRedisServerName    *string
	nonPersistentRedisAddr          *string
	nonPersistentRedisPoolMaxIdle   *int
	nonPersistentRedisPoolMaxActive *int
	nonPersistentRedisMode          *string
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start subscriber server")
	server := &server{
		CmdClause:        cmd,
		port:             cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:          cmd.Flag("project", "Google Cloud project name.").String(),
		demoSiteEnabled:  cmd.Flag("demo-site-enabled", "Enable demo site.").Default("false").Bool(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		webURL:           cmd.Flag("web-url", "Web console URL.").Required().String(),
		emailConfigPath:  cmd.Flag("email-config-path", "Path to email config.").Required().String(),
		operationalDatabaseType: cmd.Flag("storage-type", "Operational database type (mysql, postgres).").
			Default("mysql").String(),
		mysqlUser:        cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:        cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:        cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:        cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName:      cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		mysqlDBOpenConns: cmd.Flag("mysql-db-open-conns", "MySQL open connections.").Required().Int(),
		postgresUser:     cmd.Flag("postgres-user", "PostgreSQL user.").String(),
		postgresPass:     cmd.Flag("postgres-pass", "PostgreSQL password.").String(),
		postgresHost:     cmd.Flag("postgres-host", "PostgreSQL host.").String(),
		postgresPort:     cmd.Flag("postgres-port", "PostgreSQL port.").Int(),
		postgresDBName:   cmd.Flag("postgres-db-name", "PostgreSQL database name.").String(),
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
		subscriberConfig: cmd.Flag(
			"subscriber-config",
			"Path to subscribers config.",
		).Required().String(),
		onDemandSubscriberConfig: cmd.Flag(
			"on-demand-subscriber-config",
			"Path to on-demand subscribers config.",
		).Required().String(),
		processorsConfig: cmd.Flag(
			"processors-config",
			"Path to processors config.",
		).Required().String(),
		onDemandProcessorsConfig: cmd.Flag(
			"on-demand-processors-config",
			"Path to on-demand processors config.",
		).Required().String(),
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
		persistentRedisMode: cmd.Flag("persistent-redis-mode",
			"Persistent Redis client mode: cluster, standalone, or auto.",
		).Default("auto").String(),
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
		nonPersistentRedisMode: cmd.Flag("non-persistent-redis-mode",
			"Non-persistent Redis client mode: cluster, standalone, or auto.",
		).Default("auto").String(),
	}
	r.RegisterCommand(server)
	return server
}

func (s *server) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	*s.serviceTokenPath = s.insertTelepresenceMountRoot(*s.serviceTokenPath)
	*s.keyPath = s.insertTelepresenceMountRoot(*s.keyPath)
	*s.certPath = s.insertTelepresenceMountRoot(*s.certPath)

	registerer := metrics.DefaultRegisterer()

	mysqlClient, err := s.createMySQLClient(ctx, registerer, logger)
	if err != nil {
		return err
	}

	// Operational-database storages (separate from the data-warehouse clients).
	var dbClient database.Client
	var postgresClient postgres.Client
	var pushStorage pushstorage.PushStorage
	var segmentStorage v2fs.SegmentStorage
	var segmentUserStorage v2fs.SegmentUserStorage
	var fluiStorage v2fs.FeatureLastUsedInfoStorage
	var accountStorage accstorage.AccountStorage
	var auditLogStorage v2als.AuditLogStorage
	var adminAuditLogStorage v2als.AdminAuditLogStorage
	var experimentStorage operationalstorage.ExperimentStorage
	var autoOpsRuleStorage operationalstorage.AutoOpsRuleStorage
	if *s.operationalDatabaseType == "postgres" {
		if *s.postgresUser == "" || *s.postgresHost == "" || *s.postgresDBName == "" {
			return fmt.Errorf("postgres-user, postgres-host, and postgres-db-name are required when storage-type=postgres")
		}
		postgresClient, err = s.createPostgresClient(ctx, registerer, logger)
		if err != nil {
			return err
		}
		dbClient = database.NewPostgresStorageClient(postgresClient)
		pushStorage = pushstorage.NewPostgresPushStorage(postgresClient)
		segmentStorage = featurepostgres.NewSegmentStorage(postgresClient)
		segmentUserStorage = featurepostgres.NewSegmentUserStorage(postgresClient)
		fluiStorage = featurepostgres.NewFeatureLastUsedInfoStorage(postgresClient)
		accountStorage = accountpostgres.NewAccountStorage(postgresClient)
		auditLogStorage = auditlogpostgres.NewAuditLogStorage(postgresClient)
		adminAuditLogStorage = auditlogpostgres.NewAdminAuditLogStorage(postgresClient)
		experimentStorage = oppostgres.NewExperimentStorage(postgresClient)
		autoOpsRuleStorage = oppostgres.NewAutoOpsRuleStorage(postgresClient)
	} else {
		dbClient = database.NewMySQLStorageClient(mysqlClient)
		pushStorage = pushstorage.NewMySQLPushStorage(mysqlClient)
		segmentStorage = featuremysql.NewSegmentStorage(mysqlClient)
		segmentUserStorage = featuremysql.NewSegmentUserStorage(mysqlClient)
		fluiStorage = featuremysql.NewFeatureLastUsedInfoStorage(mysqlClient)
		accountStorage = accountmysql.NewAccountStorage(mysqlClient)
		auditLogStorage = auditlogmysql.NewAuditLogStorage(mysqlClient)
		adminAuditLogStorage = auditlogmysql.NewAdminAuditLogStorage(mysqlClient)
		experimentStorage = opmysql.NewExperimentStorage(mysqlClient)
		autoOpsRuleStorage = opmysql.NewAutoOpsRuleStorage(mysqlClient)
	}

	creds, err := client.NewPerRPCCredentials(*s.serviceTokenPath)
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

	nonPersistentRedisClient, err := redisv3.NewClient(
		*s.nonPersistentRedisAddr,
		redisv3.WithPoolSize(*s.nonPersistentRedisPoolMaxActive),
		redisv3.WithMinIdleConns(*s.nonPersistentRedisPoolMaxIdle),
		redisv3.WithServerName(*s.nonPersistentRedisServerName),
		redisv3.WithRedisMode(redisv3.RedisMode(*s.nonPersistentRedisMode)),
		redisv3.WithMetrics(registerer),
		redisv3.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	persistentRedisClient, err := redisv3.NewClient(
		*s.persistentRedisAddr,
		redisv3.WithPoolSize(*s.persistentRedisPoolMaxActive),
		redisv3.WithMinIdleConns(*s.persistentRedisPoolMaxIdle),
		redisv3.WithServerName(*s.persistentRedisServerName),
		redisv3.WithRedisMode(redisv3.RedisMode(*s.persistentRedisMode)),
		redisv3.WithMetrics(registerer),
		redisv3.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	slackNotifier := notifier.NewSlackNotifier(*s.webURL)

	notificationSender := notificationsender.NewSender(
		notificationClient,
		[]notifier.Notifier{slackNotifier},
		notificationsender.WithMetrics(registerer),
		notificationsender.WithLogger(logger),
	)

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

	// Load subscribers config once. We need it both to construct the
	// cache-invalidation publisher (which reuses the cacheRefresher
	// processor's PubSub backend settings) and to drive the multi-pubsub
	// dispatch in startMultiPubSub. Tolerate a missing/unreadable file
	// the same way startMultiPubSub historically did.
	subscriberConfigs, err := s.loadSubscriberConfigurations(logger)
	if err != nil {
		return err
	}

	// Create the cache-invalidation publisher up front so it can be wired
	// into the cache-refresher processor below. The PubSub backend is
	// taken from the cacheRefresher processor's subscribers config — that
	// processor is the one consuming domain events and triggering the
	// refresh, so its publisher inherently uses the same backend in the
	// same pod. This avoids parallel CLI flags for the same backend.
	cacheInvalidationPublisher, cacheInvalidationCleanup, err := s.createCacheInvalidationPublisher(
		ctx, subscriberConfigs, registerer, logger,
	)
	if err != nil {
		return err
	}
	if cacheInvalidationCleanup != nil {
		defer cacheInvalidationCleanup()
	}

	pubSubProcessors, dwhCleanup, err := s.registerPubSubProcessorMap(
		ctx,
		environmentClient,
		mysqlClient,
		dbClient,
		segmentStorage,
		segmentUserStorage,
		fluiStorage,
		pushStorage,
		auditLogStorage,
		adminAuditLogStorage,
		accountStorage,
		experimentStorage,
		autoOpsRuleStorage,
		persistentRedisClient,
		nonPersistentRedisClient,
		experimentClient,
		featureClient,
		batchClient,
		autoOpsClient,
		notificationSender,
		cacheInvalidationPublisher,
		registerer,
		logger,
	)
	if err != nil {
		return err
	}
	if dwhCleanup != nil {
		defer dwhCleanup()
	}

	multiPubSub, err := s.startMultiPubSub(ctx, pubSubProcessors, subscriberConfigs, registerer, logger)
	if err != nil {
		return err
	}

	// healthCheckService
	// Use a dedicated context so we can stop the health checker goroutine cleanly during shutdown
	healthCheckCtx, healthCheckCancel := context.WithCancel(ctx)
	defer healthCheckCancel()

	restHealthChecker := health.NewRestChecker(
		"", "",
		health.WithTimeout(healthCheckTimeout),
		health.WithCheck("metrics", metrics.Check),
	)
	go restHealthChecker.Run(healthCheckCtx)
	// healthcheckService
	healthcheckServer := rest.NewServer(
		*s.certPath, *s.keyPath,
		rest.WithLogger(logger),
		rest.WithService(restHealthChecker),
		rest.WithMetrics(registerer),
		rest.WithPort(*s.port),
	)
	go healthcheckServer.Run()

	defer func() {
		shutdownStartTime := time.Now()

		// Mark as unhealthy so readiness probes fail
		// This ensures Kubernetes readiness probe fails on next check,
		// preventing new traffic from being routed to this pod.
		restHealthChecker.Stop()

		// Stop PubSub subscription
		// This stops receiving new messages and allows in-flight messages to be processed.
		multiPubSub.Stop()

		// Close clients
		// These are fast cleanup operations that can run asynchronously.
		go notificationClient.Close()
		go experimentClient.Close()
		go environmentClient.Close()
		go featureClient.Close()
		go autoOpsClient.Close()
		go batchClient.Close()
		go mysqlClient.Close()
		if postgresClient != nil {
			go postgresClient.Close()
		}
		go nonPersistentRedisClient.Close()
		go persistentRedisClient.Close()

		// Log total shutdown duration
		logger.Info("Graceful shutdown sequence completed",
			zap.Duration("total_elapsed", time.Since(shutdownStartTime)),
		)
	}()

	<-ctx.Done()
	return nil
}

// loadSubscriberConfigurations reads the subscribers config file and
// parses it into a map keyed by processor name. A missing or empty file
// is tolerated (returns an empty map and a logged warning) to match the
// historical behaviour of startMultiPubSub.
func (s *server) loadSubscriberConfigurations(
	logger *zap.Logger,
) (map[string]subscriber.Configuration, error) {
	bytes, err := os.ReadFile(*s.subscriberConfig)
	if err != nil {
		logger.Warn("subscriber: failed to read subscriber config",
			zap.String("path", *s.subscriberConfig),
			zap.Error(err),
		)
		return map[string]subscriber.Configuration{}, nil
	}
	configs := map[string]subscriber.Configuration{}
	if err := json.Unmarshal(bytes, &configs); err != nil {
		logger.Error("subscriber: failed to unmarshal subscriber config",
			zap.Error(err),
		)
		return nil, err
	}
	return configs, nil
}

// resolveCacheInvalidationConfig returns the cacheRefresher subscriber
// configuration if cache-invalidation announcements should be enabled
// for this pod. Announcements are enabled iff the cacheRefresher block
// is present AND its CacheInvalidationTopic is non-empty.
//
// Extracted as a free function so it can be unit-tested without spinning
// up a real PubSub backend.
func resolveCacheInvalidationConfig(
	subscriberConfigs map[string]subscriber.Configuration,
) (subscriber.Configuration, bool) {
	conf, ok := subscriberConfigs[processor.CacheRefresherName]
	if !ok {
		return subscriber.Configuration{}, false
	}
	if conf.CacheInvalidationTopic == "" {
		return subscriber.Configuration{}, false
	}
	return conf, true
}

// createCacheInvalidationPublisher builds a publisher for the
// cache-invalidation announcement topic, reusing the PubSub backend
// settings from the cacheRefresher processor's subscribers config so the
// publisher and the consumer that triggers it always agree on backend.
// The destination topic itself comes from the same config block
// (CacheInvalidationTopic field).
//
// The returned cleanup function (may be nil) stops the publisher, closes
// the factory client, and releases any backend-specific resources (e.g.
// a dedicated Redis client for the Redis Streams backend).
//
// Returns (nil, nil, nil) when announcements are disabled (no
// cacheRefresher block, or its CacheInvalidationTopic is empty).
func (s *server) createCacheInvalidationPublisher(
	ctx context.Context,
	subscriberConfigs map[string]subscriber.Configuration,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (publisher.Publisher, func(), error) {
	conf, enabled := resolveCacheInvalidationConfig(subscriberConfigs)
	if !enabled {
		logger.Info(
			"subscriber: cache-invalidation announcements disabled (no cacheRefresher block or empty cacheInvalidationTopic)",
		)
		return nil, nil, nil
	}
	pubSubType := factory.PubSubType(conf.PubSubType)
	if pubSubType == "" {
		pubSubType = factory.PubSubType(subscriber.DefaultPubSubType)
	}
	factoryOpts := []factory.Option{
		factory.WithPubSubType(pubSubType),
		factory.WithMetrics(registerer),
		factory.WithLogger(logger),
	}
	// backendCleanup releases backend-specific resources (currently the
	// Redis client for the Redis Streams backend). The factory client
	// itself is closed by the returned cleanup below, regardless of
	// backend.
	var backendCleanup func()
	switch pubSubType {
	case factory.Google:
		factoryOpts = append(factoryOpts, factory.WithProjectID(conf.Project))
	case factory.RedisStream:
		redisClient, err := redisv3.NewClient(
			conf.RedisAddr,
			redisv3.WithPoolSize(conf.RedisPoolSize),
			redisv3.WithMinIdleConns(conf.RedisMinIdle),
			redisv3.WithServerName(conf.RedisServerName),
			redisv3.WithRedisMode(redisv3.RedisMode(conf.RedisMode)),
			redisv3.WithMetrics(registerer),
			redisv3.WithLogger(logger),
		)
		if err != nil {
			return nil, nil, err
		}
		factoryOpts = append(factoryOpts, factory.WithRedisClient(redisClient))
		if conf.RedisPartitionCount > 0 {
			factoryOpts = append(factoryOpts, factory.WithPartitionCount(conf.RedisPartitionCount))
		}
		backendCleanup = func() { _ = redisClient.Close() }
	}
	client, err := factory.NewClient(ctx, factoryOpts...)
	if err != nil {
		if backendCleanup != nil {
			backendCleanup()
		}
		return nil, nil, err
	}
	pub, err := client.CreatePublisher(conf.CacheInvalidationTopic)
	if err != nil {
		if cerr := client.Close(); cerr != nil {
			logger.Error("subscriber: failed to close cache invalidation pubsub client during error cleanup",
				zap.Error(cerr),
			)
		}
		if backendCleanup != nil {
			backendCleanup()
		}
		logger.Error("subscriber: failed to create cache invalidation publisher",
			zap.String("topic", conf.CacheInvalidationTopic),
			zap.Error(err),
		)
		return nil, nil, err
	}
	cleanup := func() {
		pub.Stop()
		if err := client.Close(); err != nil {
			logger.Error("subscriber: failed to close cache invalidation pubsub client",
				zap.Error(err),
			)
		}
		if backendCleanup != nil {
			backendCleanup()
		}
	}
	return pub, cleanup, nil
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

func (s *server) createPostgresClient(
	ctx context.Context,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (postgres.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return postgres.NewClient(
		ctx,
		*s.postgresUser, *s.postgresPass, *s.postgresHost,
		*s.postgresPort,
		*s.postgresDBName,
		postgres.WithLogger(logger),
		postgres.WithMetrics(registerer),
	)
}

// closeDWHWriter closes a data-warehouse event writer if it implements io.Closer,
// logging any error. BigQuery writers hold a live managed stream and must be closed.
func closeDWHWriter(w interface{}, name string, logger *zap.Logger) {
	if c, ok := w.(io.Closer); ok {
		if err := c.Close(); err != nil {
			logger.Error("subscriber: failed to close "+name, zap.Error(err))
		}
	}
}

// initDataWarehouseStorages initializes every data-warehouse storage based on the resolved
// data-warehouse config. The data warehouse is always separate from the operational database;
// the processor receives ready storage interfaces and depends on no DWH client or dialect.
// Returns the evaluation/goal event writers, the goal event storage (for retries), and a
// cleanup closing any dedicated client. cleanup is nil when nothing needs closing.
func (s *server) initDataWarehouseStorages(
	ctx context.Context,
	dwhConfig DataWarehouseConfig,
	mainMySQLClient mysql.Client,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (dwhstorage.EvalEventWriter, dwhstorage.GoalEventWriter, ecdwh.EventStorage, func(), error) {
	switch dwhConfig.Type {
	case "mysql":
		client := mainMySQLClient
		var cleanup func()
		if dwhConfig.MySQL.UseMainConnection {
			logger.Info("Using main MySQL connection for data warehouse")
		} else {
			dedicated, err := createDWHMySQLClient(ctx, &dwhConfig.MySQL, logger)
			if err != nil {
				return nil, nil, nil, nil, fmt.Errorf("failed to create dedicated MySQL client: %w", err)
			}
			client = dedicated
			cleanup = func() { dedicated.Close() }
		}
		return dwhstorage.NewEvalEventWriter(dwhmysql.NewEvaluationEventStorage(client)),
			dwhstorage.NewGoalEventWriter(dwhmysql.NewGoalEventStorage(client)),
			ecmysql.NewMySQLEventStorage(client, logger),
			cleanup,
			nil
	case "postgres":
		client, err := createDWHPostgresClient(ctx, &dwhConfig.Postgres, logger)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("failed to create dedicated Postgres client: %w", err)
		}
		return dwhstorage.NewEvalEventWriter(dwhpostgres.NewEvaluationEventStorage(client)),
			dwhstorage.NewGoalEventWriter(dwhpostgres.NewGoalEventStorage(client)),
			ecpostgres.NewPostgresEventStorage(client, logger),
			func() { client.Close() },
			nil
	default:
		// BigQuery (default)
		project := dwhConfig.BigQuery.Project
		dataset := dwhConfig.BigQuery.Dataset
		batchSize := dwhConfig.BatchSize
		evalWriter, err := dwhbigquery.NewEvaluationEventWriter(ctx, logger, project, dataset, batchSize, registerer)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		goalWriter, err := dwhbigquery.NewGoalEventWriter(ctx, logger, project, dataset, batchSize, registerer)
		if err != nil {
			// evalWriter already holds a live BigQuery managed stream; close it before bailing out.
			closeDWHWriter(evalWriter, "evaluation event writer", logger)
			return nil, nil, nil, nil, err
		}
		eventQuerier, err := bqquerier.NewClient(
			ctx,
			project,
			dwhConfig.BigQuery.Location,
			bqquerier.WithLogger(logger),
			bqquerier.WithMetrics(registerer),
		)
		if err != nil {
			// Both writers are open at this point; close them before bailing out.
			closeDWHWriter(evalWriter, "evaluation event writer", logger)
			closeDWHWriter(goalWriter, "goal event writer", logger)
			return nil, nil, nil, nil, err
		}
		cleanup := func() {
			closeDWHWriter(evalWriter, "evaluation event writer", logger)
			closeDWHWriter(goalWriter, "goal event writer", logger)
			if err := eventQuerier.Close(); err != nil {
				logger.Error("subscriber: failed to close data warehouse querier", zap.Error(err))
			}
		}
		return evalWriter,
			goalWriter,
			ecbigquery.NewBigQueryEventStorage(eventQuerier, dataset, logger),
			cleanup,
			nil
	}
}

// createDWHMySQLClient creates a dedicated MySQL client for the data warehouse.
func createDWHMySQLClient(
	ctx context.Context,
	config *MySQLConfig,
	logger *zap.Logger,
) (mysql.Client, error) {
	if config == nil {
		return nil, fmt.Errorf("mysql config is nil")
	}
	if config.Host == "" || config.Database == "" || config.User == "" {
		return nil, fmt.Errorf("mysql host, database, and user are required for dedicated connection")
	}
	port := config.Port
	if port == 0 {
		port = 3306 // Default MySQL port
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	client, err := mysql.NewClient(
		ctx,
		config.User,
		config.Password,
		config.Host,
		port,
		config.Database,
		mysql.WithLogger(logger),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create MySQL client: %w", err)
	}
	logger.Info("Created dedicated MySQL client for data warehouse",
		zap.String("host", config.Host),
		zap.Int("port", port),
		zap.String("database", config.Database),
		zap.String("user", config.User),
	)
	return client, nil
}

// createDWHPostgresClient creates a dedicated Postgres client for the data warehouse.
func createDWHPostgresClient(
	ctx context.Context,
	config *PostgresConfig,
	logger *zap.Logger,
) (postgres.Client, error) {
	if config == nil {
		return nil, fmt.Errorf("postgres config is nil")
	}
	if config.Host == "" || config.Database == "" || config.User == "" {
		return nil, fmt.Errorf("postgres host, database, and user are required for dedicated connection")
	}
	port := config.Port
	if port == 0 {
		port = 5432 // Default Postgres port
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	client, err := postgres.NewClient(
		ctx,
		config.User,
		config.Password,
		config.Host,
		port,
		config.Database,
		postgres.WithLogger(logger),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Postgres client: %w", err)
	}
	logger.Info("Created dedicated Postgres client for data warehouse",
		zap.String("host", config.Host),
		zap.Int("port", port),
		zap.String("database", config.Database),
		zap.String("user", config.User),
	)
	return client, nil
}

func (s *server) startMultiPubSub(
	ctx context.Context,
	processors *processor.PubSubProcessors,
	subscriberConfigs map[string]subscriber.Configuration,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (*subscriber.MultiSubscriber, error) {
	multiSubscriber := subscriber.NewMultiSubscriber(
		subscriber.WithLogger(logger),
		subscriber.WithMetrics(registerer),
	)
	for name, config := range subscriberConfigs {
		p, err := processors.GetProcessorByName(name)
		if err != nil {
			logger.Warn(
				"subscriber: processor not found during startup. It could be because the processor is not registered yet.",
				zap.String("name", name),
				zap.Error(err),
			)
			// since we will keep old and new configmap at the same time during canary release,
			// we should skip the error, just log it here
			continue
		}
		multiSubscriber.AddSubscriber(subscriber.NewPubSubSubscriber(
			name, config, p,
			subscriber.WithLogger(logger),
			subscriber.WithMetrics(registerer),
		))
	}
	onDemandSubscriberConfigBytes, err := os.ReadFile(*s.onDemandSubscriberConfig)
	if err != nil {
		logger.Error("subscriber: failed to read subscriber config", zap.Error(err))
	} else {
		var onDemandConfigMap map[string]subscriber.OnDemandConfiguration
		if err := json.Unmarshal(onDemandSubscriberConfigBytes, &onDemandConfigMap); err != nil {
			logger.Error("subscriber: failed to unmarshal onDemand subscriber config",
				zap.Error(err),
			)
			return nil, err
		}
		for name, config := range onDemandConfigMap {
			p, err := processors.GetProcessorByName(name)
			if err != nil {
				logger.Error("subscriber: onDemand processor not found",
					zap.String("name", name),
					zap.Error(err),
				)
				// since we will keep old and new configmap at the same time during canary release,
				// we should skip the error, just log it here
				continue
			}
			multiSubscriber.AddSubscriber(subscriber.NewOnDemandSubscriber(
				name, config, p.(subscriber.OnDemandProcessor),
				subscriber.WithLogger(logger),
				subscriber.WithMetrics(registerer),
			))
		}
	}

	multiSubscriber.Start(ctx)
	return multiSubscriber, nil
}

func (s *server) registerPubSubProcessorMap(
	ctx context.Context,
	environmentClient environmentclient.Client,
	mysqlClient mysql.Client,
	dbClient database.Client,
	segmentStorage v2fs.SegmentStorage,
	segmentUserStorage v2fs.SegmentUserStorage,
	fluiStorage v2fs.FeatureLastUsedInfoStorage,
	pushStorage pushstorage.PushStorage,
	auditLogStorage v2als.AuditLogStorage,
	adminAuditLogStorage v2als.AdminAuditLogStorage,
	accountStorage accstorage.AccountStorage,
	experimentStorage operationalstorage.ExperimentStorage,
	autoOpsRuleStorage operationalstorage.AutoOpsRuleStorage,
	persistentRedisClient redisv3.Client,
	nonPersistentRedisClient redisv3.Client,
	exClient experimentclient.Client,
	ftClient featureclient.Client,
	batchClient btclient.Client,
	opsClient autoopsclient.Client,
	sender notificationsender.Sender,
	cacheInvalidationPublisher publisher.Publisher,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (*processor.PubSubProcessors, func(), error) {
	processors := processor.NewPubSubProcessors(registerer)
	// dwhCleanup closes any dedicated data-warehouse client created below.
	var dwhCleanup func()

	processorsConfigBytes, err := os.ReadFile(*s.processorsConfig)
	if err != nil {
		logger.Error("subscriber: failed to read processors config", zap.Error(err))
	} else {
		var processorsConfigMap map[string]interface{}
		if err := json.Unmarshal(processorsConfigBytes, &processorsConfigMap); err != nil {
			logger.Error("subscriber: failed to unmarshal processors config",
				zap.Error(err),
			)
			return nil, nil, err
		}
		auditLogPersister, err := processor.NewAuditLogPersister(
			processorsConfigMap[processor.AuditLogPersisterName],
			auditLogStorage,
			adminAuditLogStorage,
			logger,
		)
		if err != nil {
			return nil, nil, err
		}
		processors.RegisterProcessor(processor.AuditLogPersisterName, auditLogPersister)

		processors.RegisterProcessor(
			processor.DomainEventInformerName,
			processor.NewDomainEventInformer(environmentClient, sender, logger),
		)

		nonPersistentRedisCache := cachev3.NewRedisCache(nonPersistentRedisClient)
		processors.RegisterProcessor(
			processor.CacheRefresherName,
			processor.NewCacheRefresher(
				ftClient,
				exClient,
				opsClient,
				accountStorage,
				cachev3.NewFeaturesCache(nonPersistentRedisCache, 0),
				cachev3.NewSegmentUsersCache(nonPersistentRedisCache, 0),
				cachev3.NewEnvironmentAPIKeyCache(nonPersistentRedisCache, 0),
				cachev3.NewExperimentsCache(nonPersistentRedisCache),
				cachev3.NewAutoOpsRulesCache(nonPersistentRedisCache),
				cacheInvalidationPublisher,
				logger,
			),
		)

		segmentPersister, err := processor.NewSegmentUserPersister(
			processorsConfigMap[processor.SegmentUserPersisterName],
			batchClient,
			dbClient,
			segmentStorage,
			segmentUserStorage,
			registerer,
			logger,
		)
		if err != nil {
			return nil, nil, err
		}
		processors.RegisterProcessor(
			processor.SegmentUserPersisterName,
			segmentPersister,
		)

		if *s.demoSiteEnabled {
			demoOrganizationCreationNotifier := processor.NewDemoOrganizationCreationNotifier(
				processorsConfigMap[processor.DemoOrganizationCreationNotifierName],
				*s.webURL,
				logger,
			)
			processors.RegisterProcessor(
				processor.DemoOrganizationCreationNotifierName,
				demoOrganizationCreationNotifier,
			)
		}

		// Email service
		emailConfig, err := s.readEmailConfig(logger)
		if err != nil {
			return nil, nil, err
		}
		emailService, err := email.NewService(*emailConfig, logger)
		if err != nil {
			logger.Error("Failed to create email service", zap.Error(err))
			return nil, nil, err
		}

		// Email sender processor
		emailSender := processor.NewEmailSender(
			processorsConfigMap[processor.EmailSenderName],
			emailService,
			logger,
		)
		processors.RegisterProcessor(
			processor.EmailSenderName,
			emailSender,
		)

		redisCache := cachev3.NewRedisCache(persistentRedisClient)
		evaluationCountEventPersister, err := processor.NewEvaluationCountEventPersister(
			ctx,
			processorsConfigMap[processor.EvaluationCountEventPersisterName],
			fluiStorage,
			redisCache,
			cachev3.NewUserAttributesCache(redisCache),
			cachev3.NewDAUCache(redisCache),
			logger,
		)
		if err != nil {
			return nil, nil, err
		}
		processors.RegisterProcessor(
			processor.EvaluationCountEventPersisterName,
			evaluationCountEventPersister,
		)

		processors.RegisterProcessor(
			processor.PushSenderName,
			processor.NewPushSender(
				ftClient,
				batchClient,
				pushStorage,
				logger,
			),
		)

		processors.RegisterProcessor(
			processor.MetricsEventPersisterName,
			processor.NewMetricsEventPersister(
				registerer,
				logger,
			),
		)
	}

	onDemandProcessorsConfigBytes, err := os.ReadFile(*s.onDemandProcessorsConfig)
	if err != nil {
		logger.Error("subscriber: failed to read onDemand processors config", zap.Error(err))
	} else {
		var onDemandProcessorsConfigMap map[string]interface{}
		if err := json.Unmarshal(onDemandProcessorsConfigBytes, &onDemandProcessorsConfigMap); err != nil {
			logger.Error("subscriber: failed to unmarshal onDemand processors config",
				zap.Error(err),
			)
			return nil, nil, err
		}

		// Initialize the data-warehouse client(s) here, separately from the operational
		// database client. The DWH backend (bigquery, mysql, postgres) is resolved from
		// the events DWH persister config; eval and goal persisters share it.
		dwhConfig, err := parseDWHConfig(
			onDemandProcessorsConfigMap[processor.EvaluationCountEventDWHPersisterName],
		)
		if err != nil {
			logger.Error("subscriber: failed to parse data warehouse config", zap.Error(err))
			return nil, nil, err
		}
		dwhLocation, err := locale.GetLocation(dwhConfig.Timezone)
		if err != nil {
			logger.Error("subscriber: failed to resolve data warehouse timezone", zap.Error(err))
			return nil, nil, err
		}
		evalEventWriter, goalEventWriter, goalEventStorage, cleanup, err := s.initDataWarehouseStorages(
			ctx, dwhConfig, mysqlClient, registerer, logger,
		)
		if err != nil {
			logger.Error("subscriber: failed to initialize data warehouse storages", zap.Error(err))
			return nil, nil, err
		}
		dwhCleanup = cleanup

		evaluationEventsDWHPersister, err := processor.NewEventsDWHPersister(
			ctx,
			onDemandProcessorsConfigMap[processor.EvaluationCountEventDWHPersisterName],
			evalEventWriter,
			goalEventWriter,
			goalEventStorage,
			experimentStorage,
			dwhLocation,
			nonPersistentRedisClient, // use non-persistent redis instance here
			persistentRedisClient,    // use persistent redis instance here for goal retry events
			exClient,
			ftClient,
			processor.EvaluationCountEventDWHPersisterName,
			logger,
		)
		if err != nil {
			if dwhCleanup != nil {
				dwhCleanup()
			}
			return nil, nil, err
		}
		processors.RegisterProcessor(
			processor.EvaluationCountEventDWHPersisterName,
			evaluationEventsDWHPersister,
		)

		goalEventsDWHPersister, err := processor.NewEventsDWHPersister(
			ctx,
			onDemandProcessorsConfigMap[processor.GoalCountEventDWHPersisterName],
			evalEventWriter,
			goalEventWriter,
			goalEventStorage,
			experimentStorage,
			dwhLocation,
			nonPersistentRedisClient, // use non-persistent redis instance here
			persistentRedisClient,    // use persistent redis instance here for goal retry events
			exClient,
			ftClient,
			processor.GoalCountEventDWHPersisterName,
			logger,
		)
		if err != nil {
			if dwhCleanup != nil {
				dwhCleanup()
			}
			return nil, nil, err
		}
		processors.RegisterProcessor(
			processor.GoalCountEventDWHPersisterName,
			goalEventsDWHPersister,
		)

		evaluationEventsOPSPersister, err := processor.NewEventsOPSPersister(
			ctx,
			onDemandProcessorsConfigMap[processor.EvaluationCountEventOPSPersisterName],
			autoOpsRuleStorage,
			persistentRedisClient, // use persistent redis instance here
			opsClient,
			ftClient,
			processor.EvaluationCountEventOPSPersisterName,
			logger,
		)
		if err != nil {
			if dwhCleanup != nil {
				dwhCleanup()
			}
			return nil, nil, err
		}
		processors.RegisterProcessor(processor.EvaluationCountEventOPSPersisterName, evaluationEventsOPSPersister)

		goalEventsOPSPersister, err := processor.NewEventsOPSPersister(
			ctx,
			onDemandProcessorsConfigMap[processor.GoalCountEventOPSPersisterName],
			autoOpsRuleStorage,
			persistentRedisClient, // use persistent redis instance here
			opsClient,
			ftClient,
			processor.GoalCountEventOPSPersisterName,
			logger,
		)
		if err != nil {
			if dwhCleanup != nil {
				dwhCleanup()
			}
			return nil, nil, err
		}
		processors.RegisterProcessor(processor.GoalCountEventOPSPersisterName, goalEventsOPSPersister)
	}

	return processors, dwhCleanup, nil
}

func (s *server) readEmailConfig(
	logger *zap.Logger,
) (*email.Config, error) {
	bytes, err := os.ReadFile(*s.emailConfigPath)
	if err != nil {
		logger.Error("Failed to read email config file",
			zap.Error(err),
		)
		return nil, err
	}
	config := email.Config{}
	if err = json.Unmarshal(bytes, &config); err != nil {
		logger.Error("Failed to unmarshal email config",
			zap.Error(err),
		)
		return nil, err
	}
	return &config, nil
}

func (s *server) insertTelepresenceMountRoot(path string) string {
	volumeRoot := os.Getenv("TELEPRESENCE_ROOT")
	if volumeRoot == "" {
		return path
	}
	return volumeRoot + path
}
