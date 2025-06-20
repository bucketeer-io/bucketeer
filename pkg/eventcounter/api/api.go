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
	"fmt"
	"strings"
	"sync"
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
	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
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
	pfMergeKind                  = "pfmerge"
	pfMergeExpiration            = 10 * time.Minute
)

type DataWarehouseConfig struct {
	Type      string                      `yaml:"type"`
	BatchSize int                         `yaml:"batchSize"`
	Timezone  string                      `yaml:"timezone"`
	BigQuery  DataWarehouseBigQueryConfig `yaml:"bigquery"`
	MySQL     DataWarehouseMySQLConfig    `yaml:"mysql"`
}

type DataWarehouseBigQueryConfig struct {
	Project  string `yaml:"project"`
	Dataset  string `yaml:"dataset"`
	Location string `yaml:"location"`
}

type DataWarehouseMySQLConfig struct {
	UseMainConnection bool   `yaml:"useMainConnection"`
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	User              string `yaml:"user"`
	Password          string `yaml:"password"`
	Database          string `yaml:"database"`
}

var (
	errUnknownTimeRange = errors.New("eventcounter: a time range is unknown")
)

type options struct {
	logger              *zap.Logger
	dataWarehouseConfig *DataWarehouseConfig
}

type Option func(*options)

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

func WithDataWarehouse(dataWarehouseType string) Option {
	return func(opts *options) {
		// Maintain backward compatibility by setting a basic config
		opts.dataWarehouseConfig = &DataWarehouseConfig{
			Type: dataWarehouseType,
		}
	}
}

func WithDataWarehouseConfig(config *DataWarehouseConfig) Option {
	return func(opts *options) {
		opts.dataWarehouseConfig = config
	}
}

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
	opts ...Option,
) rpc.Service {
	dopts := &options{
		logger: l,
		dataWarehouseConfig: &DataWarehouseConfig{
			Type: "bigquery", // default
		},
	}
	for _, opt := range opts {
		opt(dopts)
	}

	registerMetrics(r)

	var eventStorage v2ecstorage.EventStorage
	switch dopts.dataWarehouseConfig.Type {
	case "mysql":
		// Use the main MySQL client if useMainConnection is true or no custom connection specified
		if dopts.dataWarehouseConfig.MySQL.UseMainConnection || dopts.dataWarehouseConfig.MySQL.Host == "" {
			eventStorage = v2ecstorage.NewMySQLEventStorage(mc, dopts.logger)
		} else {
			// Create custom MySQL client with the specified connection details
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			customMySQLClient, err := createCustomMySQLClient(
				ctx,
				dopts.dataWarehouseConfig.MySQL,
				dopts.logger,
			)
			if err != nil {
				dopts.logger.Error("Failed to create custom MySQL client for data warehouse",
					zap.Error(err),
					zap.String("host", dopts.dataWarehouseConfig.MySQL.Host),
					zap.String("database", dopts.dataWarehouseConfig.MySQL.Database),
				)
				// Return nil to cause service initialization to fail
				// This prevents data inconsistency by ensuring we don't silently fall back
				return nil
			}

			dopts.logger.Info("Using custom MySQL connection for data warehouse",
				zap.String("host", dopts.dataWarehouseConfig.MySQL.Host),
				zap.String("database", dopts.dataWarehouseConfig.MySQL.Database),
			)
			eventStorage = v2ecstorage.NewMySQLEventStorage(customMySQLClient, dopts.logger)
		}
	case "bigquery":
		eventStorage = v2ecstorage.NewEventStorage(b, bigQueryDataSet, dopts.logger)
	default:
		// Default to BigQuery for backward compatibility
		eventStorage = v2ecstorage.NewEventStorage(b, bigQueryDataSet, dopts.logger)
	}

	return &eventCounterService{
		experimentClient:             e,
		featureClient:                f,
		accountClient:                a,
		eventStorage:                 eventStorage,
		mysqlExperimentResultStorage: v2ecstorage.NewExperimentResultStorage(mc),
		mysqlMAUSummaryStorage:       v2ecstorage.NewMAUSummaryStorage(mc),
		userCountStorage:             v2ecstorage.NewUserCountStorage(mc),
		metrics:                      r,
		evaluationCountCacher:        cachev3.NewEventCountCache(redis),
		location:                     loc,
		logger:                       dopts.logger.Named("api"),
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
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
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
		req.EnvironmentId,
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
				zap.String("environmentId", req.EnvironmentId),
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
	// we should keep the order of variationIDs
	for _, vid := range variationIDs {
		vcs = append(vcs, vcsMap[vid])
	}
	return vcs
}

func (s *eventCounterService) GetEvaluationTimeseriesCount(
	ctx context.Context,
	req *ecproto.GetEvaluationTimeseriesCountRequest,
) (*ecproto.GetEvaluationTimeseriesCountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateGetEvaluationTimeseriesCount(req, localizer); err != nil {
		return nil, err
	}
	resp, err := s.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		EnvironmentId: req.EnvironmentId,
		Id:            req.FeatureId,
	})
	if err != nil {
		s.logger.Error(
			"Failed to get feature",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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

	variationEventCountsMap := make(map[string][]float64, len(vIDs))
	variationUserCountsMap := make(map[string][]float64, len(vIDs))
	variationTotalEventCountsMap := make(map[string]int64, len(vIDs))
	variationTotalUserCountsMap := make(map[string]int64, len(vIDs))

	var mu sync.Mutex
	var eg errgroup.Group

	for _, variationID := range vIDs {
		eg.Go(func() error {

			// Get event data
			eventCountKeys := s.getEventCountKeys(
				hourlyTimeStamps,
				req.EnvironmentId,
				req.FeatureId,
				variationID,
			)
			eventCounts, err := s.getEventCounts(eventCountKeys, timestampUnit)
			if err != nil {
				s.logCountError(
					ctx,
					err,
					"Failed to get event counts", req.EnvironmentId, req.FeatureId, variationID,
					resp.Feature.Version,
					timestampUnit, req.TimeRange,
				)
				return err
			}

			// Get user data
			userCountKeys := s.getUserCountKeys(
				hourlyTimeStamps,
				req.EnvironmentId,
				req.FeatureId,
				variationID,
			)
			userCounts, err := s.getUserCounts(
				userCountKeys,
				req.FeatureId,
				req.EnvironmentId,
				variationID,
				timestampUnit,
			)
			if err != nil {
				s.logCountError(
					ctx,
					err,
					"Failed to get user counts", req.EnvironmentId, req.FeatureId, variationID,
					resp.Feature.Version,
					timestampUnit, req.TimeRange,
				)
				return err
			}

			totalUserCounts, err := s.getTotalUserCounts(
				userCountKeys,
				req.FeatureId,
				req.EnvironmentId,
				variationID,
			)
			if err != nil {
				s.logCountError(
					ctx,
					err,
					"Failed to get total user counts", req.EnvironmentId, req.FeatureId, variationID,
					resp.Feature.Version,
					timestampUnit, req.TimeRange,
				)
				return err
			}

			mu.Lock()
			variationEventCountsMap[variationID] = eventCounts
			variationUserCountsMap[variationID] = userCounts
			variationTotalEventCountsMap[variationID] = s.getTotalEventCounts(eventCounts)
			variationTotalUserCountsMap[variationID] = totalUserCounts
			mu.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		dt, errDt := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if errDt != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}

	variationTSEvents := make([]*ecproto.VariationTimeseries, 0, len(vIDs))
	variationTSUsers := make([]*ecproto.VariationTimeseries, 0, len(vIDs))

	for _, vID := range vIDs {
		variationTSEvents = append(variationTSEvents, &ecproto.VariationTimeseries{
			VariationId: vID,
			Timeseries: &ecproto.Timeseries{
				Timestamps:  timestamps,
				Values:      variationEventCountsMap[vID],
				Unit:        timestampUnit,
				TotalCounts: variationTotalEventCountsMap[vID],
			},
		})

		variationTSUsers = append(variationTSUsers, &ecproto.VariationTimeseries{
			VariationId: vID,
			Timeseries: &ecproto.Timeseries{
				Timestamps:  timestamps,
				Values:      variationUserCountsMap[vID],
				Unit:        timestampUnit,
				TotalCounts: variationTotalUserCountsMap[vID],
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
	featureID, environmentId, variationID string,
	unit ecproto.Timeseries_Unit,
) ([]float64, error) {
	if unit == ecproto.Timeseries_HOUR {
		return s.getHourlyUserCounts(keys, featureID, environmentId)
	}
	return s.getDailyUserCounts(keys, featureID, environmentId, variationID)
}

func (s *eventCounterService) getHourlyUserCounts(
	days [][]string,
	featureID, environmentId string,
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
	featureID, environmentId, variationID string,
) ([]float64, error) {
	counts := make([]float64, 0, len(days))
	for _, day := range days {
		c, err := s.countUniqueUser(
			day,
			featureID,
			environmentId,
			variationID,
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
	featureID, environmentId, variationID string,
) (int64, error) {
	flat := s.flattenAry(userCountKeys)
	count, err := s.countUniqueUser(
		flat,
		featureID, environmentId,
		variationID,
	)
	if err != nil {
		return 0, err
	}
	return int64(count), nil
}

func (s *eventCounterService) countUniqueUser(
	userCountKeys []string,
	featureID, environmentId, variationID string,
) (count float64, err multiError) {
	key := newPFMergeKey(
		UserCountPrefix,
		featureID,
		environmentId,
		variationID,
	)
	// We need to count the number of unique users in the target term.
	if e := s.evaluationCountCacher.MergeMultiKeys(key, userCountKeys, pfMergeExpiration); e != nil {
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
	environmentId string,
	featureID string,
	vID string,
) [][]string {
	eventCountKeys := make([][]string, 0, len(hourlyTimeStamps))
	for _, twentyFourHours := range hourlyTimeStamps {
		ecHourlyKeys := make([]string, 0, len(twentyFourHours))
		for _, hour := range twentyFourHours {
			ec := newEvaluationCountkey(EventCountPrefix, featureID, vID, environmentId, hour)
			ecHourlyKeys = append(ecHourlyKeys, ec)
		}
		eventCountKeys = append(eventCountKeys, ecHourlyKeys)
	}
	return eventCountKeys
}

func (*eventCounterService) getUserCountKeys(
	hourlyTimeStamps [][]int64,
	environmentId string,
	featureID string,
	vID string,
) [][]string {
	userCountKeys := make([][]string, 0, len(hourlyTimeStamps))
	for _, twentyFourHours := range hourlyTimeStamps {
		ucHourlyKeys := make([]string, 0, len(twentyFourHours))
		for _, hour := range twentyFourHours {
			uc := newEvaluationCountkey(UserCountPrefix, featureID, vID, environmentId, hour)
			ucHourlyKeys = append(ucHourlyKeys, uc)
		}
		userCountKeys = append(userCountKeys, ucHourlyKeys)
	}
	return userCountKeys
}

func (s *eventCounterService) logCountError(
	ctx context.Context,
	err error,
	msg, environmentId, featureID, vID string,
	featureVersion int32,
	unit ecproto.Timeseries_Unit,
	timeRange ecproto.GetEvaluationTimeseriesCountRequest_TimeRange,
) {
	s.logger.Error(
		msg,
		log.FieldsFromImcomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentId", environmentId),
			zap.String("unit", unit.String()),
			zap.String("timeRange", timeRange.String()),
			zap.String("featureId", featureID),
			zap.Int32("featureVersion", featureVersion),
			zap.String("variationId", vID),
		)...,
	)
}

func newPFMergeKey(
	kind, featureID, environmentId, variationID string,
) string {
	return cache.MakeKey(
		fmt.Sprintf("%s:%s", pfMergeKind, kind),
		fmt.Sprintf("%s:%s:%s", pfMergeKey, featureID, variationID),
		environmentId,
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
	kind, featureID, variationID, environmentId string,
	ts int64,
) string {
	return cache.MakeKey(
		kind,
		fmt.Sprintf("%d:%s:%s", ts, featureID, variationID),
		environmentId,
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
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
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
	result, err := s.mysqlExperimentResultStorage.GetExperimentResult(ctx, req.ExperimentId, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2ecstorage.ErrExperimentResultNotFound) {
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
				zap.String("environmentId", req.EnvironmentId),
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
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
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
	experiments, err := s.listExperiments(ctx, req.FeatureId, req.FeatureVersion, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, storage.ErrKeyNotFound) {
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
				zap.String("environmentId", req.EnvironmentId),
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
		er, err := s.getExperimentResultMySQL(ctx, e.Id, req.EnvironmentId)
		if err != nil {
			if errors.Is(err, v2ecstorage.ErrExperimentResultNotFound) {
				getExperimentCountsCounter.WithLabelValues(codeSuccess).Inc()
			} else {
				s.logger.Error(
					"Failed to get Experiment result",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentId", req.EnvironmentId),
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
	environmentId string,
) ([]*experimentproto.Experiment, error) {
	experiments := []*experimentproto.Experiment{}
	cursor := ""
	for {
		resp, err := s.experimentClient.ListExperiments(ctx, &experimentproto.ListExperimentsRequest{
			FeatureId:      featureID,
			FeatureVersion: featureVersion,
			PageSize:       listRequestPageSize,
			Cursor:         cursor,
			EnvironmentId:  environmentId,
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
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
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
		req.EnvironmentId,
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
				zap.String("environmentId", req.EnvironmentId),
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
	// we should keep the order of variationIDs
	for _, vid := range variationIDs {
		vcs = append(vcs, vcsMap[vid])
	}
	return vcs
}

func (s *eventCounterService) GetMAUCount(
	ctx context.Context,
	req *ecproto.GetMAUCountRequest,
) (*ecproto.GetMAUCountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
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
	userCount, eventCount, err := s.userCountStorage.GetMAUCount(ctx, req.EnvironmentId, req.YearMonth)
	if err != nil {
		s.logger.Error(
			"Failed to get the mau count",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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
	_, err := s.checkSystemAdminRole(ctx, localizer)
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
	id, environmentId string,
) (*ecproto.ExperimentResult, error) {
	result, err := s.mysqlExperimentResultStorage.GetExperimentResult(ctx, id, environmentId)
	if err != nil {
		if errors.Is(err, v2ecstorage.ErrExperimentResultNotFound) {
			return nil, err
		}
		s.logger.Error(
			"Failed to get experiment count",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
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
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
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
		req.EnvironmentId,
	)
	userCount, err := s.evaluationCountCacher.GetUserCount(cacheKey)
	if err != nil {
		s.logger.Error(
			"Failed to get ops evaluation user count",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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
	variationID, environmentId string,
) string {
	return cache.MakeKey(
		kind,
		fmt.Sprintf("%s:%d:%s:%s:%s", featureID, featureVersion, opsRuleID, clauseID, variationID),
		environmentId,
	)
}

func (s *eventCounterService) GetOpsGoalUserCount(
	ctx context.Context,
	req *ecproto.GetOpsGoalUserCountRequest,
) (*ecproto.GetOpsGoalUserCountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
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
		req.EnvironmentId,
	)
	userCount, err := s.evaluationCountCacher.GetUserCount(cacheKey)
	if err != nil {
		s.logger.Error(
			"Failed to get ops goal user count",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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
	variationID, environmentId string,
) string {
	return cache.MakeKey(
		kind,
		fmt.Sprintf("%s:%d:%s:%s:%s", featureID, featureVersion, opsRuleID, clauseID, variationID),
		environmentId,
	)
}

func (s *eventCounterService) checkEnvironmentRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Environment,
	environmentId string,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckEnvironmentRole(
		ctx,
		requiredRole,
		environmentId,
		func(email string) (*accountproto.AccountV2, error) {
			resp, err := s.accountClient.GetAccountV2ByEnvironmentID(ctx, &accountproto.GetAccountV2ByEnvironmentIDRequest{
				Email:         email,
				EnvironmentId: environmentId,
			})
			if err != nil {
				return nil, err
			}
			return resp.Account, nil
		})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
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
			s.logger.Error(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
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
					zap.String("environmentId", environmentId),
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

func (s *eventCounterService) checkSystemAdminRole(
	ctx context.Context,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckSystemAdminRole(ctx)
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
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
			s.logger.Error(
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

// createCustomMySQLClient creates a dedicated MySQL client with the specified connection details
func createCustomMySQLClient(
	ctx context.Context,
	config DataWarehouseMySQLConfig,
	logger *zap.Logger,
) (mysql.Client, error) {
	// Validate required fields
	if config.Host == "" || config.Database == "" || config.User == "" {
		return nil, fmt.Errorf("mysql host, database, and user are required for custom connection")
	}

	// Set default port if not specified
	port := config.Port
	if port == 0 {
		port = 3306 // Default MySQL port
	}

	// Create MySQL client with custom connection
	client, err := mysql.NewClient(
		ctx,
		config.User,
		config.Password,
		config.Host,
		port,
		config.Database,
		mysql.WithLogger(logger),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create MySQL client: %w", err)
	}

	logger.Info("Created custom MySQL client for data warehouse",
		zap.String("host", config.Host),
		zap.Int("port", port),
		zap.String("database", config.Database),
		zap.String("user", config.User),
	)

	return client, nil
}
