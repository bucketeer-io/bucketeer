// Copyright 2022 The Bucketeer Authors.
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

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/eventcounter/domain"
	ecdruid "github.com/bucketeer-io/bucketeer/pkg/eventcounter/druid"
	dmock "github.com/bucketeer-io/bucketeer/pkg/eventcounter/druid/mock"
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

func TestNewEventCounterService(t *testing.T) {
	metrics := metrics.NewMetrics(
		9999,
		"/metrics",
	)
	reg := metrics.DefaultRegisterer()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	g := NewEventCounterService(nil, nil, nil, nil, nil, reg, logger)
	assert.IsType(t, &eventCounterService{}, g)
}

func TestGetEvaluationCountV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	now := time.Now()

	patterns := []struct {
		desc        string
		setup       func(*eventCounterService)
		input       *ecproto.GetEvaluationCountV2Request
		expected    *ecproto.GetEvaluationCountV2Response
		expectedErr error
	}{
		{
			desc: "error: ErrStartAtRequired",
			input: &ecproto.GetEvaluationCountV2Request{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: localizedError(statusStartAtRequired, locale.JaJP),
		},
		{
			desc: "error: ErrEndAtRequired",
			input: &ecproto.GetEvaluationCountV2Request{
				EnvironmentNamespace: "ns0",
				StartAt:              now.Add(-7 * 24 * time.Hour).Unix(),
			},
			expectedErr: localizedError(statusEndAtRequired, locale.JaJP),
		},
		{
			desc: "error: ErrStartAtIsAfterEndAt",
			input: &ecproto.GetEvaluationCountV2Request{
				EnvironmentNamespace: "ns0",
				StartAt:              now.Unix(),
				EndAt:                now.Add(-31 * 24 * time.Hour).Unix(),
			},
			expectedErr: localizedError(statusStartAtIsAfterEndAt, locale.JaJP),
		},
		{
			desc: "error: ErrFeatureIDRequired",
			input: &ecproto.GetEvaluationCountV2Request{
				EnvironmentNamespace: "ns0",
				StartAt:              now.Add(-30 * 24 * time.Hour).Unix(),
				EndAt:                now.Unix(),
			},
			expectedErr: localizedError(statusFeatureIDRequired, locale.JaJP),
		},
		{
			desc: "success: one variation",
			setup: func(s *eventCounterService) {
				s.druidQuerier.(*dmock.MockQuerier).EXPECT().QueryEvaluationCount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&ecproto.Row{Cells: []*ecproto.Cell{
						{Value: ecdruid.ColumnVariation},
						{Value: ecdruid.ColumnEvaluationUser},
						{Value: ecdruid.ColumnEvaluationTotal},
					}},
					[]*ecproto.Row{
						{Cells: []*ecproto.Cell{
							{Value: "vid0", Type: ecproto.Cell_STRING},
							{ValueDouble: float64(1), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(2), Type: ecproto.Cell_DOUBLE},
						}},
						{Cells: []*ecproto.Cell{
							{Value: "vid1", Type: ecproto.Cell_STRING},
							{ValueDouble: float64(12), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(123), Type: ecproto.Cell_DOUBLE},
						}},
					},
					nil)
			},
			input: &ecproto.GetEvaluationCountV2Request{
				EnvironmentNamespace: "ns0",
				StartAt:              now.Add(-30 * 24 * time.Hour).Unix(),
				EndAt:                now.Unix(),
				FeatureId:            "fid",
				FeatureVersion:       int32(1),
				VariationIds:         []string{"vid1"},
			},
			expected: &ecproto.GetEvaluationCountV2Response{
				Count: &ecproto.EvaluationCount{
					FeatureId:      "fid",
					FeatureVersion: int32(1),
					RealtimeCounts: []*ecproto.VariationCount{
						{
							VariationId: "vid1",
							UserCount:   int64(12),
							EventCount:  int64(123),
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: all variations",
			setup: func(s *eventCounterService) {
				s.druidQuerier.(*dmock.MockQuerier).EXPECT().QueryEvaluationCount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&ecproto.Row{Cells: []*ecproto.Cell{
						{Value: ecdruid.ColumnVariation},
						{Value: ecdruid.ColumnEvaluationUser},
						{Value: ecdruid.ColumnEvaluationTotal},
					}},
					[]*ecproto.Row{
						{Cells: []*ecproto.Cell{
							{Value: "vid0", Type: ecproto.Cell_STRING},
							{ValueDouble: float64(1), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(2), Type: ecproto.Cell_DOUBLE},
						}},
						{Cells: []*ecproto.Cell{
							{Value: "vid1", Type: ecproto.Cell_STRING},
							{ValueDouble: float64(12), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(123), Type: ecproto.Cell_DOUBLE},
						}},
					},
					nil)
			},
			input: &ecproto.GetEvaluationCountV2Request{
				EnvironmentNamespace: "ns0",
				StartAt:              now.Add(-30 * 24 * time.Hour).Unix(),
				EndAt:                now.Unix(),
				FeatureId:            "fid",
				FeatureVersion:       int32(1),
				VariationIds:         []string{"vid0", "vid1"},
			},
			expected: &ecproto.GetEvaluationCountV2Response{
				Count: &ecproto.EvaluationCount{
					FeatureId:      "fid",
					FeatureVersion: int32(1),
					RealtimeCounts: []*ecproto.VariationCount{
						{
							VariationId: "vid0",
							UserCount:   int64(1),
							EventCount:  int64(2),
						},
						{
							VariationId: "vid1",
							UserCount:   int64(12),
							EventCount:  int64(123),
						},
					},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		gs := newEventCounterService(t, mockController)
		if p.setup != nil {
			p.setup(gs)
		}
		actual, err := gs.GetEvaluationCountV2(createContextWithToken(t, accountproto.Account_UNASSIGNED), p.input)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
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

	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
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
			expectedErr: localizedError(statusExperimentIDRequired, locale.JaJP),
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
		actual, err := gs.GetExperimentResult(createContextWithToken(t, accountproto.Account_UNASSIGNED), p.input)
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

	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
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
			expectedErr: localizedError(statusFeatureIDRequired, locale.JaJP),
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
			expectedErr: localizedError(statusInternal, locale.JaJP),
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
			actual, err := s.ListExperimentResults(createContextWithToken(t, accountproto.Account_UNASSIGNED), p.input)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetGoalCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	now := time.Now()

	patterns := []struct {
		desc        string
		setup       func(*eventCounterService)
		input       *ecproto.GetGoalCountRequest
		expected    *ecproto.GetGoalCountResponse
		expectedErr error
	}{
		{
			desc: "error: ErrStartAtRequired",
			input: &ecproto.GetGoalCountRequest{
				EnvironmentNamespace: "ns0",
				GoalId:               "gid",
			},
			expectedErr: localizedError(statusStartAtRequired, locale.JaJP),
		},
		{
			desc: "error: ErrEndAtRequired",
			input: &ecproto.GetGoalCountRequest{
				EnvironmentNamespace: "ns0",
				GoalId:               "gid",
				StartAt:              now.Add(-7 * 24 * time.Hour).Unix(),
			},
			expectedErr: localizedError(statusEndAtRequired, locale.JaJP),
		},
		{
			desc: "error: ErrStartAtIsAfterEndAt",
			input: &ecproto.GetGoalCountRequest{
				EnvironmentNamespace: "ns0",
				GoalId:               "gid",
				StartAt:              now.Unix(),
				EndAt:                now.Add(-31 * 24 * time.Hour).Unix(),
			},
			expectedErr: localizedError(statusStartAtIsAfterEndAt, locale.JaJP),
		},
		{
			desc: "error: ErrPeriodOutOfRange",
			input: &ecproto.GetGoalCountRequest{
				EnvironmentNamespace: "ns0",
				GoalId:               "gid",
				StartAt:              now.Add(-32 * 24 * time.Hour).Unix(),
				EndAt:                now.Unix(),
			},
			expectedErr: localizedError(statusPeriodOutOfRange, locale.JaJP),
		},
		{
			desc: "error: ErrGoalIDRequired",
			input: &ecproto.GetGoalCountRequest{
				EnvironmentNamespace: "ns0",
				StartAt:              now.Add(-30 * 24 * time.Hour).Unix(),
				EndAt:                now.Unix(),
			},
			expectedErr: localizedError(statusGoalIDRequired, locale.JaJP),
		},
		{
			desc: "success",
			setup: func(s *eventCounterService) {
				s.druidQuerier.(*dmock.MockQuerier).EXPECT().QueryCount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&ecproto.Row{Cells: []*ecproto.Cell{{Value: "val"}}}, []*ecproto.Row{{Cells: []*ecproto.Cell{{Value: "123"}}}}, nil)
			},
			input: &ecproto.GetGoalCountRequest{
				EnvironmentNamespace: "ns0",
				GoalId:               "gid",
				StartAt:              now.Add(-30 * 24 * time.Hour).Unix(),
				EndAt:                now.Unix(),
			},
			expected: &ecproto.GetGoalCountResponse{
				Headers: &ecproto.Row{Cells: []*ecproto.Cell{{Value: "val"}}},
				Rows:    []*ecproto.Row{{Cells: []*ecproto.Cell{{Value: "123"}}}},
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
			actual, err := s.GetGoalCount(createContextWithToken(t, accountproto.Account_UNASSIGNED), p.input)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetGoalCountV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	now := time.Now()

	patterns := []struct {
		desc        string
		setup       func(*eventCounterService)
		input       *ecproto.GetGoalCountV2Request
		expected    *ecproto.GetGoalCountV2Response
		expectedErr error
	}{
		{
			desc: "error: ErrStartAtRequired",
			input: &ecproto.GetGoalCountV2Request{
				EnvironmentNamespace: "ns0",
				GoalId:               "gid",
			},
			expectedErr: localizedError(statusStartAtRequired, locale.JaJP),
		},
		{
			desc: "error: ErrEndAtRequired",
			input: &ecproto.GetGoalCountV2Request{
				EnvironmentNamespace: "ns0",
				GoalId:               "gid",
				StartAt:              now.Add(-7 * 24 * time.Hour).Unix(),
			},
			expectedErr: localizedError(statusEndAtRequired, locale.JaJP),
		},
		{
			desc: "error: ErrStartAtIsAfterEndAt",
			input: &ecproto.GetGoalCountV2Request{
				EnvironmentNamespace: "ns0",
				GoalId:               "gid",
				StartAt:              now.Unix(),
				EndAt:                now.Add(-31 * 24 * time.Hour).Unix(),
			},
			expectedErr: localizedError(statusStartAtIsAfterEndAt, locale.JaJP),
		},
		{
			desc: "error: ErrGoalIDRequired",
			input: &ecproto.GetGoalCountV2Request{
				EnvironmentNamespace: "ns0",
				StartAt:              now.Add(-31 * 24 * time.Hour).Unix(),
				EndAt:                now.Unix(),
			},
			expectedErr: localizedError(statusGoalIDRequired, locale.JaJP),
		},
		{
			desc: "success: one variation",
			setup: func(s *eventCounterService) {
				s.druidQuerier.(*dmock.MockQuerier).EXPECT().QueryGoalCount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&ecproto.Row{Cells: []*ecproto.Cell{
						{Value: ecdruid.ColumnVariation},
						{Value: ecdruid.ColumnGoalUser},
						{Value: ecdruid.ColumnGoalTotal},
						{Value: ecdruid.ColumnGoalValueTotal},
						{Value: ecdruid.ColumnGoalValueMean},
						{Value: ecdruid.ColumnGoalValueVariance},
					}},
					[]*ecproto.Row{
						{Cells: []*ecproto.Cell{
							{Value: "vid0", Type: ecproto.Cell_STRING},
							{ValueDouble: float64(1), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(2), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(1.23), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(1.234), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(1.2345), Type: ecproto.Cell_DOUBLE},
						}},
						{Cells: []*ecproto.Cell{
							{Value: "vid1", Type: ecproto.Cell_STRING},
							{ValueDouble: float64(12), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(123), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(123.45), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(123.456), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(123.4567), Type: ecproto.Cell_DOUBLE},
						}},
					},
					nil)
			},
			input: &ecproto.GetGoalCountV2Request{
				EnvironmentNamespace: "ns0",
				GoalId:               "gid",
				FeatureId:            "fid",
				FeatureVersion:       int32(1),
				VariationIds:         []string{"vid1"},
				StartAt:              now.Add(-30 * 24 * time.Hour).Unix(),
				EndAt:                now.Unix(),
			},
			expected: &ecproto.GetGoalCountV2Response{
				GoalCounts: &ecproto.GoalCounts{
					GoalId: "gid",
					RealtimeCounts: []*ecproto.VariationCount{
						{
							VariationId:             "vid1",
							UserCount:               int64(12),
							EventCount:              int64(123),
							ValueSum:                float64(123.45),
							ValueSumPerUserMean:     float64(123.456),
							ValueSumPerUserVariance: float64(123.4567),
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: all variations",
			setup: func(s *eventCounterService) {
				s.druidQuerier.(*dmock.MockQuerier).EXPECT().QueryGoalCount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&ecproto.Row{Cells: []*ecproto.Cell{
						{Value: ecdruid.ColumnVariation},
						{Value: ecdruid.ColumnGoalUser},
						{Value: ecdruid.ColumnGoalTotal},
						{Value: ecdruid.ColumnGoalValueTotal},
						{Value: ecdruid.ColumnGoalValueMean},
						{Value: ecdruid.ColumnGoalValueVariance},
					}},
					[]*ecproto.Row{
						{Cells: []*ecproto.Cell{
							{Value: "vid0", Type: ecproto.Cell_STRING},
							{ValueDouble: float64(1), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(2), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(1.23), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(1.234), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(1.2345), Type: ecproto.Cell_DOUBLE},
						}},
						{Cells: []*ecproto.Cell{
							{Value: "vid1", Type: ecproto.Cell_STRING},
							{ValueDouble: float64(12), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(123), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(123.45), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(123.456), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(123.4567), Type: ecproto.Cell_DOUBLE},
						}},
					},
					nil)
			},
			input: &ecproto.GetGoalCountV2Request{
				EnvironmentNamespace: "ns0",
				GoalId:               "gid",
				FeatureId:            "fid",
				FeatureVersion:       int32(1),
				VariationIds:         []string{"vid0", "vid1"},
				StartAt:              now.Add(-30 * 24 * time.Hour).Unix(),
				EndAt:                now.Unix(),
			},
			expected: &ecproto.GetGoalCountV2Response{
				GoalCounts: &ecproto.GoalCounts{
					GoalId: "gid",
					RealtimeCounts: []*ecproto.VariationCount{
						{
							VariationId:             "vid0",
							UserCount:               int64(1),
							EventCount:              int64(2),
							ValueSum:                float64(1.23),
							ValueSumPerUserMean:     float64(1.234),
							ValueSumPerUserVariance: float64(1.2345),
						},
						{
							VariationId:             "vid1",
							UserCount:               int64(12),
							EventCount:              int64(123),
							ValueSum:                float64(123.45),
							ValueSumPerUserMean:     float64(123.456),
							ValueSumPerUserVariance: float64(123.4567),
						},
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
			actual, err := s.GetGoalCountV2(createContextWithToken(t, accountproto.Account_UNASSIGNED), p.input)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetUserCountV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	now := time.Now()

	patterns := []struct {
		desc        string
		setup       func(*eventCounterService)
		input       *ecproto.GetUserCountV2Request
		expected    *ecproto.GetUserCountV2Response
		expectedErr error
	}{
		{
			desc: "error: ErrStartAtRequired",
			input: &ecproto.GetUserCountV2Request{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: localizedError(statusStartAtRequired, locale.JaJP),
		},
		{
			desc: "error: ErrEndAtRequired",
			input: &ecproto.GetUserCountV2Request{
				EnvironmentNamespace: "ns0",
				StartAt:              now.Add(-7 * 24 * time.Hour).Unix(),
			},
			expectedErr: localizedError(statusEndAtRequired, locale.JaJP),
		},
		{
			desc: "error: ErrStartAtIsAfterEndAt",
			input: &ecproto.GetUserCountV2Request{
				EnvironmentNamespace: "ns0",
				StartAt:              now.Unix(),
				EndAt:                now.Add(-31 * 24 * time.Hour).Unix(),
			},
			expectedErr: localizedError(statusStartAtIsAfterEndAt, locale.JaJP),
		},
		{
			desc: "success",
			setup: func(s *eventCounterService) {
				s.druidQuerier.(*dmock.MockQuerier).EXPECT().QueryUserCount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&ecproto.Row{Cells: []*ecproto.Cell{
						{Value: ecdruid.ColumnUserTotal},
						{Value: ecdruid.ColumnUserCount},
					}},
					[]*ecproto.Row{
						{Cells: []*ecproto.Cell{
							{ValueDouble: float64(4), Type: ecproto.Cell_DOUBLE},
							{ValueDouble: float64(2), Type: ecproto.Cell_DOUBLE},
						}},
					},
					nil)
			},
			input: &ecproto.GetUserCountV2Request{
				EnvironmentNamespace: "ns0",
				StartAt:              now.Add(-30 * 24 * time.Hour).Unix(),
				EndAt:                now.Unix(),
			},
			expected: &ecproto.GetUserCountV2Response{
				EventCount: 4,
				UserCount:  2,
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
			actual, err := s.GetUserCountV2(createContextWithToken(t, accountproto.Account_UNASSIGNED), p.input)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListUserMetadata(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*eventCounterService)
		input       *ecproto.ListUserMetadataRequest
		expected    *ecproto.ListUserMetadataResponse
		expectedErr error
	}{
		{
			desc: "success",
			setup: func(s *eventCounterService) {
				s.druidQuerier.(*dmock.MockQuerier).EXPECT().QuerySegmentMetadata(gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"d1", "d2"}, nil)
			},
			input: &ecproto.ListUserMetadataRequest{
				EnvironmentNamespace: "ns0",
			},
			expected: &ecproto.ListUserMetadataResponse{
				Data: []string{"d1", "d2"},
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
			actual, err := s.ListUserMetadata(createContextWithToken(t, accountproto.Account_UNASSIGNED), p.input)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGenInterval(t *testing.T) {
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
			actual, err := genInterval(p.inputLocation, p.inputEndAt, p.inputDurationDays)
			assert.Equal(t, p.expected.Unix(), actual.Unix())
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetEvaluationTimeseriesCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*eventCounterService)
		input       *ecproto.GetEvaluationTimeseriesCountRequest
		expected    *ecproto.GetEvaluationTimeseriesCountResponse
		expectedErr error
	}{
		{
			desc: "error: ErrFeatureIDRequired",
			input: &ecproto.GetEvaluationTimeseriesCountRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: localizedError(statusFeatureIDRequired, locale.JaJP),
		},
		{
			desc: "success",
			setup: func(s *eventCounterService) {
				s.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(gomock.Any(), gomock.Any()).Return(
					&featureproto.GetFeatureResponse{
						Feature: &featureproto.Feature{
							Id:         "fid",
							Variations: []*featureproto.Variation{{Id: "vid0"}, {Id: "vid1"}},
						},
					}, nil)
				s.druidQuerier.(*dmock.MockQuerier).EXPECT().QueryEvaluationTimeseriesCount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					map[string]*ecproto.VariationTimeseries{
						ecdruid.ColumnEvaluationTotal: {
							VariationId: "vid0",
							Timeseries: &ecproto.Timeseries{
								Timestamps: []int64{int64(1)},
								Values:     []float64{float64(1.2)},
							},
						},
						ecdruid.ColumnEvaluationUser: {
							VariationId: "vid0",
							Timeseries: &ecproto.Timeseries{
								Timestamps: []int64{int64(2)},
								Values:     []float64{float64(2.3)},
							},
						},
					}, nil)
				s.druidQuerier.(*dmock.MockQuerier).EXPECT().QueryEvaluationTimeseriesCount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					map[string]*ecproto.VariationTimeseries{
						ecdruid.ColumnEvaluationTotal: {
							VariationId: "vid1",
							Timeseries: &ecproto.Timeseries{
								Timestamps: []int64{int64(3)},
								Values:     []float64{float64(3.4)},
							},
						},
						ecdruid.ColumnEvaluationUser: {
							VariationId: "vid1",
							Timeseries: &ecproto.Timeseries{
								Timestamps: []int64{int64(4)},
								Values:     []float64{float64(4.5)},
							},
						},
					}, nil)
			},
			input: &ecproto.GetEvaluationTimeseriesCountRequest{
				EnvironmentNamespace: "ns0",
				FeatureId:            "fid",
			},
			expected: &ecproto.GetEvaluationTimeseriesCountResponse{
				EventCounts: []*ecproto.VariationTimeseries{
					{
						VariationId: "vid0",
						Timeseries: &ecproto.Timeseries{
							Timestamps: []int64{int64(1)},
							Values:     []float64{float64(1.2)},
						},
					},
					{
						VariationId: "vid1",
						Timeseries: &ecproto.Timeseries{
							Timestamps: []int64{int64(3)},
							Values:     []float64{float64(3.4)},
						},
					},
				},
				UserCounts: []*ecproto.VariationTimeseries{
					{
						VariationId: "vid0",
						Timeseries: &ecproto.Timeseries{
							Timestamps: []int64{int64(2)},
							Values:     []float64{float64(2.3)},
						},
					},
					{
						VariationId: "vid1",
						Timeseries: &ecproto.Timeseries{
							Timestamps: []int64{int64(4)},
							Values:     []float64{float64(4.5)},
						},
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
			actual, err := s.GetEvaluationTimeseriesCount(createContextWithToken(t, accountproto.Account_UNASSIGNED), p.input)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
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
		druidQuerier:                 dmock.NewMockQuerier(mockController),
		metrics:                      reg,
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
