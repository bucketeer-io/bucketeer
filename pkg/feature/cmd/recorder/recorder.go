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

package recorder

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	featurerecorder "github.com/bucketeer-io/bucketeer/pkg/feature/recorder"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

const command = "recorder"

type recorder struct {
	*kingpin.CmdClause
	port                         *int
	project                      *string
	mysqlUser                    *string
	mysqlPass                    *string
	mysqlHost                    *string
	mysqlPort                    *int
	mysqlDBName                  *string
	subscription                 *string
	topic                        *string
	maxMPS                       *int
	flushInterval                *time.Duration
	startupInterval              *time.Duration
	certPath                     *string
	keyPath                      *string
	pullerNumGoroutines          *int
	pullerMaxOutstandingMessages *int
	pullerMaxOutstandingBytes    *int
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start feature recorder")
	recorder := &recorder{
		CmdClause:    cmd,
		port:         cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:      cmd.Flag("project", "Google Cloud project name.").String(),
		mysqlUser:    cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:    cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:    cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:    cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName:  cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		subscription: cmd.Flag("subscription", "Google PubSub subscription name.").String(),
		topic:        cmd.Flag("topic", "Google PubSub topic name.").String(),
		maxMPS: cmd.Flag(
			"max-mps",
			"Maximum messages should be handled in a second.",
		).Default("5000").Int(),
		flushInterval:   cmd.Flag("flush-interval", "Interval between two flushes.").Default("1m").Duration(),
		startupInterval: cmd.Flag("startup-interval", "Interval to start workers").Default("1s").Duration(),
		certPath:        cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:         cmd.Flag("key", "Path to TLS key.").Required().String(),
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
	r.RegisterCommand(recorder)
	return recorder
}

func (r *recorder) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	*r.keyPath = r.insertTelepresenceMountRoot(*r.keyPath)
	*r.certPath = r.insertTelepresenceMountRoot(*r.certPath)

	puller, err := r.createPuller(ctx, logger)
	if err != nil {
		return err
	}

	registerer := metrics.DefaultRegisterer()
	mysqlClient, err := mysql.NewClient(
		ctx,
		*r.mysqlUser,
		*r.mysqlPass,
		*r.mysqlHost,
		*r.mysqlPort,
		*r.mysqlDBName,
		mysql.WithLogger(logger),
	)
	if err != nil {
		logger.Error("Failed to create mysql client", zap.Error(err))
		return err
	}
	defer mysqlClient.Close()

	recorder := featurerecorder.NewRecorder(puller, mysqlClient,
		featurerecorder.WithMaxMPS(*r.maxMPS),
		featurerecorder.WithLogger(logger),
		featurerecorder.WithMetrics(registerer),
		featurerecorder.WithFlushInterval(*r.flushInterval),
		featurerecorder.WithStartupInterval(*r.startupInterval),
	)
	defer recorder.Stop()
	go recorder.Run() // nolint:errcheck

	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
		health.WithCheck("recorder", recorder.Check),
	)
	go healthChecker.Run(ctx)

	server := rpc.NewServer(healthChecker, *r.certPath, *r.keyPath,
		"feature-recorder",
		rpc.WithPort(*r.port),
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

func (r *recorder) createPuller(ctx context.Context, logger *zap.Logger) (puller.Puller, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	client, err := pubsub.NewClient(ctx, *r.project, pubsub.WithLogger(logger))
	if err != nil {
		return nil, err
	}
	return client.CreatePuller(*r.subscription, *r.topic,
		pubsub.WithNumGoroutines(*r.pullerNumGoroutines),
		pubsub.WithMaxOutstandingMessages(*r.pullerMaxOutstandingMessages),
		pubsub.WithMaxOutstandingBytes(*r.pullerMaxOutstandingBytes),
	)
}

// for telepresence --swap-deployment
func (r *recorder) insertTelepresenceMountRoot(path string) string {
	volumeRoot := os.Getenv("TELEPRESENCE_ROOT")
	if volumeRoot == "" {
		return path
	}
	return volumeRoot + path
}
