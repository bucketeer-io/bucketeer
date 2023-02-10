package api

import (
	"errors"

	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"

	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
)

const (
	ErrorTypeBadRequest          = "BadRequest"
	ErrorTypeUnauthenticated     = "Unauthenticated"
	ErrorTypeForbidden           = "Forbidden"
	ErrorTypeNotFound            = "NotFound"
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
	badRequestErrorMetricsEventP          = &eventproto.BadRequestErrorMetricsEvent{}
	unauthorizedErrorMetricsEventP        = &eventproto.UnauthorizedErrorMetricsEvent{}
	forbiddenErrorMetricsEventP           = &eventproto.ForbiddenErrorMetricsEvent{}
	notFoundErrorMetricsEventP            = &eventproto.NotFoundErrorMetricsEvent{}
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
	metricsEvents []*eventproto.MetricsEvent, environmentNamespace string,
) {
	// TODO: using buffered channel to reduce the number of go routines
	go func() {
		for i := range metricsEvents {
			if err := s.saveMetrics(metricsEvents[i], environmentNamespace); err != nil {
				s.logger.Error("Failed to store metrics event to prometheus client", zap.Error(err))
				eventCounter.WithLabelValues(callerGatewayService, typeMetrics, codeNonRepeatableError).Inc()
			} else {
				eventCounter.WithLabelValues(callerGatewayService, typeMetrics, codeOK).Inc()
			}
		}
	}()
}

func (s *grpcGatewayService) saveMetrics(event *eventproto.MetricsEvent, environmentNamespace string) error {
	if ptypes.Is(event.Event, getEvaluationLatencyMetricsEventP) {
		return s.saveGetEvaluationLatencyMetricsEvent(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, getEvaluationSizeMetricsEventP) {
		return s.saveGetEvaluationSizeMetricsEvent(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, timeoutErrorCountMetricsEventP) {
		return s.saveTimeoutErrorCountMetricsEvent(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, internalErrorCountMetricsEventP) {
		return s.saveInternalErrorCountMetricsEvent(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, latencyMetricsEventP) {
		return s.saveLatencyMetricsEvent(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, sizeMetricsEventP) {
		return s.saveSizeMetricsEvent(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, badRequestErrorMetricsEventP) {
		return s.saveBadRequestError(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, unauthorizedErrorMetricsEventP) {
		return s.saveUnauthorizedError(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, forbiddenErrorMetricsEventP) {
		return s.saveForbiddenError(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, notFoundErrorMetricsEventP) {
		return s.saveNotFoundError(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, clientClosedRequestErrorMetricsEventP) {
		return s.saveClientClosedRequestError(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, internalServerErrorMetricsEventP) {
		return s.saveInternalServerError(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, serviceUnavailableErrorMetricsEventP) {
		return s.saveServiceUnavailableError(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, timeoutErrorMetricsEventP) {
		return s.saveTimeoutError(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, internalErrorMetricsEventP) {
		return s.saveInternalError(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, networkErrorMetricsEventP) {
		return s.saveNetworkError(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, internalSdkErrorMetricsEventP) {
		return s.saveInternalSdkError(event, environmentNamespace)
	}
	if ptypes.Is(event.Event, unknownErrorMetricsEventP) {
		return s.saveUnknownError(event, environmentNamespace)
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

func (s *grpcGatewayService) saveLatencyMetricsEvent(event *eventproto.MetricsEvent, env string) error {
	ev := &eventproto.LatencyMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	if ev.Duration == nil {
		return MetricsSaveErrInvalidDuration
	}
	if ev.ApiId == eventproto.ApiId_UNKNOWN_API {
		return MetricsSaveErrUnknownApiId
	}
	var tag string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
	}
	dur, err := ptypes.Duration(ev.Duration)
	if err != nil {
		return MetricsSaveErrInvalidDuration
	}
	sdkLatencyHistogram.WithLabelValues(
		env,
		tag,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Observe(dur.Seconds())
	return nil
}

func (s *grpcGatewayService) saveSizeMetricsEvent(event *eventproto.MetricsEvent, env string) error {
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
		env,
		tag,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Observe(float64(ev.SizeByte))
	return nil
}

func (s *grpcGatewayService) saveBadRequestError(event *eventproto.MetricsEvent, env string) error {
	errorType := ErrorTypeBadRequest
	ev := &eventproto.BadRequestErrorMetricsEvent{}
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
	sdkErrorCounter.WithLabelValues(
		env,
		tag,
		errorType,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Inc()
	return nil
}

func (s *grpcGatewayService) saveUnauthorizedError(event *eventproto.MetricsEvent, env string) error {
	errorType := ErrorTypeUnauthenticated
	ev := &eventproto.UnauthorizedErrorMetricsEvent{}
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
	sdkErrorCounter.WithLabelValues(
		env,
		tag,
		errorType,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Inc()
	return nil
}

func (s *grpcGatewayService) saveForbiddenError(event *eventproto.MetricsEvent, env string) error {
	errorType := ErrorTypeForbidden
	ev := &eventproto.ForbiddenErrorMetricsEvent{}
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
	sdkErrorCounter.WithLabelValues(
		env,
		tag,
		errorType,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Inc()
	return nil
}

func (s *grpcGatewayService) saveNotFoundError(event *eventproto.MetricsEvent, env string) error {
	errorType := ErrorTypeNotFound
	ev := &eventproto.NotFoundErrorMetricsEvent{}
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
	sdkErrorCounter.WithLabelValues(
		env,
		tag,
		errorType,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Inc()
	return nil
}

func (s *grpcGatewayService) saveClientClosedRequestError(event *eventproto.MetricsEvent, env string) error {
	errorType := ErrorTypeClientClosedRequest
	ev := &eventproto.ClientClosedRequestErrorMetricsEvent{}
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
	sdkErrorCounter.WithLabelValues(
		env,
		tag,
		errorType,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Inc()
	return nil
}

func (s *grpcGatewayService) saveInternalServerError(event *eventproto.MetricsEvent, env string) error {
	errorType := ErrorTypeInternalServerError
	ev := &eventproto.InternalServerErrorMetricsEvent{}
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
	sdkErrorCounter.WithLabelValues(
		env,
		tag,
		errorType,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Inc()
	return nil
}

func (s *grpcGatewayService) saveServiceUnavailableError(event *eventproto.MetricsEvent, env string) error {
	errorType := ErrorTypeServiceUnavailable
	ev := &eventproto.ServiceUnavailableErrorMetricsEvent{}
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
	sdkErrorCounter.WithLabelValues(
		env,
		tag,
		errorType,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Inc()
	return nil
}

func (s *grpcGatewayService) saveTimeoutError(event *eventproto.MetricsEvent, env string) error {
	errorType := ErrorTypeTimeout
	ev := &eventproto.TimeoutErrorMetricsEvent{}
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
	sdkErrorCounter.WithLabelValues(
		env,
		tag,
		errorType,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Inc()
	return nil
}

func (s *grpcGatewayService) saveInternalError(event *eventproto.MetricsEvent, env string) error {
	errorType := ErrorTypeInternal
	ev := &eventproto.InternalErrorMetricsEvent{}
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
	sdkErrorCounter.WithLabelValues(
		env,
		tag,
		errorType,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Inc()
	return nil
}

func (s *grpcGatewayService) saveNetworkError(event *eventproto.MetricsEvent, env string) error {
	errorType := ErrorTypeNetwork
	ev := &eventproto.NetworkErrorMetricsEvent{}
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
	sdkErrorCounter.WithLabelValues(
		env,
		tag,
		errorType,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Inc()
	return nil
}

func (s *grpcGatewayService) saveInternalSdkError(event *eventproto.MetricsEvent, env string) error {
	errorType := ErrorTypeSDKInternal
	ev := &eventproto.InternalSdkErrorMetricsEvent{}
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
	sdkErrorCounter.WithLabelValues(
		env,
		tag,
		errorType,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Inc()
	return nil
}

func (s *grpcGatewayService) saveUnknownError(event *eventproto.MetricsEvent, env string) error {
	errorType := ErrorTypeUnknown
	ev := &eventproto.UnknownErrorMetricsEvent{}
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
	sdkErrorCounter.WithLabelValues(
		env,
		tag,
		errorType,
		ev.ApiId.String(),
		event.SdkVersion,
		event.SourceId.String(),
	).Inc()
	return nil
}
