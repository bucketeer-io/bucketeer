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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	featuredomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/rest"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	serviceeventproto "github.com/bucketeer-io/bucketeer/proto/event/service"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

type gatewayService struct {
	featureClient          featureclient.Client
	accountClient          accountclient.Client
	goalPublisher          publisher.Publisher
	evaluationPublisher    publisher.Publisher
	userPublisher          publisher.Publisher
	metricsPublisher       publisher.Publisher
	segmentUsersCache      cachev3.SegmentUsersCache
	featuresCache          cachev3.FeaturesCache
	environmentAPIKeyCache cachev3.EnvironmentAPIKeyCache
	flightgroup            singleflight.Group
	opts                   *options
	logger                 *zap.Logger
}

func NewGatewayService(
	featureClient featureclient.Client,
	accountClient accountclient.Client,
	gp publisher.Publisher,
	ep publisher.Publisher,
	up publisher.Publisher,
	mp publisher.Publisher,
	redisV3Cache cache.MultiGetCache,
	inMemoryCache cache.Cache,
	opts ...Option,
) *gatewayService {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.metrics != nil {
		registerMetrics(options.metrics)
	}
	return &gatewayService{
		featureClient:          featureClient,
		accountClient:          accountClient,
		goalPublisher:          gp,
		evaluationPublisher:    ep,
		userPublisher:          up,
		metricsPublisher:       mp,
		featuresCache:          cachev3.NewFeaturesCache(redisV3Cache),
		segmentUsersCache:      cachev3.NewSegmentUsersCache(redisV3Cache),
		environmentAPIKeyCache: cachev3.NewEnvironmentAPIKeyCache(inMemoryCache),
		opts:                   &options,
		logger:                 options.logger.Named("api"),
	}
}

const (
	Version          = "/v1"
	Service          = "/gateway"
	pingAPI          = "/ping"
	evaluationsAPI   = "/evaluations"
	evaluationAPI    = "/evaluation"
	eventAPI         = "/events"
	authorizationKey = "authorization"
)

var (
	errContextCanceled    = rest.NewErrStatus(http.StatusBadRequest, "gateway: context canceled")
	errMissingAPIKey      = rest.NewErrStatus(http.StatusUnauthorized, "gateway: missing APIKey")
	errInvalidAPIKey      = rest.NewErrStatus(http.StatusUnauthorized, "gateway: invalid APIKey")
	errInternal           = rest.NewErrStatus(http.StatusInternalServerError, "gateway: internal")
	errInvalidHttpMethod  = rest.NewErrStatus(http.StatusMethodNotAllowed, "gateway: invalid http method")
	errTagRequired        = rest.NewErrStatus(http.StatusBadRequest, "gateway: tag is required")
	errUserRequired       = rest.NewErrStatus(http.StatusBadRequest, "gateway: user is required")
	errUserIDRequired     = rest.NewErrStatus(http.StatusBadRequest, "gateway: user id is required")
	errBadRole            = rest.NewErrStatus(http.StatusUnauthorized, "gateway: bad role")
	errDisabledAPIKey     = rest.NewErrStatus(http.StatusUnauthorized, "gateway: disabled APIKey")
	errFeatureNotFound    = rest.NewErrStatus(http.StatusNotFound, "gateway: feature not found")
	errFeatureIDRequired  = rest.NewErrStatus(http.StatusBadRequest, "gateway: feature id is required")
	errMissingEventID     = rest.NewErrStatus(http.StatusBadRequest, "gateway: missing event id")
	errMissingEvents      = rest.NewErrStatus(http.StatusBadRequest, "gateway: missing events")
	errBodyRequired       = rest.NewErrStatus(http.StatusBadRequest, "gateway: body is required")
	errEvaluationNotFound = rest.NewErrStatus(http.StatusNotFound, "gateway: evaluation not found")
)

var (
	errInvalidType = errors.New("gateway: invalid message type")
)

func (s *gatewayService) Register(mux *http.ServeMux) {
	s.regist(mux, pingAPI, s.ping)
	s.regist(mux, evaluationsAPI, s.getEvaluations)
	s.regist(mux, evaluationAPI, s.getEvaluation)
	s.regist(mux, eventAPI, s.registerEvents)
}

func (*gatewayService) regist(mux *http.ServeMux, path string, handler func(http.ResponseWriter, *http.Request)) {
	mux.HandleFunc(fmt.Sprintf("%s%s%s", Version, Service, path), handler)
}

type pingResponse struct {
	Time int64 `json:"time,omitempty"`
}

type getEvaluationsRequest struct {
	Tag               string              `json:"tag,omitempty"`
	User              *userproto.User     `json:"user,omitempty"`
	UserEvaluationsID string              `json:"user_evaluations_id,omitempty"`
	SourceID          eventproto.SourceId `json:"source_id,omitempty"`
}

type getEvaluationsResponse struct {
	Evaluations       *featureproto.UserEvaluations `json:"evaluations,omitempty"`
	UserEvaluationsID string                        `json:"user_evaluations_id,omitempty"`
}

type getEvaluationRequest struct {
	Tag       string              `json:"tag,omitempty"`
	User      *userproto.User     `json:"user,omitempty"`
	FeatureID string              `json:"feature_id,omitempty"`
	SourceId  eventproto.SourceId `json:"source_id,omitempty"`
}

type registerEventsRequest struct {
	Events []event `json:"events,omitempty"`
}

type registerEventsResponse struct {
	Errors map[string]*registerEventsResponseError `json:"errors,omitempty"`
}

type registerEventsResponseError struct {
	Retriable bool   `json:"retriable"` // omitempty is not used intentionally
	Message   string `json:"message,omitempty"`
}

type getEvaluationResponse struct {
	Evaluation *featureproto.Evaluation `json:"evaluations,omitempty"`
}

type event struct {
	ID                   string          `json:"id,omitempty"`
	Event                json.RawMessage `json:"event,omitempty"`
	EnvironmentNamespace string          `json:"environment_namespace,omitempty"`
	Type                 EventType       `json:"type,omitempty"`
}

type metricsEvent struct {
	Timestamp  int64                  `json:"timestamp,omitempty"`
	Event      json.RawMessage        `json:"event,omitempty"`
	SdkVersion string                 `json:"sdk_version,omitempty"`
	Type       metricsDetailEventType `json:"type,omitempty"`
}

type latencyMetricsEvent struct {
	ApiId    eventproto.ApiId  `json:"api_id,omitempty"`
	Labels   map[string]string `json:"labels,omitempty"`
	Duration time.Duration     `json:"duration,omitempty"`
}

type sizeMetricsEvent struct {
	ApiId    eventproto.ApiId  `json:"api_id,omitempty"`
	Labels   map[string]string `json:"labels,omitempty"`
	SizeByte int32             `json:"size_byte,omitempty"`
}

type timeoutErrorMetricsEvent struct {
	ApiId  eventproto.ApiId  `json:"api_id,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

type internalErrorMetricsEvent struct {
	ApiId  eventproto.ApiId  `json:"api_id,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

type networkErrorMetricsEvent struct {
	ApiId  eventproto.ApiId  `json:"api_id,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

type internalSdkErrorMetricsEvent struct {
	ApiId  eventproto.ApiId  `json:"api_id,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

func (s *gatewayService) ping(w http.ResponseWriter, req *http.Request) {
	rest.ReturnSuccessResponse(
		w,
		&pingResponse{
			Time: time.Now().Unix(),
		},
	)
}

func (s *gatewayService) getEvaluations(w http.ResponseWriter, req *http.Request) {
	envAPIKey, reqBody, err := s.checkGetEvaluationsRequest(req)
	if err != nil {
		rest.ReturnFailureResponse(w, err)
		return
	}
	s.publishUser(req.Context(), envAPIKey.Environment.Id, reqBody.Tag, reqBody.User, reqBody.SourceID)
	f, err, _ := s.flightgroup.Do(
		envAPIKey.Environment.Id,
		func() (interface{}, error) {
			return s.getFeatures(req.Context(), envAPIKey.Environment.Id)
		},
	)
	if err != nil {
		rest.ReturnFailureResponse(w, err)
		return
	}
	features := s.filterOutArchivedFeatures(f.([]*featureproto.Feature))
	if len(features) == 0 {
		rest.ReturnSuccessResponse(
			w,
			&getEvaluationsResponse{
				Evaluations:       s.emptyUserEvaluations(),
				UserEvaluationsID: "no_evaluations",
			},
		)
		return
	}
	ueid := featuredomain.UserEvaluationsID(reqBody.User.Id, reqBody.User.Data, features)
	if reqBody.UserEvaluationsID == ueid {
		rest.ReturnSuccessResponse(
			w,
			&getEvaluationsResponse{
				Evaluations:       s.emptyUserEvaluations(),
				UserEvaluationsID: ueid,
			},
		)
		return
	}
	evaluations, err := s.evaluateFeatures(
		req.Context(),
		reqBody.User,
		features,
		envAPIKey.Environment.Id,
		reqBody.Tag,
	)
	if err != nil {
		s.logger.Error(
			"Failed to evaluate features",
			log.FieldsFromImcomingContext(req.Context()).AddFields(
				zap.Error(err),
				zap.String("environmentID", envAPIKey.Environment.Id),
				zap.String("userId", reqBody.User.Id),
			)...,
		)
		rest.ReturnFailureResponse(w, err)
		return
	}
	rest.ReturnSuccessResponse(
		w,
		&getEvaluationsResponse{
			Evaluations:       evaluations,
			UserEvaluationsID: ueid,
		},
	)
}

func (s *gatewayService) getEvaluation(w http.ResponseWriter, req *http.Request) {
	envAPIKey, reqBody, err := s.checkGetEvaluationRequest(req)
	if err != nil {
		rest.ReturnFailureResponse(w, err)
		return
	}
	s.publishUser(req.Context(), envAPIKey.Environment.Id, reqBody.Tag, reqBody.User, reqBody.SourceId)
	f, err, _ := s.flightgroup.Do(
		envAPIKey.Environment.Id,
		func() (interface{}, error) {
			return s.getFeatures(req.Context(), envAPIKey.Environment.Id)
		},
	)
	if err != nil {
		rest.ReturnFailureResponse(w, err)
		return
	}
	fs := s.filterOutArchivedFeatures(f.([]*featureproto.Feature))
	features, err := s.getTargetFeatures(fs, reqBody.FeatureID)
	if err != nil {
		rest.ReturnFailureResponse(w, err)
		return
	}
	evaluations, err := s.evaluateFeatures(
		req.Context(),
		reqBody.User,
		features,
		envAPIKey.Environment.Id,
		reqBody.Tag,
	)
	if err != nil {
		s.logger.Error(
			"Failed to evaluate features",
			log.FieldsFromImcomingContext(req.Context()).AddFields(
				zap.Error(err),
				zap.String("environmentID", envAPIKey.Environment.Id),
				zap.String("userId", reqBody.User.Id),
				zap.String("featureId", reqBody.FeatureID),
			)...,
		)
		rest.ReturnFailureResponse(w, errInternal)
		return
	}
	eval, err := s.findEvaluation(evaluations.Evaluations, reqBody.FeatureID)
	if err != nil {
		rest.ReturnFailureResponse(w, err)
		return
	}
	rest.ReturnSuccessResponse(
		w,
		&getEvaluationResponse{
			Evaluation: eval,
		},
	)
}

func (s *gatewayService) getTargetFeatures(fs []*featureproto.Feature, id string) ([]*featureproto.Feature, error) {
	feature, err := s.findFeature(fs, id)
	if err != nil {
		return nil, err
	}
	if len(feature.Prerequisites) > 0 {
		// If we select only the prerequisite feature flags, we have to get them recursively.
		// Thus, we evaluate all features here to avoid complex logic.
		return fs, nil
	}
	return []*featureproto.Feature{feature}, nil
}

func (*gatewayService) findFeature(fs []*featureproto.Feature, id string) (*featureproto.Feature, error) {
	for _, f := range fs {
		if f.Id == id {
			return f, nil
		}
	}
	return nil, errFeatureNotFound
}

func (*gatewayService) findEvaluation(
	evals []*featureproto.Evaluation,
	id string,
) (*featureproto.Evaluation, error) {
	for _, e := range evals {
		if e.FeatureId == id {
			return e, nil
		}
	}
	return nil, errEvaluationNotFound
}

func (s *gatewayService) checkGetEvaluationsRequest(
	req *http.Request,
) (*accountproto.EnvironmentAPIKey, getEvaluationsRequest, error) {
	if req.Method != http.MethodPost {
		return nil, getEvaluationsRequest{}, errInvalidHttpMethod
	}
	envAPIKey, err := s.checkRequest(req.Context(), req)
	if err != nil {
		return nil, getEvaluationsRequest{}, err
	}
	var body getEvaluationsRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		s.logger.Error(
			"Failed to decode request body",
			log.FieldsFromImcomingContext(req.Context()).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, getEvaluationsRequest{}, errInternal
	}
	if err := s.validateGetEvaluationsRequest(&body); err != nil {
		return nil, getEvaluationsRequest{}, err
	}
	return envAPIKey, body, nil
}

func (s *gatewayService) checkGetEvaluationRequest(
	req *http.Request,
) (*accountproto.EnvironmentAPIKey, getEvaluationRequest, error) {
	if req.Method != http.MethodPost {
		return nil, getEvaluationRequest{}, errInvalidHttpMethod
	}
	envAPIKey, err := s.checkRequest(req.Context(), req)
	if err != nil {
		return nil, getEvaluationRequest{}, err
	}
	var body getEvaluationRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		s.logger.Error(
			"Failed to decode request body",
			log.FieldsFromImcomingContext(req.Context()).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, getEvaluationRequest{}, errInternal
	}
	if err := s.validateGetEvaluationRequest(&body); err != nil {
		return nil, getEvaluationRequest{}, err
	}
	return envAPIKey, body, nil
}

func (*gatewayService) validateGetEvaluationsRequest(body *getEvaluationsRequest) error {
	if body.Tag == "" {
		return errTagRequired
	}
	if body.User == nil {
		return errUserRequired
	}
	if body.User.Id == "" {
		return errUserIDRequired
	}
	return nil
}

func (*gatewayService) validateGetEvaluationRequest(body *getEvaluationRequest) error {
	if body.Tag == "" {
		return errTagRequired
	}
	if body.User == nil {
		return errUserRequired
	}
	if body.User.Id == "" {
		return errUserIDRequired
	}
	if body.FeatureID == "" {
		return errFeatureIDRequired
	}
	return nil
}

func (s *gatewayService) publishUser(
	ctx context.Context,
	environmentId, tag string,
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

func (s *gatewayService) publishUserEvent(
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

func (s *gatewayService) checkRequest(
	ctx context.Context,
	req *http.Request,
) (*accountproto.EnvironmentAPIKey, error) {
	if isContextCanceled(ctx) {
		s.logger.Warn(
			"Request was canceled",
			log.FieldsFromImcomingContext(ctx)...,
		)
		return nil, errContextCanceled
	}
	envAPIKey, err := s.findEnvironmentAPIKey(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := s.checkEnvironmentAPIKey(envAPIKey, accountproto.APIKey_SDK); err != nil {
		return nil, err
	}
	return envAPIKey, nil
}

func (*gatewayService) checkEnvironmentAPIKey(
	environmentAPIKey *accountproto.EnvironmentAPIKey,
	role accountproto.APIKey_Role,
) error {
	if environmentAPIKey.ApiKey.Role != role {
		return errBadRole
	}
	if environmentAPIKey.EnvironmentDisabled {
		return errDisabledAPIKey
	}
	if environmentAPIKey.ApiKey.Disabled {
		return errDisabledAPIKey
	}
	return nil
}

func (s *gatewayService) findEnvironmentAPIKey(
	ctx context.Context,
	req *http.Request,
) (*accountproto.EnvironmentAPIKey, error) {
	id := req.Header.Get(authorizationKey)
	if id == "" {
		return nil, errMissingAPIKey
	}
	envAPIKey, err := getEnvironmentAPIKeyFromCache(
		ctx,
		id,
		s.environmentAPIKeyCache,
		callerGatewayService,
		cacheLayerInMemory,
	)
	if err == nil {
		return envAPIKey, nil
	}
	k, err, _ := s.flightgroup.Do(
		id,
		func() (interface{}, error) {
			return s.getEnvironmentAPIKey(
				ctx,
				id,
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

func (s *gatewayService) getEnvironmentAPIKey(
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
			return nil, errInvalidAPIKey
		}
		logger.Error(
			"Failed to get environment APIKey from account service",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, errInternal
	}
	envAPIKey := resp.EnvironmentApiKey
	if err := environmentAPIKeyCache.Put(envAPIKey); err != nil {
		logger.Error(
			"Failed to cache environment APIKey",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", envAPIKey.Environment.Id),
			)...,
		)
	}
	return envAPIKey, nil
}

func (s *gatewayService) evaluateFeatures(
	ctx context.Context,
	user *userproto.User,
	features []*featureproto.Feature,
	environmentId, tag string,
) (*featureproto.UserEvaluations, error) {
	mapIDs := make(map[string]struct{})
	for _, f := range features {
		feature := &featuredomain.Feature{Feature: f}
		for _, id := range feature.ListSegmentIDs() {
			mapIDs[id] = struct{}{}
		}
	}
	mapSegmentUsers, err := s.listSegmentUsers(ctx, user.Id, mapIDs, environmentId)
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
	userEvaluations, err := featuredomain.EvaluateFeatures(features, user, mapSegmentUsers, tag)
	if err != nil {
		s.logger.Error(
			"Failed to evaluate",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentId),
			)...,
		)
	}
	return userEvaluations, nil
}

func (s *gatewayService) listSegmentUsers(
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

func (s *gatewayService) segmentFlightID(environmentId, segmentID string) string {
	return fmt.Sprintf("%s:%s", environmentId, segmentID)
}

func (s *gatewayService) getSegmentUsers(
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
		return nil, errInternal
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

func (s *gatewayService) getSegmentUsersFromCache(
	segmentID, environmentId string,
) ([]*featureproto.SegmentUser, error) {
	segment, err := s.segmentUsersCache.Get(segmentID, environmentId)
	if err == nil {
		return segment.Users, nil
	}
	return nil, err
}

func (s *gatewayService) getFeatures(
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
		return nil, errInternal
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

func (s *gatewayService) getFeaturesFromCache(
	ctx context.Context,
	environmentId string,
) (*featureproto.Features, error) {
	features, err := s.featuresCache.Get(environmentId)
	if err == nil {
		restCacheCounter.WithLabelValues(callerGatewayService, typeFeatures, cacheLayerExternal, codeHit).Inc()
		return features, nil
	}
	restCacheCounter.WithLabelValues(callerGatewayService, typeFeatures, cacheLayerExternal, codeMiss).Inc()
	return nil, err
}

func (s *gatewayService) listFeatures(
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

func (s *gatewayService) registerEvents(w http.ResponseWriter, req *http.Request) {
	envAPIKey, reqBody, err := s.checkRegisterEvents(req)
	if err != nil {
		rest.ReturnFailureResponse(w, err)
		return
	}
	errs := make(map[string]*registerEventsResponseError)
	goalMessages := make([]publisher.Message, 0)
	evaluationMessages := make([]publisher.Message, 0)
	metricsMessages := make([]publisher.Message, 0)
	publish := func(p publisher.Publisher, messages []publisher.Message, typ string) {
		errors := p.PublishMulti(req.Context(), messages)
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
				log.FieldsFromImcomingContext(req.Context()).AddFields(
					zap.Error(err),
					zap.String("environmentID", envAPIKey.Environment.Id),
					zap.String("id", id),
				)...,
			)
			errs[id] = &registerEventsResponseError{
				Retriable: retriable,
				Message:   "Failed to publish event",
			}
		}
		restEventCounter.WithLabelValues(callerGatewayService, typ, codeNonRepeatableError).Add(nonRepeateableErrors)
		restEventCounter.WithLabelValues(callerGatewayService, typ, codeRepeatableError).Add(repeatableErrors)
		restEventCounter.WithLabelValues(callerGatewayService, typ, codeOK).Add(float64(len(messages) - len(errors)))
	}
	for _, event := range reqBody.Events {
		event.EnvironmentNamespace = envAPIKey.Environment.Id
		if event.ID == "" {
			rest.ReturnFailureResponse(w, errMissingEventID)
			return
		}
		switch event.Type {
		case GoalEventType:
			goal, errCode, err := s.getGoalEvent(req.Context(), event)
			if err != nil {
				restEventCounter.WithLabelValues(callerGatewayService, typeMetrics, errCode).Inc()
				errs[event.ID] = &registerEventsResponseError{
					Retriable: false,
					Message:   err.Error(),
				}
				continue
			}
			goalAny, err := ptypes.MarshalAny(goal)
			if err != nil {
				restEventCounter.WithLabelValues(callerGatewayService, typeGoal, codeMarshalAnyFailed).Inc()
				errs[event.ID] = &registerEventsResponseError{
					Retriable: false,
					Message:   err.Error(),
				}
				continue
			}
			goalMessages = append(goalMessages, &eventproto.Event{
				Id:                   event.ID,
				Event:                goalAny,
				EnvironmentNamespace: event.EnvironmentNamespace,
			})
		case EvaluationEventType:
			eval, errCode, err := s.getEvaluationEvent(req.Context(), event)
			if err != nil {
				restEventCounter.WithLabelValues(callerGatewayService, typeMetrics, errCode).Inc()
				errs[event.ID] = &registerEventsResponseError{
					Retriable: false,
					Message:   err.Error(),
				}
				continue
			}
			evalAny, err := ptypes.MarshalAny(eval)
			if err != nil {
				restEventCounter.WithLabelValues(callerGatewayService, typeEvaluation, codeMarshalAnyFailed).Inc()
				errs[event.ID] = &registerEventsResponseError{
					Retriable: false,
					Message:   err.Error(),
				}
				continue
			}
			evaluationMessages = append(evaluationMessages, &eventproto.Event{
				Id:                   event.ID,
				Event:                evalAny,
				EnvironmentNamespace: event.EnvironmentNamespace,
			})
		case MetricsEventType:
			metrics, errCode, err := s.getMetricsEvent(req.Context(), event)
			if err != nil {
				restEventCounter.WithLabelValues(callerGatewayService, typeMetrics, errCode).Inc()
				errs[event.ID] = &registerEventsResponseError{
					Retriable: false,
					Message:   err.Error(),
				}
				continue
			}
			metricsAny, err := ptypes.MarshalAny(metrics)
			if err != nil {
				restEventCounter.WithLabelValues(callerGatewayService, typeMetrics, codeMarshalAnyFailed).Inc()
				errs[event.ID] = &registerEventsResponseError{
					Retriable: false,
					Message:   err.Error(),
				}
				continue
			}
			metricsMessages = append(metricsMessages, &eventproto.Event{
				Id:                   event.ID,
				Event:                metricsAny,
				EnvironmentNamespace: event.EnvironmentNamespace,
			})
		default:
			errs[event.ID] = &registerEventsResponseError{
				Retriable: false,
				Message:   errInvalidType.Error(),
			}
			restEventCounter.WithLabelValues(callerGatewayService, typeUnknown, codeInvalidType).Inc()
			continue
		}
	}
	publish(s.goalPublisher, goalMessages, typeGoal)
	publish(s.evaluationPublisher, evaluationMessages, typeEvaluation)
	publish(s.metricsPublisher, metricsMessages, typeMetrics)
	if len(errs) > 0 {
		if s.containsInvalidTimestampError(errs) {
			restEventCounter.WithLabelValues(callerGatewayService, typeRegisterEvent, codeInvalidTimestampRequest).Inc()
		}
	} else {
		restEventCounter.WithLabelValues(callerGatewayService, typeRegisterEvent, codeOK).Inc()
	}
	rest.ReturnSuccessResponse(
		w,
		registerEventsResponse{Errors: errs},
	)
}

/* Because we got the following error, `nolint` is added. After solving it, we'll remove it.

pkg/gateway/api/api.go:829:47: cannot use ev
(variable of type *"github.com/bucketeer-io/bucketeer/proto/event/client".GoalEvent)
as protoreflect.ProtoMessage value in argument to protojson.Unmarshal:
missing method ProtoReflect (typecheck)
                        if err := protojson.Unmarshal(event.Event, ev); err != nil {
                                                                   ^
*/
//nolint:typecheck
func (s *gatewayService) getGoalEvent(ctx context.Context, event event) (*eventproto.GoalEvent, string, error) {
	ev := &eventproto.GoalEvent{}
	if err := protojson.Unmarshal(event.Event, ev); err != nil {
		s.logger.Error(
			"Failed to extract goal event",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", event.ID),
			)...,
		)
		return nil, codeUnmarshalFailed, errUnmarshalFailed
	}
	errorCode, err := s.validateGoalEvent(ctx, event.ID, ev.Timestamp)
	if err != nil {
		return nil, errorCode, err
	}
	return ev, "", nil
}

/* Because we got the following error, `nolint` is added. After solving it, we'll remove it.

pkg/gateway/api/api.go:829:47: cannot use ev
(variable of type *"github.com/bucketeer-io/bucketeer/proto/event/client".GoalEvent)
as protoreflect.ProtoMessage value in argument to protojson.Unmarshal:
missing method ProtoReflect (typecheck)
                        if err := protojson.Unmarshal(event.Event, ev); err != nil {
                                                                   ^
*/
//nolint:typecheck
func (s *gatewayService) getEvaluationEvent(
	ctx context.Context,
	event event,
) (*eventproto.EvaluationEvent, string, error) {
	ev := &eventproto.EvaluationEvent{}
	if err := protojson.Unmarshal(event.Event, ev); err != nil {
		s.logger.Error(
			"Failed to extract evaluation event",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", event.ID),
			)...,
		)
		return nil, codeUnmarshalFailed, errUnmarshalFailed
	}
	errorCode, err := s.validateEvaluationEvent(ctx, event.ID, ev.Timestamp)
	if err != nil {
		return nil, errorCode, err
	}
	return ev, "", nil
}

/* Because we got the following error, `nolint` is added. After solving it, we'll remove it.

pkg/gateway/api/api.go:829:47: cannot use ev
(variable of type *"github.com/bucketeer-io/bucketeer/proto/event/client".GoalEvent)
as protoreflect.ProtoMessage value in argument to protojson.Unmarshal:
missing method ProtoReflect (typecheck)
                        if err := protojson.Unmarshal(event.Event, ev); err != nil {
                                                                   ^
*/
//nolint:typecheck
func (s *gatewayService) getMetricsEvent(
	ctx context.Context,
	event event,
) (*eventproto.MetricsEvent, string, error) {
	metricsEvt := &metricsEvent{}
	if err := json.Unmarshal(event.Event, metricsEvt); err != nil {
		s.logger.Error(
			"Failed to extract metrics event",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", event.ID),
			)...,
		)
		return nil, codeUnmarshalFailed, errUnmarshalFailed
	}
	errorCode, err := s.validateMetricsEvent(ctx, event.ID)
	if err != nil {
		return nil, errorCode, err
	}
	var eventAny *anypb.Any
	switch metricsEvt.Type {
	case latencyMetricsEventType:
		latency := &latencyMetricsEvent{}
		if err := json.Unmarshal(metricsEvt.Event, latency); err != nil {
			s.logger.Error(
				"Failed to extract latencyMetrics event",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", event.ID),
				)...,
			)
			return nil, codeUnmarshalFailed, errUnmarshalFailed
		}
		apiID := getApiID(latency.ApiId)
		eventAny, err = ptypes.MarshalAny(&eventproto.LatencyMetricsEvent{
			ApiId:    apiID,
			Labels:   latency.Labels,
			Duration: ptypes.DurationProto(latency.Duration),
		})
		if err != nil {
			return nil, codeMarshalAnyFailed, err
		}
	case sizeMetricsEventType:
		size := &sizeMetricsEvent{}
		if err := json.Unmarshal(metricsEvt.Event, size); err != nil {
			s.logger.Error(
				"Failed to extract sizeMetrics event",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", event.ID),
				)...,
			)
			return nil, codeUnmarshalFailed, errUnmarshalFailed
		}
		apiID := getApiID(size.ApiId)
		eventAny, err = ptypes.MarshalAny(&eventproto.SizeMetricsEvent{
			ApiId:    apiID,
			Labels:   size.Labels,
			SizeByte: size.SizeByte,
		})
		if err != nil {
			return nil, codeMarshalAnyFailed, err
		}
	case timeoutErrorMetricsEventType:
		timeout := &timeoutErrorMetricsEvent{}
		if err := json.Unmarshal(event.Event, timeout); err != nil {
			s.logger.Error(
				"Failed to extract timeoutErrorMetrics event",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", event.ID),
				)...,
			)
			return nil, codeUnmarshalFailed, errUnmarshalFailed
		}
		eventAny, err = ptypes.MarshalAny(&eventproto.TimeoutErrorMetricsEvent{
			ApiId:  timeout.ApiId,
			Labels: timeout.Labels,
		})
		if err != nil {
			return nil, codeMarshalAnyFailed, err
		}
	case internalErrorMetricsEventType:
		internal := &internalErrorMetricsEvent{}
		if err := json.Unmarshal(event.Event, internal); err != nil {
			s.logger.Error(
				"Failed to extract internalErrorMetrics event",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", event.ID),
				)...,
			)
			return nil, codeUnmarshalFailed, errUnmarshalFailed
		}
		eventAny, err = ptypes.MarshalAny(&eventproto.InternalErrorMetricsEvent{
			ApiId:  internal.ApiId,
			Labels: internal.Labels,
		})
		if err != nil {
			return nil, codeMarshalAnyFailed, err
		}
	case networkErrorMetricsEventType:
		network := &networkErrorMetricsEvent{}
		if err := json.Unmarshal(event.Event, network); err != nil {
			s.logger.Error(
				"Failed to extract networkErrorMetrics event",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", event.ID),
				)...,
			)
			return nil, codeUnmarshalFailed, errUnmarshalFailed
		}
		eventAny, err = ptypes.MarshalAny(&eventproto.NetworkErrorMetricsEvent{
			ApiId:  network.ApiId,
			Labels: network.Labels,
		})
		if err != nil {
			return nil, codeMarshalAnyFailed, err
		}
	case internalSdkErrorMetricsEventType:
		internalSdk := &internalSdkErrorMetricsEvent{}
		if err := json.Unmarshal(event.Event, internalSdk); err != nil {
			s.logger.Error(
				"Failed to extract internalSdkErrorMetrics event",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", event.ID),
				)...,
			)
			return nil, codeUnmarshalFailed, errUnmarshalFailed
		}
		eventAny, err = ptypes.MarshalAny(&eventproto.InternalErrorMetricsEvent{
			ApiId:  internalSdk.ApiId,
			Labels: internalSdk.Labels,
		})
		if err != nil {
			return nil, codeMarshalAnyFailed, err
		}
	default:
		return nil, codeInvalidType, errInvalidType
	}

	return &eventproto.MetricsEvent{
		Timestamp:  metricsEvt.Timestamp,
		Event:      eventAny,
		SdkVersion: metricsEvt.SdkVersion,
	}, "", nil
}

// TODO: Remove this assignment after updating current SDK
func getApiID(apiID eventproto.ApiId) eventproto.ApiId {
	if apiID == eventproto.ApiId_UNKNOWN_API {
		return eventproto.ApiId_GET_EVALUATION
	}
	return apiID
}

func (s *gatewayService) checkRegisterEvents(
	req *http.Request,
) (*accountproto.EnvironmentAPIKey, registerEventsRequest, error) {
	if req.Method != http.MethodPost {
		return nil, registerEventsRequest{}, errInvalidHttpMethod
	}
	envAPIKey, err := s.checkRequest(req.Context(), req)
	if err != nil {
		return nil, registerEventsRequest{}, err
	}
	var body registerEventsRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		if err == io.EOF {
			return nil, registerEventsRequest{}, errBodyRequired
		}
		s.logger.Error(
			"Failed to decode request body",
			log.FieldsFromImcomingContext(req.Context()).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, registerEventsRequest{}, errInternal
	}
	if len(body.Events) == 0 {
		return nil, registerEventsRequest{}, errMissingEvents
	}
	return envAPIKey, body, nil
}

func (s *gatewayService) containsInvalidTimestampError(errs map[string]*registerEventsResponseError) bool {
	for _, v := range errs {
		if v.Message == errInvalidTimestamp.Error() {
			return true
		}
	}
	return false
}

func (s *gatewayService) filterOutArchivedFeatures(fs []*featureproto.Feature) []*featureproto.Feature {
	result := make([]*featureproto.Feature, 0)
	for _, f := range fs {
		if f.Archived {
			continue
		}
		result = append(result, f)
	}
	return result
}

func (s *gatewayService) emptyUserEvaluations() *featureproto.UserEvaluations {
	return &featureproto.UserEvaluations{
		Id:                 "no_evaluations",
		Evaluations:        []*featureproto.Evaluation{},
		CreatedAt:          time.Now().Unix(),
		ArchivedFeatureIds: []string{},
		ForceUpdate:        false,
	}
}
