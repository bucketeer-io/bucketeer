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

package autoops

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	autoopsclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const totalVariationWeight = int32(100000)

func TestCreateAndListProgressiveRollout(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createDisabledFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	schedules := createProgressiveRolloutSchedule()
	createProgressiveRollout(
		ctx,
		t,
		autoOpsClient,
		featureID,
		&autoopsproto.ProgressiveRolloutManualScheduleClause{
			Schedules:   schedules,
			VariationId: feature.Variations[0].Id,
		},
		nil,
	)
	progressiveRollout := listProgressiveRollouts(t, autoOpsClient, featureID)
	if len(progressiveRollout) != 1 {
		t.Fatal("not enough rules")
	}
	actual := progressiveRollout[0]
	if actual.Id == "" {
		t.Fatal("id is empty")
	}
	if actual.FeatureId != featureID {
		t.Fatalf("different feature ID, expected: %v, actual: %v", featureID, actual.FeatureId)
	}
	if actual.Type != autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE {
		t.Fatalf("different ops type, expected: %v, actual: %v", autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE, actual.Type)
	}
	actualClause := unmarshalProgressiveRolloutManualClause(t, actual.Clause)
	if actualClause.VariationId != feature.Variations[0].Id {
		t.Fatalf("different variation id, expected: %v, actual: %v", feature.Variations[0].Id, actualClause.VariationId)
	}
	if len(actualClause.Schedules) != len(schedules) {
		t.Fatalf("different length of schedules, expected: %v, actual: %v", len(actualClause.Schedules), len(schedules))
	}
}

func TestGetProgressiveRollout(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createDisabledFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	schedules := createProgressiveRolloutSchedule()
	createProgressiveRollout(
		ctx,
		t,
		autoOpsClient,
		featureID,
		&autoopsproto.ProgressiveRolloutManualScheduleClause{
			Schedules:   schedules,
			VariationId: feature.Variations[0].Id,
		},
		nil,
	)
	progressiveRollouts := listProgressiveRollouts(t, autoOpsClient, featureID)
	if len(progressiveRollouts) != 1 {
		t.Fatal("not enough rules")
	}
	actual := getProgressiveRollout(t, progressiveRollouts[0].Id)
	if actual.Id == "" {
		t.Fatal("id is empty")
	}
	if actual.FeatureId != featureID {
		t.Fatalf("different feature ID, expected: %v, actual: %v", featureID, actual.FeatureId)
	}
	if actual.Type != autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE {
		t.Fatalf("different ops type, expected: %v, actual: %v", autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE, actual.Type)
	}
	actualClause := unmarshalProgressiveRolloutManualClause(t, actual.Clause)
	if actualClause.VariationId != feature.Variations[0].Id {
		t.Fatalf("different variation id, expected: %v, actual: %v", feature.Variations[0].Id, actualClause.VariationId)
	}
	if len(actualClause.Schedules) != len(schedules) {
		t.Fatalf("different length of schedules, expected: %v, actual: %v", len(actualClause.Schedules), len(schedules))
	}
}

func TestStopProgressiveRollout(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createDisabledFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	schedules := createProgressiveRolloutSchedule()
	createProgressiveRollout(
		ctx,
		t,
		autoOpsClient,
		featureID,
		&autoopsproto.ProgressiveRolloutManualScheduleClause{
			Schedules:   schedules,
			VariationId: feature.Variations[0].Id,
		},
		nil,
	)
	progressiveRollouts := listProgressiveRollouts(t, autoOpsClient, featureID)
	if len(progressiveRollouts) != 1 {
		t.Fatal("not enough rules")
	}
	stopProgressiveRollout(t, autoOpsClient, progressiveRollouts[0].Id)
	resp, err := autoOpsClient.GetProgressiveRollout(ctx, &autoopsproto.GetProgressiveRolloutRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   progressiveRollouts[0].Id,
	})
	assert.NoError(t, err)
	assert.True(t, resp.ProgressiveRollout.StoppedAt > time.Now().Add(time.Second*-10).Unix())
	assert.Equal(t, autoopsproto.ProgressiveRollout_STOPPED, resp.ProgressiveRollout.Status)
	assert.Equal(t, autoopsproto.ProgressiveRollout_USER, resp.ProgressiveRollout.StoppedBy)
}

func TestDeleteProgressiveRollout(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createDisabledFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	schedules := createProgressiveRolloutSchedule()
	createProgressiveRollout(
		ctx,
		t,
		autoOpsClient,
		featureID,
		&autoopsproto.ProgressiveRolloutManualScheduleClause{
			Schedules:   schedules,
			VariationId: feature.Variations[0].Id,
		},
		nil,
	)
	progressiveRollouts := listProgressiveRollouts(t, autoOpsClient, featureID)
	if len(progressiveRollouts) != 1 {
		t.Fatal("not enough rules")
	}
	deleteProgressiveRollout(t, autoOpsClient, progressiveRollouts[0].Id)
	resp, err := autoOpsClient.GetProgressiveRollout(ctx, &autoopsproto.GetProgressiveRolloutRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   progressiveRollouts[0].Id,
	})
	if resp != nil {
		t.Fatal("progressiveRollout is not deleted")
	}
	if err == nil {
		t.Fatal("err is empty")
	}
	if status.Code(err) != codes.NotFound {
		t.Fatalf("different error code, expected: %s, actual: %s", codes.NotFound, status.Code(err))
	}
}

func TestExecuteProgressiveRollout(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createDisabledFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	schedules := createProgressiveRolloutSchedule()
	createProgressiveRollout(
		ctx,
		t,
		autoOpsClient,
		featureID,
		&autoopsproto.ProgressiveRolloutManualScheduleClause{
			Schedules:   schedules,
			VariationId: feature.Variations[0].Id,
		},
		nil,
	)
	progressiveRollouts := listProgressiveRollouts(t, autoOpsClient, featureID)
	if len(progressiveRollouts) != 1 {
		t.Fatal("not enough rules")
	}
	clause := unmarshalProgressiveRolloutManualClause(t, progressiveRollouts[0].Clause)
	_, err := autoOpsClient.ExecuteProgressiveRollout(ctx, &autoopsproto.ExecuteProgressiveRolloutRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   progressiveRollouts[0].Id,
		ChangeProgressiveRolloutTriggeredAtCommand: &autoopsproto.ChangeProgressiveRolloutScheduleTriggeredAtCommand{
			ScheduleId: clause.Schedules[0].ScheduleId,
		},
	})
	if err != nil {
		t.Fatalf("failed to execute progressive rollout: %s", err.Error())
	}
	feature = getFeature(t, featureClient, featureID)
	expectedStrategy := &featureproto.RolloutStrategy{
		Variations: []*featureproto.RolloutStrategy_Variation{
			{
				Variation: feature.Variations[0].Id,
				Weight:    schedules[0].Weight,
			},
			{
				Variation: feature.Variations[1].Id,
				Weight:    totalVariationWeight - schedules[0].Weight,
			},
		},
	}
	if !feature.Enabled {
		t.Fatalf("Flag shouldn't be disabled at this point")
	}
	if !proto.Equal(feature.DefaultStrategy.RolloutStrategy, expectedStrategy) {
		t.Fatalf("Strategy is not equal. Expected: %s actual: %s", expectedStrategy, feature.Rules[0].Strategy.RolloutStrategy)
	}
	actual := listProgressiveRollouts(t, autoOpsClient, featureID)
	if actual[0].Status != autoopsproto.ProgressiveRollout_RUNNING {
		t.Fatalf("different status, expected: %v, actual: %v", actual[0].Status, autoopsproto.ProgressiveRollout_RUNNING)
	}
	actualClause := unmarshalProgressiveRolloutManualClause(t, actual[0].Clause)
	if actualClause.VariationId != feature.Variations[0].Id {
		t.Fatalf("different variation id, expected: %v, actual: %v", feature.Variations[0].Id, actualClause.VariationId)
	}
	if len(actualClause.Schedules) != len(schedules) {
		t.Fatalf("different length of schedules, expected: %v, actual: %v", len(actualClause.Schedules), len(schedules))
	}
	if actualClause.Schedules[0].TriggeredAt == 0 {
		t.Fatalf("triggered at is empty")
	}
}

func TestProgressiveRolloutBatch(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	autoOpsClient := newAutoOpsClient(t)
	defer autoOpsClient.Close()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	featureID := createFeatureID(t)
	createDisabledFeature(ctx, t, featureClient, featureID)
	feature := getFeature(t, featureClient, featureID)
	now := time.Now()
	schedules := []*autoopsproto.ProgressiveRolloutSchedule{
		{
			Weight:    20000,
			ExecuteAt: now.Add(5 * time.Second).Unix(),
		},
	}
	createProgressiveRollout(
		ctx,
		t,
		autoOpsClient,
		featureID,
		&autoopsproto.ProgressiveRolloutManualScheduleClause{
			Schedules:   schedules,
			VariationId: feature.Variations[0].Id,
		},
		nil,
	)
	progressiveRollouts := listProgressiveRollouts(t, autoOpsClient, featureID)
	if len(progressiveRollouts) != 1 {
		t.Fatal("not enough rules")
	}

	maxRetryCount := 18 // 3 minutes
	for i := 0; i < maxRetryCount-1; i++ {
		if i >= maxRetryCount {
			t.Fatalf("Retry count has reached the limit")
		}
		time.Sleep(10 * time.Second)
		feature = getFeature(t, featureClient, featureID)
		expectedStrategy := &featureproto.RolloutStrategy{
			Variations: []*featureproto.RolloutStrategy_Variation{
				{
					Variation: feature.Variations[0].Id,
					Weight:    schedules[0].Weight,
				},
				{
					Variation: feature.Variations[1].Id,
					Weight:    totalVariationWeight - schedules[0].Weight,
				},
			},
		}
		if !proto.Equal(feature.DefaultStrategy.RolloutStrategy, expectedStrategy) {
			continue
		}
		if !feature.Enabled {
			t.Fatalf("Flag shouldn't be disabled at this point")
		}
		actual := listProgressiveRollouts(t, autoOpsClient, featureID)
		if actual[0].Status != autoopsproto.ProgressiveRollout_FINISHED {
			t.Fatalf("different status, expected: %v, actual: %v", actual[0].Status, autoopsproto.ProgressiveRollout_FINISHED)
		}
		actualClause := unmarshalProgressiveRolloutManualClause(t, actual[0].Clause)
		if actualClause.VariationId != feature.Variations[0].Id {
			t.Fatalf("different variation id, expected: %v, actual: %v", feature.Variations[0].Id, actualClause.VariationId)
		}
		if actualClause.Schedules[0].TriggeredAt == 0 {
			t.Fatalf("triggered at is empty")
		}
		break
	}
}

func createProgressiveRollout(
	ctx context.Context,
	t *testing.T,
	client autoopsclient.Client,
	featureID string,
	manual *autoopsproto.ProgressiveRolloutManualScheduleClause,
	template *autoopsproto.ProgressiveRolloutTemplateScheduleClause,
) {
	t.Helper()
	cmd := &autoopsproto.CreateProgressiveRolloutCommand{
		FeatureId:                                featureID,
		ProgressiveRolloutManualScheduleClause:   manual,
		ProgressiveRolloutTemplateScheduleClause: template,
	}
	_, err := client.CreateProgressiveRollout(ctx, &autoopsproto.CreateProgressiveRolloutRequest{
		EnvironmentNamespace: *environmentNamespace,
		Command:              cmd,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func listProgressiveRollouts(t *testing.T, client autoopsclient.Client, featureID string) []*autoopsproto.ProgressiveRollout {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.ListProgressiveRollouts(ctx, &autoopsproto.ListProgressiveRolloutsRequest{
		EnvironmentNamespace: *environmentNamespace,
		PageSize:             0,
		FeatureIds:           []string{featureID},
	})
	if err != nil {
		t.Fatal("Failed to list progressive rollout", err)
	}
	return resp.ProgressiveRollouts
}

func unmarshalProgressiveRolloutManualClause(t *testing.T, clause *anypb.Any) *autoopsproto.ProgressiveRolloutManualScheduleClause {
	c := &autoopsproto.ProgressiveRolloutManualScheduleClause{}
	if err := ptypes.UnmarshalAny(clause, c); err != nil {
		t.Fatal(err)
	}
	return c
}

func createProgressiveRolloutSchedule() []*autoopsproto.ProgressiveRolloutSchedule {
	now := time.Now()
	return []*autoopsproto.ProgressiveRolloutSchedule{
		{
			Weight:    20000,
			ExecuteAt: now.AddDate(0, 0, 3).Unix(),
		},
		{
			Weight:    40000,
			ExecuteAt: now.AddDate(0, 0, 6).Unix(),
		},
		{
			Weight:    60000,
			ExecuteAt: now.AddDate(0, 0, 9).Unix(),
		},
	}
}

func getProgressiveRollout(t *testing.T, id string) *autoopsproto.ProgressiveRollout {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newAutoOpsClient(t)
	defer c.Close()
	resp, err := c.GetProgressiveRollout(ctx, &autoopsproto.GetProgressiveRolloutRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
	})
	if err != nil {
		t.Fatal("Failed to get progressive rollout", err)
	}
	return resp.ProgressiveRollout
}

func stopProgressiveRollout(t *testing.T, client autoopsclient.Client, id string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.StopProgressiveRollout(ctx, &autoopsproto.StopProgressiveRolloutRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
		Command: &autoopsproto.StopProgressiveRolloutCommand{
			StoppedBy: autoopsproto.ProgressiveRollout_USER,
		},
	})
	if err != nil {
		t.Fatal("Failed to stop progressive rollout", err)
	}
}

func deleteProgressiveRollout(t *testing.T, client autoopsclient.Client, id string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.DeleteProgressiveRollout(ctx, &autoopsproto.DeleteProgressiveRolloutRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
		Command:              &autoopsproto.DeleteProgressiveRolloutCommand{},
	})
	if err != nil {
		t.Fatal("Failed to delete progressive rollout", err)
	}
}
