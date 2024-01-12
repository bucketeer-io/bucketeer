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
	"sort"
	"strings"
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
	EventCountPrefix             = "ec"
	UserCountPrefix              = "uc"
	opsEvaluationUserCountPrefix = "autoops:evaluation"
	opsGoalUserCountPrefix       = "autoops:goal"
	defaultVariationID           = "default"
	twentyFourHours              = 24
	pfMergeKey                   = "pfmerge-key"
)

var (
	errUnknownTimeRange = errors.New("eventcounter: a time range is unknown")
)

type eventCounterService struct {
	experimentClient             experimentclient.Client
	featureClient                featureclient.Client
	accountClient                accountclient.Client
	eventStorage                 v2ecstorage.EventStorage
	mysqlExperimentResultStorage v2ecstorage.ExperimentResultStorage
	mysqlMAUSummaryStorage       v2ecstorage.MAUSummaryStorage
	userCountStorage             v2ecstorage.UserCountStorage
	metrics                      metrics.Registerer
	evaluationCountCacher        cachev3.EventCounterCache
	location                     *time.Location
	logger                       *zap.Logger
}

func NewEventCounterService(
	mc mysql.Client,
	e experimentclient.Client,
	f featureclient.Client,
	a accountclient.Client,
	b bqquerier.Client,
	bigQueryDataSet string,
	r metrics.Registerer,
	redis cache.MultiGetDeleteCountCache,
	loc *time.Location,
	l *zap.Logger,
) rpc.Service {
	registerMetrics(r)
	return &eventCounterService{
		experimentClient:             e,
		featureClient:                f,
		accountClient:                a,
		eventStorage:                 v2ecstorage.NewEventStorage(b, bigQueryDataSet, l),
		mysqlExperimentResultStorage: v2ecstorage.NewExperimentResultStorage(mc),
		mysqlMAUSummaryStorage:       v2ecstorage.NewMAUSummaryStorage(mc),
		userCountStorage:             v2ecstorage.NewUserCountStorage(mc),
		metrics:                      r,
		evaluationCountCacher:        cachev3.NewEventCountCache(redis),
		location:                     loc,
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
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
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
			Message: localizer.MustLocalizeWithTemplate(locale.StartAtIsAfterEndAt),
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

func (s *eventCounterService) GetEvaluationTimeseriesCount(
	ctx context.Context,
	req *ecproto.GetEvaluationTimeseriesCountRequest,
) (*ecproto.GetEvaluationTimeseriesCountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateGetEvaluationTimeseriesCount(req, localizer); err != nil {
		return nil, err
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
	// This timestamp will be used as `Timestamps` field in ecproto.Timeseries.
	timestamps, timestampUnit, err := s.getTimestamps(req.TimeRange)
	if err != nil {
		dt, err := statusUnknownTimeRange.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time_range"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	hourlyTimeStamps := getHourlyTimeStamps(timestamps, timestampUnit)
	vIDs := getVariationIDs(resp.Feature.Variations)
	variationTSEvents := make([]*ecproto.VariationTimeseries, 0, len(vIDs))
	variationTSUsers := make([]*ecproto.VariationTimeseries, 0, len(vIDs))
	for _, vID := range vIDs {
		eventCountKeys := s.getEventCountKeys(
			hourlyTimeStamps,
			req.EnvironmentNamespace,
			req.FeatureId,
			vID,
		)
		userCountKeys := s.getUserCountKeys(
			hourlyTimeStamps,
			req.EnvironmentNamespace,
			req.FeatureId,
			vID,
		)
		eventCounts, err := s.getEventCounts(eventCountKeys, timestampUnit)
		if err != nil {
			s.logCountError(
				ctx,
				err,
				"Failed to get event counts", req.EnvironmentNamespace, req.FeatureId, vID,
				resp.Feature.Version,
				timestampUnit, req.TimeRange,
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
		totalEventCounts := s.getTotalEventCounts(eventCounts)
		userCounts, err := s.getUserCounts(
			userCountKeys,
			req.FeatureId,
			req.EnvironmentNamespace,
			timestampUnit,
		)
		if err != nil {
			s.logCountError(
				ctx,
				err,
				"Failed to get user counts", req.EnvironmentNamespace, req.FeatureId, vID,
				resp.Feature.Version,
				timestampUnit, req.TimeRange,
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
		totalUserCounts, err := s.getTotalUserCounts(
			userCountKeys,
			req.FeatureId,
			req.EnvironmentNamespace,
		)
		if err != nil {
			s.logCountError(
				ctx,
				err,
				"Failed to get user counts", req.EnvironmentNamespace, req.FeatureId, vID,
				resp.Feature.Version,
				timestampUnit, req.TimeRange,
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
				Timestamps:  timestamps,
				Values:      userCounts,
				Unit:        timestampUnit,
				TotalCounts: totalUserCounts,
			},
		})
		variationTSEvents = append(variationTSEvents, &ecproto.VariationTimeseries{
			VariationId: vID,
			Timeseries: &ecproto.Timeseries{
				Timestamps:  timestamps,
				Values:      eventCounts,
				Unit:        timestampUnit,
				TotalCounts: totalEventCounts,
			},
		})
	}
	return &ecproto.GetEvaluationTimeseriesCountResponse{
		EventCounts: variationTSEvents,
		UserCounts:  variationTSUsers,
	}, nil
}

func (s *eventCounterService) validateGetEvaluationTimeseriesCount(
	req *ecproto.GetEvaluationTimeseriesCountRequest,
	localizer locale.Localizer,
) error {
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
	if req.TimeRange == ecproto.GetEvaluationTimeseriesCountRequest_UNKNOWN {
		dt, err := statusUnknownTimeRange.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "time_range"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

type multiError []error

func (m multiError) Error() string {
	str := make([]string, 0, len(m))
	for _, e := range m {
		if e != nil {
			s := e.Error()
			str = append(str, s)
		}
	}
	return fmt.Sprintf("%d errors: %s", len(str), strings.Join(str, ", "))
}

func (s *eventCounterService) getEventCounts(
	keys [][]string,
	unit ecproto.Timeseries_Unit,
) ([]float64, error) {
	if unit == ecproto.Timeseries_HOUR {
		return s.evaluationCountCacher.GetEventCounts(keys[0])
	}
	return s.evaluationCountCacher.GetEventCountsV2(keys)
}

func (s *eventCounterService) getTotalEventCounts(
	eventCounts []float64,
) int64 {
	total := float64(0)
	for _, count := range eventCounts {
		total += count
	}
	return int64(total)
}

func (s *eventCounterService) getUserCounts(
	keys [][]string,
	featureID, environmentNamespace string,
	unit ecproto.Timeseries_Unit,
) ([]float64, error) {
	if unit == ecproto.Timeseries_HOUR {
		return s.getHourlyUserCounts(keys, featureID, environmentNamespace)
	}
	return s.getDailyUserCounts(keys, featureID, environmentNamespace)
}

func (s *eventCounterService) getHourlyUserCounts(
	days [][]string,
	featureID, environmentNamespace string,
) ([]float64, error) {
	hours := days[0]
	counts := make([]float64, 0, len(hours))
	for _, hour := range hours {
		c, err := s.evaluationCountCacher.GetUserCount(hour)
		if err != nil {
			return nil, err
		}
		counts = append(counts, float64(c))
	}
	return counts, nil
}

func (s *eventCounterService) getDailyUserCounts(
	days [][]string,
	featureID, environmentNamespace string,
) ([]float64, error) {
	counts := make([]float64, 0, len(days))
	for _, day := range days {
		c, err := s.countUniqueUser(
			day,
			featureID, environmentNamespace,
		)
		if err != nil {
			return nil, err
		}
		counts = append(counts, c)
	}
	return counts, nil
}

func (s *eventCounterService) flattenAry(
	keys [][]string,
) []string {
	flat := []string{}
	for _, k := range keys {
		flat = append(flat, k...)
	}
	return flat
}

func (s *eventCounterService) getTotalUserCounts(
	userCountKeys [][]string,
	featureID, environmentNamespace string,
) (int64, error) {
	flat := s.flattenAry(userCountKeys)
	count, err := s.countUniqueUser(
		flat,
		featureID, environmentNamespace,
	)
	if err != nil {
		return 0, err
	}
	return int64(count), nil
}

func (s *eventCounterService) countUniqueUser(
	userCountKeys []string,
	featureID, environmentNamespace string,
) (count float64, err multiError) {
	key := newPFMergeKey(
		UserCountPrefix,
		featureID,
		environmentNamespace,
	)
	// We need to count the number of unique users in the target term.
	if e := s.evaluationCountCacher.MergeMultiKeys(key, userCountKeys); e != nil {
		err = append(err, e)
		return
	}
	defer func() {
		if e := s.evaluationCountCacher.DeleteKey(key); e != nil {
			err = append(err, e)
		}
	}()
	c, e := s.evaluationCountCacher.GetUserCount(key)
	count = float64(c)
	if e != nil {
		err = append(err, e)
		return
	}
	return
}

func (*eventCounterService) getEventCountKeys(
	hourlyTimeStamps [][]int64,
	environmentNamespace string,
	featureID string,
	vID string,
) [][]string {
	eventCountKeys := make([][]string, 0, len(hourlyTimeStamps))
	for _, twentyFourHours := range hourlyTimeStamps {
		ecHourlyKeys := make([]string, 0, len(twentyFourHours))
		for _, hour := range twentyFourHours {
			ec := newEvaluationCountkey(EventCountPrefix, featureID, vID, environmentNamespace, hour)
			ecHourlyKeys = append(ecHourlyKeys, ec)
		}
		eventCountKeys = append(eventCountKeys, ecHourlyKeys)
	}
	return eventCountKeys
}

func (*eventCounterService) getUserCountKeys(
	hourlyTimeStamps [][]int64,
	environmentNamespace string,
	featureID string,
	vID string,
) [][]string {
	userCountKeys := make([][]string, 0, len(hourlyTimeStamps))
	for _, twentyFourHours := range hourlyTimeStamps {
		ucHourlyKeys := make([]string, 0, len(twentyFourHours))
		for _, hour := range twentyFourHours {
			uc := newEvaluationCountkey(UserCountPrefix, featureID, vID, environmentNamespace, hour)
			ucHourlyKeys = append(ucHourlyKeys, uc)
		}
		userCountKeys = append(userCountKeys, ucHourlyKeys)
	}
	return userCountKeys
}

func (s *eventCounterService) logCountError(
	ctx context.Context,
	err error,
	msg, environmentNamespace, featureID, vID string,
	featureVersion int32,
	unit ecproto.Timeseries_Unit,
	timeRange ecproto.GetEvaluationTimeseriesCountRequest_TimeRange,
) {
	s.logger.Error(
		msg,
		log.FieldsFromImcomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("unit", unit.String()),
			zap.String("timeRange", timeRange.String()),
			zap.String("featureId", featureID),
			zap.Int32("featureVersion", featureVersion),
			zap.String("variationId", vID),
		)...,
	)
}

func newPFMergeKey(
	kind, featureID, environmentNamespace string,
) string {
	return cache.MakeKey(
		kind,
		fmt.Sprintf("%s:%s", pfMergeKey, featureID),
		environmentNamespace,
	)
}

func truncateDate(loc *time.Location, t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

func truncateHour(loc *time.Location, t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, loc)
}

func getStartTime(loc *time.Location, endAt time.Time, durationDays int) time.Time {
	return endAt.In(loc).AddDate(0, 0, -durationDays)
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

func getOneDayTimestamps(timestamp time.Time) []int64 {
	timestamps := make([]int64, 0, twentyFourHours)
	for i := 0; i < twentyFourHours; i++ {
		ts := timestamp.Add(time.Duration(i) * time.Hour).Unix()
		timestamps = append(timestamps, ts)
	}
	return timestamps
}

func (s *eventCounterService) getTimestamps(
	timeRange ecproto.GetEvaluationTimeseriesCountRequest_TimeRange,
) ([]int64, ecproto.Timeseries_Unit, error) {
	endAt := time.Now()
	switch timeRange {
	case ecproto.GetEvaluationTimeseriesCountRequest_TWENTY_FOUR_HOURS:
		startAt := getStartTime(s.location, endAt, 1)
		truncated := truncateHour(s.location, startAt).Add(time.Hour * 1)
		return getOneDayTimestamps(truncated), ecproto.Timeseries_HOUR, nil
	case ecproto.GetEvaluationTimeseriesCountRequest_SEVEN_DAYS:
		startAt := getStartTime(s.location, endAt, 6)
		truncated := truncateDate(s.location, startAt)
		return getDailyTimestamps(truncated, 6), ecproto.Timeseries_DAY, nil
	case ecproto.GetEvaluationTimeseriesCountRequest_FOURTEEN_DAYS:
		startAt := getStartTime(s.location, endAt, 13)
		truncated := truncateDate(s.location, startAt)
		return getDailyTimestamps(truncated, 13), ecproto.Timeseries_DAY, nil
	case ecproto.GetEvaluationTimeseriesCountRequest_THIRTY_DAYS:
		startAt := getStartTime(s.location, endAt, 29)
		truncated := truncateDate(s.location, startAt)
		return getDailyTimestamps(truncated, 29), ecproto.Timeseries_DAY, nil
	default:
		return nil, 0, errUnknownTimeRange
	}
}

func getDailyTimestamps(startAt time.Time, limit int) []int64 {
	timestamps := make([]int64, 0, limit)
	for i := 0; i <= limit; i++ {
		ts := startAt.AddDate(0, 0, i).Unix()
		timestamps = append(timestamps, ts)
	}
	return timestamps
}

/*
getHourlyTimeStamps returns a two-dimensional array. For example,

[

	["2014-11-01 00:00:00", "2014-11-01 01:00:00", "2014-11-01 02:00:00", "2014-11-01 03:00:00", ...],
	["2014-11-02 00:00:00", "2014-11-02 01:00:00", "2014-11-02 02:00:00", "2014-11-02 03:00:00", ...],
	["2014-11-03 00:00:00", "2014-11-03 01:00:00", "2014-11-03 02:00:00", "2014-11-03 03:00:00", ...],
	...

]
*/
func getHourlyTimeStamps(days []int64, unit ecproto.Timeseries_Unit) [][]int64 {
	if unit == ecproto.Timeseries_HOUR {
		return [][]int64{days}
	}
	timestamps := make([][]int64, 0, len(days))
	for _, day := range days {
		t := time.Unix(int64(day), 0)
		timestamps = append(timestamps, getOneDayTimestamps(t))
	}
	return timestamps
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
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if req.ExperimentId == "" {
		dt, err := statusExperimentIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "experiment_id"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
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
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if req.FeatureId == "" {
		dt, err := statusFeatureIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
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

func (s *eventCounterService) GetExperimentGoalCount(
	ctx context.Context,
	req *ecproto.GetExperimentGoalCountRequest,
) (*ecproto.GetExperimentGoalCountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
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
		req.GoalId,
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
			Message: localizer.MustLocalizeWithTemplate(locale.StartAtIsAfterEndAt),
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

func (s *eventCounterService) GetMAUCount(
	ctx context.Context,
	req *ecproto.GetMAUCountRequest,
) (*ecproto.GetMAUCountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
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

func (s *eventCounterService) SummarizeMAUCounts(
	ctx context.Context,
	req *ecproto.SummarizeMAUCountsRequest,
) (*ecproto.SummarizeMAUCountsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkAdminRole(ctx, localizer)
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
	summaries := make([]*ecproto.MAUSummary, 0)
	// Get the mau counts grouped by sourceID and environmentID.
	groupBySourceID, err := s.userCountStorage.GetMAUCountsGroupBySourceID(ctx, req.YearMonth)
	if err != nil {
		s.logger.Error(
			"Failed to get the mau counts by sourceID",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	summaries = append(summaries, groupBySourceID...)
	// Get the mau counts grouped by environmentID.
	groupByEnvID, err := s.userCountStorage.GetMAUCounts(ctx, req.YearMonth)
	if err != nil {
		s.logger.Error(
			"Failed to get the mau counts",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	summaries = append(summaries, groupByEnvID...)
	s.logger.Debug("SummarizeMAUCounts result", zap.Any("summaries", summaries))
	for _, summary := range summaries {
		summary.IsFinished = req.IsFinished
		summary.CreatedAt = time.Now().Unix()
		summary.UpdatedAt = time.Now().Unix()
		err := s.mysqlMAUSummaryStorage.UpsertMAUSummary(ctx, summary)
		if err != nil {
			s.logger.Error(
				"Failed to upsert the mau summary",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.Any("summary", summary),
				)...,
			)
			return nil, err
		}
	}
	return &ecproto.SummarizeMAUCountsResponse{}, nil
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

func (s *eventCounterService) GetOpsEvaluationUserCount(
	ctx context.Context,
	req *ecproto.GetOpsEvaluationUserCountRequest,
) (*ecproto.GetOpsEvaluationUserCountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetOpsEvaluationUserCountRequest(req, localizer); err != nil {
		return nil, err
	}
	cacheKey := newOpsEvaluationUserCountKey(
		opsEvaluationUserCountPrefix,
		req.OpsRuleId,
		req.ClauseId,
		req.FeatureId,
		int(req.FeatureVersion),
		req.VariationId,
		req.EnvironmentNamespace,
	)
	userCount, err := s.evaluationCountCacher.GetUserCount(cacheKey)
	if err != nil {
		s.logger.Error(
			"Failed to get ops evaluation user count",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.String("opsRuleId", req.OpsRuleId),
				zap.String("clauseId", req.ClauseId),
				zap.String("featureId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
				zap.String("variationId", req.VariationId),
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
	return &ecproto.GetOpsEvaluationUserCountResponse{
		OpsRuleId: req.OpsRuleId,
		ClauseId:  req.ClauseId,
		Count:     userCount,
	}, nil
}

func validateGetOpsEvaluationUserCountRequest(
	req *ecproto.GetOpsEvaluationUserCountRequest,
	localizer locale.Localizer,
) error {
	if req.OpsRuleId == "" {
		dt, err := statusAutoOpsRuleIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "ops_rule_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.ClauseId == "" {
		dt, err := statusClauseIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause_id"),
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
	if req.FeatureVersion == 0 {
		dt, err := statusFeatureVersionRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_version"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.VariationId == "" {
		dt, err := statusVariationIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func newOpsEvaluationUserCountKey(
	kind, opsRuleID, clauseID, featureID string,
	featureVersion int,
	variationID, environmentNamespace string,
) string {
	return cache.MakeKey(
		kind,
		fmt.Sprintf("%s:%d:%s:%s:%s", featureID, featureVersion, opsRuleID, clauseID, variationID),
		environmentNamespace,
	)
}

func (s *eventCounterService) GetOpsGoalUserCount(
	ctx context.Context,
	req *ecproto.GetOpsGoalUserCountRequest,
) (*ecproto.GetOpsGoalUserCountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetOpsGoalUserCountRequest(req, localizer); err != nil {
		return nil, err
	}
	cacheKey := newOpsGoalUserCountKey(
		opsGoalUserCountPrefix,
		req.OpsRuleId,
		req.ClauseId,
		req.FeatureId,
		int(req.FeatureVersion),
		req.VariationId,
		req.EnvironmentNamespace,
	)
	userCount, err := s.evaluationCountCacher.GetUserCount(cacheKey)
	if err != nil {
		s.logger.Error(
			"Failed to get ops goal user count",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.String("opsRuleId", req.OpsRuleId),
				zap.String("clauseId", req.ClauseId),
				zap.String("featureId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
				zap.String("variationId", req.VariationId),
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
	return &ecproto.GetOpsGoalUserCountResponse{
		OpsRuleId: req.OpsRuleId,
		ClauseId:  req.ClauseId,
		Count:     userCount,
	}, nil
}

func validateGetOpsGoalUserCountRequest(
	req *ecproto.GetOpsGoalUserCountRequest,
	localizer locale.Localizer,
) error {
	if req.OpsRuleId == "" {
		dt, err := statusAutoOpsRuleIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "ops_rule_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.ClauseId == "" {
		dt, err := statusClauseIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause_id"),
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
	if req.FeatureVersion == 0 {
		dt, err := statusFeatureVersionRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_version"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.VariationId == "" {
		dt, err := statusVariationIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func newOpsGoalUserCountKey(
	kind, opsRuleID, clauseID, featureID string,
	featureVersion int,
	variationID, environmentNamespace string,
) string {
	return cache.MakeKey(
		kind,
		fmt.Sprintf("%s:%d:%s:%s:%s", featureID, featureVersion, opsRuleID, clauseID, variationID),
		environmentNamespace,
	)
}

func (s *eventCounterService) checkRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Environment,
	environmentNamespace string,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckRole(
		ctx,
		requiredRole,
		environmentNamespace,
		func(email string) (*accountproto.AccountV2, error) {
			resp, err := s.accountClient.GetAccountV2ByEnvironmentID(ctx, &accountproto.GetAccountV2ByEnvironmentIDRequest{
				Email:         email,
				EnvironmentId: environmentNamespace,
			})
			if err != nil {
				return nil, err
			}
			return resp.Account, nil
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

func (s *eventCounterService) checkAdminRole(
	ctx context.Context,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckAdminRole(ctx)
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Info(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
