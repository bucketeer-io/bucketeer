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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
)

func TestSaveMetrics(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	projectID := "project0"
	ns := "ns0"
	patterns := []struct {
		desc        string
		inputEvent  func() *eventproto.MetricsEvent
		expectedErr error
	}{
		{
			desc: "error: SizeMetricsEvent MetricsSaveErrUnknownApiId",
			inputEvent: func() *eventproto.MetricsEvent {
				size, err := anypb.New(&eventproto.SizeMetricsEvent{
					ApiId:    eventproto.ApiId_UNKNOWN_API,
					Labels:   map[string]string{"tag": "iOS"},
					SizeByte: 100,
				})
				if err != nil {
					t.Fatal(err)
				}
				return &eventproto.MetricsEvent{
					Timestamp:  time.Now().Unix(),
					Event:      size,
					SdkVersion: "v0.0.1-unit-test",
					SourceId:   eventproto.SourceId_IOS,
				}
			},
			expectedErr: MetricsSaveErrUnknownApiId,
		},
		{
			desc: "success: SizeMetricsEvent",
			inputEvent: func() *eventproto.MetricsEvent {
				size, err := anypb.New(&eventproto.SizeMetricsEvent{
					ApiId:    eventproto.ApiId_GET_EVALUATIONS,
					Labels:   map[string]string{"tag": "iOS"},
					SizeByte: 100,
				})
				if err != nil {
					t.Fatal(err)
				}
				return &eventproto.MetricsEvent{
					Timestamp:  time.Now().Unix(),
					Event:      size,
					SdkVersion: "v0.0.1-unit-test",
					SourceId:   eventproto.SourceId_IOS,
				}
			},
			expectedErr: nil,
		},
		{
			desc: "error: LatencyMetricsEvent MetricsSaveErrInvalidDuration",
			inputEvent: func() *eventproto.MetricsEvent {
				latency, err := anypb.New(&eventproto.LatencyMetricsEvent{
					ApiId:  eventproto.ApiId_GET_EVALUATIONS,
					Labels: map[string]string{"tag": "iOS"},
				})
				if err != nil {
					t.Fatal(err)
				}
				return &eventproto.MetricsEvent{
					Timestamp:  time.Now().Unix(),
					Event:      latency,
					SdkVersion: "v0.0.1-unit-test",
					SourceId:   eventproto.SourceId_IOS,
				}
			},
			expectedErr: MetricsSaveErrInvalidDuration,
		},
		{
			desc: "success: LatencyMetricsEvent",
			inputEvent: func() *eventproto.MetricsEvent {
				latency, err := anypb.New(&eventproto.LatencyMetricsEvent{
					ApiId:    eventproto.ApiId_GET_EVALUATIONS,
					Labels:   map[string]string{"tag": "iOS"},
					Duration: durationpb.New(time.Duration(5)),
				})
				if err != nil {
					t.Fatal(err)
				}
				return &eventproto.MetricsEvent{
					Timestamp:  time.Now().Unix(),
					Event:      latency,
					SdkVersion: "v0.0.1-unit-test",
					SourceId:   eventproto.SourceId_IOS,
				}
			},
			expectedErr: nil,
		},
		{
			desc: "error: UnknownEvent",
			inputEvent: func() *eventproto.MetricsEvent {
				return &eventproto.MetricsEvent{
					Timestamp:  time.Now().Unix(),
					SdkVersion: "v0.0.1-unit-test",
					SourceId:   eventproto.SourceId_UNKNOWN,
				}
			},
			expectedErr: MetricsSaveErrUnknownEvent,
		},

		{
			desc: "success: ErrorMetricsEvent",
			inputEvent: func() *eventproto.MetricsEvent {
				internalSDKErr, err := anypb.New(&eventproto.InternalSdkErrorMetricsEvent{
					ApiId:  eventproto.ApiId_GET_EVALUATIONS,
					Labels: map[string]string{"tag": "iOS"},
				})
				if err != nil {
					t.Fatal(err)
				}
				return &eventproto.MetricsEvent{
					Timestamp:  time.Now().Unix(),
					Event:      internalSDKErr,
					SdkVersion: "v0.0.1-unit-test",
					SourceId:   eventproto.SourceId_IOS,
				}
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			err := gs.saveMetrics(p.inputEvent(), projectID, ns)
			if p.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, p.expectedErr)
			}
		})
	}
}

func TestSaveMetricsEventsAsyncEnqueuesJob(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	gs := newGrpcGatewayServiceWithMock(t, mockController)
	gs.metricsJobs = make(chan metricsJob, 1)

	events := []*eventproto.MetricsEvent{{SdkVersion: "v0.0.1-unit-test"}}
	gs.saveMetricsEventsAsync(events, "project-id", "env-url-code")

	select {
	case job := <-gs.metricsJobs:
		assert.Equal(t, events, job.events)
		assert.Equal(t, "project-id", job.projectID)
		assert.Equal(t, "env-url-code", job.environmentUrlCode)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for metrics job")
	}
}

func TestSaveMetricsEventsAsyncSpillsToGoroutineWhenQueueIsFull(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	gs := newGrpcGatewayServiceWithMock(t, mockController)
	gs.metricsJobs = make(chan metricsJob, 1)
	// Pre-fill the queue so the next submit cannot enqueue.
	gs.metricsJobs <- metricsJob{projectID: "queued"}

	processed := make(chan string, 1)
	gs.metricsJobProcessor = func(job metricsJob) {
		processed <- job.projectID
	}

	gs.saveMetricsEventsAsync([]*eventproto.MetricsEvent{{SdkVersion: "v0.0.1-unit-test"}}, "overflow", "env-url-code")

	// The overflow batch should be processed on a one-off goroutine, not dropped.
	select {
	case projectID := <-processed:
		assert.Equal(t, "overflow", projectID)
	case <-time.After(time.Second):
		t.Fatal("overflow batch was not processed")
	}

	// The originally queued job must still be in the queue (overflow does not displace it).
	assert.Equal(t, 1, len(gs.metricsJobs))
	job := <-gs.metricsJobs
	assert.Equal(t, "queued", job.projectID)
}

func TestMetricsWorkerRecoversFromPanicAndContinues(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	gs := newGrpcGatewayServiceWithMock(t, mockController)
	processed := make(chan string, 1)
	gs.metricsJobProcessor = func(job metricsJob) {
		if job.projectID == "panic" {
			panic("poison metrics job")
		}
		processed <- job.projectID
	}
	gs.startMetricsWorkers(1, 2)
	defer gs.ShutdownMetricsPool()

	gs.saveMetricsEventsAsync([]*eventproto.MetricsEvent{{SdkVersion: "v0.0.1-unit-test"}}, "panic", "env-url-code")
	gs.saveMetricsEventsAsync([]*eventproto.MetricsEvent{{SdkVersion: "v0.0.1-unit-test"}}, "ok", "env-url-code")

	select {
	case projectID := <-processed:
		assert.Equal(t, "ok", projectID)
	case <-time.After(time.Second):
		t.Fatal("worker did not continue after recovered panic")
	}
}

func TestShutdownMetricsPoolDrainsQueuedJobs(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	gs := newGrpcGatewayServiceWithMock(t, mockController)
	processed := make(chan string, 3)
	gs.metricsJobProcessor = func(job metricsJob) {
		processed <- job.projectID
	}
	gs.startMetricsWorkers(1, 3)

	gs.saveMetricsEventsAsync([]*eventproto.MetricsEvent{{SdkVersion: "v0.0.1-unit-test"}}, "project-1", "env-url-code")
	gs.saveMetricsEventsAsync([]*eventproto.MetricsEvent{{SdkVersion: "v0.0.1-unit-test"}}, "project-2", "env-url-code")
	gs.saveMetricsEventsAsync([]*eventproto.MetricsEvent{{SdkVersion: "v0.0.1-unit-test"}}, "project-3", "env-url-code")
	gs.ShutdownMetricsPool()

	close(processed)
	assert.ElementsMatch(t, []string{"project-1", "project-2", "project-3"}, channelValues(processed))
}

func TestShutdownMetricsPoolWaitsForOverflowJob(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	gs := newGrpcGatewayServiceWithMock(t, mockController)
	gs.metricsJobs = make(chan metricsJob, 1)
	gs.metricsJobs <- metricsJob{projectID: "queued"}

	overflowStarted := make(chan struct{})
	releaseOverflow := make(chan struct{})
	processed := make(chan string, 1)
	gs.metricsJobProcessor = func(job metricsJob) {
		close(overflowStarted)
		<-releaseOverflow
		processed <- job.projectID
	}

	gs.saveMetricsEventsAsync(
		[]*eventproto.MetricsEvent{{SdkVersion: "v0.0.1-unit-test"}},
		"overflow",
		"env-url-code",
	)

	select {
	case <-overflowStarted:
	case <-time.After(time.Second):
		t.Fatal("overflow job did not start")
	}

	shutdownDone := make(chan struct{})
	go func() {
		gs.ShutdownMetricsPool()
		close(shutdownDone)
	}()

	select {
	case <-shutdownDone:
		t.Fatal("shutdown returned before overflow job completed")
	case <-time.After(50 * time.Millisecond):
	}

	close(releaseOverflow)

	select {
	case <-shutdownDone:
	case <-time.After(time.Second):
		t.Fatal("shutdown did not complete after overflow job finished")
	}
	assert.Equal(t, "overflow", <-processed)
}

func channelValues(ch <-chan string) []string {
	var values []string
	for value := range ch {
		values = append(values, value)
	}
	return values
}
