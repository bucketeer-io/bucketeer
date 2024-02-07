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
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const opsGoalKeyPrefix = "autoops:goal"

type linkGoalOpsRule struct {
	ruleID    string
	featureID string
	clauses   map[string]*aoproto.OpsEventRateClause
}

type evalGoalUpdater struct {
	ctx               context.Context
	featureClient     featureclient.Client
	autoOpsClient     aoclient.Client
	eventCounterCache cachev3.EventCounterCache
	autoOpsRulesCache cachev3.AutoOpsRulesCache
	flightgroup       singleflight.Group
	logger            *zap.Logger
}

func NewGoalUserCountUpdater(
	ctx context.Context,
	featureClient featureclient.Client,
	autoOpsClient aoclient.Client,
	eventCounterCache cachev3.EventCounterCache,
	autoOpsRulesCache cachev3.AutoOpsRulesCache,
	logger *zap.Logger,
) Updater {
	return &evalGoalUpdater{
		ctx:               ctx,
		featureClient:     featureClient,
		autoOpsClient:     autoOpsClient,
		eventCounterCache: eventCounterCache,
		autoOpsRulesCache: autoOpsRulesCache,
		logger:            logger,
	}
}

func (u *evalGoalUpdater) UpdateUserCounts(ctx context.Context, evt environmentEventMap) map[string]bool {
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
			case *eventproto.GoalEvent:
				retriable, err := u.updateUserCount(ctx, environmentNamespace, evt, listAutoOpsRules)
				if err != nil {
					if err == ErrAutoOpsRuleNotFound {
						// If there is nothing to link, we don't report it as an error
						handledCounter.WithLabelValues(codeAutoOpsRuleNotFound).Inc()
						u.logger.Debug(
							"There is no auto ops rules to link the goal event",
							zap.Error(err),
							zap.String("eventId", id),
							zap.String("environmentNamespace", environmentNamespace),
						)
						continue
					}
					if !retriable {
						u.logger.Error(
							"Failed to persister goal event for auto ops",
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
					"Unexpected message type while trying to persister a goal event",
					zap.String("eventId", id),
					zap.String("environmentNamespace", environmentNamespace),
				)
				fails[id] = false
			}
		}
	}
	return fails
}

func (u *evalGoalUpdater) updateUserCount(
	ctx context.Context,
	environmentNamespace string,
	event *eventproto.GoalEvent,
	listAutoOpsRules []*aoproto.AutoOpsRule,
) (bool, error) {
	// Link the rules
	linkedRules := u.linkOpsRulesByGoalID(event.GoalId, listAutoOpsRules)
	if len(linkedRules) == 0 {
		return false, ErrAutoOpsRuleNotFound
	}
	featureIDs := u.getUniqueFeatureIDs(linkedRules)
	// Get the latest feature version
	resp, err := u.featureClient.GetFeatures(ctx, &featureproto.GetFeaturesRequest{
		EnvironmentNamespace: environmentNamespace,
		Ids:                  featureIDs,
	})
	if err != nil {
		handledCounter.WithLabelValues(codeFailedToGetFeatures).Inc()
		return true, err
	}
	// At this point get features can't be empty
	if len(resp.Features) == 0 {
		handledCounter.WithLabelValues(codeGetFeaturesReturnedEmpty).Inc()
		return true, ErrFeatureEmptyList
	}
	for _, r := range linkedRules {
		// Get the latest feature version
		fVersion, err := u.getFeatureVersion(r.featureID, resp.Features)
		if err != nil {
			u.logger.Error(
				"Failed to find the feature version",
				zap.Error(ErrFeatureVersionNotFound),
				zap.String("featureId", r.featureID),
				zap.String("environmentNamespace", environmentNamespace),
			)
			handledCounter.WithLabelValues(codeFailedToFindFeatureVersion).Inc()
			return false, err
		}
		// Update the user count per clause
		err = u.updateUserCountPerClause(
			environmentNamespace,
			r.featureID,
			fVersion,
			event.UserId,
			r,
		)
		if err != nil {
			return true, err
		}
	}
	return false, nil
}

func (u *evalGoalUpdater) listAutoOpsRules(
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

func (u *evalGoalUpdater) linkOpsRulesByGoalID(
	goalID string,
	listAutoOpsRules []*aoproto.AutoOpsRule,
) []*linkGoalOpsRule {
	linkedRules := []*linkGoalOpsRule{}
	for _, aor := range listAutoOpsRules {
		autoOpsRule := &aodomain.AutoOpsRule{AutoOpsRule: aor}
		// We ignore the rules that are already triggered
		if autoOpsRule.AlreadyTriggered() {
			continue
		}
		clauses, err := autoOpsRule.ExtractOpsEventRateClauses()
		if err != nil {
			handledCounter.WithLabelValues(codeFailedToExtractOpsEventRateClauses).Inc()
			continue
		}
		if len(clauses) == 0 {
			continue
		}
		// Link the clauses that contain the goal ID from the the goal event
		linkedClauses := make(map[string]*aoproto.OpsEventRateClause)
		for id, clause := range clauses {
			if clause.GoalId == goalID {
				linkedClauses[id] = clause
			}
		}
		if len(linkedClauses) == 0 {
			continue
		}
		linkedRules = append(linkedRules, &linkGoalOpsRule{
			ruleID:    autoOpsRule.Id,
			featureID: autoOpsRule.FeatureId,
			clauses:   linkedClauses,
		})
	}
	return linkedRules
}

func (u *evalGoalUpdater) getUniqueFeatureIDs(rules []*linkGoalOpsRule) []string {
	ids := []string{}
	for _, rule := range rules {
		ids = append(ids, rule.featureID)
	}
	return ids
}

func (u *evalGoalUpdater) getFeatureVersion(
	featureID string,
	features []*featureproto.Feature,
) (int32, error) {
	for _, f := range features {
		if f.Id == featureID {
			return f.Version, nil
		}
	}
	return 0, ErrFeatureVersionNotFound
}

func (u *evalGoalUpdater) updateUserCountPerClause(
	environmentNamespace,
	featureID string,
	featureVersion int32,
	userID string,
	rule *linkGoalOpsRule,
) error {
	for id, clause := range rule.clauses {
		key := u.newUserCountKey(
			environmentNamespace,
			rule.ruleID,
			id,
			featureID,
			clause.VariationId,
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

func (u *evalGoalUpdater) newUserCountKey(
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
		opsGoalKeyPrefix,
		key,
		environmentNamespace,
	)
}
