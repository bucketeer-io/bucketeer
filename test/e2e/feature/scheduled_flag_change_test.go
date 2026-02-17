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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	auditlogclient "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	rpcclient "github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestCreateScheduledFlagChange(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	// Create a feature first
	featureID := newFeatureID(t)
	createFeature(t, client, newCreateFeatureReq(featureID))

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
	createFeature(t, client, newCreateFeatureReq(featureID))

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
	createFeature(t, client, newCreateFeatureReq(featureID))

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
	createFeature(t, client, newCreateFeatureReq(featureID))

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
	createFeature(t, client, newCreateFeatureReq(featureID))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create multiple scheduled flag changes
	createResp1 := createScheduledFlagChange(t, client, featureID, &featureproto.ScheduledChangePayload{
		Enabled: wrapperspb.Bool(true),
	})
	time.Sleep(1 * time.Second) // Ensure different created_at times
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
	createFeature(t, client, newCreateFeatureReq(featureID))

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
	createFeature(t, client, newCreateFeatureReq(featureID))

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
	createFeature(t, client, newCreateFeatureReq(featureID))

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
	createFeature(t, client, newCreateFeatureReq(featureID))

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

// TestArchiveFeatureCancelsPendingScheduledChanges verifies that when a feature is archived,
// all pending and conflicting scheduled changes for that feature are automatically cancelled.
// This prevents orphaned schedules from attempting to execute on archived flags.
func TestArchiveFeatureCancelsPendingScheduledChanges(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	// Create a feature first
	featureID := newFeatureID(t)
	createFeature(t, client, newCreateFeatureReq(featureID))

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

	// Archive the feature using UpdateFeature
	updateReq := &featureproto.UpdateFeatureRequest{
		Id:            featureID,
		EnvironmentId: *environmentID,
		Archived:      wrapperspb.Bool(true),
		Comment:       "Archiving feature to test scheduled change cancellation",
	}
	_, err := client.UpdateFeature(ctx, updateReq)
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

func TestScheduledFlagChangeAuditLogs(t *testing.T) {
	t.Parallel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()

	auditLogClient := newAuditLogClient(t)
	defer auditLogClient.Close()

	// Record start time for filtering audit logs
	startTime := time.Now().Unix() - 60 // 60 seconds buffer for clock skew

	// Create a feature first
	featureID := newFeatureID(t)
	createFeature(t, featureClient, newCreateFeatureReq(featureID))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 1. Create a scheduled flag change
	createResp := createScheduledFlagChange(t, featureClient, featureID, &featureproto.ScheduledChangePayload{
		Enabled: wrapperspb.Bool(true),
	})
	scheduleID := createResp.ScheduledFlagChange.Id
	t.Logf("Created scheduled flag change with ID: %s", scheduleID)

	// 2. Update the scheduled flag change
	newScheduledAt := time.Now().Add(2 * time.Hour).Unix()
	updateReq := &featureproto.UpdateScheduledFlagChangeRequest{
		EnvironmentId: *environmentID,
		Id:            scheduleID,
		ScheduledAt:   wrapperspb.Int64(newScheduledAt),
		Comment:       wrapperspb.String("Updated schedule time"),
	}
	_, err := featureClient.UpdateScheduledFlagChange(ctx, updateReq)
	require.NoError(t, err)
	t.Log("Updated scheduled flag change")

	// 3. Delete (cancel) the scheduled flag change
	deleteReq := &featureproto.DeleteScheduledFlagChangeRequest{
		EnvironmentId: *environmentID,
		Id:            scheduleID,
	}
	_, err = featureClient.DeleteScheduledFlagChange(ctx, deleteReq)
	require.NoError(t, err)
	t.Log("Cancelled scheduled flag change")

	// Wait for audit logs to be processed
	time.Sleep(5 * time.Second)

	// 4. List audit logs and verify all three events are recorded
	toTime := time.Now().Unix() + 60 // 60 seconds buffer

	// Search for audit logs with SCHEDULED_FLAG_CHANGE entity type
	var foundCreated, foundUpdated, foundCancelled bool
	maxRetries := 30
	for retry := 0; retry < maxRetries; retry++ {
		time.Sleep(2 * time.Second)

		listResp, err := auditLogClient.ListAuditLogs(ctx, &auditlog.ListAuditLogsRequest{
			EnvironmentId:  *environmentID,
			EntityType:     wrapperspb.Int32(int32(eventproto.Event_SCHEDULED_FLAG_CHANGE)),
			PageSize:       50,
			Cursor:         "0",
			OrderBy:        auditlog.ListAuditLogsRequest_TIMESTAMP,
			OrderDirection: auditlog.ListAuditLogsRequest_DESC,
			From:           startTime,
			To:             toTime,
		})
		if err != nil {
			t.Logf("Retry %d: Failed to list audit logs: %v", retry+1, err)
			continue
		}

		for _, log := range listResp.AuditLogs {
			if log.EntityId == scheduleID {
				switch log.Type {
				case eventproto.Event_SCHEDULED_FLAG_CHANGE_CREATED:
					foundCreated = true
					t.Logf("Found CREATED audit log: %s", log.Id)
				case eventproto.Event_SCHEDULED_FLAG_CHANGE_UPDATED:
					foundUpdated = true
					t.Logf("Found UPDATED audit log: %s", log.Id)
				case eventproto.Event_SCHEDULED_FLAG_CHANGE_CANCELLED:
					foundCancelled = true
					t.Logf("Found CANCELLED audit log: %s", log.Id)
				}
			}
		}

		if foundCreated && foundUpdated && foundCancelled {
			t.Logf("All audit logs found after %d retries", retry+1)
			break
		}
	}

	// Assert all events were found
	assert.True(t, foundCreated, "SCHEDULED_FLAG_CHANGE_CREATED audit log should exist")
	assert.True(t, foundUpdated, "SCHEDULED_FLAG_CHANGE_UPDATED audit log should exist")
	assert.True(t, foundCancelled, "SCHEDULED_FLAG_CHANGE_CANCELLED audit log should exist")
}

// TestScheduledFlagChangeConflict_DependencyMissing verifies that when creating
// a new schedule that references a variation which an earlier schedule deletes,
// the response includes a DEPENDENCY_MISSING conflict.
func TestScheduledFlagChangeConflict_DependencyMissing(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	// Create a feature
	featureID := newFeatureID(t)
	createFeature(t, client, newCreateFeatureReq(featureID))

	// Get the feature to get variation IDs
	feature := getFeature(t, featureID, client)
	require.GreaterOrEqual(t, len(feature.Variations), 2)
	varToDelete := feature.Variations[2] // Variation C (not used by default/off)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Schedule A at +1h: delete variation C
	scheduleAResp := createScheduledFlagChangeAt(t, client, featureID,
		&featureproto.ScheduledChangePayload{
			VariationChanges: []*featureproto.VariationChange{
				{
					ChangeType: featureproto.ChangeType_DELETE,
					Variation:  &featureproto.Variation{Id: varToDelete.Id},
				},
			},
		},
		time.Now().Add(1*time.Hour).Unix(),
	)
	require.NotEmpty(t, scheduleAResp.ScheduledFlagChange.Id)

	// Schedule B at +2h: update the same variation C (which Schedule A deletes)
	scheduleBReq := &featureproto.CreateScheduledFlagChangeRequest{
		EnvironmentId: *environmentID,
		FeatureId:     featureID,
		ScheduledAt:   time.Now().Add(2 * time.Hour).Unix(),
		Timezone:      "UTC",
		Payload: &featureproto.ScheduledChangePayload{
			VariationChanges: []*featureproto.VariationChange{
				{
					ChangeType: featureproto.ChangeType_UPDATE,
					Variation: &featureproto.Variation{
						Id:    varToDelete.Id,
						Value: "updated-value",
						Name:  varToDelete.Name,
					},
				},
			},
		},
		Comment: "e2e test - update variation that earlier schedule deletes",
	}
	scheduleBResp, err := client.CreateScheduledFlagChange(ctx, scheduleBReq)
	require.NoError(t, err)
	require.NotEmpty(t, scheduleBResp.ScheduledFlagChange.Id)

	// Verify DEPENDENCY_MISSING conflict is detected
	require.NotEmpty(t, scheduleBResp.DetectedConflicts,
		"Expected DEPENDENCY_MISSING conflict when earlier schedule deletes referenced variation")
	foundDependencyMissing := false
	for _, c := range scheduleBResp.DetectedConflicts {
		if c.Type == featureproto.ScheduledChangeConflict_CONFLICT_TYPE_DEPENDENCY_MISSING {
			foundDependencyMissing = true
			assert.Contains(t, c.ConflictingScheduleId, scheduleAResp.ScheduledFlagChange.Id)
			t.Logf("DEPENDENCY_MISSING conflict: %s (field: %s)", c.Description, c.ConflictingField)
		}
	}
	assert.True(t, foundDependencyMissing, "Expected at least one DEPENDENCY_MISSING conflict")
}

// TestScheduledFlagChangeConflict_InvalidReferenceOnFlagUpdate verifies that when
// a flag is updated directly (e.g., a variation is deleted), pending schedules that
// reference the deleted variation are marked as CONFLICT with INVALID_REFERENCE.
func TestScheduledFlagChangeConflict_InvalidReferenceOnFlagUpdate(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	// Create a feature
	featureID := newFeatureID(t)
	createFeature(t, client, newCreateFeatureReq(featureID))

	// Get the feature to get variation IDs
	feature := getFeature(t, featureID, client)
	require.GreaterOrEqual(t, len(feature.Variations), 4)
	// Use variation C (index 2) — not the default on/off variation
	varToDelete := feature.Variations[2]

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a schedule that updates variation C
	scheduleResp := createScheduledFlagChangeAt(t, client, featureID,
		&featureproto.ScheduledChangePayload{
			VariationChanges: []*featureproto.VariationChange{
				{
					ChangeType: featureproto.ChangeType_UPDATE,
					Variation: &featureproto.Variation{
						Id:    varToDelete.Id,
						Value: "scheduled-updated-value",
						Name:  varToDelete.Name,
					},
				},
			},
		},
		time.Now().Add(1*time.Hour).Unix(),
	)
	require.NotEmpty(t, scheduleResp.ScheduledFlagChange.Id)

	// Verify the schedule is PENDING
	getResp := getScheduledFlagChange(t, client, scheduleResp.ScheduledFlagChange.Id)
	assert.Equal(t,
		featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
		getResp.ScheduledFlagChange.Status,
	)

	// Now directly delete variation C via UpdateFeature
	updateReq := &featureproto.UpdateFeatureRequest{
		Id:            featureID,
		EnvironmentId: *environmentID,
		VariationChanges: []*featureproto.VariationChange{
			{
				ChangeType: featureproto.ChangeType_DELETE,
				Variation:  &featureproto.Variation{Id: varToDelete.Id},
			},
		},
		Comment: "e2e test - delete variation to trigger conflict",
	}
	updateResp, err := client.UpdateFeature(ctx, updateReq)
	require.NoError(t, err)

	// Verify UpdateFeatureResponse indicates conflicts were detected
	assert.True(t, updateResp.ScheduleConflictsDetected,
		"Expected ScheduleConflictsDetected=true after deleting a variation referenced by a schedule")
	assert.GreaterOrEqual(t, updateResp.ConflictCount, int32(1))
	t.Logf("UpdateFeature detected %d schedule conflict(s)", updateResp.ConflictCount)

	// Verify the schedule is now marked as CONFLICT
	getAfter := getScheduledFlagChange(t, client, scheduleResp.ScheduledFlagChange.Id)
	assert.Equal(t,
		featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT,
		getAfter.ScheduledFlagChange.Status,
		"Schedule should be CONFLICT after referenced variation was deleted",
	)

	// Verify INVALID_REFERENCE conflict details
	require.NotEmpty(t, getAfter.ScheduledFlagChange.Conflicts)
	foundInvalidRef := false
	for _, c := range getAfter.ScheduledFlagChange.Conflicts {
		if c.Type == featureproto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE {
			foundInvalidRef = true
			t.Logf("INVALID_REFERENCE conflict: %s (field: %s)", c.Description, c.ConflictingField)
		}
	}
	assert.True(t, foundInvalidRef, "Expected INVALID_REFERENCE conflict on the schedule")
}

// TestScheduledFlagChangeConflict_NoConflictOnUnrelatedFlagUpdate verifies that
// updating an unrelated field on the flag (e.g., description) does NOT mark
// pending schedules as CONFLICT.
func TestScheduledFlagChangeConflict_NoConflictOnUnrelatedFlagUpdate(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()

	// Create a feature
	featureID := newFeatureID(t)
	createFeature(t, client, newCreateFeatureReq(featureID))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a schedule that enables the flag
	scheduleResp := createScheduledFlagChange(t, client, featureID,
		&featureproto.ScheduledChangePayload{
			Enabled: wrapperspb.Bool(true),
		},
	)
	require.NotEmpty(t, scheduleResp.ScheduledFlagChange.Id)

	// Update an unrelated field (description)
	updateReq := &featureproto.UpdateFeatureRequest{
		Id:            featureID,
		EnvironmentId: *environmentID,
		Description:   wrapperspb.String("updated description - should not cause conflict"),
		Comment:       "e2e test - unrelated update",
	}
	updateResp, err := client.UpdateFeature(ctx, updateReq)
	require.NoError(t, err)

	// Verify no conflicts detected
	assert.False(t, updateResp.ScheduleConflictsDetected,
		"Unrelated flag update should NOT cause schedule conflicts")
	assert.Equal(t, int32(0), updateResp.ConflictCount)

	// Verify the schedule is still PENDING
	getAfter := getScheduledFlagChange(t, client, scheduleResp.ScheduledFlagChange.Id)
	assert.Equal(t,
		featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
		getAfter.ScheduledFlagChange.Status,
		"Schedule should still be PENDING after unrelated flag update",
	)
}

// Helper functions

func createScheduledFlagChangeAt(
	t *testing.T,
	client featureclient.Client,
	featureID string,
	payload *featureproto.ScheduledChangePayload,
	scheduledAt int64,
) *featureproto.CreateScheduledFlagChangeResponse {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := &featureproto.CreateScheduledFlagChangeRequest{
		EnvironmentId: *environmentID,
		FeatureId:     featureID,
		ScheduledAt:   scheduledAt,
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

func newAuditLogClient(t *testing.T) auditlogclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := auditlogclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create auditlog client:", err)
	}
	return client
}

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
