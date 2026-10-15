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
	"fmt"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/errgroup"
	dwhdatabase "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/storage/v2/dwh_database"
	operationaldatabase "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/storage/v2/operational_database"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	ecproto "github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const listRequestPageSize = 500

const (
	EventCountPrefix             = "ec"
	UserCountPrefix              = "uc"
	opsEvaluationUserCountPrefix = "autoops:evaluation"
	opsGoalUserCountPrefix       = "autoops:goal"
	defaultVariationID           = "default"
	twentyFourHours              = 24
	oneHourMinutes               = 60
	// countFetchConcurrency bounds how many per-day (or per-hour) count
	// fetches run in parallel per variation, so long time ranges cannot
	// flood Redis with concurrent requests.
	countFetchConcurrency = 8
)

var (
	errUnknownTimeRange = errors.New("eventcounter: a time range is unknown")
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

type eventCounterService struct {
	experimentClient        experimentclient.Client
	featureClient           featureclient.Client
	accountClient           accountclient.Client
	eventStorage            dwhdatabase.EventStorage
	experimentResultStorage operationaldatabase.ExperimentResultStorage
	metrics                 metrics.Registerer
	evaluationCountCacher   cachev3.EventCounterCache
	// counterCache is raw access to the counter Redis, used to read the
	// daily-tier cutover marker (cachev3.EvaluationDailyCutoverKey).
	counterCache cache.Getter
	// location is a display-only concern. Counter bucket boundaries are
	// always truncated in UTC to match the persister's keys.
	location *time.Location
	logger   *zap.Logger
}

func NewEventCounterService(
	eventStorage dwhdatabase.EventStorage,
	experimentResultStorage operationaldatabase.ExperimentResultStorage,
	e experimentclient.Client,
	f featureclient.Client,
	a accountclient.Client,
	r metrics.Registerer,
	redis cache.MultiGetDeleteCountCache,
	loc *time.Location,
	l *zap.Logger,
	opts ...Option,
) rpc.Service {
	dopts := &options{
		logger: l,
	}
	for _, opt := range opts {
		opt(dopts)
	}

	registerMetrics(r)

	return &eventCounterService{
		experimentClient:        e,
		featureClient:           f,
		accountClient:           a,
		eventStorage:            eventStorage,
		experimentResultStorage: experimentResultStorage,
		metrics:                 r,
		evaluationCountCacher:   cachev3.NewEventCountCache(redis),
		counterCache:            redis,
		location:                loc,
		logger:                  dopts.logger.Named("api"),
	}
}

func (s *eventCounterService) Register(server *grpc.Server) {
	ecproto.RegisterEventCounterServiceServer(server, s)
}

func (s *eventCounterService) GetExperimentEvaluationCount(
	ctx context.Context,
	req *ecproto.GetExperimentEvaluationCountRequest,
) (*ecproto.GetExperimentEvaluationCountResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err = validateGetExperimentEvaluationCountRequest(req); err != nil {
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.Time("startAt", startAt),
				zap.Time("endAt", endAt),
				zap.String("featureId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
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
) error {
	if req.StartAt == 0 {
		return statusStartAtRequired.Err()
	}
	if req.EndAt == 0 {
		return statusEndAtRequired.Err()
	}
	if req.StartAt > req.EndAt {
		return statusStartAtIsAfterEndAt.Err()
	}
	if req.FeatureId == "" {
		return statusFeatureIDRequired.Err()
	}
	return nil
}

func (s *eventCounterService) convertEvaluationCounts(
	rows []*dwhdatabase.EvaluationEventCount,
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
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := s.validateGetEvaluationTimeseriesCount(req); err != nil {
		return nil, err
	}
	resp, err := s.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		EnvironmentId: req.EnvironmentId,
		Id:            req.FeatureId,
	})
	if err != nil {
		s.logger.Error(
			"Failed to get feature",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("featureId", req.FeatureId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	// This timestamp will be used as `Timestamps` field in ecproto.Timeseries.
	timestamps, timestampUnit, err := s.getTimestamps(req.TimeRange)
	if err != nil {
		return nil, statusUnknownTimeRange.Err()
	}
	// Daily ranges mix two sources: days on/after the cutover are served from
	// the daily-tier keys, older days fall back to their 24 hourly keys.
	var dailyCutover int64
	if timestampUnit == ecproto.Timeseries_DAY {
		dailyCutover = s.getDailyCutover()
	}
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
				timestamps,
				timestampUnit,
				dailyCutover,
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
				timestamps,
				timestampUnit,
				dailyCutover,
				req.EnvironmentId,
				req.FeatureId,
				variationID,
			)
			userCounts, err := s.getUserCounts(
				userCountKeys,
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

			totalUserCounts, err := s.getTotalUserCounts(userCountKeys)
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
		return nil, statusInternal.Err()
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
) error {
	if req.FeatureId == "" {
		return statusFeatureIDRequired.Err()
	}
	if req.TimeRange == ecproto.GetEvaluationTimeseriesCountRequest_UNKNOWN {
		return statusUnknownTimeRange.Err()
	}
	return nil
}

func (s *eventCounterService) getEventCounts(
	keys [][]string,
	unit ecproto.Timeseries_Unit,
) ([]float64, error) {
	if unit == ecproto.Timeseries_HOUR || unit == ecproto.Timeseries_MINUTE {
		// One key per point: a single pipelined GetMulti.
		return s.evaluationCountCacher.GetEventCounts(keys[0])
	}
	// One key set per point (a daily key or a day's hourly fallback keys),
	// each summed into one point value.
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
	unit ecproto.Timeseries_Unit,
) ([]float64, error) {
	if unit == ecproto.Timeseries_HOUR || unit == ecproto.Timeseries_MINUTE {
		return s.getSingleKeyUserCounts(keys)
	}
	return s.getDailyUserCounts(keys)
}

// getSingleKeyUserCounts counts the ranges whose points map to exactly one
// key each (hour and minute units).
func (s *eventCounterService) getSingleKeyUserCounts(
	days [][]string,
) ([]float64, error) {
	// The buckets are independent, so count them in parallel.
	hours := days[0]
	counts := make([]float64, len(hours))
	// The local errgroup recovers panics inside the goroutines, which plain
	// x/sync/errgroup would let crash the process.
	var eg errgroup.Group
	eg.SetLimit(countFetchConcurrency)
	for i, hour := range hours {
		eg.Go(func() error {
			c, err := s.evaluationCountCacher.GetUserCount(hour)
			if err != nil {
				return err
			}
			counts[i] = float64(c)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return counts, nil
}

// getDailyUserCounts counts one value per day. A day served from the daily
// tier has a single uc1d key; a fallback day has 24 hourly keys unioned by
// one CountUniqueUsers call. Both are HLL unions, so the sources can never
// double count a user.
func (s *eventCounterService) getDailyUserCounts(
	days [][]string,
) ([]float64, error) {
	// The days are independent, so count them in parallel.
	counts := make([]float64, len(days))
	var eg errgroup.Group
	eg.SetLimit(countFetchConcurrency)
	for i, day := range days {
		eg.Go(func() error {
			c, err := s.evaluationCountCacher.CountUniqueUsers(day)
			if err != nil {
				return err
			}
			counts[i] = float64(c)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
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

// getTotalUserCounts unions every key the range used — daily keys and any
// fallback hourly keys together — in ONE CountUniqueUsers call. HLL union
// over mixed granularities is still an exact union, so mixing tiers can never
// double count a user. Cross-slot keys are handled transparently by the
// cluster-safe PFCount in pkg/redis/v3.
func (s *eventCounterService) getTotalUserCounts(
	userCountKeys [][]string,
) (int64, error) {
	flat := s.flattenAry(userCountKeys)
	return s.evaluationCountCacher.CountUniqueUsers(flat)
}

func (s *eventCounterService) getEventCountKeys(
	timestamps []int64,
	unit ecproto.Timeseries_Unit,
	dailyCutover int64,
	environmentId string,
	featureID string,
	vID string,
) [][]string {
	return s.getCountKeys(
		cachev3.EventCountMinuteKind,
		EventCountPrefix,
		cachev3.EventCountDailyKind,
		timestamps, unit, dailyCutover, environmentId, featureID, vID,
	)
}

func (s *eventCounterService) getUserCountKeys(
	timestamps []int64,
	unit ecproto.Timeseries_Unit,
	dailyCutover int64,
	environmentId string,
	featureID string,
	vID string,
) [][]string {
	return s.getCountKeys(
		cachev3.UserCountMinuteKind,
		UserCountPrefix,
		cachev3.UserCountDailyKind,
		timestamps, unit, dailyCutover, environmentId, featureID, vID,
	)
}

// getCountKeys returns the Redis keys backing each point of the range:
//
//   - MINUTE unit: one minute-tier key per point.
//   - HOUR unit: one hourly key per point (byte-identical to the historical
//     24-hours behavior).
//   - DAY unit: per day, either the single daily-tier key (when the UTC day
//     starts on/after the cutover marker) or the day's 24 hourly keys (the
//     fallback for days that predate the daily tier). Each day is served from
//     exactly one source, so no double counting path exists. A zero cutover
//     (marker absent) serves every day from hourly keys.
func (*eventCounterService) getCountKeys(
	minuteKind, hourlyKind, dailyKind string,
	timestamps []int64,
	unit ecproto.Timeseries_Unit,
	dailyCutover int64,
	environmentId string,
	featureID string,
	vID string,
) [][]string {
	switch unit {
	case ecproto.Timeseries_MINUTE:
		keys := make([]string, 0, len(timestamps))
		for _, ts := range timestamps {
			keys = append(keys, newEvaluationCountkey(minuteKind, featureID, vID, environmentId, ts))
		}
		return [][]string{keys}
	case ecproto.Timeseries_HOUR:
		keys := make([]string, 0, len(timestamps))
		for _, ts := range timestamps {
			keys = append(keys, newEvaluationCountkey(hourlyKind, featureID, vID, environmentId, ts))
		}
		return [][]string{keys}
	default: // ecproto.Timeseries_DAY
		pointKeys := make([][]string, 0, len(timestamps))
		for _, day := range timestamps {
			if dailyCutover > 0 && day >= dailyCutover {
				key := newEvaluationCountkey(dailyKind, featureID, vID, environmentId, day)
				pointKeys = append(pointKeys, []string{key})
				continue
			}
			hourlyKeys := make([]string, 0, twentyFourHours)
			for _, hour := range getOneDayTimestamps(time.Unix(day, 0).UTC()) {
				hourlyKeys = append(hourlyKeys, newEvaluationCountkey(hourlyKind, featureID, vID, environmentId, hour))
			}
			pointKeys = append(pointKeys, hourlyKeys)
		}
		return pointKeys
	}
}

// getDailyCutover returns the unix timestamp of the first UTC midnight served
// from the daily tier, or 0 when the marker is absent or unreadable — in that
// case every day falls back to hourly keys, which is exactly the
// pre-daily-tier behavior and always safe.
func (s *eventCounterService) getDailyCutover() int64 {
	value, err := s.counterCache.Get(cachev3.EvaluationDailyCutoverKey)
	if err != nil {
		if !errors.Is(err, cache.ErrNotFound) {
			s.logger.Warn("Failed to get the daily cutover marker; falling back to hourly keys",
				zap.Error(err),
			)
		}
		return 0
	}
	b, err := cache.Bytes(value)
	if err != nil {
		s.logger.Warn("Daily cutover marker has an unexpected type; falling back to hourly keys",
			zap.Error(err),
		)
		return 0
	}
	cutover, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		s.logger.Warn("Failed to parse the daily cutover marker; falling back to hourly keys",
			zap.Error(err),
			zap.ByteString("value", b),
		)
		return 0
	}
	return cutover
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
		log.FieldsFromIncomingContext(ctx).AddFields(
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

// All bucket truncation happens in UTC to match the persister's keys. The
// previous location-based truncation only worked because whole-hour-offset
// zones share hour boundaries with UTC; fractional-offset zones (e.g.
// UTC+5:30) would miss every key, and the daily tier is materialized at UTC
// midnights forever.
func truncateDate(t time.Time) time.Time {
	u := t.UTC()
	return time.Date(u.Year(), u.Month(), u.Day(), 0, 0, 0, 0, time.UTC)
}

func truncateHour(t time.Time) time.Time {
	return t.UTC().Truncate(time.Hour)
}

func truncateMinute(t time.Time) time.Time {
	return t.UTC().Truncate(time.Minute)
}

func getStartTime(endAt time.Time, durationDays int) time.Time {
	return endAt.UTC().AddDate(0, 0, -durationDays)
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
	case ecproto.GetEvaluationTimeseriesCountRequest_ONE_HOUR:
		// 60 one-minute points ending at the current (partial) minute.
		truncated := truncateMinute(endAt.Add(-time.Hour)).Add(time.Minute)
		return getMinuteTimestamps(truncated, oneHourMinutes), ecproto.Timeseries_MINUTE, nil
	case ecproto.GetEvaluationTimeseriesCountRequest_TWENTY_FOUR_HOURS:
		startAt := getStartTime(endAt, 1)
		truncated := truncateHour(startAt).Add(time.Hour * 1)
		return getOneDayTimestamps(truncated), ecproto.Timeseries_HOUR, nil
	case ecproto.GetEvaluationTimeseriesCountRequest_SEVEN_DAYS:
		startAt := getStartTime(endAt, 6)
		truncated := truncateDate(startAt)
		return getDailyTimestamps(truncated, 6), ecproto.Timeseries_DAY, nil
	case ecproto.GetEvaluationTimeseriesCountRequest_FOURTEEN_DAYS:
		startAt := getStartTime(endAt, 13)
		truncated := truncateDate(startAt)
		return getDailyTimestamps(truncated, 13), ecproto.Timeseries_DAY, nil
	case ecproto.GetEvaluationTimeseriesCountRequest_THIRTY_DAYS:
		startAt := getStartTime(endAt, 29)
		truncated := truncateDate(startAt)
		return getDailyTimestamps(truncated, 29), ecproto.Timeseries_DAY, nil
	case ecproto.GetEvaluationTimeseriesCountRequest_SIXTY_DAYS:
		startAt := getStartTime(endAt, 59)
		truncated := truncateDate(startAt)
		return getDailyTimestamps(truncated, 59), ecproto.Timeseries_DAY, nil
	case ecproto.GetEvaluationTimeseriesCountRequest_NINETY_DAYS:
		startAt := getStartTime(endAt, 89)
		truncated := truncateDate(startAt)
		return getDailyTimestamps(truncated, 89), ecproto.Timeseries_DAY, nil
	default:
		return nil, 0, errUnknownTimeRange
	}
}

func getMinuteTimestamps(startAt time.Time, count int) []int64 {
	timestamps := make([]int64, 0, count)
	for i := 0; i < count; i++ {
		timestamps = append(timestamps, startAt.Add(time.Duration(i)*time.Minute).Unix())
	}
	return timestamps
}

func getDailyTimestamps(startAt time.Time, limit int) []int64 {
	timestamps := make([]int64, 0, limit)
	for i := 0; i <= limit; i++ {
		ts := startAt.AddDate(0, 0, i).Unix()
		timestamps = append(timestamps, ts)
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
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if req.ExperimentId == "" {
		return nil, statusExperimentIDRequired.Err()
	}
	result, err := s.experimentResultStorage.GetExperimentResult(ctx, req.ExperimentId, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, operationaldatabase.ErrExperimentResultNotFound) {
			return nil, statusExperimentResultNotFound.Err()
		}
		s.logger.Error(
			"Failed to get experiment result",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("experimentId", req.ExperimentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &ecproto.GetExperimentResultResponse{
		ExperimentResult: result.ExperimentResult,
	}, nil
}

func (s *eventCounterService) ListExperimentResults(
	ctx context.Context,
	req *ecproto.ListExperimentResultsRequest,
) (*ecproto.ListExperimentResultsResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if req.FeatureId == "" {
		return nil, statusFeatureIDRequired.Err()
	}
	experiments, err := s.listExperiments(ctx, req.FeatureId, req.FeatureVersion, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, storage.ErrKeyNotFound) {
			listExperimentCountsCounter.WithLabelValues(codeSuccess).Inc()
			return nil, statusExperimentResultNotFound.Err()
		}
		var featureVersion *int32
		if req.FeatureVersion != nil {
			featureVersion = &req.FeatureVersion.Value
		}
		s.logger.Error(
			"Failed to get Experiment list",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("featureID", req.FeatureId),
				zap.Int32p("featureVersion", featureVersion),
			)...,
		)
		listExperimentCountsCounter.WithLabelValues(codeFail).Inc()
		return nil, api.NewGRPCStatus(err).Err()
	}
	results := make(map[string]*ecproto.ExperimentResult, len(experiments))
	for _, e := range experiments {
		er, err := s.getExperimentResult(ctx, e.Id, req.EnvironmentId)
		if err != nil {
			if errors.Is(err, operationaldatabase.ErrExperimentResultNotFound) {
				getExperimentCountsCounter.WithLabelValues(codeSuccess).Inc()
			} else {
				s.logger.Error(
					"Failed to get Experiment result",
					log.FieldsFromIncomingContext(ctx).AddFields(
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
	featureVersion *wrapperspb.Int32Value,
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
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err = validateGetExperimentGoalCountRequest(req); err != nil {
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.Time("startAt", startAt),
				zap.Time("endAt", endAt),
				zap.String("featureId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
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
) error {
	if req.StartAt == 0 {
		return statusStartAtRequired.Err()
	}
	if req.EndAt == 0 {
		return statusEndAtRequired.Err()
	}
	if req.StartAt > req.EndAt {
		return statusStartAtIsAfterEndAt.Err()
	}
	if req.FeatureId == "" {
		return statusFeatureIDRequired.Err()
	}
	if req.GoalId == "" {
		return statusGoalIDRequired.Err()
	}
	return nil
}

func (s *eventCounterService) convertGoalCounts(
	rows []*dwhdatabase.GoalEventCount,
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

func (s *eventCounterService) getExperimentResult(
	ctx context.Context,
	id, environmentId string,
) (*ecproto.ExperimentResult, error) {
	result, err := s.experimentResultStorage.GetExperimentResult(ctx, id, environmentId)
	if err != nil {
		if errors.Is(err, operationaldatabase.ErrExperimentResultNotFound) {
			return nil, err
		}
		s.logger.Error(
			"Failed to get experiment count",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateGetOpsEvaluationUserCountRequest(req); err != nil {
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("opsRuleId", req.OpsRuleId),
				zap.String("clauseId", req.ClauseId),
				zap.String("featureId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
				zap.String("variationId", req.VariationId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &ecproto.GetOpsEvaluationUserCountResponse{
		OpsRuleId: req.OpsRuleId,
		ClauseId:  req.ClauseId,
		Count:     userCount,
	}, nil
}

func validateGetOpsEvaluationUserCountRequest(
	req *ecproto.GetOpsEvaluationUserCountRequest,
) error {
	if req.OpsRuleId == "" {
		return statusAutoOpsRuleIDRequired.Err()
	}
	if req.ClauseId == "" {
		return statusClauseIDRequired.Err()
	}
	if req.FeatureId == "" {
		return statusFeatureIDRequired.Err()
	}
	if req.FeatureVersion == 0 {
		return statusFeatureVersionRequired.Err()
	}
	if req.VariationId == "" {
		return statusVariationIDRequired.Err()
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
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateGetOpsGoalUserCountRequest(req); err != nil {
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("opsRuleId", req.OpsRuleId),
				zap.String("clauseId", req.ClauseId),
				zap.String("featureId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
				zap.String("variationId", req.VariationId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &ecproto.GetOpsGoalUserCountResponse{
		OpsRuleId: req.OpsRuleId,
		ClauseId:  req.ClauseId,
		Count:     userCount,
	}, nil
}

func validateGetOpsGoalUserCountRequest(
	req *ecproto.GetOpsGoalUserCountRequest,
) error {
	if req.OpsRuleId == "" {
		return statusAutoOpsRuleIDRequired.Err()
	}
	if req.ClauseId == "" {
		return statusClauseIDRequired.Err()
	}
	if req.FeatureId == "" {
		return statusFeatureIDRequired.Err()
	}
	if req.FeatureVersion == 0 {
		return statusFeatureVersionRequired.Err()
	}
	if req.VariationId == "" {
		return statusVariationIDRequired.Err()
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
) (*eventproto.Editor, error) {
	return role.CheckEnvironmentRoleWithLog(
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
		},
		s.logger,
		statusUnauthenticated.Err(),
		statusPermissionDenied.Err(),
		func(err error) error { return api.NewGRPCStatus(err).Err() },
	)
}
