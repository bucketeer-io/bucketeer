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
	"errors"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	insightsstorage "github.com/bucketeer-io/bucketeer/v2/pkg/insights/storage/v2"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/insights/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	prometheusmock "github.com/bucketeer-io/bucketeer/v2/pkg/prometheus/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	clientproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	insightsproto "github.com/bucketeer-io/bucketeer/v2/proto/insights"
)

func TestNewInsightsService(t *testing.T) {
	t.Parallel()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	s := NewInsightsService(nil, nil, nil, WithLogger(logger))
	assert.IsType(t, &insightsService{}, s)
}

func TestGetInsightsMonthlySummary(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*insightsService)
		noProm      bool
		orgRole     *accountproto.AccountV2_Role_Organization
		envRole     *accountproto.AccountV2_Role_Environment
		input       *insightsproto.GetInsightsMonthlySummaryRequest
		expected    *insightsproto.GetInsightsMonthlySummaryResponse
		expectedErr error
	}{
		{
			desc:        "error: ErrEnvironmentIDRequired",
			input:       &insightsproto.GetInsightsMonthlySummaryRequest{},
			expectedErr: statusEnvironmentIDRequired.Err(),
		},
		{
			desc:    "error: ErrPermissionDenied",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_UNASSIGNED),
			input: &insightsproto.GetInsightsMonthlySummaryRequest{
				EnvironmentIds: []string{"env1"},
				SourceIds:      []clientproto.SourceId{clientproto.SourceId_ANDROID},
			},
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc:    "error: storage error",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *insightsService) {
				s.monthlySummaryStorage.(*storagemock.MockMonthlySummaryStorage).EXPECT().
					ListMonthlySummaries(gomock.Any(), []string{"env1"}, []string{"ANDROID"}).
					Return(nil, errors.New("storage error"))
			},
			input: &insightsproto.GetInsightsMonthlySummaryRequest{
				EnvironmentIds: []string{"env1"},
				SourceIds:      []clientproto.SourceId{clientproto.SourceId_ANDROID},
			},
			expectedErr: statusInternal.Err(),
		},
		{
			desc:    "success: empty result",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *insightsService) {
				s.monthlySummaryStorage.(*storagemock.MockMonthlySummaryStorage).EXPECT().
					ListMonthlySummaries(gomock.Any(), []string{"env1"}, []string{"ANDROID"}).
					Return([]insightsstorage.ListMonthlySummaryResult{}, nil)
			},
			input: &insightsproto.GetInsightsMonthlySummaryRequest{
				EnvironmentIds: []string{"env1"},
				SourceIds:      []clientproto.SourceId{clientproto.SourceId_ANDROID},
			},
			expected: &insightsproto.GetInsightsMonthlySummaryResponse{
				Series: []*insightsproto.MonthlySummarySeries{},
			},
			expectedErr: nil,
		},
		{
			desc:    "success: with data",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *insightsService) {
				s.monthlySummaryStorage.(*storagemock.MockMonthlySummaryStorage).EXPECT().
					ListMonthlySummaries(gomock.Any(), []string{"env1"}, []string{"ANDROID"}).
					Return([]insightsstorage.ListMonthlySummaryResult{
						{
							Yearmonth:       "202601",
							EnvironmentID:   "env1",
							EnvironmentName: "Environment 1",
							ProjectName:     "Project 1",
							SourceID:        "ANDROID",
							MAU:             1000,
							Requests:        50000,
						},
						{
							Yearmonth:       "202602",
							EnvironmentID:   "env1",
							EnvironmentName: "Environment 1",
							ProjectName:     "Project 1",
							SourceID:        "ANDROID",
							MAU:             1200,
							Requests:        60000,
						},
					}, nil)
			},
			input: &insightsproto.GetInsightsMonthlySummaryRequest{
				EnvironmentIds: []string{"env1"},
				SourceIds:      []clientproto.SourceId{clientproto.SourceId_ANDROID},
			},
			expected: &insightsproto.GetInsightsMonthlySummaryResponse{
				Series: []*insightsproto.MonthlySummarySeries{
					{
						EnvironmentId:   "env1",
						EnvironmentName: "Environment 1",
						ProjectName:     "Project 1",
						SourceId:        clientproto.SourceId_ANDROID,
						Data: []*insightsproto.MonthlySummaryDataPoint{
							{Yearmonth: "202601", Mau: 1000, Requests: 50000},
							{Yearmonth: "202602", Mau: 1200, Requests: 60000},
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc:    "success: without source filter returns all sources",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *insightsService) {
				s.monthlySummaryStorage.(*storagemock.MockMonthlySummaryStorage).EXPECT().
					ListMonthlySummaries(gomock.Any(), []string{"env1"}, allSourceIDStrings()).
					Return([]insightsstorage.ListMonthlySummaryResult{
						{
							Yearmonth:       "202601",
							EnvironmentID:   "env1",
							EnvironmentName: "Environment 1",
							ProjectName:     "Project 1",
							SourceID:        "ANDROID",
							MAU:             1000,
							Requests:        50000,
						},
						{
							Yearmonth:       "202601",
							EnvironmentID:   "env1",
							EnvironmentName: "Environment 1",
							ProjectName:     "Project 1",
							SourceID:        "IOS",
							MAU:             800,
							Requests:        40000,
						},
					}, nil)
			},
			input: &insightsproto.GetInsightsMonthlySummaryRequest{
				EnvironmentIds: []string{"env1"},
				// SourceIds are not set
			},
			expected: &insightsproto.GetInsightsMonthlySummaryResponse{
				Series: []*insightsproto.MonthlySummarySeries{
					{
						EnvironmentId:   "env1",
						EnvironmentName: "Environment 1",
						ProjectName:     "Project 1",
						SourceId:        clientproto.SourceId_ANDROID,
						Data: []*insightsproto.MonthlySummaryDataPoint{
							{Yearmonth: "202601", Mau: 1000, Requests: 50000},
						},
					},
					{
						EnvironmentId:   "env1",
						EnvironmentName: "Environment 1",
						ProjectName:     "Project 1",
						SourceId:        clientproto.SourceId_IOS,
						Data: []*insightsproto.MonthlySummaryDataPoint{
							{Yearmonth: "202601", Mau: 800, Requests: 40000},
						},
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			var s *insightsService
			if p.noProm {
				s = newInsightsServiceWithoutProm(t, mockController, p.orgRole, p.envRole)
			} else {
				s = newInsightsServiceForTest(t, mockController, p.orgRole, p.envRole)
			}
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.GetInsightsMonthlySummary(ctx, p.input)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
			if p.expected != nil && actual != nil {
				assert.Len(t, actual.Series, len(p.expected.Series))
			}
		})
	}
}

func TestGetInsightsLatency(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	now := time.Now()
	startAt := now.Add(-24 * time.Hour).Unix()
	endAt := now.Unix()

	patterns := []struct {
		desc        string
		setup       func(*insightsService)
		noProm      bool
		orgRole     *accountproto.AccountV2_Role_Organization
		envRole     *accountproto.AccountV2_Role_Environment
		input       *insightsproto.GetInsightsTimeSeriesRequest
		expected    *insightsproto.GetInsightsTimeSeriesResponse
		expectedErr error
	}{
		{
			desc:        "error: data source not configured",
			noProm:      true,
			input:       &insightsproto.GetInsightsTimeSeriesRequest{EnvironmentIds: []string{"env1"}, StartAt: startAt, EndAt: endAt},
			expectedErr: statusDataSourceNotConfigured.Err(),
		},
		{
			desc:        "error: ErrEnvironmentIDRequired",
			input:       &insightsproto.GetInsightsTimeSeriesRequest{},
			expectedErr: statusEnvironmentIDRequired.Err(),
		},
		{
			desc: "error: ErrStartAtRequired",
			input: &insightsproto.GetInsightsTimeSeriesRequest{
				EnvironmentIds: []string{"env1"},
				EndAt:          endAt,
			},
			expectedErr: statusStartAtRequired.Err(),
		},
		{
			desc: "error: ErrEndAtRequired",
			input: &insightsproto.GetInsightsTimeSeriesRequest{
				EnvironmentIds: []string{"env1"},
				StartAt:        startAt,
			},
			expectedErr: statusEndAtRequired.Err(),
		},
		{
			desc: "error: ErrStartAtIsAfterEndAt",
			input: &insightsproto.GetInsightsTimeSeriesRequest{
				EnvironmentIds: []string{"env1"},
				StartAt:        endAt,
				EndAt:          startAt,
			},
			expectedErr: statusStartAtIsAfterEndAt.Err(),
		},
		{
			desc: "error: ErrQueryRangeTooLarge",
			input: &insightsproto.GetInsightsTimeSeriesRequest{
				EnvironmentIds: []string{"env1"},
				StartAt:        now.Add(-32 * 24 * time.Hour).Unix(),
				EndAt:          now.Unix(),
			},
			expectedErr: statusQueryRangeTooLarge.Err(),
		},
		{
			desc:    "error: ErrPermissionDenied",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_UNASSIGNED),
			input: &insightsproto.GetInsightsTimeSeriesRequest{
				EnvironmentIds: []string{"env1"},
				StartAt:        startAt,
				EndAt:          endAt,
			},
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc:    "error: prometheus error",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *insightsService) {
				s.promClient.(*prometheusmock.MockClient).EXPECT().
					QueryRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("prometheus error"))
			},
			input: &insightsproto.GetInsightsTimeSeriesRequest{
				EnvironmentIds: []string{"env1"},
				StartAt:        startAt,
				EndAt:          endAt,
			},
			expectedErr: statusInternal.Err(),
		},
		{
			desc:    "success",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *insightsService) {
				s.promClient.(*prometheusmock.MockClient).EXPECT().
					QueryRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(model.Matrix{
						&model.SampleStream{
							Metric: model.Metric{
								"environment_id": "env1",
								"source_id":      "ANDROID",
								"method":         "GetEvaluations",
							},
							Values: []model.SamplePair{
								{Timestamp: 1704067200000, Value: 0.05},
								{Timestamp: 1704070800000, Value: 0.06},
							},
						},
					}, nil)
			},
			input: &insightsproto.GetInsightsTimeSeriesRequest{
				EnvironmentIds: []string{"env1"},
				StartAt:        startAt,
				EndAt:          endAt,
			},
			expected: &insightsproto.GetInsightsTimeSeriesResponse{
				Timeseries: []*insightsproto.InsightsTimeSeries{
					{
						EnvironmentId: "env1",
						SourceId:      clientproto.SourceId_ANDROID,
						ApiId:         clientproto.ApiId_GET_EVALUATIONS,
						Data: []*insightsproto.InsightsDataPoint{
							{Timestamp: 1704067200, Value: 0.05},
							{Timestamp: 1704070800, Value: 0.06},
						},
					},
				},
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			var s *insightsService
			if p.noProm {
				s = newInsightsServiceWithoutProm(t, mockController, p.orgRole, p.envRole)
			} else {
				s = newInsightsServiceForTest(t, mockController, p.orgRole, p.envRole)
			}
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.GetInsightsLatency(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if p.expected != nil {
				assert.Equal(t, p.expected, actual)
			}
		})
	}
}

func TestGetInsightsErrorRates(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	now := time.Now()
	startAt := now.Add(-24 * time.Hour).Unix()
	endAt := now.Unix()

	patterns := []struct {
		desc        string
		setup       func(*insightsService)
		orgRole     *accountproto.AccountV2_Role_Organization
		envRole     *accountproto.AccountV2_Role_Environment
		input       *insightsproto.GetInsightsTimeSeriesRequest
		expected    *insightsproto.GetInsightsTimeSeriesResponse
		expectedErr error
	}{
		{
			desc: "error: prometheus not configured",
			input: &insightsproto.GetInsightsTimeSeriesRequest{
				EnvironmentIds: []string{"env1"},
				StartAt:        startAt,
				EndAt:          endAt,
			},
			expectedErr: statusDataSourceNotConfigured.Err(),
		},
		{
			desc:    "success",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *insightsService) {
				s.promClient.(*prometheusmock.MockClient).EXPECT().
					QueryRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(model.Matrix{
						&model.SampleStream{
							Metric: model.Metric{
								"environment_id": "env1",
								"source_id":      "ANDROID",
								"method":         "GetEvaluations",
							},
							Values: []model.SamplePair{
								{Timestamp: 1704067200000, Value: 0.02},
							},
						},
					}, nil)
			},
			input: &insightsproto.GetInsightsTimeSeriesRequest{
				EnvironmentIds: []string{"env1"},
				StartAt:        startAt,
				EndAt:          endAt,
			},
			expected: &insightsproto.GetInsightsTimeSeriesResponse{
				Timeseries: []*insightsproto.InsightsTimeSeries{
					{
						EnvironmentId: "env1",
						SourceId:      clientproto.SourceId_ANDROID,
						ApiId:         clientproto.ApiId_GET_EVALUATIONS,
						Data: []*insightsproto.InsightsDataPoint{
							{Timestamp: 1704067200, Value: 0.02},
						},
					},
				},
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			var s *insightsService
			if p.desc == "error: prometheus not configured" {
				s = newInsightsServiceWithoutProm(t, mockController, p.orgRole, p.envRole)
			} else {
				s = newInsightsServiceForTest(t, mockController, p.orgRole, p.envRole)
			}
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.GetInsightsErrorRates(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if p.expected != nil {
				assert.Equal(t, p.expected, actual)
			}
		})
	}
}

func TestConvertMatrixToTimeSeries(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc      string
		input     model.Matrix
		labelKeys []string
		expected  []*insightsproto.InsightsTimeSeries
	}{
		{
			desc:     "empty matrix",
			input:    model.Matrix{},
			expected: []*insightsproto.InsightsTimeSeries{},
		},
		{
			desc: "single stream",
			input: model.Matrix{
				&model.SampleStream{
					Metric: model.Metric{
						"environment_id": "env1",
						"source_id":      "ANDROID",
						"method":         "GetEvaluations",
					},
					Values: []model.SamplePair{
						{Timestamp: 1704067200000, Value: 0.05},
						{Timestamp: 1704070800000, Value: 0.06},
					},
				},
			},
			expected: []*insightsproto.InsightsTimeSeries{
				{
					EnvironmentId: "env1",
					SourceId:      clientproto.SourceId_ANDROID,
					ApiId:         clientproto.ApiId_GET_EVALUATIONS,
					Data: []*insightsproto.InsightsDataPoint{
						{Timestamp: 1704067200, Value: 0.05},
						{Timestamp: 1704070800, Value: 0.06},
					},
				},
			},
		},
		{
			desc: "multi stream",
			input: model.Matrix{
				&model.SampleStream{
					Metric: model.Metric{
						"environment_id": "env1",
						"source_id":      "ANDROID",
						"method":         "GetEvaluations",
					},
					Values: []model.SamplePair{
						{Timestamp: 1704067200000, Value: 0.05},
						{Timestamp: 1704070800000, Value: 0.06},
					},
				},
				&model.SampleStream{
					Metric: model.Metric{
						"environment_id": "env1",
						"source_id":      "IOS",
						"method":         "GetEvaluations",
					},
					Values: []model.SamplePair{
						{Timestamp: 1704067200000, Value: 0.10},
						{Timestamp: 1704070800000, Value: 0.12},
					},
				},
				&model.SampleStream{
					Metric: model.Metric{
						"environment_id": "env2",
						"source_id":      "ANDROID",
						"method":         "GetEvaluations",
					},
					Values: []model.SamplePair{
						{Timestamp: 1704067200000, Value: 0.03},
					},
				},
			},
			expected: []*insightsproto.InsightsTimeSeries{
				{
					EnvironmentId: "env1",
					SourceId:      clientproto.SourceId_ANDROID,
					ApiId:         clientproto.ApiId_GET_EVALUATIONS,
					Data: []*insightsproto.InsightsDataPoint{
						{Timestamp: 1704067200, Value: 0.05},
						{Timestamp: 1704070800, Value: 0.06},
					},
				},
				{
					EnvironmentId: "env1",
					SourceId:      clientproto.SourceId_IOS,
					ApiId:         clientproto.ApiId_GET_EVALUATIONS,
					Data: []*insightsproto.InsightsDataPoint{
						{Timestamp: 1704067200, Value: 0.10},
						{Timestamp: 1704070800, Value: 0.12},
					},
				},
				{
					EnvironmentId: "env2",
					SourceId:      clientproto.SourceId_ANDROID,
					ApiId:         clientproto.ApiId_GET_EVALUATIONS,
					Data: []*insightsproto.InsightsDataPoint{
						{Timestamp: 1704067200, Value: 0.03},
					},
				},
			},
		},
		{
			desc: "with labelKeys extracts only specified labels",
			input: model.Matrix{
				&model.SampleStream{
					Metric: model.Metric{
						"environment_id":  "env1",
						"source_id":       "ANDROID",
						"method":          "GetEvaluations",
						"evaluation_type": "diff",
						"pod":             "api-server-abc123", // ignored
					},
					Values: []model.SamplePair{
						{Timestamp: 1704067200000, Value: 100},
					},
				},
			},
			labelKeys: []string{"evaluation_type"},
			expected: []*insightsproto.InsightsTimeSeries{
				{
					EnvironmentId: "env1",
					SourceId:      clientproto.SourceId_ANDROID,
					ApiId:         clientproto.ApiId_GET_EVALUATIONS,
					Data: []*insightsproto.InsightsDataPoint{
						{Timestamp: 1704067200, Value: 100},
					},
					Labels: map[string]string{"evaluation_type": "diff"},
				},
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			actual := convertMatrixToTimeSeries(p.input, p.labelKeys)
			assert.ElementsMatch(t, p.expected, actual)
		})
	}
}

func TestConvertToMonthlySummarySeries(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		input    []insightsstorage.ListMonthlySummaryResult
		expected []*insightsproto.MonthlySummarySeries
	}{
		{
			desc:     "empty records",
			input:    []insightsstorage.ListMonthlySummaryResult{},
			expected: []*insightsproto.MonthlySummarySeries{},
		},
		{
			desc: "single series with multiple yearmonths",
			input: []insightsstorage.ListMonthlySummaryResult{
				{Yearmonth: "202601", EnvironmentID: "env1", EnvironmentName: "Env1", ProjectName: "Proj1", SourceID: "ANDROID", MAU: 100, Requests: 1000},
				{Yearmonth: "202602", EnvironmentID: "env1", EnvironmentName: "Env1", ProjectName: "Proj1", SourceID: "ANDROID", MAU: 120, Requests: 1200},
				{Yearmonth: "202603", EnvironmentID: "env1", EnvironmentName: "Env1", ProjectName: "Proj1", SourceID: "ANDROID", MAU: 150, Requests: 1500},
			},
			expected: []*insightsproto.MonthlySummarySeries{
				{
					EnvironmentId:   "env1",
					EnvironmentName: "Env1",
					ProjectName:     "Proj1",
					SourceId:        clientproto.SourceId_ANDROID,
					Data: []*insightsproto.MonthlySummaryDataPoint{
						{Yearmonth: "202601", Mau: 100, Requests: 1000},
						{Yearmonth: "202602", Mau: 120, Requests: 1200},
						{Yearmonth: "202603", Mau: 150, Requests: 1500},
					},
				},
			},
		},
		{
			desc: "multiple series",
			input: []insightsstorage.ListMonthlySummaryResult{
				{Yearmonth: "202601", EnvironmentID: "env1", EnvironmentName: "Env1", ProjectName: "Proj1", SourceID: "ANDROID", MAU: 100, Requests: 1000},
				{Yearmonth: "202601", EnvironmentID: "env1", EnvironmentName: "Env1", ProjectName: "Proj1", SourceID: "GO_SERVER", MAU: 50, Requests: 500},
				{Yearmonth: "202602", EnvironmentID: "env1", EnvironmentName: "Env1", ProjectName: "Proj1", SourceID: "ANDROID", MAU: 120, Requests: 1200},
				{Yearmonth: "202602", EnvironmentID: "env1", EnvironmentName: "Env1", ProjectName: "Proj1", SourceID: "GO_SERVER", MAU: 60, Requests: 600},
				{Yearmonth: "202602", EnvironmentID: "env2", EnvironmentName: "Env2", ProjectName: "Proj2", SourceID: "ANDROID", MAU: 2, Requests: 20},
			},
			expected: []*insightsproto.MonthlySummarySeries{
				{
					EnvironmentId:   "env1",
					EnvironmentName: "Env1",
					ProjectName:     "Proj1",
					SourceId:        clientproto.SourceId_ANDROID,
					Data: []*insightsproto.MonthlySummaryDataPoint{
						{Yearmonth: "202601", Mau: 100, Requests: 1000},
						{Yearmonth: "202602", Mau: 120, Requests: 1200},
					},
				},
				{
					EnvironmentId:   "env1",
					EnvironmentName: "Env1",
					ProjectName:     "Proj1",
					SourceId:        clientproto.SourceId_GO_SERVER,
					Data: []*insightsproto.MonthlySummaryDataPoint{
						{Yearmonth: "202601", Mau: 50, Requests: 500},
						{Yearmonth: "202602", Mau: 60, Requests: 600},
					},
				},
				{
					EnvironmentId:   "env2",
					EnvironmentName: "Env2",
					ProjectName:     "Proj2",
					SourceId:        clientproto.SourceId_ANDROID,
					Data: []*insightsproto.MonthlySummaryDataPoint{
						{Yearmonth: "202602", Mau: 2, Requests: 20},
					},
				},
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			s := &insightsService{}
			actual := s.convertToMonthlySummarySeries(p.input)
			assert.ElementsMatch(t, p.expected, actual)
		})
	}
}

func createContextWithToken(t *testing.T) context.Context {
	t.Helper()
	tk := &token.AccessToken{
		Email:         "test@example.com",
		IsSystemAdmin: false,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.AccessTokenKey, tk)
}

func newInsightsServiceForTest(
	t *testing.T,
	mockController *gomock.Controller,
	specifiedOrgRole *accountproto.AccountV2_Role_Organization,
	specifiedEnvRole *accountproto.AccountV2_Role_Environment,
) *insightsService {
	t.Helper()

	var or accountproto.AccountV2_Role_Organization
	var er accountproto.AccountV2_Role_Environment
	if specifiedOrgRole != nil {
		or = *specifiedOrgRole
	} else {
		or = accountproto.AccountV2_Role_Organization_ADMIN
	}
	if specifiedEnvRole != nil {
		er = *specifiedEnvRole
	} else {
		er = accountproto.AccountV2_Role_Environment_EDITOR
	}

	logger, err := log.NewLogger()
	require.NoError(t, err)

	accountClientMock := accountclientmock.NewMockClient(mockController)
	ar := &accountproto.GetAccountV2ByEnvironmentIDResponse{
		Account: &accountproto.AccountV2{
			Email:            "email",
			OrganizationRole: or,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "env1",
					Role:          er,
				},
			},
		},
	}
	accountClientMock.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()

	promClientMock := prometheusmock.NewMockClient(mockController)
	storageMock := storagemock.NewMockMonthlySummaryStorage(mockController)

	return &insightsService{
		accountClient:         accountClientMock,
		promClient:            promClientMock,
		monthlySummaryStorage: storageMock,
		logger:                logger,
	}
}

func newInsightsServiceWithoutProm(
	t *testing.T,
	mockController *gomock.Controller,
	specifiedOrgRole *accountproto.AccountV2_Role_Organization,
	specifiedEnvRole *accountproto.AccountV2_Role_Environment,
) *insightsService {
	t.Helper()

	var or accountproto.AccountV2_Role_Organization
	var er accountproto.AccountV2_Role_Environment
	if specifiedOrgRole != nil {
		or = *specifiedOrgRole
	} else {
		or = accountproto.AccountV2_Role_Organization_ADMIN
	}
	if specifiedEnvRole != nil {
		er = *specifiedEnvRole
	} else {
		er = accountproto.AccountV2_Role_Environment_EDITOR
	}

	logger, err := log.NewLogger()
	require.NoError(t, err)

	accountClientMock := accountclientmock.NewMockClient(mockController)
	ar := &accountproto.GetAccountV2ByEnvironmentIDResponse{
		Account: &accountproto.AccountV2{
			Email:            "email",
			OrganizationRole: or,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "env1",
					Role:          er,
				},
			},
		},
	}
	accountClientMock.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()

	storageMock := storagemock.NewMockMonthlySummaryStorage(mockController)

	return &insightsService{
		accountClient:         accountClientMock,
		promClient:            nil, // prometheus not configured
		monthlySummaryStorage: storageMock,
		logger:                logger,
	}
}

func toPtr[T any](value T) *T {
	return &value
}
