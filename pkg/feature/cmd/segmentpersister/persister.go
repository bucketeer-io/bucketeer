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

package segmentpersister

import (
	"context"
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	fsp "github.com/bucketeer-io/bucketeer/pkg/feature/segmentpersister"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

const command = "segment-persister"

type Persister interface {
	Run(context.Context, metrics.Metrics, *zap.Logger) error
}

type persister struct {
	*kingpin.CmdClause
	port                                      *int
	project                                   *string
	mysqlUser                                 *string
	mysqlPass                                 *string
	mysqlHost                                 *string
	mysqlPort                                 *int
	mysqlDBName                               *string
	domainEventTopic                          *string
	bulkSegmentUsersReceivedEventTopic        *string
	bulkSegmentUsersReceivedEventSubscription *string
	maxMPS                                    *int
	numWorkers                                *int
	flushSize                                 *int
	flushInterval                             *time.Duration
	pullerNumGoroutines                       *int
	pullerMaxOutstandingMessages              *int
	pullerMaxOutstandingBytes                 *int
	redisServerName                           *string
	redisAddr                                 *string
	redisPoolMaxIdle                          *int
	redisPoolMaxActive                        *int
	certPath                                  *string
	keyPath                                   *string
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start segment persister")
	persister := &persister{
		CmdClause:        cmd,
		port:             cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:          cmd.Flag("project", "Google Cloud project name.").String(),
		mysqlUser:        cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:        cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:        cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:        cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName:      cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		domainEventTopic: cmd.Flag("domain-event-topic", "PubSub topic to publish domain events.").Required().String(),
		bulkSegmentUsersReceivedEventTopic: cmd.Flag(
			"bulk-segment-users-received-event-topic",
			"PubSub topic to subscribe bulk segment users received events.",
		).Required().String(),
		bulkSegmentUsersReceivedEventSubscription: cmd.Flag(
			"bulk-segment-users-received-event-subscription",
			"PubSub subscription to subscribe bulk segment users received events.",
		).Required().String(),
		maxMPS:        cmd.Flag("max-mps", "Maximum messages should be handled in a second.").Default("100").Int(),
		numWorkers:    cmd.Flag("num-workers", "Number of workers.").Default("2").Int(),
		flushSize:     cmd.Flag("flush-size", "Maximum number of messages in one flush.").Default("2").Int(),
		flushInterval: cmd.Flag("flush-interval", "Maximum interval between two flushes.").Default("10s").Duration(),
		pullerNumGoroutines: cmd.Flag(
			"puller-num-goroutines",
			"Number of goroutines will be spawned to pull messages.",
		).Int(),
		pullerMaxOutstandingMessages: cmd.Flag(
			"puller-max-outstanding-messages",
			"Maximum number of unprocessed messages.",
		).Int(),
		pullerMaxOutstandingBytes: cmd.Flag("puller-max-outstanding-bytes", "Maximum size of unprocessed messages.").Int(),
		redisServerName:           cmd.Flag("redis-server-name", "Name of the redis.").Required().String(),
		redisAddr:                 cmd.Flag("redis-addr", "Address of the redis.").Required().String(),
		redisPoolMaxIdle: cmd.Flag(
			"redis-pool-max-idle",
			"Maximum number of idle connections in the pool.",
		).Default("5").Int(),
		redisPoolMaxActive: cmd.Flag(
			"redis-pool-max-active",
			"Maximum number of connections allocated by the pool at a given time.",
		).Default("10").Int(),
		certPath: cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:  cmd.Flag("key", "Path to TLS key.").Required().String(),
	}
	r.RegisterCommand(persister)
	return persister
}

func (p *persister) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()

	mysqlClient, err := p.createMySQLClient(ctx, registerer, logger)
	if err != nil {
		return err
	}
	defer mysqlClient.Close()

	redisV3Client, err := redisv3.NewClient(
		*p.redisAddr,
		redisv3.WithPoolSize(*p.redisPoolMaxActive),
		redisv3.WithMinIdleConns(*p.redisPoolMaxIdle),
		redisv3.WithServerName(*p.redisServerName),
		redisv3.WithMetrics(registerer),
		redisv3.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer redisV3Client.Close()
	redisV3Cache := cachev3.NewRedisCache(redisV3Client)

	pubsubClient, err := p.createPubsubClient(ctx, logger)
	if err != nil {
		return err
	}
	segmentUsersPuller, err := pubsubClient.CreatePuller(
		*p.bulkSegmentUsersReceivedEventSubscription,
		*p.bulkSegmentUsersReceivedEventTopic,
		pubsub.WithNumGoroutines(*p.pullerNumGoroutines),
		pubsub.WithMaxOutstandingMessages(*p.pullerMaxOutstandingMessages),
		pubsub.WithMaxOutstandingBytes(*p.pullerMaxOutstandingBytes),
	)
	if err != nil {
		return err
	}

	domainPublisher, err := pubsubClient.CreatePublisher(*p.domainEventTopic)
	if err != nil {
		return err
	}
	defer domainPublisher.Stop()

	persister := fsp.NewPersister(
		segmentUsersPuller,
		domainPublisher,
		mysqlClient,
		redisV3Cache,
		fsp.WithMaxMPS(*p.maxMPS),
		fsp.WithNumWorkers(*p.numWorkers),
		fsp.WithFlushSize(*p.flushSize),
		fsp.WithFlushInterval(*p.flushInterval),
		fsp.WithMetrics(registerer),
		fsp.WithLogger(logger),
	)
	defer persister.Stop()
	go persister.Run() // nolint:errcheck

	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
		health.WithCheck("segment-persister", persister.Check),
		health.WithCheck("redis", redisV3Client.Check),
	)
	go healthChecker.Run(ctx)

	server := rpc.NewServer(healthChecker, *p.certPath, *p.keyPath,
		"feature-persister",
		rpc.WithPort(*p.port),
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

func (p *persister) createMySQLClient(
	ctx context.Context,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (mysql.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return mysql.NewClient(
		ctx,
		*p.mysqlUser, *p.mysqlPass, *p.mysqlHost,
		*p.mysqlPort,
		*p.mysqlDBName,
		mysql.WithLogger(logger),
		mysql.WithMetrics(registerer),
	)
}
func (p *persister) createPubsubClient(ctx context.Context, logger *zap.Logger) (*pubsub.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	client, err := pubsub.NewClient(ctx, *p.project, pubsub.WithLogger(logger))
	if err != nil {
		return nil, err
	}
	return client, nil
}
