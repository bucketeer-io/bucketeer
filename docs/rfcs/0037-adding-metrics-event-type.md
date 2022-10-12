# Plan for adding metrics event type

## Summary

Currently, Bucketeer has [several metrics event type](https://github.com/bucketeer-io/bucketeer/blob/main/proto/event/client/event.proto). However, there are some problems as follows:

* It doesn't report each event's internal error and timeout error. Thus, we can't check which event has a problem when error occurs.
* It doesn't report any metics about latency and size except for register events.

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

### 2. How far should we abstruct event type?

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
    Event string
    Labels   map[string]string
	MetaData map[string]string   
}
```

### 3. How far should we abstruct error event type?

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
    Tag string
    MetaData map[string]string 
}
```
