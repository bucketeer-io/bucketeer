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
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	gmetadata "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	accstorage "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	auditlogclient "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/client"
	autoopsclient "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	coderefclient "github.com/bucketeer-io/bucketeer/v2/pkg/coderef/client"
	environmentclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	eventcounterclient "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/client"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	featuredomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	notificationclient "github.com/bucketeer-io/bucketeer/v2/pkg/notification/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	pushclient "github.com/bucketeer-io/bucketeer/v2/pkg/push/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	tagclient "github.com/bucketeer-io/bucketeer/v2/pkg/tag/client"
	teamclient "github.com/bucketeer-io/bucketeer/v2/pkg/team/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

const (
	listRequestSize         = 500
	secondsToReturnAllFlags = 30 * 24 * 60 * 60 // 30 days
	obfuscateAPIKeyLength   = 4
	// singleflightFetchTimeout bounds the shared fetch executed inside a
	// singleflight group when callers' request contexts are detached. It must
	// be longer than the worst-case downstream latency (DB / feature service)
	// and shorter than the SDK / load-balancer timeout so that runaway work
	// is eventually released.
	singleflightFetchTimeout = 10 * time.Second
)

var (
	ErrSDKVersionRequired      = status.Error(codes.InvalidArgument, "gateway: sdk version is required")
	ErrSourceIDRequired        = status.Error(codes.InvalidArgument, "gateway: source id is required")
	ErrUserRequired            = status.Error(codes.InvalidArgument, "gateway: user is required")
	ErrUserIDRequired          = status.Error(codes.InvalidArgument, "gateway: user id is required")
	ErrGoalIDRequired          = status.Error(codes.InvalidArgument, "gateway: goal id is required")
	ErrFeatureIDRequired       = status.Error(codes.InvalidArgument, "gateway: feature id is required")
	ErrTagRequired             = status.Error(codes.InvalidArgument, "gateway: tag is required")
	ErrMissingEvents           = status.Error(codes.InvalidArgument, "gateway: missing events")
	ErrMissingEventID          = status.Error(codes.InvalidArgument, "gateway: missing event id")
	ErrInvalidTimestamp        = status.Error(codes.InvalidArgument, "gateway: invalid timestamp")
	ErrContextCanceled         = status.Error(codes.Canceled, "gateway: context canceled")
	ErrContextDeadlineExceeded = status.Error(codes.DeadlineExceeded, "gateway: context deadline exceeded")
	ErrFeatureNotFound         = status.Error(codes.NotFound, "gateway: feature not found")
	ErrEvaluationNotFound      = status.Error(codes.NotFound, "gateway: evaluation not found")
	ErrPushNotFound            = status.Error(codes.NotFound, "gateway: push not found")
	ErrAccountNotFound         = status.Error(codes.NotFound, "gateway: account not found")
	ErrMissingAPIKey           = status.Error(codes.Unauthenticated, "gateway: missing APIKey")
	ErrInvalidAPIKey           = status.Error(codes.PermissionDenied, "gateway: invalid APIKey")
	ErrDisabledAPIKey          = status.Error(codes.PermissionDenied, "gateway: disabled APIKey")
	ErrBadRole                 = status.Error(codes.PermissionDenied, "gateway: bad role")
	ErrInternal                = status.Error(codes.Internal, "gateway: internal")
	ErrNotFound                = status.Error(codes.NotFound, "gateway: not found")

	// errCallerCanceled is wrapped around the underlying gRPC status error when
	// singleflightFetch returns because the caller's request context was
	// canceled (or its deadline expired). Call sites use errors.Is to detect
	// this case and avoid bumping internal-error metrics for a client-side
	// disconnect that is not a server-side failure.
	errCallerCanceled = errors.New("gateway: caller context canceled")

	grpcGoalEvent       = &eventproto.GoalEvent{}
	grpcEvaluationEvent = &eventproto.EvaluationEvent{}
	grpcMetricsEvent    = &eventproto.MetricsEvent{}
)

type options struct {
	apiKeyMemoryCacheTTL              time.Duration
	apiKeyMemoryCacheEvictionInterval time.Duration
	featuresMemoryCacheTTL            time.Duration
	segmentUsersMemoryCacheTTL        time.Duration
	pubsubTimeout                     time.Duration
	oldestEventTimestamp              time.Duration
	furthestEventTimestamp            time.Duration
	inMemoryCache                     *cachev3.InMemoryCache
	metrics                           metrics.Registerer
	logger                            *zap.Logger
}

var defaultOptions = options{
	apiKeyMemoryCacheTTL:              1 * time.Minute,
	apiKeyMemoryCacheEvictionInterval: 30 * time.Second,
	featuresMemoryCacheTTL:            1 * time.Minute,
	segmentUsersMemoryCacheTTL:        1 * time.Minute,
	pubsubTimeout:                     20 * time.Second,
	// 31 days - aligns with 30-day DB retention + 1 day buffer
	oldestEventTimestamp: 744 * time.Hour,
	// 1 hour - handles legitimate clock skew while preventing malicious timestamps
	furthestEventTimestamp: 1 * time.Hour,
	logger:                 zap.NewNop(),
}

type Option func(*options)

func WithOldestEventTimestamp(d time.Duration) Option {
	return func(opts *options) {
		opts.oldestEventTimestamp = d
	}
}

func WithFurthestEventTimestamp(d time.Duration) Option {
	return func(opts *options) {
		opts.furthestEventTimestamp = d
	}
}

func WithAPIKeyMemoryCacheTTL(ttl time.Duration) Option {
	return func(opts *options) {
		opts.apiKeyMemoryCacheTTL = ttl
	}
}

func WithAPIKeyMemoryCacheEvictionInterval(interval time.Duration) Option {
	return func(opts *options) {
		opts.apiKeyMemoryCacheEvictionInterval = interval
	}
}

func WithFeaturesMemoryCacheTTL(ttl time.Duration) Option {
	return func(opts *options) {
		opts.featuresMemoryCacheTTL = ttl
	}
}

func WithSegmentUsersMemoryCacheTTL(ttl time.Duration) Option {
	return func(opts *options) {
		opts.segmentUsersMemoryCacheTTL = ttl
	}
}

func WithInMemoryCache(c *cachev3.InMemoryCache) Option {
	return func(opts *options) {
		opts.inMemoryCache = c
	}
}

func WithMetrics(r metrics.Registerer) Option {
	return func(opts *options) {
		opts.metrics = r
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type grpcGatewayService struct {
	featureClient               featureclient.Client
	accountClient               accountclient.Client
	pushClient                  pushclient.Client
	codeRefClient               coderefclient.Client
	auditLogClient              auditlogclient.Client
	autoOpsClient               autoopsclient.Client
	tagClient                   tagclient.Client
	teamClient                  teamclient.Client
	notificationClient          notificationclient.Client
	experimentClient            experimentclient.Client
	eventCounterClient          eventcounterclient.Client
	environmentClient           environmentclient.Client
	mysqlClient                 mysql.Client
	accountStorage              accstorage.AccountStorage
	goalPublisher               publisher.Publisher
	evaluationPublisher         publisher.Publisher
	userPublisher               publisher.Publisher
	featuresCache               cachev3.FeaturesCache
	featuresRedisCache          cachev3.FeaturesCache
	segmentUsersCache           cachev3.SegmentUsersCache
	segmentUsersRedisCache      cachev3.SegmentUsersCache
	environmentAPIKeyCache      cachev3.EnvironmentAPIKeyCache
	environmentAPIKeyRedisCache cachev3.EnvironmentAPIKeyCache
	apiKeyLastUsedInfoCacher    sync.Map
	flightgroup                 singleflight.Group
	opts                        *options
	logger                      *zap.Logger
}

func NewGrpcGatewayService(
	ctx context.Context,
	featureClient featureclient.Client,
	accountClient accountclient.Client,
	pushClient pushclient.Client,
	codeRefClient coderefclient.Client,
	auditLogClient auditlogclient.Client,
	autoOpsClient autoopsclient.Client,
	tagClient tagclient.Client,
	teamClient teamclient.Client,
	notificationClient notificationclient.Client,
	experimentClient experimentclient.Client,
	eventCounterClient eventcounterclient.Client,
	environmentClient environmentclient.Client,
	mysqlClient mysql.Client,
	gp publisher.Publisher,
	ep publisher.Publisher,
	up publisher.Publisher,
	redisV3Cache cache.MultiGetCache,
	opts ...Option,
) rpc.Service {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.metrics != nil {
		registerMetrics(options.metrics)
	}
	inMemoryCache := options.inMemoryCache
	if inMemoryCache == nil {
		inMemoryCache = cachev3.NewInMemoryCache(
			cachev3.WithEvictionInterval(options.apiKeyMemoryCacheEvictionInterval),
		)
	}
	s := &grpcGatewayService{
		featureClient:               featureClient,
		accountClient:               accountClient,
		pushClient:                  pushClient,
		codeRefClient:               codeRefClient,
		auditLogClient:              auditLogClient,
		autoOpsClient:               autoOpsClient,
		tagClient:                   tagClient,
		teamClient:                  teamClient,
		notificationClient:          notificationClient,
		experimentClient:            experimentClient,
		eventCounterClient:          eventCounterClient,
		environmentClient:           environmentClient,
		mysqlClient:                 mysqlClient,
		accountStorage:              accstorage.NewAccountStorage(mysqlClient),
		goalPublisher:               gp,
		evaluationPublisher:         ep,
		userPublisher:               up,
		featuresCache:               cachev3.NewFeaturesCache(inMemoryCache, options.featuresMemoryCacheTTL),
		featuresRedisCache:          cachev3.NewFeaturesCache(redisV3Cache, 0),
		segmentUsersCache:           cachev3.NewSegmentUsersCache(inMemoryCache, options.segmentUsersMemoryCacheTTL),
		segmentUsersRedisCache:      cachev3.NewSegmentUsersCache(redisV3Cache, 0),
		environmentAPIKeyCache:      cachev3.NewEnvironmentAPIKeyCache(inMemoryCache, options.apiKeyMemoryCacheTTL),
		environmentAPIKeyRedisCache: cachev3.NewEnvironmentAPIKeyCache(redisV3Cache, 0),
		apiKeyLastUsedInfoCacher:    sync.Map{},
		opts:                        &options,
		logger:                      options.logger.Named("api_grpc"),
	}

	go s.writeAPIKeyLastUsedAtCacheToDatabase(ctx)

	return s
}

func (s *grpcGatewayService) Register(server *grpc.Server) {
	gwproto.RegisterGatewayServer(server, s)
}

func (s *grpcGatewayService) Ping(ctx context.Context, req *gwproto.PingRequest) (*gwproto.PingResponse, error) {
	return &gwproto.PingResponse{Time: time.Now().Unix()}, nil
}

func (s *grpcGatewayService) Track(ctx context.Context, req *gwproto.TrackRequest) (*gwproto.TrackResponse, error) {
	ctx, span := trace.StartSpan(ctx, "bucketeerGRPCGatewayService.Track")
	defer span.End()
	if err := s.validateTrackRequest(req); err != nil {
		eventCounter.WithLabelValues(callerGatewayService, typeTrack, codeInvalidURLParams)
		s.logger.Error("Failed to validate Track request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("apiKey", obfuscateString(req.Apikey, obfuscateAPIKeyLength)),
				zap.String("tag", req.Tag),
				zap.Any("userId", req.Userid),
				zap.String("goalId", req.Goalid),
				zap.Int64("timestamp", req.Timestamp),
				zap.Float64("value", req.Value),
			)...,
		)
		return nil, err
	}
	envAPIKey, err := s.checkTrackRequest(ctx, req.Apikey)
	if err != nil {
		s.logger.Error("Failed to check Track request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("apiKey", obfuscateString(req.Apikey, obfuscateAPIKeyLength)),
				zap.String("tag", req.Tag),
				zap.Any("userId", req.Userid),
				zap.String("goalId", req.Goalid),
				zap.Int64("timestamp", req.Timestamp),
				zap.Float64("value", req.Value),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodTrack, "").Inc()
	goalEvent := &eventproto.GoalEvent{
		GoalId:    req.Goalid,
		UserId:    req.Userid,
		User:      &userproto.User{Id: req.Userid},
		Value:     req.Value,
		Timestamp: req.Timestamp,
		Tag:       req.Tag,
	}
	id, err := uuid.NewUUID()
	if err != nil {
		s.logger.Error(
			"Failed to generate uuid for goal event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", envAPIKey.Environment.Id),
				zap.String("goalId", goalEvent.GoalId),
			)...,
		)
		return nil, ErrInternal
	}
	goal, err := ptypes.MarshalAny(goalEvent)
	if err != nil {
		eventCounter.WithLabelValues(callerGatewayService, typeGoal, codeNonRepeatableError)
		s.logger.Error(
			"Failed to marshal goal event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", envAPIKey.Environment.Id),
				zap.String("goalId", goalEvent.GoalId),
			)...,
		)
		return nil, ErrInternal
	}
	event := &eventproto.Event{
		Id:            id.String(),
		Event:         goal,
		EnvironmentId: envAPIKey.Environment.Id,
	}
	if err := s.goalPublisher.Publish(ctx, event); err != nil {
		if err == publisher.ErrBadMessage {
			eventCounter.WithLabelValues(callerGatewayService, typeGoal, codeNonRepeatableError)
		} else {
			eventCounter.WithLabelValues(callerGatewayService, typeGoal, codeRepeatableError)
		}
		s.logger.Error(
			"Failed to publish goal event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", envAPIKey.Environment.Id),
				zap.String("eventId", event.Id),
				zap.String("goalId", goalEvent.GoalId),
			)...,
		)
		return nil, ErrInternal
	}
	eventCounter.WithLabelValues(callerGatewayService, typeGoal, codeOK)
	return &gwproto.TrackResponse{}, nil
}

func (s *grpcGatewayService) validateTrackRequest(req *gwproto.TrackRequest) error {
	if req.Apikey == "" {
		return ErrMissingAPIKey
	}
	if req.Userid == "" {
		return ErrUserIDRequired
	}
	if req.Goalid == "" {
		return ErrGoalIDRequired
	}
	if req.Tag == "" {
		return ErrTagRequired
	}
	if !validateTimestamp(req.Timestamp, s.opts.oldestEventTimestamp, s.opts.furthestEventTimestamp) {
		return ErrInvalidTimestamp
	}
	return nil
}

func (s *grpcGatewayService) GetEvaluations(
	ctx context.Context,
	req *gwproto.GetEvaluationsRequest,
) (*gwproto.GetEvaluationsResponse, error) {
	ctx, span := trace.StartSpan(ctx, "bucketeerGRPCGatewayService.GetEvaluations")
	defer span.End()
	startTime := time.Now()
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{accountproto.APIKey_SDK_CLIENT})
	if err != nil {
		if !isCallerContextErr(err) && !errors.Is(err, ErrInvalidAPIKey) && !errors.Is(err, ErrMissingAPIKey) {
			s.logger.Error("Failed to check GetEvaluations request",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("tag", req.Tag),
					zap.Any("user", req.User),
					zap.Any("sourceId", req.SourceId),
					zap.String("sdkVersion", req.SdkVersion),
				)...,
			)
		}
		return nil, err
	}
	projectID := envAPIKey.ProjectId
	environmentId := envAPIKey.Environment.Id
	sourceID := req.SourceId.String()
	requestTotal.WithLabelValues(envAPIKey.Environment.OrganizationId, projectID, envAPIKey.ProjectUrlCode,
		environmentId, envAPIKey.Environment.UrlCode, methodGetEvaluations, sourceID).Inc()
	defer func() {
		handledSecondsHistogram.WithLabelValues(environmentId, sourceID, methodGetEvaluations).
			Observe(time.Since(startTime).Seconds())
	}()
	if err := s.validateGetEvaluationsRequest(req); err != nil {
		s.logger.Error("Failed to validate GetEvaluations request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("projectId", projectID),
				zap.String("projectUrlCode", envAPIKey.ProjectUrlCode),
				zap.String("environmentId", environmentId),
				zap.String("tag", req.Tag),
				zap.Any("user", req.User),
				zap.Any("sourceId", req.SourceId),
				zap.String("sdkVersion", req.SdkVersion),
			)...,
		)
		evaluationsCounter.WithLabelValues(
			environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeBadRequest, sourceID).Inc()
		return nil, err
	}

	ctx, spanGetFeatures := trace.StartSpan(ctx, "bucketeerGRPCGatewayService.GetEvaluations.GetFeatures")
	f, err := s.singleflightFetch(ctx, environmentId, func(ctx context.Context) (interface{}, error) {
		return s.getFeatures(ctx, environmentId)
	})
	spanGetFeatures.End()
	if err != nil {
		if errors.Is(err, errCallerCanceled) {
			evaluationsCounter.WithLabelValues(
				environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeCanceled, sourceID).Inc()
			return nil, status.FromContextError(ctx.Err()).Err()
		}
		evaluationsCounter.WithLabelValues(
			environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeInternalError, sourceID).Inc()
		apiErrorCounter.WithLabelValues(environmentId, sourceID, methodGetEvaluations).Inc()
		return nil, err
	}
	features := f.([]*featureproto.Feature)
	filteredByTag := s.filterByTag(features, req.Tag)

	if len(features) == 0 {
		evaluationsCounter.WithLabelValues(
			environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeNoFeatures, sourceID).Inc()
		return &gwproto.GetEvaluationsResponse{
			State:             featureproto.UserEvaluations_FULL,
			Evaluations:       s.emptyUserEvaluations(),
			UserEvaluationsId: "no_evaluations",
		}, nil
	}
	ueid := evaluation.UserEvaluationsID(req.User.Id, req.User.Data, filteredByTag)
	if req.UserEvaluationsId == ueid {
		evaluationsCounter.WithLabelValues(
			environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeNone, sourceID).Inc()
		s.logger.Debug(
			"Features length when UEID is the same",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("environmentID", environmentId),
				zap.String("tag", req.Tag),
				zap.Int("featuresLength", len(features)),
				zap.Int("filteredByTagLength", len(filteredByTag)),
			)...,
		)
		return &gwproto.GetEvaluationsResponse{
			State:             featureproto.UserEvaluations_FULL,
			Evaluations:       s.emptyUserEvaluations(),
			UserEvaluationsId: ueid,
		}, nil
	}

	segmentUsersMap, err := s.getSegmentUsersMap(ctx, features, environmentId)
	if err != nil {
		if errors.Is(err, errCallerCanceled) {
			evaluationsCounter.WithLabelValues(
				environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeCanceled, sourceID).Inc()
			return nil, status.FromContextError(ctx.Err()).Err()
		}
		evaluationsCounter.WithLabelValues(
			environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeInternalError, sourceID).Inc()
		apiErrorCounter.WithLabelValues(environmentId, sourceID, methodGetEvaluations).Inc()
		s.logger.Error(
			"Failed to get segment users map",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentId),
			)...,
		)
		return nil, err
	}
	evaluator := evaluation.NewEvaluator()
	var evaluations *featureproto.UserEvaluations
	// FIXME Remove s.getEvaluations once all SDKs use UserEvaluationCondition.
	// New SDKs always use UserEvaluationCondition.
	if req.UserEvaluationCondition == nil {
		// Old evaluation requires tag to be set.
		if req.Tag == "" {
			evaluationsCounter.WithLabelValues(
				environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeBadRequest, sourceID).Inc()
			return nil, ErrTagRequired
		}
		evaluations, err = evaluator.EvaluateFeatures(
			features,
			req.User,
			segmentUsersMap,
			req.Tag,
		)
		if err != nil {
			evaluationsCounter.WithLabelValues(
				environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeInternalError, sourceID).Inc()
			apiErrorCounter.WithLabelValues(environmentId, sourceID, methodGetEvaluations).Inc()

			// Extract feature IDs for debugging dependency issues
			featureIDs := make([]string, len(features))
			archivedFeatureIDs := make([]string, 0)
			for i, f := range features {
				featureIDs[i] = f.Id
				if f.Archived {
					archivedFeatureIDs = append(archivedFeatureIDs, f.Id)
				}
			}

			s.logger.Error(
				"Failed to evaluate",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("userId", req.User.Id),
					zap.String("environmentID", environmentId),
					zap.String("tag", req.Tag),
					zap.Int("totalFeatures", len(features)),
					zap.Strings("featureIDs", featureIDs),
					zap.Strings("archivedFeatureIDs", archivedFeatureIDs),
				)...,
			)
			return nil, ErrInternal
		}
		evaluationsCounter.WithLabelValues(
			environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeOld, sourceID).Inc()
	} else {
		evaluations, err = evaluator.EvaluateFeaturesByEvaluatedAt(
			features,
			req.User,
			segmentUsersMap,
			req.UserEvaluationsId,
			req.UserEvaluationCondition.EvaluatedAt,
			req.UserEvaluationCondition.UserAttributesUpdated,
			req.Tag,
		)
		if err != nil {
			evaluationsCounter.WithLabelValues(
				environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeInternalError, sourceID).Inc()
			apiErrorCounter.WithLabelValues(environmentId, sourceID, methodGetEvaluations).Inc()

			// Extract feature IDs for debugging dependency issues
			featureIDs := make([]string, len(features))
			archivedFeatureIDs := make([]string, 0)
			for i, f := range features {
				featureIDs[i] = f.Id
				if f.Archived {
					archivedFeatureIDs = append(archivedFeatureIDs, f.Id)
				}
			}

			s.logger.Error(
				"Failed to evaluate",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("userId", req.User.Id),
					zap.String("environmentID", environmentId),
					zap.String("tag", req.Tag),
					zap.Int("totalFeatures", len(features)),
					zap.Strings("featureIDs", featureIDs),
					zap.Strings("archivedFeatureIDs", archivedFeatureIDs),
					zap.String("userEvaluationsId", req.UserEvaluationsId),
					zap.Int64("evaluatedAt", req.UserEvaluationCondition.EvaluatedAt),
					zap.Bool("userAttributesUpdated", req.UserEvaluationCondition.UserAttributesUpdated),
				)...,
			)
			return nil, ErrInternal
		}
		if evaluations.ForceUpdate {
			evaluationsCounter.WithLabelValues(
				environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeAll, sourceID).Inc()
		} else {
			evaluationsCounter.WithLabelValues(
				environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeDiff, sourceID).Inc()
		}
	}
	s.logger.Debug(
		"Features length when UEID is different",
		log.FieldsFromIncomingContext(ctx).AddFields(
			zap.String("environmentID", environmentId),
			zap.String("tag", req.Tag),
			zap.Int("featuresLength", len(features)),
			zap.Int("activeFeaturesLength", len(features)),
			zap.Int("filteredByTagLength", len(filteredByTag)),
			zap.Int("evaluationsLength", len(evaluations.Evaluations)),
		)...,
	)
	return &gwproto.GetEvaluationsResponse{
		State:             featureproto.UserEvaluations_FULL,
		Evaluations:       evaluations,
		UserEvaluationsId: ueid,
	}, nil
}

func (s *grpcGatewayService) validateGetEvaluationsRequest(req *gwproto.GetEvaluationsRequest) error {
	if req.User == nil {
		return ErrUserRequired
	}
	if req.User.Id == "" {
		return ErrUserIDRequired
	}
	return nil
}

func (s *grpcGatewayService) GetEvaluation(
	ctx context.Context,
	req *gwproto.GetEvaluationRequest,
) (*gwproto.GetEvaluationResponse, error) {
	ctx, span := trace.StartSpan(ctx, "bucketeerGRPCGatewayService.GetEvaluation")
	defer span.End()
	startTime := time.Now()
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{accountproto.APIKey_SDK_CLIENT})
	if err != nil {
		if !isCallerContextErr(err) && !errors.Is(err, ErrInvalidAPIKey) && !errors.Is(err, ErrMissingAPIKey) {
			s.logger.Error("Failed to check GetEvaluation request",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("tag", req.Tag),
					zap.Any("user", req.User),
					zap.String("featureId", req.FeatureId),
					zap.Any("sourceId", req.SourceId),
					zap.String("sdkVersion", req.SdkVersion),
				)...,
			)
		}
		return nil, err
	}
	sourceID := req.SourceId.String()
	defer func() {
		handledSecondsHistogram.WithLabelValues(envAPIKey.Environment.Id, sourceID, methodGetEvaluation).
			Observe(time.Since(startTime).Seconds())
	}()
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetEvaluation, sourceID).Inc()
	if err := s.validateGetEvaluationRequest(req); err != nil {
		s.logger.Error("Failed to validate GetEvaluation request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("projectId", envAPIKey.ProjectId),
				zap.String("projectUrlCode", envAPIKey.ProjectUrlCode),
				zap.String("environmentId", envAPIKey.Environment.Id),
				zap.String("tag", req.Tag),
				zap.Any("user", req.User),
				zap.Any("sourceId", req.SourceId),
				zap.String("sdkVersion", req.SdkVersion),
			)...,
		)
		return nil, err
	}

	ctx, spanGetFeatures := trace.StartSpan(ctx, "bucketeerGRPCGatewayService.GetEvaluation.GetFeatures")
	f, err := s.singleflightFetch(ctx, envAPIKey.Environment.Id, func(ctx context.Context) (interface{}, error) {
		return s.getFeatures(ctx, envAPIKey.Environment.Id)
	})
	spanGetFeatures.End()
	if err != nil {
		if errors.Is(err, errCallerCanceled) {
			return nil, status.FromContextError(ctx.Err()).Err()
		}
		apiErrorCounter.WithLabelValues(envAPIKey.Environment.Id, sourceID, methodGetEvaluation).Inc()
		return nil, err
	}
	fs := s.filterOutArchivedFeatures(f.([]*featureproto.Feature))
	features, err := s.getTargetFeatures(fs, req.FeatureId)
	if err != nil {
		return nil, err
	}
	segmentUsersMap, err := s.getSegmentUsersMap(ctx, features, envAPIKey.Environment.Id)
	if err != nil {
		if errors.Is(err, errCallerCanceled) {
			return nil, status.FromContextError(ctx.Err()).Err()
		}
		s.logger.Error(
			"Failed to get segment users map",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", envAPIKey.Environment.Id),
			)...,
		)
		apiErrorCounter.WithLabelValues(envAPIKey.Environment.Id, sourceID, methodGetEvaluation).Inc()
		return nil, err
	}
	evaluator := evaluation.NewEvaluator()
	evaluations, err := evaluator.EvaluateFeatures(features, req.User, segmentUsersMap, req.Tag)
	if err != nil {
		s.logger.Error(
			"Failed to evaluate features",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", envAPIKey.Environment.Id),
				zap.String("userId", req.User.Id),
				zap.String("featureId", req.FeatureId),
			)...,
		)
		apiErrorCounter.WithLabelValues(envAPIKey.Environment.Id, sourceID, methodGetEvaluation).Inc()
		return nil, ErrInternal
	}
	eval, err := s.findEvaluation(evaluations.Evaluations, req.FeatureId)
	if err != nil {
		s.logger.Error("Failed to find evaluation",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("projectId", envAPIKey.ProjectId),
				zap.String("projectUrlCode", envAPIKey.ProjectUrlCode),
				zap.String("environmentId", envAPIKey.Environment.Id),
				zap.String("tag", req.Tag),
				zap.Any("user", req.User),
				zap.Any("sourceId", req.SourceId),
				zap.String("sdkVersion", req.SdkVersion),
			)...,
		)
		return nil, err
	}
	return &gwproto.GetEvaluationResponse{
		Evaluation: eval,
	}, nil
}

func (s *grpcGatewayService) GetFeatureFlags(
	ctx context.Context,
	req *gwproto.GetFeatureFlagsRequest,
) (*gwproto.GetFeatureFlagsResponse, error) {
	ctx, span := trace.StartSpan(ctx, "bucketeerGRPCGatewayService.GetFeatureFlags")
	defer span.End()
	startTime := time.Now()
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{accountproto.APIKey_SDK_SERVER})
	if err != nil {
		if !isCallerContextErr(err) && !errors.Is(err, ErrInvalidAPIKey) && !errors.Is(err, ErrMissingAPIKey) {
			s.logger.Error("Failed to check GetFeatureFlags request",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("tag", req.Tag),
					zap.String("featureFlagsId", req.FeatureFlagsId),
					zap.Any("sourceId", req.SourceId),
					zap.String("sdkVersion", req.SdkVersion),
				)...,
			)
		}
		return nil, err
	}
	projectID := envAPIKey.ProjectId
	environmentId := envAPIKey.Environment.Id
	sourceID := req.SourceId.String()
	defer func() {
		handledSecondsHistogram.WithLabelValues(environmentId, sourceID, methodGetFeatureFlags).
			Observe(time.Since(startTime).Seconds())
	}()
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		environmentId, envAPIKey.Environment.UrlCode, methodGetFeatureFlags, sourceID).Inc()

	if err := s.validateGetFeatureFlagsRequest(req); err != nil {
		s.logger.Error("Failed to validate GetFeatureFlags request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("projectId", projectID),
				zap.String("projectUrlCode", envAPIKey.ProjectUrlCode),
				zap.String("environmentId", environmentId),
				zap.String("apiKey", obfuscateString(envAPIKey.ApiKey.Id, obfuscateAPIKeyLength)),
				zap.Any("sourceId", req.SourceId),
				zap.String("sdkVersion", req.SdkVersion),
			)...,
		)
		getFeatureFlagsCounter.WithLabelValues(projectID, envAPIKey.ProjectUrlCode,
			environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeBadRequest).Inc()
		return nil, err
	}
	ctx, spanGetFeatures := trace.StartSpan(ctx, "bucketeerGRPCGatewayService.GetFeatureFlags.GetFeatures")
	f, err := s.singleflightFetch(ctx, environmentId, func(ctx context.Context) (interface{}, error) {
		return s.getFeatures(ctx, environmentId)
	})
	spanGetFeatures.End()
	if err != nil {
		if errors.Is(err, errCallerCanceled) {
			getFeatureFlagsCounter.WithLabelValues(projectID, envAPIKey.ProjectUrlCode,
				environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeCanceled).Inc()
			return nil, status.FromContextError(ctx.Err()).Err()
		}
		getFeatureFlagsCounter.WithLabelValues(projectID, envAPIKey.ProjectUrlCode,
			environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeInternalError).Inc()
		apiErrorCounter.WithLabelValues(environmentId, sourceID, methodGetFeatureFlags).Inc()
		return nil, err
	}
	// Filter flags by tag if needed
	features := f.([]*featureproto.Feature)
	var targetFeatures []*featureproto.Feature
	if req.Tag == "" {
		targetFeatures = features
	} else {
		targetFeatures = s.filterByTag(features, req.Tag)
	}
	now := time.Now()
	if len(targetFeatures) == 0 {
		getFeatureFlagsCounter.WithLabelValues(projectID, envAPIKey.ProjectUrlCode,
			environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeNoFeatures).Inc()
		return &gwproto.GetFeatureFlagsResponse{
			FeatureFlagsId:         "",
			Features:               []*featureproto.Feature{},
			ArchivedFeatureFlagIds: make([]string, 0),
			RequestedAt:            now.Unix(),
		}, nil
	}
	// We don't include archived flags when generating the Feature Flag IDs
	filteredArchivedFlags := s.filterOutArchivedFeatures(targetFeatures)
	ffID := evaluation.GenerateFeaturesID(filteredArchivedFlags)
	// Return an empty response because nothing changed.
	// We preserve req.RequestedAt (clamped to now to handle clock skew)
	// instead of advancing to now.Unix(), so the SDK's time cursor stays
	// anchored to when it last received actual data.
	if req.FeatureFlagsId == ffID {
		getFeatureFlagsCounter.WithLabelValues(projectID, envAPIKey.ProjectUrlCode,
			environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeNone).Inc()
		s.logger.Debug(
			"Feature Flags ID is the same",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("environmentID", environmentId),
				zap.String("tag", req.Tag),
				zap.Int("featuresLength", len(targetFeatures)),
				zap.Int("targetFeaturesLength", len(targetFeatures)),
			)...,
		)
		return &gwproto.GetFeatureFlagsResponse{
			FeatureFlagsId:         ffID,
			Features:               []*featureproto.Feature{},
			ArchivedFeatureFlagIds: make([]string, 0),
			RequestedAt:            min(req.RequestedAt, now.Unix()),
		}, nil
	}
	s.logger.Debug(
		"Feature Flags ID is different",
		log.FieldsFromIncomingContext(ctx).AddFields(
			zap.String("environmentID", environmentId),
			zap.String("tag", req.Tag),
			zap.Int("featuresLength", len(targetFeatures)),
			zap.Int("targetFeaturesLength", len(targetFeatures)),
		)...,
	)
	// Return all flags when: first request, cache older than 30 days,
	// or future requestedAt (clock skew) to avoid missing updates in the Diff filter.
	if req.FeatureFlagsId == "" ||
		req.RequestedAt < now.Unix()-secondsToReturnAllFlags ||
		req.RequestedAt > now.Unix() {
		getFeatureFlagsCounter.WithLabelValues(projectID, envAPIKey.ProjectUrlCode,
			environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeAll).Inc()
		return &gwproto.GetFeatureFlagsResponse{
			FeatureFlagsId:         ffID,
			Features:               filteredArchivedFlags,
			ArchivedFeatureFlagIds: make([]string, 0),
			RequestedAt:            now.Unix(),
			ForceUpdate:            true,
		}, nil
	}
	// Diff path: only reached when req.RequestedAt is within [now-30days, now],
	// so req.RequestedAt can be used directly without clamping.
	updatedFeatures := make([]*featureproto.Feature, 0, len(targetFeatures))
	archivedIDs := make([]string, 0)
	for _, feature := range targetFeatures {
		if s.isArchivedBeforeLastThirtyDays(feature) {
			archivedIDs = append(archivedIDs, feature.Id)
			continue
		}
		if feature.UpdatedAt >= req.RequestedAt {
			updatedFeatures = append(updatedFeatures, feature)
		}
	}
	getFeatureFlagsCounter.WithLabelValues(projectID, envAPIKey.ProjectUrlCode,
		environmentId, envAPIKey.Environment.UrlCode, req.Tag, codeDiff).Inc()
	return &gwproto.GetFeatureFlagsResponse{
		FeatureFlagsId:         ffID,
		Features:               updatedFeatures,
		ArchivedFeatureFlagIds: archivedIDs,
		RequestedAt:            now.Unix(),
		ForceUpdate:            false,
	}, nil
}

func (s *grpcGatewayService) validateGetFeatureFlagsRequest(req *gwproto.GetFeatureFlagsRequest) error {
	if req.SourceId == eventproto.SourceId_UNKNOWN {
		return ErrSourceIDRequired
	}
	if req.SdkVersion == "" {
		return ErrSDKVersionRequired
	}
	return nil
}

// To keep the response size small, the feature flags archived more than 30 days are excluded
func (s *grpcGatewayService) isArchivedBeforeLastThirtyDays(feature *featureproto.Feature) bool {
	return feature.Archived && feature.UpdatedAt > time.Now().Unix()-secondsToReturnAllFlags
}

func (s *grpcGatewayService) GetSegmentUsers(
	ctx context.Context,
	req *gwproto.GetSegmentUsersRequest,
) (*gwproto.GetSegmentUsersResponse, error) {
	ctx, span := trace.StartSpan(ctx, "bucketeerGRPCGatewayService.GetSegmentUsers")
	defer span.End()
	startTime := time.Now()
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{accountproto.APIKey_SDK_SERVER})
	if err != nil {
		if !isCallerContextErr(err) && !errors.Is(err, ErrInvalidAPIKey) && !errors.Is(err, ErrMissingAPIKey) {
			s.logger.Error("Failed to check GetSegmentUsers request",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.Strings("segmentIds", req.SegmentIds),
					zap.Any("sourceId", req.SourceId),
					zap.String("sdkVersion", req.SdkVersion),
				)...,
			)
		}
		return nil, err
	}
	projectID := envAPIKey.ProjectId
	environmentId := envAPIKey.Environment.Id
	sourceID := req.SourceId.String()
	defer func() {
		handledSecondsHistogram.WithLabelValues(environmentId, sourceID, methodGetSegmentUsers).
			Observe(time.Since(startTime).Seconds())
	}()
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		environmentId, envAPIKey.Environment.UrlCode, methodGetSegmentUsers, sourceID).Inc()

	if err := s.validateGetSegmentUsersRequest(req); err != nil {
		s.logger.Error("Failed to validate GetSegmentUsers request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("projectId", projectID),
				zap.String("projectUrlCode", envAPIKey.ProjectUrlCode),
				zap.String("environmentId", environmentId),
				zap.String("apiKey", obfuscateString(envAPIKey.ApiKey.Id, obfuscateAPIKeyLength)),
				zap.Strings("segmentIds", req.SegmentIds),
				zap.Any("sourceId", req.SourceId),
				zap.String("sdkVersion", req.SdkVersion),
			)...,
		)
		getSegmentUsersCounter.WithLabelValues(projectID, envAPIKey.ProjectUrlCode, environmentId,
			envAPIKey.Environment.UrlCode, sourceID, req.GetSdkVersion(), codeBadRequest).Inc()
		return nil, err
	}

	// Get the feature flags from the cache
	ctx, spanGetFeatures := trace.StartSpan(ctx, "bucketeerGRPCGatewayService.GetSegmentUsers.GetFeatures")
	f, err := s.singleflightFetch(ctx, environmentId, func(ctx context.Context) (interface{}, error) {
		return s.getFeatures(ctx, environmentId)
	})
	spanGetFeatures.End()
	if err != nil {
		if errors.Is(err, errCallerCanceled) {
			getSegmentUsersCounter.WithLabelValues(
				projectID, envAPIKey.ProjectUrlCode, environmentId,
				envAPIKey.Environment.UrlCode, sourceID, req.GetSdkVersion(), codeCanceled).Inc()
			return nil, status.FromContextError(ctx.Err()).Err()
		}
		getSegmentUsersCounter.WithLabelValues(
			projectID, envAPIKey.ProjectUrlCode, environmentId,
			envAPIKey.Environment.UrlCode, sourceID, req.GetSdkVersion(), codeInternalError).Inc()
		apiErrorCounter.WithLabelValues(environmentId, sourceID, methodGetSegmentUsers).Inc()
		return nil, err
	}

	// Return an empty response when there is no feature flags
	targetFeatures := s.filterOutArchivedFeatures(f.([]*featureproto.Feature))
	if len(targetFeatures) == 0 {
		getSegmentUsersCounter.WithLabelValues(projectID, envAPIKey.ProjectUrlCode, environmentId,
			envAPIKey.Environment.UrlCode, sourceID, req.GetSdkVersion(), codeNoFeatures).Inc()
		return s.emptyGetSegmentUsersResponse()
	}

	// Get the segment IDs in used
	targetSegmentIDs := make([]string, 0)
	for _, feature := range targetFeatures {
		f := &featuredomain.Feature{Feature: feature}
		ids := f.ListSegmentIDs()
		if len(ids) > 0 {
			targetSegmentIDs = append(targetSegmentIDs, ids...)
		}
	}

	// Return an empty response when there is no segments
	if len(targetSegmentIDs) == 0 {
		getSegmentUsersCounter.WithLabelValues(projectID, envAPIKey.ProjectUrlCode, environmentId,
			envAPIKey.Environment.UrlCode, sourceID, req.GetSdkVersion(), codeNoSegments).Inc()
		return s.emptyGetSegmentUsersResponse()
	}

	// Get the segment users
	targetSegmentUsers := make([]*featureproto.SegmentUsers, 0, len(targetSegmentIDs))
	for _, sID := range targetSegmentIDs {
		ctx, spanGetSegmentUsers := trace.StartSpan(
			ctx,
			"bucketeerGRPCGatewayService.GetSegmentUsers.GetSegmentUsersBySegmentID",
		)
		su, err := s.singleflightFetch(
			ctx,
			s.segmentFlightID(environmentId, sID),
			func(ctx context.Context) (interface{}, error) {
				return s.getSegmentUsersBySegmentID(ctx, sID, environmentId)
			},
		)
		spanGetSegmentUsers.End()
		if err != nil {
			if errors.Is(err, errCallerCanceled) {
				getSegmentUsersCounter.WithLabelValues(
					projectID, envAPIKey.ProjectUrlCode, environmentId,
					envAPIKey.Environment.UrlCode, sourceID, req.GetSdkVersion(), codeCanceled).Inc()
				return nil, status.FromContextError(ctx.Err()).Err()
			}
			getSegmentUsersCounter.WithLabelValues(
				projectID, envAPIKey.ProjectUrlCode, environmentId,
				envAPIKey.Environment.UrlCode, sourceID, req.GetSdkVersion(), codeInternalError).Inc()
			apiErrorCounter.WithLabelValues(environmentId, sourceID, methodGetSegmentUsers).Inc()
			return nil, err
		}
		segmentUsers := su.(*featureproto.SegmentUsers)
		targetSegmentUsers = append(targetSegmentUsers, segmentUsers)
	}

	now := time.Now().Unix()
	// Return all segments when: cache older than 30 days,
	// or future requestedAt (clock skew) to avoid missing updates in the Diff filter.
	if req.RequestedAt < now-secondsToReturnAllFlags || req.RequestedAt > now {
		getSegmentUsersCounter.WithLabelValues(projectID, envAPIKey.ProjectUrlCode, environmentId,
			envAPIKey.Environment.UrlCode, sourceID, req.GetSdkVersion(), codeAll).Inc()
		return &gwproto.GetSegmentUsersResponse{
			SegmentUsers:      targetSegmentUsers,
			DeletedSegmentIds: make([]string, 0),
			RequestedAt:       now,
			ForceUpdate:       true,
		}, nil
	}

	// Find deleted segments
	deletedSegmentIDs := make([]string, 0)
	for _, id := range req.SegmentIds {
		if !contains(targetSegmentIDs, id) {
			deletedSegmentIDs = append(deletedSegmentIDs, id)
		}
	}
	// Diff path: only reached when req.RequestedAt is within [now-30days, now],
	// so req.RequestedAt can be used directly without clamping.
	updatedSegments := make([]*featureproto.SegmentUsers, 0, len(targetSegmentUsers))
	for _, su := range targetSegmentUsers {
		if su.UpdatedAt >= req.RequestedAt {
			updatedSegments = append(updatedSegments, su)
		}
	}

	// When nothing changed, preserve the SDK's requestedAt so its time cursor
	// stays anchored to the last actual data sync.
	if len(updatedSegments) == 0 && len(deletedSegmentIDs) == 0 {
		getSegmentUsersCounter.WithLabelValues(projectID, envAPIKey.ProjectUrlCode, environmentId,
			envAPIKey.Environment.UrlCode, sourceID, req.GetSdkVersion(), codeNone).Inc()
		return &gwproto.GetSegmentUsersResponse{
			SegmentUsers:      updatedSegments,
			DeletedSegmentIds: deletedSegmentIDs,
			RequestedAt:       req.RequestedAt,
			ForceUpdate:       false,
		}, nil
	}
	getSegmentUsersCounter.WithLabelValues(projectID, envAPIKey.ProjectUrlCode, environmentId,
		envAPIKey.Environment.UrlCode, sourceID, req.GetSdkVersion(), codeDiff).Inc()
	return &gwproto.GetSegmentUsersResponse{
		SegmentUsers:      updatedSegments,
		DeletedSegmentIds: deletedSegmentIDs,
		RequestedAt:       now,
		ForceUpdate:       false,
	}, nil
}

func (s *grpcGatewayService) validateGetSegmentUsersRequest(req *gwproto.GetSegmentUsersRequest) error {
	if req.SourceId == eventproto.SourceId_UNKNOWN {
		return ErrSourceIDRequired
	}
	if req.SdkVersion == "" {
		return ErrSDKVersionRequired
	}
	return nil
}

func (s *grpcGatewayService) emptyGetSegmentUsersResponse() (*gwproto.GetSegmentUsersResponse, error) {
	return &gwproto.GetSegmentUsersResponse{
		SegmentUsers:      []*featureproto.SegmentUsers{},
		DeletedSegmentIds: make([]string, 0),
		RequestedAt:       time.Now().Unix(),
		ForceUpdate:       true,
	}, nil
}

func (s *grpcGatewayService) getTargetFeatures(fs []*featureproto.Feature, id string) ([]*featureproto.Feature, error) {
	feature, err := s.findFeature(fs, id)
	if err != nil {
		return nil, err
	}
	// Check if the flag depends on other flags.
	// If not, we return only the target flag
	df := &featuredomain.Feature{Feature: feature}
	if len(df.FeatureIDsDependsOn()) == 0 {
		return []*featureproto.Feature{feature}, nil
	}
	// Otherwise, we evaluate all features here to avoid complex logic.
	evaluator := evaluation.NewEvaluator()
	return evaluator.GetPrerequisiteDownwards([]*featureproto.Feature{feature}, fs)
}

func (*grpcGatewayService) findFeature(fs []*featureproto.Feature, id string) (*featureproto.Feature, error) {
	for _, f := range fs {
		if f.Id == id {
			return f, nil
		}
	}
	return nil, ErrFeatureNotFound
}

func (*grpcGatewayService) findEvaluation(
	evals []*featureproto.Evaluation,
	id string,
) (*featureproto.Evaluation, error) {
	for _, e := range evals {
		if e.FeatureId == id {
			return e, nil
		}
	}
	return nil, ErrEvaluationNotFound
}

func (s *grpcGatewayService) validateGetEvaluationRequest(req *gwproto.GetEvaluationRequest) error {
	if req.Tag == "" {
		return ErrTagRequired
	}
	if req.User == nil {
		return ErrUserRequired
	}
	if req.User.Id == "" {
		return ErrUserIDRequired
	}
	if req.FeatureId == "" {
		return ErrFeatureIDRequired
	}
	return nil
}

// singleflightFetch shares concurrent fetches via singleflight while ensuring
// that a single caller's context cancellation does not abort the work for the
// other waiters. The shared work runs with a context derived from the first
// caller's context (preserving values such as logging fields and trace
// metadata) but with cancellation detached and a fresh timeout applied. Each
// caller still honors its own ctx.Done(), so a slow caller can give up
// locally without affecting the in-flight call.
//
// This avoids the thundering-herd cancellation that occurs when many callers
// fan-in on a cache-miss path: previously, if the first caller's gRPC context
// was canceled, the shared downstream call (DB / feature service) was canceled
// too and every waiter received "context canceled".
//
// When the caller's own context is canceled while waiting, the returned error
// wraps errCallerCanceled so call sites can distinguish a client disconnect
// (not a server-side failure) from a downstream error and avoid bumping
// internal-error metrics for it.
func (s *grpcGatewayService) singleflightFetch(
	ctx context.Context,
	key string,
	fn func(ctx context.Context) (interface{}, error),
) (interface{}, error) {
	ch := s.flightgroup.DoChan(key, func() (interface{}, error) {
		innerCtx, cancel := context.WithTimeout(
			context.WithoutCancel(ctx),
			singleflightFetchTimeout,
		)
		defer cancel()
		return fn(innerCtx)
	})
	select {
	case res := <-ch:
		return res.Val, res.Err
	case <-ctx.Done():
		return nil, fmt.Errorf("%w: %w", errCallerCanceled, status.FromContextError(ctx.Err()).Err())
	}
}

// isCallerContextErr reports whether err is one of the gateway sentinels
// produced when the caller's request context terminated (canceled or its
// deadline was exceeded). Used by the public-API entry points to suppress
// noisy "Failed to check ... request" logs for client-side disconnects.
func isCallerContextErr(err error) bool {
	return errors.Is(err, ErrContextCanceled) || errors.Is(err, ErrContextDeadlineExceeded)
}

// translateCallerCanceledErr converts a singleflightFetch caller-cancellation
// error into the package's well-known sentinels, preserving the underlying
// gRPC status (Canceled vs DeadlineExceeded) so the SDK sees the right code
// and upstream log-suppression keeps working. Errors that do not wrap
// errCallerCanceled are returned unchanged.
func translateCallerCanceledErr(ctx context.Context, err error) error {
	if !errors.Is(err, errCallerCanceled) {
		return err
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return ErrContextDeadlineExceeded
	}
	return ErrContextCanceled
}

func (s *grpcGatewayService) getFeatures(
	ctx context.Context,
	environmentId string,
) ([]*featureproto.Feature, error) {
	// L1: in-memory cache
	fs, err := getFeaturesFromCache(
		environmentId,
		s.featuresCache,
		callerGatewayService,
		cacheLayerInMemory,
	)
	if err == nil {
		return fs.Features, nil
	}
	// L2: Redis cache (kept warm by batch cacher)
	fs, err = getFeaturesFromCache(
		environmentId,
		s.featuresRedisCache,
		callerGatewayService,
		cacheLayerExternal,
	)
	if err == nil {
		putFeaturesCache(ctx, fs, environmentId, s.featuresCache, s.logger)
		return fs.Features, nil
	}
	// L3: feature service (DB)
	s.logger.Warn(
		"No cached data for Features",
		log.FieldsFromIncomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentID", environmentId),
		)...,
	)
	features, err := s.listFeatures(ctx, environmentId)
	if err != nil {
		s.logger.Error(
			"Failed to retrieve features from storage",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentId),
			)...,
		)
		return nil, ErrInternal
	}
	return features, nil
}

func (s *grpcGatewayService) listFeatures(
	ctx context.Context,
	environmentId string,
) ([]*featureproto.Feature, error) {
	features := []*featureproto.Feature{}
	cursor := ""
	for {
		resp, err := s.featureClient.ListFeatures(ctx, &featureproto.ListFeaturesRequest{
			PageSize:      listRequestSize,
			Cursor:        cursor,
			EnvironmentId: environmentId,
		})
		if err != nil {
			return nil, err
		}
		for _, f := range resp.Features {
			ff := featuredomain.Feature{Feature: f}
			if ff.IsDisabledAndOffVariationEmpty() {
				continue
			}
			// To keep the cache size small, we exclude feature flags archived more than thirty days ago.
			if ff.IsArchivedBeforeLastThirtyDays() {
				continue
			}
			features = append(features, f)
		}
		featureSize := len(resp.Features)
		if featureSize == 0 || featureSize < listRequestSize {
			return features, nil
		}
		cursor = resp.Cursor
	}
}

func getFeaturesFromCache(
	environmentId string,
	c cachev3.FeaturesCache,
	caller, layer string,
) (*featureproto.Features, error) {
	features, err := c.Get(environmentId)
	if err == nil {
		cacheCounter.WithLabelValues(caller, typeFeatures, layer, codeHit).Inc()
		return features, nil
	}
	cacheCounter.WithLabelValues(caller, typeFeatures, layer, codeMiss).Inc()
	return nil, err
}

func putFeaturesCache(
	ctx context.Context,
	features *featureproto.Features,
	environmentId string,
	featuresCache cachev3.FeaturesCache,
	logger *zap.Logger,
) {
	if err := featuresCache.Put(features, environmentId); err != nil {
		logger.Error(
			"Failed to cache features",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentId),
			)...,
		)
	}
}

func (s *grpcGatewayService) getSegmentUsersMap(
	ctx context.Context,
	features []*featureproto.Feature,
	environmentId string,
) (map[string][]*featureproto.SegmentUser, error) {
	evaluator := evaluation.NewEvaluator()
	mapIDs := make(map[string]struct{})
	for _, f := range features {
		for _, id := range evaluator.ListSegmentIDs(f) {
			mapIDs[id] = struct{}{}
		}
	}
	segmentUsersMap, err := s.listSegmentUsers(ctx, mapIDs, environmentId)
	if err != nil {
		s.logger.Error(
			"Failed to list segments",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentId),
			)...,
		)
		return nil, err
	}
	return segmentUsersMap, nil
}

func (s *grpcGatewayService) listSegmentUsers(
	ctx context.Context,
	mapSegmentIDs map[string]struct{},
	environmentId string,
) (map[string][]*featureproto.SegmentUser, error) {
	if len(mapSegmentIDs) == 0 {
		return nil, nil
	}
	users := make(map[string][]*featureproto.SegmentUser)
	for segmentID := range mapSegmentIDs {
		su, err := s.singleflightFetch(
			ctx,
			s.segmentFlightID(environmentId, segmentID),
			func(ctx context.Context) (interface{}, error) {
				return s.getSegmentUsersBySegmentID(ctx, segmentID, environmentId)
			},
		)
		if err != nil {
			return nil, err
		}
		segmentUsers := su.(*featureproto.SegmentUsers)
		users[segmentID] = segmentUsers.Users
	}
	return users, nil
}

func (s *grpcGatewayService) segmentFlightID(environmentId, segmentID string) string {
	return environmentId + ":" + segmentID
}

func (s *grpcGatewayService) getSegmentUsersBySegmentID(
	ctx context.Context,
	segmentID, environmentId string,
) (*featureproto.SegmentUsers, error) {
	// L1: in-memory cache
	segmentUsers, err := getSegmentUsersFromCache(
		segmentID,
		environmentId,
		s.segmentUsersCache,
		callerGatewayService,
		cacheLayerInMemory,
	)
	if err == nil {
		return segmentUsers, nil
	}
	// L2: Redis cache (kept warm by batch cacher)
	segmentUsers, err = getSegmentUsersFromCache(
		segmentID,
		environmentId,
		s.segmentUsersRedisCache,
		callerGatewayService,
		cacheLayerExternal,
	)
	if err == nil {
		putSegmentUsersCache(ctx, segmentUsers, environmentId, s.segmentUsersCache, s.logger)
		return segmentUsers, nil
	}
	// L3: feature service (DB)
	s.logger.Warn(
		"No cached data for SegmentUsers",
		log.FieldsFromIncomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentID", environmentId),
			zap.String("segmentId", segmentID),
		)...,
	)
	req := &featureproto.ListSegmentUsersRequest{
		SegmentId:     segmentID,
		EnvironmentId: environmentId,
	}
	res, err := s.featureClient.ListSegmentUsers(ctx, req)
	if err != nil {
		s.logger.Error(
			"Failed to retrieve segment users from database",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentId),
				zap.String("segmentId", segmentID),
			)...,
		)
		return nil, ErrInternal
	}
	reqGet := &featureproto.GetSegmentRequest{
		Id:            segmentID,
		EnvironmentId: environmentId,
	}
	respGet, err := s.featureClient.GetSegment(ctx, reqGet)
	if err != nil {
		s.logger.Error(
			"Failed to retrieve segment from database",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentId),
				zap.String("segmentId", segmentID),
			)...,
		)
		return nil, ErrInternal
	}
	segmentUsers = &featureproto.SegmentUsers{
		SegmentId: segmentID,
		Users:     res.Users,
		UpdatedAt: respGet.Segment.UpdatedAt,
	}
	putSegmentUsersCache(ctx, segmentUsers, environmentId, s.segmentUsersCache, s.logger)
	return segmentUsers, nil
}

func getSegmentUsersFromCache(
	segmentID, environmentId string,
	c cachev3.SegmentUsersCache,
	caller, layer string,
) (*featureproto.SegmentUsers, error) {
	segmentUsers, err := c.Get(segmentID, environmentId)
	if err == nil {
		cacheCounter.WithLabelValues(caller, typeSegmentUsers, layer, codeHit).Inc()
		return segmentUsers, nil
	}
	cacheCounter.WithLabelValues(caller, typeSegmentUsers, layer, codeMiss).Inc()
	return nil, err
}

func putSegmentUsersCache(
	ctx context.Context,
	segmentUsers *featureproto.SegmentUsers,
	environmentId string,
	segmentUsersCache cachev3.SegmentUsersCache,
	logger *zap.Logger,
) {
	if err := segmentUsersCache.Put(segmentUsers, environmentId); err != nil {
		logger.Error(
			"Failed to cache segment users",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentId),
				zap.String("segmentId", segmentUsers.SegmentId),
			)...,
		)
	}
}

func (s *grpcGatewayService) RegisterEvents(
	ctx context.Context,
	req *gwproto.RegisterEventsRequest,
) (*gwproto.RegisterEventsResponse, error) {
	ctx, span := trace.StartSpan(ctx, "bucketeerGRPCGatewayService.RegisterEvents")
	defer span.End()
	startTime := time.Now()
	allowedRoles := []accountproto.APIKey_Role{accountproto.APIKey_SDK_CLIENT, accountproto.APIKey_SDK_SERVER}
	envAPIKey, err := s.checkRequest(ctx, allowedRoles)
	if err != nil {
		if !isCallerContextErr(err) && !errors.Is(err, ErrInvalidAPIKey) && !errors.Is(err, ErrMissingAPIKey) {
			s.logger.Error("Failed to check RegisterEvents request",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.Any("events", req.Events),
					zap.Any("sourceId", req.SourceId),
					zap.String("sdkVersion", req.SdkVersion),
				)...,
			)
		}
		return nil, err
	}
	sourceID := req.SourceId.String()
	defer func() {
		handledSecondsHistogram.WithLabelValues(envAPIKey.Environment.Id, sourceID, methodRegisterEvents).
			Observe(time.Since(startTime).Seconds())
	}()
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodRegisterEvents, sourceID).Inc()
	if len(req.Events) == 0 {
		s.logger.Error("Failed to validate RegisterEvents request. Missing events.",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", envAPIKey.Environment.Id),
				zap.String("apiKey", obfuscateString(envAPIKey.ApiKey.Id, obfuscateAPIKeyLength)),
				zap.Any("sourceId", req.SourceId),
				zap.String("sdkVersion", req.SdkVersion),
			)...,
		)
		return nil, ErrMissingEvents
	}
	errs := make(map[string]*gwproto.RegisterEventsResponse_Error)
	goalMessages := make([]publisher.Message, 0)
	evaluationMessages := make([]publisher.Message, 0)
	metricsEvents := make([]*eventproto.MetricsEvent, 0)
	publish := func(
		p publisher.Publisher,
		messages []publisher.Message,
		typ string,
	) map[string]*gwproto.RegisterEventsResponse_Error {
		errs := make(map[string]*gwproto.RegisterEventsResponse_Error)
		multiErrs := p.PublishMulti(ctx, messages)
		var repeatableErrors, nonRepeateableErrors float64
		for id, err := range multiErrs {
			retriable := err != publisher.ErrBadMessage
			if retriable {
				repeatableErrors++
			} else {
				nonRepeateableErrors++
			}
			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				s.logger.Error(
					"Failed to publish event",
					log.FieldsFromIncomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentID", envAPIKey.Environment.Id),
						zap.String("eventType", typ),
						zap.String("id", id),
						zap.Any("sourceId", req.SourceId),
					)...,
				)
			}
			errs[id] = &gwproto.RegisterEventsResponse_Error{
				Retriable: retriable,
				Message:   "Failed to publish event",
			}
		}
		eventCounter.WithLabelValues(callerGatewayService, typ, codeNonRepeatableError).Add(nonRepeateableErrors)
		eventCounter.WithLabelValues(callerGatewayService, typ, codeRepeatableError).Add(repeatableErrors)
		eventCounter.WithLabelValues(callerGatewayService, typ, codeOK).Add(float64(len(messages) - len(multiErrs)))
		return errs
	}
	for i, event := range req.Events {
		event.EnvironmentId = envAPIKey.Environment.Id
		if event.Id == "" {
			return nil, ErrMissingEventID
		}
		validator := newEventValidator(event, s.opts.oldestEventTimestamp, s.opts.furthestEventTimestamp, s.logger)
		if validator == nil {
			errs[event.Id] = &gwproto.RegisterEventsResponse_Error{
				Retriable: false,
				Message:   "Invalid message type",
			}
			eventCounter.WithLabelValues(callerGatewayService, typeUnknown, codeInvalidType).Inc()
			s.logger.Warn("Received invalid type event",
				zap.String("apiKey", obfuscateString(envAPIKey.ApiKey.Id, obfuscateAPIKeyLength)),
				zap.String("projectID", envAPIKey.ProjectId),
				zap.String("eventID", event.Id),
				zap.String("environmentId", event.EnvironmentId),
				zap.Any("event", event.Event),
				zap.Any("sourceId", req.SourceId),
			)
			continue
		}
		if ptypes.Is(event.Event, grpcGoalEvent) {
			errorCode, err := validator.validate(ctx)
			if err != nil {
				eventCounter.WithLabelValues(callerGatewayService, typeGoal, errorCode).Inc()
				errs[event.Id] = &gwproto.RegisterEventsResponse_Error{
					Retriable: false,
					Message:   err.Error(),
				}
				continue
			}
			goalMessages = append(goalMessages, event)
			continue
		}
		if ptypes.Is(event.Event, grpcEvaluationEvent) {
			errorCode, err := validator.validate(ctx)
			if err != nil {
				eventCounter.WithLabelValues(callerGatewayService, typeEvaluation, errorCode).Inc()
				errs[event.Id] = &gwproto.RegisterEventsResponse_Error{
					Retriable: false,
					Message:   err.Error(),
				}
				continue
			}
			evaluationMessages = append(evaluationMessages, event)
			// Report evaluation events with error reasons for monitoring.
			if evValidator, ok := validator.(*eventEvaluationValidator); ok &&
				evValidator.lastUnmarshaledEvent != nil &&
				isEvaluationEventErrorReason(evValidator.lastUnmarshaledEvent.Reason) {
				ev := evValidator.lastUnmarshaledEvent
				evaluationEventErrorReasonCounter.WithLabelValues(
					envAPIKey.ProjectId,
					envAPIKey.Environment.UrlCode,
					ev.Tag,
					ev.Reason.Type.String(),
					ev.SdkVersion,
					ev.SourceId.String(),
				).Inc()
			}
			continue
		}
		if ptypes.Is(event.Event, grpcMetricsEvent) {
			errorCode, err := validator.validate(ctx)
			if err != nil {
				eventCounter.WithLabelValues(callerGatewayService, typeMetrics, errorCode).Inc()
				errs[event.Id] = &gwproto.RegisterEventsResponse_Error{
					Retriable: false,
					Message:   err.Error(),
				}
				continue
			}
			m := &eventproto.MetricsEvent{}
			if err := ptypes.UnmarshalAny(req.Events[i].Event, m); err != nil {
				eventCounter.WithLabelValues(callerGatewayService, typeMetrics, codeUnmarshalFailed).Inc()
				errs[event.Id] = &gwproto.RegisterEventsResponse_Error{
					Retriable: false,
					Message:   err.Error(),
				}
				continue
			}
			metricsEvents = append(metricsEvents, m)
		}
	}
	// MetricsEvents are saved asynchronously for performance, since there is no user impact even if they are lost.
	s.saveMetricsEventsAsync(metricsEvents, envAPIKey.ProjectId, envAPIKey.Environment.UrlCode)
	goalErrors := publish(s.goalPublisher, goalMessages, typeGoal)
	evalErrors := publish(s.evaluationPublisher, evaluationMessages, typeEvaluation)
	errs = s.mergeMaps(errs, goalErrors, evalErrors)
	if len(errs) > 0 {
		if s.containsInvalidTimestampError(errs) {
			eventCounter.WithLabelValues(callerGatewayService, typeRegisterEvent, codeInvalidTimestampRequest).Inc()
		}
	} else {
		eventCounter.WithLabelValues(callerGatewayService, typeRegisterEvent, codeOK).Inc()
	}
	return &gwproto.RegisterEventsResponse{Errors: errs}, nil
}

func (s *grpcGatewayService) containsInvalidTimestampError(errs map[string]*gwproto.RegisterEventsResponse_Error) bool {
	for _, v := range errs {
		if v.Message == errInvalidTimestamp.Error() {
			return true
		}
	}
	return false
}

func (*grpcGatewayService) mergeMaps(
	maps ...map[string]*gwproto.RegisterEventsResponse_Error,
) map[string]*gwproto.RegisterEventsResponse_Error {
	result := make(map[string]*gwproto.RegisterEventsResponse_Error)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func (s *grpcGatewayService) checkTrackRequest(
	ctx context.Context,
	apiKey string,
) (*accountproto.EnvironmentAPIKey, error) {
	if isContextCanceled(ctx) {
		s.logger.Warn(
			"Request was canceled",
			log.FieldsFromIncomingContext(ctx)...,
		)
		return nil, ErrContextCanceled
	}
	envAPIKey, err := s.getEnvironmentAPIKey(ctx, apiKey)
	if err != nil {
		s.logger.Error("Failed to get environment API key",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("apiKey", obfuscateString(apiKey, obfuscateAPIKeyLength)),
			)...,
		)
		return nil, err
	}
	if err := checkEnvironmentAPIKey(envAPIKey, []accountproto.APIKey_Role{accountproto.APIKey_SDK_CLIENT}); err != nil {
		s.logger.Error("Failed to check environment API key",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("apiKey", obfuscateString(envAPIKey.ApiKey.Id, obfuscateAPIKeyLength)),
			)...,
		)
		return nil, err
	}
	return envAPIKey, nil
}

func (s *grpcGatewayService) checkRequest(
	ctx context.Context,
	roles []accountproto.APIKey_Role,
) (*accountproto.EnvironmentAPIKey, error) {
	if isContextCanceled(ctx) {
		return nil, ErrContextCanceled
	}
	apiKey, err := s.extractAPIKey(ctx)
	if err != nil {
		return nil, err
	}
	envAPIKey, err := s.getEnvironmentAPIKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	if err := checkEnvironmentAPIKey(envAPIKey, roles); err != nil {
		s.logger.Error("Failed to check environment API key",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("apiKey", obfuscateString(envAPIKey.ApiKey.Id, obfuscateAPIKeyLength)),
			)...,
		)
		return nil, err
	}

	go s.cacheAPIKeyLastUsedAt(envAPIKey, time.Now().Unix())

	return envAPIKey, nil
}

func (s *grpcGatewayService) getEnvironmentAPIKey(
	ctx context.Context,
	apiKey string,
) (*accountproto.EnvironmentAPIKey, error) {
	k, err := s.singleflightFetch(
		ctx,
		environmentAPIKeyFlightID(apiKey),
		func(ctx context.Context) (interface{}, error) {
			// L1: in-memory cache
			envAPIKey, err := getEnvironmentAPIKeyFromCache(
				apiKey,
				s.environmentAPIKeyCache,
				callerGatewayService,
				cacheLayerInMemory,
			)
			if err == nil {
				return envAPIKey, nil
			}
			// L2: Redis cache (kept warm by batch cacher)
			envAPIKey, err = getEnvironmentAPIKeyFromCache(
				apiKey,
				s.environmentAPIKeyRedisCache,
				callerGatewayService,
				cacheLayerExternal,
			)
			if err == nil {
				putEnvironmentAPIKeyCache(
					ctx,
					envAPIKey,
					s.environmentAPIKeyCache,
					s.logger,
				)
				return envAPIKey, nil
			}
			// L3: direct DB query
			s.logger.Warn(
				"API key not found in cache",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("apiKey", obfuscateString(apiKey, obfuscateAPIKeyLength)),
				)...,
			)
			// Since the Get and List APIs for the API keys are obsfucated,
			// we need to directly query the database.
			domainEnvAPIKey, err := s.accountStorage.GetEnvironmentAPIKey(ctx, apiKey)
			if err != nil {
				if errors.Is(err, accstorage.ErrAPIKeyNotFound) {
					return nil, ErrInvalidAPIKey
				}
				s.logger.Error(
					"Failed to get environment APIKey from storage",
					log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
				)
				return nil, ErrInternal
			}
			envAPIKey = domainEnvAPIKey.EnvironmentAPIKey
			putEnvironmentAPIKeyCache(
				ctx,
				envAPIKey,
				s.environmentAPIKeyCache,
				s.logger,
			)
			return envAPIKey, nil
		},
	)
	if err != nil {
		return nil, translateCallerCanceledErr(ctx, err)
	}
	envAPIKey := k.(*accountproto.EnvironmentAPIKey)
	return envAPIKey, nil
}

func (s *grpcGatewayService) extractAPIKey(ctx context.Context) (string, error) {
	md, ok := gmetadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMissingAPIKey
	}
	keys, ok := md["authorization"]
	if !ok || len(keys) == 0 || keys[0] == "" {
		return "", ErrMissingAPIKey
	}
	return keys[0], nil
}

func environmentAPIKeyFlightID(id string) string {
	return id
}

func getEnvironmentAPIKeyFromCache(
	apikey string,
	c cachev3.EnvironmentAPIKeyCache,
	caller, layer string,
) (*accountproto.EnvironmentAPIKey, error) {
	envAPIKey, err := c.Get(apikey)
	if err == nil {
		cacheCounter.WithLabelValues(caller, typeAPIKey, layer, codeHit).Inc()
		return envAPIKey, nil
	}
	cacheCounter.WithLabelValues(caller, typeAPIKey, layer, codeMiss).Inc()
	return nil, err
}

func putEnvironmentAPIKeyCache(
	ctx context.Context,
	envAPIKey *accountproto.EnvironmentAPIKey,
	environmentAPIKeyCache cachev3.EnvironmentAPIKeyCache,
	logger *zap.Logger,
) {
	if err := environmentAPIKeyCache.Put(envAPIKey); err != nil {
		logger.Error(
			"Failed to cache environment APIKey",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", envAPIKey.Environment.Id),
			)...,
		)
	}
}

func checkEnvironmentAPIKey(
	environmentAPIKey *accountproto.EnvironmentAPIKey,
	roles []accountproto.APIKey_Role,
) error {

	if !contains(roles, environmentAPIKey.ApiKey.Role) {
		return ErrBadRole
	}
	if environmentAPIKey.EnvironmentDisabled {
		return ErrDisabledAPIKey
	}
	if environmentAPIKey.ApiKey.Disabled {
		return ErrDisabledAPIKey
	}
	return nil
}

func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// obfuscateString obfuscates the input string by showing the first n characters,
// replacing the middle with dots, and showing the last n characters.
func obfuscateString(input string, showLength int) string {
	// If the string length is exactly 2 * showLength, obfuscate with dots in the middle
	if len(input) > showLength*2 {
		return input[:showLength] + "...." + input[len(input)-showLength:]
	}
	return input
}

func isContextCanceled(ctx context.Context) bool {
	return ctx.Err() == context.Canceled
}

func (s *grpcGatewayService) filterOutArchivedFeatures(fs []*featureproto.Feature) []*featureproto.Feature {
	result := make([]*featureproto.Feature, 0)
	for _, f := range fs {
		if f.Archived {
			continue
		}
		result = append(result, f)
	}
	return result
}

func (s *grpcGatewayService) filterByTag(fs []*featureproto.Feature, tag string) []*featureproto.Feature {
	result := make([]*featureproto.Feature, 0)
	for _, f := range fs {
		for _, t := range f.Tags {
			if t == tag {
				result = append(result, f)
				break
			}
		}
	}
	return result
}

func (s *grpcGatewayService) emptyUserEvaluations() *featureproto.UserEvaluations {
	return &featureproto.UserEvaluations{
		Id:                 "no_evaluations",
		Evaluations:        []*featureproto.Evaluation{},
		CreatedAt:          time.Now().Unix(),
		ArchivedFeatureIds: []string{},
		ForceUpdate:        false,
	}
}
