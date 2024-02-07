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
	"fmt"
	"strings"
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
	certPath           *string
	serviceTokenPath   *string
	webGatewayAddress  *string
	name               *string
	urlCode            *string
	description        *string
	createEnvironments *[]string
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("create", "Create a new project")
	command := &command{
		CmdClause:          cmd,
		certPath:           cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		serviceTokenPath:   cmd.Flag("service-token", "Path to service token file.").Required().String(),
		webGatewayAddress:  cmd.Flag("web-gateway", "Address of web-gateway.").Required().String(),
		name:               cmd.Flag("name", "Name of an project.").Required().String(),
		urlCode:            cmd.Flag("url-code", "URL code for the project.").Required().String(),
		description:        cmd.Flag("description", "(optional) Description of an project.").String(),
		createEnvironments: cmd.Flag("create-environment", "A list of environment names to be created").Strings(),
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
	// create project
	req := &environmentproto.CreateProjectRequest{
		Command: &environmentproto.CreateProjectCommand{
			Name:        *c.name,
			UrlCode:     strings.ToLower(strings.ReplaceAll(*c.urlCode, " ", "-")),
			Description: *c.description,
		},
	}
	resp, err := client.CreateProject(ctx, req)
	if err != nil {
		logger.Error("Failed to create project", zap.Error(err))
		return err
	}
	projID := resp.Project.Id
	logger.Info(fmt.Sprintf("%s project created. ID: %s", resp.Project.Name, projID))
	// create environments (optional)
	for _, envName := range *c.createEnvironments {
		req := &environmentproto.CreateEnvironmentV2Request{
			Command: &environmentproto.CreateEnvironmentV2Command{
				Name:      envName,
				UrlCode:   strings.ToLower(strings.ReplaceAll(envName, " ", "-")),
				ProjectId: projID,
			},
		}
		resp, err := client.CreateEnvironmentV2(ctx, req)
		if err != nil {
			logger.Error("Failed to create environment", zap.Error(err))
			return err
		}
		logger.Info(fmt.Sprintf(
			"%s environment created. ID: %s",
			resp.Environment.Name, resp.Environment.Id,
		))
	}
	logger.Info("Project created successfully")
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
