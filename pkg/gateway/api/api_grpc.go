// Copyright 2024 The Bucketeer Authors.
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
	"time"

	"github.com/golang/protobuf/ptypes"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	gmetadata "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	featuredomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	serviceeventproto "github.com/bucketeer-io/bucketeer/proto/event/service"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

const (
	listRequestSize = 500
)

var (
	ErrUserRequired       = status.Error(codes.InvalidArgument, "gateway: user is required")
	ErrUserIDRequired     = status.Error(codes.InvalidArgument, "gateway: user id is required")
	ErrGoalIDRequired     = status.Error(codes.InvalidArgument, "gateway: goal id is required")
	ErrFeatureIDRequired  = status.Error(codes.InvalidArgument, "gateway: feature id is required")
	ErrTagRequired        = status.Error(codes.InvalidArgument, "gateway: tag is required")
	ErrMissingEvents      = status.Error(codes.InvalidArgument, "gateway: missing events")
	ErrMissingEventID     = status.Error(codes.InvalidArgument, "gateway: missing event id")
	ErrInvalidTimestamp   = status.Error(codes.InvalidArgument, "gateway: invalid timestamp")
	ErrContextCanceled    = status.Error(codes.Canceled, "gateway: context canceled")
	ErrFeatureNotFound    = status.Error(codes.NotFound, "gateway: feature not found")
	ErrEvaluationNotFound = status.Error(codes.NotFound, "gateway: evaluation not found")
	ErrMissingAPIKey      = status.Error(codes.Unauthenticated, "gateway: missing APIKey")
	ErrInvalidAPIKey      = status.Error(codes.PermissionDenied, "gateway: invalid APIKey")
	ErrDisabledAPIKey     = status.Error(codes.PermissionDenied, "gateway: disabled APIKey")
	ErrBadRole            = status.Error(codes.PermissionDenied, "gateway: bad role")
	ErrInternal           = status.Error(codes.Internal, "gateway: internal")

	grpcGoalEvent       = &eventproto.GoalEvent{}
	grpcEvaluationEvent = &eventproto.EvaluationEvent{}
	grpcMetricsEvent    = &eventproto.MetricsEvent{}
)

type options struct {
	apiKeyMemoryCacheTTL              time.Duration
	apiKeyMemoryCacheEvictionInterval time.Duration
	pubsubTimeout                     time.Duration
	oldestEventTimestamp              time.Duration
	furthestEventTimestamp            time.Duration
	metrics                           metrics.Registerer
	logger                            *zap.Logger
}

var defaultOptions = options{
	apiKeyMemoryCacheTTL:              5 * time.Minute,
	apiKeyMemoryCacheEvictionInterval: 30 * time.Second,
	pubsubTimeout:                     20 * time.Second,
	oldestEventTimestamp:              24 * time.Hour,
	furthestEventTimestamp:            24 * time.Hour,
	logger:                            zap.NewNop(),
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
	featureClient          featureclient.Client
	accountClient          accountclient.Client
	goalPublisher          publisher.Publisher
	evaluationPublisher    publisher.Publisher
	userPublisher          publisher.Publisher
	featuresCache          cachev3.FeaturesCache
	segmentUsersCache      cachev3.SegmentUsersCache
	environmentAPIKeyCache cachev3.EnvironmentAPIKeyCache
	flightgroup            singleflight.Group
	opts                   *options
	logger                 *zap.Logger
}

func NewGrpcGatewayService(
	featureClient featureclient.Client,
	accountClient accountclient.Client,
	gp publisher.Publisher,
	ep publisher.Publisher,
	up publisher.Publisher,
	redisV3Cache cache.MultiGetCache,
	inMemoryCache cache.Cache,
	opts ...Option,
) rpc.Service {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.metrics != nil {
		registerMetrics(options.metrics)
	}
	return &grpcGatewayService{
		featureClient:          featureClient,
		accountClient:          accountClient,
		goalPublisher:          gp,
		evaluationPublisher:    ep,
		userPublisher:          up,
		featuresCache:          cachev3.NewFeaturesCache(redisV3Cache),
		segmentUsersCache:      cachev3.NewSegmentUsersCache(redisV3Cache),
		environmentAPIKeyCache: cachev3.NewEnvironmentAPIKeyCache(inMemoryCache),
		opts:                   &options,
		logger:                 options.logger.Named("api_grpc"),
	}
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
		s.logger.Warn(
			"Invalid track url parameters",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	envAPIKey, err := s.checkTrackRequest(ctx, req.Apikey)
	if err != nil {
		s.logger.Error(
			"Failed to get environment api key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("apiKey", req.Apikey),
			)...,
		)
		return nil, err
	}
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
			log.FieldsFromImcomingContext(ctx).AddFields(
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
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", envAPIKey.Environment.Id),
				zap.String("goalId", goalEvent.GoalId),
			)...,
		)
		return nil, ErrInternal
	}
	event := &eventproto.Event{
		Id:                   id.String(),
		Event:                goal,
		EnvironmentNamespace: envAPIKey.Environment.Id,
	}
	if err := s.goalPublisher.Publish(ctx, event); err != nil {
		if err == publisher.ErrBadMessage {
			eventCounter.WithLabelValues(callerGatewayService, typeGoal, codeNonRepeatableError)
		} else {
			eventCounter.WithLabelValues(callerGatewayService, typeGoal, codeRepeatableError)
		}
		s.logger.Error(
			"Failed to publish goal event",
			log.FieldsFromImcomingContext(ctx).AddFields(
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
	envAPIKey, err := s.checkRequest(ctx)
	if err != nil {
		return nil, err
	}
	projectID := envAPIKey.ProjectId
	environmentId := envAPIKey.Environment.Id
	if err := s.validateGetEvaluationsRequest(req); err != nil {
		evaluationsCounter.WithLabelValues(projectID, environmentId, req.Tag, evaluationBadRequest).Inc()
		return nil, err
	}
	s.publishUser(ctx, environmentId, req.Tag, req.User, req.SourceId)
	ctx, spanGetFeatures := trace.StartSpan(ctx, "bucketeerGRPCGatewayService.GetEvaluations.GetFeatures")
	f, err, _ := s.flightgroup.Do(
		environmentId,
		func() (interface{}, error) {
			return s.getFeatures(ctx, environmentId)
		},
	)
	if err != nil {
		evaluationsCounter.WithLabelValues(projectID, environmentId, req.Tag, evaluationInternalError).Inc()
		return nil, err
	}
	spanGetFeatures.End()
	features := f.([]*featureproto.Feature)
	activeFeatures := s.filterOutArchivedFeatures(features)
	filteredByTag := s.filterByTag(activeFeatures, req.Tag)
	if len(features) == 0 {
		evaluationsCounter.WithLabelValues(projectID, environmentId, req.Tag, evaluationNoFeatures).Inc()
		return &gwproto.GetEvaluationsResponse{
			State:             featureproto.UserEvaluations_FULL,
			Evaluations:       s.emptyUserEvaluations(),
			UserEvaluationsId: "no_evaluations",
		}, nil
	}
	ueid := featuredomain.UserEvaluationsID(req.User.Id, req.User.Data, filteredByTag)
	if req.UserEvaluationsId == ueid {
		evaluationsCounter.WithLabelValues(projectID, environmentId, req.Tag, evaluationNone).Inc()
		s.logger.Debug(
			"Features length when UEID is the same",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("environmentID", environmentId),
				zap.String("tag", req.Tag),
				zap.Int("featuresLength", len(features)),
				zap.Int("activeFeaturesLength", len(activeFeatures)),
				zap.Int("filteredByTagLength", len(filteredByTag)),
			)...,
		)
		return &gwproto.GetEvaluationsResponse{
			State:             featureproto.UserEvaluations_FULL,
			Evaluations:       s.emptyUserEvaluations(),
			UserEvaluationsId: ueid,
		}, nil
	}
	segmentUsersMap, err := s.getSegmentUsersMap(ctx, req.User, features, environmentId)
	if err != nil {
		evaluationsCounter.WithLabelValues(projectID, environmentId, req.Tag, evaluationInternalError).Inc()
		s.logger.Error(
			"Failed to get segment users map",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentId),
			)...,
		)
		return nil, err
	}
	var evaluations *featureproto.UserEvaluations
	// FIXME Remove s.getEvaluations once all SDKs use UserEvaluationCondition.
	// New SDKs always use UserEvaluationCondition.
	if req.UserEvaluationCondition == nil {
		// Old evaluation requires tag to be set.
		if req.Tag == "" {
			evaluationsCounter.WithLabelValues(projectID, environmentId, req.Tag, evaluationBadRequest).Inc()
			return nil, ErrTagRequired
		}
		evaluations, err = featuredomain.EvaluateFeatures(
			activeFeatures,
			req.User,
			segmentUsersMap,
			req.Tag,
		)
		if err != nil {
			evaluationsCounter.WithLabelValues(projectID, environmentId, req.Tag, evaluationInternalError).Inc()
			s.logger.Error(
				"Failed to evaluate",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("userId", req.User.Id),
					zap.String("environmentID", environmentId),
				)...,
			)
			return nil, ErrInternal
		}
		evaluationsCounter.WithLabelValues(projectID, environmentId, req.Tag, evaluationOld).Inc()
	} else {
		evaluations, err = featuredomain.EvaluateFeaturesByEvaluatedAt(
			features,
			req.User,
			segmentUsersMap,
			req.UserEvaluationsId,
			req.UserEvaluationCondition.EvaluatedAt,
			req.UserEvaluationCondition.UserAttributesUpdated,
			req.Tag,
		)
		if err != nil {
			evaluationsCounter.WithLabelValues(projectID, environmentId, req.Tag, evaluationInternalError).Inc()
			s.logger.Error(
				"Failed to evaluate",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("userId", req.User.Id),
					zap.String("environmentID", environmentId),
				)...,
			)
			return nil, ErrInternal
		}
		if evaluations.ForceUpdate {
			evaluationsCounter.WithLabelValues(projectID, environmentId, req.Tag, evaluationAll).Inc()
		} else {
			evaluationsCounter.WithLabelValues(projectID, environmentId, req.Tag, evaluationDiff).Inc()
		}
	}
	s.logger.Debug(
		"Features length when UEID is different",
		log.FieldsFromImcomingContext(ctx).AddFields(
			zap.String("environmentID", environmentId),
			zap.String("tag", req.Tag),
			zap.Int("featuresLength", len(features)),
			zap.Int("activeFeaturesLength", len(activeFeatures)),
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
	envAPIKey, err := s.checkRequest(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.validateGetEvaluationRequest(req); err != nil {
		return nil, err
	}
	s.publishUser(ctx, envAPIKey.Environment.Id, req.Tag, req.User, req.SourceId)
	ctx, spanGetFeatures := trace.StartSpan(ctx, "bucketeerGRPCGatewayService.GetEvaluation.GetFeatures")
	f, err, _ := s.flightgroup.Do(
		envAPIKey.Environment.Id,
		func() (interface{}, error) {
			return s.getFeatures(ctx, envAPIKey.Environment.Id)
		},
	)
	if err != nil {
		return nil, err
	}
	spanGetFeatures.End()
	fs := s.filterOutArchivedFeatures(f.([]*featureproto.Feature))
	features, err := s.getTargetFeatures(fs, req.FeatureId)
	if err != nil {
		return nil, err
	}
	segmentUsersMap, err := s.getSegmentUsersMap(ctx, req.User, features, envAPIKey.Environment.Id)
	if err != nil {
		s.logger.Error(
			"Failed to get segment users map",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", envAPIKey.Environment.Id),
			)...,
		)
		return nil, err
	}
	evaluations, err := featuredomain.EvaluateFeatures(features, req.User, segmentUsersMap, req.Tag)
	if err != nil {
		s.logger.Error(
			"Failed to evaluate",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", envAPIKey.Environment.Id),
			)...,
		)
		return nil, err
	}
	if err != nil {
		s.logger.Error(
			"Failed to evaluate features",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", envAPIKey.Environment.Id),
				zap.String("userId", req.User.Id),
				zap.String("featureId", req.FeatureId),
			)...,
		)
		return nil, ErrInternal
	}
	eval, err := s.findEvaluation(evaluations.Evaluations, req.FeatureId)
	if err != nil {
		return nil, err
	}
	return &gwproto.GetEvaluationResponse{
		Evaluation: eval,
	}, nil
}

func (s *grpcGatewayService) getTargetFeatures(fs []*featureproto.Feature, id string) ([]*featureproto.Feature, error) {
	feature, err := s.findFeature(fs, id)
	if err != nil {
		return nil, err
	}
	if len(feature.Prerequisites) == 0 {
		return []*featureproto.Feature{feature}, nil
	}
	return featuredomain.GetPrerequisiteDownwards([]*featureproto.Feature{feature}, fs)
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

func (s *grpcGatewayService) publishUser(
	ctx context.Context,
	environmentId,
	tag string,
	user *userproto.User,
	sourceID eventproto.SourceId,
) {
	// TODO: using buffered channel to reduce the number of go routines
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.opts.pubsubTimeout)
		defer cancel()
		if err := s.publishUserEvent(ctx, user, tag, environmentId, sourceID); err != nil {
			s.logger.Error(
				"Failed to publish UserEvent",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentID", environmentId),
				)...,
			)
		}
	}()
}

func (s *grpcGatewayService) publishUserEvent(
	ctx context.Context,
	user *userproto.User,
	tag, environmentId string,
	sourceID eventproto.SourceId,
) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	userEvent := &serviceeventproto.UserEvent{
		Id:                   id.String(),
		SourceId:             sourceID,
		Tag:                  tag,
		UserId:               user.Id,
		LastSeen:             time.Now().Unix(),
		Data:                 nil, // We set nil until we decide again what to do with the user metadata.
		EnvironmentNamespace: environmentId,
	}
	ue, err := ptypes.MarshalAny(userEvent)
	if err != nil {
		return err
	}
	event := &eventproto.Event{
		Id:                   id.String(),
		Event:                ue,
		EnvironmentNamespace: environmentId,
	}
	return s.userPublisher.Publish(ctx, event)
}

func (s *grpcGatewayService) getFeatures(
	ctx context.Context,
	environmentId string,
) ([]*featureproto.Feature, error) {
	fs, err := s.getFeaturesFromCache(ctx, environmentId)
	if err == nil {
		return fs.Features, nil
	}
	s.logger.Info(
		"No cached data for Features",
		log.FieldsFromImcomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentID", environmentId),
		)...,
	)
	features, err := s.listFeatures(ctx, environmentId)
	if err != nil {
		s.logger.Error(
			"Failed to retrieve features from storage",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentId),
			)...,
		)
		return nil, ErrInternal
	}
	if err := s.featuresCache.Put(&featureproto.Features{Features: features}, environmentId); err != nil {
		s.logger.Error(
			"Failed to cache features",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentId),
			)...,
		)
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
			PageSize:             listRequestSize,
			Cursor:               cursor,
			EnvironmentNamespace: environmentId,
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

func (s *grpcGatewayService) getFeaturesFromCache(
	ctx context.Context,
	environmentId string,
) (*featureproto.Features, error) {
	features, err := s.featuresCache.Get(environmentId)
	if err == nil {
		cacheCounter.WithLabelValues(callerGatewayService, typeFeatures, cacheLayerExternal, codeHit).Inc()
		return features, nil
	}
	cacheCounter.WithLabelValues(callerGatewayService, typeFeatures, cacheLayerExternal, codeMiss).Inc()
	return nil, err
}

func (s *grpcGatewayService) getSegmentUsersMap(
	ctx context.Context,
	user *userproto.User,
	features []*featureproto.Feature,
	environmentId string,
) (map[string][]*featureproto.SegmentUser, error) {
	mapIDs := make(map[string]struct{})
	for _, f := range features {
		feature := &featuredomain.Feature{Feature: f}
		for _, id := range feature.ListSegmentIDs() {
			mapIDs[id] = struct{}{}
		}
	}
	segmentUsersMap, err := s.listSegmentUsers(ctx, user.Id, mapIDs, environmentId)
	if err != nil {
		s.logger.Error(
			"Failed to list segments",
			log.FieldsFromImcomingContext(ctx).AddFields(
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
	userID string,
	mapSegmentIDs map[string]struct{},
	environmentId string,
) (map[string][]*featureproto.SegmentUser, error) {
	if len(mapSegmentIDs) == 0 {
		return nil, nil
	}
	users := make(map[string][]*featureproto.SegmentUser)
	for segmentID := range mapSegmentIDs {
		s, err, _ := s.flightgroup.Do(s.segmentFlightID(environmentId, segmentID), func() (interface{}, error) {
			return s.getSegmentUsers(ctx, segmentID, environmentId)
		})
		if err != nil {
			return nil, err
		}
		segmentUsers := s.([]*featureproto.SegmentUser)
		users[segmentID] = segmentUsers
	}
	return users, nil
}

func (s *grpcGatewayService) segmentFlightID(environmentId, segmentID string) string {
	return environmentId + ":" + segmentID
}

func (s *grpcGatewayService) getSegmentUsers(
	ctx context.Context,
	segmentID, environmentId string,
) ([]*featureproto.SegmentUser, error) {
	segmentUsers, err := s.getSegmentUsersFromCache(segmentID, environmentId)
	if err == nil {
		return segmentUsers, nil
	}
	s.logger.Info(
		"No cached data for SegmentUsers",
		log.FieldsFromImcomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentID", environmentId),
			zap.String("segmentId", segmentID),
		)...,
	)
	req := &featureproto.ListSegmentUsersRequest{
		SegmentId:            segmentID,
		EnvironmentNamespace: environmentId,
	}
	res, err := s.featureClient.ListSegmentUsers(ctx, req)
	if err != nil {
		s.logger.Error(
			"Failed to retrieve segment users from storage",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentId),
				zap.String("segmentId", segmentID),
			)...,
		)
		return nil, ErrInternal
	}
	su := &featureproto.SegmentUsers{
		SegmentId: segmentID,
		Users:     res.Users,
	}
	if err := s.segmentUsersCache.Put(su, environmentId); err != nil {
		s.logger.Error(
			"Failed to cache segment users",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentId),
				zap.String("segmentId", segmentID),
			)...,
		)
	}
	return res.Users, nil
}

func (s *grpcGatewayService) getSegmentUsersFromCache(
	segmentID, environmentId string,
) ([]*featureproto.SegmentUser, error) {
	segment, err := s.segmentUsersCache.Get(segmentID, environmentId)
	if err == nil {
		cacheCounter.WithLabelValues(callerGatewayService, typeSegmentUsers, cacheLayerExternal, codeHit).Inc()
		return segment.Users, nil
	}
	cacheCounter.WithLabelValues(callerGatewayService, typeSegmentUsers, cacheLayerExternal, codeMiss).Inc()
	return nil, err
}

func (s *grpcGatewayService) RegisterEvents(
	ctx context.Context,
	req *gwproto.RegisterEventsRequest,
) (*gwproto.RegisterEventsResponse, error) {
	ctx, span := trace.StartSpan(ctx, "bucketeerGRPCGatewayService.RegisterEvents")
	defer span.End()
	envAPIKey, err := s.checkRequest(ctx)
	if err != nil {
		return nil, err
	}
	if len(req.Events) == 0 {
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
		errors := p.PublishMulti(ctx, messages)
		var repeatableErrors, nonRepeateableErrors float64
		for id, err := range errors {
			retriable := err != publisher.ErrBadMessage
			if retriable {
				repeatableErrors++
			} else {
				nonRepeateableErrors++
			}
			s.logger.Error(
				"Failed to publish event",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentID", envAPIKey.Environment.Id),
					zap.String("eventType", typ),
					zap.String("id", id),
				)...,
			)
			errs[id] = &gwproto.RegisterEventsResponse_Error{
				Retriable: retriable,
				Message:   "Failed to publish event",
			}
		}
		eventCounter.WithLabelValues(callerGatewayService, typ, codeNonRepeatableError).Add(nonRepeateableErrors)
		eventCounter.WithLabelValues(callerGatewayService, typ, codeRepeatableError).Add(repeatableErrors)
		eventCounter.WithLabelValues(callerGatewayService, typ, codeOK).Add(float64(len(messages) - len(errors)))
		return errs
	}
	for i, event := range req.Events {
		event.EnvironmentNamespace = envAPIKey.Environment.Id
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
				zap.String("apiKey", envAPIKey.ApiKey.Id),
				zap.String("projectID", envAPIKey.ProjectId),
				zap.String("eventID", event.Id),
				zap.String("environmentNamespace", event.EnvironmentNamespace),
				zap.Any("event", event.Event),
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
			log.FieldsFromImcomingContext(ctx)...,
		)
		return nil, ErrContextCanceled
	}
	envAPIKey, err := s.getEnvironmentAPIKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	if err := checkEnvironmentAPIKey(envAPIKey, accountproto.APIKey_SDK); err != nil {
		return nil, err
	}
	return envAPIKey, nil
}

func (s *grpcGatewayService) checkRequest(ctx context.Context) (*accountproto.EnvironmentAPIKey, error) {
	if isContextCanceled(ctx) {
		s.logger.Warn(
			"Request was canceled",
			log.FieldsFromImcomingContext(ctx)...,
		)
		return nil, ErrContextCanceled
	}
	id, err := s.extractAPIKeyID(ctx)
	if err != nil {
		return nil, err
	}
	envAPIKey, err := s.getEnvironmentAPIKey(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := checkEnvironmentAPIKey(envAPIKey, accountproto.APIKey_SDK); err != nil {
		return nil, err
	}
	return envAPIKey, nil
}

func (s *grpcGatewayService) getEnvironmentAPIKey(
	ctx context.Context,
	apiKey string,
) (*accountproto.EnvironmentAPIKey, error) {
	envAPIKey, err := getEnvironmentAPIKeyFromCache(
		ctx,
		apiKey,
		s.environmentAPIKeyCache,
		callerGatewayService,
		cacheLayerInMemory,
	)
	if err == nil {
		return envAPIKey, nil
	}
	k, err, _ := s.flightgroup.Do(
		environmentAPIKeyFlightID(apiKey),
		func() (interface{}, error) {
			return getEnvironmentAPIKey(
				ctx,
				apiKey,
				s.accountClient,
				s.environmentAPIKeyCache,
				s.logger,
			)
		},
	)
	if err != nil {
		return nil, err
	}
	envAPIKey = k.(*accountproto.EnvironmentAPIKey)
	return envAPIKey, nil
}

func (s *grpcGatewayService) extractAPIKeyID(ctx context.Context) (string, error) {
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

func getEnvironmentAPIKey(
	ctx context.Context,
	id string,
	accountClient accountclient.Client,
	environmentAPIKeyCache cachev3.EnvironmentAPIKeyCache,
	logger *zap.Logger,
) (*accountproto.EnvironmentAPIKey, error) {
	resp, err := accountClient.GetAPIKeyBySearchingAllEnvironments(
		ctx,
		&accountproto.GetAPIKeyBySearchingAllEnvironmentsRequest{Id: id},
	)
	if err != nil {
		if code := status.Code(err); code == codes.NotFound {
			return nil, ErrInvalidAPIKey
		}
		logger.Error(
			"Failed to get environment APIKey from account service",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, ErrInternal
	}
	envAPIKey := resp.EnvironmentApiKey
	if err := environmentAPIKeyCache.Put(envAPIKey); err != nil {
		logger.Error(
			"Failed to cache environment APIKey",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", envAPIKey.Environment.Id),
			)...,
		)
	}
	return envAPIKey, nil
}

func getEnvironmentAPIKeyFromCache(
	ctx context.Context,
	id string,
	c cachev3.EnvironmentAPIKeyCache,
	caller, layer string,
) (*accountproto.EnvironmentAPIKey, error) {
	envAPIKey, err := c.Get(id)
	if err == nil {
		cacheCounter.WithLabelValues(caller, typeAPIKey, layer, codeHit).Inc()
		return envAPIKey, nil
	}
	cacheCounter.WithLabelValues(caller, typeAPIKey, layer, codeMiss).Inc()
	return nil, err
}

func checkEnvironmentAPIKey(environmentAPIKey *accountproto.EnvironmentAPIKey, role accountproto.APIKey_Role) error {
	if environmentAPIKey.ApiKey.Role != role {
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
