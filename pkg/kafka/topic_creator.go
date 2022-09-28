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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"

	storagekafka "github.com/bucketeer-io/bucketeer/pkg/storage/kafka"
)

type TopicCreator interface {
	CreateTopics(ctx context.Context) error
}

type options struct {
	logger            *zap.Logger
	partitionNum      int32
	replicationFactor int16
	minInSyncReplicas string
}

type Option func(*options)

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

func WithPartitionNum(n int32) Option {
	return func(opts *options) {
		opts.partitionNum = n
	}
}

func WithReplicationFactor(f int16) Option {
	return func(opts *options) {
		opts.replicationFactor = f
	}
}

func WithMinInSyncReplicas(n string) Option {
	return func(opts *options) {
		opts.minInSyncReplicas = n
	}
}

type topicCreator struct {
	client      *storagekafka.ClusterAdmin
	topicPrefix string
	opts        *options
	logger      *zap.Logger
}

func NewTopicCreator(client *storagekafka.ClusterAdmin, topicPrefix string, opts ...Option) TopicCreator {
	dopts := &options{
		logger:            zap.NewNop(),
		partitionNum:      3,
		replicationFactor: 3,
		minInSyncReplicas: "2",
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &topicCreator{
		client:      client,
		topicPrefix: topicPrefix,
		opts:        dopts,
		logger:      dopts.logger.Named("kafka"),
	}
}

func (tc *topicCreator) CreateTopics(ctx context.Context) error {
	for _, topic := range topics {
		topicName := storagekafka.TopicName(tc.topicPrefix, topic)
		topicDetail := &sarama.TopicDetail{
			NumPartitions:     tc.opts.partitionNum,
			ReplicationFactor: tc.opts.replicationFactor,
			ConfigEntries:     map[string]*string{"min.insync.replicas": &tc.opts.minInSyncReplicas},
		}
		if err := tc.client.CreateTopic(topicName, topicDetail); err != nil {
			tc.logger.Error("Failed to create topic", zap.Error(err),
				zap.String("topic", topicName))
			return err
		}
		tc.logger.Info("Suceeded to create topic",
			zap.String("topic", topicName))
	}
	return nil
}
