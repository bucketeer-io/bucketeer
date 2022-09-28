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
	"regexp"
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/auth/api"
	"github.com/bucketeer-io/bucketeer/pkg/auth/oidc"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/token"
)

const command = "server"

type server struct {
	*kingpin.CmdClause
	port                *int
	accountService      *string
	certPath            *string
	keyPath             *string
	serviceTokenPath    *string
	oauthPrivateKeyPath *string
	oauthClientID       *string
	oauthClientSecret   *string
	oauthRedirectURLs   *[]string
	oauthIssuer         *string
	oauthIssuerCertPath *string
	emailFilter         *string
}

func RegisterCommand(r cli.CommandRegistry, p cli.ParentCommand) cli.Command {
	cmd := p.Command(command, "Start the server")
	server := &server{
		CmdClause: cmd,
		port:      cmd.Flag("port", "Port to bind to.").Default("9090").Int(),
		accountService: cmd.Flag(
			"account-service",
			"bucketeer-account-service address.",
		).Default("account:9090").String(),
		certPath:         cmd.Flag("cert", "Path to TLS certificate.").Required().String(),
		keyPath:          cmd.Flag("key", "Path to TLS key.").Required().String(),
		serviceTokenPath: cmd.Flag("service-token", "Path to service token.").Required().String(),
		oauthPrivateKeyPath: cmd.Flag(
			"oauth-private-key",
			"Path to private key for signing oauth token.",
		).Required().String(),
		oauthClientID: cmd.Flag("oauth-client-id", "The oauth clientID registered at dex.").Required().String(),
		oauthClientSecret: cmd.Flag(
			"oauth-client-secret",
			"The oauth client secret registered at Dex.",
		).Required().String(),
		oauthRedirectURLs:   cmd.Flag("oauth-redirect-urls", "The redirect urls registered at Dex.").Required().Strings(),
		oauthIssuer:         cmd.Flag("oauth-issuer", "The url of dex issuer.").Required().String(),
		oauthIssuerCertPath: cmd.Flag("oauth-issuer-cert", "Path to TLS certificate of issuer.").Required().String(),
		emailFilter:         cmd.Flag("email-filter", "Regexp pattern for filtering email.").String(),
	}
	r.RegisterCommand(server)
	return server
}

func (s *server) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	registerer := metrics.DefaultRegisterer()

	oidc, err := oidc.NewOIDC(
		ctx,
		*s.oauthIssuer,
		*s.oauthIssuerCertPath,
		*s.oauthClientID,
		*s.oauthClientSecret,
		*s.oauthRedirectURLs,
		oidc.WithLogger(logger))
	if err != nil {
		return err
	}

	signer, err := token.NewSigner(*s.oauthPrivateKeyPath)
	if err != nil {
		return err
	}

	creds, err := client.NewPerRPCCredentials(*s.serviceTokenPath)
	if err != nil {
		return err
	}

	accountClient, err := accountclient.NewClient(*s.accountService, *s.certPath,
		client.WithPerRPCCredentials(creds),
		client.WithDialTimeout(30*time.Second),
		client.WithBlock(),
		client.WithMetrics(registerer),
		client.WithLogger(logger))
	if err != nil {
		return err
	}
	defer accountClient.Close()

	serviceOptions := []api.Option{
		api.WithLogger(logger),
	}
	if *s.emailFilter != "" {
		filter, err := regexp.Compile(*s.emailFilter)
		if err != nil {
			return err
		}
		serviceOptions = append(serviceOptions, api.WithEmailFilter(filter))
	}
	service := api.NewAuthService(oidc, signer, accountClient, serviceOptions...)

	healthChecker := health.NewGrpcChecker(
		health.WithTimeout(time.Second),
		health.WithCheck("metrics", metrics.Check),
	)
	go healthChecker.Run(ctx)

	server := rpc.NewServer(service, *s.certPath, *s.keyPath,
		rpc.WithPort(*s.port),
		rpc.WithMetrics(registerer),
		rpc.WithLogger(logger),
		rpc.WithService(healthChecker),
		rpc.WithHandler("/health", healthChecker),
	)
	defer server.Stop(10 * time.Second)
	go server.Run()

	<-ctx.Done()
	return nil
}
