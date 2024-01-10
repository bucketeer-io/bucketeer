// Copyright 2023 The Bucketeer Authors.
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
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	authclientmock "github.com/bucketeer-io/bucketeer/pkg/auth/client/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/pkg/experiment/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
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

var testWebhookURL = func() *url.URL {
	u, err := url.Parse("https://bucketeer.io/hook")
	if err != nil {
		panic(err)
	}
	return u
}()

type dummyWebhookCryptoUtil struct{}

func (u *dummyWebhookCryptoUtil) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	return []byte(data), nil
}

func (u *dummyWebhookCryptoUtil) Decrypt(ctx context.Context, data []byte) ([]byte, error) {
	return []byte(data), nil
}

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
		testWebhookURL,
		&dummyWebhookCryptoUtil{},
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
			desc:        "err: ErrNoCommand",
			req:         &autoopsproto.CreateAutoOpsRuleRequest{},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
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
					OpsType:   autoopsproto.OpsType_ENABLE_FEATURE,
				},
			},
			expectedErr: createError(statusClauseRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause")),
		},
		{
			desc: "err: ErrIncompatibleOpsType",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_ENABLE_FEATURE,
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
					OpsType:   autoopsproto.OpsType_DISABLE_FEATURE,
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
					OpsType:   autoopsproto.OpsType_DISABLE_FEATURE,
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
					OpsType:   autoopsproto.OpsType_DISABLE_FEATURE,
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
					OpsType:   autoopsproto.OpsType_DISABLE_FEATURE,
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
					OpsType:   autoopsproto.OpsType_DISABLE_FEATURE,
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
					OpsType:   autoopsproto.OpsType_ENABLE_FEATURE,
					DatetimeClauses: []*autoopsproto.DatetimeClause{
						{Time: 0},
					},
				},
			},
			expectedErr: createError(statusDatetimeClauseInvalidTime, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time")),
		},
		{
			desc: "err: ErrWebhookClauseWebhookIDRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_DISABLE_FEATURE,
					WebhookClauses: []*autoopsproto.WebhookClause{
						{
							WebhookId: "",
							Conditions: []*autoopsproto.WebhookClause_Condition{
								{
									Filter:   ".foo.bar",
									Value:    "foobaz",
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
							},
						},
					},
				},
			},
			expectedErr: createError(statusWebhookClauseWebhookIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "webhook_id")),
		},
		{
			desc: "err: ErrWebhookClauseWebhookClauseConditionRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_DISABLE_FEATURE,
					WebhookClauses: []*autoopsproto.WebhookClause{
						{
							WebhookId:  "webhook-1",
							Conditions: []*autoopsproto.WebhookClause_Condition{},
						},
					},
				},
			},
			expectedErr: createError(statusWebhookClauseConditionRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "condition")),
		},
		{
			desc: "err: ErrWebhookClauseWebhookClauseConditionFilterRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_DISABLE_FEATURE,
					WebhookClauses: []*autoopsproto.WebhookClause{
						{
							WebhookId: "foo-id",
							Conditions: []*autoopsproto.WebhookClause_Condition{
								{
									Filter:   "",
									Value:    "foobaz",
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
							},
						},
					},
				},
			},
			expectedErr: createError(statusWebhookClauseConditionFilterRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "condition_filter")),
		},
		{
			desc: "err: internal error",
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
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_DISABLE_FEATURE,
					OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
						{
							VariationId:     "vid",
							GoalId:          "gid",
							MinCount:        10,
							ThreadsholdRate: 0.5,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						},
					},
					DatetimeClauses: []*autoopsproto.DatetimeClause{
						{Time: time.Now().AddDate(0, 0, 1).Unix()},
					},
					WebhookClauses: []*autoopsproto.WebhookClause{
						{
							WebhookId: "foo-id",
							Conditions: []*autoopsproto.WebhookClause_Condition{
								{
									Filter:   ".foo.bar",
									Value:    "foobaz",
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
							},
						},
					},
				},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().GetGoal(
					gomock.Any(), gomock.Any(),
				).Return(&experimentproto.GetGoalResponse{}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
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
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				Command: &autoopsproto.CreateAutoOpsRuleCommand{
					FeatureId: "fid",
					OpsType:   autoopsproto.OpsType_DISABLE_FEATURE,
					OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
						{
							VariationId:     "vid",
							GoalId:          "gid",
							MinCount:        10,
							ThreadsholdRate: 0.5,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						},
					},
					DatetimeClauses: []*autoopsproto.DatetimeClause{
						{Time: time.Now().AddDate(0, 0, 1).Unix()},
					},
					WebhookClauses: []*autoopsproto.WebhookClause{
						{
							WebhookId: "foo-id",
							Conditions: []*autoopsproto.WebhookClause_Condition{
								{
									Filter:   ".foo.bar",
									Value:    "foobaz",
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
							},
						},
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
			desc: "err: ErrWebhookClauseRequired",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id:                       "aid1",
				AddWebhookClauseCommands: []*autoopsproto.AddWebhookClauseCommand{{}},
			},
			expected:    nil,
			expectedErr: createError(statusWebhookClauseRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "webhook_clause")),
		},
		{
			desc: "err: ChangeWebhookClauseCommand: ErrWebhookClauseWebhookClauseConditionRequired",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
				ChangeWebhookClauseCommands: []*autoopsproto.ChangeWebhookClauseCommand{
					{
						Id: "aid",
						WebhookClause: &autoopsproto.WebhookClause{
							WebhookId:  "foo-id",
							Conditions: []*autoopsproto.WebhookClause_Condition{},
						},
					},
				},
			},
			expected:    nil,
			expectedErr: createError(statusWebhookClauseConditionRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "condition")),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().GetGoal(
					gomock.Any(), gomock.Any(),
				).Return(&experimentproto.GetGoalResponse{}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id:                              "aid1",
				EnvironmentNamespace:            "ns0",
				ChangeAutoOpsRuleOpsTypeCommand: &autoopsproto.ChangeAutoOpsRuleOpsTypeCommand{OpsType: autoopsproto.OpsType_DISABLE_FEATURE},
				AddOpsEventRateClauseCommands: []*autoopsproto.AddOpsEventRateClauseCommand{{
					OpsEventRateClause: &autoopsproto.OpsEventRateClause{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					},
				}},
				DeleteClauseCommands: []*autoopsproto.DeleteClauseCommand{{
					Id: "cid",
				}},
				AddDatetimeClauseCommands: []*autoopsproto.AddDatetimeClauseCommand{{
					DatetimeClause: &autoopsproto.DatetimeClause{
						Time: time.Now().AddDate(0, 0, 1).Unix(),
					},
				}},
			},
			expected:    &autoopsproto.UpdateAutoOpsRuleResponse{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.UpdateAutoOpsRule(ctx, p.req)
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.DeleteAutoOpsRuleRequest{
				Id:                   "aid1",
				EnvironmentNamespace: "ns0",
				Command:              &autoopsproto.DeleteAutoOpsRuleCommand{},
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
		setup       func(*AutoOpsService)
		req         *autoopsproto.GetAutoOpsRuleRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &autoopsproto.GetAutoOpsRuleRequest{EnvironmentNamespace: "ns0"},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
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
			req:         &autoopsproto.GetAutoOpsRuleRequest{Id: "wrongid", EnvironmentNamespace: "ns0"},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
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
			req:         &autoopsproto.GetAutoOpsRuleRequest{Id: "aid1", EnvironmentNamespace: "ns0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController, nil)
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

	patterns := []struct {
		setup       func(*AutoOpsService)
		req         *autoopsproto.ListAutoOpsRulesRequest
		expectedErr error
	}{
		{
			setup: func(s *AutoOpsService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			req:         &autoopsproto.ListAutoOpsRulesRequest{EnvironmentNamespace: "ns0", Cursor: ""},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createAutoOpsService(mockController, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.ListAutoOpsRules(createContextWithTokenRoleUnassigned(t), p.req)
		assert.Equal(t, p.expectedErr, err)
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
			desc: "err: ErrNotFound",
			setup: func(s *AutoOpsService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req: &autoopsproto.ExecuteAutoOpsRequest{
				Id:                                  "aid1",
				EnvironmentNamespace:                "ns0",
				ChangeAutoOpsRuleTriggeredAtCommand: &autoopsproto.ChangeAutoOpsRuleTriggeredAtCommand{},
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil).AnyTimes()
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.ExecuteAutoOpsRequest{
				Id:                                  "aid1",
				EnvironmentNamespace:                "ns0",
				ChangeAutoOpsRuleTriggeredAtCommand: &autoopsproto.ChangeAutoOpsRuleTriggeredAtCommand{},
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
			s := createAutoOpsService(mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.existGoal(context.Background(), "ns-0", p.goalID)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func createAutoOpsService(c *gomock.Controller, db storage.Client) *AutoOpsService {
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
	return NewAutoOpsService(
		mysqlClientMock,
		featureClientMock,
		experimentClientMock,
		accountClientMock,
		authClientMock,
		p,
		testWebhookURL,
		&dummyWebhookCryptoUtil{},
		WithLogger(logger),
	)
}

func createContextWithTokenRoleUnassigned(t *testing.T) context.Context {
	t.Helper()
	token := &token.IDToken{
		Issuer:    "issuer",
		Subject:   "sub",
		Audience:  "audience",
		Expiry:    time.Now().AddDate(100, 0, 0),
		IssuedAt:  time.Now(),
		Email:     "email",
		AdminRole: accountproto.Account_UNASSIGNED,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}

func createContextWithTokenRoleOwner(t *testing.T) context.Context {
	t.Helper()
	token := &token.IDToken{
		Issuer:    "issuer",
		Subject:   "sub",
		Audience:  "audience",
		Expiry:    time.Now().AddDate(100, 0, 0),
		IssuedAt:  time.Now(),
		Email:     "email",
		AdminRole: accountproto.Account_OWNER,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}
