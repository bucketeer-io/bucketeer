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
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestCreateScheduledFlagChange_ValidationErrors(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		setup       func(*FeatureService)
		req         *featureproto.CreateScheduledFlagChangeRequest
		expectedErr error
	}{
		{
			desc:  "missing feature id",
			setup: func(s *FeatureService) {},
			req: &featureproto.CreateScheduledFlagChangeRequest{
				EnvironmentId: "ns0",
				FeatureId:     "",
				ScheduledAt:   time.Now().Add(time.Hour).Unix(),
				Payload:       &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(true)},
			},
			expectedErr: statusMissingFeatureID.Err(),
		},
		{
			desc:  "missing scheduled at",
			setup: func(s *FeatureService) {},
			req: &featureproto.CreateScheduledFlagChangeRequest{
				EnvironmentId: "ns0",
				FeatureId:     "feature-id",
				ScheduledAt:   0,
				Payload:       &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(true)},
			},
			expectedErr: statusMissingScheduledAt.Err(),
		},
		{
			desc:  "missing payload",
			setup: func(s *FeatureService) {},
			req: &featureproto.CreateScheduledFlagChangeRequest{
				EnvironmentId: "ns0",
				FeatureId:     "feature-id",
				ScheduledAt:   time.Now().Add(time.Hour).Unix(),
				Payload:       nil,
			},
			expectedErr: statusMissingPayload.Err(),
		},
		{
			desc:  "scheduled time too soon",
			setup: func(s *FeatureService) {},
			req: &featureproto.CreateScheduledFlagChangeRequest{
				EnvironmentId: "ns0",
				FeatureId:     "feature-id",
				ScheduledAt:   time.Now().Add(30 * time.Second).Unix(), // 30s from now, less than 1 minute
				Payload:       &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(true)},
			},
			expectedErr: statusScheduledTimeTooSoon.Err(),
		},
		{
			desc:  "empty payload",
			setup: func(s *FeatureService) {},
			req: &featureproto.CreateScheduledFlagChangeRequest{
				EnvironmentId: "ns0",
				FeatureId:     "feature-id",
				ScheduledAt:   time.Now().Add(time.Hour).Unix(),
				Payload:       &featureproto.ScheduledChangePayload{}, // Empty payload
			},
			expectedErr: statusEmptyPayload.Err(),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createFeatureServiceWithGetAccountByEnvironmentMock(gomock.NewController(t), accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_EDITOR)
			p.setup(service)
			ctx := createContextWithToken()
			_, err := service.CreateScheduledFlagChange(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCreateScheduledFlagChange_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := createFeatureServiceWithGetAccountByEnvironmentMock(ctrl, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_EDITOR)

	mysqlClient := service.mysqlClient.(*mysqlmock.MockClient)
	featureStorage := service.featureStorage.(*mock.MockFeatureStorage)
	scheduledStorage := service.scheduledFlagChangeStorage.(*mock.MockScheduledFlagChangeStorage)
	domainPublisher := service.domainPublisher.(*publishermock.MockPublisher)

	feature := &domain.Feature{
		Feature: &featureproto.Feature{
			Id:      "feature-id",
			Name:    "Test Feature",
			Version: 1,
			Variations: []*featureproto.Variation{
				{Id: "var-1", Name: "Variation 1", Value: "true"},
				{Id: "var-2", Name: "Variation 2", Value: "false"},
			},
		},
	}

	// Mock transaction - execute the function passed to it
	mysqlClient.EXPECT().RunInTransactionV2(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(context.Context, mysql.Transaction) error) error {
			return f(ctx, nil)
		},
	)
	featureStorage.EXPECT().GetFeature(gomock.Any(), "feature-id", "ns0").Return(feature, nil)
	// ListScheduledFlagChanges is called twice: count check and conflict detection
	scheduledStorage.EXPECT().ListScheduledFlagChanges(gomock.Any(), gomock.Any()).
		Return([]*featureproto.ScheduledFlagChange{}, 0, int64(0), nil).
		Times(2)
	scheduledStorage.EXPECT().CreateScheduledFlagChange(gomock.Any(), gomock.Any()).Return(nil)
	// Mock domain event publishing
	domainPublisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)

	ctx := createContextWithToken()
	req := &featureproto.CreateScheduledFlagChangeRequest{
		EnvironmentId: "ns0",
		FeatureId:     "feature-id",
		ScheduledAt:   time.Now().Add(time.Hour).Unix(),
		Timezone:      "Asia/Tokyo",
		Payload:       &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(true)},
		Comment:       "Enable flag in 1 hour",
	}

	resp, err := service.CreateScheduledFlagChange(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, resp.ScheduledFlagChange)
	assert.Equal(t, "feature-id", resp.ScheduledFlagChange.FeatureId)
	assert.Equal(t, featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING, resp.ScheduledFlagChange.Status)
	// Check that change summaries contain the enable flag message key
	require.Len(t, resp.ScheduledFlagChange.ChangeSummaries, 1)
	assert.Equal(t, "ScheduledChange.EnableFlag", resp.ScheduledFlagChange.ChangeSummaries[0].MessageKey)
}

func TestGetScheduledFlagChange_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := createFeatureServiceWithGetAccountByEnvironmentMock(ctrl, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER)

	scheduledStorage := service.scheduledFlagChangeStorage.(*mock.MockScheduledFlagChangeStorage)
	featureStorage := service.featureStorage.(*mock.MockFeatureStorage)

	sfc := &domain.ScheduledFlagChange{
		ScheduledFlagChange: &featureproto.ScheduledFlagChange{
			Id:            "sfc-id",
			FeatureId:     "feature-id",
			EnvironmentId: "ns0",
			ScheduledAt:   time.Now().Add(time.Hour).Unix(),
			Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
			Payload:       &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(true)},
		},
	}

	feature := &domain.Feature{
		Feature: &featureproto.Feature{
			Id:      "feature-id",
			Version: 1,
		},
	}

	scheduledStorage.EXPECT().GetScheduledFlagChange(gomock.Any(), "sfc-id", "ns0").Return(sfc, nil)
	featureStorage.EXPECT().GetFeature(gomock.Any(), "feature-id", "ns0").Return(feature, nil)

	ctx := createContextWithToken()
	req := &featureproto.GetScheduledFlagChangeRequest{
		EnvironmentId: "ns0",
		Id:            "sfc-id",
	}

	resp, err := service.GetScheduledFlagChange(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, resp.ScheduledFlagChange)
	assert.Equal(t, "sfc-id", resp.ScheduledFlagChange.Id)
}

func TestListScheduledFlagChanges_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := createFeatureServiceWithGetAccountByEnvironmentMock(ctrl, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER)

	scheduledStorage := service.scheduledFlagChangeStorage.(*mock.MockScheduledFlagChangeStorage)
	featureStorage := service.featureStorage.(*mock.MockFeatureStorage)

	sfcs := []*featureproto.ScheduledFlagChange{
		{
			Id:            "sfc-1",
			FeatureId:     "feature-id",
			EnvironmentId: "ns0",
			ScheduledAt:   time.Now().Add(time.Hour).Unix(),
			Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
			Payload:       &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(true)},
		},
		{
			Id:            "sfc-2",
			FeatureId:     "feature-id",
			EnvironmentId: "ns0",
			ScheduledAt:   time.Now().Add(2 * time.Hour).Unix(),
			Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
			Payload:       &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(false)},
		},
	}

	feature := &domain.Feature{
		Feature: &featureproto.Feature{
			Id:      "feature-id",
			Version: 1,
		},
	}

	scheduledStorage.EXPECT().ListScheduledFlagChanges(gomock.Any(), gomock.Any()).Return(sfcs, 2, int64(2), nil)
	featureStorage.EXPECT().GetFeature(gomock.Any(), "feature-id", "ns0").Return(feature, nil).Times(2)

	ctx := createContextWithToken()
	req := &featureproto.ListScheduledFlagChangesRequest{
		EnvironmentId: "ns0",
		FeatureId:     "feature-id",
		PageSize:      10,
	}

	resp, err := service.ListScheduledFlagChanges(ctx, req)
	require.NoError(t, err)
	assert.Len(t, resp.ScheduledFlagChanges, 2)
	assert.Equal(t, int64(2), resp.TotalCount)
}

func TestDeleteScheduledFlagChange_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := createFeatureServiceWithGetAccountByEnvironmentMock(ctrl, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_EDITOR)

	scheduledStorage := service.scheduledFlagChangeStorage.(*mock.MockScheduledFlagChangeStorage)
	featureStorage := service.featureStorage.(*mock.MockFeatureStorage)
	mysqlClient := service.mysqlClient.(*mysqlmock.MockClient)
	domainPublisher := service.domainPublisher.(*publishermock.MockPublisher)

	sfc := &domain.ScheduledFlagChange{
		ScheduledFlagChange: &featureproto.ScheduledFlagChange{
			Id:            "sfc-id",
			FeatureId:     "feature-id",
			EnvironmentId: "ns0",
			Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
		},
	}

	feature := &domain.Feature{
		Feature: &featureproto.Feature{
			Id:   "feature-id",
			Name: "Test Feature",
		},
	}

	mysqlClient.EXPECT().RunInTransactionV2(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(context.Context, mysql.Transaction) error) error {
			return f(ctx, nil)
		},
	)
	scheduledStorage.EXPECT().GetScheduledFlagChange(gomock.Any(), "sfc-id", "ns0").Return(sfc, nil)
	featureStorage.EXPECT().GetFeature(gomock.Any(), "feature-id", "ns0").Return(feature, nil)
	scheduledStorage.EXPECT().UpdateScheduledFlagChange(gomock.Any(), gomock.Any()).Return(nil)
	// Mock domain event publishing
	domainPublisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)

	ctx := createContextWithToken()
	req := &featureproto.DeleteScheduledFlagChangeRequest{
		EnvironmentId: "ns0",
		Id:            "sfc-id",
	}

	resp, err := service.DeleteScheduledFlagChange(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestGetScheduledFlagChangeSummary_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := createFeatureServiceWithGetAccountByEnvironmentMock(ctrl, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER)

	scheduledStorage := service.scheduledFlagChangeStorage.(*mock.MockScheduledFlagChangeStorage)

	scheduledAt := time.Now().Add(time.Hour).Unix()
	sfcs := []*featureproto.ScheduledFlagChange{
		{
			Id:            "sfc-1",
			FeatureId:     "feature-id",
			EnvironmentId: "ns0",
			ScheduledAt:   scheduledAt,
			Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
			Payload:       &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(true)},
		},
		{
			Id:            "sfc-2",
			FeatureId:     "feature-id",
			EnvironmentId: "ns0",
			ScheduledAt:   time.Now().Add(2 * time.Hour).Unix(),
			Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT,
			Payload:       &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(false)},
		},
	}

	scheduledStorage.EXPECT().ListScheduledFlagChanges(gomock.Any(), gomock.Any()).Return(sfcs, 0, int64(2), nil)

	ctx := createContextWithToken()
	req := &featureproto.GetScheduledFlagChangeSummaryRequest{
		EnvironmentId: "ns0",
		FeatureId:     "feature-id",
	}

	resp, err := service.GetScheduledFlagChangeSummary(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, resp.Summary)
	assert.Equal(t, int32(1), resp.Summary.PendingCount)
	assert.Equal(t, int32(1), resp.Summary.ConflictCount)
	assert.Equal(t, scheduledAt, resp.Summary.NextScheduledAt)
}

func TestValidateScheduledTime(t *testing.T) {
	t.Parallel()

	now := time.Now().Unix()

	patterns := []struct {
		desc        string
		scheduledAt int64
		expectError bool
	}{
		{
			desc:        "too soon - 30 seconds from now",
			scheduledAt: now + 30,
			expectError: true,
		},
		{
			desc:        "valid - exactly 1 minute from now",
			scheduledAt: now + 60,
			expectError: false,
		},
		{
			desc:        "valid - 2 minutes from now",
			scheduledAt: now + 2*60,
			expectError: false,
		},
		{
			desc:        "valid - 1 hour from now",
			scheduledAt: now + 60*60,
			expectError: false,
		},
		{
			desc:        "too far - more than 1 year",
			scheduledAt: now + 366*24*60*60,
			expectError: true,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := validateScheduledTime(p.scheduledAt)
			if p.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCountChanges(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc          string
		payload       *featureproto.ScheduledChangePayload
		expectedCount int
	}{
		{
			desc:          "nil payload",
			payload:       nil,
			expectedCount: 0,
		},
		{
			desc:          "empty payload",
			payload:       &featureproto.ScheduledChangePayload{},
			expectedCount: 0,
		},
		{
			desc: "single enabled change",
			payload: &featureproto.ScheduledChangePayload{
				Enabled: wrapperspb.Bool(true),
			},
			expectedCount: 1,
		},
		{
			desc: "multiple changes",
			payload: &featureproto.ScheduledChangePayload{
				Enabled:     wrapperspb.Bool(true),
				Name:        wrapperspb.String("new-name"),
				Description: wrapperspb.String("new-desc"),
				RuleChanges: []*featureproto.RuleChange{
					{ChangeType: featureproto.ChangeType_CREATE},
				},
			},
			expectedCount: 4,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sfc := &domain.ScheduledFlagChange{
				ScheduledFlagChange: &featureproto.ScheduledFlagChange{Payload: p.payload},
			}
			assert.Equal(t, p.expectedCount, sfc.CountChanges())
		})
	}
}

func TestValidateScheduleGap(t *testing.T) {
	t.Parallel()

	baseTime := time.Now().Add(2 * time.Hour).Unix()
	gapSec := int64(minScheduleGapBetweenMinutes * 60) // 300s = 5 min

	existing := []*featureproto.ScheduledFlagChange{
		{
			Id:          "sfc-1",
			ScheduledAt: baseTime,
		},
		{
			Id:          "sfc-2",
			ScheduledAt: baseTime + 10*60, // +10 min
		},
	}

	patterns := []struct {
		desc        string
		scheduledAt int64
		excludeID   string
		expectedErr error
	}{
		{
			desc:        "exactly at existing schedule",
			scheduledAt: baseTime,
			expectedErr: statusScheduledTimeTooClose.Err(),
		},
		{
			desc:        "1 minute after existing schedule",
			scheduledAt: baseTime + 60,
			expectedErr: statusScheduledTimeTooClose.Err(),
		},
		{
			desc:        "4 minutes after existing schedule",
			scheduledAt: baseTime + 4*60,
			expectedErr: statusScheduledTimeTooClose.Err(),
		},
		{
			desc:        "exactly 5 minutes after existing schedule",
			scheduledAt: baseTime + gapSec,
			expectedErr: nil,
		},
		{
			desc:        "6 minutes after existing (but within 5min of sfc-2)",
			scheduledAt: baseTime + 6*60,
			expectedErr: statusScheduledTimeTooClose.Err(),
		},
		{
			desc:        "between the two with enough gap from both",
			scheduledAt: baseTime + gapSec, // 5 min after sfc-1, 5 min before sfc-2
			expectedErr: nil,
		},
		{
			desc:        "well after all existing schedules",
			scheduledAt: baseTime + 20*60,
			expectedErr: nil,
		},
		{
			desc:        "too close but excluded (updating same schedule)",
			scheduledAt: baseTime + 60,
			excludeID:   "sfc-1",
			expectedErr: nil, // sfc-1 is excluded; 1min+60s from sfc-2 at +10min = ~9min gap, OK
		},
		{
			desc:        "no existing schedules",
			scheduledAt: baseTime,
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			input := existing
			if p.desc == "no existing schedules" {
				input = nil
			}
			err := validateScheduleGap(p.scheduledAt, input, p.excludeID)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateScheduledFlagChange_ScheduleGapTooClose(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := createFeatureServiceWithGetAccountByEnvironmentMock(
		ctrl,
		accountproto.AccountV2_Role_Organization_MEMBER,
		accountproto.AccountV2_Role_Environment_EDITOR,
	)

	mysqlClient := service.mysqlClient.(*mysqlmock.MockClient)
	featureStorage := service.featureStorage.(*mock.MockFeatureStorage)
	scheduledStorage := service.scheduledFlagChangeStorage.(*mock.MockScheduledFlagChangeStorage)

	scheduledAt := time.Now().Add(2 * time.Hour).Unix()

	feature := &domain.Feature{
		Feature: &featureproto.Feature{
			Id:      "feature-id",
			Name:    "Test Feature",
			Version: 1,
			Variations: []*featureproto.Variation{
				{Id: "var-1", Value: "true", Name: "V1"},
				{Id: "var-2", Value: "false", Name: "V2"},
			},
		},
	}

	// Existing schedule at scheduledAt + 2 minutes (within the 5-min gap)
	existingSchedule := &featureproto.ScheduledFlagChange{
		Id:          "sfc-existing",
		FeatureId:   "feature-id",
		ScheduledAt: scheduledAt + 2*60,
		Status:      featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
	}

	mysqlClient.EXPECT().RunInTransactionV2(
		gomock.Any(), gomock.Any(),
	).DoAndReturn(
		func(
			ctx context.Context,
			f func(context.Context, mysql.Transaction) error,
		) error {
			return f(ctx, nil)
		},
	)
	featureStorage.EXPECT().GetFeature(
		gomock.Any(), "feature-id", "ns0",
	).Return(feature, nil)
	// listPendingSchedulesForFeature returns one schedule too close
	scheduledStorage.EXPECT().ListScheduledFlagChanges(
		gomock.Any(), gomock.Any(),
	).Return(
		[]*featureproto.ScheduledFlagChange{existingSchedule},
		1, int64(1), nil,
	)

	ctx := createContextWithToken()
	req := &featureproto.CreateScheduledFlagChangeRequest{
		EnvironmentId: "ns0",
		FeatureId:     "feature-id",
		ScheduledAt:   scheduledAt,
		Timezone:      "UTC",
		Payload: &featureproto.ScheduledChangePayload{
			Enabled: wrapperspb.Bool(true),
		},
		Comment: "Too close to existing schedule",
	}

	_, err := service.CreateScheduledFlagChange(ctx, req)
	assert.Equal(t, statusScheduledTimeTooClose.Err(), err)
}

func TestCreateScheduledFlagChange_CircularPrerequisite(t *testing.T) {
	t.Parallel()

	flagA := &domain.Feature{
		Feature: &featureproto.Feature{
			Id:      "flag-a",
			Name:    "Flag A",
			Version: 1,
			Variations: []*featureproto.Variation{
				{Id: "var-a1", Value: "true", Name: "A-True"},
				{Id: "var-a2", Value: "false", Name: "A-False"},
			},
		},
	}
	flagBNoDeps := &featureproto.Feature{
		Id:      "flag-b",
		Name:    "Flag B",
		Version: 1,
		Variations: []*featureproto.Variation{
			{Id: "var-b1", Value: "true", Name: "B-True"},
			{Id: "var-b2", Value: "false", Name: "B-False"},
		},
	}
	flagBDependsOnA := &featureproto.Feature{
		Id:      "flag-b",
		Name:    "Flag B",
		Version: 1,
		Variations: []*featureproto.Variation{
			{Id: "var-b1", Value: "true", Name: "B-True"},
			{Id: "var-b2", Value: "false", Name: "B-False"},
		},
		Prerequisites: []*featureproto.Prerequisite{
			{FeatureId: "flag-a", VariationId: "var-a1"},
		},
	}
	flagC := &featureproto.Feature{
		Id:      "flag-c",
		Name:    "Flag C",
		Version: 1,
		Variations: []*featureproto.Variation{
			{Id: "var-c1", Value: "true", Name: "C-True"},
			{Id: "var-c2", Value: "false", Name: "C-False"},
		},
		Prerequisites: []*featureproto.Prerequisite{
			{FeatureId: "flag-b", VariationId: "var-b1"},
		},
	}

	patterns := []struct {
		desc              string
		flagB             *featureproto.Feature
		prerequisiteFlag  string
		prerequisiteVarID string
		setupMocks        func(*mock.MockFeatureStorage, *mock.MockScheduledFlagChangeStorage, *publishermock.MockPublisher)
		expectedErr       error
	}{
		{
			desc:              "cycle detected: Flag A -> Flag B -> Flag A",
			flagB:             flagBDependsOnA,
			prerequisiteFlag:  "flag-b",
			prerequisiteVarID: "var-b1",
			setupMocks: func(
				fs *mock.MockFeatureStorage,
				ss *mock.MockScheduledFlagChangeStorage,
				_ *publishermock.MockPublisher,
			) {
				// GetFeature for flag-a (inside transaction)
				fs.EXPECT().GetFeature(
					gomock.Any(), "flag-a", "ns0",
				).Return(flagA, nil)
				// Prerequisite validation: GetFeature for flag-b
				fs.EXPECT().GetFeature(
					gomock.Any(), "flag-b", "ns0",
				).Return(
					&domain.Feature{Feature: flagBDependsOnA}, nil,
				)
				// Circular check: ListFeatures
				fs.EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return(
					[]*featureproto.Feature{
						flagA.Feature, flagBDependsOnA,
					},
					0, int64(0), nil,
				)
				// Validation fails → count/conflict/create never called
			},
			expectedErr: statusCircularPrerequisiteDetected.Err(),
		},
		{
			desc:              "no cycle: Flag A -> Flag B (B has no deps)",
			flagB:             flagBNoDeps,
			prerequisiteFlag:  "flag-b",
			prerequisiteVarID: "var-b1",
			setupMocks: func(
				fs *mock.MockFeatureStorage,
				ss *mock.MockScheduledFlagChangeStorage,
				dp *publishermock.MockPublisher,
			) {
				// GetFeature for flag-a
				fs.EXPECT().GetFeature(
					gomock.Any(), "flag-a", "ns0",
				).Return(flagA, nil)
				// Prerequisite validation: GetFeature for flag-b
				fs.EXPECT().GetFeature(
					gomock.Any(), "flag-b", "ns0",
				).Return(
					&domain.Feature{Feature: flagBNoDeps}, nil,
				)
				// Circular check: ListFeatures
				fs.EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return(
					[]*featureproto.Feature{
						flagA.Feature, flagBNoDeps,
					},
					0, int64(0), nil,
				)
				// Count pending schedules
				ss.EXPECT().ListScheduledFlagChanges(
					gomock.Any(), gomock.Any(),
				).Return(
					[]*featureproto.ScheduledFlagChange{},
					0, int64(0), nil,
				)
				// Conflict detection
				ss.EXPECT().ListScheduledFlagChanges(
					gomock.Any(), gomock.Any(),
				).Return(
					[]*featureproto.ScheduledFlagChange{},
					0, int64(0), nil,
				)
				// Create + publish
				ss.EXPECT().CreateScheduledFlagChange(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				dp.EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			expectedErr: nil,
		},
		{
			desc:              "transitive cycle detected: Flag A -> Flag C -> Flag B -> Flag A",
			flagB:             flagC,
			prerequisiteFlag:  "flag-c",
			prerequisiteVarID: "var-c1",
			setupMocks: func(
				fs *mock.MockFeatureStorage,
				ss *mock.MockScheduledFlagChangeStorage,
				_ *publishermock.MockPublisher,
			) {
				// GetFeature for flag-a (inside transaction)
				fs.EXPECT().GetFeature(
					gomock.Any(), "flag-a", "ns0",
				).Return(flagA, nil)
				// Prerequisite validation: GetFeature for flag-c
				fs.EXPECT().GetFeature(
					gomock.Any(), "flag-c", "ns0",
				).Return(
					&domain.Feature{Feature: flagC}, nil,
				)
				// Circular check: ListFeatures (returns all flags in environment)
				fs.EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(),
				).Return(
					[]*featureproto.Feature{
						flagA.Feature, flagBDependsOnA, flagC,
					},
					0, int64(0), nil,
				)
				// Validation fails → count/conflict/create never called
			},
			expectedErr: statusCircularPrerequisiteDetected.Err(),
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := createFeatureServiceWithGetAccountByEnvironmentMock(
				ctrl,
				accountproto.AccountV2_Role_Organization_MEMBER,
				accountproto.AccountV2_Role_Environment_EDITOR,
			)

			mysqlClient := service.mysqlClient.(*mysqlmock.MockClient)
			featureStorage := service.featureStorage.(*mock.MockFeatureStorage)
			scheduledStorage := service.scheduledFlagChangeStorage.(*mock.MockScheduledFlagChangeStorage)
			domainPublisher := service.domainPublisher.(*publishermock.MockPublisher)

			mysqlClient.EXPECT().RunInTransactionV2(
				gomock.Any(), gomock.Any(),
			).DoAndReturn(
				func(
					ctx context.Context,
					f func(context.Context, mysql.Transaction) error,
				) error {
					return f(ctx, nil)
				},
			)
			p.setupMocks(featureStorage, scheduledStorage, domainPublisher)

			ctx := createContextWithToken()
			req := &featureproto.CreateScheduledFlagChangeRequest{
				EnvironmentId: "ns0",
				FeatureId:     "flag-a",
				ScheduledAt:   time.Now().Add(time.Hour).Unix(),
				Timezone:      "UTC",
				Payload: &featureproto.ScheduledChangePayload{
					PrerequisiteChanges: []*featureproto.PrerequisiteChange{
						{
							ChangeType: featureproto.ChangeType_CREATE,
							Prerequisite: &featureproto.Prerequisite{
								FeatureId:   p.prerequisiteFlag,
								VariationId: p.prerequisiteVarID,
							},
						},
					},
				},
				Comment: "prerequisite cycle test",
			}

			resp, err := service.CreateScheduledFlagChange(ctx, req)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp.ScheduledFlagChange)
				assert.Equal(t, "flag-a", resp.ScheduledFlagChange.FeatureId)
			}
		})
	}
}
