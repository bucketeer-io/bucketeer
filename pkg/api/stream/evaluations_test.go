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

package stream

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

func TestValidateStreamEvaluationsRequest(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		req         *gatewayproto.StreamEvaluationsRequest
		expectedErr error
	}{
		{
			desc:        "tag missing",
			req:         &gatewayproto.StreamEvaluationsRequest{User: &userproto.User{Id: "u"}},
			expectedErr: errTagRequired,
		},
		{
			desc:        "user missing",
			req:         &gatewayproto.StreamEvaluationsRequest{Tag: "t"},
			expectedErr: errUserRequired,
		},
		{
			desc:        "user id missing",
			req:         &gatewayproto.StreamEvaluationsRequest{Tag: "t", User: &userproto.User{}},
			expectedErr: errUserIDRequired,
		},
		{
			desc:        "valid",
			req:         &gatewayproto.StreamEvaluationsRequest{Tag: "t", User: &userproto.User{Id: "u"}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expectedErr, validateStreamEvaluationsRequest(p.req))
		})
	}
}

type stubFlusher struct{}

func (stubFlusher) Flush() {}

func TestSendSSEEvent(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	msg := &gatewayproto.StreamEvaluationsEvent{
		Evaluations: &featureproto.UserEvaluations{
			Id:          "test",
			Evaluations: []*featureproto.Evaluation{},
		},
	}
	err := sendSSEEvent(&buf, stubFlusher{}, "put", msg)
	require.NoError(t, err)
	data, err := sseMarshalOpts.Marshal(msg)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("event: put\ndata: %s\n\n", data), buf.String())
}

func TestSendHeartbeat(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	err := sendHeartbeat(&buf, stubFlusher{})
	require.NoError(t, err)
	assert.Equal(t, ":\n\n", buf.String())
}

func TestSendErrorEvent(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	sendErrorEvent(&buf, stubFlusher{}, gatewayproto.StreamErrorEvent_INTERNAL, "something broke")
	data, err := sseMarshalOpts.Marshal(&gatewayproto.StreamErrorEvent{
		Code:    gatewayproto.StreamErrorEvent_INTERNAL,
		Message: "something broke",
	})
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("event: error\ndata: %s\n\n", data), buf.String())
}

func TestInitialPut(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc      string
		evals     *featureproto.UserEvaluations
		evalErr   error
		expectErr bool
	}{
		{
			desc: "success",
			evals: &featureproto.UserEvaluations{
				Id:          "eval-1",
				Evaluations: []*featureproto.Evaluation{},
			},
		},
		{
			desc:      "evaluate error",
			evalErr:   errors.New("eval failed"),
			expectErr: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			h := &EvaluationsHandler{
				evaluate: func(_ context.Context, _ *userproto.User, _, _ string, prevUEID string, evaluatedAt int64) (string, *featureproto.UserEvaluations, error) {
					assert.Equal(t, "", prevUEID)
					assert.Equal(t, int64(0), evaluatedAt)
					return "ueid-1", p.evals, p.evalErr
				},
			}
			var buf bytes.Buffer
			ueid, evalAt, err := h.sendInitialPut(context.Background(), &buf, stubFlusher{}, &userproto.User{Id: "u"}, "env1", "tag1", "source1", "", 0)
			if p.expectErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "ueid-1", ueid)
			assert.Greater(t, evalAt, int64(0))
			lines := strings.SplitN(buf.String(), "\n", 3)
			assert.Equal(t, "event: put", lines[0])
			var got gatewayproto.StreamEvaluationsEvent
			require.NoError(t, sseUnmarshalOpts.Unmarshal([]byte(strings.TrimPrefix(lines[1], "data: ")), &got))
			assert.True(t, proto.Equal(p.evals, got.Evaluations))
			assert.Equal(t, "ueid-1", got.GetUserEvaluationsId())
		})
	}
}

func TestPatch(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		evaluatedAt int64
		evals       *featureproto.UserEvaluations
		evalErr     error
		expectErr   bool
		expectSend  bool
	}{
		{
			desc:        "diff: sends patch event",
			evaluatedAt: 1000,
			evals: &featureproto.UserEvaluations{
				Id: "eval-2",
				Evaluations: []*featureproto.Evaluation{
					{FeatureId: "f1"},
				},
			},
			expectSend: true,
		},
		{
			desc:        "none: skips sending",
			evaluatedAt: 1000,
			evals: &featureproto.UserEvaluations{
				Id:          "eval-2",
				Evaluations: []*featureproto.Evaluation{},
			},
			expectSend: false,
		},
		{
			desc:        "evaluate error",
			evaluatedAt: 1000,
			evalErr:     errors.New("eval failed"),
			expectErr:   true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			h := &EvaluationsHandler{
				evaluate: func(_ context.Context, _ *userproto.User, _, _ string, prevUEID string, evaluatedAt int64) (string, *featureproto.UserEvaluations, error) {
					assert.Equal(t, "prev-ueid", prevUEID)
					assert.Equal(t, p.evaluatedAt, evaluatedAt)
					return "new-ueid", p.evals, p.evalErr
				},
			}
			var buf bytes.Buffer
			ueid, newEvalAt, err := h.sendPatch(context.Background(), &buf, stubFlusher{}, &userproto.User{Id: "u"}, "env1", "tag1", "source1", "prev-ueid", p.evaluatedAt)
			if p.expectErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "new-ueid", ueid)
			assert.Greater(t, newEvalAt, p.evaluatedAt)
			if !p.expectSend {
				assert.Empty(t, buf.String())
				return
			}
			lines := strings.SplitN(buf.String(), "\n", 3)
			assert.Equal(t, "event: patch", lines[0])
			var got gatewayproto.StreamEvaluationsEvent
			require.NoError(t, sseUnmarshalOpts.Unmarshal([]byte(strings.TrimPrefix(lines[1], "data: ")), &got))
			assert.True(t, proto.Equal(p.evals, got.Evaluations))
			assert.Equal(t, "new-ueid", got.GetUserEvaluationsId())
		})
	}
}

// Handle blocks in its event loop while the SSE connection is active.
// Shutdown must unblock it so httpServer.Shutdown is not delayed.
func TestHandleExitsOnDispatcherShutdown(t *testing.T) {
	t.Parallel()
	d := NewDispatcher(zap.NewNop())
	handleReturned := startBlockedStreamHandle(t, d)

	d.Shutdown()

	select {
	case <-handleReturned:
		// Success: Shutdown unblocked the handler.
	case <-time.After(time.Second):
		t.Fatal("Handle is still blocked in its event loop after Dispatcher.Shutdown")
	}
}

// startBlockedStreamHandle runs Handle for one SSE connection in a goroutine and waits
// until the connection is registered. The returned channel is closed when Handle returns.
func startBlockedStreamHandle(t *testing.T, d *Dispatcher) <-chan struct{} {
	t.Helper()
	h := NewEvaluationsHandler(
		d,
		time.Hour, // ensure the heartbeat ticker never fires during the test
		func(_ context.Context, _ *userproto.User, _, _ string, _ string, _ int64) (string, *featureproto.UserEvaluations, error) {
			return "ueid-1", &featureproto.UserEvaluations{
				Id:          "ueid-1",
				Evaluations: []*featureproto.Evaluation{},
			}, nil
		},
		func(_ context.Context, _ *http.Request) (*accountproto.EnvironmentAPIKey, error) {
			return &accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "env-1"},
			}, nil
		},
		prometheus.NewCounterVec(prometheus.CounterOpts{Name: "test_request_total"},
			[]string{"organization_id", "project_id", "project_url_code",
				"environment_id", "environment_url_code", "method", "source_id"}),
		zap.NewNop(),
	)

	body, err := protojson.Marshal(&gatewayproto.StreamEvaluationsRequest{
		Tag:  "tag-A",
		User: &userproto.User{Id: "u-1"},
	})
	require.NoError(t, err)
	httpReq := httptest.NewRequest(http.MethodPost, "/stream_evaluations", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handleReturned := make(chan struct{})
	go func() {
		defer close(handleReturned)
		h.Handle(rec, httpReq)
	}()

	require.Eventually(t, func() bool {
		return len(snapshotConnCounts(d)) > 0
	}, time.Second, 10*time.Millisecond)

	return handleReturned
}
