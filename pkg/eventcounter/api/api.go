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
	"fmt"
	"sort"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	ecdruid "github.com/bucketeer-io/bucketeer/pkg/eventcounter/druid"
	v2ecstorage "github.com/bucketeer-io/bucketeer/pkg/eventcounter/storage/v2"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/role"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	bqquerier "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/querier"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const listRequestPageSize = 500

const (
	eventCountPrefix   = "ec"
	userCountPrefix    = "uc"
	defaultVariationID = "default"
)

var (
	jpLocation = time.FixedZone("Asia/Tokyo", 9*60*60)
)

type eventCounterService struct {
	experimentClient             experimentclient.Client
	featureClient                featureclient.Client
	accountClient                accountclient.Client
	druidQuerier                 ecdruid.Querier
	eventStorage                 v2ecstorage.EventStorage
	mysqlExperimentResultStorage v2ecstorage.ExperimentResultStorage
	userCountStorage             v2ecstorage.UserCountStorage
	metrics                      metrics.Registerer
	evaluationCountCacher        cachev3.EventCounterCache
	logger                       *zap.Logger
}

func NewEventCounterService(
	mc mysql.Client,
	e experimentclient.Client,
	f featureclient.Client,
	a accountclient.Client,
	d ecdruid.Querier,
	b bqquerier.Client,
	bigQueryDataSet string,
	r metrics.Registerer,
	redis cache.MultiGetDeleteCountCache,
	l *zap.Logger,
) rpc.Service {
	registerMetrics(r)
	return &eventCounterService{
		experimentClient:             e,
		featureClient:                f,
		accountClient:                a,
		druidQuerier:                 d,
		eventStorage:                 v2ecstorage.NewEventStorage(b, bigQueryDataSet, l),
		mysqlExperimentResultStorage: v2ecstorage.NewExperimentResultStorage(mc),
		userCountStorage:             v2ecstorage.NewUserCountStorage(mc),
		metrics:                      r,
		evaluationCountCacher:        cachev3.NewEventCountCache(redis),
		logger:                       l.Named("api"),
	}
}

func (s *eventCounterService) Register(server *grpc.Server) {
	ecproto.RegisterEventCounterServiceServer(server, s)
}

func (s *eventCounterService) GetExperimentEvaluationCount(
	ctx context.Context,
	req *ecproto.GetExperimentEvaluationCountRequest,
) (*ecproto.GetExperimentEvaluationCountResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err = validateGetExperimentEvaluationCountRequest(req, localizer); err != nil {
		return nil, err
	}
	startAt := time.Unix(req.StartAt, 0)
	endAt := time.Unix(req.EndAt, 0)
	evaluationCounts, err := s.eventStorage.QueryEvaluationCount(
		ctx,
		req.EnvironmentNamespace,
		startAt,
		endAt,
		req.FeatureId,
		req.FeatureVersion,
	)
	if err != nil {
		s.logger.Error(
			"Failed to query experiment evaluation counts",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.Time("startAt", startAt),
				zap.Time("endAt", endAt),
				zap.String("featureId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	variationCounts := s.convertEvaluationCounts(evaluationCounts, req.VariationIds)
	s.logger.Debug("GetExperimentEvaluationCount result", zap.Any("rows", variationCounts))
	return &ecproto.GetExperimentEvaluationCountResponse{
		FeatureId:       req.FeatureId,
		FeatureVersion:  req.FeatureVersion,
		VariationCounts: variationCounts,
	}, nil
}

func validateGetExperimentEvaluationCountRequest(
	req *ecproto.GetExperimentEvaluationCountRequest,
	localizer locale.Localizer,
) error {
	if req.StartAt == 0 {
		dt, err := statusStartAtRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "start_at"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.EndAt == 0 {
		dt, err := statusEndAtRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "end_at"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.StartAt > req.EndAt {
		dt, err := statusStartAtIsAfterEndAt.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.StartAtIsAfterEnd),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.FeatureId == "" {
		dt, err := statusFeatureIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *eventCounterService) convertEvaluationCounts(
	rows []*v2ecstorage.EvaluationEventCount,
	variationIDs []string,
) []*ecproto.VariationCount {
	vcsMap := map[string]*ecproto.VariationCount{}
	for _, id := range variationIDs {
		vcsMap[id] = &ecproto.VariationCount{VariationId: id}
	}
	for _, row := range rows {
		vc, ok := vcsMap[row.VariationID]
		if !ok {
			continue
		}
		vc.UserCount = row.EvaluationUser
		vc.EventCount = row.EvaluationTotal
		vcsMap[row.VariationID] = vc
	}
	vcs := make([]*ecproto.VariationCount, 0, len(vcsMap))
	for _, vc := range vcsMap {
		vcs = append(vcs, vc)
	}
	sort.SliceStable(vcs, func(i, j int) bool { return vcs[i].VariationId < vcs[j].VariationId })
	return vcs
}

func (s *eventCounterService) GetEvaluationCountV2(
	ctx context.Context,
	req *ecproto.GetEvaluationCountV2Request,
) (*ecproto.GetEvaluationCountV2Response, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err = validateGetEvaluationCountV2Request(req); err != nil {
		return nil, err
	}
	startAt := time.Unix(req.StartAt, 0)
	endAt := time.Unix(req.EndAt, 0)
	headers, rows, err := s.druidQuerier.QueryEvaluationCount(
		ctx,
		req.EnvironmentNamespace,
		startAt,
		endAt,
		req.FeatureId,
		req.FeatureVersion,
		"",
		[]string{}, []*ecproto.Filter{},
	)
	if err != nil {
		s.logger.Error(
			"Failed to query evaluation counts",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.Time("startAt", startAt),
				zap.Time("endAt", endAt),
				zap.String("featureId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	vcs, err := convToVariationCounts(headers, rows, req.VariationIds)
	if err != nil {
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &ecproto.GetEvaluationCountV2Response{
		Count: &ecproto.EvaluationCount{
			FeatureId:      req.FeatureId,
			FeatureVersion: req.FeatureVersion,
			RealtimeCounts: vcs,
		},
	}, nil
}

func validateGetEvaluationCountV2Request(req *ecproto.GetEvaluationCountV2Request) error {
	if req.StartAt == 0 {
		return localizedError(statusStartAtRequired, locale.JaJP)
	}
	if req.EndAt == 0 {
		return localizedError(statusEndAtRequired, locale.JaJP)
	}
	if req.StartAt > req.EndAt {
		return localizedError(statusStartAtIsAfterEndAt, locale.JaJP)
	}
	if req.FeatureId == "" {
		return localizedError(statusFeatureIDRequired, locale.JaJP)
	}
	return nil
}

func convToVariationCounts(
	headers *ecproto.Row,
	rows []*ecproto.Row,
	variationIDs []string,
) ([]*ecproto.VariationCount, error) {
	vcsMap := map[string]*ecproto.VariationCount{}
	for _, id := range variationIDs {
		vcsMap[id] = &ecproto.VariationCount{VariationId: id}
	}
	varIdx, err := variationIdx(headers)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		vid := row.Cells[varIdx].Value
		vc, ok := vcsMap[vid]
		if !ok {
			continue
		}
		for i, cell := range row.Cells {
			switch headers.Cells[i].Value {
			// Evaluation.
			case ecdruid.ColumnEvaluationTotal:
				vc.EventCount = int64(cell.ValueDouble)
			case ecdruid.ColumnEvaluationUser:
				vc.UserCount = int64(cell.ValueDouble)
			// Goal.
			case ecdruid.ColumnGoalTotal:
				vc.EventCount = int64(cell.ValueDouble)
			case ecdruid.ColumnGoalUser:
				vc.UserCount = int64(cell.ValueDouble)
			case ecdruid.ColumnGoalValueTotal:
				vc.ValueSum = cell.ValueDouble
			case ecdruid.ColumnGoalValueMean:
				vc.ValueSumPerUserMean = cell.ValueDouble
			case ecdruid.ColumnGoalValueVariance:
				vc.ValueSumPerUserVariance = cell.ValueDouble
			}
		}
		vcsMap[vid] = vc
	}
	vcs := []*ecproto.VariationCount{}
	for _, vc := range vcsMap {
		vcs = append(vcs, vc)
	}
	sort.SliceStable(vcs, func(i, j int) bool { return vcs[i].VariationId < vcs[j].VariationId })
	return vcs, nil
}

func variationIdx(headers *ecproto.Row) (int, error) {
	for i, cell := range headers.Cells {
		if cell.Value == ecdruid.ColumnVariation {
			return i, nil
		}
	}
	return 0, errors.New("eventcounter: variation header not found")
}

func (s *eventCounterService) GetEvaluationTimeseriesCount(
	ctx context.Context,
	req *ecproto.GetEvaluationTimeseriesCountRequest,
) (*ecproto.GetEvaluationTimeseriesCountResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if req.FeatureId == "" {
		return nil, localizedError(statusFeatureIDRequired, locale.JaJP)
	}
	resp, err := s.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		EnvironmentNamespace: req.EnvironmentNamespace,
		Id:                   req.FeatureId,
	})
	if err != nil {
		s.logger.Error(
			"Failed to get feature",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.String("featureId", req.FeatureId),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	endAt := time.Now()
	startAt, err := genInterval(jpLocation, endAt, 30)
	if err != nil {
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	variationTSEvents := []*ecproto.VariationTimeseries{}
	variationTSUsers := []*ecproto.VariationTimeseries{}
	for _, variation := range resp.Feature.Variations {
		varTS, err := s.druidQuerier.QueryEvaluationTimeseriesCount(
			ctx,
			req.EnvironmentNamespace,
			startAt,
			endAt,
			req.FeatureId,
			0,
			variation.Id,
		)
		if err != nil {
			s.logger.Error(
				"Failed to query goal counts",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
					zap.Time("startAt", startAt),
					zap.Time("endAt", endAt),
					zap.String("featureId", req.FeatureId),
					zap.Int32("featureVersion", resp.Feature.Version),
					zap.String("variationId", variation.Id),
				)...,
			)
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		variationTSEvents = append(variationTSEvents, varTS[ecdruid.ColumnEvaluationTotal])
		variationTSUsers = append(variationTSUsers, varTS[ecdruid.ColumnEvaluationUser])
	}
	return &ecproto.GetEvaluationTimeseriesCountResponse{
		EventCounts: variationTSEvents,
		UserCounts:  variationTSUsers,
	}, nil
}

func (s *eventCounterService) GetEvaluationTimeseriesCountV2(
	ctx context.Context,
	req *ecproto.GetEvaluationTimeseriesCountRequest,
) (*ecproto.GetEvaluationTimeseriesCountResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if req.FeatureId == "" {
		return nil, localizedError(statusFeatureIDRequired, locale.JaJP)
	}
	resp, err := s.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		EnvironmentNamespace: req.EnvironmentNamespace,
		Id:                   req.FeatureId,
	})
	if err != nil {
		s.logger.Error(
			"Failed to get feature",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.String("featureId", req.FeatureId),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	endAt := time.Now()
	startAt, err := genInterval(jpLocation, endAt, 30)
	if err != nil {
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	timeStamps := getOneMonthTimeStamps(startAt)
	vIDs := getVariationIDs(resp.Feature.Variations)
	variationTSEvents := []*ecproto.VariationTimeseries{}
	variationTSUsers := []*ecproto.VariationTimeseries{}
	for _, vID := range vIDs {
		eventCountKeys := []string{}
		userCountKeys := []string{}
		for _, ts := range timeStamps {
			ec := newEvaluationCountkey(eventCountPrefix, req.FeatureId, vID, req.EnvironmentNamespace, ts)
			eventCountKeys = append(eventCountKeys, ec)
			uc := newEvaluationCountkey(userCountPrefix, req.FeatureId, vID, req.EnvironmentNamespace, ts)
			userCountKeys = append(userCountKeys, uc)
		}
		eventCounts, err := s.evaluationCountCacher.GetEventCounts(eventCountKeys)
		if err != nil {
			s.logger.Error(
				"Failed to get event counts",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
					zap.Time("startAt", startAt),
					zap.Time("endAt", endAt),
					zap.String("featureId", req.FeatureId),
					zap.Int32("featureVersion", resp.Feature.Version),
					zap.String("variationId", vID),
				)...,
			)
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		userCounts, err := s.evaluationCountCacher.GetUserCounts(userCountKeys)
		if err != nil {
			s.logger.Error(
				"Failed to get user counts",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
					zap.Time("startAt", startAt),
					zap.Time("endAt", endAt),
					zap.String("featureId", req.FeatureId),
					zap.Int32("featureVersion", resp.Feature.Version),
					zap.String("variationId", vID),
				)...,
			)
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		variationTSUsers = append(variationTSUsers, &ecproto.VariationTimeseries{
			VariationId: vID,
			Timeseries: &ecproto.Timeseries{
				Timestamps: timeStamps,
				Values:     userCounts,
			},
		})
		variationTSEvents = append(variationTSEvents, &ecproto.VariationTimeseries{
			VariationId: vID,
			Timeseries: &ecproto.Timeseries{
				Timestamps: timeStamps,
				Values:     eventCounts,
			},
		})
	}
	return &ecproto.GetEvaluationTimeseriesCountResponse{
		EventCounts: variationTSEvents,
		UserCounts:  variationTSUsers,
	}, nil
}

func genInterval(loc *time.Location, endAt time.Time, durationDays int) (time.Time, error) {
	year, month, day := endAt.In(loc).AddDate(0, 0, -durationDays).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, loc), nil
}

func newEvaluationCountkey(
	kind, featureID, variationID, environmentNamespace string,
	ts int64,
) string {
	return cache.MakeKey(
		kind,
		fmt.Sprintf("%d:%s:%s", ts, featureID, variationID),
		environmentNamespace,
	)
}

func getOneMonthTimeStamps(startAt time.Time) []int64 {
	limit := 31
	timeStamps := make([]int64, 0, limit)
	for i := 0; i < limit; i++ {
		ts := startAt.AddDate(0, 0, i).Unix()
		timeStamps = append(timeStamps, ts)
	}
	return timeStamps
}

func getVariationIDs(vs []*featureproto.Variation) []string {
	vIDs := []string{}
	for _, v := range vs {
		vIDs = append(vIDs, v.Id)
	}
	vIDs = append(vIDs, defaultVariationID)
	return vIDs
}

func (s *eventCounterService) GetExperimentResult(
	ctx context.Context,
	req *ecproto.GetExperimentResultRequest,
) (*ecproto.GetExperimentResultResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if req.ExperimentId == "" {
		return nil, localizedError(statusExperimentIDRequired, locale.JaJP)
	}
	result, err := s.mysqlExperimentResultStorage.GetExperimentResult(ctx, req.ExperimentId, req.EnvironmentNamespace)
	if err != nil {
		if err == v2ecstorage.ErrExperimentResultNotFound {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to get experiment result",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.String("experimentId", req.ExperimentId),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &ecproto.GetExperimentResultResponse{
		ExperimentResult: result.ExperimentResult,
	}, nil
}

func (s *eventCounterService) ListExperimentResults(
	ctx context.Context,
	req *ecproto.ListExperimentResultsRequest,
) (*ecproto.ListExperimentResultsResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if req.FeatureId == "" {
		return nil, localizedError(statusFeatureIDRequired, locale.JaJP)
	}
	experiments, err := s.listExperiments(ctx, req.FeatureId, req.FeatureVersion, req.EnvironmentNamespace)
	if err != nil {
		if err == storage.ErrKeyNotFound {
			listExperimentCountsCounter.WithLabelValues(codeSuccess).Inc()
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to get Experiment list",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.String("featureID", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion.Value),
			)...,
		)
		listExperimentCountsCounter.WithLabelValues(codeFail).Inc()
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	results := make(map[string]*ecproto.ExperimentResult, len(experiments))
	for _, e := range experiments {
		er, err := s.getExperimentResultMySQL(ctx, e.Id, req.EnvironmentNamespace)
		if err != nil {
			if err == v2ecstorage.ErrExperimentResultNotFound {
				getExperimentCountsCounter.WithLabelValues(codeSuccess).Inc()
			} else {
				s.logger.Error(
					"Failed to get Experiment result",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentNamespace", req.EnvironmentNamespace),
						zap.String("experimentID", e.Id),
					)...,
				)
				getExperimentCountsCounter.WithLabelValues(codeFail).Inc()
			}
			continue
		}
		getExperimentCountsCounter.WithLabelValues(codeSuccess).Inc()
		results[e.Id] = er
	}
	listExperimentCountsCounter.WithLabelValues(codeSuccess).Inc()
	return &ecproto.ListExperimentResultsResponse{Results: results}, nil
}

func (s *eventCounterService) listExperiments(
	ctx context.Context,
	featureID string,
	featureVersion *wrappers.Int32Value,
	environmentNamespace string,
) ([]*experimentproto.Experiment, error) {
	experiments := []*experimentproto.Experiment{}
	cursor := ""
	for {
		resp, err := s.experimentClient.ListExperiments(ctx, &experimentproto.ListExperimentsRequest{
			FeatureId:            featureID,
			FeatureVersion:       featureVersion,
			PageSize:             listRequestPageSize,
			Cursor:               cursor,
			EnvironmentNamespace: environmentNamespace,
		})
		if err != nil {
			return nil, err
		}
		experiments = append(experiments, resp.Experiments...)
		featureSize := len(resp.Experiments)
		if featureSize == 0 || featureSize < listRequestPageSize {
			return experiments, nil
		}
		cursor = resp.Cursor
	}
}

func (s *eventCounterService) GetGoalCount(
	ctx context.Context,
	req *ecproto.GetGoalCountRequest,
) (*ecproto.GetGoalCountResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetGoalCountsRequest(req); err != nil {
		return nil, err
	}
	startAt := time.Unix(req.StartAt, 0)
	endAt := time.Unix(req.EndAt, 0)
	headers, rows, err := s.druidQuerier.QueryCount(
		ctx,
		req.EnvironmentNamespace,
		startAt,
		endAt,
		req.GoalId,
		req.FeatureId,
		req.FeatureVersion,
		req.Reason,
		req.Segments,
		req.Filters,
	)
	if err != nil {
		s.logger.Error(
			"Failed to query goal counts",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.Time("startAt", startAt),
				zap.Time("endAt", endAt),
				zap.String("featureId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
				zap.Strings("segments", req.Segments),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &ecproto.GetGoalCountResponse{Headers: headers, Rows: rows}, nil
}

func validateGetGoalCountsRequest(req *ecproto.GetGoalCountRequest) error {
	if req.StartAt == 0 {
		return localizedError(statusStartAtRequired, locale.JaJP)
	}
	if req.EndAt == 0 {
		return localizedError(statusEndAtRequired, locale.JaJP)
	}
	if req.StartAt > req.EndAt {
		return localizedError(statusStartAtIsAfterEndAt, locale.JaJP)
	}
	if req.StartAt < time.Now().Add(-31*24*time.Hour).Unix() {
		return localizedError(statusPeriodOutOfRange, locale.JaJP)
	}
	if req.GoalId == "" {
		return localizedError(statusGoalIDRequired, locale.JaJP)
	}
	return nil
}

func (s *eventCounterService) GetExperimentGoalCount(
	ctx context.Context,
	req *ecproto.GetExperimentGoalCountRequest,
) (*ecproto.GetExperimentGoalCountResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err = validateGetExperimentGoalCountRequest(req, localizer); err != nil {
		return nil, err
	}
	startAt := time.Unix(req.StartAt, 0)
	endAt := time.Unix(req.EndAt, 0)
	goalCounts, err := s.eventStorage.QueryGoalCount(
		ctx,
		req.EnvironmentNamespace,
		startAt,
		endAt,
		req.FeatureId,
		req.FeatureVersion,
	)
	if err != nil {
		s.logger.Error(
			"Failed to query experiment goal counts",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.Time("startAt", startAt),
				zap.Time("endAt", endAt),
				zap.String("featureId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	variationCounts := s.convertGoalCounts(goalCounts, req.VariationIds)
	s.logger.Debug("GetExperimentGoalCount result", zap.Any("rows", variationCounts))
	return &ecproto.GetExperimentGoalCountResponse{
		GoalId:          req.GoalId,
		VariationCounts: variationCounts,
	}, nil
}

func validateGetExperimentGoalCountRequest(
	req *ecproto.GetExperimentGoalCountRequest,
	localizer locale.Localizer,
) error {
	if req.StartAt == 0 {
		dt, err := statusStartAtRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "start_at"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.EndAt == 0 {
		dt, err := statusEndAtRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "end_at"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.StartAt > req.EndAt {
		dt, err := statusStartAtIsAfterEndAt.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.StartAtIsAfterEnd),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.FeatureId == "" {
		dt, err := statusFeatureIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.GoalId == "" {
		dt, err := statusGoalIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *eventCounterService) convertGoalCounts(
	rows []*v2ecstorage.GoalEventCount,
	variationIDs []string,
) []*ecproto.VariationCount {
	vcsMap := map[string]*ecproto.VariationCount{}
	for _, id := range variationIDs {
		vcsMap[id] = &ecproto.VariationCount{VariationId: id}
	}
	for _, row := range rows {
		vc, ok := vcsMap[row.VariationID]
		if !ok {
			continue
		}
		vc.UserCount = row.GoalUser
		vc.EventCount = row.GoalTotal
		vc.ValueSum = row.GoalValueTotal
		vc.ValueSumPerUserMean = row.GoalValueMean
		vc.ValueSumPerUserVariance = row.GoalValueVariance
		vcsMap[row.VariationID] = vc
	}
	vcs := make([]*ecproto.VariationCount, 0, len(vcsMap))
	for _, vc := range vcsMap {
		vcs = append(vcs, vc)
	}
	sort.SliceStable(vcs, func(i, j int) bool { return vcs[i].VariationId < vcs[j].VariationId })
	return vcs
}

func (s *eventCounterService) GetGoalCountV2(
	ctx context.Context,
	req *ecproto.GetGoalCountV2Request,
) (*ecproto.GetGoalCountV2Response, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err = validateGetGoalCountV2Request(req); err != nil {
		return nil, err
	}
	startAt := time.Unix(req.StartAt, 0)
	endAt := time.Unix(req.EndAt, 0)
	headers, rows, err := s.druidQuerier.QueryGoalCount(
		ctx,
		req.EnvironmentNamespace,
		startAt,
		endAt,
		req.GoalId,
		req.FeatureId,
		req.FeatureVersion,
		"",
		[]string{}, []*ecproto.Filter{},
	)
	if err != nil {
		s.logger.Error(
			"Failed to query goal counts",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.Time("startAt", startAt),
				zap.Time("endAt", endAt),
				zap.String("featureId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	vcs, err := convToVariationCounts(headers, rows, req.VariationIds)
	if err != nil {
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &ecproto.GetGoalCountV2Response{
		GoalCounts: &ecproto.GoalCounts{
			GoalId:         req.GoalId,
			RealtimeCounts: vcs,
		},
	}, nil
}

func validateGetGoalCountV2Request(req *ecproto.GetGoalCountV2Request) error {
	if req.StartAt == 0 {
		return localizedError(statusStartAtRequired, locale.JaJP)
	}
	if req.EndAt == 0 {
		return localizedError(statusEndAtRequired, locale.JaJP)
	}
	if req.StartAt > req.EndAt {
		return localizedError(statusStartAtIsAfterEndAt, locale.JaJP)
	}
	if req.GoalId == "" {
		return localizedError(statusGoalIDRequired, locale.JaJP)
	}
	return nil
}

func (s *eventCounterService) GetUserCountV2(
	ctx context.Context,
	req *ecproto.GetUserCountV2Request,
) (*ecproto.GetUserCountV2Response, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err = validateGetUserCountV2Request(req); err != nil {
		return nil, err
	}
	startAt := time.Unix(req.StartAt, 0)
	endAt := time.Unix(req.EndAt, 0)
	headers, rows, err := s.druidQuerier.QueryUserCount(ctx, req.EnvironmentNamespace, startAt, endAt)
	if err != nil {
		s.logger.Error(
			"Failed to query user count",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.Time("startAt", startAt),
				zap.Time("endAt", endAt),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	eventCount, userCount := convToUserCount(headers, rows)
	return &ecproto.GetUserCountV2Response{
		EventCount: eventCount,
		UserCount:  userCount,
	}, nil
}

func (s *eventCounterService) GetMAUCount(
	ctx context.Context,
	req *ecproto.GetMAUCountRequest,
) (*ecproto.GetMAUCountResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if req.YearMonth == "" {
		dt, err := statusMAUYearMonthRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "year_month"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	userCount, eventCount, err := s.userCountStorage.GetMAUCount(ctx, req.EnvironmentNamespace, req.YearMonth)
	if err != nil {
		s.logger.Error(
			"Failed to get the mau count",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.String("yearMonth", req.YearMonth),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &ecproto.GetMAUCountResponse{
		UserCount:  userCount,
		EventCount: eventCount,
	}, nil
}

func validateGetUserCountV2Request(req *ecproto.GetUserCountV2Request) error {
	if req.StartAt == 0 {
		return localizedError(statusStartAtRequired, locale.JaJP)
	}
	if req.EndAt == 0 {
		return localizedError(statusEndAtRequired, locale.JaJP)
	}
	if req.StartAt > req.EndAt {
		return localizedError(statusStartAtIsAfterEndAt, locale.JaJP)
	}
	return nil
}

func convToUserCount(headers *ecproto.Row, rows []*ecproto.Row) (eventCount, userCount int64) {
	for _, row := range rows {
		for i, cell := range row.Cells {
			switch headers.Cells[i].Value {
			case ecdruid.ColumnUserTotal:
				eventCount = int64(cell.ValueDouble)
			case ecdruid.ColumnUserCount:
				userCount = int64(cell.ValueDouble)
			}
		}
	}
	return
}

func (s *eventCounterService) ListUserMetadata(
	ctx context.Context,
	req *ecproto.ListUserMetadataRequest,
) (*ecproto.ListUserMetadataResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	data, err := s.druidQuerier.QuerySegmentMetadata(ctx, req.EnvironmentNamespace, ecdruid.DataTypeGoalEvents)
	if err != nil {
		s.logger.Error(
			"Failed to query segment metadata",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &ecproto.ListUserMetadataResponse{Data: data}, nil
}

func (s *eventCounterService) getExperimentResultMySQL(
	ctx context.Context,
	id, environmentNamespace string,
) (*ecproto.ExperimentResult, error) {
	result, err := s.mysqlExperimentResultStorage.GetExperimentResult(ctx, id, environmentNamespace)
	if err != nil {
		if err == v2ecstorage.ErrExperimentResultNotFound {
			return nil, err
		}
		s.logger.Error(
			"Failed to get experiment count",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
			)...,
		)
		return nil, err
	}
	return result.ExperimentResult, nil
}

func (s *eventCounterService) checkRole(
	ctx context.Context,
	requiredRole accountproto.Account_Role,
	environmentNamespace string,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckRole(ctx, requiredRole, func(email string) (*accountproto.GetAccountResponse, error) {
		return s.accountClient.GetAccount(ctx, &accountproto.GetAccountRequest{
			Email:                email,
			EnvironmentNamespace: environmentNamespace,
		})
	})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Info(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.UnauthenticatedError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		case codes.PermissionDenied:
			s.logger.Info(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			dt, err := statusPermissionDenied.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.PermissionDenied),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
	}
	return editor, nil
}
