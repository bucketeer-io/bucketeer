// Copyright 2025 The Bucketeer Authors.
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

	"github.com/bucketeer-io/bucketeer/pkg/cli"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	v3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
)

type command struct {
	*kingpin.CmdClause
	srcAddress      *string
	destAddress     *string
	srcPassword     *string
	destPassword    *string
	overrideDestKey *bool
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
			return fmt.Errorf("error scanning keys from source Redis: %v", err)
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
			return copiedKeys, fmt.Errorf("error dumping key %s: %v", key, err)
		}

		exists, err := dest.Exists(key)
		if err != nil {
			return copiedKeys, fmt.Errorf("error checking key existence %s: %v", key, err)
		}

		if exists == 1 {
			if *c.overrideDestKey {
				if err := dest.Del(key); err != nil {
					return copiedKeys, fmt.Errorf("error deleting existing key %s: %v", key, err)
				}
			} else {
				logger.Info("Skipping existing key", zap.String("key", key))
				continue
			}
		}

		err = dest.Restore(key, 0, dumpedValue)

		if err != nil {
			return copiedKeys, fmt.Errorf("error restoring key %s: %v", key, err)
		}
		copiedKeys++
	}

	return copiedKeys, nil
}
