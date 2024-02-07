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

package persister

import (
	"context"
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	pst "github.com/bucketeer-io/bucketeer/pkg/metricsevent/persister"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
)

const command = "persister"

type Persister interface {
	Run(context.Context, metrics.Metrics, *zap.Logger) error
}

type persister struct {
	*kingpin.CmdClause
	port                         *int
	project                      *string
	subscription                 *string
	maxMPS                       *int
	numWorkers                   *int
	topic                        *string
	flushSize                    *int
	flushInterval                *time.Duration
	certPath                     *string
	keyPath                      *string
	pullerNumGoroutines          *int
	pullerMaxOutstandingMessages *int
	pullerMaxOutstandingBytes    *int
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start metricsevent persister")
	persister := &persister{
		CmdClause:     cmd,
		port:          cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:       cmd.Flag("project", "Google Cloud project name.").String(),
		subscription:  cmd.Flag("subscription", "Google PubSub subscription name.").Required().String(),
		topic:         cmd.Flag("topic", "Google PubSub topic name.").Required().String(),
		maxMPS:        cmd.Flag("max-mps", "Maximum messages should be handled in a second.").Default("1000").Int(),
		numWorkers:    cmd.Flag("num-workers", "Number of workers.").Default("2").Int(),
		flushSize:     cmd.Flag("flush-size", "Maximum number of messages in one flush.").Default("100").Int(),
		flushInterval: cmd.Flag("flush-interval", "Maximum interval between two flushes.").Default("2s").Duration(),
		certPath:      cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:       cmd.Flag("key", "Path to TLS key.").Required().String(),
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
	}
	r.RegisterCommand(persister)
	return persister
}

func (p *persister) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()

	puller, err := p.createPuller(ctx, logger)
	if err != nil {
		return err
	}

	persister := pst.NewPersister(
		puller,
		pst.WithMaxMPS(*p.maxMPS),
		pst.WithNumWorkers(*p.numWorkers),
		pst.WithMetrics(registerer),
		pst.WithLogger(logger),
	)
	defer persister.Stop()
	go persister.Run() // nolint:errcheck

	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
		health.WithCheck("persister", persister.Check),
	)
	go healthChecker.Run(ctx)

	server := rpc.NewServer(healthChecker, *p.certPath, *p.keyPath,
		"metrics-event-persister",
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

func (p *persister) createPuller(ctx context.Context, logger *zap.Logger) (puller.Puller, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	client, err := pubsub.NewClient(ctx, *p.project, pubsub.WithLogger(logger))
	if err != nil {
		return nil, err
	}
	return client.CreatePuller(*p.subscription, *p.topic,
		pubsub.WithNumGoroutines(*p.pullerNumGoroutines),
		pubsub.WithMaxOutstandingMessages(*p.pullerMaxOutstandingMessages),
		pubsub.WithMaxOutstandingBytes(*p.pullerMaxOutstandingBytes),
	)
}
