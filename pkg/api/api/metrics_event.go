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
	"errors"
	"strings"

	"go.uber.org/zap"

	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
)

const (
	ErrRedirectionRequest        = "ErrRedirection"
	ErrorTypeBadRequest          = "BadRequest"
	ErrorTypeUnauthenticated     = "Unauthenticated"
	ErrorTypeForbidden           = "Forbidden"
	ErrorTypeNotFound            = "NotFound"
	ErrPayloadTooLargeRequest    = "ErrPayloadTooLarge"
	ErrorTypeClientClosedRequest = "ClientClosedRequest"
	ErrorTypeInternalServerError = "InternalServerError"
	ErrorTypeServiceUnavailable  = "ServiceUnavailable"
	ErrorTypeTimeout             = "Timeout"
	ErrorTypeInternal            = "Internal"
	ErrorTypeNetwork             = "Network"
	ErrorTypeSDKInternal         = "SDKInternal"
	ErrorTypeUnknown             = "Unknown"
)

var (
	MetricsSaveErrUnknownEvent    = errors.New("gateway: unknown metrics event")
	MetricsSaveErrInvalidDuration = errors.New("gateway: metrics event has invalid duration")
	MetricsSaveErrUnknownApiId    = errors.New("gateway: metrics event has unknown api id")

	getEvaluationLatencyMetricsEventP     = &eventproto.GetEvaluationLatencyMetricsEvent{}
	getEvaluationSizeMetricsEventP        = &eventproto.GetEvaluationSizeMetricsEvent{}
	timeoutErrorCountMetricsEventP        = &eventproto.TimeoutErrorCountMetricsEvent{}
	internalErrorCountMetricsEventP       = &eventproto.InternalErrorCountMetricsEvent{}
	latencyMetricsEventP                  = &eventproto.LatencyMetricsEvent{}
	sizeMetricsEventP                     = &eventproto.SizeMetricsEvent{}
	redirectionRequestExceptionEventP     = &eventproto.RedirectionRequestExceptionEvent{}
	badRequestErrorMetricsEventP          = &eventproto.BadRequestErrorMetricsEvent{}
	unauthorizedErrorMetricsEventP        = &eventproto.UnauthorizedErrorMetricsEvent{}
	forbiddenErrorMetricsEventP           = &eventproto.ForbiddenErrorMetricsEvent{}
	notFoundErrorMetricsEventP            = &eventproto.NotFoundErrorMetricsEvent{}
	payloadTooLargeExceptionEventP        = &eventproto.PayloadTooLargeExceptionEvent{}
	clientClosedRequestErrorMetricsEventP = &eventproto.ClientClosedRequestErrorMetricsEvent{}
	internalServerErrorMetricsEventP      = &eventproto.InternalServerErrorMetricsEvent{}
	serviceUnavailableErrorMetricsEventP  = &eventproto.ServiceUnavailableErrorMetricsEvent{}
	timeoutErrorMetricsEventP             = &eventproto.TimeoutErrorMetricsEvent{}
	internalErrorMetricsEventP            = &eventproto.InternalErrorMetricsEvent{}
	networkErrorMetricsEventP             = &eventproto.NetworkErrorMetricsEvent{}
	internalSdkErrorMetricsEventP         = &eventproto.InternalSdkErrorMetricsEvent{}
	unknownErrorMetricsEventP             = &eventproto.UnknownErrorMetricsEvent{}
)

func (s *grpcGatewayService) saveMetricsEventsAsync(
	metricsEvents []*eventproto.MetricsEvent, projectID, environmentUrlCode string,
) {
	// TODO: using buffered channel to reduce the number of go routines
	go func() {
		for i := range metricsEvents {
			if err := s.saveMetrics(metricsEvents[i], projectID, environmentUrlCode); err != nil {
				s.logger.Error("Failed to store metrics event to prometheus client", zap.Error(err))
				eventCounter.WithLabelValues(callerGatewayService, typeMetrics, codeNonRepeatableError).Inc()
			} else {
				eventCounter.WithLabelValues(callerGatewayService, typeMetrics, codeOK).Inc()
			}
		}
	}()
}

func (s *grpcGatewayService) saveMetrics(event *eventproto.MetricsEvent, projectID, environmentUrlCode string) error {
	// TODO: Remove after deleting the api-gateway REST server
	if event.Event.MessageIs(getEvaluationLatencyMetricsEventP) {
		return s.saveGetEvaluationLatencyMetricsEvent(event, environmentUrlCode)
	}
	// TODO: Remove after deleting the api-gateway REST server
	if event.Event.MessageIs(getEvaluationSizeMetricsEventP) {
		return s.saveGetEvaluationSizeMetricsEvent(event, environmentUrlCode)
	}
	// TODO: Remove after deleting the api-gateway REST server
	if event.Event.MessageIs(timeoutErrorCountMetricsEventP) {
		return s.saveTimeoutErrorCountMetricsEvent(event, environmentUrlCode)
	}
	// TODO: Remove after deleting the api-gateway REST server
	if event.Event.MessageIs(internalErrorCountMetricsEventP) {
		return s.saveInternalErrorCountMetricsEvent(event, environmentUrlCode)
	}
	if event.Event.MessageIs(latencyMetricsEventP) {
		return s.saveLatencyMetricsEvent(event, projectID, environmentUrlCode)
	}
	if event.Event.MessageIs(sizeMetricsEventP) {
		return s.saveSizeMetricsEvent(event, projectID, environmentUrlCode)
	}
	if event.Event.MessageIs(badRequestErrorMetricsEventP) {
		return s.saveBadRequestError(event, projectID, environmentUrlCode)
	}
	if event.Event.MessageIs(redirectionRequestExceptionEventP) {
		return s.saveRedirectionRequestError(event, projectID, environmentUrlCode)
	}
	if event.Event.MessageIs(unauthorizedErrorMetricsEventP) {
		return s.saveUnauthorizedError(event, projectID, environmentUrlCode)
	}
	if event.Event.MessageIs(forbiddenErrorMetricsEventP) {
		return s.saveForbiddenError(event, projectID, environmentUrlCode)
	}
	if event.Event.MessageIs(notFoundErrorMetricsEventP) {
		return s.saveNotFoundError(event, projectID, environmentUrlCode)
	}
	if event.Event.MessageIs(payloadTooLargeExceptionEventP) {
		return s.payloadTooLargeRequestError(event, projectID, environmentUrlCode)
	}
	if event.Event.MessageIs(clientClosedRequestErrorMetricsEventP) {
		return s.saveClientClosedRequestError(event, projectID, environmentUrlCode)
	}
	if event.Event.MessageIs(internalServerErrorMetricsEventP) {
		return s.saveInternalServerError(event, projectID, environmentUrlCode)
	}
	if event.Event.MessageIs(serviceUnavailableErrorMetricsEventP) {
		return s.saveServiceUnavailableError(event, projectID, environmentUrlCode)
	}
	if event.Event.MessageIs(timeoutErrorMetricsEventP) {
		return s.saveTimeoutError(event, projectID, environmentUrlCode)
	}
	if event.Event.MessageIs(internalErrorMetricsEventP) {
		return s.saveInternalError(event, projectID, environmentUrlCode)
	}
	if event.Event.MessageIs(networkErrorMetricsEventP) {
		return s.saveNetworkError(event, projectID, environmentUrlCode)
	}
	if event.Event.MessageIs(internalSdkErrorMetricsEventP) {
		return s.saveInternalSdkError(event, projectID, environmentUrlCode)
	}
	if event.Event.MessageIs(unknownErrorMetricsEventP) {
		return s.saveUnknownError(event, projectID, environmentUrlCode)
	}
	return MetricsSaveErrUnknownEvent
}

func (s *grpcGatewayService) saveGetEvaluationLatencyMetricsEvent(event *eventproto.MetricsEvent, env string) error {
	ev := &eventproto.GetEvaluationLatencyMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	if ev.Duration == nil {
		s.logger.Warn("Invalid duration. Duration is nil",
			zap.String("environmentUrlCode", env),
			zap.String("sdkVersion", event.SdkVersion),
			zap.String("sourceId", event.SourceId.String()),
			zap.String("tag", ev.Labels["tag"]),
			zap.Any("labels", ev.Labels),
		)
		return MetricsSaveErrInvalidDuration
	}
	var tag, status string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
		status = ev.Labels["state"]
	}
	if err := ev.Duration.CheckValid(); err != nil {
		s.logger.Warn("Invalid duration. Duration is not valid",
			zap.Error(err),
			zap.String("environmentUrlCode", env),
			zap.String("sdkVersion", event.SdkVersion),
			zap.String("sourceId", event.SourceId.String()),
			zap.String("tag", ev.Labels["tag"]),
			zap.Any("labels", ev.Labels),
		)
		return MetricsSaveErrInvalidDuration
	}
	dur := ev.Duration.AsDuration()
	sdkGetEvaluationsLatencyHistogram.WithLabelValues(env, tag, status).Observe(dur.Seconds())
	return nil
}

func (s *grpcGatewayService) saveGetEvaluationSizeMetricsEvent(event *eventproto.MetricsEvent, env string) error {
	ev := &eventproto.GetEvaluationSizeMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	var tag, status string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
		status = ev.Labels["state"]
	}
	sdkGetEvaluationsSizeHistogram.WithLabelValues(env, tag, status).Observe(float64(ev.SizeByte))
	return nil
}

func (s *grpcGatewayService) saveTimeoutErrorCountMetricsEvent(event *eventproto.MetricsEvent, env string) error {
	ev := &eventproto.TimeoutErrorCountMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	sdkTimeoutErrorCounter.WithLabelValues(env, ev.Tag).Inc()
	return nil
}

func (s *grpcGatewayService) saveInternalErrorCountMetricsEvent(event *eventproto.MetricsEvent, env string) error {
	ev := &eventproto.InternalErrorCountMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	sdkInternalErrorCounter.WithLabelValues(env, ev.Tag).Inc()
	return nil
}

func (s *grpcGatewayService) saveLatencyMetricsEvent(event *eventproto.MetricsEvent, projectID, env string) error {
	ev := &eventproto.LatencyMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	// TODO: When updated to the SDK that uses ev.LatencySecond, we must remove the implementation that use ev.Duration.
	if ev.Duration == nil && ev.LatencySecond == 0 {
		return MetricsSaveErrInvalidDuration
	}
	if ev.ApiId == eventproto.ApiId_UNKNOWN_API {
		return MetricsSaveErrUnknownApiId
	}
	var tag string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
	}
	// TODO: When updated to the SDK that uses ev.LatencySecond, we must remove the implementation that use ev.Duration.
	if ev.LatencySecond != 0 {
		sdkLatencyHistogram.WithLabelValues(
			projectID,
			env,
			tag,
			ev.ApiId.String(),
			event.SdkVersion,
			event.SourceId.String(),
		).Observe(ev.LatencySecond)
		return nil
	}
	if err := ev.Duration.CheckValid(); err != nil {
		return MetricsSaveErrInvalidDuration
	}
	dur := ev.Duration.AsDuration()
	sdkLatencyHistogram.WithLabelValues(
		projectID,
		env,
		tag,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Observe(dur.Seconds())
	return nil
}

func (s *grpcGatewayService) saveSizeMetricsEvent(event *eventproto.MetricsEvent, projectID, env string) error {
	ev := &eventproto.SizeMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	if ev.ApiId == eventproto.ApiId_UNKNOWN_API {
		return MetricsSaveErrUnknownApiId
	}
	var tag string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
	}
	sdkSizeHistogram.WithLabelValues(
		projectID,
		env,
		tag,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Observe(float64(ev.SizeByte))
	return nil
}

func (s *grpcGatewayService) saveBadRequestError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeBadRequest
	ev := &eventproto.BadRequestErrorMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveRedirectionRequestError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrRedirectionRequest
	ev := &eventproto.RedirectionRequestExceptionEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveUnauthorizedError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeUnauthenticated
	ev := &eventproto.UnauthorizedErrorMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveForbiddenError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeForbidden
	ev := &eventproto.ForbiddenErrorMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveNotFoundError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeNotFound
	ev := &eventproto.NotFoundErrorMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) payloadTooLargeRequestError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrPayloadTooLargeRequest
	ev := &eventproto.PayloadTooLargeExceptionEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveClientClosedRequestError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeClientClosedRequest
	ev := &eventproto.ClientClosedRequestErrorMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveInternalServerError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeInternalServerError
	ev := &eventproto.InternalServerErrorMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveServiceUnavailableError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeServiceUnavailable
	ev := &eventproto.ServiceUnavailableErrorMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveTimeoutError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeTimeout
	ev := &eventproto.TimeoutErrorMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveInternalError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeInternal
	ev := &eventproto.InternalErrorMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveNetworkError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeNetwork
	ev := &eventproto.NetworkErrorMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveInternalSdkError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeSDKInternal
	ev := &eventproto.InternalSdkErrorMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	// Log the error message for debugging SDK internal errors
	if ev.Labels != nil {
		errorMessage := ev.Labels["error_message"]
		// Skip logging cache not found errors
		if errorMessage != "" && !isCacheNotFoundError(errorMessage) {
			s.logger.Warn("SDK internal error received",
				zap.String("projectID", projectID),
				zap.String("environmentUrlCode", env),
				zap.String("apiId", ev.ApiId.String()),
				zap.String("sdkVersion", event.SdkVersion),
				zap.String("sourceId", event.SourceId.String()),
				zap.String("tag", ev.Labels["tag"]),
				zap.String("errorMessage", errorMessage),
				zap.Any("labels", ev.Labels),
			)
		}
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func isCacheNotFoundError(errorMessage string) bool {
	return strings.Contains(errorMessage, "cache: not found")
}

func (s *grpcGatewayService) saveUnknownError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeUnknown
	ev := &eventproto.UnknownErrorMetricsEvent{}
	if err := event.Event.UnmarshalTo(ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveErrorCount(
	event *eventproto.MetricsEvent,
	projectID, env, errorType string,
	apiID eventproto.ApiId,
	labels map[string]string,
) error {
	if apiID == eventproto.ApiId_UNKNOWN_API {
		return MetricsSaveErrUnknownApiId
	}
	var tag string
	if labels != nil {
		tag = labels["tag"]
	}
	sdkErrorCounter.WithLabelValues(
		projectID,
		env,
		tag,
		errorType,
		apiID.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Inc()
	return nil
}
