// Copyright 2024 The Bucketeer Authors.
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

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/durationpb"

	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
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
				size, err := ptypes.MarshalAny(&eventproto.SizeMetricsEvent{
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
				size, err := ptypes.MarshalAny(&eventproto.SizeMetricsEvent{
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
				latency, err := ptypes.MarshalAny(&eventproto.LatencyMetricsEvent{
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
				latency, err := ptypes.MarshalAny(&eventproto.LatencyMetricsEvent{
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
				internalSDKErr, err := ptypes.MarshalAny(&eventproto.InternalSdkErrorMetricsEvent{
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
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
