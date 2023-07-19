// Copyright 2023 The Bucketeer Authors.
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

	autoopsclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client"
	ftclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/job"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor"
	opseventjob "github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/job"
	"github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/targetstore"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

const command = "batch"

type batch struct {
	*kingpin.CmdClause
	port                    *int
	project                 *string
	mysqlUser               *string
	mysqlPass               *string
	mysqlHost               *string
	mysqlPort               *int
	mysqlDBName             *string
	environmentService      *string
	autoOpsService          *string
	eventCounterService     *string
	featureService          *string
	certPath                *string
	keyPath                 *string
	serviceTokenPath        *string
	refreshInterval         *time.Duration
	scheduleCountWatcher    *string
	scheduleDatetimeWatcher *string
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start batch layer")
	batch := &batch{
		CmdClause:   cmd,
		port:        cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:     cmd.Flag("project", "Google Cloud project name.").Required().String(),
		mysqlUser:   cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:   cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:   cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:   cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName: cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		environmentService: cmd.Flag(
			"environment-service",
			"bucketeer-environment-service address.",
		).Default("environment:9090").String(),
		autoOpsService: cmd.Flag(
			"auto-ops-service",
			"bucketeer-auto-ops-service address.",
		).Default("auto-ops:9090").String(),
		eventCounterService: cmd.Flag(
			"event-counter-service",
			"bucketeer-event-counter-service address.",
		).Default("event-counter-server:9090").String(),
		featureService: cmd.Flag(
			"feature-service",
			"bucketeer-feature-service address.",
		).Default("feature:9090").String(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		refreshInterval: cmd.Flag(
			"refresh-interval",
			"Interval between refreshing target objects.",
		).Default("10m").Duration(),
		scheduleCountWatcher: cmd.Flag(
			"schedule-count-watcher",
			"Cron style schedule for count watcher.",
		).Default("0,10,20,30,40,50 * * * * *").String(),
		scheduleDatetimeWatcher: cmd.Flag(
			"schedule-datetime-watcher",
			"Cron style schedule for datetime watcher.",
		).Default("0,10,20,30,40,50 * * * * *").String(),
	}
	r.RegisterCommand(batch)
	return batch
}

func (b *batch) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	*b.serviceTokenPath = b.insertTelepresenceMountRoot(*b.serviceTokenPath)
	*b.keyPath = b.insertTelepresenceMountRoot(*b.keyPath)
	*b.certPath = b.insertTelepresenceMountRoot(*b.certPath)

	registerer := metrics.DefaultRegisterer()

	mysqlClient, err := b.createMySQLClient(ctx, registerer, logger)
	if err != nil {
		return err
	}
	defer mysqlClient.Close()

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

	autoOpsClient, err := autoopsclient.NewClient(*b.autoOpsService, *b.certPath, clientOptions...)
	if err != nil {
		return err
	}
	defer autoOpsClient.Close()

	eventCounterClient, err := ecclient.NewClient(*b.eventCounterService, *b.certPath, clientOptions...)
	if err != nil {
		return err
	}
	defer eventCounterClient.Close()

	featureClient, err := ftclient.NewClient(*b.featureService, *b.certPath,
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

	targetStore := targetstore.NewTargetStore(
		environmentClient,
		autoOpsClient,
		targetstore.WithRefreshInterval(*b.refreshInterval),
		targetstore.WithMetrics(registerer),
		targetstore.WithLogger(logger),
	)
	defer targetStore.Stop()
	go targetStore.Run()

	autoOpsExecutor := executor.NewAutoOpsExecutor(
		autoOpsClient,
		executor.WithLogger(logger),
	)

	manager := job.NewManager(registerer, "ops_event_batch", logger)
	defer manager.Stop()
	err = b.registerJobs(manager,
		mysqlClient,
		targetStore,
		eventCounterClient,
		featureClient,
		autoOpsExecutor,
		logger,
	)
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
		"ops-event-batch",
		rpc.WithPort(*b.port),
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

func (b *batch) registerJobs(
	m *job.Manager,
	mysqlClient mysql.Client,
	targetStore targetstore.TargetStore,
	eventCounterClient ecclient.Client,
	featureClient ftclient.Client,
	autoOpsExecutor executor.AutoOpsExecutor,
	logger *zap.Logger) error {

	jobs := []struct {
		name string
		cron string
		job  job.Job
	}{
		{
			cron: *b.scheduleCountWatcher,
			name: "ops_event_count_watcher",
			job: opseventjob.NewCountWatcher(
				mysqlClient,
				targetStore,
				eventCounterClient,
				featureClient,
				autoOpsExecutor,
				opseventjob.WithTimeout(5*time.Minute),
				opseventjob.WithLogger(logger)),
		},
		{
			cron: *b.scheduleDatetimeWatcher,
			name: "datetime_watcher",
			job: opseventjob.NewDatetimeWatcher(
				targetStore,
				autoOpsExecutor,
				opseventjob.WithTimeout(5*time.Minute),
				opseventjob.WithLogger(logger)),
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

func (b *batch) createMySQLClient(
	ctx context.Context,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (mysql.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return mysql.NewClient(
		ctx,
		*b.mysqlUser, *b.mysqlPass, *b.mysqlHost,
		*b.mysqlPort,
		*b.mysqlDBName,
		mysql.WithLogger(logger),
		mysql.WithMetrics(registerer),
	)
}
