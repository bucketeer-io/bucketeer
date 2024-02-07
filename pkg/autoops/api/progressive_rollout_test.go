// Copyright 2024 The Bucketeer Authors.
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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	v2as "github.com/bucketeer-io/bucketeer/pkg/autoops/storage/v2"
	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestCreateProgressiveRolloutMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	ti := time.Date(2020, 12, 15, 0, 0, 0, 0, time.UTC)
	invalidSpanSchedules := []*autoopsproto.ProgressiveRolloutSchedule{
		{
			ScheduleId: "sid-1",
			ExecuteAt:  ti.Unix(),
			Weight:     20,
		},
		{
			ScheduleId: "sid-1",
			ExecuteAt:  ti.Add(time.Minute * 3).Unix(),
			Weight:     40,
		},
	}

	validSchedules := []*autoopsproto.ProgressiveRolloutSchedule{
		{
			ScheduleId: "sid-1",
			ExecuteAt:  ti.Unix(),
			Weight:     20,
		},
		{
			ScheduleId: "sid-1",
			ExecuteAt:  ti.AddDate(0, 0, 3).Unix(),
			Weight:     40,
		},
	}

	executedAtRequiredSchedules := []*autoopsproto.ProgressiveRolloutSchedule{
		{
			ScheduleId: "sid-1",
			ExecuteAt:  0,
		},
	}

	invalidWeightSchedules := []*autoopsproto.ProgressiveRolloutSchedule{
		{
			ScheduleId: "sid-1",
			ExecuteAt:  ti.Unix(),
			Weight:     -1,
		},
	}

	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		req         *autoopsproto.CreateProgressiveRolloutRequest
		expectedErr error
	}{
		{
			desc: "err: ErrNoCommand",
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: nil,
			},
			expectedErr: createError(statusProgressiveRolloutNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "err: ErrFeatureIDRequired",
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{},
			},
			expectedErr: createError(statusProgressiveRolloutFeatureIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id")),
		},
		{
			desc: "err: Internal",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
				},
			},
			expectedErr: createError(statusProgressiveRolloutInternal, localizer.MustLocalizeWithTemplate(locale.InternalServerError)),
		},
		{
			desc: "err: InvalidVariationSize",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
						{
							Id: "vid-3",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
				},
			},
			expectedErr: createError(statusProgressiveRolloutInvalidVariationSize, localizer.MustLocalizeWithTemplate(locale.AutoOpsInvalidVariationSize)),
		},
		{
			desc: "err: ErrClauseRequired",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
				},
			},
			expectedErr: createError(statusProgressiveRolloutClauseRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause")),
		},
		{
			desc: "err: IncorrecctProgressiveRolloutClause",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId:                                "fid",
					ProgressiveRolloutManualScheduleClause:   &autoopsproto.ProgressiveRolloutManualScheduleClause{},
					ProgressiveRolloutTemplateScheduleClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{},
				},
			},
			expectedErr: createError(statusIncorrectProgressiveRolloutClause, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "clause")),
		},
		{
			desc: "err: manual ErrVariationIdRequired",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId:                              "fid",
					ProgressiveRolloutManualScheduleClause: &autoopsproto.ProgressiveRolloutManualScheduleClause{},
				},
			},
			expectedErr: createError(statusProgressiveRolloutClauseVariationIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			desc: "err: manual ErrInvalidVariationId",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutManualScheduleClause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
						VariationId: "invalid",
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutClauseInvalidVariationID, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_id")),
		},
		{
			desc: "err: template ErrVariationIdRequired",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId:                                "fid",
					ProgressiveRolloutTemplateScheduleClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{},
				},
			},
			expectedErr: createError(statusProgressiveRolloutClauseVariationIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			desc: "err: manual ErrSchedulesRequired",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutManualScheduleClause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
						VariationId: "vid-1",
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutClauseSchedulesRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "schedule")),
		},
		{
			desc: "err: template ErrSchedulesRequired",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutTemplateScheduleClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
						VariationId: "vid-1",
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutClauseSchedulesRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "schedule")),
		},
		{
			desc: "err: template ErrInvalidIncrements",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutTemplateScheduleClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
						VariationId: "vid-1",
						Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
							{
								ScheduleId: "sid-1",
								ExecuteAt:  time.Now().AddDate(3, 0, 0).Unix(),
								Weight:     60,
							},
						},
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutClauseInvalidIncrements, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "increments")),
		},
		{
			desc: "err: template ErrUnknownInterval",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutTemplateScheduleClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
						VariationId: "vid-1",
						Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
							{
								ScheduleId: "sid-1",
								ExecuteAt:  time.Now().AddDate(3, 0, 0).Unix(),
								Weight:     60,
							},
						},
						Increments: 30,
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutClauseUnknownInterval, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "interval")),
		},
		{
			desc: "err: manual ErrExecutedAtRequired",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutTemplateScheduleClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
						VariationId: "vid-1",
						Schedules:   executedAtRequiredSchedules,
						Interval:    autoopsproto.ProgressiveRolloutTemplateScheduleClause_DAILY,
						Increments:  2,
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutScheduleExecutedAtRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "execute_at")),
		},
		{
			desc: "err: manual ErrExecutedAtRequired",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutManualScheduleClause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
						VariationId: "vid-1",
						Schedules:   executedAtRequiredSchedules,
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutScheduleExecutedAtRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "execute_at")),
		},
		{
			desc: "err: manual ErrExecutedAtRequired",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutTemplateScheduleClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
						VariationId: "vid-1",
						Schedules:   invalidWeightSchedules,
						Interval:    autoopsproto.ProgressiveRolloutTemplateScheduleClause_DAILY,
						Increments:  2,
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutScheduleInvalidWeight, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "weight")),
		},
		{
			desc: "err: manual ErrExecutedAtRequired",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutManualScheduleClause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
						VariationId: "vid-1",
						Schedules:   invalidWeightSchedules,
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutScheduleInvalidWeight, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "weight")),
		},
		{
			desc: "err: template ErrInvalidScheduleSpans",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutTemplateScheduleClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
						VariationId: "vid-1",
						Schedules:   invalidSpanSchedules,
						Interval:    autoopsproto.ProgressiveRolloutTemplateScheduleClause_DAILY,
						Increments:  2,
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutInvalidScheduleSpans, localizer.MustLocalize(locale.AutoOpsInvalidScheduleSpans)),
		},
		{
			desc: "err: manual ErrInvalidScheduleSpans",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutManualScheduleClause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
						VariationId: "vid-1",
						Schedules:   invalidSpanSchedules,
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutInvalidScheduleSpans, localizer.MustLocalize(locale.AutoOpsInvalidScheduleSpans)),
		},
		{
			desc: "err: begin transaction error",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
				aos.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, errors.New("error"))
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutTemplateScheduleClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
						VariationId: "vid-1",
						Schedules:   validSchedules,
						Interval:    autoopsproto.ProgressiveRolloutTemplateScheduleClause_DAILY,
						Increments:  2,
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "err: transaction error",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
				aos.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				aos.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutTemplateScheduleClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
						VariationId: "vid-1",
						Schedules:   validSchedules,
						Interval:    autoopsproto.ProgressiveRolloutTemplateScheduleClause_DAILY,
						Increments:  2,
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "err: ErrProgressiveRolloutAlreadyExists",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
				aos.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				aos.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2as.ErrProgressiveRolloutAlreadyExists)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutTemplateScheduleClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
						VariationId: "vid-1",
						Schedules:   validSchedules,
						Interval:    autoopsproto.ProgressiveRolloutTemplateScheduleClause_DAILY,
						Increments:  2,
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutAlreadyExists, localizer.MustLocalize(locale.AlreadyExistsError)),
		},
		{
			desc: "err: AutoOpsHasDatetime",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
				aos.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				aos.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errProgressiveRolloutAutoOpsHasDatetime)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutTemplateScheduleClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
						VariationId: "vid-1",
						Schedules:   validSchedules,
						Interval:    autoopsproto.ProgressiveRolloutTemplateScheduleClause_DAILY,
						Increments:  2,
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutAutoOpsHasDatetime, localizer.MustLocalize(locale.AutoOpsHasDatetime)),
		},
		{
			desc: "err: AutoOpsHasWebhook",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
				aos.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				aos.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errProgressiveRolloutAutoOpsHasWebhook)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutTemplateScheduleClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
						VariationId: "vid-1",
						Schedules:   validSchedules,
						Interval:    autoopsproto.ProgressiveRolloutTemplateScheduleClause_DAILY,
						Increments:  2,
					},
				},
			},
			expectedErr: createError(statusProgressiveRolloutAutoOpsHasWebhook, localizer.MustLocalize(locale.AutoOpsHasWebhook)),
		},
		{
			desc: "success",
			setup: func(aos *AutoOpsService) {
				aos.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetFeatureResponse{Feature: &featureproto.Feature{
					Variations: []*featureproto.Variation{
						{
							Id: "vid-1",
						},
						{
							Id: "vid-2",
						},
					},
					Enabled: true,
				}}, nil)
				aos.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				aos.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.CreateProgressiveRolloutRequest{
				Command: &autoopsproto.CreateProgressiveRolloutCommand{
					FeatureId: "fid",
					ProgressiveRolloutTemplateScheduleClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
						VariationId: "vid-1",
						Schedules:   validSchedules,
						Interval:    autoopsproto.ProgressiveRolloutTemplateScheduleClause_DAILY,
						Increments:  2,
					},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.CreateProgressiveRollout(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetProgressiveRolloutMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleUnassigned(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		req         *autoopsproto.GetProgressiveRolloutRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &autoopsproto.GetProgressiveRolloutRequest{EnvironmentNamespace: "ns0"},
			expectedErr: createError(statusProgressiveRolloutIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrNotFound",
			setup: func(s *AutoOpsService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req:         &autoopsproto.GetProgressiveRolloutRequest{Id: "wrongid", EnvironmentNamespace: "ns0"},
			expectedErr: createError(statusProgressiveRolloutNotFound, localizer.MustLocalizeWithTemplate(locale.NotFoundError, locale.ProgressiveRollout)),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req:         &autoopsproto.GetProgressiveRolloutRequest{Id: "aid1", EnvironmentNamespace: "ns0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.GetProgressiveRollout(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestStopProgressiveRolloutMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		req         *autoopsproto.StopProgressiveRolloutRequest
		expectedErr error
	}{
		{
			desc:        "err: id is required",
			req:         &autoopsproto.StopProgressiveRolloutRequest{},
			expectedErr: createError(statusProgressiveRolloutIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:        "err: command is reuired",
			req:         &autoopsproto.StopProgressiveRolloutRequest{Id: "id", EnvironmentNamespace: "ns"},
			expectedErr: createError(statusProgressiveRolloutNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "err: failed to begin transaction",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, errors.New("error"))
			},
			req: &autoopsproto.StopProgressiveRolloutRequest{
				Id:                   "id",
				EnvironmentNamespace: "ns",
				Command: &autoopsproto.StopProgressiveRolloutCommand{
					StoppedBy: autoopsproto.ProgressiveRollout_USER,
				},
			},
			expectedErr: createError(statusProgressiveRolloutInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "err: internal error during transaction",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req: &autoopsproto.StopProgressiveRolloutRequest{
				Id:                   "id",
				EnvironmentNamespace: "ns",
				Command: &autoopsproto.StopProgressiveRolloutCommand{
					StoppedBy: autoopsproto.ProgressiveRollout_USER,
				},
			},
			expectedErr: createError(statusProgressiveRolloutInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "err: not found",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2as.ErrProgressiveRolloutNotFound)
			},
			req: &autoopsproto.StopProgressiveRolloutRequest{
				Id:                   "id",
				EnvironmentNamespace: "ns",
				Command: &autoopsproto.StopProgressiveRolloutCommand{
					StoppedBy: autoopsproto.ProgressiveRollout_USER,
				},
			},
			expectedErr: createError(statusProgressiveRolloutNotFound, localizer.MustLocalizeWithTemplate(locale.NotFoundError, locale.ProgressiveRollout)),
		},
		{
			desc: "err: unexpected affected rows",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2as.ErrProgressiveRolloutUnexpectedAffectedRows)
			},
			req: &autoopsproto.StopProgressiveRolloutRequest{
				Id:                   "id",
				EnvironmentNamespace: "ns",
				Command: &autoopsproto.StopProgressiveRolloutCommand{
					StoppedBy: autoopsproto.ProgressiveRollout_USER,
				},
			},
			expectedErr: createError(statusProgressiveRolloutNotFound, localizer.MustLocalizeWithTemplate(locale.NotFoundError, locale.ProgressiveRollout)),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.StopProgressiveRolloutRequest{
				Id:                   "id",
				EnvironmentNamespace: "ns",
				Command: &autoopsproto.StopProgressiveRolloutCommand{
					StoppedBy: autoopsproto.ProgressiveRollout_USER,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.StopProgressiveRollout(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteProgressiveRolloutMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		req         *autoopsproto.DeleteProgressiveRolloutRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &autoopsproto.DeleteProgressiveRolloutRequest{},
			expectedErr: createError(statusProgressiveRolloutIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: failed to begin transaction",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, errors.New("error"))
			},
			req:         &autoopsproto.DeleteProgressiveRolloutRequest{Id: "wrongid", EnvironmentNamespace: "ns0"},
			expectedErr: createError(statusProgressiveRolloutInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "err: internal error during transaction",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req:         &autoopsproto.DeleteProgressiveRolloutRequest{Id: "wrongid", EnvironmentNamespace: "ns0"},
			expectedErr: createError(statusProgressiveRolloutInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "err: internal error during transaction",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			req:         &autoopsproto.DeleteProgressiveRolloutRequest{Id: "wrongid", EnvironmentNamespace: "ns0"},
			expectedErr: createError(statusProgressiveRolloutInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "err: ErrProgressiveRolloutNotFound",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2as.ErrProgressiveRolloutNotFound)
			},
			req:         &autoopsproto.DeleteProgressiveRolloutRequest{Id: "wrongid", EnvironmentNamespace: "ns0"},
			expectedErr: createError(statusProgressiveRolloutNotFound, localizer.MustLocalizeWithTemplate(locale.NotFoundError, locale.ProgressiveRollout)),
		},
		{
			desc: "err: ErrProgressiveRolloutUnexpectedAffectedRows",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2as.ErrProgressiveRolloutUnexpectedAffectedRows)
			},
			req:         &autoopsproto.DeleteProgressiveRolloutRequest{Id: "wrongid", EnvironmentNamespace: "ns0"},
			expectedErr: createError(statusProgressiveRolloutNotFound, localizer.MustLocalizeWithTemplate(locale.NotFoundError, locale.ProgressiveRollout)),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req:         &autoopsproto.DeleteProgressiveRolloutRequest{Id: "aid1", EnvironmentNamespace: "ns0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.DeleteProgressiveRollout(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListProgressiveRolloutsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleUnassigned(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc                 string
		setup                func(*AutoOpsService)
		orderBy              autoopsproto.ListProgressiveRolloutsRequest_OrderBy
		environmentNamespace string
		expected             error
	}{
		{
			desc:                 "err: InvalidOrderBy",
			setup:                nil,
			orderBy:              autoopsproto.ListProgressiveRolloutsRequest_OrderBy(999),
			environmentNamespace: "ns0",
			expected:             createError(statusProgressiveRolloutInvalidOrderBy, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "order_by")),
		},
		{
			desc: "err: interal error",
			setup: func(s *AutoOpsService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			orderBy:              autoopsproto.ListProgressiveRolloutsRequest_DEFAULT,
			environmentNamespace: "ns0",
			expected:             createError(statusProgressiveRolloutInternal, localizer.MustLocalizeWithTemplate(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			orderBy:              autoopsproto.ListProgressiveRolloutsRequest_DEFAULT,
			environmentNamespace: "ns0",
			expected:             nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createAutoOpsService(mockController, nil)
			if p.setup != nil {
				p.setup(service)
			}
			req := &autoopsproto.ListProgressiveRolloutsRequest{
				OrderBy:              p.orderBy,
				EnvironmentNamespace: "ns0",
			}
			_, err := service.ListProgressiveRollouts(ctx, req)
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestExecuteProgressiveRolloutMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		req         *autoopsproto.ExecuteProgressiveRolloutRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &autoopsproto.ExecuteProgressiveRolloutRequest{},
			expectedErr: createError(statusProgressiveRolloutIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrNoCommand",
			req: &autoopsproto.ExecuteProgressiveRolloutRequest{
				Id: "aid",
			},
			expectedErr: createError(statusProgressiveRolloutNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "err: ErrNoCommand",
			req: &autoopsproto.ExecuteProgressiveRolloutRequest{
				Id: "aid",
				ChangeProgressiveRolloutTriggeredAtCommand: &autoopsproto.ChangeProgressiveRolloutScheduleTriggeredAtCommand{},
			},
			expectedErr: createError(statusProgressiveRolloutScheduleIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "schedule_id")),
		},
		{
			desc: "err: begin transaction error",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, errors.New("error"))
			},
			req: &autoopsproto.ExecuteProgressiveRolloutRequest{
				Id:                   "aid1",
				EnvironmentNamespace: "ns0",
				ChangeProgressiveRolloutTriggeredAtCommand: &autoopsproto.ChangeProgressiveRolloutScheduleTriggeredAtCommand{
					ScheduleId: "sid1",
				},
			},
			expectedErr: createError(statusProgressiveRolloutInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.ExecuteProgressiveRolloutRequest{
				Id:                   "aid1",
				EnvironmentNamespace: "ns0",
				ChangeProgressiveRolloutTriggeredAtCommand: &autoopsproto.ChangeProgressiveRolloutScheduleTriggeredAtCommand{
					ScheduleId: "sid1",
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.ExecuteProgressiveRollout(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
