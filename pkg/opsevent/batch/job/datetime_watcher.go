// Copyright 2023 The Bucketeer Authors.
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

package job

import (
	"context"
	"time"

	"go.uber.org/zap"

	autoopsdomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	environmentdomain "github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/job"
	"github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor"
	"github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/targetstore"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

type datetimeWatcher struct {
	environmentLister targetstore.EnvironmentLister
	autoOpsRuleLister targetstore.AutoOpsRuleLister
	autoOpsExecutor   executor.AutoOpsExecutor
	opts              *options
	logger            *zap.Logger
}

func NewDatetimeWatcher(
	targetStore targetstore.TargetStore,
	autoOpsExecutor executor.AutoOpsExecutor,
	opts ...Option) job.Job {

	dopts := &options{
		timeout: 5 * time.Minute,
		logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &datetimeWatcher{
		environmentLister: targetStore,
		autoOpsRuleLister: targetStore,
		autoOpsExecutor:   autoOpsExecutor,
		opts:              dopts,
		logger:            dopts.logger.Named("datetime-watcher"),
	}
}

func (w *datetimeWatcher) Run(ctx context.Context) (lastErr error) {
	ctx, cancel := context.WithTimeout(ctx, w.opts.timeout)
	defer cancel()
	environments := w.environmentLister.GetEnvironments(ctx)
	for _, env := range environments {
		autoOpsRules := w.autoOpsRuleLister.GetAutoOpsRules(ctx, env.Namespace)
		for _, a := range autoOpsRules {
			asmt, err := w.assessAutoOpsRule(ctx, env, a)
			if err != nil {
				lastErr = err
			}
			if !asmt {
				continue
			}
			if err = w.autoOpsExecutor.Execute(ctx, env.Namespace, a.Id); err != nil {
				lastErr = err
			}
		}
	}
	return
}

func (w *datetimeWatcher) assessAutoOpsRule(
	ctx context.Context,
	env *environmentdomain.Environment,
	a *autoopsdomain.AutoOpsRule,
) (bool, error) {
	datetimeClauses, err := a.ExtractDatetimeClauses()
	if err != nil {
		w.logger.Error("Failed to extract datetime clauses", zap.Error(err),
			zap.String("environmentNamespace", env.Namespace),
			zap.String("featureId", a.FeatureId),
			zap.String("autoOpsRuleId", a.Id),
		)
		return false, err
	}
	var lastErr error
	nowTimestamp := time.Now().Unix()
	for _, c := range datetimeClauses {
		if asmt := w.assessRule(c, nowTimestamp); asmt {
			w.logger.Info("Clause satisfies condition",
				zap.String("environmentNamespace", env.Namespace),
				zap.String("featureId", a.FeatureId),
				zap.String("autoOpsRuleId", a.Id),
				zap.Any("datetimeClause", c),
			)
			return true, nil
		}
	}
	return false, lastErr
}

func (w *datetimeWatcher) assessRule(datetimeClause *autoopsproto.DatetimeClause, nowTimestamp int64) bool {
	return datetimeClause.Time <= nowTimestamp
}
