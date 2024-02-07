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

type ProgressiveRolloutExecutor interface {
	ExecuteProgressiveRollout(
		ctx context.Context,
		environmentNamespace, ruleID, scheduleID string,
	) error
}

type progressiveRolloutExecutor struct {
	autoOpsClient autoopsclient.Client
	logger        *zap.Logger
}

func NewProgressiveRolloutExecutor(autoOpsClient autoopsclient.Client, opts ...Option) ProgressiveRolloutExecutor {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &progressiveRolloutExecutor{
		autoOpsClient: autoOpsClient,
		logger:        dopts.logger.Named("progressive-rollout-executor"),
	}
}

func (e *progressiveRolloutExecutor) ExecuteProgressiveRollout(
	ctx context.Context,
	environmentNamespace, progressiveRolloutID, scheduleID string,
) error {
	_, err := e.autoOpsClient.ExecuteProgressiveRollout(ctx, &autoopsproto.ExecuteProgressiveRolloutRequest{
		EnvironmentNamespace: environmentNamespace,
		Id:                   progressiveRolloutID,
		ChangeProgressiveRolloutTriggeredAtCommand: &autoopsproto.ChangeProgressiveRolloutScheduleTriggeredAtCommand{
			ScheduleId: scheduleID,
		},
	})
	if err != nil {
		e.logger.Error("Failed to execute ProgressiveRollout", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("progressiveRolloutID", progressiveRolloutID),
		)
		return err
	}
	return nil
}
