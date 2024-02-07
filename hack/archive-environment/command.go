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

package main

import (
	"context"
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

type command struct {
	*kingpin.CmdClause
	certPath          *string
	serviceTokenPath  *string
	webGatewayAddress *string
	id                *string
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("archive", "Archive a environment")
	command := &command{
		CmdClause:         cmd,
		certPath:          cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		serviceTokenPath:  cmd.Flag("service-token", "Path to service token file.").Required().String(),
		webGatewayAddress: cmd.Flag("web-gateway", "Address of web-gateway.").Required().String(),
		id:                cmd.Flag("id", "Id of an environment.").Required().String(),
	}
	r.RegisterCommand(command)
	return command
}

func (c *command) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	client, err := createEnvironmentClient(*c.webGatewayAddress, *c.certPath, *c.serviceTokenPath, logger)
	if err != nil {
		logger.Error("Failed to create environment client", zap.Error(err))
		return err
	}
	defer client.Close()
	req := &environmentproto.ArchiveEnvironmentV2Request{
		Id:      *c.id,
		Command: &environmentproto.ArchiveEnvironmentV2Command{},
	}
	if _, err = client.ArchiveEnvironmentV2(ctx, req); err != nil {
		logger.Error("Failed to delete environment", zap.Error(err))
		return err
	}
	logger.Info("Environment created")
	return nil
}

func createEnvironmentClient(addr, cert, serviceToken string, logger *zap.Logger) (environmentclient.Client, error) {
	creds, err := client.NewPerRPCCredentials(serviceToken)
	if err != nil {
		return nil, err
	}
	return environmentclient.NewClient(addr, cert,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(10*time.Second),
		client.WithBlock(),
		client.WithLogger(logger),
	)
}
