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
	"os"
	"time"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	accstorage "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	autoopsclient "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client"
	btclient "github.com/bucketeer-io/bucketeer/v2/pkg/batch/client"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cli"
	"github.com/bucketeer-io/bucketeer/v2/pkg/email"
	environmentclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/health"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	notificationclient "github.com/bucketeer-io/bucketeer/v2/pkg/notification/client"
	notificationsender "github.com/bucketeer-io/bucketeer/v2/pkg/notification/sender"
	"github.com/bucketeer-io/bucketeer/v2/pkg/notification/sender/notifier"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/factory"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	redisv3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rest"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/processor"
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
	// MySQL
	mysqlUser        *string
	mysqlPass        *string
	mysqlHost        *string
	mysqlPort        *int
	mysqlDBName      *string
	mysqlDBOpenConns *int
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
	// PubSub for the cache-invalidation announcement topic.
	// The subscriber service publishes to this topic after refreshing L2,
	// so each api pod can drop its L1 entry without re-reading MySQL.
	pubSubType                *string
	pubSubRedisServerName     *string
	pubSubRedisAddr           *string
	pubSubRedisPoolSize       *int
	pubSubRedisMinIdle        *int
	pubSubRedisPartitionCount *int
	pubSubRedisMode           *string
	cacheInvalidationTopic    *string
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
		mysqlUser:        cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:        cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:        cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:        cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName:      cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		mysqlDBOpenConns: cmd.Flag("mysql-db-open-conns", "MySQL open connections.").Required().Int(),
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
		pubSubType: cmd.Flag("pubsub-type",
			"Type of PubSub to use (google or redis-stream) for the cache-invalidation publisher.",
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
		pubSubRedisMode: cmd.Flag("pubsub-redis-mode",
			"PubSub Redis client mode: cluster, standalone, or auto.",
		).Default("auto").String(),
		cacheInvalidationTopic: cmd.Flag("cache-invalidation-topic",
			"PubSub topic on which the subscriber announces L2 cache refreshes "+
				"so api pods can drop their L1 (in-memory) entries. Both backends "+
				"auto-create the topic on first publish.",
		).Default("cache-invalidation").String(),
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

	// Create the cache-invalidation publisher up front so it can be wired
	// into the cache-refresher processor below. It is intentionally a
	// separate PubSub client from the per-processor consumers (each
	// processor brings its own client via the subscribers config); having
	// one shared publisher avoids creating N publishers per refresher type.
	cacheInvalidationPubSubClient, cacheInvalidationCleanup, err := s.createCacheInvalidationPubSubClient(
		ctx, registerer, logger,
	)
	if err != nil {
		return err
	}
	if cacheInvalidationCleanup != nil {
		defer cacheInvalidationCleanup()
	}

	var cacheInvalidationPublisher publisher.Publisher
	if cacheInvalidationPubSubClient != nil && *s.cacheInvalidationTopic != "" {
		cacheInvalidationPublisher, err = cacheInvalidationPubSubClient.CreatePublisher(*s.cacheInvalidationTopic)
		if err != nil {
			logger.Error("subscriber: failed to create cache invalidation publisher",
				zap.String("topic", *s.cacheInvalidationTopic),
				zap.Error(err),
			)
			return err
		}
		defer cacheInvalidationPublisher.Stop()
	}

	pubSubProcessors, err := s.registerPubSubProcessorMap(
		ctx,
		environmentClient,
		mysqlClient,
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

	multiPubSub, err := s.startMultiPubSub(ctx, pubSubProcessors, registerer, logger)
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

// createCacheInvalidationPubSubClient builds a PubSub client suitable for
// publishing cache-invalidation announcements. The returned cleanup
// function closes both the factory client (e.g. the Google Pub/Sub gRPC
// connection) and any backend-specific resources such as a dedicated
// Redis client.
//
// Returns (nil, nil, nil) when the cache-invalidation topic is empty
// (announcements disabled).
func (s *server) createCacheInvalidationPubSubClient(
	ctx context.Context,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (factory.Client, func(), error) {
	if *s.cacheInvalidationTopic == "" {
		return nil, nil, nil
	}
	pubSubType := factory.PubSubType(*s.pubSubType)
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
		factoryOpts = append(factoryOpts, factory.WithProjectID(*s.project))
	case factory.RedisStream:
		redisClient, err := redisv3.NewClient(
			*s.pubSubRedisAddr,
			redisv3.WithPoolSize(*s.pubSubRedisPoolSize),
			redisv3.WithMinIdleConns(*s.pubSubRedisMinIdle),
			redisv3.WithServerName(*s.pubSubRedisServerName),
			redisv3.WithRedisMode(redisv3.RedisMode(*s.pubSubRedisMode)),
			redisv3.WithMetrics(registerer),
			redisv3.WithLogger(logger),
		)
		if err != nil {
			return nil, nil, err
		}
		factoryOpts = append(factoryOpts, factory.WithRedisClient(redisClient))
		factoryOpts = append(factoryOpts, factory.WithPartitionCount(*s.pubSubRedisPartitionCount))
		backendCleanup = func() { _ = redisClient.Close() }
	}
	client, err := factory.NewClient(ctx, factoryOpts...)
	if err != nil {
		if backendCleanup != nil {
			backendCleanup()
		}
		return nil, nil, err
	}
	cleanup := func() {
		if err := client.Close(); err != nil {
			logger.Error("subscriber: failed to close cache invalidation pubsub client",
				zap.Error(err),
			)
		}
		if backendCleanup != nil {
			backendCleanup()
		}
	}
	return client, cleanup, nil
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

func (s *server) startMultiPubSub(
	ctx context.Context,
	processors *processor.PubSubProcessors,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (*subscriber.MultiSubscriber, error) {
	multiSubscriber := subscriber.NewMultiSubscriber(
		subscriber.WithLogger(logger),
		subscriber.WithMetrics(registerer),
	)
	subscriberConfigBytes, err := os.ReadFile(*s.subscriberConfig)
	if err != nil {
		logger.Error("subscriber: failed to read subscriber config", zap.Error(err))
	} else {
		var configMap map[string]subscriber.Configuration
		if err := json.Unmarshal(subscriberConfigBytes, &configMap); err != nil {
			logger.Error("subscriber: failed to unmarshal subscriber config",
				zap.Error(err),
			)
			return nil, err
		}
		for name, config := range configMap {
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
) (*processor.PubSubProcessors, error) {
	processors := processor.NewPubSubProcessors(registerer)

	processorsConfigBytes, err := os.ReadFile(*s.processorsConfig)
	if err != nil {
		logger.Error("subscriber: failed to read processors config", zap.Error(err))
	} else {
		var processorsConfigMap map[string]interface{}
		if err := json.Unmarshal(processorsConfigBytes, &processorsConfigMap); err != nil {
			logger.Error("subscriber: failed to unmarshal processors config",
				zap.Error(err),
			)
			return nil, err
		}
		auditLogPersister, err := processor.NewAuditLogPersister(
			processorsConfigMap[processor.AuditLogPersisterName],
			mysqlClient,
			logger,
		)
		if err != nil {
			return nil, err
		}
		processors.RegisterProcessor(processor.AuditLogPersisterName, auditLogPersister)

		processors.RegisterProcessor(
			processor.DomainEventInformerName,
			processor.NewDomainEventInformer(environmentClient, sender, logger),
		)

		nonPersistentRedisCache := cachev3.NewRedisCache(nonPersistentRedisClient)
		processors.RegisterProcessor(
			processor.CacheEvictionName,
			processor.NewCacheRefresher(
				ftClient,
				exClient,
				opsClient,
				accstorage.NewAccountStorage(mysqlClient),
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
			mysqlClient,
			registerer,
			logger,
		)
		if err != nil {
			return nil, err
		}
		processors.RegisterProcessor(
			processor.SegmentUserPersisterName,
			segmentPersister,
		)

		userEventPersister, err := processor.NewUserEventPersister(
			processorsConfigMap[processor.UserEventPersisterName],
			mysqlClient,
			logger,
		)
		if err != nil {
			return nil, err
		}
		processors.RegisterProcessor(
			processor.UserEventPersisterName,
			userEventPersister,
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
			return nil, err
		}
		emailService, err := email.NewService(*emailConfig, logger)
		if err != nil {
			logger.Error("Failed to create email service", zap.Error(err))
			return nil, err
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
			mysqlClient,
			redisCache,
			cachev3.NewUserAttributesCache(redisCache),
			cachev3.NewDAUCache(redisCache),
			logger,
		)
		if err != nil {
			return nil, err
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
				mysqlClient,
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
			return nil, err
		}

		evaluationEventsDWHPersister, err := processor.NewEventsDWHPersister(
			ctx,
			onDemandProcessorsConfigMap[processor.EvaluationCountEventDWHPersisterName],
			mysqlClient,
			nonPersistentRedisClient, // use non-persistent redis instance here
			persistentRedisClient,    // use persistent redis instance here for goal retry events
			exClient,
			ftClient,
			processor.EvaluationCountEventDWHPersisterName,
			registerer,
			logger,
		)
		if err != nil {
			return nil, err
		}
		processors.RegisterProcessor(
			processor.EvaluationCountEventDWHPersisterName,
			evaluationEventsDWHPersister,
		)

		goalEventsDWHPersister, err := processor.NewEventsDWHPersister(
			ctx,
			onDemandProcessorsConfigMap[processor.GoalCountEventDWHPersisterName],
			mysqlClient,
			nonPersistentRedisClient, // use non-persistent redis instance here
			persistentRedisClient,    // use persistent redis instance here for goal retry events
			exClient,
			ftClient,
			processor.GoalCountEventDWHPersisterName,
			registerer,
			logger,
		)
		if err != nil {
			return nil, err
		}
		processors.RegisterProcessor(
			processor.GoalCountEventDWHPersisterName,
			goalEventsDWHPersister,
		)

		evaluationEventsOPSPersister, err := processor.NewEventsOPSPersister(
			ctx,
			onDemandProcessorsConfigMap[processor.EvaluationCountEventOPSPersisterName],
			mysqlClient,
			persistentRedisClient, // use persistent redis instance here
			opsClient,
			ftClient,
			processor.EvaluationCountEventOPSPersisterName,
			logger,
		)
		if err != nil {
			return nil, err
		}
		processors.RegisterProcessor(processor.EvaluationCountEventOPSPersisterName, evaluationEventsOPSPersister)

		goalEventsOPSPersister, err := processor.NewEventsOPSPersister(
			ctx,
			onDemandProcessorsConfigMap[processor.GoalCountEventOPSPersisterName],
			mysqlClient,
			persistentRedisClient, // use persistent redis instance here
			opsClient,
			ftClient,
			processor.GoalCountEventOPSPersisterName,
			logger,
		)
		if err != nil {
			return nil, err
		}
		processors.RegisterProcessor(processor.GoalCountEventOPSPersisterName, goalEventsOPSPersister)
	}

	return processors, nil
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
