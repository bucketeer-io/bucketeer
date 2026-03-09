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
	"fmt"
	"sort"
	"time"

	"github.com/prometheus/common/model"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	insightsstorage "github.com/bucketeer-io/bucketeer/v2/pkg/insights/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/prometheus"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	clientproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	insightsproto "github.com/bucketeer-io/bucketeer/v2/proto/insights"
)

const (
	defaultQueryStep  = 5 * time.Minute
	maxQueryRangeDays = 30
)

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type insightsService struct {
	accountClient         accountclient.Client
	promClient            prometheus.Client // nil when prometheus-url is not configured
	monthlySummaryStorage insightsstorage.MonthlySummaryStorage
	logger                *zap.Logger
}

func NewInsightsService(
	accountClient accountclient.Client,
	promClient prometheus.Client,
	monthlySummaryStorage insightsstorage.MonthlySummaryStorage,
	opts ...Option,
) rpc.Service {
	dopts := &options{logger: zap.NewNop()}
	for _, opt := range opts {
		opt(dopts)
	}
	return &insightsService{
		accountClient:         accountClient,
		promClient:            promClient,
		monthlySummaryStorage: monthlySummaryStorage,
		logger:                dopts.logger.Named("insights-api"),
	}
}

func (s *insightsService) Register(server *grpc.Server) {
	insightsproto.RegisterInsightsServiceServer(server, s)
}

func (s *insightsService) GetInsightsMonthlySummary(
	ctx context.Context,
	req *insightsproto.GetInsightsMonthlySummaryRequest,
) (*insightsproto.GetInsightsMonthlySummaryResponse, error) {
	if err := s.validateMonthlySummaryRequest(req); err != nil {
		return nil, err
	}
	if err := s.checkEnvironmentRoles(ctx, req.EnvironmentIds); err != nil {
		return nil, err
	}

	sourceIDs := sourceIDsToStrings(req.SourceIds)
	if len(sourceIDs) == 0 {
		sourceIDs = allSourceIDStrings()
	}

	records, err := s.monthlySummaryStorage.ListMonthlySummaries(
		ctx,
		req.EnvironmentIds,
		sourceIDs,
	)
	if err != nil {
		s.logger.Error("Failed to list monthly summary",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, statusInternal.Err()
	}

	series := s.convertToMonthlySummarySeries(records)
	return &insightsproto.GetInsightsMonthlySummaryResponse{
		Series: series,
	}, nil
}

func (s *insightsService) GetInsightsLatency(
	ctx context.Context,
	req *insightsproto.GetInsightsTimeSeriesRequest,
) (*insightsproto.GetInsightsTimeSeriesResponse, error) {
	query := latencyQuery(
		req.EnvironmentIds,
		sourceIDsToStrings(req.SourceIds),
		apiIDsToMethods(req.ApiIds),
	)
	return s.queryTimeSeries(ctx, req, query, nil)
}

func (s *insightsService) GetInsightsRequests(
	ctx context.Context,
	req *insightsproto.GetInsightsTimeSeriesRequest,
) (*insightsproto.GetInsightsTimeSeriesResponse, error) {
	query := requestCountQuery(
		req.EnvironmentIds,
		sourceIDsToStrings(req.SourceIds),
		apiIDsToMethods(req.ApiIds),
	)
	return s.queryTimeSeries(ctx, req, query, nil)
}

func (s *insightsService) GetInsightsEvaluations(
	ctx context.Context,
	req *insightsproto.GetInsightsTimeSeriesRequest,
) (*insightsproto.GetInsightsTimeSeriesResponse, error) {
	query := evaluationsQuery(
		req.EnvironmentIds,
		sourceIDsToStrings(req.SourceIds),
	)
	return s.queryTimeSeries(ctx, req, query, []string{"evaluation_type"})
}

func (s *insightsService) GetInsightsErrorRates(
	ctx context.Context,
	req *insightsproto.GetInsightsTimeSeriesRequest,
) (*insightsproto.GetInsightsTimeSeriesResponse, error) {
	query := errorRatesQuery(
		req.EnvironmentIds,
		sourceIDsToStrings(req.SourceIds),
		apiIDsToMethods(req.ApiIds),
	)
	return s.queryTimeSeries(ctx, req, query, nil)
}

func (s *insightsService) queryTimeSeries(
	ctx context.Context,
	req *insightsproto.GetInsightsTimeSeriesRequest,
	query string,
	labelKeys []string,
) (*insightsproto.GetInsightsTimeSeriesResponse, error) {
	if s.promClient == nil {
		return nil, statusDataSourceNotConfigured.Err()
	}
	if err := s.validateTimeSeriesRequest(req); err != nil {
		return nil, err
	}
	if err := s.checkEnvironmentRoles(ctx, req.EnvironmentIds); err != nil {
		return nil, err
	}

	matrix, err := s.promClient.QueryRange(
		ctx,
		query,
		time.Unix(req.StartAt, 0),
		time.Unix(req.EndAt, 0),
		defaultQueryStep,
	)
	if err != nil {
		s.logger.Error("Failed to query prometheus",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, statusInternal.Err()
	}

	return &insightsproto.GetInsightsTimeSeriesResponse{
		Timeseries: convertMatrixToTimeSeries(matrix, labelKeys),
	}, nil
}

func (s *insightsService) convertToMonthlySummarySeries(
	rows []insightsstorage.ListMonthlySummaryResult,
) []*insightsproto.MonthlySummarySeries {
	seriesMap := make(map[string]*insightsproto.MonthlySummarySeries)

	for _, row := range rows {
		envSrc := fmt.Sprintf("%s:%s", row.EnvironmentID, row.SourceID)

		if _, ok := seriesMap[envSrc]; !ok {
			seriesMap[envSrc] = &insightsproto.MonthlySummarySeries{
				EnvironmentId:   row.EnvironmentID,
				EnvironmentName: row.EnvironmentName,
				ProjectName:     row.ProjectName,
				SourceId:        clientproto.SourceId(clientproto.SourceId_value[row.SourceID]),
			}
		}

		seriesMap[envSrc].Data = append(seriesMap[envSrc].Data, &insightsproto.MonthlySummaryDataPoint{
			Yearmonth: row.Yearmonth,
			Mau:       row.MAU,
			Requests:  row.Requests,
		})
	}

	series := make([]*insightsproto.MonthlySummarySeries, 0, len(seriesMap))
	for _, s := range seriesMap {
		series = append(series, s)
	}
	return series
}

func convertMatrixToTimeSeries(matrix model.Matrix, labelKeys []string) []*insightsproto.InsightsTimeSeries {
	result := make([]*insightsproto.InsightsTimeSeries, 0, len(matrix))
	for _, stream := range matrix {
		envID := string(stream.Metric[model.LabelName("environment_id")])
		sourceIDStr := string(stream.Metric[model.LabelName("source_id")])
		method := string(stream.Metric[model.LabelName("method")])

		sourceID := clientproto.SourceId(clientproto.SourceId_value[sourceIDStr])
		apiID := methodToApiID[method]

		dataPoints := make([]*insightsproto.InsightsDataPoint, 0, len(stream.Values))
		for _, v := range stream.Values {
			dataPoints = append(dataPoints, &insightsproto.InsightsDataPoint{
				Timestamp: int64(v.Timestamp) / 1000,
				Value:     float64(v.Value),
			})
		}

		result = append(result, &insightsproto.InsightsTimeSeries{
			EnvironmentId: envID,
			SourceId:      sourceID,
			ApiId:         apiID,
			Data:          dataPoints,
			Labels:        extractLabels(stream.Metric, labelKeys),
		})
	}
	return result
}

func extractLabels(metric model.Metric, keys []string) map[string]string {
	if len(keys) == 0 {
		return nil
	}
	labels := make(map[string]string, len(keys))
	for _, key := range keys {
		labels[key] = string(metric[model.LabelName(key)])
	}
	return labels
}

func (s *insightsService) validateMonthlySummaryRequest(
	req *insightsproto.GetInsightsMonthlySummaryRequest,
) error {
	if len(req.EnvironmentIds) == 0 {
		return statusEnvironmentIDRequired.Err()
	}
	return nil
}

func (s *insightsService) validateTimeSeriesRequest(
	req *insightsproto.GetInsightsTimeSeriesRequest,
) error {
	if len(req.EnvironmentIds) == 0 {
		return statusEnvironmentIDRequired.Err()
	}
	if req.StartAt == 0 {
		return statusStartAtRequired.Err()
	}
	if req.EndAt == 0 {
		return statusEndAtRequired.Err()
	}
	if req.StartAt > req.EndAt {
		return statusStartAtIsAfterEndAt.Err()
	}
	maxRange := int64(maxQueryRangeDays * 24 * 60 * 60)
	if req.EndAt-req.StartAt > maxRange {
		return statusQueryRangeTooLarge.Err()
	}
	return nil
}

func (s *insightsService) checkEnvironmentRoles(ctx context.Context, envIDs []string) error {
	for _, envID := range envIDs {
		_, err := role.CheckEnvironmentRole(
			ctx,
			accountproto.AccountV2_Role_Environment_VIEWER,
			envID,
			func(email string) (*accountproto.AccountV2, error) {
				resp, err := s.accountClient.GetAccountV2ByEnvironmentID(
					ctx,
					&accountproto.GetAccountV2ByEnvironmentIDRequest{
						Email:         email,
						EnvironmentId: envID,
					},
				)
				if err != nil {
					return nil, err
				}
				return resp.Account, nil
			},
		)
		if err != nil {
			switch status.Code(err) {
			case codes.Unauthenticated:
				s.logger.Info("Unauthenticated",
					log.FieldsFromIncomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentId", envID),
					)...,
				)
				return statusUnauthenticated.Err()
			case codes.PermissionDenied:
				s.logger.Info("Permission denied",
					log.FieldsFromIncomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentId", envID),
					)...,
				)
				return statusPermissionDenied.Err()
			default:
				s.logger.Error("Failed to check role",
					log.FieldsFromIncomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentId", envID),
					)...,
				)
				return statusInternal.Err()
			}
		}
	}
	return nil
}

func sourceIDsToStrings(sourceIDs []clientproto.SourceId) []string {
	strs := make([]string, len(sourceIDs))
	for i, sid := range sourceIDs {
		strs[i] = sid.String()
	}
	return strs
}

func allSourceIDStrings() []string {
	ids := make([]int, 0, len(clientproto.SourceId_name))
	for id := range clientproto.SourceId_name {
		if id == int32(clientproto.SourceId_UNKNOWN) {
			continue
		}
		ids = append(ids, int(id))
	}
	sort.Ints(ids)

	sourceIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		sourceIDs = append(sourceIDs, clientproto.SourceId_name[int32(id)])
	}
	return sourceIDs
}

var apiIDToMethod = map[clientproto.ApiId]string{
	clientproto.ApiId_GET_EVALUATION:    "GetEvaluation",
	clientproto.ApiId_GET_EVALUATIONS:   "GetEvaluations",
	clientproto.ApiId_REGISTER_EVENTS:   "RegisterEvents",
	clientproto.ApiId_GET_FEATURE_FLAGS: "GetFeatureFlags",
	clientproto.ApiId_GET_SEGMENT_USERS: "GetSegmentUsers",
}

var methodToApiID = map[string]clientproto.ApiId{
	"GetEvaluation":   clientproto.ApiId_GET_EVALUATION,
	"GetEvaluations":  clientproto.ApiId_GET_EVALUATIONS,
	"RegisterEvents":  clientproto.ApiId_REGISTER_EVENTS,
	"GetFeatureFlags": clientproto.ApiId_GET_FEATURE_FLAGS,
	"GetSegmentUsers": clientproto.ApiId_GET_SEGMENT_USERS,
}

func apiIDsToMethods(apiIDs []clientproto.ApiId) []string {
	methods := make([]string, len(apiIDs))
	for i, aid := range apiIDs {
		if name, ok := apiIDToMethod[aid]; ok {
			methods[i] = name
		} else {
			methods[i] = aid.String()
		}
	}
	return methods
}
