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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/wrapperspb"

	envclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	featuredomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	featurestoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	holdoutdomain "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/domain"
	holdoutmock "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	databasemock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestCreateHoldoutGroup(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()

	t.Run("org member is denied", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service := createFeatureServiceWithGetAccountByEnvironmentMock(
			controller,
			accountproto.AccountV2_Role_Organization_MEMBER,
			accountproto.AccountV2_Role_Environment_EDITOR,
		)
		_, err := service.CreateHoldoutGroup(createContextWithTokenRoleUnassigned(), &featureproto.CreateHoldoutGroupRequest{
			EnvironmentId: "ns0",
			Name:          "holdout",
		})
		assert.Equal(t, statusPermissionDenied.Err(), err)
	})

	t.Run("client assignment is rejected", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, _ := newHoldoutService(t, controller)
		_, err := service.CreateHoldoutGroup(createContextWithTokenRoleUnassigned(), hashCreateRequest(now, &wrapperspb.DoubleValue{Value: 1}))
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("hash group copies the environment seed", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectGoal(service, "g1")
		expectTx(service)
		expectHoldoutEvents(t, service, controller, 1, eventproto.Event_HOLDOUT_GROUP_CREATED)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().ListActiveGroups(gomock.Any(), "ns0").Return(nil, nil)
		storage.EXPECT().GetEnvironmentSamplingSeed(gomock.Any(), "ns0").Return("env-seed", nil)
		storage.EXPECT().CreateHoldoutGroup(gomock.Any(), gomock.Any()).Return(nil)

		req := hashCreateRequest(now, nil)
		req.GoalIds = []string{"g1"}
		resp, err := service.CreateHoldoutGroup(createContextWithTokenRoleUnassigned(), req)
		require.NoError(t, err)
		assert.Equal(t, featureproto.HoldoutGroup_WAITING, resp.HoldoutGroup.Status)
		assert.Equal(t, "env-seed", resp.HoldoutGroup.SamplingSeed)
		assert.GreaterOrEqual(t, resp.HoldoutGroup.BucketStart, 0.0)
		assert.Less(t, resp.HoldoutGroup.BucketStart, 100.0)
		assert.Equal(t, []string{"g1"}, resp.HoldoutGroup.GoalIds)
		assert.Empty(t, resp.HoldoutGroup.ExperimentId)
		assert.Empty(t, resp.HoldoutGroup.AnalysisExperimentId)
	})

	t.Run("segment group requires an existing segment", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectTx(service)
		expectHoldoutEvents(t, service, controller, 1, eventproto.Event_HOLDOUT_GROUP_CREATED)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().ListActiveGroups(gomock.Any(), "ns0").Return(nil, nil)
		service.segmentStorage.(*featurestoragemock.MockSegmentStorage).EXPECT().GetSegment(
			gomock.Any(), "seg", "ns0",
		).Return(&featuredomain.Segment{Segment: &featureproto.Segment{Id: "seg"}}, nil, nil)
		storage.EXPECT().CreateHoldoutGroup(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := service.CreateHoldoutGroup(createContextWithTokenRoleUnassigned(), &featureproto.CreateHoldoutGroupRequest{
			EnvironmentId:  "ns0",
			Name:           "segment holdout",
			AssignmentType: featureproto.HoldoutGroup_SEGMENT,
			SegmentId:      "seg",
			StartAt:        now + 100,
			StopAt:         now + 1_000,
		})
		require.NoError(t, err)
		assert.Equal(t, "seg", resp.HoldoutGroup.SegmentId)
		assert.Empty(t, resp.HoldoutGroup.SamplingSeed)
	})

	t.Run("segment group stores the canonical segment id", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectTx(service)
		expectHoldoutEvents(t, service, controller, 1, eventproto.Event_HOLDOUT_GROUP_CREATED)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().ListActiveGroups(gomock.Any(), "ns0").Return(nil, nil)
		service.segmentStorage.(*featurestoragemock.MockSegmentStorage).EXPECT().GetSegment(
			gomock.Any(), "SEG", "ns0",
		).Return(&featuredomain.Segment{Segment: &featureproto.Segment{Id: "seg"}}, nil, nil)
		storage.EXPECT().CreateHoldoutGroup(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := service.CreateHoldoutGroup(createContextWithTokenRoleUnassigned(), &featureproto.CreateHoldoutGroupRequest{
			EnvironmentId:  "ns0",
			Name:           "segment holdout",
			AssignmentType: featureproto.HoldoutGroup_SEGMENT,
			SegmentId:      "SEG",
			StartAt:        now + 100,
			StopAt:         now + 1_000,
		})
		require.NoError(t, err)
		assert.Equal(t, "seg", resp.HoldoutGroup.SegmentId)
	})

	t.Run("overlapping period is rejected", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectTx(service)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().ListActiveGroups(gomock.Any(), "ns0").Return([]holdoutdomain.GroupSnapshot{{
			ID: "other", Name: "other", Status: featureproto.HoldoutGroup_WAITING,
			StartAt: now + 50, StopAt: now + 500,
		}}, nil)

		_, err := service.CreateHoldoutGroup(createContextWithTokenRoleUnassigned(), hashCreateRequest(now, nil))
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})
}

func TestGetAndListHoldoutGroups(t *testing.T) {
	t.Parallel()

	t.Run("unauthenticated", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, _ := newHoldoutService(t, controller)
		_, err := service.GetHoldoutGroup(context.Background(), &featureproto.GetHoldoutGroupRequest{
			Id: "hg", EnvironmentId: "ns0",
		})
		assert.Equal(t, statusUnauthenticated.Err(), err)
	})

	t.Run("viewer can get and list", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service := createFeatureServiceWithGetAccountByEnvironmentMock(
			controller,
			accountproto.AccountV2_Role_Organization_MEMBER,
			accountproto.AccountV2_Role_Environment_VIEWER,
		)
		storage := holdoutmock.NewMockHoldoutStorage(controller)
		service.holdoutStorage = storage
		group := waitingHoldout(time.Now().Unix())
		storage.EXPECT().GetHoldoutGroup(gomock.Any(), "hg", "ns0").Return(group, nil)
		storage.EXPECT().ListCurrentFeatures(gomock.Any(), "ns0", []string{"hg"}).Return([]*holdoutdomain.HoldoutGroupFeature{{
			HoldoutGroupFeature: &featureproto.HoldoutGroupFeature{
				HoldoutGroupId: "hg", FeatureId: "flag", HoldoutVariationId: "var",
			},
		}}, nil)
		storage.EXPECT().ListHoldoutGroups(gomock.Any(), gomock.Any()).Return([]*holdoutdomain.HoldoutGroup{group}, 1, int64(1), nil)
		storage.EXPECT().ListCurrentFeatures(gomock.Any(), "ns0", []string{"hg"}).Return(nil, nil)

		ctx := createContextWithTokenRoleUnassigned()
		got, err := service.GetHoldoutGroup(ctx, &featureproto.GetHoldoutGroupRequest{Id: "hg", EnvironmentId: "ns0"})
		require.NoError(t, err)
		require.Len(t, got.HoldoutGroup.Features, 1)
		assert.Equal(t, "flag", got.HoldoutGroup.Features[0].FeatureId)

		listed, err := service.ListHoldoutGroups(ctx, &featureproto.ListHoldoutGroupsRequest{
			EnvironmentId: "ns0",
			Archived:      wrapperspb.Bool(false),
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), listed.TotalCount)
		assert.Equal(t, "1", listed.Cursor)
	})

	t.Run("page size above the maximum is rejected", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, _ := newHoldoutService(t, controller)
		_, err := service.ListHoldoutGroups(createContextWithTokenRoleUnassigned(), &featureproto.ListHoldoutGroupsRequest{
			EnvironmentId: "ns0",
			PageSize:      int64(maxPageSizePerRequest + 1),
		})
		assert.Equal(t, statusExceededMaxPageSizePerRequest.Err(), err)
	})
}

func TestUpdateHoldoutGroup(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	service, storage := newHoldoutService(t, controller)
	expectTx(service)
	expectHoldoutEvents(t, service, controller, 1, eventproto.Event_HOLDOUT_GROUP_UPDATED)
	storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
	storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").Return(waitingHoldout(time.Now().Unix()), nil)
	storage.EXPECT().UpdateHoldoutGroup(gomock.Any(), gomock.Any()).Return(nil)
	storage.EXPECT().ListCurrentFeatures(gomock.Any(), "ns0", []string{"hg"}).Return(nil, nil)

	resp, err := service.UpdateHoldoutGroup(createContextWithTokenRoleUnassigned(), &featureproto.UpdateHoldoutGroupRequest{
		Id:            "hg",
		EnvironmentId: "ns0",
		Name:          wrapperspb.String("renamed"),
	})
	require.NoError(t, err)
	assert.Equal(t, "renamed", resp.HoldoutGroup.Name)
}

func TestHoldoutGroupLifecycle(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()

	t.Run("manual start clamps a future start", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectTransition(t, service, storage, controller, waitingHoldout(now), nil, eventproto.Event_HOLDOUT_GROUP_STARTED)
		resp, err := service.StartHoldoutGroup(createContextWithTokenRoleUnassigned(), &featureproto.StartHoldoutGroupRequest{
			Id: "hg", EnvironmentId: "ns0",
		})
		require.NoError(t, err)
		assert.Equal(t, featureproto.HoldoutGroup_RUNNING, resp.HoldoutGroup.Status)
		assert.LessOrEqual(t, resp.HoldoutGroup.StartAt, time.Now().Unix())
	})

	t.Run("start rejects a missing segment", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		group := waitingHoldout(now)
		group.AssignmentType = featureproto.HoldoutGroup_SEGMENT
		group.SegmentId = "seg"
		expectTx(service)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").Return(group, nil)
		storage.EXPECT().ListActiveGroups(gomock.Any(), "ns0").Return(nil, nil)
		service.segmentStorage.(*featurestoragemock.MockSegmentStorage).EXPECT().GetSegmentForUpdate(
			gomock.Any(), "seg", "ns0",
		).Return(&featuredomain.Segment{Segment: &featureproto.Segment{Id: "seg", Deleted: true}}, nil, nil)
		_, err := service.StartHoldoutGroup(createContextWithTokenRoleUnassigned(), &featureproto.StartHoldoutGroupRequest{
			Id: "hg", EnvironmentId: "ns0",
		})
		assert.Equal(t, codes.NotFound, status.Code(err))
	})

	t.Run("start rejects a segment that is still uploading", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		group := waitingHoldout(now)
		group.AssignmentType = featureproto.HoldoutGroup_SEGMENT
		group.SegmentId = "seg"
		expectTx(service)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").Return(group, nil)
		storage.EXPECT().ListActiveGroups(gomock.Any(), "ns0").Return(nil, nil)
		service.segmentStorage.(*featurestoragemock.MockSegmentStorage).EXPECT().GetSegmentForUpdate(
			gomock.Any(), "seg", "ns0",
		).Return(&featuredomain.Segment{Segment: &featureproto.Segment{
			Id: "seg", Status: featureproto.Segment_UPLOADING,
		}}, nil, nil)
		_, err := service.StartHoldoutGroup(createContextWithTokenRoleUnassigned(), &featureproto.StartHoldoutGroupRequest{
			Id: "hg", EnvironmentId: "ns0",
		})
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})

	t.Run("start rejected when another group is serving", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectTx(service)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").Return(waitingHoldout(now), nil)
		storage.EXPECT().ListActiveGroups(gomock.Any(), "ns0").Return([]holdoutdomain.GroupSnapshot{{
			ID: "live", Status: featureproto.HoldoutGroup_RUNNING, StartAt: now - 10, StopAt: now + 10,
		}}, nil)
		_, err := service.StartHoldoutGroup(createContextWithTokenRoleUnassigned(), &featureproto.StartHoldoutGroupRequest{
			Id: "hg", EnvironmentId: "ns0",
		})
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})

	t.Run("manual analysis", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		group := waitingHoldout(now)
		group.Status = featureproto.HoldoutGroup_RUNNING
		group.StartAt = now - 10
		expectTransition(t, service, storage, controller, group, nil, eventproto.Event_HOLDOUT_GROUP_ANALYSIS_STARTED)
		resp, err := service.StartHoldoutGroupAnalysis(createContextWithTokenRoleUnassigned(), &featureproto.StartHoldoutGroupAnalysisRequest{
			Id: "hg", EnvironmentId: "ns0",
		})
		require.NoError(t, err)
		assert.Equal(t, featureproto.HoldoutGroup_ANALYSIS, resp.HoldoutGroup.Status)
		assert.NotZero(t, resp.HoldoutGroup.AnalysisStartedAt)
	})

	t.Run("manual stop", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		group := waitingHoldout(now)
		group.Status = featureproto.HoldoutGroup_RUNNING
		group.StartAt = now - 10
		expectTransition(t, service, storage, controller, group, nil, eventproto.Event_HOLDOUT_GROUP_STOPPED)
		resp, err := service.StopHoldoutGroup(createContextWithTokenRoleUnassigned(), &featureproto.StopHoldoutGroupRequest{
			Id: "hg", EnvironmentId: "ns0",
		})
		require.NoError(t, err)
		assert.Equal(t, featureproto.HoldoutGroup_FORCE_STOPPED, resp.HoldoutGroup.Status)
		assert.NotZero(t, resp.HoldoutGroup.StoppedAt)
	})

	t.Run("scheduled stop waits for stop_at", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		group := waitingHoldout(now)
		group.Status = featureproto.HoldoutGroup_RUNNING
		group.StartAt = now - 10
		expectTx(service)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").Return(group, nil)
		storage.EXPECT().ListActiveGroups(gomock.Any(), "ns0").Return(nil, nil)
		_, err := service.StopHoldoutGroup(createContextWithTokenRoleUnassigned(), &featureproto.StopHoldoutGroupRequest{
			Id: "hg", EnvironmentId: "ns0", Mode: featureproto.StopHoldoutGroupRequest_SCHEDULED,
		})
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})

	t.Run("scheduled stop finishes a due group", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		group := waitingHoldout(now)
		group.Status = featureproto.HoldoutGroup_ANALYSIS
		group.StartAt = now - 100
		group.StopAt = now - 1
		expectTransition(t, service, storage, controller, group, nil, eventproto.Event_HOLDOUT_GROUP_STOPPED)
		resp, err := service.StopHoldoutGroup(createContextWithTokenRoleUnassigned(), &featureproto.StopHoldoutGroupRequest{
			Id: "hg", EnvironmentId: "ns0", Mode: featureproto.StopHoldoutGroupRequest_SCHEDULED,
		})
		require.NoError(t, err)
		assert.Equal(t, featureproto.HoldoutGroup_STOPPED, resp.HoldoutGroup.Status)
	})

	t.Run("archive a stopped group", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		group := waitingHoldout(now)
		group.Status = featureproto.HoldoutGroup_STOPPED
		expectTransition(t, service, storage, controller, group, nil, eventproto.Event_HOLDOUT_GROUP_ARCHIVED)
		resp, err := service.ArchiveHoldoutGroup(createContextWithTokenRoleUnassigned(), &featureproto.ArchiveHoldoutGroupRequest{
			Id: "hg", EnvironmentId: "ns0",
		})
		require.NoError(t, err)
		assert.True(t, resp.HoldoutGroup.Archived)
	})

	t.Run("delete a waiting group", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectTransition(t, service, storage, controller, waitingHoldout(now), nil, eventproto.Event_HOLDOUT_GROUP_DELETED)
		_, err := service.DeleteHoldoutGroup(createContextWithTokenRoleUnassigned(), &featureproto.DeleteHoldoutGroupRequest{
			Id: "hg", EnvironmentId: "ns0",
		})
		require.NoError(t, err)
	})
}

func TestAddAndRemoveHoldoutGroupFeatures(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()

	t.Run("add a flag", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		group := waitingHoldout(now)
		group.Status = featureproto.HoldoutGroup_RUNNING
		group.StartAt = now - 10
		expectTx(service)
		expectHoldoutEvents(t, service, controller, 2, eventproto.Event_HOLDOUT_GROUP_FEATURE_ADDED)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").Return(group, nil)
		expectFeature(service, "flag", "var", false)
		storage.EXPECT().ListLiveFeatureTargets(gomock.Any(), "ns0", "flag").Return(nil, nil)
		storage.EXPECT().CreateHoldoutGroupFeatures(gomock.Any(), gomock.Any()).Return(nil)
		storage.EXPECT().UpdateHoldoutGroup(gomock.Any(), gomock.Any()).Return(nil)
		// The current targets are loaded before the change, inside the transaction.
		storage.EXPECT().ListCurrentFeatures(gomock.Any(), "ns0", []string{"hg"}).Return(nil, nil)

		resp, err := service.AddHoldoutGroupFeatures(createContextWithTokenRoleUnassigned(), &featureproto.AddHoldoutGroupFeaturesRequest{
			Id: "hg", EnvironmentId: "ns0",
			Features: []*featureproto.HoldoutGroupFeatureInput{{FeatureId: "flag", HoldoutVariationId: "var"}},
		})
		require.NoError(t, err)
		require.Len(t, resp.HoldoutGroup.Features, 1)
		assert.Equal(t, "flag", resp.HoldoutGroup.Features[0].FeatureId)
		assert.Equal(t, "var", resp.HoldoutGroup.Features[0].HoldoutVariationId)
	})

	t.Run("noncanonical feature id is rejected", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectTx(service)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").Return(waitingHoldout(now), nil)
		service.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().GetFeatureForUpdate(
			gomock.Any(), "FLAG", "ns0",
		).Return(&featuredomain.Feature{Feature: &featureproto.Feature{
			Id:         "flag",
			Variations: []*featureproto.Variation{{Id: "var"}},
		}}, nil)
		_, err := service.AddHoldoutGroupFeatures(createContextWithTokenRoleUnassigned(), &featureproto.AddHoldoutGroupFeaturesRequest{
			Id: "hg", EnvironmentId: "ns0",
			Features: []*featureproto.HoldoutGroupFeatureInput{{FeatureId: "FLAG", HoldoutVariationId: "var"}},
		})
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("add is frozen during analysis", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		group := waitingHoldout(now)
		group.Status = featureproto.HoldoutGroup_ANALYSIS
		expectTx(service)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").Return(group, nil)
		_, err := service.AddHoldoutGroupFeatures(createContextWithTokenRoleUnassigned(), &featureproto.AddHoldoutGroupFeaturesRequest{
			Id: "hg", EnvironmentId: "ns0",
			Features: []*featureproto.HoldoutGroupFeatureInput{{FeatureId: "flag", HoldoutVariationId: "var"}},
		})
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})

	t.Run("archived flag is rejected", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectTx(service)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").Return(waitingHoldout(now), nil)
		expectFeature(service, "flag", "var", true)
		_, err := service.AddHoldoutGroupFeatures(createContextWithTokenRoleUnassigned(), &featureproto.AddHoldoutGroupFeaturesRequest{
			Id: "hg", EnvironmentId: "ns0",
			Features: []*featureproto.HoldoutGroupFeatureInput{{FeatureId: "flag", HoldoutVariationId: "var"}},
		})
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})

	t.Run("variation must belong to the flag", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectTx(service)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").Return(waitingHoldout(now), nil)
		expectFeature(service, "flag", "var", false)
		_, err := service.AddHoldoutGroupFeatures(createContextWithTokenRoleUnassigned(), &featureproto.AddHoldoutGroupFeaturesRequest{
			Id: "hg", EnvironmentId: "ns0",
			Features: []*featureproto.HoldoutGroupFeatureInput{{FeatureId: "flag", HoldoutVariationId: "missing"}},
		})
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("flag already in this group", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectTx(service)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").Return(waitingHoldout(now), nil)
		expectFeature(service, "flag", "var", false)
		storage.EXPECT().ListLiveFeatureTargets(gomock.Any(), "ns0", "flag").Return([]holdoutdomain.FeatureTarget{{
			GroupID: "hg", FeatureID: "flag", Status: featureproto.HoldoutGroup_WAITING,
		}}, nil)
		_, err := service.AddHoldoutGroupFeatures(createContextWithTokenRoleUnassigned(), &featureproto.AddHoldoutGroupFeaturesRequest{
			Id: "hg", EnvironmentId: "ns0",
			Features: []*featureproto.HoldoutGroupFeatureInput{{FeatureId: "flag", HoldoutVariationId: "var"}},
		})
		assert.Equal(t, codes.AlreadyExists, status.Code(err))
	})

	t.Run("flag in another active group", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectTx(service)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").Return(waitingHoldout(now), nil)
		expectFeature(service, "flag", "var", false)
		storage.EXPECT().ListLiveFeatureTargets(gomock.Any(), "ns0", "flag").Return([]holdoutdomain.FeatureTarget{{
			GroupID: "other", GroupName: "other", FeatureID: "flag", Status: featureproto.HoldoutGroup_RUNNING,
		}}, nil)
		_, err := service.AddHoldoutGroupFeatures(createContextWithTokenRoleUnassigned(), &featureproto.AddHoldoutGroupFeaturesRequest{
			Id: "hg", EnvironmentId: "ns0",
			Features: []*featureproto.HoldoutGroupFeatureInput{{FeatureId: "flag", HoldoutVariationId: "var"}},
		})
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})

	t.Run("remove a flag", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		group := waitingHoldout(now)
		group.Status = featureproto.HoldoutGroup_ANALYSIS
		group.StartAt = now - 10
		expectTx(service)
		expectHoldoutEvents(t, service, controller, 2, eventproto.Event_HOLDOUT_GROUP_FEATURE_REMOVED)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").Return(group, nil)
		// Removal only reads the flag for the event; it does not need the row lock.
		service.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().GetFeature(
			gomock.Any(), "flag", "ns0",
		).Return(&featuredomain.Feature{Feature: &featureproto.Feature{Id: "flag"}}, nil)
		storage.EXPECT().RemoveHoldoutGroupFeature(gomock.Any(), "ns0", "hg", "flag", gomock.Any()).Return(nil)
		storage.EXPECT().UpdateHoldoutGroup(gomock.Any(), gomock.Any()).Return(nil)
		// The current targets are loaded before the change, inside the transaction.
		storage.EXPECT().ListCurrentFeatures(gomock.Any(), "ns0", []string{"hg"}).Return([]*holdoutdomain.HoldoutGroupFeature{
			{HoldoutGroupFeature: &featureproto.HoldoutGroupFeature{
				HoldoutGroupId: "hg", FeatureId: "flag", HoldoutVariationId: "var",
			}},
			{HoldoutGroupFeature: &featureproto.HoldoutGroupFeature{
				HoldoutGroupId: "hg", FeatureId: "other", HoldoutVariationId: "var",
			}},
		}, nil)
		resp, err := service.RemoveHoldoutGroupFeatures(createContextWithTokenRoleUnassigned(), &featureproto.RemoveHoldoutGroupFeaturesRequest{
			Id: "hg", EnvironmentId: "ns0", FeatureIds: []string{"flag"},
		})
		require.NoError(t, err)
		require.Len(t, resp.HoldoutGroup.Features, 1)
		assert.Equal(t, "other", resp.HoldoutGroup.Features[0].FeatureId)
	})
}

func TestUpdateEnableFeatureLocksFlagRow(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	service, _ := newHoldoutService(t, controller)
	expectTx(service)
	// The whole-row save reads the flag with a row lock, never a plain read.
	service.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().GetFeatureForUpdate(
		gomock.Any(), "fid", "ns0",
	).Return(&featuredomain.Feature{Feature: &featureproto.Feature{Id: "fid"}}, nil)
	service.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().UpdateFeature(
		gomock.Any(), gomock.Any(), "ns0",
	).Return(nil)
	require.NoError(t, service.updateEnableFeature(context.Background(), "fid", "ns0", true))
}

func pendingScheduledChange(payload *featureproto.ScheduledChangePayload) *featuredomain.ScheduledFlagChange {
	return &featuredomain.ScheduledFlagChange{ScheduledFlagChange: &featureproto.ScheduledFlagChange{
		Id:            "sfc",
		FeatureId:     "fid",
		EnvironmentId: "ns0",
		Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
		Payload:       payload,
	}}
}

func TestExecuteScheduledFlagChangeLocksFlagBeforeSchedule(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	service, _ := newHoldoutService(t, controller)
	expectTx(service)
	scheduled := service.scheduledFlagChangeStorage.(*featurestoragemock.MockScheduledFlagChangeStorage)
	gomock.InOrder(
		// The flag ID is resolved with a plain read outside the transaction...
		scheduled.EXPECT().GetScheduledFlagChange(gomock.Any(), "sfc", "ns0").
			Return(pendingScheduledChange(nil), nil),
		// ...so the flag row can be locked before the schedule row, matching
		// the order archive uses (flag lock, then cancel its schedules).
		service.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().GetFeatureForUpdate(
			gomock.Any(), "fid", "ns0",
		).Return(nil, v2fs.ErrFeatureNotFound),
	)
	scheduled.EXPECT().UpdateScheduledFlagChange(gomock.Any(), gomock.Any()).Return(nil)
	_, err := service.ExecuteScheduledFlagChange(
		createContextWithTokenRoleUnassigned(),
		&featureproto.ExecuteScheduledFlagChangeRequest{Id: "sfc", EnvironmentId: "ns0"},
	)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestExecuteScheduledFlagChangeRechecksScheduleUnderLock(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	service, _ := newHoldoutService(t, controller)
	expectTx(service)
	scheduled := service.scheduledFlagChangeStorage.(*featurestoragemock.MockScheduledFlagChangeStorage)
	executed := pendingScheduledChange(nil)
	executed.Status = featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_EXECUTED
	gomock.InOrder(
		scheduled.EXPECT().GetScheduledFlagChange(gomock.Any(), "sfc", "ns0").
			Return(pendingScheduledChange(nil), nil),
		service.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().GetFeatureForUpdate(
			gomock.Any(), "fid", "ns0",
		).Return(&featuredomain.Feature{Feature: &featureproto.Feature{Id: "fid"}}, nil),
		// Another executor finished this schedule between the plain read and the lock.
		scheduled.EXPECT().GetScheduledFlagChangeForUpdate(gomock.Any(), "sfc", "ns0").
			Return(executed, nil),
	)
	_, err := service.ExecuteScheduledFlagChange(
		createContextWithTokenRoleUnassigned(),
		&featureproto.ExecuteScheduledFlagChangeRequest{Id: "sfc", EnvironmentId: "ns0"},
	)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestExecuteScheduledFlagChangeHoldoutErrorIsFailedPrecondition(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	service, storage := newHoldoutService(t, controller)
	expectTx(service)
	archive := &featureproto.ScheduledChangePayload{Archived: wrapperspb.Bool(true)}
	scheduled := service.scheduledFlagChangeStorage.(*featurestoragemock.MockScheduledFlagChangeStorage)
	scheduled.EXPECT().GetScheduledFlagChange(gomock.Any(), "sfc", "ns0").
		Return(pendingScheduledChange(archive), nil)
	scheduled.EXPECT().GetScheduledFlagChangeForUpdate(gomock.Any(), "sfc", "ns0").
		Return(pendingScheduledChange(archive), nil)
	// Once by the executor, once by updateFeatureWithinTransaction; both are locking reads.
	service.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().GetFeatureForUpdate(
		gomock.Any(), "fid", "ns0",
	).Return(&featuredomain.Feature{Feature: &featureproto.Feature{Id: "fid"}}, nil).Times(2)
	service.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().ListFeatures(
		gomock.Any(), gomock.Any(),
	).Return([]*featureproto.Feature{{Id: "fid"}}, 0, int64(0), nil)
	storage.EXPECT().ListLiveFeatureTargets(gomock.Any(), "ns0", "fid").Return([]holdoutdomain.FeatureTarget{{
		GroupID: "hg", GroupName: "holdout", FeatureID: "fid",
	}}, nil)
	scheduled.EXPECT().UpdateScheduledFlagChange(gomock.Any(), gomock.Any()).Return(nil)

	_, err := service.ExecuteScheduledFlagChange(
		createContextWithTokenRoleUnassigned(),
		&featureproto.ExecuteScheduledFlagChangeRequest{Id: "sfc", EnvironmentId: "ns0"},
	)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestHoldoutGuardsOnFeatureAndSegment(t *testing.T) {
	t.Parallel()
	ctx := createContextWithTokenRoleUnassigned()

	t.Run("nil storage allows variation removal", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service := createFeatureService(controller)
		err := service.rejectHoldoutVariationRemoval(context.Background(), "ns0", "fid", []*featureproto.VariationChange{{
			ChangeType: featureproto.ChangeType_DELETE,
			Variation:  &featureproto.Variation{Id: "var"},
		}})
		assert.NoError(t, err)
	})

	t.Run("archive is rejected while the flag is a holdout target", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectEnvironment(service)
		expectTx(service)
		service.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().GetFeatureForUpdate(
			gomock.Any(), "fid", "ns0",
		).Return(&featuredomain.Feature{Feature: &featureproto.Feature{Id: "fid"}}, nil)
		service.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().ListFeatures(
			gomock.Any(), gomock.Any(),
		).Return([]*featureproto.Feature{{Id: "fid"}}, 0, int64(0), nil)
		storage.EXPECT().ListLiveFeatureTargets(gomock.Any(), "ns0", "fid").Return([]holdoutdomain.FeatureTarget{{
			GroupID: "hg", GroupName: "holdout", FeatureID: "fid",
		}}, nil)
		_, err := service.UpdateFeature(ctx, &featureproto.UpdateFeatureRequest{
			Id: "fid", EnvironmentId: "ns0", Comment: "archive", Archived: wrapperspb.Bool(true),
		})
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})

	t.Run("variation deletion is rejected when it is the holdout variation", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectEnvironment(service)
		expectTx(service)
		service.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().GetFeatureForUpdate(
			gomock.Any(), "fid", "ns0",
		).Return(&featuredomain.Feature{Feature: &featureproto.Feature{Id: "fid"}}, nil)
		service.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().ListFeatures(
			gomock.Any(), gomock.Any(),
		).Return([]*featureproto.Feature{{Id: "fid"}}, 0, int64(0), nil)
		storage.EXPECT().ListLiveFeatureTargets(gomock.Any(), "ns0", "fid").Return([]holdoutdomain.FeatureTarget{{
			GroupID: "hg", GroupName: "holdout", FeatureID: "fid", VariationID: "var",
		}}, nil)
		_, err := service.UpdateFeature(ctx, &featureproto.UpdateFeatureRequest{
			Id: "fid", EnvironmentId: "ns0", Comment: "delete variation",
			VariationChanges: []*featureproto.VariationChange{{
				ChangeType: featureproto.ChangeType_DELETE,
				Variation:  &featureproto.Variation{Id: "var"},
			}},
		})
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})

	t.Run("delete is rejected while the flag is a holdout target", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectEnvironment(service)
		expectTx(service)
		service.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().GetFeatureForUpdate(
			gomock.Any(), "fid", "ns0",
		).Return(&featuredomain.Feature{Feature: &featureproto.Feature{Id: "fid"}}, nil)
		storage.EXPECT().ListLiveFeatureTargets(gomock.Any(), "ns0", "fid").Return([]holdoutdomain.FeatureTarget{{
			GroupID: "hg", GroupName: "holdout", FeatureID: "fid",
		}}, nil)
		_, err := service.DeleteFeature(ctx, &featureproto.DeleteFeatureRequest{
			Id: "fid", EnvironmentId: "ns0",
		})
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})

	t.Run("serving holdout locks the segment", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectTx(service)
		expectTx(service)
		expectTx(service)
		// Each writer locks the segment row first, then checks serving groups.
		service.segmentStorage.(*featurestoragemock.MockSegmentStorage).EXPECT().GetSegmentForUpdate(
			gomock.Any(), "seg", "ns0",
		).Return(&featuredomain.Segment{Segment: &featureproto.Segment{Id: "seg"}}, nil, nil).Times(3)
		storage.EXPECT().ListActiveGroups(gomock.Any(), "ns0").Return([]holdoutdomain.GroupSnapshot{{
			ID: "hg", Name: "holdout", SegmentID: "seg", Status: featureproto.HoldoutGroup_RUNNING,
		}}, nil).Times(3)
		service.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().ListFeatures(
			gomock.Any(), gomock.Any(),
		).Return([]*featureproto.Feature{}, 0, int64(0), nil)

		_, err := service.DeleteSegment(ctx, &featureproto.DeleteSegmentRequest{Id: "seg", EnvironmentId: "ns0"})
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
		_, err = service.UpdateSegment(ctx, &featureproto.UpdateSegmentRequest{
			Id: "seg", EnvironmentId: "ns0", Name: wrapperspb.String("renamed"),
		})
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
		_, err = service.BulkUploadSegmentUsers(ctx, &featureproto.BulkUploadSegmentUsersRequest{
			SegmentId: "seg", EnvironmentId: "ns0", Data: []byte("user"), State: featureproto.SegmentUser_INCLUDED,
		})
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	})
}

func newHoldoutService(t *testing.T, controller *gomock.Controller) (*FeatureService, *holdoutmock.MockHoldoutStorage) {
	t.Helper()
	service := createFeatureService(controller)
	storage := holdoutmock.NewMockHoldoutStorage(controller)
	service.holdoutStorage = storage
	return service, storage
}

func expectTx(service *FeatureService) {
	service.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
		gomock.Any(), gomock.Any(),
	).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})
}

func expectHoldoutEvents(
	t *testing.T,
	service *FeatureService,
	controller *gomock.Controller,
	count int,
	eventType eventproto.Event_Type,
) {
	t.Helper()
	pub := publishermock.NewMockPublisher(controller)
	service.domainPublisher = pub
	pub.EXPECT().PublishMulti(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, messages []publisher.Message) map[string]error {
			require.Len(t, messages, count)
			for _, message := range messages {
				event := message.(*eventproto.Event)
				assert.Equal(t, eventType, event.Type)
				switch event.EntityType {
				case eventproto.Event_FEATURE:
					var decoded featureproto.Feature
					require.NoError(t, protojson.Unmarshal([]byte(event.EntityData), &decoded))
					assert.NotEmpty(t, decoded.Id)
					// The flag itself did not change: both snapshots must be the same.
					assert.Equal(t, event.EntityData, event.PreviousEntityData)
				case eventproto.Event_HOLDOUT_GROUP:
					assertHoldoutTargetSnapshots(t, event)
				}
			}
			return nil
		},
	)
}

// assertHoldoutTargetSnapshots checks that a target add/remove event records
// the membership change in its before/after group snapshots.
func assertHoldoutTargetSnapshots(t *testing.T, event *eventproto.Event) {
	t.Helper()
	if event.Type != eventproto.Event_HOLDOUT_GROUP_FEATURE_ADDED &&
		event.Type != eventproto.Event_HOLDOUT_GROUP_FEATURE_REMOVED {
		return
	}
	var current, previous featureproto.HoldoutGroup
	require.NoError(t, protojson.Unmarshal([]byte(event.EntityData), &current))
	require.NoError(t, protojson.Unmarshal([]byte(event.PreviousEntityData), &previous))
	assert.NotEqual(t, previous.UpdatedAt, current.UpdatedAt)
	if event.Type == eventproto.Event_HOLDOUT_GROUP_FEATURE_ADDED {
		assert.Len(t, current.Features, len(previous.Features)+1)
		return
	}
	assert.Len(t, current.Features, len(previous.Features)-1)
}

func expectTransition(
	t *testing.T,
	service *FeatureService,
	storage *holdoutmock.MockHoldoutStorage,
	controller *gomock.Controller,
	group *holdoutdomain.HoldoutGroup,
	active []holdoutdomain.GroupSnapshot,
	eventType eventproto.Event_Type,
) {
	t.Helper()
	expectTx(service)
	expectHoldoutEvents(t, service, controller, 1, eventType)
	storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
	storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").Return(group, nil)
	storage.EXPECT().ListActiveGroups(gomock.Any(), "ns0").Return(active, nil)
	storage.EXPECT().UpdateHoldoutGroup(gomock.Any(), gomock.Any()).Return(nil)
	storage.EXPECT().ListCurrentFeatures(gomock.Any(), "ns0", []string{"hg"}).Return(nil, nil)
}

func expectGoal(service *FeatureService, goalID string) {
	service.experimentClient.(*experimentclientmock.MockClient).EXPECT().GetGoal(
		gomock.Any(), gomock.Any(),
	).Return(&experimentproto.GetGoalResponse{Goal: &experimentproto.Goal{Id: goalID}}, nil)
}

func expectFeature(service *FeatureService, featureID, variationID string, archived bool) {
	service.featureStorage.(*featurestoragemock.MockFeatureStorage).EXPECT().GetFeatureForUpdate(
		gomock.Any(), featureID, "ns0",
	).Return(&featuredomain.Feature{Feature: &featureproto.Feature{
		Id:         featureID,
		Archived:   archived,
		Variations: []*featureproto.Variation{{Id: variationID}},
	}}, nil)
}

func expectEnvironment(service *FeatureService) {
	service.environmentClient.(*envclientmock.MockClient).EXPECT().GetEnvironmentV2(
		gomock.Any(), gomock.Any(),
	).Return(&envproto.GetEnvironmentV2Response{Environment: &envproto.EnvironmentV2{}}, nil)
}

func hashCreateRequest(now int64, bucket *wrapperspb.DoubleValue) *featureproto.CreateHoldoutGroupRequest {
	return &featureproto.CreateHoldoutGroupRequest{
		EnvironmentId:  "ns0",
		Name:           "Q1 holdout",
		AssignmentType: featureproto.HoldoutGroup_HASH,
		Percentage:     10,
		StartAt:        now + 100,
		StopAt:         now + 1_000,
		BucketStart:    bucket,
	}
}

func waitingHoldout(now int64) *holdoutdomain.HoldoutGroup {
	return &holdoutdomain.HoldoutGroup{HoldoutGroup: &featureproto.HoldoutGroup{
		Id:             "hg",
		EnvironmentId:  "ns0",
		Name:           "holdout",
		Status:         featureproto.HoldoutGroup_WAITING,
		AssignmentType: featureproto.HoldoutGroup_HASH,
		Percentage:     10,
		SamplingSeed:   "seed",
		StartAt:        now + 100,
		StopAt:         now + 1_000,
		GoalIds:        []string{},
	}}
}
