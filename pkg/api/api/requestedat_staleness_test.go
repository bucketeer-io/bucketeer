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

// Tests for the RequestedAt staleness bug and its fix.
//
// Affected APIs:
//   - GetFeatureFlags: uses FeatureFlagsId + RequestedAt for None/Diff/All responses
//   - GetSegmentUsers: uses SegmentIds + RequestedAt for None/Diff/All responses
//
// NOT affected:
//   - GetEvaluation / GetEvaluations: client SDK APIs that evaluate per-request
//     using evaluatedAt, not the RequestedAt/Diff caching mechanism.

package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// ---------------------------------------------------------------------------
// GetFeatureFlags simulation
// ---------------------------------------------------------------------------

// featureFlagsResponse mirrors the fields of GetFeatureFlagsResponse
// relevant to the cache sync logic.
type featureFlagsResponse struct {
	FeatureFlagsId         string
	Features               []*featureproto.Feature
	ArchivedFeatureFlagIds []string
	RequestedAt            int64
	ForceUpdate            bool
	ResponseType           string // "none", "all", "diff" — for logging only
}

// simulateGetFeatureFlagsBuggy replicates the BUGGY logic (before fix) from
// grpcGatewayService.GetFeatureFlags where "None" responses advance
// requestedAt to now and "Diff" uses secondsForAdjustment.
func simulateGetFeatureFlagsBuggy(
	serverFeatures []*featureproto.Feature,
	reqFeatureFlagsId string,
	reqRequestedAt int64,
	now int64,
) featureFlagsResponse {
	if len(serverFeatures) == 0 {
		return featureFlagsResponse{
			Features:               []*featureproto.Feature{},
			ArchivedFeatureFlagIds: []string{},
			RequestedAt:            now,
			ResponseType:           "empty",
		}
	}

	ffID := evaluation.GenerateFeaturesID(serverFeatures)

	if reqFeatureFlagsId == ffID {
		return featureFlagsResponse{
			FeatureFlagsId:         ffID,
			Features:               []*featureproto.Feature{},
			ArchivedFeatureFlagIds: []string{},
			RequestedAt:            now, // BUG: advances requestedAt
			ResponseType:           "none",
		}
	}

	if reqFeatureFlagsId == "" || reqRequestedAt < now-secondsToReturnAllFlags {
		return featureFlagsResponse{
			FeatureFlagsId:         ffID,
			Features:               serverFeatures,
			ArchivedFeatureFlagIds: []string{},
			RequestedAt:            now,
			ForceUpdate:            true,
			ResponseType:           "all",
		}
	}

	const secondsForAdjustment = 10 // old constant, removed from production
	adjustedRequestedAt := reqRequestedAt - secondsForAdjustment
	updatedFeatures := make([]*featureproto.Feature, 0)
	for _, feature := range serverFeatures {
		if feature.UpdatedAt > adjustedRequestedAt {
			updatedFeatures = append(updatedFeatures, feature)
		}
	}

	return featureFlagsResponse{
		FeatureFlagsId:         ffID,
		Features:               updatedFeatures,
		ArchivedFeatureFlagIds: []string{},
		RequestedAt:            now,
		ResponseType:           "diff",
	}
}

// simulateGetFeatureFlagsFixed replicates the production logic from
// grpcGatewayService.GetFeatureFlags after the fix.
func simulateGetFeatureFlagsFixed(
	serverFeatures []*featureproto.Feature,
	reqFeatureFlagsId string,
	reqRequestedAt int64,
	now int64,
) featureFlagsResponse {
	if len(serverFeatures) == 0 {
		return featureFlagsResponse{
			Features:               []*featureproto.Feature{},
			ArchivedFeatureFlagIds: []string{},
			RequestedAt:            now,
			ResponseType:           "empty",
		}
	}

	ffID := evaluation.GenerateFeaturesID(serverFeatures)

	// None: preserve requestedAt (clamped to now for future values)
	if reqFeatureFlagsId == ffID {
		return featureFlagsResponse{
			FeatureFlagsId:         ffID,
			Features:               []*featureproto.Feature{},
			ArchivedFeatureFlagIds: []string{},
			RequestedAt:            min(reqRequestedAt, now),
			ResponseType:           "none",
		}
	}

	// All: first request, very old cache, or future requestedAt (clock skew)
	if reqFeatureFlagsId == "" ||
		reqRequestedAt < now-secondsToReturnAllFlags ||
		reqRequestedAt > now {
		return featureFlagsResponse{
			FeatureFlagsId:         ffID,
			Features:               serverFeatures,
			ArchivedFeatureFlagIds: []string{},
			RequestedAt:            now,
			ForceUpdate:            true,
			ResponseType:           "all",
		}
	}

	// Diff: reqRequestedAt is guaranteed within [now-30days, now]
	updatedFeatures := make([]*featureproto.Feature, 0)
	for _, feature := range serverFeatures {
		if feature.UpdatedAt >= reqRequestedAt {
			updatedFeatures = append(updatedFeatures, feature)
		}
	}

	return featureFlagsResponse{
		FeatureFlagsId:         ffID,
		Features:               updatedFeatures,
		ArchivedFeatureFlagIds: []string{},
		RequestedAt:            now,
		ResponseType:           "diff",
	}
}

// featureFlagsSDKCache represents the Go SDK's in-memory state for feature flags.
type featureFlagsSDKCache struct {
	featureFlagsID string
	requestedAt    int64
	features       map[string]*featureproto.Feature
}

func (c *featureFlagsSDKCache) applyResponse(resp featureFlagsResponse) {
	if resp.ForceUpdate {
		c.features = make(map[string]*featureproto.Feature)
		for _, f := range resp.Features {
			c.features[f.Id] = f
		}
	} else {
		for _, f := range resp.Features {
			c.features[f.Id] = f
		}
		for _, archivedID := range resp.ArchivedFeatureFlagIds {
			delete(c.features, archivedID)
		}
	}
	c.featureFlagsID = resp.FeatureFlagsId
	c.requestedAt = resp.RequestedAt
}

type featureFlagsServerFunc func(
	serverFeatures []*featureproto.Feature,
	reqFeatureFlagsId string,
	reqRequestedAt int64,
	now int64,
) featureFlagsResponse

// ---------------------------------------------------------------------------
// GetSegmentUsers simulation
// ---------------------------------------------------------------------------

// segmentUsersResponse mirrors the fields of GetSegmentUsersResponse
// relevant to the cache sync logic.
type segmentUsersResponse struct {
	SegmentUsers      []*featureproto.SegmentUsers
	DeletedSegmentIds []string
	RequestedAt       int64
	ForceUpdate       bool
	ResponseType      string
}

// simulateGetSegmentUsersBuggy replicates the BUGGY GetSegmentUsers logic
// where the response always advances requestedAt to now and uses
// secondsForAdjustment in the diff filter.
func simulateGetSegmentUsersBuggy(
	serverSegments []*featureproto.SegmentUsers,
	reqSegmentIds []string,
	reqRequestedAt int64,
	now int64,
) segmentUsersResponse {
	if reqRequestedAt < now-secondsToReturnAllFlags {
		return segmentUsersResponse{
			SegmentUsers:      serverSegments,
			DeletedSegmentIds: []string{},
			RequestedAt:       now,
			ForceUpdate:       true,
			ResponseType:      "all",
		}
	}

	serverIDs := make(map[string]bool)
	for _, s := range serverSegments {
		serverIDs[s.SegmentId] = true
	}
	deletedIDs := make([]string, 0)
	for _, id := range reqSegmentIds {
		if !serverIDs[id] {
			deletedIDs = append(deletedIDs, id)
		}
	}

	const secondsForAdjustment = 10
	adjustedRequestedAt := reqRequestedAt - secondsForAdjustment
	updated := make([]*featureproto.SegmentUsers, 0)
	for _, su := range serverSegments {
		if su.UpdatedAt > adjustedRequestedAt {
			updated = append(updated, su)
		}
	}

	return segmentUsersResponse{
		SegmentUsers:      updated,
		DeletedSegmentIds: deletedIDs,
		RequestedAt:       now, // BUG: always advances
		ResponseType:      responseType(updated, deletedIDs),
	}
}

// simulateGetSegmentUsersFixed replicates the production logic from
// grpcGatewayService.GetSegmentUsers after the fix.
func simulateGetSegmentUsersFixed(
	serverSegments []*featureproto.SegmentUsers,
	reqSegmentIds []string,
	reqRequestedAt int64,
	now int64,
) segmentUsersResponse {
	// All: very old cache or future requestedAt (clock skew)
	if reqRequestedAt < now-secondsToReturnAllFlags || reqRequestedAt > now {
		return segmentUsersResponse{
			SegmentUsers:      serverSegments,
			DeletedSegmentIds: []string{},
			RequestedAt:       now,
			ForceUpdate:       true,
			ResponseType:      "all",
		}
	}

	serverIDs := make(map[string]bool)
	for _, s := range serverSegments {
		serverIDs[s.SegmentId] = true
	}
	deletedIDs := make([]string, 0)
	for _, id := range reqSegmentIds {
		if !serverIDs[id] {
			deletedIDs = append(deletedIDs, id)
		}
	}

	// Diff: reqRequestedAt is guaranteed within [now-30days, now]
	updated := make([]*featureproto.SegmentUsers, 0)
	for _, su := range serverSegments {
		if su.UpdatedAt >= reqRequestedAt {
			updated = append(updated, su)
		}
	}

	// None: preserve requestedAt when nothing changed
	if len(updated) == 0 && len(deletedIDs) == 0 {
		return segmentUsersResponse{
			SegmentUsers:      updated,
			DeletedSegmentIds: deletedIDs,
			RequestedAt:       reqRequestedAt,
			ResponseType:      "none",
		}
	}
	return segmentUsersResponse{
		SegmentUsers:      updated,
		DeletedSegmentIds: deletedIDs,
		RequestedAt:       now,
		ResponseType:      responseType(updated, deletedIDs),
	}
}

func responseType(updated []*featureproto.SegmentUsers, deleted []string) string {
	if len(updated) == 0 && len(deleted) == 0 {
		return "none"
	}
	return "diff"
}

// segmentUsersSDKCache represents the Go SDK's in-memory state for segment users.
type segmentUsersSDKCache struct {
	requestedAt int64
	segmentIDs  []string
	segments    map[string]*featureproto.SegmentUsers
}

func (c *segmentUsersSDKCache) applyResponse(resp segmentUsersResponse) {
	if resp.ForceUpdate {
		c.segments = make(map[string]*featureproto.SegmentUsers)
		c.segmentIDs = nil
		for _, su := range resp.SegmentUsers {
			c.segments[su.SegmentId] = su
			c.segmentIDs = append(c.segmentIDs, su.SegmentId)
		}
	} else {
		for _, su := range resp.SegmentUsers {
			c.segments[su.SegmentId] = su
		}
		for _, id := range resp.DeletedSegmentIds {
			delete(c.segments, id)
			for i, sid := range c.segmentIDs {
				if sid == id {
					c.segmentIDs = append(c.segmentIDs[:i], c.segmentIDs[i+1:]...)
					break
				}
			}
		}
	}
	c.requestedAt = resp.RequestedAt
}

type segmentUsersServerFunc func(
	serverSegments []*featureproto.SegmentUsers,
	reqSegmentIds []string,
	reqRequestedAt int64,
	now int64,
) segmentUsersResponse

// ---------------------------------------------------------------------------
// GetFeatureFlags tests
// ---------------------------------------------------------------------------

// TestGetFeatureFlagsStaleCacheBuggyVsFixed runs the same scenario against
// both the buggy and fixed server logic to demonstrate the difference.
func TestGetFeatureFlagsStaleCacheBuggyVsFixed(t *testing.T) {
	t.Parallel()

	baseTime := int64(1710000000)

	flagEnabled := &featureproto.Feature{
		Id:        "flag-1",
		Version:   1,
		Enabled:   true,
		UpdatedAt: baseTime - 3600,
	}
	flagDisabled := &featureproto.Feature{
		Id:        "flag-1",
		Version:   2,
		Enabled:   false,
		UpdatedAt: baseTime + 125,
	}

	patterns := []struct {
		desc            string
		serverFn        featureFlagsServerFunc
		expectStale     bool
		expectPermanent bool
	}{
		{
			desc:            "buggy: None responses advance requestedAt, causing permanent staleness",
			serverFn:        simulateGetFeatureFlagsBuggy,
			expectStale:     true,
			expectPermanent: true,
		},
		{
			desc:            "fixed: None responses preserve requestedAt, SDK picks up changes",
			serverFn:        simulateGetFeatureFlagsFixed,
			expectStale:     false,
			expectPermanent: false,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			serverFeatures := []*featureproto.Feature{flagEnabled}
			sdk := &featureFlagsSDKCache{features: make(map[string]*featureproto.Feature)}

			// Phase 1: Initial sync (T=0) → "All" response
			resp := p.serverFn(serverFeatures, sdk.featureFlagsID, sdk.requestedAt, baseTime)
			require.Equal(t, "all", resp.ResponseType)
			sdk.applyResponse(resp)
			require.True(t, sdk.features["flag-1"].Enabled)

			// Phase 2: Several "None" polls (T=60, T=120)
			for _, offset := range []int64{60, 120} {
				resp = p.serverFn(serverFeatures, sdk.featureFlagsID, sdk.requestedAt, baseTime+offset)
				require.Equal(t, "none", resp.ResponseType)
				sdk.applyResponse(resp)
			}

			// Phase 3: Flag updated at T=125, but SDK polls at T=180 before Redis refreshes
			resp = p.serverFn(serverFeatures, sdk.featureFlagsID, sdk.requestedAt, baseTime+180)
			require.Equal(t, "none", resp.ResponseType)
			sdk.applyResponse(resp)

			// Phase 4: Redis refreshes at T=181
			serverFeatures = []*featureproto.Feature{flagDisabled}

			// Phase 5: SDK polls at T=240 → ffID mismatch → "Diff"
			resp = p.serverFn(serverFeatures, sdk.featureFlagsID, sdk.requestedAt, baseTime+240)
			require.Equal(t, "diff", resp.ResponseType)
			sdk.applyResponse(resp)

			if p.expectStale {
				assert.Empty(t, resp.Features,
					"buggy: Diff missed the updated flag")
				assert.True(t, sdk.features["flag-1"].Enabled,
					"buggy: SDK still has old Enabled=true")
			} else {
				assert.Len(t, resp.Features, 1,
					"fixed: Diff includes the updated flag")
				assert.False(t, sdk.features["flag-1"].Enabled,
					"fixed: SDK correctly has Enabled=false")
			}

			// Phase 6: Check permanence — subsequent polls
			for _, offset := range []int64{300, 360, 420} {
				resp = p.serverFn(serverFeatures, sdk.featureFlagsID, sdk.requestedAt, baseTime+offset)
				assert.Equal(t, "none", resp.ResponseType,
					fmt.Sprintf("T=%d: subsequent polls should be 'none'", offset))
			}

			if p.expectPermanent {
				assert.True(t, sdk.features["flag-1"].Enabled,
					"buggy: permanently stale — never picks up the change")
			} else {
				assert.False(t, sdk.features["flag-1"].Enabled,
					"fixed: correctly reflects the disabled flag")
			}
		})
	}
}

// TestGetFeatureFlagsStaleCacheWithinSinglePollingInterval demonstrates that
// even without a Redis refresh delay, the buggy code can miss updates because
// 24h of "None" responses pushed requestedAt too far ahead.
func TestGetFeatureFlagsStaleCacheWithinSinglePollingInterval(t *testing.T) {
	t.Parallel()

	baseTime := int64(1710000000)
	pollingInterval := int64(60)

	flag := &featureproto.Feature{
		Id:        "flag-1",
		Version:   1,
		Enabled:   true,
		UpdatedAt: baseTime - 3600,
	}
	flagUpdated := &featureproto.Feature{
		Id:        "flag-1",
		Version:   2,
		Enabled:   false,
		UpdatedAt: baseTime + 5,
	}

	patterns := []struct {
		desc     string
		serverFn featureFlagsServerFunc
		wantHit  bool
	}{
		{
			desc:     "buggy: 24h of None polls causes missed update",
			serverFn: simulateGetFeatureFlagsBuggy,
			wantHit:  false,
		},
		{
			desc:     "fixed: requestedAt stays anchored, update is detected",
			serverFn: simulateGetFeatureFlagsFixed,
			wantHit:  true,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			serverFeatures := []*featureproto.Feature{flag}
			sdk := &featureFlagsSDKCache{features: make(map[string]*featureproto.Feature)}

			resp := p.serverFn(serverFeatures, sdk.featureFlagsID, sdk.requestedAt, baseTime)
			sdk.applyResponse(resp)

			for i := int64(1); i <= 1440; i++ {
				now := baseTime + i*pollingInterval
				resp = p.serverFn(serverFeatures, sdk.featureFlagsID, sdk.requestedAt, now)
				require.Equal(t, "none", resp.ResponseType)
				sdk.applyResponse(resp)
			}

			serverFeatures = []*featureproto.Feature{flagUpdated}

			now := baseTime + 1441*pollingInterval
			resp = p.serverFn(serverFeatures, sdk.featureFlagsID, sdk.requestedAt, now)
			require.Equal(t, "diff", resp.ResponseType)

			if p.wantHit {
				assert.Len(t, resp.Features, 1,
					"fixed: updated flag is included in diff")
				sdk.applyResponse(resp)
				assert.False(t, sdk.features["flag-1"].Enabled,
					"fixed: SDK correctly reflects the disabled flag")
			} else {
				assert.Empty(t, resp.Features,
					"buggy: updated flag is missed by diff filter")
				sdk.applyResponse(resp)
				assert.True(t, sdk.features["flag-1"].Enabled,
					"buggy: SDK still shows stale Enabled=true")
			}
		})
	}
}

// TestGetFeatureFlagsNoneResponsePreservesRequestedAt verifies that "None"
// responses do not advance the SDK's requestedAt, while "All" and "Diff"
// responses do.
func TestGetFeatureFlagsNoneResponsePreservesRequestedAt(t *testing.T) {
	t.Parallel()

	baseTime := int64(1710000000)

	features := []*featureproto.Feature{
		{
			Id:        "flag-1",
			Version:   1,
			Enabled:   true,
			UpdatedAt: baseTime - 3600,
		},
	}

	sdk := &featureFlagsSDKCache{features: make(map[string]*featureproto.Feature)}

	resp := simulateGetFeatureFlagsFixed(features, sdk.featureFlagsID, sdk.requestedAt, baseTime)
	require.Equal(t, "all", resp.ResponseType)
	sdk.applyResponse(resp)
	assert.Equal(t, baseTime, sdk.requestedAt, "All response sets requestedAt to now")

	for _, offset := range []int64{60, 120, 180, 240, 300} {
		now := baseTime + offset
		resp = simulateGetFeatureFlagsFixed(features, sdk.featureFlagsID, sdk.requestedAt, now)
		require.Equal(t, "none", resp.ResponseType)
		sdk.applyResponse(resp)
		assert.Equal(t, baseTime, sdk.requestedAt,
			fmt.Sprintf("T=%d: None response must preserve requestedAt at %d", offset, baseTime))
	}

	updatedFeatures := []*featureproto.Feature{
		{
			Id:        "flag-1",
			Version:   2,
			Enabled:   false,
			UpdatedAt: baseTime + 350,
		},
	}
	diffNow := baseTime + 360
	resp = simulateGetFeatureFlagsFixed(updatedFeatures, sdk.featureFlagsID, sdk.requestedAt, diffNow)
	require.Equal(t, "diff", resp.ResponseType)
	sdk.applyResponse(resp)
	assert.Equal(t, diffNow, sdk.requestedAt, "Diff response advances requestedAt to now")
	assert.False(t, sdk.features["flag-1"].Enabled, "SDK reflects the updated flag")
}

// ---------------------------------------------------------------------------
// GetSegmentUsers tests
// ---------------------------------------------------------------------------

// TestGetSegmentUsersStaleCacheBuggyVsFixed runs the same staleness scenario
// for GetSegmentUsers, verifying the fix works for segment data too.
func TestGetSegmentUsersStaleCacheBuggyVsFixed(t *testing.T) {
	t.Parallel()

	baseTime := int64(1710000000)

	segmentV1 := &featureproto.SegmentUsers{
		SegmentId: "seg-1",
		Users: []*featureproto.SegmentUser{
			{SegmentId: "seg-1", UserId: "user-1"},
		},
		UpdatedAt: baseTime - 3600,
	}
	segmentV2 := &featureproto.SegmentUsers{
		SegmentId: "seg-1",
		Users: []*featureproto.SegmentUser{
			{SegmentId: "seg-1", UserId: "user-1"},
			{SegmentId: "seg-1", UserId: "user-2"},
		},
		UpdatedAt: baseTime + 125,
	}

	patterns := []struct {
		desc        string
		serverFn    segmentUsersServerFunc
		expectStale bool
	}{
		{
			desc:        "buggy: None responses advance requestedAt, segment update is missed",
			serverFn:    simulateGetSegmentUsersBuggy,
			expectStale: true,
		},
		{
			desc:        "fixed: None responses preserve requestedAt, segment update is detected",
			serverFn:    simulateGetSegmentUsersFixed,
			expectStale: false,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			serverSegments := []*featureproto.SegmentUsers{segmentV1}
			sdk := &segmentUsersSDKCache{segments: make(map[string]*featureproto.SegmentUsers)}

			// Phase 1: Initial sync (T=0) → "All" (requestedAt=0 < now-30days)
			resp := p.serverFn(serverSegments, sdk.segmentIDs, sdk.requestedAt, baseTime)
			require.Equal(t, "all", resp.ResponseType)
			sdk.applyResponse(resp)
			require.Len(t, sdk.segments["seg-1"].Users, 1)

			// Phase 2: "None" polls advance requestedAt (T=60, T=120)
			for _, offset := range []int64{60, 120} {
				resp = p.serverFn(serverSegments, sdk.segmentIDs, sdk.requestedAt, baseTime+offset)
				require.Equal(t, "none", resp.ResponseType)
				sdk.applyResponse(resp)
			}

			// Phase 3: Segment updated at T=125, but Redis hasn't refreshed at T=180
			resp = p.serverFn(serverSegments, sdk.segmentIDs, sdk.requestedAt, baseTime+180)
			require.Equal(t, "none", resp.ResponseType)
			sdk.applyResponse(resp)

			// Phase 4: Redis refreshes at T=181
			serverSegments = []*featureproto.SegmentUsers{segmentV2}

			// Phase 5: SDK polls at T=240
			resp = p.serverFn(serverSegments, sdk.segmentIDs, sdk.requestedAt, baseTime+240)

			if p.expectStale {
				assert.Equal(t, "none", resp.ResponseType,
					"buggy: updated segment is missed because requestedAt was advanced past UpdatedAt")
				sdk.applyResponse(resp)
				assert.Len(t, sdk.segments["seg-1"].Users, 1,
					"buggy: SDK still has old segment data (1 user)")
			} else {
				assert.Equal(t, "diff", resp.ResponseType,
					"fixed: updated segment is detected")
				sdk.applyResponse(resp)
				assert.Len(t, sdk.segments["seg-1"].Users, 2,
					"fixed: SDK has new segment data (2 users)")
			}
		})
	}
}

// TestGetSegmentUsersNoneResponsePreservesRequestedAt verifies that "None"
// responses in GetSegmentUsers preserve requestedAt.
func TestGetSegmentUsersNoneResponsePreservesRequestedAt(t *testing.T) {
	t.Parallel()

	baseTime := int64(1710000000)

	segments := []*featureproto.SegmentUsers{
		{
			SegmentId: "seg-1",
			Users:     []*featureproto.SegmentUser{{SegmentId: "seg-1", UserId: "user-1"}},
			UpdatedAt: baseTime - 3600,
		},
	}

	sdk := &segmentUsersSDKCache{segments: make(map[string]*featureproto.SegmentUsers)}

	// "All" response sets requestedAt
	resp := simulateGetSegmentUsersFixed(segments, sdk.segmentIDs, sdk.requestedAt, baseTime)
	require.Equal(t, "all", resp.ResponseType)
	sdk.applyResponse(resp)
	assert.Equal(t, baseTime, sdk.requestedAt, "All response sets requestedAt to now")

	// Multiple "None" responses should NOT advance requestedAt
	for _, offset := range []int64{60, 120, 180, 240, 300} {
		now := baseTime + offset
		resp = simulateGetSegmentUsersFixed(segments, sdk.segmentIDs, sdk.requestedAt, now)
		require.Equal(t, "none", resp.ResponseType)
		sdk.applyResponse(resp)
		assert.Equal(t, baseTime, sdk.requestedAt,
			fmt.Sprintf("T=%d: None response must preserve requestedAt at %d", offset, baseTime))
	}

	// "Diff" response DOES advance requestedAt
	updatedSegments := []*featureproto.SegmentUsers{
		{
			SegmentId: "seg-1",
			Users: []*featureproto.SegmentUser{
				{SegmentId: "seg-1", UserId: "user-1"},
				{SegmentId: "seg-1", UserId: "user-2"},
			},
			UpdatedAt: baseTime + 350,
		},
	}
	diffNow := baseTime + 360
	resp = simulateGetSegmentUsersFixed(updatedSegments, sdk.segmentIDs, sdk.requestedAt, diffNow)
	require.Equal(t, "diff", resp.ResponseType)
	sdk.applyResponse(resp)
	assert.Equal(t, diffNow, sdk.requestedAt, "Diff response advances requestedAt to now")
	assert.Len(t, sdk.segments["seg-1"].Users, 2, "SDK reflects the updated segment")
}

// ---------------------------------------------------------------------------
// Same-second boundary tests (UpdatedAt == RequestedAt)
// ---------------------------------------------------------------------------

// TestGetFeatureFlagsSameSecondUpdate verifies that a feature updated in the
// exact same Unix second as the previous response's RequestedAt is included
// in the diff (>= comparison).
func TestGetFeatureFlagsSameSecondUpdate(t *testing.T) {
	t.Parallel()

	baseTime := int64(1710000000)

	flagV1 := &featureproto.Feature{
		Id:        "flag-1",
		Version:   1,
		Enabled:   true,
		UpdatedAt: baseTime - 3600,
	}
	flagSameSecond := &featureproto.Feature{
		Id:        "flag-2",
		Version:   1,
		Enabled:   false,
		UpdatedAt: baseTime, // same second as RequestedAt from the "All" response
	}

	serverFeatures := []*featureproto.Feature{flagV1}
	sdk := &featureFlagsSDKCache{features: make(map[string]*featureproto.Feature)}

	// Initial sync at T=baseTime → requestedAt = baseTime
	resp := simulateGetFeatureFlagsFixed(serverFeatures, sdk.featureFlagsID, sdk.requestedAt, baseTime)
	require.Equal(t, "all", resp.ResponseType)
	sdk.applyResponse(resp)
	require.Equal(t, baseTime, sdk.requestedAt)

	// flag-2 is added at exactly baseTime (same second as requestedAt)
	serverFeatures = []*featureproto.Feature{flagV1, flagSameSecond}

	// SDK polls at T=baseTime+60 → ffID mismatch → Diff
	resp = simulateGetFeatureFlagsFixed(serverFeatures, sdk.featureFlagsID, sdk.requestedAt, baseTime+60)
	require.Equal(t, "diff", resp.ResponseType)
	assert.Len(t, resp.Features, 1,
		fmt.Sprintf("flag-2 (UpdatedAt == RequestedAt) is included with >= comparison"))
	assert.Equal(t, "flag-2", resp.Features[0].Id)
}

// TestGetFeatureFlagsSameSecondAfterDiff verifies the same-second edge case
// specifically after a Diff response (where RequestedAt = now and a new feature
// could be updated at exactly that second).
func TestGetFeatureFlagsSameSecondAfterDiff(t *testing.T) {
	t.Parallel()

	baseTime := int64(1710000000)

	flagV1 := &featureproto.Feature{
		Id:        "flag-1",
		Version:   1,
		Enabled:   true,
		UpdatedAt: baseTime - 3600,
	}
	flagV2 := &featureproto.Feature{
		Id:        "flag-1",
		Version:   2,
		Enabled:   false,
		UpdatedAt: baseTime + 100,
	}

	sdk := &featureFlagsSDKCache{features: make(map[string]*featureproto.Feature)}

	// Initial "All" at T=baseTime
	serverFeatures := []*featureproto.Feature{flagV1}
	resp := simulateGetFeatureFlagsFixed(serverFeatures, sdk.featureFlagsID, sdk.requestedAt, baseTime)
	sdk.applyResponse(resp)

	// Flag updated at T=100, SDK gets "Diff" at T=100 (same second)
	serverFeatures = []*featureproto.Feature{flagV2}
	diffTime := baseTime + 100
	resp = simulateGetFeatureFlagsFixed(serverFeatures, sdk.featureFlagsID, sdk.requestedAt, diffTime)
	require.Equal(t, "diff", resp.ResponseType)
	require.Len(t, resp.Features, 1, "flag-1 is included in diff")
	sdk.applyResponse(resp)
	assert.False(t, sdk.features["flag-1"].Enabled)
	assert.Equal(t, diffTime, sdk.requestedAt, "requestedAt advances to diffTime")

	// Now another flag is updated at exactly diffTime (same second as previous RequestedAt)
	flagNew := &featureproto.Feature{
		Id:        "flag-2",
		Version:   1,
		Enabled:   true,
		UpdatedAt: diffTime, // same second as sdk.requestedAt
	}
	serverFeatures = []*featureproto.Feature{flagV2, flagNew}

	// SDK polls at T=diffTime+60
	resp = simulateGetFeatureFlagsFixed(serverFeatures, sdk.featureFlagsID, sdk.requestedAt, diffTime+60)
	require.Equal(t, "diff", resp.ResponseType)

	// With >= comparison, flag-2 (UpdatedAt == RequestedAt) IS included
	// With > comparison, it would be missed
	foundFlag2 := false
	for _, f := range resp.Features {
		if f.Id == "flag-2" {
			foundFlag2 = true
		}
	}
	assert.True(t, foundFlag2,
		"flag-2 updated at exactly RequestedAt must be included (>= comparison)")
}
