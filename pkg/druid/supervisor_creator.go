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
package druid

import (
	"context"

	"go.uber.org/zap"

	storagedruid "github.com/bucketeer-io/bucketeer/pkg/storage/druid"
)

type SupervisorCreator interface {
	CreateSupervisors(ctx context.Context) error
}

type options struct {
	maxRowsPerSegment int
	retentionPeriod   string
	logger            *zap.Logger
}

type Option func(*options)

func WithMaxRowsPerSegment(r int) Option {
	return func(opts *options) {
		opts.maxRowsPerSegment = r
	}
}

func WithRetentionPeriod(r string) Option {
	return func(opts *options) {
		opts.retentionPeriod = r
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type supervisorCreator struct {
	coordinatorClient *storagedruid.CoordinatorClient
	overlordClient    *storagedruid.OverlordClient
	datasourcePrefix  string
	kafkaURL          string
	kafkaTopicPrefix  string
	kafkaUsername     string
	kafkaPassword     string
	opts              *options
	logger            *zap.Logger
}

func NewSupervisorCreator(
	coordinatorClient *storagedruid.CoordinatorClient,
	overlordClient *storagedruid.OverlordClient,
	datasourcePrefix,
	kafkaURL,
	kafkaTopicPrefix,
	kafkaUsername,
	kafkaPassword string,
	opts ...Option) SupervisorCreator {

	dopts := &options{
		maxRowsPerSegment: 3000000,
		retentionPeriod:   "P2M",
		logger:            zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &supervisorCreator{
		coordinatorClient: coordinatorClient,
		overlordClient:    overlordClient,
		datasourcePrefix:  datasourcePrefix,
		kafkaURL:          kafkaURL,
		kafkaTopicPrefix:  kafkaTopicPrefix,
		kafkaUsername:     kafkaUsername,
		kafkaPassword:     kafkaPassword,
		opts:              dopts,
		logger:            dopts.logger.Named("druid"),
	}
}

func (c *supervisorCreator) CreateSupervisors(ctx context.Context) error {
	for _, sv := range EventSupervisors(
		c.datasourcePrefix,
		c.kafkaTopicPrefix,
		c.kafkaURL,
		c.kafkaUsername,
		c.kafkaPassword,
		c.opts.maxRowsPerSegment,
	) {
		if err := c.overlordClient.CreateOrUpdateSupervisor(ctx, sv, ""); err != nil {
			c.logger.Error("Failed to create suprvisor", zap.Error(err),
				zap.Any("supervisor", sv))
			return err
		}
		datasource := sv.DataSchema.Datasource
		retentionRule := retentionRule(c.opts.retentionPeriod)
		if err := c.coordinatorClient.CreateOrUpdateRetentionRule(ctx, datasource, retentionRule, ""); err != nil {
			c.logger.Error("Failed to set retention rule", zap.Error(err),
				zap.String("datastore", datasource))
			return err
		}
		c.logger.Info("Succeeded to create suprvisor",
			zap.Any("supervisor", sv))
	}
	return nil
}
