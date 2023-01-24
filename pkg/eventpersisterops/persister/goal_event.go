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

type goalTargetRule struct {
	id        string
	featureID string
	clauses   map[string]*aoproto.OpsEventRateClause
}

type evalGoalUpdater struct {
	ctx               context.Context
	featureClient     featureclient.Client
	autoOpsClient     aoclient.Client
	eventCounterCache cachev3.EventCounterCache
	flightgroup       singleflight.Group
	logger            *zap.Logger
}

func NewGoalUserCountUpdater(
	ctx context.Context,
	featureClient featureclient.Client,
	autoOpsClient aoclient.Client,
	eventCounterCache cachev3.EventCounterCache,
	logger *zap.Logger,
) Updater {
	return &evalGoalUpdater{
		ctx:               ctx,
		featureClient:     featureClient,
		autoOpsClient:     autoOpsClient,
		eventCounterCache: eventCounterCache,
		logger:            logger,
	}
}

func (u *evalGoalUpdater) UpdateUserCounts(ctx context.Context, evt environmentEventMap) map[string]bool {
	fails := map[string]bool{}
	for environmentNamespace, events := range evt {
		for id, event := range events {
			switch evt := event.(type) {
			case *eventproto.GoalEvent:
				retriable, err := u.updateUserCount(ctx, environmentNamespace, evt)
				if err != nil {
					if err == ErrNoAutoOpsRules || err == ErrAutoOpsRulesNotFound {
						handledCounter.WithLabelValues(codeNoLink).Inc()
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
) (bool, error) {
	// List all the auto ops rules
	list, err := u.listAutoOpsRules(ctx, environmentNamespace)
	if err != nil {
		handledCounter.WithLabelValues(codeFailedToListAutoOpsRules).Inc()
		return true, err
	}
	if len(list) == 0 {
		return false, ErrNoAutoOpsRules
	}
	// Find the rules
	featureIDs, targetRules := u.findOpsRules(event.GoalId, list)
	if len(featureIDs) == 0 {
		return false, ErrAutoOpsRulesNotFound
	}
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
	for _, tr := range targetRules {
		// Find the latest feature version
		fVersion, err := u.findFeatureVersion(tr.featureID, resp.Features)
		if err != nil {
			u.logger.Error(
				"Failed to find the feature version",
				zap.Error(ErrFailedToFindFeatureVersion),
				zap.String("featureId", tr.featureID),
				zap.String("environmentNamespace", environmentNamespace),
			)
			handledCounter.WithLabelValues(codeFailedToFindFeatureVersion).Inc()
			return false, err
		}
		// Update the user count by rule
		err = u.updateUserCountByRule(
			environmentNamespace,
			tr.featureID,
			fVersion,
			event.UserId,
			tr,
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
			aor := []*aoproto.AutoOpsRule{}
			cursor := ""
			for {
				resp, err := u.autoOpsClient.ListAutoOpsRules(ctx, &aoproto.ListAutoOpsRulesRequest{
					EnvironmentNamespace: environmentNamespace,
					PageSize:             listRequestSize,
					Cursor:               cursor,
				})
				if err != nil {
					return nil, err
				}
				aor = append(aor, resp.AutoOpsRules...)
				aorSize := len(resp.AutoOpsRules)
				if aorSize == 0 || aorSize < listRequestSize {
					return aor, nil
				}
				cursor = resp.Cursor
			}
		},
	)
	if err != nil {
		return nil, err
	}
	return exp.([]*aoproto.AutoOpsRule), nil
}

func (u *evalGoalUpdater) findOpsRules(
	goalID string,
	listAutoOpsRules []*aoproto.AutoOpsRule,
) ([]string, []*goalTargetRule) {
	featureIDsMap := make(map[string]struct{})
	targetRules := []*goalTargetRule{}
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
		// Find the clauses that contain the goal ID from the the goal event
		targetClauses := make(map[string]*aoproto.OpsEventRateClause)
		for id, clause := range clauses {
			if clause.GoalId == goalID {
				featureIDsMap[autoOpsRule.FeatureId] = struct{}{}
				targetClauses[id] = clause
			}
		}
		if len(targetClauses) == 0 {
			continue
		}
		targetRules = append(targetRules, &goalTargetRule{
			id:        autoOpsRule.Id,
			featureID: autoOpsRule.FeatureId,
			clauses:   targetClauses,
		})
	}
	// Convert map to slice
	featureIDs := make([]string, 0, len(featureIDsMap))
	for id := range featureIDsMap {
		featureIDs = append(featureIDs, id)
	}
	return featureIDs, targetRules
}

func (u *evalGoalUpdater) findFeatureVersion(
	featureID string,
	features []*featureproto.Feature,
) (int32, error) {
	for _, f := range features {
		if f.Id == featureID {
			return f.Version, nil
		}
	}
	return 0, ErrFailedToFindFeatureVersion
}

func (u *evalGoalUpdater) updateUserCountByRule(
	environmentNamespace,
	featureID string,
	featureVersion int32,
	userID string,
	targetRule *goalTargetRule,
) error {
	for id, clause := range targetRule.clauses {
		key := u.newUserCountKey(
			environmentNamespace,
			targetRule.id,
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
	key := fmt.Sprintf("%s:%s:%s:%d:%s",
		ruleID,
		clauseID,
		featureID,
		featureVersion,
		variationID,
	)
	return cache.MakeKey(
		opsGoalKeyPrefix,
		key,
		environmentNamespace,
	)
}
