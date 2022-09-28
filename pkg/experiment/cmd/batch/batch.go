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

package batch

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	experimentjob "github.com/bucketeer-io/bucketeer/pkg/experiment/batch/job"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/job"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
)

const command = "batch"

type batch struct {
	*kingpin.CmdClause
	port               *int
	project            *string
	environmentService *string
	experimentService  *string
	certPath           *string
	keyPath            *string
	serviceTokenPath   *string
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start batch layer")
	batch := &batch{
		CmdClause: cmd,
		port:      cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:   cmd.Flag("project", "Google Cloud project name.").String(),
		environmentService: cmd.Flag(
			"environment-service",
			"bucketeer-environment-service address.",
		).Default("environment:9090").String(),
		experimentService: cmd.Flag(
			"experiment-service",
			"bucketeer-experiment-service address.",
		).Default("experiment:9090").String(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
	}
	r.RegisterCommand(batch)
	return batch
}

func (b *batch) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	*b.serviceTokenPath = b.insertTelepresenceMountRoot(*b.serviceTokenPath)
	*b.keyPath = b.insertTelepresenceMountRoot(*b.keyPath)
	*b.certPath = b.insertTelepresenceMountRoot(*b.certPath)

	registerer := metrics.DefaultRegisterer()

	creds, err := client.NewPerRPCCredentials(*b.serviceTokenPath)
	if err != nil {
		return err
	}

	clientOptions := []client.Option{
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30 * time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	}
	environmentClient, err := environmentclient.NewClient(*b.environmentService, *b.certPath, clientOptions...)
	if err != nil {
		return err
	}
	defer environmentClient.Close()

	experimentClient, err := experimentclient.NewClient(*b.experimentService, *b.certPath, clientOptions...)
	if err != nil {
		return err
	}
	defer experimentClient.Close()

	manager := job.NewManager(
		registerer,
		"experiment_batch",
		logger,
	)
	defer manager.Stop()
	err = b.registerJobs(manager, environmentClient, experimentClient, logger)
	if err != nil {
		return err
	}
	go manager.Run() // nolint:errcheck

	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
	)
	go healthChecker.Run(ctx)

	server := rpc.NewServer(healthChecker, *b.certPath, *b.keyPath,
		rpc.WithPort(*b.port),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
		rpc.WithHandler("/health", healthChecker),
	)
	defer server.Stop(10 * time.Second)
	go server.Run()

	<-ctx.Done()
	return nil
}

func (b *batch) registerJobs(
	m *job.Manager,
	environmentClient environmentclient.Client,
	experimentClient experimentclient.Client,
	logger *zap.Logger) error {

	jobs := []struct {
		name string
		cron string
		job  job.Job
	}{
		{
			cron: "0 * * * * *",
			name: "experiment_status_updater",
			job: experimentjob.NewExperimentStatusUpdater(
				environmentClient,
				experimentClient,
				experimentjob.WithLogger(logger)),
		},
	}
	for i := range jobs {
		if err := m.AddCronJob(jobs[i].name, jobs[i].cron, jobs[i].job); err != nil {
			logger.Error("Failed to add cron job",
				zap.String("name", jobs[i].name),
				zap.String("cron", jobs[i].cron),
				zap.Error(err))
			return err
		}
	}
	return nil
}

// for telepresence --swap-deployment
func (b *batch) insertTelepresenceMountRoot(path string) string {
	volumeRoot := os.Getenv("TELEPRESENCE_ROOT")
	if volumeRoot == "" {
		return path
	}
	return volumeRoot + path
}
