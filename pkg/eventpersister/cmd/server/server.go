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

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/eventpersister/datastore"
	"github.com/bucketeer-io/bucketeer/pkg/eventpersister/persister"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/storage/kafka"
	bigtable "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigtable"
)

const (
	command = "server"
)

type server struct {
	*kingpin.CmdClause
	port                         *int
	project                      *string
	bigtableInstance             *string
	subscription                 *string
	topic                        *string
	maxMPS                       *int
	numWorkers                   *int
	kafkaURL                     *string
	kafkaTopicPrefix             *string
	kafkaTopicDataType           *string
	kafkaUsername                *string
	kafkaPassword                *string
	numWriters                   *int
	flushSize                    *int
	flushInterval                *time.Duration
	flushTimeout                 *time.Duration
	featureService               *string
	certPath                     *string
	keyPath                      *string
	serviceTokenPath             *string
	pullerNumGoroutines          *int
	pullerMaxOutstandingMessages *int
	pullerMaxOutstandingBytes    *int
	alloyDBRegion                *string
	alloyDBClusterID             *string
	alloyDBInstanceID            *string
	alloyDBUser                  *string
	alloyDBPass                  *string
	alloyDBName                  *string
}

func RegisterServerCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the server")
	server := &server{
		CmdClause:          cmd,
		port:               cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:            cmd.Flag("project", "Google Cloud project name.").String(),
		bigtableInstance:   cmd.Flag("bigtable-instance", "Instance name to use Bigtable.").Required().String(),
		subscription:       cmd.Flag("subscription", "Google PubSub subscription name.").String(),
		topic:              cmd.Flag("topic", "Google PubSub topic name.").String(),
		maxMPS:             cmd.Flag("max-mps", "Maximum messages should be handled in a second.").Default("1000").Int(),
		numWorkers:         cmd.Flag("num-workers", "Number of workers.").Default("2").Int(),
		kafkaURL:           cmd.Flag("kafka-url", "Kafka URL.").String(),
		kafkaTopicPrefix:   cmd.Flag("kafka-topic-prefix", "Kafka topic dataset section prefix.").String(),
		kafkaTopicDataType: cmd.Flag("kafka-topic-data-type", "Kafka topic data type.").String(),
		kafkaUsername:      cmd.Flag("kafka-username", "Kafka username.").String(),
		kafkaPassword:      cmd.Flag("kafka-password", "Kafka password.").String(),
		numWriters:         cmd.Flag("num-writers", "Number of writers.").Default("2").Int(),
		flushSize: cmd.Flag(
			"flush-size",
			"Maximum number of messages to batch before writing to datastore.",
		).Default("50").Int(),
		flushInterval:    cmd.Flag("flush-interval", "Maximum interval between two flushes.").Default("5s").Duration(),
		flushTimeout:     cmd.Flag("flush-timeout", "Maximum time for a flush to finish.").Default("20s").Duration(),
		featureService:   cmd.Flag("feature-service", "bucketeer-feature-service address.").Default("feature:9090").String(),
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
		alloyDBRegion:     cmd.Flag("alloydb-region", "").Required().String(),
		alloyDBClusterID:  cmd.Flag("alloydb-cluster-id", "").Required().String(),
		alloyDBInstanceID: cmd.Flag("alloydb-instance-id", "").Required().String(),
		alloyDBUser:       cmd.Flag("alloydb-user", "").Required().String(),
		alloyDBPass:       cmd.Flag("alloydb-pass", "").Required().String(),
		alloyDBName:       cmd.Flag("alloydb-name", "").Required().String(),
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

	datastore, err := s.createWriters(ctx, registerer, logger)
	if err != nil {
		return err
	}
	defer datastore.Close()

	btClient, err := s.createBigtableClient(ctx, registerer, logger)
	if err != nil {
		return err
	}
	defer btClient.Close()

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

	if err != nil {
		return err
	}

	// TODO: Return error after postgreSQL is stable.
	// postgresClient, _ := s.createPostgresqlClient(ctx, logger)
	// if err != nil {
	// 	return err
	// }
	// if postgresClient != nil {
	// 	defer postgresClient.Close()
	// }

	p := persister.NewPersister(
		featureClient,
		puller,
		datastore,
		btClient,
		nil, // Disable PostgreSQL temporarily due to instability issues on the Google side.
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
	)
	go healthChecker.Run(ctx)

	server := rpc.NewServer(healthChecker, *s.certPath, *s.keyPath,
		rpc.WithPort(*s.port),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
		rpc.WithHandler("/health", healthChecker),
	)
	defer server.Stop(10 * time.Second)
	go server.Run()

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

func (s *server) createWriters(
	ctx context.Context,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (datastore.Writer, error) {
	writers := make([]datastore.Writer, 0, *s.numWriters)
	for i := 0; i < *s.numWriters; i++ {
		writer, err := s.createKafkaWriter(ctx, registerer, logger)
		if err != nil {
			return nil, err
		}
		writers = append(writers, writer)
	}
	if len(writers) == 1 {
		logger.Info("Created a single writer", zap.Int("numWriters", *s.numWriters))
		return writers[0], nil
	}
	logger.Info("Created a writer pool", zap.Int("numWriters", *s.numWriters), zap.Int("poolSize", len(writers)))
	return datastore.NewWriterPool(writers), nil
}

func (s *server) createKafkaWriter(
	ctx context.Context,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (datastore.Writer, error) {
	logger.Debug("createKafkaWriter")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	kafkaProducer, err := kafka.NewProducer(
		ctx,
		*s.project,
		*s.kafkaURL,
		*s.kafkaUsername,
		*s.kafkaPassword)
	if err != nil {
		logger.Error("Failed to create Kafka producer", zap.Error(err))
		return nil, err
	}
	writer, err := datastore.NewKafkaWriter(kafkaProducer,
		*s.kafkaTopicPrefix,
		*s.kafkaTopicDataType,
		datastore.WithMetrics(registerer),
		datastore.WithLogger(logger),
	)
	if err != nil {
		logger.Error("Failed to create Kafka writer", zap.Error(err))
		return nil, err
	}
	return writer, nil
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

// func (s *server) createPostgresqlClient(
// 	ctx context.Context,
// 	logger *zap.Logger,
// ) (postgres.Client, error) {
// 	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()
// 	return postgres.NewClient(
// 		ctx,
// 		*s.project, *s.alloyDBRegion, *s.alloyDBClusterID, *s.alloyDBInstanceID,
// 		*s.alloyDBUser, *s.alloyDBPass, *s.alloyDBName,
// 		postgres.WithLogger(logger),
// 	)
// }
