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
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	eccachemock "github.com/bucketeer-io/bucketeer/pkg/cache/v3/mock"
	"github.com/bucketeer-io/bucketeer/pkg/eventcounter/domain"
	v2ecs "github.com/bucketeer-io/bucketeer/pkg/eventcounter/storage/v2"
	v2ecsmock "github.com/bucketeer-io/bucketeer/pkg/eventcounter/storage/v2/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/pkg/experiment/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	jpLocation = time.FixedZone("Asia/Tokyo", 9*60*60)
)

func TestNewEventCounterService(t *testing.T) {
	metrics := metrics.NewMetrics(
		9999,
		"/metrics",
	)
	reg := metrics.DefaultRegisterer()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	g := NewEventCounterService(nil, nil, nil, nil, nil, "", reg, nil, jpLocation, logger)
	assert.IsType(t, &eventCounterService{}, g)
}

func TestGetExperimentEvaluationCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	now := time.Now()
	ctx := createContextWithToken(t, accountproto.Account_UNASSIGNED)
	correctStartAtUnix := now.Add(-30 * 24 * time.Hour).Unix()
	correctStartAt := time.Unix(correctStartAtUnix, 0)
	correctEndAtUnix := now.Unix()
	correctEndAt := time.Unix(correctEndAtUnix, 0)
	ns := "ns0"
	fID := "fid"
	fVersion := int32(1)
	vID1 := "vid01"
	vID2 := "vid02"
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
		setup       func(*eventCounterService)
		input       *ecproto.GetExperimentEvaluationCountRequest
		expected    *ecproto.GetExperimentEvaluationCountResponse
		expectedErr error
	}{
		{
			desc: "error: ErrStartAtRequired",
			input: &ecproto.GetExperimentEvaluationCountRequest{
				EnvironmentNamespace: ns,
				FeatureId:            fID,
				EndAt:                correctEndAtUnix,
			},
			expectedErr: createError(statusStartAtRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "start_at")),
		},
		{
			desc: "error: ErrEndAtRequired",
			input: &ecproto.GetExperimentEvaluationCountRequest{
				EnvironmentNamespace: ns,
				FeatureId:            fID,
				StartAt:              correctStartAtUnix,
			},
			expectedErr: createError(statusEndAtRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "end_at")),
		},
		{
			desc: "error: ErrStartAtIsAfterEndAt",
			input: &ecproto.GetExperimentEvaluationCountRequest{
				EnvironmentNamespace: ns,
				FeatureId:            fID,
				StartAt:              now.Unix(),
				EndAt:                now.Add(-31 * 24 * time.Hour).Unix(),
			},
			expectedErr: createError(statusStartAtIsAfterEndAt, localizer.MustLocalizeWithTemplate(locale.StartAtIsAfterEndAt)),
		},
		{
			desc: "error: ErrFeatureIDRequired",
			input: &ecproto.GetExperimentEvaluationCountRequest{
				EnvironmentNamespace: ns,
				StartAt:              correctStartAtUnix,
				EndAt:                correctEndAtUnix,
			},
			expectedErr: createError(statusFeatureIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id")),
		},
		{
			desc: "success: one variation",
			setup: func(s *eventCounterService) {
				s.eventStorage.(*v2ecsmock.MockEventStorage).EXPECT().QueryEvaluationCount(ctx, ns, correctStartAt, correctEndAt, fID, fVersion).Return(
					[]*v2ecs.EvaluationEventCount{
						{
							VariationID:     vID1,
							EvaluationUser:  int64(1),
							EvaluationTotal: int64(2),
						},
					},
					nil,
				)
			},
			input: &ecproto.GetExperimentEvaluationCountRequest{
				EnvironmentNamespace: ns,
				StartAt:              correctStartAtUnix,
				EndAt:                correctEndAtUnix,
				FeatureId:            fID,
				FeatureVersion:       fVersion,
				VariationIds:         []string{vID1},
			},
			expected: &ecproto.GetExperimentEvaluationCountResponse{
				FeatureId:      fID,
				FeatureVersion: fVersion,
				VariationCounts: []*ecproto.VariationCount{
					{
						VariationId: vID1,
						UserCount:   int64(1),
						EventCount:  int64(2),
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: all variations",
			setup: func(s *eventCounterService) {
				s.eventStorage.(*v2ecsmock.MockEventStorage).EXPECT().QueryEvaluationCount(ctx, ns, correctStartAt, correctEndAt, fID, fVersion).Return(
					[]*v2ecs.EvaluationEventCount{
						{
							VariationID:     vID1,
							EvaluationUser:  int64(1),
							EvaluationTotal: int64(2),
						},
						{
							VariationID:     vID2,
							EvaluationUser:  int64(12),
							EvaluationTotal: int64(123),
						},
					},
					nil)
			},
			input: &ecproto.GetExperimentEvaluationCountRequest{
				EnvironmentNamespace: ns,
				StartAt:              correctStartAtUnix,
				EndAt:                correctEndAtUnix,
				FeatureId:            fID,
				FeatureVersion:       fVersion,
				VariationIds:         []string{vID1, vID2},
			},
			expected: &ecproto.GetExperimentEvaluationCountResponse{
				FeatureId:      fID,
				FeatureVersion: fVersion,
				VariationCounts: []*ecproto.VariationCount{
					{
						VariationId: vID1,
						UserCount:   int64(1),
						EventCount:  int64(2),
					},
					{
						VariationId: vID2,
						UserCount:   int64(12),
						EventCount:  int64(123),
					},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newEventCounterService(t, mockController)
			if p.setup != nil {
				p.setup(gs)
			}
			actual, err := gs.GetExperimentEvaluationCount(ctx, p.input)
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestListExperiments(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc                 string
		setup                func(*eventCounterService)
		inputFeatureID       string
		inputFeatureVersion  *wrappers.Int32Value
		expected             []*experimentproto.Experiment
		environmentNamespace string
		expectedErr          error
	}{
		{
			desc: "no error",
			setup: func(s *eventCounterService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().ListExperiments(gomock.Any(), &experimentproto.ListExperimentsRequest{
					FeatureId:            "fid",
					FeatureVersion:       &wrappers.Int32Value{Value: int32(1)},
					PageSize:             listRequestPageSize,
					Cursor:               "",
					EnvironmentNamespace: "ns0",
				}).Return(&experimentproto.ListExperimentsResponse{}, nil)
			},
			inputFeatureID:       "fid",
			inputFeatureVersion:  &wrappers.Int32Value{Value: int32(1)},
			environmentNamespace: "ns0",
			expected:             []*experimentproto.Experiment{},
			expectedErr:          nil,
		},
		{
			desc: "error",
			setup: func(s *eventCounterService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().ListExperiments(gomock.Any(), &experimentproto.ListExperimentsRequest{
					FeatureId:            "fid",
					FeatureVersion:       &wrappers.Int32Value{Value: int32(1)},
					PageSize:             listRequestPageSize,
					Cursor:               "",
					EnvironmentNamespace: "ns0",
				}).Return(nil, errors.New("test"))
			},
			inputFeatureID:       "fid",
			inputFeatureVersion:  &wrappers.Int32Value{Value: int32(1)},
			environmentNamespace: "ns0",
			expected:             nil,
			expectedErr:          errors.New("test"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newEventCounterService(t, mockController)
			p.setup(s)
			actual, err := s.listExperiments(context.Background(), p.inputFeatureID, p.inputFeatureVersion, p.environmentNamespace)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetExperimentResultMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, accountproto.Account_UNASSIGNED)
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
		setup       func(*eventCounterService)
		input       *ecproto.GetExperimentResultRequest
		expectedErr error
	}{
		{
			desc:        "error: ErrExperimentIDRequired",
			input:       &ecproto.GetExperimentResultRequest{EnvironmentNamespace: "ns0"},
			expectedErr: createError(statusExperimentIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "experiment_id")),
		},
		{
			desc: "err: ErrNotFound",
			setup: func(s *eventCounterService) {
				s.mysqlExperimentResultStorage.(*v2ecsmock.MockExperimentResultStorage).EXPECT().GetExperimentResult(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, v2ecs.ErrExperimentResultNotFound)
			},
			input: &ecproto.GetExperimentResultRequest{
				ExperimentId:         "eid",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success: get the result from storage",
			setup: func(s *eventCounterService) {
				s.mysqlExperimentResultStorage.(*v2ecsmock.MockExperimentResultStorage).EXPECT().GetExperimentResult(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.ExperimentResult{}, nil)
			},
			input: &ecproto.GetExperimentResultRequest{
				ExperimentId:         "eid",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		gs := newEventCounterService(t, mockController)
		if p.setup != nil {
			p.setup(gs)
		}
		actual, err := gs.GetExperimentResult(ctx, p.input)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		if err == nil {
			assert.NotNil(t, actual)
		}
	}
}

func TestListExperimentResultsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, accountproto.Account_UNASSIGNED)
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
		setup       func(*eventCounterService)
		input       *ecproto.ListExperimentResultsRequest
		expected    *ecproto.ListExperimentResultsResponse
		expectedErr error
	}{
		{
			desc:        "error: ErrFeatureIDRequired",
			input:       &ecproto.ListExperimentResultsRequest{EnvironmentNamespace: "ns0"},
			expectedErr: createError(statusFeatureIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id")),
		},
		{
			desc: "err: ErrNotFound",
			setup: func(s *eventCounterService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().ListExperiments(
					gomock.Any(), gomock.Any(),
				).Return(nil, storage.ErrKeyNotFound)
			},
			input: &ecproto.ListExperimentResultsRequest{
				FeatureId:            "fid",
				EnvironmentNamespace: "ns0",
			},
			expected:    nil,
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *eventCounterService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().ListExperiments(
					gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("test"))
			},
			input: &ecproto.ListExperimentResultsRequest{
				FeatureId:            "fid",
				FeatureVersion:       &wrappers.Int32Value{Value: int32(1)},
				EnvironmentNamespace: "ns0",
			},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success: no results",
			setup: func(s *eventCounterService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().ListExperiments(
					gomock.Any(), gomock.Any(),
				).Return(
					&experimentproto.ListExperimentsResponse{
						Experiments: []*experimentproto.Experiment{
							{
								Id:             "eid",
								GoalId:         "gid",
								FeatureId:      "fid",
								FeatureVersion: int32(1),
							},
						},
					},
					nil,
				)
				s.mysqlExperimentResultStorage.(*v2ecsmock.MockExperimentResultStorage).EXPECT().GetExperimentResult(
					gomock.Any(), "eid", gomock.Any(),
				).Return(nil, v2ecs.ErrExperimentResultNotFound)
			},
			input: &ecproto.ListExperimentResultsRequest{
				FeatureId:            "fid",
				FeatureVersion:       &wrappers.Int32Value{Value: int32(1)},
				EnvironmentNamespace: "ns0",
			},
			expected: &ecproto.ListExperimentResultsResponse{
				Results: make(map[string]*ecproto.ExperimentResult, 0),
			},
			expectedErr: nil,
		},
		{
			desc: "success: get results from storage",
			setup: func(s *eventCounterService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().ListExperiments(
					gomock.Any(), gomock.Any(),
				).Return(
					&experimentproto.ListExperimentsResponse{
						Experiments: []*experimentproto.Experiment{
							{
								Id:             "eid",
								GoalId:         "gid",
								FeatureId:      "fid",
								FeatureVersion: int32(1),
							},
						},
					},
					nil,
				)
				s.mysqlExperimentResultStorage.(*v2ecsmock.MockExperimentResultStorage).EXPECT().GetExperimentResult(
					gomock.Any(), "eid", gomock.Any(),
				).Return(
					&domain.ExperimentResult{
						ExperimentResult: &ecproto.ExperimentResult{
							Id:          "eid",
							GoalResults: []*ecproto.GoalResult{{GoalId: "gid"}},
						},
					},
					nil,
				)
			},
			input: &ecproto.ListExperimentResultsRequest{
				FeatureId:            "fid",
				FeatureVersion:       &wrappers.Int32Value{Value: int32(1)},
				EnvironmentNamespace: "ns0",
			},
			expected: &ecproto.ListExperimentResultsResponse{
				Results: map[string]*ecproto.ExperimentResult{
					"eid": {
						Id:          "eid",
						GoalResults: []*ecproto.GoalResult{{GoalId: "gid"}},
					},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newEventCounterService(t, mockController)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.ListExperimentResults(ctx, p.input)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetExperimentGoalCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	now := time.Now()
	ctx := createContextWithToken(t, accountproto.Account_UNASSIGNED)
	correctStartAtUnix := now.Add(-30 * 24 * time.Hour).Unix()
	correctStartAt := time.Unix(correctStartAtUnix, 0)
	correctEndAtUnix := now.Unix()
	correctEndAt := time.Unix(correctEndAtUnix, 0)
	ns := "ns0"
	fID := "fid"
	fVersion := int32(1)
	vID1 := "vid01"
	vID2 := "vid02"
	gID := "gid"
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
		setup       func(*eventCounterService)
		input       *ecproto.GetExperimentGoalCountRequest
		expected    *ecproto.GetExperimentGoalCountResponse
		expectedErr error
	}{
		{
			desc: "error: ErrStartAtRequired",
			input: &ecproto.GetExperimentGoalCountRequest{
				EnvironmentNamespace: ns,
				FeatureId:            fID,
				GoalId:               gID,
				EndAt:                correctEndAtUnix,
			},
			expectedErr: createError(statusStartAtRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "start_at")),
		},
		{
			desc: "error: ErrEndAtRequired",
			input: &ecproto.GetExperimentGoalCountRequest{
				EnvironmentNamespace: ns,
				FeatureId:            fID,
				GoalId:               gID,
				StartAt:              correctStartAtUnix,
			},
			expectedErr: createError(statusEndAtRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "end_at")),
		},
		{
			desc: "error: ErrStartAtIsAfterEndAt",
			input: &ecproto.GetExperimentGoalCountRequest{
				EnvironmentNamespace: ns,
				FeatureId:            fID,
				GoalId:               gID,
				StartAt:              now.Unix(),
				EndAt:                now.Add(-30 * 24 * time.Hour).Unix(),
			},
			expectedErr: createError(statusStartAtIsAfterEndAt, localizer.MustLocalizeWithTemplate(locale.StartAtIsAfterEndAt)),
		},
		{
			desc: "error: ErrFeatureIDRequired",
			input: &ecproto.GetExperimentGoalCountRequest{
				EnvironmentNamespace: ns,
				GoalId:               gID,
				StartAt:              correctStartAtUnix,
				EndAt:                correctEndAtUnix,
			},
			expectedErr: createError(statusFeatureIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id")),
		},
		{
			desc: "error: ErrGoalIDRequired",
			input: &ecproto.GetExperimentGoalCountRequest{
				EnvironmentNamespace: ns,
				FeatureId:            fID,
				StartAt:              correctStartAtUnix,
				EndAt:                correctEndAtUnix,
			},
			expectedErr: createError(statusGoalIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id")),
		},
		{
			desc: "success: one variation",
			setup: func(s *eventCounterService) {
				s.eventStorage.(*v2ecsmock.MockEventStorage).EXPECT().QueryGoalCount(ctx, ns, correctStartAt, correctEndAt, gID, fID, fVersion).Return(
					[]*v2ecs.GoalEventCount{
						{
							VariationID:       vID1,
							GoalUser:          int64(1),
							GoalTotal:         int64(2),
							GoalValueTotal:    1.23,
							GoalValueMean:     1.234,
							GoalValueVariance: 1.2345,
						},
					},
					nil,
				)
			},
			input: &ecproto.GetExperimentGoalCountRequest{
				EnvironmentNamespace: ns,
				GoalId:               gID,
				FeatureId:            fID,
				FeatureVersion:       fVersion,
				VariationIds:         []string{vID1},
				StartAt:              correctStartAtUnix,
				EndAt:                correctEndAtUnix,
			},
			expected: &ecproto.GetExperimentGoalCountResponse{
				GoalId: gID,
				VariationCounts: []*ecproto.VariationCount{
					{
						VariationId:             vID1,
						UserCount:               int64(1),
						EventCount:              int64(2),
						ValueSum:                1.23,
						ValueSumPerUserMean:     1.234,
						ValueSumPerUserVariance: 1.2345,
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: all variations",
			setup: func(s *eventCounterService) {
				s.eventStorage.(*v2ecsmock.MockEventStorage).EXPECT().QueryGoalCount(ctx, ns, correctStartAt, correctEndAt, gID, fID, fVersion).Return(
					[]*v2ecs.GoalEventCount{
						{
							VariationID:       vID1,
							GoalUser:          int64(1),
							GoalTotal:         int64(2),
							GoalValueTotal:    1.23,
							GoalValueMean:     1.234,
							GoalValueVariance: 1.2345,
						},
						{
							VariationID:       vID2,
							GoalUser:          int64(12),
							GoalTotal:         int64(123),
							GoalValueTotal:    123.45,
							GoalValueMean:     123.456,
							GoalValueVariance: 123.4567,
						},
					},
					nil,
				)
			},
			input: &ecproto.GetExperimentGoalCountRequest{
				EnvironmentNamespace: ns,
				GoalId:               gID,
				FeatureId:            fID,
				FeatureVersion:       fVersion,
				VariationIds:         []string{vID1, vID2},
				StartAt:              correctStartAtUnix,
				EndAt:                correctEndAtUnix,
			},
			expected: &ecproto.GetExperimentGoalCountResponse{
				GoalId: gID,
				VariationCounts: []*ecproto.VariationCount{
					{
						VariationId:             vID1,
						UserCount:               int64(1),
						EventCount:              int64(2),
						ValueSum:                1.23,
						ValueSumPerUserMean:     1.234,
						ValueSumPerUserVariance: 1.2345,
					},
					{
						VariationId:             vID2,
						UserCount:               int64(12),
						EventCount:              int64(123),
						ValueSum:                123.45,
						ValueSumPerUserMean:     123.456,
						ValueSumPerUserVariance: 123.4567,
					},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newEventCounterService(t, mockController)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.GetExperimentGoalCount(ctx, p.input)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetMAUCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithToken(t, accountproto.Account_UNASSIGNED)
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
	input := &ecproto.GetMAUCountRequest{
		EnvironmentNamespace: "ns0",
		YearMonth:            "201212",
	}
	patterns := []struct {
		desc        string
		setup       func(*eventCounterService)
		input       *ecproto.GetMAUCountRequest
		expected    *ecproto.GetMAUCountResponse
		expectedErr error
	}{
		{
			desc:     "error: mau year month is required",
			input:    &ecproto.GetMAUCountRequest{EnvironmentNamespace: "ns0"},
			expected: nil,
			expectedErr: createError(
				statusMAUYearMonthRequired,
				localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "year_month"),
			),
		},
		{
			desc: "err: internal",
			setup: func(s *eventCounterService) {
				s.userCountStorage.(*v2ecsmock.MockUserCountStorage).EXPECT().GetMAUCount(
					ctx, input.EnvironmentNamespace, input.YearMonth,
				).Return(int64(0), int64(0), errors.New("internal"))
			},
			input:       input,
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(s *eventCounterService) {
				s.userCountStorage.(*v2ecsmock.MockUserCountStorage).EXPECT().GetMAUCount(
					ctx, input.EnvironmentNamespace, input.YearMonth,
				).Return(int64(2), int64(4), nil)
			},
			input: input,
			expected: &ecproto.GetMAUCountResponse{
				UserCount:  2,
				EventCount: 4,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newEventCounterService(t, mockController)
			if p.setup != nil {
				p.setup(gs)
			}
			actual, err := gs.GetMAUCount(ctx, p.input)
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestSummarizeMAUCounts(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithToken(t, accountproto.Account_VIEWER)
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
	input := &ecproto.SummarizeMAUCountsRequest{
		YearMonth:  "201212",
		IsFinished: false,
	}
	patterns := []struct {
		desc        string
		setup       func(*eventCounterService)
		input       *ecproto.SummarizeMAUCountsRequest
		expected    *ecproto.SummarizeMAUCountsResponse
		expectedErr error
	}{
		{
			desc:     "error: mau year month is required",
			input:    &ecproto.SummarizeMAUCountsRequest{},
			expected: nil,
			expectedErr: createError(
				statusMAUYearMonthRequired,
				localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "year_month"),
			),
		},
		{
			desc: "err: internal",
			setup: func(s *eventCounterService) {
				s.userCountStorage.(*v2ecsmock.MockUserCountStorage).EXPECT().GetMAUCountsGroupBySourceID(
					ctx, input.YearMonth,
				).Return([]*ecproto.MAUSummary{}, nil)
				s.userCountStorage.(*v2ecsmock.MockUserCountStorage).EXPECT().GetMAUCounts(
					ctx, input.YearMonth,
				).Return(nil, errors.New("internal"))
			},
			input:       input,
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success get mau counts",
			setup: func(s *eventCounterService) {
				s.userCountStorage.(*v2ecsmock.MockUserCountStorage).EXPECT().GetMAUCountsGroupBySourceID(
					ctx, input.YearMonth,
				).Return([]*ecproto.MAUSummary{}, nil)
				s.userCountStorage.(*v2ecsmock.MockUserCountStorage).EXPECT().GetMAUCounts(
					ctx, input.YearMonth,
				).Return([]*ecproto.MAUSummary{}, nil)
			},
			input:       input,
			expected:    &ecproto.SummarizeMAUCountsResponse{},
			expectedErr: nil,
		},
		{
			desc: "success upsert mau summary",
			setup: func(s *eventCounterService) {
				s.userCountStorage.(*v2ecsmock.MockUserCountStorage).EXPECT().GetMAUCountsGroupBySourceID(
					ctx, input.YearMonth,
				).Return([]*ecproto.MAUSummary{{Yearmonth: input.YearMonth}}, nil)
				s.userCountStorage.(*v2ecsmock.MockUserCountStorage).EXPECT().GetMAUCounts(
					ctx, input.YearMonth,
				).Return([]*ecproto.MAUSummary{}, nil)
				s.mysqlMAUSummaryStorage.(*v2ecsmock.MockMAUSummaryStorage).EXPECT().UpsertMAUSummary(
					ctx, gomock.Any(),
				).Return(nil)
			},
			input:       input,
			expected:    &ecproto.SummarizeMAUCountsResponse{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newEventCounterService(t, mockController)
			if p.setup != nil {
				p.setup(gs)
			}
			actual, err := gs.SummarizeMAUCounts(ctx, p.input)
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGetStartTime(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc              string
		inputLocation     *time.Location
		inputEndAt        time.Time
		inputDurationDays int
		expected          time.Time
		expectedErr       error
	}{
		{
			desc:              "success",
			inputLocation:     jpLocation,
			inputEndAt:        time.Date(2020, 12, 25, 0, 0, 0, 0, time.UTC),
			inputDurationDays: 10,
			expected:          time.Date(2020, 12, 15, 0, 0, 0, 0, jpLocation),
			expectedErr:       nil,
		},
		{
			desc:              "over prime meridian",
			inputLocation:     jpLocation,
			inputEndAt:        time.Date(2020, 12, 25, 23, 0, 0, 0, time.UTC),
			inputDurationDays: 10,
			expected:          time.Date(2020, 12, 16, 0, 0, 0, 0, jpLocation),
			expectedErr:       nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := truncateDate(p.inputLocation, getStartTime(p.inputLocation, p.inputEndAt, p.inputDurationDays))
			assert.Equal(t, p.expected.Unix(), actual.Unix())
		})
	}
}

func TestGetDailyTimestamps(t *testing.T) {
	t.Parallel()

	endAt := time.Date(2020, 12, 25, 8, 0, 0, 0, jpLocation)
	startAt := truncateDate(jpLocation, getStartTime(jpLocation, endAt, 3))

	patterns := []struct {
		desc             string
		startAt          time.Time
		expectedElements []int64
		expectedLen      int
	}{
		{
			desc:    "success",
			startAt: startAt,
			expectedElements: []int64{
				getDate(startAt),
				getDate(startAt.AddDate(0, 0, 1)),
				getDate(startAt.AddDate(0, 0, 2)),
				getDate(endAt),
			},
			expectedLen: 4,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := getDailyTimestamps(p.startAt, 3)
			assert.Equal(t, p.expectedElements, actual)
			assert.Len(t, actual, p.expectedLen)
		})
	}
}

func TestGetOneDayTimestamps(t *testing.T) {
	t.Parallel()

	endAt := time.Now()
	startAt := truncateDate(jpLocation, getStartTime(jpLocation, endAt, 30))

	patterns := []struct {
		desc             string
		startAt          time.Time
		expectedElements []int64
		expectedLen      int
	}{
		{
			desc:    "success",
			startAt: startAt,
			expectedElements: []int64{
				getDate(startAt),
				getDate(startAt.Add(1 * time.Hour)),
				getDate(startAt.Add(23 * time.Hour)),
			},
			expectedLen: 24,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := getOneDayTimestamps(p.startAt)
			for _, e := range p.expectedElements {
				assert.Contains(t, actual, int64(e))
			}
			assert.Len(t, actual, p.expectedLen)
		})
	}
}

func TestGetHourlyTimestamps(t *testing.T) {
	t.Parallel()

	endAt := time.Now()
	startAt := getStartTime(jpLocation, endAt, 3)
	expected := [][]int64{}
	expected = append(expected, getOneDayTimestamps(startAt))
	expected = append(expected, getOneDayTimestamps(startAt.AddDate(0, 0, 1)))
	expected = append(expected, getOneDayTimestamps(startAt.AddDate(0, 0, 2)))
	expected = append(expected, getOneDayTimestamps(startAt.AddDate(0, 0, 3)))
	daily := getDailyTimestamps(startAt, 3)
	actual := getHourlyTimeStamps(daily, ecproto.Timeseries_DAY)
	assert.Equal(t, expected, actual)
}

func TestGetTotalEventCounts(t *testing.T) {
	mockController := gomock.NewController(t)
	patterns := []struct {
		desc     string
		input    []float64
		expected int64
	}{
		{
			desc: "success: integer",
			input: []float64{
				1,
				3,
				2,
			},
			expected: 6,
		},
		{
			desc: "success: float",
			input: []float64{
				1.3,
				3.9,
				2.0,
			},
			expected: 7,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newEventCounterService(t, mockController)
			actual := gs.getTotalEventCounts(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestGetTimestamps(t *testing.T) {
	mockController := gomock.NewController(t)
	patterns := []struct {
		desc         string
		input        ecproto.GetEvaluationTimeseriesCountRequest_TimeRange
		expectedLen  int
		expectedUnit ecproto.Timeseries_Unit
		expectedErr  error
	}{
		{
			desc:        "fail: errUnknownTimeRange",
			input:       ecproto.GetEvaluationTimeseriesCountRequest_TimeRange(100),
			expectedErr: errUnknownTimeRange,
		},
		{
			desc:         "success: TWENTY_FOUR_HOURS",
			input:        ecproto.GetEvaluationTimeseriesCountRequest_TWENTY_FOUR_HOURS,
			expectedLen:  24,
			expectedUnit: ecproto.Timeseries_HOUR,
		},
		{
			desc:         "success: SEVEN_DAYS",
			input:        ecproto.GetEvaluationTimeseriesCountRequest_SEVEN_DAYS,
			expectedLen:  7,
			expectedUnit: ecproto.Timeseries_DAY,
		},
		{
			desc:         "success: FOURTEEN_DAYS",
			input:        ecproto.GetEvaluationTimeseriesCountRequest_FOURTEEN_DAYS,
			expectedLen:  14,
			expectedUnit: ecproto.Timeseries_DAY,
		},
		{
			desc:         "success: THIRTY_DAYS",
			input:        ecproto.GetEvaluationTimeseriesCountRequest_THIRTY_DAYS,
			expectedLen:  30,
			expectedUnit: ecproto.Timeseries_DAY,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newEventCounterService(t, mockController)
			timestamps, timestampUnit, err := gs.getTimestamps(p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Len(t, timestamps, p.expectedLen)
			assert.Equal(t, timestampUnit, p.expectedUnit)
		})
	}
}

func TestGetVariationIDs(t *testing.T) {
	t.Parallel()

	vID1 := newUUID(t)
	vID2 := newUUID(t)

	patterns := []struct {
		desc       string
		variations []*featureproto.Variation
		expected   []string
	}{
		{
			desc: "success",
			variations: []*featureproto.Variation{
				{
					Id:    vID1,
					Value: "true",
				},
				{
					Id:    vID2,
					Value: "false",
				},
			},
			expected: []string{
				vID1, vID2, defaultVariationID,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := getVariationIDs(p.variations)
			assert.Equal(t, actual, p.expected)
		})
	}
}

func TestGetEvaluationTimeseriesCount(t *testing.T) {
	t.Parallel()

	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithToken(t, accountproto.Account_UNASSIGNED)
	environmentNamespace := "ns0"
	fID := "fid"
	vID0 := "vid0"
	vID1 := "vid1"
	randomNumberGroup := getRandomNumberGroup(3)
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
		setup       func(context.Context, *eventCounterService)
		input       *ecproto.GetEvaluationTimeseriesCountRequest
		expected    *ecproto.GetEvaluationTimeseriesCountResponse
		expectedErr error
	}{
		{
			desc: "error: ErrFeatureIDRequired",
			input: &ecproto.GetEvaluationTimeseriesCountRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(statusFeatureIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id")),
		},
		{
			desc: "error: ErrUnknownTimeRange",
			input: &ecproto.GetEvaluationTimeseriesCountRequest{
				EnvironmentNamespace: "ns0",
				FeatureId:            fID,
			},
			expectedErr: createError(statusUnknownTimeRange, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time_range")),
		},
		{
			desc: "error: get feature failed",
			setup: func(ctx context.Context, s *eventCounterService) {
				s.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(ctx, &featureproto.GetFeatureRequest{
					EnvironmentNamespace: environmentNamespace,
					Id:                   fID,
				}).Return(
					&featureproto.GetFeatureResponse{
						Feature: &featureproto.Feature{
							Id:         "fid",
							Variations: []*featureproto.Variation{{Id: "vid0"}, {Id: "vid1"}},
						},
					}, errors.New("error"))
			},
			input: &ecproto.GetEvaluationTimeseriesCountRequest{
				EnvironmentNamespace: environmentNamespace,
				FeatureId:            fID,
				TimeRange:            ecproto.GetEvaluationTimeseriesCountRequest_FOURTEEN_DAYS,
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "error: get event counts failed",
			setup: func(ctx context.Context, s *eventCounterService) {
				s.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(ctx, &featureproto.GetFeatureRequest{
					EnvironmentNamespace: environmentNamespace,
					Id:                   fID,
				}).Return(
					&featureproto.GetFeatureResponse{
						Feature: &featureproto.Feature{
							Id:         "fid",
							Variations: []*featureproto.Variation{{Id: "vid0"}, {Id: "vid1"}},
						},
					}, nil)
				s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().GetEventCountsV2(gomock.Any()).Return(
					nil, errors.New("error"))
			},
			input: &ecproto.GetEvaluationTimeseriesCountRequest{
				EnvironmentNamespace: environmentNamespace,
				FeatureId:            fID,
				TimeRange:            ecproto.GetEvaluationTimeseriesCountRequest_FOURTEEN_DAYS,
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "error: MergeMultiKeys failed",
			setup: func(ctx context.Context, s *eventCounterService) {
				s.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(ctx, &featureproto.GetFeatureRequest{
					EnvironmentNamespace: environmentNamespace,
					Id:                   fID,
				}).Return(
					&featureproto.GetFeatureResponse{
						Feature: &featureproto.Feature{
							Id:         "fid",
							Variations: []*featureproto.Variation{{Id: "vid0"}, {Id: "vid1"}},
						},
					}, nil)
				s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().GetEventCountsV2(gomock.Any()).Return(
					[]float64{
						1, 3, 5,
					}, nil)
				s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().MergeMultiKeys(gomock.Any(), gomock.Any()).Return(errors.New("error1"))
			},
			input: &ecproto.GetEvaluationTimeseriesCountRequest{
				EnvironmentNamespace: environmentNamespace,
				FeatureId:            fID,
				TimeRange:            ecproto.GetEvaluationTimeseriesCountRequest_FOURTEEN_DAYS,
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "error: GetUserCountsV2 failed",
			setup: func(ctx context.Context, s *eventCounterService) {
				s.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(ctx, &featureproto.GetFeatureRequest{
					EnvironmentNamespace: environmentNamespace,
					Id:                   fID,
				}).Return(
					&featureproto.GetFeatureResponse{
						Feature: &featureproto.Feature{
							Id:         "fid",
							Variations: []*featureproto.Variation{{Id: "vid0"}, {Id: "vid1"}},
						},
					}, nil)
				s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().GetEventCountsV2(gomock.Any()).Return(
					[]float64{
						1, 3, 5,
					}, nil)
				s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().MergeMultiKeys(gomock.Any(), gomock.Any()).Return(nil)
				s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().GetUserCount(gomock.Any()).Return(int64(0), errors.New("error"))
				s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().DeleteKey(gomock.Any()).Return(nil)
			},
			input: &ecproto.GetEvaluationTimeseriesCountRequest{
				EnvironmentNamespace: environmentNamespace,
				FeatureId:            fID,
				TimeRange:            ecproto.GetEvaluationTimeseriesCountRequest_FOURTEEN_DAYS,
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "error: DeleteKey failed",
			setup: func(ctx context.Context, s *eventCounterService) {
				s.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(ctx, &featureproto.GetFeatureRequest{
					EnvironmentNamespace: environmentNamespace,
					Id:                   fID,
				}).Return(
					&featureproto.GetFeatureResponse{
						Feature: &featureproto.Feature{
							Id:         "fid",
							Variations: []*featureproto.Variation{{Id: "vid0"}, {Id: "vid1"}},
						},
					}, nil)
				s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().GetEventCountsV2(gomock.Any()).Return(
					[]float64{
						1, 3, 5,
					}, nil)
				s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().MergeMultiKeys(gomock.Any(), gomock.Any()).Return(nil)
				s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().GetUserCount(gomock.Any()).Return(int64(0), nil)
				s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().DeleteKey(gomock.Any()).Return(errors.New("error1"))
			},
			input: &ecproto.GetEvaluationTimeseriesCountRequest{
				EnvironmentNamespace: environmentNamespace,
				FeatureId:            fID,
				TimeRange:            ecproto.GetEvaluationTimeseriesCountRequest_FOURTEEN_DAYS,
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "success",
			setup: func(ctx context.Context, s *eventCounterService) {
				s.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(ctx, &featureproto.GetFeatureRequest{
					EnvironmentNamespace: environmentNamespace,
					Id:                   fID,
				}).Return(
					&featureproto.GetFeatureResponse{
						Feature: &featureproto.Feature{
							Id:         "fid",
							Variations: []*featureproto.Variation{{Id: vID0}, {Id: vID1}},
						},
					}, nil)
				vIDs := []string{vID0, vID1, defaultVariationID}
				hourlyTimeStamps := getFourteenDaysTimestamps()
				for idx, vID := range vIDs {
					ec := getEventCountKeysV2(vID, fID, environmentNamespace, hourlyTimeStamps)
					val := randomNumberGroup[idx]
					s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().GetEventCountsV2(ec).Return(
						val, nil)
					uc := getUserCountKeysV2(vID, fID, environmentNamespace, hourlyTimeStamps)
					pfMergeKey := newPFMergeKey(
						UserCountPrefix,
						fID,
						environmentNamespace,
					)
					for idx := range hourlyTimeStamps {
						s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().MergeMultiKeys(pfMergeKey, uc[idx]).Return(nil)
						s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().DeleteKey(pfMergeKey).Return(nil)
						s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().GetUserCount(pfMergeKey).Return(int64(0), nil)
					}
					s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().MergeMultiKeys(pfMergeKey, s.flattenAry(uc)).Return(nil)
					s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().DeleteKey(pfMergeKey).Return(nil)
					s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().GetUserCount(pfMergeKey).Return(int64(0), nil)
				}
			},
			input: &ecproto.GetEvaluationTimeseriesCountRequest{
				EnvironmentNamespace: environmentNamespace,
				FeatureId:            fID,
				TimeRange:            ecproto.GetEvaluationTimeseriesCountRequest_FOURTEEN_DAYS,
			},
			expected: &ecproto.GetEvaluationTimeseriesCountResponse{
				EventCounts: []*ecproto.VariationTimeseries{
					{
						VariationId: vID0,
					},
					{
						VariationId: vID1,
					},
					{
						VariationId: defaultVariationID,
					},
				},
				UserCounts: []*ecproto.VariationTimeseries{
					{
						VariationId: vID0,
					},
					{
						VariationId: vID1,
					},
					{
						VariationId: defaultVariationID,
					},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newEventCounterService(t, mockController)
			if p.setup != nil {
				p.setup(ctx, s)
			}
			actual, err := s.GetEvaluationTimeseriesCount(ctx, p.input)
			if p.expectedErr == nil {
				for idx := range p.expected.EventCounts {
					actualTs := actual.EventCounts[idx]
					assert.Equal(t, p.expected.EventCounts[idx].VariationId, actualTs.VariationId)
					assert.Equal(t, randomNumberGroup[idx], actualTs.Timeseries.Values)
					assert.Len(t, actualTs.Timeseries.Timestamps, 14)
				}
				for idx := range p.expected.UserCounts {
					actualTs := actual.EventCounts[idx]
					assert.Equal(t, p.expected.UserCounts[idx].VariationId, actualTs.VariationId)
					assert.Equal(t, randomNumberGroup[idx], actualTs.Timeseries.Values)
					assert.Len(t, actualTs.Timeseries.Timestamps, 14)
				}
			}
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetOpsEvaluationUserCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithToken(t, accountproto.Account_UNASSIGNED)
	environmentNamespace := "ns0"
	opsRuleID := "rule0"
	clauseID := "clause0"
	fID := "fid0"
	fVersion := 2
	vID0 := "vid0"
	cacheKey := "ns0:autoops:evaluation:fid0:2:rule0:clause0:vid0"
	cacheKeyWithoutNS := "autoops:evaluation:fid0:2:rule0:clause0:vid0"
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
		setup       func(*eventCounterService)
		input       *ecproto.GetOpsEvaluationUserCountRequest
		expected    *ecproto.GetOpsEvaluationUserCountResponse
		expectedErr error
	}{
		{
			desc: "error: ErrOpsRuleIDRequired",
			input: &ecproto.GetOpsEvaluationUserCountRequest{
				EnvironmentNamespace: environmentNamespace,
				ClauseId:             clauseID,
				FeatureId:            fID,
				FeatureVersion:       int32(fVersion),
				VariationId:          vID0,
			},
			expectedErr: createError(statusAutoOpsRuleIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "ops_rule_id")),
		},
		{
			desc: "error: ErrClauseIDRequired",
			input: &ecproto.GetOpsEvaluationUserCountRequest{
				EnvironmentNamespace: environmentNamespace,
				OpsRuleId:            opsRuleID,
				FeatureId:            fID,
				FeatureVersion:       int32(fVersion),
				VariationId:          vID0,
			},
			expectedErr: createError(statusClauseIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause_id")),
		},
		{
			desc: "error: ErrFeatureIDRequired",
			input: &ecproto.GetOpsEvaluationUserCountRequest{
				EnvironmentNamespace: environmentNamespace,
				OpsRuleId:            opsRuleID,
				ClauseId:             clauseID,
				FeatureVersion:       int32(fVersion),
				VariationId:          vID0,
			},
			expectedErr: createError(statusFeatureIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id")),
		},
		{
			desc: "error: ErrFeatureVersionRequired",
			input: &ecproto.GetOpsEvaluationUserCountRequest{
				EnvironmentNamespace: environmentNamespace,
				OpsRuleId:            opsRuleID,
				ClauseId:             clauseID,
				FeatureId:            fID,
				VariationId:          vID0,
			},
			expectedErr: createError(statusFeatureVersionRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_version")),
		},
		{
			desc: "error: ErrVariationIDRequired",
			input: &ecproto.GetOpsEvaluationUserCountRequest{
				EnvironmentNamespace: environmentNamespace,
				OpsRuleId:            opsRuleID,
				ClauseId:             clauseID,
				FeatureId:            fID,
				FeatureVersion:       int32(fVersion),
			},
			expectedErr: createError(statusVariationIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			desc: "success",
			setup: func(s *eventCounterService) {
				s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().
					GetUserCount(cacheKey).Return(int64(1234), nil)
			},
			input: &ecproto.GetOpsEvaluationUserCountRequest{
				EnvironmentNamespace: environmentNamespace,
				OpsRuleId:            opsRuleID,
				ClauseId:             clauseID,
				FeatureId:            fID,
				FeatureVersion:       int32(fVersion),
				VariationId:          vID0,
			},
			expected: &ecproto.GetOpsEvaluationUserCountResponse{
				OpsRuleId: opsRuleID,
				ClauseId:  clauseID,
				Count:     1234,
			},
			expectedErr: nil,
		},
		{
			desc: "success: without environment_namespace",
			setup: func(s *eventCounterService) {
				s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().
					GetUserCount(cacheKeyWithoutNS).Return(int64(9876), nil)
			},
			input: &ecproto.GetOpsEvaluationUserCountRequest{
				OpsRuleId:      opsRuleID,
				ClauseId:       clauseID,
				FeatureId:      fID,
				FeatureVersion: int32(fVersion),
				VariationId:    vID0,
			},
			expected: &ecproto.GetOpsEvaluationUserCountResponse{
				OpsRuleId: opsRuleID,
				ClauseId:  clauseID,
				Count:     9876,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newEventCounterService(t, mockController)
			if p.setup != nil {
				p.setup(gs)
			}
			actual, err := gs.GetOpsEvaluationUserCount(ctx, p.input)
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGetOpsGoalUserCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithToken(t, accountproto.Account_UNASSIGNED)
	environmentNamespace := "ns0"
	opsRuleID := "rule0"
	clauseID := "clause0"
	fID := "fid0"
	fVersion := 2
	vID0 := "vid0"
	cacheKey := "ns0:autoops:goal:fid0:2:rule0:clause0:vid0"
	cacheKeyWithoutNS := "autoops:goal:fid0:2:rule0:clause0:vid0"
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
		setup       func(*eventCounterService)
		input       *ecproto.GetOpsGoalUserCountRequest
		expected    *ecproto.GetOpsGoalUserCountResponse
		expectedErr error
	}{
		{
			desc: "error: ErrOpsRuleIDRequired",
			input: &ecproto.GetOpsGoalUserCountRequest{
				EnvironmentNamespace: environmentNamespace,
				ClauseId:             clauseID,
				FeatureId:            fID,
				FeatureVersion:       int32(fVersion),
				VariationId:          vID0,
			},
			expectedErr: createError(statusAutoOpsRuleIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "ops_rule_id")),
		},
		{
			desc: "error: ErrClauseIDRequired",
			input: &ecproto.GetOpsGoalUserCountRequest{
				EnvironmentNamespace: environmentNamespace,
				OpsRuleId:            opsRuleID,
				FeatureId:            fID,
				FeatureVersion:       int32(fVersion),
				VariationId:          vID0,
			},
			expectedErr: createError(statusClauseIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause_id")),
		},
		{
			desc: "error: ErrFeatureIDRequired",
			input: &ecproto.GetOpsGoalUserCountRequest{
				EnvironmentNamespace: environmentNamespace,
				OpsRuleId:            opsRuleID,
				ClauseId:             clauseID,
				FeatureVersion:       int32(fVersion),
				VariationId:          vID0,
			},
			expectedErr: createError(statusFeatureIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id")),
		},
		{
			desc: "error: ErrFeatureVersionRequired",
			input: &ecproto.GetOpsGoalUserCountRequest{
				EnvironmentNamespace: environmentNamespace,
				OpsRuleId:            opsRuleID,
				ClauseId:             clauseID,
				FeatureId:            fID,
				VariationId:          vID0,
			},
			expectedErr: createError(statusFeatureVersionRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_version")),
		},
		{
			desc: "error: ErrVariationIDRequired",
			input: &ecproto.GetOpsGoalUserCountRequest{
				EnvironmentNamespace: environmentNamespace,
				OpsRuleId:            opsRuleID,
				ClauseId:             clauseID,
				FeatureId:            fID,
				FeatureVersion:       int32(fVersion),
			},
			expectedErr: createError(statusVariationIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id")),
		},
		{
			desc: "success",
			setup: func(s *eventCounterService) {
				s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().
					GetUserCount(cacheKey).Return(int64(1234), nil)
			},
			input: &ecproto.GetOpsGoalUserCountRequest{
				EnvironmentNamespace: environmentNamespace,
				OpsRuleId:            opsRuleID,
				ClauseId:             clauseID,
				FeatureId:            fID,
				FeatureVersion:       int32(fVersion),
				VariationId:          vID0,
			},
			expected: &ecproto.GetOpsGoalUserCountResponse{
				OpsRuleId: opsRuleID,
				ClauseId:  clauseID,
				Count:     1234,
			},
			expectedErr: nil,
		},
		{
			desc: "success: without environment_namespace",
			setup: func(s *eventCounterService) {
				s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().
					GetUserCount(cacheKeyWithoutNS).Return(int64(9876), nil)
			},
			input: &ecproto.GetOpsGoalUserCountRequest{
				OpsRuleId:      opsRuleID,
				ClauseId:       clauseID,
				FeatureId:      fID,
				FeatureVersion: int32(fVersion),
				VariationId:    vID0,
			},
			expected: &ecproto.GetOpsGoalUserCountResponse{
				OpsRuleId: opsRuleID,
				ClauseId:  clauseID,
				Count:     9876,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newEventCounterService(t, mockController)
			if p.setup != nil {
				p.setup(gs)
			}
			actual, err := gs.GetOpsGoalUserCount(ctx, p.input)
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func getRandomNumberGroup(size int) [][]float64 {
	group := make([][]float64, 0, size)
	for i := 0; i < size; i++ {
		group = append(group, getRandomNumbers())
	}
	return group
}

func getRandomNumbers() []float64 {
	size := 31
	nums := make([]float64, 0, size)
	for i := 0; i < size; i++ {
		nums = append(nums, rand.Float64())
	}
	return nums
}

func getEventCountKeys(vID, fID, environmentNamespace string, timeStamps []int64) []string {
	eventCountKeys := []string{}
	for _, ts := range timeStamps {
		ec := newEvaluationCountkey(EventCountPrefix, fID, vID, environmentNamespace, ts)
		eventCountKeys = append(eventCountKeys, ec)
	}
	return eventCountKeys
}

func getUserCountKeys(vID, fid, environmentNamespace string, timeStamps []int64) []string {
	userCountKeys := []string{}
	for _, ts := range timeStamps {
		uc := newEvaluationCountkey(UserCountPrefix, fid, vID, environmentNamespace, ts)
		userCountKeys = append(userCountKeys, uc)
	}
	return userCountKeys
}

func getEventCountKeysV2(vID, fID, environmentNamespace string, timeStamps [][]int64) [][]string {
	eventCountKeys := [][]string{}
	for _, day := range timeStamps {
		hourly := []string{}
		for _, hour := range day {
			ec := newEvaluationCountkey(EventCountPrefix, fID, vID, environmentNamespace, hour)
			hourly = append(hourly, ec)
		}
		eventCountKeys = append(eventCountKeys, hourly)
	}
	return eventCountKeys
}

func getUserCountKeysV2(vID, fID, environmentNamespace string, timeStamps [][]int64) [][]string {
	userCountKeys := [][]string{}
	for _, day := range timeStamps {
		hourly := []string{}
		for _, hour := range day {
			ec := newEvaluationCountkey(UserCountPrefix, fID, vID, environmentNamespace, hour)
			hourly = append(hourly, ec)
		}
		userCountKeys = append(userCountKeys, hourly)
	}
	return userCountKeys
}

func getDate(t time.Time) int64 {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, jpLocation).Unix()
}

func newEventCounterService(t *testing.T, mockController *gomock.Controller) *eventCounterService {
	logger, err := log.NewLogger()
	require.NoError(t, err)
	metrics := metrics.NewMetrics(
		9999,
		"/metrics",
	)
	reg := metrics.DefaultRegisterer()
	accountClientMock := accountclientmock.NewMockClient(mockController)
	ar := &accountproto.GetAccountResponse{
		Account: &accountproto.Account{
			Email: "email",
			Role:  accountproto.Account_VIEWER,
		},
	}
	accountClientMock.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	return &eventCounterService{
		experimentClient:             experimentclientmock.NewMockClient(mockController),
		featureClient:                featureclientmock.NewMockClient(mockController),
		accountClient:                accountClientMock,
		mysqlExperimentResultStorage: v2ecsmock.NewMockExperimentResultStorage(mockController),
		mysqlMAUSummaryStorage:       v2ecsmock.NewMockMAUSummaryStorage(mockController),
		userCountStorage:             v2ecsmock.NewMockUserCountStorage(mockController),
		evaluationCountCacher:        eccachemock.NewMockEventCounterCache(mockController),
		eventStorage:                 v2ecsmock.NewMockEventStorage(mockController),
		metrics:                      reg,
		location:                     jpLocation,
		logger:                       logger.Named("api"),
	}
}

func createContextWithToken(t *testing.T, role accountproto.Account_Role) context.Context {
	t.Helper()
	token := &token.IDToken{
		Email:     "test@example.com",
		AdminRole: role,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}

func newUUID(t *testing.T) string {
	t.Helper()
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}

func TestMultiError(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		err      multiError
		expected string
	}{
		{
			desc: "2 errors",
			err: multiError{
				errors.New("foobar"),
				errors.New("hoge"),
			},
			expected: "2 errors: foobar, hoge",
		},
		{
			desc: "1 error",
			err: multiError{
				errors.New("foobar"),
			},
			expected: "1 errors: foobar",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := fmt.Errorf("%v", p.err).Error()
			assert.Equal(t, p.expected, actual)
		})
	}
}

func getFourteenDaysTimestamps() [][]int64 {
	endAt := time.Now()
	startAt := truncateDate(jpLocation, getStartTime(jpLocation, endAt, 13))
	dailyTimeStamps := getDailyTimestamps(startAt, 13)
	return getHourlyTimeStamps(dailyTimeStamps, ecproto.Timeseries_DAY)
}

func getTwentyFourHoursTimestamps() []int64 {
	endAt := time.Now()
	startAt := truncateHour(jpLocation, getStartTime(jpLocation, endAt, 1))
	return getOneDayTimestamps(startAt)
}

func TestGetUserCounts(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	featureID := "fID"
	environmentNamespace := "en"
	vID := "vID"
	fourteenDaysKeys := getUserCountKeysV2(vID, featureID, environmentNamespace, getFourteenDaysTimestamps())
	twentyFourHoursKeys := getUserCountKeysV2(vID, featureID, environmentNamespace, [][]int64{getTwentyFourHoursTimestamps()})
	patterns := []struct {
		desc        string
		unit        ecproto.Timeseries_Unit
		keys        [][]string
		setup       func(*eventCounterService)
		expectedLen int
	}{
		{
			desc: "success: 14 days timestamp",
			unit: ecproto.Timeseries_DAY,
			keys: fourteenDaysKeys,
			setup: func(s *eventCounterService) {
				for _, day := range fourteenDaysKeys {
					s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().
						GetUserCount(gomock.Any()).Return(int64(1234), nil)
					s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().DeleteKey(gomock.Any()).Return(nil)
					s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().MergeMultiKeys(gomock.Any(), day).Return(nil)
				}
			},
			expectedLen: 14,
		},
		{
			desc: "success: 24 hours timestamp",
			unit: ecproto.Timeseries_HOUR,
			keys: twentyFourHoursKeys,
			setup: func(s *eventCounterService) {
				for _, hour := range twentyFourHoursKeys[0] {
					s.evaluationCountCacher.(*eccachemock.MockEventCounterCache).EXPECT().
						GetUserCount(hour).Return(int64(1234), nil)
				}
			},
			expectedLen: 24,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newEventCounterService(t, mockController)
			if p.setup != nil {
				p.setup(gs)
			}
			actual, _ := gs.getUserCounts(p.keys, featureID, environmentNamespace, p.unit)
			assert.Len(t, actual, p.expectedLen)
		})
	}
}

func TestCheckAdminRole(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := context.Background()
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
		inputCtx    context.Context
		expectedErr error
	}{
		{
			desc:        "error: Unauthenticated",
			inputCtx:    context.Background(),
			expectedErr: createError(statusUnauthenticated, localizer.MustLocalizeWithTemplate(locale.UnauthenticatedError)),
		},
		{
			desc:        "error: PermissionDenied",
			inputCtx:    createContextWithToken(t, accountproto.Account_UNASSIGNED),
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalizeWithTemplate(locale.PermissionDenied)),
		},
		{
			desc:        "success",
			inputCtx:    createContextWithToken(t, accountproto.Account_EDITOR),
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newEventCounterService(t, mockController)
			_, actualErr := gs.checkAdminRole(p.inputCtx, localizer)
			assert.Equal(t, actualErr, p.expectedErr)
		})
	}
}
