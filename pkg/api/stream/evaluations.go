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
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rest"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
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

var sseMarshalOpts = protojson.MarshalOptions{EmitUnpopulated: true}

type CheckRequestFunc func(ctx context.Context, req *http.Request) (*accountproto.EnvironmentAPIKey, error)

type EvaluateFunc func(
	ctx context.Context,
	user *userproto.User,
	environmentID, tag string,
	evaluatedAt int64,
) (*featureproto.UserEvaluations, error)

// EvaluationsHandler handles the SSE stream_evaluations endpoint.
type EvaluationsHandler struct {
	dispatcher        *Dispatcher
	heartbeatInterval time.Duration
	evaluate          EvaluateFunc
	checkRequest      CheckRequestFunc
	requestCounter    *prometheus.CounterVec
	logger            *zap.Logger
}

func NewEvaluationsHandler(
	dispatcher *Dispatcher,
	heartbeatInterval time.Duration,
	evaluate EvaluateFunc,
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
		evaluate:          evaluate,
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

	flusher, ok := w.(http.Flusher)
	if !ok {
		rest.ReturnFailureResponse(w, errInternal)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	// Disable proxy buffering so heartbeat and patch events reach the client immediately.
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	ctx := httpReq.Context()

	// Register before the initial `put` so updates dispatched during the PUT
	// computation are delivered as the first patch.
	events, deregister := h.dispatcher.register(envID, req.Tag)
	defer deregister()

	evaluatedAt, err := h.sendInitialPut(ctx, w, flusher, req.User, envID, req.Tag)
	if err != nil {
		h.logger.Error("Failed to send initial put",
			zap.Error(err),
			zap.String("environmentID", envID),
			zap.String("tag", req.Tag),
			zap.String("userID", req.User.Id),
		)
		sendErrorEvent(w, flusher, gatewayproto.StreamErrorEvent_INTERNAL, "evaluation failed")
		return
	}

	ticker := time.NewTicker(h.heartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := sendHeartbeat(w, flusher); err != nil {
				// We cannot send an error event here because
				// write failure means the client connection is already gone.
				return
			}
		case <-events:
			newEvalAt, err := h.sendPatch(ctx, w, flusher, req.User, envID, req.Tag, evaluatedAt)
			if err != nil {
				h.logger.Error("Failed to send patch",
					zap.Error(err),
					zap.String("environmentID", envID),
					zap.String("tag", req.Tag),
					zap.String("userID", req.User.Id),
				)
				sendErrorEvent(w, flusher, gatewayproto.StreamErrorEvent_INTERNAL, "evaluation failed")
				return
			}
			evaluatedAt = newEvalAt
		}
	}
}

func (h *EvaluationsHandler) sendInitialPut(
	ctx context.Context,
	w io.Writer,
	flusher http.Flusher,
	user *userproto.User,
	envID, tag string,
) (evaluatedAt int64, err error) {
	evaluatedAt = time.Now().Unix()
	evals, err := h.evaluate(ctx, user, envID, tag, 0)
	if err != nil {
		return 0, err
	}
	evt := &gatewayproto.StreamEvaluationsEvent{
		Evaluations: evals,
	}
	if err := sendSSEEvent(w, flusher, "put", evt); err != nil {
		return 0, err
	}
	return evaluatedAt, nil
}

func (h *EvaluationsHandler) sendPatch(
	ctx context.Context,
	w io.Writer,
	flusher http.Flusher,
	user *userproto.User,
	envID, tag string,
	prevEvaluatedAt int64,
) (newEvaluatedAt int64, err error) {
	newEvaluatedAt = time.Now().Unix()
	evals, err := h.evaluate(ctx, user, envID, tag, prevEvaluatedAt)
	if err != nil {
		return 0, err
	}
	evt := &gatewayproto.StreamEvaluationsEvent{
		Evaluations: evals,
	}
	if err := sendSSEEvent(w, flusher, "patch", evt); err != nil {
		return 0, err
	}
	return newEvaluatedAt, nil
}

func sendSSEEvent(w io.Writer, flusher http.Flusher, eventType string, msg proto.Message) error {
	data, err := sseMarshalOpts.Marshal(msg)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, data); err != nil {
		return err
	}
	flusher.Flush()
	return nil
}

func sendHeartbeat(w io.Writer, flusher http.Flusher) error {
	if _, err := fmt.Fprintf(w, ":\n\n"); err != nil {
		return err
	}
	flusher.Flush()
	return nil
}

func sendErrorEvent(w io.Writer, flusher http.Flusher, code gatewayproto.StreamErrorEvent_Code, msg string) {
	evt := &gatewayproto.StreamErrorEvent{
		Code:    code,
		Message: msg,
	}
	_ = sendSSEEvent(w, flusher, "error", evt)
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
