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
	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client"
	ftclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/job"
	"github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor"
	"github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/targetstore"
	opseventdomain "github.com/bucketeer-io/bucketeer/pkg/opsevent/domain"
	v2os "github.com/bucketeer-io/bucketeer/pkg/opsevent/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
	ftproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type countWatcher struct {
	mysqlClient        mysql.Client
	environmentLister  targetstore.EnvironmentLister
	autoOpsRuleLister  targetstore.AutoOpsRuleLister
	eventCounterClient ecclient.Client
	featureClient      ftclient.Client
	autoOpsExecutor    executor.AutoOpsExecutor
	opts               *options
	logger             *zap.Logger
}

func NewCountWatcher(
	mysqlClient mysql.Client,
	targetStore targetstore.TargetStore,
	eventCounterClient ecclient.Client,
	featureClient ftclient.Client,
	autoOpsExecutor executor.AutoOpsExecutor,
	opts ...Option,
) job.Job {
	dopts := &options{
		timeout: 5 * time.Minute,
		logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &countWatcher{
		mysqlClient:        mysqlClient,
		environmentLister:  targetStore,
		autoOpsRuleLister:  targetStore,
		eventCounterClient: eventCounterClient,
		featureClient:      featureClient,
		autoOpsExecutor:    autoOpsExecutor,
		opts:               dopts,
		logger:             dopts.logger.Named("count-watcher"),
	}
}

func (w *countWatcher) Run(ctx context.Context) (lastErr error) {
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

func (w *countWatcher) assessAutoOpsRule(
	ctx context.Context,
	env *environmentdomain.Environment,
	a *autoopsdomain.AutoOpsRule,
) (bool, error) {
	opsEventRateClauses, err := a.ExtractOpsEventRateClauses()
	if err != nil {
		w.logger.Error("Failed to extract ops event rate clauses", zap.Error(err),
			zap.String("environmentNamespace", env.Namespace),
			zap.String("featureId", a.FeatureId),
			zap.String("autoOpsRuleId", a.Id),
		)
		return false, err
	}
	featureVersion, err := w.getLatestFeatureVersion(ctx, a.FeatureId, env.Namespace)
	if err != nil {
		w.logger.Error("Failed to get the latest feature version", zap.Error(err),
			zap.String("environmentNamespace", env.Namespace),
			zap.String("featureId", a.FeatureId),
			zap.String("autoOpsRuleId", a.Id),
		)
		return false, err
	}
	var lastErr error
	for id, c := range opsEventRateClauses {
		logFunc := func(msg string) {
			w.logger.Debug(msg,
				zap.String("environmentNamespace", env.Namespace),
				zap.String("featureId", a.FeatureId),
				zap.String("autoOpsRuleId", a.Id),
				zap.Any("opsEventRateClause", c),
			)
		}
		evaluationCount, err := w.getTargetOpsEvaluationCount(ctx,
			logFunc,
			env.Namespace,
			a.Id,
			id,
			a.FeatureId,
			c.VariationId,
			featureVersion,
		)
		if err != nil {
			lastErr = err
			continue
		}
		if evaluationCount == 0 {
			continue
		}
		opsEventCount, err := w.getTargetOpsGoalEventCount(
			ctx,
			logFunc,
			env.Namespace,
			a.Id,
			id,
			a.FeatureId,
			c.VariationId,
			featureVersion,
		)
		if err != nil {
			lastErr = err
			continue
		}
		if opsEventCount == 0 {
			continue
		}
		opsCount := opseventdomain.NewOpsCount(a.FeatureId, a.Id, id, opsEventCount, evaluationCount)
		if err = w.persistOpsCount(ctx, env.Namespace, opsCount); err != nil {
			lastErr = err
			continue
		}
		if asmt := w.assessRule(c, evaluationCount, opsEventCount); asmt {
			w.logger.Info("Clause satisfies condition",
				zap.String("environmentNamespace", env.Namespace),
				zap.String("featureId", a.FeatureId),
				zap.String("autoOpsRuleId", a.Id),
				zap.Any("opsEventRateClause", c),
			)
			return true, nil
		}
	}
	return false, lastErr
}

func (w *countWatcher) getLatestFeatureVersion(
	ctx context.Context,
	featureID, environmentNamespace string,
) (int32, error) {
	req := &ftproto.GetFeatureRequest{
		Id:                   featureID,
		EnvironmentNamespace: environmentNamespace,
	}
	resp, err := w.featureClient.GetFeature(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.Feature.Version, nil
}

func (w *countWatcher) assessRule(
	opsEventRateClause *autoopsproto.OpsEventRateClause,
	evaluationCount, opsCount int64,
) bool {
	rate := float64(opsCount) / float64(evaluationCount)
	if opsCount < opsEventRateClause.MinCount {
		return false
	}
	switch opsEventRateClause.Operator {
	case autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL:
		if rate >= opsEventRateClause.ThreadsholdRate {
			return true
		}
	case autoopsproto.OpsEventRateClause_LESS_OR_EQUAL:
		if rate <= opsEventRateClause.ThreadsholdRate {
			return true
		}
	}
	return false
}

func (w *countWatcher) getTargetOpsEvaluationCount(
	ctx context.Context,
	logFunc func(string),
	environmentNamespace, ruleID, clauseID, FeatureID, variationID string,
	featureVersion int32,
) (int64, error) {
	count, err := w.getEvaluationCount(
		ctx,
		environmentNamespace,
		ruleID,
		clauseID,
		FeatureID,
		variationID,
		featureVersion,
	)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		logFunc("Ops evaluation user count is zero")
	}
	return count, nil
}

func (w *countWatcher) getEvaluationCount(
	ctx context.Context,
	environmentNamespace, ruleID, clauseID, FeatureID, variationID string,
	featureVersion int32,
) (int64, error) {
	resp, err := w.eventCounterClient.GetOpsEvaluationUserCount(ctx, &ecproto.GetOpsEvaluationUserCountRequest{
		EnvironmentNamespace: environmentNamespace,
		OpsRuleId:            ruleID,
		ClauseId:             clauseID,
		FeatureId:            FeatureID,
		FeatureVersion:       featureVersion,
		VariationId:          variationID,
	})
	if err != nil {
		w.logger.Error("Failed to get ops evaluation count", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("ruleId", ruleID),
			zap.String("clauseId", clauseID),
			zap.String("featureId", FeatureID),
			zap.Int32("featureVersion", featureVersion),
			zap.String("variationId", variationID),
		)
		return 0, err
	}
	return resp.Count, nil
}

func (w *countWatcher) getTargetOpsGoalEventCount(
	ctx context.Context,
	logFunc func(string),
	environmentNamespace, ruleID, clauseID, FeatureID, variationID string,
	featureVersion int32,
) (int64, error) {
	count, err := w.getOpsGoalEventCount(
		ctx,
		environmentNamespace,
		ruleID,
		clauseID,
		FeatureID,
		variationID,
		featureVersion,
	)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		logFunc("Ops goal user count is zero")
	}
	return count, nil
}

func (w *countWatcher) getOpsGoalEventCount(
	ctx context.Context,
	environmentNamespace, ruleID, clauseID, FeatureID, variationID string,
	featureVersion int32,
) (int64, error) {
	resp, err := w.eventCounterClient.GetOpsGoalUserCount(ctx, &ecproto.GetOpsGoalUserCountRequest{
		EnvironmentNamespace: environmentNamespace,
		OpsRuleId:            ruleID,
		ClauseId:             clauseID,
		FeatureId:            FeatureID,
		FeatureVersion:       featureVersion,
		VariationId:          variationID,
	})
	if err != nil {
		w.logger.Error("Failed to get ops goal count", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("ruleId", ruleID),
			zap.String("clauseId", clauseID),
			zap.String("featureId", FeatureID),
			zap.Int32("featureVersion", featureVersion),
			zap.String("variationId", variationID),
		)
		return 0, err
	}
	return resp.Count, nil
}

func (w *countWatcher) persistOpsCount(
	ctx context.Context,
	environmentNamespace string,
	oc *opseventdomain.OpsCount,
) error {
	opsCountStorage := v2os.NewOpsCountStorage(w.mysqlClient)
	if err := opsCountStorage.UpsertOpsCount(ctx, environmentNamespace, oc); err != nil {
		w.logger.Error("Failed to upsert ops count", zap.Error(err),
			zap.String("autoOpsRuleId", oc.AutoOpsRuleId),
			zap.String("clauseId", oc.ClauseId),
			zap.String("environmentNamespace", environmentNamespace))
		return err
	}
	return nil
}
