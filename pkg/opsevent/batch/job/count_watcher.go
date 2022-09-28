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

const (
	queryTimeRange = -30 * 24 * time.Hour
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
		evaluationCount, err := w.getTargetEvaluationCount(ctx,
			logFunc,
			env.Namespace,
			a.FeatureId,
			c.VariationId,
			featureVersion,
		)
		if err != nil {
			lastErr = err
			continue
		}
		if evaluationCount == nil {
			continue
		}
		opsEventCount, err := w.getTargetOpsEventCount(
			ctx,
			logFunc,
			env.Namespace,
			a.FeatureId,
			c.VariationId,
			c.GoalId,
			featureVersion,
		)
		if err != nil {
			lastErr = err
			continue
		}
		if opsEventCount == nil {
			continue
		}
		opsCount := opseventdomain.NewOpsCount(a.FeatureId, a.Id, id, opsEventCount.UserCount, evaluationCount.UserCount)
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
	evaluationCount,
	opsCount *ecproto.VariationCount,
) bool {
	rate := float64(opsCount.UserCount) / float64(evaluationCount.UserCount)
	if opsCount.UserCount < opsEventRateClause.MinCount {
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

func (w *countWatcher) getTargetEvaluationCount(
	ctx context.Context,
	logFunc func(string),
	environmentNamespace, FeatureID, variationID string,
	featureVersion int32,
) (*ecproto.VariationCount, error) {
	evaluationCount, err := w.getEvaluationCount(
		ctx,
		environmentNamespace,
		FeatureID,
		variationID,
		featureVersion,
	)
	if err != nil {
		return nil, err
	}
	if evaluationCount == nil {
		logFunc("evaluationCount is nil")
		return nil, nil
	}
	if evaluationCount.UserCount == 0 {
		logFunc("evaluationCount.UserCount is zero")
		return nil, nil
	}
	return evaluationCount, nil
}

func (w *countWatcher) getEvaluationCount(
	ctx context.Context,
	environmentNamespace, FeatureID, variationID string,
	featureVersion int32,
) (*ecproto.VariationCount, error) {
	endAt := time.Now()
	startAt := endAt.Add(queryTimeRange)
	resp, err := w.eventCounterClient.GetEvaluationCountV2(ctx, &ecproto.GetEvaluationCountV2Request{
		EnvironmentNamespace: environmentNamespace,
		StartAt:              startAt.Unix(),
		EndAt:                endAt.Unix(),
		FeatureId:            FeatureID,
		FeatureVersion:       featureVersion,
		VariationIds:         []string{variationID},
	})
	if err != nil {
		w.logger.Error("Failed to get evaluation realtime count", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("featureId", FeatureID),
			zap.Int32("featureVersion", featureVersion),
			zap.String("variationId", variationID),
		)
		return nil, err
	}
	if len(resp.Count.RealtimeCounts) == 0 {
		return nil, nil
	}
	for _, vc := range resp.Count.RealtimeCounts {
		if vc.VariationId == variationID {
			return vc, nil
		}
	}
	return nil, nil
}

func (w *countWatcher) getTargetOpsEventCount(
	ctx context.Context,
	logFunc func(string),
	environmentNamespace, FeatureID, variationID, goalID string,
	featureVersion int32,
) (*ecproto.VariationCount, error) {
	opsCount, err := w.getOpsEventCount(
		ctx,
		environmentNamespace,
		FeatureID,
		variationID,
		goalID,
		featureVersion,
	)
	if err != nil {
		return nil, err
	}
	if opsCount == nil {
		logFunc("opsCount is nil")
		return nil, nil
	}
	return opsCount, nil
}

func (w *countWatcher) getOpsEventCount(
	ctx context.Context,
	environmentNamespace, FeatureID, variationID, goalID string,
	featureVersion int32,
) (*ecproto.VariationCount, error) {
	endAt := time.Now()
	startAt := endAt.Add(queryTimeRange)
	resp, err := w.eventCounterClient.GetGoalCountV2(ctx, &ecproto.GetGoalCountV2Request{
		EnvironmentNamespace: environmentNamespace,
		StartAt:              startAt.Unix(),
		EndAt:                endAt.Unix(),
		FeatureId:            FeatureID,
		FeatureVersion:       featureVersion,
		VariationIds:         []string{variationID},
		GoalId:               goalID,
	})
	if err != nil {
		w.logger.Error("Failed to get ops realtime variation count", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("featureId", FeatureID),
			zap.Int32("featureVersion", featureVersion),
			zap.String("variationId", variationID),
			zap.String("goalId", goalID),
		)
		return nil, err
	}
	for _, vc := range resp.GoalCounts.RealtimeCounts {
		if vc.VariationId == variationID {
			return vc, nil
		}
	}
	return nil, nil
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
