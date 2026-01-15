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
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"

	btclient "github.com/bucketeer-io/bucketeer/v2/pkg/batch/client"
	environmentclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	rpcclient "github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	btproto "github.com/bucketeer-io/bucketeer/v2/proto/batch"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	"github.com/bucketeer-io/bucketeer/v2/proto/feature"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

const (
	autoArchiveTimeout       = 3 * time.Minute
	autoArchiveRetryTimes    = 30
	autoArchiveRetryInterval = 5 * time.Second
	// thresholdForArchiving is set to 0 to immediately archive features with LastUsedInfo.
	// This allows E2E tests to verify archiving behavior without waiting for days to pass.
	thresholdForArchiving = 0
	// thresholdToPreventArchiving is set very high to ensure features are NOT archived.
	// Used to test that recently used features are not archived.
	thresholdToPreventArchiving = 365
)

// TestFeatureAutoArchiver_BasicAutoArchive tests that a feature flag with LastUsedInfo
// is automatically archived when the threshold is set to 0 days.
func TestFeatureAutoArchiver_BasicAutoArchive(t *testing.T) {
	t.Parallel()

	// Setup clients
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	envClient := newEnvironmentClient(t)
	defer envClient.Close()
	batchClient := newBatchClient(t)
	defer batchClient.Close()

	// Create a test feature
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, featureClient, cmd)

	// Enable the feature to generate last used info
	enableFeature(t, featureID, featureClient)
	f := getFeature(t, featureID, featureClient)

	// Register an evaluation event with current timestamp
	// The threshold is set to 0, so any feature with LastUsedInfo is archivable
	registerEvaluationEventWithTimestamp(t, f, time.Now())

	// Wait for last used info to be recorded
	waitForLastUsedInfo(t, featureClient, featureID)

	// Get original auto-archive settings
	originalSettings := getAutoArchiveSettings(t, envClient)

	// Enable auto-archive with threshold=0 (any feature with LastUsedInfo is archivable)
	enableAutoArchive(t, envClient, thresholdForArchiving, false)
	defer restoreAutoArchiveSettings(t, envClient, originalSettings)

	// Execute the auto-archive batch job
	executeBatchJob(t, batchClient, btproto.BatchJob_FeatureAutoArchiver)

	// Verify the feature is archived
	var archived bool
	for i := 0; i < autoArchiveRetryTimes; i++ {
		feature := getFeature(t, featureID, featureClient)
		if feature.Archived {
			archived = true
			break
		}
		time.Sleep(autoArchiveRetryInterval)
	}
	assert.True(t, archived, "Feature should be archived after auto-archive batch job")
}

// TestFeatureAutoArchiver_PrerequisiteDependency tests that a feature flag that is
// referenced by other flags (as a prerequisite) is NOT archived, even if it meets
// the archival criteria.
func TestFeatureAutoArchiver_PrerequisiteDependency(t *testing.T) {
	t.Parallel()

	// Setup clients
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	envClient := newEnvironmentClient(t)
	defer envClient.Close()
	batchClient := newBatchClient(t)
	defer batchClient.Close()

	// Create parent feature (will be depended upon)
	parentFeatureID := fmt.Sprintf("%s-parent-%s", prefixID, newUUID(t))
	parentCmd := newCreateFeatureCommand(parentFeatureID)
	createFeature(t, featureClient, parentCmd)
	enableFeature(t, parentFeatureID, featureClient)
	parentFeature := getFeature(t, parentFeatureID, featureClient)

	// Create child feature that depends on parent
	childFeatureID := fmt.Sprintf("%s-child-%s", prefixID, newUUID(t))
	childCmd := newCreateFeatureCommand(childFeatureID)
	createFeature(t, featureClient, childCmd)
	enableFeature(t, childFeatureID, featureClient)

	// Add prerequisite: child depends on parent
	addPrerequisite(t, featureClient, childFeatureID, parentFeatureID, parentFeature.Variations[0].Id)

	// Register evaluation events with current timestamp for both features
	now := time.Now()
	registerEvaluationEventWithTimestamp(t, parentFeature, now)
	childFeature := getFeature(t, childFeatureID, featureClient)
	registerEvaluationEventWithTimestamp(t, childFeature, now)

	// Wait for last used info to be recorded
	waitForLastUsedInfo(t, featureClient, parentFeatureID)
	waitForLastUsedInfo(t, featureClient, childFeatureID)

	// Get original auto-archive settings
	originalSettings := getAutoArchiveSettings(t, envClient)

	// Enable auto-archive with threshold=0 (any feature with LastUsedInfo is archivable)
	enableAutoArchive(t, envClient, thresholdForArchiving, false)
	defer restoreAutoArchiveSettings(t, envClient, originalSettings)

	// Execute the auto-archive batch job
	executeBatchJob(t, batchClient, btproto.BatchJob_FeatureAutoArchiver)

	// Wait a bit for processing
	time.Sleep(autoArchiveRetryInterval * 2)

	// Verify parent feature is NOT archived (because it's a dependency)
	parentResult := getFeature(t, parentFeatureID, featureClient)
	assert.False(t, parentResult.Archived, "Parent feature should NOT be archived because it is a prerequisite for another feature")

	// Verify child feature IS archived (it depends on parent but is not depended upon)
	var childArchived bool
	for i := 0; i < autoArchiveRetryTimes; i++ {
		childResult := getFeature(t, childFeatureID, featureClient)
		if childResult.Archived {
			childArchived = true
			break
		}
		time.Sleep(autoArchiveRetryInterval)
	}
	assert.True(t, childArchived, "Child feature should be archived")
}

// TestFeatureAutoArchiver_BulkArchive tests that multiple feature flags that meet
// the archival criteria are all archived in a single batch job execution.
func TestFeatureAutoArchiver_BulkArchive(t *testing.T) {
	t.Parallel()

	// Setup clients
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	envClient := newEnvironmentClient(t)
	defer envClient.Close()
	batchClient := newBatchClient(t)
	defer batchClient.Close()

	// Create multiple test features
	numFeatures := 3
	featureIDs := make([]string, numFeatures)
	for i := 0; i < numFeatures; i++ {
		featureID := fmt.Sprintf("%s-bulk-%d-%s", prefixID, i, newUUID(t))
		featureIDs[i] = featureID
		cmd := newCreateFeatureCommand(featureID)
		createFeature(t, featureClient, cmd)
		enableFeature(t, featureID, featureClient)
	}

	// Register evaluation events with current timestamp for all features
	now := time.Now()
	for _, featureID := range featureIDs {
		f := getFeature(t, featureID, featureClient)
		registerEvaluationEventWithTimestamp(t, f, now)
	}

	// Wait for last used info to be recorded for all features
	for _, featureID := range featureIDs {
		waitForLastUsedInfo(t, featureClient, featureID)
	}

	// Get original auto-archive settings
	originalSettings := getAutoArchiveSettings(t, envClient)

	// Enable auto-archive with threshold=0 (any feature with LastUsedInfo is archivable)
	enableAutoArchive(t, envClient, thresholdForArchiving, false)
	defer restoreAutoArchiveSettings(t, envClient, originalSettings)

	// Execute the auto-archive batch job
	executeBatchJob(t, batchClient, btproto.BatchJob_FeatureAutoArchiver)

	// Verify all features are archived
	for _, featureID := range featureIDs {
		var archived bool
		for i := 0; i < autoArchiveRetryTimes; i++ {
			feature := getFeature(t, featureID, featureClient)
			if feature.Archived {
				archived = true
				break
			}
			time.Sleep(autoArchiveRetryInterval)
		}
		assert.True(t, archived, "Feature %s should be archived after auto-archive batch job", featureID)
	}
}

// TestFeatureAutoArchiver_DisabledEnvironment tests that when auto-archive is disabled
// for an environment, no features are archived even if they would otherwise be archivable.
func TestFeatureAutoArchiver_DisabledEnvironment(t *testing.T) {
	t.Parallel()

	// Setup clients
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	envClient := newEnvironmentClient(t)
	defer envClient.Close()
	batchClient := newBatchClient(t)
	defer batchClient.Close()

	// Create a test feature
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, featureClient, cmd)
	enableFeature(t, featureID, featureClient)
	f := getFeature(t, featureID, featureClient)

	// Register an evaluation event with current timestamp
	registerEvaluationEventWithTimestamp(t, f, time.Now())

	// Wait for last used info to be recorded
	waitForLastUsedInfo(t, featureClient, featureID)

	// Get original auto-archive settings
	originalSettings := getAutoArchiveSettings(t, envClient)

	// Ensure auto-archive is DISABLED for the environment
	disableAutoArchive(t, envClient)
	defer restoreAutoArchiveSettings(t, envClient, originalSettings)

	// Execute the auto-archive batch job
	executeBatchJob(t, batchClient, btproto.BatchJob_FeatureAutoArchiver)

	// Wait a bit for processing
	time.Sleep(autoArchiveRetryInterval * 2)

	// Verify the feature is NOT archived
	feature := getFeature(t, featureID, featureClient)
	assert.False(t, feature.Archived, "Feature should NOT be archived when auto-archive is disabled")
}

// TestFeatureAutoArchiver_RecentlyUsedNotArchived tests that a feature flag that has been
// recently used (within the threshold) is NOT archived.
func TestFeatureAutoArchiver_RecentlyUsedNotArchived(t *testing.T) {
	t.Parallel()

	// Setup clients
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	envClient := newEnvironmentClient(t)
	defer envClient.Close()
	batchClient := newBatchClient(t)
	defer batchClient.Close()

	// Create a test feature
	featureID := newFeatureID(t)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, featureClient, cmd)
	enableFeature(t, featureID, featureClient)
	f := getFeature(t, featureID, featureClient)

	// Register an evaluation event with current timestamp (just used)
	registerEvaluationEventWithTimestamp(t, f, time.Now())

	// Wait for last used info to be recorded
	waitForLastUsedInfo(t, featureClient, featureID)

	// Get original auto-archive settings
	originalSettings := getAutoArchiveSettings(t, envClient)

	// Enable auto-archive with very high threshold (365 days)
	// Since the feature was just used (UnusedDays=0), it should NOT be archived
	enableAutoArchive(t, envClient, thresholdToPreventArchiving, false)
	defer restoreAutoArchiveSettings(t, envClient, originalSettings)

	// Execute the auto-archive batch job
	executeBatchJob(t, batchClient, btproto.BatchJob_FeatureAutoArchiver)

	// Wait a bit for processing
	time.Sleep(autoArchiveRetryInterval * 2)

	// Verify the feature is NOT archived (recently used)
	feature := getFeature(t, featureID, featureClient)
	assert.False(t, feature.Archived, "Feature should NOT be archived when recently used")
}

// Helper functions

type autoArchiveSettings struct {
	enabled       bool
	unusedDays    int32
	checkCodeRefs bool
}

func newEnvironmentClient(t *testing.T) environmentclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := environmentclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create environment client:", err)
	}
	return client
}

func newBatchClient(t *testing.T) btclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := btclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create batch client:", err)
	}
	return client
}

func getAutoArchiveSettings(t *testing.T, client environmentclient.Client) *autoArchiveSettings {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := client.GetEnvironmentV2(ctx, &environmentproto.GetEnvironmentV2Request{
		Id: *environmentID,
	})
	if err != nil {
		t.Fatal("Failed to get auto-archive settings:", err)
	}

	return &autoArchiveSettings{
		enabled:       resp.Environment.AutoArchiveEnabled,
		unusedDays:    resp.Environment.AutoArchiveUnusedDays,
		checkCodeRefs: resp.Environment.AutoArchiveCheckCodeRefs,
	}
}

func enableAutoArchive(t *testing.T, client environmentclient.Client, unusedDays int32, checkCodeRefs bool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := client.UpdateEnvironmentV2(ctx, &environmentproto.UpdateEnvironmentV2Request{
		Id:                       *environmentID,
		AutoArchiveEnabled:       wrapperspb.Bool(true),
		AutoArchiveUnusedDays:    wrapperspb.Int32(unusedDays),
		AutoArchiveCheckCodeRefs: wrapperspb.Bool(checkCodeRefs),
	})
	if err != nil {
		t.Fatal("Failed to enable auto-archive:", err)
	}
}

func disableAutoArchive(t *testing.T, client environmentclient.Client) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := client.UpdateEnvironmentV2(ctx, &environmentproto.UpdateEnvironmentV2Request{
		Id:                 *environmentID,
		AutoArchiveEnabled: wrapperspb.Bool(false),
	})
	if err != nil {
		t.Fatal("Failed to disable auto-archive:", err)
	}
}

func restoreAutoArchiveSettings(t *testing.T, client environmentclient.Client, settings *autoArchiveSettings) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// When restoring to disabled state, we need to update settings first (while enabled),
	// then disable. When restoring to enabled state, we can update everything at once.
	if settings.enabled {
		// Restore to enabled state with all settings
		_, err := client.UpdateEnvironmentV2(ctx, &environmentproto.UpdateEnvironmentV2Request{
			Id:                       *environmentID,
			AutoArchiveEnabled:       wrapperspb.Bool(true),
			AutoArchiveUnusedDays:    wrapperspb.Int32(settings.unusedDays),
			AutoArchiveCheckCodeRefs: wrapperspb.Bool(settings.checkCodeRefs),
		})
		if err != nil {
			t.Logf("Warning: Failed to restore auto-archive settings: %v", err)
		}
	} else {
		// Just disable auto-archive (don't try to update other settings when disabling)
		_, err := client.UpdateEnvironmentV2(ctx, &environmentproto.UpdateEnvironmentV2Request{
			Id:                 *environmentID,
			AutoArchiveEnabled: wrapperspb.Bool(false),
		})
		if err != nil {
			t.Logf("Warning: Failed to restore auto-archive settings: %v", err)
		}
	}
}

func executeBatchJob(t *testing.T, client btclient.Client, job btproto.BatchJob) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), autoArchiveTimeout)
	defer cancel()

	numRetries := 5
	var err error
	for i := 0; i < numRetries; i++ {
		_, err = client.ExecuteBatchJob(ctx, &btproto.BatchJobRequest{Job: job})
		if err == nil {
			return
		}
		// FailedPrecondition errors are expected when the environment contains
		// features with prerequisites that can't be archived. This is not a test failure.
		if strings.Contains(err.Error(), "FailedPrecondition") ||
			strings.Contains(err.Error(), "used as a prerequsite") {
			t.Logf("Batch job completed with expected prerequisite warning: %v", err)
			return
		}
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		t.Fatal("Failed to execute batch job after retries:", err)
	}
}

func addPrerequisite(
	t *testing.T,
	client featureclient.Client,
	featureID, prerequisiteFeatureID, prerequisiteVariationID string,
) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.UpdateFeature(ctx, &featureproto.UpdateFeatureRequest{
		EnvironmentId: *environmentID,
		Id:            featureID,
		PrerequisiteChanges: []*featureproto.PrerequisiteChange{
			{
				ChangeType: featureproto.ChangeType_CREATE,
				Prerequisite: &featureproto.Prerequisite{
					FeatureId:   prerequisiteFeatureID,
					VariationId: prerequisiteVariationID,
				},
			},
		},
	})
	if err != nil {
		t.Fatal("Failed to add prerequisite:", err)
	}
}

func registerEvaluationEventWithTimestamp(t *testing.T, f *feature.Feature, timestamp time.Time) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	evaluation, err := ptypes.MarshalAny(&eventproto.EvaluationEvent{
		Timestamp:      timestamp.Unix(),
		FeatureId:      f.Id,
		FeatureVersion: f.Version,
		UserId:         "e2e-auto-archive-test-user",
		VariationId:    f.Variations[0].Id,
		User:           &userproto.User{Id: "e2e-auto-archive-test-user"},
		Reason:         &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
		Tag:            f.Tags[0],
	})
	if err != nil {
		t.Fatal("Failed to marshal evaluation event:", err)
	}

	events := []*eventproto.Event{
		{
			Id:    newUUID(t),
			Event: evaluation,
		},
	}
	req := &gatewayproto.RegisterEventsRequest{Events: events}
	_, err = c.RegisterEvents(ctx, req)
	if err != nil {
		t.Fatal("Failed to register events:", err)
	}
}

func waitForLastUsedInfo(t *testing.T, client featureclient.Client, featureID string) {
	t.Helper()
	for i := 0; i < featureRecorderRetryTimes; i++ {
		f := getFeature(t, featureID, client)
		if f.LastUsedInfo != nil {
			return
		}
		time.Sleep(time.Second)
	}
	t.Logf("Warning: LastUsedInfo not recorded for feature %s after retries", featureID)
}
