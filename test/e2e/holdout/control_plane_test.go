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

package holdout

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	btproto "github.com/bucketeer-io/bucketeer/v2/proto/batch"
	domainevent "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// RFC §4.1, §4.2.1, §4.5: create/get/list/update round trip in HASH mode,
// server-assigned sampling_seed and bucket_start, updated_at bump.
func TestCreateGetListUpdateHash(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	spec := waitingHashSpec()
	spec.percentage = 12.5
	spec.startAnalysisAt = spec.startAt + 3600
	req := spec.request(t)
	created := createHoldoutGroup(t, client, req)

	assert.NotEmpty(t, created.Id)
	assert.Equal(t, holdoutEnv.id, created.EnvironmentId)
	assert.Equal(t, req.Name, created.Name)
	assert.Equal(t, featureproto.HoldoutGroup_WAITING, created.Status)
	assert.Equal(t, featureproto.HoldoutGroup_HASH, created.AssignmentType)
	assert.Equal(t, 12.5, created.Percentage)
	assert.NotEmpty(t, created.SamplingSeed, "server assigns sampling_seed")
	assert.GreaterOrEqual(t, created.BucketStart, 0.0)
	assert.Less(t, created.BucketStart, 100.0)
	assert.Equal(t, spec.startAt, created.StartAt)
	assert.Equal(t, spec.stopAt, created.StopAt)
	assert.Equal(t, spec.startAnalysisAt, created.StartAnalysisAt)
	assert.Equal(t, created.CreatedAt, created.UpdatedAt)

	got := getHoldoutGroup(t, client, created.Id)
	assert.Equal(t, created.Id, got.Id)
	assert.Equal(t, created.SamplingSeed, got.SamplingSeed)
	assert.Equal(t, created.BucketStart, got.BucketStart)
	assert.Empty(t, got.Features)

	// A second HASH group in the same environment shares the environment
	// seed and gets its own bucket range (RFC §4.2.1).
	second := createHoldoutGroup(t, client, waitingHashSpec().request(t))
	assert.Equal(t, created.SamplingSeed, second.SamplingSeed, "environment-scoped seed is copied")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	list, err := client.ListHoldoutGroups(ctx, &featureproto.ListHoldoutGroupsRequest{
		EnvironmentId: holdoutEnv.id,
		PageSize:      100,
		Statuses:      []featureproto.HoldoutGroup_Status{featureproto.HoldoutGroup_WAITING},
		SearchKeyword: req.Name,
	})
	require.NoError(t, err)
	require.Len(t, list.HoldoutGroups, 1)
	assert.Equal(t, created.Id, list.HoldoutGroups[0].Id)

	bumpSecond(t, created)
	newStop := created.StopAt + 3600
	updated := updateHoldoutGroup(t, client, &featureproto.UpdateHoldoutGroupRequest{
		Id:          created.Id,
		Name:        wrapperspb.String(req.Name + "-renamed"),
		Description: wrapperspb.String("updated"),
		StopAt:      wrapperspb.Int64(newStop),
	})
	assert.Equal(t, req.Name+"-renamed", updated.Name)
	assert.Equal(t, "updated", updated.Description)
	assert.Equal(t, newStop, updated.StopAt)
	assert.Greater(t, updated.UpdatedAt, created.UpdatedAt, "updated_at bumps on update")
	// Assignment fields are immutable through the API: the request has no
	// field for them, and they survive the update unchanged.
	assert.Equal(t, created.SamplingSeed, updated.SamplingSeed)
	assert.Equal(t, created.BucketStart, updated.BucketStart)
	assert.Equal(t, created.Percentage, updated.Percentage)
}

// RFC §4.2.2: SEGMENT mode round trip.
func TestCreateSegmentMode(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	segment := createSegment(t, client)
	start, stop := farFutureWindow()
	req := groupSpec{
		assignment: featureproto.HoldoutGroup_SEGMENT,
		segmentID:  segment.Id,
		startAt:    start,
		stopAt:     stop,
	}.request(t)
	created := createHoldoutGroup(t, client, req)
	assert.Equal(t, featureproto.HoldoutGroup_SEGMENT, created.AssignmentType)
	assert.Equal(t, segment.Id, created.SegmentId)
	assert.Empty(t, created.SamplingSeed)
	assert.Zero(t, created.Percentage)
	got := getHoldoutGroup(t, client, created.Id)
	assert.Equal(t, segment.Id, got.SegmentId)

	t.Run("PercentageRejectedInSegmentMode", func(t *testing.T) {
		start, stop := farFutureWindow()
		bad := groupSpec{
			assignment: featureproto.HoldoutGroup_SEGMENT, segmentID: segment.Id, percentage: 10,
			startAt: start, stopAt: stop,
		}.request(t)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := client.CreateHoldoutGroup(ctx, bad)
		requireCode(t, err, codes.InvalidArgument)
	})
	t.Run("UnknownSegmentRejected", func(t *testing.T) {
		start, stop := farFutureWindow()
		bad := groupSpec{
			assignment: featureproto.HoldoutGroup_SEGMENT, segmentID: "missing-" + newUUID(t),
			startAt: start, stopAt: stop,
		}.request(t)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := client.CreateHoldoutGroup(ctx, bad)
		requireCode(t, err, codes.NotFound)
	})
}

// RFC §4.5 validation on CreateHoldoutGroup.
func TestCreateValidation(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	// Each case gets its own window: the overlap check runs before the
	// argument checks, so a shared window would mask the expected code.
	base := func() *featureproto.CreateHoldoutGroupRequest {
		start, stop := farFutureWindow()
		return groupSpec{assignment: featureproto.HoldoutGroup_HASH, percentage: 10, startAt: start, stopAt: stop}.request(t)
	}
	existing := createHoldoutGroup(t, client, base())

	cases := []struct {
		desc string
		req  func() *featureproto.CreateHoldoutGroupRequest
		code codes.Code
	}{
		{"client bucket_start rejected", func() *featureproto.CreateHoldoutGroupRequest {
			r := base()
			r.BucketStart = wrapperspb.Double(10)
			return r
		}, codes.InvalidArgument},
		{"client sampling_seed rejected", func() *featureproto.CreateHoldoutGroupRequest {
			r := base()
			r.SamplingSeed = wrapperspb.String("seed")
			return r
		}, codes.InvalidArgument},
		{"percentage 0 rejected", func() *featureproto.CreateHoldoutGroupRequest {
			r := base()
			r.Percentage = 0
			return r
		}, codes.InvalidArgument},
		{"percentage above 50 rejected", func() *featureproto.CreateHoldoutGroupRequest {
			r := base()
			r.Percentage = 50.001
			return r
		}, codes.InvalidArgument},
		{"negative percentage rejected", func() *featureproto.CreateHoldoutGroupRequest {
			r := base()
			r.Percentage = -1
			return r
		}, codes.InvalidArgument},
		{"start_at >= stop_at rejected", func() *featureproto.CreateHoldoutGroupRequest {
			r := base()
			r.StartAt, r.StopAt = r.StopAt, r.StartAt
			return r
		}, codes.InvalidArgument},
		{"start_analysis_at before start_at rejected", func() *featureproto.CreateHoldoutGroupRequest {
			r := base()
			r.StartAnalysisAt = r.StartAt - 1
			return r
		}, codes.InvalidArgument},
		{"start_analysis_at after stop_at rejected", func() *featureproto.CreateHoldoutGroupRequest {
			r := base()
			r.StartAnalysisAt = r.StopAt
			return r
		}, codes.InvalidArgument},
		{"stop_at in the past rejected", func() *featureproto.CreateHoldoutGroupRequest {
			r := base()
			r.StartAt = time.Now().Add(-2 * time.Hour).Unix()
			r.StopAt = time.Now().Add(-time.Hour).Unix()
			return r
		}, codes.InvalidArgument},
		{"empty name rejected", func() *featureproto.CreateHoldoutGroupRequest {
			r := base()
			r.Name = ""
			return r
		}, codes.InvalidArgument},
		{"segment_id in hash mode rejected", func() *featureproto.CreateHoldoutGroupRequest {
			r := base()
			r.SegmentId = "segment"
			return r
		}, codes.InvalidArgument},
		{"window overlapping a WAITING group rejected", func() *featureproto.CreateHoldoutGroupRequest {
			// Starts inside the window of `existing`, which is WAITING.
			r := base()
			r.StartAt, r.StopAt = existing.StartAt+60, existing.StopAt+60
			return r
		}, codes.FailedPrecondition},
		{"window touching a WAITING group end is accepted as half-open", nil, codes.OK},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			if tc.req == nil {
				// [stop, stop+1h) does not intersect [start, stop).
				r := base()
				r.StartAt, r.StopAt = existing.StopAt, existing.StopAt+3600
				g := createHoldoutGroup(t, client, r)
				assert.Equal(t, featureproto.HoldoutGroup_WAITING, g.Status)
				return
			}
			_, err := client.CreateHoldoutGroup(ctx, tc.req())
			requireCode(t, err, tc.code)
		})
	}

	t.Run("percentage is normalized to the DECIMAL(6,3) grid", func(t *testing.T) {
		s, e := farFutureWindow()
		r := groupSpec{assignment: featureproto.HoldoutGroup_HASH, percentage: 33.3334, startAt: s, stopAt: e}.request(t)
		g := createHoldoutGroup(t, client, r)
		assert.Equal(t, 33.333, g.Percentage)
	})

	t.Run("update that moves the window onto another WAITING group is rejected", func(t *testing.T) {
		s, e := farFutureWindow()
		other := createHoldoutGroup(t, client, groupSpec{
			assignment: featureproto.HoldoutGroup_HASH, percentage: 10, startAt: s, stopAt: e,
		}.request(t))
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := client.UpdateHoldoutGroup(ctx, &featureproto.UpdateHoldoutGroupRequest{
			Id:            other.Id,
			EnvironmentId: holdoutEnv.id,
			StartAt:       wrapperspb.Int64(existing.StartAt),
			StopAt:        wrapperspb.Int64(existing.StopAt),
		})
		requireCode(t, err, codes.FailedPrecondition)
		_, err = client.UpdateHoldoutGroup(ctx, &featureproto.UpdateHoldoutGroupRequest{
			Id:              other.Id,
			EnvironmentId:   holdoutEnv.id,
			StartAnalysisAt: wrapperspb.Int64(other.StopAt + 1),
		})
		requireCode(t, err, codes.InvalidArgument, "start_analysis_at outside the window")
		_, err = client.UpdateHoldoutGroup(ctx, &featureproto.UpdateHoldoutGroupRequest{
			Id: other.Id, EnvironmentId: holdoutEnv.id,
		})
		requireCode(t, err, codes.InvalidArgument, "nothing to update")
	})
}

// RFC §4.5 validation on Add/RemoveHoldoutGroupFeatures and the flag-side
// guards of §4.6 (variation deletion and archive), all on WAITING groups.
func TestAddRemoveFeaturesValidation(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	flag := createEnabledFeature(t, client)
	group := createHoldoutGroup(t, client, waitingHashSpec().request(t))

	t.Run("unknown flag", func(t *testing.T) {
		err := tryAddHoldoutFeatures(t, client, group.Id, &featureproto.HoldoutGroupFeatureInput{
			FeatureId: "missing-" + newUUID(t), HoldoutVariationId: "v",
		})
		requireCode(t, err, codes.NotFound)
	})
	t.Run("holdout_variation_id not on the flag", func(t *testing.T) {
		err := tryAddHoldoutFeatures(t, client, group.Id, &featureproto.HoldoutGroupFeatureInput{
			FeatureId: flag.Id, HoldoutVariationId: "not-a-variation",
		})
		requireCode(t, err, codes.InvalidArgument)
	})
	t.Run("empty inputs", func(t *testing.T) {
		err := tryAddHoldoutFeatures(t, client, group.Id)
		requireCode(t, err, codes.InvalidArgument)
		err = tryAddHoldoutFeatures(t, client, group.Id, &featureproto.HoldoutGroupFeatureInput{FeatureId: flag.Id})
		requireCode(t, err, codes.InvalidArgument)
	})
	t.Run("archived flag", func(t *testing.T) {
		archived := createEnabledFeature(t, client)
		updateFeature(t, client, &featureproto.UpdateFeatureRequest{Id: archived.Id, Archived: wrapperspb.Bool(true)})
		err := tryAddHoldoutFeatures(t, client, group.Id, target(archived, "B"))
		requireCode(t, err, codes.FailedPrecondition)
	})

	bumpSecond(t, group)
	added := addHoldoutFeatures(t, client, group.Id, target(flag, "B"))
	require.Len(t, added.Features, 1)
	assert.Equal(t, flag.Id, added.Features[0].FeatureId)
	assert.Equal(t, variationID(flag, "B"), added.Features[0].HoldoutVariationId)
	assert.NotZero(t, added.Features[0].AddedAt)
	assert.Zero(t, added.Features[0].RemovedAt)
	assert.Greater(t, added.UpdatedAt, group.UpdatedAt, "updated_at bumps on add")

	t.Run("adding the same flag twice", func(t *testing.T) {
		err := tryAddHoldoutFeatures(t, client, group.Id, target(flag, "B"))
		requireCode(t, err, codes.AlreadyExists)
	})
	t.Run("flag targeted by another active group", func(t *testing.T) {
		other := createHoldoutGroup(t, client, waitingHashSpec().request(t))
		err := tryAddHoldoutFeatures(t, client, other.Id, target(flag, "A"))
		requireCode(t, err, codes.FailedPrecondition)
		assert.Contains(t, err.Error(), group.Name, "error names the owning group")
	})
	t.Run("deleting the holdout variation is rejected naming the group", func(t *testing.T) {
		err := tryUpdateFeature(t, client, &featureproto.UpdateFeatureRequest{
			Id: flag.Id,
			VariationChanges: []*featureproto.VariationChange{{
				ChangeType: featureproto.ChangeType_DELETE,
				Variation:  &featureproto.Variation{Id: variationID(flag, "B")},
			}},
		})
		requireCode(t, err, codes.FailedPrecondition)
		assert.Contains(t, err.Error(), group.Name)
		// Deleting the spare variation is still allowed on a target flag.
		updateFeature(t, client, &featureproto.UpdateFeatureRequest{
			Id: flag.Id,
			VariationChanges: []*featureproto.VariationChange{{
				ChangeType: featureproto.ChangeType_DELETE,
				Variation:  &featureproto.Variation{Id: variationID(flag, "C")},
			}},
		})
	})
	t.Run("archiving a target flag is rejected naming the group", func(t *testing.T) {
		err := tryUpdateFeature(t, client, &featureproto.UpdateFeatureRequest{
			Id: flag.Id, Archived: wrapperspb.Bool(true),
		})
		requireCode(t, err, codes.FailedPrecondition)
		assert.Contains(t, err.Error(), group.Name)
	})
	t.Run("other flag edits stay free", func(t *testing.T) {
		updateFeature(t, client, &featureproto.UpdateFeatureRequest{
			Id: flag.Id, Enabled: wrapperspb.Bool(false), Description: wrapperspb.String("still editable"),
		})
	})

	bumpSecond(t, added)
	removed := removeHoldoutFeatures(t, client, group.Id, flag.Id)
	assert.Empty(t, removed.Features)
	assert.Greater(t, removed.UpdatedAt, added.UpdatedAt, "updated_at bumps on remove")
	t.Run("removing a flag that is not a target", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := client.RemoveHoldoutGroupFeatures(ctx, &featureproto.RemoveHoldoutGroupFeaturesRequest{
			Id: group.Id, EnvironmentId: holdoutEnv.id, FeatureIds: []string{flag.Id},
		})
		requireCode(t, err, codes.NotFound)
	})
	t.Run("released flag can be archived and re-added elsewhere", func(t *testing.T) {
		other := createHoldoutGroup(t, client, waitingHashSpec().request(t))
		addHoldoutFeatures(t, client, other.Id, target(flag, "A"))
		removeHoldoutFeatures(t, client, other.Id, flag.Id)
		updateFeature(t, client, &featureproto.UpdateFeatureRequest{Id: flag.Id, Archived: wrapperspb.Bool(true)})
	})
}

// RFC §4.5: every mutation needs Organization_ADMIN or above; reads need
// Environment_VIEWER.
func TestAuthorization(t *testing.T) {
	t.Parallel()
	owner := newFeatureClient(t)
	editor := newFeatureClientAs(t, *envEditorAccessTokenPath)
	viewer := newFeatureClientAs(t, *envViewerAccessTokenPath)
	flag := createEnabledFeature(t, owner)
	group := createHoldoutGroup(t, owner, waitingHashSpec().request(t))
	env := holdoutEnv.id

	for _, tc := range []struct {
		desc   string
		client featureproto.FeatureServiceClient
	}{{"editor", editor}, {"viewer", viewer}} {
		t.Run(tc.desc+" is denied every mutation", func(t *testing.T) {
			c := tc.client
			calls := map[string]func(ctx context.Context) error{
				"CreateHoldoutGroup": func(ctx context.Context) error {
					_, err := c.CreateHoldoutGroup(ctx, waitingHashSpec().request(t))
					return err
				},
				"UpdateHoldoutGroup": func(ctx context.Context) error {
					_, err := c.UpdateHoldoutGroup(ctx, &featureproto.UpdateHoldoutGroupRequest{
						Id: group.Id, EnvironmentId: env, Name: wrapperspb.String("x"),
					})
					return err
				},
				"AddHoldoutGroupFeatures": func(ctx context.Context) error {
					_, err := c.AddHoldoutGroupFeatures(ctx, &featureproto.AddHoldoutGroupFeaturesRequest{
						Id: group.Id, EnvironmentId: env, Features: []*featureproto.HoldoutGroupFeatureInput{target(flag, "B")},
					})
					return err
				},
				"RemoveHoldoutGroupFeatures": func(ctx context.Context) error {
					_, err := c.RemoveHoldoutGroupFeatures(ctx, &featureproto.RemoveHoldoutGroupFeaturesRequest{
						Id: group.Id, EnvironmentId: env, FeatureIds: []string{flag.Id},
					})
					return err
				},
				"StartHoldoutGroup": func(ctx context.Context) error {
					_, err := c.StartHoldoutGroup(ctx, &featureproto.StartHoldoutGroupRequest{Id: group.Id, EnvironmentId: env})
					return err
				},
				"StartHoldoutGroupAnalysis": func(ctx context.Context) error {
					_, err := c.StartHoldoutGroupAnalysis(ctx, &featureproto.StartHoldoutGroupAnalysisRequest{Id: group.Id, EnvironmentId: env})
					return err
				},
				"StopHoldoutGroup": func(ctx context.Context) error {
					_, err := c.StopHoldoutGroup(ctx, &featureproto.StopHoldoutGroupRequest{Id: group.Id, EnvironmentId: env})
					return err
				},
				"ArchiveHoldoutGroup": func(ctx context.Context) error {
					_, err := c.ArchiveHoldoutGroup(ctx, &featureproto.ArchiveHoldoutGroupRequest{Id: group.Id, EnvironmentId: env})
					return err
				},
				"DeleteHoldoutGroup": func(ctx context.Context) error {
					_, err := c.DeleteHoldoutGroup(ctx, &featureproto.DeleteHoldoutGroupRequest{Id: group.Id, EnvironmentId: env})
					return err
				},
			}
			for name, call := range calls {
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				err := call(ctx)
				cancel()
				requireCode(t, err, codes.PermissionDenied, name)
			}
		})
	}

	t.Run("viewer and editor can read", func(t *testing.T) {
		for _, c := range []featureproto.FeatureServiceClient{viewer, editor} {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			got, err := c.GetHoldoutGroup(ctx, &featureproto.GetHoldoutGroupRequest{Id: group.Id, EnvironmentId: env})
			require.NoError(t, err)
			assert.Equal(t, group.Id, got.HoldoutGroup.Id)
			list, err := c.ListHoldoutGroups(ctx, &featureproto.ListHoldoutGroupsRequest{EnvironmentId: env, PageSize: 10})
			require.NoError(t, err)
			assert.NotEmpty(t, list.HoldoutGroups)
			cancel()
		}
	})

	t.Run("owner succeeds on the same mutations", func(t *testing.T) {
		// The group is still untouched after the denied calls.
		got := getHoldoutGroup(t, owner, group.Id)
		assert.Equal(t, group.UpdatedAt, got.UpdatedAt)
		assert.Empty(t, got.Features)
		updateHoldoutGroup(t, owner, &featureproto.UpdateHoldoutGroupRequest{Id: group.Id, Name: wrapperspb.String(group.Name + "-owner")})
		addHoldoutFeatures(t, owner, group.Id, target(flag, "B"))
		removeHoldoutFeatures(t, owner, group.Id, flag.Id)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := owner.DeleteHoldoutGroup(ctx, &featureproto.DeleteHoldoutGroupRequest{Id: group.Id, EnvironmentId: env})
		require.NoError(t, err)
		_, err = owner.GetHoldoutGroup(ctx, &featureproto.GetHoldoutGroupRequest{Id: group.Id, EnvironmentId: env})
		requireCode(t, err, codes.NotFound, "deleted group is gone")
	})
}

// RFC §4.5: Delete is accepted only while WAITING.
func TestDeleteWhileWaiting(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	group := createHoldoutGroup(t, client, waitingHashSpec().request(t))
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.ArchiveHoldoutGroup(ctx, &featureproto.ArchiveHoldoutGroupRequest{Id: group.Id, EnvironmentId: holdoutEnv.id})
	requireCode(t, err, codes.FailedPrecondition, "a WAITING group is deleted, not archived")
	_, err = client.StopHoldoutGroup(ctx, &featureproto.StopHoldoutGroupRequest{Id: group.Id, EnvironmentId: holdoutEnv.id})
	requireCode(t, err, codes.FailedPrecondition, "a WAITING group cannot be stopped")
	_, err = client.StartHoldoutGroupAnalysis(ctx, &featureproto.StartHoldoutGroupAnalysisRequest{Id: group.Id, EnvironmentId: holdoutEnv.id})
	requireCode(t, err, codes.FailedPrecondition, "a WAITING group cannot enter analysis")
	_, err = client.DeleteHoldoutGroup(ctx, &featureproto.DeleteHoldoutGroupRequest{Id: group.Id, EnvironmentId: holdoutEnv.id})
	require.NoError(t, err)
	_, err = client.GetHoldoutGroup(ctx, &featureproto.GetHoldoutGroupRequest{Id: group.Id, EnvironmentId: holdoutEnv.id})
	requireCode(t, err, codes.NotFound)
}

// RFC §4.6 manual lifecycle: Start clamps start_at, StartAnalysis freezes
// the flag list, Stop records FORCE_STOPPED and stopped_at, a group that has
// served is archived rather than deleted. RFC §4.5: every mutation writes an
// audit log entry and Add/Remove also write a feature-scoped entry.
func TestLifecycleManual(t *testing.T) {
	servingMu.Lock()
	defer servingMu.Unlock()
	client := newFeatureClient(t)
	flagA := createEnabledFeature(t, client)
	flagB := createEnabledFeature(t, client)
	since := time.Now().Unix() - 60

	// Future start_at: manual Start must clamp it to now.
	spec := servingHashSpec()
	spec.startAt = time.Now().Add(30 * time.Minute).Unix()
	group := createHoldoutGroup(t, client, spec.request(t))
	bumpSecond(t, group)
	group = updateHoldoutGroup(t, client, &featureproto.UpdateHoldoutGroupRequest{
		Id: group.Id, Description: wrapperspb.String("lifecycle"),
	})
	bumpSecond(t, group)
	withA := addHoldoutFeatures(t, client, group.Id, target(flagA, "B"))
	assert.Greater(t, withA.UpdatedAt, group.UpdatedAt)

	bumpSecond(t, withA)
	before := time.Now().Unix()
	running := startHoldoutGroup(t, client, group.Id)
	assert.Equal(t, featureproto.HoldoutGroup_RUNNING, running.Status)
	assert.LessOrEqual(t, running.StartAt, time.Now().Unix(), "start_at clamped to now")
	assert.GreaterOrEqual(t, running.StartAt, before)
	assert.Greater(t, running.UpdatedAt, withA.UpdatedAt, "updated_at bumps on start")
	assert.Len(t, running.Features, 1)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.DeleteHoldoutGroup(ctx, &featureproto.DeleteHoldoutGroupRequest{Id: group.Id, EnvironmentId: holdoutEnv.id})
	requireCode(t, err, codes.FailedPrecondition, "a RUNNING group cannot be deleted")
	_, err = client.StartHoldoutGroup(ctx, &featureproto.StartHoldoutGroupRequest{Id: group.Id, EnvironmentId: holdoutEnv.id})
	requireCode(t, err, codes.FailedPrecondition, "a RUNNING group cannot be started again")
	_, err = client.UpdateHoldoutGroup(ctx, &featureproto.UpdateHoldoutGroupRequest{
		Id: group.Id, EnvironmentId: holdoutEnv.id, StartAt: wrapperspb.Int64(running.StartAt - 10),
	})
	requireCode(t, err, codes.FailedPrecondition, "start_at is immutable while serving")

	t.Run("a second group cannot serve at the same time", func(t *testing.T) {
		second := createHoldoutGroup(t, client, groupSpec{
			assignment: featureproto.HoldoutGroup_HASH, percentage: 5,
			startAt: running.StopAt, stopAt: running.StopAt + 3600,
		}.request(t))
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := client.StartHoldoutGroup(ctx, &featureproto.StartHoldoutGroupRequest{Id: second.Id, EnvironmentId: holdoutEnv.id})
		requireCode(t, err, codes.FailedPrecondition)
		assert.Contains(t, err.Error(), group.Name)
	})

	// Adding while RUNNING is allowed.
	bumpSecond(t, running)
	withB := addHoldoutFeatures(t, client, group.Id, target(flagB, "B"))
	assert.Len(t, withB.Features, 2)

	bumpSecond(t, withB)
	analysis := startHoldoutGroupAnalysis(t, client, group.Id)
	assert.Equal(t, featureproto.HoldoutGroup_ANALYSIS, analysis.Status)
	assert.NotZero(t, analysis.AnalysisStartedAt)
	assert.Greater(t, analysis.UpdatedAt, withB.UpdatedAt, "updated_at bumps on analysis start")

	t.Run("flag list is frozen during ANALYSIS", func(t *testing.T) {
		flagC := createEnabledFeature(t, client)
		err := tryAddHoldoutFeatures(t, client, group.Id, target(flagC, "B"))
		requireCode(t, err, codes.FailedPrecondition)
		// Remove stays available as the escape hatch (RFC §4.3.1).
		bumpSecond(t, analysis)
		removed := removeHoldoutFeatures(t, client, group.Id, flagB.Id)
		assert.Len(t, removed.Features, 1)
		assert.Greater(t, removed.UpdatedAt, analysis.UpdatedAt)
	})

	current := getHoldoutGroup(t, client, group.Id)
	bumpSecond(t, current)
	before = time.Now().Unix()
	stopped := stopHoldoutGroup(t, client, group.Id)
	assert.Equal(t, featureproto.HoldoutGroup_FORCE_STOPPED, stopped.Status)
	assert.GreaterOrEqual(t, stopped.StoppedAt, before)
	assert.Greater(t, stopped.UpdatedAt, current.UpdatedAt, "updated_at bumps on stop")

	_, err = client.DeleteHoldoutGroup(ctx, &featureproto.DeleteHoldoutGroupRequest{Id: group.Id, EnvironmentId: holdoutEnv.id})
	requireCode(t, err, codes.FailedPrecondition, "a group that has served cannot be deleted")
	_, err = client.StopHoldoutGroup(ctx, &featureproto.StopHoldoutGroupRequest{Id: group.Id, EnvironmentId: holdoutEnv.id})
	requireCode(t, err, codes.FailedPrecondition, "a stopped group cannot be stopped again")
	err = tryAddHoldoutFeatures(t, client, group.Id, target(flagB, "B"))
	requireCode(t, err, codes.FailedPrecondition, "a stopped group accepts no flags")
	// The stopped group's row no longer counts as a current target, so the
	// flag can be enrolled elsewhere (RFC §4.1 removed_at semantics).
	other := createHoldoutGroup(t, client, waitingHashSpec().request(t))
	addHoldoutFeatures(t, client, other.Id, target(flagA, "A"))

	bumpSecond(t, stopped)
	archived := archiveHoldoutGroup(t, client, group.Id)
	assert.True(t, archived.Archived)
	assert.Greater(t, archived.UpdatedAt, stopped.UpdatedAt, "updated_at bumps on archive")
	_, err = client.UpdateHoldoutGroup(ctx, &featureproto.UpdateHoldoutGroupRequest{
		Id: group.Id, EnvironmentId: holdoutEnv.id, Name: wrapperspb.String("x"),
	})
	requireCode(t, err, codes.FailedPrecondition, "an archived group is read-only")

	t.Run("audit log has an entry per mutation", func(t *testing.T) {
		want := []domainevent.Event_Type{
			domainevent.Event_HOLDOUT_GROUP_CREATED,
			domainevent.Event_HOLDOUT_GROUP_UPDATED,
			domainevent.Event_HOLDOUT_GROUP_FEATURE_ADDED,
			domainevent.Event_HOLDOUT_GROUP_STARTED,
			domainevent.Event_HOLDOUT_GROUP_ANALYSIS_STARTED,
			domainevent.Event_HOLDOUT_GROUP_FEATURE_REMOVED,
			domainevent.Event_HOLDOUT_GROUP_STOPPED,
			domainevent.Event_HOLDOUT_GROUP_ARCHIVED,
		}
		requireAuditLogTypes(t, domainevent.Event_HOLDOUT_GROUP, group.Id, since, want...)
		// Add/Remove also land on the flag's own history (RFC §4.5).
		requireAuditLogTypes(t, domainevent.Event_FEATURE, flagB.Id, since,
			domainevent.Event_HOLDOUT_GROUP_FEATURE_ADDED,
			domainevent.Event_HOLDOUT_GROUP_FEATURE_REMOVED,
		)
	})
}

// requireAuditLogTypes waits until the audit log for entityID contains every
// wanted event type, in the order they were performed.
func requireAuditLogTypes(
	t *testing.T,
	entityType domainevent.Event_EntityType,
	entityID string,
	since int64,
	want ...domainevent.Event_Type,
) {
	t.Helper()
	client := newAuditLogClient(t)
	var got []domainevent.Event_Type
	require.Eventually(t, func() bool {
		got = got[:0]
		cursor := ""
		for {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			resp, err := client.ListAuditLogs(ctx, &auditlog.ListAuditLogsRequest{
				EnvironmentId:  holdoutEnv.id,
				EntityType:     wrapperspb.Int32(int32(entityType)),
				PageSize:       200,
				Cursor:         cursor,
				OrderBy:        auditlog.ListAuditLogsRequest_TIMESTAMP,
				OrderDirection: auditlog.ListAuditLogsRequest_ASC,
				From:           since,
				To:             time.Now().Unix() + 60,
			})
			cancel()
			if err != nil {
				t.Logf("ListAuditLogs: %v", err)
				return false
			}
			for _, l := range resp.AuditLogs {
				if l.EntityId == entityID {
					got = append(got, l.Type)
				}
			}
			if resp.Cursor == "" || len(resp.AuditLogs) == 0 || resp.Cursor == cursor {
				break
			}
			cursor = resp.Cursor
		}
		return isSubsequence(want, got)
	}, auditLogTimeout, 2*time.Second, "audit log for %s: want %v in order, got %v", entityID, want, got)
}

// isSubsequence reports whether want appears in got in order (not
// necessarily contiguous: a test may issue extra updates in between).
func isSubsequence(want, got []domainevent.Event_Type) bool {
	i := 0
	for _, g := range got {
		if i < len(want) && g == want[i] {
			i++
		}
	}
	return i == len(want)
}

// RFC §4.6 scheduled path: the HoldoutGroupStatusUpdater advances a group by
// at most one step per run (WAITING→RUNNING→ANALYSIS→STOPPED), and a WAITING
// group whose whole window was missed goes straight to STOPPED and releases
// its flags.
//
// The timestamps are placed seconds apart and the job is run through
// BatchService.ExecuteBatchJob as soon as each one has passed, so every step
// is observed on its own. The dev cron runs the same job every minute; it
// can only take the same step, never a further one, so it cannot change the
// outcome of an assertion.
func TestLifecycleScheduled(t *testing.T) {
	servingMu.Lock()
	defer servingMu.Unlock()
	client := newFeatureClient(t)
	flag := createEnabledFeature(t, client)

	t.Run("one step per run through RUNNING, ANALYSIS and STOPPED", func(t *testing.T) {
		now := time.Now().Unix()
		spec := groupSpec{
			assignment: featureproto.HoldoutGroup_HASH, percentage: 50,
			startAt: now + 10, startAnalysisAt: now + 25, stopAt: now + 40,
		}
		group := createHoldoutGroup(t, client, spec.request(t))
		addHoldoutFeatures(t, client, group.Id, target(flag, "B"))

		// Before start_at nothing happens.
		runBatchJob(t, btproto.BatchJob_HoldoutGroupStatusUpdater)
		assert.Equal(t, featureproto.HoldoutGroup_WAITING, getHoldoutGroup(t, client, group.Id).Status)

		waitUntilUnix(t, spec.startAt)
		runBatchJob(t, btproto.BatchJob_HoldoutGroupStatusUpdater)
		running := getHoldoutGroup(t, client, group.Id)
		assert.Equal(t, featureproto.HoldoutGroup_RUNNING, running.Status, "first due run: start only")
		assert.Equal(t, spec.startAt, running.StartAt, "scheduled start does not clamp start_at")
		assert.Zero(t, running.AnalysisStartedAt)

		waitUntilUnix(t, spec.startAnalysisAt)
		runBatchJob(t, btproto.BatchJob_HoldoutGroupStatusUpdater)
		analysis := getHoldoutGroup(t, client, group.Id)
		assert.Equal(t, featureproto.HoldoutGroup_ANALYSIS, analysis.Status, "second due run: analysis only")
		assert.GreaterOrEqual(t, analysis.AnalysisStartedAt, spec.startAnalysisAt)
		assert.Zero(t, analysis.StoppedAt)

		waitUntilUnix(t, spec.stopAt)
		runBatchJob(t, btproto.BatchJob_HoldoutGroupStatusUpdater)
		stopped := getHoldoutGroup(t, client, group.Id)
		assert.Equal(t, featureproto.HoldoutGroup_STOPPED, stopped.Status, "third due run: scheduled stop records STOPPED")
		assert.GreaterOrEqual(t, stopped.StoppedAt, spec.stopAt)
		assert.GreaterOrEqual(t, stopped.UpdatedAt, analysis.UpdatedAt)

		requireAuditLogTypes(t, domainevent.Event_HOLDOUT_GROUP, group.Id, now-60,
			domainevent.Event_HOLDOUT_GROUP_STARTED,
			domainevent.Event_HOLDOUT_GROUP_ANALYSIS_STARTED,
			domainevent.Event_HOLDOUT_GROUP_STOPPED,
		)
	})

	t.Run("SCHEDULED mode is stricter than MANUAL", func(t *testing.T) {
		spec := waitingHashSpec()
		group := createHoldoutGroup(t, client, spec.request(t))
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := client.StartHoldoutGroup(ctx, &featureproto.StartHoldoutGroupRequest{
			Id: group.Id, EnvironmentId: holdoutEnv.id, Mode: featureproto.StartHoldoutGroupRequest_SCHEDULED,
		})
		requireCode(t, err, codes.FailedPrecondition, "start_at has not passed")
		_, err = client.StopHoldoutGroup(ctx, &featureproto.StopHoldoutGroupRequest{
			Id: group.Id, EnvironmentId: holdoutEnv.id, Mode: featureproto.StopHoldoutGroupRequest_SCHEDULED,
		})
		requireCode(t, err, codes.FailedPrecondition, "stop_at has not passed")
	})

	t.Run("a missed window is stopped by the cron and releases its flags", func(t *testing.T) {
		now := time.Now().Unix()
		spec := groupSpec{assignment: featureproto.HoldoutGroup_HASH, percentage: 50, startAt: now + 1, stopAt: now + 3}
		group := createHoldoutGroup(t, client, spec.request(t))
		addHoldoutFeatures(t, client, group.Id, target(flag, "B"))
		waitUntilUnix(t, spec.stopAt)
		// No explicit trigger here: this leg proves the dev cron
		// (every minute) runs the job. statusUpdaterTimeout covers one run.
		var got *featureproto.HoldoutGroup
		require.Eventually(t, func() bool {
			got = getHoldoutGroup(t, client, group.Id)
			return got.Status == featureproto.HoldoutGroup_STOPPED
		}, statusUpdaterTimeout, pollInterval, "missed WAITING group was not stopped by the cron; status %v", got.GetStatus())
		assert.Zero(t, got.AnalysisStartedAt, "never served")
		assert.NotZero(t, got.StoppedAt)
		// The flag is released: another group can enroll it.
		other := createHoldoutGroup(t, client, waitingHashSpec().request(t))
		addHoldoutFeatures(t, client, other.Id, target(flag, "A"))
	})
}
