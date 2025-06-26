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
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/api/api"
	auditlogclient "github.com/bucketeer-io/bucketeer/pkg/auditlog/client"
	autoopsclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	coderefclient "github.com/bucketeer-io/bucketeer/pkg/coderef/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/factory"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	pushclient "github.com/bucketeer-io/bucketeer/pkg/push/client"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/pkg/rest"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/gateway"
	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
)

const command = "server"

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
	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
	)
	go healthChecker.Run(ctx)

	server := rpc.NewServer(service, *s.certPath, *s.keyPath,
		"api-gateway",
		rpc.WithPort(*s.port),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
		rpc.WithService(healthChecker),
		rpc.WithHandler("/health", healthChecker),
	)
	defer server.Stop(10 * time.Second)
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
	defer apiGateway.Stop(10 * time.Second)

	restHealthChecker := health.NewRestChecker(
		api.Version, api.Service,
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
	)
	go restHealthChecker.Run(ctx)

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
	defer httpServer.Stop(10 * time.Second)
	go httpServer.Run()

	// Ensure to stop the health check before stopping the application
	// so the Kubernetes Readiness can detect it faster and remove the pod
	// from the service load balancer.
	defer healthChecker.Stop()
	defer restHealthChecker.Stop()

	<-ctx.Done()
	return nil
}
