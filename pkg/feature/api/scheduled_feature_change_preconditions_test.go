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
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	envclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	databasemock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// expectNoExperiments reports that no experiment exists for the flag.
func expectNoExperiments(service *FeatureService) {
	service.experimentClient.(*experimentclientmock.MockClient).EXPECT().ListExperiments(
		gomock.Any(), gomock.Any(),
	).Return(&experimentproto.ListExperimentsResponse{}, nil)
}

func TestExecuteScheduledFlagChangeAppliesUpdatePreconditions(t *testing.T) {
	t.Parallel()

	const (
		featureID  = "feature-id"
		envID      = "ns0"
		scheduleID = "sfc-id"
	)
	now := time.Now().Unix()

	flag := func() *domain.Feature {
		return &domain.Feature{Feature: &featureproto.Feature{
			Id:            featureID,
			Name:          "flag",
			Version:       1,
			Enabled:       false,
			VariationType: featureproto.Feature_BOOLEAN,
			Variations: []*featureproto.Variation{
				{Id: "var-true", Value: "true", Name: "true"},
				{Id: "var-false", Value: "false", Name: "false"},
			},
			OffVariation: "var-false",
			DefaultStrategy: &featureproto.Strategy{
				Type:          featureproto.Strategy_FIXED,
				FixedStrategy: &featureproto.FixedStrategy{Variation: "var-true"},
			},
		}}
	}
	// A due schedule that enables the flag.
	schedule := func(comment string) *domain.ScheduledFlagChange {
		return &domain.ScheduledFlagChange{ScheduledFlagChange: &featureproto.ScheduledFlagChange{
			Id:            scheduleID,
			FeatureId:     featureID,
			EnvironmentId: envID,
			ScheduledAt:   now - 60,
			Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
			Payload:       &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(true)},
			Comment:       comment,
		}}
	}
	experiment := func(startAt, stopAt, stoppedAt int64) *experimentproto.Experiment {
		return &experimentproto.Experiment{
			Id:        "experiment-id",
			FeatureId: featureID,
			StartAt:   startAt,
			StopAt:    stopAt,
			StoppedAt: stoppedAt,
		}
	}

	cases := []struct {
		name string
		// Comment of the schedule as returned by the plain read, and (when
		// set) by the locking read inside the transaction, as after a
		// concurrent edit of the schedule.
		comment           string
		commentUnderLock  *string
		experiments       []*experimentproto.Experiment
		experimentsErr    error
		environment       *envproto.EnvironmentV2
		environmentLookup bool
		wantErr           error
		wantCode          codes.Code
		// wantFailed means the rejected schedule is persisted as FAILED with a
		// reason so the batch executor does not retry it; otherwise a rejected
		// schedule must stay PENDING.
		wantFailed bool
	}{
		{
			name:        "rejects a schedule while an experiment on the flag is running",
			experiments: []*experimentproto.Experiment{experiment(now-3600, now+3600, math.MaxInt64)},
			wantErr:     statusWaitingOrRunningExperimentExists.Err(),
			wantFailed:  true,
		},
		{
			name:        "rejects a schedule while an experiment on the flag is waiting to start",
			experiments: []*experimentproto.Experiment{experiment(now+3600, now+7200, math.MaxInt64)},
			wantErr:     statusWaitingOrRunningExperimentExists.Err(),
			wantFailed:  true,
		},
		{
			name:              "rejects a schedule without a comment when the environment requires one",
			comment:           "",
			environment:       &envproto.EnvironmentV2{Id: envID, RequireComment: true},
			environmentLookup: true,
			wantErr:           statusCommentRequiredForUpdating.Err(),
			wantFailed:        true,
		},
		{
			name:              "rejects a schedule whose comment was cleared before the row was locked",
			comment:           "enable for the launch",
			commentUnderLock:  new(string),
			environment:       &envproto.EnvironmentV2{Id: envID, RequireComment: true},
			environmentLookup: true,
			wantErr:           statusCommentRequiredForUpdating.Err(),
			wantFailed:        true,
		},
		{
			name:           "leaves the schedule pending when the experiments cannot be listed",
			experimentsErr: status.Error(codes.Unavailable, "experiment service unavailable"),
			wantCode:       codes.Unavailable,
		},
		{
			name:              "applies a schedule once the experiment on the flag has finished",
			experiments:       []*experimentproto.Experiment{experiment(now-7200, now-3600, now-3600)},
			environment:       &envproto.EnvironmentV2{Id: envID},
			environmentLookup: true,
		},
		{
			name:              "applies a schedule with a comment when the environment requires one",
			comment:           "enable for the launch",
			environment:       &envproto.EnvironmentV2{Id: envID, RequireComment: true},
			environmentLookup: true,
		},
		{
			name:              "applies a schedule without a comment when the environment does not require one",
			comment:           "",
			environment:       &envproto.EnvironmentV2{Id: envID},
			environmentLookup: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			controller := gomock.NewController(t)
			defer controller.Finish()

			service := createFeatureServiceNew(controller)
			featureStorage := service.featureStorage.(*mock.MockFeatureStorage)
			scheduledStorage := service.scheduledFlagChangeStorage.(*mock.MockScheduledFlagChangeStorage)

			// Mark the transaction context so writes can be told apart from
			// the FAILED write, which must happen after the rollback.
			type txCtxKey struct{}
			service.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
				gomock.Any(), gomock.Any(),
			).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(context.WithValue(ctx, txCtxKey{}, true))
			})

			// The plain read resolves the flag, then the flag and the schedule
			// rows are locked and the schedule is re-read under its lock.
			scheduledStorage.EXPECT().GetScheduledFlagChange(gomock.Any(), scheduleID, envID).
				Return(schedule(tc.comment), nil)
			commentUnderLock := tc.comment
			if tc.commentUnderLock != nil {
				commentUnderLock = *tc.commentUnderLock
			}
			scheduledStorage.EXPECT().GetScheduledFlagChangeForUpdate(gomock.Any(), scheduleID, envID).
				Return(schedule(commentUnderLock), nil)
			experiments := service.experimentClient.(*experimentclientmock.MockClient).EXPECT().
				ListExperiments(gomock.Any(), &experimentproto.ListExperimentsRequest{
					FeatureId:     featureID,
					EnvironmentId: envID,
					PageSize:      listRequestSize,
				})
			if tc.experimentsErr != nil {
				experiments.Return(nil, tc.experimentsErr)
			} else {
				experiments.Return(&experimentproto.ListExperimentsResponse{Experiments: tc.experiments}, nil)
			}
			// The environment is only looked up once the experiment check passed.
			if tc.environmentLookup {
				service.environmentClient.(*envclientmock.MockClient).EXPECT().GetEnvironmentV2(
					gomock.Any(), &envproto.GetEnvironmentV2Request{Id: envID},
				).Return(&envproto.GetEnvironmentV2Response{Environment: tc.environment}, nil)
			}

			rejected := tc.wantErr != nil || tc.wantCode != codes.OK
			if rejected {
				// A rejected schedule changes nothing: the flag is locked but
				// neither read again for the update nor written.
				featureStorage.EXPECT().GetFeatureForUpdate(gomock.Any(), featureID, envID).Return(flag(), nil)
				featureStorage.EXPECT().ListFeatures(gomock.Any(), gomock.Any()).Times(0)
				featureStorage.EXPECT().UpdateFeature(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				if tc.wantFailed {
					scheduledStorage.EXPECT().UpdateScheduledFlagChange(gomock.Any(), gomock.Any()).
						DoAndReturn(func(updateCtx context.Context, updated *domain.ScheduledFlagChange) error {
							// Written after the rollback, not inside the transaction.
							assert.Nil(t, updateCtx.Value(txCtxKey{}))
							assert.Equal(
								t,
								featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_FAILED,
								updated.Status,
							)
							assert.NotEmpty(t, updated.FailureReason)
							return nil
						})
				} else {
					scheduledStorage.EXPECT().UpdateScheduledFlagChange(gomock.Any(), gomock.Any()).Times(0)
				}
			} else {
				// Once by the executor, once by updateFeatureWithinTransaction.
				featureStorage.EXPECT().GetFeatureForUpdate(gomock.Any(), featureID, envID).
					DoAndReturn(func(context.Context, string, string) (*domain.Feature, error) {
						return flag(), nil
					}).Times(2)
				featureStorage.EXPECT().ListFeatures(gomock.Any(), gomock.Any()).
					Return([]*featureproto.Feature{flag().Feature}, 0, int64(0), nil)
				featureStorage.EXPECT().UpdateFeature(gomock.Any(), gomock.Any(), envID).
					DoAndReturn(func(_ context.Context, updated *domain.Feature, _ string) error {
						assert.True(t, updated.Enabled)
						return nil
					})
				scheduledStorage.EXPECT().UpdateScheduledFlagChange(gomock.Any(), gomock.Any()).
					DoAndReturn(func(updateCtx context.Context, updated *domain.ScheduledFlagChange) error {
						assert.NotNil(t, updateCtx.Value(txCtxKey{}))
						assert.Equal(
							t,
							featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_EXECUTED,
							updated.Status,
						)
						return nil
					})
				// The flag update is recorded with the schedule's own comment.
				service.domainPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, messages []publisher.Message) map[string]error {
						require.Len(t, messages, 1)
						event := messages[0].(*eventproto.Event)
						assert.Equal(t, eventproto.Event_FEATURE_UPDATED, event.Type)
						assert.Equal(t, "Applied from scheduled change: "+tc.comment, event.Options.Comment)
						return nil
					})
			}

			resp, err := service.ExecuteScheduledFlagChange(
				createContextWithToken(),
				&featureproto.ExecuteScheduledFlagChangeRequest{Id: scheduleID, EnvironmentId: envID},
			)
			switch {
			case tc.wantErr != nil:
				assert.Equal(t, tc.wantErr, err)
				assert.Nil(t, resp)
			case tc.wantCode != codes.OK:
				assert.Equal(t, tc.wantCode, status.Code(err))
				assert.Nil(t, resp)
			default:
				require.NoError(t, err)
				assert.Equal(
					t,
					featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_EXECUTED,
					resp.ScheduledFlagChange.Status,
				)
			}
		})
	}
}
