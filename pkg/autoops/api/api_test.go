// Copyright 2025 The Bucketeer Authors.
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	authclientmock "github.com/bucketeer-io/bucketeer/pkg/auth/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	v2ao "github.com/bucketeer-io/bucketeer/pkg/autoops/storage/v2"
	mockAutoOpsStorage "github.com/bucketeer-io/bucketeer/pkg/autoops/storage/v2/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/pkg/experiment/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	mockOpsCountStorage "github.com/bucketeer-io/bucketeer/pkg/opsevent/storage/v2/mock"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

func TestNewAutoOpsService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	mysqlClientMock := mysqlmock.NewMockClient(mockController)
	featureClientMock := featureclientmock.NewMockClient(mockController)
	experimentClientMock := experimentclientmock.NewMockClient(mockController)
	accountClientMock := accountclientmock.NewMockClient(mockController)
	authClientMock := authclientmock.NewMockClient(mockController)
	p := publishermock.NewMockPublisher(mockController)
	logger := zap.NewNop()
	s := NewAutoOpsService(
		mysqlClientMock,
		featureClientMock,
		experimentClientMock,
		accountClientMock,
		authClientMock,
		p,
		WithLogger(logger),
	)
	assert.IsType(t, &AutoOpsService{}, s)
}

func TestCreateAutoOpsRuleMySQL(t *testing.T) {
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
		req         *autoopsproto.CreateAutoOpsRuleRequest
		expectedErr error
	}{
		{
			desc: "err: ErrFeatureIDRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{},
			},
			expectedErr: createError(statusFeatureIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id")),
		},
		{
			desc: "err: ErrClauseRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_SCHEDULE,
				},
			},
			expectedErr: createError(statusClauseRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause")),
		},
		{
			desc: "err: ErrIncompatibleOpsType",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_TYPE_UNKNOWN,
					OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
						{
							VariationId:     "",
							GoalId:          "gid",
							MinCount:        10,
							ThreadsholdRate: 0.5,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						},
					},
				},
			},
			expectedErr: createError(statusIncompatibleOpsType, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "ops_type")),
		},
		{
			desc: "err: ErrOpsEventRateClauseVariationIDRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_EVENT_RATE,
					OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
						{
							VariationId:     "",
							GoalId:          "gid1",
							MinCount:        10,
							ThreadsholdRate: 0.5,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						},
					},
				},
			},
			expectedErr: createError(statusOpsEventRateClauseVariationIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			desc: "err: ErrOpsEventRateClauseGoalIDRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_EVENT_RATE,
					OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
						{
							VariationId:     "vid",
							GoalId:          "",
							MinCount:        10,
							ThreadsholdRate: 0.5,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						},
					},
				},
			},
			expectedErr: createError(statusOpsEventRateClauseGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			desc: "err: ErrOpsEventRateClauseMinCountRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_EVENT_RATE,
					OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
						{
							VariationId:     "vid",
							GoalId:          "gid",
							MinCount:        0,
							ThreadsholdRate: 0.5,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						},
					},
				},
			},
			expectedErr: createError(statusOpsEventRateClauseMinCountRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "min_count")),
		},
		{
			desc: "err: ErrOpsEventRateClauseInvalidThredshold: less",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_EVENT_RATE,
					OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
						{
							VariationId:     "vid",
							GoalId:          "gid",
							MinCount:        10,
							ThreadsholdRate: -0.1,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						},
					},
				},
			},
			expectedErr: createError(statusOpsEventRateClauseInvalidThredshold, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "threshold")),
		},
		{
			desc: "err: ErrOpsEventRateClauseInvalidThredshold: greater",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_EVENT_RATE,
					OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
						{
							VariationId:     "vid",
							GoalId:          "gid",
							MinCount:        10,
							ThreadsholdRate: 1.1,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						},
					},
				},
			},
			expectedErr: createError(statusOpsEventRateClauseInvalidThredshold, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "threshold")),
		},
		{
			desc: "err: ErrDatetimeClauseInvalidTime",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_SCHEDULE,
					DatetimeClauses: []*autoopsproto.DatetimeClause{
						{Time: 0},
					},
				},
			},
			expectedErr: createError(statusDatetimeClauseInvalidTime, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time")),
		},
		{
			desc: "err: ErrDatetimeClauseDuplicateTime",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_SCHEDULE,
					DatetimeClauses: []*autoopsproto.DatetimeClause{
						{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
						{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
					},
				},
			},
			expectedErr: createError(statusDatetimeClauseDuplicateTime, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time")),
		},
		{
			desc: "err: ErrDatetimeClauseMustSpecified",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId:       "fid",
					OpsType:         autoopsproto.OpsType_SCHEDULE,
					DatetimeClauses: nil,
					OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
						{
							VariationId:     "vid",
							GoalId:          "gid",
							MinCount:        10,
							ThreadsholdRate: 0.5,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
							ActionType:      autoopsproto.ActionType_DISABLE,
						},
					},
				},
			},
			expectedErr: createError(statusClauseRequiredForDateTime, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "clause")),
		},
		{
			desc: "err: ErrDatetimeClauseMustSpecified",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_SCHEDULE,
					DatetimeClauses: []*autoopsproto.DatetimeClause{
						{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
					},
					OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
						{
							VariationId:     "vid",
							GoalId:          "gid",
							MinCount:        10,
							ThreadsholdRate: 0.5,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
							ActionType:      autoopsproto.ActionType_DISABLE,
						},
					},
				},
			},
			expectedErr: createError(statusIncompatibleOpsType, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "ops_type")),
		},
		{
			desc: "err: ErrOpsEventRateClauseMustSpecified",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_EVENT_RATE,
					DatetimeClauses: []*autoopsproto.DatetimeClause{
						{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
					},
					OpsEventRateClauses: nil,
				},
			},
			expectedErr: createError(statusClauseRequiredForEventDate, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "clause")),
		},
		{
			desc: "err: ErrDatetimeClauseMustNotBeSpecified",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_EVENT_RATE,
					DatetimeClauses: []*autoopsproto.DatetimeClause{
						{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
					},
					OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
						{
							VariationId:     "vid",
							GoalId:          "gid",
							MinCount:        10,
							ThreadsholdRate: 0.5,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
							ActionType:      autoopsproto.ActionType_DISABLE,
						},
					},
				},
			},
			expectedErr: createError(statusIncompatibleOpsType, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "ops_type")),
		},
		{
			desc: "err: internal error",
			setup: func(s *AutoOpsService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().GetGoal(
					gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_EVENT_RATE,
					OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
						{
							VariationId:     "vid",
							GoalId:          "gid",
							MinCount:        10,
							ThreadsholdRate: 0.5,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
							ActionType:      autoopsproto.ActionType_DISABLE,
						},
					},
				},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success event rate",
			setup: func(s *AutoOpsService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().GetGoal(
					gomock.Any(), gomock.Any(),
				).Return(&experimentproto.GetGoalResponse{}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().CreateAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_EVENT_RATE,
					OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
						{
							VariationId:     "vid",
							GoalId:          "gid",
							MinCount:        10,
							ThreadsholdRate: 0.5,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
							ActionType:      autoopsproto.ActionType_DISABLE,
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success schedule",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().CreateAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_SCHEDULE,
					DatetimeClauses: []*autoopsproto.DatetimeClause{
						{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
					},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.CreateAutoOpsRule(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCreateAutoOpsRuleMySQLNoCommand(t *testing.T) {
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
		req         *autoopsproto.CreateAutoOpsRuleRequest
		expectedErr error
	}{
		{
			desc: "err: ErrFeatureIDRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				EnvironmentId: "env-id",
			},
			expectedErr: createError(statusFeatureIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id")),
		},
		{
			desc: "err: ErrClauseRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_SCHEDULE,
			},
			expectedErr: createError(statusClauseRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause")),
		},
		{
			desc: "err: ErrIncompatibleOpsType",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_TYPE_UNKNOWN,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					},
				},
			},
			expectedErr: createError(statusIncompatibleOpsType, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "ops_type")),
		},
		{
			desc: "err: ErrOpsEventRateClauseVariationIDRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "",
						GoalId:          "gid1",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					},
				},
			},
			expectedErr: createError(statusOpsEventRateClauseVariationIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			desc: "err: ErrOpsEventRateClauseGoalIDRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					},
				},
			},
			expectedErr: createError(statusOpsEventRateClauseGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			desc: "err: ErrOpsEventRateClauseMinCountRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        0,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					},
				},
			},
			expectedErr: createError(statusOpsEventRateClauseMinCountRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "min_count")),
		},
		{
			desc: "err: ErrOpsEventRateClauseInvalidThredshold: less",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: -0.1,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					},
				},
			},
			expectedErr: createError(statusOpsEventRateClauseInvalidThredshold, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "threshold")),
		},
		{
			desc: "err: ErrOpsEventRateClauseInvalidThredshold: greater",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: 1.1,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					},
				},
			},
			expectedErr: createError(statusOpsEventRateClauseInvalidThredshold, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "threshold")),
		},
		{
			desc: "err: ErrDatetimeClauseInvalidTime",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_SCHEDULE,
				DatetimeClauses: []*autoopsproto.DatetimeClause{
					{Time: 0},
				},
			},
			expectedErr: createError(statusDatetimeClauseInvalidTime, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time")),
		},
		{
			desc: "err: ErrDatetimeClauseDuplicateTime",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_SCHEDULE,
				DatetimeClauses: []*autoopsproto.DatetimeClause{
					{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
					{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
				},
			},
			expectedErr: createError(statusDatetimeClauseDuplicateTime, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time")),
		},
		{
			desc: "err: ErrDatetimeClauseMustSpecified",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId:       "fid",
				OpsType:         autoopsproto.OpsType_SCHEDULE,
				DatetimeClauses: nil,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						ActionType:      autoopsproto.ActionType_DISABLE,
					},
				},
			},
			expectedErr: createError(statusClauseRequiredForDateTime, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "clause")),
		},
		{
			desc: "err: ErrDatetimeClauseMustSpecified",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_SCHEDULE,
				DatetimeClauses: []*autoopsproto.DatetimeClause{
					{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
				},
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						ActionType:      autoopsproto.ActionType_DISABLE,
					},
				},
			},
			expectedErr: createError(statusIncompatibleOpsType, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "ops_type")),
		},
		{
			desc: "err: ErrOpsEventRateClauseMustSpecified",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				DatetimeClauses: []*autoopsproto.DatetimeClause{
					{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
				},
				OpsEventRateClauses: nil,
			},
			expectedErr: createError(statusClauseRequiredForEventDate, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "clause")),
		},
		{
			desc: "err: ErrDatetimeClauseMustNotBeSpecified",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				DatetimeClauses: []*autoopsproto.DatetimeClause{
					{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
				},
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						ActionType:      autoopsproto.ActionType_DISABLE,
					},
				},
			},
			expectedErr: createError(statusIncompatibleOpsType, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "ops_type")),
		},
		{
			desc: "err: internal error",
			setup: func(s *AutoOpsService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().GetGoal(
					gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						ActionType:      autoopsproto.ActionType_DISABLE,
					},
				},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success event rate",
			setup: func(s *AutoOpsService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().GetGoal(
					gomock.Any(), gomock.Any(),
				).Return(&experimentproto.GetGoalResponse{}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().CreateAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						ActionType:      autoopsproto.ActionType_DISABLE,
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success schedule",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().CreateAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_SCHEDULE,
				DatetimeClauses: []*autoopsproto.DatetimeClause{
					{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.CreateAutoOpsRule(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAutoOpsRuleMySQL(t *testing.T) {
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
		req         *autoopsproto.UpdateAutoOpsRuleRequest
		expected    *autoopsproto.UpdateAutoOpsRuleResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &autoopsproto.UpdateAutoOpsRuleRequest{},
			expected:    nil,
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrNoCommand",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
			},
			expected:    nil,
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "err: ErrOpsEventRateClauseRequired",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id:                            "aid1",
				AddOpsEventRateClauseCommands: []*autoopsproto.AddOpsEventRateClauseCommand{{}},
			},
			expected:    nil,
			expectedErr: createError(statusOpsEventRateClauseRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "ops_event_rate_clause")),
		},
		{
			desc: "err: DeleteClauseCommand: ErrClauseIdRequired",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id:                   "aid1",
				DeleteClauseCommands: []*autoopsproto.DeleteClauseCommand{{}},
			},
			expected:    nil,
			expectedErr: createError(statusClauseIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause_id")),
		},
		{
			desc: "err: ChangeOpsEventRateClauseCommand: ErrClauseIdRequired",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id:                               "aid1",
				ChangeOpsEventRateClauseCommands: []*autoopsproto.ChangeOpsEventRateClauseCommand{{}},
			},
			expected:    nil,
			expectedErr: createError(statusClauseIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause_id")),
		},
		{
			desc: "err: ChangeOpsEventRateClauseCommand: ErrOpsEventRateClauseRequired",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
				ChangeOpsEventRateClauseCommands: []*autoopsproto.ChangeOpsEventRateClauseCommand{{
					Id: "aid",
				}},
			},
			expected:    nil,
			expectedErr: createError(statusOpsEventRateClauseRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "ops_event_rate_clause")),
		},
		{
			desc: "err: ErrDatetimeClauseReqired",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id:                        "aid1",
				AddDatetimeClauseCommands: []*autoopsproto.AddDatetimeClauseCommand{{}},
			},
			expected:    nil,
			expectedErr: createError(statusDatetimeClauseRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "datetime_clause")),
		},
		{
			desc: "err: ChangeDatetimeClauseCommand: ErrDatetimeClauseInvalidTime",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
				ChangeDatetimeClauseCommands: []*autoopsproto.ChangeDatetimeClauseCommand{{
					Id:             "aid",
					DatetimeClause: &autoopsproto.DatetimeClause{Time: 0},
				}},
			},
			expected:    nil,
			expectedErr: createError(statusDatetimeClauseInvalidTime, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time")),
		},
		{
			desc: "err: ChangeDatetimeClauseCommand: ErrDatetimeClauseDuplicateTime",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
				ChangeDatetimeClauseCommands: []*autoopsproto.ChangeDatetimeClauseCommand{
					{
						Id:             "aid",
						DatetimeClause: &autoopsproto.DatetimeClause{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
					},
					{
						Id:             "aid2",
						DatetimeClause: &autoopsproto.DatetimeClause{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
					},
				},
			},
			expected:    nil,
			expectedErr: createError(statusDatetimeClauseDuplicateTime, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time")),
		},
		{
			desc: "err: AddDatetimeClauseCommand: ErrDatetimeClauseInvalidTime",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
				AddDatetimeClauseCommands: []*autoopsproto.AddDatetimeClauseCommand{
					{
						DatetimeClause: &autoopsproto.DatetimeClause{Time: 0, ActionType: autoopsproto.ActionType_DISABLE},
					},
				},
			},
			expected:    nil,
			expectedErr: createError(statusDatetimeClauseInvalidTime, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time")),
		},
		{
			desc: "err: AddDatetimeClauseCommand: ErrDatetimeClauseDuplicateTime",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
				AddDatetimeClauseCommands: []*autoopsproto.AddDatetimeClauseCommand{
					{
						DatetimeClause: &autoopsproto.DatetimeClause{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_DISABLE},
					},
					{
						DatetimeClause: &autoopsproto.DatetimeClause{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_DISABLE},
					},
				},
			},
			expected:    nil,
			expectedErr: createError(statusDatetimeClauseDuplicateTime, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time")),
		},
		{
			desc: "err: AddDatetimeClauseCommands: ErrDatetimeClauseDuplicateTime",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
				AddDatetimeClauseCommands: []*autoopsproto.AddDatetimeClauseCommand{
					{
						DatetimeClause: &autoopsproto.DatetimeClause{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_DISABLE},
					},
				},
				ChangeDatetimeClauseCommands: []*autoopsproto.ChangeDatetimeClauseCommand{
					{
						Id:             "aid",
						DatetimeClause: &autoopsproto.DatetimeClause{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_DISABLE},
					},
				},
			},
			expected:    nil,
			expectedErr: createError(statusDatetimeClauseDuplicateTime, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time")),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.All(),
				).Return(&domain.AutoOpsRule{
					AutoOpsRule: &autoopsproto.AutoOpsRule{
						Id: "aid1", OpsType: autoopsproto.OpsType_SCHEDULE, AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING, Deleted: false, Clauses: []*autoopsproto.Clause{
							{Id: "cid", ActionType: autoopsproto.ActionType_ENABLE, Clause: &anypb.Any{}},
							{Id: "cid2", ActionType: autoopsproto.ActionType_ENABLE, Clause: &anypb.Any{}},
						}},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil).AnyTimes()
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().UpdateAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id:            "aid1",
				EnvironmentId: "ns0",
				AddDatetimeClauseCommands: []*autoopsproto.AddDatetimeClauseCommand{{
					DatetimeClause: &autoopsproto.DatetimeClause{
						ActionType: autoopsproto.ActionType_ENABLE,
						Time:       time.Now().AddDate(0, 0, 1).Unix(),
					},
				}},
				DeleteClauseCommands: []*autoopsproto.DeleteClauseCommand{{
					Id: "cid",
				}},
				ChangeDatetimeClauseCommands: []*autoopsproto.ChangeDatetimeClauseCommand{
					{
						Id: "cid2",
						DatetimeClause: &autoopsproto.DatetimeClause{
							ActionType: autoopsproto.ActionType_DISABLE,
							Time:       time.Now().AddDate(0, 0, 2).Unix(),
						},
					},
				},
			},
			expected:    &autoopsproto.UpdateAutoOpsRuleResponse{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.UpdateAutoOpsRule(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestStopAutoOpsRuleMySQL(t *testing.T) {
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
		req         *autoopsproto.StopAutoOpsRuleRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &autoopsproto.StopAutoOpsRuleRequest{},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrNoCommand",
			req: &autoopsproto.StopAutoOpsRuleRequest{
				Id: "aid1",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AutoOpsRule{
					AutoOpsRule: &autoopsproto.AutoOpsRule{
						Id: "aid1", OpsType: autoopsproto.OpsType_SCHEDULE, AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING, Deleted: false, Clauses: []*autoopsproto.Clause{
							{Id: "cid", ActionType: autoopsproto.ActionType_ENABLE, Clause: &anypb.Any{}},
						}},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().UpdateAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.StopAutoOpsRuleRequest{
				Id:            "aid1",
				EnvironmentId: "ns0",
				Command:       &autoopsproto.StopAutoOpsRuleCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.StopAutoOpsRule(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteAutoOpsRuleMySQL(t *testing.T) {
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
		req         *autoopsproto.DeleteAutoOpsRuleRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &autoopsproto.DeleteAutoOpsRuleRequest{},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrNoCommand",
			req: &autoopsproto.DeleteAutoOpsRuleRequest{
				Id: "aid1",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AutoOpsRule{
					AutoOpsRule: &autoopsproto.AutoOpsRule{
						Id: "aid1", OpsType: autoopsproto.OpsType_SCHEDULE, AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING, Deleted: false, Clauses: []*autoopsproto.Clause{
							{Id: "cid", ActionType: autoopsproto.ActionType_ENABLE, Clause: &anypb.Any{}},
						}},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().UpdateAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.DeleteAutoOpsRuleRequest{
				Id:            "aid1",
				EnvironmentId: "ns0",
				Command:       &autoopsproto.DeleteAutoOpsRuleCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.DeleteAutoOpsRule(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAutoOpsRuleMySQL(t *testing.T) {
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
		service     *AutoOpsService
		setup       func(*AutoOpsService)
		req         *autoopsproto.GetAutoOpsRuleRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			service:     createAutoOpsService(mockController),
			req:         &autoopsproto.GetAutoOpsRuleRequest{EnvironmentId: "ns0"},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:    "err: ErrNotFound",
			service: createAutoOpsService(mockController),
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, v2ao.ErrAutoOpsRuleNotFound)
			},
			req:         &autoopsproto.GetAutoOpsRuleRequest{Id: "wrongid", EnvironmentId: "ns0"},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc:        "errPermissionDenied",
			service:     createServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			setup:       func(s *AutoOpsService) {},
			req:         &autoopsproto.GetAutoOpsRuleRequest{Id: "aid1", EnvironmentId: "ns0"},
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc:    "success",
			service: createAutoOpsService(mockController),
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AutoOpsRule{
					AutoOpsRule: &autoopsproto.AutoOpsRule{
						Id: "aid1", OpsType: autoopsproto.OpsType_SCHEDULE, AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING, Deleted: false, Clauses: []*autoopsproto.Clause{
							{Id: "cid", ActionType: autoopsproto.ActionType_ENABLE, Clause: &anypb.Any{}},
						}},
				}, nil)
			},
			req:         &autoopsproto.GetAutoOpsRuleRequest{Id: "aid1", EnvironmentId: "ns0"},
			expectedErr: nil,
		},
		{
			desc:    "success with view account",
			service: createServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AutoOpsRule{
					AutoOpsRule: &autoopsproto.AutoOpsRule{
						Id: "aid1", OpsType: autoopsproto.OpsType_SCHEDULE, AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING, Deleted: false, Clauses: []*autoopsproto.Clause{
							{Id: "cid", ActionType: autoopsproto.ActionType_ENABLE, Clause: &anypb.Any{}},
						}},
				}, nil)
			},
			req:         &autoopsproto.GetAutoOpsRuleRequest{Id: "aid1", EnvironmentId: "ns0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := p.service
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.GetAutoOpsRule(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAutoOpsRulesMySQL(t *testing.T) {
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
		service     *AutoOpsService
		setup       func(*AutoOpsService)
		req         *autoopsproto.ListAutoOpsRulesRequest
		expectedErr error
	}{
		{
			desc:    "success",
			service: createAutoOpsService(mockController),
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().ListAutoOpsRules(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*autoopsproto.AutoOpsRule{}, 0, nil)
			},
			req:         &autoopsproto.ListAutoOpsRulesRequest{EnvironmentId: "ns0", Cursor: ""},
			expectedErr: nil,
		},
		{
			desc:        "errPermissionDenied",
			service:     createServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			setup:       func(s *AutoOpsService) {},
			req:         &autoopsproto.ListAutoOpsRulesRequest{EnvironmentId: "ns0", Cursor: ""},
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc:    "success with viewer",
			service: createServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().ListAutoOpsRules(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*autoopsproto.AutoOpsRule{}, 0, nil)
			},
			req:         &autoopsproto.ListAutoOpsRulesRequest{EnvironmentId: "ns0", Cursor: ""},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := p.service
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.ListAutoOpsRules(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListOpsCountsMySQL(t *testing.T) {
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
		service     *AutoOpsService
		setup       func(*AutoOpsService)
		req         *autoopsproto.ListOpsCountsRequest
		expectedErr error
	}{
		{
			desc:    "success",
			service: createAutoOpsService(mockController),
			setup: func(s *AutoOpsService) {
				s.opsCountStorage.(*mockOpsCountStorage.MockOpsCountStorage).EXPECT().ListOpsCounts(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*autoopsproto.OpsCount{}, 0, nil)
			},
			req:         &autoopsproto.ListOpsCountsRequest{EnvironmentId: "ns0", Cursor: ""},
			expectedErr: nil,
		},
		{
			desc:        "errPermissionDenied",
			service:     createServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			setup:       func(s *AutoOpsService) {},
			req:         &autoopsproto.ListOpsCountsRequest{EnvironmentId: "ns0", Cursor: ""},
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc:    "success with view ",
			service: createServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *AutoOpsService) {
				s.opsCountStorage.(*mockOpsCountStorage.MockOpsCountStorage).EXPECT().ListOpsCounts(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*autoopsproto.OpsCount{}, 0, nil)
			},
			req:         &autoopsproto.ListOpsCountsRequest{EnvironmentId: "ns0", Cursor: ""},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := p.service
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.ListOpsCounts(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestExecuteAutoOpsRuleMySQL(t *testing.T) {
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
		req         *autoopsproto.ExecuteAutoOpsRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &autoopsproto.ExecuteAutoOpsRequest{},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrNoCommand",
			req: &autoopsproto.ExecuteAutoOpsRequest{
				Id: "aid",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "err: ErrNoExecuteAutoOpsRuleCommand_ClauseId",
			req: &autoopsproto.ExecuteAutoOpsRequest{
				Id:            "aid1",
				EnvironmentId: "ns0",
				ExecuteAutoOpsRuleCommand: &autoopsproto.ExecuteAutoOpsRuleCommand{
					ClauseId: "",
				},
			},
			expectedErr: createError(statusClauseRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause_id")),
		},
		{
			desc: "err: ErrNotFound",
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, v2ao.ErrAutoOpsRuleNotFound)
			},
			req: &autoopsproto.ExecuteAutoOpsRequest{
				Id:            "aid1",
				EnvironmentId: "ns0",
				ExecuteAutoOpsRuleCommand: &autoopsproto.ExecuteAutoOpsRuleCommand{
					ClauseId: "id",
				},
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AutoOpsRule{
					AutoOpsRule: &autoopsproto.AutoOpsRule{
						Id: "aid1", OpsType: autoopsproto.OpsType_SCHEDULE, AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING, Deleted: false, Clauses: []*autoopsproto.Clause{
							{Id: "testClauseId", ActionType: autoopsproto.ActionType_ENABLE, Clause: &anypb.Any{}},
						}},
				}, nil).AnyTimes()
			},
			req: &autoopsproto.ExecuteAutoOpsRequest{
				Id:            "aid1",
				EnvironmentId: "ns0",
				ExecuteAutoOpsRuleCommand: &autoopsproto.ExecuteAutoOpsRuleCommand{
					ClauseId: "testClauseId",
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.ExecuteAutoOps(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestExistGoal(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		goalID      string
		expected    bool
		expectedErr error
	}{
		{
			desc: "not found",
			setup: func(s *AutoOpsService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().GetGoal(gomock.Any(), gomock.Any()).Return(nil, storage.ErrKeyNotFound)
			},
			goalID:      "gid-0",
			expected:    false,
			expectedErr: nil,
		},
		{
			desc: "fails",
			setup: func(s *AutoOpsService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().GetGoal(gomock.Any(), gomock.Any()).Return(nil, errors.New("test"))
			},
			goalID:      "gid-0",
			expected:    false,
			expectedErr: errors.New("test"),
		},
		{
			desc: "exists",
			setup: func(s *AutoOpsService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().GetGoal(gomock.Any(), gomock.Any()).Return(&experimentproto.GetGoalResponse{}, nil)
			},
			goalID:      "gid-0",
			expected:    true,
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.existGoal(context.Background(), "ns-0", p.goalID)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func createAutoOpsService(c *gomock.Controller) *AutoOpsService {
	mysqlClientMock := mysqlmock.NewMockClient(c)
	featureClientMock := featureclientmock.NewMockClient(c)
	accountClientMock := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountV2ByEnvironmentIDResponse{
		Account: &accountproto.AccountV2{
			Email:            "email",
			OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "ns0",
					Role:          accountproto.AccountV2_Role_Environment_EDITOR,
				},
				{
					EnvironmentId: "",
					Role:          accountproto.AccountV2_Role_Environment_EDITOR,
				},
			},
		},
	}
	accountClientMock.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	experimentClientMock := experimentclientmock.NewMockClient(c)
	authClientMock := authclientmock.NewMockClient(c)
	p := publishermock.NewMockPublisher(c)
	logger := zap.NewNop()
	return &AutoOpsService{
		mysqlClient:      mysqlClientMock,
		autoOpsStorage:   mockAutoOpsStorage.NewMockAutoOpsRuleStorage(c),
		prStorage:        mockAutoOpsStorage.NewMockProgressiveRolloutStorage(c),
		opsCountStorage:  mockOpsCountStorage.NewMockOpsCountStorage(c),
		featureClient:    featureClientMock,
		experimentClient: experimentClientMock,
		accountClient:    accountClientMock,
		authClient:       authClientMock,
		publisher:        p,
		opts: &options{
			logger: zap.NewNop(),
		},
		logger: logger,
	}
}

func createServiceWithGetAccountByEnvironmentMock(c *gomock.Controller, ro accountproto.AccountV2_Role_Organization, re accountproto.AccountV2_Role_Environment) *AutoOpsService {
	mysqlClientMock := mysqlmock.NewMockClient(c)
	featureClientMock := featureclientmock.NewMockClient(c)
	accountClientMock := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountV2ByEnvironmentIDResponse{
		Account: &accountproto.AccountV2{
			Email:            "email",
			OrganizationRole: ro,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "ns0",
					Role:          re,
				},
			},
		},
	}
	accountClientMock.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	experimentClientMock := experimentclientmock.NewMockClient(c)
	authClientMock := authclientmock.NewMockClient(c)
	p := publishermock.NewMockPublisher(c)
	logger := zap.NewNop()
	return &AutoOpsService{
		mysqlClient:      mysqlClientMock,
		autoOpsStorage:   mockAutoOpsStorage.NewMockAutoOpsRuleStorage(c),
		prStorage:        mockAutoOpsStorage.NewMockProgressiveRolloutStorage(c),
		opsCountStorage:  mockOpsCountStorage.NewMockOpsCountStorage(c),
		featureClient:    featureClientMock,
		experimentClient: experimentClientMock,
		accountClient:    accountClientMock,
		authClient:       authClientMock,
		publisher:        p,
		opts: &options{
			logger: zap.NewNop(),
		},
		logger: logger,
	}
}

func createContextWithTokenRoleUnassigned(t *testing.T) context.Context {
	t.Helper()
	token := &token.AccessToken{
		Issuer:   "issuer",
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "email",
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}

func createContextWithTokenRoleOwner(t *testing.T) context.Context {
	t.Helper()
	token := &token.AccessToken{
		Issuer:        "issuer",
		Audience:      "audience",
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         "email",
		IsSystemAdmin: true,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}
