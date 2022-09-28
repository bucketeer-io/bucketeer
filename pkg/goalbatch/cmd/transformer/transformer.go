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

package transformer

import (
	"context"
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	tf "github.com/bucketeer-io/bucketeer/pkg/goalbatch/transformer"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	userclient "github.com/bucketeer-io/bucketeer/pkg/user/client"
)

const command = "transformer"

type transformer struct {
	*kingpin.CmdClause
	port                         *int
	metricsTopic                 *string
	project                      *string
	userService                  *string
	goalBatchTopic               *string
	goalBatchSubscription        *string
	goalTopic                    *string
	maxMPS                       *int
	numWorkers                   *int
	certPath                     *string
	keyPath                      *string
	serviceTokenPath             *string
	pullerNumGoroutines          *int
	pullerMaxOutstandingMessages *int
	pullerMaxOutstandingBytes    *int
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start transformer server")
	s := &transformer{
		CmdClause:      cmd,
		port:           cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		metricsTopic:   cmd.Flag("metrics-topic", "Topic to use for publishing MetricsEvent.").String(),
		project:        cmd.Flag("project", "Google Cloud project name.").String(),
		userService:    cmd.Flag("user-service", "bucketeer-user-service address.").Default("user:9090").String(),
		goalBatchTopic: cmd.Flag("goal-batch-topic", "Google PubSub topic name of incoming goal batch events.").String(),
		goalBatchSubscription: cmd.Flag(
			"goal-batch-subscription",
			"Google PubSub subscription name of incoming goal batch event.",
		).String(),
		goalTopic:        cmd.Flag("goal-topic", "Google PubSub topic name of outgoing goal events.").String(),
		maxMPS:           cmd.Flag("max-mps", "Maximum messages should be handled in a second.").Default("5000").Int(),
		numWorkers:       cmd.Flag("num-workers", "Number of workers.").Default("1").Int(),
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
		pullerMaxOutstandingBytes: cmd.Flag("puller-max-outstanding-bytes", "Maximum size of unprocessed messages.").Int(),
	}
	r.RegisterCommand(s)
	return s
}

func (t *transformer) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()

	goalPublisher, goalBatchPuller, err := t.createPublisherPuller(ctx, registerer, logger)
	if err != nil {
		return err
	}
	defer goalPublisher.Stop()

	creds, err := client.NewPerRPCCredentials(*t.serviceTokenPath)
	if err != nil {
		return err
	}

	userClient, err := userclient.NewClient(*t.userService, *t.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer userClient.Close()

	goalBatchTransformer := tf.NewTransformer(
		userClient,
		goalBatchPuller,
		goalPublisher,
		tf.WithMaxMPS(*t.maxMPS),
		tf.WithNumWorkers(*t.numWorkers),
		tf.WithMetrics(registerer),
		tf.WithLogger(logger),
	)
	defer goalBatchTransformer.Stop()
	go goalBatchTransformer.Run() // nolint:errcheck

	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
		health.WithCheck("transformer", goalBatchTransformer.Check),
	)
	go healthChecker.Run(ctx)

	server := rpc.NewServer(healthChecker, *t.certPath, *t.keyPath,
		rpc.WithPort(*t.port),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
		rpc.WithHandler("/health", healthChecker),
	)
	defer server.Stop(10 * time.Second)
	go server.Run()

	<-ctx.Done()
	return nil
}

func (t *transformer) createPublisherPuller(
	ctx context.Context,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (publisher.Publisher, puller.Puller, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	client, err := pubsub.NewClient(
		ctx,
		*t.project,
		pubsub.WithMetrics(registerer),
		pubsub.WithLogger(logger),
	)
	if err != nil {
		return nil, nil, err
	}
	goalBatchPuller, err := client.CreatePuller(*t.goalBatchSubscription, *t.goalBatchTopic,
		pubsub.WithNumGoroutines(*t.pullerNumGoroutines),
		pubsub.WithMaxOutstandingMessages(*t.pullerMaxOutstandingMessages),
		pubsub.WithMaxOutstandingBytes(*t.pullerMaxOutstandingBytes),
	)
	if err != nil {
		return nil, nil, err
	}
	goalPublisher, err := client.CreatePublisher(*t.goalTopic)
	if err != nil {
		return nil, nil, err
	}
	return goalPublisher, goalBatchPuller, nil
}
