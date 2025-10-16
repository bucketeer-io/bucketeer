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

package cmd

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	auditlogclient "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/client"
	autoopsclient "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cli"
	coderefclient "github.com/bucketeer-io/bucketeer/v2/pkg/coderef/client"
	environmentclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	eventcounterclient "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/client"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/health"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	notificationclient "github.com/bucketeer-io/bucketeer/v2/pkg/notification/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/factory"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	pushclient "github.com/bucketeer-io/bucketeer/v2/pkg/push/client"
	redisv3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rest"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc/gateway"
	tagclient "github.com/bucketeer-io/bucketeer/v2/pkg/tag/client"
	teamclient "github.com/bucketeer-io/bucketeer/v2/pkg/team/client"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

const (
	command = "server"

	// Shutdown timing for graceful termination:
	// During normal operations (rolling updates, scale down), the pod gets the full
	// terminationGracePeriodSeconds (60s). During Spot VM preemption, kubelet enforces
	// a best-effort 15s limit. We optimize for the common case (normal operations).
	// See: https://cloud.google.com/kubernetes-engine/docs/concepts/spot-vms
	propagationDelay      = 10 * time.Second
	serverShutDownTimeout = 30 * time.Second
	grpcStopTimeout       = 5 * time.Second
)

type server struct {
	*kingpin.CmdClause
	port                   *int
	grpcGatewayPort        *int
	project                *string
	goalTopic              *string
	goalTopicProject       *string
	evaluationTopic        *string
	evaluationTopicProject *string
	userTopic              *string
	metricsTopic           *string
	publishNumGoroutines   *int
	publishTimeout         *time.Duration
	featureService         *string
	accountService         *string
	codeRefService         *string
	pushService            *string
	auditLogService        *string
	tagService             *string
	teamService            *string
	notificationService    *string
	experimentService      *string
	environmentService     *string
	eventCounterService    *string
	redisServerName        *string
	redisAddr              *string
	certPath               *string
	keyPath                *string
	serviceTokenPath       *string
	redisPoolMaxIdle       *int
	redisPoolMaxActive     *int
	oldestEventTimestamp   *time.Duration
	furthestEventTimestamp *time.Duration
	// PubSub configurations
	pubSubType                *string
	pubSubRedisServerName     *string
	pubSubRedisAddr           *string
	pubSubRedisPoolSize       *int
	pubSubRedisMinIdle        *int
	pubSubRedisPartitionCount *int
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the gRPC server")
	server := &server{
		CmdClause:       cmd,
		port:            cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		grpcGatewayPort: cmd.Flag("grpc-gateway-port", "Port to bind to for gRPC-gateway.").Default("9089").Int(),
		project:         cmd.Flag("project", "GCP Project id to use for PubSub.").Required().String(),
		goalTopic:       cmd.Flag("goal-topic", "Topic to use for publishing GoalEvent.").Required().String(),
		goalTopicProject: cmd.Flag(
			"goal-topic-project",
			"GCP Project id to use for PubSub to publish GoalEvent.",
		).String(),
		evaluationTopic: cmd.Flag(
			"evaluation-topic",
			"Topic to use for publishing EvaluationEvent.",
		).Required().String(),
		evaluationTopicProject: cmd.Flag(
			"evaluation-topic-project",
			"GCP Project id to use for PubSub to publish EvaluationEvent.",
		).String(),
		// FIXME: This flag will be required once user feature is fully released.
		userTopic:    cmd.Flag("user-topic", "Topic to use for publishing UserEvent.").String(),
		metricsTopic: cmd.Flag("metrics-topic", "Topic to use for publishing MetricsEvent.").String(),
		publishNumGoroutines: cmd.Flag(
			"publish-num-goroutines",
			"The number of goroutines for publishing.",
		).Default("0").Int(),
		publishTimeout: cmd.Flag(
			"publish-timeout",
			"The maximum time to publish a bundle of messages.",
		).Default("1m").Duration(),
		featureService: cmd.Flag(
			"feature-service",
			"bucketeer-feature-service address.",
		).Default("feature:9090").String(),
		accountService: cmd.Flag(
			"account-service",
			"bucketeer-account-service address.",
		).Default("account:9090").String(),
		pushService: cmd.Flag(
			"push-service",
			"bucketeer-push-service address.",
		).Default("push:9090").String(),
		codeRefService: cmd.Flag(
			"code-ref-service",
			"bucketeer-code-ref-service address.",
		).Default("code-ref:9090").String(),
		auditLogService: cmd.Flag(
			"audit-log-service",
			"bucketeer-audit-log-service address.",
		).Default("audit-log:9090").String(),
		tagService: cmd.Flag(
			"tag-service",
			"bucketeer-tag-service address.",
		).Default("tag:9090").String(),
		teamService: cmd.Flag(
			"team-service",
			"bucketeer-team-service address.",
		).Default("team:9090").String(),
		notificationService: cmd.Flag(
			"notification-service",
			"bucketeer-notification-service address.",
		).Default("notification:9090").String(),
		experimentService: cmd.Flag(
			"experiment-service",
			"bucketeer-experiment-service address.",
		).Default("experiment:9090").String(),
		environmentService: cmd.Flag(
			"environment-service",
			"bucketeer-environment-service address.",
		).Default("environment:9090").String(),
		eventCounterService: cmd.Flag(
			"event-counter-service",
			"bucketeer-event-counter-service address.",
		).Default("event-counter:9090").String(),
		redisServerName:  cmd.Flag("redis-server-name", "Name of the redis.").Required().String(),
		redisAddr:        cmd.Flag("redis-addr", "Address of the redis.").Required().String(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		redisPoolMaxIdle: cmd.Flag(
			"redis-pool-max-idle",
			"Maximum number of idle connections in the pool.",
		).Default("5").Int(),
		redisPoolMaxActive: cmd.Flag(
			"redis-pool-max-active",
			"Maximum number of connections allocated by the pool at a given time.",
		).Default("10").Int(),
		oldestEventTimestamp: cmd.Flag(
			"oldest-event-timestamp",
			"The duration of oldest event timestamp from processing time to allow.",
		).Default("744h").Duration(),
		furthestEventTimestamp: cmd.Flag(
			"furthest-event-timestamp",
			"The duration of furthest event timestamp from processing time to allow.",
		).Default("1h").Duration(),
		// PubSub configurations
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

	pubsubCtx, pubsubCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pubsubCancel()

	// Create PubSub client using the factory
	pubSubType := factory.PubSubType(*s.pubSubType)
	factoryOpts := []factory.Option{
		factory.WithPubSubType(pubSubType),
		factory.WithMetrics(registerer),
		factory.WithLogger(logger),
	}

	// Add provider-specific options
	switch pubSubType {
	case factory.Google:
		factoryOpts = append(factoryOpts, factory.WithProjectID(*s.project))
	case factory.RedisStream:
		redisClient, err := redisv3.NewClient(
			*s.pubSubRedisAddr,
			redisv3.WithPoolSize(*s.pubSubRedisPoolSize),
			redisv3.WithMinIdleConns(*s.pubSubRedisMinIdle),
			redisv3.WithServerName(*s.pubSubRedisServerName),
			redisv3.WithMetrics(registerer),
			redisv3.WithLogger(logger),
		)
		if err != nil {
			return err
		}
		factoryOpts = append(factoryOpts, factory.WithRedisClient(redisClient))
		factoryOpts = append(factoryOpts, factory.WithPartitionCount(*s.pubSubRedisPartitionCount))
	}

	pubsubClient, err := factory.NewClient(pubsubCtx, factoryOpts...)
	if err != nil {
		return err
	}

	var goalTopicProject string
	if *s.goalTopicProject == "" {
		goalTopicProject = *s.project
	} else {
		goalTopicProject = *s.goalTopicProject
	}
	goalPublisher, err := pubsubClient.CreatePublisherInProject(*s.goalTopic, goalTopicProject)
	if err != nil {
		return err
	}
	defer goalPublisher.Stop()

	var evaluationTopicProject string
	if *s.evaluationTopicProject == "" {
		evaluationTopicProject = *s.project
	} else {
		evaluationTopicProject = *s.evaluationTopicProject
	}
	evaluationPublisher, err := pubsubClient.CreatePublisherInProject(
		*s.evaluationTopic,
		evaluationTopicProject,
	)
	if err != nil {
		return nil
	}
	defer evaluationPublisher.Stop()

	// FIXME: This condition won't be necessary once user feature is fully released.
	var userPublisher publisher.Publisher
	if *s.userTopic != "" {
		userPublisher, err = pubsubClient.CreatePublisherInProject(*s.userTopic, *s.project)
		if err != nil {
			return err
		}
		defer userPublisher.Stop()
	}

	// FIXME: This condition won't be necessary once user feature is fully released.
	var metricsPublisher publisher.Publisher
	if *s.metricsTopic != "" {
		metricsPublisher, err = pubsubClient.CreatePublisherInProject(*s.metricsTopic, *s.project)
		if err != nil {
			return err
		}
		defer metricsPublisher.Stop()
	}

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

	codeRefClient, err := coderefclient.NewClient(*s.codeRefService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer codeRefClient.Close()

	auditLogClient, err := auditlogclient.NewClient(*s.auditLogService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer auditLogClient.Close()

	autoOpsClient, err := autoopsclient.NewClient(*s.auditLogService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer auditLogClient.Close()

	tagClient, err := tagclient.NewClient(*s.tagService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer tagClient.Close()

	teamClient, err := teamclient.NewClient(*s.teamService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer teamClient.Close()

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
	defer notificationClient.Close()

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

	eventCounterClient, err := eventcounterclient.NewClient(*s.eventCounterService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer eventCounterClient.Close()

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

	service := api.NewGrpcGatewayService(
		featureClient,
		accountClient,
		pushClient,
		codeRefClient,
		auditLogClient,
		autoOpsClient,
		tagClient,
		teamClient,
		notificationClient,
		experimentClient,
		eventCounterClient,
		environmentClient,
		goalPublisher,
		evaluationPublisher,
		userPublisher,
		redisV3Cache,
		api.WithOldestEventTimestamp(*s.oldestEventTimestamp),
		api.WithFurthestEventTimestamp(*s.furthestEventTimestamp),
		api.WithMetrics(registerer),
		api.WithLogger(logger),
	)

	// We don't check the Redis health status because if the check fails,
	// the Kubernetes will restart the container and it might cause internal errors.
	// Use a dedicated context so we can stop the health checker goroutine cleanly during shutdown
	healthCheckCtx, healthCheckCancel := context.WithCancel(ctx)
	defer healthCheckCancel() // Ensure cleanup on all paths (including early returns)

	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(5*time.Second),
		health.WithCheck("metrics", metrics.Check),
	)
	go healthChecker.Run(healthCheckCtx)

	server := rpc.NewServer(service, *s.certPath, *s.keyPath,
		"api-gateway",
		rpc.WithPort(*s.port),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
		rpc.WithService(healthChecker),
		rpc.WithHandler("/health", healthChecker), // Liveness probe
		rpc.WithHandler("/ready", healthChecker),  // Readiness probe
	)
	go server.Run()

	// Set up gRPC Gateway for API service
	grpcGatewayAddr := fmt.Sprintf(":%d", *s.grpcGatewayPort)
	grpcAddr := fmt.Sprintf("localhost:%d", *s.port)

	// Create a HandlerRegistrar adapter function that matches gateway.HandlerRegistrar signature
	gatewayHandler := func(ctx context.Context,
		mux *runtime.ServeMux,
		opts []grpc.DialOption,
	) error {
		return gwproto.RegisterGatewayHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	}

	apiGateway, err := gateway.NewGateway(
		grpcGatewayAddr,
		gateway.WithLogger(logger.Named("api-grpc-gateway")),
		gateway.WithMetrics(registerer),
		gateway.WithCertPath(*s.certPath),
		gateway.WithKeyPath(*s.keyPath),
	)
	if err != nil {
		return fmt.Errorf("failed to create API gateway: %v", err)
	}

	if err := apiGateway.Start(ctx, gatewayHandler); err != nil {
		return fmt.Errorf("failed to start API gateway: %v", err)
	}

	restHealthChecker := health.NewRestChecker(
		api.Version, api.Service,
		health.WithTimeout(5*time.Second),
		health.WithCheck("metrics", metrics.Check),
	)
	go restHealthChecker.Run(healthCheckCtx)

	gatewayService := api.NewGatewayService(
		featureClient,
		accountClient,
		pushClient,
		goalPublisher,
		evaluationPublisher,
		userPublisher,
		metricsPublisher,
		redisV3Cache,
		api.WithOldestEventTimestamp(*s.oldestEventTimestamp),
		api.WithFurthestEventTimestamp(*s.furthestEventTimestamp),
		api.WithMetrics(registerer),
		api.WithLogger(logger),
	)

	httpServer := rest.NewServer(
		*s.certPath, *s.keyPath,
		rest.WithLogger(logger),
		rest.WithService(gatewayService),
		rest.WithService(restHealthChecker),
		rest.WithMetrics(registerer),
	)
	go httpServer.Run()

	defer func() {
		shutdownStartTime := time.Now()
		logger.Info("Starting graceful shutdown sequence")

		// Cancel the health checker goroutines to prevent connection errors during shutdown
		healthCheckCancel()
		// Mark as unhealthy so readiness probes fail
		// This ensures Kubernetes readiness probe fails on next check,
		// preventing new traffic from being routed to this pod.
		healthChecker.Stop()
		restHealthChecker.Stop()

		// Wait for K8s endpoint propagation
		// This prevents "context deadline exceeded" errors during high traffic.
		time.Sleep(propagationDelay)
		logger.Info("Starting HTTP/gRPC server shutdown")

		// CRITICAL: Shutdown order matters due to dependencies:
		// 1. apiGateway/httpServer make gRPC calls to the backend server
		// 2. We MUST drain them BEFORE stopping the backend
		// 3. Otherwise their handlers hang waiting for a dead backend
		// We run apiGateway and httpServer in parallel since they don't depend on each other
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			apiGateway.Stop(serverShutDownTimeout)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			httpServer.Stop(serverShutDownTimeout)
		}()

		// Wait for HTTP/REST traffic to fully drain
		wg.Wait()

		// Now it's safe to stop the gRPC server (no more HTTP→gRPC calls)
		server.Stop(grpcStopTimeout)

		// Close clients
		// These are fast cleanup operations that can run asynchronously.
		go goalPublisher.Stop()
		go evaluationPublisher.Stop()
		if userPublisher != nil {
			go userPublisher.Stop()
		}
		if metricsPublisher != nil {
			go metricsPublisher.Stop()
		}
		go featureClient.Close()
		go accountClient.Close()
		go pushClient.Close()
		go codeRefClient.Close()
		go auditLogClient.Close()
		go autoOpsClient.Close()
		go tagClient.Close()
		go teamClient.Close()
		go notificationClient.Close()
		go experimentClient.Close()
		go eventCounterClient.Close()
		go environmentClient.Close()
		go redisV3Client.Close()

		// Log total shutdown duration
		logger.Info("Graceful shutdown sequence completed",
			zap.Duration("total_elapsed", time.Since(shutdownStartTime)),
			zap.Duration("propagation_delay", propagationDelay),
			zap.Duration("server_shutdown_timeout", serverShutDownTimeout),
			zap.Duration("grpc_stop_timeout", grpcStopTimeout),
		)
	}()

	<-ctx.Done()
	return nil
}
