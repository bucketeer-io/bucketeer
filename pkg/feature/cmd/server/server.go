// Copyright 2022 The Bucketeer Authors.
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

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	"github.com/bucketeer-io/bucketeer/pkg/feature/api"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	bigtable "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigtable"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/token"
)

const command = "server"

type server struct {
	*kingpin.CmdClause
	port                               *int
	project                            *string
	mysqlUser                          *string
	mysqlPass                          *string
	mysqlHost                          *string
	mysqlPort                          *int
	mysqlDBName                        *string
	bigtableInstance                   *string
	accountService                     *string
	experimentService                  *string
	redisServerName                    *string
	redisAddr                          *string
	redisPoolMaxIdle                   *int
	redisPoolMaxActive                 *int
	bulkSegmentUsersReceivedEventTopic *string
	domainEventTopic                   *string
	certPath                           *string
	keyPath                            *string
	serviceTokenPath                   *string

	oauthKeyPath  *string
	oauthClientID *string
	oauthIssuer   *string
}

func RegisterServerCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the gRPC server")
	server := &server{
		CmdClause:        cmd,
		port:             cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:          cmd.Flag("project", "Google Cloud project name.").Required().String(),
		mysqlUser:        cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:        cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:        cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:        cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName:      cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		bigtableInstance: cmd.Flag("bigtable-instance", "Instance name to use Bigtable.").Required().String(),
		accountService: cmd.Flag(
			"account-service",
			"bucketeer-account-service address.",
		).Default("account:9090").String(),
		experimentService: cmd.Flag(
			"experiment-service",
			"bucketeer-experiment-service address.",
		).Default("experiment:9090").String(),
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
		bulkSegmentUsersReceivedEventTopic: cmd.Flag(
			"bulk-segment-users-received-event-topic",
			"PubSub topic to publish bulk segment users received events.",
		).Required().String(),
		domainEventTopic: cmd.Flag("domain-event-topic", "PubSub topic to publish domain events.").Required().String(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		oauthKeyPath:     cmd.Flag("oauth-key", "Path to public key used to verify oauth token.").Required().String(),
		oauthClientID:    cmd.Flag("oauth-client-id", "The oauth clientID registered at dex.").Required().String(),
		oauthIssuer:      cmd.Flag("oauth-issuer", "The url of dex issuer.").Required().String(),
	}
	r.RegisterCommand(server)
	return server
}

func (s *server) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()

	mysqlClient, err := s.createMySQLClient(ctx, registerer, logger)
	if err != nil {
		return err
	}
	defer mysqlClient.Close()

	btClient, err := s.createBigtableClient(ctx, registerer, logger)
	if err != nil {
		return err
	}
	defer btClient.Close()

	creds, err := client.NewPerRPCCredentials(*s.serviceTokenPath)
	if err != nil {
		return err
	}

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

	pubsubClient, err := s.createPubsubClient(ctx, registerer, logger)
	if err != nil {
		return err
	}
	segmentUsersPublisher, err := pubsubClient.CreatePublisher(*s.bulkSegmentUsersReceivedEventTopic)
	if err != nil {
		return err
	}
	defer segmentUsersPublisher.Stop()
	domainPublisher, err := pubsubClient.CreatePublisher(*s.domainEventTopic)
	if err != nil {
		return err
	}
	defer domainPublisher.Stop()

	service := api.NewFeatureService(
		mysqlClient,
		btClient,
		accountClient,
		experimentClient,
		redisV3Cache,
		segmentUsersPublisher,
		domainPublisher,
		api.WithLogger(logger),
	)

	verifier, err := token.NewVerifier(*s.oauthKeyPath, *s.oauthIssuer, *s.oauthClientID)
	if err != nil {
		return err
	}

	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
	)
	go healthChecker.Run(ctx)

	server := rpc.NewServer(service, *s.certPath, *s.keyPath,
		rpc.WithPort(*s.port),
		rpc.WithVerifier(verifier),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
		rpc.WithService(healthChecker),
		rpc.WithHandler("/health", healthChecker),
	)
	defer server.Stop(10 * time.Second)
	go server.Run()

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

func (s *server) createPubsubClient(
	ctx context.Context,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (*pubsub.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	client, err := pubsub.NewClient(
		ctx,
		*s.project,
		pubsub.WithMetrics(registerer),
		pubsub.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *server) createBigtableClient(
	ctx context.Context,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (bigtable.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return bigtable.NewBigtableClient(ctx, *s.project, *s.bigtableInstance,
		bigtable.WithMetrics(registerer),
		bigtable.WithLogger(logger),
	)
}
