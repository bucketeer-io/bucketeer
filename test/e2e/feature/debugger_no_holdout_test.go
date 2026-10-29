// Copyright 2026 The Bucketeer Authors.
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

package feature

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

// The debugger evaluates from the feature cache, which the FeatureFlagCacher
// cron refreshes every minute; one TTL of slack is added.
const debuggerCacheTimeout = 150 * time.Second

// Regression guard for RFC 0048 §4.8: the console debugger shares the
// evaluator with the gateway, and the shared e2e environment never has a
// holdout group, so no evaluation it returns may carry Reason.HOLDOUT.
// TestEvaluateFeatures and TestDebugEvaluateFeatures cover the shape of the
// responses; this test only pins the reason.
func TestDebuggerNoHoldoutReason(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	featureID := newFeatureID(t)
	createFeature(t, client, newCreateFeatureReq(featureID))
	enableFeature(t, featureID, client)
	userID := "user-id-01"

	var evals *feature.EvaluateFeaturesResponse
	require.Eventually(t, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		var err error
		evals, err = client.EvaluateFeatures(ctx, &feature.EvaluateFeaturesRequest{
			User:          &userproto.User{Id: userID},
			EnvironmentId: *environmentID,
			Tag:           tags[0],
		})
		return err == nil && hasEvaluationOf(evals.UserEvaluations.Evaluations, featureID)
	}, debuggerCacheTimeout, 2*time.Second, "EvaluateFeatures never returned the new flag")
	for _, e := range evals.UserEvaluations.Evaluations {
		assert.NotEqual(t, feature.Reason_HOLDOUT, e.GetReason().GetType(), "EvaluateFeatures %s", e.FeatureId)
		assert.Empty(t, e.GetReason().GetHoldoutGroupId())
	}

	var debug *feature.DebugEvaluateFeaturesResponse
	require.Eventually(t, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		var err error
		debug, err = client.DebugEvaluateFeatures(ctx, &feature.DebugEvaluateFeaturesRequest{
			Users:         []*userproto.User{{Id: userID}},
			FeatureIds:    []string{featureID},
			EnvironmentId: *environmentID,
		})
		return err == nil && hasEvaluationOf(debug.Evaluations, featureID)
	}, debuggerCacheTimeout, 2*time.Second, "DebugEvaluateFeatures never returned the new flag")
	for _, e := range debug.Evaluations {
		assert.NotEqual(t, feature.Reason_HOLDOUT, e.GetReason().GetType(), "DebugEvaluateFeatures %s", e.FeatureId)
		assert.Empty(t, e.GetReason().GetHoldoutGroupId())
	}
}

func hasEvaluationOf(evaluations []*feature.Evaluation, featureID string) bool {
	for _, e := range evaluations {
		if e.FeatureId == featureID {
			return true
		}
	}
	return false
}
