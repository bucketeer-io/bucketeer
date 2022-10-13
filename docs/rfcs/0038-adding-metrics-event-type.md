# Plan for adding metrics event type

## Summary

Currently, Bucketeer has [several metrics event type](https://github.com/bucketeer-io/bucketeer/blob/main/proto/event/client/event.proto). However, there are some problems as follows:

* It doesn't report each event's internal error and timeout error. Thus, we can't check which event has a problem when error occurs.
* It doesn't report any metrics about latency and size except for register events.

Therefore, we need to fix the above problems.

## TODO

To add metrics event type, we have the following TODOs:

* Support added metrics event type in Bucketeer server. Especially, we need to modify pkg/gateway/api/api.go.
* Support added metrics event type in all SDKs.

I plan to add the following metrics event type.

* LatencyMetricsEvent for each events
* SizeMetricsEvent for each events
* InternalError for each events
* TimeoutError for each events

## Controversial topics

### 1. Do we need Labels field?

We can divide Labels field into some fields.

#### Plan A

```go
type LatencyMetricsEvent struct {
    Event string
	Labels   map[string]string
	Duration time.Duration    
}

type SizeMetricsEvent struct {
    Event string
	Labels   map[string]string
	SizeByte int32
}
...

```

#### Plan B

```go
type LatencyMetricsEvent struct {
    Event string
	Tag string
    Status string
	Duration time.Duration    
}

type SizeMetricsEvent struct {
    Event string
	Tag string
    Status string
	SizeByte int32
}

...

```

### 2. How far should we abstract event type?

#### Plan A

```go
type LatencyMetricsEvent struct {
    Event string
	Labels   map[string]string
	Duration time.Duration    
}

type SizeMetricsEvent struct {
    Event string
	Labels   map[string]string
	SizeByte int32
}
```

#### Plan B

```go
type MetricsEvent struct {
    Event string // GetEvaluations, GetEvaluation, ...
    Type string // Size or Latency, ...
	Labels map[string]string   
}
```

### 3. How far should we abstract error event type?

#### Plan A

```go
type InternalErrorCountMetricsEvent struct {
    Event string
    Tag string
}

type TimeoutErrorCountMetricsEvent struct {
    Event string
    Tag string
}
```

#### Plan B

```go
type ErrorMetricsEvent struct {
    Event string
    Error string
    Labels map[string]string 
}
```

## Conclustion

## About 1

We decided A because of following reasons:

* We may remove Tag field and Status field.

## About 2 and 3

We decided A.
Personally, I prefer to B because we can simplify the logic around JSON. However, we can't decide it because we released new version in go server SDK and Android SDK in the previous format.

## Summary

As a summary, we decided the following format:

```go
type LatencyMetricsEvent struct {
    Event string
	Metadata   map[string]string
	Duration time.Duration    
}

type SizeMetricsEvent struct {
    Event string
	Metadata   map[string]string
	SizeByte int32
}

type InternalErrorCountMetricsEvent struct {
    Event string
    Metadata   map[string]string
}

type TimeoutErrorCountMetricsEvent struct {
    Event string
    Metadata   map[string]string
}
```

## Other

We change Labels field to Metadata field.
