// Copyright 2026 The Bucketeer Authors.
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
	"fmt"
	"time"

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cli"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	v3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
)

type command struct {
	*kingpin.CmdClause
	srcAddress      *string
	destAddress     *string
	srcPassword     *string
	destPassword    *string
	overrideDestKey *bool

	srcTLSEnabled            *bool
	srcTLSCACert             *string
	srcTLSCert               *string
	srcTLSKey                *string
	srcTLSInsecureSkipVerify *bool

	destTLSEnabled            *bool
	destTLSCACert             *string
	destTLSCert               *string
	destTLSKey                *string
	destTLSInsecureSkipVerify *bool
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("copy", "Copy data from source Redis to destination Redis")
	command := &command{
		CmdClause:    cmd,
		srcAddress:   cmd.Flag("src-address", "Source Redis address").Required().String(),
		destAddress:  cmd.Flag("dest-address", "Destination Redis address").Required().String(),
		srcPassword:  cmd.Flag("src-password", "Source Redis password").String(),
		destPassword: cmd.Flag("dest-password", "Destination Redis password").String(),
		overrideDestKey: cmd.Flag("override-dest-key", "Override existing keys in the destination Redis").
			Default("false").
			Bool(),
		srcTLSEnabled: cmd.Flag(
			"src-tls-enabled",
			"Enable TLS when connecting to the source Redis server.",
		).Default("false").Bool(),
		srcTLSCACert: cmd.Flag(
			"src-tls-ca-cert",
			"Path to the source Redis TLS CA certificate file. Uses the system CA pool if unset.",
		).String(),
		srcTLSCert: cmd.Flag(
			"src-tls-cert",
			"Path to the source Redis TLS client certificate file (for mutual TLS).",
		).String(),
		srcTLSKey: cmd.Flag(
			"src-tls-key",
			"Path to the source Redis TLS client private key file (for mutual TLS).",
		).String(),
		srcTLSInsecureSkipVerify: cmd.Flag(
			"src-tls-insecure-skip-verify",
			"Skip source Redis server certificate verification. Not recommended for production.",
		).Default("false").Bool(),
		destTLSEnabled: cmd.Flag(
			"dest-tls-enabled",
			"Enable TLS when connecting to the destination Redis server.",
		).Default("false").Bool(),
		destTLSCACert: cmd.Flag(
			"dest-tls-ca-cert",
			"Path to the destination Redis TLS CA certificate file. Uses the system CA pool if unset.",
		).String(),
		destTLSCert: cmd.Flag(
			"dest-tls-cert",
			"Path to the destination Redis TLS client certificate file (for mutual TLS).",
		).String(),
		destTLSKey: cmd.Flag(
			"dest-tls-key",
			"Path to the destination Redis TLS client private key file (for mutual TLS).",
		).String(),
		destTLSInsecureSkipVerify: cmd.Flag(
			"dest-tls-insecure-skip-verify",
			"Skip destination Redis server certificate verification. Not recommended for production.",
		).Default("false").Bool(),
	}
	r.RegisterCommand(command)
	return command
}

func (c *command) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	srcClient, err := v3.NewClient(*c.srcAddress,
		v3.WithLogger(logger),
		v3.WithPassword(*c.srcPassword),
		v3.WithPoolSize(10),
		v3.WithMinIdleConns(5),
		v3.WithMaxRetries(3),
		v3.WithDialTimeout(10*time.Second),
		v3.WithTLS(v3.TLSConfig{
			Enabled:            *c.srcTLSEnabled,
			CACert:             *c.srcTLSCACert,
			Cert:               *c.srcTLSCert,
			Key:                *c.srcTLSKey,
			InsecureSkipVerify: *c.srcTLSInsecureSkipVerify,
		}),
	)
	if err != nil {
		logger.Error("Error creating source Redis client", zap.Error(err))
		return err
	}
	defer srcClient.Close()

	destClient, err := v3.NewClient(*c.destAddress,
		v3.WithLogger(logger),
		v3.WithPassword(*c.destPassword),
		v3.WithPoolSize(10),
		v3.WithMinIdleConns(5),
		v3.WithMaxRetries(3),
		v3.WithDialTimeout(10*time.Second),
		v3.WithTLS(v3.TLSConfig{
			Enabled:            *c.destTLSEnabled,
			CACert:             *c.destTLSCACert,
			Cert:               *c.destTLSCert,
			Key:                *c.destTLSKey,
			InsecureSkipVerify: *c.destTLSInsecureSkipVerify,
		}),
	)
	if err != nil {
		logger.Error("Error creating destination Redis client", zap.Error(err))
		return err
	}
	defer destClient.Close()

	if err := c.scanAndCopyBatch(srcClient, destClient, logger); err != nil {
		logger.Error("Error during scan and copy process", zap.Error(err))
		return err
	}

	logger.Info("Data copy completed")
	return nil
}

func (c *command) scanAndCopyBatch(src, dest v3.Client, logger *zap.Logger) error {
	var cursor uint64
	batchSize := int64(1000)
	totalCopied := 0

	for {
		nextCursor, keys, err := src.Scan(cursor, "*", batchSize)
		if err != nil {
			logger.Error(
				"Error scanning keys from source Redis",
				zap.Error(err),
				zap.Uint64("cursor", cursor),
			)
			return fmt.Errorf("error scanning keys from source Redis: %w", err)
		}

		copiedKeys, err := c.copyBatch(src, dest, keys, logger)
		if err != nil {
			logger.Error(
				"Error copying batch",
				zap.Error(err),
				zap.Uint64("cursor", cursor),
				zap.Int("copiedKeys", copiedKeys),
			)
		} else {
			totalCopied += copiedKeys
			logger.Info(
				"Successfully copied batch",
				zap.Uint64("cursor", cursor),
				zap.Int("copiedKeys", copiedKeys),
				zap.Int("totalCopied", totalCopied),
			)
		}
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	logger.Info(
		"Successfully copied total keys",
		zap.Int("totalCopied", totalCopied),
	)
	return nil
}

func (c *command) copyBatch(src, dest v3.Client, keys []string, logger *zap.Logger) (int, error) {
	copiedKeys := 0
	for _, key := range keys {
		dumpedValue, err := src.Dump(key)
		if err != nil {
			if errors.Is(err, v3.ErrNil) {
				logger.Info("Key not found", zap.String("key", key))
				continue
			}
			return copiedKeys, fmt.Errorf("error dumping key %s: %w", key, err)
		}

		exists, err := dest.Exists(key)
		if err != nil {
			return copiedKeys, fmt.Errorf("error checking key existence %s: %w", key, err)
		}

		if exists == 1 {
			if *c.overrideDestKey {
				if err := dest.Del(key); err != nil {
					return copiedKeys, fmt.Errorf("error deleting existing key %s: %w", key, err)
				}
			} else {
				logger.Info("Skipping existing key", zap.String("key", key))
				continue
			}
		}

		err = dest.Restore(key, 0, dumpedValue)

		if err != nil {
			return copiedKeys, fmt.Errorf("error restoring key %s: %w", key, err)
		}
		copiedKeys++
	}

	return copiedKeys, nil
}
