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

package opsevent

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	aoclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	autoopsdomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	envclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	"github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	envproto "github.com/bucketeer-io/bucketeer/proto/environment"
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
	w.logger.Debug("Start datetime watcher")
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
			w.logger.Debug("[datetime] auto ops rule",
				zap.String("autoOpsRuleId", a.Id),
				zap.String("environmentNamespace", env.Id),
				zap.String("featureId", a.FeatureId),
				zap.String("status", a.AutoOpsStatus.String()),
				zap.Any("clauses", a.Clauses))

			if !aor.HasExecuteClause() {
				continue
			}
			executeClause, err := w.getExecuteDateTimeClause(ctx, env.Id, aor)
			if err != nil {
				lastErr = err
			}
			if executeClause == nil {
				continue
			}
			w.logger.Debug("Execute auto ops rule",
				zap.String("featureId", a.FeatureId),
				zap.String("autoOpsRuleId", a.Id),
				zap.String("executeClauseId", executeClause.Id),
				zap.String("executeActionType", executeClause.ActionType.String()),
			)
			if err = w.autoOpsExecutor.Execute(ctx, env.Id, a.Id, executeClause); err != nil {
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
	environmentNamespace string,
) ([]*autoopsproto.AutoOpsRule, error) {
	resp, err := w.aoClient.ListAutoOpsRules(ctx, &autoopsproto.ListAutoOpsRulesRequest{
		PageSize:             0,
		EnvironmentNamespace: environmentNamespace,
	})
	if err != nil {
		return nil, err
	}
	return resp.AutoOpsRules, nil
}

func (w *datetimeWatcher) getExecuteDateTimeClause(
	ctx context.Context,
	environmentNamespace string,
	a *autoopsdomain.AutoOpsRule,
) (*autoopsproto.Clause, error) {
	datetimeClauses, err := a.ExtractDatetimeClauses()
	if err != nil {
		w.logger.Error("Failed to extract datetime clauses", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("featureId", a.FeatureId),
			zap.String("autoOpsRuleId", a.Id),
		)
		return nil, err
	}
	nowTimestamp := time.Now().Unix()
	var latestExecuteClause *autoopsproto.Clause
	var latestDatetime = int64(0)
	for _, c := range datetimeClauses {
		datetimeClause, err := a.UnmarshalDatetimeClause(c)
		if err != nil {
			w.logger.Error("Failed to unmarshal datetime clauses", zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
				zap.String("featureId", a.FeatureId),
				zap.String("clauseId", c.Id),
			)
			return nil, err
		}

		if w.assessRule(datetimeClause, nowTimestamp) {
			if datetimeClause.Time >= latestDatetime {
				latestDatetime = datetimeClause.Time
				latestExecuteClause = c
			}
		}
	}

	if latestExecuteClause != nil {
		w.logger.Info("Clause satisfies condition",
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("featureId", a.FeatureId),
			zap.String("autoOpsRuleId", a.Id),
			zap.Any("datetimeClauseId", latestExecuteClause.Id),
		)
		return latestExecuteClause, nil
	} else {
		return nil, nil
	}
}

func (w *datetimeWatcher) assessRule(datetimeClause *autoopsproto.DatetimeClause, nowTimestamp int64) bool {
	return datetimeClause.Time <= nowTimestamp
}
