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
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
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
	ftstorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	bigtable "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigtable"
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
	ErrUserRequired      = status.Error(codes.InvalidArgument, "gateway: user is required")
	ErrUserIDRequired    = status.Error(codes.InvalidArgument, "gateway: user id is required")
	ErrGoalIDRequired    = status.Error(codes.InvalidArgument, "gateway: goal id is required")
	ErrFeatureIDRequired = status.Error(codes.InvalidArgument, "gateway: feature id is required")
	ErrTagRequired       = status.Error(codes.InvalidArgument, "gateway: tag is required")
	ErrMissingEvents     = status.Error(codes.InvalidArgument, "gateway: missing events")
	ErrMissingEventID    = status.Error(codes.InvalidArgument, "gateway: missing event id")
	ErrInvalidTimestamp  = status.Error(codes.InvalidArgument, "gateway: invalid timestamp")
	ErrContextCanceled   = status.Error(codes.Canceled, "gateway: context canceled")
	ErrFeatureNotFound   = status.Error(codes.NotFound, "gateway: feature not found")
	ErrMissingAPIKey     = status.Error(codes.Unauthenticated, "gateway: missing APIKey")
	ErrInvalidAPIKey     = status.Error(codes.PermissionDenied, "gateway: invalid APIKey")
	ErrDisabledAPIKey    = status.Error(codes.PermissionDenied, "gateway: disabled APIKey")
	ErrBadRole           = status.Error(codes.PermissionDenied, "gateway: bad role")
	ErrInternal          = status.Error(codes.Internal, "gateway: internal")

	grpcGoalEvent       = &eventproto.GoalEvent{}
	grpcGoalBatchEvent  = &eventproto.GoalBatchEvent{}
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
	userEvaluationStorage  ftstorage.UserEvaluationsStorage
	featureClient          featureclient.Client
	accountClient          accountclient.Client
	goalPublisher          publisher.Publisher
	goalBatchPublisher     publisher.Publisher
	evaluationPublisher    publisher.Publisher
	userPublisher          publisher.Publisher
	metricsPublisher       publisher.Publisher
	featuresCache          cachev3.FeaturesCache
	segmentUsersCache      cachev3.SegmentUsersCache
	environmentAPIKeyCache cachev3.EnvironmentAPIKeyCache
	flightgroup            singleflight.Group
	opts                   *options
	logger                 *zap.Logger
}

func NewGrpcGatewayService(
	bt bigtable.Client,
	featureClient featureclient.Client,
	accountClient accountclient.Client,
	gp publisher.Publisher,
	gbp publisher.Publisher,
	ep publisher.Publisher,
	up publisher.Publisher,
	mp publisher.Publisher,
	v3Cache cache.MultiGetCache,
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
		userEvaluationStorage:  ftstorage.NewUserEvaluationsStorage(bt),
		featureClient:          featureClient,
		accountClient:          accountClient,
		goalPublisher:          gp,
		goalBatchPublisher:     gbp,
		evaluationPublisher:    ep,
		userPublisher:          up,
		metricsPublisher:       mp,
		featuresCache:          cachev3.NewFeaturesCache(v3Cache),
		segmentUsersCache:      cachev3.NewSegmentUsersCache(v3Cache),
		environmentAPIKeyCache: cachev3.NewEnvironmentAPIKeyCache(v3Cache),
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
				zap.String("environmentNamespace", envAPIKey.EnvironmentNamespace),
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
				zap.String("environmentNamespace", envAPIKey.EnvironmentNamespace),
				zap.String("goalId", goalEvent.GoalId),
			)...,
		)
		return nil, ErrInternal
	}
	event := &eventproto.Event{
		Id:                   id.String(),
		Event:                goal,
		EnvironmentNamespace: envAPIKey.EnvironmentNamespace,
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
				zap.String("environmentNamespace", envAPIKey.EnvironmentNamespace),
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
	envAPIKey, err := s.checkRequest(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.validateGetEvaluationsRequest(req); err != nil {
		return nil, err
	}
	s.publishUser(ctx, envAPIKey.EnvironmentNamespace, req.Tag, req.User, req.SourceId)
	f, err, _ := s.flightgroup.Do(
		envAPIKey.EnvironmentNamespace,
		func() (interface{}, error) {
			return s.getFeatures(ctx, envAPIKey.EnvironmentNamespace)
		},
	)
	if err != nil {
		return nil, err
	}
	features := f.([]*featureproto.Feature)
	if len(features) == 0 {
		return &gwproto.GetEvaluationsResponse{
			State:       featureproto.UserEvaluations_FULL,
			Evaluations: nil,
		}, nil
	}
	ueid := featuredomain.UserEvaluationsID(req.User.Id, req.User.Data, features)
	if req.UserEvaluationsId == ueid {
		return &gwproto.GetEvaluationsResponse{
			State:             featureproto.UserEvaluations_FULL,
			Evaluations:       nil,
			UserEvaluationsId: ueid,
		}, nil
	}
	evaluations, err := s.evaluateFeatures(ctx, req.User, features, envAPIKey.EnvironmentNamespace, req.Tag)
	if err != nil {
		s.logger.Error(
			"Failed to evaluate features",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", envAPIKey.EnvironmentNamespace),
				zap.String("userId", req.User.Id),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.GetEvaluationsResponse{
		State:             featureproto.UserEvaluations_FULL,
		Evaluations:       evaluations,
		UserEvaluationsId: ueid,
	}, nil
}

func (s *grpcGatewayService) validateGetEvaluationsRequest(req *gwproto.GetEvaluationsRequest) error {
	if req.Tag == "" {
		return ErrTagRequired
	}
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
	envAPIKey, err := s.checkRequest(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.validateGetEvaluationRequest(req); err != nil {
		return nil, err
	}
	s.publishUser(ctx, envAPIKey.EnvironmentNamespace, req.Tag, req.User, req.SourceId)
	f, err, _ := s.flightgroup.Do(
		envAPIKey.EnvironmentNamespace,
		func() (interface{}, error) {
			return s.getFeatures(ctx, envAPIKey.EnvironmentNamespace)
		},
	)
	if err != nil {
		return nil, err
	}
	fs := f.([]*featureproto.Feature)
	var features []*featureproto.Feature
	for _, f := range fs {
		if f.Id == req.FeatureId {
			features = append(features, f)
			break
		}
	}
	if len(features) == 0 {
		return nil, ErrFeatureNotFound
	}
	evaluations, err := s.evaluateFeatures(ctx, req.User, features, envAPIKey.EnvironmentNamespace, req.Tag)
	if err != nil {
		s.logger.Error(
			"Failed to evaluate features",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", envAPIKey.EnvironmentNamespace),
				zap.String("userId", req.User.Id),
				zap.String("featureId", req.FeatureId),
			)...,
		)
		return nil, ErrInternal
	}
	if err := s.upsertUserEvaluation(
		ctx,
		envAPIKey.EnvironmentNamespace,
		req.Tag,
		evaluations.Evaluations[0],
	); err != nil {
		eventCounter.WithLabelValues(callerGatewayService, typeMetrics, codeUpsertUserEvaluationFailed).Inc()
		s.logger.Error(
			"Failed to upsert user evaluation while trying to get evaluation",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", envAPIKey.EnvironmentNamespace),
				zap.String("userId", req.User.Id),
				zap.String("featureId", req.FeatureId),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.GetEvaluationResponse{
		Evaluation: evaluations.Evaluations[0],
	}, nil
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
	environmentNamespace,
	tag string,
	user *userproto.User,
	sourceID eventproto.SourceId,
) {
	// TODO: using buffered channel to reduce the number of go routines
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.opts.pubsubTimeout)
		defer cancel()
		if err := s.publishUserEvent(ctx, user, tag, environmentNamespace, sourceID); err != nil {
			s.logger.Error(
				"Failed to publish UserEvent",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
		}
	}()
}

func (s *grpcGatewayService) publishUserEvent(
	ctx context.Context,
	user *userproto.User,
	tag, environmentNamespace string,
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
		Data:                 user.Data,
		EnvironmentNamespace: environmentNamespace,
	}
	ue, err := ptypes.MarshalAny(userEvent)
	if err != nil {
		return err
	}
	event := &eventproto.Event{
		Id:                   id.String(),
		Event:                ue,
		EnvironmentNamespace: environmentNamespace,
	}
	return s.userPublisher.Publish(ctx, event)
}

func (s *grpcGatewayService) getFeatures(
	ctx context.Context,
	environmentNamespace string,
) ([]*featureproto.Feature, error) {
	fs, err := s.getFeaturesFromCache(ctx, environmentNamespace)
	if err == nil {
		return fs.Features, nil
	}
	s.logger.Info(
		"No cached data for Features",
		log.FieldsFromImcomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
		)...,
	)
	features, err := s.listFeatures(ctx, environmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to retrieve features from storage",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
			)...,
		)
		return nil, ErrInternal
	}
	if err := s.featuresCache.Put(&featureproto.Features{Features: features}, environmentNamespace); err != nil {
		s.logger.Error(
			"Failed to cache features",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
			)...,
		)
	}
	return features, nil
}

func (s *grpcGatewayService) listFeatures(
	ctx context.Context,
	environmentNamespace string,
) ([]*featureproto.Feature, error) {
	features := []*featureproto.Feature{}
	cursor := ""
	for {
		resp, err := s.featureClient.ListFeatures(ctx, &featureproto.ListFeaturesRequest{
			PageSize:             listRequestSize,
			Cursor:               cursor,
			EnvironmentNamespace: environmentNamespace,
			Archived:             &wrappers.BoolValue{Value: false},
		})
		if err != nil {
			return nil, err
		}
		for _, f := range resp.Features {
			if !f.Enabled && f.OffVariation == "" {
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
	environmentNamespace string,
) (*featureproto.Features, error) {
	features, err := s.featuresCache.Get(environmentNamespace)
	if err == nil {
		cacheCounter.WithLabelValues(callerGatewayService, typeFeatures, cacheLayerExternal, codeHit).Inc()
		return features, nil
	}
	cacheCounter.WithLabelValues(callerGatewayService, typeFeatures, cacheLayerExternal, codeMiss).Inc()
	return nil, err
}

func (s *grpcGatewayService) evaluateFeatures(
	ctx context.Context,
	user *userproto.User,
	features []*featureproto.Feature,
	environmentNamespace, tag string,
) (*featureproto.UserEvaluations, error) {
	mapIDs := make(map[string]struct{})
	for _, f := range features {
		feature := &featuredomain.Feature{Feature: f}
		for _, id := range feature.ListSegmentIDs() {
			mapIDs[id] = struct{}{}
		}
	}
	mapSegmentUsers, err := s.listSegmentUsers(ctx, user.Id, mapIDs, environmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to list segments",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
			)...,
		)
		return nil, err
	}
	userEvaluations, err := featuredomain.EvaluateFeatures(features, user, mapSegmentUsers, tag)
	if err != nil {
		s.logger.Error(
			"Failed to evaluate",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
			)...,
		)
	}
	return userEvaluations, nil
}

func (s *grpcGatewayService) listSegmentUsers(
	ctx context.Context,
	userID string,
	mapSegmentIDs map[string]struct{},
	environmentNamespace string,
) (map[string][]*featureproto.SegmentUser, error) {
	if len(mapSegmentIDs) == 0 {
		return nil, nil
	}
	users := make(map[string][]*featureproto.SegmentUser)
	for segmentID := range mapSegmentIDs {
		s, err, _ := s.flightgroup.Do(s.segmentFlightID(environmentNamespace, segmentID), func() (interface{}, error) {
			return s.getSegmentUsers(ctx, segmentID, environmentNamespace)
		})
		if err != nil {
			return nil, err
		}
		segmentUsers := s.([]*featureproto.SegmentUser)
		users[segmentID] = segmentUsers
	}
	return users, nil
}

func (s *grpcGatewayService) segmentFlightID(environmentNamespace, segmentID string) string {
	return environmentNamespace + ":" + segmentID
}

func (s *grpcGatewayService) getSegmentUsers(
	ctx context.Context,
	segmentID, environmentNamespace string,
) ([]*featureproto.SegmentUser, error) {
	segmentUsers, err := s.getSegmentUsersFromCache(segmentID, environmentNamespace)
	if err == nil {
		return segmentUsers, nil
	}
	s.logger.Info(
		"No cached data for SegmentUsers",
		log.FieldsFromImcomingContext(ctx).AddFields(
			zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("segmentId", segmentID),
		)...,
	)
	req := &featureproto.ListSegmentUsersRequest{
		SegmentId:            segmentID,
		EnvironmentNamespace: environmentNamespace,
	}
	res, err := s.featureClient.ListSegmentUsers(ctx, req)
	if err != nil {
		s.logger.Error(
			"Failed to retrieve segment users from storage",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
				zap.String("segmentId", segmentID),
			)...,
		)
		return nil, ErrInternal
	}
	su := &featureproto.SegmentUsers{
		SegmentId: segmentID,
		Users:     res.Users,
	}
	if err := s.segmentUsersCache.Put(su, environmentNamespace); err != nil {
		s.logger.Error(
			"Failed to cache segment users",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
				zap.String("segmentId", segmentID),
			)...,
		)
	}
	return res.Users, nil
}

func (s *grpcGatewayService) getSegmentUsersFromCache(
	segmentID, environmentNamespace string,
) ([]*featureproto.SegmentUser, error) {
	segment, err := s.segmentUsersCache.Get(segmentID, environmentNamespace)
	if err == nil {
		cacheCounter.WithLabelValues(callerGatewayService, typeSegmentUsers, cacheLayerExternal, codeHit).Inc()
		return segment.Users, nil
	}
	cacheCounter.WithLabelValues(callerGatewayService, typeSegmentUsers, cacheLayerExternal, codeMiss).Inc()
	return nil, err
}

func (s *grpcGatewayService) upsertUserEvaluation(
	ctx context.Context,
	environmentNamespace, tag string,
	evaluation *featureproto.Evaluation,
) error {
	if err := s.userEvaluationStorage.UpsertUserEvaluation(
		ctx,
		evaluation,
		environmentNamespace,
		tag,
	); err != nil {
		return err
	}
	return nil
}

func (s *grpcGatewayService) convToEvaluation(
	ctx context.Context,
	event *eventproto.Event,
) (*featureproto.Evaluation, string, error) {
	ev := &eventproto.EvaluationEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		s.logger.Error(
			"Failed to extract evaluation event for converting evaluation",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", event.Id),
			)...,
		)
		return nil, "", errUnmarshalFailed
	}
	evaluation := &featureproto.Evaluation{
		Id: featuredomain.EvaluationID(
			ev.FeatureId,
			ev.FeatureVersion,
			ev.UserId,
		),
		FeatureId:      ev.FeatureId,
		FeatureVersion: ev.FeatureVersion,
		UserId:         ev.UserId,
		VariationId:    ev.VariationId,
		Reason:         ev.Reason,
	}
	// For requests that doesn't have the tag info,
	// it will insert none instead, until all SDK clients are updated
	var tag string
	if ev.Tag == "" {
		tag = "none"
	} else {
		tag = ev.Tag
	}
	return evaluation, tag, nil
}

func (s *grpcGatewayService) RegisterEvents(
	ctx context.Context,
	req *gwproto.RegisterEventsRequest,
) (*gwproto.RegisterEventsResponse, error) {
	envAPIKey, err := s.checkRequest(ctx)
	if err != nil {
		return nil, err
	}
	if len(req.Events) == 0 {
		return nil, ErrMissingEvents
	}
	errs := make(map[string]*gwproto.RegisterEventsResponse_Error)
	goalMessages := make([]publisher.Message, 0)
	goalBatchMessages := make([]publisher.Message, 0)
	evaluationMessages := make([]publisher.Message, 0)
	metricsMessages := make([]publisher.Message, 0)
	publish := func(p publisher.Publisher, messages []publisher.Message, typ string) {
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
					zap.String("environmentNamespace", envAPIKey.EnvironmentNamespace),
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
	}
	for _, event := range req.Events {
		event.EnvironmentNamespace = envAPIKey.EnvironmentNamespace
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
		if ptypes.Is(event.Event, grpcGoalBatchEvent) {
			errorCode, err := validator.validate(ctx)
			if err != nil {
				eventCounter.WithLabelValues(callerGatewayService, typeGoalBatch, errorCode).Inc()
				errs[event.Id] = &gwproto.RegisterEventsResponse_Error{
					Retriable: false,
					Message:   err.Error(),
				}
				continue
			}
			goalBatchMessages = append(goalBatchMessages, event)
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
			evaluation, tag, err := s.convToEvaluation(ctx, event)
			if err != nil {
				eventCounter.WithLabelValues(callerGatewayService, typeEvaluation, codeEvaluationConversionFailed).Inc()
				errs[event.Id] = &gwproto.RegisterEventsResponse_Error{
					Retriable: false,
					Message:   err.Error(),
				}
				continue
			}
			if err := s.upsertUserEvaluation(ctx, envAPIKey.EnvironmentNamespace, tag, evaluation); err != nil {
				eventCounter.WithLabelValues(callerGatewayService, typeEvaluation, codeUpsertUserEvaluationFailed).Inc()
				errs[event.Id] = &gwproto.RegisterEventsResponse_Error{
					Retriable: true,
					Message:   "Failed to upsert user evaluation",
				}
				continue
			}
			evaluationMessages = append(evaluationMessages, event)
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
			metricsMessages = append(metricsMessages, event)
		}
	}
	publish(s.goalPublisher, goalMessages, typeGoal)
	publish(s.goalBatchPublisher, goalBatchMessages, typeGoalBatch)
	publish(s.evaluationPublisher, evaluationMessages, typeEvaluation)
	publish(s.metricsPublisher, metricsMessages, typeMetrics)
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
	k, err, _ := s.flightgroup.Do(
		environmentAPIKeyFlightID(apiKey),
		func() (interface{}, error) {
			return getEnvironmentAPIKey(
				ctx,
				apiKey,
				s.accountClient,
				s.environmentAPIKeyCache,
				callerGatewayService,
				s.logger,
			)
		},
	)
	if err != nil {
		return nil, err
	}
	envAPIKey := k.(*accountproto.EnvironmentAPIKey)
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
	caller string,
	logger *zap.Logger,
) (*accountproto.EnvironmentAPIKey, error) {
	envAPIKey, err := getEnvironmentAPIKeyFromCache(ctx, id, environmentAPIKeyCache, caller, cacheLayerExternal)
	if err == nil {
		return envAPIKey, nil
	}
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
	envAPIKey = resp.EnvironmentApiKey
	if err := environmentAPIKeyCache.Put(envAPIKey); err != nil {
		logger.Error(
			"Failed to cache environment APIKey",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", envAPIKey.EnvironmentNamespace),
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
