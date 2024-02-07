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

package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/profiler"
	octrace "go.opencensus.io/trace"
	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/trace"
)

var (
	errCommandNotFound = errors.New("command not found")

	healthCheckSpanName     = "grpc.health.v1.Health.Check"
	pubsubAckSpanName       = "google.pubsub.v1.Subscriber.Acknowledge"
	pubsubModifyAckSpanName = "google.pubsub.v1.Subscriber.ModifyAckDeadline"
)

type App struct {
	name    string
	version string
	cmds    map[string]Command
	app     *kingpin.Application
}

func NewApp(name, desc, version, build string) *App {
	app := &App{
		name:    name,
		version: fmt.Sprintf("%s-%s", version, build),
		app:     kingpin.New(name, desc),
		cmds:    make(map[string]Command),
	}
	app.app.Version(app.version)
	app.app.DefaultEnvars()
	return app
}

func (a *App) Command(name string, desc string) *kingpin.CmdClause {
	return a.app.Command(name, desc)
}

func (a *App) RegisterCommand(cmd Command) {
	a.cmds[cmd.FullCommand()] = cmd
}

func (a *App) Run() error {
	logLevel := a.app.Flag("log-level", "The level of logging.").Default("info").Enum(log.Levels...)
	profile := a.app.Flag("profile", "If true enables uploading the profiles to Stackdriver.").Default("true").Bool()
	metricsPort := a.app.Flag("metrics-port", "Port to bind metrics server to.").Default("9002").Int()
	traceSamplingProbability := a.app.Flag(
		"trace-sampling-probability",
		"How offten we send traces to exporters.",
	).Default("0.01").Float()
	tracePubsubAckSamplingProbability := a.app.Flag(
		"trace-pubsub-ack-sampling-probability",
		"How offten we send traces of pubsub ack to exporters.",
	).Default("0.0001").Float()
	gcpTraceEnabled := a.app.Flag(
		"gcp-trace-enabled",
		"Enables sending trace data to GCP Trace service.",
	).Default("true").Bool()

	cmd, err := a.app.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	if a.cmds[cmd] == nil {
		return errCommandNotFound
	}

	serviceName := fmt.Sprintf("%s.%s", a.name, cmd)
	logger, err := log.NewLogger(
		log.WithLevel(*logLevel),
		log.WithServiceContext(serviceName, a.version),
	)
	if err != nil {
		return err
	}
	defer logger.Sync() // nolint:errcheck

	if *profile {
		err = profiler.Start(profiler.Config{
			Service:        serviceName,
			ServiceVersion: a.version},
		)
		if err != nil {
			logger.Error("Failed to start profiler", zap.Error(err))
			return err
		}
	}

	metrics := metrics.NewMetrics(
		*metricsPort,
		"/metrics",
		metrics.WithLogger(logger),
	)
	defer metrics.Stop()
	go metrics.Run() // nolint:errcheck

	if *gcpTraceEnabled {
		sd, err := trace.NewStackdriverExporter(serviceName, a.version, logger)
		if err != nil {
			logger.Error("Failed to create the Stackdriver exporter", zap.Error(err))
			return err
		}
		defer sd.Flush()
		octrace.RegisterExporter(sd)
	}
	octrace.ApplyConfig(octrace.Config{
		DefaultSampler: trace.NewSampler(
			trace.WithDefaultProbability(*traceSamplingProbability),
			trace.WithFilteringSampler(healthCheckSpanName, octrace.NeverSample()),
			trace.WithFilteringSampler(
				pubsubAckSpanName,
				octrace.ProbabilitySampler(*tracePubsubAckSamplingProbability),
			),
			trace.WithFilteringSampler(
				pubsubModifyAckSpanName,
				octrace.ProbabilitySampler(*tracePubsubAckSamplingProbability),
			),
		),
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	go func() {
		select {
		case s := <-ch:
			logger.Info("App is stopping due to signal", zap.Stringer("signal", s))
			cancel()
		case <-ctx.Done():
		}
	}()
	logger.Info(fmt.Sprintf("Running %s", serviceName))
	return a.cmds[cmd].Run(ctx, metrics, logger)
}
