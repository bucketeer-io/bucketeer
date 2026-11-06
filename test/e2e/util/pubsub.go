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

package util

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/factory"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	redisv3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

// PubSubConfig points the e2e tests at the Pub/Sub backend and topics the services use.
// It is read from environment variables, so the e2e packages that do not seed events are unaffected.
type PubSubConfig struct {
	Type                string
	Project             string
	RedisAddr           string
	RedisPartitionCount int
	EvaluationTopic     string
	GoalTopic           string
	EvaluationDWHTopic  string
	GoalDWHTopic        string
	EvaluationOPSTopic  string
	GoalOPSTopic        string
}

// PubSubConfigFromEnv skips the test when E2E_PUBSUB_TYPE is not set.
func PubSubConfigFromEnv(t *testing.T) PubSubConfig {
	t.Helper()
	cfg := PubSubConfig{
		Type:                os.Getenv("E2E_PUBSUB_TYPE"),
		Project:             os.Getenv("E2E_PUBSUB_PROJECT"),
		RedisAddr:           os.Getenv("E2E_PUBSUB_REDIS_ADDR"),
		RedisPartitionCount: 16,
		EvaluationTopic:     os.Getenv("E2E_EVALUATION_TOPIC"),
		GoalTopic:           os.Getenv("E2E_GOAL_TOPIC"),
		EvaluationDWHTopic:  os.Getenv("E2E_EVALUATION_DWH_TOPIC"),
		GoalDWHTopic:        os.Getenv("E2E_GOAL_DWH_TOPIC"),
		EvaluationOPSTopic:  os.Getenv("E2E_EVALUATION_OPS_TOPIC"),
		GoalOPSTopic:        os.Getenv("E2E_GOAL_OPS_TOPIC"),
	}
	if cfg.Type == "" {
		t.Skip("E2E_PUBSUB_TYPE is not set")
	}
	if v := os.Getenv("E2E_PUBSUB_REDIS_PARTITION_COUNT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			t.Fatalf("Invalid E2E_PUBSUB_REDIS_PARTITION_COUNT: %v", err)
		}
		cfg.RedisPartitionCount = n
	}
	return cfg
}

// PublishEvents publishes events straight to a topic, bypassing the gateway.
func PublishEvents(
	ctx context.Context,
	t *testing.T,
	cfg PubSubConfig,
	topic string,
	events []*eventproto.Event,
) {
	t.Helper()
	opts := []factory.Option{factory.WithPubSubType(factory.PubSubType(cfg.Type))}
	switch factory.PubSubType(cfg.Type) {
	case factory.Google:
		opts = append(opts, factory.WithProjectID(cfg.Project))
	case factory.RedisStream:
		redisClient, err := redisv3.NewClient(cfg.RedisAddr)
		if err != nil {
			t.Fatalf("Failed to create redis client: %v", err)
		}
		defer redisClient.Close()
		opts = append(opts,
			factory.WithRedisClient(redisClient),
			factory.WithPartitionCount(cfg.RedisPartitionCount),
		)
	}
	client, err := factory.NewClient(ctx, opts...)
	if err != nil {
		t.Fatalf("Failed to create pubsub client: %v", err)
	}
	defer client.Close()
	p, err := client.CreatePublisherInProject(topic, cfg.Project)
	if err != nil {
		t.Fatalf("Failed to create publisher for topic %s: %v", topic, err)
	}
	defer p.Stop()
	messages := make([]publisher.Message, 0, len(events))
	for _, e := range events {
		messages = append(messages, e)
	}
	if errs := p.PublishMulti(ctx, messages); len(errs) > 0 {
		t.Fatalf("Failed to publish events to topic %s: %v", topic, errs)
	}
}

// NewEvaluationEvent builds the event the gateway publishes for an evaluation.
func NewEvaluationEvent(
	t *testing.T,
	environmentID, featureID string,
	featureVersion int32,
	userID, variationID, tag string,
	reason *featureproto.Reason,
) *eventproto.Event {
	t.Helper()
	evaluation, err := anypb.New(&eventproto.EvaluationEvent{
		Timestamp:      time.Now().Unix(),
		FeatureId:      featureID,
		FeatureVersion: featureVersion,
		UserId:         userID,
		VariationId:    variationID,
		User: &userproto.User{
			Id:   userID,
			Data: map[string]string{"appVersion": "0.1.0"}},
		Reason: reason,
		Tag:    tag,
	})
	if err != nil {
		t.Fatal(err)
	}
	return newEvent(t, environmentID, evaluation)
}

// NewGoalEvent builds the event the gateway publishes for a goal.
func NewGoalEvent(
	t *testing.T,
	environmentID, goalID, userID, tag string,
	value float64,
	timestamp int64,
) *eventproto.Event {
	t.Helper()
	goal, err := anypb.New(&eventproto.GoalEvent{
		Timestamp: timestamp,
		GoalId:    goalID,
		UserId:    userID,
		Value:     value,
		User: &userproto.User{
			Id:   userID,
			Data: map[string]string{"appVersion": "0.1.0"}},
		Tag: tag,
	})
	if err != nil {
		t.Fatal(err)
	}
	return newEvent(t, environmentID, goal)
}

func newEvent(t *testing.T, environmentID string, event *anypb.Any) *eventproto.Event {
	t.Helper()
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	return &eventproto.Event{Id: id.String(), Event: event, EnvironmentId: environmentID}
}
