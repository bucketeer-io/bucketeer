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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestCreateScheduledFlagChange(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	// Create a feature first
	featureID := newFeatureID(t)
	createFeatureNoCmd(t, client, newCreateFeatureReq(featureID))

	// Create scheduled flag change
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	scheduledAt := time.Now().Add(1 * time.Hour).Unix()
	req := &featureproto.CreateScheduledFlagChangeRequest{
		EnvironmentId: *environmentID,
		FeatureId:     featureID,
		ScheduledAt:   scheduledAt,
		Timezone:      "UTC",
		Payload: &featureproto.ScheduledChangePayload{
			Enabled: wrapperspb.Bool(true),
		},
		Comment: "e2e test - enable flag in 1 hour",
	}

	resp, err := client.CreateScheduledFlagChange(ctx, req)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.ScheduledFlagChange.Id)
	assert.Equal(t, featureID, resp.ScheduledFlagChange.FeatureId)
	assert.Equal(t, *environmentID, resp.ScheduledFlagChange.EnvironmentId)
	assert.Equal(t, scheduledAt, resp.ScheduledFlagChange.ScheduledAt)
	assert.Equal(t, "UTC", resp.ScheduledFlagChange.Timezone)
	assert.Equal(t, featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING, resp.ScheduledFlagChange.Status)
	assert.NotEmpty(t, resp.ScheduledFlagChange.ChangeSummaries)
	assert.Equal(t, "ScheduledChange.EnableFlag", resp.ScheduledFlagChange.ChangeSummaries[0].MessageKey)
}

func TestGetScheduledFlagChange(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	// Create a feature first
	featureID := newFeatureID(t)
	createFeatureNoCmd(t, client, newCreateFeatureReq(featureID))

	// Create scheduled flag change
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	createResp := createScheduledFlagChange(t, client, featureID, &featureproto.ScheduledChangePayload{
		Enabled: wrapperspb.Bool(true),
	})

	// Get the scheduled flag change
	getReq := &featureproto.GetScheduledFlagChangeRequest{
		EnvironmentId: *environmentID,
		Id:            createResp.ScheduledFlagChange.Id,
	}
	getResp, err := client.GetScheduledFlagChange(ctx, getReq)
	require.NoError(t, err)
	assert.Equal(t, createResp.ScheduledFlagChange.Id, getResp.ScheduledFlagChange.Id)
	assert.Equal(t, featureID, getResp.ScheduledFlagChange.FeatureId)
	assert.Equal(t, featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING, getResp.ScheduledFlagChange.Status)
}

func TestGetScheduledFlagChangeNotFound(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	getReq := &featureproto.GetScheduledFlagChangeRequest{
		EnvironmentId: *environmentID,
		Id:            "non-existent-id",
	}
	_, err := client.GetScheduledFlagChange(ctx, getReq)
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestUpdateScheduledFlagChange(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	// Create a feature first
	featureID := newFeatureID(t)
	createFeatureNoCmd(t, client, newCreateFeatureReq(featureID))

	// Create scheduled flag change
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	createResp := createScheduledFlagChange(t, client, featureID, &featureproto.ScheduledChangePayload{
		Enabled: wrapperspb.Bool(true),
	})

	// Update the scheduled flag change
	newScheduledAt := time.Now().Add(2 * time.Hour).Unix()
	updateReq := &featureproto.UpdateScheduledFlagChangeRequest{
		EnvironmentId: *environmentID,
		Id:            createResp.ScheduledFlagChange.Id,
		ScheduledAt:   wrapperspb.Int64(newScheduledAt),
		Payload: &featureproto.ScheduledChangePayload{
			Enabled: wrapperspb.Bool(false),
		},
		Comment: wrapperspb.String("e2e test - updated to disable flag"),
	}
	_, err := client.UpdateScheduledFlagChange(ctx, updateReq)
	require.NoError(t, err)

	// Verify the update
	getResp := getScheduledFlagChange(t, client, createResp.ScheduledFlagChange.Id)
	assert.Equal(t, newScheduledAt, getResp.ScheduledFlagChange.ScheduledAt)
	assert.Equal(t, "ScheduledChange.DisableFlag", getResp.ScheduledFlagChange.ChangeSummaries[0].MessageKey)
}

func TestDeleteScheduledFlagChange(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	// Create a feature first
	featureID := newFeatureID(t)
	createFeatureNoCmd(t, client, newCreateFeatureReq(featureID))

	// Create scheduled flag change
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	createResp := createScheduledFlagChange(t, client, featureID, &featureproto.ScheduledChangePayload{
		Enabled: wrapperspb.Bool(true),
	})

	// Delete (cancel) the scheduled flag change
	deleteReq := &featureproto.DeleteScheduledFlagChangeRequest{
		EnvironmentId: *environmentID,
		Id:            createResp.ScheduledFlagChange.Id,
	}
	_, err := client.DeleteScheduledFlagChange(ctx, deleteReq)
	require.NoError(t, err)

	// Verify it's cancelled (not deleted, soft delete)
	getResp := getScheduledFlagChange(t, client, createResp.ScheduledFlagChange.Id)
	assert.Equal(t, featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CANCELLED, getResp.ScheduledFlagChange.Status)
}

func TestListScheduledFlagChanges(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	// Create a feature first
	featureID := newFeatureID(t)
	createFeatureNoCmd(t, client, newCreateFeatureReq(featureID))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create multiple scheduled flag changes
	createResp1 := createScheduledFlagChange(t, client, featureID, &featureproto.ScheduledChangePayload{
		Enabled: wrapperspb.Bool(true),
	})
	time.Sleep(100 * time.Millisecond) // Ensure different created_at times
	createResp2 := createScheduledFlagChange(t, client, featureID, &featureproto.ScheduledChangePayload{
		Enabled: wrapperspb.Bool(false),
	})

	// List scheduled flag changes
	listReq := &featureproto.ListScheduledFlagChangesRequest{
		EnvironmentId: *environmentID,
		FeatureId:     featureID,
		PageSize:      10,
	}
	listResp, err := client.ListScheduledFlagChanges(ctx, listReq)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(listResp.ScheduledFlagChanges), 2)

	// Verify the created schedules are in the list
	foundIDs := make(map[string]bool)
	for _, sfc := range listResp.ScheduledFlagChanges {
		foundIDs[sfc.Id] = true
	}
	assert.True(t, foundIDs[createResp1.ScheduledFlagChange.Id], "First schedule not found")
	assert.True(t, foundIDs[createResp2.ScheduledFlagChange.Id], "Second schedule not found")
}

func TestListScheduledFlagChangesWithStatusFilter(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	// Create a feature first
	featureID := newFeatureID(t)
	createFeatureNoCmd(t, client, newCreateFeatureReq(featureID))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a scheduled flag change
	createResp := createScheduledFlagChange(t, client, featureID, &featureproto.ScheduledChangePayload{
		Enabled: wrapperspb.Bool(true),
	})

	// Cancel it
	_, err := client.DeleteScheduledFlagChange(ctx, &featureproto.DeleteScheduledFlagChangeRequest{
		EnvironmentId: *environmentID,
		Id:            createResp.ScheduledFlagChange.Id,
	})
	require.NoError(t, err)

	// List only pending schedules - should not include cancelled
	listReq := &featureproto.ListScheduledFlagChangesRequest{
		EnvironmentId: *environmentID,
		FeatureId:     featureID,
		Statuses:      []featureproto.ScheduledFlagChangeStatus{featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING},
		PageSize:      10,
	}
	listResp, err := client.ListScheduledFlagChanges(ctx, listReq)
	require.NoError(t, err)

	// Verify the cancelled schedule is not in the pending list
	for _, sfc := range listResp.ScheduledFlagChanges {
		assert.NotEqual(t, createResp.ScheduledFlagChange.Id, sfc.Id)
	}
}

func TestGetScheduledFlagChangeSummary(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	// Create a feature first
	featureID := newFeatureID(t)
	createFeatureNoCmd(t, client, newCreateFeatureReq(featureID))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create scheduled flag changes
	createScheduledFlagChange(t, client, featureID, &featureproto.ScheduledChangePayload{
		Enabled: wrapperspb.Bool(true),
	})
	createScheduledFlagChange(t, client, featureID, &featureproto.ScheduledChangePayload{
		Enabled: wrapperspb.Bool(false),
	})

	// Get summary
	summaryReq := &featureproto.GetScheduledFlagChangeSummaryRequest{
		EnvironmentId: *environmentID,
		FeatureId:     featureID,
	}
	summaryResp, err := client.GetScheduledFlagChangeSummary(ctx, summaryReq)
	require.NoError(t, err)
	assert.NotNil(t, summaryResp.Summary)
	assert.Equal(t, featureID, summaryResp.Summary.FeatureId)
	assert.GreaterOrEqual(t, summaryResp.Summary.PendingCount, int32(2))
	assert.NotZero(t, summaryResp.Summary.NextScheduledAt)
}

func TestCreateScheduledFlagChangeWithVariationChanges(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	// Create a feature first
	featureID := newFeatureID(t)
	createFeatureNoCmd(t, client, newCreateFeatureReq(featureID))

	// Get the feature to get variation IDs
	feature := getFeature(t, featureID, client)
	require.GreaterOrEqual(t, len(feature.Variations), 2)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create scheduled flag change with variation update
	req := &featureproto.CreateScheduledFlagChangeRequest{
		EnvironmentId: *environmentID,
		FeatureId:     featureID,
		ScheduledAt:   time.Now().Add(1 * time.Hour).Unix(),
		Timezone:      "Asia/Tokyo",
		Payload: &featureproto.ScheduledChangePayload{
			VariationChanges: []*featureproto.VariationChange{
				{
					ChangeType: featureproto.ChangeType_UPDATE,
					Variation: &featureproto.Variation{
						Id:    feature.Variations[0].Id,
						Value: "updated-value",
						Name:  feature.Variations[0].Name,
					},
				},
			},
		},
		Comment: "e2e test - update variation value",
	}

	resp, err := client.CreateScheduledFlagChange(ctx, req)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.ScheduledFlagChange.Id)
	assert.Equal(t, featureproto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_VARIATIONS, resp.ScheduledFlagChange.Category)
}

func TestCreateScheduledFlagChangeWithMultipleChanges(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	// Create a feature first
	featureID := newFeatureID(t)
	createFeatureNoCmd(t, client, newCreateFeatureReq(featureID))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create scheduled flag change with multiple changes (mixed category)
	req := &featureproto.CreateScheduledFlagChangeRequest{
		EnvironmentId: *environmentID,
		FeatureId:     featureID,
		ScheduledAt:   time.Now().Add(1 * time.Hour).Unix(),
		Timezone:      "UTC",
		Payload: &featureproto.ScheduledChangePayload{
			Enabled:     wrapperspb.Bool(true),
			Name:        wrapperspb.String("new-feature-name"),
			Description: wrapperspb.String("new description"),
			TagChanges: []*featureproto.TagChange{
				{
					ChangeType: featureproto.ChangeType_CREATE,
					Tag:        "new-tag",
				},
			},
		},
		Comment: "e2e test - multiple changes",
	}

	resp, err := client.CreateScheduledFlagChange(ctx, req)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.ScheduledFlagChange.Id)
	// Multiple settings changes should be SETTINGS category
	assert.Equal(t, featureproto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_SETTINGS, resp.ScheduledFlagChange.Category)
	// Should have multiple summaries
	assert.GreaterOrEqual(t, len(resp.ScheduledFlagChange.ChangeSummaries), 3)
}

func TestArchiveFeatureCancelsPendingScheduledChanges(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	// Create a feature first
	featureID := newFeatureID(t)
	createFeatureNoCmd(t, client, newCreateFeatureReq(featureID))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create multiple pending scheduled flag changes
	createResp1 := createScheduledFlagChange(t, client, featureID, &featureproto.ScheduledChangePayload{
		Enabled: wrapperspb.Bool(true),
	})
	createResp2 := createScheduledFlagChange(t, client, featureID, &featureproto.ScheduledChangePayload{
		Enabled: wrapperspb.Bool(false),
	})

	// Verify they are pending
	getResp1 := getScheduledFlagChange(t, client, createResp1.ScheduledFlagChange.Id)
	assert.Equal(t, featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING, getResp1.ScheduledFlagChange.Status)
	getResp2 := getScheduledFlagChange(t, client, createResp2.ScheduledFlagChange.Id)
	assert.Equal(t, featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING, getResp2.ScheduledFlagChange.Status)

	// Archive the feature
	archiveReq := &featureproto.ArchiveFeatureRequest{
		Id:            featureID,
		Command:       &featureproto.ArchiveFeatureCommand{},
		EnvironmentId: *environmentID,
	}
	_, err := client.ArchiveFeature(ctx, archiveReq)
	require.NoError(t, err)

	// Verify both scheduled changes are now CANCELLED
	getResp1After := getScheduledFlagChange(t, client, createResp1.ScheduledFlagChange.Id)
	assert.Equal(t, featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CANCELLED, getResp1After.ScheduledFlagChange.Status)

	getResp2After := getScheduledFlagChange(t, client, createResp2.ScheduledFlagChange.Id)
	assert.Equal(t, featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CANCELLED, getResp2After.ScheduledFlagChange.Status)

	// Verify listing pending returns none for this feature
	listReq := &featureproto.ListScheduledFlagChangesRequest{
		EnvironmentId: *environmentID,
		FeatureId:     featureID,
		Statuses:      []featureproto.ScheduledFlagChangeStatus{featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING},
		PageSize:      10,
	}
	listResp, err := client.ListScheduledFlagChanges(ctx, listReq)
	require.NoError(t, err)
	assert.Empty(t, listResp.ScheduledFlagChanges, "No pending scheduled changes should exist after archiving")
}

// Helper functions

func createScheduledFlagChange(
	t *testing.T,
	client featureclient.Client,
	featureID string,
	payload *featureproto.ScheduledChangePayload,
) *featureproto.CreateScheduledFlagChangeResponse {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := &featureproto.CreateScheduledFlagChangeRequest{
		EnvironmentId: *environmentID,
		FeatureId:     featureID,
		ScheduledAt:   time.Now().Add(1 * time.Hour).Unix(),
		Timezone:      "UTC",
		Payload:       payload,
		Comment:       "e2e test scheduled change",
	}

	resp, err := client.CreateScheduledFlagChange(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create scheduled flag change: %v", err)
	}
	return resp
}

func getScheduledFlagChange(
	t *testing.T,
	client featureclient.Client,
	id string,
) *featureproto.GetScheduledFlagChangeResponse {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := &featureproto.GetScheduledFlagChangeRequest{
		EnvironmentId: *environmentID,
		Id:            id,
	}

	resp, err := client.GetScheduledFlagChange(ctx, req)
	if err != nil {
		t.Fatalf("Failed to get scheduled flag change: %v", err)
	}
	return resp
}
