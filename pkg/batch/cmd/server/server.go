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

package server

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	autoopsclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	"github.com/bucketeer-io/bucketeer/pkg/batch/api"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/experiment"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/notification"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/opsevent"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	notificationclient "github.com/bucketeer-io/bucketeer/pkg/notification/client"
	notificationsender "github.com/bucketeer-io/bucketeer/pkg/notification/sender"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender/notifier"
	opsexecutor "github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor"
	"github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/targetstore"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

const command = "server"

var serverShutDownTimeout = 10 * time.Second

type server struct {
	*kingpin.CmdClause
	// Common
	port             *int
	project          *string
	certPath         *string
	keyPath          *string
	serviceTokenPath *string
	timezone         *string
	refreshInterval  *time.Duration
	webURL           *string
	// MySQL
	mysqlUser   *string
	mysqlPass   *string
	mysqlHost   *string
	mysqlPort   *int
	mysqlDBName *string
	// gRPC service
	environmentService  *string
	experimentService   *string
	autoOpsService      *string
	eventCounterService *string
	featureService      *string
	notificationService *string
	// pubsub config
	domainSubscription           *string
	domainTopic                  *string
	pullerNumGoroutines          *int
	pullerMaxOutstandingMessages *int
	pullerMaxOutstandingBytes    *int
	runningDurationPerBatch      *time.Duration
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start batch server")
	server := &server{
		CmdClause:        cmd,
		port:             cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:          cmd.Flag("project", "Google Cloud project name.").String(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		timezone:         cmd.Flag("timezone", "Time zone").Required().String(),
		refreshInterval: cmd.Flag(
			"refresh-interval",
			"Interval between refreshing target objects.",
		).Default("1m").Duration(),
		webURL:      cmd.Flag("web-url", "Web console URL.").Required().String(),
		mysqlUser:   cmd.Flag("mysql-user", "MySQL user.").Required().String(),
		mysqlPass:   cmd.Flag("mysql-pass", "MySQL password.").Required().String(),
		mysqlHost:   cmd.Flag("mysql-host", "MySQL host.").Required().String(),
		mysqlPort:   cmd.Flag("mysql-port", "MySQL port.").Required().Int(),
		mysqlDBName: cmd.Flag("mysql-db-name", "MySQL database name.").Required().String(),
		environmentService: cmd.Flag(
			"environment-service",
			"bucketeer-environment-service address.",
		).Default("environment:9090").String(),
		experimentService: cmd.Flag(
			"experiment-service",
			"bucketeer-experiment-service address.",
		).Default("experiment:9090").String(),
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
		notificationService: cmd.Flag(
			"notification-service",
			"bucketeer-notification-service address.",
		).Default("notification:9090").String(),
		domainTopic: cmd.Flag(
			"domain-topic",
			"Google PubSub topic name of incoming domain events.").Required().String(),
		domainSubscription: cmd.Flag(
			"domain-subscription",
			"Google PubSub subscription name of incoming domain event.",
		).Required().String(),
		pullerNumGoroutines: cmd.Flag(
			"puller-num-goroutines",
			"Number of goroutines will be spawned to pull messages.",
		).Required().Int(),
		pullerMaxOutstandingMessages: cmd.Flag(
			"puller-max-outstanding-messages",
			"Maximum number of unprocessed messages.",
		).Required().Int(),
		pullerMaxOutstandingBytes: cmd.Flag(
			"puller-max-outstanding-bytes",
			"Maximum size of unprocessed messages.").Int(),
		runningDurationPerBatch: cmd.Flag(
			"running-duration-per-batch",
			"Duration of running domain event informer per batch.",
		).Required().Duration(),
	}
	r.RegisterCommand(server)
	return server
}

func (s *server) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	*s.serviceTokenPath = s.insertTelepresenceMountRoot(*s.serviceTokenPath)
	*s.keyPath = s.insertTelepresenceMountRoot(*s.keyPath)
	*s.certPath = s.insertTelepresenceMountRoot(*s.certPath)

	registerer := metrics.DefaultRegisterer()

	mysqlClient, err := s.createMySQLClient(ctx, registerer, logger)
	if err != nil {
		return err
	}

	creds, err := client.NewPerRPCCredentials(*s.serviceTokenPath)
	if err != nil {
		return err
	}

	notificationClient, err := notificationclient.NewClient(*s.notificationService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	environmentClient, err := environmentclient.NewClient(*s.environmentService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	eventCounterClient, err := ecclient.NewClient(*s.eventCounterService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
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

	autoOpsClient, err := autoopsclient.NewClient(*s.autoOpsService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	targetStore := targetstore.NewTargetStore(
		environmentClient,
		autoOpsClient,
		targetstore.WithRefreshInterval(*s.refreshInterval),
		targetstore.WithMetrics(registerer),
		targetstore.WithLogger(logger),
	)
	go targetStore.Run()

	autoOpsExecutor := opsexecutor.NewAutoOpsExecutor(
		autoOpsClient,
		opsexecutor.WithLogger(logger),
	)

	slackNotifier := notifier.NewSlackNotifier(*s.webURL)

	notificationSender := notificationsender.NewSender(
		notificationClient,
		[]notifier.Notifier{slackNotifier},
		notificationsender.WithLogger(logger),
	)

	domainEventPuller, err := s.createPuller(ctx, registerer, logger)
	if err != nil {
		return err
	}

	location, err := locale.GetLocation(*s.timezone)
	if err != nil {
		return err
	}

	service := api.NewBatchService(
		experiment.NewExperimentStatusUpdater(
			environmentClient,
			experimentClient,
			jobs.WithLogger(logger),
		),
		notification.NewExperimentRunningWatcher(
			environmentClient,
			experimentClient,
			notificationSender,
			jobs.WithTimeout(1*time.Minute),
			jobs.WithLogger(logger),
		),
		notification.NewFeatureStaleWatcher(
			environmentClient,
			featureClient,
			notificationSender,
			jobs.WithTimeout(1*time.Minute),
			jobs.WithLogger(logger),
		),
		notification.NewMAUCountWatcher(
			environmentClient,
			eventCounterClient,
			notificationSender,
			location,
			jobs.WithTimeout(60*time.Minute),
			jobs.WithLogger(logger),
		),
		opsevent.NewDatetimeWatcher(
			targetStore,
			autoOpsExecutor,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		opsevent.NewEventCountWatcher(
			mysqlClient,
			targetStore,
			eventCounterClient,
			featureClient,
			autoOpsExecutor,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		environmentClient,
		domainEventPuller,
		notificationSender,
		logger,
		notification.WithRunningDurationPerBatch(*s.runningDurationPerBatch),
	)

	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
	)
	go healthChecker.Run(ctx)

	server := rpc.NewServer(service, *s.certPath, *s.keyPath,
		"batch-server",
		rpc.WithPort(*s.port),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
		rpc.WithService(healthChecker),
		rpc.WithHandler("/health", healthChecker),
	)
	go server.Run()

	defer func() {
		server.Stop(serverShutDownTimeout)
		time.Sleep(serverShutDownTimeout)
		targetStore.Stop()
		notificationClient.Close()
		experimentClient.Close()
		environmentClient.Close()
		eventCounterClient.Close()
		featureClient.Close()
		autoOpsClient.Close()
		mysqlClient.Close()
	}()

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

func (s *server) insertTelepresenceMountRoot(path string) string {
	volumeRoot := os.Getenv("TELEPRESENCE_ROOT")
	if volumeRoot == "" {
		return path
	}
	return volumeRoot + path
}

func (s *server) createPuller(ctx context.Context,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (puller.Puller, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	pubsubClient, err := pubsub.NewClient(
		ctx,
		*s.project,
		pubsub.WithMetrics(registerer),
		pubsub.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}
	pubsubPuller, err := pubsubClient.CreatePuller(*s.domainSubscription, *s.domainTopic,
		pubsub.WithNumGoroutines(*s.pullerNumGoroutines),
		pubsub.WithMaxOutstandingMessages(*s.pullerMaxOutstandingMessages),
		pubsub.WithMaxOutstandingBytes(*s.pullerMaxOutstandingBytes),
	)
	if err != nil {
		return nil, err
	}
	return pubsubPuller, nil
}
