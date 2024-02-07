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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package executor

import (
	"context"

	"go.uber.org/zap"

	autoopsclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type AutoOpsExecutor interface {
	Execute(ctx context.Context, environmentNamespace, ruleID string) error
}

type autoOpsExecutor struct {
	autoOpsClient autoopsclient.Client
	logger        *zap.Logger
}

func NewAutoOpsExecutor(autoOpsClient autoopsclient.Client, opts ...Option) AutoOpsExecutor {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &autoOpsExecutor{
		autoOpsClient: autoOpsClient,
		logger:        dopts.logger.Named("auto-ops-executor"),
	}
}

func (e *autoOpsExecutor) Execute(ctx context.Context, environmentNamespace, ruleID string) error {
	resp, err := e.autoOpsClient.ExecuteAutoOps(ctx, &autoopsproto.ExecuteAutoOpsRequest{
		EnvironmentNamespace:                environmentNamespace,
		Id:                                  ruleID,
		ChangeAutoOpsRuleTriggeredAtCommand: &autoopsproto.ChangeAutoOpsRuleTriggeredAtCommand{},
	})
	if err != nil {
		e.logger.Error("Failed to execute auto ops", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("ruleID", ruleID),
		)
		return err
	}
	if resp.AlreadyTriggered {
		e.logger.Debug("autoOpsRule has already triggered",
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("ruleID", ruleID),
		)
	}
	return nil
}
