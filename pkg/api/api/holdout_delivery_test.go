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

package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/blang/semver"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

const holdoutDeliveryAPIKey = "holdout-delivery-api-key"

// storedHoldoutGroup is a serving group as the cacher stores it: every
// column of the row, and a target that has since been removed.
func storedHoldoutGroup(environmentID string, updatedAt int64) *featureproto.HoldoutGroup {
	return &featureproto.HoldoutGroup{
		Id:                   "hg-1",
		EnvironmentId:        environmentID,
		Name:                 "Q3 holdout",
		Description:          "Everything shipped in Q3",
		AssignmentType:       featureproto.HoldoutGroup_HASH,
		SamplingSeed:         "seed",
		BucketStart:          12.5,
		Percentage:           5,
		Status:               featureproto.HoldoutGroup_RUNNING,
		StartAt:              100,
		StopAt:               200,
		StartAnalysisAt:      150,
		GoalIds:              []string{"goal-1"},
		ExperimentId:         "experiment-1",
		AnalysisExperimentId: "experiment-2",
		CreatedBy:            "creator@example.com",
		UpdatedBy:            "editor@example.com",
		CreatedAt:            90,
		UpdatedAt:            updatedAt,
		Features: []*featureproto.HoldoutGroupFeature{
			{
				Id:                 "row-a",
				HoldoutGroupId:     "hg-1",
				FeatureId:          "flag-a",
				EnvironmentId:      environmentID,
				HoldoutVariationId: "variation-off",
				AddedAt:            120,
			},
			{
				Id:                 "row-b",
				HoldoutGroupId:     "hg-1",
				FeatureId:          "flag-b",
				EnvironmentId:      environmentID,
				HoldoutVariationId: "variation-off",
				AddedAt:            120,
				RemovedAt:          130,
			},
		},
	}
}

// deliveredHoldoutGroup is what a server SDK receives for storedHoldoutGroup.
func deliveredHoldoutGroup(environmentID string, updatedAt int64) *featureproto.HoldoutGroup {
	return &featureproto.HoldoutGroup{
		Id:             "hg-1",
		EnvironmentId:  environmentID,
		AssignmentType: featureproto.HoldoutGroup_HASH,
		SamplingSeed:   "seed",
		BucketStart:    12.5,
		Percentage:     5,
		Status:         featureproto.HoldoutGroup_RUNNING,
		UpdatedAt:      updatedAt,
		Features: []*featureproto.HoldoutGroupFeature{
			{
				Id:                 "row-a",
				HoldoutGroupId:     "hg-1",
				FeatureId:          "flag-a",
				EnvironmentId:      environmentID,
				HoldoutVariationId: "variation-off",
				AddedAt:            120,
			},
		},
	}
}

func TestSDKHoldoutGroups(t *testing.T) {
	t.Parallel()

	t.Run("nothing is serving", func(t *testing.T) {
		t.Parallel()
		stopped := storedHoldoutGroup("env", 50)
		stopped.Status = featureproto.HoldoutGroup_STOPPED
		waiting := storedHoldoutGroup("env", 50)
		waiting.Status = featureproto.HoldoutGroup_WAITING
		archived := storedHoldoutGroup("env", 50)
		archived.Archived = true
		for desc, snapshot := range map[string]*featureproto.HoldoutGroup{
			"no snapshot":           nil,
			"never stopped a group": evaluation.ClearanceSnapshot(0),
			"clearance":             evaluation.ClearanceSnapshot(1000),
			"unavailable":           evaluation.UnavailableHoldoutSnapshot(),
			"stopped":               stopped,
			"waiting":               waiting,
			"archived":              archived,
		} {
			assert.Nil(t, sdkHoldoutGroups(snapshot), desc)
		}
	})

	t.Run("a serving group is cut down to what local evaluation reads", func(t *testing.T) {
		t.Parallel()
		stored := storedHoldoutGroup("env", 50)
		before := proto.Clone(stored)

		groups := sdkHoldoutGroups(stored)

		require.Len(t, groups, 1)
		assert.True(t, proto.Equal(deliveredHoldoutGroup("env", 50), groups[0]), "got %v", groups[0])
		assert.True(t, proto.Equal(before, stored), "the cached group must not be modified")
		assert.NotSame(t, stored, groups[0])
		assert.NotSame(t, stored.Features[0], groups[0].Features[0])
	})

	t.Run("analysis keeps serving", func(t *testing.T) {
		t.Parallel()
		stored := storedHoldoutGroup("env", 50)
		stored.Status = featureproto.HoldoutGroup_ANALYSIS
		groups := sdkHoldoutGroups(stored)
		require.Len(t, groups, 1)
		assert.Equal(t, featureproto.HoldoutGroup_ANALYSIS, groups[0].Status)
	})

	t.Run("a segment group keeps its segment", func(t *testing.T) {
		t.Parallel()
		stored := storedHoldoutGroup("env", 50)
		stored.AssignmentType = featureproto.HoldoutGroup_SEGMENT
		stored.SegmentId = "segment-1"
		groups := sdkHoldoutGroups(stored)
		require.Len(t, groups, 1)
		assert.Equal(t, featureproto.HoldoutGroup_SEGMENT, groups[0].AssignmentType)
		assert.Equal(t, "segment-1", groups[0].SegmentId)
	})

	t.Run("a group without current targets is still delivered", func(t *testing.T) {
		t.Parallel()
		stored := storedHoldoutGroup("env", 50)
		stored.Features = stored.Features[1:]
		groups := sdkHoldoutGroups(stored)
		require.Len(t, groups, 1)
		assert.Empty(t, groups[0].Features)
	})

	t.Run("the delivered group has the same id input as the stored one", func(t *testing.T) {
		t.Parallel()
		stored := storedHoldoutGroup("env", 50)
		assert.Equal(t,
			evaluation.GenerateFeaturesID(nil, stored),
			evaluation.GenerateFeaturesID(nil, sdkHoldoutGroups(stored)...),
		)
	})
}

func TestHoldoutCapableSDK(t *testing.T) {
	t.Parallel()
	minVersions := map[eventproto.SourceId]semver.Version{
		eventproto.SourceId_GO_SERVER:   semver.MustParse("1.9.0"),
		eventproto.SourceId_NODE_SERVER: semver.MustParse("0.5.0"),
	}
	patterns := []struct {
		desc       string
		sourceID   eventproto.SourceId
		sdkVersion string
		expected   bool
	}{
		{desc: "minimum version", sourceID: eventproto.SourceId_GO_SERVER, sdkVersion: "1.9.0", expected: true},
		{desc: "newer patch", sourceID: eventproto.SourceId_GO_SERVER, sdkVersion: "1.9.1", expected: true},
		{desc: "newer major", sourceID: eventproto.SourceId_GO_SERVER, sdkVersion: "2.0.0", expected: true},
		{desc: "v prefix", sourceID: eventproto.SourceId_GO_SERVER, sdkVersion: "v1.9.0", expected: true},
		{desc: "older", sourceID: eventproto.SourceId_GO_SERVER, sdkVersion: "1.8.9", expected: false},
		{desc: "older with v prefix", sourceID: eventproto.SourceId_GO_SERVER, sdkVersion: "v1.8.0", expected: false},
		{desc: "pre-release of the minimum", sourceID: eventproto.SourceId_GO_SERVER, sdkVersion: "1.9.0-rc.1", expected: false},
		{desc: "numeric, not lexical", sourceID: eventproto.SourceId_GO_SERVER, sdkVersion: "1.10.0", expected: true},
		{desc: "another source has its own minimum", sourceID: eventproto.SourceId_NODE_SERVER, sdkVersion: "0.5.0", expected: true},
		{desc: "another source below its minimum", sourceID: eventproto.SourceId_NODE_SERVER, sdkVersion: "0.4.6", expected: false},
		{desc: "not a version", sourceID: eventproto.SourceId_GO_SERVER, sdkVersion: "latest", expected: false},
		{desc: "empty version", sourceID: eventproto.SourceId_GO_SERVER, sdkVersion: "", expected: false},
		{desc: "a wrapper is not listed", sourceID: eventproto.SourceId_OPEN_FEATURE_GO, sdkVersion: "9.9.9", expected: false},
	}
	for _, p := range patterns {
		assert.Equal(t, p.expected, holdoutCapableSDK(minVersions, p.sourceID, p.sdkVersion), p.desc)
	}

}

func TestHoldoutMinSDKVersions(t *testing.T) {
	t.Parallel()
	// Only the SDKs that have released holdout support are listed. A new entry
	// must come with the release that ships it, so the table is pinned here.
	assert.Equal(t,
		map[eventproto.SourceId]semver.Version{
			eventproto.SourceId_GO_SERVER: semver.MustParse("1.9.0"),
		},
		holdoutMinSDKVersions,
	)

	patterns := []struct {
		desc       string
		sourceID   eventproto.SourceId
		sdkVersion string
		expected   bool
	}{
		{desc: "go: first release with holdout", sourceID: eventproto.SourceId_GO_SERVER, sdkVersion: "1.9.0", expected: true},
		{desc: "go: later release", sourceID: eventproto.SourceId_GO_SERVER, sdkVersion: "1.10.0", expected: true},
		{desc: "go: last release without holdout", sourceID: eventproto.SourceId_GO_SERVER, sdkVersion: "1.8.0", expected: false},
		// A source that is not listed is not capable, whatever it reports.
		{desc: "node: not listed", sourceID: eventproto.SourceId_NODE_SERVER, sdkVersion: "999.0.0", expected: false},
		// An OpenFeature provider reports its own version, not the SDK's.
		{desc: "openfeature go", sourceID: eventproto.SourceId_OPEN_FEATURE_GO, sdkVersion: "999.0.0", expected: false},
		{desc: "openfeature node", sourceID: eventproto.SourceId_OPEN_FEATURE_NODE, sdkVersion: "999.0.0", expected: false},
	}
	for _, p := range patterns {
		assert.Equal(t, p.expected, holdoutCapableSDK(holdoutMinSDKVersions, p.sourceID, p.sdkVersion), p.desc)
	}
}

func TestCountHoldoutIncompatibleSDK(t *testing.T) {
	t.Parallel()
	const environmentID = "holdout-incompatible-unit"
	minVersions := map[eventproto.SourceId]semver.Version{
		eventproto.SourceId_GO_SERVER: semver.MustParse("1.9.0"),
	}
	count := func(sourceID eventproto.SourceId, sdkVersion string) float64 {
		return testutil.ToFloat64(
			holdoutIncompatibleSDKCounter.WithLabelValues(environmentID, sourceID.String(), sdkVersion),
		)
	}
	targeted := sdkHoldoutGroups(storedHoldoutGroup(environmentID, 50))
	untargeted := []*featureproto.HoldoutGroup{{Id: "hg-1", Status: featureproto.HoldoutGroup_RUNNING}}

	// Nothing is held out, so an old SDK contaminates nothing.
	countHoldoutIncompatibleSDK(nil, minVersions, environmentID, eventproto.SourceId_GO_SERVER, "1.8.0")
	countHoldoutIncompatibleSDK(untargeted, minVersions, environmentID, eventproto.SourceId_GO_SERVER, "1.8.0")
	assert.Equal(t, float64(0), count(eventproto.SourceId_GO_SERVER, "1.8.0"))

	// A capable SDK is not counted.
	countHoldoutIncompatibleSDK(targeted, minVersions, environmentID, eventproto.SourceId_GO_SERVER, "1.9.0")
	assert.Equal(t, float64(0), count(eventproto.SourceId_GO_SERVER, "1.9.0"))

	// An old SDK and an unlisted source are counted per source and version.
	countHoldoutIncompatibleSDK(targeted, minVersions, environmentID, eventproto.SourceId_GO_SERVER, "1.8.0")
	countHoldoutIncompatibleSDK(targeted, minVersions, environmentID, eventproto.SourceId_GO_SERVER, "1.8.0")
	countHoldoutIncompatibleSDK(targeted, minVersions, environmentID, eventproto.SourceId_GO_SERVER, "1.7.0")
	countHoldoutIncompatibleSDK(targeted, minVersions, environmentID, eventproto.SourceId_NODE_SERVER, "0.4.6")
	assert.Equal(t, float64(2), count(eventproto.SourceId_GO_SERVER, "1.8.0"))
	assert.Equal(t, float64(1), count(eventproto.SourceId_GO_SERVER, "1.7.0"))
	assert.Equal(t, float64(1), count(eventproto.SourceId_NODE_SERVER, "0.4.6"))
}

// holdoutDeliveryService returns a gateway whose in-memory caches hold the
// given flags and holdout snapshot for environmentID. A nil snapshotErr with a
// nil snapshot is an environment that never had a group.
func holdoutDeliveryService(
	t *testing.T,
	controller *gomock.Controller,
	environmentID string,
	features []*featureproto.Feature,
	snapshot *featureproto.HoldoutGroup,
	snapshotErr error,
) (*grpcGatewayService, context.Context) {
	t.Helper()
	gs := newGrpcGatewayServiceWithMock(t, controller)
	gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().
		Get(holdoutDeliveryAPIKey).Return(
		&accountproto.EnvironmentAPIKey{
			Environment: &environmentproto.EnvironmentV2{Id: environmentID},
			ApiKey: &accountproto.APIKey{
				Id:   holdoutDeliveryAPIKey,
				Role: accountproto.APIKey_SDK_SERVER,
			},
		}, nil)
	gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().
		Get(environmentID).Return(&featureproto.Features{Features: features}, nil)
	holdoutCache := cachev3mock.NewMockHoldoutGroupCache(controller)
	holdoutCache.EXPECT().Get(environmentID).Return(snapshot, snapshotErr)
	gs.holdoutGroupCache = holdoutCache
	// Cache reads run on a singleflight goroutine. A mock failure there ends
	// that goroutine only, so the deadline is what fails the test.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"authorization": []string{holdoutDeliveryAPIKey},
	})
	return gs, ctx
}

func TestGrpcGetFeatureFlagsHoldout(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)

	const environmentID = "holdout-feature-flags"
	now := time.Now()
	ago := func(d time.Duration) int64 { return now.Add(-d).Unix() }
	// Both flags were last edited long before any cursor below, outside the
	// grace window of the flag diff.
	newFlags := func() []*featureproto.Feature {
		return []*featureproto.Feature{
			{Id: "flag-a", Version: 1, Tags: []string{"web"}, UpdatedAt: ago(2 * time.Hour)},
			{Id: "flag-b", Version: 1, Tags: []string{"web"}, UpdatedAt: ago(2 * time.Hour)},
		}
	}
	noFlags := []*featureproto.Feature{}
	noArchived := make([]string, 0)
	flagsOnlyID := evaluation.GenerateFeaturesID(newFlags())
	groupID := func(updatedAt int64) string {
		return evaluation.GenerateFeaturesID(newFlags(), deliveredHoldoutGroup(environmentID, updatedAt))
	}
	groups := func(updatedAt int64) []*featureproto.HoldoutGroup {
		return []*featureproto.HoldoutGroup{deliveredHoldoutGroup(environmentID, updatedAt)}
	}
	require.NotEqual(t, flagsOnlyID, groupID(ago(time.Minute)))

	// A second target added in the same second as the previous revision.
	retargeted := storedHoldoutGroup(environmentID, ago(time.Minute))
	retargeted.Features = append(retargeted.Features, &featureproto.HoldoutGroupFeature{
		Id:                 "row-c",
		HoldoutGroupId:     "hg-1",
		FeatureId:          "flag-b",
		EnvironmentId:      environmentID,
		HoldoutVariationId: "variation-off",
		AddedAt:            140,
	})
	retargetedDelivered := deliveredHoldoutGroup(environmentID, ago(time.Minute))
	retargetedDelivered.Features = append(retargetedDelivered.Features, &featureproto.HoldoutGroupFeature{
		Id:                 "row-c",
		HoldoutGroupId:     "hg-1",
		FeatureId:          "flag-b",
		EnvironmentId:      environmentID,
		HoldoutVariationId: "variation-off",
		AddedAt:            140,
	})
	retargetedID := evaluation.GenerateFeaturesID(newFlags(), retargetedDelivered)
	require.NotEqual(t, groupID(ago(time.Minute)), retargetedID)

	// flag-a edited a minute ago: inside the grace window of the cursors below.
	editedFlags := newFlags()
	editedFlags[0].UpdatedAt = ago(time.Minute)
	editedFlagsID := func(holdoutGroups ...*featureproto.HoldoutGroup) string {
		return evaluation.GenerateFeaturesID(
			[]*featureproto.Feature{editedFlags[0], editedFlags[1]}, holdoutGroups...,
		)
	}

	// flag-a edited an hour ago. An SDK whose cursor is five minutes old and
	// that still holds the two-hour-old flag-a has missed that edit: it is
	// older than the cursor minus the grace window, so the flag diff does not
	// find it. This is what a poll answered by a stale replica leaves behind.
	missedFlags := newFlags()
	missedFlags[0].UpdatedAt = ago(time.Hour)
	missedFlagsID := func(holdoutGroups ...*featureproto.HoldoutGroup) string {
		return evaluation.GenerateFeaturesID(
			[]*featureproto.Feature{missedFlags[0], missedFlags[1]}, holdoutGroups...,
		)
	}
	require.NotEqual(t, flagsOnlyID, missedFlagsID())

	patterns := []struct {
		desc        string
		flags       []*featureproto.Feature
		snapshot    *featureproto.HoldoutGroup
		snapshotErr error
		tag         string
		ffID        string
		requestedAt int64
		expected    *gwproto.GetFeatureFlagsResponse
		// RequestedAt is stamped by the server instead of echoing the request.
		requestedAtNow bool
	}{
		{
			desc:        "first request: all flags and the serving group",
			flags:       newFlags(),
			snapshot:    storedHoldoutGroup(environmentID, ago(time.Hour)),
			ffID:        "",
			requestedAt: 0,
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         groupID(ago(time.Hour)),
				Features:               newFlags(),
				ArchivedFeatureFlagIds: noArchived,
				ForceUpdate:            true,
				HoldoutGroups:          groups(ago(time.Hour)),
			},
			requestedAtNow: true,
		},
		{
			desc:        "same id: nothing changed and the group is still listed",
			flags:       newFlags(),
			snapshot:    storedHoldoutGroup(environmentID, ago(time.Hour)),
			ffID:        groupID(ago(time.Hour)),
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         groupID(ago(time.Hour)),
				Features:               noFlags,
				ArchivedFeatureFlagIds: noArchived,
				RequestedAt:            ago(5 * time.Minute),
				HoldoutGroups:          groups(ago(time.Hour)),
			},
		},
		{
			desc:        "start: the group is sent with an empty flag diff",
			flags:       newFlags(),
			snapshot:    storedHoldoutGroup(environmentID, ago(time.Minute)),
			ffID:        flagsOnlyID,
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         groupID(ago(time.Minute)),
				Features:               noFlags,
				ArchivedFeatureFlagIds: noArchived,
				RequestedAt:            ago(5 * time.Minute),
				ForceUpdate:            false,
				HoldoutGroups:          groups(ago(time.Minute)),
			},
		},
		{
			// The SDK was last answered by a replica that had not seen the
			// group yet. Its flags part proves it holds the current flags, so
			// the age of the holdout revision does not matter.
			desc:        "start seen late: still an empty flag diff",
			flags:       newFlags(),
			snapshot:    storedHoldoutGroup(environmentID, ago(time.Hour)),
			ffID:        flagsOnlyID,
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         groupID(ago(time.Hour)),
				Features:               noFlags,
				ArchivedFeatureFlagIds: noArchived,
				RequestedAt:            ago(5 * time.Minute),
				ForceUpdate:            false,
				HoldoutGroups:          groups(ago(time.Hour)),
			},
		},
		{
			desc:        "stop: an empty list with an empty flag diff",
			flags:       newFlags(),
			snapshot:    evaluation.ClearanceSnapshot(ago(time.Minute)),
			ffID:        groupID(ago(time.Hour)),
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         flagsOnlyID,
				Features:               noFlags,
				ArchivedFeatureFlagIds: noArchived,
				RequestedAt:            ago(5 * time.Minute),
				ForceUpdate:            false,
				HoldoutGroups:          nil,
			},
		},
		{
			desc:        "target change in the same second: the new targets with an empty flag diff",
			flags:       newFlags(),
			snapshot:    retargeted,
			ffID:        groupID(ago(time.Minute)),
			requestedAt: ago(30 * time.Second),
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         retargetedID,
				Features:               noFlags,
				ArchivedFeatureFlagIds: noArchived,
				RequestedAt:            ago(30 * time.Second),
				ForceUpdate:            false,
				HoldoutGroups:          []*featureproto.HoldoutGroup{retargetedDelivered},
			},
		},
		{
			// The SDK already holds the edited flag, so the grace window of
			// the flag diff must not send it again.
			desc:        "start after a recent flag edit the SDK holds: an empty flag diff",
			flags:       editedFlags,
			snapshot:    storedHoldoutGroup(environmentID, ago(30*time.Second)),
			ffID:        editedFlagsID(),
			requestedAt: ago(45 * time.Second),
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         editedFlagsID(deliveredHoldoutGroup(environmentID, ago(30*time.Second))),
				Features:               noFlags,
				ArchivedFeatureFlagIds: noArchived,
				RequestedAt:            ago(45 * time.Second),
				ForceUpdate:            false,
				HoldoutGroups:          groups(ago(30 * time.Second)),
			},
		},
		{
			desc:        "flag change: the diff carries the group too",
			flags:       editedFlags,
			snapshot:    storedHoldoutGroup(environmentID, ago(time.Hour)),
			ffID:        groupID(ago(time.Hour)),
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         editedFlagsID(deliveredHoldoutGroup(environmentID, ago(time.Hour))),
				Features:               []*featureproto.Feature{editedFlags[0]},
				ArchivedFeatureFlagIds: noArchived,
				ForceUpdate:            false,
				HoldoutGroups:          groups(ago(time.Hour)),
			},
			requestedAtNow: true,
		},
		{
			desc:        "flag change and start together: the diff carries the group",
			flags:       editedFlags,
			snapshot:    storedHoldoutGroup(environmentID, ago(time.Minute)),
			ffID:        flagsOnlyID,
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         editedFlagsID(deliveredHoldoutGroup(environmentID, ago(time.Minute))),
				Features:               []*featureproto.Feature{editedFlags[0]},
				ArchivedFeatureFlagIds: noArchived,
				ForceUpdate:            false,
				HoldoutGroups:          groups(ago(time.Minute)),
			},
			requestedAtNow: true,
		},
		{
			// A recent holdout change must not hide a flag update the SDK
			// missed: an empty diff here would hand the SDK the complete ID
			// and the missed edit would never be sent.
			desc:        "missed flag update and start together: full resend with the group",
			flags:       missedFlags,
			snapshot:    storedHoldoutGroup(environmentID, ago(time.Minute)),
			ffID:        flagsOnlyID,
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         missedFlagsID(deliveredHoldoutGroup(environmentID, ago(time.Minute))),
				Features:               []*featureproto.Feature{missedFlags[0], missedFlags[1]},
				ArchivedFeatureFlagIds: noArchived,
				ForceUpdate:            true,
				HoldoutGroups:          groups(ago(time.Minute)),
			},
			requestedAtNow: true,
		},
		{
			desc:        "missed flag update while the SDK holds the serving group: full resend",
			flags:       missedFlags,
			snapshot:    storedHoldoutGroup(environmentID, ago(time.Minute)),
			ffID:        groupID(ago(time.Minute)),
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         missedFlagsID(deliveredHoldoutGroup(environmentID, ago(time.Minute))),
				Features:               []*featureproto.Feature{missedFlags[0], missedFlags[1]},
				ArchivedFeatureFlagIds: noArchived,
				ForceUpdate:            true,
				HoldoutGroups:          groups(ago(time.Minute)),
			},
			requestedAtNow: true,
		},
		{
			desc:        "missed flag update and stop together: full resend without a group",
			flags:       missedFlags,
			snapshot:    evaluation.ClearanceSnapshot(ago(time.Minute)),
			ffID:        groupID(ago(time.Hour)),
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         missedFlagsID(),
				Features:               []*featureproto.Feature{missedFlags[0], missedFlags[1]},
				ArchivedFeatureFlagIds: noArchived,
				ForceUpdate:            true,
				HoldoutGroups:          nil,
			},
			requestedAtNow: true,
		},
		{
			desc:        "no group ever: an unexplained id mismatch still forces a full resend",
			flags:       newFlags(),
			snapshot:    evaluation.ClearanceSnapshot(0),
			ffID:        "unknown-id",
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         flagsOnlyID,
				Features:               newFlags(),
				ArchivedFeatureFlagIds: noArchived,
				ForceUpdate:            true,
				HoldoutGroups:          nil,
			},
			requestedAtNow: true,
		},
		{
			desc:        "unreadable snapshot: the request is answered without the group",
			flags:       newFlags(),
			snapshotErr: errors.New("memory cache failure"),
			ffID:        groupID(ago(time.Hour)),
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         flagsOnlyID,
				Features:               noFlags,
				ArchivedFeatureFlagIds: noArchived,
				RequestedAt:            ago(5 * time.Minute),
				ForceUpdate:            false,
				HoldoutGroups:          nil,
			},
		},
		{
			desc:        "unreadable snapshot: an SDK without a group is left alone",
			flags:       newFlags(),
			snapshotErr: errors.New("memory cache failure"),
			ffID:        flagsOnlyID,
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         flagsOnlyID,
				Features:               noFlags,
				ArchivedFeatureFlagIds: noArchived,
				RequestedAt:            ago(5 * time.Minute),
				HoldoutGroups:          nil,
			},
		},
		{
			desc:        "unreadable snapshot: a first request still gets all flags",
			flags:       newFlags(),
			snapshotErr: errors.New("memory cache failure"),
			ffID:        "",
			requestedAt: 0,
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         flagsOnlyID,
				Features:               newFlags(),
				ArchivedFeatureFlagIds: noArchived,
				ForceUpdate:            true,
				HoldoutGroups:          nil,
			},
			requestedAtNow: true,
		},
		{
			desc:        "no flag for the tag: the group is still listed",
			flags:       newFlags(),
			snapshot:    storedHoldoutGroup(environmentID, ago(time.Hour)),
			tag:         "ios",
			ffID:        "",
			requestedAt: 0,
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         "",
				Features:               noFlags,
				ArchivedFeatureFlagIds: noArchived,
				HoldoutGroups:          groups(ago(time.Hour)),
			},
			requestedAtNow: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs, ctx := holdoutDeliveryService(t, controller, environmentID, p.flags, p.snapshot, p.snapshotErr)
			var before proto.Message
			if p.snapshot != nil {
				before = proto.Clone(p.snapshot)
			}

			actual, err := gs.GetFeatureFlags(ctx, &gwproto.GetFeatureFlagsRequest{
				Tag:            p.tag,
				FeatureFlagsId: p.ffID,
				RequestedAt:    p.requestedAt,
				SourceId:       eventproto.SourceId_GO_SERVER,
				SdkVersion:     "1.8.0",
			})

			require.NoError(t, err)
			require.NotNil(t, actual)
			if p.requestedAtNow {
				assert.InDelta(t, time.Now().Unix(), actual.RequestedAt, 5)
				p.expected.RequestedAt = actual.RequestedAt
			}
			assert.True(t, proto.Equal(p.expected, actual), "expected %v\nactual   %v", p.expected, actual)
			assert.Len(t, actual.HoldoutGroups, len(p.expected.HoldoutGroups))
			if p.snapshot != nil {
				assert.True(t, proto.Equal(before, p.snapshot), "the cached snapshot must not be modified")
			}
		})
	}
}

func TestGrpcGetFeatureFlagsCountsIncompatibleSDK(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)

	const environmentID = "holdout-feature-flags-compat"
	const sdkVersion = "1.8.0"
	metric := holdoutIncompatibleSDKCounter.WithLabelValues(
		environmentID, eventproto.SourceId_GO_SERVER.String(), sdkVersion,
	)
	flags := func() []*featureproto.Feature {
		return []*featureproto.Feature{{Id: "flag-a", Version: 1, UpdatedAt: 10}}
	}
	request := func() *gwproto.GetFeatureFlagsRequest {
		return &gwproto.GetFeatureFlagsRequest{
			SourceId:   eventproto.SourceId_GO_SERVER,
			SdkVersion: sdkVersion,
		}
	}

	// No group is serving.
	gs, ctx := holdoutDeliveryService(t, controller, environmentID, flags(), evaluation.ClearanceSnapshot(100), nil)
	_, err := gs.GetFeatureFlags(ctx, request())
	require.NoError(t, err)
	assert.Equal(t, float64(0), testutil.ToFloat64(metric))

	// A group is serving, but it holds no flag out yet.
	untargeted := storedHoldoutGroup(environmentID, 100)
	untargeted.Features = nil
	gs, ctx = holdoutDeliveryService(t, controller, environmentID, flags(), untargeted, nil)
	_, err = gs.GetFeatureFlags(ctx, request())
	require.NoError(t, err)
	assert.Equal(t, float64(0), testutil.ToFloat64(metric))

	// A group with targets is serving and the SDK is older than the first
	// release that applies it. The request is counted and still answered in
	// full.
	gs, ctx = holdoutDeliveryService(t, controller, environmentID, flags(), storedHoldoutGroup(environmentID, 100), nil)
	resp, err := gs.GetFeatureFlags(ctx, request())
	require.NoError(t, err)
	assert.Len(t, resp.Features, 1)
	assert.Len(t, resp.HoldoutGroups, 1)
	assert.Equal(t, float64(1), testutil.ToFloat64(metric))

	// The first release that applies the group is not counted, and gets the
	// same answer.
	const capableSDKVersion = "1.9.0"
	capableMetric := holdoutIncompatibleSDKCounter.WithLabelValues(
		environmentID, eventproto.SourceId_GO_SERVER.String(), capableSDKVersion,
	)
	gs, ctx = holdoutDeliveryService(t, controller, environmentID, flags(), storedHoldoutGroup(environmentID, 100), nil)
	resp, err = gs.GetFeatureFlags(ctx, &gwproto.GetFeatureFlagsRequest{
		SourceId:   eventproto.SourceId_GO_SERVER,
		SdkVersion: capableSDKVersion,
	})
	require.NoError(t, err)
	assert.Len(t, resp.Features, 1)
	assert.Len(t, resp.HoldoutGroups, 1)
	assert.Equal(t, float64(0), testutil.ToFloat64(capableMetric))
	assert.Equal(t, float64(1), testutil.ToFloat64(metric))

}

func TestGrpcGetSegmentUsersHoldout(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)

	const environmentID = "holdout-segment-users"
	now := time.Now()
	ago := func(d time.Duration) int64 { return now.Add(-d).Unix() }

	plainFlag := &featureproto.Feature{Id: "flag-a", Version: 1, UpdatedAt: ago(2 * time.Hour)}
	segmentFlag := func(segmentID string) *featureproto.Feature {
		return &featureproto.Feature{
			Id:        "flag-" + segmentID,
			Version:   1,
			UpdatedAt: ago(2 * time.Hour),
			Rules: []*featureproto.Rule{{
				Id: "rule",
				Clauses: []*featureproto.Clause{{
					Id:       "clause",
					Operator: featureproto.Clause_SEGMENT,
					Values:   []string{segmentID},
				}},
			}},
		}
	}
	segmentGroup := func(segmentID string) *featureproto.HoldoutGroup {
		group := storedHoldoutGroup(environmentID, ago(time.Minute))
		group.AssignmentType = featureproto.HoldoutGroup_SEGMENT
		group.SegmentId = segmentID
		return group
	}
	// Both segments were last edited long before any cursor below.
	holdoutSegment := &featureproto.SegmentUsers{
		SegmentId: "segment-holdout",
		Users:     []*featureproto.SegmentUser{{SegmentId: "segment-holdout", UserId: "user-1"}},
		UpdatedAt: ago(24 * time.Hour),
	}
	flagSegment := &featureproto.SegmentUsers{
		SegmentId: "segment-flag",
		Users:     []*featureproto.SegmentUser{{SegmentId: "segment-flag", UserId: "user-2"}},
		UpdatedAt: ago(24 * time.Hour),
	}
	editedHoldoutSegment := &featureproto.SegmentUsers{
		SegmentId: "segment-holdout",
		Users:     []*featureproto.SegmentUser{{SegmentId: "segment-holdout", UserId: "user-3"}},
		UpdatedAt: ago(time.Minute),
	}
	noSegments := make([]*featureproto.SegmentUsers, 0)
	noDeleted := make([]string, 0)

	patterns := []struct {
		desc     string
		flags    []*featureproto.Feature
		snapshot *featureproto.HoldoutGroup
		// cached lists the segments the gateway is expected to read, once each.
		cached      []*featureproto.SegmentUsers
		segmentIDs  []string
		requestedAt int64
		expected    *gwproto.GetSegmentUsersResponse
		// RequestedAt is stamped by the server instead of echoing the request.
		requestedAtNow bool
	}{
		{
			desc:        "the holdout segment is delivered without a flag rule referencing it",
			flags:       []*featureproto.Feature{plainFlag},
			snapshot:    segmentGroup("segment-holdout"),
			cached:      []*featureproto.SegmentUsers{holdoutSegment},
			segmentIDs:  nil,
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetSegmentUsersResponse{
				SegmentUsers:      []*featureproto.SegmentUsers{holdoutSegment},
				DeletedSegmentIds: noDeleted,
				ForceUpdate:       false,
			},
			requestedAtNow: true,
		},
		{
			desc:        "the holdout segment is delivered next to the flag segments the SDK holds",
			flags:       []*featureproto.Feature{plainFlag, segmentFlag("segment-flag")},
			snapshot:    segmentGroup("segment-holdout"),
			cached:      []*featureproto.SegmentUsers{flagSegment, holdoutSegment},
			segmentIDs:  []string{"segment-flag"},
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetSegmentUsersResponse{
				SegmentUsers:      []*featureproto.SegmentUsers{holdoutSegment},
				DeletedSegmentIds: noDeleted,
				ForceUpdate:       false,
			},
			requestedAtNow: true,
		},
		{
			desc:        "an SDK that already holds the holdout segment gets nothing",
			flags:       []*featureproto.Feature{plainFlag},
			snapshot:    segmentGroup("segment-holdout"),
			cached:      []*featureproto.SegmentUsers{holdoutSegment},
			segmentIDs:  []string{"segment-holdout"},
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetSegmentUsersResponse{
				SegmentUsers:      noSegments,
				DeletedSegmentIds: noDeleted,
				RequestedAt:       ago(5 * time.Minute),
				ForceUpdate:       false,
			},
		},
		{
			desc:        "an edit of the holdout segment reaches an SDK that holds it",
			flags:       []*featureproto.Feature{plainFlag},
			snapshot:    segmentGroup("segment-holdout"),
			cached:      []*featureproto.SegmentUsers{editedHoldoutSegment},
			segmentIDs:  []string{"segment-holdout"},
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetSegmentUsersResponse{
				SegmentUsers:      []*featureproto.SegmentUsers{editedHoldoutSegment},
				DeletedSegmentIds: noDeleted,
				ForceUpdate:       false,
			},
			requestedAtNow: true,
		},
		{
			desc:        "a full resend includes the holdout segment",
			flags:       []*featureproto.Feature{plainFlag, segmentFlag("segment-flag")},
			snapshot:    segmentGroup("segment-holdout"),
			cached:      []*featureproto.SegmentUsers{flagSegment, holdoutSegment},
			segmentIDs:  nil,
			requestedAt: 0,
			expected: &gwproto.GetSegmentUsersResponse{
				SegmentUsers:      []*featureproto.SegmentUsers{flagSegment, holdoutSegment},
				DeletedSegmentIds: noDeleted,
				ForceUpdate:       true,
			},
			requestedAtNow: true,
		},
		{
			desc:        "a segment used by both a flag and the group is read and sent once",
			flags:       []*featureproto.Feature{segmentFlag("segment-holdout")},
			snapshot:    segmentGroup("segment-holdout"),
			cached:      []*featureproto.SegmentUsers{holdoutSegment},
			segmentIDs:  nil,
			requestedAt: 0,
			expected: &gwproto.GetSegmentUsersResponse{
				SegmentUsers:      []*featureproto.SegmentUsers{holdoutSegment},
				DeletedSegmentIds: noDeleted,
				ForceUpdate:       true,
			},
			requestedAtNow: true,
		},
		{
			desc:        "after the group stops its segment is reported as deleted",
			flags:       []*featureproto.Feature{plainFlag, segmentFlag("segment-flag")},
			snapshot:    evaluation.ClearanceSnapshot(ago(time.Minute)),
			cached:      []*featureproto.SegmentUsers{flagSegment},
			segmentIDs:  []string{"segment-flag", "segment-holdout"},
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetSegmentUsersResponse{
				SegmentUsers:      noSegments,
				DeletedSegmentIds: []string{"segment-holdout"},
				ForceUpdate:       false,
			},
			requestedAtNow: true,
		},
		{
			desc:        "a hash group adds no segment",
			flags:       []*featureproto.Feature{plainFlag},
			snapshot:    storedHoldoutGroup(environmentID, ago(time.Minute)),
			cached:      nil,
			segmentIDs:  nil,
			requestedAt: ago(5 * time.Minute),
			expected: &gwproto.GetSegmentUsersResponse{
				SegmentUsers:      noSegments,
				DeletedSegmentIds: noDeleted,
				ForceUpdate:       true,
			},
			requestedAtNow: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs, ctx := holdoutDeliveryService(t, controller, environmentID, p.flags, p.snapshot, nil)
			for _, segment := range p.cached {
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().
					Get(segment.SegmentId, environmentID).Return(segment, nil).Times(1)
			}

			actual, err := gs.GetSegmentUsers(ctx, &gwproto.GetSegmentUsersRequest{
				SegmentIds:  p.segmentIDs,
				RequestedAt: p.requestedAt,
				SourceId:    eventproto.SourceId_GO_SERVER,
				SdkVersion:  "1.8.0",
			})

			require.NoError(t, err)
			require.NotNil(t, actual)
			if p.requestedAtNow {
				assert.InDelta(t, time.Now().Unix(), actual.RequestedAt, 5)
				p.expected.RequestedAt = actual.RequestedAt
			}
			assert.True(t, proto.Equal(p.expected, actual), "expected %v\nactual   %v", p.expected, actual)
			assert.Len(t, actual.SegmentUsers, len(p.expected.SegmentUsers))
			assert.Len(t, actual.DeletedSegmentIds, len(p.expected.DeletedSegmentIds))
		})
	}
}
