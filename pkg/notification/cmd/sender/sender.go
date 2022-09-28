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

package sender

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	notificationclient "github.com/bucketeer-io/bucketeer/pkg/notification/client"
	notificationsender "github.com/bucketeer-io/bucketeer/pkg/notification/sender"
	batchinformer "github.com/bucketeer-io/bucketeer/pkg/notification/sender/informer/batch"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender/informer/batch/job"
	domaineventinformer "github.com/bucketeer-io/bucketeer/pkg/notification/sender/informer/domainevent"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender/notifier"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
)

const command = "sender"

type sender struct {
	*kingpin.CmdClause
	port                             *int
	project                          *string
	domainTopic                      *string
	domainSubscription               *string
	notificationService              *string
	environmentService               *string
	eventCounterService              *string
	featureService                   *string
	experimentService                *string
	scheduleFeatureStaleWatcher      *string
	scheduleExperimentRunningWatcher *string
	scheduleMAUCountWatcher          *string
	maxMPS                           *int
	numWorkers                       *int
	certPath                         *string
	keyPath                          *string
	serviceTokenPath                 *string
	pullerNumGoroutines              *int
	pullerMaxOutstandingMessages     *int
	pullerMaxOutstandingBytes        *int
	webURL                           *string
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the Notification Sender")
	sender := &sender{
		CmdClause:   cmd,
		port:        cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		project:     cmd.Flag("project", "Google Cloud project name.").Required().String(),
		domainTopic: cmd.Flag("domain-topic", "Google PubSub topic name of incoming domain events.").String(),
		domainSubscription: cmd.Flag(
			"domain-subscription",
			"Google PubSub subscription name of incoming domain event.",
		).String(),
		notificationService: cmd.Flag(
			"notification-service",
			"bucketeer-notification-service address.",
		).Default("notification:9090").String(),
		environmentService: cmd.Flag(
			"environment-service",
			"bucketeer-environment-service address.",
		).Default("environment:9090").String(),
		eventCounterService: cmd.Flag(
			"event-counter-service",
			"bucketeer-event-counter-service address.",
		).Default("event-counter:9090").String(),
		featureService: cmd.Flag(
			"feature-service",
			"bucketeer-feature-service address.",
		).Default("feature:9090").String(),
		experimentService: cmd.Flag(
			"experiment-service",
			"bucketeer-experiment-service address.",
		).Default("experiment:9090").String(),
		scheduleFeatureStaleWatcher: cmd.Flag(
			"schedule-feature-stale-watcher",
			"Cron format schedule for feature stale watcher.",
		).Default("0 0 1 * * MON").String(), // on every Monday 10:00am JST
		scheduleExperimentRunningWatcher: cmd.Flag(
			"schedule-experiment-running-watcher",
			"Cron format schedule for experiment running watcher.",
		).Default("0 0 1 * * *").String(), // on every day 10:00am JST
		scheduleMAUCountWatcher: cmd.Flag(
			"schedule-mau-count-watcher",
			"Cron format schedule for mau count watcher.",
		).Default("0 0 1 1 * *").String(), // on every month 1st 10:00am JST
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
		webURL:                    cmd.Flag("web-url", "Web console URL.").Required().String(),
	}
	r.RegisterCommand(sender)
	return sender
}

func (s *sender) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()

	*s.serviceTokenPath = s.insertTelepresenceMoutRoot(*s.serviceTokenPath)
	*s.keyPath = s.insertTelepresenceMoutRoot(*s.keyPath)
	*s.certPath = s.insertTelepresenceMoutRoot(*s.certPath)

	domainEventPuller, err := s.createPuller(ctx, registerer, logger)
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
	defer notificationClient.Close()

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
	defer environmentClient.Close()

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
	defer eventCounterClient.Close()

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
	defer experimentClient.Close()

	slackNotifier := notifier.NewSlackNotifier(*s.webURL)

	notificationSender := notificationsender.NewSender(
		notificationClient,
		[]notifier.Notifier{slackNotifier},
		notificationsender.WithLogger(logger),
	)

	domainEventInformer := domaineventinformer.NewDomainEventInformer(
		environmentClient,
		domainEventPuller,
		notificationSender,
		domaineventinformer.WithMetrics(registerer),
		domaineventinformer.WithLogger(logger),
	)
	defer domainEventInformer.Stop()
	go domainEventInformer.Run() // nolint:errcheck

	jobs := s.createJobs(
		environmentClient,
		featureClient,
		experimentClient,
		eventCounterClient,
		notificationSender,
		logger,
	)
	batchInformer, err := batchinformer.NewJobInformer(
		jobs,
		batchinformer.WithMetrics(registerer),
		batchinformer.WithLogger(logger),
	)
	if err != nil {
		return err
	}
	defer batchInformer.Stop()
	go batchInformer.Run() // nolint:errcheck

	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
		health.WithCheck("domain_event_informer", domainEventInformer.Check),
		health.WithCheck("batch_informer", batchInformer.Check),
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

func (s *sender) createPuller(
	ctx context.Context,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (puller.Puller, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	client, err := pubsub.NewClient(
		ctx,
		*s.project,
		pubsub.WithMetrics(registerer),
		pubsub.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}
	puller, err := client.CreatePuller(*s.domainSubscription, *s.domainTopic,
		pubsub.WithNumGoroutines(*s.pullerNumGoroutines),
		pubsub.WithMaxOutstandingMessages(*s.pullerMaxOutstandingMessages),
		pubsub.WithMaxOutstandingBytes(*s.pullerMaxOutstandingBytes),
	)
	if err != nil {
		return nil, err
	}
	return puller, nil
}

func (s *sender) createJobs(
	environmentClient environmentclient.Client,
	featureClient featureclient.Client,
	experimentClient experimentclient.Client,
	eventCounterClient ecclient.Client,
	notificationSender notificationsender.Sender,
	logger *zap.Logger) []*batchinformer.Job {
	return []*batchinformer.Job{
		{
			Cron: *s.scheduleFeatureStaleWatcher,
			Name: "feature_stale_watcher",
			Job: job.NewFeatureWatcher(
				environmentClient,
				featureClient,
				notificationSender,
				job.WithTimeout(1*time.Minute),
				job.WithLogger(logger)),
		},
		{
			Cron: *s.scheduleExperimentRunningWatcher,
			Name: "experiment_running_watcher",
			Job: job.NewExperimentRunningWatcher(
				environmentClient,
				experimentClient,
				notificationSender,
				job.WithTimeout(1*time.Minute),
				job.WithLogger(logger)),
		},
		{
			Cron: *s.scheduleMAUCountWatcher,
			Name: "mau_count",
			Job: job.NewMAUCountWatcher(
				environmentClient,
				eventCounterClient,
				notificationSender,
				job.WithTimeout(60*time.Minute),
				job.WithLogger(logger)),
		},
	}
}

func (s *sender) insertTelepresenceMoutRoot(path string) string {
	volumeRoot := os.Getenv("TELEPRESENCE_ROOT")
	if volumeRoot == "" {
		return path
	}
	return volumeRoot + path
}
