// Copyright 2024 The Bucketeer Authors.
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
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/gateway/api"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/pkg/rest"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
)

const command = "server"

type server struct {
	*kingpin.CmdClause
	port                   *int
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
	redisServerName        *string
	redisAddr              *string
	certPath               *string
	keyPath                *string
	serviceTokenPath       *string
	redisPoolMaxIdle       *int
	redisPoolMaxActive     *int
	oldestEventTimestamp   *time.Duration
	furthestEventTimestamp *time.Duration
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the gRPC server")
	server := &server{
		CmdClause: cmd,
		port:      cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:   cmd.Flag("project", "GCP Project id to use for PubSub.").Required().String(),
		goalTopic: cmd.Flag("goal-topic", "Topic to use for publishing GoalEvent.").Required().String(),
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
		).Default("24h").Duration(),
		furthestEventTimestamp: cmd.Flag(
			"furthest-event-timestamp",
			"The duration of furthest event timestamp from processing time to allow.",
		).Default("24h").Duration(),
	}
	r.RegisterCommand(server)
	return server
}

func (s *server) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()

	pubsubCtx, pubsubCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pubsubCancel()
	pubsubClient, err := pubsub.NewClient(
		pubsubCtx,
		*s.project,
		pubsub.WithMetrics(registerer),
		pubsub.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	publishOptions := []pubsub.PublishOption{pubsub.WithPublishTimeout(*s.publishTimeout)}
	if *s.publishNumGoroutines > 0 {
		publishOptions = append(publishOptions, pubsub.WithPublishNumGoroutines(*s.publishNumGoroutines))
	}

	var goalTopicProject string
	if *s.goalTopicProject == "" {
		goalTopicProject = *s.project
	} else {
		goalTopicProject = *s.goalTopicProject
	}
	goalPublisher, err := pubsubClient.CreatePublisherInProject(*s.goalTopic, goalTopicProject, publishOptions...)
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
		publishOptions...,
	)
	if err != nil {
		return nil
	}
	defer evaluationPublisher.Stop()

	// FIXME: This condition won't be necessary once user feature is fully released.
	var userPublisher publisher.Publisher
	if *s.userTopic != "" {
		userPublisher, err = pubsubClient.CreatePublisherInProject(*s.userTopic, *s.project, publishOptions...)
		if err != nil {
			return err
		}
		defer userPublisher.Stop()
	}

	// FIXME: This condition won't be necessary once user feature is fully released.
	var metricsPublisher publisher.Publisher
	if *s.metricsTopic != "" {
		metricsPublisher, err = pubsubClient.CreatePublisherInProject(*s.metricsTopic, *s.project, publishOptions...)
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

	inMemoryCache := cachev3.NewInMemoryCache(cachev3.WithEvictionInterval(cachev3.EnvironmentAPIKeyEvictionInterval))

	service := api.NewGrpcGatewayService(
		featureClient,
		accountClient,
		goalPublisher,
		evaluationPublisher,
		userPublisher,
		redisV3Cache,
		inMemoryCache,
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

	restHealthChecker := health.NewRestChecker(
		api.Version, api.Service,
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
	)
	go restHealthChecker.Run(ctx)

	gatewayService := api.NewGatewayService(
		featureClient,
		accountClient,
		goalPublisher,
		evaluationPublisher,
		userPublisher,
		metricsPublisher,
		redisV3Cache,
		inMemoryCache,
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
