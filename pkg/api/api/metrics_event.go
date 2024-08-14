package api

import (
	"errors"

	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"

	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
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
	if ptypes.Is(event.Event, getEvaluationLatencyMetricsEventP) {
		return s.saveGetEvaluationLatencyMetricsEvent(event, environmentUrlCode)
	}
	// TODO: Remove after deleting the api-gateway REST server
	if ptypes.Is(event.Event, getEvaluationSizeMetricsEventP) {
		return s.saveGetEvaluationSizeMetricsEvent(event, environmentUrlCode)
	}
	// TODO: Remove after deleting the api-gateway REST server
	if ptypes.Is(event.Event, timeoutErrorCountMetricsEventP) {
		return s.saveTimeoutErrorCountMetricsEvent(event, environmentUrlCode)
	}
	// TODO: Remove after deleting the api-gateway REST server
	if ptypes.Is(event.Event, internalErrorCountMetricsEventP) {
		return s.saveInternalErrorCountMetricsEvent(event, environmentUrlCode)
	}
	if ptypes.Is(event.Event, latencyMetricsEventP) {
		return s.saveLatencyMetricsEvent(event, projectID, environmentUrlCode)
	}
	if ptypes.Is(event.Event, sizeMetricsEventP) {
		return s.saveSizeMetricsEvent(event, projectID, environmentUrlCode)
	}
	if ptypes.Is(event.Event, badRequestErrorMetricsEventP) {
		return s.saveBadRequestError(event, projectID, environmentUrlCode)
	}
	if ptypes.Is(event.Event, redirectionRequestExceptionEventP) {
		return s.saveRedirectionRequestError(event, projectID, environmentUrlCode)
	}
	if ptypes.Is(event.Event, unauthorizedErrorMetricsEventP) {
		return s.saveUnauthorizedError(event, projectID, environmentUrlCode)
	}
	if ptypes.Is(event.Event, forbiddenErrorMetricsEventP) {
		return s.saveForbiddenError(event, projectID, environmentUrlCode)
	}
	if ptypes.Is(event.Event, notFoundErrorMetricsEventP) {
		return s.saveNotFoundError(event, projectID, environmentUrlCode)
	}
	if ptypes.Is(event.Event, payloadTooLargeExceptionEventP) {
		return s.payloadTooLargeRequestError(event, projectID, environmentUrlCode)
	}
	if ptypes.Is(event.Event, clientClosedRequestErrorMetricsEventP) {
		return s.saveClientClosedRequestError(event, projectID, environmentUrlCode)
	}
	if ptypes.Is(event.Event, internalServerErrorMetricsEventP) {
		return s.saveInternalServerError(event, projectID, environmentUrlCode)
	}
	if ptypes.Is(event.Event, serviceUnavailableErrorMetricsEventP) {
		return s.saveServiceUnavailableError(event, projectID, environmentUrlCode)
	}
	if ptypes.Is(event.Event, timeoutErrorMetricsEventP) {
		return s.saveTimeoutError(event, projectID, environmentUrlCode)
	}
	if ptypes.Is(event.Event, internalErrorMetricsEventP) {
		return s.saveInternalError(event, projectID, environmentUrlCode)
	}
	if ptypes.Is(event.Event, networkErrorMetricsEventP) {
		return s.saveNetworkError(event, projectID, environmentUrlCode)
	}
	if ptypes.Is(event.Event, internalSdkErrorMetricsEventP) {
		return s.saveInternalSdkError(event, projectID, environmentUrlCode)
	}
	if ptypes.Is(event.Event, unknownErrorMetricsEventP) {
		return s.saveUnknownError(event, projectID, environmentUrlCode)
	}
	return MetricsSaveErrUnknownEvent
}

func (s *grpcGatewayService) saveGetEvaluationLatencyMetricsEvent(event *eventproto.MetricsEvent, env string) error {
	ev := &eventproto.GetEvaluationLatencyMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	if ev.Duration == nil {
		return MetricsSaveErrInvalidDuration
	}
	var tag, status string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
		status = ev.Labels["state"]
	}
	dur, err := ptypes.Duration(ev.Duration)
	if err != nil {
		return MetricsSaveErrInvalidDuration
	}
	sdkGetEvaluationsLatencyHistogram.WithLabelValues(env, tag, status).Observe(dur.Seconds())
	return nil
}

func (s *grpcGatewayService) saveGetEvaluationSizeMetricsEvent(event *eventproto.MetricsEvent, env string) error {
	ev := &eventproto.GetEvaluationSizeMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
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
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	sdkTimeoutErrorCounter.WithLabelValues(env, ev.Tag).Inc()
	return nil
}

func (s *grpcGatewayService) saveInternalErrorCountMetricsEvent(event *eventproto.MetricsEvent, env string) error {
	ev := &eventproto.InternalErrorCountMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	sdkInternalErrorCounter.WithLabelValues(env, ev.Tag).Inc()
	return nil
}

func (s *grpcGatewayService) saveLatencyMetricsEvent(event *eventproto.MetricsEvent, projectID, env string) error {
	ev := &eventproto.LatencyMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
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
	dur, err := ptypes.Duration(ev.Duration)
	if err != nil {
		return MetricsSaveErrInvalidDuration
	}
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
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
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
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveRedirectionRequestError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrRedirectionRequest
	ev := &eventproto.RedirectionRequestExceptionEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveUnauthorizedError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeUnauthenticated
	ev := &eventproto.UnauthorizedErrorMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveForbiddenError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeForbidden
	ev := &eventproto.ForbiddenErrorMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveNotFoundError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeNotFound
	ev := &eventproto.NotFoundErrorMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) payloadTooLargeRequestError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrPayloadTooLargeRequest
	ev := &eventproto.PayloadTooLargeExceptionEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveClientClosedRequestError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeClientClosedRequest
	ev := &eventproto.ClientClosedRequestErrorMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveInternalServerError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeInternalServerError
	ev := &eventproto.InternalServerErrorMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveServiceUnavailableError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeServiceUnavailable
	ev := &eventproto.ServiceUnavailableErrorMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveTimeoutError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeTimeout
	ev := &eventproto.TimeoutErrorMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveInternalError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeInternal
	ev := &eventproto.InternalErrorMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveNetworkError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeNetwork
	ev := &eventproto.NetworkErrorMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveInternalSdkError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeSDKInternal
	ev := &eventproto.InternalSdkErrorMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	return s.saveErrorCount(event, projectID, env, errorType, ev.ApiId, ev.Labels)
}

func (s *grpcGatewayService) saveUnknownError(event *eventproto.MetricsEvent, projectID, env string) error {
	errorType := ErrorTypeUnknown
	ev := &eventproto.UnknownErrorMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
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
