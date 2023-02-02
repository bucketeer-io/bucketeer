package api

type EventType int

const (
	GoalEventType EventType = iota + 1 // eventType starts from 1 for validation.
	// Do NOT remove the goalBatchEventType because the go-server-sdk depends on the same order
	// https://github.com/ca-dp/bucketeer-go-server-sdk/blob/master/pkg/bucketeer/api/rest.go#L35
	GoalBatchEventType // nolint:deadcode,unused,varcheck
	EvaluationEventType
	MetricsEventType
)

type metricsDetailEventType int

const (
	latencyMetricsEventType metricsDetailEventType = iota + 1
	sizeMetricsEventType
	timeoutErrorMetricsEventType
	internalErrorMetricsEventType
	networkErrorMetricsEventType
	internalSdkErrorMetricsEventType
	badRequestErrorMetricsEventType
	unauthorizedErrorMetricsEventType
	forbiddenErrorMetricsEventType
	notFoundErrorMetricsEventType
	clientClosedRequestErrorMetricsEventType
	internalServerErrorMetricsEventType
	serviceUnavailableErrorMetricsEventType
	unknownErrorMetricsEventType
)

var metricsEventTypeValueMap = map[metricsDetailEventType]string{
	latencyMetricsEventType:                  "LatencyMetricsEvent",
	sizeMetricsEventType:                     "SizeMetricsEvent",
	timeoutErrorMetricsEventType:             "TimeoutErrorMetricsEvent",
	internalErrorMetricsEventType:            "InternalErrorMetricsEvent",
	networkErrorMetricsEventType:             "NetworkErrorMetricsEvent",
	internalSdkErrorMetricsEventType:         "InternalSdkErrorMetricsEvent",
	badRequestErrorMetricsEventType:          "badRequestErrorMetricsEvent",
	unauthorizedErrorMetricsEventType:        "unauthorizedErrorMetricsEvent",
	forbiddenErrorMetricsEventType:           "forbiddenErrorMetricsEvent",
	notFoundErrorMetricsEventType:            "notFoundErrorMetricsEvent",
	clientClosedRequestErrorMetricsEventType: "clientClosedRequestErrorMetricsEvent",
	internalServerErrorMetricsEventType:      "internalServerErrorMetricsEvent",
	serviceUnavailableErrorMetricsEventType:  "serviceUnavailableErrorMetricsEvent",
	unknownErrorMetricsEventType:             "unknownErrorMetricsEvent",
}

func (m metricsDetailEventType) String() string {
	return metricsEventTypeValueMap[m]
}
