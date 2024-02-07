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

package persister

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	aoclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	aodomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	aoproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
)

const opsEvalKeyPrefix = "autoops:evaluation"

type evalEvtUpdater struct {
	ctx               context.Context
	featureClient     featureclient.Client
	autoOpsClient     aoclient.Client
	eventCounterCache cachev3.EventCounterCache
	autoOpsRulesCache cachev3.AutoOpsRulesCache
	flightgroup       singleflight.Group
	logger            *zap.Logger
}

func NewEvalUserCountUpdater(
	ctx context.Context,
	featureClient featureclient.Client,
	autoOpsClient aoclient.Client,
	eventCounterCache cachev3.EventCounterCache,
	autoOpsRulesCache cachev3.AutoOpsRulesCache,
	logger *zap.Logger,
) Updater {
	return &evalEvtUpdater{
		ctx:               ctx,
		featureClient:     featureClient,
		autoOpsClient:     autoOpsClient,
		eventCounterCache: eventCounterCache,
		autoOpsRulesCache: autoOpsRulesCache,
		logger:            logger,
	}
}

func (u *evalEvtUpdater) UpdateUserCounts(ctx context.Context, evt environmentEventMap) map[string]bool {
	fails := map[string]bool{}
	for environmentNamespace, events := range evt {
		listAutoOpsRules, err := u.listAutoOpsRules(ctx, environmentNamespace)
		if err != nil {
			u.logger.Error("failed to list auto ops rules",
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
			)
			handledCounter.WithLabelValues(codeFailedToListAutoOpsRules).Inc()
			// Make sure to retry all the events in the next pulling
			for id := range events {
				fails[id] = true
			}
			continue
		}
		if len(listAutoOpsRules) == 0 {
			continue
		}
		for id, event := range events {
			switch evt := event.(type) {
			case *eventproto.EvaluationEvent:
				retriable, err := u.updateUserCount(ctx, environmentNamespace, evt, listAutoOpsRules)
				if err != nil {
					if err == ErrAutoOpsRuleNotFound {
						// If there is nothing to link, we don't report it as an error
						handledCounter.WithLabelValues(codeAutoOpsRuleNotFound).Inc()
						continue
					}
					if !retriable {
						u.logger.Error(
							"Failed to persister evaluation event for auto ops",
							zap.Error(err),
							zap.String("eventId", id),
							zap.String("environmentNamespace", environmentNamespace),
						)
					}
					fails[id] = retriable
					continue
				}
				handledCounter.WithLabelValues(codeLinked).Inc()
			default:
				u.logger.Error(
					"Unexpected message type while trying to persister an evaluation event",
					zap.String("eventId", id),
					zap.String("environmentNamespace", environmentNamespace),
				)
				fails[id] = false
			}
		}
	}
	return fails
}

func (u *evalEvtUpdater) updateUserCount(
	ctx context.Context,
	environmentNamespace string,
	event *eventproto.EvaluationEvent,
	listAutoOpsRules []*aoproto.AutoOpsRule,
) (bool, error) {
	rules := u.linkOpsRulesByFeatureID(event.FeatureId, listAutoOpsRules)
	if len(rules) == 0 {
		return false, ErrAutoOpsRuleNotFound
	}
	// Link the event rate clauses by variation ID
	linkedOpsRules := make(map[string][]string, len(rules))
	for _, rule := range rules {
		clauseIDs, err := u.linkOpsEventRateByVariationID(rule, event.VariationId)
		if err != nil {
			return false, err
		}
		linkedOpsRules[rule.Id] = clauseIDs
	}
	// Update the user count per rule
	for ruleID, clauseIDs := range linkedOpsRules {
		err := u.updateUserCountPerClause(
			environmentNamespace,
			event.FeatureId,
			event.FeatureVersion,
			event.VariationId,
			event.UserId,
			ruleID,
			clauseIDs,
		)
		if err != nil {
			return true, err
		}
	}
	return false, nil
}

func (u *evalEvtUpdater) listAutoOpsRules(
	ctx context.Context,
	environmentNamespace string,
) ([]*aoproto.AutoOpsRule, error) {
	exp, err, _ := u.flightgroup.Do(
		fmt.Sprintf("%s:%s", environmentNamespace, "listAutoOpsRules"),
		func() (interface{}, error) {
			// Get the auto ops rules cache
			aorList, err := u.autoOpsRulesCache.Get(environmentNamespace)
			if err == nil {
				return aorList.AutoOpsRules, nil
			}
			// We don't use the feature ID to filter the results in the request
			// because it will increase access to the DB, which also will increase the costs.
			// So we list all rules and use the singleflight implementation to share the response
			resp, err := u.autoOpsClient.ListAutoOpsRules(ctx, &aoproto.ListAutoOpsRulesRequest{
				EnvironmentNamespace: environmentNamespace,
				PageSize:             0,
			})
			if err != nil {
				return nil, err
			}
			// Cache the auto ops rules for the next request
			autoOpsRules := &aoproto.AutoOpsRules{
				AutoOpsRules: resp.AutoOpsRules,
			}
			if err := u.autoOpsRulesCache.Put(autoOpsRules, environmentNamespace); err != nil {
				u.logger.Error("Failed to cache auto ops rules",
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)
			}
			return resp.AutoOpsRules, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return exp.([]*aoproto.AutoOpsRule), nil
}

func (u *evalEvtUpdater) linkOpsRulesByFeatureID(
	featureID string,
	listAutoOpsRules []*aoproto.AutoOpsRule,
) []*aoproto.AutoOpsRule {
	rules := []*aoproto.AutoOpsRule{}
	for _, aor := range listAutoOpsRules {
		// Ignore already triggered ops rules
		if aor.FeatureId == featureID && aor.TriggeredAt == 0 {
			rules = append(rules, aor)
		}
	}
	return rules
}

func (u *evalEvtUpdater) linkOpsEventRateByVariationID(
	rule *aoproto.AutoOpsRule,
	variationID string,
) ([]string, error) {
	r := &aodomain.AutoOpsRule{AutoOpsRule: rule}
	clauses, err := r.ExtractOpsEventRateClauses()
	if err != nil {
		handledCounter.WithLabelValues(codeFailedToExtractOpsEventRateClauses).Inc()
		return nil, err
	}
	ids := make([]string, 0, len(clauses))
	for id, clause := range clauses {
		// The variation must match to link
		if clause.VariationId == variationID {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func (u *evalEvtUpdater) updateUserCountPerClause(
	environmentNamespace,
	featureID string,
	featureVersion int32,
	variationID,
	userID,
	ruleID string,
	clauseIDs []string,
) error {
	for _, clauseID := range clauseIDs {
		key := u.newUserCountKey(
			environmentNamespace,
			ruleID,
			clauseID,
			featureID,
			variationID,
			featureVersion,
		)
		if err := u.eventCounterCache.UpdateUserCount(key, userID); err != nil {
			handledCounter.WithLabelValues(codeFailedToUpdateUserCount).Inc()
			return err
		}
		u.logger.Debug(
			"User count updated successfully",
			zap.String("pfcountKey", key),
			zap.String("environmentNamespace", environmentNamespace),
		)
	}
	return nil
}

func (u *evalEvtUpdater) newUserCountKey(
	environmentNamespace,
	ruleID, clauseID, featureID, variationID string,
	featureVersion int32,
) string {
	key := fmt.Sprintf("%s:%d:%s:%s:%s",
		featureID,
		featureVersion,
		ruleID,
		clauseID,
		variationID,
	)
	return cache.MakeKey(
		opsEvalKeyPrefix,
		key,
		environmentNamespace,
	)
}
