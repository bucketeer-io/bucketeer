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
	"net/http"
	"strconv"

	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
)

const (
	urlParamKeyAPIKey    = "apikey"
	urlParamKeyUserID    = "userid"
	urlParamKeyGoalID    = "goalid"
	urlParamKeyTag       = "tag"
	urlParamKeyTimestamp = "timestamp"
	urlParamKeyValue     = "value"
)

var (
	errAPIKeyEmpty      = errors.New("gateway: api key is empty")
	errUserIDEmpty      = errors.New("gateway: user id is empty")
	errGoalIDEmpty      = errors.New("gateway: goal id is empty")
	errTagEmpty         = errors.New("gateway: tag is empty")
	errTimestampEmpty   = errors.New("gateway: timestamp is empty")
	errTimestampInvalid = errors.New("gateway: timestamp is invalid")
	errValueInvalid     = errors.New("gateway: value is invalid")
)

type TrackHandler struct {
	accountClient          accountclient.Client
	goalBatchPublisher     publisher.Publisher
	environmentAPIKeyCache cachev3.EnvironmentAPIKeyCache
	flightgroup            singleflight.Group
	opts                   *options
	logger                 *zap.Logger
}

func NewTrackHandler(
	accountClient accountclient.Client,
	gbp publisher.Publisher,
	v3Cache cache.MultiGetCache,
	opts ...Option) *TrackHandler {

	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.metrics != nil {
		registerMetrics(options.metrics)
	}
	return &TrackHandler{
		accountClient:          accountClient,
		goalBatchPublisher:     gbp,
		environmentAPIKeyCache: cachev3.NewEnvironmentAPIKeyCache(v3Cache),
		opts:                   &options,
		logger:                 options.logger.Named("trackhandler"),
	}
}

func (h *TrackHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if isContextCanceled(ctx) {
		h.logger.Warn(
			"Request was canceled",
			log.FieldsFromImcomingContext(ctx)...,
		)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	params, err := h.validateParams(req)
	if err != nil {
		h.logger.Warn(
			"Invalid url parameters",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	envAPIKey, err := h.getEnvironmentAPIKey(ctx, params.apiKey)
	if err != nil {
		h.logger.Error(
			"Failed to get environment api key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", envAPIKey.EnvironmentNamespace),
				zap.String("apiKey", params.apiKey),
			)...,
		)
		eventCounter.WithLabelValues(callerTrackHandler, typeHTTPTrack, codeNonRepeatableError).Inc()
		if err == ErrInvalidAPIKey {
			resp.WriteHeader(http.StatusForbidden)
			return
		}
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := checkEnvironmentAPIKey(envAPIKey, accountproto.APIKey_SDK); err != nil {
		eventCounter.WithLabelValues(callerTrackHandler, typeHTTPTrack, codeNonRepeatableError).Inc()
		resp.WriteHeader(http.StatusForbidden)
		return
	}
	goalBatchEvent, err := h.createGoalBatchEvent(envAPIKey.EnvironmentNamespace, params)
	if err != nil {
		h.logger.Error(
			"Failed to create goal batch event",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", envAPIKey.EnvironmentNamespace),
				zap.String("apiKey", params.apiKey),
				zap.String("userId", params.userID),
				zap.String("goalId", params.goalID),
				zap.Int64("timestamp", params.timestamp),
				zap.Float64("value", params.value),
				zap.String("tag", params.tag),
			)...,
		)
		eventCounter.WithLabelValues(callerTrackHandler, typeHTTPTrack, codeNonRepeatableError).Inc()
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := h.goalBatchPublisher.Publish(ctx, goalBatchEvent); err != nil {
		eventCounter.WithLabelValues(callerTrackHandler, typeHTTPTrack, codeNonRepeatableError).Inc()
		h.logger.Error(
			"Failed to publish event",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", envAPIKey.EnvironmentNamespace),
			)...,
		)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	eventCounter.WithLabelValues(callerTrackHandler, typeHTTPTrack, codeOK).Inc()
	resp.WriteHeader(http.StatusOK)
}

type params struct {
	apiKey    string
	userID    string
	goalID    string
	tag       string
	timestamp int64
	value     float64
}

func (h *TrackHandler) validateParams(req *http.Request) (*params, error) {
	params := &params{}
	q := req.URL.Query()
	apikey := q.Get(urlParamKeyAPIKey)
	if apikey == "" {
		eventCounter.WithLabelValues(callerTrackHandler, typeHTTPTrack, codeInvalidURLParams).Inc()
		return nil, errAPIKeyEmpty
	}
	params.apiKey = apikey
	userID := q.Get(urlParamKeyUserID)
	if userID == "" {
		eventCounter.WithLabelValues(callerTrackHandler, typeHTTPTrack, codeInvalidURLParams).Inc()
		return nil, errUserIDEmpty
	}
	params.userID = userID
	goalID := q.Get(urlParamKeyGoalID)
	if goalID == "" {
		eventCounter.WithLabelValues(callerTrackHandler, typeHTTPTrack, codeInvalidURLParams).Inc()
		return nil, errGoalIDEmpty
	}
	params.goalID = goalID
	tag := q.Get(urlParamKeyTag)
	if tag == "" {
		eventCounter.WithLabelValues(callerTrackHandler, typeHTTPTrack, codeInvalidURLParams).Inc()
		return nil, errTagEmpty
	}
	params.tag = tag
	timestampStr := q.Get(urlParamKeyTimestamp)
	if timestampStr == "" {
		eventCounter.WithLabelValues(callerTrackHandler, typeHTTPTrack, codeInvalidURLParams).Inc()
		return nil, errTimestampEmpty
	}
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		eventCounter.WithLabelValues(callerTrackHandler, typeHTTPTrack, codeInvalidURLParams).Inc()
		return nil, errTimestampInvalid
	}
	if !validateTimestamp(timestamp, h.opts.oldestEventTimestamp, h.opts.furthestEventTimestamp) {
		eventCounter.WithLabelValues(callerTrackHandler, typeHTTPTrack, codeInvalidTimestamp).Inc()
		return nil, errTimestampInvalid
	}
	params.timestamp = timestamp
	valueStr := q.Get(urlParamKeyValue)
	if valueStr != "" {
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			eventCounter.WithLabelValues(callerTrackHandler, typeHTTPTrack, codeInvalidURLParams).Inc()
			return nil, errValueInvalid
		}
		params.value = value
	}
	return params, nil
}

func (h *TrackHandler) createGoalBatchEvent(environmentNamespace string, params *params) (*eventproto.Event, error) {
	goalBatchEvent := &eventproto.GoalBatchEvent{
		UserId: params.userID,
		UserGoalEventsOverTags: []*eventproto.UserGoalEventsOverTag{{
			Tag: params.tag,
			UserGoalEvents: []*eventproto.UserGoalEvent{{
				Timestamp: params.timestamp,
				GoalId:    params.goalID,
				Value:     params.value,
			}},
		}},
	}
	any, err := ptypes.MarshalAny(goalBatchEvent)
	if err != nil {
		return nil, err
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &eventproto.Event{
		Id:                   id.String(),
		Event:                any,
		EnvironmentNamespace: environmentNamespace,
	}, nil
}

func (h *TrackHandler) getEnvironmentAPIKey(ctx context.Context, id string) (*accountproto.EnvironmentAPIKey, error) {
	k, err, _ := h.flightgroup.Do(
		environmentAPIKeyFlightID(id),
		func() (interface{}, error) {
			return getEnvironmentAPIKey(
				ctx,
				id,
				h.accountClient,
				h.environmentAPIKeyCache,
				callerTrackHandler,
				h.logger,
			)
		},
	)
	if err != nil {
		return nil, err
	}
	envAPIKey := k.(*accountproto.EnvironmentAPIKey)
	return envAPIKey, nil
}
