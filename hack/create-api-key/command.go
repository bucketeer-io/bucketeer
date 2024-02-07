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
	"errors"
	"io/ioutil"
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
)

type command struct {
	*kingpin.CmdClause
	certPath             *string
	serviceTokenPath     *string
	webGatewayAddress    *string
	name                 *string
	role                 *string
	output               *string
	environmentNamespace *string
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("create", "Create a new api key")
	command := &command{
		CmdClause:         cmd,
		certPath:          cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		serviceTokenPath:  cmd.Flag("service-token", "Path to service token file.").Required().String(),
		webGatewayAddress: cmd.Flag("web-gateway", "Address of web-gateway.").Required().String(),
		name:              cmd.Flag("name", "The name of key.").Required().String(),
		role:              cmd.Flag("role", "The role of key.").Default("SDK").Enum("SDK", "SERVICE"),
		output:            cmd.Flag("output", "Path of file to write api key.").Required().String(),
		environmentNamespace: cmd.Flag(
			"environment-namespace",
			"The environment namespace to store api key",
		).Required().String(),
	}
	r.RegisterCommand(command)
	return command
}

func (c *command) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	client, err := createAccountClient(*c.webGatewayAddress, *c.certPath, *c.serviceTokenPath, logger)
	if err != nil {
		logger.Error("Failed to create account client", zap.Error(err))
		return err
	}
	role, ok := accountproto.APIKey_Role_value[*c.role]
	if !ok {
		logger.Error("Wrong role parameter", zap.String("role", *c.role))
		return errors.New("wrong role parameter")
	}
	resp, err := client.CreateAPIKey(ctx, &accountproto.CreateAPIKeyRequest{
		Command: &accountproto.CreateAPIKeyCommand{
			Name: *c.name,
			Role: accountproto.APIKey_Role(role),
		},
		EnvironmentNamespace: *c.environmentNamespace,
	})
	if err != nil {
		logger.Error("Failed to create api key", zap.Error(err))
		return err
	}
	if err := ioutil.WriteFile(*c.output, []byte(resp.ApiKey.Id), 0644); err != nil {
		logger.Error("Failed to write key to file", zap.Error(err), zap.String("output", *c.output))
		return err
	}
	logger.Info("Key generated")
	return nil
}

func createAccountClient(addr, cert, serviceToken string, logger *zap.Logger) (accountclient.Client, error) {
	creds, err := client.NewPerRPCCredentials(serviceToken)
	if err != nil {
		return nil, err
	}
	return accountclient.NewClient(addr, cert,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithLogger(logger),
	)
}
