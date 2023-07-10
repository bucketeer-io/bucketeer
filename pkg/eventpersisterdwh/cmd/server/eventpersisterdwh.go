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
	"errors"
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/eventpersisterdwh/persister"
	ec "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	ft "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
)

const (
	command         = "server"
	evalEvtSvcName  = "event-persister-evaluation-events-dwh"
	evalGoalSvcName = "event-persister-goal-events-dwh"
)

var errUnknownSvcName = errors.New("persister: unknown service name")

type server struct {
	*kingpin.CmdClause
	serviceName *string
	// option
	maxMPS        *int
	numWorkers    *int
	flushSize     *int
	flushInterval *time.Duration
	flushTimeout  *time.Duration
	timezone      *string
	// pubsub
	project                      *string
	subscription                 *string
	topic                        *string
	pullerNumGoroutines          *int
	pullerMaxOutstandingMessages *int
	pullerMaxOutstandingBytes    *int
	// rpc
	port              *int
	serviceTokenPath  *string
	certPath          *string
	keyPath           *string
	experimentService *string
	featureService    *string
	// bigquery
	bigQueryDataSet   *string
	bigQueryBatchSize *int
}

func RegisterServerCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the server")
	server := &server{
		CmdClause:   cmd,
		serviceName: cmd.Flag("service-name", "Service name.").Required().String(),
		port:        cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:     cmd.Flag("project", "Google Cloud project name.").Required().String(),
		maxMPS:      cmd.Flag("max-mps", "Maximum messages should be handled in a second.").Default("1000").Int(),
		numWorkers:  cmd.Flag("num-workers", "Number of workers.").Default("2").Int(),
		flushSize: cmd.Flag(
			"flush-size",
			"Maximum number of messages to batch before writing to datastore.",
		).Default("50").Int(),
		flushInterval: cmd.Flag("flush-interval", "Maximum interval between two flushes.").Default("5s").Duration(),
		flushTimeout:  cmd.Flag("flush-timeout", "Maximum time for a flush to finish.").Default("20s").Duration(),
		timezone:      cmd.Flag("timezone", "Time zone").Required().String(),
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
		featureService: cmd.Flag(
			"feature-service",
			"bucketeer-feature-service address.",
		).Default("featureService:9090").String(),
		bigQueryDataSet:   cmd.Flag("bigquery-data-set", "BigQuery DataSet Name").Required().String(),
		bigQueryBatchSize: cmd.Flag("bigquery-batch-size", "BigQuery Size of rows to be sent at once").Default("10").Int(),
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
	featureClient, err := ft.NewClient(*s.featureService, *s.certPath,
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
	location, err := locale.GetLocation(*s.timezone)
	if err != nil {
		return err
	}
	writer, err := s.newBigQueryWriter(
		ctx,
		registerer,
		logger,
		experimentClient,
		featureClient,
		location,
	)
	if err != nil {
		return err
	}
	p := persister.NewPersisterDWH(
		puller,
		registerer,
		writer,
		persister.WithMaxMPS(*s.maxMPS),
		persister.WithNumWorkers(*s.numWorkers),
		persister.WithFlushSize(*s.flushSize),
		persister.WithFlushInterval(*s.flushInterval),
		persister.WithFlushTimeout(*s.flushTimeout),
		persister.WithMetrics(registerer),
		persister.WithLogger(logger),
		persister.WithBatchSize(*s.bigQueryBatchSize),
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
		"event-persister-dwh",
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

func (s *server) newBigQueryWriter(
	ctx context.Context,
	r metrics.Registerer,
	logger *zap.Logger,
	exClient ec.Client,
	ftClient ft.Client,
	location *time.Location,
) (persister.Writer, error) {
	var writer persister.Writer
	var err error
	switch *s.serviceName {
	case evalEvtSvcName:
		writer, err = persister.NewEvalEventWriter(
			ctx,
			r,
			logger,
			exClient,
			*s.project,
			*s.bigQueryDataSet,
			*s.bigQueryBatchSize,
			location,
		)
	case evalGoalSvcName:
		writer, err = persister.NewGoalEventWriter(
			ctx,
			r,
			logger,
			exClient,
			ftClient,
			*s.project,
			*s.bigQueryDataSet,
			*s.bigQueryBatchSize,
			location,
		)
	default:
		return nil, errUnknownSvcName
	}
	return writer, err
}
