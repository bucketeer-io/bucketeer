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
	"os"
	"time"

	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/token"
)

type command struct {
	*kingpin.CmdClause
	keyPath  *string
	issuer   *string
	sub      *string
	audience *string
	email    *string
	role     *string
	output   *string
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("generate", "Generate a new service token")
	command := &command{
		CmdClause: cmd,
		keyPath:   cmd.Flag("key", "Path to the private keys.").Required().String(),
		issuer:    cmd.Flag("issuer", "Issuer url set in dex config.").Required().String(),
		sub:       cmd.Flag("sub", "Subject id.").Required().String(),
		audience:  cmd.Flag("audience", "Client id set in dex config.").Required().String(),
		email:     cmd.Flag("email", "Email will be set in token.").Required().String(),
		// FIXME: This should be removed in the future
		role:   cmd.Flag("role", "Role will be set in token.").Default("VIEWER").Enum("VIEWER", "EDITOR", "OWNER"),
		output: cmd.Flag("output", "Path of file to write service token.").Required().String(),
	}
	r.RegisterCommand(command)
	return command
}

func (c *command) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	signer, err := token.NewSigner(*c.keyPath)
	if err != nil {
		logger.Error("Failed to create signer", zap.Error(err))
		return err
	}
	idToken := &token.IDToken{
		Issuer:        *c.issuer,
		Subject:       *c.sub,
		Audience:      *c.audience,
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         *c.email,
		IsSystemAdmin: true,
	}
	signedIDToken, err := signer.Sign(idToken)
	if err != nil {
		logger.Error("Failed to sign token", zap.Error(err))
		return err
	}
	if err := os.WriteFile(*c.output, []byte(signedIDToken), 0644); err != nil {
		logger.Error("Failed to write token to file", zap.Error(err), zap.String("output", *c.output))
		return err
	}
	logger.Info("Token generated")
	return nil
}
