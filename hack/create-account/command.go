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

package main

import (
	"context"
	"errors"
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
	email                *string
	role                 *string
	environmentNamespace *string
	isAdmin              *bool
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("create", "Create a new account")
	command := &command{
		CmdClause:         cmd,
		certPath:          cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		serviceTokenPath:  cmd.Flag("service-token", "Path to service token file.").Required().String(),
		webGatewayAddress: cmd.Flag("web-gateway", "Address of web-gateway.").Required().String(),
		email:             cmd.Flag("email", "The email of an account.").Required().String(),
		role:              cmd.Flag("role", "The role of an account.").Required().Enum("VIEWER", "EDITOR", "OWNER"),
		environmentNamespace: cmd.Flag(
			"environment-namespace",
			"The environment namespace for Datestore namespace",
		).Required().String(),
		isAdmin: cmd.Flag("is-admin", "Is an account admin or not.").Default("false").Bool(),
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
	role, ok := accountproto.Account_Role_value[*c.role]
	if !ok {
		logger.Error("Wrong role parameter", zap.String("role", *c.role))
		return errors.New("wrong role parameter")
	}
	if *c.isAdmin {
		err := c.createAdminAccount(ctx, client, accountproto.Account_Role(role))
		if err != nil {
			logger.Error("Failed to create admin account", zap.Error(err))
			return err
		}
		logger.Info("Admin account created")
		return nil
	}
	err = c.createAccount(ctx, client, accountproto.Account_Role(role))
	if err != nil {
		logger.Error("Failed to create account", zap.Error(err),
			zap.String("environmentNamespace", *c.environmentNamespace))
		return err
	}
	logger.Info("Account created")
	return nil
}

func (c *command) createAdminAccount(
	ctx context.Context,
	client accountclient.Client,
	role accountproto.Account_Role,
) error {
	req := &accountproto.CreateAdminAccountRequest{
		Command: &accountproto.CreateAdminAccountCommand{Email: *c.email},
	}
	if _, err := client.CreateAdminAccount(ctx, req); err != nil {
		return err
	}
	return nil
}

func (c *command) createAccount(
	ctx context.Context,
	client accountclient.Client,
	role accountproto.Account_Role,
) error {
	req := &accountproto.CreateAccountRequest{
		Command: &accountproto.CreateAccountCommand{
			Email: *c.email,
			Role:  role,
		},
		EnvironmentNamespace: *c.environmentNamespace,
	}
	if _, err := client.CreateAccount(ctx, req); err != nil {
		return err
	}
	return nil
}

func createAccountClient(addr, cert, serviceToken string, logger *zap.Logger) (accountclient.Client, error) {
	creds, err := client.NewPerRPCCredentials(serviceToken)
	if err != nil {
		return nil, err
	}
	return accountclient.NewClient(addr, cert,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(10*time.Second),
		client.WithBlock(),
		client.WithLogger(logger),
	)
}
