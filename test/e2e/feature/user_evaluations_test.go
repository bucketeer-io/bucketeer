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

package feature

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestUpsertAndGetUserEvaluations(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := newFeatureClient(t)
	timestamp := time.Now().Unix()
	userID := fmt.Sprintf("e2e-test-ue-user-%s", newUUID(t))
	reqUpsert := &featureproto.UpsertUserEvaluationRequest{
		EnvironmentNamespace: *environmentNamespace,
		Tag:                  tags[0],
		Evaluation:           createEvaluation(t, userID, int32(timestamp)),
	}
	respUpsert, err := client.UpsertUserEvaluation(ctx, reqUpsert)
	assert.NoError(t, err)
	assert.NotNil(t, respUpsert)
	req := &featureproto.GetUserEvaluationsRequest{
		EnvironmentNamespace: *environmentNamespace,
		Tag:                  tags[0],
		UserId:               userID,
	}
	resp, err := client.GetUserEvaluations(ctx, req)
	assert.NoError(t, err)
	assert.True(t, containsEvaluation(t, reqUpsert.Evaluation, resp.Evaluations))
}

func createEvaluation(t *testing.T, userID string, featureVersion int32) *featureproto.Evaluation {
	t.Helper()
	return &featureproto.Evaluation{
		Id: domain.EvaluationID(
			"feature-id",
			featureVersion,
			userID,
		),
		FeatureId:      "feature-id",
		FeatureVersion: featureVersion,
		UserId:         userID,
		VariationId:    "variation-id",
		VariationValue: "variation-value",
		Reason: &featureproto.Reason{
			Type: featureproto.Reason_DEFAULT,
		},
	}
}

func containsEvaluation(
	t *testing.T,
	evaluation *featureproto.Evaluation,
	evaluations []*featureproto.Evaluation,
) bool {
	t.Helper()
	for _, e := range evaluations {
		if proto.Equal(e, evaluation) {
			return true
		}
	}
	return false
}
