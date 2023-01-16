package api

type EventType int

type metricsDetailEventType int

const (
	GoalEventType EventType = iota + 1 // eventType starts from 1 for validation.
	// Do NOT remove the goalBatchEventType because the go-server-sdk depends on the same order
	// https://github.com/ca-dp/bucketeer-go-server-sdk/blob/master/pkg/bucketeer/api/rest.go#L35
	GoalBatchEventType // nolint:deadcode,unused,varcheck
	EvaluationEventType
	MetricsEventType
)

const (
	LatencyMetricsEventType metricsDetailEventType = iota + 1
	SizeMetricsEventType
	TimeoutErrorMetricsEventType
	InternalErrorMetricsEventType
	NetworkErrorMetricsEventType
	InternalSdkErrorMetricsEventType
)
