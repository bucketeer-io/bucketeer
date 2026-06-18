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

package stream

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rest"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

var (
	errInvalidHttpMethod = rest.NewErrStatus(http.StatusMethodNotAllowed, "gateway: invalid http method")
	errInternal          = rest.NewErrStatus(http.StatusInternalServerError, "gateway: internal")
	errTagRequired       = rest.NewErrStatus(http.StatusBadRequest, "gateway: tag is required")
	errUserRequired      = rest.NewErrStatus(http.StatusBadRequest, "gateway: user is required")
	errUserIDRequired    = rest.NewErrStatus(http.StatusBadRequest, "gateway: user id is required")
)

// sseUnmarshalOpts ignores unknown fields like the polling endpoints' encoding/json decoder.
var sseUnmarshalOpts = protojson.UnmarshalOptions{DiscardUnknown: true}

type CheckRequestFunc func(ctx context.Context, req *http.Request) (*accountproto.EnvironmentAPIKey, error)

// EvaluationsHandler handles the SSE stream_evaluations endpoint.
type EvaluationsHandler struct {
	dispatcher        *Dispatcher
	heartbeatInterval time.Duration
	checkRequest      CheckRequestFunc
	requestCounter    *prometheus.CounterVec
	logger            *zap.Logger
}

func NewEvaluationsHandler(
	dispatcher *Dispatcher,
	heartbeatInterval time.Duration,
	checkRequest CheckRequestFunc,
	requestCounter *prometheus.CounterVec,
	logger *zap.Logger,
) *EvaluationsHandler {
	if heartbeatInterval <= 0 {
		heartbeatInterval = 25 * time.Second
	}

	return &EvaluationsHandler{
		dispatcher:        dispatcher,
		heartbeatInterval: heartbeatInterval,
		checkRequest:      checkRequest,
		requestCounter:    requestCounter,
		logger:            logger.Named("stream-evaluations"),
	}
}

func (h *EvaluationsHandler) Handle(w http.ResponseWriter, httpReq *http.Request) {
	envAPIKey, req, err := h.checkStreamEvaluationsRequest(httpReq)
	if err != nil {
		rest.ReturnFailureResponse(w, err)
		return
	}
	envID := envAPIKey.Environment.Id
	h.requestCounter.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envID, envAPIKey.Environment.UrlCode,
		methodStreamEvaluations, req.SourceId.String()).Inc()

	// TODO: Prepare http.Flusher and write the SSE response headers.

	ctx := httpReq.Context()

	// Register before the initial `put` so updates dispatched during the PUT
	// computation are delivered as the first patch.
	events, deregister := h.dispatcher.register(envID, req.Tag)
	defer deregister()

	// TODO: Send the initial `put` event (full snapshot, mirroring getEvaluations).

	ticker := time.NewTicker(h.heartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// TODO: Write an SSE heartbeat comment and flush.
		case <-events:
			// TODO: Send a `patch` event.
		}
	}
}

func (h *EvaluationsHandler) checkStreamEvaluationsRequest(
	httpReq *http.Request,
) (*accountproto.EnvironmentAPIKey, *gatewayproto.StreamEvaluationsRequest, error) {
	if httpReq.Method != http.MethodPost {
		return nil, nil, errInvalidHttpMethod
	}
	rawBody, err := io.ReadAll(httpReq.Body)
	if err != nil {
		h.logger.Error(
			"Failed to read stream evaluations request body",
			log.FieldsFromIncomingContext(httpReq.Context()).AddFields(zap.Error(err))...,
		)
		return nil, nil, errInternal
	}
	req := &gatewayproto.StreamEvaluationsRequest{}
	if err := sseUnmarshalOpts.Unmarshal(rawBody, req); err != nil {
		h.logger.Error(
			"Failed to decode stream evaluations request body",
			log.FieldsFromIncomingContext(httpReq.Context()).AddFields(zap.Error(err))...,
		)
		return nil, nil, errInternal
	}
	envAPIKey, err := h.checkRequest(httpReq.Context(), httpReq)
	if err != nil {
		h.logger.Error("Failed to check StreamEvaluations request",
			zap.Error(err),
			zap.String("tag", req.Tag),
			zap.Any("user", req.User),
			zap.Any("sourceId", req.SourceId),
		)
		return nil, nil, err
	}
	if err := validateStreamEvaluationsRequest(req); err != nil {
		h.logger.Error("Failed to validate StreamEvaluations request",
			zap.Error(err),
			zap.String("tag", req.Tag),
			zap.Any("user", req.User),
			zap.Any("sourceId", req.SourceId),
		)
		return nil, nil, err
	}
	return envAPIKey, req, nil
}

func validateStreamEvaluationsRequest(req *gatewayproto.StreamEvaluationsRequest) error {
	if req.Tag == "" {
		return errTagRequired
	}
	if req.User == nil {
		return errUserRequired
	}
	if req.User.Id == "" {
		return errUserIDRequired
	}
	return nil
}
