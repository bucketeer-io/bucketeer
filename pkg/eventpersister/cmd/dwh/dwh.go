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

package dwh

import (
	"context"
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/eventpersister/datastore"
	"github.com/bucketeer-io/bucketeer/pkg/eventpersister/persister"
	ec "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigtable"
)

const (
	command = "server"
)

type server struct {
	*kingpin.CmdClause
	port             *int
	project          *string
	bigtableInstance *string
	// option
	maxMPS        *int
	numWorkers    *int
	flushSize     *int
	flushInterval *time.Duration
	flushTimeout  *time.Duration
	// pubsub
	subscription                 *string
	topic                        *string
	pullerNumGoroutines          *int
	pullerMaxOutstandingMessages *int
	pullerMaxOutstandingBytes    *int
	// rpc
	serviceTokenPath  *string
	certPath          *string
	keyPath           *string
	experimentService *string
	// bigquery
	bigQueryDataSet *string
}

func RegisterServerCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the server")
	server := &server{
		CmdClause:        cmd,
		port:             cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:          cmd.Flag("project", "Google Cloud project name.").String(),
		bigtableInstance: cmd.Flag("bigtable-instance", "Instance name to use Bigtable.").Required().String(),
		maxMPS:           cmd.Flag("max-mps", "Maximum messages should be handled in a second.").Default("1000").Int(),
		numWorkers:       cmd.Flag("num-workers", "Number of workers.").Default("2").Int(),
		flushSize: cmd.Flag(
			"flush-size",
			"Maximum number of messages to batch before writing to datastore.",
		).Default("50").Int(),
		flushInterval: cmd.Flag("flush-interval", "Maximum interval between two flushes.").Default("5s").Duration(),
		flushTimeout:  cmd.Flag("flush-timeout", "Maximum time for a flush to finish.").Default("20s").Duration(),
		subscription:  cmd.Flag("subscription", "Google PubSub subscription name.").String(),
		topic:         cmd.Flag("topic", "Google PubSub topic name.").String(),
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
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		experimentService: cmd.Flag(
			"experiment-service",
			"bucketeer-experiment-service address.",
		).Default("experiment:9090").String(),
		bigQueryDataSet: cmd.Flag("bigquery-data-set", "BigQuery DataSet Name").String(),
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

	btClient, err := s.createBigtableClient(ctx, registerer, logger)
	if err != nil {
		return err
	}
	defer btClient.Close()

	creds, err := client.NewPerRPCCredentials(*s.serviceTokenPath)
	if err != nil {
		return err
	}

	experimentClient, err := ec.NewClient(*s.experimentService, *s.certPath,
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
	evalEventWriter, err := datastore.NewEvalEventWriter(ctx, *s.project, *s.bigQueryDataSet)
	if err != nil {
		return err
	}
	goalEventWriter, err := datastore.NewGoalEventWriter(ctx, *s.project, *s.bigQueryDataSet)
	if err != nil {
		return err
	}
	p := persister.NewPersisterDwh(
		experimentClient,
		puller,
		evalEventWriter,
		goalEventWriter,
		btClient,
		persister.WithMaxMPS(*s.maxMPS),
		persister.WithNumWorkers(*s.numWorkers),
		persister.WithFlushSize(*s.flushSize),
		persister.WithFlushInterval(*s.flushInterval),
		persister.WithFlushTimeout(*s.flushTimeout),
		persister.WithMetrics(registerer),
		persister.WithLogger(logger),
	)
	if err != nil {
		return err
	}
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
