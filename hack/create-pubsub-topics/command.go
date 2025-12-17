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
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cli"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
)

type command struct {
	*kingpin.CmdClause
	pubsubEmulatorHost *string
	project            *string
	topics             *[]string
}

func registerCommand(r cli.CommandRegistry, p cli.ParentCommand) *command {
	cmd := p.Command("create", "Create topics in PubSub Emulator")
	command := &command{
		CmdClause:          cmd,
		pubsubEmulatorHost: cmd.Flag("pubsub-emulator-host", "PubSub Emulator Host").Default("localhost:8089").String(),
		project:            cmd.Flag("project", "Project ID").Default("bucketeer-test").String(),
		topics:             cmd.Flag("topic", "Topic name to create (can be specified multiple times)").Strings(),
	}
	r.RegisterCommand(command)
	return command
}

func (c *command) Run(ctx context.Context, metrics metrics.Metrics, logger *zap.Logger) error {
	// Set PUBSUB_EMULATOR_HOST environment variable
	emulatorHost := *c.pubsubEmulatorHost
	if !strings.Contains(emulatorHost, "://") {
		emulatorHost = "http://" + emulatorHost
	}
	os.Setenv("PUBSUB_EMULATOR_HOST", emulatorHost)

	// Default topics if none specified
	topics := *c.topics
	if len(topics) == 0 {
		topics = []string{"domain", "goal", "evaluation", "user", "metrics"}
	}

	// Create PubSub client
	client, err := pubsub.NewClient(ctx, *c.project, option.WithoutAuthentication())
	if err != nil {
		logger.Error("Failed to create PubSub client", zap.Error(err))
		return err
	}
	defer client.Close()

	// Create topics
	for _, topicName := range topics {
		if err := c.createTopic(ctx, client, topicName, logger); err != nil {
			if isAlreadyExistsError(err) {
				logger.Info("Topic already exists, skipping", zap.String("topic", topicName))
				continue
			}
			logger.Error("Failed to create topic", zap.String("topic", topicName), zap.Error(err))
			return err
		}
		logger.Info("Successfully created topic", zap.String("topic", topicName))
	}

	logger.Info("All topics created successfully")
	return nil
}

func (c *command) createTopic(ctx context.Context, client *pubsub.Client, topicName string, logger *zap.Logger) error {
	// Use a longer timeout for topic operations
	topicCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	topic := client.Topic(topicName)

	// Try to create the topic directly (idempotent for emulator)
	// The emulator will return an error if it already exists, which we handle
	_, err := client.CreateTopic(topicCtx, topicName)
	if err != nil {
		// Check if topic already exists
		existsCtx, existsCancel := context.WithTimeout(ctx, 30*time.Second)
		defer existsCancel()
		exists, existsErr := topic.Exists(existsCtx)
		if existsErr != nil {
			return fmt.Errorf("failed to check if topic exists: %w (create error: %v)", existsErr, err)
		}
		if exists {
			logger.Info("Topic already exists", zap.String("topic", topicName))
			return nil
		}
		// If it doesn't exist and create failed, return the error
		return fmt.Errorf("failed to create topic: %w", err)
	}

	logger.Info("Successfully created topic", zap.String("topic", topicName))
	return nil
}

func isAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	// Check for common "already exists" patterns
	return strings.Contains(errStr, "already exists") ||
		strings.Contains(errStr, "already created") ||
		strings.Contains(errStr, "duplicate") ||
		strings.Contains(errStr, "409") // HTTP 409 Conflict
}
