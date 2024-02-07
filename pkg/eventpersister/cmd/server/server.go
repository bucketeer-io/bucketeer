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

package server

import (
	"context"
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/eventpersister/persister"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
)

const (
	command = "server"
)

type server struct {
	*kingpin.CmdClause
	port                         *int
	project                      *string
	subscription                 *string
	topic                        *string
	maxMPS                       *int
	numWorkers                   *int
	flushSize                    *int
	flushInterval                *time.Duration
	flushTimeout                 *time.Duration
	certPath                     *string
	keyPath                      *string
	serviceTokenPath             *string
	pullerNumGoroutines          *int
	pullerMaxOutstandingMessages *int
	pullerMaxOutstandingBytes    *int
	mysqlUser                    *string
	mysqlPass                    *string
	mysqlHost                    *string
	mysqlPort                    *int
	mysqlDbName                  *string
	redisServerName              *string
	redisAddr                    *string
	redisPoolMaxIdle             *int
	redisPoolMaxActive           *int
}

func RegisterServerCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the server")
	server := &server{
		CmdClause:    cmd,
		port:         cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:      cmd.Flag("project", "Google Cloud project name.").String(),
		subscription: cmd.Flag("subscription", "Google PubSub subscription name.").String(),
		topic:        cmd.Flag("topic", "Google PubSub topic name.").String(),
		maxMPS:       cmd.Flag("max-mps", "Maximum messages should be handled in a second.").Default("1000").Int(),
		numWorkers:   cmd.Flag("num-workers", "Number of workers.").Default("2").Int(),
		flushSize: cmd.Flag(
			"flush-size",
			"Maximum number of messages to batch before writing to datastore.",
		).Default("50").Int(),
		flushInterval:    cmd.Flag("flush-interval", "Maximum interval between two flushes.").Default("5s").Duration(),
		flushTimeout:     cmd.Flag("flush-timeout", "Maximum time for a flush to finish.").Default("20s").Duration(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		pullerNumGoroutines: cmd.Flag(
			"puller-num-goroutines",
			"Number of goroutines will be spawned to pull messages.",
		).Int(),
		pullerMaxOutstandingMessages: cmd.Flag(
			"puller-max-outstanding-messages",
			"Maximum number of unprocessed messages.",
		).Int(),
		pullerMaxOutstandingBytes: cmd.Flag(
			"puller-max-outstanding-bytes",
			"Maximum size of unprocessed messages.",
		).Int(),
		mysqlUser:       cmd.Flag("mysql-user", "").String(),
		mysqlPass:       cmd.Flag("mysql-pass", "").String(),
		mysqlHost:       cmd.Flag("mysql-host", "").String(),
		mysqlPort:       cmd.Flag("mysql-port", "").Int(),
		mysqlDbName:     cmd.Flag("mysql-dbname", "").String(),
		redisServerName: cmd.Flag("redis-server-name", "Name of the redis.").Required().String(),
		redisAddr:       cmd.Flag("redis-addr", "Address of the redis.").Required().String(),
		redisPoolMaxIdle: cmd.Flag(
			"redis-pool-max-idle",
			"Maximum number of idle connections in the pool.",
		).Default("5").Int(),
		redisPoolMaxActive: cmd.Flag(
			"redis-pool-max-active",
			"Maximum number of connections allocated by the pool at a given time.",
		).Default("10").Int(),
	}
	r.RegisterCommand(server)
	return server
}

func (s *server) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()

	puller, err := s.createPuller(ctx, logger)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

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

	p := persister.NewPersister(
		puller,
		redisV3Cache,
		persister.WithMaxMPS(*s.maxMPS),
		persister.WithNumWorkers(*s.numWorkers),
		persister.WithFlushSize(*s.flushSize),
		persister.WithFlushInterval(*s.flushInterval),
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
	)
	go healthChecker.Run(ctx)

	server := rpc.NewServer(healthChecker, *s.certPath, *s.keyPath,
		"event-persister",
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

func (s *server) createPuller(ctx context.Context, logger *zap.Logger) (puller.Puller, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	client, err := pubsub.NewClient(ctx, *s.project, pubsub.WithLogger(logger))
	if err != nil {
		logger.Error("Failed to create PubSub client", zap.Error(err))
		return nil, err
	}
	return client.CreatePuller(*s.subscription, *s.topic,
		pubsub.WithNumGoroutines(*s.pullerNumGoroutines),
		pubsub.WithMaxOutstandingMessages(*s.pullerMaxOutstandingMessages),
		pubsub.WithMaxOutstandingBytes(*s.pullerMaxOutstandingBytes),
	)
}
