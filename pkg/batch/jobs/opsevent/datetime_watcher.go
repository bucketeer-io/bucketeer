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

package opsevent

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	aoclient "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client"
	autoopsdomain "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	envclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/batch/executor"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

type datetimeWatcher struct {
	envClient       envclient.Client
	aoClient        aoclient.Client
	autoOpsExecutor executor.AutoOpsExecutor
	opts            *jobs.Options
	logger          *zap.Logger
}

func NewDatetimeWatcher(
	envClient envclient.Client,
	aoClient aoclient.Client,
	autoOpsExecutor executor.AutoOpsExecutor,
	opts ...jobs.Option) jobs.Job {

	dopts := &jobs.Options{
		Timeout: 5 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &datetimeWatcher{
		envClient:       envClient,
		aoClient:        aoClient,
		autoOpsExecutor: autoOpsExecutor,
		opts:            dopts,
		logger:          dopts.Logger.Named("datetime-watcher"),
	}
}

func (w *datetimeWatcher) Run(ctx context.Context) (lastErr error) {
	ctx, cancel := context.WithTimeout(ctx, w.opts.Timeout)
	defer cancel()
	envs, err := w.listEnvironments(ctx)
	if err != nil {
		lastErr = err
		return
	}
	for _, env := range envs {
		autoOpsRules, err := w.listAutoOpsRules(ctx, env.Id)
		if err != nil {
			lastErr = err
			return
		}
		for _, a := range autoOpsRules {
			aor := &autoopsdomain.AutoOpsRule{AutoOpsRule: a}
			if aor.IsFinished() || aor.IsStopped() || aor.OpsType != autoopsproto.OpsType_SCHEDULE {
				continue
			}
			executeClauseID, err := w.getExecuteClauseId(ctx, env.Id, aor)
			if err != nil {
				lastErr = err
			}
			if executeClauseID == "" {
				continue
			}
			if err = w.autoOpsExecutor.Execute(ctx, env.Id, a.Id, executeClauseID); err != nil {
				lastErr = err
			}
		}
	}
	return
}

func (w *datetimeWatcher) listEnvironments(ctx context.Context) ([]*envproto.EnvironmentV2, error) {
	resp, err := w.envClient.ListEnvironmentsV2(ctx, &envproto.ListEnvironmentsV2Request{
		PageSize: 0,
		Archived: wrapperspb.Bool(false),
	})
	if err != nil {
		return nil, err
	}
	return resp.Environments, nil
}

func (w *datetimeWatcher) listAutoOpsRules(
	ctx context.Context,
	environmentId string,
) ([]*autoopsproto.AutoOpsRule, error) {
	resp, err := w.aoClient.ListAutoOpsRules(ctx, &autoopsproto.ListAutoOpsRulesRequest{
		PageSize:      0,
		EnvironmentId: environmentId,
	})
	if err != nil {
		return nil, err
	}
	return resp.AutoOpsRules, nil
}

func (w *datetimeWatcher) getExecuteClauseId(
	ctx context.Context,
	environmentId string,
	a *autoopsdomain.AutoOpsRule,
) (string, error) {
	var nextClauseID string
	for _, clause := range a.Clauses {
		// Skip clauses that were executed
		if clause.ExecutedAt != 0 {
			continue
		}
		// Set the next scheduled clause
		nextClauseID = clause.Id
		break
	}
	dateClauses, err := a.ExtractDatetimeClauses()
	if err != nil {
		return "", err
	}
	dateClause := dateClauses[nextClauseID]
	if dateClause == nil {
		w.logger.Error("Datetime clause is nil",
			zap.String("environmentId", environmentId),
			zap.String("featureId", a.FeatureId),
			zap.String("autoOpsRuleId", a.Id),
			zap.String("clauseId", nextClauseID),
		)
		return "", fmt.Errorf("datetime clause not found")
	}
	// Check if the scheduled time is older than now
	if dateClause.Time <= time.Now().Unix() {
		w.logger.Info("Scheduled operation satisfies the time condition",
			zap.String("environmentId", environmentId),
			zap.String("featureId", a.FeatureId),
			zap.String("autoOpsRuleId", a.Id),
			zap.String("clauseId", nextClauseID),
			zap.Int64("clauseTime", dateClause.Time),
			zap.Any("dateClauses", dateClause),
		)
		return nextClauseID, nil
	}
	w.logger.Debug("Scheduled operation does not satisfy the time condition",
		zap.String("environmentId", environmentId),
		zap.String("featureId", a.FeatureId),
		zap.String("autoOpsRuleId", a.Id),
		zap.String("clauseId", nextClauseID),
		zap.Int64("clauseTime", dateClause.Time),
		zap.Any("dateClauses", dateClause),
	)
	// Nothing to execute
	return "", nil
}
