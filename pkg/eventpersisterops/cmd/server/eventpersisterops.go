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
	"errors"
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	aoclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/eventpersisterops/persister"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

const (
	command         = "server"
	evalEvtSvcName  = "event-persister-evaluation-events-ops"
	evalGoalSvcName = "event-persister-goal-events-ops"
)

var errUnknownSvcName = errors.New("persister: unknown service name")

type server struct {
	*kingpin.CmdClause
	serviceName *string
	// egress services
	featureService *string
	autoOpsService *string
	// rpc
	port             *int
	serviceTokenPath *string
	certPath         *string
	keyPath          *string
	// pubsub
	project                      *string
	subscription                 *string
	topic                        *string
	pullerNumGoroutines          *int
	pullerMaxOutstandingMessages *int
	pullerMaxOutstandingBytes    *int
	// redis
	redisServerName    *string
	redisAddr          *string
	redisPoolMaxIdle   *int
	redisPoolMaxActive *int
	// mysql
	mysqlUser   *string
	mysqlPass   *string
	mysqlHost   *string
	mysqlPort   *int
	mysqlDbName *string
	// batch options
	maxMPS        *int
	numWorkers    *int
	flushSize     *int
	checkInterval *time.Duration
	flushInterval *time.Duration
	flushTimeout  *time.Duration
}

func RegisterServerCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the server")
	server := &server{
		CmdClause:        cmd,
		serviceName:      cmd.Flag("service-name", "Service name.").Required().String(),
		project:          cmd.Flag("project", "Google Cloud project name.").Required().String(),
		port:             cmd.Flag("port", "Port to bind to.").Required().Int(),
		featureService:   cmd.Flag("feature-service", "bucketeer-feature-service address.").Required().String(),
		autoOpsService:   cmd.Flag("auto-ops-service", "bucketeer-auto-ops-service address.").Required().String(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		subscription:     cmd.Flag("subscription", "Google PubSub subscription name.").Required().String(),
		topic:            cmd.Flag("topic", "Google PubSub topic name.").Required().String(),
		pullerNumGoroutines: cmd.Flag(
			"puller-num-goroutines",
			"Number of goroutines will be spawned to pull messages.",
		).Required().Int(),
		pullerMaxOutstandingMessages: cmd.Flag(
			"puller-max-outstanding-messages",
			"Maximum number of unprocessed messages.",
		).Required().Int(),
		pullerMaxOutstandingBytes: cmd.Flag(
			"puller-max-outstanding-bytes",
			"Maximum size of unprocessed messages.",
		).Required().Int(),
		redisServerName: cmd.Flag("redis-server-name", "Name of the redis.").Required().String(),
		redisAddr:       cmd.Flag("redis-addr", "Address of the redis.").Required().String(),
		redisPoolMaxIdle: cmd.Flag(
			"redis-pool-max-idle",
			"Maximum number of idle connections in the pool.",
		).Required().Int(),
		redisPoolMaxActive: cmd.Flag(
			"redis-pool-max-active",
			"Maximum number of connections allocated by the pool at a given time.",
		).Required().Int(),
		mysqlUser:   cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:   cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:   cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:   cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDbName: cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		maxMPS:      cmd.Flag("max-mps", "Maximum messages should be handled in a second.").Required().Int(),
		numWorkers:  cmd.Flag("num-workers", "Number of workers.").Required().Int(),
		flushSize: cmd.Flag(
			"flush-size",
			"Maximum number of messages to batch before writing to datastore.",
		).Required().Int(),
		checkInterval: cmd.Flag(
			"check-interval",
			"Interval to check if there are auto ops rules to be handled.",
		).Required().Duration(),
		flushInterval: cmd.Flag("flush-interval", "Maximum interval between two flushes.").Required().Duration(),
		flushTimeout:  cmd.Flag("flush-timeout", "Maximum time for a flush to finish.").Required().Duration(),
	}
	r.RegisterCommand(server)
	return server
}

func (s *server) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()

	pubsubClient, err := s.createPubsubClient(ctx, logger)
	if err != nil {
		return err
	}

	// mysqlClient
	mysqlClient, err := s.createMySQLClient(ctx, registerer, logger)
	if err != nil {
		return err
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

	autoOpsClient, err := aoclient.NewClient(*s.autoOpsService, *s.certPath,
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

	updater, err := s.newUserCountUpdater(
		ctx,
		featureClient,
		autoOpsClient,
		redisV3Cache,
		registerer,
		logger,
	)
	if err != nil {
		return err
	}

	p := persister.NewPersister(
		updater,
		mysqlClient,
		pubsubClient,
		*s.subscription,
		*s.topic,
		*s.pullerNumGoroutines,
		*s.pullerMaxOutstandingMessages,
		*s.pullerMaxOutstandingBytes,
		persister.WithMaxMPS(*s.maxMPS),
		persister.WithNumWorkers(*s.numWorkers),
		persister.WithFlushSize(*s.flushSize),
		persister.WithCheckInterval(*s.checkInterval),
		persister.WithFlushInterval(*s.flushInterval),
		persister.WithFlushTimeout(*s.flushTimeout),
		persister.WithMetrics(registerer),
		persister.WithLogger(logger),
	)
	defer p.Stop()
	go p.Run() // nolint:errcheck

	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
		health.WithCheck("persister", p.Check),
		health.WithCheck("redis", redisV3Client.Check),
		health.WithCheck("subscription", func(ctx context.Context) health.Status {
			exists, err := pubsubClient.SubscriptionExists(*s.subscription)
			if err != nil {
				logger.Error("Failed to check subscription",
					zap.Error(err),
					zap.String("subscription", *s.subscription),
				)
				return health.Unhealthy
			}
			if !exists {
				logger.Info("Subscription does not exist",
					zap.String("subscription", *s.subscription),
				)
				return health.Unhealthy
			}
			return health.Healthy
		}),
	)
	go healthChecker.Run(ctx)

	server := rpc.NewServer(healthChecker, *s.certPath, *s.keyPath,
		"event-persister-ops",
		rpc.WithPort(*s.port),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
		rpc.WithHandler("/health", healthChecker),
	)
	defer server.Stop(10 * time.Second)
	go server.Run()

	// Ensure to stop the health check before stopping the application
	// so the Kubernetes Readiness can detect it faster and remove the pod
	// from the service load balancer.
	defer healthChecker.Stop()

	<-ctx.Done()
	return nil
}

func (s *server) createPubsubClient(ctx context.Context, logger *zap.Logger) (*pubsub.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	pubsubClient, err := pubsub.NewClient(ctx, *s.project, pubsub.WithLogger(logger))
	if err != nil {
		logger.Error("Failed to create PubSub client", zap.Error(err))
		return nil, err
	}
	return pubsubClient, nil
}

func (s *server) newUserCountUpdater(
	ctx context.Context,
	featureClient featureclient.Client,
	autoOpsClient aoclient.Client,
	redis cache.MultiGetDeleteCountCache,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (persister.Updater, error) {
	var updater persister.Updater
	var err error
	switch *s.serviceName {
	case evalEvtSvcName:
		updater = persister.NewEvalUserCountUpdater(
			ctx,
			featureClient,
			autoOpsClient,
			cachev3.NewEventCountCache(redis),
			cachev3.NewAutoOpsRulesCache(redis),
			logger,
		)
	case evalGoalSvcName:
		updater = persister.NewGoalUserCountUpdater(
			ctx,
			featureClient,
			autoOpsClient,
			cachev3.NewEventCountCache(redis),
			cachev3.NewAutoOpsRulesCache(redis),
			logger,
		)
	default:
		return nil, errUnknownSvcName
	}
	return updater, err
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
		*s.mysqlDbName,
		mysql.WithLogger(logger),
		mysql.WithMetrics(registerer),
	)
}
