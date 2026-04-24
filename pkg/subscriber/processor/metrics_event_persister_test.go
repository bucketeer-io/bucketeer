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

package processor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/mock"
	clientevent "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
)

func TestUnmarshalMessage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*testing.T) (*clientevent.MetricsEvent, *puller.Message)
		expectedErr bool
	}{
		{
			desc: "getEvaluationLatencyMetricsEvent: success",
			setup: func(t *testing.T) (*clientevent.MetricsEvent, *puller.Message) {
				e, err := anypb.New(&clientevent.GetEvaluationLatencyMetricsEvent{
					Labels:   map[string]string{"tag": "test", "status": "success"},
					Duration: &durationpb.Duration{Seconds: time.Now().Unix()},
				})
				require.NoError(t, err)
				me := &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
				any, err := anypb.New(me)
				assert.NoError(t, err)
				event := &clientevent.Event{Event: any}
				data, err := proto.Marshal(event)
				assert.NoError(t, err)
				return me, &puller.Message{Data: data}
			},
			expectedErr: false,
		},
		{
			desc: "getEvaluationLatencyMetricsEvent: invalid message data",
			setup: func(t *testing.T) (*clientevent.MetricsEvent, *puller.Message) {
				me := &clientevent.GoalEvent{}
				data, err := proto.Marshal(me)
				assert.NoError(t, err)
				return nil, &puller.Message{Data: data}
			},
			expectedErr: true,
		},
		{
			desc: "getEvaluationLatencyMetricsEvent: invalid metrics event",
			setup: func(t *testing.T) (*clientevent.MetricsEvent, *puller.Message) {
				me := &clientevent.GoalEvent{}
				any, err := anypb.New(me)
				assert.NoError(t, err)
				event := &clientevent.Event{Event: any}
				data, err := proto.Marshal(event)
				assert.NoError(t, err)
				return nil, &puller.Message{Data: data}
			},
			expectedErr: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			pst := newMetricsEventPersister(t, mockController)
			expected, input := p.setup(t)
			e, err := pst.unmarshalMessage(input)
			assert.Equal(t, p.expectedErr, err != nil)
			if !p.expectedErr {
				assert.Equal(t, expected.Timestamp, e.Timestamp)
			}
		})
	}
}

func newMetricsEventPersister(t *testing.T, mockController *gomock.Controller) *metricsEventPersister {
	t.Helper()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &metricsEventPersister{
		storage: storagemock.NewMockStorage(mockController),
		logger:  logger.Named("experiment-cacher"),
	}
}

func TestSaveMetrics(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*testing.T, *metricsEventPersister) *clientevent.MetricsEvent
		expectedErr error
	}{
		{
			desc: "error: unknown event",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				e, err := anypb.New(&clientevent.GoalEvent{})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: ErrUnknownEvent,
		},
		{
			desc: "getEvaluationLatencyMetricsEvent: error: invalid duration",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				e, err := anypb.New(&clientevent.GetEvaluationLatencyMetricsEvent{
					Labels:   map[string]string{"tag": "test", "state": "FULL"},
					Duration: nil,
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: ErrInvalidDuration,
		},
		{
			desc: "LatencyMetricsEvent: error: invalid duration",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				e, err := anypb.New(&clientevent.LatencyMetricsEvent{
					ApiId:    clientevent.ApiId_GET_EVALUATION,
					Labels:   map[string]string{"tag": "test", "state": "FULL"},
					Duration: nil,
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: ErrInvalidDuration,
		},
		{
			desc: "LatencyMetricsEvent: error: unknown api id",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				e, err := anypb.New(&clientevent.LatencyMetricsEvent{
					Labels:   map[string]string{"tag": "test", "state": "FULL"},
					Duration: &durationpb.Duration{Seconds: time.Now().Unix()},
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: ErrUnknownApiId,
		},
		{
			desc: "SizeMetricsEvent: error: unknown api id",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				e, err := anypb.New(&clientevent.SizeMetricsEvent{
					Labels: map[string]string{"tag": "test", "state": "FULL"},
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: ErrUnknownApiId,
		},
		{
			desc: "TimeoutErrorMetricsEvent: error: unknown api id",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				e, err := anypb.New(&clientevent.TimeoutErrorMetricsEvent{
					Labels: map[string]string{"tag": "test"},
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: ErrUnknownApiId,
		},
		{
			desc: "InternalErrorMetricsEvent: error: unknown api id",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				e, err := anypb.New(&clientevent.InternalErrorMetricsEvent{
					Labels: map[string]string{"tag": "test"},
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: ErrUnknownApiId,
		},
		{
			desc: "NetworkErrorMetrics: error: unknown api id",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				e, err := anypb.New(&clientevent.NetworkErrorMetricsEvent{
					Labels: map[string]string{"tag": "test"},
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: ErrUnknownApiId,
		},
		{
			desc: "InternalSdkErrorMetricsEvent: error: unknown api id",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				e, err := anypb.New(&clientevent.InternalSdkErrorMetricsEvent{
					Labels: map[string]string{"tag": "test"},
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: ErrUnknownApiId,
		},
		{
			desc: "getEvaluationLatencyMetricsEvent: success",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				pst.storage.(*storagemock.MockStorage).EXPECT().SaveGetEvaluationLatencyMetricsEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return().Times(1)
				e, err := anypb.New(&clientevent.GetEvaluationLatencyMetricsEvent{
					Labels:   map[string]string{"tag": "test", "state": "FULL"},
					Duration: &durationpb.Duration{Seconds: time.Now().Unix()},
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: nil,
		},
		{
			desc: "getEvaluationSizeMetricsEvent: success",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				pst.storage.(*storagemock.MockStorage).EXPECT().SaveGetEvaluationSizeMetricsEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return().Times(1)
				e, err := anypb.New(&clientevent.GetEvaluationSizeMetricsEvent{
					Labels:   map[string]string{"tag": "test", "state": "FULL"},
					SizeByte: 100,
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: nil,
		},
		{
			desc: "TimeoutErrorCountMetricsEvent: success",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				pst.storage.(*storagemock.MockStorage).EXPECT().SaveTimeoutErrorCountMetricsEvent(gomock.Any()).Return().Times(1)
				e, err := anypb.New(&clientevent.TimeoutErrorCountMetricsEvent{
					Tag: "test",
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: nil,
		},
		{
			desc: "InternalErrorCountMetricsEvent: success",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				pst.storage.(*storagemock.MockStorage).EXPECT().SaveInternalErrorCountMetricsEvent(gomock.Any()).Return().Times(1)
				e, err := anypb.New(&clientevent.InternalErrorCountMetricsEvent{
					Tag: "test",
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: nil,
		},
		{
			desc: "LatencyMetricsEvent: success",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				pst.storage.(*storagemock.MockStorage).EXPECT().SaveLatencyMetricsEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return().Times(1)
				e, err := anypb.New(&clientevent.LatencyMetricsEvent{
					ApiId:    clientevent.ApiId_GET_EVALUATIONS,
					Labels:   map[string]string{"tag": "test", "state": "FULL"},
					Duration: &durationpb.Duration{Seconds: time.Now().Unix()},
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: nil,
		},
		{
			desc: "SizeMetricsEvent: error: unknown api id",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				pst.storage.(*storagemock.MockStorage).EXPECT().SaveSizeMetricsEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return().Times(1)
				e, err := anypb.New(&clientevent.SizeMetricsEvent{
					ApiId:  clientevent.ApiId_GET_EVALUATIONS,
					Labels: map[string]string{"tag": "test", "state": "FULL"},
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: nil,
		},
		{
			desc: "TimeoutErrorMetricsEvent: error: unknown api id",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				pst.storage.(*storagemock.MockStorage).EXPECT().SaveTimeoutErrorMetricsEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return().Times(1)
				e, err := anypb.New(&clientevent.TimeoutErrorMetricsEvent{
					ApiId:  clientevent.ApiId_GET_EVALUATIONS,
					Labels: map[string]string{"tag": "test"},
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: nil,
		},
		{
			desc: "InternalErrorMetricsEvent: error: unknown api id",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				pst.storage.(*storagemock.MockStorage).EXPECT().SaveInternalErrorMetricsEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return().Times(1)
				e, err := anypb.New(&clientevent.InternalErrorMetricsEvent{
					ApiId:  clientevent.ApiId_GET_EVALUATIONS,
					Labels: map[string]string{"tag": "test"},
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: nil,
		},
		{
			desc: "NetworkErrorMetrics: error: unknown api id",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				pst.storage.(*storagemock.MockStorage).EXPECT().SaveNetworkErrorMetricsEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return().Times(1)
				e, err := anypb.New(&clientevent.NetworkErrorMetricsEvent{
					ApiId:  clientevent.ApiId_GET_EVALUATIONS,
					Labels: map[string]string{"tag": "test"},
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: nil,
		},
		{
			desc: "InternalSdkErrorMetricsEvent: error: unknown api id",
			setup: func(t *testing.T, pst *metricsEventPersister) *clientevent.MetricsEvent {
				pst.storage.(*storagemock.MockStorage).EXPECT().SaveInternalSdkErrorMetricsEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return().Times(1)
				e, err := anypb.New(&clientevent.InternalSdkErrorMetricsEvent{
					ApiId:  clientevent.ApiId_GET_EVALUATIONS,
					Labels: map[string]string{"tag": "test"},
				})
				require.NoError(t, err)
				return &clientevent.MetricsEvent{
					Timestamp: time.Now().Unix(),
					Event:     e,
				}
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			pst := newMetricsEventPersister(t, mockController)
			input := p.setup(t, pst)
			err := pst.saveMetrics(input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
