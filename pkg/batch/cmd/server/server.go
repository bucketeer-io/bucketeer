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
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/calculator"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/experiment"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/mau"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/notification"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/opsevent"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/rediscounter"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	experimentcalculatorclient "github.com/bucketeer-io/bucketeer/pkg/experimentcalculator/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	notificationclient "github.com/bucketeer-io/bucketeer/pkg/notification/client"
	notificationsender "github.com/bucketeer-io/bucketeer/pkg/notification/sender"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender/notifier"
	"github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor"
	opsexecutor "github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
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
	environmentService          *string
	experimentService           *string
	autoOpsService              *string
	eventCounterService         *string
	featureService              *string
	notificationService         *string
	experimentCalculatorService *string
	// PubSub config
	domainSubscription           *string
	domainTopic                  *string
	pullerNumGoroutines          *int
	pullerMaxOutstandingMessages *int
	pullerMaxOutstandingBytes    *int
	runningDurationPerBatch      *time.Duration
	maxMPS                       *int
	// Redis
	redisServerName    *string
	redisAddr          *string
	redisPoolMaxIdle   *int
	redisPoolMaxActive *int
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
		experimentCalculatorService: cmd.Flag(
			"experiment-calculator-service",
			"bucketeer-experiment-calculator-service address.",
		).Default("experiment-calculator:9090").String(),
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
		maxMPS: cmd.Flag(
			"max-mps",
			"Maximum number of messages per second.",
		).Required().Int(),
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

	autoOpsExecutor := opsexecutor.NewAutoOpsExecutor(
		autoOpsClient,
		opsexecutor.WithLogger(logger),
	)

	progressiveRolloutExecutor := opsexecutor.NewProgressiveRolloutExecutor(
		autoOpsClient,
		executor.WithLogger(logger),
	)

	experimentCalculatorClient, err := experimentcalculatorclient.NewClient(*s.experimentCalculatorService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	slackNotifier := notifier.NewSlackNotifier(*s.webURL)

	notificationSender := notificationsender.NewSender(
		notificationClient,
		[]notifier.Notifier{slackNotifier},
		notificationsender.WithMetrics(registerer),
		notificationsender.WithLogger(logger),
	)

	location, err := locale.GetLocation(*s.timezone)
	if err != nil {
		return err
	}

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
			environmentClient,
			autoOpsClient,
			autoOpsExecutor,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		opsevent.NewEventCountWatcher(
			mysqlClient,
			environmentClient,
			autoOpsClient,
			eventCounterClient,
			featureClient,
			autoOpsExecutor,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		opsevent.NewProgressiveRolloutWacher(
			environmentClient,
			autoOpsClient,
			progressiveRolloutExecutor,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		rediscounter.NewRedisCounterDeleter(
			redisV3Cache,
			environmentClient,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		calculator.NewExperimentCalculate(
			environmentClient,
			experimentClient,
			experimentCalculatorClient,
			location,
			jobs.WithTimeout(5*time.Minute),
			jobs.WithLogger(logger),
		),
		mau.NewMAUSummarizer(
			mysqlClient,
			eventCounterClient,
			location,
			jobs.WithTimeout(30*time.Minute),
			jobs.WithLogger(logger),
		),
		mau.NewMAUPartitionDeleter(
			mysqlClient,
			location,
			jobs.WithTimeout(60*time.Minute),
			jobs.WithLogger(logger),
		),
		mau.NewMAUPartitionCreator(
			mysqlClient,
			location,
			jobs.WithTimeout(60*time.Minute),
			jobs.WithLogger(logger),
		),
		notification.NewDomainEventInformer(
			environmentClient,
			notificationSender,
			*s.maxMPS,
			*s.runningDurationPerBatch,
			*s.project,
			*s.domainSubscription,
			*s.domainTopic,
			*s.pullerNumGoroutines,
			*s.pullerMaxOutstandingMessages,
			*s.pullerMaxOutstandingBytes,
			notification.WithLogger(logger),
			notification.WithMetrics(registerer),
		),
		logger,
	)

	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
		health.WithCheck("redis", redisV3Client.Check),
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
		notificationClient.Close()
		experimentClient.Close()
		environmentClient.Close()
		eventCounterClient.Close()
		featureClient.Close()
		autoOpsClient.Close()
		experimentCalculatorClient.Close()
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
