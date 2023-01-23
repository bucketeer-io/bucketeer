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
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/eventpersisterops/persister"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
)

const (
	command = "server"
)

type server struct {
	*kingpin.CmdClause
	port *int
	// option
	maxMPS        *int
	numWorkers    *int
	flushSize     *int
	flushInterval *time.Duration
	flushTimeout  *time.Duration
	// rpc
	serviceTokenPath *string
	certPath         *string
	keyPath          *string
}

func RegisterServerCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the server")
	server := &server{
		CmdClause:  cmd,
		port:       cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		maxMPS:     cmd.Flag("max-mps", "Maximum messages should be handled in a second.").Default("1000").Int(),
		numWorkers: cmd.Flag("num-workers", "Number of workers.").Default("2").Int(),
		flushSize: cmd.Flag(
			"flush-size",
			"Maximum number of messages to batch before writing to datastore.",
		).Default("50").Int(),
		flushInterval:    cmd.Flag("flush-interval", "Maximum interval between two flushes.").Default("5s").Duration(),
		flushTimeout:     cmd.Flag("flush-timeout", "Maximum time for a flush to finish.").Default("20s").Duration(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
	}
	r.RegisterCommand(server)
	return server
}

func (s *server) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()

	p := persister.NewPersister()
	defer p.Stop()

	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
		health.WithCheck("persister", p.Check),
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
