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

	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cli"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	redisv3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
)

const (
	defaultScanCount = 100
	retryKeyPattern  = "goal_event_retry"
)

type command struct {
	*kingpin.CmdClause
	redisAddr     *string
	redisPassword *string
	environmentID *string
	scanCount     *int64
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("delete", "Delete Redis goal event retry keys for a specific environment")
	command := &command{
		CmdClause:     cmd,
		redisAddr:     cmd.Flag("redis-addr", "Redis server address (host:port).").Required().String(),
		redisPassword: cmd.Flag("redis-password", "Redis password.").Default("").String(),
		environmentID: cmd.Flag("environment-id", "Environment ID to delete retry keys for.").Required().String(),
		scanCount:     cmd.Flag("scan-count", "Number of keys to scan per iteration.").Default("100").Int64(),
	}
	r.RegisterCommand(command)
	return command
}

func (c *command) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	// Construct the key pattern: {environmentID}:goal_event_retry:*
	keyPattern := *c.environmentID + ":" + retryKeyPattern + ":*"

	logger.Info("Starting Redis retry keys deletion",
		zap.String("redis-addr", *c.redisAddr),
		zap.String("environment-id", *c.environmentID),
		zap.String("key-pattern", keyPattern),
	)

	opts := []redisv3.Option{
		redisv3.WithLogger(logger),
	}
	if *c.redisPassword != "" {
		opts = append(opts, redisv3.WithPassword(*c.redisPassword))
	}

	client, err := redisv3.NewClient(*c.redisAddr, opts...)
	if err != nil {
		logger.Error("Failed to create Redis client", zap.Error(err))
		return err
	}
	defer client.Close()

	totalDeleted := 0
	var cursor uint64 = 0

	for {
		nextCursor, keys, err := client.Scan(cursor, keyPattern, *c.scanCount)
		if err != nil {
			logger.Error("Failed to scan Redis keys", zap.Error(err))
			return err
		}

		for _, key := range keys {
			if err := client.Del(key); err != nil {
				logger.Error("Failed to delete key",
					zap.String("key", key),
					zap.Error(err),
				)
				return err
			}
			totalDeleted++
			logger.Debug("Deleted key", zap.String("key", key))
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	if totalDeleted == 0 {
		logger.Info("No keys found matching pattern",
			zap.String("pattern", keyPattern),
		)
	} else {
		logger.Info("Successfully deleted Redis retry keys",
			zap.Int("total-deleted", totalDeleted),
			zap.String("pattern", keyPattern),
		)
	}

	return nil
}
