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
				ScheduledAt:   time.Now().Add(time.Minute).Unix(), // 1 minute from now, less than 5 minutes
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

	featureStorage := service.featureStorage.(*mock.MockFeatureStorage)
	scheduledStorage := service.scheduledFlagChangeStorage.(*mock.MockScheduledFlagChangeStorage)

	feature := &domain.Feature{
		Feature: &featureproto.Feature{
			Id:      "feature-id",
			Version: 1,
			Variations: []*featureproto.Variation{
				{Id: "var-1", Name: "Variation 1", Value: "true"},
				{Id: "var-2", Name: "Variation 2", Value: "false"},
			},
		},
	}

	featureStorage.EXPECT().GetFeature(gomock.Any(), "feature-id", "ns0").Return(feature, nil)
	scheduledStorage.EXPECT().ListScheduledFlagChanges(gomock.Any(), gomock.Any()).Return([]*featureproto.ScheduledFlagChange{}, 0, int64(0), nil)
	scheduledStorage.EXPECT().CreateScheduledFlagChange(gomock.Any(), gomock.Any()).Return(nil)

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
	mysqlClient := service.mysqlClient.(*mysqlmock.MockClient)

	sfc := &domain.ScheduledFlagChange{
		ScheduledFlagChange: &featureproto.ScheduledFlagChange{
			Id:            "sfc-id",
			FeatureId:     "feature-id",
			EnvironmentId: "ns0",
			Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
		},
	}

	mysqlClient.EXPECT().RunInTransactionV2(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(context.Context, mysql.Transaction) error) error {
			return f(ctx, nil)
		},
	)
	scheduledStorage.EXPECT().GetScheduledFlagChange(gomock.Any(), "sfc-id", "ns0").Return(sfc, nil)
	scheduledStorage.EXPECT().UpdateScheduledFlagChange(gomock.Any(), gomock.Any()).Return(nil)

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
			desc:        "too soon - 1 minute from now",
			scheduledAt: now + 60,
			expectError: true,
		},
		{
			desc:        "too soon - 4 minutes from now",
			scheduledAt: now + 4*60,
			expectError: true,
		},
		{
			desc:        "valid - 6 minutes from now",
			scheduledAt: now + 6*60,
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
